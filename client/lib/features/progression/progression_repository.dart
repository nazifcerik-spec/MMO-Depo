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
  );

  final String track;
  final String category;
  final int xp;
  final int tier;
  final int maxTier;

  /// Cumulative XP at which the next tier unlocks; 0 at max tier.
  final int nextTierXp;
  final List<BonusEffect> bonuses;

  double get progress => tier >= maxTier || nextTierXp <= 0 ? 1 : (xp / nextTierXp).clamp(0, 1).toDouble();
}

class Reputation {
  const Reputation({
    required this.faction,
    required this.rep,
    required this.tier,
    required this.tierName,
    required this.dailyLeft,
  });

  factory Reputation.fromJson(Map<String, dynamic> j) => Reputation(
    faction: j['faction'] as String,
    rep: (j['rep_xp'] as num?)?.toInt() ?? 0,
    tier: (j['tier'] as num?)?.toInt() ?? 0,
    tierName: j['tier_name'] as String? ?? '',
    dailyLeft: (j['daily_gain_left'] as num?)?.toInt() ?? 0,
  );

  final String faction;
  final int rep;
  final int tier;
  final String tierName;
  final int dailyLeft;
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
    unlockedAt: DateTime.tryParse(j['unlocked_at'] as String? ?? ''),
    rewardTitle: j['reward_title_code'] as String?,
  );

  final String code;
  final String category;
  final int points;
  final int progress;
  final int target;
  final int rewardGold;
  final DateTime? unlockedAt;
  final String? rewardTitle;

  bool get unlocked => unlockedAt != null;
}

class OwnedTitle {
  const OwnedTitle(this.code, {required this.active});

  final String code;
  final bool active;
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

  Future<List<OwnedTitle>> titles() async => [
    for (final t in ((await api.get('$_c/titles'))['titles'] as List? ?? const []).cast<Map<String, dynamic>>())
      OwnedTitle(t['code'] as String, active: t['active'] == true),
  ];

  /// Empty [code] clears the active title.
  Future<void> setActiveTitle(String code) => api.post('$_c/titles/active', {'title_code': code});
}
