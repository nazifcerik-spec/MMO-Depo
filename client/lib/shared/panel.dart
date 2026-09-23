import 'package:flutter/material.dart';

import '../app/theme.dart';
import '../core/api_client.dart';
import '../l10n/app_localizations.dart';
import 'widgets.dart';

/// Shared load/error/busy plumbing for tab panels: [fetch] fills state,
/// [load] wraps it with error capture, [run] wraps a mutation (busy flag +
/// error toast + reload), [frame] renders spinner / error / pull-to-refresh list.
mixin PanelState<T extends StatefulWidget> on State<T> {
  bool loaded = false;
  Object? error;
  bool busy = false;

  Future<void> fetch();

  Future<void> load() async {
    try {
      await fetch();
      if (mounted) {
        setState(() {
          loaded = true;
          error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => error = e);
    }
  }

  Future<void> run(Future<void> Function() action) async {
    setState(() => busy = true);
    try {
      await action();
    } on ApiError catch (e) {
      if (mounted && !e.isSessionError) toast(context, Text(errorText(AppLocalizations.of(context), e)));
      await load();
    } finally {
      if (mounted) setState(() => busy = false);
    }
  }

  Widget frame(List<Widget> Function() children) {
    final l = AppLocalizations.of(context);
    if (!loaded) {
      return error != null
          ? MessageView(message: errorText(l, error!), onRetry: load)
          : const Center(child: CircularProgressIndicator());
    }
    return RefreshIndicator(
      onRefresh: load,
      child: AbsorbPointer(
        absorbing: busy,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
          children: children(),
        ),
      ),
    );
  }
}
