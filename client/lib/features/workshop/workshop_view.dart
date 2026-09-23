import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../market/market_repository.dart';
import '../market/market_view.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_widgets.dart';
import 'workshop_repository.dart';

String professionName(AppLocalizations l, String code) => switch (code) {
  'blacksmithing' => l.profBlacksmithing,
  'alchemy' => l.profAlchemy,
  'cooking' => l.profCooking,
  _ => humanize(code),
};

/// Workshop tab: crafting and the general vendor behind one segmented control.
/// Both panels reload on [dataVersion] (materials and gold change elsewhere).
class WorkshopView extends StatefulWidget {
  const WorkshopView({super.key, required this.repo, required this.dataVersion});

  final WorkshopRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<WorkshopView> createState() => _WorkshopViewState();
}

class _WorkshopViewState extends State<WorkshopView> {
  int _seg = 0;
  final _visited = <int>{0};

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final panels = <Widget Function()>[
      () => _CraftPanel(repo: widget.repo, dataVersion: widget.dataVersion),
      () => _ShopPanel(repo: widget.repo, dataVersion: widget.dataVersion),
      () => MarketPanel(
        repo: MarketRepository(widget.repo.api, widget.repo.characterId),
        dataVersion: widget.dataVersion,
        gold: widget.repo.gold,
      ),
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
                ButtonSegment(value: 0, label: Text(l.segCraft), icon: const Icon(Icons.handyman_outlined, size: 18)),
                ButtonSegment(value: 1, label: Text(l.segShop), icon: const Icon(Icons.storefront_outlined, size: 18)),
                ButtonSegment(value: 2, label: Text(l.segMarket), icon: const Icon(Icons.balance_outlined, size: 18)),
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

// --- Crafting -----------------------------------------------------------------

class _CraftPanel extends StatefulWidget {
  const _CraftPanel({required this.repo, required this.dataVersion});

  final WorkshopRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<_CraftPanel> createState() => _CraftPanelState();
}

class _CraftPanelState extends State<_CraftPanel> with PanelState {
  List<String> _professions = const [];
  List<Recipe> _recipes = const [];
  Map<String, Skill> _skills = const {};
  Map<String, int> _have = const {};
  int _gold = 0;

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
      widget.repo.professions(),
      widget.repo.recipes(),
      widget.repo.skills(),
      widget.repo.materials(),
      widget.repo.gold(),
    ]);
    _professions = r[0] as List<String>;
    _recipes = r[1] as List<Recipe>;
    _skills = r[2] as Map<String, Skill>;
    _have = r[3] as Map<String, int>;
    _gold = r[4] as int;
  }

  bool _canCraft(Recipe rec) {
    final skill = _skills[rec.profession];
    return skill != null &&
        skill.level >= rec.minSkill &&
        rec.goldCost <= _gold &&
        rec.inputs.entries.every((e) => (_have[e.key] ?? 0) >= e.value);
  }

  Future<void> _learn(String p) => run(() async {
    await widget.repo.learn(p);
    await fetch();
  });

  Future<void> _craft(Recipe rec) => run(() async {
    final l = AppLocalizations.of(context);
    final res = await widget.repo.craft(rec.code);
    if (!mounted) return;
    toast(
      context,
      Text(
        [
          l.crafted(itemName(res.output), res.stack),
          if (res.leveledUp) l.skillUp(professionName(l, rec.profession), res.skillAfter),
        ].join(' · '),
      ),
    );
    widget.dataVersion.value++; // bag + gold changed; this panel reloads via the listener
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(
      () => [
        for (final p in _professions) ...[
          Card(
            child: Padding(
              padding: const EdgeInsets.all(Space.lg),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Row(
                    children: [
                      Expanded(child: Text(professionName(l, p), style: text.titleMedium)),
                      if (_skills[p] == null)
                        FilledButton.tonal(style: compactButton, onPressed: () => _learn(p), child: Text(l.learn)),
                    ],
                  ),
                  if (_skills[p] case final s?) ...[
                    const SizedBox(height: Space.sm),
                    StatBar(
                      label: l.levelShort(s.level),
                      value: s.xpToNext <= 0 ? 1 : s.xp / s.xpToNext,
                      color: Palette.epic,
                      valueText: l.xpValue(s.xp, s.xpToNext),
                    ),
                    for (final rec in _recipes.where((r) => r.profession == p)) ...[
                      const Divider(height: Space.xl, color: Palette.borderSubtle),
                      _RecipeRow(
                        recipe: rec,
                        have: _have,
                        skill: s.level,
                        canCraft: _canCraft(rec),
                        onCraft: () => _craft(rec),
                      ),
                    ],
                  ] else ...[
                    const SizedBox(height: Space.xs),
                    Text(l.learnHint, style: text.bodySmall),
                  ],
                ],
              ),
            ),
          ),
          const SizedBox(height: Space.md),
        ],
      ],
    );
  }
}

class _RecipeRow extends StatelessWidget {
  const _RecipeRow({
    required this.recipe,
    required this.have,
    required this.skill,
    required this.canCraft,
    required this.onCraft,
  });

  final Recipe recipe;
  final Map<String, int> have;
  final int skill;
  final bool canCraft;
  final VoidCallback onCraft;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final small = Theme.of(context).textTheme.bodySmall;
    final r = recipe;
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                r.outputStack > 1 ? '${itemName(r.output)} ×${r.outputStack}' : itemName(r.output),
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 2),
              Wrap(
                spacing: Space.sm,
                children: [
                  for (final e in r.inputs.entries)
                    Text(
                      '${itemName(e.key)} ${have[e.key] ?? 0}/${e.value}',
                      style: TextStyle(
                        fontSize: 12,
                        fontFeatures: numeric,
                        color: (have[e.key] ?? 0) >= e.value ? Palette.success : Palette.danger,
                      ),
                    ),
                ],
              ),
              if (r.goldCost > 0 || r.minSkill > skill)
                Wrap(
                  spacing: Space.sm,
                  crossAxisAlignment: WrapCrossAlignment.center,
                  children: [
                    if (r.goldCost > 0) GoldText(r.goldCost, size: 12),
                    if (r.minSkill > skill)
                      Text(l.requiresSkill(r.minSkill), style: small?.copyWith(color: Palette.danger)),
                  ],
                ),
            ],
          ),
        ),
        const SizedBox(width: Space.sm),
        FilledButton.tonal(style: compactButton, onPressed: canCraft ? onCraft : null, child: Text(l.craft)),
      ],
    );
  }
}

// --- Shop ---------------------------------------------------------------------

class _ShopPanel extends StatefulWidget {
  const _ShopPanel({required this.repo, required this.dataVersion});

  final WorkshopRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<_ShopPanel> createState() => _ShopPanelState();
}

class _ShopPanelState extends State<_ShopPanel> with PanelState {
  List<Offer> _stock = const [];
  List<Buyback> _buybacks = const [];
  int _gold = 0;

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
    final r = await Future.wait<Object>([widget.repo.stock(), widget.repo.buybacks(), widget.repo.gold()]);
    _stock = r[0] as List<Offer>;
    _buybacks = r[1] as List<Buyback>;
    _gold = r[2] as int;
  }

  Future<void> _buy(Offer o) => run(() async {
    final l = AppLocalizations.of(context);
    await widget.repo.buy(o);
    if (!mounted) return;
    toast(context, Text(l.bought(itemName(o.templateCode))));
    widget.dataVersion.value++;
  });

  Future<void> _buyBack(Buyback b) => run(() async {
    final l = AppLocalizations.of(context);
    await widget.repo.buyBack(b);
    if (!mounted) return;
    toast(context, Text(l.bought(itemName(b.templateCode, b.prefix, b.suffix))));
    widget.dataVersion.value++;
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final now = DateTime.now();
    return frame(
      () => [
        Row(
          children: [
            Expanded(child: Text(humanize(generalVendor), style: text.titleMedium)),
            Semantics(label: l.gold, child: GoldText(_gold)),
          ],
        ),
        const SizedBox(height: Space.md),
        for (final o in _stock)
          Card(
            margin: const EdgeInsets.only(bottom: Space.sm),
            child: ListTile(
              leading: Icon(itemIcon(o.itemClass, o.itemSubclass), color: Palette.textSecondary),
              title: Text(itemName(o.templateCode)),
              subtitle: GoldText(o.price, size: 13),
              trailing: FilledButton.tonal(
                style: compactButton,
                onPressed: o.price <= _gold ? () => _buy(o) : null,
                child: Text(l.buy),
              ),
            ),
          ),
        if (_buybacks.isNotEmpty) ...[
          const SizedBox(height: Space.lg),
          Text(l.buybacks, style: text.titleMedium),
          const SizedBox(height: Space.xs),
          Text(l.buybackHint, style: text.bodySmall),
          const SizedBox(height: Space.md),
          for (final b in _buybacks)
            Card(
              margin: const EdgeInsets.only(bottom: Space.sm),
              child: ListTile(
                title: Text(
                  b.stack > 1
                      ? '${itemName(b.templateCode, b.prefix, b.suffix)} ×${b.stack}'
                      : itemName(b.templateCode, b.prefix, b.suffix),
                  style: TextStyle(color: rarityColor(b.rarity)),
                ),
                subtitle: Text(l.remaining(shortDuration(l, b.expiresAt.difference(now)))),
                trailing: FilledButton.tonal(
                  style: compactButton,
                  onPressed: b.price <= _gold ? () => _buyBack(b) : null,
                  child: GoldText(b.price, size: 13),
                ),
              ),
            ),
        ],
      ],
    );
  }
}
