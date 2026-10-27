import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;

import 'config.dart';

/// A build newer than the running one, from the update feed (DD-304).
class Release {
  const Release({required this.version, required this.apk, required this.size});

  final String version;
  final Uri apk;

  /// Bytes; 0 when the feed does not say.
  final int size;
}

/// In-app updates of the sideloaded Android build (DD-304): the feed is a
/// small JSON (`{"version": "1.25.0", "apk": "mena-world-1.25.0.apk",
/// "size": 41000000}`), the APK sits next to it. The system installer accepts
/// the download only when it is signed like the installed app.
class AppUpdate {
  AppUpdate({
    String feedUrl = AppConfig.updateFeedUrl,
    this.current = AppConfig.appVersion,
    http.Client? client,
    this.channel = const MethodChannel('mena/update'),
  }) : _feed = Uri.tryParse(feedUrl),
       _http = client ?? http.Client();

  final Uri? _feed;
  final String current;
  final http.Client _http;

  /// The Android side (MainActivity).
  final MethodChannel channel;

  /// Android builds that know their version and feed.
  bool get available =>
      !kIsWeb &&
      defaultTargetPlatform == TargetPlatform.android &&
      _feed != null &&
      _feed.isScheme('https') &&
      current.isNotEmpty;

  /// The newer release, or null when this build is the latest.
  Future<Release?> check() async {
    final feed = _feed!;
    // A fresh read: the feed is cached by its host for a few minutes.
    final res = await _http
        .get(feed.replace(queryParameters: {'t': '${DateTime.now().millisecondsSinceEpoch}'}))
        .timeout(const Duration(seconds: 15));
    if (res.statusCode != 200) throw HttpException('feed ${res.statusCode}');
    final j = jsonDecode(res.body) as Map<String, dynamic>;
    final version = j['version'] as String? ?? '';
    if (compareVersions(version, current) <= 0) return null;
    final apk = feed.resolve(j['apk'] as String? ?? '');
    // Only from where the feed lives, over https.
    if (!apk.isScheme('https') || apk.host != feed.host || !apk.path.endsWith('.apk')) {
      throw const FormatException('bad apk url');
    }
    return Release(version: version, apk: apk, size: (j['size'] as num?)?.toInt() ?? 0);
  }

  /// Downloads [r] and opens the system installer. Returns false when the
  /// player must first allow installs from this app (the settings page opens).
  Future<bool> install(Release r, {ValueChanged<double>? onProgress}) async {
    final dir = await channel.invokeMethod<String>('dir');
    final file = File('$dir/update.apk');
    final res = await _http.send(http.Request('GET', r.apk)).timeout(const Duration(seconds: 30));
    if (res.statusCode != 200) throw HttpException('apk ${res.statusCode}');
    final total = res.contentLength ?? r.size;
    final sink = file.openWrite();
    var got = 0;
    try {
      await for (final chunk in res.stream.timeout(const Duration(seconds: 30))) {
        sink.add(chunk);
        got += chunk.length;
        if (total > 0) onProgress?.call(got / total);
      }
    } finally {
      await sink.close();
    }
    if (r.size > 0 && got != r.size) throw const HttpException('apk size');
    return await channel.invokeMethod<String>('install') == 'started';
  }
}

/// "1.25.0" vs "1.24.3": negative, zero or positive. Missing parts are 0,
/// non-numeric parts sort first.
int compareVersions(String a, String b) {
  List<int> parts(String v) => [for (final p in v.split('.')) int.tryParse(p) ?? -1];
  final x = parts(a), y = parts(b);
  for (var i = 0; i < x.length || i < y.length; i++) {
    final d = (i < x.length ? x[i] : 0) - (i < y.length ? y[i] : 0);
    if (d != 0) return d;
  }
  return 0;
}
