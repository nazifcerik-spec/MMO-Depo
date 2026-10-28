import 'package:flutter/material.dart';

import '../../app/theme.dart';
import '../../l10n/app_localizations.dart';
import '../../shared/content_names.dart';
import '../../shared/widgets.dart';

/// One side's numbers as it fought (DD-308).
class FighterStats {
  const FighterStats({
    this.level = 0,
    this.hpMax = 0,
    this.attack = 0,
    this.defense = 0,
    this.hitPct = 0,
    this.evasionPct = 0,
    this.critPct = 0,
  });

  factory FighterStats.fromJson(Map<String, dynamic>? j) {
    int n(String k) => (j?[k] as num?)?.toInt() ?? 0;
    return FighterStats(
      level: n('level'),
      hpMax: n('hp_max'),
      attack: n('attack'),
      defense: n('defense'),
      hitPct: n('hit_pct'),
      evasionPct: n('evasion_pct'),
      critPct: n('crit_pct'),
    );
  }

  final int level, hpMax, attack, defense, hitPct, evasionPct, critPct;
}

/// One side's swings summed over the run's fights.
class SwingSums {
  const SwingSums({this.swings = 0, this.hits = 0, this.crits = 0, this.damage = 0, this.maxHit = 0});

  factory SwingSums.fromJson(Map<String, dynamic>? j) {
    int n(String k) => (j?[k] as num?)?.toInt() ?? 0;
    return SwingSums(
      swings: n('swings'),
      hits: n('hits'),
      crits: n('crits'),
      damage: n('damage'),
      maxHit: n('max_hit'),
    );
  }

  final int swings, hits, crits, damage, maxHit;

  int get hitPct => swings == 0 ? 0 : (hits * 100 / swings).round();
  int get avg => hits == 0 ? 0 : (damage / hits).round();
}

/// Why a hunt, dungeon or raid ended and how the fights went (DD-308).
class FightReport {
  const FightReport({
    required this.end,
    this.endMob = '',
    required this.enemyCode,
    required this.you,
    required this.enemy,
    required this.hitChance,
    required this.enemyHitChance,
    required this.fights,
    required this.yours,
    required this.theirs,
    required this.lowestHpPct,
    this.restSeconds = 0,
  });

  /// [restSeconds]: a hunt's camping time; 0 elsewhere.
  static FightReport? fromJson(Object? j, {int restSeconds = 0}) {
    if (j is! Map<String, dynamic>) return null;
    int n(String k) => (j[k] as num?)?.toInt() ?? 0;
    return FightReport(
      end: j['end'] as String? ?? '',
      endMob: j['end_mob'] as String? ?? '',
      enemyCode: j['enemy_code'] as String? ?? '',
      you: FighterStats.fromJson(j['you'] as Map<String, dynamic>?),
      enemy: FighterStats.fromJson(j['enemy'] as Map<String, dynamic>?),
      hitChance: n('hit_chance'),
      enemyHitChance: n('enemy_hit_chance'),
      fights: n('fights'),
      yours: SwingSums.fromJson(j['yours'] as Map<String, dynamic>?),
      theirs: SwingSums.fromJson(j['theirs'] as Map<String, dynamic>?),
      lowestHpPct: n('lowest_hp_pct'),
      restSeconds: restSeconds,
    );
  }

  final String end, endMob, enemyCode;
  final FighterStats you, enemy;
  final int hitChance, enemyHitChance, fights, lowestHpPct, restSeconds;
  final SwingSums yours, theirs;

  /// The reason line; [party]: a raid, where the whole party fell.
  String endText(AppLocalizations l, {bool party = false}) {
    final mob = ContentNames.of('mob', endMob);
    return switch (end) {
      'time_up' => l.endTimeUp,
      'collected_early' => l.endEarly,
      'fled' => l.endFled(mob),
      'defeat' => party ? l.endPartyDefeat(mob) : l.endDefeat(mob),
      'timeout' => l.endTimeout(mob),
      'cleared' => l.endCleared,
      _ => '',
    };
  }

  bool get bad => end == 'fled' || end == 'defeat' || end == 'timeout';
}

/// The report under a reveal's loot: the end reason, the two sides' stats
/// and how the swings went. [party]: a raid (your stats, the party's swings).
class FightReportCard extends StatelessWidget {
  const FightReportCard({super.key, required this.report, this.party = false, this.hunt = false});

  final FightReport report;
  final bool party;

  /// Explains why a hunt seldom ends at 5% HP: it camps first.
  final bool hunt;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    final r = report;
    final reason = r.endText(l, party: party);
    final muted = text.bodySmall?.copyWith(color: Palette.textSecondary);
    TableRow row(String label, String you, String enemy, {bool head = false}) {
      final style = head ? text.labelMedium : text.bodySmall;
      return TableRow(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 2),
            child: Text(label, style: muted),
          ),
          Text(you, style: style, textAlign: TextAlign.end),
          Text(enemy, style: style, textAlign: TextAlign.end),
        ],
      );
    }

    String swings(SwingSums s) => l.repSwings(s.hits, s.swings, s.hitPct, s.avg, s.maxHit);

    return Card(
      key: const Key('fight-report'),
      color: Palette.bgElevated,
      child: Padding(
        padding: const EdgeInsets.all(Space.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(l.fightReport, style: text.titleSmall),
            if (reason.isNotEmpty) ...[
              const SizedBox(height: Space.xs),
              Text(reason, style: text.bodyMedium?.copyWith(color: r.bad ? Palette.danger : Palette.success)),
            ],
            if (hunt && (r.end == 'time_up' || r.end == 'collected_early')) ...[
              const SizedBox(height: Space.xs),
              Text(l.repCampHint, style: muted),
            ],
            const SizedBox(height: Space.md),
            Table(
              columnWidths: const {0: FlexColumnWidth(1.4), 1: FlexColumnWidth(), 2: FlexColumnWidth()},
              children: [
                row('', l.repYou, ContentNames.of('mob', r.enemyCode), head: true),
                row(l.repLevel, '${r.you.level}', '${r.enemy.level}'),
                row(l.repHp, '${r.you.hpMax}', '${r.enemy.hpMax}'),
                row(l.repAttack, '${r.you.attack}', '${r.enemy.attack}'),
                row(l.repDefense, '${r.you.defense}', '${r.enemy.defense}'),
                row(l.repEvasion, '${r.you.evasionPct}', '${r.enemy.evasionPct}'),
                row(l.repHitChance, l.percent(r.hitChance), l.percent(r.enemyHitChance)),
                row(l.repCrit, l.percent(r.you.critPct), l.percent(r.enemy.critPct)),
              ],
            ),
            const SizedBox(height: Space.md),
            Text(l.repFights(r.fights), style: text.bodySmall),
            const SizedBox(height: Space.xs),
            Text('${party ? l.repPartySwings : l.repYourSwings}: ${swings(r.yours)}', style: text.bodySmall),
            Text('${l.repEnemySwings}: ${swings(r.theirs)}', style: text.bodySmall),
            Text(l.repLowest(r.lowestHpPct), style: text.bodySmall),
            if (r.restSeconds > 0)
              Text(l.repRest(shortDuration(l, Duration(seconds: r.restSeconds))), style: text.bodySmall),
          ],
        ),
      ),
    );
  }
}
