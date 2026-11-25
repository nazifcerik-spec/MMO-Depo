import 'dart:js_interop';
import 'dart:js_interop_unsafe';

import 'package:web/web.dart' as web;

/// Runs [action] while no other tab runs one under [name].
Future<void> browserLock(String name, Future<void> Function() action) async {
  // Web Locks need a secure context (https or localhost).
  if (!(web.window.navigator as JSObject).has('locks')) return action();
  Object? error;
  StackTrace? trace;
  await web.window.navigator.locks
      .request(
        name,
        ((JSAny? _) {
          return action()
              .then<JSAny?>(
                (_) => null,
                onError: (Object e, StackTrace s) {
                  error = e;
                  trace = s;
                  return null;
                },
              )
              .toJS;
        }).toJS,
      )
      .toDart;
  if (error case final Object e) Error.throwWithStackTrace(e, trace!);
}

/// A value kept for this tab only; null removes it.
void tabStore(String key, String? value) {
  try {
    if (value == null) {
      web.window.sessionStorage.removeItem(key);
    } else {
      web.window.sessionStorage.setItem(key, value);
    }
  } catch (_) {
    // Storage blocked: a reload lands on the character list.
  }
}

String? tabRead(String key) {
  try {
    return web.window.sessionStorage.getItem(key);
  } catch (_) {
    return null;
  }
}

/// The web address of a player's page; null outside a browser.
String? profileUrl(String name) => '${web.document.baseURI}@${Uri.encodeComponent(name)}';

const _linkKey = 'mena.profile';

/// The player named by the address the page was opened with (`…/@Name`),
/// once; null without one. The address goes back to the game's, and the name
/// waits in the tab's storage while the player signs in.
String? takeProfileLink() {
  final last = Uri.tryParse(web.window.location.href)?.pathSegments.lastOrNull ?? '';
  if (last.length > 1 && last.startsWith('@')) {
    tabStore(_linkKey, last.substring(1));
    web.window.history.replaceState(null, '', web.document.baseURI);
  }
  final name = tabRead(_linkKey);
  tabStore(_linkKey, null);
  return name;
}
