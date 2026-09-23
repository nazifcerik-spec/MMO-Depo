# Guild Event Notification

## 1. Amaç
Guild aktivitelerine bağlı push: guild boss summon, raid, war, guild mail, MOTD change (opt), invite.

## 2. Trigger'lar

| Event | Kind | Notlar |
|---|---|---|
| Guild boss summoned | `guild_event` (subtype `boss_summoned`) | Officer summon, tüm guild |
| Guild raid formed | `guild_event` (subtype `raid_forming`) | Members raid create'de |
| Guild war declared | `guild_event` (subtype `war_declared`) | Tüm guild |
| Guild MOTD change | (v1'de push YOK) | Login'de otomatik göster |
| Guild invite | `guild_invite` | Target character (bkz. §5) |
| Guild leader change | `guild_event` (subtype `leader_change`) | Tüm guild |
| Guild disband initiated | `guild_event` (subtype `disband_grace`) | Tüm guild (48h grace) |

## 3. Payload Örnek (Boss Summoned)

```json
{
  "title_key": "notif.guild.boss_summoned.title",
  "body_key": "notif.guild.boss_summoned.body",
  "body_vars": {
    "boss_name": "Ashen Titan",
    "summoner_name": "Grim",
    "guild_name": "Iron Fist"
  },
  "category": "event",
  "route": "guild/<guild_id>/boss/<spawn_id>",
  "badge_delta": 1,
  "priority": "high",
  "collapse_key": "guild_boss_<guild_id>"
}
```

## 4. Rate Limit

Bkz. push-notification.md: `guild_event` 4/saat, cooldown 5dk.

Officer/leader daily action limit sistem seviyesinde de var (guild permission bank withdraw vs benzer) — spam abuse önlem.

## 5. Guild Invite

Ayrı kind `guild_invite`:
```json
{
  "title_key": "notif.guild.invite.title",
  "body_key": "notif.guild.invite.body",
  "body_vars": { "guild_name": "Iron Fist", "inviter_name": "Grim" },
  "category": "social",
  "route": "guild/invite/<invite_id>",
  "badge_delta": 1
}
```

Priority `normal`; ancak block check zorunlu (bkz. [social/block.md](../social/block.md)).

## 6. Cross-Reference
- [guild/boss.md](../guild/boss.md), [guild/raid.md](../guild/raid.md), [guild/war.md](../guild/war.md)
- [guild/members.md](../guild/members.md) — invite
- [social/block.md](../social/block.md)
- [preferences.md](preferences.md), [push-notification.md](push-notification.md)

## 7. Bekleyen Kararlar
- MOTD change push v1 yok; Phase 02 opt-in.
- Guild bank withdraw büyük miktar push (officer + leader) Phase 02.
- Alliance-wide event push (Phase 03).
- Rich guild-boss encounter progress push (mid-fight update) Phase 02.
