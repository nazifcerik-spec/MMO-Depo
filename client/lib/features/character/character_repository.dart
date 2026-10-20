import '../inventory/inventory_models.dart' show FoodBuff;
import '../../core/api_client.dart';
import '../progression/progression_repository.dart';
import 'models.dart';

/// One class skill (DD-222): actives auto-cast in combat, passives are stat
/// modifiers; both unlock by level.
class ClassSkill {
  const ClassSkill({
    required this.code,
    required this.active,
    required this.unlockLevel,
    this.effect = '',
    this.power = 0,
    this.mpCost = 0,
    this.cooldownTurns = 0,
    this.modifiers = const [],
    this.spec = '',
    this.durationTurns = 0,
    this.category = '',
  });

  factory ClassSkill.fromJson(Map<String, dynamic> j) => ClassSkill(
    code: j['code'] as String,
    spec: j['spec'] as String? ?? '',
    durationTurns: ((j['config'] as Map<String, dynamic>?)?['duration_turns'] as num?)?.toInt() ?? 0,
    category: (j['config'] as Map<String, dynamic>?)?['category'] as String? ?? '',
    active: j['kind'] == 'active',
    unlockLevel: (j['unlock_level'] as num?)?.toInt() ?? 1,
    effect: j['effect'] as String? ?? '',
    power: (j['power'] as num?)?.toInt() ?? 0,
    mpCost: (j['mp_cost'] as num?)?.toInt() ?? 0,
    cooldownTurns: (j['cooldown_turns'] as num?)?.toInt() ?? 0,
    modifiers: [
      for (final m in (j['modifiers'] as List? ?? const []).cast<Map<String, dynamic>>()) BonusEffect.fromJson(m),
    ],
  );

  final String code;
  final bool active;
  final int unlockLevel;
  final String effect; // strike | heal | guard | empower

  /// Owning class node (DD-250): the base class or a promotion.
  final String spec;
  final int durationTurns;

  /// skill | passive | awakening | ultimate (DD-250); empty on old servers.
  final String category;
  final int power;
  final int mpCost;
  final int cooldownTurns;
  final List<BonusEffect> modifiers;
}

/// Character endpoints. Stateless; screens own their loading state.
class CharacterRepository {
  const CharacterRepository(this.api);

  final ApiClient api;

  Future<CharacterList> list() async => CharacterList.fromJson(await api.get('/v1/character'));

  Future<CharacterSummary> create(String name, String race, String cls) async =>
      CharacterSummary.fromJson(await api.post('/v1/character', {'name': name, 'race_code': race, 'class_code': cls}));

  /// Marks the character as played and scopes later calls to it.
  Future<CharacterSheet> enter(String id) async {
    final body = await api.post('/v1/character/$id/enter', const {});
    api.characterId = id;
    return CharacterSheet.fromJson(body['sheet'] as Map<String, dynamic>);
  }

  /// Deletes the character at once (DD-283); the server checks [confirmName].
  Future<void> delete(String id, String confirmName) =>
      api.post('/v1/character/$id/delete', {'confirm_name': confirmName});

  /// Restores a character pending deletion by an account deletion.
  Future<void> cancelDelete(String id) => api.post('/v1/character/$id/delete-cancel', const {});

  Future<CharacterSheet> sheet(String id) async => CharacterSheet.fromJson(await api.get('/v1/character/$id/sheet'));

  /// The running food buff (DD-269), or null.
  Future<FoodBuff?> buff(String id) async => FoodBuff.fromActive((await api.get('/v1/character/$id/buff'))['buff']);

  static final _skills = <String, List<ClassSkill>>{};

  /// Class skill catalog (content; cached for the app session).
  Future<List<ClassSkill>> classSkills(String classCode) async => _skills[classCode] ??= [
    for (final s
        in ((await api.get('/v1/classes/$classCode/skills'))['skills'] as List? ?? const [])
            .cast<Map<String, dynamic>>())
      ClassSkill.fromJson(s),
  ];

  Future<int> goldCopper(String id) async {
    final w = await api.get('/v1/character/$id/wallet');
    return (w['gold_copper'] as num?)?.toInt() ?? 0;
  }
}
