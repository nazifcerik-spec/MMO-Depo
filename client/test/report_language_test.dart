import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/core/app_locale.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/character/models.dart';
import 'package:otter/features/social/social_repository.dart';
import 'package:otter/features/social/social_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

Widget _app(Widget Function(BuildContext) child) => MaterialApp(
  localizationsDelegates: const [
    AppLocalizations.delegate,
    GlobalMaterialLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
  ],
  supportedLocales: AppLocalizations.supportedLocales,
  home: Scaffold(body: Builder(builder: child)),
);

void main() {
  testWidgets('a report needs a reason and sends the chat line (DD-355)', (tester) async {
    final be = FakeBackend();
    final repo = SocialRepository(be.api(MemoryTokenStore()), 'char-1');
    await tester.pumpWidget(
      _app(
        (context) => TextButton(
          onPressed: () => showReportDialog(context, repo, 'Şule', messageId: 7),
          child: const Text('open'),
        ),
      ),
    );
    await tester.tap(find.text('open'));
    await tester.pumpAndSettle();
    expect(find.text('Report Şule'), findsOneWidget);
    expect(tester.widget<FilledButton>(find.byKey(const ValueKey('report-send'))).onPressed, isNull);

    await tester.tap(find.byKey(const ValueKey('report-spam')));
    await tester.pump();
    await tester.enterText(find.byKey(const ValueKey('report-note')), '  ads  ');
    await tester.tap(find.byKey(const ValueKey('report-send')));
    await tester.pumpAndSettle();
    expect(be.lastReport, {
      'character_id': 'char-1',
      'target_name': 'Şule',
      'category': 'spam',
      'note': 'ads',
      'chat_message_id': 7,
    });
    expect(find.text('Report sent. Thank you.'), findsOneWidget);
  });

  test('names in every game language pass the first check; the server has the last word', () {
    for (final name in ['Grim', 'Şule', 'Çağrı', 'Иван', '민수', '李明', 'さくら', 'José', 'Hero_42']) {
      expect(characterNamePattern.hasMatch(name), isTrue, reason: name);
    }
    for (final name in ['', 'a', '1bad', '_bad', 'a b', 'a-b', 'seventeen_letters', '😀😀']) {
      expect(characterNamePattern.hasMatch(name), isFalse, reason: name);
    }
  });

  test('seven languages, each with its own translations', () {
    expect(appLanguages.map((l) => l.$1), ['tr', 'en', 'ru', 'ko', 'zh', 'ja', 'es']);
    for (final (code, _) in appLanguages) {
      expect(AppLocalizations.supportedLocales.map((l) => l.languageCode), contains(code));
    }
  });
}
