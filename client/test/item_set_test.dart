import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';
import 'package:otter/features/inventory/item_widgets.dart';
import 'package:otter/l10n/app_localizations.dart';

void main() {
  final set = ItemSet.fromJson({
    'code': 'north_keep_siege_plate',
    'pieces': ['undying_helm', 'undying_chestplate', 'undying_chain'],
    'bonuses': [
      {
        'pieces': 2,
        'modifiers': [
          {'kind': 'mult_add', 'target': 'hp_max', 'amount': 0.04},
        ],
      },
      {
        'pieces': 3,
        'modifiers': [
          {'kind': 'flat', 'target': 'dmg_reduction_pct', 'amount': 3},
        ],
      },
    ],
  });

  test('a set counts distinct worn pieces (DD-277)', () {
    expect(set.wornCount(['undying_helm', 'undying_helm', 'sword_iron']), 1);
    expect(set.wornCount(['undying_helm', 'undying_chain']), 2);
    expect(set.bonuses.first.$1, 2);
  });

  testWidgets('set section: name, pieces worn, bonuses by count', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: Scaffold(
          body: SetSection(set: set, worn: const ['undying_helm', 'undying_chestplate']),
        ),
      ),
    );
    expect(find.textContaining('2/3'), findsOneWidget);
    expect(find.text('(2) +4% HP'), findsOneWidget);
    expect(find.text('(3) +3 % Damage reduction'), findsOneWidget);
  });
}
