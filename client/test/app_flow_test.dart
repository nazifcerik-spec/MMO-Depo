import 'dart:convert';

import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/adventure/adventure_models.dart';
import 'package:otter/shared/content_names.dart';
import 'package:otter/features/character/class_catalog.dart';
import 'package:otter/app/app.dart';
import 'package:otter/features/activity/activity_view.dart';
import 'package:otter/features/activity/activity_models.dart';
import 'package:otter/features/inventory/item_images.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/auth/auth_controller.dart';
import 'package:otter/features/character/models.dart';
import 'package:otter/shared/links.dart';
import 'package:otter/shared/loot_sheet.dart';

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
// The page's vertical list (a text field inside it scrolls sideways).
final _list = find
    .descendant(
      of: find.byType(ListView),
      matching: find.byWidgetPredicate((w) => w is Scrollable && w.axisDirection == AxisDirection.down),
    )
    .last;

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

    await tester.enterText(find.byType(TextFormField).at(1), 'busy-horse');
    await tester.tap(find.text('Log in'));
    await tester.pumpAndSettle();
    expect(find.text('The server is busy right now. Try again in a few seconds.'), findsOneWidget);

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
    await tester.pumpAndSettle();
    expect(find.text('Terms of Service'), findsOneWidget, reason: 'policy links next to the checkbox');
    await tester.ensureVisible(find.text('Create account'));
    await tester.tap(find.text('Create account'));
    await tester.pumpAndSettle();
    expect(find.text('Characters'), findsOneWidget);
    final reg = be.requests.firstWhere((r) => r.url.path == '/v1/auth/register');
    expect(reg.body, contains('"accept_terms":true'));
  });

  testWidgets('creation screen explains the class: traits, skill bubble, pyramid path, node skills', (tester) async {
    tester.view.physicalSize = const Size(360 * 3, 780 * 3); // phone width: the four-node pyramid row must fit
    tester.view.devicePixelRatio = 3;
    addTearDown(tester.view.reset);
    ClassCatalog.reset();
    ContentNames.use({
      'class.warrior': 'Warrior',
      'class.warrior.role': 'Melee / Physical',
      'class.warrior.intro': 'A heavily armored frontline fighter.',
      'class.warrior.tagline': 'Master of steel and fury',
      'class.warrior.quote': "Shields break. Warriors don't.",
      'class.guardian': 'Guardian',
      'class.executioner': 'Executioner',
      'class.iron_fortress': 'Iron Fortress',
      'class.warlord': 'Warlord',
      'class.berserker': 'Berserker',
      'class.weapon_master': 'Weapon Master',
      'resource.rage': 'Rage',
      'attribute.strength': 'Strength',
      'attribute.vitality': 'Vitality',
      'equipment.heavy_armor': 'Heavy armor',
      'skill.warrior_strike': 'Strike',
      'skill.warrior_strike.desc': 'Basic attack that builds Rage.',
      'skill.iron_fortress_unbreakable': 'Unbreakable',
      'skill.iron_fortress_unbreakable.desc': '+150% armor for 3 turns.',
      'title.the_unbreakable': 'The Unbreakable',
    });
    addTearDown(() {
      ContentNames.use(const {});
      ClassCatalog.reset();
    });
    final be = FakeBackend()
      ..classCatalog = {
        'attribute_points_total': 3022,
        'milestones': {'awakening': 600, 'ultimate': 700, 'capstone': 850, 'mastery': 1000},
        'classes': [
          {
            'code': 'warrior',
            'tier': 0,
            'resource': {'kind': 'rage', 'max': 100, 'start_pct': 0, 'per_hit': 12, 'per_turn': 0},
            'base_stats': {'strength': 10, 'vitality': 10},
            'primary_attribute': 'strength',
            'secondary_attributes': ['vitality'],
            'equipment': ['heavy_armor'],
          },
          for (final (code, parent, lvl) in [
            ('guardian', 'warrior', 100),
            ('executioner', 'warrior', 100),
            ('iron_fortress', 'guardian', 300),
            ('warlord', 'guardian', 300),
            ('berserker', 'executioner', 300),
            ('weapon_master', 'executioner', 300),
          ])
            {
              'code': code,
              'tier': lvl == 100 ? 1 : 2,
              'parent': parent,
              'required_level': lvl,
              'equipment': ['heavy_armor'],
              if (code == 'iron_fortress') 'mastery_title': 'the_unbreakable',
            },
        ],
      }
      ..warriorSkills = [
        {
          'code': 'warrior_strike',
          'kind': 'active',
          'unlock_level': 1,
          'effect': 'strike',
          'power': 110,
          'mp_cost': 0,
          'cooldown_turns': 1,
        },
        {
          'code': 'iron_fortress_unbreakable',
          'spec': 'iron_fortress',
          'kind': 'active',
          'unlock_level': 700,
          'effect': 'guard',
          'power': 150,
          'mp_cost': 50,
          'cooldown_turns': 20,
          'config': {'category': 'ultimate', 'duration_turns': 3},
        },
      ];
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await _boot(tester, be, store: store);
    await tester.tap(find.text('New character'));
    await tester.pumpAndSettle();

    // Hero (DD-266): tagline under the name, role chips, the saying under all.
    expect(find.text('Master of steel and fury'), findsOneWidget);
    expect(find.text('Melee'), findsOneWidget);
    expect(find.text('Physical'), findsOneWidget);
    expect(find.text('A heavily armored frontline fighter.'), findsOneWidget);
    expect(find.text('"Shields break. Warriors don\'t."'), findsOneWidget);

    final list = find
        .descendant(of: find.byKey(const ValueKey('create-character-list')), matching: find.byType(Scrollable))
        .first;
    await tester.scrollUntilVisible(find.text('Difficulty'), 200, scrollable: list);
    expect(find.text('Resource: Rage'), findsOneWidget);
    expect(find.text('Starts empty · +12 per hit'), findsOneWidget);
    expect(find.text('Main attribute: Strength · second: Vitality'), findsOneWidget);
    expect(find.text('Heavy armor'), findsOneWidget);
    expect(find.text('Easy — forgiving, safe on AFK hunts'), findsOneWidget);

    // Skills: base only; a tap opens the bubble, a second tap closes it.
    await tester.scrollUntilVisible(find.text('Strike'), 200, scrollable: list);
    await tester.ensureVisible(find.text('Strike'));
    await tester.pumpAndSettle();
    expect(find.text('Unbreakable'), findsNothing);
    final bubble = find.textContaining('Basic attack that builds Rage.', findRichText: true);
    expect(bubble, findsNothing);
    await tester.tap(find.text('Strike'));
    await tester.pumpAndSettle();
    expect(bubble, findsOneWidget);
    expect(find.textContaining('active · Lvl 1', findRichText: true), findsOneWidget);
    await tester.tap(find.text('Strike'));
    await tester.pumpAndSettle();
    expect(bubble, findsNothing);

    await tester.scrollUntilVisible(find.text('Weapon Master'), 200, scrollable: list);
    await tester.ensureVisible(find.text('Iron Fortress'));
    await tester.pumpAndSettle();
    for (final n in ['Guardian', 'Executioner', 'Iron Fortress', 'Warlord', 'Berserker', 'Weapon Master']) {
      expect(find.text(n), findsOneWidget, reason: 'pyramid node $n');
    }
    await tester.tap(find.text('Iron Fortress'));
    await tester.pumpAndSettle();
    expect(find.text('Unbreakable'), findsOneWidget);
    expect(find.text('ultimate'), findsOneWidget);
    expect(find.text('Cost 50 · cooldown 20 turns'), findsOneWidget);
    expect(find.text('Lv 1000 mastery: The Unbreakable'), findsOneWidget);
    await tester.tap(find.text('Close'));
    await tester.pumpAndSettle();
    expect(find.text('Unbreakable'), findsNothing);
  });

  testWidgets('create a character, enter it, and read the sheet', (tester) async {
    final be = FakeBackend();
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse'); // persisted session → restore on boot
    await _boot(tester, be, store: store);
    expect(find.text('Characters'), findsOneWidget);

    await tester.tap(find.text('New character'));
    await tester.pumpAndSettle();
    // Class strip, race cards, then the name (DD-266).
    final strip = find
        .descendant(of: find.byKey(const ValueKey('class-strip')), matching: find.byType(Scrollable))
        .first;
    await tester.scrollUntilVisible(find.text('Shaman'), 200, scrollable: strip);
    expect(find.text('Shaman'), findsOneWidget, reason: 'all ten classes are offered');
    await tester.scrollUntilVisible(find.text('Ranger'), -200, scrollable: strip);
    await tester.tap(find.text('Ranger').first);
    await tester.pumpAndSettle();
    final list = find
        .descendant(of: find.byKey(const ValueKey('create-character-list')), matching: find.byType(Scrollable))
        .first;
    await tester.scrollUntilVisible(find.text('Gnome'), 200, scrollable: list);
    for (final r in ['Human', 'Elf', 'Dwarf', 'Orc', 'Troll', 'Undead', 'Dark Elf', 'Drakari', 'Tauren', 'Gnome']) {
      expect(find.text(r), findsOneWidget, reason: 'race $r (DD-267)');
    }
    // One size for every card, however long the name or lore.
    final cards = find.byWidgetPredicate((w) => w.runtimeType.toString() == '_RaceCard');
    final sizes = {for (final e in cards.evaluate()) tester.getSize(find.byWidget(e.widget))};
    expect(sizes, hasLength(1), reason: '$sizes');
    expect(find.text('Your race sets your look; it does not change attributes.'), findsOneWidget);
    await tester.ensureVisible(find.text('Dark Elf'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Dark Elf'));
    await tester.pump();
    await tester.tap(find.text('Continue as Ranger'));
    await tester.pumpAndSettle();
    expect(find.text('Name your character'), findsOneWidget);
    await tester.enterText(find.byType(TextFormField), '1bad');
    await tester.tap(find.text('Create'));
    await tester.pump();
    expect(find.text('3–16 characters, starts with a letter (A–Z, 0–9, _)'), findsWidgets);

    await tester.enterText(find.byType(TextFormField), 'Grim');
    await tester.tap(find.text('Create'));
    await tester.pumpAndSettle();
    expect(be.characters.single['race_code'], 'dark_elf');
    expect(find.text('Lvl 1 · Ranger · Dark Elf'), findsOneWidget);
    expect(find.text('1/2 slots'), findsOneWidget);

    await tester.tap(find.text('Grim'));
    await tester.pumpAndSettle();
    expect(find.text('80/120'), findsOneWidget); // HP
    expect(find.text('FOCUS'), findsOneWidget, reason: 'the class resource replaces the MP bar');
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

  testWidgets('delete a character at once after typing its name; a pending one can be restored', (tester) async {
    final be = FakeBackend()
      ..characters.addAll([
        {'id': 'char-1', 'name': 'Grim', 'race_code': 'orc', 'class_code': 'warrior', 'level': 7, 'state': 'active'},
        {
          'id': 'char-2',
          'name': 'Vex',
          'race_code': 'elf',
          'class_code': 'mage',
          'level': 3,
          'state': 'delete_grace', // left by an account deletion (DD-231)
          'delete_grace_at': DateTime.now().toUtc().add(const Duration(days: 20)).toIso8601String(),
        },
      ]);
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await _boot(tester, be, store: store);

    expect(find.byTooltip('Delete character'), findsOneWidget, reason: 'only the active character');
    await tester.tap(find.byTooltip('Delete character'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Grim will be deleted right away'), findsOneWidget);
    FilledButton confirm() => tester.widget<FilledButton>(find.byKey(const Key('delete-character-confirm')));
    expect(confirm().onPressed, isNull, reason: 'off until the name is typed');
    await tester.enterText(find.byKey(const Key('delete-character-name')), 'Gri');
    await tester.pump();
    expect(confirm().onPressed, isNull);
    await tester.enterText(find.byKey(const Key('delete-character-name')), ' grim ');
    await tester.pump();
    expect(confirm().onPressed, isNotNull, reason: 'case and spaces do not matter');
    await tester.tap(find.byKey(const Key('delete-character-confirm')));
    await tester.pumpAndSettle();
    final sent = be.requests.lastWhere((r) => r.url.path.endsWith('/delete'));
    expect(sent.url.path, '/v1/character/char-1/delete');
    expect(jsonDecode(sent.body), {'confirm_name': 'grim'});
    expect(be.characters.map((c) => c['id']), ['char-2']);
    expect(find.text('Grim was deleted.'), findsOneWidget);
    expect(find.text('Grim'), findsNothing, reason: 'gone from the list at once');

    expect(find.text('Deletion pending'), findsOneWidget);
    expect(find.textContaining('Deleted on '), findsOneWidget);
    await tester.tap(find.text('Vex'));
    await tester.pumpAndSettle();
    expect(find.text('Characters'), findsOneWidget, reason: 'a character pending deletion cannot be entered');

    await tester.tap(find.text('Cancel deletion'));
    await tester.pumpAndSettle();
    expect(be.characters.single['state'], 'active');
    expect(find.text('Vex will not be deleted.'), findsOneWidget);
    expect(find.byTooltip('Delete character'), findsOneWidget);
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
    await tester.pumpAndSettle();
    // Gathering profession (DD-244): level, the unlocked resources (best
    // first, preselected) and the next one still locked.
    expect(find.text('Fishing · Lvl 20'), findsOneWidget);
    expect(find.text('10 / 64 XP'), findsOneWidget);
    expect(find.text('Unlocks at level 70'), findsOneWidget);
    await tester.tap(find.text('Salmon'));
    await _gatherScrollTo(tester, '3h 00m');
    await tester.tap(find.text('3h 00m'));
    await tester.pump();
    await _gatherScrollTo(tester, 'Start Fishing');
    await tester.tap(find.text('Start Fishing'));
    await tester.pumpAndSettle();
    expect(be.activeActivity!['kind'], 'fishing');
    expect(be.activeActivity!['target'], 'salmon');
    expect(be.activeActivity!['duration_seconds'], 10800);
    expect(find.textContaining('left'), findsOneWidget);

    final sheetsBefore = be.sheetGets;
    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    expect(find.text('Collect early?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last);
    await tester.pumpAndSettle();

    expect(find.text('Copper Ore'), findsOneWidget);
    expect(find.text('Flawless Raw Sapphire'), findsOneWidget);
    expect(find.text('Fishing +36 XP · Fishing reached level 21!'), findsOneWidget);
    await tester.tap(find.text('Nice!'));
    await tester.pumpAndSettle();
    expect(find.text('Choose an activity'), findsOneWidget);
    expect(be.activeActivity, isNull);
    expect(be.sheetGets, greaterThan(sheetsBefore), reason: 'reveal bumps other tabs');
  });

  testWidgets('daily reward: shown on entering while unclaimed, claim once (DD-311)', (tester) async {
    final be = FakeBackend()
      ..dailyClaimable = true
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
    expect(find.text('Daily login reward'), findsOneWidget);
    expect(find.text('Day 7'), findsOneWidget);
    expect(find.text('3× Potion Hp Small'), findsOneWidget);
    await tester.tap(find.byKey(const Key('daily-claim')));
    await tester.pumpAndSettle();
    expect(be.dailyClaims, 1);
    expect(find.text('Come back tomorrow'), findsOneWidget);
    await tester.tap(find.text('Close'));
    await tester.pumpAndSettle();
    expect(find.text('Daily login reward'), findsNothing);
  });

  testWidgets('gathering restart: collect so far, the same run starts again (DD-310)', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.text('Gather'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Fishing'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Salmon'));
    await _gatherScrollTo(tester, 'Start Fishing');
    await tester.tap(find.text('Start Fishing'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Restart'));
    await tester.pumpAndSettle();
    expect(find.text('Collect and start again?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Restart'));
    await tester.pumpAndSettle();
    expect(be.restarts, 1);
    expect(find.text('Copper Ore'), findsOneWidget, reason: 'the haul so far');
    await tester.tap(find.text('Nice!'));
    await tester.pumpAndSettle();
    expect(be.activeActivity!['id'], 'act-2');
    expect(be.activeActivity!['target'], 'salmon');
    expect(find.text('Restart'), findsOneWidget, reason: 'the new run is on screen');
  });

  testWidgets('hunting resources show the animal when its picture ships', (tester) async {
    ItemImages.animals = {'forest_rat'};
    addTearDown(() => ItemImages.animals = const {});
    await _inGame(tester);
    await tester.tap(find.text('Gather'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Hunting'));
    await tester.pumpAndSettle();
    expect(huntAnimal('forest_rat_pelt'), 'forest_rat');
    expect(huntAnimal('dragon_scale'), 'dragon');
    final img = tester.widget<Image>(find.byType(Image).first);
    expect(((img.image as ResizeImage).imageProvider as AssetImage).assetName, 'assets/animals/forest_rat.png');
  });

  testWidgets('abandon asks first and forfeits', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.text('Gather'));
    await tester.pumpAndSettle();
    await _gatherScrollTo(tester, 'Start Mining');
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

  testWidgets('bag: item detail, equip/unequip, class gate, sell, destroy, sort, multi-select', (tester) async {
    final semantics = tester.ensureSemantics();
    final (be, ws, _) = await _inGame(tester);
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
    expect(find.text('That slot is empty: equipping adds all of these.'), findsOneWidget);
    final sheets = be.sheetGets;
    await tester.tap(find.text('Equip'));
    await tester.pumpAndSettle();
    expect(be.equipped['mainhand']?['instance_id'], 's1');
    expect(be.sheetGets, greaterThan(sheets), reason: 'gear change reloads the sheet');

    // Another main-hand piece is compared with the worn sword, stat by stat.
    await tester.tap(find.bySemanticsLabel('Staff Oak ×1'));
    await tester.pumpAndSettle();
    expect(find.text('Compared with what you wear'), findsOneWidget);
    expect(find.text('Worn: Sharp Sword Iron'), findsOneWidget);
    expect(find.text('-10'), findsOneWidget, reason: 'attack lost');
    expect(find.text('-10%'), findsOneWidget, reason: 'STR lost');
    Navigator.of(tester.element(find.text('Compared with what you wear'))).pop(); // close the sheet
    await tester.pumpAndSettle();

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

    // Sort by name (display only): Boar Hide moves to the first cell.
    be.bag.addAll([item('b1', 'boar_hide', 20), item('w1', 'wolf_pelt', 21)]);
    ws.last.push('hunt.reveal', {}); // any data bump reloads the bag
    await tester.pumpAndSettle();
    await tester.tap(find.text('Bag order'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Name').last);
    await tester.pumpAndSettle();
    Offset at(String label) => tester.getCenter(find.bySemanticsLabel(label));
    expect(at('Boar Hide ×1').dx, lessThan(at('Sharp Sword Iron ×1').dx));
    expect(at('Boar Hide ×1').dy, at('Sharp Sword Iron ×1').dy);

    // Long press, then drag across: Boar Hide, Sharp Sword Iron, Staff Oak.
    final g = await tester.startGesture(at('Boar Hide ×1'));
    await tester.pump(kLongPressTimeout + const Duration(milliseconds: 50));
    await g.moveTo(at('Sharp Sword Iron ×1'));
    await g.moveTo(at('Staff Oak ×1'));
    await g.up();
    await tester.pumpAndSettle();
    expect(find.text('3 selected'), findsOneWidget);
    // A tap toggles while selecting (Wolf Pelt in, then out again).
    await tester.tap(find.bySemanticsLabel('Wolf Pelt ×1'));
    await tester.pump();
    expect(find.text('4 selected'), findsOneWidget);
    await tester.tap(find.bySemanticsLabel('Wolf Pelt ×1'));
    await tester.pump();

    await tester.tap(find.byTooltip('Sell'));
    await tester.pumpAndSettle();
    expect(find.text('Sell 3 items?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Sell'));
    await tester.pumpAndSettle();
    expect(be.bag.map((i) => i['instance_id']), ['w1']);
    expect(find.textContaining('Sold for'), findsOneWidget);
    expect(find.textContaining('selected'), findsNothing, reason: 'selection ends after the sale');
    semantics.dispose();
  });

  testWidgets('bag: gear the class cannot wear stands out, with its kind and why', (tester) async {
    final (be, _, _) = await _inGame(tester);
    Map<String, dynamic> gear(
      String id,
      String code,
      int slot,
      String sub,
      String equip,
      String wear, {
      int level = 1,
    }) => {
      'instance_id': id,
      'template_code': code,
      'item_class': 'armor',
      'item_subclass': sub,
      'equip_slot': equip,
      'slot': slot,
      'stack_quantity': 1,
      'max_stack': 1,
      'rarity': 'common',
      'min_level': level,
      'wear': wear,
    };
    be.bag.addAll([
      gear('a', 'chest_novice', 0, 'plate_armor', 'chest', 'ok'),
      gear('b', 'robe_novice', 1, 'cloth_armor', 'chest', 'class'),
      gear('c', 'chest_adamant', 2, 'plate_armor', 'chest', 'level', level: 50),
    ]);
    await tester.tap(find.byIcon(Icons.backpack_outlined));
    await tester.pumpAndSettle();
    // A legend explains the marks; the cloth robe is dimmed.
    expect(find.text('not for your class'), findsOneWidget);
    expect(find.text('not yet (level or attributes)'), findsOneWidget);
    expect(find.byWidgetPredicate((w) => w is Opacity && w.opacity < 1), findsOneWidget, reason: 'only the robe');

    await tester.tap(find.bySemanticsLabel("Robe Novice ×1 · Your class can't wear this"));
    await tester.pumpAndSettle();
    expect(find.textContaining('Chest · Cloth Armor'), findsOneWidget, reason: 'slot and armor kind');
    expect(find.text("Your class can't wear this"), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Equip')).onPressed, isNull);
    Navigator.of(tester.element(find.textContaining('Chest · Cloth Armor'))).pop();
    await tester.pumpAndSettle();

    await tester.tap(find.bySemanticsLabel('Chest Adamant ×1 · Needs level 50'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Chest · Plate Armor'), findsOneWidget);
    expect(find.text('Needs level 50'), findsOneWidget);
    Navigator.of(tester.element(find.textContaining('Chest · Plate Armor'))).pop();
    await tester.pumpAndSettle();

    await tester.tap(find.bySemanticsLabel('Chest Novice ×1'));
    await tester.pumpAndSettle();
    expect(find.text('Your class can wear this'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Equip')).onPressed, isNotNull);
  });

  testWidgets('bag: an enchanting scroll goes on fitting gear, replacing an old enchant', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.bag.addAll([
      {
        'instance_id': 'scroll',
        'template_code': 'enchant_weapon_attack_1',
        'item_class': 'consumable',
        'item_subclass': 'enchant_scroll',
        'slot': 0,
        'stack_quantity': 2,
        'max_stack': 20,
        'rarity': 'common',
        'enchant': {
          'slots': ['weapon'],
          'modifiers': [
            {'kind': 'mult_add', 'target': 'attack', 'amount': 0.01},
          ],
        },
      },
      {
        'instance_id': 'sword',
        'template_code': 'sword_training',
        'item_class': 'weapon',
        'item_subclass': 'sword',
        'equip_slot': 'mainhand',
        'slot': 1,
        'stack_quantity': 1,
        'max_stack': 1,
        'rarity': 'common',
        'enchant_code': 'enchant_weapon_crit_1',
        'enchant_effect': {
          'slots': ['weapon'],
          'modifiers': [
            {'kind': 'flat', 'target': 'crit_pct', 'amount': 1},
          ],
        },
      },
      {
        'instance_id': 'helm',
        'template_code': 'helm_novice',
        'item_class': 'armor',
        'item_subclass': 'plate_armor',
        'equip_slot': 'helm',
        'slot': 2,
        'stack_quantity': 1,
        'max_stack': 1,
        'rarity': 'common',
      },
    ]);
    await tester.tap(find.byIcon(Icons.backpack_outlined));
    await tester.pumpAndSettle();

    await tester.tap(find.bySemanticsLabel('Enchant Weapon Attack 1 ×2'));
    await tester.pumpAndSettle();
    expect(find.text('Goes on: weapons\n+1% Attack'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Enchant'));
    await tester.pumpAndSettle();
    // Only the sword fits a weapon scroll; its current enchant is shown.
    expect(find.text('Choose the gear to enchant'), findsOneWidget);
    expect(find.text('Sword Training'), findsOneWidget);
    expect(find.text('Helm Novice'), findsNothing);
    // The current enchant says what it adds (DD-295).
    expect(find.text('Enchant: Enchant Weapon Crit 1 · +1 Crit'), findsOneWidget);
    await tester.tap(find.text('Sword Training'));
    await tester.pumpAndSettle();
    expect(find.text('Replace the enchant?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Enchant'));
    await tester.pumpAndSettle();
    expect(find.text('Sword Training enchanted'), findsOneWidget);
    expect(be.bag.firstWhere((i) => i['instance_id'] == 'sword')['enchant_code'], 'enchant_weapon_attack_1');
    expect(be.bag.firstWhere((i) => i['instance_id'] == 'scroll')['stack_quantity'], 1);

    await tester.tap(find.bySemanticsLabel('Sword Training ×1'));
    await tester.pumpAndSettle();
    expect(find.text('Enchant: Enchant Weapon Attack 1 · +1% Attack'), findsOneWidget);
  });

  testWidgets('hunt: level-aware default, start, early collect, result with level-up', (tester) async {
    final (be, _, _) = await _inGame(tester);
    Map<String, dynamic> consumable(String id, String code, int slot, int qty, Map<String, dynamic> effect) => {
      'instance_id': id,
      'template_code': code,
      'item_class': 'consumable',
      'item_subclass': code.startsWith('food') ? 'food' : 'potion',
      'slot': slot,
      'stack_quantity': qty,
      'max_stack': 20,
      'rarity': 'common',
      'use_effect': effect,
    };
    be.bag.addAll([
      consumable('p1', 'potion_hp_small', 0, 3, {'kind': 'heal_pct', 'amount': 20}),
      consumable('p2', 'potion_hp_small', 1, 2, {'kind': 'heal_pct', 'amount': 20}),
      consumable('m1', 'potion_mp_small', 2, 4, {'kind': 'mana_pct', 'amount': 30}),
      consumable('f1', 'food_grilled_guppy', 3, 2, {
        'kind': 'heal_pct',
        'amount': 15,
        'buff': {'attack_pct': 2, 'defense_pct': 2, 'minutes': 60},
      }),
      {
        'instance_id': 'ore',
        'template_code': 'copper_ore',
        'item_class': 'material',
        'item_subclass': 'ore',
        'slot': 4,
        'stack_quantity': 9,
        'max_stack': 999,
        'rarity': 'common',
      },
    ]);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    expect(find.text('Choose your prey'), findsOneWidget);
    // The hunt panel's own list (the page around it holds other scrollables),
    // found from a widget on screen now.
    Finder panelOf(Finder f) =>
        find.byWidget(tester.widget<Scrollable>(find.ancestor(of: f, matching: find.byType(Scrollable)).first));
    Future<void> startHunt() async {
      final anchor = [
        'Choose your prey',
        'Duration',
        'Supplies to take',
      ].map(find.text).firstWhere((f) => f.evaluate().isNotEmpty);
      await tester.scrollUntilVisible(find.text('Start hunt'), 200, scrollable: panelOf(anchor));
      await tester.ensureVisible(find.text('Start hunt'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Start hunt'));
      await tester.pumpAndSettle();
    }

    expect(find.text('Cave Troll'), findsOneWidget);
    // Only mobs in reach are huntable; the rest are listed by region (DD-273).
    expect(find.textContaining('Hunting from here:'), findsOneWidget);
    expect(find.text('Sand Wyrm'), findsNothing);
    await tester.scrollUntilVisible(
      find.text('Monsters in other regions'),
      200,
      scrollable: panelOf(find.text('Choose your prey')),
    );
    await tester.ensureVisible(find.text('Monsters in other regions')); // clear of the bottom bar
    await tester.pumpAndSettle();
    await tester.tap(find.text('Monsters in other regions'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.text('Salt Desert'));
    await tester.pumpAndSettle();
    expect(find.text('Salt Desert'), findsOneWidget);
    expect(find.textContaining('Sand Wyrm'), findsOneWidget);
    await tester.scrollUntilVisible(find.text('Choose your prey'), -200, scrollable: panelOf(find.text('Salt Desert')));

    // Default = toughest mob at or below level 3.
    await startHunt();
    expect(be.activeHunt!['mob_code'], 'grey_wolf');
    await tester.tap(find.text('Abandon'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Abandon'));
    await tester.pumpAndSettle();
    expect(be.activeHunt, isNull);

    await tester.tap(find.text('Giant Rat'));
    // AFK lengths are 1 / 3 / 6 hours (DD-249).
    for (final d in ['1h 00m', '3h 00m', '6h 00m']) {
      expect(find.text(d), findsOneWidget);
    }
    expect(find.text('30m'), findsNothing);
    await tester.tap(find.text('6h 00m'));
    await tester.pump();
    // Supplies from the bag (DD-288): stacks add up, materials are not offered.
    await tester.scrollUntilVisible(find.text('Food'), 200, scrollable: panelOf(find.text('6h 00m')));
    await tester.ensureVisible(find.byKey(const Key('supply-food_grilled_guppy')));
    await tester.pumpAndSettle();
    expect(find.text('Supplies to take'), findsOneWidget);
    expect(find.text('Healing potions'), findsOneWidget);
    expect(find.text('Mana potions'), findsOneWidget);
    expect(find.text('Refills 30% of your mana'), findsOneWidget);
    expect(find.byKey(const Key('supply-potion_hp_small')), findsOneWidget);
    expect(find.byKey(const Key('supply-copper_ore')), findsNothing);
    expect(find.text('/ 5'), findsOneWidget, reason: 'stacks add up');
    // Typed counts; more than the stock snaps to the stock (DD-292).
    await tester.enterText(find.byKey(const Key('supply-potion_hp_small')), '9');
    await tester.pump();
    expect(
      tester
          .widget<TextField>(
            find.descendant(of: find.byKey(const Key('supply-potion_hp_small')), matching: find.byType(TextField)),
          )
          .controller!
          .text,
      '5',
    );
    await tester.enterText(find.byKey(const Key('supply-potion_hp_small')), '2');
    await tester.enterText(find.byKey(const Key('supply-food_grilled_guppy')), '1');
    await tester.pump();
    await startHunt();
    expect(be.activeHunt!['mob_code'], 'giant_rat');
    expect(be.lastHuntSupplies, [
      {'template_code': 'potion_hp_small', 'quantity': 2},
      {'template_code': 'food_grilled_guppy', 'quantity': 1},
    ]);

    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last);
    await tester.pumpAndSettle();
    expect(find.text('7 won · 1 lost · Used: 2× Potion Hp Small · Level up! Now level 4'), findsOneWidget);
    expect(find.text('+160 XP'), findsOneWidget);
    expect(find.text('Rare × 1'), findsOneWidget);
    // The fight report (DD-308): why it ended and how the sides measured up.
    final report = find.byKey(const Key('fight-report'));
    final sheet = find.descendant(of: find.byType(LootSheet), matching: find.byType(Scrollable)).first;
    await tester.scrollUntilVisible(report, 200, scrollable: sheet);
    expect(find.text('HP fell below 5% against Cave Troll: you broke off and came back.'), findsOneWidget);
    expect(find.descendant(of: report, matching: find.text('87%')), findsOneWidget);
    expect(find.text('Your swings: 35/40 hits (88%), avg 60, max 88'), findsOneWidget);
    expect(find.text('Enemy swings: 30/38 hits (79%), avg 30, max 45'), findsOneWidget);
    expect(find.text('Lowest HP in a fight: 3%'), findsOneWidget);
    expect(find.textContaining('Camped: '), findsOneWidget);
    await tester.scrollUntilVisible(find.text('Nice!'), 200, scrollable: sheet);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Nice!'));
    await tester.pumpAndSettle();
    expect(find.text('Choose your prey'), findsOneWidget);
  });

  testWidgets('dungeon: level gate, collect locked until the run ends, cleared result', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Dungeon'));
    await tester.pumpAndSettle();

    // Only dungeons entered from here (DD-302); each says what it pays this
    // character (DD-303).
    expect(find.text('Far Keep'), findsNothing);
    expect(find.text('≈ 6000 XP for you'), findsOneWidget);
    final enters = find.widgetWithText(FilledButton, 'Enter');
    expect(tester.widget<FilledButton>(enters.at(1)).onPressed, isNull, reason: 'Sunken Crypt needs level 18');

    // A level-up elsewhere (here: a hunt revealed on the server) unlocks it
    // at once — the sheet reload carries the new level to every tab.
    be.characters.first['level'] = 18;
    ws.last.push('hunt.reveal', {'leveled_up': true});
    await tester.pumpAndSettle();
    expect(tester.widget<FilledButton>(enters.at(1)).onPressed, isNotNull);
    be.characters.first['level'] = 3;
    ws.last.push('hunt.reveal', {});
    await tester.pumpAndSettle();

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

  testWidgets('raid: leader opens a lobby, members join live, start with confirm, collect at the end', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Dungeon'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Raids'));
    await tester.pumpAndSettle();
    expect(find.text('Party of 3–5 · 1h 00m'), findsOneWidget);
    expect(find.text('Boss: 2 gear pieces guaranteed'), findsOneWidget);

    await tester.tap(find.text('Open raid'));
    await tester.pumpAndSettle();
    expect(find.text('Members (1/5)'), findsOneWidget);
    expect(find.text('At least 3 members to start'), findsOneWidget);
    final start = find.widgetWithText(FilledButton, 'Start raid');
    expect(tester.widget<FilledButton>(start).onPressed, isNull, reason: 'not enough members yet');

    // Two party members join on their devices: the push reloads the lobby.
    (be.raidRun!['members'] as List).addAll(<Map<String, dynamic>>[
      be.raidMember('c2', 'Ada'),
      be.raidMember('c3', 'Bo'),
    ]);
    ws.last.push('raid.changed', {'raid_id': 'raid-1', 'state': 'forming'});
    await tester.pumpAndSettle();
    expect(find.text('Members (3/5)'), findsOneWidget);
    await tester.tap(start);
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Start raid').last); // the confirm dialog
    await tester.pumpAndSettle();
    expect(be.raidRun!['state'], 'running');
    expect(find.text('Rewards open when the raid ends; every member collects their own.'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Collect')).onPressed, isNull);

    be.raidRun!['expires_at'] = DateTime.now().toUtc().toIso8601String();
    await tester.fling(find.byType(ListView).last, const Offset(0, 400), 1000);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Collect'));
    await tester.pumpAndSettle();
    expect(find.text('Cleared! · Boss defeated · 9 won · 0 lost'), findsOneWidget);
    Navigator.of(tester.element(find.text('Cleared! · Boss defeated · 9 won · 0 lost'))).pop(); // close the loot sheet
    await tester.pumpAndSettle();
    expect(find.text('Cleared this week'), findsOneWidget);
    expect(find.text('Open raid'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Open raid')).onPressed, isNull);

    // Out of a party there is nothing to open.
    be.raidInParty = false;
    await tester.fling(find.byType(ListView).last, const Offset(0, 400), 1000);
    await tester.pumpAndSettle();
    expect(
      find.text(
        'Raids are for a party. Queue below for a random one, join a listing on the board, or advertise your own.',
      ),
      findsOneWidget,
    );
    expect(find.text('Open raid'), findsNothing);
  });

  testWidgets('raid finding: queue and leave it, join a listed party, advertise as leader (DD-299)', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.raidInParty = false;
    be.raidBoard.add({
      'id': 'l-1',
      'raid_code': 'north_keep_siege',
      'members': [
        {'name': 'Aria', 'level': 120, 'class_code': 'cleric'},
      ],
    });
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Dungeon'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Raids'));
    await tester.pumpAndSettle();

    // Out of a party: queue or advertise the chosen raid; the board lists parties.
    expect(find.text('Find a party'), findsOneWidget);
    expect(find.text('Listing board'), findsOneWidget);
    expect(find.text('Aria · Lvl 120 · Cleric'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Queue'));
    await tester.pumpAndSettle();
    expect(be.raidCalls, ['queue']);
    expect(find.text('Queued for North Keep Siege · 1 waiting'), findsOneWidget);
    await tester.tap(find.widgetWithText(OutlinedButton, 'Leave queue'));
    await tester.pumpAndSettle();
    expect(find.text('Queued for North Keep Siege · 1 waiting'), findsNothing);

    await tester.ensureVisible(find.widgetWithText(FilledButton, 'Join party'));
    await tester.tap(find.widgetWithText(FilledButton, 'Join party'));
    await tester.pumpAndSettle();
    expect(be.raidCalls.last, 'listings/l-1/join');
    expect(find.text('Listing board'), findsNothing, reason: 'in a party now');

    // A leader advertises; the card then offers to take it down.
    await tester.ensureVisible(find.widgetWithText(FilledButton, 'Advertise'));
    await tester.tap(find.widgetWithText(FilledButton, 'Advertise'));
    await tester.pumpAndSettle();
    expect(find.text('Your party is on the board for North Keep Siege'), findsOneWidget);
    await tester.tap(find.widgetWithText(OutlinedButton, 'Remove listing'));
    await tester.pumpAndSettle();
    expect(be.raidCalls.last, 'listing/close');
  });

  testWidgets('travel: pick any destination on the map, level gate, depart, turn back (DD-279, DD-280)', (
    tester,
  ) async {
    final (be, _, _) = await _inGame(tester);
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Travel'));
    await tester.pumpAndSettle();
    expect(find.text('Ironhold City'), findsWidgets);
    // Every region on the map carries its level.
    expect(find.text('Lvl 40'), findsOneWidget);
    expect(find.text('Lvl 90'), findsOneWidget);

    Future<void> openRegion(String code) async {
      await tester.tap(find.byKey(Key('map-$code')));
      await tester.pumpAndSettle();
    }

    Future<void> closeSheet() async {
      await tester.tapAt(const Offset(10, 10));
      await tester.pumpAndSettle();
    }

    await openRegion('dragon_peaks');
    expect(tester.widget<FilledButton>(find.byKey(const Key('map-go'))).onPressed, isNull, reason: 'needs level 40');
    await closeSheet();
    // Far regions are one journey away too (DD-280): the shortest road's time.
    await openRegion('void_coast');
    expect(find.text('Travel time: 25m'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.byKey(const Key('map-go'))).onPressed, isNull, reason: 'needs level 90');
    await closeSheet();

    await openRegion('whisperwood');
    expect(find.text('Travel time: 5m'), findsOneWidget);
    await tester.tap(find.byKey(const Key('map-go')));
    await tester.pumpAndSettle();
    expect(be.travel!['to_region_code'], 'whisperwood');
    expect(find.text('Travelling to Whisperwood'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Arrive')).onPressed, isNull);

    await tester.tap(find.text('Turn back'));
    await tester.pumpAndSettle();
    expect(be.travel, isNull);
    expect(find.textContaining('any region your level allows'), findsOneWidget);
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
    expect(find.widgetWithText(FilledButton, 'Learn'), findsNWidgets(3));

    await tester.tap(find.widgetWithText(FilledButton, 'Learn').first);
    await tester.pumpAndSettle();
    expect(be.skills, {'alchemy': 1});
    expect(find.text('Linseed 3/2'), findsNothing, reason: 'recipes live on the profession page');
    expect(find.text('2 recipes'), findsOneWidget);

    // Tapping the profession opens its recipes; tapping a name unfolds it.
    await tester.tap(find.text('Alchemy'));
    await tester.pumpAndSettle();
    expect(find.text('Restores 60 HP'), findsNothing);
    await tester.tap(find.text('Potion Hp Small'));
    await tester.pumpAndSettle();
    expect(find.text('Restores 60 HP'), findsOneWidget);
    expect(find.text('Requires level'), findsOneWidget);
    expect(find.text('Stack'), findsOneWidget);
    expect(find.text('Linseed 3/2'), findsOneWidget);
    expect(find.text('Needs skill 3'), findsOneWidget);

    final crafts = find.widgetWithText(FilledButton, 'Craft');
    expect(tester.widget<FilledButton>(crafts.at(1)).onPressed, isNull, reason: 'skill 1 < min 3');
    expect(find.text('30s each'), findsOneWidget);
    // Recipe colors (DD-294): at skill 1 the skill-1 potion is orange, full XP.
    expect(find.byKey(const Key('difficulty-legend')), findsOneWidget);
    expect(tester.widget<Text>(find.byKey(const Key('diff-recipe_potion_hp_small'))).data, '100% XP');
    expect(find.text('1m each'), findsOneWidget, reason: 'no time from the server: 60 s');
    await tester.tap(crafts.first);
    await tester.pumpAndSettle();
    // A timed job (DD-264): materials taken, nothing made yet, one job at a time.
    expect(find.text('Crafting Potion Hp Small ×1'), findsOneWidget);
    expect(find.text('0/1 ready'), findsOneWidget);
    expect(find.text('Linseed 1/2'), findsOneWidget, reason: 'materials reload after starting');
    expect(tester.widget<FilledButton>(find.widgetWithText(FilledButton, 'Craft').first).onPressed, isNull);

    // The job is over server-side: collect without a question.
    be.craftJob!['started_at'] = DateTime.now().toUtc().subtract(const Duration(minutes: 1)).toIso8601String();
    be.craftJob!['ready_at'] = DateTime.now().toUtc().toIso8601String();
    await tester.fling(find.byType(ListView).last, const Offset(0, 400), 1000);
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Collect'));
    await tester.pumpAndSettle();
    expect(find.text('Crafted Potion Hp Small ×1 · Alchemy is now level 2'), findsOneWidget);
    expect(find.text('Crafting Potion Hp Small ×1'), findsNothing);
  });

  testWidgets('workshop: recipes sit under headings that fold; scrolls say what they add', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.skills.addAll({'blacksmithing': 1, 'enchanting': 1});
    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();

    // Blacksmithing: helms and gloves apart; a tapped heading folds.
    await tester.tap(find.text('Blacksmithing'));
    await tester.pumpAndSettle();
    expect(find.byKey(const Key('group-slot:helm')), findsOneWidget);
    expect(find.byKey(const Key('group-slot:gloves')), findsOneWidget);
    expect(find.text('Gloves Iron'), findsOneWidget);
    await tester.tap(find.byKey(const Key('group-slot:gloves')));
    await tester.pumpAndSettle();
    expect(find.text('Gloves Iron'), findsNothing);
    expect(find.text('Helm Iron'), findsOneWidget);
    await tester.pageBack();
    await tester.pumpAndSettle();

    // Enchanting: weapon and profession enchants apart, each with its effect.
    await tester.tap(find.text('Enchanting'));
    await tester.pumpAndSettle();
    expect(find.text('Weapon enchants'), findsOneWidget);
    expect(find.text('Profession enchants'), findsOneWidget);
    expect(find.text('+1% Attack'), findsOneWidget);
    expect(find.text('+2% Mining yield'), findsOneWidget);
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
    // The game's own line: a mute for swearing (DD-312).
    ws.last.push('chat.message', {...push(97, 'Brak', 'char-9', ''), 'system': 'profanity_mute'});
    await tester.pumpAndSettle();
    expect(find.text('Brak was muted for 5 minutes for swearing.'), findsOneWidget);

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

  testWidgets('account menu: data export link, account deletion signs out with the date', (tester) async {
    final be = FakeBackend();
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await _boot(tester, be, store: store);

    await tester.tap(find.byTooltip('Account'));
    await tester.pumpAndSettle();
    expect(find.text('Privacy Policy'), findsOneWidget);
    final opened = <Uri>[];
    launchExternal = (uri) async {
      opened.add(uri);
      return true;
    };
    addTearDown(() => launchExternal = (_) async => false);
    await tester.tap(find.text('Download my data'));
    await tester.pumpAndSettle();
    expect(be.exportLinks, 1);
    expect(opened.single.path, '/v1/me/export');
    expect(find.text('The download link opened in your browser (valid for 15 minutes).'), findsOneWidget);

    await tester.tap(find.byTooltip('Account'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Delete my account'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const Key('delete-password')), 'wrong-horse');
    await tester.tap(find.text('Delete account'));
    await tester.pumpAndSettle();
    expect(find.text('Wrong email/username or password.'), findsOneWidget);
    expect(be.deletionRequested, isFalse);
    await tester.enterText(find.byKey(const Key('delete-password')), 'correct-horse');
    await tester.tap(find.text('Delete account'));
    await tester.pumpAndSettle();
    expect(be.deletionRequested, isTrue);
    expect(find.textContaining('Your account will be deleted on'), findsOneWidget);
    expect(find.text('New here? Create an account'), findsOneWidget, reason: 'back on the sign-in screen');
  });

  testWidgets('support: answered ticket shows the reply, new ticket validates and sends client info', (tester) async {
    final be = FakeBackend();
    be.tickets.add({
      'id': 'tkt-0',
      'category': 'account',
      'subject': 'Lost my sword',
      'body': 'It vanished after a hunt.',
      'state': 'answered',
      'reply': 'Restored it, sorry!',
      'created_at': '2026-10-01T10:00:00Z',
    });
    final store = MemoryTokenStore();
    await be.api(store).login('grim', 'correct-horse');
    await _boot(tester, be, store: store);

    await tester.tap(find.byTooltip('Account'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Support'));
    await tester.pumpAndSettle();
    expect(find.text('Lost my sword'), findsOneWidget);
    expect(find.text('Answered'), findsOneWidget);
    expect(find.text('Restored it, sorry!'), findsOneWidget);

    await tester.tap(find.text('New request'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const Key('support-subject')), 'Bug');
    await tester.enterText(find.byKey(const Key('support-body')), 'short');
    await tester.tap(find.text('Send'));
    await tester.pumpAndSettle();
    expect(find.text('Subject needs at least 3 characters, the message at least 10.'), findsOneWidget);
    expect(be.lastTicket, isNull);

    await tester.enterText(find.byKey(const Key('support-body')), 'The bag shows an empty slot after selling.');
    await tester.tap(find.text('Send'));
    await tester.pumpAndSettle();
    expect(be.lastTicket!['category'], 'bug');
    expect((be.lastTicket!['client_info'] as Map)['locale'], 'en');
    expect(find.text("Request received. We'll let you know when there's an answer."), findsOneWidget);
    expect(find.text('Open'), findsOneWidget, reason: 'the new ticket is listed');
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

    await tester.tap(find.byTooltip('Account'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Change password'));
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

  testWidgets('party: on the character sheet: accept, activities, form by invite, promote, chat (DD-301)', (
    tester,
  ) async {
    final (be, ws, _) = await _inGame(tester);
    // The party lives on the character sheet; an invite push reloads it.
    be.partyInvites.add({'id': 'pi-1', 'invited_by_name': 'Aria'});
    ws.last.push('party.invite');
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(find.text('Invitations'), 200, scrollable: _list);
    await tester.ensureVisible(find.byTooltip('Accept'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Accept'));
    await tester.pumpAndSettle();
    expect(find.text('Party (2/5)'), findsOneWidget);
    expect(find.byTooltip('More'), findsNothing, reason: 'members cannot manage the party');

    // What each member is doing (DD-301).
    (be.party!['members'] as List).first['activity'] = {'kind': 'gathering', 'detail': 'mining'};
    ws.last.push('party.changed');
    await tester.pumpAndSettle();
    expect(find.textContaining('Gathering: Mining'), findsOneWidget);
    expect(find.textContaining('Idle'), findsOneWidget);

    await tester.ensureVisible(find.widgetWithText(OutlinedButton, 'Leave'));
    await tester.tap(find.widgetWithText(OutlinedButton, 'Leave'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Leave'));
    await tester.pumpAndSettle();
    // No party: the first invite forms one (DD-297).
    expect(find.text('Create party'), findsNothing);
    await tester.ensureVisible(find.widgetWithText(FilledButton, 'Invite'));
    await tester.enterText(find.byType(TextField), 'Brak');
    await tester.tap(find.widgetWithText(FilledButton, 'Invite'));
    await tester.pumpAndSettle();
    expect(be.groupCalls, contains('party/invite {"target_name":"Brak"}'));
    expect(find.text('Party (1/5)'), findsOneWidget);

    // Someone joins: the server pushes party.changed and the panel reloads.
    (be.party!['members'] as List).add({'character_id': 'char-7', 'name': 'Brak', 'level': 2, 'class_code': 'mage'});
    ws.last.push('party.changed');
    await tester.pumpAndSettle();
    expect(find.text('Party (2/5)'), findsOneWidget);

    await tester.ensureVisible(find.byTooltip('More'));
    await tester.tap(find.byTooltip('More'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Make leader'));
    await tester.pumpAndSettle();
    expect(be.party!['leader_character_id'], 'char-7');
    expect(find.byTooltip('More'), findsNothing, reason: 'no longer leader');

    // Party chat opens Social on the party channel.
    await tester.ensureVisible(find.text('Party chat'));
    await tester.tap(find.text('Party chat'));
    await tester.pumpAndSettle();
    expect(be.chatRead.last, 'party');
    expect(find.widgetWithText(Tab, 'Party'), findsNothing, reason: 'no party tab in Social any more');
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

  testWidgets('rewards tab lists collected loot as stacks, apart from the inbox', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.rewardLog.add({
      'id': 'r-1',
      'kind': 'reward_collected',
      'read_at': DateTime.now().toUtc().toIso8601String(),
      'created_at': DateTime.now().toUtc().subtract(const Duration(minutes: 2)).toIso8601String(),
      'body': {
        'kind': 'reward_collected',
        'body_vars': {
          'source': 'hunt',
          'target': 'young_boar',
          'gold_copper': 120,
          'xp': 45,
          'new_level': 6,
          'items': [
            {'code': 'boar_hide', 'rarity': 'common', 'count': 5},
            {'code': 'sword_iron', 'rarity': 'rare', 'prefix': 'ancient', 'count': 1},
          ],
        },
      },
    });
    await tester.tap(find.byTooltip('Notifications'));
    await tester.pumpAndSettle();
    expect(find.text('Hunt: Young Boar'), findsNothing, reason: 'rewards are not in the regular inbox');

    await tester.tap(find.text('Rewards'));
    await tester.pumpAndSettle();
    expect(find.text('Hunt: Young Boar'), findsOneWidget);
    expect(find.text('Boar Hide × 5'), findsOneWidget);
    expect(find.text('+45 XP'), findsOneWidget);
    expect(find.text('Level 6!'), findsOneWidget);
    expect(find.textContaining('Sword Iron'), findsOneWidget);
    expect(find.textContaining('× 1'), findsNothing, reason: 'a single drop shows no count');

    await tester.tap(find.byTooltip('Delete'));
    await tester.pumpAndSettle();
    expect(be.rewardLog, isEmpty);
    expect(find.text('Hunt: Young Boar'), findsNothing);
  });

  testWidgets('workshop: choose how many to craft', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.skills['alchemy'] = 1;
    be.bag.add({
      'instance_id': 'seed',
      'template_code': 'linseed',
      'item_class': 'material',
      'item_subclass': 'herb',
      'slot': 0,
      'stack_quantity': 7,
      'max_stack': 99,
      'rarity': 'common',
    });
    await tester.tap(find.byIcon(Icons.handyman_outlined).first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Alchemy'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Craft').first);
    await tester.pumpAndSettle();
    expect(find.text('How many Potion Hp Small?'), findsOneWidget);
    await tester.enterText(find.byKey(const Key('craft-count')), '2');
    await tester.pump();
    expect(find.text('Craft ×2'), findsOneWidget);
    await tester.tap(find.text('max 3'));
    await tester.pump();
    expect(find.text('Total 1m'), findsOneWidget);
    await tester.tap(find.text('Craft ×3'));
    await tester.pumpAndSettle();
    expect(be.lastCraftCount, 3);
    expect(find.text('Crafting Potion Hp Small ×3'), findsOneWidget);
    expect(find.text('Linseed 1/2'), findsOneWidget);

    // Collecting early asks first and refunds the unfinished units.
    await tester.tap(find.widgetWithText(FilledButton, 'Collect'));
    await tester.pumpAndSettle();
    expect(
      find.text('The 0 finished items are delivered; the materials and gold of the other 3 come back.'),
      findsOneWidget,
    );
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last);
    await tester.pumpAndSettle();
    expect(find.text('0 made, 3 refunded'), findsOneWidget);
    expect(find.text('Linseed 7/2'), findsOneWidget);
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
    expect(find.byTooltip('Notifications'), findsOneWidget);
    expect(find.descendant(of: find.byType(Badge), matching: find.text('2')), findsNothing);
  });

  testWidgets('progression: achievements sorted, mastery bonus, reputation, title select + clear', (tester) async {
    final (be, _, _) = await _inGame(tester);
    await tester.scrollUntilVisible(find.text('Achievements & mastery'), 200, scrollable: _list);
    await tester.ensureVisible(find.text('Achievements & mastery'));
    await tester.pumpAndSettle();
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

  testWidgets('one action at a time: a busy character is told what to finish first', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.busyWith = 'crafting';
    await tester.tap(find.byIcon(Icons.explore_outlined));
    await tester.pumpAndSettle();
    final panel = find.byWidget(
      tester.widget<Scrollable>(
        find.ancestor(of: find.text('Choose your prey'), matching: find.byType(Scrollable)).first,
      ),
    );
    await tester.scrollUntilVisible(find.text('Start hunt'), 200, scrollable: panel);
    await tester.ensureVisible(find.text('Start hunt'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Start hunt'));
    await tester.pumpAndSettle();
    expect(find.text('Finish what your character is doing first: crafting.'), findsOneWidget);
    expect(be.activeHunt, isNull);
  });

  testWidgets('events: tapping one explains how to take part and what it pays', (tester) async {
    final (be, _, _) = await _inGame(tester);
    final now = DateTime.now().toUtc();
    be.worldEvents.add({
      'id': 's-1',
      'template_code': 'greenwood_incursion',
      'kind': 'invasion',
      'min_level': 1,
      'region_code': 'greenwood',
      'goal_score': 500000,
      'total_score': 0,
      'starts_at': now.subtract(const Duration(minutes: 10)).toIso8601String(),
      'ends_at': now.add(const Duration(minutes: 50)).toIso8601String(),
      'state': 'active',
      'tiers': [
        {'min_score': 1000, 'gold_copper': 500, 'event_tokens': 5},
        {'min_score': 3000, 'gold_copper': 2000, 'event_tokens': 10},
        {'min_score': 6000, 'gold_copper': 5000, 'event_tokens': 20},
      ],
      'afk_cap_tier': 2,
      'power_now': 100,
      'enlisted': false,
    });
    await tester.tap(find.text('Adventure'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Events'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Greenwood Incursion'));
    await tester.pumpAndSettle();
    expect(find.text('How to take part'), findsOneWidget);
    expect(find.text('Be in Greenwood when you enlist.'), findsOneWidget);
    expect(find.textContaining('Yours now: 100.'), findsOneWidget);
    expect(find.textContaining('It stops at tier 2.'), findsOneWidget);
    final sheet = find.byType(Scrollable).last; // the sheet sits on top
    await tester.scrollUntilVisible(find.text('Join now: about 3000 points → tier 2'), 200, scrollable: sheet);
    expect(find.text('Tier 3 · 6000+ points'), findsOneWidget);
    expect(find.text('+20 event tokens'), findsOneWidget);
    // 100 power × ~50 minutes = 5000, capped at tier 2's 3000.
    expect(find.text('Join now: about 3000 points → tier 2'), findsOneWidget);
    Navigator.of(tester.element(find.text('Join now: about 3000 points → tier 2'))).pop();
    await tester.pumpAndSettle();

    await tester.tap(find.text('Harvest Festival 2026'));
    await tester.pumpAndSettle();
    expect(find.textContaining('claim 15 event tokens once a day'), findsOneWidget);
    expect(find.text('Tokens stay on your account (balance: 0).'), findsOneWidget);
  });

  testWidgets('quests: accept, live completion badge, turn in with reward + level up, abandon daily', (tester) async {
    final (be, ws, _) = await _inGame(tester);
    await tester.tap(find.byTooltip('Quests'));
    await tester.pumpAndSettle();
    expect(find.text('No open quests. Pick one below.'), findsOneWidget);
    expect(find.text('First Hunt'), findsOneWidget);
    expect(find.text('• Defeat Young Boar ×5'), findsOneWidget);
    expect(find.textContaining('Resets in'), findsOneWidget);

    await tester.tap(find.widgetWithText(FilledButton, 'Accept').first);
    await tester.pumpAndSettle();
    expect(be.questCalls, ['accept main_01_first_hunt']);
    expect(find.text('0/5'), findsOneWidget);
    expect(find.widgetWithText(FilledButton, 'Turn in'), findsNothing, reason: 'not complete yet');

    await tester.pageBack();
    await tester.pumpAndSettle();
    be.completeQuest('main_01_first_hunt');
    ws.last.push('quest.completed', {'code': 'main_01_first_hunt'});
    await tester.pumpAndSettle();
    expect(find.text('Quest objectives complete — turn it in!'), findsOneWidget);
    expect(find.descendant(of: find.byType(Badge), matching: find.text('1')), findsOneWidget);

    await tester.tap(find.byTooltip('Quests'));
    await tester.pumpAndSettle();
    final sheetsBefore = be.sheetGets;
    await tester.tap(find.widgetWithText(FilledButton, 'Turn in'));
    await tester.pumpAndSettle();
    expect(find.text('Quest complete: +80 XP · +2s 0c\nLevel up! Now level 4'), findsOneWidget);
    expect(be.sheetGets, greaterThan(sheetsBefore), reason: 'XP and gold refresh the sheet');

    await tester.tap(find.widgetWithText(FilledButton, 'Accept'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Abandon'));
    await tester.pumpAndSettle();
    expect(find.text("Progress is lost and this daily can't be taken again today."), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Abandon'));
    await tester.pumpAndSettle();
    expect(be.questCalls.last, 'abandon daily_miner');
    expect(find.text('No open quests. Pick one below.'), findsOneWidget);
  });

  testWidgets('class skills: unlocked active + passive described, higher tier locked', (tester) async {
    await _inGame(tester);
    // Folded until opened (DD-301).
    await tester.scrollUntilVisible(find.byKey(const Key('skills-section')), 200, scrollable: _list);
    expect(find.text('Backstab · Active'), findsNothing);
    final header = find.descendant(of: find.byKey(const Key('skills-section')), matching: find.byType(InkWell)).first;
    await tester.ensureVisible(header);
    await tester.pumpAndSettle();
    await tester.tap(header);
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(find.text('Eviscerate · Active'), 200, scrollable: _list);
    expect(find.text('Backstab · Active'), findsOneWidget);
    expect(find.text('200% damage · 12 MP · every 3 turns'), findsOneWidget);
    expect(find.text('Precision · Passive'), findsOneWidget);
    expect(find.text('+5 Crit'), findsOneWidget);
    expect(find.text('Lvl 6+'), findsOneWidget, reason: 'level 3 rogue: eviscerate still locked');
  });

  testWidgets('leveling: promotion at the gate, attributes, talents, MAX at the cap', (tester) async {
    final (be, _, _) = await _inGame(tester);
    be.leveling = {
      'level': 100,
      'level_cap': 1000,
      'xp': 0,
      'xp_to_next': 6719,
      'maxed': false,
      'class_code': 'warrior',
      'spec_code': 'warrior',
      'path': ['warrior'],
      'promotion': {
        'available': true,
        'overflow_xp': 25000,
        'choices': [
          {
            'code': 'guardian',
            'role_tags': ['tank'],
          },
          {
            'code': 'executioner',
            'role_tags': ['physical_dps'],
          },
        ],
      },
      'attributes': {
        'available': 30,
        'raw': {'strength': 10, 'vitality': 10, 'dexterity': 7, 'intelligence': 4, 'wisdom': 5, 'luck': 5},
        'gains': {
          'strength': {'attack': 2.2, 'block': 0.03},
          'vitality': {'hp': 18, 'defense': 1.2, 'hp_regen': 0.05},
          'intelligence': <String, dynamic>{},
        },
      },
      'talents': {'available': 2},
      'title': 'veteran',
      'mastery': {'level': 1000},
      'next_milestone': {'level': 100, 'kind': 'promotion'},
    };
    be.sheetGets = 0;
    await tester.drag(_list, const Offset(0, 300)); // pull to refresh: the card loads
    await tester.pumpAndSettle();
    expect(find.text('PROMOTION AVAILABLE'), findsOneWidget);
    expect(find.text('25000 XP is waiting for your choice'), findsOneWidget);

    await tester.tap(find.text('PROMOTION AVAILABLE'));
    await tester.pumpAndSettle();
    expect(find.text('Guardian'), findsOneWidget);
    expect(find.text('Executioner'), findsOneWidget);
    await tester.tap(find.text('Guardian'));
    await tester.pumpAndSettle();
    expect(find.text('Become Guardian?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Choose'));
    await tester.pumpAndSettle();
    expect(be.promotedTo, 'guardian');
    expect(find.text('PROMOTION AVAILABLE'), findsNothing, reason: 'the card reloads after choosing');

    await tester.tap(find.text('Attributes · 30 points'));
    await tester.pumpAndSettle();
    // What a point adds to this class (DD-285).
    expect(find.text('Per point: +2.2 Attack · +0.03% Block'), findsOneWidget);
    expect(find.text('Per point: +18 HP · +1.2 Defense · +0.05 HP/s'), findsOneWidget);
    expect(find.text('No effect for your class'), findsWidgets);
    // Points are typed into each attribute's box (DD-292).
    await tester.enterText(find.byKey(const Key('attr-strength')), '2');
    await tester.enterText(find.byKey(const Key('attr-vitality')), '1');
    await tester.pump();
    expect(find.text('Adding: +4.4 Attack · +0.06% Block'), findsOneWidget);
    expect(find.text('Attributes · 27 points'), findsOneWidget);

    // More than is left snaps to what is left; an emptied box counts as 0.
    await tester.enterText(find.byKey(const Key('attr-vitality')), '99');
    await tester.pump();
    expect(find.text('Attributes · 0 points'), findsOneWidget);
    expect(
      tester
          .widget<TextField>(
            find.descendant(of: find.byKey(const Key('attr-vitality')), matching: find.byType(TextField)),
          )
          .controller!
          .text,
      '28',
    );
    expect(find.text('Adding: +504 HP · +33.6 Defense · +1.4 HP/s'), findsOneWidget, reason: '28 points');
    await tester.enterText(find.byKey(const Key('attr-vitality')), '');
    await tester.pump();
    expect(find.text('Attributes · 28 points'), findsOneWidget);
    await tester.enterText(find.byKey(const Key('attr-vitality')), '17');
    await tester.pump();

    await tester.tap(find.widgetWithText(FilledButton, 'Save'));
    await tester.pumpAndSettle();
    expect(be.lastAllocation, {'strength': 2, 'vitality': 17});

    await tester.tap(find.text('Talents · 2 points'));
    await tester.pumpAndSettle();
    expect(find.text('Needs level 850 and 25 points in this tree'), findsOneWidget, reason: 'capstone locked');
    await tester.tap(find.byIcon(Icons.add_circle_outline).first);
    await tester.pumpAndSettle();
    expect(be.learnedTalents, ['warrior_defense_toughness']);
    await tester.pageBack();
    await tester.pumpAndSettle();

    be.leveling!['level'] = 1000;
    be.leveling!['maxed'] = true;
    be.leveling!['xp_to_next'] = 0;
    be.leveling!['mastery'] = {'level': 1000, 'unlocked': true, 'title': 'the_unbreakable'};
    be.leveling!.remove('next_milestone');
    await tester.drag(_list, const Offset(0, 300));
    await tester.pumpAndSettle();
    expect(find.text('MAX'), findsOneWidget);
    expect(find.text('Class mastery: The Unbreakable'), findsOneWidget);
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

  test('hunting XP falls off with the level gap either way (DD-281)', () {
    const f = HuntXPFalloff(freePct: 0.05, power: 4, floor: 0.1);
    expect(f.xp(200, 40, 42), 200, reason: '5% of 40: two levels free');
    expect(f.xp(200, 40, 38), 200);
    expect(f.xp(200, 40, 46), 131, reason: 'four past free: (1 - 0.1)^4');
    expect(f.xp(200, 40, 34), 131);
    expect(f.xp(200, 10, 3), 20, reason: 'level 10 farming level 3: the 10% floor');
    expect(f.xp(200, 400, 1), 20, reason: 'the floor');
    expect(const HuntXPFalloff().xp(200, 400, 1), 200, reason: 'older server: no falloff');
    expect(
      HuntXPFalloff.fromJson({'free_pct': 0.05, 'power': 4, 'floor': 0.1}).factor(500, 460),
      closeTo(0.88529281, 1e-9),
    );
  });

  test('splitCopper', () {
    expect(splitCopper(123456), (gold: 12, silver: 34, copper: 56));
    expect(splitCopper(-5), (gold: 0, silver: 0, copper: 0));
  });
}

/// The Gather picker outgrows the test screen once resources are listed.
Future<void> _gatherScrollTo(WidgetTester tester, String text) async {
  await tester.scrollUntilVisible(
    find.text(text),
    200,
    scrollable: find.descendant(of: find.byType(ActivityView), matching: find.byType(Scrollable)).first,
  );
  await tester.ensureVisible(find.text(text)); // fully, not just its top edge
  await tester.pumpAndSettle();
}
