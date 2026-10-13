import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/activity/activity_models.dart';
import 'package:otter/l10n/app_localizations.dart';
import 'package:otter/shared/loot_sheet.dart';

LootItem drop(String code, String rarity, [String? prefix]) =>
    LootItem(instanceId: '', templateCode: code, rarity: rarity, prefix: prefix);

void main() {
  test('identical drops stack: rarest first, then first-seen order, affixes apart', () {
    final stacks = stackLoot([
      drop('iron_ore', 'common'),
      drop('copper_ore', 'uncommon'),
      drop('iron_ore', 'common'),
      drop('iron_ore', 'rare'),
      drop('iron_ore', 'common'),
      drop('sword_iron', 'rare', 'ancient'),
      drop('copper_ore', 'uncommon'),
    ]);
    expect(
      [for (final s in stacks) '${s.item.templateCode}/${s.item.rarity}${s.item.prefix ?? ''}x${s.count}'],
      ['iron_ore/rarex1', 'sword_iron/rareancientx1', 'copper_ore/uncommonx2', 'iron_ore/commonx3'],
    );
    expect(stackLoot(const []), isEmpty);
  });

  testWidgets('the loot sheet shows one row per stack with its count', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('en'),
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: Scaffold(
          body: LootSheet(
            title: 'Mining',
            loot: Loot(items: [drop('iron_ore', 'common'), drop('iron_ore', 'common'), drop('gold_ore', 'rare')]),
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();
    expect(find.text('Iron Ore × 2'), findsOneWidget);
    expect(find.text('Gold Ore'), findsOneWidget);
  });
}
