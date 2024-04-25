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

func (u *user) GetMessageUserById(userId int64) (*chat_model.User, error) {
	//获取当前用户信息
	qUser := helper.Db.User
	mUser := chat_model.User{}
	err := qUser.Select(qUser.ID, qUser.Nickname, qUser.UserName, qUser.Phone, qUser.Avatar).Where(qUser.ID.Eq(userId)).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	mUser.Avatar = helper.GenerateStaticUrl(mUser.Avatar) //生成消息的发送人的头像
	return &mUser, nil
}

func (u *user) GetMessageReceiverUsers(groupId int64, receiver int64) ([]*chat_model.User, error) {
	users := make([]*chat_model.User, 0)
	qUser := helper.Db.User

	switch {
	case receiver > 0:
		//私聊消息
		userInfo, err := u.GetMessageUserById(receiver)
		if userInfo != nil {
			users = append(users, userInfo)
		}
		return users, err
	case groupId > 0:
		//群消息：暂时当做广播消息发送出去 TODO 查询全部用户
		err := qUser.Scan(&users)
		return users, err
	default:
		//广播消息
		err := qUser.Scan(&users)
		return users, err
	}
}

// GetFriendContact 获取好友联系人列表
func (u *user) GetFriendContact(userId int64) ([]*chat_model.User, error) {
	mUsers := make([]*chat_model.User, 0)
	qUser := helper.Db.User
	qUserContact := helper.Db.UserContact
	err := qUser.Join(qUserContact, qUserContact.UserID.EqCol(qUser.ID)).
		Select(qUser.ID, qUser.UserName, qUser.Nickname, qUser.Phone, qUser.Avatar).
		Where(qUserContact.Status.Eq(consts.UserFriendStatusIsFriend)).
		Where(qUser.Where(qUserContact.UserID.Eq(userId)).Or(qUserContact.FriendUserID.Eq(userId))).
		Scan(&mUsers)
	for k, mUser := range mUsers {
		mUsers[k].Avatar = helper.GenerateStaticUrl(mUser.Avatar)
	}
	return mUsers, err
}
