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
	Code string `form:"code" binding:"required"`
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
	fmt.Println(err, mUserInfo)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	if mUserInfo.ID == 0 {
		//创建新用户
		mUserInfo = &chat_model.User{
			WxOpenid:     wxResult.OpenId,
			WxUnionid:    wxResult.UnionId,
			WxSessionKey: wxResult.SessionKey,
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
			WxUnionid:    wxResult.UnionId,
			WxSessionKey: wxResult.SessionKey,
		})
	}

	//生成token，此时还未获取到昵称和用户名
	token, err := helper.NewJwtToken(mUserInfo.ID, "", "")
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	//处理头像
	mUserInfo.Avatar = helper.GenerateStaticUrl(mUserInfo.Avatar)

	helper.ResponseOkWithMessageData(c, gin.H{
		"user_id":  mUserInfo.ID,
		"token":    token,
		"userInfo": mUserInfo,
	}, "ok")
}

type PhoneLoginForm struct {
	Phone    string `form:"phone" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
	Avatar   string `form:"avatar"`
}

// PhoneLogin 手机号登录
func PhoneLogin(c *gin.Context) {
	var loginForm PhoneLoginForm
	err := c.ShouldBind(&loginForm)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
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
			Phone:    loginForm.Phone,
			Nickname: loginForm.Nickname,
			Avatar:   loginForm.Avatar,
		}
		err = helper.Db.WithContext(c).User.Create(mUserInfo)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	} else {
		//保存信息数据
		_, err = helper.Db.WithContext(c).User.Where(qUser.ID.Eq(mUserInfo.ID)).Updates(&chat_model.User{
			Nickname: loginForm.Nickname,
			Phone:    loginForm.Phone,
			Avatar:   loginForm.Avatar,
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
	EncryptedData string `form:"encryptedData"`
	RawData       string `form:"rawData"`
	Signature     string `form:"signature"`
	Iv            string `form:"iv"`
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

	fmt.Println("updateUser......", updateUser.Avatar, updateUser.Nickname, updateUser.Gender)

	//返回数据
	helper.ResponseOkWithData(c, gin.H{
		"wxUserForm": form,
		"user_id":    c.GetInt64(consts.UserId),
	})
}

type UserAvatarForm struct {
	Avatar   string `form:"avatar"`
	Nickname string `form:"nickname"`
	Phone    string `form:"phone"`
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
	Content  string `form:"content" binding:"required"`  //消息类型
	Receiver int64  `form:"receiver" binding:"required"` //消息接收的用户ID
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
	Keyword string `form:"keyword" binding:"required"`
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

type AddFriendForm struct {
	UserId int64 `form:"user_id" binding:"required"`
}

func AddFriend(c *gin.Context) {
	var form AddFriendForm
	if err := c.ShouldBind(&form); err != nil {
		helper.ResponseError(c, "参数错误2")
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

	//查询好友关系
	//contact
	qContact := helper.Db.UserContact
	mContact := chat_model.UserContact{}
	err = qContact.Where(
		qContact.Where(qContact.UserID.Eq(user.ID)).Or().Where(qContact.FriendUserID.Eq(friend.ID)),
		qContact.Where(qContact.FriendUserID.Eq(user.ID)).Or().Where(qContact.UserID.Eq(friend.ID)),
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
	if mContact.FriendUserID == user.ID && mContact.Status == consts.UserFriendStatusIsApplying {
		//同意对方的申请添加好友
		_, err = qContact.WithContext(c).Select(qContact.Status).Where(qContact.ID.Eq(mContact.ID)).Update(qContact.Status, consts.UserFriendStatusIsFriend)
		if err != nil {
			helper.ResponseError(c, err.Error())
			return
		}
	}
	//添加好友申请
	err = qContact.WithContext(c).Select(qContact.UserID, qContact.FriendUserID, qContact.Status).Create(&chat_model.UserContact{
		UserID:       user.ID,
		FriendUserID: friend.ID,
		Status:       consts.UserFriendStatusIsApplying,
	})
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}
	helper.ResponseOkWithMessage(c, "添加好友处理成功，请耐心等待您的好友同意")

}
