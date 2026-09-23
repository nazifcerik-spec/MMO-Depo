# Mobile Lifecycle (Android + iOS Background)

## 1. Amaç
Client (Flutter) mobile platformda background/foreground state, push wake, WS reconnect, background fetch politikası.

**Kritik prensip:** oyunun offline "AFK engine" tamamen **server-side**; client kapalı olsa da aktivite server'da devam eder. Client sadece push tetiklendiğinde açılır ve reveal isteği yapar.

## 2. Client State Modeli

| State | Açıklama | WS | Push |
|---|---|---|---|
| `foreground_active` | Uygulama açık, kullanıcı etkileşimde | Bağlı | Sessiz (in-app banner) |
| `foreground_inactive` | Uygulama açık ama arka plan sekmesi (multitasking) | Bağlı | Sessiz banner |
| `background_recent` | Son 5dk arka plana atıldı | Bağlı kalır (grace) | OS push (silent) |
| `background_deep` | 5dk+ arka planda / minimize | Kapalı | OS push (visible) |
| `terminated` | Uygulama kapatıldı / OS kill | Kapalı | OS push (visible + wake) |

## 3. Android Lifecycle

### 3.1 Foreground Service Yok
V1'de foreground service **kullanılmaz** — battery ve Play Store policy. Client sadece push wake ile aktifleşir.

### 3.2 Doze Mode + App Standby
- **Doze**: kullanıcı telefonu bırakınca (screen off + still) OS "deep sleep"; network suspend.
- **App Standby**: uzun süre kullanılmayan app suspend.
- **Etki**: WS bağlantı düşer; push (FCM high priority) hâlâ ulaşır.

### 3.3 FCM Delivery
- **High priority** (`android.priority=high`) — Doze'u aşar, cihazı wake eder (limited quota).
- **Normal priority** — Doze bitene kadar bekler (batched delivery).
- V1 politika: `activity_completed`, `dungeon_completed`, `world_event_*`, `whisper_received`, `guild_event` **high**; `energy_restored`, `friend_online`, `market_expiring` **normal**.

### 3.4 Notification Channel Kayıt
Client başlangıçta channel kayıt eder (Android 8+):
```dart
[
  { id: 'activity', name: 'Aktivite', importance: HIGH, sound: default },
  { id: 'event', name: 'Etkinlik', importance: HIGH, sound: default },
  { id: 'social', name: 'Sosyal', importance: DEFAULT, sound: default },
  { id: 'system', name: 'Sistem', importance: HIGH, sound: default },
  { id: 'silent', name: 'Sessiz', importance: LOW, sound: null },
]
```

Kullanıcı OS Ayarlar'dan tek tek kapatabilir — bu server-side preference'ı **override etmez**; server hâlâ gönderir, OS filtreler.

## 4. iOS Lifecycle

### 4.1 Background Modes
Info.plist:
```xml
<key>UIBackgroundModes</key>
<array>
  <string>remote-notification</string>   <!-- silent push için -->
  <string>fetch</string>                 <!-- background fetch -->
</array>
```

**Kullanılan modlar:**
- `remote-notification` — silent push (`content-available: 1`) ile app kısa süre wake (30 sn budget).
- `fetch` — sistem periyodik fetch (opsiyonel, opportunistic).

**Kullanılmayan:**
- `audio` (v1'de audio play arka plan yok).
- `location` (idle MMO'da location tracking yok).
- `voip` (v1'de yok).

### 4.2 APNs Priority
- `apns-priority: 10` — kritik push (activity, dungeon, world event).
- `apns-priority: 5` — silent/low (energy, market_expiring).

### 4.3 Background Task Budget
Silent push wake budget ~30 sn. Client:
1. Reveal API çağırma **yok** (server zaten cek — reveal user açtığında).
2. Basit sync (unread count, badge update).
3. Return.

## 5. WS Reconnect Stratejisi

Foreground return:
1. WS state check.
2. Bağlı değilse reconnect (exp backoff 1s, 2s, 5s, 15s, cap 60s).
3. Reconnect'te `subscribe:{channels}` + `resume_from:<last_msg_id>` (misseD event replay son 5 dk).
4. `character.presence` update (heartbeat).

Background transition:
1. WS grace period 5 dk (Android/iOS keep-alive sınırında).
2. 5 dk sonra WS close (client-side); server presence 5 dk sonra "offline".

## 6. Deep Link Routing

Notification tıklama → client açılır + route:
```
game://route/character/activity/<contract_id>/reveal
game://route/guild/<guild_id>/boss/<spawn_id>
game://route/mail/<mail_id>
```

Client router karşılık kuralı; auth'suz açılırsa login sonrası queued deep link resolve.

## 7. Cold Start Optimizasyonu

Notification tap → cold start (app terminated) → hedef ekrana varış hedef ≤ 3 saniye. Adımlar:
1. Splash + auth restore (token cache).
2. Reveal API prefetch (deep link parse).
3. UI reveal ekran açılır (loot animation).

Slow cold start (5+ sn) → drop-off riski. Client team KPI.

## 8. Battery / Data

- Foreground WS bandwidth ~ 5-15 KB/dk (chat + activity ping).
- Background zero (WS kapalı).
- Push payload ≤ 1 KB target.
- Cold start network warm-up ~50 KB (character state).

## 9. Metrikler (Client-side reported)
- `client_cold_start_ms{route}`
- `client_ws_reconnect_success_ratio`
- `client_background_transition_total`
- `client_deep_link_open_total{route}`
- `client_notification_tap_total{kind}`

## 10. Cross-Reference
- [push-notification.md](push-notification.md)
- [backend/websocket.md](../backend/websocket.md)
- [auth/session.md](../auth/session.md), [auth/device.md](../auth/device.md)
- [architecture/flutter.md](../architecture/flutter.md)

## 11. Bekleyen Kararlar
- Foreground service Android (Phase 03 — persistent connection için opt-in, battery cost).
- Background fetch iOS priority (v1 opportunistic; Phase 02 tune).
- WS keep-alive interval (v1: 25 sn ping; Phase 02 mobile-optimize adaptive).
- Cold start budget hedef Phase 02 UX pass (2 sn?).
