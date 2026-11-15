import '../../core/api_client.dart';

/// A story quest of a scene (DD-367): state "" (not taken), active,
/// completed or turned_in.
class StoryQuest {
  const StoryQuest(this.code, this.state, {this.optional = false});

  final String code;
  final String state;
  final bool optional;
}

/// An answer of the current scene; [requiresQuest] is set while the quest it
/// needs is not turned in.
class StoryOption {
  const StoryOption(this.code, {this.requiresQuest});

  final String code;
  final String? requiresQuest;
}

class StoryScene {
  const StoryScene({
    required this.code,
    required this.chapter,
    required this.minLevel,
    required this.state,
    required this.extras,
    required this.afterExtras,
    required this.hasAfter,
    required this.quests,
    required this.options,
    required this.ready,
    this.chosen = '',
  });

  factory StoryScene.fromJson(Map<String, dynamic> j) {
    List<String> strings(String k) => [for (final s in (j[k] as List? ?? const [])) '$s'];
    return StoryScene(
      code: j['code'] as String? ?? '',
      chapter: j['chapter'] as String? ?? '',
      minLevel: (j['min_level'] as num?)?.toInt() ?? 1,
      state: j['state'] as String? ?? 'current',
      chosen: j['chosen'] as String? ?? '',
      extras: strings('extras'),
      afterExtras: strings('after_extras'),
      hasAfter: j['after'] == true,
      ready: j['ready'] == true,
      quests: [
        for (final q in (j['quests'] as List? ?? const []).cast<Map<String, dynamic>>())
          StoryQuest(q['code'] as String? ?? '', q['state'] as String? ?? '', optional: q['optional'] == true),
      ],
      options: [
        for (final o in (j['options'] as List? ?? const []).cast<Map<String, dynamic>>())
          StoryOption(o['code'] as String? ?? '', requiresQuest: o['requires_quest'] as String?),
      ],
    );
  }

  final String code;
  final String chapter;
  final int minLevel;

  /// done | current | locked (current, but the level is too low).
  final String state;

  /// The option chosen ("" = continued); done scenes only.
  final String chosen;
  final List<String> extras;
  final List<String> afterExtras;
  final bool hasAfter;
  final List<StoryQuest> quests;
  final List<StoryOption> options;

  /// The required quests are turned in.
  final bool ready;

  bool get done => state == 'done';
}

class StoryChapterState {
  const StoryChapterState(this.code, this.minLevel, this.state);

  final String code;
  final int minLevel;
  final String state; // done | open | locked
}

class StoryJournal {
  const StoryJournal({required this.chapters, required this.scenes, required this.finished, required this.textVersion});

  factory StoryJournal.fromJson(Map<String, dynamic> j) => StoryJournal(
    chapters: [
      for (final c in (j['chapters'] as List? ?? const []).cast<Map<String, dynamic>>())
        StoryChapterState(
          c['code'] as String? ?? '',
          (c['min_level'] as num?)?.toInt() ?? 1,
          c['state'] as String? ?? '',
        ),
    ],
    scenes: [for (final s in (j['scenes'] as List? ?? const []).cast<Map<String, dynamic>>()) StoryScene.fromJson(s)],
    finished: j['finished'] == true,
    textVersion: j['text_version'] as String? ?? '',
  );

  final List<StoryChapterState> chapters;
  final List<StoryScene> scenes;
  final bool finished;
  final String textVersion;
}

class StoryRepository {
  const StoryRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  static String _textKey = '';
  static Map<String, String> _text = const {};

  String get _base => '/v1/character/$characterId/story';

  Future<StoryJournal> journal() async => StoryJournal.fromJson(await api.get(_base));

  /// Passes the current scene with [option] ("" = continue).
  Future<StoryJournal> advance(String scene, String option) async =>
      StoryJournal.fromJson(await api.post('$_base/$scene/advance', {'option': option}));

  /// The scene text for [lang], fetched once per text version.
  Future<Map<String, String>> text(String lang, String version) async {
    final key = '$lang/$version';
    if (key == _textKey && _text.isNotEmpty) return _text;
    final j = await api.get('/v1/story/text?lang=$lang');
    final raw = j['text'];
    _text = raw is Map ? raw.map((k, v) => MapEntry('$k', '$v')) : const {};
    _textKey = key;
    return _text;
  }
}
