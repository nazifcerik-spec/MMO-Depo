import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';

/// A story chapter (DD-313); its title is the content name story.<code>,
/// its text story.<code>.body.
class StoryChapter {
  const StoryChapter({required this.code, required this.minLevel, this.requires = ''});

  factory StoryChapter.fromJson(Map<String, dynamic> j) => StoryChapter(
    code: j['code'] as String? ?? '',
    minLevel: (j['min_level'] as num?)?.toInt() ?? 1,
    requires: j['requires'] as String? ?? '',
  );

  final String code;
  final int minLevel;
  final String requires;

  String get title => ContentNames.of('story', code);
  String get body => ContentNames.of('story', '$code.body', '');
}

/// Story mode: the chapter list, locked by level; the content comes later
/// and until then the page says so.
class StoryScreen extends StatefulWidget {
  const StoryScreen({super.key, required this.api, required this.level});

  final ApiClient api;
  final int level;

  @override
  State<StoryScreen> createState() => _StoryScreenState();
}

class _StoryScreenState extends State<StoryScreen> {
  List<StoryChapter>? _chapters;
  Object? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _error = null);
    try {
      final j = await widget.api.get('/v1/story');
      final list = [
        for (final c in (j['chapters'] as List? ?? const []).cast<Map<String, dynamic>>()) StoryChapter.fromJson(c),
      ];
      if (mounted) setState(() => _chapters = list);
    } on ApiError catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final chapters = _chapters;
    return Scaffold(
      appBar: AppBar(title: Text(l.storyTitle)),
      body: SafeArea(
        child: PageFrame(
          child: switch ((chapters, _error)) {
            (_, final Object e) => MessageView(message: errorText(l, e), onRetry: _load),
            (null, _) => const Center(child: CircularProgressIndicator()),
            (final List<StoryChapter> list, _) when list.isEmpty => MessageView(
              message: l.storyComingSoon,
              icon: Icons.auto_stories_outlined,
            ),
            (final List<StoryChapter> list, _) => ListView(
              padding: const EdgeInsets.all(Space.screen),
              children: [
                for (final (i, c) in list.indexed)
                  Card(
                    child: ListTile(
                      leading: CircleAvatar(child: Text('${i + 1}')),
                      title: Text(c.title),
                      subtitle: c.minLevel > widget.level ? Text(l.storyLocked(c.minLevel)) : null,
                      enabled: c.minLevel <= widget.level,
                      trailing: const Icon(Icons.chevron_right),
                      onTap: () => Navigator.of(context).push(
                        MaterialPageRoute<void>(
                          builder: (_) => Scaffold(
                            appBar: AppBar(title: Text(c.title)),
                            body: SafeArea(
                              child: SingleChildScrollView(
                                padding: const EdgeInsets.all(Space.screen),
                                child: Text(c.body, style: Theme.of(context).textTheme.bodyLarge),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          },
        ),
      ),
    );
  }
}
