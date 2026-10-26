import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';

EquippedItem worn(String code) =>
    EquippedItem(instanceId: code, templateCode: code, rarity: 'common', itemClass: 'armor');

void main() {
  test('compareBonuses sums each stat and keeps flat and percent apart', () {
    final d = compareBonuses(
      [
        const StatBonus('flat', 'attack', 12),
        const StatBonus('flat', 'attack', 3), // an affix on the same stat
        const StatBonus('mult_add', 'attack', 0.05),
        const StatBonus('flat', 'crit_pct', 2),
        const StatBonus('clamp', 'hit_pct', 95), // not a bonus a player compares
      ],
      [const StatBonus('flat', 'attack', 10), const StatBonus('flat', 'defense', 8)],
    );
    expect(
      [for (final x in d) '${x.kind} ${x.target} ${x.worn} ${x.next} ${x.change}'],
      [
        'flat attack 10.0 15.0 5.0',
        'mult_add attack 0.0 0.05 0.05',
        'flat crit_pct 0.0 2.0 2.0',
        'flat defense 8.0 0.0 -8.0', // only the worn piece has it
      ],
    );
  });

  test('equipTarget: a ring or trinket replaces what the server would swap out', () {
    expect(equipTarget('gloves', {}), 'gloves');
    expect(equipTarget('ring', {}), 'ring1');
    expect(equipTarget('ring', {'ring1': worn('a')}), 'ring2', reason: 'the free one of the pair');
    expect(equipTarget('ring', {'ring2': worn('b')}), 'ring1');
    expect(equipTarget('ring', {'ring1': worn('a'), 'ring2': worn('b')}), 'ring1', reason: 'both taken: the first');
    expect(equipTarget('trinket', {'trinket1': worn('a')}), 'trinket2');
  });
}
