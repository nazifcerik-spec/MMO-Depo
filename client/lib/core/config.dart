/// Build-time configuration (`--dart-define`). No secrets live here.
///
///   flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080
class AppConfig {
  const AppConfig._();

  static const flavor = String.fromEnvironment('FLAVOR', defaultValue: 'dev');
  static const apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
}
