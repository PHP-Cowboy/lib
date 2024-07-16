package consts

type GiveMoneyType int

const (
	GiveMoneyTypeNormal       GiveMoneyType = 1 //正常提现
	GiveMoneyTypeActNewPlayer GiveMoneyType = 2 //嘉年华提现
	GiveMoneyTypeFxlb         GiveMoneyType = 3 //分享裂变
)
