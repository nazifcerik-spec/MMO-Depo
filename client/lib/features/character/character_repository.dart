import '../../core/api_client.dart';
import 'models.dart';

/// Character endpoints. Stateless; screens own their loading state.
class CharacterRepository {
  const CharacterRepository(this.api);

  final ApiClient api;

  Future<CharacterList> list() async => CharacterList.fromJson(await api.get('/v1/character'));

  Future<CharacterSummary> create(String name, String race, String cls) async =>
      CharacterSummary.fromJson(await api.post('/v1/character', {'name': name, 'race_code': race, 'class_code': cls}));

  /// Marks the character as played and scopes later calls to it.
  Future<CharacterSheet> enter(String id) async {
    final body = await api.post('/v1/character/$id/enter', const {});
    api.characterId = id;
    return CharacterSheet.fromJson(body['sheet'] as Map<String, dynamic>);
  }

  Future<CharacterSheet> sheet(String id) async => CharacterSheet.fromJson(await api.get('/v1/character/$id/sheet'));

  Future<int> goldCopper(String id) async {
    final w = await api.get('/v1/character/$id/wallet');
    return (w['gold_copper'] as num?)?.toInt() ?? 0;
  }
}
