import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/shared/motion.dart';

/// One-shot animations (DD-337): they settle on their final state, play only
/// on change, and are skipped when the platform asks for reduced motion.
void main() {
  Widget host(Widget child, {bool reduce = false}) => MaterialApp(
    home: MediaQuery(
      data: MediaQueryData(disableAnimations: reduce),
      child: Scaffold(body: Center(child: child)),
    ),
  );

  BoxDecoration? glow(WidgetTester t) =>
      t
              .widget<DecoratedBox>(find.descendant(of: find.byType(GlowPulse), matching: find.byType(DecoratedBox)))
              .decoration
          as BoxDecoration?;

  testWidgets('GlowPulse flares on change and settles without a shadow', (t) async {
    await t.pumpWidget(host(const GlowPulse(trigger: 1, color: Colors.amber, child: SizedBox.square(dimension: 20))));
    expect(glow(t)?.boxShadow, isNull);
    await t.pumpWidget(host(const GlowPulse(trigger: 2, color: Colors.amber, child: SizedBox.square(dimension: 20))));
    await t.pump(const Duration(milliseconds: 300));
    expect(glow(t)?.boxShadow, isNotEmpty);
    await t.pumpAndSettle();
    expect(glow(t)?.boxShadow, isNull);
  });

  testWidgets('GlowPulse stays still with reduced motion', (t) async {
    await t.pumpWidget(
      host(const GlowPulse(trigger: 1, color: Colors.amber, child: SizedBox.square(dimension: 20)), reduce: true),
    );
    await t.pumpWidget(
      host(const GlowPulse(trigger: 2, color: Colors.amber, child: SizedBox.square(dimension: 20)), reduce: true),
    );
    await t.pump(const Duration(milliseconds: 300));
    expect(glow(t)?.boxShadow, isNull);
  });

  testWidgets('CountUp ends on its value; instantly with reduced motion', (t) async {
    await t.pumpWidget(host(CountUp(value: 250, builder: (v) => Text('$v'))));
    expect(find.text('0'), findsOneWidget);
    await t.pumpAndSettle();
    expect(find.text('250'), findsOneWidget);

    await t.pumpWidget(host(CountUp(value: 90, builder: (v) => Text('$v')), reduce: true));
    await t.pump();
    expect(find.text('90'), findsOneWidget);
  });

  testWidgets('SwapIn shows the new child once the swap settles', (t) async {
    await t.pumpWidget(host(const SwapIn(swapKey: 1, child: Text('egg'))));
    await t.pumpWidget(host(const SwapIn(swapKey: 2, child: Text('pet'))));
    await t.pumpAndSettle();
    expect(find.text('pet'), findsOneWidget);
    expect(find.text('egg'), findsNothing);
  });

  testWidgets('PopIn settles at full size', (t) async {
    await t.pumpWidget(host(const PopIn(child: Text('Level 5'))));
    await t.pumpAndSettle();
    final scale = t.widget<Transform>(find.ancestor(of: find.text('Level 5'), matching: find.byType(Transform)).first);
    expect(scale.transform.getMaxScaleOnAxis(), closeTo(1, 1e-6));
  });
}
