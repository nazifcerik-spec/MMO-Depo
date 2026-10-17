import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/// Item art bundled as `assets/items/<template_code>.png` (optional per item),
/// and hunting animals as `assets/animals/<animal>.png`. The lists are read
/// once from the asset manifest, so a missing picture costs no failed load —
/// the item keeps its class icon.
class ItemImages {
  static Set<String> _codes = const {};
  static Set<String> _animals = const {};

  static const _dir = 'assets/items/';
  static const _animalDir = 'assets/animals/';

  static Future<void> load([AssetBundle? bundle]) async {
    try {
      final manifest = await AssetManifest.loadFromAssetBundle(bundle ?? rootBundle);
      Set<String> under(String dir) => {
        for (final a in manifest.listAssets())
          if (a.startsWith(dir) && a.endsWith('.png')) a.substring(dir.length, a.length - 4),
      };
      _codes = under(_dir);
      _animals = under(_animalDir);
    } catch (_) {
      _codes = const {}; // no art is a valid state: icons everywhere
      _animals = const {};
    }
  }

  static String? path(String code) => _codes.contains(code) ? '$_dir$code.png' : null;

  static String? animalPath(String animal) => _animals.contains(animal) ? '$_animalDir$animal.png' : null;

  @visibleForTesting
  static set animals(Set<String> a) => _animals = a;

  @visibleForTesting
  static set codes(Set<String> c) => _codes = c;
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
