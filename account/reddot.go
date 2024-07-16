package account

import (
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"git.dev666.cc/external/dreamgo/xy"
)

/**红点公共函数*/

/*
*添加红点
@param uid,uidstr : 用户id
@param rdid : 红点id
@param count : 数量
@param bclear : 是否清空
@param bsync : 是否通知客户端
*/
func AddRedDot(uid uint64, uidstr string, rdid int, count int, bclear bool, bsync bool) {
	if bsync {
		if CB_SendNoticeEvent == nil {
			logger.Logger.Errorf("AddRedDot: no register SendNoticeEvent!")
			return
		}

		//发送通知
		dicJson := make(map[string]interface{}, 0)
		dicJson["rdid"] = rdid
		dicJson["count"] = count
		CB_SendNoticeEvent("reddot", dicJson, xy.UserId_(uid))
	}
}
