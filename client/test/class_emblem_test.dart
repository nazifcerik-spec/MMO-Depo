import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/character/class_catalog.dart';
import 'package:otter/features/shop/shop_screen.dart';

const _bases = ['warrior', 'rogue', 'ranger', 'mage', 'monk', 'cleric', 'paladin', 'druid', 'bard', 'shaman'];

void main() {
  test('every base class has bundled emblem art (DD-325)', () {
    for (final c in _bases) {
      expect(File(classEmblemAsset(c)).existsSync(), isTrue, reason: c);
    }
  });

  test('every XP boost has bundled emblem art (DD-326)', () {
    for (final c in ['gather', 'adventure', 'craft']) {
      expect(File(boostAsset(c)).existsSync(), isTrue, reason: c);
    }
  });

  testWidgets('ClassEmblem: the art in a circle at its size', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Center(child: ClassEmblem(root: 'druid', size: 40)),
      ),
    );
    final img = tester.widget<Image>(find.byType(Image));
    expect(((img.image as ResizeImage).imageProvider as AssetImage).assetName, 'assets/class_icons/druid.webp');
    expect(find.byType(ClipOval), findsOneWidget);
    expect(tester.getSize(find.byType(ClassEmblem)), const Size(40, 40));
  });
}
