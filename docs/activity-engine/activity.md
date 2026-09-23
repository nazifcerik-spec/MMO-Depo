# Activity (Ortak Abstraction)

## 1. Amaç
Tüm aktivite tiplerinin ortak sözleşmesi. Mining, hunting, travel, combat — hepsi `Activity` interface'i implement eder.

## 2. Interface

```go
package activityengine

type Activity interface {
    // Metadata
    Type() ActivityType
    Category() Category  // gathering, crafting, hunting, travel, combat, dungeon, adventure

    // Validation
    ValidateStart(ctx context.Context, char Character, params ActivityParams, loc Location) error

    // Snapshot
    Snapshot(ctx context.Context, char Character, params ActivityParams) (Snapshot, error)

    // Rate compute — analytical, no simulation
    ComputeRate(ctx context.Context, snapshot Snapshot, params ActivityParams) (Rate, error)

    // Reveal (analytical yield)
    ComputeYield(ctx context.Context, contract Contract, snapshot Snapshot, elapsedMs int64) (Yield, error)

    // Threshold detection (envanter dolar mı, stamina biter mi, mob tükenir mi)
    DetectThreshold(ctx context.Context, contract Contract, snapshot Snapshot, elapsedMs int64) (*Threshold, error)

    // Pause / resume / cancel side effects
    OnPause(ctx context.Context, contract Contract, reason PauseReason) error
    OnResume(ctx context.Context, contract Contract) error
    OnCancel(ctx context.Context, contract Contract, mode CancelMode) (Settlement, error)
    OnComplete(ctx context.Context, contract Contract) (Settlement, error)
}
```

## 3. Ortak Tipler

```go
type ActivityType string  // "mining", "fishing", "hunting", "travel", "combat", "crafting", "dungeon_idle", "adventure"

type Category string      // "gathering", "crafting", "hunting", "travel", "combat", "dungeon", "adventure"

type Channel string       // "primary", "pet", "mount_travel"

type ActivityParams map[string]any   // aktivite-özel serialized

type Rate struct {
    YieldPerHour        float64
    XPPerHour           float64
    ProfessionXPPerHour float64
    MasteryXPPerHour    float64
    LootTableSampleRate float64   // loot rolls per hour (rare drop)
    ResourceDrainPerHour float64  // stamina, durability (v1 stamina kapalı)
}

type Yield struct {
    XPCharacter       int64
    XPProfession      map[uuid.UUID]int64  // profession_id → XP
    XPMastery         map[uuid.UUID]int64  // track_id → XP
    ItemDrops         []ItemDrop
    CurrencyDeltas    map[string]int64
    ReputationDeltas  map[uuid.UUID]int32
    DeathCount        int32
    XPDebt            int64
}

type Threshold struct {
    Kind         ThresholdKind   // inventory_full, stamina_zero, mob_dry, death_cascade, hp_zero_permastop
    HitAtMs      int64           // contract'ın başından itibaren ne kadar sonra
    Notes        string
}

type Settlement struct {
    Yield        Yield
    ContractEnd  ContractEnd     // paused, completed, cancelled
    ClosedAt     time.Time
}
```

## 4. Registry

```go
type Registry interface {
    Register(a Activity)
    Get(t ActivityType) (Activity, error)
    All() []Activity
}
```

Startup'ta registrasyon:
```go
reg := activityengine.NewRegistry()
reg.Register(mining.NewActivity(deps))
reg.Register(fishing.NewActivity(deps))
reg.Register(hunting.NewActivity(deps))
reg.Register(travel.NewActivity(deps))
reg.Register(combat.NewActivity(deps))
reg.Register(crafting.NewActivity(deps))
reg.Register(dungeon.NewActivity(deps))
reg.Register(adventure.NewActivity(deps))
```

Yeni aktivite tipi eklemek → yeni package + interface implement + registry'de register.

## 5. Params Schema

Her aktivitenin `ActivityParams` schema'sı JSON:

```yaml
# mining
{ "spot_id": "loc-iron-valley", "autoloot": true, "target_node_tier_min": 1 }

# hunting
{ "spot_id": "loc-whispering-forest", "mob_type_codes": ["ashen_wolf"], "autopotion": true, "stop_on_hp_below_pct": 30 }

# travel
{ "destination_location_id": "loc-emerald-coast", "method": "mount", "route_json": [...] }

# crafting
{ "recipe_id": "...", "iterations": 20, "station_location_id": "..." }

# dungeon_idle
{ "dungeon_template_id": "...", "difficulty": "normal", "party_id": null }
```

Loader/validator per activity tip.

## 6. Snapshot İçeriği

```go
type Snapshot struct {
    // Core
    CharacterID    uuid.UUID
    ContentVersion string
    Level          int32
    AscensionLevel int32

    // Race + Class + Spec
    RaceCode        string
    ClassCode       string
    SpecTier1Code   string
    SpecTier2Code   string

    // Stat derived (kilitli)
    DerivedStats    engine.Derived

    // Equipment ids (referenced instances)
    EquippedItemInstanceIDs map[EquipmentSlot]uuid.UUID

    // Skill rank
    SkillRanks      map[uuid.UUID]int16  // learned skills

    // Professions
    ProfessionLevels map[uuid.UUID]int16
    ActiveProfessions []uuid.UUID

    // Mastery levels (activity yield etkileyen)
    MasteryLevels   map[uuid.UUID]int16

    // Pet + Mount
    ActivePetInstanceID   *uuid.UUID
    ActiveMountInstanceID *uuid.UUID

    // Buffs (uzun-süreli food, elixir)
    ActiveBuffCodes []string

    // Seed (RNG parent)
    ContractSeed    Seed
}
```

Snapshot alındıktan sonra karakter değişse bile contract süresince bu değerler kullanılır.

**Trade-off**: karakter equip değiştirirse ve reveal olursa, snapshot güncelleneceğinden sonuç yeni equipment'la olur. Yani snapshot'ın "kesinliği" reveal aralığıyla sınırlı.

**Karar**: snapshot **contract lifetime** boyunca; her reveal'da re-snapshot fresh.

## 7. Contract vs Instance

- **Contract**: uzun-süreli AFK; kalıcı row (`activity.contract`).
- **Instance**: kısa-süreli session-scope (dungeon run, adventure run); `activity.instance`.

Bu doküman **contract**-based aktivitelere odaklanır. Instance-based (dungeon, adventure) kısmen contract benzeri ancak farklı lifecycle (bkz. [dungeon.md](dungeon.md), [adventure.md](../game-design/adventure.md)).

## 8. Cross-Character Fairness

- Snapshot'ı geç almak avantaj sağlar mı? Hayır — Rate compute snapshot güncel; yeni gear yeni snapshot'la yansır (yeni contract veya reveal sonrası).
- Snapshot alma zamanı = contract creation. Aynı loop'ta reveal → immediate; snapshot refresh.

## 9. Content Reference (Loaded Content Freezes at Snapshot)

`content_version` snapshot'a bağlı. Contract süresi içinde content patch olursa:
- Contract yaşamaya devam.
- Reveal yaparken snapshot content_version kullanılır (mekanik consistency).
- Yeni contract yeni content ile başlar.

## 10. Test Stratejisi

- Activity interface implementation her tip için.
- Registry lookup.
- Params schema validation.
- Snapshot immutability.
- Cross-activity contract kanalı çakışma.

## 11. Metrics
- `activity_started_total{type}`.
- `activity_paused_total{type,reason}`.
- `activity_completed_total{type}`.
- `activity_active_contracts_total{type}` gauge.
- `activity_channel_conflict_total`.

## 12. Cross-Reference
- [Lifecycle](lifecycle.md), [Offline Progress](offline-progress.md), [Rewards](rewards.md), [XP](xp.md)
- Aktivite dokümanları: [mining.md](mining.md), [fishing.md](fishing.md), [crafting.md](crafting.md), [hunting.md](hunting.md), [travel.md](travel.md), [combat.md](combat.md), [dungeon.md](dungeon.md)
- [database/activity.md](../database/activity.md)
- [systems/activity-engine.md](../systems/activity-engine.md)

## 13. Bekleyen Kararlar
- Adventure aktivitesinin bu klasörde ayrı doküman mı (v1 mixed — Phase 02'de detay).
- Cross-channel dependency (mount_travel primary'yi bloke etmesi) v1: allowed.
- Content patch mid-contract compensation policy.
