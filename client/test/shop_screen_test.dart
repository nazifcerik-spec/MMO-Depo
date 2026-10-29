import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/shop/shop_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

Widget _app(Widget home) => MaterialApp(
  localizationsDelegates: AppLocalizations.localizationsDelegates,
  supportedLocales: AppLocalizations.supportedLocales,
  home: Scaffold(body: home),
);

void main() {
  testWidgets('shop: credits buy a boost and a premium avatar, then wear it (DD-321)', (tester) async {
    final be = FakeBackend();
    String? worn;
    await tester.pumpWidget(
      _app(ShopScreen(repo: ShopRepository(be.api(), 'char-1'), onProfile: (p) => worn = p.avatar)),
    );
    await tester.pumpAndSettle();
    expect(find.text('50 credits'), findsOneWidget);
    expect(find.text('Buying credits with money is coming soon.'), findsOneWidget);

    await tester.tap(find.byKey(const Key('buy-boost_gather')));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Buy'));
    await tester.pumpAndSettle();
    expect(find.text('30 credits'), findsOneWidget);

    await tester.tap(find.text('Avatars'));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const Key('cos-av07'))); // free: worn at once
    await tester.pumpAndSettle();
    expect(worn, 'av07');
    await tester.drag(find.byType(GridView).first, const Offset(0, -900));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const Key('cos-av41')));
    await tester.pumpAndSettle();
    expect(find.text('Spend 15 credits?'), findsOneWidget);
    await tester.tap(find.widgetWithText(FilledButton, 'Buy'));
    await tester.pumpAndSettle();
    expect(be.ownedCosmetics, ['av41']);
    expect(worn, 'av41');
    expect(be.credits, 15);
  });

  testWidgets('boost bar: three boxes, watch an ad, the box turns active (DD-320)', (tester) async {
    final be = FakeBackend();
    await tester.pumpWidget(_app(BoostBar(repo: ShopRepository(be.api(), 'char-1'))));
    await tester.pumpAndSettle();
    expect(find.text('Gathering'), findsOneWidget);
    expect(find.text('Combat'), findsOneWidget);
    expect(find.text('Professions'), findsOneWidget);
    await tester.tap(find.text('Watch ad (3 left)').first);
    await tester.pumpAndSettle();
    expect(find.textContaining('Active,'), findsOneWidget);
  });
}
