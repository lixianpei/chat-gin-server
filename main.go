package main

import (
	"GoChatServer/im"
	"GoChatServer/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//创建gin实例
	engine := gin.Default()

	//初始化API路由
	router.InitRoute(engine)

	//初始化Ws
	im.InitWebsocket(engine)

	//启动服务 TODO : 优雅启动
	err := engine.Run(":8081")
	if err != nil {
		log.Panicln(err.Error())
	}
}
