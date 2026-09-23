import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'quest_repository.dart';

/// "main_04_wolves" → "Wolves", "daily_miner" → "Miner" (content names are
/// derived from codes until content localisation lands).
String questName(String code) => humanize(code.replaceFirst(RegExp(r'^(main_\d+|side|daily)_'), ''));

/// Objective line from its progression event kind (DD-221 grammar).
String objectiveText(AppLocalizations l, String event) {
  final i = event.indexOf(':');
  final kind = i < 0 ? event : event.substring(0, i);
  final target = i < 0 ? null : event.substring(i + 1);
  if (kind.startsWith('gather_')) return l.objGather(activityName(l, kind.substring(7)));
  return switch (kind) {
    'mob_killed' => target == null ? l.objKillAny : l.objKill(humanize(target)),
    'boss_killed' => l.objBoss,
    'dungeon_cleared' => target == null ? l.objDungeonAny : l.objDungeon(humanize(target)),
    'item_crafted' => target == null ? l.objCraftAny : l.objCraft(itemName(target)),
    _ => humanize(event),
  };
}

/// Quest log: open quests first (turn-in ready on top), then what can be
/// accepted grouped by story / side / daily. Bumps [dataVersion] after a
/// turn-in so the sheet picks up XP and gold.
class QuestsScreen extends StatefulWidget {
  const QuestsScreen({super.key, required this.repo, required this.dataVersion});

  final QuestRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<QuestsScreen> createState() => _QuestsScreenState();
}

class _QuestsScreenState extends State<QuestsScreen> with PanelState {
  QuestBoard? _board;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _board = await widget.repo.board();

  Future<void> _mutate(Future<void> Function() action) => run(() async {
    await action();
    await load();
  });

  Future<void> _abandon(Quest q) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.abandonQuestTitle,
      body: q.category == 'daily' ? l.abandonDailyBody : l.abandonQuestBody,
      action: l.abandon,
      danger: true,
    );
    if (ok) await _mutate(() => widget.repo.abandon(q.code));
  }

  Future<void> _turnIn(Quest q) => _mutate(() async {
    final r = await widget.repo.turnIn(q.code);
    widget.dataVersion.value++;
    if (!mounted) return;
    final l = AppLocalizations.of(context);
    final parts = [if (r.xp > 0) l.xpGain(r.xp), if (r.gold > 0) '+${goldString(r.gold)}'];
    toast(
      context,
      Text([l.questRewarded(parts.join(' · ')), if (r.newLevel != null) l.levelUp(r.newLevel!)].join('\n')),
    );
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Scaffold(
      appBar: AppBar(title: Text(l.quests)),
      body: SafeArea(
        child: PageFrame(
          child: frame(() {
            final b = _board!;
            final open = [...b.open]
              ..sort((a, c) => (a.state == 'completed' ? 0 : 1) - (c.state == 'completed' ? 0 : 1));
            final byCat = {
              for (final cat in const ['main', 'side', 'daily'])
                cat: b.available.where((q) => q.category == cat).toList(),
            };
            final text = Theme.of(context).textTheme;
            Widget header(String s, [String? trailing]) => Padding(
              padding: const EdgeInsets.fromLTRB(0, Space.lg, 0, Space.sm),
              child: Row(
                children: [
                  Expanded(
                    child: Text(s, style: text.titleSmall?.copyWith(color: Palette.accent)),
                  ),
                  if (trailing != null) Text(trailing, style: text.bodySmall),
                ],
              ),
            );
            return [
              header(l.questsInProgress),
              if (open.isEmpty) Text(l.questsNoneOpen, style: text.bodySmall),
              for (final q in open)
                _QuestCard(
                  q: q,
                  onPrimary: q.state == 'completed' ? () => _turnIn(q) : null,
                  primaryLabel: l.turnIn,
                  onAbandon: () => _abandon(q),
                ),
              for (final (cat, title) in [('main', l.questsStory), ('side', l.questsSide), ('daily', l.questsDaily)])
                if (byCat[cat]!.isNotEmpty) ...[
                  header(
                    title,
                    cat == 'daily' && b.dailyResetsAt != null
                        ? l.resetsIn(shortDuration(l, b.dailyResetsAt!.difference(DateTime.now())))
                        : null,
                  ),
                  for (final q in byCat[cat]!)
                    _QuestCard(
                      q: q,
                      onPrimary: () => _mutate(() => widget.repo.accept(q.code)),
                      primaryLabel: l.accept,
                    ),
                ],
            ];
          }),
        ),
      ),
    );
  }
}

class _QuestCard extends StatelessWidget {
  const _QuestCard({required this.q, required this.primaryLabel, this.onPrimary, this.onAbandon});

  final Quest q;
  final String primaryLabel;
  final VoidCallback? onPrimary;
  final VoidCallback? onAbandon;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final ready = q.state == 'completed';
    final open = q.state != 'available';
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      shape: ready
          ? RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: const BorderSide(color: Palette.success),
            )
          : null,
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  ready ? Icons.task_alt : Icons.assignment_outlined,
                  size: 20,
                  color: ready ? Palette.success : Palette.textMuted,
                ),
                const SizedBox(width: Space.sm),
                Expanded(child: Text(questName(q.code), style: text.titleMedium)),
                if (onAbandon != null)
                  IconButton(
                    tooltip: l.abandon,
                    visualDensity: VisualDensity.compact,
                    icon: const Icon(Icons.close, size: 18),
                    onPressed: onAbandon,
                  ),
              ],
            ),
            for (final o in q.objectives)
              Padding(
                padding: const EdgeInsets.only(top: Space.sm),
                child: open
                    ? StatBar(
                        label: objectiveText(l, o.event),
                        value: (o.progress / o.count).clamp(0, 1).toDouble(),
                        color: o.progress >= o.count ? Palette.success : Palette.accent,
                        valueText: '${o.progress}/${o.count}',
                      )
                    : Text('• ${objectiveText(l, o.event)} ×${o.count}', style: text.bodyMedium),
              ),
            const SizedBox(height: Space.sm),
            Row(
              children: [
                Expanded(
                  child: Text(
                    [
                      if (q.rewardXp > 0) l.xpGain(q.rewardXp),
                      if (q.rewardGold > 0) goldString(q.rewardGold),
                    ].join(' · '),
                    style: text.bodySmall,
                  ),
                ),
                if (!open || ready) FilledButton(style: compactButton, onPressed: onPrimary, child: Text(primaryLabel)),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
