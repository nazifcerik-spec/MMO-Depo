import 'dart:async';

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/panel.dart';
import '../../shared/quantity_field.dart';
import '../../shared/widgets.dart';
import '../inventory/inventory_models.dart';
import '../inventory/inventory_repository.dart';

/// One item put up in a trade, with its tax (DD-315).
class TradeItem {
  const TradeItem({
    required this.instanceId,
    required this.templateCode,
    required this.quantity,
    required this.rarity,
    required this.taxCopper,
    this.prefix,
    this.suffix,
  });

  factory TradeItem.fromJson(Map<String, dynamic> j) => TradeItem(
    instanceId: j['instance_id'] as String,
    templateCode: j['template_code'] as String,
    quantity: (j['quantity'] as num?)?.toInt() ?? 1,
    rarity: j['rarity'] as String? ?? 'common',
    taxCopper: (j['tax_copper'] as num?)?.toInt() ?? 0,
    prefix: j['prefix'] as String?,
    suffix: j['suffix'] as String?,
  );

  final String instanceId;
  final String templateCode;
  final int quantity;
  final String rarity;
  final int taxCopper;
  final String? prefix;
  final String? suffix;
}

class TradeSide {
  const TradeSide({required this.name, required this.items, required this.taxCopper, required this.ready});

  factory TradeSide.fromJson(Map<String, dynamic> j) => TradeSide(
    name: j['name'] as String? ?? '',
    items: [for (final i in (j['items'] as List? ?? const []).cast<Map<String, dynamic>>()) TradeItem.fromJson(i)],
    taxCopper: (j['tax_copper'] as num?)?.toInt() ?? 0,
    ready: j['ready'] as bool? ?? false,
  );

  final String name;
  final List<TradeItem> items;
  final int taxCopper;
  final bool ready;
}

class Trade {
  const Trade({required this.id, required this.state, required this.me, required this.them});

  factory Trade.fromJson(Map<String, dynamic> j) => Trade(
    id: j['id'] as String,
    state: j['state'] as String? ?? 'open',
    me: TradeSide.fromJson(j['me'] as Map<String, dynamic>),
    them: TradeSide.fromJson(j['them'] as Map<String, dynamic>),
  );

  final String id;
  final String state;
  final TradeSide me;
  final TradeSide them;
}

class TradeRepository {
  TradeRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId/trade';

  Trade? _parse(Map<String, dynamic> j) =>
      j['trade'] is Map<String, dynamic> ? Trade.fromJson(j['trade'] as Map<String, dynamic>) : null;

  Future<Trade?> current() async => _parse(await api.get(_base));
  Future<Trade?> open(String name) async => _parse(await api.post(_base, {'target_name': name}));
  Future<Trade?> offer(String id, String itemId, int qty) async =>
      _parse(await api.post('$_base/$id/offer', {'item_id': itemId, 'quantity': qty}));
  Future<Trade?> remove(String id, String itemId) async =>
      _parse(await api.post('$_base/$id/remove', {'item_id': itemId}));
  Future<Trade?> ready(String id) async => _parse(await api.post('$_base/$id/ready', const {}));
  Future<void> cancel(String id) => api.post('$_base/$id/cancel', const {});
}

/// Trading with another player (DD-315): name a player, put items up, both
/// get ready; each side pays 8% of what it gives at the average market price.
class TradePanel extends StatefulWidget {
  const TradePanel({super.key, required this.repo, required this.bag, required this.live});

  final TradeRepository repo;
  final InventoryRepository bag;
  final LiveEvents live;

  @override
  State<TradePanel> createState() => _TradePanelState();
}

class _TradePanelState extends State<TradePanel> with PanelState, AutomaticKeepAliveClientMixin {
  Trade? _trade;
  final _name = TextEditingController();
  StreamSubscription<LiveEvent>? _sub;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _sub = widget.live.events.where((e) => e.kind == 'trade.changed').listen((_) => load());
    load();
  }

  @override
  void dispose() {
    _sub?.cancel();
    _name.dispose();
    super.dispose();
  }

  @override
  Future<void> fetch() async => _trade = await widget.repo.current();

  Future<void> _set(Future<Trade?> Function() action) => run(() async {
    final t = await action();
    if (mounted) setState(() => _trade = t);
  });

  Future<void> _add(Trade t) async {
    final l = AppLocalizations.of(context);
    final bag = await widget.bag.bag();
    if (!mounted) return;
    final offered = {for (final i in t.me.items) i.instanceId};
    final pick = await showModalBottomSheet<BagItem>(
      context: context,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => SafeArea(
        child: ListView(
          shrinkWrap: true,
          children: [
            for (final i in bag.items)
              if (!i.isLocked && !i.isBound && !offered.contains(i.instanceId))
                ListTile(
                  key: Key('pick-${i.instanceId}'),
                  title: Text(
                    stackLabel(itemName(i.templateCode, i.prefix, i.suffix), i.stack),
                    style: TextStyle(color: rarityColor(i.rarity)),
                  ),
                  onTap: () => Navigator.of(context).pop(i),
                ),
          ],
        ),
      ),
    );
    if (pick == null || !mounted) return;
    var qty = pick.stack;
    if (pick.stack > 1) {
      final n = await showDialog<int>(
        context: context,
        builder: (_) => _QtyDialog(max: pick.stack, title: l.tradeAdd),
      );
      if (n == null) return;
      qty = n;
    }
    await _set(() => widget.repo.offer(t.id, pick.instanceId, qty));
  }

  Future<void> _cancel(Trade t) async {
    final l = AppLocalizations.of(context);
    if (!await confirmAction(context, title: l.tradeCancel, body: l.tradeCancelBody, action: l.tradeCancel)) return;
    await run(() async {
      await widget.repo.cancel(t.id);
      if (mounted) setState(() => _trade = null);
    });
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(() {
      final t = _trade;
      if (t == null || t.state != 'open') {
        return [
          if (t?.state == 'done') ...[
            Text(l.tradeDone, style: text.bodyMedium?.copyWith(color: Palette.success)),
            const SizedBox(height: Space.md),
          ],
          Text(l.tradeHint, style: text.bodySmall),
          const SizedBox(height: Space.md),
          TextField(
            controller: _name,
            decoration: InputDecoration(labelText: l.tradeWith),
            textInputAction: TextInputAction.go,
            onSubmitted: (_) => _set(() => widget.repo.open(_name.text.trim())),
          ),
          const SizedBox(height: Space.sm),
          FilledButton(
            key: const Key('trade-open'),
            onPressed: () => _set(() => widget.repo.open(_name.text.trim())),
            child: Text(l.tradeOpen),
          ),
        ];
      }
      Widget side(String title, TradeSide s, {required bool mine}) => Card(
        child: Padding(
          padding: const EdgeInsets.all(Space.md),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(child: Text(title, style: text.titleSmall)),
                  Chip(
                    label: Text(s.ready ? l.tradeReady : l.tradeWaiting),
                    backgroundColor: s.ready ? Palette.success.withValues(alpha: 0.2) : null,
                  ),
                ],
              ),
              if (s.items.isEmpty) Text(l.tradeNothing, style: text.bodySmall),
              for (final i in s.items)
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        stackLabel(itemName(i.templateCode, i.prefix, i.suffix), i.quantity),
                        style: TextStyle(color: rarityColor(i.rarity)),
                      ),
                    ),
                    if (mine)
                      IconButton(
                        tooltip: l.tradeRemove,
                        icon: const Icon(Icons.close, size: 18),
                        onPressed: () => _set(() => widget.repo.remove(t.id, i.instanceId)),
                      ),
                  ],
                ),
              if (mine) ...[
                TextButton.icon(
                  key: const Key('trade-add'),
                  onPressed: () => _add(t),
                  icon: const Icon(Icons.add),
                  label: Text(l.tradeAdd),
                ),
                Row(
                  children: [
                    Text('${l.tradeTax}: ', style: text.bodySmall),
                    GoldText(s.taxCopper, size: 13),
                  ],
                ),
              ],
            ],
          ),
        ),
      );
      return [
        Text(l.tradeTitle(t.them.name), style: text.titleMedium),
        const SizedBox(height: Space.sm),
        side(l.tradeYours, t.me, mine: true),
        side(l.tradeTheirs(t.them.name), t.them, mine: false),
        const SizedBox(height: Space.md),
        FilledButton(
          key: const Key('trade-ready'),
          onPressed: t.me.ready ? null : () => _set(() => widget.repo.ready(t.id)),
          child: Text(t.me.ready ? l.tradeWaitingThem : l.tradeConfirm),
        ),
        const SizedBox(height: Space.sm),
        TextButton(
          onPressed: () => _cancel(t),
          style: TextButton.styleFrom(foregroundColor: Palette.danger),
          child: Text(l.tradeCancel),
        ),
      ];
    });
  }
}

class _QtyDialog extends StatefulWidget {
  const _QtyDialog({required this.max, required this.title});

  final int max;
  final String title;

  @override
  State<_QtyDialog> createState() => _QtyDialogState();
}

class _QtyDialogState extends State<_QtyDialog> {
  late int _qty = widget.max;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(widget.title),
      content: QuantityField(
        key: const Key('trade-qty'),
        value: _qty,
        min: 1,
        max: widget.max,
        onChanged: (v) => setState(() => _qty = v),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.cancel)),
        FilledButton(onPressed: () => Navigator.of(context).pop(_qty), child: Text(widget.title)),
      ],
    );
  }
}
