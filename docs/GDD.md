# Game Design Document (GDD)

**Proje Kod Adı:** Otter (çalışma adı)
**Tür:** Persistent Idle MMORPG (Kalıcı Idle/AFK odaklı Çevrimiçi Rol Yapma Oyunu)
**Platform:** Cross-platform (Web öncelikli, ileride Mobile ve Desktop istemcileri)
**Faz:** Phase 00 — Design Foundation
**Doküman Dili:** Türkçe (birincil) + İngilizce teknik terimler
**Son Güncelleme:** 2026-09-07

> Bu doküman oyunun tek bir "kuzey yıldızı" referansıdır. Her alt sistem ayrı bir dokümanda detaylandırılır; bu dosya sadece felsefeyi, temel kararları ve dokümanlar arası köprüyü kurar.

---

## 0. Doküman Haritası (Document Map)

Bu GDD aşağıdaki alt dokümanlardan beslenir. Bir sistem değiştiğinde ilgili doküman güncellenir; ardından bu dosyadaki ilgili bölüm yeniden gözden geçirilir.

### /docs/game-design/ — Oyun Tasarımı
- [Core Gameplay](game-design/core-gameplay.md) — Temel gameplay döngüsü
- [Character](game-design/character.md) — Karakter modeli
- [Races](game-design/races.md) — Irk sistemi
- [Classes](game-design/classes.md) — Sınıf sistemi
- [Specializations](game-design/specializations.md) — Uzmanlık dalları
- [Stats](game-design/stats.md) — İstatistik sistemi
- [Skills](game-design/skills.md) — Yetenek sistemi
- [Professions](game-design/professions.md) — Meslek sistemi
- [Items](game-design/items.md) — Eşya sistemi
- [Inventory](game-design/inventory.md) — Envanter
- [Equipment](game-design/equipment.md) — Ekipman
- [Combat](game-design/combat.md) — Muharebe sistemi
- [Hunting](game-design/hunting.md) — Avcılık / Mob farming
- [Adventure](game-design/adventure.md) — Macera aktiviteleri
- [Dungeons](game-design/dungeons.md) — Zindanlar
- [Pets](game-design/pets.md) — Evcil hayvanlar
- [Mounts](game-design/mounts.md) — Binekler
- [World](game-design/world.md) — Dünya yapısı
- [Map](game-design/map.md) — Harita sistemi
- [Cities](game-design/cities.md) — Şehirler
- [Travel](game-design/travel.md) — Seyahat sistemi
- [Quests](game-design/quests.md) — Görev sistemi
- [Party](game-design/party.md) — Parti sistemi
- [Guilds](game-design/guilds.md) — Loncalar
- [Economy](game-design/economy.md) — Ekonomi
- [Marketplace](game-design/marketplace.md) — Pazar/Auction
- [Events](game-design/events.md) — Etkinlikler
- [World Bosses](game-design/world-bosses.md) — Dünya boss'ları
- [Achievements](game-design/achievements.md) — Başarımlar
- [Mastery](game-design/mastery.md) — Ustalık
- [Ascension](game-design/ascension.md) — Yükseliş
- [Social](game-design/social.md) — Sosyal sistemler
- [Progression](game-design/progression.md) — İlerleme meta-sistemi

### /docs/systems/ — Motor Seviyesi Sistemler
- [AFK Engine](systems/afk-engine.md) — Offline ilerleme motoru
- [Activity Engine](systems/activity-engine.md) — Aktivite yürütücüsü
- [Reward System](systems/reward-system.md) — Ödül dağıtımı
- [Loot System](systems/loot-system.md) — Ganimet tabloları
- [Random System](systems/random-system.md) — RNG & seed
- [Notification System](systems/notification-system.md) — Bildirimler

### /docs/architecture/ — Teknik Mimari
- [System Architecture](architecture/system-architecture.md) — Uç uca teknik yığın (Flutter → API → Go → Domain Modules → PostgreSQL / Redis / Workers)
- [Flutter (Client)](architecture/flutter.md) — İstemci mimarisi, katman ayrımı, state management, WS/HTTP sözleşmesi
- [Go (Server)](architecture/go.md) — Repo yapısı, domain module sözleşmesi, concurrency, observability
- [PostgreSQL](architecture/postgresql.md) — Şema grupları, transaction politikaları, ledger discipline, backup/PITR
- [Redis](architecture/redis.md) — Ephemeral state, pub/sub, rate limit, lock, leaderboard
- [Docker](architecture/docker.md) — Containerization, multi-stage build, image standards
- [Git](architecture/git.md) — Repo yapısı, trunk-based branching, Conventional Commits, PR kuralları
- [CI/CD](architecture/cicd.md) — Pipeline türleri, deployment stratejisi, rollback
- [Environment](architecture/environment.md) — local/ci/dev/staging/prod ortam politikaları
- [Config](architecture/config.md) — 12-factor config, secret management, feature flag
- [Logging](architecture/logging.md) — Structured JSON logs, correlation, log vs audit vs metric

### /docs/database/ — Şema
- [Database README](database/README.md) — Şema grupları, konvansiyonlar, ERD, migration disiplini
- [auth-accounts](database/auth-accounts.md) — user, account, session, jwt_key
- [character](database/character.md) — character, stats, inventory, equipment, skill_learn, professions, mastery, ascension, wallet
- [content-templates](database/content-templates.md) — race, class, spec, skill, item template+instance, recipe, dungeon, adventure, npc, world_boss, loot, resource, pet/mount template, quest, achievement, mastery track, faction, currency
- [world](database/world.md) — continent → region → territory → settlement → location → poi + resource_node + waygate + travel_route
- [activity](database/activity.md) — contract, snapshot, instance, reveal_log, event_batch, death_event, rng_version
- [economy-market](database/economy-market.md) — ledger_entry, item_transfer_log, market_listing, market_transaction, mail, watchlist, price_index
- [social-guilds-parties](database/social-guilds-parties.md) — guild + rank + member + bank + invite + research, party + member + lfg, friendship + ignore, chat, notification, report
- [companions](database/companions.md) — pet_owned, mount_owned, compendium
- [progression-quests-achievements](database/progression-quests-achievements.md) — quest_progress, achievement_progress, event_participation, world_boss_spawn+participation, dungeon_lockout, ascension_history
- [audit-system](database/audit-system.md) — event_log, admin_action, security_signal, login_history, schema_migrations, feature_flag, job_queue, scheduled_task

### /docs/backend/ — Go Backend Core (Phase 01 platform)
- [Backend README](backend/README.md) — Platform katmanları, paket haritası, ortak ilkeler
- [HTTP Server](backend/http-server.md) — Bootstrap, timeouts, TLS, graceful shutdown, healthz
- [Router](backend/router.md) — chi organizasyon, versioning, path convention, pagination
- [Middleware](backend/middleware.md) — RequestID, Recover, OtelTrace, Logger, Metrics, RateLimit, Auth, CharacterScope
- [Error System](backend/error-system.md) — Sentinel + wrap + AppError + HTTP/WS/worker mapping
- [Logger](backend/logger.md) — slog, structured, context enricher, sampling
- [Database](backend/database.md) — pgx pool, sqlc pattern, repository interface, retry
- [Redis](backend/redis.md) — client wrapper, key convention, Lua scripts (rate limit + lock)
- [Transactions](backend/transactions.md) — Tx helper, isolation, retry, outbox pattern
- [WebSocket](backend/websocket.md) — Hub, multiplex, heartbeat, backpressure, cross-node fanout
- [Worker](backend/worker.md) — Job queue consumer, SKIP LOCKED, scheduler, handler katalog
- [Event Bus](backend/event-bus.md) — In-process + Redis pub/sub, outbox discipline, event kind registry

### /docs/auth/ — Faz 5: Account & Auth (Cross-Platform)
- [Auth README](auth/README.md) — Genel bakış, endpoint listesi, token model, ilkeler
- [Register](auth/register.md) — Kayıt akışı + e-mail verification + auto-login
- [Login](auth/login.md) — Password verify + MFA + token issue + anomaly
- [Logout](auth/logout.md) — Single / all-devices / admin-forced
- [Refresh Token](auth/refresh-token.md) — Opaque + rotation + replay detection
- [Password](auth/password.md) — argon2id + reset + change + common-password check
- [Session](auth/session.md) — Server-side session model + cap + listing + revoke
- [Device](auth/device.md) — Install id + push token binding + trusted device
- [Account](auth/account.md) — User↔Account model + silme + ban + GDPR + suspend

### /docs/character/ — Faz 6: Character Subsystem
- [Character README](character/README.md) — Alt sistem haritası + resource system özet
- [Character](character/character.md) — Kimlik + creation + rename + delete + session claim
- [Level](character/level.md) — XP threshold formulü + level up flow + ascension bridge
- [XP](character/xp.md) — XP kaynakları + multipliers + XP debt + XP ledger
- [Attributes](character/attributes.md) — STR/DEX/INT/WIS/STA/CHA allocation + respec + effective calc
- [Stats](character/stats.md) — Derived + tertiary formulleri + cap discipline + cache
- [Health](character/health.md) — HP pool + regen + damage + death + revive flow
- [Mana](character/mana.md) — Caster resource (Mage/Cleric)
- [Stamina](character/stamina.md) — Endurance model (Phase 02 opt-in — soft cap alternative)
- [Energy](character/energy.md) — Rogue resource + rage/focus/nature özet
- [Reputation](character/reputation.md) — Faction rep tiers + grant + tier effects
- [Titles](character/titles.md) — Owned/active title + server-first + display composition

### /docs/race/ — Faz 7: Race System (Data-Driven)
- [Race README](race/README.md) — Data-driven ilkeleri + content pack format + v1 race seti
- [Race](race/race.md) — Entity model + creation integration + change (Phase 02) + deprecation
- [Base Stats](race/base-stats.md) — Stat modifier budget + v1 tablosu + nerf compensation
- [Passive](race/passive.md) — Effect DSL (kind/target/scope) + v1 passive'ler + combiner
- [Racial Skills](race/racial-skills.md) — Aktif ability + skill bar slot 11 + cooldown + AoE
- [Racial Bonuses](race/racial-bonuses.md) — Profession affinity + exploration + unique mechanic handler pattern

### /docs/class/ — Faz 8: Class + Specialization (Data-Driven)
- [Class README](class/README.md) — Tree yapısı + data-driven ilkeleri + content pack format
- [Class](class/class.md) — Base class entity + creation + starter equipment pack
- [Specialization](class/specialization.md) — Tier I & II unlock + tree constraint + skill grant
- [Respec](class/respec.md) — Spec değişim akışı + cost curve + cooldown + free token
- [Placeholder Content](class/placeholder-content.md) — **Geçici v1 class + spec listesi (revize edilecek)**

### /docs/skill-engine/ — Faz 9: Skill + Stat Engine (Combat İçin Altyapı)
- [Engine README](skill-engine/README.md) — Katmanlar + engine flow + data flow diagramları
- [Stats](skill-engine/stats.md) — Compute pipeline + 10 katman + cap/clamp + cache
- [Modifiers](skill-engine/modifiers.md) — Modifier stack + apply order (flat/mult_add/mult_prod/set) + scope
- [Skills](skill-engine/skills.md) — Registry + effect DSL + effect kinds/targets + learn/forget
- [Passives](skill-engine/passives.md) — Aura vs trigger-based, ICD, trigger event kataloğu
- [Active Skills](skill-engine/active-skills.md) — Cast pipeline (17 adım), interrupt, target validation
- [Cooldown](skill-engine/cooldown.md) — Redis key convention + GCD + CDR + category CD + ICD
- [Mastery](skill-engine/mastery.md) — Use-based skill rank + mastery track + milestone reward
- [Scaling](skill-engine/scaling.md) — Damage/heal formula + level scale + PvE/PvP split + armor DR

### /docs/item-engine/ — Faz 10: Item Engine
- [Item Engine README](item-engine/README.md) — Engine katmanları + ana interface + data flow
- [Item Template](item-engine/item-template.md) — Data-driven template + versioning + snapshot hash
- [Item Instance](item-engine/item-instance.md) — Instance lifecycle + ownership polymorphism + transfer
- [Rarity](item-engine/rarity.md) — 7 rarity band + drop distribution + affix budget + rarity multiplier
- [Affix](item-engine/affix.md) — Pool + roll algoritması + mutual exclusion + reroll
- [Prefix](item-engine/prefix.md) — Prefix katalog + semantic rol + class bias + naming convention
- [Suffix](item-engine/suffix.md) — Suffix katalog + naming convention + grammar
- [Stats](item-engine/stats.md) — Rolled base + affix + socket + enchant + set + unique → Modifier stack
- [Sockets](item-engine/sockets.md) — Slot type + install/remove + rarity default
- [Gems](item-engine/gems.md) — Template + instance + cut states + color coding (Phase 02)
- [Enchantments](item-engine/enchantments.md) — Enchant scroll + tier + apply/remove + upgrade Phase 02
- [Set](item-engine/set.md) — 2/4/6 tier bonus + evaluation + trigger-based signature
- [Unique Effect](item-engine/unique-effect.md) — Handler-registered mekanik + legendary cap (2)

### /docs/inventory/ — Faz 11: Inventory + Equipment
- [Inventory README](inventory/README.md) — Storage türleri + ortak endpoint + ilkeler
- [Inventory](inventory/inventory.md) — Karakter bag + move/split/merge/salvage/destroy/vendor sell + overflow mailbox
- [Bank](inventory/bank.md) — Tab-based bank + location check + purchase curve + rename
- [Storage](inventory/storage.md) — Ortak Storage soyutlaması + access control matrix + overflow zinciri
- [Equipment](inventory/equipment.md) — Equip/unequip flow + 2H off-hand + legendary cap + stat recompute
- [Slots](inventory/slots.md) — Slot mekaniği + capacity + expansion + integrity worker
- [Stacking](inventory/stacking.md) — Stack rules (6 kriter) + split/merge/consolidate
- [Weight](inventory/weight.md) — Design hazır ancak **v1'de kapalı** (DD-021)
- [Item Transfer](inventory/item-transfer.md) — Cross-owner transfer + atomicity + bound/lock validation
- [Item Lock](inventory/item-lock.md) — Manuel + auto-lock + reason katalog + UI badge

### /docs/activity-engine/ — 🔥 Faz 12: Activity / AFK Engine (Oyunun Kalbi)
- [Activity Engine README](activity-engine/README.md) — **⚠️ NEVER simulate every offline tick** + genel akış
- [Activity](activity-engine/activity.md) — Ortak abstraction + Registry + Snapshot
- [Mining](activity-engine/mining.md) — Gathering: node yield + rare drop
- [Fishing](activity-engine/fishing.md) — Gathering: catch distribution + bait
- [Crafting](activity-engine/crafting.md) — Recipe iterations + success rate + item instance
- [Hunting](activity-engine/hunting.md) — AFK combat: DPS estimate + kill count + death cascade
- [Travel](activity-engine/travel.md) — Konum değişimi + duration + arrival
- [Combat](activity-engine/combat.md) — Manuel + rotation AFK + party + world boss
- [Dungeon](activity-engine/dungeon.md) — Manuel run (instance) + idle-run (contract) + weekly lockout
- [Lifecycle](activity-engine/lifecycle.md) — Start / Pause / Resume / Cancel / Complete + state machine
- [Offline Progress](activity-engine/offline-progress.md) — ⭐ Analitik hesap + Poisson/Binomial variance + threshold detection
- [Rewards](activity-engine/rewards.md) — Yield → ledger + inventory + item instance dispatch
- [XP](activity-engine/xp.md) — Character + profession + mastery XP + debt + level cascading

### /docs/world/ — Faz 13: World / Map
- [World README](world/README.md) — Hiyerarşi + content-driven pipeline + endpoint kataloğu
- [World](world/world.md) — Top-level singleton + weekly reset scheduler
- [Continent](world/continent.md) — Kıta + unlock + starter mapping per race
- [Region](world/region.md) — Biyom + climate + content ihtiyaç listesi + region mastery
- [Territory](world/territory.md) — Politik + faction control + guild war hazır
- [City](world/city.md) — Capital/regional + servis + spawn + home city + specialty
- [Village](world/village.md) — Küçük yerleşim + kısıtlı servis + specialty drop
- [Outpost](world/outpost.md) — Uzmanlaşmış hub + purpose (crafting/quest/waypoint)
- [Location](world/location.md) — POI anchor + activity_hooks + fog of war Phase 02
- [Resource Node](world/resource-node.md) — Template + instance + personal Redis TTL
- [Monster Area](world/monster-area.md) — Spawn pool + danger tier + rare mob
- [Dungeon Area](world/dungeon-area.md) — Entrance + manual/idle mode + weekly lockout

### /docs/travel/ — Faz 14: Travel (Speed Layers + AFK)
- [Travel README](travel/README.md) — Efektif süre formülü + katman özet + config
- [Travel](travel/travel.md) — Entity + lifecycle + estimate/start/current/cancel endpoint
- [Distance](travel/distance.md) — Adjacency graph + pathfind + cache
- [Base Time](travel/base-time.md) — Ana formül + segment-based compute
- [Mount](travel/mount.md) — Mount speed multiplier + type-terrain compatibility
- [Race](travel/race.md) — Racial speed passive + scope-aware
- [Weather](travel/weather.md) — Weather modifier (v1 kapalı; Phase 02 design ready)
- [Terrain](travel/terrain.md) — Terrain multiplier + mount/race interaction + road bonus
- [Equipment](travel/equipment.md) — Boots + cape + trinket + barding + speed potion
- [AFK Integration](travel/afk-integration.md) — Offline arrival + soft cap + snapshot determinism

### /docs/professions/ — Faz 15: Gathering + Professions
- [Professions README](professions/README.md) — Katalog + ortak model + XP curve + tool + race affinity
- [Mining](professions/mining.md) — Ore/gem nodes + dwarven affinity + tool tier
- [Woodcutting](professions/woodcutting.md) — Tree nodes + elven affinity
- [Fishing](professions/fishing.md) — Spot + bait + merfolk affinity
- [Hunting](professions/hunting.md) — Combat yan-effect + bestiary + orc affinity + bounty Phase 02
- [Herbalism](professions/herbalism.md) — Herb/mushroom nodes + biome distribution
- [Farming](professions/farming.md) — **Phase 02** — housing plot + crop cycle
- [Skinning](professions/skinning.md) — Hunting kill'inden pasif material
- [Archaeology](professions/archaeology.md) — **Phase 02** — dig site + culture artifact + fragment assemble

### /docs/crafting/ — Faz 16: Crafting
- [Crafting README](crafting/README.md) — 7 crafting profession + ana interface + endpoint
- [Recipes](crafting/recipes.md) — Recipe entity + öğrenme + station + kategori
- [Materials](crafting/materials.md) — Material tier + convert + cross-profession flow
- [Crafting Queue](crafting/crafting-queue.md) — Batch iteration + analytical yield + threshold
- [Quality](crafting/quality.md) — Poor/normal/exceptional + mastery bias + item_instance marking
- [Critical Craft](crafting/critical-craft.md) — Crit ×2 qty + quality bump + big miss
- [Mastery](crafting/mastery.md) — Crafting mastery track + milestone reward + grandmaster signature
- [Crafting XP](crafting/crafting-xp.md) — Character + profession + mastery XP + fail penalty

### /docs/combat/ — 🔥 Faz 17: Combat Engine
- [Combat README](combat/README.md) — Ana interface + 5 combat mod + entegrasyon
- [Combat Loop](combat/combat-loop.md) — 250ms tick + scheduler + AFK ayrımı
- [Auto-Attack](combat/auto-attack.md) — Weapon-based auto + dual wield + haste
- [Damage](combat/damage.md) — 15-adım pipeline + armor DR + crit + evade + block + absorb
- [Healing](combat/healing.md) — Heal + shield + regen + resurrect
- [Status Effects](combat/status-effects.md) — Buff/debuff/DoT/HoT/CC/shield + DR + dispel
- [Threat](combat/threat.md) — Aggro table Redis + threat modifier + taunt + swap
- [Death & Revive](combat/death-revive.md) — Death sequence + 4 revive mode + Weakness debuff
- [PvE Encounters](combat/pve-encounters.md) — AI profile + boss phases + enrage + world boss
- [Combat Events](combat/combat-events.md) — Event bus katalog + subscriber pattern

### /docs/hunting-adventure/ — Faz 18: Hunting + Adventure
- [Hunting-Adventure README](hunting-adventure/README.md) — Katman diyagramı + content pipeline
- [Hunting](hunting-adventure/hunting.md) — Deep-dive hunt aktivite + filter + party + pet
- [Encounter](hunting-adventure/encounter.md) — Random encounter framework + injection
- [Rare Enemy](hunting-adventure/rare-enemy.md) — Rare mob spawn + first-server kill + cooldown
- [Elite Enemy](hunting-adventure/elite-enemy.md) — Elite tier'lar + champion + mini-boss + group
- [Treasure](hunting-adventure/treasure.md) — Chest + world drop + lockpick Phase 02
- [Random Event](hunting-adventure/random-event.md) — Region/world event + broadcast + participation
- [Adventure](hunting-adventure/adventure.md) — Story stage-based + narration reveal + check types
- [Exploration](hunting-adventure/exploration.md) — Discovery + region completion + waypoint

### /docs/dungeon/ — 🔥 Faz 19: Dungeon Engine (AFK Çalışacak)
- [Dungeon README](dungeon/README.md) — AFK-first prensip + dungeon yapısı diyagram
- [Dungeon](dungeon/dungeon.md) — Template + instance + weekly lockout + difficulty modes
- [Floor](dungeon/floor.md) — Floor progression + checkpoint + boss per floor
- [Room](dungeon/room.md) — 8+ room kind + content pack layout
- [Encounter](dungeon/encounter.md) — Encounter framework + lifecycle
- [Enemy](dungeon/enemy.md) — Standard mob + dungeon scaling
- [Elite](dungeon/elite.md) — Elite + group + champion + mini-boss
- [Boss](dungeon/boss.md) — Phase system + enrage timer + first server clear
- [Treasure](dungeon/treasure.md) — Dungeon chest + personal instance
- [Trap](dungeon/trap.md) — Trap check + damage/CC + rogue bypass
- [Random Event](dungeon/random-event.md) — Dungeon-inject event pool + buff/curse
- [Loot](dungeon/loot.md) — Distribution + weekly lockout penalty + boss loot priority
- [AFK Run](dungeon/afk-run.md) — Idle mode analitik yield + party + lockout share

### /docs/pets/ — Faz 20: Pet System
- [Pets README](pets/README.md) — Katalog + 4 pet türü + ana yapı
- [Species](pets/species.md) — Template + 8 family + rarity + evolution chain
- [Instance](pets/instance.md) — Owned pet + Phase 20 schema (happiness/stamina/traits/skills)
- [Level](pets/level.md) — Cap 50 + XP curve (base 80, factor 1.12) + rarity stat mult
- [XP](pets/xp.md) — Kaynaklar + happiness mult formula + per-hour cap
- [Traits](pets/traits.md) — Species pool roll + effect DSL + reroll Phase 02
- [Skills](pets/skills.md) — 2 slot (active+passive) + species base + rank up
- [Evolution](pets/evolution.md) — Stage progression + level reset + trait preserve
- [Happiness](pets/happiness.md) — 0-100 stat + XP mult + decay worker + feed cooldown
- [Stamina](pets/stamina.md) — 0-100 + drain per activity + rest recovery + exhaust threshold
- [Combat](pets/combat.md) — Combat loop + auto-attack + threat 0.5 + downed 30s
- [Hunting](pets/hunting.md) — DPS boost + species affinity + skinning bonus + rare detection Phase 02

### /docs/mounts/ — Faz 21: Mount System
- [Mounts README](mounts/README.md) — 5 mount türü + pet vs mount ayrım + ilkeler
- [Mount](mounts/mount.md) — Species template + instance + Phase 21 schema (level/xp/stamina/skill_ranks)
- [Level](mounts/level.md) — Cap 30 + XP curve (base 100, factor 1.10) + travel-based XP
- [Stats](mounts/stats.md) — Speed/cargo/HP/stamina + rarity mult + cap 3.5×
- [Skills](mounts/skills.md) — Active + passive + species signature + rank up
- [Travel](mounts/travel.md) — Travel integration + snapshot lock + combat dismount
- [Terrain](mounts/terrain.md) — Compatibility multiplier per species + mount type bias
- [Carry Capacity](mounts/carry-capacity.md) — Cargo bonus slots + dismount overflow

### /docs/quests/ — Faz 22: Quest System
- [Quests README](quests/README.md) — Katalog + ortak sözleşme + objective türleri
- [Main](quests/main.md) — Story chain + prerequisite + reward-heavy
- [Side](quests/side.md) — Region-bound + opsiyonel + mini-chain
- [Daily](quests/daily.md) — Rotasyonlu havuz + 3 slot + faction rep farm
- [Weekly](quests/weekly.md) — 5-10 slot + Salı 06:00 UTC reset + weekly_seal
- [Dynamic](quests/dynamic.md) — Procedural/bounty board + rare mob auto-quest
- [Event](quests/event.md) — Event-tied + event_token + participation
- [Guild](quests/guild.md) — **Phase 02** guild-shared + contribution ledger

### /docs/party/ — Faz 23: Party System
- [Party README](party/README.md) — Katalog + endpoint özet
- [Party](party/party.md) — Entity + lifecycle + Phase 23 schema (size_cap/persistent/name)
- [Leader](party/leader.md) — Yetkiler + promote + auto-migration
- [Member](party/member.md) — Join/leave/kick + ready + cross-region
- [Invite](party/invite.md) — Invite flow + LFG listing + matchmaking Phase 02
- [Roles](party/roles.md) — Tank/heal/dps/support + class default + composition Phase 02
- [Combat](party/combat.md) — Aura sync + party bonus +5% + kill credit
- [Hunt](party/hunt.md) — Per-character contract + party bonus + rare mob
- [Dungeon](party/dungeon.md) — Manuel instance + idle-run + wipe recovery
- [Loot](party/loot.md) — 4 rule (personal default) + boss priority + contribution reward Phase 02

### /docs/guild/ — 🔥 Faz 24: Guild System
- [Guild README](guild/README.md) — Katalog + endpoint özet
- [Guild](guild/guild.md) — Entity + lifecycle + disband 48h grace + create fee 10k gold
- [Members](guild/members.md) — Invite 7d + join + kick + contribution + cap 50
- [Roles](guild/roles.md) — 5 default rank (leader/officer/veteran/member/initiate)
- [Permissions](guild/permissions.md) — Rank matrix + bank withdraw daily cap
- [Level](guild/level.md) — Cap 10 v1 (20 Phase 02) + member cap expansion
- [XP](guild/xp.md) — 9 XP kaynak + per-member 5000/day cap
- [Bank](guild/bank.md) — 5 tab (level-gated) + treasury + audit log 12ay
- [Research](guild/research.md) — **Phase 02** — passive perks (combat/gather/social/economic)
- [Quest](guild/quest.md) — **Phase 02** — guild-shared contribution
- [Boss](guild/boss.md) — **Phase 02** — guild-summoned boss + weekly lockout
- [Raid](guild/raid.md) — **Phase 02** — 10-25 kişi + master loot
- [War](guild/war.md) — **Phase 02+** — territory control + PvP + declaration flow

### /docs/economy/ — 🔥 Faz 25: Economy Engine
- [Economy README](economy/README.md) — Katalog + altın kural (ledger-first)
- [Gold](economy/gold.md) — Ana currency (copper base unit) + wallet + cap + kaynak/sink
- [Currency](economy/currency.md) — Taksonomi (gold/gem/honor/faction/event/guild_contribution)
- [Vendor](economy/vendor.md) — NPC buy/sell + buyback 2h TTL + fiyat formülü
- [Trade](economy/trade.md) — P2P trade window + anti-scam auto-unlock + finalize TX
- [Transaction](economy/transaction.md) — Double-entry ledger + SERIALIZABLE + reconciliation
- [Tax](economy/tax.md) — Vendor spread + trade tax (Phase 02) + auction (Faz 26) + repair tax
- [Fee](economy/fee.md) — Service fee (guild create/rename/respec/repair/mail/teleport)

### /docs/marketplace/ — 🔥 Faz 26: Player Market + Auction
- [Marketplace README](marketplace/README.md) — Katalog + global market kısıtları
- [Listing](marketplace/listing.md) — Lifecycle (active/sold/expired/cancelled) + escrow model
- [Sell](marketplace/sell.md) — Post akış + escrow lock + snapshot + fee tahsili
- [Buy](marketplace/buy.md) — Instant buy atomic swap + race guard + same-account block
- [Cancel](marketplace/cancel.md) — Iptal + kademeli fee iade + overflow claim
- [Search](marketplace/search.md) — Query API + Redis cache 30sn + rate limit
- [Filter](marketplace/filter.md) — 15+ filter dimension + snapshot JSONB + GIN index
- [Price](marketplace/price.md) — Min/max clamp + suggested price (median 7d)
- [Tax](marketplace/tax.md) — Listing fee %5 × duration mult + sale cut %5
- [Auction](marketplace/auction.md) — **Phase 02 primary** bid-based + anti-snipe 5dk×6
- [Market History](marketplace/market-history.md) — Sale history + percentile stats 7d

### /docs/world-events/ — 🔥 Faz 27: World Events
- [World Events README](world-events/README.md) — Katalog + AFK katılım (shared timer + personal contribution)
- [World Boss](world-events/world-boss.md) — Global + attendance contract 24h + 4-tier reward
- [Invasion](world-events/invasion.md) — Region multi-wave + auto-eligible + success/fail + 48h buff
- [Regional Event](world-events/regional-event.md) — Periodic activity multiplier + region buff
- [Seasonal Event](world-events/seasonal-event.md) — 2-6 hafta + event token BoA + AFK login %50 credit
- [Dynamic Event](world-events/dynamic-event.md) — Emergent chain event (5 objective + max 3 derinlik)

### /docs/social/ — 🟢 Faz 28: Social
- [Social README](social/README.md) — Katalog + kanal taksonomisi
- [Friends](social/friends.md) — Friend request + presence + cap 100
- [Block](social/block.md) — Account-level block + interaction gate
- [Chat](social/chat.md) — Chat engine + WS delivery + rate limit + moderation
- [Private Message](social/private-message.md) — Whisper 1-to-1 + offline queue + unread
- [Guild Chat](social/guild-chat.md) — Guild + officer + MOTD + presence feed
- [Party Chat](social/party-chat.md) — Party low-latency (< 200ms) + system message
- [Global Chat](social/global-chat.md) — Region-scoped global/trade/lfg/region + slow mode
- [Mail](social/mail.md) — Attachment (12 item + gold) + escrow + 30d expire iade

### /docs/progression/ — 🟢 Faz 29: Long-Term Progression
- [Progression README](progression/README.md) — Katalog + endgame felsefe + 8 sistem
- [Mastery](progression/mastery.md) — Cross-system mastery (combat/gathering/crafting/…) 10 tier
- [Ascension](progression/ascension.md) — Rebirth + prestige tree + 60-100h/cycle
- [Achievements](progression/achievements.md) — Trigger + puan + meta + retroactive
- [Titles](progression/titles.md) — Extended title (retirement + rarity + favorite)
- [Collections](progression/collections.md) — Pet/mount/appearance/tabard + set bonus (transmog Phase 02)
- [Milestones](progression/milestones.md) — Zaman/hacim sayaç + chain + retroactive
- [Reputation](progression/reputation.md) — 8-tier + daily cap 1500 + hostility gate
- [Faction](progression/faction.md) — Faction taxonomy + relationship graph + cross-effect

### /docs/notification/ — 🟢 Faz 30: Notification + Mobile Background
- [Notification README](notification/README.md) — Katalog + 17 notification türü + prensipler
- [Push Notification](notification/push-notification.md) — FCM (Android) + APNs (iOS) + token yönetim + dedup
- [Activity Completed](notification/activity-completed.md) — AFK activity reveal push + grouped notification
- [Dungeon Completed](notification/dungeon-completed.md) — Dungeon AFK run + success/fail
- [Craft Completed](notification/craft-completed.md) — Queue per_slot vs queue_drained mode
- [Travel Completed](notification/travel-completed.md) — Arrived + interrupted + < 60sn suppress
- [Energy Restored](notification/energy-restored.md) — Opt-in silent push + energy_full_at prediction
- [World Event](notification/world-event.md) — World boss/invasion/seasonal + eligibility filter
- [Guild Event](notification/guild-event.md) — Boss/raid/war/invite + officer action
- [Mobile Lifecycle](notification/mobile-lifecycle.md) — Android Doze + iOS Background Modes + WS reconnect + deep link
- [Preferences](notification/preferences.md) — Granular opt-out + quiet hours + global mute + uncontrollable

### /docs/ui/ — 🔵 Faz 31: Flutter UI/UX
- [UI README](ui/README.md) — Katalog + prensipler + bilgi mimarisi
- [Design System](ui/design-system.md) — Token + tipografi + motion + breakpoints
- [Login](ui/login.md) — Login/register/OTP/reset
- [Character](ui/character.md) — Select + create wizard + delete
- [Character Sheet](ui/character-sheet.md) — Stats + resource bars + buff/debuff
- [Inventory](ui/inventory.md) — Grid + filter/sort + batch + item detail sheet
- [Equipment](ui/equipment.md) — Paper doll + set indicator + compare tooltip
- [Map](ui/map.md) — World/region/territory navigator + POI + travel initiate
- [Activity](ui/activity.md) — AFK contract + reveal + estimate + concurrent cap
- [Profession](ui/profession.md) — Multi-profession dashboard + mastery
- [Crafting](ui/crafting.md) — Recipe + queue + critical craft alert
- [Combat](ui/combat.md) — Live combat panel + AFK combat state + death screen
- [Dungeon](ui/dungeon.md) — Dungeon list + entry + AFK progress + reveal
- [Pet](ui/pet.md) — Pet list + active + detail
- [Mount](ui/mount.md) — Mount list + terrain bonus + skin
- [Party](ui/party.md) — Party HUD + full screen + invite + roles
- [Guild](ui/guild.md) — Guild dashboard + members + bank + chat + MOTD
- [Market](ui/market.md) — Browse + my listings + history chart
- [Quest](ui/quest.md) — Log + tracker + turn-in
- [Social](ui/social.md) — Friends + mail + chat + block sub-tabs

### /docs/admin/ — 🔵 Faz 32: Admin Panel
- [Admin README](admin/README.md) — Katalog + prensipler (read-first, audit-first)
- [Auth + RBAC](admin/auth-rbac.md) — SSO + WebAuthn 2FA + 6 rol + step-up
- [Players](admin/players.md) — Account view + PII reveal audit + mute/ban
- [Characters](admin/characters.md) — Character inspect + freeze + item restore
- [Items](admin/items.md) — Template + instance browser + forensic history
- [Economy](admin/economy.md) — Ledger + wallet + sink/source + reversal
- [Guild](admin/guild.md) — Guild list + members + bank audit + disband override
- [Market](admin/market.md) — Listing viewer + admin cancel + RMT flag Phase 02
- [Combat](admin/combat.md) — Session inspect + suspicious pattern
- [Activities](admin/activities.md) — Active AFK contract + force-reveal
- [Dungeons](admin/dungeons.md) — Active runs + template + lockout reset
- [Events](admin/events.md) — World event scheduler + manual trigger
- [Logs](admin/logs.md) — Aggregate audit viewer + correlation + export

### /docs/security/ — 🔴 Faz 33: Security + Anti-Cheat
- [Security README](security/README.md) — Threat model + 11 tehdit vektörü + prensipler
- [Rate Limit](security/rate-limit.md) — Token bucket + endpoint kataloğu + circuit breaker
- [Request Validation](security/request-validation.md) — Schema + sanitization + UTF-8 normalize
- [Transaction Validation](security/transaction-validation.md) — Ledger integrity + outbox + reconciler
- [Inventory Validation](security/inventory-validation.md) — 12 check + lock + fingerprint
- [Currency Validation](security/currency-validation.md) — Cap + tradable + scope + system hesabı
- [Trade Validation](security/trade-validation.md) — Bait-and-switch guard + finalize recheck
- [Replay Protection](security/replay-protection.md) — Idempotency key + nonce + WS dedup
- [Session Security](security/session-security.md) — JWT EdDSA + refresh rotation + geo anomaly
- [Audit + Response](security/audit.md) — Anomaly detection + alert sınıfı + incident runbook

### /docs/testing/ — 🔴 Faz 34: Automated Testing
- [Testing README](testing/README.md) — Test piramidi + prensipler + CI aşamaları + fixture yönetimi
- [AFK](testing/afk.md) — Analytical hesap + snapshot + idempotent reveal + threshold inverse
- [Combat](testing/combat.md) — 15-step pipeline + threat + CC diminishing + death penalty
- [Loot](testing/loot.md) — Deterministic RNG + personal loot + chi-square dist + rarity
- [Crafting](testing/crafting.md) — Queue lifecycle + critical craft + material atomic
- [Inventory](testing/inventory.md) — Ownership + lock + weight cap + stack + dupe prevent
- [Economy](testing/economy.md) — Ledger sum=0 + idempotency + wallet cap + reconciler + reversal
- [Trade](testing/trade.md) — State machine + auto-unlock + confirm hash + bind bypass + block
- [Dungeon](testing/dungeon.md) — Floor progression + weekly lockout + AFK run + first-kill
- [Character](testing/character.md) — Lifecycle + XP/level + freeze + delete grace + ascension
- [Class](testing/class.md) — Class immutable + spec pick + skill tree + respec + restriction

### /docs/performance/ — 🔴 Faz 35: Performance + Load Test
- [Performance README](performance/README.md) — CCU tier T1-T5 + SLO + tooling (k6)
- [SLI + SLO](performance/slo-slis.md) — Full endpoint tablo + error budget + burn rate
- [CCU Scenarios](performance/ccu-scenarios.md) — Karışım (3:7 active:AFK) + burst modeli + fixture
- [Login](performance/load-login.md) — Steady + burst + storm + Argon2 tuning
- [Combat](performance/load-combat.md) — Skill cast + boss peak + party fanout + damage pipeline
- [Market](performance/load-market.md) — Search cache + buy race + list storm + SERIALIZABLE retry
- [Guild](performance/load-guild.md) — Home compound + chat fanout + bank concurrent
- [AFK Calculation](performance/load-afk-calc.md) — Reveal throughput + wave + notification push
- [Dungeon](performance/load-dungeon.md) — Enter + AFK reveal + weekly reset burst
- [WebSocket](performance/load-websocket.md) — 10-50k conn + fanout + reconnect wave + backpressure
- [Capacity Planning](performance/capacity-planning.md) — Infrastructure sizing T1-T4 + cost projection

### /docs/infra/ — 🔴 Faz 36: Production Infrastructure
- [Infra README](infra/README.md) — Prod philosophy + immutable + everything-as-code
- [Docker](infra/docker.md) — Distroless nonroot + multi-stage + cosign + Trivy scan
- [Reverse Proxy](infra/reverse-proxy.md) — HAProxy + TLS termination + WS proxy + rate limit
- [CDN](infra/cdn.md) — Cloudflare + static asset + WAF + DDoS mitigation
- [SSL](infra/ssl.md) — ACME + cert-manager + TLS 1.3 + HSTS + rotation
- [Postgres Prod](infra/postgres-prod.md) — HA + replication + PgBouncer + tuning + migration
- [Redis Prod](infra/redis-prod.md) — Sentinel/Cluster + persistence + eviction + pubsub scale
- [Backup](infra/backup.md) — WAL archive + PITR + Redis snapshot + monthly restore drill
- [Monitoring](infra/monitoring.md) — Prometheus + Grafana + Jaeger/Tempo + SLO dashboard
- [Logging](infra/logging.md) — Structured JSON + Loki + correlation_id + PII redact
- [Alerting](infra/alerting.md) — Alertmanager + PagerDuty + 4-tier + runbook + post-mortem
- [CI/CD Prod](infra/ci-cd-prod.md) — Canary Argo Rollouts + feature flag + zero-downtime + DORA
- [Rollback](infra/rollback.md) — App/flag/migration/content/DR restore + decision tree + drill

### /docs/release/ — 🟣 Faz 37: Closed Beta / Release
- [Release README](release/README.md) — Launch roadmap + phase gates + KPI + timeline (3-5 ay)
- [Production Readiness](release/production-readiness.md) — Go-live checklist (tech + biz + legal + ops)
- [Closed Beta](release/closed-beta.md) — 3 wave (200 → 5k) + telemetry + feedback + wipe politikası
- [Bug Fix](release/bug-fix.md) — Triage P0-P3 + hotfix path + rollback vs hotfix
- [Balance](release/balance.md) — Content-driven pass + weekly cycle + A/B test Phase 02
- [Performance Tuning](release/performance-tuning.md) — Regression detection + bottleneck kalıpları + RUM
- [Open Beta](release/open-beta.md) — Public 10-25k + marketing kick-off + T4 stress + wipe
- [Release v1.0](release/release.md) — Launch day playbook + war room + hotfix protocol
- [Post-Launch](release/post-launch.md) — First 30d + Phase 02 roadmap kick-off + monetization

### /docs/decisions/ — Kararlar
- [Design Decisions](decisions/design-decisions.md) — Alınmış tasarım kararlarının günlüğü

---

## 1. Game Vision (Oyun Vizyonu)

**Kısa cümle:** *"Oyuncunun karakteri, o klavyeden kalktığında bile yaşamaya devam eden; giriş yaptığında ilerlemenin biriktiğini görmenin tatmin verici olduğu, derin bir Idle MMORPG."*

### 1.1 Uzun Vizyon
- Karakterin varlığı, oyuncunun aktif olmasına bağlı değildir. **Zaman en değerli kaynaktır** ve oyunun her sistemi zamanı akıllıca yatırım yapmayı ödüllendirir.
- **Aktif oynayan** ve **bırakıp gelen** oyuncular aynı ekonomide barış içinde var olur. Aktif oyun bir çarpan sağlar; ancak oyunun **kapısı** değildir.
- Oyuncu 15 dakika içinde giriş yapıp anlamlı bir karar (aktivite değiştirme, market işlemi, guild etkileşimi) verip çıkabilmelidir.
- MMO hissi (paylaşılan dünya, ekonomi, sosyal katman) korunur; ancak *"5 saatini akşam ayırmak zorunda değilsin"* prensibi vardır.

### 1.2 Vizyon Dışı (Non-Goals)
- Reaction-based / twitch skill oynanış değildir. Muharebe **karar odaklıdır**, refleks odaklı değildir.
- Pay-to-win değildir. Zaman ve bilgi kazandırır; para en fazla **kozmetik ve konfor** satın alabilir.
- "Log in every 4h or lose stuff" tarzı sömürücü zorlayıcı loop içermez.
- Karmaşık, yüz saatlik tutorial gerektiren rakamsal iç mekanik değildir. Karmaşıklık **derinliğe** hizmet eder; sürtünme yaratmak için değildir.

---

## 2. Game Pillars (Oyunun Ayakları)

Her tasarım kararı bu üç sütunla test edilir. Bir öneri sütunları güçlendiriyorsa evet, zayıflatıyorsa hayırdır.

### 2.1 Persistent Progression (Kalıcı İlerleme)
Karakter her zaman bir şey yapıyor. Offline geçen zaman kaybolmuş zaman değildir; hesap sırasında **ödüle** dönüşür. AFK bir "afet önleme" mekaniği değil, **birinci sınıf gameplay** modudur.

### 2.2 Meaningful Choices (Anlamlı Seçimler)
Her seçim (class, spec, gear, aktivite, rota, market) fırsat maliyeti yaratır. "Hepsi eşit güçte, sadece görsel fark" tuzağına düşülmez; ancak "yanlış build seçtim, silmem lazım" acımasızlığı da tercih edilmez. **Respec** yolları vardır, ama bedeli vardır.

### 2.3 Shared, Living World (Paylaşılan, Yaşayan Dünya)
Ekonomi tek dünya çapındadır. Loncalar bölgesel etki edinir. Dünya boss'ları sunucu genelinde olay olur. Oyuncu, tek başına oynasa bile *"burada başka insanlar da var"* hissini kaybetmez.

---

## 3. Core Gameplay Loop (Ana Gameplay Döngüsü)

Ayrıntı: [core-gameplay.md](game-design/core-gameplay.md).

**Kısa döngü (Session Loop, ~5–15 dk):**
1. Login → offline ödülleri topla (AFK Engine)
2. Envanteri temizle / market'e ürün koy
3. Aktif aktiviteyi değiştir (mining → hunting, quest teslim, dungeon planlama)
4. Kısa aktif etkileşim (isteğe bağlı: manuel muharebe, boss katkısı, guild sohbet)
5. Yeni aktivite başlat → logout

**Orta döngü (Daily Loop, günlük):**
- Günlük görevler, dungeon reset'i, weekly market rotasyonu, event katılımı

**Uzun döngü (Meta Loop, haftalar):**
- Class specialization ilerlemesi, mastery, ekipman upgrade, guild progression, ascension hazırlığı

---

## 4. Player Experience (Oyuncu Deneyimi)

Hedeflenen duygular, öncelik sırasına göre:

1. **"Karakterim beni bekliyordu."** — Login anında biriken ödülleri toplamak.
2. **"Karar verdim, dünya buna cevap verdi."** — Aktivite seçiminin gerçekten ilerlemeyi değiştirdiğini görmek.
3. **"Başkalarıyla aynı dünyadayım."** — Ekonomi, guild, boss anlarında paylaşılan alan hissi.
4. **"Kendi ritmimde ilerliyorum."** — Baskı yok; sistem oyuncuyu sömürmüyor.

### 4.1 Anti-Patterns (Yaratmak İstemediğimiz Duygular)
- FOMO paniği (kaçırılan bir şey için giriş yapmak zorunda kalmak).
- "İşe gitmek gibi" günlük checklist zorunluluğu.
- Grind için grind (ödülü olmayan tekrar).
- Yalnızlık: bir MMO'da geçirilen 20 saatte tek bir başka oyuncu ismi görmemek.

---

## 5. Character Progression (Karakter İlerlemesi)

Detay: [progression.md](game-design/progression.md), [character.md](game-design/character.md), [mastery.md](game-design/mastery.md), [ascension.md](game-design/ascension.md).

Üç paralel ilerleme ekseni:

| Eksen | Nasıl Artar | Sıfırlanır mı? |
|---|---|---|
| **Character Level** | XP (aktivite + combat) | Ascension'da resetlenir |
| **Profession Level** | İlgili meslek aktiviteleri | Genellikle hayır |
| **Mastery** | Uzun vadeli kullanım, koleksiyon | Kalıcı |

Buna ek olarak:
- **Class Specialization Tree** — seçimlere bağlı; respec ile değiştirilebilir.
- **Ascension Layer** — endgame; karakter reset karşılığında kalıcı meta bonuslar.

---

## 6. World Structure (Dünya Yapısı)

Detay: [world.md](game-design/world.md), [map.md](game-design/map.md), [cities.md](game-design/cities.md), [travel.md](game-design/travel.md).

Hiyerarşi:
```
World
 └─ Continent
     └─ Region
         └─ Territory
             ├─ City
             ├─ Village
             ├─ Outpost
             └─ Location (resource node, dungeon entrance, quest hub, POI)
```

- Şehirler pasif "menü ekranı" değildir; her şehir kendi profession, vendor, faction, market erişimini sunar.
- Bölgeler seviye bantları değil, **kimlikler** taşır (klimat, kaynaklar, kültür, tehditler).
- Travel gerçek zamanlıdır ve AFK bir aktivite tipidir; teleport ekonomi maliyetlidir.

---

## 7. Combat Philosophy (Muharebe Felsefesi)

Detay: [combat.md](game-design/combat.md), [skills.md](game-design/skills.md).

- **Karar odaklı**, tick-tabanlı. Server otoritedir; istemci sunum yapar.
- Auto-attack varsayılan; skill'ler manuel veya otomasyon (rota/priority list) ile kullanılabilir.
- AFK muharebe: karakter, tanımlı rotation'ı server üzerinde çalıştırır.
- Ölüm cezası vardır ama karakter silinmez; XP debt / durability loss / opportunity cost gibi katmanlı bedeller.

---

## 8. Idle / AFK Philosophy (Idle Felsefesi)

Detay: [afk-engine.md](systems/afk-engine.md), [activity-engine.md](systems/activity-engine.md).

- AFK **birinci sınıf gameplay**. Oyunun tüm ana aktiviteleri (mining, fishing, woodcutting, hunting, crafting, travel, training, dungeon-run) offline yürütülebilir.
- **Cap yerine yumuşak azalma** (soft cap): Belirli bir sürenin ötesinde verim azalır; ancak hiçbir zaman "zamanını çöpe attın" olmaz.
- Server authoritative: Oyuncu offline iken karakterin durumu server tarafında bilinir; login sırasında hesaplanır, istemciye "reveal" edilir.
- Tüm aktiviteler tek bir Activity Engine soyutlaması altındadır — yeni aktivite eklemek sistemi çatlatmaz.

---

## 9. Economy Philosophy (Ekonomi Felsefesi)

Detay: [economy.md](game-design/economy.md), [marketplace.md](game-design/marketplace.md), [items.md](game-design/items.md).

- **Faucets & Sinks**: Sisteme giren para/eşya kadar sistemden çıkan para/eşya olmalı. Her yeni içerik iki tarafı da düşünmeli.
- **Sıkı takip edilen üç metrik:** Toplam para arzı, günlük para dolaşımı, temel commodity fiyatları.
- Auction House global; region-specific alt pazarlar ileride açılabilir.
- Player-to-player trade var; ancak dolandırıcılığı önlemek için escrow ve confirm-two-sides kuralları.

---

## 10. Social Philosophy (Sosyal Felsefe)

Detay: [social.md](game-design/social.md), [guilds.md](game-design/guilds.md), [party.md](game-design/party.md).

- Yalnız oynanabilir bir MMO. Ancak sosyal katmana adım atmak **her zaman daha iyi** hissettirir; **zorunlu** hissettirmez.
- Guild sadece chat kanalı değildir; ortak progression, guild hall, guild research/craft, guild bank.
- Party küçük ölçekli (2–5) ve olaya-özel (dungeon, hunt, event); guild kalıcı sosyal birimdir.

---

## 11. Content Philosophy (İçerik Felsefesi)

- Her yeni içerik parçası (yeni bölge, dungeon, item seti, meslek dalı) mevcut sistemlerin **kombinasyonlarından** oluşabilmelidir. Data-driven içerik.
- Content patch'leri "raid + gear tier" değil; **yeni oyun tarzı** eklemelidir (yeni activity type, yeni social loop, yeni economy sink).
- **Evergreen content** > seasonal churn. Season'lar var ama karakteri sıfırlamaz; sadece kozmetik/leaderboard ekler.

---

## 12. Endgame Philosophy (Endgame Felsefesi)

Detay: [ascension.md](game-design/ascension.md), [mastery.md](game-design/mastery.md), [world-bosses.md](game-design/world-bosses.md).

- Level cap oyunun bitişi değil, endgame'in başlangıcıdır.
- Endgame üç sacayak: **Ascension** (dikey ilerleme), **Mastery** (yatay derinlik), **Social/Competitive** (guild, leaderboard, dünya boss, ileride PvP).

---

## 13. Monetization Principles (Para Kazanma İlkeleri)

> Phase 00'da monetization uygulanmaz; ancak tasarım şu ilkelerle uyumlu kalmalıdır.

- **Never sell power.** Stat, XP boost oranı, ekipman satılmaz.
- **Konfor & kozmetik satılabilir:** ekstra karakter slotu, cosmetic transmog, kozmetik pet/mount, isim değiştirme, envanter genişlemesi (mantıklı sınırlar içinde).
- **Time is not the product to sell.** "Beklemeyi atla" mekaniği yoktur; çünkü beklemek zaten oyunun kendisidir.
- Battle pass olursa: sadece kozmetik + biraz konfor; içerik kilitleme yok.

---

## 14. Cross-Platform Philosophy (Çapraz Platform Felsefesi)

- Otorite server'dadır; istemci tipi belirleyici değildir. Web istemcisi ve gelecekteki native istemciler aynı state'i gösterir.
- UI tasarımı **mobile-first thinking, desktop-scaled**: küçük ekranda okunabilen bir sistem her yerde çalışır.
- Push notification desteği (opt-in): "AFK aktiviten bitti", "market'te ürünün satıldı", "guild boss başlıyor".

---

## 15. Technical Design Principles (Teknik Tasarım İlkeleri)

Detaylı teknik gereksinimler: bkz. Bölüm 18.

- **Server Authoritative**: Tüm state, para, envanter, XP hesabı serverda tutulur; istemci öneri gönderir, server doğrular.
- **Deterministic Simulation for AFK**: Offline hesaplama, seed + timestamp ile deterministik olarak tekrarlanabilir.
- **Data-Driven**: Yeni class, race, item, aktivite eklemek için kod değişikliği yerine data tanımı.
- **Time Source of Truth**: Server saati mutlak referans. İstemci saati asla ödül hesabında kullanılmaz.
- **Event-Sourced Progress**: Karakter değişimleri (XP kazanımı, item drop, ölüm) olay olarak loglanır; state olayların projeksiyonudur.

---

## 16. Future Expansion Philosophy (Gelecek Genişleme Felsefesi)

Şu anda geliştirilmeyecek, ancak mimaride "yer ayrılmış" sistemler:

- **PvP** (structured & open-world)
- **Arena** (ranked, brackets)
- **Raid** (10+ kişi, weekly reset)
- **Guild War** & **Territory War**
- **Player Housing**
- **Farming / Homestead**
- **Seasonal World Events**
- **Achievements & Collections** (v1'de light versiyon)
- **Factions & Reputation Grinding**
- **Transmog** sistemi (kozmetik override)
- **Enchanting Expansion** (rune, catalyst, upgrade path)
- **Legendary Progression** (uzun soluklu craft quest'ler)
- **Endgame Prestige Systems** (Ascension'ın ötesinde)

Bu sistemlerin her biri için: veri şeması, isim alanı, entity ilişkileri şimdiden **çakışmayacak** şekilde planlanır (`character.legendary_progress` alanı boş kalabilir, ama vardır).

---

## 17. Cross-System Consistency Map (Sistem Arası Bağlantı Haritası)

Bir sistem değiştiğinde etkilenecek diğer sistemleri hızlı görmek için:

```
Race          → Stats, Skills, Professions, Combat
Class         → Stats, Skills, Combat, Progression
Specialization→ Class, Skills, Combat
Skills        → Combat, Class, Specialization
Stats         → Combat, Items, Equipment, Race, Class
Items         → Combat, Economy, Marketplace, Loot, Inventory, Equipment
Equipment     → Stats, Combat, Items
Combat        → Skills, Items, Stats, Pets, Dungeons, Hunting, World Bosses
Pets          → Combat, Activity, AFK Engine
Mounts        → Travel, World, Activity
Travel        → World, Map, AFK Engine, Activity Engine
Activity      → AFK Engine, Reward System, Loot System
AFK Engine    → Activity Engine, Random System, Reward System, Notification
Professions   → Items, Economy, Activity Engine
Dungeons      → Combat, Party, Loot, Activity Engine
Party         → Combat, Dungeons, Hunting, Social
Guild         → Social, Economy, Progression, (future) Guild War
Economy       → Items, Marketplace, Professions, Loot
Marketplace   → Economy, Items, Notification
Quests        → Reward System, World, Progression
World Boss    → Combat, Loot, Notification, Guild
Ascension     → Progression, Character, Stats
Mastery       → Progression, Activity, Character
```

Ayrıntılı grafiği ve tutarsızlık takibi için bölüm 19'a bakınız.

---

## 18. Phase 01 Teknik Gereksinim Özeti (Kısa Liste)

> Bu bölüm Phase 00'ın çıktısı olarak Phase 01'in gireceği zemini tanımlar. Detaylı teknik tasarım Phase 01'in görevidir.

**Yığın (locked, 2026-09-08):**
```
Flutter (client) → API (HTTPS + WSS) → Go Application Server
  → Domain Modules → PostgreSQL (durable) + Redis (ephemeral) + Workers (async)
```
Detay: [System Architecture](architecture/system-architecture.md). Karar: [DD-026](decisions/design-decisions.md).

### 18.1 Zorunlu Altyapı (Must-Have)
1. **Server-authoritative game loop** — istekleri doğrulayan, state'i tutan sunucu.
2. **Time service** — Tek monotonik zaman kaynağı; NTP-safe; istemci clock skew tolere edilebilir ama otorite değildir.
3. **AFK Engine calculator** — Aktivite tipi + başlangıç zamanı + son zaman + karakter durumu → ödüller.
4. **Activity Engine** — Aktivite tanımı, life-cycle (start/tick/interrupt/end), event yayını.
5. **Deterministic RNG** — Seed + tick indeksi tabanlı; log ve replay için tekrarlanabilir.
6. **Persistence** — Karakter, envanter, meslek seviyeleri, aktif aktivite, offline snapshot.
7. **Auth & Session** — Hesap, karakter slot, login/logout event'leri.
8. **Content pipeline** — Race/Class/Item/Skill data'sını hot-reloadable veri olarak sunar.
9. **Event log / audit trail** — Ekonomi olayları için yasal-forensic seviyede iz.
10. **Notification queue** — Push, in-game, e-posta hedeflerine yönlendirme.

### 18.2 Öncelikli İstenen (Should-Have)
- WebSocket / gRPC-web canlı push
- Redis benzeri ephemeral state (aktif aktivite ticker'ı)
- SQL (durable) — postgres önerilir
- Rate limiting & anti-cheat sinyal katmanı

### 18.3 Ertelenebilir (Could-Have)
- Sharding stratejisi (v1 tek dünya, tek shard varsayımı)
- Cross-region latency çözümleri

---

## 19. Consistency & Gap Tracker (Tutarlılık ve Boşluk Takibi)

Phase 00 sonu itibarıyla tespit edilen ve **ileride kararı verilmesi gereken** açık noktalar:

1. **Death penalty ölçeği** — XP debt mi, durability loss mı, ikisi mi? İlk teklif: kombine, ama ölçek Phase 01 balancing pass'ine bırakıldı.
2. **Auction House vergisi %** — 5% listing + 5% sale önerisi ekonomide test edilecek.
3. **Max concurrent AFK activities** — 1 ana + 1 pet + 1 travel önerisi; kilit karar Phase 01.
4. **Race stat delta büyüklüğü** — ±%3 civarı hedef; ancak class ile toplam etkileşim Phase 01 sim'iyle test edilecek.
5. **Ascension frequency** — hedef: "aktif bir oyuncu için 60–100 saat aktif oyun sonrası ilk ascension".
6. **Guild max size** — 50 başlangıç önerisi; ölçek Phase 02'de guild war ile revizyona açık.

Bu maddelerin her biri, kararlaştıklarında [Design Decisions](decisions/design-decisions.md) dosyasına taşınır.

---

## 20. Doküman Konvansiyonları

- **Turkish primary**, İngilizce teknik terim korunur.
- Her sistem dokümanı standart iskeleti izler: *Amaç, Fayda, Gameplay Loop içindeki yer, Kurallar, Kullanıcı akışı, İlişkiler, Veri, I/O, Progression, Reward, Failure, Edge cases, AFK davranışı, Server authoritative gereksinimleri, Genişletilebilirlik, Örnek senaryolar, Balance considerations.*
- Sistem çakıştığında: **AFK Engine > Activity Engine > tekil aktivite kuralı**. Yani AFK Engine bir kuralı override ederse alt sistem uyum sağlar.
- Cross-reference: `[metin](path)` markdown link ile.
- Kararlar: mutlaka [decisions/design-decisions.md](decisions/design-decisions.md) dosyasında tarih ve gerekçe ile.

---

## 21. Değişiklik Günlüğü

| Tarih | Değişiklik | Not |
|---|---|---|
| 2026-09-07 | İlk sürüm (Phase 00 kapanışı) | Tüm alt dokümanlarla birlikte oluşturuldu |
