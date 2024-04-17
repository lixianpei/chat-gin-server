package router

import (
	"GoChatServer/api"
	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) {
	//定义路由
	apiRouter := e.Group("/api")
	{
		apiRouter.GET("/login", api.Login)
	}
}
