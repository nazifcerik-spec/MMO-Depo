import '../../core/api_client.dart';

/// One stat modifier granted by a mastery tier (`kind` flat | mult_add | …).
class BonusEffect {
  const BonusEffect(this.kind, this.target, this.amount);

  factory BonusEffect.fromJson(Map<String, dynamic> j) =>
      BonusEffect(j['kind'] as String? ?? '', j['target'] as String? ?? '', (j['amount'] as num?)?.toDouble() ?? 0);

  final String kind;
  final String target;
  final double amount;

  bool get percent => kind.startsWith('mult');
}

class MasteryTrack {
  const MasteryTrack({
    required this.track,
    required this.category,
    required this.xp,
    required this.tier,
    required this.maxTier,
    required this.nextTierXp,
    required this.bonuses,
    this.sources = const [],
  });

  factory MasteryTrack.fromJson(Map<String, dynamic> j) => MasteryTrack(
    track: j['track'] as String,
    category: j['category'] as String? ?? '',
    xp: (j['xp'] as num?)?.toInt() ?? 0,
    tier: (j['tier'] as num?)?.toInt() ?? 0,
    maxTier: (j['max_tier'] as num?)?.toInt() ?? 0,
    nextTierXp: (j['next_tier_xp'] as num?)?.toInt() ?? 0,
    bonuses: [
      for (final b in (j['active_bonuses'] as List? ?? const []).cast<Map<String, dynamic>>()) BonusEffect.fromJson(b),
    ],
    sources: (j['sources'] as List? ?? const []).cast<String>(),
  );

  final String track;
  final String category;
  final int xp;
  final int tier;
  final int maxTier;

  /// Cumulative XP at which the next tier unlocks; 0 at max tier.
  final int nextTierXp;
  final List<BonusEffect> bonuses;

  /// The events that raise it (DD-372).
  final List<String> sources;

  double get progress => tier >= maxTier || nextTierXp <= 0 ? 1 : (xp / nextTierXp).clamp(0, 1).toDouble();
}

class Reputation {
  const Reputation({
    required this.faction,
    required this.rep,
    required this.tier,
    required this.tierName,
    required this.dailyLeft,
    this.sources = const [],
  });

  factory Reputation.fromJson(Map<String, dynamic> j) => Reputation(
    faction: j['faction'] as String,
    rep: (j['rep_xp'] as num?)?.toInt() ?? 0,
    tier: (j['tier'] as num?)?.toInt() ?? 0,
    tierName: j['tier_name'] as String? ?? '',
    dailyLeft: (j['daily_gain_left'] as num?)?.toInt() ?? 0,
    sources: (j['sources'] as List? ?? const []).cast<String>(),
  );

  final String faction;
  final int rep;
  final int tier;
  final String tierName;
  final int dailyLeft;

  /// The events that raise it (DD-372).
  final List<String> sources;
}

class ProgressionSummary {
  const ProgressionSummary({
    required this.points,
    required this.activeTitle,
    required this.mastery,
    required this.reputation,
  });

  factory ProgressionSummary.fromJson(Map<String, dynamic> j) => ProgressionSummary(
    points: (j['achievement_points'] as num?)?.toInt() ?? 0,
    activeTitle: j['active_title_code'] as String?,
    mastery: [
      for (final m in (j['mastery'] as List? ?? const []).cast<Map<String, dynamic>>()) MasteryTrack.fromJson(m),
    ],
    reputation: [
      for (final r in (j['reputation'] as List? ?? const []).cast<Map<String, dynamic>>()) Reputation.fromJson(r),
    ],
  );

  final int points;
  final String? activeTitle;
  final List<MasteryTrack> mastery;
  final List<Reputation> reputation;
}

class Achievement {
  const Achievement({
    required this.code,
    required this.category,
    required this.points,
    required this.progress,
    required this.target,
    required this.rewardGold,
    this.event = '',
    this.unlockedAt,
    this.rewardTitle,
  });

  factory Achievement.fromJson(Map<String, dynamic> j) => Achievement(
    code: j['code'] as String,
    category: j['category'] as String? ?? '',
    points: (j['points'] as num?)?.toInt() ?? 0,
    progress: (j['progress'] as num?)?.toInt() ?? 0,
    target: (j['target'] as num?)?.toInt() ?? 1,
    rewardGold: (j['reward_gold_copper'] as num?)?.toInt() ?? 0,
    event: j['event'] as String? ?? '',
    unlockedAt: DateTime.tryParse(j['unlocked_at'] as String? ?? ''),
    rewardTitle: j['reward_title_code'] as String?,
  );

  final String code;
  final String category;
  final int points;
  final int progress;
  final int target;
  final int rewardGold;

  /// What it counts ("mob_killed", "class_mastery:warlord").
  final String event;
  final DateTime? unlockedAt;
  final String? rewardTitle;

  bool get unlocked => unlockedAt != null;
}

/// A title and how it is earned (DD-372).
class TitleGoal {
  const TitleGoal({
    required this.code,
    required this.source,
    required this.event,
    required this.target,
    required this.progress,
    required this.owned,
    required this.active,
    this.achievement = '',
    this.rewardGold = 0,
    this.points = 0,
    this.ownClass = false,
  });

  factory TitleGoal.fromJson(Map<String, dynamic> j) => TitleGoal(
    code: j['code'] as String,
    source: j['source'] as String? ?? 'achievement',
    event: j['event'] as String? ?? '',
    target: (j['target'] as num?)?.toInt() ?? 0,
    progress: (j['progress'] as num?)?.toInt() ?? 0,
    owned: j['owned'] == true,
    active: j['active'] == true,
    achievement: j['achievement'] as String? ?? '',
    rewardGold: (j['reward_gold_copper'] as num?)?.toInt() ?? 0,
    points: (j['points'] as num?)?.toInt() ?? 0,
    ownClass: j['own_class'] == true,
  );

  final String code;

  /// achievement | level | class_mastery
  final String source;
  final String event;
  final int target;
  final int progress;
  final bool owned;
  final bool active;

  /// The achievement that grants it; empty when none does.
  final String achievement;
  final int rewardGold;
  final int points;

  /// A mastery of the character's own class tree.
  final bool ownClass;
}

/// Achievements, mastery, reputation and titles for one character.
class ProgressionRepository {
  const ProgressionRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _c => '/v1/character/$characterId';

  Future<ProgressionSummary> summary() async => ProgressionSummary.fromJson(await api.get('$_c/progression'));

  Future<List<Achievement>> achievements() async => [
    for (final a
        in ((await api.get('$_c/achievements'))['achievements'] as List? ?? const []).cast<Map<String, dynamic>>())
      Achievement.fromJson(a),
  ];

  Future<List<TitleGoal>> titles() async => [
    for (final t in ((await api.get('$_c/titles'))['all'] as List? ?? const []).cast<Map<String, dynamic>>())
      TitleGoal.fromJson(t),
  ];

  /// Empty [code] clears the active title.
  Future<void> setActiveTitle(String code) => api.post('$_c/titles/active', {'title_code': code});
}
