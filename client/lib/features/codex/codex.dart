import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../activity/activity_view.dart' show activityIcon;
import '../character/class_catalog.dart' show classIcon;
import '../helpers/helpers_screen.dart' show professionLabel;
import '../inventory/inventory_models.dart' show StatBonus;
import '../inventory/item_images.dart';
import '../inventory/item_widgets.dart';
import '../quests/quests_screen.dart' show questName;

/// Linked tags (DD-372): every piece of content names everything connected to
/// it as [LinkTag]s, and every tag opens the page of what it names
/// (`GET /v1/codex/{kind}/{code}`).

/// What a tag points at. [sub] narrows the kind (raid, item class, base
/// class); [skill] is the profession level it needs, [count] a quantity.
class CodexRef {
  const CodexRef(this.kind, this.code, {this.sub = '', this.level = 0, this.skill = 0, this.count = 0});

  factory CodexRef.fromJson(Map<String, dynamic> j) => CodexRef(
    j['kind'] as String,
    j['code'] as String,
    sub: j['sub'] as String? ?? '',
    level: (j['level'] as num?)?.toInt() ?? 0,
    skill: (j['skill'] as num?)?.toInt() ?? 0,
    count: (j['count'] as num?)?.toInt() ?? 0,
  );

  final String kind;
  final String code;
  final String sub;
  final int level;
  final int skill;
  final int count;
}

typedef CodexGoal = ({String event, int count});
typedef CodexGroup = ({String key, List<CodexRef> refs});

class CodexPage {
  const CodexPage({
    required this.ref,
    this.facts = const {},
    this.traits = const {},
    this.goals = const [],
    this.bonuses = const [],
    this.groups = const [],
  });

  factory CodexPage.fromJson(Map<String, dynamic> j) {
    List<Map<String, dynamic>> rows(String k) => (j[k] as List? ?? const []).cast<Map<String, dynamic>>();
    final groups = [
      for (final g in rows('groups'))
        (
          key: g['key'] as String,
          refs: [for (final r in (g['refs'] as List).cast<Map<String, dynamic>>()) CodexRef.fromJson(r)],
        ),
    ];
    int order(String key) {
      final i = _groupOrder.indexOf(key);
      return i < 0 ? _groupOrder.length : i;
    }

    return CodexPage(
      ref: CodexRef.fromJson(j),
      facts: {
        for (final e in (j['facts'] as Map<String, dynamic>? ?? const {}).entries) e.key: (e.value as num).toInt(),
      },
      traits: {for (final e in (j['traits'] as Map<String, dynamic>? ?? const {}).entries) e.key: '${e.value}'},
      goals: [for (final g in rows('goals')) (event: g['event'] as String, count: (g['count'] as num).toInt())],
      bonuses: [
        for (final m in rows('modifiers'))
          (
            StatBonus(m['kind'] as String, m['target'] as String, (m['amount'] as num).toDouble()),
            (m['at'] as num?)?.toInt() ?? 0,
          ),
      ],
      groups: groups..sort((a, b) => order(a.key) - order(b.key)),
    );
  }

  final CodexRef ref;
  final Map<String, int> facts;
  final Map<String, String> traits;
  final List<CodexGoal> goals;

  /// Each bonus with the piece count, tier or level it needs (0 = none).
  final List<(StatBonus, int)> bonuses;
  final List<CodexGroup> groups;
}

/// The order relations are listed in: what it is, where it is, how to get
/// it, what it leads to.
const _groupOrder = [
  'classes',
  'equipment',
  'family',
  'set',
  'pieces',
  'egg',
  'hatches',
  'promotes_from',
  'branches',
  'races',
  'skills', //
  'requires', 'leads_to', 'boss', 'enemies', 'hunted_in', 'fought_in', 'boss_of', 'region', 'grows', 'within',
  'contains', 'roads', 'mobs', 'dungeons', 'raids', 'dropped_by', 'dropped_in', 'crafted_by', 'materials', 'gathered',
  'sold_by', 'drops', 'sets', 'sells', 'nodes', 'makes', 'used_by', 'used_in', 'professions', 'sources', 'achievements',
  'titles', 'quests', 'mastery', 'reputation', 'events', 'items',
];

/// Loads pages; content does not change while the app runs, so each page is
/// asked for once. The game sets the client when a character enters.
class Codex {
  Codex._();

  static ApiClient? _api;
  static final _pages = <String, Future<CodexPage?>>{};

  static void use(ApiClient api) {
    if (identical(api, _api)) return;
    _api = api;
    _pages.clear();
  }

  static Future<CodexPage?> page(String kind, String code) => _pages['$kind/$code'] ??= () async {
    try {
      return CodexPage.fromJson(await _api!.get('/v1/codex/$kind/${Uri.encodeComponent(code)}'));
    } catch (_) {
      _pages.remove('$kind/$code'); // try again the next time it is opened
      return null;
    }
  }();
}

/// A content piece's display name.
String codexName(AppLocalizations l, String kind, String code) => switch (kind) {
  'item' => itemName(code),
  'quest' => questName(code),
  'profession' => professionLabel(l, code),
  'race' => raceName(l, code),
  'class' => className(l, code) != code ? className(l, code) : ContentNames.of('class', code),
  'list' => l.codexList(code),
  _ => ContentNames.of(kind, code),
};

/// "Raid", "Monster", ... for a page's subtitle.
String codexKindName(AppLocalizations l, CodexRef r) => l.codexKind(switch ((r.kind, r.sub)) {
  ('dungeon', 'raid') => 'raid',
  ('mob', 'boss') => 'boss',
  _ => r.kind,
});

IconData _icon(CodexRef r) => switch (r.kind) {
  'item' => itemIcon(r.sub),
  'mob' => Icons.pets,
  'dungeon' => r.sub == 'raid' ? Icons.groups_3 : Icons.castle_outlined,
  'region' => Icons.place_outlined,
  'profession' => r.sub == 'gathering' ? activityIcon(r.code) : Icons.handyman_outlined,
  'vendor' => Icons.storefront,
  'quest' => Icons.assignment_outlined,
  'title' => Icons.workspace_premium_outlined,
  'achievement' => Icons.emoji_events_outlined,
  'class' => classIcon(r.sub.isEmpty ? r.code : r.sub),
  'skill' => Icons.flash_on,
  'pet' => Icons.pets,
  'itemset' => Icons.style_outlined,
  'track' => Icons.auto_awesome_outlined,
  'faction' => Icons.flag_outlined,
  'worldevent' => Icons.public,
  'race' => Icons.person_outline,
  'equipment' => Icons.shield_outlined,
  _ => Icons.list,
};

/// The content's own art when it has one, its kind's icon otherwise.
Widget codexArt(CodexRef r, double size) {
  final icon = Icon(_icon(r), size: size * 0.8, color: Palette.textSecondary);
  Widget circle(String asset) => CircleArt(asset: asset, size: size, fallback: icon);
  return SizedBox.square(
    dimension: size,
    child: switch (r.kind) {
      'item' => FittedBox(
        child: ItemArt(code: r.code, icon: _icon(r), color: Palette.textSecondary, size: size / 1.5),
      ),
      'mob' => mobArt(r.code, size),
      'dungeon' => dungeonArt(r.code, size),
      'region' => circle('assets/regions/${r.code}.webp'),
      'profession' => circle('assets/professions/${r.code}.webp'),
      'quest' => circle('assets/quests/${r.code}.webp'),
      'class' when r.sub.isNotEmpty && r.sub != r.code => circle('assets/class_nodes/${r.code}.webp'),
      'class' => circle('assets/class_icons/${r.code}.webp'),
      'skill' => circle('assets/skills/${r.code}.webp'),
      'pet' => circle('assets/pets/${r.code}_1.webp'),
      _ => icon,
    },
  );
}

/// The page a progression event names (`mob_killed:grey_wolf`, `gather_mining`,
/// `dungeon_cleared`); null when it names nothing (`level`).
CodexRef? eventRef(String event) {
  final i = event.indexOf(':');
  final kind = i < 0 ? event : event.substring(0, i);
  final code = i < 0 ? null : event.substring(i + 1);
  if (kind.startsWith('gather_')) return CodexRef('profession', kind.substring(7), sub: 'gathering');
  return switch (kind) {
    'mob_killed' => code == null ? const CodexRef('list', 'mobs') : CodexRef('mob', code),
    'dungeon_cleared' => code == null ? const CodexRef('list', 'dungeons') : CodexRef('dungeon', code),
    'raid_cleared' => code == null ? const CodexRef('list', 'dungeons') : CodexRef('dungeon', code, sub: 'raid'),
    'boss_killed' => const CodexRef('list', 'dungeons'),
    'item_crafted' => code == null ? const CodexRef('list', 'crafts') : CodexRef('item', code),
    'class_mastery' when code != null => CodexRef('class', code, sub: 'spec'),
    _ => null,
  };
}

/// What to do, in a sentence: "Kill 1000 monsters. Hunts, ... all count."
String goalText(AppLocalizations l, String event, int n) {
  final i = event.indexOf(':');
  final kind = i < 0 ? event : event.substring(0, i);
  final code = i < 0 ? null : event.substring(i + 1);
  if (kind.startsWith('gather_')) return l.goalGather(n, professionLabel(l, kind.substring(7)));
  return switch (kind) {
    'level' => l.goalLevel(n),
    'class_mastery' => l.goalMastery(n, codexName(l, 'class', code ?? '')),
    'mob_killed' => code == null ? l.goalKillAny(n) : l.goalKill(n, ContentNames.of('mob', code)),
    'boss_killed' => l.goalBoss(n),
    'dungeon_cleared' => code == null ? l.goalDungeonAny(n) : l.goalDungeon(n, ContentNames.of('dungeon', code)),
    'raid_cleared' => code == null ? l.goalRaidAny(n) : l.goalRaid(n, ContentNames.of('dungeon', code)),
    'item_crafted' => code == null ? l.goalCraftAny(n) : l.goalCraft(n, itemName(code)),
    _ => '${humanize(event)} ×$n',
  };
}

/// Opens the page of [r] over whatever is shown.
Future<void> showCodex(BuildContext context, CodexRef r) => showModalBottomSheet<void>(
  context: context,
  isScrollControlled: true,
  showDragHandle: true,
  backgroundColor: Palette.bgSurface,
  builder: (_) => CodexSheet(start: r),
);

/// A tappable tag: art, name and, when known, the level, profession level or
/// quantity. Opens the page of what it names.
class LinkTag extends StatelessWidget {
  const LinkTag(this.target, {super.key, this.label, this.color, this.onOpen});

  final CodexRef target;

  /// Replaces the name and numbers.
  final String? label;
  final Color? color;

  /// Replaces [showCodex] (a page follows its own links in place).
  final ValueChanged<CodexRef>? onOpen;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final c = color ?? Palette.accent;
    final text =
        label ??
        [
          target.count > 0
              ? '${codexName(l, target.kind, target.code)} ×${target.count}'
              : codexName(l, target.kind, target.code),
          if (target.skill > 0)
            l.tagSkill(target.skill)
          else if (target.level > 0 && target.kind != 'class')
            l.levelShort(target.level),
        ].join(' · ');
    return Material(
      color: Palette.bgElevated,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(Radii.sm),
        side: BorderSide(color: c.withValues(alpha: 0.5)),
      ),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () => onOpen != null ? onOpen!(target) : showCodex(context, target),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: 6),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              codexArt(target, 20),
              const SizedBox(width: 6),
              Flexible(
                child: Text(text, style: TextStyle(fontSize: 12, color: color ?? Palette.textPrimary)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Opens a page from a row whose own tap does something else.
class CodexButton extends StatelessWidget {
  const CodexButton(this.target, {super.key});

  final CodexRef target;

  @override
  Widget build(BuildContext context) => IconButton(
    tooltip: AppLocalizations.of(context).codexOpen,
    visualDensity: VisualDensity.compact,
    icon: const Icon(Icons.info_outline, size: 20, color: Palette.textSecondary),
    onPressed: () => showCodex(context, target),
  );
}

/// Tags side by side.
class LinkTags extends StatelessWidget {
  const LinkTags(this.refs, {super.key, this.onOpen});

  final Iterable<CodexRef> refs;
  final ValueChanged<CodexRef>? onOpen;

  @override
  Widget build(BuildContext context) => Wrap(
    spacing: Space.xs,
    runSpacing: Space.xs,
    children: [for (final r in refs) LinkTag(r, onOpen: onOpen)],
  );
}

/// The relations of a page as titled tag groups; long groups fold.
class CodexGroups extends StatefulWidget {
  const CodexGroups({super.key, required this.groups, this.onOpen});

  final List<CodexGroup> groups;
  final ValueChanged<CodexRef>? onOpen;

  @override
  State<CodexGroups> createState() => _CodexGroupsState();
}

class _CodexGroupsState extends State<CodexGroups> {
  /// Tags shown before a group folds.
  static const _shown = 12;

  final _open = <String>{};

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final g in widget.groups) ...[
          const SizedBox(height: Space.md),
          Text(l.codexGroup(g.key), style: text.bodySmall),
          const SizedBox(height: Space.xs),
          if (g.refs.length <= _shown + 1 || _open.contains(g.key))
            LinkTags(g.refs, onOpen: widget.onOpen)
          else ...[
            LinkTags(g.refs.take(_shown), onOpen: widget.onOpen),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                onPressed: () => setState(() => _open.add(g.key)),
                icon: const Icon(Icons.expand_more, size: 18),
                label: Text(l.codexMore(g.refs.length - _shown)),
              ),
            ),
          ],
        ],
      ],
    );
  }
}

/// The tags of one piece of content, inside another page (item details).
/// Shows nothing until (and unless) they load.
class CodexTags extends StatelessWidget {
  const CodexTags({super.key, required this.kind, required this.code});

  final String kind;
  final String code;

  @override
  Widget build(BuildContext context) => FutureBuilder<CodexPage?>(
    future: Codex.page(kind, code),
    builder: (context, snap) {
      final page = snap.data;
      if (page == null || page.groups.isEmpty) return const SizedBox.shrink();
      return CodexGroups(key: const Key('item-tags'), groups: page.groups);
    },
  );
}

/// One thing to do and what it is about.
class GoalLine extends StatelessWidget {
  const GoalLine({super.key, required this.event, required this.count, this.onOpen});

  final String event;
  final int count;
  final ValueChanged<CodexRef>? onOpen;

  @override
  Widget build(BuildContext context) {
    final ref = eventRef(event);
    return Padding(
      padding: const EdgeInsets.only(top: Space.xs),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(goalText(AppLocalizations.of(context), event, count)),
          if (ref != null) ...[const SizedBox(height: Space.xs), LinkTag(ref, onOpen: onOpen)],
        ],
      ),
    );
  }
}

/// A content page; its tags open in place with a way back.
class CodexSheet extends StatefulWidget {
  const CodexSheet({super.key, required this.start});

  final CodexRef start;

  @override
  State<CodexSheet> createState() => _CodexSheetState();
}

class _CodexSheetState extends State<CodexSheet> {
  late final _path = [widget.start];

  void _open(CodexRef r) => setState(() => _path.add(r));

  void _back() => setState(_path.removeLast);

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final here = _path.last;
    return PopScope(
      canPop: _path.length == 1,
      onPopInvokedWithResult: (didPop, _) {
        if (!didPop) _back();
      },
      child: SafeArea(
        child: SizedBox(
          height: MediaQuery.sizeOf(context).height * 0.75, // same height for every page
          child: PageFrame(
            child: FutureBuilder<CodexPage?>(
              key: ValueKey('${here.kind}/${here.code}#${_path.length}'),
              future: Codex.page(here.kind, here.code),
              builder: (context, snap) => ListView(
                key: const Key('codex-page'),
                padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
                children: [
                  _header(l, snap.data?.ref ?? here),
                  if (snap.data case final page?)
                    ..._body(l, page)
                  else if (snap.connectionState == ConnectionState.done)
                    Padding(
                      padding: const EdgeInsets.only(top: Space.xl),
                      child: Text(l.codexFailed, textAlign: TextAlign.center),
                    )
                  else
                    const Padding(
                      padding: EdgeInsets.only(top: Space.xxl),
                      child: Center(child: CircularProgressIndicator()),
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _header(AppLocalizations l, CodexRef r) {
    final text = Theme.of(context).textTheme;
    return Row(
      children: [
        if (_path.length > 1)
          IconButton(
            key: const Key('codex-back'),
            tooltip: MaterialLocalizations.of(context).backButtonTooltip,
            icon: const Icon(Icons.arrow_back),
            onPressed: _back,
          ),
        codexArt(r, 56),
        const SizedBox(width: Space.lg),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(codexName(l, r.kind, r.code), style: text.titleLarge),
              Text(
                [codexKindName(l, r), if (r.level > 0) l.levelShort(r.level)].join(' · '),
                style: text.bodySmall?.copyWith(color: Palette.accent),
              ),
            ],
          ),
        ),
        IconButton(
          tooltip: MaterialLocalizations.of(context).closeButtonTooltip,
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ],
    );
  }

  List<Widget> _body(AppLocalizations l, CodexPage p) {
    final text = Theme.of(context).textTheme;
    final about = switch (p.ref.kind) {
      'skill' => ContentNames.of('skill', '${p.ref.code}.desc', ''),
      'class' => ContentNames.of('class', '${p.ref.code}.intro', ContentNames.of('class', '${p.ref.code}.role', '')),
      'pet' => ContentNames.of('pet', '${p.ref.code}.story', ''),
      _ => '',
    };
    Widget header(String s) => Padding(
      padding: const EdgeInsets.only(top: Space.lg),
      child: Text(s, style: text.titleSmall?.copyWith(color: Palette.accent)),
    );
    return [
      if (about.isNotEmpty) ...[const SizedBox(height: Space.md), Text(about)],
      for (final f in _facts(l, p))
        Padding(
          padding: const EdgeInsets.only(top: Space.sm),
          child: Text(f, style: const TextStyle(fontFeatures: numeric)),
        ),
      if (p.goals.isNotEmpty) ...[
        header(p.ref.kind == 'quest' ? l.codexGoals : l.codexHow),
        for (final g in p.goals) GoalLine(event: g.event, count: g.count, onOpen: _open),
      ],
      if (p.bonuses.isNotEmpty) ...[
        header(l.codexBonuses),
        for (final (b, at) in p.bonuses)
          if (bonusText(l, b) case final t?)
            Text(
              at == 0
                  ? t
                  : '${switch (p.ref.kind) {
                      'itemset' => l.codexPieces(at),
                      'track' => l.codexTier(at),
                      _ => l.levelShort(at),
                    }} $t',
              style: const TextStyle(color: Palette.success, fontWeight: FontWeight.w600),
            ),
      ],
      CodexGroups(groups: p.groups, onOpen: _open),
    ];
  }

  /// The numbers of a page, readable.
  List<String> _facts(AppLocalizations l, CodexPage p) {
    final f = p.facts;
    int n(String k) => f[k] ?? 0;
    final rewards = [
      if (n('points') > 0) l.pointsShort(n('points')),
      if (p.ref.kind == 'quest' && n('xp') > 0) l.xpGain(n('xp')),
      if (n('gold') > 0) goldString(n('gold')),
    ];
    return [
      if (p.traits['slot'] case final slot?)
        gearKind(l, slot, p.traits['subclass'] ?? '')
      else if (p.traits['element'] case final e?)
        l.petElement(ContentNames.of('element', e)),
      if (p.ref.kind == 'mob') l.mobStats(n('hp'), n('attack'), n('defense')),
      if (p.ref.kind == 'mob' && n('xp') > 0) l.xpPerKill(n('xp')),
      if (n('group_max') > 1) l.raidGroup(n('group_min'), n('group_max')),
      if (p.ref.kind == 'dungeon') shortDuration(l, Duration(seconds: n('seconds'))),
      if (n('heal_hp') > 0) l.restoresHp(n('heal_hp')),
      if (n('heal_pct') > 0) l.restoresHpPct(n('heal_pct')),
      if (n('mana_pct') > 0) l.restoresMpPct(n('mana_pct')),
      if (n('buff_minutes') > 0) l.foodBuffLine(n('buff_attack_pct'), n('buff_defense_pct'), n('buff_minutes')),
      if (rewards.isNotEmpty) l.codexRewards(rewards.join(' · ')),
    ];
  }
}
