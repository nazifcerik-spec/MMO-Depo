import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/ranking/ranking_screen.dart';
import 'package:otter/features/social/social_repository.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

Map<String, dynamic> _row(int rank, String id, String name, int value) => {
  'rank': rank,
  'character_id': id,
  'name': name,
  'class_code': 'warrior',
  'level': 400,
  'value': value,
  'score': value,
  'achievement_points': 2,
  'mastery_tiers': 3,
  'profession_levels': 4,
  'dungeons_cleared': 5,
  'raids_cleared': 6,
};

void main() {
  testWidgets('tables, my rank, the score breakdown and a report from a row (DD-357)', (tester) async {
    final be = FakeBackend()
      ..rankingPage = {
        'entries': [_row(1, 'other', 'Şule', 9000), _row(2, 'me', 'Grim', 4321)],
        'me': _row(2, 'me', 'Grim', 4321),
        'points': {
          'level': 10,
          'achievement': 5,
          'mastery_tier': 25,
          'profession_level': 3,
          'dungeon': 20,
          'raid': 100,
        },
      };
    final api = be.api(MemoryTokenStore());
    await api.login('grim', 'correct-horse');
    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('en'),
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: AppLocalizations.supportedLocales,
        home: RankingScreen(repo: RankingRepository(api, 'me'), social: SocialRepository(api, 'me')),
      ),
    );
    await tester.pumpAndSettle();
    expect(find.text('Şule'), findsOneWidget);
    expect(be.lastRankingQuery?['board'], 'score');
    expect(be.lastRankingQuery?['character_id'], 'me');
    expect(find.byKey(const ValueKey('ranking-me')), findsOneWidget);

    // Another table asks the server for it.
    await tester.tap(find.text('Level').first);
    await tester.pumpAndSettle();
    expect(be.lastRankingQuery?['board'], 'level');

    // My card explains the score.
    await tester.tap(find.byKey(const ValueKey('ranking-me')));
    await tester.pumpAndSettle();
    expect(find.text('How the score adds up'), findsOneWidget);
    expect(find.text('6 × 100'), findsOneWidget, reason: 'raids at 100 points');
    await tester.tapAt(const Offset(10, 10));
    await tester.pumpAndSettle();

    // Holding someone else's row reports them; mine does nothing.
    await tester.longPress(find.text('Şule'));
    await tester.pumpAndSettle();
    expect(find.text('Report Şule'), findsOneWidget);
  });
}
