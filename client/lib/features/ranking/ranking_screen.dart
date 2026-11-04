import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';
import '../character/class_catalog.dart';
import '../character/models.dart';
import '../social/social_repository.dart';
import '../social/social_screen.dart' show showReportDialog;

/// One row of a ranking table (DD-357).
class RankEntry {
  const RankEntry.fromJson(this.j);

  final Map<String, dynamic> j;

  int get rank => (j['rank'] as num?)?.toInt() ?? 0;
  String get id => j['character_id'] as String? ?? '';
  String get name => j['name'] as String? ?? '';
  String get classCode => j['class_code'] as String? ?? '';
  int get level => (j['level'] as num?)?.toInt() ?? 0;
  int get value => (j['value'] as num?)?.toInt() ?? 0;
  int get score => (j['score'] as num?)?.toInt() ?? 0;
  int count(String key) => (j[key] as num?)?.toInt() ?? 0;
}

class RankPage {
  const RankPage(this.entries, this.me, this.points);

  final List<RankEntry> entries;
  final RankEntry? me;
  final Map<String, int> points;
}

class RankingRepository {
  const RankingRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  Future<RankPage> page(String board, String? cls, int offset) async {
    final q = {'board': board, 'character_id': characterId, 'limit': '$_pageSize', 'offset': '$offset', 'class': ?cls};
    final j = await api.get(Uri(path: '/v1/leaderboard', queryParameters: q).toString());
    return RankPage(
      [for (final e in (j['entries'] as List? ?? const []).cast<Map<String, dynamic>>()) RankEntry.fromJson(e)],
      j['me'] is Map<String, dynamic> ? RankEntry.fromJson(j['me'] as Map<String, dynamic>) : null,
      {for (final e in (j['points'] as Map<String, dynamic>? ?? const {}).entries) e.key: (e.value as num).toInt()},
    );
  }
}

const _pageSize = 50;

/// The ranking tables: character score and one table per measure. There is no
/// player-versus-player here, only what characters have achieved.
class RankingScreen extends StatefulWidget {
  const RankingScreen({super.key, required this.repo, required this.social});

  final RankingRepository repo;
  final SocialRepository social;

  @override
  State<RankingScreen> createState() => _RankingScreenState();
}

class _RankingScreenState extends State<RankingScreen> {
  String _board = 'score';
  String? _class;
  RankPage? _page;
  List<RankEntry> _rows = const [];
  String? _error;
  bool _more = false;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load({bool append = false}) async {
    final board = _board;
    final cls = _class;
    try {
      final p = await widget.repo.page(board, cls, append ? _rows.length : 0);
      if (!mounted || board != _board || cls != _class) return;
      setState(() {
        _error = null;
        _page = p;
        _rows = append ? [..._rows, ...p.entries] : p.entries;
        _more = p.entries.length == _pageSize;
      });
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = errorText(AppLocalizations.of(context), e));
    }
  }

  void _pick(String board) {
    setState(() {
      _board = board;
      _page = null;
    });
    _load();
  }

  void _breakdown(RankEntry me, Map<String, int> points) {
    final l = AppLocalizations.of(context);
    final rows = [
      (l.rankingLevel, me.level, points['level'] ?? 0),
      (l.rankingAchievements, me.count('achievement_points'), points['achievement'] ?? 0),
      (l.rankingMastery, me.count('mastery_tiers'), points['mastery_tier'] ?? 0),
      (l.rankingProfessions, me.count('profession_levels'), points['profession_level'] ?? 0),
      (l.rankingDungeons, me.count('dungeons_cleared'), points['dungeon'] ?? 0),
      (l.rankingRaids, me.count('raids_cleared'), points['raid'] ?? 0),
    ];
    showModalBottomSheet<void>(
      context: context,
      builder: (sheet) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(Space.screen),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l.rankingHowScore, style: Theme.of(sheet).textTheme.titleMedium),
              const SizedBox(height: Space.sm),
              for (final (name, count, each) in rows)
                Padding(
                  padding: const EdgeInsets.symmetric(vertical: 3),
                  child: Row(
                    children: [
                      Expanded(child: Text(name)),
                      Text('$count × $each', style: const TextStyle(color: Palette.textMuted)),
                      const SizedBox(width: Space.md),
                      SizedBox(width: 72, child: Text('${count * each}', textAlign: TextAlign.end)),
                    ],
                  ),
                ),
              const Divider(),
              Row(
                children: [
                  Expanded(
                    child: Text(l.rankingScore, style: const TextStyle(fontWeight: FontWeight.w700)),
                  ),
                  Text('${me.score}', style: const TextStyle(fontWeight: FontWeight.w700)),
                ],
              ),
              const SizedBox(height: Space.sm),
              Text(l.rankingScoreNote, style: const TextStyle(fontSize: 12, color: Palette.textMuted)),
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final boards = {
      'score': l.rankingScore,
      'level': l.rankingLevel,
      'achievements': l.rankingAchievements,
      'mastery': l.rankingMastery,
      'professions': l.rankingProfessions,
      'clears': l.rankingClears,
    };
    final me = _page?.me;
    return Scaffold(
      appBar: AppBar(title: Text(l.ranking)),
      body: SafeArea(
        child: PageFrame(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              SizedBox(
                height: 52,
                child: ListView(
                  scrollDirection: Axis.horizontal,
                  padding: const EdgeInsets.symmetric(horizontal: Space.screen, vertical: Space.xs),
                  children: [
                    for (final e in boards.entries)
                      Padding(
                        padding: const EdgeInsets.only(right: Space.xs),
                        child: ChoiceChip(
                          key: ValueKey('board-${e.key}'),
                          label: Text(e.value),
                          selected: _board == e.key,
                          onSelected: (_) => _pick(e.key),
                        ),
                      ),
                  ],
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: Space.screen),
                child: DropdownButton<String?>(
                  key: const ValueKey('ranking-class'),
                  isExpanded: true,
                  value: _class,
                  items: [
                    DropdownMenuItem(child: Text(l.rankingAllClasses)),
                    for (final c in classes)
                      DropdownMenuItem(value: c, child: Text(ContentNames.of('class', c, className(l, c)))),
                  ],
                  onChanged: (v) {
                    setState(() {
                      _class = v;
                      _page = null;
                    });
                    _load();
                  },
                ),
              ),
              if (me != null)
                Padding(
                  padding: const EdgeInsets.fromLTRB(Space.screen, Space.xs, Space.screen, Space.sm),
                  child: _Row(
                    key: const ValueKey('ranking-me'),
                    entry: me,
                    highlighted: true,
                    onTap: () => _breakdown(me, _page?.points ?? const {}),
                  ),
                ),
              const Divider(height: 1),
              Expanded(
                child: _error != null && _page == null
                    ? MessageView(message: _error!, onRetry: _load)
                    : _page == null
                    ? const Center(child: CircularProgressIndicator())
                    : _rows.isEmpty
                    ? Center(child: Text(l.rankingEmpty))
                    : RefreshIndicator(
                        onRefresh: _load,
                        child: ListView.builder(
                          key: const ValueKey('ranking-list'),
                          padding: const EdgeInsets.all(Space.screen),
                          itemCount: _rows.length + (_more ? 1 : 0),
                          itemBuilder: (context, i) {
                            if (i == _rows.length) {
                              return TextButton(onPressed: () => _load(append: true), child: Text(l.loadMore));
                            }
                            final e = _rows[i];
                            final mine = e.id == widget.repo.characterId;
                            return Padding(
                              padding: const EdgeInsets.only(bottom: Space.xs),
                              child: _Row(
                                entry: e,
                                highlighted: mine,
                                onLongPress: mine ? null : () => showReportDialog(context, widget.social, e.name),
                              ),
                            );
                          },
                        ),
                      ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _Row extends StatelessWidget {
  const _Row({super.key, required this.entry, this.highlighted = false, this.onTap, this.onLongPress});

  final RankEntry entry;
  final bool highlighted;
  final VoidCallback? onTap;
  final VoidCallback? onLongPress;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final medal = switch (entry.rank) {
      1 => const Color(0xFFE8C14F),
      2 => const Color(0xFFC0C4CC),
      3 => const Color(0xFFCD8B56),
      _ => Palette.textMuted,
    };
    return Material(
      color: highlighted ? Palette.accent.withValues(alpha: 0.12) : Palette.bgSurface,
      borderRadius: BorderRadius.circular(Radii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(Radii.md),
        onTap: onTap,
        onLongPress: onLongPress,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: Space.sm, vertical: Space.xs),
          child: Row(
            children: [
              SizedBox(
                width: 44,
                child: Text(
                  '${entry.rank}',
                  textAlign: TextAlign.center,
                  style: TextStyle(fontWeight: FontWeight.w700, color: medal),
                ),
              ),
              ClassEmblem(root: entry.classCode, size: 32),
              const SizedBox(width: Space.sm),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(entry.name, overflow: TextOverflow.ellipsis),
                    Text(
                      '${l.levelShort(entry.level)} · ${ContentNames.of('class', entry.classCode, className(l, entry.classCode))}',
                      style: const TextStyle(fontSize: 12, color: Palette.textMuted),
                    ),
                  ],
                ),
              ),
              Text('${entry.value}', style: const TextStyle(fontWeight: FontWeight.w700)),
            ],
          ),
        ),
      ),
    );
  }
}
