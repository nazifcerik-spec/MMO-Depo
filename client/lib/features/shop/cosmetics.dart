import 'package:flutter/material.dart';

import '../../app/theme.dart';

/// Cosmetics (DD-321): 30 page backgrounds (DD-323) and 80 avatar portraits
/// (DD-324), bundled as art. Codes match the server's (bg01..bg30,
/// av01..av80); the first 10 backgrounds and 40 avatars are free.
const backgroundCount = 30, freeBackgrounds = 10, avatarCount = 80, freeAvatars = 40;

String backgroundCode(int n) => 'bg${n.toString().padLeft(2, '0')}';
String avatarCode(int n) => 'av${n.toString().padLeft(2, '0')}';

int _number(String code) => int.tryParse(code.substring(2)) ?? 1;

bool isFree(String code) => code.startsWith('bg') ? _number(code) <= freeBackgrounds : _number(code) <= freeAvatars;

int _background(String code) => _number(code).clamp(1, backgroundCount);

/// Background [code]'s art, bundled under the same code
/// (`scripts/bg_assets.py`).
String backgroundAsset(String code) => 'assets/backgrounds/${backgroundCode(_background(code))}.webp';

/// Each background's accent and secondary colour, picked from its art and
/// readable on the darkened page: buttons, bars and the selected tab.
const _backgroundColors = <(Color, Color)>[
  (Color(0xFFE3A857), Color(0xFF6B7A8F)), // bg01 Ironhold
  (Color(0xFF7BC67E), Color(0xFF3E6B48)), // bg02 Greenwood
  (Color(0xFF8FB8C9), Color(0xFF4A6670)), // bg03 Mist Forest
  (Color(0xFF9DBF6A), Color(0xFF3F5A4A)), // bg04 Misty Marshes
  (Color(0xFFA9B8C8), Color(0xFF55606E)), // bg05 Iron Mountains
  (Color(0xFFE8B864), Color(0xFF8A5A3C)), // bg06 Crown Plains
  (Color(0xFF6FA8DC), Color(0xFF2F4858)), // bg07 Storm Coast
  (Color(0xFFE6D3B3), Color(0xFF9C7A5B)), // bg08 Salt Desert
  (Color(0xFF6EC9E0), Color(0xFF3B4F8C)), // bg09 Crystal Caves
  (Color(0xFF7FD1B9), Color(0xFF4B5D7A)), // bg10 Northern Keep
  (Color(0xFFFF7A3D), Color(0xFF8C2F1F)), // bg11 Emberpeak
  (Color(0xFFE8743B), Color(0xFF4A3530)), // bg12 Burnt Forest
  (Color(0xFFB48EDB), Color(0xFF3D2E4F)), // bg13 Royal Crypt
  (Color(0xFFC9C3E6), Color(0xFF4E4A6B)), // bg14 Crown Graveyard
  (Color(0xFFF2C14E), Color(0xFF7A5230)), // bg15 Buried Palace
  (Color(0xFF4FD1C5), Color(0xFF1F4E66)), // bg16 Sunken Temple
  (Color(0xFF5ED6A0), Color(0xFF1E3D3A)), // bg17 Sunken Fleet
  (Color(0xFFF29E4C), Color(0xFF2E5E7E)), // bg18 Dragonspine Isles
  (Color(0xFFE0533D), Color(0xFF5A2A2A)), // bg19 Dragon Slopes
  (Color(0xFFFFB347), Color(0xFF6E1F1F)), // bg20 Dragon Throne
  (Color(0xFFFF5E3A), Color(0xFF2B1A1A)), // bg21 Volcano Heart
  (Color(0xFF8ED1FC), Color(0xFF5A7FB5)), // bg22 Floating Isles
  (Color(0xFFF5D77A), Color(0xFF3A4A8C)), // bg23 Celestial Dome
  (Color(0xFFB57BFF), Color(0xFF2E2366)), // bg24 Sky Rift
  (Color(0xFFA7B1C2), Color(0xFF1F2430)), // bg25 Starless Lands
  (Color(0xFF9FA8DA), Color(0xFF2C2F4A)), // bg26 Darkened Observatory
  (Color(0xFFD96C4F), Color(0xFF3A3330)), // bg27 Ash Sea
  (Color(0xFFC77DFF), Color(0xFF14101F)), // bg28 Shore of the Void
  (Color(0xFF7AE7FF), Color(0xFF24164A)), // bg29 Edge of Eternity
  (Color(0xFFFFE08A), Color(0xFF2A1F4D)), // bg30 Eternal Throne
];

/// Background [code]'s (accent, secondary) colours.
(Color, Color) backgroundColors(String code) => _backgroundColors[_background(code) - 1];

/// [base] in background [code]'s colours.
ThemeData themedFor(ThemeData base, String code) {
  final (accent, secondary) = backgroundColors(code);
  return base.copyWith(
    colorScheme: base.colorScheme.copyWith(
      primary: accent,
      secondary: secondary,
      secondaryContainer: secondary,
      onSecondaryContainer: Palette.textPrimary,
    ),
  );
}

/// Background [code] filling its box: the art covers it (edges a phone's
/// aspect does not show are cropped), darkened so text stays readable.
/// [thumbnail] decodes it small and undarkened for the shop's tiles.
class BackgroundArt extends StatelessWidget {
  const BackgroundArt({super.key, required this.code, this.child, this.thumbnail = false});

  final String code;
  final Widget? child;
  final bool thumbnail;

  @override
  Widget build(BuildContext context) => Stack(
    fit: StackFit.expand,
    children: [
      Image.asset(
        backgroundAsset(code),
        fit: BoxFit.cover,
        cacheWidth: thumbnail ? 160 : null,
        gaplessPlayback: true,
        errorBuilder: (_, _, _) => const ColoredBox(color: Palette.bgPrimary),
      ),
      if (!thumbnail)
        const DecoratedBox(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [Color(0x99000000), Color(0x73000000), Color(0xBF000000)],
            ),
          ),
        ),
      ?child,
    ],
  );
}

/// Avatar [code]'s portrait, bundled under the same code
/// (`scripts/avatar_assets.py`).
String avatarAsset(String code) => 'assets/avatars/${avatarCode(_number(code).clamp(1, avatarCount))}.webp';

/// Avatar [code]: its portrait in a circle, decoded at the size shown.
class AvatarBadge extends StatelessWidget {
  const AvatarBadge({super.key, required this.code, this.size = 40});

  final String code;
  final double size;

  @override
  Widget build(BuildContext context) => ClipOval(
    child: Image.asset(
      avatarAsset(code),
      width: size,
      height: size,
      fit: BoxFit.cover,
      cacheWidth: (size * MediaQuery.devicePixelRatioOf(context)).round(),
      errorBuilder: (_, _, _) => CircleAvatar(
        radius: size / 2,
        backgroundColor: Palette.bgElevated,
        child: Icon(Icons.person, size: size * 0.55, color: Palette.textSecondary),
      ),
    ),
  );
}
