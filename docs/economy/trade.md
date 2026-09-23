# Trade (Player-to-Player)

## 1. Amaç
İki oyuncu arasında güvenli, atomik değişim: item + gold ↔ item + gold.

## 2. Prensip
- **Atomik**: iki taraf da onayladıktan sonra tek TX'te swap yapılır. Yarım transfer imkansız.
- **Server-authoritative**: client sadece intent gönderir; validation ve swap serverda.
- **Audit-first**: trade window + confirm + finalize tüm adımlar `audit.trade_log`'a yazılır.
- **Anti-scam**: son değişiklikten sonra karşı taraf onayı **reset** olur.

## 3. Model

```sql
CREATE TABLE economy.trade_session (
  id UUID PRIMARY KEY,
  initiator_character_id UUID NOT NULL,
  partner_character_id UUID NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'open',
  -- open | initiator_locked | partner_locked | both_locked
  -- initiator_confirmed | partner_confirmed | both_confirmed
  -- finalized | cancelled | expired
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '10 minutes',
  finalized_at TIMESTAMPTZ NULL,
  cancel_reason VARCHAR(48) NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE economy.trade_offer (
  trade_session_id UUID NOT NULL,
  side VARCHAR(16) NOT NULL,          -- 'initiator' | 'partner'
  gold_copper BIGINT NOT NULL DEFAULT 0 CHECK (gold_copper >= 0),
  items JSONB NOT NULL DEFAULT '[]',  -- [{item_instance_id, quantity}]
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (trade_session_id, side)
);
```

## 4. State Machine

```
open
 ├── (initiator update / partner update) → open (unlocked; version++)
 ├── initiator lock → initiator_locked
 ├── partner lock  → partner_locked
 └── both lock     → both_locked
       ├── either confirm → *_confirmed
       └── both confirm   → finalized (TX swap)

anywhere:
 ├── cancel (user or auto expire) → cancelled / expired
 └── update on locked side        → auto-unlock both sides (anti-scam reset)
```

## 5. Akış

**1. Trade request**
```
POST /v1/trade/request
  body: { partner_character_id }
```
Karşı taraf pending trade request alır (WS notify). Kabul → session oluştu (`open`).

**2. Update offer**
```
PUT /v1/trade/{session_id}/offer
  body: {
    gold_copper,
    items: [{item_instance_id, quantity}],
    idempotency_key
  }
```
Server:
- Item ownership + `is_trade_locked=false` (bkz. [inventory/item-lock.md](../inventory/item-lock.md)).
- Item bind kısıtı (BoP → trade edilemez, `E_ITEM_NOT_TRADABLE`).
- gold_copper ≤ wallet.gold_copper.
- Update sonrası **her iki taraf otomatik unlock** olur.
- `version++`, WS her iki tarafa güncel offer push.

**3. Lock**
```
POST /v1/trade/{session_id}/lock  (kendi tarafın)
```
- Karşı taraf offer değiştirirse benim lock'um auto-drop olur.

**4. Confirm** (her iki taraf lock'tan sonra)
```
POST /v1/trade/{session_id}/confirm
  body: { seen_partner_offer_hash }
```
- `seen_partner_offer_hash` client'ın son gördüğü offer hash'i; mismatch → `E_TRADE_OFFER_CHANGED` (anti bait-and-switch).

**5. Finalize** (her iki taraf confirmed olunca server otomatik):
Tek SERIALIZABLE TX:
1. Her iki tarafın item + gold hâlâ mevcut mu ve lockable mı → recheck.
2. Item lock (idempotent).
3. Ledger yaz: A→B gold delta, B→A gold delta.
4. Item ownership swap (`character_id` update + inventory slot temiz).
5. Trade tax uygula (bkz. [tax.md](tax.md)).
6. `trade_session.state='finalized'`, `finalized_at=now()`.
7. audit.trade_log entry.
8. WS her iki tarafa "trade completed".

**6. Cancel / Expire**
Herhangi biri cancel → state='cancelled', reason yazılır.
10 dk boyunca finalize edilmezse → `expired` (worker).

## 6. Anti-Fraud / Güvenlik

- **Offer değişimi otomatik unlock** — Fitzgerald bait-and-switch engeli.
- **Confirm'de hash** — client'ın gördüğü offer, server'ın state'i ile eşleşmeli.
- **Trade cooldown**: aynı iki karakter arası cancel sonrası 30 sn (spam koruma).
- **Same account block**: aynı account'un iki karakteri arasında trade v1'de **serbest** (self-fund önlem yok); Phase 02'de anti-farm review.
- **Rate limit**: karakter başına dakikada max 3 trade session açılışı.
- **Bind rules**:
  - BoP item: trade edilemez.
  - BoA item: sadece aynı account (v1'de trade UI yok, mail transfer Faz 26).
  - Quest / soulbound item: trade edilemez.

## 7. Trade Tax
V1'de trade tax **yok** (dostluk-friendly). Phase 02'de high-value trade için opsiyonel flat tax (bkz. [tax.md](tax.md) §Trade Tax).

## 8. Endpoint

**POST `/v1/trade/request`** / **POST `/v1/trade/request/{id}/accept`**
**PUT `/v1/trade/{session_id}/offer`**
**POST `/v1/trade/{session_id}/lock`** / **POST `/v1/trade/{session_id}/unlock`**
**POST `/v1/trade/{session_id}/confirm`**
**POST `/v1/trade/{session_id}/cancel`**
**GET `/v1/trade/{session_id}`**
**GET `/v1/character/{id}/trade/history`** (kişisel)

## 9. Metrikler
- `trade_session_opened_total`, `trade_session_finalized_total`, `trade_session_cancelled_total{reason}`
- `trade_gold_transferred_copper_total`
- `trade_offer_hash_mismatch_total`
- `trade_finalize_latency_ms`

## 10. Cross-Reference
- [transaction.md](transaction.md) — SERIALIZABLE swap TX
- [inventory/item-lock.md](../inventory/item-lock.md), [inventory/item-transfer.md](../inventory/item-transfer.md)
- [tax.md](tax.md)
- [notification](../systems/notification.md) — trade request WS
- [database/audit.md](../database/audit.md) — audit.trade_log

## 11. Bekleyen Kararlar
- Trade tax formülü Phase 02 (flat vs %).
- Same-account trade Phase 02 anti-farm review.
- Trade log görüntüleme UI (kullanıcı kendi history v1; başkasınınki yok).
- Mail item transfer Faz 26.
