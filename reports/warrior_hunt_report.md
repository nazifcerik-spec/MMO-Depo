# Savaşçı av istatistiği — Sv 1–50

`WARRIOR_REPORT=1 go test ./internal/balance -run TestWarriorHuntReport`; üretim av kodu (`hunt.Roll`, DD-281), 60 tohum, bir saatlik av: savaşlar peş peşe, tur 3 sn, can %50'nin altına düşünce kamp (saniyede azami canın %1'i). Karakter: seviyenin ekipmanı (her yuvada ulaşılan en yüksek kademe, sıradan), dağıtılmış nitelik puanları, sınıf pasifleri ve aktifleri. Canavar: Sv−10 … Sv+5 arası her av canavarı; XP canavarın kendi değeri, seviye farkına göre düşer.

Sütunlar: **XP** = canavarın öldürme başı XP'si (seviye farkı uygulanmış); **Tur** = kazanılan savaşın ortalama turu; **Sn/öldürme** = savaş + kamp dahil öldürme başı av süresi; **Kamp** = sürenin kampta geçen payı; **Ölüm** = bir saatlik avın ölümle bitme oranı; **Öldürme/sa**, **XP/sa**; **Dk/sv** = bu canavarla bir sonraki seviyeye dakika.

## Özet

En verimli güvenli canavar = ölümle biten saat ≤ %10 olanlar içinde saatte en çok XP veren. Ölümcül = ≥ %50.

| Sv | Saldırı | En verimli güvenli canavar | Sn/öldürme | XP/sa | Dk/sv | İlk ölümcül canavar |
|---:|---:|---|---:|---:|---:|---|
| 1 | 52 | young_boar (Sv 1) | 20 | 39008 | 6 | — |
| 2 | 61 | young_boar (Sv 1) | 14 | 55208 | 4 | — |
| 3 | 69 | young_boar (Sv 1) | 12 | 63801 | 4 | — |
| 4 | 78 | young_boar (Sv 1) | 11 | 68070 | 4 | — |
| 5 | 98 | young_boar (Sv 1) | 6 | 117490 | 2 | — |
| 6 | 107 | young_boar (Sv 1) | 6 | 117490 | 2 | — |
| 7 | 115 | young_boar (Sv 1) | 6 | 117490 | 3 | — |
| 8 | 123 | young_boar (Sv 1) | 6 | 117490 | 3 | — |
| 9 | 132 | young_boar (Sv 1) | 6 | 117490 | 3 | — |
| 10 | 158 | young_boar (Sv 1) | 6 | 122960 | 3 | — |
| 11 | 166 | young_boar (Sv 1) | 5 | 139538 | 2 | — |
| 12 | 175 |  (Sv 0) | 0 | 0 | 0 | — |
| 13 | 183 |  (Sv 0) | 0 | 0 | 0 | — |
| 14 | 192 |  (Sv 0) | 0 | 0 | 0 | — |
| 15 | 217 | shadow_stalker (Sv 20) | 43 | 23904 | 17 | — |
| 16 | 226 | shadow_stalker (Sv 20) | 33 | 37418 | 11 | — |
| 17 | 234 | shadow_stalker (Sv 20) | 27 | 44854 | 9 | — |
| 18 | 243 | shadow_stalker (Sv 20) | 24 | 55299 | 8 | — |
| 19 | 251 | shadow_stalker (Sv 20) | 20 | 63813 | 7 | — |
| 20 | 274 | shadow_stalker (Sv 20) | 16 | 83710 | 6 | — |
| 21 | 282 | shadow_stalker (Sv 20) | 15 | 86357 | 6 | — |
| 22 | 291 | shadow_stalker (Sv 20) | 15 | 88872 | 6 | — |
| 23 | 299 | shadow_stalker (Sv 20) | 14 | 92476 | 5 | — |
| 24 | 308 | shadow_stalker (Sv 20) | 13 | 97801 | 5 | — |
| 25 | 343 | shadow_stalker (Sv 20) | 7 | 173178 | 3 | — |
| 26 | 351 | shadow_stalker (Sv 20) | 7 | 177702 | 3 | — |
| 27 | 360 | shadow_stalker (Sv 20) | 7 | 179074 | 3 | — |
| 28 | 368 | shadow_stalker (Sv 20) | 7 | 179152 | 3 | — |
| 29 | 376 | shadow_stalker (Sv 20) | 7 | 179206 | 3 | — |
| 30 | 385 | shadow_stalker (Sv 20) | 7 | 179255 | 3 | — |
| 31 | 393 |  (Sv 0) | 0 | 0 | 0 | — |
| 32 | 402 |  (Sv 0) | 0 | 0 | 0 | — |
| 33 | 410 |  (Sv 0) | 0 | 0 | 0 | — |
| 34 | 418 |  (Sv 0) | 0 | 0 | 0 | — |
| 35 | 458 | bog_witch (Sv 40) | 25 | 72289 | 10 | — |
| 36 | 466 | bog_witch (Sv 40) | 23 | 81328 | 9 | — |
| 37 | 475 | bog_witch (Sv 40) | 21 | 90827 | 8 | — |
| 38 | 483 | bog_witch (Sv 40) | 19 | 100403 | 7 | — |
| 39 | 491 | bog_witch (Sv 40) | 17 | 112173 | 7 | — |
| 40 | 500 | bog_witch (Sv 40) | 15 | 126655 | 6 | — |
| 41 | 507 | bog_witch (Sv 40) | 14 | 132037 | 6 | — |
| 42 | 515 | bog_witch (Sv 40) | 13 | 138944 | 6 | — |
| 43 | 522 | bog_witch (Sv 40) | 13 | 147125 | 6 | — |
| 44 | 530 | bog_witch (Sv 40) | 12 | 156087 | 5 | — |
| 45 | 537 | bog_witch (Sv 40) | 11 | 167657 | 5 | — |
| 46 | 545 | bog_witch (Sv 40) | 10 | 179105 | 5 | — |
| 47 | 552 | bog_witch (Sv 40) | 10 | 188067 | 5 | — |
| 48 | 560 | bog_witch (Sv 40) | 9 | 198059 | 5 | — |
| 49 | 567 | bog_witch (Sv 40) | 9 | 206171 | 5 | — |
| 50 | 650 | bog_witch (Sv 40) | 8 | 231105 | 4 | — |

## Saldırı gücü taraması

Aynı karakter, yalnız saldırı değeri çarpılarak; canavar seviyeye uygun en güçlü av canavarı.

| Sv | Canavar | Çarpan | Saldırı | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 10 | young_boar | ×0.6 | 95 | 2.2 | 6 | %0 | %0 | 556 | 117953 | 3 |
| 10 | young_boar | ×0.8 | 126 | 2.1 | 6 | %0 | %0 | 566 | 120077 | 3 |
| 10 | young_boar | ×1.0 | 158 | 2.1 | 6 | %0 | %0 | 580 | 122960 | 3 |
| 10 | young_boar | ×1.2 | 190 | 1.2 | 4 | %0 | %0 | 1011 | 214261 | 2 |
| 10 | young_boar | ×1.4 | 221 | 1.1 | 3 | %0 | %0 | 1103 | 233921 | 1 |
| 10 | young_boar | ×1.6 | 253 | 1.1 | 3 | %0 | %0 | 1103 | 233921 | 1 |
| 25 | shadow_stalker | ×0.6 | 206 | 3.9 | 12 | %0 | %0 | 305 | 110273 | 5 |
| 25 | shadow_stalker | ×0.8 | 274 | 3.3 | 10 | %0 | %0 | 365 | 131723 | 4 |
| 25 | shadow_stalker | ×1.0 | 343 | 2.5 | 7 | %0 | %0 | 480 | 173178 | 3 |
| 25 | shadow_stalker | ×1.2 | 412 | 2.4 | 7 | %0 | %0 | 497 | 179405 | 3 |
| 25 | shadow_stalker | ×1.4 | 480 | 2.4 | 7 | %0 | %0 | 497 | 179405 | 3 |
| 25 | shadow_stalker | ×1.6 | 549 | 2.4 | 7 | %0 | %0 | 497 | 179405 | 3 |
| 50 | bog_witch | ×0.6 | 390 | 3.9 | 12 | %1 | %0 | 301 | 156745 | 6 |
| 50 | bog_witch | ×0.8 | 520 | 3.5 | 11 | %2 | %0 | 333 | 173195 | 6 |
| 50 | bog_witch | ×1.0 | 650 | 2.6 | 8 | %4 | %0 | 444 | 231105 | 4 |
| 50 | bog_witch | ×1.2 | 780 | 2.6 | 8 | %4 | %0 | 444 | 231105 | 4 |
| 50 | bog_witch | ×1.4 | 910 | 2.6 | 8 | %4 | %0 | 444 | 231105 | 4 |
| 50 | bog_witch | ×1.6 | 1040 | 1.8 | 5 | %0 | %0 | 683 | 355030 | 3 |

## Sv 1 — saldırı 52 · savunma 31 · can 414 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 4.3 | 20 | %34 | %0 | 184 | 39008 | 6 |

## Sv 2 — saldırı 61 · savunma 34 · can 448 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 3.4 | 14 | %26 | %0 | 260 | 55208 | 4 |

## Sv 3 — saldırı 69 · savunma 37 · can 482 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 3.2 | 12 | %19 | %0 | 301 | 63801 | 4 |

## Sv 4 — saldırı 78 · savunma 40 · can 516 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 3.2 | 11 | %14 | %0 | 321 | 68070 | 4 |

## Sv 5 — saldırı 98 · savunma 49 · can 566 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.2 | 6 | %0 | %0 | 554 | 117490 | 2 |

## Sv 6 — saldırı 107 · savunma 52 · can 600 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.2 | 6 | %0 | %0 | 554 | 117490 | 2 |

## Sv 7 — saldırı 115 · savunma 55 · can 634 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.2 | 6 | %0 | %0 | 554 | 117490 | 3 |

## Sv 8 — saldırı 123 · savunma 58 · can 668 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.2 | 6 | %0 | %0 | 554 | 117490 | 3 |

## Sv 9 — saldırı 132 · savunma 61 · can 702 · kritik %5 · isabet %95

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.2 | 6 | %0 | %0 | 554 | 117490 | 3 |

## Sv 10 — saldırı 158 · savunma 73 · can 810 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 2.1 | 6 | %0 | %0 | 580 | 122960 | 3 |

## Sv 11 — saldırı 166 · savunma 76 · can 844 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| young_boar | 1 | 172 | 29 | 2 | 212 | 1.8 | 5 | %0 | %0 | 658 | 139538 | 2 |

## Sv 12 — saldırı 175 · savunma 79 · can 878 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 13 — saldırı 183 · savunma 82 · can 912 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 14 — saldırı 192 · savunma 85 · can 946 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 15 — saldırı 217 · savunma 96 · can 1036 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 6.5 | 43 | %54 | %3 | 66 | 23904 | 17 |

## Sv 16 — saldırı 226 · savunma 99 · can 1070 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 5.5 | 33 | %49 | %2 | 104 | 37418 | 11 |

## Sv 17 — saldırı 234 · savunma 102 · can 1104 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 5.0 | 27 | %45 | %0 | 124 | 44854 | 9 |

## Sv 18 — saldırı 243 · savunma 105 · can 1138 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.8 | 24 | %39 | %0 | 153 | 55299 | 8 |

## Sv 19 — saldırı 251 · savunma 108 · can 1172 · kritik %5 · isabet %97

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.5 | 20 | %34 | %0 | 177 | 63813 | 7 |

## Sv 20 — saldırı 274 · savunma 121 · can 1280 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.3 | 16 | %17 | %0 | 232 | 83710 | 6 |

## Sv 21 — saldırı 282 · savunma 124 · can 1314 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.3 | 15 | %14 | %0 | 239 | 86357 | 6 |

## Sv 22 — saldırı 291 · savunma 127 · can 1348 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.4 | 15 | %10 | %0 | 246 | 88872 | 6 |

## Sv 23 — saldırı 299 · savunma 130 · can 1382 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.4 | 14 | %5 | %0 | 256 | 92476 | 5 |

## Sv 24 — saldırı 308 · savunma 133 · can 1416 · kritik %6 · isabet %98

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 4.4 | 13 | %1 | %0 | 271 | 97801 | 5 |

## Sv 25 — saldırı 343 · savunma 152 · can 1564 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.5 | 7 | %0 | %0 | 480 | 173178 | 3 |

## Sv 26 — saldırı 351 · savunma 155 · can 1598 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.4 | 7 | %0 | %0 | 492 | 177702 | 3 |

## Sv 27 — saldırı 360 · savunma 158 · can 1632 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.4 | 7 | %0 | %0 | 496 | 179074 | 3 |

## Sv 28 — saldırı 368 · savunma 161 · can 1666 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.4 | 7 | %0 | %0 | 496 | 179152 | 3 |

## Sv 29 — saldırı 376 · savunma 164 · can 1700 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.4 | 7 | %0 | %0 | 496 | 179206 | 3 |

## Sv 30 — saldırı 385 · savunma 167 · can 1734 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| shadow_stalker | 20 | 624 | 82 | 17 | 361 | 2.4 | 7 | %0 | %0 | 497 | 179255 | 3 |

## Sv 31 — saldırı 393 · savunma 170 · can 1768 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 32 — saldırı 402 · savunma 173 · can 1802 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 33 — saldırı 410 · savunma 176 · can 1836 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 34 — saldırı 418 · savunma 179 · can 1870 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|

## Sv 35 — saldırı 458 · savunma 202 · can 2046 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 4.0 | 25 | %52 | %0 | 139 | 72289 | 10 |

## Sv 36 — saldırı 466 · savunma 205 · can 2080 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.9 | 23 | %49 | %0 | 156 | 81328 | 9 |

## Sv 37 — saldırı 475 · savunma 208 · can 2114 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.8 | 21 | %45 | %0 | 175 | 90827 | 8 |

## Sv 38 — saldırı 483 · savunma 211 · can 2148 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.6 | 19 | %42 | %0 | 193 | 100403 | 7 |

## Sv 39 — saldırı 491 · savunma 214 · can 2182 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.5 | 17 | %37 | %0 | 216 | 112173 | 7 |

## Sv 40 — saldırı 500 · savunma 217 · can 2216 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.4 | 15 | %32 | %0 | 244 | 126655 | 6 |

## Sv 41 — saldırı 507 · savunma 220 · can 2250 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.3 | 14 | %30 | %0 | 254 | 132037 | 6 |

## Sv 42 — saldırı 515 · savunma 223 · can 2284 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.2 | 13 | %29 | %0 | 267 | 138944 | 6 |

## Sv 43 — saldırı 522 · savunma 226 · can 2318 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.1 | 13 | %27 | %0 | 283 | 147125 | 6 |

## Sv 44 — saldırı 530 · savunma 229 · can 2352 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 3.0 | 12 | %26 | %0 | 300 | 156087 | 5 |

## Sv 45 — saldırı 537 · savunma 232 · can 2386 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.8 | 11 | %24 | %0 | 322 | 167657 | 5 |

## Sv 46 — saldırı 545 · savunma 235 · can 2420 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.7 | 10 | %22 | %0 | 344 | 179105 | 5 |

## Sv 47 — saldırı 552 · savunma 238 · can 2454 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.6 | 10 | %21 | %0 | 362 | 188067 | 5 |

## Sv 48 — saldırı 560 · savunma 241 · can 2488 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.5 | 9 | %19 | %0 | 381 | 198059 | 5 |

## Sv 49 — saldırı 567 · savunma 244 · can 2522 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.5 | 9 | %17 | %0 | 396 | 206171 | 5 |

## Sv 50 — saldırı 650 · savunma 271 · can 2742 · kritik %6 · isabet %100

| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| bog_witch | 40 | 1035 | 168 | 25 | 520 | 2.6 | 8 | %4 | %0 | 444 | 231105 | 4 |

