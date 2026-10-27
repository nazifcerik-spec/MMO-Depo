import '../inventory/inventory_models.dart' show FoodBuff;

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../progression/progression_screen.dart' show bonusText;
import '../social/group_panels.dart';
import '../social/group_repository.dart';
import 'character_repository.dart';
import 'leveling.dart';
import 'models.dart';
import 'update_tile.dart';
import '../../shared/content_names.dart';

/// Character sheet tab (docs/ui/character-sheet.md): resource bars, XP,
/// attributes, combat stats, gold. Reloads on pull-to-refresh and whenever
/// [dataVersion] changes (after a reveal or a relevant push event).
class CharacterSheetView extends StatefulWidget {
  const CharacterSheetView({
    super.key,
    required this.repo,
    required this.initial,
    required this.dataVersion,
    this.onOpenProgress,
    this.onSheet,
    this.live,
    this.onPartyChat,
  });

  /// Pushes for the party section (DD-301); null = no party section.
  final LiveEvents? live;

  /// Opens the party chat.
  final VoidCallback? onPartyChat;

  final CharacterRepository repo;
  final CharacterSheet initial;
  final ValueNotifier<int> dataVersion;
  final VoidCallback? onOpenProgress;

  /// Every fresh sheet (the shell tracks the level from it).
  final ValueChanged<CharacterSheet>? onSheet;

  @override
  State<CharacterSheetView> createState() => _CharacterSheetViewState();
}

class _CharacterSheetViewState extends State<CharacterSheetView> {
  late CharacterSheet _sheet = widget.initial;
  int? _gold;
  List<ClassSkill>? _skills;
  LevelingSummary? _leveling;
  FoodBuff? _buff;
  late final _levelingRepo = LevelingRepository(widget.repo.api, widget.initial.id);
  late final _groups = GroupRepository(widget.repo.api, widget.initial.id);

  /// Pull-to-refresh reloads the party section too (members' activities).
  final _partyReload = ValueNotifier(0);

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(_reload);
    _reload();
    _loadSkills();
  }

  /// Content only — a failure just hides the section.
  Future<void> _loadSkills() async {
    try {
      final list = await widget.repo.classSkills(_sheet.classCode);
      if (mounted) setState(() => _skills = list);
    } on ApiError {
      // keep hidden
    }
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
      widget.onSheet?.call(_sheet);
    }
    try {
      final lv = await _levelingRepo.summary();
      if (mounted) setState(() => _leveling = lv);
    } on ApiError {
      // Optional card: the sheet stays usable without it.
    }
    try {
      final b = await widget.repo.buff(_sheet.id);
      if (mounted) setState(() => _buff = b);
    } on ApiError {
      // Optional line.
    }
  }

  /// Skills of the character's class path (base class → promotions);
  /// other branches stay hidden.
  bool _onPath(ClassSkill k) {
    final path = _leveling?.path ?? [_sheet.classCode];
    return k.spec.isEmpty || path.contains(k.spec);
  }

  /// After a progression change: reload here and let every tab follow.
  void _changed() {
    widget.dataVersion.value++;
  }

  Future<void> _promote() async {
    final lv = _leveling;
    if (lv == null) return;
    final choice = await showPromotionSheet(context, lv);
    if (choice == null || !mounted) return;
    if (await levelingAction(context, () => _levelingRepo.promote(choice))) _changed();
  }

  Future<void> _attributes() async {
    final lv = _leveling;
    if (lv == null) return;
    final add = await showModalBottomSheet<Map<String, int>>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => AttributeSheet(summary: lv),
    );
    if (add == null || !mounted) return;
    if (await levelingAction(context, () => _levelingRepo.allocate(add))) _changed();
  }

  Future<void> _talents() async {
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => TalentsScreen(repo: _levelingRepo, dataVersion: widget.dataVersion),
      ),
    );
  }

  Future<void> _refresh() async {
    _partyReload.value++;
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
              if (s.resourceKind.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(bottom: Space.sm),
                  child: StatBar(
                    label: ContentNames.of('resource', s.resourceKind).toUpperCase(),
                    value: 1,
                    color: Palette.info,
                    valueText: '${s.resourceMax}',
                  ),
                )
              else
                bar(l.mp, 'mp', Palette.info),
              bar(l.stamina, 'stamina', Palette.success),
              bar(l.energy, 'energy', Palette.warning),
              if (_buff case FoodBuff(expiresAt: final until?) when until.isAfter(DateTime.now()))
                Row(
                  children: [
                    const Icon(Icons.restaurant_outlined, size: 16, color: Palette.success),
                    const SizedBox(width: Space.xs),
                    Expanded(
                      child: Text(
                        l.fedActive(
                          _buff!.attackPct,
                          _buff!.defensePct,
                          shortDuration(l, until.difference(DateTime.now())),
                        ),
                        style: text.bodySmall?.copyWith(color: Palette.success),
                      ),
                    ),
                  ],
                ),
              const Divider(height: Space.xl, color: Palette.borderSubtle),
              if (_leveling case final lv?)
                LevelingCard(summary: lv, onPromote: _promote, onAttributes: _attributes, onTalents: _talents)
              else
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
          if (widget.live case final live?) ...[
            const SizedBox(height: Space.lg),
            _Section(
              title: l.tabParty,
              children: [PartyPanel(repo: _groups, live: live, onChat: widget.onPartyChat, reload: _partyReload)],
            ),
          ],
          if (_skills?.isNotEmpty ?? false) ...[
            const SizedBox(height: Space.lg),
            _Section(
              key: const Key('skills-section'),
              title: l.skills,
              collapsible: true, // folded until opened (DD-301)
              children: [
                for (final k in _skills!)
                  if (_onPath(k)) _SkillRow(skill: k, level: s.level),
              ],
            ),
          ],
          if (widget.onOpenProgress != null) ...[
            const SizedBox(height: Space.lg),
            OutlinedButton.icon(
              onPressed: widget.onOpenProgress,
              icon: const Icon(Icons.emoji_events_outlined),
              label: Text(l.progress),
            ),
          ],
          const UpdateTile(), // in-app updates (DD-304)
        ],
      ),
    );
  }
}

/// "mage_fire_bolt" → "Fire Bolt".
String skillName(String code) => ContentNames.of('skill', code, humanize(code.substring(code.indexOf('_') + 1)));

String skillDescription(AppLocalizations l, ClassSkill k) {
  final authored = ContentNames.of('skill', '${k.code}.desc', '');
  if (authored.isNotEmpty) return authored;
  if (!k.active) return k.modifiers.map((m) => bonusText(l, m)).join(' · ');
  final effect = switch (k.effect) {
    'heal' => l.skillHeal(k.power),
    'guard' => l.skillGuard(k.power, k.durationTurns),
    'empower' => l.skillEmpower(k.power, k.durationTurns),
    _ => l.skillStrike(k.power),
  };
  return [effect, if (k.mpCost > 0) l.skillMp(k.mpCost), l.skillCooldown(k.cooldownTurns)].join(' · ');
}

class _SkillRow extends StatelessWidget {
  const _SkillRow({required this.skill, required this.level});

  final ClassSkill skill;
  final int level;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final locked = skill.unlockLevel > level;
    return Opacity(
      opacity: locked ? 0.6 : 1,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: Space.xs),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(
              locked ? Icons.lock_outline : (skill.active ? Icons.flash_on : Icons.shield_moon_outlined),
              size: 18,
              color: locked ? Palette.textMuted : Palette.accent,
            ),
            const SizedBox(width: Space.sm),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${skillName(skill.code)} · ${skill.active ? l.skillActive : l.skillPassive}',
                    style: const TextStyle(fontWeight: FontWeight.w600),
                  ),
                  Text(skillDescription(l, skill), style: Theme.of(context).textTheme.bodySmall),
                ],
              ),
            ),
            if (locked) Text(l.levelReq(skill.unlockLevel), style: Theme.of(context).textTheme.bodySmall),
          ],
        ),
      ),
    );
  }
}

class _Section extends StatefulWidget {
  const _Section({super.key, this.title, required this.children, this.collapsible = false});

  final String? title;
  final List<Widget> children;

  /// Folded until its title is tapped.
  final bool collapsible;

  @override
  State<_Section> createState() => _SectionState();
}

class _SectionState extends State<_Section> {
  late bool _open = !widget.collapsible;

  @override
  Widget build(BuildContext context) {
    final title = widget.title;
    final label = title == null
        ? null
        : Text(
            title,
            style: const TextStyle(
              fontSize: 11,
              letterSpacing: 1.2,
              fontWeight: FontWeight.w700,
              color: Palette.textMuted,
            ),
          );
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            if (label != null)
              widget.collapsible
                  ? InkWell(
                      onTap: () => setState(() => _open = !_open),
                      child: Row(
                        children: [
                          Expanded(child: label),
                          Icon(_open ? Icons.expand_less : Icons.expand_more, color: Palette.textSecondary),
                        ],
                      ),
                    )
                  : label,
            if (label != null && _open) const SizedBox(height: Space.md),
            if (_open) ...widget.children,
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
