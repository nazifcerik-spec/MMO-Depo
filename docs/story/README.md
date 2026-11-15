# MeNa World — Hikâye: "Yedi Parçalı Taç"

Bu klasör oyunun ana hikâyesidir. Kod içermez; oyun içi Hikâye ekranına (`content/story.yaml`, DD-313) ve görevlere (`content/quests.yaml`) aktarılmak üzere yazıldı.

Hikâye ikinci şahısla anlatılır ("sen"). Oyuncu hangi ırk ve sınıftan olursa olsun aynı metni okur; ırka ya da sınıfa özel küçük satırlar ayrıca işaretlidir.

## Dosyalar

| Dosya | Bölüm | Seviye | Kıta |
|---|---|---|---|
| [00-onsoz.md](00-onsoz.md) | Önsöz: Külden Uyanış | 1 | Sundermere, Demirkale |
| [01-sundermere.md](01-sundermere.md) | I. Bölüm: Kıvılcımın Taşıyıcısı | 1–100 | Sundermere |
| [02-kirik-tac.md](02-kirik-tac.md) | II. Bölüm: Taçsız Ülke | 100–300 | Kırık Taç Krallığı |
| [03-ejdersirti.md](03-ejdersirti.md) | III. Bölüm: Ejderlerin Yası | 300–600 | Ejdersırtı Adaları |
| [04-yildizsiz.md](04-yildizsiz.md) | IV. Bölüm: Saati Durduran | 600–700 | Yıldızsız Topraklar |
| [05-gok-yarigi.md](05-gok-yarigi.md) | V. Bölüm: Gökteki Yara | 700–850 | Gök Yarığı |
| [06-boslugun-kiyisi.md](06-boslugun-kiyisi.md) | VI. Bölüm: Son Kalenin Bekçisi | 850–1000 | Boşluğun Kıyısı |
| [07-sonlar.md](07-sonlar.md) | Sonlar ve Sonsöz | 1000 | Ebediyet Kıyısı |

## Hikâyenin özü

Dünyanın adı **Mena**: toprağın, denizin ve yıldızların içinden akan canlı bir nefes. Bin yıl önce bu nefesin bekçisi tanrı **Aethon**'du. Başında, yıldızları gökte tutan **Yıldız Tacı** vardı.

Aethon ölümlü bir kadını, şifacı **Neris**'i sevdi. Neris öldüğünde Aethon ölümü durdurmak istedi ve zamanı dondurmaya, dünyayı tek bir ana çivilemeye kalktı. Başmelek **Uriel** ve kardeşleri onu tahtından indirdi. Taç yedi parçaya bölündü, parçalar dünyaya düştü:

- Parçaların düştüğü yerlerde krallıklar çöktü (Kırık Taç).
- Yıldızlar söndü (Yıldızsız Topraklar).
- Gök yarıldı (Gök Yarığı).
- Aethon dünyanın kenarına, Ebediyet Kıyısı'na zincirlendi.

Oyuncu, Yeşilorman'ın kıyısına düşen yedinci ve en küçük parçayı, **Kıvılcım**'ı avucunda taşıyarak uyanır. O andan sonra her parça onu çağırır. Kıvılcımın içinde bir ses konuşur.

## Tercih sistemi

Tercihler üç tür iz bırakır.

### 1. Parça tercihleri (hikâyenin sonunu belirler)

Her bölümün sonunda bir Taç Parçası ele geçer. Oyuncu parçayla ne yapacağına karar verir:

| Tercih | Anlamı | Oyun içi his |
|---|---|---|
| **Mühürle** | Parçayı Uriel'in mührüyle kilitle. Tanrı zincirde kalsın. | Düzen, güvenlik, durgunluk |
| **Yut** | Parçanın gücünü Kıvılcıma kat. Güç senin olsun. | Hırs, güç, yalnızlık |
| **Bırak** | Parçayı kırıp ışığını Mena'ya geri ver. | Özgürlük, kayıp, yeniden doğuş |

Altı parça vardır (Varn, Aldric, Azhorath, Kaedrin, Uriel, Morvath). Yedinci parça Kıvılcımın kendisidir ve son sahnede verilir. Sayaçlar: `muhur`, `yut`, `birak` (0–6).

### 2. Merhamet

Bazı tercihler `merhamet +1` verir: bir düşmanı bağışlamak, bir laneti kırmak, zor yolu seçmek. Hikâye boyunca 18 merhamet fırsatı vardır. Gizli son en az 8 merhamet ister; bazı geçitler (Uriel, Leviathan) de merhamete bakar.

### 3. Kişi ve dünya bayrakları

Bir NPC'nin yaşayıp yaşamadığı, kimin yanında durduğu, hangi hizibe katıldığın. Bunlar sonraki sahneleri, yardıma gelenleri ve sonsözü değiştirir. Tam liste aşağıda.

**Oyun kuralıyla uyum:** Zindan ve baskın boss'ları oyunda tekrar tekrar yenilir. Hikâyede bir boss'u ilk yenişin bir "düello" ya da "sınav" olarak anlatılır; bağışlama tercihi ancak bu ilk yenilgiden sonra gelir. Böylece boss'un tekrar oynanabilmesi hikâyeyle çelişmez (yankı, anı, bekçi büyüsü gibi gerekçeler bölümlerde verilmiştir).

## Bayrak tablosu

| Bayrak | Değerler | Belirlendiği yer | Etkilediği yer |
|---|---|---|---|
| `ses_sir` | `anlatti` / `sakladi` | 0.3 | 1.4, 1.7, 5.6 |
| `mirel` | `yoldas` / `serbest` / `tutuklu` | 1.1 | 1.3, 1.9, 2.1, 2.4, 2.7, 6.5, sonsöz |
| `taraf` | `demirkale` / `yesilorman` / `kor_ahdi` | 1.4 | 1.5, 2.3, 3.3, 3.4, 3.5, 6.5 |
| `ihra` | `yoldas` / `yok` | 1.5 | 1.6, 3.4, 5.3, 5.4, 6.2, 6.4, gizli son |
| `kul_olusu` | `arindi` / `yok_edildi` | 1.5 | 6.2 |
| `duvar_resmi` | `alindi` / `birakildi` | 1.6 | 3.4, 5.5 |
| `hedda_klan` | `kurtarildi` / `yikildi` | 1.7 | 2.5, 6.5, sonsöz |
| `isolde` | `hukumdar` / `surgun` | 1.8 | 2.2, 2.5, 2.6, 6.5, sonsöz |
| `varn_ruhu` | `huzurda` / `bagli` | 1.8 | 1.9 (baskın), 6.3 |
| `ova` | `asiler` / `sadiklar` / `tarafsiz` | 2.2 | 2.6, 6.5 |
| `ana_agac` | `kurtarildi` / `yandi` | 2.3 | 3.3, 3.5, 6.5, sonsöz |
| `seth_ka` | `ozgur` / `bagli` | 2.4 | 2.5, 6.1, 6.2, sonsöz |
| `tac_varisi` | `mirel` / `orrin` / `bos` | 2.7 | 2.8, 6.5, sonsöz |
| `brannoc` | `yasiyor` / `oldu` | 2.6 | 6.5, sonsöz |
| `yssa` | `dost` / `rakip` | 3.1 | 3.7, 6.1, 6.5, sonsöz |
| `yumurta` | `kurtarildi` / `satildi` / `kirildi` | 3.2 | 3.6, 3.7, 3.8, 6.5, sonsöz |
| `sylwen` | `uyudu` / `uyandi` | 3.3 | 3.6, 6.5, sonsöz |
| `ani_kristali` | `alindi` / `kirildi` | 3.4 | 5.4, 5.5, gizli son |
| `ansel` | `durduruldu` / `kazandi` / `kurtarildi` | 3.5 | 6.5, sonsöz |
| `ilvar` | `yasiyor` / `kayip` | 4.2 | 4.6, 5.6, 5.7, 6.1, sonsöz |
| `pip` | `pip` / `kaedrin` / `ikisi_bir` / `yok` | 4.5 | 4.6, 5.1, 5.7, 6.5, 6.6, sonlar, sonsöz |
| `zephyros` | `affedildi` / `dustu` | 5.2 | 5.7, 6.5, sonsöz |
| `kaos_ana` | `teselli` / `parcalandi` | 5.4 | gizli son |
| `uriel` | `muttefik` / `dusman` | 5.6 | 5.7, 6.5, 6.6, sonlar |
| `leviathan` | `gecit` / `savas` | 6.1 | 6.2, sonsöz |
| `morvath` | `huzurda` / `yikildi` | 6.4 | 6.6, sonlar, sonsöz |

## Önemli kişiler

### Yoldaşlar ve dostlar
| Kişi | Irk | Kim | İlk görünüş |
|---|---|---|---|
| **Hedda Taşelli** | Cüce | Demirkale'nin baş demircisi. Seni ormanın kıyısında bulan, sert ama sıcak akıl hocası. | 0.1 |
| **Pip Çarkdişi** | Gnom | Kendinden büyük icatlar taşıyan, çok konuşan mucit. Göğün yarığını "tamir etmek" istiyor. | 1.3 |
| **Mirel Kuşdili** | İnsan | Sokakta büyümüş genç bir hırsız. Kanında kimsenin bilmediği bir taç var. | 1.1 |
| **Ihra** | Hortlak | Bin yıllık bir ölü; Aethon'un tapınağında rahibeydi. Her şeyi hatırlıyor ama hiçbir şeyi anlatmak istemiyor. | 1.5 |
| **Kaptan Yssa Kızılpul** | Drakari | Fırtına Kıyısı'nın korsan kaptanı. Altına, rüzgâra ve kendine inanır. | 3.1 |
| **Gözlemci Ilvar** | Karanlık Elf | Sönmüş yıldızların haritasını çıkaran yaşlı astronom. | 4.1 |

### Hizip önderleri
| Kişi | Irk | Hizip | İstediği |
|---|---|---|---|
| **Komutan Brannoc Ağırkalkan** | İnsan | Demirkale Muhafızları | Kıvılcımı kasaya kilitlemek. Şehri korumak, ne pahasına olursa olsun. |
| **Liora Yaprakses** | Elf | Yeşilorman Bekçileri | Kıvılcımı ormanın köklerine geri vermek. |
| **Başrahip Ansel** | İnsan | Kor Ahdi | Kıvılcımla Ateş Lordu Ignar'ı uyandırıp "dünyayı yeniden ısıtmak". |

### Tarihin kişileri ve boss'lar
| Kişi | Nerede | Kim |
|---|---|---|
| **Düşmüş Lord Varn** | Kuzey Kalesi | İlk parçayı korumak için ölümsüzlüğe yemin etmiş lord. |
| **Isolde Varn** | Kuzey Kalesi | Varn'ın kızı; babasının yeminini bozmak isteyen genç kadın. |
| **Asi Kral Orrin** | Asiler Kampı | Kralsız ülkenin kendini kral ilan etmiş asisi. Zalim mi, kurtarıcı mı? |
| **Sera Boynuzyürek** | Taç Ovaları | Orrin'in tauren yardımcısı. |
| **Firavun Seth-Ka** | Gömülü Saray | Bin yıldır kendi mezarında uyanık bir kral. |
| **Mezar Lordu Malachar** | Kraliyet Kriptası | Son kralın kardeşi; ölümü kandırmış bir prens. |
| **Kral Aldric, Taçsız Kral** | Taç Mezarlığı | İkinci parçayı başına geçirmiş son kral. |
| **Amiral Drask** | Batık Filo | Fırtınaya meydan okuyup batan amiral. |
| **Ana Ejder Vyrna** | Kuluçka Mağarası | Son ejder yumurtalarının annesi. |
| **Ejder Kral Azhorath** | Ejder Tahtı | Vyrna'nın eşi; üçüncü parça onu delirtmiş. |
| **Ruh Kraliçesi Sylwen** | Kayıp Koru | Liora'nın soyunun atası; ormanın rüyası. |
| **Kristal Titan** | Kristal Kalp | Aethon'un anılarını taşıyan canlı kristal. |
| **Ateş Lordu Ignar** | Volkan Kalbi | Kor Ahdi'nin taptığı ateş varlığı. |
| **Yıldız Yiyen Omh** | Söndürülmüş Gözlemevi | Sönen yıldızların açlığından doğmuş varlık. |
| **Zaman Lordu Kaedrin** | Saat Kulesi | Zamanı durdurmuş büyücü. Yüzü tanıdık. |
| **Fırtına Kralı Zephyros** | Fırtına Kalesi | Uriel'in eski generali. |
| **Kaos Ana** | Yarık Çekirdeği | Aethon'un yasından doğmuş canavar. |
| **Başmelek Uriel** | Göksel Kubbe | Aethon'u zincire vuran gardiyan. |
| **Leviathan Kral** | Batık Taht | Ölülerin ruhlarını boşluğa taşıyan deniz devi. |
| **Kıyamet Komutanı Morvath** | Son Kale | Bin yıl önceki ilk Taşıyıcı. |
| **Düşmüş Tanrı Aethon** | Ebedi Taht | Zincirdeki tanrı. |

### Yan karakterler
| Kişi | Irk | Nerede | Kim |
|---|---|---|---|
| **Nell Kazanbaş** | İnsan | Demirkale, "Paslı Kepçe" hanı | Her dedikoduyu bilen han sahibi. |
| **Grosk Tuzdişi** | Ork | Fısıltı Vadisi | Kervan tüccarı; her şeyin fiyatını bilir, değerini değil. |
| **Toran Külgözü** | Trol | Közdoruk | Dağın ruhlarıyla konuşan şaman. |
| **Usta Bodri** | Cüce | Demir Dağlar | Hedda'nın kuzeni, son madenci. |
| **Kara Mızrak Leydi Corva** | İnsan | Taç Ovaları | Sadık birliklerin komutanı. |
| **Rahibe Tamsin** | İnsan | Kor Ahdi | Ansel'in şüpheci yardımcısı. |
| **Yaşlı Moru** | Tauren | Ejder Yamaçları | Ejder çobanı, son ejder dostu. |
| **Selin Rüzgârdoğan** | Gnom | Uçan Adalar | Pip'in eski hocası. |
| **Kâtip Oen** | Hortlak | Son Kale | Morvath'ın bin yıllık günlüğünü tutan kâtip. |

## Sonlar (özet)

| Son | Koşul | Özet |
|---|---|---|
| **Zincir** | Son tercih: Mühürle | Aethon sonsuza dek zincirde. Yıldızlar sönük kalır, dünya güvende ama donuk. |
| **Yeni Taç** | Son tercih: Yut | Sen tanrı olursun. Yut sayacına göre merhametli ya da korkunç bir tanrı. |
| **Mena'nın Nefesi** | Son tercih: Bırak | Taç sonsuza dek kırılır, tanrıların çağı biter. Yıldızlar döner, büyü azalır. |
| **Son Ninni (gizli)** | Merhamet ≥ 8, `ihra = yoldas`, `ani_kristali = alindi`, `kaos_ana = teselli`, `pip` ≠ `yok` | Aethon'a Neris'in anısı geri verilir. Tanrı ölümlü olarak ölmeyi seçer; parçalar yeni yıldızlar olur. |

Ayrıntılar ve her kişinin sonsözü [07-sonlar.md](07-sonlar.md) içinde.
