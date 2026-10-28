import '../../core/api_client.dart';

const chatChannels = ['global', 'guild', 'party'];

class ChatMessage {
  const ChatMessage({
    required this.id,
    required this.channel,
    required this.fromId,
    required this.fromName,
    required this.body,
    required this.sentAt,
    this.targetId,
    this.system = '',
  });

  factory ChatMessage.fromJson(Map<String, dynamic> j) => ChatMessage(
    id: (j['id'] as num?)?.toInt() ?? 0,
    channel: j['channel'] as String? ?? 'global',
    fromId: j['from_character_id'] as String? ?? '',
    fromName: j['from_name'] as String? ?? '',
    targetId: j['target_character_id'] as String?,
    body: j['body'] as String? ?? '',
    sentAt: DateTime.tryParse(j['sent_at'] as String? ?? '') ?? DateTime.now(),
    system: j['system'] as String? ?? '',
  );

  /// A line the game wrote (DD-312): "profanity_mute" = [fromName] was
  /// muted for swearing.
  final String system;

  final int id;
  final String channel;
  final String fromId;
  final String fromName;
  final String? targetId;
  final String body;
  final DateTime sentAt;
}

class Friend {
  const Friend({required this.id, required this.name, required this.level, required this.classCode});

  final String id;
  final String name;
  final int level;
  final String classCode;
}

class FriendRequest {
  const FriendRequest({required this.id, required this.fromName});

  final String id;
  final String fromName;
}

class Blocked {
  const Blocked({required this.id, required this.name});

  final String id;
  final String name;
}

class MailSummary {
  const MailSummary({
    required this.id,
    required this.senderName,
    required this.subject,
    required this.gold,
    required this.items,
    required this.state,
    required this.createdAt,
  });

  factory MailSummary.fromJson(Map<String, dynamic> j) => MailSummary(
    id: j['id'] as String,
    senderName: j['sender_name'] as String? ?? '',
    subject: j['subject'] as String? ?? '',
    gold: (j['gold_attachment_copper'] as num?)?.toInt() ?? 0,
    items: (j['item_attachment_count'] as num?)?.toInt() ?? 0,
    state: j['state'] as String? ?? 'sent',
    createdAt: DateTime.parse(j['created_at'] as String),
  );

  final String id;
  final String senderName;
  final String subject;
  final int gold;
  final int items;

  /// sent (unread) | read | claimed
  final String state;
  final DateTime createdAt;

  bool get unread => state == 'sent';
  bool get hasAttachments => state != 'claimed' && (gold > 0 || items > 0);
}

class MailAttachment {
  const MailAttachment({required this.templateCode, required this.rarity, required this.quantity});

  final String templateCode;
  final String rarity;
  final int quantity;
}

class Mail {
  const Mail(this.summary, this.body, this.attachments);

  final MailSummary summary;
  final String body;
  final List<MailAttachment> attachments;
}

class SocialRepository {
  const SocialRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _c => '/v1/character/$characterId';

  List<Map<String, dynamic>> _rows(Map<String, dynamic> j, String key) =>
      (j[key] as List? ?? const []).cast<Map<String, dynamic>>();

  // --- chat -----------------------------------------------------------------

  /// Oldest first. [peer] (character name) is required for whispers.
  Future<List<ChatMessage>> history(String channel, {String? peer}) async {
    final q = peer == null ? '' : '?target=${Uri.encodeQueryComponent(peer)}';
    final j = await api.get('$_c/chat/$channel$q');
    return [for (final m in _rows(j, 'messages')) ChatMessage.fromJson(m)].reversed.toList();
  }

  Future<ChatMessage> send(String channel, String body, {String? peer}) async =>
      ChatMessage.fromJson(await api.post('$_c/chat', {'channel': channel, 'body': body, 'target_name': ?peer}));

  // --- friends ----------------------------------------------------------------

  Future<List<Friend>> friends() async => [
    for (final f in _rows(await api.get('$_c/friends'), 'friends'))
      Friend(
        id: f['character_id'] as String,
        name: f['name'] as String,
        level: (f['level'] as num?)?.toInt() ?? 1,
        classCode: f['class_code'] as String? ?? '',
      ),
  ];

  Future<List<FriendRequest>> requests() async => [
    for (final r in _rows(await api.get('$_c/friends/requests'), 'requests'))
      FriendRequest(id: r['id'] as String, fromName: r['from_name'] as String? ?? ''),
  ];

  Future<List<Blocked>> blocks() async => [
    for (final b in _rows(await api.get('$_c/blocks'), 'blocked'))
      Blocked(id: b['character_id'] as String, name: b['name'] as String? ?? ''),
  ];

  /// Returns true when the other side had already asked (instant friendship).
  Future<bool> addFriend(String name) async =>
      (await api.post('$_c/friends/request', {'target_name': name}))['friends'] as bool? ?? false;

  Future<void> answer(String requestId, {required bool accept}) =>
      api.post('$_c/friends/requests/$requestId/${accept ? 'accept' : 'decline'}', const {});

  Future<void> removeFriend(String id) => api.post('$_c/friends/remove', {'target_character_id': id});

  Future<void> block(String name) => api.post('$_c/blocks', {'target_name': name});

  Future<void> unblock(String id) => api.post('$_c/blocks/remove', {'target_character_id': id});

  // --- mail -------------------------------------------------------------------

  Future<List<MailSummary>> inbox() async => [
    for (final m in _rows(await api.get('$_c/mail'), 'mail')) MailSummary.fromJson(m),
  ];

  /// Opening a mail marks it read on the server.
  Future<Mail> read(String id) async {
    final j = await api.get('$_c/mail/$id');
    return Mail(MailSummary.fromJson(j), j['body'] as String? ?? '', [
      for (final a in _rows(j, 'attachments'))
        MailAttachment(
          templateCode: a['template_code'] as String? ?? '',
          rarity: a['rarity'] as String? ?? 'common',
          quantity: (a['stack_quantity'] as num?)?.toInt() ?? 1,
        ),
    ]);
  }

  /// Gold attachments only in v1 client; returns the postage charged.
  Future<int> sendMail({required String to, required String subject, required String body, int gold = 0}) async {
    final j = await api.post('$_c/mail', {
      'recipient_name': to,
      'subject': subject,
      'body': body,
      if (gold > 0) 'gold_copper': gold,
    });
    return (j['postage_copper'] as num?)?.toInt() ?? 0;
  }

  Future<void> claim(String id) => api.post('$_c/mail/$id/claim', const {});

  Future<void> deleteMail(String id) => api.post('$_c/mail/$id/delete', const {});
}
