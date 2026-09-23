# System Architecture (Sistem Mimarisi)

**Karar tarihi:** 2026-09-08
**Faz:** Phase 00 → Phase 01 köprüsü
**Karar sahibi:** Proje sahibi (kullanıcı)

Bu doküman oyunun uç uca teknik omurgasını tanımlar. Alt sistem dokümanlarındaki (AFK Engine, Activity Engine, vb.) tüm "server authoritative" gereksinimleri burada tanımlı katmanlar üzerinde koşar.

---

## 1. Yüksek Seviye Diyagram

```
┌────────────────────────────────────────────┐
│              Flutter Client                │   ← UI, presentation, input
│  (Web / iOS / Android / Desktop shells)    │
└──────────────────┬─────────────────────────┘
                   │  HTTPS / WSS
                   ▼
┌────────────────────────────────────────────┐
│              API Gateway (HTTP + WS)       │   ← auth, routing, rate-limit
│      REST for state, WebSocket for push    │
└──────────────────┬─────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────┐
│         Go Application Server              │   ← business logic entry
│  (goroutine-based; horizontally scalable)  │
└──────────────────┬─────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────┐
│            Domain Modules (Go)             │
│  ┌──────────┬──────────┬──────────────┐    │
│  │ Character│ Activity │ AFK Engine   │    │
│  │ Combat   │ Economy  │ Marketplace  │    │
│  │ Guild    │ Party    │ Notification │    │
│  │ Profession│ Loot    │ RNG          │    │
│  └──────────┴──────────┴──────────────┘    │
└──────┬────────────────────────────┬────────┘
       │                            │
       ▼                            ▼
┌──────────────────┐        ┌──────────────────┐
│    PostgreSQL    │        │      Redis       │
│  (durable state) │        │  (ephemeral,     │
│  characters,     │        │   session,       │
│  inventory,      │        │   active ticks,  │
│  ledger,         │        │   pub/sub,       │
│  audit,          │        │   rate-limit,    │
│  events log      │        │   leaderboards)  │
└──────────────────┘        └──────────────────┘
       ▲                            ▲
       │                            │
       └──────────┬─────────────────┘
                  │
                  ▼
        ┌────────────────────────┐
        │        Workers          │   ← async tasks
        │  (Go background procs)  │
        │  • Reveal batch          │
        │  • Marketplace expiry    │
        │  • World-boss spawns     │
        │  • Notification dispatch │
        │  • Economy metrics       │
        │  • Cleanup / archive     │
        └────────────────────────┘
```

---

## 2. Katman Sözleşmesi

### 2.1 Flutter Client
- Tek codebase; web, iOS, Android, desktop shell.
- **Sadece sunum**: state gösterir, input toplar. Otorite değildir.
- Push kanal: WebSocket (activity updates, notifications, chat).
- Offline dayanıklılık: cache okuma; yazma her zaman server confirm bekler.
- Server clock skew UI'da tolere edilir; ödül hesabında istemci saati kullanılmaz.

### 2.2 API Layer
- **REST**: state okuma/yazma (inventory, market, character sheet).
- **WebSocket / WSS**: canlı push (activity tick, notification, chat, boss spawn).
- Auth: JWT (access + refresh); session revocation Redis üzerinden.
- Rate limit: per-route + per-account + per-IP.
- Idempotency-Key header desteği (kritik POST'lar için).
- Versioning: `/v1/...`; breaking change'ler yeni versiyon.

### 2.3 Go Application Server
- **Dil**: Go (goroutine-based concurrency, low latency, native binary).
- Horizontal scale: stateless app node'ları; sticky session yok (session state Redis'te).
- Graceful shutdown; long-running goroutine'ler için context cancellation.
- Observability: OpenTelemetry (trace, metric, log).
- Deployment: container (Docker) → orchestrator (Kubernetes veya benzeri; Phase 01'de karar).

### 2.4 Domain Modules (Go)
Her domain modülü:
- Kendi paketi (`internal/character`, `internal/afk`, `internal/economy`, ...).
- Kendi domain type'ları, use-case fonksiyonları, repository interface'i.
- Cross-module etkileşim: **event bus** veya explicit function call (module olarak; global singleton yok).
- Test edilebilirlik: repository interface + in-memory impl.

Ana modüller Phase 01'de bu ana havuza karşılık gelir:
- `character` (creation, profile, level, stats)
- `race`, `class`, `specialization` (data-driven registry)
- `skills`, `combat` (skill exec, damage calc, tick loop)
- `activity`, `afk` (Activity Engine + AFK Engine)
- `inventory`, `equipment`, `items`
- `economy`, `marketplace`, `ledger`
- `professions`, `loot`
- `world`, `travel`, `map`, `city`
- `dungeons`, `world_bosses`, `events`
- `party`, `guild`, `social`
- `quests`, `achievements`, `mastery`, `ascension`
- `notification`
- `rng` (deterministic random)

### 2.5 PostgreSQL (Durable State)
- **Kaynak-of-truth** her kalıcı veri için.
- İşlem-güvenli (ACID) transaction'lar: currency ledger, item transfer, marketplace, character mutation.
- Ana şemalar:
  - `character.*` (profile, stats, inventory, equipment, activity_contract, professions, mastery, ascension)
  - `content.*` (item_template, race, class, spec, skill, recipe, dungeon_template, ...)
  - `economy.*` (currency_ledger, market_listing, market_transaction, mail)
  - `social.*` (guild, guild_member, guild_bank, party, friendship, chat_message, report)
  - `world.*` (continent, region, territory, city, location, poi, npc_template, resource_node)
  - `audit.*` (event log, admin actions, security signals)
- Migrations versionlanır (ör. `sql-migrate` veya `goose`).
- Indexler her sıcak query için tasarlanır (character_id, item_instance_id, listing_id, ...).

### 2.6 Redis (Ephemeral State)
- Session state, active connection map.
- Rate limiter buckets.
- Aktif tick counter'lar (activity live cache).
- Pub/Sub: chat channel fan-out, notification fan-out, world-boss announcement.
- Sorted sets: leaderboards (Phase 02 for PvP), boss contribution live.
- Distributed lock (Redlock) — kısa süreli critical section (opt-in).
- **Kural**: Redis kaybolursa oyun **canlılığı** etkilenir (chat, live push) ama **veri** kaybolmaz.

### 2.7 Workers (Async / Background)
- Go binary (aynı codebase, farklı entrypoint).
- İş türleri:
  - **Reveal batch worker** — offline'dan dönen büyük reveal'ları arka planda hazırlar.
  - **Marketplace expiry** — listing bitiş, iade mail.
  - **World-boss scheduler** — spawn window rasgele seçim.
  - **Event scheduler** — daily/weekly/seasonal event start/stop.
  - **Notification dispatcher** — push, email, in-game routing.
  - **Economy metrics** — commodity price index, gold supply monitoring.
  - **Cleanup / archive** — soft-deleted characters, expired mails, contract'ların hard cap ötesi arşivi.
  - **Audit reconciliation** — currency ledger vs character wallet fark denetimi.
- Job kuyruğu: PostgreSQL-based queue (SKIP LOCKED) veya Redis stream (Phase 01 karar).
- Idempotent job pattern zorunlu; retry backoff.

---

## 3. Cross-Cutting Kararlar

### 3.1 Time Source
- **UTC epoch (ms)**, monotonic; NTP sync all nodes.
- Timezone yalnız sunumda; asla hesabın parçası değil.

### 3.2 RNG
- Deterministic (bkz. [Random System](../systems/random-system.md)).
- Server-side only; go module (`internal/rng`).
- Seed persist per contract.

### 3.3 Event Sourcing (Selektif)
- Ekonomi + ödül olayları event log'a yazılır (append-only).
- Karakter state event'lerin projeksiyonu; snapshot + append-log hibriti.
- Reconstruction: gerektiğinde event'lerden state yeniden inşa edilebilir (audit + rollback safety).

### 3.4 Migration Stratejisi
- Data-driven içerik (item, race, class): hot-reloadable data table + version.
- Schema migration: forward-only tercih; nadiren rollback (down-migration test edilmiş olmalı).
- Content patch'i canlı sunucuya deploy: staged rollout + shadow read.

### 3.5 Observability
- **Metrics**: Prometheus-compat (Go exporter).
- **Traces**: OpenTelemetry (uzun kuyruk request tracing).
- **Logs**: structured JSON, correlation_id her request'te.
- **Dashboards**: player concurrency, active contract count, reveal latency p50/p99, marketplace TPS, gold supply.

### 3.6 Security
- **Auth**: JWT + refresh; refresh rotation.
- **TLS** her yerde (mTLS internal opsiyonel).
- **Secrets**: vault benzeri, env variable değil.
- **Anti-cheat sinyaller**: reveal delta anomali, market flood, RMT pattern.
- **Rate limits** her endpoint için tanımlı ve enforce.
- **PII minimizasyon**: e-mail sadece auth ve notification için; oyun içi ekranda gösterilmez.

### 3.7 Testing
- Unit test: domain modül seviyesinde saf fonksiyonlar.
- Integration test: Postgres testcontainers ile.
- Determinism test: RNG + AFK Engine golden replay.
- Load test: reveal throughput, market TPS, chat fanout.

---

## 4. Non-Functional Hedefler (Phase 01 baseline)

| Metrik | Hedef |
|---|---|
| Login p50 | < 300 ms (reveal dahil) |
| Login p99 | < 1500 ms |
| Reveal throughput | ≥ 5000 karakter/sn (batch) |
| API RPS / node | ≥ 3000 |
| WebSocket concurrent conn / node | ≥ 20000 |
| Marketplace TPS (peak) | ≥ 500 |
| Postgres write TPS (peak) | ≥ 4000 |
| Data durability | RPO ≤ 1 dk, RTO ≤ 15 dk |
| Uptime | ≥ %99.5 (v1); %99.9 (v2 hedef) |

---

## 5. Yığın Seçim Gerekçesi (özet)

- **Flutter**: tek codebase ile web + mobile + desktop; MMO'nun cross-platform vizyonuna uygun (bkz. [GDD §14](../GDD.md)).
- **Go**: high concurrency, low GC pause, native binary, deployment kolaylığı, event-driven mimariler için verimli.
- **PostgreSQL**: ACID, JSONB desteği (esnek content), extension ekosistemi, mature.
- **Redis**: ephemeral state + pub/sub + rate limit için endüstri standardı.
- **Workers (Go)**: aynı dil, aynı binary → operational simplicity.

Detaylı alternatif değerlendirmesi: [Design Decisions DD-026](../decisions/design-decisions.md).

---

## 6. Bekleyen Teknik Kararlar (Phase 01)

- Orchestrator: Kubernetes vs Nomad vs simple systemd (v1 scale bağlı).
- Job queue: Postgres SKIP LOCKED vs Redis stream vs harici (NATS).
- Container registry & CI/CD pipeline.
- Config management: env + config file + Vault mix.
- Log aggregation: Loki vs ELK vs cloud provider.
- Blue/green vs canary deployment stratejisi.
- Multi-region persistence (Phase 02).

---

## 7. İlişkiler
- [GDD §15 Technical Design Principles](../GDD.md)
- [GDD §18 Phase 01 Teknik Gereksinim Özeti](../GDD.md)
- [AFK Engine](../systems/afk-engine.md)
- [Activity Engine](../systems/activity-engine.md)
- [Random System](../systems/random-system.md)
- [Reward System](../systems/reward-system.md)
- [Notification System](../systems/notification-system.md)
- [Economy](../game-design/economy.md), [Marketplace](../game-design/marketplace.md)
