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
    this.regionCode = '',
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
      regionCode: j['region_code'] as String? ?? '',
    );
  }

  /// Where the mob is hunted (DD-273); empty = anywhere.
  final String regionCode;

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
  const Location({required this.region, required this.neighbors, this.travel, this.huntRegions = const []});

  factory Location.fromJson(Map<String, dynamic> j) => Location(
    region: RegionInfo.fromJson(j['region'] as Map<String, dynamic>),
    neighbors: [
      for (final n in (j['neighbors'] as List? ?? const []).cast<Map<String, dynamic>>())
        Neighbor(RegionInfo.fromJson(n), Duration(seconds: (n['travel_seconds'] as num?)?.toInt() ?? 0)),
    ],
    travel: j['travel'] is Map<String, dynamic> ? TravelInfo.fromJson(j['travel'] as Map<String, dynamic>) : null,
    huntRegions: [for (final r in (j['hunt_regions'] as List? ?? const [])) r as String],
  );

  final RegionInfo region;

  /// Regions whose mobs can be hunted from here: this one and its road links (DD-273).
  final List<String> huntRegions;
  final List<Neighbor> neighbors;
  final TravelInfo? travel;
}

/// Raid template (DD-261): a dungeon for a party of [groupMin]..[groupMax].
class RaidInfo {
  const RaidInfo({required this.dungeon, required this.groupMin, required this.groupMax, required this.bossGear});

  factory RaidInfo.fromJson(Map<String, dynamic> j) => RaidInfo(
    dungeon: DungeonInfo.fromJson(j),
    groupMin: (j['group_min'] as num?)?.toInt() ?? 1,
    groupMax: (j['group_max'] as num?)?.toInt() ?? 1,
    bossGear: (j['boss_gear'] as num?)?.toInt() ?? 0,
  );

  final DungeonInfo dungeon;
  final int groupMin;
  final int groupMax;
  final int bossGear;

  String get code => dungeon.code;
}

class RaidMember {
  const RaidMember({
    required this.characterId,
    required this.name,
    required this.classCode,
    required this.level,
    required this.collected,
  });

  factory RaidMember.fromJson(Map<String, dynamic> j) => RaidMember(
    characterId: j['character_id'] as String,
    name: j['name'] as String? ?? '',
    classCode: j['class_code'] as String? ?? '',
    level: (j['level'] as num?)?.toInt() ?? 1,
    collected: j['collected'] as bool? ?? false,
  );

  final String characterId;
  final String name;
  final String classCode;
  final int level;
  final bool collected;
}

/// A party's raid: `forming` (lobby), `running` (finished once [expiresAt]
/// has passed) or `cancelled`.
class RaidRun {
  const RaidRun({
    required this.id,
    required this.code,
    required this.leaderId,
    required this.state,
    required this.members,
    required this.groupMin,
    required this.groupMax,
    this.startedAt,
    this.expiresAt,
  });

  factory RaidRun.fromJson(Map<String, dynamic> j) {
    DateTime? at(String k) => j[k] is String ? DateTime.parse(j[k] as String) : null;
    return RaidRun(
      id: j['id'] as String,
      code: j['raid_code'] as String? ?? '',
      leaderId: j['leader_character_id'] as String? ?? '',
      state: j['state'] as String? ?? '',
      members: [
        for (final m in (j['members'] as List? ?? const []).cast<Map<String, dynamic>>()) RaidMember.fromJson(m),
      ],
      groupMin: (j['group_min'] as num?)?.toInt() ?? 1,
      groupMax: (j['group_max'] as num?)?.toInt() ?? 1,
      startedAt: at('started_at'),
      expiresAt: at('expires_at'),
    );
  }

  final String id;
  final String code;
  final String leaderId;
  final String state;
  final List<RaidMember> members;
  final int groupMin;
  final int groupMax;
  final DateTime? startedAt;
  final DateTime? expiresAt;

  bool get forming => state == 'forming';
  bool get running => state == 'running';
}

/// What the raid screen shows: the raid the character is in, or its party's
/// open lobby, plus this week's cleared raids.
class RaidCurrent {
  const RaidCurrent({this.raid, this.lobby, required this.inParty, required this.isLeader, required this.lockouts});

  factory RaidCurrent.fromJson(Map<String, dynamic> j) => RaidCurrent(
    raid: j['raid'] is Map<String, dynamic> ? RaidRun.fromJson(j['raid'] as Map<String, dynamic>) : null,
    lobby: j['lobby'] is Map<String, dynamic> ? RaidRun.fromJson(j['lobby'] as Map<String, dynamic>) : null,
    inParty: j['in_party'] as bool? ?? false,
    isLeader: j['is_leader'] as bool? ?? false,
    lockouts: {...(j['lockouts'] as List? ?? const []).cast<String>()},
  );

  final RaidRun? raid;
  final RaidRun? lobby;
  final bool inParty;
  final bool isLeader;
  final Set<String> lockouts;
}

/// A member's raid collection as a [CombatResult]: the shared fight's outcome
/// with the member's own reward.
CombatResult raidResult(Map<String, dynamic> j) {
  final r = j['result'] as Map<String, dynamic>? ?? const {};
  final rw = j['reward'] as Map<String, dynamic>? ?? const {};
  int n(Map<String, dynamic> m, String k) => (m[k] as num?)?.toInt() ?? 0;
  final items = [for (final i in (j['items'] as List? ?? const [])) LootItem.fromJson(i as Map<String, dynamic>)];
  final xp = j['xp'] as Map<String, dynamic>?;
  return CombatResult(
    loot: Loot(
      goldCopper: n(rw, 'gold_copper'),
      xp: n(rw, 'xp'),
      items: items,
      counts: items.isEmpty
          ? {'common': n(rw, 'common_drops'), 'rare': n(rw, 'rare_drops'), 'legendary': n(rw, 'legendary_drops')}
          : const {},
    ),
    victories: n(r, 'victories'),
    defeats: r['wiped'] == true ? 1 : 0,
    leveledUp: xp?['leveled_up'] as bool? ?? false,
    newLevel: (xp?['level'] as num?)?.toInt(),
    cleared: r['cleared'] as bool? ?? false,
    bossKilled: r['boss_killed'] as bool? ?? false,
  );
}
