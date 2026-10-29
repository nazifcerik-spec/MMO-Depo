import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/panel.dart';
import '../../shared/quantity_field.dart';
import '../../shared/widgets.dart';
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
  Future<List<HelperInfo>> startGather(String code, Duration d) async =>
      _parse(await api.post('$_base/$code/start', {'duration_seconds': d.inSeconds}));
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
      final d = await showDialog<Duration>(
        context: context,
        builder: (c) => SimpleDialog(
          title: Text(AppLocalizations.of(context).helperSend),
          children: [
            for (final d in const [Duration(hours: 1), Duration(hours: 3), Duration(hours: 6)])
              SimpleDialogOption(
                onPressed: () => Navigator.of(c).pop(d),
                child: Text(shortDuration(AppLocalizations.of(context), d)),
              ),
          ],
        ),
      );
      if (d != null) await _set(() => widget.repo.startGather(h.code, d));
      return;
    }
    final skills = await widget.workshop.skills();
    final level = skills[h.profession]?.level ?? 0;
    final recipes = [
      for (final r in await widget.workshop.recipes(h.profession))
        if (r.minSkill <= level) r,
    ];
    if (!mounted) return;
    final pick = await showDialog<(Recipe, int)>(
      context: context,
      builder: (_) => _CraftDialog(recipes: recipes),
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
              Text(
                job.recipe.isNotEmpty
                    ? '${ContentNames.item(job.recipe.replaceFirst('recipe_', ''))} ×${job.count}'
                    : ContentNames.item(job.target),
                style: text.bodySmall,
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

class _CraftDialog extends StatefulWidget {
  const _CraftDialog({required this.recipes});

  final List<Recipe> recipes;

  @override
  State<_CraftDialog> createState() => _CraftDialogState();
}

class _CraftDialogState extends State<_CraftDialog> {
  Recipe? _recipe;
  int _count = 1;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.helperSend),
      content: widget.recipes.isEmpty
          ? Text(l.helperNoRecipes)
          : Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                DropdownButton<Recipe>(
                  key: const Key('helper-recipe'),
                  isExpanded: true,
                  value: _recipe,
                  hint: Text(l.helperPickRecipe),
                  items: [
                    for (final r in widget.recipes)
                      DropdownMenuItem(
                        value: r,
                        child: Text(ContentNames.item(r.output), overflow: TextOverflow.ellipsis),
                      ),
                  ],
                  onChanged: (r) => setState(() => _recipe = r),
                ),
                const SizedBox(height: Space.sm),
                QuantityField(value: _count, min: 1, max: 100, onChanged: (v) => setState(() => _count = v)),
              ],
            ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.cancel)),
        FilledButton(
          onPressed: _recipe == null ? null : () => Navigator.of(context).pop((_recipe!, _count)),
          child: Text(l.helperSend),
        ),
      ],
    );
  }
}
