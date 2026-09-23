import '../../core/api_client.dart';

class GroupMember {
  const GroupMember({required this.id, required this.name, required this.level, required this.classCode, this.rank});

  factory GroupMember.fromJson(Map<String, dynamic> j) => GroupMember(
    id: j['character_id'] as String,
    name: j['name'] as String? ?? '',
    level: (j['level'] as num?)?.toInt() ?? 1,
    classCode: j['class_code'] as String? ?? '',
    rank: (j['rank_index'] as num?)?.toInt(),
  );

  final String id;
  final String name;
  final int level;
  final String classCode;

  /// Guild rank index (0 leader … 4 initiate); null for party members.
  final int? rank;
}

class Party {
  const Party({required this.id, required this.leaderId, required this.sizeCap, required this.members});

  factory Party.fromJson(Map<String, dynamic> j) => Party(
    id: j['id'] as String,
    leaderId: j['leader_character_id'] as String,
    sizeCap: (j['size_cap'] as num?)?.toInt() ?? 5,
    members: [
      for (final m in (j['members'] as List? ?? const []).cast<Map<String, dynamic>>()) GroupMember.fromJson(m),
    ],
  );

  final String id;
  final String leaderId;
  final int sizeCap;
  final List<GroupMember> members;
}

class Guild {
  const Guild({
    required this.id,
    required this.name,
    required this.tag,
    required this.leaderId,
    required this.motd,
    required this.level,
    required this.memberCount,
    required this.memberCap,
    required this.members,
    this.myRank,
    this.disbandAt,
  });

  factory Guild.fromJson(Map<String, dynamic> j) => Guild(
    id: j['id'] as String,
    name: j['name'] as String? ?? '',
    tag: j['tag'] as String? ?? '',
    leaderId: j['leader_character_id'] as String? ?? '',
    motd: j['motd'] as String? ?? '',
    level: (j['level'] as num?)?.toInt() ?? 1,
    memberCount: (j['member_count'] as num?)?.toInt() ?? 0,
    memberCap: (j['member_cap'] as num?)?.toInt() ?? 0,
    members: [for (final m in (j['members'] as List? ?? const []).cast<Map<String, dynamic>>()) GroupMember.fromJson(m)]
      ..sort((a, b) => (a.rank ?? 9).compareTo(b.rank ?? 9)),
    myRank: (j['my_rank'] as num?)?.toInt(),
    disbandAt: DateTime.tryParse(j['disband_scheduled_at'] as String? ?? ''),
  );

  final String id;
  final String name;
  final String tag;
  final String leaderId;
  final String motd;
  final int level;
  final int memberCount;
  final int memberCap;
  final List<GroupMember> members;
  final int? myRank;
  final DateTime? disbandAt;
}

class GroupInvite {
  const GroupInvite({required this.id, required this.fromName, this.guildName});

  final String id;
  final String fromName;

  /// "Name [TAG]" for guild invites; null for party invites.
  final String? guildName;
}

/// Guild creation fee (server guild.CreationFeeCopper, 10 gold).
const guildCreationFeeCopper = 100000;

/// Guild permission matrix — mirrors server guild.Can* (docs/guild/permissions.md §3).
abstract final class GuildRanks {
  static const leader = 0, officer = 1, veteran = 2, member = 3, initiate = 4;

  static bool canInvite(int rank) => rank <= veteran;
  static bool canSetMotd(int rank) => rank <= officer;
  static bool canKick(int actor, int target) => actor <= officer && target > actor;
  static bool canSetRank(int actor, int target, int to) =>
      to >= officer && to <= initiate && target > actor && (actor == leader || (actor == officer && to >= veteran));
}

class GroupRepository {
  const GroupRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _c => '/v1/character/$characterId';

  List<Map<String, dynamic>> _rows(Map<String, dynamic> j) =>
      (j['invites'] as List? ?? const []).cast<Map<String, dynamic>>();

  // --- party ----------------------------------------------------------------

  Future<Party?> party() async {
    final p = (await api.get('$_c/party'))['party'];
    return p is Map<String, dynamic> ? Party.fromJson(p) : null;
  }

  Future<List<GroupInvite>> partyInvites() async => [
    for (final i in _rows(await api.get('$_c/party/invites')))
      GroupInvite(id: i['id'] as String, fromName: i['invited_by_name'] as String? ?? ''),
  ];

  Future<void> createParty() => api.post('$_c/party', const {});
  Future<void> inviteToParty(String name) => api.post('$_c/party/invite', {'target_name': name});
  Future<void> answerPartyInvite(String id, {required bool accept}) =>
      api.post('$_c/party/invites/$id/${accept ? 'accept' : 'decline'}', const {});
  Future<void> leaveParty() => api.post('$_c/party/leave', const {});
  Future<void> kickFromParty(String id) => api.post('$_c/party/kick', {'target_character_id': id});
  Future<void> promoteInParty(String id) => api.post('$_c/party/promote', {'target_character_id': id});
  Future<void> disbandParty() => api.post('$_c/party/disband', const {});

  // --- guild ----------------------------------------------------------------

  Future<Guild?> guild() async {
    final g = (await api.get('$_c/guild'))['guild'];
    return g is Map<String, dynamic> ? Guild.fromJson(g) : null;
  }

  Future<List<GroupInvite>> guildInvites() async => [
    for (final i in _rows(await api.get('$_c/guild/invites')))
      GroupInvite(
        id: i['id'] as String,
        fromName: i['invited_by_name'] as String? ?? '',
        guildName: '${i['guild_name']} [${i['guild_tag']}]',
      ),
  ];

  Future<void> createGuild(String name, String tag) => api.post('$_c/guild', {'name': name, 'tag': tag});
  Future<void> inviteToGuild(String name) => api.post('$_c/guild/invite', {'target_name': name});
  Future<void> answerGuildInvite(String id, {required bool accept}) =>
      api.post('$_c/guild/invites/$id/${accept ? 'accept' : 'decline'}', const {});
  Future<void> leaveGuild() => api.post('$_c/guild/leave', const {});
  Future<void> kickFromGuild(String id) => api.post('$_c/guild/kick', {'target_character_id': id});
  Future<void> setGuildRank(String id, int rank) =>
      api.post('$_c/guild/rank', {'target_character_id': id, 'rank_index': rank});
  Future<void> transferGuild(String id) => api.post('$_c/guild/transfer', {'target_character_id': id});
  Future<void> setMotd(String motd) => api.post('$_c/guild/motd', {'motd': motd});
  Future<void> disbandGuild() => api.post('$_c/guild/disband', const {});
  Future<void> cancelDisband() => api.post('$_c/guild/disband/cancel', const {});
}
