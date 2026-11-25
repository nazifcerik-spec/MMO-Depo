import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../codex/codex.dart';
import 'progression_repository.dart';
import '../../shared/content_names.dart';

/// Achievements, mastery tracks, faction reputation and title selection.
/// One load fills every tab; [dataVersion] is bumped when the active title
/// changes so the sheet can pick it up.
class ProgressionScreen extends StatefulWidget {
  const ProgressionScreen({super.key, required this.repo, required this.dataVersion});

  final ProgressionRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<ProgressionScreen> createState() => _ProgressionScreenState();
}

class _ProgressionScreenState extends State<ProgressionScreen> with PanelState {
  ProgressionSummary? _summary;
  List<Achievement> _achievements = const [];
  List<TitleGoal> _titles = const [];

  /// Other classes' masteries are folded until asked for.
  bool _allMasteries = false;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait([widget.repo.summary(), widget.repo.achievements(), widget.repo.titles()]);
    _summary = r[0] as ProgressionSummary;
    // Unlocked first (newest on top), then closest to completion.
    _achievements = (r[1] as List<Achievement>)
      ..sort((a, b) {
        if (a.unlocked != b.unlocked) return a.unlocked ? -1 : 1;
        if (a.unlocked) return b.unlockedAt!.compareTo(a.unlockedAt!);
        return (b.progress / b.target).compareTo(a.progress / a.target);
      });
    _titles = r[2] as List<TitleGoal>;
  }

  Future<void> _setTitle(String code) => run(() async {
    await widget.repo.setActiveTitle(code);
    widget.dataVersion.value++;
    await load();
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return DefaultTabController(
      length: 4,
      child: Scaffold(
        appBar: AppBar(
          title: Text(l.progress),
          bottom: TabBar(
            isScrollable: true,
            tabAlignment: TabAlignment.start,
            tabs: [
              Tab(text: l.tabAchievements),
              Tab(text: l.tabMastery),
              Tab(text: l.tabReputation),
              Tab(text: l.tabTitles),
            ],
          ),
        ),
        body: SafeArea(
          child: PageFrame(
            child: TabBarView(
              children: [frame(_achievementRows), frame(_masteryRows), frame(_reputationRows), frame(_titleRows)],
            ),
          ),
        ),
      ),
    );
  }

  List<Widget> _achievementRows() {
    final l = AppLocalizations.of(context);
    return [
      Padding(
        padding: const EdgeInsets.only(bottom: Space.md),
        child: Text(
          l.achievementPoints(_summary?.points ?? 0),
          style: Theme.of(context).textTheme.titleMedium?.copyWith(color: Palette.accent),
        ),
      ),
      for (final a in _achievements) _AchievementTile(a: a),
    ];
  }

  List<Widget> _masteryRows() {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return [
      for (final m in _summary?.mastery ?? const <MasteryTrack>[])
        Card(
          margin: const EdgeInsets.only(bottom: Space.md),
          child: Padding(
            padding: const EdgeInsets.all(Space.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(child: Text(ContentNames.of('track', m.track), style: text.titleMedium)),
                    Text(l.masteryTier(m.tier, m.maxTier), style: text.bodySmall),
                  ],
                ),
                const SizedBox(height: Space.sm),
                StatBar(
                  label: l.xp,
                  value: m.progress,
                  color: Palette.epic,
                  valueText: m.tier >= m.maxTier ? '${m.xp}' : l.xpValue(m.xp, m.nextTierXp),
                ),
                for (final b in m.bonuses)
                  Text(bonusText(l, b), style: const TextStyle(color: Palette.success, fontSize: 13)),
                _Sources(events: m.sources),
              ],
            ),
          ),
        ),
    ];
  }

  List<Widget> _reputationRows() {
    final l = AppLocalizations.of(context);
    return [
      for (final r in _summary?.reputation ?? const <Reputation>[])
        Card(
          margin: const EdgeInsets.only(bottom: Space.md),
          child: Padding(
            padding: const EdgeInsets.all(Space.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        ContentNames.of('faction', r.faction),
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ),
                    Text(l.repDailyLeft(r.dailyLeft), style: Theme.of(context).textTheme.bodySmall),
                  ],
                ),
                Text('${ContentNames.of('reptier', r.tierName)} · ${r.rep}'),
                _Sources(events: r.sources),
              ],
            ),
          ),
        ),
    ];
  }

  List<Widget> _titleRows() {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    List<TitleGoal> of(String source, {bool? own}) => [
      for (final t in _titles)
        if (t.source == source && (own == null || t.ownClass == own)) t,
    ];
    final others = of('class_mastery', own: false);
    Widget header(String s) => Padding(
      padding: const EdgeInsets.fromLTRB(0, Space.lg, 0, Space.sm),
      child: Text(s, style: text.titleSmall?.copyWith(color: Palette.accent)),
    );
    Widget cards(List<TitleGoal> list) => TileGrid(
      spacing: Space.md,
      children: [for (final t in list) _TitleCard(key: ValueKey('title-${t.code}'), t: t, onSet: _setTitle)],
    );
    return [
      Text(l.titlesEarned(_titles.where((t) => t.owned).length, _titles.length), style: text.titleMedium),
      for (final (name, list) in [
        (l.titlesAchievement, of('achievement')),
        (l.titlesLevel, of('level')),
        (l.titlesMastery, of('class_mastery', own: true)),
      ])
        if (list.isNotEmpty) ...[header(name), cards(list)],
      if (others.isNotEmpty) ...[
        InkWell(
          key: const Key('titles-other'),
          onTap: () => setState(() => _allMasteries = !_allMasteries),
          child: Row(
            children: [
              Expanded(child: header(l.titlesOther)),
              Text('${others.length}', style: text.bodySmall),
              Icon(_allMasteries ? Icons.expand_less : Icons.expand_more, color: Palette.textSecondary),
            ],
          ),
        ),
        if (_allMasteries) cards(others),
      ],
    ];
  }
}

/// A mastery is reached at its level, not "once".
int goalCount(String event, int target) => event.startsWith('class_mastery:') ? 1000 : target;

/// What raises a track or a faction, as tags.
class _Sources extends StatelessWidget {
  const _Sources({required this.events});

  final List<String> events;

  @override
  Widget build(BuildContext context) {
    // "any dungeon" and "any boss" open the same list: once is enough.
    final refs = <String, CodexRef>{};
    for (final e in events) {
      final r = eventRef(e);
      if (r != null) refs['${r.kind}/${r.code}'] = r;
    }
    if (refs.isEmpty) return const SizedBox.shrink();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: Space.sm),
        Text(AppLocalizations.of(context).codexGroup('sources'), style: Theme.of(context).textTheme.bodySmall),
        const SizedBox(height: Space.xs),
        LinkTags(refs.values),
      ],
    );
  }
}

/// One title: what earns it in words, how far the character is, what it is
/// about as tags, and Use once earned. Every card has the same size.
class _TitleCard extends StatelessWidget {
  const _TitleCard({super.key, required this.t, required this.onSet});

  final TitleGoal t;
  final ValueChanged<String> onSet;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final done = t.owned || t.target <= 0;
    const goal = TextStyle(fontSize: 13, height: 1.3);
    final tags = [?eventRef(t.event), if (t.achievement.isNotEmpty) CodexRef('achievement', t.achievement)];
    final rewards = [if (t.points > 0) l.pointsShort(t.points), if (t.rewardGold > 0) goldString(t.rewardGold)];
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Icon(
                  t.owned ? Icons.workspace_premium : Icons.workspace_premium_outlined,
                  size: 20,
                  color: t.owned ? Palette.accent : Palette.textMuted,
                ),
                const SizedBox(width: Space.sm),
                Expanded(
                  child: InkWell(
                    onTap: () => showCodex(context, CodexRef('title', t.code)),
                    child: Text(
                      ContentNames.of('title', t.code),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: text.titleMedium,
                    ),
                  ),
                ),
                if (rewards.isNotEmpty) Text(rewards.join(' · '), style: text.bodySmall),
              ],
            ),
            const SizedBox(height: Space.sm),
            SizedBox(
              height: 3 * 1.3 * MediaQuery.textScalerOf(context).scale(13),
              child: Text(goalText(l, t.event, t.target), maxLines: 3, overflow: TextOverflow.ellipsis, style: goal),
            ),
            const SizedBox(height: Space.sm),
            StatBar(
              label: t.owned ? l.titleOwned : l.titleLocked,
              value: done ? 1 : t.progress / t.target,
              color: done ? Palette.success : Palette.accent,
              valueText: done ? '${t.target}/${t.target}' : '${t.progress}/${t.target}',
              stacked: true,
            ),
            const SizedBox(height: Space.sm),
            Row(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: SizedBox(
                      height: 48, // the button's height, shown or not
                      child: Row(
                        children: [
                          for (final r in tags) ...[
                            ConstrainedBox(constraints: const BoxConstraints(maxWidth: 220), child: LinkTag(r)),
                            const SizedBox(width: Space.xs),
                          ],
                        ],
                      ),
                    ),
                  ),
                ),
                if (t.owned) ...[
                  const SizedBox(width: Space.sm),
                  t.active
                      ? OutlinedButton(
                          style: OutlinedButton.styleFrom(minimumSize: const Size(72, 40)),
                          onPressed: () => onSet(''),
                          child: Text(l.titleRemove),
                        )
                      : FilledButton.tonal(
                          style: compactButton,
                          onPressed: () => onSet(t.code),
                          child: Text(l.titleWear),
                        ),
                ],
              ],
            ),
          ],
        ),
      ),
    );
  }
}

/// "+2 Attack", "+5% Mining yield", "+1/h Mining yield" (DD-220 targets).
String bonusText(AppLocalizations l, BonusEffect b) {
  const gather = 'gather_yield_';
  final isGather = b.target.startsWith(gather);
  final name = isGather
      ? l.gatherYield(activityName(l, b.target.substring(gather.length)))
      : switch (b.target) {
          'attack' => l.attack,
          'defense' => l.defense,
          'crit_pct' => l.crit,
          'hit_pct' => l.hit,
          'str' => l.statStr,
          'agi' => l.statAgi,
          'int' => l.statInt,
          'vit' => l.statVit,
          'cha' => l.statCha,
          'luck' => l.statLuck,
          'hp_max' => l.hp,
          'mp_max' => l.mp,
          _ => humanize(b.target),
        };
  final n = b.amount == b.amount.roundToDouble() ? b.amount.round().toString() : b.amount.toString();
  final value = b.percent ? '+${(b.amount * 100).round()}%' : (isGather ? '+${l.perHour(n)}' : '+$n');
  return '$value $name';
}

class _AchievementTile extends StatelessWidget {
  const _AchievementTile({required this.a});

  final Achievement a;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final rewards = [
      if (a.rewardGold > 0) goldString(a.rewardGold),
      if (a.rewardTitle != null) l.rewardTitle(ContentNames.of('title', a.rewardTitle!)),
    ].join(' · ');
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  a.unlocked ? Icons.emoji_events : Icons.emoji_events_outlined,
                  color: a.unlocked ? Palette.accent : Palette.textMuted,
                  size: 20,
                ),
                const SizedBox(width: Space.sm),
                Expanded(child: Text(ContentNames.of('achievement', a.code), style: text.titleSmall)),
                Text(l.pointsShort(a.points), style: text.bodySmall),
              ],
            ),
            if (!a.unlocked) ...[
              const SizedBox(height: Space.sm),
              StatBar(
                label: ContentNames.of('category', a.category),
                value: a.target == 0 ? 0 : (a.progress / a.target).clamp(0, 1).toDouble(),
                color: Palette.accent,
                valueText: '${a.progress}/${a.target}',
              ),
            ],
            if (a.event.isNotEmpty) ...[
              const SizedBox(height: Space.xs),
              Text(goalText(l, a.event, goalCount(a.event, a.target)), style: text.bodySmall),
            ],
            if (rewards.isNotEmpty) ...[const SizedBox(height: Space.xs), Text(rewards, style: text.bodySmall)],
            if (a.event.isNotEmpty || a.rewardTitle != null) ...[
              const SizedBox(height: Space.sm),
              LinkTags([?eventRef(a.event), if (a.rewardTitle case final t?) CodexRef('title', t)]),
            ],
          ],
        ),
      ),
    );
  }
}
