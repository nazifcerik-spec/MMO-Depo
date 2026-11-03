import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_images.dart';
import 'appearance_picker.dart';
import 'models.dart';
import 'paper_doll.dart';

/// The head's place among the front layers: necklaces go under it, helms over.
const headZ = 10;

/// Level bands of the chest piece, one outfit each (scripts/body_art.py BANDS).
const bodyBands = [1, 100, 300, 500, 700, 850];

/// The outfit the chest piece puts on the figure (DD-344): its armor type
/// (cloth, leather, plate) and level band, `none` without one.
String bodyOutfit(EquippedItem? chest) {
  final type = switch (chest?.itemSubclass) {
    'cloth_armor' => 'cloth',
    'leather_armor' => 'leather',
    'plate_armor' => 'plate',
    _ => null,
  };
  if (chest == null || type == null) return 'none';
  var band = 1;
  for (var i = 0; i < bodyBands.length; i++) {
    if (chest.minLevel >= bodyBands[i]) band = i + 1;
  }
  return '$type$band';
}

/// The character drawn full-body (DD-344): race and body figure in the
/// chest piece's outfit with the chosen head on it (DD-346), framed in its
/// rarity, wearing its equipment's sprites as paper doll layers (DD-347). Races
/// without figures yet show their portrait.
class BodyPortrait extends StatelessWidget {
  const BodyPortrait({
    super.key,
    required this.race,
    required this.look,
    this.equipped = const {},
    this.width = 118,
    this.onEdit,
  });

  final String race;
  final Appearance look;
  final Map<String, EquippedItem> equipped;
  final double width;

  /// Opens the look editor; null hides the button.
  final VoidCallback? onEdit;

  // Where the figure's eyes are in the frame, and how wide a head frame must
  // be to fit its head.
  static const _eyes = Offset(0.505, 0.122);
  static const _head = 0.34;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final height = width * 1.6;
    final chest = equipped['chest'];
    final outfit = bodyOutfit(chest);
    final body = look.body;
    final figure = ItemImages.bodyPath('${race}_${body}_$outfit') ?? ItemImages.bodyPath('${race}_${body}_none');
    final frame = chest == null ? Palette.borderStrong : rarityColor(chest.rarity);
    final dpr = MediaQuery.devicePixelRatioOf(context);

    // Equipment layers from the paper doll metadata (DD-347), under or over
    // the figure, in z order. A missing sprite simply draws nothing.
    final layers = [
      if (figure != null)
        for (final MapEntry(key: slot, value: item) in equipped.entries)
          ?PaperDoll.current.layerOf(
            equipmentId: item.templateCode,
            slot: slot,
            subclass: item.itemSubclass,
            race: race,
            body: body,
          ),
    ]..sort((a, b) => a.zIndex.compareTo(b.zIndex));
    Widget sprite(DollLayer d, {bool mirror = false}) {
      final size = d.scale * width;
      final anchor = mirror ? d.mirrorAnchor! : d.anchor;
      final side = mirror ? -1.0 : 1.0;
      final pivotX = mirror ? 1 - d.pivot.dx : d.pivot.dx;
      return Positioned(
        key: ValueKey('doll-${d.slot}${mirror ? '-mirror' : ''}'),
        left: (anchor.dx + side * d.offsetX) * width - pivotX * size,
        top: (anchor.dy + d.offsetY) * height - d.pivot.dy * size,
        width: size,
        height: size,
        child: Transform.rotate(
          angle: side * d.rotation * math.pi / 180,
          alignment: Alignment(pivotX * 2 - 1, d.pivot.dy * 2 - 1),
          child: Transform.flip(
            flipX: d.flipX != mirror,
            child: Image.asset(
              d.avatarAsset,
              cacheWidth: (size * dpr).round(),
              filterQuality: FilterQuality.medium,
              gaplessPlayback: true,
              errorBuilder: (_, _, _) => const SizedBox.shrink(),
            ),
          ),
        ),
      );
    }

    List<Widget> sprites(DollLayer d) => [sprite(d), if (d.mirrorAnchor != null) sprite(d, mirror: true)];

    final raceArt = ItemImages.racePath(race);
    return Semantics(
      label: chest == null ? l.portraitNoArmor : null,
      child: Container(
        key: const Key('body-portrait'),
        width: width,
        height: height,
        decoration: BoxDecoration(
          color: Palette.bgPrimary,
          borderRadius: BorderRadius.circular(Radii.md),
          border: Border.all(color: frame.withValues(alpha: 0.9), width: 1.5),
          boxShadow: [
            if (chest != null && (chest.rarity == 'legendary' || chest.rarity == 'mythic'))
              BoxShadow(color: frame.withValues(alpha: 0.35), blurRadius: 14),
          ],
        ),
        clipBehavior: Clip.antiAlias,
        child: Stack(
          clipBehavior: Clip.none,
          children: [
            for (final d in layers)
              if (d.behind) ...sprites(d),
            if (figure != null)
              Positioned.fill(
                child: Image.asset(
                  figure,
                  key: ValueKey(figure),
                  fit: BoxFit.cover,
                  cacheWidth: (width * dpr).round(),
                  filterQuality: FilterQuality.medium,
                ),
              )
            else if (raceArt != null)
              Positioned.fill(
                child: Image.asset(raceArt, fit: BoxFit.cover, cacheWidth: (width * 2 * dpr).round()),
              )
            else
              const Center(child: Icon(Icons.person, color: Palette.textMuted, size: 48)),
            for (final d in layers)
              if (!d.behind && d.zIndex < headZ) ...sprites(d),
            if (figure != null && customizableRaces.contains(race))
              Positioned(
                left: _eyes.dx * width - width * _head / 2,
                top: _eyes.dy * height - width * _head / 2,
                child: HeadLayers(
                  look: look,
                  size: width * _head,
                  onFigure: true,
                  hideHair: layers.any((d) => d.hidesHair),
                ),
              ),
            for (final d in layers)
              if (!d.behind && d.zIndex >= headZ) ...sprites(d),
            if (onEdit != null)
              Positioned(
                right: 2,
                top: 2,
                child: Material(
                  color: Colors.black45,
                  shape: const CircleBorder(),
                  child: IconButton(
                    key: const Key('look-edit'),
                    tooltip: l.bodyChange,
                    visualDensity: VisualDensity.compact,
                    iconSize: 16,
                    constraints: const BoxConstraints.tightFor(width: 28, height: 28),
                    padding: EdgeInsets.zero,
                    onPressed: onEdit,
                    icon: const Icon(Icons.face_retouching_natural, color: Palette.textSecondary),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
