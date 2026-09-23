# Crafting

## 1. Amaç
Recipe tarama, malzeme kontrol, kuyruğa ekleme, aktif kuyruğu izleme, critical craft alerti, mastery progresi.

## 2. Layout

```
┌────────────────────────────┐
│ ← Crafting · Blacksmith    │
│ Lvl 22 · Mastery T2         │
├────────────────────────────┤
│ [Tümü][Silah][Zırh][Con]   │
│ [Filtre ▼] [Sırala ▼]      │
├────────────────────────────┤
│ ┌───────────────────────┐  │
│ │ [icon] Iron Sword     │  │
│ │ Lvl 15 · Uncommon     │  │
│ │ Malz: 3 iron / var 5  │  │
│ │ Süre: 4 dk            │  │
│ │ [Yap] [+Kuyruğa]      │  │
│ └───────────────────────┘  │
│ ...                        │
├────────────────────────────┤
│ Aktif Kuyruk (3/5 slot):   │
│ 1. Iron Sword ▓▓░░ 2dk     │
│ 2. Iron Plate — bekliyor   │
│ 3. Iron Boots — bekliyor   │
└────────────────────────────┘
```

## 3. Bileşenler

- **Header**: profession + level + mastery tier.
- **Filter/sort**: category, level, craftable-only, mastery-required.
- **Recipe card**: icon, name, level, rarity, material check (green/red), time.
- **Queue widget**: aktif slot listesi + progress bar + drag to reorder.

## 4. Craft Modal (Detail)

Tap recipe → modal:
- Malzeme listesi (var/gerek).
- Beklenen quality range (normal / masterwork / critical).
- Critical craft ihtimali (%).
- XP earn.
- Adet slider (batch craft).
- **Kuyruğa Ekle** / **Şimdi Başlat**.

## 5. Critical Craft Alert

Craft complete event `critical_craft=true` → animation + rare sound → toast "★ Kritik!". Item detail otomatik açık.

## 6. State

- Recipe provider — `GET /v1/crafting/recipes?profession=blacksmith`.
- Queue provider — `GET /v1/character/{id}/crafting/queue` + WS `crafting:{id}` events.
- Material availability — inventory join.

## 7. API

- `GET /v1/crafting/recipes?profession=&level=`
- `POST /v1/crafting/queue/enqueue` body: `{ recipe_code, quantity, idempotency_key }`
- `POST /v1/crafting/queue/{slot_id}/cancel`
- `GET /v1/character/{id}/crafting/queue`

## 8. Metrikler
- `ui_crafting_view_total`
- `ui_crafting_enqueue_total{recipe_code}`
- `ui_crafting_critical_alert_total`

## 9. Cross-Reference
- [crafting](../crafting/README.md), [professions](../professions/README.md)
- [progression/mastery.md](../progression/mastery.md)
