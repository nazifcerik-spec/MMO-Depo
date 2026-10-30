import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/helpers/helpers_screen.dart';
import 'package:otter/features/workshop/workshop_repository.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

Future<FakeBackend> _open(WidgetTester tester, [void Function(FakeBackend)? setup]) async {
  final be = FakeBackend();
  setup?.call(be);
  final api = be.api();
  await tester.runAsync(() => api.login('grim', 'correct-horse'));
  await tester.pumpWidget(
    MaterialApp(
      localizationsDelegates: AppLocalizations.localizationsDelegates,
      supportedLocales: AppLocalizations.supportedLocales,
      home: HelpersScreen(
        repo: HelperRepository(api, 'char-1'),
        workshop: WorkshopRepository(api, 'char-1'),
        dataVersion: ValueNotifier(0),
      ),
    ),
  );
  await tester.pumpAndSettle();
  return be;
}

void main() {
  testWidgets('helpers: hire, send to gather a chosen resource, collect (DD-319, DD-327)', (tester) async {
    final be = await _open(tester);
    expect(find.text('Miner · Mining'), findsOneWidget);
    expect(find.text('+15% yield'), findsOneWidget);

    await tester.tap(find.textContaining('Hire ·'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Hire'));
    await tester.pumpAndSettle();
    expect(be.helpers.first['hired'], isTrue);

    await tester.tap(find.text('Send to work'));
    await tester.pumpAndSettle();
    // The resources the level unlocks, with the best preselected.
    expect(find.byKey(const Key('helper-node-raw_stone')), findsOneWidget);
    await tester.tap(find.byKey(const Key('helper-node-raw_stone')));
    await tester.tap(find.text('3h 00m'));
    await tester.pump();
    await tester.tap(find.byKey(const Key('helper-gather-send')));
    await tester.pumpAndSettle();
    final spec = (be.helpers.first['job'] as Map)['spec'] as Map;
    expect(spec['target'], 'raw_stone');
    expect(spec['duration_seconds'], 3 * 3600);
    expect(find.textContaining('Working,'), findsOneWidget);

    await tester.tap(find.widgetWithText(FilledButton, 'Collect'));
    await tester.pumpAndSettle();
    expect(find.text('1 items, +12 Mining XP'), findsOneWidget);
    expect(find.text('Send to work'), findsOneWidget);
  });

  testWidgets('crafting helper: how many it can make, max picks that, item details (DD-327)', (tester) async {
    final be = await _open(tester, (be) {
      be.helpers.add({
        'code': 'apothecary',
        'kind': 'craft',
        'profession': 'alchemy',
        'bonus': 'double',
        'bonus_pct': 10,
        'price_copper': 100000,
        'hired': true,
      });
      be.skills['alchemy'] = 1;
      be.characters.add({'id': 'char-1', 'name': 'Grim', 'class_code': 'warrior', 'level': 7, 'state': 'active'});
      be.bag.add({
        'instance_id': 'seed',
        'template_code': 'linseed',
        'item_class': 'material',
        'item_subclass': 'herb',
        'slot': 0,
        'stack_quantity': 9,
        'max_stack': 99,
        'rarity': 'common',
      });
    });
    await tester.tap(
      find.descendant(of: find.byKey(const Key('helper-apothecary')), matching: find.text('Send to work')),
    );
    await tester.pumpAndSettle();

    // 9 linseed at 2 a potion: 4; the mana potion's skill is not reached.
    expect(find.text('Can make: 4'), findsOneWidget);
    expect(find.byKey(const Key('can-make-recipe_potion_mp_small')), findsNothing);
    expect(find.text('1m each'), findsOneWidget, reason: 'half speed: twice the 30 s');

    await tester.tap(find.text('Potion Hp Small'));
    await tester.pumpAndSettle();
    expect(find.text('Restores 60 HP'), findsOneWidget);

    await tester.tap(find.widgetWithText(FilledButton, 'Send to work'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('max 4'));
    await tester.pump();
    await tester.tap(find.widgetWithText(FilledButton, 'Craft ×4'));
    await tester.pumpAndSettle();
    final spec = (be.helpers.last['job'] as Map)['spec'] as Map;
    expect(spec['recipe'], 'recipe_potion_hp_small');
    expect(spec['count'], 4);
  });
}
