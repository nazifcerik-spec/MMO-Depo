import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/app/app.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/auth/auth_controller.dart';
import 'package:otter/features/character/models.dart';

import 'fake_backend.dart';
import 'fake_ws.dart';

Future<FakeSocketServer> _boot(
  WidgetTester tester,
  FakeBackend be, {
  TokenStore? store,
  Locale locale = const Locale('en'),
}) async {
  tester.view.physicalSize = const Size(375 * 3, 812 * 3); // mobile-first viewport
  tester.view.devicePixelRatio = 3;
  addTearDown(tester.view.reset);
  final ws = FakeSocketServer();
  await tester.pumpWidget(
    OtterApp(
      auth: AuthController(be.api(store), wsConnect: ws.connect),
      locale: locale,
    ),
  );
  await tester.pumpAndSettle();
  return ws;
}

/// Signed-in player with one character, entered and on the sheet tab.
Future<(FakeBackend, FakeSocketServer, MemoryTokenStore)> _inGame(WidgetTester tester) async {
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
  final ws = await _boot(tester, be, store: store);
  await tester.tap(find.text('Grim'));
  await tester.pumpAndSettle();
  return (be, ws, store);
}

/// The vertical list inside the visible panel (tab views also scroll sideways).
final _list = find.descendant(of: find.byType(ListView), matching: find.byType(Scrollable)).last;

void main() {
  testWidgets('login validates, shows server errors, then lands on character select', (tester) async {
    final be = FakeBackend();
    await _boot(tester, be);

    await tester.tap(find.text('Log in'));
    await tester.pump();
    expect(find.text('Required'), findsNWidgets(2));

    await tester.enterText(find.byType(TextFormField).at(0), 'grim');
    await tester.enterText(find.byType(TextFormField).at(1), 'wrong');
    await tester.tap(find.text('Log in'));
    await tester.pumpAndSettle();
    expect(find.text('Wrong email/username or password.'), findsOneWidget);

    await tester.enterText(find.byType(TextFormField).at(1), 'correct-horse');
    await tester.tap(find.text('Log in'));
    await tester.pumpAndSettle();
    expect(find.text('Characters'), findsOneWidget);
    expect(find.text('No characters yet. Create your first hero.'), findsOneWidget);
  });

  testWidgets('register form enforces confirm + terms client-side', (tester) async {
    final be = FakeBackend();
    await _boot(tester, be);
    await tester.tap(find.text('New here? Create an account'));
    await tester.pumpAndSettle();

    final fields = find.byType(TextFormField);
    await tester.enterText(fields.at(0), 'grim@example.com');
    await tester.enterText(fields.at(1), 'grim');
    await tester.enterText(fields.at(2), 'correct-horse');
    await tester.enterText(fields.at(3), 'different');
    await tester.tap(find.text('Create account'));
    await tester.pump();
    expect(find.text('Passwords do not match'), findsOneWidget);

    await tester.enterText(fields.at(3), 'correct-horse');
    await tester.tap(find.text('Create account'));
    await tester.pump();
    expect(find.text('Please accept the terms to continue'), findsOneWidget);
    expect(be.requests, isEmpty, reason: 'nothing sent until the form is valid');

    await tester.tap(find.byType(Checkbox));
    await tester.tap(find.text('Create account'));
    await tester.pumpAndSettle();
    expect(find.text('Characters'), findsOneWidget);
  });

  testWidgets('create a character, enter it, and read the sheet', (tester) async {
    final be = FakeBackend();
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse'); // persisted session → restore on boot
    await _boot(tester, be, store: store);
    expect(find.text('Characters'), findsOneWidget);

    await tester.tap(find.text('New character'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextFormField), '1bad');
    await tester.tap(find.text('Create'));
    await tester.pump();
    expect(find.text('3–16 characters, starts with a letter (A–Z, 0–9, _)'), findsWidgets);

    await tester.enterText(find.byType(TextFormField), 'Grim');
    await tester.tap(find.text('Dwarf'));
    await tester.tap(find.text('Hunter'));
    await tester.tap(find.text('Create'));
    await tester.pumpAndSettle();
    expect(be.characters.single['race_code'], 'dwarf');
    expect(find.text('Lvl 1 · Hunter · Dwarf'), findsOneWidget);
    expect(find.text('1/2 slots'), findsOneWidget);

    await tester.tap(find.text('Grim'));
    await tester.pumpAndSettle();
    expect(find.text('80/120'), findsOneWidget); // HP
    expect(find.text('40 / 100'), findsOneWidget); // XP
    expect(find.text('Attributes'), findsOneWidget);
    expect(find.text('90%'), findsOneWidget); // hit
    expect(find.textContaining('12', findRichText: true), findsWidgets); // 12g 34s 56c
    expect(be.requests.last.headers['X-Character-Id'], 'char-1');

    await tester.tap(find.byTooltip('Log out'));
    await tester.pumpAndSettle();
    expect(find.text('Log in'), findsOneWidget);
    expect(store.token, isNull);
  });

  testWidgets('forced sign-out (refresh rejected) returns to login with a notice', (tester) async {
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
    await _boot(tester, be, store: store);
    await tester.tap(find.text('Grim'));
    await tester.pumpAndSettle();

    be.validAccess.clear();
    be.rejectRefresh = true;
    await tester.fling(find.byType(ListView), const Offset(0, 400), 1000);
    await tester.pumpAndSettle();

    expect(find.text('Log in'), findsOneWidget);
    expect(find.text('Your session expired. Please log in again.'), findsOneWidget);
  });

  testWidgets('gathering loop: start, collect early with confirm, loot sheet, sheet reload', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.text('Gather'));
    await tester.pumpAndSettle();
    expect(find.text('Choose an activity'), findsOneWidget);

    await tester.tap(find.text('Fishing'));
    await tester.tap(find.text('2h 00m'));
    await tester.pump();
    await tester.tap(find.text('Start Fishing'));
    await tester.pumpAndSettle();
    expect(be.activeActivity!['kind'], 'fishing');
    expect(be.activeActivity!['duration_seconds'], 7200);
    expect(find.textContaining('left'), findsOneWidget);

    final sheetsBefore = be.sheetGets;
    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    expect(find.text('Collect early?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last);
    await tester.pumpAndSettle();

    expect(find.text('Copper Ore'), findsOneWidget);
    expect(find.text('Flawless Raw Sapphire'), findsOneWidget);
    await tester.tap(find.text('Nice!'));
    await tester.pumpAndSettle();
    expect(find.text('Choose an activity'), findsOneWidget);
    expect(be.activeActivity, isNull);
    expect(be.sheetGets, greaterThan(sheetsBefore), reason: 'reveal bumps other tabs');
  });

  testWidgets('abandon asks first and forfeits', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.text('Gather'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Start Mining'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Abandon'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Cancel'));
    await tester.pumpAndSettle();
    expect(be.activeActivity, isNotNull);

    await tester.tap(find.text('Abandon'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Abandon'));
    await tester.pumpAndSettle();
    expect(be.activeActivity, isNull);
    expect(find.text('Choose an activity'), findsOneWidget);
  });

  testWidgets('live socket: token handshake, push toast + refresh, closed on logout', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    final socket = ws.last;
    expect(socket.uri.scheme, 'ws');
    expect(socket.uri.path, '/v1/ws');
    expect(socket.uri.queryParameters['token'], startsWith('access-'));

    final before = be.sheetGets;
    socket.push('achievement.unlocked', {'code': 'first_steps'});
    await tester.pumpAndSettle();
    expect(find.text('Achievement unlocked!'), findsOneWidget);
    expect(be.sheetGets, greaterThan(before));

    await tester.tap(find.byTooltip('Log out'));
    await tester.pumpAndSettle();
    expect(socket.closed, isTrue);
  });

  testWidgets('bag: item detail, equip/unequip, class gate, sell, destroy', (tester) async {
    final semantics = tester.ensureSemantics();
    final (be, _, _) = await _inGame(tester);
    Map<String, dynamic> item(String id, String code, int slot, {Map<String, dynamic> extra = const {}}) => {
      'instance_id': id,
      'template_code': code,
      'item_class': 'material',
      'item_subclass': 'ore',
      'slot': slot,
      'stack_quantity': 1,
      'max_stack': 1,
      'rarity': 'common',
      'is_bound': false,
      'is_locked': false,
      'weight_grams_per': 200,
      'base_value_copper': 100,
      'min_level': 1,
      ...extra,
    };
    be.bag.addAll([
      item(
        's1',
        'sword_iron',
        0,
        extra: {
          'item_class': 'weapon',
          'rarity': 'rare',
          'prefix': 'sharp',
          'equip_slot': 'mainhand',
          'stat_bonuses': [
            {'kind': 'flat', 'target': 'attack', 'amount': 10},
            {'kind': 'mult_add', 'target': 'str', 'amount': 0.1},
          ],
        },
      ),
      item('o1', 'iron_ore', 1, extra: {'stack_quantity': 12, 'max_stack': 99, 'base_value_copper': 5}),
      item(
        'm1',
        'staff_oak',
        2,
        extra: {'item_class': 'weapon', 'equip_slot': 'mainhand', 'class_restriction': 'mage'},
      ),
      item('j1', 'rat_tail', 3),
    ]);

    await tester.tap(find.byIcon(Icons.backpack_outlined));
    await tester.pumpAndSettle();
    expect(find.text('Equipment'), findsOneWidget);

    // Detail → equip.
    await tester.tap(find.bySemanticsLabel('Sharp Sword Iron ×1'));
    await tester.pumpAndSettle();
    expect(find.text('+10 Attack'), findsOneWidget);
    expect(find.text('+10% STR'), findsOneWidget);
    final sheets = be.sheetGets;
    await tester.tap(find.text('Equip'));
    await tester.pumpAndSettle();
    expect(be.equipped['mainhand']?['instance_id'], 's1');
    expect(be.sheetGets, greaterThan(sheets), reason: 'gear change reloads the sheet');

    // Equipped slot → unequip back to the bag.
    await tester.tap(find.bySemanticsLabel('Main hand: Sharp Sword Iron'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Unequip'));
    await tester.pumpAndSettle();
    expect(be.equipped, isEmpty);
    expect(be.bag.any((i) => i['instance_id'] == 's1'), isTrue);

    // Server-side class gate surfaces as a readable error.
    await tester.tap(find.bySemanticsLabel('Staff Oak ×1'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Equip'));
    await tester.pumpAndSettle();
    expect(find.text("Your class can't use this item."), findsOneWidget);

    // Sell with confirm; server price is shown (12 × 5c × 0.25 = 15c).
    await tester.tap(find.bySemanticsLabel('Iron Ore ×12'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Sell'));
    await tester.pumpAndSettle();
    expect(find.text('Sell Iron Ore?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Sell'));
    await tester.pumpAndSettle();
    expect(be.bag.any((i) => i['instance_id'] == 'o1'), isFalse);
    expect(find.textContaining('Sold for'), findsOneWidget);

    // Destroy with confirm.
    await tester.tap(find.bySemanticsLabel('Rat Tail ×1'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Destroy'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Destroy'));
    await tester.pumpAndSettle();
    expect(be.bag.any((i) => i['instance_id'] == 'j1'), isFalse);
    semantics.dispose();
  });

  testWidgets('hunt: level-aware default, start, early collect, result with level-up', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    expect(find.text('Choose your prey'), findsOneWidget);
    expect(find.text('Cave Troll'), findsOneWidget);

    // Default = toughest mob at or below level 3.
    await tester.tap(find.text('Start hunt'));
    await tester.pumpAndSettle();
    expect(be.activeHunt!['mob_code'], 'grey_wolf');
    await tester.tap(find.text('Abandon'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Abandon'));
    await tester.pumpAndSettle();
    expect(be.activeHunt, isNull);

    await tester.tap(find.text('Giant Rat'));
    await tester.tap(find.text('30m'));
    await tester.pump();
    await tester.tap(find.text('Start hunt'));
    await tester.pumpAndSettle();
    expect(be.activeHunt!['mob_code'], 'giant_rat');

    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last);
    await tester.pumpAndSettle();
    expect(find.text('7 won · 1 lost · Level up! Now level 4'), findsOneWidget);
    expect(find.text('+160 XP'), findsOneWidget);
    expect(find.text('Rare × 1'), findsOneWidget);
    await tester.tap(find.text('Nice!'));
    await tester.pumpAndSettle();
    expect(find.text('Choose your prey'), findsOneWidget);
  });

  testWidgets('dungeon: level gate, collect locked until the run ends, cleared result', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Dungeon'));
    await tester.pumpAndSettle();

    final enters = find.widgetWithText(FilledButton, 'Enter');
    expect(tester.widget<FilledButton>(enters.at(1)).onPressed, isNull, reason: 'Sunken Crypt needs level 18');

    await tester.tap(enters.first);
    await tester.pumpAndSettle();
    expect(be.activeDungeon!['dungeon_code'], 'rat_warrens');
    expect(find.text('Rewards unlock when the run ends.'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Collect')).onPressed, isNull);

    // Fast-forward: the run is over server-side.
    be.activeDungeon!['expires_at'] = DateTime.now().toUtc().toIso8601String();
    await tester.fling(find.byType(ListView).last, const Offset(0, 400), 1000);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    expect(find.text('Cleared! · Boss defeated · 7 won · 0 lost'), findsOneWidget);
    expect(find.text('Rat Crown'), findsOneWidget);
  });

  testWidgets('travel: routes with level gate, depart, turn back', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Travel'));
    await tester.pumpAndSettle();
    expect(find.text('Ironhold City'), findsOneWidget);

    final go = find.widgetWithText(FilledButton, 'Go');
    expect(tester.widget<FilledButton>(go.at(1)).onPressed, isNull, reason: 'Dragon Peaks needs level 40');
    await tester.tap(go.first);
    await tester.pumpAndSettle();
    expect(be.travel!['to_region_code'], 'whisperwood');
    expect(find.text('Travelling to Whisperwood'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Arrive')).onPressed, isNull);

    await tester.tap(find.text('Turn back'));
    await tester.pumpAndSettle();
    expect(be.travel, isNull);
    expect(find.text('Routes'), findsOneWidget);
  });

  testWidgets('workshop: learn, material-aware craft, skill gate, skill-up toast', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.bag.add({
      'instance_id': 'seed',
      'template_code': 'linseed',
      'item_class': 'material',
      'item_subclass': 'herb',
      'slot': 0,
      'stack_quantity': 3,
      'max_stack': 99,
      'rarity': 'common',
    });
    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();
    expect(find.text('Alchemy'), findsOneWidget);
    expect(find.widgetWithText(FilledButton, 'Learn'), findsNWidgets(2));

    await tester.tap(find.widgetWithText(FilledButton, 'Learn').first);
    await tester.pumpAndSettle();
    expect(be.skills, {'alchemy': 1});
    expect(find.text('Linseed 3/2'), findsOneWidget);
    expect(find.text('Needs skill 3'), findsOneWidget);

    final crafts = find.widgetWithText(FilledButton, 'Craft');
    expect(tester.widget<FilledButton>(crafts.at(1)).onPressed, isNull, reason: 'skill 1 < min 3');
    await tester.tap(crafts.first);
    await tester.pumpAndSettle();
    expect(find.text('Crafted Potion Hp Small ×1 · Alchemy is now level 2'), findsOneWidget);
    expect(find.text('Linseed 1/2'), findsOneWidget, reason: 'materials reload after crafting');
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Craft').first).onPressed, isNull);
  });

  testWidgets('shop: affordable buy with price guard, buy back a sold item', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.buybackRows.add({
      'id': 'bb-1',
      'vendor_code': 'general_ironhold',
      'template_code': 'iron_ore',
      'rarity': 'common',
      'stack_quantity': 4,
      'sold_price_copper': 12,
      'expires_at': DateTime.now().toUtc().add(const Duration(minutes: 90)).toIso8601String(),
    });
    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Shop'));
    await tester.pumpAndSettle();

    final buys = find.widgetWithText(FilledButton, 'Buy');
    expect(tester.widget<FilledButton>(buys.at(1)).onPressed, isNull, reason: 'Sword Steel costs more than we have');
    await tester.tap(buys.first);
    await tester.pumpAndSettle();
    expect(be.bought, ['potion_hp_small']);
    expect(find.text('Bought Potion Hp Small'), findsOneWidget);

    expect(find.text('Iron Ore ×4'), findsOneWidget);
    final row = find.ancestor(of: find.text('Iron Ore ×4'), matching: find.byType(ListTile));
    await tester.tap(find.descendant(of: row, matching: find.byType(FilledButton)));
    await tester.pumpAndSettle();
    expect(be.buybackRows, isEmpty);
    expect(find.text('Iron Ore ×4'), findsNothing);
  });

  testWidgets('market: filter, own listing not buyable, confirmed buy with price echo', (tester) async {
    final (be, _, _) = await _inGame(tester);
    Map<String, dynamic> listing(String id, String seller, String code, String cls, int qty, int unit) => {
      'id': id,
      'seller_character_id': seller,
      'seller_name': seller == 'char-1' ? 'Grim' : 'Brak',
      'template_code': code,
      'item_class': cls,
      'rarity': 'common',
      'stack_quantity': qty,
      'unit_price_copper': unit,
      'total_price_copper': qty * unit,
      'expires_at': DateTime.now().toUtc().add(const Duration(hours: 20)).toIso8601String(),
      'state': 'active',
    };
    be.marketListings.addAll([
      listing('l-1', 'char-9', 'potion_hp_small', 'consumable', 3, 20),
      listing('l-2', 'char-1', 'sword_iron', 'weapon', 1, 500),
    ]);
    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Market'));
    await tester.pumpAndSettle();

    Finder buyButton(String name) => find.descendant(
      of: find.ancestor(of: find.text(name), matching: find.byType(ListTile)),
      matching: find.byType(FilledButton),
    );
    expect(tester.widget<FilledButton>(buyButton('Sword Iron')).onPressed, isNull, reason: 'own listing');

    await tester.tap(buyButton('Potion Hp Small ×3'));
    await tester.pumpAndSettle();
    expect(find.text('Buy Potion Hp Small ×3?'), findsOneWidget);
    expect(find.text('Pay 60c to Brak.'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Buy'));
    await tester.pumpAndSettle();
    expect(be.marketBought, ['l-1']);
    expect(find.text('Bought Potion Hp Small ×3'), findsOneWidget);

    await tester.tap(find.text('Weapons'));
    await tester.pumpAndSettle();
    expect(be.requests.lastWhere((r) => r.url.path == '/v1/market/listings').url.queryParameters['class'], 'weapon');
  });

  testWidgets('market: list a stack from the bag with price suggestion + fee preview, then cancel', (tester) async {
    final semantics = tester.ensureSemantics();
    final (be, _, _) = await _inGame(tester);
    be.bag.add({
      'instance_id': 'ore',
      'template_code': 'iron_ore',
      'item_class': 'material',
      'item_subclass': 'ore',
      'slot': 0,
      'stack_quantity': 12,
      'max_stack': 99,
      'rarity': 'common',
    });
    await tester.tap(find.byIcon(Icons.backpack_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.bySemanticsLabel('Iron Ore ×12'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'List on market'));
    await tester.pumpAndSettle();

    expect(find.text('Market price: 35c'), findsOneWidget);
    expect(find.text('Quantity: 12 of 12'), findsOneWidget);
    await tester.tap(find.text('12h 00m'));
    await tester.pump();
    // 35c × 12 × 5% × 1.0 = 21c fee; 420c − 5% cut = 399c.
    expect(find.textContaining('21', findRichText: true), findsWidgets);
    await tester.tap(find.widgetWithText(FilledButton, 'List on market'));
    await tester.pumpAndSettle();
    expect(be.lastSell, {
      'item_instance_id': 'ore',
      'stack_quantity': 12,
      'unit_price_copper': 35,
      'duration_hours': 12,
    });
    expect(find.text('Iron Ore listed on the market'), findsOneWidget);

    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Market'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('My listings (1)'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Cancel listing'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Cancel listing'));
    await tester.pumpAndSettle();
    expect(be.myListings, isEmpty);
    semantics.dispose();
  });

  testWidgets('chat: history, send, live push with dedup, blocked hidden, channel gate', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    be.chat.add({
      'id': 1,
      'channel': 'global',
      'from_character_id': 'char-9',
      'from_name': 'Brak',
      'body': 'hi all',
      'sent_at': DateTime.now().toUtc().toIso8601String(),
    });
    be.blocked.add({'character_id': 'char-7', 'name': 'Troll'});
    await tester.tap(find.byTooltip('Social'));
    await tester.pumpAndSettle();
    expect(find.textContaining('hi all', findRichText: true), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'hello there');
    await tester.tap(find.byTooltip('Send'));
    await tester.pumpAndSettle();
    expect(be.chatSent, ['hello there']);
    expect(find.textContaining('hello there', findRichText: true), findsOneWidget);

    Map<String, dynamic> push(int id, String from, String fromId, String text) => {
      'id': id,
      'channel': 'global',
      'from_character_id': fromId,
      'from_name': from,
      'body': text,
      'sent_at': DateTime.now().toUtc().toIso8601String(),
    };
    ws.last.push('chat.message', push(99, 'Brak', 'char-9', 'live one'));
    ws.last.push('chat.message', push(99, 'Brak', 'char-9', 'live one')); // duplicate delivery
    ws.last.push('chat.message', push(98, 'Troll', 'char-7', 'spam'));
    await tester.pumpAndSettle();
    expect(find.textContaining('live one', findRichText: true), findsOneWidget);
    expect(find.textContaining('spam', findRichText: true), findsNothing, reason: 'blocked sender');

    await tester.tap(find.widgetWithText(ChoiceChip, 'Guild'));
    await tester.pumpAndSettle();
    expect(find.text("You're not in that group."), findsOneWidget);
  });

  testWidgets('friends: accept request, add by name errors, whisper jumps to chat', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.friendRequests.add({'id': 'fr-1', 'from_character_id': 'char-5', 'from_name': 'Aria'});
    await tester.tap(find.byTooltip('Social'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Friends'));
    await tester.pumpAndSettle();

    expect(find.text('Friend requests'), findsOneWidget);
    await tester.tap(find.byTooltip('Accept'));
    await tester.pumpAndSettle();
    expect(be.friends.single['name'], 'Aria');
    expect(find.text('Friends (1)'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Nobody');
    await tester.tap(find.widgetWithText(FilledButton, 'Add'));
    await tester.pumpAndSettle();
    expect(find.text('No character with that name.'), findsOneWidget);
    await tester.enterText(find.byType(TextField), 'Brak');
    await tester.tap(find.widgetWithText(FilledButton, 'Add'));
    await tester.pumpAndSettle();
    expect(find.text('Friend request sent to Brak'), findsOneWidget);

    await tester.tap(find.byTooltip('More'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Whisper'));
    await tester.pumpAndSettle();
    expect(find.widgetWithText(ChoiceChip, 'Aria'), findsOneWidget);
    expect(be.requests.where((r) => r.url.path.endsWith('/chat/whisper')).last.url.queryParameters['target'], 'Aria');
  });

  testWidgets('mail: open marks read, take gold attachment, compose with gold', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.mail.add({
      'id': 'm-1',
      'sender_name': 'Ironhold Bank',
      'subject': 'Welcome',
      'gold_attachment_copper': 5000,
      'item_attachment_count': 0,
      'state': 'sent',
      'created_at': DateTime.now().toUtc().toIso8601String(),
      'expires_at': DateTime.now().toUtc().add(const Duration(days: 30)).toIso8601String(),
    });
    await tester.tap(find.byTooltip('Social'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Mail'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Welcome'));
    await tester.pumpAndSettle();
    expect(find.text('Welcome to Ironhold!'), findsOneWidget);
    expect(be.mail.single['state'], 'read');
    await tester.tap(find.text('Take attachments'));
    await tester.pumpAndSettle();
    expect(be.mail.single['state'], 'claimed');
    expect(find.text('Attachments moved to your bag and wallet'), findsOneWidget);

    await tester.tap(find.text('Write mail'));
    await tester.pumpAndSettle();
    final fields = find.byType(TextField);
    await tester.enterText(fields.at(0), 'Aria');
    await tester.enterText(fields.at(1), 'Gift');
    await tester.enterText(fields.at(2), 'For you');
    await tester.enterText(fields.at(3), '1');
    await tester.tap(find.widgetWithText(FilledButton, 'Send'));
    await tester.pumpAndSettle();
    expect(be.lastMail, {'recipient_name': 'Aria', 'subject': 'Gift', 'body': 'For you', 'gold_copper': 10000});
    expect(find.text('Mail sent (postage 7c)'), findsOneWidget);
  });

  testWidgets('use a potion: effect shown, HP toast, stack shrinks, refused at full HP', (tester) async {
    final semantics = tester.ensureSemantics();
    final (be, _, _) = await _inGame(tester);
    be.bag.add({
      'instance_id': 'pot',
      'template_code': 'potion_hp_small',
      'item_class': 'consumable',
      'item_subclass': 'potion',
      'slot': 0,
      'stack_quantity': 3,
      'max_stack': 20,
      'rarity': 'common',
      'use_effect': {'kind': 'heal_hp', 'amount': 60},
    });
    await tester.tap(find.byIcon(Icons.backpack_outlined));
    await tester.pumpAndSettle();

    await tester.tap(find.bySemanticsLabel('Potion Hp Small ×3'));
    await tester.pumpAndSettle();
    expect(find.text('Restores 60 HP'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Use'));
    await tester.pumpAndSettle();
    expect(find.text('Used Potion Hp Small · HP 90/120'), findsOneWidget);
    expect(be.bag.single['stack_quantity'], 2);

    be.hp = 120;
    await tester.tap(find.bySemanticsLabel('Potion Hp Small ×2'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Use'));
    await tester.pumpAndSettle();
    expect(find.text('Your health is already full.'), findsOneWidget);
    expect(be.bag.single['stack_quantity'], 2, reason: 'refused use consumes nothing');
    semantics.dispose();
  });

  testWidgets('forgot password: dialog prefills e-mail and shows the neutral confirmation', (tester) async {
    final be = FakeBackend();
    await _boot(tester, be);
    await tester.enterText(find.byType(TextFormField).at(0), 'grim@example.com');
    await tester.tap(find.text('Forgot password?'));
    await tester.pumpAndSettle();
    expect(
      find.descendant(of: find.byType(AlertDialog), matching: find.widgetWithText(TextField, 'grim@example.com')),
      findsOneWidget,
    );
    await tester.tap(find.text('Send link'));
    await tester.pumpAndSettle();
    expect(be.resetRequestedFor, ['grim@example.com']);
    expect(find.textContaining('If an account exists'), findsOneWidget);
  });

  testWidgets('unverified e-mail banner resends; change password rejects wrong current without refreshing', (
    tester,
  ) async {
    final be = FakeBackend()..emailVerified = false;
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await _boot(tester, be, store: store);

    expect(find.text('Please confirm your e-mail address.'), findsOneWidget);
    await tester.tap(find.text('Resend'));
    await tester.pumpAndSettle();
    expect(be.verificationResends, 1);
    expect(find.text('Verification e-mail sent.'), findsOneWidget);

    await tester.tap(find.byTooltip('Change password'));
    await tester.pumpAndSettle();
    final fields = find.byType(TextField);
    await tester.enterText(fields.at(0), 'wrong-horse');
    await tester.enterText(fields.at(1), 'new-horse-1');
    await tester.enterText(fields.at(2), 'new-horse-1');
    final refreshesBefore = be.refreshCalls;
    await tester.tap(find.text('Save'));
    await tester.pumpAndSettle();
    expect(find.text('Wrong email/username or password.'), findsOneWidget);
    expect(be.refreshCalls, refreshesBefore, reason: 'credential 401 must not trigger a token refresh');

    await tester.enterText(fields.at(0), 'correct-horse');
    await tester.tap(find.text('Save'));
    await tester.pumpAndSettle();
    expect(be.passwordChangedTo, 'new-horse-1');
    expect(find.text('Password changed. Your other devices were signed out.'), findsOneWidget);
  });

  testWidgets('party: accept invite, then create/invite/promote as leader; push reloads', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    be.partyInvites.add({'id': 'pi-1', 'invited_by_name': 'Aria'});
    await tester.tap(find.byTooltip('Social'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.widgetWithText(Tab, 'Party'));
    await tester.tap(find.widgetWithText(Tab, 'Party'));
    await tester.pumpAndSettle();

    expect(find.text('Invitations'), findsOneWidget);
    await tester.tap(find.byTooltip('Accept'));
    await tester.pumpAndSettle();
    expect(find.text('Party (2/5)'), findsOneWidget);
    expect(find.byTooltip('More'), findsNothing, reason: 'members cannot manage the party');

    await tester.tap(find.text('Leave'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Leave'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Create party'));
    await tester.pumpAndSettle();
    expect(find.text('Party (1/5)'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Brak');
    await tester.tap(find.widgetWithText(FilledButton, 'Invite'));
    await tester.pumpAndSettle();
    expect(be.groupCalls, contains('party/invite {"target_name":"Brak"}'));

    // Someone joins: the server pushes party.changed and the panel reloads.
    (be.party!['members'] as List).add({'character_id': 'char-7', 'name': 'Brak', 'level': 2, 'class_code': 'mage'});
    ws.last.push('party.changed');
    await tester.pumpAndSettle();
    expect(find.text('Party (2/5)'), findsOneWidget);

    await tester.tap(find.byTooltip('More'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Make leader'));
    await tester.pumpAndSettle();
    expect(be.party!['leader_character_id'], 'char-7');
    expect(find.byTooltip('More'), findsNothing, reason: 'no longer leader');
  });

  testWidgets('guild: found with fee confirm, MOTD, rank change within permissions, schedule + cancel disband', (
    tester,
  ) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.byTooltip('Social'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.widgetWithText(Tab, 'Guild'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(Tab, 'Guild'));
    await tester.pumpAndSettle();

    await tester.enterText(find.widgetWithText(TextField, 'Guild name'), 'Iron Wolves');
    await tester.enterText(find.widgetWithText(TextField, 'Tag'), 'iwf');
    await tester.tap(find.text('Found guild (10g 0s 0c)'));
    await tester.pumpAndSettle();
    expect(find.text('Founding a guild costs 10g 0s 0c.'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Create'));
    await tester.pumpAndSettle();
    expect(be.groupCalls.last, 'guild {"name":"Iron Wolves","tag":"IWF"}');
    expect(find.text('Iron Wolves [IWF]'), findsOneWidget);

    await tester.tap(find.byTooltip('Edit message of the day'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.descendant(of: find.byType(AlertDialog), matching: find.byType(TextField)),
      'Raid at 9',
    );
    await tester.tap(find.text('Save'));
    await tester.pumpAndSettle();
    expect(find.text('Raid at 9'), findsOneWidget);

    await tester.tap(find.byTooltip('More'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Change rank'));
    await tester.pumpAndSettle();
    expect(find.text('Leader'), findsNothing, reason: 'leadership only moves via transfer');
    await tester.tap(find.text('Officer'));
    await tester.pumpAndSettle();
    expect(be.groupCalls.last, 'guild/rank {"target_character_id":"char-7","rank_index":1}');

    await tester.tap(find.text('Disband'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Disband'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Guild disbands in'), findsOneWidget);
    await tester.tap(find.text('Cancel disbanding'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Guild disbands in'), findsNothing);
  });

  testWidgets('notifications: badge, read one / all, live bump, settings keep quiet hours', (tester) async {
    Map<String, dynamic> note(String id, String kind, Map<String, dynamic> vars) => {
      'id': id,
      'kind': kind,
      'body': {'kind': kind, 'body_vars': vars},
      'created_at': DateTime.now().toUtc().subtract(const Duration(minutes: 5)).toIso8601String(),
    };
    final (be, ws, _) = await _inGame(tester);
    be.inbox.addAll([
      note('n-1', 'market_sold', {'template_code': 'sword_iron', 'net_copper': 12000}),
      note('n-2', 'guild_invite', {'guild_name': 'Iron Wolves'}),
    ]);
    ws.last.push('mail.received'); // any inbox-backed push refreshes the badge
    await tester.pumpAndSettle();
    expect(find.descendant(of: find.byType(Badge), matching: find.text('2')), findsOneWidget);

    await tester.tap(find.byTooltip('Notifications'));
    await tester.pumpAndSettle();
    expect(find.text('Item sold'), findsOneWidget);
    expect(find.textContaining('Sword Iron · +1g 20s 0c'), findsOneWidget);
    expect(find.textContaining('Iron Wolves'), findsOneWidget);

    await tester.tap(find.text('Item sold'));
    await tester.pumpAndSettle();
    expect(be.inbox.first['read_at'], isNotNull);
    await tester.tap(find.byTooltip('Mark all as read'));
    await tester.pumpAndSettle();
    expect(find.byTooltip('Mark all as read'), findsNothing);

    await tester.tap(find.byTooltip('Notification settings'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Pause push notifications'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Paused for'), findsOneWidget);
    await tester.tap(
      find.descendant(of: find.widgetWithText(ListTile, 'Item sold'), matching: find.text('Push + inbox')),
    );
    await tester.pumpAndSettle();
    await tester.tap(find.text('Inbox only').last);
    await tester.pumpAndSettle();
    expect(be.notifKinds['market_sold'], 'inbox_only');
    expect(be.lastPrefsPut!['quiet_hours_start'], '22:00', reason: 'PUT replaces quiet hours wholesale');

    await tester.pageBack();
    await tester.pumpAndSettle();
    await tester.pageBack();
    await tester.pumpAndSettle();
    expect(find.byType(Badge), findsOneWidget);
    expect(find.descendant(of: find.byType(Badge), matching: find.text('2')), findsNothing);
  });

  testWidgets('progression: achievements sorted, mastery bonus, reputation, title select + clear', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.scrollUntilVisible(find.text('Achievements & mastery'), 200);
    await tester.tap(find.text('Achievements & mastery'));
    await tester.pumpAndSettle();

    expect(find.text('5 achievement points'), findsOneWidget);
    expect(
      tester.getTopLeft(find.text('First Blood')).dy,
      lessThan(tester.getTopLeft(find.text('Slayer 100')).dy),
      reason: 'unlocked first',
    );
    expect(find.text('37/100'), findsOneWidget);
    await tester.scrollUntilVisible(find.textContaining('Title: Prospector'), 200, scrollable: _list);

    await tester.tap(find.text('Mastery'));
    await tester.pumpAndSettle();
    expect(find.text('Tier 1/10'), findsOneWidget);
    expect(find.text('150 / 500'), findsOneWidget);
    expect(find.text('+1/h Mining yield'), findsOneWidget);

    await tester.tap(find.text('Reputation'));
    await tester.pumpAndSettle();
    expect(find.text('Neutral · 40'), findsOneWidget);
    expect(find.text('+1460 today'), findsOneWidget);

    await tester.tap(find.text('Titles'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Prospector'));
    await tester.pumpAndSettle();
    expect(be.activeTitle, 'prospector');
    await tester.tap(find.text('No title'));
    await tester.pumpAndSettle();
    expect(be.activeTitle, isNull);
  });

  testWidgets('events: seasonal claim once, enlist with confirm, level gate, closed result', (tester) async {
    final (be, _, _) = await _inGame(tester);
    DateTime at(Duration d) => DateTime.now().toUtc().add(d);
    Map<String, dynamic> spawn(String id, String code, String kind, String state, int minLevel, Duration start) => {
      'id': id,
      'template_code': code,
      'kind': kind,
      'min_level': minLevel,
      'goal_score': 500000,
      'total_score': 125000,
      'starts_at': at(start).toIso8601String(),
      'ends_at': at(start + const Duration(hours: 1)).toIso8601String(),
      'state': state,
      'tiers': <Object>[],
      'enlisted': false,
    };
    be.worldEvents.addAll([
      spawn('s-1', 'ashen_colossus', 'world_boss', 'active', 1, const Duration(minutes: -10)),
      spawn('s-2', 'greenwood_incursion', 'invasion', 'scheduled', 30, const Duration(hours: 5))
        ..['region_code'] = 'greenwood',
      spawn('s-3', 'ashen_colossus', 'world_boss', 'succeeded', 1, const Duration(days: -1))
        ..addAll({'enlisted': true, 'my_power': 300, 'my_score': 1800, 'my_reward_tier': 2}),
    ]);
    await tester.tap(find.text('Adventure'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Events'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Claim 15'));
    await tester.pumpAndSettle();
    expect(find.text('+15 event tokens'), findsOneWidget);
    expect(find.text('Claimed'), findsOneWidget);
    expect(find.text('15 event tokens'), findsOneWidget);

    expect(find.textContaining('Invasion · Greenwood only · Starts in'), findsOneWidget);
    final enlist = find.widgetWithText(FilledButton, 'Enlist');
    expect(enlist, findsNWidgets(2));
    expect(tester.widget<FilledButton>(enlist.last).onPressed, isNull, reason: 'level 30 gate');

    await tester.tap(enlist.first);
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(of: find.byType(AlertDialog), matching: find.text('Enlist')));
    await tester.pumpAndSettle();
    expect(be.worldEvents.first['enlisted'], isTrue);
    expect(find.text('Enlisted · power 420'), findsOneWidget);

    await tester.scrollUntilVisible(find.text('Your score 1800 · reward tier 2'), 200, scrollable: _list);
    expect(find.textContaining('World boss · Victory'), findsOneWidget);
  });

  testWidgets('Turkish locale renders translated copy', (tester) async {
    await _boot(tester, FakeBackend(), locale: const Locale('tr'));
    expect(find.text('Giriş yap'), findsOneWidget);
    expect(find.text('E-posta veya kullanıcı adı'), findsOneWidget);
  });

  testWidgets('Turkish in-game chrome fits a phone (4 adventure segments, progress button)', (tester) async {
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
    await _boot(tester, be, store: store, locale: const Locale('tr'));
    await tester.tap(find.text('Grim'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Macera'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Etkinlik'));
    await tester.pumpAndSettle();
    expect(find.text('Şu an planlı dünya etkinliği yok.'), findsOneWidget); // layout overflow would fail the test
  });

  test('splitCopper', () {
    expect(splitCopper(123456), (gold: 12, silver: 34, copper: 56));
    expect(splitCopper(-5), (gold: 0, silver: 0, copper: 0));
  });
}
