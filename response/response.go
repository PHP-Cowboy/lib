package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Success struct {
	Code int
	Msg  string
	Data interface{}
}

type ServiceData struct {
	TotalNum  int         `json:"totalNum"`
	TotalPage int         `json:"totalPage"`
	Data      interface{} `json:"data"`
}

type NoPageServiceData = []map[string]interface{}

type OptionServiceData map[string]interface{}

func (success *Response) Result(request *gin.Context) {
	request.JSON(http.StatusOK, Response{
		Code: success.Code,
		// Msg:  success.Msg, //取消返回信息
		Data: success.Data,
	})
}

func (fail *Response) ResultError(request *gin.Context) {
	request.JSON(http.StatusOK, Response{
		Code: fail.Code,
		// Msg:  fail.Msg,//取消返回信息
	})
}

func Error(code int, msg string) Response {
	return Response{
		code,
		msg,
		nil,
	}
}

func Ok(request *gin.Context) {
	success := Response{http.StatusOK, "Success", nil}
	success.Result(request)
}

func OkWithData(request *gin.Context, data interface{}) {
	success := Response{http.StatusOK, "Success", data}
	success.Result(request)
}

func OkCodeWithData(request *gin.Context, code int, data interface{}) {
	success := Response{code, "Success", data}
	success.Result(request)
}

func OkWithMsg(request *gin.Context, msg string) {
	success := Response{http.StatusOK, msg, nil}
	success.Result(request)
}

func Fail(request *gin.Context, err Response) {
	fail := Response{err.Code, err.Msg, nil}
	fail.ResultError(request)
}

func FailWithMsg(request *gin.Context, err Response, msg string) {
	fail := Response{err.Code, msg, nil}
	fail.ResultError(request)
}
