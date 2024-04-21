package middleware

import (
	"GoChatServer/consts"
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

var LoginAuthUriWhiteList = map[string]bool{
	"/api/im/login":      true,
	"/api/im/phoneLogin": true,
	"/api/im/upload":     true,
}

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前
		uri := c.Request.URL.Path
		if _, ok := LoginAuthUriWhiteList[uri]; ok {
			//白名单跳过验证
			fmt.Println("白名单跳过验证", uri)
			return
		}

		//获取token
		token := c.GetHeader("Authorization")
		claims, err := helper.JwtParseChecking(token) // claims
		if err != nil {
			fmt.Println("鉴权失败：", err.Error())
			helper.ResponseError(c, err.Error())
			c.Abort()
			return
		}

		//鉴权通过，记录当前用户信息
		c.Set(consts.UserId, claims.UserId)
		c.Set(consts.UserNickname, claims.Nickname)
		c.Set(consts.UserPhone, claims.Phone)

		//fmt.Println("c.Get.GetInt64", c.GetInt64(consts.UserId))
		//fmt.Println("c.Get.UserNickname", c.GetString(consts.UserNickname))
		//fmt.Println("c.Get.UserPhone", c.GetString(consts.UserPhone))

		c.Next()
	}
}
