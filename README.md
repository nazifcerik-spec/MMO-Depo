# Otter

Persistent Idle MMORPG. Server-authoritative + AFK-first + ledger-based ekonomi.

## Yapı

```
docs/       Faz 00 tasarım dokümanları (415 doc, 212 DD).
server/     Go application server (Faz 01+).
deploy/     docker-compose + Dockerfile.
scripts/    Dev/ops yardımcı script'ler.
```

## Hızlı Başlangıç (Dev)

Gerekli: Docker + Docker Compose, Go 1.22+.

```bash
# 1. Dev environment aç (Postgres 16 + Redis 7 + migrator + API)
cd deploy
docker compose up -d

# 2. Health check
curl http://localhost:8080/live
# → {"status":"ok"}

curl http://localhost:8080/ready
# → {"status":"ok","db":"ok","redis":"ok"}

# 3. Log takibi
docker compose logs -f api
```

Server-only geliştirme (Docker'sız):

```bash
cd server
cp ../.env.example .env
# .env düzenle: DATABASE_URL, REDIS_ADDR
make dev
```

## Rehber

Faz 00 dokümantasyonu: [docs/GDD.md](docs/GDD.md)
Design Decisions: [docs/decisions/design-decisions.md](docs/decisions/design-decisions.md)
Claude Code oturumları için: [CLAUDE.md](CLAUDE.md)

## Uygulama Faz Haritası

- **Faz 00 (bitti):** GDD — 415 doc, 212 karar.
- **Faz 01 (aktif):** Backend + Client implementation.
- **Faz 02+:** Auction, Guild advanced, Store, PvP.
