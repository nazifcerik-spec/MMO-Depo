import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../daily/daily_reward.dart';
import '../notifications/notification_repository.dart';
import '../notifications/notifications_screen.dart' show RewardTile;
import '../player/player_repository.dart';
import '../player/player_screen.dart' show PlayerBanner;
import 'game_menu.dart' show GameMenu, JobCard;

/// The first page on wide screens (DD-370): the character's banner, what it is
/// doing, what waits today and the latest rewards.
class OverviewView extends StatefulWidget {
  const OverviewView({
    super.key,
    required this.players,
    required this.name,
    required this.daily,
    required this.notifications,
    required this.dataVersion,
    required this.job,
    required this.onJob,
    required this.onJobDone,
    required this.questsReady,
    required this.unread,
    required this.onQuests,
    required this.onNotifications,
    required this.onDaily,
    required this.boosts,
  });

  final PlayerRepository players;
  final String name;
  final DailyRepository daily;
  final NotificationRepository notifications;
  final ValueNotifier<int> dataVersion;
  final CurrentAction? job;
  final VoidCallback onJob;
  final VoidCallback onJobDone;
  final int questsReady;
  final int unread;
  final VoidCallback onQuests;
  final VoidCallback onNotifications;
  final VoidCallback onDaily;
  final Widget boosts;

  @override
  State<OverviewView> createState() => _OverviewViewState();
}

class _OverviewViewState extends State<OverviewView> with PanelState {
  static const _shown = 5;

  PlayerProfile? _me;
  DailyStatus? _daily;
  List<InboxItem> _rewards = const [];

  @override
  void initState() {
    super.initState();
    widget.dataVersion.addListener(load);
    load();
  }

  @override
  void dispose() {
    widget.dataVersion.removeListener(load);
    super.dispose();
  }

  /// Each part on its own: one that fails keeps its last state.
  @override
  Future<void> fetch() async {
    Future<void> part(Future<void> Function() get) async {
      try {
        await get();
      } on ApiError {
        // the rest of the page still shows
      }
    }

    await Future.wait([
      part(() async => _me = await widget.players.profile(widget.name)),
      part(() async => _daily = await widget.daily.status()),
      part(() async => _rewards = (await widget.notifications.inbox(rewards: true)).$1.take(_shown).toList()),
    ]);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(
      () => [
        if (_me case final me?) ...[PlayerBanner(profile: me), const SizedBox(height: Space.lg)],
        Panes(
          main: [
            // The full menu shows it already.
            if (GameMenu.isCompact(context)) ...[
              JobCard(action: widget.job, onOpen: widget.onJob, onDone: widget.onJobDone),
              const SizedBox(height: Space.xl),
            ],
            Text(l.overviewRecent, style: text.titleMedium),
            const SizedBox(height: Space.md),
            if (_rewards.isEmpty) Text(l.overviewNoRewards, style: text.bodySmall),
            for (final r in _rewards) RewardTile(entry: r),
          ],
          side: [
            Text(l.overviewToday, style: text.titleMedium),
            const SizedBox(height: Space.md),
            if (_daily case final d?)
              _Today(
                icon: Icons.card_giftcard_outlined,
                label: d.claimable ? l.overviewDailyReady : l.overviewDailyDone,
                lit: d.claimable,
                onTap: widget.onDaily,
              ),
            _Today(
              icon: Icons.assignment_outlined,
              label: l.overviewQuests(widget.questsReady),
              lit: widget.questsReady > 0,
              onTap: widget.onQuests,
            ),
            _Today(
              icon: Icons.notifications_outlined,
              label: l.overviewUnread(widget.unread),
              lit: widget.unread > 0,
              onTap: widget.onNotifications,
            ),
            const SizedBox(height: Space.md),
            widget.boosts,
          ],
        ),
      ],
    );
  }
}

class _Today extends StatelessWidget {
  const _Today({required this.icon, required this.label, required this.lit, required this.onTap});

  final IconData icon;
  final String label;

  /// Something waits here.
  final bool lit;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final c = lit ? Palette.success : Palette.textSecondary;
    return Card(
      margin: const EdgeInsets.only(bottom: Space.sm),
      child: ListTile(
        leading: Icon(icon, color: c),
        title: Text(label, style: TextStyle(color: lit ? Palette.textPrimary : Palette.textSecondary)),
        trailing: const Icon(Icons.chevron_right, color: Palette.textMuted),
        onTap: onTap,
      ),
    );
  }
}
