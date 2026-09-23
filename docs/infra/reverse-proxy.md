# Reverse Proxy (LB + TLS Termination + WS)

## 1. Amaç
Public trafiğin app node'lara dağıtımı. TLS termination + WebSocket proxy + basit routing + rate limit (edge) + graceful drain.

## 2. Seçim

- **v1**: **HAProxy 2.9+** (mature LB, WS destek, dinamik config).
- Alternatif: nginx, Envoy (Phase 02 mesh için).
- Cloud LB (AWS ALB / GCP LB) sağlayıcıya göre kabul.

## 3. Görev Dağılımı

| Katman | Görev |
|---|---|
| Edge (Cloudflare / CDN) | DDoS mitigation, WAF, static cache |
| L7 LB (HAProxy) | TLS termination, WS upgrade, path routing, retry |
| Service (K8s) | Internal ClusterIP; L4 sticky yok |

## 4. TLS Termination

- Public: TLS 1.3 preferred, 1.2 fallback.
- Certs: ACME (Let's Encrypt) via cert-manager (K8s).
- Kesin: TLS **only at edge/LB**; internal service-to-service mTLS Phase 02.
- Cipher suite modern (`TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`).
- HSTS `max-age=63072000; includeSubDomains; preload`.

## 5. WebSocket Proxy

HAProxy config kalıbı:
```
frontend fe_https
  bind :443 ssl crt /etc/certs/api.pem alpn h2,http/1.1
  http-request set-header X-Forwarded-Proto https
  http-request set-header X-Real-IP %[src]
  acl is_ws hdr(Upgrade) -i websocket
  use_backend be_ws if is_ws
  default_backend be_api

backend be_ws
  option http-keep-alive
  timeout tunnel 1h
  timeout server 1h
  server app1 10.0.1.10:8080 check inter 5s
  server app2 10.0.1.11:8080 check inter 5s

backend be_api
  option httpchk GET /live
  http-check expect status 200
  balance roundrobin
  server app1 10.0.1.10:8080 check
  server app2 10.0.1.11:8080 check
```

WS tunnel timeout 1 saat (client heartbeat 25 sn); drop olursa client reconnect.

## 6. Routing

| Path | Backend | Notlar |
|---|---|---|
| `/v1/*` | be_api | Public API |
| `/internal/*` | 401 external; internal only | Firewall + service token |
| `/admin/*` | be_admin | VPN/SSO gate (bkz. [admin](../admin/)) |
| `/metrics` | 401 external; scraper only | Prometheus |
| `/ws` (Upgrade) | be_ws | WebSocket |

## 7. Health Check

- **Backend health**: `/live` her 5s; 3 fail → mark down.
- **Slow start**: yeni node önce %25 traffic → 60 sn ramp.
- **Graceful drain**: SIGTERM app → LB no new conn + wait active drain (30s max) → stop.

## 8. Rate Limit (Edge)

HAProxy stick-table:
```
stick-table type ip size 100k expire 30s store http_req_rate(10s)
http-request deny deny_status 429 if { src_http_req_rate() gt 100 }
```

Edge: coarse IP rate (DDoS baseline).
Application: fine-grained (account, character, endpoint — bkz. [security/rate-limit.md](../security/rate-limit.md)).

## 9. Retry Policy

- Idempotent GET: retry once on 5xx (max 2 attempts).
- POST: **no retry** (idempotency key ile client responsible).
- Timeout: connect 5s, server 30s, tunnel 1h (WS).

## 10. HA Setup

- **2+ HAProxy nodes** behind cloud LB (ALB) veya VRRP (keepalived).
- Config sync GitOps veya Consul KV.
- No session state (stateless LB).

## 11. Observability

- HAProxy stats page (internal only).
- Prometheus haproxy_exporter → metrics.
- Access log JSON → Loki.

## 12. Metrikler
- `haproxy_backend_up{backend}`
- `haproxy_backend_response_errors_total`
- `haproxy_frontend_requests_total`
- `haproxy_backend_current_queue`
- `haproxy_ssl_handshake_duration_seconds`

## 13. Cross-Reference
- [ssl.md](ssl.md), [cdn.md](cdn.md), [monitoring.md](monitoring.md)
- [performance/load-websocket.md](../performance/load-websocket.md)

## 14. Bekleyen Kararlar
- Cloud LB (ALB/GCLB) vs self-managed HAProxy — ops karar.
- Envoy migration Phase 02 (mesh + xDS dynamic).
- Slow start default 60s tune Phase 02.
