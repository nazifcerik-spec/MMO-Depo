import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/app/app.dart';
import 'package:otter/core/api_client.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/auth/auth_controller.dart';
import 'package:otter/features/guide/guide.dart';
import 'package:otter/features/inventory/bank_screen.dart';
import 'package:otter/features/inventory/inventory_models.dart' show equipSlots;
import 'package:otter/features/inventory/item_widgets.dart';
import 'package:otter/features/pets/pets_screen.dart' show PetArt;
import 'package:otter/features/player/player_repository.dart';
import 'package:otter/features/player/player_screen.dart';
import 'package:otter/features/ranking/ranking_screen.dart';
import 'package:otter/l10n/app_localizations.dart';
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

  testWidgets('wide screens: the side menu opens every place beside it (DD-370)', (tester) async {
    _size(tester, 1440, 900);
    final now = DateTime.now().toUtc();
    final be = FakeBackend()
      ..characters.add({
        'id': 'char-1',
        'name': 'Grim',
        'race_code': 'orc',
        'class_code': 'rogue',
        'level': 3,
        'state': 'active',
      })
      ..action = {
        'kind': 'hunt',
        'detail': 'young_boar',
        'started_at': now.subtract(const Duration(minutes: 10)).toIso8601String(),
        'ends_at': now.add(const Duration(hours: 1)).toIso8601String(),
      }
      ..profiles['Grim'] = {
        'character_id': 'char-1',
        'name': 'Grim',
        'race_code': 'orc',
        'class_code': 'rogue',
        'level': 3,
      };
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
    await tester.pump();
    // The job card's clock never settles: pump through instead.
    Future<void> settle() async {
      for (var i = 0; i < 12; i++) {
        await tester.pump(const Duration(milliseconds: 100));
      }
    }

    await settle();
    expect(find.byKey(const Key('game-menu')), findsOneWidget);
    expect(find.byType(NavigationBar), findsNothing);
    expect(find.byType(AppBar), findsNothing, reason: 'the menu replaces the top bar');
    expect(find.byType(PlayerBanner), findsOneWidget, reason: 'the overview comes first');
    expect(find.byKey(const Key('job-card')).hitTestable(), findsOneWidget);
    expect(find.textContaining('Boar'), findsWidgets);

    Finder entry(String code) => find.byKey(ValueKey('menu-$code'));
    Future<void> open(String code) async {
      await tester.ensureVisible(entry(code));
      await tester.pump();
      await tester.tap(entry(code));
      await settle();
      expect(tester.takeException(), isNull, reason: code);
      expect(find.byKey(const Key('game-menu')), findsOneWidget, reason: '$code keeps the menu');
    }

    for (final code in [
      'character', 'bag', 'bank', 'quests', 'story', 'hunt', 'dungeon', 'travel', 'events', 'gather', 'craft', //
      'helpers', 'vendor', 'market', 'trade', 'credits', 'chat', 'friends', 'mail', 'guild', 'ranking', 'progress',
      'pets', 'profile', 'notifications', 'guide', 'support', 'overview',
    ]) {
      await open(code);
    }
    expect(find.byType(SegmentedButton<int>).hitTestable(), findsNothing);

    // A page opened from the menu replaces the last one, and goes back to the tab.
    await open('ranking');
    await open('bank');
    expect(find.byType(RankingScreen), findsNothing);
    await tester.tap(find.byType(BackButton).hitTestable());
    await settle();
    expect(find.byType(BankScreen), findsNothing);
    expect(find.byType(PlayerBanner), findsOneWidget);

    // The job card goes where the job is collected.
    await tester.tap(find.byKey(const Key('job-card')).first);
    await settle();
    expect(find.text('Hunt').hitTestable(), findsWidgets);

    // Less room: icons only. Less still: the phone layout, on the same tab.
    tester.view.physicalSize = const Size(1000, 900);
    await settle();
    expect(tester.getSize(find.byKey(const Key('game-menu'))).width, 76);
    tester.view.physicalSize = const Size(500, 900);
    await settle();
    expect(find.byKey(const Key('game-menu')), findsNothing);
    expect(tester.widget<NavigationBar>(find.byType(NavigationBar)).selectedIndex, 3);
    expect(tester.takeException(), isNull);
  });

  for (final w in [390.0, 1280.0]) {
    testWidgets('a player page shows the character and its deeds at $w (DD-371)', (tester) async {
      _size(tester, w, 900);
      final be = FakeBackend()
        ..profiles['Şule'] = {
          'character_id': 'other',
          'name': 'Şule',
          'race_code': 'human',
          'class_code': 'mage',
          'spec_code': 'pyromancer',
          'level': 712,
          'title': 'x',
          'avatar': 'av07',
          'background': 'bg03',
          'guild_name': 'Gece Nöbeti',
          'guild_tag': 'GN',
          'score': 12345,
          'rank': 4,
          'achievement_points': 10,
          'professions': [
            {'code': 'mining', 'level': 40},
          ],
          'equipment': {
            'mainhand': {'template_code': 'sword_iron', 'item_class': 'weapon', 'rarity': 'epic'},
          },
          'pet': {'species': 'flamewing', 'level': 40},
        };
      final api = be.api(MemoryTokenStore());
      await api.login('grim', 'correct-horse');
      Widget app(String name) => MaterialApp(
        locale: const Locale('en'),
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: PlayerScreen(key: ValueKey(name), repo: PlayerRepository(api, 'char-1'), name: name),
      );
      await tester.pumpWidget(app('Şule'));
      await tester.pumpAndSettle();
      expect(find.text('Şule'), findsOneWidget);
      expect(find.text('<GN> Gece Nöbeti'), findsOneWidget);
      expect(find.text('Rank #4'), findsOneWidget);
      expect(find.text('12345'), findsOneWidget);
      expect(find.byType(ItemCell), findsNWidgets(equipSlots.length));
      await tester.scrollUntilVisible(find.byType(PetArt), 200, scrollable: find.byType(Scrollable).first);
      expect(tester.takeException(), isNull);

      await tester.pumpWidget(app('Nobody'));
      await tester.pumpAndSettle();
      expect(find.text('No such player.'), findsOneWidget);
    });
  }

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
