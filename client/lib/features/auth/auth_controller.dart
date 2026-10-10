import 'package:flutter/foundation.dart';

import '../../core/api_client.dart';
import '../../core/live_events.dart';

enum AuthStatus { unknown, offline, signedOut, signedIn }

/// Session state for the whole app. Screens call these methods and render
/// [status]; any [ApiError] from login/register is rethrown for inline display.
class AuthController extends ChangeNotifier {
  AuthController(this.api, {ChannelFactory? wsConnect}) : live = LiveEvents(api, connect: wsConnect) {
    api.onSessionLost = () => _set(AuthStatus.signedOut, expired: true);
  }

  final ApiClient api;

  /// Session-scoped push channel; the in-game shell starts it, sign-out stops it.
  final LiveEvents live;
  AuthStatus status = AuthStatus.unknown;

  /// True when the last sign-out was forced by the server (show a notice).
  bool sessionExpired = false;

  /// Set after an account deletion request: when the data will be erased.
  DateTime? deletionScheduledFor;

  Future<void> restore() async {
    _set(AuthStatus.unknown);
    try {
      _set(await api.restore() ? AuthStatus.signedIn : AuthStatus.signedOut);
    } on ApiError {
      _set(AuthStatus.offline);
    }
  }

  Future<void> login(String emailOrUsername, String password) async {
    await api.login(emailOrUsername.trim(), password);
    _set(AuthStatus.signedIn);
  }

  Future<void> register(String email, String username, String password, {required bool acceptTerms}) async {
    await api.register(email.trim(), username.trim(), password, acceptTerms: acceptTerms);
    _set(AuthStatus.signedIn);
  }

  bool get emailVerified => api.emailVerified;

  /// Always "succeeds" (202) — the server never reveals whether the e-mail exists.
  Future<void> requestPasswordReset(String email) =>
      api.post('/v1/auth/password/reset-request', {'email': email.trim()}, auth: false);

  Future<void> resendVerification() => api.post('/v1/auth/resend-verification', const {});

  /// Other devices are signed out; this one stays.
  Future<void> changePassword(String current, String next) =>
      api.post('/v1/auth/password/change', {'current_password': current, 'new_password': next});

  /// Account rights (DD-231). Deletion signs every device out; signing in
  /// again within 30 days cancels it.
  Future<void> deleteAccount(String password) async {
    final j = await api.delete('/v1/me', {'current_password': password});
    deletionScheduledFor = DateTime.tryParse('${j['purge_after']}')?.toLocal();
    await api.forgetSession();
    _set(AuthStatus.signedOut);
  }

  /// A one-time link that downloads all of the account's data.
  Future<String> exportLink() async => (await api.post('/v1/me/export-link', const {}))['url'] as String;

  Future<void> logout() async {
    await api.logout();
    _set(AuthStatus.signedOut);
  }

  void _set(AuthStatus s, {bool expired = false}) {
    if (s != AuthStatus.signedIn) live.stop();
    if (status == s && sessionExpired == expired) return;
    status = s;
    sessionExpired = expired;
    notifyListeners();
  }
}
