import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/progression/progression_repository.dart';
import 'package:otter/features/progression/progression_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

void main() {
  final en = lookupAppLocalizations(const Locale('en'));
  final tr = lookupAppLocalizations(const Locale('tr'));

  test('mastery bonus lines', () {
    expect(bonusText(en, const BonusEffect('flat', 'attack', 2)), '+2 Attack');
    expect(bonusText(en, const BonusEffect('mult_add', 'gather_yield_mining', 0.05)), '+5% Mining yield');
    expect(bonusText(en, const BonusEffect('flat', 'gather_yield_fishing', 1)), '+1/h Fishing yield');
    expect(bonusText(en, const BonusEffect('flat', 'unknown_stat', 1.5)), '+1.5 Unknown Stat');
    expect(bonusText(tr, const BonusEffect('flat', 'gather_yield_mining', 1)), startsWith('+1/sa '));
  });
}
