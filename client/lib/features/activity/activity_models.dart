/// Wire models for gathering activities (server/internal/activity).
const activityKinds = ['mining', 'fishing', 'woodcutting', 'herbalism', 'hunting'];

/// AFK run lengths (DD-249): 1, 3 or 6 hours (server cap 6h). The
/// 1-minute option exists for dev builds.
const activityDurations = [Duration(hours: 1), Duration(hours: 3), Duration(hours: 6)];

class LootItem {
  const LootItem({
    required this.instanceId,
    required this.templateCode,
    required this.rarity,
    this.prefix,
    this.suffix,
  });

  factory LootItem.fromJson(Map<String, dynamic> j) => LootItem(
    instanceId: j['instance_id'] as String? ?? '',
    templateCode: j['template_code'] as String? ?? '',
    rarity: j['rarity'] as String? ?? 'common',
    prefix: j['prefix'] as String?,
    suffix: j['suffix'] as String?,
  );

  final String instanceId;
  final String templateCode;
  final String rarity;
  final String? prefix;
  final String? suffix;
}

/// Identical drops (template, rarity, affixes) shown once with a count —
/// "Iron Ore × 5" (DD-242).
class LootStack {
  const LootStack(this.item, this.count);

  final LootItem item;
  final int count;
}

const _rarityRank = {'common': 0, 'uncommon': 1, 'rare': 2, 'epic': 3, 'legendary': 4, 'mythic': 5};

/// Groups drops into stacks: rarest first, then in order of first appearance.
List<LootStack> stackLoot(List<LootItem> items) {
  final order = <String>[];
  final first = <String, LootItem>{};
  final counts = <String, int>{};
  for (final i in items) {
    final key = '${i.templateCode}|${i.rarity}|${i.prefix ?? ''}|${i.suffix ?? ''}';
    if (counts.containsKey(key)) {
      counts[key] = counts[key]! + 1;
    } else {
      order.add(key);
      first[key] = i;
      counts[key] = 1;
    }
  }
  final stacks = [for (final k in order) LootStack(first[k]!, counts[k]!)];
  // List.sort is not stable: tie-break on first appearance.
  final pos = {for (var i = 0; i < stacks.length; i++) stacks[i]: i};
  stacks.sort((a, b) {
    final r = (_rarityRank[b.item.rarity] ?? 0) - (_rarityRank[a.item.rarity] ?? 0);
    return r != 0 ? r : pos[a]! - pos[b]!;
  });
  return stacks;
}

/// Reward summary shared by activity, hunt and dungeon reveals.
class Loot {
  const Loot({this.goldCopper = 0, this.xp = 0, this.items = const [], this.counts = const {}});

  final int goldCopper;
  final int xp;
  final List<LootItem> items;

  /// Rarity → count when the server reports counts but no instances.
  final Map<String, int> counts;

  bool get isEmpty => goldCopper == 0 && xp == 0 && items.isEmpty && counts.values.every((c) => c == 0);

  List<LootStack> get stacks => stackLoot(items);
}

/// The animal a hunting resource comes from: `wolf_pelt` → `wolf`
/// (its picture is `assets/animals/wolf.png`).
String huntAnimal(String nodeCode) => nodeCode.replaceFirst(RegExp(r'_(pelt|hide|feather|scale)$'), '');

/// One resource of a gathering profession (DD-244).
class GatherNode {
  const GatherNode({required this.code, required this.level, required this.tier});

  factory GatherNode.fromJson(Map<String, dynamic> j) => GatherNode(
    code: j['code'] as String,
    level: (j['level'] as num).toInt(),
    tier: (j['tier'] as num?)?.toInt() ?? 0,
  );

  final String code;
  final int level;
  final int tier;

  /// Mirrors item.TierRarity: T0–1 common … T7 mythic.
  String get rarity => switch (tier) {
    >= 7 => 'mythic',
    6 => 'legendary',
    5 => 'epic',
    4 => 'rare',
    >= 2 => 'uncommon',
    _ => 'common',
  };
}

/// A gathering profession's level plus its resources, lowest first.
class GatherSkill {
  const GatherSkill({required this.kind, required this.level, this.xp = 0, this.xpToNext = 0, this.nodes = const []});

  factory GatherSkill.fromJson(Map<String, dynamic> j) => GatherSkill(
    kind: j['kind'] as String,
    level: (j['level'] as num).toInt(),
    xp: (j['xp'] as num?)?.toInt() ?? 0,
    xpToNext: (j['xp_to_next'] as num?)?.toInt() ?? 0,
    nodes: [for (final n in (j['nodes'] as List? ?? const [])) GatherNode.fromJson(n as Map<String, dynamic>)],
  );

  final String kind;
  final int level;
  final int xp;
  final int xpToNext;
  final List<GatherNode> nodes;

  bool unlocks(GatherNode n) => n.level <= level;

  /// The best resource the level unlocks (the server's default).
  GatherNode? get best {
    GatherNode? b;
    for (final n in nodes) {
      if (unlocks(n)) b = n;
    }
    return b;
  }
}

/// Gathering XP a reveal paid.
class GatherProgress {
  const GatherProgress({required this.level, required this.gained, required this.leveledUp});

  final int level;
  final int gained;
  final bool leveledUp;

  static GatherProgress? fromJson(Object? j) => j is Map<String, dynamic>
      ? GatherProgress(
          level: (j['level'] as num?)?.toInt() ?? 1,
          gained: (j['gained'] as num?)?.toInt() ?? 0,
          leveledUp: j['leveled_up'] == true,
        )
      : null;
}

class ActivityYield {
  const ActivityYield({
    required this.completionRatio,
    required this.goldCopper,
    required this.common,
    required this.rare,
    required this.legendary,
    this.items = const [],
    this.target = '',
  });

  factory ActivityYield.fromJson(Map<String, dynamic>? j) {
    j ??= const {};
    int n(String k) => (j![k] as num?)?.toInt() ?? 0;
    return ActivityYield(
      completionRatio: (j['completion_ratio'] as num?)?.toDouble() ?? 0,
      goldCopper: n('gold_copper'),
      common: n('common_count'),
      rare: n('rare_count'),
      legendary: n('legendary_count'),
      items: [for (final i in (j['items'] as List? ?? const [])) LootItem.fromJson(i as Map<String, dynamic>)],
      target: j['target'] as String? ?? '',
    );
  }

  final double completionRatio;
  final int goldCopper;
  final int common;
  final int rare;
  final int legendary;
  final List<LootItem> items;

  /// Gathered resource (DD-244); empty on older contracts.
  final String target;

  int get units => common + rare + legendary;

  Loot toLoot() => Loot(
    goldCopper: goldCopper,
    items: items,
    counts: items.isEmpty ? {'common': common, 'rare': rare, 'legendary': legendary} : const {},
  );
}

class ActivityContract {
  const ActivityContract({
    required this.id,
    required this.kind,
    required this.startedAt,
    required this.expiresAt,
    required this.state,
  });

  factory ActivityContract.fromJson(Map<String, dynamic> j) => ActivityContract(
    id: j['id'] as String,
    kind: j['kind'] as String,
    startedAt: DateTime.parse(j['started_at'] as String),
    expiresAt: DateTime.parse(j['expires_at'] as String),
    state: j['state'] as String,
  );

  final String id;
  final String kind;
  final DateTime startedAt;
  final DateTime expiresAt;
  final String state;

  Duration get duration => expiresAt.difference(startedAt);

  /// Display-only progress from the device clock; rewards use server time.
  double progressAt(DateTime now) {
    final total = duration.inMilliseconds;
    if (total <= 0) return 1;
    return (now.difference(startedAt).inMilliseconds / total).clamp(0, 1).toDouble();
  }

  bool completeAt(DateTime now) => !now.isBefore(expiresAt);
}

/// Active contract + the server's yield-so-far preview.
class ActiveActivity {
  const ActiveActivity(this.contract, this.current);

  final ActivityContract contract;
  final ActivityYield current;
}
