import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/links.dart';
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

/// Account menu (DD-231): password, data export, deletion, legal pages.
class AccountMenuButton extends StatelessWidget {
  const AccountMenuButton({super.key, required this.auth});

  final AuthController auth;

  Future<void> _export(BuildContext context) async {
    final l = AppLocalizations.of(context);
    try {
      final url = await auth.exportLink();
      if (context.mounted && await openExternal(context, url) && context.mounted) {
        toast(context, Text(l.exportOpened));
      }
    } on ApiError catch (e) {
      if (context.mounted) toast(context, Text(errorText(l, e)));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return PopupMenuButton<String>(
      tooltip: l.accountMenu,
      icon: const Icon(Icons.manage_accounts_outlined),
      onSelected: (v) => switch (v) {
        'password' => showChangePasswordDialog(context, auth),
        'export' => _export(context),
        'delete' => showDialog<void>(
          context: context,
          builder: (_) => _DeleteAccountDialog(auth: auth),
        ),
        'terms' => openExternal(context, legalUrl(context, 'terms')),
        _ => openExternal(context, legalUrl(context, 'privacy')),
      },
      itemBuilder: (_) => [
        PopupMenuItem(value: 'password', child: Text(l.changePassword)),
        PopupMenuItem(value: 'export', child: Text(l.downloadMyData)),
        PopupMenuItem(value: 'terms', child: Text(l.termsOfService)),
        PopupMenuItem(value: 'privacy', child: Text(l.privacyPolicy)),
        const PopupMenuDivider(),
        PopupMenuItem(
          value: 'delete',
          child: Text(l.deleteAccount, style: const TextStyle(color: Palette.danger)),
        ),
      ],
    );
  }
}

class _DeleteAccountDialog extends StatefulWidget {
  const _DeleteAccountDialog({required this.auth});

  final AuthController auth;

  @override
  State<_DeleteAccountDialog> createState() => _DeleteAccountDialogState();
}

class _DeleteAccountDialogState extends State<_DeleteAccountDialog> {
  final _password = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _password.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final l = AppLocalizations.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await widget.auth.deleteAccount(_password.text);
      // Signing out already dismissed every route; popping again would close
      // the root screen and leave a blank app.
      if (mounted && (ModalRoute.of(context)?.isCurrent ?? false)) Navigator.pop(context);
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
      title: Text(l.deleteAccount),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(l.deleteAccountBody),
          const SizedBox(height: Space.md),
          TextField(
            key: const Key('delete-password'),
            controller: _password,
            obscureText: true,
            autofillHints: const [AutofillHints.password],
            decoration: InputDecoration(labelText: l.currentPassword),
          ),
          if (_error != null) ...[
            const SizedBox(height: Space.sm),
            Text(_error!, style: const TextStyle(color: Palette.danger)),
          ],
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(
          style: FilledButton.styleFrom(backgroundColor: Palette.danger),
          onPressed: _busy ? null : _submit,
          child: Text(l.deleteAccountConfirm),
        ),
      ],
    );
  }
}
