# Design Decisions Log

Bu dosya Phase 00 sırasında alınan tüm önemli tasarım kararlarının kaydıdır. Her karar için:

- **Karar (Decision)**
- **Gerekçe (Reason)**
- **Değerlendirilen alternatifler (Alternatives)**
- **Sonuçlar (Consequences)**
- **Tarih**

formatı kullanılır.

Kararlar tarih sırasına göre yeniden yazılmaz; append-only'dir. Yeni bir karar bir eskisini geçersiz kılıyorsa: eski kararın yanına "SUPERSEDED BY DD-XXX (tarih)" notu düşülür.

---

## DD-001 — Doküman Dili
**Karar:** Ana doküman dili Türkçe, teknik terimler İngilizce korunur.
**Reason:** Kullanıcı hem Türkçe hem İngilizce dil tercihi belirtti. Türkçe okunabilirlik + İngilizce endüstri standardı en iyi kompromis.
**Alternatives:**
1. Tam çift dil (her doküman iki kez) — reddedildi: bakım maliyeti yüksek.
2. Tamamen İngilizce — reddedildi: kullanıcının primary dili Türkçe.
3. Tamamen Türkçe — reddedildi: teknik terim tercümesi (buff, tick, DPS) doğallık kaybettirir.
**Consequences:** Yeni doküman kabul kriteri: Türkçe okunabilirlik + İngilizce teknik terim korunması.
**Tarih:** 2026-09-07

---

## DD-002 — Server-Authoritative Mimari
**Karar:** Tüm state, ödül hesabı, RNG ve time source server tarafındadır. İstemci öneri gönderir, server doğrular.
**Reason:** MMO ekonomisinin ve AFK sisteminin bütünlüğü için istemciye güvenilemez.
**Alternatives:** Hybrid (client bazı hesapları yapar) — reddedildi: exploit riski.
**Consequences:** Tüm doküman şablonlarında "Server Authoritative Gereksinimleri" başlığı zorunlu.
**Tarih:** 2026-09-07

---

## DD-003 — AFK Birinci Sınıf Gameplay
**Karar:** AFK / offline progression, yardımcı bir feature değil, ana gameplay pillar'ıdır.
**Reason:** Oyun vizyonu "persistent idle MMO". Ana ayakların ilki.
**Alternatives:** AFK sadece küçük bonus verir, ana progression aktif oynama gerektirir — reddedildi: vizyona aykırı.
**Consequences:** Tüm aktiviteler AFK-uyumlu tasarlanır (Activity Engine + AFK Engine soyutlaması).
**Tarih:** 2026-09-07

---

## DD-004 — Class Değişmez, Spec Değişebilir
**Karar:** Karakter yaratıldıktan sonra class değiştirilemez; specialization değiştirilebilir (respec).
**Reason:** Class kimliği taahhüde bağlı olmalı; spec, denemeye açık esneklik alanı.
**Alternatives:** Class de değişebilir olsun (dual class) — Phase 03'e ertelendi; şimdilik kimlik önceliği.
**Consequences:** Class rework olursa mevcut karakterlere free respec verilir; olası "class silme" hallerinde legacy koruması.
**Tarih:** 2026-09-07

---

## DD-005 — Race Class'ı Kilitlemez
**Karar:** Herhangi bir ırk herhangi bir class'ı oynayabilir.
**Reason:** "Best race for class" tuzağı optionality'yi öldürür; kimlik ve mekanik ayrı katmanlar.
**Alternatives:** Race-lock (belirli class'lar sadece belirli ırklara) — reddedildi: modern MMO tercihi.
**Consequences:** Race stat/passive'leri class'a bağımsız balanced olmalı; sinerji öneri olur, zorunluluk olmaz.
**Tarih:** 2026-09-07

---

## DD-006 — Data-Driven Class ve Race
**Karar:** Yeni class, race, item, aktivite eklemek için kod değişikliği değil, data tanımı yeterli olacak şekilde tasarım.
**Reason:** İçerik pipeline'ının hız kazanması; class/race sayısının uzun vadede büyümesi.
**Alternatives:** Hard-coded (Diablo tarzı) — reddedildi: genişletme yavaş.
**Consequences:** Class ve race'in ability, stat modifier, spec ağacı gibi tüm parçaları schemas ile ifade edilir. Kod sadece "combiner" olur.
**Tarih:** 2026-09-07

---

## DD-007 — Deterministic RNG
**Karar:** Tüm RNG seed + counter tabanlı, versiyonlu; replay ve audit için deterministic.
**Reason:** AFK ödüllerinin auditability'si + oyuncu güveni.
**Alternatives:** Live true-random kaynak — reddedildi: replay ve audit imkansız.
**Consequences:** RNG algoritma versiyonu contract'a bağlanır; algoritma değişince yeni contract yeni versiyon kullanır, eskiler bozulmaz.
**Tarih:** 2026-09-07

---

## DD-008 — Global Marketplace
**Karar:** Marketplace kıta bazlı global.
**Reason:** Likidite; adil fiyat oluşumu; sosyal ölçek.
**Alternatives:** Region-specific market — reddedildi: sunucu erken evrede yeterli oyuncu yoğunluğuna sahip olmayabilir; parçalamak zararlı.
**Consequences:** Fiyat rekabeti global; region-specific pazar Phase 02 opsiyonu olarak açık kalır.
**Tarih:** 2026-09-07

---

## DD-009 — Never Sell Power
**Karar:** Monetization pathway'inde asla stat, XP boost, ekipman satılmaz.
**Reason:** "Fair play" oyuncu güveni; oyunun ana ekonomisinin sağlığı.
**Alternatives:** Time-skip / stat boost — reddedildi: sömürücü.
**Consequences:** Monetization sadece kozmetik, karakter slot, isim değişikliği, envanter genişleme (mantıklı sınırlar).
**Tarih:** 2026-09-07

---

## DD-010 — Karakter Silinmez (Death Penalty Yumuşak)
**Karar:** Ölüm karakter silmez. Ceza XP debt + durability loss + kısa spawn.
**Reason:** Idle oyunda drastic loss oyuncu kaybına yol açar.
**Alternatives:** Hardcore mode karakter permadeath — Phase 03 opsiyonu; şu an ana modda hayır.
**Consequences:** Penalty ölçeği (XP debt %) balance pass'i Phase 01'e devredildi.
**Tarih:** 2026-09-07

---

## DD-011 — Ascension Reset ile Kalıcı Meta Bonus
**Karar:** Cap sonrası dikey ilerleme mekaniği "ascension": karakter level 1'e döner, envanter/spec/ekipman korunur, ascension_point +N.
**Reason:** Endgame kilit; klasik "prestige" mekaniği yeniden başlangıç hissi ile long-term ilerleme.
**Alternatives:** Sadece mastery/horizontal — reddedildi: dikey ilerleme özlemi ciddi.
**Consequences:** Ascension frequency: aktif oyuncu ~60-100 saat aktif oyun/ascension hedefi. Balance Phase 01.
**Tarih:** 2026-09-07

---

## DD-012 — Personal Loot Default
**Karar:** Party loot default = personal loot (herkes kendi roll'unu alır).
**Reason:** Adil, drama önleyici, casual friendly.
**Alternatives:** FFA / round-robin / master loot — hepsi opsiyonel Phase 02, ancak default personal.
**Consequences:** Loot table hesabı server tarafında party içi paralel roll üretmeli.
**Tarih:** 2026-09-07

---

## DD-013 — Aktivite Modeli: 1 Ana + Side Kanallar
**Karar:** Karakter aynı anda 1 ana aktivite + 1 pet + 1 travel side-activity çalıştırabilir. Ana aktivite tek.
**Reason:** Karar odaklılık; "her şeyi aynı anda yap" kaosunu engellemek.
**Alternatives:** Sınırsız paralel aktivite — reddedildi: sistemin özünü bulanıklaştırır.
**Consequences:** Aktivite değişimi transaction; oyuncu bir aktivite başlatırken diğerini kapatır.
**Tarih:** 2026-09-07

---

## DD-014 — Soft Cap Curve, Hard Cap 30 Gün
**Karar:** AFK yield 24h tam; 24-48h %60; 48-168h %30; 168h+ 0. 30 gün ötesinde arşive taşınır.
**Reason:** Daily engagement teşviki + tatildeki oyuncuların "yıkıcı" hissetmemesi.
**Alternatives:** Fixed 8h cap (Diablo Immortal tarzı) — reddedildi: casual friendly değil. Unlimited cap — reddedildi: exploit + ekonomi bozulur.
**Consequences:** Curve parametreleri aktivite bazlı override edilebilir.
**Tarih:** 2026-09-07

---

## DD-015 — Marketplace Fee Modeli
**Karar:** Listing fee %2 (upfront, non-refundable) + sale fee %5 (satışta kesinti).
**Reason:** Sink; spam engelleyici; ekonomik denge.
**Alternatives:** Sadece sale fee (Blizzard AH) — reddedildi: listing spam. Sadece listing fee — reddedildi: yeterli sink değil.
**Consequences:** Phase 01 balance pass'i ile revize edilebilir; monitoring metrikleri kuruldu.
**Tarih:** 2026-09-07

---

## DD-016 — Weekly Reset Zamanı
**Karar:** Weekly reset Salı 06:00 UTC.
**Reason:** Global çoğunluk için "hafta içi orta" — hafta sonu content'ini hazır tutar; hiçbir bölge için "gece yarısı" değil.
**Alternatives:** Pazartesi 00:00 UTC (WoW style) — pek çok TZ için hafta sonu düşer, olumsuz.
**Consequences:** Weekly lockouts, event start/end, daily quest rotation buna bağlanır.
**Tarih:** 2026-09-07

---

## DD-017 — Home City Sistemi
**Karar:** Her karakter bir home city belirler; respawn burada; değiştirme cooldown ile mümkün.
**Reason:** Kimlik + logistical anchor.
**Alternatives:** Free respawn (herhangi bir şehir) — reddedildi: MMO scale hissi kaybolur.
**Consequences:** Home city change fee + cooldown Phase 01.
**Tarih:** 2026-09-07

---

## DD-018 — Guild Size Başlangıç 50
**Karar:** Guild max size 50, guild level ile 100'e kadar.
**Reason:** Küçük/orta guild kültürü + guild progression için genişleme yolu.
**Alternatives:** Sınırsız — reddedildi: mega-guild dominance riski. 20 — çok küçük.
**Consequences:** Guild war (Phase 02) tasarımı bu ölçeğe göre yapılır.
**Tarih:** 2026-09-07

---

## DD-019 — Ölüm Sırasında Item Kaybı Yok
**Karar:** Ölüm karakter item'ını düşürmez.
**Reason:** Adventure MMORPG psikolojisi + casual friendly.
**Alternatives:** UO/EQ tarzı item loss — reddedildi: hedef kitleye aykırı.
**Consequences:** Death penalty daha zayıf olduğundan diğer kanallardan (XP debt, durability) beslenir.
**Tarih:** 2026-09-07

---

## DD-020 — Tek Dünya, Tek Shard (v1)
**Karar:** Phase 01'de single-shard MMO. Sharding stratejisi Phase 02+ için ertelendi.
**Reason:** Ekonomi ve sosyal katmanın tam ölçekli hissedilmesi.
**Alternatives:** Multi-shard v1 — reddedildi: karmaşıklık; sharded MMO'nun "boş dünya" hissi.
**Consequences:** Scaling stratejisi seçildi (soft-sharding şehir instance'ları Phase 01 opsiyonel).
**Tarih:** 2026-09-07

---

## DD-021 — Envanter Slot Bazlı, Ağırlık Yok
**Karar:** Envanter sadece slot sınırıyla çalışır (Phase 00); item weight sistemi yok.
**Reason:** Sadelik; slot doluluk kavramı UX'te net.
**Alternatives:** WoW-style bag + slot + weight — reddedildi: karmaşıklık artışı düşük değer.
**Consequences:** İleride ağırlık ekleneceğinde `weight` alanı zaten schema'da; bag upgrade path buna göre planlanır.
**Tarih:** 2026-09-07

---

## DD-022 — Resource Node Personal Instance
**Karar:** Kaynak node'ları personal instance (her oyuncu kendi node'unu görür).
**Reason:** MMO yoğun olduğunda gathering competition frustration yaratır.
**Alternatives:** Shared world node (klasik) — reddedildi: idle MMO'da farming kolluğa dönüşür.
**Consequences:** Rare node (event) shared kalabilir; ekonomik trade-off korunur.
**Tarih:** 2026-09-07

---

## DD-023 — Cross-Character Trade: Marketplace + Mail
**Karar:** Karakterler arası ticaret sadece marketplace ve mail üzerinden; direct trade yok.
**Reason:** Scam engelleme + ekonomi audit trail.
**Alternatives:** Direct trade window (WoW) — Phase 02 opsiyonu, ancak escrow ile.
**Consequences:** Adil ticaret; ancak "arkadaşa hediye" mail üzerinden yapılır.
**Tarih:** 2026-09-07

---

## DD-024 — Pets ve Mounts BoP
**Karar:** Pet ve mount instance'ları karaktere bound; trade edilemez.
**Reason:** Koleksiyon kimliği; monetization koruması.
**Alternatives:** Trade edilebilir mount — reddedildi: kozmetik kimliğin ekonomik metaya dönüşmesi.
**Consequences:** Ancak mount unlock item'ları BoE olabilir; item satılır, mount unlock alıcıda gerçekleşir.
**Tarih:** 2026-09-07

---

## DD-025 — Doküman Standardı: 17-madde Şablon
**Karar:** Her sistem dokümanı 17 başlıklı standart iskeleti izler (bkz. GDD §20).
**Reason:** Tutarlı, karşılaştırılabilir, eksik başlık = kritik bilgi eksiği görünürlüğü.
**Alternatives:** Sistemine göre serbest format — reddedildi: uzun vadede parçalanır.
**Consequences:** Yeni sistem dokümanı eklenirken şablon takip zorunlu.
**Tarih:** 2026-09-07

---

## DD-026 — Teknik Yığın (Stack)
**Karar:** Flutter (client) → API (REST + WebSocket) → Go Application Server → Domain Modules → PostgreSQL (durable state) + Redis (ephemeral state) + Workers (Go, async).
**Reason:**
- **Flutter**: tek codebase → web + iOS + Android + desktop; MMO'nun cross-platform vizyonuna (GDD §14) doğrudan uygun.
- **Go**: goroutine-based concurrency, düşük GC pause, native binary, event-driven MMO backend'ler için verimli; ekosistemi olgun.
- **PostgreSQL**: ACID, JSONB (esnek content), audit-friendly, mature — currency ledger ve marketplace transaction'ları için doğal ev.
- **Redis**: ephemeral session, live tick cache, pub/sub, rate limit için endüstri standardı; kaybı verinin değil canlılığın etkilenmesi ilkesi.
- **Workers (Go)**: aynı dil/binary — operasyonel sadelik; reveal batch, marketplace expiry, event scheduling, notification dispatch.
**Alternatives:**
1. **Node.js / TypeScript backend** — reddedildi: Go'ya kıyasla long-lived tick loop'larda GC/concurrency dezavantajı.
2. **Elixir / Erlang (BEAM)** — güçlü aday; reddedildi: geliştirici havuzu darlığı ve içerik pipeline'ı için Go daha pragmatik.
3. **Unity / .NET backend** — reddedildi: web-first idle MMO için ağır; native mobile shell zorunlu değil.
4. **React Native / native (Kotlin+Swift) client** — reddedildi: Flutter tek codebase avantajı + web desteği.
5. **NoSQL (Mongo) primary store** — reddedildi: ekonomi için ACID zorunlu; JSONB Postgres yeterli esneklik veriyor.
6. **RabbitMQ / Kafka job queue v1** — Phase 01'de küçük ölçek için Postgres SKIP LOCKED veya Redis stream yeterli; harici broker Phase 02 kararı.
**Consequences:**
- Domain Modules pattern: her modül kendi paketi, repository interface + test edilebilirlik.
- Deployment: containerized Go binary; orchestrator seçimi (K8s vs simpler) Phase 01 açık kararı.
- Observability standardı: OpenTelemetry + Prometheus + structured JSON logs.
- Non-functional hedefler (login p50 < 300ms, reveal ≥ 5000 char/sn) bu yığına göre kalibre edildi.
- İçerik pipeline: hot-reloadable data table + versionlu schema migration.
**Tarih:** 2026-09-08

---

## DD-027 — Trunk-Based Development + Conventional Commits
**Karar:** Branch stratejisi trunk-based (kısa ömürlü feature branch → PR → squash merge to `main`). Commit mesaj formatı Conventional Commits.
**Reason:** Small PR + fast feedback + linear history + otomatik changelog gen; MMO'nun sık iterasyon ihtiyacına uygun.
**Alternatives:** GitFlow (develop + release branches) — reddedildi: karmaşıklık ve merge overhead küçük takım için orantısız.
**Consequences:** Feature flag disiplini zorunlu (yarım feature trunk'ta olabilir). Uzun ömürlü branch anti-pattern.
**Tarih:** 2026-09-08

---

## DD-028 — Aynı Artifact Promotion (dev → staging → prod)
**Karar:** Bir git-sha için üretilen Docker image dev'de test edilir, staging'e promote edilir, prod'a promote edilir. Yeniden build edilmez.
**Reason:** Reproducibility, environment parity, güvenlik (staging'de doğrulanan image prod'a değişmez halde gider).
**Alternatives:** Her ortama re-build — reddedildi: parity kaybı, "works on staging, breaks in prod".
**Consequences:** Image tagging strategy (git-sha zorunlu), container registry retention policy.
**Tarih:** 2026-09-08

---

## DD-029 — 12-Factor Config, Env-First
**Karar:** Konfigürasyon 12-factor prensiplerine uygun; primary override kanalı env variable; secret'lar secret manager'dan runtime enjekte.
**Reason:** Immutable image + runtime injection + secret rotation kolaylığı.
**Alternatives:** Config file per env (ansible/chef stili) — reddedildi: rotation zor, image parity bozar.
**Consequences:** Config struct + validation zorunlu (Go boot fail-fast). Ortam davranış farkı için feature flag kullan; `if env == "prod"` yasak.
**Tarih:** 2026-09-08

---

## DD-030 — Structured JSON Logging + Correlation
**Karar:** Tüm server log'ları structured JSON; her satırda `request_id`, `trace_id`, `character_id` (varsa). `msg` alanı sabit; parametreler alan olarak.
**Reason:** Grep/search verimliliği; observability platform native support; audit + debug + monitoring aynı kanaldan.
**Alternatives:** Free-text log — reddedildi: aranabilirlik düşük, PII risk yüksek.
**Consequences:** Log vs audit vs metric ayrı kanallar; PII asla log'a; retention level bazlı (bkz. [logging.md](../architecture/logging.md)).
**Tarih:** 2026-09-08

---

## DD-031 — Monorepo (v1)
**Karar:** Tek repo (`mmo`) altında client/, server/, content/, docs/, deploy/, scripts/. Client-server split Phase 02+ değerlendirmesi.
**Reason:** Küçük takım için cross-domain PR kolaylığı, atomic content+code change, tek CI pipeline.
**Alternatives:** Multi-repo (client + server + content ayrı) — Phase 02 hedef; şu an fragmentation overhead.
**Consequences:** CI'da path-based filter (client-only PR sadece client stage'i çalıştırır, opsiyonel). Monorepo tool (Nx/Turborepo) v1'de gerekmiyor.
**Tarih:** 2026-09-08

---

## DD-032 — Distroless Nonroot Container
**Karar:** Runtime container base = `gcr.io/distroless/base-debian12:nonroot`. No shell, no package manager, non-root user (65532).
**Reason:** Attack surface min, image size < 40 MB, security scan clean.
**Alternatives:**
1. Alpine — reddedildi: musl libc edge case + shell available.
2. Ubuntu slim — reddedildi: gereksiz büyük.
**Consequences:** Debug ephemeral sidecar ile (kubectl debug ephemeral container). Init script yerine Go binary embedded logic.
**Tarih:** 2026-09-08

---

## DD-033 — Schema-Bazlı Namespace (PostgreSQL)
**Karar:** DB objeleri PostgreSQL schema (namespace) ile ayrılır: `auth`, `character`, `content`, `world`, `activity`, `economy`, `social`, `companion`, `progression`, `audit`, `system`.
**Reason:** Ownership + permission granularity + okunabilir query prefix + domain module ile 1:1 eşleşme.
**Alternatives:** Tek `public` schema + prefix'li isim (`auth_user`, `character_stats`) — reddedildi: uzun isim + role management daha zor.
**Consequences:** Migration'da schema önce yaratılır; role permission per-schema; cross-schema JOIN mümkün ancak dikkat.
**Tarih:** 2026-09-08

---

## DD-034 — UUID PK Default, BIGSERIAL Append Tabloları İçin
**Karar:** Çoğu tablo `UUID` PK (v4 başlangıç, v7 hedef); yüksek-throughput append tabloları (ledger, event_log, chat_message, transfer_log, login_history, job_queue) `BIGSERIAL`.
**Reason:** UUID cross-region/shard-safe; ancak sequential id append tabloda index cache locality + storage saving sağlar.
**Alternatives:** Tam UUID everywhere — reddedildi: append tablo write throughput'u için sub-optimal.
**Consequences:** Application tarafında iki id türü ile çalışma; foreign key şeması buna göre.
**Tarih:** 2026-09-08

---

## DD-035 — Template ↔ Instance Ayrımı (Item + Pet + Mount + Node + NPC)
**Karar:** Tüm sahiplenilebilir/spawn-lanabilir entity'ler için `_template` (data) vs `_instance` (state) ikili tablo pattern'i.
**Reason:** Data-driven içerik (DD-006) + oyuncuya özgü state ayrımı; template versionluma karşı instance snapshot hash ile korunur.
**Alternatives:** Denormalize (item satırı tüm alanlarıyla) — reddedildi: content patch migration cehennemi, ekonomi bütünlüğü kaybı.
**Consequences:** Instance tablolarında `template_id` + `template_snapshot_hash`; deprecated template'ler korunur; migration content_version ile.
**Tarih:** 2026-09-08

---

## DD-036 — Append-Only Ledger Discipline
**Karar:** `economy.ledger_entry`, `audit.event_log`, `audit.admin_action`, `audit.security_signal`, `audit.login_history`, `progression.title_earned` — UPDATE/DELETE yasak (role permission). Sadece INSERT.
**Reason:** Forensic audit + reconcile capability + rollback safety (compensating entry). Ekonomik bütünlüğün kalbi.
**Alternatives:** Trigger-based audit (silme dahil) — reddedildi: opaque + trust bound.
**Consequences:** Fiziksel silme yok; compensating entries ile "geri alma". `app_rw` role INSERT-only bu tablolarda; `admin_ro` SELECT için.
**Tarih:** 2026-09-08

---

## DD-037 — Currency Copper Cinsinden (Integer)
**Karar:** Tüm parasal alanlar `BIGINT` ve copper cinsinden. 1 gold = 10000 copper. Float asla kullanılmaz.
**Reason:** Precision safety (float rounding tehlikesi); alt-birim UI'da mikro-işlem gösterimi (silver, copper).
**Alternatives:** DECIMAL(19,4) gold cinsinden — reddedildi: performans + gereksiz precision.
**Consequences:** UI layer copper → gold/silver/copper split rendering. Config'de tüm price/fee copper.
**Tarih:** 2026-09-08

---

## DD-038 — Optimistic Locking (version INT)
**Karar:** Concurrency kritik tablolarda `version INT` sütunu; UPDATE `WHERE version = :expected`; mismatch → retry.
**Reason:** Row-level lock uzun tutmadan concurrency safety. Migration + rolling deploy uyumlu.
**Alternatives:** SELECT FOR UPDATE her yerde — reddedildi: transaction süresi uzar, deadlock riski.
**Consequences:** Uygulama layer'da version increment + retry policy standart. Domain repository interface pattern zorunlu.
**Tarih:** 2026-09-08

---

## DD-039 — Job Queue: Postgres SKIP LOCKED (v1)
**Karar:** v1 background job kuyruğu `system.job_queue` tablosu + `FOR UPDATE SKIP LOCKED` pattern. Redis Streams / NATS Phase 02 karar.
**Reason:** Operasyonel sadelik (aynı DB); transaction safety; küçük ölçekte yeterli throughput.
**Alternatives:** Redis Streams (durability düşük başlangıçta), NATS (harici bağımlılık artışı) — Phase 02'e ertelendi.
**Consequences:** Worker `job_queue` polling loop; heavy load'ta ölçüm sonrası alternatif geçiş planı. Idempotency key + max_attempts + dead letter state.
**Tarih:** 2026-09-08

---

## DD-040 — HTTP Router: chi
**Karar:** Go HTTP router olarak `go-chi/chi` v5 kullanılır.
**Reason:** stdlib-uyumlu (`http.Handler`), minimal, middleware zincirlemesi net, subroute + route grouping ergonomik, mature.
**Alternatives:** `gin` (fast, opinionated, custom context zorluğu), `echo` (middleware model chi'ya benzer, ancak ek katman); `net/http` + custom mux (fazla yeniden icat).
**Consequences:** Handler'lar `http.HandlerFunc` uyumlu; `chi.Router` interface RouteRegistrar pattern'inin temeli. `chi.RouteContext` route pattern metric labelı için kullanılır.
**Tarih:** 2026-09-08

---

## DD-041 — WebSocket Kütüphanesi: nhooyr.io/websocket
**Karar:** WebSocket için `nhooyr.io/websocket` kullanılır.
**Reason:** Modern, context-aware (ctx cancel → conn close), TCP write/read timeouts via ctx, JSON+binary primary API, hafif; net/http entegrasyonu direct.
**Alternatives:** `gorilla/websocket` (mature ama context-native değil, deprecated risk), `gobwas/ws` (low-level, verbose).
**Consequences:** Handler'lar `ctx` passing ile yönetilir; heartbeat/close simply via `ctx.WithTimeout` + `conn.Close(status, reason)`.
**Tarih:** 2026-09-08

---

## DD-042 — Outbox Pattern for Cross-Store Atomicity
**Karar:** Postgres + Redis/notification/webhook birleşimi gereken side-effect'ler için **outbox pattern**: tx içinde `audit.event_log` (v1) veya dedicated `outbox` (v2) tablosuna event insert; commit sonrası worker outbox'ı okuyup ilgili kanallara publish eder. **2PC yok.**
**Reason:** Cross-store atomic işlemi yok; consistency guarantee at-least-once + idempotent handler ile sağlanır. Basit, güvenilir, mature.
**Alternatives:** XA/2PC (Postgres + Redis destek fena), best-effort direct publish (dual-write anomaly riski).
**Consequences:** Notification/chat/webhook subscribers idempotent olmak zorunda. `audit.event_log` v1'de outbox rolü de görür; v2'de ayrı `outbox` tablo değerlendirilir.
**Tarih:** 2026-09-08

---

## DD-043 — Event Bus: In-Process + Redis Pub/Sub (v1)
**Karar:** Domain event bus iki katmanlı — in-process (Go channel) synchronous/async fan-out + Redis pub/sub cross-node. Kalıcı gereksinim outbox worker → publish flow'undan gelir.
**Reason:** v1 için at-most-once yeterli (kalıcı olan zaten outbox); operational simplicity; extra broker (NATS/Kafka) v2 kararı.
**Alternatives:** NATS JetStream v1 — reddedildi: harici bağımlılık ve operational overhead küçük ölçekte orantısız. Redis Streams v1 — mümkün; ancak `PUBLISH` yeterli, upgrade path v2 açık.
**Consequences:** Cross-node güvenilir teslim şart olduğu haller için outbox + worker; direct bus fire-and-forget yalnız non-critical. Event kind registry disiplini gerekli (`docs/backend/event-registry.md` Phase 01).
**Tarih:** 2026-09-08

---

## DD-044 — Cursor-Based Pagination (Offset Yasak)
**Karar:** Tüm listeleme endpoint'leri cursor-based pagination kullanır (`?cursor=<opaque>&limit=N`). Offset yasak.
**Reason:** Deep offset scan performans katili; concurrent insert/update inconsistency yaratır; cursor stable-order + index-friendly.
**Alternatives:** Offset-based (klasik ama scale-hostile), keyset (cursor'ın alt-uygulaması).
**Consequences:** Client tarafında infinite scroll pattern; UI page number kavramı yok. Cursor encoding opaque (base64 struct); server internal.
**Tarih:** 2026-09-08

---

## DD-045 — Access Token: Short-Lived JWT (EdDSA), 15 dk TTL
**Karar:** Access token JWT formatı, EdDSA imzalı (kid rotation destekli, `auth.jwt_key` tablosu), TTL 15 dk. Refresh her ~14 dk client tarafından proactive.
**Reason:** Stateless doğrulama (network hop ↓); revoke etkisi kısa (blocklist entry hızlı expire); EdDSA modern + fast + short signature.
**Alternatives:** Opaque bearer + server-side lookup her request — reddedildi: her request DB/Redis touch, latency ve yük. RS256 — kabul edilir, ancak EdDSA daha modern.
**Consequences:** JWT kid registry + validation cache. Refresh flow (DD-046) short-TTL'i tolere eder.
**Tarih:** 2026-09-08

---

## DD-046 — Refresh Token: Opaque + Rotation + Replay Detection
**Karar:** Refresh token opaque (32 byte random, base64url); server sadece sha256 hash tutar; rotation zorunlu (her refresh yeni token + eski invalidate); reused revoked token gelirse account family compromise → tüm session'lar revoke + security alert.
**Reason:** JWT refresh riski (long TTL leak = uzun süre kötüye kullanım); rotation + replay detection theft ve session hijack'e karşı savunma katmanı.
**Alternatives:** JWT refresh (revoke zor); non-rotating opaque (bir çalınma tüm süre kötüye kullanılabilir).
**Consequences:** Client side race condition için "refresh in-flight" queue pattern zorunlu; concurrent refresh optimistic lock ile birinci kazanır.
**Tarih:** 2026-09-08

---

## DD-047 — Password Hashing: argon2id (m=64MB, t=3, p=1)
**Karar:** Password hash argon2id; parametreler v1 config'te; login sırasında eski param varsa rehash + persist.
**Reason:** Modern; GPU/ASIC-resistant (memory-hard); OWASP önerisi.
**Alternatives:** bcrypt (yaygın ama CPU-only, ASIC-vulnerable), scrypt (argon2id alternatifi; ayarlı ama daha az mainstream), PBKDF2 (yavaş, brute-force karşı zayıf).
**Consequences:** Sunucu CPU/memory pressure (concurrent hash) — semaphore ile bounded; rehash-on-login upgrade path.
**Tarih:** 2026-09-08

---

## DD-048 — Cross-Platform: Aynı Account Sistemi (Web + iOS + Android + Desktop)
**Karar:** Android/iOS/Web/Desktop istemcileri **aynı** account sistemini kullanır: aynı endpoint'ler, aynı token formatı, aynı session tablosu. Cihaza özel: `auth.session.client_kind`, `auth.device`.
**Reason:** Tek MMO, tek karakter, tek ekonomi — hesap parçalanamaz. Cross-device play (mobile başlat, PC devam et) core UX vaadi (bkz. GDD §14).
**Alternatives:** Platform başına ayrı auth (App Store login vs Google Play vs kendi backend) — reddedildi: cross-play imkansız hale gelir.
**Consequences:** IAP/subscription platform-specific ancak entitlement account-level (server-authoritative). Sosyal login (Apple Sign-In gibi) opsiyon olarak eklenir ancak email/password birinci sınıf yol.
**Tarih:** 2026-09-08

---

## DD-049 — Refresh Token Bearer (v1), Cookie Path v2 Opsiyonel
**Karar:** v1'de tüm platformlar refresh token'ı bearer olarak body/header üzerinden gönderir. Secure storage istemci sorumluluğunda (Keychain/Keystore/OS store/memory). Web için HttpOnly cookie yolu Phase 02 opsiyonu.
**Reason:** Tek pattern → daha az edge case + mobile öncelikli hedef.
**Alternatives:** Web cookie zorunlu v1 — reddedildi: CORS + CSRF katmanı erken karmaşıklık.
**Consequences:** Web'de XSS mitigasyon frontend disiplini önemli; refresh token asla `localStorage`'a değil, memory + secure store (opsiyonel).
**Tarih:** 2026-09-08

---

## DD-050 — Enumeration Response Policy
**Karar:** `POST /auth/register` gerçek 409 döner (e-mail taken); `POST /auth/password/reset-request` generic 202 döner (e-mail var/yok fark etmez). Rate limit sıkı.
**Reason:** Register akışı UX kritik ("e-mail taken" bilgisi kullanıcı için değerli, kaydolamayan hesap yaratma tuzağı önlenir); password reset akışı enumeration için düşük değerli (attacker rate limit'e takılır).
**Alternatives:** İkisi de generic — UX kayıp; ikisi de spesifik — reset enumeration çıkarır.
**Consequences:** Rate limit ve audit sinyalleri register endpoint'inde daha sıkı. E-mail change endpoint'i register benzeri.
**Tarih:** 2026-09-08

---

## DD-051 — Password Change / Reset → All Sessions Revoke
**Karar:** Password değişimi (change veya reset) tüm session'ları revoke eder (logout-all). Kullanıcı tüm cihazlarda yeniden login olur.
**Reason:** Password compromise senaryosu: değişim aksiyonu tüm çalınmış oturumları kapatmalı.
**Alternatives:** Sadece current device sürsün — reddedildi: compromise senaryosunda güvenlik boşluğu.
**Consequences:** UX friction (yeniden login) — trade-off kabul; "keep this device" opt-in Phase 02.
**Tarih:** 2026-09-08

---

## DD-052 — XP Curve: Geometric (BASE=100, FACTOR=1.14)
**Karar:** Character XP threshold = `100 * pow(1.14, level-1)`; level 50 cumulative ~190k. Ascension her seferinde threshold `*1.05`.
**Reason:** Erken level'da hızlı ilerleme (retention), geç level'da anlamlı grind (endgame kimliği), ascension'da hafif dampener (dikey ilerleme gücünü rein).
**Alternatives:** Lineer (uzun oyunda motivasyon kaybı), quadratic (erken level yavaş), tablo-based (bakım maliyeti).
**Consequences:** Content pipeline'ın 190k XP'i "1 ascension = 60-100h aktif oyun" hedefiyle uyumlu tune edilmesi gerekli.
**Tarih:** 2026-09-08

---

## DD-053 — XP Debt, Level Down Yok
**Karar:** Ölüm cezası **XP debt** (level 20+, %5 of level XP req), asla level azalması yok. Debt cap = current level XP req; next grants önce debt öder.
**Reason:** Level down demoralize eder; XP debt "geri adım gibi ama level'ı korur" — casual friendly + hala anlamlı ceza.
**Alternatives:** Level down (Everquest), no penalty (WoW retail modern) — reddedildi: ilki brutal, ikincisi anlamsız risk.
**Consequences:** `character.character.xp_debt` alanı schema eklemesi.
**Tarih:** 2026-09-08

---

## DD-054 — Stamina: v1'de Etkin Değil (Soft Cap Yeterli)
**Karar:** Stamina sistemi tasarlandı (`docs/character/stamina.md`), ancak v1'de **aktive edilmedi**. AFK Engine soft cap (DD-014) idle rate limiting işlevini zaten karşılıyor.
**Reason:** Persistent Progression pillar'ı (DD-003) ile stamina "wall" hissi çelişme riski. Soft cap yumuşak azalma yeterli.
**Alternatives:** v1'de aktive et (idle MMO'da ritim yaratır ancak "wait wall" hissettirebilir), asla ekleme (schema olmadan ileride ekleme migration maliyeti).
**Consequences:** Schema hazır (Phase 02 aktivasyon opsiyonel). Balance simülasyonu Phase 02'de yapılacak.
**Tarih:** 2026-09-08

---

## DD-055 — HP Ephemeral Redis, Snapshot Postgres
**Karar:** Current HP/Mana/Resource runtime state Redis'te (`char:hp:<char_id>`); save trigger'larda Postgres snapshot'a yazılır (session end / activity change / death / heartbeat 5 dk).
**Reason:** Combat tick 250ms; per-tick Postgres write imkansız throughput. Redis kaybı: snapshot'tan restore (biraz eski ama safe).
**Alternatives:** Full Postgres (write amplification), full Redis (data loss risk).
**Consequences:** Redis-Postgres write-back queue; snapshot version alanı; Redis flush job.
**Tarih:** 2026-09-08

---

## DD-056 — Server-First Title Kilidi Recycle Yasak
**Karar:** Server-first title'lar (`content.title.is_server_first=true`) bir kez sahiplenildikten sonra `unique_holder_character_id` sabitlenir. Sahiplenmiş karakter silinse bile title recyclelanmaz.
**Reason:** Prestij bütünlüğü; "silip yeniden almak" exploit'i önlenir; kalıcı MMO folklore.
**Alternatives:** Recycle (yeniden karakter → yeni sahiplenme) — reddedildi: değeri erimeli değil.
**Consequences:** Silinmiş karakterin server-first title'ı "unclaimable" statüsünde audit'te kalır; UI'de tarihçe gösterebilir.
**Tarih:** 2026-09-08

---

## DD-057 — Race Data-Driven, Effect DSL Whitelisted
**Karar:** Race tanımı YAML content pack ile eklenir. Passive/active davranışlar JSON-encoded effect list (kind + target + amount + scope). Effect combiner çözer. Kod tarafında race code'u referansı YASAK (switch/if in domain code).
**Reason:** İçerik pipeline'ı hızlandırır; yeni race eklemek migration + kod değişikliği değil sadece content pack. Whitelisted target/kind set exploit önler.
**Alternatives:** Full-code race (hard-coded switch) — reddedildi: DD-006 data-driven prensibiyle çelişir. Sözcüksüz freeform DSL — reddedildi: güvenlik ve determinism kaybı.
**Consequences:** Effect target whitelist maintenance gerekli (yeni derived stat eklendiğinde eklenmeli). Content loader validation zorunlu (build-time CI check).
**Tarih:** 2026-09-08

---

## DD-058 — Race Change Kilit v1'de Kapalı
**Karar:** Race change endpoint `POST /v1/characters/{id}/race-change` implementasyonu yapılır ancak feature flag `race_change_enabled=false`. Karakter yaratımında seçilen race kalıcı.
**Reason:** Race kimlik bağlamı taahhüde bağlıdır (DD-005 kimlik önceliği); race change ekonomik/dengeleyici karar Phase 02 balancing pass sonrası.
**Alternatives:** v1'de aktif (kimlik erimesi riski), tamamen yasak (Phase 02'de nerf compensation seçenek kaybı).
**Consequences:** Nerf senaryosu için "compensation token" (1x free race change 30 gün TTL) worker altyapısı hazır.
**Tarih:** 2026-09-08

---

## DD-059 — Unique Mechanic: Handler-Registered (Yaml Referanslı Kod)
**Karar:** Race'in unique mechanic'i (`dwarven_dark_vision`, `merfolk_water_travel` gibi) YAML'da kod referansı olarak durur; ancak gerçek davranış Go handler tablosunda kayıtlı (`racemechanic.handlers`). Yeni unique mechanic ekleme = handler PR + yaml.
**Reason:** Effect DSL'in mantıklı sınırı var; niş "custom" davranışları DSL'e sığdırmak DSL karmaşıklığı ↑. Handler kayıt tablosu güvenli + testable.
**Alternatives:** Tam data-driven DSL (interpreter overhead + security surface), tam kod (data-driven prensibi tamamen ihlal).
**Consequences:** Yeni race ekleme: eğer sadece stat/passive/active + affinity kullanıyorsa **kod yok**. Custom mechanic isteniyorsa **handler PR** gerekir. Content loader unknown mechanic code → build fail.
**Tarih:** 2026-09-08

---

## DD-060 — Race Slot Bar Konumu (11 — Reserved)
**Karar:** Character skill bar'ının 11. slotu (0-indexed, son slot) racial active için reserved default. Oyuncu istersinse bu slotu boşaltıp başka skill koyabilir (kişisel tercih).
**Reason:** Racial ability'nin oyuncu tarafından "yerini bulmak zor" olmaması; öntanımlı erişim.
**Alternatives:** Auto slot (herhangi bir boş), no slot (oyuncu manuel yerleştirir) — reddedildi: UX friction.
**Consequences:** Character creation transaction'ında slot 11 insert. Oyuncu customize edebilir.
**Tarih:** 2026-09-08

---

## DD-061 — Class + Specialization Data-Driven (YAML + Effect DSL)
**Karar:** Class ve specialization tanımı YAML content pack ile. Passive/active/signature skill davranışları content.skill_template + effect DSL üzerinden (DD-057 race sistemiyle aynı prensip). Domain kod'unda class code referansı YASAK.
**Reason:** Yeni class/spec eklemek migration + kod PR değil sadece content pack. Rework iterasyonu hızlı. Kullanıcı revize edeceği için esneklik kritik.
**Alternatives:** Hard-coded class enum + go switch — reddedildi: iterasyon yavaş, DD-006 ile çelişir.
**Consequences:** Content loader zorunlu validation (skill code varlığı, tree constraint, stat budget). Effect DSL target whitelist maintenance gerekli.
**Tarih:** 2026-09-08

---

## DD-062 — Spec Tree Kesin: 2 Tier I × 3 Tier II per Class
**Karar:** Her base class tam olarak **2 Spec I + 6 Spec II** taşır (her Tier I'in 3 Tier II çocuğu). Content loader bu constraint'i zorlar.
**Reason:** Tutarlı player deneyimi (her class'ın 6 leaf build'i olur → 6 base × 6 = 36 unique identity). UI predictable.
**Alternatives:** Esnek (2-4 seçenek) — reddedildi: UI karmaşıklığı + balance zorluğu. Sabit tek Spec II per Tier I — reddedildi: derinlik kaybı.
**Consequences:** Yeni class eklemek 2+6=8 spec + skill'ler gerektirir (content bar yüksek). Phase 02'de esnek '2-4' opsiyonu tekrar değerlendirilir.
**Tarih:** 2026-09-08

---

## DD-063 — Class Değişmez, Sadece Spec Respec
**Karar:** Class kalıcı (DD-004 reaffirm). Yalnız specialization respec ile değişir. Spec I değişince Spec II otomatik reset.
**Reason:** Class kimlik taahhüdü + oyun içi kimlik bağı; Spec derinlik esnekliği.
**Alternatives:** Class change ücretli (Phase 03+ tartışılabilir; şu an hayır), dual class hybrid (Phase 03+ değerlendirme).
**Consequences:** Class rework durumu için free respec token altyapısı yeter; class değiştirme workflow yok.
**Tarih:** 2026-09-08

---

## DD-064 — Respec Cost Curve: Exponential + Cooldown
**Karar:** İlk respec ücretsiz. Sonraki: Spec II `500 * 2^n copper`, Spec I `2000 * 2^n copper`; cap 2M copper (20k gold). Cooldown: Spec II 24h, Spec I 72h.
**Reason:** Denemeye açık ilk seçim; sonraki tekrar iteration ekonomiye sink + spam engel. Cooldown paylaşımı Spec I sonrası Spec II lock.
**Alternatives:** Sabit maliyet (ekonomiye zayıf katkı), tamamen ücretsiz (respec spam), token-only monetization (paywall algısı).
**Consequences:** `character.respec_count`, `respec_cooldown_until`, `free_respec_tokens` schema. Content rework free token grant workflow. Placement Phase 01 balance pass'inde revize edilebilir.
**Tarih:** 2026-09-08

---

## DD-065 — Placeholder Class + Spec İçeriği v1
**Karar:** Warrior/Ranger/Mage/Cleric/Rogue/Druid — 6 base class placeholder. Warrior için user'ın verdiği isimler kullanılır (Berserker/Guardian → Ravager/Warlord/Blood Reaver + Guardian tier II TBD). Diğer sınıflar için placeholder tree hazır ancak final isimler ve mekanikler **kullanıcı tarafından revize edilecek**.
**Reason:** Content pipeline test edilebilmesi için concrete content gerekli; ancak class kimliği subjektif tercih. İlk taslak çerçeve.
**Alternatives:** Boş bırak (pipeline test edilemez), tam final yap (kullanıcı revizyon opsiyonu kaybolur).
**Consequences:** [class/placeholder-content.md](../class/placeholder-content.md) tüm placeholder içerik + revize checklist. Kullanıcı final ile revizyon PR'ı gelene kadar isim/mekanik freeze YOK.
**Tarih:** 2026-09-08

---

## DD-066 — Skill + Stat Engine Combat'tan Bağımsız
**Karar:** SkillStatEngine combat'ın parçası değildir; combat, AFK sim, aktivite yield hesabı, vendor pricing — hepsi aynı engine sözleşmesini kullanır (`engine.ComputeDerived`, `engine.Cast`, `engine.ApplyBuff`).
**Reason:** DRY + tutarlılık: aynı hesap her yerde aynı sonuç. AFK reveal ile manuel play arasında determinism köprüsü. Test kolaylığı (tek golden replay suite).
**Alternatives:** Combat-özel engine + ayrı gathering/afk hesabı — reddedildi: drift + double bookkeeping.
**Consequences:** `internal/platform/engine` (veya `internal/engine`) paketi; her domain modülü buna interface üzerinden bağlı. Combat modülü engine'in "en büyük tüketicisi", "sahibi" değil.
**Tarih:** 2026-09-08

---

## DD-067 — Modifier Kind Set Sabit + Apply Order Kesin
**Karar:** Modifier kind set kesin: `flat / mult_add / mult_prod / set / min / max`. Apply order: flat → mult_add (tek çarpım) → mult_prod (her biri ayrı çarpım) → set → clamp.
**Reason:** Formul öngörülebilirliği; content team modifier davranışını sim edilebilir; balance debugging kolaylaşır.
**Alternatives:** Serbest priority (chaos), sadece additive (esneklik kaybı).
**Consequences:** Content loader modifier kind whitelist enforce; test golden replay bu order'a bağlı; yeni kind eklemek CI-controlled.
**Tarih:** 2026-09-08

---

## DD-068 — Cooldown Redis-Only, DB Persist Yok
**Karar:** Skill CD, GCD, category CD, ICD (passive proc) — hepsi Redis TTL-based; Postgres persist yok. Redis kayıp durumunda tüm CD sıfırlanmış olur (rare avantaj); DB'de "last_cast_at" tutulmaz.
**Reason:** Latency (Redis micro-sec vs Postgres ms); simplicity (TTL doğal expiry); combat tick 250ms constraint.
**Alternatives:** Postgres persist + write-through cache — reddedildi: write amplification, gereksiz complexity.
**Consequences:** Dungeon lockout, daily quest CD gibi uzun-süreli cooldown'lar ayrı domain (progression schema); engine sadece ephemeral combat CD.
**Tarih:** 2026-09-08

---

## DD-069 — CDR Cast Anında Uygulanır (Snapshot)
**Karar:** Cooldown reduction (CDR) skill cast anında hesaplanır; Redis'e giden `ms` değeri CDR sonrası snapshot. Cast sonrası CDR değişse bile mevcut CD etkilenmez.
**Reason:** Fair play + öngörülebilirlik; buff pop → cooldown "bites" edilemez.
**Alternatives:** Runtime hesap (buff aktifken remaining_ms değişir) — reddedildi: exploit fırsatı, complex UI.
**Consequences:** CDR pop timing önemli (cast öncesi buff aç); rotation editor bu bilgiyi UI'da açıklar.
**Tarih:** 2026-09-08

---

## DD-070 — DoT/HoT: Snapshot (v1)
**Karar:** DoT/HoT effect'ler buff apply anında caster stat'larını **kilitler** (snapshot). Her tick fresh recompute değil.
**Reason:** Combat determinism + performance (tick başına recompute overhead); simplicity.
**Alternatives:** Re-calc per tick — Phase 02 experiment; buff enfoque durumuna daha realistic ancak complexity.
**Consequences:** DoT tick amount `character.stats_buff` içinde persist; caster stat sonradan değişse bile aynı damage.
**Tarih:** 2026-09-08

---

## DD-071 — Skill Rank v1'de max_rank=1 Sabit, Sistem Aktif
**Karar:** Skill rank sistemi engine tarafında implement edilir (`character.skill_learn.rank` + XP + rank up flow); ancak v1 content'te tüm skill'ler `max_rank=1`. Phase 02 content patch rank 5'e çıkarır.
**Reason:** Engine değişikliği Phase 02'de gerekmez (sadece content); erken doğru mimari.
**Alternatives:** v1'de tam olarak devre dışı (Phase 02'de tam yeniden yaz) — reddedildi: yeniden yaz maliyeti.
**Consequences:** Skill rank up event'leri v1'de tetiklenmez (max_rank=1); ancak XP accumulate eder ve Phase 02 content patch ile rank up başlar.
**Tarih:** 2026-09-08

---

## DD-072 — Item Engine: Ortak Create API (Loot + Craft + Quest + Admin)
**Karar:** Item instance yaratımı tek `ItemEngine.Create(req CreateRequest)` çağrısıyla; source (loot/craft/quest/admin/event) parametre olarak. Deterministic seed her drop için.
**Reason:** DRY + tutarlılık (aynı roll formulü). Test golden replay tek sistem. Audit trail birleşik.
**Alternatives:** Source başına ayrı factory — reddedildi: drift + double bookkeeping.
**Consequences:** Tüm domain modülleri (combat drop, crafting output, quest reward, admin grant) `item.Engine.Create` çağırır.
**Tarih:** 2026-09-08

---

## DD-073 — Rarity Distribution 7 Band + Multiplier Curve
**Karar:** Rarity band'ları: poor / common / uncommon / rare / epic / legendary / mythic. Baz mob drop dağılımı: 5/60/25/8/1.8/0.18/0.02 %. Rarity multiplier: poor 0.6× → mythic 1.75×.
**Reason:** Loot doyumu curve (common çok, rare heyecan, legendary/mythic hayal). Multiplier ile aynı template yüksek rarity'de belirgin fark.
**Alternatives:** 5 band (basit), 10 band (fragmentation).
**Consequences:** Content team drop distribution loot table meta'da override edebilir. Bad luck protection Phase 02 opsiyonu.
**Tarih:** 2026-09-08

---

## DD-074 — Affix Budget Rarity-Driven, Prefix + Suffix Ayrı
**Karar:** Affix count: common 0-1, uncommon 1-2, rare 2-3, epic 3-4, legendary 4, mythic 4-5. Prefix + suffix max 1 each; ekstra affix rasgele tarafa. Mutual exclusion tag-based.
**Reason:** Öngörülebilir tooltip; kimlik ayrımı prefix (bonus stat) vs suffix (utility). Tag exclusion ile çift STR bonus çakışması önlenir.
**Alternatives:** Serbest affix sayısı, kategori ayrımı yok.
**Consequences:** Content pool tasarımı rarity-per-slot ayrılabilir; tag whitelist maintained.
**Tarih:** 2026-09-08

---

## DD-075 — Socket + Gem: Removable Pattern (v1: Color=Any)
**Karar:** Socketed gem kalıcı **kaybolmaz** — remove endpoint ile gem geri envanter'e döner (small gold cost). v1'de color coding devre dışı (type=any); Phase 02'de renk sistemi opsiyonel.
**Reason:** Player deneyimi (yanlış gem tak, geri al); rework friction düşük. Color sistem Phase 02'de eklenmesi opsiyon.
**Alternatives:** Consumable gem (once socketed, destroyed). Reddedildi: sinir bozucu, item churn yüksek.
**Consequences:** Gem instance `owner_kind='socketed_in_item'` net ownership. Remove cost economy sink küçük.
**Tarih:** 2026-09-08

---

## DD-076 — Enchant Scroll Model (Direct-Apply Değil)
**Karar:** Enchanter item yerine **scroll** üretir; item'a scroll consume ile apply edilir. Scroll BoE (trade edilebilir).
**Reason:** Enchanting ekonomiye entegre (scroll market'te satılır); item bound olsa da enchant yeteneği ekonomiden gelir.
**Alternatives:** Direct enchanter apply (workbench recipe), reagent + apply.
**Consequences:** Enchant scroll ayrı item template category='consumable', subcategory='enchant_scroll'. Enchanter recipes scroll'lar için.
**Tarih:** 2026-09-08

---

## DD-077 — Set Bonus Tiers 2/4/6 Kesin, Trigger 6p
**Karar:** Set bonuses 3 tier: 2 piece (konfor bonus), 4 piece (sinergi bonus), 6 piece (trigger-based signature). Ring1+Ring2 aynı set 2 sayılır.
**Reason:** Standard MMO pattern; player için öngörülebilir tier hedefi.
**Alternatives:** Free tier count (2/3/5/6), sadece 2/4.
**Consequences:** Content loader `pieces_required` in [2, 4, 6] enforce; Phase 02'de esneklik değerlendirme.
**Tarih:** 2026-09-08

---

## DD-078 — Unique Effect Handler-Registered + Legendary Cap 2
**Karar:** Legendary+ item'ların "unique effect"i Go handler registry ile (DD-059 race pattern paralel). Karakter eş zamanlı **max 2 legendary+** item takabilir.
**Reason:** Custom mekanikleri effect DSL'in kaldıramadığı derinlik gerektirir; handler + params combo. Cap ile combo OP riski sınırlı.
**Alternatives:** Effect DSL genişletme (interpreter risk), no cap (broken combos).
**Consequences:** Yeni legendary item ekleme handler PR + YAML. Loader unknown handler_code → build fail. Cap konfigüre edilebilir.
**Tarih:** 2026-09-08

---

## DD-079 — Storage Ortak Abstraction (7 Kind)
**Karar:** Inventory, Bank, Equipment, GuildBank, Mailbox, MarketListing, Void — hepsi ortak `Storage` interface implementasyonu. Access control her kind kendi `CanAccess` handler'ında.
**Reason:** Item transfer, capacity, listing UI, migration — hepsi tek pattern; drift önlenir; test surface küçük.
**Alternatives:** Storage başına özel API — reddedildi: kod duplication + drift.
**Consequences:** Yeni storage tipi eklemek (Phase 02 Loadout, cross-account) interface implementation.
**Tarih:** 2026-09-08

---

## DD-080 — Inventory Overflow → Mailbox (7 Gün TTL)
**Karar:** Envanter dolu iken rare/BoP drop varsa item **system mail**'e overflow (`sender_kind='system'`, 7 gün TTL); expire olursa void (audit).
**Reason:** Loot heyecanı korunur (drop kaybolmaz); envanter yönetimi disiplini teşvik edilir.
**Alternatives:** Drop refuse (kayıp), ground drop (v1'de pozisyonel değil, uygulanabilir değil).
**Consequences:** Mailbox cap 100 → aşılırsa hard-drop with alarm (rare senaryo); worker overflow monitor.
**Tarih:** 2026-09-08

---

## DD-081 — Bank Access Location-Gated
**Karar:** Character bank sadece şehir/köy'lerin `services.bank == true` POI'lerinde erişilir. Herhangi bir yerden bank açılamaz.
**Reason:** Şehir ziyaretini teşvik + gathering/travel ekonomik ritim; safe zone hissi.
**Alternatives:** Cross-map bank (WoW modern), phone bank (Phase 03 konfor).
**Consequences:** Content'te tüm başlangıç şehirleri bank servisiyle donatılır; outpost'larda opsiyonel bank.
**Tarih:** 2026-09-08

---

## DD-082 — Transfer Atomicity: SERIALIZABLE + item_transfer_log Zorunlu
**Karar:** Cross-storage item transfer'i SERIALIZABLE tx isolation'da; her transfer `economy.item_transfer_log`'a satır (append-only). Idempotency-Key header'ı ile duplicate safe.
**Reason:** Ekonomi bütünlüğü kritik; double-transfer/duplication exploit sıfır tolerans.
**Alternatives:** Read Committed (yetersiz race window), event sourcing full (over-engineering v1).
**Consequences:** Retry policy 40001 için 3× jitter; reconcile worker daily (item_instance ownership vs storage row consistency).
**Tarih:** 2026-09-08

---

## DD-083 — Item Lock Manuel + Auto-Lock on Pickup Content-Driven
**Karar:** Item lock manuel; içerik `content.item_template.auto_lock_on_pickup=true` ile drop anında lock (unique legendary, quest item için). Locked item: sell/salvage/destroy/mail/list bloklu; unlock manuel.
**Reason:** Accidental sell/destroy kayıp kritik; unique/first-time drop koruma. Manuel toggle basit UX.
**Alternatives:** Server-side rule engine (auto-lock rarity threshold Phase 02), no lock (kayıp riski).
**Consequences:** Client tooltip badge (🔒); event `item.action_blocked_by_lock` audit; Phase 02 rule engine üzerine inşa edilir.
**Tarih:** 2026-09-08

---

## DD-084 — NEVER Simulate Every Offline Tick (Altın Kural)
**Karar:** Activity/AFK Engine per-tick simülasyon YAPMAZ. Tüm yield analitik hesap (rate × elapsed = expected) + aggregate distribution sample (Poisson/Binomial/Normal) + threshold detection (linear inverse) ile O(1) çözer.
**Reason:** Ölçek (100k char × 8h = 2.88B tick imkansız), latency (login p50 < 300ms hedefi), determinism (test + audit + replay), throughput (≥5000 char/sn reveal hedefi). Per-tick loop mimarinin ölüm sebebi olur.
**Alternatives:** Simulate loop küçük N için (100 iteration cap) — reddedildi: kod complexity + drift + karışıklık. Full-analytical + no variance — reddedildi: dopamine kaybı (fixed yield boring).
**Consequences:** Content team ve Claude Code direktif: aktivite implementer bu prensibe uymalı; CI lint (grep `for.*elapsed.*tick`) engelleme. Aktivite rate compute pure function; test golden replay.
**Tarih:** 2026-09-08

---

## DD-085 — Aggregate Variance: Poisson (Gathering/Kill) + Binomial (Craft/Drop)
**Karar:** RNG variance için tek distribution sample:
- **Poisson** for gathering yield, kill count, rare proc — parametre λ = rate × elapsed.
- **Binomial** for craft success rate, drop chance — parametre (n, p).
- **Normal approximation** büyük parametre için (λ > 20 veya n > 20).
- **Deterministic mean only** çok büyük parametre için (λ > 1000).
**Reason:** Distribution sample = O(1) veya O(√n) küçük N için, per-tick roll = O(N). Mathematically equivalent expected value; variance realistic.
**Alternatives:** Only-mean deterministic (variance yok, boring), per-tick RNG (imkansız scale).
**Consequences:** RNG stream seed'li sub-stream; determinism test golden replay. Cap değerler config'te.
**Tarih:** 2026-09-08

---

## DD-086 — Threshold Detection: Linear Inverse (Analytical Solve)
**Karar:** Inventory-full, stamina-zero, consumable-out, death-cascade gibi threshold'lar iterasyon değil analitik çözüm:
`ms_until_threshold = current_capacity / rate_per_ms × 3600000`.
En küçük threshold `HitAtMs` kazanır; `elapsed_effective = HitAtMs`; sonra yield hesabı.
**Reason:** Threshold hit anını simulate etmeden formülle bulmak lineer time. Reveal karmaşıklığı düşük tutulur.
**Alternatives:** Binary search over simulated state (log N ancak simulate() bağımlı), full simulate (yasak).
**Consequences:** Rate constant assumption during window (rate mid-contract değişmez, snapshot immutable — kararlı hesap). Non-linear threshold (death cascade) approximation.
**Tarih:** 2026-09-08

---

## DD-087 — Reveal Atomic + Idempotent + reveal_log Zorunlu
**Karar:** Reveal → transaction: snapshot compute yield → dispatch (ledger + inventory + xp) → contract status update → reveal_log insert. Aynı `idempotency_key` unique → duplicate no-op safe. `activity.reveal_log` her reveal için satır (append-only).
**Reason:** Ekonomi bütünlüğü + double-grant prevention + retry safe. Reveal audit-forensic.
**Alternatives:** Best-effort dispatch (drift), external ledger sync — reddedildi.
**Consequences:** Reveal handler retry-safe; client de POST /reveal idempotent. reveal_log partition (Phase 02) storage optimization.
**Tarih:** 2026-09-08

---

## DD-088 — Snapshot Contract Lifetime, Not Per-Reveal
**Karar:** Contract snapshot bir kez alınır (contract create anında); reveal sırasında bu snapshot kullanılır. Karakter equipment/stat değişse bile aktif contract snapshot güncellenmez. Yeni contract oluştuğunda fresh snapshot.
**Reason:** Reveal determinism (aynı contract + elapsed → aynı yield); mid-contract stat drift işlem karmaşıklığı yaratmaz.
**Alternatives:** Per-reveal fresh snapshot — reddedildi: reveal determinism kaybı, mid-contract exploits.
**Consequences:** Karakter equipment değiştirirse yield yeni gear'la olmaz **contract sonuna kadar**. Kullanıcı yeni contract başlatabilir (auto-close). UX tooltip açıklaması.
**Tarih:** 2026-09-08

---

## DD-089 — Dungeon Iki Mod: Manuel Instance + Idle Contract (Yield ×0.5)
**Karar:** Dungeon iki tür — manuel (activity.instance, party, aktif oyun), idle-run (activity.contract, primary channel, AFK). Idle-run yield manuel'in ~%50'si.
**Reason:** Casual friendly (idle path); competitive/full loot manuel path; weekly lockout paylaşımlı.
**Alternatives:** Sadece manuel (idle MMO değil), sadece idle (dungeon deneyimi kaybı).
**Consequences:** İki lifecycle (instance vs contract); loot table paylaşımlı; weekly lockout single counter.
**Tarih:** 2026-09-08

---

## DD-090 — Location Parent Polimorfizmi (region/territory/settlement)
**Karar:** `world.location.parent_id + parent_kind` polymorphic; location parent'ı üç seviye olabilir (wilderness → region parent; urban → settlement parent).
**Reason:** Esneklik — bir POI şehir içi veya wilderness'ta olabilir; hiyerarşiyi zorunlu tek yol yapmak (region → territory → settlement → location) esnekliği kısıtlar.
**Alternatives:** Tek parent tip (zorunlu region) — reddedildi: settlement içi POI (shop, quest giver) modellemesi zor.
**Consequences:** FK enforcement application layer; migration + query filter dikkat. Loader validation zorunlu.
**Tarih:** 2026-09-08

---

## DD-091 — Personal Resource Node Instance via Redis TTL
**Karar:** Resource node depletion **personal** — karakter perspektifinden. Server tarafında tek shared instance row (`world.resource_node`), ancak karakter-özel "yielded_at → next_available_at" Redis'te (`char:node:<char>:<node>` TTL respawn_ms). Rare node opsiyonel Phase 02 global depletion.
**Reason:** DD-022 (personal instance ilkesi); idle MMO gathering competition frustration önlenir; performance (Redis TTL vs DB query per gather).
**Alternatives:** Global shared depletion (competition MMO klasik; reddedildi: idle-first design contra), full DB persist (write amplification).
**Consequences:** Redis kayıp = node available (avantaj rare); reveal logic Redis-safe fallback (analytical rate compute node availability düşünmez — yield perspektifinden).
**Tarih:** 2026-09-08

---

## DD-092 — Monster Area Spawn Personal (v1)
**Karar:** Monster area mob spawn **personal** — her karakter kendi spawn hızını görür (Redis-based). Shared spawn Phase 02 opsiyonu (competition modu).
**Reason:** AFK hunting yield determinism + fairness. Aynı spot'ta 100 karakter olsa da bireysel yield etkilenmez.
**Alternatives:** Shared spawn pool (competitive MMO); reddedildi: idle-first vizyona aykırı.
**Consequences:** Rare mob shared (nadir spawn, ilk vuran alır); world boss dedicated system.
**Tarih:** 2026-09-08

---

## DD-093 — Content World YAML Pipeline (Continent/Region/Territory/Settlement/Location)
**Karar:** Dünya içeriği YAML content pack'te; `content/world/{continents|regions|territories|settlements|locations|...}/*.yaml`. Loader → DB upsert → `content_version` bump → hot reload.
**Reason:** DD-006 data-driven ilkesi ile tutarlı; içerik ekleme migration + kod PR değil sadece YAML.
**Alternatives:** DB-first editing (admin panel gerekir; erken evre için over-kill), hard-coded (esneklik yok).
**Consequences:** CI content validation: cross-ref (region.continent_code exists, etc.) zorunlu; deploy pipeline content build stage.
**Tarih:** 2026-09-08

---

## DD-094 — Travel Efektif Süre Multiplikatif Formül
**Karar:** Travel duration formula: `duration = (distance × terrain × 1000) / (base_speed × mount × race × equipment × road × weather)`. Katmanlar çarpımsal (multiplicative); toplam boost cap 3.5×; terrain penalty cap dışı.
**Reason:** Predictable UI (her katman kendi bonus'unu tanımlı); segment-based hesap; balance tune noktaları izole.
**Alternatives:** Additive (kolay ama flat curve), sabit tier'lar (esneklik yok).
**Consequences:** Content pack her katman için multipliers config'te; content team balance tuning.
**Tarih:** 2026-09-08

---

## DD-095 — Travel Snapshot Contract Lifetime
**Karar:** Travel contract snapshot'ında race/mount/equipment/buff dondurulur; travel süresince değişimler yok sayılır. Buff bitse veya item değişse duration korunur.
**Reason:** DD-088 (activity snapshot pattern) ile tutarlı. Deterministic reveal.
**Alternatives:** Dynamic mid-travel re-compute — reddedildi: complexity + non-deterministic reveal.
**Consequences:** UX: kullanıcı "yeni gear taktım hız artacak" beklerken travel bitene kadar eski hız. Yeni travel yeni snapshot.
**Tarih:** 2026-09-08

---

## DD-096 — Distance Model: Location Adjacency Graph + Dijkstra
**Karar:** v1'de positional coord yok; distance = `world.location_adjacency` graph üzerinden Dijkstra pathfind. Cache Redis 5 dk.
**Reason:** Abstract simplicity; content team edge-by-edge tanımlar; hesap deterministic + cache-friendly.
**Alternatives:** Positional coord + Euclidean — Phase 02 opsiyonu; v1 için over-engineering.
**Consequences:** Content pack adjacency YAML'lar zorunlu; missing edge = unreachable. Route cache prefix content_version bump.
**Tarih:** 2026-09-08

---

## DD-097 — Weather v1'de Kapalı (Phase 02 Aktivasyon)
**Karar:** Weather sistem tasarımı hazır (`world.region_weather_state` schema) ancak v1'de aktif değil (weather_multiplier=1.0 sabit).
**Reason:** Complexity ekler (worker, UI, content, race resist table); v1 core loop öncelik. Phase 02 open-world zenginlik için hazır.
**Alternatives:** v1'de aktive et — reddedildi: playtest önceliği farklı.
**Consequences:** Formul yeri hazır; enable flag ile aktive. Weather-race resist table Phase 02 content.
**Tarih:** 2026-09-08

---

## DD-098 — v1 Profession Seti: 6 Aktif (Mining/Woodcutting/Fishing/Hunting/Herbalism/Skinning); Farming + Archaeology Phase 02
**Karar:** v1'de 6 gathering profession aktif: mining, woodcutting, fishing, hunting, herbalism, skinning. Farming ve archaeology **Phase 02** (schema hazır, aktivasyon sonra). Crafting profession set (blacksmithing, leatherworking, tailoring, alchemy, cooking, enchanting, jewelcrafting) ayrı faz.
**Reason:** v1 core loop odak; farming housing bağlı (Phase 02), archaeology content-heavy (Phase 02). 6 aktif profession balance + ekonomi test için yeter.
**Alternatives:** All 8 v1 (content overload), sadece 3-4 (ekonomi dar), no gathering v1 (aktif oyun eksik).
**Consequences:** 6 profession × N node/mob = ~40 gathering node type v1 content ihtiyaç. 3 slot cap oyuncu için tercih zorlar.
**Tarih:** 2026-09-08

---

## DD-099 — Skinning: Standalone Aktivite Değil, Hunting Kill Yan-Effect
**Karar:** Skinning ayrı aktivite başlatılmaz; hunting kill'inden pasif olarak drop üretir (skinning profession aktif ise). Kill başına 3 skinning XP baseline; family-based loot table.
**Reason:** Simplicity: her hunting aktivitesi ikili yield (combat loot + skinning materials); ayrı contract yönetimi gereksiz.
**Alternatives:** Standalone skinning (kill sonrası carcass process) — reddedildi: workflow karmaşıklığı, mekanik değer düşük.
**Consequences:** Skinning aktivite tipi yok; profession XP hunting yield içinde hesaplanır; content skinning_tables per mob family.
**Tarih:** 2026-09-08

---

---

## DD-100 — Crafting Quality: Poor/Normal/Exceptional (v1)
**Karar:** Crafting output 3 quality tier (poor/normal/exceptional). Mastery + tool + recipe quality_bias distribution etkiler. Exceptional +1 rarity + stat ×1.2; Poor -1 rarity + ×0.7. Item instance'a `craft_quality` alanı.
**Reason:** Craft'in "same output every time" boringlığını kırar; mastery + luck ödüllendirir; ekonomiye exceptional premium item.
**Alternatives:** Sadece binary (success/fail), zengin tier (poor→masterwork→legendary craft) — v1 3 tier basit; Phase 02 expand.
**Consequences:** UI tooltip badge + rarity glow; marketplace filter Phase 02.
**Tarih:** 2026-09-08

---

## DD-101 — Critical Craft ×2 Qty + Quality Bump (Kombine)
**Karar:** Critical iteration output: qty ×2 + byproduct guaranteed + quality tier +1 (kombine). Base %5 crit; mastery + tool ile cap %30.
**Reason:** Loot moment (dopamine); crafter için "big score" hissi; mastery grind incentivize.
**Alternatives:** Sadece qty (boring quality), sadece quality (qty flat). Kombine v1 dopamine max.
**Consequences:** Item instance `craft_critical` boolean; UI "Crit Craft" badge (Phase 02); event `crafting.critical_craft`.
**Tarih:** 2026-09-08

---

## DD-102 — Crafting XP Fail Penalty %50 + Big Miss Zero XP
**Karar:** Fail iteration → %50 XP (character/profession/mastery). Big miss (%1 nadir) → zero XP + full material loss + notification.
**Reason:** Fail totally void olmasın (öğrenme değeri); big miss dramatik olay (nadir; balance için mastery koruma).
**Alternatives:** Fail zero XP (harsh), full XP fail (no penalty), no big miss (drama kaybı). v1 kombine.
**Consequences:** Community feedback tolerance (big miss frustration); mastery-based reduction cap 0.5% floor.
**Tarih:** 2026-09-08

---

## DD-103 — Crafting AFK Analytical: Binomial Success + Multinomial Quality + Binomial Critical
**Karar:** Batch crafting yield **DD-084 uyumlu** — per-iteration loop yok. Success: binomial(iterations, craft_success_pct). Quality distribution: multinomial(successCount, {poor/normal/exceptional}). Critical: binomial(successCount, critical_pct). Threshold: linear inverse (input exhaust, output inv full).
**Reason:** 200 iteration = 200 loop yerine 3 sample; determinism seed'li; O(1) reveal.
**Alternatives:** Per-iteration loop — DD-084 ihlali.
**Consequences:** Sampling library (Poisson, Binomial, Multinomial) tests golden replay zorunlu.
**Tarih:** 2026-09-08

---

## DD-104 — Combat Session Ephemeral (Redis), Manuel + AFK Aynı Motor
**Karar:** Live combat session state Redis'te (ephemeral); tek node sticky. AFK combat contract analitik hesap (per-tick simülasyon YOK — DD-084 uyumlu). Aynı formula (damage/heal/scaling) hem manuel hem AFK.
**Reason:** Consistency: aynı build manuel + AFK aynı sonuç. Redis latency düşük (250ms tick). AFK aggregate O(1).
**Alternatives:** DB persist session (write amplification), farklı formula (drift + karışıklık).
**Consequences:** Node down = session migrate gerekli (Phase 02); Redis kayıp = session reset (rare). AFK DPS estimate accuracy Phase 02 refinement.
**Tarih:** 2026-09-09

---

## DD-105 — Threat Table Redis Sorted Set + Swap %10 Buffer
**Karar:** Mob threat table Redis sorted set (`mob:threat:<id>` ZSET). Target swap: yeni tehdit mevcut hedefin %10 üstünde ise. Taunt = threat + 1000 forced target duration.
**Reason:** Basit + performant + tank rotation baskısı (tanka her zaman en yüksek, sürekli micromanage gerekmez).
**Alternatives:** Threat decay (v1 kapalı; Phase 02 opt), taunt DR (Phase 02).
**Consequences:** Combat session cleanup Redis TTL 5 min post-combat.
**Tarih:** 2026-09-09

---

## DD-106 — Death Cascade: 5min Auto-Revive + Weakness 60s + XP Debt %5 (Level 20+)
**Karar:** Ölüm ceza kombine: 30sn self respawn / 5min auto-revive graveyard; 60sn Weakness debuff (-10% dmg + heal); XP debt %5 of level (level 20+); durability -5% each equipped; item loss YOK (DD-019).
**Reason:** Casual friendly (auto-revive), meaningful penalty (XP debt + weakness), item preservation (yıkıcı olmama). Ceza katmanlı.
**Alternatives:** Harsh (item loss + level down — reddedildi DD-010/019), soft (no penalty — trivializes).
**Consequences:** Weakness stack cap 3 (-30% max). Auto-revive 5min AFK-friendly. Weekly monitoring cascade rate.
**Tarih:** 2026-09-09

---

## DD-107 — CC Diminishing Returns 4 Level (100/50/25/0) + 20s Reset
**Karar:** CC (stun/root/silence) diminishing: aynı kategori ardışık apply → duration ×1.0 → ×0.5 → ×0.25 → immune. 20 sn no CC → reset. DR per category (stun/root/silence/disorient/disable).
**Reason:** CC chain-lock önlenir; tank/PvP balance; adil UX (immunity hissedilir).
**Alternatives:** Sabit duration (chain-lock riski), DR per-target-per-caster (karmaşık — WoW modern).
**Consequences:** Client UI DR indicator önerilir (Phase 02). CC balance Phase 02 playtest.
**Tarih:** 2026-09-09

---

## DD-108 — Kill Credit: Last Hit Attributor (v1); Party Assist Bonus
**Karar:** Mob kill credit **last-hit attributor** kazanır. Party member >10% total damage yaptıysa assist credit → XP share (v1 winner-takes-XP + party XP bonus +%10).
**Reason:** Simple + deterministic. Assist system Phase 02 (highest damage vs last hit debate).
**Alternatives:** Highest damage (attributor confusion), fully shared (loot debate).
**Consequences:** Assist credit v1'de basit party bonus; Phase 02 damage share formula.
**Tarih:** 2026-09-09

---

## DD-109 — Encounter Framework: Content-Driven Injection + Poisson Trigger
**Karar:** Encounter (rare mob, treasure, npc dialog, ambush) framework `content.encounter_template` + `character.encounter_history`. Reveal içinde Poisson sample per template (`expected = chance_per_hour × elapsed_hours`); character-per-encounter cooldown enforce.
**Reason:** Content-driven; her aktivite tipine injection; DD-084 analitik uyumlu (per-tick simulate yok).
**Alternatives:** Hard-coded encounter (esneklik yok), per-tick roll (yasak).
**Consequences:** Content team per region encounter pool tanımlar; cooldown per template character-özel (spam koruma).
**Tarih:** 2026-09-09

---

## DD-110 — Rare Mob v1'de Personal (Server-Wide Anons Phase 02)
**Karar:** Rare mob spawn v1'de **personal** — karakter kendi rare mob'unu görür + öldürür + rare drop alır. Server-wide anons Phase 02 opsiyonu (opt-in event framework).
**Reason:** Idle MMO fairness (herkes rare drop şansı eşit); AFK reveal deterministic; server-wide competition Phase 02 event olarak.
**Alternatives:** Global rare (first-come first-serve) — reddedildi v1: AFK oyuncular dezavantajlı.
**Consequences:** Server-first kill tracking (title reward) v1'de local first (character-first); global-first Phase 02.
**Tarih:** 2026-09-09

---

## DD-111 — Elite Tier 4 Level (Normal/Elite/Champion/Mini-Boss)
**Karar:** Mob elite_tier 4 seviye: 0=normal, 1=elite (HP ×3), 2=champion (×5 + rare guaranteed + weekly lockout), 3=mini-boss (×10 + region unique + phase system).
**Reason:** Content zenginlik + party incentive; mini-boss dungeon-outside boss encounter.
**Alternatives:** Sadece elite (variety kaybı), sadece boss (grind wall).
**Consequences:** Content tier balance; champion/mini-boss weekly lockout; content ihtiyaç per region.
**Tarih:** 2026-09-09

---

## DD-112 — Adventure Stage Deterministic Check (Combat/Skill/Resource/Luck)
**Karar:** Adventure stage check 4 tip: combat (DPS threshold), skill (stat min), resource (item consume), luck (RNG seeded). Snapshot deterministic; reveal analytical.
**Reason:** DD-088 (snapshot lifetime) + narrative variety; per-stage simulate yok.
**Alternatives:** Full RNG per stage (unfair — luck heavy), sadece combat (variety kaybı).
**Consequences:** Content team her stage check_type tanımlar; narration_key i18n; failure soft (adventure continues).
**Tarih:** 2026-09-09

---

## DD-113 — Dungeon AFK-First: Analytical Per-Floor/Per-Room Bounded Loop
**Karar:** Dungeon idle mode AFK contract; reveal analitik hesap (DD-084 uyumlu). Per-room ve per-boss bounded loop (dungeon başına ~50 room + ~3 boss = küçük N iteration OK). Per-tick simulate YOK.
**Reason:** Ölçek — 100k karakter × 8h AFK dungeon reveal O(1) hedefi. Bounded loop OK; küçük N sabittir (dungeon size template-driven).
**Alternatives:** Per-mob simulate (yasak DD-084), full deterministic aggregate (variance kaybı).
**Consequences:** Dungeon size cap (max floor 10, max room 100 per dungeon v1); mega dungeon Phase 02+ ayrı mimari değerlendirmesi.
**Tarih:** 2026-09-09

---

## DD-114 — Dungeon Idle Yield Multiplier 0.5 (Content-Override)
**Karar:** Dungeon idle run yield rate ×0.5 default (`afk_yield_multiplier`); content template override edebilir. Loot quality (rare drop rate) aynı per clear — sadece rate yavaşlar.
**Reason:** Manuel oynama incentive (2× hızlı); idle path fair (kalite düşmez, süre uzar).
**Alternatives:** ×1.0 (manuel value kaybı), ×0.3 (idle disincentive).
**Consequences:** Content team her dungeon için idle balance test; monetization Phase 02 "auto-clear scroll" opsiyonu.
**Tarih:** 2026-09-09

---

## DD-115 — Personal Chest Instance (v1 Dungeon)
**Karar:** Dungeon treasure/chest per character personal instance (v1). Party members her biri kendi chest'ini açar; no competition.
**Reason:** DD-012 personal loot ile tutarlı; party fairness; idle-run auto-open safe.
**Alternatives:** Shared chest (competition), round-robin (Phase 02 opt).
**Consequences:** Chest content persist-free (procedural per-character on interact); AFK reveal auto-open per member.
**Tarih:** 2026-09-09

---

## DD-116 — Dungeon Weekly Lockout Manuel + Idle Paylaşımlı
**Karar:** Character × dungeon × difficulty weekly lockout counter manuel run + idle run **paylaşımlı**. Örnek: manuel'de 3 clear → idle'da rare loot penalty aktif. Weekly reset Salı 06:00 UTC (DD-016).
**Reason:** Idle exploit önleme (manuel farm + idle farm double-dip yasak). Fair ekonomi.
**Alternatives:** Ayrı counter (double reward), sadece manuel counter (idle free-play).
**Consequences:** Character-özel weekly counter Postgres; idle activity reveal öncesi lockout check.
**Tarih:** 2026-09-09

---

## DD-117 — Trap: Skill/Luck Check + Damage; Advance Guaranteed
**Karar:** Dungeon trap check DEX/INT/WIS/luck bazlı; fail → damage veya CC; ancak room advance her zaman (trap wipe-prevention). Rogue class bypass flag (Phase 02 disarm skill).
**Reason:** Trap variety + strategy; wipe risk trap'dan gelmesin (frustration).
**Alternatives:** Wipe trap (harsh, unfair), no damage (meaningless).
**Consequences:** Content team trap difficulty tune; Phase 02 rogue trap detection passive.
**Tarih:** 2026-09-09

---

## DD-118 — Pet Stamina v1 AKTİF (Character Stamina Kapalı Rağmen)
**Karar:** Pet stamina v1'de aktif (0..100 pool, drain/rest cycle). Character stamina kapalı (DD-054); ancak pet için mekanik meaningful — pet care + activation loop.
**Reason:** Pet frequently activated/deactivated; stamina care burden idle MMO'da hoş bir "care mekaniği"; character wall değil (pet exhaust yalnız pet contribution durdurur, character oyna devam).
**Alternatives:** Pet stamina kapalı (character ile paralel; care depth kaybı), stamina hard block (frustration).
**Consequences:** Rest endpoint city-gated (v1); pet exhaust threshold AFK reveal içinde detection.
**Tarih:** 2026-09-09

---

## DD-119 — Pet Happiness Multiplier: XP Share 0.5-1.0 (Care Ödül)
**Karar:** Pet XP share owner activity XP × 20% × happiness_multiplier (0.5 + happiness/100 × 0.5, range 0.5..1.0). Happy pet full share; unhappy half.
**Reason:** Care mekaniği ödüllendirilir; unattended pet ceza (yumuşak, ceza değil "eksik ödül").
**Alternatives:** Flat XP share (care meaningless), harsh multiplier (frustration).
**Consequences:** Client tooltip happiness display + feed prompt; feed cooldown 6h over-feed protection.
**Tarih:** 2026-09-09

---

## DD-120 — Pet Evolution: Level Reset + Trait Preserve
**Karar:** Pet evolution (species change) sırasında pet **level 1'e reset**; XP=0. **Traits preserved** (existing rolled traits korunur). **Skills reset** (yeni species base skills). Bond preserved (Phase 02).
**Reason:** New species higher ceiling ancak grinding fresh; traits kimlik korunur (attachment protection).
**Alternatives:** Preserve all (evolution meaningless), reset all (attachment kaybı).
**Consequences:** Content team evolution chain balance (post-evolve stats fazla; level grind tekrarı worthwhile).
**Tarih:** 2026-09-09

---

## DD-121 — Pet Ölmez (Downed 30s Auto-Revive)
**Karar:** Pet HP 0 → downed state (Redis TTL 30s); auto-revive full HP. Pet permanent damage yok. DD-024 ile tutarlı (bound + BoP + not destroyable).
**Reason:** Pet attachment protection; casual friendly; combat frustration engellenir.
**Alternatives:** Pet ölür (hardcore mode Phase 03 opt), permanent stat penalty (frustration).
**Consequences:** Pet respawn kolay; combat pet cap %15 DPS contribution — balance risk azalır.
**Tarih:** 2026-09-09

---

## DD-122 — Mount Level Cap 30, XP Travel-Based (1 XP/unit)
**Karar:** Mount level cap 30 (v1); XP kaynağı travel distance (1 XP per distance unit). Uzun travel = mount level up otomatik. Legendary mount cap 40; mythic 50.
**Reason:** Mount identity = travel companion; travel XP source doğal + AFK-friendly. Cap 30 v1 balanced (level 30 ~15k XP → ~15k unit travel ≈ 30-50 hour AFK travel).
**Alternatives:** XP from combat (mount combat gerektirir Phase 03), fixed level (progression kaybı), high cap 100 (grind wall).
**Consequences:** Content team level 30 cap için stat scale (+50% speed at cap); AFK travel reveal içinde mount XP auto grant.
**Tarih:** 2026-09-09

---

## DD-123 — Mount Terrain Compatibility per Species (Content JSON)
**Karar:** Her mount species `terrain_compatibility_json` — terrain başına multiplier (ör. `mountain: 1.6` penalty, `water: 5.0` slow). Aquatic mount water bonus (0.5); non-aquatic water severe penalty.
**Reason:** Species identity meaningful; oyuncu terrain-based mount seçimi yapar; content depth.
**Alternatives:** Uniform terrain (species diversity kaybı), harsh forbid (route reroute).
**Consequences:** Route pathfind mount terrain penalty aware; content team her species için terrain profile.
**Tarih:** 2026-09-09

---

## DD-124 — Mount 1 Aktif Slot v1 (Multi-Mount Phase 03)
**Karar:** Karakter aynı anda **1 aktif mount** (v1). Multi-mount (2+ concurrent) Phase 03 opsiyon.
**Reason:** Simplicity; travel single-track logic; combat dismount rules unambiguous.
**Alternatives:** Cargo + speed mount combo (Phase 03), infinite (chaotic).
**Consequences:** Character mount seçimi (cargo vs speed) trade-off. Multi-mount ihtiyacı Phase 03 balance.
**Tarih:** 2026-09-09

---

## DD-125 — Quest Kategoriler: 7 Tür (Main/Side/Daily/Weekly/Dynamic/Event/Guild)
**Karar:** Quest sistemi 7 kategori: main (sıralı story), side (opsiyonel bölgesel), daily (3 slot 24h reset), weekly (5-10 slot Salı 06:00 UTC), dynamic (procedural bounty), event (event-tied), guild (Phase 02). Aynı `content.quest_template.category` field.
**Reason:** Variety + retention loop; casual (daily/weekly) + hardcore (main chain) + emergent (dynamic/event) + social (guild).
**Alternatives:** Sadece main+side (WoW classic), sadece daily/weekly (mobile MMO), farmable-only (grind fatigue).
**Consequences:** Content team her kategori için tune ve reward balance; 7 kategori endpoint filter/UI ayrımı.
**Tarih:** 2026-09-09

---

## DD-126 — Daily/Weekly Missed = No Makeup (Cadence Discipline)
**Karar:** Daily/weekly quest 24h/7 gün içinde tamamlanmazsa expire; no makeup. Next reset yeni pool.
**Reason:** Cadence discipline (idle MMO'da "kaçırdım catch-up" grind spirali önlenir); casual friendly (baskı yok).
**Alternatives:** Grace period + makeup, accumulate reward — reddedildi: FOMO artırır.
**Consequences:** UI clear: "expires at X"; missed daily normal; season achievement missing 1 daily kabul edilebilir.
**Tarih:** 2026-09-09

---

## DD-127 — Objective AFK-Tick: Aktivite Reveal İçinde Auto-Progress
**Karar:** Kill/gather/discover objective aktivite reveal içinde otomatik tick (kill counter, gather counter, POI discovery). Deliver/dialog objective manuel player action.
**Reason:** AFK-first design (DD-003); passive quest progression casual friendly; explicit action objective narrative olmadan devam etmesin.
**Alternatives:** Manuel-only objective tracking (idle MMO'da anti-pattern), tam auto everything (dialog narrative kaybı).
**Consequences:** Reveal payload objective progress delta; client tooltip AFK progress hint.
**Tarih:** 2026-09-09

---

## DD-128 — Party Size Cap 5 v1 (Raid Phase 02+)
**Karar:** Party max 5 üye v1. Raid (10+) Phase 02+ ayrı sistem.
**Reason:** MMO klasik small party; dungeon balance 5-man; UI/roster manageability.
**Alternatives:** 10 (raid semi-support), 3 (küçük grup), 8 (arbitrary).
**Consequences:** Content dungeon party_size_max=5 default; raid content Phase 02 mimarisi yeniden değerlendirilir.
**Tarih:** 2026-09-09

---

## DD-129 — Party Contract v1: Per-Character Independent
**Karar:** Party ile hunt/dungeon-idle → **her karakter kendi contract**'ını çalıştırır. Party bonus modifier snapshot'a eklenir. Shared party contract (leader-anchored) Phase 02 opsiyonu.
**Reason:** Simplicity; per-character reveal independent (member offline farklı zaman); loot personal fair.
**Alternatives:** Shared contract (member fair sync ancak reveal timing complexity), leader-only credit (member fair değil).
**Consequences:** Party disband mid-contract → party bonus modifier next snapshot refresh kaybolur (v1 contract lifetime lock); Phase 02 dynamic.
**Tarih:** 2026-09-09

---

## DD-130 — Auto-Leader Migration 5min Timeout
**Karar:** Leader disconnected 5 dakika → next longest-active member auto-promote. Vote kick threshold 2/3.
**Reason:** Party continuity; leader AFK/disconnect party paralizasyonu önlenir.
**Alternatives:** Manual only (party stuck), instant (unfair for brief disconnects).
**Consequences:** Notification "New leader"; leader reconnect membership korunur ama leadership değişmiş.
**Tarih:** 2026-09-09

---

## DD-131 — Guild Cap 50 v1 (Level ile 100'e Kadar)
**Karar:** Guild başlangıç member cap 50 (DD-018 reaffirm); guild level ile 100'e kadar expand (level 5: +5, level 10: +10, level 25: +25, level 50: +50 = 100 cap). Level cap v1: 10 (Phase 02: 20).
**Reason:** Küçük/orta guild kültürü; guild progression member cap ile ödüllendirilir; mega-guild dominance önlenir.
**Alternatives:** Sabit 100 (küçük guild disincentive), 20-30 (yetersiz), 500 (chaos).
**Consequences:** Content team level-cap-expansion balance; player pipe motivation for guild leveling.
**Tarih:** 2026-09-09

---

## DD-132 — Guild 5 Sabit Rank v1 (Custom Rank Phase 02)
**Karar:** Guild rank hierarchy v1 sabit 5 rank (leader/officer/veteran/member/initiate) + sabit default permissions. Custom rank (rename + custom permissions) Phase 02.
**Reason:** Simplicity + predictable UX; content team permission testing zor değil; complex customization gerekmez v1.
**Alternatives:** Full custom v1 (over-engineered), 3 rank (yetersiz).
**Consequences:** Guild yönetim UI 5 rank fixed; Phase 02 admin panel custom.
**Tarih:** 2026-09-09

---

## DD-133 — Guild Bank BoP Block v1 (Owner-Lock Phase 02)
**Karar:** Guild bank BoP (bind on pickup) item deposit **yasak** v1. Phase 02: BoP allowed but only original owner can withdraw.
**Reason:** BoP simplicity v1 (item ownership tracking complexity); attachment loss önlenir.
**Alternatives:** BoP allowed with dedup (complex), full block v1 (missing feature).
**Consequences:** Content team crafting output BoE tercih; BoP quest item guild bank ret; Phase 02 owner-lock flag schema.
**Tarih:** 2026-09-09

---

## DD-134 — Guild Disband 48h Grace + Cancel Window
**Karar:** Leader disband trigger 48 saat grace period; leader cancel edebilir. Grace expire: bank items → member mail (Phase 02) veya leader claim (v1), treasury leader claim.
**Reason:** Accidental disband koruma; members warning window; treasury/bank fair.
**Alternatives:** Instant disband (harsh), 7 gün (uzun bekleyiş).
**Consequences:** Worker `guild_disband_expire` daily; UI clear countdown; audit trail.
**Tarih:** 2026-09-09

---

## DD-135 — Ekonomi Ledger-First (Double-Entry + SERIALIZABLE)
**Karar:** Her değer transferi tek servis (`TxnService.Write`) üzerinden `economy.ledger_entry`'ye double-entry (debit + credit, sum=0) yazılır; SERIALIZABLE isolation'da wallet/balance projeksiyonu ve `audit.ledger_audit` aynı TX'te güncellenir. Domain modülleri wallet/balance sütununa doğrudan UPDATE yazamaz.
**Reason:** Ekonomik state için tek hakikat kaynağı; drift/duplicate/loss forensic olarak tespit edilebilir; reconciliation mümkün; audit compliance.
**Alternatives:** Direct wallet UPDATE + async log (drift kaçınılmaz), event-sourced full CQRS (aşırı karmaşık v1).
**Consequences:** Tüm domainler `TxnService` bağımlılık; retry budget (serialization_failure exp backoff x3); nightly reconciler drift alarm; ledger 12 ay hot 12+ ay cold.
**Tarih:** 2026-09-09

---

## DD-136 — Gold Base Unit = Copper (Integer BIGINT)
**Karar:** Server ve DB'de tüm gold değerleri `copper` cinsinden integer (BIGINT). `1 gold = 100 silver = 10.000 copper`. UI dönüşümü client tarafında. Softcap 10¹² copper (10 milyar gold) cüzdan üstlimiti.
**Reason:** Float rounding ledger'da tolere edilemez; integer copper ile deterministic sum; UI trivyally 3 sütun.
**Alternatives:** Decimal(20,4) (yavaş + hâlâ zorlayıcı), float64 (rounding cehennemi), silver base (çok küçük precision).
**Consequences:** Fiyat formülleri integer copper üzerinden; `round HALF_UP`; cap aşımı `E_GOLD_CAP_EXCEEDED` (overflow → Phase 02 mail).
**Tarih:** 2026-09-09

---

## DD-137 — v1'de Sadece Gold Trade Edilebilir
**Karar:** Currency taksonomi: `gold` (tradable), `gem` (premium/BoA), `honor` (Phase 02 PvP/BoA), `faction_<code>` (BoA), `event_<code>` (BoA), `guild_contribution` (guild-scoped). V1'de yalnızca `gold` P2P trade edilebilir; diğerleri Bind-on-Account veya karakter-özel.
**Reason:** Premium/reward currency piyasada gold ile hedge'lenmesin; RMT ve monetization erozyonu önlensin; ekonomik denge basit.
**Alternatives:** Tüm currency tradable (RMT riski), currency exchange (arbitraj + sink/source dengesiz).
**Consequences:** Trade offer'da yalnızca `gold_copper` alanı; auction sadece `gold` fiyat (Faz 26); currency conversion yasağı Phase 02 review'e kadar sürer.
**Tarih:** 2026-09-09

---

## DD-138 — Trade Window: Değişiklik Sonrası Auto-Unlock (Anti Bait-and-Switch)
**Karar:** Trade session'da bir taraf offer'ı güncelledikten sonra her iki tarafın lock'u ve confirm'u otomatik sıfırlanır. Confirm request'i `seen_partner_offer_hash` içerir; server state hash'i ile mismatch → `E_TRADE_OFFER_CHANGED`.
**Reason:** Klasik "son saniye item swap" scam'ini engelle; oyuncu her zaman gördüğü offer için taahhüt eder.
**Alternatives:** Manuel her değişiklikte tekrar lock talep (kullanıcı deneyimi kötü), guard yok (scam garantili).
**Consequences:** Trade UI değişiklikte "confirm reset" görsel bildirim; cancel sonrası 30 sn cooldown (spam koruma); rate limit karakter başına 3 session/dk.
**Tarih:** 2026-09-09

---

## DD-139 — Vendor Sell Multiplier 0.25 (İmplicit Sink)
**Karar:** Vendor sell_multiplier varsayılan `0.25` (satılan itemin base fiyatının %25'i oyuncuya). Kalan %75 `system.vendor_sink`'e implicit ledger entry; ekstra tax entry yazılmaz. Rarity multiplier: common 1.0 → legendary 2.0.
**Reason:** V1'de ana gold sink; vendor floor stabil; auction house olmadan (Faz 26'ya kadar) net sink garanti.
**Alternatives:** 0.50 (yetersiz sink), 0.10 (oyuncu satmaz — yığılma), dinamik pricing (Phase 03).
**Consequences:** Rarity multiplier config'de; oyuncu tercihi: legendary'yi disenchant vs vendor karar; buyback 2h TTL (accidental sell koruma).
**Tarih:** 2026-09-09

---

## DD-140 — Buyback Penceresi 2 Saat + Yeni UUID
**Karar:** Player→vendor satılan item 2 saat boyunca `economy.vendor_buyback` tablosunda saklanır, aynı fiyatla geri alınabilir. Geri alınan instance yeni UUID alır; orijinal UUID resurrect edilmez; history ledger'da linkli.
**Reason:** Accidental sell koruma; ancak instance identity restore edilirse dupe/audit riski (attached enchant, socket state).
**Alternatives:** Sınırsız buyback (dupe riski), buyback yok (frustration), aynı UUID (dedup problemleri).
**Consequences:** Instance snapshot JSONB tutulur; expired buyback worker temizler; UI vendor ekranında ayrı "buyback" tab.
**Tarih:** 2026-09-09

---

## DD-141 — Trade Tax V1'de Kapalı (Phase 02 Opt %2 >10k g)
**Karar:** V1'de P2P trade tax `%0` (yeni oyuncu / arkadaş takası dostluğu). Phase 02'de opsiyonel: gold_leg > 10.000 g için `%2` tax → `system.tax_sink`. Config `trade_tax.enabled=false` v1.
**Reason:** V1'de sink yeterince agresif değil (vendor spread + fee'ler yeterli); trade tax whale swap için ileride opt.
**Alternatives:** Sabit %5 tax v1 (sosyal takas caydırıcı), progressive tax (v1 için karmaşık).
**Consequences:** Config hot reload; Phase 02 aç/kapa; explicit karar Phase 02 balance pass sonucu.
**Tarih:** 2026-09-09

---

## DD-142 — Reversal Sadece Ops (Kullanıcı Reversal Yok)
**Karar:** Ledger append-only; iptal = `kind='reversal'` yeni entry group. Kullanıcı UI'sinden reversal yolu **yok**; yalnızca admin `/internal/ops/economy/reversal` endpoint (admin JWT + 2FA + audit).
**Reason:** Kullanıcı-triggered reversal dupe/rollback abuse'a açık; support ihtiyacında op onaylı reversal yeterli.
**Alternatives:** Kullanıcı iade butonu (abuse), reversal yok (support için imkansız).
**Consequences:** Support runbook Phase 02; her reversal `audit.ledger_audit`'e "operator_id + reason" yazar; SLA Phase 02.
**Tarih:** 2026-09-09

---

## DD-143 — Marketplace v1: Sadece Fixed-Price + Global Single Market
**Karar:** Faz 26 v1'de marketplace **fixed-price listing only** (bid-based auction Phase 02'ye ertelenir). Dünya çapında **tek global market** (shard/region bölünmesi yok). Sadece `gold` currency (DD-137 uyumlu).
**Reason:** Auction karmaşıklığı (bid escrow, anti-snipe, outbid iade) v1 scope'unu şişirir; fixed-price piyasa keşfi için yeterli veri üretir; global market likidite garantisi.
**Alternatives:** Auction v1'de (Phase 02'ye ertele), regional market (likidite parçalanır).
**Consequences:** Auction dokümanı hazır (`marketplace/auction.md`) ancak endpoint disabled v1; regional market Phase 03 tartışması.
**Tarih:** 2026-09-09

---

## DD-144 — Marketplace Escrow: Item Lock + Ownership Sabit (v1)
**Karar:** Listing süresince item satıcı üzerinde kalır ancak `is_locked=true, lock_owner='market', locked_ref_id=listing_id` ile lock'lanır; `market.escrow_item` join tablosu tutulur. Ownership değişmez (envanterde "listed" badge grey görünür ama filter dışı).
**Reason:** Ownership değişikliği (system UUID'ye transfer) audit karmaşıklığı yaratır; lock-based basit + reversible; expire/cancel'de sadece unlock yeterli.
**Alternatives:** Ownership transfer to `MARKET_ESCROW_UUID` (system item bloat), envanterden fiziksel çıkarma (envanter reconciliation zor).
**Consequences:** Envanter UI "listed" gösterimi; overflow (envanter dolu expire'da) için `market.overflow_item` 7 gün claim window.
**Tarih:** 2026-09-09

---

## DD-145 — Listing Fee %5 Base × Duration Multiplier + Sale Cut %5
**Karar:** Listing fee = `unit_price × qty × 0.05 × duration_mult` (12h=1.0, 24h=1.5, 48h=2.5, 168h=5.0); min 10c. Sale cut = `total_price × 0.05` alıcının payment'ından. İki fee de `system.tax_sink`.
**Reason:** Duration multiplier uzun listeleme dengeleyicisi; sale cut kesin sink; %5 rate WoW/EVE aralığında; Phase 02 balance pass'te tune edilir.
**Alternatives:** Flat fee (küçük item spam), sadece sale cut (unsold spam), %10 (yeni oyuncu prohibitive).
**Consequences:** Config hot reload; `market_listing_fee_collected_copper_total` + `market_sale_cut_collected_copper_total` metrikleri; expected sink volume Phase 02 modellenir.
**Tarih:** 2026-09-09

---

## DD-146 — Marketplace Same-Account Block (Post Değil, Buy Zamanında)
**Karar:** Marketplace listing post sırasında **kontrol yok**. Buy zamanında `listing.seller_account_id == buyer.account_id` → `E_SAME_ACCOUNT_BLOCK`. Aynı karakter kendi listing'ini alamaz (`E_SELF_BUY_BLOCK`).
**Reason:** Post zamanında bilinemez (bir daha post edilebilir); buy zamanında zorunlu; self-fund (multi-char aynı account arası fee kaçırma) engeli.
**Alternatives:** Aynı account trade tamamen serbest (self-fund risk), post zamanında block (yararsız — buyer bilinmez).
**Consequences:** `seller_account_id` listing satırında denorm tutulur (join cost yok); trade P2P'de v1'de aynı block yok (DD ayrı — Phase 02 review).
**Tarih:** 2026-09-09

---

## DD-147 — Price Clamp: Min = Vendor Sell × 1.10, Max = Vendor Buy × 10.000
**Karar:** Listing unit_price sınırları: `min = max(1c, base_price × vendor.sell_multiplier × 1.10)`, `max = base_price × vendor.buy_multiplier × 10.000`. Dışında `E_PRICE_OUT_OF_RANGE`.
**Reason:** Min → vendor arbitraj engeli (marketten al, vendor'a sat kâr edemesin); max → fat-finger + inflate koruma; free market ama uçlar kilitli.
**Alternatives:** Sadece min (fat-finger boşta), sadece max (arbitraj boşta), sabit clamp (item bazında dengesiz).
**Consequences:** Item base_price content'inde doğru olması kritik; legendary premium clamp Phase 02.
**Tarih:** 2026-09-09

---

## DD-148 — Cancel Fee İade Kademeli (1h/6h/expire)
**Karar:** Cancel fee iadesi: <1h %50, 1-6h %25, ≥6h %0. Expire'de iade **yok**.
**Reason:** Fat-finger post iyi niyet iadesi; ancak "post-cancel-post" arbitraj (fiyat probe) engellenmeli; expire'de tamamen yanmış → aktif listing pool temiz.
**Alternatives:** Tam iade her zaman (spam), hiç iade yok (frustration).
**Consequences:** Iade `system.tax_sink`'ten negatif entry (net sink metriği doğru); overflow item envanter doluysa 7 gün claim.
**Tarih:** 2026-09-09

---

## DD-149 — Search Server-Side + Redis Cache 30sn + Rate Limit 20/dk
**Karar:** Marketplace search **server-only** (client 10k+ listing çekmez). Postgres primary index + Redis result cache (normalized-param hash key, TTL 30sn). Karakter başına rate limit 20 req/dk.
**Reason:** Client-side impossible (data hacmi); 30sn cache staleness kabul edilebilir (sold item hızla düşer); rate limit spam bot koruma.
**Alternatives:** Realtime (cache yok — DB CPU patlar), long cache 5dk (stale sold data), no rate limit (bot flood).
**Consequences:** Phase 02'de OpenSearch/Elastic sekonder index (facet + fuzzy); JSONB GIN index snapshot filter için.
**Tarih:** 2026-09-09

---

## DD-150 — World Events: Shared Timer + Personal Contribution
**Karar:** Tüm world event türleri (world boss, invasion, regional, seasonal, dynamic) **shared timer** (aynı UTC starts_at/ends_at herkeste görünür) + **personal contribution ledger** (`progression.*_attendance` / `*_participation`) modelini kullanır. Rekabet (top damager kill vurucu, node yarışı) yok; reward tier contribution eşiklerine göre bağımsız.
**Reason:** AFK MMO'da online rekabet frustration; DD-091/115 (personal instance) uyumlu; leaderboard olsa bile ödül bağımsız = herkes tam paylaşımdan yararlanır.
**Alternatives:** Top-N reward (rekabet), first-tag steal (klasik MMO), shared HP + top damager loot (AFK dostu değil).
**Consequences:** Contribution normalizasyon (damage + heal×0.5 + AFK score) formülü; cosmetic leaderboard opsiyonel (title/tabard) ödüle etki etmez; personal instance boss visualization Phase 02 (client illusion "birlikte savaştık").
**Tarih:** 2026-09-09

---

## DD-151 — World Boss AFK Attendance Contract (Registration 24h Window)
**Karar:** World boss'a `attend` registration `starts_at - 24h → starts_at` window'unda açık; kayıt olan karakter offline/AFK olsa bile `activity-engine` "world_boss_attend" contract çalıştırır. AFK contribution `min(afk_score, tier_3_threshold)` cap — sadece online katılım tier_4 açar.
**Reason:** AFK-first oyuncu event kaçırmasın; ancak online katılım avantajı korunur (tier_4 exclusive) — engagement teşviki.
**Alternatives:** AFK katılım yok (frustration), AFK full tier (online katılım anlamsız), auto-attend (max concurrent AFK çakışması).
**Consequences:** Contract snapshot (DD-088) at registration; max concurrent AFK slot (DD-OPEN-3) çakışma kararı; tier eşikleri content'te.
**Tarih:** 2026-09-09

---

## DD-152 — Invasion Region-Active Auto-Eligibility (24h Recent Regions)
**Karar:** Invasion event trigger anında region'da bulunan veya son 24 saat içinde region'da aktif olan karakter otomatik participation kaydı alır. Explicit "join" yok. AFK contract event window'da otomatik oluşturulur.
**Reason:** Region loyalty ödüllendirilsin; oyuncu her region'a bir kez uğrayıp "farm eligibility" yapamasın; sürtünmesiz katılım.
**Alternatives:** Explicit join (sürtünme), 7d horizon (bölge-hopping abuse), current region only (offline karakter frustration).
**Consequences:** `character.recent_regions` history 24-72h tutulur (bkz. character schema Phase 01); dynamic trigger Phase 02.
**Tarih:** 2026-09-09

---

## DD-153 — Seasonal Event Token: Account-Scoped BoA + 14 Gün Post-Event TTL
**Karar:** Seasonal event `event_<code>` token **account scope** (tüm karakterler paylaşır), **Bind-on-Account** (trade edilemez, DD-137 uyumlu). Event ends_at + 14 gün vendor açık ve token geçerli; sonra token otomatik silinir.
**Reason:** Alt karakter oyuncusu tek token pool'la etkileşim (adalet); post-event 14 gün grace player'ın vendor'da harcama planlamasına izin verir; sürekli birikme (yıllar arası) inflation yaratır.
**Alternatives:** Character-scoped (alt karakter mağdur), permanent token (inflation), event bittiği an temizleme (frustration).
**Consequences:** `economy.account_currency_balance` schema yeterli (DD-137); worker `seasonal_event_expire` T+14d cleanup; content bundle rotation Phase 02.
**Tarih:** 2026-09-09

---

## DD-154 — Seasonal Event Auto-AFK Login Credit %50 (Max 7 Gün)
**Karar:** Seasonal event window'da account hiç login yapmazsa günlük login token'ın %50'si otomatik AFK credit ile `account_currency_balance`'a yazılır (max 7 gün retroaktif). 7 günü aşan kayıp günler kaybolur.
**Reason:** AFK-first oyuncu event kaçırmasın (tatil, uzun süreli offline); yine de login teşviki (%100 login-only bonus); 7 gün cap kötü niyet önleme.
**Alternatives:** Login zorunlu (AFK-first ihlali), %100 auto (login incentive sıfırlanır), sınırsız retroaktif (abuse).
**Consequences:** Worker `seasonal_login_auto_credit` daily; config `auto_afk_credit_ratio` + `auto_afk_credit_max_days` hot reload; UI'da "auto credit earned" görünür.
**Tarih:** 2026-09-09

---

## DD-155 — Dynamic Event Chain Max Derinlik 3 (v1)
**Karar:** Dynamic event chain (parent success → child spawn) v1'de **max 3 seviye derinlik**. 3. seviye başarısı `chain_bonus_final` (title + unique item) verir.
**Reason:** Uzun chain oyuncu takibi zor; content authoring maliyeti; balance testable; anlamlı narrative yay.
**Alternatives:** Sınırsız chain (complexity + confusion), 1 seviye (chain hissiyatı yok), 5 seviye (Phase 02 opt).
**Consequences:** `chain_success_unlock` JSONB delay-based scheduler; `dynamic_event_chain_length_histogram` metric; Phase 02'de 5 seviye + player-triggered.
**Tarih:** 2026-09-09

---

## DD-156 — Block Account-Level (Karakter-Level Değil)
**Karar:** Block ilişkisi **account ↔ account** düzeyinde tutulur (`social.block_list` account UUID'leri ile). Blocker'ın block ettiği account'un **tüm karakterleri** (mevcut + gelecek) etkilenir.
**Reason:** Toxic oyuncu alt karakter yaratıp aynı taciz döngüsünü sürdürmesin; block bir kez kesin kesin.
**Alternatives:** Karakter-karakter (bypass kolay), IP-based (dinamik IP, aile paylaşımı false positive), device-based (v1 kapsam dışı).
**Consequences:** Karakter→account lookup her check'te; Redis cache `block:blocker:{account}` set (15 dk TTL); block/unblock cache invalidate; friend + block çelişki auto-remove friend.
**Tarih:** 2026-09-09

---

## DD-157 — Global/Trade/LFG/Region Chat Persist Edilmez (Redis Rolling Only)
**Karar:** Public global kanallar (`global`/`trade`/`lfg`/`region`) `chat.message` tablosuna yazılmaz. Yalnızca Redis rolling window (MAXLEN 200, TTL 24h) + WS live delivery. Compliance için opt-in %1 örnek sampling Phase 02 (`audit.chat_sample`).
**Reason:** Public chat mesaj hacmi (10+ req/sn per region) DB'yi şişirir; kişisel scroll'da global mesaj yok zaten; forensic sampling yeterli.
**Alternatives:** Full persist (DB pressure + retention maliyeti), zero sample (moderation forensic zayıf), 100% örnek (hâlâ hacim).
**Consequences:** Yeni join'e history yalnızca Redis'ten (last 200); global mesaj audit için sampling infra Phase 02; slow mode + trust score + auto-mute agresif.
**Tarih:** 2026-09-09

---

## DD-158 — Whisper/Guild/Party/System Chat: chat.message + 30d TTL + Personal Scroll
**Karar:** Whisper, guild, guild_officer, party, system kanalları `chat.message` tablosuna persist edilir; retention 30 gün; kişisel scroll (client) + server history endpoint. Whisper için `read_at` unread badge sayaç.
**Reason:** Kişisel/dar-kapsam kanallarda hacim yönetilebilir; kişisel scroll oyun deneyiminin parçası (mail-benzeri asenkron whisper).
**Alternatives:** Whisper ephemeral (kullanıcı beklentisine ters), guild ephemeral (yeni join scroll yok), sürekli persist (maliyet).
**Consequences:** Partial index (`WHERE channel_kind='whisper'`); nightly retention job 30+ gün silme; unread WS badge push.
**Tarih:** 2026-09-09

---

## DD-159 — Mail 30 Gün Expire + Attachment Auto-Return + Postage Yanmış
**Karar:** Mail 30 gün expire; opened=false ve attachment varsa `state='returned'` yeni iade mail sender'a (7 gün expires), aynı item/gold transfer. Postage `system.fee_sink`'te kalır (iade edilmez). Iade mail claim edilmezse `social.mail_overflow_item` 30 gün nihai window.
**Reason:** Uzun süreli offline oyuncu attachment kaybetmesin; postage'ın hep yakılması RMT bot maliyet artırır; overflow queue nihai kaybı yavaşlatır.
**Alternatives:** Zero return (frustration + item kaybı), sınırsız iade zinciri (dupe riski), postage iade (spam mail atarı bedava).
**Consequences:** Worker `mail_expire` her 5 dk `SKIP LOCKED`; iade mail postage'sız (sender ceza ödemez); overflow schema `mail_overflow_item` benzer market overflow.
**Tarih:** 2026-09-09

---

## DD-160 — Blocked Recipient'a Mail: Silent Drop, Attach Sender'da Kalır (v1)
**Karar:** Sender, kendisini block etmiş bir recipient'a mail atmaya çalıştığında: mesaj sender'a "sent" görünür (silent drop pattern), ancak attachment ve gold sender'da hiç düşmez (envanterinde kalır), postage da alınmaz. Server response 202 (fake success).
**Reason:** Silent drop toxic sender'a "block'landın" sinyali vermez (harassment escalation önlemi); attachment düşürmemek = spam bot para/item yakmasın; gerçek recipient bilmez.
**Alternatives:** Explicit `E_TARGET_BLOCKED` (harassment feedback), attach al postage ye (spam bot yem), attach al escrow'a (drift + admin cleanup).
**Consequences:** Debug için `moderation_flag='blocked_target'` mail'e yazılsa da attach reject inline; sender history'sinde "sent" görünür ama recipient inbox'ta yok; kendine block ettiği kişiye mail atarsa `E_TARGET_BLOCKED_SELF` explicit hata.
**Tarih:** 2026-09-09

---

## DD-161 — Chat WS Fanout Redis Pubsub + Region-Scoped Sub Keys
**Karar:** Chat delivery Redis pubsub üzerinden fanout. Sub key şeması: `chat:live:{kind}:{scope_key}` (örn `chat:live:global:region:<id>`, `chat:live:party:<id>`, `chat:live:guild:<id>`, `chat:live:whisper:<character_id>`). Client region değişince eski region key'lerinden unsub + yeni region sub.
**Reason:** Redis pubsub low-latency (< 50ms fanout); region-scoped sub connection başına küçük subscription set (bant tasarrufu); travel/region-change'te clean transition.
**Alternatives:** NATS (yeni infra), direct WS broadcast tüm subscribers (bandwidth patlar), global sub filter client-side (bant israfı).
**Consequences:** Redis pubsub reconnect handling; connection kaybında re-sub logic backend/websocket.md'de; region switch WS "region_changed" event → client sub update.
**Tarih:** 2026-09-09

---

## DD-162 — Long-Term Progression: 8 Sistem Server-Auth Event Stream + Karar-Yumuşak
**Karar:** Endgame progression 8 sistem (mastery, ascension, achievement, title, collection, milestone, reputation, faction) tek `event.progression.*` bus üzerinden yazılır → double-write yok. Her sistem "karar-yumuşak" (kaçırılan hedef alternatif yol sunar); AFK-first uyumlu (mastery/rep ana kazanç AFK aktiviteden).
**Reason:** 8 sistem birbirini tetikliyor (mastery tier up → title unlock → achievement); merkezi event bus double-count önler; karar-yumuşak yaklaşım FOMO azaltır (idle MMO'da haftalarca offline oyuncu için kritik).
**Alternatives:** Her sistem kendi kaynak polling'i (double-count riski), FOMO-driven (idle oyuncu kaybı), tek monolit progression tablosu (bakım imkansız).
**Consequences:** `event.progression.*` schema Phase 01 kritik; her sistem consumer worker; personal ledger append-only history.
**Tarih:** 2026-09-09

---

## DD-163 — Mastery Ascension'da Sıfırlanmaz (Kalıcı Horizontal)
**Karar:** Ascension level + XP + class skill puanı sıfırlar; **mastery XP ve tier'ları KALIR**. Mastery yatay progression olarak kalıcı; ascension sadece dikey (level) reset.
**Reason:** Yüzlerce saatte biriken mastery XP kaybı player rage; horizontal progression MMO'da "ascension'a değer" hissini artırır (ascension = cesaret, kayıp değil).
**Alternatives:** Full reset (rage + oyuncu kaybı), partial mastery reset (karmaşık iletişim), mastery ascension bonus (dengesizlik).
**Consequences:** Ascension UI clear ("kaybedecekleriniz: level, XP, class puanı; koruyacaklarınız: mastery, rep, kozmetik, gold, item, guild"); prestige tree ayrı puan sistem.
**Tarih:** 2026-09-09

---

## DD-164 — Ascension: Level 60 + 60h Playtime + Path Questi + 10k Gold + Orb
**Karar:** V1 ascension koşulları: `min_level=60`, `min_playtime_hours=60`, `path_of_ascension_v1` quest line complete, cost `10.000 g + 1 orb_of_ascension`. Frequency hedef 60-100h aktif oyun/ascension.
**Reason:** Playtime gate rush önler; questline story motivator; cost fee sink; orb item event/craft/reputation kaynağından — çoklu path.
**Alternatives:** Sadece level (rush + shallow), sadece playtime (AFK farm), cost yok (spam ascension).
**Consequences:** `path_of_ascension_v1` quest content zorunlu; orb kaynağı (event drop + craft + rep vendor) 3 path; frequency final tune DD-OPEN-5.
**Tarih:** 2026-09-09

---

## DD-165 — Achievement Puan Sistemi: 5/10/25/50/100 + Meta + Retroactive Opt-In
**Karar:** Achievement point value rarity göre: common 5, uncommon 10, rare 25, epic 50, legendary 100. Meta-achievement (alt sette tümü unlock) parent unlock. Yeni release'te retroactive default `false` (content bazında opt-in).
**Reason:** Basit skala oyuncu iletişimi; meta zincir uzun vadeli motivasyon; retroactive default false çünkü eski event log her achievement kind için garanti değil (ör. player state).
**Alternatives:** Puan yok (leaderboard input yok), sürekli retroactive (log arkeoloji karmaşık), tek katman achievement (meta yok).
**Consequences:** `character.achievement_points` denorm; retroactive Phase 02 ops runbook decision per content release; leaderboard Phase 02.
**Tarih:** 2026-09-09

---

## DD-166 — Collection Scope: Pet/Mount/Skin Account, Appearance Character (v1)
**Karar:** Collection kategorisi bazında scope: `pet/mount/mount_skin/tabard` **account** (paylaşımlı); `appearance/transmog` **character** (rol içi immersion) v1. Phase 02 review.
**Reason:** Pet/mount collection alt karakter tekrar farm istemez (adalet); appearance RP oyuncuları için karakter-özgü kimlik.
**Alternatives:** Hepsi account (RP kaybı), hepsi character (grind fatigue alt char).
**Consequences:** `progression.account_collection` + `progression.character_collection` iki tablo; set bonus scope-uyumlu; Phase 02'de appearance opt account.
**Tarih:** 2026-09-09

---

## DD-167 — Reputation 8-Tier + Daily Cap 1500 + Neutral Start (No Decay v1)
**Karar:** Reputation tier 0..7 (`hated -6000, hostile -3000, unfriendly -1000, neutral 0, friendly 3000, honored 9000, revered 21000, exalted 42000`). Daily rep gain cap **1500/faction/gün** (reset 00:00 UTC — DD-016). Yeni karakter tüm faction ile **neutral start**. V1'de **decay yok**.
**Reason:** Klasik MMO 8-tier familiarity; daily cap grinding koruma; neutral start yumuşak intro; decay v1 fatigue (aktif oyunca özel disincentive).
**Alternatives:** Slow decay v1 (idle MMO'da rage), günlük cap yok (24h farm meta), sabit start (rp/hikaye çakışması).
**Consequences:** `character_reputation.daily_gain_used` + `daily_reset_at` tracker; ekonomi sink (rep vendor); decay opt Phase 02 (90 gün inaktif → -1 tier gibi).
**Tarih:** 2026-09-09

---

## DD-168 — Faction Cross-Effect (Depth 1) + Rival Rep Loss
**Karar:** Faction relationship graph: ally (`+0.25 gain cross`), neutral (0), enemy (`-0.5 loss cross`). Rep gain propagation depth 1 (v1); Phase 02 depth 2. Rival faction hostile → NPC saldırgan + vendor/quest kilit.
**Reason:** Faction seçimi anlamlı olmalı — rep kazanmak bedeli olsun; propagation depth 1 hesap ucuz; hostility gate PvE derinlik.
**Alternatives:** Cross-effect yok (faction seçimi anlamsız), tüm faction depth-N traverse (perf + confusion), sadece binary hostility (WoW gibi).
**Consequences:** `content.faction_relationship` tablosu; rep_gain worker graph traverse depth 1; hostility gate combat threat kaydı + city ban logic.
**Tarih:** 2026-09-09

---

## DD-169 — Push Provider: FCM (Android) + APNs (iOS) v1 (WebPush Phase 02)
**Karar:** V1 push transport: **FCM HTTP v1** (Android) + **APNs HTTP/2 + JWT** (iOS). Provider-agnostic payload iç şema; platform adapter dispatcher'da. WebPush (browser) Phase 02'ye ertelendi.
**Reason:** İki büyük mobil platform Play/App Store zorunluluğu; Firebase + APNs endüstri standard; abstraction gelecekte OneSignal/Airship gibi aggregator veya WebPush eklemeye açık.
**Alternatives:** OneSignal/Airship (vendor lock + ekstra maliyet), sadece FCM (iOS Silent limit), sadece APNs (Android desteklenmez).
**Consequences:** `notification.device_push_token` schema + provider adapter; FCM analytics_label + APNs collapse-id/thread-id; provider credentials secret rotation Phase 02.
**Tarih:** 2026-09-09

---

## DD-170 — Notification Dedup Key + Idempotent Delivery
**Karar:** Her push mesajı `dedup_key` ile UNIQUE `(account_id, dedup_key)` constraint altında; retry aynı key düşürülür (`ON CONFLICT DO NOTHING`). Örnek: `activity_completed:<contract_id>`, `world_event_starting:<spawn_id>`.
**Reason:** Reveal retry + worker restart + provider retry idempotent olmalı; kullanıcı aynı bildirimi ikinci kez almasın.
**Alternatives:** At-most-once ack (loss riski), at-least-once no dedup (duplicate push → user rage), Redis distributed lock (karmaşık).
**Consequences:** `notification.push_message` unique index; storage minor cost 30 gün retention; retry policy exp backoff safe.
**Tarih:** 2026-09-09

---

## DD-171 — Notification 17 Kind + Granular Opt-Out + Uncontrollable Set
**Karar:** V1'de 17 notification kind (activity/dungeon/craft/travel/energy/world/guild/mail/whisper/friend/market×2/party_invite/guild_invite/achievement/system/security). Her kind kullanıcı opt-out (4 mode: push+inbox / inbox_only / silent_push / off). Uncontrollable set: `system_broadcast`, `account_security`, `payment_receipt` (Phase 02).
**Reason:** Granular kontrol retention için kritik; ancak security/maintenance opt-out riski (compliance + user harm) — override yasak.
**Alternatives:** Kaba category (retention loss), tümü opt-out (security açığı), tümü zorunlu (retention loss + kullanıcı rage).
**Consequences:** `notification.preference_kind` per-account tablo; UI disable toggle uncontrollable + tooltip; default values `content.notification_kind_defaults`.
**Tarih:** 2026-09-09

---

## DD-172 — Quiet Hours + Global Mute + Kritik Bypass
**Karar:** Kullanıcı quiet hours (start/end + timezone) tanımlar; bu pencerede push priority `silent` (badge/inbox only). `Global mute` X saat için tüm push kesilir (inbox delivery devam). Her ikisi de **critical bypass**: `system_broadcast` + `account_security` her koşulda ses açık.
**Reason:** Uyku/toplantı vs. spam koruma; tatil için tek toggle; kritik güvenlik/maintenance mesajı asla kaybolmasın.
**Alternatives:** Sadece kind opt-out (granuler ama tüm iş), no quiet hours (uyku spam), critical bypass yok (security açığı).
**Consequences:** Push dispatcher'da timezone-aware pencere check; global_mute_until expire; critical kind allowlist config `notif.critical_kinds`.
**Tarih:** 2026-09-09

---

## DD-173 — Android FCM `high` Priority Doze Bypass Yalnızca Kritik Kind İçin
**Karar:** FCM message `android.priority` allocation:
- **high**: `activity_completed, dungeon_completed, travel_completed, whisper_received, world_event_*, guild_event, party_invite, guild_invite, mail_received, system_broadcast, account_security`
- **normal**: `craft_completed, energy_restored, friend_online, market_expiring, achievement_unlocked`

**Reason:** FCM high priority Doze/App Standby'ı aşar ama Google quota'lı (abuse abandonment risk); non-time-critical (energy, friend, market_expiring) normal priority batched delivery.
**Alternatives:** Hepsi high (quota abuse + battery drain), hepsi normal (activity_completed delivery 30dk+ gecikir).
**Consequences:** iOS eşdeğer `apns-priority=10` vs `5`; priority tune Phase 02 A/B test open rate ile.
**Tarih:** 2026-09-09

---

## DD-174 — Mobile Lifecycle: Foreground Service YOK v1 (Battery + Store Policy)
**Karar:** V1'de Android **foreground service kullanılmaz** (persistent connection için background service çalıştırma yok). Client sadece FCM push wake ile aktive olur; server-side AFK engine tüm state'i tutar. iOS'ta `remote-notification` + `fetch` background modes; `audio/location/voip` YOK.
**Reason:** Battery/Play Store policy — foreground service persistent MMO connection için gereksiz (server AFK engine yeterli); user acquisition için Play Store rejection riski önlenir.
**Alternatives:** Foreground service (battery drain + Play policy warning), WorkManager periodic (5+ dk latency), push-only (v1 seçim).
**Consequences:** WS foreground'da bağlı; background transition 5 dk grace sonra WS close; reconnect exp backoff + `resume_from` replay son 5dk; Phase 03'te opt-in "high engagement mode" foreground service değerlendirme.
**Tarih:** 2026-09-09

---

## DD-175 — Notification Fallback = Inbox (Push Fail Etse Bile Ulaşır)
**Karar:** Her push mesajı için `notification.inbox_message` fallback yazımı (in-app notifications tab). Push provider fail veya kind mode 'inbox_only'/'off' olsa bile inbox delivery yapılır (unless kind explicit 'off'). Inbox retention 30 gün.
**Reason:** Push provider unreliable (FCM/APNs delivery guarantee yok); kullanıcı in-app açtığında kaçırdığı bildirimi görsün; retention safety net.
**Alternatives:** Fallback yok (kayıp bildirim + confusion), fallback her zaman (opt-out anlamsızlaşır) — bu yüzden 'off' mode inbox'ı da kesiyor.
**Consequences:** `inbox_message` schema unread badge; nightly cleanup > 30d; opt-out 'off' kesin sessizlik.
**Tarih:** 2026-09-09

---

## DD-176 — Flutter UI: Mobile-First + Riverpod + Data-Driven Reactive
**Karar:** Flutter client mobile-first (hedef viewport 375×812), tablet + web scale-up. State management **Riverpod** (öneri; final Phase 01 client team). Server hakikat kaynağı; UI reactive rebuild WS event ile (character/{id}/*, party/{id}, guild/{id}, inventory/{id}). Cached snapshot ile offline-first paint; server sync geldiğinde reconcile.
**Reason:** Mobile-first stack tutarlılığı; Riverpod compile-safe + testable + provider composition; data-driven UI backend'in tek hakikat prensibi (DD-135 ledger-first + character server-auth) ile hizalı.
**Alternatives:** BLoC (verbose), Provider (deprecated pattern), GetX (opinionated + lock-in), MobX (Dart tooling sınırlı).
**Consequences:** Provider mimari test-friendly; WS reconnect + resume pattern client team responsibility; cached snapshot secure storage (Keychain / EncryptedSharedPreferences).
**Tarih:** 2026-09-09

---

## DD-177 — Dark-First Design System + Rarity Palette + 4pt Grid
**Karar:** UI dark-first (light theme Phase 02). Color palette 12+ semantic role (bg×3, text×3, accent×2, state×3, rarity×6). 4pt spacing grid. Typography Inter (400/500/600/700) + Inter Mono (numeric HUD). Radius sm=6/md=10/lg=16/pill=999. Motion durations 150/200/300/400/600ms.
**Reason:** MMO gameplay uzun süreli oturum — dark theme göz yorgunluğu azaltır; rarity palette WoW/Diablo standardı oyuncu familiarity; 4pt grid Material/Apple HIG uyumlu.
**Alternatives:** Light-first (uzun oturum yorucu), custom fantasy font v1 (rendering hız), 8pt grid (mobile 375px'te dense değil).
**Consequences:** Tek token dosyası (design-system.md); Figma kaynak Phase 01 art team; theme switch API hazır (Phase 02 light eklemek trivial).
**Tarih:** 2026-09-09

---

## DD-178 — Deep Link Routing: game://route/... Her Ekran Adreslenebilir
**Karar:** Her UI ekranı deep link ile açılabilir: `game://route/character/activity/<id>/reveal`, `game://route/mail/<id>`, `game://route/guild/<id>/boss/<spawn_id>` gibi. Notification tap → route resolve; auth'suz açılırsa login sonrası queued deep link.
**Reason:** Push notification (DD-170) → deep link → cold start ekran açılışı 3 sn hedef; retention için kritik akış; universal link (iOS) + Android intent filter kolay.
**Alternatives:** Ekran sabit navigation stack (deep link yok — notification anlamsız), custom scheme parser client-only (auth loss).
**Consequences:** Router registry client-side; Universal Links + Android App Links config; queued deep link auth restore sonrası pop; test coverage her route için.
**Tarih:** 2026-09-09

---

## DD-179 — UI Optimistic Updates + Server Reconcile
**Karar:** Client-side aksiyonlar (inventory move, skill cast, chat send) optimistic UI update; server ack veya reject sonrası reconcile/rollback. Inventory drop/drag, skill cooldown, chat mesaj feed'i bu pattern'i kullanır.
**Reason:** Perceived latency < 100ms hissi; MMO'da UI donma frustration; server-authority korunur (reject rollback ile).
**Alternatives:** Full server round-trip (300-500ms latency hissiyat), local-only (drift + server ret cascade zor).
**Consequences:** Optimistic state overlay stream; rollback animation (200ms) + toast; test coverage her optimistic aksiyon; idempotency key zorunlu.
**Tarih:** 2026-09-09

---

## DD-180 — Bottom Tab Bilgi Mimarisi: 5 Tab (Home/Character/World/Combat/Social)
**Karar:** Ana navigation 5 bottom tab (mobile): **Home** (dashboard: aktif contract, energy, günlük ödül), **Character** (sheet+inv+equip), **World** (map+activity+travel+profession+crafting), **Combat** (combat+dungeon+party), **Social** (guild+friends+mail+chat+market+quest).
**Reason:** MMO 15+ ekran; 5-tab bottom nav Material/HIG max önerilen; her tab'in içinde sub-nav (tab drawer / segment control); overflow tab yok.
**Alternatives:** 3-tab (aşırı nested), 7-tab (bottom crowded), side drawer only (mobile bir eller kolay değil).
**Consequences:** Tab bar 5 slot; her tab kendi navigation stack; deep link her tab'i hedefleyebilir; onboarding tour tab hierarchy açıklar.
**Tarih:** 2026-09-09

---

## DD-181 — Admin Panel: Read-First, Server API Only, Full Audit
**Karar:** Admin panel default read-only; write aksiyonlar `/internal/ops/*` endpoint üzerinden ledger/audit yazar (DB'ye direct query yasak). Her sayfa view + her aksiyon `audit.admin_action`'a operator_id + target + before/after diff + reason ile yazılır.
**Reason:** Panel DB'ye direct write yaparsa ledger tutarlılığı (DD-135) kaybolur, drift kaçınılmaz olur; audit-first support forensic + compliance zorunlu.
**Alternatives:** DB direct query (drift + hesap verilemez), audit opsiyonel (compliance riski), yalnızca write audit (view leak forensic zayıf).
**Consequences:** Her admin endpoint domain servisi çağırır (mevcut `TxnService` gibi); audit retention 24 ay hot; correlation_id propagation her aksiyon zorunlu.
**Tarih:** 2026-09-09

---

## DD-182 — Admin Auth: SSO + WebAuthn 2FA + Step-Up Kritik Aksiyon
**Karar:** Admin login **SSO (OIDC)** (Google Workspace/Okta/Azure) + **WebAuthn 2FA zorunlu**. Kritik aksiyonlar (reversal, ban_permanent, disband_override, config_write) **step-up** WebAuthn re-verification (5dk TTL token). Session TTL 8h.
**Reason:** Admin panel oyun credentials'tan tamamen ayrı (compromise cascade önlemi); WebAuthn phishing-resistant; step-up en tehlikeli aksiyonlar için ek layer.
**Alternatives:** Password + TOTP (phishing riski), sadece SSO (2FA yok — SSO tenant compromise), step-up yok (kazara reversal riski).
**Consequences:** SSO provider ops karar; WebAuthn credential kaydı (hardware key veya platform authenticator); step-up token Redis TTL.
**Tarih:** 2026-09-09

---

## DD-183 — Admin RBAC: 6 Rol × Least-Privilege
**Karar:** V1'de 6 admin rol: `viewer, support_l1, support_l2, game_designer, ops_engineer, admin`. Her rol permission set açık (bkz. auth-rbac.md matrix). En az yetki prensibi — support_l1 mute<=24h ama ban yok; support_l2 reversal<=100k gold; admin tüm aksiyon + role grant.
**Reason:** Support ekibinin küçük hataları büyük etki yaratmasın; permission granular audit trail'e uygun.
**Alternatives:** Tek admin rol (blast radius), 15+ rol (bakım karmaşık), custom permission per user (yönetilemez).
**Consequences:** `ops.admin_role.permissions` JSONB; permission check middleware; role grant `admin` role only + audit.
**Tarih:** 2026-09-09

---

## DD-184 — PII Reveal Audit + Otomatik Re-Mask (60sn)
**Karar:** Player detail'de e-posta ve diğer PII **default masked**; admin `player.pii_reveal` permission ile açar; her reveal `audit.admin_action` (`action=pii_reveal`) satırı yazar; 60 sn sonra otomatik re-mask.
**Reason:** GDPR/KVKK gereklilik (PII minimizasyon); admin ekibi yanlışlıkla PII screenshot/paylaşım riskini azaltır; audit trail sızıntı forensic'i için kritik.
**Alternatives:** Her zaman açık (compliance riski), hiç açılmaz (support imkansız), sınırsız süre (kazara paylaşım artar).
**Consequences:** UI countdown timer; PII reveal frequency KPI (aşırı reveal manager review); export akışında PII masked default (unmask ayrı approve Phase 02).
**Tarih:** 2026-09-09

---

## DD-185 — Ledger Reversal Rate Limit 5/gün per Operator + Reason Zorunlu
**Karar:** `ledger.reversal` aksiyonu operator başına günde max 5 defa; reason 32+ char zorunlu; step-up gerekli. Aşım → admin escalation. Support_l2 sadece <100k gold; üstü admin.
**Reason:** Reversal ekonomiye direkt müdahale — abuse (arkadaşa iyilik) veya kazara mass reversal riskini yavaşlatır; audit + rate limit iki katman.
**Alternatives:** Sınırsız (drift + abuse), tek reversal/gün (support acil case tıkanır), reason opsiyonel (audit değersiz).
**Consequences:** `admin_ledger_reversal_total{operator}` metric; günlük rate reset UTC 00:00 (DD-016); reversal SLA runbook Phase 02.
**Tarih:** 2026-09-09

---

## DD-186 — Güvenlik: Server-Auth + Fail-Closed + Defense in Depth
**Karar:** Hiçbir client input tam güven görmez; her değer server tarafından yeniden hesap veya validate edilir. Belirsizlikte **fail-closed** (unknown role → viewer, invalid token → 401). **Defense in depth**: WAF + rate limit + auth + input validate + business rule + reconcile. Zero-trust internal (Phase 02 mTLS).
**Reason:** MMO'da client tampering (modded client, injected packet) sık; tek katman kırıldığında diğerleri yakalar; fail-closed silent-permissive'den güvenli.
**Alternatives:** Client trust (hile garantisi), fail-open (yanlış permission ile leak), tek katman (single point of failure).
**Consequences:** Her endpoint validation + rate limit + auth + business rule zincirinden geçer; reconciler drift yakalar; anomaly rule config'de hot reload.
**Tarih:** 2026-09-09

---

## DD-187 — Rate Limit: Token Bucket + IP/Account/Character 3 Boyut
**Karar:** Rate limit algoritma **token bucket** (Redis Lua atomic). Her istek 3 boyut düşer (IP + account + character); herhangi biri tükenirse 429. Circuit breaker CPU>85% / DB pool>90%'de non-critical endpoint reject.
**Reason:** Token bucket burst dostu (client kısa süre yoğun ok); 3-boyut NAT sorununu (aynı IP çoklu account) + bot flood'ı ayırır; circuit breaker cascade failure önler.
**Alternatives:** Fixed window (burst yok / spike sorun), tek boyut (NAT false positive veya account flood), CB yok (cascade collapse).
**Consequences:** Redis Lua script atomic; 429 response `Retry-After` header + client backoff jitter; circuit breaker Prometheus metric feed.
**Tarih:** 2026-09-09

---

## DD-188 — Idempotency Key Zorunlu Tüm Mutating Endpoint (24h Retention)
**Karar:** Tüm mutating endpoint (activity start/reveal, market, trade, mail, bank, craft, prestige, TxnService.Write invoke edenler) client `idempotency_key` (32 alphanumeric) zorunlu. `api.idempotency_record` 24 saat cache; aynı key + aynı request hash → cached response, farklı hash → `E_IDEMPOTENT_MISMATCH` (409).
**Reason:** Network retry safety (mobil bağlantı flaky); dupe/double-charge önleme; DD-135 ledger UNIQUE constraint ile hizalı.
**Alternatives:** Idempotency opsiyonel (client hatası dupe), retention 5dk (client offline retry olamaz), retention 30d (storage).
**Consequences:** Client SDK auto-generate key önerisi; retention worker cleanup > 24h; mismatch suspicious flag + security event.
**Tarih:** 2026-09-09

---

## DD-189 — JWT EdDSA + Refresh Rotation + Reuse = Session Family Revoke
**Karar:** Access token JWT **EdDSA (Ed25519)** algoritma, TTL 15 dk. Refresh token 30 gün, tek kullanımlık + rotation (yeni refresh her use'da); **reuse detection** aynı refresh 2. kez → tüm session tree (parent chain) revoke + inbox uyarı "unusual activity".
**Reason:** EdDSA RSA'dan hızlı + kısa key; rotation compromise window kısaltır; reuse detection çalınmış token evi mahvetmeye izin vermez.
**Alternatives:** RS256 (yavaş), long-lived access (compromise window uzar), rotation yok (çalınmış refresh sınırsız kullanım).
**Consequences:** Key rotation 90 gün + 30 gün overlap; HSM/KMS private key Phase 02; refresh_token tablo parent_token_id chain; UI'da tüm cihazlardan çıkış prompt.
**Tarih:** 2026-09-09

---

## DD-190 — Anomaly Detection Rule-Based v1 (Config Hot Reload)
**Karar:** V1'de anomaly detection **rule-based** (credential stuffing, brute force, geo impossibility, refresh reuse, RMT bulk trade, ledger drift, sum violation, chat spam, item flow spike). Config `security/rules.yml` hot reload; ML pipeline Phase 03.
**Reason:** Rule-based tuning + audit trail açık; ML modelin false positive/negatif interpret zor; v1 için yeterli.
**Alternatives:** ML v1 (data yok — cold start), rule yok (blind), commercial fraud API (cost + vendor lock).
**Consequences:** Config-driven rule engine; her tetiklenen kural `audit.security_event`; severity 4-tier (info/warn/high/critical); auto-action (IP block, session revoke, character freeze) + support queue.
**Tarih:** 2026-09-09

---

## DD-191 — Alert Sınıflandırma 4-Tier + PagerDuty Critical
**Karar:** Security alert 4 severity: `info` (log only), `warn` (inbox notify player), `high` (Slack #sec + admin dashboard), `critical` (PagerDuty on-call). Critical örnekler: ledger sum violation, mass ban trigger, drift ≥ 1c.
**Reason:** Sınıflandırma yorulma önler (her şey P0 → immune); on-call yalnızca gerçek kritik için uyanır.
**Alternatives:** Tek tier (P0 fatigue), 2-tier (nuance yok), 6+ tier (karar tembelliği).
**Consequences:** Alertmanager Prometheus routing per severity; incident runbook her tier için ayrı; post-mortem critical/high zorunlu.
**Tarih:** 2026-09-09

---

## DD-192 — Test Piramidi: Unit 70% / Integration 20% / Contract + E2E 10%
**Karar:** Test piramidi hedef: unit tests %70 (Go `testing` + `testify`), integration %20 (testcontainers-go real Postgres + Redis), contract + E2E %10 (JSON schema snapshot + k6 smoke). Core domain modülleri **coverage ≥ 80%**, ledger/RNG kritik ≥ 90-95%. CI fail on coverage drop > 2%.
**Reason:** Fast feedback (unit < 100ms) TDD comfort; integration real infra ile drift yakalar; E2E deployment smoke; her katman uygun maliyet/benefit.
**Alternatives:** All-unit (integration bug kaçar), all-E2E (yavaş + flaky), no coverage gate (regression riski).
**Consequences:** Testcontainers-go infra Phase 01; property-based test `gopter` critical path (RNG, ledger); coverage gate CI; fixture pack `testdata/fixtures/` content baseline.
**Tarih:** 2026-09-09

---

## DD-193 — Deterministic Test: Seeded RNG + Clock Injection Zorunlu
**Karar:** Tüm test'ler deterministic olmalı. RNG seed injection (`RNG(seed)`) + clock injection (`clock.Freeze(t)`); test aynı komuta 100× çalışsın aynı sonucu versin. Flaky test CI fail; retry değil root-cause fix.
**Reason:** Flaky test CI güvenini bozar (herkes "yine flaky" der, gerçek bug'ı kaçırır); deterministic = reproducible = debug fast.
**Alternatives:** Retry on failure (bug gizler), no clock injection (time-dependent flake), production RNG in test (non-repro).
**Consequences:** RNG interface injection tüm domain; `clock.Clock` interface (real + test impl); flake dashboard `test_flake_total{test_name}` metric (Phase 02).
**Tarih:** 2026-09-09

---

## DD-194 — Property-Based Testing Kritik Path (Ledger, RNG, Modifier Stack)
**Karar:** Property-based test (`gopter`) zorunlu path: ledger sum invariant (SUM(delta per currency)==0), RNG chi-square distribution, modifier stack apply order idempotent, item conservation across trades, skill point conservation across respec. CI: 1.000 iteration; nightly 10.000.
**Reason:** Table-driven test edge case kaçırır; property test random input + assertion invariant = daha geniş coverage; matematik invariant sağlam kalmalı.
**Alternatives:** Sadece table-driven (edge case blind), fuzz sadece (invariant assertion yok), no property (regression riski).
**Consequences:** `gopter` Go framework; nightly ayrı CI job (yavaş); property fail → seed dump reproducer.
**Tarih:** 2026-09-09

---

## DD-195 — Golden Fixture + Regression Snapshot Kritik Path
**Karar:** Kritik path (loot roll seed 101, AFK reveal mining 1h lvl 15, damage pipeline fireball lvl 30, crafting iron_sword seed 42) için **golden fixture** JSON snapshot. Output değişirse test fail + PR reviewer manual approve.
**Reason:** Formula veya balance değişiklikleri unintended regression yaratır; snapshot audit trail — "bu değişikliği bilerek yaptık" onayı.
**Alternatives:** Sadece assertion (kayması fark edilmez), snapshot yok (regression), snapshot auto-update (drift).
**Consequences:** `testdata/fixtures/*.json` git tracked; snapshot update PR review gerektirir; snapshot diff CI output.
**Tarih:** 2026-09-09

---

## DD-196 — Testcontainer Fresh DB Per Test Class + Parallel Safe
**Karar:** Integration test her test class için fresh Postgres + Redis container (testcontainers-go). Parallel run safe (isolate DB); shared state yok. Container reuse yalnızca aynı test class içindeki test method'lar arası (schema init 1×).
**Reason:** Real infra ile test şart (mock drift); fresh isolate flake önler; parallel run CI wall-clock kısa tutar.
**Alternatives:** Shared DB (test order dependency), mock DB (SQL drift), sequential run (CI yavaş).
**Consequences:** CI runtime memory + docker; test class başına ~5s setup; parallel test count = CI worker count; nightly heavy integration slower.
**Tarih:** 2026-09-09

---

## DD-197 — Performance: SLO-Driven + Realistic Mix + Nightly Baseline Regression
**Karar:** Her endpoint için latency (p50/p95/p99) + error rate SLO tanımlı; error budget aylık track. Load test **realistic user mix** (active:afk = 3:7) ile; steady + burst kombine; nightly main branch baseline; regression > %5 CI fail.
**Reason:** Test kod değil, tüketici hedefli olur; realistic mix bot-script'ten üstün gerçekçi bottleneck yakalar; regression sürekli fark edilir (spec değişikliği unintended perf drop yaratabilir).
**Alternatives:** Ad-hoc load (regression gecikir), synthetic hafif mix (bottleneck kaçar), sadece steady (burst weakness gizlenir).
**Consequences:** k6 script `perf/k6/`; Prometheus + Grafana dashboard; nightly CI compare + Slack report; error budget tüketimi release freeze tetikler.
**Tarih:** 2026-09-09

---

## DD-198 — CCU Tier Planning: T1 1k → T4 50k (V1 Launch), T5 100k Phase 03
**Karar:** V1 launch capacity hedef **T3 (10.000 CCU)**; peak season / event day için **T4 (50.000)** hazırlığı zorunlu. T5 100k Phase 03 (multi-region / shard karar).
**Reason:** İdeal launch scale realistic; peak event surprise oldu ise T4 auto-scale ile geçici scale; T5 için mimari değişim (shard) v1 kapsam dışı.
**Alternatives:** T3 kilit (event day crash), T5 v1 (aşırı capex), no target (belirsiz sizing).
**Consequences:** Capacity planning per tier ([capacity-planning.md](performance/capacity-planning.md)); auto-scaling threshold Phase 02 tune; T5 shard tasarım Phase 03 karar.
**Tarih:** 2026-09-09

---

## DD-199 — User Behavior Mix Hipotezi: Active:AFK = 3:7 (Peak)
**Karar:** Load test'lerde varsayılan user mix **30% active + 70% AFK** (peak saat). Non-peak 20:80. Bu oran v1 hipotez; Phase 02 real telemetry ile revize.
**Reason:** AFK-first MMO idle genre benchmark; AFK player'ın network overhead düşük (ping + occasional reveal) → capacity budgeting rasyonel.
**Alternatives:** 50:50 (over-estimate active load; over-provision), 10:90 (under-estimate active; under-provision), no assumption (test tasarımı imkansız).
**Consequences:** Test fixture bu oranla generate; SLO validation bu mix ile; Phase 02'de telemetry-driven revize.
**Tarih:** 2026-09-09

---

## DD-200 — WebSocket Node Sizing: 15k Conn Per Node + Redis Pubsub (v1)
**Karar:** V1'de WebSocket connection cap per app node **~15.000 conn** (memory 700-1000 MB, CPU 30-40% baseline). Multi-node fanout Redis pubsub cluster-wide. NATS migration Phase 03 karar (pubsub scale limit gördüğümüzde).
**Reason:** Go net/http + gorilla/websocket 15k conn per node ölçek denenmiş; Redis pubsub v1 için yeterli (~2k msg/sn); NATS ek complexity v1'de gereksiz.
**Alternatives:** 30k/node (memory + GC pressure), NATS v1 (yeni infra), sticky session LB (WS reset drama).
**Consequences:** T3 (10k conn) 1 node; T4 (50k conn) 3-4 node; multi-node fanout Redis cluster; message compression Phase 02.
**Tarih:** 2026-09-09

---

## DD-201 — Error Budget Tüketimi > 80% → Release Freeze
**Karar:** Aylık error budget tüketimi > %80'e ulaşırsa yeni feature release **donduruldu**; sadece hotfix + reliability work merge edilir. > %100'de PagerDuty critical + post-mortem.
**Reason:** SLO'ya bağlılık soyut sözden ekşiyip somut karara dönmelidir; feature velocity ile reliability arasında explicit tradeoff.
**Alternatives:** Sadece uyarı (aksiyon yok), sıfır tolerans (aşırı katı), post-hoc analiz (proactive değil).
**Consequences:** Alertmanager threshold + Slack #eng notify; freeze policy runbook Phase 02; reliability work backlog priority-boost.
**Tarih:** 2026-09-09

---

## DD-202 — Production Infra: Immutable + Everything-as-Code + GitOps
**Karar:** Prod ortam **immutable infrastructure** (SSH ile canlı değişiklik yasak); tüm konfigürasyon kod olarak Terraform (infra) + Helm/Kustomize (deploy). GitOps (ArgoCD) Phase 02. 12-factor + secret manager env inject.
**Reason:** Drift eliminate; reproducibility (aynı commit → aynı infra); audit trail git history; rollback trivial (revert commit).
**Alternatives:** Manuel config (drift + audit yok), Ansible push (state divergence), no IaC (compliance sıkıntı).
**Consequences:** Terraform state remote (S3 + DynamoDB lock); Helm chart per service; secret Vault; ArgoCD Phase 02 auto-sync + manual approve prod.
**Tarih:** 2026-09-09

---

## DD-203 — Container: Distroless Nonroot + Multi-Stage + Cosign Signed + Trivy Scan
**Karar:** Prod container image: distroless nonroot base, multi-stage build (Go static binary), cosign imzalı, Trivy vulnerability scan CI'de; critical CVE fail build. Registry retention prod tag 1 yıl (rollback için).
**Reason:** Minimum attack surface (no shell, no package manager), signed image supply chain security, scan preemptive CVE detect; rollback tag korunmalı.
**Alternatives:** Alpine base (glibc/musl uyum sorunu), Ubuntu base (büyük + CVE yüzeyi), no signing (supply chain attack risk).
**Consequences:** Image size < 50 MB; cosign KMS-backed key Phase 02; Kyverno admission Phase 02 (imzasız/CVE image reject); registry retention policy config.
**Tarih:** 2026-09-09

---

## DD-204 — Backup: Postgres WAL Continuous + Nightly Base + PITR + Aylık Restore Drill
**Karar:** Postgres backup: continuous WAL archive S3 (15dk RPO) + nightly base backup + PITR capability. Redis snapshot hourly + AOF. **Zorunlu aylık restore drill** fresh env üzerinde. Cross-region backup bucket (regional outage protection).
**Reason:** Data corruption/exploit/ransomware'e karşı zaman-noktasına dönüş; restore drill drill test → real incident hızlı; cross-region regional outage total loss önler.
**Alternatives:** Sadece daily backup (24h RPO — kabul edilemez), no drill (restore çalışmıyor bulunur ancak incident'te), tek region backup (regional outage risk).
**Consequences:** wal-g/pgBackRest tool; managed RDS backup 7-35 gün retention; monthly drill runbook; RTO 1-2 saat / RPO 15 dk hedef; S3 Object Lock ransomware defense.
**Tarih:** 2026-09-09

---

## DD-205 — Monitoring: Prometheus + Grafana + OpenTelemetry Tempo + Correlation ID
**Karar:** Metrics **Prometheus** (pull scrape 15s); traces **OpenTelemetry SDK + Jaeger/Tempo** (10% sampling baseline, 100% error/slow); dashboards **Grafana** (unified metric+log+trace). Trace ID = correlation_id her log satırında; Tempo one-click trace ↔ log navigate.
**Reason:** Prometheus endüstri standard + hafif; OpenTelemetry vendor-neutral; correlation_id incident debug için hayati (audit action + log + trace tek pencere).
**Alternatives:** Datadog (kapalı ekosistem + pahalı), self-metrics (invent wheel), no correlation (incident debug cehennem).
**Consequences:** Prometheus long-term Thanos/Cortex Phase 02; Grafana Cloud vs self-hosted karar; Trace context W3C traceparent propagation zorunlu tüm service.
**Tarih:** 2026-09-09

---

## DD-206 — CI/CD Prod: Canary Argo Rollouts + Feature Flag + SRE Approval + Deploy Freeze Cuma pm
**Karar:** Prod deploy: canary rollout (Argo Rollouts) 5% → 30% → 100% with auto-analysis (error+latency); feature flag (GrowthBook öneri) deploy≠release; SRE approval gate + PR review; Cuma 15:00 sonrası ve maintenance window sırasında deploy freeze.
**Reason:** Canary blast radius küçük; feature flag risk-free release (toggle < 1dk rollback); SRE gate audit + son kontrol; Cuma pm freeze incident on-call'u yıpratmasın.
**Alternatives:** Big-bang deploy (yüksek risk), no flag (rollback code redeploy 5-10 dk), no approval (accidental deploy).
**Consequences:** Argo Rollouts K8s CR; GrowthBook self-hosted / LaunchDarkly SaaS; deploy metric DORA track; auto-rollback trigger tanımlı.
**Tarih:** 2026-09-09

---

## DD-207 — Rollback Karar Ağacı: Flag→App→Migration→Content→DR (MTTR Katları)
**Karar:** Rollback decision tree: **Feature flag toggle** (1dk) > **App version undo** (5dk) > **Reversible migration** (5dk) > **Content pack git revert** (10dk) > **Irreversible migration** = **Backup restore** (60dk) > **Full DR** (1-2h). Her rollback runbook wiki + quarterly drill.
**Reason:** İncident sırasında karar hızı kritik; sıralı decision tree operator'un doğru yol seçimini garanti eder; blameless post-mortem action item feed'i.
**Alternatives:** Ad-hoc rollback (yavaş + risk), tek yol (birkaç senaryo kapsanmaz), no drill (production'da öğreniriz — geç).
**Consequences:** Runbook wiki per scenario; status page (Statuspage / Betteruptime) update automation; DORA MTTR < 30 dk hedef; quarterly drill.
**Tarih:** 2026-09-09

---

## DD-208 — Release Pipeline: Kademeli Gate + Data-Driven + Reversible Her Adım
**Karar:** Launch pipeline kademeli 7 gate: **Production Readiness → Closed Beta W1 → Bug Fix + Balance → CB W2 → CB W3 → Open Beta → Release v1.0**. Her gate önceki başarı olmadan geçilmez; her karar metric + community feedback ile; her adım rollback path açık.
**Reason:** MMO launch tek-shot değil — gradual real-user feedback en iyi risk mitigation; kademeli capacity ramp (200 → 5k → 25k) infra scale doğrular; her wave öğrenilen ders sonraki'ne uygulanır.
**Alternatives:** Big-bang launch (yüksek risk), sadece OB (CB verimli data alamaz), no wipe (v1'de balance stability için gerekli).
**Consequences:** 3-5 ay Production Readiness → v1.0 timeline; wave arası 1-2 hafta bug+balance+perf cycle; Founder tester rewards Release'te retain.
**Tarih:** 2026-09-09

---

## DD-209 — Beta Wipe Politikası: W1→W2 Zorunlu, OB→Release Zorunlu, Release Sonrası Yok
**Karar:** Karakter wipe zorunlu noktalar: Closed Beta W1 → W2 (balance data invalidate), Open Beta → Release v1.0 (fresh start). Release sonrası **wipe YOK** (progression korunur). Founder tester rewards (title + cosmetic + starter currency) her wipe'ta grant edilir.
**Reason:** Balance test için fresh economic state şart; Release fresh start herkes eşit başlar (early tester avantajı yok — sadece cosmetic hatıra); post-launch wipe player trust katliamı olur.
**Alternatives:** No wipe (economy drift), sadece OB wipe (CB inconsistent), post-launch wipe (kesinlikle yasak).
**Consequences:** Wipe duyuru 1 hafta önce; Founder migration script pre-tested; blog + Discord + in-app banner communication.
**Tarih:** 2026-09-09

---

## DD-210 — Balance Content-Driven: Kod Değil, Content Pack + Config Hot Reload
**Karar:** Balance değişiklikleri **content pack update** (item stat, drop rate, XP curve, recipe cost, skill damage) veya **config hot reload** (fee, tax) ile yapılır. Formula değişimi (rare) yalnızca Phase 02 review sonrası kod deploy.
**Reason:** Content-driven balance minute-cycle (staging → prod dakikalarla); kod deploy overhead + risk; formula değişiklik player trust sarsar.
**Alternatives:** Hard-coded balance (deploy gerekli), no versioning content (rollback zor).
**Consequences:** Content pack git versioned + signed; hot reload `content_version` header; balance pass haftalık cycle; formula change rare + review-gated.
**Tarih:** 2026-09-09

---

## DD-211 — Launch War Room 24h + T-24h Freeze + T+72h Retrospective
**Karar:** Release launch günü **24h war room** (Discord voice + Zoom); attendees SRE lead + eng lead + PO + community + support + DBA; 8h shift × 3. **T-24h freeze** (non-critical change ban), **T+72h retrospective**.
**Reason:** İlk 24 saat en kritik (peak CCU + tespit edilmemiş bug); freeze surprise değişikliği önler; retrospective öğrenilen dersleri Phase 02'ye taşır.
**Alternatives:** No war room (incident response yavaş), no freeze (deploy chaos), no retrospective (öğrenim kaybı).
**Consequences:** On-call rotation reinforced 3 shift; pre-approved hotfix list; deploy freeze T-24h → T+72h enforced tooling; post-launch report public opt.
**Tarih:** 2026-09-09

---

## DD-212 — Release v1.0 Free-to-Play + Cosmetic Store Phase 02
**Karar:** Release v1.0 **free-to-play**; monetization YOK (payment integration Phase 02). Phase 02'de cosmetic store (mount skin, tabard, aura), gem currency activate, Founder Pack pre-order. **No pay-to-win**, **no gambling**; transparent odds bundle.
**Reason:** Launch scope minimize; monetization complexity Phase 02'de dedicated pass; ethics first (idle MMO player trust critical).
**Alternatives:** Store day-1 (scope + payment infra risk), Battle Pass v1 (design + balance load), pay-to-win (community suicide).
**Consequences:** Store Phase 02 launch; Stripe / Play Billing / App Store IAP; TR local (PayTR/Iyzico) opt; compliance ethics doc; retention KPI free-to-play retention model.
**Tarih:** 2026-09-09

---

## DD-213 — Admin v1: Ayrı Hesap + Argon2id + TOTP 2FA (SSO/WebAuthn Ertelendi)
**Karar:** Faz 01 admin backend'i (`/internal/admin`, `/internal/ops`) SSO (OIDC) + WebAuthn yerine oyuncu hesaplarından ayrı `ops.admin_user` + argon2id şifre + TOTP (RFC 6238) 2FA kullanır. Step-up (reversal, kalıcı ban, disband override, rol verme) taze bir TOTP kodu ile session'a 5 dk'lık `stepup_until` yazar. Oturum opak token (SHA-256 hash DB'de), TTL 8 saat. Rol→izin matrisi kodda (docs/admin/auth-rbac.md §4). İlk admin `cmd/adminctl` ile oluşturulur. React admin arayüzü ayrı iş.
**Reason:** SSO sağlayıcı seçimi (DD açık karar) ve WebAuthn cihaz kaydı Faz 01'de mevcut değil; TOTP dış bağımlılık gerektirmeden ikinci faktör sağlar; audit + RBAC + step-up modeli değişmeden SSO'ya taşınabilir.
**Alternatives:** SSO bekleyip admin'i ertelemek (moderasyon aracı olmadan beta riski), yalnız şifre (2FA yok — reddedildi).
**Consequences:** SSO/WebAuthn geldiğinde login + step-up doğrulayıcısı değişir, izin/audit katmanı aynı kalır; TOTP secret'ları DB'de düz saklanır (DB erişimi kısıtlı) — KMS şifreleme Phase 02.
**Tarih:** 2026-09-23

---

## DD-214 — Observability v1: Private /metrics Port + OTLP/HTTP + Collector-Side Tail Sampling
**Karar:** Prometheus metrikleri (`prometheus/client_golang`, global olmayan registry) API'de `METRICS_ADDR` (:9090), worker'da `WORKER_METRICS_ADDR` (:9091) ayrı listener'da; public router'dan geçmez. Etiketler chi route pattern'i (`/v1/character/{id}/sheet`), eşleşmeyen yol `unmatched` — ham path asla etiket olmaz. Trace: OpenTelemetry SDK + OTLP/HTTP exporter, yalnızca `OTEL_EXPORTER_OTLP_ENDPOINT` set ise aktif; W3C `traceparent` her zaman propagate edilir. İstemci `X-Request-Id` göndermezse correlation_id = trace_id (DD-205). Her `*Context` log satırı `trace_id`/`span_id` taşır. SDK head sampling varsayılan %100 (`TRACE_SAMPLE_RATIO`); DD-205'teki %10 baseline + %100 error/slow kuralı collector'da tail sampling ile uygulanır. Worker job'ları `worker.<job>` span'ı + run/latency/rows metrikleri üretir.
**Reason:** Head sampling hata/yavaş istekleri önceden bilemez; tail sampling doğru yer. OTLP/HTTP gRPC sunucusu gerektirmez. Ayrı metrics portu reverse proxy kuralına bağımlı olmadan iç ağda kalır.
**Alternatives:** `otelhttp` contrib (ek bağımlılık, route pattern'i sonradan bilmiyor), Prometheus default registry (test izolasyonu yok), `/metrics`'i ana router'a koymak (yanlış proxy config'inde sızıntı).
**Consequences:** pgx query span'ları + Redis span'ları Phase 02; dev'de `docker compose --profile observability` Prometheus + Jaeger getirir; prod'da collector + Tempo.
**Tarih:** 2026-09-23

---

## DD-215 — Flutter Client v1: Minimal Bağımlılık (ChangeNotifier + http + Navigator)
**Karar:** İlk client sprint'i (auth + karakter seçimi + character sheet) Flutter SDK dışında yalnızca `http`, `flutter_secure_storage`, `intl` (+ `flutter_localizations`) kullanır. State: `ChangeNotifier` (`AuthController`) + ekran-lokal `State`; routing: `Navigator` + auth gate; JSON: elle `fromJson`. docs/architecture/flutter.md'deki `riverpod`, `go_router`, `dio`, `freezed`, `mocktail`, `golden_toolkit` ihtiyaç doğduğunda (WS event'leri, deep link, 10+ ekran) eklenir. Ağ sözleşmesi dokümandaki gibi: Bearer + single-flight refresh (reuse detection'ı tetiklememek için), her POST'ta `Idempotency-Key` (refresh sonrası replay aynı key ile), `X-Request-Id`, 10 sn timeout, yalnız GET'te 3 deneme exponential backoff + jitter. Sadece refresh token secure storage'da; access token bellekte. `API_BASE_URL` `--dart-define` ile. Metin i18n (TR + EN `.arb`); Türkçe büyük harf sorunu (`i→I`) nedeniyle UI'da `toUpperCase` kullanılmaz. Server'a Flutter web için CORS middleware eklendi (`CORS_ALLOWED_ORIGINS`, dev'de localhost:* serbest, asla `*`).
**Reason:** Üç ekranlık bir akış için state/routing/codegen kütüphaneleri bakım ve build süresi ekler, pratik fayda sağlamaz; mevcut API sözleşmesi elle modellenecek kadar küçük.
**Alternatives:** Dokümandaki tam stack'i baştan kurmak (erken soyutlama), `bloc` (daha fazla boilerplate).
**Consequences:** Feature sayısı büyüdüğünde Riverpod + go_router'a geçiş ayrı DD ile; mevcut controller'lar provider'a taşınabilir şekilde UI'dan ayrık tutuldu.
**Tarih:** 2026-09-23

---

## DD-216 — CI v1: GitHub Actions + Coverage Ratchet (80% Hedefi Integration Suite ile)
**Karar:** CI sağlayıcı GitHub Actions (`.github/workflows/ci.yml`, cicd.md'deki açık karar kapatıldı). Üç job: **server** (gofmt, vet, `go test -race -cover`, paket başı coverage ratchet), **migrations** (Postgres 16 service: tüm migration'lar up → down -all → up, ardından `-tags integration` suite: content pack upsert ×2, tüm katalog yüklemeleri, ledger yazma + idempotent replay), **client** (dart format, flutter analyze, flutter test --coverage, toplam satır ≥ %80). Coverage gate `scripts/coverage_gate.sh`: server için `server/coverage-thresholds.txt` paket tabanları ölçülen değerin ~2 puan altında (docs/testing "drop > 2% fail"); yükseltilir, DD'siz düşürülmez.
**Reason:** Core domain paketlerinin çoğu Postgres-backed servis — memory-repo unit testleriyle %80'e çıkmaz; bugünkü değerler %2–58. Sahte bir %80 gate'i ya CI'ı kalıcı kırmızı bırakır ya da anlamsız testlerle şişirilir. Ratchet gerilemeyi hemen durdurur; gerçek artış integration suite'in genişlemesiyle gelir. Yerel Docker çalışmadığı için migration SQL'i ilk kez CI'da gerçek Postgres'e karşı doğrulanır.
**Alternatives:** Hemen %80 gate (tüm core paketler kırmızı), testcontainers-go (runner'da Docker gerekir; service container daha basit), coverage gate'siz CI (gerileme görünmez).
**Consequences:** Faz 01 çıkışında core domain (activity, economy, inventory, item, hunt, dungeon, crafting, market, character, auth) için %80 hâlâ zorunlu — integration suite servis katmanını kapsayacak şekilde büyütülmeli; eşikler her artışta yükseltilir. Race detector yalnız CI'da koşar (yerelde cgo/gcc yok).
**Tarih:** 2026-09-23

---

## DD-217 — Karakter Sahipliği Router Seviyesinde (IDOR Kapatma) + 5xx Detay Gizleme
**Karar:** `/v1/character/{id}/...` bloğunun tamamı `RequireCharacterParam` middleware'i ile korunur: `{id}` çağıranın hesabına ait değilse 404 `E_NOT_FOUND` (403 değil — varlık sızdırılmaz), handler'a hiç ulaşmaz. Body'de `character_id` alan vendor buy/sell/buyback aynı kontrolü handler'da (`vendor.CharacterOwner`) yapar. `apierr.Write` 5xx mesajlarını istemciye "internal error" olarak gönderir; detay correlation id ile loglanır.
**Reason:** Güvenlik incelemesi: items, inventory (view/move/split/delete), equipment (equip/unequip), buybacks handler'ları sahiplik kontrolü yapmıyordu (yorumlar var olmayan bir middleware zincirine güveniyordu); vendor buy başka oyuncunun altınını harcatabiliyor, sell/buyback başkasının eşyasını satabiliyordu. Her handler'ın kendi kontrolünü hatırlaması yerine tek bir guard unutulamaz. 18 handler 500'de ham `err.Error()` (SQL/driver detayı) dönüyordu.
**Alternatives:** Her handler'a ayrı kontrol (unutulmaya açık — sorunun kaynağı), `X-Character-Id` header'lı `RequireCharacter` (path id ile header uyuşmazlığı ayrı açık yaratır).
**Consequences:** Karakter kapsamlı her istek bir PK lookup ekler (gerekirse kısa TTL cache Phase 02). Contract-id'li rotalar (reveal/stop/complete/cancel) sahipliği serviste kontrol etmeye devam eder (doğrulandı). Regresyon testleri: middleware (foreign/unknown → 404, handler çağrılmaz), vendor (foreign character → 404), apierr (5xx detay gizli).
**Tarih:** 2026-09-23

---

## DD-218 — Kalıcı HP + Analitik Yenilenme + Tüketimlik Etkileri
**Karar:** HP artık karakterde saklanır (`hp_current`, `hp_updated_at`; NULL = dolu). Yenilenme tick'lenmez: okuyan taraf `min(max, kayıtlı + ⌊geçen_sn/2⌋ × (1 + 0.1×VIT))` hesaplar (docs/character/health.md §3.2 formülü; STA yerine VIT). Av/zindan güncel HP ile başlar (0 HP → `E_TOO_WOUNDED`), reveal'da savaşın bittiği an (av: `started_at + tamamlanan × cadence`, zindan: `expires_at`) sonucu HP'si geri yazılır — değer deterministik olduğu için MarkRevealed'dan önce yazılır ve tekrar güvenlidir. Tüketimlikler `content.item_template.use_effect` (`{kind: heal_hp, amount}`) ile tanımlanır; `POST /character/{id}/items/{item_id}/use` tek transaction'da eşya + karakter satırını kilitler, etkiyi uygular, stack'ten 1 düşer; aktif av/zindanda (`E_IN_COMBAT`) ve can doluyken (`E_HP_FULL`) reddeder. Idempotency middleware tekrar denemede ikinci birimi tüketmez.
**Reason:** DD-084 (AFK analitik) ve CLAUDE.md (timer/polling yok) ile uyumlu; Redis + 2sn tick modeli (health.md §2) AFK oyunda binlerce karakter için gereksiz yük. İksirler ve can barı bu karara kadar işlevsizdi.
**Alternatives:** Redis ephemeral HP + worker tick (sürekli iş, çökme anında kayıp), her savaşın tam canla başlaması (mevcut hal — iksir/ölüm anlamsız).
**Consequences:** MP/enerji/dayanıklılık henüz saklanmıyor (oynanışta kullanılmıyor); `potion_mp_small` bu yüzden kullanılamaz. Ölüm cezası (XP borcu, dayanıklılık kaybı, diriliş) Phase 02; şimdilik 0 HP'den analitik yenilenme. Aktif savaş sırasında sheet'teki HP yenilenmiş görünür, reveal'da savaş sonucu yazılır.
**Tarih:** 2026-09-23

---

## DD-219 — Hesap Kurtarma: Tek Token Tablosu + E-posta Outbox + Server Sayfaları
**Karar:** Şifre sıfırlama (1 sa) ve e-posta doğrulama (24 sa) tek `auth.account_token` tablosunda (`kind`), yalnız SHA-256 hash saklanır; tüketim koşullu UPDATE (tek kullanım) ve aynı türdeki açık kardeş token'ları da yakar. E-postalar `ops.email_outbox` → worker `email_dispatch` (stdlib net/smtp, 587 STARTTLS zorunlu / 465 TLS; `SMTP_HOST` yoksa dev'de log, prod'da gönderilmez) ile gider; gövde (tek kullanımlık link) gönderimden/vazgeçilmeden sonra NULL'lanır. Uç noktalar: `password/reset-request` (her durumda 202 + 400ms taban süre → hesap sızdırmaz; saatte hesap başı 3 e-posta, IP başı 5 istek), `password/reset-confirm` (tüm oturumları kapatır), `verify-email`, `resend-verification`, `password/change` (bu oturum hariç hepsini kapatır). Deep link henüz olmadığı için e-postadaki bağlantılar server'daki `/verify-email` ve `/reset-password` sayfalarını açar (dış kaynaksız, nonce'lı sıkı CSP, `no-referrer`, `no-store`). Doğrulanmamış hesabın karakter oluşturması `AUTH_REQUIRE_VERIFIED_EMAIL` ile açılır (varsayılan kapalı; SMTP kurulunca açılmalı).
**Reason:** docs/auth/password.md §5–8 ve register.md §7–8 gereksinimleri; iki ayrı tablo aynı yaşam döngüsünü tekrar ederdi. Outbox, API yolunda SMTP gecikmesi/hatası olmadan kalıcı yeniden deneme sağlar.
**Alternatives:** API'den senkron SMTP (gecikme + zamanlama sızıntısı + kayıp), harici e-posta SDK'sı (gereksiz bağımlılık), mobil deep link (henüz yok — Phase 01 sonu).
**Consequences:** Prod'da `PUBLIC_BASE_URL`, `SMTP_*` secret manager'dan verilmeli. `password/change` docs'taki "logout-all" yerine mevcut oturumu korur (UX; saldırgan oturumları yine kapanır). OTP/2FA ve yaygın şifre kara listesi (password.md §9) Phase 02.
**Tarih:** 2026-09-23

---

## DD-220 — Ustalık Bonusları Oyuna Etki Eder: Sheet Modifier Stack + Aktivite Snapshot'ı
**Karar:** Açılmış ustalık kademe bonusları (`content.mastery_track.tier_bonuses`, her kademe ≤ mevcut) iki yoldan etkili olur. (1) Sheet hedefleri (`attack`, `defense`, `crit_pct`, `hit_pct`, primary, kaynak max) `skill.CharacterAdapter.WithModifierSource` ile ekipman bonuslarıyla aynı DD-107 stack'ine girer; bu yüzden av/zindan combatant'ı, dünya etkinliği güç snapshot'ı ve sheet otomatik yansıtır. (2) `gather_yield_<kind>` hedefleri aktivite başlangıcında `progression.GatherBonus` ile okunup contract snapshot'ına dondurulur (DD-088): `yield_flat_per_hour` (saatlik ek kaynak) + `yield_mult` (additive yüzde). `ComputeYield` oranı `(base + flat) × (1 + mult)` olarak kullanır. Bonus kaynağı hata verirse sheet hesabı/aktivite başlangıcı hata döner (sessizce bonussuz snapshot üretilmez). Eski contract'larda alan yok → sıfır bonus, reveal deterministik kalır.
**Reason:** Sprint 24'te bonuslar yalnız API'de gösteriliyordu; oyuncuya vaat edilen ilerleme etkisizdi. Mevcut modifier stack ve snapshot mekanizması yeniden kullanıldı — yeni hesap yolu yok.
**Alternatives:** Reveal anında güncel ustalığı okumak (AFK sırasında kademe atlayınca geçmiş süreye geriye dönük bonus → DD-088 ihlali), bonusu gold çarpanı olarak uygulamak (item drop sayısını etkilemezdi).
**Consequences:** Her sheet hesabı bir ek indexli sorgu (`progression.character_mastery` by character_id). Bonus hedefleri migration seed'inde; content pack'e taşınınca target allowlist'i (`skill` target'ları + `gather_yield_*`) validate edilmeli.
**Tarih:** 2026-09-23

---

## DD-221 — Görevler v1: Content Pack + Progression Olay Hattı + Deterministik Günlük Seçim
**Karar:** Main/side/daily görevler `content/quests.yaml` → `content.quest_template` (contentpack validate: kategori, 1–5 hedef, `requires_quest` zinciri + döngü yok, dailies ön koşulsuz; hedef olay grameri `mob_killed[:mob]`, `boss_killed`, `dungeon_cleared[:dungeon]`, `gather_<kind>`, `item_crafted[:item]` ve kod referansları aynı paketteki mob/dungeon/item'a karşı doğrulanır). Karakter durumu `progression.character_quest` (active/completed/turned_in/abandoned/expired; açık örnek, main/side tek teslim ve günlük gün-başı tekliği partial UNIQUE ile). Hedefler mastery/başarımı besleyen aynı `progression.Recorder` hattından ilerler (`progression.Recorders` fan-out); hunt artık `mob_killed:<mob>`, crafting `item_crafted:<template>` olaylarını da yayar (mevcut eşleşmeler tam-eşleşme olduğu için çift sayım yok). İlerleme kabulden sonra başlar. Günlük teklif: seviyeye uygun havuzdan `fnv(karakter, UTC gün, kod)` en küçük 3 — saf fonksiyon, satır/scheduler yok; geçmiş günün açık günlüğü sayılmaz, bir sonraki kabulde `expired` olur (kaçan günlük = kaçan ödül). Teslim sırası: idempotent ledger (`quest-reward-<id>`, reward_source → karakter) → koşullu `completed→turned_in` UPDATE → yalnız kazanan istek XP verir (çift ödeme imkânsız; çok nadir çökme durumunda XP kaybı loglanır). NPC teslimi yok (v1 her yerden). Açık görev tavanı 25. WS `quest.completed`.
**Reason:** docs/quests/* v1 kapsamı; ayrı tick yolu yazmak yerine mevcut olay hattı ve ledger/idempotency desenleri yeniden kullanıldı. Deterministik seçim cron + tablo gerektirmez ve yeniden başlatmada değişmez.
**Alternatives:** Worker ile günlük satır üretimi (ek job + tablo + saat dilimi hataları), reveal içinde senkron görev güncelleme (reveal tx'ini büyütür), kabul öncesi ilerlemeyi saymak (geçmiş olayları saklamayı gerektirir).
**Consequences:** Eşya ödülü, itibar ödülü, `discover`/`deliver`/`use_skill` hedefleri, haftalık/etkinlik/lonca görevleri, NPC teslimi ve görev isimlerinin yerelleştirmesi (şimdilik koddan türetiliyor) sonraki iş. `Record` diğer progression kayıtları gibi best-effort (hata loglanır).
**Tarih:** 2026-09-23

---

## DD-222 — Sınıf Yetenekleri v1: Seviye ile Açılan Pasif + Otomatik Aktif
**Karar:** `content/skills.yaml` → `content.class_skill` (sınıf başı 2 aktif + 2 pasif). Karakter başına durum yok: sahip olunanlar sınıf + seviyeden türetilir (öğrenme/respec/uzmanlık Phase 02). Pasifler `skill.CharacterAdapter` modifier kaynağı olarak DD-107 stack'ine girer (ModifierSource artık base sheet'i de alır → sınıf/seviye için DB sorgusu yok). Aktifler av/zindan başlangıç snapshot'ına dondurulur (DD-088) ve deterministik tur tabanlı savaşta oyuncu sırası geldiğinde öncelik sırasıyla cooldown'u dolmuş ve MP'si yeten ilk yetenek kullanılır; aksi hâlde temel saldırı. Etkiler: `strike` (normal isabet/kritik zarı, hasar × power/100) ve `heal` (HP max'ın power%'i; yalnız HP ≤ %50, zar çekmez). MP her karşılaşmada dolu başlar (kalıcı MP yok); cooldown tur cinsinden ve karşılaşma başına sıfırlanır. Yeteneği olmayan snapshot'larda RNG akışı ve sonuç birebir aynı (eski kontratlar etkilenmez). Roller'lar artık Combatant'ı `Clone()` ile kopyalar (alan alan kopya yeni alanları düşürürdü). Content doğrulaması modifier tür/hedeflerini `skill.KnownKinds/KnownTargets`'tan alır (tek doğruluk kaynağı; eski listedeki `en_max` yazım hatası kaldırıldı). Public `GET /v1/classes/{class}/skills`.
**Reason:** docs/skill-engine/active-skills.md'deki tam cast pipeline'ı (GCD, kaynak türleri, interrupt, rotasyon editörü, manuel cast) AFK/analitik modelde karşılığı olmayan karmaşıklık; v1'de sınıf kimliğini ve güç artışını veren en küçük model seçildi, mevcut modifier stack ve snapshot yeniden kullanıldı.
**Alternatives:** Karakter başına öğrenilen yetenek tablosu + yetenek puanı (şimdilik ekonomi/UX yükü), zar tabanlı yetenek seçimi (determinizm korunur ama dengelemesi zor), manuel cast endpoint'i (AFK oyunda anlamsız).
**Consequences:** Denge (power/MP/cooldown) içerik pass'inde ayarlanacak. Mob yetenekleri, buff/debuff, DoT, kalıcı MP ve spec/respec Phase 02. Dünya etkinliği gücü pasifleri sheet üzerinden alır, aktifleri kullanmaz.
**Tarih:** 2026-09-23

---

## DD-223 — Kapsam Kapısı Birim + Entegrasyon Testlerinin Birleşik Profilinden Ölçülür
**Karar:** CI'daki paket başı kapsam ratchet'i (DD-216) artık yalnız birim testlerinin `go test -cover` çıktısını değil, Postgres'li job'da koşan birim + entegrasyon testlerinin `-coverpkg=./...` ile birleşik coverprofile'ını ölçer (`scripts/coverage_gate.sh profile`). Bir kod bloğu herhangi bir test binary'sinde çalıştıysa kapsanmış sayılır. Server job'u hızlı geri bildirim için race testlerini koşmaya devam eder; kapı Postgres job'unda. Eşikler ölçülen değerin ~2 puan altında tutulur ve yalnız yükselir; `go` modu yerel hızlı kontrol için kalır.
**Reason:** Domain servislerinin çoğu SQL'dir ve yalnız entegrasyon testlerinde çalışır; eski kapı bu kodu hiç saymıyordu (ör. quest birimde %38.9, birleşik %86.5). CLAUDE.md'deki "%80 çekirdek kapsam integration suite büyütülerek karşılanacak" hedefi ancak birleşik ölçümle izlenebilir. Ayrıca 846fa4d'deki CI düşüşü, eşiklerin Go 1.27 ile ölçülüp CI'da Go 1.25 ile değerlendirilmesinden ve 2 puan payın bırakılmamasından kaynaklandı.
**Alternatives:** Entegrasyon testlerini her pakete `_integration_test.go` olarak dağıtmak (her paket kendi kapsamını raporlar ama ortak fixture'lar kopyalanır), testcontainers ile birim job'unda DB açmak (CI süresi + Docker-in-Docker).
**Consequences:** Kapı Postgres'siz yerelde birebir üretilemez; yerel komut: `DATABASE_URL=... go test -tags integration -coverpkg=./... -coverprofile=cover.out ./... && sh scripts/coverage_gate.sh profile server/coverage-thresholds.txt cover.out`. %80 hedefi paket paket entegrasyon testleriyle kapatılacak.
**Tarih:** 2026-10-08

---

## DD-224 — Vendor Alım/Geri Alımı Kaynak Başına Tek Eşya Üretir (DB Tekilliği)
**Karar:** Vendor alımı ve geri alımı eşyayı kaynak referansı başına bir kez üretir: alımda ledger txn'i, geri alımda buyback kaydı. Migration 0028 `inventory.item_instance(source_ref_id)` üzerinde `source_kind IN ('vendor_buy','vendor_buyback')` için partial UNIQUE index ekler; servis önce aynı kaynaktan üretilmiş eşyayı arar, eşzamanlı ikiz index'e takılırsa mevcut eşyayı döndürür (`item.Service.BySource`). Geri alımın ledger anahtarı istek değil kayıt başınadır (`vendor-buyback-<entry>`), "claimed" işareti en son atılır.
**Reason:** Gerçek DB testleri iki hata buldu: (1) aynı Idempotency-Key ile tekrarlanan alım ledger'da bir kez ücretlendirip her seferinde yeni eşya üretiyordu (`MaterializeDrops` idempotent değildi; yorum `Grant`'in idempotent olduğunu söylüyordu) — üretimden sonra/yanıt kaydından önce çökme ya da eşzamanlı çift gönderim bedava eşya demekti. (2) Geri alımda farklı anahtarlı iki eşzamanlı istek iki kez ücret alıp tek eşya veriyordu; claim üretimden önce olduğu için başarısız üretim ödenmiş ama verilmemiş eşya bırakıyordu.
**Alternatives:** `MaterializeDrops`'u genel olarak kaynak başına idempotent yapmak (avlar aynı kaynak için birden çok eşya üretir — bozar), yalnız uygulama düzeyinde kontrol (eşzamanlılıkta yarış kalır), advisory lock (DB index'inden daha karmaşık).
**Consequences:** 0028 mevcut bir yinelenmiş vendor eşyası varsa uygulanamaz (prod'da yok; dev DB'de testin ürettiği tek kopya silindi). Eşzamanlı ikizler hata yerine aynı sonucu alır. Regresyon: `TestVendorRules`, `TestVendorConcurrentMintsOnce`.
**Tarih:** 2026-10-08

## DD-225 — İstemci IP'si Yalnız Güvenilen Proxy'den Gelen X-Forwarded-For ile Belirlenir
**Karar:** `RealIP` middleware'i (CORS'tan hemen sonra) istemci IP'sini istek başına bir kez çözer ve `httpctx.ClientIP` ile sunar; oran sınırı, auth oturum IP hash'i ve admin girişi bunu kullanır. Forwarded başlıklar yalnızca bağlantının karşı ucu `TRUSTED_PROXIES` (IP/CIDR listesi) içindeyse okunur; `X-Forwarded-For` sağdan sola gezilir, güvenilen hop'lar atlanır, ilk güvenilmeyen adres istemcidir. Liste boşsa (varsayılan) her zaman bağlantı adresi kullanılır.
**Reason:** Eski `ClientIP` XFF'nin ilk girdisine koşulsuz güveniyordu: istemci her istekte farklı bir değer yazarak IP başı sınırları (login 5/dk kaba kuvvet koruması, kayıt, şifre sıfırlama, admin login) tamamen aşabiliyordu. LB'ler XFF'ye ekleme yaptığı için ilk girdi LB arkasında da istemcinin uydurduğu değerdir. docs/backend/middleware.md §3.2 zaten proxy allowlist'i öngörüyordu; auth paketindeki kopya `clientIP` de kaldırıldı.
**Alternatives:** Sabit "sağdan N hop" sayısı (topoloji değişince sessizce bozulur), chi `RealIP` (başlıklara koşulsuz güvenir), başlıkları edge'de silmeyi operasyona bırakmak (unutulursa açık).
**Consequences:** Prod LB/ingress arkasında `TRUSTED_PROXIES` ayarlanmazsa tüm istemciler LB adresinden tek IP görünür (IP sınırları ortaklaşır) — deploy kontrol listesinde. Testler: `TestRealIP` (sahte önek, güvenilen hop, bozuk girdi, IPv6).
**Tarih:** 2026-10-08

## DD-226 — Pazar İlanı Eşya Kaydından Bağımsız Yaşar
**Karar:** Migration 0029: `market.listing.item_instance_id` nullable, FK `ON DELETE SET NULL` (önceden `NOT NULL ... ON DELETE CASCADE`). `Listing.ItemInstanceID` `*uuid.UUID`.
**Reason:** Eşya emanetten çıktıktan sonra (iptal ya da satış) çantadaki aynı türden yığına birleşip silinebiliyor; CASCADE ilan satırını da siliyordu. Sonuç: iptal başarılı olduğu hâlde 404 `E_LISTING_NOT_FOUND` dönüyordu (istemci hata gösteriyordu), alıcı yığını birleşince satış kaydı satıcının listesinden yok oluyordu. İlan şablon/adet/fiyatı zaten kendi kolonlarında tutuyor; eşya referansı yalnız aktifken gerekli.
**Alternatives:** Birleşmeyi emanetten dönen eşyada kapatmak (çantada parçalı yığınlar), soft delete (tüm eşya sorgularına filtre).
**Consequences:** Kapanmış ilanların `item_instance_id`'si null olabilir (istemci kullanmıyor). Regresyon: `TestHTTPInventoryOps` (iptal + alıcı birleşmesi).
**Tarih:** 2026-10-08

## DD-227 — Ledger Yazımı Çakışmada Rastgele Gecikmeli Geri Çekilmeyle 8 Kez Dener
**Karar:** `PgTxnService.Write` SERIALIZABLE çakışmasında (40001) en fazla 8 kez dener; denemeler arasında üstel + rastgele gecikme (4 ms, 8 ms, … ×1–2; en kötü durumda toplam ≈ 0.8 sn), ctx iptaline uyar. DD-135'teki "3 deneme" bunun yerini aldı.
**Reason:** Gerçek DB yarış testi (aynı eşyayı 8 eşzamanlı listeleme) ücret iadelerinin kaybolduğunu gösterdi: denemeler beklemeden yapıldığı için aynı hesaba yazan istekler her turda yeniden çakışıyor, 3 deneme bitince telafi yazımı ("market: fee refund failed — manual reconcile needed") düşüyor ve oyuncu altın kaybediyordu. Aynı risk tüm telafi yazımlarında (lonca kuruluş iadesi, üretim iadesi) vardı.
**Alternatives:** Hesap başına advisory lock (tüm yazımları seri hâle getirir, kilitlenme riski), telafileri kuyruğa alıp worker'la yeniden denemek (yeni tablo + iş; şimdilik gereksiz), READ COMMITTED + satır kilidi (ledger değişmezlerini zayıflatır).
**Consequences:** Yoğun çakışmada istek gecikmesi en fazla ~0.8 sn artar. Bütçe yine tükenirse hata loglanır (aynı davranış). Regresyon: `TestMarketConcurrentSellSameItem`.
**Tarih:** 2026-10-08

## DD-228 — Ürün Adı: New World MMO (Kod Adı Otter Kalır)
**Karar:** Oyuncuya görünen ad "New World MMO": uygulama başlığı (TR/EN), Android etiketi, iOS görünen ad, web manifest/başlık, kurtarma e-postalarının konusu ve sayfaları, admin TOTP issuer, dev console. İç kod adı Otter kalır (Go modülü, paket/sınıf adları, veritabanı, Docker servisleri, dokümanlar).
**Reason:** Ürün sahibinin kararı. İç adı değiştirmek yüzlerce dosyaya dokunur, oyuncuya hiçbir şey göstermez.
**Alternatives:** Tam yeniden adlandırma (yüksek risk, sıfır kullanıcı faydası).
**Consequences:** "New World", Amazon Games'in yayımlanmış MMO'sunun tescilli markasıdır; mağaza yayını öncesi hukuki marka kontrolü gerekir (store/legal adımı). Ad tek yerden (l10n `appTitle` + platform dosyaları) değiştiği için gerekirse kolayca güncellenir.
**Tarih:** 2026-10-08

## DD-229 — Denge Geçişi v1: Ölçülen Hedefler, Seviye Büyümesi, Ekipman Merdiveni
**Karar:** Denge simülasyonla ölçülür ve testle sabitlenir: `internal/balance` her sınıfı her seviyede, içeriğin kendi moblarıyla ve o seviyenin ekipmanıyla, üretimdeki sheet + savaş koduyla dövüştürür (`go test ./internal/balance -v` matrisi basar). Hedefler (`TestBalanceTargets`): önerilen mobda seviye ekipmanıyla 1 saatlik av ölüm oranı ≤ %10; 5+ seviye üstü mob ≥ %25 ölüm; Sv 1 ≤ 15 dk, hiçbir seviye > 90 dk; zindan min seviyede her sınıf ≥ %30, +4 seviyede ≥ %90 temizler; sınıf farkı ≤ 40 puan. Mekanikler: (1) seviye büyümesi — sınıf başı yarım puan adımlı artış (`character.BasePrimary`); (2) saldırı sınıfın ana statından (savaşçı STR, okçu/haydut AGI, büyücü/rahip INT; `skill.ComputeSheetFor`), savunma = VIT (×2 değil); (3) av/zindanda karşılaşmalar arası dinlenme — DD-218 yenilenmesi aralık süresince (zindanda üçte biri); (4) seviye farkı kuralı — mob g seviye yüksekse oyuncu −4g isabet, −%6g savunma (`combat` tek yerde; seviyesi 0 olan eski anlık görüntüler etkilenmez); (5) zindan boss'u mob şablonunun 2× canı, 1.2× saldırısı; (6) XP eğrisi 50·L^2.2. İçerik: `content/gear.yaml` (5 sınıf silahı + 5 zırh parçası × 6 kademe, Sv 1/5/10/15/20/25), mob değerleri tek eğri (can 140+32L, saldırı 26+3L, savunma 1+L, xp 10·L^1.6) × arketip (tank/kırılgan/kaçamak/kritikçi), +2 mob (Sv 10, 22), zindan kompozisyonları seviye bandına çekildi. Ekipman kaynakları: dükkân 1. kademe, demircilik deri (beceri 1) / demir + demir silahlar (3–5) / pul (10), av zaferi başına %4 ve zindan boss'unda garanti (sıradan değil) düşüş — mobun kademesinden, karakterin sınıfına uygun (`item.GearPool/RollGear`, tek idempotent `Grant` partisi).
**Reason:** Ölçüm öncesi oyun Sv 12'den sonra çöküyordu: seviye karakteri güçlendirmiyordu, saldırı herkes için STR'den geliyordu, avda can hiç yenilenmiyordu (12 saatlik AFK av kaçınılmaz ölüm), ekipman Sv 15'te bitiyordu ve hiçbir yerden düşmüyordu; Sv 18'den itibaren rahip dışında her sınıf %100 ölüyordu. Demircilik tarifleri beceri 5 istiyordu ama beceri 1'den ilerletecek tarif yoktu. Kriptin Sv 3 girişi Sv 8 boss'una çıkıyordu.
**Alternatives:** Sayıları elle ayarlamak (ölçülemez, regresyon korumasız), mob gücünü oyuncu seviyesine ölçeklemek (ilerleme hissini yok eder), seviye farkını yalnız stat eğrisiyle sağlamak (dinlenme ve kademeli ekipman farkı eritiyordu — denendi).
**Consequences:** Tüm sayılar yer tutucu; hedefler değişirse önce test, sonra içerik. Eski seviye-farksız anlık görüntüler eski kurallarla oynar. Yeni içerik (ekipman, mob, zindan) `TestBalanceTargets`'tan geçmek zorunda. İçerik adları hâlâ koddan türetiliyor (sonraki adım: TR/EN adlar).
**Tarih:** 2026-10-08

## DD-230 — İçerik Adları İçerikte Yaşar, Sunucu Dile Göre Verir
**Karar:** Oyuncuya görünen içerik adları `content/i18n/<dil>.yaml` dosyalarında `tür.kod: Ad` biçiminde tutulur (eşya, ek, mob, zindan, bölge, meslek, dükkân, görev, yetenek, başarım, kategori, unvan, grup, ustalık, itibar kademesi, dünya/sezon etkinliği; `affix.format` eklerin dile göre dizilişi, TR: `{prefix} {name} ({suffix})`). `contentpack.Read` dosyaları okur ve doğrular: tüm diller İngilizceyle aynı anahtarlara sahip, boş ad yok, YAML'daki her eşya/mob/zindan/bölge/meslek/dükkân/görev/yetenek kodunun adı var. Sunucu son başarılı yüklemeyi bellekte tutar ve `GET /v1/content/strings?lang=` (oturumlu, diğer kataloglar gibi; ETag = içerik sürümü + dil, `max-age=300`; bilinmeyen dil → İngilizce) ile verir. İstemci oyuna girişte cihaz dilinde bir kez çeker (`ContentNames`); tablo yoksa ya da ad eksikse koddan türetilen ad (`humanize`) gösterilir.
**Reason:** İstemci tüm içerik adlarını koddan türetiyordu ("Iron Ore", "Main 01 First Hunt") — Türkçe arayüzde İngilizce, bazen anlamsız adlar. Adlar istemciye gömülseydi her içerik eklemesi uygulama sürümü isterdi; içerikte tutulunca içerik ve adı birlikte yayımlanır.
**Alternatives:** Her API yanıtına ad alanı eklemek (onlarca uç nokta, Accept-Language her yerde), istemci ARB dosyaları (içerik başına uygulama sürümü), içeriği DB'de çok dilli tutmak (gereksiz şema; içerik zaten YAML'dan geliyor).
**Consequences:** DB'de tohumlanan kodlar (başarım, unvan, grup, ustalık, etkinlik) için tamlığı entegrasyon testi, ek havuzları için `item` testi denetler. Türkçe ek dizilişi parantezli ("Keskin Demir Kılıç (Ayı)"); tam dil bilgisi uyumu sonraki iş. Yeni dil = yeni YAML dosyası.
**Tarih:** 2026-10-08

## DD-231 — Hesap Hakları: Koşul Kabulü + 13+, Kendi Kendine Silme ve Veri Dışa Aktarma, Politika Sayfaları
**Karar:** Kayıt `accept_terms` olmadan reddedilir (tek onay: 13 yaş üstü + Kullanım Şartları + Gizlilik Politikası; self-declared, docs/security/audit.md §COPPA); kabul edilen sürüm ve zaman hesaba yazılır (`terms_version`, `terms_accepted_at`, sürüm `auth.DefaultTermsVersion`). Silme (docs/auth/account.md §4.4): `DELETE /v1/me` + şifre → `deletion_requested_at`, tüm oturumlar kapanır, bilgi e-postası; 30 gün içinde **giriş yapmak iptal eder** (yanıtta `deletion_cancelled`; ayrı iptal ucu gerekmedi). Worker `account_purge` (saatlik) süresi dolanı anonimleştirir (e-posta/kullanıcı adı `deleted…`, şifre özeti kullanılamaz, `state='deleted'`), token/cihaz/bildirim/tercih/gönderilmemiş e-postaları siler, karakterleri mevcut karakter silme işine verir (eşyalar, lonca liderliği orada). Ledger ve denetim kayıtları anonim hesaba bağlı kalır. Dışa aktarma (§4.6, Phase 02'den öne alındı, e-posta yerine): `POST /v1/me/export-link` 15 dakikalık tek kullanımlık bağlantı → tarayıcı `GET /v1/me/export?token=` ile hesabın tüm verisini JSON indirir (`privacy.Export`, salt okunur tek anlık görüntü). Politika metinleri `content/legal/{terms,privacy}.{tr,en}.html`, `GET /terms` ve `/privacy` (dil `?lang` / Accept-Language; işletmeci `LEGAL_OPERATOR`, destek `SUPPORT_EMAIL`; boşsa sayfada "taslak" bandı). İstemci: kayıt onayı sunucuya gider + politika bağlantıları; hesap menüsü (şifre, verilerimi indir, şartlar, gizlilik, hesabımı sil).
**Reason:** Yayın kontrol listesinin NO-GO maddeleri (ToS + Gizlilik yayında, yaş kapısı) ve Google Play'in uygulama içi + web hesap silme şartı; KVKK m.11 / GDPR erişim ve silme hakları. Onay kutusu istemcide vardı ama sunucuya hiç gitmiyordu.
**Alternatives:** Dışa aktarmayı e-postayla göndermek (SMTP kurulu değil, dosya eki büyür), hemen sert silme (yanlışlıkla silmeye geri dönüş yok; doküman 30 gün istiyor), politika metinlerini istemciye gömmek (güncelleme sürüm ister, mağaza URL'si gerekir).
**Consequences:** Metinler taslaktır; yayın öncesi hukuki inceleme ve `LEGAL_OPERATOR`/`SUPPORT_EMAIL` zorunlu. Koşullar değişince sürüm artar; mevcut oyunculardan yeniden onay Phase 02. Yedeklerden silme (docs/infra/backup.md §13) yedek saklama süresiyle (≤ 35 gün) sağlanır. Eski istemciler `accept_terms` göndermediği için kayıt olamaz (giriş etkilenmez).
**Tarih:** 2026-10-10

## DD-232 — Tek Sunuculu Üretim Kurulumu, Yedekleme Tatbikatı, Geri Alma
**Karar:** v1 üretimi tek sunucuda Docker Compose (`deploy/prod/`): Caddy (otomatik TLS, tek açık port 80/443, WebSocket geçişi, HSTS, 1 MB gövde sınırı) → API; worker; Postgres ve Redis (şifreli) yalnız iç ağda; tek seferlik `migrate`; `backup` servisi günlük `pg_dump` (özel biçim, okunabilirlik kontrolü, 35 gün saklama, `LAST_SUCCESS`). API yalnız Caddy'nin sabit IP'sinden gelen `X-Forwarded-For`'a güvenir (DD-225). İmajlar git commit'iyle etiketlenir; `scripts/deploy.sh` (temiz checkout → build → migrate → up → sağlık bekleme → `deployed.log`), `scripts/rollback.sh [tag]` (yalnız uygulama; migration'lar ileri uyumlu — genişlet önce, daralt sonraki sürümde). `scripts/restore-drill.sh <dump>`: atılabilir konteynerde geri yükler, şema sürümünü, kirli migration'ı, hesap varlığını ve defterin para birimi başına dengesini (DD-135) denetler. İmaj düzeltmeleri: bağımlılıklar `go.sum` ile doğrulanır (önceden yalnız `go.mod` kopyalanıp `go mod tidy` çalışıyordu), içerik paketi imaja gömülür, distroless imajlarda `wget` olmadığı için sağlık kontrolü binary'nin kendisi (`/app healthcheck`), kökte beyaz liste `.dockerignore` (önceden `deploy/` içindeydi ve uygulanmıyordu; her derlemede ~325 MB istemci dosyası gönderiliyordu).
**Reason:** Yayın kontrol listesinin NO-GO maddeleri (yedek + geri yükleme tatbikatı, geri alma prosedürü) ve hesap gerektirmeden çalışan, tek komutla kurulan bir üretim ortamı. Dev'deki API konteyneri sağlık kontrolü yüzünden sürekli "unhealthy" görünüyordu.
**Alternatives:** Kubernetes/Terraform (docs/infra hedefi; tek sunucu için gereksiz karmaşa, bulut hesabı ister), WAL arşivi + S3 (15 dk RPO; nesne depolama hesabı açılınca eklenecek), yönetilen veritabanı (hesap ve maliyet).
**Consequences:** RPO 24 saat (günlük döküm) — oyuncu ilerlemesinin en fazla bir günü kaybolabilir; `deploy/prod/backups` sunucu dışına kopyalanmadıkça felaket kurtarma yok. Tatbikat aylık koşulmalı. Tek sunucu = tek arıza noktası; ölçek ihtiyacında docs/infra hedef mimarisine geçilir.
**Tarih:** 2026-10-10

## DD-233 — Android Sürüm İmzalama Depo Dışında
**Karar:** Sürüm derlemesi `android/key.properties` (+ yükleme keystore'u) varsa onunla imzalanır; ikisi de git dışıdır (`key.properties.example` nasıl oluşturulacağını gösterir). Yoksa yerel derlemeler hata ayıklama anahtarına düşer (bu derlemeler mağazaya yüklenemez). Play App Signing kullanılır: Google yayın anahtarını tutar, bizdeki yalnız yükleme anahtarıdır. `applicationId` (`com.otter.otter`, kod adı) mağazaya ilk yüklemeden önce kalıcı bir değere çevrilmelidir — ürün kararı.
**Reason:** Mağaza yayını imzalı derleme ister; anahtar ve şifreleri depoya girmemeli (Kritik Kural 6).
**Alternatives:** Anahtarı CI gizlisinde tutmak (CI'dan yayın kurulunca eklenebilir), her geliştiricinin kendi anahtarı (mağazada tek yükleme anahtarı olmalı).
**Consequences:** Yükleme anahtarı kaybolursa Play Console'dan sıfırlama gerekir — güvenli yedeklenmeli.
**Tarih:** 2026-10-10

## DD-234 — Oyun İçi Destek Talepleri
**Karar:** Oyuncu hesap menüsünden destek talebi açar (`POST /v1/support/tickets`: kategori bug/account/abuse/feedback/other, başlık 3–120, mesaj 10–4000 karakter, isteğe bağlı kendi karakteri). İstemci yalnız allowlist'li `client_info` gönderir (platform, dil, flavor, sürüm; anahtar başı 64 karakter); bilinmeyen anahtarlar atılır. Hesap başı en fazla 10 açık talep (hesap satırı kilidi altında sayılır) + saatte 5 oluşturma. Operatörler `/internal/ops/support/tickets` kuyruğunu okur ve yanıtlar/kapatır (`support.tickets` izni: support_l1, support_l2, admin); her okuma ve yanıt `audit.admin_action`'a yazılır. Yanıt oyuncuya `support_reply` bildirimi (push + inbox) olarak gider. Talepler veri dışa aktarımına girer, hesap silinince silinir.
**Reason:** Mağaza yayınından sonra oyuncunun hata/hesap sorununu iletebileceği ve yanıtı oyun içinde görebileceği bir kanal gerekiyor; e-posta adresi tek başına talep takibi ve kötüye kullanım şikâyeti için yetersiz.
**Alternatives:** Harici yardım masası (Zendesk vb. — hesap + maliyet + kişisel veri üçüncü tarafa), yalnız `SUPPORT_EMAIL` (takip yok, yanıt oyunda görünmez).
**Consequences:** Ekler (ekran görüntüsü) v1'de yok. Yanıt tek mesaj; konuşma dizisi gerekirse sonraki iş. Kuyruk yalnız ops API'si üzerinden (admin arayüzü yok).
**Tarih:** 2026-10-10

## DD-235 — Birinci Taraf İstemci Hata Toplayıcı
**Karar:** Uygulama yakalanmamış Dart hatalarını (`FlutterError.onError` + `PlatformDispatcher.onError`) kendi sunucumuza gönderir: `POST /v1/client/errors` (kimlik gerekmez — çöküş girişten önce olabilir; IP başı saatte 30). İstemci aynı hatayı bir kez, çalıştırma başına en fazla 10 rapor yollar, gönderim hatasını yutar. Sunucu e-posta, Bearer/JWT ve `token/code/password/key` sorgu değerlerini maskeler, mesajı 1000 / yığını 40 satır–8000 karakterle sınırlar, platform/sürüm/flavor'ı allowlist'ler. Raporlar parmak izine (tür + mesajın ilk satırı + ilk 3 çerçeve; sayılar, satır:sütun, adresler atılır) göre tek satırda toplanır (`ops.client_error`: sayaç + ilk/son görülme + son sürüm); 10 000 farklı grup tavanı, 30 gün görülmeyen grup worker `client_error_gc` ile silinir. Hesap ve IP saklanmaz. Operatör: `GET /internal/ops/client-errors` (`ops.client_errors`: support_l2, ops_engineer, admin; denetimli). Metrik `otter_client_errors_total{platform,kind}` (uyarı kuralları için). Uygulama sürümü `--dart-define=APP_VERSION`.
**Reason:** Yayın öncesi istemci hata görünürlüğü gerekiyor (production-readiness); Sentry ve Crashlytics hesap ister ve kişisel veriyi üçüncü tarafa taşır. Gruplama satır sayısını sınırlı tutar.
**Alternatives:** Sentry (hesap + DSN, SaaS), Firebase Crashlytics (Firebase projesi bekliyor; native çöküşleri yakalayan tek seçenek), yalnız log satırı (gruplama/sayaç yok, log toplama kararı DD-OPEN-11 açık).
**Consequences:** Native (NDK/iOS) çöküşler ve obfuscate edilmiş yığınların sembolleştirilmesi kapsam dışı — Crashlytics Firebase hesabıyla eklenecek. Anonim uç kötüye kullanılırsa yalnız sayaçlar şişer (tavan + oran sınırı).
**Tarih:** 2026-10-10

## DD-236 — Uyarı Kuralları ve Tek Sunucu İzleme
**Karar:** Prometheus uyarı kuralları depoda (`deploy/alerts.yml`) ve `promtool test rules` ile birim testli (`deploy/alerts_test.yml`, `scripts/check-alerts.sh`, CI `alerts` işi). Kurallar yalnız mevcut metriklerden: hedef erişilemez (critical), 5xx oranı >%5 (high) / >%1 (warn) — saniyede 0.2 istek hacim kapısıyla —, toplam p95 >500 ms (warn), worker işi hata (warn) / her çalışmada hata (high) / 3 aralık + 2 dk başarısız (high; replikalar arası `max by (job)`), defter mutabakatı kayma buldu (critical), DB havuzu >%90 (warn), filo genelinde dakikada >100 başarısız giriş (high), saatte >100 istemci hata raporu (warn). Durma tespiti için yeni metrikler: `otter_worker_job_last_success_timestamp_seconds`, `otter_worker_job_interval_seconds`. Üretimde `monitoring` profili: Prometheus (15 gün saklama) + Alertmanager, arayüzler yalnız sunucunun loopback'inde (SSH tüneli). Teslim v1'de e-posta; `alertmanager.yml` git dışı (SMTP şifresi), örnek dosya amtool ile doğrulanır; gürültü bastırma: yüksek olan düşük eşini, `InstanceDown{worker}` worker iş uyarılarını susturur. Dev Prometheus aynı kuralları yükler (teslim yok).
**Reason:** Yayın öncesi kesinti, defter bütünlüğü ve kimlik bilgisi saldırılarının birine ulaşması gerekiyor; kuralların test edilmesi yanlış alarm/sessiz alarm riskini düşürür. Tek sunucu kurulumuna (DD-232) ek altyapı gerektirmeyen en basit yol.
**Alternatives:** Yönetilen izleme (Grafana Cloud vb. — hesap), PagerDuty/Slack (hesap; alıcı olarak sonra eklenir), Grafana panoları (v1'de Prometheus arayüzü yeterli).
**Consequences:** Sunucu/disk/Postgres/Redis düzeyi metrikler yok (node/postgres exporter eklenince disk >%85, yedek tazeliği, sertifika süresi uyarıları eklenecek). Eşikler lansman varsayılanı; gerçek trafikle ayarlanmalı. Prometheus sunucunun kendisi çökerse uyarı da gelmez — dış uptime kontrolü (ör. `/live`) ayrıca önerilir.
**Tarih:** 2026-10-10

## DD-237 — Bağımlılık Açık Taraması ve Desteklenen Go Sürümü
**Karar:** CI her push'ta `govulncheck` (v1.8.0) koşar ve yalnız kodumuzun gerçekten çağırdığı açıklarda kırılır; derleme ve CI desteklenen en yeni Go ana sürümünün son yamasıyla yapılır (imaj `golang:1.27-alpine`, CI `go-version: 1.27.x` + `check-latest`), `go.mod` `go` satırı bağımlılıkların istediği asgari sürümdür. Statik analiz staticcheck (v0.8.1) ile ayrı `lint` işinde; staticcheck Go 1.27 export verisini okuyana kadar o iş Go 1.26 ile koşar.
**Reason:** İlk taramada 21 erişilebilir açık vardı (JWT başlık ayrıştırmada bellek, pgx, go-redis, OTel, x/net HTTP/2, grpc, stdlib net/http · crypto/tls · html/template); imajlar artık güvenlik yaması almayan Go 1.25'e bağlıydı. Yeni açık yayımlanınca CI'ın kırılması güncellemeyi zorunlu kılar.
**Alternatives:** Dependabot/Renovate (GitHub ayarı — kullanıcı açabilir; bu karar CI kapısıyla çelişmez), golangci-lint (tek ikili ama yapılandırma yükü; staticcheck kontrolleri yeterli), Go sürümünü sabit yamaya kilitlemek (yamaları kaçırır).
**Consequences:** Upstream'de yeni açık çıktığında ilgisiz bir push da kırmızı olabilir — beklenen davranış (önce güncelle). Yerel geliştirme Go sürümü de güncel tutulmalı (taramalar yerel toolchain'in stdlib'ini değerlendirir).
**Tarih:** 2026-10-10

## DD-238 — Yük Testi (k6) ve Argon2id Sınırı + Parametreleri
**Karar:** Yük senaryoları `perf/k6/` (dokümandaki araç), Docker'la `scripts/loadtest.sh` üzerinden: `mix.js` (aktif %30 HTTP — ccu-scenarios ağırlıkları, 300 oyuncu başına ~40 rps; AFK %70 WebSocket + dakikalık nabız + global sohbet fanout gecikmesi) ve `login.js` (sabit varış hızında giriş; ≤20/sn normal: p95 < 300 ms, hata < %0.5; üstü fırtına: fazlası hızlı 503 ile atılır, kabul edilenler p95 < 800 ms, başka hata yok). Eşikler slo-slis.md SLO'ları; kaçırılırsa koşu başarısız. Her sanal oyuncu kendi istemci IP'siyle gelir: koşu süresince API yalnız compose ağına `TRUSTED_PROXIES` verir, oran sınırları açık kalır (atlatma kodu yok). Argon2id: aynı anda en fazla `AUTH_HASH_CONCURRENCY` (varsayılan CPU sayısı) türetme; 500 ms içinde yer açılmazsa `ErrBusy` → 503 `E_BUSY` + `Retry-After: 1` (giriş, kayıt, şifre değişimi, hesap silme, admin girişi); bilinmeyen hesap yolu da aynı yanıtı verir (yük altında hesap varlığı sızmaz). Parametreler 64 MiB × 3 → OWASP tabanı 19 MiB × 2 × 1; bilinmeyen hesap zamanlama sahtesi servisin kendi parametrelerinden üretilir (sabit eski sahte, yeni parametrelerle hesap numaralandırma kahini olurdu). İstemci `E_BUSY` için okunur mesaj gösterir.
**Reason:** Ölçüm (bu dizüstünde): sınırsız 64 MiB×3 ile 50 giriş/sn API'yi 2.6 GB'a çıkarıp OOM ile öldürdü (Docker motoru da dondu). Sınırla bellek sabitlendi ama hash ~140 ms (yükte ~280 ms) olduğundan giriş SLO'su tutmadı. Yeni parametrelerle: 10/sn p95 59 ms (önce 365), tepe bellek 69 MB (önce 543); 4 slotla 60/sn fırtınada ~46/sn kabul, kabul p95 638 ms, %23 503, başka hata yok, tepe 194 MB. Karışım smoke (100 oyuncu): sheet p95 8 ms, envanter 9, reveal 60, dükkân 44, pazar 9, sohbet fanout p95 19 ms; hata 0, 429 yok.
**Alternatives:** 64 MiB×3'ü koruyup daha büyük sunucu (giriş maliyeti/kapasitesi 5× kötü), 46 MiB×1 (OWASP eşdeğeri; bellek 2.4× fazla), kuyruk sınırsız bekleme (fırtınada bellek/gecikme patlar), hesap başına yeniden hash'leme (henüz prod kullanıcısı yok — eski hash'ler zaten kendi parametreleriyle doğrulanıyor).
**Consequences:** T1 (1000 oyuncu) ve üstü bu makinede koşulmadı — yük üreticisi sunucuyla aynı 7.5 GB'ı paylaşıyor; ayrı yük makinesi/perf ortamı gerekir. Prod'da `AUTH_HASH_CONCURRENCY` × 19 MiB bellek ayrılmalı. Canlı muharebe (skill cast) uç noktası yok, karışımdan çıkarıldı.
**Tarih:** 2026-10-11

## DD-239 — Ücretsiz Barındırma: Render (API) + Supabase (PostgreSQL)
**Karar:** Kapalı test/oyun için ücretsiz kurulum: veritabanı Supabase Free (Frankfurt), API Render Free web servisi (Frankfurt, `render.yaml` Blueprint), Redis Render Free Key Value (yalnız iç ağ). Ücretsiz katmana uyum: (1) API açılışta gömülü migration'ları uygular (`MIGRATE_ON_START`; `server/migrations` `embed`, bağımlılıksız `db.Migrate` — golang-migrate'in `schema_migrations` tablosu, dosya başına tek transaction, advisory lock, CLI ile birbirinin yerine geçebilir); (2) worker işleri API sürecinde (`RUN_WORKER`, `worker.FromConfig` ortak kurulum), aynı anda tek iş; (3) hesap tabloları `auth` → `identity` şeması (Supabase `auth`'u kendi servisine ayırıyor, yazmaya izin vermiyor; eski migration'lar `identity` oluşturur, 0033 mevcut veritabanlarını taşır); (4) `schema_migrations`'ta RLS (Supabase Data API public şemayı yayınlıyor); (5) JWT tohumu Render `generateValue` ile, URL-safe/dolgusuz base64 de kabul. `DATABASE_URL` (Supabase session pooler, şifreli) yalnız Render panelinde, depoya girmez.
**Reason:** Kullanıcı telefondan, bilgisayardan bağımsız ve ücretsiz oynamak istiyor. Render Free'de arka plan işçisi ve deploy öncesi komut yok; Supabase'de `auth` şeması kullanılamaz.
**Alternatives:** Ücretli Render (Starter ×2 ≈ 14 $/ay; ayrı worker + pre-deploy migrate), VPS + `deploy/prod` (DD-232; ödeme gerekir), bilgisayar + tünel (bilgisayar açık kalmalı, adres değişir), golang-migrate kütüphanesi (bağımlılık ağacı; 90 satırlık çalıştırıcı yeterli).
**Consequences:** Bulunan hata: worker havuz kilitlenmesi — her iş kilit için bir, iş için bir-iki bağlantı tutuyor; açılışta 17 iş aynı anda havuzu (8) doldurup birbirini bekliyordu, API istekleri de bağlantı alamıyordu (ayrı worker'da 25'lik havuzda da teorik risk). `Runner.WithMaxParallel` + `ParallelFor(maxConns)=maxConns/3`; API modunda 1. Render benzeri yerel koşu (0.1 CPU, 512 MB, boş veritabanı): 33 migration, 5 sn'de hazır, 25 iş çalışması hatasız, giriş 0.57 sn, oyun istekleri 5–20 ms. Ücretsiz katman sınırları: 15 dk boşta uyku (~50 sn ilk istek; işler uyurken durur, uyanınca toplar), yedek yok, e-posta/izleme yok, birkaç düzine oyuncu kapasitesi. Yayında bulunanlar (2026-10-12, canlı: https://new-world-mmo-api.onrender.com, Supabase `new-world-mmo` Frankfurt): (a) Render Key Value Blueprint'te yalnız `connectionString` verir → `REDIS_ADDR` `redis://` URL'si kabul eder, log'a yalnız host:port; (b) `DATABASE_URL` hatalıyken pgx hatası değeri (şifreyi) olduğu gibi log'a basıyordu → hata metni değeri göstermez; (c) istek zinciri oyuncu → Cloudflare → Render kenarı (`::1`) → uygulama; yalnız `10.0.0.0/8` güvenilince herkes `::1`, sonra Cloudflare kenar IP'si görünüyordu (IP başı sınırlar herkese ortak) → `TRUSTED_PROXIES` loopback + Cloudflare aralıkları (yılda bir kontrol), istek log'una `client_ip`; doğrulandı: gerçek istemci IP'si. Uyumayı önleyen Supabase `pg_cron` işi kullanıcı onayına bırakıldı.
**Tarih:** 2026-10-12

## DD-240 — Eşya Görselleri: Uygulamaya Gömülü, Koda Göre Adlandırılmış
**Karar:** Eşya görselleri `client/assets/items/<template_code>.png` olarak uygulamaya gömülür (kare, şeffaf, 128×128). Eşleştirme yalnız dosya adıyla: açılışta asset manifest'i bir kez okunur (`ItemImages.load`), görseli olan kodlar kümeye girer; `ItemArt` görsel varsa onu ekran boyutunda çözerek (`cacheWidth`) gösterir, yoksa sınıf ikonu kalır. Nadirlik hücre çerçevesinin rengi, önek/sonekli eşyalar temel eşyanın görselini kullanır. Çanta, ekipman, eşya detayı ve dükkân listesi görseli gösterir.
**Reason:** Kullanıcı görselleri koddan isimle eklemek istiyor; ek eşleştirme tablosu veya sunucu değişikliği gerektirmeyen en basit yol. Eksik görsel hata üretmez (manifest kontrolü), görseller kademeli eklenebilir.
**Alternatives:** Sunucudan/CDN'den indirme (APK'sız güncelleme ama önbellek, sürümleme ve ağ maliyeti), içerik paketinde `icon:` alanı (kod adıyla aynı bilgiyi tekrarlar).
**Consequences:** Yeni görsel yeni APK gerektirir. 84 eşya × ~10 KB ≈ 1 MB uygulama boyutu. Görsel adı koddan farklı yazılırsa (büyük harf, boşluk, .jpg) eşleşmez ve ikon görünür.
**Tarih:** 2026-10-13

## DD-241 — Can İksiri Kademeleri; Mana İksirleri Mana Sistemine Kadar Ertelendi
**Karar:** Dört yeni can iksiri: Orta Boy (Sv 5, +150, 60b), Büyük Boy (Sv 12, +280, 1g 50b), Ekstra Büyük (Sv 20, +450, 3g 20b), Çok Büyük (Sv 28, +700, 6g); yığın 20. Miktarlar ilgili seviyede tipik can tavanının ~%40–60'ı (can = 100 + 10×VIT + ekipman). Orta ve Büyük Ironhold dükkânında; dördü de simyayla üretilir (beceri 5/10/20/30) — üst kademeler o ana dek hiçbir tarifte kullanılmayan Kral Çiçeği ve Anka Yaprağı'nı tüketir. Mana iksirleri (kullanıcının istediği Orta/Büyük/Ekstra Büyük/Çok Çok Büyük) eklenmedi: mana her savaşın başında dolduğu için (DD-222) etkileri olamaz; kullanıcı "şimdilik sadece can" dedi. Mevcut Küçük Mana İksiri dükkândan çıkarıldı, tarifi migration 0034 ile emekliye ayrıldı (loader yalnız upsert eder), şablon kalır (eldekiler satılabilir).
**Reason:** Oyuncunun işe yaramayan eşyaya altın harcamaması; can iksirleri seviyeyle ölçeklenmeliydi (tek iksir 60 can veriyordu).
**Alternatives:** Kalıcı mana (can gibi taşınır/yenilenir) veya savaşta otomatik iksir — mana iksirlerini anlamlı kılar, denge yeniden ayarı gerektirir; ileride.
**Consequences:** Mana sistemi değişince dört mana iksiri + Küçük Mana İksiri tarifi geri gelir (kodlar `potion_mp_{medium,large,xlarge,huge}` için yer ayrılmış sayılır).
**Tarih:** 2026-10-13

## DD-242 — Ganimet Yığınlama ve Kazanım Geçmişi
**Karar:** (1) Toplama/av/zindan sonuç ekranı aynı eşyaları (şablon + nadirlik + önek/sonek) tek satırda sayıyla gösterir ("Demir Cevheri × 5"; tek düşüşte sayı yok), en nadir üstte, sonra ilk görülme sırası. (2) Her toplanan kontrat için sunucu bir `reward_collected` bildirimi yazar: kaynak (activity/hunt/dungeon), hedef (iş türü/mob/zindan), altın, XP, yeni seviye, yığınlanmış eşya listesi (+ av zafer/yenilgi, zindan temizlendi/boss). Yalnız gelen kutusu (push yok), **okunmuş olarak** yazılır (oyuncu ganimet ekranını az önce gördü; rozet şişmez), kontrat başına tek kayıt (`reward:<contract>` dedup — tekrar açma çoğaltmaz). Bildirimler ekranında ayrı "Kazanımlar" sekmesi (`GET /v1/notifications?rewards=1`); normal liste bunları içermez, böylece sık toplama diğer bildirimleri son-100 penceresinden itmez. Gelen kutusu kuralı gereği 30 gün saklanır.
**Reason:** Kullanıcı isteği: ganimet tek tek satır yerine yığın halinde, kazanımlar bildirimlerde kalıcı liste olarak.
**Alternatives:** Ayrı kazanım tablosu + uç nokta (ikinci saklama/temizlik yolu; gelen kutusu zaten saklıyor, dedup ve 30 gün GC'si var), istemcide yerel kayıt (cihaz değişince kaybolur).
**Consequences:** Kazanımlar 30 gün sonra silinir (gelen kutusu `expires_at`); daha uzun saklama istenirse tür başı süre eklenir. Sekme son 100 kazanımı gösterir.
**Tarih:** 2026-10-13

## DD-243 — Cevher Merdiveni, Ergitme Mesleği ve Düşme Kademelerinin Tamamı
**Karar:** Madencilik kullanıcının 8 kademeli cevher listesine göre: Çok Yaygın (Kömür, Taş Cevheri) ve Yaygın (Bakır, Demir, Kalay, Kurşun, Çinko) → Sıradan havuz (çok yaygınlar ×2 ağırlık); Orta (Kobalt, Nikel, ×2) ve Değerli (Altın, Gümüş) → Seçkin; Nadir (Paladyum, Rodyum, Titanyum) → Nadir; Çok Nadir (Platin, İridyum, Osmiyum) → Destansı; Efsanevi (Orikalkum, Mitril) → Efsanevi; Mitik (Adamantit) → Mitik. 14 yeni cevher; mevcut cevherlerin seviyeleri merdivene hizalandı (Sv 1/1/5/10/15/20/25/30). Yeni **Ergitme** mesleği: 18 metal için `<metal>_bar` külçesi = cevher + kömür (yaygın 2+1 beceri 1–3, orta 2+1 8, değerli 2+2 12, nadir 3+2 18, çok nadir 3+3 25, efsanevi 4+3 32, mitik 5+4 40; altın maliyeti ve meslek XP'si kademeyle artar; külçe değeri girdilerin ~1.6 katı). Kömür yakıttır, Taş Cevheri'nin külçesi yok. **Düşme kodu düzeltmesi:** verim modeli sıradan/nadir/efsanevi sayar; `RollDrops` her buluntuyu %25 (efsanevide %15) ihtimalle bir üst nadirliğe (seçkin/destansı/mitik) çıkarır — önceden bu kademeler hiç kullanılmıyordu.
**Reason:** Kullanıcı isteği (cevher + külçe). Düzeltme: Bakır/Kalay Cevheri, Levrek, Karaağaç, Adamotu, Kurt Postu hiç düşmüyordu; Pul zırh ve Orta/Büyük Can İksiri tarifleri yapılamıyordu.
**Alternatives:** Verim modeline 6 ayrı sayaç (analitik model ve eski kontrat determinizmi değişir), malzemeye sabit nadirlik alanı (içerik + şema değişikliği).
**Consequences:** Dağıtım anında aktif toplama kontratlarının düşmeleri (aynı tohumla) yeni bölünmeyle hesaplanır. Külçeler henüz hiçbir üretimde kullanılmıyor — demircilik tariflerinin cevher yerine külçeye geçmesi ayrı karar. Yeni test: düşme tablosundaki her kod içerikte olmalı (bilinmeyen kod sessizce düşmüyordu). Eşya listesi artık `scripts/item_catalog.py` ile üretilir.
**Tarih:** 2026-10-13

---

## DD-244 — Toplama Meslekleri (Sv 1–100) ve Kaynak Seçimi
**Karar:** Kullanıcının kaynak listesi (Book.xlsx) esas alınır. Beş toplama türü (madencilik, şifacılık, odunculuk, balıkçılık, avcılık) birer **toplama mesleği**dir: seviye 1–100, `crafting.character_profession`'da tutulur (ilk kazanımda açılır, elle öğrenilmez — `content.profession.kind = 'gathering'`, migration 0035). `content/gathering.yaml` 150 kaynak tanımlar; her kaynak `gather: {kind, level, tier}` taşır (`content.item_template.gather_*`). Oyuncu başlangıçta kaynağı seçer (`target`; boşsa açılmış en iyi kaynak); seviye kapısı `E_GATHER_LEVEL` (409). Seçim kontrat snapshot'ına donar. Kademe T0–T7: hız çarpanı 1.0→0.3, nadirlik (T0–1 sıradan, T2–3 seçkin, T4 nadir, T5 destansı, T6 efsanevi, T7 mitik), birim başı meslek XP'si 8→90. Her buluntu seçilen kaynağın bir birimidir. Meslek eğrisi tüm mesleklerde `round(60 × 1.06^L)` (önceden `100 × L`). Av/zindan zaferleri, mob seviyesi × 4'e (tavan 100) kadar açılan hayvanların derisini bırakır (sıradanlar en iyi 4'ten, nadir/efsanevi en iyisi). Listede olmayan eski malzemeler emekli (şablonları duruyor), tarifler yeni kaynaklara bağlandı; ergitme 24 külçe (yeni: yıldızçelik, tungsten, elektrum, krom, bronz, alüminyum), ergitme beceri kapısı `cevher seviyesi − 13`. `GET /v1/character/{id}/gathering` meslek seviyeleri + kaynak listesi; reveal yanıtında `gather` (kazanılan XP, seviye, seviye atlama); kazanım geçmişi kaydına `gather_xp`/`gather_level`.
**Reason:** Kullanıcı isteği: belgedeki listelerdeki eşyalar meslek sınıflarına; seviye = toplama meslek seviyesi 1–100; liste esas. DD-243'ün kademe-rastgele havuzu oyuncuya hedef seçtirmiyordu.
**Alternatives:** Karakter seviyesiyle kapı (meslek derinliği yok), DD-243 havuzuna ekleme (liste büyüdükçe istenen kaynak neredeyse hiç düşmez), ayrı toplama tablosu (meslek altyapısı zaten var).
**Consequences:** Eski (hedefsiz) kontratlar reveal'da türün başlangıç kaynağını verir. Avcılık mesleği yalnız toplamayla ilerler (av/zindan derisi XP vermez). Zanaat meslekleri yeni eğriyle orta seviyelerde daha hızlı ilerler. Yeni kaynakların görselleri `client/assets/items/<code>.png` ile eklenebilir. İstemci: Toplama sekmesinde seviye/XP çubuğu, kaynak listesi (açılmış en iyi üstte, sonraki 2 kilitli), sonuçta meslek XP'si; Atölye toplama mesleklerini göstermez.
**Tarih:** 2026-10-15

---

## DD-245 — Çanta Sıralama/Çoklu Seçim, Bağlı Eşya Satışı, Bildirim Silme, Toplu Üretim, Anlık Seviye
**Karar:** (1) Çanta görünümü sıralanabilir (çanta düzeni / nadirlik / tür / ad / değer); yalnız gösterim, sunucudaki konumlar değişmez. (2) Uzun basıp sürükleyerek çoklu seçim (sürükleme ilk hücreye göre ekler ya da çıkarır; seçim modunda dokunma seçer/bırakır; kilitli eşya seçilmez); seçilenler toplu satılır ya da yok edilir (onaylı, istemci sırayla tek tek çağırır). (3) Bağlı eşyalar (ör. bir kez kuşanılmış BoE) dükkâna satılabilir; geri alınan eşya yeniden bağlanır (`vendor.Repository.Rebind`), böylece sat-geri al döngüsü bağı çözmez; pazar bağlı eşyayı hâlâ kabul etmez. (4) `DELETE /v1/notifications/{id}`: hesabın kendi gelen kutusu / kazanım kaydı girdisini siler (yabancı ya da bilinmeyen id → 0, hata değil); istemci her girdinin sol üst köşesinde çöp kutusu. (5) `POST /craft` `count` (1..100): tüm partinin malzemesi önce kontrol edilir (yetmezse hiçbir şey yapılmaz), sonra her üretim kendi garantileriyle tek tek; sonradan biri düşerse (ör. altın biter) parti durur ve yapılan kadarı döner (`count`). İstemci: malzeme/altının izin verdiği en fazla adede kadar −/+ / kaydırıcı / max seçimi. (6) Seviye atlaması anında yansır: kabuk karakter seviyesini her yeni sheet'ten bir `ValueNotifier`'da tutar; Macera panelleri (av, zindan, yolculuk, etkinlik kapıları ve varsayılan mob) onunla yeniden kurulur, toplama sekmesi meslek seviyelerini her veri değişiminde sessizce yeniler.
**Reason:** Kullanıcı istekleri. (3): kuşanıp çıkarılan eşya `E_ITEM_BOUND` ile satılamıyordu. (6): Macera sekmesi seviyeyi girişteki değerle sabit tutuyordu; yeni açılan zindan/bölge karaktere yeniden girmeden açılmıyordu.
**Alternatives:** Sunucu tarafı toplu satış ucu (istemci döngüsü v1 için yeterli, mevcut idempotent satışı yeniden kullanır), sunucuda kalıcı sıralama (eşya taşıma — gereksiz yazım), bağlı eşya satışında ceza.
**Consequences:** Bağlı eşya satışından gelen altın normal satış gibi (değer × çarpan). Toplu üretimde her birim ayrı defter kaydı. `vendor.ErrItemBound` artık dönmüyor (eski istemciler için eşlemesi duruyor).
**Tarih:** 2026-10-15

---

## DD-246 — Çanta 100 Slot; Dolu Çanta Yeni Eşyayı Reddeder
**Karar:** `inventory.BagSlotCap` 60 → 100 (DB sınırı zaten 240). `inventory.Service.HasRoom(char, şablon, adet)`: boş slot varsa evet; yoksa yalnız mevcut yığınların boş kapasitesi adet kadarsa evet (boş şablon = yalnız boş slot). `FreeSlots`. Eşya ekleyen her akış işlemden ve ödemeden ÖNCE kontrol eder, yer yoksa 409 `E_BAG_FULL`: dükkân alımı ve geri alım, üretim (her birim — toplu üretim yer bitince durur ve yapılanı döner), pazar alımı, posta eki alma (ek sayısı kadar boş slot), toplama başlangıcı (seçilen kaynak için yer), av ve zindan başlangıcı (boş slot). Kuşanılan eşyayı çıkarma zaten `E_BAG_FULL` veriyordu. İstemci: çanta dolunca sayaç uyarı renginde; hata metni "Çantan dolu.".
**Reason:** Kullanıcı isteği: 100 eşyalık çanta, 100/100'de üretim/satın alma vb. engellensin. Önceden dolu çantada yeni eşyalar slotsuz (görünmez) taşma olarak kalıyordu.
**Alternatives:** Doluyken her şeyi reddetmek (yığına sığan iksir alımını da engellerdi — gereksiz katı), taşmayı postaya göndermek (Phase 02).
**Consequences:** Kontrol ile ekleme arasında yarış (eşzamanlı iki istek) nadiren taşmaya düşer; taşma davranışı yedek olarak duruyor (yer açılınca çanta görünümünde yerleşir). Başlamış bir av/toplama sırasında çanta dolarsa teslimdeki ganimet de taşmaya düşer — teslim engellenmez (kontrat kilitlenmesin). Pazar ilanı iptali eşyayı satıcıya her durumda geri verir.
**Tarih:** 2026-10-17

---

## DD-247 — E-posta Teslimi Brevo HTTPS API'si + Şifre Değişti Bildirimi
**Karar:** `mail.BrevoSender` (`POST https://api.brevo.com/v3/smtp/email`, `api-key` başlığı, düz metin gövde). `BREVO_API_KEY` doluysa worker `email_dispatch` SMTP yerine bunu kullanır; gönderen `MAIL_FROM` (Brevo'da doğrulanmış adres; yoksa `SMTP_FROM`), ad `MAIL_FROM_NAME` (varsayılan "New World MMO"). Mevcut akışlar aynen: kayıtta doğrulama bağlantısı, şifremi unuttum → sıfırlama bağlantısı, hesap silme bildirimi. Yeni: şifre değiştirme ve sıfırlama sonrası güvenlik bildirimi ("Şifren değişti" — bağlantı/token yok; best-effort, kuyruğa yazılamazsa değişiklik geri alınmaz). `render.yaml`'a `BREVO_API_KEY` + `MAIL_FROM` (panelden girilir).
**Reason:** Kullanıcı isteği (kayıt, şifre unutma/değiştirme için e-posta). Render ücretsiz planı giden SMTP portlarını (25/465/587) engelliyor; SMTP göndericisi orada çalışamaz. Brevo ücretsiz planı günde 300 e-posta, alan adı gerektirmeden tek doğrulanmış gönderici adresiyle çalışır.
**Alternatives:** Ücretli Render planı + SMTP, Resend (başkalarına göndermek için alan adı doğrulaması ister), Mailgun/SendGrid (ücretsiz planları kısıtlı).
**Consequences:** Hesap (Brevo) ve API anahtarı kullanıcının işidir; anahtar yalnız Render ortam değişkeninde durur, repo'ya girmez. Gmail gibi ücretsiz bir adres gönderici yapılırsa DMARC yüzünden e-postalar spam'e düşebilir — kalıcı çözüm kendi alan adı. Gönderilemeyen e-postalar kuyrukta 5 deneme / 24 saat bekler.
**Tarih:** 2026-10-17

---

## DD-248 — Atölye Meslek Sayfaları, Yığın 999, Av Hayvanı Görselleri
**Karar:** (1) Atölye → Üretim yalnız meslek listesini gösterir (seviye, tarif sayısı, öğren); öğrenilmiş mesleğe dokunmak ayrı sayfa açar (`ProfessionScreen`): o mesleğin tarifleri, eşya adına dokununca aşağı açılan özellik paneli (tür/yuva, bonuslar, kullanım etkisi, seviye/sınıf şartı, ağırlık, yığın, beceri XP'si). Bunun için `GET /v1/recipes` satırlarına `output` (çıktı şablonunun özellikleri; `crafting.Handler.WithItems`). (2) Yığın sınırı 99 olan tüm şablonlar (toplanan 150 kaynak + 29 eşya) 999; 20'lik iksir yığınları değişmedi. (3) Toplama → Avcılık'ta kaynak satırı, görseli varsa hayvanın resmini gösterir: `client/assets/animals/<hayvan>.png`, hayvan kodu = kaynak kodu eksi `_pelt/_hide/_feather/_scale` (`huntAnimal`); dosya listesi `assets/animals/README.md`. Avdan çıkacak eşyalar (hayvan başı ganimet tablosu) henüz değişmedi — kullanıcıyla belirlenecek.
**Reason:** Kullanıcı istekleri.
**Alternatives:** Tarifleri açılır meslek kartlarında tutmak (liste yine uzun kalır), ayrı eşya şablonu ucu (tarif yanıtına eklemek tek istek).
**Consequences:** Ağırlık sınırı yumuşak olduğu için 999'luk yığınlar ağırlık göstergesini kolayca aşar (ret yok). Hayvan resimleri APK'ya gömülür (yeni resim = yeni APK).
**Tarih:** 2026-10-17

---

## DD-249 — AFK Süreleri 1 / 3 / 6 Saat
**Karar:** Toplama ve av için seçilebilen süreler yalnız 1 sa, 3 sa, 6 sa (varsayılan 1 sa; dev sürümünde ek 1/5 dk test seçenekleri durur). Sunucu tavanı 12 sa → 6 sa (`activity.MaxDurationSeconds`, `hunt.MaxHuntSeconds`); alt sınır değişmedi (erken toplama zaten süreyle orantılı ödüllendirir). Zindan süreleri şablon başına sabit kalır (15–30 dk).
**Reason:** Kullanıcı isteği.
**Alternatives:** Sunucuda yalnız bu üç değeri kabul etmek (testler ve dev kısa süreleri bozulur; kısa süre ekonomik avantaj sağlamaz).
**Consequences:** 6 saatten uzun kontratlar reddedilir (`E_VALIDATION`); mevcut uzun kontratlar normal biter.
**Tarih:** 2026-10-17

---

## DD-250 — Seviye 1–1000 + Sınıf İlerlemesi (Savaşçı referans uygulaması)
**Karar:** Veri odaklı ilerleme motoru `internal/leveling` (`content/progression.yaml` + `skills.yaml` açılımları; içerik yüklemesinde atomik kurulur, yokken yerleşik varsayılanlar). XP eğrisi `dakika = 3.75006670 + 0.00121846427·L^1.69337967`, `xp = ceil(dakika × 1000)` (formül ya da tablo; 1..999 tablosu önbellekte). Seviye tavanı 1000 (DB CHECK); tavanda XP `postcap_xp`'de birikir. Tek XP yolu `character.Service.GrantXP(char, miktar, XPSource{kind,id,metadata})`: doğrulama → `xp_multiplier` → satır kilidi → (kaynak id'si varsa `character.xp_grant` ile bir kez) → çoklu seviye → terfi kapısında durup kalanı `promotion_overflow_xp`'de tutma → seviye ödül hattı (`onLevelUp`: nitelik/yetenek puanı, unvan, terfi, uyanış/nihai/zirve/ustalık, yol açılımları) → `character.progression_event` (tür+ref benzersiz, bir kez) — hepsi tek tx. Hunt/zindan/görev bu yoldan öder; AFK savaş XP'si `afk_combat_xp_per_minute` (1000) × savaşılan dakika. Terfi kapıları Sv 100 (Muhafız/Cellat) ve 300 (Demir Kale/Savaş Lordu, Gözü Dönmüş/Silah Ustası); seçim (`POST /character/{id}/promotion`) bekleyen XP'yi aynı motordan işler. Nitelik puanı her 2 seviyede 1 + 100/300/600/850/1000'de +5 (toplam 525), sunucuda dağıtılır (`POST .../attributes`); savaş ölçeği azalan getirili (100/250/400 basamakları), gereksinimler ham değer. Yetenek puanı 100'den 980'e her 20 seviyede 1 (45); üç Savaşçı ağacı (Savunma/Silahlar/Kan) 25 puan + zirve (Sv 850 + ağaçta 25 puan). Seviye unvanları (Kıdemli…Ebedi) ve sınıf ustalığı (Sv 1000: unvan, rozet, çerçeve, küçük bonus). Savaşçı: Öfke kaynağı (boş başlar, isabet başına dolar), taban statlar (otomatik büyüme yok), yapılandırılabilir türetilmiş formül + yüzde tavanları, 49 yol yeteneği (`spec`, `auto_trigger`, `activation_mode`, `resource_gain`, `duration_turns`; etkiler strike/heal/guard/empower). Savaş motoruna Öfke, koruma/güçlenme, HP eşik tetikleri ve blok/kritik hasar/infaz/eksik can hasarı/can çalma/hasar azaltma (hepsi sıfırken eski kurallar birebir, ek RNG yok). Ekipman gereksinimleri (`requirements`: dal + ham nitelik), uzmanlık değişimi sınırı (meşgulken yok, aynı kademe, geçersiz ekipman çantaya). Uçlar: `GET /character/{id}/leveling` (+`/events`), `/talents`, `GET /v1/progression/config`, ops `GET /internal/ops/progression`, dev `POST /v1/dev/character/{id}/{xp|respec}`. İçerik adları en/tr tam, zh-CN/es kısmi (ilerleme anahtarları zorunlu, kalanı en'e düşer). İstemci: ilerleme kartı (MAKS, unvan, sıradaki hedef + süre, bekleyen XP), terfi seçimi, nitelik dağıtımı, yetenek ağacı ekranı, uygulama dilleri en/tr/es/zh.
**Reason:** Kullanıcının Savaşçı ilerleme belgesi (DD-250 kapsamı).
**Alternatives:** Seviye mantığını her kaynakta ayrı tutmak (çift ödül riski), sınıf dallarını koda gömmek, otomatik stat büyümesini korumak.
**Consequences:** Yaklaşımlar: blok karşı saldırısı, bariyer, tehdit/grup aurası, kontrol direnci, silah aileleri ve "kalkan gerekir" savaşta henüz mekanik değil — eşdeğer stat değiştiricileri ve koruma/güçlenme etkileriyle temsil edilir (veri hazır). Savaşçı formülüne referansta olmayan taban/seviye terimleri denge simülasyonuyla eklendi (Sv 1–25 içeriği oynanabilir kalsın). İçerik Sv 25'te bitiyor; görev XP'leri yeni eğriye göre küçük. Diğer sınıflar eski otomatik büyümeyle sürer (puanlar ek olarak verilir) — kendi ağaçları geldiğinde aynı motor. Ön koşul: istemci üretilmiş çeviri sınıfları kapsam ölçümünden hariç.
**Tarih:** 2026-10-17

---

## DD-251 — Karakter Slotu 5
**Karar:** Hesap başı karakter slotu 2 → 5 (`CHARACTER_SLOT_CAP` varsayılanı ve servis varsayılanı). DD-OPEN-8'i kapatır.
**Reason:** Kullanıcı isteği.
**Alternatives:** Ortam değişkeniyle yalnız canlıda artırmak (varsayılan ile canlı ayrışır).
**Consequences:** Canlıda `CHARACTER_SLOT_CAP` elle ayarlanmışsa o değer geçerli kalır. İstemci slot sayısını sunucudan okur.
**Tarih:** 2026-10-17

---

## DD-252 — On temel sınıf, 20 terfi, 40 uzmanlık (tek motor, çok veri)
**Karar:** DD-250 motoru genelleştirildi; sınıfa özel motor yok. Temel sınıflar: Savaşçı, Haydut (`rogue`), Korucu (`ranger`, eski `hunter`), Büyücü, Keşiş, Rahip (`cleric`, eski `priest`), Paladin, Druid, Ozan, Şaman — her biri 2 terfi (Sv 100) × 2 uzmanlık (Sv 300). Migration 0037: `hunter→ranger`, `priest→cleric` (karakter + eşya kısıtı), `class_skill.effect` yeni etkiler. Ortak kilometre taşları: temel Sv 1/10/25/75 + Sv 50 pasif; terfi pasifi + 100/150/200/250; uzmanlık 300 pasif, 350, 400, 450 pasif, 500, 550 pasif, 600 Uyanış (otomatik), 700 tek Nihai (uzun bekleme), 850 zirve (yetenek ağacında 25 puan), 1000 ustalık (unvan, rozet, çerçeve, başarım `mastery_<uzmanlık>`, prestij bayrağı, ≤%3 bonus). Sınıf başı 3 yetenek ağacı (5×5 + zirve), toplam 45 puan. **Motor:** türetilmiş stat formülü sınıf başı genel terimler (`base + per_level·L + Σ katsayı·etkin nitelik`, can/saldırı/zırh/kritik/isabet/kaçınma/blok/mana); kaynak modeli `{kind, max, start_pct, per_hit, per_turn}` (Öfke, Enerji, Odak, Mana, Chi, İnanç, Doğa Enerjisi, Ritim, Ruh); savaş etkileri strike/heal/guard/empower + **barrier, dot, hot, weaken, stun** (sersemletme zinciri yasak: bekleme ≥ süre+2); güçlendirme/zayıflatma hedef %40 canın altındayken kullanılmaz, öncelikte vuruşlardan sonra gelir. Yetenek veri modeli: `type` (DAMAGE, HEAL, BARRIER, BUFF, DEBUFF, DOT, HOT, STUN, ROOT, SILENCE, TAUNT, DISPEL, RESOURCE_GAIN, …; etkiyle uyumu doğrulanır), `damage_type`, `target_type`, `range`, `cast_time`, `tags`; her aktifin AFK kuralı zorunlu (`auto_trigger`: `on_cooldown` | `self_hp_below` | `target_hp_below` | `ally_hp_below` — tek başına oyuncu kendi müttefikidir). Pasif modeli `{trigger (ON_HIT … ON_COMBAT_END), condition, internal_cooldown, stack_limit, duration, scaling}`; etkisi modifier'larla sabit-durum değeri olarak uygulanır. Rol etiketleri ortak küme (tank, melee_dps, healer, …, afk_safe, group_support); uzmanlık başı AFK profili (1–5) + güçlü/zayıf içerik. Ekipman izinleri sınıf düğümünde (`equipment_families` + `armor_types` ile doğrulanır), en derin tanımlı düğüm geçerli; kuşanırken `E_EQUIPMENT_PERMISSION`. **Solo Adaptasyon:** destek yolları (`support_power`) partide değilken gücün %35'i saldırıya döner (`solo_adaptation.ratio`, app ModifierSource; parti = başka aktif üye). İçerik `scripts/gen_classes.py` + `scripts/class_data.py` ile üretilir (Savaşçı elle, aynı şema): `progression.yaml` (63 düğüm, 30 ağaç), `skills.yaml` (471 yetenek), yeni silah hatları (keşiş yumruk, paladin çift el topuz, druid orak, ozan lavta, şaman totem; 1. kademe dükkânda), i18n en/tr/zh-CN/es (sınıf/terfi/uzmanlık + rol, yetenek ad+açıklama, yetenek ağacı, unvan, kaynak, ustalık başarımı). İçerik paketi seviye/ustalık unvanlarını ve ustalık başarımlarını kaydeder. Denge: `internal/balance` 1–25 hedefleri 10 sınıfla; `TestSpecializationBalance` 40 uzmanlığı Sv 700/1000'de alan ortalamasına kalibre sentetik moblarla karşılaştırır (DPS en hızlı öldürür, tank en uzun dayanır, destek en çok grup değeri, hiçbiri tüm kategorilerde üstte değil); rapor `reports/class_balance_report.{json,md}` (`BALANCE_REPORT=1`).
**Reason:** Kullanıcının "ALL CLASSES PROGRESSION SYSTEM" belgesi.
**Alternatives:** Sınıf başı ayrı motor; AFK için ayrı yetenek tanımları; destekler için ayrı solo kopyalar.
**Consequences:** Yaklaşımlar: savaş tek hedefli — AOE, parti auraları, tehdit, evcil hayvan, tuzak, yük/yığın tetikleyicileri ve "dost" iyileştirmeleri mekanik olarak tek karşılaşmada temsil edilir (eşdeğer modifier/etki; AOE ve grup değeri raporda tasarım profilinden). zh-CN/es'de ilerleme anahtarları tam, diğer içerik en'e düşer. Savaşçı Öfke'si isabet başı 8→12, Savaşçı/diğer güçlendirmeler öncelikte vuruşlardan sonra; `berserker_death_defiance` düşük can pasifi oldu. Rahip topuzları INT yerine BİL verir. Yeniden uzmanlaşma sınırı DD-250'deki gibi (dev uç), ödemeli respec yok. Yüksek seviye denge sentetik moblarla — gerçek içerik Sv 25'te bitiyor.
**Tarih:** 2026-10-18

---

## DD-253 — Karakter oluşturmada sınıf bilgisi
**Karar:** Karakter oluşturma ekranında seçilen sınıf için bilgi kartı: ikon, ad, roller, tanıtım metni (`class.<sınıf>.intro`, 10 sınıf × en/tr/zh-CN/es, içerik i18n'inde; ilerleme anahtarı olarak zorunlu), kaynak (tür, büyüklük, başlangıç payı, isabet/tur kazancı), nitelik odağı, taban nitelik çubukları (toplam nitelik puanı sunucudan), ekipman izinleri (`equipment.<aile>` adları), Sv 1–75 yetenekleri ve **piramit sınıf yolu** (tepede temel sınıf, ortada Sv 100 terfileri, tabanda Sv 300 uzmanlıkları, bağlantı çizgileri). Bir düğüme dokununca küçük pencere o düğümün yeteneklerini (seviye, tür etiketi, etki, maliyet/bekleme) ve uzmanlıkta Sv 1000 ustalık unvanını gösterir. Veri yeni `GET /v1/classes` ucundan (70 düğüm: taban statlar, kaynak, etkin ekipman izni, AFK profili, ustalık unvanı, `attribute_points_total`, kilometre taşları; 5 dk önbellek) ve mevcut `GET /v1/classes/{class}/skills`'ten; adlar `ContentNames` (oluşturma ekranı yüklenmemişse kendisi yükler). İstemci oturum boyunca önbellekler. Sınıf ikonları gelene kadar Material ikonları + sınıf rengi.
**Reason:** Oyuncunun sınıfı seçerken nitelikleri, rolü ve terfi yollarını görmesi (kullanıcı isteği).
**Alternatives:** Sınıf verisini istemciye gömmek (içerikle ayrışır); tam `progression/config`'i indirmek (gereksiz XP tablosu ve açılımlar).
**Consequences:** Bilgi kartı tamamen içerikten gelir; yetenek/nitelik değişikliği istemci güncellemesi gerektirmez. Yükleme başarısızsa ekran bir uyarıyla oluşturmaya devam eder.
**Tarih:** 2026-10-18

---

## DD-254 — Seviye farkına göre XP + Sv 25–100 dünyası
**Karar:** (1) Savaş XP'si mobun karaktere göre seviyesine bağlı: en fazla 10 seviye aşağısı tam, 11–25 aşağısı %50, daha aşağısı %10; kendi seviyesi ve üstü her zaman tam (`progression.yaml` `mob_level_xp` bantları, `leveling.MobXPFactor/ScaleMobXP`). Av ve zindan (dalga + boss bonusu) ödülü sözleşmedeki oyuncu anlık görüntüsünün seviyesine göre ölçeklenir; seviyesi bilinmeyen eski görüntüler tam alır. `GET /v1/mobs` `xp_bands` döner; istemci av listesinde mob başına "XP %…" gösterir (hesap sunucuda). (2) Seviye farkı kuralına 5 seviyelik muafiyetten sonra seviye başı %6 saldırı kaybı eklendi (en az %10); zindan boss'ları (+2…+4) etkilenmez, çok üstteki moblar şifacı sınıflar için de ölümcül olur. (3) Sv 25–100 içeriği `scripts/gen_world.py` ile: 3 bölge (Sisli Bataklıklar 25, Demir Dağlar 50, Kuzey Kalesi Harabeleri 75 + bağlantılar), 18 mob (5 seviyede bir; mevcut eğri hp 150+25L, saldırı 28+2,1L, savunma 1+0,8L × arketip: normal/tank/glass/evasive/crit/elite/chief), 3 zindan (Batık Tapınak 38, Unutulmuş Maden 63, Kuzey Kalesi 88; boss girişin 4 seviye üstünde, `chief`), adlar en/tr/zh-CN/es. (4) Ekipman merdiveni `scripts/gen_gear_ladder.py` ile Sv 35/50/65/80/95 kademeleri (10 sınıf silahı + 5 zırh = 75 eşya, en/tr adları). Ekipman düşüşü artık kaynağın ulaştığı en yüksek kademeden (`item.GearPool`; kademeler seyrekleşince de düşer). Denge: `internal/balance` seviyeleri 100'e kadar; hedefler aynı, "dışarıdan seviye" kuralı Sv < 25'te +5, sonrasında +10 (beceri setleri ve iyileştirme +5'i aşılabilir kılıyor). `TestMobCurve` (MOBCURVE=1) eğrinin ekipmanlı ortalama oyuncuya uyduğunu gösterir.
**Reason:** XP mob seviyesinden bağımsızdı (Sv 500 karakter Sv 1 mobdan tam XP alıyordu); içerik Sv 25'te bitiyordu.
**Alternatives:** Mob başına sabit XP (AFK oranından kopar); seviye farkını oyuncu tarafında engellemek (av başlatmayı yasaklamak).
**Consequences:** Oyuncu her ~10 seviyede yeni moba geçer. Yeni zindanlar giriş seviyesinde eskilerden kolay (%74–100 temizleme). İçerik Sv 100'de bitiyor; sonraki dönemler aynı betiklerle eklenecek. Baskınlar ertelendi.
**Tarih:** 2026-10-18

---

## DD-255 — Sv 100–300 dünyası, ganimet seviye sınırı, ekipman adları 4 dil
**Karar:** (1) Kırık Taç Krallığı (yeni kıta, Sv 100): Taç Ovaları 100, Yanık Orman 150, Tuz Çölü 200, Taç Mezarlığı 250; 24 mob (10 seviyede bir + zindan boss'ları), 4 zindan (Asiler Kampı 128, Ateş Kökü 178, Gömülü Saray 228, Kraliyet Kriptası 278; boss `elite`, girişin 4 üstünde), Kuzey Kalesi Harabeleri ↔ Taç Ovaları bağlantısı, adlar en/tr/zh-CN/es (`scripts/gen_world.py`, ortak `scripts/i18n_util.py`). (2) Ekipman merdiveni 110–285 (25 seviyede bir 8 kademe, 120 eşya); 35–95 kademeleri aynen korunur; tüm merdiven eşyalarının adları artık 4 dilde (`scripts/gen_gear_ladder.py`). (3) **Ganimet seviye sınırı:** av/zindan ekipmanı ve derisi mobun seviyesinden, en çok karakterin seviyesinden düşer (`activity.DropLevel`). Ekipman havuzu ulaşılan en yüksek kademenin kendisidir (önceki "5 seviyelik pencere" seyrek kademelerde 1. kademe eşyaları da içeriyordu); şart taşıyan özel eşyalar rastgele düşmez. (4) Seviye farkı saldırı kaybı 5 seviyelik muafiyetten sonra seviye başı %10. (5) Denge simülasyonu Sv 1–300, Sv 100'den itibaren her sınıfın iki terfi yolu ayrı ölçülür (beceriler dahil). "Dışarıdan seviye mob tehlikeli" hedefi yalnız Sv < 25: sonrasında tank/şifacı yollar çok üstteki moblara dayanabilir, ama XP dakika başı sabit ve ganimet karakter seviyesiyle sınırlı olduğundan üst seviyede avlanmak kâr getirmez.
**Reason:** Kullanıcı isteği (adlar + 100–300 içeriği); ölçüm, güvenli yolların üst seviye mobdan daha iyi ganimet topladığını gösterdi.
**Alternatives:** Üst seviye avı yasaklamak; seviye farkı cezasını şifacıları da öldürecek kadar sertleştirmek (zindan boss'larını bozar).
**Consequences:** 100–300 zindanları giriş seviyesinde %78–100 temizlenir. Mevcut eşyaların zh/es adları (merdiven dışı) hâlâ İngilizceye düşer. İçerik Sv 300'de bitiyor.
**Tarih:** 2026-10-18

---

## DD-256 — Sv 300–600 dünyası: Ejdersırtı Adaları
**Karar:** Yeni kıta Ejdersırtı Adaları (Sv 300, Taç Mezarlığı ↔ Fırtına Kıyısı bağlantısı): Fırtına Kıyısı 300, Ejder Yamaçları 350, Sis Ormanı 400, Kristal Mağaralar 450, Volkanik Ada 500, Taht Zirveleri 550; 36 mob (10 seviyede bir + 6 zindan boss'u), 6 zindan (Batık Filo 328, Kuluçka Mağarası 378, Kayıp Koru 428, Kristal Kalp 478, Volkan Kalbi 528, Ejder Tahtı 578; 60 dk), adlar en/tr/zh-CN/es. Sv 300+ zindan boss'ları yeni `overlord` arketipi (can ×1,85, saldırı ×1,3, savunma ×1,25; zindan bunu ayrıca ×2 can ×1,2 saldırı yapar). Ekipman merdiveni 310–585 (12 kademe, 180 eşya, 4 dil). Denge simülasyonu Sv 1–590; Sv 300'den itibaren 40 uzmanlığın her biri ayrı ölçülür. Mob eğrisi değişmedi: Sv 600'de ortalama oyuncu eğrinin %5 can / %14 saldırı önünde, hedefler tutuyor.
**Reason:** Kullanıcı isteği (300–600 içeriği).
**Alternatives:** Mob eğrisini oyuncu ortalamasına yeniden oturtmak (300–600 avları zorlaşırdı; hedefler zaten tutuyor).
**Consequences:** 300–600 zindanları giriş seviyesinde ortalama %98, en zayıf uzmanlıkta %71–92 temizlenir (Peygamber gibi destek uzmanlıkları en zayıf). Uyanış (Sv 600) sonrası içerik yok. Entegrasyon testlerinde lonca etiketleri artık 5 karakter [A-Z0-9] (paylaşılan dev DB'de birikmiş loncalarla çakışıyordu).
**Tarih:** 2026-10-18

---

## DD-257 — Sv 600–700 dünyası: Yıldızsız Topraklar
**Karar:** Yeni kıta Yıldızsız Topraklar (Sv 600, Taht Zirveleri ↔ Karanlık Ova bağlantısı): Karanlık Ova 600, Kayıp Şehir 650; 12 mob (10 seviyede bir + 2 zindan boss'u, `overlord`), 2 zindan (Söndürülmüş Gözlemevi 628, Saat Kulesi 678; 60 dk), adlar en/tr/zh-CN/es. Ekipman merdiveni 610/635/660/685 (60 eşya, 4 dil). Denge simülasyonu Sv 1–690 (uyanış becerileri dahil). Yıldız Kemiren `crit` arketipi: `glass` hâliyle Arbalet Ustası Sv 630'da saatte %60 ölüyordu.
**Reason:** Kullanıcı isteği (600–700 içeriği).
**Alternatives:** Arbalet Ustası'nı güçlendirmek (diğer seviyelerde hedefler tutuyor; içerik tarafında tek mob yeterli).
**Consequences:** 600–700 zindanları giriş seviyesinde en zayıf uzmanlıkta (Arbalet Ustası) %87–89, en güçlüde %100. Nihai yetenek (Sv 700) sonrası içerik yok.
**Tarih:** 2026-10-18

---

## DD-258 — Sv 700–850 dünyası: Gök Yarığı
**Karar:** Yeni kıta Gök Yarığı (Sv 700, Kayıp Şehir ↔ Uçan Adalar bağlantısı): Uçan Adalar 700, Yarık Kenarı 750, Göksel Tapınak 800; 18 mob (10 seviyede bir + 3 `overlord` zindan boss'u), 3 zindan (Fırtına Kalesi 728, Yarık Çekirdeği 778, Göksel Kubbe 828; 60 dk), adlar en/tr/zh-CN/es. Ekipman merdiveni 710–835 (6 kademe, 90 eşya, 4 dil). Denge simülasyonu Sv 1–830 (nihai yetenekler dahil); ek ayar gerekmedi.
**Reason:** Kullanıcı isteği (700–850 içeriği).
**Alternatives:** —
**Consequences:** 700–850 zindanları giriş seviyesinde en zayıf uzmanlıkta %92–93, en güçlüde %100. Yetenek ağacı zirvesi (Sv 850) sonrası içerik yok.
**Tarih:** 2026-10-18

---

## DD-259 — Sv 850–1000 dünyası: Boşluğun Kıyısı (içerik 1000'e tamam)
**Karar:** Yeni kıta Boşluğun Kıyısı (Sv 850, Göksel Tapınak ↔ Kül Denizi bağlantısı): Kül Denizi 850, Son Kale Surları 900, Ebediyet Kıyısı 950; 18 mob (10 seviyede bir + 3 `overlord` zindan boss'u; son mob Yıkım Avatarı Sv 1000), 3 zindan (Batık Taht 878, Son Kale 928, Ebedi Taht 978; 60 dk), adlar en/tr/zh-CN/es. Ekipman merdiveni 860–985 (6 kademe, 90 eşya, 4 dil). Ruh Biçen `crit` arketipi (`glass` hâliyle Arbalet Ustası Sv 930'da %33 ölüyordu). Denge simülasyonu Sv 1–990. `-race` derlemesinde ağır denge testleri atlanır (tek goroutine, deterministik; kapsam işinde yarış dedektörsüz koşar), yoksa CI'da paket zaman aşımını geçerdi.
**Reason:** Kullanıcı isteği (850–1000 içeriği).
**Alternatives:** Arbalet Ustası'nı güçlendirmek (iki seviyede aynı desen: yüksek saldırılı `glass` mob; ileride sınıf ayarı düşünülebilir).
**Consequences:** Seviye 1–1000'in her 10 seviyesinde mob, her ~50 seviyesinde zindan, her 25 seviyede ekipman kademesi var (140 mob, 25 zindan, 54 bölge). 850–1000 zindanları giriş seviyesinde en zayıf uzmanlıkta %93. Baskınlar (grup içeriği) hâlâ ertelendi. Arbalet Ustası yüksek saldırılı moblarda en kırılgan uzmanlık.
**Tarih:** 2026-10-18

---

## DD-260 — Arbalet Ustası güçlendirmesi
**Karar:** Arbalet Ustası: Ağır Arbalet Okları +%5 → +%8 saldırı ve +%3 hasar azaltma; Ezici Ok %280 → %300; Hızlı Doldurma +%20 → +%25 saldırı ve kullanımda +20 Odak; Güçlendirilmiş Mekanizma +%5 can eklendi (hasar azaltma 4 → 5); Balistik Uzmanlığı +%12 → +%15 kritik hasar, +%3 → +%4 saldırı, +%3 can çalma. DD-257/259'da onu korumak için `crit` yapılan Yıldız Kemiren (630) ve Ruh Biçen (930) yeniden `glass`.
**Reason:** Uzmanlık raporunda en düşük AFK güvenliği (%43) ve DPS'ler içinde yavaş öldürme (5,3 tur); yüksek saldırılı moblarda iki seviyede av hedefini tutturamıyordu.
**Alternatives:** Mobları yumuşak bırakmak (sorunu içerikte saklar).
**Consequences:** Sv 1000 karşılaştırmasında 4,5 tur öldürme, 31 tur dayanma, %100 boss ve AFK güvenliği; Keskin Nişancı daha hızlı (2,9 tur) ama daha kırılgan kalır. Yıldız Kemiren ve Ruh Biçen'e karşı ölüm %0. `reports/class_balance_report.{json,md}` yenilendi.
**Tarih:** 2026-10-18

---

## DD-261 — Baskınlar v1: 3–5 kişilik parti zindanları
**Karar:** Baskın, `group_max > 1` olan bir zindan şablonudur (`content/dungeons.yaml`, `group_min`/`group_max`); haftalık kilidi zindanlarla paylaşır (`dungeon.weekly_lockout`, `contract_id` = üye satırı). Altı baskın, her sınıf kilometre taşında bir: Kuzey Kalesi Kuşatması (100), Taçsız Kral'ın Sarayı (300), Ejder Kral'ın İni (600), Zaman Lordu'nun Kulesi (700), Başmeleklerin Düşüşü (850), Düşmüş Tanrı'nın Tahtı (1000); 60 dk, 3–5 kişi. Akış: parti lideri lobi açar → parti üyeleri katılır (her biri için seviye, haftalık kilit, can, çanta kontrolü ve savaş anlık görüntüsünün dondurulması) → lider en az `group_min` kişiyle başlatır (partiden ayrılanlar düşer; sonra kimse giremez/çıkamaz) → süre dolunca her üye kendi ödülünü toplar. Savaş: yeni `combat.Group` — her tur yaşayan her oyuncu sırayla hareket eder; iyileştirme, iyileştirme-zamanla ve kalkan en yaralı müttefike gider (tek kişilik savaşta kendisi; 1 kişilik grup tek kişilik savaşla birebir aynı); mob en yüksek savunmalı oyuncuya (tank) vurur, boss her 3. vuruşta tüm gruba %60 güçte alan saldırısı yapar ve sersemletilemez; düşen oyuncu baskın boyunca düşük kalır. Dalga mobları grup için ×3 can; boss'lar yeni `titan` (Sv ≤300: can ×30, saldırı ×3,2) / `elder_titan` (600+: saldırı ×4,8) arketipi, `raid_boss: true` — av listesinde yok, av başlatmayı reddeder, yalnız baskında boss olabilir (contentpack doğrular). Savaş seed + dondurulmuş anlık görüntülerden deterministik: her üye aynı kaydı görür. Ödül kişisel: her öldürme AFK savaş XP oranı (üyenin seviyesine göre ölçekli) + mobun altını (bölünmez), boss bonusu (zindanın 3 katı), garanti nadir eşya ve 2 garanti ekipman (boss kademesi, `DropLevel` ile karakter seviyesine sınırlı), dalga başına %4 ekipman; kendi zar atışları. Toplama idempotent (ledger `raid-reveal-<üye>`, eşya kaynağı üye satırı, saklı sonuç). Migration 0038 (`raid.run`, `raid.member`, parti başı tek açık lobi, karakter başı tek etkin baskın, `group_min/max`, `raid_boss`). Uçlar: `GET /v1/raids`, `GET/POST /v1/character/{id}/raid`, `POST /v1/character/{id}/raid/{raid_id}/{join|leave|start|reveal}`; lobi değişince partiye WS `raid.changed`. İstemci: Macera → Zindan içinde "Zindanlar | Baskınlar" seçimi.
**Reason:** Dünya senaryosundaki kilometre taşı baskınları grup savaşı istiyordu; motor tek kişilikti. Ayrı bir şablon/kilit sistemi yerine zindanı genişletmek içerik, doğrulama, ad ve kilit altyapısını yeniden kullanır. Hazır-onay yerine katılma adımı: kimsenin zamanı onu katılmadan bağlanmaz.
**Alternatives:** 10–25 kişilik lonca baskını (docs/guild/raid.md, Phase 02 — parti tavanı 5, DD-128); lider dağıtımlı ganimet (master loot) — v1 kişisel ganimet; konumsal mekanikler (Phase 03).
**Consequences:** `TestRaidTargets`: her baskın kendi seviyesinde 10 karışık 5'li grupla ortalama ≥ %70, en kötü grup ≥ %40 temizler; 3 kişi ortalama ≤ %30 (ölçülen ortalamalar: 5'li %94–100, 3'lü %0–26). Baskın boss'unu tek başına avlamak mümkün değil. Bilinen sınırlama v1: baskın sürerken üye başka av/zindan başlatabilir (zindanlar da birbirini engellemiyor) — DD-263 ile kapandı.
**Tarih:** 2026-10-18

---

## DD-262 — Solo zindanlar sınırsız; ödül seviyeye göre
**Karar:** Solo zindanların haftalık kilidi kalktı: başlatma kilide bakmaz, toplama kilit satırı yazmaz (`dungeon.weekly_lockout` yalnız baskınlarda, DD-261). Ödül avdaki gibi karakter seviyesine göre: öldürme başına AFK savaş XP'si ve boss bonus XP'si `mob_level_xp` bantlarıyla (DD-254), boss bonus altını da aynı çarpanla ölçeklenir (yeni); eşya ve ekipman düşüşü karakter seviyesiyle sınırlı (`DropLevel`, DD-255). Seviye kapısı aynen kalır: `min_level` altındaki karakter zindana giremez, baskında her üye için ayrı kontrol edilir (`E_LEVEL_REQUIREMENT`).
**Reason:** Kullanıcı isteği: zindan tekrar oynanabilir olsun, ama aşağı seviye içerikte kazanç avdaki gibi düşsün.
**Alternatives:** Günlük giriş hakkı (yine sınır koyar); bonus altını ölçeklememek (yüksek seviyeli karakter düşük zindanı altın için döndürürdü).
**Consequences:** Zindan süresi (30–60 dk) ve tek işlem kuralı (DD-263) tekrar hızını sınırlar. `TestHTTPDungeonsUnlimitedAndLevelGated`: aynı hafta iki temizleme, seviye kapısı, baskın üyesi seviye kapısı.
**Tarih:** 2026-10-18

---

## DD-263 — Tek işlem kuralı
**Karar:** Karakter aynı anda yalnız bir iş yapar: toplama, av, zindan, baskın (lobiye katılmak dahil, ödül toplanana kadar), yolculuk veya üretim. Kural veritabanında: `character.active_action` görünümü her iş tablosunun etkin satırlarının birleşimi; her iş tablosuna `BEFORE INSERT` tetikleyicisi `character.claim_action(kind)` — `pg_advisory_xact_lock(karakter)` ile eşzamanlı başlatmaları sıraya koyar, başka türde etkin iş varsa SQLSTATE `OT001` (HINT = meşgul tür) fırlatır. Serbest bırakma yok: iş kendi satırı etkin durumdan çıkınca sayılmaz. Yeni `internal/action` (`ErrBusy`, `FromDB`, `Write`) → 409 `E_CHARACTER_BUSY`, `field` = meşgul tür; istemci "Önce karakterinin yaptığı işi bitir: av." der. Dünya etkinliğine kaydolmak iş sayılmaz (yalnız güç anlık görüntüsü).
**Reason:** Kullanıcı isteği; ayrıca paralel AFK işlerle çift ilerleme açığını kapatır.
**Alternatives:** Ortak "iş yuvası" tablosu (her bitiş yolunda serbest bırakma gerekir — unutulan bir yol karakteri kilitler); servislerde ön kontrol (yarışta iki iş birden girer).
**Consequences:** Migration 0039. `TestHTTPOneActionAtATime`: her tür diğerini engeller; dört farklı iş aynı anda başlatılınca tam biri kazanır. Bitmiş ama toplanmamış iş (av, zindan, baskın) de engeller: önce topla.
**Tarih:** 2026-10-18

---

## DD-264 — Süreli üretim
**Karar:** Üretim anlık değil, iş: `content/recipes.yaml` `craft_seconds` (birim başı; ergitme 30 sn + beceri×1,5 sn, iksir 30 sn + beceri×3 sn, ekipman (120 sn + beceri×30 sn) × parça çarpanı — göğüs/silah 1,5, bacak 1,2; 30 sn – 10,5 dk). `POST /craft {recipe_code, count}` işin tamamını önden kontrol eder (beceri, malzeme, altın, çanta, toplam ≤ 6 sa), `crafting.job` satırını açar (tek işlem, DD-263; tarif donmuş JSON), altını ve malzemeyi o an alır; ödeyemeyen başlangıç kendini geri alır. `POST /craft/{job_id}/collect` her an: biten birimler basılır ve meslek XP'si verir, bitmeyenlerin altını (`craft-refund-<iş>`) ve malzemesi iade edilir; iş bir kez `collected` olur, sonuç saklanır (tekrar = aynı yanıt). `GET /craft/active`. İstemci: tarifte "Tanesi 30 sn", adet seçicide toplam süre, süren iş kartı (canlı "x/y hazır"), erken toplamada iade onayı; üretim sürerken üret düğmeleri pasif.
**Reason:** Kullanıcı isteği; üretim de diğer işler gibi zaman ister.
**Alternatives:** Malzemeyi toplamada almak (iş sırasında satılırsa yarım kalır); iptalde iade yok (100 adetlik iş yanlışlıkla başlarsa kayıp).
**Consequences:** Meslek XP'si birim başı aynı. Bilinen sınırlama v1: üretim bitince bildirim gönderilmez (istemci ilerlemeyi gösterir).
**Tarih:** 2026-10-18

---

## DD-265 — Etkinlik açıklaması
**Karar:** Macera → Etkinlik'te karta dokununca alt pencere: nasıl katılınır (kayıt penceresi + seviye, istila bölgesi, güç formülü ve şu anki güç, puan = güç × dakika ve AFK tavan kademesi, topluluk hedefi → kademe sınırı, ödül kapanışta otomatik, kaydolmak diğer işleri engellemez), kademe ödülleri (puan eşiği, altın, jeton) ve şimdi katılırsa tahmini puan/kademe. Sezonluk etkinlik: günlük jeton ve bakiye. Sunucu etkinlik listesine `afk_cap_tier` ve `power_now` ekler.
**Reason:** Kullanıcı isteği: oyuncu ne yapması gerektiğini ve ne kazanacağını görsün.
**Alternatives:** Statik yardım metni (gerçek ödülleri/tavanı göstermez).
**Consequences:** Tahmin istemcide sunucunun `Contribution` kuralının kopyası, yalnız gösterim.
**Tarih:** 2026-10-18

---

## DD-266 — Karakter oluşturma ekranı şablonu
**Karar:** Karakter oluşturma alt pencere yerine tam sayfa (`CreateCharacterScreen`), onaylanan sade koyu şablonla: üstte 10 sınıflık yuvarlak ikon şeridi; portre (`assets/classes/<sınıf>.png`, yoksa sınıf ikonu) yanında ad, slogan (`class.<sınıf>.tagline`), rol etiketleri ve tanıtım; ikisinin altında sınıf sözü (`class.<sınıf>.quote`); Sınıf özellikleri (kaynak + dolma şekli, ana/ikinci nitelik + taban nitelik çubukları, ekipman, zorluk); İlk yetenekler ikon olarak — dokununca oku yeteneğe dönük baloncuk (tür, seviye, maliyet, açıklama), tekrar dokununca kapanır; Sınıf yolu piramidi (düğüme dokun → yetenek penceresi, DD-253); Irkını seç (4 kart, `assets/races/<ırk>.png`, yoksa ikon); alttaki "X ile devam et" ad penceresini açar ve karakteri oluşturur. Slogan ve söz 10 sınıf × 4 dil (`scripts/class_data.py` `TAGLINES`/`QUOTES`, contentpack zorunlu anahtar). Zorluk istemci sabiti (`classDifficulty`: Savaşçı/Paladin/Rahip kolay, Haydut/Büyücü zor, diğerleri orta).
**Reason:** Kullanıcının seçtiği taslak (v3).
**Alternatives:** Siyah-altın ve açık turuncu poster stilleri (iyi tam boy görseller ister).
**Consequences:** Irk kartlarında stat bonusu yazmaz: ırklar şu an yalnız görünüş ("Irk görünüşünü belirler; nitelikleri değiştirmez."). Görseller gelene kadar ikonlar. Eski `ClassInfoCard` silindi.
**Tarih:** 2026-10-19

---

## DD-267 — Altı yeni ırk
**Karar:** Trol, Hortlak, Karanlık Elf, Drakari, Tauren ve Gnom eklendi (kodlar `troll`, `undead`, `dark_elf`, `drakari`, `tauren`, `gnome`; sunucu izin listesi `character.knownRaces`). İlk dördü gibi yalnız görünüş: nitelik etkisi yok. Adlar ve tek satırlık tanımlar istemcide 4 dilde; oluşturma ekranında satırda 4 kart (4 + 4 + 2). Portreler `client/assets/races/<kod>.png`.
**Reason:** Kullanıcı isteği.
**Alternatives:** Irklara nitelik bonusu (denge simülasyonlarını ve sınıf eğrisini değiştirir; ayrı karar).
**Consequences:** Migration yok (`race_code` serbest metin, doğrulama serviste).
**Tarih:** 2026-10-19

---

## DD-268 — Aşçılık tarifleri
**Karar:** Aşçılık 18 yemekle dolduruldu (`python scripts/gen_cooking.py` → `items.yaml` + `recipes.yaml` üretilmiş blok, 4 dilde adlar). Her tarif: o bantta 2 balık + 1 ot + 1 odun (ateş) → 2 yemek; balıkçılık seviyesine paralel (balık Sv 1–100, beceri ≈ 0,9 × balık seviyesi). Yemek `consumable/food`, `use_effect: heal_hp` ≈ aynı seviyedeki mob canının %35'i (Izgara Lepistes 60 → Leviathan Kraliyet Şöleni 930), `min_level` = balık seviyesi, yığın 50. Altın ve meslek XP'si ergitme gibi kademe başı; süre 20 sn + 2 sn × beceri. Aşçılık meslek tavanı 60 → 100.
**Reason:** Kullanıcı isteği; meslek vardı ama hiç tarifi yoktu. Avcılık et değil post verdiği için malzeme balıkçılık + bitki toplama + odunculuk.
**Alternatives:** Geçici güç bufları (henüz tüketilebilir yalnız `heal_hp` destekliyor; ayrı iş).
**Consequences:** Yemekler iksirlerden zayıf ama toplama malzemesinden ucuza çıkar. İstemcide yemek ikonu. Hayvan görselleri 256×256'ya küçültüldü (75 MB → 3,5 MB; orijinaller git dışında `client/art_originals/`), iki dosya adı düzeltildi (`ancient_phoenix`, `black_panther`).
**Tarih:** 2026-10-19

---

## DD-269 — Yemek bonusu (tokluk)
**Karar:** Yemeğin `use_effect`'i can yenilemeye ek olarak `buff: {attack_pct, defense_pct, minutes}` taşır; 18 yemek 60 dakika boyunca seviyeye göre +%2..9 saldırı ve savunma verir (`gen_cooking.py`). Karakter başı tek yemek bonusu: migration 0040 `character.buff` (yeni yemek eskisinin yerini alır). Canlı bonus (`expires_at` gelecekte) `consumable.Service.Modifiers` ile sheet'in modifier yığınına `mult_add` olarak girer; av/zindan/baskın başlangıçta sheet'i dondurduğu için tok başlanan iş bonusu sonuna kadar taşır. Can doluyken bonuslu yemek yenebilir (iksir hâlâ `E_HP_FULL`). `GET /v1/character/{id}/buff`; kullanma yanıtında `buff`. İstemci: eşya detayında bonus satırı, yeme bildirimi, karakter sayfasında "Tok · +%x saldırı · +%y savunma · süre".
**Reason:** Kullanıcı isteği.
**Alternatives:** Bonusu sözleşme anlık görüntüsüne ayrıca yazmak (sheet zaten donuyor; gereksiz).
**Consequences:** Contentpack bonusu doğrular (0..50 %, en az biri > 0, 1..1440 dk). `TestHTTPFoodBuff`.
**Tarih:** 2026-10-19

---

## DD-270 — MeNa World adı ve logosu
**Karar:** Oyunun adı her yerde "MeNa World" (uygulama etiketi, 4 dilde başlık, web, e-posta konuları ve gönderen adı varsayılanı, yasal sayfalar, admin TOTP, dev console). Logo `client/assets/branding/logo.png` (şeffaf, 720 px): giriş ekranında başlık yerine, karakterler ekranının üst barında küçük. Android: adaptif ikon (koyu `#0F1013` zemin + logo), eski tip dairesel ikon, açılış ekranı (koyu zemin + logo; Android 12+ sistem splash'ı aynı renk), uygulama penceresi zemini koyu (beyaz flaş yok). iOS ve web ikonları aynı logodan (iOS opak). Kod adı Otter ve `com.otter.otter` değişmedi.
**Reason:** Kullanıcının logosu ve onayı (A: koyu zemin, üst barda logo, ad MeNa World).
**Alternatives:** Lacivert ikon zemini; yazısız sade ikon (çizim gerekir).
**Consequences:** 48 px ikonda logo yazısı okunması zor. Irk portreleri 300×400'e küçültüldü (27 MB → 2,6 MB; orijinaller `client/art_originals/races/`), ırk kartları 3:4 portre.
**Tarih:** 2026-10-19

---

## DD-271 — Zırh türleri: kumaş, deri, plaka
**Karar:**
- Zırh türleri `cloth`, `leather` ve `plate` (`progression.yaml` `armor_types`). Önceki `light`, `medium` ve `heavy` kaldırıldı.
- Zırhın alt sınıfı türüdür (`<tür>_armor`); slot `equip_slot`'ta kalır.
- Mevcut merdiven plaka sayılır; Sv 1 deri seti deri sayılır. Her kademeye eksik iki tür eklendi (`scripts/gen_gear_types.py`). Statlar, değer ve kod tabanı plaka karşılığıyla aynı; ağırlık kumaşta ×0,4, deride ×0,7.
- Her sınıf tek tür giyer:
  - kumaş: büyücü, rahip, druid
  - deri: düzenbaz, korucu, keşiş, ozan
  - plaka: savaşçı, paladin, şaman
- `leveling.IsFamily` artık zırh türlerini de yönetir. Önceden zırh hiç denetlenmiyordu.
- `item.GearPool` yalnız sınıfın kullanabileceği türü düşürür; denge simülasyonunun kitleri de aynı izni uygular.
- İçerik doğrulaması: silah ve zırhın alt sınıfı yönetilen bir aile ya da tür olmalı.
- Dükkân Sv 1 kumaş ve plaka setini de satar.

**Reason:** Kullanıcı kararı (K2, içerik denetimi). Statlar eşit tutuldu, böylece denge testleri değişmedi; türün işlevi şimdilik kimin giyeceğini belirlemek. Tür farkı setlerle gelebilir.
**Alternatives:** Türe göre farklı savunma/can (denge simülasyonunu yeniden ayarlamak gerekir). Zırhı sınıfsız bırakmak (önceki öneri).
**Consequences:** Zırh şablonu 235'ten 705'e çıktı. Önceden giyilmiş, artık izinli olmayan zırh çıkarılmaz; yalnız yeni kuşanma denetlenir. Testler: `TestResourcesAndEquipmentPermissions`, `TestShippedGearPoolHonoursPermissions`.
**Tarih:** 2026-10-19

---

## DD-272 — Sınıf başı iki silah ailesi
**Karar:**
- Her sınıfa ikinci bir silah ailesi eklendi; izin listesi tam olarak bir zırh türü ve bu iki aileden oluşur:
  - savaşçı: kılıç, balta
  - paladin: çift el topuz, kılıç
  - şaman: totem, ayin baltası
  - düzenbaz: hançer, kısa kılıç
  - korucu: yay, arbalet
  - keşiş: muşta, tonfa
  - ozan: lavta, meç
  - büyücü: asa, değnek
  - rahip: topuz, kutsal kitap
  - druid: orak, doğa odağı
- Kalkan yalnız Koruyucu kolunda (Koruyucu, Demir Kale, Savaş Lordu).
- Kullanılmayan 12 aile (`greatsword`, `spear`, `orb` …) ailelerden ve izinlerden çıktı.
- Yeni hatlar kademe başı sınıfın ilk silahının kopyası (saldırı ve statlar aynı), sınıfa kilitli; toplam 470 şablon. Sv 1 parçaları dükkânda.

**Reason:** Kullanıcı kararı (K3, içerik denetimi: 21/32 aile boştu).
**Alternatives:** Tüm aileleri doldurmak (≈1000 şablon daha); izinleri tek aileye indirmek.
**Consequences:** Terfisiz savaşçı Titan Siper Kalkanı'nı artık izin hatasıyla (`E_EQUIPMENT_PERMISSION`) alamaz.
**Tarih:** 2026-10-19

---

## DD-273 — Bölgesel av
**Karar:**
- Her av canavarının bir bölgesi var (`mob_template.region_code`, migration 0041). Baskın boss'larının bölgesi yok. Contentpack ikisini de zorunlu tutar ve bölgeyi doğrular.
- Bir canavar, karakter onun bölgesinde ya da o bölgeye yol bağlantısı (`region_link`) olan bir bölgede dururken avlanabilir. Ebeveyn/çocuk kenarları sayılmaz; yoksa kıtada durmak bütün bölgelerine erişim verirdi.
- Aksi halde `409 E_MOB_REGION`.
- Sv 1–25 canavarları:
  - Yeşilorman: Sv 1–6
  - Fısıltı Kıyısı: Sv 8–15
  - Kor Tepe: Sv 18–25
- Üretilmiş canavarlar `gen_world.py` bölgesini alır.
- Yeni yollar: Demirkale ↔ Yeşilorman (yeni oyuncu yolculuk etmeden avlanır), Fısıltı Kıyısı ↔ Kor Tepe.
- API:
  - `GET /location` → `hunt_regions`
  - `GET /v1/mobs` → `region_code`
- İstemci:
  - Av listesi yalnız erişilebilir canavarları gösterir; "Buradan avlanılan bölgeler" satırı var.
  - "Diğer bölgelerdeki canavarlar" bölge bölge listelenir.
  - Varışta liste yenilenir.

**Reason:** Kullanıcı kararı (K4): bölgeler ve yolculuk anlam kazansın.
**Alternatives:** Yalnız bulunulan bölge (fazla sürtünme); kapı yok (önceki durum).
**Consequences:** Zindan ve baskınlar bölgeye bağlı değil. Dünya etkinlikleri kendi bölge kuralını (`E_WRONG_REGION`) korur. Testler: `TestHTTPHuntRegionGate`, `TestStart_RegionGate`, `TestCatalog_HuntRegions`.
**Tarih:** 2026-10-19

---

## DD-274 — Meslek ölçeği ve tüm merdiveni kapsayan üretim
**Karar:**
- Meslek seviyesi 1–100, karakter seviyesi 1–1000'e karşılık gelir (K1).
- Sv L ekipmanın tarifi beceri ⌈L/10⌉ ister. Malzemeler o beceri + 1'e kadar toplanabilen en iyi kaynaklardır.
- Rastgele düşen her ekipmanın tarifi var (1645 tarif, `scripts/gen_crafting.py`):
  - Demircilik: plaka zırh ve metal silahlar. O kademenin külçesi bir alt külçeyle alaşım yapılır; ayrıca zırha deri kayış, silaha kütük sap.
  - Deri İşçiliği (yeni): deri zırh. Deri/post + ot.
  - Terzilik (yeni): kumaş zırh ve kutsal kitap. Ot lifi + deri astar.
  - Ağaç İşçiliği (yeni): asa, değnek, yay, arbalet, lavta, totem, doğa odağı, tonfa. Kütük + külçe; yay ve arbalete tüy de girer.
- Demircilik ve simya tavanı 100.
- Ücret: eşya değerinin %10'u. XP: beceri eğrisinin %25'i × parça ağırlığı. Süre: (90 + 6 × beceri) × parça ağırlığı.
- Yeni `heal_pct` etkisi azami canın yüzdesini yeniler (yukarı yuvarlanır).
- İksirler:
  - Küçük Can İksiri: Sv 1, %20
  - Orta Boy: Sv 100, %30
  - Büyük Boy: Sv 250, %40
  - Ekstra Büyük: Sv 450, %55
  - Çok Büyük: Sv 650, %70
  - Yüce (yeni): Sv 850, %85
- Yemekler balık seviyesi × 10'da açılır ve %15–40 can yeniler. Kalan 12 balığa ızgara yemeği eklendi (30 yemek).
- Deri ganimetinin avcılık seviyesi ⌈canavar Sv / 10⌉ (önceden ×4, Sv 25'te tavana vuruyordu).
- İstemci:
  - Tarifler meslek başına çekilir (`?profession=`).
  - Meslek sayfası açık tarifleri en yeniden eskiye, ardından bir sonraki kilitli adımı gösterir.
  - İyileştirme satırı "Azami canın %X kadarını yeniler".

**Reason:** Kullanıcı kararı (K1) ve içerik denetimi. Önceki durum:
- Üretim Sv 10'da bitiyordu.
- 24 külçenin ve 66 toplama kaynağının hiçbiri tarifte kullanılmıyordu.
- Deri ganimeti Sv 25'te tavandaydı.
- İksirler sabit can verdiği için üst seviyede anlamsızdı.

**Alternatives:** Meslekleri 1–1000'e genişletmek (kaynak listesi ×10). Yalnız demircilik.
**Consequences:**
- Elle yazılmış simya ve demircilik tarifleri kaldırıldı; aynı kodlar (`recipe_<çıktı>`) üretilmiş tariflerle güncellendi. Örneğin `recipe_helm_leather` artık Deri İşçiliği'nde.
- Ores Sv 12'de başladığı için ilk metal tarifleri madencilik 12 ister.
- Kullanılmayan kaynak 66'dan 7'ye indi (4 taş kuyumculuğu bekliyor).
- Testler: `TestShippedCraftingCoversGear`, `TestHealAmount`, `TestCombatHides_ByMobLevel`.

**Tarih:** 2026-10-19

---

## DD-275 — Nadirlik ve ekler güç verir
**Karar:**
- Nadirlik, eşyanın kendi düz (`flat`) bonuslarını çarpar ve sonuç yuvarlanır: seçkin ×1,05, nadir ×1,10, destansı ×1,20, efsanevi ×1,30, mitik ×1,50 (`item.StatMultiplier`). Yüzde türündeki bonuslar değişmez.
- Her ek tek bir stat verir: `taban + seviye başı × eşya Sv`, yuvarlanır, en az 1 (`item/affix.go` `affixStats`).
  - Silah önekleri: keskin/alevli saldırı, acımasız kritik hasar, çevik isabet, donmuş kritik, lanetli can çalma.
  - Zırh önekleri: güçlendirilmiş savunma, yaldızlı can, gölgeli kaçınma, rünik mana, muskalı hasar azaltma.
  - Kadim (ikisinde de): dayanıklılık.
  - Son ekler: ayı güç, kurt çeviklik, bilgelik zekâ, öküz dayanıklılık, isabet, yıkım, hız, kaplumbağa, kalkan blok, dayanıklılık can.
  - Yüzde statları yavaş büyür; Sv 1000'de en fazla birkaç puan, tavanlar korunur.
- Hesap `inventory.EffectiveModifiers`'ta yapılır. Kuşanılı eşyalar sheet'e bu değerlerle girer (`skill.EquippedBonus.Modifiers`).
- Çanta ve ekipman yanıtında `effective_bonuses` alanı var; istemci bonusları buradan gösterir.

**Reason:** İçerik denetimi: ekler yalnız isimdi, nadirlik yalnız satış değerini etkiliyordu; ganimet döngüsü güç kazandırmıyordu.
**Alternatives:** Ek kataloğunu içerik dosyasına taşımak (`affix.md`); 22 ek için Go tablosu daha basit. Nadirliğe göre ek sayısını artırmak.
**Consequences:** Denge simülasyonu sıradan ekipmanla ölçer; düşen nadir eşya oyuncuyu hedeflerin biraz üstüne taşır. Var olan eşyalar ek kodlarını zaten taşıdığı için anında güç kazanır. Testler: `TestAffixBonuses`, `TestEffectiveModifiers`, `TestHTTPRarityAndAffixesAddPower`.
**Tarih:** 2026-10-19

---

## DD-276 — Takılar, Koruyucu kalkanları, Kuyumculuk
**Karar:**
- Yeni eşya sınıfı `jewelry`. Sınıfsızdır; her sınıf takar. Sv 25 kademesinden itibaren her kademede üç takı var (`scripts/gen_jewelry.py`, 126 şablon):
  - Kolye: can
  - Yüzük: saldırı + isabet
  - Tılsım: savunma + kaçınma
- Yüzük ve tılsım şablonunun slotu `ring`/`trinket`. Çiftin boş olanına takılır; ikisi de doluysa ilki değişir (`inventory.PairedSlots`, `ValidTemplateSlot`).
- Sv 25 altında takı yok: denge hedefi orada aşılmış canavarların tehlikeli kalmasını istiyor.
- Koruyucu kolu için kalkan merdiveni: Sv 110+, 36 şablon, savunma + blok, `requirements: {spec: guardian}`. Kalkanlar düşmez, yalnız üretilir.
- Yeni meslek Kuyumculuk (tavan 100): takılar için külçe + alt külçe + taş. Kalkanları Demircilik yapar.
- Ganimet havuzu takıları da içerir. Denge kitleri yüzük ve tılsımı iki kez sayar.
- Baskın boss'u `elder_titan` saldırı çarpanı 4,8 → 5,3. Takılarla 3 kişilik baskın temizleme Sv 1000'de %51'e çıkmıştı; hedef ≤%30. 5 kişilik grup hâlâ ≥%88 temizler.
- İstemci: takı ikonu, pazar sınıf filtresinde Takı, `ring`/`trinket` slot adları.

**Reason:** Yol haritası Faz C (içerik denetimi: 12 slottan 5'i boştu, taşların kullanımı yoktu). Yan el için K3 kararına uyuldu: yalnız kalkan izni olan Koruyucu kolu.
**Alternatives:** Sınıfa göre takı (ana statlı). Bütün sınıflara yan el (K3 ile çelişir). Takıyı Sv 1'den başlatmak (erken oyun dengesi bozulur).
**Consequences:** Ekipman 2024 şablon. Taş kullanımı 4'ten 1'e indi. Testler: `TestHTTPJewelryPairs`, `TestShippedGearPoolHonoursPermissions` (takı), `TestRaidTargets`.
**Tarih:** 2026-10-19

---

## DD-277 — Baskın setleri
**Karar:**
- Yeni `internal/itemset` paketi bellekte bir katalog tutar; contentpack her yüklemede kurar (`itemset.Install`, progression gibi).
- `content/sets.yaml` iki bölümden oluşur:
  - `items`: set parçaları.
  - `sets`: kod, baskın kaynağı, parçalar ve parça sayısındaki bonuslar.
- Doğrulama: parça tek sette olur, eşikler artan sırada; kaynak baskındır; parçalar kuşanılabilir; modifier'lar bilinen türde.
- İçerik (`scripts/gen_sets.py`): 6 baskın × kumaş/deri/plaka = 18 set, 108 parça.
  - Her set 5 zırh ve 1 kolyeden oluşur.
  - Parçalar baskın seviyesindedir ve altındaki kademenin statlarından %15 güçlüdür.
  - Bonuslar: 2 parça +%4 can, 4 parça +%5 saldırı, 6 parça türe özgü (kumaş +10 kritik hasar, deri +3 kritik, plaka +3 hasar azaltma).
  - Adlar baskın boss'unun sıfatıyla (Ölümsüz, Tiran, Ejder Kral, Zamansız, Düşmüş, Zincirsiz), 4 dilde.
- Düşüş: baskın boss'u, üyenin zırh türüne uyan setten bir parça düşürür (`GearDrop.SetSource`, boss nadirlik tablosu).
- Set parçaları ganimet havuzuna, üretime ve denge kitlerine girmez.
- Sheet: giyilen farklı parçalar sayılır, eşiği geçen bonuslar ekipman bonuslarına eklenir (`inventory.EquippedBonuses`).
- `GET /v1/item-sets` herkese açık ve önbelleklenebilir.
- İstemci: eşya detayında set adı, giyili/toplam parça ve bonuslar; kazanılan bonuslar yeşil.

**Reason:** Yol haritası Faz D: baskın ödülüne ve uzun vadeli hedefe set eklemek. Tasarım `docs/item-engine/set.md` sadeleştirildi: imza tetikleyici ve çapraz set etkileşimi yok.
**Alternatives:** DB tablosu + `set_code` sütunu (sets.yaml bellekte yeterli). Zindan setleri (sonraki adım).
**Consequences:** Ekipman 2132 şablon. Solo zindanların seti yok. Testler: `itemset` birim testleri, `TestShippedRaidSetPieces`, `TestHTTPSetBonus`, `TestValidateSets_Rejects`; kapsam kapısına `internal/itemset`.
**Tarih:** 2026-10-19

---

## Bekleyen (Open) Kararlar — Phase 01'e Devredildi

Aşağıdakiler GDD §19'da işaretli, Phase 01 balance pass'ine bağlı:

- **DD-OPEN-1**: XP debt yüzdesi (death penalty scale).
- **DD-OPEN-2**: Auction House vergi oranı final tune (%5 sale mi %7 mi).
- **DD-OPEN-3**: Max concurrent AFK activities (1 ana + 1 pet + 1 travel önerisi).
- **DD-OPEN-4**: Race stat delta bütçesi son değer (±3 vs ±5).
- **DD-OPEN-5**: Ascension frequency final tune (60-100h aktif oyun/ascension).
- **DD-OPEN-6**: Guild max size progression Phase 02'de nerede kapsar (100 mi 150 mi).
- **DD-OPEN-7**: RNG algoritma final seçimi (xoshiro256++ vs ChaCha8).
- ~~**DD-OPEN-8**: Character slot sayısı başlangıç~~ → DD-251 (5).
- **DD-OPEN-9**: Orchestrator seçimi (Kubernetes vs Nomad vs systemd) — v1 ölçeğine göre.
- **DD-OPEN-10**: Job queue (Postgres SKIP LOCKED vs Redis stream vs NATS) — reveal batch throughput'una göre.
- **DD-OPEN-11**: Log aggregation (Loki vs ELK vs cloud provider).
- **DD-OPEN-12**: Deployment stratejisi (blue/green vs canary) — v1 için basit rolling yeterli olabilir.

Bu maddeler kesinleştikçe DD-XXX'e taşınacak.
