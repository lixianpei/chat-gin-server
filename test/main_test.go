package test

import (
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化配置
	helper.InitConfig("../config")
	helper.InitLogger()
	helper.InitRedis()

	helper.InitChatDatabase()

	fmt.Println("TestMain...")
	os.Exit(m.Run())
}

func TestDb(t *testing.T) {

	qUser := helper.Db.User
	mUser := chat_model.User{}
	err := helper.Db.User.Where(qUser.WxOpenid.Eq("dddd")).Scan(&mUser)
	fmt.Println(err, mUser)
	fmt.Println(mUser.ID)
}

func TestRedis(t *testing.T) {
	c := &gin.Context{Request: &http.Request{}}
	//fmt.Println(helper.Redis.Lock(c, "a", time.Hour))
	//fmt.Println(helper.Redis.Lock(c, "b", time.Hour))
	//fmt.Println(helper.Redis.Lock(c, "c", time.Hour))

	//fmt.Println(helper.Redis.Lock(c, "a1", time.Hour))

	fmt.Println(helper.Redis.Del(c, "a", "a1"))
}
