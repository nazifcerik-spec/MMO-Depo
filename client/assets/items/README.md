# Item art (DD-240)

One picture per item template, named exactly after its code:

    boar_hide.png        → Domuz Postu
    sword_training.png   → Eğitim Kılıcı

- Code = the `code:` field in `content/items.yaml` / `content/gear.yaml`
  (list: `docs/content-item-catalog.md`). Lowercase, underscores, `.png`.
- Square, transparent background, 128×128 px (shown at ~33 px; decoded at
  display size).
- One picture per item: rarity is the cell's border colour, affixes
  ("Kadim …") use the base item's picture.
- An item without a file keeps its class icon, so art can arrive gradually.
- Pictures ship inside the app: new art needs a new APK build.
