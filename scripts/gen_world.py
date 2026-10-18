"""World content past level 25 (DD-254, DD-255): zones, mobs and dungeons, with
names in en, tr, zh-CN and es. Mob stats follow the content's level curve
(hp 150+25L, attack 28+2.1L, defense 1+0.8L, xp 10·L^1.6) × archetype — the
same curve the level 1-25 mobs use, which the balance package confirms still
matches the class-average player with the gear ladder to level 100.

Run from the repository root: python scripts/gen_world.py
"""
import io
import os
import sys

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

ROOT = os.path.join(os.path.dirname(__file__), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

# archetype: hp, attack, defense multipliers; crit, evasion
ARCH = {
    'normal': (1.0, 1.0, 1.0, 5, 5),
    'tank': (1.75, 1.0, 1.9, 5, 2),
    'glass': (0.9, 1.5, 0.75, 11, 5),
    'evasive': (0.95, 1.15, 1.0, 9, 20),
    'crit': (1.05, 1.18, 1.0, 17, 10),
    'elite': (1.4, 1.15, 1.2, 9, 5),
    'chief': (1.2, 1.08, 1.1, 7, 5),  # dungeon bosses (the dungeon doubles HP, +20% attack)
    'overlord': (1.85, 1.3, 1.25, 10, 5),  # level 300+ dungeon bosses
}

ZONES = [
    dict(code='misty_marshes', level=25, travel=420, link='emberpeak',
         name=L('Misty Marshes', 'Sisli Bataklıklar', '迷雾沼泽', 'Marismas Brumosas'),
         mobs=[
             ('mist_leech', 30, 'evasive', L('Mist Leech', 'Sis Sülüğü', '雾水蛭', 'Sanguijuela de Niebla')),
             ('mud_golem', 35, 'tank', L('Mud Golem', 'Çamur Golemi', '泥魔像', 'Gólem de Barro')),
             ('bog_witch', 40, 'glass', L('Bog Witch', 'Bataklık Cadısı', '沼泽女巫', 'Bruja del Pantano')),
             ('marsh_queen', 42, 'chief', L('Marsh Queen', 'Bataklık Kraliçesi', '沼泽女王', 'Reina del Pantano')),
             ('venom_lizard', 45, 'crit', L('Venom Lizard', 'Zehirli Kertenkele', '毒蜥', 'Lagarto Venenoso')),
             ('drowned_knight', 50, 'elite', L('Drowned Knight', 'Batık Şövalye', '溺亡骑士', 'Caballero Ahogado')),
         ],
         dungeon=dict(code='sunken_temple', level=38, minutes=45, waves=['mist_leech', 'mud_golem', 'bog_witch'],
                      boss='marsh_queen', gold=8000, xp=12000,
                      name=L('Sunken Temple', 'Batık Tapınak', '沉没神殿', 'Templo Hundido'))),
    dict(code='iron_mountains', level=50, travel=480, link='misty_marshes',
         name=L('Iron Mountains', 'Demir Dağlar', '铁山', 'Montañas de Hierro'),
         mobs=[
             ('mountain_troll', 55, 'tank', L('Mountain Troll', 'Dağ Trolü', '山岭巨魔', 'Trol de Montaña')),
             ('dwarf_ghost', 60, 'evasive', L('Dwarf Ghost', 'Cüce Hayaleti', '矮人幽魂', 'Fantasma Enano')),
             ('rock_eagle', 65, 'glass', L('Rock Eagle', 'Kaya Kartalı', '岩鹰', 'Águila de Roca')),
             ('iron_king', 67, 'chief', L('Iron King', 'Demir Kral', '钢铁之王', 'Rey de Hierro')),
             ('iron_spider', 70, 'crit', L('Iron Spider', 'Demir Örümcek', '铁蛛', 'Araña de Hierro')),
             ('mine_warden', 75, 'elite', L('Mine Warden', 'Maden Bekçisi', '矿井看守', 'Guardián de la Mina')),
         ],
         dungeon=dict(code='forgotten_mine', level=63, minutes=50, waves=['mountain_troll', 'dwarf_ghost', 'rock_eagle'],
                      boss='iron_king', gold=15000, xp=20000,
                      name=L('Forgotten Mine', 'Unutulmuş Maden', '被遗忘的矿井', 'Mina Olvidada'))),
    dict(code='northern_keep_ruins', level=75, travel=540, link='iron_mountains',
         name=L('Northern Keep Ruins', 'Kuzey Kalesi Harabeleri', '北方要塞废墟', 'Ruinas de la Fortaleza del Norte'),
         mobs=[
             ('ghost_guard', 80, 'normal', L('Ghost Guard', 'Hayalet Muhafız', '幽灵卫兵', 'Guardia Fantasma')),
             ('cursed_archer', 85, 'glass', L('Cursed Archer', 'Lanetli Okçu', '诅咒弓手', 'Arquero Maldito')),
             ('bone_rider', 90, 'crit', L('Bone Rider', 'Kemik Atlı', '骸骨骑手', 'Jinete de Hueso')),
             ('fallen_lord_varn', 92, 'chief', L('Fallen Lord Varn', 'Düşmüş Lord Varn', '堕落领主瓦恩', 'Señor Caído Varn')),
             ('keep_gargoyle', 95, 'tank', L('Keep Gargoyle', 'Kale Gargoyli', '要塞石像鬼', 'Gárgola de la Fortaleza')),
             ('fallen_commander', 100, 'elite', L('Fallen Commander', 'Düşmüş Komutan', '堕落指挥官', 'Comandante Caído')),
         ],
         dungeon=dict(code='northern_keep', level=88, minutes=60, waves=['ghost_guard', 'cursed_archer', 'bone_rider'],
                      boss='fallen_lord_varn', gold=25000, xp=30000,
                      name=L('Northern Keep', 'Kuzey Kalesi', '北方要塞', 'Fortaleza del Norte'))),
]


# Continents beyond Sundermere (root regions; zones hang under them).
CONTINENTS = [
    dict(code='broken_crown', level=100,
         name=L('Kingdom of the Broken Crown', 'Kırık Taç Krallığı', '断冠王国', 'Reino de la Corona Rota')),
    dict(code='dragonspine_isles', level=300,
         name=L('Dragonspine Isles', 'Ejdersırtı Adaları', '龙脊群岛', 'Islas Espinadragón')),
    dict(code='starless_lands', level=600,
         name=L('Starless Lands', 'Yıldızsız Topraklar', '无星之地', 'Tierras sin Estrellas')),
    dict(code='sky_rift', level=700,
         name=L('Sky Rift', 'Gök Yarığı', '天空裂隙', 'Grieta Celeste')),
    dict(code='void_shore', level=850,
         name=L('Shore of the Void', 'Boşluğun Kıyısı', '虚空之岸', 'Orilla del Vacío')),
]

ZONES += [
    dict(code='crown_plains', parent='broken_crown', level=100, travel=300, link='northern_keep_ruins',
         name=L('Crown Plains', 'Taç Ovaları', '王冠平原', 'Llanuras de la Corona'),
         mobs=[
             ('rebel_rider', 110, 'normal', L('Rebel Rider', 'Asi Süvari', '叛军骑兵', 'Jinete Rebelde')),
             ('plains_lion', 120, 'crit', L('Plains Lion', 'Ova Aslanı', '平原狮', 'León de las Llanuras')),
             ('raider_mage', 130, 'evasive', L('Raider Mage', 'Yağmacı Büyücü', '掠夺法师', 'Mago Saqueador')),
             ('rebel_king_orrin', 132, 'elite', L('Rebel King Orrin', 'Asi Kral Orrin', '叛王奥林', 'Rey Rebelde Orrin')),
             ('war_bear', 140, 'tank', L('War Bear', 'Savaş Ayısı', '战熊', 'Oso de Guerra')),
             ('black_lancer', 150, 'elite', L('Black Lancer', 'Kara Mızrakçı', '黑枪骑士', 'Lancero Negro')),
         ],
         dungeon=dict(code='rebel_camp', level=128, minutes=60, waves=['rebel_rider', 'plains_lion', 'raider_mage'],
                      boss='rebel_king_orrin', gold=40000, xp=50000,
                      name=L('Rebel Camp', 'Asiler Kampı', '叛军营地', 'Campamento Rebelde'))),
    dict(code='burnt_forest', parent='broken_crown', level=150, travel=360, link='crown_plains',
         name=L('Burnt Forest', 'Yanık Orman', '焦土森林', 'Bosque Calcinado'),
         mobs=[
             ('ash_wolf', 160, 'evasive', L('Ash Wolf', 'Kül Kurdu', '灰烬狼', 'Lobo de Ceniza')),
             ('flame_spirit', 170, 'glass', L('Flame Spirit', 'Alev Ruhu', '火焰之灵', 'Espíritu de Llama')),
             ('cinder_treant', 180, 'tank', L('Cinder Treant', 'Kömür Treanti', '余烬树人', 'Ent de Ceniza')),
             ('burnt_mother_tree', 182, 'elite', L('Burnt Mother Tree', 'Yanık Ana Ağaç', '焦灼母树', 'Árbol Madre Calcinado')),
             ('crimson_harpy', 190, 'crit', L('Crimson Harpy', 'Kızıl Harpya', '赤红鹰身女妖', 'Arpía Carmesí')),
             ('scorched_shaman', 200, 'elite', L('Scorched Shaman', 'Yanık Şaman', '焦灼萨满', 'Chamán Abrasado')),
         ],
         dungeon=dict(code='fire_root', level=178, minutes=60, waves=['ash_wolf', 'flame_spirit', 'cinder_treant'],
                      boss='burnt_mother_tree', gold=60000, xp=80000,
                      name=L('Fire Root', 'Ateş Kökü', '火之根', 'Raíz de Fuego'))),
    dict(code='salt_desert', parent='broken_crown', level=200, travel=420, link='burnt_forest',
         name=L('Salt Desert', 'Tuz Çölü', '盐之沙漠', 'Desierto de Sal'),
         mobs=[
             ('sand_scorpion', 210, 'crit', L('Sand Scorpion', 'Kum Akrebi', '沙蝎', 'Escorpión de Arena')),
             ('desert_viper', 220, 'evasive', L('Desert Viper', 'Çöl Yılanı', '沙漠蝰蛇', 'Víbora del Desierto')),
             ('salt_elemental', 230, 'tank', L('Salt Elemental', 'Tuz Elementali', '盐元素', 'Elemental de Sal')),
             ('pharaoh_seth_ka', 232, 'elite', L('Pharaoh Seth-Ka', 'Firavun Seth-Ka', '法老塞特卡', 'Faraón Seth-Ka')),
             ('mummy_priest', 240, 'glass', L('Mummy Priest', 'Mumya Rahip', '木乃伊祭司', 'Sacerdote Momia')),
             ('sand_worm', 250, 'elite', L('Sand Worm', 'Kum Solucanı', '沙虫', 'Gusano de Arena')),
         ],
         dungeon=dict(code='buried_palace', level=228, minutes=60, waves=['sand_scorpion', 'desert_viper', 'salt_elemental'],
                      boss='pharaoh_seth_ka', gold=85000, xp=120000,
                      name=L('Buried Palace', 'Gömülü Saray', '掩埋宫殿', 'Palacio Enterrado'))),
    dict(code='crown_graveyard', parent='broken_crown', level=250, travel=480, link='salt_desert',
         name=L('Crown Graveyard', 'Taç Mezarlığı', '王冠墓地', 'Cementerio de la Corona'),
         mobs=[
             ('royal_skeleton', 260, 'normal', L('Royal Skeleton', 'Kraliyet İskeleti', '王室骷髅', 'Esqueleto Real')),
             ('tomb_keeper', 270, 'tank', L('Tomb Keeper', 'Mezar Bekçisi', '墓穴看守', 'Guardián de Tumbas')),
             ('soul_drinker', 280, 'glass', L('Soul Drinker', 'Ruh Emici', '噬魂者', 'Bebedor de Almas')),
             ('crypt_lord_malachar', 282, 'elite', L('Crypt Lord Malachar', 'Mezar Lordu Malachar', '墓穴领主玛拉卡', 'Señor de la Cripta Malachar')),
             ('black_knight', 290, 'crit', L('Black Knight', 'Kara Şövalye', '黑骑士', 'Caballero Negro')),
             ('crownless_king', 300, 'elite', L('Crownless King', 'Taçsız Kral', '无冕之王', 'Rey sin Corona')),
         ],
         dungeon=dict(code='royal_crypt', level=278, minutes=60, waves=['royal_skeleton', 'tomb_keeper', 'soul_drinker'],
                      boss='crypt_lord_malachar', gold=115000, xp=170000,
                      name=L('Royal Crypt', 'Kraliyet Kriptası', '王室地穴', 'Cripta Real'))),
]


ZONES += [
    dict(code='storm_coast', parent='dragonspine_isles', level=300, travel=300, link='crown_graveyard',
         name=L('Storm Coast', 'Fırtına Kıyısı', '风暴海岸', 'Costa de las Tormentas'),
         mobs=[
             ('sea_raider', 310, 'normal', L('Sea Raider', 'Deniz Haydudu', '海上掠夺者', 'Saqueador del Mar')),
             ('storm_bird', 320, 'evasive', L('Storm Bird', 'Fırtına Kuşu', '风暴鸟', 'Pájaro de Tormenta')),
             ('coral_golem', 330, 'tank', L('Coral Golem', 'Mercan Golemi', '珊瑚魔像', 'Gólem de Coral')),
             ('admiral_drask', 332, 'overlord', L('Admiral Drask', 'Amiral Drask', '德拉斯克上将', 'Almirante Drask')),
             ('siren', 340, 'glass', L('Siren', 'Siren', '塞壬', 'Sirena')),
             ('kraken_spawn', 350, 'elite', L('Kraken Spawn', 'Kraken Yavrusu', '海妖幼崽', 'Cría de Kraken')),
         ],
         dungeon=dict(code='sunken_fleet', level=328, minutes=60, waves=['sea_raider', 'storm_bird', 'coral_golem'],
                      boss='admiral_drask', gold=150000, xp=230000,
                      name=L('Sunken Fleet', 'Batık Filo', '沉没舰队', 'Flota Hundida'))),
    dict(code='dragon_slopes', parent='dragonspine_isles', level=350, travel=360, link='storm_coast',
         name=L('Dragon Slopes', 'Ejder Yamaçları', '龙之山坡', 'Laderas del Dragón'),
         mobs=[
             ('young_dragon', 360, 'crit', L('Young Dragon', 'Genç Ejder', '幼龙', 'Dragón Joven')),
             ('dragon_cultist', 370, 'glass', L('Dragon Cultist', 'Ejder Kültisti', '龙之信徒', 'Cultista del Dragón')),
             ('lava_lizard', 380, 'tank', L('Lava Lizard', 'Lav Kertenkelesi', '熔岩蜥蜴', 'Lagarto de Lava')),
             ('mother_dragon_vyrna', 382, 'overlord', L('Mother Dragon Vyrna', 'Ana Ejder Vyrna', '母龙薇尔娜', 'Dragona Madre Vyrna')),
             ('winged_wyvern', 390, 'evasive', L('Winged Wyvern', 'Kanatlı Wyvern', '飞翼双足龙', 'Guiverno Alado')),
             ('dragon_guard', 400, 'elite', L('Dragon Guard', 'Ejder Muhafızı', '龙之卫士', 'Guardia del Dragón')),
         ],
         dungeon=dict(code='brood_cavern', level=378, minutes=60, waves=['young_dragon', 'dragon_cultist', 'lava_lizard'],
                      boss='mother_dragon_vyrna', gold=190000, xp=320000,
                      name=L('Brood Cavern', 'Kuluçka Mağarası', '巢穴洞窟', 'Caverna de la Nidada'))),
    dict(code='mist_forest', parent='dragonspine_isles', level=400, travel=420, link='dragon_slopes',
         name=L('Mist Forest', 'Sis Ormanı', '迷雾森林', 'Bosque de la Niebla'),
         mobs=[
             ('fey_wolf', 410, 'evasive', L('Fey Wolf', 'Peri Kurdu', '妖精狼', 'Lobo Feérico')),
             ('mist_stalker', 420, 'crit', L('Mist Stalker', 'Sis Avcısı', '迷雾潜行者', 'Acechador de la Niebla')),
             ('ancient_treant', 430, 'tank', L('Ancient Treant', 'Kadim Treant', '远古树人', 'Ent Ancestral')),
             ('spirit_queen_sylwen', 432, 'overlord', L('Spirit Queen Sylwen', 'Ruh Kraliçesi Sylwen', '灵魂女王希尔温', 'Reina Espiritual Sylwen')),
             ('shadow_elf', 440, 'glass', L('Shadow Elf', 'Gölge Elfi', '暗影精灵', 'Elfo Sombrío')),
             ('forest_spirit', 450, 'elite', L('Forest Spirit', 'Orman Ruhu', '森林之灵', 'Espíritu del Bosque')),
         ],
         dungeon=dict(code='lost_grove', level=428, minutes=60, waves=['fey_wolf', 'mist_stalker', 'ancient_treant'],
                      boss='spirit_queen_sylwen', gold=240000, xp=430000,
                      name=L('Lost Grove', 'Kayıp Koru', '失落林地', 'Arboleda Perdida'))),
    dict(code='crystal_caves', parent='dragonspine_isles', level=450, travel=480, link='mist_forest',
         name=L('Crystal Caves', 'Kristal Mağaralar', '水晶洞穴', 'Cuevas de Cristal'),
         mobs=[
             ('crystal_spider', 460, 'crit', L('Crystal Spider', 'Kristal Örümcek', '水晶蜘蛛', 'Araña de Cristal')),
             ('stone_eater', 470, 'tank', L('Stone Eater', 'Taş Yiyen', '食石者', 'Devorador de Piedra')),
             ('light_elemental', 480, 'glass', L('Light Elemental', 'Işık Elementali', '光元素', 'Elemental de Luz')),
             ('crystal_titan', 482, 'overlord', L('Crystal Titan', 'Kristal Titan', '水晶泰坦', 'Titán de Cristal')),
             ('crystal_dragon', 490, 'evasive', L('Crystal Dragon', 'Kristal Ejder', '水晶龙', 'Dragón de Cristal')),
             ('cave_titan', 500, 'elite', L('Cave Titan', 'Mağara Titanı', '洞穴泰坦', 'Titán de las Cavernas')),
         ],
         dungeon=dict(code='crystal_heart', level=478, minutes=60, waves=['crystal_spider', 'stone_eater', 'light_elemental'],
                      boss='crystal_titan', gold=300000, xp=560000,
                      name=L('Crystal Heart', 'Kristal Kalp', '水晶之心', 'Corazón de Cristal'))),
    dict(code='volcanic_isle', parent='dragonspine_isles', level=500, travel=540, link='crystal_caves',
         name=L('Volcanic Isle', 'Volkanik Ada', '火山岛', 'Isla Volcánica'),
         mobs=[
             ('magma_golem', 510, 'tank', L('Magma Golem', 'Magma Golemi', '岩浆魔像', 'Gólem de Magma')),
             ('fire_giant', 520, 'normal', L('Fire Giant', 'Ateş Devi', '火巨人', 'Gigante de Fuego')),
             ('ash_dragon', 530, 'crit', L('Ash Dragon', 'Kül Ejderi', '灰烬龙', 'Dragón de Ceniza')),
             ('fire_lord_ignar', 532, 'overlord', L('Fire Lord Ignar', 'Ateş Lordu Ignar', '火焰领主伊格纳', 'Señor del Fuego Ignar')),
             ('volcano_shaman', 540, 'glass', L('Volcano Shaman', 'Volkan Şamanı', '火山萨满', 'Chamán del Volcán')),
             ('volcano_warden', 550, 'elite', L('Volcano Warden', 'Volkan Bekçisi', '火山守卫', 'Guardián del Volcán')),
         ],
         dungeon=dict(code='volcano_heart', level=528, minutes=60, waves=['magma_golem', 'fire_giant', 'ash_dragon'],
                      boss='fire_lord_ignar', gold=370000, xp=720000,
                      name=L('Volcano Heart', 'Volkan Kalbi', '火山之心', 'Corazón del Volcán'))),
    dict(code='throne_peaks', parent='dragonspine_isles', level=550, travel=600, link='volcanic_isle',
         name=L('Throne Peaks', 'Taht Zirveleri', '王座峰', 'Picos del Trono'),
         mobs=[
             ('dragon_knight', 560, 'normal', L('Dragon Knight', 'Ejder Şövalyesi', '龙骑士', 'Caballero Dragón')),
             ('ancient_wyvern', 570, 'evasive', L('Ancient Wyvern', 'Kadim Wyvern', '远古双足龙', 'Guiverno Ancestral')),
             ('dragon_sorcerer', 580, 'glass', L('Dragon Sorcerer', 'Ejder Büyücüsü', '龙之术士', 'Hechicero Dragón')),
             ('dragon_king_azhorath', 582, 'overlord', L('Dragon King Azhorath', 'Ejder Kral Azhorath', '龙王阿兹霍拉斯', 'Rey Dragón Azhorath')),
             ('throne_guardian', 590, 'tank', L('Throne Guardian', 'Taht Muhafızı', '王座守护者', 'Guardián del Trono')),
             ('elder_dragon', 600, 'elite', L('Elder Dragon', 'Yaşlı Ejder', '古龙', 'Dragón Anciano')),
         ],
         dungeon=dict(code='dragon_throne', level=578, minutes=60, waves=['dragon_knight', 'ancient_wyvern', 'dragon_sorcerer'],
                      boss='dragon_king_azhorath', gold=450000, xp=900000,
                      name=L('Dragon Throne', 'Ejder Tahtı', '龙之王座', 'Trono del Dragón'))),
]


ZONES += [
    dict(code='dark_plains', parent='starless_lands', level=600, travel=300, link='throne_peaks',
         name=L('Dark Plains', 'Karanlık Ova', '黑暗平原', 'Llanuras Oscuras'),
         mobs=[
             ('shadow_rider', 610, 'normal', L('Shadow Rider', 'Gölge Atlı', '暗影骑手', 'Jinete de Sombras')),
             ('void_wolf', 620, 'evasive', L('Void Wolf', 'Boşluk Kurdu', '虚空狼', 'Lobo del Vacío')),
             ('star_gnawer', 630, 'crit', L('Star Gnawer', 'Yıldız Kemiren', '噬星者', 'Roedor de Estrellas')),
             ('star_eater_omh', 632, 'overlord', L('Star Eater Omh', 'Yıldız Yiyen Omh', '吞星者奥姆', 'Devoraestrellas Omh')),
             ('dark_angel', 640, 'crit', L('Dark Angel', 'Kara Melek', '黑暗天使', 'Ángel Oscuro')),
             ('night_giant', 650, 'elite', L('Night Giant', 'Gece Devi', '暗夜巨人', 'Gigante de la Noche')),
         ],
         dungeon=dict(code='darkened_observatory', level=628, minutes=60, waves=['shadow_rider', 'void_wolf', 'star_gnawer'],
                      boss='star_eater_omh', gold=550000, xp=1150000,
                      name=L('Darkened Observatory', 'Söndürülmüş Gözlemevi', '熄灭的天文台', 'Observatorio Apagado'))),
    dict(code='lost_city', parent='starless_lands', level=650, travel=360, link='dark_plains',
         name=L('Lost City', 'Kayıp Şehir', '失落之城', 'Ciudad Perdida'),
         mobs=[
             ('ghost_noble', 660, 'normal', L('Ghost Noble', 'Hayalet Soylu', '幽灵贵族', 'Noble Fantasma')),
             ('time_thief', 670, 'evasive', L('Time Thief', 'Zaman Hırsızı', '时间窃贼', 'Ladrón del Tiempo')),
             ('cursed_sorcerer', 680, 'glass', L('Cursed Sorcerer', 'Lanetli Büyücü', '诅咒术士', 'Hechicero Maldito')),
             ('time_lord_kaedrin', 682, 'overlord', L('Time Lord Kaedrin', 'Zaman Lordu Kaedrin', '时间领主凯德林', 'Señor del Tiempo Kaedrin')),
             ('stone_gargoyle', 690, 'tank', L('Stone Gargoyle', 'Taş Gargoyl', '石像鬼', 'Gárgola de Piedra')),
             ('city_guardian', 700, 'elite', L('City Guardian', 'Şehir Koruyucusu', '城市守护者', 'Guardián de la Ciudad')),
         ],
         dungeon=dict(code='clock_tower', level=678, minutes=60, waves=['ghost_noble', 'time_thief', 'cursed_sorcerer'],
                      boss='time_lord_kaedrin', gold=650000, xp=1400000,
                      name=L('Clock Tower', 'Saat Kulesi', '钟楼', 'Torre del Reloj'))),
]


ZONES += [
    dict(code='floating_isles', parent='sky_rift', level=700, travel=300, link='lost_city',
         name=L('Floating Isles', 'Uçan Adalar', '浮空群岛', 'Islas Flotantes'),
         mobs=[
             ('cloud_giant', 710, 'tank', L('Cloud Giant', 'Bulut Devi', '云巨人', 'Gigante de las Nubes')),
             ('storm_dragon', 720, 'crit', L('Storm Dragon', 'Fırtına Ejderi', '风暴龙', 'Dragón de Tormenta')),
             ('sky_knight', 730, 'normal', L('Sky Knight', 'Gök Şövalyesi', '天空骑士', 'Caballero del Cielo')),
             ('storm_king_zephyros', 732, 'overlord', L('Storm King Zephyros', 'Fırtına Kralı Zephyros', '风暴之王泽菲罗斯', 'Rey de la Tormenta Zephyros')),
             ('wind_elemental', 740, 'evasive', L('Wind Elemental', 'Rüzgâr Elementali', '风元素', 'Elemental de Viento')),
             ('thunderbird', 750, 'elite', L('Thunderbird', 'Yıldırım Kuşu', '雷鸟', 'Ave del Trueno')),
         ],
         dungeon=dict(code='storm_fortress', level=728, minutes=60, waves=['cloud_giant', 'storm_dragon', 'sky_knight'],
                      boss='storm_king_zephyros', gold=780000, xp=1700000,
                      name=L('Storm Fortress', 'Fırtına Kalesi', '风暴要塞', 'Fortaleza de la Tormenta'))),
    dict(code='rift_edge', parent='sky_rift', level=750, travel=360, link='floating_isles',
         name=L('Rift Edge', 'Yarık Kenarı', '裂隙边缘', 'Borde de la Grieta'),
         mobs=[
             ('void_creature', 760, 'normal', L('Void Creature', 'Boşluk Yaratığı', '虚空生物', 'Criatura del Vacío')),
             ('rift_warden', 770, 'tank', L('Rift Warden', 'Yarık Bekçisi', '裂隙守卫', 'Celador de la Grieta')),
             ('chaos_elemental', 780, 'glass', L('Chaos Elemental', 'Kaos Elementali', '混沌元素', 'Elemental del Caos')),
             ('chaos_mother', 782, 'overlord', L('Chaos Mother', 'Kaos Ana', '混沌之母', 'Madre del Caos')),
             ('rift_dragon', 790, 'crit', L('Rift Dragon', 'Yarık Ejderi', '裂隙龙', 'Dragón de la Grieta')),
             ('chaos_lord', 800, 'elite', L('Chaos Lord', 'Kaos Lordu', '混沌领主', 'Señor del Caos')),
         ],
         dungeon=dict(code='rift_core', level=778, minutes=60, waves=['void_creature', 'rift_warden', 'chaos_elemental'],
                      boss='chaos_mother', gold=920000, xp=2100000,
                      name=L('Rift Core', 'Yarık Çekirdeği', '裂隙核心', 'Núcleo de la Grieta'))),
    dict(code='celestial_temple', parent='sky_rift', level=800, travel=420, link='rift_edge',
         name=L('Celestial Temple', 'Göksel Tapınak', '天界神殿', 'Templo Celestial'),
         mobs=[
             ('fallen_angel', 810, 'normal', L('Fallen Angel', 'Düşmüş Melek', '堕天使', 'Ángel Caído')),
             ('temple_guardian', 820, 'tank', L('Temple Guardian', 'Tapınak Muhafızı', '神殿守护者', 'Guardián del Templo')),
             ('celestial_dragon', 830, 'crit', L('Celestial Dragon', 'Göksel Ejder', '天界龙', 'Dragón Celestial')),
             ('archangel_uriel', 832, 'overlord', L('Archangel Uriel', 'Başmelek Uriel', '大天使乌列尔', 'Arcángel Uriel')),
             ('light_hunter', 840, 'evasive', L('Light Hunter', 'Işık Avcısı', '光之猎手', 'Cazador de Luz')),
             ('seraph', 850, 'elite', L('Seraph', 'Seraf', '炽天使', 'Serafín')),
         ],
         dungeon=dict(code='celestial_dome', level=828, minutes=60, waves=['fallen_angel', 'temple_guardian', 'celestial_dragon'],
                      boss='archangel_uriel', gold=1080000, xp=2600000,
                      name=L('Celestial Dome', 'Göksel Kubbe', '天界穹顶', 'Cúpula Celestial'))),
]


ZONES += [
    dict(code='ash_sea', parent='void_shore', level=850, travel=300, link='celestial_temple',
         name=L('Ash Sea', 'Kül Denizi', '灰烬之海', 'Mar de Ceniza'),
         mobs=[
             ('ash_leviathan', 860, 'tank', L('Ash Leviathan', 'Kül Leviathanı', '灰烬利维坦', 'Leviatán de Ceniza')),
             ('null_hunter', 870, 'evasive', L('Null Hunter', 'Hiçlik Avcısı', '虚无猎手', 'Cazador de la Nada')),
             ('undying_soldier', 880, 'normal', L('Undying Soldier', 'Ölümsüz Asker', '不死士兵', 'Soldado Inmortal')),
             ('leviathan_king', 882, 'overlord', L('Leviathan King', 'Leviathan Kral', '利维坦之王', 'Rey Leviatán')),
             ('ash_wyrm_lord', 890, 'crit', L('Ash Wyrm Lord', 'Kül Ejderi Lordu', '灰烬龙王', 'Señor Dragón de Ceniza')),
             ('eternal_sentinel', 900, 'elite', L('Eternal Sentinel', 'Ebedi Bekçi', '永恒哨兵', 'Centinela Eterno')),
         ],
         dungeon=dict(code='sunken_throne', level=878, minutes=60, waves=['ash_leviathan', 'null_hunter', 'undying_soldier'],
                      boss='leviathan_king', gold=1250000, xp=3100000,
                      name=L('Sunken Throne', 'Batık Taht', '沉没王座', 'Trono Hundido'))),
    dict(code='last_bastion', parent='void_shore', level=900, travel=360, link='ash_sea',
         name=L('Last Bastion Walls', 'Son Kale Surları', '最后堡垒城墙', 'Murallas del Último Bastión'),
         mobs=[
             ('doom_rider', 910, 'crit', L('Doom Rider', 'Kıyamet Atlısı', '末日骑手', 'Jinete del Fin')),
             ('void_titan', 920, 'tank', L('Void Titan', 'Boşluk Titanı', '虚空泰坦', 'Titán del Vacío')),
             ('soul_reaper', 930, 'crit', L('Soul Reaper', 'Ruh Biçen', '灵魂收割者', 'Segador de Almas')),
             ('doom_commander_morvath', 932, 'overlord', L('Doom Commander Morvath', 'Kıyamet Komutanı Morvath', '末日指挥官莫瓦斯', 'Comandante del Fin Morvath')),
             ('ancient_golem', 940, 'normal', L('Ancient Golem', 'Kadim Golem', '远古魔像', 'Gólem Ancestral')),
             ('last_commander', 950, 'elite', L('Last Commander', 'Son Komutan', '最后的指挥官', 'Último Comandante')),
         ],
         dungeon=dict(code='bastion_keep', level=928, minutes=60, waves=['doom_rider', 'void_titan', 'soul_reaper'],
                      boss='doom_commander_morvath', gold=1450000, xp=3700000,
                      name=L('Last Bastion', 'Son Kale', '最后堡垒', 'Último Bastión'))),
    dict(code='edge_of_eternity', parent='void_shore', level=950, travel=420, link='last_bastion',
         name=L('Edge of Eternity', 'Ebediyet Kıyısı', '永恒边缘', 'Borde de la Eternidad'),
         mobs=[
             ('void_dragon', 960, 'crit', L('Void Dragon', 'Hiçlik Ejderi', '虚无之龙', 'Dragón de la Nada')),
             ('throne_sentinel', 970, 'tank', L('Throne Sentinel', 'Taht Gözcüsü', '王座哨兵', 'Centinela del Trono')),
             ('eternal_mage', 980, 'glass', L('Eternal Mage', 'Ebedi Büyücü', '永恒法师', 'Mago Eterno')),
             ('fallen_god_aethon', 982, 'overlord', L('Fallen God Aethon', 'Düşmüş Tanrı Aethon', '堕落之神埃松', 'Dios Caído Aethon')),
             ('god_shadow', 990, 'evasive', L('Shadow of a God', 'Tanrı Gölgesi', '神之阴影', 'Sombra de un Dios')),
             ('avatar_of_ruin', 1000, 'elite', L('Avatar of Ruin', 'Yıkım Avatarı', '毁灭化身', 'Avatar de la Ruina')),
         ],
         dungeon=dict(code='eternal_throne', level=978, minutes=60, waves=['void_dragon', 'throne_sentinel', 'eternal_mage'],
                      boss='fallen_god_aethon', gold=1700000, xp=4400000,
                      name=L('Eternal Throne', 'Ebedi Taht', '永恒王座', 'Trono Eterno'))),
]


def mob_yaml(code, level, arch):
    hm, am, dm, crit, eva = ARCH[arch]
    gold = round(503 * (level / 25) ** 1.5)
    return '\n'.join([
        '  - code: ' + code, '    display_name_key: mob.' + code, '    level: %d' % level,
        '    hp_max: %d' % round((150 + 25 * level) * hm), '    attack: %d' % round((28 + 2.1 * level) * am),
        '    defense: %d' % round((1 + 0.8 * level) * dm), '    crit_pct: %d' % crit,
        '    hit_pct: %d' % min(92, 82 + level // 5), '    evasion_pct: %d' % eva,
        '    xp_reward: %d' % round(10 * level ** 1.6), '    gold_copper_min: %d' % gold,
        '    gold_copper_max: %d' % round(gold * 2.5)])


mobs, dungeons, regions, links = [], [], [], []
names = {k: {} for k in ('en', 'tr', 'zh', 'es')}


def put(key, n):
    for lang in names:
        names[lang][key] = n[lang]


for c in CONTINENTS:
    put('region.' + c['code'], c['name'])
    regions.append('\n'.join(['  - code: ' + c['code'], '    display_name_key: region.' + c['code'], '    kind: continent',
                              '    min_level: %d' % c['level'], '    travel_seconds_from_parent: 0', '    safe: true']))
for z in ZONES:
    put('region.' + z['code'], z['name'])
    regions.append('\n'.join(['  - code: ' + z['code'], '    display_name_key: region.' + z['code'], '    kind: zone',
                              '    parent_code: ' + z.get('parent', 'sundermere'), '    min_level: %d' % z['level'],
                              '    travel_seconds_from_parent: %d' % z['travel'], '    safe: false']))
    a, b = sorted([z['code'], z['link']])
    links.append('  - { region_a_code: %s, region_b_code: %s, travel_seconds: 300 }' % (a, b))
    for code, level, arch, n in z['mobs']:
        put('mob.' + code, n)
        mobs.append(mob_yaml(code, level, arch))
    d = z['dungeon']
    put('dungeon.' + d['code'], d['name'])
    put('region.' + d['code'], d['name'])
    regions.append('\n'.join(['  - code: ' + d['code'], '    display_name_key: region.' + d['code'], '    kind: dungeon',
                              '    parent_code: ' + z['code'], '    min_level: %d' % d['level'],
                              '    travel_seconds_from_parent: 300', '    safe: false']))
    dungeons.append('\n'.join(['  - code: ' + d['code'], '    display_name_key: dungeon.' + d['code'],
                               '    min_level: %d' % d['level'],
                               '    duration_seconds: %d         # %d minutes' % (d['minutes'] * 60, d['minutes']),
                               '    waves:',
                               '      - { mob_code: %s, count: 3 }' % d['waves'][0],
                               '      - { mob_code: %s, count: 3 }' % d['waves'][1],
                               '      - { mob_code: %s, count: 2 }' % d['waves'][2],
                               '    boss_mob_code: ' + d['boss'], '    bonus_gold_copper: %d' % d['gold'],
                               '    bonus_xp: %d' % d['xp']]))


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        s = f.read()
    return s.replace('\r\n', '\n'), '\r\n' in s


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


def block(s, name, body, before=None):
    start, end = '  # >>> generated: %s (scripts/gen_world.py)\n' % name, '  # <<< generated: %s\n' % name
    b = start + body.rstrip('\n') + '\n' + end
    if start in s:
        return s[:s.index(start)] + b + s[s.index(end) + len(end):]
    if before:
        assert s.count(before) == 1, before
        return s.replace(before, b + '\n' + before)
    return s.rstrip('\n') + '\n\n' + b


for fname, name, body, before in [
    ('mobs.yaml', 'mobs', '\n\n'.join(mobs), None),
    ('dungeons.yaml', 'dungeons', '\n\n'.join(dungeons), None),
]:
    p = os.path.join(ROOT, fname)
    s, crlf = read(p)
    write(p, block(s, name, body, before), crlf)
p = os.path.join(ROOT, 'regions.yaml')
s, crlf = read(p)
s = block(s, 'regions', '\n\n'.join(regions), before='# Non-parent shortcuts.')
s = block(s, 'links', '\n'.join(links))
write(p, s, crlf)

merge(names)
print('mobs', len(mobs), 'dungeons', len(dungeons), 'regions', len(regions))
