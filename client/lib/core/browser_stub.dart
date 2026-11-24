/// Runs [action] while no other tab runs one under [name].
Future<void> browserLock(String name, Future<void> Function() action) => action();

/// A value kept for this tab only; null removes it.
void tabStore(String key, String? value) {}

String? tabRead(String key) => null;
