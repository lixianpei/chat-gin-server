package helper

import (
	"fmt"
	"github.com/spf13/viper"
)

var Configs *ConfigData

func InitConfig(path string) {
	Configs = newConfig(path)
}

func newConfig(configPath string) (c *ConfigData) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取失败:", err.Error())
		return nil
	}

	err = v.Unmarshal(&c)
	if err != nil {
		fmt.Println("配置信息绑定结构体失败：", err.Error())
		return nil
	}

	fmt.Println("配置读取成功...")
	fmt.Println(fmt.Sprintf("配置信息：%+v", c))
	return c
}

type ConfigData struct {
	Server    configServer // 服务相关配置
	Db        configDb     // 数据库相关配置
	WeiXin    configWeiXin //微信相关配置
	Redis     configRedis
	Websocket configWebsocket
}

// 服务相关配置
type configServer struct {
	Address               string
	Env                   string
	Host                  string
	UploadFilePath        string
	StaticFileServerPath  string
	DefaultAvatar         []string
	MaxUploadFileSizeMb   int64    //上传文件的最大容量
	AllowUploadExtensions []string //允许上传的文件后缀
}

// 数据库相关配置
type configDb struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	ConnMaxIdleTime int
	ConnMaxLifetime int
	MaxIdleConns    int
}

// 微信相关配置
type configWeiXin struct {
	Appid  string
	Secret string
}

// Redis相关配置
type configRedis struct {
	Address  string
	Password string
	Prefix   string
}

// Websocket相关配置
type configWebsocket struct {
	WriteWait      int64
	PongWait       int64
	MaxMessageSize int64
}
