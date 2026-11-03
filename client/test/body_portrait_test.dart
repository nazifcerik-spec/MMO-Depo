import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/character/body_portrait.dart';
import 'package:otter/features/inventory/inventory_models.dart';

EquippedItem chest(String subclass, int level) => EquippedItem(
  instanceId: 'c',
  templateCode: 'x',
  rarity: 'rare',
  itemClass: 'armor',
  itemSubclass: subclass,
  minLevel: level,
);

void main() {
  test('the chest piece picks the outfit: armor type and level band (DD-344)', () {
    expect(bodyOutfit(null), 'none');
    expect(bodyOutfit(chest('cloth_armor', 1)), 'cloth1');
    expect(bodyOutfit(chest('leather_armor', 99)), 'leather1');
    expect(bodyOutfit(chest('leather_armor', 100)), 'leather2');
    expect(bodyOutfit(chest('plate_armor', 640)), 'plate4');
    expect(bodyOutfit(chest('plate_armor', 1000)), 'plate6');
    expect(bodyOutfit(chest('jewelry', 300)), 'none', reason: 'not body armor');
  });

  test('every outfit of the drawn figures ships', () {
    for (final body in ['male', 'female']) {
      for (final o in [
        'none',
        for (final t in ['cloth', 'leather', 'plate'])
          for (var b = 1; b <= bodyBands.length; b++) '$t$b',
      ]) {
        expect(File('assets/bodies/human_${body}_$o.webp').existsSync(), isTrue, reason: 'human $body $o');
      }
    }
  });
}
