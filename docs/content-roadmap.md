# İçerik Geliştirme Yol Haritası

Kaynak: [content-audit.md](content-audit.md) (`python scripts/content_audit.py`, 23 bulgu). Bu belge öneridir; her adım uygulanırken DD açılır.

## Durum (2026-10-20)

Denetim bulguları 23'ten 9'a indi.

| Adım | Durum | DD |
|---|---|---|
| K1 meslek ölçeği, K2 kumaş/deri/plaka, K3 sınıf başı iki silah ailesi, K4 bölgesel av | bitti | 271–274 |
| A — üretim zinciri: 1645 ekipman tarifi, 3 yeni meslek, yüzdeli iyileştirme, deri ganimeti ölçeği | bitti | 274 |
| B — nadirlik ve ekler güç verir | bitti | 275 |
| C — takılar, Koruyucu kalkanları, Kuyumculuk | bitti | 276 |
| D — baskın setleri (18 set) | bitti | 277 |
| E4 — görev zincirleri Sv 1000'e | bitti | 278 |
| D3 — solo zindan setleri | açık | |
| D4 — mitik parça ve benzersiz etki | açık (tasarım kararı gerekir) | |
| E1 — canavara özgü malzemeler | açık (önce bir kullanım alanı gerekir) | |
| E2 — komşu kademe şansı | açık (düşük öncelik) | |
| E5 — bölge başına daha çok canavar | açık (düşük öncelik) | |
| F — soket, büyü, üretim kalitesi | açık (Faz 02) | |

## Temel sorun

Seviye 1–1000 dünyası (140 canavar, 31 zindan/baskın, 706 ekipman) tamam. Ama onu besleyen ekonomi Sv 10–100 civarında kalmış:

- **Üretim** Sv 10'da bitiyor.
- **Toplama meslekleri** 100'de bitiyor.
- **Deri ganimeti** Sv 25'te sabitleniyor.
- **İksirler** Sv 28'de bitiyor.

Ekipman tarafında:

- 12 slottan yalnız 7'si dolu.
- Ekler ve nadirlik güç vermiyor.
- Set yok.

Sonuç olarak Sv 10 sonrasında oyuncunun tek ilerleme yolu %4 ihtimalli ekipman ganimeti.

## Önce verilecek kararlar

| # | Karar | Seçenekler | Öneri |
|---|---|---|---|
| K1 | Meslek seviyesi ↔ karakter seviyesi | (a) meslekleri 1–1000'e genişlet; (b) meslek 1–100 = karakter 1–1000 (×10) | **(b)**: mevcut 150 kaynak ve 8 kademe korunur; her meslek seviyesi 10 karakter seviyesine denk gelir. Az içerik gerektirir. |
| K2 | Zırh türleri (kumaş/hafif/orta/ağır) | (a) zırhı 4 türe böl (zırh şablonu ×4); (b) türleri izinlerden kaldır, zırh sınıfsız kalsın | **(b)** şimdilik: tutarlı ve sade. Sınıf kimliği silah, takı ve setten gelir. (a) daha sonra seçilebilir. |
| K3 | Boş silah aileleri (21/32) | (a) hepsine içerik; (b) sınıf başı +1 aile, kalanları izinden çıkar | **(b)**: ör. savaşçıya balta, büyücüye asa+küre, okçuya arbalet. |
| K4 | Av bölgeye bağlı mı | (a) yalnız bulunulan bölge ve komşuları; (b) serbest kalsın | **(a)**: yolculuk ve dünya anlam kazanır. Sv 1–25 canavarlarına bölge atanır. UX'e sürtünme ekler, bu yüzden kararı oyuncu deneyimine göre ver. |

## Faz A — Üretim zincirini bağla (1 sprint, çoğu içerik)

1. **Demircilik külçeye bağlanır.**
   - Tarifler ham cevher yerine külçe ister.
   - Ergitme artık çıkmaz sokak değil.
2. **Demircilik tüm ekipman merdivenini kapsar.**
   - `gen_gear_ladder.py` her kademe için tarif bloğu da üretir: silah + 5 zırh, girdiler o kademenin külçesi + kütük + deri.
   - K1'e göre beceri = ⌈Sv/10⌉.
   - Demircilik tavanı 100 olur.
   - Ganimet olmadan da ilerleme sağlanır; üretilen eşya sıradan nadirlikte kalır, ganimet daha nadir olabilir.
3. **Deri ganimeti ölçeklenir.**
   - `CombatHides` erişimi = canavar Sv / 10 (K1).
   - Sv 1000 canavarı en iyi deriyi, Sv 200 canavarı Sv 20 hayvanını düşürür.
4. **Kullanılmayan kaynaklara kullanım verilir.**
   - Deri, post ve pul: zırh tariflerinde.
   - Kütük: silah sapı, yay, asa.
   - Ot: iksir ve yemek.
   - Balık: kalan 12 balık için yemek.
   - Taş: Faz C'deki kuyumculuk.
   - Hedef: 66 kullanılmayan kaynak → 0.
5. **İyileştirme yüzdeye geçer.**
   - Yeni `heal_pct` etkisi; iksir ve yemek azami canın yüzdesini yeniler.
   - Tek tablo her seviyeye ölçeklenir.
   - Üst kademe iksirler daha yüksek yüzde verir.
6. **Simya beceri 30–60 aralığı** dolar (direnç/hasar iksirleri yemek bonusu modeliyle). K1 (b) seçilirse simya tavanı 100 olur.

**Doğrulama:** denetimde "kullanılmayan kaynak", "külçe" ve "son tarif" bulguları sıfıra iner. `TestBalanceTargets` geçer.

## Faz B — Ekler ve nadirlik güç versin (1 sprint, kod)

1. **Ek kataloğu** (`docs/item-engine/affix.md`): `content/affixes.yaml`.
   - Her ön/son ek bir modifier verir (ör. `sharp` +%x saldırı, `of_the_bear` +güç).
   - Ek büyüklüğü eşya kademesiyle ölçeklenir.
   - Mevcut eşyalardaki ek adları korunur; yeni kod gerekmez, adlar katalogdan çözülür.
2. **Nadirlik taban stat çarpanı:** seçkin ×1.05 … efsanevi ×1.25.
   - Mevcut ekipman bonusu yolundan (`EquippedFetcher`) uygulanır.
3. **Denge testine nadirlik eklenir:** ortalama ganimet beklentisiyle sheet hesaplanır.

**Neden erken:** en çok oynanan döngü (av → ganimet) bugün yalnız satış değeri üretiyor. Bu adım az kodla en büyük ilerleme hissini verir.

## Faz C — Boş slotlar: takı + yan el (1–2 sprint)

1. **Takı merdiveni.**
   - Kolye, yüzük ×2, tılsım ×2: kademe başı 3 şablon (yüzük ve tılsım iki slota takılır) → 47 × 3 = 141 şablon.
   - Statlar: kolye can + ana stat, yüzük kritik/isabet, tılsım kaçınma/hasar azaltma.
   - `gen_gear_ladder.py` ile üretilir.
2. **Kuyumculuk mesleği** (tavan 100): taş + külçe + değerli taş → takı. Taşlar böylece kullanım kazanır.
3. **Yan el.**
   - Kalkan (savaşçı/paladin/rahip) ve odak/kitap (büyücü/rahip/druid/şaman) merdiveni: 2 aile × 47.
   - Titan Siper Kalkanı uzmanlık eşyası olarak kalır.
4. **K3 sonucu silah aileleri** eklenir ya da izinlerden çıkarılır.

**Doğrulama:** denetimde boş slot sayısı sıfır olur. Denge simülasyonu yeni güç bütçesiyle yeniden ayarlanır (mob eğrisi ya da ekipman çarpanı).

## Faz D — Setler (1–2 sprint)

1. **Set motoru** (`docs/item-engine/set.md`).
   - `content/item_sets.yaml`: parça listesi, 2/4/5 parça bonusları.
   - Sheet modifier kaynağı (yemek bonusuyla aynı yol).
   - İstemcide set paneli.
2. **İlk içerik:** 6 baskın seti (Sv 100/300/600/700/850/1000), her biri 5 zırh + 1 takı.
   - Yalnız baskın boss'undan düşer.
   - Sınıf rolüne göre 2 varyant: savunma ve hasar.
3. **Zindan setleri:** bölge başı 1 set, 25 zindan.
   - Zindan boss'u setten parça düşürür.
   - Böylece "zindana özgü ganimet yok" bulgusu kapanır.
4. **Mitik:** baskın boss'unda %1 mitik set parçası + benzersiz etki (`unique-effect.md`), en sonda.

## Faz E — Ganimet ve dünya derinliği (1–2 sprint)

1. **Canavar/arketip ganimet tabloları.**
   - Boss mührü, ejder pulu gibi özel malzemeler.
   - Kullanım: set yükseltme ve kuyumculuk.
   - `gen_world.py` bölge başı 1 özel malzeme üretir.
2. **Komşu kademe:** %15 bir alt, %5 bir üst kademe (elit ve şef canavarlarda). Kademe sınırında ani değişim yumuşar.
3. **K4 (a) seçilirse:** av bölge kapısı + Sv 1–25 canavarlarına bölge.
4. **Görevler** Sv 10'da bitiyor: `gen_world.py` bölge başı 3–5 görevlik zincir üretir (54 bölge).
5. Bölge başına 6 canavar az; Sv 100 sonrası bölge başı +2 canavar düşünülebilir (düşük öncelik).

## Faz F — Eşya motoru tasarımları (Faz 02)

Tasarımı hazır, uygulanmamış: soket + mücevher, büyü, üretim kalitesi / kritik üretim. Setler ve takılar oturduktan sonra ele alınmalı. Her biri yeni bir güç ekseni ekler ve dengeyi yeniden açar.

## Sıra özeti

| Sıra | Faz | Etki | Maliyet |
|---:|---|---|---|
| 1 | A — üretim zinciri + yüzde iyileştirme | Sv 10 sonrası ilerleme yolu, 66 kaynak kullanıma girer | Düşük (üreteç + içerik) |
| 2 | B — ek/nadirlik gücü | Ganimet anlamlı olur | Düşük–orta |
| 3 | C — takı + yan el + kuyumculuk | 5 boş slot dolar | Orta (içerik + denge) |
| 4 | D — setler | Zindan/baskın hedefi | Orta (motor + içerik) |
| 5 | E — ganimet tabloları, bölge, görev | Dünya tutarlılığı | Orta |
| 6 | F — soket, büyü, kalite | Derinlik | Yüksek |

Temizlik (her zaman): 6 emekli eşya (`potion_mp_small`, `linseed`, `mandrake`, `king_bloom`, `elm_log`, `yew_log`) oyuncu envanterlerinde yoksa migration ile silinebilir.
