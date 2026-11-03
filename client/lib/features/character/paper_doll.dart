import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';

/// How one piece of equipment is drawn on the figure (DD-347). Fractions are
/// of the portrait frame's width (x, scale) and height (y); [pivot] is the
/// point of the sprite that lands on the anchor and that it turns around.
class DollLayer {
  const DollLayer({
    required this.equipmentId,
    required this.slot,
    required this.avatarAsset,
    required this.anchor,
    this.pivot = const Offset(0.5, 0.5),
    this.offsetX = 0,
    this.offsetY = 0,
    this.scale = 0.5,
    this.rotation = 0,
    this.zIndex = 0,
    this.behind = false,
    this.flipX = false,
    this.hidesHair = false,
    this.mirrorAnchor,
  });

  final String equipmentId;
  final String slot;
  final String avatarAsset;

  /// Position of the named body anchor, in frame fractions.
  final Offset anchor;
  final Offset pivot;
  final double offsetX;
  final double offsetY;
  final double scale;

  /// Degrees, clockwise.
  final double rotation;
  final int zIndex;

  /// renderingLayer `back`: under the figure instead of over it.
  final bool behind;
  final bool flipX;

  /// A helm or hood replaces the hair layer.
  final bool hidesHair;

  /// Paired gear (gloves): drawn again, mirrored, on this anchor.
  final Offset? mirrorAnchor;
}

/// The paper doll metadata (`assets/paperdoll.json`, written by
/// scripts/paperdoll_editor.py). A piece's rule merges, later winning: its
/// slot, `<slot>.<subclass>` (a type of gear), the item itself. Without a rule nothing is drawn, so
/// new items need no code and a missing or broken file only hides the layers.
class PaperDoll {
  PaperDoll(Map<String, dynamic> data)
    : _anchors = _map(data['anchors']),
      _anchorOverrides = _map(data['anchorOverrides']),
      _slots = _map(data['slots']),
      _subclasses = _map(data['subclasses']),
      _items = _map(data['items']);

  static const assetPath = 'assets/paperdoll.json';

  /// The loaded metadata; empty until [load] (no layers).
  static PaperDoll current = PaperDoll(const {});

  static Future<void> load([AssetBundle? bundle]) async {
    try {
      final raw = await (bundle ?? rootBundle).loadString(assetPath);
      current = PaperDoll(jsonDecode(raw) as Map<String, dynamic>);
    } catch (e) {
      debugPrint('paper doll metadata not loaded: $e');
    }
  }

  final Map<String, dynamic> _anchors;
  final Map<String, dynamic> _anchorOverrides;
  final Map<String, dynamic> _slots;
  final Map<String, dynamic> _subclasses;
  final Map<String, dynamic> _items;

  static Map<String, dynamic> _map(Object? v) => v is Map<String, dynamic> ? v : const {};

  static Offset? _point(Object? v) => v is List && v.length == 2 && v[0] is num && v[1] is num
      ? Offset((v[0] as num).toDouble(), (v[1] as num).toDouble())
      : null;

  /// The anchor on this figure: `<race>_<body>` first, then `<body>`, then
  /// the shared one, so a differently built race only overrides what moves.
  Offset? anchorOf(String name, {String race = '', String body = ''}) =>
      _point(_map(_anchorOverrides['${race}_$body'])[name]) ??
      _point(_map(_anchorOverrides[body])[name]) ??
      _point(_anchors[name]);

  /// The piece's layer on this figure; null when it draws nothing: no rule,
  /// `renderingLayer: none`, or an anchor this file does not know.
  DollLayer? layerOf({
    required String equipmentId,
    required String slot,
    String subclass = '',
    String race = '',
    String body = '',
  }) {
    final rule = <String, dynamic>{
      ..._map(_slots[slot]),
      ..._map(_subclasses['$slot.$subclass']),
      ..._map(_items[equipmentId]),
    };
    if (rule.isEmpty || rule['renderingLayer'] == 'none') return null;
    final asset = rule['avatarAsset'];
    final anchor = anchorOf('${rule['anchorPoint']}', race: race, body: body);
    if (asset is! String || anchor == null) return null;
    double n(String k, double fallback) => (rule[k] as num?)?.toDouble() ?? fallback;
    return DollLayer(
      equipmentId: equipmentId,
      slot: slot,
      avatarAsset: asset.replaceAll('{id}', equipmentId),
      anchor: anchor,
      pivot: _point(rule['pivot']) ?? const Offset(0.5, 0.5),
      offsetX: n('offsetX', 0),
      offsetY: n('offsetY', 0),
      scale: n('scale', 0.5),
      rotation: n('rotation', 0),
      zIndex: (rule['zIndex'] as num?)?.toInt() ?? 0,
      behind: rule['renderingLayer'] == 'back',
      flipX: rule['flipX'] == true,
      hidesHair: rule['hidesHair'] == true,
      mirrorAnchor: rule['mirrorAt'] == null ? null : anchorOf('${rule['mirrorAt']}', race: race, body: body),
    );
  }
}
