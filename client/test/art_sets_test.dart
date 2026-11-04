import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

/// Every promotion, specialization, region, crafting profession, helper and
/// the six late materials ship art (DD-332); a new one needs its picture.
void main() {
  Iterable<String> codes(String file, RegExp re) =>
      re.allMatches(File('../content/$file').readAsStringSync()).map((m) => m[1]!);

  void expectArt(String dir, Iterable<String> codes) {
    final missing = [
      for (final c in codes)
        if (!File('assets/$dir/$c.webp').existsSync() && !File('assets/$dir/$c.png').existsSync()) c,
    ];
    expect(missing, isEmpty, reason: dir);
  }

  test('class tree nodes, regions, professions, helpers and materials have art (DD-332)', () {
    final nodes = codes('progression.yaml', RegExp(r'- code: (\w+)\s+parent:')).toList();
    expect(nodes, hasLength(60));
    expectArt('class_nodes', nodes);
    final regions = codes('regions.yaml', RegExp(r'^  - code: (\w+)', multiLine: true)).toList();
    expect(regions, isNotEmpty);
    expectArt('regions', regions);
    expectArt('professions', [
      'blacksmithing',
      'leatherworking',
      'tailoring',
      'woodworking',
      'jewelcrafting',
      'alchemy',
      'cooking',
      'smelting',
      'enchanting',
    ]);
    expectArt('helpers', [
      'miner',
      'woodsman',
      'angler',
      'herbalist',
      'trapper',
      'smith',
      'tanner',
      'weaver',
      'carver',
      'apothecary',
      'cook',
      'smelter',
      'scribe',
      'gemcutter',
    ]);
    expectArt('items', ['elm_log', 'yew_log', 'linseed', 'mandrake', 'king_bloom', 'mercury_ore']);
  });

  test('every enchanting scroll has art (DD-333)', () {
    final scrolls = codes('items.yaml', RegExp(r'code: (enchant_\w+)')).toSet();
    expect(scrolls, hasLength(210));
    expectArt('items', scrolls);
  });

  test('pets, dungeons, quests, gathering trades, class portraits and consumables have art (DD-335)', () {
    final species = codes(
      'pets.yaml',
      RegExp(r'^  - code: (\w+)\n    display_name_key: pet\.', multiLine: true),
    ).toList();
    expect(species, hasLength(50));
    expectArt('pets', [
      for (final s in species)
        for (final stage in [1, 2, 3]) '${s}_$stage',
    ]);
    expectArt('items', species.map((s) => 'pet_egg_$s'));
    expectArt('dungeons', codes('dungeons.yaml', RegExp(r'^  - code: (\w+)', multiLine: true)));
    expectArt('quests', codes('quests.yaml', RegExp(r'^  - code: (\w+)', multiLine: true)));
    expectArt('professions', ['mining', 'herbalism', 'woodcutting', 'fishing', 'hunting']);
    expectArt('classes', [
      'warrior',
      'rogue',
      'ranger',
      'mage',
      'monk',
      'cleric',
      'paladin',
      'druid',
      'bard',
      'shaman',
    ]);
    final consumables = [
      for (final f in ['items.yaml', 'gathering.yaml'])
        for (final block in File('../content/$f').readAsStringSync().split('\n  - code: '))
          if (block.contains('item_class: consumable')) block.split(RegExp(r'\s')).first,
    ];
    expect(consumables, isNotEmpty);
    expectArt('items', consumables);
  });

  // Art in a folder pubspec does not list never reaches the app (the quest
  // medallions shipped that way until DD-338).
  test('every asset folder is bundled', () {
    final pubspec = File('pubspec.yaml').readAsStringSync();
    final missing = [
      for (final d in Directory('assets').listSync().whereType<Directory>())
        if (!pubspec.contains('- assets/${d.uri.pathSegments.where((s) => s.isNotEmpty).last}/') &&
            !d.path.endsWith('held')) // the paper doll's sprites: switched off, not shipped (DD-348)
          d.path,
    ];
    expect(missing, isEmpty);
  });

  test('every monster and class skill has art (DD-336)', () {
    final mobs = codes('mobs.yaml', RegExp(r'^  - code: (\w+)', multiLine: true)).toList();
    expect(mobs, hasLength(146));
    expectArt('mobs', mobs);
    final skills = codes('skills.yaml', RegExp(r'^  - code: (\w+)', multiLine: true)).toList();
    expect(skills, hasLength(471));
    expectArt('skills', skills);
  });
}
