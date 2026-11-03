import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import '../inventory/item_images.dart';
import '../inventory/item_widgets.dart';
import 'pets_screen.dart';

/// A species' card from the pet guide (DD-341): who it is, how it evolves,
/// how strong it gets and how to find it — the essentials, nothing more.
Future<void> showPetSpecies(BuildContext context, PetSpecies s, {int powerPerLevel = 10}) => showModalBottomSheet<void>(
  context: context,
  isScrollControlled: true,
  showDragHandle: true,
  backgroundColor: Palette.bgSurface,
  builder: (_) => PetSpeciesSheet(species: s, powerPerLevel: powerPerLevel),
);

class PetSpeciesSheet extends StatelessWidget {
  const PetSpeciesSheet({super.key, required this.species, this.powerPerLevel = 10});

  final PetSpecies species;
  final int powerPerLevel;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final s = species;
    final rarity = petRarity(s.rank);
    final color = rarityColor(rarity);

    Widget section(String title, List<Widget> children) => Padding(
      padding: const EdgeInsets.only(top: Space.lg),
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
    Widget stage(Widget art, String label) => Expanded(
      child: Column(
        children: [
          art,
          const SizedBox(height: Space.xs),
          Text(label, style: text.bodySmall, textAlign: TextAlign.center),
        ],
      ),
    );
    const arrow = Padding(
      padding: EdgeInsets.only(bottom: 18),
      child: Icon(Icons.chevron_right, size: 18, color: Palette.textMuted),
    );

    return SafeArea(
      child: PageFrame(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxHeight: MediaQuery.sizeOf(context).height * 0.85),
          child: ListView(
            key: const Key('pet-species-sheet'),
            shrinkWrap: true,
            padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
            children: [
              Row(
                children: [
                  PetArt(species: s.code, stage: 3, size: 96, color: color),
                  const SizedBox(width: Space.lg),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          petName(s.code),
                          style: text.titleLarge?.copyWith(color: color, fontWeight: FontWeight.w700),
                        ),
                        const SizedBox(height: Space.xs),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 2),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(Radii.sm),
                            border: Border.all(color: color),
                          ),
                          child: Text(
                            rarityName(l, rarity),
                            style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.w600),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              section(l.petEvolution, [
                Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    stage(
                      ItemArt(code: 'pet_egg_${s.code}', icon: Icons.egg_outlined, color: color, size: 52),
                      l.petHatching,
                    ),
                    for (final (n, from, to) in [(1, 1, 33), (2, 34, 66), (3, 67, 100)]) ...[
                      arrow,
                      stage(PetArt(species: s.code, stage: n, size: 52, color: color), l.petStageLevels(from, to)),
                    ],
                  ],
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
              section(l.petFind, [
                line(Icons.groups_3, l.petFindLine(s.minMobLevel, s.eggPctText)),
                line(Icons.egg_outlined, l.petHatchTime(shortDuration(l, Duration(seconds: s.hatchSeconds)))),
              ]),
            ],
          ),
        ),
      ),
    );
  }
}
