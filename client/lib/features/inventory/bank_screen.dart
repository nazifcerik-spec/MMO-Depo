import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/quantity_field.dart';
import '../../shared/widgets.dart';
import 'inventory_models.dart';
import 'inventory_repository.dart';
import 'item_images.dart';
import 'item_widgets.dart';

/// The account bank (DD-314): 200 slots every character of the account
/// shares. Tap a bag item to deposit it, a banked one to withdraw it.
class BankScreen extends StatefulWidget {
  const BankScreen({super.key, required this.repo, required this.dataVersion});

  final InventoryRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<BankScreen> createState() => _BankScreenState();
}

class _BankScreenState extends State<BankScreen> {
  Bank? _bank;
  Bag? _bag;
  Object? _error;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final r = await Future.wait<Object>([widget.repo.bank(), widget.repo.bag()]);
      if (mounted) {
        setState(() {
          _bank = r[0] as Bank;
          _bag = r[1] as Bag;
          _error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _move(BagItem item, {required bool deposit}) async {
    final l = AppLocalizations.of(context);
    var qty = item.stack;
    if (item.stack > 1) {
      final picked = await showDialog<int>(
        context: context,
        builder: (_) => _QuantityDialog(item: item, title: deposit ? l.bankDeposit : l.bankWithdraw),
      );
      if (picked == null) return;
      qty = picked;
    }
    setState(() => _busy = true);
    try {
      final b = deposit
          ? await widget.repo.deposit(item.instanceId, qty)
          : await widget.repo.withdraw(item.instanceId, qty);
      final bag = await widget.repo.bag();
      widget.dataVersion.value++;
      if (mounted) {
        setState(() {
          _bank = b;
          _bag = bag;
        });
      }
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final bank = _bank, bag = _bag;
    return Scaffold(
      appBar: AppBar(title: Text(l.bankTitle)),
      body: SafeArea(
        child: PageFrame(
          child: _error != null
              ? MessageView(message: errorText(l, _error!), onRetry: _load)
              : bank == null || bag == null
              ? const Center(child: CircularProgressIndicator())
              : AbsorbPointer(
                  absorbing: _busy,
                  child: ListView(
                    padding: const EdgeInsets.all(Space.screen),
                    children: [
                      Text(l.bankHint, style: text.bodySmall),
                      const SizedBox(height: Space.md),
                      Text('${l.bankTitle} · ${l.slotsUsed(bank.slotsUsed, bank.slotCap)}', style: text.titleMedium),
                      const SizedBox(height: Space.sm),
                      if (bank.items.isEmpty) Text(l.bankEmpty, style: text.bodySmall),
                      for (final i in bank.items)
                        _ItemRow(key: Key('bank-${i.instanceId}'), item: i, onTap: () => _move(i, deposit: false)),
                      const SizedBox(height: Space.lg),
                      Text('${l.bag} · ${l.slotsUsed(bag.slotsUsed, bag.slotCap)}', style: text.titleMedium),
                      const SizedBox(height: Space.sm),
                      for (final i in bag.items)
                        if (!i.isLocked && !i.isBound)
                          _ItemRow(key: Key('bag-${i.instanceId}'), item: i, onTap: () => _move(i, deposit: true)),
                    ],
                  ),
                ),
        ),
      ),
    );
  }
}

class _ItemRow extends StatelessWidget {
  const _ItemRow({super.key, required this.item, required this.onTap});

  final BagItem item;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) => ListTile(
    dense: true,
    contentPadding: EdgeInsets.zero,
    leading: ItemArt(
      code: item.templateCode,
      icon: itemIcon(item.itemClass, item.itemSubclass),
      color: rarityColor(item.rarity),
    ),
    title: Text(
      stackLabel(itemName(item.templateCode, item.prefix, item.suffix), item.stack),
      style: TextStyle(color: rarityColor(item.rarity)),
    ),
    onTap: onTap,
  );
}

class _QuantityDialog extends StatefulWidget {
  const _QuantityDialog({required this.item, required this.title});

  final BagItem item;
  final String title;

  @override
  State<_QuantityDialog> createState() => _QuantityDialogState();
}

class _QuantityDialogState extends State<_QuantityDialog> {
  late int _qty = widget.item.stack;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(widget.title),
      content: Row(
        children: [
          Expanded(child: Text(itemName(widget.item.templateCode, widget.item.prefix, widget.item.suffix))),
          QuantityField(
            key: const Key('bank-qty'),
            value: _qty,
            min: 1,
            max: widget.item.stack,
            onChanged: (v) => setState(() => _qty = v),
          ),
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.cancel)),
        FilledButton(onPressed: () => Navigator.of(context).pop(_qty), child: Text(widget.title)),
      ],
    );
  }
}
