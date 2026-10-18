import '../activity/activity_models.dart';

/// AFK hunt lengths (DD-249), same as gathering: 1, 3 or 6 hours (server cap
/// 6h). Dungeons have a fixed duration per template.
const huntDurations = activityDurations;

/// Combat XP share for a mob [maxBelow] levels under the character
/// (DD-254); maxBelow 0 closes the list.
class MobXPBand {
  const MobXPBand(this.maxBelow, this.factor);

  factory MobXPBand.fromJson(Map<String, dynamic> j) =>
      MobXPBand((j['max_below'] as num?)?.toInt() ?? 0, (j['factor'] as num?)?.toDouble() ?? 1);

  final int maxBelow;
  final double factor;
}

/// Percent of full combat XP a mob pays (display only; the server decides).
int mobXPPercent(List<MobXPBand> bands, int playerLevel, int mobLevel) {
  final gap = playerLevel - mobLevel;
  if (gap <= 0) return 100;
  for (final b in bands) {
    if (b.maxBelow == 0 || gap <= b.maxBelow) return (b.factor * 100).round();
  }
  return 100;
}

class MobInfo {
  const MobInfo({
    required this.code,
    required this.level,
    required this.hpMax,
    required this.attack,
    required this.defense,
    required this.xpReward,
    required this.goldMin,
    required this.goldMax,
  });

  factory MobInfo.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    return MobInfo(
      code: j['code'] as String,
      level: n('level'),
      hpMax: n('hp_max'),
      attack: n('attack'),
      defense: n('defense'),
      xpReward: n('xp_reward'),
      goldMin: n('gold_copper_min'),
      goldMax: n('gold_copper_max'),
    );
  }

  final String code;
  final int level;
  final int hpMax;
  final int attack;
  final int defense;
  final int xpReward;
  final int goldMin;
  final int goldMax;
}

class DungeonInfo {
  const DungeonInfo({
    required this.code,
    required this.minLevel,
    required this.duration,
    required this.encounters,
    required this.bossMobCode,
    required this.bonusGold,
    required this.bonusXp,
  });

  factory DungeonInfo.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    return DungeonInfo(
      code: j['code'] as String,
      minLevel: n('min_level'),
      duration: Duration(seconds: n('duration_seconds')),
      encounters: n('encounters'),
      bossMobCode: j['boss_mob_code'] as String? ?? '',
      bonusGold: n('bonus_gold_copper'),
      bonusXp: n('bonus_xp'),
    );
  }

  final String code;
  final int minLevel;
  final Duration duration;
  final int encounters;
  final String bossMobCode;
  final int bonusGold;
  final int bonusXp;
}

/// Running hunt or dungeon run — both are "target + timed contract".
class CombatContract {
  const CombatContract({required this.id, required this.target, required this.startedAt, required this.expiresAt});

  factory CombatContract.fromJson(Map<String, dynamic> j, String targetKey) => CombatContract(
    id: j['id'] as String,
    target: j[targetKey] as String? ?? '',
    startedAt: DateTime.parse(j['started_at'] as String),
    expiresAt: DateTime.parse(j['expires_at'] as String),
  );

  final String id;

  /// mob_code (hunt) or dungeon_code (dungeon).
  final String target;
  final DateTime startedAt;
  final DateTime expiresAt;
}

/// Reveal outcome for hunts and dungeons.
class CombatResult {
  const CombatResult({
    required this.loot,
    required this.victories,
    required this.defeats,
    required this.leveledUp,
    this.newLevel,
    this.cleared,
    this.bossKilled = false,
  });

  /// [runKey] is `hunt` or `run`; drop counts use hunt's `*_drops` names.
  factory CombatResult.fromJson(Map<String, dynamic> j, String runKey) {
    final r = j[runKey] as Map<String, dynamic>? ?? const {};
    int n(String k) => (r[k] as num?)?.toInt() ?? 0;
    final items = [for (final i in (j['items'] as List? ?? const [])) LootItem.fromJson(i as Map<String, dynamic>)];
    final xp = j['xp'] as Map<String, dynamic>?;
    return CombatResult(
      loot: Loot(
        goldCopper: n('total_gold_copper'),
        xp: n('total_xp'),
        items: items,
        counts: items.isEmpty
            ? {'common': n('common_drops'), 'rare': n('rare_drops'), 'legendary': n('legendary_drops')}
            : const {},
      ),
      victories: n('victories'),
      defeats: n('defeats'),
      leveledUp: xp?['leveled_up'] as bool? ?? false,
      newLevel: (xp?['level'] as num?)?.toInt(),
      cleared: r['cleared'] as bool?,
      bossKilled: r['boss_killed'] as bool? ?? false,
    );
  }

  final Loot loot;
  final int victories;
  final int defeats;
  final bool leveledUp;
  final int? newLevel;

  /// Dungeon only: every wave + boss beaten.
  final bool? cleared;
  final bool bossKilled;
}

class RegionInfo {
  const RegionInfo({required this.code, required this.kind, required this.minLevel, required this.safe});

  factory RegionInfo.fromJson(Map<String, dynamic> j) => RegionInfo(
    code: j['code'] as String,
    kind: j['kind'] as String? ?? '',
    minLevel: (j['min_level'] as num?)?.toInt() ?? 0,
    safe: j['safe'] as bool? ?? false,
  );

  final String code;
  final String kind;
  final int minLevel;
  final bool safe;
}

class Neighbor {
  const Neighbor(this.region, this.travel);

  final RegionInfo region;
  final Duration travel;
}

class TravelInfo {
  const TravelInfo({required this.id, required this.to, required this.startedAt, required this.arrivesAt});

  factory TravelInfo.fromJson(Map<String, dynamic> j) => TravelInfo(
    id: j['id'] as String,
    to: j['to_region_code'] as String,
    startedAt: DateTime.parse(j['started_at'] as String),
    arrivesAt: DateTime.parse(j['arrives_at'] as String),
  );

  final String id;
  final String to;
  final DateTime startedAt;
  final DateTime arrivesAt;
}

class Location {
  const Location({required this.region, required this.neighbors, this.travel});

  factory Location.fromJson(Map<String, dynamic> j) => Location(
    region: RegionInfo.fromJson(j['region'] as Map<String, dynamic>),
    neighbors: [
      for (final n in (j['neighbors'] as List? ?? const []).cast<Map<String, dynamic>>())
        Neighbor(RegionInfo.fromJson(n), Duration(seconds: (n['travel_seconds'] as num?)?.toInt() ?? 0)),
    ],
    travel: j['travel'] is Map<String, dynamic> ? TravelInfo.fromJson(j['travel'] as Map<String, dynamic>) : null,
  );

  final RegionInfo region;
  final List<Neighbor> neighbors;
  final TravelInfo? travel;
}
