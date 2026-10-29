import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/app/theme.dart';
import 'package:otter/features/shop/cosmetics.dart';

void main() {
  test('every background has bundled art and a pair of colours (DD-323)', () {
    for (var n = 1; n <= backgroundCount; n++) {
      final code = backgroundCode(n);
      expect(File(backgroundAsset(code)).existsSync(), isTrue, reason: code);
      final (accent, secondary) = backgroundColors(code);
      expect(accent, isNot(secondary), reason: code);
    }
    expect(backgroundAsset('bg99'), 'assets/backgrounds/bg30.webp', reason: 'unknown codes clamp');
  });

  test('themedFor recolours the accent only', () {
    final base = buildTheme();
    final t = themedFor(base, 'bg11');
    expect(t.colorScheme.primary, const Color(0xFFFF7A3D));
    expect(t.colorScheme.secondary, const Color(0xFF8C2F1F));
    expect(t.colorScheme.surface, base.colorScheme.surface);
  });

  testWidgets('BackgroundArt: the art covers, darkened behind the page', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: BackgroundArt(code: 'bg07', child: Text('page')),
      ),
    );
    final img = tester.widget<Image>(find.byType(Image));
    expect((img.image as AssetImage).assetName, 'assets/backgrounds/bg07.webp');
    expect(img.fit, BoxFit.cover);
    expect(find.byType(DecoratedBox), findsOneWidget);
    expect(find.text('page'), findsOneWidget);

    await tester.pumpWidget(const MaterialApp(home: BackgroundArt(code: 'bg07', thumbnail: true)));
    expect(find.byType(DecoratedBox), findsNothing, reason: 'shop tiles show the art as is');
  });
}
