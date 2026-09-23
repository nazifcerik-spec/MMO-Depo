import 'dart:async';
import 'dart:convert';

import 'package:stream_channel/stream_channel.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// In-memory WebSocket for tests. [server] pushes frames to the client.
class FakeSocket with StreamChannelMixin<dynamic> implements WebSocketChannel {
  FakeSocket(this.uri, {this.failConnect = false});

  final Uri uri;
  final bool failConnect;
  final _in = StreamController<dynamic>();
  final _out = StreamController<dynamic>();
  bool closed = false;

  void push(String kind, [Map<String, dynamic> payload = const {}]) =>
      _in.add(jsonEncode({'kind': kind, 'payload': payload}));

  /// Server-side drop.
  void drop() => _in.close();

  @override
  Future<void> get ready => failConnect ? Future.error(WebSocketChannelException('refused')) : Future.value();

  @override
  Stream<dynamic> get stream => _in.stream;

  @override
  WebSocketSink get sink => _Sink(this);

  @override
  String? get protocol => null;

  @override
  int? get closeCode => null;

  @override
  String? get closeReason => null;
}

class _Sink implements WebSocketSink {
  _Sink(this.s);

  final FakeSocket s;

  @override
  void add(dynamic data) => s._out.add(data);

  @override
  void addError(Object error, [StackTrace? stackTrace]) {}

  @override
  Future<void> addStream(Stream<dynamic> stream) => stream.forEach(add);

  @override
  Future<void> close([int? closeCode, String? closeReason]) async {
    s.closed = true;
  }

  @override
  Future<void> get done => Future.value();
}

/// Records every socket the client opens.
class FakeSocketServer {
  final sockets = <FakeSocket>[];
  int failNext = 0;

  WebSocketChannel connect(Uri uri) {
    final s = FakeSocket(uri, failConnect: failNext > 0);
    if (failNext > 0) failNext--;
    sockets.add(s);
    return s;
  }

  FakeSocket get last => sockets.last;
}
