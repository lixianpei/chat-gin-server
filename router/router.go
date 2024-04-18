package router

import (
	"GoChatServer/api"
	"GoChatServer/helper"
	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) {
	//定义路由
	apiRouter := e.Group("/api")
	{
		apiRouter.GET("/im/ping", func(c *gin.Context) {
			helper.ResponseOk(c)
		})
		apiRouter.POST("/im/login", api.Login)
		apiRouter.POST("/im/GetOnlineList", api.GetOnlineList)
	}
}
