/// Wire models for bag + equipment (server/internal/inventory).
library;

/// Server slot allowlist (inventory.KnownEquipSlots), in paper-doll order.
const equipSlots = [
  'helm',
  'necklace',
  'chest',
  'gloves',
  'legs',
  'boots',
  'mainhand',
  'offhand',
  'ring1',
  'ring2',
  'trinket1',
  'trinket2',
];

/// Vendor that buys anything (content/vendors.yaml).
const generalVendor = 'general_ironhold';

class StatBonus {
  const StatBonus(this.kind, this.target, this.amount);

  final String kind;
  final String target;
  final double amount;
}

class BagItem {
  const BagItem({
    required this.instanceId,
    required this.templateCode,
    required this.itemClass,
    required this.itemSubclass,
    required this.slot,
    required this.stack,
    required this.maxStack,
    required this.rarity,
    required this.isBound,
    required this.isLocked,
    required this.weightGramsPer,
    required this.baseValueCopper,
    required this.minLevel,
    this.prefix,
    this.suffix,
    this.equipSlot,
    this.classRestriction,
    this.bonuses = const [],
    this.healHp = 0,
  });

  factory BagItem.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    final raw = j['stat_bonuses'];
    return BagItem(
      instanceId: j['instance_id'] as String,
      templateCode: j['template_code'] as String,
      itemClass: j['item_class'] as String? ?? '',
      itemSubclass: j['item_subclass'] as String? ?? '',
      slot: n('slot'),
      stack: n('stack_quantity'),
      maxStack: n('max_stack'),
      rarity: j['rarity'] as String? ?? 'common',
      isBound: j['is_bound'] as bool? ?? false,
      isLocked: j['is_locked'] as bool? ?? false,
      weightGramsPer: n('weight_grams_per'),
      baseValueCopper: n('base_value_copper'),
      minLevel: n('min_level'),
      prefix: j['prefix'] as String?,
      suffix: j['suffix'] as String?,
      equipSlot: j['equip_slot'] as String?,
      classRestriction: j['class_restriction'] as String?,
      healHp: switch (j['use_effect']) {
        {'kind': 'heal_hp', 'amount': final num a} => a.toInt(),
        _ => 0,
      },
      bonuses: [
        if (raw is List)
          for (final b in raw.whereType<Map<String, dynamic>>())
            StatBonus(b['kind'] as String? ?? '', b['target'] as String? ?? '', (b['amount'] as num?)?.toDouble() ?? 0),
      ],
    );
  }

  final String instanceId;
  final String templateCode;
  final String itemClass;
  final String itemSubclass;
  final int slot;
  final int stack;
  final int maxStack;
  final String rarity;
  final bool isBound;
  final bool isLocked;
  final int weightGramsPer;
  final int baseValueCopper;
  final int minLevel;
  final String? prefix;
  final String? suffix;

  /// Template slot this item equips into; null for materials/consumables.
  final String? equipSlot;
  final String? classRestriction;
  final List<StatBonus> bonuses;

  /// HP restored by using one unit (content use_effect heal_hp); 0 = not usable.
  final int healHp;

  bool get usable => healHp > 0;

  bool get equippable => equipSlot != null && equipSlot!.isNotEmpty;
}

class Bag {
  const Bag({
    required this.items,
    required this.slotCap,
    required this.slotsUsed,
    required this.weightGrams,
    required this.weightCapGrams,
  });

  factory Bag.fromJson(Map<String, dynamic> j) => Bag(
    items: [for (final i in (j['items'] as List? ?? const [])) BagItem.fromJson(i as Map<String, dynamic>)],
    slotCap: (j['slot_cap'] as num?)?.toInt() ?? 60,
    slotsUsed: (j['slots_used'] as num?)?.toInt() ?? 0,
    weightGrams: (j['weight_grams'] as num?)?.toInt() ?? 0,
    weightCapGrams: (j['weight_cap_grams'] as num?)?.toInt() ?? 0,
  );

  final List<BagItem> items;
  final int slotCap;
  final int slotsUsed;
  final int weightGrams;
  final int weightCapGrams;

  /// Slot index → item (bag cells are positional).
  Map<int, BagItem> get bySlot => {for (final i in items) i.slot: i};
}

/// One equipped item (the equipment endpoint returns a flattened row).
class EquippedItem {
  const EquippedItem({
    required this.instanceId,
    required this.templateCode,
    required this.rarity,
    required this.itemClass,
    this.prefix,
    this.suffix,
  });

  factory EquippedItem.fromJson(Map<String, dynamic> j) => EquippedItem(
    instanceId: j['instance_id'] as String,
    templateCode: j['template_code'] as String,
    rarity: j['rarity'] as String? ?? 'common',
    itemClass: j['item_class'] as String? ?? '',
    prefix: j['prefix'] as String?,
    suffix: j['suffix'] as String?,
  );

  final String instanceId;
  final String templateCode;
  final String rarity;
  final String itemClass;
  final String? prefix;
  final String? suffix;
}
