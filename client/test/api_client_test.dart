import 'dart:math';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:otter/core/api_client.dart';
import 'package:otter/core/token_store.dart';

import 'fake_backend.dart';

void main() {
  late FakeBackend be;
  late MemoryTokenStore store;
  late ApiClient api;

  setUp(() {
    be = FakeBackend();
    store = MemoryTokenStore();
    api = be.api(store);
  });

  test('login persists only the refresh token and signs later calls', () async {
    await api.login('grim', 'correct-horse');
    expect(store.token, 'refresh-1');
    expect(api.accountId, 'acct-1');

    await api.get('/v1/character');
    final req = be.requests.last;
    expect(req.headers['Authorization'], 'Bearer access-1');
    expect(req.headers['X-Request-Id'], matches(RegExp(r'^[0-9a-f-]{36}$')));
    expect(req.headers.containsKey('Idempotency-Key'), isFalse, reason: 'GETs carry no idempotency key');
  });

  test('every POST carries a valid Idempotency-Key', () async {
    await api.login('grim', 'correct-horse');
    await api.post('/v1/character', {'name': 'Grim', 'race_code': 'human', 'class_code': 'warrior'});
    final key = be.requests.last.headers['Idempotency-Key']!;
    expect(key, matches(RegExp(r'^[A-Za-z0-9_.\-]{16,128}$')));
  });

  test('wrong password surfaces the server error code', () async {
    await expectLater(
      api.login('grim', 'nope'),
      throwsA(
        isA<ApiError>().having((e) => e.code, 'code', 'E_INVALID_CREDENTIALS').having((e) => e.status, 'status', 401),
      ),
    );
    expect(api.hasSession, isFalse);
  });

  test('expired access token refreshes once and replays with the same Idempotency-Key', () async {
    await api.login('grim', 'correct-horse');
    be.validAccess.clear(); // server-side expiry

    await api.post('/v1/character', {'name': 'Grim', 'race_code': 'human', 'class_code': 'warrior'});

    expect(be.refreshCalls, 1);
    final posts = be.requests.where((r) => r.url.path == '/v1/character').toList();
    expect(posts, hasLength(2));
    expect(posts[0].headers['Idempotency-Key'], posts[1].headers['Idempotency-Key']);
    expect(posts[1].headers['Authorization'], 'Bearer access-2');
    expect(store.token, 'refresh-2', reason: 'rotated refresh token persisted');
  });

  test('concurrent 401s share a single refresh (no reuse detection trip)', () async {
    await api.login('grim', 'correct-horse');
    be.validAccess.clear();

    await Future.wait([api.get('/v1/character'), api.get('/v1/character'), api.get('/v1/character')]);
    expect(be.refreshCalls, 1);
  });

  test('rejected refresh clears the session and notifies', () async {
    var lost = 0;
    api.onSessionLost = () => lost++;
    await api.login('grim', 'correct-horse');
    be.validAccess.clear();
    be.rejectRefresh = true;

    await expectLater(api.get('/v1/character'), throwsA(isA<ApiError>()));
    expect(lost, 1);
    expect(api.hasSession, isFalse);
    expect(store.token, isNull);
  });

  test('GET retries transient 503s; POST never auto-retries', () async {
    await api.login('grim', 'correct-horse');
    be.failNextGets503 = 2;
    final res = await api.get('/v1/character');
    expect(res['slot_cap'], 2);

    var calls = 0;
    final flaky = ApiClient(
      baseUrl: 'http://otter.test',
      store: MemoryTokenStore(),
      retryDelay: Duration.zero,
      client: MockClient((_) async {
        calls++;
        return http.Response('{"error":{"code":"E_UNAVAILABLE"}}', 503);
      }),
    );
    await expectLater(
      flaky.post('/v1/auth/login', {}, auth: false),
      throwsA(isA<ApiError>().having((e) => e.status, 'status', 503)),
    );
    expect(calls, 1);
  });

  test('transport failure maps to E_NETWORK after retries', () async {
    var calls = 0;
    final offline = ApiClient(
      baseUrl: 'http://otter.test',
      store: MemoryTokenStore(),
      retryDelay: Duration.zero,
      client: MockClient((_) async {
        calls++;
        throw http.ClientException('down');
      }),
    );
    await expectLater(
      offline.get('/v1/world/regions'),
      throwsA(isA<ApiError>().having((e) => e.isNetwork, 'net', true)),
    );
    expect(calls, 3);
  });

  test('restore: no token → false; valid token → session; offline → rethrows', () async {
    expect(await api.restore(), isFalse);

    await api.login('grim', 'correct-horse');
    final cold = be.api(store); // app restart with the persisted token
    expect(await cold.restore(), isTrue);
    expect(cold.hasSession, isTrue);

    final offline = ApiClient(
      baseUrl: 'http://otter.test',
      store: store,
      retryDelay: Duration.zero,
      client: MockClient((_) async => throw http.ClientException('down')),
    );
    await expectLater(offline.restore(), throwsA(isA<ApiError>()));
    expect(store.token, isNotNull, reason: 'offline must not log the player out');
  });

  test('logout clears local state even when the server is unreachable', () async {
    await api.login('grim', 'correct-horse');
    api.characterId = 'char-1';
    await api.logout();
    expect(api.hasSession, isFalse);
    expect(api.characterId, isNull);
    expect(store.token, isNull);
  });

  test('uuidV4 is RFC 4122 v4', () {
    final id = uuidV4(Random(1));
    expect(id, matches(RegExp(r'^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$')));
  });
}
