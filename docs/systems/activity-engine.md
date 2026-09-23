# Activity Engine (Aktivite Motoru)

## 1. Amaç
Tüm oyun içi aktivitelerin (mining, hunting, travel, adventure, dungeon idle, crafting) tek bir soyutlama altında yürütülmesini sağlamak. Yeni aktivite tipini eklemek sistemi çatlatmadan mümkün olmalı.

## 2. Oyuncuya Fayda
- Aktiviteler arası tutarlı UX (start/stop, cancel, reveal).
- Yeni aktivite tipi eklendiğinde tanıdık davranış.

## 3. Kavramlar
- **Activity Type**: enum + spec (params şeması).
- **Activity Instance**: bir karakterin başlattığı, running/paused/completed state.
- **Life-cycle**: `start → tick → interrupt / pause / complete → reveal on login`.
- **Contract vs Instance**: contract kalıcıdır (AFK yaşam boyu); instance session-scope olabilir (aktif dungeon run).

## 4. Life-cycle
```
start(char, activity_type, params)
  → validate (character eligibility, location, requirements)
  → create instance + contract (if AFK-eligible)
  → persist
  → emit "activity_started"
tick()
  → per activity type: simulate delta
  → check stop conditions
  → emit events
pause(reason)
  → set status paused; keep contract
resume()
  → set status active; new start_at
interrupt(reason)
  → cancel; partial yield settlement
complete()
  → final yield; contract closed
```

## 5. Kurallar
- Bir karakter aynı anda **1 ana aktivite**; ek olarak 1 pet, 1 travel side-activity.
- Aktivite değişimi: server transaction; eski contract closed → new contract created.
- Location/requirement değişince (ör. travel bitti): aktif contract yeni location'a re-validated.
- Aktivite tanımları veri-tabanlıdır; kod yalnızca "simulate" ve "validate" implementation'ı sağlar.

## 6. Kullanıcı Akışı
```
UI → "Start Mining at Iron Valley"
 → server validate
 → create contract
 → running
 → user closes app
 → later reveal
```

## 7. İlişkiler
- [AFK Engine](afk-engine.md), [Reward System](reward-system.md), [Loot System](loot-system.md), [Random System](random-system.md), [Notification System](notification-system.md)
- Her aktivite domain sistemi: [Hunting](../game-design/hunting.md), [Adventure](../game-design/adventure.md), [Professions](../game-design/professions.md), [Travel](../game-design/travel.md), [Dungeons](../game-design/dungeons.md)

## 8. Gereken Veriler
```
activity_type_registry
  code, display_name, category, params_schema
  simulate_ref (fn ref), validate_ref (fn ref)
  soft_cap_override_ms (nullable)

activity_instance
  id, character_id, activity_type
  params_json, seed, started_at, last_tick_at, status
  contract_id (nullable link)

activity_event_log (compressed)
  instance_id, tick_index, event_type, payload
```

## 9. I/O
- Input: start/pause/resume/interrupt/complete.
- Output: event log, yield deltas, notification triggers.

## 10. Failure Durumları
- Validation fail: activity başlatılmaz.
- Simulate exception (data drift): safe fallback yield + admin notification.
- Contract orphaned: cleanup job.

## 11. Edge Cases
- Aktivite tipi kaldırıldı (patch): orphan instance → auto-complete, mail refund.
- Parametre değişikliği (spot moved): compensation.

## 12. Server Authoritative Gereksinimleri
- Tüm life-cycle transaction.
- Simulate fonksiyonu pure & deterministic.
- Contract state single source of truth.

## 13. Genişletilebilirlik
- Yeni activity type: registry'e ekle + simulate + validate implement + params_schema.
- Sub-activities (pet, mount travel) parallel channel.

## 14. Örnek — Simulate Sözleşmesi (pseudo)
```
simulate(instance, from, to, character_snapshot) -> YieldReport
  ticks = compute_ticks(from, to, tick_rate)
  yields = []
  for tick in ticks:
    if stop_condition(character_snapshot): break
    yields.append(one_tick(seed, tick_index, character_snapshot))
    apply_yield_to_snapshot(character_snapshot, yields[-1])
  return aggregate(yields)
```

## 15. Balance Considerations
- Tick rate default 30 sn AFK-scale aktiviteler için (mining, gathering, hunting bulk).
- Combat AFK için 1s bulk (dungeon idle) yeterli.
- Simulate performansı: 10k karakter/sec reveal hedeflenir.
