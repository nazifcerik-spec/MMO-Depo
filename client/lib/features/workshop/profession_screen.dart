import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/quantity_field.dart';
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
        if (res.party != null) res.party!.text(l),
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

/// The slot order of recipe headings after the weapons (DD-295).
const _slotOrder = ['offhand', 'helm', 'chest', 'legs', 'gloves', 'boots', 'necklace', 'ring', 'trinket'];
const _enchantOrder = ['weapon', 'armor', 'jewelry', 'profession'];

/// A recipe's heading on its profession page (DD-295) as (rank, key):
/// weapons and shields by kind, armor and jewelry by slot, enchanting
/// scrolls by the gear they go on (profession enchants apart), anything
/// else by its subclass.
(int, String) recipeGroup(Recipe r) {
  final info = r.info;
  if (info == null) return (100, '');
  if (info.enchant case final e?) {
    final profession = e.modifiers.any((m) => m.target.startsWith('gather_yield_') || m.target == 'craft_speed');
    final g = profession ? 'profession' : (e.slots.isEmpty ? 'any' : e.slots.first);
    final i = _enchantOrder.indexOf(g);
    return (20 + (i < 0 ? _enchantOrder.length : i), 'enchant:$g');
  }
  final slot = info.equipSlot ?? '';
  if (info.itemClass == 'weapon') return (0, 'equipment:${info.itemSubclass}');
  if (slot == 'offhand') return (1, 'equipment:${info.itemSubclass}');
  final i = _slotOrder.indexOf(slot);
  if (slot.isNotEmpty) return (2 + (i < 0 ? _slotOrder.length : i), 'slot:$slot');
  return (50, 'sub:${info.itemSubclass}');
}

/// [recipes] under their headings, in heading order; each keeps its order.
List<(String, List<Recipe>)> groupRecipes(List<Recipe> recipes) {
  final groups = <String, List<Recipe>>{};
  final rank = <String, int>{};
  for (final r in recipes) {
    final (n, key) = recipeGroup(r);
    rank[key] = n;
    (groups[key] ??= []).add(r);
  }
  final keys = groups.keys.toList()..sort((a, b) => rank[a] != rank[b] ? rank[a]!.compareTo(rank[b]!) : a.compareTo(b));
  return [for (final k in keys) (k, groups[k]!)];
}

String recipeGroupName(AppLocalizations l, String key) {
  final i = key.indexOf(':');
  final (kind, code) = i < 0 ? ('', key) : (key.substring(0, i), key.substring(i + 1));
  return switch (kind) {
    'equipment' => ContentNames.of('equipment', code),
    'slot' => slotName(l, code),
    'enchant' => switch (code) {
      'weapon' => l.recipeGroupEnchantWeapon,
      'armor' => l.recipeGroupEnchantArmor,
      'jewelry' => l.recipeGroupEnchantJewelry,
      'profession' => l.recipeGroupEnchantProfession,
      _ => l.enchantGroupAny,
    },
    _ => humanize(code),
  };
}

/// Recipe colors (DD-294), like the classic trade-skill ones.
Color difficultyColor(String color) => switch (color) {
  'orange' => const Color(0xFFFF9800),
  'yellow' => const Color(0xFFFFD54F),
  'green' => Palette.success,
  _ => Palette.textMuted,
};

String difficultyLabel(AppLocalizations l, double xpFactor) =>
    xpFactor <= 0 ? l.diffNoXp : l.diffXp((xpFactor * 100).round());

/// What the colors mean: each band's XP share, then grey.
class _DifficultyLegend extends StatelessWidget {
  const _DifficultyLegend({required this.bands});

  final List<DifficultyBand> bands;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final small = Theme.of(context).textTheme.bodySmall;
    Widget dot(String color, String text) => Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(Icons.circle, size: 10, color: difficultyColor(color)),
        const SizedBox(width: 4),
        Text(text, style: small),
      ],
    );
    return Column(
      key: const Key('difficulty-legend'),
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(l.diffLegend, style: small),
        const SizedBox(height: 2),
        Wrap(
          spacing: Space.md,
          runSpacing: 2,
          children: [for (final b in bands) dot(b.color, difficultyLabel(l, b.xpFactor)), dot('grey', l.diffNoXp)],
        ),
      ],
    );
  }
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
  List<(String, List<Recipe>)> _groups = const [];

  /// Folded headings (DD-295).
  final _closed = <String>{};
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
      widget.repo.professions(), // the difficulty bands (DD-294)
    ]);
    _skill = (r[1] as Map<String, Skill>)[widget.profession];
    _groups = groupRecipes(visibleRecipes(r[0] as List<Recipe>, _skill?.level ?? 0));
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
                const SizedBox(height: Space.xs),
                _DifficultyLegend(bands: widget.repo.difficulty),
                const SizedBox(height: Space.md),
              ],
              if (_job case final job?)
                CraftJobCard(job: job, onCollect: () => _collect(job), onComplete: () => setState(() {})),
              for (final (key, recipes) in _groups) ...[
                // One heading alone says nothing; they show from two up.
                if (_groups.length > 1)
                  _GroupHeader(
                    key: Key('group-$key'),
                    name: recipeGroupName(l, key),
                    count: recipes.length,
                    open: !_closed.contains(key),
                    onTap: () => setState(() => _closed.contains(key) ? _closed.remove(key) : _closed.add(key)),
                  ),
                if (_groups.length == 1 || !_closed.contains(key))
                  for (final rec in recipes)
                    _RecipeTile(
                      key: ValueKey(rec.code),
                      recipe: rec,
                      have: _have,
                      skill: skill?.level ?? 0,
                      difficulty: recipeDifficulty(widget.repo.difficulty, rec.minSkill, skill?.level ?? 0),
                      // One job at a time (DD-263).
                      onCraft: _job == null && _maxCrafts(rec) > 0 ? () => _craft(rec) : null,
                    ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

/// A recipe heading: name, recipe count, fold arrow (DD-295).
class _GroupHeader extends StatelessWidget {
  const _GroupHeader({super.key, required this.name, required this.count, required this.open, required this.onTap});

  final String name;
  final int count;
  final bool open;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final text = Theme.of(context).textTheme;
    return InkWell(
      borderRadius: BorderRadius.circular(Radii.sm),
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(Space.xs, Space.md, Space.xs, Space.sm),
        child: Row(
          children: [
            Expanded(
              child: Text(name, style: text.titleSmall?.copyWith(color: Palette.accent)),
            ),
            Text('$count', style: text.bodySmall),
            Icon(open ? Icons.expand_less : Icons.expand_more, color: Palette.textSecondary),
          ],
        ),
      ),
    );
  }
}

/// A recipe: output name (tap to unfold its properties), inputs with
/// have/need, gold cost, skill gate, and the Craft button.
class _RecipeTile extends StatefulWidget {
  const _RecipeTile({
    super.key,
    required this.recipe,
    required this.have,
    required this.skill,
    required this.onCraft,
    this.difficulty = ('orange', 1),
  });

  final Recipe recipe;
  final Map<String, int> have;
  final int skill;

  /// Color and XP share at the character's skill (DD-294).
  final (String, double) difficulty;
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
                          child: Text(
                            name,
                            style: TextStyle(
                              fontWeight: FontWeight.w600,
                              color: r.minSkill > widget.skill ? null : difficultyColor(widget.difficulty.$1),
                            ),
                          ),
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
            // What a scroll adds, without unfolding (DD-295).
            if (info?.enchant case final e?)
              Padding(
                padding: const EdgeInsets.only(top: Space.xs),
                child: Text(
                  enchantEffects(l, e),
                  style: small?.copyWith(color: Palette.epic, fontWeight: FontWeight.w600),
                ),
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
                  Text(l.requiresSkill(r.minSkill), style: small?.copyWith(color: Palette.danger))
                else
                  Text(
                    difficultyLabel(l, widget.difficulty.$2),
                    key: Key('diff-${r.code}'),
                    style: small?.copyWith(color: difficultyColor(widget.difficulty.$1), fontWeight: FontWeight.w600),
                  ),
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
            if (info.manaPct > 0) Text(l.restoresMpPct(info.manaPct), style: good),
            if (info.buff case final b?) Text(l.foodBuffLine(b.attackPct, b.defensePct, b.minutes), style: good),
            if (info.enchant case final e?) ...[
              Text(l.enchantFits(enchantGroups(l, e))),
              for (final b in e.modifiers)
                if (bonusText(l, b) case final t?) Text(t, style: good),
            ],
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

/// Quantity picker: a number box, a slider and a "max" shortcut.
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
              QuantityField(
                key: const Key('craft-count'),
                value: _n,
                min: 1,
                max: widget.max,
                width: 80,
                onChanged: _set,
              ),
              const SizedBox(width: Space.sm),
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
