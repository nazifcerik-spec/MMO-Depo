import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app/theme.dart';

/// A number box for a quantity (attribute points, supplies, craft count,
/// listing size): digits only, kept within [min]..[max]. A typed number past
/// [max] snaps to it; an emptied box counts as [min] until something is
/// typed. Tapping it selects the number so it is replaced as you type.
class QuantityField extends StatefulWidget {
  const QuantityField({
    super.key,
    required this.value,
    required this.max,
    required this.onChanged,
    this.min = 0,
    this.width = 64,
    this.semanticLabel,
  });

  final int value;
  final int min;
  final int max;
  final ValueChanged<int> onChanged;
  final double width;
  final String? semanticLabel;

  @override
  State<QuantityField> createState() => _QuantityFieldState();
}

class _QuantityFieldState extends State<QuantityField> {
  late final _ctl = TextEditingController(text: '${widget.value}');

  @override
  void didUpdateWidget(QuantityField old) {
    super.didUpdateWidget(old);
    final shown = int.tryParse(_ctl.text);
    final blankAtMin = _ctl.text.isEmpty && widget.value == widget.min;
    if (shown != widget.value && !blankAtMin) _show(widget.value);
  }

  @override
  void dispose() {
    _ctl.dispose();
    super.dispose();
  }

  void _show(int v) => _ctl.value = TextEditingValue(
    text: '$v',
    selection: TextSelection.collapsed(offset: '$v'.length),
  );

  void _changed(String text) {
    if (text.isEmpty) {
      widget.onChanged(widget.min);
      return;
    }
    final typed = int.parse(text);
    final v = typed.clamp(widget.min, widget.max < widget.min ? widget.min : widget.max);
    if (v != typed) _show(v);
    widget.onChanged(v);
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.width,
      child: Semantics(
        label: widget.semanticLabel,
        textField: true,
        child: TextField(
          controller: _ctl,
          keyboardType: TextInputType.number,
          inputFormatters: [FilteringTextInputFormatter.digitsOnly, LengthLimitingTextInputFormatter(5)],
          textAlign: TextAlign.center,
          style: const TextStyle(fontWeight: FontWeight.w700, fontFeatures: numeric),
          decoration: const InputDecoration(
            isDense: true,
            contentPadding: EdgeInsets.symmetric(horizontal: Space.xs, vertical: Space.sm),
          ),
          onTap: () => _ctl.selection = TextSelection(baseOffset: 0, extentOffset: _ctl.text.length),
          onChanged: _changed,
        ),
      ),
    );
  }
}
