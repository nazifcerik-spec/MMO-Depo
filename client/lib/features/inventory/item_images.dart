import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/// Item art bundled as `assets/items/<template_code>.png` or `.webp` (painted
/// gear icons, DD-329; optional per item),
/// hunting animals as `assets/animals/<animal>`, class and race portraits
/// as `assets/classes/<class>` and `assets/races/<race>` (DD-266; `.png` or
/// `.webp`).
/// The lists are read
/// once from the asset manifest, so a missing picture costs no failed load —
/// the item keeps its class icon.
class ItemImages {
  static Map<String, String> _codes = const {}; // template code → asset
  static Map<String, String> _animals = const {};
  static Map<String, String> _classes = const {};
  static Map<String, String> _races = const {};
  static Map<String, String> _bodies = const {};
  static Map<String, String> _held = const {};

  static const _dir = 'assets/items/';
  static const _animalDir = 'assets/animals/';
  static const _classDir = 'assets/classes/';
  static const _raceDir = 'assets/races/';
  static const _bodyDir = 'assets/bodies/';
  static const _heldDir = 'assets/held/';

  static Future<void> load([AssetBundle? bundle]) async {
    try {
      final manifest = await AssetManifest.loadFromAssetBundle(bundle ?? rootBundle);
      Map<String, String> under(String dir) => {
        for (final a in manifest.listAssets())
          if (a.startsWith(dir) && (a.endsWith('.png') || a.endsWith('.webp')))
            a.substring(dir.length, a.lastIndexOf('.')): a,
      };
      _codes = under(_dir);
      _animals = under(_animalDir);
      _classes = under(_classDir);
      _races = under(_raceDir);
      _bodies = under(_bodyDir);
      _held = under(_heldDir);
    } catch (_) {
      _codes = const {}; // no art is a valid state: icons everywhere
      _animals = const {};
      _classes = const {};
      _races = const {};
      _bodies = const {};
      _held = const {};
    }
  }

  static String? path(String code) => _codes[code];

  static String? animalPath(String animal) => _animals[animal];

  static String? classPath(String code) => _classes[code];

  static String? racePath(String code) => _races[code];

  /// Full-body portrait `<race>_<body>_<outfit>` (DD-344).
  static String? bodyPath(String key) => _bodies[key];

  /// A weapon or shield cut out for the portrait's hands (DD-344).
  static String? heldPath(String code) => _held[code];

  @visibleForTesting
  static set bodies(Set<String> b) => _bodies = {for (final x in b) x: '$_bodyDir$x.webp'};

  @visibleForTesting
  static set animals(Set<String> a) => _animals = {for (final x in a) x: '$_animalDir$x.png'};

  @visibleForTesting
  static set codes(Set<String> c) => _codes = {for (final code in c) code: '$_dir$code.png'};
}

/// The item's picture when one ships, else its class icon.
class ItemArt extends StatelessWidget {
  const ItemArt({super.key, required this.code, required this.icon, this.color, this.size = 22});

  final String code;
  final IconData icon;
  final Color? color;
  final double size;

  @override
  Widget build(BuildContext context) {
    final fallback = Icon(icon, color: color, size: size);
    final p = ItemImages.path(code);
    if (p == null) return fallback;
    final box = size * 1.5;
    return Image.asset(
      p,
      width: box,
      height: box,
      // Decode at display size, not the file's (memory on long bag grids).
      cacheWidth: (box * MediaQuery.devicePixelRatioOf(context)).round(),
      filterQuality: FilterQuality.medium,
      errorBuilder: (_, _, _) => fallback,
    );
  }
}
