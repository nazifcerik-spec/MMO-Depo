import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';
import 'package:otter/features/workshop/profession_screen.dart';
import 'package:otter/features/workshop/workshop_repository.dart';

Recipe _r(String code, int skill) => Recipe(
  code: code,
  profession: 'blacksmithing',
  minSkill: skill,
  inputs: const {},
  goldCost: 0,
  output: code,
  outputStack: 1,
  xpReward: 1,
);

void main() {
  test('profession page shows what is learnable, newest first, then the next step (DD-274)', () {
    final all = [_r('a1', 1), _r('b1', 1), _r('c2', 2), _r('d5', 5), _r('e5', 5), _r('f9', 9)];
    expect(visibleRecipes(all, 2).map((r) => r.code), ['c2', 'a1', 'b1', 'd5', 'e5']);
    expect(visibleRecipes(all, 0).map((r) => r.code), ['a1', 'b1'], reason: 'unlearned: the first step only');
    expect(visibleRecipes(all, 9).map((r) => r.code).first, 'f9');
    expect(visibleRecipes(all, 9), hasLength(6));
  });

  test('heal_pct potions read as a share of max HP and are usable', () {
    final item = BagItem.fromJson({
      'instance_id': 'i-1',
      'template_code': 'potion_hp_large',
      'item_class': 'consumable',
      'use_effect': {'kind': 'heal_pct', 'amount': 40},
    });
    expect(item.healPct, 40);
    expect(item.healHp, 0);
    expect(item.usable, isTrue);
    final info = ItemInfo.fromJson({
      'item_class': 'consumable',
      'use_effect': {'kind': 'heal_pct', 'amount': 20},
    });
    expect(info.healPct, 20);
  });

  test('bag items show their effective bonuses: rarity and affixes applied (DD-275)', () {
    final item = BagItem.fromJson({
      'instance_id': 'i-2',
      'template_code': 'sword_iron',
      'item_class': 'weapon',
      'rarity': 'epic',
      'equip_slot': 'mainhand',
      'stat_bonuses': [
        {'kind': 'flat', 'target': 'attack', 'amount': 16},
      ],
      'effective_bonuses': [
        {'kind': 'flat', 'target': 'attack', 'amount': 19},
        {'kind': 'flat', 'target': 'attack', 'amount': 3},
      ],
    });
    expect(item.bonuses.map((b) => b.amount), [19, 3]);
  });
}
