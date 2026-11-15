import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart';
import '../../shared/motion.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../quests/quest_repository.dart';
import '../quests/quests_screen.dart';
import 'story_repository.dart';

/// Story mode (DD-313, DD-367): the chapters as a journal. Passed scenes fold
/// up with the answer given; the current scene shows its text, its story
/// quests (accepted and turned in right here) and its choice.
class StoryScreen extends StatefulWidget {
  const StoryScreen({super.key, required this.repo, required this.quests, required this.dataVersion});

  final StoryRepository repo;
  final QuestRepository quests;
  final ValueNotifier<int> dataVersion;

  @override
  State<StoryScreen> createState() => _StoryScreenState();
}

class _StoryScreenState extends State<StoryScreen> with PanelState {
  StoryJournal? _journal;
  Map<String, String> _text = const {};
  QuestBoard? _board;
  String _lang = 'en';
  bool _started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _lang = Localizations.localeOf(context).languageCode;
    if (!_started) {
      _started = true;
      load();
    }
  }

  @override
  Future<void> fetch() async {
    // Future.wait rethrows the ApiError itself (PanelState catches it).
    final r = await Future.wait<Object>([widget.repo.journal(), widget.quests.board()]);
    final j = r[0] as StoryJournal, b = r[1] as QuestBoard;
    _text = await widget.repo.text(_lang, j.textVersion);
    _journal = j;
    _board = b;
  }

  String _t(String key) => _text[key] ?? '';

  String _optionLabel(String scene, String option) => _t('$scene.$option');

  Future<void> _quest(Future<void> Function() action) => run(() async {
    await action();
    await load();
  });

  Future<void> _turnIn(Quest q) => _quest(() async {
    final r = await widget.quests.turnIn(q.code);
    widget.dataVersion.value++;
    if (!mounted) return;
    final l = AppLocalizations.of(context);
    final parts = [if (r.xp > 0) l.xpGain(r.xp), if (r.gold > 0) '+${goldString(r.gold)}'];
    toast(
      context,
      Text([l.questRewarded(parts.join(' · ')), if (r.newLevel != null) l.levelUp(r.newLevel!)].join('\n')),
    );
  });

  Future<void> _advance(StoryScene sc, StoryOption? o) async {
    final l = AppLocalizations.of(context);
    if (o != null) {
      final ok = await confirmAction(
        context,
        title: l.storyChoiceTitle,
        body: '${_optionLabel(sc.code, o.code)}\n\n${l.storyChoiceConfirm}',
        action: l.storyChoose,
      );
      if (!ok || !mounted) return;
    }
    await run(() async {
      await widget.repo.advance(sc.code, o?.code ?? '');
      await load(); // the next scene offers its own quests on the board
      final result = o == null ? '' : _t('${sc.code}.${o.code}.result');
      if (result.isNotEmpty && mounted) await _showResult(_t(sc.code), result);
    });
  }

  Future<void> _showResult(String title, String body) => showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (c) => SafeArea(
      child: SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
        child: PopIn(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(title, style: Theme.of(c).textTheme.titleLarge),
              const SizedBox(height: Space.md),
              StoryText(body),
              const SizedBox(height: Space.lg),
              FilledButton(onPressed: () => Navigator.pop(c), child: Text(AppLocalizations.of(c).storyContinue)),
            ],
          ),
        ),
      ),
    ),
  );

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Scaffold(
      appBar: AppBar(title: Text(l.storyTitle)),
      body: SafeArea(
        child: PageFrame(
          child: frame(() {
            final j = _journal!;
            if (j.chapters.isEmpty) {
              return [MessageView(message: l.storyComingSoon, icon: Icons.auto_stories_outlined)];
            }
            return [
              for (final c in j.chapters)
                _chapter(l, c, [
                  for (final s in j.scenes)
                    if (s.chapter == c.code) s,
                ]),
            ];
          }),
        ),
      ),
    );
  }

  Widget _chapter(AppLocalizations l, StoryChapterState c, List<StoryScene> scenes) {
    final text = Theme.of(context).textTheme;
    final title = ContentNames.of('story', c.code);
    final blurb = ContentNames.of('story', '${c.code}.body', '');
    if (c.state == 'locked') {
      return Card(
        margin: const EdgeInsets.only(bottom: Space.md),
        child: ListTile(
          enabled: false,
          leading: const Icon(Icons.lock_outline),
          title: Text(title),
          subtitle: Text(l.storyLocked(c.minLevel)),
        ),
      );
    }
    final done = c.state == 'done';
    return Card(
      margin: const EdgeInsets.only(bottom: Space.md),
      clipBehavior: Clip.antiAlias,
      child: ExpansionTile(
        initiallyExpanded: !done,
        leading: Icon(done ? Icons.task_alt : Icons.auto_stories, color: done ? Palette.success : Palette.accent),
        title: Text(title, style: text.titleMedium),
        subtitle: blurb.isEmpty ? null : Text(blurb, style: text.bodySmall),
        childrenPadding: const EdgeInsets.fromLTRB(Space.md, 0, Space.md, Space.md),
        children: [for (final s in scenes) s.done ? _doneScene(l, s) : _currentScene(l, s)],
      ),
    );
  }

  List<Widget> _paragraphs(StoryScene s, {required bool withAfter}) => [
    StoryText(_t('${s.code}.body')),
    for (final x in s.extras) StoryText(_t('${s.code}.x.$x')),
    if (withAfter && s.hasAfter) StoryText(_t('${s.code}.after')),
    if (withAfter)
      for (final x in s.afterExtras) StoryText(_t('${s.code}.x.$x')),
  ];

  Widget _doneScene(AppLocalizations l, StoryScene s) {
    final text = Theme.of(context).textTheme;
    final result = s.chosen.isEmpty ? '' : _t('${s.code}.${s.chosen}.result');
    return Card(
      color: Palette.bgElevated,
      margin: const EdgeInsets.only(top: Space.sm),
      clipBehavior: Clip.antiAlias,
      child: ExpansionTile(
        dense: true,
        title: Text(_t(s.code), style: text.titleSmall),
        subtitle: s.chosen.isEmpty
            ? null
            : Text(l.storyYourChoice(_optionLabel(s.code, s.chosen)), style: text.bodySmall),
        childrenPadding: const EdgeInsets.fromLTRB(Space.md, 0, Space.md, Space.md),
        expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          ..._paragraphs(s, withAfter: true),
          if (result.isNotEmpty) ...[const Divider(), StoryText(result)],
        ],
      ),
    );
  }

  Widget _currentScene(AppLocalizations l, StoryScene s) {
    final text = Theme.of(context).textTheme;
    final b = _board!;
    Quest? boardQuest(String code) {
      for (final q in [...b.open, ...b.available]) {
        if (q.code == code) return q;
      }
      return null;
    }

    final locked = s.state == 'locked';
    return Padding(
      padding: const EdgeInsets.only(top: Space.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(_t(s.code), style: text.titleLarge?.copyWith(color: Palette.accent)),
          const SizedBox(height: Space.sm),
          ..._paragraphs(s, withAfter: s.ready && !locked),
          if (locked)
            Padding(
              padding: const EdgeInsets.symmetric(vertical: Space.md),
              child: Row(
                children: [
                  const Icon(Icons.hourglass_bottom, size: 18, color: Palette.textMuted),
                  const SizedBox(width: Space.sm),
                  Expanded(child: Text(l.storyLocked(s.minLevel), style: text.bodyMedium)),
                ],
              ),
            ),
          if (!locked && s.quests.isNotEmpty) ...[
            Padding(
              padding: const EdgeInsets.fromLTRB(0, Space.md, 0, Space.sm),
              child: Text(l.storyQuests, style: text.titleSmall?.copyWith(color: Palette.accent)),
            ),
            for (final sq in s.quests) _sceneQuest(l, sq, boardQuest(sq.code)),
          ],
          if (!locked && s.ready) ...[
            const SizedBox(height: Space.md),
            if (s.options.isEmpty)
              FilledButton(onPressed: busy ? null : () => _advance(s, null), child: Text(l.storyContinue))
            else
              for (final o in s.options)
                Padding(
                  padding: const EdgeInsets.only(bottom: Space.sm),
                  child: OutlinedButton(
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.all(Space.md),
                      alignment: Alignment.centerLeft,
                    ),
                    onPressed: busy || o.requiresQuest != null ? null : () => _advance(s, o),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(_optionLabel(s.code, o.code)),
                        if (o.requiresQuest != null)
                          Text(l.storyNeedsQuest(questName(o.requiresQuest!)), style: text.bodySmall),
                      ],
                    ),
                  ),
                ),
          ],
        ],
      ),
    );
  }

  Widget _sceneQuest(AppLocalizations l, StoryQuest sq, Quest? q) {
    if (sq.state == 'turned_in' || q == null) {
      return ListTile(
        dense: true,
        contentPadding: EdgeInsets.zero,
        leading: Icon(
          sq.state == 'turned_in' ? Icons.task_alt : Icons.assignment_outlined,
          color: sq.state == 'turned_in' ? Palette.success : Palette.textMuted,
        ),
        title: Text(questName(sq.code)),
        subtitle: sq.state == 'turned_in' ? Text(l.storyQuestDone) : null,
      );
    }
    final available = q.state == 'available';
    return QuestCard(
      key: ValueKey('story-${q.code}-${q.state}'),
      q: q,
      primaryLabel: available ? l.accept : l.turnIn,
      onPrimary: busy
          ? null
          : available
          ? () => _quest(() => widget.quests.accept(q.code))
          : () => _turnIn(q),
    );
  }
}

/// Story text: blank lines split paragraphs, **x** is bold (the speaker),
/// *x* italic (the inner voice).
class StoryText extends StatelessWidget {
  const StoryText(this.text, {super.key});

  final String text;

  static final _span = RegExp(r'\*\*([^*]+)\*\*|\*([^*]+)\*');

  static List<TextSpan> spans(String s) {
    final out = <TextSpan>[];
    var at = 0;
    for (final m in _span.allMatches(s)) {
      if (m.start > at) out.add(TextSpan(text: s.substring(at, m.start)));
      out.add(
        m.group(1) != null
            ? TextSpan(
                text: m.group(1),
                style: const TextStyle(fontWeight: FontWeight.w700, color: Palette.accent),
              )
            : TextSpan(
                text: m.group(2),
                style: const TextStyle(fontStyle: FontStyle.italic, color: Palette.info),
              ),
      );
      at = m.end;
    }
    if (at < s.length) out.add(TextSpan(text: s.substring(at)));
    return out;
  }

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(context).textTheme.bodyLarge?.copyWith(height: 1.45);
    final paragraphs = text.trim().split(RegExp(r'\n\s*\n'));
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final p in paragraphs)
          if (p.trim().isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(bottom: Space.md),
              child: Text.rich(TextSpan(style: style, children: spans(p.trim()))),
            ),
      ],
    );
  }
}
