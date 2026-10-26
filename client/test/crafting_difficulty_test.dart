import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';
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
