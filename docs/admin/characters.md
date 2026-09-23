# Characters (Character Inspect)

## 1. Amaç
Bir karakterin state'ini derinlemesine görmek — stats, inventory, equipment, aktif contract, quest, guild, party, achievement, mail, chat log (yetkiye göre).

## 2. Detail View

```
Character: Grim (char_ab12…)
Level 42 · Warrior · Human · Ironhold
State: active [Freeze] [Restore]

┌─ Stats ────────────────┐
│ HP 850/950 · MP 200/300│
│ STR 45 AGI 28 INT 12  │
└────────────────────────┘

┌─ Resources ────────────┐
│ Wallet: 12,340c        │
│ Energy: 87/100         │
│ Achievement pts: 285   │
└────────────────────────┘

Tabs:
[Inventory] [Equipment] [Quests] [Activity] [Guild] [Party]
[Achievements] [Mail] [Chat Log] [Ledger] [Audit]
```

## 3. Sub-tabs

### Inventory
Item slot grid; her item için detail (source_ref, obtained_at, current lock state, market listing var mı).
Aksiyonlar: **restore item** (soft-deleted item resurrect), **remove item** (dupe/exploit sonucu).

### Equipment
Paper doll; slot bazlı equipped item.
Aksiyon: **unequip** (env transfer).

### Quests
Aktif + tamamlanan quest listesi.
Aksiyon: **reset quest** (support use case), **skip quest** (progression stuck).

### Activity
Aktif AFK contract snapshot + progress.
Aksiyon: **force-reveal** (server bug workaround), **cancel contract** (yield kayıp).

### Guild
Guild membership + rank + contribution.
Aksiyon: **remove from guild** (destructive; audit).

### Party
Aktif party info + role.

### Achievements
Unlocked list + points.
Aksiyon: **grant achievement** (support restoration; audit; opsiyonel Phase 02).

### Mail
Inbox list (subject görünür; body full read `player.pii_reveal` gerekli).
Aksiyon: **delete mail** (support; attachment iade).

### Chat Log
Karakterin gönderdiği chat mesajları (whisper + guild + party — kendi kanalları).
Retention 30 gün.
Aksiyon: yok (read-only forensic).

### Ledger
Ekonomi ledger karaktere ait entries; filter kind + tarih range.

### Audit
Karaktere yapılmış admin aksiyonların history'si.

## 4. Freeze / Restore

**Freeze**: karakter tüm mutation reddedilir (login OK, hareket YOK). Kullanım: ekonomi drift tespit → onay bekleyene kadar donduruma.
```
POST /internal/admin/character/{id}/freeze { reason }
POST /internal/admin/character/{id}/restore { reason }
```

## 5. Item Restore

Kullanıcı destroy ettiği (v1'de delete yok — item state history retain) veya bug ile kaybolan item için:
- Item instance UUID biliniyorsa `POST /internal/admin/character/{id}/item/restore { item_instance_id, reason }`.
- Yeni instance oluşturmak (item template yaratımı) support l2 için **yasak**; sadece existing instance state restore.

## 6. Endpoint (Ana)

- `GET /internal/admin/character/{id}` compound.
- `GET /internal/admin/character/{id}/inventory`
- `GET /internal/admin/character/{id}/equipment`
- `GET /internal/admin/character/{id}/quests`
- `GET /internal/admin/character/{id}/activity`
- `GET /internal/admin/character/{id}/ledger?kind=&from=&to=`
- `GET /internal/admin/character/{id}/chat?channel=&limit=&cursor=`
- `POST /internal/admin/character/{id}/freeze` / `restore`
- `POST /internal/admin/character/{id}/item/restore`
- `POST /internal/admin/character/{id}/quest/{code}/reset`

## 7. Cross-Reference
- [character](../character/README.md), [inventory](../inventory/README.md)
- [ops.character.freeze — DD ekonomi reconciler'la ilgili]
- [economy/transaction.md](../economy/transaction.md)

## 8. Bekleyen Kararlar
- Achievement grant support scope (v1 yok; Phase 02 opt).
- Chat log tam retention (30d) mı uzatılsın (Phase 02).
- Item restore rate limit / dual approval Phase 02.
