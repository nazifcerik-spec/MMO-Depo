import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../../shared/content_names.dart';
import 'character_repository.dart';
import 'class_catalog.dart';
import 'class_info.dart';
import 'models.dart';

/// Bottom sheet: name + race + class, with the chosen class explained
/// (DD-253). Pops with the created character.
class CreateCharacterSheet extends StatefulWidget {
  const CreateCharacterSheet({super.key, required this.repo});

  final CharacterRepository repo;

  @override
  State<CreateCharacterSheet> createState() => _CreateCharacterSheetState();
}

class _CreateCharacterSheetState extends State<CreateCharacterSheet> {
  final _form = GlobalKey<FormState>();
  final _name = TextEditingController();
  String _race = races.first;
  String _class = classes.first;
  bool _busy = false;
  String? _error;

  ClassCatalog? _catalog;
  bool _catalogFailed = false;
  final _skills = <String, List<ClassSkill>>{};
  bool _started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_started) return;
    _started = true;
    _loadCatalog(Localizations.localeOf(context).languageCode);
  }

  Future<void> _loadCatalog(String lang) async {
    final api = widget.repo.api;
    try {
      final loads = await Future.wait([ClassCatalog.load(api), if (!ContentNames.loaded) ContentNames.load(api, lang)]);
      if (!mounted) return;
      setState(() => _catalog = loads.first as ClassCatalog);
      await _loadSkills(_class);
    } on ApiError {
      if (mounted) setState(() => _catalogFailed = true);
    }
  }

  Future<void> _loadSkills(String cls) async {
    if (_skills.containsKey(cls)) return;
    try {
      final list = await widget.repo.classSkills(cls);
      if (mounted) setState(() => _skills[cls] = list);
    } on ApiError {
      if (mounted) setState(() => _skills[cls] = const []);
    }
  }

  void _pickClass(String c) {
    _class = c;
    _loadSkills(c);
  }

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    if (!_form.currentState!.validate()) return;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final c = await widget.repo.create(_name.text.trim(), _race, _class);
      if (mounted) Navigator.of(context).pop(c);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;

    Widget choices(
      String label,
      List<String> codes,
      String selected,
      String Function(AppLocalizations, String) name,
      ValueChanged<String> onPick, {
      IconData Function(String)? icon,
    }) => Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: text.bodySmall),
        const SizedBox(height: Space.sm),
        Wrap(
          spacing: Space.sm,
          runSpacing: Space.sm,
          children: [
            for (final c in codes)
              ChoiceChip(
                avatar: icon == null ? null : Icon(icon(c), size: 16),
                label: Text(name(l, c)),
                selected: c == selected,
                onSelected: _busy ? null : (_) => setState(() => onPick(c)),
                selectedColor: Palette.accent.withValues(alpha: 0.2),
                side: BorderSide(color: c == selected ? Palette.accent : Palette.borderSubtle),
              ),
          ],
        ),
      ],
    );

    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: PageFrame(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Flexible(
              child: SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.sm),
                child: Form(
                  key: _form,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(l.newCharacter, style: text.titleLarge),
                      const SizedBox(height: Space.lg),
                      TextFormField(
                        controller: _name,
                        autofocus: true,
                        maxLength: 16,
                        decoration: InputDecoration(labelText: l.characterName, helperText: l.nameRule),
                        validator: (v) => characterNamePattern.hasMatch((v ?? '').trim()) ? null : l.nameRule,
                      ),
                      const SizedBox(height: Space.md),
                      choices(l.race, races, _race, raceName, (c) => _race = c),
                      const SizedBox(height: Space.lg),
                      choices(l.characterClass, classes, _class, className, _pickClass, icon: classIcon),
                      const SizedBox(height: Space.md),
                      if (_catalog case final cat?)
                        ClassInfoCard(catalog: cat, classCode: _class, skills: _skills[_class])
                      else if (_catalogFailed)
                        Text(l.classInfoFailed, style: text.bodySmall?.copyWith(color: Palette.textMuted))
                      else
                        const Padding(
                          padding: EdgeInsets.all(Space.md),
                          child: Center(
                            child: SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5)),
                          ),
                        ),
                    ],
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.screen),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (_error != null) ...[
                    Text(_error!, style: const TextStyle(color: Palette.danger)),
                    const SizedBox(height: Space.sm),
                  ],
                  FilledButton(
                    onPressed: _busy ? null : _submit,
                    child: _busy
                        ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
                        : Text(l.create),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
