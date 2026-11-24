import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/app/app.dart';
import 'package:otter/core/api_client.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/auth/auth_controller.dart';
import 'package:otter/features/guide/guide.dart';
import 'package:otter/shared/layout.dart';

import 'fake_backend.dart';
import 'fake_ws.dart';

void _size(WidgetTester tester, double w, double h) {
  tester.view.physicalSize = Size(w, h);
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
}

Widget _grid() => MaterialApp(
  home: Scaffold(
    body: TileGrid(
      children: [for (var i = 0; i < 5; i++) SizedBox(key: ValueKey(i), height: i.isEven ? 40 : 80)],
    ),
  ),
);

void main() {
  setUpAll(() => TutorialState.debugPending = false);

  testWidgets('tile grid: one column on phones, equal columns and row heights when wide (DD-369)', (tester) async {
    _size(tester, 375, 812);
    await tester.pumpWidget(_grid());
    expect(tester.getSize(find.byKey(const ValueKey(0))), const Size(375, 40));
    expect(tester.getTopLeft(find.byKey(const ValueKey(1))).dy, 48);

    _size(tester, 1200, 800);
    await tester.pumpWidget(_grid());
    final first = tester.getRect(find.byKey(const ValueKey(0)));
    final second = tester.getRect(find.byKey(const ValueKey(1)));
    final last = tester.getRect(find.byKey(const ValueKey(4)));
    expect(first.size, second.size);
    expect(first.top, second.top);
    expect(last.width, first.width, reason: 'an incomplete row keeps the column width');
    expect(last.top, greaterThan(first.bottom));
  });

  testWidgets('wide screens: rail instead of the bottom bar, panes in every tab (DD-369)', (tester) async {
    _size(tester, 1440, 900);
    final be = FakeBackend()
      ..characters.add({
        'id': 'char-1',
        'name': 'Grim',
        'race_code': 'orc',
        'class_code': 'rogue',
        'level': 3,
        'state': 'active',
      });
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await tester.pumpWidget(
      OtterApp(
        auth: AuthController(be.api(store), wsConnect: FakeSocketServer().connect),
        locale: const Locale('en'),
      ),
    );
    await tester.pumpAndSettle();
    await tester.tap(find.text('Grim'));
    await tester.pumpAndSettle();

    expect(find.byType(NavigationRail), findsOneWidget);
    expect(find.byType(NavigationBar), findsNothing);
    for (final tab in ['Gather', 'Bag', 'Adventure', 'Workshop', 'Character']) {
      await tester.tap(find.descendant(of: find.byType(NavigationRail), matching: find.text(tab)));
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull, reason: tab);
    }
    // Every Adventure and Workshop segment lays out too.
    for (final (tab, segments) in [
      ('Adventure', ['Dungeon', 'Travel', 'Events', 'Hunt']),
      ('Workshop', ['Shop', 'Market', 'Craft']),
    ]) {
      await tester.tap(find.descendant(of: find.byType(NavigationRail), matching: find.text(tab)));
      await tester.pumpAndSettle();
      for (final s in segments) {
        await tester.tap(find.descendant(of: find.byType(SegmentedButton<int>), matching: find.text(s)).hitTestable());
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull, reason: s);
      }
    }

    // Narrowing the window goes back to the phone layout.
    tester.view.physicalSize = const Size(500, 900);
    await tester.pumpAndSettle();
    expect(find.byType(NavigationBar), findsOneWidget);
    expect(find.byType(NavigationRail), findsNothing);
  });

  test('browser tabs share the session: the stored token is used, another account is left alone', () async {
    final be = FakeBackend();
    final store = MemoryTokenStore();
    final a = be.api(store), b = be.api(store);
    await a.login('grim', 'correct-horse');
    expect(await b.restore(), isTrue); // rotates what tab a holds in memory
    be.validAccess.clear();
    await a.get('/v1/character');
    expect(a.hasSession, isTrue);
    expect(store.token, 'refresh-${be.refreshCalls + 1}');

    // Tab b signs in to another account: tab a ends, b's session stays.
    be.accountId = 'acct-2';
    await b.login('mira', 'correct-horse');
    be.validAccess.clear();
    var lost = 0;
    a.onSessionLost = () => lost++;
    await expectLater(a.get('/v1/character'), throwsA(isA<ApiError>()));
    expect((lost, a.hasSession), (1, false));
    expect(store.token, isNotNull);
    expect(await be.api(store).restore(), isTrue, reason: 'the other account is still signed in');
  });
}
