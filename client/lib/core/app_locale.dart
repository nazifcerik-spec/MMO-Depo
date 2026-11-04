import 'package:flutter/widgets.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// The game's languages (DD-354), each named in itself.
const appLanguages = <(String, String)>[
  ('tr', 'Türkçe'),
  ('en', 'English'),
  ('ru', 'Русский'),
  ('ko', '한국어'),
  ('zh', '中文'),
  ('ja', '日本語'),
  ('es', 'Español'),
];

/// The language the player chose; null follows the device. Kept on the device.
class AppLocale extends ValueNotifier<Locale?> {
  AppLocale._() : super(null);

  static final instance = AppLocale._();
  static const _key = 'mena.language';
  static const _storage = FlutterSecureStorage();

  Future<void> load() async {
    try {
      final code = await _storage.read(key: _key);
      if (code != null && appLanguages.any((l) => l.$1 == code)) value = Locale(code);
    } catch (_) {
      // Unreadable storage: follow the device.
    }
  }

  Future<void> choose(String code) async {
    value = Locale(code);
    try {
      await _storage.write(key: _key, value: code);
    } catch (_) {
      // Applies for this run.
    }
  }
}
