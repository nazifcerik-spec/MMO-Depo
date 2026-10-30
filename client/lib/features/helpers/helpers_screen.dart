import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../activity/activity_models.dart';
import '../activity/activity_repository.dart';
import '../activity/activity_view.dart' show GatherNodeTile, activityIcon;
import '../inventory/item_images.dart';
import '../workshop/profession_screen.dart';
import '../workshop/workshop_repository.dart';
import '../workshop/workshop_view.dart' show professionName;

/// One hired-helper type and the character's state of it (DD-319).
class HelperInfo {
  const HelperInfo({
    required this.code,
    required this.kind,
    required this.profession,
    required this.bonus,
    required this.bonusPct,
    required this.price,
    required this.hired,
    this.job,
  });

  factory HelperInfo.fromJson(Map<String, dynamic> j) {
    final job = j['job'] as Map<String, dynamic>?;
    return HelperInfo(
      code: j['code'] as String,
      kind: j['kind'] as String,
      profession: j['profession'] as String,
      bonus: j['bonus'] as String,
      bonusPct: (j['bonus_pct'] as num).toInt(),
      price: (j['price_copper'] as num).toInt(),
      hired: j['hired'] as bool? ?? false,
      job: job == null ? null : HelperJob.fromJson(job),
    );
  }

  final String code, kind, profession, bonus;
  final int bonusPct, price;
  final bool hired;
  final HelperJob? job;

  bool get gathers => kind == 'gather';

  /// helper.Type.pct: the bonus applied to [v] when the helper has it.
  double _pct(String b, double v) {
    if (bonus != b) return v;
    return b == 'speed' || b == 'gold_save' ? v * (100 - bonusPct) / 100 : v * (100 + bonusPct) / 100;
  }

  /// One unit's time: twice the character's, less a speed bonus.
  int unitSeconds(int craftSeconds) => _pct('speed', craftSeconds * 2.0).ceil();

  /// What [count] crafts cost in gold (helper.Type.GoldCost).
  int goldCost(int perUnit, int count) => _pct('gold_save', (perUnit * count).toDouble()).round();

  /// How many of [r] the bag and purse pay for now, capped like the
  /// server (crafting.MaxCraftCount).
  int canMake(Recipe r, Map<String, int> have, int gold) {
    var n = maxCraftCount;
    for (final e in r.inputs.entries) {
      final can = (have[e.key] ?? 0) ~/ e.value;
      if (can < n) n = can;
    }
    while (n > 0 && goldCost(r.goldCost, n) > gold) {
      n--;
    }
    return n;
  }
}

class HelperJob {
  const HelperJob({required this.startedAt, required this.readyAt, this.target = '', this.recipe = '', this.count = 0});

  factory HelperJob.fromJson(Map<String, dynamic> j) {
    final spec = j['spec'] as Map<String, dynamic>? ?? const {};
    return HelperJob(
      startedAt: DateTime.parse(j['started_at'] as String),
      readyAt: DateTime.parse(j['ready_at'] as String),
      target: spec['target'] as String? ?? '',
      recipe: spec['recipe'] as String? ?? '',
      count: (spec['count'] as num?)?.toInt() ?? 0,
    );
  }

  final DateTime startedAt, readyAt;
  final String target, recipe;
  final int count;
}

class HelperRepository {
  HelperRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId/helpers';

  List<HelperInfo> _parse(Map<String, dynamic> j) => [
    for (final h in (j['helpers'] as List? ?? const []).cast<Map<String, dynamic>>()) HelperInfo.fromJson(h),
  ];

  Future<List<HelperInfo>> list() async => _parse(await api.get(_base));
  Future<List<HelperInfo>> hire(String code) async => _parse(await api.post('$_base/$code/hire', const {}));
  Future<List<HelperInfo>> startGather(String code, String target, Duration d) async =>
      _parse(await api.post('$_base/$code/start', {'target': target, 'duration_seconds': d.inSeconds}));
  Future<List<HelperInfo>> startCraft(String code, String recipe, int count) async =>
      _parse(await api.post('$_base/$code/start', {'recipe_code': recipe, 'count': count}));
  Future<List<HelperInfo>> cancel(String code) async => _parse(await api.post('$_base/$code/cancel', const {}));

  /// Collects; returns the XP and the items it brought.
  Future<(int xp, int items)> collect(String code) async {
    final r = (await api.post('$_base/$code/collect', const {}))['result'] as Map<String, dynamic>? ?? const {};
    final items = (r['items'] as List? ?? const []).length + ((r['made'] as num?)?.toInt() ?? 0);
    return ((r['xp'] as num?)?.toInt() ?? 0, items);
  }
}

String helperBonus(AppLocalizations l, HelperInfo h) => switch (h.bonus) {
  'yield' => l.helperBonusYield(h.bonusPct),
  'speed' => l.helperBonusSpeed(h.bonusPct),
  'xp' => l.helperBonusXp(h.bonusPct),
  'gold_save' => l.helperBonusGold(h.bonusPct),
  'double' => l.helperBonusDouble(h.bonusPct),
  _ => '',
};

/// Hired helpers (DD-319): buy one per profession, send it on gathering
/// runs or crafting batches beside the character; half speed and XP.
class HelpersScreen extends StatefulWidget {
  const HelpersScreen({super.key, required this.repo, required this.workshop, required this.dataVersion});

  final HelperRepository repo;
  final WorkshopRepository workshop;
  final ValueNotifier<int> dataVersion;

  @override
  State<HelpersScreen> createState() => _HelpersScreenState();
}

class _HelpersScreenState extends State<HelpersScreen> with PanelState {
  List<HelperInfo> _helpers = const [];

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _helpers = await widget.repo.list();

  Future<void> _set(Future<List<HelperInfo>> Function() action) => run(() async {
    final list = await action();
    widget.dataVersion.value++;
    if (mounted) setState(() => _helpers = list);
  });

  Future<void> _hire(HelperInfo h) async {
    final l = AppLocalizations.of(context);
    if (!await confirmAction(
      context,
      title: l.helperHire,
      body: l.helperHireBody(helperName(l, h.code), goldString(h.price)),
      action: l.helperHire,
    )) {
      return;
    }
    await _set(() => widget.repo.hire(h.code));
  }

  Future<void> _start(HelperInfo h) async {
    if (h.gathers) {
      final pick = await Navigator.of(context).push<(String, Duration)>(
        MaterialPageRoute(
          builder: (_) =>
              _GatherPickScreen(helper: h, activity: ActivityRepository(widget.repo.api), repo: widget.repo),
        ),
      );
      if (pick != null) await _set(() => widget.repo.startGather(h.code, pick.$1, pick.$2));
      return;
    }
    final pick = await Navigator.of(context).push<(Recipe, int)>(
      MaterialPageRoute(
        builder: (_) => _HelperCraftScreen(helper: h, workshop: widget.workshop),
      ),
    );
    if (pick != null) await _set(() => widget.repo.startCraft(h.code, pick.$1.code, pick.$2));
  }

  Future<void> _collect(HelperInfo h) => run(() async {
    final l = AppLocalizations.of(context);
    final (xp, items) = await widget.repo.collect(h.code);
    final list = await widget.repo.list();
    widget.dataVersion.value++;
    if (!mounted) return;
    setState(() => _helpers = list);
    toast(context, Text(l.helperCollected(items, xp, professionLabel(l, h.profession))));
  });

  Future<void> _cancel(HelperInfo h) async {
    final l = AppLocalizations.of(context);
    if (!await confirmAction(context, title: l.helperCancel, body: l.helperCancelBody, action: l.helperCancel)) return;
    await _set(() => widget.repo.cancel(h.code));
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text(l.helpersTitle)),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              Text(l.helpersHint, style: text.bodySmall),
              for (final group in [true, false]) ...[
                const SizedBox(height: Space.lg),
                Text(group ? l.helpersGather : l.helpersCraft, style: text.titleMedium),
                for (final h in _helpers)
                  if (h.gathers == group) _card(l, text, h),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _card(AppLocalizations l, TextTheme text, HelperInfo h) {
    final job = h.job;
    final done = job != null && !DateTime.now().isBefore(job.readyAt);
    return Card(
      key: Key('helper-${h.code}'),
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('${helperName(l, h.code)} · ${professionLabel(l, h.profession)}', style: text.titleSmall),
            Text(helperBonus(l, h), style: text.bodySmall?.copyWith(color: Palette.success)),
            if (job != null) ...[
              const SizedBox(height: Space.xs),
              Row(
                children: [
                  ItemArt(
                    code: job.recipe.isNotEmpty ? job.recipe.replaceFirst('recipe_', '') : job.target,
                    icon: h.gathers ? activityIcon(h.profession) : Icons.handyman_outlined,
                    color: Palette.textSecondary,
                  ),
                  const SizedBox(width: Space.sm),
                  Expanded(
                    child: Text(
                      job.recipe.isNotEmpty
                          ? '${ContentNames.item(job.recipe.replaceFirst('recipe_', ''))} ×${job.count}'
                          : ContentNames.item(job.target),
                      style: text.bodySmall,
                    ),
                  ),
                ],
              ),
              Text(
                done ? l.helperDone : l.helperWorking(shortDuration(l, job.readyAt.difference(DateTime.now()))),
                style: text.bodySmall,
              ),
            ],
            const SizedBox(height: Space.sm),
            Wrap(
              spacing: Space.sm,
              children: [
                if (!h.hired)
                  FilledButton(onPressed: () => _hire(h), child: Text('${l.helperHire} · ${goldString(h.price)}'))
                else if (job == null)
                  FilledButton(onPressed: () => _start(h), child: Text(l.helperSend))
                else ...[
                  FilledButton(onPressed: done || h.gathers ? () => _collect(h) : null, child: Text(l.collect)),
                  TextButton(onPressed: () => _cancel(h), child: Text(l.helperCancel)),
                ],
              ],
            ),
          ],
        ),
      ),
    );
  }
}

/// A profession's display name: gathering kinds and crafting professions.
String professionLabel(AppLocalizations l, String code) => switch (code) {
  'mining' || 'woodcutting' || 'fishing' || 'herbalism' || 'hunting' => activityName(l, code),
  _ => professionName(l, code),
};

String helperName(AppLocalizations l, String code) => switch (code) {
  'miner' => l.helperMiner,
  'woodsman' => l.helperWoodsman,
  'angler' => l.helperAngler,
  'herbalist' => l.helperHerbalist,
  'trapper' => l.helperTrapper,
  'smith' => l.helperSmith,
  'tanner' => l.helperTanner,
  'weaver' => l.helperWeaver,
  'carver' => l.helperCarver,
  'apothecary' => l.helperApothecary,
  'cook' => l.helperCook,
  'smelter' => l.helperSmelter,
  'scribe' => l.helperScribe,
  'gemcutter' => l.helperGemcutter,
  _ => code,
};

/// A gathering helper's run (DD-327): the resources the character's level
/// unlocks, best first, with their pictures, and the length.
class _GatherPickScreen extends StatefulWidget {
  const _GatherPickScreen({required this.helper, required this.activity, required this.repo});

  final HelperInfo helper;
  final ActivityRepository activity;
  final HelperRepository repo;

  @override
  State<_GatherPickScreen> createState() => _GatherPickScreenState();
}

class _GatherPickScreenState extends State<_GatherPickScreen> with PanelState {
  List<GatherNode> _nodes = const [];
  String? _target;
  Duration _duration = const Duration(hours: 1);

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final skills = await widget.activity.skills(widget.repo.characterId);
    final g = skills.where((s) => s.kind == widget.helper.profession).firstOrNull;
    _nodes = g == null ? const [] : g.nodes.where(g.unlocks).toList().reversed.toList();
    _target ??= _nodes.firstOrNull?.code;
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final kind = widget.helper.profession;
    return Scaffold(
      appBar: AppBar(title: Text('${helperName(l, widget.helper.code)} · ${professionLabel(l, kind)}')),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              Text(l.durationLabel, style: text.titleSmall),
              const SizedBox(height: Space.xs),
              Wrap(
                spacing: Space.sm,
                children: [
                  for (final d in const [Duration(hours: 1), Duration(hours: 3), Duration(hours: 6)])
                    ChoiceChip(
                      label: Text(shortDuration(l, d)),
                      selected: d == _duration,
                      onSelected: (_) => setState(() => _duration = d),
                    ),
                ],
              ),
              const SizedBox(height: Space.md),
              FilledButton(
                key: const Key('helper-gather-send'),
                onPressed: _target == null ? null : () => Navigator.of(context).pop((_target!, _duration)),
                child: Text(l.helperSend),
              ),
              const SizedBox(height: Space.lg),
              Text(l.helperPickResource, style: text.titleSmall),
              const SizedBox(height: Space.xs),
              for (final n in _nodes)
                GatherNodeTile(
                  key: Key('helper-node-${n.code}'),
                  node: n,
                  kind: kind,
                  trailing: l.levelShort(n.level),
                  selected: n.code == _target,
                  onTap: () => setState(() => _target = n.code),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

/// A crafting helper's batch (DD-327): the recipes the character knows,
/// grouped like the profession page, each with its picture, properties
/// (tap the name) and how many the bag and purse pay for at the helper's
/// time and gold; the count picker's max is that number.
class _HelperCraftScreen extends StatefulWidget {
  const _HelperCraftScreen({required this.helper, required this.workshop});

  final HelperInfo helper;
  final WorkshopRepository workshop;

  @override
  State<_HelperCraftScreen> createState() => _HelperCraftScreenState();
}

class _HelperCraftScreenState extends State<_HelperCraftScreen> with PanelState {
  List<(String, List<Recipe>)> _groups = const [];
  final _closed = <String>{};
  Map<String, int> _have = const {};
  int _gold = 0;
  int _skill = 0;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final profession = widget.helper.profession;
    final r = await Future.wait<Object>([
      widget.workshop.recipes(profession),
      widget.workshop.skills(),
      widget.workshop.materials(),
      widget.workshop.gold(),
      widget.workshop.professions(), // the difficulty bands (DD-294)
    ]);
    _skill = (r[1] as Map<String, Skill>)[profession]?.level ?? 0;
    final open = [
      for (final rec in r[0] as List<Recipe>)
        if (rec.minSkill <= _skill) rec,
    ]..sort((a, b) => b.minSkill.compareTo(a.minSkill));
    _groups = groupRecipes(open);
    _have = r[2] as Map<String, int>;
    _gold = r[3] as int;
  }

  Future<void> _pick(Recipe rec, int max) async {
    final l = AppLocalizations.of(context);
    final count = await showDialog<int>(
      context: context,
      builder: (_) => CraftCountDialog(
        name: itemName(rec.output),
        max: max,
        unitSeconds: widget.helper.unitSeconds(rec.craftSeconds),
        hint: l.helperCraftHint,
      ),
    );
    if (count != null && count > 0 && mounted) Navigator.of(context).pop((rec, count));
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final h = widget.helper;
    return Scaffold(
      appBar: AppBar(title: Text('${helperName(l, h.code)} · ${professionLabel(l, h.profession)}')),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              Text(helperBonus(l, h), style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Palette.success)),
              const SizedBox(height: Space.sm),
              if (_groups.isEmpty) Text(l.helperNoRecipes),
              for (final (key, recipes) in _groups) ...[
                if (_groups.length > 1)
                  GroupHeader(
                    name: recipeGroupName(l, key),
                    count: recipes.length,
                    open: !_closed.contains(key),
                    onTap: () => setState(() => _closed.contains(key) ? _closed.remove(key) : _closed.add(key)),
                  ),
                if (_groups.length == 1 || !_closed.contains(key))
                  for (final rec in recipes) _tile(l, rec),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _tile(AppLocalizations l, Recipe rec) {
    final h = widget.helper;
    final n = h.canMake(rec, _have, _gold);
    return RecipeTile(
      key: ValueKey(rec.code),
      recipe: rec,
      have: _have,
      skill: _skill,
      difficulty: recipeDifficulty(widget.workshop.difficulty, rec.minSkill, _skill),
      unitSeconds: h.unitSeconds(rec.craftSeconds),
      goldCost: h.goldCost(rec.goldCost, 1),
      canMake: n,
      action: l.helperSend,
      onCraft: n > 0 ? () => _pick(rec, n) : null,
    );
  }
}
