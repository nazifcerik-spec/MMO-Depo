import 'package:flutter/material.dart';

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

/// Bag + equipment tab: 12 equipment slots above the positional bag grid.
/// Tapping an item opens its detail sheet; the chosen action runs here and
/// both lists reload (the server moves items between bag and slots).
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
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

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
    if (!mounted) return;
    final action = await showModalBottomSheet<_ItemAction>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => _ItemDetail(item: item, sellMult: _sellMult),
    );
    if (action == null || !mounted) return;
    final l = AppLocalizations.of(context);
    final name = displayName(item);
    switch (action) {
      case _ItemAction.use:
        await _run(() async {
          final (hp, max) = await widget.repo.use(item.instanceId);
          if (mounted) toast(context, Text(l.usedItem(name, hp, max)));
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
                subtitle: slotName(l, slot),
              ),
              const SizedBox(height: Space.xl),
              FilledButton(onPressed: () => Navigator.pop(c, _ItemAction.unequip), child: Text(l.unequip)),
            ],
          ),
        ),
      ),
    );
    if (action == _ItemAction.unequip) await _run(() => widget.repo.unequip(slot));
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
    final bySlot = bag.bySlot;
    final text = Theme.of(context).textTheme;

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
                      semantic: '${slotName(l, slot)}: ${itemName(e.templateCode, e.prefix, e.suffix)}',
                      onTap: () => _openEquipped(slot, e),
                    )
                  else
                    ItemCell(label: slotName(l, slot)),
              ],
            ),
            const SizedBox(height: Space.xl),
            Row(
              children: [
                Expanded(child: Text(l.bag, style: text.titleMedium)),
                Text(
                  '${l.slotsUsed(bag.slotsUsed, bag.slotCap)} · ${l.weightKg((bag.weightGrams / 1000).toStringAsFixed(1), (bag.weightCapGrams / 1000).toStringAsFixed(0))}',
                  style: text.bodySmall?.copyWith(
                    color: bag.weightGrams > bag.weightCapGrams ? Palette.warning : null,
                    fontFeatures: numeric,
                  ),
                ),
              ],
            ),
            const SizedBox(height: Space.md),
            GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: bag.slotCap,
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 6,
                mainAxisSpacing: Space.sm,
                crossAxisSpacing: Space.sm,
              ),
              itemBuilder: (_, i) {
                final item = bySlot[i];
                if (item == null) return const ItemCell();
                return ItemCell(
                  rarity: item.rarity,
                  icon: itemIcon(item.itemClass, item.itemSubclass),
                  stack: item.stack,
                  locked: item.isLocked,
                  semantic: '${displayName(item)} ×${item.stack}',
                  onTap: () => _openBagItem(item),
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}

class _ItemHeader extends StatelessWidget {
  const _ItemHeader({required this.name, required this.rarity, required this.icon, required this.subtitle});

  final String name;
  final String rarity;
  final IconData icon;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final c = rarityColor(rarity);
    return Row(
      children: [
        SizedBox.square(
          dimension: 56,
          child: ItemCell(rarity: rarity, icon: icon),
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
  const _ItemDetail({required this.item, required this.sellMult});

  final BagItem item;
  final double? sellMult;

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
                    l.restoresHp(item.healHp),
                    style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
                  ),
                ),
              if (bonuses.isNotEmpty) const SizedBox(height: Space.md),
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
