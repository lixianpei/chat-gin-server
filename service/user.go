package service

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/types"
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

func (u *user) GetMessageUserById(userId int64) (*types.UserItem, error) {
	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := types.UserItem{}
	err := qUser.Select(qUser.ID, qUser.Nickname, qUser.UserName, qUser.Phone, qUser.Avatar).Where(qUser.ID.Eq(userId)).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	mUser.AvatarUrl = helper.GenerateStaticUrl(mUser.Avatar) //生成消息的发送人的头像
	return &mUser, nil
}

func (u *user) GetUsersByRoomId(roomId int64) ([]*types.UserItem, error) {
	qru := helper.DbQuery.RoomUser
	qU := helper.DbQuery.User
	users := make([]*types.UserItem, 0)
	err := qru.Join(qU, qU.ID.EqCol(qru.UserID)).
		Select(qU.ALL).
		Where(qru.RoomID.Eq(roomId)).
		Scan(&users)
	return users, err
}

func (u *user) GetUsersMapByRoomIds(c *gin.Context, roomIds []int64) (usersMap map[int64][]*types.RoomUserItem, err error) {
	if len(roomIds) == 0 {
		return
	}
	qru := helper.DbQuery.RoomUser
	qu := helper.DbQuery.User
	users := make([]*types.RoomUserItem, 0)
	usersMap = make(map[int64][]*types.RoomUserItem)
	err = qru.WithContext(c).
		Join(qu, qu.ID.EqCol(qru.UserID)).
		Select(qru.UserID, qu.Phone, qu.UserName, qu.Nickname, qu.Gender, qu.Avatar, qru.RoomID.As("roomId")).
		Where(qru.RoomID.In(roomIds...)).
		Scan(&users)
	for k, v := range users {
		if v.UserID > 0 {
			users[k].AvatarUrl = helper.GenerateStaticUrl(v.Avatar)
			if _, ok := usersMap[v.RoomId]; !ok {
				usersMap[v.RoomId] = make([]*types.RoomUserItem, 0)
			}
			usersMap[v.RoomId] = append(usersMap[v.RoomId], v)
		}
	}
	return usersMap, err
}

func (u *user) GetMessageReceiverUsers(roomId int64) ([]*types.UserItem, error) {
	users := make([]*types.UserItem, 0)
	qUser := helper.DbQuery.User

	switch {
	case roomId > 0:
		//群聊消息
		return u.GetUsersByRoomId(roomId)
	default:
		//广播消息
		err := qUser.Scan(&users)
		return users, err
	}
}

// UserList 用户信息
type UserList struct {
	ID            int64               `json:"id"`            // 自增
	Phone         string              `json:"phone"`         // 用户手机号
	UserName      string              `json:"userName"`      // 用户名称
	Nickname      string              `json:"nickname"`      // 用户昵称
	Gender        int32               `json:"gender"`        // 性别
	Avatar        string              `json:"avatar"`        // 头像
	LastLoginTime string              `json:"lastLoginTime"` // 最后登录时间
	UnreadCount   int64               `json:"unreadCount"`
	AvatarUrl     string              `json:"avatarUrl"`
	LastMessage   *chat_model.Message `gorm:"-" json:"lastMessage"` //新增的自带不在数据表中需要添加 gorm:"-" 避免提示错误
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
