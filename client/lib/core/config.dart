/// Build-time configuration (`--dart-define`). No secrets live here.
///
///   flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080
class AppConfig {
  const AppConfig._();

  static const flavor = String.fromEnvironment('FLAVOR', defaultValue: 'dev');
  static const apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
  static const appVersion = String.fromEnvironment('APP_VERSION');

  /// JSON describing the latest Android build (DD-304); empty = no in-app
  /// updates (dev, web, tests).
  static const updateFeedUrl = String.fromEnvironment('UPDATE_FEED_URL');
}
