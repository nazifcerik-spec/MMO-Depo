import '../../core/api_client.dart';

/// What the character is doing now (DD-370). A raid lobby has no times yet.
class CurrentAction {
  const CurrentAction({required this.kind, required this.detail, this.startedAt, this.endsAt});

  factory CurrentAction.fromJson(Map<String, dynamic> j) => CurrentAction(
    kind: j['kind'] as String? ?? '',
    detail: j['detail'] as String? ?? '',
    startedAt: DateTime.tryParse(j['started_at'] as String? ?? '')?.toLocal(),
    endsAt: DateTime.tryParse(j['ends_at'] as String? ?? '')?.toLocal(),
  );

  final String kind;
  final String detail;
  final DateTime? startedAt;
  final DateTime? endsAt;
}

/// An item someone wears, as others see it.
class ShownGear {
  const ShownGear.fromJson(this.j);

  final Map<String, dynamic> j;

  String get templateCode => j['template_code'] as String? ?? '';
  String get itemClass => j['item_class'] as String? ?? '';
  String get itemSubclass => j['item_subclass'] as String? ?? '';
  String get rarity => j['rarity'] as String? ?? 'common';
  String? get prefix => j['prefix'] as String?;
  String? get suffix => j['suffix'] as String?;
}

/// A character's public page (DD-371).
class PlayerProfile {
  const PlayerProfile.fromJson(this.j);

  final Map<String, dynamic> j;

  String _s(String k) => j[k] as String? ?? '';
  int _n(String k) => (j[k] as num?)?.toInt() ?? 0;

  String get id => _s('character_id');
  String get name => _s('name');
  String get raceCode => _s('race_code');
  String get classCode => _s('class_code');

  /// The specialization, else the base class.
  String get pathCode => _s('spec_code').isEmpty ? classCode : _s('spec_code');
  int get level => _n('level');
  String get title => _s('title');
  String get avatar => j['avatar'] as String? ?? 'av01';
  String get background => j['background'] as String? ?? 'bg01';
  String get guildName => _s('guild_name');
  String get guildTag => _s('guild_tag');
  int get score => _n('score');

  /// 0 until the ranking has counted the character.
  int get rank => _n('rank');
  int get achievementPoints => _n('achievement_points');
  int get masteryTiers => _n('mastery_tiers');
  int get dungeonsCleared => _n('dungeons_cleared');
  int get raidsCleared => _n('raids_cleared');

  /// Profession code and level, highest first.
  List<(String, int)> get professions => [
    for (final p in (j['professions'] as List? ?? const []).cast<Map<String, dynamic>>())
      (p['code'] as String? ?? '', (p['level'] as num?)?.toInt() ?? 0),
  ];

  Map<String, ShownGear> get equipment => {
    for (final e in (j['equipment'] as Map<String, dynamic>? ?? const {}).entries)
      e.key: ShownGear.fromJson(e.value as Map<String, dynamic>),
  };

  /// Species and level of the active pet.
  (String, int)? get pet => switch (j['pet']) {
    {'species': final String s, 'level': final num l} => (s, l.toInt()),
    _ => null,
  };
}

class PlayerRepository {
  const PlayerRepository(this.api, this.characterId);

  final ApiClient api;

  /// The viewing character.
  final String characterId;

  /// Null when idle.
  Future<CurrentAction?> action() async {
    final a = (await api.get('/v1/character/$characterId/action'))['action'];
    return a is Map<String, dynamic> ? CurrentAction.fromJson(a) : null;
  }

  Future<PlayerProfile> profile(String name) async =>
      PlayerProfile.fromJson(await api.get('/v1/players/${Uri.encodeComponent(name)}'));
}
