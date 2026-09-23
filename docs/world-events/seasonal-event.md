# Seasonal Event

## 1. Amaç
Uzun süreli (2-6 hafta) tema event'i. Login incentive + event currency + kozmetik ödül + geçici mini-game/quest track.

Örnek: **Winter Feast**, **Anniversary Week**, **Summer Games**, **Harvest Moon**.

## 2. Karakteristik

- **Süre**: 2-6 hafta.
- **Global**: tüm oyunculara açık, tüm bölgelerde asset (dekorasyon vs).
- **Event currency**: dedicated `event_<code>` token (bkz. [economy/currency.md](../economy/currency.md)).
- **Quest track**: 20-40 quest (daily / one-off) event boyunca.
- **Event vendor**: token karşılığı özel item (kozmetik, geçici mount, minipet, boost consumable).
- **AFK login bonus**: hesap seviyesinde günlük event token drop; giriş = full ödül.
- **Kalıcı ödül**: kozmetik / minipet / title (permanent); consumable geçici (event bitince kalır ama vendor kapanır).

## 3. Model

```sql
CREATE TABLE content.seasonal_event_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  token_currency_code VARCHAR(48) NOT NULL,
  quest_pack_id UUID NOT NULL,
  vendor_pack_id UUID NOT NULL,
  cosmetic_bundle_id UUID NULL,
  daily_login_reward JSONB NOT NULL,
  version INT DEFAULT 1
);

CREATE TABLE progression.seasonal_event_participation (
  account_id UUID NOT NULL,
  seasonal_event_id UUID NOT NULL,
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  event_token_earned BIGINT NOT NULL DEFAULT 0,
  event_token_spent BIGINT NOT NULL DEFAULT 0,
  quest_completed_count INT NOT NULL DEFAULT 0,
  cosmetic_unlocked JSONB NOT NULL DEFAULT '[]',
  PRIMARY KEY (account_id, seasonal_event_id)
);
```

## 4. Token (Currency)

- `event_<event_code>` (örn. `event_winter_feast`).
- Scope: **account** (tüm karakterler paylaşır — DD-137 currency scope 'account').
- Tradable: **hayır** (BoA).
- Cap: 10^5 (100k).
- **TTL**: event ends_at + 14 gün → otomatik silme; vendor 14 gün alışverişe açık kalır post-event.

Ledger yolu ekonomi ile aynı (`TxnService.Write` → `economy.account_currency_balance` + `audit.ledger_audit`).

## 5. Kaynaklar

| Kaynak | Yaklaşık Token |
|---|---|
| Daily login | 5-15 |
| Quest completion (daily) | 5-20 |
| Quest completion (one-off epic) | 50-200 |
| World boss attendance | 5-40 |
| Regional event | 2-10 |
| AFK contract "seasonal_bonus" | 1-5 / hour |

## 6. Vendor (Event Store)

Event window boyunca (+14 gün post) açık:
- Kozmetik (armor skin, weapon skin, tabard, mount skin).
- Geçici minipet / mount (seasonal, permanent unlock).
- Event food (buff consumable).
- Rare mat pack (limited stock).

Vendor `content.vendor_template` içinde `vendor_type='seasonal'`, `currency_code='event_<code>'`. Bkz. [economy/vendor.md](../economy/vendor.md).

## 7. Quest Pack

- Daily quest 3-5 tanesi (daily reset 00:00 UTC — DD-016 uyumlu).
- One-off quest 15-25 tanesi (event boyunca bir kez).
- Ana line ("Ashen King's Return" gibi) 5-10 quest, sonunda title + cosmetic.

Bkz. [quests](../quests/README.md).

## 8. AFK Login Bonus

Hesap event window'da her gün ilk giriş → daily login reward (event token + minor item).

**AFK versiyonu**: hiç login yoksa da daily token'un yarısı otomatik `event_<code>` currency'ye yazılır (attendance credit) — event'in "gerçekten kaçırılmasın" garantisi.

Config:
```yaml
seasonal_event:
  daily_login_reward: { event_token: 15, item: "festival_food_stub" }
  auto_afk_credit_ratio: 0.50   # login yoksa günlük tokenin %50'si
  auto_afk_credit_max_days: 7   # 7 günün üstünde kaybolur
```

## 9. Cosmetic / Permanent Ödül

- Kozmetik (skin) unlock: account bound, permanent.
- Title unlock: character bound.
- Minipet / mount: account bound.

Event bitince vendor kapanır ama unlock kalır. Yeni oyuncu event kaçırdıysa gelecek yıl aynı event'te aynı ödül teklif edilir (rotation).

## 10. Endpoint

**GET `/v1/seasonal/active`** — açık event(s).
**GET `/v1/seasonal/{event_code}/vendor`**
**GET `/v1/seasonal/{event_code}/quests`**
**GET `/v1/account/{id}/seasonal/{event_code}/progress`**

## 11. Metrikler
- `seasonal_event_active_gauge`
- `seasonal_token_earned_total{event, source}`
- `seasonal_token_spent_total{event, sku}`
- `seasonal_cosmetic_unlock_total{event, sku}`
- `seasonal_daily_login_ratio`

## 12. Cross-Reference
- [economy/currency.md](../economy/currency.md), [economy/vendor.md](../economy/vendor.md)
- [quests/event.md](../quests/event.md)
- [systems/activity-engine.md](../systems/activity-engine.md)
- [character/titles.md](../character/titles.md)

## 13. Bekleyen Kararlar
- Event pack yayın mekanizması (hot bundle vs full patch) Phase 02.
- Rotation takvimi (yıllık) content pack Phase 02.
- Premium event boost (gem → event token dönüşümü) yasağı v1 (DD-137 uyumlu); Phase 02 review.
- AFK auto-credit oranı final tune Phase 02.
