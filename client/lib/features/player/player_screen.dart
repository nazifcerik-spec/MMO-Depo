import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../codex/codex.dart';
import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../core/browser.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../character/class_catalog.dart';
import '../inventory/inventory_models.dart' show equipSlots;
import '../inventory/item_widgets.dart';
import '../pets/pets_screen.dart' show PetArt, petName, petStage;
import '../shop/cosmetics.dart';
import '../social/social_repository.dart';
import '../social/social_screen.dart' show showReportDialog;
import '../workshop/workshop_view.dart' show professionName;
import 'player_repository.dart';

/// Opens [name]'s public page (DD-371); [me] is the viewing character.
Future<void> openPlayer(BuildContext context, ApiClient api, String me, String name) => Navigator.of(context).push(
  MaterialPageRoute<void>(
    builder: (_) => PlayerScreen(repo: PlayerRepository(api, me), name: name),
  ),
);

/// Asks for a name and opens that player's page.
Future<void> findPlayer(BuildContext context, ApiClient api, String me) async {
  final name = await showDialog<String>(context: context, builder: (_) => const _FindDialog());
  if (name != null && name.isNotEmpty && context.mounted) await openPlayer(context, api, me, name);
}

class _FindDialog extends StatefulWidget {
  const _FindDialog();

  @override
  State<_FindDialog> createState() => _FindDialogState();
}

class _FindDialogState extends State<_FindDialog> {
  final _name = TextEditingController();

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  void _submit() => Navigator.pop(context, _name.text.trim());

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return AlertDialog(
      title: Text(l.profileFind),
      content: TextField(
        key: const Key('find-player'),
        controller: _name,
        autofocus: true,
        maxLength: 16,
        textInputAction: TextInputAction.search,
        onSubmitted: (_) => _submit(),
        decoration: InputDecoration(hintText: l.characterName, counterText: ''),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: Text(l.cancel)),
        FilledButton(onPressed: _submit, child: Text(l.jobGo)),
      ],
    );
  }
}

/// A player's public page: who they are, what they wear and what they have
/// achieved. Never anything of the account.
class PlayerScreen extends StatefulWidget {
  const PlayerScreen({super.key, required this.repo, required this.name});

  final PlayerRepository repo;
  final String name;

  @override
  State<PlayerScreen> createState() => _PlayerScreenState();
}

class _PlayerScreenState extends State<PlayerScreen> {
  PlayerProfile? _p;
  ApiError? _error;
  late final _social = SocialRepository(widget.repo.api, widget.repo.characterId);

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final p = await widget.repo.profile(widget.name);
      if (mounted) {
        setState(() {
          _p = p;
          _error = null;
        });
      }
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  Future<void> _addFriend(PlayerProfile p) async {
    final l = AppLocalizations.of(context);
    try {
      final instant = await _social.addFriend(p.name);
      if (mounted) toast(context, Text(instant ? l.friendAdded(p.name) : l.friendRequested(p.name)));
    } on ApiError catch (e) {
      if (mounted) toast(context, Text(errorText(l, e)));
    }
  }

  Future<void> _copyLink(String url) async {
    await Clipboard.setData(ClipboardData(text: url));
    if (mounted) toast(context, Text(AppLocalizations.of(context).profileLinkCopied));
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final p = _p;
    final link = p == null ? null : profileUrl(p.name);
    final mine = p?.id == widget.repo.characterId;
    return Scaffold(
      appBar: AppBar(
        title: Text(l.profileTitle),
        actions: [
          if (link != null)
            IconButton(tooltip: l.profileCopyLink, icon: const Icon(Icons.link), onPressed: () => _copyLink(link)),
          if (p != null && !mine) ...[
            IconButton(tooltip: l.addFriend, icon: const Icon(Icons.person_add_alt), onPressed: () => _addFriend(p)),
            IconButton(
              tooltip: l.report,
              icon: const Icon(Icons.flag_outlined),
              onPressed: () => showReportDialog(context, _social, p.name),
            ),
          ],
        ],
      ),
      body: SafeArea(
        child: PageFrame(
          maxWidth: gameWidth,
          child: p == null
              ? (_error == null
                    ? const Center(child: CircularProgressIndicator())
                    : MessageView(
                        message: _error!.status == 404 ? l.profileNotFound : errorText(l, _error!),
                        icon: Icons.person_search_outlined,
                        onRetry: _error!.status == 404 ? null : _load,
                      ))
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView(
                    padding: const EdgeInsets.fromLTRB(Space.screen, Space.sm, Space.screen, Space.xxl),
                    children: [
                      PlayerBanner(profile: p),
                      const SizedBox(height: Space.lg),
                      _Deeds(profile: p),
                      const SizedBox(height: Space.lg),
                      Panes(
                        main: [_Worn(profile: p)],
                        side: [
                          if (!isWide(context)) const SizedBox(height: Space.lg),
                          _Work(profile: p),
                          if (p.pet case final pet?) ...[const SizedBox(height: Space.lg), _PetCard(pet: pet)],
                        ],
                      ),
                    ],
                  ),
                ),
        ),
      ),
    );
  }
}

/// The page's head: the player's background, portrait and names.
class PlayerBanner extends StatelessWidget {
  const PlayerBanner({super.key, required this.profile});

  final PlayerProfile profile;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final p = profile;
    final accent = backgroundColors(p.background).$1;
    final wide = isWide(context);
    return Container(
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(borderRadius: BorderRadius.circular(Radii.lg)),
      foregroundDecoration: BoxDecoration(
        borderRadius: BorderRadius.circular(Radii.lg),
        border: Border.all(color: accent.withValues(alpha: 0.6)),
      ),
      child: SizedBox(
        height: wide ? 220 : 168,
        child: BackgroundArt(
          code: p.background,
          child: Padding(
            padding: const EdgeInsets.all(Space.lg),
            child: Align(
              alignment: Alignment.bottomLeft,
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(3),
                    decoration: BoxDecoration(shape: BoxShape.circle, color: accent),
                    child: AvatarBadge(code: p.avatar, size: wide ? 112 : 84),
                  ),
                  const SizedBox(width: Space.lg),
                  Expanded(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (p.title.isNotEmpty)
                          InkWell(
                            onTap: () => showCodex(context, CodexRef('title', p.title)),
                            child: Text(
                              ContentNames.of('title', p.title),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: text.bodyMedium?.copyWith(color: accent, fontWeight: FontWeight.w600),
                            ),
                          ),
                        Text(p.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: text.headlineMedium),
                        Row(
                          children: [
                            ClassEmblem(root: p.classCode, size: 22),
                            const SizedBox(width: Space.xs),
                            Expanded(
                              child: Text(
                                l.levelLine(
                                  p.level,
                                  ContentNames.of('class', p.pathCode, className(l, p.classCode)),
                                  raceName(l, p.raceCode),
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                          ],
                        ),
                        if (p.guildName.isNotEmpty)
                          Text(
                            '<${p.guildTag}> ${p.guildName}',
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: text.bodyMedium?.copyWith(color: Palette.textSecondary),
                          ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

/// The ranking's numbers, equal boxes.
class _Deeds extends StatelessWidget {
  const _Deeds({required this.profile});

  final PlayerProfile profile;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final p = profile;
    final boxes = [
      (Icons.leaderboard_outlined, l.rankingScore, '${p.score}', p.rank > 0 ? l.profileRank(p.rank) : ''),
      (Icons.emoji_events_outlined, l.rankingAchievements, '${p.achievementPoints}', ''),
      (Icons.auto_awesome_outlined, l.rankingMastery, '${p.masteryTiers}', ''),
      (Icons.castle_outlined, l.rankingDungeons, '${p.dungeonsCleared}', ''),
      (Icons.groups_outlined, l.rankingRaids, '${p.raidsCleared}', ''),
    ];
    return LayoutBuilder(
      builder: (context, box) {
        final cols = box.maxWidth >= 700 ? boxes.length : (box.maxWidth >= 420 ? 3 : 2);
        final w = (box.maxWidth - (cols - 1) * Space.sm) / cols;
        return Wrap(
          spacing: Space.sm,
          runSpacing: Space.sm,
          children: [
            for (final (icon, label, value, note) in boxes)
              Container(
                width: w,
                height: 92,
                padding: const EdgeInsets.all(Space.md),
                decoration: BoxDecoration(
                  color: Palette.bgSurface,
                  borderRadius: BorderRadius.circular(Radii.md),
                  border: Border.all(color: Palette.borderSubtle),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        Icon(icon, size: 16, color: Palette.textSecondary),
                        const SizedBox(width: Space.xs),
                        Expanded(
                          child: Text(
                            label,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: Theme.of(context).textTheme.bodySmall,
                          ),
                        ),
                      ],
                    ),
                    Text(
                      value,
                      style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w700, fontFeatures: numeric),
                    ),
                    Text(
                      note,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(fontSize: 12, color: Palette.textMuted),
                    ),
                  ],
                ),
              ),
          ],
        );
      },
    );
  }
}

class _Worn extends StatelessWidget {
  const _Worn({required this.profile});

  final PlayerProfile profile;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final gear = profile.equipment;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l.equipment, style: Theme.of(context).textTheme.titleMedium),
        const SizedBox(height: Space.md),
        GridView.count(
          crossAxisCount: isWide(context) ? 6 : 4,
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          mainAxisSpacing: Space.sm,
          crossAxisSpacing: Space.sm,
          children: [
            for (final slot in equipSlots)
              if (gear[slot] case final g?)
                Tooltip(
                  message: itemName(g.templateCode, g.prefix, g.suffix),
                  child: ItemCell(
                    onTap: () => showCodex(context, CodexRef('item', g.templateCode)),
                    rarity: g.rarity,
                    icon: itemIcon(g.itemClass, g.itemSubclass),
                    code: g.templateCode,
                    semantic: '${slotName(l, slot)}: ${itemName(g.templateCode, g.prefix, g.suffix)}',
                  ),
                )
              else
                ItemCell(label: slotName(l, slot)),
          ],
        ),
      ],
    );
  }
}

class _Work extends StatelessWidget {
  const _Work({required this.profile});

  final PlayerProfile profile;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final list = profile.professions;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l.rankingProfessions, style: text.titleMedium),
        const SizedBox(height: Space.md),
        if (list.isEmpty) Text(l.profileNoWork, style: text.bodySmall),
        for (final (code, level) in list)
          InkWell(
            onTap: () => showCodex(context, CodexRef('profession', code)),
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: Space.xs),
              child: Row(
                children: [
                  CircleArt(asset: 'assets/professions/$code.webp', size: 32, fallback: const SizedBox()),
                  const SizedBox(width: Space.sm),
                  Expanded(child: Text(professionName(l, code), overflow: TextOverflow.ellipsis)),
                  Text(
                    l.levelShort(level),
                    style: const TextStyle(fontWeight: FontWeight.w700, fontFeatures: numeric),
                  ),
                ],
              ),
            ),
          ),
      ],
    );
  }
}

class _PetCard extends StatelessWidget {
  const _PetCard({required this.pet});

  final (String, int) pet;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final (species, level) = pet;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l.profilePet, style: Theme.of(context).textTheme.titleMedium),
        const SizedBox(height: Space.md),
        Row(
          children: [
            PetArt(species: species, stage: petStage(level), size: 64, color: Palette.accent),
            const SizedBox(width: Space.md),
            Expanded(child: Text(petName(species))),
            CodexButton(CodexRef('pet', species)),
            Text(
              l.levelShort(level),
              style: const TextStyle(fontWeight: FontWeight.w700, fontFeatures: numeric),
            ),
          ],
        ),
      ],
    );
  }
}
