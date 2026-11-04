import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import '../core/app_locale.dart';
import '../features/auth/auth_controller.dart';
import '../features/auth/auth_screen.dart';
import '../features/character/character_repository.dart';
import '../features/character/character_select_screen.dart';
import '../l10n/app_localizations.dart';
import '../shared/widgets.dart';
import 'theme.dart';

class OtterApp extends StatefulWidget {
  const OtterApp({super.key, required this.auth, this.locale});

  final AuthController auth;

  /// Forces a locale (tests); null follows the device.
  final Locale? locale;

  @override
  State<OtterApp> createState() => _OtterAppState();
}

class _OtterAppState extends State<OtterApp> {
  final _navigator = GlobalKey<NavigatorState>();
  late final _repo = CharacterRepository(widget.auth.api);

  @override
  void initState() {
    super.initState();
    widget.auth.addListener(_onAuth);
    widget.auth.restore();
  }

  @override
  void dispose() {
    widget.auth.removeListener(_onAuth);
    super.dispose();
  }

  /// Any sign-out (manual or forced by the server) drops pushed screens so
  /// the gate below is visible again.
  void _onAuth() {
    if (widget.auth.status != AuthStatus.signedIn) {
      _navigator.currentState?.popUntil((r) => r.isFirst);
    }
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder(
      valueListenable: AppLocale.instance,
      builder: (context, chosen, _) => _app(context, widget.locale ?? chosen),
    );
  }

  Widget _app(BuildContext context, Locale? locale) {
    return MaterialApp(
      navigatorKey: _navigator,
      onGenerateTitle: (c) => AppLocalizations.of(c).appTitle,
      debugShowCheckedModeBanner: false,
      theme: buildTheme(),
      locale: locale,
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      home: ListenableBuilder(
        listenable: widget.auth,
        builder: (context, _) => switch (widget.auth.status) {
          AuthStatus.unknown => const Scaffold(body: Center(child: CircularProgressIndicator())),
          AuthStatus.offline => Scaffold(
            body: MessageView(message: AppLocalizations.of(context).offline, onRetry: widget.auth.restore),
          ),
          AuthStatus.signedOut => AuthScreen(key: const ValueKey('auth'), auth: widget.auth),
          AuthStatus.signedIn => CharacterSelectScreen(auth: widget.auth, repo: _repo),
        },
      ),
    );
  }
}
