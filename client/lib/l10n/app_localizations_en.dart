// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Otter';

  @override
  String get tagline => 'Adventures that keep going while you\'re away.';

  @override
  String get emailOrUsername => 'Email or username';

  @override
  String get email => 'Email';

  @override
  String get username => 'Username';

  @override
  String get password => 'Password';

  @override
  String get passwordConfirm => 'Confirm password';

  @override
  String get showPassword => 'Show password';

  @override
  String get hidePassword => 'Hide password';

  @override
  String get login => 'Log in';

  @override
  String get register => 'Create account';

  @override
  String get noAccount => 'New here? Create an account';

  @override
  String get haveAccount => 'Already have an account? Log in';

  @override
  String get acceptTerms => 'I accept the Terms of Service and Privacy Policy';

  @override
  String get fieldRequired => 'Required';

  @override
  String get invalidEmail => 'Enter a valid email';

  @override
  String get usernameRule => '3–32 letters, digits, _ or .';

  @override
  String get passwordMin => 'At least 8 characters';

  @override
  String get passwordMismatch => 'Passwords do not match';

  @override
  String get termsRequired => 'Please accept the terms to continue';

  @override
  String get logout => 'Log out';

  @override
  String get retry => 'Retry';

  @override
  String get offline => 'Can\'t reach the server. Check your connection.';

  @override
  String get characters => 'Characters';

  @override
  String slotsUsed(int used, int cap) {
    return '$used/$cap slots';
  }

  @override
  String get newCharacter => 'New character';

  @override
  String get noCharacters => 'No characters yet. Create your first hero.';

  @override
  String get create => 'Create';

  @override
  String get characterName => 'Name';

  @override
  String get nameRule => '3–16 characters, starts with a letter (A–Z, 0–9, _)';

  @override
  String get race => 'Race';

  @override
  String get characterClass => 'Class';

  @override
  String get raceHuman => 'Human';

  @override
  String get raceElf => 'Elf';

  @override
  String get raceDwarf => 'Dwarf';

  @override
  String get raceOrc => 'Orc';

  @override
  String get classWarrior => 'Warrior';

  @override
  String get classMage => 'Mage';

  @override
  String get classHunter => 'Hunter';

  @override
  String get classPriest => 'Priest';

  @override
  String get classRogue => 'Rogue';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Lvl $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Lvl $level';
  }

  @override
  String get hp => 'HP';

  @override
  String get mp => 'MP';

  @override
  String get stamina => 'ST';

  @override
  String get energy => 'EN';

  @override
  String get xp => 'XP';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => 'Attributes';

  @override
  String get combat => 'Combat';

  @override
  String get statStr => 'STR';

  @override
  String get statAgi => 'AGI';

  @override
  String get statInt => 'INT';

  @override
  String get statVit => 'VIT';

  @override
  String get statCha => 'CHA';

  @override
  String get statLuck => 'LUCK';

  @override
  String get attack => 'Attack';

  @override
  String get defense => 'Defense';

  @override
  String get crit => 'Crit';

  @override
  String get hit => 'Hit';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => 'Gold';

  @override
  String get switchCharacter => 'Switch character';

  @override
  String get statePendingDelete => 'Deletion pending';

  @override
  String get stateFrozen => 'Frozen';

  @override
  String get errInvalidCredentials => 'Wrong email/username or password.';

  @override
  String get errAccountLocked => 'Your account is temporarily locked.';

  @override
  String get errAccountBanned => 'This account has been banned.';

  @override
  String get errAccountExists => 'That email or username is already taken.';

  @override
  String get errRateLimited => 'Too many attempts. Please wait a moment.';

  @override
  String get errNameTaken => 'That name is already taken.';

  @override
  String get errSlotCap => 'All character slots are in use.';

  @override
  String get errNotPlayable => 'This character can\'t be played right now.';

  @override
  String get errSessionExpired => 'Your session expired. Please log in again.';

  @override
  String errGeneric(String code) {
    return 'Something went wrong ($code).';
  }

  @override
  String get navCharacter => 'Character';

  @override
  String get navActivity => 'Gather';

  @override
  String get actMining => 'Mining';

  @override
  String get actFishing => 'Fishing';

  @override
  String get actWoodcutting => 'Woodcutting';

  @override
  String get actHerbalism => 'Herbalism';

  @override
  String get actHunting => 'Hunting';

  @override
  String get activityChoose => 'Choose an activity';

  @override
  String get durationLabel => 'Duration';

  @override
  String startActivity(String name) {
    return 'Start $name';
  }

  @override
  String get collect => 'Collect';

  @override
  String get abandon => 'Abandon';

  @override
  String get cancel => 'Cancel';

  @override
  String get soFar => 'Gathered so far';

  @override
  String get activityDone => 'Complete — ready to collect';

  @override
  String remaining(String time) {
    return '$time left';
  }

  @override
  String durHm(int h, String m) {
    return '${h}h ${m}m';
  }

  @override
  String durM(int m) {
    return '${m}m';
  }

  @override
  String durS(int s) {
    return '${s}s';
  }

  @override
  String get collectEarlyTitle => 'Collect early?';

  @override
  String get collectEarlyBody => 'You get what has been gathered so far and the activity ends.';

  @override
  String get abandonTitle => 'Abandon activity?';

  @override
  String get abandonBody => 'Everything gathered so far will be lost.';

  @override
  String get lootCollect => 'Nice!';

  @override
  String get lootNothing => 'Nothing this time.';

  @override
  String xpGain(int xp) {
    return '+$xp XP';
  }

  @override
  String get rarityCommon => 'Common';

  @override
  String get rarityUncommon => 'Uncommon';

  @override
  String get rarityRare => 'Rare';

  @override
  String get rarityEpic => 'Epic';

  @override
  String get rarityLegendary => 'Legendary';

  @override
  String get rarityMythic => 'Mythic';

  @override
  String get errSlotBusy => 'Another activity is already running.';

  @override
  String get errNotActive => 'This activity has already ended.';

  @override
  String get evAchievement => 'Achievement unlocked!';

  @override
  String get evMasteryTier => 'Mastery tier up!';

  @override
  String get evMarketSold => 'One of your listings sold.';

  @override
  String get evMailReceived => 'You have new mail.';

  @override
  String get evPartyInvite => 'You were invited to a party.';

  @override
  String get evGuildInvite => 'You were invited to a guild.';

  @override
  String get navBag => 'Bag';

  @override
  String get equipment => 'Equipment';

  @override
  String get bag => 'Bag';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap kg';
  }

  @override
  String get slotHelm => 'Helm';

  @override
  String get slotNecklace => 'Neck';

  @override
  String get slotChest => 'Chest';

  @override
  String get slotGloves => 'Gloves';

  @override
  String get slotLegs => 'Legs';

  @override
  String get slotBoots => 'Boots';

  @override
  String get slotMainhand => 'Main hand';

  @override
  String get slotOffhand => 'Off hand';

  @override
  String get slotRing => 'Ring';

  @override
  String get slotTrinket => 'Trinket';

  @override
  String get equip => 'Equip';

  @override
  String get unequip => 'Unequip';

  @override
  String get sell => 'Sell';

  @override
  String get destroy => 'Destroy';

  @override
  String sellTitle(String name) {
    return 'Sell $name?';
  }

  @override
  String get sellBody => 'You can buy it back from the vendor for 2 hours.';

  @override
  String destroyTitle(String name) {
    return 'Destroy $name?';
  }

  @override
  String get destroyBody => 'This cannot be undone.';

  @override
  String get soldFor => 'Sold for';

  @override
  String get quantity => 'Quantity';

  @override
  String get requiresLevel => 'Requires level';

  @override
  String get weight => 'Weight';

  @override
  String weightGrams(int g) {
    return '$g g';
  }

  @override
  String get sellValue => 'Sells for';

  @override
  String get bound => 'Soulbound';

  @override
  String get itemLocked => 'This item is held by a trade, mail or listing.';

  @override
  String get errBagFull => 'Your bag is full.';

  @override
  String get errItemLocked => 'This item is locked.';

  @override
  String get errClassRestriction => 'Your class can\'t use this item.';

  @override
  String get errLevelRequirement => 'Your level is too low.';

  @override
  String get errNotEquippable => 'This item can\'t be equipped.';

  @override
  String get navAdventure => 'Adventure';

  @override
  String get segHunt => 'Hunt';

  @override
  String get segDungeon => 'Dungeon';

  @override
  String get segTravel => 'Travel';

  @override
  String get huntChoose => 'Choose your prey';

  @override
  String get startHunt => 'Start hunt';

  @override
  String get huntDone => 'Hunt over — collect your spoils';

  @override
  String get huntEarlyBody => 'You keep the fights won so far and the hunt ends.';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'HP $hp · ATK $atk · DEF $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '$won won · $lost lost';
  }

  @override
  String levelUp(int level) {
    return 'Level up! Now level $level';
  }

  @override
  String levelReq(int level) {
    return 'Lvl $level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '$count fights · boss $boss';
  }

  @override
  String get dungeonEnter => 'Enter';

  @override
  String get dungeonDone => 'Run complete — claim rewards';

  @override
  String get dungeonWaitHint => 'Rewards unlock when the run ends.';

  @override
  String get dungeonCleared => 'Cleared!';

  @override
  String get dungeonFailed => 'Party wiped';

  @override
  String get bossKilled => 'Boss defeated';

  @override
  String get bonus => 'Clear bonus';

  @override
  String get youAreIn => 'You are in';

  @override
  String get routes => 'Routes';

  @override
  String get noRoutes => 'No routes from here.';

  @override
  String get travelTo => 'Go';

  @override
  String travelling(String place) {
    return 'Travelling to $place';
  }

  @override
  String get arrive => 'Arrive';

  @override
  String get arrived => 'You have arrived';

  @override
  String get cancelTravel => 'Turn back';

  @override
  String get kindCity => 'City';

  @override
  String get kindZone => 'Zone';

  @override
  String get kindWilderness => 'Wilderness';

  @override
  String get kindDungeon => 'Dungeon';

  @override
  String get kindContinent => 'Continent';

  @override
  String get safeZone => 'safe';

  @override
  String get errBusy => 'Something is already in progress.';

  @override
  String get errNotReady => 'Not finished yet.';

  @override
  String get errWeeklyLocked => 'You already cleared this dungeon this week.';

  @override
  String get errNoRoute => 'There is no direct route there.';

  @override
  String get navWorkshop => 'Workshop';

  @override
  String get segCraft => 'Craft';

  @override
  String get segShop => 'Shop';

  @override
  String get profBlacksmithing => 'Blacksmithing';

  @override
  String get profAlchemy => 'Alchemy';

  @override
  String get profCooking => 'Cooking';

  @override
  String get learn => 'Learn';

  @override
  String get learnHint => 'Learn this profession to unlock its recipes.';

  @override
  String get craft => 'Craft';

  @override
  String crafted(String name, int count) {
    return 'Crafted $name ×$count';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession is now level $level';
  }

  @override
  String requiresSkill(int level) {
    return 'Needs skill $level';
  }

  @override
  String get buy => 'Buy';

  @override
  String bought(String name) {
    return 'Bought $name';
  }

  @override
  String get buybacks => 'Buy back';

  @override
  String get buybackHint => 'Items you sold can be bought back for 2 hours.';

  @override
  String get errInsufficientGold => 'Not enough gold.';

  @override
  String get errInsufficientMaterials => 'You\'re missing materials.';

  @override
  String get errSkillTooLow => 'Your skill is too low for this recipe.';

  @override
  String get errProfessionNotLearned => 'Learn the profession first.';

  @override
  String get errPriceChanged => 'The price changed — check it and try again.';

  @override
  String get errCannotSell => 'This item can\'t be sold.';

  @override
  String get errBuybackExpired => 'The buy-back window has closed.';

  @override
  String get segMarket => 'Market';

  @override
  String get classWeapon => 'Weapons';

  @override
  String get classArmor => 'Armor';

  @override
  String get classConsumable => 'Consumables';

  @override
  String get classMaterial => 'Materials';

  @override
  String get filterAll => 'All';

  @override
  String get sortCheapest => 'Cheapest';

  @override
  String get sortNewest => 'Newest';

  @override
  String get marketEmpty => 'No listings match.';

  @override
  String listingLine(String seller, String price) {
    return '$seller · $price each';
  }

  @override
  String get loadMore => 'Load more';

  @override
  String myListings(int count) {
    return 'My listings ($count)';
  }

  @override
  String get cancelListing => 'Cancel listing';

  @override
  String cancelListingTitle(String name) {
    return 'Cancel $name?';
  }

  @override
  String get cancelListingBody =>
      'The item returns to your bag. Part of the listing fee is refunded only within the first 6 hours.';

  @override
  String buyListingTitle(String name) {
    return 'Buy $name?';
  }

  @override
  String buyListingBody(String price, String seller) {
    return 'Pay $price to $seller.';
  }

  @override
  String get listOnMarket => 'List on market';

  @override
  String listed(String name) {
    return '$name listed on the market';
  }

  @override
  String quantityOf(int count, int max) {
    return 'Quantity: $count of $max';
  }

  @override
  String get unitPrice => 'Price per item';

  @override
  String suggestedPrice(String price) {
    return 'Market price: $price';
  }

  @override
  String get listingFee => 'Listing fee (non-refundable)';

  @override
  String get youReceive => 'You receive when sold';

  @override
  String get errPriceOutOfRange => 'That price is outside the allowed range.';

  @override
  String get errListingCap => 'You have too many active listings.';

  @override
  String get errSelfBuy => 'You can\'t buy your own listing.';

  @override
  String get errListingUnavailable => 'This listing is no longer available.';

  @override
  String get errNotTradable => 'This item can\'t be traded.';

  @override
  String get social => 'Social';

  @override
  String get tabChat => 'Chat';

  @override
  String get tabFriends => 'Friends';

  @override
  String get tabMail => 'Mail';

  @override
  String get chanGlobal => 'World';

  @override
  String get chanGuild => 'Guild';

  @override
  String get chanParty => 'Party';

  @override
  String get chatEmpty => 'No messages yet. Say hello!';

  @override
  String chatHint(String channel) {
    return 'Message $channel';
  }

  @override
  String get send => 'Send';

  @override
  String get addFriend => 'Add';

  @override
  String friendAdded(String name) {
    return 'You and $name are now friends';
  }

  @override
  String friendRequested(String name) {
    return 'Friend request sent to $name';
  }

  @override
  String get friendRequests => 'Friend requests';

  @override
  String friendsCount(int count) {
    return 'Friends ($count)';
  }

  @override
  String get noFriends => 'No friends yet — add someone by name.';

  @override
  String get more => 'More';

  @override
  String get whisper => 'Whisper';

  @override
  String get remove => 'Remove';

  @override
  String get block => 'Block';

  @override
  String get unblock => 'Unblock';

  @override
  String get blockedList => 'Blocked';

  @override
  String removeFriendTitle(String name) {
    return 'Remove $name?';
  }

  @override
  String get removeFriendBody => 'You can send a new friend request later.';

  @override
  String blockTitle(String name) {
    return 'Block $name?';
  }

  @override
  String get blockBody =>
      'They can\'t whisper you, their chat is hidden and their mail is dropped. This also ends your friendship.';

  @override
  String get accept => 'Accept';

  @override
  String get decline => 'Decline';

  @override
  String get newMail => 'Write mail';

  @override
  String get inboxEmpty => 'Your mailbox is empty.';

  @override
  String get noSubject => '(no subject)';

  @override
  String mailFrom(String name) {
    return 'From $name';
  }

  @override
  String get attachments => 'Attachments';

  @override
  String get claim => 'Take attachments';

  @override
  String get mailClaimed => 'Attachments moved to your bag and wallet';

  @override
  String get deleteMailTitle => 'Delete this mail?';

  @override
  String get deleteMailBody => 'Deleted mail can\'t be recovered.';

  @override
  String mailSent(String postage) {
    return 'Mail sent (postage $postage)';
  }

  @override
  String get mailTo => 'To (character name)';

  @override
  String get mailSubject => 'Subject';

  @override
  String get mailBody => 'Message';

  @override
  String get attachGold => 'Attach gold';

  @override
  String get errAlreadyFriends => 'You\'re already friends.';

  @override
  String get errAlreadyRequested => 'A friend request is already pending.';

  @override
  String get errBlocked => 'You can\'t interact with this player.';

  @override
  String get errFriendCap => 'Your friend list is full.';

  @override
  String get errTargetNotFound => 'No character with that name.';

  @override
  String get errSelfTarget => 'You can\'t target yourself.';

  @override
  String get errInboxFull => 'Their mailbox is full.';

  @override
  String get errTooLong => 'That\'s too long.';

  @override
  String get errEmpty => 'Write something first.';

  @override
  String get errNotInChannel => 'You\'re not in that group.';

  @override
  String get errMuted => 'You are muted and can\'t chat right now.';

  @override
  String get errHasAttachments => 'Take the attachments before deleting.';

  @override
  String get useItem => 'Use';

  @override
  String restoresHp(int amount) {
    return 'Restores $amount HP';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return 'Used $name · HP $hp/$max';
  }

  @override
  String get errHpFull => 'Your health is already full.';

  @override
  String get errInCombat => 'Finish your hunt or dungeon run first.';

  @override
  String get errNotUsable => 'This item can\'t be used.';

  @override
  String get errTooWounded => 'You\'re too wounded to fight. Rest or drink a potion.';

  @override
  String get forgotPassword => 'Forgot password?';

  @override
  String get resetHint => 'Enter your account e-mail and we\'ll send you a link to set a new password.';

  @override
  String get resetSent => 'If an account exists for that e-mail, a reset link is on its way. It expires in 1 hour.';

  @override
  String get sendLink => 'Send link';

  @override
  String get close => 'Close';

  @override
  String get changePassword => 'Change password';

  @override
  String get currentPassword => 'Current password';

  @override
  String get newPassword => 'New password';

  @override
  String get save => 'Save';

  @override
  String get passwordChanged => 'Password changed. Your other devices were signed out.';

  @override
  String get verifyEmailBanner => 'Please confirm your e-mail address.';

  @override
  String get resend => 'Resend';

  @override
  String get verificationSent => 'Verification e-mail sent.';

  @override
  String get errEmailNotVerified => 'Confirm your e-mail before creating a character.';

  @override
  String get errAlreadyVerified => 'Your e-mail is already confirmed.';

  @override
  String get tabParty => 'Party';

  @override
  String get tabGuild => 'Guild';

  @override
  String get invite => 'Invite';

  @override
  String invitedBy(String name) {
    return 'Invited by $name';
  }

  @override
  String get noParty => 'You\'re not in a party. Create one and invite friends — fights go better together.';

  @override
  String get createParty => 'Create party';

  @override
  String get pendingInvites => 'Invitations';

  @override
  String partyMembers(int count, int cap) {
    return 'Party ($count/$cap)';
  }

  @override
  String get makeLeader => 'Make leader';

  @override
  String get kick => 'Remove';

  @override
  String get leave => 'Leave';

  @override
  String get disband => 'Disband';

  @override
  String get leavePartyTitle => 'Leave the party?';

  @override
  String get leavePartyBody => 'If you are the leader, the longest-standing member takes over.';

  @override
  String get disbandPartyTitle => 'Disband the party?';

  @override
  String get disbandPartyBody => 'Everyone is removed from the party.';

  @override
  String get rankLeader => 'Leader';

  @override
  String get rankOfficer => 'Officer';

  @override
  String get rankVeteran => 'Veteran';

  @override
  String get rankMember => 'Member';

  @override
  String get rankInitiate => 'Initiate';

  @override
  String get createGuild => 'Found a guild';

  @override
  String createGuildBody(String fee) {
    return 'Founding a guild costs $fee.';
  }

  @override
  String createGuildFee(String fee) {
    return 'Found guild ($fee)';
  }

  @override
  String get guildName => 'Guild name';

  @override
  String get guildTag => 'Tag';

  @override
  String get guildTagHint => '2–5 letters or digits';

  @override
  String guildMembers(int count, int cap) {
    return '$count/$cap members';
  }

  @override
  String disbandScheduled(String time) {
    return 'Guild disbands in $time';
  }

  @override
  String get motd => 'Message of the day';

  @override
  String get editMotd => 'Edit message of the day';

  @override
  String get setRank => 'Change rank';

  @override
  String setRankFor(String name) {
    return 'New rank for $name';
  }

  @override
  String transferTitle(String name) {
    return 'Make $name guild leader?';
  }

  @override
  String get transferBody => 'You become an officer. This can\'t be undone by you.';

  @override
  String kickTitle(String name) {
    return 'Remove $name from the guild?';
  }

  @override
  String get leaveGuildTitle => 'Leave the guild?';

  @override
  String get leaveGuildBody => 'You\'ll need a new invitation to rejoin.';

  @override
  String get disbandGuildTitle => 'Disband the guild?';

  @override
  String get disbandGuildBody => 'The guild disbands in 48 hours; you can cancel until then.';

  @override
  String get cancelDisband => 'Cancel disbanding';

  @override
  String get errAlreadyInGroup => 'They\'re already in a group.';

  @override
  String get errNotLeader => 'Only the leader can do that.';

  @override
  String get errGroupFull => 'The group is full.';

  @override
  String get errInviteInvalid => 'That invitation has expired.';

  @override
  String get errAlreadyInvited => 'An invitation is already pending.';

  @override
  String get errGuildForbidden => 'Your rank can\'t do that.';

  @override
  String get errLeaderCannotLeave => 'Hand over leadership or disband the guild first.';

  @override
  String get errInvalidGuildName => '3–24 letters or digits.';

  @override
  String get errInvalidTag => 'The tag must be 2–5 letters or digits.';

  @override
  String get errTagTaken => 'That tag is taken.';
}
