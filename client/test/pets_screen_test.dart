import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/adventure/adventure_repository.dart';
import 'package:otter/features/pets/pets_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

void main() {
  testWidgets('pets: incubate an egg, make it active, hunt, collect (DD-328)', (tester) async {
    final be = FakeBackend();
    be.petEggs.add({'item_id': 'egg-1', 'template_code': 'pet_egg_ember_fox', 'species': 'ember_fox', 'quantity': 1});
    final api = be.api();
    await tester.runAsync(() => api.login('grim', 'correct-horse'));
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: PetsScreen(
          repo: PetRepository(api, 'char-1'),
          adventure: AdventureRepository(api, 'char-1'),
          dataVersion: ValueNotifier(0),
        ),
      ),
    );
    await tester.pumpAndSettle();
    expect(find.byKey(const Key('egg-pet_egg_ember_fox')), findsOneWidget);

    await tester.tap(find.text('Incubate'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Incubate').last);
    await tester.pumpAndSettle();
    expect(find.textContaining('In its egg'), findsOneWidget);
    expect(find.text('Send hunting'), findsNothing, reason: 'not hatched yet');

    be.pets.first['hatch_ready_at'] = DateTime.now().toUtc().subtract(const Duration(seconds: 1)).toIso8601String();
    // Hatched: open the screen again.
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: PetsScreen(
          key: const Key('reloaded'),
          repo: PetRepository(api, 'char-1'),
          adventure: AdventureRepository(api, 'char-1'),
          dataVersion: ValueNotifier(0),
        ),
      ),
    );
    await tester.pumpAndSettle();
    await tester.tap(find.text('Make active'));
    await tester.pumpAndSettle();
    expect(be.pets.first['active'], isTrue);
    expect(find.text('Active'), findsOneWidget);
    expect(find.text('Combat assist: +3% damage'), findsOneWidget, reason: 'the active pet fights alongside (DD-341)');

    // The guide opens a species card: evolution, strength, where to find it.
    await tester.tap(find.text('All pets'));
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(
      find.byKey(const Key('pet-guide-ember_fox')),
      200,
      scrollable: find.byType(Scrollable).first,
    );
    await tester.tap(find.byKey(const Key('pet-guide-ember_fox')));
    await tester.pumpAndSettle();
    expect(find.byKey(const Key('pet-species-sheet')), findsOneWidget);
    for (final t in [
      'Evolution',
      'Lvl 34–66',
      'Combat assist: +3% – +12% damage on every hit',
      'Hunts like a level 10 – 1000 monster',
      'Raid bosses of level 430+ · 50% chance',
    ]) {
      expect(find.text(t), findsOneWidget, reason: t);
    }
    Navigator.of(tester.element(find.byKey(const Key('pet-species-sheet')))).pop();
    await tester.pumpAndSettle();

    // Only monsters within its power (level 10 + 5) are offered.
    await tester.tap(find.text('Send hunting'));
    await tester.pumpAndSettle();
    expect(find.byKey(const Key('pet-mob-giant_rat')), findsOneWidget);
    await tester.tap(find.byKey(const Key('pet-mob-giant_rat')));
    await tester.tap(find.byKey(const Key('pet-hunt-send')));
    await tester.pumpAndSettle();
    expect((be.pets.first['hunt'] as Map)['mob_code'], 'giant_rat');

    await tester.tap(find.widgetWithText(FilledButton, 'Collect'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Collect').last); // early: confirm
    await tester.pumpAndSettle();
    expect(find.textContaining('80 wins, +1200 pet XP'), findsOneWidget);
    expect(find.textContaining('reached level 2'), findsOneWidget);
  });
}
