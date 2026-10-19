import '../../core/api_client.dart';
import 'inventory_models.dart';

class InventoryRepository {
  const InventoryRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId';

  Future<Bag> bag() async => Bag.fromJson(await api.get('$_base/inventory'));

  Future<Map<String, EquippedItem>> equipment() async {
    final j = await api.get('$_base/equipment');
    final slots = j['slots'] as Map<String, dynamic>? ?? const {};
    return {for (final e in slots.entries) e.key: EquippedItem.fromJson(e.value as Map<String, dynamic>)};
  }

  /// Player→vendor rate, used only to preview a sale price.
  Future<double> sellMultiplier() async {
    final j = await api.get('/v1/vendor/$generalVendor');
    return (j['sell_multiplier'] as num?)?.toDouble() ?? 0;
  }

  Future<void> equip(String itemId) => api.post('$_base/equip/$itemId', const {});

  Future<void> unequip(String slot) => api.post('$_base/unequip/$slot', const {});

  /// Returns the copper actually received (server-computed).
  Future<int> sell(String itemId) async {
    final j = await api.post('/v1/vendor/$generalVendor/sell', {
      'item_instance_id': itemId,
      'character_id': characterId,
    });
    return (j['received_copper'] as num?)?.toInt() ?? 0;
  }

  /// Consumes one unit; returns HP after the effect and the food buff it
  /// started (DD-269).
  Future<(int, int, FoodBuff?)> use(String itemId) async {
    final j = await api.post('$_base/items/$itemId/use', const {});
    return ((j['hp_current'] as num).toInt(), (j['hp_max'] as num).toInt(), FoodBuff.fromActive(j['buff']));
  }

  Future<void> destroy(String itemId) => api.delete('$_base/inventory/$itemId');

  static Map<String, ItemSet>? _setsByPiece;

  /// piece template code → its set (DD-277); fetched once, they only change
  /// on a content reload.
  Future<Map<String, ItemSet>> setsByPiece() async => _setsByPiece ??= {
    for (final s in ((await api.get('/v1/item-sets'))['sets'] as List? ?? const []).cast<Map<String, dynamic>>())
      for (final p in ItemSet.fromJson(s).pieces) p: ItemSet.fromJson(s),
  };
}
