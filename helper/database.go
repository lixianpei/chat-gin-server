package helper

import (
	"GoChatServer/dal/query/chat_query"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var Db *chat_query.Query

func InitChatDatabase() {
	dbConfig := Configs.Db
	//连接DB的结构体信息
	DSNConfig := &mysql2.Config{
		User:                 dbConfig.User,
		Passwd:               dbConfig.Password,
		Net:                  "tcp",
		Addr:                 dbConfig.Host + ":" + dbConfig.Port,
		DBName:               dbConfig.Database,
		AllowNativePasswords: true, //设置allowNativePasswords=true，以启用 MySQL 数据库的原生密码认证方法
		Loc:                  time.Now().Local().Location(),
		ConnectionAttributes: "charset=utf8mb4",
		ParseTime:            true, //自动解析time类型字段为字符串，若设置为false会报错：sql: Scan error on column index 9, name "created_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSNConfig: DSNConfig,
	}), &gorm.Config{
		Logger: SqlGormLogger, //定义一个日志接收器
	})
	if err != nil {
		Logger.Error("数据库初始化失败：", err.Error())
		return
	}
	Db = chat_query.Use(db)
	Logger.Info("ChatDb数据库连接初始化成功")
}
