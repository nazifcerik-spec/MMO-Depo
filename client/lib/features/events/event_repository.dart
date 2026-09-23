import '../../core/api_client.dart';

class RewardTier {
  const RewardTier(this.minScore, this.goldCopper, this.tokens);

  final int minScore;
  final int goldCopper;
  final int tokens;
}

/// One world boss / invasion spawn with this character's attendance.
class WorldEvent {
  const WorldEvent({
    required this.id,
    required this.code,
    required this.kind,
    required this.minLevel,
    required this.goal,
    required this.total,
    required this.startsAt,
    required this.endsAt,
    required this.state,
    required this.tiers,
    required this.enlisted,
    this.region,
    this.myPower,
    this.myScore,
    this.myTier,
  });

  factory WorldEvent.fromJson(Map<String, dynamic> j) => WorldEvent(
    id: j['id'] as String,
    code: j['template_code'] as String? ?? '',
    kind: j['kind'] as String? ?? '',
    minLevel: (j['min_level'] as num?)?.toInt() ?? 1,
    goal: (j['goal_score'] as num?)?.toInt() ?? 0,
    total: (j['total_score'] as num?)?.toInt() ?? 0,
    startsAt: DateTime.parse(j['starts_at'] as String),
    endsAt: DateTime.parse(j['ends_at'] as String),
    state: j['state'] as String? ?? '',
    tiers: [
      for (final t in (j['tiers'] as List? ?? const []).cast<Map<String, dynamic>>())
        RewardTier(
          (t['min_score'] as num?)?.toInt() ?? 0,
          (t['gold_copper'] as num?)?.toInt() ?? 0,
          (t['event_tokens'] as num?)?.toInt() ?? 0,
        ),
    ],
    enlisted: j['enlisted'] == true,
    region: j['region_code'] as String?,
    myPower: (j['my_power'] as num?)?.toInt(),
    myScore: (j['my_score'] as num?)?.toInt(),
    myTier: (j['my_reward_tier'] as num?)?.toInt(),
  );

  final String id;
  final String code;
  final String kind; // world_boss | invasion
  final int minLevel;
  final int goal;
  final int total;
  final DateTime startsAt;
  final DateTime endsAt;
  final String state; // scheduled | active | succeeded | failed
  final List<RewardTier> tiers;
  final bool enlisted;
  final String? region;
  final int? myPower;
  final int? myScore;
  final int? myTier;

  bool get open => state == 'scheduled' || state == 'active';
  double get progress => goal <= 0 ? 0 : (total / goal).clamp(0, 1).toDouble();
}

class SeasonalEvent {
  const SeasonalEvent({
    required this.code,
    required this.endsAt,
    required this.dailyTokens,
    required this.claimedToday,
    required this.balance,
  });

  factory SeasonalEvent.fromJson(Map<String, dynamic> j) => SeasonalEvent(
    code: j['code'] as String,
    endsAt: DateTime.parse(j['ends_at'] as String),
    dailyTokens: (j['daily_login_tokens'] as num?)?.toInt() ?? 0,
    claimedToday: j['claimed_today'] == true,
    balance: (j['event_token_balance'] as num?)?.toInt() ?? 0,
  );

  final String code;
  final DateTime endsAt;
  final int dailyTokens;
  final bool claimedToday;
  final int balance;
}

class EventRepository {
  const EventRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _c => '/v1/character/$characterId';

  Future<List<WorldEvent>> worldEvents() async => [
    for (final e in ((await api.get('$_c/world-events'))['events'] as List? ?? const []).cast<Map<String, dynamic>>())
      WorldEvent.fromJson(e),
  ];

  Future<void> enlist(String spawnId) => api.post('$_c/world-events/$spawnId/attend', const {});

  Future<List<SeasonalEvent>> seasonal() async => [
    for (final e
        in ((await api.get('/v1/events/seasonal'))['events'] as List? ?? const []).cast<Map<String, dynamic>>())
      SeasonalEvent.fromJson(e),
  ];

  Future<int> claimDaily(String code) async =>
      ((await api.post('/v1/events/seasonal/$code/claim', const {}))['event_tokens'] as num?)?.toInt() ?? 0;
}
