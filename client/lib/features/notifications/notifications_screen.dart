import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'notification_repository.dart';
import '../../shared/content_names.dart';

String notificationKindName(AppLocalizations l, String kind) => switch (kind) {
  'activity_completed' => l.nkActivity,
  'dungeon_completed' => l.nkDungeon,
  'travel_completed' => l.nkTravel,
  'world_event_active' => l.nkWorldEvent,
  'market_sold' => l.nkMarketSold,
  'mail_received' => l.nkMail,
  'party_invite' => l.nkPartyInvite,
  'guild_invite' => l.nkGuildInvite,
  'achievement_unlocked' => l.nkAchievement,
  'support_reply' => l.nkSupport,
  _ => l.nkSystem,
};

String _modeName(AppLocalizations l, String mode) => switch (mode) {
  'push_and_inbox' => l.modePushInbox,
  'inbox_only' => l.modeInbox,
  'silent_push' => l.modeSilent,
  _ => l.modeOff,
};

/// Detail line from the producer's body_vars (server-side notification.Notice).
/// Values are rendered as plain text only.
String? _detail(InboxItem n) {
  final v = n.vars;
  String? s(String k) => v[k] is String && (v[k] as String).isNotEmpty ? v[k] as String : null;
  return switch (n.kind) {
    'market_sold' when s('template_code') != null =>
      '${itemName(s('template_code')!)} · +${goldString((v['net_copper'] as num?)?.toInt() ?? 0)}',
    'mail_received' => s('subject'),
    'guild_invite' => s('guild_name'),
    'support_reply' => s('subject'),
    'travel_completed' when s('region') != null => ContentNames.of('region', s('region')!),
    'world_event_active' when s('event') != null => ContentNames.of('worldevent', s('event')!),
    'achievement_unlocked' when s('code') != null => ContentNames.of('achievement', s('code')!),
    _ => null,
  };
}

/// Inbox (newest first, capped at 100 by the server). Opening an unread row
/// marks it read; the shell refreshes its badge when this route pops.
class NotificationsScreen extends StatefulWidget {
  const NotificationsScreen({super.key, required this.repo});

  final NotificationRepository repo;

  @override
  State<NotificationsScreen> createState() => _NotificationsScreenState();
}

class _NotificationsScreenState extends State<NotificationsScreen> with PanelState {
  List<InboxItem> _items = const [];
  int _unread = 0;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final (items, unread) = await widget.repo.inbox();
    _items = items;
    _unread = unread;
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(
        title: Text(l.notifications),
        actions: [
          if (_unread > 0)
            IconButton(
              tooltip: l.markAllRead,
              icon: const Icon(Icons.done_all),
              onPressed: busy
                  ? null
                  : () => run(() async {
                      await widget.repo.markAllRead();
                      await load();
                    }),
            ),
          IconButton(
            tooltip: l.notifSettings,
            icon: const Icon(Icons.tune),
            onPressed: () =>
                Navigator.of(context)
                    .push(MaterialPageRoute(builder: (_) => NotificationSettingsScreen(repo: widget.repo))),
          ),
        ],
      ),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              if (_items.isEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Space.xxl),
                  child: MessageView(message: l.inboxNothing, icon: Icons.notifications_none),
                ),
              for (final n in _items)
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  leading: Icon(
                    n.unread ? Icons.circle : Icons.circle_outlined,
                    size: 10,
                    color: n.unread ? Palette.accent : Palette.textMuted,
                  ),
                  minLeadingWidth: 10,
                  title: Text(
                    notificationKindName(l, n.kind),
                    style: n.unread ? const TextStyle(fontWeight: FontWeight.w600) : null,
                  ),
                  subtitle: Text(
                    [?_detail(n), l.timeAgo(shortDuration(l, DateTime.now().difference(n.createdAt)))].join('\n'),
                    style: text.bodySmall,
                  ),
                  onTap: n.unread
                      ? () => run(() async {
                          await widget.repo.markRead(n.id);
                          await load();
                        })
                      : null,
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class NotificationSettingsScreen extends StatefulWidget {
  const NotificationSettingsScreen({super.key, required this.repo});

  final NotificationRepository repo;

  @override
  State<NotificationSettingsScreen> createState() => _NotificationSettingsScreenState();
}

class _NotificationSettingsScreenState extends State<NotificationSettingsScreen> with PanelState {
  NotificationPrefs? _prefs;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _prefs = await widget.repo.prefs();

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Scaffold(
      appBar: AppBar(title: Text(l.notifSettings)),
      body: SafeArea(
        child: PageFrame(
          child: frame(() {
            final p = _prefs!;
            return [
              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(l.muteAll),
                subtitle: Text(
                  p.muted ? l.mutedFor(shortDuration(l, p.muteUntil!.difference(DateTime.now()))) : l.muteHint,
                ),
                value: p.muted,
                onChanged: (on) => run(() async {
                  await (on ? widget.repo.mute(24) : widget.repo.unmute());
                  await load();
                }),
              ),
              const Divider(color: Palette.borderSubtle),
              for (final kind in tunableKinds)
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(notificationKindName(l, kind)),
                  trailing: DropdownButton<String>(
                    value: p.kinds[kind] ?? 'push_and_inbox',
                    underline: const SizedBox(),
                    items: [
                      for (final m in notificationModes) DropdownMenuItem(value: m, child: Text(_modeName(l, m))),
                    ],
                    onChanged: (m) {
                      if (m == null || m == p.kinds[kind]) return;
                      run(() async {
                        final updated = await widget.repo.setMode(p, kind, m);
                        if (mounted) setState(() => _prefs = updated);
                      });
                    },
                  ),
                ),
            ];
          }),
        ),
      ),
    );
  }
}
