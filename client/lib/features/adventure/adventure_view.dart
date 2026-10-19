import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/config.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/loot_sheet.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../events/event_repository.dart';
import '../events/events_panel.dart';
import 'adventure_models.dart';
import 'adventure_repository.dart';
import '../../shared/content_names.dart';

part 'raid_panel.dart';

/// Adventure tab: hunts, dungeons, travel and events behind one segmented control.
/// Panels are built on first selection and kept alive.
class AdventureView extends StatefulWidget {
  const AdventureView({super.key, required this.repo, required this.dataVersion, required this.level, this.live});

  final AdventureRepository repo;
  final ValueNotifier<int> dataVersion;

  /// Push events (raid lobby changes); null in tests.
  final LiveEvents? live;

  /// Live character level; panels rebuild with it on a level-up.
  final ValueListenable<int> level;

  @override
  State<AdventureView> createState() => _AdventureViewState();
}

class _AdventureViewState extends State<AdventureView> {
  int _seg = 0;
  final _visited = <int>{0};

  /// Bumped when the character arrives somewhere: the hunt list follows (DD-273).
  final _moved = ValueNotifier<int>(0);

  @override
  void dispose() {
    _moved.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => ValueListenableBuilder<int>(valueListenable: widget.level, builder: _build);

  Widget _build(BuildContext context, int level, Widget? _) {
    final l = AppLocalizations.of(context);
    final panels = <Widget Function()>[
      () => _HuntPanel(repo: widget.repo, dataVersion: widget.dataVersion, level: level, moved: _moved),
      () => _DungeonsTab(repo: widget.repo, dataVersion: widget.dataVersion, level: level, live: widget.live),
      () => _TravelPanel(repo: widget.repo, level: level, moved: _moved),
      () => EventsPanel(repo: EventRepository(widget.repo.api, widget.repo.characterId), level: level),
    ];
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.sm),
          child: SizedBox(
            width: double.infinity,
            child: SegmentedButton<int>(
              showSelectedIcon: false,
              segments: [
                ButtonSegment(value: 0, label: Text(l.segHunt), icon: const Icon(Icons.gps_fixed, size: 18)),
                ButtonSegment(value: 1, label: Text(l.segDungeon), icon: const Icon(Icons.castle_outlined, size: 18)),
                ButtonSegment(value: 2, label: Text(l.segTravel), icon: const Icon(Icons.map_outlined, size: 18)),
                ButtonSegment(value: 3, label: Text(l.segEvents), icon: const Icon(Icons.public, size: 18)),
              ],
              selected: {_seg},
              onSelectionChanged: (s) => setState(() {
                _seg = s.first;
                _visited.add(_seg);
              }),
            ),
          ),
        ),
        Expanded(
          child: IndexedStack(
            index: _seg,
            children: [for (var i = 0; i < panels.length; i++) _visited.contains(i) ? panels[i]() : const SizedBox()],
          ),
        ),
      ],
    );
  }
}

/// Running contract card with collect/abandon actions.
class _RunningCard extends StatelessWidget {
  const _RunningCard({
    required this.icon,
    required this.title,
    required this.startedAt,
    required this.endsAt,
    required this.doneLabel,
    required this.onCollect,
    required this.onAbandon,
    required this.onComplete,
    this.hint,
    this.collectBeforeEnd = true,
  });

  final IconData icon;
  final String title;
  final DateTime startedAt;
  final DateTime endsAt;
  final String doneLabel;
  final String? hint;

  /// Hunts pay out partial progress; dungeons only at the end.
  final bool collectBeforeEnd;
  final VoidCallback onCollect;
  final VoidCallback onAbandon;
  final VoidCallback onComplete;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final done = ContractProgress.isDone(endsAt);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(Space.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      radius: 22,
                      backgroundColor: Palette.danger.withValues(alpha: 0.15),
                      child: Icon(icon, color: Palette.danger),
                    ),
                    const SizedBox(width: Space.md),
                    Expanded(child: Text(title, style: Theme.of(context).textTheme.titleMedium)),
                  ],
                ),
                const SizedBox(height: Space.lg),
                ContractProgress(startedAt: startedAt, endsAt: endsAt, doneLabel: doneLabel, onComplete: onComplete),
                if (hint != null && !done) ...[
                  const SizedBox(height: Space.sm),
                  Text(hint!, style: Theme.of(context).textTheme.bodySmall),
                ],
              ],
            ),
          ),
        ),
        const SizedBox(height: Space.lg),
        FilledButton(onPressed: done || collectBeforeEnd ? onCollect : null, child: Text(l.collect)),
        const SizedBox(height: Space.sm),
        TextButton(
          onPressed: onAbandon,
          style: TextButton.styleFrom(foregroundColor: Palette.danger),
          child: Text(l.abandon),
        ),
      ],
    );
  }
}

String _resultSubtitle(AppLocalizations l, CombatResult r) => [
  if (r.cleared != null) r.cleared! ? l.dungeonCleared : l.dungeonFailed,
  if (r.bossKilled) l.bossKilled,
  l.huntSummary(r.victories, r.defeats),
  if (r.leveledUp && r.newLevel != null) l.levelUp(r.newLevel!),
].join(' · ');

// --- Hunt ---------------------------------------------------------------------

class _HuntPanel extends StatefulWidget {
  const _HuntPanel({required this.repo, required this.dataVersion, required this.level, required this.moved});

  final AdventureRepository repo;
  final ValueNotifier<int> dataVersion;
  final int level;
  final ValueListenable<int> moved;

  @override
  State<_HuntPanel> createState() => _HuntPanelState();
}

class _HuntPanelState extends State<_HuntPanel> with PanelState {
  List<MobInfo> _mobs = const [];
  CombatContract? _active;
  String? _mob;

  /// Regions in reach from where the character stands (DD-273).
  List<String> _reach = const [];

  /// The player picked [_mob] themselves; otherwise it follows the level.
  bool _mobPicked = false;
  Duration _duration = huntDurations.first;

  List<Duration> get _durations => [if (AppConfig.flavor == 'dev') const Duration(minutes: 5), ...huntDurations];

  @override
  void initState() {
    super.initState();
    widget.moved.addListener(_onMoved);
    load();
  }

  @override
  void dispose() {
    widget.moved.removeListener(_onMoved);
    super.dispose();
  }

  void _onMoved() {
    _mobPicked = false;
    load();
  }

  @override
  void didUpdateWidget(_HuntPanel old) {
    super.didUpdateWidget(old);
    if (old.level != widget.level && !_mobPicked) _mob = _defaultMob();
  }

  bool _inReach(MobInfo m) => m.regionCode.isEmpty || _reach.contains(m.regionCode);

  /// The toughest mob in reach at or below the character's level.
  String? _defaultMob() {
    final here = _mobs.where(_inReach);
    return here.where((m) => m.level <= widget.level).lastOrNull?.code ?? here.firstOrNull?.code;
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object?>([widget.repo.mobs(), widget.repo.activeHunt(), widget.repo.location()]);
    _mobs = r[0]! as List<MobInfo>;
    _active = r[1] as CombatContract?;
    _reach = (r[2]! as Location).huntRegions;
    if (!_mobPicked || !_mobs.any((m) => m.code == _mob && _inReach(m))) _mob = _defaultMob();
  }

  Future<void> _start() => run(() async {
    final c = await widget.repo.startHunt(_mob!, _duration);
    if (mounted) setState(() => _active = c);
  });

  Future<void> _collect(CombatContract c) async {
    final l = AppLocalizations.of(context);
    if (!ContractProgress.isDone(c.expiresAt) &&
        !await confirmAction(context, title: l.collectEarlyTitle, body: l.huntEarlyBody, action: l.collect)) {
      return;
    }
    await run(() async {
      final r = await widget.repo.revealHunt(c.id);
      if (!mounted) return;
      setState(() => _active = null);
      widget.dataVersion.value++;
      await showLootSheet(
        context,
        title: ContentNames.of('mob', c.target),
        subtitle: _resultSubtitle(l, r),
        loot: r.loot,
      );
    });
  }

  Future<void> _abandon(CombatContract c) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.abandonTitle,
      body: l.abandonBody,
      action: l.abandon,
      danger: true,
    );
    if (!ok) return;
    await run(() async {
      await widget.repo.stopHunt(c.id);
      if (mounted) setState(() => _active = null);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return frame(() {
      final a = _active;
      if (a != null) {
        return [
          _RunningCard(
            icon: Icons.gps_fixed,
            title: ContentNames.of('mob', a.target),
            startedAt: a.startedAt,
            endsAt: a.expiresAt,
            doneLabel: l.huntDone,
            onCollect: () => _collect(a),
            onAbandon: () => _abandon(a),
            onComplete: () => setState(() {}),
          ),
        ];
      }
      final text = Theme.of(context).textTheme;
      final here = _mobs.where(_inReach).toList();
      final away = _mobs.where((m) => !_inReach(m)).toList();
      final hunting = [for (final r in _reach) ContentNames.of('region', r)];
      return [
        Text(l.huntChoose, style: text.titleMedium),
        const SizedBox(height: Space.xs),
        Text(
          here.isEmpty ? l.huntNoMobsHere : l.huntReach(hunting.join(', ')),
          key: const Key('hunt-reach'),
          style: text.bodySmall,
        ),
        const SizedBox(height: Space.md),
        for (final m in here) ...[
          _MobTile(
            mob: m,
            selected: m.code == _mob,
            overLevel: m.level > widget.level,
            xpPercent: mobXPPercent(widget.repo.xpBands, widget.level, m.level),
            onTap: () => setState(() {
              _mob = m.code;
              _mobPicked = true;
            }),
          ),
          const SizedBox(height: Space.sm),
        ],
        const SizedBox(height: Space.lg),
        Text(l.durationLabel, style: text.bodySmall),
        const SizedBox(height: Space.sm),
        Wrap(
          spacing: Space.sm,
          runSpacing: Space.sm,
          children: [
            for (final d in _durations)
              ChoiceChip(
                label: Text(shortDuration(l, d)),
                selected: d == _duration,
                onSelected: (_) => setState(() => _duration = d),
                selectedColor: Palette.accent.withValues(alpha: 0.2),
                side: BorderSide(color: d == _duration ? Palette.accent : Palette.borderSubtle),
              ),
          ],
        ),
        const SizedBox(height: Space.xl),
        FilledButton.icon(
          onPressed: _mob == null ? null : _start,
          icon: const Icon(Icons.gps_fixed),
          label: Text(l.startHunt),
        ),
        if (away.isNotEmpty) ...[const SizedBox(height: Space.lg), _AwayMobs(mobs: away)],
      ];
    });
  }
}

/// Where the out-of-reach mobs live, so a player knows where to travel (DD-273).
class _AwayMobs extends StatelessWidget {
  const _AwayMobs({required this.mobs});

  final List<MobInfo> mobs;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final byRegion = <String, List<MobInfo>>{};
    for (final m in mobs) {
      byRegion.putIfAbsent(m.regionCode, () => []).add(m);
    }
    const muted = TextStyle(fontSize: 12, color: Palette.textMuted);
    return ExpansionTile(
      key: const Key('hunt-away'),
      tilePadding: EdgeInsets.zero,
      childrenPadding: const EdgeInsets.only(bottom: Space.md),
      expandedCrossAxisAlignment: CrossAxisAlignment.start,
      title: Text(l.huntOtherRegions, style: Theme.of(context).textTheme.titleSmall),
      subtitle: Text(l.huntTravelHint, style: muted),
      children: [
        for (final e in byRegion.entries) ...[
          const SizedBox(height: Space.sm),
          Text(ContentNames.of('region', e.key), style: const TextStyle(fontWeight: FontWeight.w600)),
          for (final m in e.value) Text('${l.levelShort(m.level)}  ${ContentNames.of('mob', m.code)}', style: muted),
        ],
      ],
    );
  }
}

class _MobTile extends StatelessWidget {
  const _MobTile({
    required this.mob,
    required this.selected,
    required this.overLevel,
    required this.onTap,
    this.xpPercent = 100,
  });

  final MobInfo mob;
  final bool selected;
  final bool overLevel;

  /// Share of full combat XP (DD-254): low-level mobs pay less.
  final int xpPercent;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Semantics(
      selected: selected,
      button: true,
      child: Material(
        color: selected ? Palette.accent.withValues(alpha: 0.12) : Palette.bgSurface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(Radii.md),
          side: BorderSide(color: selected ? Palette.accent : Palette.borderSubtle, width: selected ? 1.5 : 1),
        ),
        child: InkWell(
          borderRadius: BorderRadius.circular(Radii.md),
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.all(Space.md),
            child: Row(
              children: [
                Container(
                  width: 40,
                  alignment: Alignment.center,
                  child: Text(
                    l.levelShort(mob.level),
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w700,
                      color: overLevel ? Palette.danger : Palette.textSecondary,
                    ),
                  ),
                ),
                const SizedBox(width: Space.sm),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(ContentNames.of('mob', mob.code), style: const TextStyle(fontWeight: FontWeight.w600)),
                      Text(
                        l.mobStats(mob.hpMax, mob.attack, mob.defense),
                        style: const TextStyle(fontSize: 11, color: Palette.textMuted, fontFeatures: numeric),
                      ),
                    ],
                  ),
                ),
                Text(
                  l.xpShare(xpPercent),
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: switch (xpPercent) {
                      >= 100 => Palette.epic,
                      >= 50 => Palette.warning,
                      _ => Palette.textMuted,
                    },
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

// --- Dungeon ------------------------------------------------------------------

class _DungeonPanel extends StatefulWidget {
  const _DungeonPanel({required this.repo, required this.dataVersion, required this.level});

  final AdventureRepository repo;
  final ValueNotifier<int> dataVersion;
  final int level;

  @override
  State<_DungeonPanel> createState() => _DungeonPanelState();
}

class _DungeonPanelState extends State<_DungeonPanel> with PanelState {
  List<DungeonInfo> _dungeons = const [];
  CombatContract? _active;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object?>([widget.repo.dungeons(), widget.repo.activeDungeon()]);
    _dungeons = r[0]! as List<DungeonInfo>;
    _active = r[1] as CombatContract?;
  }

  Future<void> _enter(DungeonInfo d) => run(() async {
    final c = await widget.repo.startDungeon(d.code);
    if (mounted) setState(() => _active = c);
  });

  Future<void> _collect(CombatContract c) => run(() async {
    final l = AppLocalizations.of(context);
    final r = await widget.repo.revealDungeon(c.id);
    if (!mounted) return;
    setState(() => _active = null);
    widget.dataVersion.value++;
    await showLootSheet(
      context,
      title: ContentNames.of('dungeon', c.target),
      subtitle: _resultSubtitle(l, r),
      loot: r.loot,
    );
  });

  Future<void> _abandon(CombatContract c) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.abandonTitle,
      body: l.abandonBody,
      action: l.abandon,
      danger: true,
    );
    if (!ok) return;
    await run(() async {
      await widget.repo.stopDungeon(c.id);
      if (mounted) setState(() => _active = null);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return frame(() {
      final a = _active;
      if (a != null) {
        return [
          _RunningCard(
            icon: Icons.castle_outlined,
            title: ContentNames.of('dungeon', a.target),
            startedAt: a.startedAt,
            endsAt: a.expiresAt,
            doneLabel: l.dungeonDone,
            hint: l.dungeonWaitHint,
            collectBeforeEnd: false,
            onCollect: () => _collect(a),
            onAbandon: () => _abandon(a),
            onComplete: () => setState(() {}),
          ),
        ];
      }
      final text = Theme.of(context).textTheme;
      return [
        for (final d in _dungeons) ...[
          Card(
            child: Padding(
              padding: const EdgeInsets.all(Space.lg),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Row(
                    children: [
                      Expanded(child: Text(ContentNames.of('dungeon', d.code), style: text.titleMedium)),
                      Text(
                        l.levelReq(d.minLevel),
                        style: TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.w700,
                          color: d.minLevel > widget.level ? Palette.danger : Palette.textSecondary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: Space.xs),
                  Text(
                    '${l.dungeonEncounters(d.encounters, ContentNames.of('mob', d.bossMobCode))} · ${shortDuration(l, d.duration)}',
                    style: text.bodySmall,
                  ),
                  const SizedBox(height: Space.sm),
                  Wrap(
                    spacing: Space.sm,
                    crossAxisAlignment: WrapCrossAlignment.center,
                    children: [
                      Text(l.bonus, style: text.bodySmall),
                      GoldText(d.bonusGold, size: 13),
                      Text(
                        l.xpGain(d.bonusXp),
                        style: const TextStyle(fontSize: 12, color: Palette.epic, fontWeight: FontWeight.w600),
                      ),
                    ],
                  ),
                  const SizedBox(height: Space.md),
                  FilledButton.tonal(
                    onPressed: d.minLevel > widget.level ? null : () => _enter(d),
                    child: Text(l.dungeonEnter),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: Space.md),
        ],
      ];
    });
  }
}

/// Solo dungeons and party raids (DD-261) under one segment.
class _DungeonsTab extends StatefulWidget {
  const _DungeonsTab({required this.repo, required this.dataVersion, required this.level, this.live});

  final AdventureRepository repo;
  final ValueNotifier<int> dataVersion;
  final int level;
  final LiveEvents? live;

  @override
  State<_DungeonsTab> createState() => _DungeonsTabState();
}

class _DungeonsTabState extends State<_DungeonsTab> {
  bool _raids = false;
  bool _raidsSeen = false;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    Widget chip(String label, bool raids) => ChoiceChip(
      label: Text(label),
      selected: _raids == raids,
      onSelected: (_) => setState(() {
        _raids = raids;
        _raidsSeen |= raids;
      }),
    );
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: Space.screen),
          child: Row(
            children: [
              chip(l.dungeonsSolo, false),
              const SizedBox(width: Space.sm),
              chip(l.raidsTab, true),
            ],
          ),
        ),
        Expanded(
          child: IndexedStack(
            index: _raids ? 1 : 0,
            children: [
              _DungeonPanel(repo: widget.repo, dataVersion: widget.dataVersion, level: widget.level),
              if (_raidsSeen)
                _RaidPanel(repo: widget.repo, dataVersion: widget.dataVersion, level: widget.level, live: widget.live)
              else
                const SizedBox(),
            ],
          ),
        ),
      ],
    );
  }
}

// --- Travel -------------------------------------------------------------------

class _TravelPanel extends StatefulWidget {
  const _TravelPanel({required this.repo, required this.level, required this.moved});

  final AdventureRepository repo;
  final int level;

  /// Bumped when the character's region changes.
  final ValueNotifier<int> moved;

  @override
  State<_TravelPanel> createState() => _TravelPanelState();
}

class _TravelPanelState extends State<_TravelPanel> with PanelState {
  Location? _loc;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final before = _loc?.region.code;
    _loc = await widget.repo.location();
    if (before != null && before != _loc!.region.code) widget.moved.value++;
  }

  Future<void> _go(Neighbor n) => run(() async {
    await widget.repo.travel(n.region.code);
    await fetch();
  });

  Future<void> _arrive(TravelInfo t) => run(() async {
    try {
      await widget.repo.arrive(t.id);
    } on ApiError catch (e) {
      if (e.code != 'E_TRAVEL_NOT_ACTIVE') rethrow; // worker already completed it
    }
    await fetch();
  });

  Future<void> _cancel(TravelInfo t) => run(() async {
    await widget.repo.cancelTravel(t.id);
    await fetch();
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return frame(() {
      final loc = _loc!;
      final text = Theme.of(context).textTheme;
      final t = loc.travel;
      return [
        Text(l.youAreIn, style: text.bodySmall),
        const SizedBox(height: Space.xs),
        Row(
          children: [
            const Icon(Icons.place, color: Palette.accent),
            const SizedBox(width: Space.sm),
            Expanded(child: Text(ContentNames.of('region', loc.region.code), style: text.titleLarge)),
            _KindChip(kind: loc.region.kind, safe: loc.region.safe),
          ],
        ),
        const SizedBox(height: Space.xl),
        if (t != null) ...[
          Card(
            child: Padding(
              padding: const EdgeInsets.all(Space.lg),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text(l.travelling(ContentNames.of('region', t.to)), style: text.titleMedium),
                  const SizedBox(height: Space.md),
                  ContractProgress(
                    startedAt: t.startedAt,
                    endsAt: t.arrivesAt,
                    doneLabel: l.arrived,
                    onComplete: () => setState(() {}),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: Space.lg),
          FilledButton(
            onPressed: ContractProgress.isDone(t.arrivesAt) ? () => _arrive(t) : null,
            child: Text(l.arrive),
          ),
          const SizedBox(height: Space.sm),
          TextButton(onPressed: () => _cancel(t), child: Text(l.cancelTravel)),
        ] else ...[
          Text(l.routes, style: text.titleMedium),
          const SizedBox(height: Space.md),
          if (loc.neighbors.isEmpty) Text(l.noRoutes, style: text.bodySmall),
          for (final n in loc.neighbors) ...[
            Card(
              child: ListTile(
                title: Text(ContentNames.of('region', n.region.code)),
                subtitle: Text(
                  [shortDuration(l, n.travel), if (n.region.minLevel > 1) l.levelReq(n.region.minLevel)].join(' · '),
                ),
                leading: _KindChip(kind: n.region.kind, safe: n.region.safe),
                trailing: FilledButton.tonal(
                  style: compactButton,
                  onPressed: n.region.minLevel > widget.level ? null : () => _go(n),
                  child: Text(l.travelTo),
                ),
              ),
            ),
            const SizedBox(height: Space.sm),
          ],
        ],
      ];
    });
  }
}

class _KindChip extends StatelessWidget {
  const _KindChip({required this.kind, required this.safe});

  final String kind;
  final bool safe;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final (icon, label) = switch (kind) {
      'city' => (Icons.location_city, l.kindCity),
      'zone' => (Icons.terrain, l.kindZone),
      'wilderness' => (Icons.forest, l.kindWilderness),
      'dungeon' => (Icons.castle_outlined, l.kindDungeon),
      'continent' => (Icons.public, l.kindContinent),
      _ => (Icons.place_outlined, kind),
    };
    return Tooltip(
      message: safe ? '$label · ${l.safeZone}' : label,
      child: Icon(icon, size: 22, color: safe ? Palette.success : Palette.textSecondary),
    );
  }
}
