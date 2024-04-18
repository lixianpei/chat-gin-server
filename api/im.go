package api

import (
	"GoChatServer/helper"
	"GoChatServer/ws"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Phone    string `form:"phone" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
}

func Login(c *gin.Context) {
	var loginForm LoginForm
	err := c.ShouldBind(&loginForm)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//生成token
	token, err := helper.NewJwtToken(loginForm.Phone, loginForm.Nickname)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	helper.ResponseOkWithMessageData(c, gin.H{
		"token": token,
	}, "ok")
}

// GetOnlineList 获取在线的所有客户端
func GetOnlineList(c *gin.Context) {
	clients := ws.IM.OnlineClients()
	helper.ResponseOkWithData(c, clients)
}
