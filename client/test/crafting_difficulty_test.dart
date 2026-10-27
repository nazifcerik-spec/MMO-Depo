import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';
import 'package:otter/features/workshop/profession_screen.dart';
import 'package:otter/features/workshop/workshop_repository.dart';

void main() {
  test('recipeDifficulty follows the server bands (DD-294)', () {
    String color(int min, int skill) => recipeDifficulty(defaultDifficulty, min, skill).$1;
    expect(color(10, 10), 'orange');
    expect(color(10, 14), 'orange');
    expect(color(10, 15), 'yellow');
    expect(color(10, 20), 'green');
    expect(color(10, 29), 'green');
    expect(color(10, 30), 'grey');
    expect(recipeDifficulty(defaultDifficulty, 10, 15).$2, 0.6);
    expect(recipeDifficulty(defaultDifficulty, 10, 30).$2, 0);
    final custom = [const DifficultyBand(2, 'orange', 1)];
    expect(recipeDifficulty(custom, 1, 3).$1, 'grey');
  });

  test('groupRecipes puts recipes under their headings in order (DD-295)', () {
    Recipe r(String code, ItemInfo? info) => Recipe(
      code: code,
      profession: 'x',
      minSkill: 1,
      inputs: const {},
      goldCost: 0,
      output: code,
      outputStack: 1,
      xpReward: 0,
      info: info,
    );
    ItemInfo gear(String cls, String sub, String slot) => ItemInfo(itemClass: cls, itemSubclass: sub, equipSlot: slot);
    ItemInfo scroll(String slot, String target) => ItemInfo(
      itemClass: 'consumable',
      itemSubclass: 'enchant_scroll',
      enchant: EnchantInfo(slots: [slot], modifiers: [StatBonus('mult_add', target, 0.01)]),
    );
    final groups = groupRecipes([
      r('gloves1', gear('armor', 'plate_armor', 'gloves')),
      r('ring', gear('jewelry', 'ring', 'ring')),
      r('mining', scroll('jewelry', 'gather_yield_mining')),
      r('crit', scroll('jewelry', 'crit_pct')),
      r('shield', gear('armor', 'shield', 'offhand')),
      r('helm', gear('armor', 'plate_armor', 'helm')),
      r('sword', gear('weapon', 'sword', 'mainhand')),
      r('axe', gear('weapon', 'axe', 'mainhand')),
      r('gloves2', gear('armor', 'plate_armor', 'gloves')),
      r('bar', const ItemInfo(itemClass: 'material', itemSubclass: 'bar')),
      r('attack', scroll('weapon', 'attack')),
    ]);
    expect(groups.map((g) => g.$1), [
      'equipment:axe',
      'equipment:sword',
      'equipment:shield',
      'slot:helm',
      'slot:gloves',
      'slot:ring',
      'enchant:weapon',
      'enchant:jewelry',
      'enchant:profession',
      'sub:bar',
    ]);
    expect(groups[4].$2.map((r) => r.code), ['gloves1', 'gloves2'], reason: 'order kept within a heading');
  });

  test('EnchantInfo.fits maps worn slots to their groups (DD-293)', () {
    const weapon = EnchantInfo(slots: ['weapon'], modifiers: []);
    expect(weapon.fits('mainhand'), isTrue);
    expect(weapon.fits('offhand'), isTrue);
    expect(weapon.fits('chest'), isFalse);
    const jewelry = EnchantInfo(slots: ['jewelry'], modifiers: []);
    expect(jewelry.fits('ring2'), isTrue, reason: 'a worn ring slot is a ring');
    expect(jewelry.fits('trinket1'), isTrue);
    expect(const EnchantInfo(slots: ['any'], modifiers: []).fits('boots'), isTrue);
    expect(EnchantInfo.fromJson(null), isNull);
    expect(
      EnchantInfo.fromJson({
        'slots': ['armor'],
        'modifiers': [
          {'kind': 'mult_add', 'target': 'defense', 'amount': 0.005},
        ],
      })!.modifiers.single.amount,
      0.005,
    );
  });
}
