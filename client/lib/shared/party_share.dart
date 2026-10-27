import '../l10n/app_localizations.dart';

/// How a finished job's XP was split with party members doing the same work
/// (DD-296): the job's XP, what the earner kept, who shared it.
class PartyShare {
  const PartyShare({required this.earned, required this.kept, required this.partners});

  static PartyShare? fromJson(Object? j) => j is Map<String, dynamic>
      ? PartyShare(
          earned: (j['earned'] as num?)?.toInt() ?? 0,
          kept: (j['kept'] as num?)?.toInt() ?? 0,
          partners: [
            for (final p in (j['partners'] as List? ?? const []).whereType<Map<String, dynamic>>())
              p['name'] as String? ?? '',
          ],
        )
      : null;

  final int earned;
  final int kept;
  final List<String> partners;

  String text(AppLocalizations l) => l.partyShared(partners.join(', '), kept, earned);
}
