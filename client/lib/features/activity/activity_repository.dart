import '../../core/api_client.dart';
import '../../shared/party_share.dart';
import 'activity_models.dart';

class ActivityRepository {
  const ActivityRepository(this.api);

  final ApiClient api;

  /// Null when the character has no running contract.
  Future<ActiveActivity?> active(String characterId) async {
    final j = await api.get('/v1/character/$characterId/activity/active');
    final c = j['contract'];
    if (c is! Map<String, dynamic>) return null;
    return ActiveActivity(
      ActivityContract.fromJson(c),
      ActivityYield.fromJson(j['current_yield'] as Map<String, dynamic>?),
    );
  }

  /// Every gathering profession's level and resources (DD-244).
  Future<List<GatherSkill>> skills(String characterId) async {
    final j = await api.get('/v1/character/$characterId/gathering');
    return [for (final p in (j['professions'] as List? ?? const [])) GatherSkill.fromJson(p as Map<String, dynamic>)];
  }

  Future<ActiveActivity> start(String characterId, String kind, Duration duration, {String? target}) async {
    final j = await api.post('/v1/character/$characterId/activity', {
      'kind': kind,
      'duration_seconds': duration.inSeconds,
      'target': ?target,
    });
    return ActiveActivity(
      ActivityContract.fromJson(j['contract'] as Map<String, dynamic>),
      const ActivityYield(completionRatio: 0, goldCopper: 0, common: 0, rare: 0, legendary: 0),
    );
  }

  /// Collects what was earned so far and ends the contract. Idempotent on the
  /// server (a retried reveal returns the cached result).
  Future<(ActivityYield, GatherProgress?, PartyShare?)> reveal(String contractId) async {
    final j = await api.post('/v1/activity/$contractId/reveal', const {});
    return (
      ActivityYield.fromJson(j['yield'] as Map<String, dynamic>?),
      GatherProgress.fromJson(j['gather']),
      PartyShare.fromJson(j['party_share']), // DD-296
    );
  }

  /// Abandons the contract; accrued yield is forfeited.
  Future<void> stop(String contractId) => api.post('/v1/activity/$contractId/stop', const {});
}
