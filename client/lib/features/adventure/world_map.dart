import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import 'adventure_models.dart';

/// The travel map (DD-279): every region at its content position with its
/// level above it, continents as tinted bands, roads between linked regions.
/// Regions above the character's level are dimmed; tapping any region
/// reports it to [onTap]. A journey under way is drawn along its [route].
class WorldMapView extends StatefulWidget {
  const WorldMapView({
    super.key,
    required this.map,
    required this.here,
    required this.level,
    required this.onTap,
    this.route = const [],
  });

  final WorldMap map;

  /// Where the character stands.
  final String here;

  final int level;

  /// The regions of a journey under way, origin first (DD-280).
  final List<String> route;
  final ValueChanged<RegionInfo> onTap;

  @override
  State<WorldMapView> createState() => _WorldMapViewState();
}

/// Band tints, one per continent in map order.
const _bandColors = [Palette.success, Palette.accent, Palette.info, Palette.epic, Palette.warning, Palette.danger];
const _scale = 0.85; // opening zoom: two map columns on a phone

class _WorldMapViewState extends State<WorldMapView> {
  final _view = TransformationController();
  String? _centeredOn;

  @override
  void dispose() {
    _view.dispose();
    super.dispose();
  }

  /// Opens on the character's region (again after it moves).
  void _center(Size viewport) {
    if (_centeredOn == widget.here) return;
    final r = widget.map.regions[widget.here];
    if (r == null) return;
    _centeredOn = widget.here;
    _view.value = Matrix4.identity()
      ..translateByDouble(viewport.width / 2 - r.mapX * _scale, viewport.height / 2 - r.mapY * _scale, 0, 1)
      ..scaleByDouble(_scale, _scale, 1, 1);
  }

  @override
  Widget build(BuildContext context) {
    final regions = widget.map.regions.values.toList();
    if (regions.isEmpty) return const SizedBox.shrink();
    final width = regions.map((r) => r.mapX).reduce((a, b) => a > b ? a : b) + 180;
    final height = regions.map((r) => r.mapY).reduce((a, b) => a > b ? a : b) + 120;
    return LayoutBuilder(
      builder: (context, box) {
        _center(box.biggest);
        return ClipRect(
          child: InteractiveViewer(
            transformationController: _view,
            constrained: false,
            minScale: 0.3,
            maxScale: 2,
            boundaryMargin: const EdgeInsets.all(240),
            child: SizedBox(
              width: width,
              height: height,
              child: Stack(
                children: [
                  Positioned.fill(
                    child: CustomPaint(painter: _MapPainter(widget.map, route: widget.route)),
                  ),
                  for (final r in regions) _node(context, r),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _node(BuildContext context, RegionInfo r) {
    final size = switch (r.kind) {
      'continent' => 44.0,
      'dungeon' => 26.0,
      'city' => 40.0,
      _ => 34.0,
    };
    const width = 128.0, levelLine = 16.0;
    return Positioned(
      left: r.mapX - width / 2,
      top: r.mapY - levelLine - size / 2,
      width: width,
      child: _MapNode(
        region: r,
        size: size,
        here: r.code == widget.here,
        locked: r.minLevel > widget.level,
        onTap: () => widget.onTap(r),
      ),
    );
  }
}

class _MapNode extends StatelessWidget {
  const _MapNode({
    required this.region,
    required this.size,
    required this.here,
    required this.locked,
    required this.onTap,
  });

  final RegionInfo region;
  final double size;
  final bool here;
  final bool locked;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final name = ContentNames.of('region', region.code);
    final ring = here
        ? Palette.accent
        : locked
        ? Palette.danger.withValues(alpha: 0.6)
        : Palette.success;
    final icon = switch (region.kind) {
      'city' => Icons.location_city,
      'dungeon' => Icons.castle_outlined,
      'continent' => Icons.public,
      'wilderness' => Icons.forest,
      _ => Icons.terrain,
    };
    return Semantics(
      button: true,
      label: '$name, ${l.levelShort(region.minLevel)}',
      selected: here,
      child: GestureDetector(
        key: Key('map-${region.code}'),
        behavior: HitTestBehavior.opaque,
        onTap: onTap,
        child: Opacity(
          opacity: locked ? 0.5 : 1,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              SizedBox(
                height: 16,
                child: Text(
                  l.levelShort(region.minLevel),
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                    color: locked ? Palette.danger : Palette.textSecondary,
                  ),
                ),
              ),
              Container(
                width: size,
                height: size,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: Palette.bgElevated,
                  border: Border.all(color: ring, width: here ? 3 : 2),
                  boxShadow: [if (here) BoxShadow(color: Palette.accent.withValues(alpha: 0.5), blurRadius: 12)],
                ),
                child: Icon(
                  locked ? Icons.lock_outline : icon,
                  size: size * 0.5,
                  color: here ? Palette.accent : Palette.textPrimary,
                ),
              ),
              const SizedBox(height: 2),
              Text(
                name,
                textAlign: TextAlign.center,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  fontSize: region.kind == 'dungeon' ? 10 : 12,
                  fontWeight: here || region.kind == 'continent' ? FontWeight.w700 : FontWeight.w500,
                  color: here ? Palette.accent : Palette.textPrimary,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Continent bands, roads (links and zone–dungeon paths) and the route of a
/// journey under way.
class _MapPainter extends CustomPainter {
  _MapPainter(this.map, {this.route = const []});

  final WorldMap map;
  final List<String> route;

  @override
  void paint(Canvas canvas, Size size) {
    final byContinent = <String, List<RegionInfo>>{};
    for (final r in map.regions.values) {
      byContinent.putIfAbsent(map.continentOf(r.code), () => []).add(r);
    }
    final continents = byContinent.keys.toList()
      ..sort((a, b) => (map.regions[a]?.mapY ?? 0).compareTo(map.regions[b]?.mapY ?? 0));
    for (final (i, c) in continents.indexed) {
      final members = byContinent[c]!;
      final xs = members.map((r) => r.mapX), ys = members.map((r) => r.mapY);
      final rect = Rect.fromLTRB(
        xs.reduce((a, b) => a < b ? a : b) - 80,
        ys.reduce((a, b) => a < b ? a : b) - 40,
        xs.reduce((a, b) => a > b ? a : b) + 80,
        ys.reduce((a, b) => a > b ? a : b) + 70,
      );
      final color = _bandColors[i % _bandColors.length];
      final band = RRect.fromRectAndRadius(rect, const Radius.circular(28));
      canvas.drawRRect(band, Paint()..color = color.withValues(alpha: 0.07));
      canvas.drawRRect(
        band,
        Paint()
          ..style = PaintingStyle.stroke
          ..strokeWidth = 1.5
          ..color = color.withValues(alpha: 0.3),
      );
    }

    final road = Paint()
      ..color = Palette.borderStrong
      ..strokeWidth = 4
      ..strokeCap = StrokeCap.round;
    for (final (a, b) in map.links) {
      final ra = map.regions[a], rb = map.regions[b];
      if (ra != null && rb != null) canvas.drawLine(Offset(ra.mapX, ra.mapY), Offset(rb.mapX, rb.mapY), road);
    }
    final path = Paint()
      ..color = Palette.borderSubtle
      ..strokeWidth = 2;
    for (final r in map.regions.values) {
      final p = r.kind == 'dungeon' ? map.regions[r.parent] : null;
      if (p != null) canvas.drawLine(Offset(p.mapX, p.mapY), Offset(r.mapX, r.mapY), path);
    }

    final stops = [
      for (final code in route)
        if (map.regions[code] case final r?) Offset(r.mapX, r.mapY),
    ];
    final ink = Paint()
      ..color = Palette.accent
      ..strokeWidth = 5
      ..strokeCap = StrokeCap.round;
    const dash = 14.0;
    for (var i = 0; i + 1 < stops.length; i++) {
      final a = stops[i], b = stops[i + 1], total = (b - a).distance;
      for (var d = 0.0; d < total; d += dash * 2) {
        canvas.drawLine(a + (b - a) * (d / total), a + (b - a) * ((d + dash).clamp(0, total) / total), ink);
      }
    }
  }

  @override
  bool shouldRepaint(_MapPainter old) => old.map != map || old.route != route;
}
