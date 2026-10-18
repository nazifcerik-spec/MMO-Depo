import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../market/market_repository.dart';
import '../market/market_view.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_images.dart';
import '../inventory/item_widgets.dart';
import 'profession_screen.dart';
import 'workshop_repository.dart';
import '../../shared/content_names.dart';

String professionName(AppLocalizations l, String code) => switch (code) {
  'blacksmithing' => l.profBlacksmithing,
  'alchemy' => l.profAlchemy,
  'cooking' => l.profCooking,
  _ => ContentNames.of('profession', code),
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

/// Profession list: level, recipe count, learn. Tapping a learned profession
/// opens its recipes on their own page ([ProfessionScreen]).
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
      widget.repo.professions(),
      widget.repo.recipes(),
      widget.repo.skills(),
      widget.repo.activeCraft().then<Object>((j) => j ?? false),
    ]);
    _professions = r[0] as List<String>;
    _recipes = r[1] as List<Recipe>;
    _skills = r[2] as Map<String, Skill>;
    _job = r[3] is CraftJob ? r[3] as CraftJob : null;
  }

  Future<void> _collect(CraftJob job) async {
    if (!await confirmEarlyCollect(context, job) || !mounted) return;
    await run(() => collectCraftJob(context, widget.repo, job, widget.dataVersion));
  }

  Future<void> _learn(String p) => run(() async {
    await widget.repo.learn(p);
    await fetch();
  });

  void _open(String p) => Navigator.of(context).push(
    MaterialPageRoute<void>(
      builder: (_) => ProfessionScreen(repo: widget.repo, dataVersion: widget.dataVersion, profession: p),
    ),
  );

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(
      () => [
        if (_job case final job?)
          CraftJobCard(job: job, onCollect: () => _collect(job), onComplete: () => setState(() {})),
        for (final p in _professions)
          Card(
            margin: const EdgeInsets.only(bottom: Space.md),
            child: InkWell(
              borderRadius: BorderRadius.circular(Radii.md),
              onTap: _skills[p] == null ? null : () => _open(p),
              child: Padding(
                padding: const EdgeInsets.all(Space.lg),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Row(
                      children: [
                        Expanded(child: Text(professionName(l, p), style: text.titleMedium)),
                        if (_skills[p] == null)
                          FilledButton.tonal(style: compactButton, onPressed: () => _learn(p), child: Text(l.learn))
                        else ...[
                          Text(l.recipeCount(_recipes.where((r) => r.profession == p).length), style: text.bodySmall),
                          const Icon(Icons.chevron_right, color: Palette.textSecondary),
                        ],
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
                    ] else ...[
                      const SizedBox(height: Space.xs),
                      Text(l.learnHint, style: text.bodySmall),
                    ],
                  ],
                ),
              ),
            ),
          ),
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
            Expanded(child: Text(ContentNames.of('vendor', generalVendor), style: text.titleMedium)),
            Semantics(label: l.gold, child: GoldText(_gold)),
          ],
        ),
        const SizedBox(height: Space.md),
        for (final o in _stock)
          Card(
            margin: const EdgeInsets.only(bottom: Space.sm),
            child: ListTile(
              leading: ItemArt(
                code: o.templateCode,
                icon: itemIcon(o.itemClass, o.itemSubclass),
                color: Palette.textSecondary,
              ),
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
