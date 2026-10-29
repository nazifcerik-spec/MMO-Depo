import 'package:flutter/material.dart';

/// Cosmetics drawn in code (DD-321): 30 page backgrounds and 80 avatars.
/// Codes match the server's (bg01..bg30, av01..av80); the first 10
/// backgrounds and 40 avatars are free. Art can replace them later under
/// the same codes.
const backgroundCount = 30, freeBackgrounds = 10, avatarCount = 80, freeAvatars = 40;

String backgroundCode(int n) => 'bg${n.toString().padLeft(2, '0')}';
String avatarCode(int n) => 'av${n.toString().padLeft(2, '0')}';

int _number(String code) => int.tryParse(code.substring(2)) ?? 1;

bool isFree(String code) => code.startsWith('bg') ? _number(code) <= freeBackgrounds : _number(code) <= freeAvatars;

/// Background n: a dark two-tone gradient walking the colour wheel; the
/// premium ones (11..30) get a third, brighter stop.
Gradient backgroundGradient(String code) {
  final n = _number(code).clamp(1, backgroundCount);
  final hue = (n * 37) % 360.0;
  final a = HSLColor.fromAHSL(1, hue, 0.45, 0.09).toColor();
  final b = HSLColor.fromAHSL(1, (hue + 40) % 360, 0.55, 0.16).toColor();
  return LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: n <= freeBackgrounds ? [a, b] : [a, b, HSLColor.fromAHSL(1, (hue + 80) % 360, 0.6, 0.24).toColor()],
  );
}

const _avatarIcons = <IconData>[
  Icons.shield,
  Icons.auto_awesome,
  Icons.local_fire_department,
  Icons.ac_unit,
  Icons.bolt,
  Icons.eco,
  Icons.pets,
  Icons.castle,
  Icons.diamond,
  Icons.star,
  Icons.nightlight,
  Icons.wb_sunny,
  Icons.water_drop,
  Icons.terrain,
  Icons.anchor,
  Icons.catching_pokemon,
  Icons.psychology,
  Icons.sports_martial_arts,
  Icons.music_note,
  Icons.spa,
];

const _avatarColors = <Color>[Color(0xFFE4A756), Color(0xFF5DA9E9), Color(0xFF5FB973), Color(0xFFA26BE8)];

/// Avatar n: one of 20 emblems in one of 4 colours.
class AvatarBadge extends StatelessWidget {
  const AvatarBadge({super.key, required this.code, this.size = 40});

  final String code;
  final double size;

  @override
  Widget build(BuildContext context) {
    final n = _number(code).clamp(1, avatarCount) - 1;
    final color = _avatarColors[n ~/ _avatarIcons.length % _avatarColors.length];
    return CircleAvatar(
      radius: size / 2,
      backgroundColor: color.withValues(alpha: 0.18),
      child: Icon(_avatarIcons[n % _avatarIcons.length], color: color, size: size * 0.55),
    );
  }
}
