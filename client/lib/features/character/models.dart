/// Wire models for the character endpoints (server/internal/domain/character).
class CharacterSummary {
  const CharacterSummary({
    required this.id,
    required this.name,
    required this.raceCode,
    required this.classCode,
    required this.level,
    required this.state,
  });

  factory CharacterSummary.fromJson(Map<String, dynamic> j) => CharacterSummary(
    id: j['id'] as String,
    name: j['name'] as String,
    raceCode: j['race_code'] as String,
    classCode: j['class_code'] as String,
    level: (j['level'] as num).toInt(),
    state: j['state'] as String,
  );

  final String id;
  final String name;
  final String raceCode;
  final String classCode;
  final int level;
  final String state;

  bool get playable => state == 'active';
}

class CharacterList {
  const CharacterList(this.characters, this.slotCap);

  factory CharacterList.fromJson(Map<String, dynamic> j) => CharacterList([
    for (final c in (j['characters'] as List? ?? const [])) CharacterSummary.fromJson(c as Map<String, dynamic>),
  ], (j['slot_cap'] as num?)?.toInt() ?? 0);

  final List<CharacterSummary> characters;
  final int slotCap;

  bool get full => characters.length >= slotCap;
}

class CharacterSheet {
  const CharacterSheet({
    required this.id,
    required this.name,
    required this.raceCode,
    required this.classCode,
    required this.level,
    required this.xp,
    required this.xpToNext,
    required this.primary,
    required this.derived,
    required this.resources,
    required this.state,
    this.resourceKind = '',
    this.resourceMax = 0,
  });

  factory CharacterSheet.fromJson(Map<String, dynamic> j) {
    Map<String, int> ints(Object? m) => {
      for (final e in (m as Map<String, dynamic>? ?? const {}).entries)
        if (e.value is num) e.key: (e.value as num).toInt(),
    };
    return CharacterSheet(
      id: j['id'] as String,
      name: j['name'] as String,
      raceCode: j['race_code'] as String,
      classCode: j['class_code'] as String,
      level: (j['level'] as num).toInt(),
      xp: (j['xp'] as num).toInt(),
      xpToNext: (j['xp_to_next'] as num?)?.toInt() ?? 0,
      primary: ints(j['primary']),
      derived: ints(j['derived']),
      resources: ints(j['resources']),
      state: j['state'] as String? ?? 'active',
      resourceKind: j['resource_kind'] as String? ?? '',
      resourceMax: (j['resource_max'] as num?)?.toInt() ?? 0,
    );
  }

  final String id;
  final String name;
  final String raceCode;
  final String classCode;
  final int level;
  final int xp;
  final int xpToNext;

  /// Keys: str, agi, int, vit, cha, luck.
  final Map<String, int> primary;

  /// Keys: attack, defense, crit_pct, hit_pct.
  final Map<String, int> derived;

  /// Keys: hp_current/hp_max, mp_*, stamina_*, energy_*.
  final Map<String, int> resources;
  final String state;

  /// Class combat resource (DD-252): rage, energy, focus, mana, chi, ...
  /// Empty on servers without one. Combat fills it per encounter.
  final String resourceKind;
  final int resourceMax;

  /// 0..1 progress toward the next level (1 at level cap).
  double get xpProgress => xpToNext <= 0 ? 1 : (xp / xpToNext).clamp(0, 1).toDouble();
}

/// Money in copper → "12g 34s 56c" parts (1g = 100s = 10 000c).
({int gold, int silver, int copper}) splitCopper(int copper) {
  final c = copper < 0 ? 0 : copper;
  return (gold: c ~/ 10000, silver: (c % 10000) ~/ 100, copper: c % 100);
}

const races = ['human', 'elf', 'dwarf', 'orc', 'troll', 'undead', 'dark_elf', 'drakari', 'tauren', 'gnome'];

/// The ten base classes (DD-252), in the server's canonical order.
const classes = ['warrior', 'rogue', 'ranger', 'mage', 'monk', 'cleric', 'paladin', 'druid', 'bard', 'shaman'];

/// Mirrors the server rule so most mistakes are caught before a round trip.
final characterNamePattern = RegExp(r'^[A-Za-z][A-Za-z0-9_]{2,15}$');
