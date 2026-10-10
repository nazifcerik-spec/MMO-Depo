import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:otter/core/error_reporter.dart';

void main() {
  test('reports once per distinct error, caps per run, never throws', () async {
    final sent = <Map<String, dynamic>>[];
    var fail = false;
    final reporter = ErrorReporter(
      baseUrl: 'http://api',
      client: MockClient((req) async {
        expect(req.url.toString(), 'http://api/v1/client/errors');
        if (fail) throw http.ClientException('offline');
        sent.add(jsonDecode(req.body) as Map<String, dynamic>);
        return http.Response('', 202);
      }),
    );

    final stack = StackTrace.current;
    await reporter.report(StateError('boom'), stack, kind: 'flutter');
    await reporter.report(StateError('boom'), stack, kind: 'flutter');
    expect(sent, hasLength(1), reason: 'duplicate dropped');
    expect(sent.single['kind'], 'flutter');
    expect(sent.single['message'], 'Bad state: boom');
    expect(sent.single['stack'], contains('error_reporter_test.dart'));
    expect(sent.single['platform'], isNotEmpty);

    await reporter.report(StateError('x' * 5000), null, kind: 'platform');
    expect((sent.last['message'] as String).length, 1000);

    fail = true;
    await reporter.report(StateError('offline'), null, kind: 'platform'); // swallowed
    fail = false;
    for (var i = 0; i < 20; i++) {
      await reporter.report(StateError('e$i'), null, kind: 'platform');
    }
    expect(sent, hasLength(ErrorReporter.maxPerRun - 1), reason: 'the failed one still counts toward the cap');
  });

  test('install routes framework and async errors to the collector', () async {
    TestWidgetsFlutterBinding.ensureInitialized();
    final kinds = <String>[];
    final reporter = ErrorReporter(
      baseUrl: 'http://api',
      client: MockClient((req) async {
        kinds.add((jsonDecode(req.body) as Map)['kind'] as String);
        return http.Response('', 202);
      }),
    );
    final prevFlutter = FlutterError.onError;
    final prevPlatform = PlatformDispatcher.instance.onError;
    final prevPresent = FlutterError.presentError;
    addTearDown(() {
      FlutterError.onError = prevFlutter;
      PlatformDispatcher.instance.onError = prevPlatform;
      FlutterError.presentError = prevPresent;
    });
    reporter.install();
    FlutterError.presentError = (_) {};
    FlutterError.onError!(FlutterErrorDetails(exception: StateError('build'), stack: StackTrace.current));
    expect(PlatformDispatcher.instance.onError!(StateError('async'), StackTrace.current), isTrue);
    await pumpEventQueue();
    expect(kinds, ['flutter', 'platform']);
  });
}
