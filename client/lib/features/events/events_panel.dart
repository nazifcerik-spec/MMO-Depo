import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'event_repository.dart';
import '../../shared/content_names.dart';

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

  Future<void> _info(Widget sheet) => showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    backgroundColor: Palette.bgSurface,
    builder: (_) => sheet,
  );

  Future<void> _claim(SeasonalEvent s) => run(() async {
    final n = await widget.repo.claimDaily(s.code);
    if (mounted) toast(context, Text(AppLocalizations.of(context).eventTokensGained(n)));
    await load();
  });

  Future<void> _enlist(WorldEvent e) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: ContentNames.of('worldevent', e.code),
      body: l.enlistHint,
      action: l.enlist,
    );
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
        for (final s in _seasonal)
          _SeasonalCard(
            s: s,
            onClaim: s.claimedToday || busy ? null : () => _claim(s),
            onInfo: () => _info(_SeasonalInfo(s: s)),
          ),
        if (_events.isEmpty)
          Padding(
            padding: const EdgeInsets.only(top: Space.xl),
            child: MessageView(message: l.eventsEmpty, icon: Icons.public),
          ),
        for (final e in _events)
          _WorldEventCard(
            e: e,
            onEnlist: e.open && !e.enlisted && widget.level >= e.minLevel && !busy ? () => _enlist(e) : null,
            onInfo: () => _info(_WorldEventInfo(e: e)),
          ),
      ],
    );
  }
}

class _SeasonalCard extends StatelessWidget {
  const _SeasonalCard({required this.s, this.onClaim, this.onInfo});

  final SeasonalEvent s;
  final VoidCallback? onClaim;
  final VoidCallback? onInfo;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      color: Palette.accent.withValues(alpha: 0.08),
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onInfo,
        child: Padding(
          padding: const EdgeInsets.all(Space.lg),
          child: Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(ContentNames.of('seasonal', s.code), style: text.titleMedium),
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
      ),
    );
  }
}

class _WorldEventCard extends StatelessWidget {
  const _WorldEventCard({required this.e, this.onEnlist, this.onInfo});

  final WorldEvent e;
  final VoidCallback? onEnlist;
  final VoidCallback? onInfo;

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
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onInfo,
        child: Padding(
          padding: const EdgeInsets.all(Space.lg),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(e.kind == 'world_boss' ? Icons.whatshot : Icons.flag_outlined, color: color, size: 20),
                  const SizedBox(width: Space.sm),
                  Expanded(child: Text(ContentNames.of('worldevent', e.code), style: text.titleMedium)),
                  Text(l.levelReq(e.minLevel), style: text.bodySmall),
                  const SizedBox(width: Space.xs),
                  Icon(Icons.info_outline, size: 18, color: Palette.textSecondary, semanticLabel: l.eventDetails),
                ],
              ),
              const SizedBox(height: Space.xs),
              Text(
                [
                  e.kind == 'world_boss' ? l.worldBoss : l.invasion,
                  if (e.region != null) l.regionOnly(ContentNames.of('region', e.region!)),
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
      ),
    );
  }
}

/// Bottom sheet: what to do to benefit from a world event, and what it pays.
class _WorldEventInfo extends StatelessWidget {
  const _WorldEventInfo({required this.e});

  final WorldEvent e;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final cap = e.afkCapTier >= 1 && e.afkCapTier <= e.tiers.length ? e.afkCapTier : e.tiers.length;
    final steps = [
      l.eventStepEnlist(e.minLevel),
      if (e.region != null) l.eventStepRegion(ContentNames.of('region', e.region!)),
      l.eventStepPower(e.powerNow),
      l.eventStepScore(cap),
      l.eventStepGoal,
      l.eventStepPaid,
    ];
    final tiers = e.sortedTiers;
    final estimate = e.estimateScore(DateTime.now());
    final tier = e.tierFor(estimate);
    return _InfoSheet(
      title: ContentNames.of('worldevent', e.code),
      subtitle: e.kind == 'world_boss' ? l.worldBoss : l.invasion,
      children: [
        Text(l.eventHowTitle, style: text.titleSmall),
        const SizedBox(height: Space.sm),
        for (var i = 0; i < steps.length; i++) _Step(n: i + 1, text: steps[i]),
        const SizedBox(height: Space.lg),
        Text(l.eventRewardsTitle, style: text.titleSmall),
        const SizedBox(height: Space.sm),
        for (var i = 0; i < tiers.length; i++)
          Padding(
            padding: const EdgeInsets.only(bottom: Space.xs),
            child: Wrap(
              spacing: Space.sm,
              crossAxisAlignment: WrapCrossAlignment.center,
              children: [
                Text(l.eventTier(i + 1, '${tiers[i].minScore}'), style: const TextStyle(fontWeight: FontWeight.w600)),
                if (tiers[i].goldCopper > 0) GoldText(tiers[i].goldCopper, size: 13),
                if (tiers[i].tokens > 0)
                  Text(l.eventTierTokens(tiers[i].tokens), style: const TextStyle(color: Palette.epic, fontSize: 13)),
              ],
            ),
          ),
        if (e.open && !e.enlisted && e.powerNow > 0) ...[
          const SizedBox(height: Space.md),
          Text(
            tier > 0 ? l.eventEstimate('$estimate', tier) : l.eventEstimateNone('$estimate'),
            style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
          ),
        ],
      ],
    );
  }
}

/// Bottom sheet for a seasonal event: the daily claim and what it gives.
class _SeasonalInfo extends StatelessWidget {
  const _SeasonalInfo({required this.s});

  final SeasonalEvent s;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return _InfoSheet(
      title: ContentNames.of('seasonal', s.code),
      subtitle: l.endsIn(shortDuration(l, s.endsAt.difference(DateTime.now()))),
      children: [
        Text(l.eventHowTitle, style: text.titleSmall),
        const SizedBox(height: Space.sm),
        _Step(n: 1, text: l.seasonalHow(s.dailyTokens)),
        const SizedBox(height: Space.lg),
        Text(l.eventRewardsTitle, style: text.titleSmall),
        const SizedBox(height: Space.sm),
        Text(l.seasonalBenefit(s.balance)),
      ],
    );
  }
}

class _InfoSheet extends StatelessWidget {
  const _InfoSheet({required this.title, required this.subtitle, required this.children});

  final String title;
  final String subtitle;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final text = Theme.of(context).textTheme;
    return SafeArea(
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: MediaQuery.sizeOf(context).height * 0.8),
        child: ListView(
          shrinkWrap: true,
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          children: [
            Text(title, style: text.titleLarge),
            Text(subtitle, style: text.bodySmall),
            const SizedBox(height: Space.lg),
            ...children,
          ],
        ),
      ),
    );
  }
}

class _Step extends StatelessWidget {
  const _Step({required this.n, required this.text});

  final int n;
  final String text;

  @override
  Widget build(BuildContext context) => Padding(
    padding: const EdgeInsets.only(bottom: Space.sm),
    child: Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        CircleAvatar(
          radius: 10,
          backgroundColor: Palette.accent.withValues(alpha: 0.2),
          child: Text('$n', style: const TextStyle(fontSize: 11, color: Palette.accent)),
        ),
        const SizedBox(width: Space.sm),
        Expanded(child: Text(text)),
      ],
    ),
  );
}
