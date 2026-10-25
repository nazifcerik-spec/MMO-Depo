import 'dart:async';

import 'package:flutter/material.dart';

/// An icon button that steps once on tap and keeps stepping, faster and
/// faster, while held. [onStep] null disables it; a held press stops once
/// the parent passes null (nothing left to add or remove).
class HoldRepeatButton extends StatefulWidget {
  const HoldRepeatButton({super.key, required this.icon, required this.tooltip, required this.onStep});

  final IconData icon;
  final String tooltip;
  final VoidCallback? onStep;

  @override
  State<HoldRepeatButton> createState() => _HoldRepeatButtonState();
}

class _HoldRepeatButtonState extends State<HoldRepeatButton> {
  Timer? _timer;
  int _repeats = 0;
  bool _held = false; // the press already stepped: its tap must not step again

  void _start() {
    _stop();
    _held = false; // a new press (a held one may end without a tap)
    _timer = Timer(const Duration(milliseconds: 350), _tick);
  }

  void _tick() {
    final step = widget.onStep;
    if (step == null || !mounted) return _stop();
    _held = true;
    step();
    _repeats++;
    // 8/s at first, 25/s after a second of holding.
    _timer = Timer(Duration(milliseconds: _repeats < 8 ? 120 : 40), _tick);
  }

  void _stop() {
    _timer?.cancel();
    _timer = null;
    _repeats = 0;
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final step = widget.onStep;
    return Listener(
      onPointerDown: step == null ? null : (_) => _start(),
      onPointerUp: (_) => _stop(),
      onPointerCancel: (_) => _stop(),
      child: IconButton(
        tooltip: widget.tooltip,
        icon: Icon(widget.icon),
        onPressed: step == null
            ? null
            : () {
                if (_held) {
                  _held = false;
                  return;
                }
                step();
              },
      ),
    );
  }
}
