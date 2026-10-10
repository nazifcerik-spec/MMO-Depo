import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/links.dart';
import '../../shared/widgets.dart';
import 'account_dialogs.dart';
import 'auth_controller.dart';

final _emailPattern = RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$');
final _usernamePattern = RegExp(r'^[a-zA-Z0-9_.]{3,32}$');

/// Login and register on one screen (docs/ui/login.md). Password reset and
/// OTP challenge arrive with their server endpoints.
class AuthScreen extends StatefulWidget {
  const AuthScreen({super.key, required this.auth});

  final AuthController auth;

  @override
  State<AuthScreen> createState() => _AuthScreenState();
}

class _AuthScreenState extends State<AuthScreen> {
  final _form = GlobalKey<FormState>();
  final _id = TextEditingController();
  final _email = TextEditingController();
  final _username = TextEditingController();
  final _password = TextEditingController();
  final _confirm = TextEditingController();

  bool _registering = false;
  bool _obscure = true;
  bool _terms = false;
  bool _busy = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    if (widget.auth.sessionExpired) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) setState(() => _error = AppLocalizations.of(context).errSessionExpired);
      });
    }
    final purge = widget.auth.deletionScheduledFor;
    if (purge != null) {
      widget.auth.deletionScheduledFor = null;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        final date = MaterialLocalizations.of(context).formatMediumDate(purge);
        setState(() => _error = AppLocalizations.of(context).deleteScheduled(date));
      });
    }
  }

  @override
  void dispose() {
    for (final c in [_id, _email, _username, _password, _confirm]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    setState(() => _error = null);
    if (!_form.currentState!.validate()) return;
    if (_registering && !_terms) {
      setState(() => _error = l.termsRequired);
      return;
    }
    setState(() => _busy = true);
    try {
      if (_registering) {
        await widget.auth.register(_email.text, _username.text, _password.text, acceptTerms: _terms);
      } else {
        await widget.auth.login(_id.text, _password.text);
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  void _toggleMode() => setState(() {
    _registering = !_registering;
    _error = null;
    _form.currentState?.reset();
  });

  String? _required(String? v, AppLocalizations l) => (v == null || v.trim().isEmpty) ? l.fieldRequired : null;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    const gap = SizedBox(height: Space.md);

    return Scaffold(
      body: SafeArea(
        child: PageFrame(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(Space.screen),
            child: AutofillGroup(
              child: Form(
                key: _form,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const SizedBox(height: Space.xxl),
                    const _Logo(),
                    const SizedBox(height: Space.lg),
                    Text(l.appTitle, textAlign: TextAlign.center, style: text.headlineMedium),
                    const SizedBox(height: Space.xs),
                    Text(l.tagline, textAlign: TextAlign.center, style: text.bodySmall),
                    const SizedBox(height: Space.xxl),
                    if (_registering) ...[
                      TextFormField(
                        controller: _email,
                        keyboardType: TextInputType.emailAddress,
                        autofillHints: const [AutofillHints.email],
                        textInputAction: TextInputAction.next,
                        decoration: InputDecoration(labelText: l.email, prefixIcon: const Icon(Icons.mail_outline)),
                        validator: (v) =>
                            _required(v, l) ?? (_emailPattern.hasMatch(v!.trim()) ? null : l.invalidEmail),
                      ),
                      gap,
                      TextFormField(
                        controller: _username,
                        autofillHints: const [AutofillHints.newUsername],
                        textInputAction: TextInputAction.next,
                        decoration: InputDecoration(
                          labelText: l.username,
                          prefixIcon: const Icon(Icons.person_outline),
                        ),
                        validator: (v) =>
                            _required(v, l) ?? (_usernamePattern.hasMatch(v!.trim()) ? null : l.usernameRule),
                      ),
                    ] else
                      TextFormField(
                        controller: _id,
                        autofillHints: const [AutofillHints.username, AutofillHints.email],
                        textInputAction: TextInputAction.next,
                        decoration: InputDecoration(
                          labelText: l.emailOrUsername,
                          prefixIcon: const Icon(Icons.person_outline),
                        ),
                        validator: (v) => _required(v, l),
                      ),
                    gap,
                    TextFormField(
                      controller: _password,
                      obscureText: _obscure,
                      autofillHints: [_registering ? AutofillHints.newPassword : AutofillHints.password],
                      textInputAction: _registering ? TextInputAction.next : TextInputAction.done,
                      onFieldSubmitted: _registering ? null : (_) => _submit(),
                      decoration: InputDecoration(
                        labelText: l.password,
                        prefixIcon: const Icon(Icons.lock_outline),
                        suffixIcon: IconButton(
                          tooltip: _obscure ? l.showPassword : l.hidePassword,
                          icon: Icon(_obscure ? Icons.visibility_outlined : Icons.visibility_off_outlined),
                          onPressed: () => setState(() => _obscure = !_obscure),
                        ),
                      ),
                      validator: (v) => _required(v, l) ?? (_registering && v!.length < 8 ? l.passwordMin : null),
                    ),
                    if (_registering) ...[
                      gap,
                      TextFormField(
                        controller: _confirm,
                        obscureText: _obscure,
                        textInputAction: TextInputAction.done,
                        decoration: InputDecoration(
                          labelText: l.passwordConfirm,
                          prefixIcon: const Icon(Icons.lock_outline),
                        ),
                        validator: (v) => v == _password.text ? null : l.passwordMismatch,
                      ),
                      const SizedBox(height: Space.sm),
                      CheckboxListTile(
                        value: _terms,
                        onChanged: (v) => setState(() => _terms = v ?? false),
                        title: Text(l.acceptTerms, style: text.bodySmall),
                        controlAffinity: ListTileControlAffinity.leading,
                        contentPadding: EdgeInsets.zero,
                        dense: true,
                      ),
                      Wrap(
                        alignment: WrapAlignment.center,
                        children: [
                          TextButton(
                            onPressed: () => openExternal(context, legalUrl(context, 'terms')),
                            child: Text(l.termsOfService),
                          ),
                          TextButton(
                            onPressed: () => openExternal(context, legalUrl(context, 'privacy')),
                            child: Text(l.privacyPolicy),
                          ),
                        ],
                      ),
                    ],
                    if (_error != null) ...[
                      const SizedBox(height: Space.md),
                      Text(
                        _error!,
                        key: const Key('auth-error'),
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Palette.danger),
                      ),
                    ],
                    const SizedBox(height: Space.xl),
                    FilledButton(
                      onPressed: _busy ? null : _submit,
                      child: _busy
                          ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
                          : Text(_registering ? l.register : l.login),
                    ),
                    const SizedBox(height: Space.md),
                    TextButton(
                      onPressed: _busy ? null : _toggleMode,
                      child: Text(_registering ? l.haveAccount : l.noAccount),
                    ),
                    if (!_registering)
                      TextButton(
                        onPressed: _busy
                            ? null
                            : () => showForgotPasswordDialog(
                                context,
                                widget.auth,
                                initialEmail: _id.text.contains('@') ? _id.text.trim() : '',
                              ),
                        child: Text(l.forgotPassword),
                      ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _Logo extends StatelessWidget {
  const _Logo();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Container(
        width: 84,
        height: 84,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          gradient: const RadialGradient(colors: [Color(0x55E4A756), Color(0x000F1013)]),
          border: Border.all(color: Palette.accent.withValues(alpha: 0.6), width: 1.5),
          boxShadow: [BoxShadow(color: Palette.accent.withValues(alpha: 0.25), blurRadius: 24)],
        ),
        child: const Icon(Icons.water_drop_outlined, size: 40, color: Palette.accent),
      ),
    );
  }
}
