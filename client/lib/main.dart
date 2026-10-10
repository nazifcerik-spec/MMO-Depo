import 'package:flutter/material.dart';

import 'app/app.dart';
import 'core/api_client.dart';
import 'core/config.dart';
import 'core/error_reporter.dart';
import 'core/token_store.dart';
import 'features/auth/auth_controller.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  ErrorReporter(baseUrl: AppConfig.apiBaseUrl).install();
  final api = ApiClient(baseUrl: AppConfig.apiBaseUrl, store: SecureTokenStore());
  runApp(OtterApp(auth: AuthController(api)));
}
