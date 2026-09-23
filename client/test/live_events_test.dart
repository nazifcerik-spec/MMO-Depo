import 'package:flutter_test/flutter_test.dart';
import 'package:otter/core/live_events.dart';

import 'fake_backend.dart';
import 'fake_ws.dart';

void main() {
  late FakeBackend be;
  late FakeSocketServer server;

  setUp(() {
    be = FakeBackend();
    server = FakeSocketServer();
  });

  Future<LiveEvents> started() async {
    final api = be.api();
    await api.login('grim', 'correct-horse');
    final live = LiveEvents(api, connect: server.connect, maxBackoff: const Duration(milliseconds: 20));
    live.start();
    await pumpEventQueue();
    return live;
  }

  test('delivers frames and skips malformed ones', () async {
    final live = await started();
    final got = <String>[];
    live.events.listen((e) => got.add(e.kind));

    server.last.push('system.welcome');
    server.last.push('hunt.reveal', {'contract_id': 'h1'});
    server.last.drop(); // closes after queued frames
    await pumpEventQueue();

    expect(got, ['system.welcome', 'hunt.reveal']);
    live.stop();
  });

  test('reconnects with backoff after a drop or failed handshake, then stops cleanly', () async {
    server.failNext = 1;
    final live = await started();
    expect(server.sockets, hasLength(1));
    expect(live.connected, isFalse);

    await Future<void>.delayed(const Duration(milliseconds: 60));
    expect(server.sockets, hasLength(2), reason: 'retried after failed handshake');
    expect(live.connected, isTrue);

    server.last.drop();
    await Future<void>.delayed(const Duration(milliseconds: 60));
    expect(server.sockets, hasLength(3), reason: 'retried after server drop');

    live.stop();
    server.last.drop();
    await Future<void>.delayed(const Duration(milliseconds: 60));
    expect(server.sockets, hasLength(3), reason: 'no reconnect once stopped');
    expect(live.connected, isFalse);
  });

  test('does not connect without a session', () async {
    final live = LiveEvents(be.api(), connect: server.connect);
    live.start();
    await pumpEventQueue();
    expect(server.sockets, isEmpty);
    live.stop();
  });
}
