# Respec (Specialization Değişim Akışı)

## 1. Amaç
Karakterin Tier I ve/veya Tier II specialization seçimini değiştirmek. Class kalıcıdır (DD-004); spec esnek katman.

## 2. Türler
- **Spec I Respec**: Tier I değişimi. Otomatik olarak Tier II de sıfırlanır (parent değiştiği için Tier II geçersiz).
- **Spec II Respec**: Sadece Tier II değişimi (aynı parent altında farklı seçenek).
- **Full Respec**: Tier I + Tier II reset; oyuncu iki spec de baştan seçer.

## 3. Cost & Cooldown

| Türü | İlk Kez | Sonraki | Cooldown |
|---|---|---|---|
| Spec II Respec | Ücretsiz | `500 * 2^n copper` (yani 5 silver, 10, 20, ...) | 24 saat |
| Spec I Respec | Ücretsiz | `2000 * 2^n copper` (yani 20 silver, 40, ...) | 72 saat |
| Full Respec (I + II) | Ücretsiz (aynı sayaç) | `Spec I cost + Spec II cost` | 72 saat |

Cost cap: 20000 gold (2000000 copper) toplam.
Cooldown paylaşılır (Spec I respec sonrası 72 saat hem I hem II kilit).

**Free respec token**: content rework, quest reward, event ödülü — cost bypass.

## 4. Endpoint

### 4.1 Respec

**POST `/v1/characters/{id}/respec-spec`**
Auth + CharacterScope.

Request:
```json
{
  "new_tier1_id": "sp-guardian",   // null → aynı Tier I korunur
  "new_tier2_id": "sp-sentinel",   // null → Tier II sadece reset (yeniden seç)
  "use_free_token": false           // true ise cost bypass; token invalidate
}
```

Server:
1. Load character FOR UPDATE.
2. Determine mode:
   - Both provided → Spec I + II change (Full).
   - Only Tier II → Spec II only.
   - Only Tier I → Spec I change; Tier II mandatorily reset.
3. Validation:
   - Character level >= 5 (Tier I) / >= 10 (Tier II).
   - New spec valid class + tier + parent.
   - Cooldown check.
   - Cost / token check.
4. Transaction:
   - Debit currency via ledger (skip if token).
   - **Remove old spec skills** (`character.skill_learn` delete for spec-granted skills).
   - **Preserve** class starter + racial + level-granted skills.
   - Update `character.spec_tier1_id`, `character.spec_tier2_id`.
   - Grant new spec skills.
   - Reset `character.skill_bar` slot 10 (signature) to new spec signature (if tier 2).
   - Recompute derived cache (`derived_cache_invalidated_at`).
   - Reset `character.rotation` to new spec default (bkz. §9).
   - `respec_count++` (schema field).
   - Set `respec_cooldown_until`.
   - Emit event `character.spec_respec`.
   - Audit event.
5. Response 200: updated character sheet.

### 4.2 Respec Preview (opsiyonel Phase 02)

**GET `/v1/characters/{id}/respec-preview?tier1=...&tier2=...`**
- Cost + cooldown + skill delta gösterir; henüz apply etmez.

## 5. Schema Eklemeleri

```sql
ALTER TABLE character.character
  ADD COLUMN respec_count SMALLINT NOT NULL DEFAULT 0,
  ADD COLUMN respec_cooldown_until TIMESTAMPTZ NULL,
  ADD COLUMN free_respec_tokens SMALLINT NOT NULL DEFAULT 0;
```

## 6. Skill Removal Logic

Spec değişince ne kalır ne gider:
- **Kalır**: class starter skills, racial skills, level-granted class skills, öğrenilmiş consumable skill'ler.
- **Silinir**: eski spec'in `passive_skill_codes` + `active_skill_unlock_codes` + `signature_skill_code`.

Query:
```sql
DELETE FROM character.skill_learn
WHERE character_id = $1
  AND skill_id IN (SELECT unnest(passive_ids || active_skill_unlock_ids || ARRAY[signature_skill_id])
                   FROM content.specialization
                   WHERE id IN ($old_tier1_id, $old_tier2_id));
```

Rank kayıp mı?
- **v1**: rank sıfır (yeniden 1). Adı korunmadığı için sorun değil.
- **Phase 02**: rank persist (aynı skill başka spec'te tekrar açılırsa eski rank).

## 7. Skill Bar Reset

Reset stratejisi (v1):
- `character.skill_bar` slot 0..9 → **korunur** eğer skill hala öğrenilmişse; silinmiş skill boşaltılır.
- Slot 10 (signature) → yeni Tier II signature ile doldurulur (Tier II değiştiyse).
- Slot 11 (racial) → dokunulmaz.

Öneri: preserved slot layout + boş slot notification "Manually assign new spec skills".

## 8. Rotation Reset

`character.rotation.priority_list` yeni spec default'una **override**.

- Default rotation content'te tanımlı: `content.specialization.default_rotation_json` (yeni alan).
- Oyuncu kendi rotation'ı bozulmuş gibi hissederse manuel edit (Rotation Editor UI).

Schema:
```sql
ALTER TABLE content.specialization
  ADD COLUMN default_rotation_json JSONB NOT NULL DEFAULT '[]';
```

## 9. Free Respec Token Kaynakları

- Content rework compensation (automatic grant).
- Story quest reward (major milestone).
- Event participation.
- Admin grant (support ticket resolution).
- Level 50 first-time cap (opsiyonel; öğrenme fırsatı 2. kez).

Token expire: 30 gün after grant (kullanılmazsa yanar).

## 10. Ledger

Respec cost:
```
economy.ledger_entry:
  reason='debit_respec'
  source_kind='character_respec'
  source_ref_id=character_id
  idempotency_key='respec:<char>:<respec_count>'
```

## 11. Cooldown Bypass Yok

Cooldown bypass yok — gold ödemek bile cooldown'ı sıfırlamaz. Amaç: spam respec engelleme (economy sink + kimlik disiplini).

Trial period (Phase 02 opsiyonel): karakter yeni bir spec seçtikten sonra 1 saat "trial" — dungeon veya PvP'de bulunmadıysa cost iade.

## 12. Fail Modları
- Cooldown → 429 `spec.respec_cooldown`.
- Insufficient funds → 402.
- No free token → normal cost path.
- Tier 2 without Tier 1 → 400.
- Same spec re-select → 400 `spec.no_change`.

## 13. Cross-Play Kısıtları
- Aktif dungeon içinde respec yasak (transaction reddedilir).
- Aktif combat → 15 sn out-of-combat sonrası respec.
- Aktif AFK activity → respec sırasında pause; sonrasında oyuncu resume.

## 14. Event
```json
{
  "_v": 1,
  "character_id": "c-...",
  "old_tier1_code": "berserker",
  "new_tier1_code": "guardian",
  "old_tier2_code": "ravager",
  "new_tier2_code": null,
  "used_free_token": false,
  "cost_copper": 500
}
```

Subscribers:
- `notification.spec_changed`.
- `party.role_reeval` (party içindeyse tank/heal/dps değişimi).

## 15. Cross-Platform
Aynı sheet; respec her platformdan tetiklenebilir.

## 16. Metrics
- `spec_respec_total{tier,paid}`.
- `spec_respec_cost_avg_copper`.
- `spec_free_token_used_total`.
- `spec_respec_cooldown_hit_total`.

## 17. Audit
- `character.spec_respec`.
- `character.spec_respec_free_token_used`.

## 18. Test Stratejisi
- Cost curve calculation.
- Cooldown enforce.
- Skill removal + regrant integrity (no orphan skills).
- Free token consume + expiry.
- Race: concurrent respec (optimistic lock).

## 19. İlişkiler
- [Class](class.md), [Specialization](specialization.md), [Placeholder Content](placeholder-content.md)
- [character/attributes.md](../character/attributes.md) — respec-stats ayrı akış
- [database/character.md](../database/character.md), [database/content-templates.md](../database/content-templates.md)
- [economy/ledger](../database/economy-market.md)

## 20. Bekleyen Kararlar
- Cost curve final (500 base + 2× exponential — simülasyon).
- Trial period feature (Phase 02).
- Spec preview endpoint (Phase 02).
- Cross-spec rank persistence (Phase 02).
