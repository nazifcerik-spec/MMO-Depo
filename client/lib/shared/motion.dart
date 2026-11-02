import 'dart:math';

import 'package:flutter/material.dart';

/// One-shot state-change animations (docs/animation-proposals.md, DD-337):
/// each plays once when its value changes, never loops, and is skipped
/// when the platform asks for reduced motion.

/// The duration to animate for, or zero when animations are off.
Duration motion(BuildContext context, Duration d) =>
    MediaQuery.maybeDisableAnimationsOf(context) ?? false ? Duration.zero : d;

/// A soft glow around [child] that flares once whenever [trigger] changes
/// (and on first build when [playOnStart]).
class GlowPulse extends StatefulWidget {
  const GlowPulse({
    super.key,
    required this.trigger,
    required this.color,
    required this.child,
    this.playOnStart = false,
    this.radius = 12,
    this.duration = const Duration(milliseconds: 600),
  });

  final Object? trigger;
  final Color color;
  final Widget child;
  final bool playOnStart;
  final double radius;
  final Duration duration;

  @override
  State<GlowPulse> createState() => _GlowPulseState();
}

class _GlowPulseState extends State<GlowPulse> with SingleTickerProviderStateMixin {
  late final AnimationController _c = AnimationController(vsync: this, duration: widget.duration);

  @override
  void initState() {
    super.initState();
    if (widget.playOnStart) WidgetsBinding.instance.addPostFrameCallback((_) => _play());
  }

  @override
  void didUpdateWidget(GlowPulse old) {
    super.didUpdateWidget(old);
    if (old.trigger != widget.trigger) _play();
  }

  void _play() {
    if (!mounted || motion(context, widget.duration) == Duration.zero) return;
    _c.forward(from: 0);
  }

  @override
  void dispose() {
    _c.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => AnimatedBuilder(
    animation: _c,
    child: widget.child,
    builder: (_, child) {
      // Always the same box, so the child keeps its state between frames.
      final t = _c.isAnimating ? sin(pi * _c.value) : 0.0;
      return DecoratedBox(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(widget.radius),
          boxShadow: t == 0
              ? null
              : [
                  BoxShadow(
                    color: widget.color.withValues(alpha: 0.7 * t),
                    blurRadius: 24 * t,
                    spreadRadius: 4 * t,
                  ),
                ],
        ),
        child: child,
      );
    },
  );
}

/// Cross-fades [child] into its replacement (keyed by [swapKey]) with a
/// slight grow, e.g. an egg hatching or a pet reaching its next stage.
class SwapIn extends StatelessWidget {
  const SwapIn({super.key, required this.swapKey, required this.child});

  final Object swapKey;
  final Widget child;

  @override
  Widget build(BuildContext context) => AnimatedSwitcher(
    duration: motion(context, const Duration(milliseconds: 600)),
    switchInCurve: Curves.easeOutBack,
    transitionBuilder: (child, a) => FadeTransition(
      opacity: a,
      child: ScaleTransition(scale: Tween(begin: 0.85, end: 1.0).animate(a), child: child),
    ),
    child: KeyedSubtree(key: ValueKey(swapKey), child: child),
  );
}

/// Grows [child] from 80 % to full size with a small overshoot when it first
/// appears, e.g. a level-up badge.
class PopIn extends StatelessWidget {
  const PopIn({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) => TweenAnimationBuilder<double>(
    tween: Tween(begin: 0.8, end: 1),
    duration: motion(context, const Duration(milliseconds: 350)),
    curve: Curves.easeOutBack,
    builder: (_, s, c) => Transform.scale(scale: s, child: c),
    child: child,
  );
}

/// Counts from zero up to [value] once, building the text with [builder].
class CountUp extends StatelessWidget {
  const CountUp({super.key, required this.value, required this.builder});

  final int value;
  final Widget Function(int) builder;

  @override
  Widget build(BuildContext context) => TweenAnimationBuilder<double>(
    tween: Tween(begin: 0, end: value.toDouble()),
    duration: motion(context, const Duration(milliseconds: 500)),
    curve: Curves.easeOut,
    builder: (_, v, _) => builder(v.round()),
  );
}
