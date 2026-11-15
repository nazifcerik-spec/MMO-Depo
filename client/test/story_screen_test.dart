import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/quests/quest_repository.dart';
import 'package:otter/features/story/story_repository.dart';
import 'package:otter/features/story/story_screen.dart';
import 'package:otter/l10n/app_localizations.dart';
import 'package:otter/shared/content_names.dart';

import 'fake_backend.dart';

Future<Widget> _app(FakeBackend be) async {
  final api = be.api();
  await api.login('grim', 'correct-horse');
  return MaterialApp(
    localizationsDelegates: AppLocalizations.localizationsDelegates,
    supportedLocales: AppLocalizations.supportedLocales,
    home: StoryScreen(
      repo: StoryRepository(api, 'c1'),
      quests: QuestRepository(api, 'c1'),
      dataVersion: ValueNotifier(0),
    ),
  );
}

Map<String, dynamic> _scene(
  String code,
  String state, {
  String chapter = 'prologue',
  List<Map<String, dynamic>> quests = const [],
  List<Map<String, dynamic>> options = const [],
  bool ready = true,
  bool after = false,
  List<String> extras = const [],
  String chosen = '',
}) => {
  'code': code,
  'chapter': chapter,
  'min_level': 1,
  'state': state,
  'chosen': chosen,
  'extras': extras,
  'after_extras': <String>[],
  'after': after,
  'quests': quests,
  'options': options,
  'ready': ready,
};

void main() {
  setUp(
    () => ContentNames.use({
      'story.prologue': 'Prologue',
      'story.prologue.body': 'A star falls.',
      'story.chapter_1': 'Chapter I',
      'quest.story_p_bread': 'Bread Money',
    }),
  );

  testWidgets('story mode: coming soon while empty (DD-313)', (tester) async {
    await tester.pumpWidget(await _app(FakeBackend()));
    await tester.pumpAndSettle();
    expect(find.textContaining('Story mode is coming soon'), findsOneWidget);
  });

  testWidgets('story: passed scene folds with the answer, current scene asks (DD-367)', (tester) async {
    final be = FakeBackend()
      ..storyText = {
        'p_star': 'The Falling Star',
        'p_star.body': 'That night **the sky** split.',
        'p_star.x.dwarf': 'Hedda narrows her eyes.',
        'p_voice': 'The Voice',
        'p_voice.body': '**Voice:** *"You are awake."*\n\nHedda looks at you.',
        'p_voice.tell': 'Tell her',
        'p_voice.tell.result': 'Hedda is silent for a long while.',
        'p_voice.hide': 'Hide it',
      }
      ..storyJournal = {
        'chapters': [
          {'code': 'prologue', 'min_level': 1, 'state': 'open'},
          {'code': 'chapter_1', 'min_level': 20, 'state': 'locked'},
        ],
        'scenes': [
          _scene('p_star', 'done', extras: ['dwarf']),
          _scene(
            'p_voice',
            'current',
            options: [
              {'code': 'tell'},
              {'code': 'hide'},
            ],
          ),
        ],
        'finished': false,
        'text_version': 'v1',
      };
    be.storyNext = {
      ...be.storyJournal,
      'scenes': [_scene('p_star', 'done'), _scene('p_voice', 'done', chosen: 'tell')],
    };
    await tester.pumpWidget(await _app(be));
    await tester.pumpAndSettle();

    expect(find.text('Prologue'), findsOneWidget);
    expect(find.text('Opens at level 20'), findsOneWidget, reason: 'locked chapter');
    expect(find.text('The Falling Star'), findsOneWidget, reason: 'passed scene, folded');
    expect(find.text('Hedda narrows her eyes.'), findsNothing);
    expect(find.text('The Voice'), findsOneWidget);
    expect(find.textContaining('Hedda looks at you.'), findsOneWidget);

    await tester.tap(find.text('Tell her'));
    await tester.pumpAndSettle();
    expect(find.text('Your choice'), findsOneWidget, reason: 'confirm dialog');
    await tester.tap(find.text('Choose'));
    await tester.pumpAndSettle();
    expect(be.storyAdvances, ['p_voice tell']);
    expect(
      be.requests.where((r) => r.url.path.endsWith('/quests')).length,
      2,
      reason: 'board reloaded for the next scene',
    );
    expect(find.text('Hedda is silent for a long while.'), findsOneWidget, reason: 'the result sheet');
    await tester.tap(find.widgetWithText(FilledButton, 'Continue'));
    await tester.pumpAndSettle();
    expect(find.text('You chose: Tell her'), findsOneWidget);
  });

  testWidgets('story: scene quests are accepted here and gate the scene', (tester) async {
    final be = FakeBackend()
      ..storyText = {'p_bread': 'Ironhold', 'p_bread.body': 'Bring me stone.', 'p_bread.after': 'Light in the stone.'}
      ..storyJournal = {
        'chapters': [
          {'code': 'prologue', 'min_level': 1, 'state': 'open'},
        ],
        'scenes': [
          _scene(
            'p_bread',
            'current',
            after: true,
            ready: false,
            quests: [
              {'code': 'story_p_bread', 'state': ''},
            ],
          ),
        ],
        'finished': false,
        'text_version': 'v1',
      };
    be.questsAvailable.add({
      'code': 'story_p_bread',
      'category': 'story',
      'state': 'available',
      'objectives': [
        {'event': 'gather_mining', 'count': 10, 'progress': 0},
      ],
      'reward_xp': 1126,
      'reward_gold_copper': 50,
    });
    await tester.pumpWidget(await _app(be));
    await tester.pumpAndSettle();

    expect(find.text('Story quests'), findsOneWidget);
    expect(find.text('Bread Money'), findsOneWidget);
    expect(find.text('Light in the stone.'), findsNothing, reason: 'after text waits for the quests');
    expect(find.widgetWithText(FilledButton, 'Continue'), findsNothing);
    await tester.tap(find.widgetWithText(FilledButton, 'Accept'));
    await tester.pumpAndSettle();
    expect(be.questCalls, contains('accept story_p_bread'));
  });

  test('story text: bold speakers, italic voice', () {
    final spans = StoryText.spans('**Hedda:** "Up." *"Wake."* end');
    expect(spans.map((s) => s.text), ['Hedda:', ' "Up." ', '"Wake."', ' end']);
    expect(spans[0].style?.fontWeight, FontWeight.w700);
    expect(spans[2].style?.fontStyle, FontStyle.italic);
  });
}
