import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/config.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/loot_sheet.dart';
import '../../shared/widgets.dart';
import 'activity_models.dart';
import 'activity_repository.dart';

IconData activityIcon(String kind) => switch (kind) {
  'mining' => Icons.construction,
  'fishing' => Icons.phishing,
  'woodcutting' => Icons.forest,
  'herbalism' => Icons.eco,
  'hunting' => Icons.pets,
  _ => Icons.help_outline,
};

/// Gathering loop: pick kind + duration → progress → collect (reveal).
/// [dataVersion] is bumped after a reveal so other tabs (sheet, bag) reload.
class ActivityView extends StatefulWidget {
  const ActivityView({super.key, required this.repo, required this.characterId, required this.dataVersion});

  final ActivityRepository repo;
  final String characterId;
  final ValueNotifier<int> dataVersion;

  @override
  State<ActivityView> createState() => _ActivityViewState();
}

class _ActivityViewState extends State<ActivityView> {
  ActiveActivity? _active;
  bool _loaded = false;
  Object? _error;
  bool _busy = false;

  String _kind = activityKinds.first;
  Duration _duration = activityDurations[2];

  List<Duration> get _durations => [if (AppConfig.flavor == 'dev') const Duration(minutes: 1), ...activityDurations];

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final a = await widget.repo.active(widget.characterId);
      if (mounted) {
        setState(() {
          _active = a;
          _loaded = true;
          _error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _run(Future<void> Function() action) async {
    setState(() => _busy = true);
    try {
      await action();
    } on ApiError catch (e) {
      if (mounted) {
        if (!e.isSessionError) toast(context, Text(errorText(AppLocalizations.of(context), e)));
        await _load(); // state may have moved on (e.g. revealed elsewhere)
      }
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _start() => _run(() async {
    final a = await widget.repo.start(widget.characterId, _kind, _duration);
    if (mounted) setState(() => _active = a);
  });

  Future<void> _collect(ActiveActivity a) async {
    final l = AppLocalizations.of(context);
    if (!a.contract.completeAt(DateTime.now()) &&
        !await confirmAction(context, title: l.collectEarlyTitle, body: l.collectEarlyBody, action: l.collect)) {
      return;
    }
    await _run(() async {
      final y = await widget.repo.reveal(a.contract.id);
      if (!mounted) return;
      setState(() => _active = null);
      widget.dataVersion.value++;
      await showLootSheet(context, title: activityName(l, a.contract.kind), loot: y.toLoot());
    });
  }

  Future<void> _abandon(ActiveActivity a) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.abandonTitle,
      body: l.abandonBody,
      action: l.abandon,
      danger: true,
    );
    if (!ok) return;
    await _run(() async {
      await widget.repo.stop(a.contract.id);
      if (mounted) setState(() => _active = null);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    if (!_loaded) {
      return _error != null
          ? MessageView(message: errorText(l, _error!), onRetry: _load)
          : const Center(child: CircularProgressIndicator());
    }
    final a = _active;
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
        children: a == null ? _picker(l) : _running(l, a),
      ),
    );
  }

  List<Widget> _picker(AppLocalizations l) {
    final text = Theme.of(context).textTheme;
    return [
      Text(l.activityChoose, style: text.titleMedium),
      const SizedBox(height: Space.md),
      GridView.count(
        crossAxisCount: 3,
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        mainAxisSpacing: Space.sm,
        crossAxisSpacing: Space.sm,
        childAspectRatio: 1.05,
        children: [
          for (final k in activityKinds)
            _KindTile(
              kind: k,
              label: activityName(l, k),
              selected: k == _kind,
              onTap: _busy ? null : () => setState(() => _kind = k),
            ),
        ],
      ),
      const SizedBox(height: Space.xl),
      Text(l.durationLabel, style: text.bodySmall),
      const SizedBox(height: Space.sm),
      Wrap(
        spacing: Space.sm,
        runSpacing: Space.sm,
        children: [
          for (final d in _durations)
            ChoiceChip(
              label: Text(shortDuration(l, d)),
              selected: d == _duration,
              onSelected: _busy ? null : (_) => setState(() => _duration = d),
              selectedColor: Palette.accent.withValues(alpha: 0.2),
              side: BorderSide(color: d == _duration ? Palette.accent : Palette.borderSubtle),
            ),
        ],
      ),
      const SizedBox(height: Space.xl),
      FilledButton.icon(
        onPressed: _busy ? null : _start,
        icon: Icon(activityIcon(_kind)),
        label: Text(l.startActivity(activityName(l, _kind))),
      ),
    ];
  }

  List<Widget> _running(AppLocalizations l, ActiveActivity a) {
    final text = Theme.of(context).textTheme;
    final y = a.current;
    return [
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
                    backgroundColor: Palette.accent.withValues(alpha: 0.15),
                    child: Icon(activityIcon(a.contract.kind), color: Palette.accent),
                  ),
                  const SizedBox(width: Space.md),
                  Expanded(child: Text(activityName(l, a.contract.kind), style: text.titleMedium)),
                ],
              ),
              const SizedBox(height: Space.lg),
              ContractProgress(
                startedAt: a.contract.startedAt,
                endsAt: a.contract.expiresAt,
                doneLabel: l.activityDone,
                onComplete: () => setState(() {}),
              ),
              const SizedBox(height: Space.lg),
              Text(l.soFar, style: text.bodySmall),
              const SizedBox(height: Space.xs),
              Wrap(
                spacing: Space.lg,
                runSpacing: Space.xs,
                crossAxisAlignment: WrapCrossAlignment.center,
                children: [
                  GoldText(y.goldCopper),
                  for (final (r, n) in [('common', y.common), ('rare', y.rare), ('legendary', y.legendary)])
                    if (n > 0)
                      Text(
                        '${rarityName(l, r)} × $n',
                        style: TextStyle(color: rarityColor(r), fontWeight: FontWeight.w600),
                      ),
                ],
              ),
            ],
          ),
        ),
      ),
      const SizedBox(height: Space.lg),
      FilledButton(
        onPressed: _busy ? null : () => _collect(a),
        child: _busy
            ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
            : Text(l.collect),
      ),
      const SizedBox(height: Space.sm),
      TextButton(
        onPressed: _busy ? null : () => _abandon(a),
        style: TextButton.styleFrom(foregroundColor: Palette.danger),
        child: Text(l.abandon),
      ),
    ];
  }
}

class _KindTile extends StatelessWidget {
  const _KindTile({required this.kind, required this.label, required this.selected, this.onTap});

  final String kind;
  final String label;
  final bool selected;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      selected: selected,
      button: true,
      child: Material(
        color: selected ? Palette.accent.withValues(alpha: 0.12) : Palette.bgSurface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(Radii.md),
          side: BorderSide(color: selected ? Palette.accent : Palette.borderSubtle, width: selected ? 1.5 : 1),
        ),
        child: InkWell(
          borderRadius: BorderRadius.circular(Radii.md),
          onTap: onTap,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(activityIcon(kind), color: selected ? Palette.accent : Palette.textSecondary, size: 28),
              const SizedBox(height: Space.sm),
              Text(label, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
            ],
          ),
        ),
      ),
    );
  }
}
