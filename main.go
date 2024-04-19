package main

import (
	"GoChatServer/helper"
	"GoChatServer/router"
	"GoChatServer/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	//Configs 配置文件相关配置
	Configs *helper.ConfigData

	// Logger 日志
	Logger *logrus.Logger
)

func main() {
	//初始化配置
	Configs = helper.NewConfig("./config")

	//初始化日志
	Logger = helper.NewLogger()

	//创建gin实例
	engine := gin.Default()

	//初始化API路由
	router.InitRoute(engine)

	//初始化Ws
	ws.InitWebsocket(engine)

	Logger.Info("Main服务已启动...0000")

	//启动服务 TODO : 优雅启动
	err := engine.Run(Configs.Server.Address)
	if err != nil {
		Logger.Error("Main服务启动异常：", err.Error())
	}
	Logger.Error("Main服务已停止....")
}
