// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'MeNa World';

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
  String get acceptTerms => 'I am 13 or older and accept the Terms of Service and Privacy Policy';

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
  String get statAgi => 'DEX';

  @override
  String get statInt => 'INT';

  @override
  String get statVit => 'VIT';

  @override
  String get statCha => 'WIS';

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
  String get gatherResource => 'Resource';

  @override
  String gatherXp(int xp, int next) {
    return '$xp / $next XP';
  }

  @override
  String get gatherMax => 'Max level';

  @override
  String gatherLocked(int level) {
    return 'Unlocks at level $level';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name +$xp XP';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '$name reached level $level!';
  }

  @override
  String get errGatherLevel => 'Your gathering level is too low for this resource.';

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
  String get sortLabel => 'Sort';

  @override
  String get sortSlot => 'Bag order';

  @override
  String get sortRarity => 'Rarity';

  @override
  String get sortType => 'Type';

  @override
  String get sortName => 'Name';

  @override
  String get sortValue => 'Value';

  @override
  String selectedCount(int count) {
    return '$count selected';
  }

  @override
  String get selectAll => 'Select all';

  @override
  String get clearSelection => 'Clear selection';

  @override
  String sellSelectedTitle(int count) {
    return 'Sell $count items?';
  }

  @override
  String destroySelectedTitle(int count) {
    return 'Destroy $count items?';
  }

  @override
  String get deleteEntry => 'Delete';

  @override
  String craftHowMany(String name) {
    return 'How many $name?';
  }

  @override
  String craftCount(int count) {
    return 'Craft ×$count';
  }

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
  String recipeCount(int count) {
    return '$count recipes';
  }

  @override
  String get itemType => 'Type';

  @override
  String get stackSize => 'Stack';

  @override
  String get skillXpReward => 'Skill XP';

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
  String get noParty => 'You\'re not in a party. Invite someone and a party forms (2–5 members).';

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
  String get leavePartyBody =>
      'If you are the leader, the longest-standing member takes over. A party left with one member disbands.';

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

  @override
  String get progress => 'Achievements & mastery';

  @override
  String get tabAchievements => 'Achievements';

  @override
  String get tabMastery => 'Mastery';

  @override
  String get tabReputation => 'Reputation';

  @override
  String get tabTitles => 'Titles';

  @override
  String achievementPoints(int points) {
    return '$points achievement points';
  }

  @override
  String pointsShort(int points) {
    return '$points pts';
  }

  @override
  String rewardTitle(String title) {
    return 'Title: $title';
  }

  @override
  String masteryTier(int tier, int max) {
    return 'Tier $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return '+$amount today';
  }

  @override
  String get noTitle => 'No title';

  @override
  String get noTitlesYet => 'Titles unlock through achievements.';

  @override
  String get errTitleNotOwned => 'You haven\'t unlocked that title.';

  @override
  String get notifications => 'Notifications';

  @override
  String get markAllRead => 'Mark all as read';

  @override
  String get inboxNothing => 'No notifications yet.';

  @override
  String get notifSettings => 'Notification settings';

  @override
  String get muteAll => 'Pause push notifications';

  @override
  String get muteHint => 'For 24 hours. The inbox keeps filling.';

  @override
  String mutedFor(String time) {
    return 'Paused for $time';
  }

  @override
  String get modePushInbox => 'Push + inbox';

  @override
  String get modeInbox => 'Inbox only';

  @override
  String get modeSilent => 'Silent';

  @override
  String get modeOff => 'Off';

  @override
  String get nkActivity => 'Activity finished';

  @override
  String get nkDungeon => 'Dungeon run finished';

  @override
  String get nkTravel => 'Arrived';

  @override
  String get nkWorldEvent => 'World event started';

  @override
  String get nkMarketSold => 'Item sold';

  @override
  String get nkMail => 'New mail';

  @override
  String get nkPartyInvite => 'Party invitation';

  @override
  String get nkGuildInvite => 'Guild invitation';

  @override
  String get nkAchievement => 'Achievement unlocked';

  @override
  String get nkSystem => 'Announcement';

  @override
  String timeAgo(String time) {
    return '$time ago';
  }

  @override
  String get segEvents => 'Events';

  @override
  String get eventsEmpty => 'No world events scheduled right now.';

  @override
  String get worldBoss => 'World boss';

  @override
  String get invasion => 'Invasion';

  @override
  String startsIn(String time) {
    return 'Starts in $time';
  }

  @override
  String endsIn(String time) {
    return 'Ends in $time';
  }

  @override
  String regionOnly(String region) {
    return '$region only';
  }

  @override
  String get eventGoal => 'Goal';

  @override
  String get enlist => 'Enlist';

  @override
  String get enlistHint =>
      'Your current power is locked in when you enlist. You score for every minute the event runs; rewards arrive when it ends.';

  @override
  String enlistedPower(int power) {
    return 'Enlisted · power $power';
  }

  @override
  String myContribution(int score, int tier) {
    return 'Your score $score · reward tier $tier';
  }

  @override
  String get eventSucceeded => 'Victory';

  @override
  String get eventFailed => 'Failed';

  @override
  String tokenBalance(int tokens) {
    return '$tokens event tokens';
  }

  @override
  String claimDaily(int tokens) {
    return 'Claim $tokens';
  }

  @override
  String get claimedToday => 'Claimed';

  @override
  String eventTokensGained(int tokens) {
    return '+$tokens event tokens';
  }

  @override
  String get errEventNotOpen => 'Enlistment for this event is closed.';

  @override
  String get errWrongRegion => 'Travel to the event\'s region first.';

  @override
  String get errAlreadyEnlisted => 'You\'re already enlisted.';

  @override
  String get errAlreadyClaimed => 'Already claimed today.';

  @override
  String get errEventNotActive => 'This event isn\'t running.';

  @override
  String gatherYield(String activity) {
    return '$activity yield';
  }

  @override
  String perHour(String value) {
    return '$value/h';
  }

  @override
  String get quests => 'Quests';

  @override
  String get questsInProgress => 'In progress';

  @override
  String get questsNoneOpen => 'No open quests. Pick one below.';

  @override
  String get questsStory => 'Story';

  @override
  String get questsSide => 'Side quests';

  @override
  String get questsDaily => 'Daily';

  @override
  String resetsIn(String time) {
    return 'Resets in $time';
  }

  @override
  String get turnIn => 'Turn in';

  @override
  String questRewarded(String rewards) {
    return 'Quest complete: $rewards';
  }

  @override
  String get abandonQuestTitle => 'Abandon quest?';

  @override
  String get abandonQuestBody => 'Progress is lost. You can take the quest again later.';

  @override
  String get abandonDailyBody => 'Progress is lost and this daily can\'t be taken again today.';

  @override
  String objGather(String activity) {
    return 'Gather ($activity)';
  }

  @override
  String get objKillAny => 'Defeat enemies';

  @override
  String objKill(String mob) {
    return 'Defeat $mob';
  }

  @override
  String get objBoss => 'Defeat a boss';

  @override
  String get objDungeonAny => 'Clear a dungeon';

  @override
  String objDungeon(String dungeon) {
    return 'Clear $dungeon';
  }

  @override
  String get objCraftAny => 'Craft items';

  @override
  String objCraft(String item) {
    return 'Craft $item';
  }

  @override
  String get evQuestComplete => 'Quest objectives complete — turn it in!';

  @override
  String get errQuestLocked => 'Finish the previous story quest first.';

  @override
  String get errQuestTaken => 'You already took this quest.';

  @override
  String get errQuestNotOffered => 'That daily isn\'t on today\'s board.';

  @override
  String get errQuestCap => 'Too many open quests. Finish or abandon one.';

  @override
  String get errQuestNotActive => 'This quest isn\'t open.';

  @override
  String get errQuestIncomplete => 'Objectives aren\'t complete yet.';

  @override
  String get skills => 'Skills';

  @override
  String get skillActive => 'Active';

  @override
  String get skillPassive => 'Passive';

  @override
  String skillStrike(int power) {
    return '$power% damage';
  }

  @override
  String skillHeal(int power) {
    return 'Heals $power% HP when below half';
  }

  @override
  String skillMp(int mp) {
    return '$mp MP';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: 'every $turns turns', one: 'every turn');
    return '$_temp0';
  }

  @override
  String get termsOfService => 'Terms of Service';

  @override
  String get privacyPolicy => 'Privacy Policy';

  @override
  String get accountMenu => 'Account';

  @override
  String get downloadMyData => 'Download my data';

  @override
  String get exportOpened => 'The download link opened in your browser (valid for 15 minutes).';

  @override
  String get linkOpenFailed => 'Could not open the link.';

  @override
  String get deleteAccount => 'Delete my account';

  @override
  String get deleteAccountBody =>
      'Your account and all characters will be permanently deleted after 30 days. Signing in again within that time cancels the deletion. Enter your password to confirm.';

  @override
  String get deleteAccountConfirm => 'Delete account';

  @override
  String deleteScheduled(String date) {
    return 'Your account will be deleted on $date. To cancel, just sign in before then.';
  }

  @override
  String get deletionCancelled => 'Your account deletion request was cancelled.';

  @override
  String get support => 'Support';

  @override
  String get supportIntro =>
      'Found a bug or need help with your account? Write to us and the team answers here; you get a notification when they do.';

  @override
  String get supportNew => 'New request';

  @override
  String get supportNone => 'No requests yet.';

  @override
  String get supportSent => 'Request received. We\'ll let you know when there\'s an answer.';

  @override
  String get supportCategory => 'Topic';

  @override
  String get supportSubject => 'Subject';

  @override
  String get supportMessage => 'What happened? Step by step if you can.';

  @override
  String get supportTooShort => 'Subject needs at least 3 characters, the message at least 10.';

  @override
  String get supportBug => 'Bug';

  @override
  String get supportAccount => 'Account';

  @override
  String get supportAbuse => 'Abuse / report a player';

  @override
  String get supportFeedback => 'Suggestion';

  @override
  String get supportOther => 'Other';

  @override
  String get supportOpen => 'Open';

  @override
  String get supportAnswered => 'Answered';

  @override
  String get supportClosed => 'Closed';

  @override
  String get supportReply => 'Team reply';

  @override
  String get nkSupport => 'Support reply';

  @override
  String get errServerBusy => 'The server is busy right now. Try again in a few seconds.';

  @override
  String get notifTabInbox => 'Notifications';

  @override
  String get notifTabRewards => 'Rewards';

  @override
  String get rewardsNone => 'No rewards yet. Everything you gather or win in fights is listed here (kept for 30 days).';

  @override
  String rewardFromActivity(String name) {
    return 'Gathering: $name';
  }

  @override
  String rewardFromHunt(String name) {
    return 'Hunt: $name';
  }

  @override
  String rewardFromDungeon(String name) {
    return 'Dungeon: $name';
  }

  @override
  String rewardLevelUp(int level) {
    return 'Level $level!';
  }

  @override
  String get levelMax => 'MAX';

  @override
  String get promotionAvailable => 'PROMOTION AVAILABLE';

  @override
  String get promotionChoose => 'Choose';

  @override
  String promotionConfirm(String name) {
    return 'Become $name?';
  }

  @override
  String get promotionBody => 'This choice shapes your class for good.';

  @override
  String nextMilestone(int level, String what) {
    return 'Next: Lv $level · $what';
  }

  @override
  String get msPromotion => 'Class promotion';

  @override
  String msAwakening(String name) {
    return 'Awakening: $name';
  }

  @override
  String msUltimate(String name) {
    return 'Ultimate: $name';
  }

  @override
  String get msCapstone => 'Talent capstone';

  @override
  String get msMastery => 'Class mastery';

  @override
  String msTitle(String name) {
    return 'Title: $name';
  }

  @override
  String estMinutes(int n) {
    return '~$n min';
  }

  @override
  String estHours(String h) {
    return '~$h h';
  }

  @override
  String overflowHeld(int xp) {
    return '$xp XP is waiting for your choice';
  }

  @override
  String attributePointsBtn(int n) {
    return 'Attributes · $n points';
  }

  @override
  String talentPointsBtn(int n) {
    return 'Talents · $n points';
  }

  @override
  String get talentsTitle => 'Talents';

  @override
  String capstoneLocked(int level, int points) {
    return 'Needs level $level and $points points in this tree';
  }

  @override
  String masteryUnlocked(String title) {
    return 'Class mastery: $title';
  }

  @override
  String get errPromotionUnavailable => 'That promotion is not available.';

  @override
  String get errNotEnoughPoints => 'Not enough points.';

  @override
  String get errTalentMaxed => 'That talent is already at max rank.';

  @override
  String get errCapstoneLocked => 'The capstone needs its level and a full tree.';

  @override
  String get errSpecRequirement => 'Your specialization cannot use this item.';

  @override
  String get errAttributeRequirement => 'Your attributes are too low for this item.';

  @override
  String get errBusyCombat => 'Not while in combat, AFK or a dungeon.';

  @override
  String skillGuard(int pct, int turns) {
    return '+$pct% armor for $turns turns';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '+$pct% attack for $turns turns';
  }

  @override
  String get classRanger => 'Ranger';

  @override
  String get classCleric => 'Cleric';

  @override
  String get classMonk => 'Monk';

  @override
  String get classPaladin => 'Paladin';

  @override
  String get classDruid => 'Druid';

  @override
  String get classBard => 'Bard';

  @override
  String get classShaman => 'Shaman';

  @override
  String get errEquipmentPermission => 'Your class cannot use this equipment type.';

  @override
  String get classResource => 'Resource';

  @override
  String get classFocus => 'Attribute focus';

  @override
  String classSecondary(String names) {
    return 'Secondary: $names';
  }

  @override
  String get resStartsEmpty => 'Starts empty';

  @override
  String get resStartsFull => 'Starts full';

  @override
  String resStartsAt(int pct) {
    return 'Starts at $pct%';
  }

  @override
  String resPerHit(int n) {
    return '+$n per hit';
  }

  @override
  String resPerTurn(int n) {
    return '+$n per turn';
  }

  @override
  String baseAttributesHint(int total) {
    return 'Base attributes · 1 point every 2 levels, $total in total';
  }

  @override
  String get classEquipment => 'Equipment';

  @override
  String get firstSkills => 'First skills';

  @override
  String get classPath => 'Class path';

  @override
  String get classPathHint => 'Tap a class to see its skills.';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Lv $awakening Awakening · Lv $ultimate Ultimate · Lv $capstone talent capstone · Lv $mastery class mastery';
  }

  @override
  String get skillTagActive => 'active';

  @override
  String get skillTagPassive => 'passive';

  @override
  String get skillTagAwakening => 'awakening';

  @override
  String get skillTagUltimate => 'ultimate';

  @override
  String skillCostCd(int cost, int cd) {
    return 'Cost $cost · cooldown $cd turns';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Lv $level mastery: $title';
  }

  @override
  String get classInfoFailed => 'Class details couldn\'t load.';

  @override
  String xpShare(int pct) {
    return 'XP $pct%';
  }

  @override
  String get dungeonsSolo => 'Dungeons';

  @override
  String get raidsTab => 'Raids';

  @override
  String raidGroup(int min, int max) {
    return 'Party of $min–$max';
  }

  @override
  String raidBossGear(int count) {
    return 'Boss: $count gear pieces guaranteed';
  }

  @override
  String get raidOpen => 'Open raid';

  @override
  String get raidNeedParty =>
      'Raids are for a party. Queue below for a random one, join a listing on the board, or advertise your own.';

  @override
  String get raidLeaderOnly => 'Your party leader opens raids.';

  @override
  String get raidLobby => 'Raid lobby';

  @override
  String get raidJoin => 'Join';

  @override
  String get raidLeave => 'Leave';

  @override
  String get raidCancel => 'Call off';

  @override
  String get raidCancelBody => 'The lobby closes for everyone.';

  @override
  String get raidStart => 'Start raid';

  @override
  String get raidStartBody => 'Nobody can join or leave once it starts.';

  @override
  String raidMembers(int count, int max) {
    return 'Members ($count/$max)';
  }

  @override
  String raidNeedMore(int min) {
    return 'At least $min members to start';
  }

  @override
  String get raidWaitHint => 'Rewards open when the raid ends; every member collects their own.';

  @override
  String get raidDone => 'Raid over — collect your rewards';

  @override
  String get raidLocked => 'Cleared this week';

  @override
  String get raidCollected => 'collected';

  @override
  String get errRaidNeedsParty => 'Join a party first.';

  @override
  String get errRaidExists => 'Your party already has an open raid.';

  @override
  String get errAlreadyInRaid => 'You are already in a raid.';

  @override
  String get errRaidNotOpen => 'This raid is no longer open.';

  @override
  String get errRaidNotEnough => 'Not enough members to start.';

  @override
  String get errRaidWeeklyLocked => 'You already cleared this raid this week.';

  @override
  String errCharacterBusy(String action) {
    return 'Finish what your character is doing first: $action.';
  }

  @override
  String get busyGathering => 'gathering';

  @override
  String get busyHunt => 'hunt';

  @override
  String get busyDungeon => 'dungeon';

  @override
  String get busyRaid => 'raid';

  @override
  String get busyTravel => 'travel';

  @override
  String get busyCrafting => 'crafting';

  @override
  String get errCraftActive => 'A crafting job is already running.';

  @override
  String get errCraftTooLong => 'That job would take more than 6 hours.';

  @override
  String craftEach(String time) {
    return '$time each';
  }

  @override
  String craftTotal(String time) {
    return 'Total $time';
  }

  @override
  String craftRunning(String item, int count) {
    return 'Crafting $item ×$count';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count ready';
  }

  @override
  String get craftDone => 'All done — collect';

  @override
  String get craftCollectEarlyTitle => 'Collect now?';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return 'The $done finished items are delivered; the materials and gold of the other $rest come back.';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '$made made, $refunded refunded';
  }

  @override
  String get craftBusyHint => 'Crafting takes time; your character does nothing else meanwhile.';

  @override
  String get eventHowTitle => 'How to take part';

  @override
  String get eventRewardsTitle => 'Rewards';

  @override
  String eventStepEnlist(int level) {
    return 'Enlist from 24 hours before it starts until it ends (level $level+).';
  }

  @override
  String eventStepRegion(String region) {
    return 'Be in $region when you enlist.';
  }

  @override
  String eventStepPower(int power) {
    return 'Your power is frozen when you enlist: attack × 2 + defense + HP ÷ 10. Yours now: $power.';
  }

  @override
  String eventStepScore(int tier) {
    return 'Score = power × minutes attended, from the start (or from when you enlist, if later). It stops at tier $tier.';
  }

  @override
  String get eventStepGoal =>
      'If everyone together reaches the goal, rewards go up to the top tier; if not, they stop at tier 2.';

  @override
  String get eventStepPaid =>
      'Rewards are paid automatically when it closes. Enlisting does not stop your other actions.';

  @override
  String eventTier(int tier, String score) {
    return 'Tier $tier · $score+ points';
  }

  @override
  String eventTierTokens(int count) {
    return '+$count event tokens';
  }

  @override
  String eventEstimate(String score, int tier) {
    return 'Join now: about $score points → tier $tier';
  }

  @override
  String eventEstimateNone(String score) {
    return 'Join now: about $score points — not enough for a tier';
  }

  @override
  String seasonalHow(int count) {
    return 'While the event runs, claim $count event tokens once a day (the day starts at 00:00 UTC).';
  }

  @override
  String seasonalBenefit(int count) {
    return 'Tokens stay on your account (balance: $count).';
  }

  @override
  String get eventDetails => 'Details';

  @override
  String get createCharacterTitle => 'Create character';

  @override
  String get classTraits => 'Class traits';

  @override
  String classResourceLine(String name) {
    return 'Resource: $name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return 'Main attribute: $primary · second: $secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return 'Main attribute: $primary';
  }

  @override
  String get classDifficulty => 'Difficulty';

  @override
  String get difficultyEasy => 'Easy — forgiving, safe on AFK hunts';

  @override
  String get difficultyMedium => 'Medium — balanced; skill order matters';

  @override
  String get difficultyHard => 'Hard — fragile; needs the right gear';

  @override
  String get skillsTapHint => 'tap for details';

  @override
  String get chooseRace => 'Choose your race';

  @override
  String get raceCosmeticHint => 'Your race sets your look; it does not change attributes.';

  @override
  String get raceHumanDesc => 'People of the cities';

  @override
  String get raceElfDesc => 'Children of the old woods';

  @override
  String get raceDwarfDesc => 'Smiths of the mountains';

  @override
  String get raceOrcDesc => 'Warriors of the steppe';

  @override
  String continueAs(String name) {
    return 'Continue as $name';
  }

  @override
  String get nameYourCharacter => 'Name your character';

  @override
  String get raceTroll => 'Troll';

  @override
  String get raceUndead => 'Undead';

  @override
  String get raceDarkElf => 'Dark Elf';

  @override
  String get raceDrakari => 'Drakari';

  @override
  String get raceTauren => 'Tauren';

  @override
  String get raceGnome => 'Gnome';

  @override
  String get raceTrollDesc => 'Hunters of the jungle isles';

  @override
  String get raceUndeadDesc => 'Risen from the grave';

  @override
  String get raceDarkElfDesc => 'Children of the deep dark';

  @override
  String get raceDrakariDesc => 'Kin of dragon blood';

  @override
  String get raceTaurenDesc => 'Guardians of the plains';

  @override
  String get raceGnomeDesc => 'Clever tinkerers';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '+$atk% attack and +$def% defense for $minutes min';
  }

  @override
  String fedToast(int atk, int def) {
    return 'Fed: +$atk% attack, +$def% defense';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return 'Fed · +$atk% attack · +$def% defense · $time left';
  }

  @override
  String huntReach(String regions) {
    return 'Hunting from here: $regions';
  }

  @override
  String get huntNoMobsHere => 'There is nothing to hunt here. Travel to a hunting ground first.';

  @override
  String get huntOtherRegions => 'Monsters in other regions';

  @override
  String get huntTravelHint => 'Travel to a monster\'s region, or a region linked to it by road, to hunt it.';

  @override
  String get errMobRegion => 'This monster lives in another region. Travel there first.';

  @override
  String restoresHpPct(int pct) {
    return 'Restores $pct% of max HP';
  }

  @override
  String get statEvasion => '% Evasion';

  @override
  String get statBlock => '% Block';

  @override
  String get statCritDmg => '% Crit damage';

  @override
  String get statLifesteal => '% Lifesteal';

  @override
  String get statDmgReduction => '% Damage reduction';

  @override
  String get classJewelry => 'Jewelry';

  @override
  String get mapHint =>
      'You can travel straight to any region your level allows; the time follows the shortest road. Drag and pinch the map; tap a region.';

  @override
  String get mapNoRoad => 'No direct road leads here. Travel through the glowing regions next to you.';

  @override
  String get youAreHere => 'You are here.';

  @override
  String travelTime(String time) {
    return 'Travel time: $time';
  }

  @override
  String xpPerKill(int xp) {
    return '$xp XP';
  }

  @override
  String get deleteCharacter => 'Delete character';

  @override
  String deleteCharacterBody(String name) {
    return '$name will be deleted right away, with all items and gold. This cannot be undone. Type the character\'s name to confirm.';
  }

  @override
  String get characterNameConfirm => 'Character name';

  @override
  String characterDeletesAt(String when) {
    return 'Deleted on $when';
  }

  @override
  String get cancelCharacterDelete => 'Cancel deletion';

  @override
  String characterDeleted(String name) {
    return '$name was deleted.';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name will not be deleted.';
  }

  @override
  String get errNameMismatch => 'The name does not match.';

  @override
  String attrGainPerPoint(String gains) {
    return 'Per point: $gains';
  }

  @override
  String attrGainPending(String gains) {
    return 'Adding: $gains';
  }

  @override
  String get attrNoEffect => 'No effect for your class';

  @override
  String get attrInputHint => 'Type how many points each attribute gets';

  @override
  String get gainHp => 'HP';

  @override
  String get gainMp => 'MP';

  @override
  String get gainHpRegen => 'HP/s';

  @override
  String get huntSupplies => 'Supplies to take';

  @override
  String get huntSuppliesHint =>
      'Potions are drunk between fights when HP drops under half, food is eaten when its bonus runs out, mana potions when mana runs low. What is left comes back to your bag.';

  @override
  String get supplyHp => 'Healing potions';

  @override
  String get supplyMp => 'Mana potions';

  @override
  String get supplyMpHint => 'For classes using mana, energy, focus, nature energy or spirit';

  @override
  String get supplyFood => 'Food';

  @override
  String get huntNoSupplies => 'No potions or food in your bag: buy them at the shop or craft them.';

  @override
  String restoresMpPct(int pct) {
    return 'Refills $pct% of your mana';
  }

  @override
  String get huntSupplyOnly => 'Used only as a hunt supply';

  @override
  String get errBadSupply => 'That item cannot be taken on a hunt.';

  @override
  String get errNotEnoughSupplies => 'Your bag holds fewer of that supply.';

  @override
  String get errHuntOnly => 'Take it on a hunt as a supply.';

  @override
  String get huntFled => 'Your HP fell under 5%: you broke off the fight';

  @override
  String suppliesUsed(String list) {
    return 'Used: $list';
  }

  @override
  String get compareTitle => 'Compared with what you wear';

  @override
  String compareWorn(String name) {
    return 'Worn: $name';
  }

  @override
  String get compareEmpty => 'That slot is empty: equipping adds all of these.';

  @override
  String get compareColStat => 'Stat';

  @override
  String get compareColWorn => 'Worn';

  @override
  String get compareColThis => 'This';

  @override
  String get compareColDiff => 'Change';

  @override
  String get compareSame => 'Same stats as what you wear.';

  @override
  String get wearOk => 'Your class can wear this';

  @override
  String get wearClass => 'Your class can\'t wear this';

  @override
  String wearLevel(int level) {
    return 'Needs level $level';
  }

  @override
  String get wearAttribute => 'Your attributes are too low';

  @override
  String get wearLegendClass => 'not for your class';

  @override
  String get wearLegendLater => 'not yet (level or attributes)';

  @override
  String get enchantAction => 'Enchant';

  @override
  String get enchantPick => 'Choose the gear to enchant';

  @override
  String get enchantNone => 'None of your gear fits this scroll.';

  @override
  String enchantFits(String groups) {
    return 'Goes on: $groups';
  }

  @override
  String get enchantGroupWeapon => 'weapons';

  @override
  String get enchantGroupArmor => 'armor';

  @override
  String get enchantGroupJewelry => 'jewelry';

  @override
  String get enchantGroupAny => 'any gear';

  @override
  String get enchantReplaceTitle => 'Replace the enchant?';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old is lost and $next takes its place.';
  }

  @override
  String enchantDone(String item) {
    return '$item enchanted';
  }

  @override
  String enchantOn(String name) {
    return 'Enchant: $name';
  }

  @override
  String get wornTag => 'worn';

  @override
  String get errNotEnchant => 'That is not an enchanting scroll.';

  @override
  String get errEnchantSlot => 'The enchant doesn\'t fit that gear.';

  @override
  String get statGatherMining => 'Mining yield';

  @override
  String get statGatherHerbalism => 'Herbalism yield';

  @override
  String get statGatherWoodcutting => 'Woodcutting yield';

  @override
  String get statGatherFishing => 'Fishing yield';

  @override
  String get statGatherHunting => 'Hunting yield';

  @override
  String get statCraftSpeed => 'Crafting speed';

  @override
  String diffXp(int pct) {
    return '$pct% XP';
  }

  @override
  String get diffNoXp => 'no XP';

  @override
  String get diffLegend => 'Recipe colors: the XP a craft teaches at your skill';

  @override
  String get recipeGroupEnchantWeapon => 'Weapon enchants';

  @override
  String get recipeGroupEnchantArmor => 'Armor enchants';

  @override
  String get recipeGroupEnchantJewelry => 'Jewelry enchants';

  @override
  String get recipeGroupEnchantProfession => 'Profession enchants';

  @override
  String get partyHint =>
      'Members doing the same work at the same time (the same gathering kind, profession, hunt or dungeon) each earn +3% XP per member working: five together earn +15%. Nobody gives XP away.';

  @override
  String get partyChat => 'Party chat';

  @override
  String partyShared(String names, int kept, int earned) {
    return 'Party bonus with $names: $earned → $kept XP';
  }

  @override
  String evRaidMatched(String raid) {
    return 'A party was found for $raid. See Raids.';
  }

  @override
  String get nkRaidMatched => 'Raid party found';

  @override
  String get workHunt => 'Hunting';

  @override
  String get workDungeon => 'Dungeon';

  @override
  String get raidFind => 'Find a party';

  @override
  String get raidQueueJoin => 'Queue';

  @override
  String get raidQueueLeave => 'Leave queue';

  @override
  String raidQueued(String raid, int waiting) {
    return 'Queued for $raid · $waiting waiting';
  }

  @override
  String get raidQueueHint => 'A random party forms once 5 are waiting, or after 5 minutes with at least 2.';

  @override
  String get raidAdvertise => 'Advertise';

  @override
  String get raidUnlist => 'Remove listing';

  @override
  String raidListed(String raid) {
    return 'Your party is on the board for $raid';
  }

  @override
  String get raidBoard => 'Listing board';

  @override
  String get raidBoardEmpty => 'No open listings.';

  @override
  String get raidJoinParty => 'Join party';

  @override
  String get raidPickRaid => 'Raid';

  @override
  String get errListingGone => 'That listing is gone or full.';

  @override
  String get workGathering => 'Gathering';

  @override
  String get workAdventure => 'Adventure';

  @override
  String get workRaid => 'Raid';

  @override
  String get workTravel => 'Travelling';

  @override
  String get workCrafting => 'Crafting';

  @override
  String get workIdle => 'Idle';

  @override
  String get errDungeonRegion => 'Travel to this place on the map first.';

  @override
  String get dungeonsElsewhere =>
      'More dungeons and raids open as you travel: each is entered from its own region on the map.';

  @override
  String get dungeonNoneHere => 'No dungeon in reach here. Travel on the map to a dungeon region.';

  @override
  String dungeonXp(int xp) {
    return '≈ $xp XP for you';
  }

  @override
  String appVersionLine(String version) {
    return 'Version $version';
  }

  @override
  String get updateButton => 'Update';

  @override
  String updateAvailable(String version) {
    return 'New version ready: $version';
  }

  @override
  String get updateLatest => 'You have the latest version';

  @override
  String updateConfirm(String version, int size) {
    return 'Download and install version $version ($size MB)?';
  }

  @override
  String get updateDownloading => 'Downloading…';

  @override
  String get updatePermission => 'Allow installs from MeNa World, then tap Update again.';

  @override
  String get updateFailed => 'Could not get the update. Check your connection.';

  @override
  String get fightReport => 'Fight report';

  @override
  String get endTimeUp => 'Time ran out: it went on to the planned end.';

  @override
  String get endEarly => 'Collected early: stopped before the planned end.';

  @override
  String endFled(String mob) {
    return 'HP fell below 5% against $mob: you broke off and came back.';
  }

  @override
  String endDefeat(String mob) {
    return 'You fell to $mob.';
  }

  @override
  String endPartyDefeat(String mob) {
    return 'The party fell to $mob.';
  }

  @override
  String endTimeout(String mob) {
    return 'The fight with $mob went past the turn limit: it stopped.';
  }

  @override
  String get endCleared => 'Cleared: every wave and the boss fell.';

  @override
  String get repCampHint =>
      'Under 50% HP the hunter camps back to full before the next fight, so a hunt rarely drops to 5%.';

  @override
  String get repYou => 'You';

  @override
  String get repLevel => 'Level';

  @override
  String get repHp => 'HP';

  @override
  String get repAttack => 'Attack';

  @override
  String get repDefense => 'Defense';

  @override
  String get repEvasion => 'Evasion';

  @override
  String get repHitChance => 'Hit chance';

  @override
  String get repCrit => 'Crit';

  @override
  String repFights(int n) {
    return '$n fights';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '$hits/$swings hits ($pct%), avg $avg, max $max';
  }

  @override
  String get repYourSwings => 'Your swings';

  @override
  String get repPartySwings => 'The party\'s swings';

  @override
  String get repEnemySwings => 'Enemy swings';

  @override
  String repLowest(int pct) {
    return 'Lowest HP in a fight: $pct%';
  }

  @override
  String repRest(String time) {
    return 'Camped: $time';
  }

  @override
  String get restart => 'Restart';

  @override
  String get restartTitle => 'Collect and start again?';

  @override
  String get restartBody => 'What came in so far goes to your bag, and the same run starts again for the same time.';

  @override
  String get dailyTitle => 'Daily login reward';

  @override
  String get dailyBody =>
      'Claim once a day on any character. Seven days in a row pay the most; a missed day starts the cycle again.';

  @override
  String get dailyClaim => 'Claim';

  @override
  String get dailyTomorrow => 'Come back tomorrow';

  @override
  String dailyDay(int n) {
    return 'Day $n';
  }

  @override
  String get dailyOpen => 'Daily reward';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name was muted for $min minutes for swearing.';
  }

  @override
  String get storyTitle => 'Story';

  @override
  String get storyOpen => 'Story mode';

  @override
  String get storyComingSoon => 'Story mode is coming soon. Chapters will appear here as they are written.';

  @override
  String storyLocked(int level) {
    return 'Opens at level $level';
  }

  @override
  String get bankTitle => 'Bank';

  @override
  String get bankHint =>
      'Shared by all your characters. Tap an item in your bag to deposit it, or a banked one to take it. Bound items stay with their character.';

  @override
  String get bankEmpty => 'The bank is empty.';

  @override
  String get bankDeposit => 'Deposit';

  @override
  String get bankWithdraw => 'Take';

  @override
  String get errBankFull => 'The bank is full.';

  @override
  String get errItemBound => 'Bound items stay with their character.';

  @override
  String get errItemEquipped => 'Take the item off first.';

  @override
  String get tabTrade => 'Trade';

  @override
  String get tradeHint =>
      'Trade items with another player. Both put items up and confirm; each side pays a tax of 8% of the average market price of what it gives, per item.';

  @override
  String get tradeWith => 'Player name';

  @override
  String get tradeOpen => 'Start trade';

  @override
  String tradeTitle(String name) {
    return 'Trade with $name';
  }

  @override
  String get tradeYours => 'You give';

  @override
  String tradeTheirs(String name) {
    return '$name gives';
  }

  @override
  String get tradeNothing => 'Nothing yet.';

  @override
  String get tradeAdd => 'Add item';

  @override
  String get tradeRemove => 'Take back';

  @override
  String get tradeTax => 'Your tax';

  @override
  String get tradeReady => 'Ready';

  @override
  String get tradeWaiting => 'Not ready';

  @override
  String get tradeConfirm => 'Confirm trade';

  @override
  String get tradeWaitingThem => 'Waiting for the other player';

  @override
  String get tradeCancel => 'Cancel trade';

  @override
  String get tradeCancelBody => 'Nothing changes hands.';

  @override
  String get tradeDone => 'Trade complete.';

  @override
  String get errTradeBusy => 'One of you is already trading.';

  @override
  String get errTradeClosed => 'That trade is over.';

  @override
  String get errTradeItem => 'That item cannot be traded now.';

  @override
  String get errTradeEmpty => 'Put something up first.';

  @override
  String get errTradeTooMany => 'Too many items in the trade.';

  @override
  String get errCraftLevel => 'Your character\'s level is below this item\'s level.';

  @override
  String get gatherNoneHere => 'Your resources grow on other maps: travel to one shown beside them.';

  @override
  String get errGatherRegion => 'That resource grows on another map.';

  @override
  String get helpersOpen => 'Helpers';

  @override
  String get helpersTitle => 'Helpers';

  @override
  String get helpersHint =>
      'Hire helpers to gather and craft beside you. They work at half your speed and earn half the XP, which goes to your profession; each kind has its own bonus.';

  @override
  String get helpersGather => 'Gathering';

  @override
  String get helpersCraft => 'Crafting';

  @override
  String get helperHire => 'Hire';

  @override
  String helperHireBody(String name, String price) {
    return 'Hire the $name for $price?';
  }

  @override
  String get helperSend => 'Send to work';

  @override
  String get helperCancel => 'Call back';

  @override
  String get helperCancelBody => 'A gathering run brings nothing back; a craft returns its materials and gold.';

  @override
  String get helperDone => 'Done — collect it';

  @override
  String helperWorking(String time) {
    return 'Working, $time left';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return '$items items, +$xp $profession XP';
  }

  @override
  String get helperNoRecipes => 'No recipe at your level in this profession yet.';

  @override
  String get helperPickRecipe => 'Pick a recipe';

  @override
  String helperBonusYield(int pct) {
    return '+$pct% yield';
  }

  @override
  String helperBonusSpeed(int pct) {
    return '$pct% faster';
  }

  @override
  String helperBonusXp(int pct) {
    return '+$pct% XP';
  }

  @override
  String helperBonusGold(int pct) {
    return '$pct% less gold cost';
  }

  @override
  String helperBonusDouble(int pct) {
    return '$pct% chance of a double output';
  }

  @override
  String get helperMiner => 'Miner';

  @override
  String get helperWoodsman => 'Woodsman';

  @override
  String get helperAngler => 'Angler';

  @override
  String get helperHerbalist => 'Herbalist';

  @override
  String get helperTrapper => 'Trapper';

  @override
  String get helperSmith => 'Smith';

  @override
  String get helperTanner => 'Tanner';

  @override
  String get helperWeaver => 'Weaver';

  @override
  String get helperCarver => 'Carver';

  @override
  String get helperApothecary => 'Apothecary';

  @override
  String get helperCook => 'Cook';

  @override
  String get helperSmelter => 'Smelter';

  @override
  String get helperScribe => 'Scribe';

  @override
  String get helperGemcutter => 'Gemcutter';

  @override
  String get errHelperBusy => 'That helper is already working.';

  @override
  String get errHelperNotReady => 'The helper is not done yet.';
}
