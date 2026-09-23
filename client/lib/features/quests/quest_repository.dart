import '../../core/api_client.dart';

class QuestObjective {
  const QuestObjective(this.event, this.count, this.progress);

  final String event; // progression event kind, e.g. "mob_killed:grey_wolf"
  final int count;
  final int progress;
}

class Quest {
  const Quest({
    required this.code,
    required this.category,
    required this.state,
    required this.objectives,
    required this.rewardXp,
    required this.rewardGold,
    this.expiresAt,
  });

  factory Quest.fromJson(Map<String, dynamic> j) => Quest(
    code: j['code'] as String,
    category: j['category'] as String? ?? '',
    state: j['state'] as String? ?? 'available',
    objectives: [
      for (final o in (j['objectives'] as List? ?? const []).cast<Map<String, dynamic>>())
        QuestObjective(
          o['event'] as String? ?? '',
          (o['count'] as num?)?.toInt() ?? 1,
          (o['progress'] as num?)?.toInt() ?? 0,
        ),
    ],
    rewardXp: (j['reward_xp'] as num?)?.toInt() ?? 0,
    rewardGold: (j['reward_gold_copper'] as num?)?.toInt() ?? 0,
    expiresAt: DateTime.tryParse(j['expires_at'] as String? ?? ''),
  );

  final String code;
  final String category; // main | side | daily
  final String state; // available | active | completed
  final List<QuestObjective> objectives;
  final int rewardXp;
  final int rewardGold;
  final DateTime? expiresAt;
}

class QuestBoard {
  const QuestBoard({required this.open, required this.available, required this.dailyResetsAt});

  factory QuestBoard.fromJson(Map<String, dynamic> j) {
    List<Quest> list(String k) => [
      for (final q in (j[k] as List? ?? const []).cast<Map<String, dynamic>>()) Quest.fromJson(q),
    ];
    return QuestBoard(
      open: list('open'),
      available: list('available'),
      dailyResetsAt: DateTime.tryParse(j['daily_resets_at'] as String? ?? ''),
    );
  }

  final List<Quest> open;
  final List<Quest> available;
  final DateTime? dailyResetsAt;

  int get readyCount => open.where((q) => q.state == 'completed').length;
}

class TurnInResult {
  const TurnInResult(this.xp, this.gold, {this.newLevel});

  final int xp;
  final int gold;

  /// Set when the reward XP levelled the character up.
  final int? newLevel;
}

class QuestRepository {
  const QuestRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _q => '/v1/character/$characterId/quests';

  Future<QuestBoard> board() async => QuestBoard.fromJson(await api.get(_q));

  Future<void> accept(String code) => api.post('$_q/$code/accept', const {});
  Future<void> abandon(String code) => api.post('$_q/$code/abandon', const {});

  Future<TurnInResult> turnIn(String code) async {
    final j = await api.post('$_q/$code/turn-in', const {});
    final xp = j['xp'];
    return TurnInResult(
      (j['reward_xp'] as num?)?.toInt() ?? 0,
      (j['reward_gold_copper'] as num?)?.toInt() ?? 0,
      newLevel: xp is Map && xp['leveled_up'] == true ? (xp['level'] as num?)?.toInt() : null,
    );
  }
}
