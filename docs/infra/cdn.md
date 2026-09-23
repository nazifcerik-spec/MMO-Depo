# CDN (Static Asset + Edge Cache)

## 1. Amaç
Client-side asset (Flutter web build, icon sprite, cosmetic art, content pack manifest) global CDN edge cache; WAF; DDoS mitigation.

## 2. Sağlayıcı

- **v1**: **Cloudflare** (öneri) — free/paid tier, WAF, DDoS, R2 origin ucuz.
- Alternatif: AWS CloudFront, Fastly.

## 3. Kapsam

| Asset Türü | Cache Stratejisi |
|---|---|
| Flutter web bundle (index.html, main.dart.js) | Immutable hash filename → 1y cache; index.html 5m |
| Icon sprite sheet (item, class) | Immutable hash → 1y cache |
| Cosmetic art (avatar, tabard, aura) | Immutable → 1y cache |
| Content pack manifest (`content/v42/manifest.json`) | 5m cache; ETag revalidate |
| Content pack blob (`content/v42/items.yaml`) | Immutable version path → 1y |
| Locale bundle | Immutable hash → 1y |
| Landing page | 5m cache; SSR/static |
| API response | **CDN cache YOK** (dynamic) — origin passthrough |

## 4. Cache Rules

Cloudflare Page Rules:
```
game.cdn.io/assets/*     → Cache Everything, TTL 1y
game.cdn.io/content/*    → Cache Everything, TTL 1y
game.cdn.io/manifest/*   → Cache 5m, revalidate
api.game.io/*            → Bypass cache
admin.game.internal/*    → Bypass cache
```

Immutable path: `/assets/v1.4.0-abc123/main.js` → hash pattern; new deploy = new URL → auto invalidation.

## 5. Purge

- **Auto**: immutable hash yeni deploy → yeni URL, eski otomatik expire.
- **Manual purge**: rare (bug fix stale content) — API `POST /zones/{id}/purge_cache`.

## 6. Origin

- Static bucket S3 / Cloudflare R2.
- Origin auth: Cloudflare access rule (signed URL Phase 02 hassas asset için).

## 7. WAF (Web Application Firewall)

Cloudflare WAF ruleset:
- OWASP core rules (SQL injection, XSS, RFI).
- Bot management (Phase 02 challenge suspect).
- Rate limit edge (100 req/dk/IP baseline).
- Country block (compliance / sanctions).

## 8. DDoS Mitigation

- Cloudflare default L3/L4/L7 DDoS mitigation (built-in).
- Attack alerts → ops Slack + PagerDuty (high sustained).

## 9. TLS at Edge

- Full (strict) mode — origin TLS zorunlu.
- HSTS enforce.
- TLS 1.3 modern.

## 10. Client Asset Loading

Client:
- Bootstrap → CDN manifest fetch (5m cache).
- Manifest version diff → download new bundle.
- Locale + content pack lazy load on-demand.
- Retry on 5xx (exp backoff 2 attempts).

## 11. Bandwidth

T3 bandwidth split (approx):
- API (dynamic): ~130 Mbps peak.
- WS: 100-200 Mbps peak (chat + events).
- CDN offload (static): 300-800 Mbps peak (asset first load).

Total: ~1 Gbps peak T3; T4 3-4 Gbps.

## 12. Cost

Cloudflare Pro ~$20/ay + R2 storage $0.015/GB + egress free.
Alternate CloudFront: pay per GB egress ~$0.085/GB.

## 13. Metrikler
- `cdn_hit_ratio` (Cloudflare analytics).
- `cdn_bandwidth_bytes` (region breakdown).
- `cdn_origin_pull_requests`.
- `waf_block_total{rule}`.

## 14. Cross-Reference
- [reverse-proxy.md](reverse-proxy.md), [ssl.md](ssl.md)
- [ui/design-system.md](../ui/design-system.md) — asset kaynağı

## 15. Bekleyen Kararlar
- CDN sağlayıcı final (Cloudflare vs CloudFront) ops karar.
- Signed URL hassas asset Phase 02.
- Content pack CDN vs direct backend serve — v1: CDN.
- Regional edge cache (Cloudflare all-region default) Phase 02 pricing.
