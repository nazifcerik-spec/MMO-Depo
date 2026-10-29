import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../auth/account_dialogs.dart';
import '../auth/auth_controller.dart';
import 'character_repository.dart';
import '../home/home_shell.dart';
import 'create_character_screen.dart';
import 'models.dart';
import 'update_button.dart';

class CharacterSelectScreen extends StatefulWidget {
  const CharacterSelectScreen({super.key, required this.auth, required this.repo});

  final AuthController auth;
  final CharacterRepository repo;

  @override
  State<CharacterSelectScreen> createState() => _CharacterSelectScreenState();
}

class _CharacterSelectScreenState extends State<CharacterSelectScreen> {
  CharacterList? _list;
  Object? _error;
  String? _entering;

  @override
  void initState() {
    super.initState();
    _load();
    if (widget.auth.api.deletionCancelled) {
      widget.auth.api.deletionCancelled = false;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) toast(context, Text(AppLocalizations.of(context).deletionCancelled));
      });
    }
  }

  Future<void> _load() async {
    setState(() => _error = null);
    try {
      final list = await widget.repo.list();
      if (mounted) setState(() => _list = list);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _create() async {
    final created = await Navigator.of(context)
        .push<CharacterSummary>(MaterialPageRoute(builder: (_) => CreateCharacterScreen(repo: widget.repo)));
    if (created != null) await _load();
  }

  Future<void> _enter(CharacterSummary c) async {
    final l = AppLocalizations.of(context);
    setState(() => _entering = c.id);
    try {
      final sheet = await widget.repo.enter(c.id);
      if (!mounted) return;
      await Navigator.of(context).push(
        MaterialPageRoute(
          builder: (_) => HomeShell(auth: widget.auth, characters: widget.repo, initial: sheet),
        ),
      );
      if (mounted) await _load(); // level/state may have changed while playing
    } on ApiError catch (e) {
      if (mounted && widget.auth.status == AuthStatus.signedIn) {
        toast(context, Text(errorText(l, e)));
      }
    } finally {
      if (mounted) setState(() => _entering = null);
    }
  }

  Future<void> _delete(CharacterSummary c) async {
    final l = AppLocalizations.of(context);
    final done = await showDialog<bool>(
      context: context,
      builder: (_) => _DeleteCharacterDialog(c: c, repo: widget.repo),
    );
    if (done != true || !mounted) return;
    toast(context, Text(l.characterDeleted(c.name)));
    await _load();
  }

  Future<void> _cancelDelete(CharacterSummary c) async {
    final l = AppLocalizations.of(context);
    try {
      await widget.repo.cancelDelete(c.id);
      if (!mounted) return;
      toast(context, Text(l.characterDeleteCancelled(c.name)));
      await _load();
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final list = _list;
    return Scaffold(
      appBar: AppBar(
        leading: Padding(
          padding: const EdgeInsets.all(Space.xs),
          child: Semantics(
            label: l.appTitle,
            image: true,
            child: Image.asset(brandLogo, cacheHeight: (48 * MediaQuery.devicePixelRatioOf(context)).round()),
          ),
        ),
        title: Text(l.characters),
        actions: [
          const UpdateButton(), // in-app updates (DD-304)
          AccountMenuButton(auth: widget.auth),
          IconButton(tooltip: l.logout, icon: const Icon(Icons.logout), onPressed: widget.auth.logout),
        ],
      ),
      floatingActionButton: list == null || list.full
          ? null
          : FloatingActionButton.extended(onPressed: _create, icon: const Icon(Icons.add), label: Text(l.newCharacter)),
      body: SafeArea(
        child: PageFrame(
          child: switch ((list, _error)) {
            (_, final Object e) when list == null => MessageView(message: errorText(l, e), onRetry: _load),
            (null, _) => const Center(child: CircularProgressIndicator()),
            (final CharacterList list, _) => RefreshIndicator(
              onRefresh: _load,
              child: ListView(
                padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, 96),
                children: [
                  if (!widget.auth.emailVerified) _VerifyBanner(auth: widget.auth),
                  Text(l.slotsUsed(list.characters.length, list.slotCap), style: Theme.of(context).textTheme.bodySmall),
                  const SizedBox(height: Space.md),
                  if (list.characters.isEmpty)
                    Padding(
                      padding: const EdgeInsets.only(top: Space.xxl),
                      child: MessageView(message: l.noCharacters, icon: Icons.shield_outlined),
                    ),
                  for (final c in list.characters) ...[
                    _CharacterCard(
                      c: c,
                      busy: _entering == c.id,
                      onTap: _entering == null ? () => _enter(c) : null,
                      onDelete: _entering == null ? () => _delete(c) : null,
                      onCancelDelete: () => _cancelDelete(c),
                    ),
                    const SizedBox(height: Space.md),
                  ],
                ],
              ),
            ),
          },
        ),
      ),
    );
  }
}

class _CharacterCard extends StatelessWidget {
  const _CharacterCard({required this.c, required this.busy, this.onTap, this.onDelete, this.onCancelDelete});

  final CharacterSummary c;
  final bool busy;
  final VoidCallback? onTap;
  final VoidCallback? onDelete;
  final VoidCallback? onCancelDelete;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final badge = switch (c.state) {
      'delete_grace' => (l.statePendingDelete, Palette.warning),
      'frozen' => (l.stateFrozen, Palette.info),
      _ => null,
    };
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: c.playable ? onTap : null,
        child: Padding(
          padding: const EdgeInsets.all(Space.lg),
          child: Row(
            children: [
              CircleAvatar(
                radius: 24,
                backgroundColor: Palette.bgElevated,
                child: Text(
                  c.name.characters.first.toUpperCase(),
                  style: const TextStyle(color: Palette.accent, fontWeight: FontWeight.w700, fontSize: 20),
                ),
              ),
              const SizedBox(width: Space.lg),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(c.name, style: text.titleMedium),
                    const SizedBox(height: 2),
                    Text(
                      l.levelLine(c.level, className(l, c.classCode), raceName(l, c.raceCode)),
                      style: text.bodySmall,
                    ),
                    if (badge != null) ...[
                      const SizedBox(height: Space.xs),
                      Text(
                        badge.$1,
                        style: TextStyle(color: badge.$2, fontSize: 12, fontWeight: FontWeight.w600),
                      ),
                    ],
                    if (c.state == 'delete_grace' && c.deleteAt != null)
                      Text(l.characterDeletesAt(_when(context, c.deleteAt!)), style: text.bodySmall),
                  ],
                ),
              ),
              if (busy)
                const SizedBox.square(dimension: 20, child: CircularProgressIndicator(strokeWidth: 2))
              else if (c.state == 'delete_grace')
                TextButton(onPressed: onCancelDelete, child: Text(l.cancelCharacterDelete))
              else if (c.playable) ...[
                IconButton(
                  tooltip: l.deleteCharacter,
                  icon: const Icon(Icons.delete_outline, color: Palette.textMuted),
                  onPressed: onDelete,
                ),
                const Icon(Icons.chevron_right, color: Palette.textMuted),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

String _when(BuildContext context, DateTime t) {
  final m = MaterialLocalizations.of(context);
  return '${m.formatMediumDate(t)} ${m.formatTimeOfDay(TimeOfDay.fromDateTime(t))}';
}

/// Deletion needs the character's name typed out; the delete button stays
/// off until it matches (the server checks it too).
class _DeleteCharacterDialog extends StatefulWidget {
  const _DeleteCharacterDialog({required this.c, required this.repo});

  final CharacterSummary c;
  final CharacterRepository repo;

  @override
  State<_DeleteCharacterDialog> createState() => _DeleteCharacterDialogState();
}

class _DeleteCharacterDialogState extends State<_DeleteCharacterDialog> {
  final _name = TextEditingController();
  bool _busy = false;
  String? _error;

  bool get _matches => _name.text.trim().toLowerCase() == widget.c.name.toLowerCase();

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await widget.repo.delete(widget.c.id, _name.text.trim());
      if (mounted) Navigator.pop(context, true);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.deleteCharacter),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(l.deleteCharacterBody(widget.c.name)),
          const SizedBox(height: Space.md),
          TextField(
            key: const Key('delete-character-name'),
            controller: _name,
            autocorrect: false,
            enableSuggestions: false,
            decoration: InputDecoration(labelText: l.characterNameConfirm, hintText: widget.c.name),
            onChanged: (_) => setState(() {}),
            onSubmitted: (_) {
              if (_matches && !_busy) _submit();
            },
          ),
          if (_error != null) ...[
            const SizedBox(height: Space.sm),
            Text(_error!, style: const TextStyle(color: Palette.danger)),
          ],
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context, false), child: Text(l.cancel)),
        FilledButton(
          key: const Key('delete-character-confirm'),
          style: FilledButton.styleFrom(backgroundColor: Palette.danger),
          onPressed: _matches && !_busy ? _submit : null,
          child: Text(l.deleteCharacter),
        ),
      ],
    );
  }
}

/// Shown until the account confirms its e-mail (DD-219).
class _VerifyBanner extends StatelessWidget {
  const _VerifyBanner({required this.auth});

  final AuthController auth;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      color: Palette.warning.withValues(alpha: 0.12),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(Space.lg, Space.md, Space.sm, Space.md),
        child: Row(
          children: [
            const Icon(Icons.mark_email_unread_outlined, color: Palette.warning),
            const SizedBox(width: Space.md),
            Expanded(child: Text(l.verifyEmailBanner)),
            TextButton(
              onPressed: () async {
                try {
                  await auth.resendVerification();
                  if (context.mounted) toast(context, Text(l.verificationSent));
                } on ApiError catch (e) {
                  if (context.mounted) toast(context, Text(errorText(l, e)));
                }
              },
              child: Text(l.resend),
            ),
          ],
        ),
      ),
    );
  }
}
