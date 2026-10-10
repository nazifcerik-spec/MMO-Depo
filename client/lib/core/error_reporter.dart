import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

import 'config.dart';

/// Sends uncaught errors to the server's collector (DD-235). Never throws,
/// reports each distinct error once and at most [maxPerRun] per app run, so a
/// crash loop cannot flood the network.
class ErrorReporter {
  ErrorReporter({required this.baseUrl, http.Client? client}) : _http = client ?? http.Client();

  final String baseUrl;
  final http.Client _http;
  final _seen = <String>{};

  static const maxPerRun = 10;

  Future<void> report(Object error, StackTrace? stack, {required String kind}) async {
    final message = _cut(error.toString(), 1000);
    final trace = _cut((stack ?? StackTrace.empty).toString().split('\n').take(40).join('\n'), 6000);
    final key = '$kind|${message.split('\n').first}|${trace.split('\n').take(3).join('|')}';
    if (_seen.length >= maxPerRun || !_seen.add(key)) return;
    try {
      await _http
          .post(
            Uri.parse('$baseUrl/v1/client/errors'),
            headers: const {'Content-Type': 'application/json'},
            body: jsonEncode({
              'kind': kind,
              'message': message,
              'stack': trace,
              'platform': kIsWeb ? 'web' : defaultTargetPlatform.name,
              'app_version': AppConfig.appVersion,
              'flavor': AppConfig.flavor,
            }),
          )
          .timeout(const Duration(seconds: 5));
    } catch (_) {
      // Reporting is best effort; the app must never fail because of it.
    }
  }

  /// Routes framework and uncaught async errors here (keeps console output).
  void install() {
    FlutterError.onError = (details) {
      FlutterError.presentError(details);
      report(details.exception, details.stack, kind: 'flutter');
    };
    PlatformDispatcher.instance.onError = (error, stack) {
      debugPrint('Uncaught: $error\n$stack');
      report(error, stack, kind: 'platform');
      return true;
    };
  }

  static String _cut(String s, int n) => s.length <= n ? s : s.substring(0, n);
}
