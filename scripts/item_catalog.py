#!/usr/bin/env python3
"""Regenerates docs/content-item-catalog.md from the content pack.

    python scripts/item_catalog.py

Reads content/{items,gear,vendors,recipes}.yaml, content/i18n/tr.yaml and the
loot tables in server/internal/item/loot.go, so the list never drifts from
what the game actually ships. Needs PyYAML.
"""
import collections
import io
import os
import re

import yaml

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def load(rel):
    with open(os.path.join(ROOT, rel), encoding="utf-8") as f:
        return yaml.safe_load(f)


items = load("content/items.yaml")["items"] + load("content/gear.yaml")["items"]
tr = load("content/i18n/tr.yaml")
vendors = load("content/vendors.yaml")["vendors"]
recipes = load("content/recipes.yaml")["recipes"]
profs = {p["code"]: p for p in load("content/professions.yaml")["professions"]}
with open(os.path.join(ROOT, "server/internal/item/loot.go"), encoding="utf-8") as f:
    loot = f.read()

KIND = {"Mining": "Madencilik", "Fishing": "Balıkçılık", "Woodcutting": "Odunculuk",
        "Herbalism": "Bitki toplama", "Hunting": "Avcılık (toplama)"}
TIER = {"Common": "Sıradan", "Uncommon": "Seçkin", "Rare": "Nadir", "Epic": "Destansı",
        "Legendary": "Efsanevi", "Mythic": "Mitik"}
STAT = {"attack": "Saldırı", "defense": "Savunma", "str": "GÜÇ", "agi": "ÇEV", "int": "ZEK", "vit": "DAY",
        "luck": "Şans", "crit_pct": "Kritik %", "hit_pct": "İsabet %", "hp_max": "Can", "mp_max": "Mana", "cha": "KAR"}
SLOT = {"mainhand": "Ana el", "offhand": "Yan el", "helm": "Kask", "chest": "Göğüs", "legs": "Bacak",
        "boots": "Çizme", "gloves": "Eldiven", "necklace": "Kolye", "ring1": "Yüzük", "trinket1": "Tılsım"}
CLASS = {"warrior": "Savaşçı", "mage": "Büyücü", "rogue": "Haydut", "priest": "Rahip", "hunter": "Avcı"}
SUB = {"ore": "cevher", "bar": "külçe", "fuel": "yakıt", "fish": "balık", "log": "kütük", "wood": "kütük",
       "herb": "bitki", "hide": "post", "pelt": "post", "potion": "iksir", "sword": "kılıç", "staff": "asa",
       "bow": "yay", "mace": "topuz", "dagger": "hançer"}


def name(code):
    return tr.get("item." + code, code)


def money(c):
    g, s, b = c // 10000, (c // 100) % 100, c % 100
    parts = [f"{g}a" if g else "", f"{s}g" if s else "", f"{b}b" if b or c == 0 else ""]
    return " ".join(p for p in parts if p)


src = collections.defaultdict(list)
for m in re.finditer(r"activity\.Kind(\w+): \{(.*?)\n\t\},", loot, re.S):
    for tier, codes in re.findall(r"(\w+):\s*\[\]string\{([^}]*)\}", m.group(2)):
        listed = re.findall(r'"(\w+)"', codes)
        for c in dict.fromkeys(listed):
            weight = f", ×{listed.count(c)} ağırlık" if listed.count(c) > 1 else ""
            src[c].append(f"{KIND[m.group(1)]} — {TIER[tier]}{weight}")
for v in vendors:
    for s in v["sells"]:
        it = next(i for i in items if i["code"] == s["template_code"])
        price = s.get("price_override_copper") or round(it["base_value_copper"] * v["buy_multiplier"])
        src[it["code"]].append(f"Dükkân ({money(price)})")
uses = collections.defaultdict(list)
for r in recipes:
    pname = tr.get("profession." + r["profession_code"], r["profession_code"])
    inputs = ", ".join(f"{x['quantity']}× {name(x['template_code'])}" for x in r["inputs"])
    cost = f", {money(r['gold_cost_copper'])}" if r.get("gold_cost_copper") else ""
    src[r["output_template_code"]].append(f"Üretim: {pname} {r['min_skill']}+ ({inputs}{cost})")
    for x in r["inputs"]:
        uses[x["template_code"]].append(name(r["output_template_code"]))
for it in items:
    if it.get("equip_slot") and it["item_class"] in ("weapon", "armor"):
        src[it["code"]].append("Av zaferi %4 / zindan boss'u")


def bonus(it):
    out = []
    for b in it.get("stat_bonuses") or []:
        t = STAT.get(b["target"], b["target"])
        out.append(f"+{b['amount']} {t}" if b["kind"] == "flat" else f"+{int(b['amount'] * 100)}% {t}")
    if it.get("use_effect"):
        out.append(f"Kullan: +{it['use_effect']['amount']} Can")
    return ", ".join(out) or "—"


lines = ["# Oyundaki Eşyalar", "",
         f"Toplam **{len(items)}** eşya şablonu. Bu dosya `python scripts/item_catalog.py` ile içerik "
         "dosyalarından üretilir; elle düzenleme. Para: 1 altın (a) = 100 gümüş (g) = 10 000 bakır (b). "
         "Satış = değer × 0.25 (genel dükkân). Toplamada kademe şansı: bir kademenin buluntusu %25 "
         "(efsanevide %15) ihtimalle bir üst nadirliğe çıkar (DD-243).", ""]
groups = [("weapon", "Silahlar"), ("armor", "Zırhlar"), ("consumable", "Tüketilebilirler"), ("material", "Malzemeler")]
for cls, title in groups:
    rows = sorted((i for i in items if i["item_class"] == cls),
                  key=lambda i: (i.get("class_restriction") or "", i.get("equip_slot") or "", i["item_subclass"],
                                 i["min_level"], i["code"]))
    if not rows:
        continue
    lines += [f"## {title} ({len(rows)})", ""]
    if cls in ("weapon", "armor"):
        lines += ["| Ad | Kod | Tür | Yuva | Sınıf | Sv | Bonuslar | Değer | Ağırlık | Nereden |",
                  "|---|---|---|---|---|---|---|---|---|---|"]
        for i in rows:
            kind = "zırh" if cls == "armor" else SUB.get(i["item_subclass"], i["item_subclass"])
            lines.append(f"| {name(i['code'])} | `{i['code']}` | {kind} | {SLOT.get(i.get('equip_slot'), '—')} | "
                         f"{CLASS.get(i.get('class_restriction'), 'Hepsi')} | {i['min_level']} | {bonus(i)} | "
                         f"{money(i['base_value_copper'])} | {i['base_weight_g'] / 1000:g} kg | "
                         f"{'; '.join(src[i['code']]) or '—'} |")
    else:
        lines += ["| Ad | Kod | Tür | Sv | Etki | Yığın | Değer | Ağırlık | Nereden | Kullanıldığı yer |",
                  "|---|---|---|---|---|---|---|---|---|---|"]
        for i in rows:
            used = ", ".join(dict.fromkeys(uses[i["code"]])) or "—"
            lines.append(f"| {name(i['code'])} | `{i['code']}` | {SUB.get(i['item_subclass'], i['item_subclass'])} | "
                         f"{i['min_level']} | {bonus(i)} | {i.get('max_stack') if i.get('stackable') else 1} | "
                         f"{money(i['base_value_copper'])} | {i['base_weight_g'] / 1000:g} kg | "
                         f"{'; '.join(src[i['code']]) or '—'} | {used} |")
    lines.append("")
with io.open(os.path.join(ROOT, "docs/content-item-catalog.md"), "w", encoding="utf-8", newline="\n") as f:
    f.write("\n".join(lines))
print(f"{len(items)} items")
