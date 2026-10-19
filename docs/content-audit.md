# İçerik Denetimi — eşya, üretim, ganimet, kademe, set, canavar

> Otomatik üretilir: `python scripts/content_audit.py`. Elle düzenleme — betiği değiştir. Kurallar `server/internal/item` (GearPool, CombatHides, TierRarity) ve `activity.DropLevel` aynasıdır. Yol haritası: [content-roadmap.md](content-roadmap.md).

## 0. Özet

| Alan | Adet | Aralık |
|---|---:|---|
| Ekipman şablonu | 706 | Sv 1–985, 47 kademe |
| Toplama kaynağı | 150 | meslek Sv 1–100, kademe T0–T7 |
| Diğer eşya (külçe, iksir, yemek, eski malzeme) | 53 | |
| Tarif | 67 | smelting 24, blacksmithing 20, cooking 18, alchemy 5 |
| Canavar | 140 | Sv 1–1000 (+6 baskın boss’u) |
| Zindan / baskın | 25 / 6 | Sv 6–1000 |
| Görev | 18 | Sv 1–10 |
| Set | 0 | tasarım var (`docs/item-engine/set.md`), uygulanmadı |
| Dükkân | 1 | 18 eşya |

## 1. Ekipman kademeleri (tier)

Her kademe: 10 sınıf silahı (sınıfa kilitli) + 5 zırh parçası (sınıfsız). Kaynak sütunu: D = dükkân, Ü = üretim, G = ganimet.

| Sv | Kademe (göğüs adı) | Silah | Zırh | Saldırı (kılıç) | Savunma (göğüs) | Kaynak | Bu kademeyi düşüren canavar Sv |
|---:|---|---:|---:|---:|---:|---|---|
| 1 | Deri Göğüslük | 10 | 5 | 6 | 2 | DGÜ | 1–3 |
| 5 | Demir Göğüslük | 10 | 5 | 16 | 4 | GÜ | 5–8 |
| 10 | Pul Göğüslük | 10 | 5 | 29 | 6 | GÜ | 10–14 |
| 15 | Çelik Göğüslük | 10 | 5 | 42 | 7 | G | 15–18 |
| 20 | Şövalye Göğüslüğü | 10 | 5 | 54 | 9 | G | 20–24 |
| 25 | Rünlü Göğüslük | 10 | 5 | 66 | 11 | G | 25–30 |
| 35 | Mithril Göğüslük | 10 | 5 | 88 | 14 | G | 35–45 |
| 50 | Adamant Göğüslük | 10 | 5 | 120 | 19 | G | 50–60 |
| 65 | Ejder Dövme Göğüslük | 10 | 5 | 152 | 24 | G | 65–75 |
| 80 | Yıldız Dövme Göğüslük | 10 | 5 | 184 | 29 | G | 80–92 |
| 95 | Titan Göğüslük | 10 | 5 | 216 | 34 | G | 95–100 |
| 110 | Kraliyet Göğüslük | 10 | 5 | 249 | 39 | G | 110–132 |
| 135 | Güneş Dövme Göğüslük | 10 | 5 | 304 | 48 | G | 140–150 |
| 160 | Obsidyen Göğüslük | 10 | 5 | 359 | 57 | G | 160–182 |
| 185 | Kum Çeliği Göğüslük | 10 | 5 | 414 | 66 | G | 190–200 |
| 210 | Rün Bağlı Göğüslük | 10 | 5 | 469 | 75 | G | 210–232 |
| 235 | Anka Göğüslük | 10 | 5 | 524 | 84 | G | 240–250 |
| 260 | Dehşet Dövme Göğüslük | 10 | 5 | 579 | 93 | G | 260–282 |
| 285 | Hükümdar Göğüslük | 10 | 5 | 634 | 102 | G | 290–300 |
| 310 | Mercan Göğüslük | 10 | 5 | 689 | 111 | G | 310–332 |
| 335 | Ejderha Kemiği Göğüslük | 10 | 5 | 744 | 120 | G | 340–350 |
| 360 | Sis Örgü Göğüslük | 10 | 5 | 799 | 129 | G | 360–382 |
| 385 | Prizmatik Göğüslük | 10 | 5 | 854 | 138 | G | 390–400 |
| 410 | Magma Dövme Göğüslük | 10 | 5 | 909 | 147 | G | 410–432 |
| 435 | Bora Göğüslük | 10 | 5 | 964 | 156 | G | 440–450 |
| 460 | Hiçlik Göğüslük | 10 | 5 | 1019 | 165 | G | 460–482 |
| 485 | Melek Göğüslük | 10 | 5 | 1074 | 174 | G | 490–500 |
| 510 | Tutulma Göğüslük | 10 | 5 | 1129 | 183 | G | 510–532 |
| 535 | Ejder Lordu Göğüslük | 10 | 5 | 1184 | 191 | G | 540–550 |
| 560 | Göksel Göğüslük | 10 | 5 | 1239 | 200 | G | 560–582 |
| 585 | Ebedi Göğüslük | 10 | 5 | 1294 | 209 | G | 590–600 |
| 610 | Yıldızsız Göğüslük | 10 | 5 | 1349 | 218 | G | 610–632 |
| 635 | Gece Çeliği Göğüslük | 10 | 5 | 1404 | 227 | G | 640–650 |
| 660 | Zaman Dövme Göğüslük | 10 | 5 | 1459 | 236 | G | 660–682 |
| 685 | Yükselmiş Göğüslük | 10 | 5 | 1514 | 245 | G | 690–700 |
| 710 | Fırtına Doğumlu Göğüslük | 10 | 5 | 1569 | 254 | G | 710–732 |
| 735 | Yarık Dövme Göğüslük | 10 | 5 | 1624 | 263 | G | 740–750 |
| 760 | Kaos Göğüslük | 10 | 5 | 1679 | 272 | G | 760–782 |
| 785 | Gök Dövme Göğüslük | 10 | 5 | 1734 | 281 | G | 790–800 |
| 810 | Başmelek Göğüslük | 10 | 5 | 1789 | 290 | G | 810–832 |
| 835 | Şafak Göğüslük | 10 | 5 | 1844 | 299 | G | 840–850 |
| 860 | Kül Göğüslük | 10 | 5 | 1899 | 308 | G | 860–882 |
| 885 | Leviathan Göğüslük | 10 | 5 | 1954 | 317 | G | 890–900 |
| 910 | Kıyamet Dövme Göğüslük | 10 | 5 | 2009 | 326 | G | 910–932 |
| 935 | Hiçlik Dövme Göğüslük | 10 | 5 | 2064 | 335 | G | 940–950 |
| 960 | Tanrı Katili Göğüslük | 10 | 5 | 2119 | 344 | G | 960–982 |
| 985 | Kadim Çağ Göğüslük | 10 | 5 | 2174 | 353 | G | 990–1000 |

Gereksinimli (rastgele düşmeyen) eşyalar: Titan Siper Kalkanı (Sv 600).

Kademe aralığı: Sv 1–25 arası 5, 25–95 arası 10–15, 95 sonrası 25 seviye. En geniş boşluk: Sv 110→135, Sv 135→160, Sv 160→185.

### 1.1 Slot doluluğu

| Slot | Şablon | Not |
|---|---:|---|
| mainhand | 470 |  |
| offhand | 1 | yalnız 1 şablon |
| helm | 47 |  |
| chest | 47 |  |
| legs | 47 |  |
| boots | 47 |  |
| gloves | 47 |  |
| necklace | 0 | **boş** |
| ring1 | 0 | **boş** |
| ring2 | 0 | **boş** |
| trinket1 | 0 | **boş** |
| trinket2 | 0 | **boş** |

### 1.2 Silah aileleri ve zırh türleri

| Aile | Şablon |
|---|---:|
| sword | 47 |
| axe | 0 |
| mace | 47 |
| greatsword | 0 |
| greataxe | 0 |
| spear | 0 |
| shield | 1 |
| dagger | 47 |
| short_sword | 0 |
| dual_dagger | 0 |
| fist_weapon | 47 |
| throwing_knife | 0 |
| bow | 47 |
| longbow | 0 |
| crossbow | 0 |
| short_spear | 0 |
| staff | 47 |
| wand | 0 |
| orb | 0 |
| grimoire | 0 |
| tonfa | 0 |
| unarmed | 0 |
| holy_tome | 0 |
| two_handed_mace | 47 |
| sickle | 47 |
| nature_focus | 0 |
| lute | 47 |
| harp | 0 |
| flute | 0 |
| rapier | 0 |
| ritual_axe | 0 |
| totem_focus | 47 |
| cloth_armor | 0 |
| light_armor | 0 |
| medium_armor | 0 |
| heavy_armor | 0 |

### 1.3 Nadirlik ve ek (affix)

| Nadirlik | Değer çarpanı | Ek kuralı | Ekipman ganimetinde |
|---|---:|---|---|
| sıradan | ×1 | yok | %60 |
| seçkin | ×1.5 | yok | %25 · boss %50 |
| nadir | ×3 | %50 tek ek | %11 · boss %35 |
| destansı | ×8 | ön ek + %50 son ek | %4 · boss %13 |
| efsanevi | ×25 | ön + son ek | boss %2 |
| mitik | ×100 | ön + son ek (+ benzersiz etki ayrılmış) | **hiç düşmez** (yalnız T7 malzeme) |

## 2. Eşya listesi

### 2.1 Toplama kaynakları

| Kod | Ad | Tür | Meslek Sv | Kademe | Nadirlik | Değer | Kullanan tarif |
|---|---|---|---:|---:|---|---:|---|
| tiny_guppy | Minik Lepistes | fishing/fish | 1 | T0 | sıradan | 2c | Izgara Lepistes |
| river_sardine | Nehir Sardalyası | fishing/fish | 4 | T0 | sıradan | 3c | Sardalya Şiş |
| minnow | Golyan | fishing/fish | 7 | T0 | sıradan | 3c | Naneli Balık Çorbası |
| small_pike | Küçük Turna | fishing/fish | 10 | T1 | sıradan | 4c | Turna Fileto |
| common_carp | Sazan | fishing/fish | 12 | T1 | sıradan | 4c | Lavantalı Sazan |
| silver_herring | Gümüş Ringa | fishing/fish | 14 | T1 | sıradan | 5c | — |
| river_trout | Dere Alabalığı | fishing/fish | 16 | T1 | sıradan | 5c | Otlu Alabalık |
| salmon | Somon | fishing/fish | 18 | T1 | sıradan | 6c | Füme Somon |
| bass | Tatlı Su Levreği | fishing/fish | 20 | T2 | nadir-altı | 6c | Levrek Güveci |
| lake_perch | Göl Levreği | fishing/fish | 22 | T2 | nadir-altı | 7c | — |
| red_snapper | Kırmızı Mercan | fishing/fish | 24 | T2 | nadir-altı | 8c | — |
| river_catfish | Nehir Yayını | fishing/fish | 26 | T2 | nadir-altı | 9c | — |
| bluefin_tuna | Mavi Yüzgeçli Ton | fishing/fish | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| rainbow_trout | Gökkuşağı Alabalığı | fishing/fish | 30 | T3 | nadir-altı | 11c | — |
| royal_carp | Kraliyet Sazanı | fishing/fish | 33 | T3 | nadir-altı | 14c | Kraliyet Sazanı Ziyafeti |
| emerald_pike | Zümrüt Turnası | fishing/fish | 36 | T3 | nadir-altı | 16c | — |
| silver_barracuda | Gümüş Barakuda | fishing/fish | 39 | T3 | nadir-altı | 19c | — |
| golden_sturgeon | Altın Mersin | fishing/fish | 42 | T3 | nadir-altı | 23c | Mantarlı Mersin Turtası |
| crimson_swordfish | Kızıl Kılıçbalığı | fishing/fish | 45 | T4 | nadir | 28c | Fırtına Izgarası Kılıçbalığı |
| stormscale_eel | Fırtına Pullu Yılanbalığı | fishing/fish | 50 | T4 | nadir | 37c | — |
| frostfin_salmon | Ayaz Yüzgeçli Somon | fishing/fish | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası |
| obsidian_tuna | Obsidyen Ton | fishing/fish | 60 | T5 | destansı | 66c | — |
| ghost_marlin | Hayalet Marlin | fishing/fish | 65 | T5 | destansı | 88c | Ruh Marlini Tabağı |
| golden_carp | İmparatorluk Altın Koisi | fishing/fish | 70 | T5 | destansı | 1s 18c | — |
| crystal_serpentfish | Kristal Yılanbalığı | fishing/fish | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı |
| void_angler | Boşluk Fenerbalığı | fishing/fish | 80 | T6 | efsanevi | 2s 12c | — |
| astral_moonfish | Yıldız Aybalığı | fishing/fish | 85 | T6 | efsanevi | 2s 83c | Yıldız Aybalığı Lezzeti |
| phoenix_koi | Anka Koisi | fishing/fish | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti |
| abyssal_dragonfish | Uçurum Ejderbalığı | fishing/fish | 95 | T7 | mitik | 5s 7c | — |
| ancient_leviathan_scale | Göksel Leviathan | fishing/fish | 100 | T7 | mitik | 6s 79c | Leviathan Kraliyet Şöleni |
| common_grass | Sıradan Ot | herbalism/herb | 1 | T0 | sıradan | 2c | Izgara Lepistes, Küçük Can İksiri |
| dandelion | Karahindiba | herbalism/herb | 4 | T0 | sıradan | 3c | Sardalya Şiş |
| wild_mint | Yabani Nane | herbalism/herb | 7 | T0 | sıradan | 3c | Naneli Balık Çorbası, Orta Boy Can İksiri |
| chamomile | Papatya | herbalism/herb | 10 | T1 | sıradan | 4c | Büyük Boy Can İksiri, Orta Boy Can İksiri, Turna Fileto, Çırak Asası |
| lavender | Lavanta | herbalism/herb | 12 | T1 | sıradan | 4c | Lavantalı Sazan |
| sunflower | Ayçiçeği | herbalism/herb | 14 | T1 | sıradan | 5c | — |
| bitterroot | Acıkök | herbalism/herb | 16 | T1 | sıradan | 5c | Otlu Alabalık |
| healing_moss | Şifalı Yosun | herbalism/herb | 18 | T1 | sıradan | 6c | Füme Somon |
| marsh_reed | Bataklık Sazı | herbalism/herb | 20 | T2 | nadir-altı | 6c | Levrek Güveci |
| blue_thistle | Mavi Devedikeni | herbalism/herb | 22 | T2 | nadir-altı | 7c | — |
| red_sage | Kızıl Adaçayı | herbalism/herb | 24 | T2 | nadir-altı | 8c | Büyük Boy Can İksiri, Ekstra Büyük Can İksiri |
| nightshade | Güzelavrat Otu | herbalism/herb | 26 | T2 | nadir-altı | 9c | — |
| wild_ginseng | Yaban Ginsengi | herbalism/herb | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| golden_clover | Altın Yonca | herbalism/herb | 30 | T3 | nadir-altı | 11c | — |
| dreamroot | Düş Kökü | herbalism/herb | 33 | T3 | nadir-altı | 14c | Kraliyet Sazanı Ziyafeti |
| silverleaf | Gümüş Yaprak | herbalism/herb | 36 | T3 | nadir-altı | 16c | Ekstra Büyük Can İksiri, Çok Büyük Can İksiri |
| bloodthorn | Kan Dikeni | herbalism/herb | 39 | T3 | nadir-altı | 19c | — |
| ghost_mushroom | Hayalet Mantarı | herbalism/herb | 42 | T3 | nadir-altı | 23c | Mantarlı Mersin Turtası |
| stormflower | Fırtına Çiçeği | herbalism/herb | 45 | T4 | nadir | 28c | Fırtına Izgarası Kılıçbalığı |
| emberleaf | Kor Yaprağı | herbalism/herb | 50 | T4 | nadir | 37c | — |
| frostbloom | Ayaz Çiçeği | herbalism/herb | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası |
| moonshade | Ay Gölgesi | herbalism/herb | 60 | T5 | destansı | 66c | — |
| spirit_lotus | Ruh Nilüferi | herbalism/herb | 65 | T5 | destansı | 88c | Ruh Marlini Tabağı, Çok Büyük Can İksiri |
| ancient_mandrake | Kadim Adamotu | herbalism/herb | 70 | T5 | destansı | 1s 18c | — |
| dragonblood_herb | Ejderkanı Otu | herbalism/herb | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı |
| astral_bloom | Yıldız Çiçeği | herbalism/herb | 80 | T6 | efsanevi | 2s 12c | — |
| void_orchid | Boşluk Orkidesi | herbalism/herb | 85 | T6 | efsanevi | 2s 83c | Yıldız Aybalığı Lezzeti |
| celestial_lotus | Göksel Nilüfer | herbalism/herb | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti |
| phoenix_petal | Anka Çiçeği | herbalism/herb | 95 | T7 | mitik | 5s 7c | — |
| worldroot | Dünya Kökü | herbalism/herb | 100 | T7 | mitik | 6s 79c | Leviathan Kraliyet Şöleni |
| forest_rat_pelt | Orman Faresi Postu | hunting/pelt | 1 | T0 | sıradan | 2c | Deri Bacaklık, Deri Eldiven, Deri Göğüslük, Deri Kask, Deri Çizme |
| rabbit_pelt | Küçük Tavşan Postu | hunting/pelt | 4 | T0 | sıradan | 3c | — |
| field_mouse_pelt | Tarla Faresi Postu | hunting/pelt | 7 | T0 | sıradan | 3c | — |
| squirrel_pelt | Sincap Postu | hunting/pelt | 10 | T1 | sıradan | 4c | — |
| hare_pelt | Yaban Tavşanı Postu | hunting/pelt | 12 | T1 | sıradan | 4c | Av Yayı, Demir Göğüslük |
| badger_pelt | Porsuk Postu | hunting/pelt | 14 | T1 | sıradan | 5c | — |
| wild_turkey_feather | Yaban Hindisi Tüyü | hunting/feather | 16 | T1 | sıradan | 5c | — |
| forest_deer_hide | Orman Geyiği Derisi | hunting/hide | 18 | T1 | sıradan | 6c | — |
| red_fox_pelt | Kızıl Tilki Postu | hunting/pelt | 20 | T2 | nadir-altı | 6c | — |
| boar_hide | Yaban Domuzu Postu | hunting/hide | 22 | T2 | nadir-altı | 7c | — |
| mountain_goat_hide | Dağ Keçisi Postu | hunting/hide | 24 | T2 | nadir-altı | 8c | — |
| wolf_pelt | Yaban Kurdu Postu | hunting/pelt | 26 | T2 | nadir-altı | 9c | Pul Eldiven, Pul Göğüslük, Pul Çizme |
| forest_bear_hide | Orman Ayısı Postu | hunting/hide | 28 | T2 | nadir-altı | 10c | — |
| wild_buffalo_hide | Yaban Mandası Derisi | hunting/hide | 30 | T3 | nadir-altı | 11c | — |
| giant_elk_hide | Dev Geyik Postu | hunting/hide | 33 | T3 | nadir-altı | 14c | — |
| dire_wolf_pelt | Dev Kurt Postu | hunting/pelt | 36 | T3 | nadir-altı | 16c | — |
| mountain_lion_pelt | Dağ Aslanı Postu | hunting/pelt | 39 | T3 | nadir-altı | 19c | — |
| black_panther_pelt | Kara Panter Postu | hunting/pelt | 42 | T3 | nadir-altı | 23c | — |
| storm_eagle_feather | Fırtına Kartalı Tüyü | hunting/feather | 45 | T4 | nadir | 28c | — |
| ember_boar_hide | Kor Domuzu Postu | hunting/hide | 50 | T4 | nadir | 37c | — |
| frost_wolf_pelt | Ayaz Kurdu Postu | hunting/pelt | 55 | T4 | nadir | 49c | — |
| dire_bear_hide | Dev Ayı Postu | hunting/hide | 60 | T5 | destansı | 66c | — |
| spirit_stag_hide | Ruh Geyiği Postu | hunting/hide | 65 | T5 | destansı | 88c | — |
| ancient_mammoth_hide | Kadim Mamut Postu | hunting/hide | 70 | T5 | destansı | 1s 18c | — |
| shadow_chimera_hide | Gölge Kimera Derisi | hunting/hide | 75 | T6 | efsanevi | 1s 58c | — |
| elder_wyvern_scale | Yaşlı Wyvern Pulu | hunting/scale | 80 | T6 | efsanevi | 2s 12c | — |
| astral_griffin_feather | Yıldız Grifonu Tüyü | hunting/feather | 85 | T6 | efsanevi | 2s 83c | — |
| void_behemoth_hide | Boşluk Behemotu Derisi | hunting/hide | 90 | T7 | mitik | 3s 79c | — |
| ancient_phoenix_feather | Kadim Anka Tüyü | hunting/feather | 95 | T7 | mitik | 5s 7c | — |
| dragon_scale | Göksel Ejder Pulu | hunting/scale | 100 | T7 | mitik | 6s 79c | — |
| raw_stone | Ham Taş | mining/stone | 1 | T0 | sıradan | 2c | — |
| stone_ore | Taş Cevheri | mining/stone | 4 | T0 | sıradan | 3c | — |
| flint | Çakmaktaşı | mining/stone | 7 | T0 | sıradan | 3c | — |
| coal | Kömür | mining/fuel | 10 | T1 | sıradan | 4c | Adamantit Külçesi, Altın Külçesi, Alüminyum Külçesi, Bakır Külçesi, Bronz Külçesi, Demir Külçesi, Elektrum Külçesi, Gümüş Külçesi, Kalay Külçesi, Kobalt Külçesi, Krom Külçesi, Kurşun Külçesi, Mitril Külçesi, Nikel Külçesi, Orikalkum Külçesi, Osmiyum Külçesi, Paladyum Külçesi, Platin Külçesi, Rodyum Külçesi, Titanyum Külçesi, Tungsten Külçesi, Yıldız Çeliği Külçesi, Çinko Külçesi, İridyum Külçesi |
| aluminium_ore | Alüminyum Cevheri | mining/ore | 12 | T1 | sıradan | 4c | Alüminyum Külçesi |
| lead_ore | Kurşun Cevheri | mining/ore | 14 | T1 | sıradan | 5c | Kurşun Külçesi |
| iron_ore | Demir Cevheri | mining/ore | 16 | T1 | sıradan | 5c | Demir Bacaklık, Demir Eldiven, Demir Göğüslük, Demir Hançer, Demir Kask, Demir Külçesi, Demir Kılıç, Demir Topuz, Demir Çizme |
| tin_ore | Kalay Cevheri | mining/ore | 18 | T1 | sıradan | 6c | Kalay Külçesi, Pul Bacaklık, Pul Göğüslük, Pul Kask |
| copper_ore | Bakır Cevheri | mining/ore | 20 | T2 | nadir-altı | 6c | Bakır Külçesi, Pul Bacaklık, Pul Eldiven, Pul Göğüslük, Pul Kask, Pul Çizme |
| bronze_ore | Bronz Cevheri | mining/ore | 22 | T2 | nadir-altı | 7c | Bronz Külçesi |
| obsidian | Obsidyen | mining/stone | 24 | T2 | nadir-altı | 8c | — |
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
| young_sapling | Fidan Dalı | woodcutting/log | 1 | T0 | sıradan | 2c | Izgara Lepistes |
| poplar_log | Kavak Kütüğü | woodcutting/log | 4 | T0 | sıradan | 3c | Sardalya Şiş |
| willow_log | Söğüt Kütüğü | woodcutting/log | 7 | T0 | sıradan | 3c | Naneli Balık Çorbası |
| fir_log | Köknar Kütüğü | woodcutting/log | 10 | T1 | sıradan | 4c | Turna Fileto |
| pine_log | Çam Kütüğü | woodcutting/log | 12 | T1 | sıradan | 4c | Lavantalı Sazan |
| birch_log | Huş Kütüğü | woodcutting/log | 14 | T1 | sıradan | 5c | — |
| cedar_log | Sedir Kütüğü | woodcutting/log | 16 | T1 | sıradan | 5c | Otlu Alabalık |
| oak_log | Meşe Kütüğü | woodcutting/log | 18 | T1 | sıradan | 6c | Av Yayı, Demir Hançer, Demir Kılıç, Demir Topuz, Füme Somon, Çırak Asası |
| beech_log | Kayın Kütüğü | woodcutting/log | 20 | T2 | nadir-altı | 6c | Levrek Güveci |
| ash_log | Dişbudak Kütüğü | woodcutting/log | 22 | T2 | nadir-altı | 7c | — |
| maple_log | Akçaağaç Kütüğü | woodcutting/log | 24 | T2 | nadir-altı | 8c | — |
| walnut_log | Ceviz Kütüğü | woodcutting/log | 26 | T2 | nadir-altı | 9c | — |
| teak_log | Tik Kütüğü | woodcutting/log | 28 | T2 | nadir-altı | 10c | Ginsengli Ton Balığı |
| mahogany_log | Maun Kütüğü | woodcutting/log | 30 | T3 | nadir-altı | 11c | — |
| ebony_log | Abanoz Kütüğü | woodcutting/log | 33 | T3 | nadir-altı | 14c | Kraliyet Sazanı Ziyafeti |
| silver_willow_log | Gümüş Söğüt Kütüğü | woodcutting/log | 36 | T3 | nadir-altı | 16c | — |
| golden_oak_log | Altın Meşe Kütüğü | woodcutting/log | 39 | T3 | nadir-altı | 19c | — |
| blackwood_log | Kara Odun Kütüğü | woodcutting/log | 42 | T3 | nadir-altı | 23c | Mantarlı Mersin Turtası |
| crystal_birch_log | Kristal Huş Kütüğü | woodcutting/log | 45 | T4 | nadir | 28c | Fırtına Izgarası Kılıçbalığı |
| emberwood_log | Korağacı Kütüğü | woodcutting/log | 50 | T4 | nadir | 37c | — |
| frostpine_log | Ayaz Çamı Kütüğü | woodcutting/log | 55 | T4 | nadir | 49c | Ayaz Somon Çorbası |
| bloodwood_log | Kanağacı Kütüğü | woodcutting/log | 60 | T5 | destansı | 66c | — |
| spiritwood_log | Ruhağacı Kütüğü | woodcutting/log | 65 | T5 | destansı | 88c | Ruh Marlini Tabağı |
| ironwood_log | Kadim Demirağaç Kütüğü | woodcutting/log | 70 | T5 | destansı | 1s 18c | — |
| starbark_log | Yıldızkabuğu Kütüğü | woodcutting/log | 75 | T6 | efsanevi | 1s 58c | Ejderkanlı Yılanbalığı |
| moonwood_log | Ayağacı Kütüğü | woodcutting/log | 80 | T6 | efsanevi | 2s 12c | — |
| dragonheart_log | Ejderyüreği Kütüğü | woodcutting/log | 85 | T6 | efsanevi | 2s 83c | Yıldız Aybalığı Lezzeti |
| voidwood_log | Boşlukağacı Kütüğü | woodcutting/log | 90 | T7 | mitik | 3s 79c | Anka Koi Ziyafeti |
| celestial_ash_log | Göksel Dişbudak Kütüğü | woodcutting/log | 95 | T7 | mitik | 5s 7c | — |
| worldtree_log | Dünya Ağacı Kütüğü | woodcutting/log | 100 | T7 | mitik | 6s 79c | Leviathan Kraliyet Şöleni |

Kullanılmayan toplama kaynakları (hiçbir tarifte girdi değil — yalnız satılır):

| Alt tür | Kullanılmayan / toplam |
|---|---|
| feather | 4 / 4 |
| fish | 12 / 30 |
| fuel | 0 / 1 |
| herb | 10 / 30 |
| hide | 12 / 12 |
| log | 12 / 30 |
| ore | 1 / 25 |
| pelt | 9 / 12 |
| scale | 2 / 2 |
| stone | 4 / 4 |

### 2.2 Diğer eşyalar

| Kod | Ad | Sınıf | Sv | Değer | Etki | Kaynak | Kullanan tarif |
|---|---|---|---:|---:|---|---|---|
| food_grilled_guppy | Izgara Lepistes | consumable/food | 1 | 6c | +60 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_sardine_skewer | Sardalya Şiş | consumable/food | 4 | 9c | +90 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_minnow_soup | Naneli Balık Çorbası | consumable/food | 7 | 9c | +115 can, +%2 sal/sav 60 dk | üretim (cooking) | — |
| food_pike_fillet | Turna Fileto | consumable/food | 10 | 12c | +140 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_lavender_carp | Lavantalı Sazan | consumable/food | 12 | 12c | +160 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_herb_trout | Otlu Alabalık | consumable/food | 16 | 15c | +190 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_smoked_salmon | Füme Somon | consumable/food | 18 | 18c | +210 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_bass_stew | Levrek Güveci | consumable/food | 20 | 18c | +225 can, +%3 sal/sav 60 dk | üretim (cooking) | — |
| food_tuna_steak | Ginsengli Ton Balığı | consumable/food | 28 | 30c | +300 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_royal_carp_feast | Kraliyet Sazanı Ziyafeti | consumable/food | 33 | 42c | +340 can, +%4 sal/sav 60 dk | üretim (cooking) | — |
| food_sturgeon_pie | Mantarlı Mersin Turtası | consumable/food | 42 | 69c | +420 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_swordfish_grill | Fırtına Izgarası Kılıçbalığı | consumable/food | 45 | 84c | +445 can, +%5 sal/sav 60 dk | üretim (cooking) | — |
| food_frostfin_chowder | Ayaz Somon Çorbası | consumable/food | 55 | 1s 47c | +535 can, +%6 sal/sav 60 dk | üretim (cooking) | — |
| food_marlin_platter | Ruh Marlini Tabağı | consumable/food | 65 | 2s 64c | +620 can, +%7 sal/sav 60 dk | üretim (cooking) | — |
| food_dragonblood_serpentfish | Ejderkanlı Yılanbalığı | consumable/food | 75 | 4s 74c | +710 can, +%7 sal/sav 60 dk | üretim (cooking) | — |
| food_moonfish_delicacy | Yıldız Aybalığı Lezzeti | consumable/food | 85 | 8s 49c | +795 can, +%8 sal/sav 60 dk | üretim (cooking) | — |
| food_phoenix_koi_feast | Anka Koi Ziyafeti | consumable/food | 90 | 11s 37c | +840 can, +%8 sal/sav 60 dk | üretim (cooking) | — |
| food_leviathan_banquet | Leviathan Kraliyet Şöleni | consumable/food | 100 | 20s 37c | +925 can, +%9 sal/sav 60 dk | üretim (cooking) | — |
| potion_hp_small | Küçük Can İksiri | consumable/potion | 1 | 20c | +60 can | dükkân, üretim (alchemy) | — |
| potion_mp_small | Küçük Mana İksiri | consumable/potion | 1 | 20c |  | **yok** | — |
| potion_hp_medium | Orta Boy Can İksiri | consumable/potion | 5 | 60c | +150 can | dükkân, üretim (alchemy) | — |
| potion_hp_large | Büyük Boy Can İksiri | consumable/potion | 12 | 1s 50c | +280 can | dükkân, üretim (alchemy) | — |
| potion_hp_xlarge | Ekstra Büyük Can İksiri | consumable/potion | 20 | 3s 20c | +450 can | üretim (alchemy) | — |
| potion_hp_huge | Çok Büyük Can İksiri | consumable/potion | 28 | 6s | +700 can | üretim (alchemy) | — |
| adamantite_bar | Adamantit Külçesi | material/bar | 1 | 54s 58c |  | üretim (smelting) | — |
| orichalcum_bar | Orikalkum Külçesi | material/bar | 1 | 40s 82c |  | üretim (smelting) | — |
| mythril_bar | Mitril Külçesi | material/bar | 1 | 30s 58c |  | üretim (smelting) | — |
| starsteel_bar | Yıldız Çeliği Külçesi | material/bar | 1 | 18s 30c |  | üretim (smelting) | — |
| iridium_bar | İridyum Külçesi | material/bar | 1 | 13s 76c |  | üretim (smelting) | — |
| osmium_bar | Osmiyum Külçesi | material/bar | 1 | 10s 30c |  | üretim (smelting) | — |
| platinum_bar | Platin Külçesi | material/bar | 1 | 5s 86c |  | üretim (smelting) | — |
| rhodium_bar | Rodyum Külçesi | material/bar | 1 | 4s 42c |  | üretim (smelting) | — |
| palladium_bar | Paladyum Külçesi | material/bar | 1 | 3s 36c |  | üretim (smelting) | — |
| titanium_bar | Titanyum Külçesi | material/bar | 1 | 2s 48c |  | üretim (smelting) | — |
| tungsten_bar | Tungsten Külçesi | material/bar | 1 | 1s 90c |  | üretim (smelting) | — |
| electrum_bar | Elektrum Külçesi | material/bar | 1 | 1s 47c |  | üretim (smelting) | — |
| gold_bar | Altın Külçesi | material/bar | 1 | 86c |  | üretim (smelting) | — |
| silver_bar | Gümüş Külçesi | material/bar | 1 | 74c |  | üretim (smelting) | — |
| cobalt_bar | Kobalt Külçesi | material/bar | 1 | 64c |  | üretim (smelting) | — |
| chromium_bar | Krom Külçesi | material/bar | 1 | 58c |  | üretim (smelting) | — |
| nickel_bar | Nikel Külçesi | material/bar | 1 | 48c |  | üretim (smelting) | — |
| zinc_bar | Çinko Külçesi | material/bar | 1 | 38c |  | üretim (smelting) | — |
| bronze_bar | Bronz Külçesi | material/bar | 1 | 29c |  | üretim (smelting) | — |
| copper_bar | Bakır Külçesi | material/bar | 1 | 26c |  | üretim (smelting) | — |
| tin_bar | Kalay Külçesi | material/bar | 1 | 26c |  | üretim (smelting) | — |
| iron_bar | Demir Külçesi | material/bar | 1 | 22c |  | üretim (smelting) | — |
| lead_bar | Kurşun Külçesi | material/bar | 1 | 22c |  | üretim (smelting) | — |
| aluminium_bar | Alüminyum Külçesi | material/bar | 1 | 19c |  | üretim (smelting) | — |
| linseed | Keten Tohumu | material/herb | 1 | 7c |  | **yok** | — |
| mandrake | Adamotu | material/herb | 5 | 12c |  | **yok** | — |
| king_bloom | Kral Çiçeği | material/herb | 20 | 60c |  | **yok** | — |
| elm_log | Karaağaç Kütüğü | material/wood | 8 | 9c |  | **yok** | — |
| yew_log | Porsuk Kütüğü | material/wood | 20 | 30c |  | **yok** | — |

## 3. Üretim listesi

### Simya — 5 tarif, beceri 1–30 (tavan 60), çıktı Sv 1–28

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_potion_hp_small | 1 | 3× Sıradan Ot | 1× Küçük Can İksiri | 1 | 35s | 15 | 5c |
| recipe_potion_hp_medium | 5 | 2× Papatya, 1× Yabani Nane | 1× Orta Boy Can İksiri | 5 | 45s | 30 | 15c |
| recipe_potion_hp_large | 10 | 2× Kızıl Adaçayı, 2× Papatya | 1× Büyük Boy Can İksiri | 12 | 60s | 50 | 40c |
| recipe_potion_hp_xlarge | 20 | 2× Gümüş Yaprak, 2× Kızıl Adaçayı | 1× Ekstra Büyük Can İksiri | 20 | 90s | 80 | 1s |
| recipe_potion_hp_huge | 30 | 1× Ruh Nilüferi, 2× Gümüş Yaprak | 1× Çok Büyük Can İksiri | 28 | 120s | 120 | 2s |

### Demircilik — 20 tarif, beceri 1–10 (tavan 60), çıktı Sv 1–10

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_helm_leather | 1 | 2× Orman Faresi Postu | 1× Deri Kask | 1 | 150s | 15 | 5c |
| recipe_chest_leather | 1 | 4× Orman Faresi Postu | 1× Deri Göğüslük | 1 | 225s | 25 | 8c |
| recipe_legs_leather | 1 | 3× Orman Faresi Postu | 1× Deri Bacaklık | 1 | 180s | 20 | 6c |
| recipe_boots_leather | 1 | 2× Orman Faresi Postu | 1× Deri Çizme | 1 | 150s | 15 | 5c |
| recipe_gloves_leather | 1 | 1× Orman Faresi Postu | 1× Deri Eldiven | 1 | 150s | 12 | 4c |
| recipe_helm_iron | 3 | 4× Demir Cevheri | 1× Demir Kask | 5 | 210s | 35 | 20c |
| recipe_chest_iron | 3 | 6× Demir Cevheri, 2× Yaban Tavşanı Postu | 1× Demir Göğüslük | 5 | 315s | 45 | 30c |
| recipe_legs_iron | 3 | 5× Demir Cevheri | 1× Demir Bacaklık | 5 | 250s | 40 | 25c |
| recipe_boots_iron | 3 | 3× Demir Cevheri | 1× Demir Çizme | 5 | 210s | 30 | 15c |
| recipe_gloves_iron | 3 | 2× Demir Cevheri | 1× Demir Eldiven | 5 | 210s | 25 | 12c |
| recipe_sword_iron | 5 | 5× Demir Cevheri, 1× Meşe Kütüğü | 1× Demir Kılıç | 5 | 405s | 40 | 25c |
| recipe_mace_iron | 5 | 5× Demir Cevheri, 2× Meşe Kütüğü | 1× Demir Topuz | 5 | 405s | 40 | 25c |
| recipe_dagger_iron | 5 | 3× Demir Cevheri, 1× Meşe Kütüğü | 1× Demir Hançer | 5 | 405s | 35 | 20c |
| recipe_bow_hunting | 5 | 4× Meşe Kütüğü, 2× Yaban Tavşanı Postu | 1× Av Yayı | 5 | 405s | 35 | 20c |
| recipe_staff_apprentice | 5 | 4× Meşe Kütüğü, 2× Papatya | 1× Çırak Asası | 5 | 405s | 35 | 20c |
| recipe_helm_scale | 10 | 4× Bakır Cevheri, 2× Kalay Cevheri | 1× Pul Kask | 10 | 420s | 60 | 80c |
| recipe_chest_scale | 10 | 6× Bakır Cevheri, 3× Kalay Cevheri, 2× Yaban Kurdu Postu | 1× Pul Göğüslük | 10 | 630s | 80 | 1s 40c |
| recipe_legs_scale | 10 | 5× Bakır Cevheri, 3× Kalay Cevheri | 1× Pul Bacaklık | 10 | 505s | 70 | 1s 10c |
| recipe_boots_scale | 10 | 3× Bakır Cevheri, 1× Yaban Kurdu Postu | 1× Pul Çizme | 10 | 420s | 55 | 70c |
| recipe_gloves_scale | 10 | 2× Bakır Cevheri, 1× Yaban Kurdu Postu | 1× Pul Eldiven | 10 | 420s | 50 | 60c |

### Aşçılık — 18 tarif, beceri 1–90 (tavan 100), çıktı Sv 1–100

| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |
|---|---:|---|---|---:|---:|---:|---:|
| recipe_cook_grilled_guppy | 1 | 2× Minik Lepistes, 1× Sıradan Ot, 1× Fidan Dalı | 2× Izgara Lepistes | 1 | 20s | 40 | 2c |
| recipe_cook_sardine_skewer | 4 | 2× Nehir Sardalyası, 1× Karahindiba, 1× Kavak Kütüğü | 2× Sardalya Şiş | 4 | 30s | 40 | 2c |
| recipe_cook_minnow_soup | 6 | 2× Golyan, 1× Yabani Nane, 1× Söğüt Kütüğü | 2× Naneli Balık Çorbası | 7 | 30s | 40 | 2c |
| recipe_cook_pike_fillet | 9 | 2× Küçük Turna, 1× Papatya, 1× Köknar Kütüğü | 2× Turna Fileto | 10 | 40s | 60 | 5c |
| recipe_cook_lavender_carp | 11 | 2× Sazan, 1× Lavanta, 1× Çam Kütüğü | 2× Lavantalı Sazan | 12 | 40s | 60 | 5c |
| recipe_cook_herb_trout | 14 | 2× Dere Alabalığı, 1× Acıkök, 1× Sedir Kütüğü | 2× Otlu Alabalık | 16 | 50s | 60 | 5c |
| recipe_cook_smoked_salmon | 16 | 2× Somon, 1× Şifalı Yosun, 1× Meşe Kütüğü | 2× Füme Somon | 18 | 50s | 60 | 5c |
| recipe_cook_bass_stew | 18 | 2× Tatlı Su Levreği, 1× Bataklık Sazı, 1× Kayın Kütüğü | 2× Levrek Güveci | 20 | 55s | 90 | 10c |
| recipe_cook_tuna_steak | 25 | 2× Mavi Yüzgeçli Ton, 1× Yaban Ginsengi, 1× Tik Kütüğü | 2× Ginsengli Ton Balığı | 28 | 70s | 90 | 10c |
| recipe_cook_royal_carp_feast | 30 | 2× Kraliyet Sazanı, 1× Düş Kökü, 1× Abanoz Kütüğü | 2× Kraliyet Sazanı Ziyafeti | 33 | 80s | 180 | 20c |
| recipe_cook_sturgeon_pie | 38 | 2× Altın Mersin, 1× Hayalet Mantarı, 1× Kara Odun Kütüğü | 2× Mantarlı Mersin Turtası | 42 | 95s | 180 | 20c |
| recipe_cook_swordfish_grill | 40 | 2× Kızıl Kılıçbalığı, 1× Fırtına Çiçeği, 1× Kristal Huş Kütüğü | 2× Fırtına Izgarası Kılıçbalığı | 45 | 100s | 260 | 40c |
| recipe_cook_frostfin_chowder | 50 | 2× Ayaz Yüzgeçli Somon, 1× Ayaz Çiçeği, 1× Ayaz Çamı Kütüğü | 2× Ayaz Somon Çorbası | 55 | 120s | 260 | 40c |
| recipe_cook_marlin_platter | 58 | 2× Hayalet Marlin, 1× Ruh Nilüferi, 1× Ruhağacı Kütüğü | 2× Ruh Marlini Tabağı | 65 | 135s | 460 | 80c |
| recipe_cook_dragonblood_serpentfish | 68 | 2× Kristal Yılanbalığı, 1× Ejderkanı Otu, 1× Yıldızkabuğu Kütüğü | 2× Ejderkanlı Yılanbalığı | 75 | 155s | 760 | 1s 50c |
| recipe_cook_moonfish_delicacy | 76 | 2× Yıldız Aybalığı, 1× Boşluk Orkidesi, 1× Ejderyüreği Kütüğü | 2× Yıldız Aybalığı Lezzeti | 85 | 170s | 760 | 1s 50c |
| recipe_cook_phoenix_koi_feast | 81 | 2× Anka Koisi, 1× Göksel Nilüfer, 1× Boşlukağacı Kütüğü | 2× Anka Koi Ziyafeti | 90 | 180s | 1100 | 2s 50c |
| recipe_cook_leviathan_banquet | 90 | 2× Göksel Leviathan, 1× Dünya Kökü, 1× Dünya Ağacı Kütüğü | 2× Leviathan Kraliyet Şöleni | 100 | 200s | 1100 | 2s 50c |

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

## 4. Ganimet (drop) kuralları

- **Ganimet seviyesi** = min(canavar Sv, karakter Sv) (`activity.DropLevel`).
- **Ekipman**: ganimet seviyesine eşit ya da altındaki **en yüksek tek kademe**, karakterin sınıfına uygun (`item.GearPool`; gereksinimli eşya düşmez).
  - Av: her zafer %4 → sıradan %60 · seçkin %25 · nadir %11 · destansı %4.
  - Zindan: dalga öldürmesi başına %4 + boss’ta 1 garanti → boss tablosu seçkin %50 · nadir %35 · destansı %13 · efsanevi %2.
  - Baskın: dalga %4 + boss’ta 2 garanti, kişisel ganimet.
- **Deri/post** (`item.CombatHides`): avcılık seviyesi = canavar Sv × 4 (tavan 100). Sıradan buluntu açılan son 4 hayvandan, nadir/efsanevi en iyisinden.
- **Toplama**: seçilen kaynağın kendisi, kademe rengiyle (`TierRarity`): T0–1 sıradan, T2–3 seçkin, T4 nadir, T5 destansı, T6 efsanevi, T7 mitik.
- **Altın / XP**: canavar tablosu; XP seviye farkı bantları (10 aşağı tam, 25 aşağı %50, ötesi %10).

## 5. Canavar listesi

Bölge, `scripts/gen_world.py` ZONES’undan; Sv 1–25 canavarları elle yazılmış (`mobs.yaml` başı), bölgesi yok. Av ve zindan başlatmak bulunulan bölgeye bağlı değil (yolculuk yalnız konum değiştiriyor).

| Sv | Kod | Ad | Arketip | Bölge | Can | Sal | Sav | XP | Altın | Ekipman kademesi | Deri (en iyi) | Rol |
|---:|---|---|---|---|---:|---:|---:|---:|---|---:|---|---|
| 1 | young_boar | Genç Yaban Domuzu | — | — | 172 | 29 | 2 | 10 | 7c–18c | 1 | Küçük Tavşan Postu | görev |
| 3 | grey_wolf | Boz Kurt | — | — | 212 | 35 | 4 | 64 | 24c–60c | 1 | Yaban Tavşanı Postu | dalga ×1, görev |
| 5 | bandit_scout | Haydut Gözcü | — | — | 300 | 41 | 6 | 131 | 48c–1s 20c | 5 | Kızıl Tilki Postu | dalga ×1, görev |
| 6 | cave_bat | Mağara Yarasası | — | — | 266 | 42 | 6 | 193 | 62c–1s 55c | 5 | Dağ Keçisi Postu | dalga ×1 |
| 8 | goblin_raider | Goblin Yağmacı | — | — | 396 | 50 | 9 | 279 | 94c–2s 35c | 5 | Yaban Mandası Derisi | boss Kül Mahzeni, dalga ×1, görev |
| 10 | marsh_lurker | Bataklık Pusucusu | — | — | 460 | 56 | 11 | 398 | 1s 29c–3s 22c | 10 | Dağ Aslanı Postu | dalga ×1 |
| 12 | forest_troll | Orman Trolü | — | — | 760 | 53 | 21 | 640 | 1s 69c–4s 22c | 10 | Fırtına Kartalı Tüyü | dalga ×1, görev |
| 14 | dire_wolf | Dev Kurt | — | — | 529 | 68 | 14 | 750 | 2s 13c–5s 32c | 10 | Ayaz Kurdu Postu | boss Orman Kovuğu, dalga ×1 |
| 15 | cultist_acolyte | Tarikat Müridi | — | — | 465 | 85 | 10 | 800 | 2s 35c–5s 88c | 15 | Dev Ayı Postu | dalga ×1 |
| 18 | rock_ogre | Kaya Devi | — | — | 1038 | 68 | 30 | 1224 | 3s 8c–7s 70c | 15 | Kadim Mamut Postu | dalga ×1 |
| 20 | shadow_stalker | Gölge Avcısı | — | — | 624 | 82 | 17 | 1328 | 3s 61c–9s 2c | 20 | Yaşlı Wyvern Pulu | boss Gölgebataklık Harabeleri, dalga ×1 |
| 22 | frost_giant | Buz Devi | — | — | 1224 | 78 | 37 | 1687 | 4s 16c–10s 40c | 20 | Yıldız Grifonu Tüyü | dalga ×1 |
| 24 | elder_wyrm | Kadim Ejder | — | — | 1317 | 83 | 40 | 1939 | 4s 73c–11s 82c | 20 | Kadim Anka Tüyü | dalga ×1 |
| 25 | ashen_lich | Kül Ölüsü | — | — | 705 | 121 | 16 | 1811 | 5s 3c–12s 58c | 25 | Göksel Ejder Pulu | boss Obsidyen Kule |
| 30 | mist_leech | Sis Sülüğü | evasive | Sisli Bataklıklar | 855 | 105 | 25 | 2309 | 6s 61c–16s 52c | 25 | Göksel Ejder Pulu | dalga ×1 |
| 35 | mud_golem | Çamur Golemi | tank | Sisli Bataklıklar | 1794 | 102 | 55 | 2955 | 8s 33c–20s 82c | 35 | Göksel Ejder Pulu | dalga ×1 |
| 40 | bog_witch | Bataklık Cadısı | glass | Sisli Bataklıklar | 1035 | 168 | 25 | 3658 | 10s 18c–25s 45c | 35 | Göksel Ejder Pulu | dalga ×1 |
| 42 | marsh_queen | Bataklık Kraliçesi | chief | Sisli Bataklıklar | 1440 | 125 | 38 | 3955 | 10s 95c–27s 38c | 35 | Göksel Ejder Pulu | boss Batık Tapınak |
| 45 | venom_lizard | Zehirli Kertenkele | crit | Sisli Bataklıklar | 1339 | 145 | 37 | 4417 | 12s 15c–30s 38c | 35 | Göksel Ejder Pulu |  |
| 50 | drowned_knight | Batık Şövalye | elite | Sisli Bataklıklar | 1960 | 153 | 49 | 5228 | 14s 23c–35s 58c | 50 | Göksel Ejder Pulu |  |
| 55 | mountain_troll | Dağ Trolü | tank | Demir Dağlar | 2669 | 144 | 86 | 6089 | 16s 41c–41s 2c | 50 | Göksel Ejder Pulu | dalga ×1 |
| 60 | dwarf_ghost | Cüce Hayaleti | evasive | Demir Dağlar | 1568 | 177 | 49 | 6999 | 18s 70c–46s 75c | 50 | Göksel Ejder Pulu | dalga ×1 |
| 65 | rock_eagle | Kaya Kartalı | glass | Demir Dağlar | 1598 | 247 | 40 | 7955 | 21s 9c–52s 72c | 65 | Göksel Ejder Pulu | dalga ×1 |
| 67 | iron_king | Demir Kral | chief | Demir Dağlar | 2190 | 182 | 60 | 8351 | 22s 7c–55s 18c | 65 | Göksel Ejder Pulu | boss Unutulmuş Maden |
| 70 | iron_spider | Demir Örümcek | crit | Demir Dağlar | 1995 | 206 | 57 | 8957 | 23s 57c–58s 92c | 65 | Göksel Ejder Pulu |  |
| 75 | mine_warden | Maden Bekçisi | elite | Demir Dağlar | 2835 | 213 | 73 | 10002 | 26s 14c–65s 35c | 65 | Göksel Ejder Pulu |  |
| 80 | ghost_guard | Hayalet Muhafız | normal | Kuzey Kalesi Harabeleri | 2150 | 196 | 65 | 11090 | 28s 79c–71s 98c | 80 | Göksel Ejder Pulu | dalga ×1 |
| 85 | cursed_archer | Lanetli Okçu | glass | Kuzey Kalesi Harabeleri | 2048 | 310 | 52 | 12220 | 31s 53c–78s 82c | 80 | Göksel Ejder Pulu | dalga ×1 |
| 90 | bone_rider | Kemik Atlı | crit | Kuzey Kalesi Harabeleri | 2520 | 256 | 73 | 13390 | 34s 36c–85s 90c | 80 | Göksel Ejder Pulu | dalga ×2 |
| 92 | fallen_lord_varn | Düşmüş Lord Varn | chief | Kuzey Kalesi Harabeleri | 2940 | 239 | 82 | 13869 | 35s 51c–88s 78c | 80 | Göksel Ejder Pulu | boss Kuzey Kalesi |
| 95 | keep_gargoyle | Kale Gargoyli | tank | Kuzey Kalesi Harabeleri | 4419 | 228 | 146 | 14600 | 37s 26c–93s 15c | 95 | Göksel Ejder Pulu | dalga ×1 |
| 100 | fallen_commander | Düşmüş Komutan | elite | Kuzey Kalesi Harabeleri | 3710 | 274 | 97 | 15849 | 40s 24c–1g 60c | 95 | Göksel Ejder Pulu | dalga ×1 |
| 102 | varn_undying | Ölümsüz Varn | titan | Kuzey Kalesi Kuşatması | 81000 | 775 | 116 | 16359 | 41s 45c–1g 3s 62c | 95 | Göksel Ejder Pulu | baskın boss’u Kuzey Kalesi Kuşatması |
| 110 | rebel_rider | Asi Süvari | normal | Taç Ovaları | 2900 | 259 | 89 | 18460 | 46s 42c–1g 16s 5c | 110 | Göksel Ejder Pulu | dalga ×1 |
| 120 | plains_lion | Ova Aslanı | crit | Taç Ovaları | 3308 | 330 | 97 | 21217 | 52s 90c–1g 32s 25c | 110 | Göksel Ejder Pulu | dalga ×1 |
| 130 | raider_mage | Yağmacı Büyücü | evasive | Taç Ovaları | 3230 | 346 | 105 | 24116 | 59s 64c–1g 49s 10c | 110 | Göksel Ejder Pulu | dalga ×1 |
| 132 | rebel_king_orrin | Asi Kral Orrin | elite | Taç Ovaları | 4830 | 351 | 128 | 24713 | 61s 3c–1g 52s 58c | 110 | Göksel Ejder Pulu | boss Asiler Kampı |
| 140 | war_bear | Savaş Ayısı | tank | Taç Ovaları | 6388 | 322 | 215 | 27152 | 66s 66c–1g 66s 65c | 135 | Göksel Ejder Pulu |  |
| 150 | black_lancer | Kara Mızrakçı | elite | Taç Ovaları | 5460 | 394 | 145 | 30321 | 73s 93c–1g 84s 82c | 135 | Göksel Ejder Pulu |  |
| 160 | ash_wolf | Kül Kurdu | evasive | Yanık Orman | 3942 | 419 | 129 | 33620 | 81s 44c–2g 3s 60c | 160 | Göksel Ejder Pulu | dalga ×1 |
| 170 | flame_spirit | Alev Ruhu | glass | Yanık Orman | 3960 | 578 | 103 | 37044 | 89s 19c–2g 22s 98c | 160 | Göksel Ejder Pulu | dalga ×1 |
| 180 | cinder_treant | Kömür Treanti | tank | Yanık Orman | 8138 | 406 | 276 | 40592 | 97s 18c–2g 42s 95c | 160 | Göksel Ejder Pulu | dalga ×1 |
| 182 | burnt_mother_tree | Yanık Ana Ağaç | elite | Yanık Orman | 6580 | 472 | 176 | 41316 | 98s 80c–2g 47s | 160 | Göksel Ejder Pulu | boss Ateş Kökü |
| 190 | crimson_harpy | Kızıl Harpya | crit | Yanık Orman | 5145 | 504 | 153 | 44259 | 1g 5s 39c–2g 63s 48c | 185 | Göksel Ejder Pulu |  |
| 200 | scorched_shaman | Yanık Şaman | elite | Yanık Orman | 7210 | 515 | 193 | 48045 | 1g 13s 82c–2g 84s 55c | 185 | Göksel Ejder Pulu |  |
| 210 | sand_scorpion | Kum Akrebi | crit | Tuz Çölü | 5670 | 553 | 169 | 51946 | 1g 22s 46c–3g 6s 15c | 210 | Göksel Ejder Pulu | dalga ×1 |
| 220 | desert_viper | Çöl Yılanı | evasive | Tuz Çölü | 5368 | 564 | 177 | 55960 | 1g 31s 31c–3g 28s 28c | 210 | Göksel Ejder Pulu | dalga ×1 |
| 230 | salt_elemental | Tuz Elementali | tank | Tuz Çölü | 10325 | 511 | 352 | 60085 | 1g 40s 36c–3g 50s 90c | 210 | Göksel Ejder Pulu | dalga ×1 |
| 232 | pharaoh_seth_ka | Firavun Seth-Ka | elite | Tuz Çölü | 8330 | 592 | 224 | 60923 | 1g 42s 20c–3g 55s 50c | 210 | Göksel Ejder Pulu | boss Gömülü Saray |
| 240 | mummy_priest | Mumya Rahip | glass | Tuz Çölü | 5535 | 798 | 145 | 64319 | 1g 49s 61c–3g 74s 2c | 235 | Göksel Ejder Pulu |  |
| 250 | sand_worm | Kum Solucanı | elite | Tuz Çölü | 8960 | 636 | 241 | 68660 | 1g 59s 6c–3g 97s 65c | 235 | Göksel Ejder Pulu |  |
| 260 | royal_skeleton | Kraliyet İskeleti | normal | Taç Mezarlığı | 6650 | 574 | 209 | 73107 | 1g 68s 70c–4g 21s 75c | 260 | Göksel Ejder Pulu | dalga ×1 |
| 270 | tomb_keeper | Mezar Bekçisi | tank | Taç Mezarlığı | 12075 | 595 | 412 | 77657 | 1g 78s 53c–4g 46s 32c | 260 | Göksel Ejder Pulu | dalga ×2 |
| 280 | soul_drinker | Ruh Emici | glass | Taç Mezarlığı | 6435 | 924 | 169 | 82310 | 1g 88s 54c–4g 71s 35c | 260 | Göksel Ejder Pulu | dalga ×1 |
| 282 | crypt_lord_malachar | Mezar Lordu Malachar | elite | Taç Mezarlığı | 10080 | 713 | 272 | 83253 | 1g 90s 56c–4g 76s 40c | 260 | Göksel Ejder Pulu | boss Kraliyet Kriptası |
| 290 | black_knight | Kara Şövalye | crit | Taç Mezarlığı | 7770 | 752 | 233 | 87064 | 1g 98s 73c–4g 96s 82c | 285 | Göksel Ejder Pulu | dalga ×1 |
| 300 | crownless_king | Taçsız Kral | elite | Taç Mezarlığı | 10710 | 757 | 289 | 91917 | 2g 9s 9c–5g 22s 72c | 285 | Göksel Ejder Pulu | dalga ×1 |
| 302 | crownless_tyrant | Taçsız Tiran | titan | Taçsız Kral'ın Sarayı | 231000 | 2119 | 340 | 92899 | 2g 11s 19c–5g 27s 98c | 285 | Göksel Ejder Pulu | baskın boss’u Taçsız Kral'ın Sarayı |
| 310 | sea_raider | Deniz Haydudu | normal | Fırtına Kıyısı | 7900 | 679 | 249 | 96868 | 2g 19s 63c–5g 49s 8c | 310 | Göksel Ejder Pulu | dalga ×1 |
| 320 | storm_bird | Fırtına Kuşu | evasive | Fırtına Kıyısı | 7742 | 805 | 257 | 101915 | 2g 30s 35c–5g 75s 88c | 310 | Göksel Ejder Pulu | dalga ×1 |
| 330 | coral_golem | Mercan Golemi | tank | Fırtına Kıyısı | 14700 | 721 | 504 | 107059 | 2g 41s 23c–6g 3s 8c | 310 | Göksel Ejder Pulu | dalga ×1 |
| 332 | admiral_drask | Amiral Drask | overlord | Fırtına Kıyısı | 15632 | 943 | 333 | 108099 | 2g 43s 42c–6g 8s 55c | 310 | Göksel Ejder Pulu | boss Batık Filo |
| 340 | siren | Siren | glass | Fırtına Kıyısı | 7785 | 1113 | 205 | 112297 | 2g 52s 28c–6g 30s 70c | 335 | Göksel Ejder Pulu |  |
| 350 | kraken_spawn | Kraken Yavrusu | elite | Fırtına Kıyısı | 12460 | 877 | 337 | 117628 | 2g 63s 49c–6g 58s 72c | 335 | Göksel Ejder Pulu |  |
| 360 | young_dragon | Genç Ejder | crit | Ejder Yamaçları | 9608 | 925 | 289 | 123051 | 2g 74s 86c–6g 87s 15c | 360 | Göksel Ejder Pulu | dalga ×1 |
| 370 | dragon_cultist | Ejder Kültisti | glass | Ejder Yamaçları | 8460 | 1208 | 223 | 128565 | 2g 86s 39c–7g 15s 98c | 360 | Göksel Ejder Pulu | dalga ×1 |
| 380 | lava_lizard | Lav Kertenkelesi | tank | Ejder Yamaçları | 16888 | 826 | 580 | 134169 | 2g 98s 8c–7g 45s 20c | 360 | Göksel Ejder Pulu | dalga ×1 |
| 382 | mother_dragon_vyrna | Ana Ejder Vyrna | overlord | Ejder Yamaçları | 17945 | 1079 | 383 | 135301 | 3g 44c–7g 51s 10c | 360 | Göksel Ejder Pulu | boss Kuluçka Mağarası |
| 390 | winged_wyvern | Kanatlı Wyvern | evasive | Ejder Yamaçları | 9405 | 974 | 313 | 139863 | 3g 9s 92c–7g 74s 80c | 385 | Göksel Ejder Pulu |  |
| 400 | dragon_guard | Ejder Muhafızı | elite | Ejder Yamaçları | 14210 | 998 | 385 | 145645 | 3g 21s 92c–8g 4s 80c | 385 | Göksel Ejder Pulu |  |
| 410 | fey_wolf | Peri Kurdu | evasive | Sis Ormanı | 9880 | 1022 | 329 | 151514 | 3g 34s 7c–8g 35s 18c | 410 | Göksel Ejder Pulu | dalga ×1 |
| 420 | mist_stalker | Sis Avcısı | crit | Sis Ormanı | 11182 | 1074 | 337 | 157470 | 3g 46s 36c–8g 65s 90c | 410 | Göksel Ejder Pulu | dalga ×1 |
| 430 | ancient_treant | Kadim Treant | tank | Sis Ormanı | 19075 | 931 | 656 | 163512 | 3g 58s 81c–8g 97s 2c | 410 | Göksel Ejder Pulu | dalga ×1 |
| 432 | spirit_queen_sylwen | Ruh Kraliçesi Sylwen | overlord | Sis Ormanı | 20258 | 1216 | 433 | 164730 | 3g 61s 31c–9g 3s 28c | 410 | Göksel Ejder Pulu | boss Kayıp Koru |
| 440 | shadow_elf | Gölge Elfi | glass | Sis Ormanı | 10035 | 1428 | 265 | 169638 | 3g 71s 40c–9g 28s 50c | 435 | Göksel Ejder Pulu |  |
| 450 | forest_spirit | Orman Ruhu | elite | Sis Ormanı | 15960 | 1119 | 433 | 175849 | 3g 84s 13c–9g 60s 32c | 435 | Göksel Ejder Pulu |  |
| 460 | crystal_spider | Kristal Örümcek | crit | Kristal Mağaralar | 12232 | 1173 | 369 | 182143 | 3g 97s–9g 92s 50c | 460 | Göksel Ejder Pulu | dalga ×1 |
| 470 | stone_eater | Taş Yiyen | tank | Kristal Mağaralar | 20825 | 1015 | 716 | 188520 | 4g 10s 2c–10g 25s 5c | 460 | Göksel Ejder Pulu | dalga ×1 |
| 480 | light_elemental | Işık Elementali | glass | Kristal Mağaralar | 10935 | 1554 | 289 | 194978 | 4g 23s 17c–10g 57s 92c | 460 | Göksel Ejder Pulu | dalga ×1 |
| 482 | crystal_titan | Kristal Titan | overlord | Kristal Mağaralar | 22570 | 1352 | 483 | 196280 | 4g 25s 82c–10g 64s 55c | 460 | Göksel Ejder Pulu | boss Kristal Kalp |
| 490 | crystal_dragon | Kristal Ejder | evasive | Kristal Mağaralar | 11780 | 1216 | 393 | 201518 | 4g 36s 47c–10g 91s 18c | 485 | Göksel Ejder Pulu |  |
| 500 | cave_titan | Mağara Titanı | elite | Kristal Mağaralar | 17710 | 1240 | 481 | 208138 | 4g 49s 90c–11g 24s 75c | 485 | Göksel Ejder Pulu |  |
| 510 | magma_golem | Magma Golemi | tank | Volkanik Ada | 22575 | 1099 | 777 | 214839 | 4g 63s 46c–11g 58s 65c | 510 | Göksel Ejder Pulu | dalga ×1 |
| 520 | fire_giant | Ateş Devi | normal | Volkanik Ada | 13150 | 1120 | 417 | 221618 | 4g 77s 16c–11g 92s 90c | 510 | Göksel Ejder Pulu | dalga ×1 |
| 530 | ash_dragon | Kül Ejderi | crit | Volkanik Ada | 14070 | 1346 | 425 | 228476 | 4g 90s 99c–12g 27s 48c | 510 | Göksel Ejder Pulu | dalga ×1 |
| 532 | fire_lord_ignar | Ateş Lordu Ignar | overlord | Volkanik Ada | 24882 | 1489 | 533 | 229857 | 4g 93s 77c–12g 34s 42c | 510 | Göksel Ejder Pulu | boss Volkan Kalbi |
| 540 | volcano_shaman | Volkan Şamanı | glass | Volkanik Ada | 12285 | 1743 | 325 | 235413 | 5g 4s 95c–12g 62s 38c | 535 | Göksel Ejder Pulu |  |
| 550 | volcano_warden | Volkan Bekçisi | elite | Volkanik Ada | 19460 | 1360 | 529 | 242427 | 5g 19s 4c–12g 97s 60c | 535 | Göksel Ejder Pulu |  |
| 560 | dragon_knight | Ejder Şövalyesi | normal | Taht Zirveleri | 14150 | 1204 | 449 | 249517 | 5g 33s 26c–13g 33s 15c | 560 | Göksel Ejder Pulu | dalga ×1 |
| 570 | ancient_wyvern | Kadim Wyvern | evasive | Taht Zirveleri | 13680 | 1409 | 457 | 256685 | 5g 47s 61c–13g 69s 2c | 560 | Göksel Ejder Pulu | dalga ×1 |
| 580 | dragon_sorcerer | Ejder Büyücüsü | glass | Taht Zirveleri | 13185 | 1869 | 349 | 263928 | 5g 62s 8c–14g 5s 20c | 560 | Göksel Ejder Pulu | dalga ×2 |
| 582 | dragon_king_azhorath | Ejder Kral Azhorath | overlord | Taht Zirveleri | 27195 | 1625 | 583 | 265385 | 5g 64s 99c–14g 12s 48c | 560 | Göksel Ejder Pulu | boss Ejder Tahtı |
| 590 | throne_guardian | Taht Muhafızı | tank | Taht Zirveleri | 26075 | 1267 | 899 | 271246 | 5g 76s 68c–14g 41s 70c | 585 | Göksel Ejder Pulu | dalga ×1 |
| 600 | elder_dragon | Yaşlı Ejder | elite | Taht Zirveleri | 21210 | 1481 | 577 | 278639 | 5g 91s 40c–14g 78s 50c | 585 | Göksel Ejder Pulu | dalga ×1 |
| 602 | azhorath_ascended | Yükselmiş Azhorath | elder_titan | Ejder Kral'ın İni | 456000 | 6203 | 676 | 280127 | 5g 94s 36c–14g 85s 90c | 585 | Göksel Ejder Pulu | baskın boss’u Ejder Kral'ın İni |
| 610 | shadow_rider | Gölge Atlı | normal | Karanlık Ova | 15400 | 1309 | 489 | 286107 | 6g 6s 25c–15g 15s 62c | 610 | Göksel Ejder Pulu | dalga ×1 |
| 620 | void_wolf | Boşluk Kurdu | evasive | Karanlık Ova | 14868 | 1529 | 497 | 293648 | 6g 21s 22c–15g 53s 5c | 610 | Göksel Ejder Pulu | dalga ×1 |
| 630 | star_gnawer | Yıldız Kemiren | glass | Karanlık Ova | 14310 | 2026 | 379 | 301262 | 6g 36s 31c–15g 90s 78c | 610 | Göksel Ejder Pulu | dalga ×1 |
| 632 | star_eater_omh | Yıldız Yiyen Omh | overlord | Karanlık Ova | 29508 | 1762 | 633 | 302794 | 6g 39s 34c–15g 98s 35c | 610 | Göksel Ejder Pulu | boss Söndürülmüş Gözlemevi |
| 640 | dark_angel | Kara Melek | crit | Karanlık Ova | 16958 | 1619 | 513 | 308950 | 6g 51s 52c–16g 28s 80c | 635 | Göksel Ejder Pulu |  |
| 650 | night_giant | Gece Devi | elite | Karanlık Ova | 22960 | 1602 | 625 | 316710 | 6g 66s 85c–16g 67s 12c | 635 | Göksel Ejder Pulu |  |
| 660 | ghost_noble | Hayalet Soylu | normal | Kayıp Şehir | 16650 | 1414 | 529 | 324542 | 6g 82s 30c–17g 5s 75c | 660 | Göksel Ejder Pulu | dalga ×1 |
| 670 | time_thief | Zaman Hırsızı | evasive | Kayıp Şehir | 16055 | 1650 | 537 | 332445 | 6g 97s 86c–17g 44s 65c | 660 | Göksel Ejder Pulu | dalga ×1 |
| 680 | cursed_sorcerer | Lanetli Büyücü | glass | Kayıp Şehir | 15435 | 2184 | 409 | 340419 | 7g 13s 54c–17g 83s 85c | 660 | Göksel Ejder Pulu | dalga ×2 |
| 682 | time_lord_kaedrin | Zaman Lordu Kaedrin | overlord | Kayıp Şehir | 31820 | 1898 | 683 | 342023 | 7g 16s 70c–17g 91s 75c | 660 | Göksel Ejder Pulu | boss Saat Kulesi |
| 690 | stone_gargoyle | Taş Gargoyl | tank | Kayıp Şehir | 30450 | 1477 | 1051 | 348465 | 7g 29s 34c–18g 23s 35c | 685 | Göksel Ejder Pulu | dalga ×1 |
| 700 | city_guardian | Şehir Koruyucusu | elite | Kayıp Şehir | 24710 | 1723 | 673 | 356580 | 7g 45s 26c–18g 63s 15c | 685 | Göksel Ejder Pulu | dalga ×1 |
| 702 | kaedrin_timeless | Zamansız Kaedrin | elder_titan | Zaman Lordu'nun Kulesi | 531000 | 7211 | 788 | 358211 | 7g 48s 45c–18g 71s 12c | 685 | Göksel Ejder Pulu | baskın boss’u Zaman Lordu'nun Kulesi |
| 710 | cloud_giant | Bulut Devi | tank | Uçan Adalar | 31325 | 1519 | 1081 | 364765 | 7g 61s 28c–19g 3s 20c | 710 | Göksel Ejder Pulu | dalga ×1 |
| 720 | storm_dragon | Fırtına Ejderi | crit | Uçan Adalar | 19058 | 1817 | 577 | 373020 | 7g 77s 42c–19g 43s 55c | 710 | Göksel Ejder Pulu | dalga ×1 |
| 730 | sky_knight | Gök Şövalyesi | normal | Uçan Adalar | 18400 | 1561 | 585 | 381344 | 7g 93s 67c–19g 84s 18c | 710 | Göksel Ejder Pulu | dalga ×1 |
| 732 | storm_king_zephyros | Fırtına Kralı Zephyros | overlord | Uçan Adalar | 34132 | 2035 | 733 | 383017 | 7g 96s 94c–19g 92s 35c | 710 | Göksel Ejder Pulu | boss Fırtına Kalesi |
| 740 | wind_elemental | Rüzgâr Elementali | evasive | Uçan Adalar | 17718 | 1819 | 593 | 389736 | 8g 10s 4c–20g 25s 10c | 735 | Göksel Ejder Pulu |  |
| 750 | thunderbird | Yıldırım Kuşu | elite | Uçan Adalar | 26460 | 1843 | 721 | 398197 | 8g 26s 51c–20g 66s 28c | 735 | Göksel Ejder Pulu |  |
| 760 | void_creature | Boşluk Yaratığı | normal | Yarık Kenarı | 19150 | 1624 | 609 | 406726 | 8g 43s 10c–21g 7s 75c | 760 | Göksel Ejder Pulu | dalga ×1 |
| 770 | rift_warden | Yarık Bekçisi | tank | Yarık Kenarı | 33950 | 1645 | 1172 | 415322 | 8g 59s 79c–21g 49s 48c | 760 | Göksel Ejder Pulu | dalga ×1 |
| 780 | chaos_elemental | Kaos Elementali | glass | Yarık Kenarı | 17685 | 2499 | 469 | 423986 | 8g 76s 60c–21g 91s 50c | 760 | Göksel Ejder Pulu | dalga ×1 |
| 782 | chaos_mother | Kaos Ana | overlord | Yarık Kenarı | 36445 | 2171 | 783 | 425727 | 8g 79s 97c–21g 99s 92c | 760 | Göksel Ejder Pulu | boss Yarık Çekirdeği |
| 790 | rift_dragon | Yarık Ejderi | crit | Yarık Kenarı | 20895 | 1991 | 633 | 432716 | 8g 93s 51c–22g 33s 78c | 785 | Göksel Ejder Pulu |  |
| 800 | chaos_lord | Kaos Lordu | elite | Yarık Kenarı | 28210 | 1964 | 769 | 441513 | 9g 10s 53c–22g 76s 32c | 785 | Göksel Ejder Pulu |  |
| 810 | fallen_angel | Düşmüş Melek | normal | Göksel Tapınak | 20400 | 1729 | 649 | 450377 | 9g 27s 65c–23g 19s 12c | 810 | Göksel Ejder Pulu | dalga ×1 |
| 820 | temple_guardian | Tapınak Muhafızı | tank | Göksel Tapınak | 36138 | 1750 | 1248 | 459306 | 9g 44s 88c–23g 62s 20c | 810 | Göksel Ejder Pulu | dalga ×1 |
| 830 | celestial_dragon | Göksel Ejder | crit | Göksel Tapınak | 21945 | 2090 | 665 | 468301 | 9g 62s 22c–24g 5s 55c | 810 | Göksel Ejder Pulu | dalga ×2 |
| 832 | archangel_uriel | Başmelek Uriel | overlord | Göksel Tapınak | 38758 | 2308 | 833 | 470108 | 9g 65s 70c–24g 14s 25c | 810 | Göksel Ejder Pulu | boss Göksel Kubbe |
| 840 | light_hunter | Işık Avcısı | evasive | Göksel Tapınak | 20092 | 2061 | 673 | 477361 | 9g 79s 66c–24g 49s 15c | 835 | Göksel Ejder Pulu | dalga ×1 |
| 850 | seraph | Seraf | elite | Göksel Tapınak | 29960 | 2085 | 817 | 486486 | 9g 97s 21c–24g 93s 2c | 835 | Göksel Ejder Pulu | dalga ×1 |
| 852 | fallen_archangels | Düşmüş Başmelekler | elder_titan | Başmeleklerin Düşüşü | 643500 | 8723 | 956 | 488319 | 10g 73c–25g 1s 82c | 835 | Göksel Ejder Pulu | baskın boss’u Başmeleklerin Düşüşü |
| 860 | ash_leviathan | Kül Leviathanı | tank | Kül Denizi | 37888 | 1834 | 1309 | 495676 | 10g 14s 86c–25g 37s 15c | 860 | Göksel Ejder Pulu | dalga ×1 |
| 870 | null_hunter | Hiçlik Avcısı | evasive | Kül Denizi | 20805 | 2133 | 697 | 504930 | 10g 32s 61c–25g 81s 52c | 860 | Göksel Ejder Pulu | dalga ×1 |
| 880 | undying_soldier | Ölümsüz Asker | normal | Kül Denizi | 22150 | 1876 | 705 | 514248 | 10g 50s 47c–26g 26s 18c | 860 | Göksel Ejder Pulu | dalga ×1 |
| 882 | leviathan_king | Leviathan Kral | overlord | Kül Denizi | 41070 | 2444 | 883 | 516119 | 10g 54s 5c–26g 35s 12c | 860 | Göksel Ejder Pulu | boss Batık Taht |
| 890 | ash_wyrm_lord | Kül Ejderi Lordu | crit | Kül Denizi | 23520 | 2238 | 713 | 523629 | 10g 68s 42c–26g 71s 5c | 885 | Göksel Ejder Pulu |  |
| 900 | eternal_sentinel | Ebedi Bekçi | elite | Kül Denizi | 31710 | 2206 | 865 | 533075 | 10g 86s 48c–27g 16s 20c | 885 | Göksel Ejder Pulu |  |
| 910 | doom_rider | Kıyamet Atlısı | crit | Son Kale Surları | 24045 | 2288 | 729 | 542583 | 11g 4s 64c–27g 61s 60c | 910 | Göksel Ejder Pulu | dalga ×1 |
| 920 | void_titan | Boşluk Titanı | tank | Son Kale Surları | 40512 | 1960 | 1400 | 552154 | 11g 22s 90c–28g 7s 25c | 910 | Göksel Ejder Pulu | dalga ×1 |
| 930 | soul_reaper | Ruh Biçen | glass | Son Kale Surları | 21060 | 2972 | 559 | 561788 | 11g 41s 25c–28g 53s 12c | 910 | Göksel Ejder Pulu | dalga ×1 |
| 932 | doom_commander_morvath | Kıyamet Komutanı Morvath | overlord | Son Kale Surları | 43382 | 2581 | 933 | 563723 | 11g 44s 94c–28g 62s 35c | 910 | Göksel Ejder Pulu | boss Son Kale |
| 940 | ancient_golem | Kadim Golem | normal | Son Kale Surları | 23650 | 2002 | 753 | 571485 | 11g 59s 71c–28g 99s 28c | 935 | Göksel Ejder Pulu |  |
| 950 | last_commander | Son Komutan | elite | Son Kale Surları | 33460 | 2326 | 913 | 581243 | 11g 78s 27c–29g 45s 68c | 935 | Göksel Ejder Pulu |  |
| 960 | void_dragon | Hiçlik Ejderi | crit | Ebediyet Kıyısı | 25358 | 2412 | 769 | 591063 | 11g 96s 92c–29g 92s 30c | 960 | Göksel Ejder Pulu | dalga ×1 |
| 970 | throne_sentinel | Taht Gözcüsü | tank | Ebediyet Kıyısı | 42700 | 2065 | 1476 | 600945 | 12g 15s 67c–30g 39s 18c | 960 | Göksel Ejder Pulu | dalga ×2 |
| 980 | eternal_mage | Ebedi Büyücü | glass | Ebediyet Kıyısı | 22185 | 3129 | 589 | 610888 | 12g 34s 52c–30g 86s 30c | 960 | Göksel Ejder Pulu | dalga ×1 |
| 982 | fallen_god_aethon | Düşmüş Tanrı Aethon | overlord | Ebediyet Kıyısı | 45695 | 2717 | 983 | 612884 | 12g 38s 30c–30g 95s 75c | 960 | Göksel Ejder Pulu | boss Ebedi Taht |
| 990 | god_shadow | Tanrı Gölgesi | evasive | Ebediyet Kıyısı | 23655 | 2423 | 793 | 620892 | 12g 53s 46c–31g 33s 65c | 985 | Göksel Ejder Pulu | dalga ×1 |
| 1000 | aethon_unbound | Zincirsiz Aethon | elder_titan | Düşmüş Tanrı'nın Tahtı | 754500 | 10214 | 1121 | 630957 | 12g 72s 50c–31g 81s 25c | 985 | Göksel Ejder Pulu | baskın boss’u Düşmüş Tanrı'nın Tahtı |
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

İçerikte set yok. Tasarım `docs/item-engine/set.md` (parça sayımı, kademeli bonus, imza tetikleyici, `content/item_sets/*.yaml`) hazır, motor ve içerik yok. Aynı malzeme adını taşıyan 5 zırh parçası (ör. demir başlık/göğüs/bacak/çizme/eldiven) set adayı.

## 8. Uyumsuzluklar ve eksikler

| # | Önem | Alan | Bulgu |
|---:|---|---|---|
| 1 | Yüksek | Ekipman | 12 ekipman slotundan 5 tanesi boş: necklace, ring1, ring2, trinket1, trinket2. Karakter ekranında hep boş görünüyor, güç bütçesinin ~%30’u kullanılmıyor. |
| 2 | Yüksek | Ekonomi | 66/150 toplama kaynağı hiçbir tarifte kullanılmıyor (deri/post/tüy/pul 27, taş 4, balık 12, ot 10, kütük 12). Avcılık mesleği ürününün tek değeri dükkân fiyatı. |
| 3 | Yüksek | Ganimet | Deri ganimeti Sv 25 canavardan itibaren tavanda (avcılık 100): Sv 25–1000 arası tüm canavarlar aynı 4 hayvanın derisini düşürüyor, en iyisi Göksel Ejder Pulu. |
| 4 | Yüksek | Nadirlik | Ekler (sharp, of_the_bear …) yalnız isim: `item/affix.go` “cosmetic only”; stat vermiyor. Nadirlik ekipmanın gücünü değiştirmiyor, yalnız satış değerini — destansı kılıç ile sıradan kılıç aynı vuruyor. |
| 5 | Yüksek | Set | Set sistemi yok (tasarım var). Zindan/baskın ödülü ve uzun vadeli hedef için en büyük boşluk. |
| 6 | Yüksek | Üretim | 24/24 külçe hiçbir tarifte kullanılmıyor (adamantite_bar, orichalcum_bar, mythril_bar, starsteel_bar, iridium_bar, osmium_bar …). Demircilik tarifleri külçe değil ham cevher istiyor; ergitme (24 tarif, beceri 100’e kadar) çıkmaz sokak. |
| 7 | Yüksek | Üretim | Simya: son tarif beceri 30, meslek tavanı 60 — son 30 beceri seviyesinde açılacak tarif yok. |
| 8 | Yüksek | Üretim | Demircilik: son tarif beceri 10, meslek tavanı 60 — son 50 beceri seviyesinde açılacak tarif yok. |
| 9 | Yüksek | Üretim | En yüksek üretilebilir ekipman Sv 10; ekipman merdiveni Sv 985. Sv 10 üstü ekipmanın tek kaynağı %4 av / boss ganimeti. |
| 10 | Orta | Canavar | Av bölgeye bağlı değil: bölgeler yalnız yolculuk ekranında var; oyuncu herhangi bir yerden herhangi bir canavarı avlayabiliyor (DD-254/255 dünyaları anlamını yitiriyor). |
| 11 | Orta | Ekipman | Yan el (offhand) yalnız 1 şablon (Titan Siper Kalkanı, Sv 600, koruyucu uzmanlığı). Kalkan izni olan savaşçı/paladin/rahip 600’e kadar kalkan bulamıyor. |
| 12 | Orta | Ekipman | 21/32 silah ailesinin hiç eşyası yok (axe, greatsword, greataxe, spear, short_sword, dual_dagger, throwing_knife, longbow …). Sınıf izinleri (DD-252) her sınıfa 4–6 aile tanıyor, içerikte sınıf başı tek aile var. |
| 13 | Orta | Ekipman | Zırh türleri (kumaş/hafif/orta/ağır) tanımlı ama hiçbir zırhın alt sınıfı bu türlerden değil; `CanUse` aile dışı alt sınıfa izin verdiği için büyücü ile savaşçı aynı zırhı giyiyor. |
| 14 | Orta | Eşya motoru | Tasarımı olup uygulanmayan: soket + mücevher (`gems.md`, `sockets.md`), büyü (`enchantments.md`), üretim kalitesi / kritik üretim (`quality.md`, `critical-craft.md`), benzersiz etki (`unique-effect.md`). |
| 15 | Orta | Ganimet | Canavara özgü ganimet tablosu yok: aynı seviyedeki ejderha ile haydut aynı şeyi düşürüyor (yalnız seviye belirliyor). Görev/zanaat için özel malzeme (ejder pulu, boss mührü) yok. |
| 16 | Orta | Ganimet | Ganimet tek kademeden: Sv 34 canavar yalnız Sv 25 kademesi düşürür; bir alt kademe/üst kademe şansı yok. Kademe sınırını geçince düşen eşya bir anda değişiyor. |
| 17 | Orta | Nadirlik | Mitik ekipman hiçbir kaynaktan gelmiyor (boss tablosu efsanevide bitiyor); `unique-effect.md` tasarımı uygulanmadı. |
| 18 | Orta | Tüketilebilir | En büyük iksir +700, en büyük yemek +925 can; Sv 1000 canavar canı ~25150. İksir Sv 28, yemek Sv 100’de bitiyor — üst seviyede iyileştirme anlamsız. |
| 19 | Orta | Zindan | Zindana özgü ganimet yok: boss ekipmanı açık dünyadaki aynı kademe; zindan/baskın seti, mührü ya da benzersiz eşyası yok. |
| 20 | Orta | Üretim | Üretilebilir silahı olmayan sınıflar: bard, druid, monk, paladin, shaman (yalnız cleric, mage, ranger, rogue, warrior silahı demircilikte). |
| 21 | Orta | Üretim | Kuyumculuk / terzilik / deri işçiliği / oymacılık yok: deri, post, tüy, kütük ve taş için üretim yolu yok; takı slotları boş. |
| 22 | Düşük | Canavar | Her bölgede 6 canavar → Sv 100 sonrası her ~10 seviyede bir canavar; yan yana seçenek az (aynı seviyede ikinci canavar yok). |
| 23 | Düşük | Eşya | 6 eşya bilinçli olarak emekli (DD-241/244; eski yığınlar satılabilsin diye duruyor), kaynağı yok: potion_mp_small, linseed, mandrake, king_bloom, elm_log, yew_log. Temizlik isteğe bağlı. |
