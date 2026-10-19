import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../inventory/item_widgets.dart';
import 'workshop_repository.dart';
import 'workshop_view.dart' show professionName;

/// Server cap on one batch (crafting.MaxCraftCount).
const maxCraftCount = 100;

/// Server cap on one job's length (crafting.MaxJobSeconds, DD-264).
const maxCraftJobSeconds = 6 * 3600;

/// Asks before an early collection, which refunds the unfinished units.
Future<bool> confirmEarlyCollect(BuildContext context, CraftJob job) async {
  final done = job.done(DateTime.now());
  if (done >= job.count) return true;
  final l = AppLocalizations.of(context);
  return confirmAction(
    context,
    title: l.craftCollectEarlyTitle,
    body: l.craftCollectEarlyBody(done, job.count - done),
    action: l.collect,
  );
}

/// Collects [job] and says what came of it; bag and gold changed.
Future<void> collectCraftJob(
  BuildContext context,
  WorkshopRepository repo,
  CraftJob job,
  ValueNotifier<int> dataVersion,
) async {
  final res = await repo.collectCraft(job.id);
  if (!context.mounted) return;
  final l = AppLocalizations.of(context);
  toast(
    context,
    Text(
      [
        if (res.made > 0) l.crafted(itemName(res.output), res.stack),
        if (res.refunded > 0) l.craftCollected(res.made, res.refunded),
        if (res.leveledUp) l.skillUp(professionName(l, job.profession), res.skillAfter),
      ].join(' · '),
    ),
  );
  dataVersion.value++;
}

/// The running crafting job: what, live progress, collect (any time).
class CraftJobCard extends StatelessWidget {
  const CraftJobCard({super.key, required this.job, required this.onCollect, this.onComplete});

  final CraftJob job;
  final VoidCallback? onCollect;
  final VoidCallback? onComplete;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                const Icon(Icons.handyman_outlined, color: Palette.accent),
                const SizedBox(width: Space.sm),
                Expanded(
                  child: Text(
                    l.craftRunning(itemName(job.output), job.count),
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                ),
              ],
            ),
            const SizedBox(height: Space.md),
            ContractProgress(
              startedAt: job.startedAt,
              endsAt: job.readyAt,
              doneLabel: l.craftDone,
              onComplete: onComplete,
              detail: (now) => l.craftProgress(job.done(now), job.count),
            ),
            const SizedBox(height: Space.md),
            FilledButton(onPressed: onCollect, child: Text(l.collect)),
          ],
        ),
      ),
    );
  }
}

/// The recipes worth showing at [skill] (DD-274): everything already
/// learnable, newest first, then the next locked step. The gear professions
/// hold hundreds of recipes; far-off tiers would only bury the useful ones.
List<Recipe> visibleRecipes(List<Recipe> all, int skill) {
  final open = [
    for (final r in all)
      if (r.minSkill <= skill) r,
  ]..sort((a, b) => b.minSkill.compareTo(a.minSkill));
  final locked = [
    for (final r in all)
      if (r.minSkill > skill) r.minSkill,
  ];
  if (locked.isEmpty) return open;
  final next = locked.reduce((a, b) => a < b ? a : b);
  return [
    ...open,
    for (final r in all)
      if (r.minSkill == next) r,
  ];
}

/// One profession's recipes. Tapping an item's name unfolds its properties;
/// Craft asks how many when more than one can be made. Reloads on
/// [dataVersion] (materials and gold change elsewhere).
class ProfessionScreen extends StatefulWidget {
  const ProfessionScreen({super.key, required this.repo, required this.dataVersion, required this.profession});

  final WorkshopRepository repo;
  final ValueNotifier<int> dataVersion;
  final String profession;

  @override
  State<ProfessionScreen> createState() => _ProfessionScreenState();
}

class _ProfessionScreenState extends State<ProfessionScreen> with PanelState {
  List<Recipe> _recipes = const [];
  Skill? _skill;
  Map<String, int> _have = const {};
  int _gold = 0;
  CraftJob? _job;

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(load);
    load();
  }

  @override
  void dispose() {
    widget.dataVersion.removeListener(load);
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object>([
      widget.repo.recipes(widget.profession),
      widget.repo.skills(),
      widget.repo.materials(),
      widget.repo.gold(),
      widget.repo.activeCraft().then<Object>((j) => j ?? false),
    ]);
    _skill = (r[1] as Map<String, Skill>)[widget.profession];
    _recipes = visibleRecipes(r[0] as List<Recipe>, _skill?.level ?? 0);
    _have = r[2] as Map<String, int>;
    _gold = r[3] as int;
    _job = r[4] is CraftJob ? r[4] as CraftJob : null;
  }

  /// How many crafts the bag and purse allow now (capped like the server).
  int _maxCrafts(Recipe rec) {
    final skill = _skill;
    if (skill == null || skill.level < rec.minSkill) return 0;
    var n = maxCraftCount;
    for (final e in rec.inputs.entries) {
      final can = (_have[e.key] ?? 0) ~/ e.value;
      if (can < n) n = can;
    }
    if (rec.goldCost > 0 && _gold ~/ rec.goldCost < n) n = _gold ~/ rec.goldCost;
    final byTime = maxCraftJobSeconds ~/ (rec.craftSeconds <= 0 ? 1 : rec.craftSeconds);
    if (byTime < n) n = byTime;
    return n;
  }

  Future<void> _craft(Recipe rec) async {
    final max = _maxCrafts(rec);
    final count = max > 1
        ? await showDialog<int>(
            context: context,
            builder: (_) => _CraftCountDialog(name: itemName(rec.output), max: max, unitSeconds: rec.craftSeconds),
          )
        : 1;
    if (count == null || count < 1 || !mounted) return;
    await run(() async {
      final job = await widget.repo.craft(rec.code, count: count);
      if (!mounted) return;
      setState(() => _job = job);
      widget.dataVersion.value++; // materials + gold taken; this page reloads via the listener
    });
  }

  Future<void> _collect(CraftJob job) async {
    if (!await confirmEarlyCollect(context, job) || !mounted) return;
    await run(() => collectCraftJob(context, widget.repo, job, widget.dataVersion));
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final skill = _skill;
    return Scaffold(
      appBar: AppBar(title: Text(professionName(l, widget.profession))),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              if (skill != null) ...[
                StatBar(
                  label: l.levelShort(skill.level),
                  value: skill.xpToNext <= 0 ? 1 : skill.xp / skill.xpToNext,
                  color: Palette.epic,
                  valueText: l.xpValue(skill.xp, skill.xpToNext),
                ),
                const SizedBox(height: Space.md),
              ],
              if (_job case final job?)
                CraftJobCard(job: job, onCollect: () => _collect(job), onComplete: () => setState(() {})),
              for (final rec in _recipes)
                _RecipeTile(
                  recipe: rec,
                  have: _have,
                  skill: skill?.level ?? 0,
                  // One job at a time (DD-263).
                  onCraft: _job == null && _maxCrafts(rec) > 0 ? () => _craft(rec) : null,
                ),
            ],
          ),
        ),
      ),
    );
  }
}

/// A recipe: output name (tap to unfold its properties), inputs with
/// have/need, gold cost, skill gate, and the Craft button.
class _RecipeTile extends StatefulWidget {
  const _RecipeTile({required this.recipe, required this.have, required this.skill, required this.onCraft});

  final Recipe recipe;
  final Map<String, int> have;
  final int skill;
  final VoidCallback? onCraft;

  @override
  State<_RecipeTile> createState() => _RecipeTileState();
}

class _RecipeTileState extends State<_RecipeTile> {
  bool _open = false;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final small = Theme.of(context).textTheme.bodySmall;
    final r = widget.recipe;
    final info = r.info;
    final name = r.outputStack > 1 ? '${itemName(r.output)} ×${r.outputStack}' : itemName(r.output);
    return Card(
      margin: const EdgeInsets.only(bottom: Space.sm),
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: InkWell(
                    borderRadius: BorderRadius.circular(Radii.sm),
                    onTap: () => setState(() => _open = !_open),
                    child: Row(
                      children: [
                        SizedBox.square(
                          dimension: 36,
                          child: ItemCell(
                            rarity: 'common',
                            icon: itemIcon(info?.itemClass ?? '', info?.itemSubclass ?? ''),
                            code: r.output,
                          ),
                        ),
                        const SizedBox(width: Space.sm),
                        Expanded(
                          child: Text(name, style: const TextStyle(fontWeight: FontWeight.w600)),
                        ),
                        Icon(_open ? Icons.expand_less : Icons.expand_more, color: Palette.textSecondary),
                      ],
                    ),
                  ),
                ),
                const SizedBox(width: Space.sm),
                FilledButton.tonal(style: compactButton, onPressed: widget.onCraft, child: Text(l.craft)),
              ],
            ),
            const SizedBox(height: Space.xs),
            Wrap(
              spacing: Space.sm,
              children: [
                for (final e in r.inputs.entries)
                  Text(
                    '${itemName(e.key)} ${widget.have[e.key] ?? 0}/${e.value}',
                    style: TextStyle(
                      fontSize: 12,
                      fontFeatures: numeric,
                      color: (widget.have[e.key] ?? 0) >= e.value ? Palette.success : Palette.danger,
                    ),
                  ),
              ],
            ),
            Wrap(
              spacing: Space.sm,
              crossAxisAlignment: WrapCrossAlignment.center,
              children: [
                Text(l.craftEach(shortDuration(l, Duration(seconds: r.craftSeconds))), style: small),
                if (r.goldCost > 0) GoldText(r.goldCost, size: 12),
                if (r.minSkill > widget.skill)
                  Text(l.requiresSkill(r.minSkill), style: small?.copyWith(color: Palette.danger)),
              ],
            ),
            AnimatedSize(
              duration: const Duration(milliseconds: 180),
              alignment: Alignment.topCenter,
              child: _open ? _Details(recipe: r) : const SizedBox(width: double.infinity),
            ),
          ],
        ),
      ),
    );
  }
}

/// The unfolded item properties.
class _Details extends StatelessWidget {
  const _Details({required this.recipe});

  final Recipe recipe;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final info = recipe.info;
    Widget row(String k, String v) => Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        children: [
          Expanded(child: Text(k, style: text.bodySmall)),
          Text(v, style: const TextStyle(fontSize: 13)),
        ],
      ),
    );
    const good = TextStyle(color: Palette.success, fontWeight: FontWeight.w600);
    return Container(
      width: double.infinity,
      margin: const EdgeInsets.only(top: Space.sm),
      padding: const EdgeInsets.all(Space.md),
      decoration: BoxDecoration(color: Palette.bgElevated, borderRadius: BorderRadius.circular(Radii.sm)),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (info != null) ...[
            row(
              l.itemType,
              [
                itemName(info.itemSubclass),
                if (info.equipSlot != null && info.equipSlot!.isNotEmpty) slotName(l, info.equipSlot!),
              ].join(' · '),
            ),
            for (final b in info.bonuses)
              if (bonusText(l, b) case final t?)
                Padding(
                  padding: const EdgeInsets.only(top: 2),
                  child: Text(t, style: good),
                ),
            if (info.healHp > 0) Text(l.restoresHp(info.healHp), style: good),
            if (info.healPct > 0) Text(l.restoresHpPct(info.healPct), style: good),
            if (info.buff case final b?) Text(l.foodBuffLine(b.attackPct, b.defensePct, b.minutes), style: good),
            if (info.minLevel > 1) row(l.requiresLevel, '${info.minLevel}'),
            if (info.classRestriction != null) row(l.characterClass, className(l, info.classRestriction!)),
            row(l.weight, l.weightGrams(info.weightGrams)),
            if (info.maxStack > 1) row(l.stackSize, '${info.maxStack}'),
          ],
          row(l.skillXpReward, '+${recipe.xpReward}'),
        ],
      ),
    );
  }
}

/// Quantity picker: slider plus −/+ and a "max" shortcut.
class _CraftCountDialog extends StatefulWidget {
  const _CraftCountDialog({required this.name, required this.max, required this.unitSeconds});

  final String name;
  final int max;
  final int unitSeconds;

  @override
  State<_CraftCountDialog> createState() => _CraftCountDialogState();
}

class _CraftCountDialogState extends State<_CraftCountDialog> {
  int _n = 1;

  void _set(int v) => setState(() => _n = v.clamp(1, widget.max));

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.craftHowMany(widget.name)),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              IconButton(icon: const Icon(Icons.remove), onPressed: _n > 1 ? () => _set(_n - 1) : null),
              SizedBox(
                width: 56,
                child: Text(
                  '$_n',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(fontFeatures: numeric),
                ),
              ),
              IconButton(icon: const Icon(Icons.add), onPressed: _n < widget.max ? () => _set(_n + 1) : null),
              TextButton(onPressed: () => _set(widget.max), child: Text('max ${widget.max}')),
            ],
          ),
          Slider(
            value: _n.toDouble(),
            min: 1,
            max: widget.max.toDouble(),
            divisions: widget.max - 1,
            onChanged: (v) => _set(v.round()),
          ),
          Text(l.craftTotal(shortDuration(l, Duration(seconds: widget.unitSeconds * _n)))),
          const SizedBox(height: Space.xs),
          Text(l.craftBusyHint, textAlign: TextAlign.center, style: Theme.of(context).textTheme.bodySmall),
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(onPressed: () => Navigator.pop(context, _n), child: Text(l.craftCount(_n))),
      ],
    );
  }
}
