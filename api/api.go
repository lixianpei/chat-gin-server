package api

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/query/chat_query"
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
	qUser := helper.Db.User
	mUserInfo := &chat_model.User{}
	err = helper.Db.WithContext(c).User.Where(qUser.WxOpenid.Eq(wxResult.OpenId)).Scan(mUserInfo)
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
		err = helper.Db.WithContext(c).User.Create(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	} else {
		//保存信息数据
		_, err = helper.Db.WithContext(c).User.Where(qUser.ID.Eq(mUserInfo.ID)).Updates(&chat_model.User{
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

	qUser := helper.Db.User
	mUserInfo := &chat_model.User{}
	//若用户已提前使用微信登录，则此时已经存在token，可以获取到登录的用户信息
	loginUser, err := service.User.GetLoginUser(c)

	if loginUser != nil {
		mUserInfo = loginUser
	} else {
		//直接使用手机号进行登录：根据手机号查询用户是否存在
		err = helper.Db.WithContext(c).User.
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
		err = helper.Db.WithContext(c).User.Create(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	} else {
		//保存信息数据
		_, err = helper.Db.WithContext(c).User.Where(qUser.ID.Eq(mUserInfo.ID)).Updates(&chat_model.User{
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
	qUser := helper.Db.User
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
	qUser := helper.Db.User
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
	err = helper.Db.Transaction(func(tx *chat_query.Query) error {
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
				Type:      ws.MessageTypeNormal,
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
	qUser := helper.Db.User
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

	qUser := helper.Db.User
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
	qContact := helper.Db.UserContact
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
		Type:     ws.MessageTypeAddFriend,
		Sender:   user.ID,
		Receiver: friend.ID,
		Data:     "",
	}
	if mContact.ID > 0 && mContact.Status == consts.UserFriendStatusIsApplying {
		//存在申请中的记录视为在处理好友的申请
		if form.Status == consts.UserFriendStatusIsReject {
			//拒绝好友申请，直接删除关联数据，就是这么简单粗暴！！！
			_, err = qContact.WithContext(c).Where(qContact.ID.Eq(mContact.ID)).Delete()
		} else {
			//同意对方的申请添加好友
			_, err = qContact.WithContext(c).Select(qContact.Status).Where(qContact.ID.Eq(mContact.ID)).Update(qContact.Status, form.Status)
		}
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
		//发送通知给申请加好友的用户
		var statusDesc string
		switch form.Status {
		case consts.UserFriendStatusIsFriend:
			statusDesc = "已同意添加为好友"
		case consts.UserFriendStatusIsReject:
			statusDesc = "已拒绝添加为好友"
		}
		notifyMessage.Data = fmt.Sprintf("您申请添加用户【%s】为好友的请求处理完成：%s", friend.Nickname, statusDesc)

	} else {
		//添加好友申请
		err = qContact.WithContext(c).Select(qContact.UserID, qContact.FriendUserID, qContact.Status).Create(&chat_model.UserContact{
			UserID:       user.ID,
			FriendUserID: friend.ID,
			Status:       consts.UserFriendStatusIsFriend, //直接添加为好友，暂时去掉审核操作
		})
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
		//发送通知给被加好友的用户
		notifyMessage.Data = fmt.Sprintf("用户【%s】请求添加您为好友", user.Nickname)
	}

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

	friends, err := service.User.GetFriendContact(mUser.ID)
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
	err = helper.Db.Transaction(func(tx *chat_query.Query) error {
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
	qGroup := helper.Db.Group
	err := qGroup.WithContext(c).Where(qGroup.ID.Eq(form.GroupId)).Scan(&group)
	if err != nil || group.ID == 0 {
		helper.ResponseError(c, "聊天群信息错误")
		return
	}

	err = helper.Db.Transaction(func(tx *chat_query.Query) error {
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

//[mysql] 2024/04/25 21:57:16 connection.go:49: read tcp 127.0.0.1:62820->127.0.0.1:3306: wsarecv: An established connection was aborted by the software in your host machine.
