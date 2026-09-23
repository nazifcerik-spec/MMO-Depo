import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/widgets.dart';
import 'auth_controller.dart';

/// "Forgot password": asks for the e-mail and requests a reset link. The
/// server answers the same way whether or not the account exists.
Future<void> showForgotPasswordDialog(BuildContext context, AuthController auth, {String initialEmail = ''}) {
  return showDialog<void>(
    context: context,
    builder: (_) => _ForgotPasswordDialog(auth: auth, initialEmail: initialEmail),
  );
}

class _ForgotPasswordDialog extends StatefulWidget {
  const _ForgotPasswordDialog({required this.auth, required this.initialEmail});

  final AuthController auth;
  final String initialEmail;

  @override
  State<_ForgotPasswordDialog> createState() => _ForgotPasswordDialogState();
}

class _ForgotPasswordDialogState extends State<_ForgotPasswordDialog> {
  late final _email = TextEditingController(text: widget.initialEmail);
  bool _busy = false;
  bool _sent = false;
  String? _error;

  @override
  void dispose() {
    _email.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    if (!_email.text.contains('@')) {
      setState(() => _error = l.invalidEmail);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await widget.auth.requestPasswordReset(_email.text);
      if (mounted) setState(() => _sent = true);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.forgotPassword),
      content: _sent
          ? Text(l.resetSent)
          : Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Text(l.resetHint, style: Theme.of(context).textTheme.bodySmall),
                const SizedBox(height: Space.md),
                TextField(
                  controller: _email,
                  keyboardType: TextInputType.emailAddress,
                  autofillHints: const [AutofillHints.email],
                  decoration: InputDecoration(labelText: l.email, errorText: _error),
                ),
              ],
            ),
      actions: _sent
          ? [TextButton(onPressed: () => Navigator.pop(context), child: Text(l.close))]
          : [
              TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
              FilledButton(onPressed: _busy ? null : _submit, child: Text(l.sendLink)),
            ],
    );
  }
}

/// Change password (signed in). Other devices are signed out by the server.
Future<void> showChangePasswordDialog(BuildContext context, AuthController auth) {
  return showDialog<void>(
    context: context,
    builder: (_) => _ChangePasswordDialog(auth: auth),
  );
}

class _ChangePasswordDialog extends StatefulWidget {
  const _ChangePasswordDialog({required this.auth});

  final AuthController auth;

  @override
  State<_ChangePasswordDialog> createState() => _ChangePasswordDialogState();
}

class _ChangePasswordDialogState extends State<_ChangePasswordDialog> {
  final _current = TextEditingController();
  final _next = TextEditingController();
  final _confirm = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    for (final c in [_current, _next, _confirm]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    final error = _next.text.length < 8
        ? l.passwordMin
        : _next.text != _confirm.text
        ? l.passwordMismatch
        : null;
    if (error != null) {
      setState(() => _error = error);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await widget.auth.changePassword(_current.text, _next.text);
      if (!mounted) return;
      Navigator.pop(context);
      toast(context, Text(l.passwordChanged));
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    Widget field(TextEditingController c, String label, List<String> hints) => Padding(
      padding: const EdgeInsets.only(bottom: Space.sm),
      child: TextField(
        controller: c,
        obscureText: true,
        autofillHints: hints,
        decoration: InputDecoration(labelText: label),
      ),
    );
    return AlertDialog(
      title: Text(l.changePassword),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          field(_current, l.currentPassword, const [AutofillHints.password]),
          field(_next, l.newPassword, const [AutofillHints.newPassword]),
          field(_confirm, l.passwordConfirm, const [AutofillHints.newPassword]),
          if (_error != null) Text(_error!, style: const TextStyle(color: Palette.danger)),
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(onPressed: _busy ? null : _submit, child: Text(l.save)),
      ],
    );
  }
}
