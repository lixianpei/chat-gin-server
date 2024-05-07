package helper

import (
	"GoChatServer/dal/query/chat_query"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DbQuery *chat_query.Query
var Db *gorm.DB

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
	sqlDb, err := db.DB()
	if err != nil {
		Logger.Error("获取SqlDb失败：", err.Error())
		return
	}

	//设置相关的数据库相关配置 https://blog.csdn.net/qq_39384184/article/details/103954821
	sqlDb.SetMaxOpenConns(dbConfig.MaxOpenConns)                      //数据库连接池的连接数，需要小于数据库的最高配置，推荐连接数 = ((核心数 * 2) + 有效磁盘数)
	sqlDb.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime)) //设置连接的最大空闲时间:通常，可以将最大空闲时间设置为数据库连接超时时间的一半左右，或者根据应用程序的预期负载来调整。
	sqlDb.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime)) //连接可重用的最大时间时长，超时后连接不可再重用:通常，可以将最大生命周期设置为数据库配置的超时时间或稍短一些，以确保连接不会长时间处于不稳定的状态。例如，如果数据库的连接超时设置为5分钟，可以将连接的最大生命周期设置为4分钟左右。
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdleConns)                      //空闲连接保留在连接池的最大连接数，当设置为 0 时，必须为每个连接从头开始创建一个新连接,通常，可以将最大空闲连接数设置为最大打开连接数的一半左右，或者根据应用程序的预期负载来调整。若此数量太小，并发量上来后，同一时间创建新的请求就会出现问题：connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.

	if err != nil {
		Logger.Error("数据库初始化失败：", err.Error())
		return
	}
	Db = db                      //支持原生sql查询
	DbQuery = chat_query.Use(db) //更安全的db查询

	//db.Raw()
	Logger.Info("ChatDb数据库连接初始化成功")
}
