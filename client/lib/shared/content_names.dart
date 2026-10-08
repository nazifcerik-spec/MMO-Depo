import '../core/api_client.dart';
import 'widgets.dart';

/// Display names for content codes, served by the API per language (DD-230)
/// and keyed `kind.code` (item, affix, mob, dungeon, region, skill, quest,
/// ...). Loaded once per session; anything missing falls back to the
/// humanised code, so a failed load only costs prettier names.
class ContentNames {
  ContentNames._();

  static Map<String, String> _names = const {};

  /// Replaces the table (tests set it directly).
  static void use(Map<String, String> names) => _names = names;

  static String of(String kind, String code, [String? fallback]) => _names['$kind.$code'] ?? fallback ?? humanize(code);

  /// An item's name with its rolled affixes, in the language's word order.
  static String item(String code, [String? prefix, String? suffix]) {
    final name = of('item', code);
    if (prefix == null && suffix == null) return name;
    final format = _names['affix.format'] ?? '{prefix} {name} {suffix}';
    return format
        .replaceFirst('{prefix}', prefix == null ? '' : of('affix', prefix))
        .replaceFirst('{name}', name)
        .replaceFirst('{suffix}', suffix == null ? '' : of('affix', suffix))
        .replaceAll('()', '')
        .replaceAll(RegExp(r'\s+'), ' ')
        .trim();
  }

  /// Fetches the table for [lang]; failures keep the current one.
  static Future<void> load(ApiClient api, String lang) async {
    try {
      final j = await api.get('/v1/content/strings?lang=$lang');
      final raw = j['strings'];
      if (raw is Map) _names = raw.map((k, v) => MapEntry('$k', '$v'));
    } catch (_) {
      // Names are cosmetic: codes stay readable through humanize.
    }
  }
}
