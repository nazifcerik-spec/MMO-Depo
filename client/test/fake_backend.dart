import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:otter/core/api_client.dart';
import 'package:otter/core/token_store.dart';

/// Minimal in-memory Otter API for client tests. Records every request.
class FakeBackend {
  final requests = <http.Request>[];
  int refreshCalls = 0;
  int accessGen = 0;

  /// Access tokens the server currently accepts.
  final validAccess = <String>{};
  String? validRefresh;
  bool rejectRefresh = false;
  int failNextGets503 = 0;

  final characters = <Map<String, dynamic>>[];
  Map<String, dynamic>? activeActivity;
  final bag = <Map<String, dynamic>>[];
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
  };

  MockClient client() => MockClient((req) async {
    requests.add(req);
    final path = req.url.path;
    final body = req.body.isEmpty ? const {} : jsonDecode(req.body) as Map;

    if (path == '/v1/auth/login') {
      if (body['password'] != 'correct-horse') return error(401, 'E_INVALID_CREDENTIALS');
      return json(200, _session());
    }
    if (path == '/v1/auth/register') return json(201, _session());
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
    final so = _social(req, path, body);
    if (so != null) return so;
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
      final seed = bag.firstWhere((i) => i['template_code'] == 'linseed');
      seed['stack_quantity'] = (seed['stack_quantity'] as int) - 2;
      return json(200, {
        'output_template_code': 'potion_hp_small',
        'output_stack': 1,
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
    final del = RegExp(r'/inventory/([^/]+)$').firstMatch(path);
    if (del != null && req.method == 'DELETE') {
      bag.removeWhere((i) => i['instance_id'] == del[1]);
      return json(200, {'ok': true});
    }
    return null;
  }
}
