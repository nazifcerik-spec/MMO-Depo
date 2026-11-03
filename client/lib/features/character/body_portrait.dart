import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_images.dart';

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
/// chest piece's outfit, framed in its rarity, holding its own main-hand
/// weapon and off-hand item. Races without figures yet show their portrait.
class BodyPortrait extends StatelessWidget {
  const BodyPortrait({
    super.key,
    required this.race,
    required this.body,
    this.equipped = const {},
    this.width = 118,
    this.onChangeBody,
  });

  final String race;
  final String body;
  final Map<String, EquippedItem> equipped;
  final double width;

  /// Switches the figure (male/female); null hides the button.
  final VoidCallback? onChangeBody;

  // Where the figure's hands are in the frame (scripts/body_art.py W x H).
  static const _hand = Offset(0.24, 0.55);

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final height = width * 1.6;
    final chest = equipped['chest'];
    final outfit = bodyOutfit(chest);
    final figure = ItemImages.bodyPath('${race}_${body}_$outfit') ?? ItemImages.bodyPath('${race}_${body}_none');
    final frame = chest == null ? Palette.borderStrong : rarityColor(chest.rarity);
    final dpr = MediaQuery.devicePixelRatioOf(context);

    Widget held(String? code, {required bool offhand}) {
      final p = code == null ? null : ItemImages.heldPath(code);
      if (p == null) return const SizedBox.shrink();
      final size = width * (offhand ? 0.36 : 0.5);
      final hand = Offset((offhand ? 1 - _hand.dx : _hand.dx) * width, _hand.dy * height);
      // Icons lie diagonally, grip low-left: turned 45° the weapon stands
      // upright with its grip at the bottom, held at the figure's right
      // hand (the viewer's left). Off-hand items (shields) sit on the other.
      return Positioned(
        left: hand.dx - size / 2,
        top: offhand ? hand.dy - size / 2 : hand.dy - size * 0.9,
        width: size,
        height: size,
        child: Transform.rotate(
          angle: offhand ? 0 : -math.pi / 4,
          child: Image.asset(p, cacheWidth: (size * dpr).round(), filterQuality: FilterQuality.medium),
        ),
      );
    }

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
                child: Image.asset(race, fit: BoxFit.cover, cacheWidth: (width * 2 * dpr).round()),
              )
            else
              const Center(child: Icon(Icons.person, color: Palette.textMuted, size: 48)),
            if (figure != null) ...[
              held(equipped['mainhand']?.templateCode, offhand: false),
              held(equipped['offhand']?.templateCode, offhand: true),
            ],
            if (onChangeBody != null)
              Positioned(
                right: 2,
                top: 2,
                child: Material(
                  color: Colors.black45,
                  shape: const CircleBorder(),
                  child: IconButton(
                    key: const Key('body-change'),
                    tooltip: l.bodyChange,
                    visualDensity: VisualDensity.compact,
                    iconSize: 16,
                    constraints: const BoxConstraints.tightFor(width: 28, height: 28),
                    padding: EdgeInsets.zero,
                    onPressed: onChangeBody,
                    icon: Icon(body == 'female' ? Icons.female : Icons.male, color: Palette.textSecondary),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
