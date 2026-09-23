import 'dart:async';
import 'dart:convert';
import 'dart:math';

import 'package:http/http.dart' as http;

import 'token_store.dart';

/// Server error envelope `{"error": {code, message, field}}` or a transport
/// failure (`status == 0`, code `E_NETWORK`). Screens never see raw exceptions:
/// controllers catch [ApiError] and render it via `errorText`.
class ApiError implements Exception {
  const ApiError(this.status, this.code, [this.message = '', this.field]);

  const ApiError.network() : this(0, 'E_NETWORK');

  final int status;
  final String code;
  final String message;
  final String? field;

  bool get isNetwork => status == 0;

  /// The session is gone; the auth gate is already routing to login.
  bool get isSessionError => status == 401;

  @override
  String toString() => 'ApiError($status $code $message)';
}

/// Thin JSON client for the Otter API (docs/architecture/flutter.md §5):
/// bearer auth with single-flight refresh, Idempotency-Key on every POST,
/// X-Request-Id per call, timeouts, and bounded retry for GETs only.
class ApiClient {
  ApiClient({
    required this.baseUrl,
    required this.store,
    http.Client? client,
    this.timeout = const Duration(seconds: 10),
    this.retryDelay = const Duration(milliseconds: 300),
  }) : _http = client ?? http.Client();

  final String baseUrl;
  final TokenStore store;
  final Duration timeout;
  final Duration retryDelay;
  final http.Client _http;
  final _rand = Random.secure();

  String? _access;
  DateTime? _accessExpires;
  String? _refresh;
  Future<void>? _refreshing;

  /// Account of the current session, null when signed out.
  String? accountId;

  /// False until the account confirms its e-mail (DD-219); refreshed with
  /// every session response.
  bool emailVerified = true;

  /// Sent as `X-Character-Id` once a character has been entered.
  String? characterId;

  /// Fired when the refresh token is rejected (expired, revoked, reused).
  void Function()? onSessionLost;

  bool get hasSession => _refresh != null;

  /// A non-expiring access token for channels that authenticate once at
  /// connect time (the WebSocket handshake takes it as `?token=`).
  Future<String> freshAccessToken() async {
    if (_refresh == null) throw const ApiError(401, 'E_SESSION_EXPIRED');
    if (_accessExpiringSoon) await _refreshSession();
    return _access!;
  }

  // --- session ---------------------------------------------------------------

  /// Restores a persisted session. Returns false when there is none or the
  /// server rejected it; throws [ApiError.network] when offline so the caller
  /// can offer a retry instead of logging the player out.
  Future<bool> restore() async {
    _refresh = await store.readRefreshToken();
    if (_refresh == null) return false;
    try {
      await _refreshSession();
      return true;
    } on ApiError catch (e) {
      if (e.isNetwork) rethrow;
      return false;
    }
  }

  Future<void> login(String emailOrUsername, String password) async {
    final body = await post('/v1/auth/login', {
      'email_or_username': emailOrUsername,
      'password': password,
    }, auth: false);
    await _applySession(body);
  }

  Future<void> register(String email, String username, String password) async {
    final body = await post('/v1/auth/register', {
      'email': email,
      'username': username,
      'password': password,
    }, auth: false);
    await _applySession(body);
  }

  /// Best-effort server revoke; the local session is always cleared.
  Future<void> logout() async {
    try {
      if (_access != null) await post('/v1/auth/logout', const {});
    } on ApiError {
      // Offline or already revoked — nothing else to do.
    }
    await _clearSession();
  }

  Future<void> _applySession(Map<String, dynamic> body) async {
    _access = body['access_token'] as String;
    _accessExpires = DateTime.parse(body['access_expires_at'] as String);
    _refresh = body['refresh_token'] as String;
    accountId = body['account_id'] as String?;
    emailVerified = body['email_verified'] as bool? ?? true;
    await store.writeRefreshToken(_refresh!);
  }

  Future<void> _clearSession() async {
    _access = null;
    _accessExpires = null;
    _refresh = null;
    accountId = null;
    characterId = null;
    await store.clear();
  }

  /// Single-flight: concurrent 401s share one refresh, otherwise the second
  /// call would present an already-rotated token and trip reuse detection.
  Future<void> _refreshSession() {
    return _refreshing ??= () async {
      try {
        final token = _refresh;
        if (token == null) throw const ApiError(401, 'E_SESSION_EXPIRED');
        final body = await post('/v1/auth/refresh', {'refresh_token': token}, auth: false);
        await _applySession(body);
      } on ApiError catch (e) {
        if (!e.isNetwork) {
          await _clearSession();
          onSessionLost?.call();
        }
        rethrow;
      } finally {
        _refreshing = null;
      }
    }();
  }

  // --- requests --------------------------------------------------------------

  Future<Map<String, dynamic>> get(String path) => _send('GET', path, null, auth: true);

  Future<Map<String, dynamic>> post(String path, Map<String, dynamic> body, {bool auth = true}) =>
      _send('POST', path, body, auth: auth);

  Future<Map<String, dynamic>> delete(String path) => _send('DELETE', path, null, auth: true);

  Future<Map<String, dynamic>> _send(
    String method,
    String path,
    Map<String, dynamic>? body, {
    required bool auth,
  }) async {
    // One key per logical call: a replay after a token refresh must be
    // recognised by the server as the same request (DD-188).
    final idempotencyKey = method == 'GET' ? null : uuidV4(_rand);
    final attempts = method == 'GET' ? 3 : 1;
    var refreshed = false;

    for (var attempt = 1; ; attempt++) {
      if (auth && _refresh != null && _accessExpiringSoon) await _refreshSession();

      final http.Response res;
      try {
        res = await _raw(method, path, body, auth: auth, idempotencyKey: idempotencyKey);
      } on ApiError {
        if (attempt < attempts) {
          await _backoff(attempt);
          continue;
        }
        rethrow;
      }

      // A 401 about the request itself (wrong current password) is not a
      // token problem: refreshing would only rotate tokens for nothing.
      if (res.statusCode == 401 && auth && !refreshed && _refresh != null && !_isCredentialError(res)) {
        refreshed = true;
        await _refreshSession();
        continue;
      }
      if (_retryable(res.statusCode) && attempt < attempts) {
        await _backoff(attempt);
        continue;
      }
      return _decode(res);
    }
  }

  bool get _accessExpiringSoon =>
      _access == null ||
      (_accessExpires != null && _accessExpires!.isBefore(DateTime.now().toUtc().add(const Duration(seconds: 30))));

  Future<http.Response> _raw(
    String method,
    String path,
    Map<String, dynamic>? body, {
    required bool auth,
    String? idempotencyKey,
  }) async {
    final req = http.Request(method, Uri.parse('$baseUrl$path'));
    req.headers['Accept'] = 'application/json';
    req.headers['X-Request-Id'] = uuidV4(_rand);
    if (auth && _access != null) req.headers['Authorization'] = 'Bearer $_access';
    if (characterId != null) req.headers['X-Character-Id'] = characterId!;
    if (idempotencyKey != null) req.headers['Idempotency-Key'] = idempotencyKey;
    if (body != null) {
      req.headers['Content-Type'] = 'application/json';
      req.body = jsonEncode(body);
    }
    try {
      final streamed = await _http.send(req).timeout(timeout);
      return await http.Response.fromStream(streamed).timeout(timeout);
    } on TimeoutException {
      throw const ApiError.network();
    } on http.ClientException {
      throw const ApiError.network();
    }
  }

  Map<String, dynamic> _decode(http.Response res) {
    Object? json;
    try {
      json = res.body.isEmpty ? const <String, dynamic>{} : jsonDecode(res.body);
    } on FormatException {
      json = null;
    }
    if (res.statusCode >= 200 && res.statusCode < 300) {
      if (json is Map<String, dynamic>) return json;
      throw ApiError(res.statusCode, 'E_BAD_RESPONSE');
    }
    final err = json is Map<String, dynamic> ? json['error'] : null;
    if (err is Map<String, dynamic>) {
      throw ApiError(
        res.statusCode,
        err['code'] as String? ?? 'E_UNKNOWN',
        err['message'] as String? ?? '',
        err['field'] as String?,
      );
    }
    throw ApiError(res.statusCode, 'E_HTTP_${res.statusCode}');
  }

  static bool _isCredentialError(http.Response res) {
    try {
      final j = jsonDecode(res.body);
      return j is Map && j['error'] is Map && (j['error'] as Map)['code'] == 'E_INVALID_CREDENTIALS';
    } on FormatException {
      return false;
    }
  }

  static bool _retryable(int status) => status == 502 || status == 503 || status == 504;

  Future<void> _backoff(int attempt) {
    final base = retryDelay.inMilliseconds * (1 << (attempt - 1));
    final jitter = base == 0 ? 0 : _rand.nextInt(base ~/ 2 + 1);
    return Future.delayed(Duration(milliseconds: base + jitter));
  }
}

/// RFC 4122 v4 UUID from a CSPRNG (idempotency + request ids).
String uuidV4(Random rand) {
  final b = List<int>.generate(16, (_) => rand.nextInt(256));
  b[6] = (b[6] & 0x0f) | 0x40;
  b[8] = (b[8] & 0x3f) | 0x80;
  final h = b.map((x) => x.toRadixString(16).padLeft(2, '0')).join();
  return '${h.substring(0, 8)}-${h.substring(8, 12)}-${h.substring(12, 16)}-'
      '${h.substring(16, 20)}-${h.substring(20)}';
}
