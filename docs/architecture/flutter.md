# Flutter (Client)

## 1. Amaç
Oyunun tüm istemci yüzeyini tek bir codebase ile sağlamak: web (öncelikli), iOS, Android, desktop. Sadece sunum ve input katmanı; server otoritesi hiç değişmez.

## 2. Sürüm ve Bağımlılıklar
- **Flutter**: en son kararlı major (Phase 01 başlangıcında pin edilecek; kanal `stable`).
- **Dart SDK**: Flutter'la birlikte gelen sürüm.
- **State management**: `riverpod` (test edilebilir, DI-native).
- **Router**: `go_router`.
- **HTTP**: `dio` (interceptor: auth, retry, telemetry).
- **WebSocket**: `web_socket_channel`.
- **JSON**: `freezed` + `json_serializable`.
- **Localization**: `intl` + `.arb` dosyaları (TR + EN).
- **Storage**: `flutter_secure_storage` (token), `shared_preferences` (settings).
- **Testing**: `flutter_test`, `mocktail`, `golden_toolkit`.

Kesin sürümler `pubspec.yaml`'da pinlenir; `pubspec.lock` commit edilir.

## 3. Katman Ayrımı
```
lib/
├── app/                # entrypoints, root widget, routing, theme
├── core/               # error, network, result types, base classes
├── data/               # datasources (rest, ws), repositories, DTOs
├── domain/             # entities, use-cases, repository interfaces
├── features/           # feature-scoped UI + state
│   ├── auth/
│   ├── character/
│   ├── activity/
│   ├── inventory/
│   ├── marketplace/
│   ├── guild/
│   └── ...
├── l10n/               # generated intl
├── shared/             # widgets, tokens, formatters
└── main_<flavor>.dart  # dev / staging / prod entrypoints
```

## 4. Kurallar
- **State bir yerde**: her feature'ın state'i Riverpod provider'larıyla temsil edilir; global mutable singleton yok.
- **Server otoritedir**: UI, cache'ten hızlı render eder ama her mutation server confirm bekler; optimistic update yalnızca reversible aksiyonlar için.
- **Time**: sunum için device time; ödül hesabında asla client time (server'dan gelen zaman kullanılır).
- **Error contract**: `Result<T, AppError>` pattern; exception UI'a taşınmaz.
- **Localization**: literal string yasak (i18n key zorunlu). CI check.
- **Accessibility**: her interactive widget'ın semantic label'ı; contrast ratio kontrolü.
- **Responsive**: layout `mobile-first` düşünülür, desktop için `LayoutBuilder` scale.
- **No secrets in code**: API base URL yalnızca build-time flag'den (`--dart-define`).

## 5. Networking Sözleşmesi
- Base URL flavor-based (`API_BASE_URL` dart define).
- Auth: `Authorization: Bearer <access_token>`; 401 → refresh flow; 401 tekrar ederse logout.
- Idempotency-Key header: kritik POST'larda otomatik UUID.
- Retry: only idempotent methods; exponential backoff (max 3 attempt).
- Telemetry: request/response trace id header (`X-Request-ID`), OTLP export.

## 6. WebSocket Kanalı
- Tek connection; multiplex message types (`activity_tick`, `notification`, `chat`, `world_event`, `reveal_ready`).
- Otomatik reconnect (exponential backoff, jitter).
- Heartbeat: 30 sn ping.
- Uygulama background'a giderse: mobile'da suspend; foreground'da resume.

## 7. Test
- **Unit**: use-case ve pure function.
- **Widget**: feature root widget'ları.
- **Golden**: kritik ekranlar (character sheet, market, inventory).
- **Integration**: build flavor per env, E2E happy path.
- Coverage hedefi: %70 (Phase 01), %80 (Phase 02).

## 8. Performance Hedefleri
- Cold start (web): < 3 sn TTI (average connection).
- Frame budget: 60 fps mobile, 120 fps desktop hedef.
- Bundle size: web initial < 2 MB gzip (Phase 01 pass).
- Memory: mobile idle < 250 MB.

## 9. Multi-Platform Notları
- **Web**: canvas rendering veya HTML mode; SEO gerekli değil (login-gated app).
- **iOS / Android**: push notification (opt-in, FCM/APNS).
- **Desktop**: opsiyonel; MVP'de Windows shell.
- Platform-specific kod `Platform.isX` yerine adapter pattern.

## 10. Build & Flavor
- Flavor'lar: `dev`, `staging`, `prod`.
- Config `--dart-define` ile: `FLAVOR`, `API_BASE_URL`, `WS_URL`, `SENTRY_DSN`.
- CI build: her flavor için ayrı artifact.

## 11. Anti-Patterns
- Business logic widget içinde.
- Provider dışı mutable global state.
- Retry'siz network call.
- Inline literal string (localization dışı).
- Sync blocking I/O.

## 12. İlişkiler
- [System Architecture](system-architecture.md), [Go Server](go.md), [Config](config.md), [Environment](environment.md), [Logging](logging.md), [Notification System](../systems/notification-system.md)

## 13. Bekleyen Kararlar
- State management alternatif: `riverpod` vs `bloc` (öneri: riverpod; final Phase 01 sprint 0).
- Web renderer: HTML mi CanvasKit mi (initial load vs render quality trade-off).
- Native mobile shell öncelik sırası (Android ilk mi iOS ilk mi).
- Push provider (FCM cross-platform vs platform-native).
