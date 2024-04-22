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

func LoginAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.URL.Path
		_, isWhiteList := LoginAuthUriWhiteList[uri]

		//获取token
		token := c.GetHeader("Authorization")
		claims, err := helper.JwtParseChecking(token) // claims
		if err != nil && !isWhiteList {
			fmt.Println("鉴权失败：", err.Error())
			helper.ResponseError(c, err.Error())
			c.Abort()
			return
		}

		//鉴权通过且存在值，记录当前用户信息
		if claims != nil && claims.UserId > 0 {
			c.Set(consts.UserId, claims.UserId)
			c.Set(consts.UserNickname, claims.Nickname)
			c.Set(consts.UserPhone, claims.Phone)
		}

		c.Next()
	}
}
