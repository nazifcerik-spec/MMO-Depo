import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import 'models.dart';

/// Customization options (DD-346), mirrored by server character.*Options. New
/// art: a higher count here and there, the layers from scripts/head_art.py.
const faceOptions = 5;
const hairStyleOptions = 5;

/// Multiplied over the grey hair and iris layers; ids are 1-based.
const hairColors = [Color(0xFFEBC878), Color(0xFF4B4B5A), Color(0xFF784B2D), Color(0xFFAA3228), Color(0xFFD2D2E1)];
const eyeColors = [Color(0xFF4682DC), Color(0xFF3CAA5A), Color(0xFF785032), Color(0xFFE6AA3C), Color(0xFF965ADC)];

/// Races with drawn faces and hair (the others keep their race portrait).
const customizableRaces = {'human'};

/// A character's head from its separate layers: face, brows, irises, hair.
/// [onFigure] hides the bust's neck and clothes for the full-body portrait.
class HeadLayers extends StatelessWidget {
  const HeadLayers({super.key, required this.look, required this.size, this.onFigure = false, this.hideHair = false});

  final Appearance look;
  final double size;
  final bool onFigure;

  /// A helm or hood is worn instead of the hair (DD-347).
  final bool hideHair;

  @override
  Widget build(BuildContext context) {
    final cache = (size * MediaQuery.devicePixelRatioOf(context)).round().clamp(32, 256);
    final hair = hairColors[(look.hairColor - 1).clamp(0, hairColors.length - 1)];
    final eyes = eyeColors[(look.eyeColor - 1).clamp(0, eyeColors.length - 1)];
    final face = 'assets/heads/${look.body}_face${look.face}';

    Widget layer(String path, {Color? tint}) => Positioned.fill(
      child: Image.asset(
        '$path.webp',
        cacheWidth: cache,
        color: tint,
        colorBlendMode: tint == null ? null : BlendMode.modulate,
        gaplessPlayback: true,
        filterQuality: FilterQuality.medium,
        errorBuilder: (_, _, _) => const SizedBox.shrink(),
      ),
    );
    Widget fading(Widget child, double from, double to) => Positioned.fill(
      child: ShaderMask(
        blendMode: BlendMode.dstIn,
        shaderCallback: (r) => LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: const [Colors.white, Colors.white, Colors.transparent],
          stops: [0, from, to],
        ).createShader(r),
        child: Stack(children: [child]),
      ),
    );

    final base = layer('${face}_base');
    final top = layer('assets/heads/${look.body}_hair${look.hairStyle}', tint: hair);
    return SizedBox.square(
      dimension: size,
      child: Stack(
        children: [
          // The chin is at about 3/4 of the frame.
          if (onFigure) fading(base, 0.76, 0.84) else base,
          layer('${face}_brows', tint: hair),
          layer('${face}_eyes', tint: eyes),
          if (!hideHair)
            if (onFigure) fading(top, 0.9, 1) else top,
        ],
      ),
    );
  }
}

/// Gender, then for customizable races face, hair style, hair colour and eye
/// colour with a live preview. Every change is reported at once.
class AppearancePicker extends StatelessWidget {
  const AppearancePicker({super.key, required this.look, required this.race, required this.onChanged});

  final Appearance look;
  final String race;
  final ValueChanged<Appearance> onChanged;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final full = customizableRaces.contains(race);

    Widget option(String key, String label, bool selected, VoidCallback onTap, Widget child) => Expanded(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 2),
        child: Semantics(
          label: label,
          selected: selected,
          button: true,
          child: Material(
            color: Palette.bgElevated,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(Radii.sm),
              side: BorderSide(color: selected ? Palette.accent : Palette.borderSubtle, width: selected ? 2 : 1),
            ),
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              key: Key(key),
              onTap: onTap,
              child: AspectRatio(aspectRatio: 1, child: child),
            ),
          ),
        ),
      ),
    );
    Widget row(String title, List<Widget> options) => Padding(
      padding: const EdgeInsets.only(top: Space.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: Theme.of(context).textTheme.bodySmall),
          const SizedBox(height: Space.xs),
          Row(children: options),
        ],
      ),
    );
    Widget heads(String title, String key, int count, int selected, Appearance Function(int) at) => row(title, [
      for (var n = 1; n <= count; n++)
        option(
          '$key-$n',
          '$title $n',
          n == selected,
          () => onChanged(at(n)),
          LayoutBuilder(
            builder: (_, box) => HeadLayers(look: at(n), size: box.maxWidth),
          ),
        ),
    ]);
    Widget colors(String title, String key, List<Color> list, int selected, Appearance Function(int) at) => row(title, [
      for (var n = 1; n <= list.length; n++)
        option(
          '$key-$n',
          '$title $n',
          n == selected,
          () => onChanged(at(n)),
          Padding(
            padding: const EdgeInsets.all(Space.sm),
            child: DecoratedBox(
              decoration: BoxDecoration(color: list[n - 1], shape: BoxShape.circle),
            ),
          ),
        ),
    ]);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (full) ...[
          Center(
            child: Container(
              key: const Key('look-preview'),
              decoration: BoxDecoration(color: Palette.bgElevated, borderRadius: BorderRadius.circular(Radii.md)),
              clipBehavior: Clip.antiAlias,
              child: HeadLayers(look: look, size: 168),
            ),
          ),
          const SizedBox(height: Space.md),
        ],
        SegmentedButton<String>(
          key: const Key('body-choice'),
          showSelectedIcon: false,
          segments: [
            ButtonSegment(value: 'male', label: Text(l.bodyMale), icon: const Icon(Icons.male, size: 18)),
            ButtonSegment(value: 'female', label: Text(l.bodyFemale), icon: const Icon(Icons.female, size: 18)),
          ],
          selected: {look.body},
          onSelectionChanged: (s) => onChanged(look.withBody(s.first)),
        ),
        if (full) ...[
          heads(l.lookFace, 'look-face', faceOptions, look.face, (n) => look.copyWith(face: n)),
          heads(l.lookHairStyle, 'look-style', hairStyleOptions, look.hairStyle, (n) => look.copyWith(hairStyle: n)),
          colors(l.lookHairColor, 'look-hair', hairColors, look.hairColor, (n) => look.copyWith(hairColor: n)),
          colors(l.lookEyeColor, 'look-eyes', eyeColors, look.eyeColor, (n) => look.copyWith(eyeColor: n)),
        ] else
          Padding(
            padding: const EdgeInsets.only(top: Space.sm),
            child: Text(
              l.lookHumanOnly,
              textAlign: TextAlign.center,
              style: const TextStyle(fontSize: 11, color: Palette.textMuted),
            ),
          ),
      ],
    );
  }
}

/// Edits an existing character's look; pops with the new one when saved.
Future<Appearance?> showAppearanceEditor(BuildContext context, Appearance look, String race) =>
    showModalBottomSheet<Appearance>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (context) {
        var current = look;
        return StatefulBuilder(
          builder: (context, setState) => SafeArea(
            child: ConstrainedBox(
              constraints: BoxConstraints(maxHeight: MediaQuery.sizeOf(context).height * 0.9),
              child: SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    AppearancePicker(look: current, race: race, onChanged: (v) => setState(() => current = v)),
                    const SizedBox(height: Space.lg),
                    FilledButton(
                      key: const Key('look-save'),
                      onPressed: () => Navigator.of(context).pop(current),
                      child: Text(AppLocalizations.of(context).lookSave),
                    ),
                  ],
                ),
              ),
            ),
          ),
        );
      },
    );
