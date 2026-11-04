import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';

/// Whether the new-player tutorial still has to run (DD-356). Kept on the
/// device. Unreadable storage counts as done: better no tutorial than one on
/// every start.
class TutorialState {
  TutorialState._();

  static const _key = 'mena.tutorial_done';
  static const _storage = FlutterSecureStorage();

  /// Beyond this level the player clearly knows the game.
  static const newPlayerLevel = 5;

  /// Replaces the stored answer (tests).
  static bool? debugPending;

  static Future<bool> pending(int level) async {
    if (debugPending != null) return debugPending!;
    if (level > newPlayerLevel) return false;
    try {
      return await _storage.read(key: _key) == null;
    } catch (_) {
      return false;
    }
  }

  static Future<void> finish() async {
    debugPending = debugPending == null ? null : false;
    try {
      await _storage.write(key: _key, value: '1');
    } catch (_) {
      // Shown again next time at worst.
    }
  }
}

/// One tutorial page: the bottom tab it shows, and its text.
class _Step {
  const _Step(this.tab, this.icon, this.title, this.body);

  final int tab;
  final IconData icon;
  final String title;
  final String body;
}

List<_Step> _steps(AppLocalizations l) => [
  _Step(0, Icons.waving_hand_outlined, l.tutorialWelcomeTitle, l.tutorialWelcomeBody),
  _Step(0, Icons.person, l.navCharacter, l.tutorialCharacterBody),
  _Step(1, Icons.hourglass_bottom, l.navActivity, l.tutorialGatherBody),
  _Step(2, Icons.backpack, l.navBag, l.tutorialBagBody),
  _Step(3, Icons.explore, l.navAdventure, l.tutorialAdventureBody),
  _Step(4, Icons.handyman, l.navWorkshop, l.tutorialWorkshopBody),
  _Step(0, Icons.assignment_outlined, l.tutorialTopBarTitle, l.tutorialTopBarBody),
  _Step(0, Icons.menu_book_outlined, l.tutorialDoneTitle, l.tutorialDoneBody),
];

/// The tutorial: a card over the game that walks the bottom tabs. [onTab]
/// shows the tab a page talks about; [onClose] ends it, finished or skipped.
class TutorialCard extends StatefulWidget {
  const TutorialCard({super.key, required this.onTab, required this.onClose});

  final ValueChanged<int> onTab;
  final VoidCallback onClose;

  @override
  State<TutorialCard> createState() => _TutorialCardState();
}

class _TutorialCardState extends State<TutorialCard> {
  int _step = 0;

  void _go(int step, int tab) {
    setState(() => _step = step);
    widget.onTab(tab);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final steps = _steps(l);
    final s = steps[_step];
    final last = _step == steps.length - 1;
    return Material(
      key: const ValueKey('tutorial-card'),
      color: Palette.bgElevated,
      elevation: 8,
      borderRadius: BorderRadius.circular(Radii.lg),
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Icon(s.icon, color: Palette.accent),
                const SizedBox(width: Space.sm),
                Expanded(child: Text(s.title, style: text.titleMedium)),
                Text('${_step + 1}/${steps.length}', style: text.bodySmall?.copyWith(color: Palette.textMuted)),
              ],
            ),
            const SizedBox(height: Space.sm),
            Text(s.body, style: text.bodyMedium?.copyWith(height: 1.4)),
            const SizedBox(height: Space.md),
            Row(
              children: [
                if (!last)
                  TextButton(
                    key: const ValueKey('tutorial-skip'),
                    onPressed: widget.onClose,
                    child: Text(l.tutorialSkip),
                  ),
                const Spacer(),
                if (_step > 0)
                  TextButton(
                    key: const ValueKey('tutorial-back'),
                    onPressed: () => _go(_step - 1, steps[_step - 1].tab),
                    child: Text(l.back),
                  ),
                const SizedBox(width: Space.xs),
                FilledButton(
                  key: const ValueKey('tutorial-next'),
                  onPressed: last ? widget.onClose : () => _go(_step + 1, steps[_step + 1].tab),
                  child: Text(last ? l.tutorialStart : l.next),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

/// One guide topic: where a thing is and what it does; [tab] opens it.
class _Topic {
  const _Topic(this.icon, this.title, this.body, {this.tab});

  final IconData icon;
  final String title;
  final String body;
  final int? tab;
}

List<_Topic> _topics(AppLocalizations l) => [
  _Topic(Icons.person, l.guideCharacterTitle, l.guideCharacterBody, tab: 0),
  _Topic(Icons.trending_up, l.guidePromotionTitle, l.guidePromotionBody, tab: 0),
  _Topic(Icons.hourglass_bottom, l.guideGatherTitle, l.guideGatherBody, tab: 1),
  _Topic(Icons.backpack, l.guideBagTitle, l.guideBagBody, tab: 2),
  _Topic(Icons.account_balance_outlined, l.guideBankTitle, l.guideBankBody, tab: 2),
  _Topic(Icons.pest_control_outlined, l.guideHuntTitle, l.guideHuntBody, tab: 3),
  _Topic(Icons.castle_outlined, l.guideDungeonTitle, l.guideDungeonBody, tab: 3),
  _Topic(Icons.map_outlined, l.guideTravelTitle, l.guideTravelBody, tab: 3),
  _Topic(Icons.celebration_outlined, l.guideEventsTitle, l.guideEventsBody, tab: 3),
  _Topic(Icons.handyman, l.guideCraftTitle, l.guideCraftBody, tab: 4),
  _Topic(Icons.store_outlined, l.guideVendorTitle, l.guideVendorBody, tab: 4),
  _Topic(Icons.gavel_outlined, l.guideMarketTitle, l.guideMarketBody, tab: 4),
  _Topic(Icons.assignment_outlined, l.guideQuestsTitle, l.guideQuestsBody),
  _Topic(Icons.notifications_outlined, l.guideNotificationsTitle, l.guideNotificationsBody),
  _Topic(Icons.forum_outlined, l.guideSocialTitle, l.guideSocialBody),
  _Topic(Icons.groups_outlined, l.guidePartyTitle, l.guidePartyBody, tab: 0),
  _Topic(Icons.shield_outlined, l.guideGuildTitle, l.guideGuildBody),
  _Topic(Icons.leaderboard_outlined, l.guideRankingTitle, l.guideRankingBody, tab: 0),
  _Topic(Icons.pets, l.guidePetsTitle, l.guidePetsBody, tab: 0),
  _Topic(Icons.groups_2_outlined, l.guideHelpersTitle, l.guideHelpersBody, tab: 0),
  _Topic(Icons.flag_outlined, l.guideReportTitle, l.guideReportBody),
  _Topic(Icons.manage_accounts_outlined, l.guideAccountTitle, l.guideAccountBody),
];

/// Rehber (DD-356): what is where, at any time. Tapping "Show me" closes the
/// guide on the tab the topic lives on.
class GuideScreen extends StatefulWidget {
  const GuideScreen({super.key, required this.onOpenTab, required this.onReplayTutorial});

  final ValueChanged<int> onOpenTab;
  final VoidCallback onReplayTutorial;

  @override
  State<GuideScreen> createState() => _GuideScreenState();
}

class _GuideScreenState extends State<GuideScreen> {
  String _query = '';

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final q = _query.trim().toLowerCase();
    final topics = [
      for (final t in _topics(l))
        if (q.isEmpty || t.title.toLowerCase().contains(q) || t.body.toLowerCase().contains(q)) t,
    ];
    return Scaffold(
      appBar: AppBar(title: Text(l.guide)),
      body: SafeArea(
        child: PageFrame(
          child: ListView(
            key: const ValueKey('guide-list'),
            padding: const EdgeInsets.all(Space.screen),
            children: [
              TextField(
                key: const ValueKey('guide-search'),
                onChanged: (v) => setState(() => _query = v),
                decoration: InputDecoration(prefixIcon: const Icon(Icons.search), hintText: l.guideSearch),
              ),
              const SizedBox(height: Space.sm),
              if (topics.isEmpty)
                Padding(
                  padding: const EdgeInsets.all(Space.lg),
                  child: Text(l.guideNoMatch, textAlign: TextAlign.center),
                ),
              for (final t in topics)
                ExpansionTile(
                  key: PageStorageKey(t.title),
                  leading: Icon(t.icon, color: Palette.accent),
                  title: Text(t.title),
                  tilePadding: EdgeInsets.zero,
                  childrenPadding: const EdgeInsets.only(bottom: Space.md),
                  expandedCrossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(t.body, style: const TextStyle(height: 1.4, color: Palette.textSecondary)),
                    if (t.tab case final tab?)
                      TextButton(
                        onPressed: () {
                          Navigator.of(context).pop();
                          widget.onOpenTab(tab);
                        },
                        child: Text(l.guideShowMe),
                      ),
                  ],
                ),
              const SizedBox(height: Space.lg),
              OutlinedButton.icon(
                key: const ValueKey('guide-replay'),
                onPressed: () {
                  Navigator.of(context).pop();
                  widget.onReplayTutorial();
                },
                icon: const Icon(Icons.replay),
                label: Text(l.tutorialReplay),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
