import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart';
import '../../shared/hold_repeat_button.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';

// Character leveling + class progression (DD-250): level 1..1000, class
// promotions at 100 and 300, attribute and talent points, milestones.

String specName(String code) => ContentNames.of('class', code);
String specRole(String code) => ContentNames.of('class', '$code.role', '');
String attributeName(String code) => ContentNames.of('attribute', code);
String talentName(String code) => ContentNames.of('talent', code);
String talentDesc(String code) => ContentNames.of('talent', '$code.desc', '');
String titleName(String code) => ContentNames.of('title', code);

class Milestone {
  const Milestone({required this.level, required this.kind, this.ref = '', this.minutes = 0});

  factory Milestone.fromJson(Map<String, dynamic> j) => Milestone(
    level: (j['level'] as num).toInt(),
    kind: j['kind'] as String? ?? '',
    ref: j['ref'] as String? ?? '',
    minutes: (j['estimated_minutes'] as num?)?.toDouble() ?? 0,
  );

  final int level;
  final String kind; // promotion | awakening | ultimate | capstone | mastery | skill | passive | title
  final String ref;
  final double minutes;
}

class LevelingSummary {
  const LevelingSummary({
    required this.level,
    required this.levelCap,
    required this.xp,
    required this.xpToNext,
    required this.maxed,
    required this.classCode,
    required this.spec,
    required this.path,
    this.roleTags = const [],
    this.promotionAvailable = false,
    this.overflowXp = 0,
    this.choices = const [],
    this.attributePoints = 0,
    this.talentPoints = 0,
    this.raw = const {},
    this.effective = const {},
    this.gains = const {},
    this.title = '',
    this.masteryTitle = '',
    this.masteryUnlocked = false,
    this.next,
  });

  factory LevelingSummary.fromJson(Map<String, dynamic> j) {
    int n(Object? v) => (v as num?)?.toInt() ?? 0;
    final promo = j['promotion'] as Map<String, dynamic>? ?? const {};
    final attrs = j['attributes'] as Map<String, dynamic>? ?? const {};
    final mastery = j['mastery'] as Map<String, dynamic>? ?? const {};
    return LevelingSummary(
      level: n(j['level']),
      levelCap: n(j['level_cap']),
      xp: n(j['xp']),
      xpToNext: n(j['xp_to_next']),
      maxed: j['maxed'] == true,
      classCode: j['class_code'] as String? ?? '',
      spec: j['spec_code'] as String? ?? '',
      path: [for (final p in (j['path'] as List? ?? const [])) p as String],
      roleTags: [for (final r in (j['role_tags'] as List? ?? const [])) r as String],
      promotionAvailable: promo['available'] == true,
      overflowXp: n(promo['overflow_xp']),
      choices: [
        for (final c in (promo['choices'] as List? ?? const []).cast<Map<String, dynamic>>()) c['code'] as String,
      ],
      attributePoints: n(attrs['available']),
      talentPoints: n((j['talents'] as Map<String, dynamic>?)?['available']),
      raw: {for (final e in (attrs['raw'] as Map<String, dynamic>? ?? const {}).entries) e.key: n(e.value)},
      effective: {
        for (final e in (attrs['effective'] as Map<String, dynamic>? ?? const {}).entries)
          e.key: (e.value as num).toDouble(),
      },
      gains: {
        for (final e in (attrs['gains'] as Map<String, dynamic>? ?? const {}).entries)
          e.key: {
            for (final g in (e.value as Map<String, dynamic>? ?? const {}).entries) g.key: (g.value as num).toDouble(),
          },
      },
      title: j['title'] as String? ?? '',
      masteryTitle: mastery['title'] as String? ?? '',
      masteryUnlocked: mastery['unlocked'] == true,
      next: j['next_milestone'] is Map<String, dynamic>
          ? Milestone.fromJson(j['next_milestone'] as Map<String, dynamic>)
          : null,
    );
  }

  final int level;
  final int levelCap;
  final int xp;
  final int xpToNext;
  final bool maxed;
  final String classCode;
  final String spec;
  final List<String> path;
  final List<String> roleTags;
  final bool promotionAvailable;
  final int overflowXp;
  final List<String> choices;
  final int attributePoints;
  final int talentPoints;
  final Map<String, int> raw;
  final Map<String, double> effective;

  /// What one more point of each attribute adds to the class now:
  /// attribute -> stat (hp, mp, attack, defense, crit, hit, dodge, block,
  /// hp_regen) -> amount.
  final Map<String, Map<String, double>> gains;
  final String title;
  final String masteryTitle;
  final bool masteryUnlocked;
  final Milestone? next;

  double get xpProgress => maxed || xpToNext <= 0 ? 1 : (xp / xpToNext).clamp(0, 1).toDouble();
}

class TalentNode {
  const TalentNode({required this.code, required this.maxRank, required this.rank, this.capstone = false});

  final String code;
  final int maxRank;
  final int rank;
  final bool capstone;
}

class TalentTree {
  const TalentTree({required this.code, required this.points, required this.talents});

  final String code;
  final int points;
  final List<TalentNode> talents;
}

class TalentBoard {
  const TalentBoard({
    required this.available,
    required this.level,
    required this.treePoints,
    required this.capstoneLevel,
    required this.trees,
  });

  final int available;
  final int level;
  final int treePoints;
  final int capstoneLevel;
  final List<TalentTree> trees;
}

/// Leveling endpoints of one character.
class LevelingRepository {
  const LevelingRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _char => '/v1/character/$characterId';

  Future<LevelingSummary> summary() async => LevelingSummary.fromJson(await api.get('$_char/leveling'));

  Future<void> promote(String choice) => api.post('$_char/promotion', {'choice': choice});

  Future<void> allocate(Map<String, int> points) => api.post('$_char/attributes', points);

  Future<TalentBoard> talents() async {
    final j = await api.get('$_char/talents');
    int n(Object? v) => (v as num?)?.toInt() ?? 0;
    return TalentBoard(
      available: n(j['available']),
      level: n(j['level']),
      treePoints: n(j['tree_points']),
      capstoneLevel: n(j['capstone_level']),
      trees: [
        for (final t in (j['trees'] as List? ?? const []).cast<Map<String, dynamic>>())
          TalentTree(
            code: t['code'] as String,
            points: n(t['points']),
            talents: [
              for (final x in (t['talents'] as List? ?? const []).cast<Map<String, dynamic>>())
                TalentNode(
                  code: x['code'] as String,
                  maxRank: n(x['max_rank']),
                  rank: n(x['rank']),
                  capstone: x['capstone'] == true,
                ),
            ],
          ),
      ],
    );
  }

  Future<void> learnTalent(String code) => api.post('$_char/talents/$code', const {});
}

String milestoneText(AppLocalizations l, Milestone m) => switch (m.kind) {
  'promotion' => l.msPromotion,
  'awakening' => l.msAwakening(ContentNames.of('skill', m.ref)),
  'ultimate' => l.msUltimate(ContentNames.of('skill', m.ref)),
  'capstone' => l.msCapstone,
  'mastery' => l.msMastery,
  'title' => l.msTitle(titleName(m.ref)),
  _ => ContentNames.of('skill', m.ref),
};

String hoursText(AppLocalizations l, double minutes) =>
    minutes < 60 ? l.estMinutes(minutes.ceil()) : l.estHours((minutes / 60).toStringAsFixed(1));

/// Level, XP (MAX at the cap), class / specialization, title, next
/// milestone, and the actions the progression offers right now.
class LevelingCard extends StatelessWidget {
  const LevelingCard({
    super.key,
    required this.summary,
    required this.onPromote,
    required this.onAttributes,
    required this.onTalents,
  });

  final LevelingSummary summary;
  final VoidCallback onPromote;
  final VoidCallback onAttributes;
  final VoidCallback onTalents;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final s = summary;
    final text = Theme.of(context).textTheme;
    final name = [if (s.title.isNotEmpty) titleName(s.title), specName(s.spec)].join(' ');
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Row(
          children: [
            Expanded(
              child: Text(name, style: text.titleMedium?.copyWith(color: Palette.accent)),
            ),
            Text(l.levelShort(s.level), style: text.titleMedium),
          ],
        ),
        if (specRole(s.spec).isNotEmpty) Text(specRole(s.spec), style: text.bodySmall),
        const SizedBox(height: Space.sm),
        StatBar(
          label: l.xp,
          value: s.xpProgress,
          color: Palette.epic,
          valueText: s.maxed ? l.levelMax : l.xpValue(s.xp, s.xpToNext),
        ),
        if (s.masteryUnlocked && s.masteryTitle.isNotEmpty)
          Padding(
            padding: const EdgeInsets.only(top: Space.xs),
            child: Text(
              l.masteryUnlocked(titleName(s.masteryTitle)),
              style: TextStyle(color: rarityColor('legendary'), fontWeight: FontWeight.w700),
            ),
          ),
        if (s.next case final m? when !s.promotionAvailable)
          Padding(
            padding: const EdgeInsets.only(top: Space.xs),
            child: Text(
              '${l.nextMilestone(m.level, milestoneText(l, m))} · ${hoursText(l, m.minutes)}',
              style: text.bodySmall,
            ),
          ),
        if (s.promotionAvailable) ...[
          const SizedBox(height: Space.sm),
          FilledButton.icon(
            onPressed: onPromote,
            icon: const Icon(Icons.auto_awesome),
            label: Text(l.promotionAvailable),
          ),
          if (s.overflowXp > 0)
            Padding(
              padding: const EdgeInsets.only(top: Space.xs),
              child: Text(l.overflowHeld(s.overflowXp), style: text.bodySmall, textAlign: TextAlign.center),
            ),
        ],
        const SizedBox(height: Space.sm),
        Row(
          children: [
            Expanded(
              child: OutlinedButton(
                onPressed: onAttributes,
                child: Text(l.attributePointsBtn(s.attributePoints), textAlign: TextAlign.center),
              ),
            ),
            const SizedBox(width: Space.sm),
            Expanded(
              child: OutlinedButton(
                onPressed: s.level >= 100 ? onTalents : null,
                child: Text(l.talentPointsBtn(s.talentPoints), textAlign: TextAlign.center),
              ),
            ),
          ],
        ),
      ],
    );
  }
}

/// Promotion choice (Level 100 / 300): names, roles, confirm.
Future<String?> showPromotionSheet(BuildContext context, LevelingSummary s) {
  return showModalBottomSheet<String>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    backgroundColor: Palette.bgSurface,
    builder: (c) {
      final l = AppLocalizations.of(c);
      return SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l.promotionAvailable, style: Theme.of(c).textTheme.titleMedium, textAlign: TextAlign.center),
              const SizedBox(height: Space.xs),
              Text(l.promotionBody, style: Theme.of(c).textTheme.bodySmall, textAlign: TextAlign.center),
              const SizedBox(height: Space.md),
              for (final code in s.choices)
                Card(
                  child: ListTile(
                    title: Text(specName(code), style: const TextStyle(fontWeight: FontWeight.w700)),
                    subtitle: Text(specRole(code)),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () async {
                      final ok = await confirmAction(
                        c,
                        title: l.promotionConfirm(specName(code)),
                        body: l.promotionBody,
                        action: l.promotionChoose,
                      );
                      if (ok && c.mounted) Navigator.pop(c, code);
                    },
                  ),
                ),
            ],
          ),
        ),
      );
    },
  );
}

/// Attribute allocation: +/- per attribute within the available points
/// (hold to repeat), each with what it adds to the class.
class AttributeSheet extends StatefulWidget {
  const AttributeSheet({super.key, required this.summary});

  final LevelingSummary summary;

  @override
  State<AttributeSheet> createState() => _AttributeSheetState();
}

const attributeOrder = ['strength', 'vitality', 'dexterity', 'intelligence', 'wisdom', 'luck'];

class _AttributeSheetState extends State<AttributeSheet> {
  final Map<String, int> _add = {};

  int get _left => widget.summary.attributePoints - _add.values.fold(0, (a, b) => a + b);

  // Checked at call time: a held button can tick again before the rebuild
  // that disables it.
  void _step(String a, int by) {
    final next = (_add[a] ?? 0) + by;
    if (next < 0 || (by > 0 && _left <= 0)) return;
    setState(() => _add[a] = next);
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final s = widget.summary;
    final text = Theme.of(context).textTheme;
    return SafeArea(
      child: SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(l.attributePointsBtn(_left), style: text.titleMedium),
            Text(l.attrHoldHint, style: text.bodySmall),
            const SizedBox(height: Space.sm),
            for (final a in attributeOrder) ...[
              Row(
                children: [
                  Expanded(child: Text(attributeName(a))),
                  Text(
                    '${(s.raw[a] ?? 0) + (_add[a] ?? 0)}',
                    style: const TextStyle(fontWeight: FontWeight.w700, fontFeatures: numeric),
                  ),
                  HoldRepeatButton(
                    tooltip: '- ${attributeName(a)}',
                    icon: Icons.remove_circle_outline,
                    onStep: (_add[a] ?? 0) > 0 ? () => _step(a, -1) : null,
                  ),
                  HoldRepeatButton(
                    tooltip: '+ ${attributeName(a)}',
                    icon: Icons.add_circle_outline,
                    onStep: _left > 0 ? () => _step(a, 1) : null,
                  ),
                ],
              ),
              _GainLine(gains: s.gains[a] ?? const {}, adding: _add[a] ?? 0),
              const SizedBox(height: Space.xs),
            ],
            const SizedBox(height: Space.sm),
            FilledButton(
              onPressed: _add.values.any((v) => v > 0)
                  ? () => Navigator.pop(context, {
                      for (final e in _add.entries)
                        if (e.value > 0) e.key: e.value,
                    })
                  : null,
              child: Text(l.save),
            ),
          ],
        ),
      ),
    );
  }
}

/// What a point of one attribute adds (and, while points are pending, what
/// they will add). Gains are per point at the current diminishing-returns
/// step, so a large pending sum is a close estimate; the server decides.
class _GainLine extends StatelessWidget {
  const _GainLine({required this.gains, required this.adding});

  final Map<String, double> gains;
  final int adding;

  static const _order = ['attack', 'hp', 'defense', 'mp', 'crit', 'hit', 'dodge', 'block', 'hp_regen'];

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final small = Theme.of(context).textTheme.bodySmall;
    if (gains.isEmpty) return Text(l.attrNoEffect, style: small?.copyWith(color: Palette.textMuted));
    final num = NumberFormat.decimalPatternDigits(locale: Localizations.localeOf(context).toString(), decimalDigits: 2);
    String fmt(double v) {
      final t = num.format(v);
      final sep = num.symbols.DECIMAL_SEP;
      return t.contains(sep) ? t.replaceFirst(RegExp('${RegExp.escape(sep)}?0+\$'), '') : t;
    }

    String list(int times) => [
      for (final stat in _order)
        if (gains[stat] case final v?) _label(l, stat, '+${fmt(v * times)}'),
    ].join(' · ');
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(l.attrGainPerPoint(list(1)), style: small),
        if (adding > 0)
          Text(
            l.attrGainPending(list(adding)),
            style: small?.copyWith(color: Palette.accent, fontWeight: FontWeight.w600),
          ),
      ],
    );
  }

  static String _label(AppLocalizations l, String stat, String v) => switch (stat) {
    'hp' => '$v ${l.gainHp}',
    'mp' => '$v ${l.gainMp}',
    'attack' => '$v ${l.attack}',
    'defense' => '$v ${l.defense}',
    'crit' => '$v% ${l.crit}',
    'hit' => '$v% ${l.hit}',
    'dodge' => '$v${l.statEvasion}',
    'block' => '$v${l.statBlock}',
    'hp_regen' => '$v ${l.gainHpRegen}',
    _ => '$v $stat',
  };
}

/// Talent trees (from Level 100): spend points, capstones at the capstone
/// level with a full tree.
class TalentsScreen extends StatefulWidget {
  const TalentsScreen({super.key, required this.repo, required this.dataVersion});

  final LevelingRepository repo;
  final ValueNotifier<int> dataVersion;

  @override
  State<TalentsScreen> createState() => _TalentsScreenState();
}

class _TalentsScreenState extends State<TalentsScreen> with PanelState {
  TalentBoard? _board;

  @override
  void initState() {
    super.initState();
    load();
  }

  @override
  Future<void> fetch() async => _board = await widget.repo.talents();

  Future<void> _learn(String code) => run(() async {
    await widget.repo.learnTalent(code);
    widget.dataVersion.value++;
    await load();
  });

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final b = _board;
    return Scaffold(
      appBar: AppBar(title: Text(l.talentsTitle)),
      body: SafeArea(
        child: PageFrame(
          child: frame(
            () => [
              if (b != null) ...[
                Text(l.talentPointsBtn(b.available), style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: Space.md),
                for (final t in b.trees)
                  Card(
                    margin: const EdgeInsets.only(bottom: Space.md),
                    child: Padding(
                      padding: const EdgeInsets.all(Space.md),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          Text(
                            '${ContentNames.of('talent_tree', t.code)} · ${t.points}/${b.treePoints}',
                            style: Theme.of(context).textTheme.titleSmall,
                          ),
                          for (final n in t.talents)
                            _TalentRow(node: n, board: b, tree: t, onLearn: () => _learn(n.code)),
                        ],
                      ),
                    ),
                  ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _TalentRow extends StatelessWidget {
  const _TalentRow({required this.node, required this.board, required this.tree, required this.onLearn});

  final TalentNode node;
  final TalentBoard board;
  final TalentTree tree;
  final VoidCallback onLearn;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final capLocked = node.capstone && (board.level < board.capstoneLevel || tree.points < board.treePoints);
    final can = board.available > 0 && node.rank < node.maxRank && !capLocked;
    return Padding(
      padding: const EdgeInsets.only(top: Space.sm),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  talentName(node.code),
                  style: TextStyle(fontWeight: FontWeight.w600, color: node.capstone ? rarityColor('legendary') : null),
                ),
                if (talentDesc(node.code).isNotEmpty)
                  Text(talentDesc(node.code), style: Theme.of(context).textTheme.bodySmall),
                if (capLocked)
                  Text(
                    l.capstoneLocked(board.capstoneLevel, board.treePoints),
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Palette.warning),
                  ),
              ],
            ),
          ),
          Text('${node.rank}/${node.maxRank}', style: const TextStyle(fontFeatures: numeric)),
          IconButton(
            tooltip: '+ ${talentName(node.code)}',
            icon: const Icon(Icons.add_circle_outline),
            onPressed: can ? onLearn : null,
          ),
        ],
      ),
    );
  }
}

/// Runs a leveling action with the error toast the app uses everywhere.
Future<bool> levelingAction(BuildContext context, Future<void> Function() action) async {
  try {
    await action();
    return true;
  } on ApiError catch (e) {
    if (context.mounted && !e.isSessionError) toast(context, Text(errorText(AppLocalizations.of(context), e)));
    return false;
  }
}
