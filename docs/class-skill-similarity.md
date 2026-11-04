# Sınıf içi benzer yetenekler

Her temel sınıfın kendi ağacı (temel sınıf, 2 terfi, 4 uzmanlık) içinde karşılaştırma. Kaynak: `content/skills.yaml`, `content/progression.yaml`.

- **Neredeyse aynı:** aktiflerde aynı etki + aynı süre + güç farkı en çok %10 + bekleme farkı en çok 2 tur; pasiflerde aynı statlar + değer farkı en çok %25.
- **Aynı kol:** iki yetenek aynı karakterde birlikte bulunabilir (temel → terfi → uzmanlık). **Farklı kol:** bir karakter ikisini birden alamaz, kolların birbirine benzemesi demektir.
- Benzerlik sayısal değerlere göredir: yetenek adı, görseli, aktif yeteneğin tetik koşulu ve pasifin tetikleyicisi hesaba katılmaz. Pasifler sürekli etki olarak hesaplandığı için (`steady_state`) tetikleyici sayıyı değiştirmez.
- Bu liste bir bulgu dökümüdür, kusur listesi değildir: aynı koldaki aynı etki bilinçli bir üst kademe olabilir.
- Anlık görüntü: DD-351 sonrası içerik (seviyeler Sv 300 / 700 kapılarına göre).

| Sınıf | Yetenek | Neredeyse aynı çift | Aynı kol | Farklı kol |
|---|---:|---:|---:|---:|
| Savaşçı | 48 | 9 | 2 | 7 |
| Haydut | 47 | 27 | 5 | 22 |
| Korucu | 47 | 12 | 2 | 10 |
| Büyücü | 47 | 24 | 8 | 16 |
| Keşiş | 47 | 20 | 11 | 9 |
| Rahip | 47 | 7 | 2 | 5 |
| Paladin | 47 | 7 | 6 | 1 |
| Druid | 47 | 10 | 3 | 7 |
| Ozan | 47 | 14 | 7 | 7 |
| Şaman | 47 | 10 | 1 | 9 |
| **Toplam** | **471** | **140** | **47** | **93** |

## Savaşçı

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Ağır Kesik | Savaşçı, 25 | Hassas Vuruş | Silah Ustası, 725 | Hasar vuruşu: güç 190, bekleme 3 / güç 200, bekleme 3 | aynı kol: üst üste biner |
| Pervasız Vuruş | Cellat, 300 | Hassas Vuruş | Silah Ustası, 725 | Hasar vuruşu: güç 220, bekleme 3 / güç 200, bekleme 3 | aynı kol: üst üste biner |
| Kalkan Darbesi | Muhafız, 300 | Yırtma | Cellat, 400 | Hasar vuruşu: güç 160, bekleme 3 / güç 150, bekleme 4 | farklı kol |
| Kan Hücumu | Cellat, 500 | Toplanma Narası | Savaş Lordu, 725 | Saldırı güçlendirme: güç 20, 3 tur, bekleme 8 / güç 20, 3 tur, bekleme 8 | farklı kol |
| Misilleme | Demir Kale, 750 | Silah Değişimi | Silah Ustası, 750 | +%5 saldırı / +%5 saldırı | farklı kol |
| Misilleme | Demir Kale, 750 | Acı Güçtür | Gözü Dönmüş, 775 | +%5 saldırı / +%6 saldırı | farklı kol |
| Misilleme | Demir Kale, 750 | Cephanelik Ustalığı | Silah Ustası, 850 | +%5 saldırı / +%6 saldırı | farklı kol |
| Acı Güçtür | Gözü Dönmüş, 775 | Silah Değişimi | Silah Ustası, 750 | +%6 saldırı / +%5 saldırı | farklı kol |
| Acı Güçtür | Gözü Dönmüş, 775 | Cephanelik Ustalığı | Silah Ustası, 850 | +%6 saldırı / +%6 saldırı | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (10): Darbe (Savaşçı, Sv 1, güç 110); Yırtma (Cellat, Sv 400, güç 150); Kalkan Darbesi (Muhafız, Sv 300, güç 160); Ağır Kesik (Savaşçı, Sv 25, güç 190); Hassas Vuruş (Silah Ustası, Sv 725, güç 200); Pervasız Vuruş (Cellat, Sv 300, güç 220); Pervasız Saldırı (Gözü Dönmüş, Sv 800, güç 260); İnfaz (Savaşçı, Sv 75, güç 300); Son Darbe (Cellat, Sv 600, güç 350); Bin Kesik (Silah Ustası, Sv 950, güç 500)
- **Hasar azaltma (duruş/koruma)** (7): Kışkırtma (Muhafız, Sv 400, güç 15); Savunma Duruşu (Savaşçı, Sv 10, güç 30); Meydan Okuma (Savaş Lordu, Sv 800, güç 30); Kale Duruşu (Demir Kale, Sv 700, güç 40); Son Direniş (Muhafız, Sv 500, güç 50); Kalkan Duvarı (Demir Kale, Sv 725, güç 70); Kırılmaz (Demir Kale, Sv 950, güç 150)
- **Saldırı güçlendirme** (6): Kan Hücumu (Cellat, Sv 500, güç 20); Toplanma Narası (Savaş Lordu, Sv 725, güç 20); Çılgınlık (Gözü Dönmüş, Sv 725, güç 25); Harp Sancağı (Savaş Lordu, Sv 900, güç 30); Fetih Yürüyüşü (Savaş Lordu, Sv 950, güç 80); Öfke Avatarı (Gözü Dönmüş, Sv 950, güç 100)

**Aynı statı veren pasifler:**

- **saldırı** (4): Misilleme (Demir Kale, Sv 750, +%5 saldırı); Acı Güçtür (Gözü Dönmüş, Sv 775, +%6 saldırı); Silah Değişimi (Silah Ustası, Sv 750, +%5 saldırı); Cephanelik Ustalığı (Silah Ustası, Sv 850, +%6 saldırı)

## Haydut

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Sırttan Bıçaklama | Haydut, 25 | Karşı Hamle | Düellocu, 300 | Hasar vuruşu: güç 200, bekleme 3 / güç 180, bekleme 2 | aynı kol: üst üste biner |
| Savaş Ritmi | Düellocu, 300 | Bıçak Temposu | Kılıç Dansçısı, 700 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Savaş Ritmi | Düellocu, 300 | Durmak Bilmeyen Hareket | Kılıç Dansçısı, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Darbe Yağmuru | Düellocu, 400 | İkiz Kesik | Kılıç Dansçısı, 725 | Hasar vuruşu: güç 260, bekleme 4 / güç 240, bekleme 2 | aynı kol: üst üste biner |
| Ölüm İşareti | Suikastçı, 600 | Gölge Darbesi | Gece Bıçağı, 725 | Hasar vuruşu: güç 320, bekleme 4 / güç 300, bekleme 3 | aynı kol: üst üste biner |
| Pusu | Suikastçı, 300 | Darbe Yağmuru | Düellocu, 400 | Hasar vuruşu: güç 240, bekleme 4 / güç 260, bekleme 4 | farklı kol |
| Pusu | Suikastçı, 300 | İkiz Kesik | Kılıç Dansçısı, 725 | Hasar vuruşu: güç 240, bekleme 4 / güç 240, bekleme 2 | farklı kol |
| Ölüm İşareti | Suikastçı, 600 | Dönen Bıçaklar | Kılıç Dansçısı, 750 | Hasar vuruşu: güç 320, bekleme 4 / güç 300, bekleme 4 | farklı kol |
| Ölüm İşareti | Suikastçı, 600 | Kirli Numara | Hilebaz, 800 | Hasar vuruşu: güç 320, bekleme 4 / güç 300, bekleme 4 | farklı kol |
| Zehir Uzmanlığı | Zehir Ustası, 700 | Savaş Ritmi | Düellocu, 300 | +%5 saldırı / +%6 saldırı | farklı kol |
| Zehir Uzmanlığı | Zehir Ustası, 700 | Bıçak Temposu | Kılıç Dansçısı, 700 | +%5 saldırı / +%6 saldırı | farklı kol |
| Zehir Uzmanlığı | Zehir Ustası, 700 | Durmak Bilmeyen Hareket | Kılıç Dansçısı, 850 | +%5 saldırı / +%6 saldırı | farklı kol |
| Gölge Darbesi | Gece Bıçağı, 725 | Dönen Bıçaklar | Kılıç Dansçısı, 750 | Hasar vuruşu: güç 300, bekleme 3 / güç 300, bekleme 4 | farklı kol |
| Gölge Darbesi | Gece Bıçağı, 725 | Katalizör | Zehir Ustası, 800 | Hasar vuruşu: güç 300, bekleme 3 / güç 280, bekleme 4 | farklı kol |
| Gölge Darbesi | Gece Bıçağı, 725 | Kirli Numara | Hilebaz, 800 | Hasar vuruşu: güç 300, bekleme 3 / güç 300, bekleme 4 | farklı kol |
| Dönen Bıçaklar | Kılıç Dansçısı, 750 | Kirli Numara | Hilebaz, 800 | Hasar vuruşu: güç 300, bekleme 4 / güç 300, bekleme 4 | farklı kol |
| Aşındırıcı Zehir | Zehir Ustası, 775 | Savaş Ritmi | Düellocu, 300 | +%6 saldırı / +%6 saldırı | farklı kol |
| Aşındırıcı Zehir | Zehir Ustası, 775 | Bıçak Temposu | Kılıç Dansçısı, 700 | +%6 saldırı / +%6 saldırı | farklı kol |
| Aşındırıcı Zehir | Zehir Ustası, 775 | Durmak Bilmeyen Hareket | Kılıç Dansçısı, 850 | +%6 saldırı / +%6 saldırı | farklı kol |
| Katalizör | Zehir Ustası, 800 | Darbe Yağmuru | Düellocu, 400 | Hasar vuruşu: güç 280, bekleme 4 / güç 260, bekleme 4 | farklı kol |
| Katalizör | Zehir Ustası, 800 | Dönen Bıçaklar | Kılıç Dansçısı, 750 | Hasar vuruşu: güç 280, bekleme 4 / güç 300, bekleme 4 | farklı kol |
| Gece Çöküşü | Gece Bıçağı, 800 | Bıçak Çağlayanı | Kılıç Dansçısı, 800 | Hasar vuruşu: güç 380, bekleme 5 / güç 360, bekleme 5 | farklı kol |
| Katalizör | Zehir Ustası, 800 | Kirli Numara | Hilebaz, 800 | Hasar vuruşu: güç 280, bekleme 4 / güç 300, bekleme 4 | farklı kol |
| Zehir Patlaması | Zehir Ustası, 900 | Savaş Ritmi | Düellocu, 300 | +%8 saldırı / +%6 saldırı | farklı kol |
| Zehir Patlaması | Zehir Ustası, 900 | Bıçak Temposu | Kılıç Dansçısı, 700 | +%8 saldırı / +%6 saldırı | farklı kol |
| Zehir Patlaması | Zehir Ustası, 900 | Durmak Bilmeyen Hareket | Kılıç Dansçısı, 850 | +%8 saldırı / +%6 saldırı | farklı kol |
| Ölüm Dansı | Gece Bıçağı, 950 | Bıçak Fırtınası | Kılıç Dansçısı, 950 | Hasar vuruşu: güç 600, bekleme 16 / güç 550, bekleme 15 | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (15): Hızlı Bıçaklama (Haydut, Sv 1, güç 120); Karşı Hamle (Düellocu, Sv 300, güç 180); Sırttan Bıçaklama (Haydut, Sv 25, güç 200); Pusu (Suikastçı, Sv 300, güç 240); İkiz Kesik (Kılıç Dansçısı, Sv 725, güç 240); Darbe Yağmuru (Düellocu, Sv 400, güç 260); Katalizör (Zehir Ustası, Sv 800, güç 280); Gölge Darbesi (Gece Bıçağı, Sv 725, güç 300); Dönen Bıçaklar (Kılıç Dansçısı, Sv 750, güç 300); Kirli Numara (Hilebaz, Sv 800, güç 300); Ölüm İşareti (Suikastçı, Sv 600, güç 320); Bıçak Çağlayanı (Kılıç Dansçısı, Sv 800, güç 360); Gece Çöküşü (Gece Bıçağı, Sv 800, güç 380); Bıçak Fırtınası (Kılıç Dansçısı, Sv 950, güç 550); Ölüm Dansı (Gece Bıçağı, Sv 950, güç 600)
- **Sersemletme** (3): Duman Bombası (Hilebaz, Sv 750, güç 120); Böbrek Darbesi (Haydut, Sv 75, güç 140); Büyük İllüzyon (Hilebaz, Sv 950, güç 350)
- **Zayıflatma** (3): Zaafı Açığa Çıkar (Suikastçı, Sv 500, güç 25); Aldatmaca (Hilebaz, Sv 725, güç 25); Zehir Bulutu (Zehir Ustası, Sv 750, güç 30)
- **Hasar azaltma (duruş/koruma)** (2): Kaçınma (Haydut, Sv 10, güç 35); Kusursuz Ayak Oyunu (Düellocu, Sv 500, güç 40)
- **Saldırı güçlendirme** (2): Gölge Adımı (Suikastçı, Sv 400, güç 20); Düellocunun Meydan Okuması (Düellocu, Sv 600, güç 25)
- **Süreli hasar** (2): Zehirli Bıçak (Zehir Ustası, Sv 725, güç 70); Veba Getiren (Zehir Ustası, Sv 950, güç 150)

**Aynı statı veren pasifler:**

- **saldırı** (6): Zehir Uzmanlığı (Zehir Ustası, Sv 700, +%5 saldırı); Aşındırıcı Zehir (Zehir Ustası, Sv 775, +%6 saldırı); Zehir Patlaması (Zehir Ustası, Sv 900, +%8 saldırı); Savaş Ritmi (Düellocu, Sv 300, +%6 saldırı); Bıçak Temposu (Kılıç Dansçısı, Sv 700, +%6 saldırı); Durmak Bilmeyen Hareket (Kılıç Dansçısı, Sv 850, +%6 saldırı)
- **kritik hasarı + kritik** (3): Fırsatçı (Haydut, Sv 50, +%10 kritik hasarı, +%4 kritik); Ölümcül Açılış (Suikastçı, Sv 300, +%15 kritik hasarı, +%2 kritik); Gölge İnfazı (Gece Bıçağı, Sv 900, +%20 kritik hasarı, +%5 kritik)

## Korucu

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Hızlı Yaylım | Nişancı, 500 | Nişanlı Atış | Keskin Nişancı, 725 | Hasar vuruşu: güç 290, bekleme 4 / güç 320, bekleme 4 | aynı kol: üst üste biner |
| Hızlı Yaylım | Nişancı, 500 | Ezici Ok | Arbalet Ustası, 725 | Hasar vuruşu: güç 290, bekleme 4 / güç 300, bekleme 3 | aynı kol: üst üste biner |
| Delici Atış | Nişancı, 300 | Saldırı Emri | Hayvan Terbiyecisi, 725 | Hasar vuruşu: güç 240, bekleme 3 / güç 260, bekleme 2 | farklı kol |
| Nişanlı Atış | Keskin Nişancı, 725 | Ezici Ok | Arbalet Ustası, 725 | Hasar vuruşu: güç 320, bekleme 4 / güç 300, bekleme 3 | farklı kol |
| Kritik Odak | Keskin Nişancı, 750 | Hızlı Doldurma | Arbalet Ustası, 750 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 7 | farklı kol |
| Kritik Odak | Keskin Nişancı, 750 | Hayvani Kükreme | Hayvan Terbiyecisi, 750 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | farklı kol |
| Hızlı Doldurma | Arbalet Ustası, 750 | Hayvani Kükreme | Hayvan Terbiyecisi, 750 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 7 / güç 25, 3 tur, bekleme 8 | farklı kol |
| Ölümcül Hassasiyet | Keskin Nişancı, 800 | Kuşatma Oku | Arbalet Ustası, 800 | Hasar vuruşu: güç 400, bekleme 5 / güç 380, bekleme 5 | farklı kol |
| Ölümcül Hassasiyet | Keskin Nişancı, 800 | Sürü Saldırısı | Hayvan Terbiyecisi, 800 | Hasar vuruşu: güç 400, bekleme 5 / güç 360, bekleme 4 | farklı kol |
| Kuşatma Oku | Arbalet Ustası, 800 | Sürü Saldırısı | Hayvan Terbiyecisi, 800 | Hasar vuruşu: güç 380, bekleme 5 / güç 360, bekleme 4 | farklı kol |
| Uzak Menzil Ustalığı | Keskin Nişancı, 850 | Hazırlanmış Zemin | Tuzakçı, 700 | +%6 saldırı / +%5 saldırı | farklı kol |
| Çelik Yağmuru | Arbalet Ustası, 950 | Vahşetin Çağrısı | Hayvan Terbiyecisi, 950 | Hasar vuruşu: güç 560, bekleme 15 / güç 520, bekleme 15 | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (13): Hızlı Atış (Korucu, Sv 1, güç 115); Güçlü Atış (Korucu, Sv 25, güç 210); Delici Atış (Nişancı, Sv 300, güç 240); Saldırı Emri (Hayvan Terbiyecisi, Sv 725, güç 260); Hızlı Yaylım (Nişancı, Sv 500, güç 290); Ezici Ok (Arbalet Ustası, Sv 725, güç 300); Nişanlı Atış (Keskin Nişancı, Sv 725, güç 320); Sürü Saldırısı (Hayvan Terbiyecisi, Sv 800, güç 360); Kuşatma Oku (Arbalet Ustası, Sv 800, güç 380); Ölümcül Hassasiyet (Keskin Nişancı, Sv 800, güç 400); Vahşetin Çağrısı (Hayvan Terbiyecisi, Sv 950, güç 520); Çelik Yağmuru (Arbalet Ustası, Sv 950, güç 560); Tek Kusursuz Atış (Keskin Nişancı, Sv 950, güç 650)
- **Saldırı güçlendirme** (5): İz Sürme (İzci, Sv 500, güç 15); Keskin Göz (Nişancı, Sv 400, güç 20); Kritik Odak (Keskin Nişancı, Sv 750, güç 25); Hızlı Doldurma (Arbalet Ustası, Sv 750, güç 25); Hayvani Kükreme (Hayvan Terbiyecisi, Sv 750, güç 25)
- **Zayıflatma** (3): İşaretli Hedef (Korucu, Sv 75, güç 20); Zırh Kırıcı (Nişancı, Sv 600, güç 30); Zırh Kırma Dizisi (Arbalet Ustası, Sv 900, güç 40)
- **Süreli hasar** (3): Patlayıcı Tuzak (Tuzakçı, Sv 725, güç 90); Zincirleme Tuzak (Tuzakçı, Sv 800, güç 120); Ölüm Sahası (Tuzakçı, Sv 950, güç 200)
- **Hasar azaltma (duruş/koruma)** (2): Avcı Hazırlığı (İzci, Sv 600, güç 30); Geri Çekilme (Korucu, Sv 10, güç 35)
- **Sersemletme** (2): Kapan Tuzağı (İzci, Sv 300, güç 130); Buz Tuzağı (Tuzakçı, Sv 750, güç 150)

**Aynı statı veren pasifler:**

- **kritik + isabet** (2): Avcı Gözü (Korucu, Sv 50, +%3 kritik, +%4 isabet); Kusursuz Nişan (Keskin Nişancı, Sv 700, +%5 kritik, +%3 isabet)
- **saldırı** (2): Uzak Menzil Ustalığı (Keskin Nişancı, Sv 850, +%6 saldırı); Hazırlanmış Zemin (Tuzakçı, Sv 700, +%5 saldırı)
- **saldırı + can** (2): Hayvan Bağı (Hayvan Terbiyecisi, Sv 700, +%6 saldırı, +%5 can); İlkel Bağ (Hayvan Terbiyecisi, Sv 900, +%6 saldırı, +%5 can)

## Büyücü

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Gizem Patlaması | Büyücü, 25 | Gizem Füzesi | Gizemci, 300 | Hasar vuruşu: güç 230, bekleme 3 / güç 240, bekleme 2 | aynı kol: üst üste biner |
| Gizem Patlaması | Büyücü, 25 | Element Oku | Element Büyücüsü, 300 | Hasar vuruşu: güç 230, bekleme 3 / güç 250, bekleme 3 | aynı kol: üst üste biner |
| Mana Akışı | Büyücü, 50 | Rün Bilgisi | Rün Ustası, 700 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Mana Akışı | Büyücü, 50 | Tutuşma | Ateş Büyücüsü, 700 | +%6 saldırı / +%5 saldırı | aynı kol: üst üste biner |
| Mana Akışı | Büyücü, 50 | Rün Ustalığı | Rün Ustası, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Mana Akışı | Büyücü, 50 | Yaşayan Alev | Ateş Büyücüsü, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Mana Akışı | Büyücü, 50 | Rün Aşırı Yüklemesi | Rün Ustası, 900 | +%6 saldırı / +%8 saldırı | aynı kol: üst üste biner |
| Element Yakınsaması | Element Büyücüsü, 600 | Ateş Topu | Ateş Büyücüsü, 725 | Hasar vuruşu: güç 340, bekleme 5 / güç 320, bekleme 3 | aynı kol: üst üste biner |
| Gizem Füzesi | Gizemci, 300 | Element Oku | Element Büyücüsü, 300 | Hasar vuruşu: güç 240, bekleme 2 / güç 250, bekleme 3 | farklı kol |
| Gizem Yankısı | Gizemci, 600 | Element Yakınsaması | Element Büyücüsü, 600 | Hasar vuruşu: güç 320, bekleme 4 / güç 340, bekleme 5 | farklı kol |
| Gizem Yankısı | Gizemci, 600 | Ateş Topu | Ateş Büyücüsü, 725 | Hasar vuruşu: güç 320, bekleme 4 / güç 320, bekleme 3 | farklı kol |
| Rün Bilgisi | Rün Ustası, 700 | Tutuşma | Ateş Büyücüsü, 700 | +%6 saldırı / +%5 saldırı | farklı kol |
| Rün Bilgisi | Rün Ustası, 700 | Yaşayan Alev | Ateş Büyücüsü, 850 | +%6 saldırı / +%6 saldırı | farklı kol |
| Rün Hapishanesi | Rün Ustası, 750 | Buz Novası | Buz Büyücüsü, 750 | Sersemletme: güç 150, 2 tur, bekleme 8 / güç 160, 2 tur, bekleme 8 | farklı kol |
| Rün Rezonansı | Rün Ustası, 775 | Çıra | Ateş Büyücüsü, 775 | +%10 kritik hasarı, +%4 kritik / +%12 kritik hasarı, +%3 kritik | farklı kol |
| Rün Patlatma | Rün Ustası, 800 | Yanma | Ateş Büyücüsü, 800 | Hasar vuruşu: güç 400, bekleme 5 / güç 420, bekleme 5 | farklı kol |
| Rün Ustalığı | Rün Ustası, 850 | Tutuşma | Ateş Büyücüsü, 700 | +%6 saldırı / +%5 saldırı | farklı kol |
| Rün Ustalığı | Rün Ustası, 850 | Yaşayan Alev | Ateş Büyücüsü, 850 | +%6 saldırı / +%6 saldırı | farklı kol |
| Rün Aşırı Yüklemesi | Rün Ustası, 900 | Tutuşma | Ateş Büyücüsü, 700 | +%8 saldırı / +%5 saldırı | farklı kol |
| Zaman Döngüsü | Zaman Büyücüsü, 900 | Derin Donma | Buz Büyücüsü, 700 | +%6 saldırı, +%3 hasar azaltma / +%4 saldırı, +%3 hasar azaltma | farklı kol |
| Rün Aşırı Yüklemesi | Rün Ustası, 900 | Yaşayan Alev | Ateş Büyücüsü, 850 | +%8 saldırı / +%6 saldırı | farklı kol |
| Zaman Döngüsü | Zaman Büyücüsü, 900 | Kalıcı Don | Buz Büyücüsü, 850 | +%6 saldırı, +%3 hasar azaltma / +%5 saldırı, +%3 hasar azaltma | farklı kol |
| Gizem Tekilliği | Rün Ustası, 950 | Meteor Fırtınası | Ateş Büyücüsü, 950 | Hasar vuruşu: güç 620, bekleme 16 / güç 640, bekleme 16 | farklı kol |
| Zamanı Durdur | Zaman Büyücüsü, 950 | Donmuş Dünya | Buz Büyücüsü, 950 | Sersemletme: güç 450, 3 tur, bekleme 16 / güç 420, 3 tur, bekleme 16 | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (12): Gizem Oku (Büyücü, Sv 1, güç 130); Gizem Patlaması (Büyücü, Sv 25, güç 230); Gizem Füzesi (Gizemci, Sv 300, güç 240); Element Oku (Element Büyücüsü, Sv 300, güç 250); Buz Mızrağı (Buz Büyücüsü, Sv 725, güç 280); Gizem Yankısı (Gizemci, Sv 600, güç 320); Ateş Topu (Ateş Büyücüsü, Sv 725, güç 320); Element Yakınsaması (Element Büyücüsü, Sv 600, güç 340); Rün Patlatma (Rün Ustası, Sv 800, güç 400); Yanma (Ateş Büyücüsü, Sv 800, güç 420); Gizem Tekilliği (Rün Ustası, Sv 950, güç 620); Meteor Fırtınası (Ateş Büyücüsü, Sv 950, güç 640)
- **Kalkan** (4): Mana Kalkanı (Büyücü, Sv 10, güç 25); Gizem Bariyeri (Gizemci, Sv 500, güç 35); Zaman Kalkanı (Zaman Büyücüsü, Sv 750, güç 40); Buz Bariyeri (Buz Büyücüsü, Sv 800, güç 45)
- **Sersemletme** (4): Rün Hapishanesi (Rün Ustası, Sv 750, güç 150); Buz Novası (Buz Büyücüsü, Sv 750, güç 160); Donmuş Dünya (Buz Büyücüsü, Sv 950, güç 420); Zamanı Durdur (Zaman Büyücüsü, Sv 950, güç 450)
- **Saldırı güçlendirme** (3): Mana Kanalı (Gizemci, Sv 400, güç 20); Büyü Dalgası (Büyücü, Sv 75, güç 25); Güç Rünü (Rün Ustası, Sv 725, güç 30)
- **Süreli hasar** (3): Element Dalgası (Element Büyücüsü, Sv 400, güç 80); Alev Dalgası (Ateş Büyücüsü, Sv 750, güç 100); Cehennem Ateşi (Ateş Büyücüsü, Sv 900, güç 160)
- **Zayıflatma** (2): Zaman Genişlemesi (Zaman Büyücüsü, Sv 725, güç 30); Mutlak Sıfır (Buz Büyücüsü, Sv 900, güç 40)

**Aynı statı veren pasifler:**

- **saldırı** (6): Mana Akışı (Büyücü, Sv 50, +%6 saldırı); Rün Bilgisi (Rün Ustası, Sv 700, +%6 saldırı); Rün Ustalığı (Rün Ustası, Sv 850, +%6 saldırı); Rün Aşırı Yüklemesi (Rün Ustası, Sv 900, +%8 saldırı); Tutuşma (Ateş Büyücüsü, Sv 700, +%5 saldırı); Yaşayan Alev (Ateş Büyücüsü, Sv 850, +%6 saldırı)
- **saldırı + hasar azaltma** (3): Zaman Döngüsü (Zaman Büyücüsü, Sv 900, +%6 saldırı, +%3 hasar azaltma); Derin Donma (Buz Büyücüsü, Sv 700, +%4 saldırı, +%3 hasar azaltma); Kalıcı Don (Buz Büyücüsü, Sv 850, +%5 saldırı, +%3 hasar azaltma)
- **saldırı + kritik** (2): Gizem Verimliliği (Gizemci, Sv 300, +%5 saldırı, +%2 kritik); Hızlandırılmış Büyü (Zaman Büyücüsü, Sv 775, +%5 saldırı, +%3 kritik)
- **kritik hasarı + kritik** (2): Rün Rezonansı (Rün Ustası, Sv 775, +%10 kritik hasarı, +%4 kritik); Çıra (Ateş Büyücüsü, Sv 775, +%12 kritik hasarı, +%3 kritik)

## Keşiş

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Chi Patlaması | Keşiş, 25 | Kombo Darbesi | Dövüşçü, 300 | Hasar vuruşu: güç 220, bekleme 3 / güç 230, bekleme 2 | aynı kol: üst üste biner |
| Chi Patlaması | Keşiş, 25 | Saptıran Avuç | Mürit, 500 | Hasar vuruşu: güç 220, bekleme 3 / güç 200, bekleme 3 | aynı kol: üst üste biner |
| Dövüş Ritmi | Dövüşçü, 300 | Kusursuz Kombinasyon | Kombo Ustası, 700 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Dövüş Ritmi | Dövüşçü, 300 | Sonsuz Dizi | Kombo Ustası, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Dövüş Ritmi | Dövüşçü, 300 | Bora Ritmi | Fırtına Yumruğu, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Dengeli Ruh | Mürit, 300 | Adamant Beden | Demir Beden, 900 | +%4 hasar azaltma, +%5 can / +%3 hasar azaltma, +%8 can | aynı kol: üst üste biner |
| Seri Yumruklar | Dövüşçü, 400 | Üçlü Darbe | Kombo Ustası, 725 | Hasar vuruşu: güç 280, bekleme 3 / güç 270, bekleme 2 | aynı kol: üst üste biner |
| Seri Yumruklar | Dövüşçü, 400 | Gök Gürültüsü Yumruğu | Fırtına Yumruğu, 725 | Hasar vuruşu: güç 280, bekleme 3 / güç 300, bekleme 3 | aynı kol: üst üste biner |
| Yükselen Tekme | Dövüşçü, 500 | Gök Gürültüsü Yumruğu | Fırtına Yumruğu, 725 | Hasar vuruşu: güç 330, bekleme 4 / güç 300, bekleme 3 | aynı kol: üst üste biner |
| Yükselen Tekme | Dövüşçü, 500 | Akıcı Tekme | Kombo Ustası, 750 | Hasar vuruşu: güç 330, bekleme 4 / güç 330, bekleme 3 | aynı kol: üst üste biner |
| Bitirici Avuç | Dövüşçü, 600 | Şimşek Yağmuru | Fırtına Yumruğu, 800 | Hasar vuruşu: güç 380, bekleme 4 / güç 420, bekleme 5 | aynı kol: üst üste biner |
| Kusursuz Kombinasyon | Kombo Ustası, 700 | Bora Ritmi | Fırtına Yumruğu, 850 | +%6 saldırı / +%6 saldırı | farklı kol |
| Üçlü Darbe | Kombo Ustası, 725 | Gök Gürültüsü Yumruğu | Fırtına Yumruğu, 725 | Hasar vuruşu: güç 270, bekleme 2 / güç 300, bekleme 3 | farklı kol |
| Demir Karşılık | Demir Beden, 750 | Seri Yumruklar | Dövüşçü, 400 | Hasar vuruşu: güç 260, bekleme 3 / güç 280, bekleme 3 | farklı kol |
| Demir Karşılık | Demir Beden, 750 | Üçlü Darbe | Kombo Ustası, 725 | Hasar vuruşu: güç 260, bekleme 3 / güç 270, bekleme 2 | farklı kol |
| Akıcı Tekme | Kombo Ustası, 750 | Gök Gürültüsü Yumruğu | Fırtına Yumruğu, 725 | Hasar vuruşu: güç 330, bekleme 3 / güç 300, bekleme 3 | farklı kol |
| Kusursuz Bitiriş | Kombo Ustası, 800 | Şimşek Yağmuru | Fırtına Yumruğu, 800 | Hasar vuruşu: güç 450, bekleme 5 / güç 420, bekleme 5 | farklı kol |
| Sonsuz Dizi | Kombo Ustası, 850 | Bora Ritmi | Fırtına Yumruğu, 850 | +%6 saldırı / +%6 saldırı | farklı kol |
| Kusursuz Teknik | Kombo Ustası, 900 | Fırtına Temposu | Fırtına Yumruğu, 700 | +%7 saldırı, +%3 kritik / +%5 saldırı, +%3 kritik | farklı kol |
| Yüz Yumruk Serisi | Kombo Ustası, 950 | Göksel Fırtına | Fırtına Yumruğu, 950 | Hasar vuruşu: güç 600, bekleme 15 / güç 620, bekleme 16 | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (15): Avuç Darbesi (Keşiş, Sv 1, güç 115); Saptıran Avuç (Mürit, Sv 500, güç 200); Chi Patlaması (Keşiş, Sv 25, güç 220); Kombo Darbesi (Dövüşçü, Sv 300, güç 230); Demir Karşılık (Demir Beden, Sv 750, güç 260); Üçlü Darbe (Kombo Ustası, Sv 725, güç 270); Seri Yumruklar (Dövüşçü, Sv 400, güç 280); Gök Gürültüsü Yumruğu (Fırtına Yumruğu, Sv 725, güç 300); Yükselen Tekme (Dövüşçü, Sv 500, güç 330); Akıcı Tekme (Kombo Ustası, Sv 750, güç 330); Bitirici Avuç (Dövüşçü, Sv 600, güç 380); Şimşek Yağmuru (Fırtına Yumruğu, Sv 800, güç 420); Kusursuz Bitiriş (Kombo Ustası, Sv 800, güç 450); Yüz Yumruk Serisi (Kombo Ustası, Sv 950, güç 600); Göksel Fırtına (Fırtına Yumruğu, Sv 950, güç 620)
- **Hasar azaltma (duruş/koruma)** (4): Ruh Adımı (Ruh Gezgini, Sv 725, güç 35); Chi Savunması (Mürit, Sv 300, güç 40); Taş Duruşu (Demir Beden, Sv 725, güç 55); Yürüyen Dağ (Demir Beden, Sv 950, güç 150)
- **Kalkan** (3): Odaklı Meditasyon (Mürit, Sv 600, güç 30); Astral Muhafız (Ruh Gezgini, Sv 800, güç 40); Güçlendirilmiş Chi (Demir Beden, Sv 800, güç 45)
- **Süreli iyileştirme** (2): Meditasyon (Keşiş, Sv 10, güç 6); Aşkın Akış (Ruh Gezgini, Sv 900, güç 8)
- **Sersemletme** (2): Uçan Tekme (Keşiş, Sv 75, güç 150); Fırtına Tekmesi (Fırtına Yumruğu, Sv 750, güç 180)
- **İyileştirme** (2): İç Şifa (Mürit, Sv 400, güç 25); Şifalı Avuç (Ruh Gezgini, Sv 750, güç 30)

**Aynı statı veren pasifler:**

- **saldırı** (4): Dövüş Ritmi (Dövüşçü, Sv 300, +%6 saldırı); Kusursuz Kombinasyon (Kombo Ustası, Sv 700, +%6 saldırı); Sonsuz Dizi (Kombo Ustası, Sv 850, +%6 saldırı); Bora Ritmi (Fırtına Yumruğu, Sv 850, +%6 saldırı)
- **hasar azaltma + can** (2): Dengeli Ruh (Mürit, Sv 300, +%4 hasar azaltma, +%5 can); Adamant Beden (Demir Beden, Sv 900, +%3 hasar azaltma, +%8 can)
- **hasar azaltma** (2): Demir Deri (Demir Beden, Sv 700, +%6 hasar azaltma); Elmas Dayanıklılık (Demir Beden, Sv 850, +%5 hasar azaltma)
- **saldırı + kritik** (2): Kusursuz Teknik (Kombo Ustası, Sv 900, +%7 saldırı, +%3 kritik); Fırtına Temposu (Fırtına Yumruğu, Sv 700, +%5 saldırı, +%3 kritik)

## Rahip

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Küçük Şifa | Rahip, 10 | Toplu Arındırma | Kâhin, 400 | İyileştirme: güç 20, bekleme 6 / güç 20, bekleme 8 | aynı kol: üst üste biner |
| Kutsal Lütuf | Aziz, 300 | Taşan Lütuf | Başrahip, 775 | +%5 can, +%2 can çalma / +%5 can, +%2 can çalma | aynı kol: üst üste biner |
| Pekişmiş İnanç | Kalkan Rahibi, 775 | Bozulmamış | Şeytan Kovucu, 850 | +%3 hasar azaltma, +%5 can / +%4 hasar azaltma, +%4 can | farklı kol |
| Sonsuz Merhamet | Başrahip, 850 | Pekişmiş İnanç | Kalkan Rahibi, 775 | +%4 hasar azaltma, +%4 can / +%3 hasar azaltma, +%5 can | farklı kol |
| Sonsuz Merhamet | Başrahip, 850 | Bozulmamış | Şeytan Kovucu, 850 | +%4 hasar azaltma, +%4 can / +%4 hasar azaltma, +%4 can | farklı kol |
| Sonsuz Merhamet | Başrahip, 850 | Kalkanın Yeniden Doğuşu | Kalkan Rahibi, 900 | +%4 hasar azaltma, +%4 can / +%3 hasar azaltma, +%6 can | farklı kol |
| Kalkanın Yeniden Doğuşu | Kalkan Rahibi, 900 | Bozulmamış | Şeytan Kovucu, 850 | +%3 hasar azaltma, +%6 can / +%4 hasar azaltma, +%4 can | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Kalkan** (6): Arındırma (Rahip, Sv 75, güç 25); Koruyucu Melek (Aziz, Sv 600, güç 40); Kutsal Kalkan (Kalkan Rahibi, Sv 725, güç 40); İlahi Müdahale (Başrahip, Sv 800, güç 50); Işık Bariyeri (Kalkan Rahibi, Sv 800, güç 55); İlahi Kale (Kalkan Rahibi, Sv 950, güç 100)
- **Hasar vuruşu** (5): Kutsal Darbe (Rahip, Sv 1, güç 125); İlahi Işık (Aziz, Sv 500, güç 260); Toplu Vahiy (Kehanetçi, Sv 800, güç 320); Şeytan Çıkarma (Şeytan Kovucu, Sv 800, güç 420); Işığın Hükmü (Şeytan Kovucu, Sv 950, güç 600)
- **İyileştirme** (5): Küçük Şifa (Rahip, Sv 10, güç 20); Toplu Arındırma (Kâhin, Sv 400, güç 20); Şifa Duası (Aziz, Sv 300, güç 30); Büyük Şifa (Başrahip, Sv 725, güç 40); Mucize (Başrahip, Sv 950, güç 90)
- **Saldırı güçlendirme** (4): Kutsama (Rahip, Sv 25, güç 20); Kehanet (Kâhin, Sv 300, güç 25); Kehanet Gücü (Kehanetçi, Sv 725, güç 30); Göklerin Sesi (Kehanetçi, Sv 950, güç 80)
- **Süreli iyileştirme** (2): Yenilenme (Aziz, Sv 400, güç 6); İyileşme Duası (Başrahip, Sv 750, güç 8)
- **Hasar azaltma (duruş/koruma)** (2): Bilgelik Kutsaması (Kehanetçi, Sv 750, güç 30); Kalkan Alanı (Kalkan Rahibi, Sv 750, güç 40)
- **Zayıflatma** (2): İlahi Kavrayış (Kâhin, Sv 600, güç 25); Kutsal Mühür (Şeytan Kovucu, Sv 725, güç 30)
- **Süreli hasar** (2): Kutsanmış Toprak (Kâhin, Sv 500, güç 80); Mutlak Arınma (Şeytan Kovucu, Sv 900, güç 150)

**Aynı statı veren pasifler:**

- **hasar azaltma + can** (4): Sonsuz Merhamet (Başrahip, Sv 850, +%4 hasar azaltma, +%4 can); Pekişmiş İnanç (Kalkan Rahibi, Sv 775, +%3 hasar azaltma, +%5 can); Kalkanın Yeniden Doğuşu (Kalkan Rahibi, Sv 900, +%3 hasar azaltma, +%6 can); Bozulmamış (Şeytan Kovucu, Sv 850, +%4 hasar azaltma, +%4 can)
- **can + can çalma** (3): Kutsal Lütuf (Aziz, Sv 300, +%5 can, +%2 can çalma); Taşan Lütuf (Başrahip, Sv 775, +%5 can, +%2 can çalma); İlahi Zarafet (Başrahip, Sv 900, +%6 can, +%3 can çalma)
- **saldırı + kritik** (3): Vahiy (Kehanetçi, Sv 700, +%4 saldırı, +%3 kritik); Gerçek Vahiy (Kehanetçi, Sv 900, +%6 saldırı, +%3 kritik); Arındıran Işık (Şeytan Kovucu, Sv 775, +%5 saldırı, +%2 kritik)
- **kritik + can** (2): Yüce Lütuf (Başrahip, Sv 700, +%3 kritik, +%5 can); Kalıcı Kutsama (Kehanetçi, Sv 850, +%2 kritik, +%5 can)

## Paladin

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| İlahi Savunma | Paladin, 10 | Muhafız Kutsaması | Işık Bekçisi, 750 | Hasar azaltma (duruş/koruma): güç 35, 3 tur, bekleme 7 / güç 35, 3 tur, bekleme 9 | aynı kol: üst üste biner |
| Kutsal İnanç | Paladin, 50 | Tapınakçı Disiplini | Tapınakçı, 850 | +%4 saldırı, +%3 hasar azaltma / +%3 saldırı, +%3 hasar azaltma | aynı kol: üst üste biner |
| Gayretli Adalet | Kutsal Savaşçı, 300 | Kutsal Saldırı | Tapınakçı, 700 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Gayretli Adalet | Kutsal Savaşçı, 300 | Haklı Öfke | Adalet Yargıcı, 775 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Gayret | Kutsal Savaşçı, 600 | Savaş Duası | Tapınakçı, 750 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| Kutsanmış Siper | Koruyucu, 600 | Kalkan Muhafızı | Kalkan Şövalyesi, 900 | +%5 blok, +%5 can / +%5 blok, +%6 can | aynı kol: üst üste biner |
| Kutsal Saldırı | Tapınakçı, 700 | Haklı Öfke | Adalet Yargıcı, 775 | +%6 saldırı / +%6 saldırı | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (6): Kutsal Vuruş (Paladin, Sv 1, güç 115); Kutsal Savaşçı Darbesi (Kutsal Savaşçı, Sv 300, güç 230); Tapınakçı Darbesi (Tapınakçı, Sv 725, güç 280); Kutsal Taarruz (Tapınakçı, Sv 800, güç 380); Işık İnfazı (Adalet Yargıcı, Sv 800, güç 450); Göklerin Gazabı (Adalet Yargıcı, Sv 950, güç 620)
- **Hasar azaltma (duruş/koruma)** (5): Kutsal Meydan Okuma (Koruyucu, Sv 500, güç 25); Muhafız Aurası (Koruyucu, Sv 400, güç 30); İlahi Savunma (Paladin, Sv 10, güç 35); Muhafız Kutsaması (Işık Bekçisi, Sv 750, güç 35); İlahi Blok (Kalkan Şövalyesi, Sv 725, güç 55)
- **İyileştirme** (4): Kutsal Burç (Kalkan Şövalyesi, Sv 800, güç 30); Şifa Işığı (Işık Bekçisi, Sv 725, güç 30); El Koyma (Paladin, Sv 75, güç 35); Kurtuluş Işığı (Işık Bekçisi, Sv 950, güç 90)
- **Kalkan** (3): Kutsal Kalkan (Koruyucu, Sv 300, güç 35); İnanç Kalkanı (Kalkan Şövalyesi, Sv 750, güç 45); Işık Kalesi (Kalkan Şövalyesi, Sv 950, güç 100)
- **Saldırı güçlendirme** (3): Gayret (Kutsal Savaşçı, Sv 600, güç 25); Savaş Duası (Tapınakçı, Sv 750, güç 25); Tapınakçı Yürüyüşü (Tapınakçı, Sv 950, güç 80)
- **Sersemletme** (2): Kutsal Hücum (Kutsal Savaşçı, Sv 400, güç 150); Adalet Çekici (Adalet Yargıcı, Sv 750, güç 200)
- **Zayıflatma** (2): Yargı (Kutsal Savaşçı, Sv 500, güç 25); Mahkûm Et (Adalet Yargıcı, Sv 725, güç 30)

**Aynı statı veren pasifler:**

- **saldırı** (3): Gayretli Adalet (Kutsal Savaşçı, Sv 300, +%6 saldırı); Kutsal Saldırı (Tapınakçı, Sv 700, +%6 saldırı); Haklı Öfke (Adalet Yargıcı, Sv 775, +%6 saldırı)
- **saldırı + hasar azaltma** (2): Kutsal İnanç (Paladin, Sv 50, +%4 saldırı, +%3 hasar azaltma); Tapınakçı Disiplini (Tapınakçı, Sv 850, +%3 saldırı, +%3 hasar azaltma)
- **blok + can** (2): Kutsanmış Siper (Koruyucu, Sv 600, +%5 blok, +%5 can); Kalkan Muhafızı (Kalkan Şövalyesi, Sv 900, +%5 blok, +%6 can)
- **hasar azaltma + can** (2): Şifalı Muhafaza (Işık Bekçisi, Sv 700, +%4 hasar azaltma, +%4 can); Işıldayan Muhafız (Işık Bekçisi, Sv 900, +%3 hasar azaltma, +%6 can)
- **saldırı + infaz** (2): Yargı İşareti (Adalet Yargıcı, Sv 700, +%3 saldırı, +%12 infaz); Son Yargı (Adalet Yargıcı, Sv 900, +%4 saldırı, +%15 infaz)

## Druid

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Doğanın Döngüsü | Druid, 50 | İlkel İçgüdü | Biçim Değiştiren, 300 | +%4 saldırı, +%4 can / +%4 saldırı, +%5 can | aynı kol: üst üste biner |
| Yabani Büyüme | Druid, 75 | Doğanın Armağanı | Koru Bekçisi, 400 | Saldırı güçlendirme: güç 20, 3 tur, bekleme 8 / güç 20, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| Çiçeklenme | Koru Bekçisi, 600 | Diken Kırbacı | Diken Bilgesi, 725 | Hasar vuruşu: güç 260, bekleme 3 / güç 280, bekleme 3 | aynı kol: üst üste biner |
| Çiçeklenme | Koru Bekçisi, 600 | Parçala | Ayı Muhafızı, 725 | Hasar vuruşu: güç 260, bekleme 3 / güç 260, bekleme 3 | farklı kol |
| Diken Kırbacı | Diken Bilgesi, 725 | Parçala | Ayı Muhafızı, 725 | Hasar vuruşu: güç 280, bekleme 3 / güç 260, bekleme 3 | farklı kol |
| Çalılık Alanı | Diken Bilgesi, 750 | Ay Ateşi | Ay Çağıran, 725 | Süreli hasar: güç 100, 4 tur, bekleme 7 / güç 90, 4 tur, bekleme 6 | farklı kol |
| Doğal Ahenk | Can Bağlayıcı, 775 | İlkel Dayanıklılık | Ayı Muhafızı, 775 | +%3 hasar azaltma, +%5 can / +%3 hasar azaltma, +%6 can | farklı kol |
| Zehirli Dikenler | Diken Bilgesi, 775 | Göksel Hizalanma | Ay Çağıran, 850 | +%5 saldırı / +%6 saldırı | farklı kol |
| Doğanın Misillemesi | Diken Bilgesi, 800 | Yıldız Yağmuru | Ay Çağıran, 800 | Hasar vuruşu: güç 400, bekleme 5 / güç 420, bekleme 5 | farklı kol |
| Dünya Kökü Bağı | Can Bağlayıcı, 900 | İlkel Dayanıklılık | Ayı Muhafızı, 775 | +%3 hasar azaltma, +%6 can / +%3 hasar azaltma, +%6 can | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (7): Doğa Oku (Druid, Sv 1, güç 125); Çiçeklenme (Koru Bekçisi, Sv 600, güç 260); Parçala (Ayı Muhafızı, Sv 725, güç 260); Diken Kırbacı (Diken Bilgesi, Sv 725, güç 280); Doğanın Misillemesi (Diken Bilgesi, Sv 800, güç 400); Yıldız Yağmuru (Ay Çağıran, Sv 800, güç 420); Göksel Tutulma (Ay Çağıran, Sv 950, güç 620)
- **Saldırı güçlendirme** (4): Yabani Büyüme (Druid, Sv 75, güç 20); Doğanın Armağanı (Koru Bekçisi, Sv 400, güç 20); Avcı Formu (Biçim Değiştiren, Sv 500, güç 25); Ay Kutsaması (Ay Çağıran, Sv 750, güç 25)
- **Süreli iyileştirme** (3): Gençleşme (Druid, Sv 10, güç 6); Hayat Çiçeği (Can Bağlayıcı, Sv 725, güç 8); Dünya Ağacı (Can Bağlayıcı, Sv 950, güç 15)
- **İyileştirme** (3): Yeniden Büyüme (Koru Bekçisi, Sv 300, güç 30); Şifalı Sarmaşıklar (Can Bağlayıcı, Sv 750, güç 35); Toplu Yeniden Büyüme (Can Bağlayıcı, Sv 800, güç 45)
- **Süreli hasar** (3): Ay Ateşi (Ay Çağıran, Sv 725, güç 90); Çalılık Alanı (Diken Bilgesi, Sv 750, güç 100); Vahşinin Gazabı (Diken Bilgesi, Sv 950, güç 200)
- **Hasar azaltma (duruş/koruma)** (3): Ayı Formu (Biçim Değiştiren, Sv 300, güç 45); Demir Post (Ayı Muhafızı, Sv 750, güç 55); Büyük Ayının Avatarı (Ayı Muhafızı, Sv 950, güç 150)
- **Sersemletme** (2): Dolanan Kökler (Druid, Sv 25, güç 140); Vahşi Hücum (Biçim Değiştiren, Sv 400, güç 160)
- **Kalkan** (2): Diken Kalkanı (Koru Bekçisi, Sv 500, güç 35); Vahşi Muhafaza (Ayı Muhafızı, Sv 800, güç 45)

**Aynı statı veren pasifler:**

- **can + can çalma** (3): Hayatın Bekçisi (Koru Bekçisi, Sv 300, +%5 can, +%2 can çalma); Yaşayan Ağ (Can Bağlayıcı, Sv 700, +%5 can, +%3 can çalma); Yaprak Dökmeyen (Can Bağlayıcı, Sv 850, +%4 can, +%3 can çalma)
- **hasar azaltma + can** (3): Doğal Ahenk (Can Bağlayıcı, Sv 775, +%3 hasar azaltma, +%5 can); Dünya Kökü Bağı (Can Bağlayıcı, Sv 900, +%3 hasar azaltma, +%6 can); İlkel Dayanıklılık (Ayı Muhafızı, Sv 775, +%3 hasar azaltma, +%6 can)
- **saldırı + can** (2): Doğanın Döngüsü (Druid, Sv 50, +%4 saldırı, +%4 can); İlkel İçgüdü (Biçim Değiştiren, Sv 300, +%4 saldırı, +%5 can)
- **saldırı** (2): Zehirli Dikenler (Diken Bilgesi, Sv 775, +%5 saldırı); Göksel Hizalanma (Ay Çağıran, Sv 850, +%6 saldırı)
- **savunma + can** (2): Muhafız Formu (Ayı Muhafızı, Sv 700, +%8 savunma, +%5 can); Yaşlı Ayı (Ayı Muhafızı, Sv 900, +%4 savunma, +%8 can)
- **saldırı + kritik** (2): Ay Döngüsü (Ay Çağıran, Sv 700, +%5 saldırı, +%3 kritik); Dolunay (Ay Çağıran, Sv 900, +%6 saldırı, +%4 kritik)

## Ozan

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Ritim | Ozan, 50 | Muharebe Ritmi | Savaş Ezgicisi, 700 | +%4 saldırı, +%3 kritik / +%3 saldırı, +%4 kritik | aynı kol: üst üste biner |
| Ritim | Ozan, 50 | Savaş Marşı | Savaş Ezgicisi, 900 | +%4 saldırı, +%3 kritik / +%6 saldırı, +%3 kritik | aynı kol: üst üste biner |
| İlham Veren Koro | Ozan, 75 | Hız Şarkısı | Âşık, 500 | Hasar azaltma (duruş/koruma): güç 30, 3 tur, bekleme 8 / güç 30, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| İlham Veren Koro | Ozan, 75 | Şarkı Örgüsü | Virtüöz, 750 | Hasar azaltma (duruş/koruma): güç 30, 3 tur, bekleme 8 / güç 30, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| Savaş Dizesi | Destancı, 300 | Savaşçının Ritmi | Savaş Ezgicisi, 850 | +%6 saldırı / +%6 saldırı | aynı kol: üst üste biner |
| Melodik İlham | Âşık, 300 | Büyük Ahenk | Virtüöz, 900 | +%4 saldırı, +%4 can / +%6 saldırı, +%4 can | aynı kol: üst üste biner |
| Hız Şarkısı | Âşık, 500 | Şarkı Örgüsü | Virtüöz, 750 | Hasar azaltma (duruş/koruma): güç 30, 3 tur, bekleme 8 / güç 30, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| Cesaret Şarkısı | Âşık, 300 | Savaş Şarkısı | Destancı, 300 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | farklı kol |
| Harmonik Rezonans | Âşık, 600 | Savaş Davulu | Destancı, 500 | Hasar vuruşu: güç 260, bekleme 3 / güç 250, bekleme 3 | farklı kol |
| Karmaşık Ahenk | Virtüöz, 700 | Ezici Ahenksizlik | Ağıtçı, 775 | +%5 saldırı, +%2 kritik / +%5 saldırı, +%2 kritik | farklı kol |
| Büyük Melodi | Virtüöz, 725 | Savaş Ezgisi | Savaş Ezgicisi, 725 | Saldırı güçlendirme: güç 30, 3 tur, bekleme 8 / güç 30, 3 tur, bekleme 8 | farklı kol |
| Sonsuz Şarkı | Virtüöz, 850 | Savaş Dizesi | Destancı, 300 | +%5 saldırı / +%6 saldırı | farklı kol |
| Sonsuz Şarkı | Virtüöz, 850 | Savaşçının Ritmi | Savaş Ezgicisi, 850 | +%5 saldırı / +%6 saldırı | farklı kol |
| Büyük Senfoni | Virtüöz, 950 | Fetih Şarkısı | Savaş Ezgicisi, 950 | Saldırı güçlendirme: güç 90, 3 tur, bekleme 20 / güç 90, 3 tur, bekleme 18 | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Saldırı güçlendirme** (7): Hızlı Tempo (Ozan, Sv 25, güç 20); Cesaret Şarkısı (Âşık, Sv 300, güç 25); Savaş Şarkısı (Destancı, Sv 300, güç 25); Büyük Melodi (Virtüöz, Sv 725, güç 30); Savaş Ezgisi (Savaş Ezgicisi, Sv 725, güç 30); Büyük Senfoni (Virtüöz, Sv 950, güç 90); Fetih Şarkısı (Savaş Ezgicisi, Sv 950, güç 90)
- **Hasar vuruşu** (6): Keskin Dize (Ozan, Sv 1, güç 115); Savaş Davulu (Destancı, Sv 500, güç 250); Harmonik Rezonans (Âşık, Sv 600, güç 260); Öfke Davulları (Savaş Ezgicisi, Sv 750, güç 300); Usta Performans (Virtüöz, Sv 800, güç 360); Son Dize (Ağıtçı, Sv 800, güç 420)
- **Hasar azaltma (duruş/koruma)** (4): İlham Veren Koro (Ozan, Sv 75, güç 30); Hız Şarkısı (Âşık, Sv 500, güç 30); Şarkı Örgüsü (Virtüöz, Sv 750, güç 30); Kahramanlık Korosu (Savaş Ezgicisi, Sv 800, güç 35)
- **İyileştirme** (3): Şifa Melodisi (Âşık, Sv 400, güç 25); Yenilenme Şarkısı (Ahenkçi, Sv 800, güç 40); Kurtuluş Şarkısı (Ahenkçi, Sv 950, güç 90)
- **Süreli iyileştirme** (2): Dinlenme Şarkısı (Ozan, Sv 10, güç 5); Onarıcı Koro (Ahenkçi, Sv 725, güç 8)
- **Zayıflatma** (2): Alaycı Dize (Destancı, Sv 400, güç 25); Kırılganlık Ağıtı (Ağıtçı, Sv 725, güç 30)
- **Süreli hasar** (2): Umutsuzluk Şarkısı (Ağıtçı, Sv 750, güç 90); Son Ağıt (Ağıtçı, Sv 950, güç 200)

**Aynı statı veren pasifler:**

- **saldırı + kritik** (5): Ritim (Ozan, Sv 50, +%4 saldırı, +%3 kritik); Karmaşık Ahenk (Virtüöz, Sv 700, +%5 saldırı, +%2 kritik); Muharebe Ritmi (Savaş Ezgicisi, Sv 700, +%3 saldırı, +%4 kritik); Savaş Marşı (Savaş Ezgicisi, Sv 900, +%6 saldırı, +%3 kritik); Ezici Ahenksizlik (Ağıtçı, Sv 775, +%5 saldırı, +%2 kritik)
- **saldırı** (3): Sonsuz Şarkı (Virtüöz, Sv 850, +%5 saldırı); Savaş Dizesi (Destancı, Sv 300, +%6 saldırı); Savaşçının Ritmi (Savaş Ezgicisi, Sv 850, +%6 saldırı)
- **saldırı + can** (2): Melodik İlham (Âşık, Sv 300, +%4 saldırı, +%4 can); Büyük Ahenk (Virtüöz, Sv 900, +%6 saldırı, +%4 can)
- **can + can çalma** (2): Şifa Ahengi (Ahenkçi, Sv 700, +%5 can, +%2 can çalma); Sonsuz Ahenk (Ahenkçi, Sv 850, +%4 can, +%3 can çalma)
- **hasar azaltma + can** (2): Yankılanan Şifa (Ahenkçi, Sv 775, +%3 hasar azaltma, +%5 can); Kusursuz Rezonans (Ahenkçi, Sv 900, +%3 hasar azaltma, +%6 can)
- **saldırı + infaz** (2): Zayıflık Şarkısı (Ağıtçı, Sv 700, +%3 saldırı, +%12 infaz); Ölümün Melodisi (Ağıtçı, Sv 900, +%6 saldırı, +%10 infaz)

## Şaman

**Neredeyse aynı yetenekler** (aynı etki, yakın değer, farklı düğüm):

| Yetenek A | Dal / Sv | Yetenek B | Dal / Sv | Değerler | Durum |
|---|---|---|---|---|---|
| Ata Sesi | Ruh Çağıran, 500 | Ataların Kutsaması | Ata Bilgesi, 725 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | aynı kol: üst üste biner |
| Koruma Totemi | Totemci, 500 | Ruh Muhafazası | Ruh Çağıran, 400 | Kalkan: güç 35, 3 tur, bekleme 10 / güç 35, 3 tur, bekleme 10 | farklı kol |
| Totem Dalgası | Totemci, 600 | Ata Sesi | Ruh Çağıran, 500 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | farklı kol |
| Totem Dalgası | Totemci, 600 | Ataların Kutsaması | Ata Bilgesi, 725 | Saldırı güçlendirme: güç 25, 3 tur, bekleme 8 / güç 25, 3 tur, bekleme 8 | farklı kol |
| Şimşek Totemi | Fırtına Totemcisi, 725 | Ruh Laneti | Lanetçi, 750 | Süreli hasar: güç 100, 4 tur, bekleme 6 / güç 100, 4 tur, bekleme 7 | farklı kol |
| Statik Rezonans | Fırtına Totemcisi, 775 | Lanet Ustalığı | Lanetçi, 850 | +%12 kritik hasarı, +%2 kritik / +%12 kritik hasarı, +%2 kritik | farklı kol |
| Dağ Ruhu | Toprak Bekçisi, 800 | Ata Koruması | Ata Bilgesi, 800 | Kalkan: güç 50, 3 tur, bekleme 11 / güç 50, 3 tur, bekleme 12 | farklı kol |
| Zincir Şimşek Totemi | Fırtına Totemcisi, 800 | Kıyamet Laneti | Lanetçi, 800 | Hasar vuruşu: güç 420, bekleme 5 / güç 420, bekleme 5 | farklı kol |
| Fırtına Kanalı Ustalığı | Fırtına Totemcisi, 850 | Birikimli Lanet | Lanetçi, 700 | +%6 saldırı / +%5 saldırı | farklı kol |
| Toprağın Koruması | Toprak Bekçisi, 900 | Kadimlerin Rehberliği | Ata Bilgesi, 775 | +%3 hasar azaltma, +%6 can / +%3 hasar azaltma, +%4 can | farklı kol |

**Aynı işi yapan aktif yetenekler:**

- **Hasar vuruşu** (7): Ruh Oku (Şaman, Sv 1, güç 125); Şimşek Çarpması (Şaman, Sv 25, güç 220); Ruh Bağı (Ruh Çağıran, Sv 600, güç 260); Gök Gürültüsü Alanı (Fırtına Totemcisi, Sv 750, güç 320); Zincir Şimşek Totemi (Fırtına Totemcisi, Sv 800, güç 420); Kıyamet Laneti (Lanetçi, Sv 800, güç 420); Ruh Borası (Fırtına Totemcisi, Sv 950, güç 620)
- **Kalkan** (5): Koruma Totemi (Totemci, Sv 500, güç 35); Ruh Muhafazası (Ruh Çağıran, Sv 400, güç 35); Dağ Ruhu (Toprak Bekçisi, Sv 800, güç 50); Ata Koruması (Ata Bilgesi, Sv 800, güç 50); Dünya Sütunu (Toprak Bekçisi, Sv 950, güç 100)
- **Saldırı güçlendirme** (4): Ata Rehberliği (Şaman, Sv 75, güç 20); Totem Dalgası (Totemci, Sv 600, güç 25); Ata Sesi (Ruh Çağıran, Sv 500, güç 25); Ataların Kutsaması (Ata Bilgesi, Sv 725, güç 25)
- **Süreli hasar** (4): Fırtına Totemi (Totemci, Sv 400, güç 80); Şimşek Totemi (Fırtına Totemcisi, Sv 725, güç 100); Ruh Laneti (Lanetçi, Sv 750, güç 100); Ruh Fırtınası (Lanetçi, Sv 950, güç 200)
- **İyileştirme** (3): Ata Şifası (Ruh Çağıran, Sv 300, güç 30); Ruh Onarımı (Ata Bilgesi, Sv 750, güç 40); Atalar Meclisi (Ata Bilgesi, Sv 950, güç 90)
- **Süreli iyileştirme** (2): Şifa Totemi (Şaman, Sv 10, güç 6); Kadim Birliktelik (Ata Bilgesi, Sv 900, güç 12)
- **Hasar azaltma (duruş/koruma)** (2): Toprak Totemi (Totemci, Sv 300, güç 40); Taş Deri Totemi (Toprak Bekçisi, Sv 725, güç 50)

**Aynı statı veren pasifler:**

- **can + can çalma** (3): Ata Bağı (Ruh Çağıran, Sv 300, +%5 can, +%2 can çalma); Ata Bilgeliği (Ata Bilgesi, Sv 700, +%5 can, +%3 can çalma); Dayanıklı Ruh (Ata Bilgesi, Sv 850, +%4 can, +%3 can çalma)
- **hasar azaltma + can** (2): Toprağın Koruması (Toprak Bekçisi, Sv 900, +%3 hasar azaltma, +%6 can); Kadimlerin Rehberliği (Ata Bilgesi, Sv 775, +%3 hasar azaltma, +%4 can)
- **saldırı + kritik** (2): Fırtına Kanalı (Fırtına Totemcisi, Sv 700, +%5 saldırı, +%3 kritik); Yaşayan Fırtına (Fırtına Totemcisi, Sv 900, +%6 saldırı, +%4 kritik)
- **kritik hasarı + kritik** (2): Statik Rezonans (Fırtına Totemcisi, Sv 775, +%12 kritik hasarı, +%2 kritik); Lanet Ustalığı (Lanetçi, Sv 850, +%12 kritik hasarı, +%2 kritik)
- **saldırı** (2): Fırtına Kanalı Ustalığı (Fırtına Totemcisi, Sv 850, +%6 saldırı); Birikimli Lanet (Lanetçi, Sv 700, +%5 saldırı)
- **saldırı + infaz** (2): Süregelen Lanet (Lanetçi, Sv 775, +%5 saldırı, +%8 infaz); Kusursuz Lanet (Lanetçi, Sv 900, +%6 saldırı, +%10 infaz)
