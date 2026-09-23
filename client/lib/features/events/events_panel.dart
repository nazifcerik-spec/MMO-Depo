import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'event_repository.dart';

/// Adventure → Events: seasonal daily token claim + world boss / invasion
/// spawns (enlist snapshots power; score accrues per minute server-side).
class EventsPanel extends StatefulWidget {
  const EventsPanel({super.key, required this.repo, required this.level});

  final EventRepository repo;
  final int level;

  @override
  State<EventsPanel> createState() => _EventsPanelState();
}

class _EventsPanelState extends State<EventsPanel> with PanelState {
  List<SeasonalEvent> _seasonal = const [];
  List<WorldEvent> _events = const [];

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait([widget.repo.seasonal(), widget.repo.worldEvents()]);
    _seasonal = r[0] as List<SeasonalEvent>;
    _events = r[1] as List<WorldEvent>;
  }

  Future<void> _claim(SeasonalEvent s) => run(() async {
    final n = await widget.repo.claimDaily(s.code);
    if (mounted) toast(context, Text(AppLocalizations.of(context).eventTokensGained(n)));
    await load();
  });

  Future<void> _enlist(WorldEvent e) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(context, title: humanize(e.code), body: l.enlistHint, action: l.enlist);
    if (ok) {
      await run(() async {
        await widget.repo.enlist(e.id);
        await load();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return frame(
      () => [
        for (final s in _seasonal) _SeasonalCard(s: s, onClaim: s.claimedToday || busy ? null : () => _claim(s)),
        if (_events.isEmpty)
          Padding(
            padding: const EdgeInsets.only(top: Space.xl),
            child: MessageView(message: l.eventsEmpty, icon: Icons.public),
          ),
        for (final e in _events)
          _WorldEventCard(
            e: e,
            onEnlist: e.open && !e.enlisted && widget.level >= e.minLevel && !busy ? () => _enlist(e) : null,
          ),
      ],
    );
  }
}

class _SeasonalCard extends StatelessWidget {
  const _SeasonalCard({required this.s, this.onClaim});

  final SeasonalEvent s;
  final VoidCallback? onClaim;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      color: Palette.accent.withValues(alpha: 0.08),
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(humanize(s.code), style: text.titleMedium),
                  Text(l.endsIn(shortDuration(l, s.endsAt.difference(DateTime.now()))), style: text.bodySmall),
                  Text(l.tokenBalance(s.balance), style: text.bodySmall),
                ],
              ),
            ),
            FilledButton(
              style: compactButton,
              onPressed: onClaim,
              child: Text(s.claimedToday ? l.claimedToday : l.claimDaily(s.dailyTokens)),
            ),
          ],
        ),
      ),
    );
  }
}

class _WorldEventCard extends StatelessWidget {
  const _WorldEventCard({required this.e, this.onEnlist});

  final WorldEvent e;
  final VoidCallback? onEnlist;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final now = DateTime.now();
    final (status, color) = switch (e.state) {
      'scheduled' => (l.startsIn(shortDuration(l, e.startsAt.difference(now))), Palette.info),
      'active' => (l.endsIn(shortDuration(l, e.endsAt.difference(now))), Palette.warning),
      'succeeded' => (l.eventSucceeded, Palette.success),
      _ => (l.eventFailed, Palette.danger),
    };
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(e.kind == 'world_boss' ? Icons.whatshot : Icons.flag_outlined, color: color, size: 20),
                const SizedBox(width: Space.sm),
                Expanded(child: Text(humanize(e.code), style: text.titleMedium)),
                Text(l.levelReq(e.minLevel), style: text.bodySmall),
              ],
            ),
            const SizedBox(height: Space.xs),
            Text(
              [
                e.kind == 'world_boss' ? l.worldBoss : l.invasion,
                if (e.region != null) l.regionOnly(humanize(e.region!)),
                status,
              ].join(' · '),
              style: text.bodySmall?.copyWith(color: color),
            ),
            const SizedBox(height: Space.sm),
            StatBar(label: l.eventGoal, value: e.progress, color: color, valueText: '${e.total}/${e.goal}'),
            if (e.enlisted)
              Text(
                e.myTier != null ? l.myContribution(e.myScore ?? 0, e.myTier!) : l.enlistedPower(e.myPower ?? 0),
                style: const TextStyle(color: Palette.success, fontSize: 13),
              )
            else if (e.open)
              Align(
                alignment: Alignment.centerRight,
                child: FilledButton(style: compactButton, onPressed: onEnlist, child: Text(l.enlist)),
              ),
          ],
        ),
      ),
    );
  }
}
