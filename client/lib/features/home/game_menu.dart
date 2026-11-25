import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart';
import '../../shared/widgets.dart';
import '../character/models.dart';
import '../player/player_repository.dart';
import '../shop/cosmetics.dart';
import '../shop/shop_screen.dart' show CosmeticProfile;
import '../social/group_panels.dart' show memberActivity;

/// One place of the side menu. [code] also names its painted icon,
/// `assets/menu/<code>.webp` (`scripts/art_sets.py` `menu`).
class MenuPlace {
  const MenuPlace(this.code, this.icon, this.label, this.onOpen, {this.badge = 0, this.badgeColor});

  final String code;
  final IconData icon;
  final String label;
  final VoidCallback onOpen;
  final int badge;
  final Color? badgeColor;
}

/// The painted icon of menu place [code]; the line icon while it is missing.
Widget menuArt(String code, IconData fallback, double size, Color color) => CircleArt(
  asset: 'assets/menu/$code.webp',
  size: size,
  fallback: Icon(fallback, size: size * 0.8, color: color),
);

/// What the character is doing: name, live time left, opens its place.
class JobCard extends StatelessWidget {
  const JobCard({super.key, required this.action, required this.onOpen, this.onDone});

  final CurrentAction? action;
  final VoidCallback onOpen;
  final VoidCallback? onDone;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final a = action;
    final text = Theme.of(context).textTheme;
    return Material(
      key: const Key('job-card'),
      color: Palette.bgElevated,
      borderRadius: BorderRadius.circular(Radii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onOpen,
        child: Padding(
          padding: const EdgeInsets.all(Space.md),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l.jobNow, style: text.bodySmall),
              const SizedBox(height: Space.xs),
              if (a == null)
                Text(l.jobIdle, style: text.bodyMedium?.copyWith(color: Palette.textSecondary))
              else ...[
                Text(
                  memberActivity(l, (kind: a.kind, detail: a.detail)),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                if ((a.startedAt, a.endsAt) case (final from?, final to?)) ...[
                  const SizedBox(height: Space.sm),
                  ContractProgress(
                    // A new job starts its own clock.
                    key: ValueKey(to),
                    startedAt: from,
                    endsAt: to,
                    doneLabel: l.jobReady,
                    onComplete: onDone,
                  ),
                ],
              ],
            ],
          ),
        ),
      ),
    );
  }
}

/// The side menu on wide screens (DD-370): who is playing, what they are
/// doing, and every place of the game under one roof. [compact] shows icons
/// only.
class GameMenu extends StatelessWidget {
  const GameMenu({
    super.key,
    required this.groups,
    required this.current,
    required this.sheet,
    required this.profile,
    required this.gold,
    required this.job,
    required this.onJob,
    required this.onJobDone,
    required this.compact,
  });

  /// Title (empty = none) and places.
  final List<(String, List<MenuPlace>)> groups;

  /// Code of the open place.
  final String current;
  final CharacterSheet sheet;
  final CosmeticProfile? profile;
  final int? gold;
  final CurrentAction? job;
  final VoidCallback onJob;
  final VoidCallback onJobDone;
  final bool compact;

  static const width = 264.0, compactWidth = 76.0;

  /// Windows with room for icons only.
  static bool isCompact(BuildContext context) => MediaQuery.sizeOf(context).width < 1180;

  @override
  Widget build(BuildContext context) {
    final accent = Theme.of(context).colorScheme.primary;
    return Container(
      key: const Key('game-menu'),
      width: compact ? compactWidth : width,
      decoration: const BoxDecoration(
        color: Palette.bgSurface,
        border: Border(right: BorderSide(color: Palette.borderSubtle)),
      ),
      child: SafeArea(
        right: false,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            compact ? _face(accent) : _head(context, accent),
            if (!compact)
              Padding(
                padding: const EdgeInsets.fromLTRB(Space.md, 0, Space.md, Space.md),
                child: JobCard(action: job, onOpen: onJob, onDone: onJobDone),
              ),
            const Divider(height: 1, color: Palette.borderSubtle),
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(vertical: Space.sm),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    for (final (i, (title, places)) in groups.indexed) ...[
                      if (i > 0) _rule(context, title, accent),
                      for (final p in places)
                        _Row(
                          place: p,
                          selected: p.code == current,
                          compact: compact,
                          accent: accent,
                          art: menuArt(p.code, p.icon, 28, p.code == current ? accent : Palette.textSecondary),
                        ),
                    ],
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// A group's name on a fading gold thread.
  Widget _rule(BuildContext context, String title, Color accent) => Padding(
    padding: EdgeInsets.fromLTRB(compact ? Space.md : Space.lg, Space.md, compact ? Space.md : 0, Space.xs),
    child: Row(
      children: [
        if (!compact) ...[
          Text(
            title,
            style: Theme.of(context).textTheme.labelSmall
                ?.copyWith(color: accent, fontSize: 11, letterSpacing: 0.6, fontWeight: FontWeight.w700),
          ),
          const SizedBox(width: Space.sm),
        ],
        Expanded(
          child: Container(
            height: 1,
            decoration: BoxDecoration(
              gradient: LinearGradient(colors: [accent.withValues(alpha: 0.5), accent.withValues(alpha: 0)]),
            ),
          ),
        ),
      ],
    ),
  );

  Widget _face(Color accent) => Padding(
    padding: const EdgeInsets.symmetric(vertical: Space.md),
    child: Center(child: _avatar(accent, 44)),
  );

  Widget _avatar(Color accent, double size) => Container(
    padding: const EdgeInsets.all(2),
    decoration: BoxDecoration(shape: BoxShape.circle, color: accent),
    child: AvatarBadge(code: profile?.avatar ?? 'av01', size: size),
  );

  Widget _head(BuildContext context, Color accent) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final s = sheet;
    final hp = s.resources['hp_current'] ?? 0, hpMax = s.resources['hp_max'] ?? 0;
    return Padding(
      padding: const EdgeInsets.all(Space.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            children: [
              _avatar(accent, 48),
              const SizedBox(width: Space.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(s.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: text.titleMedium),
                    Text(
                      '${l.levelShort(s.level)} · ${ContentNames.of('class', s.classCode, className(l, s.classCode))}',
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: text.bodySmall,
                    ),
                    if (gold case final g?) GoldText(g, size: 13),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: Space.md),
          StatBar(
            label: l.hp,
            value: hpMax == 0 ? 0 : hp / hpMax,
            color: Palette.danger,
            valueText: '$hp/$hpMax',
            stacked: true,
          ),
          const SizedBox(height: Space.sm),
          StatBar(
            label: l.xp,
            value: s.xpToNext <= 0 ? 1 : s.xp / s.xpToNext,
            color: accent,
            valueText: s.xpToNext <= 0 ? '' : '${(s.xp * 100 / s.xpToNext).clamp(0, 100).floor()}%',
            stacked: true,
          ),
        ],
      ),
    );
  }
}

class _Row extends StatelessWidget {
  const _Row({
    required this.place,
    required this.selected,
    required this.compact,
    required this.accent,
    required this.art,
  });

  final MenuPlace place;
  final bool selected;
  final bool compact;
  final Color accent;
  final Widget art;

  @override
  Widget build(BuildContext context) {
    final badge = place.badge <= 0
        ? null
        : Badge(backgroundColor: place.badgeColor, label: Text(place.badge > 99 ? '99+' : '${place.badge}'));
    final row = Semantics(
      button: true,
      selected: selected,
      label: place.label,
      excludeSemantics: true,
      child: Material(
        color: selected ? accent.withValues(alpha: 0.14) : Colors.transparent,
        child: InkWell(
          key: ValueKey('menu-${place.code}'),
          onTap: place.onOpen,
          child: SizedBox(
            height: 40,
            child: Row(
              children: [
                // The open place's mark.
                Container(width: 3, color: selected ? accent : Colors.transparent),
                if (compact)
                  Expanded(
                    child: Stack(
                      alignment: Alignment.center,
                      children: [
                        art,
                        if (badge != null) Positioned(right: Space.sm, top: 2, child: badge),
                      ],
                    ),
                  )
                else ...[
                  const SizedBox(width: Space.md),
                  art,
                  const SizedBox(width: Space.md),
                  Expanded(
                    child: Text(
                      place.label,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: selected ? accent : Palette.textPrimary,
                        fontWeight: selected ? FontWeight.w700 : FontWeight.w500,
                      ),
                    ),
                  ),
                  ?badge,
                  const SizedBox(width: Space.md),
                ],
              ],
            ),
          ),
        ),
      ),
    );
    return compact ? Tooltip(message: place.label, child: row) : row;
  }
}
