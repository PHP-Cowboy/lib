/**用于注册一个公共的接口*/

package account

import "git.dev666.cc/external/dreamgo/xy"

var (
	CB_SendNoticeEvent func(ev string, dic any, userid ...xy.UserId_)
)

/**注册通知函数-在每个项目中单独调用*/
func RegisterSendNoticeEvent(f func(ev string, dic any, userid ...xy.UserId_)) {
	CB_SendNoticeEvent = f
}
