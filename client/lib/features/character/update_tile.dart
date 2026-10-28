import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/app_update.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/widgets.dart';

/// Version line + Update button under the character list (DD-304). Checks the
/// feed once per run; a tap checks again, asks, downloads and opens the
/// system installer. Hidden where in-app updates do not apply.
class UpdateTile extends StatefulWidget {
  const UpdateTile({super.key, this.updater});

  /// Tests inject one; the app uses [shared].
  final AppUpdate? updater;

  static final shared = AppUpdate();

  /// The run's first check, shared by every sheet.
  static Future<Release?>? _first;

  @override
  State<UpdateTile> createState() => _UpdateTileState();
}

class _UpdateTileState extends State<UpdateTile> {
  Release? _release;
  double? _progress; // non-null while downloading
  bool _busy = false;
  late final _u = widget.updater ?? UpdateTile.shared;

  @override
  void initState() {
    super.initState();
    if (!_u.available) return;
    (UpdateTile._first ??= _u.check())
        .then((r) {
          if (mounted) setState(() => _release = r);
        })
        .catchError((Object _) {
          UpdateTile._first = null; // offline: try again next time
        });
  }

  Future<void> _update() async {
    final l = AppLocalizations.of(context);
    setState(() => _busy = true);
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
      setState(() => _progress = 0);
      final started = await _u.install(r, onProgress: (p) => mounted ? setState(() => _progress = p) : null);
      if (!started && mounted) toast(context, Text(l.updatePermission));
    } on Exception {
      if (mounted) toast(context, Text(l.updateFailed));
    } finally {
      if (mounted) {
        setState(() {
          _busy = false;
          _progress = null;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_u.available) return const SizedBox.shrink();
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final r = _release;
    return Padding(
      padding: const EdgeInsets.only(top: Space.lg),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  r == null ? l.appVersionLine(_u.current) : l.updateAvailable(r.version),
                  style: text.bodySmall?.copyWith(color: r == null ? Palette.textMuted : Palette.success),
                ),
              ),
              FilledButton.tonalIcon(
                onPressed: _busy ? null : _update,
                icon: const Icon(Icons.system_update_outlined, size: 18),
                label: Text(l.updateButton),
              ),
            ],
          ),
          if (_progress case final p?) ...[
            const SizedBox(height: Space.sm),
            Text(l.updateDownloading, style: text.bodySmall),
            const SizedBox(height: Space.xs),
            LinearProgressIndicator(value: p),
          ],
        ],
      ),
    );
  }
}
