import 'package:flutter/material.dart';

import '../../core/api_client.dart';
import '../../shared/widgets.dart';

/// One node of the class tree (DD-253): a base class (tier 0), a promotion
/// (tier 1, level 100) or a final specialization (tier 2, level 300).
class ClassNode {
  const ClassNode({
    required this.code,
    required this.tier,
    this.parent = '',
    this.requiredLevel = 0,
    this.roleTags = const [],
    this.resourceKind = '',
    this.resourceMax = 0,
    this.resourceStartPct = 0,
    this.resourcePerHit = 0,
    this.resourcePerTurn = 0,
    this.baseStats = const {},
    this.primaryAttribute = '',
    this.secondaryAttributes = const [],
    this.equipment = const [],
    this.masteryTitle = '',
  });

  factory ClassNode.fromJson(Map<String, dynamic> j) {
    final res = j['resource'] as Map<String, dynamic>? ?? const {};
    int n(Map<String, dynamic> m, String k) => (m[k] as num?)?.toInt() ?? 0;
    List<String> strings(Object? v) => [for (final s in (v as List? ?? const [])) '$s'];
    return ClassNode(
      code: j['code'] as String,
      tier: (j['tier'] as num?)?.toInt() ?? 0,
      parent: j['parent'] as String? ?? '',
      requiredLevel: (j['required_level'] as num?)?.toInt() ?? 0,
      roleTags: strings(j['role_tags']),
      resourceKind: res['kind'] as String? ?? '',
      resourceMax: n(res, 'max'),
      resourceStartPct: n(res, 'start_pct'),
      resourcePerHit: n(res, 'per_hit'),
      resourcePerTurn: n(res, 'per_turn'),
      baseStats: {
        for (final e in (j['base_stats'] as Map<String, dynamic>? ?? const {}).entries) e.key: (e.value as num).toInt(),
      },
      primaryAttribute: j['primary_attribute'] as String? ?? '',
      secondaryAttributes: strings(j['secondary_attributes']),
      equipment: strings(j['equipment']),
      masteryTitle: j['mastery_title'] as String? ?? '',
    );
  }

  final String code;
  final int tier;
  final String parent;
  final int requiredLevel;
  final List<String> roleTags;
  final String resourceKind;
  final int resourceMax;
  final int resourceStartPct;
  final int resourcePerHit;
  final int resourcePerTurn;

  /// Raw base attributes keyed by attribute name (strength, vitality, ...).
  final Map<String, int> baseStats;
  final String primaryAttribute;
  final List<String> secondaryAttributes;

  /// Effective equipment permissions (families and `<type>_armor`).
  final List<String> equipment;
  final String masteryTitle;
}

/// The class tree, fetched once per app session (content; changes only on
/// a content reload).
class ClassCatalog {
  ClassCatalog(this.nodes, {this.attributePointsTotal = 0, this.milestones = const {}});

  factory ClassCatalog.fromJson(Map<String, dynamic> j) => ClassCatalog(
    [for (final c in (j['classes'] as List? ?? const []).cast<Map<String, dynamic>>()) ClassNode.fromJson(c)],
    attributePointsTotal: (j['attribute_points_total'] as num?)?.toInt() ?? 0,
    milestones: {
      for (final e in (j['milestones'] as Map<String, dynamic>? ?? const {}).entries) e.key: (e.value as num).toInt(),
    },
  );

  final List<ClassNode> nodes;

  /// Attribute points a character earns by the level cap.
  final int attributePointsTotal;

  /// awakening / ultimate / capstone / mastery levels.
  final Map<String, int> milestones;

  static ClassCatalog? _cached;

  static Future<ClassCatalog> load(ApiClient api) async =>
      _cached ??= ClassCatalog.fromJson(await api.get('/v1/classes'));

  /// Forgets the cached tree (tests).
  static void reset() => _cached = null;

  ClassNode? node(String code) {
    for (final n in nodes) {
      if (n.code == code) return n;
    }
    return null;
  }

  /// Direct children of a node, in code order (stable layout).
  List<ClassNode> children(String code) => [
    for (final n in nodes)
      if (n.parent == code) n,
  ]..sort((a, b) => a.code.compareTo(b.code));

  /// Highest base attribute of any class: the full length of a stat bar.
  int get maxBaseStat {
    var m = 1;
    for (final n in nodes) {
      for (final v in n.baseStats.values) {
        if (v > m) m = v;
      }
    }
    return m;
  }
}

/// Material icon for a base class: the stand-in when its emblem is missing.
IconData classIcon(String root) => switch (root) {
  'warrior' => Icons.shield_outlined,
  'rogue' => Icons.content_cut,
  'ranger' => Icons.my_location,
  'mage' => Icons.auto_awesome,
  'monk' => Icons.sports_martial_arts,
  'cleric' => Icons.healing,
  'paladin' => Icons.security,
  'druid' => Icons.eco_outlined,
  'bard' => Icons.music_note_outlined,
  'shaman' => Icons.local_fire_department_outlined,
  _ => Icons.person_outline,
};

/// Base class [root]'s emblem art (`scripts/emblem_assets.py`).
String classEmblemAsset(String root) => 'assets/class_icons/$root.webp';

/// A class tree node's emblem: the base class's own, promotions' and
/// specializations' under assets/class_nodes (`scripts/art_sets.py`).
String classNodeAsset(ClassNode n) => n.tier == 0 ? classEmblemAsset(n.code) : 'assets/class_nodes/${n.code}.webp';

/// Base class [root]'s emblem in a circle, decoded at the size shown; its
/// icon when the art is missing.
class ClassEmblem extends StatelessWidget {
  const ClassEmblem({super.key, required this.root, this.size = 24});

  final String root;
  final double size;

  @override
  Widget build(BuildContext context) => CircleArt(
    asset: classEmblemAsset(root),
    size: size,
    fallback: Icon(classIcon(root), size: size * 0.6, color: classColor(root)),
  );
}

/// How forgiving a base class is to play (1 easy, 2 medium, 3 hard), from
/// the balance simulations' AFK safety (DD-266).
int classDifficulty(String root) => switch (root) {
  'warrior' || 'paladin' || 'cleric' => 1,
  'rogue' || 'mage' => 3,
  _ => 2,
};

/// One accent hue per base class, used on the creation screen.
Color classColor(String root) => switch (root) {
  'warrior' => const Color(0xFFE0735A),
  'rogue' => const Color(0xFFB4B8C8),
  'ranger' => const Color(0xFF86C97D),
  'mage' => const Color(0xFF8D9CF2),
  'monk' => const Color(0xFFECC055),
  'cleric' => const Color(0xFFDCC078),
  'paladin' => const Color(0xFFF2A64A),
  'druid' => const Color(0xFF5FC9A2),
  'bard' => const Color(0xFFE785CF),
  'shaman' => const Color(0xFF6BBBE8),
  _ => const Color(0xFFE4A756),
};
