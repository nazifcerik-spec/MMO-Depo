part of 'adventure_view.dart';

// --- Raids (DD-261) -------------------------------------------------------------

/// Party raids: the leader opens a lobby, members join, the leader starts it;
/// once it has run every member collects their own rewards.
class _RaidPanel extends StatefulWidget {
  const _RaidPanel({required this.repo, required this.dataVersion, required this.level, this.live});

  final AdventureRepository repo;
  final ValueNotifier<int> dataVersion;
  final int level;

  /// Lobby changes of the party arrive as `raid.changed`.
  final LiveEvents? live;

  @override
  State<_RaidPanel> createState() => _RaidPanelState();
}

class _RaidPanelState extends State<_RaidPanel> with PanelState {
  List<RaidInfo> _raids = const [];
  RaidCurrent _cur = const RaidCurrent(inParty: false, isLeader: false, lockouts: {});
  StreamSubscription<LiveEvent>? _sub;

  String get _me => widget.repo.characterId;

  @override
  void initState() {
    super.initState();
    load();
    _sub = widget.live?.events.where((e) => e.kind == 'raid.changed').listen((_) => load());
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object>([widget.repo.raids(), widget.repo.raid()]);
    _raids = r[0] as List<RaidInfo>;
    _cur = r[1] as RaidCurrent;
  }

  RaidInfo? _info(String code) => _raids.where((r) => r.code == code).firstOrNull;

  Future<void> _act(Future<void> Function() action) => run(() async {
    await action();
    await load();
  });

  Future<void> _confirm(
    String title,
    String body,
    String action,
    Future<void> Function() fn, {
    bool danger = false,
  }) async {
    if (await confirmAction(context, title: title, body: body, action: action, danger: danger)) await _act(fn);
  }

  Future<void> _collect(RaidRun r) => run(() async {
    final l = AppLocalizations.of(context);
    final res = await widget.repo.revealRaid(r.id);
    if (!mounted) return;
    widget.dataVersion.value++;
    await showLootSheet(
      context,
      title: ContentNames.of('dungeon', r.code),
      subtitle: _resultSubtitle(l, res),
      loot: res.loot,
    );
    await load();
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(() {
      final active = _cur.raid ?? _cur.lobby;
      return [
        if (!_cur.inParty && _cur.raid == null)
          Padding(
            padding: const EdgeInsets.only(bottom: Space.md),
            child: Text(l.raidNeedParty, style: text.bodySmall),
          )
        else if (active == null && !_cur.isLeader)
          Padding(
            padding: const EdgeInsets.only(bottom: Space.md),
            child: Text(l.raidLeaderOnly, style: text.bodySmall),
          ),
        if (active != null) ...[_runCard(l, text, active), const SizedBox(height: Space.lg)],
        for (final r in _raids) ...[
          _raidCard(l, text, r, canOpen: active == null && _cur.isLeader),
          const SizedBox(height: Space.md),
        ],
      ];
    });
  }

  Widget _raidCard(AppLocalizations l, TextTheme text, RaidInfo r, {required bool canOpen}) {
    final d = r.dungeon;
    final locked = _cur.lockouts.contains(r.code);
    final low = d.minLevel > widget.level;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(Space.lg),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Expanded(child: Text(ContentNames.of('dungeon', r.code), style: text.titleMedium)),
                Text(
                  l.levelReq(d.minLevel),
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w700,
                    color: low ? Palette.danger : Palette.textSecondary,
                  ),
                ),
              ],
            ),
            const SizedBox(height: Space.xs),
            Text('${l.raidGroup(r.groupMin, r.groupMax)} · ${shortDuration(l, d.duration)}', style: text.bodySmall),
            Text(l.dungeonEncounters(d.encounters, ContentNames.of('mob', d.bossMobCode)), style: text.bodySmall),
            const SizedBox(height: Space.sm),
            Wrap(
              spacing: Space.sm,
              crossAxisAlignment: WrapCrossAlignment.center,
              children: [
                Text(l.bonus, style: text.bodySmall),
                GoldText(d.bonusGold, size: 13),
                Text(
                  l.xpGain(d.bonusXp),
                  style: const TextStyle(fontSize: 12, color: Palette.epic, fontWeight: FontWeight.w600),
                ),
              ],
            ),
            Text(
              l.raidBossGear(r.bossGear),
              style: const TextStyle(fontSize: 12, color: Palette.accent, fontWeight: FontWeight.w600),
            ),
            if (locked) ...[
              const SizedBox(height: Space.sm),
              Text(
                l.raidLocked,
                style: const TextStyle(fontSize: 12, color: Palette.success, fontWeight: FontWeight.w700),
              ),
            ],
            if (canOpen) ...[
              const SizedBox(height: Space.md),
              FilledButton.tonal(
                onPressed: low || locked ? null : () => _act(() => widget.repo.openRaid(r.code)),
                child: Text(l.raidOpen),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _runCard(AppLocalizations l, TextTheme text, RaidRun r) {
    final info = _info(r.code);
    final mine = r.members.any((m) => m.characterId == _me);
    final leader = r.leaderId == _me;
    final enough = r.members.length >= r.groupMin;
    final low = info != null && info.dungeon.minLevel > widget.level;
    final locked = _cur.lockouts.contains(r.code);
    final done = r.running && r.expiresAt != null && ContractProgress.isDone(r.expiresAt!);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(Space.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      radius: 22,
                      backgroundColor: Palette.danger.withValues(alpha: 0.15),
                      child: const Icon(Icons.groups, color: Palette.danger),
                    ),
                    const SizedBox(width: Space.md),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(ContentNames.of('dungeon', r.code), style: text.titleMedium),
                          if (r.forming) Text(l.raidLobby, style: text.bodySmall),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: Space.md),
                Text(l.raidMembers(r.members.length, r.groupMax), style: text.labelLarge),
                for (final m in r.members)
                  Padding(
                    padding: const EdgeInsets.only(top: Space.xs),
                    child: Row(
                      children: [
                        Icon(
                          m.characterId == r.leaderId ? Icons.star : Icons.person_outline,
                          size: 16,
                          color: m.characterId == r.leaderId ? Palette.accent : Palette.textSecondary,
                        ),
                        const SizedBox(width: Space.sm),
                        Expanded(
                          child: Text(
                            '${m.name} · ${l.levelShort(m.level)} · ${className(l, m.classCode)}',
                            style: text.bodyMedium,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                        if (m.collected)
                          Text(l.raidCollected, style: const TextStyle(fontSize: 12, color: Palette.success)),
                      ],
                    ),
                  ),
                if (r.forming && !enough) ...[
                  const SizedBox(height: Space.sm),
                  Text(l.raidNeedMore(r.groupMin), style: text.bodySmall),
                ],
                if (r.running && r.startedAt != null && r.expiresAt != null) ...[
                  const SizedBox(height: Space.lg),
                  ContractProgress(
                    startedAt: r.startedAt!,
                    endsAt: r.expiresAt!,
                    doneLabel: l.raidDone,
                    onComplete: () => setState(() {}),
                  ),
                  if (!done) ...[const SizedBox(height: Space.sm), Text(l.raidWaitHint, style: text.bodySmall)],
                ],
              ],
            ),
          ),
        ),
        const SizedBox(height: Space.md),
        if (r.running && mine)
          FilledButton(onPressed: done ? () => _collect(r) : null, child: Text(l.collect))
        else if (r.forming && !mine)
          FilledButton(
            onPressed: low || locked || r.members.length >= r.groupMax
                ? null
                : () => _act(() => widget.repo.joinRaid(r.id)),
            child: Text(l.raidJoin),
          )
        else if (r.forming && leader) ...[
          FilledButton(
            onPressed: enough
                ? () => _confirm(l.raidStart, l.raidStartBody, l.raidStart, () => widget.repo.startRaid(r.id))
                : null,
            child: Text(l.raidStart),
          ),
          TextButton(
            onPressed: () =>
                _confirm(l.raidCancel, l.raidCancelBody, l.raidCancel, () => widget.repo.leaveRaid(r.id), danger: true),
            style: TextButton.styleFrom(foregroundColor: Palette.danger),
            child: Text(l.raidCancel),
          ),
        ] else if (r.forming)
          TextButton(
            onPressed: () => _act(() => widget.repo.leaveRaid(r.id)),
            style: TextButton.styleFrom(foregroundColor: Palette.danger),
            child: Text(l.raidLeave),
          ),
      ],
    );
  }
}
