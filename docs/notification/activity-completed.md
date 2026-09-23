# Activity Completed Notification

## 1. Amaç
AFK activity contract reveal → oyuncuya "aktivite tamamlandı" push.

Bu bildirim retention'ın temel taşı — oyuncu kontrat başlar ("2 saat mining"), süre dolar → push → geri döner → reveal + reward → yeni aktivite.

## 2. Trigger

`event.activity.contract_completed`
```json
{ "character_id": "...", "contract_id": "...", "kind": "mining", "yield_summary": { "ore": 42, "gem_rare": 1 } }
```

Push dispatcher:
1. Preferences check (kind='activity_completed' opt-in? default ON).
2. Quiet hours check.
3. Rate limit (bkz. push-notification.md; 6/saat, cooldown 5dk).
4. Dedup key: `activity_completed:<contract_id>`.
5. Payload build.
6. Send FCM/APNs + inbox fallback.

## 3. Payload

```json
{
  "title_key": "notif.activity_completed.title",
  "body_key": "notif.activity_completed.body.mining",
  "body_vars": { "yield_summary": "42 iron ore, 1 rare gem" },
  "category": "activity",
  "route": "character/activity/<contract_id>/reveal",
  "badge_delta": 1,
  "collapse_key": "activity_<character_id>"
}
```

- **collapse_key** aynı karakter için birden fazla contract kısa süre içinde biterse OS notification merge (Android FCM özelliği). iOS `thread-id="activity"`.

## 4. Reveal Flow (Client)

1. Kullanıcı push'a tıklar → deep link `character/activity/<id>/reveal`.
2. Client açılır (background → foreground) → `POST /v1/activity/reveal` çağrısı.
3. Server reveal hesabı (analitik — DD-084) → response reward summary.
4. Client "loot animation" → reward apply.

Idempotent: aynı reveal ikinci kez çağrılırsa cached result.

## 5. Silent Case

Bazı kullanıcılar sadece badge/inbox ister. Preference `notif.activity_completed.mode`:
- `push_and_inbox` (default)
- `inbox_only` (in-app, no push)
- `silent_push` (background wake, badge only)
- `off`

## 6. Chained Notification

Ardışık aynı karakter için 3+ contract tamamlanırsa (rare — çoklu concurrent AFK) **grouped notification**:
- Title: "3 aktivite tamamlandı — Grim"
- Body: özet (kısa)
- Client açılınca reveal listesi.

## 7. Örnek Notification Metni (TR)

```
Başlık: Aktivite Tamamlandı — Grim
İçerik: Madencilik tamamlandı: 42 demir cevheri, 1 nadir taş.
```

Localization key'ler `notif.activity_completed.body.<kind>` (mining, fishing, hunting, ...).

## 8. Metrikler
- `notif_activity_completed_sent_total{kind}`
- `notif_activity_completed_click_open_total`
- `notif_activity_completed_reveal_delay_ms`  -- notification → reveal click süresi

## 9. Cross-Reference
- [push-notification.md](push-notification.md), [preferences.md](preferences.md)
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)

## 10. Bekleyen Kararlar
- Yield summary tekst uzunluk (mobile OS truncation) Phase 02 UX pass.
- "Auto-start next" quick-action (Phase 02 iOS/Android action button).
- Preview loot image thumbnail Phase 02 rich notification.
