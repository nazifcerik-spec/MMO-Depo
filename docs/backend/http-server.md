# HTTP Server

## 1. Amaç
Bir tek `http.Server` instance'ı — routing (`chi`), TLS (edge'de), timeout, graceful shutdown, healthz. Domain kodu bu katmanı görmez.

## 2. Sözleşme
```go
package httpx

type Server struct {
    cfg    Config
    router *chi.Mux
    srv    *http.Server
    log    *slog.Logger
}

type Config struct {
    Addr              string        // ":8080"
    ReadHeaderTimeout time.Duration // 5s
    ReadTimeout       time.Duration // 10s
    WriteTimeout      time.Duration // 15s
    IdleTimeout       time.Duration // 60s
    MaxHeaderBytes    int           // 1<<20 (1MB)
    ShutdownGrace     time.Duration // 30s
}

func NewServer(cfg Config, router *chi.Mux, log *slog.Logger) *Server
func (s *Server) Run(ctx context.Context) error
```

## 3. Boot Akışı
```
1. Config load & validate  (bkz. Config)
2. Init telemetry (OTEL + slog + metrics endpoint on :9090)
3. Init DB pool  (bkz. Database)
4. Init Redis pool (bkz. Redis)
5. Wire domain modules → repositories → services → handlers
6. Build router (bkz. Router)
7. Attach middlewares
8. Start metrics server (goroutine)
9. Start HTTP server (Run(ctx))
10. Wait for SIGTERM/SIGINT → cancel ctx → Server.Shutdown(grace)
```

## 4. Örnek Skeleton
```go
func (s *Server) Run(ctx context.Context) error {
    s.srv = &http.Server{
        Addr:              s.cfg.Addr,
        Handler:           s.router,
        ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
        ReadTimeout:       s.cfg.ReadTimeout,
        WriteTimeout:      s.cfg.WriteTimeout,
        IdleTimeout:       s.cfg.IdleTimeout,
        MaxHeaderBytes:    s.cfg.MaxHeaderBytes,
        BaseContext:       func(net.Listener) context.Context { return ctx },
    }

    errCh := make(chan error, 1)
    go func() {
        s.log.Info("http server listening", "addr", s.cfg.Addr)
        if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            errCh <- err
        }
        close(errCh)
    }()

    select {
    case <-ctx.Done():
        shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownGrace)
        defer cancel()
        s.log.Info("http server shutting down", "grace", s.cfg.ShutdownGrace)
        return s.srv.Shutdown(shutdownCtx)
    case err, ok := <-errCh:
        if !ok {
            return nil
        }
        return fmt.Errorf("http server: %w", err)
    }
}
```

## 5. Timeouts (Default)
| Setting | Value | Neden |
|---|---|---|
| ReadHeaderTimeout | 5s | Slowloris koruma |
| ReadTimeout | 10s | Normal request; body okuma |
| WriteTimeout | 15s | Handler yanıt bütçesi |
| IdleTimeout | 60s | Keep-alive |
| ShutdownGrace | 30s | In-flight request drain |

WebSocket handler'ı bu limitlerden bağımsız kendi context yönetimi ile çalışır (bkz. [WebSocket](websocket.md)).

## 6. TLS
- **Edge terminate**: TLS proxy/ingress (nginx, Envoy, ALB) tarafından; app plain HTTP (localhost).
- Internal mTLS gerekirse `srv.TLSConfig` set edilir; edge → app tunnel.

## 7. Healthz
İki endpoint (bkz. [Router](router.md)):
- `GET /healthz/live` — süreç yaşıyor (200 boş).
- `GET /healthz/ready` — DB + Redis ping OK → 200; aksi 503.

## 8. Graceful Shutdown
- SIGTERM alındığında ctx cancel.
- `Server.Shutdown` yeni bağlantıları kapatır, aktif request'lerin bitmesini bekler.
- WebSocket hub kendi shutdown'ını dinler (aktif conn'lara `close` frame gönderir).
- Worker'lar ayrı process (bkz. Worker), ayrı shutdown.

## 9. Panic Handling
- Recover middleware her handler etrafında (bkz. Middleware).
- Panic → 500 + structured log + metric increment.

## 10. Fail Modları
- Port in use → boot fail, log.
- TLS cert missing → boot fail (edge TLS bile isteniyorsa health check).
- DB unreachable at boot → boot fail (fail-fast).
- Ready probe DB ping fail → 503 (Kubernetes traffic'i keser).

## 11. Test Stratejisi
- Integration test: `httptest.NewServer(router)` ile end-to-end handler test.
- Shutdown test: SIGTERM sim + in-flight request'in tamamlanması.

## 12. Metrics
- `http_requests_total{method,route,status}` counter.
- `http_request_duration_seconds{method,route}` histogram.
- `http_in_flight_requests` gauge.
- `http_open_connections` gauge.

## 13. İlişkiler
- [Router](router.md), [Middleware](middleware.md), [Logger](logger.md), [Database](database.md), [Redis](redis.md), [WebSocket](websocket.md)

## 14. Bekleyen Kararlar
- HTTP/2 & HTTP/3 support: edge halleder v1; app tarafı stdlib default.
- Compression: middleware layer'da gzip/brotli (Phase 02 tune).
- Response streaming pattern (progressive JSON) — reveal payload büyükse.
