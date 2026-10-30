import 'dart:async';

import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../core/api_client.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/contract_progress.dart' show confirmAction;
import '../../shared/panel.dart';
import '../../shared/widgets.dart';
import '../adventure/adventure_models.dart';
import '../adventure/adventure_repository.dart';
import '../inventory/inventory_models.dart';
import '../inventory/item_images.dart';
import '../inventory/item_widgets.dart';

/// One pet species (DD-328), best first in the catalog.
class PetSpecies {
  const PetSpecies({
    required this.code,
    required this.rank,
    required this.minMobLevel,
    required this.eggPct,
    required this.hatchSeconds,
    required this.bonus,
  });

  factory PetSpecies.fromJson(Map<String, dynamic> j) => PetSpecies(
    code: j['code'] as String,
    rank: (j['rank'] as num).toInt(),
    minMobLevel: (j['min_mob_level'] as num).toInt(),
    eggPct: (j['egg_pct'] as num).toDouble(),
    hatchSeconds: (j['hatch_seconds'] as num).toInt(),
    bonus: _bonus(j['bonus'] as Map<String, dynamic>),
  );

  final String code;
  final int rank, minMobLevel, hatchSeconds;
  final double eggPct;
  final StatBonus bonus;
}

StatBonus _bonus(Map<String, dynamic> b) =>
    StatBonus(b['kind'] as String, b['target'] as String, (b['amount'] as num).toDouble());

/// A pet's rank as an item rarity: its colour.
String petRarity(int rank) => switch (rank) {
  20 => 'mythic',
  >= 16 => 'legendary',
  >= 11 => 'epic',
  >= 6 => 'rare',
  _ => 'uncommon',
};

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

class PetRepository {
  PetRepository(this.api, this.characterId);

  final ApiClient api;
  final String characterId;

  String get _base => '/v1/character/$characterId/pets';

  List<PetSpecies>? _species;

  /// Pet levels × this = the monster level a pet fights like.
  int mobReach = 5;

  Future<List<PetSpecies>> species() async {
    if (_species case final s?) return s;
    final j = await api.get('/v1/pets');
    mobReach = (j['mob_reach'] as num?)?.toInt() ?? mobReach;
    return _species = [
      for (final s in (j['species'] as List? ?? const []).cast<Map<String, dynamic>>()) PetSpecies.fromJson(s),
    ];
  }

  Future<(List<Pet>, List<PetEgg>, int)> list() async {
    final j = await api.get(_base);
    return (
      [for (final p in (j['pets'] as List? ?? const []).cast<Map<String, dynamic>>()) Pet.fromJson(p)],
      [for (final e in (j['eggs'] as List? ?? const []).cast<Map<String, dynamic>>()) PetEgg.fromJson(e)],
      (j['max_pets'] as num?)?.toInt() ?? 50,
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
  int _max = 50;
  bool _guide = false;

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
    final (pets, eggs, max) = await widget.repo.list();
    _pets = pets;
    _eggs = eggs;
    _max = max;
  }

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
      toast(
        context,
        Text(
          [
            l.petCollected(r.wins, r.xp),
            if (r.pet.level > p.level) l.petLevelUp(petName(p.species), r.pet.level),
            if (r.died) l.petDied,
          ].join(' · '),
        ),
      );
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
          child: frame(
            () => [
              Text(l.petsHint, style: text.bodySmall),
              const SizedBox(height: Space.lg),
              Text(l.petsEggs, style: text.titleMedium),
              const SizedBox(height: Space.xs),
              if (_eggs.isEmpty) Text(l.petsNoEggs, style: text.bodySmall),
              for (final e in _eggs) _eggTile(l, e),
              const SizedBox(height: Space.lg),
              Text(l.petsMine(_pets.length, _max), style: text.titleMedium),
              const SizedBox(height: Space.xs),
              for (final p in _pets)
                _PetCard(
                  key: Key('pet-${p.id}'),
                  pet: p,
                  onActive: () => _act(() => widget.repo.setActive(p.id, !p.active)),
                  onHunt: () => _hunt(p),
                  onCollect: () => _collect(p),
                ),
              const SizedBox(height: Space.lg),
              InkWell(
                onTap: () => setState(() => _guide = !_guide),
                child: Row(
                  children: [
                    Expanded(child: Text(l.petGuide, style: text.titleMedium)),
                    Icon(_guide ? Icons.expand_less : Icons.expand_more, color: Palette.textSecondary),
                  ],
                ),
              ),
              if (_guide)
                for (final s in _species) _guideTile(l, text, s),
            ],
          ),
        ),
      ),
    );
  }

  Widget _eggTile(AppLocalizations l, PetEgg e) {
    final s = _of(e.species);
    final color = rarityColor(petRarity(s?.rank ?? 1));
    return Card(
      key: Key('egg-${e.template}'),
      child: ListTile(
        leading: ItemArt(code: e.template, icon: Icons.egg_outlined, color: color, size: 28),
        title: Text('${itemName(e.template)} ×${e.quantity}', style: TextStyle(color: color)),
        subtitle: Text(shortDuration(l, Duration(seconds: s?.hatchSeconds ?? 0))),
        trailing: FilledButton.tonal(style: compactButton, onPressed: () => _hatch(e), child: Text(l.petHatch)),
      ),
    );
  }

  Widget _guideTile(AppLocalizations l, TextTheme text, PetSpecies s) {
    final pct = s.eggPct == s.eggPct.roundToDouble() ? '${s.eggPct.toInt()}' : s.eggPct.toStringAsFixed(2);
    return ListTile(
      dense: true,
      contentPadding: EdgeInsets.zero,
      leading: Icon(Icons.pets, color: rarityColor(petRarity(s.rank))),
      title: Text(petName(s.code), style: TextStyle(color: rarityColor(petRarity(s.rank)))),
      subtitle: Text(
        [
          l.petGuideLine(pct, s.minMobLevel, shortDuration(l, Duration(seconds: s.hatchSeconds))),
          ?bonusText(l, s.bonus),
        ].join('\n'),
        style: text.bodySmall,
      ),
    );
  }
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
                Icon(hatched ? Icons.pets : Icons.egg_outlined, color: color),
                const SizedBox(width: Space.sm),
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
