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
  String get tagline => '离线时，冒险依然继续。';

  @override
  String get emailOrUsername => '邮箱或用户名';

  @override
  String get email => '邮箱';

  @override
  String get username => '用户名';

  @override
  String get password => '密码';

  @override
  String get passwordConfirm => '确认密码';

  @override
  String get showPassword => '显示密码';

  @override
  String get hidePassword => '隐藏密码';

  @override
  String get login => '登录';

  @override
  String get register => '创建账号';

  @override
  String get noAccount => '新玩家？创建账号';

  @override
  String get haveAccount => '已有账号？登录';

  @override
  String get acceptTerms => '我已年满 13 岁，并接受服务条款和隐私政策';

  @override
  String get fieldRequired => '必填';

  @override
  String get invalidEmail => '请输入有效的邮箱';

  @override
  String get usernameRule => '3–32 个字母、数字、_ 或 .';

  @override
  String get passwordMin => '至少 8 个字符';

  @override
  String get passwordMismatch => '两次输入的密码不一致';

  @override
  String get termsRequired => '请接受条款后继续';

  @override
  String get logout => '退出登录';

  @override
  String get retry => '重试';

  @override
  String get offline => '无法连接服务器，请检查网络。';

  @override
  String get characters => '角色';

  @override
  String slotsUsed(int used, int cap) {
    return '槽位 $used/$cap';
  }

  @override
  String get newCharacter => '新建角色';

  @override
  String get noCharacters => '还没有角色，创建你的第一位英雄吧。';

  @override
  String get create => '创建';

  @override
  String get characterName => '名称';

  @override
  String get nameRule => '最多 16 个字符，只用一种文字，以字母或文字开头（可含数字和 _）';

  @override
  String get race => '种族';

  @override
  String get characterClass => '职业';

  @override
  String get raceHuman => '人类';

  @override
  String get raceElf => '精灵';

  @override
  String get raceDwarf => '矮人';

  @override
  String get raceOrc => '兽人';

  @override
  String get classWarrior => '战士';

  @override
  String get classMage => '法师';

  @override
  String get classRogue => '盗贼';

  @override
  String levelLine(int level, String cls, String race) {
    return '$level 级 · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return '$level 级';
  }

  @override
  String get hp => '生命';

  @override
  String get mp => '法力';

  @override
  String get stamina => '耐力';

  @override
  String get energy => '能量';

  @override
  String get xp => '经验';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => '属性';

  @override
  String get combat => '战斗';

  @override
  String get statStr => '力量';

  @override
  String get statAgi => '敏捷';

  @override
  String get statInt => '智力';

  @override
  String get statVit => '体质';

  @override
  String get statCha => '智慧';

  @override
  String get statLuck => '幸运';

  @override
  String get attack => '攻击';

  @override
  String get defense => '防御';

  @override
  String get crit => '暴击';

  @override
  String get hit => '命中';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => '金币';

  @override
  String get switchCharacter => '切换角色';

  @override
  String get statePendingDelete => '等待删除';

  @override
  String get stateFrozen => '已冻结';

  @override
  String get errInvalidCredentials => '邮箱/用户名或密码错误。';

  @override
  String get errAccountLocked => '你的账号已被临时锁定。';

  @override
  String get errAccountBanned => '此账号已被封禁。';

  @override
  String get errAccountExists => '该邮箱或用户名已被占用。';

  @override
  String get errRateLimited => '尝试次数过多，请稍候。';

  @override
  String get errNameTaken => '该名称已被占用。';

  @override
  String get errSlotCap => '角色槽位已满。';

  @override
  String get errNotPlayable => '此角色暂时无法游玩。';

  @override
  String get errSessionExpired => '登录已过期，请重新登录。';

  @override
  String errGeneric(String code) {
    return '出现问题（$code）。';
  }

  @override
  String get navCharacter => '角色';

  @override
  String get navActivity => '采集';

  @override
  String get actMining => '采矿';

  @override
  String get actFishing => '钓鱼';

  @override
  String get actWoodcutting => '伐木';

  @override
  String get actHerbalism => '采药';

  @override
  String get actHunting => '狩猎';

  @override
  String get activityChoose => '选择活动';

  @override
  String get durationLabel => '时长';

  @override
  String get gatherResource => '资源';

  @override
  String gatherXp(int xp, int next) {
    return '$xp / $next 经验';
  }

  @override
  String get gatherMax => '最高等级';

  @override
  String gatherLocked(int level) {
    return '$level 级解锁';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name +$xp 经验';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '$name达到 $level 级！';
  }

  @override
  String get errGatherLevel => '你的采集等级不足以采集此资源。';

  @override
  String startActivity(String name) {
    return '开始$name';
  }

  @override
  String get collect => '领取';

  @override
  String get abandon => '放弃';

  @override
  String get cancel => '取消';

  @override
  String get soFar => '目前已采集';

  @override
  String get activityDone => '已完成，可以领取';

  @override
  String remaining(String time) {
    return '剩余 $time';
  }

  @override
  String durHm(int h, String m) {
    return '$h 小时 $m 分';
  }

  @override
  String durM(int m) {
    return '$m 分';
  }

  @override
  String durS(int s) {
    return '$s 秒';
  }

  @override
  String get collectEarlyTitle => '提前领取？';

  @override
  String get collectEarlyBody => '你将获得目前已采集的物品，活动随之结束。';

  @override
  String get abandonTitle => '放弃活动？';

  @override
  String get abandonBody => '目前采集的物品都会丢失。';

  @override
  String get lootCollect => '太好了！';

  @override
  String get lootNothing => '这次什么也没有。';

  @override
  String xpGain(int xp) {
    return '+$xp 经验';
  }

  @override
  String get rarityCommon => '普通';

  @override
  String get rarityUncommon => '优秀';

  @override
  String get rarityRare => '稀有';

  @override
  String get rarityEpic => '史诗';

  @override
  String get rarityLegendary => '传说';

  @override
  String get rarityMythic => '神话';

  @override
  String get errSlotBusy => '已有其他活动在进行。';

  @override
  String get errNotActive => '此活动已结束。';

  @override
  String get evAchievement => '解锁成就！';

  @override
  String get evMasteryTier => '精通阶段提升！';

  @override
  String get evMarketSold => '你上架的物品已售出。';

  @override
  String get evMailReceived => '你有新邮件。';

  @override
  String get evPartyInvite => '你收到了队伍邀请。';

  @override
  String get evGuildInvite => '你收到了公会邀请。';

  @override
  String get navBag => '背包';

  @override
  String get equipment => '装备';

  @override
  String get bag => '背包';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap 千克';
  }

  @override
  String get slotHelm => '头盔';

  @override
  String get slotNecklace => '项链';

  @override
  String get slotChest => '胸甲';

  @override
  String get slotGloves => '手套';

  @override
  String get slotLegs => '腿甲';

  @override
  String get slotBoots => '靴子';

  @override
  String get slotMainhand => '主手';

  @override
  String get slotOffhand => '副手';

  @override
  String get slotRing => '戒指';

  @override
  String get slotTrinket => '饰品';

  @override
  String get equip => '装备';

  @override
  String get unequip => '卸下';

  @override
  String get sell => '出售';

  @override
  String get sortLabel => '排序';

  @override
  String get sortSlot => '背包顺序';

  @override
  String get sortRarity => '稀有度';

  @override
  String get sortType => '类型';

  @override
  String get sortName => '名称';

  @override
  String get sortValue => '价值';

  @override
  String selectedCount(int count) {
    return '已选择 $count 个';
  }

  @override
  String get selectAll => '全选';

  @override
  String get clearSelection => '清除选择';

  @override
  String sellSelectedTitle(int count) {
    return '出售 $count 件物品？';
  }

  @override
  String destroySelectedTitle(int count) {
    return '销毁 $count 件物品？';
  }

  @override
  String get deleteEntry => '删除';

  @override
  String craftHowMany(String name) {
    return '制作多少个$name？';
  }

  @override
  String craftCount(int count) {
    return '制作 ×$count';
  }

  @override
  String get destroy => '销毁';

  @override
  String sellTitle(String name) {
    return '出售$name？';
  }

  @override
  String get sellBody => '2 小时内可从商人处购回。';

  @override
  String destroyTitle(String name) {
    return '销毁$name？';
  }

  @override
  String get destroyBody => '此操作无法撤销。';

  @override
  String get soldFor => '售价';

  @override
  String get quantity => '数量';

  @override
  String get requiresLevel => '需要等级';

  @override
  String get weight => '重量';

  @override
  String weightGrams(int g) {
    return '$g 克';
  }

  @override
  String get sellValue => '出售价格';

  @override
  String get bound => '灵魂绑定';

  @override
  String get itemLocked => '此物品正被交易、邮件或上架占用。';

  @override
  String get errBagFull => '背包已满。';

  @override
  String get errItemLocked => '此物品已锁定。';

  @override
  String get errClassRestriction => '你的职业无法使用此物品。';

  @override
  String get errLevelRequirement => '你的等级太低。';

  @override
  String get errNotEquippable => '此物品无法装备。';

  @override
  String get navAdventure => '冒险';

  @override
  String get segHunt => '狩猎';

  @override
  String get segDungeon => '地下城';

  @override
  String get segTravel => '旅行';

  @override
  String get huntChoose => '选择猎物';

  @override
  String get startHunt => '开始狩猎';

  @override
  String get huntDone => '狩猎结束，领取战利品';

  @override
  String get huntEarlyBody => '你将保留已获胜的战斗，狩猎随之结束。';

  @override
  String mobStats(int hp, int atk, int def) {
    return '生命 $hp · 攻击 $atk · 防御 $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '胜 $won · 负 $lost';
  }

  @override
  String levelUp(int level) {
    return '升级！现在是 $level 级';
  }

  @override
  String levelReq(int level) {
    return '$level 级+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '$count 场战斗 · 首领 $boss';
  }

  @override
  String get dungeonEnter => '进入';

  @override
  String get dungeonDone => '探索完成，领取奖励';

  @override
  String get dungeonWaitHint => '探索结束后解锁奖励。';

  @override
  String get dungeonCleared => '通关！';

  @override
  String get dungeonFailed => '队伍全灭';

  @override
  String get bossKilled => '首领已击败';

  @override
  String get bonus => '通关奖励';

  @override
  String get youAreIn => '你位于';

  @override
  String get routes => '路线';

  @override
  String get noRoutes => '此处没有路线。';

  @override
  String get travelTo => '出发';

  @override
  String travelling(String place) {
    return '正在前往$place';
  }

  @override
  String get arrive => '到达';

  @override
  String get arrived => '你已到达';

  @override
  String get cancelTravel => '返回';

  @override
  String get kindCity => '城市';

  @override
  String get kindZone => '区域';

  @override
  String get kindWilderness => '荒野';

  @override
  String get kindDungeon => '地下城';

  @override
  String get kindContinent => '大陆';

  @override
  String get safeZone => '安全';

  @override
  String get errBusy => '已有事项在进行。';

  @override
  String get errNotReady => '尚未完成。';

  @override
  String get errWeeklyLocked => '你本周已通关此地下城。';

  @override
  String get errNoRoute => '没有直达那里的路线。';

  @override
  String get navWorkshop => '工坊';

  @override
  String get segCraft => '制作';

  @override
  String get segShop => '商店';

  @override
  String get profBlacksmithing => '锻造';

  @override
  String get profAlchemy => '炼金';

  @override
  String get profCooking => '烹饪';

  @override
  String get learn => '学习';

  @override
  String recipeCount(int count) {
    return '$count 个配方';
  }

  @override
  String get itemType => '类型';

  @override
  String get stackSize => '堆叠';

  @override
  String get skillXpReward => '技能经验';

  @override
  String get learnHint => '学习此专业以解锁配方。';

  @override
  String get craft => '制作';

  @override
  String crafted(String name, int count) {
    return '已制作 $name ×$count';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession现在是 $level 级';
  }

  @override
  String requiresSkill(int level) {
    return '需要技能 $level';
  }

  @override
  String get buy => '购买';

  @override
  String bought(String name) {
    return '已购买$name';
  }

  @override
  String get buybacks => '购回';

  @override
  String get buybackHint => '已售出的物品可在 2 小时内购回。';

  @override
  String get errInsufficientGold => '金币不足。';

  @override
  String get errInsufficientMaterials => '材料不足。';

  @override
  String get errSkillTooLow => '你的技能不足以使用此配方。';

  @override
  String get errProfessionNotLearned => '请先学习该专业。';

  @override
  String get errPriceChanged => '价格已变动，请确认后重试。';

  @override
  String get errCannotSell => '此物品无法出售。';

  @override
  String get errBuybackExpired => '购回期限已过。';

  @override
  String get segMarket => '市场';

  @override
  String get classWeapon => '武器';

  @override
  String get classArmor => '护甲';

  @override
  String get classConsumable => '消耗品';

  @override
  String get classMaterial => '材料';

  @override
  String get filterAll => '全部';

  @override
  String get sortCheapest => '最便宜';

  @override
  String get sortNewest => '最新';

  @override
  String get marketEmpty => '没有符合条件的商品。';

  @override
  String listingLine(String seller, String price) {
    return '$seller · 单价 $price';
  }

  @override
  String get loadMore => '加载更多';

  @override
  String myListings(int count) {
    return '我的上架（$count）';
  }

  @override
  String get cancelListing => '下架';

  @override
  String cancelListingTitle(String name) {
    return '下架$name？';
  }

  @override
  String get cancelListingBody => '物品将返回背包。上架费仅在前 6 小时内部分退还。';

  @override
  String buyListingTitle(String name) {
    return '购买$name？';
  }

  @override
  String buyListingBody(String price, String seller) {
    return '向$seller支付 $price。';
  }

  @override
  String get listOnMarket => '上架到市场';

  @override
  String listed(String name) {
    return '$name已上架市场';
  }

  @override
  String quantityOf(int count, int max) {
    return '数量：$count / $max';
  }

  @override
  String get unitPrice => '单价';

  @override
  String suggestedPrice(String price) {
    return '市场价：$price';
  }

  @override
  String get listingFee => '上架费（不可退还）';

  @override
  String get youReceive => '售出后到手';

  @override
  String get errPriceOutOfRange => '该价格超出允许范围。';

  @override
  String get errListingCap => '你的在售商品过多。';

  @override
  String get errSelfBuy => '不能购买自己的商品。';

  @override
  String get errListingUnavailable => '此商品已不可用。';

  @override
  String get errNotTradable => '此物品无法交易。';

  @override
  String get social => '社交';

  @override
  String get tabChat => '聊天';

  @override
  String get tabFriends => '好友';

  @override
  String get tabMail => '邮件';

  @override
  String get chanGlobal => '世界';

  @override
  String get chanGuild => '公会';

  @override
  String get chanParty => '队伍';

  @override
  String get chatEmpty => '还没有消息，打个招呼吧！';

  @override
  String chatHint(String channel) {
    return '发送到$channel';
  }

  @override
  String get send => '发送';

  @override
  String get addFriend => '添加';

  @override
  String friendAdded(String name) {
    return '你和$name已成为好友';
  }

  @override
  String friendRequested(String name) {
    return '已向$name发送好友请求';
  }

  @override
  String get friendRequests => '好友请求';

  @override
  String friendsCount(int count) {
    return '好友（$count）';
  }

  @override
  String get noFriends => '还没有好友，通过名称添加吧。';

  @override
  String get more => '更多';

  @override
  String get whisper => '私聊';

  @override
  String get remove => '移除';

  @override
  String get block => '屏蔽';

  @override
  String get unblock => '取消屏蔽';

  @override
  String get blockedList => '已屏蔽';

  @override
  String removeFriendTitle(String name) {
    return '移除$name？';
  }

  @override
  String get removeFriendBody => '之后可以重新发送好友请求。';

  @override
  String blockTitle(String name) {
    return '屏蔽$name？';
  }

  @override
  String get blockBody => '对方无法私聊你，其聊天内容会被隐藏，邮件也会被丢弃，好友关系同时解除。';

  @override
  String get accept => '接受';

  @override
  String get decline => '拒绝';

  @override
  String get newMail => '写邮件';

  @override
  String get inboxEmpty => '邮箱是空的。';

  @override
  String get noSubject => '（无主题）';

  @override
  String mailFrom(String name) {
    return '来自$name';
  }

  @override
  String get attachments => '附件';

  @override
  String get claim => '领取附件';

  @override
  String get mailClaimed => '附件已移入背包和钱包';

  @override
  String get deleteMailTitle => '删除这封邮件？';

  @override
  String get deleteMailBody => '删除的邮件无法恢复。';

  @override
  String mailSent(String postage) {
    return '邮件已发送（邮资 $postage）';
  }

  @override
  String get mailTo => '收件人（角色名）';

  @override
  String get mailSubject => '主题';

  @override
  String get mailBody => '内容';

  @override
  String get attachGold => '附加金币';

  @override
  String get errAlreadyFriends => '你们已经是好友。';

  @override
  String get errAlreadyRequested => '好友请求正在等待处理。';

  @override
  String get errBlocked => '你无法与此玩家互动。';

  @override
  String get errFriendCap => '好友列表已满。';

  @override
  String get errTargetNotFound => '没有这个名称的角色。';

  @override
  String get errSelfTarget => '不能以自己为目标。';

  @override
  String get errInboxFull => '对方邮箱已满。';

  @override
  String get errTooLong => '内容太长。';

  @override
  String get errEmpty => '请先输入内容。';

  @override
  String get errNotInChannel => '你不在该群组中。';

  @override
  String get errMuted => '你已被禁言，暂时无法聊天。';

  @override
  String get errHasAttachments => '删除前请先领取附件。';

  @override
  String get useItem => '使用';

  @override
  String restoresHp(int amount) {
    return '恢复 $amount 生命';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return '已使用$name · 生命 $hp/$max';
  }

  @override
  String get errHpFull => '生命值已满。';

  @override
  String get errInCombat => '请先结束狩猎或地下城探索。';

  @override
  String get errNotUsable => '此物品无法使用。';

  @override
  String get errTooWounded => '你伤势过重无法战斗。请休息或喝药水。';

  @override
  String get forgotPassword => '忘记密码？';

  @override
  String get resetHint => '输入账号邮箱，我们会发送设置新密码的链接。';

  @override
  String get resetSent => '如果该邮箱存在账号，重置链接已发出，1 小时后失效。';

  @override
  String get sendLink => '发送链接';

  @override
  String get close => '关闭';

  @override
  String get changePassword => '修改密码';

  @override
  String get currentPassword => '当前密码';

  @override
  String get newPassword => '新密码';

  @override
  String get save => '保存';

  @override
  String get passwordChanged => '密码已修改，其他设备已退出登录。';

  @override
  String get verifyEmailBanner => '请确认你的邮箱地址。';

  @override
  String get resend => '重新发送';

  @override
  String get verificationSent => '验证邮件已发送。';

  @override
  String get errEmailNotVerified => '创建角色前请先确认邮箱。';

  @override
  String get errAlreadyVerified => '你的邮箱已确认。';

  @override
  String get tabParty => '队伍';

  @override
  String get tabGuild => '公会';

  @override
  String get invite => '邀请';

  @override
  String invitedBy(String name) {
    return '$name邀请了你';
  }

  @override
  String get noParty => '你不在队伍中。邀请他人即可组队（2–5 人）。';

  @override
  String get createParty => '创建队伍';

  @override
  String get pendingInvites => '邀请';

  @override
  String partyMembers(int count, int cap) {
    return '队伍（$count/$cap）';
  }

  @override
  String get makeLeader => '设为队长';

  @override
  String get kick => '移出';

  @override
  String get leave => '离开';

  @override
  String get disband => '解散';

  @override
  String get leavePartyTitle => '离开队伍？';

  @override
  String get leavePartyBody => '如果你是队长，由资历最久的成员接任。只剩一人时队伍解散。';

  @override
  String get disbandPartyTitle => '解散队伍？';

  @override
  String get disbandPartyBody => '所有人都会离开队伍。';

  @override
  String get rankLeader => '会长';

  @override
  String get rankOfficer => '干部';

  @override
  String get rankVeteran => '资深成员';

  @override
  String get rankMember => '成员';

  @override
  String get rankInitiate => '新人';

  @override
  String get createGuild => '创建公会';

  @override
  String createGuildBody(String fee) {
    return '创建公会需要 $fee。';
  }

  @override
  String createGuildFee(String fee) {
    return '创建公会（$fee）';
  }

  @override
  String get guildName => '公会名称';

  @override
  String get guildTag => '标签';

  @override
  String get guildTagHint => '2–5 个字母或数字';

  @override
  String guildMembers(int count, int cap) {
    return '成员 $count/$cap';
  }

  @override
  String disbandScheduled(String time) {
    return '公会将在 $time 后解散';
  }

  @override
  String get motd => '每日公告';

  @override
  String get editMotd => '编辑每日公告';

  @override
  String get setRank => '更改职位';

  @override
  String setRankFor(String name) {
    return '$name的新职位';
  }

  @override
  String transferTitle(String name) {
    return '让$name成为会长？';
  }

  @override
  String get transferBody => '你将成为干部，且无法自行撤销。';

  @override
  String kickTitle(String name) {
    return '将$name移出公会？';
  }

  @override
  String get leaveGuildTitle => '离开公会？';

  @override
  String get leaveGuildBody => '重新加入需要新的邀请。';

  @override
  String get disbandGuildTitle => '解散公会？';

  @override
  String get disbandGuildBody => '公会将在 48 小时后解散，在此之前可以取消。';

  @override
  String get cancelDisband => '取消解散';

  @override
  String get errAlreadyInGroup => '对方已在群组中。';

  @override
  String get errNotLeader => '只有首领可以这样做。';

  @override
  String get errGroupFull => '群组已满。';

  @override
  String get errInviteInvalid => '邀请已过期。';

  @override
  String get errAlreadyInvited => '邀请正在等待处理。';

  @override
  String get errGuildForbidden => '你的职位无法这样做。';

  @override
  String get errLeaderCannotLeave => '请先移交会长或解散公会。';

  @override
  String get errInvalidGuildName => '3–24 个字母或数字。';

  @override
  String get errInvalidTag => '标签必须是 2–5 个字母或数字。';

  @override
  String get errTagTaken => '该标签已被占用。';

  @override
  String get progress => '成就与精通';

  @override
  String get tabAchievements => '成就';

  @override
  String get tabMastery => '精通';

  @override
  String get tabReputation => '声望';

  @override
  String get tabTitles => '称号';

  @override
  String achievementPoints(int points) {
    return '$points 成就点数';
  }

  @override
  String pointsShort(int points) {
    return '$points 点';
  }

  @override
  String rewardTitle(String title) {
    return '称号：$title';
  }

  @override
  String masteryTier(int tier, int max) {
    return '阶段 $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return '今日 +$amount';
  }

  @override
  String get noTitle => '无称号';

  @override
  String get noTitlesYet => '称号通过成就解锁。';

  @override
  String get errTitleNotOwned => '你尚未解锁该称号。';

  @override
  String get notifications => '通知';

  @override
  String get markAllRead => '全部标为已读';

  @override
  String get inboxNothing => '还没有通知。';

  @override
  String get notifSettings => '通知设置';

  @override
  String get muteAll => '暂停推送通知';

  @override
  String get muteHint => '暂停 24 小时，收件箱仍会接收。';

  @override
  String mutedFor(String time) {
    return '已暂停 $time';
  }

  @override
  String get modePushInbox => '推送 + 收件箱';

  @override
  String get modeInbox => '仅收件箱';

  @override
  String get modeSilent => '静音';

  @override
  String get modeOff => '关闭';

  @override
  String get nkActivity => '活动已结束';

  @override
  String get nkDungeon => '地下城探索已结束';

  @override
  String get nkTravel => '已到达';

  @override
  String get nkWorldEvent => '世界事件已开始';

  @override
  String get nkMarketSold => '物品已售出';

  @override
  String get nkMail => '新邮件';

  @override
  String get nkPartyInvite => '队伍邀请';

  @override
  String get nkGuildInvite => '公会邀请';

  @override
  String get nkAchievement => '成就已解锁';

  @override
  String get nkSystem => '公告';

  @override
  String timeAgo(String time) {
    return '$time前';
  }

  @override
  String get segEvents => '活动';

  @override
  String get eventsEmpty => '目前没有安排世界事件。';

  @override
  String get worldBoss => '世界首领';

  @override
  String get invasion => '入侵';

  @override
  String startsIn(String time) {
    return '$time后开始';
  }

  @override
  String endsIn(String time) {
    return '$time后结束';
  }

  @override
  String regionOnly(String region) {
    return '仅限$region';
  }

  @override
  String get eventGoal => '目标';

  @override
  String get enlist => '报名';

  @override
  String get enlistHint => '报名时会锁定你当前的战力。活动进行的每一分钟都会计分，结束时发放奖励。';

  @override
  String enlistedPower(int power) {
    return '已报名 · 战力 $power';
  }

  @override
  String myContribution(int score, int tier) {
    return '你的积分 $score · 奖励阶段 $tier';
  }

  @override
  String get eventSucceeded => '胜利';

  @override
  String get eventFailed => '失败';

  @override
  String tokenBalance(int tokens) {
    return '$tokens 活动代币';
  }

  @override
  String claimDaily(int tokens) {
    return '领取 $tokens';
  }

  @override
  String get claimedToday => '已领取';

  @override
  String eventTokensGained(int tokens) {
    return '+$tokens 活动代币';
  }

  @override
  String get errEventNotOpen => '此活动的报名已关闭。';

  @override
  String get errWrongRegion => '请先前往活动所在地区。';

  @override
  String get errAlreadyEnlisted => '你已报名。';

  @override
  String get errAlreadyClaimed => '今天已领取。';

  @override
  String get errEventNotActive => '此活动未在进行。';

  @override
  String gatherYield(String activity) {
    return '$activity产量';
  }

  @override
  String perHour(String value) {
    return '$value/小时';
  }

  @override
  String get quests => '任务';

  @override
  String get questsInProgress => '进行中';

  @override
  String get questsNoneOpen => '没有进行中的任务，请在下方选择。';

  @override
  String get questsStory => '主线';

  @override
  String get questsSide => '支线任务';

  @override
  String get questsDaily => '每日';

  @override
  String resetsIn(String time) {
    return '$time后重置';
  }

  @override
  String get turnIn => '提交';

  @override
  String questRewarded(String rewards) {
    return '任务完成：$rewards';
  }

  @override
  String get abandonQuestTitle => '放弃任务？';

  @override
  String get abandonQuestBody => '进度会丢失，之后可以重新接取。';

  @override
  String get abandonDailyBody => '进度会丢失，今天无法再次接取此每日任务。';

  @override
  String objGather(String activity) {
    return '采集（$activity）';
  }

  @override
  String get objKillAny => '击败敌人';

  @override
  String objKill(String mob) {
    return '击败$mob';
  }

  @override
  String get objBoss => '击败首领';

  @override
  String get objDungeonAny => '通关地下城';

  @override
  String objDungeon(String dungeon) {
    return '通关$dungeon';
  }

  @override
  String get objCraftAny => '制作物品';

  @override
  String objCraft(String item) {
    return '制作$item';
  }

  @override
  String get evQuestComplete => '任务目标已完成，快去提交！';

  @override
  String get errQuestLocked => '请先完成上一个主线任务。';

  @override
  String get errQuestTaken => '你已接取此任务。';

  @override
  String get errQuestNotOffered => '此每日任务不在今天的任务板上。';

  @override
  String get errQuestCap => '进行中的任务过多，请完成或放弃一个。';

  @override
  String get errQuestNotActive => '此任务未开启。';

  @override
  String get errQuestIncomplete => '目标尚未完成。';

  @override
  String get skills => '技能';

  @override
  String get skillActive => '主动';

  @override
  String get skillPassive => '被动';

  @override
  String skillStrike(int power) {
    return '$power% 伤害';
  }

  @override
  String skillHeal(int power) {
    return '生命低于一半时恢复 $power% 生命';
  }

  @override
  String skillMp(int mp) {
    return '$mp 法力';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: '每 $turns 回合', one: '每回合');
    return '$_temp0';
  }

  @override
  String get termsOfService => '服务条款';

  @override
  String get privacyPolicy => '隐私政策';

  @override
  String get accountMenu => '账号';

  @override
  String get downloadMyData => '下载我的数据';

  @override
  String get exportOpened => '下载链接已在浏览器中打开（15 分钟内有效）。';

  @override
  String get linkOpenFailed => '无法打开链接。';

  @override
  String get deleteAccount => '删除我的账号';

  @override
  String get deleteAccountBody => '你的账号和所有角色将在 30 天后永久删除。在此期间重新登录可取消删除。请输入密码确认。';

  @override
  String get deleteAccountConfirm => '删除账号';

  @override
  String deleteScheduled(String date) {
    return '你的账号将于 $date 删除。如需取消，请在此之前登录。';
  }

  @override
  String get deletionCancelled => '账号删除请求已取消。';

  @override
  String get support => '客服';

  @override
  String get supportIntro => '发现错误或账号需要帮助？给我们留言，团队会在这里回复，并向你发送通知。';

  @override
  String get supportNew => '新请求';

  @override
  String get supportNone => '还没有请求。';

  @override
  String get supportSent => '请求已收到。有回复时会通知你。';

  @override
  String get supportCategory => '主题';

  @override
  String get supportSubject => '标题';

  @override
  String get supportMessage => '发生了什么？尽量逐步说明。';

  @override
  String get supportTooShort => '标题至少 3 个字符，内容至少 10 个字符。';

  @override
  String get supportBug => '错误';

  @override
  String get supportAccount => '账号';

  @override
  String get supportAbuse => '滥用 / 举报玩家';

  @override
  String get supportFeedback => '建议';

  @override
  String get supportOther => '其他';

  @override
  String get supportOpen => '待处理';

  @override
  String get supportAnswered => '已回复';

  @override
  String get supportClosed => '已关闭';

  @override
  String get supportReply => '团队回复';

  @override
  String get nkSupport => '客服回复';

  @override
  String get errServerBusy => '服务器繁忙，请几秒后重试。';

  @override
  String get notifTabInbox => '通知';

  @override
  String get notifTabRewards => '奖励';

  @override
  String get rewardsNone => '还没有奖励。采集或战斗获得的一切都会列在这里（保留 30 天）。';

  @override
  String rewardFromActivity(String name) {
    return '采集：$name';
  }

  @override
  String rewardFromHunt(String name) {
    return '狩猎：$name';
  }

  @override
  String rewardFromDungeon(String name) {
    return '地下城：$name';
  }

  @override
  String rewardLevelUp(int level) {
    return '$level 级！';
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
  String get raidNeedParty => '团队副本需要队伍。可在下方排队随机组队、加入招募或发布自己的招募。';

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

  @override
  String get mapHint => '你可以直接前往等级允许的任何区域，时间按最短道路计算。拖动或缩放地图，点击区域。';

  @override
  String get mapNoRoad => '没有直达这里的道路。请经由你身边发光的区域前往。';

  @override
  String get youAreHere => '你在这里。';

  @override
  String travelTime(String time) {
    return '旅行时间：$time';
  }

  @override
  String xpPerKill(int xp) {
    return '$xp 经验';
  }

  @override
  String get deleteCharacter => '删除角色';

  @override
  String deleteCharacterBody(String name) {
    return '$name 将连同所有物品和金币立即删除，且无法恢复。请输入角色名以确认。';
  }

  @override
  String get characterNameConfirm => '角色名';

  @override
  String characterDeletesAt(String when) {
    return '删除时间：$when';
  }

  @override
  String get cancelCharacterDelete => '取消删除';

  @override
  String characterDeleted(String name) {
    return '$name 已删除。';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name 不会被删除。';
  }

  @override
  String get errNameMismatch => '名称不匹配。';

  @override
  String attrGainPerPoint(String gains) {
    return '每点：$gains';
  }

  @override
  String attrGainPending(String gains) {
    return '将增加：$gains';
  }

  @override
  String get attrNoEffect => '对你的职业无效';

  @override
  String get attrInputHint => '输入每项属性要分配的点数';

  @override
  String get gainHp => '生命';

  @override
  String get gainMp => '法力';

  @override
  String get gainHpRegen => '生命/秒';

  @override
  String get huntSupplies => '携带补给';

  @override
  String get huntSuppliesHint => '战斗间生命低于一半时喝药水，增益结束时吃食物，法力不足时喝法力药水。剩余的会回到背包。';

  @override
  String get supplyHp => '治疗药水';

  @override
  String get supplyMp => '法力药水';

  @override
  String get supplyMpHint => '适用于使用法力、能量、专注、自然能量或灵力的职业';

  @override
  String get supplyFood => '食物';

  @override
  String get huntNoSupplies => '背包里没有药水或食物：可在商店购买或自行制作。';

  @override
  String restoresMpPct(int pct) {
    return '恢复 $pct% 法力';
  }

  @override
  String get huntSupplyOnly => '仅可作为狩猎补给使用';

  @override
  String get errBadSupply => '该物品无法带去狩猎。';

  @override
  String get errNotEnoughSupplies => '背包中该补给不足。';

  @override
  String get errHuntOnly => '请将其作为狩猎补给带上。';

  @override
  String get huntFled => '生命低于 5%：你撤离了战斗';

  @override
  String suppliesUsed(String list) {
    return '已使用：$list';
  }

  @override
  String get compareTitle => '与当前装备对比';

  @override
  String compareWorn(String name) {
    return '当前装备：$name';
  }

  @override
  String get compareEmpty => '该栏位为空：装备后将获得全部属性。';

  @override
  String get compareColStat => '属性';

  @override
  String get compareColWorn => '当前';

  @override
  String get compareColThis => '此件';

  @override
  String get compareColDiff => '变化';

  @override
  String get compareSame => '属性与当前装备相同。';

  @override
  String get wearOk => '你的职业可以装备';

  @override
  String get wearClass => '你的职业无法装备';

  @override
  String wearLevel(int level) {
    return '需要等级 $level';
  }

  @override
  String get wearAttribute => '属性不足';

  @override
  String get wearLegendClass => '职业不可用';

  @override
  String get wearLegendLater => '暂不可用（等级或属性）';

  @override
  String get enchantAction => '附魔';

  @override
  String get enchantPick => '选择要附魔的装备';

  @override
  String get enchantNone => '没有适合此卷轴的装备。';

  @override
  String enchantFits(String groups) {
    return '适用：$groups';
  }

  @override
  String get enchantGroupWeapon => '武器';

  @override
  String get enchantGroupArmor => '护甲';

  @override
  String get enchantGroupJewelry => '饰品';

  @override
  String get enchantGroupAny => '任意装备';

  @override
  String get enchantReplaceTitle => '替换附魔？';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old将消失，由$next取代。';
  }

  @override
  String enchantDone(String item) {
    return '$item已附魔';
  }

  @override
  String enchantOn(String name) {
    return '附魔：$name';
  }

  @override
  String get wornTag => '已装备';

  @override
  String get errNotEnchant => '这不是附魔卷轴。';

  @override
  String get errEnchantSlot => '该附魔不适用于此装备。';

  @override
  String get statGatherMining => '采矿产量';

  @override
  String get statGatherHerbalism => '采药产量';

  @override
  String get statGatherWoodcutting => '伐木产量';

  @override
  String get statGatherFishing => '钓鱼产量';

  @override
  String get statGatherHunting => '狩猎产量';

  @override
  String get statCraftSpeed => '制作速度';

  @override
  String diffXp(int pct) {
    return '$pct% 经验';
  }

  @override
  String get diffNoXp => '无经验';

  @override
  String get diffLegend => '配方颜色：按你的技能获得的经验';

  @override
  String get recipeGroupEnchantWeapon => '武器附魔';

  @override
  String get recipeGroupEnchantArmor => '护甲附魔';

  @override
  String get recipeGroupEnchantJewelry => '饰品附魔';

  @override
  String get recipeGroupEnchantProfession => '专业附魔';

  @override
  String get partyHint => '同时进行相同工作（同一采集、职业、狩猎或地下城）的成员，每位参与成员使每人获得 +3% 经验：五人同时为 +15%。不会分走任何人的经验。';

  @override
  String get partyChat => '队伍聊天';

  @override
  String partyShared(String names, int kept, int earned) {
    return '与$names的队伍加成：$earned → $kept 经验';
  }

  @override
  String evRaidMatched(String raid) {
    return '已为$raid找到队伍，请查看团队副本。';
  }

  @override
  String get nkRaidMatched => '团队匹配成功';

  @override
  String get workHunt => '狩猎';

  @override
  String get workDungeon => '地下城';

  @override
  String get raidFind => '寻找队伍';

  @override
  String get raidQueueJoin => '排队';

  @override
  String get raidQueueLeave => '离开队列';

  @override
  String raidQueued(String raid, int waiting) {
    return '正在排队：$raid · $waiting 人等待';
  }

  @override
  String get raidQueueHint => '等待满 5 人，或 5 分钟后至少 2 人时，会自动随机组队。';

  @override
  String get raidAdvertise => '发布招募';

  @override
  String get raidUnlist => '撤销招募';

  @override
  String raidListed(String raid) {
    return '你的队伍正在招募：$raid';
  }

  @override
  String get raidBoard => '招募列表';

  @override
  String get raidBoardEmpty => '暂无招募。';

  @override
  String get raidJoinParty => '加入队伍';

  @override
  String get raidPickRaid => '团队副本';

  @override
  String get errListingGone => '该招募已结束或已满。';

  @override
  String get workGathering => '采集';

  @override
  String get workAdventure => '冒险';

  @override
  String get workRaid => '团队副本';

  @override
  String get workTravel => '旅行';

  @override
  String get workCrafting => '制作';

  @override
  String get workIdle => '空闲';

  @override
  String get errDungeonRegion => '请先在地图上前往此处。';

  @override
  String get dungeonsElsewhere => '旅行时会开放更多地下城和团队副本：每个都需从地图上的所在区域进入。';

  @override
  String get dungeonNoneHere => '此处没有可进入的地下城，请在地图上前往地下城区域。';

  @override
  String dungeonXp(int xp) {
    return '你约获得 $xp 经验';
  }

  @override
  String appVersionLine(String version) {
    return '版本 $version';
  }

  @override
  String get updateButton => '更新';

  @override
  String updateAvailable(String version) {
    return '新版本可用：$version';
  }

  @override
  String get updateLatest => '已是最新版本';

  @override
  String updateConfirm(String version, int size) {
    return '下载并安装 $version 版本（$size MB）？';
  }

  @override
  String get updateDownloading => '正在下载…';

  @override
  String get updatePermission => '请允许 MeNa World 安装应用，然后再次点击更新。';

  @override
  String get updateFailed => '无法获取更新，请检查网络连接。';

  @override
  String get fightReport => '战斗报告';

  @override
  String get endTimeUp => '时间到：一直持续到计划结束。';

  @override
  String get endEarly => '提前领取：在计划结束前停止。';

  @override
  String endFled(String mob) {
    return '对战$mob时生命低于5%：你撤退了。';
  }

  @override
  String endDefeat(String mob) {
    return '你被$mob击败。';
  }

  @override
  String endPartyDefeat(String mob) {
    return '队伍被$mob击败。';
  }

  @override
  String endTimeout(String mob) {
    return '与$mob的战斗超过回合上限：已停止。';
  }

  @override
  String get endCleared => '已通关：所有波次和首领都被击败。';

  @override
  String get repCampHint => '生命低于50%时，猎人会在下一场战斗前扎营回满，因此狩猎很少降到5%。';

  @override
  String get repYou => '你';

  @override
  String get repLevel => '等级';

  @override
  String get repHp => '生命';

  @override
  String get repAttack => '攻击';

  @override
  String get repDefense => '防御';

  @override
  String get repEvasion => '闪避';

  @override
  String get repHitChance => '命中率';

  @override
  String get repCrit => '暴击';

  @override
  String repFights(int n) {
    return '$n 场战斗';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '命中 $hits/$swings（$pct%），平均 $avg，最高 $max';
  }

  @override
  String get repYourSwings => '你的攻击';

  @override
  String get repPartySwings => '队伍的攻击';

  @override
  String get repEnemySwings => '敌人的攻击';

  @override
  String repLowest(int pct) {
    return '战斗中最低生命：$pct%';
  }

  @override
  String repRest(String time) {
    return '扎营休息：$time';
  }

  @override
  String get restart => '重新开始';

  @override
  String get restartTitle => '领取并重新开始？';

  @override
  String get restartBody => '目前的收获放入背包，并以相同时长重新开始。';

  @override
  String get dailyTitle => '每日登录奖励';

  @override
  String get dailyBody => '每天可用任意角色领取一次。连续每周收益递增，第 29 天最高；中断一天将重新开始。';

  @override
  String get dailyClaim => '领取';

  @override
  String get dailyTomorrow => '明天再来';

  @override
  String dailyDay(int n) {
    return '第$n天';
  }

  @override
  String get dailyOpen => '每日奖励';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name 因说脏话被禁言 $min 分钟。';
  }

  @override
  String get storyTitle => '剧情';

  @override
  String get storyOpen => '剧情模式';

  @override
  String get storyComingSoon => '剧情模式即将推出。章节写好后会在这里出现。';

  @override
  String storyLocked(int level) {
    return '$level 级开放';
  }

  @override
  String get bankTitle => '银行';

  @override
  String get bankHint => '所有角色共享。点击背包物品存入，点击银行物品取出。绑定物品留在角色身上。';

  @override
  String get bankEmpty => '银行是空的。';

  @override
  String get bankDeposit => '存入';

  @override
  String get bankWithdraw => '取出';

  @override
  String get errBankFull => '银行已满。';

  @override
  String get tabTrade => '交易';

  @override
  String get tradeHint => '与其他玩家交易物品。双方放入物品并确认；每方按所给物品市场均价的8%（按件）缴税。';

  @override
  String get tradeWith => '玩家名称';

  @override
  String get tradeOpen => '发起交易';

  @override
  String tradeTitle(String name) {
    return '与$name交易';
  }

  @override
  String get tradeYours => '你给出';

  @override
  String tradeTheirs(String name) {
    return '$name给出';
  }

  @override
  String get tradeNothing => '暂无物品。';

  @override
  String get tradeAdd => '添加物品';

  @override
  String get tradeRemove => '撤回';

  @override
  String get tradeTax => '你的税费';

  @override
  String get tradeReady => '已确认';

  @override
  String get tradeWaiting => '未确认';

  @override
  String get tradeConfirm => '确认交易';

  @override
  String get tradeWaitingThem => '等待对方确认';

  @override
  String get tradeCancel => '取消交易';

  @override
  String get tradeCancelBody => '不会交换任何物品。';

  @override
  String get tradeDone => '交易完成。';

  @override
  String get errTradeBusy => '其中一方已在交易中。';

  @override
  String get errTradeClosed => '该交易已结束。';

  @override
  String get errTradeItem => '该物品现在无法交易。';

  @override
  String get errTradeEmpty => '请先放入物品。';

  @override
  String get errTradeTooMany => '交易物品过多。';

  @override
  String get errCraftLevel => '角色等级低于该物品等级。';

  @override
  String get gatherNoneHere => '你的资源在其他地图：前往旁边标注的区域。';

  @override
  String get errGatherRegion => '该资源在其他地图。';

  @override
  String get helpersOpen => '助手';

  @override
  String get helpersTitle => '助手';

  @override
  String get helpersHint => '雇佣助手与你一起采集和制造。他们以你一半的速度工作，获得一半经验并计入你的职业；每种助手有独特加成。';

  @override
  String get helpersGather => '采集';

  @override
  String get helpersCraft => '制造';

  @override
  String get helperHire => '雇佣';

  @override
  String helperHireBody(String name, String price) {
    return '以$price雇佣$name？';
  }

  @override
  String get helperSend => '派去工作';

  @override
  String get helperCancel => '召回';

  @override
  String get helperCancelBody => '采集任务无收获；制造任务退还材料和金币。';

  @override
  String get helperDone => '已完成，可领取';

  @override
  String helperWorking(String time) {
    return '工作中，剩余$time';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return '$items件物品，$profession经验+$xp';
  }

  @override
  String get helperNoRecipes => '该职业暂无适合你等级的配方。';

  @override
  String get helperPickRecipe => '选择配方';

  @override
  String helperBonusYield(int pct) {
    return '产量+$pct%';
  }

  @override
  String helperBonusSpeed(int pct) {
    return '速度+$pct%';
  }

  @override
  String helperBonusXp(int pct) {
    return '经验+$pct%';
  }

  @override
  String helperBonusGold(int pct) {
    return '金币消耗-$pct%';
  }

  @override
  String helperBonusDouble(int pct) {
    return '$pct%几率双倍产出';
  }

  @override
  String get helperMiner => '矿工';

  @override
  String get helperWoodsman => '伐木工';

  @override
  String get helperAngler => '渔夫';

  @override
  String get helperHerbalist => '草药师';

  @override
  String get helperTrapper => '陷阱猎人';

  @override
  String get helperSmith => '铁匠';

  @override
  String get helperTanner => '制革匠';

  @override
  String get helperWeaver => '织工';

  @override
  String get helperCarver => '木雕匠';

  @override
  String get helperApothecary => '药剂师';

  @override
  String get helperCook => '厨师';

  @override
  String get helperSmelter => '冶炼工';

  @override
  String get helperScribe => '抄写员';

  @override
  String get helperGemcutter => '宝石匠';

  @override
  String get errHelperBusy => '该助手正在工作。';

  @override
  String get errHelperNotReady => '助手尚未完成。';

  @override
  String get shopOpen => '商店';

  @override
  String get shopTitle => '点券商店';

  @override
  String get shopCredits => '点券';

  @override
  String get shopAvatars => '头像';

  @override
  String get shopBackgrounds => '背景';

  @override
  String shopBalance(int n) {
    return '$n 点券';
  }

  @override
  String get shopBuyCredits => '购买点券以解锁加成和外观。';

  @override
  String get shopPaymentsSoon => '现金购买点券即将推出。';

  @override
  String get shopBoosts => '经验加成';

  @override
  String get shopBuy => '购买';

  @override
  String shopBuyBody(int n) {
    return '花费 $n 点券？';
  }

  @override
  String shopPrice(int n) {
    return '$n 点';
  }

  @override
  String get boostGather => '采集';

  @override
  String get boostAdventure => '战斗';

  @override
  String get boostCraft => '职业';

  @override
  String get boostDesc => '2小时经验+100%';

  @override
  String get boostShort => '经验+100%，2小时';

  @override
  String boostLeft(String time) {
    return '生效中，剩余$time';
  }

  @override
  String boostWatchAd(int n) {
    return '看广告（剩$n次）';
  }

  @override
  String get boostAdsSoon => '广告即将推出';

  @override
  String get errInsufficientCredits => '点券不足。';

  @override
  String get errBoostActive => '该加成已生效。';

  @override
  String get errAdLimit => '今日广告次数已用完。';

  @override
  String helperCanMake(int n) {
    return '可制作：$n';
  }

  @override
  String get helperPickResource => '采集目标';

  @override
  String get helperCraftHint => '材料和金币现在扣除；期间你可以继续游戏。';

  @override
  String get petsOpen => '宠物';

  @override
  String get petsHint => '通关400级以上的团队副本时会掉落宠物蛋。蛋在孵化器中孵化；宠物通过自己的狩猎升级，出战宠物为你提供加成。 出战时它还会在战斗中协助你，为每次命中增加伤害。';

  @override
  String get petsEggs => '宠物蛋';

  @override
  String get petsNoEggs => '没有宠物蛋。通关400级以上的团队副本来获取。';

  @override
  String petsMine(int n, int max) {
    return '我的宠物（$n/$max）';
  }

  @override
  String get petHatch => '孵化';

  @override
  String petHatchBody(String egg, String time) {
    return '将$egg放入孵化器？$time后孵化。';
  }

  @override
  String get petHatching => '孵化中';

  @override
  String get petHatched => '已孵化！';

  @override
  String petPower(int n) {
    return '战斗力相当于$n级怪物';
  }

  @override
  String get petActive => '出战中';

  @override
  String get petActivate => '设为出战';

  @override
  String get petDeactivate => '休息';

  @override
  String get petHunt => '派去狩猎';

  @override
  String petHunting(String mob) {
    return '正在狩猎$mob';
  }

  @override
  String get petHuntDone => '狩猎结束';

  @override
  String get petCollectEarly => '现在结束只结算已完成的战斗。';

  @override
  String petCollected(int wins, int xp) {
    return '$wins场胜利，宠物经验+$xp';
  }

  @override
  String petLevelUp(String name, int level) {
    return '$name升到了$level级！';
  }

  @override
  String get petDied => '宠物战败，狩猎提前结束。';

  @override
  String get petPickMob => '狩猎目标';

  @override
  String petXpPerWin(int n) {
    return '每场胜利+$n经验';
  }

  @override
  String get errPetCap => '宠物数量已达上限。';

  @override
  String get errPetHatching => '宠物尚未孵化。';

  @override
  String get errPetHunting => '宠物已在狩猎中。';

  @override
  String petStageLevels(int from, int to) {
    return '$from–$to级';
  }

  @override
  String get regionMobs => '怪物';

  @override
  String get regionDungeons => '地下城';

  @override
  String get regionGather => '采集资源';

  @override
  String get regionEmpty => '此区域没有专属的怪物、地下城、资源或任务。';

  @override
  String get regionRaids => '团队副本';

  @override
  String get tagUse => '使用';

  @override
  String get tagAllClasses => '所有职业';

  @override
  String get tagDrops => '掉落来源';

  @override
  String get tagRaid => '团队副本';

  @override
  String get tagDungeon => '地下城';

  @override
  String get tagCrafted => '制作';

  @override
  String tagSkill(int n) {
    return '技能 $n';
  }

  @override
  String get tagSold => '出售';

  @override
  String get tagGathered => '采集';

  @override
  String get tagUsedIn => '作为材料用于';

  @override
  String get petEvolution => '进化';

  @override
  String get petStrength => '实力详情';

  @override
  String petAssistRange(int min, int max) {
    return '战斗援助：每次命中伤害 +$min% – +$max%';
  }

  @override
  String petBonusRange(String min, String max) {
    return '加成：$min – $max';
  }

  @override
  String petFightsRange(int min, int max) {
    return '狩猎实力相当于 $min – $max 级怪物';
  }

  @override
  String get petFind => '获取途径';

  @override
  String petAssist(int pct) {
    return '战斗援助：伤害 +$pct%';
  }

  @override
  String get petLevelNote => '1 级至 100 级';

  @override
  String get petIndex => '宠物图鉴';

  @override
  String petIndexCount(int n, int max) {
    return '已发现 $n / $max';
  }

  @override
  String slotCount(int n, int max) {
    return '$n / $max';
  }

  @override
  String get petEggTapHint => '点击蛋即可孵化。';

  @override
  String get petTapHint => '点击宠物查看详情 · 绿点表示出战宠物。';

  @override
  String get petIndexHint => '点击宠物打开其页面。';

  @override
  String get petUnknown => '尚未发现';

  @override
  String get petAbout => '简介';

  @override
  String get petHabitat => '栖息地';

  @override
  String petElement(String name) {
    return '元素：$name';
  }

  @override
  String get petStage1 => '幼年';

  @override
  String get petStage2 => '成长';

  @override
  String get petStage3 => '成年';

  @override
  String get bodyMale => '男';

  @override
  String get bodyFemale => '女';

  @override
  String get bodyChoice => '外观';

  @override
  String get bodyChange => '更改外观';

  @override
  String get portraitNoArmor => '无护甲';

  @override
  String get combatPower => '战斗力';

  @override
  String get combatTapHint => '点按查看说明';

  @override
  String statCap(int cap) {
    return '上限 $cap%';
  }

  @override
  String get cEvasion => '闪避';

  @override
  String get cBlock => '格挡';

  @override
  String get cCritDmg => '暴击伤害';

  @override
  String get cDmgReduction => '减伤';

  @override
  String get cLifesteal => '吸血';

  @override
  String get cExecute => '斩杀';

  @override
  String get cMissingHp => '背水一战';

  @override
  String get cPetAssist => '宠物协助';

  @override
  String get xPower => '世界事件使用的强度：2 × 攻击 + 防御 + 最大生命 / 10。';

  @override
  String get xAttack => '每次攻击在计算敌方防御前的基础伤害。';

  @override
  String get xDefense => '你受到的每次攻击会减去你防御的一半。';

  @override
  String xCrit(int cap) {
    return '攻击暴击的几率，最多 $cap%。';
  }

  @override
  String get xCritDmg => '暴击造成普通攻击的该倍数伤害。';

  @override
  String get xHit => '命中几率，减去敌人闪避（15–95%）。';

  @override
  String xEvasion(int cap) {
    return '从敌人命中几率中扣除，最多 $cap%。';
  }

  @override
  String xBlock(int cap) {
    return '将受到的攻击减半的几率，最多 $cap%。';
  }

  @override
  String xDmgReduction(int cap) {
    return '受到的每次攻击按此比例降低，最多 $cap%。';
  }

  @override
  String get xLifesteal => '按造成伤害的该比例回复生命。';

  @override
  String xExecute(int pct) {
    return '对生命值不高于 $pct% 的敌人造成额外伤害。';
  }

  @override
  String get xMissingHp => '生命越低，额外伤害越高，最多为此值。';

  @override
  String get xPetAssist => '出战宠物为你的每次攻击增加该比例伤害。';

  @override
  String get lookFace => '脸型';

  @override
  String get lookHairStyle => '发型';

  @override
  String get lookHairColor => '发色';

  @override
  String get lookEyeColor => '瞳色';

  @override
  String get lookSave => '保存';

  @override
  String get lookHumanOnly => '该种族的脸型和发型选项即将推出。';

  @override
  String get raceAbilities => '种族能力';

  @override
  String get raceAbilitiesHint => '均为被动';

  @override
  String get chooseClass => '选择职业';

  @override
  String raceAttrBonus(int pct, String name) {
    return '$name +$pct%';
  }

  @override
  String raceProfessionXp(int pct, String name) {
    return '$name技能经验 +$pct%';
  }

  @override
  String get language => '语言';

  @override
  String get report => '举报';

  @override
  String reportTitle(String name) {
    return '举报 $name';
  }

  @override
  String get reportWhy => '发生了什么？';

  @override
  String get reportHarassment => '骚扰或辱骂';

  @override
  String get reportCheating => '作弊或外挂';

  @override
  String get reportName => '不当名称';

  @override
  String get reportSpam => '刷屏或广告';

  @override
  String get reportScam => '诈骗';

  @override
  String get reportOther => '其他';

  @override
  String get reportNote => '详情（可选）';

  @override
  String get reportSent => '举报已提交，谢谢。';

  @override
  String get errReportSelf => '不能举报自己的角色。';

  @override
  String get nkReportReviewed => '你的举报已处理';

  @override
  String get guide => '指南';

  @override
  String get guideSearch => '搜索指南';

  @override
  String get guideNoMatch => '没有找到结果。';

  @override
  String get guideShowMe => '带我去';

  @override
  String get next => '下一步';

  @override
  String get back => '上一步';

  @override
  String get tutorialSkip => '跳过教程';

  @override
  String get tutorialStart => '开始游戏';

  @override
  String get tutorialReplay => '重新观看教程';

  @override
  String get tutorialWelcomeTitle => '欢迎来到 MeNa World';

  @override
  String get tutorialWelcomeBody => '即使关闭应用，你的角色也会继续工作。给它安排任务，稍后回来领取收益。';

  @override
  String get tutorialCharacterBody => '这里有你的属性、属性点和技能。分配每级获得的点数，并在 300 级和 700 级选择转职。';

  @override
  String get tutorialGatherBody => '选择资源和时长，之后回来领取。采矿、钓鱼、伐木、采药和狩猎各自独立升级。';

  @override
  String get tutorialBagBody => '点击物品可以装备、使用、出售或上架。对比表会显示它是否优于当前装备。';

  @override
  String get tutorialAdventureBody => '狩猎怪物、攻略地下城、参加团本和活动，并在地图上旅行。怪物和地下城取决于你所在的区域。';

  @override
  String get tutorialWorkshopBody => '用职业技能制作物品，在商店购买，并在市场上与其他玩家交易。';

  @override
  String get tutorialTopBarTitle => '顶部栏';

  @override
  String get tutorialTopBarBody => '任务、通知和社交（聊天、好友、邮件、公会、交易）在右上角。角色一次只能做一件事，请先完成或领取当前任务再开始新的。';

  @override
  String get tutorialDoneTitle => '准备就绪';

  @override
  String get tutorialDoneBody => '忘记某个功能在哪里？打开角色页底部的指南。';

  @override
  String get guideCharacterTitle => '属性、属性点和技能';

  @override
  String get guideCharacterBody => '角色页。在进度卡中分配属性点。技能位于可折叠区域；作用相同的技能只保留最强的一个。';

  @override
  String get guidePromotionTitle => '转职';

  @override
  String get guidePromotionBody => '角色页的进度卡。到 300 和 700 级时，需要从两条路线中选择一条才能继续获得经验。期间获得的经验会在选择后发放。';

  @override
  String get guideGatherTitle => '采集';

  @override
  String get guideGatherBody => '采集页。选择职业、可到达的资源以及 1、3 或 6 小时。资源属于不同区域，想要更好的就去旅行。';

  @override
  String get guideBagTitle => '背包和装备';

  @override
  String get guideBagBody => '背包页。点击物品查看详情和对比。长按并拖动可多选，一起出售或摧毁。背包满时无法开始新任务。';

  @override
  String get guideBankTitle => '账号银行';

  @override
  String get guideBankBody => '背包页的银行按钮。你的所有角色共用，可在角色间转移物品。绑定和已装备的物品不能存入。';

  @override
  String get guideHuntTitle => '狩猎';

  @override
  String get guideHuntBody => '冒险页的狩猎。选择可到达的怪物、时长以及携带的药水和食物。角色重伤时会自动撤退。';

  @override
  String get guideDungeonTitle => '地下城和团本';

  @override
  String get guideDungeonBody => '冒险页的地下城。地下城为单人。团本需要 2 到 5 人队伍，可用“寻找队伍”匹配。团本首领掉落套装部件和宠物蛋。';

  @override
  String get guideTravelTitle => '旅行和地图';

  @override
  String get guideTravelBody => '冒险页的旅行。点击区域查看怪物、地下城、资源和任务，然后前往。只有目的地有等级要求。';

  @override
  String get guideEventsTitle => '活动';

  @override
  String get guideEventsBody => '冒险页的活动。报名世界首领和入侵，并领取每日赛季代币。点击卡片查看规则和奖励。';

  @override
  String get guideCraftTitle => '职业和制作';

  @override
  String get guideCraftBody => '工坊页的制作。打开职业查看配方，颜色表示配方还能提供多少经验。附魔卷轴在背包中使用。';

  @override
  String get guideVendorTitle => '商店和回购';

  @override
  String get guideVendorBody => '工坊页的商店。在背包中出售。售出的物品 2 小时内可以回购。';

  @override
  String get guideMarketTitle => '市场和拍卖';

  @override
  String get guideMarketBody => '工坊页的市场。可一口价购买或参与拍卖。在背包中上架物品：一口价，或设置起拍价和可选一口价的拍卖。最后几分钟出价会延长拍卖。';

  @override
  String get guideQuestsTitle => '任务';

  @override
  String get guideQuestsBody => '右上角的剪贴板图标。包含主线、支线和每日任务；徽标显示可提交的数量。';

  @override
  String get guideNotificationsTitle => '通知和收益';

  @override
  String get guideNotificationsBody => '右上角的铃铛图标。收益页保留 30 天内领取的物品。在设置中选择接收哪些通知。';

  @override
  String get guideSocialTitle => '聊天、好友、邮件和交易';

  @override
  String get guideSocialBody => '右上角的对话气泡。交易需双方确认并收税。屏蔽会将玩家从聊天、好友和邮件中移除。';

  @override
  String get guidePartyTitle => '队伍';

  @override
  String get guidePartyBody => '角色页的队伍区域。按名称邀请。做同类工作的成员可获得额外经验。';

  @override
  String get guideGuildTitle => '公会、金库和公会银行';

  @override
  String get guideGuildBody => '右上角对话气泡的公会页。成员可存入金币和物品；职位决定谁可以取用以及每日上限。所有操作都有记录。';

  @override
  String get guideRankingTitle => '排行榜';

  @override
  String get guideRankingBody => '角色页的排行榜按钮。积分来自等级、成就、专精、职业以及地下城和团本通关。可按职业或榜单筛选。';

  @override
  String get guidePetsTitle => '宠物';

  @override
  String get guidePetsBody => '角色页的宠物按钮。蛋来自高等级团本首领。孵化后通过狩猎升级，并激活一只以获得加成和战斗支援。';

  @override
  String get guideHelpersTitle => '助手';

  @override
  String get guideHelpersBody => '角色页的助手按钮。为某一职业雇用助手。无论角色在做什么，它都会独立工作。';

  @override
  String get guideReportTitle => '举报玩家';

  @override
  String get guideReportBody => '长按聊天内容，或使用好友菜单。选择原因后，团队会审核并在处理后通知你。';

  @override
  String get guideAccountTitle => '账号、语言和客服';

  @override
  String get guideAccountBody => '角色列表顶部的账号图标。可更改语言或密码、联系客服、下载数据或删除账号。';

  @override
  String get ranking => '排行榜';

  @override
  String get rankingScore => '积分';

  @override
  String get rankingLevel => '等级';

  @override
  String get rankingAchievements => '成就';

  @override
  String get rankingMastery => '专精';

  @override
  String get rankingProfessions => '职业技能';

  @override
  String get rankingClears => '通关';

  @override
  String get rankingDungeons => '地下城';

  @override
  String get rankingRaids => '团本';

  @override
  String get rankingAllClasses => '所有职业';

  @override
  String get rankingHowScore => '积分如何计算';

  @override
  String get rankingScoreNote => '每个地下城和团本只计一次。每隔几分钟更新。';

  @override
  String get rankingEmpty => '榜单正在统计中。';

  @override
  String get guildVault => '金库和银行';

  @override
  String get guildVaultTreasury => '公会金库';

  @override
  String get guildVaultBank => '公会银行';

  @override
  String get guildVaultDepositGold => '存入金币';

  @override
  String get guildVaultWithdrawGold => '取出金币';

  @override
  String get guildVaultGoldAmount => '金币数量';

  @override
  String get guildVaultNoLimit => '你可以无限额取出。';

  @override
  String guildVaultGoldToday(String used, String limit) {
    return '今日已取：$used / $limit';
  }

  @override
  String guildVaultItemsToday(int used, int limit) {
    return '今日取出物品：$used / $limit';
  }

  @override
  String get guildVaultNoItems => '你的职位暂时不能使用公会银行。';

  @override
  String get guildVaultLog => '记录';

  @override
  String get guildVaultLogEmpty => '暂无记录。';

  @override
  String guildVaultLogIn(String name, String what) {
    return '$name 存入了 $what';
  }

  @override
  String guildVaultLogOut(String name, String what) {
    return '$name 取出了 $what';
  }

  @override
  String get errGuildVaultNotEmpty => '解散前请先清空金库和银行。';

  @override
  String get errGuildVaultClosing => '公会正在解散，无法存入。';

  @override
  String get errGuildVaultLimit => '超出了你的每日上限。';

  @override
  String get errItemBound => '绑定物品不能存入银行。';

  @override
  String get errItemEquipped => '请先卸下该物品。';

  @override
  String get bid => '出价';

  @override
  String get bidPlaced => '出价成功。';

  @override
  String get marketFixed => '一口价';

  @override
  String get marketAuction => '拍卖';

  @override
  String auctionNoBids(String price) {
    return '暂无出价，起拍价 $price';
  }

  @override
  String auctionTopBid(String price, int count) {
    return '最高出价 $price（$count 次）';
  }

  @override
  String get auctionYouLead => '你领先';

  @override
  String get auctionStartBid => '整批起拍价';

  @override
  String get auctionBuyoutOptional => '一口价（可选）';

  @override
  String get auctionRules => '时间结束时最高出价者获胜。最后 5 分钟内出价会延长拍卖。有人出价后不能取消。';

  @override
  String get auctionYouReceiveAtLeast => '至少可得';

  @override
  String auctionMinBid(String price) {
    return '最低出价：$price';
  }

  @override
  String get auctionBidCopper => '你的出价（铜币）';

  @override
  String get auctionHeldNote => '金币会被暂扣，如有人出价更高会立即退还。';

  @override
  String auctionBuyoutFor(String price) {
    return '立即购买：$price';
  }

  @override
  String get errBidTooLow => '出价低于最低要求。';

  @override
  String get errTopBidder => '你已是最高出价者。';

  @override
  String get errAuctionBusy => '另一笔出价正在处理中，请稍后再试。';

  @override
  String get errAuctionHasBids => '已有出价的拍卖不能取消。';

  @override
  String get nkMarketOutbid => '你的出价被超过';

  @override
  String get nkMarketWon => '你赢得了拍卖';

  @override
  String get evMarketOutbid => '你的出价被超过，金币已退还。';
}
