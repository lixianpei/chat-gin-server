package middleware

import (
	"GoChatServer/helper"
	"github.com/gin-gonic/gin"
	"runtime"
)

// RecoveryHandler 异常恢复
func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 4096)
				runtime.Stack(stack, true)           //通过堆栈或错误信息
				helper.Logger.Println(string(stack)) //记录错误日志
				helper.ResponseError(c, "服务异常")      //友好输出
				return
			}
		}()

		c.Next()
	}
}
