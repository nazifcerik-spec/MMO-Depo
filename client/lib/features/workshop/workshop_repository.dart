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
}

/// A learned profession's progress.
class Skill {
  const Skill({required this.level, required this.xp, required this.xpToNext});

  final int level;
  final int xp;
  final int xpToNext;
}

class CraftResult {
  const CraftResult({required this.output, required this.stack, required this.skillAfter, required this.leveledUp});

  final String output;
  final int stack;
  final int skillAfter;
  final bool leveledUp;
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

  /// Crafts [count] at once (DD-245); the server checks the whole batch.
  Future<CraftResult> craft(String recipe, {int count = 1}) async {
    final j = await api.post('$_char/craft', {'recipe_code': recipe, if (count > 1) 'count': count});
    return CraftResult(
      output: j['output_template_code'] as String,
      stack: (j['output_stack'] as num?)?.toInt() ?? 1,
      skillAfter: (j['skill_after'] as num?)?.toInt() ?? 0,
      leveledUp: j['leveled_up'] as bool? ?? false,
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
