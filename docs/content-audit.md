# İçerik Denetimi — eşya, üretim, ganimet, kademe, set, canavar

> Otomatik üretilir: `python scripts/content_audit.py`. Elle düzenleme — betiği değiştir. Kurallar `server/internal/item` (GearPool, CombatHides, TierRarity) ve `activity.DropLevel` aynasıdır. Yol haritası: [content-roadmap.md](content-roadmap.md).

## 0. Özet

| Alan | Adet | Aralık |
|---|---:|---|
| Ekipman şablonu | 1808 | Sv 1–985, 47 kademe |
| Toplama kaynağı | 150 | meslek Sv 1–100, kademe T0–T7 |
| Diğer eşya (külçe, iksir, yemek, eski malzeme) | 66 | |
| Tarif | 1867 | blacksmithing 788, woodworking 376, tailoring 282, leatherworking 235, jewelcrafting 126, cooking 30, smelting 24, alchemy 6 |
| Canavar | 140 | Sv 1–1000 (+6 baskın boss’u) |
| Zindan / baskın | 25 / 6 | Sv 6–1000 |
| Görev | 18 | Sv 1–10 |
| Set | 0 | tasarım var (`docs/item-engine/set.md`), uygulanmadı |
| Dükkân | 1 | 38 eşya |

## 1. Ekipman kademeleri (tier)

Her kademe: 10 sınıf silahı (sınıfa kilitli) + 5 zırh parçası (sınıfsız). Kaynak sütunu: D = dükkân, Ü = üretim, G = ganimet.

| Sv | Kademe (göğüs adı) | Silah | Zırh | Saldırı (kılıç) | Savunma (göğüs) | Kaynak | Bu kademeyi düşüren canavar Sv |
|---:|---|---:|---:|---:|---:|---|---|
| 1 | Deri Göğüslük | 20 | 15 | 6 | 2 | DGÜ | 1–3 |
| 5 | Demir Göğüslük | 20 | 15 | 16 | 4 | GÜ | 5–8 |
| 10 | Pul Göğüslük | 20 | 15 | 29 | 6 | GÜ | 10–14 |
| 15 | Çelik Göğüslük | 20 | 15 | 42 | 7 | GÜ | 15–18 |
| 20 | Şövalye Göğüslüğü | 20 | 15 | 54 | 9 | GÜ | 20–24 |
| 25 | Rünlü Göğüslük | 20 | 15 | 66 | 11 | GÜ | 25–30 |
| 35 | Mithril Göğüslük | 20 | 15 | 88 | 14 | GÜ | 35–45 |
| 50 | Adamant Göğüslük | 20 | 15 | 120 | 19 | GÜ | 50–60 |
| 65 | Ejder Dövme Göğüslük | 20 | 15 | 152 | 24 | GÜ | 65–75 |
| 80 | Yıldız Dövme Göğüslük | 20 | 15 | 184 | 29 | GÜ | 80–92 |
| 95 | Titan Göğüslük | 20 | 15 | 216 | 34 | GÜ | 95–100 |
| 110 | Kraliyet Göğüslük | 20 | 16 | 249 | 39 | GÜ | 110–132 |
| 135 | Güneş Dövme Göğüslük | 20 | 16 | 304 | 48 | GÜ | 140–150 |
| 160 | Obsidyen Göğüslük | 20 | 16 | 359 | 57 | GÜ | 160–182 |
| 185 | Kum Çeliği Göğüslük | 20 | 16 | 414 | 66 | GÜ | 190–200 |
| 210 | Rün Bağlı Göğüslük | 20 | 16 | 469 | 75 | GÜ | 210–232 |
| 235 | Anka Göğüslük | 20 | 16 | 524 | 84 | GÜ | 240–250 |
| 260 | Dehşet Dövme Göğüslük | 20 | 16 | 579 | 93 | GÜ | 260–282 |
| 285 | Hükümdar Göğüslük | 20 | 16 | 634 | 102 | GÜ | 290–300 |
| 310 | Mercan Göğüslük | 20 | 16 | 689 | 111 | GÜ | 310–332 |
| 335 | Ejderha Kemiği Göğüslük | 20 | 16 | 744 | 120 | GÜ | 340–350 |
| 360 | Sis Örgü Göğüslük | 20 | 16 | 799 | 129 | GÜ | 360–382 |
| 385 | Prizmatik Göğüslük | 20 | 16 | 854 | 138 | GÜ | 390–400 |
| 410 | Magma Dövme Göğüslük | 20 | 16 | 909 | 147 | GÜ | 410–432 |
| 435 | Bora Göğüslük | 20 | 16 | 964 | 156 | GÜ | 440–450 |
| 460 | Hiçlik Göğüslük | 20 | 16 | 1019 | 165 | GÜ | 460–482 |
| 485 | Melek Göğüslük | 20 | 16 | 1074 | 174 | GÜ | 490–500 |
| 510 | Tutulma Göğüslük | 20 | 16 | 1129 | 183 | GÜ | 510–532 |
| 535 | Ejder Lordu Göğüslük | 20 | 16 | 1184 | 191 | GÜ | 540–550 |
| 560 | Göksel Göğüslük | 20 | 16 | 1239 | 200 | GÜ | 560–582 |
| 585 | Ebedi Göğüslük | 20 | 16 | 1294 | 209 | GÜ | 590–600 |
| 610 | Yıldızsız Göğüslük | 20 | 16 | 1349 | 218 | GÜ | 610–632 |
| 635 | Gece Çeliği Göğüslük | 20 | 16 | 1404 | 227 | GÜ | 640–650 |
| 660 | Zaman Dövme Göğüslük | 20 | 16 | 1459 | 236 | GÜ | 660–682 |
| 685 | Yükselmiş Göğüslük | 20 | 16 | 1514 | 245 | GÜ | 690–700 |
| 710 | Fırtına Doğumlu Göğüslük | 20 | 16 | 1569 | 254 | GÜ | 710–732 |
| 735 | Yarık Dövme Göğüslük | 20 | 16 | 1624 | 263 | GÜ | 740–750 |
| 760 | Kaos Göğüslük | 20 | 16 | 1679 | 272 | GÜ | 760–782 |
| 785 | Gök Dövme Göğüslük | 20 | 16 | 1734 | 281 | GÜ | 790–800 |
| 810 | Başmelek Göğüslük | 20 | 16 | 1789 | 290 | GÜ | 810–832 |
| 835 | Şafak Göğüslük | 20 | 16 | 1844 | 299 | GÜ | 840–850 |
| 860 | Kül Göğüslük | 20 | 16 | 1899 | 308 | GÜ | 860–882 |
| 885 | Leviathan Göğüslük | 20 | 16 | 1954 | 317 | GÜ | 890–900 |
| 910 | Kıyamet Dövme Göğüslük | 20 | 16 | 2009 | 326 | GÜ | 910–932 |
| 935 | Hiçlik Dövme Göğüslük | 20 | 16 | 2064 | 335 | GÜ | 940–950 |
| 960 | Tanrı Katili Göğüslük | 20 | 16 | 2119 | 344 | GÜ | 960–982 |
| 985 | Kadim Çağ Göğüslük | 20 | 16 | 2174 | 353 | GÜ | 990–1000 |

Gereksinimli (rastgele düşmeyen) eşyalar: Titan Siper Kalkanı (Sv 600), Kraliyet Kalkan (Sv 110), Güneş Dövme Kalkan (Sv 135), Obsidyen Kalkan (Sv 160), Kum Çeliği Kalkan (Sv 185), Rün Bağlı Kalkan (Sv 210), Anka Kalkan (Sv 235), Dehşet Dövme Kalkan (Sv 260), Hükümdar Kalkan (Sv 285), Mercan Kalkan (Sv 310), Ejderha Kemiği Kalkan (Sv 335), Sis Örgü Kalkan (Sv 360), Prizmatik Kalkan (Sv 385), Magma Dövme Kalkan (Sv 410), Bora Kalkan (Sv 435), Hiçlik Kalkan (Sv 460), Melek Kalkan (Sv 485), Tutulma Kalkan (Sv 510), Ejder Lordu Kalkan (Sv 535), Göksel Kalkan (Sv 560), Ebedi Kalkan (Sv 585), Yıldızsız Kalkan (Sv 610), Gece Çeliği Kalkan (Sv 635), Zaman Dövme Kalkan (Sv 660), Yükselmiş Kalkan (Sv 685), Fırtına Doğumlu Kalkan (Sv 710), Yarık Dövme Kalkan (Sv 735), Kaos Kalkan (Sv 760), Gök Dövme Kalkan (Sv 785), Başmelek Kalkan (Sv 810), Şafak Kalkan (Sv 835), Kül Kalkan (Sv 860), Leviathan Kalkan (Sv 885), Kıyamet Dövme Kalkan (Sv 910), Hiçlik Dövme Kalkan (Sv 935), Tanrı Katili Kalkan (Sv 960), Kadim Çağ Kalkan (Sv 985).

Kademe aralığı: Sv 1–25 arası 5, 25–95 arası 10–15, 95 sonrası 25 seviye. En geniş boşluk: Sv 110→135, Sv 135→160, Sv 160→185.

### 1.1 Slot doluluğu

| Slot | Şablon | Not |
|---|---:|---|
| mainhand | 940 |  |
| offhand | 37 | yalnız 37 şablon |
| helm | 141 |  |
| chest | 141 |  |
| legs | 141 |  |
| boots | 141 |  |
| gloves | 141 |  |
| necklace | 42 | yalnız 42 şablon |
| ring1 | 42 | yalnız 42 şablon |
| ring2 | 42 | yalnız 42 şablon |
| trinket1 | 42 | yalnız 42 şablon |
| trinket2 | 42 | yalnız 42 şablon |

### 1.2 Silah aileleri ve zırh türleri

| Aile | Şablon |
|---|---:|
| sword | 94 |
| axe | 47 |
| mace | 47 |
| shield | 37 |
| dagger | 47 |
| short_sword | 47 |
| fist_weapon | 47 |
| tonfa | 47 |
| bow | 47 |
| crossbow | 47 |
| staff | 47 |
| wand | 47 |
| holy_tome | 47 |
| two_handed_mace | 47 |
| sickle | 47 |
| nature_focus | 47 |
| lute | 47 |
| rapier | 47 |
| ritual_axe | 47 |
| totem_focus | 47 |
| cloth_armor | 235 |
| leather_armor | 235 |
| plate_armor | 235 |

### 1.3 Nadirlik ve ek (affix)

| Nadirlik | Değer çarpanı | Stat çarpanı (DD-275) | Ek kuralı | Ekipman ganimetinde |
|---|---:|---:|---|---|
| sıradan | ×1 | ×1 | yok | %60 |
| seçkin | ×1.5 | ×1.05 | yok | %25 · boss %50 |
| nadir | ×3 | ×1.10 | %50 tek ek | %11 · boss %35 |
| destansı | ×8 | ×1.20 | ön ek + %50 son ek | %4 · boss %13 |
| efsanevi | ×25 | ×1.30 | ön + son ek | boss %2 |
| mitik | ×100 | ×1.50 | ön + son ek (+ benzersiz etki ayrılmış) | **hiç düşmez** (yalnız T7 malzeme) |

Her ek bir stat verir: taban + seviye başı × eşya Sv (`item/affix.go` `affixStats`).

## 2. Eşya listesi

### 2.1 Toplama kaynakları

| Kod | Ad | Tür | Meslek Sv | Kademe | Nadirlik | Değer | Kullanan tarif |
|---|---|---|---:|---:|---|---:|---|
| tiny_guppy | Minik Lepistes | fishing/fish | 1 | T0 | sıradan | 2c | Izgara Lepistes |
| river_sardine | Nehir Sardalyası | fishing/fish | 4 | T0 | sıradan | 3c | Sardalya Şiş |
| minnow | Golyan | fishing/fish | 7 | T0 | sıradan | 3c | Naneli Balık Çorbası |
| small_pike | Küçük Turna | fishing/fish | 10 | T1 | sıradan | 4c | Turna Fileto |
| common_carp | Sazan | fishing/fish | 12 | T1 | sıradan | 4c | Lavantalı Sazan |
| silver_herring | Gümüş Ringa | fishing/fish | 14 | T1 | sıradan | 5c | Izgara Gümüş Ringa |
| river_trout | Dere Alabalığı | fishing/fish | 16 | T1 | sıradan | 5c | Otlu Alabalık |
| salmon | Somon | fishing/fish | 18 | T1 | sıradan | 6c | Füme Somon |
| bass | Tatlı Su Levreği | fishing/fish | 20 | T2 | nadir-altı | 6c | Levrek Güveci |
| lake_perch | Göl Levreği | fishing/fish | 22 | T2 | nadir-altı | 7c | Izgara Göl Levreği |
| red_snapper | Kırmızı Mercan | fishing/fish | 24 | T2 | nadir-altı | 8c | Izgara Kırmızı Mercan |
| river_catfish | Nehir Yayını | fishing/fish | 26 | T2 | nadir-altı | 9c | Izgara Nehir Yayını |
| bluefin_tuna | Mavi Yüzgeçli Ton | fishing/fish | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| rainbow_trout | Gökkuşağı Alabalığı | fishing/fish | 30 | T3 | nadir-altı | 11c | Izgara Gökkuşağı Alabalığı |
| royal_carp | Kraliyet Sazanı | fishing/fish | 33 | T3 | nadir-altı | 14c | Kraliyet Sazanı Ziyafeti |
| emerald_pike | Zümrüt Turnası | fishing/fish | 36 | T3 | nadir-altı | 16c | Izgara Zümrüt Turnası |
| silver_barracuda | Gümüş Barakuda | fishing/fish | 39 | T3 | nadir-altı | 19c | Izgara Gümüş Barakuda |
| golden_sturgeon | Altın Mersin | fishing/fish | 42 | T3 | nadir-altı | 23c | Mantarlı Mersin Turtası |
| crimson_swordfish | Kızıl Kılıçbalığı | fishing/fish | 45 | T4 | nadir | 28c | Fırtına Izgarası Kılıçbalığı |
| stormscale_eel | Fırtına Pullu Yılanbalığı | fishing/fish | 50 | T4 | nadir | 37c | Izgara Fırtına Pullu Yılanbalığı |
| frostfin_salmon | Ayaz Yüzgeçli Somon | fishing/fish | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası |
| obsidian_tuna | Obsidyen Ton | fishing/fish | 60 | T5 | destansı | 66c | Izgara Obsidyen Ton |
| ghost_marlin | Hayalet Marlin | fishing/fish | 65 | T5 | destansı | 88c | Ruh Marlini Tabağı |
| golden_carp | İmparatorluk Altın Koisi | fishing/fish | 70 | T5 | destansı | 1s 18c | Izgara İmparatorluk Altın Koisi |
| crystal_serpentfish | Kristal Yılanbalığı | fishing/fish | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı |
| void_angler | Boşluk Fenerbalığı | fishing/fish | 80 | T6 | efsanevi | 2s 12c | Izgara Boşluk Fenerbalığı |
| astral_moonfish | Yıldız Aybalığı | fishing/fish | 85 | T6 | efsanevi | 2s 83c | Yıldız Aybalığı Lezzeti |
| phoenix_koi | Anka Koisi | fishing/fish | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti |
| abyssal_dragonfish | Uçurum Ejderbalığı | fishing/fish | 95 | T7 | mitik | 5s 7c | Izgara Uçurum Ejderbalığı |
| ancient_leviathan_scale | Göksel Leviathan | fishing/fish | 100 | T7 | mitik | 6s 79c | Leviathan Kraliyet Şöleni |
| common_grass | Sıradan Ot | herbalism/herb | 1 | T0 | sıradan | 2c | Acemi Cübbe, Acemi Kukuleta, Acemi Kutsal Kitap, Acemi Sandalet, Acemi Sargı, Acemi Şalvar, Deri Bacaklık, Deri Eldiven, Deri Göğüslük, Deri Kask, Deri Çizme, Izgara Lepistes, Kalfa Cübbe, Kalfa Deri Başlık, Kalfa Deri Eldiven, Kalfa Deri Tozluk, Kalfa Deri Yelek, Kalfa Deri Çizme, Kalfa Kukuleta, Kalfa Kutsal Kitap, Kalfa Sandalet, Kalfa Sargı, Kalfa Şalvar, Kaliteli Cübbe, Kaliteli Deri Başlık, Kaliteli Deri Eldiven, Kaliteli Deri Tozluk, Kaliteli Deri Yelek, Kaliteli Deri Çizme, Kaliteli Kukuleta, Kaliteli Kutsal Kitap, Kaliteli Sandalet, Kaliteli Sargı, Kaliteli Şalvar, Küçük Can İksiri, Kıdemli Cübbe, Kıdemli Deri Başlık, Kıdemli Deri Eldiven, Kıdemli Deri Tozluk, Kıdemli Deri Yelek, Kıdemli Deri Çizme, Kıdemli Kukuleta, Kıdemli Kutsal Kitap, Kıdemli Sandalet, Kıdemli Sargı, Kıdemli Şalvar, Şampiyon Cübbe, Şampiyon Deri Başlık, Şampiyon Deri Eldiven, Şampiyon Deri Tozluk, Şampiyon Deri Yelek, Şampiyon Deri Çizme, Şampiyon Kukuleta, Şampiyon Kutsal Kitap, Şampiyon Sandalet, Şampiyon Sargı, Şampiyon Şalvar |
| dandelion | Karahindiba | herbalism/herb | 4 | T0 | sıradan | 3c | Adamant Cübbe, Adamant Deri Başlık, Adamant Deri Eldiven, Adamant Deri Tozluk, Adamant Deri Yelek, Adamant Deri Çizme, Adamant Kukuleta, Adamant Kutsal Kitap, Adamant Sandalet, Adamant Sargı, Adamant Şalvar, Mithril Cübbe, Mithril Deri Başlık, Mithril Deri Eldiven, Mithril Deri Tozluk, Mithril Deri Yelek, Mithril Deri Çizme, Mithril Kukuleta, Mithril Kutsal Kitap, Mithril Sandalet, Mithril Sargı, Mithril Şalvar, Orta Boy Can İksiri, Rünlü Cübbe, Rünlü Deri Başlık, Rünlü Deri Eldiven, Rünlü Deri Tozluk, Rünlü Deri Yelek, Rünlü Deri Çizme, Rünlü Kukuleta, Rünlü Kutsal Kitap, Rünlü Sandalet, Rünlü Sargı, Rünlü Şalvar, Sardalya Şiş |
| wild_mint | Yabani Nane | herbalism/herb | 7 | T0 | sıradan | 3c | Ejder Dövme Cübbe, Ejder Dövme Deri Başlık, Ejder Dövme Deri Eldiven, Ejder Dövme Deri Tozluk, Ejder Dövme Deri Yelek, Ejder Dövme Deri Çizme, Ejder Dövme Kukuleta, Ejder Dövme Kutsal Kitap, Ejder Dövme Sandalet, Ejder Dövme Sargı, Ejder Dövme Şalvar, Naneli Balık Çorbası, Yıldız Dövme Cübbe, Yıldız Dövme Deri Başlık, Yıldız Dövme Deri Eldiven, Yıldız Dövme Deri Tozluk, Yıldız Dövme Deri Yelek, Yıldız Dövme Deri Çizme, Yıldız Dövme Kukuleta, Yıldız Dövme Kutsal Kitap, Yıldız Dövme Sandalet, Yıldız Dövme Sargı, Yıldız Dövme Şalvar |
| chamomile | Papatya | herbalism/herb | 10 | T1 | sıradan | 4c | Orta Boy Can İksiri, Titan Cübbe, Titan Deri Başlık, Titan Deri Eldiven, Titan Deri Tozluk, Titan Deri Yelek, Titan Deri Çizme, Titan Kukuleta, Titan Kutsal Kitap, Titan Sandalet, Titan Sargı, Titan Şalvar, Turna Fileto |
| lavender | Lavanta | herbalism/herb | 12 | T1 | sıradan | 4c | Büyük Boy Can İksiri, Kraliyet Cübbe, Kraliyet Deri Başlık, Kraliyet Deri Eldiven, Kraliyet Deri Tozluk, Kraliyet Deri Yelek, Kraliyet Deri Çizme, Kraliyet Kukuleta, Kraliyet Kutsal Kitap, Kraliyet Sandalet, Kraliyet Sargı, Kraliyet Şalvar, Lavantalı Sazan |
| sunflower | Ayçiçeği | herbalism/herb | 14 | T1 | sıradan | 5c | Güneş Dövme Cübbe, Güneş Dövme Deri Başlık, Güneş Dövme Deri Eldiven, Güneş Dövme Deri Tozluk, Güneş Dövme Deri Yelek, Güneş Dövme Deri Çizme, Güneş Dövme Kukuleta, Güneş Dövme Kutsal Kitap, Güneş Dövme Sandalet, Güneş Dövme Sargı, Güneş Dövme Şalvar, Izgara Gümüş Ringa |
| bitterroot | Acıkök | herbalism/herb | 16 | T1 | sıradan | 5c | Obsidyen Cübbe, Obsidyen Deri Başlık, Obsidyen Deri Eldiven, Obsidyen Deri Tozluk, Obsidyen Deri Yelek, Obsidyen Deri Çizme, Obsidyen Kukuleta, Obsidyen Kutsal Kitap, Obsidyen Sandalet, Obsidyen Sargı, Obsidyen Şalvar, Otlu Alabalık |
| healing_moss | Şifalı Yosun | herbalism/herb | 18 | T1 | sıradan | 6c | Füme Somon |
| marsh_reed | Bataklık Sazı | herbalism/herb | 20 | T2 | nadir-altı | 6c | Kum Çeliği Cübbe, Kum Çeliği Deri Başlık, Kum Çeliği Deri Eldiven, Kum Çeliği Deri Tozluk, Kum Çeliği Deri Yelek, Kum Çeliği Deri Çizme, Kum Çeliği Kukuleta, Kum Çeliği Kutsal Kitap, Kum Çeliği Sandalet, Kum Çeliği Sargı, Kum Çeliği Şalvar, Levrek Güveci |
| blue_thistle | Mavi Devedikeni | herbalism/herb | 22 | T2 | nadir-altı | 7c | Ekstra Büyük Can İksiri, Izgara Göl Levreği, Rün Bağlı Cübbe, Rün Bağlı Deri Başlık, Rün Bağlı Deri Eldiven, Rün Bağlı Deri Tozluk, Rün Bağlı Deri Yelek, Rün Bağlı Deri Çizme, Rün Bağlı Kukuleta, Rün Bağlı Kutsal Kitap, Rün Bağlı Sandalet, Rün Bağlı Sargı, Rün Bağlı Şalvar |
| red_sage | Kızıl Adaçayı | herbalism/herb | 24 | T2 | nadir-altı | 8c | Anka Cübbe, Anka Deri Başlık, Anka Deri Eldiven, Anka Deri Tozluk, Anka Deri Yelek, Anka Deri Çizme, Anka Kukuleta, Anka Kutsal Kitap, Anka Sandalet, Anka Sargı, Anka Şalvar, Büyük Boy Can İksiri, Izgara Kırmızı Mercan |
| nightshade | Güzelavrat Otu | herbalism/herb | 26 | T2 | nadir-altı | 9c | Dehşet Dövme Cübbe, Dehşet Dövme Deri Başlık, Dehşet Dövme Deri Eldiven, Dehşet Dövme Deri Tozluk, Dehşet Dövme Deri Yelek, Dehşet Dövme Deri Çizme, Dehşet Dövme Kukuleta, Dehşet Dövme Kutsal Kitap, Dehşet Dövme Sandalet, Dehşet Dövme Sargı, Dehşet Dövme Şalvar, Izgara Nehir Yayını |
| wild_ginseng | Yaban Ginsengi | herbalism/herb | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| golden_clover | Altın Yonca | herbalism/herb | 30 | T3 | nadir-altı | 11c | Hükümdar Cübbe, Hükümdar Deri Başlık, Hükümdar Deri Eldiven, Hükümdar Deri Tozluk, Hükümdar Deri Yelek, Hükümdar Deri Çizme, Hükümdar Kukuleta, Hükümdar Kutsal Kitap, Hükümdar Sandalet, Hükümdar Sargı, Hükümdar Şalvar, Izgara Gökkuşağı Alabalığı, Mercan Cübbe, Mercan Deri Başlık, Mercan Deri Eldiven, Mercan Deri Tozluk, Mercan Deri Yelek, Mercan Deri Çizme, Mercan Kukuleta, Mercan Kutsal Kitap, Mercan Sandalet, Mercan Sargı, Mercan Şalvar, Çok Büyük Can İksiri |
| dreamroot | Düş Kökü | herbalism/herb | 33 | T3 | nadir-altı | 14c | Ejderha Kemiği Cübbe, Ejderha Kemiği Deri Başlık, Ejderha Kemiği Deri Eldiven, Ejderha Kemiği Deri Tozluk, Ejderha Kemiği Deri Yelek, Ejderha Kemiği Deri Çizme, Ejderha Kemiği Kukuleta, Ejderha Kemiği Kutsal Kitap, Ejderha Kemiği Sandalet, Ejderha Kemiği Sargı, Ejderha Kemiği Şalvar, Kraliyet Sazanı Ziyafeti |
| silverleaf | Gümüş Yaprak | herbalism/herb | 36 | T3 | nadir-altı | 16c | Izgara Zümrüt Turnası, Sis Örgü Cübbe, Sis Örgü Deri Başlık, Sis Örgü Deri Eldiven, Sis Örgü Deri Tozluk, Sis Örgü Deri Yelek, Sis Örgü Deri Çizme, Sis Örgü Kukuleta, Sis Örgü Kutsal Kitap, Sis Örgü Sandalet, Sis Örgü Sargı, Sis Örgü Şalvar |
| bloodthorn | Kan Dikeni | herbalism/herb | 39 | T3 | nadir-altı | 19c | Izgara Gümüş Barakuda, Prizmatik Cübbe, Prizmatik Deri Başlık, Prizmatik Deri Eldiven, Prizmatik Deri Tozluk, Prizmatik Deri Yelek, Prizmatik Deri Çizme, Prizmatik Kukuleta, Prizmatik Kutsal Kitap, Prizmatik Sandalet, Prizmatik Sargı, Prizmatik Şalvar |
| ghost_mushroom | Hayalet Mantarı | herbalism/herb | 42 | T3 | nadir-altı | 23c | Magma Dövme Cübbe, Magma Dövme Deri Başlık, Magma Dövme Deri Eldiven, Magma Dövme Deri Tozluk, Magma Dövme Deri Yelek, Magma Dövme Deri Çizme, Magma Dövme Kukuleta, Magma Dövme Kutsal Kitap, Magma Dövme Sandalet, Magma Dövme Sargı, Magma Dövme Şalvar, Mantarlı Mersin Turtası, Yüce Can İksiri |
| stormflower | Fırtına Çiçeği | herbalism/herb | 45 | T4 | nadir | 28c | Bora Cübbe, Bora Deri Başlık, Bora Deri Eldiven, Bora Deri Tozluk, Bora Deri Yelek, Bora Deri Çizme, Bora Kukuleta, Bora Kutsal Kitap, Bora Sandalet, Bora Sargı, Bora Şalvar, Ekstra Büyük Can İksiri, Fırtına Izgarası Kılıçbalığı, Hiçlik Cübbe, Hiçlik Deri Başlık, Hiçlik Deri Eldiven, Hiçlik Deri Tozluk, Hiçlik Deri Yelek, Hiçlik Deri Çizme, Hiçlik Kukuleta, Hiçlik Kutsal Kitap, Hiçlik Sandalet, Hiçlik Sargı, Hiçlik Şalvar |
| emberleaf | Kor Yaprağı | herbalism/herb | 50 | T4 | nadir | 37c | Izgara Fırtına Pullu Yılanbalığı, Melek Cübbe, Melek Deri Başlık, Melek Deri Eldiven, Melek Deri Tozluk, Melek Deri Yelek, Melek Deri Çizme, Melek Kukuleta, Melek Kutsal Kitap, Melek Sandalet, Melek Sargı, Melek Şalvar, Tutulma Cübbe, Tutulma Deri Başlık, Tutulma Deri Eldiven, Tutulma Deri Tozluk, Tutulma Deri Yelek, Tutulma Deri Çizme, Tutulma Kukuleta, Tutulma Kutsal Kitap, Tutulma Sandalet, Tutulma Sargı, Tutulma Şalvar |
| frostbloom | Ayaz Çiçeği | herbalism/herb | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası, Ejder Lordu Cübbe, Ejder Lordu Deri Başlık, Ejder Lordu Deri Eldiven, Ejder Lordu Deri Tozluk, Ejder Lordu Deri Yelek, Ejder Lordu Deri Çizme, Ejder Lordu Kukuleta, Ejder Lordu Kutsal Kitap, Ejder Lordu Sandalet, Ejder Lordu Sargı, Ejder Lordu Şalvar, Göksel Cübbe, Göksel Deri Başlık, Göksel Deri Eldiven, Göksel Deri Tozluk, Göksel Deri Yelek, Göksel Deri Çizme, Göksel Kukuleta, Göksel Kutsal Kitap, Göksel Sandalet, Göksel Sargı, Göksel Şalvar |
| moonshade | Ay Gölgesi | herbalism/herb | 60 | T5 | destansı | 66c | Ebedi Cübbe, Ebedi Deri Başlık, Ebedi Deri Eldiven, Ebedi Deri Tozluk, Ebedi Deri Yelek, Ebedi Deri Çizme, Ebedi Kukuleta, Ebedi Kutsal Kitap, Ebedi Sandalet, Ebedi Sargı, Ebedi Şalvar, Izgara Obsidyen Ton, Yıldızsız Cübbe, Yıldızsız Deri Başlık, Yıldızsız Deri Eldiven, Yıldızsız Deri Tozluk, Yıldızsız Deri Yelek, Yıldızsız Deri Çizme, Yıldızsız Kukuleta, Yıldızsız Kutsal Kitap, Yıldızsız Sandalet, Yıldızsız Sargı, Yıldızsız Şalvar |
| spirit_lotus | Ruh Nilüferi | herbalism/herb | 65 | T5 | destansı | 88c | Gece Çeliği Cübbe, Gece Çeliği Deri Başlık, Gece Çeliği Deri Eldiven, Gece Çeliği Deri Tozluk, Gece Çeliği Deri Yelek, Gece Çeliği Deri Çizme, Gece Çeliği Kukuleta, Gece Çeliği Kutsal Kitap, Gece Çeliği Sandalet, Gece Çeliği Sargı, Gece Çeliği Şalvar, Ruh Marlini Tabağı, Zaman Dövme Cübbe, Zaman Dövme Deri Başlık, Zaman Dövme Deri Eldiven, Zaman Dövme Deri Tozluk, Zaman Dövme Deri Yelek, Zaman Dövme Deri Çizme, Zaman Dövme Kukuleta, Zaman Dövme Kutsal Kitap, Zaman Dövme Sandalet, Zaman Dövme Sargı, Zaman Dövme Şalvar, Çok Büyük Can İksiri |
| ancient_mandrake | Kadim Adamotu | herbalism/herb | 70 | T5 | destansı | 1s 18c | Fırtına Doğumlu Cübbe, Fırtına Doğumlu Deri Başlık, Fırtına Doğumlu Deri Eldiven, Fırtına Doğumlu Deri Tozluk, Fırtına Doğumlu Deri Yelek, Fırtına Doğumlu Deri Çizme, Fırtına Doğumlu Kukuleta, Fırtına Doğumlu Kutsal Kitap, Fırtına Doğumlu Sandalet, Fırtına Doğumlu Sargı, Fırtına Doğumlu Şalvar, Izgara İmparatorluk Altın Koisi, Yükselmiş Cübbe, Yükselmiş Deri Başlık, Yükselmiş Deri Eldiven, Yükselmiş Deri Tozluk, Yükselmiş Deri Yelek, Yükselmiş Deri Çizme, Yükselmiş Kukuleta, Yükselmiş Kutsal Kitap, Yükselmiş Sandalet, Yükselmiş Sargı, Yükselmiş Şalvar |
| dragonblood_herb | Ejderkanı Otu | herbalism/herb | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı, Kaos Cübbe, Kaos Deri Başlık, Kaos Deri Eldiven, Kaos Deri Tozluk, Kaos Deri Yelek, Kaos Deri Çizme, Kaos Kukuleta, Kaos Kutsal Kitap, Kaos Sandalet, Kaos Sargı, Kaos Şalvar, Yarık Dövme Cübbe, Yarık Dövme Deri Başlık, Yarık Dövme Deri Eldiven, Yarık Dövme Deri Tozluk, Yarık Dövme Deri Yelek, Yarık Dövme Deri Çizme, Yarık Dövme Kukuleta, Yarık Dövme Kutsal Kitap, Yarık Dövme Sandalet, Yarık Dövme Sargı, Yarık Dövme Şalvar |
| astral_bloom | Yıldız Çiçeği | herbalism/herb | 80 | T6 | efsanevi | 2s 12c | Başmelek Cübbe, Başmelek Deri Başlık, Başmelek Deri Eldiven, Başmelek Deri Tozluk, Başmelek Deri Yelek, Başmelek Deri Çizme, Başmelek Kukuleta, Başmelek Kutsal Kitap, Başmelek Sandalet, Başmelek Sargı, Başmelek Şalvar, Gök Dövme Cübbe, Gök Dövme Deri Başlık, Gök Dövme Deri Eldiven, Gök Dövme Deri Tozluk, Gök Dövme Deri Yelek, Gök Dövme Deri Çizme, Gök Dövme Kukuleta, Gök Dövme Kutsal Kitap, Gök Dövme Sandalet, Gök Dövme Sargı, Gök Dövme Şalvar, Izgara Boşluk Fenerbalığı |
| void_orchid | Boşluk Orkidesi | herbalism/herb | 85 | T6 | efsanevi | 2s 83c | Kül Cübbe, Kül Deri Başlık, Kül Deri Eldiven, Kül Deri Tozluk, Kül Deri Yelek, Kül Deri Çizme, Kül Kukuleta, Kül Kutsal Kitap, Kül Sandalet, Kül Sargı, Kül Şalvar, Yüce Can İksiri, Yıldız Aybalığı Lezzeti, Şafak Cübbe, Şafak Deri Başlık, Şafak Deri Eldiven, Şafak Deri Tozluk, Şafak Deri Yelek, Şafak Deri Çizme, Şafak Kukuleta, Şafak Kutsal Kitap, Şafak Sandalet, Şafak Sargı, Şafak Şalvar |
| celestial_lotus | Göksel Nilüfer | herbalism/herb | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti, Kıyamet Dövme Cübbe, Kıyamet Dövme Deri Başlık, Kıyamet Dövme Deri Eldiven, Kıyamet Dövme Deri Tozluk, Kıyamet Dövme Deri Yelek, Kıyamet Dövme Deri Çizme, Kıyamet Dövme Kukuleta, Kıyamet Dövme Kutsal Kitap, Kıyamet Dövme Sandalet, Kıyamet Dövme Sargı, Kıyamet Dövme Şalvar, Leviathan Cübbe, Leviathan Deri Başlık, Leviathan Deri Eldiven, Leviathan Deri Tozluk, Leviathan Deri Yelek, Leviathan Deri Çizme, Leviathan Kukuleta, Leviathan Kutsal Kitap, Leviathan Sandalet, Leviathan Sargı, Leviathan Şalvar |
| phoenix_petal | Anka Çiçeği | herbalism/herb | 95 | T7 | mitik | 5s 7c | Hiçlik Dövme Cübbe, Hiçlik Dövme Deri Başlık, Hiçlik Dövme Deri Eldiven, Hiçlik Dövme Deri Tozluk, Hiçlik Dövme Deri Yelek, Hiçlik Dövme Deri Çizme, Hiçlik Dövme Kukuleta, Hiçlik Dövme Kutsal Kitap, Hiçlik Dövme Sandalet, Hiçlik Dövme Sargı, Hiçlik Dövme Şalvar, Izgara Uçurum Ejderbalığı, Tanrı Katili Cübbe, Tanrı Katili Deri Başlık, Tanrı Katili Deri Eldiven, Tanrı Katili Deri Tozluk, Tanrı Katili Deri Yelek, Tanrı Katili Deri Çizme, Tanrı Katili Kukuleta, Tanrı Katili Kutsal Kitap, Tanrı Katili Sandalet, Tanrı Katili Sargı, Tanrı Katili Şalvar |
| worldroot | Dünya Kökü | herbalism/herb | 100 | T7 | mitik | 6s 79c | Kadim Çağ Cübbe, Kadim Çağ Deri Başlık, Kadim Çağ Deri Eldiven, Kadim Çağ Deri Tozluk, Kadim Çağ Deri Yelek, Kadim Çağ Deri Çizme, Kadim Çağ Kukuleta, Kadim Çağ Kutsal Kitap, Kadim Çağ Sandalet, Kadim Çağ Sargı, Kadim Çağ Şalvar, Leviathan Kraliyet Şöleni |
| forest_rat_pelt | Orman Faresi Postu | hunting/pelt | 1 | T0 | sıradan | 2c | Acemi Bacaklık, Acemi Cübbe, Acemi Eldiven, Acemi Göğüslük, Acemi Kask, Acemi Kukuleta, Acemi Kutsal Kitap, Acemi Sandalet, Acemi Sargı, Acemi Çizme, Acemi Şalvar, Demir Bacaklık, Demir Eldiven, Demir Göğüslük, Demir Kask, Demir Çizme, Deri Bacaklık, Deri Eldiven, Deri Göğüslük, Deri Kask, Deri Çizme, Kalfa Cübbe, Kalfa Deri Başlık, Kalfa Deri Eldiven, Kalfa Deri Tozluk, Kalfa Deri Yelek, Kalfa Deri Çizme, Kalfa Kukuleta, Kalfa Kutsal Kitap, Kalfa Sandalet, Kalfa Sargı, Kalfa Şalvar, Kaliteli Cübbe, Kaliteli Deri Başlık, Kaliteli Deri Eldiven, Kaliteli Deri Tozluk, Kaliteli Deri Yelek, Kaliteli Deri Çizme, Kaliteli Kukuleta, Kaliteli Kutsal Kitap, Kaliteli Sandalet, Kaliteli Sargı, Kaliteli Şalvar, Kıdemli Cübbe, Kıdemli Deri Başlık, Kıdemli Deri Eldiven, Kıdemli Deri Tozluk, Kıdemli Deri Yelek, Kıdemli Deri Çizme, Kıdemli Kukuleta, Kıdemli Kutsal Kitap, Kıdemli Sandalet, Kıdemli Sargı, Kıdemli Şalvar, Pul Bacaklık, Pul Eldiven, Pul Göğüslük, Pul Kask, Pul Çizme, Çelik Bacaklık, Çelik Eldiven, Çelik Göğüslük, Çelik Kask, Çelik Çizme, Şampiyon Cübbe, Şampiyon Deri Başlık, Şampiyon Deri Eldiven, Şampiyon Deri Tozluk, Şampiyon Deri Yelek, Şampiyon Deri Çizme, Şampiyon Kukuleta, Şampiyon Kutsal Kitap, Şampiyon Sandalet, Şampiyon Sargı, Şampiyon Şalvar, Şövalye Bacaklığı, Şövalye Eldiveni, Şövalye Göğüslüğü, Şövalye Kaskı, Şövalye Çizmesi |
| rabbit_pelt | Küçük Tavşan Postu | hunting/pelt | 4 | T0 | sıradan | 3c | Adamant Bacaklık, Adamant Cübbe, Adamant Deri Başlık, Adamant Deri Eldiven, Adamant Deri Tozluk, Adamant Deri Yelek, Adamant Deri Çizme, Adamant Eldiven, Adamant Göğüslük, Adamant Kask, Adamant Kukuleta, Adamant Kutsal Kitap, Adamant Sandalet, Adamant Sargı, Adamant Çizme, Adamant Şalvar, Mithril Bacaklık, Mithril Cübbe, Mithril Deri Başlık, Mithril Deri Eldiven, Mithril Deri Tozluk, Mithril Deri Yelek, Mithril Deri Çizme, Mithril Eldiven, Mithril Göğüslük, Mithril Kask, Mithril Kukuleta, Mithril Kutsal Kitap, Mithril Sandalet, Mithril Sargı, Mithril Çizme, Mithril Şalvar, Rünlü Bacaklık, Rünlü Cübbe, Rünlü Deri Başlık, Rünlü Deri Eldiven, Rünlü Deri Tozluk, Rünlü Deri Yelek, Rünlü Deri Çizme, Rünlü Eldiven, Rünlü Göğüslük, Rünlü Kask, Rünlü Kukuleta, Rünlü Kutsal Kitap, Rünlü Sandalet, Rünlü Sargı, Rünlü Çizme, Rünlü Şalvar |
| field_mouse_pelt | Tarla Faresi Postu | hunting/pelt | 7 | T0 | sıradan | 3c | Ejder Dövme Bacaklık, Ejder Dövme Cübbe, Ejder Dövme Deri Başlık, Ejder Dövme Deri Eldiven, Ejder Dövme Deri Tozluk, Ejder Dövme Deri Yelek, Ejder Dövme Deri Çizme, Ejder Dövme Eldiven, Ejder Dövme Göğüslük, Ejder Dövme Kask, Ejder Dövme Kukuleta, Ejder Dövme Kutsal Kitap, Ejder Dövme Sandalet, Ejder Dövme Sargı, Ejder Dövme Çizme, Ejder Dövme Şalvar, Yıldız Dövme Bacaklık, Yıldız Dövme Cübbe, Yıldız Dövme Deri Başlık, Yıldız Dövme Deri Eldiven, Yıldız Dövme Deri Tozluk, Yıldız Dövme Deri Yelek, Yıldız Dövme Deri Çizme, Yıldız Dövme Eldiven, Yıldız Dövme Göğüslük, Yıldız Dövme Kask, Yıldız Dövme Kukuleta, Yıldız Dövme Kutsal Kitap, Yıldız Dövme Sandalet, Yıldız Dövme Sargı, Yıldız Dövme Çizme, Yıldız Dövme Şalvar |
| squirrel_pelt | Sincap Postu | hunting/pelt | 10 | T1 | sıradan | 4c | Titan Bacaklık, Titan Cübbe, Titan Deri Başlık, Titan Deri Eldiven, Titan Deri Tozluk, Titan Deri Yelek, Titan Deri Çizme, Titan Eldiven, Titan Göğüslük, Titan Kask, Titan Kukuleta, Titan Kutsal Kitap, Titan Sandalet, Titan Sargı, Titan Çizme, Titan Şalvar |
| hare_pelt | Yaban Tavşanı Postu | hunting/pelt | 12 | T1 | sıradan | 4c | Kraliyet Bacaklık, Kraliyet Cübbe, Kraliyet Deri Başlık, Kraliyet Deri Eldiven, Kraliyet Deri Tozluk, Kraliyet Deri Yelek, Kraliyet Deri Çizme, Kraliyet Eldiven, Kraliyet Göğüslük, Kraliyet Kalkan, Kraliyet Kask, Kraliyet Kukuleta, Kraliyet Kutsal Kitap, Kraliyet Sandalet, Kraliyet Sargı, Kraliyet Çizme, Kraliyet Şalvar |
| badger_pelt | Porsuk Postu | hunting/pelt | 14 | T1 | sıradan | 5c | Güneş Dövme Bacaklık, Güneş Dövme Cübbe, Güneş Dövme Deri Başlık, Güneş Dövme Deri Eldiven, Güneş Dövme Deri Tozluk, Güneş Dövme Deri Yelek, Güneş Dövme Deri Çizme, Güneş Dövme Eldiven, Güneş Dövme Göğüslük, Güneş Dövme Kalkan, Güneş Dövme Kask, Güneş Dövme Kukuleta, Güneş Dövme Kutsal Kitap, Güneş Dövme Sandalet, Güneş Dövme Sargı, Güneş Dövme Çizme, Güneş Dövme Şalvar, Obsidyen Bacaklık, Obsidyen Cübbe, Obsidyen Deri Başlık, Obsidyen Deri Eldiven, Obsidyen Deri Tozluk, Obsidyen Deri Yelek, Obsidyen Deri Çizme, Obsidyen Eldiven, Obsidyen Göğüslük, Obsidyen Kalkan, Obsidyen Kask, Obsidyen Kukuleta, Obsidyen Kutsal Kitap, Obsidyen Sandalet, Obsidyen Sargı, Obsidyen Çizme, Obsidyen Şalvar |
| wild_turkey_feather | Yaban Hindisi Tüyü | hunting/feather | 16 | T1 | sıradan | 5c | Acemi Arbalet, Adamant Arbalet, Adamant Yay, Anka Arbalet, Anka Yay, Av Yayı, Bileşik Yay, Dehşet Dövme Arbalet, Dehşet Dövme Yay, Ejder Dövme Arbalet, Ejder Dövme Yay, Ejder Kemiği Yayı, Ejderha Kemiği Arbalet, Ejderha Kemiği Yay, Güneş Dövme Arbalet, Güneş Dövme Yay, Hükümdar Arbalet, Hükümdar Yay, Kalfa Arbalet, Kaliteli Arbalet, Kraliyet Arbalet, Kraliyet Yay, Kum Çeliği Arbalet, Kum Çeliği Yay, Kıdemli Arbalet, Kısa Yay, Magma Dövme Arbalet, Magma Dövme Yay, Mercan Arbalet, Mercan Yay, Mithril Arbalet, Mithril Yay, Obsidyen Arbalet, Obsidyen Yay, Porsuk Yayı, Prizmatik Arbalet, Prizmatik Yay, Rün Bağlı Arbalet, Rün Bağlı Yay, Rünlü Arbalet, Sis Örgü Arbalet, Sis Örgü Yay, Titan Arbalet, Titan Yay, Uzun Yay, Yıldız Dövme Arbalet, Yıldız Dövme Yay, Şampiyon Arbalet |
| forest_deer_hide | Orman Geyiği Derisi | hunting/hide | 18 | T1 | sıradan | 6c | — |
| red_fox_pelt | Kızıl Tilki Postu | hunting/pelt | 20 | T2 | nadir-altı | 6c | Kum Çeliği Bacaklık, Kum Çeliği Cübbe, Kum Çeliği Deri Başlık, Kum Çeliği Deri Eldiven, Kum Çeliği Deri Tozluk, Kum Çeliği Deri Yelek, Kum Çeliği Deri Çizme, Kum Çeliği Eldiven, Kum Çeliği Göğüslük, Kum Çeliği Kalkan, Kum Çeliği Kask, Kum Çeliği Kukuleta, Kum Çeliği Kutsal Kitap, Kum Çeliği Sandalet, Kum Çeliği Sargı, Kum Çeliği Çizme, Kum Çeliği Şalvar |
| boar_hide | Yaban Domuzu Postu | hunting/hide | 22 | T2 | nadir-altı | 7c | Rün Bağlı Bacaklık, Rün Bağlı Cübbe, Rün Bağlı Deri Başlık, Rün Bağlı Deri Eldiven, Rün Bağlı Deri Tozluk, Rün Bağlı Deri Yelek, Rün Bağlı Deri Çizme, Rün Bağlı Eldiven, Rün Bağlı Göğüslük, Rün Bağlı Kalkan, Rün Bağlı Kask, Rün Bağlı Kukuleta, Rün Bağlı Kutsal Kitap, Rün Bağlı Sandalet, Rün Bağlı Sargı, Rün Bağlı Çizme, Rün Bağlı Şalvar |
| mountain_goat_hide | Dağ Keçisi Postu | hunting/hide | 24 | T2 | nadir-altı | 8c | Anka Bacaklık, Anka Cübbe, Anka Deri Başlık, Anka Deri Eldiven, Anka Deri Tozluk, Anka Deri Yelek, Anka Deri Çizme, Anka Eldiven, Anka Göğüslük, Anka Kalkan, Anka Kask, Anka Kukuleta, Anka Kutsal Kitap, Anka Sandalet, Anka Sargı, Anka Çizme, Anka Şalvar |
| wolf_pelt | Yaban Kurdu Postu | hunting/pelt | 26 | T2 | nadir-altı | 9c | Dehşet Dövme Bacaklık, Dehşet Dövme Cübbe, Dehşet Dövme Deri Başlık, Dehşet Dövme Deri Eldiven, Dehşet Dövme Deri Tozluk, Dehşet Dövme Deri Yelek, Dehşet Dövme Deri Çizme, Dehşet Dövme Eldiven, Dehşet Dövme Göğüslük, Dehşet Dövme Kalkan, Dehşet Dövme Kask, Dehşet Dövme Kukuleta, Dehşet Dövme Kutsal Kitap, Dehşet Dövme Sandalet, Dehşet Dövme Sargı, Dehşet Dövme Çizme, Dehşet Dövme Şalvar |
| forest_bear_hide | Orman Ayısı Postu | hunting/hide | 28 | T2 | nadir-altı | 10c | — |
| wild_buffalo_hide | Yaban Mandası Derisi | hunting/hide | 30 | T3 | nadir-altı | 11c | Hükümdar Bacaklık, Hükümdar Cübbe, Hükümdar Deri Başlık, Hükümdar Deri Eldiven, Hükümdar Deri Tozluk, Hükümdar Deri Yelek, Hükümdar Deri Çizme, Hükümdar Eldiven, Hükümdar Göğüslük, Hükümdar Kalkan, Hükümdar Kask, Hükümdar Kukuleta, Hükümdar Kutsal Kitap, Hükümdar Sandalet, Hükümdar Sargı, Hükümdar Çizme, Hükümdar Şalvar, Mercan Bacaklık, Mercan Cübbe, Mercan Deri Başlık, Mercan Deri Eldiven, Mercan Deri Tozluk, Mercan Deri Yelek, Mercan Deri Çizme, Mercan Eldiven, Mercan Göğüslük, Mercan Kalkan, Mercan Kask, Mercan Kukuleta, Mercan Kutsal Kitap, Mercan Sandalet, Mercan Sargı, Mercan Çizme, Mercan Şalvar |
| giant_elk_hide | Dev Geyik Postu | hunting/hide | 33 | T3 | nadir-altı | 14c | Ejderha Kemiği Bacaklık, Ejderha Kemiği Cübbe, Ejderha Kemiği Deri Başlık, Ejderha Kemiği Deri Eldiven, Ejderha Kemiği Deri Tozluk, Ejderha Kemiği Deri Yelek, Ejderha Kemiği Deri Çizme, Ejderha Kemiği Eldiven, Ejderha Kemiği Göğüslük, Ejderha Kemiği Kalkan, Ejderha Kemiği Kask, Ejderha Kemiği Kukuleta, Ejderha Kemiği Kutsal Kitap, Ejderha Kemiği Sandalet, Ejderha Kemiği Sargı, Ejderha Kemiği Çizme, Ejderha Kemiği Şalvar |
| dire_wolf_pelt | Dev Kurt Postu | hunting/pelt | 36 | T3 | nadir-altı | 16c | Sis Örgü Bacaklık, Sis Örgü Cübbe, Sis Örgü Deri Başlık, Sis Örgü Deri Eldiven, Sis Örgü Deri Tozluk, Sis Örgü Deri Yelek, Sis Örgü Deri Çizme, Sis Örgü Eldiven, Sis Örgü Göğüslük, Sis Örgü Kalkan, Sis Örgü Kask, Sis Örgü Kukuleta, Sis Örgü Kutsal Kitap, Sis Örgü Sandalet, Sis Örgü Sargı, Sis Örgü Çizme, Sis Örgü Şalvar |
| mountain_lion_pelt | Dağ Aslanı Postu | hunting/pelt | 39 | T3 | nadir-altı | 19c | Prizmatik Bacaklık, Prizmatik Cübbe, Prizmatik Deri Başlık, Prizmatik Deri Eldiven, Prizmatik Deri Tozluk, Prizmatik Deri Yelek, Prizmatik Deri Çizme, Prizmatik Eldiven, Prizmatik Göğüslük, Prizmatik Kalkan, Prizmatik Kask, Prizmatik Kukuleta, Prizmatik Kutsal Kitap, Prizmatik Sandalet, Prizmatik Sargı, Prizmatik Çizme, Prizmatik Şalvar |
| black_panther_pelt | Kara Panter Postu | hunting/pelt | 42 | T3 | nadir-altı | 23c | Bora Bacaklık, Bora Cübbe, Bora Deri Başlık, Bora Deri Eldiven, Bora Deri Tozluk, Bora Deri Yelek, Bora Deri Çizme, Bora Eldiven, Bora Göğüslük, Bora Kalkan, Bora Kask, Bora Kukuleta, Bora Kutsal Kitap, Bora Sandalet, Bora Sargı, Bora Çizme, Bora Şalvar, Hiçlik Bacaklık, Hiçlik Cübbe, Hiçlik Deri Başlık, Hiçlik Deri Eldiven, Hiçlik Deri Tozluk, Hiçlik Deri Yelek, Hiçlik Deri Çizme, Hiçlik Eldiven, Hiçlik Göğüslük, Hiçlik Kalkan, Hiçlik Kask, Hiçlik Kukuleta, Hiçlik Kutsal Kitap, Hiçlik Sandalet, Hiçlik Sargı, Hiçlik Çizme, Hiçlik Şalvar, Magma Dövme Bacaklık, Magma Dövme Cübbe, Magma Dövme Deri Başlık, Magma Dövme Deri Eldiven, Magma Dövme Deri Tozluk, Magma Dövme Deri Yelek, Magma Dövme Deri Çizme, Magma Dövme Eldiven, Magma Dövme Göğüslük, Magma Dövme Kalkan, Magma Dövme Kask, Magma Dövme Kukuleta, Magma Dövme Kutsal Kitap, Magma Dövme Sandalet, Magma Dövme Sargı, Magma Dövme Çizme, Magma Dövme Şalvar |
| storm_eagle_feather | Fırtına Kartalı Tüyü | hunting/feather | 45 | T4 | nadir | 28c | Başmelek Arbalet, Başmelek Yay, Bora Arbalet, Bora Yay, Ebedi Arbalet, Ebedi Yay, Ejder Lordu Arbalet, Ejder Lordu Yay, Fırtına Doğumlu Arbalet, Fırtına Doğumlu Yay, Gece Çeliği Arbalet, Gece Çeliği Yay, Gök Dövme Arbalet, Gök Dövme Yay, Göksel Arbalet, Göksel Yay, Hiçlik Arbalet, Hiçlik Yay, Kaos Arbalet, Kaos Yay, Melek Arbalet, Melek Yay, Tutulma Arbalet, Tutulma Yay, Yarık Dövme Arbalet, Yarık Dövme Yay, Yükselmiş Arbalet, Yükselmiş Yay, Yıldızsız Arbalet, Yıldızsız Yay, Zaman Dövme Arbalet, Zaman Dövme Yay |
| ember_boar_hide | Kor Domuzu Postu | hunting/hide | 50 | T4 | nadir | 37c | Melek Bacaklık, Melek Cübbe, Melek Deri Başlık, Melek Deri Eldiven, Melek Deri Tozluk, Melek Deri Yelek, Melek Deri Çizme, Melek Eldiven, Melek Göğüslük, Melek Kalkan, Melek Kask, Melek Kukuleta, Melek Kutsal Kitap, Melek Sandalet, Melek Sargı, Melek Çizme, Melek Şalvar, Tutulma Bacaklık, Tutulma Cübbe, Tutulma Deri Başlık, Tutulma Deri Eldiven, Tutulma Deri Tozluk, Tutulma Deri Yelek, Tutulma Deri Çizme, Tutulma Eldiven, Tutulma Göğüslük, Tutulma Kalkan, Tutulma Kask, Tutulma Kukuleta, Tutulma Kutsal Kitap, Tutulma Sandalet, Tutulma Sargı, Tutulma Çizme, Tutulma Şalvar |
| frost_wolf_pelt | Ayaz Kurdu Postu | hunting/pelt | 55 | T4 | nadir | 49c | Ejder Lordu Bacaklık, Ejder Lordu Cübbe, Ejder Lordu Deri Başlık, Ejder Lordu Deri Eldiven, Ejder Lordu Deri Tozluk, Ejder Lordu Deri Yelek, Ejder Lordu Deri Çizme, Ejder Lordu Eldiven, Ejder Lordu Göğüslük, Ejder Lordu Kalkan, Ejder Lordu Kask, Ejder Lordu Kukuleta, Ejder Lordu Kutsal Kitap, Ejder Lordu Sandalet, Ejder Lordu Sargı, Ejder Lordu Çizme, Ejder Lordu Şalvar, Göksel Bacaklık, Göksel Cübbe, Göksel Deri Başlık, Göksel Deri Eldiven, Göksel Deri Tozluk, Göksel Deri Yelek, Göksel Deri Çizme, Göksel Eldiven, Göksel Göğüslük, Göksel Kalkan, Göksel Kask, Göksel Kukuleta, Göksel Kutsal Kitap, Göksel Sandalet, Göksel Sargı, Göksel Çizme, Göksel Şalvar |
| dire_bear_hide | Dev Ayı Postu | hunting/hide | 60 | T5 | destansı | 66c | Ebedi Bacaklık, Ebedi Cübbe, Ebedi Deri Başlık, Ebedi Deri Eldiven, Ebedi Deri Tozluk, Ebedi Deri Yelek, Ebedi Deri Çizme, Ebedi Eldiven, Ebedi Göğüslük, Ebedi Kalkan, Ebedi Kask, Ebedi Kukuleta, Ebedi Kutsal Kitap, Ebedi Sandalet, Ebedi Sargı, Ebedi Çizme, Ebedi Şalvar, Yıldızsız Bacaklık, Yıldızsız Cübbe, Yıldızsız Deri Başlık, Yıldızsız Deri Eldiven, Yıldızsız Deri Tozluk, Yıldızsız Deri Yelek, Yıldızsız Deri Çizme, Yıldızsız Eldiven, Yıldızsız Göğüslük, Yıldızsız Kalkan, Yıldızsız Kask, Yıldızsız Kukuleta, Yıldızsız Kutsal Kitap, Yıldızsız Sandalet, Yıldızsız Sargı, Yıldızsız Çizme, Yıldızsız Şalvar |
| spirit_stag_hide | Ruh Geyiği Postu | hunting/hide | 65 | T5 | destansı | 88c | Gece Çeliği Bacaklık, Gece Çeliği Cübbe, Gece Çeliği Deri Başlık, Gece Çeliği Deri Eldiven, Gece Çeliği Deri Tozluk, Gece Çeliği Deri Yelek, Gece Çeliği Deri Çizme, Gece Çeliği Eldiven, Gece Çeliği Göğüslük, Gece Çeliği Kalkan, Gece Çeliği Kask, Gece Çeliği Kukuleta, Gece Çeliği Kutsal Kitap, Gece Çeliği Sandalet, Gece Çeliği Sargı, Gece Çeliği Çizme, Gece Çeliği Şalvar, Zaman Dövme Bacaklık, Zaman Dövme Cübbe, Zaman Dövme Deri Başlık, Zaman Dövme Deri Eldiven, Zaman Dövme Deri Tozluk, Zaman Dövme Deri Yelek, Zaman Dövme Deri Çizme, Zaman Dövme Eldiven, Zaman Dövme Göğüslük, Zaman Dövme Kalkan, Zaman Dövme Kask, Zaman Dövme Kukuleta, Zaman Dövme Kutsal Kitap, Zaman Dövme Sandalet, Zaman Dövme Sargı, Zaman Dövme Çizme, Zaman Dövme Şalvar |
| ancient_mammoth_hide | Kadim Mamut Postu | hunting/hide | 70 | T5 | destansı | 1s 18c | Fırtına Doğumlu Bacaklık, Fırtına Doğumlu Cübbe, Fırtına Doğumlu Deri Başlık, Fırtına Doğumlu Deri Eldiven, Fırtına Doğumlu Deri Tozluk, Fırtına Doğumlu Deri Yelek, Fırtına Doğumlu Deri Çizme, Fırtına Doğumlu Eldiven, Fırtına Doğumlu Göğüslük, Fırtına Doğumlu Kalkan, Fırtına Doğumlu Kask, Fırtına Doğumlu Kukuleta, Fırtına Doğumlu Kutsal Kitap, Fırtına Doğumlu Sandalet, Fırtına Doğumlu Sargı, Fırtına Doğumlu Çizme, Fırtına Doğumlu Şalvar, Yükselmiş Bacaklık, Yükselmiş Cübbe, Yükselmiş Deri Başlık, Yükselmiş Deri Eldiven, Yükselmiş Deri Tozluk, Yükselmiş Deri Yelek, Yükselmiş Deri Çizme, Yükselmiş Eldiven, Yükselmiş Göğüslük, Yükselmiş Kalkan, Yükselmiş Kask, Yükselmiş Kukuleta, Yükselmiş Kutsal Kitap, Yükselmiş Sandalet, Yükselmiş Sargı, Yükselmiş Çizme, Yükselmiş Şalvar |
| shadow_chimera_hide | Gölge Kimera Derisi | hunting/hide | 75 | T6 | efsanevi | 1s 58c | Kaos Bacaklık, Kaos Cübbe, Kaos Deri Başlık, Kaos Deri Eldiven, Kaos Deri Tozluk, Kaos Deri Yelek, Kaos Deri Çizme, Kaos Eldiven, Kaos Göğüslük, Kaos Kalkan, Kaos Kask, Kaos Kukuleta, Kaos Kutsal Kitap, Kaos Sandalet, Kaos Sargı, Kaos Çizme, Kaos Şalvar, Yarık Dövme Bacaklık, Yarık Dövme Cübbe, Yarık Dövme Deri Başlık, Yarık Dövme Deri Eldiven, Yarık Dövme Deri Tozluk, Yarık Dövme Deri Yelek, Yarık Dövme Deri Çizme, Yarık Dövme Eldiven, Yarık Dövme Göğüslük, Yarık Dövme Kalkan, Yarık Dövme Kask, Yarık Dövme Kukuleta, Yarık Dövme Kutsal Kitap, Yarık Dövme Sandalet, Yarık Dövme Sargı, Yarık Dövme Çizme, Yarık Dövme Şalvar |
| elder_wyvern_scale | Yaşlı Wyvern Pulu | hunting/scale | 80 | T6 | efsanevi | 2s 12c | Başmelek Bacaklık, Başmelek Cübbe, Başmelek Deri Başlık, Başmelek Deri Eldiven, Başmelek Deri Tozluk, Başmelek Deri Yelek, Başmelek Deri Çizme, Başmelek Eldiven, Başmelek Göğüslük, Başmelek Kalkan, Başmelek Kask, Başmelek Kukuleta, Başmelek Kutsal Kitap, Başmelek Sandalet, Başmelek Sargı, Başmelek Çizme, Başmelek Şalvar, Gök Dövme Bacaklık, Gök Dövme Cübbe, Gök Dövme Deri Başlık, Gök Dövme Deri Eldiven, Gök Dövme Deri Tozluk, Gök Dövme Deri Yelek, Gök Dövme Deri Çizme, Gök Dövme Eldiven, Gök Dövme Göğüslük, Gök Dövme Kalkan, Gök Dövme Kask, Gök Dövme Kukuleta, Gök Dövme Kutsal Kitap, Gök Dövme Sandalet, Gök Dövme Sargı, Gök Dövme Çizme, Gök Dövme Şalvar, Kül Bacaklık, Kül Cübbe, Kül Deri Başlık, Kül Deri Eldiven, Kül Deri Tozluk, Kül Deri Yelek, Kül Deri Çizme, Kül Eldiven, Kül Göğüslük, Kül Kalkan, Kül Kask, Kül Kukuleta, Kül Kutsal Kitap, Kül Sandalet, Kül Sargı, Kül Çizme, Kül Şalvar, Şafak Bacaklık, Şafak Cübbe, Şafak Deri Başlık, Şafak Deri Eldiven, Şafak Deri Tozluk, Şafak Deri Yelek, Şafak Deri Çizme, Şafak Eldiven, Şafak Göğüslük, Şafak Kalkan, Şafak Kask, Şafak Kukuleta, Şafak Kutsal Kitap, Şafak Sandalet, Şafak Sargı, Şafak Çizme, Şafak Şalvar |
| astral_griffin_feather | Yıldız Grifonu Tüyü | hunting/feather | 85 | T6 | efsanevi | 2s 83c | Kül Arbalet, Kül Yay, Kıyamet Dövme Arbalet, Kıyamet Dövme Yay, Leviathan Arbalet, Leviathan Yay, Şafak Arbalet, Şafak Yay |
| void_behemoth_hide | Boşluk Behemotu Derisi | hunting/hide | 90 | T7 | mitik | 3s 79c | Hiçlik Dövme Bacaklık, Hiçlik Dövme Cübbe, Hiçlik Dövme Deri Başlık, Hiçlik Dövme Deri Eldiven, Hiçlik Dövme Deri Tozluk, Hiçlik Dövme Deri Yelek, Hiçlik Dövme Deri Çizme, Hiçlik Dövme Eldiven, Hiçlik Dövme Göğüslük, Hiçlik Dövme Kalkan, Hiçlik Dövme Kask, Hiçlik Dövme Kukuleta, Hiçlik Dövme Kutsal Kitap, Hiçlik Dövme Sandalet, Hiçlik Dövme Sargı, Hiçlik Dövme Çizme, Hiçlik Dövme Şalvar, Kıyamet Dövme Bacaklık, Kıyamet Dövme Cübbe, Kıyamet Dövme Deri Başlık, Kıyamet Dövme Deri Eldiven, Kıyamet Dövme Deri Tozluk, Kıyamet Dövme Deri Yelek, Kıyamet Dövme Deri Çizme, Kıyamet Dövme Eldiven, Kıyamet Dövme Göğüslük, Kıyamet Dövme Kalkan, Kıyamet Dövme Kask, Kıyamet Dövme Kukuleta, Kıyamet Dövme Kutsal Kitap, Kıyamet Dövme Sandalet, Kıyamet Dövme Sargı, Kıyamet Dövme Çizme, Kıyamet Dövme Şalvar, Leviathan Bacaklık, Leviathan Cübbe, Leviathan Deri Başlık, Leviathan Deri Eldiven, Leviathan Deri Tozluk, Leviathan Deri Yelek, Leviathan Deri Çizme, Leviathan Eldiven, Leviathan Göğüslük, Leviathan Kalkan, Leviathan Kask, Leviathan Kukuleta, Leviathan Kutsal Kitap, Leviathan Sandalet, Leviathan Sargı, Leviathan Çizme, Leviathan Şalvar, Tanrı Katili Bacaklık, Tanrı Katili Cübbe, Tanrı Katili Deri Başlık, Tanrı Katili Deri Eldiven, Tanrı Katili Deri Tozluk, Tanrı Katili Deri Yelek, Tanrı Katili Deri Çizme, Tanrı Katili Eldiven, Tanrı Katili Göğüslük, Tanrı Katili Kalkan, Tanrı Katili Kask, Tanrı Katili Kukuleta, Tanrı Katili Kutsal Kitap, Tanrı Katili Sandalet, Tanrı Katili Sargı, Tanrı Katili Çizme, Tanrı Katili Şalvar |
| ancient_phoenix_feather | Kadim Anka Tüyü | hunting/feather | 95 | T7 | mitik | 5s 7c | Hiçlik Dövme Arbalet, Hiçlik Dövme Yay, Kadim Çağ Arbalet, Kadim Çağ Yay, Tanrı Katili Arbalet, Tanrı Katili Yay |
| dragon_scale | Göksel Ejder Pulu | hunting/scale | 100 | T7 | mitik | 6s 79c | Kadim Çağ Bacaklık, Kadim Çağ Cübbe, Kadim Çağ Deri Başlık, Kadim Çağ Deri Eldiven, Kadim Çağ Deri Tozluk, Kadim Çağ Deri Yelek, Kadim Çağ Deri Çizme, Kadim Çağ Eldiven, Kadim Çağ Göğüslük, Kadim Çağ Kalkan, Kadim Çağ Kask, Kadim Çağ Kukuleta, Kadim Çağ Kutsal Kitap, Kadim Çağ Sandalet, Kadim Çağ Sargı, Kadim Çağ Çizme, Kadim Çağ Şalvar |
| raw_stone | Ham Taş | mining/stone | 1 | T0 | sıradan | 2c | — |
| stone_ore | Taş Cevheri | mining/stone | 4 | T0 | sıradan | 3c | Adamant Kolye, Adamant Tılsım, Adamant Yüzük, Mithril Kolye, Mithril Tılsım, Mithril Yüzük, Rünlü Kolye, Rünlü Tılsım, Rünlü Yüzük |
| flint | Çakmaktaşı | mining/stone | 7 | T0 | sıradan | 3c | Ejder Dövme Kolye, Ejder Dövme Tılsım, Ejder Dövme Yüzük, Güneş Dövme Kolye, Güneş Dövme Tılsım, Güneş Dövme Yüzük, Kraliyet Kolye, Kraliyet Tılsım, Kraliyet Yüzük, Kum Çeliği Kolye, Kum Çeliği Tılsım, Kum Çeliği Yüzük, Obsidyen Kolye, Obsidyen Tılsım, Obsidyen Yüzük, Rün Bağlı Kolye, Rün Bağlı Tılsım, Rün Bağlı Yüzük, Titan Kolye, Titan Tılsım, Titan Yüzük, Yıldız Dövme Kolye, Yıldız Dövme Tılsım, Yıldız Dövme Yüzük |
| coal | Kömür | mining/fuel | 10 | T1 | sıradan | 4c | Adamantit Külçesi, Altın Külçesi, Alüminyum Külçesi, Bakır Külçesi, Bronz Külçesi, Demir Külçesi, Elektrum Külçesi, Gümüş Külçesi, Kalay Külçesi, Kobalt Külçesi, Krom Külçesi, Kurşun Külçesi, Mitril Külçesi, Nikel Külçesi, Orikalkum Külçesi, Osmiyum Külçesi, Paladyum Külçesi, Platin Külçesi, Rodyum Külçesi, Titanyum Külçesi, Tungsten Külçesi, Yıldız Çeliği Külçesi, Çinko Külçesi, İridyum Külçesi |
| aluminium_ore | Alüminyum Cevheri | mining/ore | 12 | T1 | sıradan | 4c | Alüminyum Külçesi |
| lead_ore | Kurşun Cevheri | mining/ore | 14 | T1 | sıradan | 5c | Kurşun Külçesi |
| iron_ore | Demir Cevheri | mining/ore | 16 | T1 | sıradan | 5c | Demir Külçesi |
| tin_ore | Kalay Cevheri | mining/ore | 18 | T1 | sıradan | 6c | Kalay Külçesi |
| copper_ore | Bakır Cevheri | mining/ore | 20 | T2 | nadir-altı | 6c | Bakır Külçesi |
| bronze_ore | Bronz Cevheri | mining/ore | 22 | T2 | nadir-altı | 7c | Bronz Külçesi |
| obsidian | Obsidyen | mining/stone | 24 | T2 | nadir-altı | 8c | Anka Kolye, Anka Tılsım, Anka Yüzük, Başmelek Kolye, Başmelek Tılsım, Başmelek Yüzük, Bora Kolye, Bora Tılsım, Bora Yüzük, Dehşet Dövme Kolye, Dehşet Dövme Tılsım, Dehşet Dövme Yüzük, Ebedi Kolye, Ebedi Tılsım, Ebedi Yüzük, Ejder Lordu Kolye, Ejder Lordu Tılsım, Ejder Lordu Yüzük, Ejderha Kemiği Kolye, Ejderha Kemiği Tılsım, Ejderha Kemiği Yüzük, Fırtına Doğumlu Kolye, Fırtına Doğumlu Tılsım, Fırtına Doğumlu Yüzük, Gece Çeliği Kolye, Gece Çeliği Tılsım, Gece Çeliği Yüzük, Gök Dövme Kolye, Gök Dövme Tılsım, Gök Dövme Yüzük, Göksel Kolye, Göksel Tılsım, Göksel Yüzük, Hiçlik Dövme Kolye, Hiçlik Dövme Tılsım, Hiçlik Dövme Yüzük, Hiçlik Kolye, Hiçlik Tılsım, Hiçlik Yüzük, Hükümdar Kolye, Hükümdar Tılsım, Hükümdar Yüzük, Kadim Çağ Kolye, Kadim Çağ Tılsım, Kadim Çağ Yüzük, Kaos Kolye, Kaos Tılsım, Kaos Yüzük, Kül Kolye, Kül Tılsım, Kül Yüzük, Kıyamet Dövme Kolye, Kıyamet Dövme Tılsım, Kıyamet Dövme Yüzük, Leviathan Kolye, Leviathan Tılsım, Leviathan Yüzük, Magma Dövme Kolye, Magma Dövme Tılsım, Magma Dövme Yüzük, Melek Kolye, Melek Tılsım, Melek Yüzük, Mercan Kolye, Mercan Tılsım, Mercan Yüzük, Prizmatik Kolye, Prizmatik Tılsım, Prizmatik Yüzük, Sis Örgü Kolye, Sis Örgü Tılsım, Sis Örgü Yüzük, Tanrı Katili Kolye, Tanrı Katili Tılsım, Tanrı Katili Yüzük, Tutulma Kolye, Tutulma Tılsım, Tutulma Yüzük, Yarık Dövme Kolye, Yarık Dövme Tılsım, Yarık Dövme Yüzük, Yükselmiş Kolye, Yükselmiş Tılsım, Yükselmiş Yüzük, Yıldızsız Kolye, Yıldızsız Tılsım, Yıldızsız Yüzük, Zaman Dövme Kolye, Zaman Dövme Tılsım, Zaman Dövme Yüzük, Şafak Kolye, Şafak Tılsım, Şafak Yüzük |
| mercury_ore | Cıva Cevheri | mining/ore | 26 | T2 | nadir-altı | 9c | — |
| zinc_ore | Çinko Cevheri | mining/ore | 28 | T2 | nadir-altı | 10c | Çinko Külçesi |
| nickel_ore | Nikel Cevheri | mining/ore | 30 | T3 | nadir-altı | 11c | Nikel Külçesi |
| chromium_ore | Krom Cevheri | mining/ore | 33 | T3 | nadir-altı | 14c | Krom Külçesi |
| cobalt_ore | Kobalt Cevheri | mining/ore | 36 | T3 | nadir-altı | 16c | Kobalt Külçesi |
| silver_ore | Gümüş Cevheri | mining/ore | 39 | T3 | nadir-altı | 19c | Gümüş Külçesi |
| gold_ore | Altın Cevheri | mining/ore | 42 | T3 | nadir-altı | 23c | Altın Külçesi |
| electrum_ore | Elektrum Cevheri | mining/ore | 45 | T4 | nadir | 28c | Elektrum Külçesi |
| tungsten_ore | Tungsten Cevheri | mining/ore | 50 | T4 | nadir | 37c | Tungsten Külçesi |
| titanium_ore | Titanyum Cevheri | mining/ore | 55 | T4 | nadir | 49c | Titanyum Külçesi |
| palladium_ore | Paladyum Cevheri | mining/ore | 60 | T5 | destansı | 66c | Paladyum Külçesi |
| rhodium_ore | Rodyum Cevheri | mining/ore | 65 | T5 | destansı | 88c | Rodyum Külçesi |
| platinum_ore | Platin Cevheri | mining/ore | 70 | T5 | destansı | 1s 18c | Platin Külçesi |
| osmium_ore | Osmiyum Cevheri | mining/ore | 75 | T6 | efsanevi | 1s 58c | Osmiyum Külçesi |
| iridium_ore | İridyum Cevheri | mining/ore | 80 | T6 | efsanevi | 2s 12c | İridyum Külçesi |
| starsteel_ore | Yıldız Çeliği Cevheri | mining/ore | 85 | T6 | efsanevi | 2s 83c | Yıldız Çeliği Külçesi |
| mythril_ore | Mitril Cevheri | mining/ore | 90 | T7 | mitik | 3s 79c | Mitril Külçesi |
| orichalcum_ore | Orikalkum Cevheri | mining/ore | 95 | T7 | mitik | 5s 7c | Orikalkum Külçesi |
| adamantite_ore | Adamantit Cevheri | mining/ore | 100 | T7 | mitik | 6s 79c | Adamantit Külçesi |
| young_sapling | Fidan Dalı | woodcutting/log | 1 | T0 | sıradan | 2c | Acemi Arbalet, Acemi Ayin Baltası, Acemi Balta, Acemi Değnek, Acemi Doğa Odağı, Acemi Kısa Kılıç, Acemi Meç, Acemi Tonfa, Acemi Uzun Kılıç, Akçaağaç Lavta, Av Yayı, Ay Orağı, Bez Sargılar, Bileşik Yay, Bronz Orak, Büyülü Lavta, Demir Hançer, Demir Kılıç, Demir Muşta, Demir Orak, Demir Savaş Çekici, Demir Topuz, Dikenli Orak, Ejder Yumrukları, Gizemli Asa, Gök Gürültüsü Putu, Gölge Hançeri, Gümüş Savaş Çekici, Gümüş Telli Lavta, Gümüş Topuz, Izgara Lepistes, Işıltılı Çekiç, Kalfa Arbalet, Kalfa Ayin Baltası, Kalfa Balta, Kalfa Değnek, Kalfa Doğa Odağı, Kalfa Kısa Kılıç, Kalfa Meç, Kalfa Tonfa, Kalfa Uzun Kılıç, Kaliteli Arbalet, Kaliteli Ayin Baltası, Kaliteli Balta, Kaliteli Değnek, Kaliteli Doğa Odağı, Kaliteli Kısa Kılıç, Kaliteli Meç, Kaliteli Tonfa, Kaliteli Uzun Kılıç, Kaplan Pençeleri, Kemik Totem, Konser Lavtası, Koru Orağı, Kristal Asa, Kutsal Topuz, Kutsanmış Topuz, Kutsanmış Çekiç, Kıdemli Arbalet, Kıdemli Ayin Baltası, Kıdemli Balta, Kıdemli Değnek, Kıdemli Doğa Odağı, Kıdemli Kısa Kılıç, Kıdemli Meç, Kıdemli Tonfa, Kıdemli Uzun Kılıç, Kısa Yay, Meşe Tokmak, Oyma Totem, Paslı Hançer, Porsuk Yayı, Ruh Totemi, Sade Lavta, Su Verilmiş Kılıç, Söğüt Asa, Tahta Topuz, Talim Kılıcı, Taş Put, Tırtıklı Hançer, Usta Asası, Uzun Yay, Çelik Cestus, Çelik Hançer, Çelik Uzun Kılıç, Çırak Asası, Şampiyon Arbalet, Şampiyon Ayin Baltası, Şampiyon Balta, Şampiyon Değnek, Şampiyon Doğa Odağı, Şampiyon Kısa Kılıç, Şampiyon Meç, Şampiyon Tonfa, Şampiyon Uzun Kılıç, Şövalye Kılıcı |
| poplar_log | Kavak Kütüğü | woodcutting/log | 4 | T0 | sıradan | 3c | Adamant Arbalet, Adamant Asa, Adamant Ayin Baltası, Adamant Balta, Adamant Değnek, Adamant Doğa Odağı, Adamant Hançer, Adamant Kılıç, Adamant Kısa Kılıç, Adamant Lavta, Adamant Meç, Adamant Muşta, Adamant Orak, Adamant Savaş Çekici, Adamant Tonfa, Adamant Topuz, Adamant Totem, Adamant Uzun Kılıç, Adamant Yay, Ata Totemi, Efsaneler Lavtası, Ejder Kemiği Yayı, Göksel Yumruklar, Kadim Orak, Kutsal İntikamcı, Mithril Arbalet, Mithril Asa, Mithril Ayin Baltası, Mithril Balta, Mithril Değnek, Mithril Doğa Odağı, Mithril Hançer, Mithril Kılıç, Mithril Kısa Kılıç, Mithril Lavta, Mithril Meç, Mithril Muşta, Mithril Orak, Mithril Savaş Çekici, Mithril Tonfa, Mithril Topuz, Mithril Totem, Mithril Uzun Kılıç, Mithril Yay, Rünlü Arbalet, Rünlü Ayin Baltası, Rünlü Balta, Rünlü Büyük Kılıç, Rünlü Değnek, Rünlü Doğa Odağı, Rünlü Kısa Kılıç, Rünlü Meç, Rünlü Tonfa, Rünlü Uzun Kılıç, Sardalya Şiş, Yıldız Asası, Zehir Dişi, İlahi Topuz |
| willow_log | Söğüt Kütüğü | woodcutting/log | 7 | T0 | sıradan | 3c | Ejder Dövme Arbalet, Ejder Dövme Asa, Ejder Dövme Ayin Baltası, Ejder Dövme Balta, Ejder Dövme Değnek, Ejder Dövme Doğa Odağı, Ejder Dövme Hançer, Ejder Dövme Kılıç, Ejder Dövme Kısa Kılıç, Ejder Dövme Lavta, Ejder Dövme Meç, Ejder Dövme Muşta, Ejder Dövme Orak, Ejder Dövme Savaş Çekici, Ejder Dövme Tonfa, Ejder Dövme Topuz, Ejder Dövme Totem, Ejder Dövme Uzun Kılıç, Ejder Dövme Yay, Naneli Balık Çorbası, Yıldız Dövme Arbalet, Yıldız Dövme Asa, Yıldız Dövme Ayin Baltası, Yıldız Dövme Balta, Yıldız Dövme Değnek, Yıldız Dövme Doğa Odağı, Yıldız Dövme Hançer, Yıldız Dövme Kılıç, Yıldız Dövme Kısa Kılıç, Yıldız Dövme Lavta, Yıldız Dövme Meç, Yıldız Dövme Muşta, Yıldız Dövme Orak, Yıldız Dövme Savaş Çekici, Yıldız Dövme Tonfa, Yıldız Dövme Topuz, Yıldız Dövme Totem, Yıldız Dövme Uzun Kılıç, Yıldız Dövme Yay |
| fir_log | Köknar Kütüğü | woodcutting/log | 10 | T1 | sıradan | 4c | Titan Arbalet, Titan Asa, Titan Ayin Baltası, Titan Balta, Titan Değnek, Titan Doğa Odağı, Titan Hançer, Titan Kılıç, Titan Kısa Kılıç, Titan Lavta, Titan Meç, Titan Muşta, Titan Orak, Titan Savaş Çekici, Titan Tonfa, Titan Topuz, Titan Totem, Titan Uzun Kılıç, Titan Yay, Turna Fileto |
| pine_log | Çam Kütüğü | woodcutting/log | 12 | T1 | sıradan | 4c | Kraliyet Arbalet, Kraliyet Asa, Kraliyet Ayin Baltası, Kraliyet Balta, Kraliyet Değnek, Kraliyet Doğa Odağı, Kraliyet Hançer, Kraliyet Kılıç, Kraliyet Kısa Kılıç, Kraliyet Lavta, Kraliyet Meç, Kraliyet Muşta, Kraliyet Orak, Kraliyet Savaş Çekici, Kraliyet Tonfa, Kraliyet Topuz, Kraliyet Totem, Kraliyet Uzun Kılıç, Kraliyet Yay, Lavantalı Sazan |
| birch_log | Huş Kütüğü | woodcutting/log | 14 | T1 | sıradan | 5c | Güneş Dövme Arbalet, Güneş Dövme Asa, Güneş Dövme Ayin Baltası, Güneş Dövme Balta, Güneş Dövme Değnek, Güneş Dövme Doğa Odağı, Güneş Dövme Hançer, Güneş Dövme Kılıç, Güneş Dövme Kısa Kılıç, Güneş Dövme Lavta, Güneş Dövme Meç, Güneş Dövme Muşta, Güneş Dövme Orak, Güneş Dövme Savaş Çekici, Güneş Dövme Tonfa, Güneş Dövme Topuz, Güneş Dövme Totem, Güneş Dövme Uzun Kılıç, Güneş Dövme Yay, Izgara Gümüş Ringa |
| cedar_log | Sedir Kütüğü | woodcutting/log | 16 | T1 | sıradan | 5c | Obsidyen Arbalet, Obsidyen Asa, Obsidyen Ayin Baltası, Obsidyen Balta, Obsidyen Değnek, Obsidyen Doğa Odağı, Obsidyen Hançer, Obsidyen Kılıç, Obsidyen Kısa Kılıç, Obsidyen Lavta, Obsidyen Meç, Obsidyen Muşta, Obsidyen Orak, Obsidyen Savaş Çekici, Obsidyen Tonfa, Obsidyen Topuz, Obsidyen Totem, Obsidyen Uzun Kılıç, Obsidyen Yay, Otlu Alabalık |
| oak_log | Meşe Kütüğü | woodcutting/log | 18 | T1 | sıradan | 6c | Füme Somon |
| beech_log | Kayın Kütüğü | woodcutting/log | 20 | T2 | nadir-altı | 6c | Kum Çeliği Arbalet, Kum Çeliği Asa, Kum Çeliği Ayin Baltası, Kum Çeliği Balta, Kum Çeliği Değnek, Kum Çeliği Doğa Odağı, Kum Çeliği Hançer, Kum Çeliği Kılıç, Kum Çeliği Kısa Kılıç, Kum Çeliği Lavta, Kum Çeliği Meç, Kum Çeliği Muşta, Kum Çeliği Orak, Kum Çeliği Savaş Çekici, Kum Çeliği Tonfa, Kum Çeliği Topuz, Kum Çeliği Totem, Kum Çeliği Uzun Kılıç, Kum Çeliği Yay, Levrek Güveci |
| ash_log | Dişbudak Kütüğü | woodcutting/log | 22 | T2 | nadir-altı | 7c | Izgara Göl Levreği, Rün Bağlı Arbalet, Rün Bağlı Asa, Rün Bağlı Ayin Baltası, Rün Bağlı Balta, Rün Bağlı Değnek, Rün Bağlı Doğa Odağı, Rün Bağlı Hançer, Rün Bağlı Kılıç, Rün Bağlı Kısa Kılıç, Rün Bağlı Lavta, Rün Bağlı Meç, Rün Bağlı Muşta, Rün Bağlı Orak, Rün Bağlı Savaş Çekici, Rün Bağlı Tonfa, Rün Bağlı Topuz, Rün Bağlı Totem, Rün Bağlı Uzun Kılıç, Rün Bağlı Yay |
| maple_log | Akçaağaç Kütüğü | woodcutting/log | 24 | T2 | nadir-altı | 8c | Anka Arbalet, Anka Asa, Anka Ayin Baltası, Anka Balta, Anka Değnek, Anka Doğa Odağı, Anka Hançer, Anka Kılıç, Anka Kısa Kılıç, Anka Lavta, Anka Meç, Anka Muşta, Anka Orak, Anka Savaş Çekici, Anka Tonfa, Anka Topuz, Anka Totem, Anka Uzun Kılıç, Anka Yay, Izgara Kırmızı Mercan |
| walnut_log | Ceviz Kütüğü | woodcutting/log | 26 | T2 | nadir-altı | 9c | Dehşet Dövme Arbalet, Dehşet Dövme Asa, Dehşet Dövme Ayin Baltası, Dehşet Dövme Balta, Dehşet Dövme Değnek, Dehşet Dövme Doğa Odağı, Dehşet Dövme Hançer, Dehşet Dövme Kılıç, Dehşet Dövme Kısa Kılıç, Dehşet Dövme Lavta, Dehşet Dövme Meç, Dehşet Dövme Muşta, Dehşet Dövme Orak, Dehşet Dövme Savaş Çekici, Dehşet Dövme Tonfa, Dehşet Dövme Topuz, Dehşet Dövme Totem, Dehşet Dövme Uzun Kılıç, Dehşet Dövme Yay, Izgara Nehir Yayını |
| teak_log | Tik Kütüğü | woodcutting/log | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| mahogany_log | Maun Kütüğü | woodcutting/log | 30 | T3 | nadir-altı | 11c | Hükümdar Arbalet, Hükümdar Asa, Hükümdar Ayin Baltası, Hükümdar Balta, Hükümdar Değnek, Hükümdar Doğa Odağı, Hükümdar Hançer, Hükümdar Kılıç, Hükümdar Kısa Kılıç, Hükümdar Lavta, Hükümdar Meç, Hükümdar Muşta, Hükümdar Orak, Hükümdar Savaş Çekici, Hükümdar Tonfa, Hükümdar Topuz, Hükümdar Totem, Hükümdar Uzun Kılıç, Hükümdar Yay, Izgara Gökkuşağı Alabalığı, Mercan Arbalet, Mercan Asa, Mercan Ayin Baltası, Mercan Balta, Mercan Değnek, Mercan Doğa Odağı, Mercan Hançer, Mercan Kılıç, Mercan Kısa Kılıç, Mercan Lavta, Mercan Meç, Mercan Muşta, Mercan Orak, Mercan Savaş Çekici, Mercan Tonfa, Mercan Topuz, Mercan Totem, Mercan Uzun Kılıç, Mercan Yay |
| ebony_log | Abanoz Kütüğü | woodcutting/log | 33 | T3 | nadir-altı | 14c | Ejderha Kemiği Arbalet, Ejderha Kemiği Asa, Ejderha Kemiği Ayin Baltası, Ejderha Kemiği Balta, Ejderha Kemiği Değnek, Ejderha Kemiği Doğa Odağı, Ejderha Kemiği Hançer, Ejderha Kemiği Kılıç, Ejderha Kemiği Kısa Kılıç, Ejderha Kemiği Lavta, Ejderha Kemiği Meç, Ejderha Kemiği Muşta, Ejderha Kemiği Orak, Ejderha Kemiği Savaş Çekici, Ejderha Kemiği Tonfa, Ejderha Kemiği Topuz, Ejderha Kemiği Totem, Ejderha Kemiği Uzun Kılıç, Ejderha Kemiği Yay, Kraliyet Sazanı Ziyafeti |
| silver_willow_log | Gümüş Söğüt Kütüğü | woodcutting/log | 36 | T3 | nadir-altı | 16c | Izgara Zümrüt Turnası, Sis Örgü Arbalet, Sis Örgü Asa, Sis Örgü Ayin Baltası, Sis Örgü Balta, Sis Örgü Değnek, Sis Örgü Doğa Odağı, Sis Örgü Hançer, Sis Örgü Kılıç, Sis Örgü Kısa Kılıç, Sis Örgü Lavta, Sis Örgü Meç, Sis Örgü Muşta, Sis Örgü Orak, Sis Örgü Savaş Çekici, Sis Örgü Tonfa, Sis Örgü Topuz, Sis Örgü Totem, Sis Örgü Uzun Kılıç, Sis Örgü Yay |
| golden_oak_log | Altın Meşe Kütüğü | woodcutting/log | 39 | T3 | nadir-altı | 19c | Izgara Gümüş Barakuda, Prizmatik Arbalet, Prizmatik Asa, Prizmatik Ayin Baltası, Prizmatik Balta, Prizmatik Değnek, Prizmatik Doğa Odağı, Prizmatik Hançer, Prizmatik Kılıç, Prizmatik Kısa Kılıç, Prizmatik Lavta, Prizmatik Meç, Prizmatik Muşta, Prizmatik Orak, Prizmatik Savaş Çekici, Prizmatik Tonfa, Prizmatik Topuz, Prizmatik Totem, Prizmatik Uzun Kılıç, Prizmatik Yay |
| blackwood_log | Kara Odun Kütüğü | woodcutting/log | 42 | T3 | nadir-altı | 23c | Magma Dövme Arbalet, Magma Dövme Asa, Magma Dövme Ayin Baltası, Magma Dövme Balta, Magma Dövme Değnek, Magma Dövme Doğa Odağı, Magma Dövme Hançer, Magma Dövme Kılıç, Magma Dövme Kısa Kılıç, Magma Dövme Lavta, Magma Dövme Meç, Magma Dövme Muşta, Magma Dövme Orak, Magma Dövme Savaş Çekici, Magma Dövme Tonfa, Magma Dövme Topuz, Magma Dövme Totem, Magma Dövme Uzun Kılıç, Magma Dövme Yay, Mantarlı Mersin Turtası |
| crystal_birch_log | Kristal Huş Kütüğü | woodcutting/log | 45 | T4 | nadir | 28c | Bora Arbalet, Bora Asa, Bora Ayin Baltası, Bora Balta, Bora Değnek, Bora Doğa Odağı, Bora Hançer, Bora Kılıç, Bora Kısa Kılıç, Bora Lavta, Bora Meç, Bora Muşta, Bora Orak, Bora Savaş Çekici, Bora Tonfa, Bora Topuz, Bora Totem, Bora Uzun Kılıç, Bora Yay, Fırtına Izgarası Kılıçbalığı, Hiçlik Arbalet, Hiçlik Asa, Hiçlik Ayin Baltası, Hiçlik Balta, Hiçlik Değnek, Hiçlik Doğa Odağı, Hiçlik Hançer, Hiçlik Kılıç, Hiçlik Kısa Kılıç, Hiçlik Lavta, Hiçlik Meç, Hiçlik Muşta, Hiçlik Orak, Hiçlik Savaş Çekici, Hiçlik Tonfa, Hiçlik Topuz, Hiçlik Totem, Hiçlik Uzun Kılıç, Hiçlik Yay |
| emberwood_log | Korağacı Kütüğü | woodcutting/log | 50 | T4 | nadir | 37c | Izgara Fırtına Pullu Yılanbalığı, Melek Arbalet, Melek Asa, Melek Ayin Baltası, Melek Balta, Melek Değnek, Melek Doğa Odağı, Melek Hançer, Melek Kılıç, Melek Kısa Kılıç, Melek Lavta, Melek Meç, Melek Muşta, Melek Orak, Melek Savaş Çekici, Melek Tonfa, Melek Topuz, Melek Totem, Melek Uzun Kılıç, Melek Yay, Tutulma Arbalet, Tutulma Asa, Tutulma Ayin Baltası, Tutulma Balta, Tutulma Değnek, Tutulma Doğa Odağı, Tutulma Hançer, Tutulma Kılıç, Tutulma Kısa Kılıç, Tutulma Lavta, Tutulma Meç, Tutulma Muşta, Tutulma Orak, Tutulma Savaş Çekici, Tutulma Tonfa, Tutulma Topuz, Tutulma Totem, Tutulma Uzun Kılıç, Tutulma Yay |
| frostpine_log | Ayaz Çamı Kütüğü | woodcutting/log | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası, Ejder Lordu Arbalet, Ejder Lordu Asa, Ejder Lordu Ayin Baltası, Ejder Lordu Balta, Ejder Lordu Değnek, Ejder Lordu Doğa Odağı, Ejder Lordu Hançer, Ejder Lordu Kılıç, Ejder Lordu Kısa Kılıç, Ejder Lordu Lavta, Ejder Lordu Meç, Ejder Lordu Muşta, Ejder Lordu Orak, Ejder Lordu Savaş Çekici, Ejder Lordu Tonfa, Ejder Lordu Topuz, Ejder Lordu Totem, Ejder Lordu Uzun Kılıç, Ejder Lordu Yay, Göksel Arbalet, Göksel Asa, Göksel Ayin Baltası, Göksel Balta, Göksel Değnek, Göksel Doğa Odağı, Göksel Hançer, Göksel Kılıç, Göksel Kısa Kılıç, Göksel Lavta, Göksel Meç, Göksel Muşta, Göksel Orak, Göksel Savaş Çekici, Göksel Tonfa, Göksel Topuz, Göksel Totem, Göksel Uzun Kılıç, Göksel Yay |
| bloodwood_log | Kanağacı Kütüğü | woodcutting/log | 60 | T5 | destansı | 66c | Ebedi Arbalet, Ebedi Asa, Ebedi Ayin Baltası, Ebedi Balta, Ebedi Değnek, Ebedi Doğa Odağı, Ebedi Hançer, Ebedi Kılıç, Ebedi Kısa Kılıç, Ebedi Lavta, Ebedi Meç, Ebedi Muşta, Ebedi Orak, Ebedi Savaş Çekici, Ebedi Tonfa, Ebedi Topuz, Ebedi Totem, Ebedi Uzun Kılıç, Ebedi Yay, Izgara Obsidyen Ton, Yıldızsız Arbalet, Yıldızsız Asa, Yıldızsız Ayin Baltası, Yıldızsız Balta, Yıldızsız Değnek, Yıldızsız Doğa Odağı, Yıldızsız Hançer, Yıldızsız Kılıç, Yıldızsız Kısa Kılıç, Yıldızsız Lavta, Yıldızsız Meç, Yıldızsız Muşta, Yıldızsız Orak, Yıldızsız Savaş Çekici, Yıldızsız Tonfa, Yıldızsız Topuz, Yıldızsız Totem, Yıldızsız Uzun Kılıç, Yıldızsız Yay |
| spiritwood_log | Ruhağacı Kütüğü | woodcutting/log | 65 | T5 | destansı | 88c | Gece Çeliği Arbalet, Gece Çeliği Asa, Gece Çeliği Ayin Baltası, Gece Çeliği Balta, Gece Çeliği Değnek, Gece Çeliği Doğa Odağı, Gece Çeliği Hançer, Gece Çeliği Kılıç, Gece Çeliği Kısa Kılıç, Gece Çeliği Lavta, Gece Çeliği Meç, Gece Çeliği Muşta, Gece Çeliği Orak, Gece Çeliği Savaş Çekici, Gece Çeliği Tonfa, Gece Çeliği Topuz, Gece Çeliği Totem, Gece Çeliği Uzun Kılıç, Gece Çeliği Yay, Ruh Marlini Tabağı, Zaman Dövme Arbalet, Zaman Dövme Asa, Zaman Dövme Ayin Baltası, Zaman Dövme Balta, Zaman Dövme Değnek, Zaman Dövme Doğa Odağı, Zaman Dövme Hançer, Zaman Dövme Kılıç, Zaman Dövme Kısa Kılıç, Zaman Dövme Lavta, Zaman Dövme Meç, Zaman Dövme Muşta, Zaman Dövme Orak, Zaman Dövme Savaş Çekici, Zaman Dövme Tonfa, Zaman Dövme Topuz, Zaman Dövme Totem, Zaman Dövme Uzun Kılıç, Zaman Dövme Yay |
| ironwood_log | Kadim Demirağaç Kütüğü | woodcutting/log | 70 | T5 | destansı | 1s 18c | Fırtına Doğumlu Arbalet, Fırtına Doğumlu Asa, Fırtına Doğumlu Ayin Baltası, Fırtına Doğumlu Balta, Fırtına Doğumlu Değnek, Fırtına Doğumlu Doğa Odağı, Fırtına Doğumlu Hançer, Fırtına Doğumlu Kılıç, Fırtına Doğumlu Kısa Kılıç, Fırtına Doğumlu Lavta, Fırtına Doğumlu Meç, Fırtına Doğumlu Muşta, Fırtına Doğumlu Orak, Fırtına Doğumlu Savaş Çekici, Fırtına Doğumlu Tonfa, Fırtına Doğumlu Topuz, Fırtına Doğumlu Totem, Fırtına Doğumlu Uzun Kılıç, Fırtına Doğumlu Yay, Izgara İmparatorluk Altın Koisi, Yükselmiş Arbalet, Yükselmiş Asa, Yükselmiş Ayin Baltası, Yükselmiş Balta, Yükselmiş Değnek, Yükselmiş Doğa Odağı, Yükselmiş Hançer, Yükselmiş Kılıç, Yükselmiş Kısa Kılıç, Yükselmiş Lavta, Yükselmiş Meç, Yükselmiş Muşta, Yükselmiş Orak, Yükselmiş Savaş Çekici, Yükselmiş Tonfa, Yükselmiş Topuz, Yükselmiş Totem, Yükselmiş Uzun Kılıç, Yükselmiş Yay |
| starbark_log | Yıldızkabuğu Kütüğü | woodcutting/log | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı, Kaos Arbalet, Kaos Asa, Kaos Ayin Baltası, Kaos Balta, Kaos Değnek, Kaos Doğa Odağı, Kaos Hançer, Kaos Kılıç, Kaos Kısa Kılıç, Kaos Lavta, Kaos Meç, Kaos Muşta, Kaos Orak, Kaos Savaş Çekici, Kaos Tonfa, Kaos Topuz, Kaos Totem, Kaos Uzun Kılıç, Kaos Yay, Yarık Dövme Arbalet, Yarık Dövme Asa, Yarık Dövme Ayin Baltası, Yarık Dövme Balta, Yarık Dövme Değnek, Yarık Dövme Doğa Odağı, Yarık Dövme Hançer, Yarık Dövme Kılıç, Yarık Dövme Kısa Kılıç, Yarık Dövme Lavta, Yarık Dövme Meç, Yarık Dövme Muşta, Yarık Dövme Orak, Yarık Dövme Savaş Çekici, Yarık Dövme Tonfa, Yarık Dövme Topuz, Yarık Dövme Totem, Yarık Dövme Uzun Kılıç, Yarık Dövme Yay |
| moonwood_log | Ayağacı Kütüğü | woodcutting/log | 80 | T6 | efsanevi | 2s 12c | Başmelek Arbalet, Başmelek Asa, Başmelek Ayin Baltası, Başmelek Balta, Başmelek Değnek, Başmelek Doğa Odağı, Başmelek Hançer, Başmelek Kılıç, Başmelek Kısa Kılıç, Başmelek Lavta, Başmelek Meç, Başmelek Muşta, Başmelek Orak, Başmelek Savaş Çekici, Başmelek Tonfa, Başmelek Topuz, Başmelek Totem, Başmelek Uzun Kılıç, Başmelek Yay, Gök Dövme Arbalet, Gök Dövme Asa, Gök Dövme Ayin Baltası, Gök Dövme Balta, Gök Dövme Değnek, Gök Dövme Doğa Odağı, Gök Dövme Hançer, Gök Dövme Kılıç, Gök Dövme Kısa Kılıç, Gök Dövme Lavta, Gök Dövme Meç, Gök Dövme Muşta, Gök Dövme Orak, Gök Dövme Savaş Çekici, Gök Dövme Tonfa, Gök Dövme Topuz, Gök Dövme Totem, Gök Dövme Uzun Kılıç, Gök Dövme Yay, Izgara Boşluk Fenerbalığı |
| dragonheart_log | Ejderyüreği Kütüğü | woodcutting/log | 85 | T6 | efsanevi | 2s 83c | Kül Arbalet, Kül Asa, Kül Ayin Baltası, Kül Balta, Kül Değnek, Kül Doğa Odağı, Kül Hançer, Kül Kılıç, Kül Kısa Kılıç, Kül Lavta, Kül Meç, Kül Muşta, Kül Orak, Kül Savaş Çekici, Kül Tonfa, Kül Topuz, Kül Totem, Kül Uzun Kılıç, Kül Yay, Yıldız Aybalığı Lezzeti, Şafak Arbalet, Şafak Asa, Şafak Ayin Baltası, Şafak Balta, Şafak Değnek, Şafak Doğa Odağı, Şafak Hançer, Şafak Kılıç, Şafak Kısa Kılıç, Şafak Lavta, Şafak Meç, Şafak Muşta, Şafak Orak, Şafak Savaş Çekici, Şafak Tonfa, Şafak Topuz, Şafak Totem, Şafak Uzun Kılıç, Şafak Yay |
| voidwood_log | Boşlukağacı Kütüğü | woodcutting/log | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti, Kıyamet Dövme Arbalet, Kıyamet Dövme Asa, Kıyamet Dövme Ayin Baltası, Kıyamet Dövme Balta, Kıyamet Dövme Değnek, Kıyamet Dövme Doğa Odağı, Kıyamet Dövme Hançer, Kıyamet Dövme Kılıç, Kıyamet Dövme Kısa Kılıç, Kıyamet Dövme Lavta, Kıyamet Dövme Meç, Kıyamet Dövme Muşta, Kıyamet Dövme Orak, Kıyamet Dövme Savaş Çekici, Kıyamet Dövme Tonfa, Kıyamet Dövme Topuz, Kıyamet Dövme Totem, Kıyamet Dövme Uzun Kılıç, Kıyamet Dövme Yay, Leviathan Arbalet, Leviathan Asa, Leviathan Ayin Baltası, Leviathan Balta, Leviathan Değnek, Leviathan Doğa Odağı, Leviathan Hançer, Leviathan Kılıç, Leviathan Kısa Kılıç, Leviathan Lavta, Leviathan Meç, Leviathan Muşta, Leviathan Orak, Leviathan Savaş Çekici, Leviathan Tonfa, Leviathan Topuz, Leviathan Totem, Leviathan Uzun Kılıç, Leviathan Yay |
| celestial_ash_log | Göksel Dişbudak Kütüğü | woodcutting/log | 95 | T7 | mitik | 5s 7c | Hiçlik Dövme Arbalet, Hiçlik Dövme Asa, Hiçlik Dövme Ayin Baltası, Hiçlik Dövme Balta, Hiçlik Dövme Değnek, Hiçlik Dövme Doğa Odağı, Hiçlik Dövme Hançer, Hiçlik Dövme Kılıç, Hiçlik Dövme Kısa Kılıç, Hiçlik Dövme Lavta, Hiçlik Dövme Meç, Hiçlik Dövme Muşta, Hiçlik Dövme Orak, Hiçlik Dövme Savaş Çekici, Hiçlik Dövme Tonfa, Hiçlik Dövme Topuz, Hiçlik Dövme Totem, Hiçlik Dövme Uzun Kılıç, Hiçlik Dövme Yay, Izgara Uçurum Ejderbalığı, Tanrı Katili Arbalet, Tanrı Katili Asa, Tanrı Katili Ayin Baltası, Tanrı Katili Balta, Tanrı Katili Değnek, Tanrı Katili Doğa Odağı, Tanrı Katili Hançer, Tanrı Katili Kılıç, Tanrı Katili Kısa Kılıç, Tanrı Katili Lavta, Tanrı Katili Meç, Tanrı Katili Muşta, Tanrı Katili Orak, Tanrı Katili Savaş Çekici, Tanrı Katili Tonfa, Tanrı Katili Topuz, Tanrı Katili Totem, Tanrı Katili Uzun Kılıç, Tanrı Katili Yay |
| worldtree_log | Dünya Ağacı Kütüğü | woodcutting/log | 100 | T7 | mitik | 6s 79c | Kadim Çağ Arbalet, Kadim Çağ Asa, Kadim Çağ Ayin Baltası, Kadim Çağ Balta, Kadim Çağ Değnek, Kadim Çağ Doğa Odağı, Kadim Çağ Hançer, Kadim Çağ Kılıç, Kadim Çağ Kısa Kılıç, Kadim Çağ Lavta, Kadim Çağ Meç, Kadim Çağ Muşta, Kadim Çağ Orak, Kadim Çağ Savaş Çekici, Kadim Çağ Tonfa, Kadim Çağ Topuz, Kadim Çağ Totem, Kadim Çağ Uzun Kılıç, Kadim Çağ Yay, Leviathan Kraliyet Şöleni |

Kullanılmayan toplama kaynakları (hiçbir tarifte girdi değil — yalnız satılır):

| Alt tür | Kullanılmayan / toplam |
|---|---|
| feather | 0 / 4 |
| fish | 0 / 30 |
| fuel | 0 / 1 |
| herb | 0 / 30 |
| hide | 2 / 12 |
| log | 0 / 30 |
| ore | 1 / 25 |
| pelt | 0 / 12 |
| scale | 0 / 2 |
| stone | 1 / 4 |

### 2.2 Diğer eşyalar

| Kod | Ad | Sınıf | Sv | Değer | Etki | Kaynak | Kullanan tarif |
|---|---|---|---:|---:|---|---|---|
| food_grilled_guppy | Izgara Lepistes | consumable/food | 1 | 6c | +15 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_sardine_skewer | Sardalya Şiş | consumable/food | 40 | 9c | +16 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_minnow_soup | Naneli Balık Çorbası | consumable/food | 70 | 9c | +17 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_pike_fillet | Turna Fileto | consumable/food | 100 | 12c | +18 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_lavender_carp | Lavantalı Sazan | consumable/food | 120 | 12c | +18 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_silver_herring | Izgara Gümüş Ringa | consumable/food | 140 | 15c | +18 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_herb_trout | Otlu Alabalık | consumable/food | 160 | 15c | +19 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_smoked_salmon | Füme Somon | consumable/food | 180 | 18c | +20 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_bass_stew | Levrek Güveci | consumable/food | 200 | 18c | +20 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_lake_perch | Izgara Göl Levreği | consumable/food | 220 | 21c | +20 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_red_snapper | Izgara Kırmızı Mercan | consumable/food | 240 | 24c | +21 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_river_catfish | Izgara Nehir Yayını | consumable/food | 260 | 27c | +22 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_tuna_steak | Ginsengli Ton Balığı | consumable/food | 280 | 30c | +22 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_rainbow_trout | Izgara Gökkuşağı Alabalığı | consumable/food | 300 | 33c | +22 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_royal_carp_feast | Kraliyet Sazanı Ziyafeti | consumable/food | 330 | 42c | +23 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_emerald_pike | Izgara Zümrüt Turnası | consumable/food | 360 | 48c | +24 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_silver_barracuda | Izgara Gümüş Barakuda | consumable/food | 390 | 57c | +25 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_sturgeon_pie | Mantarlı Mersin Turtası | consumable/food | 420 | 69c | +26 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_swordfish_grill | Fırtına Izgarası Kılıçbalığı | consumable/food | 450 | 84c | +26 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_stormscale_eel | Izgara Fırtına Pullu Yılanbalığı | consumable/food | 500 | 1s 11c | +28 can, +%6 sal/sav 60 dk | üretim (cooking) | — |
| food_frostfin_chowder | Ayaz Somon Çorbası | consumable/food | 550 | 1s 47c | +29 can, +%6 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_obsidian_tuna | Izgara Obsidyen Ton | consumable/food | 600 | 1s 98c | +30 can, +%6 sal/sav 60 dk | üretim (cooking) | — |
| food_marlin_platter | Ruh Marlini Tabağı | consumable/food | 650 | 2s 64c | +31 can, +%7 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_golden_carp | Izgara İmparatorluk Altın Koisi | consumable/food | 700 | 3s 54c | +32 can, +%7 sal/sav 60 dk | üretim (cooking) | — |
| food_dragonblood_serpentfish | Ejderkanlı Yılanbalığı | consumable/food | 750 | 4s 74c | +34 can, +%7 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_void_angler | Izgara Boşluk Fenerbalığı | consumable/food | 800 | 6s 36c | +35 can, +%8 sal/sav 60 dk | üretim (cooking) | — |
| food_moonfish_delicacy | Yıldız Aybalığı Lezzeti | consumable/food | 850 | 8s 49c | +36 can, +%8 sal/sav 60 dk | üretim (cooking) | — |
| food_phoenix_koi_feast | Anka Koi Ziyafeti | consumable/food | 900 | 11s 37c | +38 can, +%8 sal/sav 60 dk | üretim (cooking) | — |
| food_grilled_abyssal_dragonfish | Izgara Uçurum Ejderbalığı | consumable/food | 950 | 15s 21c | +39 can, +%9 sal/sav 60 dk | üretim (cooking) | — |
| food_leviathan_banquet | Leviathan Kraliyet Şöleni | consumable/food | 1000 | 20s 37c | +40 can, +%9 sal/sav 60 dk | üretim (cooking) | — |
| potion_mp_small | Küçük Mana İksiri | consumable/potion | 1 | 20c |  | **yok** | — |
| potion_hp_small | Küçük Can İksiri | consumable/potion | 1 | 20c | +20 can | dükkân, üretim (alchemy) | — |
| potion_hp_medium | Orta Boy Can İksiri | consumable/potion | 100 | 60c | +30 can | dükkân, üretim (alchemy) | — |
| potion_hp_large | Büyük Boy Can İksiri | consumable/potion | 250 | 1s 50c | +40 can | dükkân, üretim (alchemy) | — |
| potion_hp_xlarge | Ekstra Büyük Can İksiri | consumable/potion | 450 | 3s 20c | +55 can | üretim (alchemy) | — |
| potion_hp_huge | Çok Büyük Can İksiri | consumable/potion | 650 | 6s | +70 can | üretim (alchemy) | — |
| potion_hp_supreme | Yüce Can İksiri | consumable/potion | 850 | 10s | +85 can | üretim (alchemy) | — |
| adamantite_bar | Adamantit Külçesi | material/bar | 1 | 54s 58c |  | üretim (smelting) | Kadim Çağ Arbalet, Kadim Çağ Asa, Kadim Çağ Ayin Baltası, Kadim Çağ Bacaklık, Kadim Çağ Balta, Kadim Çağ Değnek, Kadim Çağ Doğa Odağı, Kadim Çağ Eldiven, Kadim Çağ Göğüslük, Kadim Çağ Hançer, Kadim Çağ Kalkan, Kadim Çağ Kask, Kadim Çağ Kolye, Kadim Çağ Kılıç, Kadim Çağ Kısa Kılıç, Kadim Çağ Lavta, Kadim Çağ Meç, Kadim Çağ Muşta, Kadim Çağ Orak, Kadim Çağ Savaş Çekici, Kadim Çağ Tonfa, Kadim Çağ Topuz, Kadim Çağ Totem, Kadim Çağ Tılsım, Kadim Çağ Uzun Kılıç, Kadim Çağ Yay, Kadim Çağ Yüzük, Kadim Çağ Çizme |
| orichalcum_bar | Orikalkum Külçesi | material/bar | 1 | 40s 82c |  | üretim (smelting) | Hiçlik Dövme Arbalet, Hiçlik Dövme Asa, Hiçlik Dövme Ayin Baltası, Hiçlik Dövme Bacaklık, Hiçlik Dövme Balta, Hiçlik Dövme Değnek, Hiçlik Dövme Doğa Odağı, Hiçlik Dövme Eldiven, Hiçlik Dövme Göğüslük, Hiçlik Dövme Hançer, Hiçlik Dövme Kalkan, Hiçlik Dövme Kask, Hiçlik Dövme Kolye, Hiçlik Dövme Kılıç, Hiçlik Dövme Kısa Kılıç, Hiçlik Dövme Lavta, Hiçlik Dövme Meç, Hiçlik Dövme Muşta, Hiçlik Dövme Orak, Hiçlik Dövme Savaş Çekici, Hiçlik Dövme Tonfa, Hiçlik Dövme Topuz, Hiçlik Dövme Totem, Hiçlik Dövme Tılsım, Hiçlik Dövme Uzun Kılıç, Hiçlik Dövme Yay, Hiçlik Dövme Yüzük, Hiçlik Dövme Çizme, Kadim Çağ Ayin Baltası, Kadim Çağ Bacaklık, Kadim Çağ Balta, Kadim Çağ Eldiven, Kadim Çağ Göğüslük, Kadim Çağ Hançer, Kadim Çağ Kalkan, Kadim Çağ Kask, Kadim Çağ Kolye, Kadim Çağ Kılıç, Kadim Çağ Kısa Kılıç, Kadim Çağ Meç, Kadim Çağ Muşta, Kadim Çağ Orak, Kadim Çağ Savaş Çekici, Kadim Çağ Topuz, Kadim Çağ Tılsım, Kadim Çağ Uzun Kılıç, Kadim Çağ Yüzük, Kadim Çağ Çizme, Tanrı Katili Arbalet, Tanrı Katili Asa, Tanrı Katili Ayin Baltası, Tanrı Katili Bacaklık, Tanrı Katili Balta, Tanrı Katili Değnek, Tanrı Katili Doğa Odağı, Tanrı Katili Eldiven, Tanrı Katili Göğüslük, Tanrı Katili Hançer, Tanrı Katili Kalkan, Tanrı Katili Kask, Tanrı Katili Kolye, Tanrı Katili Kılıç, Tanrı Katili Kısa Kılıç, Tanrı Katili Lavta, Tanrı Katili Meç, Tanrı Katili Muşta, Tanrı Katili Orak, Tanrı Katili Savaş Çekici, Tanrı Katili Tonfa, Tanrı Katili Topuz, Tanrı Katili Totem, Tanrı Katili Tılsım, Tanrı Katili Uzun Kılıç, Tanrı Katili Yay, Tanrı Katili Yüzük, Tanrı Katili Çizme |
| mythril_bar | Mitril Külçesi | material/bar | 1 | 30s 58c |  | üretim (smelting) | Hiçlik Dövme Ayin Baltası, Hiçlik Dövme Bacaklık, Hiçlik Dövme Balta, Hiçlik Dövme Eldiven, Hiçlik Dövme Göğüslük, Hiçlik Dövme Hançer, Hiçlik Dövme Kalkan, Hiçlik Dövme Kask, Hiçlik Dövme Kolye, Hiçlik Dövme Kılıç, Hiçlik Dövme Kısa Kılıç, Hiçlik Dövme Meç, Hiçlik Dövme Muşta, Hiçlik Dövme Orak, Hiçlik Dövme Savaş Çekici, Hiçlik Dövme Topuz, Hiçlik Dövme Tılsım, Hiçlik Dövme Uzun Kılıç, Hiçlik Dövme Yüzük, Hiçlik Dövme Çizme, Kıyamet Dövme Arbalet, Kıyamet Dövme Asa, Kıyamet Dövme Ayin Baltası, Kıyamet Dövme Bacaklık, Kıyamet Dövme Balta, Kıyamet Dövme Değnek, Kıyamet Dövme Doğa Odağı, Kıyamet Dövme Eldiven, Kıyamet Dövme Göğüslük, Kıyamet Dövme Hançer, Kıyamet Dövme Kalkan, Kıyamet Dövme Kask, Kıyamet Dövme Kolye, Kıyamet Dövme Kılıç, Kıyamet Dövme Kısa Kılıç, Kıyamet Dövme Lavta, Kıyamet Dövme Meç, Kıyamet Dövme Muşta, Kıyamet Dövme Orak, Kıyamet Dövme Savaş Çekici, Kıyamet Dövme Tonfa, Kıyamet Dövme Topuz, Kıyamet Dövme Totem, Kıyamet Dövme Tılsım, Kıyamet Dövme Uzun Kılıç, Kıyamet Dövme Yay, Kıyamet Dövme Yüzük, Kıyamet Dövme Çizme, Leviathan Arbalet, Leviathan Asa, Leviathan Ayin Baltası, Leviathan Bacaklık, Leviathan Balta, Leviathan Değnek, Leviathan Doğa Odağı, Leviathan Eldiven, Leviathan Göğüslük, Leviathan Hançer, Leviathan Kalkan, Leviathan Kask, Leviathan Kolye, Leviathan Kılıç, Leviathan Kısa Kılıç, Leviathan Lavta, Leviathan Meç, Leviathan Muşta, Leviathan Orak, Leviathan Savaş Çekici, Leviathan Tonfa, Leviathan Topuz, Leviathan Totem, Leviathan Tılsım, Leviathan Uzun Kılıç, Leviathan Yay, Leviathan Yüzük, Leviathan Çizme, Tanrı Katili Ayin Baltası, Tanrı Katili Bacaklık, Tanrı Katili Balta, Tanrı Katili Eldiven, Tanrı Katili Göğüslük, Tanrı Katili Hançer, Tanrı Katili Kalkan, Tanrı Katili Kask, Tanrı Katili Kolye, Tanrı Katili Kılıç, Tanrı Katili Kısa Kılıç, Tanrı Katili Meç, Tanrı Katili Muşta, Tanrı Katili Orak, Tanrı Katili Savaş Çekici, Tanrı Katili Topuz, Tanrı Katili Tılsım, Tanrı Katili Uzun Kılıç, Tanrı Katili Yüzük, Tanrı Katili Çizme |
| starsteel_bar | Yıldız Çeliği Külçesi | material/bar | 1 | 18s 30c |  | üretim (smelting) | Kül Arbalet, Kül Asa, Kül Ayin Baltası, Kül Bacaklık, Kül Balta, Kül Değnek, Kül Doğa Odağı, Kül Eldiven, Kül Göğüslük, Kül Hançer, Kül Kalkan, Kül Kask, Kül Kolye, Kül Kılıç, Kül Kısa Kılıç, Kül Lavta, Kül Meç, Kül Muşta, Kül Orak, Kül Savaş Çekici, Kül Tonfa, Kül Topuz, Kül Totem, Kül Tılsım, Kül Uzun Kılıç, Kül Yay, Kül Yüzük, Kül Çizme, Kıyamet Dövme Ayin Baltası, Kıyamet Dövme Bacaklık, Kıyamet Dövme Balta, Kıyamet Dövme Eldiven, Kıyamet Dövme Göğüslük, Kıyamet Dövme Hançer, Kıyamet Dövme Kalkan, Kıyamet Dövme Kask, Kıyamet Dövme Kolye, Kıyamet Dövme Kılıç, Kıyamet Dövme Kısa Kılıç, Kıyamet Dövme Meç, Kıyamet Dövme Muşta, Kıyamet Dövme Orak, Kıyamet Dövme Savaş Çekici, Kıyamet Dövme Topuz, Kıyamet Dövme Tılsım, Kıyamet Dövme Uzun Kılıç, Kıyamet Dövme Yüzük, Kıyamet Dövme Çizme, Leviathan Ayin Baltası, Leviathan Bacaklık, Leviathan Balta, Leviathan Eldiven, Leviathan Göğüslük, Leviathan Hançer, Leviathan Kalkan, Leviathan Kask, Leviathan Kolye, Leviathan Kılıç, Leviathan Kısa Kılıç, Leviathan Meç, Leviathan Muşta, Leviathan Orak, Leviathan Savaş Çekici, Leviathan Topuz, Leviathan Tılsım, Leviathan Uzun Kılıç, Leviathan Yüzük, Leviathan Çizme, Şafak Arbalet, Şafak Asa, Şafak Ayin Baltası, Şafak Bacaklık, Şafak Balta, Şafak Değnek, Şafak Doğa Odağı, Şafak Eldiven, Şafak Göğüslük, Şafak Hançer, Şafak Kalkan, Şafak Kask, Şafak Kolye, Şafak Kılıç, Şafak Kısa Kılıç, Şafak Lavta, Şafak Meç, Şafak Muşta, Şafak Orak, Şafak Savaş Çekici, Şafak Tonfa, Şafak Topuz, Şafak Totem, Şafak Tılsım, Şafak Uzun Kılıç, Şafak Yay, Şafak Yüzük, Şafak Çizme |
| iridium_bar | İridyum Külçesi | material/bar | 1 | 13s 76c |  | üretim (smelting) | Başmelek Arbalet, Başmelek Asa, Başmelek Ayin Baltası, Başmelek Bacaklık, Başmelek Balta, Başmelek Değnek, Başmelek Doğa Odağı, Başmelek Eldiven, Başmelek Göğüslük, Başmelek Hançer, Başmelek Kalkan, Başmelek Kask, Başmelek Kolye, Başmelek Kılıç, Başmelek Kısa Kılıç, Başmelek Lavta, Başmelek Meç, Başmelek Muşta, Başmelek Orak, Başmelek Savaş Çekici, Başmelek Tonfa, Başmelek Topuz, Başmelek Totem, Başmelek Tılsım, Başmelek Uzun Kılıç, Başmelek Yay, Başmelek Yüzük, Başmelek Çizme, Gök Dövme Arbalet, Gök Dövme Asa, Gök Dövme Ayin Baltası, Gök Dövme Bacaklık, Gök Dövme Balta, Gök Dövme Değnek, Gök Dövme Doğa Odağı, Gök Dövme Eldiven, Gök Dövme Göğüslük, Gök Dövme Hançer, Gök Dövme Kalkan, Gök Dövme Kask, Gök Dövme Kolye, Gök Dövme Kılıç, Gök Dövme Kısa Kılıç, Gök Dövme Lavta, Gök Dövme Meç, Gök Dövme Muşta, Gök Dövme Orak, Gök Dövme Savaş Çekici, Gök Dövme Tonfa, Gök Dövme Topuz, Gök Dövme Totem, Gök Dövme Tılsım, Gök Dövme Uzun Kılıç, Gök Dövme Yay, Gök Dövme Yüzük, Gök Dövme Çizme, Kül Ayin Baltası, Kül Bacaklık, Kül Balta, Kül Eldiven, Kül Göğüslük, Kül Hançer, Kül Kalkan, Kül Kask, Kül Kolye, Kül Kılıç, Kül Kısa Kılıç, Kül Meç, Kül Muşta, Kül Orak, Kül Savaş Çekici, Kül Topuz, Kül Tılsım, Kül Uzun Kılıç, Kül Yüzük, Kül Çizme, Şafak Ayin Baltası, Şafak Bacaklık, Şafak Balta, Şafak Eldiven, Şafak Göğüslük, Şafak Hançer, Şafak Kalkan, Şafak Kask, Şafak Kolye, Şafak Kılıç, Şafak Kısa Kılıç, Şafak Meç, Şafak Muşta, Şafak Orak, Şafak Savaş Çekici, Şafak Topuz, Şafak Tılsım, Şafak Uzun Kılıç, Şafak Yüzük, Şafak Çizme |
| osmium_bar | Osmiyum Külçesi | material/bar | 1 | 10s 30c |  | üretim (smelting) | Başmelek Ayin Baltası, Başmelek Bacaklık, Başmelek Balta, Başmelek Eldiven, Başmelek Göğüslük, Başmelek Hançer, Başmelek Kalkan, Başmelek Kask, Başmelek Kolye, Başmelek Kılıç, Başmelek Kısa Kılıç, Başmelek Meç, Başmelek Muşta, Başmelek Orak, Başmelek Savaş Çekici, Başmelek Topuz, Başmelek Tılsım, Başmelek Uzun Kılıç, Başmelek Yüzük, Başmelek Çizme, Gök Dövme Ayin Baltası, Gök Dövme Bacaklık, Gök Dövme Balta, Gök Dövme Eldiven, Gök Dövme Göğüslük, Gök Dövme Hançer, Gök Dövme Kalkan, Gök Dövme Kask, Gök Dövme Kolye, Gök Dövme Kılıç, Gök Dövme Kısa Kılıç, Gök Dövme Meç, Gök Dövme Muşta, Gök Dövme Orak, Gök Dövme Savaş Çekici, Gök Dövme Topuz, Gök Dövme Tılsım, Gök Dövme Uzun Kılıç, Gök Dövme Yüzük, Gök Dövme Çizme, Kaos Arbalet, Kaos Asa, Kaos Ayin Baltası, Kaos Bacaklık, Kaos Balta, Kaos Değnek, Kaos Doğa Odağı, Kaos Eldiven, Kaos Göğüslük, Kaos Hançer, Kaos Kalkan, Kaos Kask, Kaos Kolye, Kaos Kılıç, Kaos Kısa Kılıç, Kaos Lavta, Kaos Meç, Kaos Muşta, Kaos Orak, Kaos Savaş Çekici, Kaos Tonfa, Kaos Topuz, Kaos Totem, Kaos Tılsım, Kaos Uzun Kılıç, Kaos Yay, Kaos Yüzük, Kaos Çizme, Yarık Dövme Arbalet, Yarık Dövme Asa, Yarık Dövme Ayin Baltası, Yarık Dövme Bacaklık, Yarık Dövme Balta, Yarık Dövme Değnek, Yarık Dövme Doğa Odağı, Yarık Dövme Eldiven, Yarık Dövme Göğüslük, Yarık Dövme Hançer, Yarık Dövme Kalkan, Yarık Dövme Kask, Yarık Dövme Kolye, Yarık Dövme Kılıç, Yarık Dövme Kısa Kılıç, Yarık Dövme Lavta, Yarık Dövme Meç, Yarık Dövme Muşta, Yarık Dövme Orak, Yarık Dövme Savaş Çekici, Yarık Dövme Tonfa, Yarık Dövme Topuz, Yarık Dövme Totem, Yarık Dövme Tılsım, Yarık Dövme Uzun Kılıç, Yarık Dövme Yay, Yarık Dövme Yüzük, Yarık Dövme Çizme |
| platinum_bar | Platin Külçesi | material/bar | 1 | 5s 86c |  | üretim (smelting) | Fırtına Doğumlu Arbalet, Fırtına Doğumlu Asa, Fırtına Doğumlu Ayin Baltası, Fırtına Doğumlu Bacaklık, Fırtına Doğumlu Balta, Fırtına Doğumlu Değnek, Fırtına Doğumlu Doğa Odağı, Fırtına Doğumlu Eldiven, Fırtına Doğumlu Göğüslük, Fırtına Doğumlu Hançer, Fırtına Doğumlu Kalkan, Fırtına Doğumlu Kask, Fırtına Doğumlu Kolye, Fırtına Doğumlu Kılıç, Fırtına Doğumlu Kısa Kılıç, Fırtına Doğumlu Lavta, Fırtına Doğumlu Meç, Fırtına Doğumlu Muşta, Fırtına Doğumlu Orak, Fırtına Doğumlu Savaş Çekici, Fırtına Doğumlu Tonfa, Fırtına Doğumlu Topuz, Fırtına Doğumlu Totem, Fırtına Doğumlu Tılsım, Fırtına Doğumlu Uzun Kılıç, Fırtına Doğumlu Yay, Fırtına Doğumlu Yüzük, Fırtına Doğumlu Çizme, Kaos Ayin Baltası, Kaos Bacaklık, Kaos Balta, Kaos Eldiven, Kaos Göğüslük, Kaos Hançer, Kaos Kalkan, Kaos Kask, Kaos Kolye, Kaos Kılıç, Kaos Kısa Kılıç, Kaos Meç, Kaos Muşta, Kaos Orak, Kaos Savaş Çekici, Kaos Topuz, Kaos Tılsım, Kaos Uzun Kılıç, Kaos Yüzük, Kaos Çizme, Yarık Dövme Ayin Baltası, Yarık Dövme Bacaklık, Yarık Dövme Balta, Yarık Dövme Eldiven, Yarık Dövme Göğüslük, Yarık Dövme Hançer, Yarık Dövme Kalkan, Yarık Dövme Kask, Yarık Dövme Kolye, Yarık Dövme Kılıç, Yarık Dövme Kısa Kılıç, Yarık Dövme Meç, Yarık Dövme Muşta, Yarık Dövme Orak, Yarık Dövme Savaş Çekici, Yarık Dövme Topuz, Yarık Dövme Tılsım, Yarık Dövme Uzun Kılıç, Yarık Dövme Yüzük, Yarık Dövme Çizme, Yükselmiş Arbalet, Yükselmiş Asa, Yükselmiş Ayin Baltası, Yükselmiş Bacaklık, Yükselmiş Balta, Yükselmiş Değnek, Yükselmiş Doğa Odağı, Yükselmiş Eldiven, Yükselmiş Göğüslük, Yükselmiş Hançer, Yükselmiş Kalkan, Yükselmiş Kask, Yükselmiş Kolye, Yükselmiş Kılıç, Yükselmiş Kısa Kılıç, Yükselmiş Lavta, Yükselmiş Meç, Yükselmiş Muşta, Yükselmiş Orak, Yükselmiş Savaş Çekici, Yükselmiş Tonfa, Yükselmiş Topuz, Yükselmiş Totem, Yükselmiş Tılsım, Yükselmiş Uzun Kılıç, Yükselmiş Yay, Yükselmiş Yüzük, Yükselmiş Çizme |
| rhodium_bar | Rodyum Külçesi | material/bar | 1 | 4s 42c |  | üretim (smelting) | Fırtına Doğumlu Ayin Baltası, Fırtına Doğumlu Bacaklık, Fırtına Doğumlu Balta, Fırtına Doğumlu Eldiven, Fırtına Doğumlu Göğüslük, Fırtına Doğumlu Hançer, Fırtına Doğumlu Kalkan, Fırtına Doğumlu Kask, Fırtına Doğumlu Kolye, Fırtına Doğumlu Kılıç, Fırtına Doğumlu Kısa Kılıç, Fırtına Doğumlu Meç, Fırtına Doğumlu Muşta, Fırtına Doğumlu Orak, Fırtına Doğumlu Savaş Çekici, Fırtına Doğumlu Topuz, Fırtına Doğumlu Tılsım, Fırtına Doğumlu Uzun Kılıç, Fırtına Doğumlu Yüzük, Fırtına Doğumlu Çizme, Gece Çeliği Arbalet, Gece Çeliği Asa, Gece Çeliği Ayin Baltası, Gece Çeliği Bacaklık, Gece Çeliği Balta, Gece Çeliği Değnek, Gece Çeliği Doğa Odağı, Gece Çeliği Eldiven, Gece Çeliği Göğüslük, Gece Çeliği Hançer, Gece Çeliği Kalkan, Gece Çeliği Kask, Gece Çeliği Kolye, Gece Çeliği Kılıç, Gece Çeliği Kısa Kılıç, Gece Çeliği Lavta, Gece Çeliği Meç, Gece Çeliği Muşta, Gece Çeliği Orak, Gece Çeliği Savaş Çekici, Gece Çeliği Tonfa, Gece Çeliği Topuz, Gece Çeliği Totem, Gece Çeliği Tılsım, Gece Çeliği Uzun Kılıç, Gece Çeliği Yay, Gece Çeliği Yüzük, Gece Çeliği Çizme, Yükselmiş Ayin Baltası, Yükselmiş Bacaklık, Yükselmiş Balta, Yükselmiş Eldiven, Yükselmiş Göğüslük, Yükselmiş Hançer, Yükselmiş Kalkan, Yükselmiş Kask, Yükselmiş Kolye, Yükselmiş Kılıç, Yükselmiş Kısa Kılıç, Yükselmiş Meç, Yükselmiş Muşta, Yükselmiş Orak, Yükselmiş Savaş Çekici, Yükselmiş Topuz, Yükselmiş Tılsım, Yükselmiş Uzun Kılıç, Yükselmiş Yüzük, Yükselmiş Çizme, Zaman Dövme Arbalet, Zaman Dövme Asa, Zaman Dövme Ayin Baltası, Zaman Dövme Bacaklık, Zaman Dövme Balta, Zaman Dövme Değnek, Zaman Dövme Doğa Odağı, Zaman Dövme Eldiven, Zaman Dövme Göğüslük, Zaman Dövme Hançer, Zaman Dövme Kalkan, Zaman Dövme Kask, Zaman Dövme Kolye, Zaman Dövme Kılıç, Zaman Dövme Kısa Kılıç, Zaman Dövme Lavta, Zaman Dövme Meç, Zaman Dövme Muşta, Zaman Dövme Orak, Zaman Dövme Savaş Çekici, Zaman Dövme Tonfa, Zaman Dövme Topuz, Zaman Dövme Totem, Zaman Dövme Tılsım, Zaman Dövme Uzun Kılıç, Zaman Dövme Yay, Zaman Dövme Yüzük, Zaman Dövme Çizme |
| palladium_bar | Paladyum Külçesi | material/bar | 1 | 3s 36c |  | üretim (smelting) | Ebedi Arbalet, Ebedi Asa, Ebedi Ayin Baltası, Ebedi Bacaklık, Ebedi Balta, Ebedi Değnek, Ebedi Doğa Odağı, Ebedi Eldiven, Ebedi Göğüslük, Ebedi Hançer, Ebedi Kalkan, Ebedi Kask, Ebedi Kolye, Ebedi Kılıç, Ebedi Kısa Kılıç, Ebedi Lavta, Ebedi Meç, Ebedi Muşta, Ebedi Orak, Ebedi Savaş Çekici, Ebedi Tonfa, Ebedi Topuz, Ebedi Totem, Ebedi Tılsım, Ebedi Uzun Kılıç, Ebedi Yay, Ebedi Yüzük, Ebedi Çizme, Gece Çeliği Ayin Baltası, Gece Çeliği Bacaklık, Gece Çeliği Balta, Gece Çeliği Eldiven, Gece Çeliği Göğüslük, Gece Çeliği Hançer, Gece Çeliği Kalkan, Gece Çeliği Kask, Gece Çeliği Kolye, Gece Çeliği Kılıç, Gece Çeliği Kısa Kılıç, Gece Çeliği Meç, Gece Çeliği Muşta, Gece Çeliği Orak, Gece Çeliği Savaş Çekici, Gece Çeliği Topuz, Gece Çeliği Tılsım, Gece Çeliği Uzun Kılıç, Gece Çeliği Yüzük, Gece Çeliği Çizme, Yıldızsız Arbalet, Yıldızsız Asa, Yıldızsız Ayin Baltası, Yıldızsız Bacaklık, Yıldızsız Balta, Yıldızsız Değnek, Yıldızsız Doğa Odağı, Yıldızsız Eldiven, Yıldızsız Göğüslük, Yıldızsız Hançer, Yıldızsız Kalkan, Yıldızsız Kask, Yıldızsız Kolye, Yıldızsız Kılıç, Yıldızsız Kısa Kılıç, Yıldızsız Lavta, Yıldızsız Meç, Yıldızsız Muşta, Yıldızsız Orak, Yıldızsız Savaş Çekici, Yıldızsız Tonfa, Yıldızsız Topuz, Yıldızsız Totem, Yıldızsız Tılsım, Yıldızsız Uzun Kılıç, Yıldızsız Yay, Yıldızsız Yüzük, Yıldızsız Çizme, Zaman Dövme Ayin Baltası, Zaman Dövme Bacaklık, Zaman Dövme Balta, Zaman Dövme Eldiven, Zaman Dövme Göğüslük, Zaman Dövme Hançer, Zaman Dövme Kalkan, Zaman Dövme Kask, Zaman Dövme Kolye, Zaman Dövme Kılıç, Zaman Dövme Kısa Kılıç, Zaman Dövme Meç, Zaman Dövme Muşta, Zaman Dövme Orak, Zaman Dövme Savaş Çekici, Zaman Dövme Topuz, Zaman Dövme Tılsım, Zaman Dövme Uzun Kılıç, Zaman Dövme Yüzük, Zaman Dövme Çizme |
| titanium_bar | Titanyum Külçesi | material/bar | 1 | 2s 48c |  | üretim (smelting) | Ebedi Ayin Baltası, Ebedi Bacaklık, Ebedi Balta, Ebedi Eldiven, Ebedi Göğüslük, Ebedi Hançer, Ebedi Kalkan, Ebedi Kask, Ebedi Kolye, Ebedi Kılıç, Ebedi Kısa Kılıç, Ebedi Meç, Ebedi Muşta, Ebedi Orak, Ebedi Savaş Çekici, Ebedi Topuz, Ebedi Tılsım, Ebedi Uzun Kılıç, Ebedi Yüzük, Ebedi Çizme, Ejder Lordu Arbalet, Ejder Lordu Asa, Ejder Lordu Ayin Baltası, Ejder Lordu Bacaklık, Ejder Lordu Balta, Ejder Lordu Değnek, Ejder Lordu Doğa Odağı, Ejder Lordu Eldiven, Ejder Lordu Göğüslük, Ejder Lordu Hançer, Ejder Lordu Kalkan, Ejder Lordu Kask, Ejder Lordu Kolye, Ejder Lordu Kılıç, Ejder Lordu Kısa Kılıç, Ejder Lordu Lavta, Ejder Lordu Meç, Ejder Lordu Muşta, Ejder Lordu Orak, Ejder Lordu Savaş Çekici, Ejder Lordu Tonfa, Ejder Lordu Topuz, Ejder Lordu Totem, Ejder Lordu Tılsım, Ejder Lordu Uzun Kılıç, Ejder Lordu Yay, Ejder Lordu Yüzük, Ejder Lordu Çizme, Göksel Arbalet, Göksel Asa, Göksel Ayin Baltası, Göksel Bacaklık, Göksel Balta, Göksel Değnek, Göksel Doğa Odağı, Göksel Eldiven, Göksel Göğüslük, Göksel Hançer, Göksel Kalkan, Göksel Kask, Göksel Kolye, Göksel Kılıç, Göksel Kısa Kılıç, Göksel Lavta, Göksel Meç, Göksel Muşta, Göksel Orak, Göksel Savaş Çekici, Göksel Tonfa, Göksel Topuz, Göksel Totem, Göksel Tılsım, Göksel Uzun Kılıç, Göksel Yay, Göksel Yüzük, Göksel Çizme, Yıldızsız Ayin Baltası, Yıldızsız Bacaklık, Yıldızsız Balta, Yıldızsız Eldiven, Yıldızsız Göğüslük, Yıldızsız Hançer, Yıldızsız Kalkan, Yıldızsız Kask, Yıldızsız Kolye, Yıldızsız Kılıç, Yıldızsız Kısa Kılıç, Yıldızsız Meç, Yıldızsız Muşta, Yıldızsız Orak, Yıldızsız Savaş Çekici, Yıldızsız Topuz, Yıldızsız Tılsım, Yıldızsız Uzun Kılıç, Yıldızsız Yüzük, Yıldızsız Çizme |
| tungsten_bar | Tungsten Külçesi | material/bar | 1 | 1s 90c |  | üretim (smelting) | Ejder Lordu Ayin Baltası, Ejder Lordu Bacaklık, Ejder Lordu Balta, Ejder Lordu Eldiven, Ejder Lordu Göğüslük, Ejder Lordu Hançer, Ejder Lordu Kalkan, Ejder Lordu Kask, Ejder Lordu Kolye, Ejder Lordu Kılıç, Ejder Lordu Kısa Kılıç, Ejder Lordu Meç, Ejder Lordu Muşta, Ejder Lordu Orak, Ejder Lordu Savaş Çekici, Ejder Lordu Topuz, Ejder Lordu Tılsım, Ejder Lordu Uzun Kılıç, Ejder Lordu Yüzük, Ejder Lordu Çizme, Göksel Ayin Baltası, Göksel Bacaklık, Göksel Balta, Göksel Eldiven, Göksel Göğüslük, Göksel Hançer, Göksel Kalkan, Göksel Kask, Göksel Kolye, Göksel Kılıç, Göksel Kısa Kılıç, Göksel Meç, Göksel Muşta, Göksel Orak, Göksel Savaş Çekici, Göksel Topuz, Göksel Tılsım, Göksel Uzun Kılıç, Göksel Yüzük, Göksel Çizme, Melek Arbalet, Melek Asa, Melek Ayin Baltası, Melek Bacaklık, Melek Balta, Melek Değnek, Melek Doğa Odağı, Melek Eldiven, Melek Göğüslük, Melek Hançer, Melek Kalkan, Melek Kask, Melek Kolye, Melek Kılıç, Melek Kısa Kılıç, Melek Lavta, Melek Meç, Melek Muşta, Melek Orak, Melek Savaş Çekici, Melek Tonfa, Melek Topuz, Melek Totem, Melek Tılsım, Melek Uzun Kılıç, Melek Yay, Melek Yüzük, Melek Çizme, Tutulma Arbalet, Tutulma Asa, Tutulma Ayin Baltası, Tutulma Bacaklık, Tutulma Balta, Tutulma Değnek, Tutulma Doğa Odağı, Tutulma Eldiven, Tutulma Göğüslük, Tutulma Hançer, Tutulma Kalkan, Tutulma Kask, Tutulma Kolye, Tutulma Kılıç, Tutulma Kısa Kılıç, Tutulma Lavta, Tutulma Meç, Tutulma Muşta, Tutulma Orak, Tutulma Savaş Çekici, Tutulma Tonfa, Tutulma Topuz, Tutulma Totem, Tutulma Tılsım, Tutulma Uzun Kılıç, Tutulma Yay, Tutulma Yüzük, Tutulma Çizme |
| electrum_bar | Elektrum Külçesi | material/bar | 1 | 1s 47c |  | üretim (smelting) | Bora Arbalet, Bora Asa, Bora Ayin Baltası, Bora Bacaklık, Bora Balta, Bora Değnek, Bora Doğa Odağı, Bora Eldiven, Bora Göğüslük, Bora Hançer, Bora Kalkan, Bora Kask, Bora Kolye, Bora Kılıç, Bora Kısa Kılıç, Bora Lavta, Bora Meç, Bora Muşta, Bora Orak, Bora Savaş Çekici, Bora Tonfa, Bora Topuz, Bora Totem, Bora Tılsım, Bora Uzun Kılıç, Bora Yay, Bora Yüzük, Bora Çizme, Hiçlik Arbalet, Hiçlik Asa, Hiçlik Ayin Baltası, Hiçlik Bacaklık, Hiçlik Balta, Hiçlik Değnek, Hiçlik Doğa Odağı, Hiçlik Eldiven, Hiçlik Göğüslük, Hiçlik Hançer, Hiçlik Kalkan, Hiçlik Kask, Hiçlik Kolye, Hiçlik Kılıç, Hiçlik Kısa Kılıç, Hiçlik Lavta, Hiçlik Meç, Hiçlik Muşta, Hiçlik Orak, Hiçlik Savaş Çekici, Hiçlik Tonfa, Hiçlik Topuz, Hiçlik Totem, Hiçlik Tılsım, Hiçlik Uzun Kılıç, Hiçlik Yay, Hiçlik Yüzük, Hiçlik Çizme, Melek Ayin Baltası, Melek Bacaklık, Melek Balta, Melek Eldiven, Melek Göğüslük, Melek Hançer, Melek Kalkan, Melek Kask, Melek Kolye, Melek Kılıç, Melek Kısa Kılıç, Melek Meç, Melek Muşta, Melek Orak, Melek Savaş Çekici, Melek Topuz, Melek Tılsım, Melek Uzun Kılıç, Melek Yüzük, Melek Çizme, Tutulma Ayin Baltası, Tutulma Bacaklık, Tutulma Balta, Tutulma Eldiven, Tutulma Göğüslük, Tutulma Hançer, Tutulma Kalkan, Tutulma Kask, Tutulma Kolye, Tutulma Kılıç, Tutulma Kısa Kılıç, Tutulma Meç, Tutulma Muşta, Tutulma Orak, Tutulma Savaş Çekici, Tutulma Topuz, Tutulma Tılsım, Tutulma Uzun Kılıç, Tutulma Yüzük, Tutulma Çizme |
| gold_bar | Altın Külçesi | material/bar | 1 | 86c |  | üretim (smelting) | Bora Ayin Baltası, Bora Bacaklık, Bora Balta, Bora Eldiven, Bora Göğüslük, Bora Hançer, Bora Kalkan, Bora Kask, Bora Kolye, Bora Kılıç, Bora Kısa Kılıç, Bora Meç, Bora Muşta, Bora Orak, Bora Savaş Çekici, Bora Topuz, Bora Tılsım, Bora Uzun Kılıç, Bora Yüzük, Bora Çizme, Hiçlik Ayin Baltası, Hiçlik Bacaklık, Hiçlik Balta, Hiçlik Eldiven, Hiçlik Göğüslük, Hiçlik Hançer, Hiçlik Kalkan, Hiçlik Kask, Hiçlik Kolye, Hiçlik Kılıç, Hiçlik Kısa Kılıç, Hiçlik Meç, Hiçlik Muşta, Hiçlik Orak, Hiçlik Savaş Çekici, Hiçlik Topuz, Hiçlik Tılsım, Hiçlik Uzun Kılıç, Hiçlik Yüzük, Hiçlik Çizme, Magma Dövme Arbalet, Magma Dövme Asa, Magma Dövme Ayin Baltası, Magma Dövme Bacaklık, Magma Dövme Balta, Magma Dövme Değnek, Magma Dövme Doğa Odağı, Magma Dövme Eldiven, Magma Dövme Göğüslük, Magma Dövme Hançer, Magma Dövme Kalkan, Magma Dövme Kask, Magma Dövme Kolye, Magma Dövme Kılıç, Magma Dövme Kısa Kılıç, Magma Dövme Lavta, Magma Dövme Meç, Magma Dövme Muşta, Magma Dövme Orak, Magma Dövme Savaş Çekici, Magma Dövme Tonfa, Magma Dövme Topuz, Magma Dövme Totem, Magma Dövme Tılsım, Magma Dövme Uzun Kılıç, Magma Dövme Yay, Magma Dövme Yüzük, Magma Dövme Çizme |
| silver_bar | Gümüş Külçesi | material/bar | 1 | 74c |  | üretim (smelting) | Magma Dövme Ayin Baltası, Magma Dövme Bacaklık, Magma Dövme Balta, Magma Dövme Eldiven, Magma Dövme Göğüslük, Magma Dövme Hançer, Magma Dövme Kalkan, Magma Dövme Kask, Magma Dövme Kolye, Magma Dövme Kılıç, Magma Dövme Kısa Kılıç, Magma Dövme Meç, Magma Dövme Muşta, Magma Dövme Orak, Magma Dövme Savaş Çekici, Magma Dövme Topuz, Magma Dövme Tılsım, Magma Dövme Uzun Kılıç, Magma Dövme Yüzük, Magma Dövme Çizme, Prizmatik Arbalet, Prizmatik Asa, Prizmatik Ayin Baltası, Prizmatik Bacaklık, Prizmatik Balta, Prizmatik Değnek, Prizmatik Doğa Odağı, Prizmatik Eldiven, Prizmatik Göğüslük, Prizmatik Hançer, Prizmatik Kalkan, Prizmatik Kask, Prizmatik Kolye, Prizmatik Kılıç, Prizmatik Kısa Kılıç, Prizmatik Lavta, Prizmatik Meç, Prizmatik Muşta, Prizmatik Orak, Prizmatik Savaş Çekici, Prizmatik Tonfa, Prizmatik Topuz, Prizmatik Totem, Prizmatik Tılsım, Prizmatik Uzun Kılıç, Prizmatik Yay, Prizmatik Yüzük, Prizmatik Çizme |
| cobalt_bar | Kobalt Külçesi | material/bar | 1 | 64c |  | üretim (smelting) | Prizmatik Ayin Baltası, Prizmatik Bacaklık, Prizmatik Balta, Prizmatik Eldiven, Prizmatik Göğüslük, Prizmatik Hançer, Prizmatik Kalkan, Prizmatik Kask, Prizmatik Kolye, Prizmatik Kılıç, Prizmatik Kısa Kılıç, Prizmatik Meç, Prizmatik Muşta, Prizmatik Orak, Prizmatik Savaş Çekici, Prizmatik Topuz, Prizmatik Tılsım, Prizmatik Uzun Kılıç, Prizmatik Yüzük, Prizmatik Çizme, Sis Örgü Arbalet, Sis Örgü Asa, Sis Örgü Ayin Baltası, Sis Örgü Bacaklık, Sis Örgü Balta, Sis Örgü Değnek, Sis Örgü Doğa Odağı, Sis Örgü Eldiven, Sis Örgü Göğüslük, Sis Örgü Hançer, Sis Örgü Kalkan, Sis Örgü Kask, Sis Örgü Kolye, Sis Örgü Kılıç, Sis Örgü Kısa Kılıç, Sis Örgü Lavta, Sis Örgü Meç, Sis Örgü Muşta, Sis Örgü Orak, Sis Örgü Savaş Çekici, Sis Örgü Tonfa, Sis Örgü Topuz, Sis Örgü Totem, Sis Örgü Tılsım, Sis Örgü Uzun Kılıç, Sis Örgü Yay, Sis Örgü Yüzük, Sis Örgü Çizme |
| chromium_bar | Krom Külçesi | material/bar | 1 | 58c |  | üretim (smelting) | Ejderha Kemiği Arbalet, Ejderha Kemiği Asa, Ejderha Kemiği Ayin Baltası, Ejderha Kemiği Bacaklık, Ejderha Kemiği Balta, Ejderha Kemiği Değnek, Ejderha Kemiği Doğa Odağı, Ejderha Kemiği Eldiven, Ejderha Kemiği Göğüslük, Ejderha Kemiği Hançer, Ejderha Kemiği Kalkan, Ejderha Kemiği Kask, Ejderha Kemiği Kolye, Ejderha Kemiği Kılıç, Ejderha Kemiği Kısa Kılıç, Ejderha Kemiği Lavta, Ejderha Kemiği Meç, Ejderha Kemiği Muşta, Ejderha Kemiği Orak, Ejderha Kemiği Savaş Çekici, Ejderha Kemiği Tonfa, Ejderha Kemiği Topuz, Ejderha Kemiği Totem, Ejderha Kemiği Tılsım, Ejderha Kemiği Uzun Kılıç, Ejderha Kemiği Yay, Ejderha Kemiği Yüzük, Ejderha Kemiği Çizme, Sis Örgü Ayin Baltası, Sis Örgü Bacaklık, Sis Örgü Balta, Sis Örgü Eldiven, Sis Örgü Göğüslük, Sis Örgü Hançer, Sis Örgü Kalkan, Sis Örgü Kask, Sis Örgü Kolye, Sis Örgü Kılıç, Sis Örgü Kısa Kılıç, Sis Örgü Meç, Sis Örgü Muşta, Sis Örgü Orak, Sis Örgü Savaş Çekici, Sis Örgü Topuz, Sis Örgü Tılsım, Sis Örgü Uzun Kılıç, Sis Örgü Yüzük, Sis Örgü Çizme |
| nickel_bar | Nikel Külçesi | material/bar | 1 | 48c |  | üretim (smelting) | Ejderha Kemiği Ayin Baltası, Ejderha Kemiği Bacaklık, Ejderha Kemiği Balta, Ejderha Kemiği Eldiven, Ejderha Kemiği Göğüslük, Ejderha Kemiği Hançer, Ejderha Kemiği Kalkan, Ejderha Kemiği Kask, Ejderha Kemiği Kolye, Ejderha Kemiği Kılıç, Ejderha Kemiği Kısa Kılıç, Ejderha Kemiği Meç, Ejderha Kemiği Muşta, Ejderha Kemiği Orak, Ejderha Kemiği Savaş Çekici, Ejderha Kemiği Topuz, Ejderha Kemiği Tılsım, Ejderha Kemiği Uzun Kılıç, Ejderha Kemiği Yüzük, Ejderha Kemiği Çizme, Hükümdar Arbalet, Hükümdar Asa, Hükümdar Ayin Baltası, Hükümdar Bacaklık, Hükümdar Balta, Hükümdar Değnek, Hükümdar Doğa Odağı, Hükümdar Eldiven, Hükümdar Göğüslük, Hükümdar Hançer, Hükümdar Kalkan, Hükümdar Kask, Hükümdar Kolye, Hükümdar Kılıç, Hükümdar Kısa Kılıç, Hükümdar Lavta, Hükümdar Meç, Hükümdar Muşta, Hükümdar Orak, Hükümdar Savaş Çekici, Hükümdar Tonfa, Hükümdar Topuz, Hükümdar Totem, Hükümdar Tılsım, Hükümdar Uzun Kılıç, Hükümdar Yay, Hükümdar Yüzük, Hükümdar Çizme, Mercan Arbalet, Mercan Asa, Mercan Ayin Baltası, Mercan Bacaklık, Mercan Balta, Mercan Değnek, Mercan Doğa Odağı, Mercan Eldiven, Mercan Göğüslük, Mercan Hançer, Mercan Kalkan, Mercan Kask, Mercan Kolye, Mercan Kılıç, Mercan Kısa Kılıç, Mercan Lavta, Mercan Meç, Mercan Muşta, Mercan Orak, Mercan Savaş Çekici, Mercan Tonfa, Mercan Topuz, Mercan Totem, Mercan Tılsım, Mercan Uzun Kılıç, Mercan Yay, Mercan Yüzük, Mercan Çizme |
| zinc_bar | Çinko Külçesi | material/bar | 1 | 38c |  | üretim (smelting) | Hükümdar Ayin Baltası, Hükümdar Bacaklık, Hükümdar Balta, Hükümdar Eldiven, Hükümdar Göğüslük, Hükümdar Hançer, Hükümdar Kalkan, Hükümdar Kask, Hükümdar Kolye, Hükümdar Kılıç, Hükümdar Kısa Kılıç, Hükümdar Meç, Hükümdar Muşta, Hükümdar Orak, Hükümdar Savaş Çekici, Hükümdar Topuz, Hükümdar Tılsım, Hükümdar Uzun Kılıç, Hükümdar Yüzük, Hükümdar Çizme, Mercan Ayin Baltası, Mercan Bacaklık, Mercan Balta, Mercan Eldiven, Mercan Göğüslük, Mercan Hançer, Mercan Kalkan, Mercan Kask, Mercan Kolye, Mercan Kılıç, Mercan Kısa Kılıç, Mercan Meç, Mercan Muşta, Mercan Orak, Mercan Savaş Çekici, Mercan Topuz, Mercan Tılsım, Mercan Uzun Kılıç, Mercan Yüzük, Mercan Çizme |
| bronze_bar | Bronz Külçesi | material/bar | 1 | 29c |  | üretim (smelting) | Anka Arbalet, Anka Asa, Anka Ayin Baltası, Anka Bacaklık, Anka Balta, Anka Değnek, Anka Doğa Odağı, Anka Eldiven, Anka Göğüslük, Anka Hançer, Anka Kalkan, Anka Kask, Anka Kolye, Anka Kılıç, Anka Kısa Kılıç, Anka Lavta, Anka Meç, Anka Muşta, Anka Orak, Anka Savaş Çekici, Anka Tonfa, Anka Topuz, Anka Totem, Anka Tılsım, Anka Uzun Kılıç, Anka Yay, Anka Yüzük, Anka Çizme, Dehşet Dövme Arbalet, Dehşet Dövme Asa, Dehşet Dövme Ayin Baltası, Dehşet Dövme Bacaklık, Dehşet Dövme Balta, Dehşet Dövme Değnek, Dehşet Dövme Doğa Odağı, Dehşet Dövme Eldiven, Dehşet Dövme Göğüslük, Dehşet Dövme Hançer, Dehşet Dövme Kalkan, Dehşet Dövme Kask, Dehşet Dövme Kolye, Dehşet Dövme Kılıç, Dehşet Dövme Kısa Kılıç, Dehşet Dövme Lavta, Dehşet Dövme Meç, Dehşet Dövme Muşta, Dehşet Dövme Orak, Dehşet Dövme Savaş Çekici, Dehşet Dövme Tonfa, Dehşet Dövme Topuz, Dehşet Dövme Totem, Dehşet Dövme Tılsım, Dehşet Dövme Uzun Kılıç, Dehşet Dövme Yay, Dehşet Dövme Yüzük, Dehşet Dövme Çizme, Rün Bağlı Arbalet, Rün Bağlı Asa, Rün Bağlı Ayin Baltası, Rün Bağlı Bacaklık, Rün Bağlı Balta, Rün Bağlı Değnek, Rün Bağlı Doğa Odağı, Rün Bağlı Eldiven, Rün Bağlı Göğüslük, Rün Bağlı Hançer, Rün Bağlı Kalkan, Rün Bağlı Kask, Rün Bağlı Kolye, Rün Bağlı Kılıç, Rün Bağlı Kısa Kılıç, Rün Bağlı Lavta, Rün Bağlı Meç, Rün Bağlı Muşta, Rün Bağlı Orak, Rün Bağlı Savaş Çekici, Rün Bağlı Tonfa, Rün Bağlı Topuz, Rün Bağlı Totem, Rün Bağlı Tılsım, Rün Bağlı Uzun Kılıç, Rün Bağlı Yay, Rün Bağlı Yüzük, Rün Bağlı Çizme |
| copper_bar | Bakır Külçesi | material/bar | 1 | 26c |  | üretim (smelting) | Anka Ayin Baltası, Anka Bacaklık, Anka Balta, Anka Eldiven, Anka Göğüslük, Anka Hançer, Anka Kalkan, Anka Kask, Anka Kolye, Anka Kılıç, Anka Kısa Kılıç, Anka Meç, Anka Muşta, Anka Orak, Anka Savaş Çekici, Anka Topuz, Anka Tılsım, Anka Uzun Kılıç, Anka Yüzük, Anka Çizme, Dehşet Dövme Ayin Baltası, Dehşet Dövme Bacaklık, Dehşet Dövme Balta, Dehşet Dövme Eldiven, Dehşet Dövme Göğüslük, Dehşet Dövme Hançer, Dehşet Dövme Kalkan, Dehşet Dövme Kask, Dehşet Dövme Kolye, Dehşet Dövme Kılıç, Dehşet Dövme Kısa Kılıç, Dehşet Dövme Meç, Dehşet Dövme Muşta, Dehşet Dövme Orak, Dehşet Dövme Savaş Çekici, Dehşet Dövme Topuz, Dehşet Dövme Tılsım, Dehşet Dövme Uzun Kılıç, Dehşet Dövme Yüzük, Dehşet Dövme Çizme, Kum Çeliği Arbalet, Kum Çeliği Asa, Kum Çeliği Ayin Baltası, Kum Çeliği Bacaklık, Kum Çeliği Balta, Kum Çeliği Değnek, Kum Çeliği Doğa Odağı, Kum Çeliği Eldiven, Kum Çeliği Göğüslük, Kum Çeliği Hançer, Kum Çeliği Kalkan, Kum Çeliği Kask, Kum Çeliği Kolye, Kum Çeliği Kılıç, Kum Çeliği Kısa Kılıç, Kum Çeliği Lavta, Kum Çeliği Meç, Kum Çeliği Muşta, Kum Çeliği Orak, Kum Çeliği Savaş Çekici, Kum Çeliği Tonfa, Kum Çeliği Topuz, Kum Çeliği Totem, Kum Çeliği Tılsım, Kum Çeliği Uzun Kılıç, Kum Çeliği Yay, Kum Çeliği Yüzük, Kum Çeliği Çizme, Rün Bağlı Ayin Baltası, Rün Bağlı Bacaklık, Rün Bağlı Balta, Rün Bağlı Eldiven, Rün Bağlı Göğüslük, Rün Bağlı Hançer, Rün Bağlı Kalkan, Rün Bağlı Kask, Rün Bağlı Kolye, Rün Bağlı Kılıç, Rün Bağlı Kısa Kılıç, Rün Bağlı Meç, Rün Bağlı Muşta, Rün Bağlı Orak, Rün Bağlı Savaş Çekici, Rün Bağlı Topuz, Rün Bağlı Tılsım, Rün Bağlı Uzun Kılıç, Rün Bağlı Yüzük, Rün Bağlı Çizme |
| tin_bar | Kalay Külçesi | material/bar | 1 | 26c |  | üretim (smelting) | Kum Çeliği Ayin Baltası, Kum Çeliği Bacaklık, Kum Çeliği Balta, Kum Çeliği Eldiven, Kum Çeliği Göğüslük, Kum Çeliği Hançer, Kum Çeliği Kalkan, Kum Çeliği Kask, Kum Çeliği Kolye, Kum Çeliği Kılıç, Kum Çeliği Kısa Kılıç, Kum Çeliği Meç, Kum Çeliği Muşta, Kum Çeliği Orak, Kum Çeliği Savaş Çekici, Kum Çeliği Topuz, Kum Çeliği Tılsım, Kum Çeliği Uzun Kılıç, Kum Çeliği Yüzük, Kum Çeliği Çizme |
| iron_bar | Demir Külçesi | material/bar | 1 | 22c |  | üretim (smelting) | Obsidyen Arbalet, Obsidyen Asa, Obsidyen Ayin Baltası, Obsidyen Bacaklık, Obsidyen Balta, Obsidyen Değnek, Obsidyen Doğa Odağı, Obsidyen Eldiven, Obsidyen Göğüslük, Obsidyen Hançer, Obsidyen Kalkan, Obsidyen Kask, Obsidyen Kolye, Obsidyen Kılıç, Obsidyen Kısa Kılıç, Obsidyen Lavta, Obsidyen Meç, Obsidyen Muşta, Obsidyen Orak, Obsidyen Savaş Çekici, Obsidyen Tonfa, Obsidyen Topuz, Obsidyen Totem, Obsidyen Tılsım, Obsidyen Uzun Kılıç, Obsidyen Yay, Obsidyen Yüzük, Obsidyen Çizme |
| lead_bar | Kurşun Külçesi | material/bar | 1 | 22c |  | üretim (smelting) | Güneş Dövme Arbalet, Güneş Dövme Asa, Güneş Dövme Ayin Baltası, Güneş Dövme Bacaklık, Güneş Dövme Balta, Güneş Dövme Değnek, Güneş Dövme Doğa Odağı, Güneş Dövme Eldiven, Güneş Dövme Göğüslük, Güneş Dövme Hançer, Güneş Dövme Kalkan, Güneş Dövme Kask, Güneş Dövme Kolye, Güneş Dövme Kılıç, Güneş Dövme Kısa Kılıç, Güneş Dövme Lavta, Güneş Dövme Meç, Güneş Dövme Muşta, Güneş Dövme Orak, Güneş Dövme Savaş Çekici, Güneş Dövme Tonfa, Güneş Dövme Topuz, Güneş Dövme Totem, Güneş Dövme Tılsım, Güneş Dövme Uzun Kılıç, Güneş Dövme Yay, Güneş Dövme Yüzük, Güneş Dövme Çizme, Obsidyen Ayin Baltası, Obsidyen Bacaklık, Obsidyen Balta, Obsidyen Eldiven, Obsidyen Göğüslük, Obsidyen Hançer, Obsidyen Kalkan, Obsidyen Kask, Obsidyen Kolye, Obsidyen Kılıç, Obsidyen Kısa Kılıç, Obsidyen Meç, Obsidyen Muşta, Obsidyen Orak, Obsidyen Savaş Çekici, Obsidyen Topuz, Obsidyen Tılsım, Obsidyen Uzun Kılıç, Obsidyen Yüzük, Obsidyen Çizme |
| aluminium_bar | Alüminyum Külçesi | material/bar | 1 | 19c |  | üretim (smelting) | Acemi Arbalet, Acemi Ayin Baltası, Acemi Bacaklık, Acemi Balta, Acemi Değnek, Acemi Doğa Odağı, Acemi Eldiven, Acemi Göğüslük, Acemi Kask, Acemi Kısa Kılıç, Acemi Meç, Acemi Tonfa, Acemi Uzun Kılıç, Acemi Çizme, Adamant Arbalet, Adamant Asa, Adamant Ayin Baltası, Adamant Bacaklık, Adamant Balta, Adamant Değnek, Adamant Doğa Odağı, Adamant Eldiven, Adamant Göğüslük, Adamant Hançer, Adamant Kask, Adamant Kolye, Adamant Kılıç, Adamant Kısa Kılıç, Adamant Lavta, Adamant Meç, Adamant Muşta, Adamant Orak, Adamant Savaş Çekici, Adamant Tonfa, Adamant Topuz, Adamant Totem, Adamant Tılsım, Adamant Uzun Kılıç, Adamant Yay, Adamant Yüzük, Adamant Çizme, Akçaağaç Lavta, Ata Totemi, Av Yayı, Ay Orağı, Bez Sargılar, Bileşik Yay, Bronz Orak, Büyülü Lavta, Demir Bacaklık, Demir Eldiven, Demir Göğüslük, Demir Hançer, Demir Kask, Demir Kılıç, Demir Muşta, Demir Orak, Demir Savaş Çekici, Demir Topuz, Demir Çizme, Dikenli Orak, Efsaneler Lavtası, Ejder Dövme Arbalet, Ejder Dövme Asa, Ejder Dövme Ayin Baltası, Ejder Dövme Bacaklık, Ejder Dövme Balta, Ejder Dövme Değnek, Ejder Dövme Doğa Odağı, Ejder Dövme Eldiven, Ejder Dövme Göğüslük, Ejder Dövme Hançer, Ejder Dövme Kask, Ejder Dövme Kolye, Ejder Dövme Kılıç, Ejder Dövme Kısa Kılıç, Ejder Dövme Lavta, Ejder Dövme Meç, Ejder Dövme Muşta, Ejder Dövme Orak, Ejder Dövme Savaş Çekici, Ejder Dövme Tonfa, Ejder Dövme Topuz, Ejder Dövme Totem, Ejder Dövme Tılsım, Ejder Dövme Uzun Kılıç, Ejder Dövme Yay, Ejder Dövme Yüzük, Ejder Dövme Çizme, Ejder Kemiği Yayı, Ejder Yumrukları, Gizemli Asa, Gök Gürültüsü Putu, Göksel Yumruklar, Gölge Hançeri, Gümüş Savaş Çekici, Gümüş Telli Lavta, Gümüş Topuz, Güneş Dövme Ayin Baltası, Güneş Dövme Bacaklık, Güneş Dövme Balta, Güneş Dövme Eldiven, Güneş Dövme Göğüslük, Güneş Dövme Hançer, Güneş Dövme Kalkan, Güneş Dövme Kask, Güneş Dövme Kolye, Güneş Dövme Kılıç, Güneş Dövme Kısa Kılıç, Güneş Dövme Meç, Güneş Dövme Muşta, Güneş Dövme Orak, Güneş Dövme Savaş Çekici, Güneş Dövme Topuz, Güneş Dövme Tılsım, Güneş Dövme Uzun Kılıç, Güneş Dövme Yüzük, Güneş Dövme Çizme, Işıltılı Çekiç, Kadim Orak, Kalfa Arbalet, Kalfa Ayin Baltası, Kalfa Balta, Kalfa Değnek, Kalfa Doğa Odağı, Kalfa Kısa Kılıç, Kalfa Meç, Kalfa Tonfa, Kalfa Uzun Kılıç, Kaliteli Arbalet, Kaliteli Ayin Baltası, Kaliteli Balta, Kaliteli Değnek, Kaliteli Doğa Odağı, Kaliteli Kısa Kılıç, Kaliteli Meç, Kaliteli Tonfa, Kaliteli Uzun Kılıç, Kaplan Pençeleri, Kemik Totem, Konser Lavtası, Koru Orağı, Kraliyet Arbalet, Kraliyet Asa, Kraliyet Ayin Baltası, Kraliyet Bacaklık, Kraliyet Balta, Kraliyet Değnek, Kraliyet Doğa Odağı, Kraliyet Eldiven, Kraliyet Göğüslük, Kraliyet Hançer, Kraliyet Kalkan, Kraliyet Kask, Kraliyet Kolye, Kraliyet Kılıç, Kraliyet Kısa Kılıç, Kraliyet Lavta, Kraliyet Meç, Kraliyet Muşta, Kraliyet Orak, Kraliyet Savaş Çekici, Kraliyet Tonfa, Kraliyet Topuz, Kraliyet Totem, Kraliyet Tılsım, Kraliyet Uzun Kılıç, Kraliyet Yay, Kraliyet Yüzük, Kraliyet Çizme, Kristal Asa, Kutsal Topuz, Kutsal İntikamcı, Kutsanmış Topuz, Kutsanmış Çekiç, Kıdemli Arbalet, Kıdemli Ayin Baltası, Kıdemli Balta, Kıdemli Değnek, Kıdemli Doğa Odağı, Kıdemli Kısa Kılıç, Kıdemli Meç, Kıdemli Tonfa, Kıdemli Uzun Kılıç, Kısa Yay, Meşe Tokmak, Mithril Arbalet, Mithril Asa, Mithril Ayin Baltası, Mithril Bacaklık, Mithril Balta, Mithril Değnek, Mithril Doğa Odağı, Mithril Eldiven, Mithril Göğüslük, Mithril Hançer, Mithril Kask, Mithril Kolye, Mithril Kılıç, Mithril Kısa Kılıç, Mithril Lavta, Mithril Meç, Mithril Muşta, Mithril Orak, Mithril Savaş Çekici, Mithril Tonfa, Mithril Topuz, Mithril Totem, Mithril Tılsım, Mithril Uzun Kılıç, Mithril Yay, Mithril Yüzük, Mithril Çizme, Oyma Totem, Paslı Hançer, Porsuk Yayı, Pul Bacaklık, Pul Eldiven, Pul Göğüslük, Pul Kask, Pul Çizme, Ruh Totemi, Rünlü Arbalet, Rünlü Ayin Baltası, Rünlü Bacaklık, Rünlü Balta, Rünlü Büyük Kılıç, Rünlü Değnek, Rünlü Doğa Odağı, Rünlü Eldiven, Rünlü Göğüslük, Rünlü Kask, Rünlü Kolye, Rünlü Kısa Kılıç, Rünlü Meç, Rünlü Tonfa, Rünlü Tılsım, Rünlü Uzun Kılıç, Rünlü Yüzük, Rünlü Çizme, Sade Lavta, Su Verilmiş Kılıç, Söğüt Asa, Tahta Topuz, Talim Kılıcı, Taş Put, Titan Arbalet, Titan Asa, Titan Ayin Baltası, Titan Bacaklık, Titan Balta, Titan Değnek, Titan Doğa Odağı, Titan Eldiven, Titan Göğüslük, Titan Hançer, Titan Kask, Titan Kolye, Titan Kılıç, Titan Kısa Kılıç, Titan Lavta, Titan Meç, Titan Muşta, Titan Orak, Titan Savaş Çekici, Titan Tonfa, Titan Topuz, Titan Totem, Titan Tılsım, Titan Uzun Kılıç, Titan Yay, Titan Yüzük, Titan Çizme, Tırtıklı Hançer, Usta Asası, Uzun Yay, Yıldız Asası, Yıldız Dövme Arbalet, Yıldız Dövme Asa, Yıldız Dövme Ayin Baltası, Yıldız Dövme Bacaklık, Yıldız Dövme Balta, Yıldız Dövme Değnek, Yıldız Dövme Doğa Odağı, Yıldız Dövme Eldiven, Yıldız Dövme Göğüslük, Yıldız Dövme Hançer, Yıldız Dövme Kask, Yıldız Dövme Kolye, Yıldız Dövme Kılıç, Yıldız Dövme Kısa Kılıç, Yıldız Dövme Lavta, Yıldız Dövme Meç, Yıldız Dövme Muşta, Yıldız Dövme Orak, Yıldız Dövme Savaş Çekici, Yıldız Dövme Tonfa, Yıldız Dövme Topuz, Yıldız Dövme Totem, Yıldız Dövme Tılsım, Yıldız Dövme Uzun Kılıç, Yıldız Dövme Yay, Yıldız Dövme Yüzük, Yıldız Dövme Çizme, Zehir Dişi, Çelik Bacaklık, Çelik Cestus, Çelik Eldiven, Çelik Göğüslük, Çelik Hançer, Çelik Kask, Çelik Uzun Kılıç, Çelik Çizme, Çırak Asası, İlahi Topuz, Şampiyon Arbalet, Şampiyon Ayin Baltası, Şampiyon Balta, Şampiyon Değnek, Şampiyon Doğa Odağı, Şampiyon Kısa Kılıç, Şampiyon Meç, Şampiyon Tonfa, Şampiyon Uzun Kılıç, Şövalye Bacaklığı, Şövalye Eldiveni, Şövalye Göğüslüğü, Şövalye Kaskı, Şövalye Kılıcı, Şövalye Çizmesi |
| linseed | Keten Tohumu | material/herb | 1 | 7c |  | **yok** | — |
| mandrake | Adamotu | material/herb | 5 | 12c |  | **yok** | — |
| king_bloom | Kral Çiçeği | material/herb | 20 | 60c |  | **yok** | — |
| elm_log | Karaağaç Kütüğü | material/wood | 8 | 9c |  | **yok** | — |
| yew_log | Porsuk Kütüğü | material/wood | 20 | 30c |  | **yok** | — |

## 3. Üretim listesi

### Simya — 6 tarif, beceri 1–85 (tavan 100), çıktı Sv 1–850

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_potion_hp_small | 1 | 3× Sıradan Ot | 1× Küçük Can İksiri | 1 | 35s | 16 | 5c |
| recipe_potion_hp_medium | 10 | 2× Papatya, 1× Karahindiba | 1× Orta Boy Can İksiri | 100 | 45s | 27 | 15c |
| recipe_potion_hp_large | 25 | 2× Kızıl Adaçayı, 1× Lavanta | 1× Büyük Boy Can İksiri | 250 | 65s | 64 | 37c |
| recipe_potion_hp_xlarge | 45 | 2× Fırtına Çiçeği, 1× Mavi Devedikeni | 1× Ekstra Büyük Can İksiri | 450 | 90s | 206 | 80c |
| recipe_potion_hp_huge | 65 | 2× Ruh Nilüferi, 1× Altın Yonca | 1× Çok Büyük Can İksiri | 650 | 115s | 662 | 1s 50c |
| recipe_potion_hp_supreme | 85 | 2× Boşluk Orkidesi, 1× Hayalet Mantarı | 1× Yüce Can İksiri | 850 | 135s | 2124 | 2s 50c |

### Demircilik — 788 tarif, beceri 1–99 (tavan 100), çıktı Sv 1–985

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_axe_novice | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Acemi Balta | 1 | 155s | 26 | 4c |
| recipe_dagger_rusty | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Paslı Hançer | 1 | 155s | 26 | 4c |
| recipe_fist_wraps | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Bez Sargılar | 1 | 155s | 26 | 4c |
| recipe_mace_wooden | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Tahta Topuz | 1 | 155s | 26 | 4c |
| recipe_boots_novice | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Acemi Çizme | 1 | 95s | 16 | 2c |
| recipe_chest_novice | 1 | 5× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Acemi Göğüslük | 1 | 155s | 26 | 3c |
| recipe_gloves_novice | 1 | 2× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Acemi Eldiven | 1 | 75s | 13 | 2c |
| recipe_helm_novice | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Acemi Kask | 1 | 95s | 16 | 2c |
| recipe_legs_novice | 1 | 4× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Acemi Bacaklık | 1 | 135s | 22 | 3c |
| recipe_rapier_novice | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Acemi Meç | 1 | 155s | 26 | 4c |
| recipe_ritualaxe_novice | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Acemi Ayin Baltası | 1 | 155s | 26 | 4c |
| recipe_shortsword_novice | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Acemi Kısa Kılıç | 1 | 155s | 26 | 4c |
| recipe_sickle_bronze | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Bronz Orak | 1 | 155s | 26 | 4c |
| recipe_longsword_novice | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Acemi Uzun Kılıç | 1 | 155s | 26 | 4c |
| recipe_sword_training | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Talim Kılıcı | 1 | 155s | 26 | 4c |
| recipe_hammer_oak | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Meşe Tokmak | 1 | 155s | 26 | 4c |
| recipe_axe_journeyman | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kalfa Balta | 5 | 155s | 26 | 16c |
| recipe_dagger_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Hançer | 5 | 155s | 26 | 16c |
| recipe_fist_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Muşta | 5 | 155s | 26 | 16c |
| recipe_mace_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Topuz | 5 | 155s | 26 | 16c |
| recipe_boots_iron | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Demir Çizme | 5 | 95s | 16 | 7c |
| recipe_chest_iron | 1 | 5× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Demir Göğüslük | 5 | 155s | 26 | 13c |
| recipe_gloves_iron | 1 | 2× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Demir Eldiven | 5 | 75s | 13 | 7c |
| recipe_helm_iron | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Demir Kask | 5 | 95s | 16 | 8c |
| recipe_legs_iron | 1 | 4× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Demir Bacaklık | 5 | 135s | 22 | 12c |
| recipe_rapier_journeyman | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kalfa Meç | 5 | 155s | 26 | 16c |
| recipe_ritualaxe_journeyman | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kalfa Ayin Baltası | 5 | 155s | 26 | 16c |
| recipe_shortsword_journeyman | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kalfa Kısa Kılıç | 5 | 155s | 26 | 16c |
| recipe_sickle_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Orak | 5 | 155s | 26 | 16c |
| recipe_longsword_journeyman | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kalfa Uzun Kılıç | 5 | 155s | 26 | 16c |
| recipe_sword_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Kılıç | 5 | 155s | 26 | 16c |
| recipe_hammer_iron | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Demir Savaş Çekici | 5 | 155s | 26 | 16c |
| recipe_axe_fine | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaliteli Balta | 10 | 155s | 26 | 41c |
| recipe_dagger_serrated | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Tırtıklı Hançer | 10 | 155s | 26 | 41c |
| recipe_fist_tiger | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaplan Pençeleri | 10 | 155s | 26 | 41c |
| recipe_mace_blessed | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kutsanmış Topuz | 10 | 155s | 26 | 41c |
| recipe_boots_scale | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Pul Çizme | 10 | 95s | 16 | 18c |
| recipe_chest_scale | 1 | 5× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Pul Göğüslük | 10 | 155s | 26 | 33c |
| recipe_gloves_scale | 1 | 2× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Pul Eldiven | 10 | 75s | 13 | 16c |
| recipe_helm_scale | 1 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Pul Kask | 10 | 95s | 16 | 20c |
| recipe_legs_scale | 1 | 4× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Pul Bacaklık | 10 | 135s | 22 | 29c |
| recipe_rapier_fine | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaliteli Meç | 10 | 155s | 26 | 41c |
| recipe_ritualaxe_fine | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaliteli Ayin Baltası | 10 | 155s | 26 | 41c |
| recipe_shortsword_fine | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaliteli Kısa Kılıç | 10 | 155s | 26 | 41c |
| recipe_sickle_moon | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Ay Orağı | 10 | 155s | 26 | 41c |
| recipe_longsword_fine | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kaliteli Uzun Kılıç | 10 | 155s | 26 | 41c |
| recipe_sword_tempered | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Su Verilmiş Kılıç | 10 | 155s | 26 | 41c |
| recipe_hammer_blessed | 1 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kutsanmış Çekiç | 10 | 155s | 26 | 41c |
| recipe_axe_veteran | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kıdemli Balta | 15 | 165s | 27 | 73c |
| recipe_dagger_steel | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Çelik Hançer | 15 | 165s | 27 | 73c |
| recipe_fist_steel | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Çelik Cestus | 15 | 165s | 27 | 73c |
| recipe_mace_silver | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Gümüş Topuz | 15 | 165s | 27 | 73c |
| recipe_boots_steel | 2 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Çelik Çizme | 15 | 100s | 17 | 33c |
| recipe_chest_steel | 2 | 5× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Çelik Göğüslük | 15 | 165s | 27 | 58c |
| recipe_gloves_steel | 2 | 2× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Çelik Eldiven | 15 | 80s | 13 | 29c |
| recipe_helm_steel | 2 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Çelik Kask | 15 | 100s | 17 | 36c |
| recipe_legs_steel | 2 | 4× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Çelik Bacaklık | 15 | 145s | 23 | 51c |
| recipe_rapier_veteran | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kıdemli Meç | 15 | 165s | 27 | 73c |
| recipe_ritualaxe_veteran | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kıdemli Ayin Baltası | 15 | 165s | 27 | 73c |
| recipe_shortsword_veteran | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kıdemli Kısa Kılıç | 15 | 165s | 27 | 73c |
| recipe_sickle_thorn | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Dikenli Orak | 15 | 165s | 27 | 73c |
| recipe_longsword_veteran | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kıdemli Uzun Kılıç | 15 | 165s | 27 | 73c |
| recipe_sword_steel | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Çelik Uzun Kılıç | 15 | 165s | 27 | 73c |
| recipe_hammer_silver | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Gümüş Savaş Çekici | 15 | 165s | 27 | 73c |
| recipe_axe_champion | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şampiyon Balta | 20 | 165s | 27 | 1s 10c |
| recipe_dagger_shadow | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Gölge Hançeri | 20 | 165s | 27 | 1s 10c |
| recipe_fist_dragon | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Ejder Yumrukları | 20 | 165s | 27 | 1s 10c |
| recipe_mace_holy | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Kutsal Topuz | 20 | 165s | 27 | 1s 10c |
| recipe_boots_knight | 2 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Şövalye Çizmesi | 20 | 100s | 17 | 50c |
| recipe_chest_knight | 2 | 5× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Şövalye Göğüslüğü | 20 | 165s | 27 | 88c |
| recipe_gloves_knight | 2 | 2× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Şövalye Eldiveni | 20 | 80s | 13 | 44c |
| recipe_helm_knight | 2 | 3× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Şövalye Kaskı | 20 | 100s | 17 | 55c |
| recipe_legs_knight | 2 | 4× Alüminyum Külçesi, 1× Orman Faresi Postu | 1× Şövalye Bacaklığı | 20 | 145s | 23 | 77c |
| recipe_rapier_champion | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şampiyon Meç | 20 | 165s | 27 | 1s 10c |
| recipe_ritualaxe_champion | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şampiyon Ayin Baltası | 20 | 165s | 27 | 1s 10c |
| recipe_shortsword_champion | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şampiyon Kısa Kılıç | 20 | 165s | 27 | 1s 10c |
| recipe_sickle_grove | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Koru Orağı | 20 | 165s | 27 | 1s 10c |
| recipe_longsword_champion | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şampiyon Uzun Kılıç | 20 | 165s | 27 | 1s 10c |
| recipe_sword_knight | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Şövalye Kılıcı | 20 | 165s | 27 | 1s 10c |
| recipe_hammer_radiant | 2 | 4× Alüminyum Külçesi, 2× Fidan Dalı | 1× Işıltılı Çekiç | 20 | 165s | 27 | 1s 10c |
| recipe_axe_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Balta | 25 | 175s | 28 | 1s 53c |
| recipe_dagger_venom | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Zehir Dişi | 25 | 175s | 28 | 1s 53c |
| recipe_fist_celestial | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Göksel Yumruklar | 25 | 175s | 28 | 1s 53c |
| recipe_mace_divine | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× İlahi Topuz | 25 | 175s | 28 | 1s 53c |
| recipe_boots_runed | 3 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Rünlü Çizme | 25 | 110s | 18 | 69c |
| recipe_chest_runed | 3 | 5× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Rünlü Göğüslük | 25 | 175s | 28 | 1s 22c |
| recipe_gloves_runed | 3 | 2× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Rünlü Eldiven | 25 | 85s | 14 | 61c |
| recipe_helm_runed | 3 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Rünlü Kask | 25 | 110s | 18 | 76c |
| recipe_legs_runed | 3 | 4× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Rünlü Bacaklık | 25 | 150s | 25 | 1s 7c |
| recipe_rapier_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Meç | 25 | 175s | 28 | 1s 53c |
| recipe_ritualaxe_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Ayin Baltası | 25 | 175s | 28 | 1s 53c |
| recipe_shortsword_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Kısa Kılıç | 25 | 175s | 28 | 1s 53c |
| recipe_sickle_ancient | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Kadim Orak | 25 | 175s | 28 | 1s 53c |
| recipe_longsword_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Uzun Kılıç | 25 | 175s | 28 | 1s 53c |
| recipe_sword_runed | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Rünlü Büyük Kılıç | 25 | 175s | 28 | 1s 53c |
| recipe_hammer_sacred | 3 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Kutsal İntikamcı | 25 | 175s | 28 | 1s 53c |
| recipe_axe_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Balta | 35 | 180s | 30 | 2s 37c |
| recipe_dagger_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Hançer | 35 | 180s | 30 | 2s 37c |
| recipe_fist_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Muşta | 35 | 180s | 30 | 2s 37c |
| recipe_mace_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Topuz | 35 | 180s | 30 | 2s 37c |
| recipe_boots_mithril | 4 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Mithril Çizme | 35 | 115s | 19 | 1s 7c |
| recipe_chest_mithril | 4 | 5× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Mithril Göğüslük | 35 | 180s | 30 | 1s 90c |
| recipe_gloves_mithril | 4 | 2× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Mithril Eldiven | 35 | 90s | 15 | 95c |
| recipe_helm_mithril | 4 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Mithril Kask | 35 | 115s | 19 | 1s 19c |
| recipe_legs_mithril | 4 | 4× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Mithril Bacaklık | 35 | 160s | 27 | 1s 66c |
| recipe_rapier_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Meç | 35 | 180s | 30 | 2s 37c |
| recipe_ritualaxe_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Ayin Baltası | 35 | 180s | 30 | 2s 37c |
| recipe_shortsword_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Kısa Kılıç | 35 | 180s | 30 | 2s 37c |
| recipe_sickle_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Orak | 35 | 180s | 30 | 2s 37c |
| recipe_longsword_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Uzun Kılıç | 35 | 180s | 30 | 2s 37c |
| recipe_sword_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Kılıç | 35 | 180s | 30 | 2s 37c |
| recipe_hammer_mithril | 4 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Mithril Savaş Çekici | 35 | 180s | 30 | 2s 37c |
| recipe_axe_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Balta | 50 | 190s | 32 | 3s 76c |
| recipe_dagger_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Hançer | 50 | 190s | 32 | 3s 76c |
| recipe_fist_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Muşta | 50 | 190s | 32 | 3s 76c |
| recipe_mace_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Topuz | 50 | 190s | 32 | 3s 76c |
| recipe_boots_adamant | 5 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Adamant Çizme | 50 | 120s | 20 | 1s 69c |
| recipe_chest_adamant | 5 | 5× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Adamant Göğüslük | 50 | 190s | 32 | 3s 1c |
| recipe_gloves_adamant | 5 | 2× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Adamant Eldiven | 50 | 95s | 16 | 1s 51c |
| recipe_helm_adamant | 5 | 3× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Adamant Kask | 50 | 120s | 20 | 1s 88c |
| recipe_legs_adamant | 5 | 4× Alüminyum Külçesi, 1× Küçük Tavşan Postu | 1× Adamant Bacaklık | 50 | 170s | 28 | 2s 64c |
| recipe_rapier_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Meç | 50 | 190s | 32 | 3s 76c |
| recipe_ritualaxe_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Ayin Baltası | 50 | 190s | 32 | 3s 76c |
| recipe_shortsword_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Kısa Kılıç | 50 | 190s | 32 | 3s 76c |
| recipe_sickle_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Orak | 50 | 190s | 32 | 3s 76c |
| recipe_longsword_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Uzun Kılıç | 50 | 190s | 32 | 3s 76c |
| recipe_sword_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Kılıç | 50 | 190s | 32 | 3s 76c |
| recipe_hammer_adamant | 5 | 4× Alüminyum Külçesi, 2× Kavak Kütüğü | 1× Adamant Savaş Çekici | 50 | 190s | 32 | 3s 76c |
| recipe_axe_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Balta | 65 | 210s | 36 | 5s 29c |
| recipe_dagger_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Hançer | 65 | 210s | 36 | 5s 29c |
| recipe_fist_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Muşta | 65 | 210s | 36 | 5s 29c |
| recipe_mace_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Topuz | 65 | 210s | 36 | 5s 29c |
| recipe_boots_wyrmforged | 7 | 3× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Ejder Dövme Çizme | 65 | 130s | 22 | 2s 38c |
| recipe_chest_wyrmforged | 7 | 5× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Ejder Dövme Göğüslük | 65 | 210s | 36 | 4s 24c |
| recipe_gloves_wyrmforged | 7 | 2× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Ejder Dövme Eldiven | 65 | 105s | 18 | 2s 12c |
| recipe_helm_wyrmforged | 7 | 3× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Ejder Dövme Kask | 65 | 130s | 22 | 2s 65c |
| recipe_legs_wyrmforged | 7 | 4× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Ejder Dövme Bacaklık | 65 | 185s | 31 | 3s 71c |
| recipe_rapier_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Meç | 65 | 210s | 36 | 5s 29c |
| recipe_ritualaxe_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Ayin Baltası | 65 | 210s | 36 | 5s 29c |
| recipe_shortsword_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Kısa Kılıç | 65 | 210s | 36 | 5s 29c |
| recipe_sickle_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Orak | 65 | 210s | 36 | 5s 29c |
| recipe_longsword_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Uzun Kılıç | 65 | 210s | 36 | 5s 29c |
| recipe_sword_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Kılıç | 65 | 210s | 36 | 5s 29c |
| recipe_hammer_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Ejder Dövme Savaş Çekici | 65 | 210s | 36 | 5s 29c |
| recipe_axe_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Balta | 80 | 220s | 38 | 6s 88c |
| recipe_dagger_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Hançer | 80 | 220s | 38 | 6s 88c |
| recipe_fist_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Muşta | 80 | 220s | 38 | 6s 88c |
| recipe_mace_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Topuz | 80 | 220s | 38 | 6s 88c |
| recipe_boots_starforged | 8 | 3× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Yıldız Dövme Çizme | 80 | 140s | 24 | 3s 10c |
| recipe_chest_starforged | 8 | 5× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Yıldız Dövme Göğüslük | 80 | 220s | 38 | 5s 51c |
| recipe_gloves_starforged | 8 | 2× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Yıldız Dövme Eldiven | 80 | 110s | 19 | 2s 75c |
| recipe_helm_starforged | 8 | 3× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Yıldız Dövme Kask | 80 | 140s | 24 | 3s 44c |
| recipe_legs_starforged | 8 | 4× Alüminyum Külçesi, 1× Tarla Faresi Postu | 1× Yıldız Dövme Bacaklık | 80 | 195s | 34 | 4s 82c |
| recipe_rapier_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Meç | 80 | 220s | 38 | 6s 88c |
| recipe_ritualaxe_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Ayin Baltası | 80 | 220s | 38 | 6s 88c |
| recipe_shortsword_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Kısa Kılıç | 80 | 220s | 38 | 6s 88c |
| recipe_sickle_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Orak | 80 | 220s | 38 | 6s 88c |
| recipe_longsword_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Uzun Kılıç | 80 | 220s | 38 | 6s 88c |
| recipe_sword_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Kılıç | 80 | 220s | 38 | 6s 88c |
| recipe_hammer_starforged | 8 | 4× Alüminyum Külçesi, 2× Söğüt Kütüğü | 1× Yıldız Dövme Savaş Çekici | 80 | 220s | 38 | 6s 88c |
| recipe_axe_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Balta | 95 | 240s | 43 | 8s 66c |
| recipe_dagger_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Hançer | 95 | 240s | 43 | 8s 66c |
| recipe_fist_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Muşta | 95 | 240s | 43 | 8s 66c |
| recipe_mace_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Topuz | 95 | 240s | 43 | 8s 66c |
| recipe_boots_titan | 10 | 3× Alüminyum Külçesi, 1× Sincap Postu | 1× Titan Çizme | 95 | 150s | 27 | 3s 89c |
| recipe_chest_titan | 10 | 5× Alüminyum Külçesi, 1× Sincap Postu | 1× Titan Göğüslük | 95 | 240s | 43 | 6s 93c |
| recipe_gloves_titan | 10 | 2× Alüminyum Külçesi, 1× Sincap Postu | 1× Titan Eldiven | 95 | 120s | 21 | 3s 46c |
| recipe_helm_titan | 10 | 3× Alüminyum Külçesi, 1× Sincap Postu | 1× Titan Kask | 95 | 150s | 27 | 4s 33c |
| recipe_legs_titan | 10 | 4× Alüminyum Külçesi, 1× Sincap Postu | 1× Titan Bacaklık | 95 | 210s | 37 | 6s 6c |
| recipe_rapier_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Meç | 95 | 240s | 43 | 8s 66c |
| recipe_ritualaxe_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Ayin Baltası | 95 | 240s | 43 | 8s 66c |
| recipe_shortsword_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Kısa Kılıç | 95 | 240s | 43 | 8s 66c |
| recipe_sickle_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Orak | 95 | 240s | 43 | 8s 66c |
| recipe_longsword_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Uzun Kılıç | 95 | 240s | 43 | 8s 66c |
| recipe_sword_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Kılıç | 95 | 240s | 43 | 8s 66c |
| recipe_hammer_titan | 10 | 4× Alüminyum Külçesi, 2× Köknar Kütüğü | 1× Titan Savaş Çekici | 95 | 240s | 43 | 8s 66c |
| recipe_axe_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Balta | 110 | 250s | 46 | 10s 50c |
| recipe_dagger_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Hançer | 110 | 250s | 46 | 10s 50c |
| recipe_fist_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Muşta | 110 | 250s | 46 | 10s 50c |
| recipe_mace_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Topuz | 110 | 250s | 46 | 10s 50c |
| recipe_boots_royal | 11 | 3× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Çizme | 110 | 155s | 28 | 4s 72c |
| recipe_chest_royal | 11 | 5× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Göğüslük | 110 | 250s | 46 | 8s 40c |
| recipe_gloves_royal | 11 | 2× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Eldiven | 110 | 125s | 23 | 4s 20c |
| recipe_helm_royal | 11 | 3× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Kask | 110 | 155s | 28 | 5s 25c |
| recipe_legs_royal | 11 | 4× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Bacaklık | 110 | 220s | 40 | 7s 35c |
| recipe_rapier_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Meç | 110 | 250s | 46 | 10s 50c |
| recipe_ritualaxe_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Ayin Baltası | 110 | 250s | 46 | 10s 50c |
| recipe_shield_royal | 11 | 4× Alüminyum Külçesi, 1× Yaban Tavşanı Postu | 1× Kraliyet Kalkan | 110 | 220s | 40 | 7s 55c |
| recipe_shortsword_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Kısa Kılıç | 110 | 250s | 46 | 10s 50c |
| recipe_sickle_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Orak | 110 | 250s | 46 | 10s 50c |
| recipe_longsword_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Uzun Kılıç | 110 | 250s | 46 | 10s 50c |
| recipe_sword_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Kılıç | 110 | 250s | 46 | 10s 50c |
| recipe_hammer_royal | 11 | 4× Alüminyum Külçesi, 2× Çam Kütüğü | 1× Kraliyet Savaş Çekici | 110 | 250s | 46 | 10s 50c |
| recipe_axe_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Balta | 135 | 280s | 54 | 13s 70c |
| recipe_dagger_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Hançer | 135 | 280s | 54 | 13s 70c |
| recipe_fist_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Muşta | 135 | 280s | 54 | 13s 70c |
| recipe_mace_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Topuz | 135 | 280s | 54 | 13s 70c |
| recipe_boots_sunforged | 14 | 2× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Çizme | 135 | 175s | 34 | 6s 16c |
| recipe_chest_sunforged | 14 | 4× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Göğüslük | 135 | 280s | 54 | 10s 96c |
| recipe_gloves_sunforged | 14 | 1× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Eldiven | 135 | 140s | 27 | 5s 48c |
| recipe_helm_sunforged | 14 | 2× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Kask | 135 | 175s | 34 | 6s 85c |
| recipe_legs_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Bacaklık | 135 | 245s | 48 | 9s 59c |
| recipe_rapier_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Meç | 135 | 280s | 54 | 13s 70c |
| recipe_ritualaxe_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Ayin Baltası | 135 | 280s | 54 | 13s 70c |
| recipe_shield_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 1× Porsuk Postu | 1× Güneş Dövme Kalkan | 135 | 245s | 48 | 9s 85c |
| recipe_shortsword_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Kısa Kılıç | 135 | 280s | 54 | 13s 70c |
| recipe_sickle_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Orak | 135 | 280s | 54 | 13s 70c |
| recipe_longsword_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Uzun Kılıç | 135 | 280s | 54 | 13s 70c |
| recipe_sword_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Kılıç | 135 | 280s | 54 | 13s 70c |
| recipe_hammer_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Huş Kütüğü | 1× Güneş Dövme Savaş Çekici | 135 | 280s | 54 | 13s 70c |
| recipe_axe_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Balta | 160 | 300s | 61 | 17s 9c |
| recipe_dagger_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Hançer | 160 | 300s | 61 | 17s 9c |
| recipe_fist_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Muşta | 160 | 300s | 61 | 17s 9c |
| recipe_mace_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Topuz | 160 | 300s | 61 | 17s 9c |
| recipe_boots_obsidian | 16 | 2× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Çizme | 160 | 185s | 38 | 7s 68c |
| recipe_chest_obsidian | 16 | 4× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Göğüslük | 160 | 300s | 61 | 13s 67c |
| recipe_gloves_obsidian | 16 | 1× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Eldiven | 160 | 150s | 30 | 6s 84c |
| recipe_helm_obsidian | 16 | 2× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Kask | 160 | 185s | 38 | 8s 54c |
| recipe_legs_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Bacaklık | 160 | 260s | 53 | 11s 96c |
| recipe_rapier_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Meç | 160 | 300s | 61 | 17s 9c |
| recipe_ritualaxe_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Ayin Baltası | 160 | 300s | 61 | 17s 9c |
| recipe_shield_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 1× Porsuk Postu | 1× Obsidyen Kalkan | 160 | 260s | 53 | 12s 29c |
| recipe_shortsword_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Kısa Kılıç | 160 | 300s | 61 | 17s 9c |
| recipe_sickle_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Orak | 160 | 300s | 61 | 17s 9c |
| recipe_longsword_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Uzun Kılıç | 160 | 300s | 61 | 17s 9c |
| recipe_sword_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Kılıç | 160 | 300s | 61 | 17s 9c |
| recipe_hammer_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Sedir Kütüğü | 1× Obsidyen Savaş Çekici | 160 | 300s | 61 | 17s 9c |
| recipe_axe_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Balta | 185 | 325s | 73 | 20s 64c |
| recipe_dagger_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Hançer | 185 | 325s | 73 | 20s 64c |
| recipe_fist_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Muşta | 185 | 325s | 73 | 20s 64c |
| recipe_mace_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Topuz | 185 | 325s | 73 | 20s 64c |
| recipe_boots_sandsteel | 19 | 2× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Çizme | 185 | 205s | 46 | 9s 28c |
| recipe_chest_sandsteel | 19 | 4× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Göğüslük | 185 | 325s | 73 | 16s 51c |
| recipe_gloves_sandsteel | 19 | 1× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Eldiven | 185 | 165s | 36 | 8s 26c |
| recipe_helm_sandsteel | 19 | 2× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Kask | 185 | 205s | 46 | 10s 32c |
| recipe_legs_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Bacaklık | 185 | 285s | 64 | 14s 45c |
| recipe_rapier_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Meç | 185 | 325s | 73 | 20s 64c |
| recipe_ritualaxe_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Ayin Baltası | 185 | 325s | 73 | 20s 64c |
| recipe_shield_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 1× Kızıl Tilki Postu | 1× Kum Çeliği Kalkan | 185 | 285s | 64 | 14s 84c |
| recipe_shortsword_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Kısa Kılıç | 185 | 325s | 73 | 20s 64c |
| recipe_sickle_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Orak | 185 | 325s | 73 | 20s 64c |
| recipe_longsword_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Uzun Kılıç | 185 | 325s | 73 | 20s 64c |
| recipe_sword_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Kılıç | 185 | 325s | 73 | 20s 64c |
| recipe_hammer_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Kayın Kütüğü | 1× Kum Çeliği Savaş Çekici | 185 | 325s | 73 | 20s 64c |
| recipe_axe_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Balta | 210 | 345s | 82 | 24s 34c |
| recipe_dagger_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Hançer | 210 | 345s | 82 | 24s 34c |
| recipe_fist_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Muşta | 210 | 345s | 82 | 24s 34c |
| recipe_mace_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Topuz | 210 | 345s | 82 | 24s 34c |
| recipe_boots_runebound | 21 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Çizme | 210 | 215s | 51 | 10s 94c |
| recipe_chest_runebound | 21 | 4× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Göğüslük | 210 | 345s | 82 | 19s 47c |
| recipe_gloves_runebound | 21 | 1× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Eldiven | 210 | 175s | 41 | 9s 74c |
| recipe_helm_runebound | 21 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Kask | 210 | 215s | 51 | 12s 17c |
| recipe_legs_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Bacaklık | 210 | 300s | 71 | 17s 4c |
| recipe_rapier_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Meç | 210 | 345s | 82 | 24s 34c |
| recipe_ritualaxe_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Ayin Baltası | 210 | 345s | 82 | 24s 34c |
| recipe_shield_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Domuzu Postu | 1× Rün Bağlı Kalkan | 210 | 300s | 71 | 17s 50c |
| recipe_shortsword_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Kısa Kılıç | 210 | 345s | 82 | 24s 34c |
| recipe_sickle_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Orak | 210 | 345s | 82 | 24s 34c |
| recipe_longsword_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Uzun Kılıç | 210 | 345s | 82 | 24s 34c |
| recipe_sword_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Kılıç | 210 | 345s | 82 | 24s 34c |
| recipe_hammer_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Dişbudak Kütüğü | 1× Rün Bağlı Savaş Çekici | 210 | 345s | 82 | 24s 34c |
| recipe_axe_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Balta | 235 | 375s | 97 | 28s 17c |
| recipe_dagger_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Hançer | 235 | 375s | 97 | 28s 17c |
| recipe_fist_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Muşta | 235 | 375s | 97 | 28s 17c |
| recipe_mace_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Topuz | 235 | 375s | 97 | 28s 17c |
| recipe_boots_phoenix | 24 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Çizme | 235 | 235s | 61 | 12s 67c |
| recipe_chest_phoenix | 24 | 4× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Göğüslük | 235 | 375s | 97 | 22s 53c |
| recipe_gloves_phoenix | 24 | 1× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Eldiven | 235 | 185s | 49 | 11s 27c |
| recipe_helm_phoenix | 24 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Kask | 235 | 235s | 61 | 14s 8c |
| recipe_legs_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Bacaklık | 235 | 330s | 85 | 19s 72c |
| recipe_rapier_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Meç | 235 | 375s | 97 | 28s 17c |
| recipe_ritualaxe_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Ayin Baltası | 235 | 375s | 97 | 28s 17c |
| recipe_shield_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Dağ Keçisi Postu | 1× Anka Kalkan | 235 | 330s | 85 | 20s 25c |
| recipe_shortsword_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Kısa Kılıç | 235 | 375s | 97 | 28s 17c |
| recipe_sickle_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Orak | 235 | 375s | 97 | 28s 17c |
| recipe_longsword_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Uzun Kılıç | 235 | 375s | 97 | 28s 17c |
| recipe_sword_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Kılıç | 235 | 375s | 97 | 28s 17c |
| recipe_hammer_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Akçaağaç Kütüğü | 1× Anka Savaş Çekici | 235 | 375s | 97 | 28s 17c |
| recipe_axe_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Balta | 260 | 395s | 109 | 32s 12c |
| recipe_dagger_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Hançer | 260 | 395s | 109 | 32s 12c |
| recipe_fist_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Muşta | 260 | 395s | 109 | 32s 12c |
| recipe_mace_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Topuz | 260 | 395s | 109 | 32s 12c |
| recipe_boots_dreadforged | 26 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Çizme | 260 | 245s | 68 | 14s 44c |
| recipe_chest_dreadforged | 26 | 4× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Göğüslük | 260 | 395s | 109 | 25s 70c |
| recipe_gloves_dreadforged | 26 | 1× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Eldiven | 260 | 195s | 55 | 12s 85c |
| recipe_helm_dreadforged | 26 | 2× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Kask | 260 | 245s | 68 | 16s 6c |
| recipe_legs_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Bacaklık | 260 | 345s | 96 | 22s 49c |
| recipe_rapier_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Meç | 260 | 395s | 109 | 32s 12c |
| recipe_ritualaxe_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Ayin Baltası | 260 | 395s | 109 | 32s 12c |
| recipe_shield_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Kalkan | 260 | 345s | 96 | 23s 10c |
| recipe_shortsword_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Kısa Kılıç | 260 | 395s | 109 | 32s 12c |
| recipe_sickle_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Orak | 260 | 395s | 109 | 32s 12c |
| recipe_longsword_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Uzun Kılıç | 260 | 395s | 109 | 32s 12c |
| recipe_sword_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Kılıç | 260 | 395s | 109 | 32s 12c |
| recipe_hammer_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Ceviz Kütüğü | 1× Dehşet Dövme Savaş Çekici | 260 | 395s | 109 | 32s 12c |
| recipe_axe_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Balta | 285 | 420s | 130 | 36s 20c |
| recipe_dagger_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Hançer | 285 | 420s | 130 | 36s 20c |
| recipe_fist_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Muşta | 285 | 420s | 130 | 36s 20c |
| recipe_mace_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Topuz | 285 | 420s | 130 | 36s 20c |
| recipe_boots_sovereign | 29 | 2× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Çizme | 285 | 265s | 81 | 16s 28c |
| recipe_chest_sovereign | 29 | 4× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Göğüslük | 285 | 420s | 130 | 28s 96c |
| recipe_gloves_sovereign | 29 | 1× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Eldiven | 285 | 210s | 65 | 14s 48c |
| recipe_helm_sovereign | 29 | 2× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Kask | 285 | 265s | 81 | 18s 10c |
| recipe_legs_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Bacaklık | 285 | 370s | 114 | 25s 34c |
| recipe_rapier_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Meç | 285 | 420s | 130 | 36s 20c |
| recipe_ritualaxe_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Ayin Baltası | 285 | 420s | 130 | 36s 20c |
| recipe_shield_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Hükümdar Kalkan | 285 | 370s | 114 | 26s 2c |
| recipe_shortsword_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Kısa Kılıç | 285 | 420s | 130 | 36s 20c |
| recipe_sickle_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Orak | 285 | 420s | 130 | 36s 20c |
| recipe_longsword_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Uzun Kılıç | 285 | 420s | 130 | 36s 20c |
| recipe_sword_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Kılıç | 285 | 420s | 130 | 36s 20c |
| recipe_hammer_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Hükümdar Savaş Çekici | 285 | 420s | 130 | 36s 20c |
| recipe_axe_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Balta | 310 | 440s | 146 | 40s 38c |
| recipe_dagger_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Hançer | 310 | 440s | 146 | 40s 38c |
| recipe_fist_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Muşta | 310 | 440s | 146 | 40s 38c |
| recipe_mace_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Topuz | 310 | 440s | 146 | 40s 38c |
| recipe_boots_coral | 31 | 2× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Çizme | 310 | 275s | 91 | 18s 16c |
| recipe_chest_coral | 31 | 4× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Göğüslük | 310 | 440s | 146 | 32s 30c |
| recipe_gloves_coral | 31 | 1× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Eldiven | 310 | 220s | 73 | 16s 15c |
| recipe_helm_coral | 31 | 2× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Kask | 310 | 275s | 91 | 20s 19c |
| recipe_legs_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Bacaklık | 310 | 385s | 128 | 28s 26c |
| recipe_rapier_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Meç | 310 | 440s | 146 | 40s 38c |
| recipe_ritualaxe_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Ayin Baltası | 310 | 440s | 146 | 40s 38c |
| recipe_shield_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 1× Yaban Mandası Derisi | 1× Mercan Kalkan | 310 | 385s | 128 | 29s 3c |
| recipe_shortsword_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Kısa Kılıç | 310 | 440s | 146 | 40s 38c |
| recipe_sickle_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Orak | 310 | 440s | 146 | 40s 38c |
| recipe_longsword_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Uzun Kılıç | 310 | 440s | 146 | 40s 38c |
| recipe_sword_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Kılıç | 310 | 440s | 146 | 40s 38c |
| recipe_hammer_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Maun Kütüğü | 1× Mercan Savaş Çekici | 310 | 440s | 146 | 40s 38c |
| recipe_axe_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Balta | 335 | 470s | 174 | 44s 66c |
| recipe_dagger_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Hançer | 335 | 470s | 174 | 44s 66c |
| recipe_fist_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Muşta | 335 | 470s | 174 | 44s 66c |
| recipe_mace_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Topuz | 335 | 470s | 174 | 44s 66c |
| recipe_boots_drakebone | 34 | 2× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Çizme | 335 | 295s | 109 | 20s 8c |
| recipe_chest_drakebone | 34 | 4× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Göğüslük | 335 | 470s | 174 | 35s 73c |
| recipe_gloves_drakebone | 34 | 1× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Eldiven | 335 | 235s | 87 | 17s 86c |
| recipe_helm_drakebone | 34 | 2× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Kask | 335 | 295s | 109 | 22s 33c |
| recipe_legs_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Bacaklık | 335 | 410s | 152 | 31s 26c |
| recipe_rapier_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Meç | 335 | 470s | 174 | 44s 66c |
| recipe_ritualaxe_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Ayin Baltası | 335 | 470s | 174 | 44s 66c |
| recipe_shield_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 1× Dev Geyik Postu | 1× Ejderha Kemiği Kalkan | 335 | 410s | 152 | 32s 11c |
| recipe_shortsword_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Kısa Kılıç | 335 | 470s | 174 | 44s 66c |
| recipe_sickle_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Orak | 335 | 470s | 174 | 44s 66c |
| recipe_longsword_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Uzun Kılıç | 335 | 470s | 174 | 44s 66c |
| recipe_sword_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Kılıç | 335 | 470s | 174 | 44s 66c |
| recipe_hammer_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Abanoz Kütüğü | 1× Ejderha Kemiği Savaş Çekici | 335 | 470s | 174 | 44s 66c |
| recipe_axe_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Balta | 360 | 490s | 196 | 49s 4c |
| recipe_dagger_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Hançer | 360 | 490s | 196 | 49s 4c |
| recipe_fist_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Muşta | 360 | 490s | 196 | 49s 4c |
| recipe_mace_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Topuz | 360 | 490s | 196 | 49s 4c |
| recipe_boots_mistwoven | 36 | 2× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Çizme | 360 | 305s | 122 | 22s 5c |
| recipe_chest_mistwoven | 36 | 4× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Göğüslük | 360 | 490s | 196 | 39s 23c |
| recipe_gloves_mistwoven | 36 | 1× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Eldiven | 360 | 245s | 98 | 19s 62c |
| recipe_helm_mistwoven | 36 | 2× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Kask | 360 | 305s | 122 | 24s 52c |
| recipe_legs_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Bacaklık | 360 | 430s | 171 | 34s 33c |
| recipe_rapier_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Meç | 360 | 490s | 196 | 49s 4c |
| recipe_ritualaxe_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Ayin Baltası | 360 | 490s | 196 | 49s 4c |
| recipe_shield_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 1× Dev Kurt Postu | 1× Sis Örgü Kalkan | 360 | 430s | 171 | 35s 26c |
| recipe_shortsword_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Kısa Kılıç | 360 | 490s | 196 | 49s 4c |
| recipe_sickle_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Orak | 360 | 490s | 196 | 49s 4c |
| recipe_longsword_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Uzun Kılıç | 360 | 490s | 196 | 49s 4c |
| recipe_sword_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Kılıç | 360 | 490s | 196 | 49s 4c |
| recipe_hammer_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Gümüş Söğüt Kütüğü | 1× Sis Örgü Savaş Çekici | 360 | 490s | 196 | 49s 4c |
| recipe_axe_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Balta | 385 | 520s | 233 | 53s 51c |
| recipe_dagger_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Hançer | 385 | 520s | 233 | 53s 51c |
| recipe_fist_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Muşta | 385 | 520s | 233 | 53s 51c |
| recipe_mace_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Topuz | 385 | 520s | 233 | 53s 51c |
| recipe_boots_prismatic | 39 | 2× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Çizme | 385 | 325s | 146 | 24s 6c |
| recipe_chest_prismatic | 39 | 4× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Göğüslük | 385 | 520s | 233 | 42s 81c |
| recipe_gloves_prismatic | 39 | 1× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Eldiven | 385 | 260s | 116 | 21s 41c |
| recipe_helm_prismatic | 39 | 2× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Kask | 385 | 325s | 146 | 26s 76c |
| recipe_legs_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Bacaklık | 385 | 455s | 204 | 37s 46c |
| recipe_rapier_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Meç | 385 | 520s | 233 | 53s 51c |
| recipe_ritualaxe_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Ayin Baltası | 385 | 520s | 233 | 53s 51c |
| recipe_shield_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 1× Dağ Aslanı Postu | 1× Prizmatik Kalkan | 385 | 455s | 204 | 38s 47c |
| recipe_shortsword_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Kısa Kılıç | 385 | 520s | 233 | 53s 51c |
| recipe_sickle_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Orak | 385 | 520s | 233 | 53s 51c |
| recipe_longsword_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Uzun Kılıç | 385 | 520s | 233 | 53s 51c |
| recipe_sword_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Kılıç | 385 | 520s | 233 | 53s 51c |
| recipe_hammer_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Altın Meşe Kütüğü | 1× Prizmatik Savaş Çekici | 385 | 520s | 233 | 53s 51c |
| recipe_axe_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Balta | 410 | 540s | 262 | 58s 8c |
| recipe_dagger_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Hançer | 410 | 540s | 262 | 58s 8c |
| recipe_fist_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Muşta | 410 | 540s | 262 | 58s 8c |
| recipe_mace_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Topuz | 410 | 540s | 262 | 58s 8c |
| recipe_boots_magmaforged | 41 | 2× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Çizme | 410 | 335s | 164 | 26s 12c |
| recipe_chest_magmaforged | 41 | 4× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Göğüslük | 410 | 540s | 262 | 46s 46c |
| recipe_gloves_magmaforged | 41 | 1× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Eldiven | 410 | 270s | 131 | 23s 23c |
| recipe_helm_magmaforged | 41 | 2× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Kask | 410 | 335s | 164 | 29s 4c |
| recipe_legs_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Bacaklık | 410 | 470s | 229 | 40s 65c |
| recipe_rapier_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Meç | 410 | 540s | 262 | 58s 8c |
| recipe_ritualaxe_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Ayin Baltası | 410 | 540s | 262 | 58s 8c |
| recipe_shield_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 1× Kara Panter Postu | 1× Magma Dövme Kalkan | 410 | 470s | 229 | 41s 75c |
| recipe_shortsword_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Kısa Kılıç | 410 | 540s | 262 | 58s 8c |
| recipe_sickle_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Orak | 410 | 540s | 262 | 58s 8c |
| recipe_longsword_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Uzun Kılıç | 410 | 540s | 262 | 58s 8c |
| recipe_sword_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Kılıç | 410 | 540s | 262 | 58s 8c |
| recipe_hammer_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Kara Odun Kütüğü | 1× Magma Dövme Savaş Çekici | 410 | 540s | 262 | 58s 8c |
| recipe_axe_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Balta | 435 | 565s | 312 | 62s 72c |
| recipe_dagger_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Hançer | 435 | 565s | 312 | 62s 72c |
| recipe_fist_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Muşta | 435 | 565s | 312 | 62s 72c |
| recipe_mace_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Topuz | 435 | 565s | 312 | 62s 72c |
| recipe_boots_tempest | 44 | 2× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Çizme | 435 | 355s | 195 | 28s 20c |
| recipe_chest_tempest | 44 | 4× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Göğüslük | 435 | 565s | 312 | 50s 18c |
| recipe_gloves_tempest | 44 | 1× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Eldiven | 435 | 285s | 156 | 25s 9c |
| recipe_helm_tempest | 44 | 2× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Kask | 435 | 355s | 195 | 31s 36c |
| recipe_legs_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Bacaklık | 435 | 495s | 273 | 43s 90c |
| recipe_rapier_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Meç | 435 | 565s | 312 | 62s 72c |
| recipe_ritualaxe_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Ayin Baltası | 435 | 565s | 312 | 62s 72c |
| recipe_shield_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Bora Kalkan | 435 | 495s | 273 | 45s 9c |
| recipe_shortsword_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Kısa Kılıç | 435 | 565s | 312 | 62s 72c |
| recipe_sickle_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Orak | 435 | 565s | 312 | 62s 72c |
| recipe_longsword_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Uzun Kılıç | 435 | 565s | 312 | 62s 72c |
| recipe_sword_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Kılıç | 435 | 565s | 312 | 62s 72c |
| recipe_hammer_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Bora Savaş Çekici | 435 | 565s | 312 | 62s 72c |
| recipe_axe_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Balta | 460 | 585s | 350 | 67s 45c |
| recipe_dagger_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Hançer | 460 | 585s | 350 | 67s 45c |
| recipe_fist_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Muşta | 460 | 585s | 350 | 67s 45c |
| recipe_mace_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Topuz | 460 | 585s | 350 | 67s 45c |
| recipe_boots_voidtouched | 46 | 2× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Çizme | 460 | 365s | 219 | 30s 33c |
| recipe_chest_voidtouched | 46 | 4× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Göğüslük | 460 | 585s | 350 | 53s 96c |
| recipe_gloves_voidtouched | 46 | 1× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Eldiven | 460 | 295s | 175 | 26s 98c |
| recipe_helm_voidtouched | 46 | 2× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Kask | 460 | 365s | 219 | 33s 72c |
| recipe_legs_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Bacaklık | 460 | 510s | 306 | 47s 21c |
| recipe_rapier_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Meç | 460 | 585s | 350 | 67s 45c |
| recipe_ritualaxe_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Ayin Baltası | 460 | 585s | 350 | 67s 45c |
| recipe_shield_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 1× Kara Panter Postu | 1× Hiçlik Kalkan | 460 | 510s | 306 | 48s 49c |
| recipe_shortsword_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Kısa Kılıç | 460 | 585s | 350 | 67s 45c |
| recipe_sickle_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Orak | 460 | 585s | 350 | 67s 45c |
| recipe_longsword_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Uzun Kılıç | 460 | 585s | 350 | 67s 45c |
| recipe_sword_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Kılıç | 460 | 585s | 350 | 67s 45c |
| recipe_hammer_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Kristal Huş Kütüğü | 1× Hiçlik Savaş Çekici | 460 | 585s | 350 | 67s 45c |
| recipe_axe_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Balta | 485 | 615s | 417 | 72s 25c |
| recipe_dagger_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Hançer | 485 | 615s | 417 | 72s 25c |
| recipe_fist_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Muşta | 485 | 615s | 417 | 72s 25c |
| recipe_mace_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Topuz | 485 | 615s | 417 | 72s 25c |
| recipe_boots_seraphic | 49 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Çizme | 485 | 385s | 261 | 32s 49c |
| recipe_chest_seraphic | 49 | 4× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Göğüslük | 485 | 615s | 417 | 57s 80c |
| recipe_gloves_seraphic | 49 | 1× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Eldiven | 485 | 305s | 209 | 28s 90c |
| recipe_helm_seraphic | 49 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Kask | 485 | 385s | 261 | 36s 12c |
| recipe_legs_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Bacaklık | 485 | 540s | 365 | 50s 58c |
| recipe_rapier_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Meç | 485 | 615s | 417 | 72s 25c |
| recipe_ritualaxe_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Ayin Baltası | 485 | 615s | 417 | 72s 25c |
| recipe_shield_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Melek Kalkan | 485 | 540s | 365 | 51s 94c |
| recipe_shortsword_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Kısa Kılıç | 485 | 615s | 417 | 72s 25c |
| recipe_sickle_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Orak | 485 | 615s | 417 | 72s 25c |
| recipe_longsword_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Uzun Kılıç | 485 | 615s | 417 | 72s 25c |
| recipe_sword_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Kılıç | 485 | 615s | 417 | 72s 25c |
| recipe_hammer_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Melek Savaş Çekici | 485 | 615s | 417 | 72s 25c |
| recipe_axe_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Balta | 510 | 635s | 469 | 77s 13c |
| recipe_dagger_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Hançer | 510 | 635s | 469 | 77s 13c |
| recipe_fist_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Muşta | 510 | 635s | 469 | 77s 13c |
| recipe_mace_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Topuz | 510 | 635s | 469 | 77s 13c |
| recipe_boots_eclipse | 51 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Çizme | 510 | 395s | 293 | 34s 68c |
| recipe_chest_eclipse | 51 | 4× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Göğüslük | 510 | 635s | 469 | 61s 70c |
| recipe_gloves_eclipse | 51 | 1× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Eldiven | 510 | 315s | 234 | 30s 85c |
| recipe_helm_eclipse | 51 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Kask | 510 | 395s | 293 | 38s 56c |
| recipe_legs_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Bacaklık | 510 | 555s | 410 | 53s 99c |
| recipe_rapier_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Meç | 510 | 635s | 469 | 77s 13c |
| recipe_ritualaxe_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Ayin Baltası | 510 | 635s | 469 | 77s 13c |
| recipe_shield_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 1× Kor Domuzu Postu | 1× Tutulma Kalkan | 510 | 555s | 410 | 55s 45c |
| recipe_shortsword_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Kısa Kılıç | 510 | 635s | 469 | 77s 13c |
| recipe_sickle_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Orak | 510 | 635s | 469 | 77s 13c |
| recipe_longsword_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Uzun Kılıç | 510 | 635s | 469 | 77s 13c |
| recipe_sword_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Kılıç | 510 | 635s | 469 | 77s 13c |
| recipe_hammer_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Korağacı Kütüğü | 1× Tutulma Savaş Çekici | 510 | 635s | 469 | 77s 13c |
| recipe_axe_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Balta | 535 | 660s | 558 | 82s 8c |
| recipe_dagger_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Hançer | 535 | 660s | 558 | 82s 8c |
| recipe_fist_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Muşta | 535 | 660s | 558 | 82s 8c |
| recipe_mace_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Topuz | 535 | 660s | 558 | 82s 8c |
| recipe_boots_wyrmlord | 54 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Çizme | 535 | 415s | 349 | 36s 91c |
| recipe_chest_wyrmlord | 54 | 4× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Göğüslük | 535 | 660s | 558 | 65s 66c |
| recipe_gloves_wyrmlord | 54 | 1× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Eldiven | 535 | 330s | 279 | 32s 83c |
| recipe_helm_wyrmlord | 54 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Kask | 535 | 415s | 349 | 41s 4c |
| recipe_legs_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Bacaklık | 535 | 580s | 488 | 57s 46c |
| recipe_rapier_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Meç | 535 | 660s | 558 | 82s 8c |
| recipe_ritualaxe_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Ayin Baltası | 535 | 660s | 558 | 82s 8c |
| recipe_shield_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Kalkan | 535 | 580s | 488 | 59s 1c |
| recipe_shortsword_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Kısa Kılıç | 535 | 660s | 558 | 82s 8c |
| recipe_sickle_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Orak | 535 | 660s | 558 | 82s 8c |
| recipe_longsword_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Uzun Kılıç | 535 | 660s | 558 | 82s 8c |
| recipe_sword_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Kılıç | 535 | 660s | 558 | 82s 8c |
| recipe_hammer_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Ejder Lordu Savaş Çekici | 535 | 660s | 558 | 82s 8c |
| recipe_axe_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Balta | 560 | 680s | 627 | 87s 10c |
| recipe_dagger_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Hançer | 560 | 680s | 627 | 87s 10c |
| recipe_fist_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Muşta | 560 | 680s | 627 | 87s 10c |
| recipe_mace_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Topuz | 560 | 680s | 627 | 87s 10c |
| recipe_boots_celestine | 56 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Çizme | 560 | 425s | 392 | 39s 17c |
| recipe_chest_celestine | 56 | 4× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Göğüslük | 560 | 680s | 627 | 69s 68c |
| recipe_gloves_celestine | 56 | 1× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Eldiven | 560 | 340s | 314 | 34s 84c |
| recipe_helm_celestine | 56 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Kask | 560 | 425s | 392 | 43s 55c |
| recipe_legs_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Bacaklık | 560 | 595s | 549 | 60s 97c |
| recipe_rapier_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Meç | 560 | 680s | 627 | 87s 10c |
| recipe_ritualaxe_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Ayin Baltası | 560 | 680s | 627 | 87s 10c |
| recipe_shield_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 1× Ayaz Kurdu Postu | 1× Göksel Kalkan | 560 | 595s | 549 | 62s 62c |
| recipe_shortsword_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Kısa Kılıç | 560 | 680s | 627 | 87s 10c |
| recipe_sickle_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Orak | 560 | 680s | 627 | 87s 10c |
| recipe_longsword_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Uzun Kılıç | 560 | 680s | 627 | 87s 10c |
| recipe_sword_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Kılıç | 560 | 680s | 627 | 87s 10c |
| recipe_hammer_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Ayaz Çamı Kütüğü | 1× Göksel Savaş Çekici | 560 | 680s | 627 | 87s 10c |
| recipe_axe_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Balta | 585 | 710s | 747 | 92s 19c |
| recipe_dagger_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Hançer | 585 | 710s | 747 | 92s 19c |
| recipe_fist_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Muşta | 585 | 710s | 747 | 92s 19c |
| recipe_mace_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Topuz | 585 | 710s | 747 | 92s 19c |
| recipe_boots_eternal | 59 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Çizme | 585 | 445s | 467 | 41s 45c |
| recipe_chest_eternal | 59 | 4× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Göğüslük | 585 | 710s | 747 | 73s 75c |
| recipe_gloves_eternal | 59 | 1× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Eldiven | 585 | 355s | 373 | 36s 88c |
| recipe_helm_eternal | 59 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Kask | 585 | 445s | 467 | 46s 9c |
| recipe_legs_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Bacaklık | 585 | 620s | 653 | 64s 53c |
| recipe_rapier_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Meç | 585 | 710s | 747 | 92s 19c |
| recipe_ritualaxe_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Ayin Baltası | 585 | 710s | 747 | 92s 19c |
| recipe_shield_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Ebedi Kalkan | 585 | 620s | 653 | 66s 28c |
| recipe_shortsword_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Kısa Kılıç | 585 | 710s | 747 | 92s 19c |
| recipe_sickle_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Orak | 585 | 710s | 747 | 92s 19c |
| recipe_longsword_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Uzun Kılıç | 585 | 710s | 747 | 92s 19c |
| recipe_sword_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Kılıç | 585 | 710s | 747 | 92s 19c |
| recipe_hammer_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Ebedi Savaş Çekici | 585 | 710s | 747 | 92s 19c |
| recipe_axe_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Balta | 610 | 730s | 839 | 97s 34c |
| recipe_dagger_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Hançer | 610 | 730s | 839 | 97s 34c |
| recipe_fist_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Muşta | 610 | 730s | 839 | 97s 34c |
| recipe_mace_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Topuz | 610 | 730s | 839 | 97s 34c |
| recipe_boots_starless | 61 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Çizme | 610 | 455s | 524 | 43s 77c |
| recipe_chest_starless | 61 | 4× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Göğüslük | 610 | 730s | 839 | 77s 87c |
| recipe_gloves_starless | 61 | 1× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Eldiven | 610 | 365s | 420 | 38s 94c |
| recipe_helm_starless | 61 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Kask | 610 | 455s | 524 | 48s 67c |
| recipe_legs_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Bacaklık | 610 | 640s | 734 | 68s 14c |
| recipe_rapier_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Meç | 610 | 730s | 839 | 97s 34c |
| recipe_ritualaxe_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Ayin Baltası | 610 | 730s | 839 | 97s 34c |
| recipe_shield_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 1× Dev Ayı Postu | 1× Yıldızsız Kalkan | 610 | 640s | 734 | 69s 98c |
| recipe_shortsword_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Kısa Kılıç | 610 | 730s | 839 | 97s 34c |
| recipe_sickle_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Orak | 610 | 730s | 839 | 97s 34c |
| recipe_longsword_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Uzun Kılıç | 610 | 730s | 839 | 97s 34c |
| recipe_sword_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Kılıç | 610 | 730s | 839 | 97s 34c |
| recipe_hammer_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Kanağacı Kütüğü | 1× Yıldızsız Savaş Çekici | 610 | 730s | 839 | 97s 34c |
| recipe_axe_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Balta | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_dagger_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Hançer | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_fist_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Muşta | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_mace_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Topuz | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_boots_nightsteel | 64 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Çizme | 635 | 475s | 625 | 46s 12c |
| recipe_chest_nightsteel | 64 | 4× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Göğüslük | 635 | 760s | 1000 | 82s 5c |
| recipe_gloves_nightsteel | 64 | 1× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Eldiven | 635 | 380s | 500 | 41s 2c |
| recipe_helm_nightsteel | 64 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Kask | 635 | 475s | 625 | 51s 28c |
| recipe_legs_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Bacaklık | 635 | 665s | 875 | 71s 79c |
| recipe_rapier_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Meç | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_ritualaxe_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Ayin Baltası | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_shield_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Kalkan | 635 | 665s | 875 | 73s 74c |
| recipe_shortsword_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Kısa Kılıç | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_sickle_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Orak | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_longsword_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Uzun Kılıç | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_sword_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Kılıç | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_hammer_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Gece Çeliği Savaş Çekici | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_axe_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Balta | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_dagger_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Hançer | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_fist_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Muşta | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_mace_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Topuz | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_boots_chronoforged | 66 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Çizme | 660 | 485s | 702 | 48s 49c |
| recipe_chest_chronoforged | 66 | 4× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Göğüslük | 660 | 780s | 1123 | 86s 27c |
| recipe_gloves_chronoforged | 66 | 1× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Eldiven | 660 | 390s | 562 | 43s 14c |
| recipe_helm_chronoforged | 66 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Kask | 660 | 485s | 702 | 53s 92c |
| recipe_legs_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Bacaklık | 660 | 680s | 983 | 75s 49c |
| recipe_rapier_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Meç | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_ritualaxe_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Ayin Baltası | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_shield_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Kalkan | 660 | 680s | 983 | 77s 53c |
| recipe_shortsword_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Kısa Kılıç | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_sickle_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Orak | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_longsword_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Uzun Kılıç | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_sword_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Kılıç | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_hammer_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Ruhağacı Kütüğü | 1× Zaman Dövme Savaş Çekici | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_axe_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Balta | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_dagger_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Hançer | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_fist_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Muşta | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_mace_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Topuz | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_boots_ascendant | 69 | 2× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Çizme | 685 | 505s | 836 | 50s 89c |
| recipe_chest_ascendant | 69 | 4× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Göğüslük | 685 | 805s | 1338 | 90s 54c |
| recipe_gloves_ascendant | 69 | 1× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Eldiven | 685 | 405s | 669 | 45s 27c |
| recipe_helm_ascendant | 69 | 2× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Kask | 685 | 505s | 836 | 56s 59c |
| recipe_legs_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Bacaklık | 685 | 705s | 1170 | 79s 23c |
| recipe_rapier_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Meç | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_ritualaxe_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Ayin Baltası | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_shield_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Yükselmiş Kalkan | 685 | 705s | 1170 | 81s 37c |
| recipe_shortsword_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Kısa Kılıç | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_sickle_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Orak | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_longsword_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Uzun Kılıç | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_sword_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Kılıç | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_hammer_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Yükselmiş Savaş Çekici | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_axe_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Balta | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_dagger_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Hançer | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_fist_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Muşta | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_mace_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Topuz | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_boots_stormborn | 71 | 2× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Çizme | 710 | 515s | 939 | 53s 32c |
| recipe_chest_stormborn | 71 | 4× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Göğüslük | 710 | 825s | 1503 | 94s 86c |
| recipe_gloves_stormborn | 71 | 1× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Eldiven | 710 | 415s | 751 | 47s 43c |
| recipe_helm_stormborn | 71 | 2× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Kask | 710 | 515s | 939 | 59s 29c |
| recipe_legs_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Bacaklık | 710 | 720s | 1315 | 83s |
| recipe_rapier_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Meç | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_ritualaxe_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Ayin Baltası | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_shield_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Kalkan | 710 | 720s | 1315 | 85s 25c |
| recipe_shortsword_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Kısa Kılıç | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_sickle_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Orak | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_longsword_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Uzun Kılıç | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_sword_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Kılıç | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_hammer_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Kadim Demirağaç Kütüğü | 1× Fırtına Doğumlu Savaş Çekici | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_axe_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Balta | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_dagger_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Hançer | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_fist_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Muşta | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_mace_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Topuz | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_boots_riftforged | 74 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Çizme | 735 | 535s | 1119 | 55s 78c |
| recipe_chest_riftforged | 74 | 4× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Göğüslük | 735 | 855s | 1790 | 99s 23c |
| recipe_gloves_riftforged | 74 | 1× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Eldiven | 735 | 425s | 895 | 49s 61c |
| recipe_helm_riftforged | 74 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Kask | 735 | 535s | 1119 | 62s 2c |
| recipe_legs_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Bacaklık | 735 | 750s | 1566 | 86s 82c |
| recipe_rapier_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Meç | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_ritualaxe_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Ayin Baltası | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_shield_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Yarık Dövme Kalkan | 735 | 750s | 1566 | 89s 18c |
| recipe_shortsword_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Kısa Kılıç | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_sickle_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Orak | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_longsword_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Uzun Kılıç | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_sword_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Kılıç | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_hammer_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Yarık Dövme Savaş Çekici | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_axe_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Balta | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_dagger_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Hançer | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_fist_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Muşta | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_mace_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Topuz | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_boots_chaos | 76 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Çizme | 760 | 545s | 1257 | 58s 25c |
| recipe_chest_chaos | 76 | 4× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Göğüslük | 760 | 875s | 2011 | 1g 3s 64c |
| recipe_gloves_chaos | 76 | 1× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Eldiven | 760 | 435s | 1006 | 51s 82c |
| recipe_helm_chaos | 76 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Kask | 760 | 545s | 1257 | 64s 77c |
| recipe_legs_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Bacaklık | 760 | 765s | 1760 | 90s 68c |
| recipe_rapier_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Meç | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_ritualaxe_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Ayin Baltası | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_shield_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 1× Gölge Kimera Derisi | 1× Kaos Kalkan | 760 | 765s | 1760 | 93s 14c |
| recipe_shortsword_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Kısa Kılıç | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_sickle_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Orak | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_longsword_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Uzun Kılıç | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_sword_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Kılıç | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_hammer_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Yıldızkabuğu Kütüğü | 1× Kaos Savaş Çekici | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_axe_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Balta | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_dagger_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Hançer | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_fist_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Muşta | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_mace_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Topuz | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_boots_skyforged | 79 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Çizme | 785 | 565s | 1497 | 60s 76c |
| recipe_chest_skyforged | 79 | 4× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Göğüslük | 785 | 900s | 2395 | 1g 8s 9c |
| recipe_gloves_skyforged | 79 | 1× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Eldiven | 785 | 450s | 1198 | 54s 5c |
| recipe_helm_skyforged | 79 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Kask | 785 | 565s | 1497 | 67s 56c |
| recipe_legs_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Bacaklık | 785 | 790s | 2096 | 94s 58c |
| recipe_rapier_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Meç | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_ritualaxe_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Ayin Baltası | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_shield_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Kalkan | 785 | 790s | 2096 | 97s 14c |
| recipe_shortsword_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Kısa Kılıç | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_sickle_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Orak | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_longsword_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Uzun Kılıç | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_sword_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Kılıç | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_hammer_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Gök Dövme Savaş Çekici | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_axe_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Balta | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_dagger_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Hançer | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_fist_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Muşta | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_mace_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Topuz | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_boots_archangel | 81 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Çizme | 810 | 575s | 1682 | 63s 28c |
| recipe_chest_archangel | 81 | 4× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Göğüslük | 810 | 920s | 2692 | 1g 12s 59c |
| recipe_gloves_archangel | 81 | 1× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Eldiven | 810 | 460s | 1346 | 56s 29c |
| recipe_helm_archangel | 81 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Kask | 810 | 575s | 1682 | 70s 37c |
| recipe_legs_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Bacaklık | 810 | 805s | 2355 | 98s 51c |
| recipe_rapier_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Meç | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_ritualaxe_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Ayin Baltası | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_shield_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Başmelek Kalkan | 810 | 805s | 2355 | 1g 1s 18c |
| recipe_shortsword_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Kısa Kılıç | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_sickle_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Orak | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_longsword_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Uzun Kılıç | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_sword_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Kılıç | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_hammer_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Ayağacı Kütüğü | 1× Başmelek Savaş Çekici | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_axe_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Balta | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_dagger_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Hançer | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_fist_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Muşta | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_mace_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Topuz | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_boots_aurora | 84 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Çizme | 835 | 595s | 2004 | 65s 84c |
| recipe_chest_aurora | 84 | 4× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Göğüslük | 835 | 950s | 3206 | 1g 17s 13c |
| recipe_gloves_aurora | 84 | 1× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Eldiven | 835 | 475s | 1603 | 58s 56c |
| recipe_helm_aurora | 84 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Kask | 835 | 595s | 2004 | 73s 20c |
| recipe_legs_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Bacaklık | 835 | 830s | 2805 | 1g 2s 48c |
| recipe_rapier_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Meç | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_ritualaxe_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Ayin Baltası | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_shield_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Kalkan | 835 | 830s | 2805 | 1g 5s 26c |
| recipe_shortsword_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Kısa Kılıç | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_sickle_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Orak | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_longsword_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Uzun Kılıç | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_sword_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Kılıç | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_hammer_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Şafak Savaş Çekici | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_axe_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Balta | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_dagger_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Hançer | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_fist_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Muşta | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_mace_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Topuz | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_boots_ashen | 86 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Çizme | 860 | 605s | 2251 | 68s 41c |
| recipe_chest_ashen | 86 | 4× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Göğüslük | 860 | 970s | 3602 | 1g 21s 70c |
| recipe_gloves_ashen | 86 | 1× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Eldiven | 860 | 485s | 1801 | 60s 85c |
| recipe_helm_ashen | 86 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Kask | 860 | 605s | 2251 | 76s 6c |
| recipe_legs_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Bacaklık | 860 | 850s | 3151 | 1g 6s 49c |
| recipe_rapier_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Meç | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_ritualaxe_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Ayin Baltası | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_shield_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 1× Yaşlı Wyvern Pulu | 1× Kül Kalkan | 860 | 850s | 3151 | 1g 9s 38c |
| recipe_shortsword_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Kısa Kılıç | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_sickle_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Orak | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_longsword_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Uzun Kılıç | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_sword_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Kılıç | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_hammer_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Ejderyüreği Kütüğü | 1× Kül Savaş Çekici | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_axe_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Balta | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_dagger_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Hançer | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_fist_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Muşta | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_mace_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Topuz | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_boots_leviathan | 89 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Çizme | 885 | 625s | 2681 | 71s 1c |
| recipe_chest_leviathan | 89 | 4× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Göğüslük | 885 | 1000s | 4290 | 1g 26s 32c |
| recipe_gloves_leviathan | 89 | 1× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Eldiven | 885 | 500s | 2145 | 63s 16c |
| recipe_helm_leviathan | 89 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Kask | 885 | 625s | 2681 | 78s 95c |
| recipe_legs_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Bacaklık | 885 | 875s | 3753 | 1g 10s 53c |
| recipe_rapier_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Meç | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_ritualaxe_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Ayin Baltası | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_shield_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Leviathan Kalkan | 885 | 875s | 3753 | 1g 13s 53c |
| recipe_shortsword_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Kısa Kılıç | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_sickle_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Orak | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_longsword_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Uzun Kılıç | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_sword_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Kılıç | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_hammer_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Leviathan Savaş Çekici | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_axe_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Balta | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_dagger_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Hançer | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_fist_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Muşta | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_mace_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Topuz | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_boots_doomforged | 91 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Çizme | 910 | 635s | 3012 | 73s 62c |
| recipe_chest_doomforged | 91 | 4× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Göğüslük | 910 | 1020s | 4820 | 1g 30s 98c |
| recipe_gloves_doomforged | 91 | 1× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Eldiven | 910 | 510s | 2410 | 65s 49c |
| recipe_helm_doomforged | 91 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Kask | 910 | 635s | 3012 | 81s 86c |
| recipe_legs_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Bacaklık | 910 | 890s | 4218 | 1g 14s 61c |
| recipe_rapier_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Meç | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_ritualaxe_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Ayin Baltası | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_shield_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Kalkan | 910 | 890s | 4218 | 1g 17s 71c |
| recipe_shortsword_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Kısa Kılıç | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_sickle_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Orak | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_longsword_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Uzun Kılıç | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_sword_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Kılıç | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_hammer_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Boşlukağacı Kütüğü | 1× Kıyamet Dövme Savaş Çekici | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_axe_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Balta | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_dagger_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Hançer | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_fist_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Muşta | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_mace_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Topuz | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_boots_voidforged | 94 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Çizme | 935 | 655s | 3588 | 76s 26c |
| recipe_chest_voidforged | 94 | 4× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Göğüslük | 935 | 1045s | 5741 | 1g 35s 68c |
| recipe_gloves_voidforged | 94 | 1× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Eldiven | 935 | 525s | 2870 | 67s 84c |
| recipe_helm_voidforged | 94 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Kask | 935 | 655s | 3588 | 84s 80c |
| recipe_legs_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Bacaklık | 935 | 915s | 5023 | 1g 18s 72c |
| recipe_rapier_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Meç | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_ritualaxe_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Ayin Baltası | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_shield_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Kalkan | 935 | 915s | 5023 | 1g 21s 93c |
| recipe_shortsword_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Kısa Kılıç | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_sickle_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Orak | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_longsword_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Uzun Kılıç | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_sword_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Kılıç | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_hammer_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Hiçlik Dövme Savaş Çekici | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_axe_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Balta | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_dagger_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Hançer | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_fist_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Muşta | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_mace_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Topuz | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_boots_godslayer | 96 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Çizme | 960 | 665s | 4032 | 78s 93c |
| recipe_chest_godslayer | 96 | 4× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Göğüslük | 960 | 1065s | 6450 | 1g 40s 42c |
| recipe_gloves_godslayer | 96 | 1× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Eldiven | 960 | 535s | 3225 | 70s 21c |
| recipe_helm_godslayer | 96 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Kask | 960 | 665s | 4032 | 87s 76c |
| recipe_legs_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Bacaklık | 960 | 930s | 5644 | 1g 22s 86c |
| recipe_rapier_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Meç | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_ritualaxe_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Ayin Baltası | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_shield_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Kalkan | 960 | 930s | 5644 | 1g 26s 19c |
| recipe_shortsword_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Kısa Kılıç | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_sickle_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Orak | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_longsword_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Uzun Kılıç | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_sword_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Kılıç | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_hammer_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Göksel Dişbudak Kütüğü | 1× Tanrı Katili Savaş Çekici | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_axe_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Balta | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_dagger_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Hançer | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_fist_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Muşta | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_mace_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Topuz | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_boots_primordial | 99 | 2× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Çizme | 985 | 685s | 4802 | 81s 61c |
| recipe_chest_primordial | 99 | 4× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Göğüslük | 985 | 1095s | 7682 | 1g 45s 19c |
| recipe_gloves_primordial | 99 | 1× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Eldiven | 985 | 545s | 3841 | 72s 59c |
| recipe_helm_primordial | 99 | 2× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Kask | 985 | 685s | 4802 | 90s 74c |
| recipe_legs_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Bacaklık | 985 | 960s | 6722 | 1g 27s 4c |
| recipe_rapier_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Meç | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_ritualaxe_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Ayin Baltası | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_shield_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 1× Göksel Ejder Pulu | 1× Kadim Çağ Kalkan | 985 | 960s | 6722 | 1g 30s 48c |
| recipe_shortsword_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Kısa Kılıç | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_sickle_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Orak | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_longsword_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Uzun Kılıç | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_sword_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Kılıç | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_hammer_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Dünya Ağacı Kütüğü | 1× Kadim Çağ Savaş Çekici | 985 | 1095s | 7682 | 1g 81s 48c |

### Aşçılık — 30 tarif, beceri 1–90 (tavan 100), çıktı Sv 1–1000

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_cook_grilled_guppy | 1 | 2× Minik Lepistes, 1× Sıradan Ot, 1× Fidan Dalı | 2× Izgara Lepistes | 1 | 20s | 40 | 2c |
| recipe_cook_sardine_skewer | 4 | 2× Nehir Sardalyası, 1× Karahindiba, 1× Kavak Kütüğü | 2× Sardalya Şiş | 40 | 30s | 40 | 2c |
| recipe_cook_minnow_soup | 6 | 2× Golyan, 1× Yabani Nane, 1× Söğüt Kütüğü | 2× Naneli Balık Çorbası | 70 | 30s | 40 | 2c |
| recipe_cook_pike_fillet | 9 | 2× Küçük Turna, 1× Papatya, 1× Köknar Kütüğü | 2× Turna Fileto | 100 | 40s | 60 | 5c |
| recipe_cook_lavender_carp | 11 | 2× Sazan, 1× Lavanta, 1× Çam Kütüğü | 2× Lavantalı Sazan | 120 | 40s | 60 | 5c |
| recipe_cook_grilled_silver_herring | 13 | 2× Gümüş Ringa, 1× Ayçiçeği, 1× Huş Kütüğü | 2× Izgara Gümüş Ringa | 140 | 45s | 60 | 5c |
| recipe_cook_herb_trout | 14 | 2× Dere Alabalığı, 1× Acıkök, 1× Sedir Kütüğü | 2× Otlu Alabalık | 160 | 50s | 60 | 5c |
| recipe_cook_smoked_salmon | 16 | 2× Somon, 1× Şifalı Yosun, 1× Meşe Kütüğü | 2× Füme Somon | 180 | 50s | 60 | 5c |
| recipe_cook_bass_stew | 18 | 2× Tatlı Su Levreği, 1× Bataklık Sazı, 1× Kayın Kütüğü | 2× Levrek Güveci | 200 | 55s | 90 | 10c |
| recipe_cook_grilled_lake_perch | 20 | 2× Göl Levreği, 1× Mavi Devedikeni, 1× Dişbudak Kütüğü | 2× Izgara Göl Levreği | 220 | 60s | 90 | 10c |
| recipe_cook_grilled_red_snapper | 22 | 2× Kırmızı Mercan, 1× Kızıl Adaçayı, 1× Akçaağaç Kütüğü | 2× Izgara Kırmızı Mercan | 240 | 65s | 90 | 10c |
| recipe_cook_grilled_river_catfish | 23 | 2× Nehir Yayını, 1× Güzelavrat Otu, 1× Ceviz Kütüğü | 2× Izgara Nehir Yayını | 260 | 65s | 90 | 10c |
| recipe_cook_tuna_steak | 25 | 2× Mavi Yüzgeçli Ton, 1× Yaban Ginsengi, 1× Tik Kütüğü | 2× Ginsengli Ton Balığı | 280 | 70s | 90 | 10c |
| recipe_cook_grilled_rainbow_trout | 27 | 2× Gökkuşağı Alabalığı, 1× Altın Yonca, 1× Maun Kütüğü | 2× Izgara Gökkuşağı Alabalığı | 300 | 75s | 180 | 20c |
| recipe_cook_royal_carp_feast | 30 | 2× Kraliyet Sazanı, 1× Düş Kökü, 1× Abanoz Kütüğü | 2× Kraliyet Sazanı Ziyafeti | 330 | 80s | 180 | 20c |
| recipe_cook_grilled_emerald_pike | 32 | 2× Zümrüt Turnası, 1× Gümüş Yaprak, 1× Gümüş Söğüt Kütüğü | 2× Izgara Zümrüt Turnası | 360 | 85s | 180 | 20c |
| recipe_cook_grilled_silver_barracuda | 35 | 2× Gümüş Barakuda, 1× Kan Dikeni, 1× Altın Meşe Kütüğü | 2× Izgara Gümüş Barakuda | 390 | 90s | 180 | 20c |
| recipe_cook_sturgeon_pie | 38 | 2× Altın Mersin, 1× Hayalet Mantarı, 1× Kara Odun Kütüğü | 2× Mantarlı Mersin Turtası | 420 | 95s | 180 | 20c |
| recipe_cook_swordfish_grill | 40 | 2× Kızıl Kılıçbalığı, 1× Fırtına Çiçeği, 1× Kristal Huş Kütüğü | 2× Fırtına Izgarası Kılıçbalığı | 450 | 100s | 260 | 40c |
| recipe_cook_grilled_stormscale_eel | 45 | 2× Fırtına Pullu Yılanbalığı, 1× Kor Yaprağı, 1× Korağacı Kütüğü | 2× Izgara Fırtına Pullu Yılanbalığı | 500 | 110s | 260 | 40c |
| recipe_cook_frostfin_chowder | 50 | 2× Ayaz Yüzgeçli Somon, 1× Ayaz Çiçeği, 1× Ayaz Çamı Kütüğü | 2× Ayaz Somon Çorbası | 550 | 120s | 260 | 40c |
| recipe_cook_grilled_obsidian_tuna | 54 | 2× Obsidyen Ton, 1× Ay Gölgesi, 1× Kanağacı Kütüğü | 2× Izgara Obsidyen Ton | 600 | 130s | 460 | 80c |
| recipe_cook_marlin_platter | 58 | 2× Hayalet Marlin, 1× Ruh Nilüferi, 1× Ruhağacı Kütüğü | 2× Ruh Marlini Tabağı | 650 | 135s | 460 | 80c |
| recipe_cook_grilled_golden_carp | 63 | 2× İmparatorluk Altın Koisi, 1× Kadim Adamotu, 1× Kadim Demirağaç Kütüğü | 2× Izgara İmparatorluk Altın Koisi | 700 | 145s | 460 | 80c |
| recipe_cook_dragonblood_serpentfish | 68 | 2× Kristal Yılanbalığı, 1× Ejderkanı Otu, 1× Yıldızkabuğu Kütüğü | 2× Ejderkanlı Yılanbalığı | 750 | 155s | 760 | 1s 50c |
| recipe_cook_grilled_void_angler | 72 | 2× Boşluk Fenerbalığı, 1× Yıldız Çiçeği, 1× Ayağacı Kütüğü | 2× Izgara Boşluk Fenerbalığı | 800 | 165s | 760 | 1s 50c |
| recipe_cook_moonfish_delicacy | 76 | 2× Yıldız Aybalığı, 1× Boşluk Orkidesi, 1× Ejderyüreği Kütüğü | 2× Yıldız Aybalığı Lezzeti | 850 | 170s | 760 | 1s 50c |
| recipe_cook_phoenix_koi_feast | 81 | 2× Anka Koisi, 1× Göksel Nilüfer, 1× Boşlukağacı Kütüğü | 2× Anka Koi Ziyafeti | 900 | 180s | 1100 | 2s 50c |
| recipe_cook_grilled_abyssal_dragonfish | 86 | 2× Uçurum Ejderbalığı, 1× Anka Çiçeği, 1× Göksel Dişbudak Kütüğü | 2× Izgara Uçurum Ejderbalığı | 950 | 190s | 1100 | 2s 50c |
| recipe_cook_leviathan_banquet | 90 | 2× Göksel Leviathan, 1× Dünya Kökü, 1× Dünya Ağacı Kütüğü | 2× Leviathan Kraliyet Şöleni | 1000 | 200s | 1100 | 2s 50c |

### Kuyumculuk — 126 tarif, beceri 3–99 (tavan 100), çıktı Sv 25–985

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_necklace_runed | 3 | 4× Alüminyum Külçesi, 2× Taş Cevheri | 1× Rünlü Kolye | 25 | 130s | 21 | 70c |
| recipe_ring_runed | 3 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Rünlü Yüzük | 25 | 110s | 18 | 60c |
| recipe_charm_runed | 3 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Rünlü Tılsım | 25 | 110s | 18 | 60c |
| recipe_necklace_mithril | 4 | 4× Alüminyum Külçesi, 2× Taş Cevheri | 1× Mithril Kolye | 35 | 135s | 23 | 1s 8c |
| recipe_ring_mithril | 4 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Mithril Yüzük | 35 | 115s | 19 | 93c |
| recipe_charm_mithril | 4 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Mithril Tılsım | 35 | 115s | 19 | 93c |
| recipe_necklace_adamant | 5 | 4× Alüminyum Külçesi, 2× Taş Cevheri | 1× Adamant Kolye | 50 | 145s | 24 | 1s 72c |
| recipe_ring_adamant | 5 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Adamant Yüzük | 50 | 120s | 20 | 1s 48c |
| recipe_charm_adamant | 5 | 3× Alüminyum Külçesi, 2× Taş Cevheri | 1× Adamant Tılsım | 50 | 120s | 20 | 1s 48c |
| recipe_necklace_wyrmforged | 7 | 4× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Ejder Dövme Kolye | 65 | 160s | 27 | 2s 42c |
| recipe_ring_wyrmforged | 7 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Ejder Dövme Yüzük | 65 | 130s | 22 | 2s 8c |
| recipe_charm_wyrmforged | 7 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Ejder Dövme Tılsım | 65 | 130s | 22 | 2s 8c |
| recipe_necklace_starforged | 8 | 4× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Yıldız Dövme Kolye | 80 | 165s | 29 | 3s 18c |
| recipe_ring_starforged | 8 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Yıldız Dövme Yüzük | 80 | 140s | 24 | 2s 72c |
| recipe_charm_starforged | 8 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Yıldız Dövme Tılsım | 80 | 140s | 24 | 2s 72c |
| recipe_necklace_titan | 10 | 4× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Titan Kolye | 95 | 180s | 32 | 3s 97c |
| recipe_ring_titan | 10 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Titan Yüzük | 95 | 150s | 27 | 3s 40c |
| recipe_charm_titan | 10 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Titan Tılsım | 95 | 150s | 27 | 3s 40c |
| recipe_necklace_royal | 11 | 4× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Kraliyet Kolye | 110 | 185s | 34 | 4s 80c |
| recipe_ring_royal | 11 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Kraliyet Yüzük | 110 | 155s | 28 | 4s 12c |
| recipe_charm_royal | 11 | 3× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Kraliyet Tılsım | 110 | 155s | 28 | 4s 12c |
| recipe_necklace_sunforged | 14 | 3× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Güneş Dövme Kolye | 135 | 210s | 41 | 6s 27c |
| recipe_ring_sunforged | 14 | 2× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Güneş Dövme Yüzük | 135 | 175s | 34 | 5s 37c |
| recipe_charm_sunforged | 14 | 2× Kurşun Külçesi, 1× Alüminyum Külçesi, 2× Çakmaktaşı | 1× Güneş Dövme Tılsım | 135 | 175s | 34 | 5s 37c |
| recipe_necklace_obsidian | 16 | 3× Demir Külçesi, 1× Kurşun Külçesi, 2× Çakmaktaşı | 1× Obsidyen Kolye | 160 | 225s | 46 | 7s 82c |
| recipe_ring_obsidian | 16 | 2× Demir Külçesi, 1× Kurşun Külçesi, 2× Çakmaktaşı | 1× Obsidyen Yüzük | 160 | 185s | 38 | 6s 70c |
| recipe_charm_obsidian | 16 | 2× Demir Külçesi, 1× Kurşun Külçesi, 2× Çakmaktaşı | 1× Obsidyen Tılsım | 160 | 185s | 38 | 6s 70c |
| recipe_necklace_sandsteel | 19 | 3× Bakır Külçesi, 1× Kalay Külçesi, 2× Çakmaktaşı | 1× Kum Çeliği Kolye | 185 | 245s | 55 | 9s 44c |
| recipe_ring_sandsteel | 19 | 2× Bakır Külçesi, 1× Kalay Külçesi, 2× Çakmaktaşı | 1× Kum Çeliği Yüzük | 185 | 205s | 46 | 8s 9c |
| recipe_charm_sandsteel | 19 | 2× Bakır Külçesi, 1× Kalay Külçesi, 2× Çakmaktaşı | 1× Kum Çeliği Tılsım | 185 | 205s | 46 | 8s 9c |
| recipe_necklace_runebound | 21 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Çakmaktaşı | 1× Rün Bağlı Kolye | 210 | 260s | 61 | 11s 13c |
| recipe_ring_runebound | 21 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Çakmaktaşı | 1× Rün Bağlı Yüzük | 210 | 215s | 51 | 9s 54c |
| recipe_charm_runebound | 21 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Çakmaktaşı | 1× Rün Bağlı Tılsım | 210 | 215s | 51 | 9s 54c |
| recipe_necklace_phoenix | 24 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Anka Kolye | 235 | 280s | 73 | 12s 89c |
| recipe_ring_phoenix | 24 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Anka Yüzük | 235 | 235s | 61 | 11s 5c |
| recipe_charm_phoenix | 24 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Anka Tılsım | 235 | 235s | 61 | 11s 5c |
| recipe_necklace_dreadforged | 26 | 3× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Dehşet Dövme Kolye | 260 | 295s | 82 | 14s 70c |
| recipe_ring_dreadforged | 26 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Dehşet Dövme Yüzük | 260 | 245s | 68 | 12s 60c |
| recipe_charm_dreadforged | 26 | 2× Bronz Külçesi, 1× Bakır Külçesi, 2× Obsidyen | 1× Dehşet Dövme Tılsım | 260 | 245s | 68 | 12s 60c |
| recipe_necklace_sovereign | 29 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Hükümdar Kolye | 285 | 315s | 98 | 16s 56c |
| recipe_ring_sovereign | 29 | 2× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Hükümdar Yüzük | 285 | 265s | 81 | 14s 20c |
| recipe_charm_sovereign | 29 | 2× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Hükümdar Tılsım | 285 | 265s | 81 | 14s 20c |
| recipe_necklace_coral | 31 | 3× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Mercan Kolye | 310 | 330s | 110 | 18s 47c |
| recipe_ring_coral | 31 | 2× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Mercan Yüzük | 310 | 275s | 91 | 15s 83c |
| recipe_charm_coral | 31 | 2× Nikel Külçesi, 1× Çinko Külçesi, 2× Obsidyen | 1× Mercan Tılsım | 310 | 275s | 91 | 15s 83c |
| recipe_necklace_drakebone | 34 | 3× Krom Külçesi, 1× Nikel Külçesi, 2× Obsidyen | 1× Ejderha Kemiği Kolye | 335 | 355s | 130 | 20s 43c |
| recipe_ring_drakebone | 34 | 2× Krom Külçesi, 1× Nikel Külçesi, 2× Obsidyen | 1× Ejderha Kemiği Yüzük | 335 | 295s | 109 | 17s 51c |
| recipe_charm_drakebone | 34 | 2× Krom Külçesi, 1× Nikel Külçesi, 2× Obsidyen | 1× Ejderha Kemiği Tılsım | 335 | 295s | 109 | 17s 51c |
| recipe_necklace_mistwoven | 36 | 3× Kobalt Külçesi, 1× Krom Külçesi, 2× Obsidyen | 1× Sis Örgü Kolye | 360 | 365s | 147 | 22s 44c |
| recipe_ring_mistwoven | 36 | 2× Kobalt Külçesi, 1× Krom Külçesi, 2× Obsidyen | 1× Sis Örgü Yüzük | 360 | 305s | 122 | 19s 23c |
| recipe_charm_mistwoven | 36 | 2× Kobalt Külçesi, 1× Krom Külçesi, 2× Obsidyen | 1× Sis Örgü Tılsım | 360 | 305s | 122 | 19s 23c |
| recipe_necklace_prismatic | 39 | 3× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Obsidyen | 1× Prizmatik Kolye | 385 | 390s | 175 | 24s 48c |
| recipe_ring_prismatic | 39 | 2× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Obsidyen | 1× Prizmatik Yüzük | 385 | 325s | 146 | 20s 99c |
| recipe_charm_prismatic | 39 | 2× Gümüş Külçesi, 1× Kobalt Külçesi, 2× Obsidyen | 1× Prizmatik Tılsım | 385 | 325s | 146 | 20s 99c |
| recipe_necklace_magmaforged | 41 | 3× Altın Külçesi, 1× Gümüş Külçesi, 2× Obsidyen | 1× Magma Dövme Kolye | 410 | 405s | 196 | 26s 57c |
| recipe_ring_magmaforged | 41 | 2× Altın Külçesi, 1× Gümüş Külçesi, 2× Obsidyen | 1× Magma Dövme Yüzük | 410 | 335s | 164 | 22s 77c |
| recipe_charm_magmaforged | 41 | 2× Altın Külçesi, 1× Gümüş Külçesi, 2× Obsidyen | 1× Magma Dövme Tılsım | 410 | 335s | 164 | 22s 77c |
| recipe_necklace_tempest | 44 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Bora Kolye | 435 | 425s | 234 | 28s 70c |
| recipe_ring_tempest | 44 | 2× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Bora Yüzük | 435 | 355s | 195 | 24s 60c |
| recipe_charm_tempest | 44 | 2× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Bora Tılsım | 435 | 355s | 195 | 24s 60c |
| recipe_necklace_voidtouched | 46 | 3× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Hiçlik Kolye | 460 | 440s | 262 | 30s 86c |
| recipe_ring_voidtouched | 46 | 2× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Hiçlik Yüzük | 460 | 365s | 219 | 26s 45c |
| recipe_charm_voidtouched | 46 | 2× Elektrum Külçesi, 1× Altın Külçesi, 2× Obsidyen | 1× Hiçlik Tılsım | 460 | 365s | 219 | 26s 45c |
| recipe_necklace_seraphic | 49 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Melek Kolye | 485 | 460s | 313 | 33s 6c |
| recipe_ring_seraphic | 49 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Melek Yüzük | 485 | 385s | 261 | 28s 33c |
| recipe_charm_seraphic | 49 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Melek Tılsım | 485 | 385s | 261 | 28s 33c |
| recipe_necklace_eclipse | 51 | 3× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Tutulma Kolye | 510 | 475s | 352 | 35s 29c |
| recipe_ring_eclipse | 51 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Tutulma Yüzük | 510 | 395s | 293 | 30s 25c |
| recipe_charm_eclipse | 51 | 2× Tungsten Külçesi, 1× Elektrum Külçesi, 2× Obsidyen | 1× Tutulma Tılsım | 510 | 395s | 293 | 30s 25c |
| recipe_necklace_wyrmlord | 54 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Ejder Lordu Kolye | 535 | 495s | 418 | 37s 55c |
| recipe_ring_wyrmlord | 54 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Ejder Lordu Yüzük | 535 | 415s | 349 | 32s 19c |
| recipe_charm_wyrmlord | 54 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Ejder Lordu Tılsım | 535 | 415s | 349 | 32s 19c |
| recipe_necklace_celestine | 56 | 3× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Göksel Kolye | 560 | 510s | 470 | 39s 85c |
| recipe_ring_celestine | 56 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Göksel Yüzük | 560 | 425s | 392 | 34s 16c |
| recipe_charm_celestine | 56 | 2× Titanyum Külçesi, 1× Tungsten Külçesi, 2× Obsidyen | 1× Göksel Tılsım | 560 | 425s | 392 | 34s 16c |
| recipe_necklace_eternal | 59 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Ebedi Kolye | 585 | 535s | 560 | 42s 18c |
| recipe_ring_eternal | 59 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Ebedi Yüzük | 585 | 445s | 467 | 36s 15c |
| recipe_charm_eternal | 59 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Ebedi Tılsım | 585 | 445s | 467 | 36s 15c |
| recipe_necklace_starless | 61 | 3× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Yıldızsız Kolye | 610 | 545s | 629 | 44s 54c |
| recipe_ring_starless | 61 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Yıldızsız Yüzük | 610 | 455s | 524 | 38s 17c |
| recipe_charm_starless | 61 | 2× Paladyum Külçesi, 1× Titanyum Külçesi, 2× Obsidyen | 1× Yıldızsız Tılsım | 610 | 455s | 524 | 38s 17c |
| recipe_necklace_nightsteel | 64 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Gece Çeliği Kolye | 635 | 570s | 750 | 46s 92c |
| recipe_ring_nightsteel | 64 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Gece Çeliği Yüzük | 635 | 475s | 625 | 40s 22c |
| recipe_charm_nightsteel | 64 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Gece Çeliği Tılsım | 635 | 475s | 625 | 40s 22c |
| recipe_necklace_chronoforged | 66 | 3× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Zaman Dövme Kolye | 660 | 585s | 842 | 49s 34c |
| recipe_ring_chronoforged | 66 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Zaman Dövme Yüzük | 660 | 485s | 702 | 42s 29c |
| recipe_charm_chronoforged | 66 | 2× Rodyum Külçesi, 1× Paladyum Külçesi, 2× Obsidyen | 1× Zaman Dövme Tılsım | 660 | 485s | 702 | 42s 29c |
| recipe_necklace_ascendant | 69 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Yükselmiş Kolye | 685 | 605s | 1003 | 51s 78c |
| recipe_ring_ascendant | 69 | 2× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Yükselmiş Yüzük | 685 | 505s | 836 | 44s 38c |
| recipe_charm_ascendant | 69 | 2× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Yükselmiş Tılsım | 685 | 505s | 836 | 44s 38c |
| recipe_necklace_stormborn | 71 | 3× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Fırtına Doğumlu Kolye | 710 | 620s | 1127 | 54s 25c |
| recipe_ring_stormborn | 71 | 2× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Fırtına Doğumlu Yüzük | 710 | 515s | 939 | 46s 50c |
| recipe_charm_stormborn | 71 | 2× Platin Külçesi, 1× Rodyum Külçesi, 2× Obsidyen | 1× Fırtına Doğumlu Tılsım | 710 | 515s | 939 | 46s 50c |
| recipe_necklace_riftforged | 74 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Yarık Dövme Kolye | 735 | 640s | 1342 | 56s 75c |
| recipe_ring_riftforged | 74 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Yarık Dövme Yüzük | 735 | 535s | 1119 | 48s 64c |
| recipe_charm_riftforged | 74 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Yarık Dövme Tılsım | 735 | 535s | 1119 | 48s 64c |
| recipe_necklace_chaos | 76 | 3× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Kaos Kolye | 760 | 655s | 1508 | 59s 27c |
| recipe_ring_chaos | 76 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Kaos Yüzük | 760 | 545s | 1257 | 50s 80c |
| recipe_charm_chaos | 76 | 2× Osmiyum Külçesi, 1× Platin Külçesi, 2× Obsidyen | 1× Kaos Tılsım | 760 | 545s | 1257 | 50s 80c |
| recipe_necklace_skyforged | 79 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Gök Dövme Kolye | 785 | 675s | 1796 | 61s 82c |
| recipe_ring_skyforged | 79 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Gök Dövme Yüzük | 785 | 565s | 1497 | 52s 99c |
| recipe_charm_skyforged | 79 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Gök Dövme Tılsım | 785 | 565s | 1497 | 52s 99c |
| recipe_necklace_archangel | 81 | 3× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Başmelek Kolye | 810 | 690s | 2019 | 64s 39c |
| recipe_ring_archangel | 81 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Başmelek Yüzük | 810 | 575s | 1682 | 55s 19c |
| recipe_charm_archangel | 81 | 2× İridyum Külçesi, 1× Osmiyum Külçesi, 2× Obsidyen | 1× Başmelek Tılsım | 810 | 575s | 1682 | 55s 19c |
| recipe_necklace_aurora | 84 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Şafak Kolye | 835 | 715s | 2404 | 66s 98c |
| recipe_ring_aurora | 84 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Şafak Yüzük | 835 | 595s | 2004 | 57s 42c |
| recipe_charm_aurora | 84 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Şafak Tılsım | 835 | 595s | 2004 | 57s 42c |
| recipe_necklace_ashen | 86 | 3× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Kül Kolye | 860 | 725s | 2701 | 69s 60c |
| recipe_ring_ashen | 86 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Kül Yüzük | 860 | 605s | 2251 | 59s 66c |
| recipe_charm_ashen | 86 | 2× Yıldız Çeliği Külçesi, 1× İridyum Külçesi, 2× Obsidyen | 1× Kül Tılsım | 860 | 605s | 2251 | 59s 66c |
| recipe_necklace_leviathan | 89 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Leviathan Kolye | 885 | 750s | 3217 | 72s 24c |
| recipe_ring_leviathan | 89 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Leviathan Yüzük | 885 | 625s | 2681 | 61s 92c |
| recipe_charm_leviathan | 89 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Leviathan Tılsım | 885 | 625s | 2681 | 61s 92c |
| recipe_necklace_doomforged | 91 | 3× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Kıyamet Dövme Kolye | 910 | 765s | 3615 | 74s 91c |
| recipe_ring_doomforged | 91 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Kıyamet Dövme Yüzük | 910 | 635s | 3012 | 64s 21c |
| recipe_charm_doomforged | 91 | 2× Mitril Külçesi, 1× Yıldız Çeliği Külçesi, 2× Obsidyen | 1× Kıyamet Dövme Tılsım | 910 | 635s | 3012 | 64s 21c |
| recipe_necklace_voidforged | 94 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Hiçlik Dövme Kolye | 935 | 785s | 4306 | 77s 60c |
| recipe_ring_voidforged | 94 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Hiçlik Dövme Yüzük | 935 | 655s | 3588 | 66s 51c |
| recipe_charm_voidforged | 94 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Hiçlik Dövme Tılsım | 935 | 655s | 3588 | 66s 51c |
| recipe_necklace_godslayer | 96 | 3× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Tanrı Katili Kolye | 960 | 800s | 4838 | 80s 30c |
| recipe_ring_godslayer | 96 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Tanrı Katili Yüzük | 960 | 665s | 4032 | 68s 83c |
| recipe_charm_godslayer | 96 | 2× Orikalkum Külçesi, 1× Mitril Külçesi, 2× Obsidyen | 1× Tanrı Katili Tılsım | 960 | 665s | 4032 | 68s 83c |
| recipe_necklace_primordial | 99 | 3× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Obsidyen | 1× Kadim Çağ Kolye | 985 | 820s | 5762 | 83s 3c |
| recipe_ring_primordial | 99 | 2× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Obsidyen | 1× Kadim Çağ Yüzük | 985 | 685s | 4802 | 71s 17c |
| recipe_charm_primordial | 99 | 2× Adamantit Külçesi, 1× Orikalkum Külçesi, 2× Obsidyen | 1× Kadim Çağ Tılsım | 985 | 685s | 4802 | 71s 17c |

### Deri İşçiliği — 235 tarif, beceri 1–99 (tavan 100), çıktı Sv 1–985

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_boots_leather | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Deri Çizme | 1 | 95s | 16 | 2c |
| recipe_chest_leather | 1 | 5× Orman Faresi Postu, 1× Sıradan Ot | 1× Deri Göğüslük | 1 | 155s | 26 | 3c |
| recipe_gloves_leather | 1 | 2× Orman Faresi Postu, 1× Sıradan Ot | 1× Deri Eldiven | 1 | 75s | 13 | 2c |
| recipe_helm_leather | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Deri Kask | 1 | 95s | 16 | 2c |
| recipe_legs_leather | 1 | 4× Orman Faresi Postu, 1× Sıradan Ot | 1× Deri Bacaklık | 1 | 135s | 22 | 3c |
| recipe_cowl_journeyman | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kalfa Deri Başlık | 5 | 95s | 16 | 8c |
| recipe_grips_journeyman | 1 | 2× Orman Faresi Postu, 1× Sıradan Ot | 1× Kalfa Deri Eldiven | 5 | 75s | 13 | 7c |
| recipe_jerkin_journeyman | 1 | 5× Orman Faresi Postu, 1× Sıradan Ot | 1× Kalfa Deri Yelek | 5 | 155s | 26 | 13c |
| recipe_leggings_journeyman | 1 | 4× Orman Faresi Postu, 1× Sıradan Ot | 1× Kalfa Deri Tozluk | 5 | 135s | 22 | 12c |
| recipe_treads_journeyman | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kalfa Deri Çizme | 5 | 95s | 16 | 7c |
| recipe_cowl_fine | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kaliteli Deri Başlık | 10 | 95s | 16 | 20c |
| recipe_grips_fine | 1 | 2× Orman Faresi Postu, 1× Sıradan Ot | 1× Kaliteli Deri Eldiven | 10 | 75s | 13 | 16c |
| recipe_jerkin_fine | 1 | 5× Orman Faresi Postu, 1× Sıradan Ot | 1× Kaliteli Deri Yelek | 10 | 155s | 26 | 33c |
| recipe_leggings_fine | 1 | 4× Orman Faresi Postu, 1× Sıradan Ot | 1× Kaliteli Deri Tozluk | 10 | 135s | 22 | 29c |
| recipe_treads_fine | 1 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kaliteli Deri Çizme | 10 | 95s | 16 | 18c |
| recipe_cowl_veteran | 2 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kıdemli Deri Başlık | 15 | 100s | 17 | 36c |
| recipe_grips_veteran | 2 | 2× Orman Faresi Postu, 1× Sıradan Ot | 1× Kıdemli Deri Eldiven | 15 | 80s | 13 | 29c |
| recipe_jerkin_veteran | 2 | 5× Orman Faresi Postu, 1× Sıradan Ot | 1× Kıdemli Deri Yelek | 15 | 165s | 27 | 58c |
| recipe_leggings_veteran | 2 | 4× Orman Faresi Postu, 1× Sıradan Ot | 1× Kıdemli Deri Tozluk | 15 | 145s | 23 | 51c |
| recipe_treads_veteran | 2 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Kıdemli Deri Çizme | 15 | 100s | 17 | 33c |
| recipe_cowl_champion | 2 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Şampiyon Deri Başlık | 20 | 100s | 17 | 55c |
| recipe_grips_champion | 2 | 2× Orman Faresi Postu, 1× Sıradan Ot | 1× Şampiyon Deri Eldiven | 20 | 80s | 13 | 44c |
| recipe_jerkin_champion | 2 | 5× Orman Faresi Postu, 1× Sıradan Ot | 1× Şampiyon Deri Yelek | 20 | 165s | 27 | 88c |
| recipe_leggings_champion | 2 | 4× Orman Faresi Postu, 1× Sıradan Ot | 1× Şampiyon Deri Tozluk | 20 | 145s | 23 | 77c |
| recipe_treads_champion | 2 | 3× Orman Faresi Postu, 1× Sıradan Ot | 1× Şampiyon Deri Çizme | 20 | 100s | 17 | 50c |
| recipe_cowl_runed | 3 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Rünlü Deri Başlık | 25 | 110s | 18 | 76c |
| recipe_grips_runed | 3 | 2× Küçük Tavşan Postu, 1× Karahindiba | 1× Rünlü Deri Eldiven | 25 | 85s | 14 | 61c |
| recipe_jerkin_runed | 3 | 5× Küçük Tavşan Postu, 1× Karahindiba | 1× Rünlü Deri Yelek | 25 | 175s | 28 | 1s 22c |
| recipe_leggings_runed | 3 | 4× Küçük Tavşan Postu, 1× Karahindiba | 1× Rünlü Deri Tozluk | 25 | 150s | 25 | 1s 7c |
| recipe_treads_runed | 3 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Rünlü Deri Çizme | 25 | 110s | 18 | 69c |
| recipe_cowl_mithril | 4 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Mithril Deri Başlık | 35 | 115s | 19 | 1s 19c |
| recipe_grips_mithril | 4 | 2× Küçük Tavşan Postu, 1× Karahindiba | 1× Mithril Deri Eldiven | 35 | 90s | 15 | 95c |
| recipe_jerkin_mithril | 4 | 5× Küçük Tavşan Postu, 1× Karahindiba | 1× Mithril Deri Yelek | 35 | 180s | 30 | 1s 90c |
| recipe_leggings_mithril | 4 | 4× Küçük Tavşan Postu, 1× Karahindiba | 1× Mithril Deri Tozluk | 35 | 160s | 27 | 1s 66c |
| recipe_treads_mithril | 4 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Mithril Deri Çizme | 35 | 115s | 19 | 1s 7c |
| recipe_cowl_adamant | 5 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Adamant Deri Başlık | 50 | 120s | 20 | 1s 88c |
| recipe_grips_adamant | 5 | 2× Küçük Tavşan Postu, 1× Karahindiba | 1× Adamant Deri Eldiven | 50 | 95s | 16 | 1s 51c |
| recipe_jerkin_adamant | 5 | 5× Küçük Tavşan Postu, 1× Karahindiba | 1× Adamant Deri Yelek | 50 | 190s | 32 | 3s 1c |
| recipe_leggings_adamant | 5 | 4× Küçük Tavşan Postu, 1× Karahindiba | 1× Adamant Deri Tozluk | 50 | 170s | 28 | 2s 64c |
| recipe_treads_adamant | 5 | 3× Küçük Tavşan Postu, 1× Karahindiba | 1× Adamant Deri Çizme | 50 | 120s | 20 | 1s 69c |
| recipe_cowl_wyrmforged | 7 | 3× Tarla Faresi Postu, 1× Yabani Nane | 1× Ejder Dövme Deri Başlık | 65 | 130s | 22 | 2s 65c |
| recipe_grips_wyrmforged | 7 | 2× Tarla Faresi Postu, 1× Yabani Nane | 1× Ejder Dövme Deri Eldiven | 65 | 105s | 18 | 2s 12c |
| recipe_jerkin_wyrmforged | 7 | 5× Tarla Faresi Postu, 1× Yabani Nane | 1× Ejder Dövme Deri Yelek | 65 | 210s | 36 | 4s 24c |
| recipe_leggings_wyrmforged | 7 | 4× Tarla Faresi Postu, 1× Yabani Nane | 1× Ejder Dövme Deri Tozluk | 65 | 185s | 31 | 3s 71c |
| recipe_treads_wyrmforged | 7 | 3× Tarla Faresi Postu, 1× Yabani Nane | 1× Ejder Dövme Deri Çizme | 65 | 130s | 22 | 2s 38c |
| recipe_cowl_starforged | 8 | 3× Tarla Faresi Postu, 1× Yabani Nane | 1× Yıldız Dövme Deri Başlık | 80 | 140s | 24 | 3s 44c |
| recipe_grips_starforged | 8 | 2× Tarla Faresi Postu, 1× Yabani Nane | 1× Yıldız Dövme Deri Eldiven | 80 | 110s | 19 | 2s 75c |
| recipe_jerkin_starforged | 8 | 5× Tarla Faresi Postu, 1× Yabani Nane | 1× Yıldız Dövme Deri Yelek | 80 | 220s | 38 | 5s 51c |
| recipe_leggings_starforged | 8 | 4× Tarla Faresi Postu, 1× Yabani Nane | 1× Yıldız Dövme Deri Tozluk | 80 | 195s | 34 | 4s 82c |
| recipe_treads_starforged | 8 | 3× Tarla Faresi Postu, 1× Yabani Nane | 1× Yıldız Dövme Deri Çizme | 80 | 140s | 24 | 3s 10c |
| recipe_cowl_titan | 10 | 3× Sincap Postu, 1× Papatya | 1× Titan Deri Başlık | 95 | 150s | 27 | 4s 33c |
| recipe_grips_titan | 10 | 2× Sincap Postu, 1× Papatya | 1× Titan Deri Eldiven | 95 | 120s | 21 | 3s 46c |
| recipe_jerkin_titan | 10 | 5× Sincap Postu, 1× Papatya | 1× Titan Deri Yelek | 95 | 240s | 43 | 6s 93c |
| recipe_leggings_titan | 10 | 4× Sincap Postu, 1× Papatya | 1× Titan Deri Tozluk | 95 | 210s | 37 | 6s 6c |
| recipe_treads_titan | 10 | 3× Sincap Postu, 1× Papatya | 1× Titan Deri Çizme | 95 | 150s | 27 | 3s 89c |
| recipe_cowl_royal | 11 | 3× Yaban Tavşanı Postu, 1× Lavanta | 1× Kraliyet Deri Başlık | 110 | 155s | 28 | 5s 25c |
| recipe_grips_royal | 11 | 2× Yaban Tavşanı Postu, 1× Lavanta | 1× Kraliyet Deri Eldiven | 110 | 125s | 23 | 4s 20c |
| recipe_jerkin_royal | 11 | 5× Yaban Tavşanı Postu, 1× Lavanta | 1× Kraliyet Deri Yelek | 110 | 250s | 46 | 8s 40c |
| recipe_leggings_royal | 11 | 4× Yaban Tavşanı Postu, 1× Lavanta | 1× Kraliyet Deri Tozluk | 110 | 220s | 40 | 7s 35c |
| recipe_treads_royal | 11 | 3× Yaban Tavşanı Postu, 1× Lavanta | 1× Kraliyet Deri Çizme | 110 | 155s | 28 | 4s 72c |
| recipe_cowl_sunforged | 14 | 3× Porsuk Postu, 1× Ayçiçeği | 1× Güneş Dövme Deri Başlık | 135 | 175s | 34 | 6s 85c |
| recipe_grips_sunforged | 14 | 2× Porsuk Postu, 1× Ayçiçeği | 1× Güneş Dövme Deri Eldiven | 135 | 140s | 27 | 5s 48c |
| recipe_jerkin_sunforged | 14 | 5× Porsuk Postu, 1× Ayçiçeği | 1× Güneş Dövme Deri Yelek | 135 | 280s | 54 | 10s 96c |
| recipe_leggings_sunforged | 14 | 4× Porsuk Postu, 1× Ayçiçeği | 1× Güneş Dövme Deri Tozluk | 135 | 245s | 48 | 9s 59c |
| recipe_treads_sunforged | 14 | 3× Porsuk Postu, 1× Ayçiçeği | 1× Güneş Dövme Deri Çizme | 135 | 175s | 34 | 6s 16c |
| recipe_cowl_obsidian | 16 | 3× Porsuk Postu, 1× Acıkök | 1× Obsidyen Deri Başlık | 160 | 185s | 38 | 8s 54c |
| recipe_grips_obsidian | 16 | 2× Porsuk Postu, 1× Acıkök | 1× Obsidyen Deri Eldiven | 160 | 150s | 30 | 6s 84c |
| recipe_jerkin_obsidian | 16 | 5× Porsuk Postu, 1× Acıkök | 1× Obsidyen Deri Yelek | 160 | 300s | 61 | 13s 67c |
| recipe_leggings_obsidian | 16 | 4× Porsuk Postu, 1× Acıkök | 1× Obsidyen Deri Tozluk | 160 | 260s | 53 | 11s 96c |
| recipe_treads_obsidian | 16 | 3× Porsuk Postu, 1× Acıkök | 1× Obsidyen Deri Çizme | 160 | 185s | 38 | 7s 68c |
| recipe_cowl_sandsteel | 19 | 3× Kızıl Tilki Postu, 1× Bataklık Sazı | 1× Kum Çeliği Deri Başlık | 185 | 205s | 46 | 10s 32c |
| recipe_grips_sandsteel | 19 | 2× Kızıl Tilki Postu, 1× Bataklık Sazı | 1× Kum Çeliği Deri Eldiven | 185 | 165s | 36 | 8s 26c |
| recipe_jerkin_sandsteel | 19 | 5× Kızıl Tilki Postu, 1× Bataklık Sazı | 1× Kum Çeliği Deri Yelek | 185 | 325s | 73 | 16s 51c |
| recipe_leggings_sandsteel | 19 | 4× Kızıl Tilki Postu, 1× Bataklık Sazı | 1× Kum Çeliği Deri Tozluk | 185 | 285s | 64 | 14s 45c |
| recipe_treads_sandsteel | 19 | 3× Kızıl Tilki Postu, 1× Bataklık Sazı | 1× Kum Çeliği Deri Çizme | 185 | 205s | 46 | 9s 28c |
| recipe_cowl_runebound | 21 | 3× Yaban Domuzu Postu, 1× Mavi Devedikeni | 1× Rün Bağlı Deri Başlık | 210 | 215s | 51 | 12s 17c |
| recipe_grips_runebound | 21 | 2× Yaban Domuzu Postu, 1× Mavi Devedikeni | 1× Rün Bağlı Deri Eldiven | 210 | 175s | 41 | 9s 74c |
| recipe_jerkin_runebound | 21 | 5× Yaban Domuzu Postu, 1× Mavi Devedikeni | 1× Rün Bağlı Deri Yelek | 210 | 345s | 82 | 19s 47c |
| recipe_leggings_runebound | 21 | 4× Yaban Domuzu Postu, 1× Mavi Devedikeni | 1× Rün Bağlı Deri Tozluk | 210 | 300s | 71 | 17s 4c |
| recipe_treads_runebound | 21 | 3× Yaban Domuzu Postu, 1× Mavi Devedikeni | 1× Rün Bağlı Deri Çizme | 210 | 215s | 51 | 10s 94c |
| recipe_cowl_phoenix | 24 | 3× Dağ Keçisi Postu, 1× Kızıl Adaçayı | 1× Anka Deri Başlık | 235 | 235s | 61 | 14s 8c |
| recipe_grips_phoenix | 24 | 2× Dağ Keçisi Postu, 1× Kızıl Adaçayı | 1× Anka Deri Eldiven | 235 | 185s | 49 | 11s 27c |
| recipe_jerkin_phoenix | 24 | 5× Dağ Keçisi Postu, 1× Kızıl Adaçayı | 1× Anka Deri Yelek | 235 | 375s | 97 | 22s 53c |
| recipe_leggings_phoenix | 24 | 4× Dağ Keçisi Postu, 1× Kızıl Adaçayı | 1× Anka Deri Tozluk | 235 | 330s | 85 | 19s 72c |
| recipe_treads_phoenix | 24 | 3× Dağ Keçisi Postu, 1× Kızıl Adaçayı | 1× Anka Deri Çizme | 235 | 235s | 61 | 12s 67c |
| recipe_cowl_dreadforged | 26 | 3× Yaban Kurdu Postu, 1× Güzelavrat Otu | 1× Dehşet Dövme Deri Başlık | 260 | 245s | 68 | 16s 6c |
| recipe_grips_dreadforged | 26 | 2× Yaban Kurdu Postu, 1× Güzelavrat Otu | 1× Dehşet Dövme Deri Eldiven | 260 | 195s | 55 | 12s 85c |
| recipe_jerkin_dreadforged | 26 | 5× Yaban Kurdu Postu, 1× Güzelavrat Otu | 1× Dehşet Dövme Deri Yelek | 260 | 395s | 109 | 25s 70c |
| recipe_leggings_dreadforged | 26 | 4× Yaban Kurdu Postu, 1× Güzelavrat Otu | 1× Dehşet Dövme Deri Tozluk | 260 | 345s | 96 | 22s 49c |
| recipe_treads_dreadforged | 26 | 3× Yaban Kurdu Postu, 1× Güzelavrat Otu | 1× Dehşet Dövme Deri Çizme | 260 | 245s | 68 | 14s 44c |
| recipe_cowl_sovereign | 29 | 3× Yaban Mandası Derisi, 1× Altın Yonca | 1× Hükümdar Deri Başlık | 285 | 265s | 81 | 18s 10c |
| recipe_grips_sovereign | 29 | 2× Yaban Mandası Derisi, 1× Altın Yonca | 1× Hükümdar Deri Eldiven | 285 | 210s | 65 | 14s 48c |
| recipe_jerkin_sovereign | 29 | 5× Yaban Mandası Derisi, 1× Altın Yonca | 1× Hükümdar Deri Yelek | 285 | 420s | 130 | 28s 96c |
| recipe_leggings_sovereign | 29 | 4× Yaban Mandası Derisi, 1× Altın Yonca | 1× Hükümdar Deri Tozluk | 285 | 370s | 114 | 25s 34c |
| recipe_treads_sovereign | 29 | 3× Yaban Mandası Derisi, 1× Altın Yonca | 1× Hükümdar Deri Çizme | 285 | 265s | 81 | 16s 28c |
| recipe_cowl_coral | 31 | 3× Yaban Mandası Derisi, 1× Altın Yonca | 1× Mercan Deri Başlık | 310 | 275s | 91 | 20s 19c |
| recipe_grips_coral | 31 | 2× Yaban Mandası Derisi, 1× Altın Yonca | 1× Mercan Deri Eldiven | 310 | 220s | 73 | 16s 15c |
| recipe_jerkin_coral | 31 | 5× Yaban Mandası Derisi, 1× Altın Yonca | 1× Mercan Deri Yelek | 310 | 440s | 146 | 32s 30c |
| recipe_leggings_coral | 31 | 4× Yaban Mandası Derisi, 1× Altın Yonca | 1× Mercan Deri Tozluk | 310 | 385s | 128 | 28s 26c |
| recipe_treads_coral | 31 | 3× Yaban Mandası Derisi, 1× Altın Yonca | 1× Mercan Deri Çizme | 310 | 275s | 91 | 18s 16c |
| recipe_cowl_drakebone | 34 | 3× Dev Geyik Postu, 1× Düş Kökü | 1× Ejderha Kemiği Deri Başlık | 335 | 295s | 109 | 22s 33c |
| recipe_grips_drakebone | 34 | 2× Dev Geyik Postu, 1× Düş Kökü | 1× Ejderha Kemiği Deri Eldiven | 335 | 235s | 87 | 17s 86c |
| recipe_jerkin_drakebone | 34 | 5× Dev Geyik Postu, 1× Düş Kökü | 1× Ejderha Kemiği Deri Yelek | 335 | 470s | 174 | 35s 73c |
| recipe_leggings_drakebone | 34 | 4× Dev Geyik Postu, 1× Düş Kökü | 1× Ejderha Kemiği Deri Tozluk | 335 | 410s | 152 | 31s 26c |
| recipe_treads_drakebone | 34 | 3× Dev Geyik Postu, 1× Düş Kökü | 1× Ejderha Kemiği Deri Çizme | 335 | 295s | 109 | 20s 8c |
| recipe_cowl_mistwoven | 36 | 3× Dev Kurt Postu, 1× Gümüş Yaprak | 1× Sis Örgü Deri Başlık | 360 | 305s | 122 | 24s 52c |
| recipe_grips_mistwoven | 36 | 2× Dev Kurt Postu, 1× Gümüş Yaprak | 1× Sis Örgü Deri Eldiven | 360 | 245s | 98 | 19s 62c |
| recipe_jerkin_mistwoven | 36 | 5× Dev Kurt Postu, 1× Gümüş Yaprak | 1× Sis Örgü Deri Yelek | 360 | 490s | 196 | 39s 23c |
| recipe_leggings_mistwoven | 36 | 4× Dev Kurt Postu, 1× Gümüş Yaprak | 1× Sis Örgü Deri Tozluk | 360 | 430s | 171 | 34s 33c |
| recipe_treads_mistwoven | 36 | 3× Dev Kurt Postu, 1× Gümüş Yaprak | 1× Sis Örgü Deri Çizme | 360 | 305s | 122 | 22s 5c |
| recipe_cowl_prismatic | 39 | 3× Dağ Aslanı Postu, 1× Kan Dikeni | 1× Prizmatik Deri Başlık | 385 | 325s | 146 | 26s 76c |
| recipe_grips_prismatic | 39 | 2× Dağ Aslanı Postu, 1× Kan Dikeni | 1× Prizmatik Deri Eldiven | 385 | 260s | 116 | 21s 41c |
| recipe_jerkin_prismatic | 39 | 5× Dağ Aslanı Postu, 1× Kan Dikeni | 1× Prizmatik Deri Yelek | 385 | 520s | 233 | 42s 81c |
| recipe_leggings_prismatic | 39 | 4× Dağ Aslanı Postu, 1× Kan Dikeni | 1× Prizmatik Deri Tozluk | 385 | 455s | 204 | 37s 46c |
| recipe_treads_prismatic | 39 | 3× Dağ Aslanı Postu, 1× Kan Dikeni | 1× Prizmatik Deri Çizme | 385 | 325s | 146 | 24s 6c |
| recipe_cowl_magmaforged | 41 | 3× Kara Panter Postu, 1× Hayalet Mantarı | 1× Magma Dövme Deri Başlık | 410 | 335s | 164 | 29s 4c |
| recipe_grips_magmaforged | 41 | 2× Kara Panter Postu, 1× Hayalet Mantarı | 1× Magma Dövme Deri Eldiven | 410 | 270s | 131 | 23s 23c |
| recipe_jerkin_magmaforged | 41 | 5× Kara Panter Postu, 1× Hayalet Mantarı | 1× Magma Dövme Deri Yelek | 410 | 540s | 262 | 46s 46c |
| recipe_leggings_magmaforged | 41 | 4× Kara Panter Postu, 1× Hayalet Mantarı | 1× Magma Dövme Deri Tozluk | 410 | 470s | 229 | 40s 65c |
| recipe_treads_magmaforged | 41 | 3× Kara Panter Postu, 1× Hayalet Mantarı | 1× Magma Dövme Deri Çizme | 410 | 335s | 164 | 26s 12c |
| recipe_cowl_tempest | 44 | 3× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Bora Deri Başlık | 435 | 355s | 195 | 31s 36c |
| recipe_grips_tempest | 44 | 2× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Bora Deri Eldiven | 435 | 285s | 156 | 25s 9c |
| recipe_jerkin_tempest | 44 | 5× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Bora Deri Yelek | 435 | 565s | 312 | 50s 18c |
| recipe_leggings_tempest | 44 | 4× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Bora Deri Tozluk | 435 | 495s | 273 | 43s 90c |
| recipe_treads_tempest | 44 | 3× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Bora Deri Çizme | 435 | 355s | 195 | 28s 20c |
| recipe_cowl_voidtouched | 46 | 3× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Hiçlik Deri Başlık | 460 | 365s | 219 | 33s 72c |
| recipe_grips_voidtouched | 46 | 2× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Hiçlik Deri Eldiven | 460 | 295s | 175 | 26s 98c |
| recipe_jerkin_voidtouched | 46 | 5× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Hiçlik Deri Yelek | 460 | 585s | 350 | 53s 96c |
| recipe_leggings_voidtouched | 46 | 4× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Hiçlik Deri Tozluk | 460 | 510s | 306 | 47s 21c |
| recipe_treads_voidtouched | 46 | 3× Kara Panter Postu, 1× Fırtına Çiçeği | 1× Hiçlik Deri Çizme | 460 | 365s | 219 | 30s 33c |
| recipe_cowl_seraphic | 49 | 3× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Melek Deri Başlık | 485 | 385s | 261 | 36s 12c |
| recipe_grips_seraphic | 49 | 2× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Melek Deri Eldiven | 485 | 305s | 209 | 28s 90c |
| recipe_jerkin_seraphic | 49 | 5× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Melek Deri Yelek | 485 | 615s | 417 | 57s 80c |
| recipe_leggings_seraphic | 49 | 4× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Melek Deri Tozluk | 485 | 540s | 365 | 50s 58c |
| recipe_treads_seraphic | 49 | 3× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Melek Deri Çizme | 485 | 385s | 261 | 32s 49c |
| recipe_cowl_eclipse | 51 | 3× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Tutulma Deri Başlık | 510 | 395s | 293 | 38s 56c |
| recipe_grips_eclipse | 51 | 2× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Tutulma Deri Eldiven | 510 | 315s | 234 | 30s 85c |
| recipe_jerkin_eclipse | 51 | 5× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Tutulma Deri Yelek | 510 | 635s | 469 | 61s 70c |
| recipe_leggings_eclipse | 51 | 4× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Tutulma Deri Tozluk | 510 | 555s | 410 | 53s 99c |
| recipe_treads_eclipse | 51 | 3× Kor Domuzu Postu, 1× Kor Yaprağı | 1× Tutulma Deri Çizme | 510 | 395s | 293 | 34s 68c |
| recipe_cowl_wyrmlord | 54 | 3× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Ejder Lordu Deri Başlık | 535 | 415s | 349 | 41s 4c |
| recipe_grips_wyrmlord | 54 | 2× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Ejder Lordu Deri Eldiven | 535 | 330s | 279 | 32s 83c |
| recipe_jerkin_wyrmlord | 54 | 5× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Ejder Lordu Deri Yelek | 535 | 660s | 558 | 65s 66c |
| recipe_leggings_wyrmlord | 54 | 4× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Ejder Lordu Deri Tozluk | 535 | 580s | 488 | 57s 46c |
| recipe_treads_wyrmlord | 54 | 3× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Ejder Lordu Deri Çizme | 535 | 415s | 349 | 36s 91c |
| recipe_cowl_celestine | 56 | 3× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Göksel Deri Başlık | 560 | 425s | 392 | 43s 55c |
| recipe_grips_celestine | 56 | 2× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Göksel Deri Eldiven | 560 | 340s | 314 | 34s 84c |
| recipe_jerkin_celestine | 56 | 5× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Göksel Deri Yelek | 560 | 680s | 627 | 69s 68c |
| recipe_leggings_celestine | 56 | 4× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Göksel Deri Tozluk | 560 | 595s | 549 | 60s 97c |
| recipe_treads_celestine | 56 | 3× Ayaz Kurdu Postu, 1× Ayaz Çiçeği | 1× Göksel Deri Çizme | 560 | 425s | 392 | 39s 17c |
| recipe_cowl_eternal | 59 | 3× Dev Ayı Postu, 1× Ay Gölgesi | 1× Ebedi Deri Başlık | 585 | 445s | 467 | 46s 9c |
| recipe_grips_eternal | 59 | 2× Dev Ayı Postu, 1× Ay Gölgesi | 1× Ebedi Deri Eldiven | 585 | 355s | 373 | 36s 88c |
| recipe_jerkin_eternal | 59 | 5× Dev Ayı Postu, 1× Ay Gölgesi | 1× Ebedi Deri Yelek | 585 | 710s | 747 | 73s 75c |
| recipe_leggings_eternal | 59 | 4× Dev Ayı Postu, 1× Ay Gölgesi | 1× Ebedi Deri Tozluk | 585 | 620s | 653 | 64s 53c |
| recipe_treads_eternal | 59 | 3× Dev Ayı Postu, 1× Ay Gölgesi | 1× Ebedi Deri Çizme | 585 | 445s | 467 | 41s 45c |
| recipe_cowl_starless | 61 | 3× Dev Ayı Postu, 1× Ay Gölgesi | 1× Yıldızsız Deri Başlık | 610 | 455s | 524 | 48s 67c |
| recipe_grips_starless | 61 | 2× Dev Ayı Postu, 1× Ay Gölgesi | 1× Yıldızsız Deri Eldiven | 610 | 365s | 420 | 38s 94c |
| recipe_jerkin_starless | 61 | 5× Dev Ayı Postu, 1× Ay Gölgesi | 1× Yıldızsız Deri Yelek | 610 | 730s | 839 | 77s 87c |
| recipe_leggings_starless | 61 | 4× Dev Ayı Postu, 1× Ay Gölgesi | 1× Yıldızsız Deri Tozluk | 610 | 640s | 734 | 68s 14c |
| recipe_treads_starless | 61 | 3× Dev Ayı Postu, 1× Ay Gölgesi | 1× Yıldızsız Deri Çizme | 610 | 455s | 524 | 43s 77c |
| recipe_cowl_nightsteel | 64 | 3× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Gece Çeliği Deri Başlık | 635 | 475s | 625 | 51s 28c |
| recipe_grips_nightsteel | 64 | 2× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Gece Çeliği Deri Eldiven | 635 | 380s | 500 | 41s 2c |
| recipe_jerkin_nightsteel | 64 | 5× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Gece Çeliği Deri Yelek | 635 | 760s | 1000 | 82s 5c |
| recipe_leggings_nightsteel | 64 | 4× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Gece Çeliği Deri Tozluk | 635 | 665s | 875 | 71s 79c |
| recipe_treads_nightsteel | 64 | 3× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Gece Çeliği Deri Çizme | 635 | 475s | 625 | 46s 12c |
| recipe_cowl_chronoforged | 66 | 3× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Zaman Dövme Deri Başlık | 660 | 485s | 702 | 53s 92c |
| recipe_grips_chronoforged | 66 | 2× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Zaman Dövme Deri Eldiven | 660 | 390s | 562 | 43s 14c |
| recipe_jerkin_chronoforged | 66 | 5× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Zaman Dövme Deri Yelek | 660 | 780s | 1123 | 86s 27c |
| recipe_leggings_chronoforged | 66 | 4× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Zaman Dövme Deri Tozluk | 660 | 680s | 983 | 75s 49c |
| recipe_treads_chronoforged | 66 | 3× Ruh Geyiği Postu, 1× Ruh Nilüferi | 1× Zaman Dövme Deri Çizme | 660 | 485s | 702 | 48s 49c |
| recipe_cowl_ascendant | 69 | 3× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Yükselmiş Deri Başlık | 685 | 505s | 836 | 56s 59c |
| recipe_grips_ascendant | 69 | 2× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Yükselmiş Deri Eldiven | 685 | 405s | 669 | 45s 27c |
| recipe_jerkin_ascendant | 69 | 5× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Yükselmiş Deri Yelek | 685 | 805s | 1338 | 90s 54c |
| recipe_leggings_ascendant | 69 | 4× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Yükselmiş Deri Tozluk | 685 | 705s | 1170 | 79s 23c |
| recipe_treads_ascendant | 69 | 3× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Yükselmiş Deri Çizme | 685 | 505s | 836 | 50s 89c |
| recipe_cowl_stormborn | 71 | 3× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Fırtına Doğumlu Deri Başlık | 710 | 515s | 939 | 59s 29c |
| recipe_grips_stormborn | 71 | 2× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Fırtına Doğumlu Deri Eldiven | 710 | 415s | 751 | 47s 43c |
| recipe_jerkin_stormborn | 71 | 5× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Fırtına Doğumlu Deri Yelek | 710 | 825s | 1503 | 94s 86c |
| recipe_leggings_stormborn | 71 | 4× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Fırtına Doğumlu Deri Tozluk | 710 | 720s | 1315 | 83s |
| recipe_treads_stormborn | 71 | 3× Kadim Mamut Postu, 1× Kadim Adamotu | 1× Fırtına Doğumlu Deri Çizme | 710 | 515s | 939 | 53s 32c |
| recipe_cowl_riftforged | 74 | 3× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Yarık Dövme Deri Başlık | 735 | 535s | 1119 | 62s 2c |
| recipe_grips_riftforged | 74 | 2× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Yarık Dövme Deri Eldiven | 735 | 425s | 895 | 49s 61c |
| recipe_jerkin_riftforged | 74 | 5× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Yarık Dövme Deri Yelek | 735 | 855s | 1790 | 99s 23c |
| recipe_leggings_riftforged | 74 | 4× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Yarık Dövme Deri Tozluk | 735 | 750s | 1566 | 86s 82c |
| recipe_treads_riftforged | 74 | 3× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Yarık Dövme Deri Çizme | 735 | 535s | 1119 | 55s 78c |
| recipe_cowl_chaos | 76 | 3× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Kaos Deri Başlık | 760 | 545s | 1257 | 64s 77c |
| recipe_grips_chaos | 76 | 2× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Kaos Deri Eldiven | 760 | 435s | 1006 | 51s 82c |
| recipe_jerkin_chaos | 76 | 5× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Kaos Deri Yelek | 760 | 875s | 2011 | 1g 3s 64c |
| recipe_leggings_chaos | 76 | 4× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Kaos Deri Tozluk | 760 | 765s | 1760 | 90s 68c |
| recipe_treads_chaos | 76 | 3× Gölge Kimera Derisi, 1× Ejderkanı Otu | 1× Kaos Deri Çizme | 760 | 545s | 1257 | 58s 25c |
| recipe_cowl_skyforged | 79 | 3× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Gök Dövme Deri Başlık | 785 | 565s | 1497 | 67s 56c |
| recipe_grips_skyforged | 79 | 2× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Gök Dövme Deri Eldiven | 785 | 450s | 1198 | 54s 5c |
| recipe_jerkin_skyforged | 79 | 5× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Gök Dövme Deri Yelek | 785 | 900s | 2395 | 1g 8s 9c |
| recipe_leggings_skyforged | 79 | 4× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Gök Dövme Deri Tozluk | 785 | 790s | 2096 | 94s 58c |
| recipe_treads_skyforged | 79 | 3× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Gök Dövme Deri Çizme | 785 | 565s | 1497 | 60s 76c |
| recipe_cowl_archangel | 81 | 3× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Başmelek Deri Başlık | 810 | 575s | 1682 | 70s 37c |
| recipe_grips_archangel | 81 | 2× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Başmelek Deri Eldiven | 810 | 460s | 1346 | 56s 29c |
| recipe_jerkin_archangel | 81 | 5× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Başmelek Deri Yelek | 810 | 920s | 2692 | 1g 12s 59c |
| recipe_leggings_archangel | 81 | 4× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Başmelek Deri Tozluk | 810 | 805s | 2355 | 98s 51c |
| recipe_treads_archangel | 81 | 3× Yaşlı Wyvern Pulu, 1× Yıldız Çiçeği | 1× Başmelek Deri Çizme | 810 | 575s | 1682 | 63s 28c |
| recipe_cowl_aurora | 84 | 3× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Şafak Deri Başlık | 835 | 595s | 2004 | 73s 20c |
| recipe_grips_aurora | 84 | 2× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Şafak Deri Eldiven | 835 | 475s | 1603 | 58s 56c |
| recipe_jerkin_aurora | 84 | 5× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Şafak Deri Yelek | 835 | 950s | 3206 | 1g 17s 13c |
| recipe_leggings_aurora | 84 | 4× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Şafak Deri Tozluk | 835 | 830s | 2805 | 1g 2s 48c |
| recipe_treads_aurora | 84 | 3× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Şafak Deri Çizme | 835 | 595s | 2004 | 65s 84c |
| recipe_cowl_ashen | 86 | 3× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Kül Deri Başlık | 860 | 605s | 2251 | 76s 6c |
| recipe_grips_ashen | 86 | 2× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Kül Deri Eldiven | 860 | 485s | 1801 | 60s 85c |
| recipe_jerkin_ashen | 86 | 5× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Kül Deri Yelek | 860 | 970s | 3602 | 1g 21s 70c |
| recipe_leggings_ashen | 86 | 4× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Kül Deri Tozluk | 860 | 850s | 3151 | 1g 6s 49c |
| recipe_treads_ashen | 86 | 3× Yaşlı Wyvern Pulu, 1× Boşluk Orkidesi | 1× Kül Deri Çizme | 860 | 605s | 2251 | 68s 41c |
| recipe_cowl_leviathan | 89 | 3× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Leviathan Deri Başlık | 885 | 625s | 2681 | 78s 95c |
| recipe_grips_leviathan | 89 | 2× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Leviathan Deri Eldiven | 885 | 500s | 2145 | 63s 16c |
| recipe_jerkin_leviathan | 89 | 5× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Leviathan Deri Yelek | 885 | 1000s | 4290 | 1g 26s 32c |
| recipe_leggings_leviathan | 89 | 4× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Leviathan Deri Tozluk | 885 | 875s | 3753 | 1g 10s 53c |
| recipe_treads_leviathan | 89 | 3× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Leviathan Deri Çizme | 885 | 625s | 2681 | 71s 1c |
| recipe_cowl_doomforged | 91 | 3× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Kıyamet Dövme Deri Başlık | 910 | 635s | 3012 | 81s 86c |
| recipe_grips_doomforged | 91 | 2× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Kıyamet Dövme Deri Eldiven | 910 | 510s | 2410 | 65s 49c |
| recipe_jerkin_doomforged | 91 | 5× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Kıyamet Dövme Deri Yelek | 910 | 1020s | 4820 | 1g 30s 98c |
| recipe_leggings_doomforged | 91 | 4× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Kıyamet Dövme Deri Tozluk | 910 | 890s | 4218 | 1g 14s 61c |
| recipe_treads_doomforged | 91 | 3× Boşluk Behemotu Derisi, 1× Göksel Nilüfer | 1× Kıyamet Dövme Deri Çizme | 910 | 635s | 3012 | 73s 62c |
| recipe_cowl_voidforged | 94 | 3× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Hiçlik Dövme Deri Başlık | 935 | 655s | 3588 | 84s 80c |
| recipe_grips_voidforged | 94 | 2× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Hiçlik Dövme Deri Eldiven | 935 | 525s | 2870 | 67s 84c |
| recipe_jerkin_voidforged | 94 | 5× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Hiçlik Dövme Deri Yelek | 935 | 1045s | 5741 | 1g 35s 68c |
| recipe_leggings_voidforged | 94 | 4× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Hiçlik Dövme Deri Tozluk | 935 | 915s | 5023 | 1g 18s 72c |
| recipe_treads_voidforged | 94 | 3× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Hiçlik Dövme Deri Çizme | 935 | 655s | 3588 | 76s 26c |
| recipe_cowl_godslayer | 96 | 3× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Tanrı Katili Deri Başlık | 960 | 665s | 4032 | 87s 76c |
| recipe_grips_godslayer | 96 | 2× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Tanrı Katili Deri Eldiven | 960 | 535s | 3225 | 70s 21c |
| recipe_jerkin_godslayer | 96 | 5× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Tanrı Katili Deri Yelek | 960 | 1065s | 6450 | 1g 40s 42c |
| recipe_leggings_godslayer | 96 | 4× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Tanrı Katili Deri Tozluk | 960 | 930s | 5644 | 1g 22s 86c |
| recipe_treads_godslayer | 96 | 3× Boşluk Behemotu Derisi, 1× Anka Çiçeği | 1× Tanrı Katili Deri Çizme | 960 | 665s | 4032 | 78s 93c |
| recipe_cowl_primordial | 99 | 3× Göksel Ejder Pulu, 1× Dünya Kökü | 1× Kadim Çağ Deri Başlık | 985 | 685s | 4802 | 90s 74c |
| recipe_grips_primordial | 99 | 2× Göksel Ejder Pulu, 1× Dünya Kökü | 1× Kadim Çağ Deri Eldiven | 985 | 545s | 3841 | 72s 59c |
| recipe_jerkin_primordial | 99 | 5× Göksel Ejder Pulu, 1× Dünya Kökü | 1× Kadim Çağ Deri Yelek | 985 | 1095s | 7682 | 1g 45s 19c |
| recipe_leggings_primordial | 99 | 4× Göksel Ejder Pulu, 1× Dünya Kökü | 1× Kadim Çağ Deri Tozluk | 985 | 960s | 6722 | 1g 27s 4c |
| recipe_treads_primordial | 99 | 3× Göksel Ejder Pulu, 1× Dünya Kökü | 1× Kadim Çağ Deri Çizme | 985 | 685s | 4802 | 81s 61c |

### Ergitme — 24 tarif, beceri 1–87 (tavan 100), çıktı Sv 1–1

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_aluminium_bar | 1 | 2× Alüminyum Cevheri, 1× Kömür | 1× Alüminyum Külçesi | 1 | 30s | 48 | 2c |
| recipe_lead_bar | 1 | 2× Kurşun Cevheri, 1× Kömür | 1× Kurşun Külçesi | 1 | 30s | 48 | 2c |
| recipe_iron_bar | 3 | 2× Demir Cevheri, 1× Kömür | 1× Demir Külçesi | 1 | 35s | 48 | 2c |
| recipe_tin_bar | 5 | 2× Kalay Cevheri, 1× Kömür | 1× Kalay Külçesi | 1 | 40s | 48 | 2c |
| recipe_copper_bar | 7 | 2× Bakır Cevheri, 1× Kömür | 1× Bakır Külçesi | 1 | 40s | 72 | 5c |
| recipe_bronze_bar | 9 | 2× Bronz Cevheri, 1× Kömür | 1× Bronz Külçesi | 1 | 45s | 72 | 5c |
| recipe_zinc_bar | 15 | 2× Çinko Cevheri, 1× Kömür | 1× Çinko Külçesi | 1 | 50s | 72 | 5c |
| recipe_nickel_bar | 17 | 2× Nikel Cevheri, 2× Kömür | 1× Nikel Külçesi | 1 | 55s | 104 | 10c |
| recipe_chromium_bar | 20 | 2× Krom Cevheri, 2× Kömür | 1× Krom Külçesi | 1 | 60s | 104 | 10c |
| recipe_cobalt_bar | 23 | 2× Kobalt Cevheri, 2× Kömür | 1× Kobalt Külçesi | 1 | 65s | 104 | 10c |
| recipe_silver_bar | 26 | 2× Gümüş Cevheri, 2× Kömür | 1× Gümüş Külçesi | 1 | 70s | 104 | 10c |
| recipe_gold_bar | 29 | 2× Altın Cevheri, 2× Kömür | 1× Altın Külçesi | 1 | 75s | 104 | 10c |
| recipe_electrum_bar | 32 | 3× Elektrum Cevheri, 2× Kömür | 1× Elektrum Külçesi | 1 | 80s | 216 | 20c |
| recipe_tungsten_bar | 37 | 3× Tungsten Cevheri, 2× Kömür | 1× Tungsten Külçesi | 1 | 85s | 216 | 20c |
| recipe_titanium_bar | 42 | 3× Titanyum Cevheri, 2× Kömür | 1× Titanyum Külçesi | 1 | 95s | 216 | 20c |
| recipe_palladium_bar | 47 | 3× Paladyum Cevheri, 3× Kömür | 1× Paladyum Külçesi | 1 | 100s | 300 | 40c |
| recipe_rhodium_bar | 52 | 3× Rodyum Cevheri, 3× Kömür | 1× Rodyum Külçesi | 1 | 110s | 300 | 40c |
| recipe_platinum_bar | 57 | 3× Platin Cevheri, 3× Kömür | 1× Platin Külçesi | 1 | 115s | 300 | 40c |
| recipe_osmium_bar | 62 | 4× Osmiyum Cevheri, 3× Kömür | 1× Osmiyum Külçesi | 1 | 125s | 544 | 80c |
| recipe_iridium_bar | 67 | 4× İridyum Cevheri, 3× Kömür | 1× İridyum Külçesi | 1 | 130s | 544 | 80c |
| recipe_starsteel_bar | 72 | 4× Yıldız Çeliği Cevheri, 3× Kömür | 1× Yıldız Çeliği Külçesi | 1 | 140s | 544 | 80c |
| recipe_mythril_bar | 77 | 5× Mitril Cevheri, 4× Kömür | 1× Mitril Külçesi | 1 | 145s | 900 | 1s 50c |
| recipe_orichalcum_bar | 82 | 5× Orikalkum Cevheri, 4× Kömür | 1× Orikalkum Külçesi | 1 | 155s | 900 | 1s 50c |
| recipe_adamantite_bar | 87 | 5× Adamantit Cevheri, 4× Kömür | 1× Adamantit Külçesi | 1 | 160s | 900 | 1s 50c |

### Terzilik — 282 tarif, beceri 1–99 (tavan 100), çıktı Sv 1–985

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_hood_novice | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Kukuleta | 1 | 95s | 16 | 2c |
| recipe_robe_novice | 1 | 10× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Cübbe | 1 | 155s | 26 | 3c |
| recipe_sandals_novice | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Sandalet | 1 | 95s | 16 | 2c |
| recipe_trousers_novice | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Şalvar | 1 | 135s | 22 | 3c |
| recipe_wraps_novice | 1 | 4× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Sargı | 1 | 75s | 13 | 2c |
| recipe_tome_novice | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Acemi Kutsal Kitap | 1 | 155s | 26 | 4c |
| recipe_hood_journeyman | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Kukuleta | 5 | 95s | 16 | 8c |
| recipe_robe_journeyman | 1 | 10× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Cübbe | 5 | 155s | 26 | 13c |
| recipe_sandals_journeyman | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Sandalet | 5 | 95s | 16 | 7c |
| recipe_trousers_journeyman | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Şalvar | 5 | 135s | 22 | 12c |
| recipe_wraps_journeyman | 1 | 4× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Sargı | 5 | 75s | 13 | 7c |
| recipe_tome_journeyman | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kalfa Kutsal Kitap | 5 | 155s | 26 | 16c |
| recipe_hood_fine | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Kukuleta | 10 | 95s | 16 | 20c |
| recipe_robe_fine | 1 | 10× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Cübbe | 10 | 155s | 26 | 33c |
| recipe_sandals_fine | 1 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Sandalet | 10 | 95s | 16 | 18c |
| recipe_trousers_fine | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Şalvar | 10 | 135s | 22 | 29c |
| recipe_wraps_fine | 1 | 4× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Sargı | 10 | 75s | 13 | 16c |
| recipe_tome_fine | 1 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kaliteli Kutsal Kitap | 10 | 155s | 26 | 41c |
| recipe_hood_veteran | 2 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Kukuleta | 15 | 100s | 17 | 36c |
| recipe_robe_veteran | 2 | 10× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Cübbe | 15 | 165s | 27 | 58c |
| recipe_sandals_veteran | 2 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Sandalet | 15 | 100s | 17 | 33c |
| recipe_trousers_veteran | 2 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Şalvar | 15 | 145s | 23 | 51c |
| recipe_wraps_veteran | 2 | 4× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Sargı | 15 | 80s | 13 | 29c |
| recipe_tome_veteran | 2 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Kıdemli Kutsal Kitap | 15 | 165s | 27 | 73c |
| recipe_hood_champion | 2 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Kukuleta | 20 | 100s | 17 | 55c |
| recipe_robe_champion | 2 | 10× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Cübbe | 20 | 165s | 27 | 88c |
| recipe_sandals_champion | 2 | 6× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Sandalet | 20 | 100s | 17 | 50c |
| recipe_trousers_champion | 2 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Şalvar | 20 | 145s | 23 | 77c |
| recipe_wraps_champion | 2 | 4× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Sargı | 20 | 80s | 13 | 44c |
| recipe_tome_champion | 2 | 8× Sıradan Ot, 1× Orman Faresi Postu | 1× Şampiyon Kutsal Kitap | 20 | 165s | 27 | 1s 10c |
| recipe_hood_runed | 3 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Kukuleta | 25 | 110s | 18 | 76c |
| recipe_robe_runed | 3 | 10× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Cübbe | 25 | 175s | 28 | 1s 22c |
| recipe_sandals_runed | 3 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Sandalet | 25 | 110s | 18 | 69c |
| recipe_trousers_runed | 3 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Şalvar | 25 | 150s | 25 | 1s 7c |
| recipe_wraps_runed | 3 | 4× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Sargı | 25 | 85s | 14 | 61c |
| recipe_tome_runed | 3 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Rünlü Kutsal Kitap | 25 | 175s | 28 | 1s 53c |
| recipe_hood_mithril | 4 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Kukuleta | 35 | 115s | 19 | 1s 19c |
| recipe_robe_mithril | 4 | 10× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Cübbe | 35 | 180s | 30 | 1s 90c |
| recipe_sandals_mithril | 4 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Sandalet | 35 | 115s | 19 | 1s 7c |
| recipe_trousers_mithril | 4 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Şalvar | 35 | 160s | 27 | 1s 66c |
| recipe_wraps_mithril | 4 | 4× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Sargı | 35 | 90s | 15 | 95c |
| recipe_tome_mithril | 4 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Mithril Kutsal Kitap | 35 | 180s | 30 | 2s 37c |
| recipe_hood_adamant | 5 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Kukuleta | 50 | 120s | 20 | 1s 88c |
| recipe_robe_adamant | 5 | 10× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Cübbe | 50 | 190s | 32 | 3s 1c |
| recipe_sandals_adamant | 5 | 6× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Sandalet | 50 | 120s | 20 | 1s 69c |
| recipe_trousers_adamant | 5 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Şalvar | 50 | 170s | 28 | 2s 64c |
| recipe_wraps_adamant | 5 | 4× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Sargı | 50 | 95s | 16 | 1s 51c |
| recipe_tome_adamant | 5 | 8× Karahindiba, 1× Küçük Tavşan Postu | 1× Adamant Kutsal Kitap | 50 | 190s | 32 | 3s 76c |
| recipe_hood_wyrmforged | 7 | 6× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Kukuleta | 65 | 130s | 22 | 2s 65c |
| recipe_robe_wyrmforged | 7 | 10× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Cübbe | 65 | 210s | 36 | 4s 24c |
| recipe_sandals_wyrmforged | 7 | 6× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Sandalet | 65 | 130s | 22 | 2s 38c |
| recipe_trousers_wyrmforged | 7 | 8× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Şalvar | 65 | 185s | 31 | 3s 71c |
| recipe_wraps_wyrmforged | 7 | 4× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Sargı | 65 | 105s | 18 | 2s 12c |
| recipe_tome_wyrmforged | 7 | 8× Yabani Nane, 1× Tarla Faresi Postu | 1× Ejder Dövme Kutsal Kitap | 65 | 210s | 36 | 5s 29c |
| recipe_hood_starforged | 8 | 6× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Kukuleta | 80 | 140s | 24 | 3s 44c |
| recipe_robe_starforged | 8 | 10× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Cübbe | 80 | 220s | 38 | 5s 51c |
| recipe_sandals_starforged | 8 | 6× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Sandalet | 80 | 140s | 24 | 3s 10c |
| recipe_trousers_starforged | 8 | 8× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Şalvar | 80 | 195s | 34 | 4s 82c |
| recipe_wraps_starforged | 8 | 4× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Sargı | 80 | 110s | 19 | 2s 75c |
| recipe_tome_starforged | 8 | 8× Yabani Nane, 1× Tarla Faresi Postu | 1× Yıldız Dövme Kutsal Kitap | 80 | 220s | 38 | 6s 88c |
| recipe_hood_titan | 10 | 6× Papatya, 1× Sincap Postu | 1× Titan Kukuleta | 95 | 150s | 27 | 4s 33c |
| recipe_robe_titan | 10 | 10× Papatya, 1× Sincap Postu | 1× Titan Cübbe | 95 | 240s | 43 | 6s 93c |
| recipe_sandals_titan | 10 | 6× Papatya, 1× Sincap Postu | 1× Titan Sandalet | 95 | 150s | 27 | 3s 89c |
| recipe_trousers_titan | 10 | 8× Papatya, 1× Sincap Postu | 1× Titan Şalvar | 95 | 210s | 37 | 6s 6c |
| recipe_wraps_titan | 10 | 4× Papatya, 1× Sincap Postu | 1× Titan Sargı | 95 | 120s | 21 | 3s 46c |
| recipe_tome_titan | 10 | 8× Papatya, 1× Sincap Postu | 1× Titan Kutsal Kitap | 95 | 240s | 43 | 8s 66c |
| recipe_hood_royal | 11 | 6× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Kukuleta | 110 | 155s | 28 | 5s 25c |
| recipe_robe_royal | 11 | 10× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Cübbe | 110 | 250s | 46 | 8s 40c |
| recipe_sandals_royal | 11 | 6× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Sandalet | 110 | 155s | 28 | 4s 72c |
| recipe_trousers_royal | 11 | 8× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Şalvar | 110 | 220s | 40 | 7s 35c |
| recipe_wraps_royal | 11 | 4× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Sargı | 110 | 125s | 23 | 4s 20c |
| recipe_tome_royal | 11 | 8× Lavanta, 1× Yaban Tavşanı Postu | 1× Kraliyet Kutsal Kitap | 110 | 250s | 46 | 10s 50c |
| recipe_hood_sunforged | 14 | 6× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Kukuleta | 135 | 175s | 34 | 6s 85c |
| recipe_robe_sunforged | 14 | 10× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Cübbe | 135 | 280s | 54 | 10s 96c |
| recipe_sandals_sunforged | 14 | 6× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Sandalet | 135 | 175s | 34 | 6s 16c |
| recipe_trousers_sunforged | 14 | 8× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Şalvar | 135 | 245s | 48 | 9s 59c |
| recipe_wraps_sunforged | 14 | 4× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Sargı | 135 | 140s | 27 | 5s 48c |
| recipe_tome_sunforged | 14 | 8× Ayçiçeği, 1× Porsuk Postu | 1× Güneş Dövme Kutsal Kitap | 135 | 280s | 54 | 13s 70c |
| recipe_hood_obsidian | 16 | 6× Acıkök, 1× Porsuk Postu | 1× Obsidyen Kukuleta | 160 | 185s | 38 | 8s 54c |
| recipe_robe_obsidian | 16 | 10× Acıkök, 1× Porsuk Postu | 1× Obsidyen Cübbe | 160 | 300s | 61 | 13s 67c |
| recipe_sandals_obsidian | 16 | 6× Acıkök, 1× Porsuk Postu | 1× Obsidyen Sandalet | 160 | 185s | 38 | 7s 68c |
| recipe_trousers_obsidian | 16 | 8× Acıkök, 1× Porsuk Postu | 1× Obsidyen Şalvar | 160 | 260s | 53 | 11s 96c |
| recipe_wraps_obsidian | 16 | 4× Acıkök, 1× Porsuk Postu | 1× Obsidyen Sargı | 160 | 150s | 30 | 6s 84c |
| recipe_tome_obsidian | 16 | 8× Acıkök, 1× Porsuk Postu | 1× Obsidyen Kutsal Kitap | 160 | 300s | 61 | 17s 9c |
| recipe_hood_sandsteel | 19 | 6× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Kukuleta | 185 | 205s | 46 | 10s 32c |
| recipe_robe_sandsteel | 19 | 10× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Cübbe | 185 | 325s | 73 | 16s 51c |
| recipe_sandals_sandsteel | 19 | 6× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Sandalet | 185 | 205s | 46 | 9s 28c |
| recipe_trousers_sandsteel | 19 | 8× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Şalvar | 185 | 285s | 64 | 14s 45c |
| recipe_wraps_sandsteel | 19 | 4× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Sargı | 185 | 165s | 36 | 8s 26c |
| recipe_tome_sandsteel | 19 | 8× Bataklık Sazı, 1× Kızıl Tilki Postu | 1× Kum Çeliği Kutsal Kitap | 185 | 325s | 73 | 20s 64c |
| recipe_hood_runebound | 21 | 6× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Kukuleta | 210 | 215s | 51 | 12s 17c |
| recipe_robe_runebound | 21 | 10× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Cübbe | 210 | 345s | 82 | 19s 47c |
| recipe_sandals_runebound | 21 | 6× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Sandalet | 210 | 215s | 51 | 10s 94c |
| recipe_trousers_runebound | 21 | 8× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Şalvar | 210 | 300s | 71 | 17s 4c |
| recipe_wraps_runebound | 21 | 4× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Sargı | 210 | 175s | 41 | 9s 74c |
| recipe_tome_runebound | 21 | 8× Mavi Devedikeni, 1× Yaban Domuzu Postu | 1× Rün Bağlı Kutsal Kitap | 210 | 345s | 82 | 24s 34c |
| recipe_hood_phoenix | 24 | 6× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Kukuleta | 235 | 235s | 61 | 14s 8c |
| recipe_robe_phoenix | 24 | 10× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Cübbe | 235 | 375s | 97 | 22s 53c |
| recipe_sandals_phoenix | 24 | 6× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Sandalet | 235 | 235s | 61 | 12s 67c |
| recipe_trousers_phoenix | 24 | 8× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Şalvar | 235 | 330s | 85 | 19s 72c |
| recipe_wraps_phoenix | 24 | 4× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Sargı | 235 | 185s | 49 | 11s 27c |
| recipe_tome_phoenix | 24 | 8× Kızıl Adaçayı, 1× Dağ Keçisi Postu | 1× Anka Kutsal Kitap | 235 | 375s | 97 | 28s 17c |
| recipe_hood_dreadforged | 26 | 6× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Kukuleta | 260 | 245s | 68 | 16s 6c |
| recipe_robe_dreadforged | 26 | 10× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Cübbe | 260 | 395s | 109 | 25s 70c |
| recipe_sandals_dreadforged | 26 | 6× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Sandalet | 260 | 245s | 68 | 14s 44c |
| recipe_trousers_dreadforged | 26 | 8× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Şalvar | 260 | 345s | 96 | 22s 49c |
| recipe_wraps_dreadforged | 26 | 4× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Sargı | 260 | 195s | 55 | 12s 85c |
| recipe_tome_dreadforged | 26 | 8× Güzelavrat Otu, 1× Yaban Kurdu Postu | 1× Dehşet Dövme Kutsal Kitap | 260 | 395s | 109 | 32s 12c |
| recipe_hood_sovereign | 29 | 6× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Kukuleta | 285 | 265s | 81 | 18s 10c |
| recipe_robe_sovereign | 29 | 10× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Cübbe | 285 | 420s | 130 | 28s 96c |
| recipe_sandals_sovereign | 29 | 6× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Sandalet | 285 | 265s | 81 | 16s 28c |
| recipe_trousers_sovereign | 29 | 8× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Şalvar | 285 | 370s | 114 | 25s 34c |
| recipe_wraps_sovereign | 29 | 4× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Sargı | 285 | 210s | 65 | 14s 48c |
| recipe_tome_sovereign | 29 | 8× Altın Yonca, 1× Yaban Mandası Derisi | 1× Hükümdar Kutsal Kitap | 285 | 420s | 130 | 36s 20c |
| recipe_hood_coral | 31 | 6× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Kukuleta | 310 | 275s | 91 | 20s 19c |
| recipe_robe_coral | 31 | 10× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Cübbe | 310 | 440s | 146 | 32s 30c |
| recipe_sandals_coral | 31 | 6× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Sandalet | 310 | 275s | 91 | 18s 16c |
| recipe_trousers_coral | 31 | 8× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Şalvar | 310 | 385s | 128 | 28s 26c |
| recipe_wraps_coral | 31 | 4× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Sargı | 310 | 220s | 73 | 16s 15c |
| recipe_tome_coral | 31 | 8× Altın Yonca, 1× Yaban Mandası Derisi | 1× Mercan Kutsal Kitap | 310 | 440s | 146 | 40s 38c |
| recipe_hood_drakebone | 34 | 6× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Kukuleta | 335 | 295s | 109 | 22s 33c |
| recipe_robe_drakebone | 34 | 10× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Cübbe | 335 | 470s | 174 | 35s 73c |
| recipe_sandals_drakebone | 34 | 6× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Sandalet | 335 | 295s | 109 | 20s 8c |
| recipe_trousers_drakebone | 34 | 8× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Şalvar | 335 | 410s | 152 | 31s 26c |
| recipe_wraps_drakebone | 34 | 4× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Sargı | 335 | 235s | 87 | 17s 86c |
| recipe_tome_drakebone | 34 | 8× Düş Kökü, 1× Dev Geyik Postu | 1× Ejderha Kemiği Kutsal Kitap | 335 | 470s | 174 | 44s 66c |
| recipe_hood_mistwoven | 36 | 6× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Kukuleta | 360 | 305s | 122 | 24s 52c |
| recipe_robe_mistwoven | 36 | 10× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Cübbe | 360 | 490s | 196 | 39s 23c |
| recipe_sandals_mistwoven | 36 | 6× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Sandalet | 360 | 305s | 122 | 22s 5c |
| recipe_trousers_mistwoven | 36 | 8× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Şalvar | 360 | 430s | 171 | 34s 33c |
| recipe_wraps_mistwoven | 36 | 4× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Sargı | 360 | 245s | 98 | 19s 62c |
| recipe_tome_mistwoven | 36 | 8× Gümüş Yaprak, 1× Dev Kurt Postu | 1× Sis Örgü Kutsal Kitap | 360 | 490s | 196 | 49s 4c |
| recipe_hood_prismatic | 39 | 6× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Kukuleta | 385 | 325s | 146 | 26s 76c |
| recipe_robe_prismatic | 39 | 10× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Cübbe | 385 | 520s | 233 | 42s 81c |
| recipe_sandals_prismatic | 39 | 6× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Sandalet | 385 | 325s | 146 | 24s 6c |
| recipe_trousers_prismatic | 39 | 8× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Şalvar | 385 | 455s | 204 | 37s 46c |
| recipe_wraps_prismatic | 39 | 4× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Sargı | 385 | 260s | 116 | 21s 41c |
| recipe_tome_prismatic | 39 | 8× Kan Dikeni, 1× Dağ Aslanı Postu | 1× Prizmatik Kutsal Kitap | 385 | 520s | 233 | 53s 51c |
| recipe_hood_magmaforged | 41 | 6× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Kukuleta | 410 | 335s | 164 | 29s 4c |
| recipe_robe_magmaforged | 41 | 10× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Cübbe | 410 | 540s | 262 | 46s 46c |
| recipe_sandals_magmaforged | 41 | 6× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Sandalet | 410 | 335s | 164 | 26s 12c |
| recipe_trousers_magmaforged | 41 | 8× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Şalvar | 410 | 470s | 229 | 40s 65c |
| recipe_wraps_magmaforged | 41 | 4× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Sargı | 410 | 270s | 131 | 23s 23c |
| recipe_tome_magmaforged | 41 | 8× Hayalet Mantarı, 1× Kara Panter Postu | 1× Magma Dövme Kutsal Kitap | 410 | 540s | 262 | 58s 8c |
| recipe_hood_tempest | 44 | 6× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Kukuleta | 435 | 355s | 195 | 31s 36c |
| recipe_robe_tempest | 44 | 10× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Cübbe | 435 | 565s | 312 | 50s 18c |
| recipe_sandals_tempest | 44 | 6× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Sandalet | 435 | 355s | 195 | 28s 20c |
| recipe_trousers_tempest | 44 | 8× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Şalvar | 435 | 495s | 273 | 43s 90c |
| recipe_wraps_tempest | 44 | 4× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Sargı | 435 | 285s | 156 | 25s 9c |
| recipe_tome_tempest | 44 | 8× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Bora Kutsal Kitap | 435 | 565s | 312 | 62s 72c |
| recipe_hood_voidtouched | 46 | 6× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Kukuleta | 460 | 365s | 219 | 33s 72c |
| recipe_robe_voidtouched | 46 | 10× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Cübbe | 460 | 585s | 350 | 53s 96c |
| recipe_sandals_voidtouched | 46 | 6× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Sandalet | 460 | 365s | 219 | 30s 33c |
| recipe_trousers_voidtouched | 46 | 8× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Şalvar | 460 | 510s | 306 | 47s 21c |
| recipe_wraps_voidtouched | 46 | 4× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Sargı | 460 | 295s | 175 | 26s 98c |
| recipe_tome_voidtouched | 46 | 8× Fırtına Çiçeği, 1× Kara Panter Postu | 1× Hiçlik Kutsal Kitap | 460 | 585s | 350 | 67s 45c |
| recipe_hood_seraphic | 49 | 6× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Kukuleta | 485 | 385s | 261 | 36s 12c |
| recipe_robe_seraphic | 49 | 10× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Cübbe | 485 | 615s | 417 | 57s 80c |
| recipe_sandals_seraphic | 49 | 6× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Sandalet | 485 | 385s | 261 | 32s 49c |
| recipe_trousers_seraphic | 49 | 8× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Şalvar | 485 | 540s | 365 | 50s 58c |
| recipe_wraps_seraphic | 49 | 4× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Sargı | 485 | 305s | 209 | 28s 90c |
| recipe_tome_seraphic | 49 | 8× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Melek Kutsal Kitap | 485 | 615s | 417 | 72s 25c |
| recipe_hood_eclipse | 51 | 6× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Kukuleta | 510 | 395s | 293 | 38s 56c |
| recipe_robe_eclipse | 51 | 10× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Cübbe | 510 | 635s | 469 | 61s 70c |
| recipe_sandals_eclipse | 51 | 6× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Sandalet | 510 | 395s | 293 | 34s 68c |
| recipe_trousers_eclipse | 51 | 8× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Şalvar | 510 | 555s | 410 | 53s 99c |
| recipe_wraps_eclipse | 51 | 4× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Sargı | 510 | 315s | 234 | 30s 85c |
| recipe_tome_eclipse | 51 | 8× Kor Yaprağı, 1× Kor Domuzu Postu | 1× Tutulma Kutsal Kitap | 510 | 635s | 469 | 77s 13c |
| recipe_hood_wyrmlord | 54 | 6× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Kukuleta | 535 | 415s | 349 | 41s 4c |
| recipe_robe_wyrmlord | 54 | 10× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Cübbe | 535 | 660s | 558 | 65s 66c |
| recipe_sandals_wyrmlord | 54 | 6× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Sandalet | 535 | 415s | 349 | 36s 91c |
| recipe_trousers_wyrmlord | 54 | 8× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Şalvar | 535 | 580s | 488 | 57s 46c |
| recipe_wraps_wyrmlord | 54 | 4× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Sargı | 535 | 330s | 279 | 32s 83c |
| recipe_tome_wyrmlord | 54 | 8× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Ejder Lordu Kutsal Kitap | 535 | 660s | 558 | 82s 8c |
| recipe_hood_celestine | 56 | 6× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Kukuleta | 560 | 425s | 392 | 43s 55c |
| recipe_robe_celestine | 56 | 10× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Cübbe | 560 | 680s | 627 | 69s 68c |
| recipe_sandals_celestine | 56 | 6× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Sandalet | 560 | 425s | 392 | 39s 17c |
| recipe_trousers_celestine | 56 | 8× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Şalvar | 560 | 595s | 549 | 60s 97c |
| recipe_wraps_celestine | 56 | 4× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Sargı | 560 | 340s | 314 | 34s 84c |
| recipe_tome_celestine | 56 | 8× Ayaz Çiçeği, 1× Ayaz Kurdu Postu | 1× Göksel Kutsal Kitap | 560 | 680s | 627 | 87s 10c |
| recipe_hood_eternal | 59 | 6× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Kukuleta | 585 | 445s | 467 | 46s 9c |
| recipe_robe_eternal | 59 | 10× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Cübbe | 585 | 710s | 747 | 73s 75c |
| recipe_sandals_eternal | 59 | 6× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Sandalet | 585 | 445s | 467 | 41s 45c |
| recipe_trousers_eternal | 59 | 8× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Şalvar | 585 | 620s | 653 | 64s 53c |
| recipe_wraps_eternal | 59 | 4× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Sargı | 585 | 355s | 373 | 36s 88c |
| recipe_tome_eternal | 59 | 8× Ay Gölgesi, 1× Dev Ayı Postu | 1× Ebedi Kutsal Kitap | 585 | 710s | 747 | 92s 19c |
| recipe_hood_starless | 61 | 6× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Kukuleta | 610 | 455s | 524 | 48s 67c |
| recipe_robe_starless | 61 | 10× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Cübbe | 610 | 730s | 839 | 77s 87c |
| recipe_sandals_starless | 61 | 6× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Sandalet | 610 | 455s | 524 | 43s 77c |
| recipe_trousers_starless | 61 | 8× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Şalvar | 610 | 640s | 734 | 68s 14c |
| recipe_wraps_starless | 61 | 4× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Sargı | 610 | 365s | 420 | 38s 94c |
| recipe_tome_starless | 61 | 8× Ay Gölgesi, 1× Dev Ayı Postu | 1× Yıldızsız Kutsal Kitap | 610 | 730s | 839 | 97s 34c |
| recipe_hood_nightsteel | 64 | 6× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Kukuleta | 635 | 475s | 625 | 51s 28c |
| recipe_robe_nightsteel | 64 | 10× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Cübbe | 635 | 760s | 1000 | 82s 5c |
| recipe_sandals_nightsteel | 64 | 6× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Sandalet | 635 | 475s | 625 | 46s 12c |
| recipe_trousers_nightsteel | 64 | 8× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Şalvar | 635 | 665s | 875 | 71s 79c |
| recipe_wraps_nightsteel | 64 | 4× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Sargı | 635 | 380s | 500 | 41s 2c |
| recipe_tome_nightsteel | 64 | 8× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Gece Çeliği Kutsal Kitap | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_hood_chronoforged | 66 | 6× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Kukuleta | 660 | 485s | 702 | 53s 92c |
| recipe_robe_chronoforged | 66 | 10× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Cübbe | 660 | 780s | 1123 | 86s 27c |
| recipe_sandals_chronoforged | 66 | 6× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Sandalet | 660 | 485s | 702 | 48s 49c |
| recipe_trousers_chronoforged | 66 | 8× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Şalvar | 660 | 680s | 983 | 75s 49c |
| recipe_wraps_chronoforged | 66 | 4× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Sargı | 660 | 390s | 562 | 43s 14c |
| recipe_tome_chronoforged | 66 | 8× Ruh Nilüferi, 1× Ruh Geyiği Postu | 1× Zaman Dövme Kutsal Kitap | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_hood_ascendant | 69 | 6× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Kukuleta | 685 | 505s | 836 | 56s 59c |
| recipe_robe_ascendant | 69 | 10× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Cübbe | 685 | 805s | 1338 | 90s 54c |
| recipe_sandals_ascendant | 69 | 6× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Sandalet | 685 | 505s | 836 | 50s 89c |
| recipe_trousers_ascendant | 69 | 8× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Şalvar | 685 | 705s | 1170 | 79s 23c |
| recipe_wraps_ascendant | 69 | 4× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Sargı | 685 | 405s | 669 | 45s 27c |
| recipe_tome_ascendant | 69 | 8× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Yükselmiş Kutsal Kitap | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_hood_stormborn | 71 | 6× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Kukuleta | 710 | 515s | 939 | 59s 29c |
| recipe_robe_stormborn | 71 | 10× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Cübbe | 710 | 825s | 1503 | 94s 86c |
| recipe_sandals_stormborn | 71 | 6× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Sandalet | 710 | 515s | 939 | 53s 32c |
| recipe_trousers_stormborn | 71 | 8× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Şalvar | 710 | 720s | 1315 | 83s |
| recipe_wraps_stormborn | 71 | 4× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Sargı | 710 | 415s | 751 | 47s 43c |
| recipe_tome_stormborn | 71 | 8× Kadim Adamotu, 1× Kadim Mamut Postu | 1× Fırtına Doğumlu Kutsal Kitap | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_hood_riftforged | 74 | 6× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Kukuleta | 735 | 535s | 1119 | 62s 2c |
| recipe_robe_riftforged | 74 | 10× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Cübbe | 735 | 855s | 1790 | 99s 23c |
| recipe_sandals_riftforged | 74 | 6× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Sandalet | 735 | 535s | 1119 | 55s 78c |
| recipe_trousers_riftforged | 74 | 8× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Şalvar | 735 | 750s | 1566 | 86s 82c |
| recipe_wraps_riftforged | 74 | 4× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Sargı | 735 | 425s | 895 | 49s 61c |
| recipe_tome_riftforged | 74 | 8× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Yarık Dövme Kutsal Kitap | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_hood_chaos | 76 | 6× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Kukuleta | 760 | 545s | 1257 | 64s 77c |
| recipe_robe_chaos | 76 | 10× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Cübbe | 760 | 875s | 2011 | 1g 3s 64c |
| recipe_sandals_chaos | 76 | 6× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Sandalet | 760 | 545s | 1257 | 58s 25c |
| recipe_trousers_chaos | 76 | 8× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Şalvar | 760 | 765s | 1760 | 90s 68c |
| recipe_wraps_chaos | 76 | 4× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Sargı | 760 | 435s | 1006 | 51s 82c |
| recipe_tome_chaos | 76 | 8× Ejderkanı Otu, 1× Gölge Kimera Derisi | 1× Kaos Kutsal Kitap | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_hood_skyforged | 79 | 6× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Kukuleta | 785 | 565s | 1497 | 67s 56c |
| recipe_robe_skyforged | 79 | 10× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Cübbe | 785 | 900s | 2395 | 1g 8s 9c |
| recipe_sandals_skyforged | 79 | 6× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Sandalet | 785 | 565s | 1497 | 60s 76c |
| recipe_trousers_skyforged | 79 | 8× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Şalvar | 785 | 790s | 2096 | 94s 58c |
| recipe_wraps_skyforged | 79 | 4× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Sargı | 785 | 450s | 1198 | 54s 5c |
| recipe_tome_skyforged | 79 | 8× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Gök Dövme Kutsal Kitap | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_hood_archangel | 81 | 6× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Kukuleta | 810 | 575s | 1682 | 70s 37c |
| recipe_robe_archangel | 81 | 10× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Cübbe | 810 | 920s | 2692 | 1g 12s 59c |
| recipe_sandals_archangel | 81 | 6× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Sandalet | 810 | 575s | 1682 | 63s 28c |
| recipe_trousers_archangel | 81 | 8× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Şalvar | 810 | 805s | 2355 | 98s 51c |
| recipe_wraps_archangel | 81 | 4× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Sargı | 810 | 460s | 1346 | 56s 29c |
| recipe_tome_archangel | 81 | 8× Yıldız Çiçeği, 1× Yaşlı Wyvern Pulu | 1× Başmelek Kutsal Kitap | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_hood_aurora | 84 | 6× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Kukuleta | 835 | 595s | 2004 | 73s 20c |
| recipe_robe_aurora | 84 | 10× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Cübbe | 835 | 950s | 3206 | 1g 17s 13c |
| recipe_sandals_aurora | 84 | 6× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Sandalet | 835 | 595s | 2004 | 65s 84c |
| recipe_trousers_aurora | 84 | 8× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Şalvar | 835 | 830s | 2805 | 1g 2s 48c |
| recipe_wraps_aurora | 84 | 4× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Sargı | 835 | 475s | 1603 | 58s 56c |
| recipe_tome_aurora | 84 | 8× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Şafak Kutsal Kitap | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_hood_ashen | 86 | 6× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Kukuleta | 860 | 605s | 2251 | 76s 6c |
| recipe_robe_ashen | 86 | 10× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Cübbe | 860 | 970s | 3602 | 1g 21s 70c |
| recipe_sandals_ashen | 86 | 6× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Sandalet | 860 | 605s | 2251 | 68s 41c |
| recipe_trousers_ashen | 86 | 8× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Şalvar | 860 | 850s | 3151 | 1g 6s 49c |
| recipe_wraps_ashen | 86 | 4× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Sargı | 860 | 485s | 1801 | 60s 85c |
| recipe_tome_ashen | 86 | 8× Boşluk Orkidesi, 1× Yaşlı Wyvern Pulu | 1× Kül Kutsal Kitap | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_hood_leviathan | 89 | 6× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Kukuleta | 885 | 625s | 2681 | 78s 95c |
| recipe_robe_leviathan | 89 | 10× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Cübbe | 885 | 1000s | 4290 | 1g 26s 32c |
| recipe_sandals_leviathan | 89 | 6× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Sandalet | 885 | 625s | 2681 | 71s 1c |
| recipe_trousers_leviathan | 89 | 8× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Şalvar | 885 | 875s | 3753 | 1g 10s 53c |
| recipe_wraps_leviathan | 89 | 4× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Sargı | 885 | 500s | 2145 | 63s 16c |
| recipe_tome_leviathan | 89 | 8× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Leviathan Kutsal Kitap | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_hood_doomforged | 91 | 6× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Kukuleta | 910 | 635s | 3012 | 81s 86c |
| recipe_robe_doomforged | 91 | 10× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Cübbe | 910 | 1020s | 4820 | 1g 30s 98c |
| recipe_sandals_doomforged | 91 | 6× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Sandalet | 910 | 635s | 3012 | 73s 62c |
| recipe_trousers_doomforged | 91 | 8× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Şalvar | 910 | 890s | 4218 | 1g 14s 61c |
| recipe_wraps_doomforged | 91 | 4× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Sargı | 910 | 510s | 2410 | 65s 49c |
| recipe_tome_doomforged | 91 | 8× Göksel Nilüfer, 1× Boşluk Behemotu Derisi | 1× Kıyamet Dövme Kutsal Kitap | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_hood_voidforged | 94 | 6× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Kukuleta | 935 | 655s | 3588 | 84s 80c |
| recipe_robe_voidforged | 94 | 10× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Cübbe | 935 | 1045s | 5741 | 1g 35s 68c |
| recipe_sandals_voidforged | 94 | 6× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Sandalet | 935 | 655s | 3588 | 76s 26c |
| recipe_trousers_voidforged | 94 | 8× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Şalvar | 935 | 915s | 5023 | 1g 18s 72c |
| recipe_wraps_voidforged | 94 | 4× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Sargı | 935 | 525s | 2870 | 67s 84c |
| recipe_tome_voidforged | 94 | 8× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Hiçlik Dövme Kutsal Kitap | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_hood_godslayer | 96 | 6× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Kukuleta | 960 | 665s | 4032 | 87s 76c |
| recipe_robe_godslayer | 96 | 10× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Cübbe | 960 | 1065s | 6450 | 1g 40s 42c |
| recipe_sandals_godslayer | 96 | 6× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Sandalet | 960 | 665s | 4032 | 78s 93c |
| recipe_trousers_godslayer | 96 | 8× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Şalvar | 960 | 930s | 5644 | 1g 22s 86c |
| recipe_wraps_godslayer | 96 | 4× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Sargı | 960 | 535s | 3225 | 70s 21c |
| recipe_tome_godslayer | 96 | 8× Anka Çiçeği, 1× Boşluk Behemotu Derisi | 1× Tanrı Katili Kutsal Kitap | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_hood_primordial | 99 | 6× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Kukuleta | 985 | 685s | 4802 | 90s 74c |
| recipe_robe_primordial | 99 | 10× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Cübbe | 985 | 1095s | 7682 | 1g 45s 19c |
| recipe_sandals_primordial | 99 | 6× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Sandalet | 985 | 685s | 4802 | 81s 61c |
| recipe_trousers_primordial | 99 | 8× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Şalvar | 985 | 960s | 6722 | 1g 27s 4c |
| recipe_wraps_primordial | 99 | 4× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Sargı | 985 | 545s | 3841 | 72s 59c |
| recipe_tome_primordial | 99 | 8× Dünya Kökü, 1× Göksel Ejder Pulu | 1× Kadim Çağ Kutsal Kitap | 985 | 1095s | 7682 | 1g 81s 48c |

### Ağaç İşçiliği — 376 tarif, beceri 1–99 (tavan 100), çıktı Sv 1–985

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_bow_short | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kısa Yay | 1 | 155s | 26 | 4c |
| recipe_crossbow_novice | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Acemi Arbalet | 1 | 155s | 26 | 4c |
| recipe_lute_simple | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Sade Lavta | 1 | 155s | 26 | 4c |
| recipe_focus_novice | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Acemi Doğa Odağı | 1 | 155s | 26 | 4c |
| recipe_staff_willow | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Söğüt Asa | 1 | 155s | 26 | 4c |
| recipe_tonfa_novice | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Acemi Tonfa | 1 | 155s | 26 | 4c |
| recipe_totem_carved | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Oyma Totem | 1 | 155s | 26 | 4c |
| recipe_wand_novice | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Acemi Değnek | 1 | 155s | 26 | 4c |
| recipe_bow_hunting | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Av Yayı | 5 | 155s | 26 | 16c |
| recipe_crossbow_journeyman | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kalfa Arbalet | 5 | 155s | 26 | 16c |
| recipe_lute_maple | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Akçaağaç Lavta | 5 | 155s | 26 | 16c |
| recipe_focus_journeyman | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kalfa Doğa Odağı | 5 | 155s | 26 | 16c |
| recipe_staff_apprentice | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Çırak Asası | 5 | 155s | 26 | 16c |
| recipe_tonfa_journeyman | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kalfa Tonfa | 5 | 155s | 26 | 16c |
| recipe_totem_bone | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kemik Totem | 5 | 155s | 26 | 16c |
| recipe_wand_journeyman | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kalfa Değnek | 5 | 155s | 26 | 16c |
| recipe_bow_yew | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Porsuk Yayı | 10 | 155s | 26 | 41c |
| recipe_crossbow_fine | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kaliteli Arbalet | 10 | 155s | 26 | 41c |
| recipe_lute_silver | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Gümüş Telli Lavta | 10 | 155s | 26 | 41c |
| recipe_focus_fine | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kaliteli Doğa Odağı | 10 | 155s | 26 | 41c |
| recipe_staff_adept | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Usta Asası | 10 | 155s | 26 | 41c |
| recipe_tonfa_fine | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kaliteli Tonfa | 10 | 155s | 26 | 41c |
| recipe_totem_spirit | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Ruh Totemi | 10 | 155s | 26 | 41c |
| recipe_wand_fine | 1 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kaliteli Değnek | 10 | 155s | 26 | 41c |
| recipe_bow_composite | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Bileşik Yay | 15 | 165s | 27 | 73c |
| recipe_crossbow_veteran | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kıdemli Arbalet | 15 | 165s | 27 | 73c |
| recipe_lute_concert | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Konser Lavtası | 15 | 165s | 27 | 73c |
| recipe_focus_veteran | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kıdemli Doğa Odağı | 15 | 165s | 27 | 73c |
| recipe_staff_crystal | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kristal Asa | 15 | 165s | 27 | 73c |
| recipe_tonfa_veteran | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kıdemli Tonfa | 15 | 165s | 27 | 73c |
| recipe_totem_stone | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Taş Put | 15 | 165s | 27 | 73c |
| recipe_wand_veteran | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Kıdemli Değnek | 15 | 165s | 27 | 73c |
| recipe_bow_longbow | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Uzun Yay | 20 | 165s | 27 | 1s 10c |
| recipe_crossbow_champion | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Şampiyon Arbalet | 20 | 165s | 27 | 1s 10c |
| recipe_lute_enchanted | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Büyülü Lavta | 20 | 165s | 27 | 1s 10c |
| recipe_focus_champion | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Şampiyon Doğa Odağı | 20 | 165s | 27 | 1s 10c |
| recipe_staff_arcane | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Gizemli Asa | 20 | 165s | 27 | 1s 10c |
| recipe_tonfa_champion | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Şampiyon Tonfa | 20 | 165s | 27 | 1s 10c |
| recipe_totem_thunder | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Gök Gürültüsü Putu | 20 | 165s | 27 | 1s 10c |
| recipe_wand_champion | 2 | 4× Fidan Dalı, 1× Alüminyum Külçesi | 1× Şampiyon Değnek | 20 | 165s | 27 | 1s 10c |
| recipe_bow_dragonbone | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Ejder Kemiği Yayı | 25 | 175s | 28 | 1s 53c |
| recipe_crossbow_runed | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Rünlü Arbalet | 25 | 175s | 28 | 1s 53c |
| recipe_lute_legend | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Efsaneler Lavtası | 25 | 175s | 28 | 1s 53c |
| recipe_focus_runed | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Rünlü Doğa Odağı | 25 | 175s | 28 | 1s 53c |
| recipe_staff_astral | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Asası | 25 | 175s | 28 | 1s 53c |
| recipe_tonfa_runed | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Rünlü Tonfa | 25 | 175s | 28 | 1s 53c |
| recipe_totem_ancestral | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Ata Totemi | 25 | 175s | 28 | 1s 53c |
| recipe_wand_runed | 3 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Rünlü Değnek | 25 | 175s | 28 | 1s 53c |
| recipe_bow_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Mithril Yay | 35 | 180s | 30 | 2s 37c |
| recipe_crossbow_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Mithril Arbalet | 35 | 180s | 30 | 2s 37c |
| recipe_lute_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Lavta | 35 | 180s | 30 | 2s 37c |
| recipe_focus_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Doğa Odağı | 35 | 180s | 30 | 2s 37c |
| recipe_staff_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Asa | 35 | 180s | 30 | 2s 37c |
| recipe_tonfa_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Tonfa | 35 | 180s | 30 | 2s 37c |
| recipe_totem_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Totem | 35 | 180s | 30 | 2s 37c |
| recipe_wand_mithril | 4 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Mithril Değnek | 35 | 180s | 30 | 2s 37c |
| recipe_bow_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Adamant Yay | 50 | 190s | 32 | 3s 76c |
| recipe_crossbow_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Adamant Arbalet | 50 | 190s | 32 | 3s 76c |
| recipe_lute_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Lavta | 50 | 190s | 32 | 3s 76c |
| recipe_focus_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Doğa Odağı | 50 | 190s | 32 | 3s 76c |
| recipe_staff_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Asa | 50 | 190s | 32 | 3s 76c |
| recipe_tonfa_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Tonfa | 50 | 190s | 32 | 3s 76c |
| recipe_totem_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Totem | 50 | 190s | 32 | 3s 76c |
| recipe_wand_adamant | 5 | 4× Kavak Kütüğü, 1× Alüminyum Külçesi | 1× Adamant Değnek | 50 | 190s | 32 | 3s 76c |
| recipe_bow_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Ejder Dövme Yay | 65 | 210s | 36 | 5s 29c |
| recipe_crossbow_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Ejder Dövme Arbalet | 65 | 210s | 36 | 5s 29c |
| recipe_lute_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Lavta | 65 | 210s | 36 | 5s 29c |
| recipe_focus_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Doğa Odağı | 65 | 210s | 36 | 5s 29c |
| recipe_staff_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Asa | 65 | 210s | 36 | 5s 29c |
| recipe_tonfa_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Tonfa | 65 | 210s | 36 | 5s 29c |
| recipe_totem_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Totem | 65 | 210s | 36 | 5s 29c |
| recipe_wand_wyrmforged | 7 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Ejder Dövme Değnek | 65 | 210s | 36 | 5s 29c |
| recipe_bow_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Yıldız Dövme Yay | 80 | 220s | 38 | 6s 88c |
| recipe_crossbow_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Yıldız Dövme Arbalet | 80 | 220s | 38 | 6s 88c |
| recipe_lute_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Lavta | 80 | 220s | 38 | 6s 88c |
| recipe_focus_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Doğa Odağı | 80 | 220s | 38 | 6s 88c |
| recipe_staff_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Asa | 80 | 220s | 38 | 6s 88c |
| recipe_tonfa_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Tonfa | 80 | 220s | 38 | 6s 88c |
| recipe_totem_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Totem | 80 | 220s | 38 | 6s 88c |
| recipe_wand_starforged | 8 | 4× Söğüt Kütüğü, 1× Alüminyum Külçesi | 1× Yıldız Dövme Değnek | 80 | 220s | 38 | 6s 88c |
| recipe_bow_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Titan Yay | 95 | 240s | 43 | 8s 66c |
| recipe_crossbow_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Titan Arbalet | 95 | 240s | 43 | 8s 66c |
| recipe_lute_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Lavta | 95 | 240s | 43 | 8s 66c |
| recipe_focus_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Doğa Odağı | 95 | 240s | 43 | 8s 66c |
| recipe_staff_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Asa | 95 | 240s | 43 | 8s 66c |
| recipe_tonfa_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Tonfa | 95 | 240s | 43 | 8s 66c |
| recipe_totem_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Totem | 95 | 240s | 43 | 8s 66c |
| recipe_wand_titan | 10 | 4× Köknar Kütüğü, 1× Alüminyum Külçesi | 1× Titan Değnek | 95 | 240s | 43 | 8s 66c |
| recipe_bow_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kraliyet Yay | 110 | 250s | 46 | 10s 50c |
| recipe_crossbow_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi, 2× Yaban Hindisi Tüyü | 1× Kraliyet Arbalet | 110 | 250s | 46 | 10s 50c |
| recipe_lute_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Lavta | 110 | 250s | 46 | 10s 50c |
| recipe_focus_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Doğa Odağı | 110 | 250s | 46 | 10s 50c |
| recipe_staff_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Asa | 110 | 250s | 46 | 10s 50c |
| recipe_tonfa_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Tonfa | 110 | 250s | 46 | 10s 50c |
| recipe_totem_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Totem | 110 | 250s | 46 | 10s 50c |
| recipe_wand_royal | 11 | 4× Çam Kütüğü, 1× Alüminyum Külçesi | 1× Kraliyet Değnek | 110 | 250s | 46 | 10s 50c |
| recipe_bow_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi, 2× Yaban Hindisi Tüyü | 1× Güneş Dövme Yay | 135 | 280s | 54 | 13s 70c |
| recipe_crossbow_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi, 2× Yaban Hindisi Tüyü | 1× Güneş Dövme Arbalet | 135 | 280s | 54 | 13s 70c |
| recipe_lute_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Lavta | 135 | 280s | 54 | 13s 70c |
| recipe_focus_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Doğa Odağı | 135 | 280s | 54 | 13s 70c |
| recipe_staff_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Asa | 135 | 280s | 54 | 13s 70c |
| recipe_tonfa_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Tonfa | 135 | 280s | 54 | 13s 70c |
| recipe_totem_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Totem | 135 | 280s | 54 | 13s 70c |
| recipe_wand_sunforged | 14 | 4× Huş Kütüğü, 1× Kurşun Külçesi | 1× Güneş Dövme Değnek | 135 | 280s | 54 | 13s 70c |
| recipe_bow_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi, 2× Yaban Hindisi Tüyü | 1× Obsidyen Yay | 160 | 300s | 61 | 17s 9c |
| recipe_crossbow_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi, 2× Yaban Hindisi Tüyü | 1× Obsidyen Arbalet | 160 | 300s | 61 | 17s 9c |
| recipe_lute_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Lavta | 160 | 300s | 61 | 17s 9c |
| recipe_focus_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Doğa Odağı | 160 | 300s | 61 | 17s 9c |
| recipe_staff_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Asa | 160 | 300s | 61 | 17s 9c |
| recipe_tonfa_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Tonfa | 160 | 300s | 61 | 17s 9c |
| recipe_totem_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Totem | 160 | 300s | 61 | 17s 9c |
| recipe_wand_obsidian | 16 | 4× Sedir Kütüğü, 1× Demir Külçesi | 1× Obsidyen Değnek | 160 | 300s | 61 | 17s 9c |
| recipe_bow_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi, 2× Yaban Hindisi Tüyü | 1× Kum Çeliği Yay | 185 | 325s | 73 | 20s 64c |
| recipe_crossbow_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi, 2× Yaban Hindisi Tüyü | 1× Kum Çeliği Arbalet | 185 | 325s | 73 | 20s 64c |
| recipe_lute_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Lavta | 185 | 325s | 73 | 20s 64c |
| recipe_focus_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Doğa Odağı | 185 | 325s | 73 | 20s 64c |
| recipe_staff_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Asa | 185 | 325s | 73 | 20s 64c |
| recipe_tonfa_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Tonfa | 185 | 325s | 73 | 20s 64c |
| recipe_totem_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Totem | 185 | 325s | 73 | 20s 64c |
| recipe_wand_sandsteel | 19 | 4× Kayın Kütüğü, 1× Bakır Külçesi | 1× Kum Çeliği Değnek | 185 | 325s | 73 | 20s 64c |
| recipe_bow_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Rün Bağlı Yay | 210 | 345s | 82 | 24s 34c |
| recipe_crossbow_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Rün Bağlı Arbalet | 210 | 345s | 82 | 24s 34c |
| recipe_lute_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Lavta | 210 | 345s | 82 | 24s 34c |
| recipe_focus_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Doğa Odağı | 210 | 345s | 82 | 24s 34c |
| recipe_staff_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Asa | 210 | 345s | 82 | 24s 34c |
| recipe_tonfa_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Tonfa | 210 | 345s | 82 | 24s 34c |
| recipe_totem_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Totem | 210 | 345s | 82 | 24s 34c |
| recipe_wand_runebound | 21 | 4× Dişbudak Kütüğü, 1× Bronz Külçesi | 1× Rün Bağlı Değnek | 210 | 345s | 82 | 24s 34c |
| recipe_bow_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Anka Yay | 235 | 375s | 97 | 28s 17c |
| recipe_crossbow_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Anka Arbalet | 235 | 375s | 97 | 28s 17c |
| recipe_lute_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Lavta | 235 | 375s | 97 | 28s 17c |
| recipe_focus_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Doğa Odağı | 235 | 375s | 97 | 28s 17c |
| recipe_staff_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Asa | 235 | 375s | 97 | 28s 17c |
| recipe_tonfa_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Tonfa | 235 | 375s | 97 | 28s 17c |
| recipe_totem_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Totem | 235 | 375s | 97 | 28s 17c |
| recipe_wand_phoenix | 24 | 4× Akçaağaç Kütüğü, 1× Bronz Külçesi | 1× Anka Değnek | 235 | 375s | 97 | 28s 17c |
| recipe_bow_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Dehşet Dövme Yay | 260 | 395s | 109 | 32s 12c |
| recipe_crossbow_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi, 2× Yaban Hindisi Tüyü | 1× Dehşet Dövme Arbalet | 260 | 395s | 109 | 32s 12c |
| recipe_lute_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Lavta | 260 | 395s | 109 | 32s 12c |
| recipe_focus_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Doğa Odağı | 260 | 395s | 109 | 32s 12c |
| recipe_staff_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Asa | 260 | 395s | 109 | 32s 12c |
| recipe_tonfa_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Tonfa | 260 | 395s | 109 | 32s 12c |
| recipe_totem_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Totem | 260 | 395s | 109 | 32s 12c |
| recipe_wand_dreadforged | 26 | 4× Ceviz Kütüğü, 1× Bronz Külçesi | 1× Dehşet Dövme Değnek | 260 | 395s | 109 | 32s 12c |
| recipe_bow_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi, 2× Yaban Hindisi Tüyü | 1× Hükümdar Yay | 285 | 420s | 130 | 36s 20c |
| recipe_crossbow_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi, 2× Yaban Hindisi Tüyü | 1× Hükümdar Arbalet | 285 | 420s | 130 | 36s 20c |
| recipe_lute_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Lavta | 285 | 420s | 130 | 36s 20c |
| recipe_focus_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Doğa Odağı | 285 | 420s | 130 | 36s 20c |
| recipe_staff_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Asa | 285 | 420s | 130 | 36s 20c |
| recipe_tonfa_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Tonfa | 285 | 420s | 130 | 36s 20c |
| recipe_totem_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Totem | 285 | 420s | 130 | 36s 20c |
| recipe_wand_sovereign | 29 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Hükümdar Değnek | 285 | 420s | 130 | 36s 20c |
| recipe_bow_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi, 2× Yaban Hindisi Tüyü | 1× Mercan Yay | 310 | 440s | 146 | 40s 38c |
| recipe_crossbow_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi, 2× Yaban Hindisi Tüyü | 1× Mercan Arbalet | 310 | 440s | 146 | 40s 38c |
| recipe_lute_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Lavta | 310 | 440s | 146 | 40s 38c |
| recipe_focus_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Doğa Odağı | 310 | 440s | 146 | 40s 38c |
| recipe_staff_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Asa | 310 | 440s | 146 | 40s 38c |
| recipe_tonfa_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Tonfa | 310 | 440s | 146 | 40s 38c |
| recipe_totem_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Totem | 310 | 440s | 146 | 40s 38c |
| recipe_wand_coral | 31 | 4× Maun Kütüğü, 1× Nikel Külçesi | 1× Mercan Değnek | 310 | 440s | 146 | 40s 38c |
| recipe_bow_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi, 2× Yaban Hindisi Tüyü | 1× Ejderha Kemiği Yay | 335 | 470s | 174 | 44s 66c |
| recipe_crossbow_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi, 2× Yaban Hindisi Tüyü | 1× Ejderha Kemiği Arbalet | 335 | 470s | 174 | 44s 66c |
| recipe_lute_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Lavta | 335 | 470s | 174 | 44s 66c |
| recipe_focus_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Doğa Odağı | 335 | 470s | 174 | 44s 66c |
| recipe_staff_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Asa | 335 | 470s | 174 | 44s 66c |
| recipe_tonfa_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Tonfa | 335 | 470s | 174 | 44s 66c |
| recipe_totem_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Totem | 335 | 470s | 174 | 44s 66c |
| recipe_wand_drakebone | 34 | 4× Abanoz Kütüğü, 1× Krom Külçesi | 1× Ejderha Kemiği Değnek | 335 | 470s | 174 | 44s 66c |
| recipe_bow_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi, 2× Yaban Hindisi Tüyü | 1× Sis Örgü Yay | 360 | 490s | 196 | 49s 4c |
| recipe_crossbow_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi, 2× Yaban Hindisi Tüyü | 1× Sis Örgü Arbalet | 360 | 490s | 196 | 49s 4c |
| recipe_lute_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Lavta | 360 | 490s | 196 | 49s 4c |
| recipe_focus_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Doğa Odağı | 360 | 490s | 196 | 49s 4c |
| recipe_staff_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Asa | 360 | 490s | 196 | 49s 4c |
| recipe_tonfa_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Tonfa | 360 | 490s | 196 | 49s 4c |
| recipe_totem_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Totem | 360 | 490s | 196 | 49s 4c |
| recipe_wand_mistwoven | 36 | 4× Gümüş Söğüt Kütüğü, 1× Kobalt Külçesi | 1× Sis Örgü Değnek | 360 | 490s | 196 | 49s 4c |
| recipe_bow_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi, 2× Yaban Hindisi Tüyü | 1× Prizmatik Yay | 385 | 520s | 233 | 53s 51c |
| recipe_crossbow_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi, 2× Yaban Hindisi Tüyü | 1× Prizmatik Arbalet | 385 | 520s | 233 | 53s 51c |
| recipe_lute_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Lavta | 385 | 520s | 233 | 53s 51c |
| recipe_focus_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Doğa Odağı | 385 | 520s | 233 | 53s 51c |
| recipe_staff_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Asa | 385 | 520s | 233 | 53s 51c |
| recipe_tonfa_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Tonfa | 385 | 520s | 233 | 53s 51c |
| recipe_totem_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Totem | 385 | 520s | 233 | 53s 51c |
| recipe_wand_prismatic | 39 | 4× Altın Meşe Kütüğü, 1× Gümüş Külçesi | 1× Prizmatik Değnek | 385 | 520s | 233 | 53s 51c |
| recipe_bow_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi, 2× Yaban Hindisi Tüyü | 1× Magma Dövme Yay | 410 | 540s | 262 | 58s 8c |
| recipe_crossbow_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi, 2× Yaban Hindisi Tüyü | 1× Magma Dövme Arbalet | 410 | 540s | 262 | 58s 8c |
| recipe_lute_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Lavta | 410 | 540s | 262 | 58s 8c |
| recipe_focus_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Doğa Odağı | 410 | 540s | 262 | 58s 8c |
| recipe_staff_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Asa | 410 | 540s | 262 | 58s 8c |
| recipe_tonfa_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Tonfa | 410 | 540s | 262 | 58s 8c |
| recipe_totem_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Totem | 410 | 540s | 262 | 58s 8c |
| recipe_wand_magmaforged | 41 | 4× Kara Odun Kütüğü, 1× Altın Külçesi | 1× Magma Dövme Değnek | 410 | 540s | 262 | 58s 8c |
| recipe_bow_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Bora Yay | 435 | 565s | 312 | 62s 72c |
| recipe_crossbow_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Bora Arbalet | 435 | 565s | 312 | 62s 72c |
| recipe_lute_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Lavta | 435 | 565s | 312 | 62s 72c |
| recipe_focus_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Doğa Odağı | 435 | 565s | 312 | 62s 72c |
| recipe_staff_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Asa | 435 | 565s | 312 | 62s 72c |
| recipe_tonfa_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Tonfa | 435 | 565s | 312 | 62s 72c |
| recipe_totem_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Totem | 435 | 565s | 312 | 62s 72c |
| recipe_wand_tempest | 44 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Bora Değnek | 435 | 565s | 312 | 62s 72c |
| recipe_bow_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Hiçlik Yay | 460 | 585s | 350 | 67s 45c |
| recipe_crossbow_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Hiçlik Arbalet | 460 | 585s | 350 | 67s 45c |
| recipe_lute_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Lavta | 460 | 585s | 350 | 67s 45c |
| recipe_focus_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Doğa Odağı | 460 | 585s | 350 | 67s 45c |
| recipe_staff_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Asa | 460 | 585s | 350 | 67s 45c |
| recipe_tonfa_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Tonfa | 460 | 585s | 350 | 67s 45c |
| recipe_totem_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Totem | 460 | 585s | 350 | 67s 45c |
| recipe_wand_voidtouched | 46 | 4× Kristal Huş Kütüğü, 1× Elektrum Külçesi | 1× Hiçlik Değnek | 460 | 585s | 350 | 67s 45c |
| recipe_bow_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi, 2× Fırtına Kartalı Tüyü | 1× Melek Yay | 485 | 615s | 417 | 72s 25c |
| recipe_crossbow_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi, 2× Fırtına Kartalı Tüyü | 1× Melek Arbalet | 485 | 615s | 417 | 72s 25c |
| recipe_lute_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Lavta | 485 | 615s | 417 | 72s 25c |
| recipe_focus_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Doğa Odağı | 485 | 615s | 417 | 72s 25c |
| recipe_staff_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Asa | 485 | 615s | 417 | 72s 25c |
| recipe_tonfa_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Tonfa | 485 | 615s | 417 | 72s 25c |
| recipe_totem_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Totem | 485 | 615s | 417 | 72s 25c |
| recipe_wand_seraphic | 49 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Melek Değnek | 485 | 615s | 417 | 72s 25c |
| recipe_bow_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi, 2× Fırtına Kartalı Tüyü | 1× Tutulma Yay | 510 | 635s | 469 | 77s 13c |
| recipe_crossbow_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi, 2× Fırtına Kartalı Tüyü | 1× Tutulma Arbalet | 510 | 635s | 469 | 77s 13c |
| recipe_lute_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Lavta | 510 | 635s | 469 | 77s 13c |
| recipe_focus_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Doğa Odağı | 510 | 635s | 469 | 77s 13c |
| recipe_staff_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Asa | 510 | 635s | 469 | 77s 13c |
| recipe_tonfa_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Tonfa | 510 | 635s | 469 | 77s 13c |
| recipe_totem_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Totem | 510 | 635s | 469 | 77s 13c |
| recipe_wand_eclipse | 51 | 4× Korağacı Kütüğü, 1× Tungsten Külçesi | 1× Tutulma Değnek | 510 | 635s | 469 | 77s 13c |
| recipe_bow_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Ejder Lordu Yay | 535 | 660s | 558 | 82s 8c |
| recipe_crossbow_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Ejder Lordu Arbalet | 535 | 660s | 558 | 82s 8c |
| recipe_lute_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Lavta | 535 | 660s | 558 | 82s 8c |
| recipe_focus_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Doğa Odağı | 535 | 660s | 558 | 82s 8c |
| recipe_staff_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Asa | 535 | 660s | 558 | 82s 8c |
| recipe_tonfa_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Tonfa | 535 | 660s | 558 | 82s 8c |
| recipe_totem_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Totem | 535 | 660s | 558 | 82s 8c |
| recipe_wand_wyrmlord | 54 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Ejder Lordu Değnek | 535 | 660s | 558 | 82s 8c |
| recipe_bow_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Göksel Yay | 560 | 680s | 627 | 87s 10c |
| recipe_crossbow_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Göksel Arbalet | 560 | 680s | 627 | 87s 10c |
| recipe_lute_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Lavta | 560 | 680s | 627 | 87s 10c |
| recipe_focus_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Doğa Odağı | 560 | 680s | 627 | 87s 10c |
| recipe_staff_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Asa | 560 | 680s | 627 | 87s 10c |
| recipe_tonfa_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Tonfa | 560 | 680s | 627 | 87s 10c |
| recipe_totem_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Totem | 560 | 680s | 627 | 87s 10c |
| recipe_wand_celestine | 56 | 4× Ayaz Çamı Kütüğü, 1× Titanyum Külçesi | 1× Göksel Değnek | 560 | 680s | 627 | 87s 10c |
| recipe_bow_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Ebedi Yay | 585 | 710s | 747 | 92s 19c |
| recipe_crossbow_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Ebedi Arbalet | 585 | 710s | 747 | 92s 19c |
| recipe_lute_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Lavta | 585 | 710s | 747 | 92s 19c |
| recipe_focus_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Doğa Odağı | 585 | 710s | 747 | 92s 19c |
| recipe_staff_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Asa | 585 | 710s | 747 | 92s 19c |
| recipe_tonfa_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Tonfa | 585 | 710s | 747 | 92s 19c |
| recipe_totem_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Totem | 585 | 710s | 747 | 92s 19c |
| recipe_wand_eternal | 59 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Ebedi Değnek | 585 | 710s | 747 | 92s 19c |
| recipe_bow_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yıldızsız Yay | 610 | 730s | 839 | 97s 34c |
| recipe_crossbow_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yıldızsız Arbalet | 610 | 730s | 839 | 97s 34c |
| recipe_lute_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Lavta | 610 | 730s | 839 | 97s 34c |
| recipe_focus_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Doğa Odağı | 610 | 730s | 839 | 97s 34c |
| recipe_staff_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Asa | 610 | 730s | 839 | 97s 34c |
| recipe_tonfa_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Tonfa | 610 | 730s | 839 | 97s 34c |
| recipe_totem_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Totem | 610 | 730s | 839 | 97s 34c |
| recipe_wand_starless | 61 | 4× Kanağacı Kütüğü, 1× Paladyum Külçesi | 1× Yıldızsız Değnek | 610 | 730s | 839 | 97s 34c |
| recipe_bow_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Gece Çeliği Yay | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_crossbow_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Gece Çeliği Arbalet | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_lute_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Lavta | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_focus_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Doğa Odağı | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_staff_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Asa | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_tonfa_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Tonfa | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_totem_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Totem | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_wand_nightsteel | 64 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Gece Çeliği Değnek | 635 | 760s | 1000 | 1g 2s 56c |
| recipe_bow_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Zaman Dövme Yay | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_crossbow_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Zaman Dövme Arbalet | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_lute_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Lavta | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_focus_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Doğa Odağı | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_staff_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Asa | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_tonfa_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Tonfa | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_totem_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Totem | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_wand_chronoforged | 66 | 4× Ruhağacı Kütüğü, 1× Rodyum Külçesi | 1× Zaman Dövme Değnek | 660 | 780s | 1123 | 1g 7s 84c |
| recipe_bow_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yükselmiş Yay | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_crossbow_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yükselmiş Arbalet | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_lute_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Lavta | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_focus_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Doğa Odağı | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_staff_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Asa | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_tonfa_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Tonfa | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_totem_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Totem | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_wand_ascendant | 69 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Yükselmiş Değnek | 685 | 805s | 1338 | 1g 13s 18c |
| recipe_bow_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi, 2× Fırtına Kartalı Tüyü | 1× Fırtına Doğumlu Yay | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_crossbow_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi, 2× Fırtına Kartalı Tüyü | 1× Fırtına Doğumlu Arbalet | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_lute_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Lavta | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_focus_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Doğa Odağı | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_staff_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Asa | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_tonfa_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Tonfa | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_totem_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Totem | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_wand_stormborn | 71 | 4× Kadim Demirağaç Kütüğü, 1× Platin Külçesi | 1× Fırtına Doğumlu Değnek | 710 | 825s | 1503 | 1g 18s 58c |
| recipe_bow_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yarık Dövme Yay | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_crossbow_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Yarık Dövme Arbalet | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_lute_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Lavta | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_focus_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Doğa Odağı | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_staff_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Asa | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_tonfa_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Tonfa | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_totem_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Totem | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_wand_riftforged | 74 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Yarık Dövme Değnek | 735 | 855s | 1790 | 1g 24s 4c |
| recipe_bow_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Kaos Yay | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_crossbow_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Kaos Arbalet | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_lute_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Lavta | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_focus_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Doğa Odağı | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_staff_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Asa | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_tonfa_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Tonfa | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_totem_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Totem | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_wand_chaos | 76 | 4× Yıldızkabuğu Kütüğü, 1× Osmiyum Külçesi | 1× Kaos Değnek | 760 | 875s | 2011 | 1g 29s 55c |
| recipe_bow_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Gök Dövme Yay | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_crossbow_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Gök Dövme Arbalet | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_lute_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Lavta | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_focus_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Doğa Odağı | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_staff_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Asa | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_tonfa_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Tonfa | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_totem_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Totem | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_wand_skyforged | 79 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Gök Dövme Değnek | 785 | 900s | 2395 | 1g 35s 11c |
| recipe_bow_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Başmelek Yay | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_crossbow_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi, 2× Fırtına Kartalı Tüyü | 1× Başmelek Arbalet | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_lute_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Lavta | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_focus_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Doğa Odağı | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_staff_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Asa | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_tonfa_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Tonfa | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_totem_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Totem | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_wand_archangel | 81 | 4× Ayağacı Kütüğü, 1× İridyum Külçesi | 1× Başmelek Değnek | 810 | 920s | 2692 | 1g 40s 73c |
| recipe_bow_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi, 2× Yıldız Grifonu Tüyü | 1× Şafak Yay | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_crossbow_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi, 2× Yıldız Grifonu Tüyü | 1× Şafak Arbalet | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_lute_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Lavta | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_focus_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Doğa Odağı | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_staff_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Asa | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_tonfa_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Tonfa | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_totem_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Totem | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_wand_aurora | 84 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Şafak Değnek | 835 | 950s | 3206 | 1g 46s 41c |
| recipe_bow_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi, 2× Yıldız Grifonu Tüyü | 1× Kül Yay | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_crossbow_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi, 2× Yıldız Grifonu Tüyü | 1× Kül Arbalet | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_lute_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Lavta | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_focus_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Doğa Odağı | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_staff_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Asa | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_tonfa_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Tonfa | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_totem_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Totem | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_wand_ashen | 86 | 4× Ejderyüreği Kütüğü, 1× Yıldız Çeliği Külçesi | 1× Kül Değnek | 860 | 970s | 3602 | 1g 52s 13c |
| recipe_bow_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi, 2× Yıldız Grifonu Tüyü | 1× Leviathan Yay | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_crossbow_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi, 2× Yıldız Grifonu Tüyü | 1× Leviathan Arbalet | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_lute_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Lavta | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_focus_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Doğa Odağı | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_staff_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Asa | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_tonfa_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Tonfa | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_totem_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Totem | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_wand_leviathan | 89 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Leviathan Değnek | 885 | 1000s | 4290 | 1g 57s 90c |
| recipe_bow_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi, 2× Yıldız Grifonu Tüyü | 1× Kıyamet Dövme Yay | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_crossbow_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi, 2× Yıldız Grifonu Tüyü | 1× Kıyamet Dövme Arbalet | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_lute_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Lavta | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_focus_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Doğa Odağı | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_staff_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Asa | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_tonfa_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Tonfa | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_totem_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Totem | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_wand_doomforged | 91 | 4× Boşlukağacı Kütüğü, 1× Mitril Külçesi | 1× Kıyamet Dövme Değnek | 910 | 1020s | 4820 | 1g 63s 73c |
| recipe_bow_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi, 2× Kadim Anka Tüyü | 1× Hiçlik Dövme Yay | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_crossbow_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi, 2× Kadim Anka Tüyü | 1× Hiçlik Dövme Arbalet | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_lute_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Lavta | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_focus_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Doğa Odağı | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_staff_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Asa | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_tonfa_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Tonfa | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_totem_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Totem | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_wand_voidforged | 94 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Hiçlik Dövme Değnek | 935 | 1045s | 5741 | 1g 69s 60c |
| recipe_bow_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi, 2× Kadim Anka Tüyü | 1× Tanrı Katili Yay | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_crossbow_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi, 2× Kadim Anka Tüyü | 1× Tanrı Katili Arbalet | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_lute_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Lavta | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_focus_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Doğa Odağı | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_staff_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Asa | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_tonfa_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Tonfa | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_totem_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Totem | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_wand_godslayer | 96 | 4× Göksel Dişbudak Kütüğü, 1× Orikalkum Külçesi | 1× Tanrı Katili Değnek | 960 | 1065s | 6450 | 1g 75s 52c |
| recipe_bow_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi, 2× Kadim Anka Tüyü | 1× Kadim Çağ Yay | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_crossbow_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi, 2× Kadim Anka Tüyü | 1× Kadim Çağ Arbalet | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_lute_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Lavta | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_focus_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Doğa Odağı | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_staff_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Asa | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_tonfa_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Tonfa | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_totem_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Totem | 985 | 1095s | 7682 | 1g 81s 48c |
| recipe_wand_primordial | 99 | 4× Dünya Ağacı Kütüğü, 1× Adamantit Külçesi | 1× Kadim Çağ Değnek | 985 | 1095s | 7682 | 1g 81s 48c |

## 4. Ganimet (drop) kuralları

- **Ganimet seviyesi** = min(canavar Sv, karakter Sv) (`activity.DropLevel`).
- **Ekipman**: ganimet seviyesine eşit ya da altındaki **en yüksek tek kademe**, karakterin sınıfına uygun (`item.GearPool`; gereksinimli eşya düşmez).
  - Av: her zafer %4 → sıradan %60 · seçkin %25 · nadir %11 · destansı %4.
  - Zindan: dalga öldürmesi başına %4 + boss’ta 1 garanti → boss tablosu seçkin %50 · nadir %35 · destansı %13 · efsanevi %2.
  - Baskın: dalga %4 + boss’ta 2 garanti, kişisel ganimet.
- **Deri/post** (`item.CombatHides`): avcılık seviyesi = ⌈canavar Sv / 10⌉ (tavan 100, DD-274). Sıradan buluntu açılan son 4 hayvandan, nadir/efsanevi en iyisinden.
- **Toplama**: seçilen kaynağın kendisi, kademe rengiyle (`TierRarity`): T0–1 sıradan, T2–3 seçkin, T4 nadir, T5 destansı, T6 efsanevi, T7 mitik.
- **Altın / XP**: canavar tablosu; XP seviye farkı bantları (10 aşağı tam, 25 aşağı %50, ötesi %10).

## 5. Canavar listesi

Bölge `mobs.yaml` `region_code`; canavar kendi bölgesinden ya da ona yolu olan bölgeden avlanır (DD-273). Zindan ve baskınlar bölgeye bağlı değil; baskın boss’unun sütununda baskın adı var.

| Sv | Kod | Ad | Arketip | Bölge | Can | Sal | Sav | XP | Altın | Ekipman kademesi | Deri (en iyi) | Rol |
|---:|---|---|---|---|---:|---:|---:|---:|---|---:|---|---|
| 1 | young_boar | Genç Yaban Domuzu | — | Yeşilorman | 172 | 29 | 2 | 10 | 7c–18c | 1 | Orman Faresi Postu | görev |
| 3 | grey_wolf | Boz Kurt | — | Yeşilorman | 212 | 35 | 4 | 64 | 24c–60c | 1 | Orman Faresi Postu | dalga ×1, görev |
| 5 | bandit_scout | Haydut Gözcü | — | Yeşilorman | 300 | 41 | 6 | 131 | 48c–1s 20c | 5 | Orman Faresi Postu | dalga ×1, görev |
| 6 | cave_bat | Mağara Yarasası | — | Yeşilorman | 266 | 42 | 6 | 193 | 62c–1s 55c | 5 | Orman Faresi Postu | dalga ×1 |
| 8 | goblin_raider | Goblin Yağmacı | — | Fısıltı Vadisi | 396 | 50 | 9 | 279 | 94c–2s 35c | 5 | Orman Faresi Postu | boss Kül Mahzeni, dalga ×1, görev |
| 10 | marsh_lurker | Bataklık Pusucusu | — | Fısıltı Vadisi | 460 | 56 | 11 | 398 | 1s 29c–3s 22c | 10 | Orman Faresi Postu | dalga ×1 |
| 12 | forest_troll | Orman Trolü | — | Fısıltı Vadisi | 760 | 53 | 21 | 640 | 1s 69c–4s 22c | 10 | Orman Faresi Postu | dalga ×1, görev |
| 14 | dire_wolf | Dev Kurt | — | Fısıltı Vadisi | 529 | 68 | 14 | 750 | 2s 13c–5s 32c | 10 | Orman Faresi Postu | boss Orman Kovuğu, dalga ×1 |
| 15 | cultist_acolyte | Tarikat Müridi | — | Fısıltı Vadisi | 465 | 85 | 10 | 800 | 2s 35c–5s 88c | 15 | Orman Faresi Postu | dalga ×1 |
| 18 | rock_ogre | Kaya Devi | — | Közdoruk | 1038 | 68 | 30 | 1224 | 3s 8c–7s 70c | 15 | Orman Faresi Postu | dalga ×1 |
| 20 | shadow_stalker | Gölge Avcısı | — | Közdoruk | 624 | 82 | 17 | 1328 | 3s 61c–9s 2c | 20 | Orman Faresi Postu | boss Gölgebataklık Harabeleri, dalga ×1 |
| 22 | frost_giant | Buz Devi | — | Közdoruk | 1224 | 78 | 37 | 1687 | 4s 16c–10s 40c | 20 | Orman Faresi Postu | dalga ×1 |
| 24 | elder_wyrm | Kadim Ejder | — | Közdoruk | 1317 | 83 | 40 | 1939 | 4s 73c–11s 82c | 20 | Orman Faresi Postu | dalga ×1 |
| 25 | ashen_lich | Kül Ölüsü | — | Közdoruk | 705 | 121 | 16 | 1811 | 5s 3c–12s 58c | 25 | Orman Faresi Postu | boss Obsidyen Kule |
| 30 | mist_leech | Sis Sülüğü | evasive | Sisli Bataklıklar | 855 | 105 | 25 | 2309 | 6s 61c–16s 52c | 25 | Orman Faresi Postu | dalga ×1 |
| 35 | mud_golem | Çamur Golemi | tank | Sisli Bataklıklar | 1794 | 102 | 55 | 2955 | 8s 33c–20s 82c | 35 | Küçük Tavşan Postu | dalga ×1 |
| 40 | bog_witch | Bataklık Cadısı | glass | Sisli Bataklıklar | 1035 | 168 | 25 | 3658 | 10s 18c–25s 45c | 35 | Küçük Tavşan Postu | dalga ×1 |
| 42 | marsh_queen | Bataklık Kraliçesi | chief | Sisli Bataklıklar | 1440 | 125 | 38 | 3955 | 10s 95c–27s 38c | 35 | Küçük Tavşan Postu | boss Batık Tapınak |
| 45 | venom_lizard | Zehirli Kertenkele | crit | Sisli Bataklıklar | 1339 | 145 | 37 | 4417 | 12s 15c–30s 38c | 35 | Küçük Tavşan Postu |  |
| 50 | drowned_knight | Batık Şövalye | elite | Sisli Bataklıklar | 1960 | 153 | 49 | 5228 | 14s 23c–35s 58c | 50 | Küçük Tavşan Postu |  |
| 55 | mountain_troll | Dağ Trolü | tank | Demir Dağlar | 2669 | 144 | 86 | 6089 | 16s 41c–41s 2c | 50 | Küçük Tavşan Postu | dalga ×1 |
| 60 | dwarf_ghost | Cüce Hayaleti | evasive | Demir Dağlar | 1568 | 177 | 49 | 6999 | 18s 70c–46s 75c | 50 | Küçük Tavşan Postu | dalga ×1 |
| 65 | rock_eagle | Kaya Kartalı | glass | Demir Dağlar | 1598 | 247 | 40 | 7955 | 21s 9c–52s 72c | 65 | Tarla Faresi Postu | dalga ×1 |
| 67 | iron_king | Demir Kral | chief | Demir Dağlar | 2190 | 182 | 60 | 8351 | 22s 7c–55s 18c | 65 | Tarla Faresi Postu | boss Unutulmuş Maden |
| 70 | iron_spider | Demir Örümcek | crit | Demir Dağlar | 1995 | 206 | 57 | 8957 | 23s 57c–58s 92c | 65 | Tarla Faresi Postu |  |
| 75 | mine_warden | Maden Bekçisi | elite | Demir Dağlar | 2835 | 213 | 73 | 10002 | 26s 14c–65s 35c | 65 | Tarla Faresi Postu |  |
| 80 | ghost_guard | Hayalet Muhafız | normal | Kuzey Kalesi Harabeleri | 2150 | 196 | 65 | 11090 | 28s 79c–71s 98c | 80 | Tarla Faresi Postu | dalga ×1 |
| 85 | cursed_archer | Lanetli Okçu | glass | Kuzey Kalesi Harabeleri | 2048 | 310 | 52 | 12220 | 31s 53c–78s 82c | 80 | Tarla Faresi Postu | dalga ×1 |
| 90 | bone_rider | Kemik Atlı | crit | Kuzey Kalesi Harabeleri | 2520 | 256 | 73 | 13390 | 34s 36c–85s 90c | 80 | Tarla Faresi Postu | dalga ×2 |
| 92 | fallen_lord_varn | Düşmüş Lord Varn | chief | Kuzey Kalesi Harabeleri | 2940 | 239 | 82 | 13869 | 35s 51c–88s 78c | 80 | Sincap Postu | boss Kuzey Kalesi |
| 95 | keep_gargoyle | Kale Gargoyli | tank | Kuzey Kalesi Harabeleri | 4419 | 228 | 146 | 14600 | 37s 26c–93s 15c | 95 | Sincap Postu | dalga ×1 |
| 100 | fallen_commander | Düşmüş Komutan | elite | Kuzey Kalesi Harabeleri | 3710 | 274 | 97 | 15849 | 40s 24c–1g 60c | 95 | Sincap Postu | dalga ×1 |
| 102 | varn_undying | Ölümsüz Varn | titan | Kuzey Kalesi Kuşatması | 81000 | 775 | 116 | 16359 | 41s 45c–1g 3s 62c | 95 | Sincap Postu | baskın boss’u Kuzey Kalesi Kuşatması |
| 110 | rebel_rider | Asi Süvari | normal | Taç Ovaları | 2900 | 259 | 89 | 18460 | 46s 42c–1g 16s 5c | 110 | Sincap Postu | dalga ×1 |
| 120 | plains_lion | Ova Aslanı | crit | Taç Ovaları | 3308 | 330 | 97 | 21217 | 52s 90c–1g 32s 25c | 110 | Yaban Tavşanı Postu | dalga ×1 |
| 130 | raider_mage | Yağmacı Büyücü | evasive | Taç Ovaları | 3230 | 346 | 105 | 24116 | 59s 64c–1g 49s 10c | 110 | Yaban Tavşanı Postu | dalga ×1 |
| 132 | rebel_king_orrin | Asi Kral Orrin | elite | Taç Ovaları | 4830 | 351 | 128 | 24713 | 61s 3c–1g 52s 58c | 110 | Porsuk Postu | boss Asiler Kampı |
| 140 | war_bear | Savaş Ayısı | tank | Taç Ovaları | 6388 | 322 | 215 | 27152 | 66s 66c–1g 66s 65c | 135 | Porsuk Postu |  |
| 150 | black_lancer | Kara Mızrakçı | elite | Taç Ovaları | 5460 | 394 | 145 | 30321 | 73s 93c–1g 84s 82c | 135 | Porsuk Postu |  |
| 160 | ash_wolf | Kül Kurdu | evasive | Yanık Orman | 3942 | 419 | 129 | 33620 | 81s 44c–2g 3s 60c | 160 | Yaban Hindisi Tüyü | dalga ×1 |
| 170 | flame_spirit | Alev Ruhu | glass | Yanık Orman | 3960 | 578 | 103 | 37044 | 89s 19c–2g 22s 98c | 160 | Yaban Hindisi Tüyü | dalga ×1 |
| 180 | cinder_treant | Kömür Treanti | tank | Yanık Orman | 8138 | 406 | 276 | 40592 | 97s 18c–2g 42s 95c | 160 | Orman Geyiği Derisi | dalga ×1 |
| 182 | burnt_mother_tree | Yanık Ana Ağaç | elite | Yanık Orman | 6580 | 472 | 176 | 41316 | 98s 80c–2g 47s | 160 | Orman Geyiği Derisi | boss Ateş Kökü |
| 190 | crimson_harpy | Kızıl Harpya | crit | Yanık Orman | 5145 | 504 | 153 | 44259 | 1g 5s 39c–2g 63s 48c | 185 | Orman Geyiği Derisi |  |
| 200 | scorched_shaman | Yanık Şaman | elite | Yanık Orman | 7210 | 515 | 193 | 48045 | 1g 13s 82c–2g 84s 55c | 185 | Kızıl Tilki Postu |  |
| 210 | sand_scorpion | Kum Akrebi | crit | Tuz Çölü | 5670 | 553 | 169 | 51946 | 1g 22s 46c–3g 6s 15c | 210 | Kızıl Tilki Postu | dalga ×1 |
| 220 | desert_viper | Çöl Yılanı | evasive | Tuz Çölü | 5368 | 564 | 177 | 55960 | 1g 31s 31c–3g 28s 28c | 210 | Yaban Domuzu Postu | dalga ×1 |
| 230 | salt_elemental | Tuz Elementali | tank | Tuz Çölü | 10325 | 511 | 352 | 60085 | 1g 40s 36c–3g 50s 90c | 210 | Yaban Domuzu Postu | dalga ×1 |
| 232 | pharaoh_seth_ka | Firavun Seth-Ka | elite | Tuz Çölü | 8330 | 592 | 224 | 60923 | 1g 42s 20c–3g 55s 50c | 210 | Dağ Keçisi Postu | boss Gömülü Saray |
| 240 | mummy_priest | Mumya Rahip | glass | Tuz Çölü | 5535 | 798 | 145 | 64319 | 1g 49s 61c–3g 74s 2c | 235 | Dağ Keçisi Postu |  |
| 250 | sand_worm | Kum Solucanı | elite | Tuz Çölü | 8960 | 636 | 241 | 68660 | 1g 59s 6c–3g 97s 65c | 235 | Dağ Keçisi Postu |  |
| 260 | royal_skeleton | Kraliyet İskeleti | normal | Taç Mezarlığı | 6650 | 574 | 209 | 73107 | 1g 68s 70c–4g 21s 75c | 260 | Yaban Kurdu Postu | dalga ×1 |
| 270 | tomb_keeper | Mezar Bekçisi | tank | Taç Mezarlığı | 12075 | 595 | 412 | 77657 | 1g 78s 53c–4g 46s 32c | 260 | Yaban Kurdu Postu | dalga ×2 |
| 280 | soul_drinker | Ruh Emici | glass | Taç Mezarlığı | 6435 | 924 | 169 | 82310 | 1g 88s 54c–4g 71s 35c | 260 | Orman Ayısı Postu | dalga ×1 |
| 282 | crypt_lord_malachar | Mezar Lordu Malachar | elite | Taç Mezarlığı | 10080 | 713 | 272 | 83253 | 1g 90s 56c–4g 76s 40c | 260 | Orman Ayısı Postu | boss Kraliyet Kriptası |
| 290 | black_knight | Kara Şövalye | crit | Taç Mezarlığı | 7770 | 752 | 233 | 87064 | 1g 98s 73c–4g 96s 82c | 285 | Orman Ayısı Postu | dalga ×1 |
| 300 | crownless_king | Taçsız Kral | elite | Taç Mezarlığı | 10710 | 757 | 289 | 91917 | 2g 9s 9c–5g 22s 72c | 285 | Yaban Mandası Derisi | dalga ×1 |
| 302 | crownless_tyrant | Taçsız Tiran | titan | Taçsız Kral'ın Sarayı | 231000 | 2119 | 340 | 92899 | 2g 11s 19c–5g 27s 98c | 285 | Yaban Mandası Derisi | baskın boss’u Taçsız Kral'ın Sarayı |
| 310 | sea_raider | Deniz Haydudu | normal | Fırtına Kıyısı | 7900 | 679 | 249 | 96868 | 2g 19s 63c–5g 49s 8c | 310 | Yaban Mandası Derisi | dalga ×1 |
| 320 | storm_bird | Fırtına Kuşu | evasive | Fırtına Kıyısı | 7742 | 805 | 257 | 101915 | 2g 30s 35c–5g 75s 88c | 310 | Yaban Mandası Derisi | dalga ×1 |
| 330 | coral_golem | Mercan Golemi | tank | Fırtına Kıyısı | 14700 | 721 | 504 | 107059 | 2g 41s 23c–6g 3s 8c | 310 | Dev Geyik Postu | dalga ×1 |
| 332 | admiral_drask | Amiral Drask | overlord | Fırtına Kıyısı | 15632 | 943 | 333 | 108099 | 2g 43s 42c–6g 8s 55c | 310 | Dev Geyik Postu | boss Batık Filo |
| 340 | siren | Siren | glass | Fırtına Kıyısı | 7785 | 1113 | 205 | 112297 | 2g 52s 28c–6g 30s 70c | 335 | Dev Geyik Postu |  |
| 350 | kraken_spawn | Kraken Yavrusu | elite | Fırtına Kıyısı | 12460 | 877 | 337 | 117628 | 2g 63s 49c–6g 58s 72c | 335 | Dev Geyik Postu |  |
| 360 | young_dragon | Genç Ejder | crit | Ejder Yamaçları | 9608 | 925 | 289 | 123051 | 2g 74s 86c–6g 87s 15c | 360 | Dev Kurt Postu | dalga ×1 |
| 370 | dragon_cultist | Ejder Kültisti | glass | Ejder Yamaçları | 8460 | 1208 | 223 | 128565 | 2g 86s 39c–7g 15s 98c | 360 | Dev Kurt Postu | dalga ×1 |
| 380 | lava_lizard | Lav Kertenkelesi | tank | Ejder Yamaçları | 16888 | 826 | 580 | 134169 | 2g 98s 8c–7g 45s 20c | 360 | Dev Kurt Postu | dalga ×1 |
| 382 | mother_dragon_vyrna | Ana Ejder Vyrna | overlord | Ejder Yamaçları | 17945 | 1079 | 383 | 135301 | 3g 44c–7g 51s 10c | 360 | Dağ Aslanı Postu | boss Kuluçka Mağarası |
| 390 | winged_wyvern | Kanatlı Wyvern | evasive | Ejder Yamaçları | 9405 | 974 | 313 | 139863 | 3g 9s 92c–7g 74s 80c | 385 | Dağ Aslanı Postu |  |
| 400 | dragon_guard | Ejder Muhafızı | elite | Ejder Yamaçları | 14210 | 998 | 385 | 145645 | 3g 21s 92c–8g 4s 80c | 385 | Dağ Aslanı Postu |  |
| 410 | fey_wolf | Peri Kurdu | evasive | Sis Ormanı | 9880 | 1022 | 329 | 151514 | 3g 34s 7c–8g 35s 18c | 410 | Dağ Aslanı Postu | dalga ×1 |
| 420 | mist_stalker | Sis Avcısı | crit | Sis Ormanı | 11182 | 1074 | 337 | 157470 | 3g 46s 36c–8g 65s 90c | 410 | Kara Panter Postu | dalga ×1 |
| 430 | ancient_treant | Kadim Treant | tank | Sis Ormanı | 19075 | 931 | 656 | 163512 | 3g 58s 81c–8g 97s 2c | 410 | Kara Panter Postu | dalga ×1 |
| 432 | spirit_queen_sylwen | Ruh Kraliçesi Sylwen | overlord | Sis Ormanı | 20258 | 1216 | 433 | 164730 | 3g 61s 31c–9g 3s 28c | 410 | Kara Panter Postu | boss Kayıp Koru |
| 440 | shadow_elf | Gölge Elfi | glass | Sis Ormanı | 10035 | 1428 | 265 | 169638 | 3g 71s 40c–9g 28s 50c | 435 | Kara Panter Postu |  |
| 450 | forest_spirit | Orman Ruhu | elite | Sis Ormanı | 15960 | 1119 | 433 | 175849 | 3g 84s 13c–9g 60s 32c | 435 | Fırtına Kartalı Tüyü |  |
| 460 | crystal_spider | Kristal Örümcek | crit | Kristal Mağaralar | 12232 | 1173 | 369 | 182143 | 3g 97s–9g 92s 50c | 460 | Fırtına Kartalı Tüyü | dalga ×1 |
| 470 | stone_eater | Taş Yiyen | tank | Kristal Mağaralar | 20825 | 1015 | 716 | 188520 | 4g 10s 2c–10g 25s 5c | 460 | Fırtına Kartalı Tüyü | dalga ×1 |
| 480 | light_elemental | Işık Elementali | glass | Kristal Mağaralar | 10935 | 1554 | 289 | 194978 | 4g 23s 17c–10g 57s 92c | 460 | Fırtına Kartalı Tüyü | dalga ×1 |
| 482 | crystal_titan | Kristal Titan | overlord | Kristal Mağaralar | 22570 | 1352 | 483 | 196280 | 4g 25s 82c–10g 64s 55c | 460 | Fırtına Kartalı Tüyü | boss Kristal Kalp |
| 490 | crystal_dragon | Kristal Ejder | evasive | Kristal Mağaralar | 11780 | 1216 | 393 | 201518 | 4g 36s 47c–10g 91s 18c | 485 | Fırtına Kartalı Tüyü |  |
| 500 | cave_titan | Mağara Titanı | elite | Kristal Mağaralar | 17710 | 1240 | 481 | 208138 | 4g 49s 90c–11g 24s 75c | 485 | Kor Domuzu Postu |  |
| 510 | magma_golem | Magma Golemi | tank | Volkanik Ada | 22575 | 1099 | 777 | 214839 | 4g 63s 46c–11g 58s 65c | 510 | Kor Domuzu Postu | dalga ×1 |
| 520 | fire_giant | Ateş Devi | normal | Volkanik Ada | 13150 | 1120 | 417 | 221618 | 4g 77s 16c–11g 92s 90c | 510 | Kor Domuzu Postu | dalga ×1 |
| 530 | ash_dragon | Kül Ejderi | crit | Volkanik Ada | 14070 | 1346 | 425 | 228476 | 4g 90s 99c–12g 27s 48c | 510 | Kor Domuzu Postu | dalga ×1 |
| 532 | fire_lord_ignar | Ateş Lordu Ignar | overlord | Volkanik Ada | 24882 | 1489 | 533 | 229857 | 4g 93s 77c–12g 34s 42c | 510 | Kor Domuzu Postu | boss Volkan Kalbi |
| 540 | volcano_shaman | Volkan Şamanı | glass | Volkanik Ada | 12285 | 1743 | 325 | 235413 | 5g 4s 95c–12g 62s 38c | 535 | Kor Domuzu Postu |  |
| 550 | volcano_warden | Volkan Bekçisi | elite | Volkanik Ada | 19460 | 1360 | 529 | 242427 | 5g 19s 4c–12g 97s 60c | 535 | Ayaz Kurdu Postu |  |
| 560 | dragon_knight | Ejder Şövalyesi | normal | Taht Zirveleri | 14150 | 1204 | 449 | 249517 | 5g 33s 26c–13g 33s 15c | 560 | Ayaz Kurdu Postu | dalga ×1 |
| 570 | ancient_wyvern | Kadim Wyvern | evasive | Taht Zirveleri | 13680 | 1409 | 457 | 256685 | 5g 47s 61c–13g 69s 2c | 560 | Ayaz Kurdu Postu | dalga ×1 |
| 580 | dragon_sorcerer | Ejder Büyücüsü | glass | Taht Zirveleri | 13185 | 1869 | 349 | 263928 | 5g 62s 8c–14g 5s 20c | 560 | Ayaz Kurdu Postu | dalga ×2 |
| 582 | dragon_king_azhorath | Ejder Kral Azhorath | overlord | Taht Zirveleri | 27195 | 1625 | 583 | 265385 | 5g 64s 99c–14g 12s 48c | 560 | Ayaz Kurdu Postu | boss Ejder Tahtı |
| 590 | throne_guardian | Taht Muhafızı | tank | Taht Zirveleri | 26075 | 1267 | 899 | 271246 | 5g 76s 68c–14g 41s 70c | 585 | Ayaz Kurdu Postu | dalga ×1 |
| 600 | elder_dragon | Yaşlı Ejder | elite | Taht Zirveleri | 21210 | 1481 | 577 | 278639 | 5g 91s 40c–14g 78s 50c | 585 | Dev Ayı Postu | dalga ×1 |
| 602 | azhorath_ascended | Yükselmiş Azhorath | elder_titan | Ejder Kral'ın İni | 456000 | 6849 | 676 | 280127 | 5g 94s 36c–14g 85s 90c | 585 | Dev Ayı Postu | baskın boss’u Ejder Kral'ın İni |
| 610 | shadow_rider | Gölge Atlı | normal | Karanlık Ova | 15400 | 1309 | 489 | 286107 | 6g 6s 25c–15g 15s 62c | 610 | Dev Ayı Postu | dalga ×1 |
| 620 | void_wolf | Boşluk Kurdu | evasive | Karanlık Ova | 14868 | 1529 | 497 | 293648 | 6g 21s 22c–15g 53s 5c | 610 | Dev Ayı Postu | dalga ×1 |
| 630 | star_gnawer | Yıldız Kemiren | glass | Karanlık Ova | 14310 | 2026 | 379 | 301262 | 6g 36s 31c–15g 90s 78c | 610 | Dev Ayı Postu | dalga ×1 |
| 632 | star_eater_omh | Yıldız Yiyen Omh | overlord | Karanlık Ova | 29508 | 1762 | 633 | 302794 | 6g 39s 34c–15g 98s 35c | 610 | Dev Ayı Postu | boss Söndürülmüş Gözlemevi |
| 640 | dark_angel | Kara Melek | crit | Karanlık Ova | 16958 | 1619 | 513 | 308950 | 6g 51s 52c–16g 28s 80c | 635 | Dev Ayı Postu |  |
| 650 | night_giant | Gece Devi | elite | Karanlık Ova | 22960 | 1602 | 625 | 316710 | 6g 66s 85c–16g 67s 12c | 635 | Ruh Geyiği Postu |  |
| 660 | ghost_noble | Hayalet Soylu | normal | Kayıp Şehir | 16650 | 1414 | 529 | 324542 | 6g 82s 30c–17g 5s 75c | 660 | Ruh Geyiği Postu | dalga ×1 |
| 670 | time_thief | Zaman Hırsızı | evasive | Kayıp Şehir | 16055 | 1650 | 537 | 332445 | 6g 97s 86c–17g 44s 65c | 660 | Ruh Geyiği Postu | dalga ×1 |
| 680 | cursed_sorcerer | Lanetli Büyücü | glass | Kayıp Şehir | 15435 | 2184 | 409 | 340419 | 7g 13s 54c–17g 83s 85c | 660 | Ruh Geyiği Postu | dalga ×2 |
| 682 | time_lord_kaedrin | Zaman Lordu Kaedrin | overlord | Kayıp Şehir | 31820 | 1898 | 683 | 342023 | 7g 16s 70c–17g 91s 75c | 660 | Ruh Geyiği Postu | boss Saat Kulesi |
| 690 | stone_gargoyle | Taş Gargoyl | tank | Kayıp Şehir | 30450 | 1477 | 1051 | 348465 | 7g 29s 34c–18g 23s 35c | 685 | Ruh Geyiği Postu | dalga ×1 |
| 700 | city_guardian | Şehir Koruyucusu | elite | Kayıp Şehir | 24710 | 1723 | 673 | 356580 | 7g 45s 26c–18g 63s 15c | 685 | Kadim Mamut Postu | dalga ×1 |
| 702 | kaedrin_timeless | Zamansız Kaedrin | elder_titan | Zaman Lordu'nun Kulesi | 531000 | 7962 | 788 | 358211 | 7g 48s 45c–18g 71s 12c | 685 | Kadim Mamut Postu | baskın boss’u Zaman Lordu'nun Kulesi |
| 710 | cloud_giant | Bulut Devi | tank | Uçan Adalar | 31325 | 1519 | 1081 | 364765 | 7g 61s 28c–19g 3s 20c | 710 | Kadim Mamut Postu | dalga ×1 |
| 720 | storm_dragon | Fırtına Ejderi | crit | Uçan Adalar | 19058 | 1817 | 577 | 373020 | 7g 77s 42c–19g 43s 55c | 710 | Kadim Mamut Postu | dalga ×1 |
| 730 | sky_knight | Gök Şövalyesi | normal | Uçan Adalar | 18400 | 1561 | 585 | 381344 | 7g 93s 67c–19g 84s 18c | 710 | Kadim Mamut Postu | dalga ×1 |
| 732 | storm_king_zephyros | Fırtına Kralı Zephyros | overlord | Uçan Adalar | 34132 | 2035 | 733 | 383017 | 7g 96s 94c–19g 92s 35c | 710 | Kadim Mamut Postu | boss Fırtına Kalesi |
| 740 | wind_elemental | Rüzgâr Elementali | evasive | Uçan Adalar | 17718 | 1819 | 593 | 389736 | 8g 10s 4c–20g 25s 10c | 735 | Kadim Mamut Postu |  |
| 750 | thunderbird | Yıldırım Kuşu | elite | Uçan Adalar | 26460 | 1843 | 721 | 398197 | 8g 26s 51c–20g 66s 28c | 735 | Gölge Kimera Derisi |  |
| 760 | void_creature | Boşluk Yaratığı | normal | Yarık Kenarı | 19150 | 1624 | 609 | 406726 | 8g 43s 10c–21g 7s 75c | 760 | Gölge Kimera Derisi | dalga ×1 |
| 770 | rift_warden | Yarık Bekçisi | tank | Yarık Kenarı | 33950 | 1645 | 1172 | 415322 | 8g 59s 79c–21g 49s 48c | 760 | Gölge Kimera Derisi | dalga ×1 |
| 780 | chaos_elemental | Kaos Elementali | glass | Yarık Kenarı | 17685 | 2499 | 469 | 423986 | 8g 76s 60c–21g 91s 50c | 760 | Gölge Kimera Derisi | dalga ×1 |
| 782 | chaos_mother | Kaos Ana | overlord | Yarık Kenarı | 36445 | 2171 | 783 | 425727 | 8g 79s 97c–21g 99s 92c | 760 | Gölge Kimera Derisi | boss Yarık Çekirdeği |
| 790 | rift_dragon | Yarık Ejderi | crit | Yarık Kenarı | 20895 | 1991 | 633 | 432716 | 8g 93s 51c–22g 33s 78c | 785 | Gölge Kimera Derisi |  |
| 800 | chaos_lord | Kaos Lordu | elite | Yarık Kenarı | 28210 | 1964 | 769 | 441513 | 9g 10s 53c–22g 76s 32c | 785 | Yaşlı Wyvern Pulu |  |
| 810 | fallen_angel | Düşmüş Melek | normal | Göksel Tapınak | 20400 | 1729 | 649 | 450377 | 9g 27s 65c–23g 19s 12c | 810 | Yaşlı Wyvern Pulu | dalga ×1 |
| 820 | temple_guardian | Tapınak Muhafızı | tank | Göksel Tapınak | 36138 | 1750 | 1248 | 459306 | 9g 44s 88c–23g 62s 20c | 810 | Yaşlı Wyvern Pulu | dalga ×1 |
| 830 | celestial_dragon | Göksel Ejder | crit | Göksel Tapınak | 21945 | 2090 | 665 | 468301 | 9g 62s 22c–24g 5s 55c | 810 | Yaşlı Wyvern Pulu | dalga ×2 |
| 832 | archangel_uriel | Başmelek Uriel | overlord | Göksel Tapınak | 38758 | 2308 | 833 | 470108 | 9g 65s 70c–24g 14s 25c | 810 | Yaşlı Wyvern Pulu | boss Göksel Kubbe |
| 840 | light_hunter | Işık Avcısı | evasive | Göksel Tapınak | 20092 | 2061 | 673 | 477361 | 9g 79s 66c–24g 49s 15c | 835 | Yaşlı Wyvern Pulu | dalga ×1 |
| 850 | seraph | Seraf | elite | Göksel Tapınak | 29960 | 2085 | 817 | 486486 | 9g 97s 21c–24g 93s 2c | 835 | Yıldız Grifonu Tüyü | dalga ×1 |
| 852 | fallen_archangels | Düşmüş Başmelekler | elder_titan | Başmeleklerin Düşüşü | 643500 | 9631 | 956 | 488319 | 10g 73c–25g 1s 82c | 835 | Yıldız Grifonu Tüyü | baskın boss’u Başmeleklerin Düşüşü |
| 860 | ash_leviathan | Kül Leviathanı | tank | Kül Denizi | 37888 | 1834 | 1309 | 495676 | 10g 14s 86c–25g 37s 15c | 860 | Yıldız Grifonu Tüyü | dalga ×1 |
| 870 | null_hunter | Hiçlik Avcısı | evasive | Kül Denizi | 20805 | 2133 | 697 | 504930 | 10g 32s 61c–25g 81s 52c | 860 | Yıldız Grifonu Tüyü | dalga ×1 |
| 880 | undying_soldier | Ölümsüz Asker | normal | Kül Denizi | 22150 | 1876 | 705 | 514248 | 10g 50s 47c–26g 26s 18c | 860 | Yıldız Grifonu Tüyü | dalga ×1 |
| 882 | leviathan_king | Leviathan Kral | overlord | Kül Denizi | 41070 | 2444 | 883 | 516119 | 10g 54s 5c–26g 35s 12c | 860 | Yıldız Grifonu Tüyü | boss Batık Taht |
| 890 | ash_wyrm_lord | Kül Ejderi Lordu | crit | Kül Denizi | 23520 | 2238 | 713 | 523629 | 10g 68s 42c–26g 71s 5c | 885 | Yıldız Grifonu Tüyü |  |
| 900 | eternal_sentinel | Ebedi Bekçi | elite | Kül Denizi | 31710 | 2206 | 865 | 533075 | 10g 86s 48c–27g 16s 20c | 885 | Boşluk Behemotu Derisi |  |
| 910 | doom_rider | Kıyamet Atlısı | crit | Son Kale Surları | 24045 | 2288 | 729 | 542583 | 11g 4s 64c–27g 61s 60c | 910 | Boşluk Behemotu Derisi | dalga ×1 |
| 920 | void_titan | Boşluk Titanı | tank | Son Kale Surları | 40512 | 1960 | 1400 | 552154 | 11g 22s 90c–28g 7s 25c | 910 | Boşluk Behemotu Derisi | dalga ×1 |
| 930 | soul_reaper | Ruh Biçen | glass | Son Kale Surları | 21060 | 2972 | 559 | 561788 | 11g 41s 25c–28g 53s 12c | 910 | Boşluk Behemotu Derisi | dalga ×1 |
| 932 | doom_commander_morvath | Kıyamet Komutanı Morvath | overlord | Son Kale Surları | 43382 | 2581 | 933 | 563723 | 11g 44s 94c–28g 62s 35c | 910 | Boşluk Behemotu Derisi | boss Son Kale |
| 940 | ancient_golem | Kadim Golem | normal | Son Kale Surları | 23650 | 2002 | 753 | 571485 | 11g 59s 71c–28g 99s 28c | 935 | Boşluk Behemotu Derisi |  |
| 950 | last_commander | Son Komutan | elite | Son Kale Surları | 33460 | 2326 | 913 | 581243 | 11g 78s 27c–29g 45s 68c | 935 | Kadim Anka Tüyü |  |
| 960 | void_dragon | Hiçlik Ejderi | crit | Ebediyet Kıyısı | 25358 | 2412 | 769 | 591063 | 11g 96s 92c–29g 92s 30c | 960 | Kadim Anka Tüyü | dalga ×1 |
| 970 | throne_sentinel | Taht Gözcüsü | tank | Ebediyet Kıyısı | 42700 | 2065 | 1476 | 600945 | 12g 15s 67c–30g 39s 18c | 960 | Kadim Anka Tüyü | dalga ×2 |
| 980 | eternal_mage | Ebedi Büyücü | glass | Ebediyet Kıyısı | 22185 | 3129 | 589 | 610888 | 12g 34s 52c–30g 86s 30c | 960 | Kadim Anka Tüyü | dalga ×1 |
| 982 | fallen_god_aethon | Düşmüş Tanrı Aethon | overlord | Ebediyet Kıyısı | 45695 | 2717 | 983 | 612884 | 12g 38s 30c–30g 95s 75c | 960 | Kadim Anka Tüyü | boss Ebedi Taht |
| 990 | god_shadow | Tanrı Gölgesi | evasive | Ebediyet Kıyısı | 23655 | 2423 | 793 | 620892 | 12g 53s 46c–31g 33s 65c | 985 | Kadim Anka Tüyü | dalga ×1 |
| 1000 | aethon_unbound | Zincirsiz Aethon | elder_titan | Düşmüş Tanrı'nın Tahtı | 754500 | 11278 | 1121 | 630957 | 12g 72s 50c–31g 81s 25c | 985 | Göksel Ejder Pulu | baskın boss’u Düşmüş Tanrı'nın Tahtı |
| 1000 | avatar_of_ruin | Yıkım Avatarı | elite | Ebediyet Kıyısı | 35210 | 2447 | 961 | 630957 | 12g 72s 50c–31g 81s 25c | 985 | Göksel Ejder Pulu | dalga ×1 |

En geniş canavar seviye boşlukları: Sv 990→1000, Sv 970→980, Sv 960→970.

Arketip dağılımı: elite 25, tank 20, crit 18, glass 17, evasive 16, el yapımı 14, overlord 14, normal 13, elder_titan 4, chief 3, titan 2.

## 6. Zindanlar ve baskınlar

| Sv | Kod | Ad | Kişi | Süre | Dalgalar | Boss | Boss Sv | Bonus altın | Bonus XP | Boss ekipman kademesi |
|---:|---|---|---|---:|---|---|---:|---|---:|---:|
| 6 | crypt_of_ashes | Kül Mahzeni | 1 | 15 dk | 3× Boz Kurt, 3× Haydut Gözcü, 3× Mağara Yarasası | Goblin Yağmacı | 8 | 2s 50c | 300 | 5 |
| 11 | forest_hollow | Orman Kovuğu | 1 | 20 dk | 3× Goblin Yağmacı, 3× Bataklık Pusucusu, 2× Orman Trolü | Dev Kurt | 14 | 7s | 1200 | 10 |
| 17 | shadowfen_ruins | Gölgebataklık Harabeleri | 1 | 30 dk | 3× Dev Kurt, 3× Tarikat Müridi, 2× Kaya Devi | Gölge Avcısı | 20 | 25s | 4000 | 20 |
| 23 | obsidian_spire | Obsidyen Kule | 1 | 40 dk | 3× Gölge Avcısı, 3× Buz Devi, 1× Kadim Ejder | Kül Ölüsü | 25 | 50s | 9000 | 25 |
| 38 | sunken_temple | Batık Tapınak | 1 | 45 dk | 3× Sis Sülüğü, 3× Çamur Golemi, 2× Bataklık Cadısı | Bataklık Kraliçesi | 42 | 80s | 12000 | 35 |
| 63 | forgotten_mine | Unutulmuş Maden | 1 | 50 dk | 3× Dağ Trolü, 3× Cüce Hayaleti, 2× Kaya Kartalı | Demir Kral | 67 | 1g 50s | 20000 | 65 |
| 88 | northern_keep | Kuzey Kalesi | 1 | 60 dk | 3× Hayalet Muhafız, 3× Lanetli Okçu, 2× Kemik Atlı | Düşmüş Lord Varn | 92 | 2g 50s | 30000 | 80 |
| 100 | north_keep_siege | Kuzey Kalesi Kuşatması | 3–5 | 60 dk | 3× Kale Gargoyli, 3× Düşmüş Komutan, 2× Kemik Atlı | Ölümsüz Varn | 102 | 9g 38s 3c | 117264 | 95 |
| 128 | rebel_camp | Asiler Kampı | 1 | 60 dk | 3× Asi Süvari, 3× Ova Aslanı, 2× Yağmacı Büyücü | Asi Kral Orrin | 132 | 4g | 50000 | 110 |
| 178 | fire_root | Ateş Kökü | 1 | 60 dk | 3× Kül Kurdu, 3× Alev Ruhu, 2× Kömür Treanti | Yanık Ana Ağaç | 182 | 6g | 80000 | 160 |
| 228 | buried_palace | Gömülü Saray | 1 | 60 dk | 3× Kum Akrebi, 3× Çöl Yılanı, 2× Tuz Elementali | Firavun Seth-Ka | 232 | 8g 50s | 120000 | 210 |
| 278 | royal_crypt | Kraliyet Kriptası | 1 | 60 dk | 3× Kraliyet İskeleti, 3× Mezar Bekçisi, 2× Ruh Emici | Mezar Lordu Malachar | 282 | 11g 50s | 170000 | 260 |
| 300 | crownless_palace | Taçsız Kral'ın Sarayı | 3–5 | 60 dk | 3× Mezar Bekçisi, 3× Kara Şövalye, 2× Taçsız Kral | Taçsız Tiran | 302 | 64g 14s 74c | 1139736 | 285 |
| 328 | sunken_fleet | Batık Filo | 1 | 60 dk | 3× Deniz Haydudu, 3× Fırtına Kuşu, 2× Mercan Golemi | Amiral Drask | 332 | 15g | 230000 | 310 |
| 378 | brood_cavern | Kuluçka Mağarası | 1 | 60 dk | 3× Genç Ejder, 3× Ejder Kültisti, 2× Lav Kertenkelesi | Ana Ejder Vyrna | 382 | 19g | 320000 | 360 |
| 428 | lost_grove | Kayıp Koru | 1 | 60 dk | 3× Peri Kurdu, 3× Sis Avcısı, 2× Kadim Treant | Ruh Kraliçesi Sylwen | 432 | 24g | 430000 | 410 |
| 478 | crystal_heart | Kristal Kalp | 1 | 60 dk | 3× Kristal Örümcek, 3× Taş Yiyen, 2× Işık Elementali | Kristal Titan | 482 | 30g | 560000 | 460 |
| 528 | volcano_heart | Volkan Kalbi | 1 | 60 dk | 3× Magma Golemi, 3× Ateş Devi, 2× Kül Ejderi | Ateş Lordu Ignar | 532 | 37g | 720000 | 510 |
| 578 | dragon_throne | Ejder Tahtı | 1 | 60 dk | 3× Ejder Şövalyesi, 3× Kadim Wyvern, 2× Ejder Büyücüsü | Ejder Kral Azhorath | 582 | 45g | 900000 | 560 |
| 600 | dragon_king_lair | Ejder Kral'ın İni | 3–5 | 60 dk | 3× Taht Muhafızı, 3× Ejder Büyücüsü, 2× Yaşlı Ejder | Yükselmiş Azhorath | 602 | 215g 76s 51c | 4785601 | 585 |
| 628 | darkened_observatory | Söndürülmüş Gözlemevi | 1 | 60 dk | 3× Gölge Atlı, 3× Boşluk Kurdu, 2× Yıldız Kemiren | Yıldız Yiyen Omh | 632 | 55g | 1150000 | 610 |
| 678 | clock_tower | Saat Kulesi | 1 | 60 dk | 3× Hayalet Soylu, 3× Zaman Hırsızı, 2× Lanetli Büyücü | Zaman Lordu Kaedrin | 682 | 65g | 1400000 | 660 |
| 700 | time_lord_tower | Zaman Lordu'nun Kulesi | 3–5 | 60 dk | 3× Taş Gargoyl, 3× Lanetli Büyücü, 2× Şehir Koruyucusu | Zamansız Kaedrin | 702 | 282g 57s 79c | 6584402 | 685 |
| 728 | storm_fortress | Fırtına Kalesi | 1 | 60 dk | 3× Bulut Devi, 3× Fırtına Ejderi, 2× Gök Şövalyesi | Fırtına Kralı Zephyros | 732 | 78g | 1700000 | 710 |
| 778 | rift_core | Yarık Çekirdeği | 1 | 60 dk | 3× Boşluk Yaratığı, 3× Yarık Bekçisi, 2× Kaos Elementali | Kaos Ana | 782 | 92g | 2100000 | 760 |
| 828 | celestial_dome | Göksel Kubbe | 1 | 60 dk | 3× Düşmüş Melek, 3× Tapınak Muhafızı, 2× Göksel Ejder | Başmelek Uriel | 832 | 108g | 2600000 | 810 |
| 850 | archangels_fall | Başmeleklerin Düşüşü | 3–5 | 60 dk | 3× Işık Avcısı, 3× Göksel Ejder, 2× Seraf | Düşmüş Başmelekler | 852 | 396g 91s 70c | 9841483 | 835 |
| 878 | sunken_throne | Batık Taht | 1 | 60 dk | 3× Kül Leviathanı, 3× Hiçlik Avcısı, 2× Ölümsüz Asker | Leviathan Kral | 882 | 125g | 3100000 | 860 |
| 928 | bastion_keep | Son Kale | 1 | 60 dk | 3× Kıyamet Atlısı, 3× Boşluk Titanı, 2× Ruh Biçen | Kıyamet Komutanı Morvath | 932 | 145g | 3700000 | 910 |
| 978 | eternal_throne | Ebedi Taht | 1 | 60 dk | 3× Hiçlik Ejderi, 3× Taht Gözcüsü, 2× Ebedi Büyücü | Düşmüş Tanrı Aethon | 982 | 170g | 4400000 | 960 |
| 1000 | fallen_god_throne | Düşmüş Tanrı'nın Tahtı | 3–5 | 60 dk | 3× Taht Gözcüsü, 3× Tanrı Gölgesi, 2× Yıkım Avatarı | Zincirsiz Aethon | 1000 | 527g 49s 29c | 13777277 | 985 |

## 7. Setler

Baskın setleri (`scripts/gen_sets.py`, DD-277): baskın boss’u sınıfın zırh türündeki setten bir parça düşürür.

| Set | Ad | Baskın | Sv | Parçalar | Bonuslar |
|---|---|---|---:|---|---|
| north_keep_siege_cloth | Ölümsüz Kumaş Seti | Kuzey Kalesi Kuşatması | 100 | Ölümsüz Kukuleta, Ölümsüz Cübbe, Ölümsüz Şalvar, Ölümsüz Sandalet, Ölümsüz Sargı, Ölümsüz Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| north_keep_siege_leather | Ölümsüz Deri Seti | Kuzey Kalesi Kuşatması | 100 | Ölümsüz Deri Başlık, Ölümsüz Deri Yelek, Ölümsüz Deri Tozluk, Ölümsüz Deri Çizme, Ölümsüz Deri Eldiven, Ölümsüz Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| north_keep_siege_plate | Ölümsüz Plaka Seti | Kuzey Kalesi Kuşatması | 100 | Ölümsüz Kask, Ölümsüz Göğüslük, Ölümsüz Bacaklık, Ölümsüz Çizme, Ölümsüz Eldiven, Ölümsüz Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |
| crownless_palace_cloth | Tiran Kumaş Seti | Taçsız Kral'ın Sarayı | 300 | Tiran Kukuleta, Tiran Cübbe, Tiran Şalvar, Tiran Sandalet, Tiran Sargı, Tiran Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| crownless_palace_leather | Tiran Deri Seti | Taçsız Kral'ın Sarayı | 300 | Tiran Deri Başlık, Tiran Deri Yelek, Tiran Deri Tozluk, Tiran Deri Çizme, Tiran Deri Eldiven, Tiran Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| crownless_palace_plate | Tiran Plaka Seti | Taçsız Kral'ın Sarayı | 300 | Tiran Kask, Tiran Göğüslük, Tiran Bacaklık, Tiran Çizme, Tiran Eldiven, Tiran Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |
| dragon_king_lair_cloth | Ejder Kral Kumaş Seti | Ejder Kral'ın İni | 600 | Ejder Kral Kukuleta, Ejder Kral Cübbe, Ejder Kral Şalvar, Ejder Kral Sandalet, Ejder Kral Sargı, Ejder Kral Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| dragon_king_lair_leather | Ejder Kral Deri Seti | Ejder Kral'ın İni | 600 | Ejder Kral Deri Başlık, Ejder Kral Deri Yelek, Ejder Kral Deri Tozluk, Ejder Kral Deri Çizme, Ejder Kral Deri Eldiven, Ejder Kral Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| dragon_king_lair_plate | Ejder Kral Plaka Seti | Ejder Kral'ın İni | 600 | Ejder Kral Kask, Ejder Kral Göğüslük, Ejder Kral Bacaklık, Ejder Kral Çizme, Ejder Kral Eldiven, Ejder Kral Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |
| time_lord_tower_cloth | Zamansız Kumaş Seti | Zaman Lordu'nun Kulesi | 700 | Zamansız Kukuleta, Zamansız Cübbe, Zamansız Şalvar, Zamansız Sandalet, Zamansız Sargı, Zamansız Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| time_lord_tower_leather | Zamansız Deri Seti | Zaman Lordu'nun Kulesi | 700 | Zamansız Deri Başlık, Zamansız Deri Yelek, Zamansız Deri Tozluk, Zamansız Deri Çizme, Zamansız Deri Eldiven, Zamansız Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| time_lord_tower_plate | Zamansız Plaka Seti | Zaman Lordu'nun Kulesi | 700 | Zamansız Kask, Zamansız Göğüslük, Zamansız Bacaklık, Zamansız Çizme, Zamansız Eldiven, Zamansız Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |
| archangels_fall_cloth | Düşmüş Kumaş Seti | Başmeleklerin Düşüşü | 850 | Düşmüş Kukuleta, Düşmüş Cübbe, Düşmüş Şalvar, Düşmüş Sandalet, Düşmüş Sargı, Düşmüş Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| archangels_fall_leather | Düşmüş Deri Seti | Başmeleklerin Düşüşü | 850 | Düşmüş Deri Başlık, Düşmüş Deri Yelek, Düşmüş Deri Tozluk, Düşmüş Deri Çizme, Düşmüş Deri Eldiven, Düşmüş Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| archangels_fall_plate | Düşmüş Plaka Seti | Başmeleklerin Düşüşü | 850 | Düşmüş Kask, Düşmüş Göğüslük, Düşmüş Bacaklık, Düşmüş Çizme, Düşmüş Eldiven, Düşmüş Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |
| fallen_god_throne_cloth | Zincirsiz Kumaş Seti | Düşmüş Tanrı'nın Tahtı | 1000 | Zincirsiz Kukuleta, Zincirsiz Cübbe, Zincirsiz Şalvar, Zincirsiz Sandalet, Zincirsiz Sargı, Zincirsiz Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_dmg_pct 10 |
| fallen_god_throne_leather | Zincirsiz Deri Seti | Düşmüş Tanrı'nın Tahtı | 1000 | Zincirsiz Deri Başlık, Zincirsiz Deri Yelek, Zincirsiz Deri Tozluk, Zincirsiz Deri Çizme, Zincirsiz Deri Eldiven, Zincirsiz Gerdanlık | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat crit_pct 3 |
| fallen_god_throne_plate | Zincirsiz Plaka Seti | Düşmüş Tanrı'nın Tahtı | 1000 | Zincirsiz Kask, Zincirsiz Göğüslük, Zincirsiz Bacaklık, Zincirsiz Çizme, Zincirsiz Eldiven, Zincirsiz Zincir Kolye | 2: mult_add hp_max 0.04; 4: mult_add attack 0.05; 6: flat dmg_reduction_pct 3 |

## 8. Uyumsuzluklar ve eksikler

| # | Önem | Alan | Bulgu |
|---:|---|---|---|
| 1 | Orta | Eşya motoru | Tasarımı olup uygulanmayan: soket + mücevher (`gems.md`, `sockets.md`), büyü (`enchantments.md`), üretim kalitesi / kritik üretim (`quality.md`, `critical-craft.md`), benzersiz etki (`unique-effect.md`). |
| 2 | Orta | Ganimet | Canavara özgü ganimet tablosu yok: aynı seviyedeki ejderha ile haydut aynı şeyi düşürüyor (yalnız seviye belirliyor). Görev/zanaat için özel malzeme (ejder pulu, boss mührü) yok. |
| 3 | Orta | Ganimet | Ganimet tek kademeden: Sv 34 canavar yalnız Sv 25 kademesi düşürür; bir alt kademe/üst kademe şansı yok. Kademe sınırını geçince düşen eşya bir anda değişiyor. |
| 4 | Orta | Nadirlik | Mitik ekipman hiçbir kaynaktan gelmiyor (boss tablosu efsanevide bitiyor); `unique-effect.md` tasarımı uygulanmadı. |
| 5 | Orta | Set | Solo zindanların seti yok; setler yalnız baskınlarda. |
| 6 | Orta | Zindan | Solo zindana özgü ganimet yok: boss ekipmanı açık dünyadaki aynı kademe (baskınlar set parçası düşürür, DD-277). |
| 7 | Düşük | Canavar | Her bölgede 6 canavar → Sv 100 sonrası her ~10 seviyede bir canavar; yan yana seçenek az (aynı seviyede ikinci canavar yok). |
| 8 | Düşük | Ekonomi | 4/150 toplama kaynağı hiçbir tarifte kullanılmıyor (forest_deer_hide, forest_bear_hide, raw_stone, mercury_ore); tek değerleri dükkân fiyatı. |
| 9 | Düşük | Eşya | 6 eşya bilinçli olarak emekli (DD-241/244; eski yığınlar satılabilsin diye duruyor), kaynağı yok: potion_mp_small, linseed, mandrake, king_bloom, elm_log, yew_log. Temizlik isteğe bağlı. |
