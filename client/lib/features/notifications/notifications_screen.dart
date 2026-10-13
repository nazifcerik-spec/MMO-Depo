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
  List<InboxItem> _rewards = const [];
  int _unread = 0;
  int _tab = 0; // 0 notifications, 1 reward history (DD-242)

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final (items, unread) = await widget.repo.inbox();
    final (rewards, _) = await widget.repo.inbox(rewards: true);
    _items = items;
    _rewards = rewards;
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
              Padding(
                padding: const EdgeInsets.only(bottom: Space.sm),
                child: SizedBox(
                  width: double.infinity,
                  child: SegmentedButton<int>(
                    showSelectedIcon: false,
                    segments: [
                      ButtonSegment(
                        value: 0,
                        label: Text(l.notifTabInbox),
                        icon: const Icon(Icons.notifications_none, size: 18),
                      ),
                      ButtonSegment(
                        value: 1,
                        label: Text(l.notifTabRewards),
                        icon: const Icon(Icons.inventory_2_outlined, size: 18),
                      ),
                    ],
                    selected: {_tab},
                    onSelectionChanged: (s) => setState(() => _tab = s.first),
                  ),
                ),
              ),
              if (_tab == 1) ...[
                if (_rewards.isEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: Space.xxl),
                    child: MessageView(message: l.rewardsNone, icon: Icons.inventory_2_outlined),
                  ),
                for (final n in _rewards) RewardTile(entry: n),
              ],
              if (_tab == 0 && _items.isEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Space.xxl),
                  child: MessageView(message: l.inboxNothing, icon: Icons.notifications_none),
                ),
              if (_tab == 0)
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

/// One collected contract: where it came from, gold / XP / level and the
/// drops as "Iron Ore × 5" chips in their rarity colour (DD-242).
class RewardTile extends StatelessWidget {
  const RewardTile({super.key, required this.entry});

  final InboxItem entry;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final v = entry.vars;
    final target = v['target'] as String? ?? '';
    final title = switch (v['source']) {
      'hunt' => l.rewardFromHunt(ContentNames.of('mob', target)),
      'dungeon' => l.rewardFromDungeon(ContentNames.of('dungeon', target)),
      _ => l.rewardFromActivity(activityName(l, target)),
    };
    final gold = (v['gold_copper'] as num?)?.toInt() ?? 0;
    final xp = (v['xp'] as num?)?.toInt() ?? 0;
    final level = (v['new_level'] as num?)?.toInt();
    final items = [
      for (final i in (v['items'] as List? ?? const []))
        if (i is Map<String, dynamic>) i,
    ];
    return Card(
      margin: const EdgeInsets.only(bottom: Space.sm),
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(child: Text(title, style: text.titleSmall)),
                Text(
                  l.timeAgo(shortDuration(l, DateTime.now().difference(entry.createdAt))),
                  style: text.bodySmall?.copyWith(color: Palette.textMuted),
                ),
              ],
            ),
            if (gold > 0 || xp > 0 || level != null) ...[
              const SizedBox(height: Space.xs),
              Wrap(
                spacing: Space.md,
                crossAxisAlignment: WrapCrossAlignment.center,
                children: [
                  if (gold > 0) GoldText(gold, size: 14),
                  if (xp > 0)
                    Text(
                      l.xpGain(xp),
                      style: const TextStyle(color: Palette.epic, fontWeight: FontWeight.w600),
                    ),
                  if (level != null)
                    Text(
                      l.rewardLevelUp(level),
                      style: const TextStyle(color: Palette.accent, fontWeight: FontWeight.w700),
                    ),
                ],
              ),
            ],
            if (items.isNotEmpty) ...[
              const SizedBox(height: Space.sm),
              Wrap(
                spacing: Space.xs,
                runSpacing: Space.xs,
                children: [
                  for (final i in items)
                    _RewardChip(
                      label: stackLabel(
                        itemName(i['code'] as String? ?? '', i['prefix'] as String?, i['suffix'] as String?),
                        (i['count'] as num?)?.toInt() ?? 1,
                      ),
                      rarity: i['rarity'] as String? ?? 'common',
                    ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _RewardChip extends StatelessWidget {
  const _RewardChip({required this.label, required this.rarity});

  final String label;
  final String rarity;

  @override
  Widget build(BuildContext context) {
    final c = rarityColor(rarity);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 2),
      decoration: BoxDecoration(
        color: Palette.bgElevated,
        borderRadius: BorderRadius.circular(Radii.sm),
        border: Border.all(color: c.withValues(alpha: 0.8)),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: 12, color: c, fontWeight: FontWeight.w600),
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
