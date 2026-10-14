import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_tr.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale) : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate = _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates = <LocalizationsDelegate<dynamic>>[
    delegate,
    GlobalMaterialLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
  ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[Locale('en'), Locale('tr')];

  /// No description provided for @appTitle.
  ///
  /// In en, this message translates to:
  /// **'New World MMO'**
  String get appTitle;

  /// No description provided for @tagline.
  ///
  /// In en, this message translates to:
  /// **'Adventures that keep going while you\'re away.'**
  String get tagline;

  /// No description provided for @emailOrUsername.
  ///
  /// In en, this message translates to:
  /// **'Email or username'**
  String get emailOrUsername;

  /// No description provided for @email.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get email;

  /// No description provided for @username.
  ///
  /// In en, this message translates to:
  /// **'Username'**
  String get username;

  /// No description provided for @password.
  ///
  /// In en, this message translates to:
  /// **'Password'**
  String get password;

  /// No description provided for @passwordConfirm.
  ///
  /// In en, this message translates to:
  /// **'Confirm password'**
  String get passwordConfirm;

  /// No description provided for @showPassword.
  ///
  /// In en, this message translates to:
  /// **'Show password'**
  String get showPassword;

  /// No description provided for @hidePassword.
  ///
  /// In en, this message translates to:
  /// **'Hide password'**
  String get hidePassword;

  /// No description provided for @login.
  ///
  /// In en, this message translates to:
  /// **'Log in'**
  String get login;

  /// No description provided for @register.
  ///
  /// In en, this message translates to:
  /// **'Create account'**
  String get register;

  /// No description provided for @noAccount.
  ///
  /// In en, this message translates to:
  /// **'New here? Create an account'**
  String get noAccount;

  /// No description provided for @haveAccount.
  ///
  /// In en, this message translates to:
  /// **'Already have an account? Log in'**
  String get haveAccount;

  /// No description provided for @acceptTerms.
  ///
  /// In en, this message translates to:
  /// **'I am 13 or older and accept the Terms of Service and Privacy Policy'**
  String get acceptTerms;

  /// No description provided for @fieldRequired.
  ///
  /// In en, this message translates to:
  /// **'Required'**
  String get fieldRequired;

  /// No description provided for @invalidEmail.
  ///
  /// In en, this message translates to:
  /// **'Enter a valid email'**
  String get invalidEmail;

  /// No description provided for @usernameRule.
  ///
  /// In en, this message translates to:
  /// **'3–32 letters, digits, _ or .'**
  String get usernameRule;

  /// No description provided for @passwordMin.
  ///
  /// In en, this message translates to:
  /// **'At least 8 characters'**
  String get passwordMin;

  /// No description provided for @passwordMismatch.
  ///
  /// In en, this message translates to:
  /// **'Passwords do not match'**
  String get passwordMismatch;

  /// No description provided for @termsRequired.
  ///
  /// In en, this message translates to:
  /// **'Please accept the terms to continue'**
  String get termsRequired;

  /// No description provided for @logout.
  ///
  /// In en, this message translates to:
  /// **'Log out'**
  String get logout;

  /// No description provided for @retry.
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get retry;

  /// No description provided for @offline.
  ///
  /// In en, this message translates to:
  /// **'Can\'t reach the server. Check your connection.'**
  String get offline;

  /// No description provided for @characters.
  ///
  /// In en, this message translates to:
  /// **'Characters'**
  String get characters;

  /// No description provided for @slotsUsed.
  ///
  /// In en, this message translates to:
  /// **'{used}/{cap} slots'**
  String slotsUsed(int used, int cap);

  /// No description provided for @newCharacter.
  ///
  /// In en, this message translates to:
  /// **'New character'**
  String get newCharacter;

  /// No description provided for @noCharacters.
  ///
  /// In en, this message translates to:
  /// **'No characters yet. Create your first hero.'**
  String get noCharacters;

  /// No description provided for @create.
  ///
  /// In en, this message translates to:
  /// **'Create'**
  String get create;

  /// No description provided for @characterName.
  ///
  /// In en, this message translates to:
  /// **'Name'**
  String get characterName;

  /// No description provided for @nameRule.
  ///
  /// In en, this message translates to:
  /// **'3–16 characters, starts with a letter (A–Z, 0–9, _)'**
  String get nameRule;

  /// No description provided for @race.
  ///
  /// In en, this message translates to:
  /// **'Race'**
  String get race;

  /// No description provided for @characterClass.
  ///
  /// In en, this message translates to:
  /// **'Class'**
  String get characterClass;

  /// No description provided for @raceHuman.
  ///
  /// In en, this message translates to:
  /// **'Human'**
  String get raceHuman;

  /// No description provided for @raceElf.
  ///
  /// In en, this message translates to:
  /// **'Elf'**
  String get raceElf;

  /// No description provided for @raceDwarf.
  ///
  /// In en, this message translates to:
  /// **'Dwarf'**
  String get raceDwarf;

  /// No description provided for @raceOrc.
  ///
  /// In en, this message translates to:
  /// **'Orc'**
  String get raceOrc;

  /// No description provided for @classWarrior.
  ///
  /// In en, this message translates to:
  /// **'Warrior'**
  String get classWarrior;

  /// No description provided for @classMage.
  ///
  /// In en, this message translates to:
  /// **'Mage'**
  String get classMage;

  /// No description provided for @classHunter.
  ///
  /// In en, this message translates to:
  /// **'Hunter'**
  String get classHunter;

  /// No description provided for @classPriest.
  ///
  /// In en, this message translates to:
  /// **'Priest'**
  String get classPriest;

  /// No description provided for @classRogue.
  ///
  /// In en, this message translates to:
  /// **'Rogue'**
  String get classRogue;

  /// No description provided for @levelLine.
  ///
  /// In en, this message translates to:
  /// **'Lvl {level} · {cls} · {race}'**
  String levelLine(int level, String cls, String race);

  /// No description provided for @levelShort.
  ///
  /// In en, this message translates to:
  /// **'Lvl {level}'**
  String levelShort(int level);

  /// No description provided for @hp.
  ///
  /// In en, this message translates to:
  /// **'HP'**
  String get hp;

  /// No description provided for @mp.
  ///
  /// In en, this message translates to:
  /// **'MP'**
  String get mp;

  /// No description provided for @stamina.
  ///
  /// In en, this message translates to:
  /// **'ST'**
  String get stamina;

  /// No description provided for @energy.
  ///
  /// In en, this message translates to:
  /// **'EN'**
  String get energy;

  /// No description provided for @xp.
  ///
  /// In en, this message translates to:
  /// **'XP'**
  String get xp;

  /// No description provided for @xpValue.
  ///
  /// In en, this message translates to:
  /// **'{xp} / {next}'**
  String xpValue(int xp, int next);

  /// No description provided for @attributes.
  ///
  /// In en, this message translates to:
  /// **'Attributes'**
  String get attributes;

  /// No description provided for @combat.
  ///
  /// In en, this message translates to:
  /// **'Combat'**
  String get combat;

  /// No description provided for @statStr.
  ///
  /// In en, this message translates to:
  /// **'STR'**
  String get statStr;

  /// No description provided for @statAgi.
  ///
  /// In en, this message translates to:
  /// **'AGI'**
  String get statAgi;

  /// No description provided for @statInt.
  ///
  /// In en, this message translates to:
  /// **'INT'**
  String get statInt;

  /// No description provided for @statVit.
  ///
  /// In en, this message translates to:
  /// **'VIT'**
  String get statVit;

  /// No description provided for @statCha.
  ///
  /// In en, this message translates to:
  /// **'CHA'**
  String get statCha;

  /// No description provided for @statLuck.
  ///
  /// In en, this message translates to:
  /// **'LUCK'**
  String get statLuck;

  /// No description provided for @attack.
  ///
  /// In en, this message translates to:
  /// **'Attack'**
  String get attack;

  /// No description provided for @defense.
  ///
  /// In en, this message translates to:
  /// **'Defense'**
  String get defense;

  /// No description provided for @crit.
  ///
  /// In en, this message translates to:
  /// **'Crit'**
  String get crit;

  /// No description provided for @hit.
  ///
  /// In en, this message translates to:
  /// **'Hit'**
  String get hit;

  /// No description provided for @percent.
  ///
  /// In en, this message translates to:
  /// **'{value}%'**
  String percent(int value);

  /// No description provided for @gold.
  ///
  /// In en, this message translates to:
  /// **'Gold'**
  String get gold;

  /// No description provided for @switchCharacter.
  ///
  /// In en, this message translates to:
  /// **'Switch character'**
  String get switchCharacter;

  /// No description provided for @statePendingDelete.
  ///
  /// In en, this message translates to:
  /// **'Deletion pending'**
  String get statePendingDelete;

  /// No description provided for @stateFrozen.
  ///
  /// In en, this message translates to:
  /// **'Frozen'**
  String get stateFrozen;

  /// No description provided for @errInvalidCredentials.
  ///
  /// In en, this message translates to:
  /// **'Wrong email/username or password.'**
  String get errInvalidCredentials;

  /// No description provided for @errAccountLocked.
  ///
  /// In en, this message translates to:
  /// **'Your account is temporarily locked.'**
  String get errAccountLocked;

  /// No description provided for @errAccountBanned.
  ///
  /// In en, this message translates to:
  /// **'This account has been banned.'**
  String get errAccountBanned;

  /// No description provided for @errAccountExists.
  ///
  /// In en, this message translates to:
  /// **'That email or username is already taken.'**
  String get errAccountExists;

  /// No description provided for @errRateLimited.
  ///
  /// In en, this message translates to:
  /// **'Too many attempts. Please wait a moment.'**
  String get errRateLimited;

  /// No description provided for @errNameTaken.
  ///
  /// In en, this message translates to:
  /// **'That name is already taken.'**
  String get errNameTaken;

  /// No description provided for @errSlotCap.
  ///
  /// In en, this message translates to:
  /// **'All character slots are in use.'**
  String get errSlotCap;

  /// No description provided for @errNotPlayable.
  ///
  /// In en, this message translates to:
  /// **'This character can\'t be played right now.'**
  String get errNotPlayable;

  /// No description provided for @errSessionExpired.
  ///
  /// In en, this message translates to:
  /// **'Your session expired. Please log in again.'**
  String get errSessionExpired;

  /// No description provided for @errGeneric.
  ///
  /// In en, this message translates to:
  /// **'Something went wrong ({code}).'**
  String errGeneric(String code);

  /// No description provided for @navCharacter.
  ///
  /// In en, this message translates to:
  /// **'Character'**
  String get navCharacter;

  /// No description provided for @navActivity.
  ///
  /// In en, this message translates to:
  /// **'Gather'**
  String get navActivity;

  /// No description provided for @actMining.
  ///
  /// In en, this message translates to:
  /// **'Mining'**
  String get actMining;

  /// No description provided for @actFishing.
  ///
  /// In en, this message translates to:
  /// **'Fishing'**
  String get actFishing;

  /// No description provided for @actWoodcutting.
  ///
  /// In en, this message translates to:
  /// **'Woodcutting'**
  String get actWoodcutting;

  /// No description provided for @actHerbalism.
  ///
  /// In en, this message translates to:
  /// **'Herbalism'**
  String get actHerbalism;

  /// No description provided for @actHunting.
  ///
  /// In en, this message translates to:
  /// **'Hunting'**
  String get actHunting;

  /// No description provided for @activityChoose.
  ///
  /// In en, this message translates to:
  /// **'Choose an activity'**
  String get activityChoose;

  /// No description provided for @durationLabel.
  ///
  /// In en, this message translates to:
  /// **'Duration'**
  String get durationLabel;

  /// No description provided for @gatherResource.
  ///
  /// In en, this message translates to:
  /// **'Resource'**
  String get gatherResource;

  /// No description provided for @gatherXp.
  ///
  /// In en, this message translates to:
  /// **'{xp} / {next} XP'**
  String gatherXp(int xp, int next);

  /// No description provided for @gatherMax.
  ///
  /// In en, this message translates to:
  /// **'Max level'**
  String get gatherMax;

  /// No description provided for @gatherLocked.
  ///
  /// In en, this message translates to:
  /// **'Unlocks at level {level}'**
  String gatherLocked(int level);

  /// No description provided for @gatherXpGained.
  ///
  /// In en, this message translates to:
  /// **'{name} +{xp} XP'**
  String gatherXpGained(String name, int xp);

  /// No description provided for @gatherLevelUp.
  ///
  /// In en, this message translates to:
  /// **'{name} reached level {level}!'**
  String gatherLevelUp(String name, int level);

  /// No description provided for @errGatherLevel.
  ///
  /// In en, this message translates to:
  /// **'Your gathering level is too low for this resource.'**
  String get errGatherLevel;

  /// No description provided for @startActivity.
  ///
  /// In en, this message translates to:
  /// **'Start {name}'**
  String startActivity(String name);

  /// No description provided for @collect.
  ///
  /// In en, this message translates to:
  /// **'Collect'**
  String get collect;

  /// No description provided for @abandon.
  ///
  /// In en, this message translates to:
  /// **'Abandon'**
  String get abandon;

  /// No description provided for @cancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// No description provided for @soFar.
  ///
  /// In en, this message translates to:
  /// **'Gathered so far'**
  String get soFar;

  /// No description provided for @activityDone.
  ///
  /// In en, this message translates to:
  /// **'Complete — ready to collect'**
  String get activityDone;

  /// No description provided for @remaining.
  ///
  /// In en, this message translates to:
  /// **'{time} left'**
  String remaining(String time);

  /// No description provided for @durHm.
  ///
  /// In en, this message translates to:
  /// **'{h}h {m}m'**
  String durHm(int h, String m);

  /// No description provided for @durM.
  ///
  /// In en, this message translates to:
  /// **'{m}m'**
  String durM(int m);

  /// No description provided for @durS.
  ///
  /// In en, this message translates to:
  /// **'{s}s'**
  String durS(int s);

  /// No description provided for @collectEarlyTitle.
  ///
  /// In en, this message translates to:
  /// **'Collect early?'**
  String get collectEarlyTitle;

  /// No description provided for @collectEarlyBody.
  ///
  /// In en, this message translates to:
  /// **'You get what has been gathered so far and the activity ends.'**
  String get collectEarlyBody;

  /// No description provided for @abandonTitle.
  ///
  /// In en, this message translates to:
  /// **'Abandon activity?'**
  String get abandonTitle;

  /// No description provided for @abandonBody.
  ///
  /// In en, this message translates to:
  /// **'Everything gathered so far will be lost.'**
  String get abandonBody;

  /// No description provided for @lootCollect.
  ///
  /// In en, this message translates to:
  /// **'Nice!'**
  String get lootCollect;

  /// No description provided for @lootNothing.
  ///
  /// In en, this message translates to:
  /// **'Nothing this time.'**
  String get lootNothing;

  /// No description provided for @xpGain.
  ///
  /// In en, this message translates to:
  /// **'+{xp} XP'**
  String xpGain(int xp);

  /// No description provided for @rarityCommon.
  ///
  /// In en, this message translates to:
  /// **'Common'**
  String get rarityCommon;

  /// No description provided for @rarityUncommon.
  ///
  /// In en, this message translates to:
  /// **'Uncommon'**
  String get rarityUncommon;

  /// No description provided for @rarityRare.
  ///
  /// In en, this message translates to:
  /// **'Rare'**
  String get rarityRare;

  /// No description provided for @rarityEpic.
  ///
  /// In en, this message translates to:
  /// **'Epic'**
  String get rarityEpic;

  /// No description provided for @rarityLegendary.
  ///
  /// In en, this message translates to:
  /// **'Legendary'**
  String get rarityLegendary;

  /// No description provided for @rarityMythic.
  ///
  /// In en, this message translates to:
  /// **'Mythic'**
  String get rarityMythic;

  /// No description provided for @errSlotBusy.
  ///
  /// In en, this message translates to:
  /// **'Another activity is already running.'**
  String get errSlotBusy;

  /// No description provided for @errNotActive.
  ///
  /// In en, this message translates to:
  /// **'This activity has already ended.'**
  String get errNotActive;

  /// No description provided for @evAchievement.
  ///
  /// In en, this message translates to:
  /// **'Achievement unlocked!'**
  String get evAchievement;

  /// No description provided for @evMasteryTier.
  ///
  /// In en, this message translates to:
  /// **'Mastery tier up!'**
  String get evMasteryTier;

  /// No description provided for @evMarketSold.
  ///
  /// In en, this message translates to:
  /// **'One of your listings sold.'**
  String get evMarketSold;

  /// No description provided for @evMailReceived.
  ///
  /// In en, this message translates to:
  /// **'You have new mail.'**
  String get evMailReceived;

  /// No description provided for @evPartyInvite.
  ///
  /// In en, this message translates to:
  /// **'You were invited to a party.'**
  String get evPartyInvite;

  /// No description provided for @evGuildInvite.
  ///
  /// In en, this message translates to:
  /// **'You were invited to a guild.'**
  String get evGuildInvite;

  /// No description provided for @navBag.
  ///
  /// In en, this message translates to:
  /// **'Bag'**
  String get navBag;

  /// No description provided for @equipment.
  ///
  /// In en, this message translates to:
  /// **'Equipment'**
  String get equipment;

  /// No description provided for @bag.
  ///
  /// In en, this message translates to:
  /// **'Bag'**
  String get bag;

  /// No description provided for @weightKg.
  ///
  /// In en, this message translates to:
  /// **'{cur}/{cap} kg'**
  String weightKg(String cur, String cap);

  /// No description provided for @slotHelm.
  ///
  /// In en, this message translates to:
  /// **'Helm'**
  String get slotHelm;

  /// No description provided for @slotNecklace.
  ///
  /// In en, this message translates to:
  /// **'Neck'**
  String get slotNecklace;

  /// No description provided for @slotChest.
  ///
  /// In en, this message translates to:
  /// **'Chest'**
  String get slotChest;

  /// No description provided for @slotGloves.
  ///
  /// In en, this message translates to:
  /// **'Gloves'**
  String get slotGloves;

  /// No description provided for @slotLegs.
  ///
  /// In en, this message translates to:
  /// **'Legs'**
  String get slotLegs;

  /// No description provided for @slotBoots.
  ///
  /// In en, this message translates to:
  /// **'Boots'**
  String get slotBoots;

  /// No description provided for @slotMainhand.
  ///
  /// In en, this message translates to:
  /// **'Main hand'**
  String get slotMainhand;

  /// No description provided for @slotOffhand.
  ///
  /// In en, this message translates to:
  /// **'Off hand'**
  String get slotOffhand;

  /// No description provided for @slotRing.
  ///
  /// In en, this message translates to:
  /// **'Ring'**
  String get slotRing;

  /// No description provided for @slotTrinket.
  ///
  /// In en, this message translates to:
  /// **'Trinket'**
  String get slotTrinket;

  /// No description provided for @equip.
  ///
  /// In en, this message translates to:
  /// **'Equip'**
  String get equip;

  /// No description provided for @unequip.
  ///
  /// In en, this message translates to:
  /// **'Unequip'**
  String get unequip;

  /// No description provided for @sell.
  ///
  /// In en, this message translates to:
  /// **'Sell'**
  String get sell;

  /// No description provided for @sortLabel.
  ///
  /// In en, this message translates to:
  /// **'Sort'**
  String get sortLabel;

  /// No description provided for @sortSlot.
  ///
  /// In en, this message translates to:
  /// **'Bag order'**
  String get sortSlot;

  /// No description provided for @sortRarity.
  ///
  /// In en, this message translates to:
  /// **'Rarity'**
  String get sortRarity;

  /// No description provided for @sortType.
  ///
  /// In en, this message translates to:
  /// **'Type'**
  String get sortType;

  /// No description provided for @sortName.
  ///
  /// In en, this message translates to:
  /// **'Name'**
  String get sortName;

  /// No description provided for @sortValue.
  ///
  /// In en, this message translates to:
  /// **'Value'**
  String get sortValue;

  /// No description provided for @selectedCount.
  ///
  /// In en, this message translates to:
  /// **'{count} selected'**
  String selectedCount(int count);

  /// No description provided for @selectAll.
  ///
  /// In en, this message translates to:
  /// **'Select all'**
  String get selectAll;

  /// No description provided for @clearSelection.
  ///
  /// In en, this message translates to:
  /// **'Clear selection'**
  String get clearSelection;

  /// No description provided for @sellSelectedTitle.
  ///
  /// In en, this message translates to:
  /// **'Sell {count} items?'**
  String sellSelectedTitle(int count);

  /// No description provided for @destroySelectedTitle.
  ///
  /// In en, this message translates to:
  /// **'Destroy {count} items?'**
  String destroySelectedTitle(int count);

  /// No description provided for @deleteEntry.
  ///
  /// In en, this message translates to:
  /// **'Delete'**
  String get deleteEntry;

  /// No description provided for @craftHowMany.
  ///
  /// In en, this message translates to:
  /// **'How many {name}?'**
  String craftHowMany(String name);

  /// No description provided for @craftCount.
  ///
  /// In en, this message translates to:
  /// **'Craft ×{count}'**
  String craftCount(int count);

  /// No description provided for @destroy.
  ///
  /// In en, this message translates to:
  /// **'Destroy'**
  String get destroy;

  /// No description provided for @sellTitle.
  ///
  /// In en, this message translates to:
  /// **'Sell {name}?'**
  String sellTitle(String name);

  /// No description provided for @sellBody.
  ///
  /// In en, this message translates to:
  /// **'You can buy it back from the vendor for 2 hours.'**
  String get sellBody;

  /// No description provided for @destroyTitle.
  ///
  /// In en, this message translates to:
  /// **'Destroy {name}?'**
  String destroyTitle(String name);

  /// No description provided for @destroyBody.
  ///
  /// In en, this message translates to:
  /// **'This cannot be undone.'**
  String get destroyBody;

  /// No description provided for @soldFor.
  ///
  /// In en, this message translates to:
  /// **'Sold for'**
  String get soldFor;

  /// No description provided for @quantity.
  ///
  /// In en, this message translates to:
  /// **'Quantity'**
  String get quantity;

  /// No description provided for @requiresLevel.
  ///
  /// In en, this message translates to:
  /// **'Requires level'**
  String get requiresLevel;

  /// No description provided for @weight.
  ///
  /// In en, this message translates to:
  /// **'Weight'**
  String get weight;

  /// No description provided for @weightGrams.
  ///
  /// In en, this message translates to:
  /// **'{g} g'**
  String weightGrams(int g);

  /// No description provided for @sellValue.
  ///
  /// In en, this message translates to:
  /// **'Sells for'**
  String get sellValue;

  /// No description provided for @bound.
  ///
  /// In en, this message translates to:
  /// **'Soulbound'**
  String get bound;

  /// No description provided for @itemLocked.
  ///
  /// In en, this message translates to:
  /// **'This item is held by a trade, mail or listing.'**
  String get itemLocked;

  /// No description provided for @errBagFull.
  ///
  /// In en, this message translates to:
  /// **'Your bag is full.'**
  String get errBagFull;

  /// No description provided for @errItemLocked.
  ///
  /// In en, this message translates to:
  /// **'This item is locked.'**
  String get errItemLocked;

  /// No description provided for @errClassRestriction.
  ///
  /// In en, this message translates to:
  /// **'Your class can\'t use this item.'**
  String get errClassRestriction;

  /// No description provided for @errLevelRequirement.
  ///
  /// In en, this message translates to:
  /// **'Your level is too low.'**
  String get errLevelRequirement;

  /// No description provided for @errNotEquippable.
  ///
  /// In en, this message translates to:
  /// **'This item can\'t be equipped.'**
  String get errNotEquippable;

  /// No description provided for @navAdventure.
  ///
  /// In en, this message translates to:
  /// **'Adventure'**
  String get navAdventure;

  /// No description provided for @segHunt.
  ///
  /// In en, this message translates to:
  /// **'Hunt'**
  String get segHunt;

  /// No description provided for @segDungeon.
  ///
  /// In en, this message translates to:
  /// **'Dungeon'**
  String get segDungeon;

  /// No description provided for @segTravel.
  ///
  /// In en, this message translates to:
  /// **'Travel'**
  String get segTravel;

  /// No description provided for @huntChoose.
  ///
  /// In en, this message translates to:
  /// **'Choose your prey'**
  String get huntChoose;

  /// No description provided for @startHunt.
  ///
  /// In en, this message translates to:
  /// **'Start hunt'**
  String get startHunt;

  /// No description provided for @huntDone.
  ///
  /// In en, this message translates to:
  /// **'Hunt over — collect your spoils'**
  String get huntDone;

  /// No description provided for @huntEarlyBody.
  ///
  /// In en, this message translates to:
  /// **'You keep the fights won so far and the hunt ends.'**
  String get huntEarlyBody;

  /// No description provided for @mobStats.
  ///
  /// In en, this message translates to:
  /// **'HP {hp} · ATK {atk} · DEF {def}'**
  String mobStats(int hp, int atk, int def);

  /// No description provided for @huntSummary.
  ///
  /// In en, this message translates to:
  /// **'{won} won · {lost} lost'**
  String huntSummary(int won, int lost);

  /// No description provided for @levelUp.
  ///
  /// In en, this message translates to:
  /// **'Level up! Now level {level}'**
  String levelUp(int level);

  /// No description provided for @levelReq.
  ///
  /// In en, this message translates to:
  /// **'Lvl {level}+'**
  String levelReq(int level);

  /// No description provided for @dungeonEncounters.
  ///
  /// In en, this message translates to:
  /// **'{count} fights · boss {boss}'**
  String dungeonEncounters(int count, String boss);

  /// No description provided for @dungeonEnter.
  ///
  /// In en, this message translates to:
  /// **'Enter'**
  String get dungeonEnter;

  /// No description provided for @dungeonDone.
  ///
  /// In en, this message translates to:
  /// **'Run complete — claim rewards'**
  String get dungeonDone;

  /// No description provided for @dungeonWaitHint.
  ///
  /// In en, this message translates to:
  /// **'Rewards unlock when the run ends.'**
  String get dungeonWaitHint;

  /// No description provided for @dungeonCleared.
  ///
  /// In en, this message translates to:
  /// **'Cleared!'**
  String get dungeonCleared;

  /// No description provided for @dungeonFailed.
  ///
  /// In en, this message translates to:
  /// **'Party wiped'**
  String get dungeonFailed;

  /// No description provided for @bossKilled.
  ///
  /// In en, this message translates to:
  /// **'Boss defeated'**
  String get bossKilled;

  /// No description provided for @bonus.
  ///
  /// In en, this message translates to:
  /// **'Clear bonus'**
  String get bonus;

  /// No description provided for @youAreIn.
  ///
  /// In en, this message translates to:
  /// **'You are in'**
  String get youAreIn;

  /// No description provided for @routes.
  ///
  /// In en, this message translates to:
  /// **'Routes'**
  String get routes;

  /// No description provided for @noRoutes.
  ///
  /// In en, this message translates to:
  /// **'No routes from here.'**
  String get noRoutes;

  /// No description provided for @travelTo.
  ///
  /// In en, this message translates to:
  /// **'Go'**
  String get travelTo;

  /// No description provided for @travelling.
  ///
  /// In en, this message translates to:
  /// **'Travelling to {place}'**
  String travelling(String place);

  /// No description provided for @arrive.
  ///
  /// In en, this message translates to:
  /// **'Arrive'**
  String get arrive;

  /// No description provided for @arrived.
  ///
  /// In en, this message translates to:
  /// **'You have arrived'**
  String get arrived;

  /// No description provided for @cancelTravel.
  ///
  /// In en, this message translates to:
  /// **'Turn back'**
  String get cancelTravel;

  /// No description provided for @kindCity.
  ///
  /// In en, this message translates to:
  /// **'City'**
  String get kindCity;

  /// No description provided for @kindZone.
  ///
  /// In en, this message translates to:
  /// **'Zone'**
  String get kindZone;

  /// No description provided for @kindWilderness.
  ///
  /// In en, this message translates to:
  /// **'Wilderness'**
  String get kindWilderness;

  /// No description provided for @kindDungeon.
  ///
  /// In en, this message translates to:
  /// **'Dungeon'**
  String get kindDungeon;

  /// No description provided for @kindContinent.
  ///
  /// In en, this message translates to:
  /// **'Continent'**
  String get kindContinent;

  /// No description provided for @safeZone.
  ///
  /// In en, this message translates to:
  /// **'safe'**
  String get safeZone;

  /// No description provided for @errBusy.
  ///
  /// In en, this message translates to:
  /// **'Something is already in progress.'**
  String get errBusy;

  /// No description provided for @errNotReady.
  ///
  /// In en, this message translates to:
  /// **'Not finished yet.'**
  String get errNotReady;

  /// No description provided for @errWeeklyLocked.
  ///
  /// In en, this message translates to:
  /// **'You already cleared this dungeon this week.'**
  String get errWeeklyLocked;

  /// No description provided for @errNoRoute.
  ///
  /// In en, this message translates to:
  /// **'There is no direct route there.'**
  String get errNoRoute;

  /// No description provided for @navWorkshop.
  ///
  /// In en, this message translates to:
  /// **'Workshop'**
  String get navWorkshop;

  /// No description provided for @segCraft.
  ///
  /// In en, this message translates to:
  /// **'Craft'**
  String get segCraft;

  /// No description provided for @segShop.
  ///
  /// In en, this message translates to:
  /// **'Shop'**
  String get segShop;

  /// No description provided for @profBlacksmithing.
  ///
  /// In en, this message translates to:
  /// **'Blacksmithing'**
  String get profBlacksmithing;

  /// No description provided for @profAlchemy.
  ///
  /// In en, this message translates to:
  /// **'Alchemy'**
  String get profAlchemy;

  /// No description provided for @profCooking.
  ///
  /// In en, this message translates to:
  /// **'Cooking'**
  String get profCooking;

  /// No description provided for @learn.
  ///
  /// In en, this message translates to:
  /// **'Learn'**
  String get learn;

  /// No description provided for @learnHint.
  ///
  /// In en, this message translates to:
  /// **'Learn this profession to unlock its recipes.'**
  String get learnHint;

  /// No description provided for @craft.
  ///
  /// In en, this message translates to:
  /// **'Craft'**
  String get craft;

  /// No description provided for @crafted.
  ///
  /// In en, this message translates to:
  /// **'Crafted {name} ×{count}'**
  String crafted(String name, int count);

  /// No description provided for @skillUp.
  ///
  /// In en, this message translates to:
  /// **'{profession} is now level {level}'**
  String skillUp(String profession, int level);

  /// No description provided for @requiresSkill.
  ///
  /// In en, this message translates to:
  /// **'Needs skill {level}'**
  String requiresSkill(int level);

  /// No description provided for @buy.
  ///
  /// In en, this message translates to:
  /// **'Buy'**
  String get buy;

  /// No description provided for @bought.
  ///
  /// In en, this message translates to:
  /// **'Bought {name}'**
  String bought(String name);

  /// No description provided for @buybacks.
  ///
  /// In en, this message translates to:
  /// **'Buy back'**
  String get buybacks;

  /// No description provided for @buybackHint.
  ///
  /// In en, this message translates to:
  /// **'Items you sold can be bought back for 2 hours.'**
  String get buybackHint;

  /// No description provided for @errInsufficientGold.
  ///
  /// In en, this message translates to:
  /// **'Not enough gold.'**
  String get errInsufficientGold;

  /// No description provided for @errInsufficientMaterials.
  ///
  /// In en, this message translates to:
  /// **'You\'re missing materials.'**
  String get errInsufficientMaterials;

  /// No description provided for @errSkillTooLow.
  ///
  /// In en, this message translates to:
  /// **'Your skill is too low for this recipe.'**
  String get errSkillTooLow;

  /// No description provided for @errProfessionNotLearned.
  ///
  /// In en, this message translates to:
  /// **'Learn the profession first.'**
  String get errProfessionNotLearned;

  /// No description provided for @errPriceChanged.
  ///
  /// In en, this message translates to:
  /// **'The price changed — check it and try again.'**
  String get errPriceChanged;

  /// No description provided for @errCannotSell.
  ///
  /// In en, this message translates to:
  /// **'This item can\'t be sold.'**
  String get errCannotSell;

  /// No description provided for @errBuybackExpired.
  ///
  /// In en, this message translates to:
  /// **'The buy-back window has closed.'**
  String get errBuybackExpired;

  /// No description provided for @segMarket.
  ///
  /// In en, this message translates to:
  /// **'Market'**
  String get segMarket;

  /// No description provided for @classWeapon.
  ///
  /// In en, this message translates to:
  /// **'Weapons'**
  String get classWeapon;

  /// No description provided for @classArmor.
  ///
  /// In en, this message translates to:
  /// **'Armor'**
  String get classArmor;

  /// No description provided for @classConsumable.
  ///
  /// In en, this message translates to:
  /// **'Consumables'**
  String get classConsumable;

  /// No description provided for @classMaterial.
  ///
  /// In en, this message translates to:
  /// **'Materials'**
  String get classMaterial;

  /// No description provided for @filterAll.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get filterAll;

  /// No description provided for @sortCheapest.
  ///
  /// In en, this message translates to:
  /// **'Cheapest'**
  String get sortCheapest;

  /// No description provided for @sortNewest.
  ///
  /// In en, this message translates to:
  /// **'Newest'**
  String get sortNewest;

  /// No description provided for @marketEmpty.
  ///
  /// In en, this message translates to:
  /// **'No listings match.'**
  String get marketEmpty;

  /// No description provided for @listingLine.
  ///
  /// In en, this message translates to:
  /// **'{seller} · {price} each'**
  String listingLine(String seller, String price);

  /// No description provided for @loadMore.
  ///
  /// In en, this message translates to:
  /// **'Load more'**
  String get loadMore;

  /// No description provided for @myListings.
  ///
  /// In en, this message translates to:
  /// **'My listings ({count})'**
  String myListings(int count);

  /// No description provided for @cancelListing.
  ///
  /// In en, this message translates to:
  /// **'Cancel listing'**
  String get cancelListing;

  /// No description provided for @cancelListingTitle.
  ///
  /// In en, this message translates to:
  /// **'Cancel {name}?'**
  String cancelListingTitle(String name);

  /// No description provided for @cancelListingBody.
  ///
  /// In en, this message translates to:
  /// **'The item returns to your bag. Part of the listing fee is refunded only within the first 6 hours.'**
  String get cancelListingBody;

  /// No description provided for @buyListingTitle.
  ///
  /// In en, this message translates to:
  /// **'Buy {name}?'**
  String buyListingTitle(String name);

  /// No description provided for @buyListingBody.
  ///
  /// In en, this message translates to:
  /// **'Pay {price} to {seller}.'**
  String buyListingBody(String price, String seller);

  /// No description provided for @listOnMarket.
  ///
  /// In en, this message translates to:
  /// **'List on market'**
  String get listOnMarket;

  /// No description provided for @listed.
  ///
  /// In en, this message translates to:
  /// **'{name} listed on the market'**
  String listed(String name);

  /// No description provided for @quantityOf.
  ///
  /// In en, this message translates to:
  /// **'Quantity: {count} of {max}'**
  String quantityOf(int count, int max);

  /// No description provided for @unitPrice.
  ///
  /// In en, this message translates to:
  /// **'Price per item'**
  String get unitPrice;

  /// No description provided for @suggestedPrice.
  ///
  /// In en, this message translates to:
  /// **'Market price: {price}'**
  String suggestedPrice(String price);

  /// No description provided for @listingFee.
  ///
  /// In en, this message translates to:
  /// **'Listing fee (non-refundable)'**
  String get listingFee;

  /// No description provided for @youReceive.
  ///
  /// In en, this message translates to:
  /// **'You receive when sold'**
  String get youReceive;

  /// No description provided for @errPriceOutOfRange.
  ///
  /// In en, this message translates to:
  /// **'That price is outside the allowed range.'**
  String get errPriceOutOfRange;

  /// No description provided for @errListingCap.
  ///
  /// In en, this message translates to:
  /// **'You have too many active listings.'**
  String get errListingCap;

  /// No description provided for @errSelfBuy.
  ///
  /// In en, this message translates to:
  /// **'You can\'t buy your own listing.'**
  String get errSelfBuy;

  /// No description provided for @errListingUnavailable.
  ///
  /// In en, this message translates to:
  /// **'This listing is no longer available.'**
  String get errListingUnavailable;

  /// No description provided for @errNotTradable.
  ///
  /// In en, this message translates to:
  /// **'This item can\'t be traded.'**
  String get errNotTradable;

  /// No description provided for @social.
  ///
  /// In en, this message translates to:
  /// **'Social'**
  String get social;

  /// No description provided for @tabChat.
  ///
  /// In en, this message translates to:
  /// **'Chat'**
  String get tabChat;

  /// No description provided for @tabFriends.
  ///
  /// In en, this message translates to:
  /// **'Friends'**
  String get tabFriends;

  /// No description provided for @tabMail.
  ///
  /// In en, this message translates to:
  /// **'Mail'**
  String get tabMail;

  /// No description provided for @chanGlobal.
  ///
  /// In en, this message translates to:
  /// **'World'**
  String get chanGlobal;

  /// No description provided for @chanGuild.
  ///
  /// In en, this message translates to:
  /// **'Guild'**
  String get chanGuild;

  /// No description provided for @chanParty.
  ///
  /// In en, this message translates to:
  /// **'Party'**
  String get chanParty;

  /// No description provided for @chatEmpty.
  ///
  /// In en, this message translates to:
  /// **'No messages yet. Say hello!'**
  String get chatEmpty;

  /// No description provided for @chatHint.
  ///
  /// In en, this message translates to:
  /// **'Message {channel}'**
  String chatHint(String channel);

  /// No description provided for @send.
  ///
  /// In en, this message translates to:
  /// **'Send'**
  String get send;

  /// No description provided for @addFriend.
  ///
  /// In en, this message translates to:
  /// **'Add'**
  String get addFriend;

  /// No description provided for @friendAdded.
  ///
  /// In en, this message translates to:
  /// **'You and {name} are now friends'**
  String friendAdded(String name);

  /// No description provided for @friendRequested.
  ///
  /// In en, this message translates to:
  /// **'Friend request sent to {name}'**
  String friendRequested(String name);

  /// No description provided for @friendRequests.
  ///
  /// In en, this message translates to:
  /// **'Friend requests'**
  String get friendRequests;

  /// No description provided for @friendsCount.
  ///
  /// In en, this message translates to:
  /// **'Friends ({count})'**
  String friendsCount(int count);

  /// No description provided for @noFriends.
  ///
  /// In en, this message translates to:
  /// **'No friends yet — add someone by name.'**
  String get noFriends;

  /// No description provided for @more.
  ///
  /// In en, this message translates to:
  /// **'More'**
  String get more;

  /// No description provided for @whisper.
  ///
  /// In en, this message translates to:
  /// **'Whisper'**
  String get whisper;

  /// No description provided for @remove.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get remove;

  /// No description provided for @block.
  ///
  /// In en, this message translates to:
  /// **'Block'**
  String get block;

  /// No description provided for @unblock.
  ///
  /// In en, this message translates to:
  /// **'Unblock'**
  String get unblock;

  /// No description provided for @blockedList.
  ///
  /// In en, this message translates to:
  /// **'Blocked'**
  String get blockedList;

  /// No description provided for @removeFriendTitle.
  ///
  /// In en, this message translates to:
  /// **'Remove {name}?'**
  String removeFriendTitle(String name);

  /// No description provided for @removeFriendBody.
  ///
  /// In en, this message translates to:
  /// **'You can send a new friend request later.'**
  String get removeFriendBody;

  /// No description provided for @blockTitle.
  ///
  /// In en, this message translates to:
  /// **'Block {name}?'**
  String blockTitle(String name);

  /// No description provided for @blockBody.
  ///
  /// In en, this message translates to:
  /// **'They can\'t whisper you, their chat is hidden and their mail is dropped. This also ends your friendship.'**
  String get blockBody;

  /// No description provided for @accept.
  ///
  /// In en, this message translates to:
  /// **'Accept'**
  String get accept;

  /// No description provided for @decline.
  ///
  /// In en, this message translates to:
  /// **'Decline'**
  String get decline;

  /// No description provided for @newMail.
  ///
  /// In en, this message translates to:
  /// **'Write mail'**
  String get newMail;

  /// No description provided for @inboxEmpty.
  ///
  /// In en, this message translates to:
  /// **'Your mailbox is empty.'**
  String get inboxEmpty;

  /// No description provided for @noSubject.
  ///
  /// In en, this message translates to:
  /// **'(no subject)'**
  String get noSubject;

  /// No description provided for @mailFrom.
  ///
  /// In en, this message translates to:
  /// **'From {name}'**
  String mailFrom(String name);

  /// No description provided for @attachments.
  ///
  /// In en, this message translates to:
  /// **'Attachments'**
  String get attachments;

  /// No description provided for @claim.
  ///
  /// In en, this message translates to:
  /// **'Take attachments'**
  String get claim;

  /// No description provided for @mailClaimed.
  ///
  /// In en, this message translates to:
  /// **'Attachments moved to your bag and wallet'**
  String get mailClaimed;

  /// No description provided for @deleteMailTitle.
  ///
  /// In en, this message translates to:
  /// **'Delete this mail?'**
  String get deleteMailTitle;

  /// No description provided for @deleteMailBody.
  ///
  /// In en, this message translates to:
  /// **'Deleted mail can\'t be recovered.'**
  String get deleteMailBody;

  /// No description provided for @mailSent.
  ///
  /// In en, this message translates to:
  /// **'Mail sent (postage {postage})'**
  String mailSent(String postage);

  /// No description provided for @mailTo.
  ///
  /// In en, this message translates to:
  /// **'To (character name)'**
  String get mailTo;

  /// No description provided for @mailSubject.
  ///
  /// In en, this message translates to:
  /// **'Subject'**
  String get mailSubject;

  /// No description provided for @mailBody.
  ///
  /// In en, this message translates to:
  /// **'Message'**
  String get mailBody;

  /// No description provided for @attachGold.
  ///
  /// In en, this message translates to:
  /// **'Attach gold'**
  String get attachGold;

  /// No description provided for @errAlreadyFriends.
  ///
  /// In en, this message translates to:
  /// **'You\'re already friends.'**
  String get errAlreadyFriends;

  /// No description provided for @errAlreadyRequested.
  ///
  /// In en, this message translates to:
  /// **'A friend request is already pending.'**
  String get errAlreadyRequested;

  /// No description provided for @errBlocked.
  ///
  /// In en, this message translates to:
  /// **'You can\'t interact with this player.'**
  String get errBlocked;

  /// No description provided for @errFriendCap.
  ///
  /// In en, this message translates to:
  /// **'Your friend list is full.'**
  String get errFriendCap;

  /// No description provided for @errTargetNotFound.
  ///
  /// In en, this message translates to:
  /// **'No character with that name.'**
  String get errTargetNotFound;

  /// No description provided for @errSelfTarget.
  ///
  /// In en, this message translates to:
  /// **'You can\'t target yourself.'**
  String get errSelfTarget;

  /// No description provided for @errInboxFull.
  ///
  /// In en, this message translates to:
  /// **'Their mailbox is full.'**
  String get errInboxFull;

  /// No description provided for @errTooLong.
  ///
  /// In en, this message translates to:
  /// **'That\'s too long.'**
  String get errTooLong;

  /// No description provided for @errEmpty.
  ///
  /// In en, this message translates to:
  /// **'Write something first.'**
  String get errEmpty;

  /// No description provided for @errNotInChannel.
  ///
  /// In en, this message translates to:
  /// **'You\'re not in that group.'**
  String get errNotInChannel;

  /// No description provided for @errMuted.
  ///
  /// In en, this message translates to:
  /// **'You are muted and can\'t chat right now.'**
  String get errMuted;

  /// No description provided for @errHasAttachments.
  ///
  /// In en, this message translates to:
  /// **'Take the attachments before deleting.'**
  String get errHasAttachments;

  /// No description provided for @useItem.
  ///
  /// In en, this message translates to:
  /// **'Use'**
  String get useItem;

  /// No description provided for @restoresHp.
  ///
  /// In en, this message translates to:
  /// **'Restores {amount} HP'**
  String restoresHp(int amount);

  /// No description provided for @usedItem.
  ///
  /// In en, this message translates to:
  /// **'Used {name} · HP {hp}/{max}'**
  String usedItem(String name, int hp, int max);

  /// No description provided for @errHpFull.
  ///
  /// In en, this message translates to:
  /// **'Your health is already full.'**
  String get errHpFull;

  /// No description provided for @errInCombat.
  ///
  /// In en, this message translates to:
  /// **'Finish your hunt or dungeon run first.'**
  String get errInCombat;

  /// No description provided for @errNotUsable.
  ///
  /// In en, this message translates to:
  /// **'This item can\'t be used.'**
  String get errNotUsable;

  /// No description provided for @errTooWounded.
  ///
  /// In en, this message translates to:
  /// **'You\'re too wounded to fight. Rest or drink a potion.'**
  String get errTooWounded;

  /// No description provided for @forgotPassword.
  ///
  /// In en, this message translates to:
  /// **'Forgot password?'**
  String get forgotPassword;

  /// No description provided for @resetHint.
  ///
  /// In en, this message translates to:
  /// **'Enter your account e-mail and we\'ll send you a link to set a new password.'**
  String get resetHint;

  /// No description provided for @resetSent.
  ///
  /// In en, this message translates to:
  /// **'If an account exists for that e-mail, a reset link is on its way. It expires in 1 hour.'**
  String get resetSent;

  /// No description provided for @sendLink.
  ///
  /// In en, this message translates to:
  /// **'Send link'**
  String get sendLink;

  /// No description provided for @close.
  ///
  /// In en, this message translates to:
  /// **'Close'**
  String get close;

  /// No description provided for @changePassword.
  ///
  /// In en, this message translates to:
  /// **'Change password'**
  String get changePassword;

  /// No description provided for @currentPassword.
  ///
  /// In en, this message translates to:
  /// **'Current password'**
  String get currentPassword;

  /// No description provided for @newPassword.
  ///
  /// In en, this message translates to:
  /// **'New password'**
  String get newPassword;

  /// No description provided for @save.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get save;

  /// No description provided for @passwordChanged.
  ///
  /// In en, this message translates to:
  /// **'Password changed. Your other devices were signed out.'**
  String get passwordChanged;

  /// No description provided for @verifyEmailBanner.
  ///
  /// In en, this message translates to:
  /// **'Please confirm your e-mail address.'**
  String get verifyEmailBanner;

  /// No description provided for @resend.
  ///
  /// In en, this message translates to:
  /// **'Resend'**
  String get resend;

  /// No description provided for @verificationSent.
  ///
  /// In en, this message translates to:
  /// **'Verification e-mail sent.'**
  String get verificationSent;

  /// No description provided for @errEmailNotVerified.
  ///
  /// In en, this message translates to:
  /// **'Confirm your e-mail before creating a character.'**
  String get errEmailNotVerified;

  /// No description provided for @errAlreadyVerified.
  ///
  /// In en, this message translates to:
  /// **'Your e-mail is already confirmed.'**
  String get errAlreadyVerified;

  /// No description provided for @tabParty.
  ///
  /// In en, this message translates to:
  /// **'Party'**
  String get tabParty;

  /// No description provided for @tabGuild.
  ///
  /// In en, this message translates to:
  /// **'Guild'**
  String get tabGuild;

  /// No description provided for @invite.
  ///
  /// In en, this message translates to:
  /// **'Invite'**
  String get invite;

  /// No description provided for @invitedBy.
  ///
  /// In en, this message translates to:
  /// **'Invited by {name}'**
  String invitedBy(String name);

  /// No description provided for @noParty.
  ///
  /// In en, this message translates to:
  /// **'You\'re not in a party. Create one and invite friends — fights go better together.'**
  String get noParty;

  /// No description provided for @createParty.
  ///
  /// In en, this message translates to:
  /// **'Create party'**
  String get createParty;

  /// No description provided for @pendingInvites.
  ///
  /// In en, this message translates to:
  /// **'Invitations'**
  String get pendingInvites;

  /// No description provided for @partyMembers.
  ///
  /// In en, this message translates to:
  /// **'Party ({count}/{cap})'**
  String partyMembers(int count, int cap);

  /// No description provided for @makeLeader.
  ///
  /// In en, this message translates to:
  /// **'Make leader'**
  String get makeLeader;

  /// No description provided for @kick.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get kick;

  /// No description provided for @leave.
  ///
  /// In en, this message translates to:
  /// **'Leave'**
  String get leave;

  /// No description provided for @disband.
  ///
  /// In en, this message translates to:
  /// **'Disband'**
  String get disband;

  /// No description provided for @leavePartyTitle.
  ///
  /// In en, this message translates to:
  /// **'Leave the party?'**
  String get leavePartyTitle;

  /// No description provided for @leavePartyBody.
  ///
  /// In en, this message translates to:
  /// **'If you are the leader, the longest-standing member takes over.'**
  String get leavePartyBody;

  /// No description provided for @disbandPartyTitle.
  ///
  /// In en, this message translates to:
  /// **'Disband the party?'**
  String get disbandPartyTitle;

  /// No description provided for @disbandPartyBody.
  ///
  /// In en, this message translates to:
  /// **'Everyone is removed from the party.'**
  String get disbandPartyBody;

  /// No description provided for @rankLeader.
  ///
  /// In en, this message translates to:
  /// **'Leader'**
  String get rankLeader;

  /// No description provided for @rankOfficer.
  ///
  /// In en, this message translates to:
  /// **'Officer'**
  String get rankOfficer;

  /// No description provided for @rankVeteran.
  ///
  /// In en, this message translates to:
  /// **'Veteran'**
  String get rankVeteran;

  /// No description provided for @rankMember.
  ///
  /// In en, this message translates to:
  /// **'Member'**
  String get rankMember;

  /// No description provided for @rankInitiate.
  ///
  /// In en, this message translates to:
  /// **'Initiate'**
  String get rankInitiate;

  /// No description provided for @createGuild.
  ///
  /// In en, this message translates to:
  /// **'Found a guild'**
  String get createGuild;

  /// No description provided for @createGuildBody.
  ///
  /// In en, this message translates to:
  /// **'Founding a guild costs {fee}.'**
  String createGuildBody(String fee);

  /// No description provided for @createGuildFee.
  ///
  /// In en, this message translates to:
  /// **'Found guild ({fee})'**
  String createGuildFee(String fee);

  /// No description provided for @guildName.
  ///
  /// In en, this message translates to:
  /// **'Guild name'**
  String get guildName;

  /// No description provided for @guildTag.
  ///
  /// In en, this message translates to:
  /// **'Tag'**
  String get guildTag;

  /// No description provided for @guildTagHint.
  ///
  /// In en, this message translates to:
  /// **'2–5 letters or digits'**
  String get guildTagHint;

  /// No description provided for @guildMembers.
  ///
  /// In en, this message translates to:
  /// **'{count}/{cap} members'**
  String guildMembers(int count, int cap);

  /// No description provided for @disbandScheduled.
  ///
  /// In en, this message translates to:
  /// **'Guild disbands in {time}'**
  String disbandScheduled(String time);

  /// No description provided for @motd.
  ///
  /// In en, this message translates to:
  /// **'Message of the day'**
  String get motd;

  /// No description provided for @editMotd.
  ///
  /// In en, this message translates to:
  /// **'Edit message of the day'**
  String get editMotd;

  /// No description provided for @setRank.
  ///
  /// In en, this message translates to:
  /// **'Change rank'**
  String get setRank;

  /// No description provided for @setRankFor.
  ///
  /// In en, this message translates to:
  /// **'New rank for {name}'**
  String setRankFor(String name);

  /// No description provided for @transferTitle.
  ///
  /// In en, this message translates to:
  /// **'Make {name} guild leader?'**
  String transferTitle(String name);

  /// No description provided for @transferBody.
  ///
  /// In en, this message translates to:
  /// **'You become an officer. This can\'t be undone by you.'**
  String get transferBody;

  /// No description provided for @kickTitle.
  ///
  /// In en, this message translates to:
  /// **'Remove {name} from the guild?'**
  String kickTitle(String name);

  /// No description provided for @leaveGuildTitle.
  ///
  /// In en, this message translates to:
  /// **'Leave the guild?'**
  String get leaveGuildTitle;

  /// No description provided for @leaveGuildBody.
  ///
  /// In en, this message translates to:
  /// **'You\'ll need a new invitation to rejoin.'**
  String get leaveGuildBody;

  /// No description provided for @disbandGuildTitle.
  ///
  /// In en, this message translates to:
  /// **'Disband the guild?'**
  String get disbandGuildTitle;

  /// No description provided for @disbandGuildBody.
  ///
  /// In en, this message translates to:
  /// **'The guild disbands in 48 hours; you can cancel until then.'**
  String get disbandGuildBody;

  /// No description provided for @cancelDisband.
  ///
  /// In en, this message translates to:
  /// **'Cancel disbanding'**
  String get cancelDisband;

  /// No description provided for @errAlreadyInGroup.
  ///
  /// In en, this message translates to:
  /// **'They\'re already in a group.'**
  String get errAlreadyInGroup;

  /// No description provided for @errNotLeader.
  ///
  /// In en, this message translates to:
  /// **'Only the leader can do that.'**
  String get errNotLeader;

  /// No description provided for @errGroupFull.
  ///
  /// In en, this message translates to:
  /// **'The group is full.'**
  String get errGroupFull;

  /// No description provided for @errInviteInvalid.
  ///
  /// In en, this message translates to:
  /// **'That invitation has expired.'**
  String get errInviteInvalid;

  /// No description provided for @errAlreadyInvited.
  ///
  /// In en, this message translates to:
  /// **'An invitation is already pending.'**
  String get errAlreadyInvited;

  /// No description provided for @errGuildForbidden.
  ///
  /// In en, this message translates to:
  /// **'Your rank can\'t do that.'**
  String get errGuildForbidden;

  /// No description provided for @errLeaderCannotLeave.
  ///
  /// In en, this message translates to:
  /// **'Hand over leadership or disband the guild first.'**
  String get errLeaderCannotLeave;

  /// No description provided for @errInvalidGuildName.
  ///
  /// In en, this message translates to:
  /// **'3–24 letters or digits.'**
  String get errInvalidGuildName;

  /// No description provided for @errInvalidTag.
  ///
  /// In en, this message translates to:
  /// **'The tag must be 2–5 letters or digits.'**
  String get errInvalidTag;

  /// No description provided for @errTagTaken.
  ///
  /// In en, this message translates to:
  /// **'That tag is taken.'**
  String get errTagTaken;

  /// No description provided for @progress.
  ///
  /// In en, this message translates to:
  /// **'Achievements & mastery'**
  String get progress;

  /// No description provided for @tabAchievements.
  ///
  /// In en, this message translates to:
  /// **'Achievements'**
  String get tabAchievements;

  /// No description provided for @tabMastery.
  ///
  /// In en, this message translates to:
  /// **'Mastery'**
  String get tabMastery;

  /// No description provided for @tabReputation.
  ///
  /// In en, this message translates to:
  /// **'Reputation'**
  String get tabReputation;

  /// No description provided for @tabTitles.
  ///
  /// In en, this message translates to:
  /// **'Titles'**
  String get tabTitles;

  /// No description provided for @achievementPoints.
  ///
  /// In en, this message translates to:
  /// **'{points} achievement points'**
  String achievementPoints(int points);

  /// No description provided for @pointsShort.
  ///
  /// In en, this message translates to:
  /// **'{points} pts'**
  String pointsShort(int points);

  /// No description provided for @rewardTitle.
  ///
  /// In en, this message translates to:
  /// **'Title: {title}'**
  String rewardTitle(String title);

  /// No description provided for @masteryTier.
  ///
  /// In en, this message translates to:
  /// **'Tier {tier}/{max}'**
  String masteryTier(int tier, int max);

  /// No description provided for @repDailyLeft.
  ///
  /// In en, this message translates to:
  /// **'+{amount} today'**
  String repDailyLeft(int amount);

  /// No description provided for @noTitle.
  ///
  /// In en, this message translates to:
  /// **'No title'**
  String get noTitle;

  /// No description provided for @noTitlesYet.
  ///
  /// In en, this message translates to:
  /// **'Titles unlock through achievements.'**
  String get noTitlesYet;

  /// No description provided for @errTitleNotOwned.
  ///
  /// In en, this message translates to:
  /// **'You haven\'t unlocked that title.'**
  String get errTitleNotOwned;

  /// No description provided for @notifications.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get notifications;

  /// No description provided for @markAllRead.
  ///
  /// In en, this message translates to:
  /// **'Mark all as read'**
  String get markAllRead;

  /// No description provided for @inboxNothing.
  ///
  /// In en, this message translates to:
  /// **'No notifications yet.'**
  String get inboxNothing;

  /// No description provided for @notifSettings.
  ///
  /// In en, this message translates to:
  /// **'Notification settings'**
  String get notifSettings;

  /// No description provided for @muteAll.
  ///
  /// In en, this message translates to:
  /// **'Pause push notifications'**
  String get muteAll;

  /// No description provided for @muteHint.
  ///
  /// In en, this message translates to:
  /// **'For 24 hours. The inbox keeps filling.'**
  String get muteHint;

  /// No description provided for @mutedFor.
  ///
  /// In en, this message translates to:
  /// **'Paused for {time}'**
  String mutedFor(String time);

  /// No description provided for @modePushInbox.
  ///
  /// In en, this message translates to:
  /// **'Push + inbox'**
  String get modePushInbox;

  /// No description provided for @modeInbox.
  ///
  /// In en, this message translates to:
  /// **'Inbox only'**
  String get modeInbox;

  /// No description provided for @modeSilent.
  ///
  /// In en, this message translates to:
  /// **'Silent'**
  String get modeSilent;

  /// No description provided for @modeOff.
  ///
  /// In en, this message translates to:
  /// **'Off'**
  String get modeOff;

  /// No description provided for @nkActivity.
  ///
  /// In en, this message translates to:
  /// **'Activity finished'**
  String get nkActivity;

  /// No description provided for @nkDungeon.
  ///
  /// In en, this message translates to:
  /// **'Dungeon run finished'**
  String get nkDungeon;

  /// No description provided for @nkTravel.
  ///
  /// In en, this message translates to:
  /// **'Arrived'**
  String get nkTravel;

  /// No description provided for @nkWorldEvent.
  ///
  /// In en, this message translates to:
  /// **'World event started'**
  String get nkWorldEvent;

  /// No description provided for @nkMarketSold.
  ///
  /// In en, this message translates to:
  /// **'Item sold'**
  String get nkMarketSold;

  /// No description provided for @nkMail.
  ///
  /// In en, this message translates to:
  /// **'New mail'**
  String get nkMail;

  /// No description provided for @nkPartyInvite.
  ///
  /// In en, this message translates to:
  /// **'Party invitation'**
  String get nkPartyInvite;

  /// No description provided for @nkGuildInvite.
  ///
  /// In en, this message translates to:
  /// **'Guild invitation'**
  String get nkGuildInvite;

  /// No description provided for @nkAchievement.
  ///
  /// In en, this message translates to:
  /// **'Achievement unlocked'**
  String get nkAchievement;

  /// No description provided for @nkSystem.
  ///
  /// In en, this message translates to:
  /// **'Announcement'**
  String get nkSystem;

  /// No description provided for @timeAgo.
  ///
  /// In en, this message translates to:
  /// **'{time} ago'**
  String timeAgo(String time);

  /// No description provided for @segEvents.
  ///
  /// In en, this message translates to:
  /// **'Events'**
  String get segEvents;

  /// No description provided for @eventsEmpty.
  ///
  /// In en, this message translates to:
  /// **'No world events scheduled right now.'**
  String get eventsEmpty;

  /// No description provided for @worldBoss.
  ///
  /// In en, this message translates to:
  /// **'World boss'**
  String get worldBoss;

  /// No description provided for @invasion.
  ///
  /// In en, this message translates to:
  /// **'Invasion'**
  String get invasion;

  /// No description provided for @startsIn.
  ///
  /// In en, this message translates to:
  /// **'Starts in {time}'**
  String startsIn(String time);

  /// No description provided for @endsIn.
  ///
  /// In en, this message translates to:
  /// **'Ends in {time}'**
  String endsIn(String time);

  /// No description provided for @regionOnly.
  ///
  /// In en, this message translates to:
  /// **'{region} only'**
  String regionOnly(String region);

  /// No description provided for @eventGoal.
  ///
  /// In en, this message translates to:
  /// **'Goal'**
  String get eventGoal;

  /// No description provided for @enlist.
  ///
  /// In en, this message translates to:
  /// **'Enlist'**
  String get enlist;

  /// No description provided for @enlistHint.
  ///
  /// In en, this message translates to:
  /// **'Your current power is locked in when you enlist. You score for every minute the event runs; rewards arrive when it ends.'**
  String get enlistHint;

  /// No description provided for @enlistedPower.
  ///
  /// In en, this message translates to:
  /// **'Enlisted · power {power}'**
  String enlistedPower(int power);

  /// No description provided for @myContribution.
  ///
  /// In en, this message translates to:
  /// **'Your score {score} · reward tier {tier}'**
  String myContribution(int score, int tier);

  /// No description provided for @eventSucceeded.
  ///
  /// In en, this message translates to:
  /// **'Victory'**
  String get eventSucceeded;

  /// No description provided for @eventFailed.
  ///
  /// In en, this message translates to:
  /// **'Failed'**
  String get eventFailed;

  /// No description provided for @tokenBalance.
  ///
  /// In en, this message translates to:
  /// **'{tokens} event tokens'**
  String tokenBalance(int tokens);

  /// No description provided for @claimDaily.
  ///
  /// In en, this message translates to:
  /// **'Claim {tokens}'**
  String claimDaily(int tokens);

  /// No description provided for @claimedToday.
  ///
  /// In en, this message translates to:
  /// **'Claimed'**
  String get claimedToday;

  /// No description provided for @eventTokensGained.
  ///
  /// In en, this message translates to:
  /// **'+{tokens} event tokens'**
  String eventTokensGained(int tokens);

  /// No description provided for @errEventNotOpen.
  ///
  /// In en, this message translates to:
  /// **'Enlistment for this event is closed.'**
  String get errEventNotOpen;

  /// No description provided for @errWrongRegion.
  ///
  /// In en, this message translates to:
  /// **'Travel to the event\'s region first.'**
  String get errWrongRegion;

  /// No description provided for @errAlreadyEnlisted.
  ///
  /// In en, this message translates to:
  /// **'You\'re already enlisted.'**
  String get errAlreadyEnlisted;

  /// No description provided for @errAlreadyClaimed.
  ///
  /// In en, this message translates to:
  /// **'Already claimed today.'**
  String get errAlreadyClaimed;

  /// No description provided for @errEventNotActive.
  ///
  /// In en, this message translates to:
  /// **'This event isn\'t running.'**
  String get errEventNotActive;

  /// No description provided for @gatherYield.
  ///
  /// In en, this message translates to:
  /// **'{activity} yield'**
  String gatherYield(String activity);

  /// No description provided for @perHour.
  ///
  /// In en, this message translates to:
  /// **'{value}/h'**
  String perHour(String value);

  /// No description provided for @quests.
  ///
  /// In en, this message translates to:
  /// **'Quests'**
  String get quests;

  /// No description provided for @questsInProgress.
  ///
  /// In en, this message translates to:
  /// **'In progress'**
  String get questsInProgress;

  /// No description provided for @questsNoneOpen.
  ///
  /// In en, this message translates to:
  /// **'No open quests. Pick one below.'**
  String get questsNoneOpen;

  /// No description provided for @questsStory.
  ///
  /// In en, this message translates to:
  /// **'Story'**
  String get questsStory;

  /// No description provided for @questsSide.
  ///
  /// In en, this message translates to:
  /// **'Side quests'**
  String get questsSide;

  /// No description provided for @questsDaily.
  ///
  /// In en, this message translates to:
  /// **'Daily'**
  String get questsDaily;

  /// No description provided for @resetsIn.
  ///
  /// In en, this message translates to:
  /// **'Resets in {time}'**
  String resetsIn(String time);

  /// No description provided for @turnIn.
  ///
  /// In en, this message translates to:
  /// **'Turn in'**
  String get turnIn;

  /// No description provided for @questRewarded.
  ///
  /// In en, this message translates to:
  /// **'Quest complete: {rewards}'**
  String questRewarded(String rewards);

  /// No description provided for @abandonQuestTitle.
  ///
  /// In en, this message translates to:
  /// **'Abandon quest?'**
  String get abandonQuestTitle;

  /// No description provided for @abandonQuestBody.
  ///
  /// In en, this message translates to:
  /// **'Progress is lost. You can take the quest again later.'**
  String get abandonQuestBody;

  /// No description provided for @abandonDailyBody.
  ///
  /// In en, this message translates to:
  /// **'Progress is lost and this daily can\'t be taken again today.'**
  String get abandonDailyBody;

  /// No description provided for @objGather.
  ///
  /// In en, this message translates to:
  /// **'Gather ({activity})'**
  String objGather(String activity);

  /// No description provided for @objKillAny.
  ///
  /// In en, this message translates to:
  /// **'Defeat enemies'**
  String get objKillAny;

  /// No description provided for @objKill.
  ///
  /// In en, this message translates to:
  /// **'Defeat {mob}'**
  String objKill(String mob);

  /// No description provided for @objBoss.
  ///
  /// In en, this message translates to:
  /// **'Defeat a boss'**
  String get objBoss;

  /// No description provided for @objDungeonAny.
  ///
  /// In en, this message translates to:
  /// **'Clear a dungeon'**
  String get objDungeonAny;

  /// No description provided for @objDungeon.
  ///
  /// In en, this message translates to:
  /// **'Clear {dungeon}'**
  String objDungeon(String dungeon);

  /// No description provided for @objCraftAny.
  ///
  /// In en, this message translates to:
  /// **'Craft items'**
  String get objCraftAny;

  /// No description provided for @objCraft.
  ///
  /// In en, this message translates to:
  /// **'Craft {item}'**
  String objCraft(String item);

  /// No description provided for @evQuestComplete.
  ///
  /// In en, this message translates to:
  /// **'Quest objectives complete — turn it in!'**
  String get evQuestComplete;

  /// No description provided for @errQuestLocked.
  ///
  /// In en, this message translates to:
  /// **'Finish the previous story quest first.'**
  String get errQuestLocked;

  /// No description provided for @errQuestTaken.
  ///
  /// In en, this message translates to:
  /// **'You already took this quest.'**
  String get errQuestTaken;

  /// No description provided for @errQuestNotOffered.
  ///
  /// In en, this message translates to:
  /// **'That daily isn\'t on today\'s board.'**
  String get errQuestNotOffered;

  /// No description provided for @errQuestCap.
  ///
  /// In en, this message translates to:
  /// **'Too many open quests. Finish or abandon one.'**
  String get errQuestCap;

  /// No description provided for @errQuestNotActive.
  ///
  /// In en, this message translates to:
  /// **'This quest isn\'t open.'**
  String get errQuestNotActive;

  /// No description provided for @errQuestIncomplete.
  ///
  /// In en, this message translates to:
  /// **'Objectives aren\'t complete yet.'**
  String get errQuestIncomplete;

  /// No description provided for @skills.
  ///
  /// In en, this message translates to:
  /// **'Skills'**
  String get skills;

  /// No description provided for @skillActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get skillActive;

  /// No description provided for @skillPassive.
  ///
  /// In en, this message translates to:
  /// **'Passive'**
  String get skillPassive;

  /// No description provided for @skillStrike.
  ///
  /// In en, this message translates to:
  /// **'{power}% damage'**
  String skillStrike(int power);

  /// No description provided for @skillHeal.
  ///
  /// In en, this message translates to:
  /// **'Heals {power}% HP when below half'**
  String skillHeal(int power);

  /// No description provided for @skillMp.
  ///
  /// In en, this message translates to:
  /// **'{mp} MP'**
  String skillMp(int mp);

  /// No description provided for @skillCooldown.
  ///
  /// In en, this message translates to:
  /// **'{turns, plural, =1{every turn} other{every {turns} turns}}'**
  String skillCooldown(int turns);

  /// No description provided for @termsOfService.
  ///
  /// In en, this message translates to:
  /// **'Terms of Service'**
  String get termsOfService;

  /// No description provided for @privacyPolicy.
  ///
  /// In en, this message translates to:
  /// **'Privacy Policy'**
  String get privacyPolicy;

  /// No description provided for @accountMenu.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get accountMenu;

  /// No description provided for @downloadMyData.
  ///
  /// In en, this message translates to:
  /// **'Download my data'**
  String get downloadMyData;

  /// No description provided for @exportOpened.
  ///
  /// In en, this message translates to:
  /// **'The download link opened in your browser (valid for 15 minutes).'**
  String get exportOpened;

  /// No description provided for @linkOpenFailed.
  ///
  /// In en, this message translates to:
  /// **'Could not open the link.'**
  String get linkOpenFailed;

  /// No description provided for @deleteAccount.
  ///
  /// In en, this message translates to:
  /// **'Delete my account'**
  String get deleteAccount;

  /// No description provided for @deleteAccountBody.
  ///
  /// In en, this message translates to:
  /// **'Your account and all characters will be permanently deleted after 30 days. Signing in again within that time cancels the deletion. Enter your password to confirm.'**
  String get deleteAccountBody;

  /// No description provided for @deleteAccountConfirm.
  ///
  /// In en, this message translates to:
  /// **'Delete account'**
  String get deleteAccountConfirm;

  /// No description provided for @deleteScheduled.
  ///
  /// In en, this message translates to:
  /// **'Your account will be deleted on {date}. To cancel, just sign in before then.'**
  String deleteScheduled(String date);

  /// No description provided for @deletionCancelled.
  ///
  /// In en, this message translates to:
  /// **'Your account deletion request was cancelled.'**
  String get deletionCancelled;

  /// No description provided for @support.
  ///
  /// In en, this message translates to:
  /// **'Support'**
  String get support;

  /// No description provided for @supportIntro.
  ///
  /// In en, this message translates to:
  /// **'Found a bug or need help with your account? Write to us and the team answers here; you get a notification when they do.'**
  String get supportIntro;

  /// No description provided for @supportNew.
  ///
  /// In en, this message translates to:
  /// **'New request'**
  String get supportNew;

  /// No description provided for @supportNone.
  ///
  /// In en, this message translates to:
  /// **'No requests yet.'**
  String get supportNone;

  /// No description provided for @supportSent.
  ///
  /// In en, this message translates to:
  /// **'Request received. We\'ll let you know when there\'s an answer.'**
  String get supportSent;

  /// No description provided for @supportCategory.
  ///
  /// In en, this message translates to:
  /// **'Topic'**
  String get supportCategory;

  /// No description provided for @supportSubject.
  ///
  /// In en, this message translates to:
  /// **'Subject'**
  String get supportSubject;

  /// No description provided for @supportMessage.
  ///
  /// In en, this message translates to:
  /// **'What happened? Step by step if you can.'**
  String get supportMessage;

  /// No description provided for @supportTooShort.
  ///
  /// In en, this message translates to:
  /// **'Subject needs at least 3 characters, the message at least 10.'**
  String get supportTooShort;

  /// No description provided for @supportBug.
  ///
  /// In en, this message translates to:
  /// **'Bug'**
  String get supportBug;

  /// No description provided for @supportAccount.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get supportAccount;

  /// No description provided for @supportAbuse.
  ///
  /// In en, this message translates to:
  /// **'Abuse / report a player'**
  String get supportAbuse;

  /// No description provided for @supportFeedback.
  ///
  /// In en, this message translates to:
  /// **'Suggestion'**
  String get supportFeedback;

  /// No description provided for @supportOther.
  ///
  /// In en, this message translates to:
  /// **'Other'**
  String get supportOther;

  /// No description provided for @supportOpen.
  ///
  /// In en, this message translates to:
  /// **'Open'**
  String get supportOpen;

  /// No description provided for @supportAnswered.
  ///
  /// In en, this message translates to:
  /// **'Answered'**
  String get supportAnswered;

  /// No description provided for @supportClosed.
  ///
  /// In en, this message translates to:
  /// **'Closed'**
  String get supportClosed;

  /// No description provided for @supportReply.
  ///
  /// In en, this message translates to:
  /// **'Team reply'**
  String get supportReply;

  /// No description provided for @nkSupport.
  ///
  /// In en, this message translates to:
  /// **'Support reply'**
  String get nkSupport;

  /// No description provided for @errServerBusy.
  ///
  /// In en, this message translates to:
  /// **'The server is busy right now. Try again in a few seconds.'**
  String get errServerBusy;

  /// No description provided for @notifTabInbox.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get notifTabInbox;

  /// No description provided for @notifTabRewards.
  ///
  /// In en, this message translates to:
  /// **'Rewards'**
  String get notifTabRewards;

  /// No description provided for @rewardsNone.
  ///
  /// In en, this message translates to:
  /// **'No rewards yet. Everything you gather or win in fights is listed here (kept for 30 days).'**
  String get rewardsNone;

  /// No description provided for @rewardFromActivity.
  ///
  /// In en, this message translates to:
  /// **'Gathering: {name}'**
  String rewardFromActivity(String name);

  /// No description provided for @rewardFromHunt.
  ///
  /// In en, this message translates to:
  /// **'Hunt: {name}'**
  String rewardFromHunt(String name);

  /// No description provided for @rewardFromDungeon.
  ///
  /// In en, this message translates to:
  /// **'Dungeon: {name}'**
  String rewardFromDungeon(String name);

  /// No description provided for @rewardLevelUp.
  ///
  /// In en, this message translates to:
  /// **'Level {level}!'**
  String rewardLevelUp(int level);
}

class _AppLocalizationsDelegate extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) => <String>['en', 'tr'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'tr':
      return AppLocalizationsTr();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
