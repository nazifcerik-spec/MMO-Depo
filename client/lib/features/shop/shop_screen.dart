import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'cosmetics.dart';

class ShopItem {
  const ShopItem({required this.code, required this.kind, required this.price, required this.owned});

  factory ShopItem.fromJson(Map<String, dynamic> j) => ShopItem(
    code: j['code'] as String,
    kind: j['kind'] as String,
    price: (j['price_credits'] as num).toInt(),
    owned: j['owned'] as bool? ?? false,
  );

  final String code, kind;
  final int price;
  final bool owned;
}

class ShopView {
  const ShopView({required this.credits, required this.items, required this.paymentsEnabled});

  factory ShopView.fromJson(Map<String, dynamic> j) => ShopView(
    credits: (j['credits'] as num?)?.toInt() ?? 0,
    items: [for (final i in (j['items'] as List? ?? const []).cast<Map<String, dynamic>>()) ShopItem.fromJson(i)],
    paymentsEnabled: j['payments_enabled'] as bool? ?? false,
  );

  final int credits;
  final List<ShopItem> items;
  final bool paymentsEnabled;
}

class CosmeticProfile {
  const CosmeticProfile({required this.avatar, required this.background, required this.owned});

  factory CosmeticProfile.fromJson(Map<String, dynamic> j) => CosmeticProfile(
    avatar: j['avatar'] as String? ?? 'av01',
    background: j['background'] as String? ?? 'bg01',
    owned: {for (final c in (j['owned'] as List? ?? const [])) c as String},
  );

  final String avatar, background;
  final Set<String> owned;

  bool usable(String code) => isFree(code) || owned.contains(code);
}

/// Boost state for one category (DD-320).
class BoostState {
  const BoostState({required this.category, this.endsAt, required this.adsLeft});

  factory BoostState.fromJson(Map<String, dynamic> j) => BoostState(
    category: j['category'] as String,
    endsAt: DateTime.tryParse(j['ends_at'] as String? ?? ''),
    adsLeft: (j['ads_left'] as num?)?.toInt() ?? 0,
  );

  final String category;
  final DateTime? endsAt;
  final int adsLeft;

  bool get active => endsAt != null && endsAt!.isAfter(DateTime.now());
}

class ShopRepository {
  ShopRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId';

  Future<ShopView> shop() async => ShopView.fromJson(await api.get('$_base/shop'));
  Future<ShopView> buy(String code) async => ShopView.fromJson(await api.post('$_base/shop/$code/buy', const {}));
  Future<CosmeticProfile> cosmetics() async => CosmeticProfile.fromJson(await api.get('$_base/cosmetics'));
  Future<CosmeticProfile> select(String code) async =>
      CosmeticProfile.fromJson(await api.post('$_base/cosmetics/select', {'code': code}));

  /// Boosts and whether ads can be watched here (DD-320).
  Future<(List<BoostState>, bool)> boosts() async {
    final j = await api.get('$_base/boosts');
    return (
      [for (final b in (j['boosts'] as List? ?? const []).cast<Map<String, dynamic>>()) BoostState.fromJson(b)],
      j['ads_enabled'] as bool? ?? false,
    );
  }

  /// A development build's stand-in for a watched ad.
  Future<void> devAd(String category) => api.post('$_base/boosts/$category/dev-ad', const {});
}

String boostName(AppLocalizations l, String category) => switch (category) {
  'gather' => l.boostGather,
  'adventure' => l.boostAdventure,
  _ => l.boostCraft,
};

/// The credit shop (DD-321): credits, boosts and cosmetics.
class ShopScreen extends StatefulWidget {
  const ShopScreen({super.key, required this.repo, this.onProfile});

  final ShopRepository repo;

  /// Called after the avatar or background changed.
  final ValueChanged<CosmeticProfile>? onProfile;

  @override
  State<ShopScreen> createState() => _ShopScreenState();
}

class _ShopScreenState extends State<ShopScreen> with PanelState, SingleTickerProviderStateMixin {
  ShopView? _shop;
  CosmeticProfile? _profile;
  late final _tabs = TabController(length: 3, vsync: this);

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  void dispose() {
    _tabs.dispose();
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object>([widget.repo.shop(), widget.repo.cosmetics()]);
    _shop = r[0] as ShopView;
    _profile = r[1] as CosmeticProfile;
  }

  Future<void> _buy(ShopItem it, String title) async {
    final l = AppLocalizations.of(context);
    if (!await confirmAction(context, title: title, body: l.shopBuyBody(it.price), action: l.shopBuy)) return;
    await run(() async {
      final s = await widget.repo.buy(it.code);
      if (mounted) setState(() => _shop = s);
      if (it.kind != 'boost') await _select(it.code);
    });
  }

  Future<void> _select(String code) async {
    final p = await widget.repo.select(code);
    if (!mounted) return;
    setState(() => _profile = p);
    widget.onProfile?.call(p);
  }

  Future<void> _tap(String code, String kind) async {
    final p = _profile!;
    if (p.usable(code)) {
      await run(() => _select(code));
      return;
    }
    final it = _shop!.items.firstWhere((i) => i.code == code);
    await _buy(it, AppLocalizations.of(context).shopBuy);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(
        title: Text(l.shopTitle),
        bottom: TabBar(
          controller: _tabs,
          tabs: [
            Tab(text: l.shopCredits),
            Tab(text: l.shopAvatars),
            Tab(text: l.shopBackgrounds),
          ],
        ),
      ),
      body: SafeArea(
        child: PageFrame(
          child: !loaded
              ? frame(() => const [])
              : TabBarView(
                  controller: _tabs,
                  children: [
                    frame(() {
                      final s = _shop!;
                      return [
                        Text(l.shopBalance(s.credits), key: const Key('credits'), style: text.titleLarge),
                        const SizedBox(height: Space.sm),
                        Card(
                          child: Padding(
                            padding: const EdgeInsets.all(Space.md),
                            child: Text(
                              s.paymentsEnabled ? l.shopBuyCredits : l.shopPaymentsSoon,
                              style: text.bodySmall,
                            ),
                          ),
                        ),
                        const SizedBox(height: Space.md),
                        Text(l.shopBoosts, style: text.titleMedium),
                        for (final it in s.items)
                          if (it.kind == 'boost')
                            ListTile(
                              contentPadding: EdgeInsets.zero,
                              title: Text(boostName(l, it.code.replaceFirst('boost_', ''))),
                              subtitle: Text(l.boostDesc),
                              trailing: FilledButton(
                                key: Key('buy-${it.code}'),
                                onPressed: s.credits >= it.price ? () => _buy(it, l.shopBuy) : null,
                                child: Text(l.shopPrice(it.price)),
                              ),
                            ),
                      ];
                    }),
                    _grid(
                      count: avatarCount,
                      code: avatarCode,
                      kind: 'avatar',
                      tile: (c) => AvatarBadge(code: c, size: 64),
                    ),
                    _grid(
                      count: backgroundCount,
                      code: backgroundCode,
                      kind: 'background',
                      columns: 3,
                      aspect: 9 / 16,
                      tile: (c) => ClipRRect(
                        borderRadius: BorderRadius.circular(Radii.sm),
                        child: BackgroundArt(code: c, thumbnail: true),
                      ),
                    ),
                  ],
                ),
        ),
      ),
    );
  }

  Widget _grid({
    required int count,
    required String Function(int) code,
    required String kind,
    required Widget Function(String) tile,
    int columns = 4,
    double aspect = 1,
  }) {
    final l = AppLocalizations.of(context);
    final p = _profile!;
    final price = _shop!.items.where((i) => i.kind == kind).firstOrNull?.price ?? 0;
    return GridView.count(
      padding: const EdgeInsets.all(Space.screen),
      crossAxisCount: columns,
      childAspectRatio: aspect,
      mainAxisSpacing: Space.sm,
      crossAxisSpacing: Space.sm,
      children: [
        for (var n = 1; n <= count; n++)
          InkWell(
            key: Key('cos-${code(n)}'),
            onTap: busy ? null : () => _tap(code(n), kind),
            borderRadius: BorderRadius.circular(Radii.md),
            child: Container(
              padding: const EdgeInsets.all(4),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(Radii.md),
                border: Border.all(
                  color: code(n) == p.avatar || code(n) == p.background ? Palette.accent : Palette.borderSubtle,
                  width: 2,
                ),
              ),
              child: Stack(
                fit: StackFit.expand,
                children: [
                  Center(child: tile(code(n))),
                  if (!p.usable(code(n)))
                    Align(
                      alignment: Alignment.bottomCenter,
                      child: ColoredBox(
                        color: Palette.bgPrimary.withValues(alpha: 0.7),
                        child: Text(l.shopPrice(price), style: const TextStyle(fontSize: 11, color: Palette.warning)),
                      ),
                    ),
                ],
              ),
            ),
          ),
      ],
    );
  }
}

/// Three boost boxes side by side (DD-320): watch an ad or spend credits
/// for two hours of +100% XP in one category.
class BoostBar extends StatefulWidget {
  const BoostBar({super.key, required this.repo});

  final ShopRepository repo;

  @override
  State<BoostBar> createState() => _BoostBarState();
}

class _BoostBarState extends State<BoostBar> {
  List<BoostState> _boosts = const [];
  bool _ads = false;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final (b, ads) = await widget.repo.boosts();
      if (mounted) {
        setState(() {
          _boosts = b;
          _ads = ads;
        });
      }
    } on ApiError {
      // optional row
    }
  }

  Future<void> _go(Future<void> Function() action) async {
    setState(() => _busy = true);
    try {
      await action();
      await _load();
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(AppLocalizations.of(context), e)));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_boosts.isEmpty) return const SizedBox.shrink();
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        for (final b in _boosts) ...[
          Expanded(
            child: Container(
              key: Key('boost-${b.category}'),
              padding: const EdgeInsets.all(Space.sm),
              decoration: BoxDecoration(
                color: b.active ? Palette.success.withValues(alpha: 0.12) : Palette.bgElevated,
                borderRadius: BorderRadius.circular(Radii.md),
                border: Border.all(color: b.active ? Palette.success : Palette.borderSubtle),
              ),
              child: Column(
                children: [
                  Text(boostName(l, b.category), style: text.labelMedium, textAlign: TextAlign.center),
                  const SizedBox(height: 2),
                  Text(
                    b.active ? l.boostLeft(shortDuration(l, b.endsAt!.difference(DateTime.now()))) : l.boostShort,
                    style: text.bodySmall?.copyWith(color: b.active ? Palette.success : null),
                    textAlign: TextAlign.center,
                  ),
                  if (!b.active) ...[
                    const SizedBox(height: Space.xs),
                    TextButton(
                      onPressed: _busy || !_ads || b.adsLeft == 0
                          ? null
                          : () => _go(() => widget.repo.devAd(b.category)),
                      child: Text(_ads ? l.boostWatchAd(b.adsLeft) : l.boostAdsSoon, textAlign: TextAlign.center),
                    ),
                  ],
                ],
              ),
            ),
          ),
          if (b != _boosts.last) const SizedBox(width: Space.sm),
        ],
      ],
    );
  }
}
