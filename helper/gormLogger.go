package helper

import (
	"GoChatServer/consts"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

var SqlGormLogger = &GormLogger{}

// GormLogger 日志记录器，便于捕获HTTP请求过程中执行的相关sql，并且可以记录相关SQL
type GormLogger struct {
	*logrus.Logger
}

func InitSqlLogger() {
	log := &GormLogger{}
	log.Logger = logrus.New()
	log.SetReportCaller(true)
	SqlGormLogger = log
}

// LogMode 实现 GORM Logger 接口的 LogMode 方法
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info 实现 GORM Logger 接口的 Info 方法
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	Logger.Errorf("GormInfo:" + msg)
}

// Warn 实现 GORM Logger 接口的 Warn 方法
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	Logger.Errorf("GormWarn:" + msg)
}

// Warn 实现 GORM Logger 接口的 Error 方法
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	Logger.Errorf("GormError:" + msg)
}

// Trace 实现 GORM Logger 接口的 Trace 方法
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	//构造SQL日志内容
	sqlLogContent := fmt.Sprintf("GormLogger:time=[%s],rows=[%d],SQL: %s", time.Since(begin).String(), rows, sql)
	//在ctx中获取到Gin的上下文，然后通过把当前sql设置上下文中，返回数据时直接读取
	c, ok := ctx.Value(gin.ContextKey).(*gin.Context)
	if ok {
		sqls := c.GetStringSlice(consts.TraceSql)
		sqls = append(sqls, sqlLogContent)
		c.Set(consts.TraceSql, sqls)
	}
	Logger.Info(sqlLogContent)
}

// RowsAffected 实现 GORM Logger 接口的 RowsAffected 方法
func (l *GormLogger) RowsAffected(ctx context.Context, rows int64) {
	fmt.Println("RowsAffected......")
}
