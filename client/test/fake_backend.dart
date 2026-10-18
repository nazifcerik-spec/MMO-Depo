import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:otter/core/api_client.dart';
import 'package:otter/core/token_store.dart';

/// Minimal in-memory Otter API for client tests. Records every request.
class FakeBackend {
  final requests = <http.Request>[];
  int exportLinks = 0;
  bool deletionRequested = false;
  int refreshCalls = 0;
  int accessGen = 0;

  /// Access tokens the server currently accepts.
  final validAccess = <String>{};
  String? validRefresh;
  bool rejectRefresh = false;
  int failNextGets503 = 0;

  final characters = <Map<String, dynamic>>[];
  bool emailVerified = true;
  final resetRequestedFor = <String>[];
  int verificationResends = 0;
  String? passwordChangedTo;
  Map<String, dynamic>? activeActivity;
  int gatherLevel = 20;

  // --- leveling (DD-250) ---
  /// null = endpoint missing (older server): the sheet falls back to its XP bar.
  Map<String, dynamic>? leveling;
  String? promotedTo;
  Map<String, dynamic>? lastAllocation;
  final learnedTalents = <String>[];
  int? lastCraftCount;
  final bag = <Map<String, dynamic>>[];
  int hp = 30; // current HP for consumable tests (max 120)
  Map<String, dynamic>? activeHunt;
  Map<String, dynamic>? activeDungeon;
  Map<String, dynamic>? travel;
  String region = 'ironhold_city';
  Duration dungeonLength = const Duration(minutes: 20);
  final equipped = <String, Map<String, dynamic>>{};
  int sheetGets = 0;
  int slotCap = 2;
  int goldCopper = 123456;

  http.Response json(int status, Object body) =>
      http.Response(jsonEncode(body), status, headers: {'content-type': 'application/json'});

  http.Response error(int status, String code, [String message = '']) => json(status, {
    'error': {'code': code, 'message': message},
  });

  Map<String, dynamic> _session() {
    accessGen++;
    final access = 'access-$accessGen';
    validAccess.add(access);
    validRefresh = 'refresh-$accessGen';
    return {
      'account_id': 'acct-1',
      'session_id': 'sess-1',
      'access_token': access,
      'access_expires_at': DateTime.now().toUtc().add(const Duration(minutes: 15)).toIso8601String(),
      'refresh_token': validRefresh,
      'refresh_expires_at': DateTime.now().toUtc().add(const Duration(days: 30)).toIso8601String(),
      'email_verified': emailVerified,
    };
  }

  Map<String, dynamic> sheetFor(Map<String, dynamic> c) => {
    ...c,
    'xp': 40,
    'xp_to_next': 100,
    'primary': {'str': 12, 'agi': 9, 'int': 5, 'vit': 11, 'cha': 4, 'luck': 3},
    'derived': {'attack': 18, 'defense': 7, 'crit_pct': 5, 'hit_pct': 90},
    'resources': {
      'hp_current': 80,
      'hp_max': 120,
      'mp_current': 20,
      'mp_max': 40,
      'stamina_current': 100,
      'stamina_max': 100,
      'energy_current': 50,
      'energy_max': 100,
    },
    if (c['class_code'] == 'ranger') ...{'resource_kind': 'focus', 'resource_max': 100},
  };

  MockClient client() => MockClient((req) async {
    requests.add(req);
    final path = req.url.path;
    final body = req.body.isEmpty ? const {} : jsonDecode(req.body) as Map;

    if (path == '/v1/auth/login') {
      if (body['password'] == 'busy-horse') return error(503, 'E_BUSY'); // argon2 slots full (DD-238)
      if (body['password'] != 'correct-horse') return error(401, 'E_INVALID_CREDENTIALS');
      return json(200, _session());
    }
    if (path == '/v1/auth/register') {
      if (body['accept_terms'] != true) return error(400, 'E_VALIDATION');
      return json(201, _session());
    }
    // Account rights (DD-231).
    if (path == '/v1/me/export-link') {
      exportLinks++;
      return json(201, {'url': 'http://localhost:8080/v1/me/export?token=t', 'expires_at': '2026-10-10T12:00:00Z'});
    }
    if (path == '/v1/me' && req.method == 'DELETE') {
      if (body['current_password'] != 'correct-horse') return error(401, 'E_INVALID_CREDENTIALS');
      deletionRequested = true;
      return json(202, {'purge_after': '2026-11-09T12:00:00Z'});
    }
    if (path == '/v1/auth/password/reset-request') {
      resetRequestedFor.add(body['email'] as String);
      return json(202, {'status': 'queued'});
    }
    if (path == '/v1/auth/refresh') {
      refreshCalls++;
      await Future<void>.delayed(Duration.zero);
      if (rejectRefresh || body['refresh_token'] != validRefresh) {
        return error(401, 'E_REFRESH_TOKEN_REUSE');
      }
      return json(200, _session());
    }

    final auth = req.headers['Authorization'];
    if (auth == null || !validAccess.contains(auth.substring('Bearer '.length))) {
      return error(401, 'E_TOKEN_EXPIRED');
    }
    if (req.method == 'GET' && failNextGets503 > 0) {
      failNextGets503--;
      return error(503, 'E_UNAVAILABLE');
    }
    if (path == '/v1/auth/logout') return json(200, {'ok': true});
    if (path == '/v1/auth/resend-verification') {
      verificationResends++;
      return json(202, {'status': 'queued'});
    }
    if (path == '/v1/auth/password/change') {
      if (body['current_password'] != 'correct-horse') return error(401, 'E_INVALID_CREDENTIALS');
      passwordChangedTo = body['new_password'] as String;
      return http.Response('', 204);
    }
    if (path == '/v1/character' && req.method == 'GET') {
      return json(200, {'characters': characters, 'slot_cap': slotCap});
    }
    if (path == '/v1/character' && req.method == 'POST') {
      final c = {
        'id': 'char-${characters.length + 1}',
        'name': body['name'],
        'race_code': body['race_code'],
        'class_code': body['class_code'],
        'level': 1,
        'state': 'active',
      };
      characters.add(c);
      return json(201, c);
    }
    final inv = _inventory(req, path, body);
    if (inv != null) return inv;
    final adv = _adventure(req, path, body);
    if (adv != null) return adv;
    final ws = _workshop(req, path, body);
    if (ws != null) return ws;
    final mk = _market(req, path, body);
    if (mk != null) return mk;
    final gr = _groups(req, path, body);
    if (gr != null) return gr;
    final lv = _leveling(req, path, body);
    if (lv != null) return lv;
    final meta = _meta(req, path, body);
    if (meta != null) return meta;
    final qs = _quests(req, path);
    if (qs != null) return qs;
    final so = _social(req, path, body);
    if (so != null) return so;
    if (RegExp(r'^/v1/character/[^/]+/gathering$').hasMatch(path)) {
      return json(200, {
        'professions': [
          for (final k in ['mining', 'herbalism', 'woodcutting', 'fishing', 'hunting'])
            {
              'kind': k,
              'level': k == 'fishing' ? gatherLevel : 1,
              'xp': 10,
              'xp_to_next': 64,
              'nodes': k == 'fishing'
                  ? [
                      {'code': 'tiny_guppy', 'level': 1, 'tier': 0},
                      {'code': 'salmon', 'level': 18, 'tier': 1},
                      {'code': 'bass', 'level': 20, 'tier': 2},
                      {'code': 'golden_carp', 'level': 70, 'tier': 5},
                    ]
                  : k == 'hunting'
                  ? [
                      {'code': 'forest_rat_pelt', 'level': 1, 'tier': 0},
                    ]
                  : [
                      {'code': 'raw_stone', 'level': 1, 'tier': 0},
                    ],
            },
        ],
      });
    }
    final act = RegExp(r'^/v1/character/([^/]+)/activity(/active)?$').firstMatch(path);
    if (act != null) {
      if (act[2] != null) {
        final a = activeActivity;
        return json(
          200,
          a == null
              ? {'active': null}
              : {
                  'contract': a,
                  'current_yield': {'gold_copper': 250, 'common_count': 3},
                },
        );
      }
      if (activeActivity != null) return error(409, 'E_SLOT_BUSY');
      final now = DateTime.now().toUtc();
      activeActivity = {
        'id': 'act-1',
        'character_id': act[1],
        'kind': body['kind'],
        'target': body['target'],
        'duration_seconds': body['duration_seconds'],
        'started_at': now.toIso8601String(),
        'expires_at': now.add(Duration(seconds: body['duration_seconds'] as int)).toIso8601String(),
        'state': 'active',
      };
      return json(201, {'contract': activeActivity, 'full_yield_peek': {}});
    }
    final done = RegExp(r'^/v1/activity/([^/]+)/(reveal|stop)$').firstMatch(path);
    if (done != null) {
      if (activeActivity == null) return error(409, 'E_NOT_ACTIVE');
      activeActivity = null;
      if (done[2] == 'stop') return json(200, {'id': done[1], 'state': 'cancelled'});
      return json(200, {
        'yield': {
          'gold_copper': 1520,
          'common_count': 2,
          'rare_count': 1,
          'items': [
            {'instance_id': 'i1', 'template_code': 'copper_ore', 'rarity': 'common'},
            {'instance_id': 'i2', 'template_code': 'raw_sapphire', 'rarity': 'rare', 'prefix': 'flawless'},
          ],
        },
        'replayed': false,
        'gather': {'level': gatherLevel + 1, 'xp': 4, 'xp_to_next': 70, 'gained': 36, 'leveled_up': true},
      });
    }
    final m = RegExp(r'^/v1/character/([^/]+)/(enter|sheet|wallet)$').firstMatch(path);
    if (m != null) {
      final c = characters.firstWhere((c) => c['id'] == m[1], orElse: () => const {});
      if (c.isEmpty) return error(404, 'E_NOT_FOUND');
      return switch (m[2]) {
        'enter' => json(200, {'character_id': c['id'], 'sheet': sheetFor(c)}),
        'sheet' => (() {
          sheetGets++;
          return json(200, sheetFor(c));
        })(),
        _ => json(200, {'character_id': c['id'], 'gold_copper': goldCopper}),
      };
    }
    return error(404, 'E_NOT_FOUND');
  });

  ApiClient api([TokenStore? store]) => ApiClient(
    baseUrl: 'http://otter.test',
    store: store ?? MemoryTokenStore(),
    client: client(),
    retryDelay: Duration.zero,
  );

  Map<String, dynamic> _contract(String id, String key, String target, Duration d) {
    final now = DateTime.now().toUtc();
    return {
      'id': id,
      key: target,
      'started_at': now.toIso8601String(),
      'expires_at': now.add(d).toIso8601String(),
      'state': 'active',
    };
  }

  http.Response? _adventure(http.Request req, String path, Map body) {
    if (path == '/v1/mobs') {
      return json(200, {
        'mobs': [
          {'code': 'giant_rat', 'level': 1, 'hp_max': 30, 'attack': 4, 'defense': 1, 'xp_reward': 8},
          {'code': 'grey_wolf', 'level': 3, 'hp_max': 60, 'attack': 9, 'defense': 3, 'xp_reward': 20},
          {'code': 'cave_troll', 'level': 20, 'hp_max': 900, 'attack': 70, 'defense': 40, 'xp_reward': 400},
        ],
      });
    }
    if (path == '/v1/dungeon/catalog') {
      return json(200, {
        'dungeons': [
          {
            'code': 'rat_warrens',
            'min_level': 3,
            'duration_seconds': dungeonLength.inSeconds,
            'encounters': 6,
            'boss_mob_code': 'rat_king',
            'bonus_gold_copper': 500,
            'bonus_xp': 120,
          },
          {
            'code': 'sunken_crypt',
            'min_level': 18,
            'duration_seconds': 3600,
            'encounters': 10,
            'boss_mob_code': 'lich',
          },
        ],
      });
    }
    if (path.endsWith('/hunt/active')) return json(200, {'contract': activeHunt});
    if (path.endsWith('/dungeon/active')) return json(200, {'contract': activeDungeon});
    if (RegExp(r'^/v1/character/[^/]+/hunt$').hasMatch(path)) {
      if (activeHunt != null) return error(409, 'E_HUNT_ACTIVE');
      activeHunt = _contract(
        'hunt-1',
        'mob_code',
        body['mob_code'] as String,
        Duration(seconds: body['duration_seconds'] as int),
      );
      return json(201, {'contract': activeHunt, 'planned_encounters_full': 30});
    }
    if (RegExp(r'^/v1/character/[^/]+/dungeon$').hasMatch(path)) {
      if (activeDungeon != null) return error(409, 'E_DUNGEON_ACTIVE');
      activeDungeon = _contract('dng-1', 'dungeon_code', body['dungeon_code'] as String, dungeonLength);
      return json(201, {'contract': activeDungeon, 'planned_encounters': 7});
    }
    if (path == '/v1/hunt/hunt-1/reveal') {
      activeHunt = null;
      return json(200, {
        'hunt': {'victories': 7, 'defeats': 1, 'total_xp': 160, 'total_gold_copper': 840, 'rare_drops': 1},
        'xp': {'level': 4, 'xp': 10, 'leveled_up': true},
      });
    }
    if (path == '/v1/dungeon/dng-1/reveal') {
      final c = activeDungeon!;
      if (DateTime.now().toUtc().isBefore(DateTime.parse(c['expires_at'] as String))) {
        return error(409, 'E_DUNGEON_NOT_READY');
      }
      activeDungeon = null;
      return json(200, {
        'run': {'cleared': true, 'boss_killed': true, 'victories': 7, 'total_xp': 300, 'total_gold_copper': 1500},
        'items': [
          {'instance_id': 'x', 'template_code': 'rat_crown', 'rarity': 'epic'},
        ],
      });
    }
    if (path == '/v1/hunt/hunt-1/stop') {
      activeHunt = null;
      return json(200, {'state': 'cancelled'});
    }
    if (path.endsWith('/location')) {
      return json(200, {
        'region': {'code': region, 'kind': 'city', 'safe': true, 'min_level': 1},
        'neighbors': [
          {'code': 'whisperwood', 'kind': 'wilderness', 'min_level': 1, 'travel_seconds': 300, 'via': 'child'},
          {'code': 'dragon_peaks', 'kind': 'zone', 'min_level': 40, 'travel_seconds': 1800, 'via': 'link'},
        ],
        'travel': ?travel,
      });
    }
    if (RegExp(r'^/v1/character/[^/]+/travel$').hasMatch(path)) {
      if (travel != null) return error(409, 'E_TRAVEL_IN_PROGRESS');
      final now = DateTime.now().toUtc();
      travel = {
        'id': 'trv-1',
        'from_region_code': region,
        'to_region_code': body['to_region_code'],
        'started_at': now.toIso8601String(),
        'arrives_at': now.add(const Duration(seconds: 300)).toIso8601String(),
        'state': 'in_progress',
      };
      return json(201, {'travel': travel});
    }
    if (path == '/v1/travel/trv-1/cancel') {
      travel = null;
      return json(200, {'ok': true});
    }
    return null;
  }

  /// profession → level (learned only).
  final skills = <String, int>{};
  final buybackRows = <Map<String, dynamic>>[];
  final bought = <String>[];

  http.Response? _workshop(http.Request req, String path, Map body) {
    if (path == '/v1/professions') {
      return json(200, {
        'professions': [
          {'code': 'alchemy', 'max_level': 50},
          {'code': 'blacksmithing', 'max_level': 50},
        ],
      });
    }
    if (path == '/v1/recipes') {
      return json(200, {
        'recipes': [
          {
            'code': 'recipe_potion_hp_small',
            'profession_code': 'alchemy',
            'min_skill': 1,
            'inputs': [
              {'template_code': 'linseed', 'quantity': 2},
            ],
            'gold_cost_copper': 5,
            'output_template_code': 'potion_hp_small',
            'output_stack': 1,
            'profession_xp_reward': 10,
            'output': {
              'item_class': 'consumable',
              'item_subclass': 'potion',
              'min_level': 2,
              'weight_grams': 100,
              'max_stack': 20,
              'use_effect': {'kind': 'heal_hp', 'amount': 60},
            },
          },
          {
            'code': 'recipe_potion_mp_small',
            'profession_code': 'alchemy',
            'min_skill': 3,
            'inputs': [
              {'template_code': 'mandrake', 'quantity': 2},
            ],
            'output_template_code': 'potion_mp_small',
            'output_stack': 1,
          },
        ],
      });
    }
    if (path.endsWith('/professions') && req.method == 'GET') {
      return json(200, {
        'professions': [
          for (final e in skills.entries) {'profession_code': e.key, 'level': e.value, 'xp': 0, 'xp_to_next': 100},
        ],
      });
    }
    if (path.endsWith('/professions/learn')) {
      skills[body['profession_code'] as String] = 1;
      return json(200, {'profession_code': body['profession_code'], 'level': 1});
    }
    if (path.endsWith('/craft')) {
      final n = body['count'] as int? ?? 1;
      lastCraftCount = n;
      final seed = bag.firstWhere((i) => i['template_code'] == 'linseed');
      seed['stack_quantity'] = (seed['stack_quantity'] as int) - 2 * n;
      return json(200, {
        'output_template_code': 'potion_hp_small',
        'output_stack': n,
        'count': n,
        'skill_before': 1,
        'skill_after': 2,
        'leveled_up': true,
      });
    }
    if (path == '/v1/vendor/general_ironhold/buy') {
      if (body['character_id'] != characters.first['id']) return error(404, 'E_NOT_FOUND');
      if (body['expected_price_copper'] != 50) return error(400, 'E_VALIDATION');
      bought.add(body['template_code'] as String);
      goldCopper -= 50;
      return json(200, {'template_code': body['template_code'], 'price_paid_copper': 50});
    }
    if (path.endsWith('/buybacks')) return json(200, {'buybacks': buybackRows});
    final bb = RegExp(r'^/v1/vendor/general_ironhold/buyback/([^/]+)$').firstMatch(path);
    if (bb != null) {
      buybackRows.removeWhere((b) => b['id'] == bb[1]);
      return json(200, {'price_paid_copper': 12});
    }
    return null;
  }

  Map<String, dynamic>? party;
  final partyInvites = <Map<String, dynamic>>[];
  Map<String, dynamic>? guild;
  final groupCalls = <String>[]; // "action body" log for assertions

  Map<String, dynamic> _member(String id, String name, [int? rank]) => {
    'character_id': id,
    'name': name,
    'level': 3,
    'class_code': 'rogue',
    'rank_index': ?rank,
  };

  http.Response? _groups(http.Request req, String path, Map body) {
    final m = RegExp(r'^/v1/character/[^/]+/(party|guild)(/.*)?$').firstMatch(path);
    if (m == null) return null;
    final kind = m[1]!, sub = m[2] ?? '';
    final me = characters.first['id'] as String;
    if (req.method == 'POST') groupCalls.add('$kind$sub ${jsonEncode(body)}');
    if (kind == 'party') {
      switch (sub) {
        case '' when req.method == 'GET':
          return json(200, {'party': party});
        case '/invites':
          return json(200, {'invites': partyInvites});
        case '':
          party = {
            'id': 'p-1',
            'leader_character_id': me,
            'size_cap': 5,
            'members': [_member(me, 'Grim')],
          };
          return json(201, {'party': party});
        case '/invite':
          return json(201, {'id': 'inv-9'});
        case '/promote':
          party!['leader_character_id'] = body['target_character_id'];
          return json(200, {});
        case '/leave' || '/disband':
          party = null;
          return json(200, {});
      }
      if (sub.startsWith('/invites/') && sub.endsWith('/accept')) {
        partyInvites.clear();
        party = {
          'id': 'p-2',
          'leader_character_id': 'char-5',
          'size_cap': 5,
          'members': [_member('char-5', 'Aria'), _member(me, 'Grim')],
        };
        return json(200, {'party': party});
      }
      return json(200, {});
    }
    switch (sub) {
      case '' when req.method == 'GET':
        return json(200, {'guild': guild});
      case '/invites':
        return json(200, {'invites': <Object>[]});
      case '':
        guild = {
          'id': 'g-1',
          'name': body['name'],
          'tag': body['tag'],
          'leader_character_id': me,
          'motd': '',
          'level': 1,
          'member_count': 2,
          'member_cap': 50,
          'my_rank': 0,
          'members': [_member(me, 'Grim', 0), _member('char-7', 'Brak', 3)],
        };
        return json(201, {'guild': guild, 'fee_copper': 100000});
      case '/motd':
        guild!['motd'] = body['motd'];
        return json(200, {'motd': body['motd']});
      case '/rank':
        for (final mem in (guild!['members'] as List).cast<Map<String, dynamic>>()) {
          if (mem['character_id'] == body['target_character_id']) mem['rank_index'] = body['rank_index'];
        }
        return json(200, {});
      case '/disband':
        guild!['disband_scheduled_at'] = DateTime.now().toUtc().add(const Duration(hours: 48)).toIso8601String();
        return json(200, {});
      case '/disband/cancel':
        guild!.remove('disband_scheduled_at');
        return json(200, {});
    }
    return json(200, {});
  }

  // --- quests ------------------------------------------------------------------

  final questsOpen = <Map<String, dynamic>>[];
  final questsAvailable = <Map<String, dynamic>>[
    {
      'code': 'main_01_first_hunt',
      'category': 'main',
      'state': 'available',
      'objectives': [
        {'event': 'mob_killed:young_boar', 'count': 5, 'progress': 0},
      ],
      'reward_xp': 80,
      'reward_gold_copper': 200,
    },
    {
      'code': 'daily_miner',
      'category': 'daily',
      'state': 'available',
      'objectives': [
        {'event': 'gather_mining', 'count': 40, 'progress': 0},
      ],
      'reward_xp': 100,
      'reward_gold_copper': 300,
    },
  ];
  final questCalls = <String>[];

  http.Response? _quests(http.Request req, String path) {
    final m = RegExp(r'^/v1/character/[^/]+/quests(?:/([^/]+)/(accept|abandon|turn-in))?$').firstMatch(path);
    if (m == null) return null;
    if (m[1] == null) {
      return json(200, {
        'open': questsOpen,
        'available': questsAvailable,
        'open_cap': 25,
        'daily_resets_at': DateTime.now().toUtc().add(const Duration(hours: 5)).toIso8601String(),
      });
    }
    final code = m[1]!;
    questCalls.add('${m[2]} $code');
    switch (m[2]) {
      case 'accept':
        final q = questsAvailable.firstWhere((q) => q['code'] == code);
        questsAvailable.remove(q);
        questsOpen.add({...q, 'state': 'active'});
        return json(200, questsOpen.last);
      case 'abandon':
        questsOpen.removeWhere((q) => q['code'] == code);
        return json(200, {'abandoned': true});
    }
    final q = questsOpen.firstWhere((q) => q['code'] == code);
    if (q['state'] != 'completed') return error(409, 'E_QUEST_INCOMPLETE');
    questsOpen.remove(q);
    return json(200, {
      'code': code,
      'reward_xp': q['reward_xp'],
      'reward_gold_copper': q['reward_gold_copper'],
      'xp': {'level': 4, 'xp': 10, 'leveled_up': true},
    });
  }

  /// Server-side tick stand-in: finishes an open quest's objectives.
  void completeQuest(String code) {
    final q = questsOpen.firstWhere((q) => q['code'] == code);
    q['state'] = 'completed';
    for (final o in (q['objectives'] as List).cast<Map<String, dynamic>>()) {
      o['progress'] = o['count'];
    }
  }

  // --- notifications / progression / events ---------------------------------

  final inbox = <Map<String, dynamic>>[];
  final rewardLog = <Map<String, dynamic>>[]; // reward history (DD-242)
  final notifKinds = <String, String>{'market_sold': 'push_and_inbox', 'mail_received': 'push_and_inbox'};
  Map<String, dynamic> quietHours = {'quiet_hours_start': '22:00', 'quiet_hours_end': '07:00', 'quiet_hours_tz': 'UTC'};
  String? muteUntil;
  Map<String, dynamic>? lastPrefsPut;
  String? activeTitle;
  final ownedTitles = <String>['prospector'];
  final worldEvents = <Map<String, dynamic>>[];
  bool seasonalClaimed = false;
  final tickets = <Map<String, dynamic>>[];
  Map<String, dynamic>? lastTicket;

  int get _unread => inbox.where((n) => n['read_at'] == null).length;

  Map<String, dynamic> _prefs() => {...quietHours, 'global_mute_until': muteUntil, 'kinds': notifKinds};

  http.Response? _meta(http.Request req, String path, Map body) {
    final now = DateTime.now().toUtc().toIso8601String();
    switch ((req.method, path)) {
      case ('GET', '/v1/support/tickets'):
        return json(200, {'tickets': tickets});
      case ('POST', '/v1/support/tickets'):
        lastTicket = Map<String, dynamic>.from(body);
        final t = {
          'id': 'tkt-${tickets.length + 1}',
          'category': body['category'],
          'subject': body['subject'],
          'body': body['body'],
          'state': 'open',
          'created_at': now,
        };
        tickets.insert(0, t);
        return json(201, t);
      case ('GET', '/v1/notifications'):
        final list = req.url.queryParameters['rewards'] == '1' ? rewardLog : inbox;
        return json(200, {'count': list.length, 'unread': _unread, 'notifications': list});
      case ('POST', '/v1/notifications/read-all'):
        for (final n in inbox) {
          n['read_at'] ??= now;
        }
        return json(200, {'marked': 1});
      case ('GET', '/v1/notification/preference'):
        return json(200, _prefs());
      case ('PUT', '/v1/notification/preference'):
        lastPrefsPut = Map<String, dynamic>.from(body);
        notifKinds.addAll((body['kinds'] as Map).cast<String, String>());
        return json(200, _prefs());
      case ('POST', '/v1/notification/preference/global-mute'):
        muteUntil = DateTime.now().toUtc().add(Duration(hours: body['hours'] as int)).toIso8601String();
        return json(200, {'global_mute_until': muteUntil});
      case ('DELETE', '/v1/notification/preference/global-mute'):
        muteUntil = null;
        return json(200, {'global_mute_until': null});
      case ('GET', '/v1/classes/rogue/skills'):
        return json(200, {
          'skills': [
            {
              'code': 'rogue_backstab',
              'kind': 'active',
              'unlock_level': 1,
              'effect': 'strike',
              'power': 200,
              'mp_cost': 12,
              'cooldown_turns': 3,
            },
            {
              'code': 'rogue_precision',
              'kind': 'passive',
              'unlock_level': 3,
              'modifiers': [
                {'kind': 'flat', 'target': 'crit_pct', 'amount': 5},
              ],
            },
            {
              'code': 'rogue_eviscerate',
              'kind': 'active',
              'unlock_level': 6,
              'effect': 'strike',
              'power': 300,
              'mp_cost': 25,
              'cooldown_turns': 5,
            },
          ],
        });
      case ('GET', '/v1/events/seasonal'):
        return json(200, {
          'events': [
            {
              'code': 'harvest_festival_2026',
              'starts_at': now,
              'ends_at': DateTime.now().toUtc().add(const Duration(days: 3)).toIso8601String(),
              'daily_login_tokens': 15,
              'claimed_today': seasonalClaimed,
              'event_token_balance': seasonalClaimed ? 15 : 0,
            },
          ],
        });
      case ('POST', '/v1/events/seasonal/harvest_festival_2026/claim'):
        if (seasonalClaimed) return error(409, 'E_ALREADY_CLAIMED');
        seasonalClaimed = true;
        return json(200, {'event_tokens': 15});
    }
    final del = RegExp(r'^/v1/notifications/([^/]+)$').firstMatch(path);
    if (del != null && req.method == 'DELETE') {
      inbox.removeWhere((n) => n['id'] == del[1]);
      rewardLog.removeWhere((n) => n['id'] == del[1]);
      return json(200, {'deleted': 1});
    }
    final read = RegExp(r'^/v1/notifications/([^/]+)/read$').firstMatch(path);
    if (read != null) {
      for (final n in inbox) {
        if (n['id'] == read[1]) n['read_at'] ??= now;
      }
      return json(200, {'marked': 1});
    }
    final m = RegExp(r'^/v1/character/[^/]+/(progression|achievements|titles|titles/active|world-events.*)$')
        .firstMatch(path);
    if (m == null) return null;
    switch (m[1]!) {
      case 'progression':
        return json(200, {
          'achievement_points': 5,
          'active_title_code': activeTitle,
          'mastery': [
            {
              'track': 'mining',
              'category': 'gathering',
              'xp': 150,
              'tier': 1,
              'max_tier': 10,
              'next_tier_xp': 500,
              'active_bonuses': [
                {'kind': 'flat', 'target': 'gather_yield_mining', 'amount': 1},
              ],
            },
          ],
          'reputation': [
            {'faction': 'ironhold_guard', 'rep_xp': 40, 'tier': 3, 'tier_name': 'neutral', 'daily_gain_left': 1460},
          ],
        });
      case 'achievements':
        return json(200, {
          'achievements': [
            {
              'code': 'slayer_100',
              'category': 'combat',
              'points': 10,
              'progress': 37,
              'target': 100,
              'reward_gold_copper': 5000,
            },
            {
              'code': 'first_blood',
              'category': 'combat',
              'points': 5,
              'progress': 1,
              'target': 1,
              'reward_gold_copper': 500,
              'unlocked_at': now,
            },
            {
              'code': 'miner_500',
              'category': 'gathering',
              'points': 10,
              'progress': 500,
              'target': 500,
              'reward_gold_copper': 5000,
              'reward_title_code': 'prospector',
              'unlocked_at': now,
            },
          ],
        });
      case 'titles':
        return json(200, {
          'titles': [
            for (final t in ownedTitles) {'code': t, 'unlocked_at': now, 'active': t == activeTitle},
          ],
        });
      case 'titles/active':
        final code = body['title_code'] as String;
        if (code.isNotEmpty && !ownedTitles.contains(code)) return error(409, 'E_TITLE_NOT_OWNED');
        activeTitle = code.isEmpty ? null : code;
        return json(200, {'active_title_code': code});
      case 'world-events':
        return json(200, {'events': worldEvents});
    }
    final attend = RegExp(r'world-events/([^/]+)/attend$').firstMatch(m[1]!);
    if (attend != null) {
      final e = worldEvents.firstWhere((e) => e['id'] == attend[1]);
      if (e['enlisted'] == true) return error(409, 'E_ALREADY_ENLISTED');
      e['enlisted'] = true;
      e['my_power'] = 420;
      return json(200, e);
    }
    return null;
  }

  final friends = <Map<String, dynamic>>[];
  final friendRequests = <Map<String, dynamic>>[];
  final blocked = <Map<String, dynamic>>[];
  final chat = <Map<String, dynamic>>[]; // oldest first
  final mail = <Map<String, dynamic>>[];
  final chatSent = <String>[];
  Map? lastMail;
  int _chatId = 100;

  http.Response? _social(http.Request req, String path, Map body) {
    final m = RegExp(r'^/v1/character/[^/]+/(.+)$').firstMatch(path);
    if (m == null) return null;
    final sub = m[1]!;
    final post = req.method == 'POST';
    if (sub == 'friends' && !post) return json(200, {'friends': friends});
    if (sub == 'friends/requests' && !post) return json(200, {'requests': friendRequests});
    if (sub == 'friends/request') {
      final name = body['target_name'];
      if (name == 'Nobody') return error(404, 'E_TARGET_NOT_FOUND');
      return json(200, {'friends': false, 'requested': true});
    }
    final ans = RegExp(r'^friends/requests/([^/]+)/(accept|decline)$').firstMatch(sub);
    if (ans != null) {
      final r = friendRequests.firstWhere((r) => r['id'] == ans[1]);
      friendRequests.remove(r);
      if (ans[2] == 'accept') {
        friends.add({'character_id': r['from_character_id'], 'name': r['from_name'], 'level': 4, 'class_code': 'mage'});
      }
      return json(200, {'accepted': ans[2] == 'accept'});
    }
    if (sub == 'friends/remove') {
      friends.removeWhere((f) => f['character_id'] == body['target_character_id']);
      return json(200, {});
    }
    if (sub == 'blocks' && !post) return json(200, {'blocked': blocked});
    final ch = RegExp(r'^chat/([a-z]+)$').firstMatch(sub);
    if (ch != null) {
      if (ch[1] == 'guild') return error(409, 'E_NOT_IN_CHANNEL');
      final target = req.url.queryParameters['target'];
      final rows = chat.where((c) => c['channel'] == ch[1] && (target == null || c['from_name'] == target)).toList();
      return json(200, {'messages': rows.reversed.toList()}); // server sends newest first
    }
    if (sub == 'chat') {
      chatSent.add(body['body'] as String);
      final msg = {
        'id': _chatId++,
        'channel': body['channel'],
        'from_character_id': characters.first['id'],
        'from_name': 'Grim',
        'body': body['body'],
        'sent_at': DateTime.now().toUtc().toIso8601String(),
      };
      return json(201, msg);
    }
    if (sub == 'mail' && !post) return json(200, {'mail': mail});
    if (sub == 'mail' && post) {
      lastMail = body;
      return json(201, {'mail_id': 'm-new', 'postage_copper': 7});
    }
    final one = RegExp(r'^mail/([^/]+)(/claim|/delete)?$').firstMatch(sub);
    if (one != null) {
      final row = mail.firstWhere((x) => x['id'] == one[1]);
      switch (one[2]) {
        case '/claim':
          row['state'] = 'claimed';
          return json(200, row);
        case '/delete':
          mail.remove(row);
          return json(200, {});
        default:
          if (row['state'] == 'sent') row['state'] = 'read';
          return json(200, {...row, 'body': 'Welcome to Ironhold!', 'attachments': []});
      }
    }
    return null;
  }

  final marketListings = <Map<String, dynamic>>[];
  final myListings = <Map<String, dynamic>>[];
  final marketBought = <String>[];
  Map? lastSell;

  http.Response? _market(http.Request req, String path, Map body) {
    if (path == '/v1/market/listings') {
      final cls = req.url.queryParameters['class'];
      return json(200, {
        'listings': [
          for (final l in marketListings)
            if (cls == null || l['item_class'] == cls) l,
        ],
      });
    }
    if (RegExp(r'^/v1/character/[^/]+/market/listings$').hasMatch(path)) return json(200, {'listings': myListings});
    if (path.startsWith('/v1/market/history/')) {
      return json(200, {'sales_30d': 3, 'avg_unit_price_copper': 40, 'lowest_active_unit_price_copper': 35});
    }
    if (path.endsWith('/market/sell')) {
      lastSell = body;
      final item = bag.firstWhere((i) => i['instance_id'] == body['item_instance_id']);
      myListings.add({
        'id': 'mine-1',
        'seller_character_id': characters.first['id'],
        'seller_name': 'Grim',
        'template_code': item['template_code'],
        'rarity': 'common',
        'stack_quantity': body['stack_quantity'],
        'unit_price_copper': body['unit_price_copper'],
        'total_price_copper': (body['unit_price_copper'] as int) * (body['stack_quantity'] as int),
        'expires_at': DateTime.now().toUtc().add(Duration(hours: body['duration_hours'] as int)).toIso8601String(),
        'state': 'active',
      });
      return json(201, myListings.last);
    }
    final buy = RegExp(r'/market/buy/([^/]+)$').firstMatch(path);
    if (buy != null) {
      final l = marketListings.firstWhere((l) => l['id'] == buy[1]);
      if (body['expected_unit_price_copper'] != l['unit_price_copper']) return error(409, 'E_PRICE_MISMATCH');
      marketListings.remove(l);
      marketBought.add(buy[1]!);
      return json(200, {...l, 'state': 'sold'});
    }
    final cancel = RegExp(r'/market/cancel/([^/]+)$').firstMatch(path);
    if (cancel != null) {
      myListings.removeWhere((l) => l['id'] == cancel[1]);
      return json(200, {'id': cancel[1], 'state': 'cancelled'});
    }
    return null;
  }

  int _freeSlot() {
    final used = {for (final i in bag) i['slot']};
    return List.generate(60, (i) => i).firstWhere((i) => !used.contains(i));
  }

  http.Response? _inventory(http.Request req, String path, Map body) {
    if (path.endsWith('/inventory') && req.method == 'GET') {
      return json(200, {
        'items': bag,
        'slot_cap': 60,
        'slots_used': bag.length,
        'weight_grams': 1200,
        'weight_cap_grams': 25000,
      });
    }
    if (path.endsWith('/equipment')) return json(200, {'slots': equipped});
    if (path == '/v1/vendor/general_ironhold') {
      return json(200, {
        'code': 'general_ironhold',
        'sell_multiplier': 0.25,
        'sells': [
          {'template_code': 'potion_hp_small', 'item_class': 'consumable', 'buy_price_copper': 50},
          {'template_code': 'sword_steel', 'item_class': 'weapon', 'buy_price_copper': 999999},
        ],
      });
    }
    final eq = RegExp(r'/equip/([^/]+)$').firstMatch(path);
    if (eq != null) {
      final item = bag.firstWhere((i) => i['instance_id'] == eq[1]);
      final restrict = item['class_restriction'];
      if (restrict != null && restrict != characters.first['class_code']) {
        return error(403, 'E_CLASS_RESTRICTION');
      }
      final slot = item['equip_slot'] as String;
      final prev = equipped.remove(slot);
      bag.remove(item);
      if (prev != null) bag.add({...prev, 'slot': item['slot']});
      equipped[slot] = item;
      return json(200, {'equipped_slot': slot});
    }
    final uneq = RegExp(r'/unequip/([^/]+)$').firstMatch(path);
    if (uneq != null) {
      final item = equipped.remove(uneq[1]);
      if (item == null) return error(404, 'E_NOTHING_EQUIPPED');
      bag.add({...item, 'slot': _freeSlot()});
      return json(200, {'ok': true});
    }
    if (path == '/v1/vendor/general_ironhold/sell') {
      final item = bag.firstWhere((i) => i['instance_id'] == body['item_instance_id']);
      bag.remove(item);
      return json(200, {
        'received_copper': ((item['base_value_copper'] as int) * 0.25 * (item['stack_quantity'] as int)).floor(),
      });
    }
    final use = RegExp(r'/items/([^/]+)/use$').firstMatch(path);
    if (use != null) {
      if (hp >= 120) return error(409, 'E_HP_FULL');
      final item = bag.firstWhere((i) => i['instance_id'] == use[1]);
      hp = (hp + 60).clamp(0, 120);
      final left = (item['stack_quantity'] as int) - 1;
      if (left == 0) {
        bag.remove(item);
      } else {
        item['stack_quantity'] = left;
      }
      return json(200, {
        'template_code': item['template_code'],
        'hp_current': hp,
        'hp_max': 120,
        'stack_remaining': left,
      });
    }
    final del = RegExp(r'/inventory/([^/]+)$').firstMatch(path);
    if (del != null && req.method == 'DELETE') {
      bag.removeWhere((i) => i['instance_id'] == del[1]);
      return json(200, {'ok': true});
    }
    return null;
  }

  http.Response? _leveling(http.Request req, String path, Map body) {
    final m = RegExp(r'^/v1/character/[^/]+/(leveling|promotion|attributes|talents(/[^/]+)?)$').firstMatch(path);
    if (m == null) return null;
    final lv = leveling;
    if (lv == null) return error(404, 'E_NOT_FOUND');
    switch ((req.method, m[1]!)) {
      case ('GET', 'leveling'):
        return json(200, lv);
      case ('POST', 'promotion'):
        promotedTo = body['choice'] as String;
        lv['spec_code'] = promotedTo;
        lv['path'] = [...(lv['path'] as List), promotedTo];
        lv['promotion'] = {'available': false, 'overflow_xp': 0};
        lv['level'] = 104;
        return json(200, {'spec_code': promotedTo, 'level': 104, 'levels_gained': 4});
      case ('POST', 'attributes'):
        lastAllocation = Map<String, dynamic>.from(body);
        final a = lv['attributes'] as Map<String, dynamic>;
        a['available'] = (a['available'] as int) - body.values.fold<int>(0, (s, v) => s + (v as int));
        return json(200, {'allocated': body});
      case ('GET', 'talents'):
        return json(200, {
          'earned': 2,
          'spent': learnedTalents.length,
          'available': 2 - learnedTalents.length,
          'tree_points': 25,
          'capstone_level': 850,
          'level': lv['level'],
          'trees': [
            {
              'code': 'warrior_defense',
              'points': learnedTalents.length,
              'talents': [
                {'code': 'warrior_defense_toughness', 'max_rank': 5, 'rank': learnedTalents.length},
                {'code': 'warrior_defense_living_fortress', 'max_rank': 1, 'rank': 0, 'capstone': true},
              ],
            },
          ],
        });
      case ('POST', final t) when t.startsWith('talents/'):
        if (learnedTalents.length >= 2) return error(409, 'E_NOT_ENOUGH_POINTS');
        learnedTalents.add(t.substring('talents/'.length));
        return json(200, {'rank': learnedTalents.length});
    }
    return error(404, 'E_NOT_FOUND');
  }
}
