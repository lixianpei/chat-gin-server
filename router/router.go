package router

import (
	"GoChatServer/api"
	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) {
	//定义路由
	apiRouter := e.Group("/api")
	{
		//鉴权中间件
		//apiRouter.Use(middleware.LoginAuth())

		apiRouter.POST("/im/login", api.WxLogin)
		apiRouter.POST("/im/wx_user_save", api.WxUserInfoSave)
		apiRouter.POST("/im/wx_user_avatar_save", api.WxUserAvatarSave)
		apiRouter.POST("/im/GetOnlineList", api.GetOnlineList)
		apiRouter.POST("/im/upload", api.UploadFile)
	}
}
