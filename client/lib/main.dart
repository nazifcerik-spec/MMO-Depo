import 'package:flutter/material.dart';

import 'app/app.dart';
import 'core/api_client.dart';
import 'core/config.dart';
import 'core/error_reporter.dart';
import 'core/token_store.dart';
import 'features/inventory/item_images.dart';
import 'features/auth/auth_controller.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  ErrorReporter(baseUrl: AppConfig.apiBaseUrl).install();
  await ItemImages.load();
  final api = ApiClient(baseUrl: AppConfig.apiBaseUrl, store: SecureTokenStore());
  runApp(OtterApp(auth: AuthController(api)));
}
