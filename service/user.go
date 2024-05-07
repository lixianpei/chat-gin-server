package service

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/structs"
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

var User = &user{}

type user struct{}

func (u *user) GetLoginUser(c *gin.Context) (*chat_model.User, error) {
	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := chat_model.User{}
	err := qUser.WithContext(c).Where(qUser.ID.Eq(c.GetInt64(consts.UserId))).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	return &mUser, nil
}

func (u *user) GetUserById(userId int64) (*chat_model.User, error) {
	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := chat_model.User{}
	err := qUser.Where(qUser.ID.Eq(userId)).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	return &mUser, nil
}

func (u *user) GetMessageUserById(userId int64) (*structs.UserItem, error) {
	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := structs.UserItem{}
	err := qUser.Select(qUser.ID, qUser.Nickname, qUser.UserName, qUser.Phone, qUser.Avatar).Where(qUser.ID.Eq(userId)).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	mUser.AvatarUrl = helper.GenerateStaticUrl(mUser.Avatar) //生成消息的发送人的头像
	return &mUser, nil
}

func (u *user) GetUsersByGroupId(groupId int64) ([]*structs.UserItem, error) {
	qGU := helper.DbQuery.GroupUser
	qU := helper.DbQuery.User
	users := make([]*structs.UserItem, 0)
	err := qGU.Join(qU, qU.ID.EqCol(qGU.UserID)).
		Select(qU.ALL).
		Where(qGU.GroupID.Eq(groupId)).Scan(&users)
	return users, err
}

func (u *user) GetMessageReceiverUsers(groupId int64, receiver int64) ([]*structs.UserItem, error) {
	users := make([]*structs.UserItem, 0)
	qUser := helper.DbQuery.User

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
		fmt.Println("群消息")
		return u.GetUsersByGroupId(groupId)
	default:
		//广播消息
		fmt.Println("广播消息")
		err := qUser.Scan(&users)
		return users, err
	}
}

// UserList 用户信息
type UserList struct {
	ID            int64               `json:"id"`              // 自增
	Phone         string              `json:"phone"`           // 用户手机号
	UserName      string              `json:"user_name"`       // 用户名称
	Nickname      string              `json:"nickname"`        // 用户昵称
	Gender        int32               `json:"gender"`          // 性别
	Avatar        string              `json:"avatar"`          // 头像
	LastLoginTime string              `json:"last_login_time"` // 最后登录时间
	UnreadCount   int64               `json:"unread_count"`
	AvatarUrl     string              `json:"avatar_url"`
	LastMessage   *chat_model.Message `gorm:"-" json:"last_message"` //新增的自带不在数据表中需要添加 gorm:"-" 避免提示错误
}

// GetFriendContact 获取好友联系人列表
func (u *user) GetFriendContact(c *gin.Context, userId int64) ([]*UserList, error) {
	mUsers := make([]*UserList, 0)
	qu := helper.DbQuery.User
	quc := helper.DbQuery.UserContact
	err := qu.WithContext(c).Join(quc, quc.FriendUserID.EqCol(qu.ID), quc.UserID.Eq(userId)).Scan(&mUsers)
	if len(mUsers) > 0 {
		for k, v := range mUsers {
			mUsers[k].AvatarUrl = helper.GenerateStaticUrl(v.Avatar)
		}
	}
	return mUsers, err
}

// IsFriendContact 判断是否为好友
func (u *user) IsFriendContact(c *gin.Context, userId int64, friendId int64) (int64, error) {
	quc := helper.DbQuery.UserContact
	return quc.WithContext(c).Where(quc.UserID.Eq(userId), quc.FriendUserID.Eq(friendId)).Count()
}
