import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:otter/core/token_store.dart';
import 'package:otter/features/inventory/inventory_repository.dart';
import 'package:otter/features/social/group_repository.dart';
import 'package:otter/features/social/guild_vault_screen.dart';
import 'package:otter/l10n/app_localizations.dart';

import 'fake_backend.dart';

void main() {
  testWidgets('deposit gold to the treasury; the log and limits show (DD-358)', (tester) async {
    final be = FakeBackend();
    final api = be.api(MemoryTokenStore());
    await api.login('grim', 'correct-horse');
    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('en'),
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: AppLocalizations.supportedLocales,
        home: GuildVaultScreen(repo: GroupRepository(api, 'char-1'), inventory: InventoryRepository(api, 'char-1')),
      ),
    );
    await tester.pumpAndSettle();
    expect(find.text('Guild treasury'), findsOneWidget);
    expect(find.textContaining('Withdrawn today:'), findsOneWidget);
    expect(find.text('Nothing has moved yet.'), findsOneWidget);

    await tester.tap(find.byKey(const Key('vault-gold-in')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const Key('vault-gold-amount')), '12');
    await tester.tap(find.byKey(const Key('vault-gold-confirm')));
    await tester.pumpAndSettle();
    expect(be.vaultCopper, 120000, reason: 'whole gold, sent in copper');
    expect(find.textContaining('Grim deposited 12'), findsOneWidget);

    // More than the treasury holds is refused in the player's words.
    await tester.tap(find.byKey(const Key('vault-gold-out')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const Key('vault-gold-amount')), '13');
    await tester.tap(find.byKey(const Key('vault-gold-confirm')));
    await tester.pumpAndSettle();
    expect(be.vaultCopper, 120000);
  });
}
