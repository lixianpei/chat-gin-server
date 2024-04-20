package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var WxApi *weiXin

type weiXin struct {
	appid     string //小程序 appId
	secret    string //小程序 appSecret
	jsCode    string //登录时获取的 code，可通过wx.login获取
	grantType string //授权类型，此处只需填写 authorization_code
}

func InitWeiXin() {
	WxApi = newWeiXin(Configs.WeiXin.Appid, Configs.WeiXin.Secret)
}

func newWeiXin(appid string, secret string) *weiXin {
	return &weiXin{
		appid:     appid,
		secret:    secret,
		jsCode:    "",
		grantType: "authorization_code",
	}
}

type LoginResult struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Login 通过code调用微信登录换取用户信息
func (w *weiXin) Login(jsCode string) (*LoginResult, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", w.appid, w.secret, jsCode)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close
	}()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	result := &LoginResult{}
	err = json.Unmarshal(body, result)

	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf(result.ErrMsg)
	}
	return result, nil
}

// CheckWxSignature 验证数据是否被篡改
func (w *weiXin) CheckWxSignature(rawData string, sessionKey string, signature string) bool {
	hasher := sha1.New()
	hasher.Write([]byte(rawData + sessionKey))
	hashed := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashed)
	return signature == hashedString
}

// DecodeWxData 解密微信数据
func (w *weiXin) DecodeWxData(encryptedData string, sessionKey string, iv string) (decodeDataString string) {
	// 从接口返回的数据
	//encryptedData := "接口返回的加密数据"
	//sessionKey := "接口返回的session_key"
	//iv := "接口返回的iv"

	// Base64 解码
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		fmt.Println("Error decoding encrypted data:", err)
		return
	}

	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		fmt.Println("Error decoding session key:", err)
		return
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		fmt.Println("Error decoding IV:", err)
		return
	}

	// 创建一个 AES 解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher:", err)
		return
	}

	// CBC 解密模式
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	// 解密数据
	decryptedData := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(decryptedData, encryptedBytes)

	// PKCS#7 去填充
	decryptedData = PKCS7Unpad(decryptedData)

	// 输出解密后的数据
	return string(decryptedData)
}

// PKCS7Unpad 去除 PKCS#7 填充
func PKCS7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
