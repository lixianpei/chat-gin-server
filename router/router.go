package router

import (
	"GoChatServer/api"
	"GoChatServer/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) {
	//定义路由
	apiRouter := e.Group("/api")
	{
		//鉴权中间件
		apiRouter.Use(middleware.LoginAuth())

		apiRouter.POST("/im/login", api.WxLogin)
		apiRouter.POST("/im/phoneLogin", api.PhoneLogin)
		apiRouter.POST("/im/wxUserSave", api.WxUserInfoSave)
		apiRouter.POST("/im/userInfoSave", api.UserInfoSave)
		apiRouter.POST("/im/getOnlineList", api.GetOnlineList)
		apiRouter.POST("/im/upload", api.UploadFile)
	}
}
