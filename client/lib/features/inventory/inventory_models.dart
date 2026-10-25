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
    this.healPct = 0,
    this.manaPct = 0,
    this.buff,
  });

  factory BagItem.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    // What the item really adds — rarity and affixes applied (DD-275); the
    // template's bonuses when the server sent none.
    final raw = j['effective_bonuses'] ?? j['stat_bonuses'];
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
      healPct: switch (j['use_effect']) {
        {'kind': 'heal_pct', 'amount': final num a} => a.toInt(),
        _ => 0,
      },
      manaPct: switch (j['use_effect']) {
        {'kind': 'mana_pct', 'amount': final num a} => a.toInt(),
        _ => 0,
      },
      buff: FoodBuff.fromUseEffect(j['use_effect']),
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

  /// HP restored by using one unit (content use_effect heal_hp); 0 = none.
  final int healHp;

  /// Share of max HP restored by one unit (use_effect heal_pct, DD-274); 0 = none.
  final int healPct;

  /// Share of a hunter's mana a mana potion refills (use_effect mana_pct,
  /// DD-288): only usable as a hunt supply.
  final int manaPct;

  /// Food buff (DD-269); null for anything else.
  final FoodBuff? buff;

  bool get usable => healHp > 0 || healPct > 0 || buff != null;

  /// What it is on a hunt (DD-288): `hp`, `mp`, `food`, or null.
  String? get supplyKind => buff != null
      ? 'food'
      : manaPct > 0
      ? 'mp'
      : healHp > 0 || healPct > 0
      ? 'hp'
      : null;

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
/// A raid set (DD-277): its pieces and the bonuses at piece counts.
class ItemSet {
  const ItemSet({required this.code, required this.pieces, required this.bonuses});

  factory ItemSet.fromJson(Map<String, dynamic> j) => ItemSet(
    code: j['code'] as String,
    pieces: [for (final p in (j['pieces'] as List? ?? const [])) p as String],
    bonuses: [
      for (final b in (j['bonuses'] as List? ?? const []).cast<Map<String, dynamic>>())
        (
          (b['pieces'] as num).toInt(),
          [
            for (final m in (b['modifiers'] as List? ?? const []).cast<Map<String, dynamic>>())
              StatBonus(
                m['kind'] as String? ?? '',
                m['target'] as String? ?? '',
                (m['amount'] as num?)?.toDouble() ?? 0,
              ),
          ],
        ),
    ],
  );

  final String code;
  final List<String> pieces;

  /// (piece count, what wearing that many adds).
  final List<(int, List<StatBonus>)> bonuses;

  /// How many distinct pieces of the set are among [worn] template codes.
  int wornCount(Iterable<String> worn) => pieces.toSet().intersection(worn.toSet()).length;
}

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

/// Bag display orders. [slot] is the server's positional layout; the others
/// only reorder what is shown (positions on the server are unchanged).
enum BagSort { slot, rarity, type, name, value }

const _rarityOrder = ['common', 'uncommon', 'rare', 'epic', 'legendary', 'mythic'];

/// [items] in [sort] order; ties keep bag position so the view is stable.
List<BagItem> sortBag(List<BagItem> items, BagSort sort, String Function(BagItem) name) {
  int rarity(BagItem i) => _rarityOrder.indexOf(i.rarity);
  final cmp = switch (sort) {
    BagSort.slot => (BagItem a, BagItem b) => 0,
    BagSort.rarity => (BagItem a, BagItem b) => rarity(b) - rarity(a),
    BagSort.type => (BagItem a, BagItem b) {
      final c = a.itemClass.compareTo(b.itemClass);
      return c != 0 ? c : a.itemSubclass.compareTo(b.itemSubclass);
    },
    BagSort.name => (BagItem a, BagItem b) => name(a).toLowerCase().compareTo(name(b).toLowerCase()),
    BagSort.value => (BagItem a, BagItem b) => (b.baseValueCopper * b.stack).compareTo(a.baseValueCopper * a.stack),
  };
  return [...items]..sort((a, b) {
    final c = cmp(a, b);
    return c != 0 ? c : a.slot - b.slot;
  });
}

/// A food buff (DD-269): attack and defense percent for [minutes], or, when
/// read from the character, until [expiresAt].
class FoodBuff {
  const FoodBuff({required this.attackPct, required this.defensePct, this.minutes = 0, this.expiresAt});

  /// From a use_effect: null unless it carries a buff.
  static FoodBuff? fromUseEffect(Object? useEffect) => switch (useEffect) {
    {'buff': final Map<String, dynamic> b} => FoodBuff(
      attackPct: (b['attack_pct'] as num?)?.toInt() ?? 0,
      defensePct: (b['defense_pct'] as num?)?.toInt() ?? 0,
      minutes: (b['minutes'] as num?)?.toInt() ?? 0,
    ),
    _ => null,
  };

  /// The character's running buff (GET /buff, use result); null when none.
  static FoodBuff? fromActive(Object? j) => switch (j) {
    final Map<String, dynamic> b => FoodBuff(
      attackPct: (b['attack_pct'] as num?)?.toInt() ?? 0,
      defensePct: (b['defense_pct'] as num?)?.toInt() ?? 0,
      expiresAt: DateTime.tryParse(b['expires_at'] as String? ?? ''),
    ),
    _ => null,
  };

  final int attackPct;
  final int defensePct;
  final int minutes;
  final DateTime? expiresAt;
}
