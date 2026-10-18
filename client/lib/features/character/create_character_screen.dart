import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../inventory/item_images.dart';
import 'character_repository.dart';
import 'character_sheet_view.dart' show skillDescription;
import 'class_catalog.dart';
import 'class_info.dart';
import 'models.dart';

/// Character creation (DD-266): pick a class — portrait, tagline, saying,
/// traits, skills (tap for a bubble), class path — and a race, then name the
/// character. Pops with the created character.
class CreateCharacterScreen extends StatefulWidget {
  const CreateCharacterScreen({super.key, required this.repo});

  final CharacterRepository repo;

  @override
  State<CreateCharacterScreen> createState() => _CreateCharacterScreenState();
}

class _CreateCharacterScreenState extends State<CreateCharacterScreen> {
  String _race = races.first;
  String _class = classes.first;
  String? _skill; // skill whose bubble is open

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
    setState(() {
      _class = c;
      _skill = null;
    });
    _loadSkills(c);
  }

  Future<void> _continue() async {
    final created = await showDialog<CharacterSummary>(
      context: context,
      builder: (_) => _NameDialog(repo: widget.repo, race: _race, classCode: _class),
    );
    if (created != null && mounted) Navigator.of(context).pop(created);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final cat = _catalog;
    final node = cat?.node(_class);
    return Scaffold(
      appBar: AppBar(title: Text(l.createCharacterTitle)),
      body: SafeArea(
        child: PageFrame(
          child: Column(
            children: [
              _ClassStrip(selected: _class, onPick: _pickClass),
              const Divider(height: 1),
              Expanded(
                child: ListView(
                  key: const ValueKey('create-character-list'),
                  padding: const EdgeInsets.only(bottom: Space.lg),
                  children: [
                    if (cat != null && node != null) ...[
                      _Hero(key: ValueKey('class-info-$_class'), classCode: _class),
                      _Section(
                        title: l.classTraits,
                        child: _Traits(catalog: cat, node: node),
                      ),
                      _Section(
                        title: l.firstSkills,
                        hint: l.skillsTapHint,
                        child: _Skills(
                          node: node,
                          skills: _skills[_class],
                          selected: _skill,
                          onTap: (code) => setState(() => _skill = _skill == code ? null : code),
                        ),
                      ),
                      _Section(
                        title: l.classPath,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.stretch,
                          children: [
                            ClassPathPyramid(
                              catalog: cat,
                              root: _class,
                              onTap: (n) => showClassNodeSkills(context, cat, n, _skills[_class] ?? const []),
                            ),
                            const SizedBox(height: Space.sm),
                            Text(
                              l.classPathHint,
                              textAlign: TextAlign.center,
                              style: text.bodySmall?.copyWith(color: Palette.textMuted),
                            ),
                          ],
                        ),
                      ),
                    ] else
                      Padding(
                        padding: const EdgeInsets.all(Space.xl),
                        child: _catalogFailed
                            ? Text(l.classInfoFailed, style: text.bodySmall?.copyWith(color: Palette.textMuted))
                            : const Center(
                                child: SizedBox.square(
                                  dimension: 22,
                                  child: CircularProgressIndicator(strokeWidth: 2.5),
                                ),
                              ),
                      ),
                    _Section(
                      title: l.chooseRace,
                      child: _Races(selected: _race, onPick: (r) => setState(() => _race = r)),
                    ),
                  ],
                ),
              ),
              const Divider(height: 1),
              Padding(
                padding: const EdgeInsets.fromLTRB(Space.screen, Space.md, Space.screen, Space.md),
                child: SizedBox(
                  width: double.infinity,
                  child: FilledButton(
                    onPressed: _continue,
                    child: Text(l.continueAs(ContentNames.of('class', _class, className(l, _class)))),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// The ten base classes as round icons; the chosen one in its colour.
class _ClassStrip extends StatelessWidget {
  const _ClassStrip({required this.selected, required this.onPick});

  final String selected;
  final ValueChanged<String> onPick;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return SizedBox(
      height: 84,
      child: ListView(
        key: const ValueKey('class-strip'),
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: Space.md, vertical: Space.sm),
        children: [
          for (final c in classes)
            Semantics(
              button: true,
              selected: c == selected,
              child: InkWell(
                borderRadius: BorderRadius.circular(Radii.md),
                onTap: () => onPick(c),
                child: SizedBox(
                  width: 68,
                  child: Column(
                    children: [
                      Container(
                        width: 44,
                        height: 44,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: c == selected ? classColor(c).withValues(alpha: 0.18) : Palette.bgElevated,
                          border: Border.all(
                            color: c == selected ? classColor(c) : Palette.borderStrong,
                            width: _border,
                          ),
                        ),
                        child: Icon(
                          classIcon(c),
                          size: 22,
                          color: c == selected ? classColor(c) : Palette.textSecondary,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        className(l, c),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(fontSize: 11, color: c == selected ? classColor(c) : Palette.textMuted),
                      ),
                    ],
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}

/// Portrait, name, tagline, role chips, introduction; the class saying under
/// both.
class _Hero extends StatelessWidget {
  const _Hero({super.key, required this.classCode});

  final String classCode;

  @override
  Widget build(BuildContext context) {
    final text = Theme.of(context).textTheme;
    final color = classColor(classCode);
    final roles = ContentNames.of('class', '$classCode.role', '').split(' / ').where((r) => r.isNotEmpty);
    final quote = ContentNames.of('class', '$classCode.quote', '');
    return Padding(
      padding: const EdgeInsets.fromLTRB(Space.screen, Space.lg, Space.screen, Space.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _Portrait(code: classCode, color: color),
              const SizedBox(width: Space.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(ContentNames.of('class', classCode), style: text.headlineSmall),
                    Text(
                      ContentNames.of('class', '$classCode.tagline', ''),
                      style: text.bodyMedium?.copyWith(color: color),
                    ),
                    const SizedBox(height: Space.sm),
                    Wrap(spacing: Space.xs, runSpacing: Space.xs, children: [for (final r in roles) _Chip(r)]),
                    const SizedBox(height: Space.sm),
                    Text(
                      ContentNames.of('class', '$classCode.intro', ''),
                      style: text.bodySmall?.copyWith(color: Palette.textSecondary, height: 1.45),
                    ),
                  ],
                ),
              ),
            ],
          ),
          if (quote.isNotEmpty) ...[
            const SizedBox(height: Space.md),
            Container(
              padding: const EdgeInsets.symmetric(vertical: Space.sm),
              decoration: const BoxDecoration(
                border: Border.symmetric(horizontal: BorderSide(color: Palette.borderSubtle)),
              ),
              child: Text(
                '"$quote"',
                textAlign: TextAlign.center,
                style: text.titleSmall?.copyWith(fontStyle: FontStyle.italic, fontFamily: 'serif'),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

/// The class portrait when one ships (`assets/classes/<code>.png`), else
/// its icon on a tinted card.
class _Portrait extends StatelessWidget {
  const _Portrait({required this.code, required this.color});

  final String code;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final fallback = Center(child: Icon(classIcon(code), size: 56, color: color));
    final path = ItemImages.classPath(code);
    return Container(
      width: 120,
      height: 160,
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.10),
        borderRadius: BorderRadius.circular(Radii.lg),
        border: Border.all(color: color.withValues(alpha: 0.45)),
      ),
      child: path == null
          ? fallback
          : Image.asset(
              path,
              fit: BoxFit.cover,
              cacheWidth: (120 * MediaQuery.devicePixelRatioOf(context)).round(),
              errorBuilder: (_, _, _) => fallback,
            ),
    );
  }
}

class _Section extends StatelessWidget {
  const _Section({required this.title, required this.child, this.hint});

  final String title;
  final String? hint;
  final Widget child;

  @override
  Widget build(BuildContext context) => Container(
    padding: const EdgeInsets.fromLTRB(Space.screen, Space.md, Space.screen, Space.md),
    decoration: const BoxDecoration(
      border: Border(top: BorderSide(color: Palette.borderSubtle)),
    ),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text.rich(
          TextSpan(
            text: title.toUpperCase(),
            style: const TextStyle(
              color: Palette.accent,
              fontSize: 12,
              fontWeight: FontWeight.w600,
              letterSpacing: 0.8,
            ),
            children: [
              if (hint != null)
                TextSpan(
                  text: ' · $hint',
                  style: const TextStyle(color: Palette.textMuted, fontWeight: FontWeight.w400, letterSpacing: 0),
                ),
            ],
          ),
        ),
        const SizedBox(height: Space.md),
        child,
      ],
    ),
  );
}

/// Resource, attribute focus with the base attributes, equipment, difficulty.
class _Traits extends StatelessWidget {
  const _Traits({required this.catalog, required this.node});

  final ClassCatalog catalog;
  final ClassNode node;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final color = classColor(node.code);
    final primary = ContentNames.of('attribute', node.primaryAttribute);
    final stats = node.baseStats.entries.toList()..sort((a, b) => b.value.compareTo(a.value));
    final (diffText, diffColor) = switch (classDifficulty(node.code)) {
      1 => (l.difficultyEasy, Palette.success),
      2 => (l.difficultyMedium, Palette.warning),
      _ => (l.difficultyHard, Palette.danger),
    };
    return Column(
      children: [
        _TraitRow(
          icon: Icons.local_fire_department_outlined,
          color: color,
          title: l.classResourceLine(ContentNames.of('resource', node.resourceKind)),
          detail: resourceText(l, node),
        ),
        _TraitRow(
          icon: Icons.fitness_center,
          color: color,
          title: node.secondaryAttributes.isEmpty
              ? l.classFocusOnly(primary)
              : l.classFocusLine(primary, ContentNames.of('attribute', node.secondaryAttributes.first)),
          child: Column(
            children: [
              for (final e in stats)
                StatLine(
                  name: ContentNames.of('attribute', e.key),
                  value: e.value,
                  max: catalog.maxBaseStat,
                  color: color,
                  primary: e.key == node.primaryAttribute || node.secondaryAttributes.contains(e.key),
                ),
            ],
          ),
        ),
        _TraitRow(
          icon: Icons.shield_outlined,
          color: color,
          title: l.classEquipment,
          detail: node.equipment.map((e) => ContentNames.of('equipment', e)).join(' · '),
        ),
        _TraitRow(
          icon: Icons.speed,
          color: color,
          title: l.classDifficulty,
          detail: diffText,
          detailColor: diffColor,
          last: true,
        ),
      ],
    );
  }
}

class _TraitRow extends StatelessWidget {
  const _TraitRow({
    required this.icon,
    required this.color,
    required this.title,
    this.detail = '',
    this.detailColor,
    this.child,
    this.last = false,
  });

  final IconData icon;
  final Color color;
  final String title;
  final String detail;
  final Color? detailColor;
  final Widget? child;
  final bool last;

  @override
  Widget build(BuildContext context) {
    final text = Theme.of(context).textTheme;
    return Padding(
      padding: EdgeInsets.only(bottom: last ? 0 : Space.md),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: Palette.bgElevated,
              border: Border.all(color: Palette.borderStrong),
            ),
            child: Icon(icon, size: 18, color: color),
          ),
          const SizedBox(width: Space.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: text.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                if (detail.isNotEmpty)
                  Text(detail, style: text.bodySmall?.copyWith(color: detailColor ?? Palette.textSecondary)),
                if (child case final c?) ...[const SizedBox(height: Space.xs), c],
              ],
            ),
          ),
        ],
      ),
    );
  }
}

/// The base class's skills as icons; the tapped one opens a bubble under
/// the row, its arrow pointing at the skill.
class _Skills extends StatelessWidget {
  const _Skills({required this.node, required this.skills, required this.selected, required this.onTap});

  final ClassNode node;
  final List<ClassSkill>? skills;
  final String? selected;
  final ValueChanged<String> onTap;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final all = skills;
    if (all == null) {
      return const Center(child: SizedBox.square(dimension: 20, child: CircularProgressIndicator(strokeWidth: 2)));
    }
    final list = [
      for (final s in all)
        if (s.spec.isEmpty || s.spec == node.code) s,
    ]..sort((a, b) => a.unlockLevel.compareTo(b.unlockLevel));
    if (list.isEmpty) return const SizedBox.shrink();
    final color = classColor(node.code);
    final index = list.indexWhere((s) => s.code == selected);
    return LayoutBuilder(
      builder: (context, box) {
        final colW = box.maxWidth / list.length;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                for (final s in list)
                  SizedBox(
                    width: colW,
                    child: InkWell(
                      borderRadius: BorderRadius.circular(Radii.md),
                      onTap: () => onTap(s.code),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(vertical: Space.xs, horizontal: 2),
                        child: Column(
                          children: [
                            Container(
                              width: 40,
                              height: 40,
                              decoration: BoxDecoration(
                                color: s.code == selected ? color.withValues(alpha: 0.18) : Palette.bgElevated,
                                borderRadius: BorderRadius.circular(Radii.md),
                                border: Border.all(
                                  color: s.code == selected ? color : Palette.borderStrong,
                                  width: _border,
                                ),
                              ),
                              child: Icon(skillIcon(s), size: 20, color: s.active ? color : Palette.epic),
                            ),
                            const SizedBox(height: 4),
                            _Lines(ContentNames.of('skill', s.code), lines: 2, style: const TextStyle(fontSize: 11)),
                            Text(
                              l.levelShort(s.unlockLevel),
                              style: const TextStyle(fontSize: 11, color: Palette.textMuted),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
              ],
            ),
            if (index >= 0)
              _Bubble(
                key: ValueKey('skill-bubble-${list[index].code}'),
                arrowX: colW * index + colW / 2,
                color: color,
                child: Text.rich(
                  TextSpan(
                    children: [
                      TextSpan(
                        text: ContentNames.of('skill', list[index].code),
                        style: const TextStyle(fontWeight: FontWeight.w600),
                      ),
                      TextSpan(
                        text: ' · ${_skillMeta(l, list[index], node)}\n',
                        style: const TextStyle(color: Palette.textMuted),
                      ),
                      TextSpan(
                        text: skillDescription(l, list[index]),
                        style: const TextStyle(color: Palette.textSecondary),
                      ),
                    ],
                  ),
                  style: text.bodySmall?.copyWith(height: 1.45),
                ),
              ),
          ],
        );
      },
    );
  }
}

String _skillMeta(AppLocalizations l, ClassSkill s, ClassNode node) => [
  switch (s.category) {
    'awakening' => l.skillTagAwakening,
    'ultimate' => l.skillTagUltimate,
    _ when !s.active => l.skillTagPassive,
    _ => l.skillTagActive,
  },
  l.levelShort(s.unlockLevel),
  if (s.active && s.mpCost > 0) '${s.mpCost} ${ContentNames.of('resource', node.resourceKind)}',
].join(' · ');

/// A rounded box with an arrow on top at [arrowX].
class _Bubble extends StatelessWidget {
  const _Bubble({super.key, required this.arrowX, required this.color, required this.child});

  final double arrowX;
  final Color color;
  final Widget child;

  static const _arrow = 8.0;

  @override
  Widget build(BuildContext context) => Padding(
    padding: const EdgeInsets.only(top: Space.xs),
    child: CustomPaint(
      painter: _BubblePainter(arrowX: arrowX, color: color, arrow: _arrow),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(Space.md, Space.sm + _arrow, Space.md, Space.sm),
        child: SizedBox(width: double.infinity, child: child),
      ),
    ),
  );
}

class _BubblePainter extends CustomPainter {
  _BubblePainter({required this.arrowX, required this.color, required this.arrow});

  final double arrowX;
  final Color color;
  final double arrow;

  @override
  void paint(Canvas canvas, Size size) {
    final x = arrowX.clamp(Radii.md + arrow, size.width - Radii.md - arrow);
    final body = RRect.fromRectAndRadius(
      Rect.fromLTWH(0, arrow, size.width, size.height - arrow),
      const Radius.circular(Radii.md),
    );
    final path = Path()
      ..addRRect(body)
      ..moveTo(x - arrow, arrow)
      ..lineTo(x, 0)
      ..lineTo(x + arrow, arrow)
      ..close();
    canvas.drawPath(path, Paint()..color = Palette.bgElevated);
    final outline = Path()
      ..moveTo(body.left + Radii.md, body.top)
      ..lineTo(x - arrow, arrow)
      ..lineTo(x, 0)
      ..lineTo(x + arrow, arrow)
      ..lineTo(body.right - Radii.md, body.top)
      ..arcToPoint(Offset(body.right, body.top + Radii.md), radius: const Radius.circular(Radii.md))
      ..lineTo(body.right, body.bottom - Radii.md)
      ..arcToPoint(Offset(body.right - Radii.md, body.bottom), radius: const Radius.circular(Radii.md))
      ..lineTo(body.left + Radii.md, body.bottom)
      ..arcToPoint(Offset(body.left, body.bottom - Radii.md), radius: const Radius.circular(Radii.md))
      ..lineTo(body.left, body.top + Radii.md)
      ..arcToPoint(Offset(body.left + Radii.md, body.top), radius: const Radius.circular(Radii.md));
    canvas.drawPath(
      outline,
      Paint()
        ..color = color
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.2,
    );
  }

  @override
  bool shouldRepaint(_BubblePainter old) => old.arrowX != arrowX || old.color != color;
}

/// The four races: portrait (or icon), name, a line of lore. Looks only.
class _Races extends StatelessWidget {
  const _Races({required this.selected, required this.onPick});

  final String selected;
  final ValueChanged<String> onPick;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    String desc(String r) => switch (r) {
      'human' => l.raceHumanDesc,
      'elf' => l.raceElfDesc,
      'dwarf' => l.raceDwarfDesc,
      'orc' => l.raceOrcDesc,
      'troll' => l.raceTrollDesc,
      'undead' => l.raceUndeadDesc,
      'dark_elf' => l.raceDarkElfDesc,
      'drakari' => l.raceDrakariDesc,
      'tauren' => l.raceTaurenDesc,
      'gnome' => l.raceGnomeDesc,
      _ => '',
    };
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // Four cards a row (ten races: 4 + 4 + 2).
        LayoutBuilder(
          builder: (context, box) {
            final w = (box.maxWidth - 3 * Space.xs) / 4;
            return Wrap(
              spacing: Space.xs,
              runSpacing: Space.xs,
              children: [
                for (final r in races)
                  SizedBox(
                    width: w,
                    child: _RaceCard(
                      code: r,
                      name: raceName(l, r),
                      desc: desc(r),
                      selected: r == selected,
                      onTap: () => onPick(r),
                    ),
                  ),
              ],
            );
          },
        ),
        const SizedBox(height: Space.sm),
        Text(
          l.raceCosmeticHint,
          textAlign: TextAlign.center,
          style: const TextStyle(fontSize: 11, color: Palette.textMuted),
        ),
      ],
    );
  }
}

class _RaceCard extends StatelessWidget {
  const _RaceCard({
    required this.code,
    required this.name,
    required this.desc,
    required this.selected,
    required this.onTap,
  });

  final String code;
  final String name;
  final String desc;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final path = ItemImages.racePath(code);
    const fallback = Icon(Icons.person_outline, color: Palette.textMuted);
    return Material(
      color: selected ? Palette.accent.withValues(alpha: 0.12) : Colors.transparent,
      borderRadius: BorderRadius.circular(Radii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(6),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(Radii.md),
            border: Border.all(color: selected ? Palette.accent : Palette.borderSubtle, width: _border),
          ),
          child: Column(
            children: [
              Container(
                height: 48,
                width: double.infinity,
                clipBehavior: Clip.antiAlias,
                decoration: BoxDecoration(color: Palette.bgElevated, borderRadius: BorderRadius.circular(Radii.sm)),
                child: path == null
                    ? fallback
                    : Image.asset(path, fit: BoxFit.cover, errorBuilder: (_, _, _) => fallback),
              ),
              const SizedBox(height: 4),
              _Lines(name, lines: 1, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600)),
              _Lines(desc, lines: 2, style: const TextStyle(fontSize: 11, color: Palette.textSecondary)),
            ],
          ),
        ),
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip(this.label);

  final String label;

  @override
  Widget build(BuildContext context) => Container(
    padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 3),
    decoration: BoxDecoration(
      color: Palette.bgElevated,
      borderRadius: BorderRadius.circular(999),
      border: Border.all(color: Palette.borderSubtle),
    ),
    child: Text(label, style: const TextStyle(fontSize: 11, color: Palette.textSecondary)),
  );
}

/// One border width for every selectable item; selection changes only the
/// colour, so nothing shifts when it is picked.
const _border = 1.5;

/// Text in a box exactly [lines] tall whatever its length, so items of the
/// same kind line up: one line scales a long word down, more lines ellipsize.
class _Lines extends StatelessWidget {
  const _Lines(this.text, {required this.lines, required this.style});

  final String text;
  final int lines;
  final TextStyle style;

  static const _height = 1.3;

  @override
  Widget build(BuildContext context) {
    final st = style.copyWith(height: _height);
    final scale = MediaQuery.textScalerOf(context).scale(st.fontSize ?? 14);
    return SizedBox(
      height: lines * scale * _height,
      width: double.infinity,
      child: lines == 1
          ? FittedBox(
              fit: BoxFit.scaleDown,
              child: Text(text, maxLines: 1, style: st),
            )
          : Text(text, textAlign: TextAlign.center, maxLines: lines, overflow: TextOverflow.ellipsis, style: st),
    );
  }
}

/// Asks for the name, then creates the character (pops with it).
class _NameDialog extends StatefulWidget {
  const _NameDialog({required this.repo, required this.race, required this.classCode});

  final CharacterRepository repo;
  final String race;
  final String classCode;

  @override
  State<_NameDialog> createState() => _NameDialogState();
}

class _NameDialogState extends State<_NameDialog> {
  final _form = GlobalKey<FormState>();
  final _name = TextEditingController();
  bool _busy = false;
  String? _error;

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
      final c = await widget.repo.create(_name.text.trim(), widget.race, widget.classCode);
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
    return AlertDialog(
      title: Text(l.nameYourCharacter),
      content: Form(
        key: _form,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              '${raceName(l, widget.race)} · ${ContentNames.of('class', widget.classCode, className(l, widget.classCode))}',
              style: TextStyle(color: classColor(widget.classCode)),
            ),
            const SizedBox(height: Space.sm),
            TextFormField(
              controller: _name,
              autofocus: true,
              maxLength: 16,
              decoration: InputDecoration(labelText: l.characterName, helperText: l.nameRule),
              validator: (v) => characterNamePattern.hasMatch((v ?? '').trim()) ? null : l.nameRule,
              onFieldSubmitted: (_) => _submit(),
            ),
            if (_error != null) Text(_error!, style: const TextStyle(color: Palette.danger)),
          ],
        ),
      ),
      actions: [
        TextButton(onPressed: _busy ? null : () => Navigator.of(context).pop(), child: Text(l.cancel)),
        FilledButton(
          onPressed: _busy ? null : _submit,
          child: _busy
              ? const SizedBox.square(dimension: 18, child: CircularProgressIndicator(strokeWidth: 2))
              : Text(l.create),
        ),
      ],
    );
  }
}
