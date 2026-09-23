import 'dart:async';

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/live_events.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/contract_progress.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import 'group_repository.dart';

String guildRankName(AppLocalizations l, int rank) => switch (rank) {
  GuildRanks.leader => l.rankLeader,
  GuildRanks.officer => l.rankOfficer,
  GuildRanks.veteran => l.rankVeteran,
  GuildRanks.member => l.rankMember,
  _ => l.rankInitiate,
};

/// Reloads a panel whenever a push event of one of [prefixes] arrives.
mixin _LiveReload<T extends StatefulWidget> on PanelState<T> {
  StreamSubscription<LiveEvent>? _sub;

  void listen(LiveEvents live, String prefix) {
    _sub = live.events.where((e) => e.kind.startsWith(prefix)).listen((_) => load());
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }
}

/// Name field + button used by both panels to invite by character name.
class _InviteField extends StatefulWidget {
  const _InviteField({required this.onInvite});

  final Future<void> Function(String name) onInvite;

  @override
  State<_InviteField> createState() => _InviteFieldState();
}

class _InviteFieldState extends State<_InviteField> {
  final _name = TextEditingController();

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final n = _name.text.trim();
    if (n.isEmpty) return;
    await widget.onInvite(n);
    _name.clear();
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Row(
      children: [
        Expanded(
          child: TextField(
            controller: _name,
            textInputAction: TextInputAction.done,
            onSubmitted: (_) => _submit(),
            decoration: InputDecoration(hintText: l.characterName),
          ),
        ),
        const SizedBox(width: Space.sm),
        FilledButton.tonal(style: compactButton, onPressed: _submit, child: Text(l.invite)),
      ],
    );
  }
}

Widget _inviteTile(AppLocalizations l, GroupInvite i, VoidCallback accept, VoidCallback decline) => ListTile(
  contentPadding: EdgeInsets.zero,
  title: Text(i.guildName ?? i.fromName),
  subtitle: i.guildName == null ? null : Text(l.invitedBy(i.fromName)),
  trailing: Row(
    mainAxisSize: MainAxisSize.min,
    children: [
      IconButton(
        tooltip: l.accept,
        icon: const Icon(Icons.check, color: Palette.success),
        onPressed: accept,
      ),
      IconButton(
        tooltip: l.decline,
        icon: const Icon(Icons.close, color: Palette.danger),
        onPressed: decline,
      ),
    ],
  ),
);

// --- Party --------------------------------------------------------------------

class PartyPanel extends StatefulWidget {
  const PartyPanel({super.key, required this.repo, required this.live});

  final GroupRepository repo;
  final LiveEvents live;

  @override
  State<PartyPanel> createState() => _PartyPanelState();
}

enum _PartyAction { promote, kick }

class _PartyPanelState extends State<PartyPanel> with PanelState, _LiveReload, AutomaticKeepAliveClientMixin {
  Party? _party;
  List<GroupInvite> _invites = const [];

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    listen(widget.live, 'party.');
    load();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object?>([widget.repo.party(), widget.repo.partyInvites()]);
    _party = r[0] as Party?;
    _invites = r[1]! as List<GroupInvite>;
  }

  Future<void> _do(Future<void> Function() action) => run(() async {
    await action();
    await fetch();
  });

  Future<void> _leave({required bool disband}) async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: disband ? l.disbandPartyTitle : l.leavePartyTitle,
      body: disband ? l.disbandPartyBody : l.leavePartyBody,
      action: disband ? l.disband : l.leave,
      danger: true,
    );
    if (ok) await _do(disband ? widget.repo.disbandParty : widget.repo.leaveParty);
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final me = widget.repo.characterId;
    return frame(() {
      final p = _party;
      if (p == null) {
        return [
          Text(l.noParty, style: text.bodySmall),
          const SizedBox(height: Space.md),
          FilledButton.tonal(onPressed: () => _do(widget.repo.createParty), child: Text(l.createParty)),
          if (_invites.isNotEmpty) ...[
            const SizedBox(height: Space.lg),
            Text(l.pendingInvites, style: text.titleMedium),
            for (final i in _invites)
              _inviteTile(
                l,
                i,
                () => _do(() => widget.repo.answerPartyInvite(i.id, accept: true)),
                () => _do(() => widget.repo.answerPartyInvite(i.id, accept: false)),
              ),
          ],
        ];
      }
      final leader = p.leaderId == me;
      return [
        Text(l.partyMembers(p.members.length, p.sizeCap), style: text.titleMedium),
        const SizedBox(height: Space.sm),
        for (final m in p.members)
          ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(
              m.id == p.leaderId ? Icons.workspace_premium : Icons.person_outline,
              color: m.id == p.leaderId ? Palette.accent : Palette.textSecondary,
            ),
            title: Text(m.name),
            subtitle: Text('${l.levelShort(m.level)} · ${className(l, m.classCode)}'),
            trailing: leader && m.id != me
                ? PopupMenuButton<_PartyAction>(
                    tooltip: l.more,
                    onSelected: (a) => _do(
                      () => a == _PartyAction.kick ? widget.repo.kickFromParty(m.id) : widget.repo.promoteInParty(m.id),
                    ),
                    itemBuilder: (_) => [
                      PopupMenuItem(value: _PartyAction.promote, child: Text(l.makeLeader)),
                      PopupMenuItem(value: _PartyAction.kick, child: Text(l.kick)),
                    ],
                  )
                : null,
          ),
        if (leader && p.members.length < p.sizeCap) ...[
          const SizedBox(height: Space.md),
          _InviteField(onInvite: (n) => _do(() => widget.repo.inviteToParty(n))),
        ],
        const SizedBox(height: Space.xl),
        OutlinedButton(
          style: OutlinedButton.styleFrom(foregroundColor: Palette.danger),
          onPressed: () => _leave(disband: false),
          child: Text(l.leave),
        ),
        if (leader)
          TextButton(
            style: TextButton.styleFrom(foregroundColor: Palette.danger),
            onPressed: () => _leave(disband: true),
            child: Text(l.disband),
          ),
      ];
    });
  }
}

// --- Guild --------------------------------------------------------------------

class GuildPanel extends StatefulWidget {
  const GuildPanel({super.key, required this.repo, required this.live});

  final GroupRepository repo;
  final LiveEvents live;

  @override
  State<GuildPanel> createState() => _GuildPanelState();
}

enum _GuildAction { rank, transfer, kick }

class _GuildPanelState extends State<GuildPanel> with PanelState, _LiveReload, AutomaticKeepAliveClientMixin {
  Guild? _guild;
  List<GroupInvite> _invites = const [];
  final _name = TextEditingController();
  final _tag = TextEditingController();

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    listen(widget.live, 'guild.');
    load();
  }

  @override
  void dispose() {
    _name.dispose();
    _tag.dispose();
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    final r = await Future.wait<Object?>([widget.repo.guild(), widget.repo.guildInvites()]);
    _guild = r[0] as Guild?;
    _invites = r[1]! as List<GroupInvite>;
  }

  Future<void> _do(Future<void> Function() action) => run(() async {
    await action();
    await fetch();
  });

  Future<void> _create() async {
    final l = AppLocalizations.of(context);
    final ok = await confirmAction(
      context,
      title: l.createGuild,
      body: l.createGuildBody(goldString(guildCreationFeeCopper)),
      action: l.create,
    );
    if (ok) await _do(() => widget.repo.createGuild(_name.text.trim(), _tag.text.trim().toUpperCase()));
  }

  Future<void> _editMotd(String current) async {
    final motd = await showDialog<String>(
      context: context,
      builder: (_) => _MotdDialog(current: current),
    );
    if (motd != null) await _do(() => widget.repo.setMotd(motd));
  }

  Future<void> _memberAction(Guild g, GroupMember m, _GuildAction a) async {
    final l = AppLocalizations.of(context);
    switch (a) {
      case _GuildAction.rank:
        final mine = g.myRank!;
        final rank = await showDialog<int>(
          context: context,
          builder: (ctx) => SimpleDialog(
            title: Text(l.setRankFor(m.name)),
            children: [
              for (var r = GuildRanks.officer; r <= GuildRanks.initiate; r++)
                if (GuildRanks.canSetRank(mine, m.rank!, r))
                  SimpleDialogOption(onPressed: () => Navigator.pop(ctx, r), child: Text(guildRankName(l, r))),
            ],
          ),
        );
        if (rank != null) await _do(() => widget.repo.setGuildRank(m.id, rank));
      case _GuildAction.transfer:
        final ok = await confirmAction(
          context,
          title: l.transferTitle(m.name),
          body: l.transferBody,
          action: l.makeLeader,
          danger: true,
        );
        if (ok) await _do(() => widget.repo.transferGuild(m.id));
      case _GuildAction.kick:
        final ok = await confirmAction(context, title: l.kickTitle(m.name), body: '', action: l.kick, danger: true);
        if (ok) await _do(() => widget.repo.kickFromGuild(m.id));
    }
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final me = widget.repo.characterId;
    return frame(() {
      final g = _guild;
      if (g == null) {
        return [
          if (_invites.isNotEmpty) ...[
            Text(l.pendingInvites, style: text.titleMedium),
            for (final i in _invites)
              _inviteTile(
                l,
                i,
                () => _do(() => widget.repo.answerGuildInvite(i.id, accept: true)),
                () => _do(() => widget.repo.answerGuildInvite(i.id, accept: false)),
              ),
            const SizedBox(height: Space.lg),
          ],
          Text(l.createGuild, style: text.titleMedium),
          const SizedBox(height: Space.sm),
          TextField(
            controller: _name,
            maxLength: 24,
            decoration: InputDecoration(labelText: l.guildName),
          ),
          TextField(
            controller: _tag,
            maxLength: 5,
            textCapitalization: TextCapitalization.characters,
            decoration: InputDecoration(labelText: l.guildTag, helperText: l.guildTagHint),
          ),
          const SizedBox(height: Space.md),
          FilledButton.tonal(onPressed: _create, child: Text(l.createGuildFee(goldString(guildCreationFeeCopper)))),
        ];
      }
      final mine = g.myRank ?? GuildRanks.initiate;
      return [
        Text('${g.name} [${g.tag}]', style: text.titleLarge),
        Text(
          '${l.levelShort(g.level)} · ${l.guildMembers(g.memberCount, g.memberCap)} · ${guildRankName(l, mine)}',
          style: text.bodySmall,
        ),
        if (g.disbandAt != null) ...[
          const SizedBox(height: Space.sm),
          Text(
            l.disbandScheduled(shortDuration(l, g.disbandAt!.difference(DateTime.now()))),
            style: const TextStyle(color: Palette.danger, fontWeight: FontWeight.w600),
          ),
        ],
        const SizedBox(height: Space.md),
        Card(
          child: ListTile(
            title: Text(l.motd, style: text.bodySmall),
            subtitle: Text(g.motd.isEmpty ? '—' : g.motd, style: text.bodyMedium),
            trailing: GuildRanks.canSetMotd(mine)
                ? IconButton(
                    tooltip: l.editMotd,
                    icon: const Icon(Icons.edit_outlined),
                    onPressed: () => _editMotd(g.motd),
                  )
                : null,
          ),
        ),
        const SizedBox(height: Space.md),
        if (GuildRanks.canInvite(mine) && g.memberCount < g.memberCap) ...[
          _InviteField(onInvite: (n) => _do(() => widget.repo.inviteToGuild(n))),
          const SizedBox(height: Space.md),
        ],
        for (final m in g.members)
          ListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(m.name),
            subtitle: Text(
              '${guildRankName(l, m.rank ?? GuildRanks.initiate)} · ${l.levelShort(m.level)} · ${className(l, m.classCode)}',
            ),
            trailing: () {
              final target = m.rank ?? GuildRanks.initiate;
              final actions = [
                if (m.id != me && GuildRanks.leader < target && _anyRank(mine, target)) _GuildAction.rank,
                if (m.id != me && mine == GuildRanks.leader) _GuildAction.transfer,
                if (m.id != me && GuildRanks.canKick(mine, target)) _GuildAction.kick,
              ];
              if (actions.isEmpty) return null;
              return PopupMenuButton<_GuildAction>(
                tooltip: l.more,
                onSelected: (a) => _memberAction(g, m, a),
                itemBuilder: (_) => [
                  for (final a in actions)
                    PopupMenuItem(
                      value: a,
                      child: Text(switch (a) {
                        _GuildAction.rank => l.setRank,
                        _GuildAction.transfer => l.makeLeader,
                        _GuildAction.kick => l.kick,
                      }),
                    ),
                ],
              );
            }(),
          ),
        const SizedBox(height: Space.xl),
        if (mine != GuildRanks.leader)
          OutlinedButton(
            style: OutlinedButton.styleFrom(foregroundColor: Palette.danger),
            onPressed: () async {
              final ok = await confirmAction(
                context,
                title: l.leaveGuildTitle,
                body: l.leaveGuildBody,
                action: l.leave,
                danger: true,
              );
              if (ok) await _do(widget.repo.leaveGuild);
            },
            child: Text(l.leave),
          )
        else if (g.disbandAt == null)
          TextButton(
            style: TextButton.styleFrom(foregroundColor: Palette.danger),
            onPressed: () async {
              final ok = await confirmAction(
                context,
                title: l.disbandGuildTitle,
                body: l.disbandGuildBody,
                action: l.disband,
                danger: true,
              );
              if (ok) await _do(widget.repo.disbandGuild);
            },
            child: Text(l.disband),
          )
        else
          TextButton(onPressed: () => _do(widget.repo.cancelDisband), child: Text(l.cancelDisband)),
      ];
    });
  }

  static bool _anyRank(int actor, int target) {
    for (var r = GuildRanks.officer; r <= GuildRanks.initiate; r++) {
      if (GuildRanks.canSetRank(actor, target, r)) return true;
    }
    return false;
  }
}

/// Owns its controller so it outlives the dialog's exit animation.
class _MotdDialog extends StatefulWidget {
  const _MotdDialog({required this.current});

  final String current;

  @override
  State<_MotdDialog> createState() => _MotdDialogState();
}

class _MotdDialogState extends State<_MotdDialog> {
  late final _c = TextEditingController(text: widget.current);

  @override
  void dispose() {
    _c.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.motd),
      content: TextField(controller: _c, autofocus: true, maxLength: 256, maxLines: 4, minLines: 2),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(onPressed: () => Navigator.pop(context, _c.text.trim()), child: Text(l.save)),
      ],
    );
  }
}
