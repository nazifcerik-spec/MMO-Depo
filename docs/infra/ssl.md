# SSL / TLS

## 1. Amaç
TLS sertifika management: ACME (Let's Encrypt) otomatik rotation, HSTS, modern cipher suite.

## 2. Sertifika Kapsamı

| Domain | Cert Type | Notlar |
|---|---|---|
| `game.io` (landing) | Cloudflare Universal SSL | CDN yönetir |
| `cdn.game.io` | Cloudflare | CDN yönetir |
| `api.game.io` | ACME (Let's Encrypt) via cert-manager | HAProxy termination |
| `ws.game.io` | ACME | HAProxy termination (WSS) |
| `admin.game.internal` | ACME (private CA opsiyonel) | VPN gate |
| `*.internal` | Internal CA | Service-to-service (Phase 02 mTLS) |

## 3. cert-manager (K8s)

`Certificate` CR + `Issuer` ACME:
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: api-game-io
spec:
  secretName: api-game-io-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames: [api.game.io, ws.game.io]
```

- DNS-01 challenge (Cloudflare API token) — wildcard support.
- 90 gün cert; auto-renew 30 gün önce.
- Rotation zero-downtime (HAProxy reload).

## 4. Cipher Suite (TLS 1.3)

Modern only:
```
TLS_AES_128_GCM_SHA256
TLS_AES_256_GCM_SHA384
TLS_CHACHA20_POLY1305_SHA256
```

TLS 1.2 fallback:
```
ECDHE-ECDSA-AES128-GCM-SHA256
ECDHE-RSA-AES128-GCM-SHA256
ECDHE-ECDSA-AES256-GCM-SHA384
```

Weak cipher **yasak**: 3DES, RC4, MD5, SHA-1.

## 5. HSTS

Response header:
```
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
```

HSTS preload list submit → yeni tarayıcı otomatik HTTPS.

## 6. Certificate Transparency

- ACME cert otomatik CT log.
- Monitor: crt.sh scan for unauthorized cert issuance (nightly script).

## 7. mTLS (Service-to-Service) — Phase 02

- Internal CA (cert-manager `Issuer` self-signed).
- Every service pod: certificate mount, verify peer cert.
- Envoy sidecar or app-native (Go crypto/tls).
- Zero-trust internal (bkz. [security/README.md](../security/README.md) DD-186).

## 8. Client Pinning

Mobile client v1'de **no cert pinning** (rotation risk).
Phase 02 opsiyonel dynamic pinning (backup pin + primary) — MITM risk minimal.

## 9. Cert Alarm

- 15 gün önceden **warn** (Slack).
- 7 gün önceden **high** (Slack + email).
- Cert expire → **critical** (PagerDuty).

## 10. Rotation Prosedürü

Otomatik (cert-manager) — insan müdahalesi yok. Manuel rotation:
1. cert-manager `Certificate` renew force.
2. Wait new secret populate.
3. HAProxy hot reload (`kill -USR2`).
4. Verify TLS handshake.

## 11. Metrikler
- `cert_expiry_seconds{domain}`
- `cert_manager_renewal_success_total`
- `tls_handshake_duration_seconds`

## 12. Cross-Reference
- [reverse-proxy.md](reverse-proxy.md), [cdn.md](cdn.md)
- [security/](../security/), [architecture/docker.md](../architecture/docker.md)

## 13. Bekleyen Kararlar
- mTLS service mesh (Envoy / Linkerd) — Phase 02.
- Internal CA provider (cert-manager / Vault PKI) — Phase 01.
- Client pinning strategy — Phase 02.
- CT monitoring tool (crt.sh script vs Sigsci) — Phase 02.
