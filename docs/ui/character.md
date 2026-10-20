# Character Select / Create / Delete

## 1. Amaç
Login sonrası karakter seçim; yeni karakter yaratım; silme (48h grace).

## 2. Ekranlar

- **Character Select** — sağa/sola swipe / grid; slot boşsa "+ New Character".
- **Character Create Wizard** — 4 adım: race → class → appearance → name.
- **Character Delete Confirm** — modal + typed confirm ("KARAKTERİ SİL" yaz).

## 3. Layout — Character Select (Mobile)

```
┌──────────────────────┐
│  ⚙  Karakterler   ⇄  │
├──────────────────────┤
│                      │
│    [ 3D preview ]    │
│      Grim (lvl 42)   │
│      Human Warrior   │
│      Ironhold        │
│                      │
│   ● ● ○   (3 slot)   │
│                      │
├──────────────────────┤
│  [ OYUNA GİR ]       │
│  Sil  |  Yeniden Adlandır │
└──────────────────────┘
```

Slot cap v1: 2 karakter/account (DD-OPEN-8 default 2).

## 4. Create Wizard

**Step 1 — Race**
- Grid 6-8 race (v1 subset); tap → detail panel (base stats + passive).
- "Devam" enabled after selection.

**Step 2 — Class**
- Class grid (WARRIOR/MAGE/HUNTER v1 subset), tap → spec preview + skills.
- Specialization choice **level 10'da**; şimdilik base class seçimi.

**Step 3 — Appearance**
- Sex + preset (art asset içerik pack limitli v1); Phase 02 detaylı customization.

**Step 4 — Name**
- Name validator: 3-16 karakter, unique.
- `POST /v1/character` → success → character select ekranına dön.

## 5. Delete

`POST /v1/character/{id}/delete` (`{"confirm_name": "<karakter adı>"}`, büyük/küçük harf duyarsız; eşleşmezse 400 `E_NAME_MISMATCH`, DD-282) → karakter eşyalarıyla birlikte hemen silinir (DD-283, bekleme yok). `delete_grace` yalnız hesap silmede (DD-231) kalır; `POST .../delete-cancel` onu geri alır.

UI: karakter kartında sil butonu → ad yazılana kadar pasif onay penceresi; bekleyen (hesap silme) kartında silinme zamanı + iptal butonu.

## 6. State

- Character list provider — GET `/v1/character` on mount.
- Create state — wizard step; local until submit.
- 3D preview: web view / lightweight 3D (Phase 02); v1 statik illustration + rig pose.

## 7. API

- `GET /v1/character` — list.
- `POST /v1/character` — create.
- `POST /v1/character/{id}/delete`
- `POST /v1/character/{id}/delete-cancel`
- `POST /v1/character/{id}/enter` — oyuna giriş (session bind character).

## 8. Error Handling

| Kod | UI |
|---|---|
| `E_SLOT_LIMIT` | Modal "Slot dolu; slot açmak için ayarlar" |
| `E_NAME_TAKEN` | Inline "Bu isim alınmış" |
| `E_NAME_INVALID` | Inline "3-16 karakter, harf ve rakam" |

## 9. Metrikler
- `ui_character_create_started_total`
- `ui_character_create_complete_total`
- `ui_character_delete_requested_total`

## 10. Cross-Reference
- [character/character.md](../character/character.md), [race](../race/README.md), [class](../class/README.md)
- [auth/session.md](../auth/session.md)
