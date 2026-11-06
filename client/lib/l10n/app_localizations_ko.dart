// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Korean (`ko`).
class AppLocalizationsKo extends AppLocalizations {
  AppLocalizationsKo([String locale = 'ko']) : super(locale);

  @override
  String get appTitle => 'MeNa World';

  @override
  String get tagline => '접속하지 않는 동안에도 모험은 계속됩니다.';

  @override
  String get emailOrUsername => '이메일 또는 사용자 이름';

  @override
  String get email => '이메일';

  @override
  String get username => '사용자 이름';

  @override
  String get password => '비밀번호';

  @override
  String get passwordConfirm => '비밀번호 확인';

  @override
  String get showPassword => '비밀번호 표시';

  @override
  String get hidePassword => '비밀번호 숨기기';

  @override
  String get login => '로그인';

  @override
  String get register => '계정 만들기';

  @override
  String get noAccount => '처음이신가요? 계정을 만드세요';

  @override
  String get haveAccount => '이미 계정이 있나요? 로그인';

  @override
  String get acceptTerms => '만 13세 이상이며 서비스 약관과 개인정보 처리방침에 동의합니다';

  @override
  String get fieldRequired => '필수 항목';

  @override
  String get invalidEmail => '올바른 이메일을 입력하세요';

  @override
  String get usernameRule => '3~32자의 문자, 숫자, _ 또는 .';

  @override
  String get passwordMin => '8자 이상';

  @override
  String get passwordMismatch => '비밀번호가 일치하지 않습니다';

  @override
  String get termsRequired => '계속하려면 약관에 동의하세요';

  @override
  String get logout => '로그아웃';

  @override
  String get retry => '다시 시도';

  @override
  String get offline => '서버에 연결할 수 없습니다. 연결 상태를 확인하세요.';

  @override
  String get characters => '캐릭터';

  @override
  String slotsUsed(int used, int cap) {
    return '슬롯 $used/$cap';
  }

  @override
  String get newCharacter => '새 캐릭터';

  @override
  String get noCharacters => '아직 캐릭터가 없습니다. 첫 영웅을 만들어 보세요.';

  @override
  String get create => '만들기';

  @override
  String get characterName => '이름';

  @override
  String get nameRule => '한 가지 문자로 16자 이내, 글자로 시작 (숫자와 _ 사용 가능)';

  @override
  String get race => '종족';

  @override
  String get characterClass => '직업';

  @override
  String get raceHuman => '인간';

  @override
  String get raceElf => '엘프';

  @override
  String get raceDwarf => '드워프';

  @override
  String get raceOrc => '오크';

  @override
  String get classWarrior => '전사';

  @override
  String get classMage => '마법사';

  @override
  String get classRogue => '도적';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Lv.$level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Lv.$level';
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
  String get xp => '경험치';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => '능력치';

  @override
  String get combat => '전투';

  @override
  String get statStr => '힘';

  @override
  String get statAgi => '민첩';

  @override
  String get statInt => '지능';

  @override
  String get statVit => '활력';

  @override
  String get statCha => '지혜';

  @override
  String get statLuck => '행운';

  @override
  String get attack => '공격력';

  @override
  String get defense => '방어력';

  @override
  String get crit => '치명타';

  @override
  String get hit => '명중';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => '골드';

  @override
  String get switchCharacter => '캐릭터 변경';

  @override
  String get statePendingDelete => '삭제 대기 중';

  @override
  String get stateFrozen => '동결됨';

  @override
  String get errInvalidCredentials => '이메일/사용자 이름 또는 비밀번호가 올바르지 않습니다.';

  @override
  String get errAccountLocked => '계정이 일시적으로 잠겼습니다.';

  @override
  String get errAccountBanned => '이 계정은 정지되었습니다.';

  @override
  String get errAccountExists => '이미 사용 중인 이메일 또는 사용자 이름입니다.';

  @override
  String get errRateLimited => '시도 횟수가 너무 많습니다. 잠시 기다려 주세요.';

  @override
  String get errNameTaken => '이미 사용 중인 이름입니다.';

  @override
  String get errSlotCap => '캐릭터 슬롯이 모두 사용 중입니다.';

  @override
  String get errNotPlayable => '지금은 이 캐릭터로 플레이할 수 없습니다.';

  @override
  String get errSessionExpired => '세션이 만료되었습니다. 다시 로그인하세요.';

  @override
  String errGeneric(String code) {
    return '문제가 발생했습니다 ($code).';
  }

  @override
  String get navCharacter => '캐릭터';

  @override
  String get navActivity => '채집';

  @override
  String get actMining => '채광';

  @override
  String get actFishing => '낚시';

  @override
  String get actWoodcutting => '벌목';

  @override
  String get actHerbalism => '약초 채집';

  @override
  String get actHunting => '사냥';

  @override
  String get activityChoose => '활동 선택';

  @override
  String get durationLabel => '시간';

  @override
  String get gatherResource => '자원';

  @override
  String gatherXp(int xp, int next) {
    return '경험치 $xp / $next';
  }

  @override
  String get gatherMax => '최대 레벨';

  @override
  String gatherLocked(int level) {
    return '레벨 $level에 해금';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name 경험치 +$xp';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '$name 레벨 $level 달성!';
  }

  @override
  String get errGatherLevel => '이 자원을 채집하기에는 채집 레벨이 낮습니다.';

  @override
  String startActivity(String name) {
    return '$name 시작';
  }

  @override
  String get collect => '수령';

  @override
  String get abandon => '포기';

  @override
  String get cancel => '취소';

  @override
  String get soFar => '지금까지 채집한 양';

  @override
  String get activityDone => '완료 — 수령 가능';

  @override
  String remaining(String time) {
    return '$time 남음';
  }

  @override
  String durHm(int h, String m) {
    return '$h시간 $m분';
  }

  @override
  String durM(int m) {
    return '$m분';
  }

  @override
  String durS(int s) {
    return '$s초';
  }

  @override
  String get collectEarlyTitle => '일찍 수령할까요?';

  @override
  String get collectEarlyBody => '지금까지 채집한 것을 받고 활동이 종료됩니다.';

  @override
  String get abandonTitle => '활동을 포기할까요?';

  @override
  String get abandonBody => '지금까지 채집한 것을 모두 잃습니다.';

  @override
  String get lootCollect => '좋아요!';

  @override
  String get lootNothing => '이번에는 아무것도 없습니다.';

  @override
  String xpGain(int xp) {
    return '경험치 +$xp';
  }

  @override
  String get rarityCommon => '일반';

  @override
  String get rarityUncommon => '고급';

  @override
  String get rarityRare => '희귀';

  @override
  String get rarityEpic => '영웅';

  @override
  String get rarityLegendary => '전설';

  @override
  String get rarityMythic => '신화';

  @override
  String get errSlotBusy => '다른 활동이 이미 진행 중입니다.';

  @override
  String get errNotActive => '이 활동은 이미 종료되었습니다.';

  @override
  String get evAchievement => '업적 달성!';

  @override
  String get evMasteryTier => '숙련 단계 상승!';

  @override
  String get evMarketSold => '등록한 물품이 판매되었습니다.';

  @override
  String get evMailReceived => '새 우편이 도착했습니다.';

  @override
  String get evPartyInvite => '파티에 초대되었습니다.';

  @override
  String get evGuildInvite => '길드에 초대되었습니다.';

  @override
  String get navBag => '가방';

  @override
  String get equipment => '장비';

  @override
  String get bag => '가방';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap kg';
  }

  @override
  String get slotHelm => '투구';

  @override
  String get slotNecklace => '목걸이';

  @override
  String get slotChest => '가슴';

  @override
  String get slotGloves => '장갑';

  @override
  String get slotLegs => '다리';

  @override
  String get slotBoots => '신발';

  @override
  String get slotMainhand => '주무기';

  @override
  String get slotOffhand => '보조 무기';

  @override
  String get slotRing => '반지';

  @override
  String get slotTrinket => '장신구';

  @override
  String get equip => '장착';

  @override
  String get unequip => '해제';

  @override
  String get sell => '판매';

  @override
  String get sortLabel => '정렬';

  @override
  String get sortSlot => '가방 순서';

  @override
  String get sortRarity => '희귀도';

  @override
  String get sortType => '종류';

  @override
  String get sortName => '이름';

  @override
  String get sortValue => '가치';

  @override
  String selectedCount(int count) {
    return '$count개 선택됨';
  }

  @override
  String get selectAll => '모두 선택';

  @override
  String get clearSelection => '선택 해제';

  @override
  String sellSelectedTitle(int count) {
    return '아이템 $count개를 판매할까요?';
  }

  @override
  String destroySelectedTitle(int count) {
    return '아이템 $count개를 파괴할까요?';
  }

  @override
  String get deleteEntry => '삭제';

  @override
  String craftHowMany(String name) {
    return '$name을(를) 몇 개 만들까요?';
  }

  @override
  String craftCount(int count) {
    return '$count개 제작';
  }

  @override
  String get destroy => '파괴';

  @override
  String sellTitle(String name) {
    return '$name을(를) 판매할까요?';
  }

  @override
  String get sellBody => '2시간 동안 상인에게서 다시 구매할 수 있습니다.';

  @override
  String destroyTitle(String name) {
    return '$name을(를) 파괴할까요?';
  }

  @override
  String get destroyBody => '되돌릴 수 없습니다.';

  @override
  String get soldFor => '판매 금액';

  @override
  String get quantity => '수량';

  @override
  String get requiresLevel => '요구 레벨';

  @override
  String get weight => '무게';

  @override
  String weightGrams(int g) {
    return '$g g';
  }

  @override
  String get sellValue => '판매가';

  @override
  String get bound => '귀속';

  @override
  String get itemLocked => '이 아이템은 거래, 우편 또는 등록에 묶여 있습니다.';

  @override
  String get errBagFull => '가방이 가득 찼습니다.';

  @override
  String get errItemLocked => '이 아이템은 잠겨 있습니다.';

  @override
  String get errClassRestriction => '현재 직업은 이 아이템을 사용할 수 없습니다.';

  @override
  String get errLevelRequirement => '레벨이 너무 낮습니다.';

  @override
  String get errNotEquippable => '이 아이템은 장착할 수 없습니다.';

  @override
  String get navAdventure => '모험';

  @override
  String get segHunt => '사냥';

  @override
  String get segDungeon => '던전';

  @override
  String get segTravel => '이동';

  @override
  String get huntChoose => '사냥감 선택';

  @override
  String get startHunt => '사냥 시작';

  @override
  String get huntDone => '사냥 종료 — 전리품 수령';

  @override
  String get huntEarlyBody => '지금까지 이긴 전투의 보상을 받고 사냥이 종료됩니다.';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'HP $hp · 공격 $atk · 방어 $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '$won승 · $lost패';
  }

  @override
  String levelUp(int level) {
    return '레벨 업! 이제 레벨 $level';
  }

  @override
  String levelReq(int level) {
    return 'Lv.$level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '전투 $count회 · 보스 $boss';
  }

  @override
  String get dungeonEnter => '입장';

  @override
  String get dungeonDone => '공략 완료 — 보상 수령';

  @override
  String get dungeonWaitHint => '공략이 끝나면 보상이 열립니다.';

  @override
  String get dungeonCleared => '클리어!';

  @override
  String get dungeonFailed => '파티 전멸';

  @override
  String get bossKilled => '보스 처치';

  @override
  String get bonus => '클리어 보너스';

  @override
  String get youAreIn => '현재 위치';

  @override
  String get routes => '경로';

  @override
  String get noRoutes => '여기서 갈 수 있는 경로가 없습니다.';

  @override
  String get travelTo => '이동';

  @override
  String travelling(String place) {
    return '$place(으)로 이동 중';
  }

  @override
  String get arrive => '도착';

  @override
  String get arrived => '도착했습니다';

  @override
  String get cancelTravel => '되돌아가기';

  @override
  String get kindCity => '도시';

  @override
  String get kindZone => '지역';

  @override
  String get kindWilderness => '야생';

  @override
  String get kindDungeon => '던전';

  @override
  String get kindContinent => '대륙';

  @override
  String get safeZone => '안전';

  @override
  String get errBusy => '이미 진행 중인 작업이 있습니다.';

  @override
  String get errNotReady => '아직 끝나지 않았습니다.';

  @override
  String get errWeeklyLocked => '이번 주에 이미 이 던전을 클리어했습니다.';

  @override
  String get errNoRoute => '그곳으로 가는 직행 경로가 없습니다.';

  @override
  String get navWorkshop => '공방';

  @override
  String get segCraft => '제작';

  @override
  String get segShop => '상점';

  @override
  String get profBlacksmithing => '대장 기술';

  @override
  String get profAlchemy => '연금술';

  @override
  String get profCooking => '요리';

  @override
  String get learn => '배우기';

  @override
  String recipeCount(int count) {
    return '레시피 $count개';
  }

  @override
  String get itemType => '종류';

  @override
  String get stackSize => '묶음';

  @override
  String get skillXpReward => '기술 경험치';

  @override
  String get learnHint => '이 직업을 배우면 레시피가 열립니다.';

  @override
  String get craft => '제작';

  @override
  String crafted(String name, int count) {
    return '$name $count개 제작 완료';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession 레벨 $level 달성';
  }

  @override
  String requiresSkill(int level) {
    return '기술 $level 필요';
  }

  @override
  String get buy => '구매';

  @override
  String bought(String name) {
    return '$name 구매 완료';
  }

  @override
  String get buybacks => '재구매';

  @override
  String get buybackHint => '판매한 아이템은 2시간 동안 다시 구매할 수 있습니다.';

  @override
  String get errInsufficientGold => '골드가 부족합니다.';

  @override
  String get errInsufficientMaterials => '재료가 부족합니다.';

  @override
  String get errSkillTooLow => '이 레시피에 필요한 기술 레벨이 부족합니다.';

  @override
  String get errProfessionNotLearned => '먼저 직업을 배우세요.';

  @override
  String get errPriceChanged => '가격이 변경되었습니다. 확인 후 다시 시도하세요.';

  @override
  String get errCannotSell => '이 아이템은 판매할 수 없습니다.';

  @override
  String get errBuybackExpired => '재구매 기간이 끝났습니다.';

  @override
  String get segMarket => '시장';

  @override
  String get classWeapon => '무기';

  @override
  String get classArmor => '방어구';

  @override
  String get classConsumable => '소모품';

  @override
  String get classMaterial => '재료';

  @override
  String get filterAll => '전체';

  @override
  String get sortCheapest => '최저가순';

  @override
  String get sortNewest => '최신순';

  @override
  String get marketEmpty => '조건에 맞는 물품이 없습니다.';

  @override
  String listingLine(String seller, String price) {
    return '$seller · 개당 $price';
  }

  @override
  String get loadMore => '더 보기';

  @override
  String myListings(int count) {
    return '내 등록 물품 ($count)';
  }

  @override
  String get cancelListing => '등록 취소';

  @override
  String cancelListingTitle(String name) {
    return '$name 등록을 취소할까요?';
  }

  @override
  String get cancelListingBody => '아이템은 가방으로 돌아옵니다. 등록 수수료는 처음 6시간 이내에만 일부 환불됩니다.';

  @override
  String buyListingTitle(String name) {
    return '$name을(를) 구매할까요?';
  }

  @override
  String buyListingBody(String price, String seller) {
    return '$seller에게 $price을(를) 지불합니다.';
  }

  @override
  String get listOnMarket => '시장에 등록';

  @override
  String listed(String name) {
    return '$name 시장에 등록됨';
  }

  @override
  String quantityOf(int count, int max) {
    return '수량: $max개 중 $count개';
  }

  @override
  String get unitPrice => '개당 가격';

  @override
  String suggestedPrice(String price) {
    return '시세: $price';
  }

  @override
  String get listingFee => '등록 수수료 (환불 불가)';

  @override
  String get youReceive => '판매 시 수령액';

  @override
  String get errPriceOutOfRange => '허용 범위를 벗어난 가격입니다.';

  @override
  String get errListingCap => '활성 등록 물품이 너무 많습니다.';

  @override
  String get errSelfBuy => '자신의 물품은 구매할 수 없습니다.';

  @override
  String get errListingUnavailable => '이 물품은 더 이상 구매할 수 없습니다.';

  @override
  String get errNotTradable => '이 아이템은 거래할 수 없습니다.';

  @override
  String get social => '소셜';

  @override
  String get tabChat => '채팅';

  @override
  String get tabFriends => '친구';

  @override
  String get tabMail => '우편';

  @override
  String get chanGlobal => '월드';

  @override
  String get chanGuild => '길드';

  @override
  String get chanParty => '파티';

  @override
  String get chatEmpty => '아직 메시지가 없습니다. 인사를 건네 보세요!';

  @override
  String chatHint(String channel) {
    return '$channel에 메시지 보내기';
  }

  @override
  String get send => '보내기';

  @override
  String get addFriend => '추가';

  @override
  String friendAdded(String name) {
    return '$name님과 친구가 되었습니다';
  }

  @override
  String friendRequested(String name) {
    return '$name님에게 친구 요청을 보냈습니다';
  }

  @override
  String get friendRequests => '친구 요청';

  @override
  String friendsCount(int count) {
    return '친구 ($count)';
  }

  @override
  String get noFriends => '아직 친구가 없습니다. 이름으로 추가해 보세요.';

  @override
  String get more => '더 보기';

  @override
  String get whisper => '귓속말';

  @override
  String get remove => '삭제';

  @override
  String get block => '차단';

  @override
  String get unblock => '차단 해제';

  @override
  String get blockedList => '차단 목록';

  @override
  String removeFriendTitle(String name) {
    return '$name님을 친구에서 삭제할까요?';
  }

  @override
  String get removeFriendBody => '나중에 다시 친구 요청을 보낼 수 있습니다.';

  @override
  String blockTitle(String name) {
    return '$name님을 차단할까요?';
  }

  @override
  String get blockBody => '상대는 귓속말을 보낼 수 없고, 채팅이 숨겨지며, 우편은 전달되지 않습니다. 친구 관계도 끝납니다.';

  @override
  String get accept => '수락';

  @override
  String get decline => '거절';

  @override
  String get newMail => '우편 쓰기';

  @override
  String get inboxEmpty => '우편함이 비어 있습니다.';

  @override
  String get noSubject => '(제목 없음)';

  @override
  String mailFrom(String name) {
    return '보낸 사람: $name';
  }

  @override
  String get attachments => '첨부물';

  @override
  String get claim => '첨부물 받기';

  @override
  String get mailClaimed => '첨부물을 가방과 지갑으로 옮겼습니다';

  @override
  String get deleteMailTitle => '이 우편을 삭제할까요?';

  @override
  String get deleteMailBody => '삭제한 우편은 복구할 수 없습니다.';

  @override
  String mailSent(String postage) {
    return '우편을 보냈습니다 (우편료 $postage)';
  }

  @override
  String get mailTo => '받는 사람 (캐릭터 이름)';

  @override
  String get mailSubject => '제목';

  @override
  String get mailBody => '내용';

  @override
  String get attachGold => '골드 첨부';

  @override
  String get errAlreadyFriends => '이미 친구입니다.';

  @override
  String get errAlreadyRequested => '친구 요청이 이미 대기 중입니다.';

  @override
  String get errBlocked => '이 플레이어와 상호작용할 수 없습니다.';

  @override
  String get errFriendCap => '친구 목록이 가득 찼습니다.';

  @override
  String get errTargetNotFound => '해당 이름의 캐릭터가 없습니다.';

  @override
  String get errSelfTarget => '자신을 대상으로 할 수 없습니다.';

  @override
  String get errInboxFull => '상대의 우편함이 가득 찼습니다.';

  @override
  String get errTooLong => '너무 깁니다.';

  @override
  String get errEmpty => '먼저 내용을 입력하세요.';

  @override
  String get errNotInChannel => '해당 그룹에 속해 있지 않습니다.';

  @override
  String get errMuted => '채팅 금지 상태라 지금은 채팅할 수 없습니다.';

  @override
  String get errHasAttachments => '삭제하기 전에 첨부물을 받으세요.';

  @override
  String get useItem => '사용';

  @override
  String restoresHp(int amount) {
    return 'HP $amount 회복';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return '$name 사용 · HP $hp/$max';
  }

  @override
  String get errHpFull => '체력이 이미 가득 찼습니다.';

  @override
  String get errInCombat => '먼저 사냥이나 던전 공략을 끝내세요.';

  @override
  String get errNotUsable => '이 아이템은 사용할 수 없습니다.';

  @override
  String get errTooWounded => '부상이 심해 싸울 수 없습니다. 쉬거나 물약을 마시세요.';

  @override
  String get forgotPassword => '비밀번호를 잊으셨나요?';

  @override
  String get resetHint => '계정 이메일을 입력하면 새 비밀번호 설정 링크를 보내 드립니다.';

  @override
  String get resetSent => '해당 이메일의 계정이 있다면 재설정 링크를 보냈습니다. 링크는 1시간 후 만료됩니다.';

  @override
  String get sendLink => '링크 보내기';

  @override
  String get close => '닫기';

  @override
  String get changePassword => '비밀번호 변경';

  @override
  String get currentPassword => '현재 비밀번호';

  @override
  String get newPassword => '새 비밀번호';

  @override
  String get save => '저장';

  @override
  String get passwordChanged => '비밀번호가 변경되었습니다. 다른 기기에서 로그아웃되었습니다.';

  @override
  String get verifyEmailBanner => '이메일 주소를 인증해 주세요.';

  @override
  String get resend => '다시 보내기';

  @override
  String get verificationSent => '인증 메일을 보냈습니다.';

  @override
  String get errEmailNotVerified => '캐릭터를 만들기 전에 이메일을 인증하세요.';

  @override
  String get errAlreadyVerified => '이메일이 이미 인증되었습니다.';

  @override
  String get tabParty => '파티';

  @override
  String get tabGuild => '길드';

  @override
  String get invite => '초대';

  @override
  String invitedBy(String name) {
    return '$name님의 초대';
  }

  @override
  String get noParty => '파티에 속해 있지 않습니다. 누군가를 초대하면 파티가 만들어집니다 (2~5명).';

  @override
  String get createParty => '파티 만들기';

  @override
  String get pendingInvites => '초대 목록';

  @override
  String partyMembers(int count, int cap) {
    return '파티 ($count/$cap)';
  }

  @override
  String get makeLeader => '파티장으로 지정';

  @override
  String get kick => '추방';

  @override
  String get leave => '나가기';

  @override
  String get disband => '해산';

  @override
  String get leavePartyTitle => '파티에서 나갈까요?';

  @override
  String get leavePartyBody => '파티장이 나가면 가장 오래된 파티원이 이어받습니다. 한 명만 남으면 파티가 해산됩니다.';

  @override
  String get disbandPartyTitle => '파티를 해산할까요?';

  @override
  String get disbandPartyBody => '모든 파티원이 파티에서 제외됩니다.';

  @override
  String get rankLeader => '길드장';

  @override
  String get rankOfficer => '간부';

  @override
  String get rankVeteran => '베테랑';

  @override
  String get rankMember => '길드원';

  @override
  String get rankInitiate => '수습';

  @override
  String get createGuild => '길드 창설';

  @override
  String createGuildBody(String fee) {
    return '길드 창설 비용은 $fee입니다.';
  }

  @override
  String createGuildFee(String fee) {
    return '길드 창설 ($fee)';
  }

  @override
  String get guildName => '길드 이름';

  @override
  String get guildTag => '태그';

  @override
  String get guildTagHint => '문자 또는 숫자 2~5자';

  @override
  String guildMembers(int count, int cap) {
    return '길드원 $count/$cap';
  }

  @override
  String disbandScheduled(String time) {
    return '$time 후 길드 해산';
  }

  @override
  String get motd => '오늘의 메시지';

  @override
  String get editMotd => '오늘의 메시지 수정';

  @override
  String get setRank => '등급 변경';

  @override
  String setRankFor(String name) {
    return '$name의 새 등급';
  }

  @override
  String transferTitle(String name) {
    return '$name님을 길드장으로 지정할까요?';
  }

  @override
  String get transferBody => '본인은 간부가 되며, 직접 되돌릴 수 없습니다.';

  @override
  String kickTitle(String name) {
    return '$name님을 길드에서 추방할까요?';
  }

  @override
  String get leaveGuildTitle => '길드를 나갈까요?';

  @override
  String get leaveGuildBody => '다시 가입하려면 새 초대가 필요합니다.';

  @override
  String get disbandGuildTitle => '길드를 해산할까요?';

  @override
  String get disbandGuildBody => '길드는 48시간 후 해산되며, 그 전까지 취소할 수 있습니다.';

  @override
  String get cancelDisband => '해산 취소';

  @override
  String get errAlreadyInGroup => '이미 그룹에 속해 있습니다.';

  @override
  String get errNotLeader => '파티장/길드장만 할 수 있습니다.';

  @override
  String get errGroupFull => '그룹이 가득 찼습니다.';

  @override
  String get errInviteInvalid => '초대가 만료되었습니다.';

  @override
  String get errAlreadyInvited => '초대가 이미 대기 중입니다.';

  @override
  String get errGuildForbidden => '현재 등급으로는 할 수 없습니다.';

  @override
  String get errLeaderCannotLeave => '먼저 길드장을 넘기거나 길드를 해산하세요.';

  @override
  String get errInvalidGuildName => '문자 또는 숫자 3~24자.';

  @override
  String get errInvalidTag => '태그는 문자 또는 숫자 2~5자여야 합니다.';

  @override
  String get errTagTaken => '이미 사용 중인 태그입니다.';

  @override
  String get progress => '업적 및 숙련';

  @override
  String get tabAchievements => '업적';

  @override
  String get tabMastery => '숙련';

  @override
  String get tabReputation => '평판';

  @override
  String get tabTitles => '칭호';

  @override
  String achievementPoints(int points) {
    return '업적 점수 $points';
  }

  @override
  String pointsShort(int points) {
    return '$points점';
  }

  @override
  String rewardTitle(String title) {
    return '칭호: $title';
  }

  @override
  String masteryTier(int tier, int max) {
    return '$tier/$max단계';
  }

  @override
  String repDailyLeft(int amount) {
    return '오늘 +$amount';
  }

  @override
  String get noTitle => '칭호 없음';

  @override
  String get noTitlesYet => '칭호는 업적으로 해금됩니다.';

  @override
  String get errTitleNotOwned => '아직 해금하지 않은 칭호입니다.';

  @override
  String get notifications => '알림';

  @override
  String get markAllRead => '모두 읽음으로 표시';

  @override
  String get inboxNothing => '아직 알림이 없습니다.';

  @override
  String get notifSettings => '알림 설정';

  @override
  String get muteAll => '푸시 알림 일시 중지';

  @override
  String get muteHint => '24시간 동안. 받은 알림함은 계속 쌓입니다.';

  @override
  String mutedFor(String time) {
    return '$time 동안 일시 중지';
  }

  @override
  String get modePushInbox => '푸시 + 알림함';

  @override
  String get modeInbox => '알림함만';

  @override
  String get modeSilent => '무음';

  @override
  String get modeOff => '끄기';

  @override
  String get nkActivity => '활동 완료';

  @override
  String get nkDungeon => '던전 공략 완료';

  @override
  String get nkTravel => '도착';

  @override
  String get nkWorldEvent => '월드 이벤트 시작';

  @override
  String get nkMarketSold => '아이템 판매';

  @override
  String get nkMail => '새 우편';

  @override
  String get nkPartyInvite => '파티 초대';

  @override
  String get nkGuildInvite => '길드 초대';

  @override
  String get nkAchievement => '업적 달성';

  @override
  String get nkSystem => '공지';

  @override
  String timeAgo(String time) {
    return '$time 전';
  }

  @override
  String get segEvents => '이벤트';

  @override
  String get eventsEmpty => '예정된 월드 이벤트가 없습니다.';

  @override
  String get worldBoss => '월드 보스';

  @override
  String get invasion => '침공';

  @override
  String startsIn(String time) {
    return '$time 후 시작';
  }

  @override
  String endsIn(String time) {
    return '$time 후 종료';
  }

  @override
  String regionOnly(String region) {
    return '$region 전용';
  }

  @override
  String get eventGoal => '목표';

  @override
  String get enlist => '참가';

  @override
  String get enlistHint => '참가하는 순간 현재 전투력이 고정됩니다. 이벤트가 진행되는 매 분마다 점수를 얻고, 보상은 종료 시 지급됩니다.';

  @override
  String enlistedPower(int power) {
    return '참가 중 · 전투력 $power';
  }

  @override
  String myContribution(int score, int tier) {
    return '내 점수 $score · 보상 등급 $tier';
  }

  @override
  String get eventSucceeded => '승리';

  @override
  String get eventFailed => '실패';

  @override
  String tokenBalance(int tokens) {
    return '이벤트 토큰 $tokens개';
  }

  @override
  String claimDaily(int tokens) {
    return '$tokens개 받기';
  }

  @override
  String get claimedToday => '수령 완료';

  @override
  String eventTokensGained(int tokens) {
    return '이벤트 토큰 +$tokens';
  }

  @override
  String get errEventNotOpen => '이 이벤트의 참가가 마감되었습니다.';

  @override
  String get errWrongRegion => '먼저 이벤트 지역으로 이동하세요.';

  @override
  String get errAlreadyEnlisted => '이미 참가했습니다.';

  @override
  String get errAlreadyClaimed => '오늘은 이미 받았습니다.';

  @override
  String get errEventNotActive => '이 이벤트는 진행 중이 아닙니다.';

  @override
  String gatherYield(String activity) {
    return '$activity 산출량';
  }

  @override
  String perHour(String value) {
    return '시간당 $value';
  }

  @override
  String get quests => '퀘스트';

  @override
  String get questsInProgress => '진행 중';

  @override
  String get questsNoneOpen => '진행 중인 퀘스트가 없습니다. 아래에서 고르세요.';

  @override
  String get questsStory => '스토리';

  @override
  String get questsSide => '서브 퀘스트';

  @override
  String get questsDaily => '일일';

  @override
  String resetsIn(String time) {
    return '$time 후 초기화';
  }

  @override
  String get turnIn => '완료 보고';

  @override
  String questRewarded(String rewards) {
    return '퀘스트 완료: $rewards';
  }

  @override
  String get abandonQuestTitle => '퀘스트를 포기할까요?';

  @override
  String get abandonQuestBody => '진행 상황을 잃습니다. 나중에 다시 받을 수 있습니다.';

  @override
  String get abandonDailyBody => '진행 상황을 잃으며 오늘은 이 일일 퀘스트를 다시 받을 수 없습니다.';

  @override
  String objGather(String activity) {
    return '채집 ($activity)';
  }

  @override
  String get objKillAny => '적 처치';

  @override
  String objKill(String mob) {
    return '$mob 처치';
  }

  @override
  String get objBoss => '보스 처치';

  @override
  String get objDungeonAny => '던전 클리어';

  @override
  String objDungeon(String dungeon) {
    return '$dungeon 클리어';
  }

  @override
  String get objCraftAny => '아이템 제작';

  @override
  String objCraft(String item) {
    return '$item 제작';
  }

  @override
  String get evQuestComplete => '퀘스트 목표 달성 — 완료 보고하세요!';

  @override
  String get errQuestLocked => '먼저 이전 스토리 퀘스트를 끝내세요.';

  @override
  String get errQuestTaken => '이미 받은 퀘스트입니다.';

  @override
  String get errQuestNotOffered => '오늘의 게시판에 없는 일일 퀘스트입니다.';

  @override
  String get errQuestCap => '진행 중인 퀘스트가 너무 많습니다. 하나를 끝내거나 포기하세요.';

  @override
  String get errQuestNotActive => '진행 중인 퀘스트가 아닙니다.';

  @override
  String get errQuestIncomplete => '아직 목표를 달성하지 못했습니다.';

  @override
  String get skills => '스킬';

  @override
  String get skillActive => '액티브';

  @override
  String get skillPassive => '패시브';

  @override
  String skillStrike(int power) {
    return '피해 $power%';
  }

  @override
  String skillHeal(int power) {
    return 'HP가 절반 이하일 때 HP $power% 회복';
  }

  @override
  String skillMp(int mp) {
    return 'MP $mp';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: '$turns턴마다', one: '매 턴');
    return '$_temp0';
  }

  @override
  String get termsOfService => '서비스 약관';

  @override
  String get privacyPolicy => '개인정보 처리방침';

  @override
  String get accountMenu => '계정';

  @override
  String get downloadMyData => '내 데이터 다운로드';

  @override
  String get exportOpened => '브라우저에서 다운로드 링크를 열었습니다 (15분간 유효).';

  @override
  String get linkOpenFailed => '링크를 열 수 없습니다.';

  @override
  String get deleteAccount => '계정 삭제';

  @override
  String get deleteAccountBody => '계정과 모든 캐릭터가 30일 후 영구 삭제됩니다. 그 기간 안에 로그인하면 삭제가 취소됩니다. 확인을 위해 비밀번호를 입력하세요.';

  @override
  String get deleteAccountConfirm => '계정 삭제';

  @override
  String deleteScheduled(String date) {
    return '계정이 $date에 삭제됩니다. 취소하려면 그 전에 로그인하세요.';
  }

  @override
  String get deletionCancelled => '계정 삭제 요청이 취소되었습니다.';

  @override
  String get support => '고객 지원';

  @override
  String get supportIntro => '버그를 발견했거나 계정 관련 도움이 필요하신가요? 문의를 남기면 팀이 이곳에 답변하고 알림을 보내 드립니다.';

  @override
  String get supportNew => '새 문의';

  @override
  String get supportNone => '아직 문의가 없습니다.';

  @override
  String get supportSent => '문의가 접수되었습니다. 답변이 등록되면 알려 드립니다.';

  @override
  String get supportCategory => '주제';

  @override
  String get supportSubject => '제목';

  @override
  String get supportMessage => '무슨 일이 있었나요? 가능하면 단계별로 적어 주세요.';

  @override
  String get supportTooShort => '제목은 3자 이상, 내용은 10자 이상이어야 합니다.';

  @override
  String get supportBug => '버그';

  @override
  String get supportAccount => '계정';

  @override
  String get supportAbuse => '악용 / 플레이어 신고';

  @override
  String get supportFeedback => '건의';

  @override
  String get supportOther => '기타';

  @override
  String get supportOpen => '접수됨';

  @override
  String get supportAnswered => '답변됨';

  @override
  String get supportClosed => '종료됨';

  @override
  String get supportReply => '운영팀 답변';

  @override
  String get nkSupport => '고객 지원 답변';

  @override
  String get errServerBusy => '서버가 혼잡합니다. 몇 초 후 다시 시도하세요.';

  @override
  String get notifTabInbox => '알림';

  @override
  String get notifTabRewards => '보상';

  @override
  String get rewardsNone => '아직 보상이 없습니다. 채집하거나 전투에서 얻은 모든 것이 여기에 표시됩니다 (30일 보관).';

  @override
  String rewardFromActivity(String name) {
    return '채집: $name';
  }

  @override
  String rewardFromHunt(String name) {
    return '사냥: $name';
  }

  @override
  String rewardFromDungeon(String name) {
    return '던전: $name';
  }

  @override
  String rewardLevelUp(int level) {
    return '레벨 $level!';
  }

  @override
  String get levelMax => 'MAX';

  @override
  String get promotionAvailable => '전직 가능';

  @override
  String get promotionChoose => '선택';

  @override
  String promotionConfirm(String name) {
    return '$name(으)로 전직할까요?';
  }

  @override
  String get promotionBody => '이 선택은 직업을 영구히 결정합니다.';

  @override
  String nextMilestone(int level, String what) {
    return '다음: Lv.$level · $what';
  }

  @override
  String get msPromotion => '직업 전직';

  @override
  String msAwakening(String name) {
    return '각성: $name';
  }

  @override
  String msUltimate(String name) {
    return '궁극기: $name';
  }

  @override
  String get msCapstone => '특성 최종 단계';

  @override
  String get msMastery => '직업 마스터리';

  @override
  String msTitle(String name) {
    return '칭호: $name';
  }

  @override
  String estMinutes(int n) {
    return '약 $n분';
  }

  @override
  String estHours(String h) {
    return '약 $h시간';
  }

  @override
  String overflowHeld(int xp) {
    return '경험치 $xp이(가) 선택을 기다립니다';
  }

  @override
  String attributePointsBtn(int n) {
    return '능력치 · $n포인트';
  }

  @override
  String talentPointsBtn(int n) {
    return '특성 · $n포인트';
  }

  @override
  String get talentsTitle => '특성';

  @override
  String capstoneLocked(int level, int points) {
    return '레벨 $level과 이 트리의 $points포인트 필요';
  }

  @override
  String masteryUnlocked(String title) {
    return '직업 마스터리: $title';
  }

  @override
  String get errPromotionUnavailable => '해당 전직은 불가능합니다.';

  @override
  String get errNotEnoughPoints => '포인트가 부족합니다.';

  @override
  String get errTalentMaxed => '이 특성은 이미 최대 등급입니다.';

  @override
  String get errCapstoneLocked => '최종 단계에는 레벨과 완성된 트리가 필요합니다.';

  @override
  String get errSpecRequirement => '현재 전문화는 이 아이템을 사용할 수 없습니다.';

  @override
  String get errAttributeRequirement => '이 아이템을 쓰기에는 능력치가 낮습니다.';

  @override
  String get errBusyCombat => '전투, 방치 활동, 던전 중에는 할 수 없습니다.';

  @override
  String skillGuard(int pct, int turns) {
    return '$turns턴 동안 방어구 +$pct%';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '$turns턴 동안 공격력 +$pct%';
  }

  @override
  String get classRanger => '레인저';

  @override
  String get classCleric => '성직자';

  @override
  String get classMonk => '수도승';

  @override
  String get classPaladin => '팔라딘';

  @override
  String get classDruid => '드루이드';

  @override
  String get classBard => '음유시인';

  @override
  String get classShaman => '주술사';

  @override
  String get errEquipmentPermission => '현재 직업은 이 장비 종류를 사용할 수 없습니다.';

  @override
  String get classResource => '자원';

  @override
  String get classFocus => '주요 능력치';

  @override
  String classSecondary(String names) {
    return '보조: $names';
  }

  @override
  String get resStartsEmpty => '0에서 시작';

  @override
  String get resStartsFull => '가득 찬 상태로 시작';

  @override
  String resStartsAt(int pct) {
    return '$pct%에서 시작';
  }

  @override
  String resPerHit(int n) {
    return '타격당 +$n';
  }

  @override
  String resPerTurn(int n) {
    return '턴당 +$n';
  }

  @override
  String baseAttributesHint(int total) {
    return '기본 능력치 · 2레벨마다 1포인트, 총 $total';
  }

  @override
  String get classEquipment => '장비';

  @override
  String get firstSkills => '초기 스킬';

  @override
  String get classPath => '직업 경로';

  @override
  String get classPathHint => '직업을 누르면 스킬을 볼 수 있습니다.';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Lv.$awakening 각성 · Lv.$ultimate 궁극기 · Lv.$capstone 특성 최종 단계 · Lv.$mastery 직업 마스터리';
  }

  @override
  String get skillTagActive => '액티브';

  @override
  String get skillTagPassive => '패시브';

  @override
  String get skillTagAwakening => '각성';

  @override
  String get skillTagUltimate => '궁극기';

  @override
  String skillCostCd(int cost, int cd) {
    return '비용 $cost · 재사용 $cd턴';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Lv.$level 마스터리: $title';
  }

  @override
  String get classInfoFailed => '직업 정보를 불러오지 못했습니다.';

  @override
  String xpShare(int pct) {
    return '경험치 $pct%';
  }

  @override
  String get dungeonsSolo => '던전';

  @override
  String get raidsTab => '레이드';

  @override
  String raidGroup(int min, int max) {
    return '$min~$max인 파티';
  }

  @override
  String raidBossGear(int count) {
    return '보스: 장비 $count개 확정';
  }

  @override
  String get raidOpen => '레이드 열기';

  @override
  String get raidNeedParty => '레이드는 파티 콘텐츠입니다. 아래에서 무작위 파티 대기열에 들어가거나, 게시판의 모집에 참여하거나, 직접 모집하세요.';

  @override
  String get raidLeaderOnly => '레이드는 파티장이 엽니다.';

  @override
  String get raidLobby => '레이드 로비';

  @override
  String get raidJoin => '참가';

  @override
  String get raidLeave => '나가기';

  @override
  String get raidCancel => '취소';

  @override
  String get raidCancelBody => '모두의 로비가 닫힙니다.';

  @override
  String get raidStart => '레이드 시작';

  @override
  String get raidStartBody => '시작하면 누구도 참가하거나 나갈 수 없습니다.';

  @override
  String raidMembers(int count, int max) {
    return '인원 ($count/$max)';
  }

  @override
  String raidNeedMore(int min) {
    return '시작하려면 최소 $min명 필요';
  }

  @override
  String get raidWaitHint => '레이드가 끝나면 보상이 열리며, 각자 자신의 보상을 받습니다.';

  @override
  String get raidDone => '레이드 종료 — 보상 수령';

  @override
  String get raidLocked => '이번 주 클리어';

  @override
  String get raidCollected => '수령 완료';

  @override
  String get errRaidNeedsParty => '먼저 파티에 가입하세요.';

  @override
  String get errRaidExists => '파티에 이미 열린 레이드가 있습니다.';

  @override
  String get errAlreadyInRaid => '이미 레이드에 참가 중입니다.';

  @override
  String get errRaidNotOpen => '이 레이드는 더 이상 열려 있지 않습니다.';

  @override
  String get errRaidNotEnough => '시작할 인원이 부족합니다.';

  @override
  String get errRaidWeeklyLocked => '이번 주에 이미 이 레이드를 클리어했습니다.';

  @override
  String errCharacterBusy(String action) {
    return '먼저 캐릭터의 현재 작업을 끝내세요: $action.';
  }

  @override
  String get busyGathering => '채집';

  @override
  String get busyHunt => '사냥';

  @override
  String get busyDungeon => '던전';

  @override
  String get busyRaid => '레이드';

  @override
  String get busyTravel => '이동';

  @override
  String get busyCrafting => '제작';

  @override
  String get errCraftActive => '제작 작업이 이미 진행 중입니다.';

  @override
  String get errCraftTooLong => '이 작업은 6시간을 넘습니다.';

  @override
  String craftEach(String time) {
    return '개당 $time';
  }

  @override
  String craftTotal(String time) {
    return '총 $time';
  }

  @override
  String craftRunning(String item, int count) {
    return '$item $count개 제작 중';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count 완료';
  }

  @override
  String get craftDone => '모두 완료 — 수령';

  @override
  String get craftCollectEarlyTitle => '지금 수령할까요?';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return '완성된 $done개는 지급되고, 나머지 $rest개의 재료와 골드는 돌려받습니다.';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '$made개 제작, $refunded개 환불';
  }

  @override
  String get craftBusyHint => '제작에는 시간이 걸리며, 그동안 캐릭터는 다른 일을 하지 않습니다.';

  @override
  String get eventHowTitle => '참가 방법';

  @override
  String get eventRewardsTitle => '보상';

  @override
  String eventStepEnlist(int level) {
    return '시작 24시간 전부터 종료 시까지 참가할 수 있습니다 (Lv.$level 이상).';
  }

  @override
  String eventStepRegion(String region) {
    return '참가할 때 $region에 있어야 합니다.';
  }

  @override
  String eventStepPower(int power) {
    return '참가 시 전투력이 고정됩니다: 공격력 × 2 + 방어력 + HP ÷ 10. 현재 전투력: $power.';
  }

  @override
  String eventStepScore(int tier) {
    return '점수 = 전투력 × 참가 시간(분). 시작 시점(또는 더 늦게 참가한 경우 그 시점)부터 계산하며 $tier단계에서 멈춥니다.';
  }

  @override
  String get eventStepGoal => '모두 함께 목표를 달성하면 최고 단계까지 보상을 받고, 아니면 2단계까지만 받습니다.';

  @override
  String get eventStepPaid => '보상은 종료 시 자동 지급됩니다. 참가해도 다른 행동은 멈추지 않습니다.';

  @override
  String eventTier(int tier, String score) {
    return '$tier단계 · $score점 이상';
  }

  @override
  String eventTierTokens(int count) {
    return '이벤트 토큰 +$count';
  }

  @override
  String eventEstimate(String score, int tier) {
    return '지금 참가 시: 약 $score점 → $tier단계';
  }

  @override
  String eventEstimateNone(String score) {
    return '지금 참가 시: 약 $score점 — 단계에 못 미침';
  }

  @override
  String seasonalHow(int count) {
    return '이벤트 기간 동안 하루 한 번 이벤트 토큰 $count개를 받을 수 있습니다 (하루는 00:00 UTC에 시작).';
  }

  @override
  String seasonalBenefit(int count) {
    return '토큰은 계정에 남습니다 (잔액: $count).';
  }

  @override
  String get eventDetails => '자세히';

  @override
  String get createCharacterTitle => '캐릭터 만들기';

  @override
  String get classTraits => '직업 특성';

  @override
  String classResourceLine(String name) {
    return '자원: $name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return '주 능력치: $primary · 보조: $secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return '주 능력치: $primary';
  }

  @override
  String get classDifficulty => '난이도';

  @override
  String get difficultyEasy => '쉬움 — 실수에 관대하고 방치 사냥에 안전함';

  @override
  String get difficultyMedium => '보통 — 균형 잡힘, 스킬 순서가 중요함';

  @override
  String get difficultyHard => '어려움 — 약함, 알맞은 장비가 필요함';

  @override
  String get skillsTapHint => '눌러서 자세히 보기';

  @override
  String get chooseRace => '종족 선택';

  @override
  String continueAs(String name) {
    return '$name(으)로 계속';
  }

  @override
  String get nameYourCharacter => '캐릭터 이름 짓기';

  @override
  String get raceTroll => '트롤';

  @override
  String get raceUndead => '언데드';

  @override
  String get raceDarkElf => '다크 엘프';

  @override
  String get raceDrakari => '드라카리';

  @override
  String get raceTauren => '타우렌';

  @override
  String get raceGnome => '노움';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '$minutes분 동안 공격력 +$atk%, 방어력 +$def%';
  }

  @override
  String fedToast(int atk, int def) {
    return '포만: 공격력 +$atk%, 방어력 +$def%';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return '포만 · 공격력 +$atk% · 방어력 +$def% · $time 남음';
  }

  @override
  String huntReach(String regions) {
    return '이곳에서 사냥 가능: $regions';
  }

  @override
  String get huntNoMobsHere => '이곳에는 사냥감이 없습니다. 먼저 사냥터로 이동하세요.';

  @override
  String get huntOtherRegions => '다른 지역의 몬스터';

  @override
  String get huntTravelHint => '몬스터를 사냥하려면 그 지역으로 이동하세요.';

  @override
  String get errMobRegion => '이 몬스터는 다른 지역에 삽니다. 먼저 그곳으로 이동하세요.';

  @override
  String restoresHpPct(int pct) {
    return '최대 HP의 $pct% 회복';
  }

  @override
  String get statEvasion => '% 회피';

  @override
  String get statBlock => '% 방어';

  @override
  String get statCritDmg => '% 치명타 피해';

  @override
  String get statLifesteal => '% 흡혈';

  @override
  String get statDmgReduction => '% 피해 감소';

  @override
  String get classJewelry => '장신구';

  @override
  String get mapHint => '레벨이 허용하는 모든 지역으로 바로 이동할 수 있으며, 시간은 최단 도로를 따릅니다. 지도를 드래그하거나 확대하고 지역을 누르세요.';

  @override
  String get mapNoRoad => '이곳으로 가는 직행 도로가 없습니다. 옆의 빛나는 지역을 거쳐 이동하세요.';

  @override
  String get youAreHere => '현재 위치입니다.';

  @override
  String travelTime(String time) {
    return '이동 시간: $time';
  }

  @override
  String xpPerKill(int xp) {
    return '경험치 $xp';
  }

  @override
  String get deleteCharacter => '캐릭터 삭제';

  @override
  String deleteCharacterBody(String name) {
    return '$name은(는) 모든 아이템, 골드와 함께 즉시 삭제됩니다. 되돌릴 수 없습니다. 확인을 위해 캐릭터 이름을 입력하세요.';
  }

  @override
  String get characterNameConfirm => '캐릭터 이름';

  @override
  String characterDeletesAt(String when) {
    return '삭제 예정: $when';
  }

  @override
  String get cancelCharacterDelete => '삭제 취소';

  @override
  String characterDeleted(String name) {
    return '$name 삭제 완료.';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name은(는) 삭제되지 않습니다.';
  }

  @override
  String get errNameMismatch => '이름이 일치하지 않습니다.';

  @override
  String attrGainPerPoint(String gains) {
    return '포인트당: $gains';
  }

  @override
  String attrGainPending(String gains) {
    return '추가 예정: $gains';
  }

  @override
  String get attrNoEffect => '현재 직업에는 효과 없음';

  @override
  String get attrInputHint => '각 능력치에 줄 포인트를 입력하세요';

  @override
  String get gainHp => 'HP';

  @override
  String get gainMp => 'MP';

  @override
  String get gainHpRegen => 'HP/초';

  @override
  String get huntSupplies => '가져갈 보급품';

  @override
  String get huntSuppliesHint => '물약은 전투 사이 HP가 절반 아래일 때, 음식은 효과가 끝났을 때, 마나 물약은 마나가 부족할 때 사용합니다. 남은 것은 가방으로 돌아옵니다.';

  @override
  String get supplyHp => '회복 물약';

  @override
  String get supplyMp => '마나 물약';

  @override
  String get supplyMpHint => '마나, 기력, 집중, 자연 에너지, 영혼을 쓰는 직업용';

  @override
  String get supplyFood => '음식';

  @override
  String get huntNoSupplies => '가방에 물약이나 음식이 없습니다. 상점에서 사거나 제작하세요.';

  @override
  String restoresMpPct(int pct) {
    return '마나 $pct% 회복';
  }

  @override
  String get huntSupplyOnly => '사냥 보급품으로만 사용';

  @override
  String get errBadSupply => '이 아이템은 사냥에 가져갈 수 없습니다.';

  @override
  String get errNotEnoughSupplies => '가방에 그만큼의 보급품이 없습니다.';

  @override
  String get errHuntOnly => '보급품으로 사냥에 가져가세요.';

  @override
  String get huntFled => 'HP가 5% 아래로 떨어져 전투에서 물러났습니다';

  @override
  String suppliesUsed(String list) {
    return '사용: $list';
  }

  @override
  String get compareTitle => '착용 장비와 비교';

  @override
  String compareWorn(String name) {
    return '착용 중: $name';
  }

  @override
  String get compareEmpty => '이 슬롯은 비어 있습니다. 장착하면 아래 수치가 모두 더해집니다.';

  @override
  String get compareColStat => '능력치';

  @override
  String get compareColWorn => '착용';

  @override
  String get compareColThis => '이 장비';

  @override
  String get compareColDiff => '변화';

  @override
  String get compareSame => '착용 장비와 능력치가 같습니다.';

  @override
  String get wearOk => '현재 직업이 착용할 수 있습니다';

  @override
  String get wearClass => '현재 직업은 착용할 수 없습니다';

  @override
  String wearLevel(int level) {
    return '레벨 $level 필요';
  }

  @override
  String get wearAttribute => '능력치가 낮습니다';

  @override
  String get wearLegendClass => '직업 불일치';

  @override
  String get wearLegendLater => '아직 불가 (레벨 또는 능력치)';

  @override
  String get enchantAction => '마법 부여';

  @override
  String get enchantPick => '마법을 부여할 장비 선택';

  @override
  String get enchantNone => '이 두루마리에 맞는 장비가 없습니다.';

  @override
  String enchantFits(String groups) {
    return '대상: $groups';
  }

  @override
  String get enchantGroupWeapon => '무기';

  @override
  String get enchantGroupArmor => '방어구';

  @override
  String get enchantGroupJewelry => '장신구';

  @override
  String get enchantGroupAny => '모든 장비';

  @override
  String get enchantReplaceTitle => '마법 부여를 교체할까요?';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old은(는) 사라지고 $next이(가) 적용됩니다.';
  }

  @override
  String enchantDone(String item) {
    return '$item에 마법 부여 완료';
  }

  @override
  String enchantOn(String name) {
    return '마법 부여: $name';
  }

  @override
  String get wornTag => '착용';

  @override
  String get errNotEnchant => '마법 부여 두루마리가 아닙니다.';

  @override
  String get errEnchantSlot => '이 장비에는 맞지 않는 마법 부여입니다.';

  @override
  String get statGatherMining => '채광 산출량';

  @override
  String get statGatherHerbalism => '약초 채집 산출량';

  @override
  String get statGatherWoodcutting => '벌목 산출량';

  @override
  String get statGatherFishing => '낚시 산출량';

  @override
  String get statGatherHunting => '사냥 산출량';

  @override
  String get statCraftSpeed => '제작 속도';

  @override
  String diffXp(int pct) {
    return '경험치 $pct%';
  }

  @override
  String get diffNoXp => '경험치 없음';

  @override
  String get diffLegend => '레시피 색상: 현재 기술에서 제작으로 얻는 경험치';

  @override
  String get recipeGroupEnchantWeapon => '무기 마법 부여';

  @override
  String get recipeGroupEnchantArmor => '방어구 마법 부여';

  @override
  String get recipeGroupEnchantJewelry => '장신구 마법 부여';

  @override
  String get recipeGroupEnchantProfession => '직업 마법 부여';

  @override
  String get partyHint => '같은 작업(같은 채집, 직업, 사냥, 던전)을 동시에 하는 파티원은 일하는 인원당 경험치 +3%를 얻습니다. 5명이면 +15%. 누구의 경험치도 줄어들지 않습니다.';

  @override
  String get partyChat => '파티 채팅';

  @override
  String partyShared(String names, int kept, int earned) {
    return '$names님과 파티 보너스: $earned → 경험치 $kept';
  }

  @override
  String evRaidMatched(String raid) {
    return '$raid 파티를 찾았습니다. 레이드를 확인하세요.';
  }

  @override
  String get nkRaidMatched => '레이드 파티 찾음';

  @override
  String get workHunt => '사냥';

  @override
  String get workDungeon => '던전';

  @override
  String get raidFind => '파티 찾기';

  @override
  String get raidQueueJoin => '대기열 참가';

  @override
  String get raidQueueLeave => '대기열 나가기';

  @override
  String raidQueued(String raid, int waiting) {
    return '$raid 대기 중 · $waiting명 대기';
  }

  @override
  String get raidQueueHint => '5명이 모이거나, 최소 2명으로 5분이 지나면 무작위 파티가 만들어집니다.';

  @override
  String get raidAdvertise => '모집 등록';

  @override
  String get raidUnlist => '모집 삭제';

  @override
  String raidListed(String raid) {
    return '파티가 $raid 게시판에 등록되었습니다';
  }

  @override
  String get raidBoard => '모집 게시판';

  @override
  String get raidBoardEmpty => '열린 모집이 없습니다.';

  @override
  String get raidJoinParty => '파티 참가';

  @override
  String get raidPickRaid => '레이드';

  @override
  String get errListingGone => '모집이 사라졌거나 가득 찼습니다.';

  @override
  String get workGathering => '채집';

  @override
  String get workAdventure => '모험';

  @override
  String get workRaid => '레이드';

  @override
  String get workTravel => '이동 중';

  @override
  String get workCrafting => '제작';

  @override
  String get workIdle => '대기';

  @override
  String get errDungeonRegion => '먼저 지도에서 이곳으로 이동하세요.';

  @override
  String get dungeonsElsewhere => '이동할수록 던전과 레이드가 더 열립니다. 각각 지도의 해당 지역에서 입장합니다.';

  @override
  String get dungeonNoneHere => '여기서 갈 수 있는 던전이 없습니다. 지도에서 던전 지역으로 이동하세요.';

  @override
  String dungeonXp(int xp) {
    return '나에게 약 $xp 경험치';
  }

  @override
  String appVersionLine(String version) {
    return '버전 $version';
  }

  @override
  String get updateButton => '업데이트';

  @override
  String updateAvailable(String version) {
    return '새 버전 준비됨: $version';
  }

  @override
  String get updateLatest => '최신 버전입니다';

  @override
  String updateConfirm(String version, int size) {
    return '버전 $version($size MB)을 다운로드하여 설치할까요?';
  }

  @override
  String get updateDownloading => '다운로드 중…';

  @override
  String get updatePermission => 'MeNa World의 설치를 허용한 뒤 업데이트를 다시 누르세요.';

  @override
  String get updateFailed => '업데이트를 받지 못했습니다. 연결을 확인하세요.';

  @override
  String get fightReport => '전투 보고서';

  @override
  String get endTimeUp => '시간 종료: 예정된 끝까지 진행되었습니다.';

  @override
  String get endEarly => '일찍 수령: 예정된 끝보다 먼저 멈췄습니다.';

  @override
  String endFled(String mob) {
    return '$mob과(와)의 전투에서 HP가 5% 아래로 떨어져 물러났습니다.';
  }

  @override
  String endDefeat(String mob) {
    return '$mob에게 쓰러졌습니다.';
  }

  @override
  String endPartyDefeat(String mob) {
    return '파티가 $mob에게 쓰러졌습니다.';
  }

  @override
  String endTimeout(String mob) {
    return '$mob과(와)의 전투가 턴 제한을 넘어 중단되었습니다.';
  }

  @override
  String get endCleared => '클리어: 모든 웨이브와 보스를 쓰러뜨렸습니다.';

  @override
  String get repCampHint => 'HP가 50% 아래면 사냥꾼은 다음 전투 전에 야영하여 회복하므로, 5%까지 떨어지는 일은 드뭅니다.';

  @override
  String get repYou => '나';

  @override
  String get repLevel => '레벨';

  @override
  String get repHp => 'HP';

  @override
  String get repAttack => '공격력';

  @override
  String get repDefense => '방어력';

  @override
  String get repEvasion => '회피';

  @override
  String get repHitChance => '명중률';

  @override
  String get repCrit => '치명타';

  @override
  String repFights(int n) {
    return '전투 $n회';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '명중 $hits/$swings ($pct%), 평균 $avg, 최대 $max';
  }

  @override
  String get repYourSwings => '내 공격';

  @override
  String get repPartySwings => '파티의 공격';

  @override
  String get repEnemySwings => '적의 공격';

  @override
  String repLowest(int pct) {
    return '전투 중 최저 HP: $pct%';
  }

  @override
  String repRest(String time) {
    return '야영: $time';
  }

  @override
  String get restart => '재시작';

  @override
  String get restartTitle => '수령하고 다시 시작할까요?';

  @override
  String get restartBody => '지금까지 얻은 것은 가방으로 들어가고, 같은 작업이 같은 시간으로 다시 시작됩니다.';

  @override
  String get dailyTitle => '일일 접속 보상';

  @override
  String get dailyBody => '어느 캐릭터로든 하루 한 번 받을 수 있습니다. 주가 이어질수록 더 많이 주며 29일째가 가장 많습니다. 하루라도 놓치면 처음부터 다시 시작합니다.';

  @override
  String get dailyClaim => '받기';

  @override
  String get dailyTomorrow => '내일 다시 오세요';

  @override
  String dailyDay(int n) {
    return '$n일째';
  }

  @override
  String get dailyOpen => '일일 보상';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name님이 욕설로 $min분간 채팅 금지되었습니다.';
  }

  @override
  String get storyTitle => '스토리';

  @override
  String get storyOpen => '스토리 모드';

  @override
  String get storyComingSoon => '스토리 모드가 곧 추가됩니다. 챕터가 작성되는 대로 여기에 표시됩니다.';

  @override
  String storyLocked(int level) {
    return '레벨 $level에 열림';
  }

  @override
  String get bankTitle => '은행';

  @override
  String get bankHint => '모든 캐릭터가 공유합니다. 가방의 아이템을 누르면 보관하고, 은행의 아이템을 누르면 꺼냅니다. 귀속 아이템은 해당 캐릭터에 남습니다.';

  @override
  String get bankEmpty => '은행이 비어 있습니다.';

  @override
  String get bankDeposit => '보관';

  @override
  String get bankWithdraw => '꺼내기';

  @override
  String get errBankFull => '은행이 가득 찼습니다.';

  @override
  String get tabTrade => '거래';

  @override
  String get tradeHint => '다른 플레이어와 아이템을 거래합니다. 양쪽이 아이템을 올리고 확인하며, 각자 내주는 아이템의 평균 시세의 8%를 개당 세금으로 냅니다.';

  @override
  String get tradeWith => '플레이어 이름';

  @override
  String get tradeOpen => '거래 시작';

  @override
  String tradeTitle(String name) {
    return '$name님과 거래';
  }

  @override
  String get tradeYours => '내가 줄 것';

  @override
  String tradeTheirs(String name) {
    return '$name님이 줄 것';
  }

  @override
  String get tradeNothing => '아직 없습니다.';

  @override
  String get tradeAdd => '아이템 추가';

  @override
  String get tradeRemove => '되돌리기';

  @override
  String get tradeTax => '내 세금';

  @override
  String get tradeReady => '준비됨';

  @override
  String get tradeWaiting => '준비 안 됨';

  @override
  String get tradeConfirm => '거래 확인';

  @override
  String get tradeWaitingThem => '상대방을 기다리는 중';

  @override
  String get tradeCancel => '거래 취소';

  @override
  String get tradeCancelBody => '아무것도 오가지 않습니다.';

  @override
  String get tradeDone => '거래 완료.';

  @override
  String get errTradeBusy => '둘 중 한 명이 이미 거래 중입니다.';

  @override
  String get errTradeClosed => '이 거래는 끝났습니다.';

  @override
  String get errTradeItem => '이 아이템은 지금 거래할 수 없습니다.';

  @override
  String get errTradeEmpty => '먼저 무언가를 올리세요.';

  @override
  String get errTradeTooMany => '거래 아이템이 너무 많습니다.';

  @override
  String get errCraftLevel => '캐릭터 레벨이 이 아이템 레벨보다 낮습니다.';

  @override
  String get gatherNoneHere => '가진 자원은 다른 지도에서 납니다. 옆에 표시된 곳으로 이동하세요.';

  @override
  String get errGatherRegion => '이 자원은 다른 지도에서 납니다.';

  @override
  String get helpersOpen => '도우미';

  @override
  String get helpersTitle => '도우미';

  @override
  String get helpersHint => '채집과 제작을 도와줄 도우미를 고용하세요. 속도는 절반이고 경험치도 절반이며, 경험치는 내 직업으로 들어옵니다. 종류마다 고유 보너스가 있습니다.';

  @override
  String get helpersGather => '채집';

  @override
  String get helpersCraft => '제작';

  @override
  String get helperHire => '고용';

  @override
  String helperHireBody(String name, String price) {
    return '$name을(를) $price에 고용할까요?';
  }

  @override
  String get helperSend => '일 보내기';

  @override
  String get helperCancel => '불러오기';

  @override
  String get helperCancelBody => '채집은 아무것도 돌아오지 않으며, 제작은 재료와 골드가 돌아옵니다.';

  @override
  String get helperDone => '완료 — 수령';

  @override
  String helperWorking(String time) {
    return '작업 중, $time 남음';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return '아이템 $items개, $profession 경험치 +$xp';
  }

  @override
  String get helperNoRecipes => '이 직업에는 아직 레벨에 맞는 레시피가 없습니다.';

  @override
  String get helperPickRecipe => '레시피 선택';

  @override
  String helperBonusYield(int pct) {
    return '산출량 +$pct%';
  }

  @override
  String helperBonusSpeed(int pct) {
    return '$pct% 빠름';
  }

  @override
  String helperBonusXp(int pct) {
    return '경험치 +$pct%';
  }

  @override
  String helperBonusGold(int pct) {
    return '골드 비용 $pct% 감소';
  }

  @override
  String helperBonusDouble(int pct) {
    return '$pct% 확률로 2배 산출';
  }

  @override
  String get helperMiner => '광부';

  @override
  String get helperWoodsman => '나무꾼';

  @override
  String get helperAngler => '낚시꾼';

  @override
  String get helperHerbalist => '약초사';

  @override
  String get helperTrapper => '덫 사냥꾼';

  @override
  String get helperSmith => '대장장이';

  @override
  String get helperTanner => '무두장이';

  @override
  String get helperWeaver => '직조공';

  @override
  String get helperCarver => '목공';

  @override
  String get helperApothecary => '약제사';

  @override
  String get helperCook => '요리사';

  @override
  String get helperSmelter => '제련공';

  @override
  String get helperScribe => '필경사';

  @override
  String get helperGemcutter => '보석 세공사';

  @override
  String get errHelperBusy => '이 도우미는 이미 일하고 있습니다.';

  @override
  String get errHelperNotReady => '도우미가 아직 끝내지 않았습니다.';

  @override
  String get shopOpen => '상점';

  @override
  String get shopTitle => '크레딧 상점';

  @override
  String get shopCredits => '크레딧';

  @override
  String get shopAvatars => '아바타';

  @override
  String get shopBackgrounds => '배경';

  @override
  String shopBalance(int n) {
    return '$n 크레딧';
  }

  @override
  String get shopBuyCredits => '크레딧을 구매해 부스트와 치장 아이템을 잠금 해제하세요.';

  @override
  String get shopPaymentsSoon => '크레딧 구매 기능이 곧 추가됩니다.';

  @override
  String get shopBoosts => '경험치 부스트';

  @override
  String get shopBuy => '구매';

  @override
  String shopBuyBody(int n) {
    return '$n 크레딧을 사용할까요?';
  }

  @override
  String shopPrice(int n) {
    return '$n cr';
  }

  @override
  String get boostGather => '채집';

  @override
  String get boostAdventure => '전투';

  @override
  String get boostCraft => '직업';

  @override
  String get boostDesc => '2시간 동안 경험치 +100%';

  @override
  String get boostShort => '경험치 +100%, 2시간';

  @override
  String boostLeft(String time) {
    return '활성, $time 남음';
  }

  @override
  String boostWatchAd(int n) {
    return '광고 보기 ($n회 남음)';
  }

  @override
  String get boostAdsSoon => '광고 준비 중';

  @override
  String get errInsufficientCredits => '크레딧이 부족합니다.';

  @override
  String get errBoostActive => '이 부스트는 이미 적용 중입니다.';

  @override
  String get errAdLimit => '오늘의 광고 부스트를 모두 사용했습니다.';

  @override
  String helperCanMake(int n) {
    return '제작 가능: $n';
  }

  @override
  String get helperPickResource => '채집할 것';

  @override
  String get helperCraftHint => '재료와 골드는 지금 차감되며, 그동안 계속 플레이할 수 있습니다.';

  @override
  String get petsOpen => '펫';

  @override
  String get petsHint =>
      '펫 알은 레벨 400 이상의 레이드를 클리어하면 나옵니다. 알은 부화기에서 부화하고, 펫은 자체 사냥으로 성장하며, 활성 펫은 보너스를 줍니다. 활성 펫은 전투에도 함께하여 모든 타격에 피해를 더합니다.';

  @override
  String get petsEggs => '알';

  @override
  String get petsNoEggs => '알이 없습니다. 레벨 400 이상의 레이드를 클리어해 찾으세요.';

  @override
  String petsMine(int n, int max) {
    return '내 펫 ($n/$max)';
  }

  @override
  String get petHatch => '부화';

  @override
  String petHatchBody(String egg, String time) {
    return '$egg을(를) 부화기에 넣을까요? $time 후 부화합니다.';
  }

  @override
  String get petHatching => '알 속';

  @override
  String get petHatched => '부화!';

  @override
  String petPower(int n) {
    return '레벨 $n 몬스터처럼 싸웁니다';
  }

  @override
  String get petActive => '활성';

  @override
  String get petActivate => '활성화';

  @override
  String get petDeactivate => '휴식';

  @override
  String get petHunt => '사냥 보내기';

  @override
  String petHunting(String mob) {
    return '$mob 사냥 중';
  }

  @override
  String get petHuntDone => '사냥 종료';

  @override
  String get petCollectEarly => '지금 사냥을 끝내면 지금까지의 전투만 정산됩니다.';

  @override
  String petCollected(int wins, int xp) {
    return '$wins승, 펫 경험치 +$xp';
  }

  @override
  String petLevelUp(String name, int level) {
    return '$name 레벨 $level 달성!';
  }

  @override
  String get petDied => '펫이 쓰러져 사냥이 일찍 끝났습니다.';

  @override
  String get petPickMob => '사냥할 몬스터';

  @override
  String petXpPerWin(int n) {
    return '승리당 경험치 +$n';
  }

  @override
  String get errPetCap => '보유할 수 있는 펫 수에 도달했습니다.';

  @override
  String get errPetHatching => '펫이 아직 부화하지 않았습니다.';

  @override
  String get errPetHunting => '펫이 이미 사냥 중입니다.';

  @override
  String petStageLevels(int from, int to) {
    return 'Lv.$from–$to';
  }

  @override
  String get regionMobs => '몬스터';

  @override
  String get regionDungeons => '던전';

  @override
  String get regionGather => '자원';

  @override
  String get regionEmpty => '이 지역에 속한 몬스터, 던전, 자원, 퀘스트가 없습니다.';

  @override
  String get regionRaids => '레이드';

  @override
  String get tagUse => '용도';

  @override
  String get tagAllClasses => '모든 직업';

  @override
  String get tagDrops => '드롭';

  @override
  String get tagRaid => '레이드';

  @override
  String get tagDungeon => '던전';

  @override
  String get tagCrafted => '제작';

  @override
  String tagSkill(int n) {
    return '기술 $n';
  }

  @override
  String get tagSold => '판매';

  @override
  String get tagGathered => '채집';

  @override
  String get tagUsedIn => '재료';

  @override
  String get petEvolution => '진화';

  @override
  String get petStrength => '강함 상세';

  @override
  String petAssistRange(int min, int max) {
    return '전투 지원: 모든 타격에 피해 +$min% ~ +$max%';
  }

  @override
  String petBonusRange(String min, String max) {
    return '보너스: $min ~ $max';
  }

  @override
  String petFightsRange(int min, int max) {
    return '레벨 $min~$max 몬스터처럼 사냥합니다';
  }

  @override
  String get petFind => '획득처';

  @override
  String petAssist(int pct) {
    return '전투 지원: 피해 +$pct%';
  }

  @override
  String get petLevelNote => '레벨 1부터 100까지';

  @override
  String get petIndex => '펫 도감';

  @override
  String petIndexCount(int n, int max) {
    return '$n / $max 발견';
  }

  @override
  String slotCount(int n, int max) {
    return '$n / $max';
  }

  @override
  String get petEggTapHint => '알을 누르면 부화기에 넣을 수 있습니다.';

  @override
  String get petTapHint => '펫을 누르면 볼 수 있습니다 · 초록 점은 활성 펫입니다.';

  @override
  String get petIndexHint => '펫을 누르면 페이지가 열립니다.';

  @override
  String get petUnknown => '아직 발견하지 못함';

  @override
  String get petAbout => '소개';

  @override
  String get petHabitat => '서식지';

  @override
  String petElement(String name) {
    return '속성: $name';
  }

  @override
  String get petStage1 => '새끼';

  @override
  String get petStage2 => '어린 개체';

  @override
  String get petStage3 => '성체';

  @override
  String get bodyMale => '남성';

  @override
  String get bodyFemale => '여성';

  @override
  String get bodyChoice => '외모';

  @override
  String get bodyChange => '외모 변경';

  @override
  String get portraitNoArmor => '방어구 없음';

  @override
  String get combatPower => '전투력';

  @override
  String get combatTapHint => '눌러서 자세히 보기';

  @override
  String statCap(int cap) {
    return '상한 $cap%';
  }

  @override
  String get cEvasion => '회피';

  @override
  String get cBlock => '방어';

  @override
  String get cCritDmg => '치명타 피해';

  @override
  String get cDmgReduction => '피해 감소';

  @override
  String get cLifesteal => '흡혈';

  @override
  String get cExecute => '처형';

  @override
  String get cMissingHp => '최후의 저항';

  @override
  String get cPetAssist => '펫 지원';

  @override
  String get xPower => '월드 이벤트에서 쓰는 힘: 공격력 2배 + 방어력 + 최대 HP / 10.';

  @override
  String get xAttack => '적의 방어력을 적용하기 전 각 타격의 기본 피해.';

  @override
  String get xDefense => '받는 타격마다 방어력의 절반만큼 줄어듭니다.';

  @override
  String xCrit(int cap) {
    return '타격이 치명타가 될 확률. 최대 $cap%.';
  }

  @override
  String get xCritDmg => '치명타는 일반 타격의 이 배수만큼 피해를 줍니다.';

  @override
  String get xHit => '타격이 맞을 확률에서 적의 회피를 뺀 값 (15~95%).';

  @override
  String xEvasion(int cap) {
    return '적의 명중률에서 빠집니다. 최대 $cap%.';
  }

  @override
  String xBlock(int cap) {
    return '받는 타격을 절반으로 줄일 확률. 최대 $cap%.';
  }

  @override
  String xDmgReduction(int cap) {
    return '받는 타격이 이 비율만큼 줄어듭니다. 최대 $cap%.';
  }

  @override
  String get xLifesteal => '준 피해의 이 비율만큼 회복합니다.';

  @override
  String xExecute(int pct) {
    return 'HP $pct% 이하인 적에게 주는 추가 피해.';
  }

  @override
  String get xMissingHp => '자신의 HP가 줄수록 이 값까지 추가 피해.';

  @override
  String get xPetAssist => '활성 펫이 내 타격마다 이 비율을 더합니다.';

  @override
  String get lookFace => '얼굴';

  @override
  String get lookHairStyle => '머리 모양';

  @override
  String get lookHairColor => '머리 색';

  @override
  String get lookEyeColor => '눈 색';

  @override
  String get lookSave => '저장';

  @override
  String get lookHumanOnly => '이 종족의 얼굴과 머리 옵션은 곧 추가됩니다.';

  @override
  String get raceAbilities => '종족 능력';

  @override
  String get raceAbilitiesHint => '모두 패시브';

  @override
  String get chooseClass => '직업 선택';

  @override
  String raceAttrBonus(int pct, String name) {
    return '$name +$pct%';
  }

  @override
  String raceProfessionXp(int pct, String name) {
    return '$name 보너스 경험치 +$pct%';
  }

  @override
  String get language => '언어';

  @override
  String get report => '신고';

  @override
  String reportTitle(String name) {
    return '$name 신고';
  }

  @override
  String get reportWhy => '무슨 일이 있었나요?';

  @override
  String get reportHarassment => '괴롭힘 또는 욕설';

  @override
  String get reportCheating => '부정행위 또는 봇';

  @override
  String get reportName => '부적절한 이름';

  @override
  String get reportSpam => '도배 또는 광고';

  @override
  String get reportScam => '사기';

  @override
  String get reportOther => '기타';

  @override
  String get reportNote => '상세 내용 (선택)';

  @override
  String get reportSent => '신고가 접수되었습니다. 감사합니다.';

  @override
  String get errReportSelf => '자신의 캐릭터는 신고할 수 없습니다.';

  @override
  String get nkReportReviewed => '신고가 검토되었습니다';

  @override
  String get guide => '안내서';

  @override
  String get guideSearch => '안내서 검색';

  @override
  String get guideNoMatch => '결과가 없습니다.';

  @override
  String get guideShowMe => '바로 가기';

  @override
  String get next => '다음';

  @override
  String get back => '이전';

  @override
  String get tutorialSkip => '튜토리얼 건너뛰기';

  @override
  String get tutorialStart => '게임 시작';

  @override
  String get tutorialReplay => '튜토리얼 다시 보기';

  @override
  String get tutorialWelcomeTitle => 'MeNa World에 오신 것을 환영합니다';

  @override
  String get tutorialWelcomeBody => '앱을 닫아도 캐릭터는 계속 일합니다. 할 일을 맡기고 나중에 돌아와 보상을 받으세요.';

  @override
  String get tutorialCharacterBody => '능력치, 특성 포인트, 스킬을 여기서 확인합니다. 레벨마다 얻는 포인트를 분배하고 300, 700레벨에서 승급을 선택하세요.';

  @override
  String get tutorialGatherBody => '자원과 시간을 고른 뒤 돌아와서 수령하세요. 채광, 낚시, 벌목, 약초 채집, 사냥은 각각 따로 레벨이 오릅니다.';

  @override
  String get tutorialBagBody => '아이템을 눌러 장착, 사용, 판매하거나 거래소에 올리세요. 비교표에서 착용 중인 장비보다 나은지 볼 수 있습니다.';

  @override
  String get tutorialAdventureBody => '몬스터 사냥, 던전 공략, 레이드와 이벤트 참여, 지도 여행을 합니다. 몬스터와 던전은 현재 지역에 따라 달라집니다.';

  @override
  String get tutorialWorkshopBody => '전문 기술로 제작하고, 상점에서 구매하고, 거래소에서 다른 플레이어와 거래하세요.';

  @override
  String get tutorialTopBarTitle => '상단 바';

  @override
  String get tutorialTopBarBody =>
      '퀘스트, 알림, 소셜(채팅, 친구, 우편, 길드, 교환)은 오른쪽 위에 있습니다. 캐릭터는 한 번에 한 가지 일만 하므로 새 일을 시작하기 전에 끝내거나 수령하세요.';

  @override
  String get tutorialDoneTitle => '준비 완료';

  @override
  String get tutorialDoneBody => '위치가 기억나지 않으면 캐릭터 탭 아래쪽의 안내서를 여세요.';

  @override
  String get guideCharacterTitle => '능력치, 특성, 스킬';

  @override
  String get guideCharacterBody => '캐릭터 탭. 진행 카드에서 특성 포인트를 분배합니다. 스킬은 접을 수 있는 영역에 있으며, 같은 역할의 스킬은 가장 강한 것만 남습니다.';

  @override
  String get guidePromotionTitle => '승급';

  @override
  String get guidePromotionBody => '캐릭터 탭의 진행 카드. 300, 700레벨에서는 두 경로 중 하나를 고를 때까지 경험치가 대기합니다. 그동안 얻은 경험치는 선택 후 지급됩니다.';

  @override
  String get guideGatherTitle => '채집';

  @override
  String get guideGatherBody => '채집 탭. 전문 기술, 도달 가능한 자원, 1·3·6시간을 고릅니다. 자원은 지역에 속하므로 더 좋은 자원은 여행해서 찾으세요.';

  @override
  String get guideBagTitle => '가방과 장비';

  @override
  String get guideBagBody =>
      '가방 탭. 아이템을 누르면 상세 정보와 비교가 나옵니다. 길게 눌러 끌면 여러 개를 골라 함께 판매하거나 파괴합니다. 가방이 가득 차면 새 작업을 시작할 수 없습니다.';

  @override
  String get guideBankTitle => '계정 은행';

  @override
  String get guideBankBody => '가방 탭의 은행 버튼. 모든 캐릭터가 함께 쓰므로 아이템을 옮길 수 있습니다. 귀속되었거나 착용 중인 아이템은 넣을 수 없습니다.';

  @override
  String get guideHuntTitle => '사냥';

  @override
  String get guideHuntBody => '모험 탭의 사냥. 도달 가능한 몬스터, 시간, 가져갈 물약과 음식을 고릅니다. 캐릭터는 크게 다치면 스스로 후퇴합니다.';

  @override
  String get guideDungeonTitle => '던전과 레이드';

  @override
  String get guideDungeonBody => '모험 탭의 던전. 던전은 혼자 합니다. 레이드는 2~5인 파티가 필요하며 파티 찾기로 매칭됩니다. 레이드 보스는 세트 아이템과 펫 알을 떨어뜨립니다.';

  @override
  String get guideTravelTitle => '여행과 지도';

  @override
  String get guideTravelBody => '모험 탭의 여행. 지역을 누르면 몬스터, 던전, 자원, 퀘스트가 보이고 이동할 수 있습니다. 레벨 조건은 목적지에만 있습니다.';

  @override
  String get guideEventsTitle => '이벤트';

  @override
  String get guideEventsBody => '모험 탭의 이벤트. 월드 보스와 침공에 참가하고 매일 시즌 토큰을 받으세요. 카드를 누르면 방법과 보상이 나옵니다.';

  @override
  String get guideCraftTitle => '전문 기술과 제작';

  @override
  String get guideCraftBody => '공방 탭의 제작. 전문 기술을 열면 제작법이 보입니다. 색은 해당 제작법이 주는 경험치를 나타냅니다. 마법 부여 두루마리는 가방에서 사용합니다.';

  @override
  String get guideVendorTitle => '상점과 재구매';

  @override
  String get guideVendorBody => '공방 탭의 상점. 판매는 가방에서 합니다. 판 물건은 2시간 동안 다시 살 수 있습니다.';

  @override
  String get guideMarketTitle => '거래소와 경매';

  @override
  String get guideMarketBody =>
      '공방 탭의 거래소. 고정가로 구매하거나 경매에 입찰하세요. 가방에서 아이템을 올립니다: 고정가 또는 시작가와 선택적 즉시 구매가가 있는 경매. 마지막 몇 분의 입찰은 경매를 연장합니다.';

  @override
  String get guideQuestsTitle => '퀘스트';

  @override
  String get guideQuestsBody => '오른쪽 위 클립보드 아이콘. 스토리, 서브, 일일 퀘스트가 있으며 배지는 완료 가능한 수를 보여 줍니다.';

  @override
  String get guideNotificationsTitle => '알림과 보상';

  @override
  String get guideNotificationsBody => '오른쪽 위 종 아이콘. 보상 탭에는 받은 내용이 30일간 남습니다. 설정에서 받을 알림을 고릅니다.';

  @override
  String get guideSocialTitle => '채팅, 친구, 우편, 교환';

  @override
  String get guideSocialBody => '오른쪽 위 말풍선. 교환은 양쪽이 확인해야 하며 세금이 붙습니다. 차단하면 채팅, 친구, 우편에서 해당 플레이어가 사라집니다.';

  @override
  String get guidePartyTitle => '파티';

  @override
  String get guidePartyBody => '캐릭터 탭의 파티 영역. 이름으로 초대합니다. 같은 종류의 일을 하는 파티원은 추가 경험치를 얻습니다.';

  @override
  String get guideGuildTitle => '길드, 금고, 길드 은행';

  @override
  String get guideGuildBody => '오른쪽 위 말풍선의 길드 탭. 길드원은 골드와 아이템을 맡기고, 등급에 따라 하루 인출 가능 여부와 한도가 정해집니다. 모든 내역이 기록됩니다.';

  @override
  String get guideRankingTitle => '랭킹';

  @override
  String get guideRankingBody => '캐릭터 탭의 랭킹 버튼. 점수는 레벨, 업적, 숙련, 전문 기술, 던전과 레이드 공략으로 쌓입니다. 직업이나 표별로 볼 수 있습니다.';

  @override
  String get guidePetsTitle => '펫';

  @override
  String get guidePetsBody => '캐릭터 탭의 펫 버튼. 알은 고레벨 레이드 보스에게서 나옵니다. 부화시켜 사냥으로 키우고, 보너스와 전투 지원을 위해 하나를 활성화하세요.';

  @override
  String get guideHelpersTitle => '도우미';

  @override
  String get guideHelpersBody => '캐릭터 탭의 도우미 버튼. 전문 기술별로 도우미를 고용합니다. 캐릭터가 무엇을 하든 따로 일합니다.';

  @override
  String get guideReportTitle => '플레이어 신고';

  @override
  String get guideReportBody => '채팅 줄을 길게 누르거나 친구 메뉴를 사용하세요. 사유를 고르면 팀이 검토한 뒤 알려 드립니다.';

  @override
  String get guideAccountTitle => '계정, 언어, 고객 지원';

  @override
  String get guideAccountBody => '캐릭터 목록 위쪽의 계정 아이콘. 언어나 비밀번호 변경, 고객 지원 문의, 데이터 다운로드, 계정 삭제를 할 수 있습니다.';

  @override
  String get ranking => '랭킹';

  @override
  String get rankingScore => '점수';

  @override
  String get rankingLevel => '레벨';

  @override
  String get rankingAchievements => '업적';

  @override
  String get rankingMastery => '숙련';

  @override
  String get rankingProfessions => '전문 기술';

  @override
  String get rankingClears => '공략';

  @override
  String get rankingDungeons => '던전';

  @override
  String get rankingRaids => '레이드';

  @override
  String get rankingAllClasses => '모든 직업';

  @override
  String get rankingHowScore => '점수 계산 방법';

  @override
  String get rankingScoreNote => '던전과 레이드는 종류별로 한 번만 계산됩니다. 몇 분마다 갱신됩니다.';

  @override
  String get rankingEmpty => '순위표를 집계하는 중입니다.';

  @override
  String get guildVault => '금고와 은행';

  @override
  String get guildVaultTreasury => '길드 금고';

  @override
  String get guildVaultBank => '길드 은행';

  @override
  String get guildVaultDepositGold => '골드 입금';

  @override
  String get guildVaultWithdrawGold => '골드 출금';

  @override
  String get guildVaultGoldAmount => '골드 수량';

  @override
  String get guildVaultNoLimit => '일일 한도 없이 출금할 수 있습니다.';

  @override
  String guildVaultGoldToday(String used, String limit) {
    return '오늘 출금: $used / $limit';
  }

  @override
  String guildVaultItemsToday(int used, int limit) {
    return '오늘 꺼낸 아이템: $used / $limit';
  }

  @override
  String get guildVaultNoItems => '현재 등급으로는 길드 은행을 사용할 수 없습니다.';

  @override
  String get guildVaultLog => '내역';

  @override
  String get guildVaultLogEmpty => '아직 내역이 없습니다.';

  @override
  String guildVaultLogIn(String name, String what) {
    return '$name님이 $what 입금';
  }

  @override
  String guildVaultLogOut(String name, String what) {
    return '$name님이 $what 출금';
  }

  @override
  String get errGuildVaultNotEmpty => '해산하기 전에 금고와 은행을 비우세요.';

  @override
  String get errGuildVaultClosing => '길드가 해산 중이라 입금할 수 없습니다.';

  @override
  String get errGuildVaultLimit => '일일 한도를 초과합니다.';

  @override
  String get errItemBound => '귀속 아이템은 옮길 수 없습니다.';

  @override
  String get errItemEquipped => '먼저 아이템을 해제하세요.';

  @override
  String get bid => '입찰';

  @override
  String get bidPlaced => '입찰했습니다.';

  @override
  String get marketFixed => '고정가';

  @override
  String get marketAuction => '경매';

  @override
  String auctionNoBids(String price) {
    return '입찰 없음, 시작가 $price';
  }

  @override
  String auctionTopBid(String price, int count) {
    return '최고 입찰 $price ($count회)';
  }

  @override
  String get auctionYouLead => '내가 선두';

  @override
  String get auctionStartBid => '전체 묶음 시작가';

  @override
  String get auctionBuyoutOptional => '즉시 구매가 (선택)';

  @override
  String get auctionRules => '시간이 끝나면 최고 입찰자가 낙찰됩니다. 마지막 5분 안의 입찰은 경매를 연장합니다. 입찰이 들어온 뒤에는 취소할 수 없습니다.';

  @override
  String get auctionYouReceiveAtLeast => '최소 수령액';

  @override
  String auctionMinBid(String price) {
    return '최소 입찰가: $price';
  }

  @override
  String get auctionBidCopper => '입찰액 (구리)';

  @override
  String get auctionHeldNote => '골드는 보류되며 다른 사람이 더 높게 입찰하면 바로 돌려받습니다.';

  @override
  String auctionBuyoutFor(String price) {
    return '즉시 구매: $price';
  }

  @override
  String get errBidTooLow => '최소 입찰가보다 낮습니다.';

  @override
  String get errTopBidder => '이미 최고 입찰자입니다.';

  @override
  String get errAuctionBusy => '다른 입찰을 처리 중입니다. 잠시 후 다시 시도하세요.';

  @override
  String get errAuctionHasBids => '입찰이 있는 경매는 취소할 수 없습니다.';

  @override
  String get nkMarketOutbid => '입찰가를 추월당했습니다';

  @override
  String get nkMarketWon => '경매에서 낙찰되었습니다';

  @override
  String get evMarketOutbid => '입찰가를 추월당해 골드가 반환되었습니다.';

  @override
  String get skillNeedsShield => '방패 필요';
}
