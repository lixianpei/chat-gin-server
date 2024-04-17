package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ResponseCodeOk    = 200 //返回成功
	ResponseCodeError = 400 //Token解析错误
	ResponseCodeLogin = 2   //普通错误
)

// ResponseOk 返回成功
func ResponseOk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeOk,
		"message": "success",
		"data":    gin.H{},
	})
	c.Abort()
}

// ResponseOkWithData 返回成功-携带数据
func ResponseOkWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeOk,
		"message": "success",
		"data":    data,
	})
	c.Abort()
}

// ResponseOkWithMessage 返回成功-携带成功消息
func ResponseOkWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeOk,
		"message": message,
		"data":    gin.H{},
	})
	c.Abort()
}

// ResponseOkWithMessageData 返回成功-携带成功消息和数据
func ResponseOkWithMessageData(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeOk,
		"message": message,
		"data":    data,
	})
	c.Abort()
}

// ResponseError 返回错误-携带错误消息
func ResponseError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeError,
		"message": message,
		"data":    gin.H{},
	})
	c.Abort()
}

// ResponseErrorWithData 返回错误-携带错误消息和数据
func ResponseErrorWithData(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    ResponseCodeError,
		"message": message,
		"data":    data,
	})
	c.Abort()
}
