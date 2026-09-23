# deploy/

Dev docker-compose + prod Dockerfile.

## Dev

```bash
cd deploy
docker compose up -d          # postgres + redis + migrator + api
docker compose logs -f api
docker compose down            # sadece durdur (volume korunur)
docker compose down -v         # volume dahil sil
```

## Volume

`deploy/.data/postgres/` ve `deploy/.data/redis/` git'ten hariç (`.gitignore`).

## Prod

`api.Dockerfile` distroless nonroot; K8s deploy Faz 01 sprint sonlarında (infra/ dokümantasyon halinde).
