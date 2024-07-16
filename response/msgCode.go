package response

var (
	OK                    = Error(200, "Success")
	CommonParam           = Error(1001, "参数或者参数格式有误")
	CommonDataException   = Error(1002, "数据异常")
	CommonOptionFail      = Error(1003, "操作失败")
	CommonIllegalRequest  = Error(1004, "非法请求")
	CommonDataWriterFail  = Error(1005, "数据写入失败")
	CommonParamBind       = Error(1006, "参数绑定失败")
	CommonStatusException = Error(1007, "状态异常")
	SignException         = Error(1008, "签名错误")

	TokenNoFound     = Error(2004, "缺少token")
	TokenNoError     = Error(2005, "token错误")
	TokenFormatError = Error(2006, "请求头中 Authorization 格式有误")
	UserInfoError    = Error(2007, "用户信息缺失")
)
