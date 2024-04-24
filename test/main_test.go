package test

import (
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/helper"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化配置
	helper.InitConfig("../config")
	helper.InitLogger()

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
