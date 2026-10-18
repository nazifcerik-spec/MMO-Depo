import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import 'inventory_models.dart';
import 'item_images.dart';

IconData itemIcon(String itemClass, [String subclass = '']) => switch ((itemClass, subclass)) {
  ('weapon', _) => Icons.hardware,
  ('armor', _) => Icons.shield_outlined,
  ('consumable', 'food') => Icons.restaurant_outlined, // cooking (DD-268)
  ('consumable', _) => Icons.local_drink_outlined,
  (_, 'ore') => Icons.landscape_outlined,
  (_, 'bar') => Icons.view_stream_outlined,
  (_, 'fuel') => Icons.local_fire_department_outlined,
  (_, 'fish') => Icons.set_meal_outlined,
  (_, 'wood' || 'log') => Icons.forest_outlined,
  (_, 'herb') => Icons.eco_outlined,
  (_, 'gem') => Icons.diamond_outlined,
  (_, 'hide' || 'leather' || 'pelt') => Icons.pets_outlined,
  (_, 'feather') => Icons.air,
  (_, 'scale') => Icons.texture,
  (_, 'stone') => Icons.terrain_outlined,
  _ => Icons.category_outlined,
};

/// Square item cell: rarity border, item art (or class icon), stack badge, lock marker.
class ItemCell extends StatelessWidget {
  const ItemCell({
    super.key,
    this.rarity,
    this.icon,
    this.code,
    this.stack = 1,
    this.locked = false,
    this.selected = false,
    this.label,
    this.semantic,
    this.onTap,
  });

  final String? rarity;
  final IconData? icon;

  /// Template code: shows `assets/items/<code>.png` when it ships.
  final String? code;
  final int stack;
  final bool locked;

  /// Multi-select highlight (bag selection mode).
  final bool selected;

  /// Caption for empty equipment slots.
  final String? label;
  final String? semantic;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final filled = icon != null;
    final c = filled ? rarityColor(rarity ?? 'common') : Palette.borderSubtle;
    return Semantics(
      button: onTap != null,
      selected: selected,
      label: semantic ?? label,
      excludeSemantics: true,
      child: Material(
        color: selected
            ? Palette.accent.withValues(alpha: 0.22)
            : filled
            ? Palette.bgElevated
            : Palette.bgSurface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(Radii.sm),
          side: selected
              ? const BorderSide(color: Palette.accent, width: 2)
              : BorderSide(color: filled ? c.withValues(alpha: 0.85) : Palette.borderSubtle, width: filled ? 1.5 : 1),
        ),
        child: InkWell(
          borderRadius: BorderRadius.circular(Radii.sm),
          onTap: onTap,
          child: Stack(
            children: [
              Center(
                child: filled
                    ? LayoutBuilder(
                        // Art scales with the cell (bag rows of 5, equipment rows of 6, detail headers).
                        builder: (_, box) {
                          final size = (box.maxWidth * 0.42).clamp(18.0, 48.0);
                          return code == null
                              ? Icon(icon, color: c, size: size)
                              : ItemArt(code: code!, icon: icon!, color: c, size: size);
                        },
                      )
                    : (label == null
                          ? null
                          : Padding(
                              padding: const EdgeInsets.all(2),
                              child: Text(
                                label!,
                                textAlign: TextAlign.center,
                                style: const TextStyle(fontSize: 9, color: Palette.textMuted),
                              ),
                            )),
              ),
              if (stack > 1)
                Positioned(
                  right: 3,
                  bottom: 1,
                  child: Text(
                    '$stack',
                    style: const TextStyle(fontSize: 10, fontWeight: FontWeight.w700, fontFeatures: numeric),
                  ),
                ),
              if (locked) const Positioned(left: 2, top: 2, child: Icon(Icons.lock, size: 10, color: Palette.warning)),
              if (selected)
                const Positioned(right: 2, top: 2, child: Icon(Icons.check_circle, size: 14, color: Palette.accent)),
            ],
          ),
        ),
      ),
    );
  }
}

String slotName(AppLocalizations l, String slot) => switch (slot) {
  'helm' => l.slotHelm,
  'necklace' => l.slotNecklace,
  'chest' => l.slotChest,
  'gloves' => l.slotGloves,
  'legs' => l.slotLegs,
  'boots' => l.slotBoots,
  'mainhand' => l.slotMainhand,
  'offhand' => l.slotOffhand,
  'ring1' || 'ring2' => l.slotRing,
  'trinket1' || 'trinket2' => l.slotTrinket,
  _ => slot,
};

String statName(AppLocalizations l, String target) => switch (target) {
  'str' => l.statStr,
  'agi' => l.statAgi,
  'int' => l.statInt,
  'vit' => l.statVit,
  'cha' => l.statCha,
  'luck' => l.statLuck,
  'attack' => l.attack,
  'defense' => l.defense,
  'crit_pct' => l.crit,
  'hit_pct' => l.hit,
  'hp_max' => l.hp,
  'mp_max' => l.mp,
  'en_max' => l.energy,
  _ => target,
};

/// "+10 Attack", "+10% STR". Clamp/set modifiers are engine details → hidden.
String? bonusText(AppLocalizations l, StatBonus b) {
  final name = statName(l, b.target);
  String num(double v) => v == v.roundToDouble() ? v.toInt().toString() : v.toStringAsFixed(1);
  return switch (b.kind) {
    'flat' => '${b.amount >= 0 ? '+' : ''}${num(b.amount)} $name',
    'mult_add' || 'mult_prod' => '${b.amount >= 0 ? '+' : ''}${num(b.amount * 100)}% $name',
    _ => null,
  };
}

String displayName(BagItem i) => itemName(i.templateCode, i.prefix, i.suffix);
