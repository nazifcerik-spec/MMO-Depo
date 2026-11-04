import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/quantity_field.dart';
import '../../shared/widgets.dart';
import '../character/models.dart';
import '../inventory/inventory_models.dart';
import 'market_repository.dart';

const _classes = ['weapon', 'armor', 'jewelry', 'consumable', 'material'];

String itemClassName(AppLocalizations l, String c) => switch (c) {
  'weapon' => l.classWeapon,
  'armor' => l.classArmor,
  'jewelry' => l.classJewelry,
  'consumable' => l.classConsumable,
  'material' => l.classMaterial,
  _ => humanize(c),
};

/// Player-to-player market: search (class filter + sort, paged), buy with a
/// price guard, and the character's own listings with cancel.
class MarketPanel extends StatefulWidget {
  const MarketPanel({super.key, required this.repo, required this.dataVersion, required this.gold});

  final MarketRepository repo;
  final ValueNotifier<int> dataVersion;

  /// Current gold, re-read on every load (buy buttons disable when short).
  final Future<int> Function() gold;

  @override
  State<MarketPanel> createState() => _MarketPanelState();
}

class _MarketPanelState extends State<MarketPanel> with PanelState {
  List<Listing> _results = const [];
  List<Listing> _mine = const [];
  bool _hasMore = false;
  int _gold = 0;
  String? _class;
  String? _kind;
  String _sort = 'price_asc';

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
      widget.repo.search(kind: _kind, itemClass: _class, sort: _sort),
      widget.repo.mine(),
      widget.gold(),
    ]);
    _results = r[0] as List<Listing>;
    _hasMore = _results.length == MarketRepository.pageSize;
    _mine = (r[1] as List<Listing>).where((l) => l.state == 'active').toList();
    _gold = r[2] as int;
  }

  Future<void> _more() => run(() async {
    final next = await widget.repo.search(kind: _kind, itemClass: _class, sort: _sort, offset: _results.length);
    if (mounted) {
      setState(() {
        _results = [..._results, ...next];
        _hasMore = next.length == MarketRepository.pageSize;
      });
    }
  });

  void _filter({String? cls, String? sort}) {
    setState(() {
      _class = cls;
      _sort = sort ?? _sort;
    });
    run(load);
  }

  void _filterKind(String? kind) {
    setState(() => _kind = kind);
    run(load);
  }

  /// A bid or the buyout on an auction (DD-359).
  Future<void> _bid(Listing x) async {
    final l = AppLocalizations.of(context);
    final min = x.minNextBid ?? 1;
    final choice = await showDialog<_BidChoice>(
      context: context,
      builder: (_) => _BidDialog(listing: x, title: _name(x), gold: _gold),
    );
    if (choice == null) return;
    await run(() async {
      if (choice.buyout) {
        await widget.repo.buyout(x);
        if (mounted) toast(context, Text(l.bought(_name(x))));
      } else {
        await widget.repo.bid(x, choice.amount < min ? min : choice.amount);
        if (mounted) toast(context, Text(l.bidPlaced));
      }
      widget.dataVersion.value++;
      await load();
    });
  }

  Future<void> _buy(Listing listing) async {
    final l = AppLocalizations.of(context);
    final name = _name(listing);
    final ok = await confirmAction(
      context,
      title: l.buyListingTitle(name),
      body: l.buyListingBody(goldString(listing.totalPrice), listing.sellerName),
      action: l.buy,
    );
    if (!ok) return;
    await run(() async {
      await widget.repo.buy(listing);
      if (!mounted) return;
      toast(context, Text(l.bought(name)));
      widget.dataVersion.value++;
    });
  }

  Future<void> _cancel(Listing listing) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.cancelListingTitle(_name(listing)),
      body: l.cancelListingBody,
      action: l.cancelListing,
      danger: true,
    );
    if (!ok) return;
    await run(() async {
      await widget.repo.cancel(listing.id);
      widget.dataVersion.value++;
    });
  }

  String _name(Listing x) {
    final n = itemName(x.templateCode, x.prefix, x.suffix);
    return x.quantity > 1 ? '$n ×${x.quantity}' : n;
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final now = DateTime.now();
    final me = widget.repo.characterId;
    return frame(
      () => [
        Row(
          children: [
            Expanded(child: Text(l.segMarket, style: text.titleMedium)),
            Semantics(label: l.gold, child: GoldText(_gold)),
          ],
        ),
        const SizedBox(height: Space.md),
        if (_mine.isNotEmpty)
          Card(
            margin: const EdgeInsets.only(bottom: Space.md),
            child: ExpansionTile(
              title: Text(l.myListings(_mine.length)),
              shape: const Border(),
              children: [
                for (final m in _mine)
                  ListTile(
                    dense: true,
                    title: Text(_name(m), style: TextStyle(color: rarityColor(m.rarity))),
                    subtitle: Text(
                      '${goldString(m.unitPrice)} · ${l.remaining(shortDuration(l, m.expiresAt.difference(now)))}',
                    ),
                    trailing: IconButton(
                      tooltip: l.cancelListing,
                      icon: const Icon(Icons.close, color: Palette.danger),
                      onPressed: () => _cancel(m),
                    ),
                  ),
              ],
            ),
          ),
        Wrap(
          spacing: Space.sm,
          runSpacing: Space.sm,
          children: [
            ChoiceChip(label: Text(l.filterAll), selected: _class == null, onSelected: (_) => _filter()),
            for (final c in _classes)
              ChoiceChip(
                label: Text(itemClassName(l, c)),
                selected: _class == c,
                onSelected: (_) => _filter(cls: c),
              ),
          ],
        ),
        const SizedBox(height: Space.sm),
        Wrap(
          spacing: Space.sm,
          children: [
            for (final (key, label) in <(String?, String)>[
              (null, l.filterAll),
              ('fixed', l.marketFixed),
              ('auction', l.marketAuction),
            ])
              ChoiceChip(
                key: ValueKey('market-kind-${key ?? 'all'}'),
                label: Text(label),
                selected: _kind == key,
                onSelected: (_) => _filterKind(key),
              ),
          ],
        ),
        const SizedBox(height: Space.sm),
        Wrap(
          spacing: Space.sm,
          children: [
            for (final (key, label) in [('price_asc', l.sortCheapest), ('newest', l.sortNewest)])
              ChoiceChip(
                label: Text(label),
                selected: _sort == key,
                onSelected: (_) => _filter(cls: _class, sort: key),
              ),
          ],
        ),
        const SizedBox(height: Space.md),
        if (_results.isEmpty)
          Padding(
            padding: const EdgeInsets.only(top: Space.xl),
            child: Text(l.marketEmpty, textAlign: TextAlign.center, style: text.bodySmall),
          ),
        for (final x in _results)
          Card(
            margin: const EdgeInsets.only(bottom: Space.sm),
            child: ListTile(
              title: Text(_name(x), style: TextStyle(color: rarityColor(x.rarity))),
              subtitle: Text(
                x.isAuction
                    ? [
                        x.currentBid == null
                            ? l.auctionNoBids(goldString(x.minNextBid ?? 0))
                            : l.auctionTopBid(goldString(x.currentBid!), x.bidCount),
                        if (x.bidderId == me) l.auctionYouLead,
                        shortDuration(l, x.expiresAt.difference(DateTime.now())),
                      ].join(' · ')
                    : l.listingLine(x.sellerName, goldString(x.unitPrice)),
              ),
              trailing: x.isAuction
                  ? FilledButton.tonal(
                      key: ValueKey('bid-${x.id}'),
                      style: compactButton,
                      onPressed: x.sellerCharacterId != me && x.bidderId != me && (x.minNextBid ?? 0) <= _gold
                          ? () => _bid(x)
                          : null,
                      child: Text(l.bid),
                    )
                  : FilledButton.tonal(
                      style: compactButton,
                      onPressed: x.sellerCharacterId != me && x.totalPrice <= _gold ? () => _buy(x) : null,
                      child: GoldText(x.totalPrice, size: 13),
                    ),
            ),
          ),
        if (_hasMore) TextButton(onPressed: _more, child: Text(l.loadMore)),
      ],
    );
  }
}

/// Bottom sheet to list a bag item. Pops `true` once the listing is created.
Future<bool> showListingSheet(BuildContext context, MarketRepository repo, BagItem item) async {
  return await showModalBottomSheet<bool>(
        context: context,
        isScrollControlled: true,
        showDragHandle: true,
        backgroundColor: Palette.bgSurface,
        builder: (_) => _ListingSheet(repo: repo, item: item),
      ) ??
      false;
}

class _ListingSheet extends StatefulWidget {
  const _ListingSheet({required this.repo, required this.item});

  final MarketRepository repo;
  final BagItem item;

  @override
  State<_ListingSheet> createState() => _ListingSheetState();
}

class _ListingSheetState extends State<_ListingSheet> {
  final _gold = TextEditingController();
  final _silver = TextEditingController();
  final _copper = TextEditingController();
  final _buyGold = TextEditingController();
  final _buySilver = TextEditingController();
  final _buyCopper = TextEditingController();
  late int _qty = widget.item.stack;
  int _hours = 24;
  bool _auction = false;
  PriceStats? _stats;
  bool _busy = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    for (final c in [_gold, _silver, _copper, _buyGold, _buySilver, _buyCopper]) {
      c.addListener(() => setState(() {}));
    }
    _loadStats();
  }

  @override
  void dispose() {
    for (final c in [_gold, _silver, _copper, _buyGold, _buySilver, _buyCopper]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _loadStats() async {
    try {
      final s = await widget.repo.stats(widget.item.templateCode);
      if (!mounted) return;
      setState(() => _stats = s);
      if (s.suggested != null && _unit == 0) _setUnit(s.suggested!);
    } on ApiError {
      // Suggestion only; the player can still type a price.
    }
  }

  void _setUnit(int copper) {
    final p = splitCopper(copper);
    _gold.text = p.gold == 0 ? '' : '${p.gold}';
    _silver.text = p.silver == 0 ? '' : '${p.silver}';
    _copper.text = '${p.copper}';
  }

  int get _unit {
    int v(TextEditingController c) => int.tryParse(c.text) ?? 0;
    return v(_gold) * 10000 + v(_silver) * 100 + v(_copper);
  }

  int get _buyout {
    int v(TextEditingController c) => int.tryParse(c.text) ?? 0;
    return v(_buyGold) * 10000 + v(_buySilver) * 100 + v(_buyCopper);
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      if (_auction) {
        await widget.repo.sellAuction(
          itemId: widget.item.instanceId,
          quantity: _qty,
          startBid: _unit,
          buyout: _buyout,
          hours: _hours,
        );
        if (mounted) Navigator.of(context).pop(true);
        return;
      }
      await widget.repo.sell(itemId: widget.item.instanceId, quantity: _qty, unitPrice: _unit, hours: _hours);
      if (mounted) Navigator.of(context).pop(true);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final item = widget.item;
    final total = _unit * _qty;
    Widget coin(TextEditingController c, String unit, Color color) => Expanded(
      child: TextField(
        controller: c,
        keyboardType: TextInputType.number,
        inputFormatters: [FilteringTextInputFormatter.digitsOnly, LengthLimitingTextInputFormatter(6)],
        decoration: InputDecoration(
          suffixText: unit,
          suffixStyle: TextStyle(color: color, fontWeight: FontWeight.w700),
        ),
      ),
    );

    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: PageFrame(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(l.listOnMarket, style: text.titleLarge),
              const SizedBox(height: Space.xs),
              Text(
                itemName(item.templateCode, item.prefix, item.suffix),
                style: TextStyle(color: rarityColor(item.rarity), fontWeight: FontWeight.w600),
              ),
              if (item.stack > 1) ...[
                const SizedBox(height: Space.lg),
                Row(
                  children: [
                    Expanded(child: Text(l.quantityOf(_qty, item.stack), style: text.bodySmall)),
                    QuantityField(
                      key: const Key('listing-qty'),
                      value: _qty,
                      min: 1,
                      max: item.stack,
                      onChanged: (v) => setState(() => _qty = v),
                    ),
                  ],
                ),
                Slider(
                  value: _qty.toDouble(),
                  min: 1,
                  max: item.stack.toDouble(),
                  divisions: item.stack - 1,
                  label: '$_qty',
                  onChanged: (v) => setState(() => _qty = v.round()),
                ),
              ],
              const SizedBox(height: Space.lg),
              Wrap(
                spacing: Space.sm,
                children: [
                  ChoiceChip(
                    key: const Key('listing-fixed'),
                    label: Text(l.marketFixed),
                    selected: !_auction,
                    onSelected: (_) => setState(() => _auction = false),
                  ),
                  ChoiceChip(
                    key: const Key('listing-auction'),
                    label: Text(l.marketAuction),
                    selected: _auction,
                    onSelected: (_) => setState(() => _auction = true),
                  ),
                ],
              ),
              const SizedBox(height: Space.lg),
              Text(_auction ? l.auctionStartBid : l.unitPrice, style: text.bodySmall),
              const SizedBox(height: Space.sm),
              Row(
                children: [
                  coin(_gold, 'g', Palette.accent),
                  const SizedBox(width: Space.sm),
                  coin(_silver, 's', const Color(0xFFC9CDD6)),
                  const SizedBox(width: Space.sm),
                  coin(_copper, 'c', const Color(0xFFC07A4A)),
                ],
              ),
              if (_stats?.suggested case final s?) ...[
                const SizedBox(height: Space.xs),
                Text(l.suggestedPrice(goldString(s)), style: text.bodySmall),
              ],
              if (_auction) ...[
                const SizedBox(height: Space.lg),
                Text(l.auctionBuyoutOptional, style: text.bodySmall),
                const SizedBox(height: Space.sm),
                Row(
                  children: [
                    coin(_buyGold, 'g', Palette.accent),
                    const SizedBox(width: Space.sm),
                    coin(_buySilver, 's', const Color(0xFFC9CDD6)),
                    const SizedBox(width: Space.sm),
                    coin(_buyCopper, 'c', const Color(0xFFC07A4A)),
                  ],
                ),
                const SizedBox(height: Space.xs),
                Text(l.auctionRules, style: text.bodySmall),
              ],
              const SizedBox(height: Space.lg),
              Text(l.durationLabel, style: text.bodySmall),
              const SizedBox(height: Space.sm),
              Wrap(
                spacing: Space.sm,
                children: [
                  for (final h in listingDurations.keys)
                    ChoiceChip(
                      label: Text(shortDuration(l, Duration(hours: h))),
                      selected: h == _hours,
                      onSelected: (_) => setState(() => _hours = h),
                    ),
                ],
              ),
              const SizedBox(height: Space.lg),
              _Line(
                label: l.listingFee,
                copper: _auction
                    ? listingFeePreview(((_unit + _qty - 1) / _qty).floor(), _qty, _hours)
                    : listingFeePreview(_unit, _qty, _hours),
              ),
              _Line(
                label: _auction ? l.auctionYouReceiveAtLeast : l.youReceive,
                copper: saleProceedsPreview(_auction ? _unit : total),
              ),
              if (_error != null) ...[
                const SizedBox(height: Space.md),
                Text(_error!, style: const TextStyle(color: Palette.danger)),
              ],
              const SizedBox(height: Space.xl),
              FilledButton(
                onPressed: _busy || _unit <= 0 ? null : _submit,
                child: _busy
                    ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
                    : Text(l.listOnMarket),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _BidChoice {
  const _BidChoice(this.amount, this.buyout);

  final int amount;
  final bool buyout;
}

/// Asks for a bid in whole gold, silver and copper; offers the buyout if any.
class _BidDialog extends StatefulWidget {
  const _BidDialog({required this.listing, required this.title, required this.gold});

  final Listing listing;
  final String title;
  final int gold;

  @override
  State<_BidDialog> createState() => _BidDialogState();
}

class _BidDialogState extends State<_BidDialog> {
  late final _amount = TextEditingController(text: '${widget.listing.minNextBid ?? 1}');

  @override
  void dispose() {
    _amount.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final x = widget.listing;
    final min = x.minNextBid ?? 1;
    final amount = int.tryParse(_amount.text) ?? 0;
    return AlertDialog(
      title: Text(widget.title),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(l.auctionMinBid(goldString(min))),
          const SizedBox(height: Space.sm),
          TextField(
            key: const Key('bid-amount'),
            controller: _amount,
            keyboardType: TextInputType.number,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            onChanged: (_) => setState(() {}),
            decoration: InputDecoration(labelText: l.auctionBidCopper, helperText: goldString(amount)),
          ),
          const SizedBox(height: Space.xs),
          Text(l.auctionHeldNote, style: const TextStyle(fontSize: 12, color: Palette.textMuted)),
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.cancel)),
        if (x.buyout != null)
          TextButton(
            key: const Key('bid-buyout'),
            onPressed: x.buyout! <= widget.gold ? () => Navigator.of(context).pop(_BidChoice(0, true)) : null,
            child: Text(l.auctionBuyoutFor(goldString(x.buyout!))),
          ),
        FilledButton(
          key: const Key('bid-confirm'),
          onPressed: amount >= min && amount <= widget.gold
              ? () => Navigator.of(context).pop(_BidChoice(amount, false))
              : null,
          child: Text(l.bid),
        ),
      ],
    );
  }
}

class _Line extends StatelessWidget {
  const _Line({required this.label, required this.copper});

  final String label;
  final int copper;

  @override
  Widget build(BuildContext context) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 2),
    child: Row(
      children: [
        Expanded(child: Text(label, style: Theme.of(context).textTheme.bodySmall)),
        GoldText(copper, size: 14),
      ],
    ),
  );
}
