import '../../core/api_client.dart';
import '../inventory/inventory_models.dart';

class Recipe {
  const Recipe({
    required this.code,
    required this.profession,
    required this.minSkill,
    required this.inputs,
    required this.goldCost,
    required this.output,
    required this.outputStack,
    required this.xpReward,
    this.craftSeconds = 60,
    this.info,
  });

  factory Recipe.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    return Recipe(
      code: j['code'] as String,
      profession: j['profession_code'] as String,
      minSkill: n('min_skill'),
      inputs: {
        for (final i in (j['inputs'] as List? ?? const []).cast<Map<String, dynamic>>())
          i['template_code'] as String: (i['quantity'] as num).toInt(),
      },
      goldCost: n('gold_cost_copper'),
      output: j['output_template_code'] as String,
      outputStack: n('output_stack'),
      xpReward: n('profession_xp_reward'),
      craftSeconds: (j['craft_seconds'] as num?)?.toInt() ?? 60,
      info: j['output'] is Map<String, dynamic> ? ItemInfo.fromJson(j['output'] as Map<String, dynamic>) : null,
    );
  }

  final String code;
  final String profession;
  final int minSkill;

  /// template_code → quantity consumed.
  final Map<String, int> inputs;
  final int goldCost;
  final String output;
  final int outputStack;
  final int xpReward;

  /// Time per unit (DD-264).
  final int craftSeconds;

  /// What the recipe makes; null when the server sent no details.
  final ItemInfo? info;
}

/// A recipe output's properties (GET /v1/recipes `output`).
class ItemInfo {
  const ItemInfo({
    required this.itemClass,
    required this.itemSubclass,
    this.minLevel = 1,
    this.equipSlot,
    this.classRestriction,
    this.bonuses = const [],
    this.healHp = 0,
    this.weightGrams = 0,
    this.maxStack = 1,
  });

  factory ItemInfo.fromJson(Map<String, dynamic> j) {
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    final raw = j['stat_bonuses'];
    return ItemInfo(
      itemClass: j['item_class'] as String? ?? '',
      itemSubclass: j['item_subclass'] as String? ?? '',
      minLevel: n('min_level'),
      equipSlot: j['equip_slot'] as String?,
      classRestriction: j['class_restriction'] as String?,
      weightGrams: n('weight_grams'),
      maxStack: n('max_stack'),
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

  final String itemClass;
  final String itemSubclass;
  final int minLevel;
  final String? equipSlot;
  final String? classRestriction;
  final List<StatBonus> bonuses;
  final int healHp;
  final int weightGrams;
  final int maxStack;
}

/// A learned profession's progress.
class Skill {
  const Skill({required this.level, required this.xp, required this.xpToNext});

  final int level;
  final int xp;
  final int xpToNext;
}

class CraftResult {
  const CraftResult({
    required this.output,
    required this.stack,
    required this.skillAfter,
    required this.leveledUp,
    this.made = 0,
    this.refunded = 0,
  });

  final String output;
  final int stack;
  final int skillAfter;
  final bool leveledUp;

  /// Units delivered and units whose materials + gold came back (DD-264).
  final int made;
  final int refunded;
}

/// A timed crafting job (DD-264): [count] units, [unitSeconds] each.
class CraftJob {
  const CraftJob({
    required this.id,
    required this.recipe,
    required this.profession,
    required this.output,
    required this.count,
    required this.unitSeconds,
    required this.startedAt,
    required this.readyAt,
  });

  factory CraftJob.fromJson(Map<String, dynamic> j) => CraftJob(
    id: j['id'] as String,
    recipe: j['recipe_code'] as String? ?? '',
    profession: j['profession_code'] as String? ?? '',
    output: j['output_template_code'] as String? ?? '',
    count: (j['count'] as num?)?.toInt() ?? 1,
    unitSeconds: (j['unit_seconds'] as num?)?.toInt() ?? 60,
    startedAt: DateTime.parse(j['started_at'] as String),
    readyAt: DateTime.parse(j['ready_at'] as String),
  );

  final String id;
  final String recipe;
  final String profession;
  final String output;
  final int count;
  final int unitSeconds;
  final DateTime startedAt;
  final DateTime readyAt;

  /// Units finished at [now] (the server counts the same way).
  int done(DateTime now) {
    final n = now.difference(startedAt).inSeconds ~/ (unitSeconds <= 0 ? 1 : unitSeconds);
    return n.clamp(0, count);
  }
}

class Offer {
  const Offer({required this.templateCode, required this.itemClass, required this.itemSubclass, required this.price});

  final String templateCode;
  final String itemClass;
  final String itemSubclass;
  final int price;
}

class Buyback {
  const Buyback({
    required this.id,
    required this.vendor,
    required this.templateCode,
    required this.rarity,
    required this.stack,
    required this.price,
    required this.expiresAt,
    this.prefix,
    this.suffix,
  });

  final String id;
  final String vendor;
  final String templateCode;
  final String rarity;
  final int stack;
  final int price;
  final DateTime expiresAt;
  final String? prefix;
  final String? suffix;
}

/// Crafting + vendor endpoints for one character. Catalogs (professions,
/// recipes, vendor stock) are cached; they only change on content reload.
class WorkshopRepository {
  WorkshopRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;
  List<String>? _professions;
  List<Recipe>? _recipes;

  String get _char => '/v1/character/$characterId';

  Future<List<String>> professions() async => _professions ??= [
    for (final p
        in ((await api.get('/v1/professions'))['professions'] as List? ?? const []).cast<Map<String, dynamic>>())
      // Gathering professions level on the Gather tab (DD-244).
      if (p['kind'] != 'gathering') p['code'] as String,
  ];

  Future<List<Recipe>> recipes() async => _recipes ??= [
    for (final r in ((await api.get('/v1/recipes'))['recipes'] as List? ?? const []).cast<Map<String, dynamic>>())
      Recipe.fromJson(r),
  ];

  /// profession_code → progress, learned professions only.
  Future<Map<String, Skill>> skills() async {
    final list = (await api.get('$_char/professions'))['professions'] as List? ?? const [];
    return {
      for (final p in list.cast<Map<String, dynamic>>())
        p['profession_code'] as String: Skill(
          level: (p['level'] as num).toInt(),
          xp: (p['xp'] as num?)?.toInt() ?? 0,
          xpToNext: (p['xp_to_next'] as num?)?.toInt() ?? 0,
        ),
    };
  }

  /// template_code → total quantity across all bag stacks.
  Future<Map<String, int>> materials() async {
    final bag = Bag.fromJson(await api.get('$_char/inventory'));
    final out = <String, int>{};
    for (final i in bag.items.where((i) => !i.isLocked)) {
      out[i.templateCode] = (out[i.templateCode] ?? 0) + i.stack;
    }
    return out;
  }

  Future<int> gold() async => ((await api.get('$_char/wallet'))['gold_copper'] as num?)?.toInt() ?? 0;

  Future<void> learn(String profession) => api.post('$_char/professions/learn', {'profession_code': profession});

  /// Starts a timed job of [count] units (DD-264); the server checks and
  /// takes the whole job's materials and gold up front.
  Future<CraftJob> craft(String recipe, {int count = 1}) async => CraftJob.fromJson(
    (await api.post('$_char/craft', {'recipe_code': recipe, if (count > 1) 'count': count}))['job']
        as Map<String, dynamic>,
  );

  /// The running job, or null.
  Future<CraftJob?> activeCraft() async {
    final j = (await api.get('$_char/craft/active'))['job'];
    return j is Map<String, dynamic> ? CraftJob.fromJson(j) : null;
  }

  /// Ends the job: finished units delivered, the rest refunded (idempotent).
  Future<CraftResult> collectCraft(String id) async {
    final j = await api.post('$_char/craft/$id/collect', const {});
    return CraftResult(
      output: j['output_template_code'] as String? ?? '',
      stack: (j['output_stack'] as num?)?.toInt() ?? 0,
      skillAfter: (j['skill_after'] as num?)?.toInt() ?? 0,
      leveledUp: j['leveled_up'] as bool? ?? false,
      made: (j['made'] as num?)?.toInt() ?? 0,
      refunded: (j['refunded'] as num?)?.toInt() ?? 0,
    );
  }

  Future<List<Offer>> stock() async {
    final j = await api.get('/v1/vendor/$generalVendor');
    return [
      for (final o in (j['sells'] as List? ?? const []).cast<Map<String, dynamic>>())
        Offer(
          templateCode: o['template_code'] as String,
          itemClass: o['item_class'] as String? ?? '',
          itemSubclass: o['item_subclass'] as String? ?? '',
          price: (o['buy_price_copper'] as num).toInt(),
        ),
    ];
  }

  /// The quoted price is sent back so a price change between view and buy
  /// is rejected by the server instead of silently charging more.
  Future<void> buy(Offer o) => api.post('/v1/vendor/$generalVendor/buy', {
    'template_code': o.templateCode,
    'expected_price_copper': o.price,
    'character_id': characterId,
  });

  Future<List<Buyback>> buybacks() async {
    final list = (await api.get('$_char/buybacks'))['buybacks'] as List? ?? const [];
    return [
      for (final b in list.cast<Map<String, dynamic>>())
        Buyback(
          id: b['id'] as String,
          vendor: b['vendor_code'] as String,
          templateCode: b['template_code'] as String? ?? '',
          rarity: b['rarity'] as String? ?? 'common',
          stack: (b['stack_quantity'] as num?)?.toInt() ?? 1,
          price: (b['sold_price_copper'] as num?)?.toInt() ?? 0,
          expiresAt: DateTime.parse(b['expires_at'] as String),
          prefix: b['prefix'] as String?,
          suffix: b['suffix'] as String?,
        ),
    ];
  }

  Future<void> buyBack(Buyback b) => api.post('/v1/vendor/${b.vendor}/buyback/${b.id}', {'character_id': characterId});
}
