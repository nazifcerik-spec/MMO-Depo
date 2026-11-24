import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Persists only the refresh token. The access token lives in memory and is
/// re-minted from the refresh token on cold start.
abstract class TokenStore {
  Future<String?> readRefreshToken();
  Future<void> writeRefreshToken(String token);
  Future<void> clear();
}

/// Keychain (iOS) / EncryptedSharedPreferences (Android) / WebCrypto (web).
/// A browser may refuse it (private window, plain-http address): the session
/// then lives only as long as the page.
class SecureTokenStore implements TokenStore {
  SecureTokenStore([FlutterSecureStorage? storage]) : _storage = storage ?? const FlutterSecureStorage();

  static const _key = 'otter.refresh_token';
  final FlutterSecureStorage _storage;

  @override
  Future<String?> readRefreshToken() async {
    try {
      return await _storage.read(key: _key);
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> writeRefreshToken(String token) async {
    try {
      await _storage.write(key: _key, value: token);
    } catch (_) {
      // Kept in memory by ApiClient.
    }
  }

  @override
  Future<void> clear() async {
    try {
      await _storage.delete(key: _key);
    } catch (_) {
      // Nothing stored.
    }
  }
}

/// In-memory store for tests.
class MemoryTokenStore implements TokenStore {
  String? token;

  @override
  Future<String?> readRefreshToken() async => token;

  @override
  Future<void> writeRefreshToken(String t) async => token = t;

  @override
  Future<void> clear() async => token = null;
}
