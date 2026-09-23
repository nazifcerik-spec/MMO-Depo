import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import 'character_repository.dart';
import 'models.dart';

/// Character sheet tab (docs/ui/character-sheet.md): resource bars, XP,
/// attributes, combat stats, gold. Reloads on pull-to-refresh and whenever
/// [dataVersion] changes (after a reveal or a relevant push event).
class CharacterSheetView extends StatefulWidget {
  const CharacterSheetView({super.key, required this.repo, required this.initial, required this.dataVersion});

  final CharacterRepository repo;
  final CharacterSheet initial;
  final ValueNotifier<int> dataVersion;

  @override
  State<CharacterSheetView> createState() => _CharacterSheetViewState();
}

class _CharacterSheetViewState extends State<CharacterSheetView> {
  late CharacterSheet _sheet = widget.initial;
  int? _gold;

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(_reload);
    _reload();
  }

  @override
  void dispose() {
    widget.dataVersion.removeListener(_reload);
    super.dispose();
  }

  /// Silent background reload; errors keep the last good values on screen.
  Future<void> _reload() async {
    try {
      await _fetch();
    } on ApiError {
      // Pull-to-refresh surfaces errors; background reloads stay quiet.
    }
  }

  Future<void> _fetch() async {
    final results = await Future.wait([widget.repo.sheet(_sheet.id), widget.repo.goldCopper(_sheet.id)]);
    if (mounted) {
      setState(() {
        _sheet = results[0] as CharacterSheet;
        _gold = results[1] as int;
      });
    }
  }

  Future<void> _refresh() async {
    try {
      await _fetch();
    } on ApiError catch (e) {
      if (mounted && !e.isSessionError) {
        toast(context, Text(errorText(AppLocalizations.of(context), e)));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final s = _sheet;
    final r = s.resources;
    final text = Theme.of(context).textTheme;

    Widget bar(String label, String key, Color color) {
      final cur = r['${key}_current'] ?? 0, max = r['${key}_max'] ?? 0;
      return Padding(
        padding: const EdgeInsets.only(bottom: Space.sm),
        child: StatBar(label: label, value: max == 0 ? 0 : cur / max, color: color, valueText: '$cur/$max'),
      );
    }

    return RefreshIndicator(
      onRefresh: _refresh,
      child: ListView(
        padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      l.levelLine(s.level, className(l, s.classCode), raceName(l, s.raceCode)),
                      style: text.bodyMedium?.copyWith(color: Palette.textSecondary),
                    ),
                  ],
                ),
              ),
              if (_gold != null) Semantics(label: l.gold, child: GoldText(_gold!)),
            ],
          ),
          const SizedBox(height: Space.lg),
          _Section(
            children: [
              bar(l.hp, 'hp', Palette.danger),
              bar(l.mp, 'mp', Palette.info),
              bar(l.stamina, 'stamina', Palette.success),
              bar(l.energy, 'energy', Palette.warning),
              const Divider(height: Space.xl, color: Palette.borderSubtle),
              StatBar(label: l.xp, value: s.xpProgress, color: Palette.epic, valueText: l.xpValue(s.xp, s.xpToNext)),
            ],
          ),
          const SizedBox(height: Space.lg),
          _Section(
            title: l.attributes,
            children: [
              _Grid(
                entries: [
                  (l.statStr, '${s.primary['str'] ?? 0}'),
                  (l.statAgi, '${s.primary['agi'] ?? 0}'),
                  (l.statInt, '${s.primary['int'] ?? 0}'),
                  (l.statVit, '${s.primary['vit'] ?? 0}'),
                  (l.statCha, '${s.primary['cha'] ?? 0}'),
                  (l.statLuck, '${s.primary['luck'] ?? 0}'),
                ],
              ),
            ],
          ),
          const SizedBox(height: Space.lg),
          _Section(
            title: l.combat,
            children: [
              _Grid(
                columns: 2,
                entries: [
                  (l.attack, '${s.derived['attack'] ?? 0}'),
                  (l.defense, '${s.derived['defense'] ?? 0}'),
                  (l.crit, l.percent(s.derived['crit_pct'] ?? 0)),
                  (l.hit, l.percent(s.derived['hit_pct'] ?? 0)),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _Section extends StatelessWidget {
  const _Section({this.title, required this.children});

  final String? title;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            if (title != null) ...[
              Text(
                title!,
                style: const TextStyle(
                  fontSize: 11,
                  letterSpacing: 1.2,
                  fontWeight: FontWeight.w700,
                  color: Palette.textMuted,
                ),
              ),
              const SizedBox(height: Space.md),
            ],
            ...children,
          ],
        ),
      ),
    );
  }
}

/// Label/value tiles in a fixed-column grid (3×2 attributes, 2×2 combat).
class _Grid extends StatelessWidget {
  const _Grid({required this.entries, this.columns = 3});

  final List<(String, String)> entries;
  final int columns;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, c) {
        final w = (c.maxWidth - (columns - 1) * Space.sm) / columns;
        return Wrap(
          spacing: Space.sm,
          runSpacing: Space.sm,
          children: [
            for (final (label, value) in entries)
              Container(
                width: w,
                padding: const EdgeInsets.symmetric(vertical: Space.md, horizontal: Space.md),
                decoration: BoxDecoration(color: Palette.bgElevated, borderRadius: BorderRadius.circular(Radii.md)),
                child: Semantics(
                  label: '$label $value',
                  excludeSemantics: true,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        label,
                        style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: Palette.textMuted),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        value,
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.w700,
                          fontFeatures: numeric,
                          color: Palette.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
          ],
        );
      },
    );
  }
}
