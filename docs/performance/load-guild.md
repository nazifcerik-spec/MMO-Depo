# Load Test: Guild

## 1. Amaç
Guild sayfa yükleme, chat throughput, bank operasyonu, member roster query.

## 2. Senaryolar

### 2.1 Guild Home View
- Yük: T3'te 3000 active × 1 guild_view/saat = **~0.8/sn**.
- Compound endpoint (guild + members + treasury + recent activity).
- Beklenti: p95 < 300ms.

### 2.2 Guild Chat Steady
- Yük: T3'te ~2500 guild üye online (50 guild × ~50 üye).
- Rate: ~10 msg/dk per guild → guild başına 0.17/sn; system-wide 8.3/sn (50 guild).
- WS fanout: her mesaj guild üye sayısına göre (avg 30 subscriber).
- Beklenti: chat send p95 < 150ms; fanout p95 < 200ms.

### 2.3 Guild Chat Burst (Event Rally)
- Boss summon anons → 30 sn içinde ~20 mesaj / guild peak.
- Test: 100 guild simultaneous → 2000 chat msg/dk = 33/sn peak.
- WS fanout scale: 33 × 30 avg subscribers = 1000 WS msg/sn.

### 2.4 Bank Concurrent Withdraw
- 5-10 officer aynı guild bank'tan concurrent withdraw.
- SERIALIZABLE TX + treasury update + audit.
- Beklenti: daily cap enforce; no double-withdraw; latency p95 < 400ms.

### 2.5 Roster Query
- Members endpoint pagination (50 üye).
- Rate: 5 view/sn system-wide.
- Beklenti: p95 < 250ms.

## 3. Kritik Path — Bank Withdraw

1. Auth + permission check (rank permission).
2. SERIALIZABLE TX begin.
3. Guild treasury SELECT FOR UPDATE.
4. Daily withdraw cap check (Redis: `guild:withdraw:daily:{gid}:{cid}` counter).
5. Ledger: guild treasury -X → character wallet +X.
6. Audit `audit.guild_bank_action`.
7. Cap counter increment.
8. WS event guild-scope.
9. COMMIT.

Bottleneck: SERIALIZABLE lock hold (~30-80ms); concurrent withdraw retry.

## 4. Metrics

- `guild_view_p95_ms`.
- `guild_chat_send_p95_ms`.
- `guild_chat_fanout_p95_ms`.
- `guild_bank_withdraw_p95_ms`.
- `guild_bank_serialization_retry_total`.
- `guild_bank_cap_hit_total`.

## 5. Redis Yapıları

- `guild:{gid}:members:online` set (WS subscriber list).
- `guild:{gid}:presence` — üye online sayacı.
- `guild:withdraw:daily:{gid}:{cid}` counter (TTL 24h from daily reset).
- `chat:live:guild:{gid}` pubsub + rolling list MAXLEN 200.

## 6. Fanout Cost

Guild chat mesajı:
- Redis PUBLISH → subscriber connection'lar iterate → serialize + send.
- 30 subscriber × 50 guild × 10 msg/dk = ~250 WS msg/dk = 4/sn baseline.
- Burst rally 33 msg/sn × 30 = 1000 WS msg/sn.

## 7. MOTD Change Event

Leader MOTD change → guild üye WS "motd_changed" event fanout.
Rate düşük (günde birkaç kez); test skip veya nominal.

## 8. Beklenen Sonuç

T3 steady:
- Guild view p95: ~200ms (compound cache).
- Chat send p95: ~100ms.
- Chat fanout p95: ~150ms.
- Bank withdraw p95: ~250ms.
- No cap breach; audit trail complete.

## 9. Cross-Reference
- [guild/guild.md](../guild/guild.md), [guild/bank.md](../guild/bank.md), [social/guild-chat.md](../social/guild-chat.md)

## 10. Bekleyen Kararlar
- Guild home compound cache TTL (v1: 30sn) Phase 02.
- Bank withdraw daily cap enforce Redis vs DB atomic (v1 Redis; Phase 02 DB).
- Guild size scaling (100+ member) Phase 02.
