import 'package:flutter/material.dart';

import '../app/theme.dart';

/// Wide-screen layout (DD-369): browsers and tablets from [wideBreakpoint]
/// get a side rail and panes; below it everything stays the phone layout.
const wideBreakpoint = 900.0;

/// Width of a column of cards: phones, then secondary pages on wide screens.
const columnWidth = 520.0;
const wideColumnWidth = 720.0;

/// Width of the in-game tabs on wide screens.
const gameWidth = 1160.0;

/// Width of the settings / equipment pane beside a tab's main list.
const sidePaneWidth = 380.0;

bool isWide(BuildContext context) => MediaQuery.sizeOf(context).width >= wideBreakpoint;

/// Page column, centred: the full width on phones, [columnWidth] on tablets;
/// on wide screens [wideColumnWidth], or [maxWidth] when given.
class PageFrame extends StatelessWidget {
  const PageFrame({super.key, required this.child, this.maxWidth});

  final Widget child;
  final double? maxWidth;

  @override
  Widget build(BuildContext context) => Align(
    alignment: Alignment.topCenter,
    child: ConstrainedBox(
      constraints: BoxConstraints(maxWidth: isWide(context) ? maxWidth ?? wideColumnWidth : columnWidth),
      child: child,
    ),
  );
}

/// Same-type cards: one column on phones, as many equal columns as fit on
/// wide screens. Cards of a row share its height, so all look alike; [spacing]
/// separates the rows (0 when the cards bring their own bottom margin).
class TileGrid extends StatelessWidget {
  const TileGrid({super.key, required this.children, this.minTileWidth = 320, this.spacing = Space.sm});

  final List<Widget> children;
  final double minTileWidth;
  final double spacing;

  @override
  Widget build(BuildContext context) {
    if (!isWide(context)) return _column(children);
    return LayoutBuilder(
      builder: (context, box) {
        final cols = ((box.maxWidth + Space.md) ~/ (minTileWidth + Space.md)).clamp(1, 4);
        if (cols == 1) return _column(children);
        return _column([
          for (var i = 0; i < children.length; i += cols)
            IntrinsicHeight(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  for (var c = 0; c < cols; c++) ...[
                    if (c > 0) const SizedBox(width: Space.md),
                    Expanded(child: i + c < children.length ? children[i + c] : const SizedBox()),
                  ],
                ],
              ),
            ),
        ]);
      },
    );
  }

  Widget _column(List<Widget> rows) => Column(
    crossAxisAlignment: CrossAxisAlignment.stretch,
    children: [
      for (final (i, r) in rows.indexed) ...[if (i > 0 && spacing > 0) SizedBox(height: spacing), r],
    ],
  );
}

/// A list and its settings: side by side on wide screens ([after] under the
/// list), one column on phones ([main], [side], [after]). [sideFirst] puts the
/// side pane on top there and on the left on wide screens.
class Panes extends StatelessWidget {
  const Panes({
    super.key,
    required this.main,
    required this.side,
    this.after = const [],
    this.sideFirst = false,
    this.sideWidth = sidePaneWidth,
  });

  final List<Widget> main;
  final List<Widget> side;
  final List<Widget> after;
  final bool sideFirst;
  final double sideWidth;

  @override
  Widget build(BuildContext context) {
    Widget column(List<Widget> c) => Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: c);
    if (!isWide(context)) {
      return column([...sideFirst ? side : main, ...sideFirst ? main : side, ...after]);
    }
    final mainPane = Expanded(child: column([...main, ...after]));
    final sidePane = SizedBox(width: sideWidth, child: column(side));
    const gap = SizedBox(width: Space.xl);
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: sideFirst ? [sidePane, gap, mainPane] : [mainPane, gap, sidePane],
    );
  }
}

/// Keeps a single card (a running job) from stretching across a wide tab.
class NarrowColumn extends StatelessWidget {
  const NarrowColumn({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) => Center(
    child: ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: columnWidth),
      child: child,
    ),
  );
}
