import '../l10n/app_localizations.dart';

/// The party bonus on a finished job's XP (DD-305): the XP before and after
/// it, and the members who worked alongside.
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
