# CCU Scenarios (Karışım + Model)

## 1. Amaç
Her CCU tier için gerçekçi kullanıcı davranış karışımı, endpoint per-user rate'i ve toplam server yükü modellemesi.

## 2. Karışım (Realistic Mix)

Bir "active player" saatlik davranış modeli (v1 hipotez; Phase 02 telemetry ile tune):

| Aksiyon | Freq / saat / player | Notlar |
|---|---|---|
| Character sheet load | 6 | Ekran açılış |
| Inventory view | 4 | — |
| Activity reveal | 2 | 30dk-1h contract |
| Activity start | 2 | Yeniden |
| Vendor buy/sell | 1.5 | — |
| Market search | 3 | Browsing |
| Market list/buy | 0.5 | — |
| Chat send | 20 | Global + party + guild karışım |
| Combat skill cast | 40 | Live combat sessions (10-15 dk her biri) |
| Dungeon enter | 0.3 | 2-3 saatte bir |
| Guild interaction | 1 | Bank/chat/members |
| Mail send/read | 0.5 | — |
| Party form/invite | 0.3 | — |
| WS ping | 120 | Her 30sn |

AFK-only player için: WS ping + presence heartbeat (60/saat) + arada 1 reveal + 1 sheet load.

## 3. Ratio Active:AFK

`active:afk = 3:7` (v1 hipotez). MMO idle genre'de ideal:
- Active peak saat (akşam prime time): 40:60.
- Non-peak (gündüz iş saati): 20:80.

Test senaryosu **realistic worst case**: peak (3:7 = %30 active).

## 4. Tier Yük Modeli

### T1 — 1.000 CCU (300 active + 700 AFK)

- **API request/sn**: 300 × (endpoint sum) / 3600 ≈ **~40 rps** (chat hariç).
- **Chat send/sn**: 300 × 20 / 3600 ≈ **1.7/sn**.
- **Combat skill/sn**: 300 × 40 / 3600 ≈ **3.3/sn**.
- **WS connection**: 1.000 concurrent.
- **WS message/sn**: 1.000 × 120 / 3600 ≈ **33/sn** (ping baseline) + chat fanout + combat event.
- **Reveal batch**: 300 × 2 / 3600 ≈ **0.17/sn** (analitik hesap yükü düşük).

### T2 — 5.000 CCU (1.500 active + 3.500 AFK)

- **API**: ~200 rps.
- **Chat**: 8.3/sn.
- **Combat skill**: 16.7/sn.
- **WS conn**: 5.000.
- **WS msg**: 165/sn ping baseline + fanout.
- **Reveal**: 0.83/sn.

### T3 — 10.000 CCU (3.000 active + 7.000 AFK)

- **API**: ~400 rps.
- **Chat**: 16.7/sn.
- **Combat skill**: 33.3/sn.
- **WS conn**: 10.000.
- **WS msg**: 333/sn ping + fanout.
- **Reveal**: 1.67/sn.
- **DB QPS**: ~800-1200 read/sec + ~200 write/sec (measured target).

### T4 — 50.000 CCU (15.000 active + 35.000 AFK)

- **API**: ~2000 rps.
- **Chat**: 83/sn.
- **Combat skill**: 167/sn.
- **WS conn**: 50.000.
- **WS msg**: 1.667/sn ping + massive fanout.
- **Reveal**: 8.3/sn.
- **DB QPS**: ~5000 read + ~1000 write.
- **Redis QPS**: ~20.000 (pubsub + cache + rate limit).

### T5 — 100.000 CCU (Phase 03 target)

- Multi-region veya shard gerekir.
- Vertical scaling limitleri v1 tek DB için aşılır → **shard karar Phase 03**.

## 5. Burst Modeli

- **Login burst** (patch günü peak +5): 1 dk içinde CCU'nun %50'si login attempt.
- **World event start** (dış patch broadcast): 30 sn içinde %30 aktif event location'a travel + attend.
- **Server maintenance return**: 5 dk içinde %90 CCU login attempt.

Test'ler bu burst'leri simüle eder; auto-scaling + rate limit doğru davranmalı.

## 6. Ramp Profil

Standard test ramp:
```
0-5 dk:   0 → hedef CCU (ramp up)
5-65 dk:  sabit yük (steady state 60 dk)
65-70 dk: burst spike (test scenario)
70-90 dk: cool down + recovery observation
```

Uzun soak (nightly): 4 saat sabit yük; memory leak + connection pool drift tespit.

## 7. Data Fixture

Test için fake account/character pool:
- T1: 5.000 karakter (rotation).
- T3: 50.000 karakter.
- T4: 250.000 karakter (fixture seed script).

Fixture pack seed cron (perf lab):
- Random level 1-60 dağılımı.
- Item envanter dolulu (25-45 slot).
- Guild membership %60.
- Friend list 5-30.

## 8. Cross-Reference
- [README.md](README.md), [capacity-planning.md](capacity-planning.md)
- Per-endpoint testler: [load-*.md](load-login.md)

## 9. Bekleyen Kararlar
- Real telemetry ile mix tune Phase 02.
- Regional split modeli Phase 03.
- Fixture seed script otomasyon Phase 01.
