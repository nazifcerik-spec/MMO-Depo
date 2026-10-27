import '../../core/api_client.dart';
import '../inventory/inventory_models.dart';
import 'adventure_models.dart';

/// Hunt, dungeon and travel endpoints for one character. Catalogs are
/// cached for the life of the repository (they only change on content reload).
class AdventureRepository {
  AdventureRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;
  List<MobInfo>? _mobs;
  List<DungeonInfo>? _dungeons;
  WorldMap? _map;

  String get _char => '/v1/character/$characterId';

  CombatContract? _contract(Map<String, dynamic> j, String targetKey) {
    final c = j['contract'];
    return c is Map<String, dynamic> ? CombatContract.fromJson(c, targetKey) : null;
  }

  // --- hunt -----------------------------------------------------------------

  /// Hunting XP fall-off by level gap (DD-281), filled by [mobs].
  HuntXPFalloff huntXP = const HuntXPFalloff();

  /// Dungeon and raid XP by level gap (DD-303), filled by [mobs].
  DungeonXPScale dungeonXP = const DungeonXPScale();

  Future<List<MobInfo>> mobs() async {
    if (_mobs case final cached?) return cached;
    final j = await api.get('/v1/mobs');
    huntXP = HuntXPFalloff.fromJson(j['hunt_level_xp'] as Map<String, dynamic>?);
    dungeonXP = DungeonXPScale.fromJson(j['dungeon_level_xp'] as Map<String, dynamic>?);
    return _mobs = [for (final m in (j['mobs'] as List? ?? const []).cast<Map<String, dynamic>>()) MobInfo.fromJson(m)];
  }

  Future<CombatContract?> activeHunt() async => _contract(await api.get('$_char/hunt/active'), 'mob_code');

  /// Starts a hunt with [supplies] (template code -> units) taken from the
  /// bag (DD-288).
  Future<CombatContract> startHunt(String mobCode, Duration d, [Map<String, int> supplies = const {}]) async =>
      _contract(
        await api.post('$_char/hunt', {
          'mob_code': mobCode,
          'duration_seconds': d.inSeconds,
          'supplies': [
            for (final e in supplies.entries)
              if (e.value > 0) {'template_code': e.key, 'quantity': e.value},
          ],
        }),
        'mob_code',
      )!;

  /// The bag's potions and dishes, one entry per template (DD-288).
  Future<List<SupplyStock>> supplies() async {
    final bag = Bag.fromJson(await api.get('$_char/inventory'));
    final first = <String, BagItem>{};
    final stock = <String, int>{};
    for (final i in bag.items) {
      if (i.isLocked || i.supplyKind == null) continue;
      first.putIfAbsent(i.templateCode, () => i);
      stock[i.templateCode] = (stock[i.templateCode] ?? 0) + i.stack;
    }
    return [for (final e in first.entries) SupplyStock(item: e.value, stock: stock[e.key]!)];
  }

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

  // --- raids (DD-261) -------------------------------------------------------

  List<RaidInfo>? _raids;

  Future<List<RaidInfo>> raids() async => _raids ??= [
    for (final r in ((await api.get('/v1/raids'))['raids'] as List? ?? const []).cast<Map<String, dynamic>>())
      RaidInfo.fromJson(r),
  ];

  Future<RaidCurrent> raid() async => RaidCurrent.fromJson(await api.get('$_char/raid'));

  /// Party leader only: opens a lobby the party's members can join.
  Future<void> openRaid(String code) => api.post('$_char/raid', {'raid_code': code});

  Future<void> joinRaid(String id) => api.post('$_char/raid/$id/join', const {});

  /// The leader leaving calls the lobby off.
  Future<void> leaveRaid(String id) => api.post('$_char/raid/$id/leave', const {});

  Future<void> startRaid(String id) => api.post('$_char/raid/$id/start', const {});

  // --- raid finding (DD-299) ---------------------------------------------------

  Future<void> queueRaid(String code) => api.post('$_char/raid/queue', {'raid_code': code});
  Future<void> leaveRaidQueue() => api.post('$_char/raid/queue/leave', const {});
  Future<void> advertiseRaid(String code) => api.post('$_char/raid/listing', {'raid_code': code});
  Future<void> unlistRaid() => api.post('$_char/raid/listing/close', const {});
  Future<void> joinRaidListing(String id) => api.post('$_char/raid/listings/$id/join', const {});

  Future<List<RaidListing>> raidListings() async => [
    for (final r
        in ((await api.get('$_char/raid/listings'))['listings'] as List? ?? const []).cast<Map<String, dynamic>>())
      RaidListing.fromJson(r),
  ];

  /// Only once the raid has ended (E_RAID_NOT_READY); idempotent.
  Future<CombatResult> revealRaid(String id) async => raidResult(await api.post('$_char/raid/$id/reveal', const {}));

  // --- travel ---------------------------------------------------------------

  Future<Location> location() async => Location.fromJson(await api.get('$_char/location'));

  /// Regions with map positions and roads (DD-279); cached.
  Future<WorldMap> worldMap() async => _map ??= WorldMap.fromJson(await api.get('/v1/world/regions'));

  Future<void> travel(String toRegion) => api.post('$_char/travel', {'to_region_code': toRegion});

  /// The worker also completes arrivals; E_TRAVEL_NOT_ACTIVE then just means
  /// "already arrived" and the caller reloads.
  Future<void> arrive(String travelId) => api.post('/v1/travel/$travelId/complete', const {});

  Future<void> cancelTravel(String travelId) => api.post('/v1/travel/$travelId/cancel', const {});
}
