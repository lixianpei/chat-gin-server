package router

import (
	"GoChatServer/api"
	"GoChatServer/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) {

	//全局中间件
	//e.Use(gin.LoggerWithConfig(gin.LoggerConfig{
	//	Formatter: func(params gin.LogFormatterParams) string {
	//		traceId := uuid.NewV4().String()
	//		fmt.Println("LoginAuthHandler...uuid", traceId)
	//
	//		helper.Logger.WithFields(logrus.Fields{
	//			consts.TraceId: traceId,
	//		})
	//		return ""
	//	},
	//}))
	e.Use(middleware.TraceHandler()) //必须第一个，便于记录traceId
	e.Use(middleware.RecoveryHandler())

	//定义路由
	apiRouter := e.Group("/api")
	{
		//鉴权中间件
		apiRouter.Use(middleware.LoginAuthHandler())

		apiRouter.POST("/im/login", api.WxLogin)
		apiRouter.POST("/im/phoneLogin", api.PhoneLogin)
		apiRouter.POST("/im/wxUserSave", api.WxUserInfoSave)
		apiRouter.POST("/im/userInfoSave", api.UserInfoSave)
		apiRouter.POST("/im/getOnlineList", api.GetOnlineList)
		apiRouter.POST("/im/upload", api.UploadFile)
		apiRouter.POST("/im/searchUser", api.SearchUser)
		apiRouter.POST("/im/addFriend", api.ApplyFriend)
	}
}
