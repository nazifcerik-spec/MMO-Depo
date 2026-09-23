# Combat (Muharebe)

## 1. Amaç
Karakter ve düşman arasındaki mekanik etkileşimi modellemek. Tick-tabanlı, server-authoritative, hem manuel hem AFK'de aynı motor üstünde çalışır.

## 2. Oyuncuya Fayda
- Karar odaklı, refleks bağımlı olmayan, ancak yine de heyecan üreten muharebe.
- Aynı motor üzerinden AFK ve aktif oynanışı deneme.
- Party dungeon, dünya boss, hunting, dungeon → aynı temele oturur.

## 3. Gameplay Loop
- Manuel: skill tuşları, target pick, positioning (Phase 02 opsiyonel).
- AFK: rotation editor + aktivite tipi (hunting/dungeon farm).

## 4. Temel Kurallar
- **Tick rate**: server 250ms tick (Phase 01 baseline).
- **Global Cooldown**: 1000ms base, haste ile azalır (min 500ms).
- **Auto-attack**: her karakter default olarak weapon speed'ine göre auto atar.
- **Threat / Aggro** sistemi: tank pull mantığı; sistem basittir ancak taunt/skill etkilerini destekler.
- **Damage calc katmanları:**
  1. base ability damage
  2. + stat scaling
  3. × class/spec multipliers
  4. × situational (crit, backstab, elemental)
  5. − target armor / resist
  6. + on-hit procs
  7. clamp minimum 1

- **Healing**: ayrı katman; overheal excess olarak istatistiğe kayıt.
- **Death**: HP ≤ 0. Karakter graveyard'da respawn; XP debt uygulanır.

## 5. Combat Modes
- **Manual PvE** (kişisel muharebe)
- **AFK PvE** (activity engine yürütür)
- **Party PvE** (dungeon, hunting)
- **Raid** (Phase 02+, 10+ oyuncu, weekly reset)
- **World Boss** (bkz. [World Bosses](world-bosses.md))
- **PvP** (Phase 02+): duel, arena, open-world.

## 6. Kullanıcı Akışı (Manuel)
```
Target lock → GCD-aware skill press → server validate → damage/heal event → floating text → HP bar update
```

## 6b. Kullanıcı Akışı (AFK)
```
Aktivite başlat (hunting X mob type)
  → Server tick loop
    → target select (aggro rules)
    → rotation eval (priority list)
    → skill cast, damage calc, loot roll
  → activity yield log
  → reveal on login
```

## 7. İlişkiler
- [Skills](skills.md), [Stats](stats.md), [Items](items.md), [Equipment](equipment.md), [Pets](pets.md), [Dungeons](dungeons.md), [Hunting](hunting.md), [World Bosses](world-bosses.md), [Activity Engine](../systems/activity-engine.md), [Random System](../systems/random-system.md)

## 8. Gereken Veriler
```
combat_event
  id, timestamp, actor_id, target_id, ability_id
  damage_amount, damage_type, crit (bool), evaded (bool)
  overheal (int, nullable), resource_change

npc_template
  id, code, name, level, hp, damage, resistances, loot_table_id, ai_profile_id

npc_instance
  id, template_id, current_hp, aggro_list, spawned_at

status_effect
  id, source_id, target_id, effect_type, magnitude, expires_at, stacks

rotation
  character_id, priority_list (json)
```

## 9. I/O
- Input: skill use, target set, movement (Phase 02).
- Output: combat_event stream, damage number, kill event, loot roll.

## 10. Progression
- Combat power = f(level, stats, gear, skills, buff, spec).
- Kill rewards: XP + loot + faction rep (Phase 02) + activity yield.

## 11. Reward
Muharebe içi ödüller: kill xp, loot drop, quest kill sayacı.

## 12. Failure Durumları
- Ölüm: HP 0. Respawn (graveyard). XP debt (kaybedilen level'ın küçük yüzdesi, level 20+).
- Durability loss ölümde: -%10 durability.
- Party wipe: dungeon boss encounter reset (bkz. Dungeons).

## 13. Edge Cases
- Simultaneous damage that would kill: last hit source belirlenir; kill credit multi-hit rule'ları (highest damage + assist).
- Retaliation (mob reflects damage): DR clamp.
- Stun/Root/Silence CC: max control %30 rule (diminishing returns).

## 14. AFK Davranışı
AFK combat, oyuncunun rotation'ını tick tick yürütür. Tüm damage/heal/loot deterministik seed. Loot AFK yield'a eklenir; login'de reveal.

## 15. Server Authoritative Gereksinimleri
- Hasar hesabı server; istemci sadece animasyon.
- RNG seed'li (bkz. Random System).
- Anti-cheat: aşırı hızlı skill spam, resource inflation → sinyal.

## 16. Genişletilebilirlik
- Element sistemi (fire/ice/nature/holy/shadow) eklenmesine hazır.
- Combo point sistemi (rogue) modeli genişletilebilir.
- Weather/terrain modifier (Phase 02+).
- PvP damage formülü ayrı katman olarak eklenir (PvP vs PvE damage split).

## 17. Örnek Senaryo
Warrior/Guardian AFK hunting: mob spawn → aggro pull → auto-attack tick → Shield Slam (rotation cond) → mob dies → loot roll (seed + luck) → xp+ item add. 8 saatte 1400 mob kill.

## 18. Balance Considerations
- TTK (time to kill) target: normal mob 4–8 sn, elite 20–40 sn, boss 3–8 dk.
- AFK vs manuel: manuel oynayanın DPS'i ~%15–25 daha yüksek, ancak grind velocity AFK ile karşılaştırılabilir. Şu prensip: **aktif oynamak "daha eğlenceli", AFK oynamak "daha temiz süre bütçesi"**.
- Critical hit variance: baz %5, cap %75 diminishing; crit dmg baz ×1.5, cap ×2.5.
- Elemental resist band ±%30.
