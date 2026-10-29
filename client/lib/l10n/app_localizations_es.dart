// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Spanish Castilian (`es`).
class AppLocalizationsEs extends AppLocalizations {
  AppLocalizationsEs([String locale = 'es']) : super(locale);

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
  String get classWarrior => 'Guerrero';

  @override
  String get classMage => 'Mago';

  @override
  String get classRogue => 'Pícaro';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Lvl $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Nv $level';
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
  String get close => 'Cerrar';

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
  String get errWrongRegion => 'Este monstruo vive en otra zona. Viaja allí primero.';

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
  String get levelMax => 'MÁX';

  @override
  String get promotionAvailable => 'ASCENSO DISPONIBLE';

  @override
  String get promotionChoose => 'Elegir';

  @override
  String promotionConfirm(String name) {
    return '¿Convertirte en $name?';
  }

  @override
  String get promotionBody => 'Esta elección define tu clase para siempre.';

  @override
  String nextMilestone(int level, String what) {
    return 'Siguiente: Nv $level · $what';
  }

  @override
  String get msPromotion => 'Ascenso de clase';

  @override
  String msAwakening(String name) {
    return 'Despertar: $name';
  }

  @override
  String msUltimate(String name) {
    return 'Definitiva: $name';
  }

  @override
  String get msCapstone => 'Talento cúspide';

  @override
  String get msMastery => 'Maestría de clase';

  @override
  String msTitle(String name) {
    return 'Título: $name';
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
    return '$xp PX esperan tu elección';
  }

  @override
  String attributePointsBtn(int n) {
    return 'Atributos · $n puntos';
  }

  @override
  String talentPointsBtn(int n) {
    return 'Talentos · $n puntos';
  }

  @override
  String get talentsTitle => 'Talentos';

  @override
  String capstoneLocked(int level, int points) {
    return 'Requiere nivel $level y $points puntos en este árbol';
  }

  @override
  String masteryUnlocked(String title) {
    return 'Maestría de clase: $title';
  }

  @override
  String get errPromotionUnavailable => 'Ese ascenso no está disponible.';

  @override
  String get errNotEnoughPoints => 'No hay puntos suficientes.';

  @override
  String get errTalentMaxed => 'Ese talento ya está al máximo.';

  @override
  String get errCapstoneLocked => 'La cúspide requiere su nivel y un árbol completo.';

  @override
  String get errSpecRequirement => 'Tu especialización no puede usar este objeto.';

  @override
  String get errAttributeRequirement => 'Tus atributos son demasiado bajos para este objeto.';

  @override
  String get errBusyCombat => 'No mientras estés en combate, AFK o en una mazmorra.';

  @override
  String skillGuard(int pct, int turns) {
    return '+$pct% de armadura durante $turns turnos';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '+$pct% de ataque durante $turns turnos';
  }

  @override
  String get classRanger => 'Explorador';

  @override
  String get classCleric => 'Clérigo';

  @override
  String get classMonk => 'Monje';

  @override
  String get classPaladin => 'Paladín';

  @override
  String get classDruid => 'Druida';

  @override
  String get classBard => 'Bardo';

  @override
  String get classShaman => 'Chamán';

  @override
  String get errEquipmentPermission => 'Tu clase no puede usar este tipo de equipo.';

  @override
  String get classResource => 'Recurso';

  @override
  String get classFocus => 'Atributo principal';

  @override
  String classSecondary(String names) {
    return 'Secundario: $names';
  }

  @override
  String get resStartsEmpty => 'Empieza vacío';

  @override
  String get resStartsFull => 'Empieza lleno';

  @override
  String resStartsAt(int pct) {
    return 'Empieza al $pct%';
  }

  @override
  String resPerHit(int n) {
    return '+$n por golpe';
  }

  @override
  String resPerTurn(int n) {
    return '+$n por turno';
  }

  @override
  String baseAttributesHint(int total) {
    return 'Atributos base · 1 punto cada 2 niveles, $total en total';
  }

  @override
  String get classEquipment => 'Equipo';

  @override
  String get firstSkills => 'Primeras habilidades';

  @override
  String get classPath => 'Ruta de clase';

  @override
  String get classPathHint => 'Toca una clase para ver sus habilidades.';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Nv $awakening Despertar · Nv $ultimate Definitiva · Nv $capstone talento cúspide · Nv $mastery maestría de clase';
  }

  @override
  String get skillTagActive => 'activa';

  @override
  String get skillTagPassive => 'pasiva';

  @override
  String get skillTagAwakening => 'despertar';

  @override
  String get skillTagUltimate => 'definitiva';

  @override
  String skillCostCd(int cost, int cd) {
    return 'Coste $cost · reutilización $cd turnos';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Maestría Nv $level: $title';
  }

  @override
  String get classInfoFailed => 'No se pudieron cargar los detalles de la clase.';

  @override
  String xpShare(int pct) {
    return 'PX $pct%';
  }

  @override
  String get dungeonsSolo => 'Mazmorras';

  @override
  String get raidsTab => 'Incursiones';

  @override
  String raidGroup(int min, int max) {
    return 'Grupo de $min a $max';
  }

  @override
  String raidBossGear(int count) {
    return 'Jefe: $count piezas de equipo garantizadas';
  }

  @override
  String get raidOpen => 'Abrir incursión';

  @override
  String get raidNeedParty =>
      'Las incursiones son para un grupo. Ponte en cola para uno al azar, únete a un anuncio o publica el tuyo.';

  @override
  String get raidLeaderOnly => 'El líder del grupo abre las incursiones.';

  @override
  String get raidLobby => 'Sala de la incursión';

  @override
  String get raidJoin => 'Unirse';

  @override
  String get raidLeave => 'Salir';

  @override
  String get raidCancel => 'Cancelar';

  @override
  String get raidCancelBody => 'La sala se cierra para todos.';

  @override
  String get raidStart => 'Iniciar incursión';

  @override
  String get raidStartBody => 'Una vez iniciada nadie puede unirse ni salir.';

  @override
  String raidMembers(int count, int max) {
    return 'Miembros ($count/$max)';
  }

  @override
  String raidNeedMore(int min) {
    return 'Se necesitan al menos $min miembros';
  }

  @override
  String get raidWaitHint => 'Las recompensas se abren al terminar; cada miembro recoge las suyas.';

  @override
  String get raidDone => 'Incursión terminada: recoge tus recompensas';

  @override
  String get raidLocked => 'Completada esta semana';

  @override
  String get raidCollected => 'recogido';

  @override
  String get errRaidNeedsParty => 'Únete primero a un grupo.';

  @override
  String get errRaidExists => 'Tu grupo ya tiene una incursión abierta.';

  @override
  String get errAlreadyInRaid => 'Ya estás en una incursión.';

  @override
  String get errRaidNotOpen => 'Esta incursión ya no está abierta.';

  @override
  String get errRaidNotEnough => 'No hay suficientes miembros para empezar.';

  @override
  String get errRaidWeeklyLocked => 'Ya completaste esta incursión esta semana.';

  @override
  String errCharacterBusy(String action) {
    return 'Termina primero lo que hace tu personaje: $action.';
  }

  @override
  String get busyGathering => 'recolección';

  @override
  String get busyHunt => 'caza';

  @override
  String get busyDungeon => 'mazmorra';

  @override
  String get busyRaid => 'incursión';

  @override
  String get busyTravel => 'viaje';

  @override
  String get busyCrafting => 'fabricación';

  @override
  String get errCraftActive => 'Ya hay una fabricación en curso.';

  @override
  String get errCraftTooLong => 'Ese trabajo duraría más de 6 horas.';

  @override
  String craftEach(String time) {
    return '$time cada uno';
  }

  @override
  String craftTotal(String time) {
    return 'Total $time';
  }

  @override
  String craftRunning(String item, int count) {
    return 'Fabricando $item ×$count';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count listos';
  }

  @override
  String get craftDone => 'Todo listo: recoge';

  @override
  String get craftCollectEarlyTitle => '¿Recoger ahora?';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return 'Se entregan los $done objetos terminados; los materiales y el oro de los otros $rest se devuelven.';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '$made fabricados, $refunded devueltos';
  }

  @override
  String get craftBusyHint => 'Fabricar lleva tiempo; mientras tanto tu personaje no hace nada más.';

  @override
  String get eventHowTitle => 'Cómo participar';

  @override
  String get eventRewardsTitle => 'Recompensas';

  @override
  String eventStepEnlist(int level) {
    return 'Inscríbete desde 24 horas antes de que empiece hasta que termine (nivel $level+).';
  }

  @override
  String eventStepRegion(String region) {
    return 'Debes estar en $region al inscribirte.';
  }

  @override
  String eventStepPower(int power) {
    return 'Tu poder se fija al inscribirte: ataque × 2 + defensa + PV ÷ 10. El tuyo ahora: $power.';
  }

  @override
  String eventStepScore(int tier) {
    return 'Puntos = poder × minutos de asistencia, desde el inicio (o desde que te inscribes, si es después). Se detiene en el nivel $tier.';
  }

  @override
  String get eventStepGoal =>
      'Si entre todos se alcanza el objetivo, las recompensas llegan al nivel máximo; si no, se detienen en el nivel 2.';

  @override
  String get eventStepPaid => 'Las recompensas se pagan solas al cerrar. Inscribirte no detiene tus otras acciones.';

  @override
  String eventTier(int tier, String score) {
    return 'Nivel $tier · $score+ puntos';
  }

  @override
  String eventTierTokens(int count) {
    return '+$count fichas de evento';
  }

  @override
  String eventEstimate(String score, int tier) {
    return 'Si te unes ahora: unos $score puntos → nivel $tier';
  }

  @override
  String eventEstimateNone(String score) {
    return 'Si te unes ahora: unos $score puntos, no alcanza un nivel';
  }

  @override
  String seasonalHow(int count) {
    return 'Mientras dure el evento, reclama $count fichas una vez al día (el día empieza a las 00:00 UTC).';
  }

  @override
  String seasonalBenefit(int count) {
    return 'Las fichas se guardan en tu cuenta (saldo: $count).';
  }

  @override
  String get eventDetails => 'Detalles';

  @override
  String get createCharacterTitle => 'Crear personaje';

  @override
  String get classTraits => 'Rasgos de la clase';

  @override
  String classResourceLine(String name) {
    return 'Recurso: $name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return 'Atributo principal: $primary · segundo: $secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return 'Atributo principal: $primary';
  }

  @override
  String get classDifficulty => 'Dificultad';

  @override
  String get difficultyEasy => 'Fácil — indulgente, segura en cacerías AFK';

  @override
  String get difficultyMedium => 'Media — equilibrada; importa el orden de habilidades';

  @override
  String get difficultyHard => 'Difícil — frágil; necesita el equipo adecuado';

  @override
  String get skillsTapHint => 'toca para ver detalles';

  @override
  String get chooseRace => 'Elige tu raza';

  @override
  String get raceCosmeticHint => 'Tu raza define tu aspecto; no cambia los atributos.';

  @override
  String get raceHumanDesc => 'Gente de las ciudades';

  @override
  String get raceElfDesc => 'Hijos de los bosques antiguos';

  @override
  String get raceDwarfDesc => 'Herreros de las montañas';

  @override
  String get raceOrcDesc => 'Guerreros de la estepa';

  @override
  String continueAs(String name) {
    return 'Continuar como $name';
  }

  @override
  String get nameYourCharacter => 'Ponle nombre a tu personaje';

  @override
  String get raceTroll => 'Trol';

  @override
  String get raceUndead => 'No muerto';

  @override
  String get raceDarkElf => 'Elfo oscuro';

  @override
  String get raceDrakari => 'Drakari';

  @override
  String get raceTauren => 'Tauren';

  @override
  String get raceGnome => 'Gnomo';

  @override
  String get raceTrollDesc => 'Cazadores de las islas selváticas';

  @override
  String get raceUndeadDesc => 'Alzados de la tumba';

  @override
  String get raceDarkElfDesc => 'Hijos de la oscuridad profunda';

  @override
  String get raceDrakariDesc => 'Linaje de sangre de dragón';

  @override
  String get raceTaurenDesc => 'Guardianes de las llanuras';

  @override
  String get raceGnomeDesc => 'Ingeniosos inventores';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '+$atk% ataque y +$def% defensa durante $minutes min';
  }

  @override
  String fedToast(int atk, int def) {
    return 'Saciado: +$atk% ataque, +$def% defensa';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return 'Saciado · +$atk% ataque · +$def% defensa · quedan $time';
  }

  @override
  String huntReach(String regions) {
    return 'Zonas de caza desde aquí: $regions';
  }

  @override
  String get huntNoMobsHere => 'Aquí no hay nada que cazar. Viaja primero a una zona de caza.';

  @override
  String get huntOtherRegions => 'Monstruos de otras zonas';

  @override
  String get huntTravelHint => 'Para cazar un monstruo, viaja a su zona o a una zona unida a ella por camino.';

  @override
  String get errMobRegion => 'Este monstruo vive en otra zona. Viaja allí primero.';

  @override
  String restoresHpPct(int pct) {
    return 'Restaura el $pct% de la vida máxima';
  }

  @override
  String get statEvasion => '% Evasión';

  @override
  String get statBlock => '% Bloqueo';

  @override
  String get statCritDmg => '% Daño crítico';

  @override
  String get statLifesteal => '% Robo de vida';

  @override
  String get statDmgReduction => '% Reducción de daño';

  @override
  String get classJewelry => 'Joyería';

  @override
  String get mapHint =>
      'Puedes viajar directamente a cualquier zona que tu nivel permita; el tiempo sigue el camino más corto. Arrastra y amplía el mapa; toca una zona.';

  @override
  String get mapNoRoad => 'Ningún camino directo lleva aquí. Viaja por las zonas que brillan a tu lado.';

  @override
  String get youAreHere => 'Estás aquí.';

  @override
  String travelTime(String time) {
    return 'Tiempo de viaje: $time';
  }

  @override
  String xpPerKill(int xp) {
    return '$xp XP';
  }

  @override
  String get deleteCharacter => 'Eliminar personaje';

  @override
  String deleteCharacterBody(String name) {
    return '$name se eliminará ahora mismo, con todos sus objetos y oro. No se puede deshacer. Escribe el nombre del personaje para confirmar.';
  }

  @override
  String get characterNameConfirm => 'Nombre del personaje';

  @override
  String characterDeletesAt(String when) {
    return 'Se elimina el $when';
  }

  @override
  String get cancelCharacterDelete => 'Cancelar eliminación';

  @override
  String characterDeleted(String name) {
    return '$name fue eliminado.';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name ya no se eliminará.';
  }

  @override
  String get errNameMismatch => 'El nombre no coincide.';

  @override
  String attrGainPerPoint(String gains) {
    return 'Por punto: $gains';
  }

  @override
  String attrGainPending(String gains) {
    return 'Se añade: $gains';
  }

  @override
  String get attrNoEffect => 'Sin efecto para tu clase';

  @override
  String get attrInputHint => 'Escribe cuántos puntos recibe cada atributo';

  @override
  String get gainHp => 'Vida';

  @override
  String get gainMp => 'Maná';

  @override
  String get gainHpRegen => 'Vida/s';

  @override
  String get huntSupplies => 'Suministros';

  @override
  String get huntSuppliesHint =>
      'Las pociones se beben entre combates con menos de la mitad de vida, la comida cuando acaba su bonificación y las de maná con poco maná. Lo que sobre vuelve a tu bolsa.';

  @override
  String get supplyHp => 'Pociones de vida';

  @override
  String get supplyMp => 'Pociones de maná';

  @override
  String get supplyMpHint => 'Para clases que usan maná, energía, concentración, energía natural o espíritu';

  @override
  String get supplyFood => 'Comida';

  @override
  String get huntNoSupplies => 'No tienes pociones ni comida: cómpralas en la tienda o fabrícalas.';

  @override
  String restoresMpPct(int pct) {
    return 'Rellena el $pct% de tu maná';
  }

  @override
  String get huntSupplyOnly => 'Solo se usa como suministro de caza';

  @override
  String get errBadSupply => 'Ese objeto no se puede llevar de caza.';

  @override
  String get errNotEnoughSupplies => 'No tienes suficientes de ese suministro.';

  @override
  String get errHuntOnly => 'Llévalo de caza como suministro.';

  @override
  String get huntFled => 'Tu vida bajó del 5 %: abandonaste el combate';

  @override
  String suppliesUsed(String list) {
    return 'Usado: $list';
  }

  @override
  String get compareTitle => 'Comparado con lo que llevas';

  @override
  String compareWorn(String name) {
    return 'Equipado: $name';
  }

  @override
  String get compareEmpty => 'Esa ranura está vacía: al equiparlo se suman todos.';

  @override
  String get compareColStat => 'Atributo';

  @override
  String get compareColWorn => 'Equipado';

  @override
  String get compareColThis => 'Este';

  @override
  String get compareColDiff => 'Cambio';

  @override
  String get compareSame => 'Mismos atributos que lo que llevas.';

  @override
  String get wearOk => 'Tu clase puede llevarlo';

  @override
  String get wearClass => 'Tu clase no puede llevarlo';

  @override
  String wearLevel(int level) {
    return 'Requiere nivel $level';
  }

  @override
  String get wearAttribute => 'Tus atributos no bastan';

  @override
  String get wearLegendClass => 'no es para tu clase';

  @override
  String get wearLegendLater => 'aún no (nivel o atributos)';

  @override
  String get enchantAction => 'Encantar';

  @override
  String get enchantPick => 'Elige el equipo a encantar';

  @override
  String get enchantNone => 'Ningún equipo encaja con este pergamino.';

  @override
  String enchantFits(String groups) {
    return 'Se aplica a: $groups';
  }

  @override
  String get enchantGroupWeapon => 'armas';

  @override
  String get enchantGroupArmor => 'armadura';

  @override
  String get enchantGroupJewelry => 'joyas';

  @override
  String get enchantGroupAny => 'cualquier equipo';

  @override
  String get enchantReplaceTitle => '¿Reemplazar el encantamiento?';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old se pierde y $next ocupa su lugar.';
  }

  @override
  String enchantDone(String item) {
    return '$item encantado';
  }

  @override
  String enchantOn(String name) {
    return 'Encantamiento: $name';
  }

  @override
  String get wornTag => 'equipado';

  @override
  String get errNotEnchant => 'No es un pergamino de encantamiento.';

  @override
  String get errEnchantSlot => 'El encantamiento no encaja en ese equipo.';

  @override
  String get statGatherMining => 'Rendimiento minero';

  @override
  String get statGatherHerbalism => 'Rendimiento de herboristería';

  @override
  String get statGatherWoodcutting => 'Rendimiento de tala';

  @override
  String get statGatherFishing => 'Rendimiento de pesca';

  @override
  String get statGatherHunting => 'Rendimiento de caza';

  @override
  String get statCraftSpeed => 'Velocidad de fabricación';

  @override
  String diffXp(int pct) {
    return '$pct % XP';
  }

  @override
  String get diffNoXp => 'sin XP';

  @override
  String get diffLegend => 'Color de receta: la XP que da a tu nivel';

  @override
  String get recipeGroupEnchantWeapon => 'Encantamientos de arma';

  @override
  String get recipeGroupEnchantArmor => 'Encantamientos de armadura';

  @override
  String get recipeGroupEnchantJewelry => 'Encantamientos de joyería';

  @override
  String get recipeGroupEnchantProfession => 'Encantamientos de oficio';

  @override
  String get partyHint =>
      'Los miembros que hacen el mismo trabajo a la vez (la misma recolección, profesión, caza o mazmorra) ganan cada uno +3 % PX por miembro: cinco juntos, +15 %. Nadie cede PX.';

  @override
  String get partyChat => 'Chat de grupo';

  @override
  String partyShared(String names, int kept, int earned) {
    return 'Bonificación de grupo con $names: $earned → $kept PX';
  }

  @override
  String evRaidMatched(String raid) {
    return 'Se encontró un grupo para $raid. Mira Incursiones.';
  }

  @override
  String get nkRaidMatched => 'Grupo de incursión';

  @override
  String get workHunt => 'Caza';

  @override
  String get workDungeon => 'Mazmorra';

  @override
  String get raidFind => 'Buscar grupo';

  @override
  String get raidQueueJoin => 'Hacer cola';

  @override
  String get raidQueueLeave => 'Salir de la cola';

  @override
  String raidQueued(String raid, int waiting) {
    return 'En cola para $raid · $waiting esperando';
  }

  @override
  String get raidQueueHint => 'Se forma un grupo al azar cuando esperan 5, o tras 5 minutos con al menos 2.';

  @override
  String get raidAdvertise => 'Anunciar';

  @override
  String get raidUnlist => 'Retirar anuncio';

  @override
  String raidListed(String raid) {
    return 'Tu grupo está en el tablón: $raid';
  }

  @override
  String get raidBoard => 'Tablón de anuncios';

  @override
  String get raidBoardEmpty => 'No hay anuncios abiertos.';

  @override
  String get raidJoinParty => 'Unirse';

  @override
  String get raidPickRaid => 'Incursión';

  @override
  String get errListingGone => 'Ese anuncio ya no existe o está lleno.';

  @override
  String get workGathering => 'Recolección';

  @override
  String get workAdventure => 'Aventura';

  @override
  String get workRaid => 'Incursión';

  @override
  String get workTravel => 'Viajando';

  @override
  String get workCrafting => 'Fabricación';

  @override
  String get workIdle => 'Inactivo';

  @override
  String get errDungeonRegion => 'Primero viaja a este lugar en el mapa.';

  @override
  String get dungeonsElsewhere =>
      'Más mazmorras e incursiones se abren al viajar: cada una se entra desde su región del mapa.';

  @override
  String get dungeonNoneHere => 'No hay mazmorras aquí. Viaja por el mapa a una región con mazmorra.';

  @override
  String dungeonXp(int xp) {
    return '≈ $xp PX para ti';
  }

  @override
  String appVersionLine(String version) {
    return 'Versión $version';
  }

  @override
  String get updateButton => 'Actualizar';

  @override
  String updateAvailable(String version) {
    return 'Nueva versión lista: $version';
  }

  @override
  String get updateLatest => 'Tienes la última versión';

  @override
  String updateConfirm(String version, int size) {
    return '¿Descargar e instalar la versión $version ($size MB)?';
  }

  @override
  String get updateDownloading => 'Descargando…';

  @override
  String get updatePermission => 'Permite instalaciones desde MeNa World y vuelve a tocar Actualizar.';

  @override
  String get updateFailed => 'No se pudo obtener la actualización. Revisa tu conexión.';

  @override
  String get fightReport => 'Informe de combate';

  @override
  String get endTimeUp => 'Se acabó el tiempo: duró hasta el final previsto.';

  @override
  String get endEarly => 'Recogido antes: se detuvo antes del final previsto.';

  @override
  String endFled(String mob) {
    return 'Tu vida bajó del 5 % contra $mob: te retiraste.';
  }

  @override
  String endDefeat(String mob) {
    return 'Caíste ante $mob.';
  }

  @override
  String endPartyDefeat(String mob) {
    return 'El grupo cayó ante $mob.';
  }

  @override
  String endTimeout(String mob) {
    return 'El combate con $mob superó el límite de turnos: se detuvo.';
  }

  @override
  String get endCleared => 'Superada: cayeron todas las oleadas y el jefe.';

  @override
  String get repCampHint =>
      'Por debajo del 50 % de vida el cazador acampa hasta llenarla antes del siguiente combate, así que rara vez baja al 5 %.';

  @override
  String get repYou => 'Tú';

  @override
  String get repLevel => 'Nivel';

  @override
  String get repHp => 'Vida';

  @override
  String get repAttack => 'Ataque';

  @override
  String get repDefense => 'Defensa';

  @override
  String get repEvasion => 'Evasión';

  @override
  String get repHitChance => 'Prob. de acierto';

  @override
  String get repCrit => 'Crítico';

  @override
  String repFights(int n) {
    return '$n combates';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '$hits/$swings aciertos ($pct %), media $avg, máx. $max';
  }

  @override
  String get repYourSwings => 'Tus golpes';

  @override
  String get repPartySwings => 'Golpes del grupo';

  @override
  String get repEnemySwings => 'Golpes del enemigo';

  @override
  String repLowest(int pct) {
    return 'Vida más baja en combate: $pct %';
  }

  @override
  String repRest(String time) {
    return 'Acampado: $time';
  }

  @override
  String get restart => 'Reiniciar';

  @override
  String get restartTitle => '¿Recoger y empezar de nuevo?';

  @override
  String get restartBody => 'Lo obtenido hasta ahora va a tu bolsa y lo mismo vuelve a empezar con la misma duración.';

  @override
  String get dailyTitle => 'Recompensa diaria';

  @override
  String get dailyBody =>
      'Reclama una vez al día con cualquier personaje. Siete días seguidos dan más; si faltas un día, el ciclo vuelve a empezar.';

  @override
  String get dailyClaim => 'Reclamar';

  @override
  String get dailyTomorrow => 'Vuelve mañana';

  @override
  String dailyDay(int n) {
    return 'Día $n';
  }

  @override
  String get dailyOpen => 'Recompensa diaria';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name fue silenciado $min minutos por insultar.';
  }

  @override
  String get storyTitle => 'Historia';

  @override
  String get storyOpen => 'Modo historia';

  @override
  String get storyComingSoon =>
      'El modo historia llegará pronto. Los capítulos aparecerán aquí a medida que se escriban.';

  @override
  String storyLocked(int level) {
    return 'Se abre en el nivel $level';
  }

  @override
  String get bankTitle => 'Banco';

  @override
  String get bankHint =>
      'Compartido por todos tus personajes. Toca un objeto de tu bolsa para depositarlo o uno del banco para sacarlo. Los objetos ligados se quedan con su personaje.';

  @override
  String get bankEmpty => 'El banco está vacío.';

  @override
  String get bankDeposit => 'Depositar';

  @override
  String get bankWithdraw => 'Sacar';

  @override
  String get errBankFull => 'El banco está lleno.';

  @override
  String get errItemBound => 'Los objetos ligados se quedan con su personaje.';

  @override
  String get errItemEquipped => 'Quítate el objeto primero.';

  @override
  String get tabTrade => 'Intercambio';

  @override
  String get tradeHint =>
      'Intercambia objetos con otro jugador. Ambos ponen objetos y confirman; cada lado paga un impuesto del 8 % del precio medio de mercado de lo que da, por unidad.';

  @override
  String get tradeWith => 'Nombre del jugador';

  @override
  String get tradeOpen => 'Iniciar intercambio';

  @override
  String tradeTitle(String name) {
    return 'Intercambio con $name';
  }

  @override
  String get tradeYours => 'Tú das';

  @override
  String tradeTheirs(String name) {
    return '$name da';
  }

  @override
  String get tradeNothing => 'Nada todavía.';

  @override
  String get tradeAdd => 'Añadir objeto';

  @override
  String get tradeRemove => 'Retirar';

  @override
  String get tradeTax => 'Tu impuesto';

  @override
  String get tradeReady => 'Listo';

  @override
  String get tradeWaiting => 'No listo';

  @override
  String get tradeConfirm => 'Confirmar intercambio';

  @override
  String get tradeWaitingThem => 'Esperando al otro jugador';

  @override
  String get tradeCancel => 'Cancelar intercambio';

  @override
  String get tradeCancelBody => 'No se intercambia nada.';

  @override
  String get tradeDone => 'Intercambio completado.';

  @override
  String get errTradeBusy => 'Uno de vosotros ya está intercambiando.';

  @override
  String get errTradeClosed => 'Ese intercambio terminó.';

  @override
  String get errTradeItem => 'Ese objeto no se puede intercambiar ahora.';

  @override
  String get errTradeEmpty => 'Primero pon algo.';

  @override
  String get errTradeTooMany => 'Demasiados objetos en el intercambio.';

  @override
  String get errCraftLevel => 'El nivel de tu personaje es inferior al del objeto.';
}
