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
	fmt.Println(loginForm)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	helper.ResponseOkWithMessageData(c, loginForm, "ok")
}
