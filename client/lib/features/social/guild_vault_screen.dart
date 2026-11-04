import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../inventory/bank_screen.dart';
import '../inventory/inventory_models.dart';
import '../inventory/inventory_repository.dart';
import 'group_repository.dart';

const _copperPerGold = 10000;

/// The guild's treasury and bank (DD-358): every member gives, ranks decide who
/// takes and how much a day. Every movement is in the log.
class GuildVaultScreen extends StatefulWidget {
  const GuildVaultScreen({super.key, required this.repo, required this.inventory});

  final GroupRepository repo;
  final InventoryRepository inventory;

  @override
  State<GuildVaultScreen> createState() => _GuildVaultScreenState();
}

class _GuildVaultScreenState extends State<GuildVaultScreen> {
  GuildVault? _vault;
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
      final r = await Future.wait<Object>([widget.repo.vault(), widget.inventory.bag()]);
      if (mounted) {
        setState(() {
          _vault = r[0] as GuildVault;
          _bag = r[1] as Bag;
          _error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _run(Future<GuildVault> Function() action) async {
    final l = AppLocalizations.of(context);
    setState(() => _busy = true);
    try {
      final v = await action();
      final bag = await widget.inventory.bag();
      if (mounted) {
        setState(() {
          _vault = v;
          _bag = bag;
        });
      }
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _gold({required bool deposit}) async {
    final l = AppLocalizations.of(context);
    final gold = await showDialog<int>(
      context: context,
      builder: (_) => _GoldDialog(title: deposit ? l.guildVaultDepositGold : l.guildVaultWithdrawGold),
    );
    if (gold == null || gold <= 0) return;
    final copper = gold * _copperPerGold;
    await _run(() => deposit ? widget.repo.depositGold(copper) : widget.repo.withdrawGold(copper));
  }

  Future<void> _item(BagItem item, {required bool deposit}) async {
    final l = AppLocalizations.of(context);
    var qty = item.stack;
    if (item.stack > 1) {
      final picked = await showDialog<int>(
        context: context,
        builder: (_) => QuantityDialog(item: item, title: deposit ? l.bankDeposit : l.bankWithdraw),
      );
      if (picked == null) return;
      qty = picked;
    }
    await _run(
      () => deposit ? widget.repo.depositItem(item.instanceId, qty) : widget.repo.withdrawItem(item.instanceId, qty),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final v = _vault, bag = _bag;
    return Scaffold(
      appBar: AppBar(title: Text(l.guildVault)),
      body: SafeArea(
        child: PageFrame(
          child: _error != null
              ? MessageView(message: errorText(l, _error!), onRetry: _load)
              : v == null || bag == null
              ? const Center(child: CircularProgressIndicator())
              : AbsorbPointer(
                  absorbing: _busy,
                  child: ListView(
                    padding: const EdgeInsets.all(Space.screen),
                    children: [
                      Text(l.guildVaultTreasury, style: text.titleMedium),
                      const SizedBox(height: Space.xs),
                      GoldText(v.treasuryCopper, key: const Key('vault-treasury'), size: 20),
                      const SizedBox(height: Space.xs),
                      Text(
                        v.goldLimit < 0
                            ? l.guildVaultNoLimit
                            : l.guildVaultGoldToday(goldString(v.goldToday), goldString(v.goldLimit)),
                        style: text.bodySmall,
                      ),
                      const SizedBox(height: Space.sm),
                      Row(
                        children: [
                          Expanded(
                            child: OutlinedButton(
                              key: const Key('vault-gold-in'),
                              onPressed: () => _gold(deposit: true),
                              child: Text(l.bankDeposit),
                            ),
                          ),
                          const SizedBox(width: Space.sm),
                          Expanded(
                            child: OutlinedButton(
                              key: const Key('vault-gold-out'),
                              onPressed: v.goldLimit == 0 ? null : () => _gold(deposit: false),
                              child: Text(l.bankWithdraw),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: Space.lg),
                      Text(
                        '${l.guildVaultBank} · ${l.slotsUsed(v.bank.slotsUsed, v.bank.slotCap)}',
                        style: text.titleMedium,
                      ),
                      const SizedBox(height: Space.xs),
                      if (!v.canUseItems)
                        Text(l.guildVaultNoItems, style: text.bodySmall)
                      else ...[
                        if (v.itemsLimit >= 0)
                          Text(l.guildVaultItemsToday(v.itemsToday, v.itemsLimit), style: text.bodySmall),
                        if (v.bank.items.isEmpty) Text(l.bankEmpty, style: text.bodySmall),
                        for (final i in v.bank.items)
                          StoreItemRow(
                            key: Key('vault-${i.instanceId}'),
                            item: i,
                            onTap: () => _item(i, deposit: false),
                          ),
                        const SizedBox(height: Space.md),
                        Text('${l.bag} · ${l.slotsUsed(bag.slotsUsed, bag.slotCap)}', style: text.titleMedium),
                        for (final i in bag.items)
                          if (!i.isLocked && !i.isBound)
                            StoreItemRow(
                              key: Key('vault-bag-${i.instanceId}'),
                              item: i,
                              onTap: () => _item(i, deposit: true),
                            ),
                      ],
                      const SizedBox(height: Space.lg),
                      Text(l.guildVaultLog, style: text.titleMedium),
                      const SizedBox(height: Space.xs),
                      if (v.log.isEmpty) Text(l.guildVaultLogEmpty, style: text.bodySmall),
                      for (final e in v.log)
                        Padding(
                          padding: const EdgeInsets.symmetric(vertical: 3),
                          child: Text(_line(l, e), style: text.bodySmall),
                        ),
                    ],
                  ),
                ),
        ),
      ),
    );
  }

  String _line(AppLocalizations l, GuildVaultEntry e) {
    final what = e.goldCopper != null
        ? goldString(e.goldCopper!)
        : stackLabel(itemName(e.templateCode ?? ''), e.quantity ?? 1);
    return e.action.endsWith('_in') ? l.guildVaultLogIn(e.name, what) : l.guildVaultLogOut(e.name, what);
  }
}

class _GoldDialog extends StatefulWidget {
  const _GoldDialog({required this.title});

  final String title;

  @override
  State<_GoldDialog> createState() => _GoldDialogState();
}

class _GoldDialogState extends State<_GoldDialog> {
  final _gold = TextEditingController();

  @override
  void dispose() {
    _gold.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(widget.title),
      content: TextField(
        key: const Key('vault-gold-amount'),
        controller: _gold,
        autofocus: true,
        keyboardType: TextInputType.number,
        inputFormatters: [FilteringTextInputFormatter.digitsOnly, LengthLimitingTextInputFormatter(9)],
        decoration: InputDecoration(labelText: l.guildVaultGoldAmount),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.cancel)),
        FilledButton(
          key: const Key('vault-gold-confirm'),
          onPressed: () => Navigator.of(context).pop(int.tryParse(_gold.text)),
          child: Text(widget.title),
        ),
      ],
    );
  }
}
