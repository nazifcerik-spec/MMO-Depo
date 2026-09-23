/// Wire models for gathering activities (server/internal/activity).
const activityKinds = ['mining', 'fishing', 'woodcutting', 'herbalism', 'hunting'];

/// Server bounds: 60s .. 12h. The 1-minute option exists for dev builds.
const activityDurations = [
  Duration(minutes: 15),
  Duration(minutes: 30),
  Duration(hours: 1),
  Duration(hours: 2),
  Duration(hours: 4),
  Duration(hours: 8),
  Duration(hours: 12),
];

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

/// Reward summary shared by activity, hunt and dungeon reveals.
class Loot {
  const Loot({this.goldCopper = 0, this.xp = 0, this.items = const [], this.counts = const {}});

  final int goldCopper;
  final int xp;
  final List<LootItem> items;

  /// Rarity → count when the server reports counts but no instances.
  final Map<String, int> counts;

  bool get isEmpty => goldCopper == 0 && xp == 0 && items.isEmpty && counts.values.every((c) => c == 0);
}

class ActivityYield {
  const ActivityYield({
    required this.completionRatio,
    required this.goldCopper,
    required this.common,
    required this.rare,
    required this.legendary,
    this.items = const [],
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
    );
  }

  final double completionRatio;
  final int goldCopper;
  final int common;
  final int rare;
  final int legendary;
  final List<LootItem> items;

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
