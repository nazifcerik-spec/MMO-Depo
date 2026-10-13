import 'dart:math';

import 'package:flutter/material.dart';

import '../app/theme.dart';
import '../features/activity/activity_models.dart';
import '../l10n/app_localizations.dart';
import 'widgets.dart';

/// Shows a reveal result (activity / hunt / dungeon) as a bottom sheet.
Future<void> showLootSheet(BuildContext context, {required String title, required Loot loot, String? subtitle}) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    backgroundColor: Palette.bgSurface,
    builder: (_) => LootSheet(title: title, subtitle: subtitle, loot: loot),
  );
}

/// Loot reveal (docs/ui/design-system.md §10): rows scale 0.8→1.0 with a
/// rarity aura, staggered; one controller drives the whole list.
class LootSheet extends StatefulWidget {
  const LootSheet({super.key, required this.title, required this.loot, this.subtitle});

  final String title;
  final String? subtitle;
  final Loot loot;

  @override
  State<LootSheet> createState() => _LootSheetState();
}

class _LootSheetState extends State<LootSheet> with SingleTickerProviderStateMixin {
  late final AnimationController _anim;

  @override
  void initState() {
    super.initState();
    final rows = widget.loot.stacks.length + widget.loot.counts.length;
    _anim = AnimationController(
      vsync: this,
      duration: Duration(milliseconds: 400 + 60 * min(rows, 10)),
    )..forward();
  }

  @override
  void dispose() {
    _anim.dispose();
    super.dispose();
  }

  Widget _stagger(int i, int n, Widget child) {
    final start = n <= 1 ? 0.0 : (i / n) * 0.5;
    final end = min(1.0, start + 0.5);
    // Scale may overshoot (easeOutBack); opacity must stay within 0..1.
    final scale = CurvedAnimation(
      parent: _anim,
      curve: Interval(start, end, curve: Curves.easeOutBack),
    );
    final fade = CurvedAnimation(
      parent: _anim,
      curve: Interval(start, end, curve: Curves.easeOut),
    );
    return FadeTransition(
      opacity: fade,
      child: ScaleTransition(scale: Tween(begin: 0.8, end: 1.0).animate(scale), child: child),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final loot = widget.loot;
    final counts = loot.counts.entries.where((e) => e.value > 0).toList();
    final stacks = loot.stacks;
    final n = stacks.length + counts.length;
    var i = 0;

    return SafeArea(
      child: PageFrame(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxHeight: MediaQuery.sizeOf(context).height * 0.8),
          child: ListView(
            shrinkWrap: true,
            padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
            children: [
              Text(widget.title, style: text.titleLarge, textAlign: TextAlign.center),
              if (widget.subtitle != null) ...[
                const SizedBox(height: Space.xs),
                Text(widget.subtitle!, style: text.bodySmall, textAlign: TextAlign.center),
              ],
              const SizedBox(height: Space.lg),
              if (loot.goldCopper > 0 || loot.xp > 0)
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    if (loot.goldCopper > 0) GoldText(loot.goldCopper, size: 20),
                    if (loot.goldCopper > 0 && loot.xp > 0) const SizedBox(width: Space.xl),
                    if (loot.xp > 0)
                      Text(
                        l.xpGain(loot.xp),
                        style: const TextStyle(color: Palette.epic, fontWeight: FontWeight.w700, fontSize: 18),
                      ),
                  ],
                ),
              const SizedBox(height: Space.lg),
              if (loot.isEmpty) Text(l.lootNothing, textAlign: TextAlign.center, style: text.bodySmall),
              for (final s in stacks) ...[
                _stagger(
                  i++,
                  n,
                  _LootRow(
                    label: stackLabel(itemName(s.item.templateCode, s.item.prefix, s.item.suffix), s.count),
                    rarity: s.item.rarity,
                  ),
                ),
                const SizedBox(height: Space.sm),
              ],
              for (final c in counts) ...[
                _stagger(i++, n, _LootRow(label: '${rarityName(l, c.key)} × ${c.value}', rarity: c.key)),
                const SizedBox(height: Space.sm),
              ],
              const SizedBox(height: Space.md),
              FilledButton(onPressed: () => Navigator.of(context).pop(), child: Text(l.lootCollect)),
            ],
          ),
        ),
      ),
    );
  }
}

class _LootRow extends StatelessWidget {
  const _LootRow({required this.label, required this.rarity});

  final String label;
  final String rarity;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final c = rarityColor(rarity);
    final glow = rarity == 'legendary' || rarity == 'mythic' || rarity == 'epic';
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: Space.lg, vertical: Space.md),
      decoration: BoxDecoration(
        color: Palette.bgElevated,
        borderRadius: BorderRadius.circular(Radii.md),
        border: Border.all(color: c.withValues(alpha: 0.8), width: 1.5),
        boxShadow: glow ? [BoxShadow(color: c.withValues(alpha: 0.35), blurRadius: 16)] : null,
      ),
      child: Row(
        children: [
          Icon(Icons.diamond_outlined, color: c, size: 20),
          const SizedBox(width: Space.md),
          Expanded(
            child: Text(label, style: const TextStyle(fontWeight: FontWeight.w600)),
          ),
          Text(
            rarityName(l, rarity),
            style: TextStyle(color: c, fontSize: 12, fontWeight: FontWeight.w700),
          ),
        ],
      ),
    );
  }
}
