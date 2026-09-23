# Design System

## 1. Amaç
Tek referans noktası: renk, tipografi, spacing, radius, shadow, icon, motion, breakpoints. Tüm ekranlar buradan besleneceğinden değişim tek yerde yapılır.

## 2. Color Palette (Dark-first)

```yaml
palette:
  bg_primary:    "#0F1013"     # ana zemin
  bg_surface:    "#1A1C22"     # kart/panel
  bg_elevated:   "#242731"     # modal/dropdown
  border_subtle: "#2E3240"
  border_strong: "#4A4F5E"
  text_primary:  "#E8EAEE"
  text_secondary:"#A8ACB8"
  text_muted:    "#6E727E"
  accent_primary:"#E4A756"     # gold (currency, primary CTA)
  accent_secondary:"#5DA9E9"   # info/water
  success:       "#5FB973"
  warning:       "#E8B84B"
  danger:        "#E85B62"
  rarity:
    common:      "#B9BAC1"
    uncommon:    "#5FB973"
    rare:        "#5DA9E9"
    epic:        "#A26BE8"
    legendary:   "#E4A756"
    mythic:      "#E85B62"
```

Light theme Phase 02.

## 3. Typography

Sans-serif geniş desteği (TR karakter dahil):
- **Primary**: Inter (400/500/600/700).
- **Numeric/HUD**: Inter Mono veya JetBrains Mono (400/500).
- **Display (rare)**: seri font opt Phase 02.

Sizes (mobile 375px):
```yaml
type_scale:
  display:   { size: 28, weight: 700, line: 34 }
  h1:        { size: 22, weight: 700, line: 28 }
  h2:        { size: 18, weight: 600, line: 24 }
  body:      { size: 14, weight: 400, line: 20 }
  body_bold: { size: 14, weight: 600, line: 20 }
  caption:   { size: 12, weight: 400, line: 16 }
  micro:     { size: 10, weight: 500, line: 14 }
```

## 4. Spacing (4pt grid)

```yaml
space: [0, 2, 4, 8, 12, 16, 20, 24, 32, 40, 48, 64]
```

Padding kuralları: kart 16, screen 20, buton yatay 16 dikey 12.

## 5. Radius / Shadow

```yaml
radius:
  sm: 6
  md: 10
  lg: 16
  pill: 999
shadow:
  card:   "0 2 8 rgba(0,0,0,0.35)"
  modal:  "0 8 24 rgba(0,0,0,0.55)"
  focus:  "0 0 0 3 rgba(228,167,86,0.35)"  # accent_primary alpha
```

## 6. Icon System

- **Line icons** 24×24 default (feather-benzeri).
- Item slot icons 64×64 (sprite sheet content pack).
- Class/skill icons 48×48.
- Rarity border stroke 2px + accent color.

## 7. Motion

```yaml
duration:
  quick:  150   # tap ripple
  fast:   200   # sheet transition
  base:   300   # page transition
  slow:   400   # reveal loot
  epic:   600   # legendary reveal
easing:
  default: easeInOut
  emphasize: cubicEmphasize (0.2, 0, 0, 1)
```

Reduced motion accessibility → tüm süreler *0.5 (setting).

## 8. Breakpoints

| Aralık | Layout |
|---|---|
| ≤ 480 | mobile (default) |
| 481–768 | large mobile / small tablet — inventory 2-col |
| 769–1024 | tablet — side rail + content |
| 1025+ | desktop/web — 3-pane (left rail + content + context) |

## 9. Components Katalog (Base)

- **Button** primary/secondary/ghost/danger, size sm/md/lg.
- **Chip / Tag** rarity badge, status pill.
- **Card** default/elevated/interactive.
- **Modal / Bottom Sheet** — mobile bottom sheet primary; desktop centered modal.
- **List Item** avatar + title + trailing.
- **Progress** linear + radial (energy, XP, quest).
- **Slot** (item) — icon + rarity border + count.
- **Toast / Snackbar** — success/warn/error, auto-dismiss 4s.
- **Skeleton** — 300ms shimmer.

## 10. Reveal / Loot Motif

Loot animation: item card scale 0.8→1.0 with rarity aura + accent color pulse. Legendary+: full-screen shine flash.

## 11. Localization Rules

- Body key + vars pattern (server payload uyumlu).
- Sayı format: locale ile (TR: `1.234,56`); currency 3 sütun (gold/silver/copper).
- Tarih format: relative ("2 saat önce") + tooltip absolute.

## 12. Cross-Reference
- [architecture/flutter.md](../architecture/flutter.md)
- [notification/mobile-lifecycle.md](../notification/mobile-lifecycle.md)

## 13. Bekleyen Kararlar
- Light theme Phase 02.
- Custom fantasy display font Phase 02.
- Rarity legendary/mythic animation art asset Phase 02.
