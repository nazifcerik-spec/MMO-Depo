import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../activity/activity_view.dart' show activityIcon;
import '../inventory/item_images.dart';
import '../quests/quests_screen.dart' show questName;
import 'adventure_models.dart';

/// Everything found in a region (DD-338), under the map's region details:
/// its monsters, dungeons, raids, gathering resources and quests. Levels the
/// character has not reached are dimmed.
class RegionContentList extends StatelessWidget {
  const RegionContentList({super.key, required this.content, required this.level});

  final RegionContent content;
  final int level;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    if (content.isEmpty) {
      return Text(l.regionEmpty, style: text.bodySmall);
    }
    Widget header(String title) => Padding(
      padding: const EdgeInsets.only(top: Space.md, bottom: Space.xs),
      child: Text(title, style: text.titleSmall?.copyWith(color: Palette.accent)),
    );
    Widget row(Widget art, String name, int lvl, {String? tag}) {
      final locked = lvl > level;
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 3),
        child: Row(
          children: [
            SizedBox.square(dimension: 32, child: Center(child: art)),
            const SizedBox(width: Space.sm),
            Expanded(
              child: Text(
                name,
                style: TextStyle(color: locked ? Palette.textMuted : Palette.textPrimary),
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (tag != null) ...[
              Text(
                tag,
                style: const TextStyle(fontSize: 12, color: Palette.epic, fontWeight: FontWeight.w600),
              ),
              const SizedBox(width: Space.sm),
            ],
            Text(
              l.levelShort(lvl),
              style: TextStyle(
                fontSize: 12,
                fontFeatures: numeric,
                color: locked ? Palette.danger : Palette.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    String category(String c) => switch (c) {
      'main' => l.questsStory,
      'side' => l.questsSide,
      'daily' => l.questsDaily,
      _ => c,
    };

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (content.mobs.isNotEmpty) ...[
          header(l.regionMobs),
          for (final m in content.mobs) row(mobArt(m.code, 32), ContentNames.of('mob', m.code), m.level),
        ],
        for (final (raid, title) in [(false, l.regionDungeons), (true, l.regionRaids)])
          if (content.dungeons.where((d) => d.raid == raid).toList() case final list when list.isNotEmpty) ...[
            header(title),
            for (final d in list) row(dungeonArt(d.code, 32), ContentNames.of('dungeon', d.code), d.level),
          ],
        if (content.gather.isNotEmpty) ...[
          header(l.regionGather),
          for (final g in content.gather)
            row(
              ItemArt(code: g.code, icon: activityIcon(g.kind), color: Palette.textSecondary, size: 28),
              itemName(g.code),
              g.level,
              tag: activityName(l, g.kind),
            ),
        ],
        if (content.quests.isNotEmpty) ...[
          header(l.quests),
          for (final q in content.quests)
            row(
              CircleArt(
                asset: 'assets/quests/${q.code}.webp',
                size: 32,
                fallback: const Icon(Icons.assignment_outlined, size: 18, color: Palette.textMuted),
              ),
              questName(q.code),
              q.level,
              tag: category(q.category),
            ),
        ],
      ],
    );
  }
}
