import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:otter/core/app_update.dart';
import 'package:otter/features/character/update_button.dart';
import 'package:otter/l10n/app_localizations.dart';

const feed = 'https://example.org/repo/releases/latest.json';
const channel = MethodChannel('mena/update-test');

/// A feed announcing [version] (APK [apk], [size] bytes) and serving the APK.
http.Client server({String version = '1.25.0', String apk = 'mena-world-1.25.0.apk', int size = 4, int status = 200}) =>
    MockClient((req) async {
      if (req.url.path.endsWith('.json')) {
        return http.Response(jsonEncode({'version': version, 'apk': apk, 'size': size}), status);
      }
      return http.Response.bytes([1, 2, 3, 4], 200);
    });

AppUpdate updater(http.Client c, {String current = '1.24.0'}) =>
    AppUpdate(feedUrl: feed, current: current, client: c, channel: channel);

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  test('compareVersions orders numerically', () {
    expect(compareVersions('1.25.0', '1.24.9'), greaterThan(0));
    expect(compareVersions('1.10.0', '1.9.0'), greaterThan(0));
    expect(compareVersions('1.2', '1.2.0'), 0);
    expect(compareVersions('', '1.0.0'), lessThan(0));
  });

  test('check finds a newer build next to the feed', () async {
    final r = await updater(server()).check();
    expect(r?.version, '1.25.0');
    expect(r?.apk.toString(), 'https://example.org/repo/releases/mena-world-1.25.0.apk');
    expect(r?.size, 4);
  });

  test('check: same or older version is no update', () async {
    expect(await updater(server(version: '1.24.0')).check(), isNull);
    expect(await updater(server(version: '1.23.5')).check(), isNull);
  });

  test('check refuses an APK from elsewhere and a failed feed', () async {
    expect(updater(server(apk: 'https://evil.example/x.apk')).check(), throwsFormatException);
    expect(updater(server(apk: 'notes.txt')).check(), throwsFormatException);
    expect(updater(server(status: 404)).check(), throwsA(isA<HttpException>()));
  });

  test('available needs an https feed and a version', () {
    expect(updater(server()).available, isTrue); // tests run as Android
    expect(AppUpdate(feedUrl: '', current: '1.0.0').available, isFalse);
    expect(AppUpdate(feedUrl: 'http://example.org/x.json', current: '1.0.0').available, isFalse);
    expect(AppUpdate(feedUrl: feed, current: '').available, isFalse);
  });

  group('install', () {
    late Directory dir;
    late List<String> calls;
    var installResult = 'started';

    setUp(() {
      dir = Directory.systemTemp.createTempSync('upd');
      calls = [];
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(channel, (c) async {
        calls.add(c.method);
        return c.method == 'dir' ? dir.path : installResult;
      });
    });
    tearDown(() => dir.deleteSync(recursive: true));

    test('downloads to update.apk and opens the installer', () async {
      installResult = 'started';
      final u = updater(server());
      final progress = <double>[];
      final ok = await u.install((await u.check())!, onProgress: progress.add);
      expect(ok, isTrue);
      expect(File('${dir.path}/update.apk').readAsBytesSync(), [1, 2, 3, 4]);
      expect(calls, ['dir', 'install']);
      expect(progress.last, 1);
    });

    test('a short download never reaches the installer', () async {
      final u = updater(server(size: 9));
      await expectLater(u.install((await u.check())!), throwsA(isA<HttpException>()));
      expect(calls, ['dir']);
    });

    testWidgets('bar button: dot when newer, tap downloads and installs', (tester) async {
      installResult = 'started';
      await tester.pumpWidget(
        MaterialApp(
          localizationsDelegates: AppLocalizations.localizationsDelegates,
          supportedLocales: AppLocalizations.supportedLocales,
          home: Scaffold(
            appBar: AppBar(actions: [UpdateButton(updater: updater(server()))]),
          ),
        ),
      );
      await tester.runAsync(() => Future<void>.delayed(const Duration(milliseconds: 50)));
      await tester.pump();
      final badge = tester.widget<Badge>(
        find.descendant(of: find.byKey(const Key('update-button')), matching: find.byType(Badge)),
      );
      expect(badge.isLabelVisible, isTrue);

      await tester.tap(find.byKey(const Key('update-button')));
      await tester.runAsync(() => Future<void>.delayed(const Duration(milliseconds: 50)));
      await tester.pumpAndSettle();
      expect(find.text('Download and install version 1.25.0 (1 MB)?'), findsOneWidget);
      await tester.tap(find.text('Update').last);
      for (var i = 0; i < 50 && calls.length < 2; i++) {
        await tester.runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
        await tester.pump();
      }
      expect(calls, ['dir', 'install']);
    });

    testWidgets('bar button: installs from unknown sources not allowed yet', (tester) async {
      installResult = 'permission';
      await tester.pumpWidget(
        MaterialApp(
          localizationsDelegates: AppLocalizations.localizationsDelegates,
          supportedLocales: AppLocalizations.supportedLocales,
          home: Scaffold(
            appBar: AppBar(actions: [UpdateButton(updater: updater(server()))]),
          ),
        ),
      );
      await tester.tap(find.byKey(const Key('update-button')));
      await tester.runAsync(() => Future<void>.delayed(const Duration(milliseconds: 50)));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Update').last);
      final perm = find.text('Allow installs from MeNa World, then tap Update again.');
      // The download is real file I/O: let it run between frames.
      for (var i = 0; i < 50 && perm.evaluate().isEmpty; i++) {
        await tester.runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
        await tester.pump();
      }
      expect(perm, findsOneWidget);
    });
  });
}
