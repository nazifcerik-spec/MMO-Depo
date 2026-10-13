import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/inventory/item_images.dart';
import 'package:otter/features/inventory/item_widgets.dart';

void main() {
  tearDown(() => ItemImages.codes = const {});

  testWidgets('an item without shipped art keeps its class icon', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: ItemCell(rarity: 'common', icon: Icons.pets_outlined, code: 'boar_hide'),
      ),
    );
    expect(find.byIcon(Icons.pets_outlined), findsOneWidget);
    expect(find.byType(Image), findsNothing);
  });

  testWidgets('an item whose <code>.png ships shows the picture', (tester) async {
    ItemImages.codes = {'boar_hide'};
    expect(ItemImages.path('boar_hide'), 'assets/items/boar_hide.png');
    expect(ItemImages.path('wolf_pelt'), isNull);
    await tester.pumpWidget(
      const MaterialApp(
        home: ItemCell(rarity: 'common', icon: Icons.pets_outlined, code: 'boar_hide'),
      ),
    );
    final image = tester.widget<Image>(find.byType(Image));
    expect((image.image as ResizeImage).imageProvider, isA<AssetImage>());
  });

  test('the manifest lists what ships under assets/items/', () async {
    TestWidgetsFlutterBinding.ensureInitialized();
    await ItemImages.load();
    expect(ItemImages.path('definitely_not_an_item'), isNull);
  });
}
