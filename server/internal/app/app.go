// Package app wires every domain module into the HTTP router. cmd/api and
// the in-process HTTP tests share it, so tests exercise production wiring.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/admin"
	"github.com/otter/server/internal/classskill"
	"github.com/otter/server/internal/clienterr"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/consumable"
	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/guild"
	httpapi "github.com/otter/server/internal/http"
	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/hunt"
	"github.com/otter/server/internal/idempotency"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/market"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/observability"
	"github.com/otter/server/internal/party"
	"github.com/otter/server/internal/privacy"
	"github.com/otter/server/internal/progression"
	"github.com/otter/server/internal/quest"
	"github.com/otter/server/internal/ratelimit"
	"github.com/otter/server/internal/skill"
	"github.com/otter/server/internal/social"
	"github.com/otter/server/internal/support"
	"github.com/otter/server/internal/vendor"
	"github.com/otter/server/internal/world"
	"github.com/otter/server/internal/worldevent"
	"github.com/otter/server/internal/wsgateway"
)

// Options are the platform pieces the router is built on.
type Options struct {
	Config    config.Config
	Pool      *pgxpool.Pool
	Redis     *redis.Client   // readiness probe only; may be nil in tests
	RateLimit ratelimit.Store // Redis store in prod
	Metrics   *observability.Metrics
	Logger    *slog.Logger
	Commit    string
}

// Build wires config → domain modules → chi router and loads the content
// pack. ctx bounds the startup work (catalog loads, content pack).
func Build(ctx context.Context, o Options) (http.Handler, error) {
	cfg, pool, rdb, metrics, logger := o.Config, o.Pool, o.Redis, o.Metrics, o.Logger
	rootCtx, CommitSHA := ctx, o.Commit
	if metrics == nil {
		metrics = observability.NewMetrics()
	}

	idemStore := idempotency.NewPgStore(pool)
	rlStore := o.RateLimit

	// --- Auth (DD-189) ---
	signer, err := buildSigner(cfg, logger)
	if err != nil {
		return nil, err
	}
	auth.SetKDFConcurrency(cfg.AuthHashConcurrency)
	authRepo := auth.NewPgRepo(pool)
	authSvc := auth.NewService(authRepo, signer, auth.ServiceConfig{
		Argon2Params:   auth.DefaultArgon2Params,
		AccessTTL:      cfg.AuthAccessTTL,
		RefreshTTL:     cfg.AuthRefreshTTL,
		SessionTTL:     cfg.AuthSessionTTL,
		Now:            time.Now,
		MinPasswordLen: cfg.AuthMinPasswordLen,
		PublicBaseURL:  cfg.PublicBaseURL,
	}, logger)
	if cfg.AuthRequireVerifiedEmail {
		logger.Info("character creation requires a verified e-mail")
	}
	logger.Info("auth wired", "access_ttl", cfg.AuthAccessTTL, "refresh_ttl", cfg.AuthRefreshTTL)

	// --- Character ---
	charRepo := character.NewPgRepo(pool)
	charSvc := character.NewService(charRepo, character.ServiceConfig{
		SlotCap:     cfg.CharacterSlotCap,
		DeleteGrace: cfg.CharacterDeleteGrace,
		Now:         time.Now,
	}, logger)
	if cfg.AuthRequireVerifiedEmail {
		charSvc.WithEmailVerification(authSvc.EmailVerified)
	}
	logger.Info("character wired", "slot_cap", cfg.CharacterSlotCap)

	// --- Economy (DD-135) ---
	currencies, err := economy.LoadCurrencies(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load currencies: %w", err)
	}
	txnSvc := economy.NewPgTxnService(pool, currencies, logger)
	logger.Info("economy wired", "currencies", len(currencies))

	// --- Item Engine (Sprint 8) ---
	itemCatalog, err := item.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load item catalog: %w", err)
	}
	itemRepo := item.NewPgRepo(pool)
	itemSvc := item.NewService(itemRepo, itemCatalog, logger)
	logger.Info("item wired", "templates", len(itemCatalog.All()))

	// --- Inventory (Sprint 9) — wired back into itemSvc as SlotAssigner ---
	invRepo := inventory.NewPgRepo(pool)
	invSvc := inventory.NewService(invRepo, logger)
	itemSvc.WithSlotAssigner(invSvc)
	logger.Info("inventory wired",
		"slot_cap", inventory.BagSlotCap,
		"weight_cap_g", inventory.WeightCapGrams,
	)

	// --- Equipment + Skill Engine (Sprint 11) ---
	// EquipmentService needs a small character read (level + class) — bridge
	// via SimpleCharFetcher so we don't couple inventory to character types.
	charFetch := inventory.SimpleCharFetcher(func(ctx context.Context, id uuid.UUID) (int16, string, error) {
		c, err := charSvc.GetOwnedInternal(ctx, id)
		if err != nil {
			return 0, "", err
		}
		return c.Level, c.ClassCode, nil
	})
	equipSvc := inventory.NewEquipmentService(invSvc, itemSvc, charFetch)
	// Wire sheet computer: character.Sheet now uses skill.ComputeSheet with equipped bonuses.
	sheetAdapter := skill.NewCharacterAdapter(equipSvc, logger)
	charSvc.WithSheetComputer(sheetAdapter)
	logger.Info("skill engine wired — character sheets use equipped stat_bonuses")

	// --- Vendor (Sprint 10) ---
	vendorRepo := vendor.NewPgRepo(pool)
	vendorSvc, err := vendor.NewService(rootCtx, vendorRepo, itemSvc, invSvc, txnSvc,
		vendor.DefaultServiceConfig(), logger)
	if err != nil {
		return nil, fmt.Errorf("vendor init: %w", err)
	}
	logger.Info("vendor wired")

	// --- Activity (DD-084) ---
	activityRepo := activity.NewPgRepo(pool)
	activitySvc := activity.NewService(activityRepo, charSvc, txnSvc,
		activity.DefaultServiceConfig(), logger).
		WithItemGranter(itemSvc) // Sprint 8: real item drops with slot assignment
	logger.Info("activity wired")

	// --- Combat (Sprint 13) ---
	// Mob catalog first — LoadCatalog will succeed even on an empty table
	// (fresh DB before contentpack runs) so we start with an empty catalog and
	// let the loader below refresh it.
	mobCat, err := combat.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load mob catalog: %w", err)
	}
	combatSvc := combat.NewService(mobCat, logger)
	combatHandler := combat.NewHandler(combatSvc, charSvc, logger)
	logger.Info("combat wired", "mobs", len(mobCat.All()))

	// --- Crafting (Sprint 17) — professions + recipes ---
	craftCat, err := crafting.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load crafting catalog: %w", err)
	}
	craftCatSvc := crafting.NewCatalogService(craftCat, logger)
	craftRepo := crafting.NewPgRepo(pool)
	craftSvc := crafting.NewService(
		craftRepo, craftCatSvc, invSvc, itemSvc, txnSvc, charSvc,
		crafting.DefaultServiceConfig(), logger,
	)
	craftHandler := crafting.NewHandler(craftSvc, logger)
	logger.Info("crafting wired",
		"professions", len(craftCat.AllProfessions()),
		"recipes", len(craftCat.AllRecipes()))

	// --- World / Travel (Sprint 16) — region hierarchy + travel state ---
	worldCat, err := world.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load world catalog: %w", err)
	}
	worldCatSvc := world.NewCatalogService(worldCat, logger)
	worldRepo := world.NewPgRepo(pool)
	worldSvc := world.NewService(worldRepo, worldCatSvc, charSvc, world.DefaultServiceConfig(), logger)
	worldHandler := world.NewHandler(worldSvc, worldCatSvc, logger)
	logger.Info("world wired", "regions", len(worldCat.All()), "starting", worldCat.StartingCode())

	// --- Dungeon (Sprint 15) — AFK run with boss + weekly lockout ---
	dungCat, err := dungeon.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load dungeon catalog: %w", err)
	}
	dungCatSvc := dungeon.NewCatalogService(dungCat, logger)
	dungRepo := dungeon.NewPgRepo(pool)
	dungSvc := dungeon.NewService(
		dungRepo, dungCatSvc, combatSvc, charSvc, charSvc, itemSvc, txnSvc,
		dungeon.DefaultServiceConfig(), logger,
		func(seed int64) combat.RNG { return activity.NewRNG(seed) },
	)
	dungHandler := dungeon.NewHandler(dungSvc, logger)
	logger.Info("dungeon wired", "dungeons", len(dungCat.All()))

	// --- Hunt (Sprint 14) — AFK combat encounter loop ---
	huntRepo := hunt.NewPgRepo(pool)
	huntSvc := hunt.NewService(
		huntRepo, combatSvc, charSvc, charSvc, itemSvc, txnSvc,
		hunt.DefaultServiceConfig(), logger,
		func(seed int64) combat.RNG { return activity.NewRNG(seed) },
	)
	huntHandler := hunt.NewHandler(huntSvc, logger)
	logger.Info("hunt wired")

	// --- WebSocket Gateway (Sprint 18) — presence + fanout ---
	// Hub lives for the process lifetime; adapter is what services import to
	// avoid a direct dependency on wsgateway.
	wsHub := wsgateway.NewHub(logger)
	wsPub := wsgateway.NewKindPayloadAdapter(wsHub)
	wsHandler := wsgateway.NewHandler(wsHub, authSvc, logger)
	metrics.GaugeFunc("otter_ws_clients", "Open WebSocket connections.",
		func() float64 { return float64(wsHub.ClientCount()) })
	huntHandler = huntHandler.WithPublisher(wsPub)
	dungHandler = dungHandler.WithPublisher(wsPub)
	logger.Info("ws gateway wired")

	// --- Notification (Sprint 23) — inbox + push queue (worker delivers) ---
	notifSvc := notification.NewService(pool, logger)
	notifHandler := notification.NewHandler(notifSvc)

	// --- Progression (Sprint 24) — mastery, achievements, titles, reputation ---
	progCat, err := progression.LoadCatalog(rootCtx, pool)
	if err != nil {
		return nil, fmt.Errorf("load progression catalog: %w", err)
	}
	progSvc := progression.NewService(pool, progCat, charSvc, txnSvc, logger).
		WithPublisher(wsPub).WithNotifier(notifSvc)
	progHandler := progression.NewHandler(progSvc)
	// Quests tick from the same event stream as mastery/achievements (DD-221).
	questSvc := quest.NewService(pool, charSvc, txnSvc, charSvc, logger).WithPublisher(wsPub)
	recorder := progression.Recorders(progSvc, questSvc)
	activitySvc.WithRecorder(recorder).WithGatherBonus(progSvc).WithNotifier(notifSvc)
	// Mastery bonuses feed the sheet's modifier stack (DD-220).
	sheetAdapter.WithModifierSource(func(ctx context.Context, id uuid.UUID, _ character.Sheet) ([]skill.Modifier, error) {
		effects, err := progSvc.ActiveBonuses(ctx, id)
		if err != nil {
			return nil, err
		}
		mods := make([]skill.Modifier, len(effects))
		for i, e := range effects {
			mods[i] = skill.Modifier{Kind: skill.ModifierKind(e.Kind), Target: skill.Target(e.Target), Amount: e.Amount}
		}
		return mods, nil
	})
	// Class skills (DD-222): passives → sheet stack, actives → combat snapshots.
	classSkills := classskill.NewService(pool)
	sheetAdapter.WithModifierSource(classSkills.PassiveModifiers)
	huntSvc.WithRecorder(recorder).WithHPWriter(charSvc).WithSkills(classSkills.ActiveSkills).WithNotifier(notifSvc)
	dungSvc.WithRecorder(recorder).WithHPWriter(charSvc).WithSkills(classSkills.ActiveSkills).WithNotifier(notifSvc)
	craftSvc.WithRecorder(recorder)

	// --- World Events (Sprint 25) — worker schedules/closes; API enlists ---
	worldEventHandler := worldevent.NewHandler(worldevent.NewService(pool, charSvc, txnSvc,
		func() string { return worldCatSvc.Catalog().StartingCode() }, logger))

	// --- Party + Guild (Sprint 20) ---
	partyHandler := party.NewHandler(party.NewService(pool, charSvc, logger).WithPublisher(wsPub).WithNotifier(notifSvc))
	guildHandler := guild.NewHandler(guild.NewService(pool, charSvc, txnSvc, logger).WithPublisher(wsPub).WithNotifier(notifSvc))
	logger.Info("party + guild wired")

	// --- Social (Sprint 21) — friends, block, mail, chat ---
	socialHandler := social.NewHandler(social.NewService(pool, charSvc, txnSvc, invSvc, logger).
		WithPublisher(wsPub).WithNotifier(notifSvc))

	// --- Marketplace (Sprint 22) ---
	marketSvc := market.NewService(pool, charSvc, txnSvc, itemSvc, invSvc, logger).
		WithPublisher(wsPub).WithNotifier(notifSvc)
	marketHandler := market.NewHandler(marketSvc)

	// --- Admin / Ops (Sprint 26, DD-213) ---
	supportSvc := support.NewService(pool, notifSvc, logger)
	clientErrSvc := clienterr.NewService(pool, metrics.ObserveClientError)
	adminHandler := admin.NewHandler(admin.NewService(pool, txnSvc, logger), marketSvc).
		WithSupport(supportSvc).WithClientErrors(clientErrSvc)

	// --- Content Pack Loader (Sprint 12) ---
	// YAML → DB upsert → refresh in-process catalogs. Runs after all services
	// that hold a catalog are wired so their Reload closures can hot-swap.
	contentLoader := contentpack.New(cfg.ContentPath, pool, logger)
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return itemSvc.ReloadCatalog(ctx, pool)
	}))
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return vendorSvc.ReloadCatalog(ctx)
	}))
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return combatSvc.ReloadCatalog(ctx, pool)
	}))
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return dungCatSvc.ReloadCatalog(ctx, pool)
	}))
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return worldCatSvc.ReloadCatalog(ctx, pool)
	}))
	contentLoader.Register(reloadFn(func(ctx context.Context) error {
		return craftCatSvc.ReloadCatalog(ctx, pool)
	}))
	contentLoader.Register(questSvc)
	contentLoader.Register(classSkills)
	if _, err := contentLoader.Load(rootCtx); err != nil {
		// Missing directory or malformed YAML must not silently swallow — the
		// migration seed is the fallback baseline; a failed loader still yields
		// a working server but no YAML overrides applied.
		logger.Warn("content pack load failed — falling back to migration seed",
			"err", err.Error(), "path", cfg.ContentPath)
		// Quests and class skills have no migration seed: serve whatever the
		// last good load stored.
		if err := questSvc.Reload(rootCtx); err != nil {
			logger.Error("quest catalog load failed", "err", err.Error())
		}
		if err := classSkills.Reload(rootCtx); err != nil {
			logger.Error("class skill catalog load failed", "err", err.Error())
		}
	}

	// HTTP
	trusted, err := appmw.ParseTrustedProxies(cfg.TrustedProxies)
	if err != nil {
		return nil, err
	}
	router := httpapi.NewRouter(httpapi.Deps{
		Logger:            logger,
		DB:                pool,
		Redis:             rdb,
		CommitSHA:         CommitSHA,
		Env:               cfg.Env,
		IdempotencyStore:  idemStore,
		RateLimitStore:    rlStore,
		AuthService:       authSvc,
		CharacterService:  charSvc,
		TxnService:        txnSvc,
		ActivityService:   activitySvc,
		ItemService:       itemSvc,
		InventoryService:  invSvc,
		EquipmentService:  equipSvc,
		VendorService:     vendorSvc,
		ContentLoader:     contentLoader,
		CombatHandler:     combatHandler,
		HuntHandler:       huntHandler,
		DungeonHandler:    dungHandler,
		WorldHandler:      worldHandler,
		CraftingHandler:   craftHandler,
		WSHandler:         wsHandler,
		PartyHandler:      partyHandler,
		GuildHandler:      guildHandler,
		SocialHandler:     socialHandler,
		MarketHandler:     marketHandler,
		NotifHandler:      notifHandler,
		ProgHandler:       progHandler,
		QuestHandler:      quest.NewHandler(questSvc),
		ClassSkills:       classSkills,
		WorldEventHandler: worldEventHandler,
		ConsumableHandler: consumable.NewHandler(consumable.NewService(pool, itemSvc, charSvc, logger)),
		AdminHandler:      adminHandler,
		SupportHandler:    support.NewHandler(supportSvc),
		ClientErrHandler:  clienterr.NewHandler(clientErrSvc),
		PrivacyHandler: privacy.NewHandler(pool, authSvc, cfg.ContentPath, privacy.Legal{
			Operator: cfg.LegalOperator, SupportEmail: cfg.SupportEmail, TermsVersion: auth.DefaultTermsVersion,
		}, logger),
		Metrics:        metrics,
		CORSOrigins:    cfg.CORSAllowedOrigins,
		TrustedProxies: trusted,
	})
	return router, nil
}

// reloadFn adapts a plain closure to contentpack.Reloader.
type reloadFn func(ctx context.Context) error

func (f reloadFn) Reload(ctx context.Context) error { return f(ctx) }

func buildSigner(cfg config.Config, logger *slog.Logger) (*auth.Signer, error) {
	if cfg.JWTPrivateKeySeed == "" {
		if cfg.Env != "dev" {
			return nil, errors.New("JWT_PRIVATE_KEY_SEED required outside dev")
		}
		seed, err := auth.GenerateSeed()
		if err != nil {
			return nil, fmt.Errorf("generate ephemeral seed: %w", err)
		}
		// Never log the seed itself: it is the token-signing private key.
		logger.Warn("JWT seed not configured — generated ephemeral key (all tokens invalidate on restart)")
		return auth.NewSigner(seed, cfg.JWTIssuer, cfg.AuthAccessTTL)
	}
	seed, err := auth.DecodeSeed(cfg.JWTPrivateKeySeed)
	if err != nil {
		return nil, fmt.Errorf("decode JWT seed: %w", err)
	}
	return auth.NewSigner(seed, cfg.JWTIssuer, cfg.AuthAccessTTL)
}
