import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/inventory_models.dart';

BagItem _item(
  String id,
  int slot, {
  String rarity = 'common',
  String cls = 'material',
  String sub = 'ore',
  int value = 1,
  int stack = 1,
}) => BagItem.fromJson({
  'instance_id': id,
  'template_code': id,
  'item_class': cls,
  'item_subclass': sub,
  'slot': slot,
  'stack_quantity': stack,
  'rarity': rarity,
  'base_value_copper': value,
});

void main() {
  final items = [
    _item('c', 0, rarity: 'rare', cls: 'weapon', sub: 'sword', value: 50),
    _item('a', 1, value: 5, stack: 20),
    _item('b', 2, rarity: 'mythic', cls: 'consumable', sub: 'potion', value: 10),
    _item('d', 3, rarity: 'rare', value: 1),
  ];
  List<String> ids(BagSort s) => [for (final i in sortBag(items, s, (i) => i.templateCode)) i.instanceId];

  test('bag display orders; ties keep bag position', () {
    expect(ids(BagSort.slot), ['c', 'a', 'b', 'd']);
    expect(ids(BagSort.rarity), ['b', 'c', 'd', 'a']);
    expect(ids(BagSort.type), ['b', 'a', 'd', 'c']);
    expect(ids(BagSort.name), ['a', 'b', 'c', 'd']);
    expect(ids(BagSort.value), ['a', 'c', 'b', 'd'], reason: 'value × stack: 100, 50, 10, 1');
  });
}
