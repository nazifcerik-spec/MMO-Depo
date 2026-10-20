# Savaşçı av istatistiği — Sv 1–50

`WARRIOR_REPORT=1 go test ./internal/balance -run TestWarriorHuntReport`; üretim savaş kodu (`hunt.Roll`), 100 tohum, saatlik AFK av (120 sn'de bir karşılaşma, saatte en çok 30), karakter seviyesinin ekipmanı (her yuvada ulaşılan en yüksek kademe), seviyeye göre dağıtılmış nitelik puanları, sınıf pasifleri ve aktifleri. Canavar: Sv−10 … Sv+5 arası her av canavarı.

Sütunlar: **Tur/öldürme** = kazanılan savaşın ortalama turu (oyuncu her tur bir kez vurur); **Can kaybı** = kazanılan savaşta azami canın ortalama yüzdesi; **Ölüm** = bir saatlik avın ölümle bitme oranı; **Öldürme/sa**, **XP/sa**; **Dk/sv** = bu canavarla bir sonraki seviyeye dakika.

## Özet

Güvenli = bir saatlik avın ölümle bitme oranı ≤ %10; ölümcül = ≥ %50. Güvenli her canavar aynı saatlik XP'yi verir; saldırı gücü savaşın süresini ve can kaybını belirler.

| Sv | Saldırı | En güçlü güvenli canavar | Tur/öldürme | Can kaybı | İlk ölümcül canavar |
|---:|---:|---|---:|---:|---|
| 1 | 52 | grey_wolf (Sv 3) | 5.2 | %21 | bandit_scout (Sv 5) |
| 2 | 59 | grey_wolf (Sv 3) | 4.7 | %16 | cave_bat (Sv 6) |
| 3 | 63 | bandit_scout (Sv 5) | 5.8 | %24 | cave_bat (Sv 6) |
| 4 | 69 | cave_bat (Sv 6) | 5.7 | %23 | goblin_raider (Sv 8) |
| 5 | 85 | cave_bat (Sv 6) | 4.6 | %14 | marsh_lurker (Sv 10) |
| 6 | 89 | goblin_raider (Sv 8) | 5.8 | %23 | marsh_lurker (Sv 10) |
| 7 | 93 | goblin_raider (Sv 8) | 5.3 | %18 | marsh_lurker (Sv 10) |
| 8 | 99 | goblin_raider (Sv 8) | 4.4 | %12 | forest_troll (Sv 12) |
| 9 | 103 | marsh_lurker (Sv 10) | 5.5 | %21 | forest_troll (Sv 12) |
| 10 | 127 | forest_troll (Sv 12) | 7.2 | %20 | dire_wolf (Sv 14) |
| 11 | 131 | forest_troll (Sv 12) | 6.4 | %15 | dire_wolf (Sv 14) |
| 12 | 135 | dire_wolf (Sv 14) | 5.4 | %22 | cultist_acolyte (Sv 15) |
| 13 | 139 | dire_wolf (Sv 14) | 4.7 | %17 | cultist_acolyte (Sv 15) |
| 14 | 145 | dire_wolf (Sv 14) | 4.4 | %14 | rock_ogre (Sv 18) |
| 15 | 167 | cultist_acolyte (Sv 15) | 3.1 | %11 | rock_ogre (Sv 18) |
| 16 | 173 | rock_ogre (Sv 18) | 7.7 | %21 | shadow_stalker (Sv 20) |
| 17 | 177 | rock_ogre (Sv 18) | 6.5 | %16 | shadow_stalker (Sv 20) |
| 18 | 181 | rock_ogre (Sv 18) | 6.1 | %12 | frost_giant (Sv 22) |
| 19 | 185 | shadow_stalker (Sv 20) | 5.5 | %18 | frost_giant (Sv 22) |
| 20 | 206 | frost_giant (Sv 22) | 7.0 | %18 | elder_wyrm (Sv 24) |
| 21 | 210 | frost_giant (Sv 22) | 6.4 | %14 | elder_wyrm (Sv 24) |
| 22 | 216 | elder_wyrm (Sv 24) | 8.2 | %21 | ashen_lich (Sv 25) |
| 23 | 220 | elder_wyrm (Sv 24) | 6.7 | %16 | ashen_lich (Sv 25) |
| 24 | 224 | ashen_lich (Sv 25) | 4.0 | %18 | — |
| 25 | 255 | ashen_lich (Sv 25) | 2.9 | %9 | mist_leech (Sv 30) |
| 26 | 261 | ashen_lich (Sv 25) | 2.9 | %8 | — |
| 27 | 265 | mist_leech (Sv 30) | 4.7 | %14 | — |
| 28 | 271 | mist_leech (Sv 30) | 4.2 | %11 | — |
| 29 | 275 | mist_leech (Sv 30) | 3.9 | %9 | — |
| 30 | 279 | mist_leech (Sv 30) | 3.7 | %7 | mud_golem (Sv 35) |
| 31 | 283 | mist_leech (Sv 30) | 3.6 | %6 | mud_golem (Sv 35) |
| 32 | 289 | mud_golem (Sv 35) | 6.5 | %16 | — |
| 33 | 293 | mud_golem (Sv 35) | 6.1 | %13 | — |
| 34 | 300 | mud_golem (Sv 35) | 5.8 | %11 | — |
| 35 | 335 | mud_golem (Sv 35) | 5.0 | %4 | bog_witch (Sv 40) |
| 36 | 339 | mud_golem (Sv 35) | 5.0 | %4 | bog_witch (Sv 40) |
| 37 | 343 | marsh_queen (Sv 42) | 5.5 | %17 | — |
| 38 | 349 | marsh_queen (Sv 42) | 4.9 | %13 | — |
| 39 | 353 | marsh_queen (Sv 42) | 4.5 | %11 | — |
| 40 | 359 | marsh_queen (Sv 42) | 4.2 | %9 | venom_lizard (Sv 45) |
| 41 | 363 | marsh_queen (Sv 42) | 4.0 | %7 | — |
| 42 | 367 | venom_lizard (Sv 45) | 4.2 | %13 | — |
| 43 | 371 | venom_lizard (Sv 45) | 3.8 | %11 | — |
| 44 | 377 | venom_lizard (Sv 45) | 3.6 | %9 | — |
| 45 | 381 | venom_lizard (Sv 45) | 3.4 | %7 | drowned_knight (Sv 50) |
| 46 | 387 | venom_lizard (Sv 45) | 3.4 | %7 | drowned_knight (Sv 50) |
| 47 | 391 | venom_lizard (Sv 45) | 3.4 | %7 | drowned_knight (Sv 50) |
| 48 | 395 | drowned_knight (Sv 50) | 5.4 | %15 | — |
| 49 | 399 | drowned_knight (Sv 50) | 4.8 | %12 | — |
| 50 | 474 | mountain_troll (Sv 55) | 6.8 | %16 | — |

## Saldırı gücü taraması

Aynı karakter, yalnız saldırı değeri çarpılarak; canavar seviyeye uygun en güçlü av canavarı.

| Sv | Canavar | Çarpan | Saldırı | Tur/öldürme | Can kaybı | Ölüm | XP/sa | Dk/sv |
|---:|---|---:|---:|---:|---:|---:|---:|---:|
| 10 | marsh_lurker | ×0.6 | 76 | 6.8 | %18 | %0 | 60000 | 4 |
| 10 | marsh_lurker | ×0.8 | 102 | 5.2 | %13 | %0 | 60000 | 4 |
| 10 | marsh_lurker | ×1.0 | 127 | 4.2 | %10 | %0 | 60000 | 4 |
| 10 | marsh_lurker | ×1.2 | 152 | 3.3 | %7 | %0 | 60000 | 4 |
| 10 | marsh_lurker | ×1.4 | 178 | 3.1 | %7 | %0 | 60000 | 4 |
| 10 | marsh_lurker | ×1.6 | 203 | 3.0 | %6 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×0.6 | 153 | 4.1 | %14 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×0.8 | 204 | 3.2 | %10 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×1.0 | 255 | 2.9 | %9 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×1.2 | 306 | 2.9 | %8 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×1.4 | 357 | 2.1 | %5 | %0 | 60000 | 4 |
| 25 | ashen_lich | ×1.6 | 408 | 2.0 | %5 | %0 | 60000 | 4 |
| 50 | drowned_knight | ×0.6 | 284 | 6.2 | %10 | %0 | 60000 | 5 |
| 50 | drowned_knight | ×0.8 | 379 | 4.9 | %8 | %0 | 60000 | 5 |
| 50 | drowned_knight | ×1.0 | 474 | 4.0 | %6 | %0 | 60000 | 5 |
| 50 | drowned_knight | ×1.2 | 569 | 3.2 | %4 | %0 | 60000 | 5 |
| 50 | drowned_knight | ×1.4 | 664 | 3.0 | %4 | %0 | 60000 | 5 |
| 50 | drowned_knight | ×1.6 | 758 | 2.9 | %4 | %0 | 60000 | 5 |

## Sv 1 — saldırı 52 · savunma 31 · can 414 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 4.0 | %8 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 5.2 | %21 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 7.5 | %39 | %100 | 8.8 | 17620 | 13 |
| cave_bat | 6 | 266 | 42 | 6 | 8.7 | %51 | %100 | 2.4 | 4820 | 47 |

## Sv 2 — saldırı 59 · savunma 33 · can 430 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 3.2 | %5 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 4.7 | %16 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 6.5 | %30 | %12 | 28.1 | 56100 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 7.6 | %40 | %100 | 5.9 | 11780 | 19 |

## Sv 3 — saldırı 63 · savunma 35 · can 446 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 3.2 | %5 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 4.4 | %13 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 5.8 | %24 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 6.9 | %33 | %85 | 14.9 | 29800 | 8 |
| goblin_raider | 8 | 396 | 50 | 9 | 8.8 | %56 | %100 | 2.0 | 3960 | 57 |

## Sv 4 — saldırı 69 · savunma 37 · can 462 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 3.2 | %4 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 3.6 | %9 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 5.4 | %20 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 5.7 | %23 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 7.4 | %42 | %100 | 4.4 | 8720 | 26 |

## Sv 5 — saldırı 85 · savunma 45 · can 494 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.2 | %1 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 3.4 | %6 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 4.3 | %10 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 4.6 | %14 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 6.1 | %28 | %37 | 25.2 | 50320 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 7.7 | %47 | %100 | 2.9 | 5780 | 39 |

## Sv 6 — saldırı 89 · savunma 48 · can 528 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %1 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 3.4 | %5 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 4.2 | %9 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 3.9 | %9 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 5.8 | %23 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 6.9 | %37 | %100 | 6.0 | 12000 | 19 |

## Sv 7 — saldırı 93 · savunma 49 · can 544 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %1 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 3.3 | %4 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 3.9 | %7 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 3.8 | %8 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 5.3 | %18 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 6.1 | %29 | %84 | 17.2 | 34340 | 7 |
| forest_troll | 12 | 760 | 53 | 21 | 11.3 | %58 | %100 | 1.7 | 3400 | 67 |

## Sv 8 — saldırı 99 · savunma 51 · can 560 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %1 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 2.8 | %3 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 3.3 | %5 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 3.8 | %7 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 4.4 | %12 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 5.8 | %25 | %13 | 28.3 | 56580 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 9.8 | %46 | %100 | 3.0 | 6080 | 37 |

## Sv 9 — saldırı 103 · savunma 53 · can 576 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %0 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 2.5 | %2 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 3.2 | %5 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 3.8 | %7 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 4.3 | %11 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 5.5 | %21 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 9.2 | %39 | %100 | 4.7 | 9320 | 24 |
| dire_wolf | 14 | 529 | 68 | 14 | 7.5 | %52 | %100 | 1.9 | 3700 | 62 |

## Sv 10 — saldırı 127 · savunma 63 · can 666 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %0 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 2.2 | %1 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 3.1 | %3 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.9 | %3 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.5 | %6 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 4.2 | %10 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 7.2 | %20 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 6.9 | %35 | %100 | 5.5 | 11020 | 21 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 5.9 | %40 | %100 | 3.8 | 7520 | 30 |

## Sv 11 — saldırı 131 · savunma 65 · can 682 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 2.1 | %0 | %0 | 30.0 | 60000 | 4 |
| grey_wolf | 3 | 212 | 35 | 4 | 2.2 | %0 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.7 | %2 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.2 | %5 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 4.2 | %9 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 6.4 | %15 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 6.1 | %27 | %98 | 13.6 | 27240 | 8 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 5.6 | %36 | %100 | 5.1 | 10140 | 23 |

## Sv 12 — saldırı 135 · savunma 68 · can 716 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| grey_wolf | 3 | 212 | 35 | 4 | 2.2 | %0 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 2.9 | %2 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.6 | %1 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.2 | %4 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 4.2 | %8 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 6.1 | %12 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 5.4 | %22 | %7 | 29.4 | 58740 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 5.5 | %31 | %100 | 8.3 | 16520 | 14 |

## Sv 13 — saldırı 139 · savunma 70 · can 732 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| grey_wolf | 3 | 212 | 35 | 4 | 2.2 | %0 | %0 | 30.0 | 60000 | 4 |
| bandit_scout | 5 | 300 | 41 | 6 | 2.7 | %1 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.5 | %1 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.1 | %4 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 4.1 | %8 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 6.1 | %11 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 4.7 | %17 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 5.3 | %27 | %97 | 15.3 | 30600 | 8 |
| rock_ogre | 18 | 1038 | 68 | 30 | 11.5 | %50 | %100 | 2.2 | 4360 | 53 |

## Sv 14 — saldırı 145 · savunma 72 · can 748 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bandit_scout | 5 | 300 | 41 | 6 | 2.3 | %1 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.5 | %1 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.1 | %3 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 3.7 | %6 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 5.9 | %10 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 4.4 | %14 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 4.4 | %20 | %12 | 29.0 | 58000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 11.1 | %44 | %100 | 3.0 | 5980 | 39 |

## Sv 15 — saldırı 167 · savunma 82 · can 820 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bandit_scout | 5 | 300 | 41 | 6 | 2.1 | %0 | %0 | 30.0 | 60000 | 4 |
| cave_bat | 6 | 266 | 42 | 6 | 2.5 | %0 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 3.1 | %3 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 5.1 | %5 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.9 | %9 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 3.1 | %11 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 9.3 | %27 | %100 | 14.2 | 28300 | 8 |
| shadow_stalker | 20 | 624 | 82 | 17 | 7.3 | %36 | %100 | 4.4 | 8840 | 26 |

## Sv 16 — saldırı 173 · savunma 84 · can 836 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| cave_bat | 6 | 266 | 42 | 6 | 2.5 | %0 | %0 | 30.0 | 60000 | 4 |
| goblin_raider | 8 | 396 | 50 | 9 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 3.1 | %3 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.5 | %8 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 3.1 | %10 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 7.7 | %21 | %2 | 30.0 | 59960 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 6.9 | %30 | %100 | 6.6 | 13140 | 18 |

## Sv 17 — saldırı 177 · savunma 86 · can 852 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| goblin_raider | 8 | 396 | 50 | 9 | 3.0 | %1 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 3.1 | %3 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.4 | %7 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 3.1 | %10 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 6.5 | %16 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 6.6 | %26 | %98 | 12.6 | 25260 | 9 |
| frost_giant | 22 | 1224 | 78 | 37 | 10.9 | %45 | %100 | 2.6 | 5160 | 45 |

## Sv 18 — saldırı 181 · savunma 89 · can 886 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| goblin_raider | 8 | 396 | 50 | 9 | 2.8 | %1 | %0 | 30.0 | 60000 | 4 |
| marsh_lurker | 10 | 460 | 56 | 11 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 4.8 | %3 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.3 | %6 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 3.1 | %9 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 6.1 | %12 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 6.2 | %21 | %25 | 27.4 | 54780 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 10.1 | %37 | %100 | 4.4 | 8780 | 27 |

## Sv 19 — saldırı 185 · savunma 90 · can 902 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| marsh_lurker | 10 | 460 | 56 | 11 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 4.5 | %3 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.3 | %6 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 3.1 | %9 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 6.1 | %11 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 5.5 | %18 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 9.5 | %30 | %100 | 7.2 | 14420 | 16 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 11.4 | %49 | %100 | 2.2 | 4340 | 54 |

## Sv 20 — saldırı 206 · savunma 102 · can 992 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| marsh_lurker | 10 | 460 | 56 | 11 | 2.9 | %1 | %0 | 30.0 | 60000 | 4 |
| forest_troll | 12 | 760 | 53 | 21 | 4.0 | %0 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.2 | %4 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.9 | %6 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 5.4 | %7 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.8 | %8 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 7.0 | %18 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 9.6 | %30 | %100 | 7.4 | 14840 | 16 |
| ashen_lich | 25 | 705 | 121 | 16 | 5.8 | %37 | %100 | 4.0 | 7960 | 30 |

## Sv 21 — saldırı 210 · savunma 104 · can 1008 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| forest_troll | 12 | 760 | 53 | 21 | 4.0 | %0 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.2 | %3 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.8 | %6 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 5.2 | %6 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.7 | %7 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 6.4 | %14 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 9.2 | %25 | %100 | 14.1 | 28100 | 8 |
| ashen_lich | 25 | 705 | 121 | 16 | 5.5 | %32 | %100 | 5.6 | 11240 | 21 |

## Sv 22 — saldırı 216 · savunma 105 · can 1024 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| forest_troll | 12 | 760 | 53 | 21 | 4.0 | %0 | %0 | 30.0 | 60000 | 4 |
| dire_wolf | 14 | 529 | 68 | 14 | 3.2 | %3 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.6 | %5 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 5.0 | %5 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.7 | %7 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 6.0 | %11 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 8.2 | %21 | %7 | 29.6 | 59300 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 5.2 | %28 | %100 | 8.7 | 17400 | 14 |

## Sv 23 — saldırı 220 · savunma 107 · can 1040 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| dire_wolf | 14 | 529 | 68 | 14 | 3.2 | %3 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.5 | %4 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 5.0 | %5 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.7 | %7 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 6.0 | %10 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 6.7 | %16 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 4.8 | %23 | %92 | 16.1 | 32140 | 7 |

## Sv 24 — saldırı 224 · savunma 110 · can 1074 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| dire_wolf | 14 | 529 | 68 | 14 | 3.2 | %3 | %0 | 30.0 | 60000 | 4 |
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.3 | %3 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.7 | %6 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 6.0 | %9 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 6.3 | %12 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 4.0 | %18 | %7 | 29.4 | 58880 | 4 |

## Sv 25 — saldırı 255 · savunma 128 · can 1204 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| cultist_acolyte | 15 | 465 | 85 | 10 | 2.0 | %2 | %0 | 30.0 | 60000 | 4 |
| rock_ogre | 18 | 1038 | 68 | 30 | 4.0 | %1 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.3 | %3 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.5 | %3 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 5.0 | %5 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %9 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 6.1 | %22 | %85 | 18.6 | 37140 | 7 |

## Sv 26 — saldırı 261 · savunma 129 · can 1220 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| rock_ogre | 18 | 1038 | 68 | 30 | 3.9 | %1 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.3 | %3 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.3 | %3 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.9 | %5 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %8 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 5.4 | %18 | %20 | 27.9 | 55820 | 4 |

## Sv 27 — saldırı 265 · savunma 131 · can 1236 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| rock_ogre | 18 | 1038 | 68 | 30 | 3.7 | %0 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.3 | %3 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.2 | %3 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.8 | %5 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %8 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 4.7 | %14 | %0 | 30.0 | 60000 | 4 |

## Sv 28 — saldırı 271 · savunma 133 · can 1252 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| rock_ogre | 18 | 1038 | 68 | 30 | 3.5 | %0 | %0 | 30.0 | 60000 | 4 |
| shadow_stalker | 20 | 624 | 82 | 17 | 3.2 | %2 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.6 | %4 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %8 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 4.2 | %11 | %0 | 30.0 | 60000 | 4 |

## Sv 29 — saldırı 275 · savunma 135 · can 1268 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 3.2 | %2 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.5 | %4 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %8 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.9 | %9 | %0 | 30.0 | 60000 | 4 |

## Sv 30 — saldırı 279 · savunma 138 · can 1302 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 3.1 | %2 | %0 | 30.0 | 60000 | 4 |
| frost_giant | 22 | 1224 | 78 | 37 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.3 | %3 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %7 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.7 | %7 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 8.3 | %24 | %100 | 11.2 | 22300 | 11 |

## Sv 31 — saldırı 283 · savunma 140 · can 1318 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| frost_giant | 22 | 1224 | 78 | 37 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.2 | %3 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %7 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.6 | %6 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 7.7 | %20 | %57 | 25.0 | 50000 | 5 |

## Sv 32 — saldırı 289 · savunma 141 · can 1334 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| frost_giant | 22 | 1224 | 78 | 37 | 4.0 | %1 | %0 | 30.0 | 60000 | 4 |
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %7 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.5 | %6 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 6.5 | %16 | %0 | 30.0 | 60000 | 4 |

## Sv 33 — saldırı 293 · savunma 143 · can 1350 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %7 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.5 | %6 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 6.1 | %13 | %0 | 30.0 | 60000 | 4 |

## Sv 34 — saldırı 300 · savunma 145 · can 1366 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| elder_wyrm | 24 | 1317 | 83 | 40 | 4.1 | %2 | %0 | 30.0 | 60000 | 4 |
| ashen_lich | 25 | 705 | 121 | 16 | 2.9 | %6 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.4 | %5 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 5.8 | %11 | %0 | 30.0 | 60000 | 4 |

## Sv 35 — saldırı 335 · savunma 167 · can 1524 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| ashen_lich | 25 | 705 | 121 | 16 | 2.4 | %3 | %0 | 30.0 | 60000 | 4 |
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %3 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 4.5 | %23 | %98 | 11.8 | 23500 | 11 |

## Sv 36 — saldırı 339 · savunma 170 · can 1558 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %3 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 4.2 | %19 | %67 | 21.6 | 43180 | 6 |

## Sv 37 — saldırı 343 · savunma 172 · can 1574 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %2 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 5.0 | %4 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.8 | %16 | %15 | 28.9 | 57700 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 5.5 | %17 | %6 | 29.7 | 59360 | 4 |

## Sv 38 — saldırı 349 · savunma 173 · can 1590 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %2 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 4.9 | %3 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.3 | %13 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 4.9 | %13 | %0 | 30.0 | 60000 | 4 |

## Sv 39 — saldırı 353 · savunma 175 · can 1606 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %2 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 4.9 | %3 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.1 | %11 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 4.5 | %11 | %0 | 30.0 | 60000 | 4 |

## Sv 40 — saldırı 359 · savunma 177 · can 1622 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mist_leech | 30 | 855 | 105 | 25 | 3.3 | %2 | %0 | 30.0 | 60000 | 4 |
| mud_golem | 35 | 1794 | 102 | 55 | 4.9 | %3 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.0 | %9 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 4.2 | %9 | %0 | 30.0 | 60000 | 4 |
| venom_lizard | 45 | 1339 | 145 | 37 | 5.4 | %21 | %96 | 14.4 | 28760 | 9 |

## Sv 41 — saldırı 363 · savunma 179 · can 1638 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mud_golem | 35 | 1794 | 102 | 55 | 4.8 | %3 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.0 | %9 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 4.0 | %7 | %0 | 30.0 | 60000 | 4 |
| venom_lizard | 45 | 1339 | 145 | 37 | 5.0 | %17 | %46 | 25.6 | 51280 | 5 |

## Sv 42 — saldırı 367 · savunma 182 · can 1672 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mud_golem | 35 | 1794 | 102 | 55 | 4.7 | %2 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.0 | %9 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.7 | %5 | %0 | 30.0 | 60000 | 4 |
| venom_lizard | 45 | 1339 | 145 | 37 | 4.2 | %13 | %0 | 30.0 | 60000 | 4 |

## Sv 43 — saldırı 371 · savunma 184 · can 1688 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mud_golem | 35 | 1794 | 102 | 55 | 4.6 | %2 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 3.0 | %8 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.6 | %5 | %0 | 30.0 | 60000 | 4 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.8 | %11 | %0 | 30.0 | 60000 | 4 |

## Sv 44 — saldırı 377 · savunma 185 · can 1704 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mud_golem | 35 | 1794 | 102 | 55 | 4.4 | %2 | %0 | 30.0 | 60000 | 4 |
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %8 | %0 | 30.0 | 60000 | 4 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.5 | %4 | %0 | 30.0 | 60000 | 4 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.6 | %9 | %0 | 30.0 | 60000 | 4 |

## Sv 45 — saldırı 381 · savunma 187 · can 1720 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| mud_golem | 35 | 1794 | 102 | 55 | 4.3 | %1 | %0 | 30.0 | 60000 | 5 |
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %8 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.4 | %4 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.4 | %7 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 7.0 | %27 | %100 | 6.4 | 12780 | 21 |

## Sv 46 — saldırı 387 · savunma 189 · can 1736 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %8 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.3 | %4 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.4 | %7 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 6.7 | %23 | %100 | 10.7 | 21320 | 13 |

## Sv 47 — saldırı 391 · savunma 191 · can 1752 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %8 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.2 | %3 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.4 | %7 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 6.4 | %19 | %84 | 19.9 | 39840 | 7 |

## Sv 48 — saldırı 395 · savunma 194 · can 1786 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %7 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.2 | %3 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.4 | %7 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 5.4 | %15 | %0 | 30.0 | 60000 | 5 |

## Sv 49 — saldırı 399 · savunma 196 · can 1802 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 2.9 | %7 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.2 | %3 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.4 | %6 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 4.8 | %12 | %0 | 30.0 | 60000 | 5 |

## Sv 50 — saldırı 474 · savunma 222 · can 2004 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 2.6 | %4 | %0 | 30.0 | 60000 | 5 |
| marsh_queen | 42 | 1440 | 125 | 38 | 3.0 | %1 | %0 | 30.0 | 60000 | 5 |
| venom_lizard | 45 | 1339 | 145 | 37 | 3.1 | %3 | %0 | 30.0 | 60000 | 5 |
| drowned_knight | 50 | 1960 | 153 | 49 | 4.0 | %6 | %0 | 30.0 | 60000 | 5 |
| mountain_troll | 55 | 2669 | 144 | 86 | 6.8 | %16 | %6 | 29.5 | 58980 | 5 |

