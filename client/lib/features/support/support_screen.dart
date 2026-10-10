import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'support_repository.dart';

String supportCategoryName(AppLocalizations l, String c) => switch (c) {
  'bug' => l.supportBug,
  'account' => l.supportAccount,
  'abuse' => l.supportAbuse,
  'feedback' => l.supportFeedback,
  _ => l.supportOther,
};

/// My tickets with the team's answers, and a form for a new one (DD-234).
class SupportScreen extends StatefulWidget {
  const SupportScreen({super.key, required this.repo});

  final SupportRepository repo;

  @override
  State<SupportScreen> createState() => _SupportScreenState();
}

class _SupportScreenState extends State<SupportScreen> with PanelState {
  List<SupportTicket> _tickets = const [];

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _tickets = await widget.repo.mine();

  Future<void> _new() async {
    final sent = await showDialog<bool>(
      context: context,
      builder: (_) => _NewTicketDialog(repo: widget.repo),
    );
    if (sent == true && mounted) {
      toast(context, Text(AppLocalizations.of(context).supportSent));
      await load();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text(l.support)),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: busy ? null : _new,
        icon: const Icon(Icons.edit_outlined),
        label: Text(l.supportNew),
      ),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              Text(l.supportIntro, style: text.bodySmall),
              const SizedBox(height: Space.md),
              if (_tickets.isEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Space.xxl),
                  child: MessageView(message: l.supportNone, icon: Icons.support_agent),
                ),
              for (final t in _tickets) _TicketCard(ticket: t),
            ],
          ),
        ),
      ),
    );
  }
}

class _TicketCard extends StatelessWidget {
  const _TicketCard({required this.ticket});

  final SupportTicket ticket;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final (label, color) = switch (ticket.state) {
      'answered' => (l.supportAnswered, Palette.accent),
      'closed' => (l.supportClosed, Palette.textMuted),
      _ => (l.supportOpen, Palette.textSecondary),
    };
    return Card(
      margin: const EdgeInsets.only(bottom: Space.sm),
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(child: Text(ticket.subject, style: text.titleSmall)),
                Text(label, style: text.labelSmall?.copyWith(color: color)),
              ],
            ),
            Text(
              '${supportCategoryName(l, ticket.category)} · ${MaterialLocalizations.of(context).formatMediumDate(ticket.createdAt)}',
              style: text.bodySmall?.copyWith(color: Palette.textMuted),
            ),
            const SizedBox(height: Space.xs),
            Text(ticket.body, style: text.bodyMedium),
            if (ticket.reply != null) ...[
              const SizedBox(height: Space.sm),
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(Space.sm),
                decoration: BoxDecoration(
                  color: Palette.accent.withValues(alpha: 0.08),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Palette.accent.withValues(alpha: 0.4)),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(l.supportReply, style: text.labelSmall?.copyWith(color: Palette.accent)),
                    Text(ticket.reply!, style: text.bodyMedium),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _NewTicketDialog extends StatefulWidget {
  const _NewTicketDialog({required this.repo});

  final SupportRepository repo;

  @override
  State<_NewTicketDialog> createState() => _NewTicketDialogState();
}

class _NewTicketDialogState extends State<_NewTicketDialog> {
  final _subject = TextEditingController();
  final _body = TextEditingController();
  String _category = 'bug';
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _subject.dispose();
    _body.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final l = AppLocalizations.of(context);
    if (_subject.text.trim().length < 3 || _body.text.trim().length < 10) {
      setState(() => _error = l.supportTooShort);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await widget.repo.create(
        category: _category,
        subject: _subject.text.trim(),
        body: _body.text.trim(),
        locale: Localizations.localeOf(context).languageCode,
      );
      if (mounted) Navigator.pop(context, true);
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
      title: Text(l.supportNew),
      content: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            DropdownButtonFormField<String>(
              initialValue: _category,
              isExpanded: true,
              decoration: InputDecoration(labelText: l.supportCategory),
              items: [
                for (final c in supportCategories) DropdownMenuItem(value: c, child: Text(supportCategoryName(l, c))),
              ],
              onChanged: (v) => setState(() => _category = v ?? _category),
            ),
            const SizedBox(height: Space.sm),
            TextField(
              key: const Key('support-subject'),
              controller: _subject,
              maxLength: 120,
              decoration: InputDecoration(labelText: l.supportSubject),
            ),
            TextField(
              key: const Key('support-body'),
              controller: _body,
              maxLength: 4000,
              minLines: 4,
              maxLines: 8,
              decoration: InputDecoration(labelText: l.supportMessage, alignLabelWithHint: true),
            ),
            if (_error != null) Text(_error!, style: const TextStyle(color: Palette.danger)),
          ],
        ),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(onPressed: _busy ? null : _send, child: Text(l.send)),
      ],
    );
  }
}
