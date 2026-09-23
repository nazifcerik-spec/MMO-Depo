import 'dart:async';
import 'dart:convert';
import 'dart:math';

import 'package:web_socket_channel/web_socket_channel.dart';

import 'api_client.dart';

/// One server push frame: `{"kind": "hunt.reveal", "payload": {...}}`.
class LiveEvent {
  const LiveEvent(this.kind, this.payload);

  final String kind;
  final Map<String, dynamic> payload;
}

typedef ChannelFactory = WebSocketChannel Function(Uri uri);

/// Single multiplexed WebSocket (docs/architecture/flutter.md §6) with
/// exponential backoff + jitter reconnect. The server pings every 25s; the
/// platform answers pongs, so no client heartbeat timer is needed.
///
/// Call [start] when a session is active and the app is foregrounded,
/// [stop] on background/sign-out — nothing runs while stopped.
class LiveEvents {
  LiveEvents(this.api, {ChannelFactory? connect, this.maxBackoff = const Duration(seconds: 30)})
    : _connect = connect ?? WebSocketChannel.connect;

  final ApiClient api;
  final Duration maxBackoff;
  final ChannelFactory _connect;
  final _events = StreamController<LiveEvent>.broadcast();
  final _rand = Random();

  WebSocketChannel? _channel;
  StreamSubscription<dynamic>? _sub;
  Timer? _retry;
  int _failures = 0;
  bool _running = false;

  Stream<LiveEvent> get events => _events.stream;
  bool get connected => _channel != null;

  void start() {
    if (_running) return;
    _running = true;
    _open();
  }

  void stop() {
    _running = false;
    _retry?.cancel();
    _retry = null;
    _close();
  }

  Future<void> dispose() async {
    stop();
    await _events.close();
  }

  Uri get _uri {
    final base = Uri.parse(api.baseUrl);
    return base.replace(scheme: base.scheme == 'https' ? 'wss' : 'ws', path: '/v1/ws');
  }

  Future<void> _open() async {
    if (!_running || _channel != null) return;
    try {
      final token = await api.freshAccessToken();
      if (!_running) return;
      final ch = _connect(_uri.replace(queryParameters: {'token': token}));
      _channel = ch;
      await ch.ready;
      _sub = ch.stream.listen(_onFrame, onDone: _onDrop, onError: (_) => _onDrop(), cancelOnError: true);
    } on ApiError catch (e) {
      if (!e.isNetwork) {
        stop(); // session gone — AuthController handles the sign-out
        return;
      }
      _onDrop();
    } catch (_) {
      _onDrop();
    }
  }

  void _onFrame(dynamic data) {
    if (data is! String) return;
    try {
      final j = jsonDecode(data);
      if (j is! Map<String, dynamic>) return;
      final kind = j['kind'];
      if (kind is! String) return;
      if (kind == 'system.welcome') _failures = 0;
      final payload = j['payload'];
      _events.add(LiveEvent(kind, payload is Map<String, dynamic> ? payload : const {}));
    } on FormatException {
      // Malformed frame — ignore; the server is authoritative, pushes are hints.
    }
  }

  void _onDrop() {
    _close();
    if (!_running) return;
    _failures++;
    final capped = min(maxBackoff.inMilliseconds, 1000 * (1 << min(_failures - 1, 5)));
    final delay = Duration(milliseconds: capped ~/ 2 + _rand.nextInt(capped ~/ 2 + 1));
    _retry?.cancel();
    _retry = Timer(delay, _open);
  }

  void _close() {
    _sub?.cancel();
    _sub = null;
    _channel?.sink.close();
    _channel = null;
  }
}
