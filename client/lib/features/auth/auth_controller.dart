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

  Future<void> register(String email, String username, String password) async {
    await api.register(email.trim(), username.trim(), password);
    _set(AuthStatus.signedIn);
  }

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
