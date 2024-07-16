package consts

// 日期格式化常量
const (
	MonthFormat      = "2006-01"
	DateFormat       = "2006-01-02"
	DateNumberFormat = "20060102"
	TimeFormat       = "2006-01-02 15:04:05"
	MinuteFormat     = "2006-01-02 15:04"
	TimeFormNoSplit  = "20060102150405"
	TimeZoneFormat   = "2006-01-02T15:04:05+08:00"
)

// // vip --已弃用
const (
	VipConfig = "hall:vipConfig" //vip配置
)

// vip --最新的vip
const (
	RedisVipConfig = "hall:vip:config" //vip配置
	RedisVipUser   = "hall:vip:user:"  //vip功能相关用户信息
)

// 二选一礼包
const (
	UserEventGift   = "hall:eventGift:user:"
	EventGiftConfig = "hall:eventGiftConfig"
)

// 救济金&&破产礼包
const (
	Benefit               = "hall:benefit:config"            //todo 这个不要修改，如果修改通知管理后台
	BenefitUserGeted      = "hall:benefit:user_geted:"       //用户领取救济金信息
	BenefitGiftPackConfig = "hall:benefit:gift_config"       //破产礼包配置
	UserBankruptcyTimes   = "hall:benefit:user_gift_times:"  //用户出发破产礼包，但未充值的次数
	BenefitGiftUserBuyNum = "hall:benefit:user_gift_buynum:" //用户出当日购买救济礼包数量
)

// onlyOne 三选一
const (
	OnlyOneConfig        = "hall:onlyOne:config"
	UserOnlyOneGetedInfo = "hall:onlyOne:getedinfo:user"
)

// 充值赠送礼包
const (
	RechargeGiftConfig        = "hall:rechargeGift:config"
	RechargeGiftUserNextPopup = "hall:rechargeGift:user_info" //用户信息
	RechargeGiftUserBuyNum    = "hall:rechargeGift:buy_num:"  //用户信息
)

// 充值礼包
const (
	RechargePackConfigId                   = "hall:rechargePackConfig:id:"
	RechargePackConfigGameRoom             = "hall:rechargePackConfig:game:room:"
	UserRechargePackConfigGameRoomDayTimes = "hall:rechargePackConfig:game:room:day:userTimes" //游戏房间内用户购买礼包次数
)

// 签到
const (
	SignUser           = "hall:sign:user:"
	UserTodayHasSigned = "hall:hasSigned:day:user:" //用户今日是否已签到
	SignConfigList     = "hall:sign:config:list"    //用户签到配置列表
	UserSignStatus     = "hall:sign:status:user:"   //用户签到状态
	SignPrizeDay       = "hall:sign:prize:day:"     //第x天签到奖励
)

// 邮件
const (
	CommonEmailList = "hall:email"              //公共邮件
	EmailPrize      = "hall:email:prize"        //邮件附件奖励
	EmailStatusUser = "hall:email:status:user:" //用户公共邮件状态
	Email           = "hall:email:"
)

// 奖品
const (
	PrizeList = "hall:prize:list" //奖品列表
)

// luckspin
const (
	LuckSpinConfig     = "hall:luckspin:config"     //luckspin配置信息
	LuckSpinGiftConfig = "hall:luckspin:giftconfig" //luckspin充值礼包配置信息
	LuckSpinUser       = "hall:luckspin:user:"      //luckspin 用户信息
)

// 人机随机信息
const (
	RobotRandInfoKey = "hall:RobotInfo:rand_info" //人机随机信息
)

// 弹窗配置信息
const (
	PopupWindowInfo = "hall:popup_window_info"
)

// 大厅入口以及弹窗配置信息
const (
	EntrancePopupConfig  = "hall:entrance_popup:config" //配置信息
	EntrancePopupUserPop = "hall:entrance_popup:user:"  //用户弹出信息
)

// 任务活动信息
const (
	RedisTaskActivityInfo      = "hall:TaskActivity:acticity_info"   //活动主体信息
	RedisTaskActivityConfig    = "hall:TaskActivity:acticity_config" //活动任务配置信息
	RedisTaskActivityUserInfo  = "hall:TaskActivity:activity_user:"  //用户参加活动的信息
	RedisTaskUserInfo          = "hall:TaskActivity:user:"           //用户任务信息
	RedisTaskUserAward         = "hall:TaskActivity:user_award:"     //用户奖励信息
	RedisActivityNewPlayerPool = "hall:act_new_player:pool"          //新手嘉年华-转盘信息
	RedisActNPUserInfo         = "hall:act_new_player:user:"         //用户新手嘉年华用户信息
)

const (
	PayConfig = "hall:PayConfig"
)

// todo redis key
const (
	UserChannel        = "UserChannel:"
	UserChannelList    = "user:channelList"
	MapUserChannelList = "MapUserChannelList"
	AdminUserChannel   = "AdminUserChannel:" //管理后台用户角色的渠道数据
	PassageList        = "PassageList"
	MapPassageList     = "MapPassageList"
	PayCfgList         = "PayCfgList"
	MapPayCfgList      = "MapPayCfgList"
)

// 红点
const (
	RedisRedDotUser = "hall:reddot:user:" //人机随机信息
)

// 总充值
const (
	RedisRechargeOrderNo = "recharge:order:" //
)

// 分享裂变
const (
	FXLB_dict          = "hall:fxlb:dict"           //分享裂变通用信息
	Fxlb_level_config  = "hall:fxlb:level_config"   //分享裂变等级配置
	Fxlb_bind_info     = "hall:fxlb:bind_info:"     //分享裂变用户绑定信息
	Fxlb_withdraw_info = "hall:fxlb:withdraw_info:" //分享裂变用户提现信息
)
const (
	UserIconMax = 20 //用户随机头像最大范围
)
const (
	UserTag_New int = 1 << iota //新手
)
