package helper

import (
	"GoChatServer/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ResponseCodeOk    = 200 //返回成功
	ResponseCodeError = 400 //错误
	ResponseCodeLogin = 2   //普通错误
)

type responseData struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	TraceId  string      `json:"trace_id"`
	TraceSql []string    `json:"trace_sql"`
}

func response(c *gin.Context, result *responseData) {
	//DEBUG调试
	result.TraceId = c.GetString(consts.TraceId)
	result.TraceSql = c.GetStringSlice(consts.TraceSql)
	c.JSON(http.StatusOK, result)
	c.Abort()
}

// ResponseOk 返回成功
func ResponseOk(c *gin.Context) {
	response(c, &responseData{
		Code:    ResponseCodeOk,
		Message: "ok",
	})
}

// ResponseOkWithData 返回成功-携带数据
func ResponseOkWithData(c *gin.Context, data interface{}) {
	response(c, &responseData{
		Code:    ResponseCodeOk,
		Message: "ok",
		Data:    data,
	})
}

// ResponseOkWithMessage 返回成功-携带成功消息
func ResponseOkWithMessage(c *gin.Context, message string) {
	response(c, &responseData{
		Code:    ResponseCodeOk,
		Message: message,
	})
}

// ResponseOkWithMessageData 返回成功-携带成功消息和数据
func ResponseOkWithMessageData(c *gin.Context, data interface{}, message string) {
	response(c, &responseData{
		Code:    ResponseCodeOk,
		Message: message,
		Data:    data,
	})
}

// ResponseError 返回错误-携带错误消息
func ResponseError(c *gin.Context, message string) {
	response(c, &responseData{
		Code:    ResponseCodeError,
		Message: message,
	})
}

// ResponseErrorWithData 返回错误-携带错误消息和数据
func ResponseErrorWithData(c *gin.Context, message string, data interface{}) {
	response(c, &responseData{
		Code:    ResponseCodeError,
		Message: message,
		Data:    data,
	})
}
