package service

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

var User = &user{}

type user struct{}

func (u *user) GetLoginUser(c *gin.Context) (*chat_model.User, error) {
	//获取当前用户信息
	qUser := helper.Db.User
	mUser := chat_model.User{}
	err := qUser.WithContext(c).Where(qUser.ID.Eq(c.GetInt64(consts.UserId))).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	return &mUser, nil
}
func (u *user) GetUserById(userId int64) (*chat_model.User, error) {
	//获取当前用户信息
	qUser := helper.Db.User
	mUser := chat_model.User{}
	err := qUser.Where(qUser.ID.Eq(userId)).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	return &mUser, nil
}
