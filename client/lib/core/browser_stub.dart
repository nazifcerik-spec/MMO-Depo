/// Runs [action] while no other tab runs one under [name].
Future<void> browserLock(String name, Future<void> Function() action) => action();

/// A value kept for this tab only; null removes it.
void tabStore(String key, String? value) {}

String? tabRead(String key) => null;

/// The web address of a player's page; null outside a browser.
String? profileUrl(String name) => null;

/// The player named by the address the page was opened with (`…/@Name`),
/// once; null without one.
String? takeProfileLink() => null;
