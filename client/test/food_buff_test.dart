import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';

void main() {
  test('food buff from a use_effect and from the running buff (DD-269)', () {
    final b = FoodBuff.fromUseEffect({
      'kind': 'heal_hp',
      'amount': 60,
      'buff': {'attack_pct': 6, 'defense_pct': 4, 'minutes': 60},
    })!;
    expect((b.attackPct, b.defensePct, b.minutes), (6, 4, 60));
    expect(FoodBuff.fromUseEffect({'kind': 'heal_hp', 'amount': 60}), isNull, reason: 'a potion has no buff');

    final live = FoodBuff.fromActive({'attack_pct': 9, 'defense_pct': 9, 'expires_at': '2026-10-19T12:00:00Z'})!;
    expect(live.expiresAt, DateTime.utc(2026, 10, 19, 12));
    expect(FoodBuff.fromActive(null), isNull);
  });
}
