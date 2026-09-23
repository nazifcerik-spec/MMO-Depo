// Package http wires the chi router with middleware order matching DD-186 (defense
// in depth): request id → logger → recoverer → cors → timeout → route.
package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/admin"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/consumable"
	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/domain/health"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/guild"
	"github.com/otter/server/internal/http/apierr"
	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/hunt"
	"github.com/otter/server/internal/idempotency"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/market"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/observability"
	"github.com/otter/server/internal/party"
	"github.com/otter/server/internal/progression"
	"github.com/otter/server/internal/ratelimit"
	"github.com/otter/server/internal/social"
	"github.com/otter/server/internal/vendor"
	"github.com/otter/server/internal/world"
	"github.com/otter/server/internal/worldevent"
	"github.com/otter/server/internal/wsgateway"
)

// Deps is the dependency bag passed to the router.
type Deps struct {
	Logger            *slog.Logger
	DB                *pgxpool.Pool
	Redis             *redis.Client
	CommitSHA         string
	Env               string
	IdempotencyStore  idempotency.Store
	RateLimitStore    ratelimit.Store
	AuthService       *auth.Service
	CharacterService  *character.Service
	TxnService        economy.TxnService
	ActivityService   *activity.Service
	ItemService       *item.Service
	InventoryService  *inventory.Service
	EquipmentService  *inventory.EquipmentService
	VendorService     *vendor.Service
	ContentLoader     *contentpack.Loader
	CombatHandler     *combat.Handler
	HuntHandler       *hunt.Handler
	DungeonHandler    *dungeon.Handler
	WorldHandler      *world.Handler
	CraftingHandler   *crafting.Handler
	WSHandler         *wsgateway.Handler
	PartyHandler      *party.Handler
	GuildHandler      *guild.Handler
	SocialHandler     *social.Handler
	MarketHandler     *market.Handler
	NotifHandler      *notification.Handler
	ProgHandler       *progression.Handler
	WorldEventHandler *worldevent.Handler
	AdminHandler      *admin.Handler
	Metrics           *observability.Metrics
	CORSOrigins       []string
	ConsumableHandler *consumable.Handler
}

// NewRouter builds the top-level chi.Mux. Every non-liveness route sits under /v1.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(observability.HTTP(d.Metrics))
	r.Use(appmw.CORS(d.CORSOrigins, d.Env == "dev"))
	r.Use(chimw.RequestID)
	r.Use(appmw.Logger(d.Logger))
	r.Use(appmw.Recoverer(d.Logger))
	r.Use(chimw.Timeout(30 * time.Second))

	healthSvc := health.NewService(d.DB, d.Redis)
	r.Get("/live", health.Live)
	r.Get("/ready", healthSvc.Ready)

	// Dev console — single-page HTML tool for exercising the API. Only in dev.
	if d.Env == "dev" {
		r.Get("/dev/console", devConsoleHandler)
	}

	authHandler := auth.NewHandler(d.AuthService, d.Logger)
	charHandler := character.NewHandler(d.CharacterService, d.Logger)
	econHandler := economy.NewHandler(d.TxnService, d.Env == "dev", d.Logger)
	activityHandler := activity.NewHandler(d.ActivityService, d.Logger)
	itemHandler := item.NewHandler(d.ItemService, d.Logger)
	invHandler := inventory.NewHandler(d.InventoryService, d.Logger)
	equipHandler := inventory.NewEquipmentHandler(d.EquipmentService, d.Logger)
	vendorHandler := vendor.NewHandler(d.VendorService, d.CharacterService, d.Logger)

	// Links from recovery e-mails (DD-219); outside /v1 like the dev console.
	r.With(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
		Endpoint: "GET /verify-email", PerIP: ptrBucket(ratelimit.PerMinute(20)),
	}, d.Logger)).Get("/verify-email", authHandler.VerifyEmailPage)
	r.Get("/reset-password", authHandler.ResetPasswordPage)

	// Admin / ops surface (DD-213) — outside /v1, own bearer sessions.
	if h := d.AdminHandler; h != nil {
		h.Routes(r,
			appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "admin.api",
				PerIP:    ptrBucket(ratelimit.PerMinute(300)),
			}, d.Logger),
			appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /internal/admin/login",
				PerIP:    ptrBucket(ratelimit.PerMinute(5)),
			}, d.Logger))
	}

	r.Route("/v1", func(r chi.Router) {
		r.Get("/version", func(w http.ResponseWriter, req *http.Request) {
			apierr.WriteJSON(w, http.StatusOK, map[string]string{"commit": d.CommitSHA})
		})

		// --- Auth endpoints ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /v1/auth/register",
				PerIP:    ptrBucket(ratelimit.PerMinute(3)),
			}, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.Post("/auth/register", authHandler.Register)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /v1/auth/login",
				PerIP:    ptrBucket(ratelimit.PerMinute(10)),
			}, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.Post("/auth/login", authHandler.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /v1/auth/refresh",
				PerIP:    ptrBucket(ratelimit.PerMinute(60)),
			}, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.Post("/auth/refresh", authHandler.Refresh)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:   "POST /v1/auth/logout",
				PerIP:      ptrBucket(ratelimit.PerMinute(30)),
				PerAccount: ptrBucket(ratelimit.PerMinute(30)),
			}, d.Logger))
			r.Post("/auth/logout", authHandler.Logout)
		})

		// --- Account recovery (DD-219) ---
		recovery := func(endpoint string, perIP ratelimit.Bucket, h http.HandlerFunc) {
			r.With(
				appmw.RateLimit(d.RateLimitStore, appmw.Policy{Endpoint: endpoint, PerIP: ptrBucket(perIP)}, d.Logger),
				appmw.Idempotency(d.IdempotencyStore, d.Logger),
			).Post(endpoint[len("POST /v1"):], h)
		}
		recovery("POST /v1/auth/password/reset-request", ratelimit.PerHour(5), authHandler.RequestPasswordReset)
		recovery("POST /v1/auth/password/reset-confirm", ratelimit.PerHour(10), authHandler.ConfirmPasswordReset)
		recovery("POST /v1/auth/verify-email", ratelimit.PerMinute(20), authHandler.VerifyEmail)
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.With(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /v1/auth/resend-verification", PerAccount: ptrBucket(ratelimit.PerHour(3)),
			}, d.Logger)).Post("/auth/resend-verification", authHandler.ResendVerification)
			r.With(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint: "POST /v1/auth/password/change", PerAccount: ptrBucket(ratelimit.PerHour(5)),
			}, d.Logger)).Post("/auth/password/change", authHandler.ChangePassword)
		})

		// --- Character endpoints + wallet read ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:   "character.rw",
				PerIP:      ptrBucket(ratelimit.PerMinute(60)),
				PerAccount: ptrBucket(ratelimit.PerMinute(120)),
			}, d.Logger))

			r.Get("/character", charHandler.List)
			r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
				Post("/character", charHandler.Create)

			r.Route("/character/{id}", func(r chi.Router) {
				// Every route below acts on {id}: enforce ownership once, here.
				r.Use(appmw.RequireCharacterParam(d.CharacterService, "id", d.Logger))
				r.Get("/sheet", charHandler.Sheet)
				r.Get("/wallet", econHandler.Wallet)
				r.Post("/enter", charHandler.Enter)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/delete-request", charHandler.RequestDelete)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/delete-cancel", charHandler.CancelDelete)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/rename", charHandler.Rename)

				// --- Activity (character-scoped) ---
				r.Get("/activity/active", activityHandler.Active)
				r.Get("/activity/history", activityHandler.History)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/activity", activityHandler.Start)

				// --- Items (character-scoped, raw list) ---
				r.Get("/items", itemHandler.Items)
				if h := d.ConsumableHandler; h != nil {
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/items/{item_id}/use", h.Use)
				}

				// --- Buybacks (character-scoped) ---
				r.Get("/buybacks", vendorHandler.ListBuybacks)

				// --- Inventory (character-scoped) ---
				r.Get("/inventory", invHandler.View)
				r.Route("/inventory/{item_id}", func(r chi.Router) {
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/move", invHandler.Move)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/split", invHandler.Split)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Delete("/", invHandler.Delete)
				})

				// --- Equipment (Sprint 11) ---
				r.Get("/equipment", equipHandler.View)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/equip/{item_id}", equipHandler.Equip)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/unequip/{slot}", equipHandler.Unequip)

				// --- Hunt (Sprint 14) — AFK combat encounter contracts ---
				if d.HuntHandler != nil {
					r.Get("/hunt/active", d.HuntHandler.Active)
					r.Get("/hunt/history", d.HuntHandler.History)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/hunt", d.HuntHandler.Start)
				}

				// --- Dungeon (Sprint 15) — AFK run with boss + weekly lockout ---
				if d.DungeonHandler != nil {
					r.Get("/dungeon/active", d.DungeonHandler.Active)
					r.Get("/dungeon/history", d.DungeonHandler.History)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/dungeon", d.DungeonHandler.Start)
				}

				// --- World / Travel (Sprint 16) ---
				if d.WorldHandler != nil {
					r.Get("/location", d.WorldHandler.Location)
					r.Get("/travel/history", d.WorldHandler.History)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/travel", d.WorldHandler.Start)
				}

				// --- Crafting (Sprint 17) — character-scoped ---
				if d.CraftingHandler != nil {
					r.Get("/professions", d.CraftingHandler.List)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/professions/learn", d.CraftingHandler.Learn)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
						Post("/craft", d.CraftingHandler.Craft)
				}

				// --- Party (Sprint 20) ---
				if h := d.PartyHandler; h != nil {
					idem := appmw.Idempotency(d.IdempotencyStore, d.Logger)
					r.Get("/party", h.Get)
					r.Get("/party/invites", h.Invites)
					r.With(idem).Post("/party", h.Create)
					r.With(idem).Post("/party/invite", h.Invite)
					r.With(idem).Post("/party/invites/{invite_id}/accept", h.Accept)
					r.With(idem).Post("/party/invites/{invite_id}/decline", h.Decline)
					r.With(idem).Post("/party/leave", h.Leave)
					r.With(idem).Post("/party/kick", h.Kick)
					r.With(idem).Post("/party/promote", h.Promote)
					r.With(idem).Post("/party/disband", h.Disband)
				}

				// --- Guild (Sprint 20) ---
				if h := d.GuildHandler; h != nil {
					idem := appmw.Idempotency(d.IdempotencyStore, d.Logger)
					r.Get("/guild", h.Mine)
					r.Get("/guild/invites", h.Invites)
					r.With(idem).Post("/guild", h.Create)
					r.With(idem).Post("/guild/invite", h.Invite)
					r.With(idem).Post("/guild/invites/{invite_id}/accept", h.Accept)
					r.With(idem).Post("/guild/invites/{invite_id}/decline", h.Decline)
					r.With(idem).Post("/guild/leave", h.Leave)
					r.With(idem).Post("/guild/kick", h.Kick)
					r.With(idem).Post("/guild/rank", h.SetRank)
					r.With(idem).Post("/guild/transfer", h.Transfer)
					r.With(idem).Post("/guild/motd", h.MOTD)
					r.With(idem).Post("/guild/disband", h.Disband)
					r.With(idem).Post("/guild/disband/cancel", h.CancelDisband)
				}

				// --- Social (Sprint 21) ---
				if h := d.SocialHandler; h != nil {
					idem := appmw.Idempotency(d.IdempotencyStore, d.Logger)
					r.Get("/friends", h.Friends())
					r.Get("/friends/requests", h.FriendRequests())
					r.With(idem).Post("/friends/request", h.RequestFriend())
					r.With(idem).Post("/friends/requests/{request_id}/accept", h.AnswerFriend(true))
					r.With(idem).Post("/friends/requests/{request_id}/decline", h.AnswerFriend(false))
					r.With(idem).Post("/friends/remove", h.RemoveFriend())
					r.Get("/blocks", h.Blocks())
					r.With(idem).Post("/blocks", h.Block())
					r.With(idem).Post("/blocks/remove", h.Unblock())
					r.Get("/mail", h.Inbox())
					r.Get("/mail/{mail_id}", h.ReadMail())
					r.With(idem).Post("/mail", h.SendMail())
					r.With(idem).Post("/mail/{mail_id}/claim", h.ClaimMail())
					r.With(idem).Post("/mail/{mail_id}/delete", h.DeleteMail())
					r.Get("/chat/{channel}", h.ChatHistory())
					r.With(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
						Endpoint:   "chat.send",
						PerAccount: ptrBucket(ratelimit.PerMinute(20)),
					}, d.Logger)).Post("/chat", h.SendChat())
				}

				// --- World events (Sprint 25) ---
				if h := d.WorldEventHandler; h != nil {
					r.Get("/world-events", h.List)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).Post("/world-events/{spawn_id}/attend", h.Attend)
				}

				// --- Progression (Sprint 24) ---
				if h := d.ProgHandler; h != nil {
					r.Get("/progression", h.Summary)
					r.Get("/achievements", h.Achievements)
					r.Get("/titles", h.Titles)
					r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).Post("/titles/active", h.SetActiveTitle)
				}

				// --- Marketplace (Sprint 22) ---
				if h := d.MarketHandler; h != nil {
					idem := appmw.Idempotency(d.IdempotencyStore, d.Logger)
					r.Get("/market/listings", h.Mine)
					r.With(idem).Post("/market/sell", h.Sell)
					r.With(idem).Post("/market/buy/{listing_id}", h.Buy)
					r.With(idem).Post("/market/cancel/{listing_id}", h.Cancel)
				}
			})

			if d.GuildHandler != nil {
				r.Get("/guilds/{guild_id}", d.GuildHandler.Public)
			}
			if d.MarketHandler != nil {
				r.Get("/market/listings", d.MarketHandler.Search)
				r.Get("/market/history/{template_code}", d.MarketHandler.History)
			}

			// --- Hunt contract-scoped (reveal/stop by contract id) ---
			if d.HuntHandler != nil {
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/hunt/{contract_id}/reveal", d.HuntHandler.Reveal)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/hunt/{contract_id}/stop", d.HuntHandler.Stop)
			}

			// --- Public mob catalog (client hunt picker) ---
			if d.CombatHandler != nil {
				r.Get("/mobs", d.CombatHandler.Mobs)
			}

			// --- Dungeon contract-scoped + public catalog ---
			if d.DungeonHandler != nil {
				r.Get("/dungeon/catalog", d.DungeonHandler.Catalog)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/dungeon/{contract_id}/reveal", d.DungeonHandler.Reveal)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/dungeon/{contract_id}/stop", d.DungeonHandler.Stop)
			}

			// --- World / Travel — public catalog + travel-scoped controls ---
			if d.WorldHandler != nil {
				r.Get("/world/regions", d.WorldHandler.Catalog)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/travel/{travel_id}/complete", d.WorldHandler.Complete)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/travel/{travel_id}/cancel", d.WorldHandler.Cancel)
			}

			// --- Crafting — public catalog ---
			if d.CraftingHandler != nil {
				r.Get("/professions", d.CraftingHandler.Professions)
				r.Get("/recipes", d.CraftingHandler.Recipes)
			}

		})

		// --- WebSocket upgrade (Sprint 18) ---
		// Placed OUTSIDE the RequireAuth group because browsers cannot send
		// custom headers on the WS handshake — Upgrade does its own token
		// verification via the `?token=` query param.
		if d.WSHandler != nil {
			r.Get("/ws", d.WSHandler.Upgrade)
		}

		// --- Notifications (Sprint 23) — account-scoped ---
		if h := d.NotifHandler; h != nil {
			r.Group(func(r chi.Router) {
				r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
				r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
					Endpoint:   "notification.rw",
					PerIP:      ptrBucket(ratelimit.PerMinute(60)),
					PerAccount: ptrBucket(ratelimit.PerMinute(120)),
				}, d.Logger))
				idem := appmw.Idempotency(d.IdempotencyStore, d.Logger)
				r.With(idem).Post("/device/push-token", h.RegisterDevice)
				r.Delete("/device/push-token/{device_id}", h.UnregisterDevice)
				r.Get("/notifications", h.Inbox)
				r.Post("/notifications/read-all", h.MarkRead)
				r.Post("/notifications/{notification_id}/read", h.MarkRead)
				r.Get("/notification/preference", h.GetPreferences)
				r.Put("/notification/preference", h.PutPreferences)
				r.Post("/notification/preference/global-mute", h.Mute)
				r.Delete("/notification/preference/global-mute", h.Unmute)
				if we := d.WorldEventHandler; we != nil {
					r.Get("/events/seasonal", we.Seasonal)
					r.With(idem).Post("/events/seasonal/{code}/claim", we.ClaimDaily)
				}
			})
		}

		// --- Presence + dev broadcast (Bearer-auth REST paths) ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:   "ws.presence",
				PerIP:      ptrBucket(ratelimit.PerMinute(30)),
				PerAccount: ptrBucket(ratelimit.PerMinute(60)),
			}, d.Logger))
			if d.WSHandler != nil {
				r.Get("/ws/presence", d.WSHandler.Presence)
				if d.Env == "dev" {
					r.Post("/dev/ws/broadcast", d.WSHandler.DevBroadcast)
				}
			}
		})

		// --- Activity + vendor contract-scoped (existing authed group) ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:   "character.rw",
				PerIP:      ptrBucket(ratelimit.PerMinute(60)),
				PerAccount: ptrBucket(ratelimit.PerMinute(120)),
			}, d.Logger))

			// --- Activity (contract-scoped) ---
			r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
				Post("/activity/{contract_id}/reveal", activityHandler.Reveal)
			r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
				Post("/activity/{contract_id}/stop", activityHandler.Stop)

			// --- Vendor (Sprint 10) ---
			r.Get("/vendor/{code}", vendorHandler.Info)
			r.Route("/vendor/{code}", func(r chi.Router) {
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/buy", vendorHandler.Buy)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/sell", vendorHandler.Sell)
				r.With(appmw.Idempotency(d.IdempotencyStore, d.Logger)).
					Post("/buyback/{buyback_id}", vendorHandler.Buyback)
			})
		})

		// --- Dev endpoints (404 outside ENV=dev) ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireAuth(d.AuthService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:   "POST /v1/dev/*",
				PerIP:      ptrBucket(ratelimit.PerMinute(30)),
				PerAccount: ptrBucket(ratelimit.PerMinute(30)),
			}, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.Post("/dev/mint", econHandler.Mint)
			if d.ContentLoader != nil {
				r.Post("/dev/content/reload", contentReloadHandler(d.ContentLoader))
			}
			if d.CombatHandler != nil {
				r.Post("/dev/combat/simulate", d.CombatHandler.Simulate)
			}
		})

		// --- Smoke echo (idempotency + rate limit + optional auth/character) ---
		r.Group(func(r chi.Router) {
			r.Use(appmw.OptionalAuth(d.AuthService, d.Logger))
			r.Use(appmw.OptionalCharacter(d.CharacterService, d.Logger))
			r.Use(appmw.RateLimit(d.RateLimitStore, appmw.Policy{
				Endpoint:     "POST /v1/echo",
				PerIP:        ptrBucket(ratelimit.PerMinute(30)),
				PerAccount:   ptrBucket(ratelimit.PerMinute(60)),
				PerCharacter: ptrBucket(ratelimit.PerMinute(90)),
			}, d.Logger))
			r.Use(appmw.Idempotency(d.IdempotencyStore, d.Logger))
			r.Post("/echo", echoHandler)
		})
	})

	return r
}

func ptrBucket(b ratelimit.Bucket) *ratelimit.Bucket { return &b }

// contentReloadHandler is the dev-only endpoint that re-reads YAML and
// refreshes item + vendor catalogs. Idempotent — every call re-scans files.
func contentReloadHandler(loader *contentpack.Loader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := loader.Load(r.Context())
		if err != nil {
			apierr.Write(w, r, http.StatusInternalServerError, "E_CONTENT_RELOAD", err.Error())
			return
		}
		apierr.WriteJSON(w, http.StatusOK, map[string]any{
			"reloaded":     true,
			"items_seen":   res.ItemsSeen,
			"vendors_seen": res.VendorsSeen,
			"source_dir":   res.SourceDir,
		})
	}
}

// echoHandler is the smoke endpoint. Reads JSON body, echoes it with a server
// timestamp so replayed calls (idempotency HIT) can be distinguished by their
// identical response.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, appmw.MaxBodyBytes+1))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_BODY_READ", "failed to read body")
		return
	}
	var payload any
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
			return
		}
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"echo":        payload,
		"server_time": time.Now().UTC().Format(time.RFC3339Nano),
	})
}
