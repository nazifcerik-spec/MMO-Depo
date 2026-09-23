import 'dart:async';

import 'package:flutter/material.dart';

import '../app/theme.dart';
import '../l10n/app_localizations.dart';
import 'widgets.dart';

/// Progress of any timed contract (gathering, hunt, dungeon, travel).
/// Owns the only timer on its screen: a 1s tick that rebuilds just this bar
/// and stops itself at completion. Device clock is display-only; the server
/// decides what was earned.
class ContractProgress extends StatefulWidget {
  const ContractProgress({
    super.key,
    required this.startedAt,
    required this.endsAt,
    required this.doneLabel,
    this.onComplete,
  });

  final DateTime startedAt;
  final DateTime endsAt;
  final String doneLabel;
  final VoidCallback? onComplete;

  static bool isDone(DateTime endsAt) => !DateTime.now().isBefore(endsAt);

  @override
  State<ContractProgress> createState() => _ContractProgressState();
}

class _ContractProgressState extends State<ContractProgress> {
  Timer? _tick;

  @override
  void initState() {
    super.initState();
    if (!ContractProgress.isDone(widget.endsAt)) {
      _tick = Timer.periodic(const Duration(seconds: 1), (_) {
        if (ContractProgress.isDone(widget.endsAt)) {
          _tick?.cancel();
          widget.onComplete?.call();
        }
        setState(() {});
      });
    }
  }

  @override
  void dispose() {
    _tick?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final now = DateTime.now();
    final total = widget.endsAt.difference(widget.startedAt).inMilliseconds;
    final value = total <= 0 ? 1.0 : (now.difference(widget.startedAt).inMilliseconds / total).clamp(0.0, 1.0);
    final done = ContractProgress.isDone(widget.endsAt);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        ClipRRect(
          borderRadius: BorderRadius.circular(Radii.sm),
          child: LinearProgressIndicator(
            value: value,
            minHeight: 10,
            color: done ? Palette.success : Palette.accent,
            backgroundColor: Palette.bgElevated,
          ),
        ),
        const SizedBox(height: Space.sm),
        Text(
          done ? widget.doneLabel : l.remaining(shortDuration(l, widget.endsAt.difference(now))),
          style: TextStyle(
            color: done ? Palette.success : Palette.textSecondary,
            fontSize: 12,
            fontFeatures: numeric,
            fontWeight: done ? FontWeight.w700 : FontWeight.w400,
          ),
        ),
      ],
    );
  }
}

/// Yes/no dialog used before irreversible or lossy actions.
Future<bool> confirmAction(
  BuildContext context, {
  required String title,
  required String body,
  required String action,
  bool danger = false,
}) async {
  final l = AppLocalizations.of(context);
  return await showDialog<bool>(
        context: context,
        builder: (c) => AlertDialog(
          title: Text(title),
          content: Text(body),
          actions: [
            TextButton(onPressed: () => Navigator.pop(c, false), child: Text(l.cancel)),
            FilledButton(
              style: danger ? FilledButton.styleFrom(backgroundColor: Palette.danger) : null,
              onPressed: () => Navigator.pop(c, true),
              child: Text(action),
            ),
          ],
        ),
      ) ??
      false;
}
