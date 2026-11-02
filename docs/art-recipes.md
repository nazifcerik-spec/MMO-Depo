# Görsel üretim tarifleri

Eksik görseller (`docs/art-missing.md`) mevcut setlerin tarzına bağlanır; yeni bir tarz açılmaz. Tarifler
`scripts/art_sets.py` içinde kod olarak durur: `python scripts/art_sets.py prompt <set> <sayfa>` Gemini komutunu
basar, `cut` sayfayı keser ve `client/assets/<klasör>/<kod>.webp` yazar. Her komut "yazı, etiket, logo yok" diye biter.

## Ortak tarz çapaları

| Çapa | Nerede kullanılıyor | Kurallar |
|---|---|---|
| **Madalyon** | sınıf/terfi amblemleri, bölgeler, meslekler | Hücre başı yuvarlak madalyon, süslü altın kenar, neredeyse siyah zemin; boyalı yüksek fantezi ikon; küçükken okunur iri şekiller, canlı renk, parıltı ve kenar ışığı. 3×3 sayfa, 160 px. |
| **Eşya ikonu** | ekipman, malzeme, parşömen | Hücre başı tek eşya, hafif çapraz, düz koyu arduvaz-mavi zemin + yumuşak vinyet; cilalı mobil MMORPG envanter ikonu, temiz siluet, yumuşak cel gölge, keskin kenar ışığı. 5×5 (ya da 3×3), 128 px. |
| **İksir şişesi** | mevcut can/mana iksirleri | Saf siyah zemin, mantar tıpalı cam şişe, altın çerçeve, parlayan sıvı; kademe arttıkça şişe büyür, kristalleşir, altın süs artar. 128 px. |
| **Portre** | avatarlar, yardımcılar | Baş-omuz, yüz ortada, koyu sade vinyet, sıcak ışık; her kişi yaş/ırk/yüz/kıyafetle ayrışır. 192 px. |
| **Pet** | pet evreleri + yumurta | Açık gri zemin, cel gölgeli yaratık, element parıltısı, yer gölgesi; 3 evre tek posterde yan yana (aynı yaratık kimliği), yumurta ayrı ve tek başına. 256 / 128 px. |

## Eksik alanlar ve tarifleri

| Alan | Eksik | Set | Çapa | Sayfa | Tarif özeti |
|---|---:|---|---|---:|---|
| Toplama meslekleri | 5 | `gather_professions` | Madalyon | 1 | Mesleğin aleti: kazma + cevher, orak + ot demeti, balta + kütük, olta + balık, yay + pati izi. |
| Sınıf portreleri | 10 | `class_portraits` | Avatar portresi, boy portresi | 10 | Dikey 9:16 kahraman çizimi, diz üstü, sınıfın renk paleti ve silahı; oluşturma ekranı için üst 3:4 kesilir. |
| Pet yumurtası + evreleri | 12 + 33 | Excel listesi (`pet_prompt.py`) | Pet | 23 | Excel komutu; poster 3 evre yan yana, yumurta "yeni ve ayrı görsel, içinde yaratık görünmez". |
| Zindan + baskın | 22 | `dungeons` | Madalyon | 3 | Madalyon içinde boss'un kendi inindeki sahnesi. |
| İksirler | 6 | `potions` | İksir şişesi | 1 | Kırmızı = can, mavi = mana; Supreme en büyük, kanatlı altın süslü kristal. |
| Yemekler | 30 | `foods` | Eşya ikonu | 2 | Tabaklanmış balık yemeği; Sv 1'de tahta tabak, merdiven yükseldikçe gümüş → altın tabak, son kademelerde hafif parıltı. |
| Görevler | 62 | `quests` | Madalyon | 7 | Madalyon içinde görevin konusu: bölge yolu + canavar ya da ödül avının hedefi. |
| Canavarlar | 146 | `mobs` | Madalyon | 17 | Canavarın baş/gövde portresi, kendi arazisinde. |
| Yetenekler | 471 | `skills` | Eşya ikonu | 19 (5×5) | Adın kelimelerinden nesne motifi + etki görünümü + sınıf paleti; ad isteme girmez (Gemini yazıyor). |

## Üretim sırası

En az eksikten başlanır; bir sayfa ~1 Gemini üretimi, günlük kota ~35–40.

1. ~~Pet 39 yumurtası + pet 40–50~~ (bitti, 150 evre + 50 yumurta)
2. ~~Toplama meslekleri → iksirler → sınıf portreleri~~ (bitti)
3. ~~Zindanlar → yemekler → görevler~~ (bitti)
4. ~~Canavarlar (17) → yetenekler (19)~~ (bitti; yetenekler ad yerine adın kelimelerinden motifle — `SKILL_MOTIF`)

Her sayfa kesildikten sonra `python scripts/art_audit.py` listeyi günceller; `client/test/art_sets_test.dart` bitmiş setleri denetler.
