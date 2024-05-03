package api

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/query/chat_query"
	"GoChatServer/dal/structs"
	"GoChatServer/helper"
	"GoChatServer/service"
	"GoChatServer/ws"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"path"
	"time"
)

type WxLoginForm struct {
	Code string `form:"code" json:"code" binding:"required"`
}

// WxLogin 微信登录
func WxLogin(c *gin.Context) {
	var loginForm WxLoginForm
	err := c.ShouldBind(&loginForm)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//调用微信登录
	wxResult, err := helper.WxApi.Login(loginForm.Code)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//查询用户是否存在
	qUser := helper.DbQuery.User
	mUserInfo := &chat_model.User{}
	err = helper.DbQuery.WithContext(c).User.Where(qUser.WxOpenid.Eq(wxResult.OpenId)).Scan(mUserInfo)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if mUserInfo.ID == 0 {
		//创建新用户
		mUserInfo = &chat_model.User{
			WxOpenid:      wxResult.OpenId,
			WxUnionid:     wxResult.UnionId,
			WxSessionKey:  wxResult.SessionKey,
			LastLoginIP:   c.ClientIP(),
			LastLoginTime: time.Now().Local(),
		}
		err = helper.DbQuery.WithContext(c).User.Create(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	} else {
		//保存信息数据
		_, err = helper.DbQuery.WithContext(c).User.Where(qUser.ID.Eq(mUserInfo.ID)).Updates(&chat_model.User{
			//WxOpenid:     wxResult.OpenId,
			WxUnionid:     wxResult.UnionId,
			WxSessionKey:  wxResult.SessionKey,
			LastLoginIP:   c.ClientIP(),
			LastLoginTime: time.Now().Local(),
		})
	}

	//生成token，此时还未获取到昵称和用户名
	token, err := helper.NewJwtToken(mUserInfo.ID, "", "")
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	helper.ResponseOkWithMessageData(c, gin.H{
		"user_id":   mUserInfo.ID,
		"token":     token,
		"avatarUrl": helper.GenerateStaticUrl(mUserInfo.Avatar),
		"userInfo":  mUserInfo,
	}, "ok")
}

type PhoneLoginForm struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
	Avatar   string `form:"avatar" json:"avatar"`
}

// PhoneLogin 手机号登录
func PhoneLogin(c *gin.Context) {
	var loginForm PhoneLoginForm
	err := c.ShouldBind(&loginForm)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if len(loginForm.Avatar) == 0 {
		loginForm.Avatar = helper.Configs.Server.DefaultAvatar[0]
	}

	qUser := helper.DbQuery.User
	mUserInfo := &chat_model.User{}
	//若用户已提前使用微信登录，则此时已经存在token，可以获取到登录的用户信息
	loginUser, err := service.User.GetLoginUser(c)

	if loginUser != nil {
		mUserInfo = loginUser
	} else {
		//直接使用手机号进行登录：根据手机号查询用户是否存在
		err = helper.DbQuery.WithContext(c).User.
			Select(qUser.ID, qUser.UserName, qUser.Nickname).
			Where(qUser.Phone.Eq(loginForm.Phone)).
			Scan(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	}

	if mUserInfo.ID == 0 {
		//创建新用户
		mUserInfo = &chat_model.User{
			Phone:         loginForm.Phone,
			Nickname:      loginForm.Nickname,
			Avatar:        loginForm.Avatar,
			LastLoginIP:   c.ClientIP(),
			LastLoginTime: time.Now().Local(),
		}
		err = helper.DbQuery.WithContext(c).User.Create(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	} else {
		//保存信息数据
		_, err = helper.DbQuery.WithContext(c).User.Where(qUser.ID.Eq(mUserInfo.ID)).Updates(&chat_model.User{
			Nickname:      loginForm.Nickname,
			Phone:         loginForm.Phone,
			Avatar:        loginForm.Avatar,
			LastLoginIP:   c.ClientIP(),
			LastLoginTime: time.Now().Local(),
		})
	}

	//生成token
	token, err := helper.NewJwtToken(mUserInfo.ID, loginForm.Phone, loginForm.Nickname)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	helper.ResponseOkWithMessageData(c, gin.H{
		"user_id":  mUserInfo.ID,
		"token":    token,
		"phone":    loginForm.Phone,
		"nickname": loginForm.Nickname,
	}, "ok")
}

type WxUserInfoForm struct {
	EncryptedData string `form:"encryptedData" json:"encryptedData"`
	RawData       string `form:"rawData" json:"rawData"`
	Signature     string `form:"signature" json:"signature"`
	Iv            string `form:"iv" json:"iv"`
}
type WxUserInfoData struct {
	Openid    string `json:"openid"`
	Nickname  string `json:"nickname"`
	Gender    int32  `json:"gender"`
	AvatarUrl string `json:"avatarUrl"`
}

// WxUserInfoSave 微信信息存储-由于信息没有可用价值，接口暂不需要
func WxUserInfoSave(c *gin.Context) {
	var form WxUserInfoForm
	err := c.ShouldBind(&form)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := chat_model.User{}
	err = qUser.WithContext(c).Where(qUser.ID.Eq(c.GetInt64(consts.UserId))).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		helper.ResponseError(c, "用户不存在")
		return
	}

	//验证数据是否被篡改
	isOk := helper.WxApi.CheckWxSignature(form.RawData, mUser.WxSessionKey, form.Signature)
	if !isOk {
		helper.ResponseError(c, "数据已被篡改，请稍后重试！")
		return
	}

	//检测和解密微信数据
	var wxUserInfo = WxUserInfoData{}
	wxDecodeDataString := helper.WxApi.DecodeWxData(form.EncryptedData, mUser.WxSessionKey, form.Iv)
	err = json.Unmarshal([]byte(wxDecodeDataString), &wxUserInfo)
	if err != nil || len(wxDecodeDataString) == 0 {
		helper.ResponseError(c, "用戶信息解密失败")
		return
	}

	//保存用户信息 TODO 获取到的新的都是空的
	updateUser := chat_model.User{
		Avatar:   wxUserInfo.AvatarUrl,
		Nickname: wxUserInfo.Nickname,
		Gender:   wxUserInfo.Gender,
	}
	_, err = qUser.WithContext(c).Where(qUser.ID.Eq(mUser.ID)).Updates(updateUser)

	//返回数据
	helper.ResponseOkWithData(c, gin.H{
		"wxUserForm": form,
		"user_id":    c.GetInt64(consts.UserId),
	})
}

type UserAvatarForm struct {
	Avatar   string `form:"avatar" json:"avatar"`
	Nickname string `form:"nickname" json:"nickname"`
	Phone    string `form:"phone" json:"phone"`
}

// UserInfoSave 微信头像存储-头像为临时头像，暂时不需要此接口
func UserInfoSave(c *gin.Context) {
	var form UserAvatarForm
	err := c.ShouldBind(&form)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//获取当前用户信息
	qUser := helper.DbQuery.User
	mUser := chat_model.User{}
	err = qUser.WithContext(c).Where(qUser.ID.Eq(c.GetInt64(consts.UserId))).Scan(&mUser)
	if err != nil || mUser.ID == 0 {
		helper.ResponseError(c, fmt.Sprintf("用户不存在：%d", c.GetInt64(consts.UserId)))
		return
	}

	//保存用户信息
	updateUser := chat_model.User{
		Avatar:   form.Avatar,
		Nickname: form.Nickname,
		Phone:    form.Phone,
	}
	_, err = qUser.WithContext(c).Where(qUser.ID.Eq(mUser.ID)).Updates(updateUser)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	helper.ResponseOk(c)
}

// GetOnlineList 获取在线的所有客户端
func GetOnlineList(c *gin.Context) {
	clients := ws.IM.OnlineClients()
	helper.ResponseOkWithData(c, clients)
}

func UploadFile(c *gin.Context) {
	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	uuider := uuid.NewV4()
	filepath := path.Join("avatars", uuider.String()+path.Ext(file.Filename))
	dst := path.Join(helper.Configs.Server.UploadFilePath, filepath)

	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	helper.ResponseOkWithData(c, gin.H{
		"filepath": filepath,                           //数据库存储的文件路径
		"url":      helper.GenerateStaticUrl(filepath), //访问文件的url
	})
}

type SendMessageForm struct {
	//Type     string `form:"type" binding:"required"`//消息类型
	Content  string `form:"content" json:"content" binding:"required"`   //消息类型
	Receiver int64  `form:"receiver" json:"receiver" binding:"required"` //消息接收的用户ID
}

func SendMessage(c *gin.Context) {
	var form SendMessageForm
	err := c.ShouldBind(&form)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	mUser, err := service.User.GetLoginUser(c)
	if err != nil {
		helper.ResponseError(c, fmt.Sprintf("用户不存在：%d", c.GetInt64(consts.UserId)))
		return
	}

	//消息内容
	mMessage := chat_model.Message{
		Sender:  mUser.ID,
		Content: form.Content,
	}
	//消息关联的用户
	messageUsers := make([]*chat_model.MessageUser, 0)

	//消息入库
	err = helper.DbQuery.Transaction(func(tx *chat_query.Query) error {
		//保存消息
		err = tx.WithContext(c).Message.Create(&mMessage)
		if err != nil {
			return err
		}

		//查询消息关联的用户,TODO 暂时查询全量用户
		users := make([]*chat_model.User, 0)
		err = tx.WithContext(c).User.Scan(&users)
		if err != nil {
			return err
		}

		for _, user := range users {
			messageUsers = append(messageUsers, &chat_model.MessageUser{
				MessageID: mMessage.ID,
				Receiver:  user.ID,
				IsRead:    0,
			})
		}
		if len(messageUsers) > 0 {
			err = tx.WithContext(c).MessageUser.Create(messageUsers...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//消息推入消息中心：
	if len(messageUsers) > 0 {
		for _, messageUser := range messageUsers {
			message := ws.Message{
				MessageId: mMessage.ID,
				Type:      consts.MessageTypeNormal,
				Sender:    mUser.ID,
				Receiver:  form.Receiver,
				GroupId:   0,
				Data:      form.Content,
				Time:      time.Now().Local().Format(time.DateTime),
			}

			messageJsonByte, err := json.Marshal(message)
			if err != nil {
				helper.Logger.Errorf("消息[%d]Marshal失败：%s", message.MessageId, err.Error())
				continue
			}
			ws.IM.SendMessageByUserId(messageJsonByte, messageUser.Receiver)
		}
	}

	helper.ResponseOkWithData(c, gin.H{
		"message_id": mMessage.ID,
	})
}

type SearchUserForm struct {
	Keyword string `form:"keyword" json:"keyword" binding:"required"`
}

func SearchUser(c *gin.Context) {
	var form SearchUserForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, "参数错误")
		return
	}
	qUser := helper.DbQuery.User
	mUser := chat_model.User{}
	err := qUser.WithContext(c).
		Select(qUser.ID, qUser.UserName, qUser.Nickname, qUser.Phone, qUser.Avatar, qUser.Gender).
		Where(
			qUser.Where(
				qUser.Phone.Like(fmt.Sprintf("%%%s%%", form.Keyword))).
				Or(qUser.Nickname.Like(fmt.Sprintf("%%%s%%", form.Keyword))),
		).Scan(&mUser)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if mUser.ID == 0 {
		helper.ResponseOkWithData(c, gin.H{})
		return
	}
	mUser.Avatar = helper.GenerateStaticUrl(mUser.Avatar)
	helper.ResponseOkWithData(c, mUser)
}

type UserDetailForm struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}
type UserDetailInfo struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增" json:"id"` // 自增
	Phone     string `gorm:"column:phone;not null;comment:用户手机号" json:"phone"`             // 用户手机号
	UserName  string `gorm:"column:user_name;not null;comment:用户名称" json:"user_name"`      // 用户名称
	Nickname  string `gorm:"column:nickname;not null;comment:用户昵称" json:"nickname"`        // 用户昵称
	Gender    int32  `gorm:"column:gender;not null;default:-1;comment:性别" json:"gender"`   // 性别
	Avatar    string `gorm:"column:avatar;not null;comment:头像" json:"avatar"`              // 头像
	AvatarUrl string `json:"avatarUrl"`                                                    // 头像
	IsFriend  int64  `json:"isFriend"`
}

func UserDetail(c *gin.Context) {
	var form UserDetailForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	qUser := helper.DbQuery.User
	mUser := UserDetailInfo{}
	err := qUser.WithContext(c).
		Select(qUser.ID, qUser.UserName, qUser.Nickname, qUser.Phone, qUser.Avatar, qUser.Gender).
		Where(qUser.ID.Eq(form.Id)).Scan(&mUser)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if mUser.ID == 0 {
		helper.ResponseOkWithData(c, gin.H{})
		return
	}
	mUser.AvatarUrl = helper.GenerateStaticUrl(mUser.Avatar)

	//检测是否为好友
	mUser.IsFriend, _ = service.User.IsFriendContact(c, c.GetInt64(consts.UserId), mUser.ID)

	helper.ResponseOkWithData(c, mUser)
}

type ApplyFriendForm struct {
	UserId int64 `form:"userId" json:"userId" binding:"required"`
	Status int   `form:"status" json:"status" binding:"oneof=2 3"`
}

func ApplyFriend(c *gin.Context) {
	var form ApplyFriendForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, "参数错误:"+err.Error())
		return
	}

	//当前用户
	user, err := service.User.GetLoginUser(c)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//申请添加的好友
	friend, err := service.User.GetUserById(form.UserId)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if user.ID == friend.ID {
		helper.ResponseError(c, "不能添加自己为好友")
		return
	}

	//查询好友关系
	qContact := helper.DbQuery.UserContact
	mContact := chat_model.UserContact{}
	err = qContact.WithContext(c).Where(
		qContact.Where(qContact.Where(qContact.UserID.Eq(user.ID), qContact.FriendUserID.Eq(friend.ID))).Or(qContact.Where(qContact.FriendUserID.Eq(user.ID), qContact.UserID.Eq(friend.ID))),
	).Scan(&mContact)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if mContact.Status == consts.UserFriendStatusIsFriend {
		helper.ResponseError(c, "已经是好友了，无需重复添加")
		return
	}
	if mContact.UserID == user.ID && mContact.Status == consts.UserFriendStatusIsApplying {
		helper.ResponseError(c, "您已经申请添加对方为好友，请勿重复操作，请耐心等待您的好友同意")
		return
	}

	notifyMessage := ws.Message{
		Type:     consts.MessageTypeAddFriend,
		Sender:   user.ID,
		Receiver: friend.ID,
		Data:     "",
	}
	err = helper.DbQuery.Transaction(func(tx *chat_query.Query) error {
		quc := tx.UserContact
		//添加自己的数据
		err = quc.Select(quc.UserID, quc.FriendUserID, quc.Status).Create(&chat_model.UserContact{
			UserID:       user.ID,
			FriendUserID: friend.ID,
			Status:       consts.UserFriendStatusIsFriend, //直接添加为好友，暂时去掉审核操作
		})
		if err != nil {
			return err
		}

		//添加好友的数据
		err = quc.Select(quc.UserID, quc.FriendUserID, quc.Status).Create(&chat_model.UserContact{
			UserID:       friend.ID,
			FriendUserID: user.ID,
			Status:       consts.UserFriendStatusIsFriend, //直接添加为好友，暂时去掉审核操作
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	//发送通知给被加好友的用户
	notifyMessage.Data = fmt.Sprintf("用户【%s】添加您为好友", user.Nickname)

	//消息通知
	notifyMessageJson, _ := json.Marshal(&notifyMessage)
	_, _ = ws.HandleMessageSaveAndSend(string(notifyMessageJson), user.ID)

	helper.ResponseOk(c)
}

// GetFriendContact 获取好友联系人
func GetFriendContact(c *gin.Context) {
	mUser, err := service.User.GetLoginUser(c)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	friends, err := service.User.GetFriendContact(c, mUser.ID)
	helper.ResponseOkWithData(c, friends)
}

type CreateGroupForm struct {
	Title string `form:"title" json:"title" binding:"required"`
}

func CreateGroup(c *gin.Context) {
	var form CreateGroupForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	mUser, err := service.User.GetLoginUser(c)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	mGroup := chat_model.Group{
		Title:         form.Title,
		CreatedUserID: mUser.ID,
	}
	err = helper.DbQuery.Transaction(func(tx *chat_query.Query) error {
		qG := tx.Group
		err = qG.WithContext(c).Select(qG.Title, qG.CreatedUserID).Create(&mGroup)
		if err != nil {
			return err
		}

		qGu := tx.GroupUser
		err = qGu.WithContext(c).Select(qGu.GroupID, qGu.UserID, qGu.IsMessageRemind).Create(&chat_model.GroupUser{
			UserID:          mUser.ID,
			GroupID:         mGroup.ID,
			IsMessageRemind: 1,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	helper.ResponseOkWithData(c, gin.H{
		"group_id": mUser.ID,
	})
}

type AddGroupUserForm struct {
	GroupId int64   `form:"group_id" json:"group_id" binding:"required"`
	UserIds []int64 `form:"user_ids" json:"user_ids" binding:"required"`
}

func AddGroupUser(c *gin.Context) {
	var form AddGroupUserForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	var group chat_model.Group
	qGroup := helper.DbQuery.Group
	err := qGroup.WithContext(c).Where(qGroup.ID.Eq(form.GroupId)).Scan(&group)
	if err != nil || group.ID == 0 {
		helper.ResponseError(c, "聊天群信息错误")
		return
	}

	err = helper.DbQuery.Transaction(func(tx *chat_query.Query) error {
		qUser := tx.User
		qGroupUser := tx.GroupUser
		users := make([]*chat_model.User, 0)
		err = qUser.WithContext(c).
			Select(qUser.ID).
			LeftJoin(qGroupUser, qGroupUser.UserID.EqCol(qUser.ID)).
			Where(qGroupUser.GroupID.Eq(form.GroupId)).
			Where(qUser.ID.In(form.UserIds...)).
			Where(qGroupUser.ID.IsNull()).
			Scan(&users)

		mGroupUsers := make([]*chat_model.GroupUser, 0)
		for _, u := range users {
			mGroupUsers = append(mGroupUsers, &chat_model.GroupUser{
				UserID:          u.ID,
				GroupID:         form.GroupId,
				IsMessageRemind: 1,
			})
		}

		qGu := tx.GroupUser
		err = qGu.WithContext(c).Select(qGu.GroupID, qGu.UserID, qGu.IsMessageRemind).Create(mGroupUsers...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	helper.ResponseOk(c)
}

// GetChatList 获取聊天列表
func GetChatList(c *gin.Context) {
	//获取发送给当前用户的好友列表
	userId := c.GetInt64(consts.UserId)

	//获取好友列表
	contacts, _ := service.User.GetFriendContact(c, userId)
	if len(contacts) == 0 {
		helper.ResponseOkWithData(c, gin.H{})
		return
	}

	//获取好友ID
	friends := make([]int64, 0)
	for _, v := range contacts {
		friends = append(friends, v.ID)
	}

	//统计用户的未读消息总数
	unReadMessageUsers, err := service.MessageService.GetUnreadMessageCount(c, userId, friends)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//整理好友发送的消息未读总数
	senderUnreadCount := map[int64]int64{}
	for _, v := range unReadMessageUsers {
		senderUnreadCount[v.Sender] = v.UnreadCount
	}
	fmt.Println(senderUnreadCount)
	//统计发送消息给我的最后一条信息
	lastMessages, err := service.MessageService.GetLastMessage(c, userId)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	lastMessageMap := map[int64]*chat_model.Message{}
	for _, v := range lastMessages {
		v := v
		lastMessageMap[v.Sender] = v
	}

	//统计好友列表中对应的未读消息总数
	for k, v := range contacts {
		contacts[k].AvatarUrl = helper.GenerateStaticUrl(v.Avatar)
		unreadCount, ok1 := senderUnreadCount[v.ID]
		if ok1 {
			contacts[k].UnreadCount = unreadCount
		}

		lastMessage, ok2 := lastMessageMap[v.ID]
		if ok2 && lastMessage.ID > 0 {
			contacts[k].LastMessage = lastMessage
		}
	}
	helper.ResponseOkWithData(c, contacts)
}

type SetMessageReadStatusForm struct {
	Sender  int64 `form:"sender" json:"sender"`
	GroupId int64 `form:"groupId" json:"groupId"`
}

// SetMessageReadStatus 设置消息阅读状态
func SetMessageReadStatus(c *gin.Context) {
	form := SetMessageReadStatusForm{}
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	if form.Sender == 0 && form.GroupId == 0 {
		helper.ResponseError(c, "参数错误")
		return
	}

	userId := c.GetInt64(consts.UserId)
	messageIds := make([]int64, 0)
	messages := make([]*chat_model.Message, 0)
	qm := helper.DbQuery.Message
	if form.GroupId > 0 {
		//群聊消息
		_ = qm.WithContext(c).Where(qm.GroupID.Eq(form.GroupId)).Scan(&messages)
	} else if form.Sender > 0 {
		//私聊消息
		_ = qm.WithContext(c).Where(qm.GroupID.Eq(0), qm.Sender.Eq(form.Sender)).Scan(&messages)
	}

	for _, v := range messages {
		messageIds = append(messageIds, v.ID)
	}

	qmu := helper.DbQuery.MessageUser
	if len(messageIds) > 0 {
		//将属于自己的消息标记为已读
		_, _ = qmu.WithContext(c).Where(qmu.MessageID.In(messageIds...), qmu.Receiver.Eq(userId)).Update(qmu.IsRead, consts.MessageReadStatusYes)
	}

	helper.ResponseOk(c)
}

type GetMessageListFrom struct {
	Sender   []int64 `form:"sender" json:"sender"`
	GroupId  int64   `json:"groupId" form:"groupId"`
	IsRead   int32   `json:"isRead" form:"isRead"`
	Page     int     `json:"page" form:"page"`
	PageSize int     `json:"pageSize" form:"pageSize"`
}
type GetMessageListRes struct {
	Id         int64             `json:"id"`
	Sender     int64             `json:"sender"`
	GroupID    int64             `json:"groupId"`
	Source     int32             `json:"source"`
	Type       int32             `json:"type"`
	Content    string            `json:"content"`
	CreatedAt  string            `json:"createdAt,type:datetime"`
	SenderInfo *structs.UserItem `json:"senderInfo" gorm:"-"`
}

func GetMessageList(c *gin.Context) {
	var form GetMessageListFrom
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	offset := (form.Page - 1) * form.PageSize
	if offset < 0 {
		offset = 0
	}
	qm := helper.DbQuery.Message
	qmu := helper.DbQuery.MessageUser
	quc := helper.DbQuery.UserContact
	list := make([]*GetMessageListRes, 0)
	count, err := qm.WithContext(c).
		Select(qm.ALL).
		Join(qmu, qmu.MessageID.EqCol(qm.ID)).
		Join(quc, quc.FriendUserID.EqCol(qm.Sender)).
		Where(qm.Sender.In(form.Sender...)).
		Order(qm.ID.Asc()).
		ScanByPage(&list, offset, form.PageSize)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//获取消息发送人的信息
	qu := helper.DbQuery.User
	senderUsers := make([]*structs.UserItem, 0)
	err = qu.WithContext(c).Where(qu.ID.In(form.Sender...)).Scan(&senderUsers)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	senderUserMap := map[int64]*structs.UserItem{}
	for _, v := range senderUsers {
		v.AvatarUrl = helper.GenerateStaticUrl(v.Avatar)
		senderUserMap[v.ID] = v
	}

	//更新消息列表的消息发送列表
	for k, v := range list {
		si, sok := senderUserMap[v.Sender]
		if sok {
			list[k].SenderInfo = si
		}

		list[k].CreatedAt = helper.FormatTimeRFC3339ToDatetime(v.CreatedAt)
	}

	helper.ResponseOkWithData(c, gin.H{
		"list":  list,
		"count": count,
	})
}
