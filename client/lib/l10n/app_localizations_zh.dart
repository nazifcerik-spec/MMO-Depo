// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Chinese (`zh`).
class AppLocalizationsZh extends AppLocalizations {
  AppLocalizationsZh([String locale = 'zh']) : super(locale);

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
  String get classWarrior => '战士';

  @override
  String get classMage => '法师';

  @override
  String get classRogue => '盗贼';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Lvl $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return '$level级';
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
  String get close => '关闭';

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
  String get errWrongRegion => '这只怪物在其他区域。请先前往那里。';

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
  String get levelMax => '满级';

  @override
  String get promotionAvailable => '可以晋升';

  @override
  String get promotionChoose => '选择';

  @override
  String promotionConfirm(String name) {
    return '成为$name？';
  }

  @override
  String get promotionBody => '此选择将永久决定你的职业。';

  @override
  String nextMilestone(int level, String what) {
    return '下一目标：$level级 · $what';
  }

  @override
  String get msPromotion => '职业晋升';

  @override
  String msAwakening(String name) {
    return '觉醒：$name';
  }

  @override
  String msUltimate(String name) {
    return '终极技：$name';
  }

  @override
  String get msCapstone => '终极天赋';

  @override
  String get msMastery => '职业精通';

  @override
  String msTitle(String name) {
    return '称号：$name';
  }

  @override
  String estMinutes(int n) {
    return '约$n分钟';
  }

  @override
  String estHours(String h) {
    return '约$h小时';
  }

  @override
  String overflowHeld(int xp) {
    return '$xp经验等待你的选择';
  }

  @override
  String attributePointsBtn(int n) {
    return '属性 · $n点';
  }

  @override
  String talentPointsBtn(int n) {
    return '天赋 · $n点';
  }

  @override
  String get talentsTitle => '天赋';

  @override
  String capstoneLocked(int level, int points) {
    return '需要$level级且此天赋树投入$points点';
  }

  @override
  String masteryUnlocked(String title) {
    return '职业精通：$title';
  }

  @override
  String get errPromotionUnavailable => '无法进行该晋升。';

  @override
  String get errNotEnoughPoints => '点数不足。';

  @override
  String get errTalentMaxed => '该天赋已满级。';

  @override
  String get errCapstoneLocked => '终极天赋需要等级和满点天赋树。';

  @override
  String get errSpecRequirement => '你的专精无法使用此物品。';

  @override
  String get errAttributeRequirement => '你的属性不足以使用此物品。';

  @override
  String get errBusyCombat => '战斗、挂机或地下城中无法进行。';

  @override
  String skillGuard(int pct, int turns) {
    return '$turns回合护甲 +$pct%';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '$turns回合攻击 +$pct%';
  }

  @override
  String get classRanger => '游侠';

  @override
  String get classCleric => '牧师';

  @override
  String get classMonk => '武僧';

  @override
  String get classPaladin => '圣骑士';

  @override
  String get classDruid => '德鲁伊';

  @override
  String get classBard => '吟游诗人';

  @override
  String get classShaman => '萨满';

  @override
  String get errEquipmentPermission => '你的职业无法使用此类装备。';

  @override
  String get classResource => '资源';

  @override
  String get classFocus => '属性侧重';

  @override
  String classSecondary(String names) {
    return '次要：$names';
  }

  @override
  String get resStartsEmpty => '开局为空';

  @override
  String get resStartsFull => '开局全满';

  @override
  String resStartsAt(int pct) {
    return '开局$pct%';
  }

  @override
  String resPerHit(int n) {
    return '每次命中+$n';
  }

  @override
  String resPerTurn(int n) {
    return '每回合+$n';
  }

  @override
  String baseAttributesHint(int total) {
    return '基础属性 · 每2级1点，共$total点';
  }

  @override
  String get classEquipment => '装备';

  @override
  String get firstSkills => '初始技能';

  @override
  String get classPath => '职业路线';

  @override
  String get classPathHint => '点击职业查看技能。';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return '$awakening级 觉醒 · $ultimate级 终极技 · $capstone级 终极天赋 · $mastery级 职业精通';
  }

  @override
  String get skillTagActive => '主动';

  @override
  String get skillTagPassive => '被动';

  @override
  String get skillTagAwakening => '觉醒';

  @override
  String get skillTagUltimate => '终极';

  @override
  String skillCostCd(int cost, int cd) {
    return '消耗$cost · 冷却$cd回合';
  }

  @override
  String masteryAt(int level, String title) {
    return '$level级精通：$title';
  }

  @override
  String get classInfoFailed => '无法加载职业信息。';

  @override
  String xpShare(int pct) {
    return '经验 $pct%';
  }

  @override
  String get dungeonsSolo => '地下城';

  @override
  String get raidsTab => '团队副本';

  @override
  String raidGroup(int min, int max) {
    return '$min–$max 人队伍';
  }

  @override
  String raidBossGear(int count) {
    return '首领：必掉 $count 件装备';
  }

  @override
  String get raidOpen => '开启团队副本';

  @override
  String get raidNeedParty => '团队副本需要队伍。请在 社交 → 队伍 中组队。';

  @override
  String get raidLeaderOnly => '由队长开启团队副本。';

  @override
  String get raidLobby => '团队副本大厅';

  @override
  String get raidJoin => '加入';

  @override
  String get raidLeave => '离开';

  @override
  String get raidCancel => '取消';

  @override
  String get raidCancelBody => '大厅将对所有人关闭。';

  @override
  String get raidStart => '开始团队副本';

  @override
  String get raidStartBody => '开始后无人可以加入或离开。';

  @override
  String raidMembers(int count, int max) {
    return '成员（$count/$max）';
  }

  @override
  String raidNeedMore(int min) {
    return '至少需要 $min 人才能开始';
  }

  @override
  String get raidWaitHint => '团队副本结束后开放奖励，每位成员各自领取。';

  @override
  String get raidDone => '团队副本结束——领取奖励';

  @override
  String get raidLocked => '本周已通关';

  @override
  String get raidCollected => '已领取';

  @override
  String get errRaidNeedsParty => '请先加入队伍。';

  @override
  String get errRaidExists => '你的队伍已有开启的团队副本。';

  @override
  String get errAlreadyInRaid => '你已在团队副本中。';

  @override
  String get errRaidNotOpen => '该团队副本已不再开放。';

  @override
  String get errRaidNotEnough => '人数不足，无法开始。';

  @override
  String get errRaidWeeklyLocked => '你本周已通关该团队副本。';

  @override
  String errCharacterBusy(String action) {
    return '请先完成角色当前的行动：$action。';
  }

  @override
  String get busyGathering => '采集';

  @override
  String get busyHunt => '狩猎';

  @override
  String get busyDungeon => '地下城';

  @override
  String get busyRaid => '团队副本';

  @override
  String get busyTravel => '旅行';

  @override
  String get busyCrafting => '制作';

  @override
  String get errCraftActive => '已有制作正在进行。';

  @override
  String get errCraftTooLong => '该制作将超过 6 小时。';

  @override
  String craftEach(String time) {
    return '每件 $time';
  }

  @override
  String craftTotal(String time) {
    return '共 $time';
  }

  @override
  String craftRunning(String item, int count) {
    return '正在制作 $item ×$count';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count 已完成';
  }

  @override
  String get craftDone => '全部完成——领取';

  @override
  String get craftCollectEarlyTitle => '现在领取？';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return '交付已完成的 $done 件；其余 $rest 件的材料和金币将退还。';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '制作 $made 件，退还 $refunded 件';
  }

  @override
  String get craftBusyHint => '制作需要时间，期间角色无法进行其他行动。';

  @override
  String get eventHowTitle => '如何参与';

  @override
  String get eventRewardsTitle => '奖励';

  @override
  String eventStepEnlist(int level) {
    return '可在开始前 24 小时至结束前报名（$level 级以上）。';
  }

  @override
  String eventStepRegion(String region) {
    return '报名时需位于 $region。';
  }

  @override
  String eventStepPower(int power) {
    return '报名时冻结战力：攻击 × 2 + 防御 + 生命 ÷ 10。你当前的战力：$power。';
  }

  @override
  String eventStepScore(int tier) {
    return '积分 = 战力 × 参与分钟数，从开始（或更晚报名时从报名起）计算，最高到第 $tier 档。';
  }

  @override
  String get eventStepGoal => '若全服共同达成目标，奖励可达最高档；否则最高第 2 档。';

  @override
  String get eventStepPaid => '活动结束时自动发放奖励。报名不影响你的其他行动。';

  @override
  String eventTier(int tier, String score) {
    return '第 $tier 档 · $score+ 积分';
  }

  @override
  String eventTierTokens(int count) {
    return '+$count 活动代币';
  }

  @override
  String eventEstimate(String score, int tier) {
    return '现在参加：约 $score 积分 → 第 $tier 档';
  }

  @override
  String eventEstimateNone(String score) {
    return '现在参加：约 $score 积分——不足一档';
  }

  @override
  String seasonalHow(int count) {
    return '活动期间每天可领取一次 $count 个活动代币（每天从 UTC 00:00 开始）。';
  }

  @override
  String seasonalBenefit(int count) {
    return '代币保存在你的账号中（余额：$count）。';
  }

  @override
  String get eventDetails => '详情';

  @override
  String get createCharacterTitle => '创建角色';

  @override
  String get classTraits => '职业特性';

  @override
  String classResourceLine(String name) {
    return '资源：$name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return '主属性：$primary · 次属性：$secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return '主属性：$primary';
  }

  @override
  String get classDifficulty => '难度';

  @override
  String get difficultyEasy => '简单——容错高，挂机狩猎安全';

  @override
  String get difficultyMedium => '中等——均衡；技能顺序很重要';

  @override
  String get difficultyHard => '困难——脆弱；需要合适的装备';

  @override
  String get skillsTapHint => '点击查看说明';

  @override
  String get chooseRace => '选择种族';

  @override
  String get raceCosmeticHint => '种族决定外观，不影响属性。';

  @override
  String get raceHumanDesc => '城邦之民';

  @override
  String get raceElfDesc => '远古森林之子';

  @override
  String get raceDwarfDesc => '群山的铁匠';

  @override
  String get raceOrcDesc => '草原的战士';

  @override
  String continueAs(String name) {
    return '以$name继续';
  }

  @override
  String get nameYourCharacter => '为角色命名';

  @override
  String get raceTroll => '巨魔';

  @override
  String get raceUndead => '亡灵';

  @override
  String get raceDarkElf => '黑暗精灵';

  @override
  String get raceDrakari => '龙裔';

  @override
  String get raceTauren => '牛头人';

  @override
  String get raceGnome => '侏儒';

  @override
  String get raceTrollDesc => '丛林群岛的猎手';

  @override
  String get raceUndeadDesc => '从坟墓中复生';

  @override
  String get raceDarkElfDesc => '幽暗地域之子';

  @override
  String get raceDrakariDesc => '流淌龙血的一族';

  @override
  String get raceTaurenDesc => '平原的守护者';

  @override
  String get raceGnomeDesc => '聪明的工匠';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '$minutes 分钟内攻击 +$atk%、防御 +$def%';
  }

  @override
  String fedToast(int atk, int def) {
    return '饱食：攻击 +$atk%，防御 +$def%';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return '饱食 · 攻击 +$atk% · 防御 +$def% · 剩余 $time';
  }

  @override
  String huntReach(String regions) {
    return '可在此狩猎的区域：$regions';
  }

  @override
  String get huntNoMobsHere => '这里没有可狩猎的怪物。请先前往狩猎区域。';

  @override
  String get huntOtherRegions => '其他区域的怪物';

  @override
  String get huntTravelHint => '要狩猎某只怪物，请前往它所在的区域或有道路相连的相邻区域。';

  @override
  String get errMobRegion => '这只怪物在其他区域。请先前往那里。';

  @override
  String restoresHpPct(int pct) {
    return '恢复最大生命值的$pct%';
  }

  @override
  String get statEvasion => '% 闪避';

  @override
  String get statBlock => '% 格挡';

  @override
  String get statCritDmg => '% 暴击伤害';

  @override
  String get statLifesteal => '% 生命偷取';

  @override
  String get statDmgReduction => '% 伤害减免';

  @override
  String get classJewelry => '饰品';
}
