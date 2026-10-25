import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/widgets.dart';
import '../market/market_repository.dart';
import '../market/market_view.dart';
import 'inventory_models.dart';
import 'inventory_repository.dart';
import 'item_widgets.dart';

enum _ItemAction { use, equip, unequip, sell, list, destroy }

/// Bag + equipment tab: 12 equipment slots above the bag grid.
/// Tapping an item opens its detail sheet; the chosen action runs here and
/// both lists reload (the server moves items between bag and slots).
/// The bag can be shown in another order (display only), and a long press —
/// dragging on across cells — selects several items to sell or destroy.
class InventoryView extends StatefulWidget {
  const InventoryView({super.key, required this.repo, required this.dataVersion});

  final InventoryRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<InventoryView> createState() => _InventoryViewState();
}

class _InventoryViewState extends State<InventoryView> {
  Bag? _bag;
  Map<String, EquippedItem> _equipped = const {};
  double? _sellMult;
  Object? _error;
  bool _busy = false;
  BagSort _sort = BagSort.slot;

  /// Selected instance ids; non-empty = selection mode.
  final Set<String> _selected = {};

  /// While a long-press drag runs: true adds the cells it crosses, false
  /// removes them (decided by the first cell).
  bool? _dragAdds;

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(_load);
    _load();
  }

  @override
  void dispose() {
    widget.dataVersion.removeListener(_load);
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final r = await Future.wait<Object>([widget.repo.bag(), widget.repo.equipment()]);
      if (mounted) {
        setState(() {
          _bag = r[0] as Bag;
          _equipped = r[1] as Map<String, EquippedItem>;
          _error = null;
          final ids = {for (final i in _bag!.items) i.instanceId};
          _selected.retainWhere(ids.contains);
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Map<String, ItemSet> _sets = const {};

  /// The set a piece belongs to (DD-277), loaded once; set info is optional.
  Future<ItemSet?> _setOf(String template) async {
    if (_sets.isEmpty) {
      try {
        _sets = await widget.repo.setsByPiece();
      } on ApiError {
        return null;
      }
    }
    return _sets[template];
  }

  Iterable<String> get _wornCodes => _equipped.values.map((e) => e.templateCode);

  Future<void> _ensureSellMult() async {
    if (_sellMult != null) return;
    try {
      _sellMult = await widget.repo.sellMultiplier();
    } on ApiError {
      // Price preview is optional; the sale itself is priced by the server.
    }
  }

  Future<void> _openBagItem(BagItem item) async {
    await _ensureSellMult();
    final set = item.equippable ? await _setOf(item.templateCode) : null;
    if (!mounted) return;
    final action = await showModalBottomSheet<_ItemAction>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => _ItemDetail(item: item, sellMult: _sellMult, set: set, worn: _wornCodes),
    );
    if (action == null || !mounted) return;
    final l = AppLocalizations.of(context);
    final name = displayName(item);
    switch (action) {
      case _ItemAction.use:
        await _run(() async {
          final (hp, max, buff) = await widget.repo.use(item.instanceId);
          if (!mounted) return;
          toast(
            context,
            Text(
              [l.usedItem(name, hp, max), if (buff != null) l.fedToast(buff.attackPct, buff.defensePct)].join(' · '),
            ),
          );
          widget.dataVersion.value++; // the sheet shows the buff
        });
      case _ItemAction.equip:
        await _run(() => widget.repo.equip(item.instanceId));
      case _ItemAction.sell:
        if (await confirmAction(context, title: l.sellTitle(name), body: l.sellBody, action: l.sell)) {
          await _run(() async {
            final got = await widget.repo.sell(item.instanceId);
            if (mounted) _toastGold(l.soldFor, got);
          });
        }
      case _ItemAction.list:
        if (await showListingSheet(context, MarketRepository(widget.repo.api, widget.repo.characterId), item)) {
          if (mounted) toast(context, Text(l.listed(name)));
          widget.dataVersion.value++;
        }
      case _ItemAction.destroy:
        if (await confirmAction(
          context,
          title: l.destroyTitle(name),
          body: l.destroyBody,
          action: l.destroy,
          danger: true,
        )) {
          await _run(() => widget.repo.destroy(item.instanceId));
        }
      case _ItemAction.unequip:
        break;
    }
  }

  Future<void> _openEquipped(String slot, EquippedItem item) async {
    final set = await _setOf(item.templateCode);
    if (!mounted) return;
    final l = AppLocalizations.of(context);
    final action = await showModalBottomSheet<_ItemAction>(
      context: context,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (c) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _ItemHeader(
                name: itemName(item.templateCode, item.prefix, item.suffix),
                rarity: item.rarity,
                icon: itemIcon(item.itemClass),
                code: item.templateCode,
                subtitle: slotName(l, slot),
              ),
              if (set != null) ...[const SizedBox(height: Space.lg), SetSection(set: set, worn: _wornCodes)],
              const SizedBox(height: Space.xl),
              FilledButton(onPressed: () => Navigator.pop(c, _ItemAction.unequip), child: Text(l.unequip)),
            ],
          ),
        ),
      ),
    );
    if (action == _ItemAction.unequip) await _run(() => widget.repo.unequip(slot));
  }

  /// The bag as shown: server positions, or every item in [_sort] order.
  List<BagItem?> _cells(Bag bag) {
    if (_sort == BagSort.slot) {
      final by = bag.bySlot;
      return [for (var i = 0; i < bag.slotCap; i++) by[i]];
    }
    final sorted = sortBag(bag.items, _sort, displayName);
    final n = sorted.length > bag.slotCap ? sorted.length : bag.slotCap;
    return [for (var i = 0; i < n; i++) i < sorted.length ? sorted[i] : null];
  }

  void _toggle(BagItem item, {bool? add}) {
    if (item.isLocked) return; // escrowed items can't be sold or destroyed
    setState(() {
      if (add ?? !_selected.contains(item.instanceId)) {
        _selected.add(item.instanceId);
      } else {
        _selected.remove(item.instanceId);
      }
    });
  }

  void _dragSelect(BagItem? item, {bool start = false}) {
    if (item == null || item.isLocked) return;
    if (start) {
      HapticFeedback.selectionClick();
      _dragAdds = !_selected.contains(item.instanceId);
    }
    final add = _dragAdds;
    if (add == null || _selected.contains(item.instanceId) == add) return;
    _toggle(item, add: add);
  }

  List<BagItem> _selectedItems() => [
    for (final i in _bag?.items ?? const <BagItem>[])
      if (_selected.contains(i.instanceId)) i,
  ];

  Future<void> _sellSelected() async {
    final l = AppLocalizations.of(context);
    final items = _selectedItems();
    if (items.isEmpty ||
        !await confirmAction(context, title: l.sellSelectedTitle(items.length), body: l.sellBody, action: l.sell)) {
      return;
    }
    var total = 0;
    await _run(() async {
      try {
        for (final i in items) {
          total += await widget.repo.sell(i.instanceId);
        }
      } finally {
        _selected.clear();
        if (total > 0 && mounted) _toastGold(l.soldFor, total);
      }
    });
  }

  Future<void> _destroySelected() async {
    final l = AppLocalizations.of(context);
    final items = _selectedItems();
    if (items.isEmpty ||
        !await confirmAction(
          context,
          title: l.destroySelectedTitle(items.length),
          body: l.destroyBody,
          action: l.destroy,
          danger: true,
        )) {
      return;
    }
    await _run(() async {
      try {
        for (final i in items) {
          await widget.repo.destroy(i.instanceId);
        }
      } finally {
        _selected.clear();
      }
    });
  }

  void _toastGold(String label, int copper) {
    toast(context, Row(children: [Text('$label '), GoldText(copper, size: 14)]));
  }

  Future<void> _run(Future<void> Function() action) async {
    setState(() => _busy = true);
    try {
      await action();
      // Gear changes the sheet and a sale changes gold: bump every tab
      // (this view reloads through the same listener).
      widget.dataVersion.value++;
    } on ApiError catch (e) {
      if (mounted && !e.isSessionError) {
        toast(context, Text(errorText(AppLocalizations.of(context), e)));
      }
      await _load();
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final bag = _bag;
    if (bag == null) {
      return _error != null
          ? MessageView(message: errorText(l, _error!), onRetry: _load)
          : const Center(child: CircularProgressIndicator());
    }
    final cells = _cells(bag);
    final text = Theme.of(context).textTheme;
    final selecting = _selected.isNotEmpty;

    return RefreshIndicator(
      onRefresh: _load,
      child: AbsorbPointer(
        absorbing: _busy,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
          children: [
            Text(l.equipment, style: text.titleMedium),
            const SizedBox(height: Space.md),
            GridView.count(
              crossAxisCount: 6,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              mainAxisSpacing: Space.sm,
              crossAxisSpacing: Space.sm,
              children: [
                for (final slot in equipSlots)
                  if (_equipped[slot] case final e?)
                    ItemCell(
                      rarity: e.rarity,
                      icon: itemIcon(e.itemClass),
                      code: e.templateCode,
                      semantic: '${slotName(l, slot)}: ${itemName(e.templateCode, e.prefix, e.suffix)}',
                      onTap: () => _openEquipped(slot, e),
                    )
                  else
                    ItemCell(label: slotName(l, slot)),
              ],
            ),
            const SizedBox(height: Space.xl),
            SizedBox(
              height: 48,
              child: selecting
                  ? _SelectionBar(
                      count: _selected.length,
                      onClear: () => setState(_selected.clear),
                      onSelectAll: () => setState(() {
                        _selected.addAll([
                          for (final i in bag.items)
                            if (!i.isLocked) i.instanceId,
                        ]);
                      }),
                      onSell: _sellSelected,
                      onDestroy: _destroySelected,
                    )
                  : Row(
                      children: [
                        Expanded(child: Text(l.bag, style: text.titleMedium)),
                        PopupMenuButton<BagSort>(
                          tooltip: l.sortLabel,
                          initialValue: _sort,
                          onSelected: (v) => setState(() => _sort = v),
                          itemBuilder: (_) => [
                            for (final v in BagSort.values) PopupMenuItem(value: v, child: Text(_sortName(l, v))),
                          ],
                          child: Padding(
                            padding: const EdgeInsets.symmetric(vertical: Space.xs),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                const Icon(Icons.sort, size: 18, color: Palette.textSecondary),
                                const SizedBox(width: Space.xs),
                                Text(_sortName(l, _sort), style: text.bodySmall),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
            ),
            Text(
              '${l.slotsUsed(bag.slotsUsed, bag.slotCap)} · ${l.weightKg((bag.weightGrams / 1000).toStringAsFixed(1), (bag.weightCapGrams / 1000).toStringAsFixed(0))}',
              style: text.bodySmall?.copyWith(
                // Full bag (DD-246): buying, crafting and new runs are refused.
                color: bag.weightGrams > bag.weightCapGrams || bag.slotsUsed >= bag.slotCap ? Palette.warning : null,
                fontFeatures: numeric,
              ),
            ),
            const SizedBox(height: Space.md),
            LayoutBuilder(
              builder: (context, box) {
                const cols = 5; // five items per row
                final step = (box.maxWidth + Space.sm) / cols; // cell + gap; cells are square
                // From the global position against the grid's current place:
                // the header changes when selection starts, so the pointer's
                // own local coordinates would be stale.
                BagItem? at(Offset global) {
                  final p = (context.findRenderObject()! as RenderBox).globalToLocal(global);
                  if (p.dx < 0 || p.dy < 0 || p.dx >= box.maxWidth) return null;
                  final i = (p.dy ~/ step) * cols + (p.dx ~/ step);
                  return i < cells.length ? cells[i] : null;
                }

                return GestureDetector(
                  onLongPressStart: (d) => _dragSelect(at(d.globalPosition), start: true),
                  onLongPressMoveUpdate: (d) => _dragSelect(at(d.globalPosition)),
                  onLongPressEnd: (_) => _dragAdds = null,
                  child: GridView.builder(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    itemCount: cells.length,
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: cols,
                      mainAxisSpacing: Space.sm,
                      crossAxisSpacing: Space.sm,
                    ),
                    itemBuilder: (_, i) {
                      final item = cells[i];
                      if (item == null) return const ItemCell();
                      return ItemCell(
                        rarity: item.rarity,
                        icon: itemIcon(item.itemClass, item.itemSubclass),
                        code: item.templateCode,
                        stack: item.stack,
                        locked: item.isLocked,
                        selected: _selected.contains(item.instanceId),
                        semantic: '${displayName(item)} ×${item.stack}',
                        onTap: () => selecting ? _toggle(item) : _openBagItem(item),
                      );
                    },
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}

String _sortName(AppLocalizations l, BagSort s) => switch (s) {
  BagSort.slot => l.sortSlot,
  BagSort.rarity => l.sortRarity,
  BagSort.type => l.sortType,
  BagSort.name => l.sortName,
  BagSort.value => l.sortValue,
};

/// Replaces the bag header while items are selected.
class _SelectionBar extends StatelessWidget {
  const _SelectionBar({
    required this.count,
    required this.onClear,
    required this.onSelectAll,
    required this.onSell,
    required this.onDestroy,
  });

  final int count;
  final VoidCallback onClear;
  final VoidCallback onSelectAll;
  final VoidCallback onSell;
  final VoidCallback onDestroy;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Row(
      children: [
        IconButton(tooltip: l.clearSelection, icon: const Icon(Icons.close), onPressed: onClear),
        Expanded(child: Text(l.selectedCount(count), style: Theme.of(context).textTheme.titleMedium)),
        TextButton(onPressed: onSelectAll, child: Text(l.selectAll)),
        IconButton(tooltip: l.sell, icon: const Icon(Icons.sell_outlined), onPressed: onSell),
        IconButton(
          tooltip: l.destroy,
          icon: const Icon(Icons.delete_outline, color: Palette.danger),
          onPressed: onDestroy,
        ),
      ],
    );
  }
}

class _ItemHeader extends StatelessWidget {
  const _ItemHeader({
    required this.name,
    required this.rarity,
    required this.icon,
    required this.code,
    required this.subtitle,
  });

  final String name;
  final String rarity;
  final IconData icon;
  final String code;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final c = rarityColor(rarity);
    return Row(
      children: [
        SizedBox.square(
          dimension: 56,
          child: ItemCell(rarity: rarity, icon: icon, code: code),
        ),
        const SizedBox(width: Space.lg),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(name, style: Theme.of(context).textTheme.titleMedium?.copyWith(color: c)),
              const SizedBox(height: 2),
              Text('${rarityName(l, rarity)} · $subtitle', style: Theme.of(context).textTheme.bodySmall),
            ],
          ),
        ),
      ],
    );
  }
}

class _ItemDetail extends StatelessWidget {
  const _ItemDetail({required this.item, required this.sellMult, this.set, this.worn = const []});

  final BagItem item;
  final double? sellMult;

  /// The raid set the item belongs to (DD-277) and the worn template codes.
  final ItemSet? set;
  final Iterable<String> worn;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final bonuses = [for (final b in item.bonuses) ?bonusText(l, b)];
    final price = sellMult == null ? null : (item.baseValueCopper * sellMult! * item.stack).floor();

    Widget row(String k, Widget v) => Padding(
      padding: const EdgeInsets.symmetric(vertical: 3),
      child: Row(
        children: [
          Expanded(child: Text(k, style: text.bodySmall)),
          v,
        ],
      ),
    );

    return SafeArea(
      child: PageFrame(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _ItemHeader(
                name: displayName(item),
                rarity: item.rarity,
                icon: itemIcon(item.itemClass, item.itemSubclass),
                code: item.templateCode,
                subtitle: item.equippable ? slotName(l, item.equipSlot!) : itemName(item.itemSubclass),
              ),
              const SizedBox(height: Space.lg),
              for (final b in bonuses)
                Padding(
                  padding: const EdgeInsets.only(bottom: 2),
                  child: Text(
                    b,
                    style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
                  ),
                ),
              if (item.usable)
                Padding(
                  padding: const EdgeInsets.only(bottom: Space.md),
                  child: Text(
                    [
                      if (item.healHp > 0) l.restoresHp(item.healHp),
                      if (item.healPct > 0) l.restoresHpPct(item.healPct),
                      if (item.buff case final b?) l.foodBuffLine(b.attackPct, b.defensePct, b.minutes),
                    ].join('\n'),
                    style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
                  ),
                ),
              if (item.manaPct > 0)
                Padding(
                  padding: const EdgeInsets.only(bottom: Space.md),
                  child: Text(
                    '${l.restoresMpPct(item.manaPct)}\n${l.huntSupplyOnly}',
                    style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
                  ),
                ),
              if (bonuses.isNotEmpty) const SizedBox(height: Space.md),
              if (set case final s?) ...[SetSection(set: s, worn: worn), const SizedBox(height: Space.md)],
              if (item.stack > 1) row(l.quantity, Text('${item.stack} / ${item.maxStack}')),
              if (item.minLevel > 1) row(l.requiresLevel, Text('${item.minLevel}')),
              if (item.classRestriction != null) row(l.characterClass, Text(className(l, item.classRestriction!))),
              row(l.weight, Text(l.weightGrams(item.weightGramsPer * item.stack))),
              if (price != null) row(l.sellValue, GoldText(price, size: 14)),
              if (item.isBound) row(l.bound, const Icon(Icons.link, size: 16, color: Palette.textSecondary)),
              const SizedBox(height: Space.xl),
              if (item.isLocked)
                Text(
                  l.itemLocked,
                  textAlign: TextAlign.center,
                  style: const TextStyle(color: Palette.warning),
                )
              else ...[
                if (item.usable)
                  FilledButton(onPressed: () => Navigator.pop(context, _ItemAction.use), child: Text(l.useItem)),
                if (item.equippable)
                  FilledButton(onPressed: () => Navigator.pop(context, _ItemAction.equip), child: Text(l.equip)),
                if (!item.isBound) ...[
                  const SizedBox(height: Space.sm),
                  FilledButton.tonal(
                    onPressed: () => Navigator.pop(context, _ItemAction.list),
                    child: Text(l.listOnMarket),
                  ),
                ],
                const SizedBox(height: Space.sm),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context, _ItemAction.sell),
                        child: Text(l.sell),
                      ),
                    ),
                    const SizedBox(width: Space.sm),
                    Expanded(
                      child: OutlinedButton(
                        style: OutlinedButton.styleFrom(foregroundColor: Palette.danger),
                        onPressed: () => Navigator.pop(context, _ItemAction.destroy),
                        child: Text(l.destroy),
                      ),
                    ),
                  ],
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
