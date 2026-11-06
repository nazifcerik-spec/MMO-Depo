// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Japanese (`ja`).
class AppLocalizationsJa extends AppLocalizations {
  AppLocalizationsJa([String locale = 'ja']) : super(locale);

  @override
  String get appTitle => 'MeNa World';

  @override
  String get tagline => '離れている間も冒険は続く。';

  @override
  String get emailOrUsername => 'メールアドレスまたはユーザー名';

  @override
  String get email => 'メールアドレス';

  @override
  String get username => 'ユーザー名';

  @override
  String get password => 'パスワード';

  @override
  String get passwordConfirm => 'パスワード（確認）';

  @override
  String get showPassword => 'パスワードを表示';

  @override
  String get hidePassword => 'パスワードを隠す';

  @override
  String get login => 'ログイン';

  @override
  String get register => 'アカウント作成';

  @override
  String get noAccount => '初めての方はアカウントを作成';

  @override
  String get haveAccount => 'アカウントをお持ちの方はログイン';

  @override
  String get acceptTerms => '13歳以上であり、利用規約とプライバシーポリシーに同意します';

  @override
  String get fieldRequired => '必須項目です';

  @override
  String get invalidEmail => '有効なメールアドレスを入力してください';

  @override
  String get usernameRule => '3～32文字の文字、数字、_ または .';

  @override
  String get passwordMin => '8文字以上';

  @override
  String get passwordMismatch => 'パスワードが一致しません';

  @override
  String get termsRequired => '続行するには規約に同意してください';

  @override
  String get logout => 'ログアウト';

  @override
  String get retry => '再試行';

  @override
  String get offline => 'サーバーに接続できません。接続を確認してください。';

  @override
  String get characters => 'キャラクター';

  @override
  String slotsUsed(int used, int cap) {
    return 'スロット $used/$cap';
  }

  @override
  String get newCharacter => '新規キャラクター';

  @override
  String get noCharacters => 'キャラクターがいません。最初の英雄を作成しましょう。';

  @override
  String get create => '作成';

  @override
  String get characterName => '名前';

  @override
  String get nameRule => '1種類の文字で16文字以内、文字で始めます（数字と _ も使えます）';

  @override
  String get race => '種族';

  @override
  String get characterClass => 'クラス';

  @override
  String get raceHuman => '人間';

  @override
  String get raceElf => 'エルフ';

  @override
  String get raceDwarf => 'ドワーフ';

  @override
  String get raceOrc => 'オーク';

  @override
  String get classWarrior => '戦士';

  @override
  String get classMage => '魔術師';

  @override
  String get classRogue => '盗賊';

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
  String get xp => '経験値';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => '能力値';

  @override
  String get combat => '戦闘';

  @override
  String get statStr => '筋力';

  @override
  String get statAgi => '器用';

  @override
  String get statInt => '知力';

  @override
  String get statVit => '体力';

  @override
  String get statCha => '賢さ';

  @override
  String get statLuck => '運';

  @override
  String get attack => '攻撃力';

  @override
  String get defense => '防御力';

  @override
  String get crit => 'クリティカル';

  @override
  String get hit => '命中';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => 'ゴールド';

  @override
  String get switchCharacter => 'キャラクター変更';

  @override
  String get statePendingDelete => '削除待ち';

  @override
  String get stateFrozen => '凍結中';

  @override
  String get errInvalidCredentials => 'メールアドレス/ユーザー名またはパスワードが違います。';

  @override
  String get errAccountLocked => 'アカウントは一時的にロックされています。';

  @override
  String get errAccountBanned => 'このアカウントは停止されています。';

  @override
  String get errAccountExists => 'そのメールアドレスまたはユーザー名は既に使われています。';

  @override
  String get errRateLimited => '試行回数が多すぎます。しばらくお待ちください。';

  @override
  String get errNameTaken => 'その名前は既に使われています。';

  @override
  String get errSlotCap => 'キャラクタースロットがすべて使用中です。';

  @override
  String get errNotPlayable => 'このキャラクターは現在プレイできません。';

  @override
  String get errSessionExpired => 'セッションが切れました。再度ログインしてください。';

  @override
  String errGeneric(String code) {
    return '問題が発生しました（$code）。';
  }

  @override
  String get navCharacter => 'キャラ';

  @override
  String get navActivity => '採集';

  @override
  String get actMining => '採掘';

  @override
  String get actFishing => '釣り';

  @override
  String get actWoodcutting => '伐採';

  @override
  String get actHerbalism => '薬草採集';

  @override
  String get actHunting => '狩猟';

  @override
  String get activityChoose => '活動を選択';

  @override
  String get durationLabel => '時間';

  @override
  String get gatherResource => '資源';

  @override
  String gatherXp(int xp, int next) {
    return '経験値 $xp / $next';
  }

  @override
  String get gatherMax => '最大レベル';

  @override
  String gatherLocked(int level) {
    return 'レベル$levelで解放';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name 経験値 +$xp';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '$nameがレベル$levelに到達！';
  }

  @override
  String get errGatherLevel => 'この資源を採集するには採集レベルが足りません。';

  @override
  String startActivity(String name) {
    return '$nameを開始';
  }

  @override
  String get collect => '受け取る';

  @override
  String get abandon => '放棄';

  @override
  String get cancel => 'キャンセル';

  @override
  String get soFar => 'これまでの採集量';

  @override
  String get activityDone => '完了 — 受け取れます';

  @override
  String remaining(String time) {
    return '残り$time';
  }

  @override
  String durHm(int h, String m) {
    return '$h時間$m分';
  }

  @override
  String durM(int m) {
    return '$m分';
  }

  @override
  String durS(int s) {
    return '$s秒';
  }

  @override
  String get collectEarlyTitle => '早めに受け取りますか？';

  @override
  String get collectEarlyBody => 'これまでに採集した分を受け取り、活動は終了します。';

  @override
  String get abandonTitle => '活動を放棄しますか？';

  @override
  String get abandonBody => 'これまでの採集分はすべて失われます。';

  @override
  String get lootCollect => 'やった！';

  @override
  String get lootNothing => '今回は何もありませんでした。';

  @override
  String xpGain(int xp) {
    return '経験値 +$xp';
  }

  @override
  String get rarityCommon => 'コモン';

  @override
  String get rarityUncommon => 'アンコモン';

  @override
  String get rarityRare => 'レア';

  @override
  String get rarityEpic => 'エピック';

  @override
  String get rarityLegendary => 'レジェンダリー';

  @override
  String get rarityMythic => 'ミシック';

  @override
  String get errSlotBusy => '別の活動が既に進行中です。';

  @override
  String get errNotActive => 'この活動は既に終了しています。';

  @override
  String get evAchievement => '実績を解除！';

  @override
  String get evMasteryTier => '熟練度の段階が上昇！';

  @override
  String get evMarketSold => '出品アイテムが売れました。';

  @override
  String get evMailReceived => '新しいメールが届きました。';

  @override
  String get evPartyInvite => 'パーティーに招待されました。';

  @override
  String get evGuildInvite => 'ギルドに招待されました。';

  @override
  String get navBag => 'バッグ';

  @override
  String get equipment => '装備';

  @override
  String get bag => 'バッグ';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap kg';
  }

  @override
  String get slotHelm => '頭';

  @override
  String get slotNecklace => '首';

  @override
  String get slotChest => '胴';

  @override
  String get slotGloves => '手';

  @override
  String get slotLegs => '脚';

  @override
  String get slotBoots => '足';

  @override
  String get slotMainhand => 'メイン';

  @override
  String get slotOffhand => 'サブ';

  @override
  String get slotRing => '指輪';

  @override
  String get slotTrinket => 'お守り';

  @override
  String get equip => '装備する';

  @override
  String get unequip => '外す';

  @override
  String get sell => '売却';

  @override
  String get sortLabel => '並べ替え';

  @override
  String get sortSlot => 'バッグの順';

  @override
  String get sortRarity => 'レアリティ';

  @override
  String get sortType => '種類';

  @override
  String get sortName => '名前';

  @override
  String get sortValue => '価値';

  @override
  String selectedCount(int count) {
    return '$count個選択中';
  }

  @override
  String get selectAll => 'すべて選択';

  @override
  String get clearSelection => '選択を解除';

  @override
  String sellSelectedTitle(int count) {
    return '$count個のアイテムを売却しますか？';
  }

  @override
  String destroySelectedTitle(int count) {
    return '$count個のアイテムを破棄しますか？';
  }

  @override
  String get deleteEntry => '削除';

  @override
  String craftHowMany(String name) {
    return '$nameをいくつ作りますか？';
  }

  @override
  String craftCount(int count) {
    return '$count個作成';
  }

  @override
  String get destroy => '破棄';

  @override
  String sellTitle(String name) {
    return '$nameを売却しますか？';
  }

  @override
  String get sellBody => '2時間以内なら商人から買い戻せます。';

  @override
  String destroyTitle(String name) {
    return '$nameを破棄しますか？';
  }

  @override
  String get destroyBody => '元に戻せません。';

  @override
  String get soldFor => '売却額';

  @override
  String get quantity => '数量';

  @override
  String get requiresLevel => '必要レベル';

  @override
  String get weight => '重量';

  @override
  String weightGrams(int g) {
    return '$g g';
  }

  @override
  String get sellValue => '売値';

  @override
  String get bound => 'ソウルバウンド';

  @override
  String get itemLocked => 'このアイテムは取引、メール、または出品に使われています。';

  @override
  String get errBagFull => 'バッグがいっぱいです。';

  @override
  String get errItemLocked => 'このアイテムはロックされています。';

  @override
  String get errClassRestriction => 'このクラスでは使用できません。';

  @override
  String get errLevelRequirement => 'レベルが足りません。';

  @override
  String get errNotEquippable => 'このアイテムは装備できません。';

  @override
  String get navAdventure => '冒険';

  @override
  String get segHunt => '狩猟';

  @override
  String get segDungeon => 'ダンジョン';

  @override
  String get segTravel => '移動';

  @override
  String get huntChoose => '獲物を選択';

  @override
  String get startHunt => '狩猟を開始';

  @override
  String get huntDone => '狩猟終了 — 戦利品を受け取る';

  @override
  String get huntEarlyBody => 'これまでに勝った戦闘の分を受け取り、狩猟は終了します。';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'HP $hp · 攻 $atk · 防 $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '$won勝 · $lost敗';
  }

  @override
  String levelUp(int level) {
    return 'レベルアップ！レベル$levelになりました';
  }

  @override
  String levelReq(int level) {
    return 'Lv.$level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '戦闘$count回 · ボス $boss';
  }

  @override
  String get dungeonEnter => '入場';

  @override
  String get dungeonDone => '攻略完了 — 報酬を受け取る';

  @override
  String get dungeonWaitHint => '報酬は攻略終了後に受け取れます。';

  @override
  String get dungeonCleared => 'クリア！';

  @override
  String get dungeonFailed => 'パーティー全滅';

  @override
  String get bossKilled => 'ボス撃破';

  @override
  String get bonus => 'クリアボーナス';

  @override
  String get youAreIn => '現在地';

  @override
  String get routes => 'ルート';

  @override
  String get noRoutes => 'ここからのルートはありません。';

  @override
  String get travelTo => '出発';

  @override
  String travelling(String place) {
    return '$placeへ移動中';
  }

  @override
  String get arrive => '到着';

  @override
  String get arrived => '到着しました';

  @override
  String get cancelTravel => '引き返す';

  @override
  String get kindCity => '都市';

  @override
  String get kindZone => '地域';

  @override
  String get kindWilderness => '荒野';

  @override
  String get kindDungeon => 'ダンジョン';

  @override
  String get kindContinent => '大陸';

  @override
  String get safeZone => '安全';

  @override
  String get errBusy => '既に進行中のものがあります。';

  @override
  String get errNotReady => 'まだ終わっていません。';

  @override
  String get errWeeklyLocked => '今週はこのダンジョンを既にクリアしています。';

  @override
  String get errNoRoute => 'そこへの直通ルートはありません。';

  @override
  String get navWorkshop => '工房';

  @override
  String get segCraft => '製作';

  @override
  String get segShop => 'ショップ';

  @override
  String get profBlacksmithing => '鍛冶';

  @override
  String get profAlchemy => '錬金術';

  @override
  String get profCooking => '料理';

  @override
  String get learn => '習得';

  @override
  String recipeCount(int count) {
    return 'レシピ$count件';
  }

  @override
  String get itemType => '種類';

  @override
  String get stackSize => 'スタック';

  @override
  String get skillXpReward => 'スキル経験値';

  @override
  String get learnHint => 'この職業を習得するとレシピが解放されます。';

  @override
  String get craft => '製作';

  @override
  String crafted(String name, int count) {
    return '$nameを$count個製作しました';
  }

  @override
  String skillUp(String profession, int level) {
    return '$professionがレベル$levelになりました';
  }

  @override
  String requiresSkill(int level) {
    return 'スキル$levelが必要';
  }

  @override
  String get buy => '購入';

  @override
  String bought(String name) {
    return '$nameを購入しました';
  }

  @override
  String get buybacks => '買い戻し';

  @override
  String get buybackHint => '売却したアイテムは2時間以内なら買い戻せます。';

  @override
  String get errInsufficientGold => 'ゴールドが足りません。';

  @override
  String get errInsufficientMaterials => '素材が足りません。';

  @override
  String get errSkillTooLow => 'このレシピにはスキルが足りません。';

  @override
  String get errProfessionNotLearned => '先に職業を習得してください。';

  @override
  String get errPriceChanged => '価格が変わりました。確認して再度お試しください。';

  @override
  String get errCannotSell => 'このアイテムは売却できません。';

  @override
  String get errBuybackExpired => '買い戻し期間が終了しました。';

  @override
  String get segMarket => 'マーケット';

  @override
  String get classWeapon => '武器';

  @override
  String get classArmor => '防具';

  @override
  String get classConsumable => '消耗品';

  @override
  String get classMaterial => '素材';

  @override
  String get filterAll => 'すべて';

  @override
  String get sortCheapest => '安い順';

  @override
  String get sortNewest => '新着順';

  @override
  String get marketEmpty => '該当する出品はありません。';

  @override
  String listingLine(String seller, String price) {
    return '$seller · 1個 $price';
  }

  @override
  String get loadMore => 'さらに読み込む';

  @override
  String myListings(int count) {
    return '自分の出品（$count）';
  }

  @override
  String get cancelListing => '出品を取り消す';

  @override
  String cancelListingTitle(String name) {
    return '$nameの出品を取り消しますか？';
  }

  @override
  String get cancelListingBody => 'アイテムはバッグに戻ります。出品手数料は最初の6時間以内のみ一部返金されます。';

  @override
  String buyListingTitle(String name) {
    return '$nameを購入しますか？';
  }

  @override
  String buyListingBody(String price, String seller) {
    return '$sellerに$priceを支払います。';
  }

  @override
  String get listOnMarket => 'マーケットに出品';

  @override
  String listed(String name) {
    return '$nameを出品しました';
  }

  @override
  String quantityOf(int count, int max) {
    return '数量：$max個中$count個';
  }

  @override
  String get unitPrice => '単価';

  @override
  String suggestedPrice(String price) {
    return '相場：$price';
  }

  @override
  String get listingFee => '出品手数料（返金不可）';

  @override
  String get youReceive => '売却時の受取額';

  @override
  String get errPriceOutOfRange => 'その価格は許容範囲外です。';

  @override
  String get errListingCap => '有効な出品が多すぎます。';

  @override
  String get errSelfBuy => '自分の出品は購入できません。';

  @override
  String get errListingUnavailable => 'この出品は現在利用できません。';

  @override
  String get errNotTradable => 'このアイテムは取引できません。';

  @override
  String get social => 'ソーシャル';

  @override
  String get tabChat => 'チャット';

  @override
  String get tabFriends => 'フレンド';

  @override
  String get tabMail => 'メール';

  @override
  String get chanGlobal => 'ワールド';

  @override
  String get chanGuild => 'ギルド';

  @override
  String get chanParty => 'パーティー';

  @override
  String get chatEmpty => 'まだメッセージはありません。挨拶してみましょう！';

  @override
  String chatHint(String channel) {
    return '$channelにメッセージ';
  }

  @override
  String get send => '送信';

  @override
  String get addFriend => '追加';

  @override
  String friendAdded(String name) {
    return '$nameさんとフレンドになりました';
  }

  @override
  String friendRequested(String name) {
    return '$nameさんにフレンド申請を送りました';
  }

  @override
  String get friendRequests => 'フレンド申請';

  @override
  String friendsCount(int count) {
    return 'フレンド（$count）';
  }

  @override
  String get noFriends => 'フレンドはまだいません。名前で追加しましょう。';

  @override
  String get more => 'その他';

  @override
  String get whisper => 'ささやき';

  @override
  String get remove => '削除';

  @override
  String get block => 'ブロック';

  @override
  String get unblock => 'ブロック解除';

  @override
  String get blockedList => 'ブロックリスト';

  @override
  String removeFriendTitle(String name) {
    return '$nameさんをフレンドから削除しますか？';
  }

  @override
  String get removeFriendBody => '後でまたフレンド申請を送れます。';

  @override
  String blockTitle(String name) {
    return '$nameさんをブロックしますか？';
  }

  @override
  String get blockBody => '相手はささやきを送れず、チャットは非表示になり、メールは届かなくなります。フレンド関係も解消されます。';

  @override
  String get accept => '承諾';

  @override
  String get decline => '辞退';

  @override
  String get newMail => 'メールを書く';

  @override
  String get inboxEmpty => 'メールボックスは空です。';

  @override
  String get noSubject => '（件名なし）';

  @override
  String mailFrom(String name) {
    return '差出人：$name';
  }

  @override
  String get attachments => '添付';

  @override
  String get claim => '添付を受け取る';

  @override
  String get mailClaimed => '添付物をバッグと財布に移しました';

  @override
  String get deleteMailTitle => 'このメールを削除しますか？';

  @override
  String get deleteMailBody => '削除したメールは復元できません。';

  @override
  String mailSent(String postage) {
    return 'メールを送信しました（送料 $postage）';
  }

  @override
  String get mailTo => '宛先（キャラクター名）';

  @override
  String get mailSubject => '件名';

  @override
  String get mailBody => '本文';

  @override
  String get attachGold => 'ゴールドを添付';

  @override
  String get errAlreadyFriends => '既にフレンドです。';

  @override
  String get errAlreadyRequested => 'フレンド申請は既に保留中です。';

  @override
  String get errBlocked => 'このプレイヤーとはやり取りできません。';

  @override
  String get errFriendCap => 'フレンドリストがいっぱいです。';

  @override
  String get errTargetNotFound => 'その名前のキャラクターはいません。';

  @override
  String get errSelfTarget => '自分を対象にはできません。';

  @override
  String get errInboxFull => '相手のメールボックスがいっぱいです。';

  @override
  String get errTooLong => '長すぎます。';

  @override
  String get errEmpty => '先に入力してください。';

  @override
  String get errNotInChannel => 'そのグループに参加していません。';

  @override
  String get errMuted => 'チャット制限中のため、現在チャットできません。';

  @override
  String get errHasAttachments => '削除する前に添付を受け取ってください。';

  @override
  String get useItem => '使用';

  @override
  String restoresHp(int amount) {
    return 'HPを$amount回復';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return '$nameを使用 · HP $hp/$max';
  }

  @override
  String get errHpFull => 'HPは既に満タンです。';

  @override
  String get errInCombat => '先に狩猟かダンジョン攻略を終えてください。';

  @override
  String get errNotUsable => 'このアイテムは使用できません。';

  @override
  String get errTooWounded => '傷が深く戦えません。休むかポーションを飲んでください。';

  @override
  String get forgotPassword => 'パスワードをお忘れですか？';

  @override
  String get resetHint => 'アカウントのメールアドレスを入力すると、新しいパスワードを設定するリンクを送信します。';

  @override
  String get resetSent => 'そのメールアドレスのアカウントがあれば、再設定リンクを送信しました。リンクの有効期限は1時間です。';

  @override
  String get sendLink => 'リンクを送信';

  @override
  String get close => '閉じる';

  @override
  String get changePassword => 'パスワード変更';

  @override
  String get currentPassword => '現在のパスワード';

  @override
  String get newPassword => '新しいパスワード';

  @override
  String get save => '保存';

  @override
  String get passwordChanged => 'パスワードを変更しました。他の端末はログアウトされました。';

  @override
  String get verifyEmailBanner => 'メールアドレスを確認してください。';

  @override
  String get resend => '再送信';

  @override
  String get verificationSent => '確認メールを送信しました。';

  @override
  String get errEmailNotVerified => 'キャラクターを作成する前にメールアドレスを確認してください。';

  @override
  String get errAlreadyVerified => 'メールアドレスは既に確認済みです。';

  @override
  String get tabParty => 'パーティー';

  @override
  String get tabGuild => 'ギルド';

  @override
  String get invite => '招待';

  @override
  String invitedBy(String name) {
    return '$nameさんからの招待';
  }

  @override
  String get noParty => 'パーティーに参加していません。誰かを招待するとパーティーが組まれます（2～5人）。';

  @override
  String get createParty => 'パーティーを作成';

  @override
  String get pendingInvites => '招待';

  @override
  String partyMembers(int count, int cap) {
    return 'パーティー（$count/$cap）';
  }

  @override
  String get makeLeader => 'リーダーにする';

  @override
  String get kick => '追放';

  @override
  String get leave => '脱退';

  @override
  String get disband => '解散';

  @override
  String get leavePartyTitle => 'パーティーを抜けますか？';

  @override
  String get leavePartyBody => 'リーダーが抜けると、最も古くからのメンバーが引き継ぎます。1人になるとパーティーは解散します。';

  @override
  String get disbandPartyTitle => 'パーティーを解散しますか？';

  @override
  String get disbandPartyBody => '全員がパーティーから外れます。';

  @override
  String get rankLeader => 'リーダー';

  @override
  String get rankOfficer => 'オフィサー';

  @override
  String get rankVeteran => 'ベテラン';

  @override
  String get rankMember => 'メンバー';

  @override
  String get rankInitiate => '新入り';

  @override
  String get createGuild => 'ギルドを設立';

  @override
  String createGuildBody(String fee) {
    return 'ギルドの設立には$feeかかります。';
  }

  @override
  String createGuildFee(String fee) {
    return 'ギルドを設立（$fee）';
  }

  @override
  String get guildName => 'ギルド名';

  @override
  String get guildTag => 'タグ';

  @override
  String get guildTagHint => '文字または数字2～5字';

  @override
  String guildMembers(int count, int cap) {
    return 'メンバー $count/$cap';
  }

  @override
  String disbandScheduled(String time) {
    return '$time後にギルド解散';
  }

  @override
  String get motd => '本日のメッセージ';

  @override
  String get editMotd => '本日のメッセージを編集';

  @override
  String get setRank => 'ランク変更';

  @override
  String setRankFor(String name) {
    return '$nameさんの新しいランク';
  }

  @override
  String transferTitle(String name) {
    return '$nameさんをギルドリーダーにしますか？';
  }

  @override
  String get transferBody => 'あなたはオフィサーになり、自分では元に戻せません。';

  @override
  String kickTitle(String name) {
    return '$nameさんをギルドから追放しますか？';
  }

  @override
  String get leaveGuildTitle => 'ギルドを脱退しますか？';

  @override
  String get leaveGuildBody => '再加入には新しい招待が必要です。';

  @override
  String get disbandGuildTitle => 'ギルドを解散しますか？';

  @override
  String get disbandGuildBody => 'ギルドは48時間後に解散します。それまでは取り消せます。';

  @override
  String get cancelDisband => '解散を取り消す';

  @override
  String get errAlreadyInGroup => '既にグループに参加しています。';

  @override
  String get errNotLeader => 'リーダーのみ可能です。';

  @override
  String get errGroupFull => 'グループがいっぱいです。';

  @override
  String get errInviteInvalid => '招待の期限が切れています。';

  @override
  String get errAlreadyInvited => '招待は既に保留中です。';

  @override
  String get errGuildForbidden => 'あなたのランクではできません。';

  @override
  String get errLeaderCannotLeave => '先にリーダーを譲るかギルドを解散してください。';

  @override
  String get errInvalidGuildName => '文字または数字3～24字。';

  @override
  String get errInvalidTag => 'タグは文字または数字2～5字にしてください。';

  @override
  String get errTagTaken => 'そのタグは使用中です。';

  @override
  String get progress => '実績と熟練度';

  @override
  String get tabAchievements => '実績';

  @override
  String get tabMastery => '熟練度';

  @override
  String get tabReputation => '名声';

  @override
  String get tabTitles => '称号';

  @override
  String achievementPoints(int points) {
    return '実績ポイント $points';
  }

  @override
  String pointsShort(int points) {
    return '${points}pt';
  }

  @override
  String rewardTitle(String title) {
    return '称号：$title';
  }

  @override
  String masteryTier(int tier, int max) {
    return '段階 $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return '今日 +$amount';
  }

  @override
  String get noTitle => '称号なし';

  @override
  String get noTitlesYet => '称号は実績で解放されます。';

  @override
  String get errTitleNotOwned => 'その称号はまだ解放していません。';

  @override
  String get notifications => '通知';

  @override
  String get markAllRead => 'すべて既読にする';

  @override
  String get inboxNothing => '通知はまだありません。';

  @override
  String get notifSettings => '通知設定';

  @override
  String get muteAll => 'プッシュ通知を一時停止';

  @override
  String get muteHint => '24時間。受信箱には引き続き届きます。';

  @override
  String mutedFor(String time) {
    return '$timeの間一時停止中';
  }

  @override
  String get modePushInbox => 'プッシュ + 受信箱';

  @override
  String get modeInbox => '受信箱のみ';

  @override
  String get modeSilent => 'サイレント';

  @override
  String get modeOff => 'オフ';

  @override
  String get nkActivity => '活動完了';

  @override
  String get nkDungeon => 'ダンジョン攻略完了';

  @override
  String get nkTravel => '到着';

  @override
  String get nkWorldEvent => 'ワールドイベント開始';

  @override
  String get nkMarketSold => 'アイテム売却';

  @override
  String get nkMail => '新着メール';

  @override
  String get nkPartyInvite => 'パーティー招待';

  @override
  String get nkGuildInvite => 'ギルド招待';

  @override
  String get nkAchievement => '実績解除';

  @override
  String get nkSystem => 'お知らせ';

  @override
  String timeAgo(String time) {
    return '$time前';
  }

  @override
  String get segEvents => 'イベント';

  @override
  String get eventsEmpty => '予定されているワールドイベントはありません。';

  @override
  String get worldBoss => 'ワールドボス';

  @override
  String get invasion => '侵攻';

  @override
  String startsIn(String time) {
    return '$time後に開始';
  }

  @override
  String endsIn(String time) {
    return '$time後に終了';
  }

  @override
  String regionOnly(String region) {
    return '$region限定';
  }

  @override
  String get eventGoal => '目標';

  @override
  String get enlist => '参加';

  @override
  String get enlistHint => '参加時に現在の戦闘力が固定されます。イベント中は1分ごとにスコアを獲得し、報酬は終了時に届きます。';

  @override
  String enlistedPower(int power) {
    return '参加中 · 戦闘力 $power';
  }

  @override
  String myContribution(int score, int tier) {
    return 'スコア $score · 報酬ランク $tier';
  }

  @override
  String get eventSucceeded => '勝利';

  @override
  String get eventFailed => '失敗';

  @override
  String tokenBalance(int tokens) {
    return 'イベントトークン $tokens';
  }

  @override
  String claimDaily(int tokens) {
    return '$tokensを受け取る';
  }

  @override
  String get claimedToday => '受取済み';

  @override
  String eventTokensGained(int tokens) {
    return 'イベントトークン +$tokens';
  }

  @override
  String get errEventNotOpen => 'このイベントの参加受付は終了しました。';

  @override
  String get errWrongRegion => '先にイベントの地域へ移動してください。';

  @override
  String get errAlreadyEnlisted => '既に参加しています。';

  @override
  String get errAlreadyClaimed => '今日は既に受け取りました。';

  @override
  String get errEventNotActive => 'このイベントは開催されていません。';

  @override
  String gatherYield(String activity) {
    return '$activityの産出量';
  }

  @override
  String perHour(String value) {
    return '$value/時';
  }

  @override
  String get quests => 'クエスト';

  @override
  String get questsInProgress => '進行中';

  @override
  String get questsNoneOpen => '進行中のクエストはありません。下から選んでください。';

  @override
  String get questsStory => 'ストーリー';

  @override
  String get questsSide => 'サブクエスト';

  @override
  String get questsDaily => 'デイリー';

  @override
  String resetsIn(String time) {
    return '$time後にリセット';
  }

  @override
  String get turnIn => '報告';

  @override
  String questRewarded(String rewards) {
    return 'クエスト完了：$rewards';
  }

  @override
  String get abandonQuestTitle => 'クエストを放棄しますか？';

  @override
  String get abandonQuestBody => '進行状況は失われます。後でまた受けられます。';

  @override
  String get abandonDailyBody => '進行状況は失われ、このデイリーは今日はもう受けられません。';

  @override
  String objGather(String activity) {
    return '採集（$activity）';
  }

  @override
  String get objKillAny => '敵を倒す';

  @override
  String objKill(String mob) {
    return '$mobを倒す';
  }

  @override
  String get objBoss => 'ボスを倒す';

  @override
  String get objDungeonAny => 'ダンジョンをクリア';

  @override
  String objDungeon(String dungeon) {
    return '$dungeonをクリア';
  }

  @override
  String get objCraftAny => 'アイテムを製作';

  @override
  String objCraft(String item) {
    return '$itemを製作';
  }

  @override
  String get evQuestComplete => 'クエスト目標を達成 — 報告しましょう！';

  @override
  String get errQuestLocked => '先に前のストーリークエストを終えてください。';

  @override
  String get errQuestTaken => 'このクエストは既に受けています。';

  @override
  String get errQuestNotOffered => 'そのデイリーは今日の掲示板にありません。';

  @override
  String get errQuestCap => '進行中のクエストが多すぎます。完了するか放棄してください。';

  @override
  String get errQuestNotActive => 'このクエストは進行中ではありません。';

  @override
  String get errQuestIncomplete => '目標はまだ達成していません。';

  @override
  String get skills => 'スキル';

  @override
  String get skillActive => 'アクティブ';

  @override
  String get skillPassive => 'パッシブ';

  @override
  String skillStrike(int power) {
    return 'ダメージ $power%';
  }

  @override
  String skillHeal(int power) {
    return 'HPが半分以下でHPを$power%回復';
  }

  @override
  String skillMp(int mp) {
    return 'MP $mp';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: '$turnsターンごと', one: '毎ターン');
    return '$_temp0';
  }

  @override
  String get termsOfService => '利用規約';

  @override
  String get privacyPolicy => 'プライバシーポリシー';

  @override
  String get accountMenu => 'アカウント';

  @override
  String get downloadMyData => 'データをダウンロード';

  @override
  String get exportOpened => 'ブラウザでダウンロードリンクを開きました（15分間有効）。';

  @override
  String get linkOpenFailed => 'リンクを開けませんでした。';

  @override
  String get deleteAccount => 'アカウントを削除';

  @override
  String get deleteAccountBody => 'アカウントとすべてのキャラクターは30日後に完全に削除されます。その間にログインすると削除は取り消されます。確認のためパスワードを入力してください。';

  @override
  String get deleteAccountConfirm => 'アカウントを削除';

  @override
  String deleteScheduled(String date) {
    return 'アカウントは$dateに削除されます。取り消すには、それまでにログインしてください。';
  }

  @override
  String get deletionCancelled => 'アカウント削除の申請を取り消しました。';

  @override
  String get support => 'サポート';

  @override
  String get supportIntro => '不具合を見つけた、またはアカウントについて困っていますか？ご連絡いただければチームがここで回答し、通知でお知らせします。';

  @override
  String get supportNew => '新しい問い合わせ';

  @override
  String get supportNone => '問い合わせはまだありません。';

  @override
  String get supportSent => '問い合わせを受け付けました。回答があればお知らせします。';

  @override
  String get supportCategory => 'トピック';

  @override
  String get supportSubject => '件名';

  @override
  String get supportMessage => '何が起きましたか？できれば順を追って書いてください。';

  @override
  String get supportTooShort => '件名は3文字以上、本文は10文字以上必要です。';

  @override
  String get supportBug => '不具合';

  @override
  String get supportAccount => 'アカウント';

  @override
  String get supportAbuse => '迷惑行為 / プレイヤー報告';

  @override
  String get supportFeedback => '提案';

  @override
  String get supportOther => 'その他';

  @override
  String get supportOpen => '受付中';

  @override
  String get supportAnswered => '回答済み';

  @override
  String get supportClosed => '終了';

  @override
  String get supportReply => 'チームの回答';

  @override
  String get nkSupport => 'サポートの回答';

  @override
  String get errServerBusy => 'サーバーが混み合っています。数秒後に再度お試しください。';

  @override
  String get notifTabInbox => '通知';

  @override
  String get notifTabRewards => '報酬';

  @override
  String get rewardsNone => '報酬はまだありません。採集や戦闘で得たものはすべてここに表示されます（30日間保存）。';

  @override
  String rewardFromActivity(String name) {
    return '採集：$name';
  }

  @override
  String rewardFromHunt(String name) {
    return '狩猟：$name';
  }

  @override
  String rewardFromDungeon(String name) {
    return 'ダンジョン：$name';
  }

  @override
  String rewardLevelUp(int level) {
    return 'レベル$level！';
  }

  @override
  String get levelMax => 'MAX';

  @override
  String get promotionAvailable => '昇格可能';

  @override
  String get promotionChoose => '選択';

  @override
  String promotionConfirm(String name) {
    return '$nameになりますか？';
  }

  @override
  String get promotionBody => 'この選択でクラスが永久に決まります。';

  @override
  String nextMilestone(int level, String what) {
    return '次：Lv.$level · $what';
  }

  @override
  String get msPromotion => 'クラス昇格';

  @override
  String msAwakening(String name) {
    return '覚醒：$name';
  }

  @override
  String msUltimate(String name) {
    return '奥義：$name';
  }

  @override
  String get msCapstone => 'タレントの頂点';

  @override
  String get msMastery => 'クラスマスタリー';

  @override
  String msTitle(String name) {
    return '称号：$name';
  }

  @override
  String estMinutes(int n) {
    return '約$n分';
  }

  @override
  String estHours(String h) {
    return '約$h時間';
  }

  @override
  String overflowHeld(int xp) {
    return '$xp経験値が選択を待っています';
  }

  @override
  String attributePointsBtn(int n) {
    return '能力値 · $nポイント';
  }

  @override
  String talentPointsBtn(int n) {
    return 'タレント · $nポイント';
  }

  @override
  String get talentsTitle => 'タレント';

  @override
  String capstoneLocked(int level, int points) {
    return 'レベル$levelとこのツリーに$pointsポイントが必要';
  }

  @override
  String masteryUnlocked(String title) {
    return 'クラスマスタリー：$title';
  }

  @override
  String get errPromotionUnavailable => 'その昇格はできません。';

  @override
  String get errNotEnoughPoints => 'ポイントが足りません。';

  @override
  String get errTalentMaxed => 'このタレントは既に最大ランクです。';

  @override
  String get errCapstoneLocked => '頂点にはレベルと完成したツリーが必要です。';

  @override
  String get errSpecRequirement => 'あなたの専門では使用できません。';

  @override
  String get errAttributeRequirement => 'このアイテムには能力値が足りません。';

  @override
  String get errBusyCombat => '戦闘、放置活動、ダンジョン中はできません。';

  @override
  String skillGuard(int pct, int turns) {
    return '$turnsターンの間 防具 +$pct%';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '$turnsターンの間 攻撃力 +$pct%';
  }

  @override
  String get classRanger => 'レンジャー';

  @override
  String get classCleric => 'クレリック';

  @override
  String get classMonk => 'モンク';

  @override
  String get classPaladin => 'パラディン';

  @override
  String get classDruid => 'ドルイド';

  @override
  String get classBard => '吟遊詩人';

  @override
  String get classShaman => 'シャーマン';

  @override
  String get errEquipmentPermission => 'このクラスはこの装備種類を使用できません。';

  @override
  String get classResource => 'リソース';

  @override
  String get classFocus => '主要能力値';

  @override
  String classSecondary(String names) {
    return '副：$names';
  }

  @override
  String get resStartsEmpty => '0から開始';

  @override
  String get resStartsFull => '満タンで開始';

  @override
  String resStartsAt(int pct) {
    return '$pct%から開始';
  }

  @override
  String resPerHit(int n) {
    return '命中ごとに +$n';
  }

  @override
  String resPerTurn(int n) {
    return 'ターンごとに +$n';
  }

  @override
  String baseAttributesHint(int total) {
    return '基本能力値 · 2レベルごとに1ポイント、合計$total';
  }

  @override
  String get classEquipment => '装備';

  @override
  String get firstSkills => '初期スキル';

  @override
  String get classPath => 'クラスの道';

  @override
  String get classPathHint => 'クラスをタップするとスキルが見られます。';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Lv.$awakening 覚醒 · Lv.$ultimate 奥義 · Lv.$capstone タレントの頂点 · Lv.$mastery クラスマスタリー';
  }

  @override
  String get skillTagActive => 'アクティブ';

  @override
  String get skillTagPassive => 'パッシブ';

  @override
  String get skillTagAwakening => '覚醒';

  @override
  String get skillTagUltimate => '奥義';

  @override
  String skillCostCd(int cost, int cd) {
    return 'コスト $cost · 再使用 $cdターン';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Lv.$levelマスタリー：$title';
  }

  @override
  String get classInfoFailed => 'クラス情報を読み込めませんでした。';

  @override
  String xpShare(int pct) {
    return '経験値 $pct%';
  }

  @override
  String get dungeonsSolo => 'ダンジョン';

  @override
  String get raidsTab => 'レイド';

  @override
  String raidGroup(int min, int max) {
    return '$min～$max人パーティー';
  }

  @override
  String raidBossGear(int count) {
    return 'ボス：装備$count個確定';
  }

  @override
  String get raidOpen => 'レイドを開く';

  @override
  String get raidNeedParty => 'レイドはパーティー向けです。下のキューでランダムパーティーに入るか、掲示板の募集に参加するか、自分で募集してください。';

  @override
  String get raidLeaderOnly => 'レイドはパーティーリーダーが開きます。';

  @override
  String get raidLobby => 'レイドロビー';

  @override
  String get raidJoin => '参加';

  @override
  String get raidLeave => '退出';

  @override
  String get raidCancel => '中止';

  @override
  String get raidCancelBody => '全員のロビーが閉じます。';

  @override
  String get raidStart => 'レイド開始';

  @override
  String get raidStartBody => '開始後は誰も参加・退出できません。';

  @override
  String raidMembers(int count, int max) {
    return 'メンバー（$count/$max）';
  }

  @override
  String raidNeedMore(int min) {
    return '開始には$min人以上必要';
  }

  @override
  String get raidWaitHint => '報酬はレイド終了後に開放され、各メンバーが自分の分を受け取ります。';

  @override
  String get raidDone => 'レイド終了 — 報酬を受け取る';

  @override
  String get raidLocked => '今週クリア済み';

  @override
  String get raidCollected => '受取済み';

  @override
  String get errRaidNeedsParty => '先にパーティーに参加してください。';

  @override
  String get errRaidExists => 'パーティーには既に開いているレイドがあります。';

  @override
  String get errAlreadyInRaid => '既にレイドに参加しています。';

  @override
  String get errRaidNotOpen => 'このレイドは既に開いていません。';

  @override
  String get errRaidNotEnough => '開始に必要なメンバーが足りません。';

  @override
  String get errRaidWeeklyLocked => '今週はこのレイドを既にクリアしています。';

  @override
  String errCharacterBusy(String action) {
    return '先にキャラクターの作業を終えてください：$action。';
  }

  @override
  String get busyGathering => '採集';

  @override
  String get busyHunt => '狩猟';

  @override
  String get busyDungeon => 'ダンジョン';

  @override
  String get busyRaid => 'レイド';

  @override
  String get busyTravel => '移動';

  @override
  String get busyCrafting => '製作';

  @override
  String get errCraftActive => '製作作業が既に進行中です。';

  @override
  String get errCraftTooLong => 'この作業は6時間を超えます。';

  @override
  String craftEach(String time) {
    return '1個 $time';
  }

  @override
  String craftTotal(String time) {
    return '合計 $time';
  }

  @override
  String craftRunning(String item, int count) {
    return '$item×$countを製作中';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count完了';
  }

  @override
  String get craftDone => 'すべて完了 — 受け取る';

  @override
  String get craftCollectEarlyTitle => '今すぐ受け取りますか？';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return '完成した$done個を受け取り、残り$rest個分の素材とゴールドは返還されます。';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '$made個製作、$refunded個返還';
  }

  @override
  String get craftBusyHint => '製作には時間がかかり、その間キャラクターは他のことをしません。';

  @override
  String get eventHowTitle => '参加方法';

  @override
  String get eventRewardsTitle => '報酬';

  @override
  String eventStepEnlist(int level) {
    return '開始24時間前から終了まで参加できます（Lv.$level以上）。';
  }

  @override
  String eventStepRegion(String region) {
    return '参加時に$regionにいる必要があります。';
  }

  @override
  String eventStepPower(int power) {
    return '参加時に戦闘力が固定されます：攻撃力 × 2 + 防御力 + HP ÷ 10。現在：$power。';
  }

  @override
  String eventStepScore(int tier) {
    return 'スコア = 戦闘力 × 参加分数。開始時（遅れて参加した場合はその時）から計算し、ランク$tierで止まります。';
  }

  @override
  String get eventStepGoal => '全員で目標を達成すれば最高ランクまで報酬を得られ、未達ならランク2までです。';

  @override
  String get eventStepPaid => '報酬は終了時に自動で支払われます。参加しても他の行動は止まりません。';

  @override
  String eventTier(int tier, String score) {
    return 'ランク$tier · $scoreポイント以上';
  }

  @override
  String eventTierTokens(int count) {
    return 'イベントトークン +$count';
  }

  @override
  String eventEstimate(String score, int tier) {
    return '今参加すると：約$scoreポイント → ランク$tier';
  }

  @override
  String eventEstimateNone(String score) {
    return '今参加すると：約$scoreポイント — ランクに届きません';
  }

  @override
  String seasonalHow(int count) {
    return 'イベント中は1日1回イベントトークンを$count個受け取れます（1日は00:00 UTCに開始）。';
  }

  @override
  String seasonalBenefit(int count) {
    return 'トークンはアカウントに残ります（残高：$count）。';
  }

  @override
  String get eventDetails => '詳細';

  @override
  String get createCharacterTitle => 'キャラクター作成';

  @override
  String get classTraits => 'クラス特性';

  @override
  String classResourceLine(String name) {
    return 'リソース：$name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return '主能力値：$primary · 副：$secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return '主能力値：$primary';
  }

  @override
  String get classDifficulty => '難易度';

  @override
  String get difficultyEasy => '易しい — ミスに寛容で放置狩猟でも安全';

  @override
  String get difficultyMedium => '普通 — バランス型、スキルの順序が重要';

  @override
  String get difficultyHard => '難しい — 脆く、適切な装備が必要';

  @override
  String get skillsTapHint => 'タップで詳細';

  @override
  String get chooseRace => '種族を選択';

  @override
  String continueAs(String name) {
    return '$nameで続ける';
  }

  @override
  String get nameYourCharacter => 'キャラクターに名前を付ける';

  @override
  String get raceTroll => 'トロール';

  @override
  String get raceUndead => 'アンデッド';

  @override
  String get raceDarkElf => 'ダークエルフ';

  @override
  String get raceDrakari => 'ドラカリ';

  @override
  String get raceTauren => 'タウレン';

  @override
  String get raceGnome => 'ノーム';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '$minutes分間 攻撃力 +$atk%、防御力 +$def%';
  }

  @override
  String fedToast(int atk, int def) {
    return '満腹：攻撃力 +$atk%、防御力 +$def%';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return '満腹 · 攻撃力 +$atk% · 防御力 +$def% · 残り$time';
  }

  @override
  String huntReach(String regions) {
    return 'ここから狩猟可能：$regions';
  }

  @override
  String get huntNoMobsHere => 'ここには獲物がいません。先に狩り場へ移動してください。';

  @override
  String get huntOtherRegions => '他の地域のモンスター';

  @override
  String get huntTravelHint => 'モンスターを狩るには、その地域へ移動してください。';

  @override
  String get errMobRegion => 'このモンスターは別の地域にいます。先にそこへ移動してください。';

  @override
  String restoresHpPct(int pct) {
    return '最大HPの$pct%を回復';
  }

  @override
  String get statEvasion => '% 回避';

  @override
  String get statBlock => '% ブロック';

  @override
  String get statCritDmg => '% クリティカルダメージ';

  @override
  String get statLifesteal => '% ライフスティール';

  @override
  String get statDmgReduction => '% ダメージ軽減';

  @override
  String get classJewelry => 'アクセサリー';

  @override
  String get mapHint => 'レベルが許す地域ならどこへでも直接移動でき、時間は最短の道で決まります。地図をドラッグやピンチで動かし、地域をタップしてください。';

  @override
  String get mapNoRoad => 'ここへの直通の道はありません。隣の光る地域を経由してください。';

  @override
  String get youAreHere => '現在地です。';

  @override
  String travelTime(String time) {
    return '移動時間：$time';
  }

  @override
  String xpPerKill(int xp) {
    return '経験値 $xp';
  }

  @override
  String get deleteCharacter => 'キャラクターを削除';

  @override
  String deleteCharacterBody(String name) {
    return '$nameはすべてのアイテムとゴールドごと直ちに削除されます。元に戻せません。確認のためキャラクター名を入力してください。';
  }

  @override
  String get characterNameConfirm => 'キャラクター名';

  @override
  String characterDeletesAt(String when) {
    return '$whenに削除';
  }

  @override
  String get cancelCharacterDelete => '削除を取り消す';

  @override
  String characterDeleted(String name) {
    return '$nameを削除しました。';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$nameの削除を取り消しました。';
  }

  @override
  String get errNameMismatch => '名前が一致しません。';

  @override
  String attrGainPerPoint(String gains) {
    return '1ポイントあたり：$gains';
  }

  @override
  String attrGainPending(String gains) {
    return '追加分：$gains';
  }

  @override
  String get attrNoEffect => 'このクラスには効果なし';

  @override
  String get attrInputHint => '各能力値に割り振るポイントを入力してください';

  @override
  String get gainHp => 'HP';

  @override
  String get gainMp => 'MP';

  @override
  String get gainHpRegen => 'HP/秒';

  @override
  String get huntSupplies => '持参する補給品';

  @override
  String get huntSuppliesHint => 'HPが半分を切ると戦闘の合間にポーションを、効果が切れると料理を、マナが減るとマナポーションを使います。残りはバッグに戻ります。';

  @override
  String get supplyHp => '回復ポーション';

  @override
  String get supplyMp => 'マナポーション';

  @override
  String get supplyMpHint => 'マナ、気力、集中、自然の力、霊力を使うクラス向け';

  @override
  String get supplyFood => '料理';

  @override
  String get huntNoSupplies => 'バッグにポーションや料理がありません。ショップで買うか製作してください。';

  @override
  String restoresMpPct(int pct) {
    return 'マナを$pct%回復';
  }

  @override
  String get huntSupplyOnly => '狩猟の補給品としてのみ使用';

  @override
  String get errBadSupply => 'このアイテムは狩猟に持参できません。';

  @override
  String get errNotEnoughSupplies => 'バッグにその補給品が足りません。';

  @override
  String get errHuntOnly => '補給品として狩猟に持参してください。';

  @override
  String get huntFled => 'HPが5%を切り、戦闘から離脱しました';

  @override
  String suppliesUsed(String list) {
    return '使用：$list';
  }

  @override
  String get compareTitle => '装備中と比較';

  @override
  String compareWorn(String name) {
    return '装備中：$name';
  }

  @override
  String get compareEmpty => 'このスロットは空です。装備するとこれらがすべて加わります。';

  @override
  String get compareColStat => '能力';

  @override
  String get compareColWorn => '装備中';

  @override
  String get compareColThis => 'これ';

  @override
  String get compareColDiff => '変化';

  @override
  String get compareSame => '装備中と同じ能力です。';

  @override
  String get wearOk => 'このクラスで装備できます';

  @override
  String get wearClass => 'このクラスでは装備できません';

  @override
  String wearLevel(int level) {
    return 'レベル$levelが必要';
  }

  @override
  String get wearAttribute => '能力値が足りません';

  @override
  String get wearLegendClass => 'クラス不一致';

  @override
  String get wearLegendLater => 'まだ不可（レベルまたは能力値）';

  @override
  String get enchantAction => 'エンチャント';

  @override
  String get enchantPick => 'エンチャントする装備を選択';

  @override
  String get enchantNone => 'この巻物に合う装備がありません。';

  @override
  String enchantFits(String groups) {
    return '対象：$groups';
  }

  @override
  String get enchantGroupWeapon => '武器';

  @override
  String get enchantGroupArmor => '防具';

  @override
  String get enchantGroupJewelry => 'アクセサリー';

  @override
  String get enchantGroupAny => 'すべての装備';

  @override
  String get enchantReplaceTitle => 'エンチャントを置き換えますか？';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$oldは失われ、$nextに置き換わります。';
  }

  @override
  String enchantDone(String item) {
    return '$itemをエンチャントしました';
  }

  @override
  String enchantOn(String name) {
    return 'エンチャント：$name';
  }

  @override
  String get wornTag => '装備中';

  @override
  String get errNotEnchant => 'エンチャントの巻物ではありません。';

  @override
  String get errEnchantSlot => 'このエンチャントはその装備に合いません。';

  @override
  String get statGatherMining => '採掘の産出量';

  @override
  String get statGatherHerbalism => '薬草採集の産出量';

  @override
  String get statGatherWoodcutting => '伐採の産出量';

  @override
  String get statGatherFishing => '釣りの産出量';

  @override
  String get statGatherHunting => '狩猟の産出量';

  @override
  String get statCraftSpeed => '製作速度';

  @override
  String diffXp(int pct) {
    return '経験値 $pct%';
  }

  @override
  String get diffNoXp => '経験値なし';

  @override
  String get diffLegend => 'レシピの色：現在のスキルで製作から得られる経験値';

  @override
  String get recipeGroupEnchantWeapon => '武器エンチャント';

  @override
  String get recipeGroupEnchantArmor => '防具エンチャント';

  @override
  String get recipeGroupEnchantJewelry => 'アクセサリーエンチャント';

  @override
  String get recipeGroupEnchantProfession => '職業エンチャント';

  @override
  String get partyHint => '同じ作業（同じ採集、職業、狩猟、ダンジョン）を同時に行うメンバーは、作業中の人数ごとに経験値+3%を得ます。5人なら+15%。誰の経験値も減りません。';

  @override
  String get partyChat => 'パーティーチャット';

  @override
  String partyShared(String names, int kept, int earned) {
    return '$namesとのパーティーボーナス：$earned → $kept経験値';
  }

  @override
  String evRaidMatched(String raid) {
    return '$raidのパーティーが見つかりました。レイドを確認してください。';
  }

  @override
  String get nkRaidMatched => 'レイドパーティー発見';

  @override
  String get workHunt => '狩猟';

  @override
  String get workDungeon => 'ダンジョン';

  @override
  String get raidFind => 'パーティーを探す';

  @override
  String get raidQueueJoin => 'キューに入る';

  @override
  String get raidQueueLeave => 'キューを出る';

  @override
  String raidQueued(String raid, int waiting) {
    return '$raidのキュー · $waiting人待機中';
  }

  @override
  String get raidQueueHint => '5人待機するか、2人以上で5分経つとランダムパーティーが組まれます。';

  @override
  String get raidAdvertise => '募集する';

  @override
  String get raidUnlist => '募集を取り消す';

  @override
  String raidListed(String raid) {
    return 'パーティーが$raidの掲示板に載りました';
  }

  @override
  String get raidBoard => '募集掲示板';

  @override
  String get raidBoardEmpty => '募集はありません。';

  @override
  String get raidJoinParty => 'パーティーに参加';

  @override
  String get raidPickRaid => 'レイド';

  @override
  String get errListingGone => 'その募集はなくなったか満員です。';

  @override
  String get workGathering => '採集';

  @override
  String get workAdventure => '冒険';

  @override
  String get workRaid => 'レイド';

  @override
  String get workTravel => '移動中';

  @override
  String get workCrafting => '製作';

  @override
  String get workIdle => '待機';

  @override
  String get errDungeonRegion => '先に地図でこの場所へ移動してください。';

  @override
  String get dungeonsElsewhere => '移動するほどダンジョンとレイドが増えます。それぞれ地図上の地域から入ります。';

  @override
  String get dungeonNoneHere => 'ここから行けるダンジョンはありません。地図でダンジョンの地域へ移動してください。';

  @override
  String dungeonXp(int xp) {
    return 'あなたには約$xp経験値';
  }

  @override
  String appVersionLine(String version) {
    return 'バージョン $version';
  }

  @override
  String get updateButton => 'アップデート';

  @override
  String updateAvailable(String version) {
    return '新しいバージョンあり：$version';
  }

  @override
  String get updateLatest => '最新バージョンです';

  @override
  String updateConfirm(String version, int size) {
    return 'バージョン$version（$size MB）をダウンロードしてインストールしますか？';
  }

  @override
  String get updateDownloading => 'ダウンロード中…';

  @override
  String get updatePermission => 'MeNa Worldからのインストールを許可し、もう一度アップデートをタップしてください。';

  @override
  String get updateFailed => 'アップデートを取得できませんでした。接続を確認してください。';

  @override
  String get fightReport => '戦闘レポート';

  @override
  String get endTimeUp => '時間切れ：予定の終わりまで続きました。';

  @override
  String get endEarly => '早めに受取：予定の終わりより前に停止しました。';

  @override
  String endFled(String mob) {
    return '$mobとの戦闘でHPが5%を切り、離脱して戻りました。';
  }

  @override
  String endDefeat(String mob) {
    return '$mobに倒されました。';
  }

  @override
  String endPartyDefeat(String mob) {
    return 'パーティーは$mobに倒されました。';
  }

  @override
  String endTimeout(String mob) {
    return '$mobとの戦闘がターン制限を超えたため停止しました。';
  }

  @override
  String get endCleared => 'クリア：すべてのウェーブとボスを倒しました。';

  @override
  String get repCampHint => 'HPが50%を切ると狩人は次の戦闘までキャンプで回復するため、5%まで落ちることはまれです。';

  @override
  String get repYou => 'あなた';

  @override
  String get repLevel => 'レベル';

  @override
  String get repHp => 'HP';

  @override
  String get repAttack => '攻撃力';

  @override
  String get repDefense => '防御力';

  @override
  String get repEvasion => '回避';

  @override
  String get repHitChance => '命中率';

  @override
  String get repCrit => 'クリティカル';

  @override
  String repFights(int n) {
    return '戦闘$n回';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '命中 $hits/$swings（$pct%）、平均 $avg、最大 $max';
  }

  @override
  String get repYourSwings => 'あなたの攻撃';

  @override
  String get repPartySwings => 'パーティーの攻撃';

  @override
  String get repEnemySwings => '敵の攻撃';

  @override
  String repLowest(int pct) {
    return '戦闘中の最低HP：$pct%';
  }

  @override
  String repRest(String time) {
    return 'キャンプ：$time';
  }

  @override
  String get restart => '再開';

  @override
  String get restartTitle => '受け取ってやり直しますか？';

  @override
  String get restartBody => 'これまでの分はバッグに入り、同じ内容が同じ時間で再び始まります。';

  @override
  String get dailyTitle => 'デイリーログインボーナス';

  @override
  String get dailyBody => 'どのキャラクターでも1日1回受け取れます。週が続くほど増え、29日目が最大です。1日でも逃すと最初からになります。';

  @override
  String get dailyClaim => '受け取る';

  @override
  String get dailyTomorrow => 'また明日';

  @override
  String dailyDay(int n) {
    return '$n日目';
  }

  @override
  String get dailyOpen => 'デイリー報酬';

  @override
  String chatMutedNotice(String name, int min) {
    return '$nameは暴言により$min分間チャット制限されました。';
  }

  @override
  String get storyTitle => 'ストーリー';

  @override
  String get storyOpen => 'ストーリーモード';

  @override
  String get storyComingSoon => 'ストーリーモードは近日公開です。章が書かれ次第ここに表示されます。';

  @override
  String storyLocked(int level) {
    return 'レベル$levelで開放';
  }

  @override
  String get bankTitle => '銀行';

  @override
  String get bankHint => '全キャラクター共通です。バッグのアイテムをタップで預け、銀行のアイテムをタップで引き出します。ソウルバウンドのアイテムはそのキャラクターに残ります。';

  @override
  String get bankEmpty => '銀行は空です。';

  @override
  String get bankDeposit => '預ける';

  @override
  String get bankWithdraw => '引き出す';

  @override
  String get errBankFull => '銀行がいっぱいです。';

  @override
  String get tabTrade => '取引';

  @override
  String get tradeHint => '他のプレイヤーとアイテムを取引します。双方がアイテムを出して確認し、それぞれ渡すアイテムの平均相場の8%を1個ごとに税として支払います。';

  @override
  String get tradeWith => 'プレイヤー名';

  @override
  String get tradeOpen => '取引を開始';

  @override
  String tradeTitle(String name) {
    return '$nameさんと取引';
  }

  @override
  String get tradeYours => 'あなたが渡す';

  @override
  String tradeTheirs(String name) {
    return '$nameさんが渡す';
  }

  @override
  String get tradeNothing => 'まだありません。';

  @override
  String get tradeAdd => 'アイテムを追加';

  @override
  String get tradeRemove => '戻す';

  @override
  String get tradeTax => 'あなたの税';

  @override
  String get tradeReady => '準備完了';

  @override
  String get tradeWaiting => '未準備';

  @override
  String get tradeConfirm => '取引を確定';

  @override
  String get tradeWaitingThem => '相手を待っています';

  @override
  String get tradeCancel => '取引を中止';

  @override
  String get tradeCancelBody => '何も交換されません。';

  @override
  String get tradeDone => '取引が完了しました。';

  @override
  String get errTradeBusy => 'どちらかが既に取引中です。';

  @override
  String get errTradeClosed => 'この取引は終了しました。';

  @override
  String get errTradeItem => 'このアイテムは今は取引できません。';

  @override
  String get errTradeEmpty => '先に何か出してください。';

  @override
  String get errTradeTooMany => '取引のアイテムが多すぎます。';

  @override
  String get errCraftLevel => 'キャラクターのレベルがこのアイテムのレベルより低いです。';

  @override
  String get gatherNoneHere => 'あなたの資源は他の地図にあります。横に表示された場所へ移動してください。';

  @override
  String get errGatherRegion => 'この資源は他の地図にあります。';

  @override
  String get helpersOpen => 'ヘルパー';

  @override
  String get helpersTitle => 'ヘルパー';

  @override
  String get helpersHint => '採集や製作を手伝うヘルパーを雇いましょう。速度は半分、経験値も半分で、経験値はあなたの職業に入ります。種類ごとに固有のボーナスがあります。';

  @override
  String get helpersGather => '採集';

  @override
  String get helpersCraft => '製作';

  @override
  String get helperHire => '雇う';

  @override
  String helperHireBody(String name, String price) {
    return '$nameを$priceで雇いますか？';
  }

  @override
  String get helperSend => '仕事に送る';

  @override
  String get helperCancel => '呼び戻す';

  @override
  String get helperCancelBody => '採集は何も戻らず、製作は素材とゴールドが戻ります。';

  @override
  String get helperDone => '完了 — 受け取る';

  @override
  String helperWorking(String time) {
    return '作業中、残り$time';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return 'アイテム$items個、$profession経験値 +$xp';
  }

  @override
  String get helperNoRecipes => 'この職業にはまだレベルに合うレシピがありません。';

  @override
  String get helperPickRecipe => 'レシピを選択';

  @override
  String helperBonusYield(int pct) {
    return '産出量 +$pct%';
  }

  @override
  String helperBonusSpeed(int pct) {
    return '$pct%高速';
  }

  @override
  String helperBonusXp(int pct) {
    return '経験値 +$pct%';
  }

  @override
  String helperBonusGold(int pct) {
    return 'ゴールド費用 $pct%減';
  }

  @override
  String helperBonusDouble(int pct) {
    return '$pct%の確率で2倍産出';
  }

  @override
  String get helperMiner => '鉱夫';

  @override
  String get helperWoodsman => '木こり';

  @override
  String get helperAngler => '釣り師';

  @override
  String get helperHerbalist => '薬草師';

  @override
  String get helperTrapper => '罠師';

  @override
  String get helperSmith => '鍛冶師';

  @override
  String get helperTanner => '革職人';

  @override
  String get helperWeaver => '織工';

  @override
  String get helperCarver => '木工師';

  @override
  String get helperApothecary => '薬剤師';

  @override
  String get helperCook => '料理人';

  @override
  String get helperSmelter => '精錬工';

  @override
  String get helperScribe => '写本師';

  @override
  String get helperGemcutter => '宝石職人';

  @override
  String get errHelperBusy => 'このヘルパーは既に作業中です。';

  @override
  String get errHelperNotReady => 'ヘルパーはまだ終わっていません。';

  @override
  String get shopOpen => 'ショップ';

  @override
  String get shopTitle => 'クレジットショップ';

  @override
  String get shopCredits => 'クレジット';

  @override
  String get shopAvatars => 'アバター';

  @override
  String get shopBackgrounds => '背景';

  @override
  String shopBalance(int n) {
    return '$nクレジット';
  }

  @override
  String get shopBuyCredits => 'クレジットを購入してブーストやコスメを解放しましょう。';

  @override
  String get shopPaymentsSoon => 'クレジットの購入は近日公開です。';

  @override
  String get shopBoosts => '経験値ブースト';

  @override
  String get shopBuy => '購入';

  @override
  String shopBuyBody(int n) {
    return '$nクレジットを使いますか？';
  }

  @override
  String shopPrice(int n) {
    return '$n cr';
  }

  @override
  String get boostGather => '採集';

  @override
  String get boostAdventure => '戦闘';

  @override
  String get boostCraft => '職業';

  @override
  String get boostDesc => '2時間 経験値 +100%';

  @override
  String get boostShort => '経験値 +100%、2時間';

  @override
  String boostLeft(String time) {
    return '有効、残り$time';
  }

  @override
  String boostWatchAd(int n) {
    return '広告を見る（残り$n回）';
  }

  @override
  String get boostAdsSoon => '広告は近日公開';

  @override
  String get errInsufficientCredits => 'クレジットが足りません。';

  @override
  String get errBoostActive => 'このブーストは既に有効です。';

  @override
  String get errAdLimit => '今日の広告ブーストは終わりました。';

  @override
  String helperCanMake(int n) {
    return '製作可能：$n';
  }

  @override
  String get helperPickResource => '採集するもの';

  @override
  String get helperCraftHint => '素材とゴールドは今支払われ、その間もプレイを続けられます。';

  @override
  String get petsOpen => 'ペット';

  @override
  String get petsHint =>
      'ペットの卵はレベル400以上のレイドをクリアすると手に入ります。卵は孵化器でかえり、ペットは自分の狩りで成長し、アクティブなペットはボーナスを与えます。さらに戦闘にも参加し、すべての攻撃にダメージを加えます。';

  @override
  String get petsEggs => '卵';

  @override
  String get petsNoEggs => '卵がありません。レベル400以上のレイドをクリアして見つけましょう。';

  @override
  String petsMine(int n, int max) {
    return 'マイペット（$n/$max）';
  }

  @override
  String get petHatch => '孵化';

  @override
  String petHatchBody(String egg, String time) {
    return '$eggを孵化器に入れますか？$timeでかえります。';
  }

  @override
  String get petHatching => '卵の中';

  @override
  String get petHatched => '孵化しました！';

  @override
  String petPower(int n) {
    return 'レベル$nのモンスターのように戦います';
  }

  @override
  String get petActive => 'アクティブ';

  @override
  String get petActivate => 'アクティブにする';

  @override
  String get petDeactivate => '休ませる';

  @override
  String get petHunt => '狩りに出す';

  @override
  String petHunting(String mob) {
    return '$mobを狩猟中';
  }

  @override
  String get petHuntDone => '狩り終了';

  @override
  String get petCollectEarly => '今狩りを終えると、これまでの戦闘分が支払われます。';

  @override
  String petCollected(int wins, int xp) {
    return '$wins勝、ペット経験値 +$xp';
  }

  @override
  String petLevelUp(String name, int level) {
    return '$nameがレベル$levelに到達！';
  }

  @override
  String get petDied => 'ペットが倒され、狩りは早く終わりました。';

  @override
  String get petPickMob => '狩るモンスター';

  @override
  String petXpPerWin(int n) {
    return '1勝ごとに経験値 +$n';
  }

  @override
  String get errPetCap => '所持できるペットの上限です。';

  @override
  String get errPetHatching => 'ペットはまだかえっていません。';

  @override
  String get errPetHunting => 'ペットは既に狩りに出ています。';

  @override
  String petStageLevels(int from, int to) {
    return 'Lv.$from–$to';
  }

  @override
  String get regionMobs => 'モンスター';

  @override
  String get regionDungeons => 'ダンジョン';

  @override
  String get regionGather => '資源';

  @override
  String get regionEmpty => 'この地域に属するモンスター、ダンジョン、資源、クエストはありません。';

  @override
  String get regionRaids => 'レイド';

  @override
  String get tagUse => '用途';

  @override
  String get tagAllClasses => '全クラス';

  @override
  String get tagDrops => 'ドロップ元';

  @override
  String get tagRaid => 'レイド';

  @override
  String get tagDungeon => 'ダンジョン';

  @override
  String get tagCrafted => '製作元';

  @override
  String tagSkill(int n) {
    return 'スキル$n';
  }

  @override
  String get tagSold => '販売元';

  @override
  String get tagGathered => '採集元';

  @override
  String get tagUsedIn => '素材の用途';

  @override
  String get petEvolution => '進化';

  @override
  String get petStrength => '強さの詳細';

  @override
  String petAssistRange(int min, int max) {
    return '戦闘支援：すべての攻撃にダメージ +$min% ～ +$max%';
  }

  @override
  String petBonusRange(String min, String max) {
    return 'ボーナス：$min ～ $max';
  }

  @override
  String petFightsRange(int min, int max) {
    return 'レベル$min～$maxのモンスターのように狩ります';
  }

  @override
  String get petFind => '入手方法';

  @override
  String petAssist(int pct) {
    return '戦闘支援：ダメージ +$pct%';
  }

  @override
  String get petLevelNote => 'レベル1から100まで';

  @override
  String get petIndex => 'ペット図鑑';

  @override
  String petIndexCount(int n, int max) {
    return '$n / $max 発見';
  }

  @override
  String slotCount(int n, int max) {
    return '$n / $max';
  }

  @override
  String get petEggTapHint => '卵をタップすると孵化器に入れられます。';

  @override
  String get petTapHint => 'ペットをタップで表示 · 緑の点がアクティブなペットです。';

  @override
  String get petIndexHint => 'ペットをタップするとページが開きます。';

  @override
  String get petUnknown => '未発見';

  @override
  String get petAbout => '概要';

  @override
  String get petHabitat => '生息地';

  @override
  String petElement(String name) {
    return '属性：$name';
  }

  @override
  String get petStage1 => '幼体';

  @override
  String get petStage2 => '若体';

  @override
  String get petStage3 => '成体';

  @override
  String get bodyMale => '男性';

  @override
  String get bodyFemale => '女性';

  @override
  String get bodyChoice => '外見';

  @override
  String get bodyChange => '外見を変更';

  @override
  String get portraitNoArmor => '防具なし';

  @override
  String get combatPower => '戦闘力';

  @override
  String get combatTapHint => 'タップで詳細';

  @override
  String statCap(int cap) {
    return '上限 $cap%';
  }

  @override
  String get cEvasion => '回避';

  @override
  String get cBlock => 'ブロック';

  @override
  String get cCritDmg => 'クリティカルダメージ';

  @override
  String get cDmgReduction => 'ダメージ軽減';

  @override
  String get cLifesteal => 'ライフスティール';

  @override
  String get cExecute => '処刑';

  @override
  String get cMissingHp => '背水';

  @override
  String get cPetAssist => 'ペット支援';

  @override
  String get xPower => 'ワールドイベントで使う強さ：攻撃力×2 + 防御力 + 最大HP / 10。';

  @override
  String get xAttack => '敵の防御を適用する前の各攻撃の基本ダメージ。';

  @override
  String get xDefense => '受ける攻撃ごとに防御力の半分だけ減ります。';

  @override
  String xCrit(int cap) {
    return '攻撃がクリティカルになる確率。最大$cap%。';
  }

  @override
  String get xCritDmg => 'クリティカルは通常攻撃のこの倍率のダメージを与えます。';

  @override
  String get xHit => '攻撃が当たる確率から敵の回避を引いた値（15～95%）。';

  @override
  String xEvasion(int cap) {
    return '敵の命中率から引かれます。最大$cap%。';
  }

  @override
  String xBlock(int cap) {
    return '受ける攻撃を半減する確率。最大$cap%。';
  }

  @override
  String xDmgReduction(int cap) {
    return '受ける攻撃はこの割合だけ減ります。最大$cap%。';
  }

  @override
  String get xLifesteal => '与えたダメージのこの割合だけ回復します。';

  @override
  String xExecute(int pct) {
    return 'HPが$pct%以下の敵への追加ダメージ。';
  }

  @override
  String get xMissingHp => '自分のHPが減るほど、この値まで追加ダメージ。';

  @override
  String get xPetAssist => 'アクティブなペットが各攻撃にこの割合を加えます。';

  @override
  String get lookFace => '顔';

  @override
  String get lookHairStyle => '髪型';

  @override
  String get lookHairColor => '髪の色';

  @override
  String get lookEyeColor => '目の色';

  @override
  String get lookSave => '保存';

  @override
  String get lookHumanOnly => 'この種族の顔と髪の選択肢は近日公開です。';

  @override
  String get raceAbilities => '種族能力';

  @override
  String get raceAbilitiesHint => 'すべてパッシブ';

  @override
  String get chooseClass => 'クラスを選択';

  @override
  String raceAttrBonus(int pct, String name) {
    return '$name +$pct%';
  }

  @override
  String raceProfessionXp(int pct, String name) {
    return '$nameのボーナス経験値 +$pct%';
  }

  @override
  String get language => '言語';

  @override
  String get report => '通報';

  @override
  String reportTitle(String name) {
    return '$name を通報';
  }

  @override
  String get reportWhy => '何がありましたか？';

  @override
  String get reportHarassment => '嫌がらせ・暴言';

  @override
  String get reportCheating => 'チート・ボット';

  @override
  String get reportName => '不適切な名前';

  @override
  String get reportSpam => 'スパム・宣伝';

  @override
  String get reportScam => '詐欺';

  @override
  String get reportOther => 'その他';

  @override
  String get reportNote => '詳細（任意）';

  @override
  String get reportSent => '通報を受け付けました。ありがとうございます。';

  @override
  String get errReportSelf => '自分のキャラクターは通報できません。';

  @override
  String get nkReportReviewed => '通報が確認されました';

  @override
  String get guide => 'ガイド';

  @override
  String get guideSearch => 'ガイドを検索';

  @override
  String get guideNoMatch => '見つかりませんでした。';

  @override
  String get guideShowMe => '見に行く';

  @override
  String get next => '次へ';

  @override
  String get back => '戻る';

  @override
  String get tutorialSkip => 'チュートリアルをスキップ';

  @override
  String get tutorialStart => 'プレイ開始';

  @override
  String get tutorialReplay => 'チュートリアルをもう一度見る';

  @override
  String get tutorialWelcomeTitle => 'MeNa World へようこそ';

  @override
  String get tutorialWelcomeBody => 'アプリを閉じていてもキャラクターは働き続けます。仕事を任せて、あとで戻って成果を受け取りましょう。';

  @override
  String get tutorialCharacterBody => 'ステータス、能力ポイント、スキルはここにあります。レベルごとに得るポイントを割り振り、レベル300と700で昇格を選びましょう。';

  @override
  String get tutorialGatherBody => '資源と時間を選び、あとで受け取りに戻りましょう。採掘・釣り・伐採・薬草採取・狩猟はそれぞれ別にレベルが上がります。';

  @override
  String get tutorialBagBody => 'アイテムをタップして装備・使用・売却・出品ができます。比較表で今の装備より良いか確認できます。';

  @override
  String get tutorialAdventureBody => 'モンスター狩り、ダンジョン攻略、レイドやイベント参加、マップの旅ができます。モンスターとダンジョンは今いる地域で変わります。';

  @override
  String get tutorialWorkshopBody => '職業で製作し、ショップで購入し、マーケットで他のプレイヤーと取引しましょう。';

  @override
  String get tutorialTopBarTitle => '上部バー';

  @override
  String get tutorialTopBarBody =>
      'クエスト、通知、ソーシャル（チャット、フレンド、メール、ギルド、トレード）は右上にあります。キャラクターは同時に1つの作業しかできないので、次を始める前に完了するか受け取りましょう。';

  @override
  String get tutorialDoneTitle => '準備完了';

  @override
  String get tutorialDoneBody => '場所を忘れたら、キャラクタータブの下にあるガイドを開いてください。';

  @override
  String get guideCharacterTitle => 'ステータス・能力・スキル';

  @override
  String get guideCharacterBody => 'キャラクタータブ。能力ポイントは進行カードで割り振ります。スキルは折りたたみ欄にあり、同じ役割のスキルは最も強いものだけが残ります。';

  @override
  String get guidePromotionTitle => '昇格';

  @override
  String get guidePromotionBody => 'キャラクタータブの進行カード。レベル300と700では2つの道から1つ選ぶまで経験値が保留され、選択時に支払われます。';

  @override
  String get guideGatherTitle => '採集';

  @override
  String get guideGatherBody => '採集タブ。職業、手の届く資源、1・3・6時間を選びます。資源は地域ごとにあるので、良いものは旅をして探しましょう。';

  @override
  String get guideBagTitle => 'バッグと装備';

  @override
  String get guideBagBody => 'バッグタブ。アイテムをタップすると詳細と比較が見られます。長押ししてドラッグで複数選択し、まとめて売却や破棄ができます。バッグが満杯だと新しい作業を始められません。';

  @override
  String get guideBankTitle => 'アカウント銀行';

  @override
  String get guideBankBody => 'バッグタブの銀行ボタン。全キャラクターで共有するので、アイテムの受け渡しに使えます。バインド済みや装備中のアイテムは入れられません。';

  @override
  String get guideHuntTitle => '狩り';

  @override
  String get guideHuntBody => '冒険タブの狩り。手の届くモンスター、時間、持っていくポーションと食事を選びます。キャラクターは重傷になると自動で撤退します。';

  @override
  String get guideDungeonTitle => 'ダンジョンとレイド';

  @override
  String get guideDungeonBody => '冒険タブのダンジョン。ダンジョンはソロです。レイドは2〜5人のパーティが必要で、パーティ検索でマッチできます。レイドボスはセット装備とペットの卵を落とします。';

  @override
  String get guideTravelTitle => '旅とマップ';

  @override
  String get guideTravelBody => '冒険タブの旅。地域をタップするとモンスター、ダンジョン、資源、クエストが表示され、移動できます。レベル条件は目的地だけです。';

  @override
  String get guideEventsTitle => 'イベント';

  @override
  String get guideEventsBody => '冒険タブのイベント。ワールドボスや侵攻に参加し、毎日のシーズントークンを受け取りましょう。カードをタップすると内容と報酬が見られます。';

  @override
  String get guideCraftTitle => '職業と製作';

  @override
  String get guideCraftBody => '工房タブの製作。職業を開くとレシピが見られます。色はそのレシピで得られる経験値を示します。付呪の巻物はバッグから使います。';

  @override
  String get guideVendorTitle => 'ショップと買い戻し';

  @override
  String get guideVendorBody => '工房タブのショップ。売却はバッグから行います。売ったものは2時間以内なら買い戻せます。';

  @override
  String get guideMarketTitle => 'マーケットとオークション';

  @override
  String get guideMarketBody =>
      '工房タブのマーケット。固定価格で購入するか、オークションに入札します。バッグから出品できます：固定価格、または開始価格と任意の即決価格を付けたオークション。終了間際の入札で時間が延長されます。';

  @override
  String get guideQuestsTitle => 'クエスト';

  @override
  String get guideQuestsBody => '右上のクリップボードアイコン。ストーリー、サブ、デイリークエストがあり、バッジは報告可能な数を示します。';

  @override
  String get guideNotificationsTitle => '通知と報酬';

  @override
  String get guideNotificationsBody => '右上のベルアイコン。報酬タブは受け取った内容を30日間保存します。設定で受け取る通知を選べます。';

  @override
  String get guideSocialTitle => 'チャット・フレンド・メール・トレード';

  @override
  String get guideSocialBody => '右上の吹き出し。トレードは双方の確認が必要で、税がかかります。ブロックするとチャット、フレンド、メールから相手が消えます。';

  @override
  String get guidePartyTitle => 'パーティ';

  @override
  String get guidePartyBody => 'キャラクタータブのパーティ欄。名前で招待します。同じ種類の作業をするメンバーはボーナス経験値を得ます。';

  @override
  String get guideGuildTitle => 'ギルド・金庫・ギルド銀行';

  @override
  String get guideGuildBody => '右上の吹き出しのギルドタブ。メンバーはゴールドとアイテムを預け、階級によって引き出しの可否と1日の上限が決まります。すべて記録されます。';

  @override
  String get guideRankingTitle => 'ランキング';

  @override
  String get guideRankingBody => 'キャラクタータブのランキングボタン。スコアはレベル、実績、熟練、職業、ダンジョンとレイドの攻略で決まります。職業や表で絞り込めます。';

  @override
  String get guidePetsTitle => 'ペット';

  @override
  String get guidePetsBody => 'キャラクタータブのペットボタン。卵は高レベルのレイドボスから入手します。孵化させて狩りで育て、ボーナスと戦闘支援のために1体を有効にしましょう。';

  @override
  String get guideHelpersTitle => 'ヘルパー';

  @override
  String get guideHelpersBody => 'キャラクタータブのヘルパーボタン。職業ごとにヘルパーを雇えます。キャラクターの作業とは別に働きます。';

  @override
  String get guideReportTitle => 'プレイヤーの通報';

  @override
  String get guideReportBody => 'チャットの行を長押しするか、フレンドのメニューを使います。理由を選ぶとチームが確認し、確認後に通知します。';

  @override
  String get guideAccountTitle => 'アカウント・言語・サポート';

  @override
  String get guideAccountBody => 'キャラクター一覧の上にあるアカウントアイコン。言語やパスワードの変更、サポートへの連絡、データのダウンロード、アカウント削除ができます。';

  @override
  String get ranking => 'ランキング';

  @override
  String get rankingScore => 'スコア';

  @override
  String get rankingLevel => 'レベル';

  @override
  String get rankingAchievements => '実績';

  @override
  String get rankingMastery => '熟練';

  @override
  String get rankingProfessions => '職業';

  @override
  String get rankingClears => '攻略';

  @override
  String get rankingDungeons => 'ダンジョン';

  @override
  String get rankingRaids => 'レイド';

  @override
  String get rankingAllClasses => 'すべてのクラス';

  @override
  String get rankingHowScore => 'スコアの内訳';

  @override
  String get rankingScoreNote => '各ダンジョンとレイドは1回だけ数えます。数分ごとに更新されます。';

  @override
  String get rankingEmpty => 'ランキングを集計中です。';

  @override
  String get guildVault => '金庫と銀行';

  @override
  String get guildVaultTreasury => 'ギルド金庫';

  @override
  String get guildVaultBank => 'ギルド銀行';

  @override
  String get guildVaultDepositGold => 'ゴールドを預ける';

  @override
  String get guildVaultWithdrawGold => 'ゴールドを引き出す';

  @override
  String get guildVaultGoldAmount => 'ゴールドの額';

  @override
  String get guildVaultNoLimit => '1日の上限なしで引き出せます。';

  @override
  String guildVaultGoldToday(String used, String limit) {
    return '本日の引き出し：$used / $limit';
  }

  @override
  String guildVaultItemsToday(int used, int limit) {
    return '本日取り出したアイテム：$used / $limit';
  }

  @override
  String get guildVaultNoItems => '今の階級ではギルド銀行を使えません。';

  @override
  String get guildVaultLog => '履歴';

  @override
  String get guildVaultLogEmpty => 'まだ履歴はありません。';

  @override
  String guildVaultLogIn(String name, String what) {
    return '$name が $what を預けました';
  }

  @override
  String guildVaultLogOut(String name, String what) {
    return '$name が $what を引き出しました';
  }

  @override
  String get errGuildVaultNotEmpty => '解散する前に金庫と銀行を空にしてください。';

  @override
  String get errGuildVaultClosing => 'ギルドは解散中のため預けられません。';

  @override
  String get errGuildVaultLimit => '1日の上限を超えています。';

  @override
  String get errItemBound => 'バインド済みのアイテムは移動できません。';

  @override
  String get errItemEquipped => '先に装備を外してください。';

  @override
  String get bid => '入札';

  @override
  String get bidPlaced => '入札しました。';

  @override
  String get marketFixed => '固定価格';

  @override
  String get marketAuction => 'オークション';

  @override
  String auctionNoBids(String price) {
    return '入札なし、開始価格 $price';
  }

  @override
  String auctionTopBid(String price, int count) {
    return '最高入札 $price（$count件）';
  }

  @override
  String get auctionYouLead => 'あなたが最高入札者';

  @override
  String get auctionStartBid => 'ロット全体の開始価格';

  @override
  String get auctionBuyoutOptional => '即決価格（任意）';

  @override
  String get auctionRules => '時間切れで最高入札者が落札します。残り5分以内の入札で時間が延長されます。入札後は取り消せません。';

  @override
  String get auctionYouReceiveAtLeast => '最低受取額';

  @override
  String auctionMinBid(String price) {
    return '最低入札額：$price';
  }

  @override
  String get auctionBidCopper => '入札額（銅貨）';

  @override
  String get auctionHeldNote => 'ゴールドは預かられ、他の人が上回るとすぐ返金されます。';

  @override
  String auctionBuyoutFor(String price) {
    return '即決：$price';
  }

  @override
  String get errBidTooLow => '最低入札額を下回っています。';

  @override
  String get errTopBidder => 'すでに最高入札者です。';

  @override
  String get errAuctionBusy => '他の入札を処理中です。少し待ってからお試しください。';

  @override
  String get errAuctionHasBids => '入札があるオークションは取り消せません。';

  @override
  String get nkMarketOutbid => '入札額を上回られました';

  @override
  String get nkMarketWon => 'オークションに落札しました';

  @override
  String get evMarketOutbid => '入札額を上回られ、ゴールドが返金されました。';

  @override
  String get skillNeedsShield => '盾が必要';
}
