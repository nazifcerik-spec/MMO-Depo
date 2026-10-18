import '../../core/api_client.dart';
import 'adventure_models.dart';

/// Hunt, dungeon and travel endpoints for one character. Catalogs are
/// cached for the life of the repository (they only change on content reload).
class AdventureRepository {
  AdventureRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;
  List<MobInfo>? _mobs;
  List<DungeonInfo>? _dungeons;

  String get _char => '/v1/character/$characterId';

  CombatContract? _contract(Map<String, dynamic> j, String targetKey) {
    final c = j['contract'];
    return c is Map<String, dynamic> ? CombatContract.fromJson(c, targetKey) : null;
  }

  // --- hunt -----------------------------------------------------------------

  /// XP fall-off for low-level mobs, filled by [mobs].
  List<MobXPBand> xpBands = const [];

  Future<List<MobInfo>> mobs() async {
    if (_mobs case final cached?) return cached;
    final j = await api.get('/v1/mobs');
    xpBands = [
      for (final b in (j['xp_bands'] as List? ?? const []).cast<Map<String, dynamic>>()) MobXPBand.fromJson(b),
    ];
    return _mobs = [for (final m in (j['mobs'] as List? ?? const []).cast<Map<String, dynamic>>()) MobInfo.fromJson(m)];
  }

  Future<CombatContract?> activeHunt() async => _contract(await api.get('$_char/hunt/active'), 'mob_code');

  Future<CombatContract> startHunt(String mobCode, Duration d) async =>
      _contract(await api.post('$_char/hunt', {'mob_code': mobCode, 'duration_seconds': d.inSeconds}), 'mob_code')!;

  /// Collects the encounters fought so far and ends the hunt (idempotent).
  Future<CombatResult> revealHunt(String id) async =>
      CombatResult.fromJson(await api.post('/v1/hunt/$id/reveal', const {}), 'hunt');

  Future<void> stopHunt(String id) => api.post('/v1/hunt/$id/stop', const {});

  // --- dungeon --------------------------------------------------------------

  Future<List<DungeonInfo>> dungeons() async => _dungeons ??= [
    for (final d
        in ((await api.get('/v1/dungeon/catalog'))['dungeons'] as List? ?? const []).cast<Map<String, dynamic>>())
      DungeonInfo.fromJson(d),
  ];

  Future<CombatContract?> activeDungeon() async => _contract(await api.get('$_char/dungeon/active'), 'dungeon_code');

  Future<CombatContract> startDungeon(String code) async =>
      _contract(await api.post('$_char/dungeon', {'dungeon_code': code}), 'dungeon_code')!;

  /// Only valid once the run's duration has elapsed (E_DUNGEON_NOT_READY).
  Future<CombatResult> revealDungeon(String id) async =>
      CombatResult.fromJson(await api.post('/v1/dungeon/$id/reveal', const {}), 'run');

  Future<void> stopDungeon(String id) => api.post('/v1/dungeon/$id/stop', const {});

  // --- travel ---------------------------------------------------------------

  Future<Location> location() async => Location.fromJson(await api.get('$_char/location'));

  Future<void> travel(String toRegion) => api.post('$_char/travel', {'to_region_code': toRegion});

  /// The worker also completes arrivals; E_TRAVEL_NOT_ACTIVE then just means
  /// "already arrived" and the caller reloads.
  Future<void> arrive(String travelId) => api.post('/v1/travel/$travelId/complete', const {});

  Future<void> cancelTravel(String travelId) => api.post('/v1/travel/$travelId/cancel', const {});
}
