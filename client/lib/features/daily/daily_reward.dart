import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';

/// One day of the seven-day login cycle (DD-311).
class DailyReward {
  const DailyReward({required this.day, required this.goldCopper, this.potion = '', this.potionQty = 0});

  factory DailyReward.fromJson(Map<String, dynamic> j) => DailyReward(
    day: (j['day'] as num?)?.toInt() ?? 0,
    goldCopper: (j['gold_copper'] as num?)?.toInt() ?? 0,
    potion: j['potion'] as String? ?? '',
    potionQty: (j['potion_qty'] as num?)?.toInt() ?? 0,
  );

  final int day;
  final int goldCopper;
  final String potion;
  final int potionQty;
}

class DailyStatus {
  const DailyStatus({required this.claimable, required this.day, required this.streak, required this.rewards});

  factory DailyStatus.fromJson(Map<String, dynamic> j) => DailyStatus(
    claimable: j['claimable'] as bool? ?? false,
    day: (j['day'] as num?)?.toInt() ?? 1,
    streak: (j['streak'] as num?)?.toInt() ?? 0,
    rewards: [
      for (final r in (j['rewards'] as List? ?? const []).cast<Map<String, dynamic>>()) DailyReward.fromJson(r),
    ],
  );

  final bool claimable;

  /// The day the next claim pays (today's when [claimable], else tomorrow's).
  final int day;

  /// Days claimed in a row so far.
  final int streak;
  final List<DailyReward> rewards;
}

class DailyRepository {
  DailyRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  Future<DailyStatus> status() async => DailyStatus.fromJson(await api.get('/v1/character/$characterId/daily-reward'));

  /// Today's reward to this character; returns it and the new status.
  Future<(DailyReward, DailyStatus)> claim() async {
    final j = await api.post('/v1/character/$characterId/daily-reward/claim', const {});
    return (
      DailyReward.fromJson(j['reward'] as Map<String, dynamic>),
      DailyStatus.fromJson(j['status'] as Map<String, dynamic>),
    );
  }
}

/// The reward window: seven day boxes, today's highlighted, and the claim
/// button. Returns true when a reward was claimed.
Future<bool> showDailyReward(BuildContext context, DailyRepository repo, DailyStatus status) async =>
    await showDialog<bool>(
      context: context,
      builder: (_) => _DailyDialog(repo: repo, status: status),
    ) ??
    false;

class _DailyDialog extends StatefulWidget {
  const _DailyDialog({required this.repo, required this.status});

  final DailyRepository repo;
  final DailyStatus status;

  @override
  State<_DailyDialog> createState() => _DailyDialogState();
}

class _DailyDialogState extends State<_DailyDialog> {
  late DailyStatus _st = widget.status;
  bool _busy = false;
  bool _claimed = false;

  Future<void> _claim() async {
    final l = AppLocalizations.of(context);
    setState(() => _busy = true);
    try {
      final (_, st) = await widget.repo.claim();
      if (!mounted) return;
      setState(() {
        _st = st;
        _claimed = true;
      });
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    // The day just paid (after a claim) or about to be.
    final current = _st.claimable ? _st.day : (_st.streak == 0 ? 0 : _st.streak);
    return AlertDialog(
      title: Text(l.dailyTitle),
      content: SizedBox(
        width: 360,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(l.dailyBody, style: text.bodySmall),
            const SizedBox(height: Space.md),
            Wrap(
              spacing: Space.sm,
              runSpacing: Space.sm,
              children: [
                for (final r in _st.rewards)
                  _DayBox(reward: r, done: _st.claimable ? r.day < _st.day : r.day <= current, today: r.day == current),
              ],
            ),
          ],
        ),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.of(context).pop(_claimed), child: Text(l.close)),
        FilledButton(
          key: const Key('daily-claim'),
          onPressed: _st.claimable && !_busy ? _claim : null,
          child: Text(_st.claimable ? l.dailyClaim : l.dailyTomorrow),
        ),
      ],
    );
  }
}

class _DayBox extends StatelessWidget {
  const _DayBox({required this.reward, required this.done, required this.today});

  final DailyReward reward;
  final bool done;
  final bool today;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Container(
      width: 96,
      padding: const EdgeInsets.all(Space.sm),
      decoration: BoxDecoration(
        color: today ? Palette.accent.withValues(alpha: 0.12) : Palette.bgElevated,
        borderRadius: BorderRadius.circular(Radii.md),
        border: Border.all(color: today ? Palette.accent : Palette.borderSubtle),
      ),
      child: Column(
        children: [
          FittedBox(
            fit: BoxFit.scaleDown,
            child: Row(
              children: [
                Text(l.dailyDay(reward.day), style: text.labelMedium),
                if (done) ...[
                  const SizedBox(width: 4),
                  const Icon(Icons.check_circle, size: 14, color: Palette.success),
                ],
              ],
            ),
          ),
          const SizedBox(height: Space.xs),
          FittedBox(fit: BoxFit.scaleDown, child: GoldText(reward.goldCopper, size: 12)),
          if (reward.potion.isNotEmpty)
            Text(
              '${reward.potionQty}× ${ContentNames.of('item', reward.potion)}',
              style: text.bodySmall,
              textAlign: TextAlign.center,
              maxLines: 2,
            ),
        ],
      ),
    );
  }
}
