import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:flutter/material.dart';
import 'package:otter/features/character/appearance_picker.dart';
import 'package:otter/features/character/body_portrait.dart';
import 'package:otter/features/character/models.dart';
import 'package:otter/features/inventory/inventory_models.dart';
import 'package:otter/l10n/app_localizations.dart';

EquippedItem chest(String subclass, int level) => EquippedItem(
  instanceId: 'c',
  templateCode: 'x',
  rarity: 'rare',
  itemClass: 'armor',
  itemSubclass: subclass,
  minLevel: level,
);

void main() {
  test('the chest piece picks the outfit: armor type and level band (DD-344)', () {
    expect(bodyOutfit(null), 'none');
    expect(bodyOutfit(chest('cloth_armor', 1)), 'cloth1');
    expect(bodyOutfit(chest('leather_armor', 99)), 'leather1');
    expect(bodyOutfit(chest('leather_armor', 100)), 'leather2');
    expect(bodyOutfit(chest('plate_armor', 640)), 'plate4');
    expect(bodyOutfit(chest('plate_armor', 1000)), 'plate6');
    expect(bodyOutfit(chest('jewelry', 300)), 'none', reason: 'not body armor');
  });

  test('every outfit of the drawn figures ships', () {
    for (final body in ['male', 'female']) {
      for (final o in [
        'none',
        for (final t in ['cloth', 'leather', 'plate'])
          for (var b = 1; b <= bodyBands.length; b++) '$t$b',
      ]) {
        expect(File('assets/bodies/human_${body}_$o.webp').existsSync(), isTrue, reason: 'human $body $o');
      }
    }
  });

  test('another body keeps the colours but not the face and hair style (DD-346)', () {
    const look = Appearance(face: 3, hairStyle: 4, hairColor: 2, eyeColor: 5);
    final other = look.withBody('female');
    expect([other.body, other.face, other.hairStyle, other.hairColor, other.eyeColor], ['female', 1, 1, 2, 5]);
    expect(identical(look.withBody('male'), look), isTrue);
    expect(look.toJson(), {'body': 'male', 'face_id': 3, 'hair_style_id': 4, 'hair_color_id': 2, 'eye_color_id': 5});
    expect(Appearance.fromJson(look.toJson()).hairStyle, 4);
    expect(Appearance.fromJson(const {}).face, 1, reason: 'an older server');
  });

  test('every layer of the 25 + 25 combinations ships', () {
    for (final body in ['male', 'female']) {
      for (var n = 1; n <= faceOptions; n++) {
        for (final part in ['base', 'brows', 'eyes']) {
          expect(File('assets/heads/${body}_face${n}_$part.webp').existsSync(), isTrue, reason: '$body face $n $part');
        }
      }
      for (var n = 1; n <= hairStyleOptions; n++) {
        expect(File('assets/heads/${body}_hair$n.webp').existsSync(), isTrue, reason: '$body hair $n');
      }
    }
  });

  testWidgets('the picker offers every option and reports each choice at once', (tester) async {
    Appearance look = const Appearance();
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        locale: const Locale('en'),
        home: Scaffold(
          body: StatefulBuilder(
            builder: (_, setState) => SingleChildScrollView(
              child: AppearancePicker(look: look, race: 'human', onChanged: (v) => setState(() => look = v)),
            ),
          ),
        ),
      ),
    );
    for (final key in ['look-face', 'look-style', 'look-hair', 'look-eyes']) {
      for (var n = 1; n <= 5; n++) {
        expect(find.byKey(Key('$key-$n')), findsOneWidget, reason: '$key $n');
      }
    }
    await tester.tap(find.byKey(const Key('look-face-3')));
    await tester.pump();
    await tester.ensureVisible(find.byKey(const Key('look-eyes-4')));
    await tester.tap(find.byKey(const Key('look-eyes-4')));
    await tester.pump();
    expect([look.face, look.eyeColor], [3, 4]);
    await tester.ensureVisible(find.text('Female'));
    await tester.tap(find.text('Female'));
    await tester.pump();
    expect([look.body, look.face, look.eyeColor], ['female', 1, 4]);
  });
}
