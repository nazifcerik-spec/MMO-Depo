// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Turkish (`tr`).
class AppLocalizationsTr extends AppLocalizations {
  AppLocalizationsTr([String locale = 'tr']) : super(locale);

  @override
  String get appTitle => 'New World MMO';

  @override
  String get tagline => 'Sen yokken de süren maceralar.';

  @override
  String get emailOrUsername => 'E-posta veya kullanıcı adı';

  @override
  String get email => 'E-posta';

  @override
  String get username => 'Kullanıcı adı';

  @override
  String get password => 'Şifre';

  @override
  String get passwordConfirm => 'Şifre (tekrar)';

  @override
  String get showPassword => 'Şifreyi göster';

  @override
  String get hidePassword => 'Şifreyi gizle';

  @override
  String get login => 'Giriş yap';

  @override
  String get register => 'Hesap oluştur';

  @override
  String get noAccount => 'Yeni misin? Kayıt ol';

  @override
  String get haveAccount => 'Hesabın var mı? Giriş yap';

  @override
  String get acceptTerms => '13 yaşından büyüğüm; Kullanım Şartları\'nı ve Gizlilik Politikası\'nı kabul ediyorum';

  @override
  String get fieldRequired => 'Zorunlu';

  @override
  String get invalidEmail => 'Geçerli bir e-posta gir';

  @override
  String get usernameRule => '3–32 harf, rakam, _ veya .';

  @override
  String get passwordMin => 'En az 8 karakter';

  @override
  String get passwordMismatch => 'Şifreler eşleşmiyor';

  @override
  String get termsRequired => 'Devam etmek için koşulları kabul et';

  @override
  String get logout => 'Çıkış yap';

  @override
  String get retry => 'Tekrar dene';

  @override
  String get offline => 'Sunucuya ulaşılamıyor. Bağlantını kontrol et.';

  @override
  String get characters => 'Karakterler';

  @override
  String slotsUsed(int used, int cap) {
    return '$used/$cap slot';
  }

  @override
  String get newCharacter => 'Yeni karakter';

  @override
  String get noCharacters => 'Henüz karakterin yok. İlk kahramanını yarat.';

  @override
  String get create => 'Oluştur';

  @override
  String get characterName => 'İsim';

  @override
  String get nameRule => '3–16 karakter, harfle başlar (A–Z, 0–9, _)';

  @override
  String get race => 'Irk';

  @override
  String get characterClass => 'Sınıf';

  @override
  String get raceHuman => 'İnsan';

  @override
  String get raceElf => 'Elf';

  @override
  String get raceDwarf => 'Cüce';

  @override
  String get raceOrc => 'Ork';

  @override
  String get classWarrior => 'Savaşçı';

  @override
  String get classMage => 'Büyücü';

  @override
  String get classHunter => 'Avcı';

  @override
  String get classPriest => 'Rahip';

  @override
  String get classRogue => 'Haydut';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Sv $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Sv $level';
  }

  @override
  String get hp => 'CAN';

  @override
  String get mp => 'MANA';

  @override
  String get stamina => 'DAY';

  @override
  String get energy => 'ENJ';

  @override
  String get xp => 'TP';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => 'Nitelikler';

  @override
  String get combat => 'Savaş';

  @override
  String get statStr => 'GÜÇ';

  @override
  String get statAgi => 'ÇEV';

  @override
  String get statInt => 'ZEK';

  @override
  String get statVit => 'DAY';

  @override
  String get statCha => 'KAR';

  @override
  String get statLuck => 'ŞANS';

  @override
  String get attack => 'Saldırı';

  @override
  String get defense => 'Savunma';

  @override
  String get crit => 'Kritik';

  @override
  String get hit => 'İsabet';

  @override
  String percent(int value) {
    return '%$value';
  }

  @override
  String get gold => 'Altın';

  @override
  String get switchCharacter => 'Karakter değiştir';

  @override
  String get statePendingDelete => 'Silinme bekliyor';

  @override
  String get stateFrozen => 'Donduruldu';

  @override
  String get errInvalidCredentials => 'E-posta/kullanıcı adı veya şifre hatalı.';

  @override
  String get errAccountLocked => 'Hesabın geçici olarak kilitli.';

  @override
  String get errAccountBanned => 'Bu hesap yasaklandı.';

  @override
  String get errAccountExists => 'Bu e-posta veya kullanıcı adı zaten alınmış.';

  @override
  String get errRateLimited => 'Çok fazla deneme. Biraz bekle.';

  @override
  String get errNameTaken => 'Bu isim zaten alınmış.';

  @override
  String get errSlotCap => 'Tüm karakter slotları dolu.';

  @override
  String get errNotPlayable => 'Bu karakter şu an oynanamaz.';

  @override
  String get errSessionExpired => 'Oturumun sona erdi. Lütfen tekrar giriş yap.';

  @override
  String errGeneric(String code) {
    return 'Bir şeyler ters gitti ($code).';
  }

  @override
  String get navCharacter => 'Karakter';

  @override
  String get navActivity => 'Toplama';

  @override
  String get actMining => 'Madencilik';

  @override
  String get actFishing => 'Balıkçılık';

  @override
  String get actWoodcutting => 'Oduncluk';

  @override
  String get actHerbalism => 'Şifacılık';

  @override
  String get actHunting => 'Avcılık';

  @override
  String get activityChoose => 'Bir aktivite seç';

  @override
  String get durationLabel => 'Süre';

  @override
  String startActivity(String name) {
    return '$name başlat';
  }

  @override
  String get collect => 'Topla';

  @override
  String get abandon => 'Bırak';

  @override
  String get cancel => 'Vazgeç';

  @override
  String get soFar => 'Şu ana kadar toplanan';

  @override
  String get activityDone => 'Tamamlandı — toplamaya hazır';

  @override
  String remaining(String time) {
    return '$time kaldı';
  }

  @override
  String durHm(int h, String m) {
    return '${h}sa ${m}dk';
  }

  @override
  String durM(int m) {
    return '${m}dk';
  }

  @override
  String durS(int s) {
    return '${s}sn';
  }

  @override
  String get collectEarlyTitle => 'Erken toplansın mı?';

  @override
  String get collectEarlyBody => 'Şu ana kadar toplananı alırsın ve aktivite biter.';

  @override
  String get abandonTitle => 'Aktivite bırakılsın mı?';

  @override
  String get abandonBody => 'Şu ana kadar toplanan her şey kaybolur.';

  @override
  String get lootCollect => 'Harika!';

  @override
  String get lootNothing => 'Bu sefer bir şey çıkmadı.';

  @override
  String xpGain(int xp) {
    return '+$xp TP';
  }

  @override
  String get rarityCommon => 'Sıradan';

  @override
  String get rarityUncommon => 'Seçkin';

  @override
  String get rarityRare => 'Nadir';

  @override
  String get rarityEpic => 'Destansı';

  @override
  String get rarityLegendary => 'Efsanevi';

  @override
  String get rarityMythic => 'Mitik';

  @override
  String get errSlotBusy => 'Zaten devam eden bir aktivite var.';

  @override
  String get errNotActive => 'Bu aktivite zaten sona ermiş.';

  @override
  String get evAchievement => 'Başarım açıldı!';

  @override
  String get evMasteryTier => 'Ustalık kademesi yükseldi!';

  @override
  String get evMarketSold => 'Pazardaki bir ilanın satıldı.';

  @override
  String get evMailReceived => 'Yeni postan var.';

  @override
  String get evPartyInvite => 'Bir gruba davet edildin.';

  @override
  String get evGuildInvite => 'Bir loncaya davet edildin.';

  @override
  String get navBag => 'Çanta';

  @override
  String get equipment => 'Ekipman';

  @override
  String get bag => 'Çanta';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap kg';
  }

  @override
  String get slotHelm => 'Miğfer';

  @override
  String get slotNecklace => 'Kolye';

  @override
  String get slotChest => 'Göğüs';

  @override
  String get slotGloves => 'Eldiven';

  @override
  String get slotLegs => 'Bacak';

  @override
  String get slotBoots => 'Çizme';

  @override
  String get slotMainhand => 'Ana el';

  @override
  String get slotOffhand => 'Yan el';

  @override
  String get slotRing => 'Yüzük';

  @override
  String get slotTrinket => 'Tılsım';

  @override
  String get equip => 'Kuşan';

  @override
  String get unequip => 'Çıkar';

  @override
  String get sell => 'Sat';

  @override
  String get destroy => 'Yok et';

  @override
  String sellTitle(String name) {
    return '$name satılsın mı?';
  }

  @override
  String get sellBody => '2 saat içinde tüccardan geri alabilirsin.';

  @override
  String destroyTitle(String name) {
    return '$name yok edilsin mi?';
  }

  @override
  String get destroyBody => 'Bu işlem geri alınamaz.';

  @override
  String get soldFor => 'Satıldı:';

  @override
  String get quantity => 'Adet';

  @override
  String get requiresLevel => 'Gereken seviye';

  @override
  String get weight => 'Ağırlık';

  @override
  String weightGrams(int g) {
    return '$g g';
  }

  @override
  String get sellValue => 'Satış değeri';

  @override
  String get bound => 'Ruha bağlı';

  @override
  String get itemLocked => 'Bu eşya bir takas, posta veya ilanda tutuluyor.';

  @override
  String get errBagFull => 'Çantan dolu.';

  @override
  String get errItemLocked => 'Bu eşya kilitli.';

  @override
  String get errClassRestriction => 'Sınıfın bu eşyayı kullanamaz.';

  @override
  String get errLevelRequirement => 'Seviyen yetersiz.';

  @override
  String get errNotEquippable => 'Bu eşya kuşanılamaz.';

  @override
  String get navAdventure => 'Macera';

  @override
  String get segHunt => 'Av';

  @override
  String get segDungeon => 'Zindan';

  @override
  String get segTravel => 'Yolculuk';

  @override
  String get huntChoose => 'Avını seç';

  @override
  String get startHunt => 'Avı başlat';

  @override
  String get huntDone => 'Av bitti — ganimetini topla';

  @override
  String get huntEarlyBody => 'Şu ana kadar kazanılan dövüşleri alırsın ve av biter.';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'CAN $hp · SAL $atk · SAV $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '$won galibiyet · $lost yenilgi';
  }

  @override
  String levelUp(int level) {
    return 'Seviye atladın! Artık seviye $level';
  }

  @override
  String levelReq(int level) {
    return 'Sv $level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '$count dövüş · boss $boss';
  }

  @override
  String get dungeonEnter => 'Gir';

  @override
  String get dungeonDone => 'Zindan bitti — ödülleri al';

  @override
  String get dungeonWaitHint => 'Ödüller zindan bitince açılır.';

  @override
  String get dungeonCleared => 'Temizlendi!';

  @override
  String get dungeonFailed => 'Grup düştü';

  @override
  String get bossKilled => 'Boss yenildi';

  @override
  String get bonus => 'Tamamlama bonusu';

  @override
  String get youAreIn => 'Bulunduğun yer';

  @override
  String get routes => 'Yollar';

  @override
  String get noRoutes => 'Buradan gidilecek yol yok.';

  @override
  String get travelTo => 'Git';

  @override
  String travelling(String place) {
    return '$place yolunda';
  }

  @override
  String get arrive => 'Var';

  @override
  String get arrived => 'Vardın';

  @override
  String get cancelTravel => 'Geri dön';

  @override
  String get kindCity => 'Şehir';

  @override
  String get kindZone => 'Bölge';

  @override
  String get kindWilderness => 'Vahşi doğa';

  @override
  String get kindDungeon => 'Zindan';

  @override
  String get kindContinent => 'Kıta';

  @override
  String get safeZone => 'güvenli';

  @override
  String get errBusy => 'Zaten devam eden bir işin var.';

  @override
  String get errNotReady => 'Henüz bitmedi.';

  @override
  String get errWeeklyLocked => 'Bu zindanı bu hafta zaten temizledin.';

  @override
  String get errNoRoute => 'Oraya doğrudan yol yok.';

  @override
  String get navWorkshop => 'Atölye';

  @override
  String get segCraft => 'Üretim';

  @override
  String get segShop => 'Dükkân';

  @override
  String get profBlacksmithing => 'Demircilik';

  @override
  String get profAlchemy => 'Simya';

  @override
  String get profCooking => 'Aşçılık';

  @override
  String get learn => 'Öğren';

  @override
  String get learnHint => 'Tariflerini açmak için bu mesleği öğren.';

  @override
  String get craft => 'Üret';

  @override
  String crafted(String name, int count) {
    return '$name ×$count üretildi';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession artık seviye $level';
  }

  @override
  String requiresSkill(int level) {
    return 'Gereken beceri $level';
  }

  @override
  String get buy => 'Satın al';

  @override
  String bought(String name) {
    return '$name satın alındı';
  }

  @override
  String get buybacks => 'Geri satın al';

  @override
  String get buybackHint => 'Sattığın eşyaları 2 saat içinde geri alabilirsin.';

  @override
  String get errInsufficientGold => 'Yeterli altının yok.';

  @override
  String get errInsufficientMaterials => 'Malzemen eksik.';

  @override
  String get errSkillTooLow => 'Becerin bu tarif için yetersiz.';

  @override
  String get errProfessionNotLearned => 'Önce mesleği öğren.';

  @override
  String get errPriceChanged => 'Fiyat değişti — kontrol edip tekrar dene.';

  @override
  String get errCannotSell => 'Bu eşya satılamaz.';

  @override
  String get errBuybackExpired => 'Geri alma süresi doldu.';

  @override
  String get segMarket => 'Pazar';

  @override
  String get classWeapon => 'Silahlar';

  @override
  String get classArmor => 'Zırhlar';

  @override
  String get classConsumable => 'Tüketimlik';

  @override
  String get classMaterial => 'Malzemeler';

  @override
  String get filterAll => 'Tümü';

  @override
  String get sortCheapest => 'En ucuz';

  @override
  String get sortNewest => 'En yeni';

  @override
  String get marketEmpty => 'Eşleşen ilan yok.';

  @override
  String listingLine(String seller, String price) {
    return '$seller · tanesi $price';
  }

  @override
  String get loadMore => 'Daha fazla';

  @override
  String myListings(int count) {
    return 'İlanlarım ($count)';
  }

  @override
  String get cancelListing => 'İlanı kaldır';

  @override
  String cancelListingTitle(String name) {
    return '$name ilanı kaldırılsın mı?';
  }

  @override
  String get cancelListingBody =>
      'Eşya çantana döner. İlan ücretinin bir kısmı yalnızca ilk 6 saat içinde iade edilir.';

  @override
  String buyListingTitle(String name) {
    return '$name satın alınsın mı?';
  }

  @override
  String buyListingBody(String price, String seller) {
    return '$seller oyuncusuna $price ödenecek.';
  }

  @override
  String get listOnMarket => 'Pazara koy';

  @override
  String listed(String name) {
    return '$name pazara kondu';
  }

  @override
  String quantityOf(int count, int max) {
    return 'Adet: $count / $max';
  }

  @override
  String get unitPrice => 'Tane fiyatı';

  @override
  String suggestedPrice(String price) {
    return 'Pazar fiyatı: $price';
  }

  @override
  String get listingFee => 'İlan ücreti (iade edilmez)';

  @override
  String get youReceive => 'Satılınca eline geçen';

  @override
  String get errPriceOutOfRange => 'Bu fiyat izin verilen aralığın dışında.';

  @override
  String get errListingCap => 'Çok fazla aktif ilanın var.';

  @override
  String get errSelfBuy => 'Kendi ilanını satın alamazsın.';

  @override
  String get errListingUnavailable => 'Bu ilan artık mevcut değil.';

  @override
  String get errNotTradable => 'Bu eşya takas edilemez.';

  @override
  String get social => 'Sosyal';

  @override
  String get tabChat => 'Sohbet';

  @override
  String get tabFriends => 'Arkadaşlar';

  @override
  String get tabMail => 'Posta';

  @override
  String get chanGlobal => 'Dünya';

  @override
  String get chanGuild => 'Lonca';

  @override
  String get chanParty => 'Grup';

  @override
  String get chatEmpty => 'Henüz mesaj yok. Merhaba de!';

  @override
  String chatHint(String channel) {
    return '$channel kanalına yaz';
  }

  @override
  String get send => 'Gönder';

  @override
  String get addFriend => 'Ekle';

  @override
  String friendAdded(String name) {
    return '$name ile artık arkadaşsınız';
  }

  @override
  String friendRequested(String name) {
    return '$name oyuncusuna arkadaşlık isteği gönderildi';
  }

  @override
  String get friendRequests => 'Arkadaşlık istekleri';

  @override
  String friendsCount(int count) {
    return 'Arkadaşlar ($count)';
  }

  @override
  String get noFriends => 'Henüz arkadaşın yok — isimle birini ekle.';

  @override
  String get more => 'Diğer';

  @override
  String get whisper => 'Fısılda';

  @override
  String get remove => 'Çıkar';

  @override
  String get block => 'Engelle';

  @override
  String get unblock => 'Engeli kaldır';

  @override
  String get blockedList => 'Engellenenler';

  @override
  String removeFriendTitle(String name) {
    return '$name çıkarılsın mı?';
  }

  @override
  String get removeFriendBody => 'Daha sonra yeniden arkadaşlık isteği gönderebilirsin.';

  @override
  String blockTitle(String name) {
    return '$name engellensin mi?';
  }

  @override
  String get blockBody => 'Sana fısıldayamaz, sohbet mesajları gizlenir, postaları düşer. Arkadaşlığınız da sona erer.';

  @override
  String get accept => 'Kabul et';

  @override
  String get decline => 'Reddet';

  @override
  String get newMail => 'Posta yaz';

  @override
  String get inboxEmpty => 'Posta kutun boş.';

  @override
  String get noSubject => '(konu yok)';

  @override
  String mailFrom(String name) {
    return 'Gönderen: $name';
  }

  @override
  String get attachments => 'Ekler';

  @override
  String get claim => 'Ekleri al';

  @override
  String get mailClaimed => 'Ekler çantana ve cüzdanına aktarıldı';

  @override
  String get deleteMailTitle => 'Bu posta silinsin mi?';

  @override
  String get deleteMailBody => 'Silinen posta geri getirilemez.';

  @override
  String mailSent(String postage) {
    return 'Posta gönderildi (posta ücreti $postage)';
  }

  @override
  String get mailTo => 'Alıcı (karakter adı)';

  @override
  String get mailSubject => 'Konu';

  @override
  String get mailBody => 'Mesaj';

  @override
  String get attachGold => 'Altın ekle';

  @override
  String get errAlreadyFriends => 'Zaten arkadaşsınız.';

  @override
  String get errAlreadyRequested => 'Bekleyen bir arkadaşlık isteği var.';

  @override
  String get errBlocked => 'Bu oyuncuyla etkileşime geçemezsin.';

  @override
  String get errFriendCap => 'Arkadaş listen dolu.';

  @override
  String get errTargetNotFound => 'Bu isimde bir karakter yok.';

  @override
  String get errSelfTarget => 'Kendini hedef alamazsın.';

  @override
  String get errInboxFull => 'Alıcının posta kutusu dolu.';

  @override
  String get errTooLong => 'Bu çok uzun.';

  @override
  String get errEmpty => 'Önce bir şey yaz.';

  @override
  String get errNotInChannel => 'Bu grubun üyesi değilsin.';

  @override
  String get errMuted => 'Susturuldun, şu an sohbet edemezsin.';

  @override
  String get errHasAttachments => 'Silmeden önce ekleri al.';

  @override
  String get useItem => 'Kullan';

  @override
  String restoresHp(int amount) {
    return '$amount can yeniler';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return '$name kullanıldı · Can $hp/$max';
  }

  @override
  String get errHpFull => 'Canın zaten dolu.';

  @override
  String get errInCombat => 'Önce avını ya da zindanını bitir.';

  @override
  String get errNotUsable => 'Bu eşya kullanılamaz.';

  @override
  String get errTooWounded => 'Savaşamayacak kadar yaralısın. Dinlen ya da iksir iç.';

  @override
  String get forgotPassword => 'Şifremi unuttum';

  @override
  String get resetHint => 'Hesabının e-posta adresini gir; yeni şifre belirlemen için bir bağlantı gönderelim.';

  @override
  String get resetSent => 'Bu e-postaya ait bir hesap varsa sıfırlama bağlantısı gönderildi. Bağlantı 1 saat geçerli.';

  @override
  String get sendLink => 'Bağlantı gönder';

  @override
  String get close => 'Kapat';

  @override
  String get changePassword => 'Şifre değiştir';

  @override
  String get currentPassword => 'Mevcut şifre';

  @override
  String get newPassword => 'Yeni şifre';

  @override
  String get save => 'Kaydet';

  @override
  String get passwordChanged => 'Şifren değişti. Diğer cihazlarındaki oturumlar kapatıldı.';

  @override
  String get verifyEmailBanner => 'Lütfen e-posta adresini doğrula.';

  @override
  String get resend => 'Tekrar gönder';

  @override
  String get verificationSent => 'Doğrulama e-postası gönderildi.';

  @override
  String get errEmailNotVerified => 'Karakter oluşturmadan önce e-postanı doğrula.';

  @override
  String get errAlreadyVerified => 'E-postan zaten doğrulanmış.';

  @override
  String get tabParty => 'Grup';

  @override
  String get tabGuild => 'Lonca';

  @override
  String get invite => 'Davet et';

  @override
  String invitedBy(String name) {
    return 'Davet eden: $name';
  }

  @override
  String get noParty => 'Bir grupta değilsin. Grup kur ve arkadaşlarını davet et — dövüşler birlikte daha kolay.';

  @override
  String get createParty => 'Grup kur';

  @override
  String get pendingInvites => 'Davetler';

  @override
  String partyMembers(int count, int cap) {
    return 'Grup ($count/$cap)';
  }

  @override
  String get makeLeader => 'Lider yap';

  @override
  String get kick => 'Çıkar';

  @override
  String get leave => 'Ayrıl';

  @override
  String get disband => 'Dağıt';

  @override
  String get leavePartyTitle => 'Gruptan ayrılınsın mı?';

  @override
  String get leavePartyBody => 'Liderse, grupta en uzun süredir olan üye liderliği devralır.';

  @override
  String get disbandPartyTitle => 'Grup dağıtılsın mı?';

  @override
  String get disbandPartyBody => 'Herkes gruptan çıkarılır.';

  @override
  String get rankLeader => 'Lider';

  @override
  String get rankOfficer => 'Subay';

  @override
  String get rankVeteran => 'Kıdemli';

  @override
  String get rankMember => 'Üye';

  @override
  String get rankInitiate => 'Acemi';

  @override
  String get createGuild => 'Lonca kur';

  @override
  String createGuildBody(String fee) {
    return 'Lonca kurmanın bedeli $fee.';
  }

  @override
  String createGuildFee(String fee) {
    return 'Lonca kur ($fee)';
  }

  @override
  String get guildName => 'Lonca adı';

  @override
  String get guildTag => 'Etiket';

  @override
  String get guildTagHint => '2–5 harf veya rakam';

  @override
  String guildMembers(int count, int cap) {
    return '$count/$cap üye';
  }

  @override
  String disbandScheduled(String time) {
    return 'Lonca $time sonra dağılacak';
  }

  @override
  String get motd => 'Günün mesajı';

  @override
  String get editMotd => 'Günün mesajını düzenle';

  @override
  String get setRank => 'Rütbe değiştir';

  @override
  String setRankFor(String name) {
    return '$name için yeni rütbe';
  }

  @override
  String transferTitle(String name) {
    return '$name lonca lideri yapılsın mı?';
  }

  @override
  String get transferBody => 'Sen subay olursun. Bunu sen geri alamazsın.';

  @override
  String kickTitle(String name) {
    return '$name loncadan çıkarılsın mı?';
  }

  @override
  String get leaveGuildTitle => 'Loncadan ayrılınsın mı?';

  @override
  String get leaveGuildBody => 'Geri dönmek için yeni bir davet gerekir.';

  @override
  String get disbandGuildTitle => 'Lonca dağıtılsın mı?';

  @override
  String get disbandGuildBody => 'Lonca 48 saat sonra dağılır; o zamana kadar iptal edebilirsin.';

  @override
  String get cancelDisband => 'Dağıtmayı iptal et';

  @override
  String get errAlreadyInGroup => 'Zaten bir grupta.';

  @override
  String get errNotLeader => 'Bunu yalnızca lider yapabilir.';

  @override
  String get errGroupFull => 'Grup dolu.';

  @override
  String get errInviteInvalid => 'Bu davetin süresi dolmuş.';

  @override
  String get errAlreadyInvited => 'Bekleyen bir davet zaten var.';

  @override
  String get errGuildForbidden => 'Rütben bunun için yeterli değil.';

  @override
  String get errLeaderCannotLeave => 'Önce liderliği devret ya da loncayı dağıt.';

  @override
  String get errInvalidGuildName => '3–24 harf veya rakam.';

  @override
  String get errInvalidTag => 'Etiket 2–5 harf veya rakam olmalı.';

  @override
  String get errTagTaken => 'Bu etiket alınmış.';

  @override
  String get progress => 'Başarımlar ve ustalık';

  @override
  String get tabAchievements => 'Başarımlar';

  @override
  String get tabMastery => 'Ustalık';

  @override
  String get tabReputation => 'İtibar';

  @override
  String get tabTitles => 'Unvanlar';

  @override
  String achievementPoints(int points) {
    return '$points başarım puanı';
  }

  @override
  String pointsShort(int points) {
    return '$points puan';
  }

  @override
  String rewardTitle(String title) {
    return 'Unvan: $title';
  }

  @override
  String masteryTier(int tier, int max) {
    return 'Kademe $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return 'Bugün +$amount';
  }

  @override
  String get noTitle => 'Unvan yok';

  @override
  String get noTitlesYet => 'Unvanlar başarımlarla açılır.';

  @override
  String get errTitleNotOwned => 'Bu unvanı henüz açmadın.';

  @override
  String get notifications => 'Bildirimler';

  @override
  String get markAllRead => 'Tümünü okundu say';

  @override
  String get inboxNothing => 'Henüz bildirim yok.';

  @override
  String get notifSettings => 'Bildirim ayarları';

  @override
  String get muteAll => 'Anlık bildirimleri duraklat';

  @override
  String get muteHint => '24 saatliğine. Bildirim kutusu dolmaya devam eder.';

  @override
  String mutedFor(String time) {
    return '$time daha duraklatıldı';
  }

  @override
  String get modePushInbox => 'Anlık + kutu';

  @override
  String get modeInbox => 'Yalnız kutu';

  @override
  String get modeSilent => 'Sessiz';

  @override
  String get modeOff => 'Kapalı';

  @override
  String get nkActivity => 'Etkinlik tamamlandı';

  @override
  String get nkDungeon => 'Zindan koşusu bitti';

  @override
  String get nkTravel => 'Varış';

  @override
  String get nkWorldEvent => 'Dünya etkinliği başladı';

  @override
  String get nkMarketSold => 'Eşya satıldı';

  @override
  String get nkMail => 'Yeni posta';

  @override
  String get nkPartyInvite => 'Parti daveti';

  @override
  String get nkGuildInvite => 'Lonca daveti';

  @override
  String get nkAchievement => 'Başarım açıldı';

  @override
  String get nkSystem => 'Duyuru';

  @override
  String timeAgo(String time) {
    return '$time önce';
  }

  @override
  String get segEvents => 'Etkinlik';

  @override
  String get eventsEmpty => 'Şu an planlı dünya etkinliği yok.';

  @override
  String get worldBoss => 'Dünya bossu';

  @override
  String get invasion => 'İstila';

  @override
  String startsIn(String time) {
    return '$time sonra başlar';
  }

  @override
  String endsIn(String time) {
    return '$time sonra biter';
  }

  @override
  String regionOnly(String region) {
    return 'Yalnız $region';
  }

  @override
  String get eventGoal => 'Hedef';

  @override
  String get enlist => 'Katıl';

  @override
  String get enlistHint =>
      'Katıldığın andaki gücün sabitlenir. Etkinlik sürdükçe her dakika puan kazanırsın; ödüller bitince gelir.';

  @override
  String enlistedPower(int power) {
    return 'Katıldın · güç $power';
  }

  @override
  String myContribution(int score, int tier) {
    return 'Puanın $score · ödül kademesi $tier';
  }

  @override
  String get eventSucceeded => 'Zafer';

  @override
  String get eventFailed => 'Başarısız';

  @override
  String tokenBalance(int tokens) {
    return '$tokens etkinlik jetonu';
  }

  @override
  String claimDaily(int tokens) {
    return '$tokens al';
  }

  @override
  String get claimedToday => 'Alındı';

  @override
  String eventTokensGained(int tokens) {
    return '+$tokens etkinlik jetonu';
  }

  @override
  String get errEventNotOpen => 'Bu etkinliğe katılım kapandı.';

  @override
  String get errWrongRegion => 'Önce etkinliğin bölgesine git.';

  @override
  String get errAlreadyEnlisted => 'Zaten katıldın.';

  @override
  String get errAlreadyClaimed => 'Bugün zaten aldın.';

  @override
  String get errEventNotActive => 'Bu etkinlik şu an sürmüyor.';

  @override
  String gatherYield(String activity) {
    return '$activity verimi';
  }

  @override
  String perHour(String value) {
    return '$value/sa';
  }

  @override
  String get quests => 'Görevler';

  @override
  String get questsInProgress => 'Devam eden';

  @override
  String get questsNoneOpen => 'Açık görev yok. Aşağıdan birini seç.';

  @override
  String get questsStory => 'Hikâye';

  @override
  String get questsSide => 'Yan görevler';

  @override
  String get questsDaily => 'Günlük';

  @override
  String resetsIn(String time) {
    return '$time sonra yenilenir';
  }

  @override
  String get turnIn => 'Teslim et';

  @override
  String questRewarded(String rewards) {
    return 'Görev tamamlandı: $rewards';
  }

  @override
  String get abandonQuestTitle => 'Görev bırakılsın mı?';

  @override
  String get abandonQuestBody => 'İlerleme kaybolur. Görevi daha sonra yeniden alabilirsin.';

  @override
  String get abandonDailyBody => 'İlerleme kaybolur ve bu günlük görev bugün yeniden alınamaz.';

  @override
  String objGather(String activity) {
    return 'Topla ($activity)';
  }

  @override
  String get objKillAny => 'Düşman yen';

  @override
  String objKill(String mob) {
    return '$mob yen';
  }

  @override
  String get objBoss => 'Bir boss yen';

  @override
  String get objDungeonAny => 'Bir zindan temizle';

  @override
  String objDungeon(String dungeon) {
    return '$dungeon temizle';
  }

  @override
  String get objCraftAny => 'Eşya üret';

  @override
  String objCraft(String item) {
    return '$item üret';
  }

  @override
  String get evQuestComplete => 'Görev hedefleri tamam — teslim et!';

  @override
  String get errQuestLocked => 'Önce önceki hikâye görevini bitir.';

  @override
  String get errQuestTaken => 'Bu görevi zaten aldın.';

  @override
  String get errQuestNotOffered => 'Bu günlük görev bugünkü panoda yok.';

  @override
  String get errQuestCap => 'Çok fazla açık görev var. Birini bitir ya da bırak.';

  @override
  String get errQuestNotActive => 'Bu görev açık değil.';

  @override
  String get errQuestIncomplete => 'Hedefler henüz tamamlanmadı.';

  @override
  String get skills => 'Yetenekler';

  @override
  String get skillActive => 'Aktif';

  @override
  String get skillPassive => 'Pasif';

  @override
  String skillStrike(int power) {
    return '%$power hasar';
  }

  @override
  String skillHeal(int power) {
    return 'HP yarının altındayken %$power iyileştirir';
  }

  @override
  String skillMp(int mp) {
    return '$mp MP';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: '$turns turda bir', one: 'her tur');
    return '$_temp0';
  }

  @override
  String get termsOfService => 'Kullanım Şartları';

  @override
  String get privacyPolicy => 'Gizlilik Politikası';

  @override
  String get accountMenu => 'Hesap';

  @override
  String get downloadMyData => 'Verilerimi indir';

  @override
  String get exportOpened => 'İndirme bağlantısı tarayıcıda açıldı (15 dakika geçerli).';

  @override
  String get linkOpenFailed => 'Bağlantı açılamadı.';

  @override
  String get deleteAccount => 'Hesabımı sil';

  @override
  String get deleteAccountBody =>
      'Hesabın ve tüm karakterlerin 30 gün sonra kalıcı olarak silinir. Bu süre içinde tekrar giriş yaparsan silme iptal olur. Onaylamak için şifreni gir.';

  @override
  String get deleteAccountConfirm => 'Hesabı sil';

  @override
  String deleteScheduled(String date) {
    return 'Hesabın $date tarihinde silinecek. Vazgeçersen bu tarihten önce giriş yapman yeterli.';
  }

  @override
  String get deletionCancelled => 'Hesap silme talebin iptal edildi.';
}
