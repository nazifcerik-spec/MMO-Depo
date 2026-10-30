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
        if (!File('assets/$dir/$c.webp').existsSync()) c,
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
}
