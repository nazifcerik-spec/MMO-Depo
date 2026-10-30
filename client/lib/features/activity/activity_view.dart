import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/config.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart';
import '../../shared/loot_sheet.dart';
import '../../shared/party_share.dart';
import '../../shared/widgets.dart';
import '../inventory/item_images.dart';
import 'activity_models.dart';
import 'activity_repository.dart';

IconData activityIcon(String kind) => switch (kind) {
  'mining' => Icons.construction,
  'fishing' => Icons.phishing,
  'woodcutting' => Icons.forest,
  'herbalism' => Icons.eco,
  'hunting' => Icons.pets,
  _ => Icons.help_outline,
};

/// Gathering loop: pick kind + resource + duration → progress → collect
/// (reveal). Each kind is a profession levelled by gathering (DD-244).
/// [dataVersion] is bumped after a reveal so other tabs (sheet, bag) reload.
class ActivityView extends StatefulWidget {
  const ActivityView({super.key, required this.repo, required this.characterId, required this.dataVersion});

  final ActivityRepository repo;
  final String characterId;
  final ValueNotifier<int> dataVersion;

  @override
  State<ActivityView> createState() => _ActivityViewState();
}

class _ActivityViewState extends State<ActivityView> {
  ActiveActivity? _active;
  bool _loaded = false;
  Object? _error;
  bool _busy = false;

  String _kind = activityKinds.first;
  Duration _duration = activityDurations.first;

  /// kind → profession; empty until loaded (or if it fails: the server then
  /// picks the best resource itself).
  Map<String, GatherSkill> _skills = const {};

  /// kind → chosen resource code.
  final Map<String, String> _target = {};

  /// Regions in reach (DD-317); null until loaded (then everything counts).
  Set<String>? _reach;

  bool _here(GatherNode n) => n.region.isEmpty || _reach == null || _reach!.contains(n.region);

  List<Duration> get _durations => [if (AppConfig.flavor == 'dev') const Duration(minutes: 1), ...activityDurations];

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(_loadSkills);
    _load();
  }

  @override
  void dispose() {
    widget.dataVersion.removeListener(_loadSkills);
    super.dispose();
  }

  Future<void> _load() async {
    _loadSkills();
    try {
      final a = await widget.repo.active(widget.characterId);
      if (mounted) {
        setState(() {
          _active = a;
          _loaded = true;
          _error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _loadSkills() async {
    try {
      final list = await widget.repo.skills(widget.characterId);
      if (mounted) setState(() => _skills = {for (final g in list) g.kind: g});
      final reach = await widget.repo.reach(widget.characterId);
      if (mounted) setState(() => _reach = reach);
    } on ApiError {
      // Optional: the picker works without levels.
    }
  }

  /// The chosen resource for [kind], falling back to the best unlocked one.
  GatherNode? _nodeFor(String kind) {
    final g = _skills[kind];
    if (g == null) return null;
    final code = _target[kind];
    for (final n in g.nodes) {
      if (n.code == code && g.unlocks(n) && _here(n)) return n;
    }
    GatherNode? best;
    for (final n in g.nodes) {
      if (g.unlocks(n) && _here(n)) best = n;
    }
    return best;
  }

  Future<void> _run(Future<void> Function() action) async {
    setState(() => _busy = true);
    try {
      await action();
    } on ApiError catch (e) {
      if (mounted) {
        if (!e.isSessionError) toast(context, Text(errorText(AppLocalizations.of(context), e)));
        await _load(); // state may have moved on (e.g. revealed elsewhere)
      }
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _start() => _run(() async {
    final a = await widget.repo.start(widget.characterId, _kind, _duration, target: _nodeFor(_kind)?.code);
    if (mounted) setState(() => _active = a);
  });

  Future<void> _collect(ActiveActivity a) async {
    final l = AppLocalizations.of(context);
    if (!a.contract.completeAt(DateTime.now()) &&
        !await confirmAction(context, title: l.collectEarlyTitle, body: l.collectEarlyBody, action: l.collect)) {
      return;
    }
    await _run(() async {
      final got = await widget.repo.reveal(a.contract.id);
      if (!mounted) return;
      setState(() => _active = null);
      widget.dataVersion.value++;
      _loadSkills();
      await _showHaul(a, got);
    });
  }

  /// Collects what came in so far and runs the same again (DD-310).
  Future<void> _restart(ActiveActivity a) async {
    final l = AppLocalizations.of(context);
    if (!a.contract.completeAt(DateTime.now()) &&
        !await confirmAction(context, title: l.restartTitle, body: l.restartBody, action: l.restart)) {
      return;
    }
    (ActivityYield, GatherProgress?, PartyShare?)? got;
    ApiError? err;
    await _run(() async {
      final (haul, next, e) = await widget.repo.restart(a.contract.id);
      if (!mounted) return;
      (got, err) = (haul, e);
      setState(() => _active = next);
      widget.dataVersion.value++;
      _loadSkills();
    });
    // Shown once the new run is on screen and idle.
    if (got == null || !mounted) return;
    await _showHaul(a, got!);
    if (err != null && mounted) toast(context, Text(errorText(l, err!)));
  }

  Future<void> _showHaul(ActiveActivity a, (ActivityYield, GatherProgress?, PartyShare?) got) {
    final l = AppLocalizations.of(context);
    final (y, g, share) = got;
    final name = activityName(l, a.contract.kind);
    return showLootSheet(
      context,
      title: name,
      loot: y.toLoot(),
      subtitle: g == null || g.gained == 0
          ? null
          : [
              l.gatherXpGained(name, g.gained),
              if (g.leveledUp) l.gatherLevelUp(name, g.level),
              if (share != null) share.text(l),
            ].join(' · '),
    );
  }

  Future<void> _abandon(ActiveActivity a) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.abandonTitle,
      body: l.abandonBody,
      action: l.abandon,
      danger: true,
    );
    if (!ok) return;
    await _run(() async {
      await widget.repo.stop(a.contract.id);
      if (mounted) setState(() => _active = null);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    if (!_loaded) {
      return _error != null
          ? MessageView(message: errorText(l, _error!), onRetry: _load)
          : const Center(child: CircularProgressIndicator());
    }
    final a = _active;
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
        children: a == null ? _picker(l) : _running(l, a),
      ),
    );
  }

  List<Widget> _picker(AppLocalizations l) {
    final text = Theme.of(context).textTheme;
    return [
      Text(l.activityChoose, style: text.titleMedium),
      const SizedBox(height: Space.md),
      GridView.count(
        crossAxisCount: 3,
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        mainAxisSpacing: Space.sm,
        crossAxisSpacing: Space.sm,
        childAspectRatio: 1.05,
        children: [
          for (final k in activityKinds)
            _KindTile(
              kind: k,
              label: activityName(l, k),
              level: _skills[k] == null ? null : l.levelShort(_skills[k]!.level),
              selected: k == _kind,
              onTap: _busy ? null : () => setState(() => _kind = k),
            ),
        ],
      ),
      if (_skills[_kind] case final g?) ..._resources(l, g),
      const SizedBox(height: Space.xl),
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
              onSelected: _busy ? null : (_) => setState(() => _duration = d),
              selectedColor: Palette.accent.withValues(alpha: 0.2),
              side: BorderSide(color: d == _duration ? Palette.accent : Palette.borderSubtle),
            ),
        ],
      ),
      const SizedBox(height: Space.xl),
      FilledButton.icon(
        onPressed: _busy ? null : _start,
        icon: Icon(activityIcon(_kind)),
        label: Text(l.startActivity(activityName(l, _kind))),
      ),
    ];
  }

  /// Level bar + the resources the level unlocks (best first) and the next
  /// two still locked.
  List<Widget> _resources(AppLocalizations l, GatherSkill g) {
    final text = Theme.of(context).textTheme;
    final chosen = _nodeFor(g.kind);
    final unlocked = g.nodes.where(g.unlocks).toList().reversed;
    final locked = g.nodes.where((n) => !g.unlocks(n)).take(2);
    final max = g.xpToNext == 0;
    return [
      const SizedBox(height: Space.xl),
      Row(
        children: [
          Expanded(child: Text('${activityName(l, g.kind)} · ${l.levelShort(g.level)}', style: text.titleSmall)),
          Text(max ? l.gatherMax : l.gatherXp(g.xp, g.xpToNext), style: text.bodySmall),
        ],
      ),
      const SizedBox(height: Space.xs),
      ClipRRect(
        borderRadius: BorderRadius.circular(Radii.sm),
        child: LinearProgressIndicator(
          value: max ? 1 : (g.xp / g.xpToNext).clamp(0, 1).toDouble(),
          minHeight: 6,
          backgroundColor: Palette.borderSubtle,
          color: Palette.accent,
        ),
      ),
      const SizedBox(height: Space.lg),
      Text(l.gatherResource, style: text.bodySmall),
      const SizedBox(height: Space.sm),
      if (unlocked.isNotEmpty && chosen == null)
        Padding(
          padding: const EdgeInsets.only(bottom: Space.sm),
          child: Text(l.gatherNoneHere, style: text.bodySmall?.copyWith(color: Palette.warning)),
        ),
      for (final n in unlocked)
        if (_here(n))
          GatherNodeTile(
            node: n,
            kind: g.kind,
            trailing: l.levelShort(n.level),
            selected: n.code == chosen?.code,
            onTap: _busy ? null : () => setState(() => _target[g.kind] = n.code),
          )
        else
          GatherNodeTile(
            key: Key('away-${n.code}'),
            node: n,
            kind: g.kind,
            trailing: ContentNames.of('region', n.region),
            locked: true,
          ),
      for (final n in locked) GatherNodeTile(node: n, kind: g.kind, trailing: l.gatherLocked(n.level), locked: true),
    ];
  }

  List<Widget> _running(AppLocalizations l, ActiveActivity a) {
    final text = Theme.of(context).textTheme;
    final y = a.current;
    return [
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
                    backgroundColor: Palette.accent.withValues(alpha: 0.15),
                    child: Icon(activityIcon(a.contract.kind), color: Palette.accent),
                  ),
                  const SizedBox(width: Space.md),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(activityName(l, a.contract.kind), style: text.titleMedium),
                        if (y.target.isNotEmpty) Text(ContentNames.item(y.target), style: text.bodySmall),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: Space.lg),
              ContractProgress(
                startedAt: a.contract.startedAt,
                endsAt: a.contract.expiresAt,
                doneLabel: l.activityDone,
                onComplete: () => setState(() {}),
              ),
              const SizedBox(height: Space.lg),
              Text(l.soFar, style: text.bodySmall),
              const SizedBox(height: Space.xs),
              Wrap(
                spacing: Space.lg,
                runSpacing: Space.xs,
                crossAxisAlignment: WrapCrossAlignment.center,
                children: [
                  GoldText(y.goldCopper),
                  if (y.target.isNotEmpty && y.units > 0)
                    Text(
                      '${ContentNames.item(y.target)} × ${y.units}',
                      style: const TextStyle(color: Palette.accent, fontWeight: FontWeight.w600),
                    ),
                  for (final (r, n) in [('common', y.common), ('rare', y.rare), ('legendary', y.legendary)])
                    if (n > 0 && y.target.isEmpty)
                      Text(
                        '${rarityName(l, r)} × $n',
                        style: TextStyle(color: rarityColor(r), fontWeight: FontWeight.w600),
                      ),
                ],
              ),
            ],
          ),
        ),
      ),
      const SizedBox(height: Space.lg),
      Row(
        children: [
          Expanded(
            child: FilledButton(
              onPressed: _busy ? null : () => _collect(a),
              child: _busy
                  ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
                  : Text(l.collect),
            ),
          ),
          const SizedBox(width: Space.sm),
          Expanded(
            child: OutlinedButton.icon(
              onPressed: _busy ? null : () => _restart(a),
              icon: const Icon(Icons.replay, size: 18),
              label: Text(l.restart),
            ),
          ),
        ],
      ),
      const SizedBox(height: Space.sm),
      TextButton(
        onPressed: _busy ? null : () => _abandon(a),
        style: TextButton.styleFrom(foregroundColor: Palette.danger),
        child: Text(l.abandon),
      ),
    ];
  }
}

class _KindTile extends StatelessWidget {
  const _KindTile({required this.kind, required this.label, required this.selected, this.level, this.onTap});

  final String kind;
  final String label;
  final String? level;
  final bool selected;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
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
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(activityIcon(kind), color: selected ? Palette.accent : Palette.textSecondary, size: 28),
              const SizedBox(height: Space.sm),
              Text(label, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
              if (level != null) Text(level!, style: const TextStyle(fontSize: 11, color: Palette.textSecondary)),
            ],
          ),
        ),
      ),
    );
  }
}

/// One gathering resource: art, rarity-coloured name, level (or lock).
class GatherNodeTile extends StatelessWidget {
  const GatherNodeTile({
    super.key,
    required this.node,
    required this.kind,
    required this.trailing,
    this.selected = false,
    this.locked = false,
    this.onTap,
  });

  final GatherNode node;
  final String kind;
  final String trailing;
  final bool selected;
  final bool locked;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final color = rarityColor(node.rarity);
    return Padding(
      padding: const EdgeInsets.only(bottom: Space.xs),
      child: Opacity(
        opacity: locked ? 0.45 : 1,
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
              padding: const EdgeInsets.symmetric(horizontal: Space.md, vertical: Space.sm),
              child: Row(
                children: [
                  if (kind == 'hunting' && ItemImages.animalPath(huntAnimal(node.code)) != null)
                    // Hunting shows the animal itself.
                    Image.asset(
                      ItemImages.animalPath(huntAnimal(node.code))!,
                      width: 44,
                      height: 44,
                      cacheWidth: (44 * MediaQuery.devicePixelRatioOf(context)).round(),
                      errorBuilder: (_, _, _) =>
                          ItemArt(code: node.code, icon: activityIcon(kind), color: color, size: 26),
                    )
                  else
                    ItemArt(
                      code: node.code,
                      icon: locked ? Icons.lock_outline : activityIcon(kind),
                      color: color,
                      size: 26,
                    ),
                  const SizedBox(width: Space.md),
                  Expanded(
                    child: Text(
                      ContentNames.item(node.code),
                      style: TextStyle(color: color, fontWeight: FontWeight.w600),
                    ),
                  ),
                  Text(trailing, style: const TextStyle(fontSize: 12, color: Palette.textSecondary)),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
