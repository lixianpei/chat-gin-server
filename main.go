package main

import (
	"GoChatServer/helper"
	"GoChatServer/router"
	"GoChatServer/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	helper.InitConfig("./config") //初始化配置
	helper.InitLogger()           //初始化日志
	helper.InitSqlLogger()        //初始化日志
	helper.InitChatDatabase()     //初始化DB
	helper.InitWeiXin()           //初始化微信实例

	//_ = gin.Default()                             //创建gin实例
	engine := gin.New()                           //创建gin实例
	router.InitRoute(engine)                      //初始化API路由
	ws.InitWebsocket(engine)                      //初始化Ws
	engine.Static("/static/uploads", "./uploads") //开启静态文件服务

	//启动服务 TODO : 优雅启动
	err := engine.Run(helper.Configs.Server.Address)
	if err != nil {
		helper.Logger.Error("Main服务启动异常：", err.Error())
	}
	helper.Logger.Error("Main服务已停止....")
}
