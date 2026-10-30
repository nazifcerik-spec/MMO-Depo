import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import 'inventory_models.dart';
import 'item_images.dart';

IconData itemIcon(String itemClass, [String subclass = '']) => switch ((itemClass, subclass)) {
  ('weapon', _) => Icons.hardware,
  ('armor', _) => Icons.shield_outlined,
  ('jewelry', _) => Icons.diamond_outlined, // DD-276
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
    this.wear,
  });

  final String? rarity;
  final IconData? icon;

  /// Template code: shows `assets/items/<code>.png` / `.webp` when it ships.
  final String? code;
  final int stack;
  final bool locked;

  /// Multi-select highlight (bag selection mode).
  final bool selected;

  /// Caption for empty equipment slots.
  final String? label;

  /// Equipment the character cannot wear (DD-290): `class` dims the cell
  /// with a red mark, `level` / `attribute` adds an amber one.
  final String? wear;
  final String? semantic;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final filled = icon != null;
    final c = filled ? rarityColor(rarity ?? 'common') : Palette.borderSubtle;
    final cell = _cell(filled, c);
    return wear == 'class' ? Opacity(opacity: 0.4, child: cell) : cell;
  }

  /// The cell-filling painted icon of [code], when one ships.
  String? get _painted {
    final p = code == null ? null : ItemImages.path(code!);
    return p != null && p.endsWith('.webp') ? p : null;
  }

  Widget _cell(bool filled, Color c) {
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
              if ((filled ? _painted : null) case final art?)
                // Painted gear icons carry their own backdrop: they fill the cell (DD-329).
                Positioned.fill(
                  child: Padding(
                    padding: const EdgeInsets.all(2),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(Radii.sm - 1),
                      child: LayoutBuilder(
                        builder: (context, box) => Image.asset(
                          art,
                          fit: BoxFit.cover,
                          cacheWidth: (box.maxWidth * MediaQuery.devicePixelRatioOf(context)).round(),
                          filterQuality: FilterQuality.medium,
                          errorBuilder: (_, _, _) => Icon(icon, color: c),
                        ),
                      ),
                    ),
                  ),
                )
              else
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
              if (wear case 'class' || 'level' || 'attribute')
                Positioned(
                  left: 2,
                  bottom: 2,
                  child: Icon(
                    wear == 'class' ? Icons.block : Icons.hourglass_bottom,
                    size: 12,
                    color: wear == 'class' ? Palette.danger : Palette.warning,
                  ),
                ),
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
  'ring' || 'ring1' || 'ring2' => l.slotRing, // a ring fills either ring slot (DD-276)
  'trinket' || 'trinket1' || 'trinket2' => l.slotTrinket,
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
  // Affix stats (DD-275); the label carries the % sign.
  'evasion_pct' => l.statEvasion,
  'block_pct' => l.statBlock,
  'crit_dmg_pct' => l.statCritDmg,
  'lifesteal_pct' => l.statLifesteal,
  'dmg_reduction_pct' => l.statDmgReduction,
  // Profession enchants (DD-293).
  'gather_yield_mining' => l.statGatherMining,
  'gather_yield_herbalism' => l.statGatherHerbalism,
  'gather_yield_woodcutting' => l.statGatherWoodcutting,
  'gather_yield_fishing' => l.statGatherFishing,
  'gather_yield_hunting' => l.statGatherHunting,
  'craft_speed' => l.statCraftSpeed,
  _ => target,
};

/// The gear groups a scroll fits, readable: "weapons", "armor · jewelry".
String enchantGroups(AppLocalizations l, EnchantInfo e) => [
  for (final g in e.slots)
    switch (g) {
      'weapon' => l.enchantGroupWeapon,
      'armor' => l.enchantGroupArmor,
      'jewelry' => l.enchantGroupJewelry,
      _ => l.enchantGroupAny,
    },
].join(' · ');

/// What an enchant adds, readable: "+1% Attack, +0.5 Crit" (DD-295).
String enchantEffects(AppLocalizations l, EnchantInfo e) => [for (final b in e.modifiers) ?bonusText(l, b)].join(', ');

/// "Enchant: Weapon Attack I · +1% Attack": the enchant a piece holds and,
/// when the server sent it, what it adds.
String enchantLabel(AppLocalizations l, String code, EnchantInfo? effect) {
  final on = l.enchantOn(ContentNames.of('item', code));
  final adds = effect == null ? '' : enchantEffects(l, effect);
  return adds.isEmpty ? on : '$on · $adds';
}

/// "+10 Attack", "+10% STR". Clamp/set modifiers are engine details → hidden.
String? bonusText(AppLocalizations l, StatBonus b) {
  final name = statName(l, b.target);
  String num(double v) =>
      v == v.roundToDouble() ? v.toInt().toString() : v.toStringAsFixed(2).replaceFirst(RegExp(r'0$'), '');
  return switch (b.kind) {
    'flat' => '${b.amount >= 0 ? '+' : ''}${num(b.amount)} $name',
    'mult_add' || 'mult_prod' => '${b.amount >= 0 ? '+' : ''}${num(b.amount * 100)}% $name',
    _ => null,
  };
}

String displayName(BagItem i) => itemName(i.templateCode, i.prefix, i.suffix);

/// An equippable item's slot and kind: "Chest · Leather armor", "Main hand · Sword".
String gearKind(AppLocalizations l, String slot, String subclass) =>
    subclass.isEmpty ? slotName(l, slot) : '${slotName(l, slot)} · ${ContentNames.of('equipment', subclass)}';

/// Why the character can or cannot wear an item (DD-290), with its colour.
(String, Color, IconData)? wearLine(AppLocalizations l, BagItem i) => switch (i.wear) {
  'ok' => (l.wearOk, Palette.success, Icons.check_circle_outline),
  'class' => (l.wearClass, Palette.danger, Icons.block),
  'level' => (l.wearLevel(i.minLevel), Palette.warning, Icons.hourglass_bottom),
  'attribute' => (l.wearAttribute, Palette.warning, Icons.hourglass_bottom),
  _ => null,
};

/// A raid set on a piece's details (DD-277): name, pieces worn, and each
/// bonus with its piece count — the earned ones lit.
class SetSection extends StatelessWidget {
  const SetSection({super.key, required this.set, required this.worn});

  final ItemSet set;
  final Iterable<String> worn;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final count = set.wornCount(worn);
    return Column(
      key: const Key('item-set'),
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '${ContentNames.of('itemset', set.code)}  $count/${set.pieces.length}',
          style: const TextStyle(color: Palette.epic, fontWeight: FontWeight.w700),
        ),
        for (final (pieces, mods) in set.bonuses)
          Text(
            '($pieces) ${[for (final m in mods) ?bonusText(l, m)].join(', ')}',
            style: TextStyle(color: count >= pieces ? Palette.success : Palette.textMuted, fontSize: 13),
          ),
      ],
    );
  }
}
