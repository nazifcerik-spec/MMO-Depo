import 'dart:async';

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/app_update.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/widgets.dart';

/// The update icon in the character list's bar (DD-304): always there on
/// Android builds, with a dot when a newer build waits. A tap checks, asks,
/// downloads (progress in a dialog) and opens the system installer.
class UpdateButton extends StatefulWidget {
  const UpdateButton({super.key, this.updater});

  /// Tests inject one; the app uses [shared].
  final AppUpdate? updater;

  static final shared = AppUpdate();

  /// The run's first check with [shared], reused when the screen returns.
  static Future<Release?>? _first;

  @override
  State<UpdateButton> createState() => _UpdateButtonState();
}

class _UpdateButtonState extends State<UpdateButton> {
  late final _u = widget.updater ?? UpdateButton.shared;
  Release? _release;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    if (!_u.available) return;
    (widget.updater == null ? UpdateButton._first ??= _u.check() : _u.check())
        .then((r) {
          if (mounted) setState(() => _release = r);
        })
        .catchError((Object _) {
          UpdateButton._first = null;
        });
  }

  Future<void> _update() async {
    final l = AppLocalizations.of(context);
    setState(() => _busy = true);
    final progress = ValueNotifier<double>(0);
    var dialog = false;
    try {
      final r = await _u.check();
      if (!mounted) return;
      setState(() => _release = r);
      if (r == null) {
        toast(context, Text(l.updateLatest));
        return;
      }
      final ok = await confirmAction(
        context,
        title: l.updateAvailable(r.version),
        body: l.updateConfirm(r.version, (r.size / 1e6).ceil()),
        action: l.updateButton,
      );
      if (!ok || !mounted) return;
      dialog = true;
      unawaited(
        showDialog<void>(
          context: context,
          barrierDismissible: false,
          builder: (_) => AlertDialog(
            title: Text(l.updateDownloading),
            content: ValueListenableBuilder<double>(
              valueListenable: progress,
              builder: (_, p, _) => LinearProgressIndicator(value: p),
            ),
          ),
        ),
      );
      final started = await _u.install(r, onProgress: (p) => progress.value = p);
      if (!started && mounted) toast(context, Text(l.updatePermission));
    } on Exception {
      if (mounted) toast(context, Text(l.updateFailed));
    } finally {
      if (dialog && mounted) Navigator.of(context, rootNavigator: true).pop();
      progress.dispose();
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_u.available) return const SizedBox.shrink();
    final l = AppLocalizations.of(context);
    return IconButton(
      key: const Key('update-button'),
      tooltip: _release == null ? l.updateButton : l.updateAvailable(_release!.version),
      onPressed: _busy ? null : _update,
      icon: Badge(
        isLabelVisible: _release != null,
        backgroundColor: Palette.success,
        child: const Icon(Icons.system_update_outlined),
      ),
    );
  }
}
