package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//const JwtSigningKeyStr = "AllYourBase"

const JwtSigningKeyStr = "Hkuq4oEDYLrb7ghygyDPEyoLHmEwT8nvmMX5jS8BrJXt4tei0Vf1speBiDlLcxuM"

type MyCustomJwtClaims struct {
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	UserId   int64  `json:"user_id"`
	jwt.StandardClaims
}

// NewJwtToken 生成Token
func NewJwtToken(userId int64, phone string, nickname string) (string, error) {
	mySigningKey := []byte(JwtSigningKeyStr)

	// Create the Claims
	claims := MyCustomJwtClaims{
		UserId:   userId,
		Phone:    phone,
		Nickname: nickname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "chatGinServer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// JwtParseChecking 解析token的数据
func JwtParseChecking(tokenString string) (*MyCustomJwtClaims, error) {
	if len(tokenString) == 0 {
		return nil, fmt.Errorf("the token cannot be empty")
	}
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSigningKeyStr), nil
	})
	if err != nil {
		return nil, err
	}

	// 提取声明
	if claims, ok := token.Claims.(*MyCustomJwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
