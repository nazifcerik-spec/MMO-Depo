import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/story/story_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

Widget _app(Widget home) => MaterialApp(
  localizationsDelegates: AppLocalizations.localizationsDelegates,
  supportedLocales: AppLocalizations.supportedLocales,
  home: home,
);

void main() {
  testWidgets('story mode: coming soon while empty (DD-313)', (tester) async {
    final be = FakeBackend();
    await tester.pumpWidget(_app(StoryScreen(api: be.api(), level: 5)));
    await tester.pumpAndSettle();
    expect(find.textContaining('Story mode is coming soon'), findsOneWidget);
  });

  testWidgets('story mode: chapters, locked by level', (tester) async {
    final be = FakeBackend()
      ..story = [
        {'code': 'prologue', 'min_level': 1},
        {'code': 'the_keep', 'min_level': 20, 'requires': 'prologue'},
      ];
    await tester.pumpWidget(_app(StoryScreen(api: be.api(), level: 5)));
    await tester.pumpAndSettle();
    expect(find.text('Prologue'), findsOneWidget);
    expect(find.text('Opens at level 20'), findsOneWidget);
    await tester.tap(find.text('Prologue'));
    await tester.pumpAndSettle();
    expect(find.widgetWithText(AppBar, 'Prologue'), findsOneWidget);
  });
}
