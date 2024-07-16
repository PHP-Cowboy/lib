package consts

// 互动类型
const (
	TaskActivityTypeUserCentric     int8 = 1 //活动时间以用户为主
	TaskActivityTypeActivityCentric int8 = 2 //活动时间以活动为主
)

// 活动id
const (
	AvtId_NewPlayer int = 1 //新手嘉年华
)

// 任务id
const (
	TaskId_SlotFruit      int = 1  //完成 %d 局slot水果
	TaskId_SlotCat        int = 2  //完成 %d 局slotmao
	TaskId_SlotIndia      int = 3  //完成 %d 局slot印第安
	TaskId_SlotRand       int = 4  //完成 %d 局任意slot
	TaskId_Tp             int = 5  //完成 %d 局TP
	TaskId_LoginAccum     int = 6  //累计登录 %d 天
	TaskId_RechargeAccum  int = 7  //累积充值大于等于 %d 卢比
	TaskId_RechargeSingle int = 8  //单笔充值大于等于 %d 卢比
	TaskId_BetTotal       int = 9  //总下注达到 %d 卢比
	TaskId_BonusToCash    int = 10 //从存钱罐中收集 %d 次cash
	TaskId_VipBankToCash  int = 11 //从vip存钱罐中收集 %d 次cash
	TaskId_PlayByBetSigle int = 12 //在单局下注大于等于100的情况下，进行游戏完成 %d 局
	TaskId_Sign7          int = 13 //完成 %d 次完整7天签到
	TaskId_King           int = 14 //完成 %d 局国王皇后
	TaskId_Dragon         int = 15 //完成 %d 局龙虎斗
	TaskId_SlotZeus       int = 16 //完成 %d 局slot宙斯
)
