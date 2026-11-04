import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/character/class_catalog.dart';
import 'package:otter/features/character/class_info.dart';
import 'package:otter/features/character/models.dart';

void main() {
  test('class pyramids narrow upward', () {
    expect(pyramidRows(3), [3]);
    expect(pyramidRows(4), [4]);
    expect(pyramidRows(5), [1, 4]);
    expect(pyramidRows(6), [2, 4]);
    expect(pyramidRows(8), [1, 3, 4]);
    expect(pyramidRows(10), [1, 2, 3, 4]);
    expect(pyramidRows(0), isEmpty);
  });

  test('race passives read in percent under the translated attribute names', () {
    final cat = ClassCatalog.fromJson({
      'races': {
        'undead': {
          'classes': ['rogue'],
          'passives': [
            {'kind': 'mult_add', 'target': 'agi', 'amount': 0.01},
            {'kind': 'mult_add', 'target': 'cha', 'amount': 0.02},
          ],
          'profession_xp': {'alchemy': 0.02},
        },
      },
    });
    final r = cat.races['undead']!;
    expect(r.attributes, {'dexterity': 1, 'wisdom': 2});
    expect(r.professionXp, {'alchemy': 2});
    expect(cat.classesOf('undead'), ['rogue']);
    expect(cat.classesOf('elf'), classes, reason: 'an unlisted race: the server decides');
  });
}
