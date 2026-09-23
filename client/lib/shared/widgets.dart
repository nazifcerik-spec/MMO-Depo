import 'package:flutter/material.dart';

import '../app/theme.dart';
import '../core/api_client.dart';
import '../features/character/models.dart';
import '../l10n/app_localizations.dart';

/// Maps a server/transport error to player-facing copy.
String errorText(AppLocalizations l, Object error) {
  if (error is! ApiError) return l.errGeneric('E_CLIENT');
  return switch (error.code) {
    'E_NETWORK' => l.offline,
    'E_INVALID_CREDENTIALS' => l.errInvalidCredentials,
    'E_ACCOUNT_LOCKED' => l.errAccountLocked,
    'E_ACCOUNT_BANNED' => l.errAccountBanned,
    'E_ACCOUNT_EXISTS' => l.errAccountExists,
    'E_RATE_LIMIT' => l.errRateLimited,
    'E_NAME_TAKEN' => l.errNameTaken,
    'E_SLOT_CAP' => l.errSlotCap,
    'E_NOT_PLAYABLE' => l.errNotPlayable,
    'E_SLOT_BUSY' => l.errSlotBusy,
    'E_NOT_ACTIVE' => l.errNotActive,
    'E_BAG_FULL' => l.errBagFull,
    'E_ITEM_LOCKED' => l.errItemLocked,
    'E_CLASS_RESTRICTION' => l.errClassRestriction,
    'E_LEVEL_REQUIREMENT' => l.errLevelRequirement,
    'E_NOT_EQUIPPABLE' || 'E_UNKNOWN_SLOT' => l.errNotEquippable,
    'E_HUNT_ACTIVE' || 'E_DUNGEON_ACTIVE' || 'E_TRAVEL_IN_PROGRESS' => l.errBusy,
    'E_HUNT_NOT_ACTIVE' || 'E_DUNGEON_NOT_ACTIVE' || 'E_TRAVEL_NOT_ACTIVE' => l.errNotActive,
    'E_DUNGEON_NOT_READY' || 'E_TRAVEL_NOT_READY' => l.errNotReady,
    'E_DUNGEON_WEEKLY_LOCKED' => l.errWeeklyLocked,
    'E_NO_ROUTE' => l.errNoRoute,
    'E_INSUFFICIENT_GOLD' => l.errInsufficientGold,
    'E_INSUFFICIENT_MATERIALS' => l.errInsufficientMaterials,
    'E_SKILL_TOO_LOW' => l.errSkillTooLow,
    'E_PROFESSION_NOT_LEARNED' => l.errProfessionNotLearned,
    'E_ITEM_BOUND' || 'E_ITEM_NO_SELL' => l.errCannotSell,
    'E_BUYBACK_EXPIRED' || 'E_BUYBACK_CLAIMED' => l.errBuybackExpired,
    'E_VALIDATION' when error.field == 'expected_price_copper' => l.errPriceChanged,
    'E_PRICE_MISMATCH' => l.errPriceChanged,
    'E_PRICE_OUT_OF_RANGE' => l.errPriceOutOfRange,
    'E_MARKET_LISTING_CAP' => l.errListingCap,
    'E_SELF_BUY' => l.errSelfBuy,
    'E_LISTING_UNAVAILABLE' || 'E_LISTING_NOT_FOUND' => l.errListingUnavailable,
    'E_ITEM_NOT_TRADABLE' => l.errNotTradable,
    'E_ALREADY_FRIENDS' => l.errAlreadyFriends,
    'E_ALREADY_REQUESTED' => l.errAlreadyRequested,
    'E_BLOCKED' => l.errBlocked,
    'E_FRIEND_CAP' => l.errFriendCap,
    'E_TARGET_NOT_FOUND' => l.errTargetNotFound,
    'E_SELF_TARGET' => l.errSelfTarget,
    'E_INBOX_FULL' => l.errInboxFull,
    'E_TOO_LONG' => l.errTooLong,
    'E_EMPTY' => l.errEmpty,
    'E_NOT_IN_CHANNEL' => l.errNotInChannel,
    'E_MUTED' => l.errMuted,
    'E_HAS_ATTACHMENTS' => l.errHasAttachments,
    'E_HP_FULL' => l.errHpFull,
    'E_IN_COMBAT' => l.errInCombat,
    'E_NOT_USABLE' => l.errNotUsable,
    'E_TOO_WOUNDED' => l.errTooWounded,
    'E_NOTHING_TO_CLAIM' => l.errNotActive,
    'E_SESSION_EXPIRED' || 'E_SESSION_REVOKED' || 'E_REFRESH_TOKEN_REUSE' => l.errSessionExpired,
    'E_VALIDATION' when error.message.isNotEmpty => error.message,
    _ => l.errGeneric(error.code),
  };
}

String raceName(AppLocalizations l, String code) => switch (code) {
  'human' => l.raceHuman,
  'elf' => l.raceElf,
  'dwarf' => l.raceDwarf,
  'orc' => l.raceOrc,
  _ => code,
};

String className(AppLocalizations l, String code) => switch (code) {
  'warrior' => l.classWarrior,
  'mage' => l.classMage,
  'hunter' => l.classHunter,
  'priest' => l.classPriest,
  'rogue' => l.classRogue,
  _ => code,
};

String rarityName(AppLocalizations l, String r) => switch (r) {
  'common' => l.rarityCommon,
  'uncommon' => l.rarityUncommon,
  'rare' => l.rarityRare,
  'epic' => l.rarityEpic,
  'legendary' => l.rarityLegendary,
  'mythic' => l.rarityMythic,
  _ => r,
};

String activityName(AppLocalizations l, String kind) => switch (kind) {
  'mining' => l.actMining,
  'fishing' => l.actFishing,
  'woodcutting' => l.actWoodcutting,
  'herbalism' => l.actHerbalism,
  'hunting' => l.actHunting,
  _ => kind,
};

/// "iron_ore" → "Iron Ore". Content-pack display names are not served yet, so
/// template codes are humanised; affixes wrap the base name.
String itemName(String templateCode, [String? prefix, String? suffix]) =>
    [if (prefix != null) humanize(prefix), humanize(templateCode), if (suffix != null) humanize(suffix)].join(' ');

/// Content code → display text ("ironhold_city" → "Ironhold City").
String humanize(String code) =>
    code.split('_').where((w) => w.isNotEmpty).map((w) => w[0].toUpperCase() + w.substring(1)).join(' ');

/// Compact duration: "2h 05m", "14m", "45s".
String shortDuration(AppLocalizations l, Duration d) {
  if (d.isNegative) d = Duration.zero;
  final h = d.inHours, m = d.inMinutes % 60, s = d.inSeconds % 60;
  if (h > 0) return l.durHm(h, m.toString().padLeft(2, '0'));
  if (m > 0) return l.durM(m);
  return l.durS(s);
}

/// "12g 34s 56c" as plain text (dialog copy, semantics); zero leading parts omitted.
String goldString(int copper) {
  final p = splitCopper(copper);
  return [if (p.gold > 0) '${p.gold}g', if (p.gold > 0 || p.silver > 0) '${p.silver}s', '${p.copper}c'].join(' ');
}

/// Shows a snackbar, replacing any visible one (the newest message wins).
void toast(BuildContext context, Widget content) {
  ScaffoldMessenger.of(context)
    ..hideCurrentSnackBar()
    ..showSnackBar(SnackBar(content: content));
}

/// Mobile-first column: full width on phones, centred 520px on wide screens.
class PageFrame extends StatelessWidget {
  const PageFrame({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) => Align(
    alignment: Alignment.topCenter,
    child: ConstrainedBox(constraints: const BoxConstraints(maxWidth: 520), child: child),
  );
}

/// Centered message with an optional retry action (errors, empty states).
class MessageView extends StatelessWidget {
  const MessageView({super.key, required this.message, this.onRetry, this.icon = Icons.cloud_off});

  final String message;
  final VoidCallback? onRetry;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(Space.xxl),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 40, color: Palette.textMuted),
            const SizedBox(height: Space.md),
            Text(
              message,
              textAlign: TextAlign.center,
              style: const TextStyle(color: Palette.textSecondary),
            ),
            if (onRetry != null) ...[
              const SizedBox(height: Space.lg),
              OutlinedButton.icon(onPressed: onRetry, icon: const Icon(Icons.refresh), label: Text(l.retry)),
            ],
          ],
        ),
      ),
    );
  }
}

/// Labelled linear bar (HP/MP/ST/EN/XP).
class StatBar extends StatelessWidget {
  const StatBar({super.key, required this.label, required this.value, required this.color, required this.valueText});

  final String label;
  final double value;
  final Color color;
  final String valueText;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: '$label $valueText',
      excludeSemantics: true,
      child: Row(
        children: [
          SizedBox(
            width: 44,
            child: Text(
              label,
              style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w700, color: Palette.textSecondary),
            ),
          ),
          Expanded(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(Radii.sm),
              child: LinearProgressIndicator(
                value: value.isNaN ? 0 : value.clamp(0, 1),
                minHeight: 10,
                color: color,
                backgroundColor: Palette.bgElevated,
              ),
            ),
          ),
          const SizedBox(width: Space.md),
          SizedBox(
            width: 92,
            child: Text(
              valueText,
              textAlign: TextAlign.right,
              style: const TextStyle(fontSize: 12, fontFeatures: numeric, color: Palette.textPrimary),
            ),
          ),
        ],
      ),
    );
  }
}

/// "12g 34s 56c" with coin colours; zero parts above copper are skipped.
class GoldText extends StatelessWidget {
  const GoldText(this.copper, {super.key, this.size = 16});

  final int copper;
  final double size;

  @override
  Widget build(BuildContext context) {
    final p = splitCopper(copper);
    TextSpan part(int v, String unit, Color c) => TextSpan(
      children: [
        TextSpan(
          text: '$v',
          style: TextStyle(color: Palette.textPrimary, fontWeight: FontWeight.w700, fontSize: size),
        ),
        TextSpan(
          text: '$unit ',
          style: TextStyle(color: c, fontWeight: FontWeight.w700, fontSize: size * 0.8),
        ),
      ],
    );
    return Text.rich(
      TextSpan(
        style: const TextStyle(fontFeatures: numeric),
        children: [
          if (p.gold > 0) part(p.gold, 'g', Palette.accent),
          if (p.gold > 0 || p.silver > 0) part(p.silver, 's', const Color(0xFFC9CDD6)),
          part(p.copper, 'c', const Color(0xFFC07A4A)),
        ],
      ),
    );
  }
}
