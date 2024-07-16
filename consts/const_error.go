package consts

import (
	lib "za.game/lib/response"
)

var (
	DataAnalysisError = lib.Error(2013, "数据解析失败")
	DataSaveError     = lib.Error(2014, "数据保存失败")
	DataFindError     = lib.Error(2015, "数据查询失败")
	CacheGetError     = lib.Error(2016, "缓存数据获取失败")
	CacheSaveError    = lib.Error(2017, "缓存数据保存失败")
	DataGetError      = lib.Error(2018, "数据获取失败")

	UserNameNotFind = lib.Error(3001, "支付渠道不存在")
	UserPwdFail     = lib.Error(3002, "订单号错误")
	PassageNotExist = lib.Error(3003, "通道不存在")

	UserUidErr        = lib.Error(3003, "用户ID获取失败")
	PayConfigErr      = lib.Error(3004, "获取支付配置列表失败")
	PayOrderErr       = lib.Error(3005, "下单失败")
	TokenExpireErr    = lib.Error(3006, "token过期")
	PayAmountErr      = lib.Error(3007, "支付金额必去是正整数")
	PayOrderMoreErr   = lib.Error(3008, "下单太频繁，稍后再试")
	PayOrderSelectErr = lib.Error(3009, "订单号异常")
	PayGiftErr        = lib.Error(3010, "获取充值礼包错误")
	PayGiftEmpty      = lib.Error(3011, "请先配置充值礼包")
	PayConfigDataErr  = lib.Error(3012, "支付配置ID异常")
	PayBalanceErr     = lib.Error(3013, "余额查询失败")
	PayOrderStatusErr = lib.Error(3014, "查询支付订单状态失败")
	PaymentAmountErr  = lib.Error(3015, "代付金额必须是10000至5000000之间的正整数")
	PaymentErr        = lib.Error(3016, "代付订单创建失败")
	PaymentStatusErr  = lib.Error(3017, "查询代付订单状态失败")
	PayBankErr        = lib.Error(3018, "用户银行卡数据异常")
	PayBankUpErr      = lib.Error(3019, "用户银行卡数据修改失败")
	PayBalanceEnough  = lib.Error(3020, "可赠送余额不足")
	AccountEnoughErr  = lib.Error(3021, "账户余额不够")
	ClickMoreErr      = lib.Error(3022, "请求太频繁，稍后再试")
	PayBuyMax         = lib.Error(3023, "达到购买上限")
	ActivityClosed    = lib.Error(3024, "活动已关闭")
	NoBroken          = lib.Error(3025, "未达到破产值")
	BonusError        = lib.Error(3026, "bonus draw fail")

	UnpurchasedRechargeGiftPackage = lib.Error(7014, "未购买充值赠送礼包")
	NoCanWithdraw                  = lib.Error(7015, "不满足提现要求")
	WithDrawTop                    = lib.Error(7016, "达到提现上限")
	NoRechargeCanNotWithDraw       = lib.Error(7017, "未充值过不能提现")
)
