import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/character/appearance_picker.dart';
import 'package:otter/features/character/body_portrait.dart';
import 'package:otter/features/character/character_repository.dart';
import 'package:otter/features/character/character_sheet_view.dart';
import 'package:otter/features/character/models.dart';
import 'package:otter/features/character/paper_doll.dart';
import 'package:otter/features/inventory/inventory_models.dart';
import 'package:otter/features/inventory/item_images.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

const _rules = {
  'anchors': {
    'hand': [0.2, 0.5],
    'other': [0.8, 0.5],
    'head': [0.5, 0.1],
  },
  'anchorOverrides': {
    'female': {
      'hand': [0.3, 0.5],
    },
    'dwarf_male': {
      'hand': [0.1, 0.7],
    },
  },
  'slots': {
    'mainhand': {'avatarAsset': 'assets/held/{id}.webp', 'anchorPoint': 'hand', 'scale': 0.5, 'zIndex': 30},
    'cape': {'avatarAsset': 'assets/capes/{id}.webp', 'anchorPoint': 'head', 'renderingLayer': 'back'},
    'gloves': {'avatarAsset': 'assets/held/{id}.webp', 'anchorPoint': 'hand', 'mirrorAt': 'other'},
    'boots': {'renderingLayer': 'none'},
    'broken': {'avatarAsset': 'x.webp', 'anchorPoint': 'nowhere'},
  },
  'subclasses': {
    'mainhand.axe': {'scale': 0.6, 'rotation': -45},
  },
  'items': {
    'axe_gold': {'scale': 0.7, 'avatarAsset': 'assets/special/gold.webp'},
  },
};

EquippedItem _item(String code, {String subclass = '', String rarity = 'common'}) =>
    EquippedItem(instanceId: code, templateCode: code, rarity: rarity, itemClass: 'weapon', itemSubclass: subclass);

Widget _app(Widget child) => MaterialApp(
  localizationsDelegates: AppLocalizations.localizationsDelegates,
  supportedLocales: AppLocalizations.supportedLocales,
  locale: const Locale('en'),
  home: Scaffold(body: child),
);

Map<String, dynamic> _shipped() => jsonDecode(File(PaperDoll.assetPath).readAsStringSync()) as Map<String, dynamic>;

/// Equippable items of the content files as (slot, subclass, code).
List<(String, String, String)> _gear() {
  final out = <(String, String, String)>[];
  for (final f in ['items.yaml', 'gear.yaml', 'sets.yaml']) {
    String code = '', sub = '';
    for (final line in File('../content/$f').readAsLinesSync()) {
      final c = RegExp(r'^  - code: (\w+)').firstMatch(line);
      if (c != null) (code, sub) = (c[1]!, '');
      final s = RegExp(r'^    item_subclass: (\w+)').firstMatch(line);
      if (s != null) sub = s[1]!;
      final e = RegExp(r'^    equip_slot: (\w+)').firstMatch(line);
      if (e != null && code.isNotEmpty) out.add((e[1]!, sub, code));
    }
  }
  return out;
}

void main() {
  final doll = PaperDoll(_rules);

  group('rules', () {
    test('slot, type and item merge with the item winning', () {
      final plain = doll.layerOf(equipmentId: 'sword_a', slot: 'mainhand', subclass: 'sword')!;
      expect([plain.scale, plain.rotation, plain.avatarAsset], [0.5, 0, 'assets/held/sword_a.webp']);
      final axe = doll.layerOf(equipmentId: 'axe_a', slot: 'mainhand', subclass: 'axe')!;
      expect([axe.scale, axe.rotation, axe.zIndex], [0.6, -45, 30]);
      final gold = doll.layerOf(equipmentId: 'axe_gold', slot: 'mainhand', subclass: 'axe')!;
      expect([gold.scale, gold.rotation, gold.avatarAsset], [0.7, -45, 'assets/special/gold.webp']);
    });

    test('a type belongs to its slot', () {
      final d = PaperDoll({
        ..._rules,
        'subclasses': {
          'helm.plate_armor': {'scale': 0.9},
        },
      });
      expect(d.layerOf(equipmentId: 'g', slot: 'gloves', subclass: 'plate_armor')!.scale, 0.5);
    });

    test('no rule, "none" or an unknown anchor draw nothing', () {
      expect(doll.layerOf(equipmentId: 'ring', slot: 'ring'), isNull);
      expect(doll.layerOf(equipmentId: 'b', slot: 'boots'), isNull);
      expect(doll.layerOf(equipmentId: 'b', slot: 'broken'), isNull);
      expect(PaperDoll(const {}).layerOf(equipmentId: 'a', slot: 'mainhand'), isNull);
    });

    test('anchors follow the figure: race and body, body, shared', () {
      Offset at(String race, String body) =>
          doll.layerOf(equipmentId: 'a', slot: 'mainhand', race: race, body: body)!.anchor;
      expect(at('human', 'male'), const Offset(0.2, 0.5));
      expect(at('human', 'female'), const Offset(0.3, 0.5));
      expect(at('dwarf', 'male'), const Offset(0.1, 0.7));
    });

    test('paired gear and the back layer', () {
      expect(doll.layerOf(equipmentId: 'g', slot: 'gloves')!.mirrorAnchor, const Offset(0.8, 0.5));
      expect(doll.layerOf(equipmentId: 'c', slot: 'cape')!.behind, isTrue);
    });

    test('a broken file leaves the doll empty, the app running', () async {
      PaperDoll.current = PaperDoll(const {});
      await PaperDoll.load(_Broken());
      expect(PaperDoll.current.layerOf(equipmentId: 'a', slot: 'mainhand'), isNull);
    });
  });

  group('shipped metadata', () {
    test('the golden axe sits in the right hand with its own picture', () {
      final d = PaperDoll(_shipped()).layerOf(equipmentId: 'axe_sunforged', slot: 'mainhand', subclass: 'axe')!;
      expect(d.avatarAsset, 'assets/held/axe_sunforged.webp');
      expect(d.anchor, const Offset(0.24, 0.55));
      expect(File(d.avatarAsset).existsSync(), isTrue);
      expect(d.zIndex, greaterThanOrEqualTo(headZ));
    });

    test('every piece either draws an existing sprite or is deliberately undrawn', () {
      final d = PaperDoll(_shipped());
      final gear = _gear();
      expect(gear.length, greaterThan(1900));
      final drawn = <String>{};
      for (final (slot, sub, code) in gear) {
        final layer = d.layerOf(equipmentId: code, slot: slot, subclass: sub);
        if (layer == null) continue;
        drawn.add(slot);
        expect(layer.avatarAsset, isNot(contains('assets/items/')), reason: 'icons are not sprites: $code');
        expect(File(layer.avatarAsset).existsSync(), isTrue, reason: '$slot $code: ${layer.avatarAsset}');
      }
      expect(drawn, {'mainhand', 'offhand', 'helm', 'gloves', 'necklace'});
    });

    test('rule keys name real slots and types, anchors exist', () {
      final raw = _shipped();
      final gear = _gear();
      final slots = gear.map((g) => g.$1).toSet();
      final types = gear.map((g) => '${g.$1}.${g.$2}').toSet();
      expect(slots, containsAll((raw['slots'] as Map).keys));
      expect(types, containsAll((raw['subclasses'] as Map).keys));
      expect(gear.map((g) => g.$3).toSet(), containsAll((raw['items'] as Map).keys));
    });
  });

  group('portrait', () {
    setUp(() {
      PaperDoll.current = doll;
      ItemImages.bodies = {'human_male_none'};
    });
    tearDown(() => PaperDoll.current = PaperDoll(const {}));

    Widget portrait(Map<String, EquippedItem> gear, {String race = 'human'}) =>
        _app(BodyPortrait(race: race, look: const Appearance(), equipped: gear));

    testWidgets('equipping and removing updates the figure at once', (tester) async {
      await tester.pumpWidget(portrait(const {}));
      expect(find.byKey(const ValueKey('doll-mainhand')), findsNothing);
      await tester.pumpWidget(portrait({'mainhand': _item('axe_gold', subclass: 'axe')}));
      expect(find.byKey(const ValueKey('doll-mainhand')), findsOneWidget);
      expect(
        tester
            .widget<Image>(
              find.descendant(of: find.byKey(const ValueKey('doll-mainhand')), matching: find.byType(Image)),
            )
            .image,
        isA<ResizeImage>(),
      );
      await tester.pumpWidget(portrait(const {}));
      expect(find.byKey(const ValueKey('doll-mainhand')), findsNothing);
    });

    testWidgets('several pieces show together, mirrored ones twice, back under the figure', (tester) async {
      await tester.pumpWidget(
        portrait({
          'mainhand': _item('axe_a', subclass: 'axe'),
          'gloves': _item('g'),
          'cape': _item('c'),
          'boots': _item('b'),
        }),
      );
      expect(find.byKey(const ValueKey('doll-mainhand')), findsOneWidget);
      expect(find.byKey(const ValueKey('doll-gloves')), findsOneWidget);
      expect(find.byKey(const ValueKey('doll-gloves-mirror')), findsOneWidget);
      expect(find.byKey(const ValueKey('doll-boots')), findsNothing);
      final stack = tester.widget<Stack>(
        find.descendant(of: find.byKey(const Key('body-portrait')), matching: find.byType(Stack)).first,
      );
      expect((stack.children.first.key as ValueKey).value, 'doll-cape', reason: 'the back layer is under the figure');
      final keys = [for (final c in stack.children) c.key];
      expect(keys.indexOf(const ValueKey('doll-gloves')), lessThan(keys.indexOf(const ValueKey('doll-mainhand'))));
    });

    testWidgets('a missing sprite draws nothing and raises nothing', (tester) async {
      await tester.pumpWidget(portrait({'mainhand': _item('no_such_item')}));
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
      expect(find.byKey(const Key('body-portrait')), findsOneWidget);
    });

    testWidgets('a race without figures shows its portrait and no layers', (tester) async {
      await tester.pumpWidget(portrait({'mainhand': _item('axe_a')}, race: 'orc'));
      expect(find.byKey(const ValueKey('doll-mainhand')), findsNothing);
    });

    testWidgets('a helm that hides hair removes the hair layer', (tester) async {
      PaperDoll.current = PaperDoll(_shipped());
      await tester.pumpWidget(portrait({'helm': _item('helm_iron', subclass: 'plate_armor')}));
      expect(tester.widget<HeadLayers>(find.byType(HeadLayers)).hideHair, isTrue);
      await tester.pumpWidget(portrait(const {}));
      expect(tester.widget<HeadLayers>(find.byType(HeadLayers)).hideHair, isFalse);
    });
  });

  testWidgets('the character tab follows equip and unequip from the server', (tester) async {
    PaperDoll.current = PaperDoll(_shipped());
    ItemImages.bodies = {'human_male_none'};
    addTearDown(() => PaperDoll.current = PaperDoll(const {}));
    final be = FakeBackend()
      ..characters.add({
        'id': 'char-1',
        'name': 'Grim',
        'race_code': 'human',
        'class_code': 'warrior',
        'level': 3,
        'state': 'active',
      });
    final api = be.api(MemoryTokenStore());
    await api.login('grim', 'correct-horse');
    final repo = CharacterRepository(api);
    final version = ValueNotifier(0);
    final sheet = await tester.runAsync(() => repo.sheet('char-1'));
    await tester.pumpWidget(_app(CharacterSheetView(repo: repo, initial: sheet!, dataVersion: version)));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('doll-mainhand')), findsNothing);

    be.equipped['mainhand'] = {
      'instance_id': 'a',
      'template_code': 'axe_sunforged',
      'rarity': 'epic',
      'item_class': 'weapon',
      'item_subclass': 'axe',
    };
    version.value++;
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('doll-mainhand')), findsOneWidget);

    be.equipped.remove('mainhand');
    version.value++;
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('doll-mainhand')), findsNothing);
  });
}

class _Broken extends AssetBundle {
  @override
  Future<ByteData> load(String key) async => throw StateError('no asset');

  @override
  Future<String> loadString(String key, {bool cache = true}) async => '{not json';
}
