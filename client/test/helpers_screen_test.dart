import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/helpers/helpers_screen.dart';
import 'package:otter/features/workshop/workshop_repository.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

void main() {
  testWidgets('helpers: hire, send to gather, collect (DD-319)', (tester) async {
    final be = FakeBackend();
    final api = be.api();
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: HelpersScreen(
          repo: HelperRepository(api, 'char-1'),
          workshop: WorkshopRepository(api, 'char-1'),
          dataVersion: ValueNotifier(0),
        ),
      ),
    );
    await tester.pumpAndSettle();
    expect(find.text('Miner · Mining'), findsOneWidget);
    expect(find.text('+15% yield'), findsOneWidget);

    await tester.tap(find.textContaining('Hire ·'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Hire'));
    await tester.pumpAndSettle();
    expect(be.helpers.first['hired'], isTrue);

    await tester.tap(find.text('Send to work'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('1h 00m'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Working,'), findsOneWidget);

    await tester.tap(find.widgetWithText(FilledButton, 'Collect'));
    await tester.pumpAndSettle();
    expect(find.text('1 items, +12 Mining XP'), findsOneWidget);
    expect(find.text('Send to work'), findsOneWidget);
  });
}
