import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../activity/activity_view.dart' show activityIcon;
import '../workshop/workshop_view.dart' show professionName;

/// An item's tags (DD-339): who can use it from what level, and every way to
/// get it or use it, from `GET /v1/items/{code}/info`.
class ItemInfo {
  const ItemInfo({
    this.level = 1,
    this.classes = const [],
    this.anyClass = false,
    this.mobs = const [],
    this.dungeons = const [],
    this.crafted = const [],
    this.vendors = const [],
    this.gather,
    this.usedBy = const [],
  });

  factory ItemInfo.fromJson(Map<String, dynamic> j) {
    List<Map<String, dynamic>> rows(String k) => (j[k] as List? ?? const []).cast<Map<String, dynamic>>();
    List<String> codes(String k) => (j[k] as List? ?? const []).cast<String>();
    final g = j['gather'] as Map<String, dynamic>?;
    return ItemInfo(
      level: (j['level'] as num?)?.toInt() ?? 1,
      classes: codes('classes'),
      anyClass: j['any_class'] as bool? ?? false,
      mobs: [for (final m in rows('mobs')) (code: m['code'] as String, level: (m['level'] as num).toInt())],
      dungeons: [for (final d in rows('dungeons')) (code: d['code'] as String, raid: d['raid'] as bool? ?? false)],
      crafted: [
        for (final c in rows('crafted')) (profession: c['profession'] as String, skill: (c['skill'] as num).toInt()),
      ],
      vendors: codes('vendors'),
      gather: g == null
          ? null
          : (kind: g['kind'] as String, level: (g['level'] as num).toInt(), region: g['region'] as String? ?? ''),
      usedBy: codes('used_by'),
    );
  }

  final int level;
  final List<String> classes;
  final bool anyClass;
  final List<({String code, int level})> mobs;
  final List<({String code, bool raid})> dungeons;
  final List<({String profession, int skill})> crafted;
  final List<String> vendors;
  final ({String kind, int level, String region})? gather;
  final List<String> usedBy;

  /// Content does not change while the app runs: one request per item.
  static final _cache = <String, Future<ItemInfo?>>{};

  static Future<ItemInfo?> load(ApiClient api, String code) => _cache[code] ??= () async {
    try {
      return ItemInfo.fromJson(await api.get('/v1/items/$code/info'));
    } on ApiError {
      _cache.remove(code); // try again next time the item is opened
      return null;
    }
  }();
}

/// A class or branch name: base classes are localized, branches come from
/// the content names.
String _classLabel(AppLocalizations l, String code) {
  final base = className(l, code);
  return base != code ? base : ContentNames.of('class', code);
}

/// The tags under an item's details, grouped: who uses it, where it drops,
/// who crafts or sells it, where it grows and which trades use it. Shows
/// nothing until (and unless) the info loads.
class ItemTags extends StatelessWidget {
  const ItemTags({super.key, required this.api, required this.code, this.equipment = false});

  final ApiClient api;
  final String code;

  /// Equipment shows who may wear it.
  final bool equipment;

  @override
  Widget build(BuildContext context) => FutureBuilder<ItemInfo?>(
    future: ItemInfo.load(api, code),
    builder: (context, snap) {
      final info = snap.data;
      if (info == null) return const SizedBox.shrink();
      final l = AppLocalizations.of(context);
      final groups = <(String, List<Widget>)>[
        (
          l.tagUse,
          [
            _Tag(icon: Icons.trending_up, text: l.levelShort(info.level), color: Palette.accent),
            if (equipment && info.anyClass) _Tag(icon: Icons.groups, text: l.tagAllClasses, color: Palette.accent),
            if (equipment)
              for (final c in info.classes) _Tag(icon: Icons.person, text: _classLabel(l, c), color: Palette.accent),
          ],
        ),
        (
          l.tagDrops,
          [
            for (final m in info.mobs)
              _Tag(icon: Icons.pets, text: '${ContentNames.of('mob', m.code)} · ${l.levelShort(m.level)}'),
            for (final d in info.dungeons)
              _Tag(
                icon: d.raid ? Icons.groups_3 : Icons.castle_outlined,
                text: '${d.raid ? l.tagRaid : l.tagDungeon}: ${ContentNames.of('dungeon', d.code)}',
              ),
          ],
        ),
        (
          l.tagCrafted,
          [
            for (final c in info.crafted)
              _Tag(icon: Icons.handyman, text: '${professionName(l, c.profession)} · ${l.tagSkill(c.skill)}'),
          ],
        ),
        (l.tagSold, [for (final v in info.vendors) _Tag(icon: Icons.storefront, text: ContentNames.of('vendor', v))]),
        (
          l.tagGathered,
          [
            if (info.gather case final g?)
              _Tag(
                icon: activityIcon(g.kind),
                text: [
                  '${activityName(l, g.kind)} ${l.levelShort(g.level)}',
                  if (g.region.isNotEmpty) ContentNames.of('region', g.region),
                ].join(' · '),
              ),
          ],
        ),
        (l.tagUsedIn, [for (final p in info.usedBy) _Tag(icon: Icons.science_outlined, text: professionName(l, p))]),
      ];
      final text = Theme.of(context).textTheme;
      return Column(
        key: const Key('item-tags'),
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          for (final (title, tags) in groups)
            if (tags.isNotEmpty) ...[
              const SizedBox(height: Space.sm),
              Text(title, style: text.bodySmall),
              const SizedBox(height: Space.xs),
              Wrap(spacing: Space.xs, runSpacing: Space.xs, children: tags),
            ],
        ],
      );
    },
  );
}

class _Tag extends StatelessWidget {
  const _Tag({required this.icon, required this.text, this.color = Palette.textSecondary});

  final IconData icon;
  final String text;
  final Color color;

  @override
  Widget build(BuildContext context) => Container(
    padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 4),
    decoration: BoxDecoration(
      color: Palette.bgElevated,
      borderRadius: BorderRadius.circular(Radii.sm),
      border: Border.all(color: color.withValues(alpha: 0.4)),
    ),
    child: Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: color),
        const SizedBox(width: 4),
        Flexible(
          child: Text(
            text,
            style: TextStyle(fontSize: 12, color: color == Palette.textSecondary ? Palette.textPrimary : color),
          ),
        ),
      ],
    ),
  );
}
