package api

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/helper"
	"GoChatServer/ws"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Phone    string `form:"phone"`
	Nickname string `form:"nickname"`
	Code     string `form:"code"`
}

// WxLogin 微信登录
func WxLogin(c *gin.Context) {
	var loginForm LoginForm
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
	err = helper.Db.WithContext(c).User.
		Select(qUser.ID, qUser.UserName, qUser.Nickname).
		Where(qUser.WxOpenid.Eq(wxResult.OpenId)).
		Scan(mUserInfo)
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

	//TODO
	loginForm.Phone = "18083198680"
	loginForm.Nickname = "LiXianPei"

	//生成token
	token, err := helper.NewJwtToken(mUserInfo.ID, loginForm.Phone, loginForm.Nickname)
	if err != nil {
		helper.ResponseError(c, err.Error())
		return
	}

	helper.ResponseOkWithMessageData(c, gin.H{
		"user_id":   mUserInfo.ID,
		"token":     token,
		"phone":     loginForm.Phone,
		"nickname":  loginForm.Nickname,
		"wxResult":  wxResult,
		"mUserInfo": mUserInfo,
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

type WxUserAvatarForm struct {
	Avatar string `form:"avatar"`
}

// WxUserAvatarSave 微信头像存储-头像为临时头像，暂时不需要此接口
func WxUserAvatarSave(c *gin.Context) {
	var form WxUserAvatarForm
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

	//保存用户信息
	updateUser := chat_model.User{
		Avatar: form.Avatar,
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
