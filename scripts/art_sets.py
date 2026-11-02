"""Picture sets beyond gear (DD-332, DD-333): class path nodes, map regions, crafting
professions, helper NPCs, a few materials, dungeons, quests and class portraits, each painted by Gemini in 3x3
sheets in the game's style and cut into client/assets/<dir>/<code>.webp.

    python scripts/art_sets.py plan
    python scripts/art_sets.py prompt <set> <sheet>   # Gemini prompt
    python scripts/art_sets.py cut <set> [sheet ...]  # client/art_originals/<set>/sheetNN.jpg

Medallion sets (nodes, regions, professions) match the class emblems
(DD-325): one round medallion with an ornate gold rim per cell, cut to the
medallion. Portraits match the avatars (DD-324); items match the gear icons
(DD-329) and land among them in assets/items.
"""
import io
import os
import re
import sys

from PIL import Image

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from emblem_assets import medallion  # noqa: E402
from item_art import gutters  # noqa: E402

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
GRID = 3
POS = ['top-left', 'top-middle', 'top-right', 'middle-left', 'center', 'middle-right',
       'bottom-left', 'bottom-middle', 'bottom-right']

# Base class palettes, as the class emblems use them.
PAL = {
    'warrior': 'steel and crimson', 'rogue': 'black and toxic green', 'ranger': 'forest green and brown leather',
    'mage': 'blue and violet', 'monk': 'orange and gold', 'cleric': 'gold and white',
    'paladin': 'gold and royal blue', 'druid': 'green and amber', 'bard': 'magenta and gold',
    'shaman': 'teal and orange',
}

# Promotions and specializations (DD-252), by base class: what each emblem shows.
NODES = [
    ('warrior', [
        ('guardian', 'a massive tower shield with a lion crest'),
        ('executioner', 'a huge executioner greataxe with a red glow'),
        ('iron_fortress', 'a fortress gate tower behind a riveted iron shield'),
        ('warlord', 'a war banner on a spear over two crossed swords'),
        ('berserker', 'a roaring horned helm wreathed in red rage flames'),
        ('weapon_master', 'a fan of five different weapons: sword, axe, spear, hammer and glaive'),
    ]),
    ('rogue', [
        ('assassin', 'a hooded skull mask with a dagger behind it'),
        ('nightblade', 'a curved black blade under a crescent moon'),
        ('poison_master', 'a bubbling poison vial wrapped by a green snake'),
        ('duelist', 'two crossed rapiers with a red rose'),
        ('blade_dancer', 'twin spinning blades inside a swirling silk ribbon'),
        ('trickster', 'a grinning jester mask with fanned playing cards'),
    ]),
    ('ranger', [
        ('marksman', 'a bullseye target pierced by an arrow'),
        ('sharpshooter', 'a fierce eagle eye inside a crosshair ring'),
        ('crossbow_master', 'a heavy loaded crossbow'),
        ('scout', 'a brass spyglass over a compass rose'),
        ('beastmaster', 'a wolf head side by side with a hawk'),
        ('trapper', 'a steel bear trap with a rope snare'),
    ]),
    ('mage', [
        ('arcanist', 'an open spellbook with floating glowing runes'),
        ('rune_master', 'a carved standing runestone glowing with runes'),
        ('chronomancer', 'a clock face with swirling sands of time'),
        ('elementalist', 'four small orbs of fire, water, earth and air circling each other'),
        ('pyromancer', 'a blazing fireball held in a clawed hand of flame'),
        ('cryomancer', 'a jagged ice crystal snowflake with frost mist'),
    ]),
    ('monk', [
        ('disciple', 'a calm meditating figure seated on a lotus'),
        ('iron_body', 'an iron-banded fist pressed on a cracked stone'),
        ('spirit_walker', 'a glowing spirit wisp above a ring of prayer beads'),
        ('fighter', 'two cloth-wrapped fists bumping together'),
        ('combo_master', 'a fist trailing a chain of three impact rings'),
        ('storm_fist', 'a fist crackling with blue lightning'),
    ]),
    ('cleric', [
        ('saint', 'a radiant halo above two praying hands'),
        ('high_priest', 'an ornate golden mitre with a sun emblem'),
        ('aegis_priest', 'a glowing dome of holy light over a small sun sigil'),
        ('oracle', 'an all-seeing eye inside a starry crystal ball'),
        ('prophet', 'an unrolled glowing scroll with a bright star'),
        ('exorcist', 'a golden chain binding a dark shadow wraith'),
    ]),
    ('paladin', [
        ('protector', 'a kite shield with spread golden wings'),
        ('aegis_knight', 'a knight great helm in front of a radiant round shield'),
        ('lightwarden', 'a lantern of holy light hanging over a shield'),
        ('crusader', 'a raised longsword with a long blue banner'),
        ('templar', 'a white temple pillar with a sword leaning on it'),
        ('justicar', 'golden scales of justice balanced on a sword'),
    ]),
    ('druid', [
        ('grove_keeper', 'a great oak tree with glowing roots'),
        ('life_binder', 'a sprouting seedling cupped in vines around a heart of light'),
        ('thorn_sage', 'a twisted thorny vine with a purple bloom'),
        ('shapeshifter', 'a face split half wolf and half human'),
        ('bearwarden', 'a roaring bear head armoured in bark'),
        ('mooncaller', 'a howling wolf under a huge full moon'),
    ]),
    ('bard', [
        ('minstrel', 'a lute with floating musical notes'),
        ('virtuoso', 'a violin and bow with golden sparkles'),
        ('harmonist', 'a small harp with glowing green healing notes'),
        ('skald', 'a curved war horn with knotwork carvings'),
        ('war_chanter', 'a war drum with two crossed drumsticks'),
        ('dirge_singer', 'a black raven perched on a broken lyre'),
    ]),
    ('shaman', [
        ('totemist', 'a carved totem pole with eagle wings'),
        ('earthwarden', 'a huge stone fist rising out of the ground'),
        ('storm_totemist', 'a totem struck by lightning under storm clouds'),
        ('spirit_caller', 'a ghostly ancestral spirit mask'),
        ('ancestor_sage', 'an elder spirit face in smoke with feathers'),
        ('hexer', 'a cursed charm of bones and beads with purple glow'),
    ]),
]

# Map regions (DD-279): a painted landscape vignette inside each medallion.
REGIONS = [
    ('sundermere', 'a green island continent with rolling hills and a castle seen from above'),
    ('ironhold_city', 'a stone fortress city with iron gates and smoking chimneys'),
    ('greenwood', 'a lush green forest with sunbeams'),
    ('whispering_reach', 'a misty river valley with tall reeds'),
    ('crypt_gates', 'a crypt gate carved with skulls'),
    ('emberpeak', 'a volcanic mountain peak glowing with embers'),
    ('ashen_hollow', 'a cave of grey ash with a red glow inside'),
    ('broken_crown', 'a cracked golden crown above a kingdom of castles'),
    ('dragonspine_isles', 'a chain of islands shaped like a dragon spine in the sea'),
    ('starless_lands', 'dark plains under a black starless sky'),
    ('sky_rift', 'floating rocks around a glowing tear in the sky'),
    ('void_shore', 'a dark shore against a purple cosmic void'),
    ('misty_marshes', 'a foggy swamp with dead trees'),
    ('sunken_temple', 'a half-sunken stone temple in green water'),
    ('iron_mountains', 'grey jagged mountains with a mine cart track'),
    ('forgotten_mine', 'an abandoned mine entrance with old lanterns'),
    ('northern_keep_ruins', 'snowy castle ruins'),
    ('northern_keep', 'a frozen keep gate with tattered banners'),
    ('crown_plains', 'golden wheat plains with a royal road'),
    ('rebel_camp', 'tents and campfires with a torn banner'),
    ('burnt_forest', 'charred black trees with glowing embers'),
    ('fire_root', 'a cavern of burning tree roots'),
    ('salt_desert', 'white salt flats with cracked ground'),
    ('buried_palace', 'a palace dome half buried in sand'),
    ('crown_graveyard', 'a graveyard of tombstones topped with crowns'),
    ('royal_crypt', 'an ornate royal crypt door with gold trim'),
    ('storm_coast', 'stormy sea cliffs with lightning'),
    ('sunken_fleet', 'old shipwrecks under blue water'),
    ('dragon_slopes', 'rocky slopes scattered with dragon bones'),
    ('brood_cavern', 'a cavern nest with dragon eggs'),
    ('mist_forest', 'a misty blue forest of ancient trees'),
    ('lost_grove', 'a hidden glowing grove with a stone circle'),
    ('crystal_caves', 'a cave full of blue crystals'),
    ('crystal_heart', 'a giant glowing crystal heart in a cavern'),
    ('volcanic_isle', 'an island volcano erupting'),
    ('volcano_heart', 'a lava chamber inside a volcano'),
    ('throne_peaks', 'snowy mountain peaks with a throne carved in the rock'),
    ('dragon_throne', 'a huge throne made of a dragon skull'),
    ('dark_plains', 'dark plains with dead grass under a black moon'),
    ('darkened_observatory', 'a ruined observatory dome with a broken telescope'),
    ('lost_city', 'overgrown ruined city towers'),
    ('clock_tower', 'a crooked clock tower with glowing gears'),
    ('floating_isles', 'floating islands with waterfalls in the sky'),
    ('storm_fortress', 'a fortress inside thunderclouds'),
    ('rift_edge', 'a cliff edge facing a glowing purple sky rift'),
    ('rift_core', 'a swirling purple vortex core'),
    ('celestial_temple', 'a white and gold temple above the clouds'),
    ('celestial_dome', 'a starlit dome with angel statues'),
    ('ash_sea', 'a grey sea of ash with burning ships'),
    ('sunken_throne', 'a throne half drowned in ash waves'),
    ('last_bastion', 'towering bastion walls at the edge of a dark void'),
    ('bastion_keep', 'a fortified keep gate glowing with blue wards'),
    ('edge_of_eternity', 'a stairway rising into endless starlight'),
    ('eternal_throne', 'a cosmic throne made of stars'),
]

# Crafting professions (Workshop list).
PROFESSIONS = [
    ('blacksmithing', 'a hammer striking a glowing sword blade on an anvil'),
    ('leatherworking', 'a tanned leather hide with an awl and stitching thread'),
    ('tailoring', 'a needle, a spool of thread and a folded cloth roll'),
    ('woodworking', 'a carving chisel and a wooden bow stave with curled shavings'),
    ('jewelcrafting', 'a cut gemstone with a jeweller loupe and a gold ring'),
    ('alchemy', 'a round bubbling flask with colourful vapour'),
    ('cooking', 'a steaming cooking pot with a ladle'),
    ('smelting', 'a crucible pouring molten metal into an ingot mould'),
    ('enchanting', 'a glowing parchment scroll above a magic circle'),
]

# Hired helpers (DD-319): head-and-shoulders portraits of the workers.
HELPERS = [
    ('miner', 'a stout bearded miner with a lamp helmet and a pickaxe over his shoulder'),
    ('woodsman', 'a burly lumberjack with a red plaid shirt and an axe'),
    ('angler', 'an old fisherman with a straw hat and a fishing rod'),
    ('herbalist', 'a young woman herbalist with a flower crown and a basket of herbs'),
    ('trapper', 'a hooded trapper in a fur cloak with a coil of rope'),
    ('smith', 'a muscular blacksmith in a leather apron with soot on his face'),
    ('tanner', 'a woman tanner in a leather apron holding a hide'),
    ('weaver', 'a woman weaver with a colourful shawl and thread spools'),
    ('carver', 'an elf woodcarver with wood shavings in her hair and a chisel'),
    ('apothecary', 'a bespectacled gnome apothecary holding a potion flask'),
    ('cook', 'a jolly round cook with a chef hat and a ladle'),
    ('smelter', 'an orc smelter with goggles lit by furnace glow'),
    ('scribe', 'an elderly scribe with a quill and a glowing scroll'),
    ('gemcutter', 'a dwarf woman gemcutter with a loupe and a sparkling gem'),
]

# Materials still without art; None = a filler cell, not cut.
MATERIALS = [
    ('elm_log', 'a short bundle of pale elm logs with light grey bark'),
    ('yew_log', 'a stack of reddish-brown yew logs with dark flaky bark'),
    ('linseed', 'a small burlap pouch spilling golden-brown flax seeds'),
    ('mandrake', 'a mandrake root with a gnarled humanlike root body and green leaves'),
    ('king_bloom', 'a regal golden and purple flower with a crown-shaped bloom'),
    ('mercury_ore', 'a chunk of grey rock with beads of liquid silver mercury'),
    (None, 'a coil of hemp rope'),
    (None, 'a lump of black coal'),
    (None, 'a bundle of white cotton'),
]


# Redraws in the last sheet's free cells: a later sheet's cut wins.
NODE_FIXES = [
    ('monk', 'disciple', 'a lotus blossom holding a floating glowing pearl'),
]


# Dungeons and raids (DD-333): each boss in its lair inside a medallion.
DUNGEONS = [
    ('crypt_of_ashes', 'a goblin raider chieftain with a torch in an ashen crypt'),
    ('forest_hollow', 'a huge snarling dire wolf in a dark forest hollow'),
    ('shadowfen_ruins', 'a shadow stalker with glowing eyes among swamp ruins'),
    ('obsidian_spire', 'an ashen lich with a staff atop a black obsidian spire'),
    ('sunken_temple', 'a marsh queen made of vines and water in a sunken temple'),
    ('forgotten_mine', 'an iron king in heavy armour in a mine cavern'),
    ('northern_keep', 'a fallen lord in frosted armour on a frozen throne'),
    ('rebel_camp', 'a rebel king with a torn cape among war tents'),
    ('fire_root', 'a burning mother tree with a face and roots of fire'),
    ('buried_palace', 'a mummified pharaoh with a golden mask in a sandy palace hall'),
    ('royal_crypt', 'a crowned crypt lord wraith in a golden crypt'),
    ('sunken_fleet', 'a drowned admiral with a cutlass among shipwrecks'),
    ('brood_cavern', 'a mother dragon guarding her eggs in a cavern'),
    ('lost_grove', 'a spirit queen of glowing leaves in a stone circle grove'),
    ('crystal_heart', 'a titan made of blue crystals in a crystal cave'),
    ('volcano_heart', 'a fire lord of molten lava in a volcano'),
    ('dragon_throne', 'a dragon king on a throne of bones'),
    ('darkened_observatory', 'a cosmic horror swallowing stars under a ruined observatory dome'),
    ('clock_tower', 'a time lord with clock gears floating around him in a clock tower'),
    ('storm_fortress', 'a storm king with a lightning crown above a fortress in the clouds'),
    ('rift_core', 'a many-eyed void being in a swirling purple vortex'),
    ('celestial_dome', 'a radiant winged warrior with a flaming sword under a starry dome'),
    ('sunken_throne', 'a leviathan king rising from grey ash waves'),
    ('bastion_keep', 'a doom commander in black spiked armour on bastion walls'),
    ('eternal_throne', 'a fallen god on a throne of stars'),
    ('north_keep_siege', 'a giant undying undead warlord besieging a snowy keep, epic raid scale'),
    ('crownless_palace', 'a colossal tyrant king with a broken crown in a ruined palace, epic raid scale'),
    ('dragon_king_lair', 'a huge ascended golden dragon king on a hoard of treasure, epic raid scale'),
    ('time_lord_tower', 'a titanic timeless lord in a tower of floating clocks, epic raid scale'),
    ('archangels_fall', 'dark winged warriors falling from a shattered golden sky, epic raid scale'),
    ('fallen_god_throne', 'a titanic god of void and stars breaking his chains, epic raid scale'),
]

# Quests (DD-333): what each one is about. Main roads show the road's
# monster before the dungeon; bounties the hunted monster in a red bounty ring.
ROADS = [
    ('shadowfen_ruins', 'a hooded cultist acolyte', 'swamp ruins'),
    ('obsidian_spire', 'a frost giant', 'a black obsidian spire'),
    ('sunken_temple', 'a giant mist leech', 'a temple sunk in a swamp'),
    ('forgotten_mine', 'a mountain troll', 'an old mine entrance'),
    ('northern_keep', 'a ghostly guard', 'a frozen keep'),
    ('rebel_camp', 'a rebel rider on horseback', 'a camp of war tents'),
    ('fire_root', 'an ash wolf with ember fur', 'a burning tree'),
    ('buried_palace', 'a giant sand scorpion', 'a palace buried in dunes'),
    ('royal_crypt', 'a crowned royal skeleton', 'a royal crypt'),
    ('sunken_fleet', 'a sea raider with a cutlass', 'wrecked ships in a bay'),
    ('brood_cavern', 'a young red dragon', 'a dragon cavern'),
    ('lost_grove', 'a glowing fey wolf', 'a hidden grove'),
    ('crystal_heart', 'a crystal spider', 'a crystal cave'),
    ('volcano_heart', 'a magma golem', 'an erupting volcano'),
    ('dragon_throne', 'a dragon knight in scaled armour', 'a mountain throne'),
    ('darkened_observatory', 'a shadow rider on a dark steed', 'a ruined observatory'),
    ('clock_tower', 'a ghost noble in old finery', 'a crooked clock tower'),
    ('storm_fortress', 'a cloud giant', 'a fortress in storm clouds'),
    ('rift_core', 'a void creature with tentacles', 'a purple sky rift'),
    ('celestial_dome', 'a fallen angel with dark wings', 'a temple above the clouds'),
    ('sunken_throne', 'an ash leviathan', 'a throne in an ash sea'),
    ('bastion_keep', 'a doom rider in black armour', 'towering bastion walls'),
    ('eternal_throne', 'a void dragon', 'a stairway into the stars'),
]
BOUNTIES = [
    ('drowned_knight', 'a drowned knight in rusted armour'), ('mine_warden', 'an armoured mine warden with a lantern'),
    ('fallen_commander', 'a fallen undead commander'), ('black_lancer', 'a black lancer knight with a long lance'),
    ('scorched_shaman', 'a scorched shaman with a burning staff'), ('sand_worm', 'a giant sand worm bursting from dunes'),
    ('crownless_king', 'a crownless undead king'), ('kraken_spawn', 'a young kraken with curling tentacles'),
    ('dragon_guard', 'a dragon guard warrior with a halberd'), ('forest_spirit', 'a forest spirit of leaves and light'),
    ('cave_titan', 'a stone cave titan'), ('volcano_warden', 'a volcano warden of obsidian and lava'),
    ('elder_dragon', 'an ancient elder dragon'), ('night_giant', 'a night giant with starry skin'),
    ('city_guardian', 'a stone city guardian statue come alive'), ('thunderbird', 'a thunderbird with lightning wings'),
    ('chaos_lord', 'a chaos lord in shifting purple armour'), ('seraph', 'a six-winged radiant seraph warrior'),
    ('eternal_sentinel', 'an eternal sentinel construct of gold and starlight'),
    ('last_commander', 'the last commander in battered silver armour'), ('avatar_of_ruin', 'an avatar of ruin, a burning colossus'),
]
QUESTS = [
    ('main_01_first_hunt', 'a young boar charging through tall grass'),
    ('main_02_supplies', 'a pickaxe and a hatchet over a pile of ore and logs'),
    ('main_03_forge', 'a glowing forge with a hammer on an anvil'),
    ('main_04_wolves', 'grey wolves before a wooden village gate'),
    ('main_05_crypt', 'a torch lighting stone stairs down into a crypt'),
    ('main_06_hollow', 'a hulking forest troll in a dark hollow'),
    ('side_bandit_trouble', 'a masked bandit scout with a short bow'),
    ('side_healers_stock', 'a red healing potion beside fresh herbs'),
    ('side_fishmonger', 'fresh fish on a market stall with a hanging scale'),
    ('side_goblin_menace', 'a snarling goblin raider with a crude axe'),
    ('daily_hunter', 'two crossed hunting spears in front of a rising sun'),
    ('daily_miner', 'a pickaxe in a rock in front of a rising sun'),
    ('daily_angler', 'a fishing rod and a leaping fish in front of a rising sun'),
    ('daily_lumberjack', 'an axe in a tree stump in front of a rising sun'),
    ('daily_herbalist', 'a sickle with a bundle of herbs in front of a rising sun'),
    ('daily_artisan', 'a hammer and a chisel in front of a rising sun'),
    ('daily_delver', 'a lantern at a dungeon door in front of a rising sun'),
    ('daily_slayer', 'a sword over a horned skull in front of a rising sun'),
] + [('main_%02d_%s' % (7 + i, code), '%s blocking a road that leads to %s' % (mob, place))
     for i, (code, mob, place) in enumerate(ROADS)] + [
    ('side_bounty_' + code, '%s inside a red hunter bounty ring' % d) for code, d in BOUNTIES]

# Enchanting scrolls (DD-333): one bold symbol per stat in its colour, painted
# large on the open scroll; the scroll itself grows richer with the tier.
SCROLL_SYMBOLS = [
    ('weapon_attack', 'two crossed red swords'),
    ('weapon_crit', 'a jagged yellow lightning bolt'),
    ('weapon_strength', 'a clenched orange fist'),
    ('weapon_agility', 'a green winged boot'),
    ('weapon_intellect', 'a glowing blue crystal orb'),
    ('weapon_wisdom', 'a purple owl head'),
    ('weapon_lifesteal', 'a crimson blood drop with two fangs'),
    ('armor_defense', 'a steel-blue shield'),
    ('armor_health', 'a red heart'),
    ('armor_vitality', 'a green oak tree'),
    ('armor_evasion', 'a teal swirl of wind'),
    ('jewelry_crit', 'a golden four-pointed burst star'),
    ('jewelry_hit', 'a bullseye target with an arrow in its centre'),
    ('jewelry_luck', 'a bright green four-leaf clover'),
    ('jewelry_crit_damage', 'a cracked skull with an orange impact burst'),
    ('jewelry_mining', 'a grey pickaxe'),
    ('jewelry_herbalism', 'a green herb sprig with flowers'),
    ('jewelry_woodcutting', 'an axe stuck in a log'),
    ('jewelry_fishing', 'a blue fish on a hook'),
    ('jewelry_hunting', 'a brown bear paw print'),
    ('jewelry_craft_speed', 'a golden hourglass with a small hammer'),
]
SCROLL_TIERS = [
    'plain tan parchment scroll held wide open, a twine tassel and a red wax seal on its edge',
    'cream parchment scroll with a leather tie and a red wax seal',
    'parchment scroll with a blue ribbon and a silver seal',
    'fine vellum scroll with a green ribbon and a silver seal',
    'vellum scroll with gold edging and a gold seal',
    'deep blue scroll with gold rollers',
    'purple scroll with gem-studded rollers',
    'crimson scroll with ornate gold rollers and a soft glow',
    'black scroll with silver filigree and a glowing aura',
    'radiant golden-white scroll with a divine glow and floating sparks',
]
SCROLLS = [('enchant_%s_%d' % (code, t + 1),
            'unrolled %s, its face one big bold symbol: %s' % (SCROLL_TIERS[t], sym))
           for code, sym in SCROLL_SYMBOLS for t in range(10)]
# Redraws in the last sheet's free cells (a later sheet's cut wins): first
# tiers that came out as closed bundles without their symbol.
SCROLL_FIXES = ['enchant_weapon_attack_1', 'enchant_weapon_crit_1']
SCROLLS += [(c, d) for c, d in SCROLLS if c in SCROLL_FIXES]

# Class portraits (DD-333): one vertical picture per base class, painted like
# the race cards (DD-266) and shown 3:4 on the creation screen.
CLASS_PORTRAITS = [
    ('warrior', 'a human man warrior in battle-worn steel and crimson plate armour with a sword and shield'),
    ('rogue', 'a dark elf woman rogue in a black hooded leather outfit with twin daggers and green poison glow'),
    ('ranger', 'a high elf man ranger in green and brown leather with a longbow and a feathered quiver'),
    ('mage', 'a human woman mage in blue and violet arcane robes casting glowing blue magic'),
    ('monk', 'a bull-headed minotaur monk in orange and gold martial wraps with prayer beads, fists raised'),
    ('cleric', 'a dwarf woman cleric in white and gold vestments with a mace and soft holy light'),
    ('paladin', 'a dragon-kin man paladin in gleaming gold and royal blue plate with a warhammer'),
    ('druid', 'a high elf woman druid in leafy robes with an antler headdress and a glowing nature staff'),
    ('bard', 'a gnome woman bard in a colourful magenta and gold doublet playing a lute'),
    ('shaman', 'a green-skinned orc man shaman in tribal mail with a totem staff and crackling lightning'),
]

MOBS = [  # content/mobs.yaml order; the scene is the mob's home region
    ('young_boar', 'a young brown boar with small tusks in a sunny green forest'),
    ('grey_wolf', 'a grey wolf baring its teeth in a green forest'),
    ('bandit_scout', 'a masked bandit scout with a short bow in a forest'),
    ('cave_bat', 'a big-eared cave bat with spread wings in a dark cave mouth'),
    ('goblin_raider', 'a snarling green goblin raider with a crude axe on a misty shore'),
    ('forest_troll', 'a hulking mossy forest troll with a club among dark trees'),
    ('dire_wolf', 'a huge black dire wolf with red eyes in the mist'),
    ('cultist_acolyte', 'a hooded cultist acolyte with a glowing red ritual dagger'),
    ('rock_ogre', 'a grey rock-skinned ogre with a boulder on a volcanic ridge'),
    ('shadow_stalker', 'a lean shadow stalker of black smoke with glowing eyes'),
    ('elder_wyrm', 'an old scaly wyrm with long horns coiled on a volcanic peak'),
    ('ashen_lich', 'a skeletal ashen lich with a burning staff and a tattered crown'),
    ('marsh_lurker', 'a slimy green marsh lurker rising from reeds and murky water'),
    ('frost_giant', 'a blue-skinned frost giant with an icy beard and a stone club'),
    ('mist_leech', 'a giant pale leech with a round toothy mouth in swamp mist'),
    ('mud_golem', 'a hulking golem of dripping swamp mud with glowing eyes'),
    ('bog_witch', 'a hunched bog witch with a gnarled staff and green swamp light'),
    ('marsh_queen', 'a regal marsh queen of vines and lilies with a reed crown'),
    ('venom_lizard', 'a green venom lizard dripping glowing poison in a swamp'),
    ('drowned_knight', 'a drowned knight in rusted armour draped with seaweed'),
    ('mountain_troll', 'a huge grey mountain troll with tusks in a rocky pass'),
    ('dwarf_ghost', 'a translucent blue dwarf ghost with a pickaxe in a mine'),
    ('rock_eagle', 'a great brown eagle with stone-grey wings on a cliff'),
    ('iron_king', 'a stern dwarf king in heavy iron armour with an iron crown'),
    ('iron_spider', 'a mechanical spider of black iron plates in a mine shaft'),
    ('mine_warden', 'an armoured mine warden with a lantern and a heavy pick'),
    ('ghost_guard', 'a pale ghost soldier with a spectral spear in ruined walls'),
    ('cursed_archer', 'an undead archer with a green cursed glow drawing a bow'),
    ('bone_rider', 'a skeleton knight riding a skeletal horse'),
    ('fallen_lord_varn', 'a fallen undead lord in black plate with a cold blue crown'),
    ('keep_gargoyle', 'a winged stone gargoyle crouched on a ruined tower'),
    ('fallen_commander', 'an undead commander in dented plate raising a broken sword'),
    ('rebel_rider', 'a rebel rider with a red scarf on a brown warhorse in golden plains'),
    ('plains_lion', 'a golden maned lion roaring in tall golden grass'),
    ('raider_mage', 'a raider battle mage in a red cloak casting fire in the plains'),
    ('rebel_king_orrin', 'a proud rebel king with a red cloak and a scarred face'),
    ('war_bear', 'a huge brown bear wearing spiked war armour'),
    ('black_lancer', 'a black-armoured lancer knight with a long lance'),
    ('ash_wolf', 'a wolf with ember-glowing fur in a burnt forest'),
    ('flame_spirit', 'a dancing spirit of orange flame among charred trees'),
    ('cinder_treant', 'a smouldering tree giant with glowing cinder cracks'),
    ('burnt_mother_tree', 'a vast burnt mother tree with an angry face and burning branches'),
    ('crimson_harpy', 'a crimson-feathered harpy with sharp talons over burnt trees'),
    ('scorched_shaman', 'a scorched shaman with a burning staff and ash-grey skin'),
    ('sand_scorpion', 'a giant tan sand scorpion with raised stinger on white salt dunes'),
    ('desert_viper', 'a striped desert viper with open fangs on salt flats'),
    ('salt_elemental', 'a jagged elemental made of white salt crystals in a desert'),
    ('pharaoh_seth_ka', 'a mummified pharaoh with a golden mask and a crook staff'),
    ('mummy_priest', 'a bandaged mummy priest with a glowing staff'),
    ('sand_worm', 'a giant sand worm bursting from the dunes with a ring of teeth'),
    ('royal_skeleton', 'a crowned royal skeleton in a torn purple robe in a graveyard'),
    ('tomb_keeper', 'a hulking tomb keeper in stone armour with a lantern in a crypt'),
    ('soul_drinker', 'a gaunt soul drinker wraith pulling glowing souls'),
    ('crypt_lord_malachar', 'a crowned crypt lord wraith with gold armour and green flame eyes'),
    ('black_knight', 'a black knight with a horned helm and a greatsword among graves'),
    ('crownless_king', 'an undead king with a broken crown on a dark throne'),
    ('sea_raider', 'a sea raider with a cutlass and a tricorn hat on a stormy shore'),
    ('storm_bird', 'a great blue storm bird with crackling feathers over waves'),
    ('coral_golem', 'a hulking golem of pink and orange coral on a beach'),
    ('admiral_drask', 'a drowned pirate admiral with a long coat and cutlass'),
    ('siren', 'an eerie siren with a fish tail on a sea rock'),
    ('kraken_spawn', 'a young kraken with curling tentacles in stormy water'),
    ('young_dragon', 'a young red dragon with small wings on a rocky slope'),
    ('dragon_cultist', 'a dragon cultist in red robes with a dragon-skull mask'),
    ('lava_lizard', 'a lizard with glowing lava cracks on its black scales'),
    ('mother_dragon_vyrna', 'a great mother dragon guarding a nest of eggs'),
    ('winged_wyvern', 'a green winged wyvern with a barbed tail in flight'),
    ('dragon_guard', 'a scaled dragonkin guard with a halberd'),
    ('fey_wolf', 'a glowing pale-green fey wolf in a misty enchanted forest'),
    ('mist_stalker', 'a long-limbed grey mist stalker half hidden in fog'),
    ('ancient_treant', 'an ancient mossy treant with a wise bark face'),
    ('spirit_queen_sylwen', 'a spirit queen made of glowing leaves and light'),
    ('shadow_elf', 'a shadow elf assassin with dark skin and twin curved blades'),
    ('forest_spirit', 'a gentle glowing forest spirit with an antlered head'),
    ('crystal_spider', 'a spider made of blue crystal in a glittering cave'),
    ('stone_eater', 'a rock-chewing beast with a huge stony jaw in a crystal cave'),
    ('light_elemental', 'a radiant elemental of white-gold light in a crystal cave'),
    ('crystal_titan', 'a giant titan made of blue crystals'),
    ('crystal_dragon', 'a dragon with clear crystal scales and glowing blue eyes'),
    ('cave_titan', 'a huge stone cave titan with mossy shoulders'),
    ('magma_golem', 'a golem of black rock and flowing magma on a volcanic island'),
    ('fire_giant', 'a fire giant with a flaming beard and a molten hammer'),
    ('ash_dragon', 'a grey ash dragon breathing smoke over a volcano'),
    ('fire_lord_ignar', 'a fire lord made of roaring flames and molten armour'),
    ('volcano_shaman', 'a volcano shaman with an obsidian mask and a lava staff'),
    ('volcano_warden', 'a warden of obsidian and lava with a great shield'),
    ('dragon_knight', 'a knight in red dragon-scale armour on a mountain peak'),
    ('ancient_wyvern', 'an ancient scarred wyvern with tattered wings on high peaks'),
    ('dragon_sorcerer', 'a horned dragonkin sorcerer casting purple fire'),
    ('dragon_king_azhorath', 'a crowned dragon king on a throne of bones'),
    ('throne_guardian', 'a giant armoured throne guardian with a tower shield'),
    ('elder_dragon', 'an enormous elder dragon with golden horns among peaks'),
    ('shadow_rider', 'a shadow rider on a dark steed under a starless sky'),
    ('void_wolf', 'a wolf of black void with star-speckled fur'),
    ('star_gnawer', 'a ravenous void beast gnawing on a glowing fallen star'),
    ('star_eater_omh', 'a cosmic horror swallowing stars with a huge maw'),
    ('dark_angel', 'a dark angel with black feathered wings and a dark sword'),
    ('night_giant', 'a night giant with deep blue starry skin'),
    ('ghost_noble', 'a ghost noble in old lace finery in a lost city'),
    ('time_thief', 'a sly hooded thief holding a glowing hourglass'),
    ('cursed_sorcerer', 'a cursed sorcerer with cracked skin and glowing purple marks'),
    ('time_lord_kaedrin', 'a time lord with floating clock gears around him'),
    ('stone_gargoyle', 'a hunched stone gargoyle with cracked wings in ancient ruins'),
    ('city_guardian', 'a giant stone guardian statue come alive in a lost city'),
    ('cloud_giant', 'a cloud giant with a white beard among floating islands'),
    ('storm_dragon', 'a blue storm dragon wreathed in lightning above clouds'),
    ('sky_knight', 'a winged sky knight in silver armour with a lance'),
    ('storm_king_zephyros', 'a storm king with a lightning crown and a trident'),
    ('wind_elemental', 'a swirling wind elemental with a face in the vortex'),
    ('thunderbird', 'a thunderbird with blazing lightning wings'),
    ('void_creature', 'a tentacled void creature by a purple sky rift'),
    ('rift_warden', 'an armoured rift warden with a glowing purple shield'),
    ('chaos_elemental', 'a shifting chaos elemental of purple and red energy'),
    ('chaos_mother', 'a many-eyed chaos mother in a swirling purple vortex'),
    ('rift_dragon', 'a dragon with torn purple void wings'),
    ('chaos_lord', 'a chaos lord in shifting purple spiked armour'),
    ('fallen_angel', 'a fallen angel with ash-grey wings and a broken halo'),
    ('temple_guardian', 'a golden temple guardian statue with a glaive above the clouds'),
    ('celestial_dragon', 'a white and gold celestial dragon with star patterns'),
    ('archangel_uriel', 'a radiant archangel with a flaming sword and gold armour'),
    ('light_hunter', 'a light hunter with a glowing white bow above the clouds'),
    ('seraph', 'a six-winged radiant seraph warrior'),
    ('ash_leviathan', 'a grey ash leviathan rising from an ash sea'),
    ('null_hunter', 'a faceless null hunter in grey cloth with a hooked blade'),
    ('undying_soldier', 'an undying soldier in ash-covered armour'),
    ('leviathan_king', 'a crowned leviathan king roaring from grey ash waves'),
    ('ash_wyrm_lord', 'a long ash wyrm lord coiled in smoke'),
    ('eternal_sentinel', 'an eternal sentinel construct of gold and starlight'),
    ('doom_rider', 'a doom rider in black armour on a burning horse'),
    ('void_titan', 'a towering void titan with a body of black stars'),
    ('soul_reaper', 'a hooded soul reaper with a great glowing scythe'),
    ('doom_commander_morvath', 'a doom commander in black spiked armour on bastion walls'),
    ('ancient_golem', 'an ancient stone golem with glowing blue cracks'),
    ('last_commander', 'the last commander in battered silver armour'),
    ('void_dragon', 'a void dragon with star-filled wings'),
    ('throne_sentinel', 'a golden throne sentinel with a long spear at the edge of the cosmos'),
    ('eternal_mage', 'an ageless mage with a robe of starlight and a crystal orb'),
    ('fallen_god_aethon', 'a fallen god on a throne of stars'),
    ('god_shadow', 'a huge dark shadow of a god with burning white eyes'),
    ('avatar_of_ruin', 'an avatar of ruin, a burning colossus of cracked stone'),
    ('varn_undying', 'a giant undying undead warlord in a cold blue crown'),
    ('crownless_tyrant', 'a colossal tyrant king with a broken crown'),
    ('azhorath_ascended', 'a huge ascended golden dragon king'),
    ('kaedrin_timeless', 'a titanic timeless lord among floating clocks'),
    ('fallen_archangels', 'dark winged archangels falling from a shattered golden sky'),
    ('aethon_unbound', 'a titanic god of void and stars breaking his chains'),
]

SKILL_LOOK = {
    'strike': 'a powerful weapon or spell attack in motion',
    'heal': 'soft golden healing light in open hands',
    'empower': 'a rising surge of power and a bright aura',
    'guard': 'a braced defensive stance behind a shield',
    'barrier': 'a glowing protective dome of energy',
    'dot': 'lingering poison, bleeding or flames eating at a target',
    'hot': 'gentle renewing light with leaves and sparkles',
    'weaken': 'a cursed mark draining a cracked enemy weapon',
    'stun': 'a crushing impact with dazed stars',
    None: 'a calm emblem of inner strength and mastery',
}


# Gemini writes a skill's name under its icon whenever the name is in the
# prompt, so each skill is described by objects taken from the words of its
# name instead (DD-336).
SKILL_MOTIF = {
    'aegis': 'a radiant round shield', 'aim': 'a crosshair over a target', 'aimed': 'a crosshair over a target',
    'ancestor': 'ghostly ancestor faces', 'ancestors': 'ghostly ancestor faces', 'ancestral': 'ghostly ancestor faces',
    'ancients': 'ghostly ancestor faces', 'ancient': 'weathered ancient stone', 'angel': 'white angel wings',
    'anthem': 'a war horn', 'arcane': 'a swirling violet arcane vortex', 'armor': 'a breastplate',
    'arsenal': 'a rack of crossed weapons', 'astral': 'a starry astral figure', 'aura': 'a glowing aura',
    'avatar': 'a towering glowing figure', 'backstab': 'a dagger in a back', 'banner': 'a tattered war banner',
    'barrage': 'a hail of blows', 'barrier': 'a shimmering wall of energy', 'bastion': 'a fortified tower',
    'beacon': 'a shining beacon fire', 'bear': 'a roaring bear head', 'beast': 'a wolf companion',
    'bestial': 'a roaring beast head', 'blade': 'a gleaming blade', 'blades': 'spinning blades',
    'blood': 'dripping blood', 'bloodthirst': 'dripping blood', 'bloom': 'a blooming flower',
    'bolt': 'a bolt of energy', 'bolts': 'heavy crossbow bolts', 'bomb': 'a smoking bomb',
    'bond': 'two linked glowing hands', 'bramble': 'thorny brambles', 'breath': 'a last breath of mist',
    'bullseye': 'an arrow in a bullseye', 'bulwark': 'a tower shield', 'burst': 'a bursting orb',
    'camouflage': 'a figure fading into leaves', 'celestial': 'stars and planets', 'chain': 'glowing chains',
    'charge': 'a charging figure', 'chi': 'a swirling chi orb', 'chorus': 'singing ghostly voices',
    'citadel': 'a shining citadel', 'cleanse': 'pure flowing water', 'cloud': 'a green toxic cloud',
    'combo': 'a chain of fist strikes', 'conduit': 'a crackling lightning rod', 'counter': 'a parrying blade',
    'counterguard': 'a parrying blade', 'crusader': 'a crusader sword', 'curse': 'a cursed skull sigil',
    'cuts': 'many slash marks', 'cutting': 'a sharp slash',
    'dance': 'a whirling dancer with blades', 'death': 'a grinning skull', 'deadeye': 'a glowing eye',
    'deflect': 'an arrow bouncing off a shield', 'deflecting': 'a palm deflecting a blow',
    'detonation': 'an explosion', 'diamond': 'a diamond', 'dirge': 'a mournful harp',
    'discord': 'jagged broken notes', 'discordant': 'jagged broken notes', 'dissonant': 'jagged broken notes',
    'drum': 'a war drum', 'drums': 'war drums', 'duelist': 'a rapier', 'earth': 'cracked earth',
    'earthbind': 'roots binding feet', 'earthen': 'rock armour', 'echo': 'echoing rings',
    'eclipse': 'a solar eclipse', 'elemental': 'fire, ice and lightning orbs', 'evasion': 'a dodging blur',
    'execute': 'an executioner axe', 'execution': 'an executioner axe', 'executioner': 'an executioner axe',
    'exorcism': 'a fleeing ghost', 'explosive': 'an exploding trap', 'eye': 'a hunter eye',
    'faith': 'a holy symbol', 'feint': 'a fake-out blade', 'fireball': 'a fireball',
    'fist': 'a clenched fist', 'fists': 'flurrying fists', 'flame': 'flames', 'fortress': 'a stone fortress',
    'flurry': 'a flurry of strikes', 'footwork': 'quick footprints', 'foresight': 'an all-seeing eye',
    'freeze': 'frozen ice', 'frost': 'frost crystals', 'frozen': 'frozen ice', 'fur': 'thick fur',
    'fury': 'raging flames', 'grace': 'a falling feather of light', 'hammer': 'a war hammer',
    'hand': 'a steady hand', 'hands': 'glowing hands', 'heal': 'a healing glow', 'healing': 'a healing glow',
    'heaven': 'parting clouds with sunbeams', 'heavenly': 'parting clouds with sunbeams', 'hex': 'a cursed doll stuck with pins',
    'hide': 'a thick hide', 'holy': 'a holy sunburst', 'hope': 'a lantern of hope',
    'ice': 'ice shards', 'ignition': 'a spark igniting', 'illusion': 'mirror images', 'inferno': 'a fire storm',
    'iron': 'iron plates', 'judgment': 'a glowing gavel of light', 'justice': 'balanced scales',
    'kick': 'a flying kick', 'kindling': 'kindling sparks', 'lance': 'a spear of ice', 'lash': 'a thorn whip',
    'life': 'a sprouting seedling', 'lifebloom': 'a glowing flower', 'light': 'a beam of light',
    'lightning': 'a lightning bolt', 'lunar': 'a crescent moon', 'mana': 'a blue mana crystal',
    'mark': 'a red target mark', 'marked': 'a red target mark', 'maul': 'claws mauling',
    'meditation': 'a meditating monk', 'melody': 'flowing music notes', 'melodic': 'flowing music notes',
    'meteor': 'falling meteors', 'missile': 'magic missiles', 'moon': 'a full moon', 'moonfire': 'a moonbeam',
    'moonlit': 'moonlight', 'mountain': 'a mountain', 'nature': 'green leaves', 'nightfall': 'a dark night sky',
    'nova': 'a ring of frost', 'palm': 'an open palm strike', 'piercing': 'a piercing arrow',
    'pillar': 'a stone pillar', 'plaguebringer': 'a plague mask', 'poisoned': 'dripping poison',
    'prayer': 'praying hands', 'prison': 'a cage of glowing light bars', 'prophecy': 'a glowing scroll of fate',
    'purification': 'purifying light', 'purify': 'purifying light', 'purging': 'purifying light',
    'rage': 'a raging face', 'rain': 'a rain of bolts', 'reload': 'a crossbow being reloaded',
    'rend': 'claw slashes', 'rewind': 'an hourglass with sand flowing upward', 'riposte': 'a parry and thrust',
    'roar': 'a roaring mouth', 'roots': 'grasping roots', 'rune': 'a glowing carved stone sigil of interlocking lines', 'runic': 'a glowing knotwork crystal',
    'seal': 'a holy seal', 'shadow': 'dark shadows', 'shadowstep': 'a figure stepping out of shadow',
    'shield': 'a shield', 'shot': 'an arrow in flight', 'siege': 'a siege bolt', 'skin': 'stone skin',
    'slam': 'a shield slam', 'slash': 'a sword slash', 'smite': 'a holy hammer strike', 'smoke': 'smoke',
    'snare': 'a snare trap', 'song': 'a lute and notes', 'songweaving': 'woven music notes',
    'stab': 'a quick dagger stab', 'stance': 'a fighting stance', 'starfall': 'falling stars',
    'static': 'static sparks', 'steel': 'steel bolts', 'stone': 'solid stone', 'stoneheart': 'a stone heart',
    'storm': 'a storm cloud', 'strike': 'a weapon strike', 'symphony': 'an orchestra of notes',
    'taunt': 'a shouting face', 'tempest': 'a whirling tempest', 'templar': 'a templar cross shield',
    'temporal': 'a glowing hourglass', 'thorn': 'thorns', 'thorns': 'thorns', 'thunder': 'thunder', 'time': 'a golden hourglass',
    'totem': 'a carved totem pole', 'totemic': 'carved totems', 'toxic': 'green toxin', 'toxin': 'green toxin',
    'trap': 'a steel trap', 'tree': 'a great tree', 'trick': 'a thrown handful of sand',
    'triple': 'three strikes', 'twin': 'twin blades', 'veil': 'a dark veil', 'venom': 'a venom drop',
    'venomous': 'a venom-dripping blade', 'verse': 'a singing bard', 'vines': 'healing vines',
    'volley': 'a volley of arrows', 'wall': 'a shield wall', 'war': 'a war horn', 'ward': 'a protective ward',
    'wave': 'a wave of fire', 'weakness': 'a cracked weapon', 'whirling': 'a whirling blade',
    'wild': 'wild claws', 'wisdom': 'an open book of light', 'world': 'a world tree', 'worldroot': 'world roots',
    'wrath': 'a furious storm',
    'guard': 'a raised kite shield', 'hardened': 'a scarred iron gauntlet', 'unyielding': 'a planted greatsword',
    'last': 'a lone battered knight', 'finishing': 'a downward killing blow', 'retaliation': 'a spiked shield',
    'unbreakable': 'an unbroken glowing shield', 'commanding': 'a commander helmet', 'rallying': 'a raised sword',
    'tactical': 'a battle map with markers', 'command': 'a pointing gauntlet', 'challenge': 'a thrown gauntlet',
    'frontline': 'a line of shields', 'march': 'marching boots', 'conquest': 'a planted flag on a hill',
    'frenzy': 'a berserker with twin axes', 'pain': 'a bleeding wound turning to fire', 'reckless': 'a wild swing',
    'onslaught': 'a flurry of axes', 'weapon': 'swapping weapons', 'flow': 'flowing water ribbons',
    'sequence': 'a chain of glowing strikes', 'opportunist': 'a dagger in an open guard',
    'lethal': 'a poisoned dagger tip', 'opening': 'a crack of light in armour', 'ambush': 'eyes in the bushes',
    'silent': 'a hooded assassin', 'killer': 'a hooded assassin', 'predator': 'a stalking panther',
    'patience': 'a coiled snake', 'catalyst': 'a bubbling flask', 'rhythm': 'pulsing beat rings',
    'relentless': 'a running blur', 'infinite': 'an infinity loop of blades', 'misdirection': 'a cloak decoy',
    'untouchable': 'a ghostly dodging figure', 'disengage': 'a backwards leap', 'critical': 'a red critical spark',
    'focus': 'a focused eye', 'deadly': 'a deadly arrowhead', 'precision': 'a needle-thin arrow',
    'range': 'a distant target', 'mechanism': 'crossbow gears', 'ballistic': 'a crossbow firing',
    'survival': 'a campfire and a knife', 'instinct': 'a wolf eye', 'tracking': 'paw tracks',
    'hunter': 'a hunter hood', 'preparation': 'laid out arrows and traps', 'attack': 'a lunging hawk',
    'pack': 'a wolf pack', 'primal': 'primal claw marks', 'connection': 'a glowing bond thread',
    'prepared': 'hidden trap triggers', 'ground': 'marked ground', 'rapid': 'speed lines', 'reset': 'a reset trap spring',
    'killing': 'a field of spikes', 'spell': 'a spell circle', 'surge': 'a surge of energy',
    'accelerated': 'a spinning hourglass', 'casting': 'casting hands', 'combustion': 'an exploding fire orb',
    'permafrost': 'frozen ground', 'zero': 'a snowflake', 'inner': 'a glowing inner core', 'balanced': 'a yin-yang',
    'spirit': 'a glowing spirit wisp', 'unshaken': 'a calm stone monk', 'adamant': 'an adamant body of stone',
    'body': 'a muscular body', 'balance': 'a yin-yang', 'step': 'a gliding step', 'harmony': 'harmonious rings',
    'flowing': 'flowing ribbons', 'transcendent': 'a rising lotus', 'ascension': 'a figure rising into light',
    'martial': 'a martial arts pose', 'combination': 'linked strike trails', 'finisher': 'a final palm strike',
    'endless': 'an endless loop', 'technique': 'a perfect form silhouette', 'blessing': 'a blessing hand',
    'divine': 'a divine halo', 'favor': 'a golden coin of light', 'renew': 'a renewing sprout of light',
    'intervention': 'a hand from the heavens', 'mercy': 'a dove', 'miracle': 'a miracle burst of light',
    'reinforced': 'riveted plates', 'insight': 'a third eye', 'sanctified': 'a glowing holy circle',
    'revelation': 'an opening sky', 'strength': 'a flexed arm', 'inspiration': 'a rising golden spark of faith',
    'enduring': 'an eternal flame', 'banish': 'a portal swallowing a demon', 'uncorrupted': 'a pure white lily',
    'consecration': 'consecrated glowing ground', 'sacred': 'a sacred chalice', 'defense': 'a raised shield',
    'block': 'a blocking shield', 'protection': 'a protective dome', 'guardian': 'a guardian statue',
    'recovery': 'a rising heart', 'radiant': 'a radiant sun', 'zeal': 'a burning heart', 'assault': 'a charging knight',
    'condemn': 'a pointing judge hand', 'sentence': 'a falling verdict sword', 'rejuvenation': 'a glowing sprout',
    'regrowth': 'regrowing leaves', 'network': 'a web of roots', 'natural': 'a leaf circle', 'evergreen': 'a pine tree',
    'form': 'a shapeshifting silhouette', 'endurance': 'a tough tortoise shell', 'savage': 'savage fangs',
    'quick': 'swift wings', 'harmonic': 'tuning fork waves', 'resonance': 'resonating sound waves',
    'pitch': 'a tuning fork', 'performance': 'a bard bowing on a candlelit stage', 'marching': 'a marching drum',
    'warrior': 'a warrior helmet', 'link': 'linked spirit chains', 'mend': 'stitched glowing light',
    'power': 'a power fist', 'presence': 'a towering silhouette', 'cry': 'a battle cry',
    'leader': 'a leader banner', 'shift': 'swapping weapons', 'motion': 'motion trails', 'tempo': 'a metronome',
}
CLASS_MOTIF = {
    'warrior': 'a sword', 'rogue': 'a dagger', 'ranger': 'a bow', 'mage': 'a staff with an orb', 'monk': 'a fist',
    'cleric': 'a holy symbol', 'paladin': 'a hammer', 'druid': 'a leaf', 'bard': 'a lute', 'shaman': 'a totem',
}


def skill_motif(cls, name):
    found = []
    for w in re.findall(r'[a-z]+', name.lower()):
        m = SKILL_MOTIF.get(w)
        if m and m not in found:
            found.append(m)
    return ' with '.join(found[:2]) or CLASS_MOTIF[cls]


def skills():
    """Class skill icons in class order, so a sheet keeps one palette."""
    import yaml
    with io.open(os.path.join(ROOT, 'content', 'skills.yaml'), encoding='utf-8') as f:
        sk = yaml.safe_load(f)['skills']
    with io.open(os.path.join(ROOT, 'content', 'i18n', 'en.yaml'), encoding='utf-8') as f:
        en = yaml.safe_load(f)
    en = en.get('strings', en)
    out = []
    for cls in PAL:
        for s in sk:
            if s['class'] != cls:
                continue
            epic = ', epic and radiant' if s.get('category') in ('ultimate', 'awakening') else ''
            out.append((s['code'], '%s, as %s, in %s%s' % (
                skill_motif(cls, en['skill.' + s['code']]), SKILL_LOOK[s.get('effect')], PAL[cls], epic)))
    return out


def nodes():
    return [(code, '%s, in %s' % (d, PAL[root])) for root, ns in NODES for code, d in ns] + [
        (code, '%s, in %s' % (d, PAL[root])) for root, code, d in NODE_FIXES]


# Gathering professions (DD-334): the Gather tab's kind tiles, beside the
# crafting emblems.
GATHER_PROFESSIONS = [
    ('mining', 'a steel pickaxe over a chunk of glittering ore'),
    ('herbalism', 'a curved sickle and a bundle of fresh green herbs with small flowers'),
    ('woodcutting', 'a woodcutter axe sunk into a cut log with wood chips'),
    ('fishing', 'a fishing rod with a silver fish leaping on the line over blue water'),
    ('hunting', 'a hunting bow and arrow above a paw print in the grass'),
]

# Cooking (DD-268, DD-274): plated fish dishes, humble at level 1 and richer
# up the ladder, ending on gold-trimmed feasts with a faint glow.
FOODS = [
    ('food_grilled_guppy', 'a tiny grilled fish on a plain wooden plank with a pinch of herbs'),
    ('food_sardine_skewer', 'three grilled sardines on a wooden skewer'),
    ('food_minnow_soup', 'a small clay bowl of clear soup with little fish and greens'),
    ('food_pike_fillet', 'a seared pike fillet on a wooden board with a lemon wedge'),
    ('food_lavender_carp', 'a whole roasted carp garnished with purple lavender sprigs'),
    ('food_grilled_silver_herring', 'two grilled silvery herrings with char marks on a tin plate'),
    ('food_herb_trout', 'a whole trout stuffed with green herbs on a wooden plate'),
    ('food_smoked_salmon', 'thick slices of pink smoked salmon on a slate board'),
    ('food_bass_stew', 'a hearty iron pot of bass stew with carrots and steam'),
    ('food_grilled_lake_perch', 'a grilled striped perch on a bed of reeds'),
    ('food_grilled_red_snapper', 'a grilled red snapper with chili slices on a ceramic plate'),
    ('food_grilled_river_catfish', 'a grilled whiskered catfish on a wooden platter'),
    ('food_tuna_steak', 'a thick seared tuna steak with pepper crust on a plate'),
    ('food_grilled_rainbow_trout', 'a grilled rainbow trout with shimmering pink and green skin'),
    ('food_royal_carp_feast', 'a golden carp on a silver platter with fruit and herbs'),
    ('food_grilled_emerald_pike', 'a grilled pike with glossy emerald-green scales on a silver plate'),
    ('food_grilled_silver_barracuda', 'a long grilled silver barracuda on a fine platter'),
    ('food_sturgeon_pie', 'a golden-crusted fish pie with a fish-shaped pastry top'),
    ('food_swordfish_grill', 'a swordfish steak with its long bill beside it on a grill plate'),
    ('food_grilled_stormscale_eel', 'a grilled dark blue eel crackling with tiny sparks on a platter'),
    ('food_frostfin_chowder', 'a bowl of creamy chowder with icy blue fish fins and frost on the rim'),
    ('food_grilled_obsidian_tuna', 'a black glossy tuna steak with glowing ember cracks on a dark plate'),
    ('food_marlin_platter', 'a grand marlin platter with sliced fish, lemons and greens'),
    ('food_grilled_golden_carp', 'a gleaming golden grilled carp on an ornate gold plate'),
    ('food_dragonblood_serpentfish', 'a crimson serpent-shaped fish in a dark red sauce on a gold platter'),
    ('food_grilled_void_angler', 'a grilled dark purple anglerfish with a softly glowing lure on a black plate'),
    ('food_moonfish_delicacy', 'a pale silver moonfish delicacy with a faint moonlight glow on fine porcelain'),
    ('food_phoenix_koi_feast', 'an orange-gold koi feast wreathed in gentle flames on a gold platter'),
    ('food_grilled_abyssal_dragonfish', 'a fearsome grilled deep-sea dragonfish with glowing teal spots on a dark platter'),
    ('food_leviathan_banquet', 'a lavish banquet platter of giant sea-beast cuts, gold trim and a radiant glow'),
]

# Potions (DD-241, DD-288): the missing rungs of the flask ladder; the same
# gold-framed flasks as the existing ones, red for health, blue for mana.
POTIONS = [
    ('potion_hp_supreme', 'the largest and most ornate crimson crystal flask with gold filigree wings and a radiant red glow'),
    ('potion_mp_medium', 'a round blue flask with a gold neck band and a cork'),
    ('potion_mp_large', 'a larger blue flask in a gold cage frame with a cork'),
    ('potion_mp_xlarge', 'a tall faceted blue crystal flask with a spiky gold frame'),
    ('potion_mp_huge', 'a big faceted blue crystal flask with many gold spikes and swirling mana inside'),
    ('potion_mp_supreme', 'the largest and most ornate sapphire crystal flask with gold filigree wings and a radiant blue glow'),
]

MEDALLION = ('Each cell is one round medallion that fills its square, with an ornate gold rim, on a near-black '
             'background. Stylized hand-painted high-fantasy MMORPG icon art, bold readable shapes that stay clear '
             'when tiny, vivid colours, dramatic glow and rim lighting.')
NO_TEXT = ('IMPORTANT: no captions, no labels, no names, no text or letters of any kind anywhere in the image. '
           'No logos, no watermark, no outer frame.')

SETS = {
    'nodes': dict(dir='class_nodes', mode='medallion', size=160, entries=nodes,
                  head='fantasy class emblem icons. ' + MEDALLION + ' Each emblem shows one symbol:'),
    'regions': dict(dir='regions', mode='medallion', size=160, entries=lambda: REGIONS,
                    head='fantasy map location icons. ' + MEDALLION + ' Inside each medallion is a small painted '
                         'landscape vignette of the place:'),
    'professions': dict(dir='professions', mode='medallion', size=160, entries=lambda: PROFESSIONS,
                        head='fantasy crafting profession emblem icons. ' + MEDALLION + ' Each emblem shows the '
                             'tools of one craft, warm workshop light:'),
    'helpers': dict(dir='helpers', mode='portrait', size=192, entries=lambda: HELPERS,
                    head='fantasy character portraits of friendly hired workers. Each portrait is head and '
                         'shoulders, face centered, looking at the viewer, on a dark simple vignette background. '
                         'Stylized hand-painted high-fantasy MMORPG art, bold readable shapes, warm lighting, '
                         'every person clearly different in age, race, face and clothes:'),
    'dungeons': dict(dir='dungeons', mode='medallion', size=160, entries=lambda: DUNGEONS,
                     head='fantasy dungeon boss icons. ' + MEDALLION + ' Inside each medallion is a dramatic '
                          'painted scene of one boss in its lair:'),
    'quests': dict(dir='quests', mode='medallion', size=160, entries=lambda: QUESTS,
                   head='fantasy quest icons. ' + MEDALLION + ' Inside each medallion is a small painted scene of '
                        'what the quest is about:'),
    'mobs': dict(dir='mobs', mode='medallion', size=160, entries=lambda: MOBS,
                 head='fantasy monster bestiary icons. ' + MEDALLION + ' Inside each medallion is a painted portrait '
                      'of one monster, head and upper body facing the viewer, in its home terrain:'),
    'skills': dict(dir='skills', mode='item', size=128, grid=5, entries=skills,
                   head='fantasy RPG class ability icons. Each cell is one square ability icon that fills its cell '
                        'edge to edge with a thin bevelled dark frame: one dramatic painted symbol of the skill, '
                        'centered, on a dark background tinted with the class colours. Stylized hand-painted '
                        'high-fantasy MMORPG spell icon art, bold readable shapes that stay clear when tiny, vivid '
                        'glow, every icon clearly different. The icons are pure pictures with nothing written on '
                        'them and no caption bars under them; sigils are abstract shapes, never letters or '
                        'numerals:'),
    'scrolls': dict(dir='items', mode='item', size=128, grid=5, entries=lambda: SCROLLS,
                    head='fantasy enchanting scroll icons. Each cell shows exactly one open scroll, whole and '
                         'centered, at a slight diagonal, on the same plain dark slate-blue backdrop with a soft '
                         'vignette. Stylized hand-painted game icon art like polished mobile MMORPG inventory icons: '
                         'clean bold silhouettes, soft cel shading, crisp rim light. The symbol is the main feature: '
                         'big, bold, high contrast, filling most of the scroll face, readable even when the icon is '
                         'tiny; the scroll face shows only that symbol, never any writing:'),
    'gather_professions': dict(dir='professions', mode='medallion', size=160, entries=lambda: GATHER_PROFESSIONS,
                               head='fantasy gathering profession emblem icons. ' + MEDALLION + ' Each emblem shows '
                                    'the tools of one gathering trade, warm outdoor light:'),
    'foods': dict(dir='items', mode='item', size=128, grid=5, entries=lambda: FOODS,
                  head='fantasy RPG cooked food icons. Each cell shows exactly one dish, whole and centered, seen '
                       'from a slight high angle, on the same plain dark slate-blue backdrop with a soft vignette. '
                       'Stylized hand-painted game icon art like polished mobile MMORPG inventory icons: clean bold '
                       'silhouettes, soft cel shading, crisp rim light, appetizing warm colours. Every dish looks '
                       'unique:'),
    'potions': dict(dir='items', mode='item', size=128, fit=True, entries=lambda: POTIONS,
                    head='fantasy RPG potion icons. Each cell shows exactly one corked potion flask, whole and '
                         'centered, on the same plain pure black background. Stylized hand-painted game icon art '
                         'like polished mobile MMORPG inventory icons: glossy glass, gold metal frames, glowing liquid, '
                         'bold silhouettes, crisp rim light:'),
    'class_portraits': dict(dir='classes', mode='single', size=(480, 640), entries=lambda: CLASS_PORTRAITS),
    'materials': dict(dir='items', mode='item', size=128, entries=lambda: MATERIALS,
                      head='fantasy RPG crafting material icons. Each cell shows exactly one item, whole and '
                           'centered, at a slight diagonal, on the same plain dark slate-blue backdrop with a soft '
                           'vignette. Stylized hand-painted game icon art like polished mobile MMORPG inventory '
                           'icons: clean bold silhouettes, soft cel shading, crisp rim light. Every item looks '
                           'unique:'),
}


def grid(name):
    return SETS[name].get('grid', GRID)


def pos(g, k):
    r, c = divmod(k, g)
    return POS[k] if g == GRID else 'row %d column %d' % (r + 1, c + 1)


def sheets(name):
    e = SETS[name]['entries']()
    per = 1 if SETS[name]['mode'] == 'single' else grid(name) ** 2
    return [e[i:i + per] for i in range(0, len(e), per)]


def prompt(name, n):
    cells = sheets(name)[n]
    if SETS[name]['mode'] == 'single':
        return ('Generate a VERTICAL PORTRAIT image (9:16): a full-figure heroic character illustration of %s, '
                'standing in a dynamic pose, seen from head to knees, filling the frame, in front of a matching '
                'painted fantasy background. Highly detailed stylized hand-painted high-fantasy MMORPG key art, '
                'vivid colours, dramatic rim lighting, ornate gear. One single character only. %s'
                % (cells[0][1], NO_TEXT))
    g = grid(name)
    parts = ['%s: %s' % (pos(g, k), d) for k, (_, d) in enumerate(cells)]
    empty = '' if len(cells) == g * g else ' The remaining cells show only the empty background.'
    return ('Generate a square image: a %dx%d grid of %d separate square %s %s.%s Divided by thin solid black '
            'lines, cells of equal size. %s' % (g, g, g * g, SETS[name]['head'], '; '.join(parts), empty, NO_TEXT))


def trim(im, dark=14):
    """Drops near-black rows/columns at the edges (grid lines)."""
    gray = im.convert('L')
    w, h = gray.size
    px = gray.load()
    col = lambda x: sum(px[x, y] for y in range(0, h, 4)) / len(range(0, h, 4))
    row = lambda y: sum(px[x, y] for x in range(0, w, 4)) / len(range(0, w, 4))
    l, r, t, b = 0, w - 1, 0, h - 1
    while l < w // 6 and col(l) < dark:
        l += 1
    while r > w * 5 // 6 and col(r) < dark:
        r -= 1
    while t < h // 6 and row(t) < dark:
        t += 1
    while b > h * 5 // 6 and row(b) < dark:
        b -= 1
    return im.crop((l, t, r + 1, b + 1))


def fit(im, light=200):
    """Whole object on a square: light grid lines are cleared from the edges,
    the art is cropped to its bounding box and padded on the backdrop."""
    m = max(4, im.width // 40)
    im = im.crop((m, m, im.width - m, im.height - m))
    gray = im.convert('L')
    box = gray.point(lambda v: 255 if 24 < v < light else 0).getbbox() or (0, 0, im.width, im.height)
    im = im.crop(box)
    side = int(max(im.size) * 1.08)
    out = Image.new('RGB', (side, side), (0, 0, 0))
    out.paste(im, ((side - im.width) // 2, (side - im.height) // 2))
    return out


def square(im):
    side = min(im.size)
    x0, y0 = (im.width - side) // 2, (im.height - side) // 2
    return im.crop((x0, y0, x0 + side, y0 + side))


def cut(name, numbers):
    s = SETS[name]
    src = os.path.join(ROOT, 'client', 'art_originals', name)
    dst = os.path.join(ROOT, 'client', 'assets', s['dir'])
    os.makedirs(dst, exist_ok=True)
    for n, cells in enumerate(sheets(name)):
        path = os.path.join(src, 'sheet%02d.jpg' % n)
        if (numbers and n not in numbers) or not os.path.exists(path):
            continue
        im = Image.open(path).convert('RGB')
        if s['mode'] == 'single':  # the upper 3:4 of the picture, head to knees
            w, h = im.size
            ch = min(h, w * 4 // 3)
            y0 = max(0, min(h - ch, (h - ch) // 4))
            im.crop((0, y0, w, y0 + ch)).resize(s['size'], Image.LANCZOS).save(
                os.path.join(dst, cells[0][0] + '.webp'), 'WEBP', quality=84, method=6)
            print(name, 'sheet%02d' % n, cells[0][0])
            continue
        gray = im.convert('L')
        w, h = im.size
        g = grid(name)
        xs = [0] + [p for a, b in gutters(gray, 0, g) for p in (a, b + 1)] + [w]
        ys = [0] + [p for a, b in gutters(gray, 1, g) for p in (a, b + 1)] + [h]
        inset = max(2, w // 200)
        for k, (code, _) in enumerate(cells):
            if code is None:
                continue
            r, c = divmod(k, g)
            cell = im.crop((xs[2 * c] + inset, ys[2 * r] + inset, xs[2 * c + 1] - inset, ys[2 * r + 1] - inset))
            if s['mode'] == 'medallion':
                cell = medallion(cell)
            else:
                cell = fit(cell) if s.get('fit') else square(trim(cell))
            cell.resize((s['size'], s['size']), Image.LANCZOS).save(
                os.path.join(dst, code + '.webp'), 'WEBP', quality=84, method=6)
        print(name, 'sheet%02d' % n, sum(1 for c, _ in cells if c))


if __name__ == '__main__':
    cmd = sys.argv[1] if len(sys.argv) > 1 else 'plan'
    if cmd == 'plan':
        for name in SETS:
            print(name, [len(x) for x in sheets(name)])
    elif cmd == 'prompt':
        print(prompt(sys.argv[2], int(sys.argv[3])))
    elif cmd == 'cut':
        cut(sys.argv[2], [int(x) for x in sys.argv[3:]])
