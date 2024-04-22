package middleware

import (
	"GoChatServer/consts"
	"GoChatServer/helper"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// TraceHandler 链路跟踪中间件，生成TraceId提供给日志记录使用
func TraceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := uuid.NewV4().String()
		c.Set(consts.TraceId, traceId)                                     //将traceId记录到上下文中，提供给其他函数使用
		helper.Logger.AddHook(&CustomLogTraceId{TraceId: traceId, ctx: c}) //将traceId记录到日志组件中，方便打印数据
	}
}

type CustomLogTraceId struct {
	ctx     *gin.Context
	TraceId string
}

func (hook *CustomLogTraceId) Fire(entry *logrus.Entry) error {
	entry.Data[consts.TraceId] = hook.TraceId
	return nil
}
func (hook *CustomLogTraceId) Levels() []logrus.Level {
	return logrus.AllLevels
}
