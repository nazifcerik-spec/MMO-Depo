import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

/// DD-361: every interface language carries every string with the same
/// placeholders, so nothing silently shows up in English.
void main() {
  Map<String, dynamic> arb(String lang) =>
      jsonDecode(File('lib/l10n/app_$lang.arb').readAsStringSync()) as Map<String, dynamic>;
  final placeholder = RegExp(r'\{([A-Za-z_]\w*)[,}]');
  Set<String> names(String s) => placeholder.allMatches(s).map((m) => m[1]!).toSet();

  final en = arb('en');
  final keys = en.keys.where((k) => !k.startsWith('@')).toList();

  for (final lang in ['tr', 'es', 'zh', 'ru', 'ko', 'ja']) {
    test('$lang has every string and placeholder', () {
      final table = arb(lang);
      expect(keys.where((k) => !table.containsKey(k)), isEmpty);
      for (final k in keys) {
        expect(names(table[k] as String), names(en[k] as String), reason: k);
      }
    });
  }
}
