import '../../core/api_client.dart';

class InboxItem {
  const InboxItem({required this.id, required this.kind, required this.vars, required this.createdAt, this.readAt});

  factory InboxItem.fromJson(Map<String, dynamic> j) {
    final body = j['body'];
    final vars = body is Map<String, dynamic> ? body['body_vars'] : null;
    return InboxItem(
      id: j['id'] as String,
      kind: j['kind'] as String? ?? '',
      vars: vars is Map<String, dynamic> ? vars : const {},
      createdAt: DateTime.tryParse(j['created_at'] as String? ?? '') ?? DateTime.now(),
      readAt: DateTime.tryParse(j['read_at'] as String? ?? ''),
    );
  }

  final String id;
  final String kind;
  final Map<String, dynamic> vars;
  final DateTime createdAt;
  final DateTime? readAt;

  bool get unread => readAt == null;
}

/// Delivery mode per kind (server notification.Mode).
const notificationModes = ['push_and_inbox', 'inbox_only', 'silent_push', 'off'];

/// Kinds the player can tune; the server has more (not produced yet) and
/// system_broadcast is locked.
const tunableKinds = [
  'activity_completed',
  'dungeon_completed',
  'travel_completed',
  'world_event_active',
  'market_sold',
  'mail_received',
  'party_invite',
  'raid_matched',
  'guild_invite',
  'achievement_unlocked',
];

class NotificationPrefs {
  const NotificationPrefs({required this.kinds, required this.quietHours, this.muteUntil});

  factory NotificationPrefs.fromJson(Map<String, dynamic> j) => NotificationPrefs(
    kinds: {for (final e in (j['kinds'] as Map? ?? const {}).entries) e.key as String: e.value as String},
    quietHours: {
      for (final k in const ['quiet_hours_start', 'quiet_hours_end', 'quiet_hours_tz']) k: j[k],
    },
    muteUntil: DateTime.tryParse(j['global_mute_until'] as String? ?? ''),
  );

  final Map<String, String> kinds;

  /// Echoed back on every update: the server replaces quiet hours wholesale.
  final Map<String, dynamic> quietHours;
  final DateTime? muteUntil;

  bool get muted => muteUntil != null && muteUntil!.isAfter(DateTime.now());
}

/// Account-scoped inbox + delivery preferences.
class NotificationRepository {
  const NotificationRepository(this.api);

  final ApiClient api;

  /// [rewards]: the reward history (DD-242) instead of the regular inbox.
  Future<(List<InboxItem>, int)> inbox({bool rewards = false}) async {
    final j = await api.get(rewards ? '/v1/notifications?rewards=1' : '/v1/notifications');
    return (
      [for (final n in (j['notifications'] as List? ?? const []).cast<Map<String, dynamic>>()) InboxItem.fromJson(n)],
      (j['unread'] as num?)?.toInt() ?? 0,
    );
  }

  Future<int> unreadCount() async => ((await api.get('/v1/notifications?unread=1'))['unread'] as num?)?.toInt() ?? 0;

  Future<void> markRead(String id) => api.post('/v1/notifications/$id/read', const {});
  Future<void> markAllRead() => api.post('/v1/notifications/read-all', const {});

  /// Removes one inbox / reward-history entry (DD-245).
  Future<void> delete(String id) => api.delete('/v1/notifications/$id');

  Future<NotificationPrefs> prefs() async => NotificationPrefs.fromJson(await api.get('/v1/notification/preference'));

  Future<NotificationPrefs> setMode(NotificationPrefs current, String kind, String mode) async =>
      NotificationPrefs.fromJson(
        await api.put('/v1/notification/preference', {
          ...current.quietHours,
          'kinds': {kind: mode},
        }),
      );

  Future<void> mute(int hours) => api.post('/v1/notification/preference/global-mute', {'hours': hours});
  Future<void> unmute() => api.delete('/v1/notification/preference/global-mute');
}
