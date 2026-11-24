/// What only a browser tab has (DD-369): the lock its sibling tabs share and
/// storage that dies with the tab. Elsewhere the lock is free and the storage
/// empty.
library;

export 'browser_stub.dart' if (dart.library.js_interop) 'browser_web.dart';
