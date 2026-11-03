import 'dart:async';

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/motion.dart';
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../adventure/adventure_models.dart';
import '../adventure/adventure_repository.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_images.dart';
import '../inventory/item_tags.dart';
import '../inventory/item_widgets.dart';
import 'pet_species_sheet.dart';

/// One pet species (DD-328), best first in the catalog.
class PetSpecies {
  const PetSpecies({
    required this.code,
    required this.rank,
    required this.minMobLevel,
    required this.eggPct,
    required this.hatchSeconds,
    required this.bonus,
    this.bonusMin,
    this.assistMin = 0,
    this.assistMax = 0,
    this.element = '',
  });

  factory PetSpecies.fromJson(Map<String, dynamic> j) => PetSpecies(
    code: j['code'] as String,
    rank: (j['rank'] as num).toInt(),
    minMobLevel: (j['min_mob_level'] as num).toInt(),
    eggPct: (j['egg_pct'] as num).toDouble(),
    hatchSeconds: (j['hatch_seconds'] as num).toInt(),
    bonus: _bonus(j['bonus'] as Map<String, dynamic>),
    bonusMin: j['bonus_min'] == null ? null : _bonus(j['bonus_min'] as Map<String, dynamic>),
    assistMin: (j['assist_min'] as num?)?.toInt() ?? 0,
    assistMax: (j['assist_max'] as num?)?.toInt() ?? 0,
    element: j['element'] as String? ?? '',
  );

  final String code;
  final int rank, minMobLevel, hatchSeconds;
  final double eggPct;

  /// The bonus at level 100 ([bonus]) and at level 1 ([bonusMin]).
  final StatBonus bonus;
  final StatBonus? bonusMin;

  /// The damage share it adds to its owner's hits while active, at level 1
  /// and level 100 (DD-341).
  final int assistMin, assistMax;

  /// The pet index's element code (DD-342); named by element.<code>.
  final String element;

  /// The egg odds as shown: whole percents without decimals.
  String get eggPctText => eggPct == eggPct.roundToDouble() ? '${eggPct.toInt()}' : eggPct.toStringAsFixed(2);
}

StatBonus _bonus(Map<String, dynamic> b) =>
    StatBonus(b['kind'] as String, b['target'] as String, (b['amount'] as num).toDouble());

/// A pet's rank as an item rarity: its colour.
String petRarity(int rank) => switch (rank) {
  50 => 'mythic',
  >= 40 => 'legendary',
  >= 27 => 'epic',
  >= 14 => 'rare',
  _ => 'uncommon',
};

/// A pet's evolution stage by level (DD-334): 1–33, 34–66, 67–100.
int petStage(int level) => level <= 33 ? 1 : (level <= 66 ? 2 : 3);

/// The species' art at [stage] (`scripts/pet_assets.py`), framed on its own
/// backdrop; a paw icon while it is missing.
class PetArt extends StatelessWidget {
  const PetArt({super.key, required this.species, required this.stage, required this.size, required this.color});

  final String species;
  final int stage;
  final double size;
  final Color color;

  @override
  Widget build(BuildContext context) => Container(
    width: size,
    height: size,
    decoration: BoxDecoration(
      borderRadius: BorderRadius.circular(Radii.md),
      border: Border.all(color: color.withValues(alpha: 0.6)),
    ),
    clipBehavior: Clip.antiAlias,
    child: Image.asset(
      'assets/pets/${species}_$stage.webp',
      fit: BoxFit.cover,
      cacheWidth: (size * MediaQuery.devicePixelRatioOf(context)).round(),
      errorBuilder: (_, _, _) => Icon(Icons.pets, size: size * 0.5, color: color),
    ),
  );
}

class Pet {
  const Pet({
    required this.id,
    required this.species,
    required this.rank,
    required this.level,
    required this.xp,
    required this.xpToNext,
    required this.power,
    required this.hatchReadyAt,
    required this.active,
    required this.bonus,
    this.assistPct = 0,
    this.huntMob,
    this.huntStartedAt,
    this.huntReadyAt,
  });

  factory Pet.fromJson(Map<String, dynamic> j) {
    final hunt = j['hunt'] as Map<String, dynamic>?;
    return Pet(
      id: j['id'] as String,
      species: j['species'] as String,
      rank: (j['rank'] as num?)?.toInt() ?? 1,
      level: (j['level'] as num).toInt(),
      xp: (j['xp'] as num).toInt(),
      xpToNext: (j['xp_to_next'] as num).toInt(),
      power: (j['power'] as num).toInt(),
      hatchReadyAt: DateTime.parse(j['hatch_ready_at'] as String),
      active: j['active'] as bool? ?? false,
      bonus: _bonus(j['bonus'] as Map<String, dynamic>),
      assistPct: (j['assist_pct'] as num?)?.toInt() ?? 0,
      huntMob: hunt?['mob_code'] as String?,
      huntStartedAt: hunt == null ? null : DateTime.parse(hunt['started_at'] as String),
      huntReadyAt: hunt == null ? null : DateTime.parse(hunt['ready_at'] as String),
    );
  }

  final String id, species;
  final int rank, level, xp, xpToNext, power;
  final DateTime hatchReadyAt;
  final bool active;
  final StatBonus bonus;

  /// Extra damage share on its owner's hits while active (DD-341).
  final int assistPct;
  final String? huntMob;
  final DateTime? huntStartedAt, huntReadyAt;

  bool hatched(DateTime now) => !now.isBefore(hatchReadyAt);
}

class PetEgg {
  const PetEgg({required this.itemId, required this.template, required this.species, required this.quantity});

  factory PetEgg.fromJson(Map<String, dynamic> j) => PetEgg(
    itemId: j['item_id'] as String,
    template: j['template_code'] as String,
    species: j['species'] as String,
    quantity: (j['quantity'] as num).toInt(),
  );

  final String itemId, template, species;
  final int quantity;
}

/// What a collected pet hunt brought.
class PetHuntResult {
  const PetHuntResult({required this.wins, required this.xp, required this.died, required this.pet});

  final int wins, xp;
  final bool died;
  final Pet pet;
}

/// The character's pets and eggs, the bag sizes and the species it has
/// discovered (DD-342).
typedef PetList = ({List<Pet> pets, List<PetEgg> eggs, int maxPets, int maxEggs, Set<String> discovered});

class PetRepository {
  PetRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId/pets';

  List<PetSpecies>? _species;

  /// How far above its power a pet may hunt.
  int mobReach = 5;

  /// Pet levels × this = the monster level a pet fights like.
  int powerPerLevel = 10;

  Future<List<PetSpecies>> species() async {
    if (_species case final s?) return s;
    final j = await api.get('/v1/pets');
    mobReach = (j['mob_reach'] as num?)?.toInt() ?? mobReach;
    powerPerLevel = (j['power_per_level'] as num?)?.toInt() ?? powerPerLevel;
    return _species = [
      for (final s in (j['species'] as List? ?? const []).cast<Map<String, dynamic>>()) PetSpecies.fromJson(s),
    ];
  }

  Future<PetList> list() async {
    final j = await api.get(_base);
    return (
      pets: [for (final p in (j['pets'] as List? ?? const []).cast<Map<String, dynamic>>()) Pet.fromJson(p)],
      eggs: [for (final e in (j['eggs'] as List? ?? const []).cast<Map<String, dynamic>>()) PetEgg.fromJson(e)],
      maxPets: (j['max_pets'] as num?)?.toInt() ?? 50,
      maxEggs: (j['max_eggs'] as num?)?.toInt() ?? 20,
      discovered: {for (final c in (j['discovered'] as List? ?? const [])) c as String},
    );
  }

  Future<void> hatch(String itemId) => api.post('$_base/hatch', {'item_id': itemId});
  Future<void> setActive(String id, bool on) => api.post('$_base/$id/active', {'active': on});
  Future<void> hunt(String id, String mob, Duration d) =>
      api.post('$_base/$id/hunt', {'mob_code': mob, 'duration_seconds': d.inSeconds});

  Future<PetHuntResult> collect(String id) async {
    final j = await api.post('$_base/$id/collect', const {});
    return PetHuntResult(
      wins: (j['victories'] as num?)?.toInt() ?? 0,
      xp: (j['xp'] as num?)?.toInt() ?? 0,
      died: j['died'] as bool? ?? false,
      pet: Pet.fromJson(j['pet'] as Map<String, dynamic>),
    );
  }
}

String petName(String code) => ContentNames.of('pet', code);

/// Pets (DD-328): eggs to incubate, the pets with their level, bonus and
/// hunts, and every species' odds.
class PetsScreen extends StatefulWidget {
  const PetsScreen({super.key, required this.repo, required this.adventure, required this.dataVersion});

  final PetRepository repo;
  final AdventureRepository adventure;
  final ValueNotifier<int> dataVersion;

  @override
  State<PetsScreen> createState() => _PetsScreenState();
}

class _PetsScreenState extends State<PetsScreen> with PanelState {
  List<PetSpecies> _species = const [];
  List<Pet> _pets = const [];
  List<PetEgg> _eggs = const [];
  int _max = 50, _maxEggs = 20;
  Set<String> _discovered = const {};

  /// Refreshes the time-left lines; hatches and hunts run for hours.
  Timer? _clock;

  @override
  void initState() {
    super.initState();
    _clock = Timer.periodic(const Duration(seconds: 30), (_) => setState(() {}));
    load();
  }

  @override
  void dispose() {
    _clock?.cancel();
    super.dispose();
  }

  @override
  Future<void> fetch() async {
    _species = await widget.repo.species();
    final got = await widget.repo.list();
    _pets = got.pets;
    _eggs = got.eggs;
    _max = got.maxPets;
    _maxEggs = got.maxEggs;
    _discovered = got.discovered;
  }

  /// A pet's card in a bottom sheet; its buttons close the sheet, then act.
  Future<void> _openPet(Pet p) async {
    final action = await showModalBottomSheet<Future<void> Function()>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Palette.bgSurface,
      builder: (c) => SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(Space.screen, 0, Space.screen, Space.screen),
          child: _PetCard(
            key: Key('pet-${p.id}'),
            pet: p,
            onActive: () => Navigator.pop(c, () => _act(() => widget.repo.setActive(p.id, !p.active))),
            onHunt: () => Navigator.pop(c, () => _hunt(p)),
            onCollect: () => Navigator.pop(c, () => _collect(p)),
          ),
        ),
      ),
    );
    if (action != null && mounted) await action();
  }

  /// The raids whose boss can drop the species' egg, from its item tags.
  Future<List<PetEggSource>> _sources(PetSpecies s) async {
    final info = await ItemInfo.load(widget.repo.api, 'pet_egg_${s.code}');
    final raids = {for (final r in await widget.adventure.raids()) r.dungeon.code: r.dungeon};
    return [
      for (final d in info?.dungeons ?? const <({String code, bool raid})>[])
        if (raids[d.code] case final r?) (boss: r.bossMobCode, region: r.regionCode),
    ];
  }

  void _openSpecies(PetSpecies s) => showPetSpecies(
    context,
    s,
    powerPerLevel: widget.repo.powerPerLevel,
    discovered: _discovered.contains(s.code),
    sources: _sources(s),
  );

  PetSpecies? _of(String code) => _species.where((s) => s.code == code).firstOrNull;

  Future<void> _act(Future<void> Function() action) => run(() async {
    await action();
    await fetch();
    widget.dataVersion.value++; // the active pet's bonus is on the sheet
  });

  Future<void> _hatch(PetEgg e) async {
    final l = AppLocalizations.of(context);
    final s = _of(e.species);
    if (!await confirmAction(
      context,
      title: l.petHatch,
      body: l.petHatchBody(itemName(e.template), shortDuration(l, Duration(seconds: s?.hatchSeconds ?? 0))),
      action: l.petHatch,
    )) {
      return;
    }
    await _act(() => widget.repo.hatch(e.itemId));
  }

  Future<void> _hunt(Pet p) async {
    final mobs = await widget.adventure.mobs();
    if (!mounted) return;
    final pick = await Navigator.of(context).push<(String, Duration)>(
      MaterialPageRoute(
        builder: (_) => _PetHuntPicker(
          pet: p,
          mobs: [
            for (final m in mobs)
              if (m.level <= p.power + widget.repo.mobReach) m,
          ]..sort((a, b) => b.level.compareTo(a.level)),
        ),
      ),
    );
    if (pick != null) await _act(() => widget.repo.hunt(p.id, pick.$1, pick.$2));
  }

  Future<void> _collect(Pet p) async {
    final l = AppLocalizations.of(context);
    if (DateTime.now().isBefore(p.huntReadyAt!) &&
        !await confirmAction(context, title: l.collect, body: l.petCollectEarly, action: l.collect)) {
      return;
    }
    await run(() async {
      final r = await widget.repo.collect(p.id);
      await fetch();
      if (!mounted) return;
      final up = r.pet.level > p.level;
      final msg = Text(
        [
          l.petCollected(r.wins, r.xp),
          if (up) l.petLevelUp(petName(p.species), r.pet.level),
          if (r.died) l.petDied,
        ].join(' · '),
      );
      toast(context, up ? PopIn(child: msg) : msg);
    });
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text(l.petsOpen)),
      body: SafeArea(
        child: PageFrame(
          child: frame(() {
            final index = [..._species]..sort((a, b) => a.rank.compareTo(b.rank));
            Widget header(String title, String count) => Padding(
              padding: const EdgeInsets.only(top: Space.lg, bottom: Space.sm),
              child: Row(
                children: [
                  Expanded(child: Text(title, style: text.titleMedium)),
                  Text(count, style: text.bodySmall),
                ],
              ),
            );
            Widget hint(String s) => Padding(
              padding: const EdgeInsets.only(top: Space.xs),
              child: Text(s, style: text.bodySmall?.copyWith(color: Palette.textMuted)),
            );
            return [
              Text(l.petsHint, style: text.bodySmall, textAlign: TextAlign.justify),
              header(l.petsEggs, l.slotCount(_eggs.length, _maxEggs)),
              _BagGrid(slots: _slots(_eggs.length, _maxEggs), children: [for (final e in _eggs) _eggCell(e)]),
              hint(_eggs.isEmpty ? l.petsNoEggs : l.petEggTapHint),
              header(l.petsMine(_pets.length, _max), ''),
              _BagGrid(slots: _slots(_pets.length, _max), children: [for (final p in _pets) _petCell(p)]),
              if (_pets.isNotEmpty) hint(l.petTapHint),
              header(l.petIndex, l.petIndexCount(_discovered.length, index.length)),
              GridView.count(
                crossAxisCount: 4,
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                mainAxisSpacing: Space.sm,
                crossAxisSpacing: Space.sm,
                childAspectRatio: 0.78,
                children: [for (final s in index) _indexTile(s)],
              ),
              hint(l.petIndexHint),
            ];
          }),
        ),
      ),
    );
  }

  /// Bag size shown: at least [cap] slots, more rows if it ever overflows.
  static int _slots(int used, int cap) => used <= cap ? cap : (used + 4) ~/ 5 * 5;

  Widget _eggCell(PetEgg e) {
    final color = rarityColor(petRarity(_of(e.species)?.rank ?? 1));
    return _Slot(
      key: Key('egg-${e.template}'),
      color: color,
      onTap: () => _hatch(e),
      badge: '×${e.quantity}',
      child: ItemArt(code: e.template, icon: Icons.egg_outlined, color: color, size: 40),
    );
  }

  Widget _petCell(Pet p) {
    final color = rarityColor(petRarity(p.rank));
    final now = DateTime.now();
    final hatched = p.hatched(now);
    final ready = p.huntReadyAt != null && !now.isBefore(p.huntReadyAt!);
    return _Slot(
      key: Key('pet-cell-${p.id}'),
      color: color,
      onTap: () => _openPet(p),
      label: hatched ? AppLocalizations.of(context).levelShort(p.level) : null,
      dot: p.active ? Palette.success : null,
      corner: !hatched
          ? Icons.hourglass_bottom
          : p.huntMob == null
          ? null
          : (ready ? Icons.check_circle : Icons.directions_run),
      child: hatched
          ? PetArt(species: p.species, stage: petStage(p.level), size: 64, color: color)
          : ItemArt(code: 'pet_egg_${p.species}', icon: Icons.egg_outlined, color: color, size: 40),
    );
  }

  Widget _indexTile(PetSpecies s) {
    final known = _discovered.contains(s.code);
    final color = rarityColor(petRarity(s.rank));
    final text = Theme.of(context).textTheme;
    return InkWell(
      key: Key('pet-guide-${s.code}'),
      borderRadius: BorderRadius.circular(Radii.md),
      onTap: () => _openSpecies(s),
      child: Container(
        padding: const EdgeInsets.all(4),
        decoration: BoxDecoration(
          color: Palette.bgSurface,
          borderRadius: BorderRadius.circular(Radii.md),
          border: Border.all(color: known ? color.withValues(alpha: 0.6) : Palette.borderSubtle),
        ),
        child: Column(
          children: [
            Expanded(
              child: LayoutBuilder(
                builder: (_, box) => known
                    ? PetArt(species: s.code, stage: 3, size: box.maxWidth, color: color)
                    : Container(
                        width: box.maxWidth,
                        decoration: BoxDecoration(
                          color: Palette.bgPrimary,
                          borderRadius: BorderRadius.circular(Radii.sm),
                        ),
                        child: const Icon(Icons.help_outline, color: Palette.textMuted),
                      ),
              ),
            ),
            const SizedBox(height: 2),
            Text('#${s.rank.toString().padLeft(2, '0')}', style: text.bodySmall?.copyWith(fontSize: 10)),
            Text(
              known ? petName(s.code) : '???',
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(fontSize: 11, color: known ? color : Palette.textMuted),
            ),
          ],
        ),
      ),
    );
  }
}

/// A five-column bag grid: the filled slots first, then empty ones up to [slots].
class _BagGrid extends StatelessWidget {
  const _BagGrid({required this.slots, required this.children});

  final int slots;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) => GridView.count(
    crossAxisCount: 5,
    shrinkWrap: true,
    physics: const NeverScrollableScrollPhysics(),
    mainAxisSpacing: Space.xs,
    crossAxisSpacing: Space.xs,
    children: [...children, for (var i = children.length; i < slots; i++) const _Slot()],
  );
}

/// One bag slot: art in a rarity frame with an optional count, level label,
/// active dot and status icon; empty when [child] is null.
class _Slot extends StatelessWidget {
  const _Slot({
    super.key,
    this.child,
    this.color = Palette.borderSubtle,
    this.onTap,
    this.badge,
    this.label,
    this.dot,
    this.corner,
  });

  final Widget? child;
  final Color color;
  final VoidCallback? onTap;
  final String? badge, label;
  final Color? dot;
  final IconData? corner;

  @override
  Widget build(BuildContext context) => Material(
    color: child == null ? Palette.bgSurface : Palette.bgElevated,
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(Radii.sm),
      side: BorderSide(
        color: child == null ? Palette.borderSubtle : color.withValues(alpha: 0.85),
        width: child == null ? 1 : 1.5,
      ),
    ),
    clipBehavior: Clip.antiAlias,
    child: InkWell(
      onTap: onTap,
      child: child == null
          ? null
          : Stack(
              children: [
                Positioned.fill(
                  child: FittedBox(fit: BoxFit.cover, child: child),
                ),
                if (label case final t?)
                  Positioned(
                    left: 2,
                    top: 2,
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 3),
                      decoration: BoxDecoration(color: Colors.black54, borderRadius: BorderRadius.circular(3)),
                      child: Text(t, style: const TextStyle(fontSize: 9, fontWeight: FontWeight.w700)),
                    ),
                  ),
                if (dot case final c?)
                  Positioned(
                    right: 3,
                    top: 3,
                    child: Container(
                      width: 9,
                      height: 9,
                      decoration: BoxDecoration(color: c, shape: BoxShape.circle),
                    ),
                  ),
                if (corner case final ic?)
                  Positioned(left: 2, bottom: 2, child: Icon(ic, size: 13, color: Palette.warning)),
                if (badge case final b?)
                  Positioned(
                    right: 3,
                    bottom: 1,
                    child: Text(b, style: const TextStyle(fontSize: 10, fontWeight: FontWeight.w700)),
                  ),
              ],
            ),
    ),
  );
}

class _PetCard extends StatelessWidget {
  const _PetCard({super.key, required this.pet, required this.onActive, required this.onHunt, required this.onCollect});

  final Pet pet;
  final VoidCallback onActive, onHunt, onCollect;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final now = DateTime.now();
    final color = rarityColor(petRarity(pet.rank));
    final hatched = pet.hatched(now);
    final hunting = pet.huntMob != null;
    return Card(
      shape: pet.active
          ? RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(Radii.md),
              side: const BorderSide(color: Palette.success, width: 1.5),
            )
          : null,
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                // Hatching and each new stage swap the art with a glow.
                GlowPulse(
                  trigger: hatched ? petStage(pet.level) : 0,
                  color: color,
                  radius: Radii.md,
                  child: SwapIn(
                    swapKey: hatched ? petStage(pet.level) : 0,
                    child: hatched
                        ? PetArt(species: pet.species, stage: petStage(pet.level), size: 64, color: color)
                        : ItemArt(code: 'pet_egg_${pet.species}', icon: Icons.egg_outlined, color: color, size: 64),
                  ),
                ),
                const SizedBox(width: Space.md),
                Expanded(
                  child: Text(
                    petName(pet.species),
                    style: text.titleSmall?.copyWith(color: color, fontWeight: FontWeight.w600),
                  ),
                ),
                if (pet.active) Text(l.petActive, style: text.bodySmall?.copyWith(color: Palette.success)),
                const SizedBox(width: Space.sm),
                Text(l.levelShort(pet.level), style: text.titleSmall),
              ],
            ),
            if (!hatched) ...[
              const SizedBox(height: Space.xs),
              Text('${l.petHatching} · ${shortDuration(l, pet.hatchReadyAt.difference(now))}', style: text.bodySmall),
            ] else ...[
              const SizedBox(height: Space.xs),
              if (pet.xpToNext > 0)
                StatBar(
                  label: 'XP',
                  value: pet.xp / pet.xpToNext,
                  color: Palette.epic,
                  valueText: '${pet.xp}/${pet.xpToNext}',
                ),
              if (bonusText(l, pet.bonus) case final b?)
                Text(
                  b,
                  style: text.bodySmall?.copyWith(color: Palette.success, fontWeight: FontWeight.w600),
                ),
              if (pet.assistPct > 0)
                Text(
                  l.petAssist(pet.assistPct),
                  style: text.bodySmall?.copyWith(
                    color: pet.active ? Palette.success : Palette.textSecondary,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              Text(l.petPower(pet.power), style: text.bodySmall),
              if (hunting)
                Text(
                  !now.isBefore(pet.huntReadyAt!)
                      ? '${l.petHunting(ContentNames.of('mob', pet.huntMob!))} · ${l.petHuntDone}'
                      : '${l.petHunting(ContentNames.of('mob', pet.huntMob!))} · '
                            '${shortDuration(l, pet.huntReadyAt!.difference(now))}',
                  style: text.bodySmall,
                ),
              const SizedBox(height: Space.sm),
              Wrap(
                spacing: Space.sm,
                children: [
                  OutlinedButton(onPressed: onActive, child: Text(pet.active ? l.petDeactivate : l.petActivate)),
                  if (hunting)
                    FilledButton(onPressed: onCollect, child: Text(l.collect))
                  else
                    FilledButton(onPressed: onHunt, child: Text(l.petHunt)),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}

/// The monsters a pet can hunt (its power + reach, strongest first) and the
/// length; a win is worth the monster's level in pet XP.
class _PetHuntPicker extends StatefulWidget {
  const _PetHuntPicker({required this.pet, required this.mobs});

  final Pet pet;
  final List<MobInfo> mobs;

  @override
  State<_PetHuntPicker> createState() => _PetHuntPickerState();
}

class _PetHuntPickerState extends State<_PetHuntPicker> {
  late String? _mob = widget.mobs.firstOrNull?.code;
  Duration _duration = const Duration(hours: 1);

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text('${petName(widget.pet.species)} · ${l.levelShort(widget.pet.level)}')),
      body: SafeArea(
        child: PageFrame(
          child: ListView(
            padding: const EdgeInsets.all(Space.lg),
            children: [
              Text(l.durationLabel, style: text.titleSmall),
              const SizedBox(height: Space.xs),
              Wrap(
                spacing: Space.sm,
                children: [
                  for (final d in const [Duration(hours: 1), Duration(hours: 3), Duration(hours: 6)])
                    ChoiceChip(
                      label: Text(shortDuration(l, d)),
                      selected: d == _duration,
                      onSelected: (_) => setState(() => _duration = d),
                    ),
                ],
              ),
              const SizedBox(height: Space.md),
              FilledButton(
                key: const Key('pet-hunt-send'),
                onPressed: _mob == null ? null : () => Navigator.of(context).pop((_mob!, _duration)),
                child: Text(l.petHunt),
              ),
              const SizedBox(height: Space.lg),
              Text(l.petPickMob, style: text.titleSmall),
              for (final m in widget.mobs)
                ListTile(
                  key: Key('pet-mob-${m.code}'),
                  contentPadding: EdgeInsets.zero,
                  selected: m.code == _mob,
                  leading: ItemArt(
                    code: m.code,
                    icon: Icons.whatshot_outlined,
                    color: m.code == _mob ? Palette.accent : Palette.textSecondary,
                  ),
                  title: Text(ContentNames.of('mob', m.code)),
                  subtitle: Text(l.petXpPerWin(m.level)),
                  trailing: Text(l.levelShort(m.level)),
                  onTap: () => setState(() => _mob = m.code),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
