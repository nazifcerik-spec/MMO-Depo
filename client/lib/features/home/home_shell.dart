import 'dart:async';

import 'package:flutter/material.dart';

import '../shop/cosmetics.dart';
import '../shop/shop_screen.dart';
import '../helpers/helpers_screen.dart';
import '../pets/pets_screen.dart';
import '../story/story_repository.dart';
import '../story/story_screen.dart';
import '../daily/daily_reward.dart';
import '../guide/guide.dart';
import '../codex/codex.dart';
import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/browser.dart';
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
import '../ranking/ranking_screen.dart';
import '../quests/quest_repository.dart';
import '../quests/quests_screen.dart';
import '../social/social_repository.dart';
import '../social/social_screen.dart';
import '../workshop/workshop_repository.dart';
import '../workshop/workshop_view.dart';
import '../inventory/bank_screen.dart';
import '../player/player_repository.dart';
import '../player/player_screen.dart';
import '../support/support_repository.dart';
import '../support/support_screen.dart';
import 'game_menu.dart';
import 'overview_view.dart';

/// In-game shell for one character: bottom navigation (the [GameMenu] on wide
/// screens, where pages open beside it, DD-370), the live event
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

  /// The wide screens' first page; phones have no such tab.
  static const _overview = 5;

  /// Open sections of the Adventure and Workshop tabs.
  final _adventure = ValueNotifier(0);
  final _workshop = ValueNotifier(0);

  /// Pages beside the side menu; [_open] is the menu code of the open one.
  final _pages = GlobalKey<NavigatorState>();
  final _root = GlobalKey();
  String? _open;
  int _opens = 0;

  bool _wide = false;
  CurrentAction? _job;
  Timer? _jobTimer;
  int? _gold;
  late final _sheet = ValueNotifier(widget.initial);
  late final _players = PlayerRepository(widget.auth.api, widget.initial.id);
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
    WidgetsBinding.instance.addPostFrameCallback((_) => _firstScreen());
    _loadProfile();
    widget.auth.api.onMutation = _jobChanged;
    _adventure.addListener(_sectionChanged);
    _workshop.addListener(_sectionChanged);
    _lifecycle = AppLifecycleListener(
      onResume: () {
        _live.start();
        _refreshUnread();
        _refreshQuests();
        _jobChanged();
      },
      onPause: _live.stop,
      onHide: _live.stop,
      onShow: _live.start,
    );
  }

  /// The tutorial first for a new player (DD-356), then the daily reward.
  bool _tutorial = false;

  Future<void> _firstScreen() async {
    final pending = await TutorialState.pending(widget.initial.level);
    if (!mounted) return;
    if (pending) {
      setState(() => _tutorial = true);
    } else {
      _openDaily(onlyClaimable: true);
    }
    // The address named a player (DD-371).
    if (takeProfileLink() case final name?) _openPlayer(name);
  }

  void _showTab(int i) {
    _pages.currentState?.popUntil((r) => r.isFirst);
    setState(() {
      _tab = i;
      _visited.add(i);
    });
  }

  void _sectionChanged() => setState(() {});

  /// Signed-in calls come in bursts (collect, start again): ask once after.
  void _jobChanged() {
    if (!_wide) return;
    _jobTimer?.cancel();
    _jobTimer = Timer(const Duration(milliseconds: 400), _refreshJob);
  }

  Future<void> _refreshJob() async {
    try {
      final a = await _players.action();
      if (mounted) setState(() => _job = a);
    } on ApiError {
      // the card keeps what it showed
    }
  }

  /// Where the running job is collected; Gathering when idle.
  void _openJob() {
    switch (_job?.kind) {
      case 'hunt':
        _showSection(3, _adventure, 0);
      case 'dungeon' || 'raid':
        _showSection(3, _adventure, 1);
      case 'travel':
        _showSection(3, _adventure, 2);
      case 'crafting':
        _showSection(4, _workshop, 0);
      default:
        _showTab(1);
    }
  }

  void _showSection(int tab, ValueNotifier<int> section, int i) {
    section.value = i;
    _showTab(tab);
  }

  void _openPlayer(String name) => _push(PlayerScreen(repo: _players, name: name), 'profile');

  void _endTutorial() {
    final first = _tutorial;
    setState(() {
      _tutorial = false;
      _tab = 0;
    });
    TutorialState.finish();
    if (first) _openDaily(onlyClaimable: true);
  }

  bool _namesRequested = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final wide = isWide(context);
    if (wide && !_wide) {
      if (!_namesRequested) {
        _tab = _overview;
        _visited.add(_overview);
      }
      _refreshJob();
    }
    _wide = wide;
    if (_namesRequested) return;
    _namesRequested = true;
    Codex.use(widget.auth.api); // linked tags (DD-372)
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
    _jobTimer?.cancel();
    if (widget.auth.api.onMutation == _jobChanged) widget.auth.api.onMutation = null;
    _dataVersion.dispose();
    _level.dispose();
    _sheet.dispose();
    _adventure.dispose();
    _workshop.dispose();
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

  /// Opens [screen]: over the game on phones, beside the menu on wide screens,
  /// where it replaces the page open there ([place] is its menu code).
  Future<void> _push(Widget screen, [String? place]) async {
    final p = _profile;
    final theme = Theme.of(context);
    final route = MaterialPageRoute<void>(
      builder: (_) => p == null ? screen : Theme(data: themedFor(theme, p.background), child: screen),
    );
    if (_pages.currentState case final pages? when _wide) {
      final n = ++_opens;
      setState(() => _open = place);
      pages.popUntil((r) => r.isFirst);
      await pages.push(route);
      if (mounted && n == _opens) setState(() => _open = null);
    } else {
      await Navigator.of(context).push(route);
    }
    if (!mounted) return;
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
      'market.outbid' => (l.evMarketOutbid, true),
      'mail.received' => (l.evMailReceived, false),
      'party.invite' => (l.evPartyInvite, false),
      'guild.invite' => (l.evGuildInvite, false),
      'quest.completed' => (l.evQuestComplete, false),
      'raid.matched' => (l.evRaidMatched(ContentNames.of('dungeon', e.payload['raid_code'] as String? ?? '')), true),
      'system.dev_broadcast' => (e.payload['message'] as String?, false),
      _ => (null, false),
    };
    if (refresh) _dataVersion.value++;
    _jobChanged();
    if (e.kind == 'quest.completed') _refreshQuests();
    if (message != null && mounted) {
      toast(context, Text(message));
      _refreshUnread(); // most toast kinds above also land in the inbox
    }
  }

  String get _id => widget.initial.id;

  void _openQuests() => _push(QuestsScreen(repo: _quests, dataVersion: _dataVersion), 'quests');

  void _openNotifications() => _push(NotificationsScreen(repo: _notifications), 'notifications');

  /// [tab]: chat, friends, mail, guild, trade.
  void _openSocial(int tab, String place, {String channel = 'global'}) => _push(
    SocialScreen(
      key: ValueKey(place),
      repo: SocialRepository(widget.auth.api, _id),
      live: _live,
      dataVersion: _dataVersion,
      channel: channel,
      tab: tab,
    ),
    place,
  );

  void _openRanking() => _push(
    RankingScreen(repo: RankingRepository(widget.auth.api, _id), social: SocialRepository(widget.auth.api, _id)),
    'ranking',
  );

  void _openGuide() =>
      _push(GuideScreen(onOpenTab: _showTab, onReplayTutorial: () => setState(() => _tutorial = true)), 'guide');

  void _openStory() => _push(
    StoryScreen(repo: StoryRepository(widget.auth.api, _id), quests: _quests, dataVersion: _dataVersion),
    'story',
  );

  void _openShop() => _push(ShopScreen(repo: _shop, onProfile: (p) => setState(() => _profile = p)), 'credits');

  void _openHelpers() => _push(
    HelpersScreen(
      repo: HelperRepository(widget.auth.api, _id),
      workshop: WorkshopRepository(widget.auth.api, _id),
      dataVersion: _dataVersion,
    ),
    'helpers',
  );

  void _openPets() => _push(
    PetsScreen(
      repo: PetRepository(widget.auth.api, _id),
      adventure: AdventureRepository(widget.auth.api, _id),
      dataVersion: _dataVersion,
    ),
    'pets',
  );

  void _openProgress() => _push(
    ProgressionScreen(repo: ProgressionRepository(widget.auth.api, _id), dataVersion: _dataVersion),
    'progress',
  );

  void _openBank() =>
      _push(BankScreen(repo: InventoryRepository(widget.auth.api, _id), dataVersion: _dataVersion), 'bank');

  /// Every place of the game, for the side menu.
  List<(String, List<MenuPlace>)> _menu(AppLocalizations l) => [
    (
      l.menuHero,
      [
        MenuPlace('overview', Icons.home_outlined, l.menuOverview, () => _showTab(_overview)),
        MenuPlace('character', Icons.person_outline, l.navCharacter, () => _showTab(0)),
        MenuPlace('bag', Icons.backpack_outlined, l.navBag, () => _showTab(2)),
        MenuPlace('bank', Icons.account_balance_outlined, l.bankTitle, _openBank),
        MenuPlace(
          'quests',
          Icons.assignment_outlined,
          l.quests,
          _openQuests,
          badge: _questsReady,
          badgeColor: Palette.success,
        ),
        MenuPlace('story', Icons.auto_stories_outlined, l.storyOpen, _openStory),
      ],
    ),
    (
      l.navAdventure,
      [
        MenuPlace('hunt', Icons.gps_fixed, l.segHunt, () => _showSection(3, _adventure, 0)),
        MenuPlace('dungeon', Icons.castle_outlined, l.segDungeon, () => _showSection(3, _adventure, 1)),
        MenuPlace('travel', Icons.map_outlined, l.segTravel, () => _showSection(3, _adventure, 2)),
        MenuPlace('events', Icons.public, l.segEvents, () => _showSection(3, _adventure, 3)),
      ],
    ),
    (
      l.menuWork,
      [
        MenuPlace('gather', Icons.hourglass_empty, l.navActivity, () => _showTab(1)),
        MenuPlace('craft', Icons.handyman_outlined, l.segCraft, () => _showSection(4, _workshop, 0)),
        MenuPlace('helpers', Icons.groups_2_outlined, l.helpersOpen, _openHelpers),
      ],
    ),
    (
      l.menuTrade,
      [
        MenuPlace('vendor', Icons.storefront_outlined, l.segShop, () => _showSection(4, _workshop, 1)),
        MenuPlace('market', Icons.balance_outlined, l.segMarket, () => _showSection(4, _workshop, 2)),
        MenuPlace('trade', Icons.swap_horiz, l.tabTrade, () => _openSocial(4, 'trade')),
        MenuPlace('credits', Icons.diamond_outlined, l.shopOpen, _openShop),
      ],
    ),
    (
      l.menuCommunity,
      [
        MenuPlace('chat', Icons.forum_outlined, l.tabChat, () => _openSocial(0, 'chat')),
        MenuPlace('friends', Icons.people_outline, l.tabFriends, () => _openSocial(1, 'friends')),
        MenuPlace('mail', Icons.mail_outline, l.tabMail, () => _openSocial(2, 'mail')),
        MenuPlace('guild', Icons.shield_outlined, l.tabGuild, () => _openSocial(3, 'guild')),
        MenuPlace('ranking', Icons.leaderboard_outlined, l.ranking, _openRanking),
      ],
    ),
    (
      l.menuGrowth,
      [
        MenuPlace('progress', Icons.emoji_events_outlined, l.progress, _openProgress),
        MenuPlace('pets', Icons.pets, l.petsOpen, _openPets),
        MenuPlace('profile', Icons.badge_outlined, l.profileMine, () => _openPlayer(widget.initial.name)),
      ],
    ),
    (
      l.accountMenu,
      [
        MenuPlace('notifications', Icons.notifications_outlined, l.notifications, _openNotifications, badge: _unread),
        MenuPlace('guide', Icons.menu_book_outlined, l.guide, _openGuide),
        MenuPlace(
          'support',
          Icons.support_agent_outlined,
          l.support,
          () => _push(SupportScreen(repo: SupportRepository(widget.auth.api)), 'support'),
        ),
        MenuPlace('switch', Icons.people_alt_outlined, l.switchCharacter, () => Navigator.of(context).pop()),
        MenuPlace('logout', Icons.logout, l.logout, widget.auth.logout),
      ],
    ),
  ];

  /// The menu code of the tab under the pages.
  String get _tabPlace => switch (_tab) {
    0 => 'character',
    1 => 'gather',
    2 => 'bag',
    3 => const ['hunt', 'dungeon', 'travel', 'events'][_adventure.value],
    4 => const ['craft', 'vendor', 'market'][_workshop.value],
    _ => 'overview',
  };

  Widget _page(int i) => switch (i) {
    0 => CharacterSheetView(
      repo: widget.characters,
      initial: widget.initial,
      dataVersion: _dataVersion,
      onSheet: (s) {
        _level.value = s.level;
        _sheet.value = s;
      },
      onGold: (g) {
        if (g != _gold) setState(() => _gold = g);
      },
      live: _live,
      onPartyChat: () => _openSocial(0, 'chat', channel: 'party'),
      boosts: BoostBar(repo: _shop),
      // On wide screens the side menu holds these.
      onOpenDaily: _openDaily,
      onOpenRanking: _wide ? null : _openRanking,
      onOpenGuide: _wide ? null : _openGuide,
      onLogout: _wide ? null : widget.auth.logout,
      onOpenStory: _wide ? null : _openStory,
      onOpenShop: _wide ? null : _openShop,
      onOpenHelpers: _wide ? null : _openHelpers,
      onOpenPets: _wide ? null : _openPets,
      onOpenProgress: _wide ? null : _openProgress,
    ),
    3 => AdventureView(
      repo: AdventureRepository(widget.auth.api, widget.initial.id),
      dataVersion: _dataVersion,
      level: _level,
      live: _live,
      section: _adventure,
    ),
    4 => WorkshopView(
      repo: WorkshopRepository(widget.auth.api, widget.initial.id),
      dataVersion: _dataVersion,
      section: _workshop,
    ),
    _overview => OverviewView(
      players: _players,
      name: widget.initial.name,
      daily: _daily,
      notifications: _notifications,
      dataVersion: _dataVersion,
      job: _job,
      onJob: _openJob,
      onJobDone: _jobChanged,
      questsReady: _questsReady,
      unread: _unread,
      onQuests: _openQuests,
      onNotifications: _openNotifications,
      onDaily: _openDaily,
      boosts: BoostBar(repo: _shop),
    ),
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
    // The chosen background's art and colours (DD-321, DD-323).
    final profile = _profile;
    final theme = profile == null ? Theme.of(context) : themedFor(Theme.of(context), profile.background);
    final accent = theme.colorScheme.primary;
    final wide = isWide(context);
    // A narrowed window has no overview.
    final tab = !wide && _tab == _overview ? 0 : _tab;
    final stack = KeyedSubtree(
      key: _root,
      child: PageFrame(
        maxWidth: gameWidth,
        child: IndexedStack(
          index: tab,
          children: [for (var i = 0; i <= _overview; i++) _visited.contains(i) ? _page(i) : const SizedBox()],
        ),
      ),
    );
    final body = wide ? _beside(l, accent, profile, stack) : _backed(profile, SafeArea(child: stack));
    return Theme(
      data: theme,
      child: Scaffold(
        appBar: wide
            ? null
            : AppBar(
                leading: IconButton(
                  tooltip: l.switchCharacter,
                  icon: const Icon(Icons.people_alt_outlined),
                  onPressed: () => Navigator.of(context).pop(),
                ),
                title: Row(
                  children: [
                    if (profile != null) ...[
                      AvatarBadge(code: profile.avatar, size: 30),
                      const SizedBox(width: Space.sm),
                    ],
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
                    onPressed: _openQuests,
                  ),
                  IconButton(
                    tooltip: l.notifications,
                    icon: Badge(
                      isLabelVisible: _unread > 0,
                      label: Text(_unread > 99 ? '99+' : '$_unread'),
                      child: const Icon(Icons.notifications_outlined),
                    ),
                    onPressed: _openNotifications,
                  ),
                  IconButton(
                    tooltip: l.social,
                    icon: const Icon(Icons.forum_outlined),
                    onPressed: () => _openSocial(0, 'chat'),
                  ),
                ],
              ),
        body: Stack(
          children: [
            Positioned.fill(child: body),
            if (_tutorial) ...[
              const Positioned.fill(child: ModalBarrier(color: Colors.black38, dismissible: false)),
              Positioned(
                left: Space.md,
                right: Space.md,
                bottom: Space.md,
                child: NarrowColumn(
                  child: TutorialCard(onTab: _showTab, onClose: _endTutorial),
                ),
              ),
            ],
          ],
        ),
        bottomNavigationBar: wide
            ? null
            : NavigationBar(
                selectedIndex: tab,
                backgroundColor: Palette.bgSurface,
                indicatorColor: accent.withValues(alpha: 0.18),
                onDestinationSelected: _showTab,
                destinations: [
                  for (final (icon, selected, label) in tabs)
                    NavigationDestination(
                      icon: Icon(icon),
                      selectedIcon: Icon(selected, color: accent),
                      label: label,
                    ),
                ],
              ),
      ),
    );
  }

  Widget _backed(CosmeticProfile? profile, Widget child) =>
      profile == null ? child : BackgroundArt(code: profile.background, child: child);

  /// Wide screens: the menu, and the pages beside it in their own navigator so
  /// the menu stays.
  Widget _beside(AppLocalizations l, Color accent, CosmeticProfile? profile, Widget stack) {
    final groups = _menu(l);
    final here = _tabPlace;
    final place = [for (final (_, places) in groups) ...places].firstWhere((p) => p.code == here);
    return Row(
      children: [
        ValueListenableBuilder(
          valueListenable: _sheet,
          builder: (context, sheet, _) => GameMenu(
            groups: groups,
            current: _open ?? here,
            sheet: sheet,
            profile: profile,
            gold: _gold,
            job: _job,
            onJob: _openJob,
            onJobDone: _jobChanged,
            compact: GameMenu.isCompact(context),
          ),
        ),
        Expanded(
          child: NavigatorPopHandler<void>(
            onPopWithResult: (_) => _pages.currentState?.maybePop(),
            child: Navigator(
              key: _pages,
              onDidRemovePage: (_) {},
              pages: [
                MaterialPage<void>(
                  key: const ValueKey('tabs'),
                  child: Material(
                    color: Palette.bgPrimary,
                    child: _backed(
                      profile,
                      SafeArea(
                        left: false,
                        child: Column(
                          children: [
                            _heading(place, accent),
                            Expanded(child: stack),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _heading(MenuPlace place, Color accent) => PageFrame(
    maxWidth: gameWidth,
    child: Padding(
      padding: const EdgeInsets.fromLTRB(Space.screen, Space.md, Space.screen, Space.sm),
      child: Row(
        children: [
          menuArt(place.code, place.icon, 36, accent),
          const SizedBox(width: Space.md),
          Text(place.label, style: Theme.of(context).textTheme.titleLarge),
        ],
      ),
    ),
  );
}
