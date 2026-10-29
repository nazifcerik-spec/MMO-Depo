import 'dart:async';

import 'package:flutter/material.dart';

import '../shop/cosmetics.dart';
import '../shop/shop_screen.dart';
import '../helpers/helpers_screen.dart';
import '../story/story_screen.dart';
import '../daily/daily_reward.dart';
import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../activity/activity_repository.dart';
import '../activity/activity_view.dart';
import '../adventure/adventure_repository.dart';
import '../adventure/adventure_view.dart';
import '../auth/auth_controller.dart';
import '../character/character_repository.dart';
import '../character/character_sheet_view.dart';
import '../character/models.dart';
import '../inventory/inventory_repository.dart';
import '../inventory/inventory_view.dart';
import '../notifications/notification_repository.dart';
import '../notifications/notifications_screen.dart';
import '../progression/progression_repository.dart';
import '../progression/progression_screen.dart';
import '../quests/quest_repository.dart';
import '../quests/quests_screen.dart';
import '../social/social_repository.dart';
import '../social/social_screen.dart';
import '../workshop/workshop_repository.dart';
import '../workshop/workshop_view.dart';

/// In-game shell for one character: bottom navigation, the live event
/// socket (foreground only), and push toasts. Tabs are built on first visit
/// and kept alive afterwards so switching tabs never refetches.
class HomeShell extends StatefulWidget {
  const HomeShell({super.key, required this.auth, required this.characters, required this.initial});

  final AuthController auth;
  final CharacterRepository characters;
  final CharacterSheet initial;

  @override
  State<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends State<HomeShell> {
  final _dataVersion = ValueNotifier(0);

  /// Current character level: level gates in other tabs follow it as soon
  /// as a refreshed sheet shows a level-up (no manual reload).
  late final _level = ValueNotifier(widget.initial.level);
  final _visited = <int>{0};
  int _tab = 0;
  int _unread = 0;
  int _questsReady = 0;
  late final _notifications = NotificationRepository(widget.auth.api);
  late final _quests = QuestRepository(widget.auth.api, widget.initial.id);

  LiveEvents get _live => widget.auth.live;
  late final AppLifecycleListener _lifecycle;
  StreamSubscription<LiveEvent>? _sub;

  @override
  void initState() {
    super.initState();
    _sub = _live.events.listen(_onEvent);
    _live.start();
    _refreshUnread();
    _refreshQuests();
    WidgetsBinding.instance.addPostFrameCallback((_) => _openDaily(onlyClaimable: true));
    _loadProfile();
    _lifecycle = AppLifecycleListener(
      onResume: () {
        _live.start();
        _refreshUnread();
        _refreshQuests();
      },
      onPause: _live.stop,
      onHide: _live.stop,
      onShow: _live.start,
    );
  }

  bool _namesRequested = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_namesRequested) return;
    _namesRequested = true;
    // Content names in the device language (DD-230); tabs rebuild with them.
    ContentNames.load(widget.auth.api, Localizations.localeOf(context).languageCode).then((_) {
      if (mounted) setState(() {});
    });
  }

  @override
  void dispose() {
    _lifecycle.dispose();
    _sub?.cancel();
    _live.stop();
    _dataVersion.dispose();
    _level.dispose();
    super.dispose();
  }

  /// Badge counts; failures keep the last value (badges are hints only).
  Future<void> _refreshUnread() async {
    try {
      final n = await _notifications.unreadCount();
      if (mounted && n != _unread) setState(() => _unread = n);
    } on ApiError {
      // ignore
    }
  }

  Future<void> _refreshQuests() async {
    try {
      final n = (await _quests.board()).readyCount;
      if (mounted && n != _questsReady) setState(() => _questsReady = n);
    } on ApiError {
      // ignore
    }
  }

  late final _daily = DailyRepository(widget.auth.api, widget.initial.id);
  late final _shop = ShopRepository(widget.auth.api, widget.initial.id);

  /// The account's avatar and page background (DD-321).
  CosmeticProfile? _profile;

  Future<void> _loadProfile() async {
    try {
      final p = await _shop.cosmetics();
      if (mounted) setState(() => _profile = p);
    } on ApiError {
      // the default look stays
    }
  }

  /// The daily login reward (DD-311): opened on entering the game while
  /// today's is unclaimed, and from the character sheet.
  Future<void> _openDaily({bool onlyClaimable = false}) async {
    try {
      final st = await _daily.status();
      if (!mounted || (onlyClaimable && !st.claimable)) return;
      if (await showDailyReward(context, _daily, st)) _dataVersion.value++;
    } on ApiError {
      // the reward waits for the next visit
    }
  }

  Future<void> _push(Widget screen) async {
    await Navigator.of(context).push(MaterialPageRoute(builder: (_) => screen));
    _refreshUnread();
    _refreshQuests();
  }

  void _onEvent(LiveEvent e) {
    final l = AppLocalizations.of(context);
    // Party events name the character they are for; others on the account skip them.
    if (e.payload['character_id'] case final String c when c != widget.initial.id) {
      if (e.kind == 'raid.matched') return;
    }
    final (message, refresh) = switch (e.kind) {
      'hunt.reveal' || 'dungeon.reveal' => (null, true),
      'achievement.unlocked' => (l.evAchievement, true),
      'mastery.tier_up' => (l.evMasteryTier, true),
      'market.sold' => (l.evMarketSold, true),
      'mail.received' => (l.evMailReceived, false),
      'party.invite' => (l.evPartyInvite, false),
      'guild.invite' => (l.evGuildInvite, false),
      'quest.completed' => (l.evQuestComplete, false),
      'raid.matched' => (l.evRaidMatched(ContentNames.of('dungeon', e.payload['raid_code'] as String? ?? '')), true),
      'system.dev_broadcast' => (e.payload['message'] as String?, false),
      _ => (null, false),
    };
    if (refresh) _dataVersion.value++;
    if (e.kind == 'quest.completed') _refreshQuests();
    if (message != null && mounted) {
      toast(context, Text(message));
      _refreshUnread(); // most toast kinds above also land in the inbox
    }
  }

  Widget _page(int i) => switch (i) {
    0 => CharacterSheetView(
      repo: widget.characters,
      initial: widget.initial,
      dataVersion: _dataVersion,
      onSheet: (s) => _level.value = s.level,
      live: _live,
      onPartyChat: () => _push(
        SocialScreen(
          repo: SocialRepository(widget.auth.api, widget.initial.id),
          live: _live,
          dataVersion: _dataVersion,
          channel: 'party',
        ),
      ),
      onOpenDaily: _openDaily,
      onOpenStory: () => _push(StoryScreen(api: widget.auth.api, level: _level.value)),
      boosts: BoostBar(repo: _shop),
      onOpenShop: () => _push(ShopScreen(repo: _shop, onProfile: (p) => setState(() => _profile = p))),
      onOpenHelpers: () => _push(
        HelpersScreen(
          repo: HelperRepository(widget.auth.api, widget.initial.id),
          workshop: WorkshopRepository(widget.auth.api, widget.initial.id),
          dataVersion: _dataVersion,
        ),
      ),
      onOpenProgress: () => _push(
        ProgressionScreen(repo: ProgressionRepository(widget.auth.api, widget.initial.id), dataVersion: _dataVersion),
      ),
    ),
    3 => AdventureView(
      repo: AdventureRepository(widget.auth.api, widget.initial.id),
      dataVersion: _dataVersion,
      level: _level,
      live: _live,
    ),
    4 => WorkshopView(repo: WorkshopRepository(widget.auth.api, widget.initial.id), dataVersion: _dataVersion),
    2 => InventoryView(repo: InventoryRepository(widget.auth.api, widget.initial.id), dataVersion: _dataVersion),
    _ => ActivityView(
      repo: ActivityRepository(widget.auth.api),
      characterId: widget.initial.id,
      dataVersion: _dataVersion,
    ),
  };

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final s = widget.initial;
    final tabs = [
      (Icons.person_outline, Icons.person, l.navCharacter),
      (Icons.hourglass_empty, Icons.hourglass_bottom, l.navActivity),
      (Icons.backpack_outlined, Icons.backpack, l.navBag),
      (Icons.explore_outlined, Icons.explore, l.navAdventure),
      (Icons.handyman_outlined, Icons.handyman, l.navWorkshop),
    ];
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          tooltip: l.switchCharacter,
          icon: const Icon(Icons.people_alt_outlined),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Row(
          children: [
            if (_profile case final p?) ...[AvatarBadge(code: p.avatar, size: 30), const SizedBox(width: Space.sm)],
            Flexible(child: Text(s.name, overflow: TextOverflow.ellipsis)),
          ],
        ),

        actions: [
          IconButton(
            tooltip: l.quests,
            icon: Badge(
              isLabelVisible: _questsReady > 0,
              backgroundColor: Palette.success,
              label: Text('$_questsReady'),
              child: const Icon(Icons.assignment_outlined),
            ),
            onPressed: () => _push(QuestsScreen(repo: _quests, dataVersion: _dataVersion)),
          ),
          IconButton(
            tooltip: l.notifications,
            icon: Badge(
              isLabelVisible: _unread > 0,
              label: Text(_unread > 99 ? '99+' : '$_unread'),
              child: const Icon(Icons.notifications_outlined),
            ),
            onPressed: () => _push(NotificationsScreen(repo: _notifications)),
          ),
          IconButton(
            tooltip: l.social,
            icon: const Icon(Icons.forum_outlined),
            onPressed: () => _push(
              SocialScreen(repo: SocialRepository(widget.auth.api, s.id), live: _live, dataVersion: _dataVersion),
            ),
          ),
          IconButton(tooltip: l.logout, icon: const Icon(Icons.logout), onPressed: widget.auth.logout),
        ],
      ),
      body: DecoratedBox(
        decoration: BoxDecoration(gradient: _profile == null ? null : backgroundGradient(_profile!.background)),
        child: SafeArea(
          child: PageFrame(
            child: IndexedStack(
              index: _tab,
              children: [for (var i = 0; i < tabs.length; i++) _visited.contains(i) ? _page(i) : const SizedBox()],
            ),
          ),
        ),
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _tab,
        backgroundColor: Palette.bgSurface,
        indicatorColor: Palette.accent.withValues(alpha: 0.18),
        onDestinationSelected: (i) => setState(() {
          _tab = i;
          _visited.add(i);
        }),
        destinations: [
          for (final (icon, selected, label) in tabs)
            NavigationDestination(
              icon: Icon(icon),
              selectedIcon: Icon(selected, color: Palette.accent),
              label: label,
            ),
        ],
      ),
    );
  }
}
