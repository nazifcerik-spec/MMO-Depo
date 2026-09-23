# Mail

## 1. Amaç
Asenkron mesaj + gold/item attachment. Karakter offline iken de alıp verebilme.

## 2. Karakteristik

- **Scope**: karakter-karakter (aynı account arası v1'de serbest, tax'sız — self mail).
- **Attachment**: max 12 item slot + gold.
- **Body**: max 512 char + 64 char subject.
- **Postage fee**: bkz. §Postage.
- **Expiration**: gönderim + **30 gün**; teslim alınmazsa (opened=false) sender'a **iade** (attachment gold/item + notification "returned").
- **Iade sonrası**: iade mesajı sender inbox'ta ek 7 gün; ondan sonra attachment overflow queue (envanter/wallet doluysa).
- **Cap**: inbox 100 mail (aşımda yeni gelen queue overflow'a; sender iade).
- **Escrow**: gold ve item mail domain escrow'unda (item lock benzer market escrow).

## 3. Model

```sql
CREATE TABLE social.mail (
  id UUID PRIMARY KEY,
  sender_character_id UUID NOT NULL,
  sender_display_name VARCHAR(64) NOT NULL,
  recipient_character_id UUID NOT NULL,
  subject VARCHAR(64) NOT NULL,
  body TEXT NOT NULL,                        -- 512 char
  gold_attachment_copper BIGINT NOT NULL DEFAULT 0,
  item_attachment_count SMALLINT NOT NULL DEFAULT 0,
  postage_paid_copper BIGINT NOT NULL DEFAULT 0,
  cod_price_copper BIGINT NULL,              -- cash on delivery (opt Phase 02)
  sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 days',
  opened_at TIMESTAMPTZ NULL,
  attachments_claimed_at TIMESTAMPTZ NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'sent', -- sent|read|claimed|expired|returned
  return_of_mail_id UUID NULL,               -- eğer bu mail iade ise orijinal ref
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE social.mail_attachment (
  mail_id UUID NOT NULL,
  slot_index SMALLINT NOT NULL,              -- 0..11
  item_instance_id UUID NOT NULL UNIQUE,
  stack_quantity INT NOT NULL DEFAULT 1,
  claimed BOOLEAN NOT NULL DEFAULT false,
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (mail_id, slot_index)
);

CREATE INDEX idx_mail_inbox ON social.mail (recipient_character_id, state, sent_at DESC);
CREATE INDEX idx_mail_expires ON social.mail (expires_at) WHERE state IN ('sent','read');
```

## 4. Postage Fee

Formül (bkz. [economy/fee.md](../economy/fee.md)):
```
postage = base + item_slot × per_item + round(gold_attach × gold_leg_rate)
base = 5c, per_item = 1c, gold_leg_rate = 0.001
```

Örnek: 5 item + 100g gold → `5 + 5×1 + round(1_000_000 × 0.001) = 1010c`.

Postage `system.fee_sink`'e ledger; iade edilmez (returned mail'de bile).

## 5. Send Akışı

`POST /v1/mail`
```json
{
  "recipient_character_name": "...",
  "subject": "...",
  "body": "...",
  "attachments": [{ "item_instance_id": "...", "stack_quantity": 1 }, ...],
  "gold_attach_copper": 0,
  "cod_price_copper": null,
  "idempotency_key": "..."
}
```

Tek SERIALIZABLE TX:
1. Idempotency check.
2. Recipient resolve.
3. Block check: recipient block sender → mail atılır (görünürde) ancak yalnızca `moderation_flag='blocked_target'`, teslim edilmez; attachment sender'a hiç düşmez (özel istisna — spam koruma), sender wallet/inv → mail escrow'a alınır ve **7 gün sender overflow'a** (spam engel; RMT bot ödemesin).
   - *Alternatif basit:* attachment mail'e gitmez, sender'da kalır, mesaj drop → 202 success. **V1 seçimi: basit.**
4. Ownership + lock: attachment item'lar `is_locked=true, lock_owner='mail'`; gold wallet -.
5. Escrow: `mail_attachment` insert + `social.mail` insert.
6. Postage tahsil: ledger `-postage → sender.wallet.gold, +postage → system.fee_sink`.
7. Recipient inbox cap check (100). Aşarsa `E_INBOX_FULL` → attachment unlock + rollback.
8. Outbox: WS "new_mail" → recipient (offline queue safe).
9. COMMIT.

## 6. Read / Claim Akışı

Read:
```
POST /v1/mail/{id}/read
```
- `state='read', opened_at=now()`.
- Read etmek attachment claim etmez; ayrı adım.

Claim (attachment):
```
POST /v1/mail/{id}/claim  body: { slot_indices: [0,1,2], claim_gold: true, idempotency_key }
```
Tek TX:
1. Mail lock; recipient == self; state IN ('read','sent').
2. Slot per item envanter kapasite check.
3. Gold cüzdan cap check.
4. Attachments transfer: item lock release + owner change; gold ledger.
5. `mail_attachment.claimed=true, claimed_at=now()`.
6. Tüm attach claimed + gold claimed ise `state='claimed', attachments_claimed_at=now()`.
7. COMMIT.

## 7. Delete

Sadece **hiç attachment kalmamış** mail delete edilebilir (`state='claimed'` veya attach zaten yoktu). Attachment varken delete engellenir (spam koruma — return of value).

```
DELETE /v1/mail/{id}
```

## 8. Expiration Worker

`mail_expire` her 5 dk:
```sql
SELECT id FROM social.mail
 WHERE state IN ('sent','read') AND expires_at <= now()
 FOR UPDATE SKIP LOCKED
 LIMIT 500;
```

Her mail:
- Attachment varsa iade mail oluştur (`state='returned', return_of_mail_id=orig`):
  - Sender'a yeni mail sub="Returned: <orig subject>", attach = orig attach (aynı item instance transfer), postage **yok**.
  - Orig mail `state='expired'`.
- Attachment yoksa: `state='expired'`, notify sender (opt).

Iade mail 7 gün expires_at; ondan sonra sender inbox_full olsa overflow queue.

## 9. Overflow (Envanter/Wallet Dolu)

Claim sırasında recipient env/wallet dolarsa:
- Partial claim: alabildiğini alır; kalanlar mail'de kalır.
- Full envanter → `E_INVENTORY_FULL`; kullanıcı bank slot açar, tekrar dener.

Iade mail sender overflow'da (7 gün) → sender claim edemezse `market.overflow_item` benzeri `social.mail_overflow_item`'a düşer 30 gün (nihai kayıp).

## 10. COD (Cash On Delivery — Phase 02)

Sender item için fiyat belirler; recipient claim ederken gold öder → sender'a otomatik döner. Marketplace alternatifi (1-to-1 satış). V1'de yok.

## 11. Block Etkileşimi

- Sender block edilmiş: mail atılır, silent drop (yukarıda §5.3).
- Recipient block edilmiş: sender-side reject `E_TARGET_BLOCKED_SELF` (kendi block ettiği kişiye mail atmaya çalışıyor).

## 12. Endpoint

**POST `/v1/mail`** — send.
**GET `/v1/mail/inbox?limit=50&cursor=`**.
**GET `/v1/mail/{id}`** — detail.
**POST `/v1/mail/{id}/read`**.
**POST `/v1/mail/{id}/claim`**.
**DELETE `/v1/mail/{id}`**.
**GET `/v1/mail/outbox?limit=50`** (Phase 02).

## 13. Rate Limit

- Send: 20/dk karakter, 60/gün account.
- Read/Claim: 60/dk (client polling ok).
- Attachment item slot: max 12 per mail.

## 14. Metrikler
- `mail_sent_total`, `mail_delivered_total`, `mail_expired_total`, `mail_returned_total`
- `mail_attachment_claim_total`
- `mail_postage_collected_copper_total`
- `mail_inbox_size_avg`
- `mail_overflow_item_gauge`

## 15. Cross-Reference
- [economy/fee.md](../economy/fee.md), [economy/transaction.md](../economy/transaction.md)
- [inventory/item-lock.md](../inventory/item-lock.md), [inventory/item-transfer.md](../inventory/item-transfer.md)
- [block.md](block.md), [chat.md](chat.md)
- [systems/notification.md](../systems/notification.md) — new_mail push

## 16. Bekleyen Kararlar
- COD (Phase 02).
- Item attachment expiration hızlandırma (< 30 gün rare item için) Phase 02.
- Outbox (gönderilen mail history) Phase 02.
- Mail alliance / cross-server Phase 03.
- Rich text mail (Phase 02 opt).
