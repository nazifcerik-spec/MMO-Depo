import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
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
  List<OwnedTitle> _titles = const [];

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
    _titles = r[2] as List<OwnedTitle>;
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
          child: ListTile(
            title: Text(ContentNames.of('faction', r.faction)),
            subtitle: Text('${ContentNames.of('reptier', r.tierName)} · ${r.rep}'),
            trailing: Text(l.repDailyLeft(r.dailyLeft), style: Theme.of(context).textTheme.bodySmall),
          ),
        ),
    ];
  }

  List<Widget> _titleRows() {
    final l = AppLocalizations.of(context);
    final active = _titles.where((t) => t.active).map((t) => t.code).firstOrNull ?? '';
    return [
      RadioGroup<String>(
        groupValue: active,
        onChanged: (v) {
          if (v != null && v != active) _setTitle(v);
        },
        child: Column(
          children: [
            RadioListTile<String>(value: '', title: Text(l.noTitle)),
            for (final t in _titles)
              RadioListTile<String>(value: t.code, title: Text(ContentNames.of('title', t.code))),
          ],
        ),
      ),
      if (_titles.isEmpty)
        Padding(
          padding: const EdgeInsets.all(Space.lg),
          child: Text(l.noTitlesYet, style: Theme.of(context).textTheme.bodySmall),
        ),
    ];
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
            if (rewards.isNotEmpty) ...[const SizedBox(height: Space.xs), Text(rewards, style: text.bodySmall)],
          ],
        ),
      ),
    );
  }
}
