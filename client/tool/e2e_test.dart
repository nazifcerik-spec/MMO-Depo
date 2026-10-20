// End-to-end contract test: the client's own repositories and models against
// a real API (docker compose dev stack). Not part of `flutter test` (lives
// outside test/). Run:
//
//   cd deploy && docker compose up -d
//   cd client && flutter test tool/e2e_test.dart --dart-define=E2E_API=http://localhost:8080
//
// Creates two throwaway accounts per run on the dev database.
import 'dart:math';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:otter/core/api_client.dart';
import 'package:otter/core/error_reporter.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/activity/activity_repository.dart';
import 'package:otter/features/adventure/adventure_repository.dart';
import 'package:otter/features/character/character_repository.dart';
import 'package:otter/features/character/class_catalog.dart';
import 'package:otter/features/character/leveling.dart';
import 'package:otter/features/character/models.dart';
import 'package:otter/features/events/event_repository.dart';
import 'package:otter/features/inventory/inventory_repository.dart';
import 'package:otter/features/market/market_repository.dart';
import 'package:otter/features/notifications/notification_repository.dart';
import 'package:otter/features/progression/progression_repository.dart';
import 'package:otter/features/quests/quest_repository.dart';
import 'package:otter/features/social/group_repository.dart';
import 'package:otter/features/social/social_repository.dart';
import 'package:otter/features/support/support_repository.dart';
import 'package:otter/features/workshop/workshop_repository.dart';
import 'package:otter/shared/content_names.dart';

const base = String.fromEnvironment('E2E_API');

Matcher apiError(String code) => throwsA(isA<ApiError>().having((e) => e.code, 'code', code));

class Player {
  Player(this.api, this.sheet);

  final ApiClient api;
  final CharacterSheet sheet;
  String get id => sheet.id;
  String get name => sheet.name;
}

Future<Player> newPlayer(String tag, String cls) async {
  final suffix = Random().nextInt(1 << 30).toRadixString(36);
  final api = ApiClient(baseUrl: base, store: MemoryTokenStore());
  await api.register('e2e_$tag$suffix@example.com', 'e2e_$tag$suffix', 'e2e-password-$suffix', acceptTerms: true);
  final chars = CharacterRepository(api);
  final name = 'E${tag.toUpperCase()}$suffix';
  final created = await chars.create(name.substring(0, min(16, name.length)), 'human', cls);
  final sheet = await chars.enter(created.id);
  return Player(api, sheet);
}

Future<void> mint(Player p, int copper) =>
    p.api.post('/v1/dev/mint', {'character_id': p.id, 'amount': copper, 'currency': 'gold'});

void main() {
  if (base.isEmpty) {
    test('e2e skipped (no --dart-define=E2E_API)', () {}, skip: true);
    return;
  }
  late Player a, b;

  setUpAll(() async {
    a = await newPlayer('a', 'warrior');
    b = await newPlayer('b', 'mage');
  });

  test('sheet, wallet, class skills', () async {
    final chars = CharacterRepository(a.api);
    final sheet = await chars.sheet(a.id);
    expect(sheet.level, 1);
    expect(sheet.resources['hp_max'], greaterThan(0));
    expect(await chars.goldCopper(a.id), 0);
    final skills = await chars.classSkills('warrior');
    expect(skills.where((s) => s.active && s.unlockLevel == 1), isNotEmpty);
    expect(skills.where((s) => !s.active).first.modifiers, isNotEmpty);
    expect(skills.map((s) => s.spec), containsAll(['warrior', 'guardian', 'iron_fortress']), reason: 'DD-250 tree');

    // Creation screen catalog (DD-253): 10 + 20 + 40 nodes, real stats.
    final cat = await ClassCatalog.load(a.api);
    expect(cat.nodes, hasLength(70));
    expect(cat.node('warrior')!.resourceKind, 'rage');
    expect(cat.node('shaman')!.baseStats['wisdom'], 11);
    expect(cat.children('ranger').map((n) => n.code), ['marksman', 'scout']);
    expect(cat.attributePointsTotal, 525);
    expect(skills.where((s) => s.category == 'ultimate'), hasLength(4), reason: 'one per specialization');

    // Leveling (DD-250): level 1 of 1000, Rage, first objective, no points yet.
    final lv = await LevelingRepository(a.api, a.id).summary();
    expect(lv.level, 1);
    expect(lv.levelCap, 1000);
    expect(lv.spec, 'warrior');
    expect(lv.next!.level, 10, reason: 'Guard unlocks at 10');
    expect(lv.attributePoints, 0);
    final board = await LevelingRepository(a.api, a.id).talents();
    expect(board.trees, hasLength(3));
    await expectLater(LevelingRepository(a.api, a.id).promote('guardian'), apiError('E_PROMOTION_UNAVAILABLE'));
  });

  test('gathering: start, active preview, early reveal, stop guard', () async {
    final repo = ActivityRepository(a.api);
    final skills = await repo.skills(a.id);
    expect(skills.map((g) => g.kind), containsAll(['mining', 'herbalism', 'woodcutting', 'fishing', 'hunting']));
    final mining = skills.firstWhere((g) => g.kind == 'mining');
    expect(mining.level, 1);
    expect(mining.best!.code, 'raw_stone');
    await expectLater(
      repo.start(a.id, 'mining', const Duration(minutes: 1), target: 'adamantite_ore'),
      apiError('E_GATHER_LEVEL'),
    );
    final started = await repo.start(a.id, 'mining', const Duration(minutes: 1), target: 'raw_stone');
    expect(started.contract.kind, 'mining');
    final active = await repo.active(a.id);
    expect(active!.contract.id, started.contract.id);
    expect(active.current.target, 'raw_stone');
    final (y, _) = await repo.reveal(started.contract.id);
    expect(y.goldCopper, greaterThanOrEqualTo(0));
    expect(y.items.every((i) => i.templateCode == 'raw_stone'), isTrue);
    expect(await repo.active(a.id), isNull);
  });

  test('quests + hunt: accept, kill ticks the objective, turn in pays', () async {
    final quests = QuestRepository(a.api, a.id);
    final board = await quests.board();
    expect(board.available.map((q) => q.code), contains('main_01_first_hunt'));
    expect(board.available.where((q) => q.category == 'daily'), hasLength(3));
    await quests.accept('main_01_first_hunt');
    await expectLater(quests.accept('main_01_first_hunt'), apiError('E_QUEST_TAKEN'));
    await expectLater(quests.turnIn('main_01_first_hunt'), apiError('E_QUEST_INCOMPLETE'));

    final adv = AdventureRepository(a.api, a.id);
    final mobs = await adv.mobs();
    expect(mobs.first.code, 'young_boar');
    // Fights run back to back (DD-281): a few seconds hold several kills.
    final started = await a.api.post('/v1/character/${a.id}/hunt', {'mob_code': 'young_boar', 'duration_seconds': 600});
    expect(started['contract'], isNotNull);
    final hunt = await adv.activeHunt();
    expect(hunt!.target, 'young_boar');
    await Future<void>.delayed(const Duration(seconds: 32));
    final res = await adv.revealHunt(hunt.id);
    expect(res.victories, 1, reason: 'a fresh warrior beats a young boar');
    expect(res.loot.xp, greaterThan(0));

    final after = await quests.board();
    final q = after.open.singleWhere((q) => q.code == 'main_01_first_hunt');
    expect(q.objectives.single.progress, 1, reason: 'mob_killed:young_boar ticked the quest');

    final progress = await ProgressionRepository(a.api, a.id).summary();
    expect(progress.mastery.singleWhere((m) => m.track == 'combat').xp, greaterThan(0));
    final achievements = await ProgressionRepository(a.api, a.id).achievements();
    expect(achievements.singleWhere((x) => x.code == 'first_blood').unlocked, isTrue);
  }, timeout: const Timeout(Duration(minutes: 2)));

  test('dungeon level gate, travel start/cancel', () async {
    final adv = AdventureRepository(a.api, a.id);
    final dungeons = await adv.dungeons();
    expect(dungeons, isNotEmpty);
    await expectLater(adv.startDungeon(dungeons.first.code), apiError('E_LEVEL_REQUIREMENT'));

    final loc = await adv.location();
    expect(loc.region.code, isNotEmpty);
    final open = loc.neighbors.where((n) => n.region.minLevel <= 1).toList();
    expect(open, isNotEmpty);
    await adv.travel(open.first.region.code);
    final moving = await adv.location();
    expect(moving.travel, isNotNull);
    // One action at a time (DD-263): travelling blocks a hunt.
    await expectLater(adv.startHunt('young_boar', const Duration(hours: 1)), apiError('E_CHARACTER_BUSY'));
    await adv.cancelTravel(moving.travel!.id);
    expect((await adv.location()).travel, isNull);
  });

  test('vendor buy/sell, bag, equip, use potion, professions + craft gate', () async {
    await mint(a, 500000);
    await ContentNames.load(a.api, 'tr');
    expect(ContentNames.of('item', 'sword_iron'), 'Demir Kılıç');
    final shop = WorkshopRepository(a.api, a.id);
    final stock = await shop.stock();
    final potion = stock.firstWhere((o) => o.templateCode == 'potion_hp_small');
    await shop.buy(potion);
    await shop.buy(potion); // one to drink, one to sell
    final chest = stock.firstWhere((o) => o.templateCode == 'chest_novice'); // tier-1 plate for the warrior (DD-271)
    await shop.buy(chest);

    final inv = InventoryRepository(a.api, a.id);
    final bag = await inv.bag();
    final potionItem = bag.items.firstWhere((i) => i.templateCode == 'potion_hp_small');
    expect(potionItem.healPct, 20); // a share of max HP (DD-274)
    // Earlier hunts may have left the character wounded: a potion heals,
    // or is refused at full HP.
    try {
      final (hp, max, _) = await inv.use(potionItem.instanceId);
      expect(hp, lessThanOrEqualTo(max));
    } on ApiError catch (e) {
      expect(e.code, 'E_HP_FULL');
    }

    final armor = bag.items.firstWhere((i) => i.templateCode == chest.templateCode);
    await inv.equip(armor.instanceId);
    final eq = await inv.equipment();
    expect(eq.values.map((e) => e.templateCode), contains(chest.templateCode));

    expect(await inv.sellMultiplier(), greaterThan(0));
    final sold = await inv.sell(potionItem.instanceId);
    expect(sold, greaterThan(0));
    final buybacks = await shop.buybacks();
    expect(buybacks, isNotEmpty);
    await shop.buyBack(buybacks.first);

    expect(await shop.professions(), contains('alchemy'));
    await shop.learn('alchemy');
    expect((await shop.skills())['alchemy']!.level, 1);
    final recipe = (await shop.recipes('alchemy')).firstWhere((r) => r.code == 'recipe_potion_hp_small');
    await expectLater(shop.craft(recipe.code), apiError('E_INSUFFICIENT_MATERIALS'));
  });

  test('market: search, stats, list and cancel', () async {
    final market = MarketRepository(a.api, a.id);
    await market.search();
    final inv = InventoryRepository(a.api, a.id);
    final tradable = (await inv.bag()).items.firstWhere((i) => !i.isBound && !i.isLocked);
    final stats = await market.stats(tradable.templateCode);
    expect(stats.sales, greaterThanOrEqualTo(0));
    await market.sell(itemId: tradable.instanceId, quantity: 1, unitPrice: 5000, hours: 12);
    final mine = await market.mine();
    expect(mine.single.unitPrice, 5000);
    final forB = await MarketRepository(b.api, b.id).search();
    expect(forB.map((l) => l.id), contains(mine.single.id));
    await market.cancel(mine.single.id);
  });

  test('social: chat, friends, mail with gold', () async {
    final sa = SocialRepository(a.api, a.id), sb = SocialRepository(b.api, b.id);
    final msg = await sa.send('global', 'hello from e2e');
    expect(msg.fromName, a.name);
    expect((await sb.history('global')).map((m) => m.id), contains(msg.id));

    await sa.addFriend(b.name);
    final req = (await sb.requests()).single;
    expect(req.fromName, a.name);
    await sb.answer(req.id, accept: true);
    expect((await sa.friends()).map((f) => f.name), contains(b.name));

    await sa.sendMail(to: b.name, subject: 'gift', body: 'for you', gold: 1000);
    final inbox = await sb.inbox();
    final mail = inbox.firstWhere((m) => m.subject == 'gift');
    expect(mail.gold, 1000);
    final read = await sb.read(mail.id);
    expect(read.body, 'for you');
    await sb.claim(mail.id);
    expect(await WorkshopRepository(b.api, b.id).gold(), 1000);
  });

  test('party + guild between two players', () async {
    final ga = GroupRepository(a.api, a.id), gb = GroupRepository(b.api, b.id);
    await ga.createParty();
    await ga.inviteToParty(b.name);
    final inv = (await gb.partyInvites()).single;
    expect(inv.fromName, a.name);
    await gb.answerPartyInvite(inv.id, accept: true);
    expect((await ga.party())!.members, hasLength(2));
    await gb.leaveParty();

    await ga.createGuild('E2E Guild ${Random().nextInt(99999)}', 'E${Random().nextInt(9999)}');
    final guild = (await ga.guild())!;
    expect(guild.myRank, GuildRanks.leader);
    await ga.inviteToGuild(b.name);
    final ginv = (await gb.guildInvites()).single;
    await gb.answerGuildInvite(ginv.id, accept: true);
    final g2 = (await ga.guild())!;
    final member = g2.members.firstWhere((m) => m.name == b.name);
    await ga.setGuildRank(member.id, GuildRanks.officer);
    await ga.setMotd('raid tonight');
    expect((await gb.guild())!.motd, 'raid tonight');
    await ga.disbandGuild();
    expect((await ga.guild())!.disbandAt, isNotNull);
    await ga.cancelDisband();
  });

  test('notifications, progression titles, events', () async {
    final n = NotificationRepository(b.api);
    final (items, _) = await n.inbox();
    expect(items.map((i) => i.kind), containsAll(['mail_received', 'party_invite', 'guild_invite']));
    expect(items.map((i) => i.kind), isNot(contains('reward_collected')), reason: 'reward history is its own list');
    // Player a gathered and hunted earlier; b did neither.
    final (rewards, _) = await NotificationRepository(a.api).inbox(rewards: true);
    expect(rewards, isNotEmpty, reason: 'the gathering and hunt collections earlier were recorded (DD-242)');
    expect((await n.inbox(rewards: true)).$1, isEmpty, reason: 'history is per account');
    expect(rewards.every((r) => r.kind == 'reward_collected' && !r.unread), isTrue);
    await n.markAllRead();
    expect(await n.unreadCount(), 0);
    final prefs = await n.prefs();
    final updated = await n.setMode(prefs, 'market_sold', 'inbox_only');
    expect(updated.kinds['market_sold'], 'inbox_only');
    await n.mute(24);
    expect((await n.prefs()).muted, isTrue);
    await n.unmute();

    final prog = ProgressionRepository(a.api, a.id);
    expect(await prog.titles(), isEmpty);
    await expectLater(prog.setActiveTitle('bossbane'), apiError('E_TITLE_NOT_OWNED'));

    final ev = EventRepository(a.api, a.id);
    await ev.worldEvents();
    final seasonal = await ev.seasonal();
    if (seasonal.isNotEmpty && !seasonal.first.claimedToday) {
      expect(await ev.claimDaily(seasonal.first.code), seasonal.first.dailyTokens);
      await expectLater(ev.claimDaily(seasonal.first.code), apiError('E_ALREADY_CLAIMED'));
    }
  });

  test('support ticket: file, list, validation', () async {
    final s = SupportRepository(a.api);
    await s.create(category: 'bug', subject: 'E2E ticket', body: 'Filed from the e2e suite.', locale: 'tr');
    final mine = await s.mine();
    expect(mine.first.subject, 'E2E ticket');
    expect(mine.first.state, 'open');
    await expectLater(s.create(category: 'bug', subject: 'x', body: 'short', locale: 'tr'), apiError('E_VALIDATION'));
  });

  test('client error report reaches the collector', () async {
    // The reporter swallows failures, so check the status the server gave.
    final client = _StatusRecorder();
    await ErrorReporter(baseUrl: base, client: client).report(StateError('e2e'), StackTrace.current, kind: 'platform');
    expect(client.codes, [202]);
  });
}

class _StatusRecorder extends http.BaseClient {
  final _inner = http.Client();
  final codes = <int>[];

  @override
  Future<http.StreamedResponse> send(http.BaseRequest request) async {
    final res = await _inner.send(request);
    codes.add(res.statusCode);
    return res;
  }
}
