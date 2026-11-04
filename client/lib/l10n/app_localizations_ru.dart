// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Russian (`ru`).
class AppLocalizationsRu extends AppLocalizations {
  AppLocalizationsRu([String locale = 'ru']) : super(locale);

  @override
  String get appTitle => 'MeNa World';

  @override
  String get tagline => 'Приключения продолжаются, пока вас нет.';

  @override
  String get emailOrUsername => 'Эл. почта или имя пользователя';

  @override
  String get email => 'Эл. почта';

  @override
  String get username => 'Имя пользователя';

  @override
  String get password => 'Пароль';

  @override
  String get passwordConfirm => 'Подтвердите пароль';

  @override
  String get showPassword => 'Показать пароль';

  @override
  String get hidePassword => 'Скрыть пароль';

  @override
  String get login => 'Войти';

  @override
  String get register => 'Создать аккаунт';

  @override
  String get noAccount => 'Впервые здесь? Создайте аккаунт';

  @override
  String get haveAccount => 'Уже есть аккаунт? Войдите';

  @override
  String get acceptTerms => 'Мне исполнилось 13 лет, и я принимаю Условия использования и Политику конфиденциальности';

  @override
  String get fieldRequired => 'Обязательное поле';

  @override
  String get invalidEmail => 'Введите корректный адрес эл. почты';

  @override
  String get usernameRule => '3–32 символа: буквы, цифры, _ или .';

  @override
  String get passwordMin => 'Минимум 8 символов';

  @override
  String get passwordMismatch => 'Пароли не совпадают';

  @override
  String get termsRequired => 'Чтобы продолжить, примите условия';

  @override
  String get logout => 'Выйти';

  @override
  String get retry => 'Повторить';

  @override
  String get offline => 'Нет связи с сервером. Проверьте подключение.';

  @override
  String get characters => 'Персонажи';

  @override
  String slotsUsed(int used, int cap) {
    return 'Слоты: $used/$cap';
  }

  @override
  String get newCharacter => 'Новый персонаж';

  @override
  String get noCharacters => 'Персонажей пока нет. Создайте первого героя.';

  @override
  String get create => 'Создать';

  @override
  String get characterName => 'Имя';

  @override
  String get nameRule => 'До 16 символов одного алфавита, начинается с буквы (можно цифры и _)';

  @override
  String get race => 'Раса';

  @override
  String get characterClass => 'Класс';

  @override
  String get raceHuman => 'Человек';

  @override
  String get raceElf => 'Эльф';

  @override
  String get raceDwarf => 'Дворф';

  @override
  String get raceOrc => 'Орк';

  @override
  String get classWarrior => 'Воин';

  @override
  String get classMage => 'Маг';

  @override
  String get classRogue => 'Разбойник';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Ур. $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Ур. $level';
  }

  @override
  String get hp => 'ОЗ';

  @override
  String get mp => 'ОМ';

  @override
  String get stamina => 'ВН';

  @override
  String get energy => 'ЭН';

  @override
  String get xp => 'Опыт';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => 'Характеристики';

  @override
  String get combat => 'Бой';

  @override
  String get statStr => 'СИЛ';

  @override
  String get statAgi => 'ЛОВ';

  @override
  String get statInt => 'ИНТ';

  @override
  String get statVit => 'ВЫН';

  @override
  String get statCha => 'МДР';

  @override
  String get statLuck => 'УДЧ';

  @override
  String get attack => 'Атака';

  @override
  String get defense => 'Защита';

  @override
  String get crit => 'Крит';

  @override
  String get hit => 'Точность';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => 'Золото';

  @override
  String get switchCharacter => 'Сменить персонажа';

  @override
  String get statePendingDelete => 'Ожидает удаления';

  @override
  String get stateFrozen => 'Заморожен';

  @override
  String get errInvalidCredentials => 'Неверная эл. почта, имя пользователя или пароль.';

  @override
  String get errAccountLocked => 'Ваш аккаунт временно заблокирован.';

  @override
  String get errAccountBanned => 'Этот аккаунт забанен.';

  @override
  String get errAccountExists => 'Эта эл. почта или имя пользователя уже заняты.';

  @override
  String get errRateLimited => 'Слишком много попыток. Подождите немного.';

  @override
  String get errNameTaken => 'Это имя уже занято.';

  @override
  String get errSlotCap => 'Все слоты персонажей заняты.';

  @override
  String get errNotPlayable => 'Этим персонажем сейчас нельзя играть.';

  @override
  String get errSessionExpired => 'Сессия истекла. Войдите снова.';

  @override
  String errGeneric(String code) {
    return 'Что-то пошло не так ($code).';
  }

  @override
  String get navCharacter => 'Персонаж';

  @override
  String get navActivity => 'Сбор';

  @override
  String get actMining => 'Добыча руды';

  @override
  String get actFishing => 'Рыбалка';

  @override
  String get actWoodcutting => 'Лесозаготовка';

  @override
  String get actHerbalism => 'Травничество';

  @override
  String get actHunting => 'Охота';

  @override
  String get activityChoose => 'Выберите занятие';

  @override
  String get durationLabel => 'Длительность';

  @override
  String get gatherResource => 'Ресурс';

  @override
  String gatherXp(int xp, int next) {
    return 'Опыт: $xp / $next';
  }

  @override
  String get gatherMax => 'Макс. уровень';

  @override
  String gatherLocked(int level) {
    return 'Откроется на уровне $level';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name: +$xp опыта';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '$name: достигнут уровень $level!';
  }

  @override
  String get errGatherLevel => 'Ваш уровень сбора слишком низок для этого ресурса.';

  @override
  String startActivity(String name) {
    return 'Начать: $name';
  }

  @override
  String get collect => 'Забрать';

  @override
  String get abandon => 'Бросить';

  @override
  String get cancel => 'Отмена';

  @override
  String get soFar => 'Собрано на данный момент';

  @override
  String get activityDone => 'Готово — можно забрать';

  @override
  String remaining(String time) {
    return 'Осталось: $time';
  }

  @override
  String durHm(int h, String m) {
    return '$h ч $m мин';
  }

  @override
  String durM(int m) {
    return '$m мин';
  }

  @override
  String durS(int s) {
    return '$s с';
  }

  @override
  String get collectEarlyTitle => 'Забрать досрочно?';

  @override
  String get collectEarlyBody => 'Вы получите то, что уже собрано, и занятие закончится.';

  @override
  String get abandonTitle => 'Бросить занятие?';

  @override
  String get abandonBody => 'Всё собранное будет потеряно.';

  @override
  String get lootCollect => 'Отлично!';

  @override
  String get lootNothing => 'На этот раз ничего.';

  @override
  String xpGain(int xp) {
    return '+$xp опыта';
  }

  @override
  String get rarityCommon => 'Обычное';

  @override
  String get rarityUncommon => 'Необычное';

  @override
  String get rarityRare => 'Редкое';

  @override
  String get rarityEpic => 'Эпическое';

  @override
  String get rarityLegendary => 'Легендарное';

  @override
  String get rarityMythic => 'Мифическое';

  @override
  String get errSlotBusy => 'Уже идёт другое занятие.';

  @override
  String get errNotActive => 'Это занятие уже завершено.';

  @override
  String get evAchievement => 'Получено достижение!';

  @override
  String get evMasteryTier => 'Новая ступень мастерства!';

  @override
  String get evMarketSold => 'Ваш лот продан.';

  @override
  String get evMailReceived => 'Вам пришла почта.';

  @override
  String get evPartyInvite => 'Вас пригласили в отряд.';

  @override
  String get evGuildInvite => 'Вас пригласили в гильдию.';

  @override
  String get navBag => 'Сумка';

  @override
  String get equipment => 'Снаряжение';

  @override
  String get bag => 'Сумка';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap кг';
  }

  @override
  String get slotHelm => 'Шлем';

  @override
  String get slotNecklace => 'Шея';

  @override
  String get slotChest => 'Грудь';

  @override
  String get slotGloves => 'Перчатки';

  @override
  String get slotLegs => 'Ноги';

  @override
  String get slotBoots => 'Сапоги';

  @override
  String get slotMainhand => 'Основная рука';

  @override
  String get slotOffhand => 'Вторая рука';

  @override
  String get slotRing => 'Кольцо';

  @override
  String get slotTrinket => 'Талисман';

  @override
  String get equip => 'Надеть';

  @override
  String get unequip => 'Снять';

  @override
  String get sell => 'Продать';

  @override
  String get sortLabel => 'Сортировка';

  @override
  String get sortSlot => 'Порядок в сумке';

  @override
  String get sortRarity => 'Редкость';

  @override
  String get sortType => 'Тип';

  @override
  String get sortName => 'Название';

  @override
  String get sortValue => 'Цена';

  @override
  String selectedCount(int count) {
    return 'Выбрано: $count';
  }

  @override
  String get selectAll => 'Выбрать все';

  @override
  String get clearSelection => 'Снять выбор';

  @override
  String sellSelectedTitle(int count) {
    return 'Продать предметы ($count)?';
  }

  @override
  String destroySelectedTitle(int count) {
    return 'Уничтожить предметы ($count)?';
  }

  @override
  String get deleteEntry => 'Удалить';

  @override
  String craftHowMany(String name) {
    return 'Сколько изготовить: $name?';
  }

  @override
  String craftCount(int count) {
    return 'Изготовить ×$count';
  }

  @override
  String get destroy => 'Уничтожить';

  @override
  String sellTitle(String name) {
    return 'Продать: $name?';
  }

  @override
  String get sellBody => 'Выкупить предмет у торговца можно в течение 2 часов.';

  @override
  String destroyTitle(String name) {
    return 'Уничтожить: $name?';
  }

  @override
  String get destroyBody => 'Это действие нельзя отменить.';

  @override
  String get soldFor => 'Продано за';

  @override
  String get quantity => 'Количество';

  @override
  String get requiresLevel => 'Требуемый уровень';

  @override
  String get weight => 'Вес';

  @override
  String weightGrams(int g) {
    return '$g г';
  }

  @override
  String get sellValue => 'Цена продажи';

  @override
  String get bound => 'Персональный';

  @override
  String get itemLocked => 'Предмет занят обменом, почтой или лотом.';

  @override
  String get errBagFull => 'Ваша сумка заполнена.';

  @override
  String get errItemLocked => 'Этот предмет заблокирован.';

  @override
  String get errClassRestriction => 'Ваш класс не может использовать этот предмет.';

  @override
  String get errLevelRequirement => 'Ваш уровень слишком низок.';

  @override
  String get errNotEquippable => 'Этот предмет нельзя надеть.';

  @override
  String get navAdventure => 'Приключения';

  @override
  String get segHunt => 'Охота';

  @override
  String get segDungeon => 'Подземелье';

  @override
  String get segTravel => 'Путешествие';

  @override
  String get huntChoose => 'Выберите добычу';

  @override
  String get startHunt => 'Начать охоту';

  @override
  String get huntDone => 'Охота окончена — заберите трофеи';

  @override
  String get huntEarlyBody => 'Вы сохраните уже выигранные бои, и охота закончится.';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'ОЗ $hp · АТК $atk · ЗАЩ $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return 'Побед: $won · Поражений: $lost';
  }

  @override
  String levelUp(int level) {
    return 'Новый уровень! Теперь уровень $level';
  }

  @override
  String levelReq(int level) {
    return 'Ур. $level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return 'Бои: $count · босс $boss';
  }

  @override
  String get dungeonEnter => 'Войти';

  @override
  String get dungeonDone => 'Забег завершён — заберите награды';

  @override
  String get dungeonWaitHint => 'Награды откроются по окончании забега.';

  @override
  String get dungeonCleared => 'Пройдено!';

  @override
  String get dungeonFailed => 'Отряд погиб';

  @override
  String get bossKilled => 'Босс повержен';

  @override
  String get bonus => 'Бонус за прохождение';

  @override
  String get youAreIn => 'Вы находитесь в';

  @override
  String get routes => 'Маршруты';

  @override
  String get noRoutes => 'Отсюда нет маршрутов.';

  @override
  String get travelTo => 'В путь';

  @override
  String travelling(String place) {
    return 'Путь в: $place';
  }

  @override
  String get arrive => 'Прибыть';

  @override
  String get arrived => 'Вы прибыли';

  @override
  String get cancelTravel => 'Повернуть назад';

  @override
  String get kindCity => 'Город';

  @override
  String get kindZone => 'Зона';

  @override
  String get kindWilderness => 'Дикие земли';

  @override
  String get kindDungeon => 'Подземелье';

  @override
  String get kindContinent => 'Континент';

  @override
  String get safeZone => 'безопасно';

  @override
  String get errBusy => 'Что-то уже выполняется.';

  @override
  String get errNotReady => 'Ещё не завершено.';

  @override
  String get errWeeklyLocked => 'Вы уже прошли это подземелье на этой неделе.';

  @override
  String get errNoRoute => 'Туда нет прямого маршрута.';

  @override
  String get navWorkshop => 'Мастерская';

  @override
  String get segCraft => 'Создание';

  @override
  String get segShop => 'Лавка';

  @override
  String get profBlacksmithing => 'Кузнечное дело';

  @override
  String get profAlchemy => 'Алхимия';

  @override
  String get profCooking => 'Кулинария';

  @override
  String get learn => 'Изучить';

  @override
  String recipeCount(int count) {
    return 'Рецептов: $count';
  }

  @override
  String get itemType => 'Тип';

  @override
  String get stackSize => 'Стопка';

  @override
  String get skillXpReward => 'Опыт навыка';

  @override
  String get learnHint => 'Изучите эту профессию, чтобы открыть рецепты.';

  @override
  String get craft => 'Создать';

  @override
  String crafted(String name, int count) {
    return 'Создано: $name ×$count';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession: теперь уровень $level';
  }

  @override
  String requiresSkill(int level) {
    return 'Нужен навык $level';
  }

  @override
  String get buy => 'Купить';

  @override
  String bought(String name) {
    return 'Куплено: $name';
  }

  @override
  String get buybacks => 'Выкуп';

  @override
  String get buybackHint => 'Проданные предметы можно выкупить в течение 2 часов.';

  @override
  String get errInsufficientGold => 'Недостаточно золота.';

  @override
  String get errInsufficientMaterials => 'Не хватает материалов.';

  @override
  String get errSkillTooLow => 'Ваш навык слишком низок для этого рецепта.';

  @override
  String get errProfessionNotLearned => 'Сначала изучите профессию.';

  @override
  String get errPriceChanged => 'Цена изменилась — проверьте её и повторите попытку.';

  @override
  String get errCannotSell => 'Этот предмет нельзя продать.';

  @override
  String get errBuybackExpired => 'Время выкупа истекло.';

  @override
  String get segMarket => 'Рынок';

  @override
  String get classWeapon => 'Оружие';

  @override
  String get classArmor => 'Броня';

  @override
  String get classConsumable => 'Расходуемое';

  @override
  String get classMaterial => 'Материалы';

  @override
  String get filterAll => 'Все';

  @override
  String get sortCheapest => 'Сначала дешёвые';

  @override
  String get sortNewest => 'Сначала новые';

  @override
  String get marketEmpty => 'Подходящих лотов нет.';

  @override
  String listingLine(String seller, String price) {
    return '$seller · $price за шт.';
  }

  @override
  String get loadMore => 'Загрузить ещё';

  @override
  String myListings(int count) {
    return 'Мои лоты ($count)';
  }

  @override
  String get cancelListing => 'Снять лот';

  @override
  String cancelListingTitle(String name) {
    return 'Снять с продажи: $name?';
  }

  @override
  String get cancelListingBody => 'Предмет вернётся в сумку. Часть комиссии возвращается только в первые 6 часов.';

  @override
  String buyListingTitle(String name) {
    return 'Купить: $name?';
  }

  @override
  String buyListingBody(String price, String seller) {
    return 'Заплатить $price продавцу $seller.';
  }

  @override
  String get listOnMarket => 'Выставить на рынок';

  @override
  String listed(String name) {
    return '$name: выставлено на рынок';
  }

  @override
  String quantityOf(int count, int max) {
    return 'Количество: $count из $max';
  }

  @override
  String get unitPrice => 'Цена за штуку';

  @override
  String suggestedPrice(String price) {
    return 'Рыночная цена: $price';
  }

  @override
  String get listingFee => 'Комиссия за лот (не возвращается)';

  @override
  String get youReceive => 'Вы получите при продаже';

  @override
  String get errPriceOutOfRange => 'Эта цена вне допустимого диапазона.';

  @override
  String get errListingCap => 'У вас слишком много активных лотов.';

  @override
  String get errSelfBuy => 'Нельзя купить собственный лот.';

  @override
  String get errListingUnavailable => 'Этот лот больше недоступен.';

  @override
  String get errNotTradable => 'Этот предмет нельзя передать.';

  @override
  String get social => 'Общение';

  @override
  String get tabChat => 'Чат';

  @override
  String get tabFriends => 'Друзья';

  @override
  String get tabMail => 'Почта';

  @override
  String get chanGlobal => 'Мир';

  @override
  String get chanGuild => 'Гильдия';

  @override
  String get chanParty => 'Отряд';

  @override
  String get chatEmpty => 'Сообщений пока нет. Поздоровайтесь!';

  @override
  String chatHint(String channel) {
    return 'Сообщение: $channel';
  }

  @override
  String get send => 'Отправить';

  @override
  String get addFriend => 'Добавить';

  @override
  String friendAdded(String name) {
    return 'Вы с $name теперь друзья';
  }

  @override
  String friendRequested(String name) {
    return 'Запрос в друзья отправлен: $name';
  }

  @override
  String get friendRequests => 'Запросы в друзья';

  @override
  String friendsCount(int count) {
    return 'Друзья ($count)';
  }

  @override
  String get noFriends => 'Друзей пока нет — добавьте кого-нибудь по имени.';

  @override
  String get more => 'Ещё';

  @override
  String get whisper => 'Шёпот';

  @override
  String get remove => 'Убрать';

  @override
  String get block => 'Заблокировать';

  @override
  String get unblock => 'Разблокировать';

  @override
  String get blockedList => 'Заблокированные';

  @override
  String removeFriendTitle(String name) {
    return 'Убрать из друзей: $name?';
  }

  @override
  String get removeFriendBody => 'Позже вы сможете отправить новый запрос.';

  @override
  String blockTitle(String name) {
    return 'Заблокировать: $name?';
  }

  @override
  String get blockBody =>
      'Этот игрок не сможет писать вам шёпотом, его сообщения в чате будут скрыты, а его почта не будет доставляться. Дружба также прекратится.';

  @override
  String get accept => 'Принять';

  @override
  String get decline => 'Отклонить';

  @override
  String get newMail => 'Написать письмо';

  @override
  String get inboxEmpty => 'Ваш почтовый ящик пуст.';

  @override
  String get noSubject => '(без темы)';

  @override
  String mailFrom(String name) {
    return 'От: $name';
  }

  @override
  String get attachments => 'Вложения';

  @override
  String get claim => 'Забрать вложения';

  @override
  String get mailClaimed => 'Вложения перемещены в сумку и кошелёк';

  @override
  String get deleteMailTitle => 'Удалить это письмо?';

  @override
  String get deleteMailBody => 'Удалённое письмо нельзя восстановить.';

  @override
  String mailSent(String postage) {
    return 'Письмо отправлено (пошлина $postage)';
  }

  @override
  String get mailTo => 'Кому (имя персонажа)';

  @override
  String get mailSubject => 'Тема';

  @override
  String get mailBody => 'Сообщение';

  @override
  String get attachGold => 'Вложить золото';

  @override
  String get errAlreadyFriends => 'Вы уже друзья.';

  @override
  String get errAlreadyRequested => 'Запрос в друзья уже отправлен.';

  @override
  String get errBlocked => 'Вы не можете взаимодействовать с этим игроком.';

  @override
  String get errFriendCap => 'Ваш список друзей заполнен.';

  @override
  String get errTargetNotFound => 'Персонажа с таким именем нет.';

  @override
  String get errSelfTarget => 'Нельзя выбрать себя.';

  @override
  String get errInboxFull => 'Почтовый ящик получателя заполнен.';

  @override
  String get errTooLong => 'Слишком длинно.';

  @override
  String get errEmpty => 'Сначала напишите что-нибудь.';

  @override
  String get errNotInChannel => 'Вы не состоите в этой группе.';

  @override
  String get errMuted => 'Вы заглушены и пока не можете писать в чат.';

  @override
  String get errHasAttachments => 'Перед удалением заберите вложения.';

  @override
  String get useItem => 'Использовать';

  @override
  String restoresHp(int amount) {
    return 'Восстанавливает $amount ОЗ';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return 'Использовано: $name · ОЗ $hp/$max';
  }

  @override
  String get errHpFull => 'Ваше здоровье уже полное.';

  @override
  String get errInCombat => 'Сначала завершите охоту или забег в подземелье.';

  @override
  String get errNotUsable => 'Этот предмет нельзя использовать.';

  @override
  String get errTooWounded => 'Вы слишком ранены для боя. Отдохните или выпейте зелье.';

  @override
  String get forgotPassword => 'Забыли пароль?';

  @override
  String get resetHint => 'Введите эл. почту аккаунта, и мы отправим ссылку для смены пароля.';

  @override
  String get resetSent => 'Если аккаунт с такой почтой существует, ссылка для сброса уже в пути. Она действует 1 час.';

  @override
  String get sendLink => 'Отправить ссылку';

  @override
  String get close => 'Закрыть';

  @override
  String get changePassword => 'Сменить пароль';

  @override
  String get currentPassword => 'Текущий пароль';

  @override
  String get newPassword => 'Новый пароль';

  @override
  String get save => 'Сохранить';

  @override
  String get passwordChanged => 'Пароль изменён. На других устройствах выполнен выход.';

  @override
  String get verifyEmailBanner => 'Подтвердите адрес эл. почты.';

  @override
  String get resend => 'Отправить снова';

  @override
  String get verificationSent => 'Письмо для подтверждения отправлено.';

  @override
  String get errEmailNotVerified => 'Подтвердите эл. почту перед созданием персонажа.';

  @override
  String get errAlreadyVerified => 'Ваша эл. почта уже подтверждена.';

  @override
  String get tabParty => 'Отряд';

  @override
  String get tabGuild => 'Гильдия';

  @override
  String get invite => 'Пригласить';

  @override
  String invitedBy(String name) {
    return 'Приглашение от: $name';
  }

  @override
  String get noParty => 'Вы не в отряде. Пригласите кого-нибудь, и отряд соберётся (2–5 участников).';

  @override
  String get createParty => 'Создать отряд';

  @override
  String get pendingInvites => 'Приглашения';

  @override
  String partyMembers(int count, int cap) {
    return 'Отряд ($count/$cap)';
  }

  @override
  String get makeLeader => 'Сделать лидером';

  @override
  String get kick => 'Исключить';

  @override
  String get leave => 'Покинуть';

  @override
  String get disband => 'Распустить';

  @override
  String get leavePartyTitle => 'Покинуть отряд?';

  @override
  String get leavePartyBody => 'Если вы лидер, им станет самый давний участник. Отряд с одним участником распадается.';

  @override
  String get disbandPartyTitle => 'Распустить отряд?';

  @override
  String get disbandPartyBody => 'Все будут исключены из отряда.';

  @override
  String get rankLeader => 'Лидер';

  @override
  String get rankOfficer => 'Офицер';

  @override
  String get rankVeteran => 'Ветеран';

  @override
  String get rankMember => 'Участник';

  @override
  String get rankInitiate => 'Новичок';

  @override
  String get createGuild => 'Основать гильдию';

  @override
  String createGuildBody(String fee) {
    return 'Основание гильдии стоит $fee.';
  }

  @override
  String createGuildFee(String fee) {
    return 'Основать гильдию ($fee)';
  }

  @override
  String get guildName => 'Название гильдии';

  @override
  String get guildTag => 'Тег';

  @override
  String get guildTagHint => '2–5 букв или цифр';

  @override
  String guildMembers(int count, int cap) {
    return 'Участники: $count/$cap';
  }

  @override
  String disbandScheduled(String time) {
    return 'Гильдия распадётся через $time';
  }

  @override
  String get motd => 'Сообщение дня';

  @override
  String get editMotd => 'Изменить сообщение дня';

  @override
  String get setRank => 'Изменить ранг';

  @override
  String setRankFor(String name) {
    return 'Новый ранг: $name';
  }

  @override
  String transferTitle(String name) {
    return 'Сделать $name лидером гильдии?';
  }

  @override
  String get transferBody => 'Вы станете офицером. Самостоятельно отменить это нельзя.';

  @override
  String kickTitle(String name) {
    return 'Исключить $name из гильдии?';
  }

  @override
  String get leaveGuildTitle => 'Покинуть гильдию?';

  @override
  String get leaveGuildBody => 'Чтобы вернуться, понадобится новое приглашение.';

  @override
  String get disbandGuildTitle => 'Распустить гильдию?';

  @override
  String get disbandGuildBody => 'Гильдия распадётся через 48 часов; до этого момента можно отменить.';

  @override
  String get cancelDisband => 'Отменить роспуск';

  @override
  String get errAlreadyInGroup => 'Этот игрок уже в группе.';

  @override
  String get errNotLeader => 'Это может только лидер.';

  @override
  String get errGroupFull => 'Группа заполнена.';

  @override
  String get errInviteInvalid => 'Срок приглашения истёк.';

  @override
  String get errAlreadyInvited => 'Приглашение уже отправлено.';

  @override
  String get errGuildForbidden => 'Вашему рангу это недоступно.';

  @override
  String get errLeaderCannotLeave => 'Сначала передайте лидерство или распустите гильдию.';

  @override
  String get errInvalidGuildName => '3–24 буквы или цифры.';

  @override
  String get errInvalidTag => 'Тег должен содержать 2–5 букв или цифр.';

  @override
  String get errTagTaken => 'Этот тег занят.';

  @override
  String get progress => 'Достижения и мастерство';

  @override
  String get tabAchievements => 'Достижения';

  @override
  String get tabMastery => 'Мастерство';

  @override
  String get tabReputation => 'Репутация';

  @override
  String get tabTitles => 'Титулы';

  @override
  String achievementPoints(int points) {
    return 'Очки достижений: $points';
  }

  @override
  String pointsShort(int points) {
    return '$points оч.';
  }

  @override
  String rewardTitle(String title) {
    return 'Титул: $title';
  }

  @override
  String masteryTier(int tier, int max) {
    return 'Ступень $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return '+$amount сегодня';
  }

  @override
  String get noTitle => 'Без титула';

  @override
  String get noTitlesYet => 'Титулы открываются за достижения.';

  @override
  String get errTitleNotOwned => 'Этот титул вам ещё не открыт.';

  @override
  String get notifications => 'Уведомления';

  @override
  String get markAllRead => 'Отметить все как прочитанные';

  @override
  String get inboxNothing => 'Уведомлений пока нет.';

  @override
  String get notifSettings => 'Настройки уведомлений';

  @override
  String get muteAll => 'Приостановить push-уведомления';

  @override
  String get muteHint => 'На 24 часа. Входящие продолжат пополняться.';

  @override
  String mutedFor(String time) {
    return 'Приостановлено на $time';
  }

  @override
  String get modePushInbox => 'Push + входящие';

  @override
  String get modeInbox => 'Только входящие';

  @override
  String get modeSilent => 'Без звука';

  @override
  String get modeOff => 'Выкл.';

  @override
  String get nkActivity => 'Занятие завершено';

  @override
  String get nkDungeon => 'Забег в подземелье завершён';

  @override
  String get nkTravel => 'Прибытие';

  @override
  String get nkWorldEvent => 'Мировое событие началось';

  @override
  String get nkMarketSold => 'Предмет продан';

  @override
  String get nkMail => 'Новое письмо';

  @override
  String get nkPartyInvite => 'Приглашение в отряд';

  @override
  String get nkGuildInvite => 'Приглашение в гильдию';

  @override
  String get nkAchievement => 'Получено достижение';

  @override
  String get nkSystem => 'Объявление';

  @override
  String timeAgo(String time) {
    return '$time назад';
  }

  @override
  String get segEvents => 'События';

  @override
  String get eventsEmpty => 'Сейчас мировых событий нет.';

  @override
  String get worldBoss => 'Мировой босс';

  @override
  String get invasion => 'Вторжение';

  @override
  String startsIn(String time) {
    return 'Начало через $time';
  }

  @override
  String endsIn(String time) {
    return 'Конец через $time';
  }

  @override
  String regionOnly(String region) {
    return 'Только $region';
  }

  @override
  String get eventGoal => 'Цель';

  @override
  String get enlist => 'Записаться';

  @override
  String get enlistHint =>
      'При записи ваша текущая сила фиксируется. Очки идут за каждую минуту события; награды придут после его окончания.';

  @override
  String enlistedPower(int power) {
    return 'Записаны · сила $power';
  }

  @override
  String myContribution(int score, int tier) {
    return 'Ваши очки $score · уровень награды $tier';
  }

  @override
  String get eventSucceeded => 'Победа';

  @override
  String get eventFailed => 'Провал';

  @override
  String tokenBalance(int tokens) {
    return 'Жетоны события: $tokens';
  }

  @override
  String claimDaily(int tokens) {
    return 'Получить $tokens';
  }

  @override
  String get claimedToday => 'Получено';

  @override
  String eventTokensGained(int tokens) {
    return '+$tokens жетонов события';
  }

  @override
  String get errEventNotOpen => 'Запись на это событие закрыта.';

  @override
  String get errWrongRegion => 'Сначала отправьтесь в регион события.';

  @override
  String get errAlreadyEnlisted => 'Вы уже записаны.';

  @override
  String get errAlreadyClaimed => 'Сегодня уже получено.';

  @override
  String get errEventNotActive => 'Это событие не идёт.';

  @override
  String gatherYield(String activity) {
    return '$activity: добыча';
  }

  @override
  String perHour(String value) {
    return '$value/ч';
  }

  @override
  String get quests => 'Задания';

  @override
  String get questsInProgress => 'В процессе';

  @override
  String get questsNoneOpen => 'Открытых заданий нет. Выберите одно ниже.';

  @override
  String get questsStory => 'Сюжет';

  @override
  String get questsSide => 'Побочные задания';

  @override
  String get questsDaily => 'Ежедневные';

  @override
  String resetsIn(String time) {
    return 'Обновление через $time';
  }

  @override
  String get turnIn => 'Сдать';

  @override
  String questRewarded(String rewards) {
    return 'Задание выполнено: $rewards';
  }

  @override
  String get abandonQuestTitle => 'Бросить задание?';

  @override
  String get abandonQuestBody => 'Прогресс будет потерян. Позже задание можно взять снова.';

  @override
  String get abandonDailyBody => 'Прогресс будет потерян, и сегодня это ежедневное задание больше не взять.';

  @override
  String objGather(String activity) {
    return 'Собрать ($activity)';
  }

  @override
  String get objKillAny => 'Победить врагов';

  @override
  String objKill(String mob) {
    return 'Победить: $mob';
  }

  @override
  String get objBoss => 'Победить босса';

  @override
  String get objDungeonAny => 'Пройти подземелье';

  @override
  String objDungeon(String dungeon) {
    return 'Пройти: $dungeon';
  }

  @override
  String get objCraftAny => 'Создать предметы';

  @override
  String objCraft(String item) {
    return 'Создать: $item';
  }

  @override
  String get evQuestComplete => 'Цели задания выполнены — сдайте его!';

  @override
  String get errQuestLocked => 'Сначала завершите предыдущее сюжетное задание.';

  @override
  String get errQuestTaken => 'Вы уже взяли это задание.';

  @override
  String get errQuestNotOffered => 'Этого ежедневного задания сегодня нет на доске.';

  @override
  String get errQuestCap => 'Слишком много открытых заданий. Завершите или бросьте одно.';

  @override
  String get errQuestNotActive => 'Это задание не открыто.';

  @override
  String get errQuestIncomplete => 'Цели ещё не выполнены.';

  @override
  String get skills => 'Навыки';

  @override
  String get skillActive => 'Активный';

  @override
  String get skillPassive => 'Пассивный';

  @override
  String skillStrike(int power) {
    return 'Урон $power%';
  }

  @override
  String skillHeal(int power) {
    return 'Лечит $power% ОЗ, если здоровья меньше половины';
  }

  @override
  String skillMp(int mp) {
    return '$mp ОМ';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: 'раз в $turns х.', one: 'каждый ход');
    return '$_temp0';
  }

  @override
  String get termsOfService => 'Условия использования';

  @override
  String get privacyPolicy => 'Политика конфиденциальности';

  @override
  String get accountMenu => 'Аккаунт';

  @override
  String get downloadMyData => 'Скачать мои данные';

  @override
  String get exportOpened => 'Ссылка для скачивания открыта в браузере (действует 15 минут).';

  @override
  String get linkOpenFailed => 'Не удалось открыть ссылку.';

  @override
  String get deleteAccount => 'Удалить мой аккаунт';

  @override
  String get deleteAccountBody =>
      'Ваш аккаунт и все персонажи будут безвозвратно удалены через 30 дней. Вход в этот период отменяет удаление. Введите пароль для подтверждения.';

  @override
  String get deleteAccountConfirm => 'Удалить аккаунт';

  @override
  String deleteScheduled(String date) {
    return 'Ваш аккаунт будет удалён $date. Чтобы отменить, просто войдите до этой даты.';
  }

  @override
  String get deletionCancelled => 'Запрос на удаление аккаунта отменён.';

  @override
  String get support => 'Поддержка';

  @override
  String get supportIntro =>
      'Нашли ошибку или нужна помощь с аккаунтом? Напишите нам — команда ответит здесь, и вы получите уведомление.';

  @override
  String get supportNew => 'Новый запрос';

  @override
  String get supportNone => 'Запросов пока нет.';

  @override
  String get supportSent => 'Запрос получен. Мы сообщим, когда появится ответ.';

  @override
  String get supportCategory => 'Тема';

  @override
  String get supportSubject => 'Заголовок';

  @override
  String get supportMessage => 'Что случилось? По возможности опишите по шагам.';

  @override
  String get supportTooShort => 'Заголовок — минимум 3 символа, сообщение — минимум 10.';

  @override
  String get supportBug => 'Ошибка';

  @override
  String get supportAccount => 'Аккаунт';

  @override
  String get supportAbuse => 'Нарушение / жалоба на игрока';

  @override
  String get supportFeedback => 'Предложение';

  @override
  String get supportOther => 'Другое';

  @override
  String get supportOpen => 'Открыт';

  @override
  String get supportAnswered => 'Отвечен';

  @override
  String get supportClosed => 'Закрыт';

  @override
  String get supportReply => 'Ответ команды';

  @override
  String get nkSupport => 'Ответ поддержки';

  @override
  String get errServerBusy => 'Сервер сейчас занят. Повторите через несколько секунд.';

  @override
  String get notifTabInbox => 'Уведомления';

  @override
  String get notifTabRewards => 'Награды';

  @override
  String get rewardsNone =>
      'Наград пока нет. Здесь отображается всё, что вы собрали или выиграли в боях (хранится 30 дней).';

  @override
  String rewardFromActivity(String name) {
    return 'Сбор: $name';
  }

  @override
  String rewardFromHunt(String name) {
    return 'Охота: $name';
  }

  @override
  String rewardFromDungeon(String name) {
    return 'Подземелье: $name';
  }

  @override
  String rewardLevelUp(int level) {
    return 'Уровень $level!';
  }

  @override
  String get levelMax => 'МАКС';

  @override
  String get promotionAvailable => 'ДОСТУПНО ПОВЫШЕНИЕ';

  @override
  String get promotionChoose => 'Выбрать';

  @override
  String promotionConfirm(String name) {
    return 'Стать: $name?';
  }

  @override
  String get promotionBody => 'Этот выбор определит ваш класс навсегда.';

  @override
  String nextMilestone(int level, String what) {
    return 'Далее: ур. $level · $what';
  }

  @override
  String get msPromotion => 'Повышение класса';

  @override
  String msAwakening(String name) {
    return 'Пробуждение: $name';
  }

  @override
  String msUltimate(String name) {
    return 'Ультимейт: $name';
  }

  @override
  String get msCapstone => 'Вершина талантов';

  @override
  String get msMastery => 'Мастерство класса';

  @override
  String msTitle(String name) {
    return 'Титул: $name';
  }

  @override
  String estMinutes(int n) {
    return '~$n мин';
  }

  @override
  String estHours(String h) {
    return '~$h ч';
  }

  @override
  String overflowHeld(int xp) {
    return '$xp опыта ждут вашего выбора';
  }

  @override
  String attributePointsBtn(int n) {
    return 'Характеристики · очков: $n';
  }

  @override
  String talentPointsBtn(int n) {
    return 'Таланты · очков: $n';
  }

  @override
  String get talentsTitle => 'Таланты';

  @override
  String capstoneLocked(int level, int points) {
    return 'Нужны уровень $level и $points очков в этой ветке';
  }

  @override
  String masteryUnlocked(String title) {
    return 'Мастерство класса: $title';
  }

  @override
  String get errPromotionUnavailable => 'Это повышение недоступно.';

  @override
  String get errNotEnoughPoints => 'Недостаточно очков.';

  @override
  String get errTalentMaxed => 'Этот талант уже на максимальном ранге.';

  @override
  String get errCapstoneLocked => 'Для вершины нужны уровень и полная ветка.';

  @override
  String get errSpecRequirement => 'Ваша специализация не может использовать этот предмет.';

  @override
  String get errAttributeRequirement => 'Для этого предмета у вас слишком низкие характеристики.';

  @override
  String get errBusyCombat => 'Нельзя во время боя, AFK или в подземелье.';

  @override
  String skillGuard(int pct, int turns) {
    return '+$pct% брони на $turns х.';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '+$pct% атаки на $turns х.';
  }

  @override
  String get classRanger => 'Следопыт';

  @override
  String get classCleric => 'Жрец';

  @override
  String get classMonk => 'Монах';

  @override
  String get classPaladin => 'Паладин';

  @override
  String get classDruid => 'Друид';

  @override
  String get classBard => 'Бард';

  @override
  String get classShaman => 'Шаман';

  @override
  String get errEquipmentPermission => 'Ваш класс не может использовать этот тип снаряжения.';

  @override
  String get classResource => 'Ресурс';

  @override
  String get classFocus => 'Ключевая характеристика';

  @override
  String classSecondary(String names) {
    return 'Вторичные: $names';
  }

  @override
  String get resStartsEmpty => 'Начинается с нуля';

  @override
  String get resStartsFull => 'Начинается полным';

  @override
  String resStartsAt(int pct) {
    return 'Начинается с $pct%';
  }

  @override
  String resPerHit(int n) {
    return '+$n за удар';
  }

  @override
  String resPerTurn(int n) {
    return '+$n за ход';
  }

  @override
  String baseAttributesHint(int total) {
    return 'Базовые характеристики · 1 очко за 2 уровня, всего $total';
  }

  @override
  String get classEquipment => 'Снаряжение';

  @override
  String get firstSkills => 'Первые навыки';

  @override
  String get classPath => 'Путь класса';

  @override
  String get classPathHint => 'Нажмите на класс, чтобы увидеть его навыки.';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Ур. $awakening Пробуждение · Ур. $ultimate Ультимейт · Ур. $capstone вершина талантов · Ур. $mastery мастерство класса';
  }

  @override
  String get skillTagActive => 'активный';

  @override
  String get skillTagPassive => 'пассивный';

  @override
  String get skillTagAwakening => 'пробуждение';

  @override
  String get skillTagUltimate => 'ультимейт';

  @override
  String skillCostCd(int cost, int cd) {
    return 'Цена $cost · перезарядка $cd х.';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Мастерство на ур. $level: $title';
  }

  @override
  String get classInfoFailed => 'Не удалось загрузить сведения о классе.';

  @override
  String xpShare(int pct) {
    return 'Опыт $pct%';
  }

  @override
  String get dungeonsSolo => 'Подземелья';

  @override
  String get raidsTab => 'Рейды';

  @override
  String raidGroup(int min, int max) {
    return 'Отряд из $min–$max';
  }

  @override
  String raidBossGear(int count) {
    return 'Босс: гарантированно снаряжения $count шт.';
  }

  @override
  String get raidOpen => 'Открыть рейд';

  @override
  String get raidNeedParty =>
      'Рейды рассчитаны на отряд. Встаньте в очередь на случайный отряд, присоединитесь к объявлению на доске или разместите своё.';

  @override
  String get raidLeaderOnly => 'Рейды открывает лидер отряда.';

  @override
  String get raidLobby => 'Лобби рейда';

  @override
  String get raidJoin => 'Присоединиться';

  @override
  String get raidLeave => 'Покинуть';

  @override
  String get raidCancel => 'Отменить';

  @override
  String get raidCancelBody => 'Лобби закроется для всех.';

  @override
  String get raidStart => 'Начать рейд';

  @override
  String get raidStartBody => 'После старта никто не сможет присоединиться или выйти.';

  @override
  String raidMembers(int count, int max) {
    return 'Участники ($count/$max)';
  }

  @override
  String raidNeedMore(int min) {
    return 'Для старта нужно минимум $min участников';
  }

  @override
  String get raidWaitHint => 'Награды откроются по окончании рейда; каждый участник забирает свои.';

  @override
  String get raidDone => 'Рейд окончен — заберите награды';

  @override
  String get raidLocked => 'Пройден на этой неделе';

  @override
  String get raidCollected => 'забрано';

  @override
  String get errRaidNeedsParty => 'Сначала вступите в отряд.';

  @override
  String get errRaidExists => 'У вашего отряда уже есть открытый рейд.';

  @override
  String get errAlreadyInRaid => 'Вы уже в рейде.';

  @override
  String get errRaidNotOpen => 'Этот рейд больше не открыт.';

  @override
  String get errRaidNotEnough => 'Недостаточно участников для старта.';

  @override
  String get errRaidWeeklyLocked => 'Вы уже прошли этот рейд на этой неделе.';

  @override
  String errCharacterBusy(String action) {
    return 'Сначала завершите текущее дело персонажа: $action.';
  }

  @override
  String get busyGathering => 'сбор';

  @override
  String get busyHunt => 'охота';

  @override
  String get busyDungeon => 'подземелье';

  @override
  String get busyRaid => 'рейд';

  @override
  String get busyTravel => 'путешествие';

  @override
  String get busyCrafting => 'создание';

  @override
  String get errCraftActive => 'Задача создания уже выполняется.';

  @override
  String get errCraftTooLong => 'Эта задача заняла бы больше 6 часов.';

  @override
  String craftEach(String time) {
    return '$time за шт.';
  }

  @override
  String craftTotal(String time) {
    return 'Всего $time';
  }

  @override
  String craftRunning(String item, int count) {
    return 'Создание: $item ×$count';
  }

  @override
  String craftProgress(int done, int count) {
    return 'Готово: $done/$count';
  }

  @override
  String get craftDone => 'Всё готово — заберите';

  @override
  String get craftCollectEarlyTitle => 'Забрать сейчас?';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return 'Готовые предметы ($done) будут выданы; материалы и золото за остальные ($rest) вернутся.';
  }

  @override
  String craftCollected(int made, int refunded) {
    return 'Создано: $made, возвращено: $refunded';
  }

  @override
  String get craftBusyHint => 'Создание требует времени; в этот период персонаж ничем другим не занимается.';

  @override
  String get eventHowTitle => 'Как участвовать';

  @override
  String get eventRewardsTitle => 'Награды';

  @override
  String eventStepEnlist(int level) {
    return 'Запись открыта за 24 часа до начала и до конца события (ур. $level+).';
  }

  @override
  String eventStepRegion(String region) {
    return 'При записи нужно находиться в регионе $region.';
  }

  @override
  String eventStepPower(int power) {
    return 'При записи сила фиксируется: атака × 2 + защита + ОЗ ÷ 10. Сейчас у вас: $power.';
  }

  @override
  String eventStepScore(int tier) {
    return 'Очки = сила × минуты участия, с начала события (или с момента записи, если позже). Останавливаются на уровне $tier.';
  }

  @override
  String get eventStepGoal =>
      'Если все вместе достигнут цели, награды идут до высшего уровня; иначе — только до уровня 2.';

  @override
  String get eventStepPaid => 'Награды выдаются автоматически при закрытии. Запись не прерывает другие действия.';

  @override
  String eventTier(int tier, String score) {
    return 'Уровень $tier · от $score очков';
  }

  @override
  String eventTierTokens(int count) {
    return '+$count жетонов события';
  }

  @override
  String eventEstimate(String score, int tier) {
    return 'Если записаться сейчас: около $score очков → уровень $tier';
  }

  @override
  String eventEstimateNone(String score) {
    return 'Если записаться сейчас: около $score очков — на уровень не хватит';
  }

  @override
  String seasonalHow(int count) {
    return 'Пока идёт событие, раз в день можно получить жетоны: $count (день начинается в 00:00 UTC).';
  }

  @override
  String seasonalBenefit(int count) {
    return 'Жетоны остаются на аккаунте (баланс: $count).';
  }

  @override
  String get eventDetails => 'Подробности';

  @override
  String get createCharacterTitle => 'Создание персонажа';

  @override
  String get classTraits => 'Особенности класса';

  @override
  String classResourceLine(String name) {
    return 'Ресурс: $name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return 'Главная характеристика: $primary · вторая: $secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return 'Главная характеристика: $primary';
  }

  @override
  String get classDifficulty => 'Сложность';

  @override
  String get difficultyEasy => 'Лёгкий — прощает ошибки, надёжен на AFK-охоте';

  @override
  String get difficultyMedium => 'Средний — сбалансирован; порядок навыков важен';

  @override
  String get difficultyHard => 'Сложный — хрупкий; нужно подходящее снаряжение';

  @override
  String get skillsTapHint => 'нажмите для подробностей';

  @override
  String get chooseRace => 'Выберите расу';

  @override
  String continueAs(String name) {
    return 'Продолжить: $name';
  }

  @override
  String get nameYourCharacter => 'Назовите персонажа';

  @override
  String get raceTroll => 'Тролль';

  @override
  String get raceUndead => 'Нежить';

  @override
  String get raceDarkElf => 'Тёмный эльф';

  @override
  String get raceDrakari => 'Дракари';

  @override
  String get raceTauren => 'Таурен';

  @override
  String get raceGnome => 'Гном';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '+$atk% атаки и +$def% защиты на $minutes мин';
  }

  @override
  String fedToast(int atk, int def) {
    return 'Сытость: +$atk% атаки, +$def% защиты';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return 'Сытость · +$atk% атаки · +$def% защиты · осталось $time';
  }

  @override
  String huntReach(String regions) {
    return 'Охота отсюда: $regions';
  }

  @override
  String get huntNoMobsHere => 'Здесь не на кого охотиться. Сначала отправьтесь в охотничьи угодья.';

  @override
  String get huntOtherRegions => 'Монстры в других регионах';

  @override
  String get huntTravelHint =>
      'Чтобы охотиться на монстра, отправьтесь в его регион или в регион, связанный с ним дорогой.';

  @override
  String get errMobRegion => 'Этот монстр живёт в другом регионе. Сначала отправьтесь туда.';

  @override
  String restoresHpPct(int pct) {
    return 'Восстанавливает $pct% макс. ОЗ';
  }

  @override
  String get statEvasion => '% уклонения';

  @override
  String get statBlock => '% блока';

  @override
  String get statCritDmg => '% крит. урона';

  @override
  String get statLifesteal => '% вампиризма';

  @override
  String get statDmgReduction => '% снижения урона';

  @override
  String get classJewelry => 'Украшения';

  @override
  String get mapHint =>
      'Можно сразу отправиться в любой регион, доступный по уровню; время зависит от кратчайшей дороги. Двигайте и масштабируйте карту, нажмите на регион.';

  @override
  String get mapNoRoad => 'Прямой дороги сюда нет. Идите через светящиеся соседние регионы.';

  @override
  String get youAreHere => 'Вы здесь.';

  @override
  String travelTime(String time) {
    return 'Время в пути: $time';
  }

  @override
  String xpPerKill(int xp) {
    return '$xp опыта';
  }

  @override
  String get deleteCharacter => 'Удалить персонажа';

  @override
  String deleteCharacterBody(String name) {
    return '$name будет удалён немедленно вместе со всеми предметами и золотом. Отменить это нельзя. Введите имя персонажа для подтверждения.';
  }

  @override
  String get characterNameConfirm => 'Имя персонажа';

  @override
  String characterDeletesAt(String when) {
    return 'Будет удалён: $when';
  }

  @override
  String get cancelCharacterDelete => 'Отменить удаление';

  @override
  String characterDeleted(String name) {
    return '$name: удалён.';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name: удаление отменено.';
  }

  @override
  String get errNameMismatch => 'Имя не совпадает.';

  @override
  String attrGainPerPoint(String gains) {
    return 'За очко: $gains';
  }

  @override
  String attrGainPending(String gains) {
    return 'Прибавится: $gains';
  }

  @override
  String get attrNoEffect => 'Для вашего класса без эффекта';

  @override
  String get attrInputHint => 'Введите, сколько очков получит каждая характеристика';

  @override
  String get gainHp => 'ОЗ';

  @override
  String get gainMp => 'ОМ';

  @override
  String get gainHpRegen => 'ОЗ/с';

  @override
  String get huntSupplies => 'Припасы с собой';

  @override
  String get huntSuppliesHint =>
      'Зелья выпиваются между боями, когда ОЗ ниже половины, еда съедается по окончании бонуса, зелья маны — при нехватке маны. Остаток вернётся в сумку.';

  @override
  String get supplyHp => 'Зелья лечения';

  @override
  String get supplyMp => 'Зелья маны';

  @override
  String get supplyMpHint => 'Для классов с маной, энергией, концентрацией, энергией природы или духом';

  @override
  String get supplyFood => 'Еда';

  @override
  String get huntNoSupplies => 'В сумке нет зелий и еды: купите их в лавке или изготовьте.';

  @override
  String restoresMpPct(int pct) {
    return 'Восполняет $pct% маны';
  }

  @override
  String get huntSupplyOnly => 'Используется только как припас для охоты';

  @override
  String get errBadSupply => 'Этот предмет нельзя взять на охоту.';

  @override
  String get errNotEnoughSupplies => 'В сумке меньше этого припаса.';

  @override
  String get errHuntOnly => 'Возьмите на охоту как припас.';

  @override
  String get huntFled => 'ОЗ упали ниже 5%: вы вышли из боя';

  @override
  String suppliesUsed(String list) {
    return 'Использовано: $list';
  }

  @override
  String get compareTitle => 'Сравнение с надетым';

  @override
  String compareWorn(String name) {
    return 'Надето: $name';
  }

  @override
  String get compareEmpty => 'Этот слот пуст: при надевании добавится всё это.';

  @override
  String get compareColStat => 'Параметр';

  @override
  String get compareColWorn => 'Надето';

  @override
  String get compareColThis => 'Этот';

  @override
  String get compareColDiff => 'Разница';

  @override
  String get compareSame => 'Те же параметры, что и у надетого.';

  @override
  String get wearOk => 'Ваш класс может это носить';

  @override
  String get wearClass => 'Ваш класс не может это носить';

  @override
  String wearLevel(int level) {
    return 'Нужен уровень $level';
  }

  @override
  String get wearAttribute => 'Слишком низкие характеристики';

  @override
  String get wearLegendClass => 'не для вашего класса';

  @override
  String get wearLegendLater => 'пока нет (уровень или характеристики)';

  @override
  String get enchantAction => 'Зачаровать';

  @override
  String get enchantPick => 'Выберите снаряжение для зачарования';

  @override
  String get enchantNone => 'Ничего из снаряжения не подходит к этому свитку.';

  @override
  String enchantFits(String groups) {
    return 'Подходит для: $groups';
  }

  @override
  String get enchantGroupWeapon => 'оружие';

  @override
  String get enchantGroupArmor => 'броня';

  @override
  String get enchantGroupJewelry => 'украшения';

  @override
  String get enchantGroupAny => 'любое снаряжение';

  @override
  String get enchantReplaceTitle => 'Заменить зачарование?';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old будет потеряно, его место займёт $next.';
  }

  @override
  String enchantDone(String item) {
    return 'Зачаровано: $item';
  }

  @override
  String enchantOn(String name) {
    return 'Зачарование: $name';
  }

  @override
  String get wornTag => 'надето';

  @override
  String get errNotEnchant => 'Это не свиток зачарования.';

  @override
  String get errEnchantSlot => 'Это зачарование не подходит к этому снаряжению.';

  @override
  String get statGatherMining => 'Добыча руды: выход';

  @override
  String get statGatherHerbalism => 'Травничество: выход';

  @override
  String get statGatherWoodcutting => 'Лесозаготовка: выход';

  @override
  String get statGatherFishing => 'Рыбалка: выход';

  @override
  String get statGatherHunting => 'Охота: выход';

  @override
  String get statCraftSpeed => 'Скорость создания';

  @override
  String diffXp(int pct) {
    return '$pct% опыта';
  }

  @override
  String get diffNoXp => 'нет опыта';

  @override
  String get diffLegend => 'Цвета рецептов: опыт за создание при вашем навыке';

  @override
  String get recipeGroupEnchantWeapon => 'Зачарования оружия';

  @override
  String get recipeGroupEnchantArmor => 'Зачарования брони';

  @override
  String get recipeGroupEnchantJewelry => 'Зачарования украшений';

  @override
  String get recipeGroupEnchantProfession => 'Зачарования профессий';

  @override
  String get partyHint =>
      'Участники, выполняющие одну и ту же работу одновременно (тот же сбор, профессия, охота или подземелье), получают +3% опыта за каждого работающего: пятеро вместе получают +15%. Никто не отдаёт свой опыт.';

  @override
  String get partyChat => 'Чат отряда';

  @override
  String partyShared(String names, int kept, int earned) {
    return 'Бонус отряда с $names: $earned → $kept опыта';
  }

  @override
  String evRaidMatched(String raid) {
    return 'Отряд для «$raid» найден. См. Рейды.';
  }

  @override
  String get nkRaidMatched => 'Отряд для рейда найден';

  @override
  String get workHunt => 'Охота';

  @override
  String get workDungeon => 'Подземелье';

  @override
  String get raidFind => 'Найти отряд';

  @override
  String get raidQueueJoin => 'Встать в очередь';

  @override
  String get raidQueueLeave => 'Выйти из очереди';

  @override
  String raidQueued(String raid, int waiting) {
    return 'Очередь на $raid · ожидают: $waiting';
  }

  @override
  String get raidQueueHint =>
      'Случайный отряд собирается, когда ждут 5 человек, или через 5 минут, если ждут хотя бы 2.';

  @override
  String get raidAdvertise => 'Разместить объявление';

  @override
  String get raidUnlist => 'Снять объявление';

  @override
  String raidListed(String raid) {
    return 'Ваш отряд на доске: $raid';
  }

  @override
  String get raidBoard => 'Доска объявлений';

  @override
  String get raidBoardEmpty => 'Открытых объявлений нет.';

  @override
  String get raidJoinParty => 'Вступить в отряд';

  @override
  String get raidPickRaid => 'Рейд';

  @override
  String get errListingGone => 'Это объявление исчезло или заполнено.';

  @override
  String get workGathering => 'Сбор';

  @override
  String get workAdventure => 'Приключение';

  @override
  String get workRaid => 'Рейд';

  @override
  String get workTravel => 'В пути';

  @override
  String get workCrafting => 'Создание';

  @override
  String get workIdle => 'Без дела';

  @override
  String get errDungeonRegion => 'Сначала отправьтесь в это место на карте.';

  @override
  String get dungeonsElsewhere =>
      'Новые подземелья и рейды открываются в пути: вход в каждый находится в его регионе на карте.';

  @override
  String get dungeonNoneHere => 'Здесь нет доступных подземелий. Отправьтесь по карте в регион с подземельем.';

  @override
  String dungeonXp(int xp) {
    return '≈ $xp опыта для вас';
  }

  @override
  String appVersionLine(String version) {
    return 'Версия $version';
  }

  @override
  String get updateButton => 'Обновить';

  @override
  String updateAvailable(String version) {
    return 'Доступна новая версия: $version';
  }

  @override
  String get updateLatest => 'У вас последняя версия';

  @override
  String updateConfirm(String version, int size) {
    return 'Скачать и установить версию $version ($size МБ)?';
  }

  @override
  String get updateDownloading => 'Загрузка…';

  @override
  String get updatePermission => 'Разрешите установку из MeNa World и снова нажмите «Обновить».';

  @override
  String get updateFailed => 'Не удалось получить обновление. Проверьте подключение.';

  @override
  String get fightReport => 'Отчёт о бое';

  @override
  String get endTimeUp => 'Время вышло: дошло до запланированного конца.';

  @override
  String get endEarly => 'Забрано досрочно: остановлено до запланированного конца.';

  @override
  String endFled(String mob) {
    return 'ОЗ упали ниже 5% в бою с $mob: вы вышли из боя и вернулись.';
  }

  @override
  String endDefeat(String mob) {
    return 'Вы пали от: $mob.';
  }

  @override
  String endPartyDefeat(String mob) {
    return 'Отряд пал от: $mob.';
  }

  @override
  String endTimeout(String mob) {
    return 'Бой с $mob превысил лимит ходов и был остановлен.';
  }

  @override
  String get endCleared => 'Пройдено: все волны и босс повержены.';

  @override
  String get repCampHint =>
      'При ОЗ ниже 50% охотник разбивает лагерь и восстанавливается перед следующим боем, поэтому до 5% охота доходит редко.';

  @override
  String get repYou => 'Вы';

  @override
  String get repLevel => 'Уровень';

  @override
  String get repHp => 'ОЗ';

  @override
  String get repAttack => 'Атака';

  @override
  String get repDefense => 'Защита';

  @override
  String get repEvasion => 'Уклонение';

  @override
  String get repHitChance => 'Шанс попадания';

  @override
  String get repCrit => 'Крит';

  @override
  String repFights(int n) {
    return 'Боёв: $n';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return 'Попадания: $hits/$swings ($pct%), средн. $avg, макс. $max';
  }

  @override
  String get repYourSwings => 'Ваши удары';

  @override
  String get repPartySwings => 'Удары отряда';

  @override
  String get repEnemySwings => 'Удары врага';

  @override
  String repLowest(int pct) {
    return 'Минимум ОЗ в бою: $pct%';
  }

  @override
  String repRest(String time) {
    return 'В лагере: $time';
  }

  @override
  String get restart => 'Повторить';

  @override
  String get restartTitle => 'Забрать и начать снова?';

  @override
  String get restartBody => 'Полученное уходит в сумку, и то же занятие начнётся снова на то же время.';

  @override
  String get dailyTitle => 'Ежедневная награда за вход';

  @override
  String get dailyBody =>
      'Получайте раз в день на любом персонаже. Каждая неделя подряд платит больше, а 29-й день — больше всего; пропуск начинает цикл заново.';

  @override
  String get dailyClaim => 'Получить';

  @override
  String get dailyTomorrow => 'Возвращайтесь завтра';

  @override
  String dailyDay(int n) {
    return 'День $n';
  }

  @override
  String get dailyOpen => 'Ежедневная награда';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name: заглушён на $min мин за брань.';
  }

  @override
  String get storyTitle => 'Сюжет';

  @override
  String get storyOpen => 'Сюжетный режим';

  @override
  String get storyComingSoon => 'Сюжетный режим скоро появится. Главы будут добавляться сюда по мере написания.';

  @override
  String storyLocked(int level) {
    return 'Откроется на уровне $level';
  }

  @override
  String get bankTitle => 'Банк';

  @override
  String get bankHint =>
      'Общий для всех ваших персонажей. Нажмите на предмет в сумке, чтобы положить его, или на предмет в банке, чтобы забрать. Персональные предметы остаются у своего персонажа.';

  @override
  String get bankEmpty => 'Банк пуст.';

  @override
  String get bankDeposit => 'Положить';

  @override
  String get bankWithdraw => 'Забрать';

  @override
  String get errBankFull => 'Банк заполнен.';

  @override
  String get tabTrade => 'Обмен';

  @override
  String get tradeHint =>
      'Обменивайтесь предметами с другим игроком. Оба выставляют предметы и подтверждают; каждая сторона платит налог 8% от средней рыночной цены отдаваемых предметов, за штуку.';

  @override
  String get tradeWith => 'Имя игрока';

  @override
  String get tradeOpen => 'Начать обмен';

  @override
  String tradeTitle(String name) {
    return 'Обмен с: $name';
  }

  @override
  String get tradeYours => 'Вы отдаёте';

  @override
  String tradeTheirs(String name) {
    return '$name отдаёт';
  }

  @override
  String get tradeNothing => 'Пока ничего.';

  @override
  String get tradeAdd => 'Добавить предмет';

  @override
  String get tradeRemove => 'Забрать обратно';

  @override
  String get tradeTax => 'Ваш налог';

  @override
  String get tradeReady => 'Готов';

  @override
  String get tradeWaiting => 'Не готов';

  @override
  String get tradeConfirm => 'Подтвердить обмен';

  @override
  String get tradeWaitingThem => 'Ждём другого игрока';

  @override
  String get tradeCancel => 'Отменить обмен';

  @override
  String get tradeCancelBody => 'Ничего не переходит из рук в руки.';

  @override
  String get tradeDone => 'Обмен завершён.';

  @override
  String get errTradeBusy => 'Один из вас уже участвует в обмене.';

  @override
  String get errTradeClosed => 'Этот обмен окончен.';

  @override
  String get errTradeItem => 'Этот предмет сейчас нельзя обменять.';

  @override
  String get errTradeEmpty => 'Сначала выставьте что-нибудь.';

  @override
  String get errTradeTooMany => 'Слишком много предметов в обмене.';

  @override
  String get errCraftLevel => 'Уровень персонажа ниже уровня этого предмета.';

  @override
  String get gatherNoneHere => 'Ваши ресурсы растут на других картах: отправьтесь туда, где они указаны.';

  @override
  String get errGatherRegion => 'Этот ресурс растёт на другой карте.';

  @override
  String get helpersOpen => 'Помощники';

  @override
  String get helpersTitle => 'Помощники';

  @override
  String get helpersHint =>
      'Нанимайте помощников для сбора и создания. Они работают вдвое медленнее и получают половину опыта, который идёт в вашу профессию; у каждого вида свой бонус.';

  @override
  String get helpersGather => 'Сбор';

  @override
  String get helpersCraft => 'Создание';

  @override
  String get helperHire => 'Нанять';

  @override
  String helperHireBody(String name, String price) {
    return 'Нанять: $name за $price?';
  }

  @override
  String get helperSend => 'Отправить на работу';

  @override
  String get helperCancel => 'Отозвать';

  @override
  String get helperCancelBody => 'При сборе ничего не вернётся; при создании вернутся материалы и золото.';

  @override
  String get helperDone => 'Готово — заберите';

  @override
  String helperWorking(String time) {
    return 'В работе, осталось $time';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return 'Предметов: $items, +$xp опыта ($profession)';
  }

  @override
  String get helperNoRecipes => 'В этой профессии пока нет рецептов для вашего уровня.';

  @override
  String get helperPickRecipe => 'Выберите рецепт';

  @override
  String helperBonusYield(int pct) {
    return '+$pct% выхода';
  }

  @override
  String helperBonusSpeed(int pct) {
    return 'Быстрее на $pct%';
  }

  @override
  String helperBonusXp(int pct) {
    return '+$pct% опыта';
  }

  @override
  String helperBonusGold(int pct) {
    return 'Меньше золота на $pct%';
  }

  @override
  String helperBonusDouble(int pct) {
    return 'Шанс $pct% на двойной результат';
  }

  @override
  String get helperMiner => 'Рудокоп';

  @override
  String get helperWoodsman => 'Лесоруб';

  @override
  String get helperAngler => 'Рыболов';

  @override
  String get helperHerbalist => 'Травник';

  @override
  String get helperTrapper => 'Траппер';

  @override
  String get helperSmith => 'Кузнец';

  @override
  String get helperTanner => 'Кожевник';

  @override
  String get helperWeaver => 'Ткач';

  @override
  String get helperCarver => 'Резчик';

  @override
  String get helperApothecary => 'Аптекарь';

  @override
  String get helperCook => 'Повар';

  @override
  String get helperSmelter => 'Плавильщик';

  @override
  String get helperScribe => 'Писарь';

  @override
  String get helperGemcutter => 'Ювелир';

  @override
  String get errHelperBusy => 'Этот помощник уже работает.';

  @override
  String get errHelperNotReady => 'Помощник ещё не закончил.';

  @override
  String get shopOpen => 'Магазин';

  @override
  String get shopTitle => 'Магазин кредитов';

  @override
  String get shopCredits => 'Кредиты';

  @override
  String get shopAvatars => 'Аватары';

  @override
  String get shopBackgrounds => 'Фоны';

  @override
  String shopBalance(int n) {
    return 'Кредитов: $n';
  }

  @override
  String get shopBuyCredits => 'Покупайте кредиты, чтобы открывать бонусы и косметику.';

  @override
  String get shopPaymentsSoon => 'Покупка кредитов за деньги скоро появится.';

  @override
  String get shopBoosts => 'Бонусы опыта';

  @override
  String get shopBuy => 'Купить';

  @override
  String shopBuyBody(int n) {
    return 'Потратить кредиты: $n?';
  }

  @override
  String shopPrice(int n) {
    return '$n кр.';
  }

  @override
  String get boostGather => 'Сбор';

  @override
  String get boostAdventure => 'Бой';

  @override
  String get boostCraft => 'Профессии';

  @override
  String get boostDesc => '+100% опыта на 2 часа';

  @override
  String get boostShort => '+100% опыта, 2 ч';

  @override
  String boostLeft(String time) {
    return 'Активно, осталось $time';
  }

  @override
  String boostWatchAd(int n) {
    return 'Смотреть рекламу (осталось: $n)';
  }

  @override
  String get boostAdsSoon => 'Реклама скоро появится';

  @override
  String get errInsufficientCredits => 'Недостаточно кредитов.';

  @override
  String get errBoostActive => 'Этот бонус уже действует.';

  @override
  String get errAdLimit => 'На сегодня бонусы за рекламу закончились.';

  @override
  String helperCanMake(int n) {
    return 'Можно создать: $n';
  }

  @override
  String get helperPickResource => 'Что собирать';

  @override
  String get helperCraftHint => 'Материалы и золото списываются сейчас; вы можете продолжать играть.';

  @override
  String get petsOpen => 'Питомцы';

  @override
  String get petsHint =>
      'Яйца питомцев выпадают при прохождении рейдов выше 400-го уровня. Яйцо вылупляется в инкубаторе; питомец растёт на собственных охотах, а активный питомец даёт вам свой бонус. Он также сражается в ваших боях, добавляя урон к каждому удару.';

  @override
  String get petsEggs => 'Яйца';

  @override
  String get petsNoEggs => 'Яиц нет. Проходите рейды выше 400-го уровня, чтобы найти их.';

  @override
  String petsMine(int n, int max) {
    return 'Мои питомцы ($n/$max)';
  }

  @override
  String get petHatch => 'В инкубатор';

  @override
  String petHatchBody(String egg, String time) {
    return 'Поместить в инкубатор: $egg? Вылупится через $time.';
  }

  @override
  String get petHatching => 'В яйце';

  @override
  String get petHatched => 'Вылупился!';

  @override
  String petPower(int n) {
    return 'Сражается как монстр $n-го уровня';
  }

  @override
  String get petActive => 'Активен';

  @override
  String get petActivate => 'Сделать активным';

  @override
  String get petDeactivate => 'Отдых';

  @override
  String get petHunt => 'Отправить на охоту';

  @override
  String petHunting(String mob) {
    return 'Охотится: $mob';
  }

  @override
  String get petHuntDone => 'Охота окончена';

  @override
  String get petCollectEarly => 'Если закончить охоту сейчас, оплачиваются уже прошедшие бои.';

  @override
  String petCollected(int wins, int xp) {
    return 'Побед: $wins, +$xp опыта питомца';
  }

  @override
  String petLevelUp(String name, int level) {
    return '$name достигает уровня $level!';
  }

  @override
  String get petDied => 'Ваш питомец побеждён; охота закончилась раньше.';

  @override
  String get petPickMob => 'Монстр для охоты';

  @override
  String petXpPerWin(int n) {
    return '+$n опыта за победу';
  }

  @override
  String get errPetCap => 'У вас максимум питомцев.';

  @override
  String get errPetHatching => 'Питомец ещё не вылупился.';

  @override
  String get errPetHunting => 'Питомец уже на охоте.';

  @override
  String petStageLevels(int from, int to) {
    return 'Ур. $from–$to';
  }

  @override
  String get regionMobs => 'Монстры';

  @override
  String get regionDungeons => 'Подземелья';

  @override
  String get regionGather => 'Ресурсы';

  @override
  String get regionEmpty => 'В этом регионе нет монстров, подземелий, ресурсов и заданий.';

  @override
  String get regionRaids => 'Рейды';

  @override
  String get tagUse => 'Использование';

  @override
  String get tagAllClasses => 'Все классы';

  @override
  String get tagDrops => 'Выпадает из';

  @override
  String get tagRaid => 'Рейд';

  @override
  String get tagDungeon => 'Подземелье';

  @override
  String get tagCrafted => 'Создаёт';

  @override
  String tagSkill(int n) {
    return 'навык $n';
  }

  @override
  String get tagSold => 'Продаёт';

  @override
  String get tagGathered => 'Собирает';

  @override
  String get tagUsedIn => 'Материал для';

  @override
  String get petEvolution => 'Эволюция';

  @override
  String get petStrength => 'Подробности о силе';

  @override
  String petAssistRange(int min, int max) {
    return 'Помощь в бою: +$min% – +$max% урона к каждому удару';
  }

  @override
  String petBonusRange(String min, String max) {
    return 'Бонус: $min – $max';
  }

  @override
  String petFightsRange(int min, int max) {
    return 'Охотится как монстр $min–$max-го уровня';
  }

  @override
  String get petFind => 'Где получить';

  @override
  String petAssist(int pct) {
    return 'Помощь в бою: +$pct% урона';
  }

  @override
  String get petLevelNote => 'От 1-го до 100-го уровня';

  @override
  String get petIndex => 'Индекс питомцев';

  @override
  String petIndexCount(int n, int max) {
    return 'Открыто: $n / $max';
  }

  @override
  String slotCount(int n, int max) {
    return '$n / $max';
  }

  @override
  String get petEggTapHint => 'Нажмите на яйцо, чтобы поместить его в инкубатор.';

  @override
  String get petTapHint => 'Нажмите на питомца для просмотра · зелёная точка отмечает активного.';

  @override
  String get petIndexHint => 'Нажмите на питомца, чтобы открыть его страницу.';

  @override
  String get petUnknown => 'Ещё не открыт';

  @override
  String get petAbout => 'О питомце';

  @override
  String get petHabitat => 'Среда обитания';

  @override
  String petElement(String name) {
    return 'Стихия: $name';
  }

  @override
  String get petStage1 => 'Детёныш';

  @override
  String get petStage2 => 'Молодой';

  @override
  String get petStage3 => 'Взрослый';

  @override
  String get bodyMale => 'Мужчина';

  @override
  String get bodyFemale => 'Женщина';

  @override
  String get bodyChoice => 'Внешность';

  @override
  String get bodyChange => 'Изменить внешность';

  @override
  String get portraitNoArmor => 'Без брони';

  @override
  String get combatPower => 'Боевая мощь';

  @override
  String get combatTapHint => 'нажмите для подробностей';

  @override
  String statCap(int cap) {
    return 'предел $cap%';
  }

  @override
  String get cEvasion => 'Уклонение';

  @override
  String get cBlock => 'Блок';

  @override
  String get cCritDmg => 'Крит. урон';

  @override
  String get cDmgReduction => 'Снижение урона';

  @override
  String get cLifesteal => 'Вампиризм';

  @override
  String get cExecute => 'Добивание';

  @override
  String get cMissingHp => 'Последний рубеж';

  @override
  String get cPetAssist => 'Помощь питомца';

  @override
  String get xPower => 'Сила в мировых событиях: 2 × атака + защита + макс. ОЗ / 10.';

  @override
  String get xAttack => 'Базовый урон каждого удара до учёта защиты врага.';

  @override
  String get xDefense => 'Каждый полученный удар ослабляется на половину вашей защиты.';

  @override
  String xCrit(int cap) {
    return 'Шанс критического удара. Не более $cap%.';
  }

  @override
  String get xCritDmg => 'Во столько раз критический удар сильнее обычного.';

  @override
  String get xHit => 'Ваш шанс попасть минус уклонение врага (15–95%).';

  @override
  String xEvasion(int cap) {
    return 'Вычитается из шанса попадания врага. Не более $cap%.';
  }

  @override
  String xBlock(int cap) {
    return 'Шанс уменьшить входящий удар вдвое. Не более $cap%.';
  }

  @override
  String xDmgReduction(int cap) {
    return 'Каждый полученный удар снижается на эту долю. Не более $cap%.';
  }

  @override
  String get xLifesteal => 'Вы исцеляетесь на эту долю нанесённого урона.';

  @override
  String xExecute(int pct) {
    return 'Дополнительный урон по врагам с $pct% ОЗ и ниже.';
  }

  @override
  String get xMissingHp => 'Дополнительный урон до этой величины по мере падения ваших ОЗ.';

  @override
  String get xPetAssist => 'Активный питомец добавляет эту долю к каждому вашему удару.';

  @override
  String get lookFace => 'Лицо';

  @override
  String get lookHairStyle => 'Причёска';

  @override
  String get lookHairColor => 'Цвет волос';

  @override
  String get lookEyeColor => 'Цвет глаз';

  @override
  String get lookSave => 'Сохранить';

  @override
  String get lookHumanOnly => 'Варианты лица и волос для этой расы скоро появятся.';

  @override
  String get raceAbilities => 'Способности расы';

  @override
  String get raceAbilitiesHint => 'все пассивные';

  @override
  String get chooseClass => 'Выберите класс';

  @override
  String raceAttrBonus(int pct, String name) {
    return '+$pct% $name';
  }

  @override
  String raceProfessionXp(int pct, String name) {
    return '+$pct% бонусного опыта: $name';
  }

  @override
  String get language => 'Язык';

  @override
  String get report => 'Пожаловаться';

  @override
  String reportTitle(String name) {
    return 'Пожаловаться на $name';
  }

  @override
  String get reportWhy => 'Что случилось?';

  @override
  String get reportHarassment => 'Оскорбления или травля';

  @override
  String get reportCheating => 'Читы или боты';

  @override
  String get reportName => 'Оскорбительное имя';

  @override
  String get reportSpam => 'Спам или реклама';

  @override
  String get reportScam => 'Мошенничество';

  @override
  String get reportOther => 'Другое';

  @override
  String get reportNote => 'Подробности (необязательно)';

  @override
  String get reportSent => 'Жалоба отправлена. Спасибо.';

  @override
  String get errReportSelf => 'Нельзя пожаловаться на своего персонажа.';

  @override
  String get nkReportReviewed => 'Ваша жалоба рассмотрена';

  @override
  String get guide => 'Справочник';

  @override
  String get guideSearch => 'Поиск по справочнику';

  @override
  String get guideNoMatch => 'Ничего не найдено.';

  @override
  String get guideShowMe => 'Показать';

  @override
  String get next => 'Далее';

  @override
  String get back => 'Назад';

  @override
  String get tutorialSkip => 'Пропустить обучение';

  @override
  String get tutorialStart => 'Начать игру';

  @override
  String get tutorialReplay => 'Пройти обучение ещё раз';

  @override
  String get tutorialWelcomeTitle => 'Добро пожаловать в MeNa World';

  @override
  String get tutorialWelcomeBody =>
      'Персонаж работает, даже когда приложение закрыто. Дайте ему задание, вернитесь позже и заберите добычу.';

  @override
  String get tutorialCharacterBody =>
      'Здесь ваши показатели, очки характеристик и умения. Распределяйте очки за каждый уровень, а на 300 и 700 уровнях выбирайте повышение.';

  @override
  String get tutorialGatherBody =>
      'Выберите ресурс и длительность, затем вернитесь за добычей. Горное дело, рыбалка, лесорубство, травничество и охота развиваются отдельно.';

  @override
  String get tutorialBagBody =>
      'Нажмите на предмет, чтобы надеть, использовать, продать или выставить его. Таблица сравнения покажет, лучше ли он надетого.';

  @override
  String get tutorialAdventureBody =>
      'Охотьтесь на монстров, проходите подземелья, участвуйте в рейдах и событиях, путешествуйте по карте. Монстры и подземелья зависят от региона.';

  @override
  String get tutorialWorkshopBody =>
      'Создавайте предметы с помощью профессий, покупайте в лавке и торгуйте с игроками на рынке.';

  @override
  String get tutorialTopBarTitle => 'Верхняя панель';

  @override
  String get tutorialTopBarBody =>
      'Задания, уведомления и общение (чат, друзья, почта, гильдия, обмен) находятся справа вверху. Персонаж занят одним делом за раз: завершите его или заберите результат, прежде чем начать новое.';

  @override
  String get tutorialDoneTitle => 'Всё готово';

  @override
  String get tutorialDoneBody => 'Забыли, где что находится? Откройте Справочник внизу вкладки «Персонаж».';

  @override
  String get guideCharacterTitle => 'Показатели, характеристики и умения';

  @override
  String get guideCharacterBody =>
      'Вкладка «Персонаж». Очки характеристик распределяются в карточке прогресса. Умения — в отдельном разделе; из умений с одной ролью остаётся самое сильное.';

  @override
  String get guidePromotionTitle => 'Повышения';

  @override
  String get guidePromotionBody =>
      'Вкладка «Персонаж», карточка прогресса. На 300 и 700 уровнях опыт ждёт, пока вы не выберете один из двух путей. Накопленный опыт начисляется после выбора.';

  @override
  String get guideGatherTitle => 'Сбор';

  @override
  String get guideGatherBody =>
      'Вкладка «Сбор». Выберите профессию, доступный ресурс и 1, 3 или 6 часов. Ресурсы привязаны к регионам: для лучших придётся путешествовать.';

  @override
  String get guideBagTitle => 'Сумка и снаряжение';

  @override
  String get guideBagBody =>
      'Вкладка «Сумка». Нажмите на предмет, чтобы увидеть детали и сравнение. Удерживайте и ведите, чтобы выбрать несколько и продать или уничтожить их разом. С полной сумкой новые дела не начать.';

  @override
  String get guideBankTitle => 'Банк аккаунта';

  @override
  String get guideBankBody =>
      'Вкладка «Сумка», кнопка «Банк». Общий для всех ваших персонажей, так можно передавать предметы между ними. Привязанные и надетые предметы не принимаются.';

  @override
  String get guideHuntTitle => 'Охота';

  @override
  String get guideHuntBody =>
      'Вкладка «Приключения», «Охота». Выберите доступного монстра, длительность, зелья и еду. Тяжело раненный персонаж отступает сам.';

  @override
  String get guideDungeonTitle => 'Подземелья и рейды';

  @override
  String get guideDungeonBody =>
      'Вкладка «Приключения», «Подземелье». Подземелья одиночные. Для рейдов нужна группа из 2–5 игроков; «Найти группу» подберёт её. С боссов рейдов падают части комплектов и яйца питомцев.';

  @override
  String get guideTravelTitle => 'Путешествия и карта';

  @override
  String get guideTravelBody =>
      'Вкладка «Приключения», «Путешествие». Нажмите на регион, чтобы увидеть монстров, подземелья, ресурсы и задания, затем «В путь». Требование к уровню есть только у пункта назначения.';

  @override
  String get guideEventsTitle => 'События';

  @override
  String get guideEventsBody =>
      'Вкладка «Приключения», «События». Записывайтесь на мировых боссов и вторжения, получайте ежедневный сезонный жетон. Нажмите на карточку, чтобы узнать правила и награды.';

  @override
  String get guideCraftTitle => 'Профессии и создание предметов';

  @override
  String get guideCraftBody =>
      'Вкладка «Мастерская», «Создание». Откройте профессию, чтобы увидеть рецепты. Цвет показывает, сколько опыта ещё даёт рецепт. Свитки зачарования используются из сумки.';

  @override
  String get guideVendorTitle => 'Лавка и выкуп';

  @override
  String get guideVendorBody =>
      'Вкладка «Мастерская», «Лавка». Продавайте из сумки. Проданное можно выкупить в течение 2 часов.';

  @override
  String get guideMarketTitle => 'Рынок и аукционы';

  @override
  String get guideMarketBody =>
      'Вкладка «Мастерская», «Рынок». Покупайте по фиксированной цене или делайте ставки на аукционе. Выставляйте предметы из сумки: по цене или на аукцион со стартовой ставкой и необязательным выкупом. Ставка в последние минуты продлевает аукцион.';

  @override
  String get guideQuestsTitle => 'Задания';

  @override
  String get guideQuestsBody =>
      'Значок планшета справа вверху. Сюжетные, побочные и ежедневные задания; счётчик показывает готовые к сдаче.';

  @override
  String get guideNotificationsTitle => 'Уведомления и награды';

  @override
  String get guideNotificationsBody =>
      'Значок колокольчика справа вверху. Вкладка «Награды» хранит полученное 30 дней. В настройках выберите, что получать.';

  @override
  String get guideSocialTitle => 'Чат, друзья, почта и обмен';

  @override
  String get guideSocialBody =>
      'Значок диалога справа вверху. Обмен подтверждают обе стороны, он облагается налогом. Блокировка скрывает игрока из чата, друзей и почты.';

  @override
  String get guidePartyTitle => 'Группа';

  @override
  String get guidePartyBody =>
      'Вкладка «Персонаж», раздел «Группа». Приглашайте по имени. Участники за одним типом дел получают бонус к опыту.';

  @override
  String get guideGuildTitle => 'Гильдия, казна и банк гильдии';

  @override
  String get guideGuildBody =>
      'Значок диалога справа вверху, вкладка «Гильдия». Участники вносят золото и предметы; звание определяет, кто и сколько может снимать в день. Все операции записываются.';

  @override
  String get guideRankingTitle => 'Рейтинг';

  @override
  String get guideRankingBody =>
      'Вкладка «Персонаж», кнопка «Рейтинг». Очки складываются из уровня, достижений, мастерства, профессий и пройденных подземелий и рейдов. Есть фильтры по классу и таблице.';

  @override
  String get guidePetsTitle => 'Питомцы';

  @override
  String get guidePetsBody =>
      'Вкладка «Персонаж», кнопка «Питомцы». Яйца падают с высокоуровневых боссов рейдов. Выведите питомца, прокачайте на охоте и активируйте ради бонуса и помощи в бою.';

  @override
  String get guideHelpersTitle => 'Помощники';

  @override
  String get guideHelpersBody =>
      'Вкладка «Персонаж», кнопка «Помощники». Наймите помощника для профессии. Он работает сам, чем бы ни был занят персонаж.';

  @override
  String get guideReportTitle => 'Жалоба на игрока';

  @override
  String get guideReportBody =>
      'Удерживайте строку чата или используйте меню друга. Выберите причину: команда рассмотрит жалобу и сообщит, когда это сделает.';

  @override
  String get guideAccountTitle => 'Аккаунт, язык и поддержка';

  @override
  String get guideAccountBody =>
      'Список персонажей, значок аккаунта вверху. Смените язык или пароль, напишите в поддержку, скачайте данные или удалите аккаунт.';

  @override
  String get ranking => 'Рейтинг';

  @override
  String get rankingScore => 'Очки';

  @override
  String get rankingLevel => 'Уровень';

  @override
  String get rankingAchievements => 'Достижения';

  @override
  String get rankingMastery => 'Мастерство';

  @override
  String get rankingProfessions => 'Профессии';

  @override
  String get rankingClears => 'Прохождения';

  @override
  String get rankingDungeons => 'Подземелья';

  @override
  String get rankingRaids => 'Рейды';

  @override
  String get rankingAllClasses => 'Все классы';

  @override
  String get rankingHowScore => 'Из чего складываются очки';

  @override
  String get rankingScoreNote => 'Каждое подземелье и рейд засчитываются один раз. Обновляется раз в несколько минут.';

  @override
  String get rankingEmpty => 'Таблица ещё подсчитывается.';

  @override
  String get guildVault => 'Казна и банк';

  @override
  String get guildVaultTreasury => 'Казна гильдии';

  @override
  String get guildVaultBank => 'Банк гильдии';

  @override
  String get guildVaultDepositGold => 'Внести золото';

  @override
  String get guildVaultWithdrawGold => 'Снять золото';

  @override
  String get guildVaultGoldAmount => 'Сумма в золоте';

  @override
  String get guildVaultNoLimit => 'Вы можете снимать без дневного лимита.';

  @override
  String guildVaultGoldToday(String used, String limit) {
    return 'Снято сегодня: $used из $limit';
  }

  @override
  String guildVaultItemsToday(int used, int limit) {
    return 'Предметов взято сегодня: $used из $limit';
  }

  @override
  String get guildVaultNoItems => 'Ваше звание пока не даёт доступа к банку гильдии.';

  @override
  String get guildVaultLog => 'Операции';

  @override
  String get guildVaultLogEmpty => 'Операций пока нет.';

  @override
  String guildVaultLogIn(String name, String what) {
    return '$name внёс(ла): $what';
  }

  @override
  String guildVaultLogOut(String name, String what) {
    return '$name снял(а): $what';
  }

  @override
  String get errGuildVaultNotEmpty => 'Перед роспуском опустошите казну и банк.';

  @override
  String get errGuildVaultClosing => 'Гильдия распускается, вносить нельзя.';

  @override
  String get errGuildVaultLimit => 'Это превышает ваш дневной лимит.';

  @override
  String get errItemBound => 'Привязанный предмет остаётся у персонажа.';

  @override
  String get errItemEquipped => 'Сначала снимите предмет.';

  @override
  String get bid => 'Ставка';

  @override
  String get bidPlaced => 'Ставка сделана.';

  @override
  String get marketFixed => 'Фиксированная цена';

  @override
  String get marketAuction => 'Аукцион';

  @override
  String auctionNoBids(String price) {
    return 'Ставок нет, старт $price';
  }

  @override
  String auctionTopBid(String price, int count) {
    return 'Лучшая ставка $price (ставок: $count)';
  }

  @override
  String get auctionYouLead => 'вы лидируете';

  @override
  String get auctionStartBid => 'Начальная ставка за весь лот';

  @override
  String get auctionBuyoutOptional => 'Цена выкупа (необязательно)';

  @override
  String get auctionRules =>
      'Когда время истечёт, побеждает лучшая ставка. Ставка в последние 5 минут продлевает аукцион. После первой ставки отменить его нельзя.';

  @override
  String get auctionYouReceiveAtLeast => 'Вы получите минимум';

  @override
  String auctionMinBid(String price) {
    return 'Минимальная ставка: $price';
  }

  @override
  String get auctionBidCopper => 'Ваша ставка в меди';

  @override
  String get auctionHeldNote => 'Золото удерживается и сразу возвращается, если вашу ставку перебьют.';

  @override
  String auctionBuyoutFor(String price) {
    return 'Выкупить: $price';
  }

  @override
  String get errBidTooLow => 'Ставка ниже минимальной.';

  @override
  String get errTopBidder => 'Ваша ставка уже лучшая.';

  @override
  String get errAuctionBusy => 'Обрабатывается другая ставка. Повторите чуть позже.';

  @override
  String get errAuctionHasBids => 'Аукцион со ставками отменить нельзя.';

  @override
  String get nkMarketOutbid => 'Вашу ставку перебили';

  @override
  String get nkMarketWon => 'Вы выиграли аукцион';

  @override
  String get evMarketOutbid => 'Вашу ставку перебили, золото возвращено.';
}
