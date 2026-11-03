import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../inventory/item_images.dart';
import '../inventory/item_widgets.dart';
import 'pets_screen.dart';

/// A raid whose boss can drop the species' egg: where the egg comes from.
typedef PetEggSource = ({String raid, String boss, int level, String region});

/// The pet index page of a species (DD-342): its art, rarity, element and
/// creature, story and habitat, strength from level 1 to 100, evolution, and
/// where its egg comes from. An undiscovered species keeps its name, story
/// and art hidden and only says where to look.
Future<void> showPetSpecies(
  BuildContext context,
  PetSpecies s, {
  int powerPerLevel = 10,
  bool discovered = true,
  Future<List<PetEggSource>>? sources,
}) => showModalBottomSheet<void>(
  context: context,
  isScrollControlled: true,
  showDragHandle: true,
  backgroundColor: Palette.bgSurface,
  builder: (_) => PetSpeciesSheet(species: s, powerPerLevel: powerPerLevel, discovered: discovered, sources: sources),
);

class PetSpeciesSheet extends StatelessWidget {
  const PetSpeciesSheet({
    super.key,
    required this.species,
    this.powerPerLevel = 10,
    this.discovered = true,
    this.sources,
  });

  final PetSpecies species;
  final int powerPerLevel;
  final bool discovered;
  final Future<List<PetEggSource>>? sources;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final s = species;
    final rarity = petRarity(s.rank);
    final color = rarityColor(rarity);

    Widget section(String title, List<Widget> children) => Container(
      margin: const EdgeInsets.only(top: Space.md),
      padding: const EdgeInsets.all(Space.md),
      decoration: BoxDecoration(
        color: Palette.bgElevated.withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(Radii.md),
        border: Border.all(color: Palette.borderSubtle),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(title, style: text.titleSmall?.copyWith(color: Palette.accent)),
          const SizedBox(height: Space.sm),
          ...children,
        ],
      ),
    );
    Widget line(IconData icon, String value) => Padding(
      padding: const EdgeInsets.only(bottom: Space.xs),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 18, color: color),
          const SizedBox(width: Space.sm),
          Expanded(child: Text(value)),
        ],
      ),
    );
    Widget chip(String label, {IconData? icon, Color? tint}) => Container(
      padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 4),
      decoration: BoxDecoration(
        color: Palette.bgElevated,
        borderRadius: BorderRadius.circular(Radii.sm),
        border: Border.all(color: (tint ?? Palette.borderStrong).withValues(alpha: 0.7)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[Icon(icon, size: 14, color: tint ?? Palette.accent), const SizedBox(width: 4)],
          Flexible(
            child: Text(label, style: TextStyle(fontSize: 12, color: tint ?? Palette.textPrimary)),
          ),
        ],
      ),
    );
    Widget art(Widget known, double size) => discovered
        ? known
        : Container(
            width: size,
            height: size,
            decoration: BoxDecoration(
              color: Palette.bgPrimary,
              borderRadius: BorderRadius.circular(Radii.md),
              border: Border.all(color: Palette.borderStrong),
            ),
            child: Icon(Icons.help_outline, size: size * 0.4, color: Palette.textMuted),
          );
    Widget stage(Widget image, String name, String levels) => Expanded(
      child: Column(
        children: [
          image,
          const SizedBox(height: Space.xs),
          Text(
            name,
            style: text.bodySmall?.copyWith(color: Palette.textPrimary),
            textAlign: TextAlign.center,
          ),
          Text(levels, style: text.bodySmall, textAlign: TextAlign.center),
        ],
      ),
    );
    const arrow = Padding(
      padding: EdgeInsets.only(bottom: 34),
      child: Icon(Icons.chevron_right, size: 18, color: Palette.textMuted),
    );
    final hatch = shortDuration(l, Duration(seconds: s.hatchSeconds));

    return SafeArea(
      child: PageFrame(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxHeight: MediaQuery.sizeOf(context).height * 0.9),
          child: SingleChildScrollView(
            key: const Key('pet-species-sheet'),
            padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Center(
                  child: art(PetArt(species: s.code, stage: 3, size: 180, color: color), 180),
                ),
                const SizedBox(height: Space.sm),
                Text(
                  discovered ? petName(s.code) : '???',
                  textAlign: TextAlign.center,
                  style: text.titleLarge?.copyWith(color: color, fontWeight: FontWeight.w700),
                ),
                Text('#${s.rank.toString().padLeft(2, '0')}', textAlign: TextAlign.center, style: text.bodySmall),
                const SizedBox(height: Space.sm),
                Wrap(
                  alignment: WrapAlignment.center,
                  spacing: Space.xs,
                  runSpacing: Space.xs,
                  children: [
                    chip(rarityName(l, rarity), icon: Icons.diamond_outlined, tint: color),
                    if (discovered && s.element.isNotEmpty)
                      chip(l.petElement(ContentNames.of('element', s.element)), icon: Icons.auto_awesome),
                    if (discovered) chip(ContentNames.of('pet', '${s.code}.creature', ''), icon: Icons.pets),
                  ],
                ),
                if (!discovered)
                  Padding(
                    padding: const EdgeInsets.only(top: Space.sm),
                    child: Text(l.petUnknown, textAlign: TextAlign.center, style: text.bodySmall),
                  )
                else
                  section(l.petAbout, [
                    Text(
                      ContentNames.of('pet', '${s.code}.story', ''),
                      textAlign: TextAlign.justify,
                      style: text.bodyMedium,
                    ),
                    const SizedBox(height: Space.sm),
                    line(
                      Icons.landscape_outlined,
                      '${l.petHabitat}: ${ContentNames.of('pet', '${s.code}.habitat', '')}',
                    ),
                  ]),
                section(l.petStrength, [
                  Text(l.petLevelNote, style: text.bodySmall),
                  const SizedBox(height: Space.xs),
                  if (s.assistMax > 0) line(Icons.flash_on, l.petAssistRange(s.assistMin, s.assistMax)),
                  if (bonusText(l, s.bonus) case final max?)
                    line(Icons.trending_up, l.petBonusRange(bonusText(l, s.bonusMin ?? s.bonus) ?? max, max)),
                  line(Icons.pets, l.petFightsRange(powerPerLevel, powerPerLevel * 100)),
                ]),
                section(l.petEvolution, [
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      stage(
                        ItemArt(code: 'pet_egg_${s.code}', icon: Icons.egg_outlined, color: color, size: 52),
                        l.petHatching,
                        hatch,
                      ),
                      for (final (n, name, from, to) in [
                        (1, l.petStage1, 1, 33),
                        (2, l.petStage2, 34, 66),
                        (3, l.petStage3, 67, 100),
                      ]) ...[
                        arrow,
                        stage(
                          art(PetArt(species: s.code, stage: n, size: 52, color: color), 52),
                          name,
                          l.petStageLevels(from, to),
                        ),
                      ],
                    ],
                  ),
                ]),
                section(l.petFind, [
                  FutureBuilder<List<PetEggSource>>(
                    future: sources,
                    builder: (context, snap) {
                      final raids = snap.data ?? const <PetEggSource>[];
                      return Wrap(
                        spacing: Space.xs,
                        runSpacing: Space.xs,
                        children: [
                          for (final r in raids) ...[
                            chip(l.petSourceRaid(ContentNames.of('dungeon', r.raid)), icon: Icons.castle_outlined),
                            if (r.boss.isNotEmpty)
                              chip(l.petSourceBoss(ContentNames.of('mob', r.boss), r.level), icon: Icons.whatshot),
                            if (r.region.isNotEmpty)
                              chip(l.petSourceRegion(ContentNames.of('region', r.region)), icon: Icons.place_outlined),
                          ],
                          if (raids.isEmpty) chip(l.petFindLine(s.minMobLevel, s.eggPctText), icon: Icons.groups_3),
                          chip(l.petSourceChance(s.eggPctText), icon: Icons.casino_outlined),
                          chip(l.petSourceHatch(hatch), icon: Icons.egg_outlined),
                        ],
                      );
                    },
                  ),
                ]),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
