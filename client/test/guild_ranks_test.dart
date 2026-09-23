import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/social/group_repository.dart';

// Mirrors server/internal/guild/guild.go Can* — the client only hides buttons,
// the server enforces; a drift would show buttons that always fail.
void main() {
  const leader = GuildRanks.leader, officer = GuildRanks.officer, veteran = GuildRanks.veteran;
  const member = GuildRanks.member, initiate = GuildRanks.initiate;

  test('invite and motd', () {
    expect([for (var r = 0; r <= 4; r++) GuildRanks.canInvite(r)], [true, true, true, false, false]);
    expect([for (var r = 0; r <= 4; r++) GuildRanks.canSetMotd(r)], [true, true, false, false, false]);
  });

  test('kick only strictly lower ranks, officers and up', () {
    expect(GuildRanks.canKick(leader, officer), isTrue);
    expect(GuildRanks.canKick(officer, officer), isFalse);
    expect(GuildRanks.canKick(officer, veteran), isTrue);
    expect(GuildRanks.canKick(veteran, initiate), isFalse);
  });

  test('set rank', () {
    expect(GuildRanks.canSetRank(leader, member, officer), isTrue);
    expect(GuildRanks.canSetRank(leader, member, leader), isFalse, reason: 'leadership moves only via transfer');
    expect(GuildRanks.canSetRank(officer, member, officer), isFalse, reason: 'officers cannot mint officers');
    expect(GuildRanks.canSetRank(officer, member, veteran), isTrue);
    expect(GuildRanks.canSetRank(officer, officer, veteran), isFalse, reason: 'no peer changes');
    expect(GuildRanks.canSetRank(veteran, initiate, member), isFalse);
  });
}
