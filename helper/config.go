package helper

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewConfig(configPath string) (c *ConfigData) {
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
	// 服务相关配置
	Server configServer

	// 数据库相关配置
	Db configDb
}

// 服务相关配置
type configServer struct {
	Address string
	Env     string
}

// 数据库相关配置
type configDb struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}
