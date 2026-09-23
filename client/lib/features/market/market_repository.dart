import '../../core/api_client.dart';

/// Allowed listing durations (hours) → fee multiplier (server: market.DurationMultiplier).
const listingDurations = {12: 1.0, 24: 1.5, 48: 2.5, 168: 5.0};

/// Client-side preview of the server's fee rule: max(10c, unit × qty × 5% × duration
/// multiplier). The server computes the real fee; this only sets expectations.
int listingFeePreview(int unit, int qty, int hours) {
  final fee = (unit * qty * 0.05 * (listingDurations[hours] ?? 1)).round();
  return fee < 10 ? 10 : fee;
}

/// 5% of every sale goes to the market (server: market.SaleCut).
int saleProceedsPreview(int total) => total - (total * 0.05).round();

class Listing {
  const Listing({
    required this.id,
    required this.sellerCharacterId,
    required this.sellerName,
    required this.templateCode,
    required this.itemClass,
    required this.rarity,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
    required this.expiresAt,
    required this.state,
    this.prefix,
    this.suffix,
  });

  factory Listing.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    return Listing(
      id: j['id'] as String,
      sellerCharacterId: j['seller_character_id'] as String? ?? '',
      sellerName: j['seller_name'] as String? ?? '',
      templateCode: j['template_code'] as String,
      itemClass: j['item_class'] as String? ?? '',
      rarity: j['rarity'] as String? ?? 'common',
      quantity: n('stack_quantity'),
      unitPrice: n('unit_price_copper'),
      totalPrice: n('total_price_copper'),
      expiresAt: DateTime.parse(j['expires_at'] as String),
      state: j['state'] as String? ?? 'active',
      prefix: j['prefix_code'] as String?,
      suffix: j['suffix_code'] as String?,
    );
  }

  final String id;
  final String sellerCharacterId;
  final String sellerName;
  final String templateCode;
  final String itemClass;
  final String rarity;
  final int quantity;
  final int unitPrice;
  final int totalPrice;
  final DateTime expiresAt;
  final String state;
  final String? prefix;
  final String? suffix;
}

/// 30-day price stats for one template (price suggestion when listing).
class PriceStats {
  const PriceStats({required this.sales, required this.avgUnit, this.lowestActive});

  final int sales;
  final int avgUnit;
  final int? lowestActive;

  /// Cheapest active listing, else the 30-day average, else nothing.
  int? get suggested => lowestActive ?? (sales > 0 ? avgUnit : null);
}

class MarketRepository {
  const MarketRepository(this.api, this.characterId);

  static const pageSize = 25;

  final ApiClient api;
  final String characterId;

  String get _char => '/v1/character/$characterId';

  List<Listing> _list(Map<String, dynamic> j) => [
    for (final l in (j['listings'] as List? ?? const []).cast<Map<String, dynamic>>()) Listing.fromJson(l),
  ];

  /// [itemClass]/[rarity] null = any. [sort]: price_asc | price_desc | newest.
  Future<List<Listing>> search({String? itemClass, String? rarity, String sort = 'price_asc', int offset = 0}) async {
    final q = Uri(
      queryParameters: {
        'class': ?itemClass,
        'rarity': ?rarity,
        'sort': sort,
        'limit': '$pageSize',
        'offset': '$offset',
      },
    ).query;
    return _list(await api.get('/v1/market/listings?$q'));
  }

  Future<List<Listing>> mine() async => _list(await api.get('$_char/market/listings'));

  Future<PriceStats> stats(String templateCode) async {
    final j = await api.get('/v1/market/history/${Uri.encodeComponent(templateCode)}');
    return PriceStats(
      sales: (j['sales_30d'] as num?)?.toInt() ?? 0,
      avgUnit: (j['avg_unit_price_copper'] as num?)?.toInt() ?? 0,
      lowestActive: (j['lowest_active_unit_price_copper'] as num?)?.toInt(),
    );
  }

  Future<void> sell({required String itemId, required int quantity, required int unitPrice, required int hours}) =>
      api.post('$_char/market/sell', {
        'item_instance_id': itemId,
        'stack_quantity': quantity,
        'unit_price_copper': unitPrice,
        'duration_hours': hours,
      });

  /// The unit price the buyer saw is echoed so a changed listing is refused.
  Future<void> buy(Listing l) => api.post('$_char/market/buy/${l.id}', {'expected_unit_price_copper': l.unitPrice});

  Future<void> cancel(String listingId) => api.post('$_char/market/cancel/$listingId', const {});
}
