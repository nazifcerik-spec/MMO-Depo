import 'package:flutter/material.dart';

/// Design tokens from docs/ui/design-system.md (dark-first; light is Phase 02).
abstract final class Palette {
  static const bgPrimary = Color(0xFF0F1013);
  static const bgSurface = Color(0xFF1A1C22);
  static const bgElevated = Color(0xFF242731);
  static const borderSubtle = Color(0xFF2E3240);
  static const borderStrong = Color(0xFF4A4F5E);
  static const textPrimary = Color(0xFFE8EAEE);
  static const textSecondary = Color(0xFFA8ACB8);
  static const textMuted = Color(0xFF6E727E);
  static const accent = Color(0xFFE4A756);
  static const info = Color(0xFF5DA9E9);
  static const success = Color(0xFF5FB973);
  static const warning = Color(0xFFE8B84B);
  static const danger = Color(0xFFE85B62);
  static const epic = Color(0xFFA26BE8);
}

/// Rarity accent (border, aura, label).
Color rarityColor(String rarity) => switch (rarity) {
  'uncommon' => Palette.success,
  'rare' => Palette.info,
  'epic' => Palette.epic,
  'legendary' => Palette.accent,
  'mythic' => Palette.danger,
  _ => const Color(0xFFB9BAC1),
};

abstract final class Space {
  static const xs = 4.0;
  static const sm = 8.0;
  static const md = 12.0;
  static const lg = 16.0;
  static const screen = 20.0;
  static const xl = 24.0;
  static const xxl = 32.0;
}

abstract final class Radii {
  static const sm = 6.0;
  static const md = 10.0;
  static const lg = 16.0;
}

/// Inline (row / list-trailing) buttons: the theme's full-width minimum
/// would be infinite there.
final compactButton = FilledButton.styleFrom(
  minimumSize: const Size(72, 40),
  padding: const EdgeInsets.symmetric(horizontal: 16),
);

/// Tabular figures so numbers don't jitter as they change.
const numeric = [FontFeature.tabularFigures()];

ThemeData buildTheme() {
  const scheme = ColorScheme.dark(
    primary: Palette.accent,
    onPrimary: Palette.bgPrimary,
    secondary: Palette.info,
    onSecondary: Palette.bgPrimary,
    error: Palette.danger,
    surface: Palette.bgSurface,
    onSurface: Palette.textPrimary,
    surfaceContainerHighest: Palette.bgElevated,
    outline: Palette.borderStrong,
    outlineVariant: Palette.borderSubtle,
  );
  final base = ThemeData(colorScheme: scheme, useMaterial3: true, brightness: Brightness.dark);
  final border = OutlineInputBorder(
    borderRadius: BorderRadius.circular(Radii.md),
    borderSide: const BorderSide(color: Palette.borderSubtle),
  );
  return base.copyWith(
    scaffoldBackgroundColor: Palette.bgPrimary,
    textTheme: base.textTheme
        .apply(bodyColor: Palette.textPrimary, displayColor: Palette.textPrimary)
        .merge(
          const TextTheme(
            headlineMedium: TextStyle(fontSize: 28, fontWeight: FontWeight.w700, height: 34 / 28),
            titleLarge: TextStyle(fontSize: 22, fontWeight: FontWeight.w700, height: 28 / 22),
            titleMedium: TextStyle(fontSize: 18, fontWeight: FontWeight.w600, height: 24 / 18),
            bodyMedium: TextStyle(fontSize: 14, height: 20 / 14),
            bodySmall: TextStyle(fontSize: 12, height: 16 / 12, color: Palette.textSecondary),
            labelSmall: TextStyle(fontSize: 10, fontWeight: FontWeight.w500, height: 14 / 10),
          ),
        ),
    appBarTheme: const AppBarTheme(
      backgroundColor: Palette.bgPrimary,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      centerTitle: false,
    ),
    cardTheme: CardThemeData(
      color: Palette.bgSurface,
      elevation: 0,
      margin: EdgeInsets.zero,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(Radii.lg),
        side: const BorderSide(color: Palette.borderSubtle),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Palette.bgSurface,
      contentPadding: const EdgeInsets.symmetric(horizontal: Space.lg, vertical: 14),
      border: border,
      enabledBorder: border,
      focusedBorder: border.copyWith(borderSide: const BorderSide(color: Palette.accent, width: 1.5)),
    ),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        minimumSize: const Size.fromHeight(52),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(Radii.md)),
        textStyle: const TextStyle(fontSize: 15, fontWeight: FontWeight.w700, letterSpacing: 0.3),
      ),
    ),
    snackBarTheme: const SnackBarThemeData(behavior: SnackBarBehavior.floating),
  );
}
