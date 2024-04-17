package api

import (
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Phone    string `form:"phone" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
}

func Login(c *gin.Context) {
	var loginForm LoginForm
	err := c.ShouldBindQuery(&loginForm)
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

	//解析数据
	claims, err := helper.JwtParseChecking(token)
	fmt.Println(err, claims)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	data := gin.H{
		"token":  token,
		"claims": claims,
	}

	helper.ResponseOkWithMessageData(c, data, "ok")
}
