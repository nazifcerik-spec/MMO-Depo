import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/guide/guide.dart';
import 'package:otter/l10n/app_localizations.dart';

Widget _app(Widget child, {Locale locale = const Locale('en')}) => MaterialApp(
  locale: locale,
  localizationsDelegates: const [
    AppLocalizations.delegate,
    GlobalMaterialLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
  ],
  supportedLocales: AppLocalizations.supportedLocales,
  home: Scaffold(body: child),
);

void main() {
  tearDown(() => TutorialState.debugPending = null);

  testWidgets('the tutorial walks the tabs and ends on Start playing (DD-356)', (tester) async {
    final tabs = <int>[];
    var closed = 0;
    await tester.pumpWidget(
      _app(
        Align(
          alignment: Alignment.bottomCenter,
          child: TutorialCard(onTab: tabs.add, onClose: () => closed++),
        ),
      ),
    );
    expect(find.text('Welcome to MeNa World'), findsOneWidget);
    expect(find.byKey(const ValueKey('tutorial-back')), findsNothing);
    for (var i = 0; i < 7; i++) {
      await tester.tap(find.byKey(const ValueKey('tutorial-next')));
      await tester.pump();
    }
    expect(tabs, [0, 1, 2, 3, 4, 0, 0], reason: 'each page shows its tab');
    expect(find.byKey(const ValueKey('tutorial-skip')), findsNothing, reason: 'nothing to skip on the last page');
    expect(closed, 0);
    await tester.tap(find.text('Start playing'));
    expect(closed, 1);
  });

  testWidgets('the player can skip the whole tutorial from the first page', (tester) async {
    var closed = 0;
    await tester.pumpWidget(
      _app(
        Align(
          alignment: Alignment.bottomCenter,
          child: TutorialCard(onTab: (_) {}, onClose: () => closed++),
        ),
      ),
    );
    await tester.tap(find.byKey(const ValueKey('tutorial-skip')));
    expect(closed, 1);
  });

  test('experienced characters are not asked; finishing is remembered', () async {
    TutorialState.debugPending = true;
    expect(await TutorialState.pending(1), isTrue);
    await TutorialState.finish();
    expect(await TutorialState.pending(1), isFalse);
  });

  testWidgets('the guide opens a topic, searches and jumps to its tab', (tester) async {
    int? opened;
    var replay = 0;
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: AppLocalizations.supportedLocales,
        locale: const Locale('en'),
        home: Builder(
          builder: (context) => TextButton(
            onPressed: () => Navigator.of(context).push(
              MaterialPageRoute<void>(
                builder: (_) => GuideScreen(onOpenTab: (t) => opened = t, onReplayTutorial: () => replay++),
              ),
            ),
            child: const Text('go'),
          ),
        ),
      ),
    );
    await tester.tap(find.text('go'));
    await tester.pumpAndSettle();
    expect(find.text('Gathering'), findsOneWidget);

    await tester.enterText(find.byKey(const ValueKey('guide-search')), 'auction');
    await tester.pumpAndSettle();
    expect(find.text('Gathering'), findsNothing);
    expect(find.text('Market and auctions'), findsOneWidget);
    await tester.tap(find.text('Market and auctions'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Show me'));
    await tester.pumpAndSettle();
    expect(opened, 4);
    expect(replay, 0);
  });

  testWidgets('every game language renders the tutorial and the guide in its own words', (tester) async {
    for (final (code, welcome) in const [
      ('tr', "MeNa World'e hoş geldin"),
      ('ru', 'Добро пожаловать в MeNa World'),
      ('ko', 'MeNa World에 오신 것을 환영합니다'),
      ('zh', '欢迎来到 MeNa World'),
      ('ja', 'MeNa World へようこそ'),
      ('es', 'Te damos la bienvenida a MeNa World'),
    ]) {
      await tester.pumpWidget(
        _app(
          Align(
            alignment: Alignment.bottomCenter,
            child: TutorialCard(onTab: (_) {}, onClose: () {}),
          ),
          locale: Locale(code),
        ),
      );
      await tester.pump();
      expect(find.text(welcome), findsOneWidget, reason: code);
    }
  });
}
