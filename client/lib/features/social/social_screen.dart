import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'group_panels.dart';
import 'group_repository.dart';
import 'social_repository.dart';

/// Chat, friends and mail for one character (opened from the home app bar).
class SocialScreen extends StatefulWidget {
  const SocialScreen({super.key, required this.repo, required this.live, required this.dataVersion});

  final SocialRepository repo;
  final LiveEvents live;

  /// Bumped after claiming mail (gold/bag changed).
  final ValueNotifier<int> dataVersion;

  @override
  State<SocialScreen> createState() => _SocialScreenState();
}

class _SocialScreenState extends State<SocialScreen> with SingleTickerProviderStateMixin {
  late final _tabs = TabController(length: 5, vsync: this);
  late final _groups = GroupRepository(widget.repo.api, widget.repo.characterId);
  final _chat = GlobalKey<_ChatPanelState>();

  @override
  void dispose() {
    _tabs.dispose();
    super.dispose();
  }

  void _whisper(String name) {
    _tabs.animateTo(0);
    _chat.currentState?.openWhisper(name);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text(l.social),
        bottom: TabBar(
          controller: _tabs,
          isScrollable: true,
          tabAlignment: TabAlignment.start,
          indicatorColor: Palette.accent,
          labelColor: Palette.accent,
          tabs: [
            Tab(text: l.tabChat),
            Tab(text: l.tabFriends),
            Tab(text: l.tabMail),
            Tab(text: l.tabParty),
            Tab(text: l.tabGuild),
          ],
        ),
      ),
      body: SafeArea(
        child: PageFrame(
          child: TabBarView(
            controller: _tabs,
            children: [
              _ChatPanel(key: _chat, repo: widget.repo, live: widget.live),
              _FriendsPanel(repo: widget.repo, onWhisper: _whisper),
              _MailPanel(repo: widget.repo, dataVersion: widget.dataVersion),
              PartyPanel(repo: _groups, live: widget.live),
              GuildPanel(repo: _groups, live: widget.live),
            ],
          ),
        ),
      ),
    );
  }
}

// --- Chat ---------------------------------------------------------------------

class _ChatPanel extends StatefulWidget {
  const _ChatPanel({super.key, required this.repo, required this.live});

  final SocialRepository repo;
  final LiveEvents live;

  @override
  State<_ChatPanel> createState() => _ChatPanelState();
}

class _ChatPanelState extends State<_ChatPanel> with AutomaticKeepAliveClientMixin {
  static const _maxMessages = 200;

  final _input = TextEditingController();
  StreamSubscription<LiveEvent>? _sub;

  String _channel = 'global';
  String? _peer; // whisper partner name
  List<ChatMessage> _messages = const []; // oldest first
  final _seen = <int>{};
  Set<String> _blocked = const {};
  Object? _error;
  bool _sending = false;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _sub = widget.live.events.where((e) => e.kind == 'chat.message').listen((e) => _receive(e.payload));
    _loadBlocked();
    _load();
  }

  @override
  void dispose() {
    _sub?.cancel();
    _input.dispose();
    super.dispose();
  }

  Future<void> _loadBlocked() async {
    try {
      final b = await widget.repo.blocks();
      if (mounted) setState(() => _blocked = {for (final x in b) x.id});
    } on ApiError {
      // Filtering pushes is best-effort; history is already filtered server-side.
    }
  }

  Future<void> _load() async {
    setState(() => _error = null);
    try {
      final list = await widget.repo.history(_channel, peer: _channel == 'whisper' ? _peer : null);
      if (!mounted) return;
      setState(() {
        _messages = list;
        _seen
          ..clear()
          ..addAll(list.map((m) => m.id));
      });
    } on ApiError catch (e) {
      if (mounted) {
        setState(() {
          _messages = const [];
          _error = e;
        });
      }
    }
  }

  bool _belongsHere(ChatMessage m) {
    if (m.channel != _channel) return false;
    if (_channel != 'whisper') return true;
    return m.fromName == _peer || m.fromId == widget.repo.characterId;
  }

  void _append(ChatMessage m) {
    if (!_seen.add(m.id) || !_belongsHere(m) || _blocked.contains(m.fromId)) return;
    setState(() {
      final next = [..._messages, m];
      _messages = next.length > _maxMessages ? next.sublist(next.length - _maxMessages) : next;
    });
  }

  void _receive(Map<String, dynamic> payload) => _append(ChatMessage.fromJson(payload));

  void _select(String channel, {String? peer}) {
    setState(() {
      _channel = channel;
      _peer = peer ?? _peer;
    });
    _load();
  }

  void openWhisper(String name) => _select('whisper', peer: name);

  Future<void> _send() async {
    final text = _input.text.trim();
    if (text.isEmpty || _sending) return;
    setState(() => _sending = true);
    try {
      final m = await widget.repo.send(_channel, text, peer: _channel == 'whisper' ? _peer : null);
      _input.clear();
      _append(m);
    } on ApiError catch (e) {
      if (mounted && !e.isSessionError) toast(context, Text(errorText(AppLocalizations.of(context), e)));
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  String _channelName(AppLocalizations l, String c) => switch (c) {
    'guild' => l.chanGuild,
    'party' => l.chanParty,
    _ => l.chanGlobal,
  };

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    final me = widget.repo.characterId;
    return Column(
      children: [
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, 0),
          child: Row(
            children: [
              for (final c in chatChannels) ...[
                ChoiceChip(label: Text(_channelName(l, c)), selected: _channel == c, onSelected: (_) => _select(c)),
                const SizedBox(width: Space.sm),
              ],
              if (_peer != null)
                ChoiceChip(
                  avatar: const Icon(Icons.alternate_email, size: 16),
                  label: Text(_peer!),
                  selected: _channel == 'whisper',
                  onSelected: (_) => _select('whisper'),
                ),
            ],
          ),
        ),
        Expanded(
          child: _error != null
              ? MessageView(message: errorText(l, _error!), onRetry: _load, icon: Icons.forum_outlined)
              : _messages.isEmpty
              ? MessageView(message: l.chatEmpty, icon: Icons.forum_outlined)
              : ListView.builder(
                  reverse: true, // newest at the bottom, stays pinned as messages arrive
                  padding: const EdgeInsets.symmetric(horizontal: Space.screen, vertical: Space.sm),
                  itemCount: _messages.length,
                  itemBuilder: (_, i) {
                    final m = _messages[_messages.length - 1 - i];
                    final mine = m.fromId == me;
                    return Padding(
                      padding: const EdgeInsets.symmetric(vertical: 3),
                      child: Text.rich(
                        TextSpan(
                          children: [
                            TextSpan(
                              text: '${m.fromName}: ',
                              style: TextStyle(
                                fontWeight: FontWeight.w700,
                                color: mine ? Palette.accent : Palette.info,
                              ),
                            ),
                            TextSpan(text: m.body),
                          ],
                        ),
                      ),
                    );
                  },
                ),
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.sm),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _input,
                  maxLength: 512,
                  maxLengthEnforcement: MaxLengthEnforcement.enforced,
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => _send(),
                  decoration: InputDecoration(hintText: l.chatHint(_channelName(l, _channel)), counterText: ''),
                ),
              ),
              IconButton(
                tooltip: l.send,
                onPressed: _sending ? null : _send,
                icon: const Icon(Icons.send, color: Palette.accent),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

// --- Friends ------------------------------------------------------------------

class _FriendsPanel extends StatefulWidget {
  const _FriendsPanel({required this.repo, required this.onWhisper});

  final SocialRepository repo;
  final void Function(String name) onWhisper;

  @override
  State<_FriendsPanel> createState() => _FriendsPanelState();
}

enum _FriendAction { whisper, remove, block }

class _FriendsPanelState extends State<_FriendsPanel> with PanelState, AutomaticKeepAliveClientMixin {
  final _name = TextEditingController();
  List<Friend> _friends = const [];
  List<FriendRequest> _requests = const [];
  List<Blocked> _blocked = const [];

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object>([widget.repo.friends(), widget.repo.requests(), widget.repo.blocks()]);
    _friends = r[0] as List<Friend>;
    _requests = r[1] as List<FriendRequest>;
    _blocked = r[2] as List<Blocked>;
  }

  Future<void> _add() async {
    final name = _name.text.trim();
    if (name.isEmpty) return;
    final l = AppLocalizations.of(context);
    await run(() async {
      final instant = await widget.repo.addFriend(name);
      _name.clear();
      if (mounted) toast(context, Text(instant ? l.friendAdded(name) : l.friendRequested(name)));
      await fetch();
    });
  }

  Future<void> _act(Friend f, _FriendAction a) async {
    final l = AppLocalizations.of(context);
    switch (a) {
      case _FriendAction.whisper:
        widget.onWhisper(f.name);
      case _FriendAction.remove:
        final ok = await confirmAction(
          context,
          title: l.removeFriendTitle(f.name),
          body: l.removeFriendBody,
          action: l.remove,
          danger: true,
        );
        if (ok) await run(() => widget.repo.removeFriend(f.id).then((_) => fetch()));
      case _FriendAction.block:
        final ok = await confirmAction(
          context,
          title: l.blockTitle(f.name),
          body: l.blockBody,
          action: l.block,
          danger: true,
        );
        if (ok) await run(() => widget.repo.block(f.name).then((_) => fetch()));
    }
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return frame(
      () => [
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _name,
                textInputAction: TextInputAction.done,
                onSubmitted: (_) => _add(),
                decoration: InputDecoration(hintText: l.characterName),
              ),
            ),
            const SizedBox(width: Space.sm),
            FilledButton.tonal(style: compactButton, onPressed: _add, child: Text(l.addFriend)),
          ],
        ),
        if (_requests.isNotEmpty) ...[
          const SizedBox(height: Space.lg),
          Text(l.friendRequests, style: text.titleMedium),
          for (final r in _requests)
            ListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(r.fromName),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    tooltip: l.accept,
                    icon: const Icon(Icons.check, color: Palette.success),
                    onPressed: () => run(() => widget.repo.answer(r.id, accept: true).then((_) => fetch())),
                  ),
                  IconButton(
                    tooltip: l.decline,
                    icon: const Icon(Icons.close, color: Palette.danger),
                    onPressed: () => run(() => widget.repo.answer(r.id, accept: false).then((_) => fetch())),
                  ),
                ],
              ),
            ),
        ],
        const SizedBox(height: Space.lg),
        Text(l.friendsCount(_friends.length), style: text.titleMedium),
        if (_friends.isEmpty)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: Space.md),
            child: Text(l.noFriends, style: text.bodySmall),
          ),
        for (final f in _friends)
          ListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(f.name),
            subtitle: Text('${l.levelShort(f.level)} · ${className(l, f.classCode)}'),
            trailing: PopupMenuButton<_FriendAction>(
              tooltip: l.more,
              onSelected: (a) => _act(f, a),
              itemBuilder: (_) => [
                PopupMenuItem(value: _FriendAction.whisper, child: Text(l.whisper)),
                PopupMenuItem(value: _FriendAction.remove, child: Text(l.remove)),
                PopupMenuItem(value: _FriendAction.block, child: Text(l.block)),
              ],
            ),
          ),
        if (_blocked.isNotEmpty) ...[
          const SizedBox(height: Space.lg),
          Text(l.blockedList, style: text.titleMedium),
          for (final b in _blocked)
            ListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(b.name, style: const TextStyle(color: Palette.textMuted)),
              trailing: TextButton(
                onPressed: () => run(() => widget.repo.unblock(b.id).then((_) => fetch())),
                child: Text(l.unblock),
              ),
            ),
        ],
      ],
    );
  }
}

// --- Mail ---------------------------------------------------------------------

class _MailPanel extends StatefulWidget {
  const _MailPanel({required this.repo, required this.dataVersion});

  final SocialRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<_MailPanel> createState() => _MailPanelState();
}

enum _MailAction { claim, delete }

class _MailPanelState extends State<_MailPanel> with PanelState, AutomaticKeepAliveClientMixin {
  List<MailSummary> _inbox = const [];

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _inbox = await widget.repo.inbox();

  Future<void> _open(MailSummary s) async {
    final l = AppLocalizations.of(context);
    Mail mail;
    try {
      mail = await widget.repo.read(s.id);
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
      return;
    }
    if (!mounted) return;
    unawaited(load()); // read state changed
    final action = await showModalBottomSheet<_MailAction>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (c) => _MailSheet(mail: mail),
    );
    if (action == null || !mounted) return;
    switch (action) {
      case _MailAction.claim:
        await run(() async {
          await widget.repo.claim(s.id);
          if (mounted) toast(context, Text(l.mailClaimed));
          widget.dataVersion.value++;
          await fetch();
        });
      case _MailAction.delete:
        final ok = await confirmAction(
          context,
          title: l.deleteMailTitle,
          body: l.deleteMailBody,
          action: l.destroy,
          danger: true,
        );
        if (ok) await run(() => widget.repo.deleteMail(s.id).then((_) => fetch()));
    }
  }

  Future<void> _compose() async {
    final l = AppLocalizations.of(context);
    final postage = await showModalBottomSheet<int>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (_) => _ComposeSheet(repo: widget.repo),
    );
    if (postage != null && mounted) {
      toast(context, Text(l.mailSent(goldString(postage))));
      widget.dataVersion.value++; // postage (and attached gold) left the wallet
    }
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    return frame(
      () => [
        FilledButton.tonal(onPressed: _compose, child: Text(l.newMail)),
        const SizedBox(height: Space.md),
        if (_inbox.isEmpty)
          Padding(
            padding: const EdgeInsets.only(top: Space.xl),
            child: Text(l.inboxEmpty, textAlign: TextAlign.center, style: Theme.of(context).textTheme.bodySmall),
          ),
        for (final m in _inbox)
          Card(
            margin: const EdgeInsets.only(bottom: Space.sm),
            child: ListTile(
              leading: Icon(
                m.unread ? Icons.mark_email_unread_outlined : Icons.drafts_outlined,
                color: m.unread ? Palette.accent : Palette.textMuted,
              ),
              title: Text(
                m.subject.isEmpty ? l.noSubject : m.subject,
                style: TextStyle(fontWeight: m.unread ? FontWeight.w700 : FontWeight.w400),
              ),
              subtitle: Text(m.senderName),
              trailing: m.hasAttachments ? const Icon(Icons.attach_file, color: Palette.accent) : null,
              onTap: () => _open(m),
            ),
          ),
      ],
    );
  }
}

class _MailSheet extends StatelessWidget {
  const _MailSheet({required this.mail});

  final Mail mail;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final s = mail.summary;
    return SafeArea(
      child: PageFrame(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(s.subject.isEmpty ? l.noSubject : s.subject, style: text.titleLarge),
              Text(l.mailFrom(s.senderName), style: text.bodySmall),
              const SizedBox(height: Space.lg),
              Text(mail.body),
              if (s.hasAttachments) ...[
                const SizedBox(height: Space.lg),
                Text(l.attachments, style: text.titleMedium),
                const SizedBox(height: Space.sm),
                if (s.gold > 0) GoldText(s.gold),
                for (final a in mail.attachments)
                  Text(
                    a.quantity > 1 ? '${itemName(a.templateCode)} ×${a.quantity}' : itemName(a.templateCode),
                    style: TextStyle(color: rarityColor(a.rarity)),
                  ),
              ],
              const SizedBox(height: Space.xl),
              if (s.hasAttachments)
                FilledButton(onPressed: () => Navigator.pop(context, _MailAction.claim), child: Text(l.claim))
              else
                OutlinedButton(
                  style: OutlinedButton.styleFrom(foregroundColor: Palette.danger),
                  onPressed: () => Navigator.pop(context, _MailAction.delete),
                  child: Text(l.destroy),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Pops with the postage charged once the mail is sent.
class _ComposeSheet extends StatefulWidget {
  const _ComposeSheet({required this.repo});

  final SocialRepository repo;

  @override
  State<_ComposeSheet> createState() => _ComposeSheetState();
}

class _ComposeSheetState extends State<_ComposeSheet> {
  final _to = TextEditingController();
  final _subject = TextEditingController();
  final _body = TextEditingController();
  final _gold = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    for (final c in [_to, _subject, _body, _gold]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _send() async {
    final l = AppLocalizations.of(context);
    if (_to.text.trim().isEmpty) {
      setState(() => _error = l.fieldRequired);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final postage = await widget.repo.sendMail(
        to: _to.text.trim(),
        subject: _subject.text.trim(),
        body: _body.text.trim(),
        gold: (int.tryParse(_gold.text) ?? 0) * 10000,
      );
      if (mounted) Navigator.of(context).pop(postage);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(l, e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    const gap = SizedBox(height: Space.md);
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: PageFrame(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(l.newMail, style: Theme.of(context).textTheme.titleLarge),
              const SizedBox(height: Space.lg),
              TextField(
                controller: _to,
                decoration: InputDecoration(labelText: l.mailTo),
              ),
              gap,
              TextField(
                controller: _subject,
                maxLength: 64,
                decoration: InputDecoration(labelText: l.mailSubject),
              ),
              TextField(
                controller: _body,
                maxLength: 2000,
                minLines: 3,
                maxLines: 6,
                decoration: InputDecoration(labelText: l.mailBody),
              ),
              TextField(
                controller: _gold,
                keyboardType: TextInputType.number,
                inputFormatters: [FilteringTextInputFormatter.digitsOnly, LengthLimitingTextInputFormatter(6)],
                decoration: InputDecoration(labelText: l.attachGold, suffixText: 'g'),
              ),
              if (_error != null) ...[gap, Text(_error!, style: const TextStyle(color: Palette.danger))],
              const SizedBox(height: Space.xl),
              FilledButton(
                onPressed: _busy ? null : _send,
                child: _busy
                    ? const SizedBox.square(dimension: 22, child: CircularProgressIndicator(strokeWidth: 2.5))
                    : Text(l.send),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
