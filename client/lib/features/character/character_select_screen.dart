import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../auth/auth_controller.dart';
import 'character_repository.dart';
import '../home/home_shell.dart';
import 'create_character_sheet.dart';
import 'models.dart';

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
    final created = await showModalBottomSheet<CharacterSummary>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => CreateCharacterSheet(repo: widget.repo),
    );
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

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final list = _list;
    return Scaffold(
      appBar: AppBar(
        title: Text(l.characters),
        actions: [IconButton(tooltip: l.logout, icon: const Icon(Icons.logout), onPressed: widget.auth.logout)],
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
                  Text(l.slotsUsed(list.characters.length, list.slotCap), style: Theme.of(context).textTheme.bodySmall),
                  const SizedBox(height: Space.md),
                  if (list.characters.isEmpty)
                    Padding(
                      padding: const EdgeInsets.only(top: Space.xxl),
                      child: MessageView(message: l.noCharacters, icon: Icons.shield_outlined),
                    ),
                  for (final c in list.characters) ...[
                    _CharacterCard(c: c, busy: _entering == c.id, onTap: _entering == null ? () => _enter(c) : null),
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
  const _CharacterCard({required this.c, required this.busy, this.onTap});

  final CharacterSummary c;
  final bool busy;
  final VoidCallback? onTap;

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
                  ],
                ),
              ),
              if (busy)
                const SizedBox.square(dimension: 20, child: CircularProgressIndicator(strokeWidth: 2))
              else if (c.playable)
                const Icon(Icons.chevron_right, color: Palette.textMuted),
            ],
          ),
        ),
      ),
    );
  }
}
