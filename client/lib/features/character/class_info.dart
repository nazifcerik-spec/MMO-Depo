import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import 'character_repository.dart';
import 'character_sheet_view.dart' show skillDescription;
import 'class_catalog.dart';

/// How the class resource fills (creation screen).
String resourceText(AppLocalizations l, ClassNode n) => [
  switch (n.resourceStartPct) {
    0 => l.resStartsEmpty,
    100 => l.resStartsFull,
    final p => l.resStartsAt(p),
  },
  if (n.resourcePerHit > 0) l.resPerHit(n.resourcePerHit),
  if (n.resourcePerTurn > 0) l.resPerTurn(n.resourcePerTurn),
].join(' · ');

/// One base attribute as a bar; [primary] ones are bold and full colour.
class StatLine extends StatelessWidget {
  const StatLine({
    super.key,
    required this.name,
    required this.value,
    required this.max,
    required this.color,
    required this.primary,
  });

  final String name;
  final int value;
  final int max;
  final Color color;
  final bool primary;

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(context).textTheme.bodySmall?.copyWith(fontWeight: primary ? FontWeight.w700 : null);
    return Padding(
      padding: const EdgeInsets.only(bottom: Space.xs),
      child: Row(
        children: [
          SizedBox(width: 96, child: Text(name, style: style)),
          Expanded(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(3),
              child: LinearProgressIndicator(
                value: (value / max).clamp(0, 1).toDouble(),
                minHeight: 6,
                color: primary ? color : color.withValues(alpha: 0.55),
                backgroundColor: Palette.bgSurface,
              ),
            ),
          ),
          SizedBox(
            width: 28,
            child: Text('$value', style: style, textAlign: TextAlign.right),
          ),
        ],
      ),
    );
  }
}

/// One skill: level, name, kind tag, effect text (and cost for actives).
class SkillLine extends StatelessWidget {
  const SkillLine({super.key, required this.skill});

  final ClassSkill skill;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final (tag, tagColor) = switch (skill.category) {
      'awakening' => (l.skillTagAwakening, Palette.warning),
      'ultimate' => (l.skillTagUltimate, Palette.danger),
      _ when !skill.active => (l.skillTagPassive, Palette.epic),
      _ => (l.skillTagActive, Palette.success),
    };
    return Container(
      padding: const EdgeInsets.symmetric(vertical: Space.sm),
      decoration: const BoxDecoration(
        border: Border(top: BorderSide(color: Palette.borderSubtle)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 52,
            child: Text(l.levelShort(skill.unlockLevel), style: text.bodySmall?.copyWith(color: Palette.textMuted)),
          ),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Wrap(
                  crossAxisAlignment: WrapCrossAlignment.center,
                  spacing: Space.xs,
                  children: [
                    Text(
                      ContentNames.of('skill', skill.code),
                      style: text.bodyMedium?.copyWith(fontWeight: FontWeight.w600),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 1),
                      decoration: BoxDecoration(
                        color: tagColor.withValues(alpha: 0.14),
                        borderRadius: BorderRadius.circular(Radii.sm),
                      ),
                      child: Text(tag, style: text.labelSmall?.copyWith(color: tagColor)),
                    ),
                  ],
                ),
                Text(skillDescription(l, skill), style: text.bodySmall?.copyWith(color: Palette.textSecondary)),
                if (skill.active)
                  Text(
                    l.skillCostCd(skill.mpCost, skill.cooldownTurns),
                    style: text.labelSmall?.copyWith(color: Palette.textMuted),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

/// Stand-in skill icon by effect until skill artwork ships.
IconData skillIcon(ClassSkill s) {
  if (s.category == 'ultimate') return Icons.whatshot;
  if (!s.active) return Icons.auto_awesome_outlined;
  return switch (s.effect) {
    'heal' => Icons.healing,
    'guard' => Icons.shield_outlined,
    'empower' => Icons.trending_up,
    'barrier' => Icons.security,
    'dot' => Icons.water_drop_outlined,
    'hot' => Icons.spa_outlined,
    'weaken' => Icons.trending_down,
    'stun' => Icons.flare,
    _ => Icons.flash_on,
  };
}

/// The class path as a pyramid: the base class on top, its two promotions
/// (level 100) below, their four specializations (level 300) at the base.
class ClassPathPyramid extends StatelessWidget {
  const ClassPathPyramid({super.key, required this.catalog, required this.root, required this.onTap});

  final ClassCatalog catalog;
  final String root;
  final ValueChanged<ClassNode> onTap;

  @override
  Widget build(BuildContext context) {
    final base = catalog.node(root);
    if (base == null) return const SizedBox.shrink();
    final promos = catalog.children(root);
    final specs = [for (final p in promos) ...catalog.children(p.code)];
    final color = classColor(root);
    Widget row(List<ClassNode> nodes, {double inset = 0}) => Padding(
      padding: EdgeInsets.symmetric(horizontal: inset),
      child: IntrinsicHeight(
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            for (var i = 0; i < nodes.length; i++) ...[
              if (i > 0) const SizedBox(width: Space.xs),
              Expanded(
                child: _PathNode(node: nodes[i], color: color, onTap: () => onTap(nodes[i])),
              ),
            ],
          ],
        ),
      ),
    );
    return LayoutBuilder(
      builder: (context, box) => Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          row([base], inset: box.maxWidth / 4),
          _Links(children: promos.length, color: color),
          row(promos),
          _Links(children: specs.length, color: color, parents: promos.length),
          row(specs),
        ],
      ),
    );
  }
}

/// Connector lines between two pyramid rows: each parent column fans out to
/// its share of the child columns.
class _Links extends StatelessWidget {
  const _Links({required this.children, required this.color, this.parents = 1});

  final int parents;
  final int children;
  final Color color;

  @override
  Widget build(BuildContext context) =>
      SizedBox(height: 18, child: CustomPaint(painter: _LinksPainter(parents, children, color.withValues(alpha: 0.6))));
}

class _LinksPainter extends CustomPainter {
  _LinksPainter(this.parents, this.children, this.color);

  final int parents;
  final int children;
  final Color color;

  @override
  void paint(Canvas canvas, Size size) {
    if (parents == 0 || children == 0) return;
    final p = Paint()
      ..color = color
      ..strokeWidth = 1.2;
    final mid = size.height / 2;
    final per = children ~/ parents;
    for (var i = 0; i < parents; i++) {
      final px = size.width * (2 * i + 1) / (2 * parents);
      canvas.drawLine(Offset(px, 0), Offset(px, mid), p);
      final first = size.width * (2 * (i * per) + 1) / (2 * children);
      final last = size.width * (2 * (i * per + per - 1) + 1) / (2 * children);
      canvas.drawLine(Offset(first, mid), Offset(last, mid), p);
      for (var j = 0; j < per; j++) {
        final cx = size.width * (2 * (i * per + j) + 1) / (2 * children);
        canvas.drawLine(Offset(cx, mid), Offset(cx, size.height), p);
      }
    }
  }

  @override
  bool shouldRepaint(_LinksPainter old) => old.parents != parents || old.children != children || old.color != color;
}

class _PathNode extends StatelessWidget {
  const _PathNode({required this.node, required this.color, required this.onTap});

  final ClassNode node;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final level = node.tier == 0 ? 1 : node.requiredLevel;
    return Material(
      color: node.tier == 0 ? color.withValues(alpha: 0.16) : Palette.bgSurface,
      borderRadius: BorderRadius.circular(Radii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: Space.sm),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(Radii.md),
            border: Border.all(color: node.tier == 2 ? Palette.borderSubtle : color.withValues(alpha: 0.5)),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              CircleArt(asset: classNodeAsset(node), size: node.tier == 2 ? 30 : 36, fallback: const SizedBox()),
              const SizedBox(height: 2),
              Text(l.levelShort(level), style: text.labelSmall?.copyWith(color: Palette.textMuted)),
              Text(
                ContentNames.of('class', node.code),
                textAlign: TextAlign.center,
                style: (node.tier == 2 ? text.bodySmall : text.bodyMedium)?.copyWith(fontWeight: FontWeight.w600),
              ),
              Text(
                ContentNames.of('class', '${node.code}.role', ''),
                textAlign: TextAlign.center,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                style: text.labelSmall?.copyWith(color: Palette.textSecondary),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// A small window with the skills a class node unlocks (base class: levels
/// 1–75, promotion: 100–250, specialization: 300–700 + its mastery title).
Future<void> showClassNodeSkills(BuildContext context, ClassCatalog catalog, ClassNode node, List<ClassSkill> all) {
  final l = AppLocalizations.of(context);
  final text = Theme.of(context).textTheme;
  final list = [
    for (final s in all)
      if (s.spec == node.code || (node.tier == 0 && s.spec.isEmpty)) s,
  ]..sort((a, b) => a.unlockLevel.compareTo(b.unlockLevel));
  final root = switch (node.tier) {
    0 => node.code,
    1 => node.parent,
    _ => catalog.node(node.parent)?.parent ?? '',
  };
  return showDialog<void>(
    context: context,
    builder: (ctx) => Dialog(
      backgroundColor: Palette.bgElevated,
      insetPadding: const EdgeInsets.all(Space.lg),
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 480, maxHeight: 560),
        child: Padding(
          padding: const EdgeInsets.fromLTRB(Space.lg, Space.lg, Space.lg, Space.sm),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Row(
                children: [
                  ClassEmblem(root: root, size: 28),
                  const SizedBox(width: Space.sm),
                  Expanded(
                    child: Text(
                      ContentNames.of('class', node.code),
                      style: text.titleMedium?.copyWith(color: classColor(root)),
                    ),
                  ),
                  Text(
                    l.levelShort(node.tier == 0 ? 1 : node.requiredLevel),
                    style: text.bodySmall?.copyWith(color: Palette.textMuted),
                  ),
                ],
              ),
              Text(
                ContentNames.of('class', '${node.code}.role', ''),
                style: text.bodySmall?.copyWith(color: Palette.textSecondary),
              ),
              const SizedBox(height: Space.sm),
              Flexible(
                child: ListView(shrinkWrap: true, children: [for (final s in list) SkillLine(skill: s)]),
              ),
              if (node.masteryTitle.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Space.sm),
                  child: Text(
                    l.masteryAt(catalog.milestones['mastery'] ?? 1000, ContentNames.of('title', node.masteryTitle)),
                    style: text.bodySmall?.copyWith(color: rarityColor('legendary')),
                  ),
                ),
              Align(
                alignment: Alignment.centerRight,
                child: TextButton(onPressed: () => Navigator.of(ctx).pop(), child: Text(l.close)),
              ),
            ],
          ),
        ),
      ),
    ),
  );
}
