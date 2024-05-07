package test

import (
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

	//helper.InitChatDatabase()

	fmt.Println("TestMain...")
	os.Exit(m.Run())
}

type messageUser struct {
	Sender       int64
	UnreadCount  int64
	MaxMessageId int64
}

func TestDb(t *testing.T) {
	//fmt.Println(helper.FormatTimeToDatetime("2024-05-02T15:24:30+08:00"))

	//c := &gin.Context{Request: &http.Request{}}

	//res, err := service.User.IsFriendContact(c, 1, 3)
	//fmt.Println(res, err)

	//qMU := helper.DbQuery.MessageUser
	//qM := helper.DbQuery.Message
	//messageUsers := make([]*messageUser, 0)
	//err := qM.
	//	Join(qMU, qMU.MessageID.EqCol(qM.ID)).
	//	Select(qM.Sender, qMU.IsRead.Count().As("unread_count"), qM.ID.Max().As("max_message_id")).
	//	Where(qMU.Receiver.Eq(1)). //
	//	Where(qM.Sender.In(2)).
	//	Group(qM.Sender).
	//	Order(qMU.IsRead.Count().Desc()).
	//	Scan(&messageUsers)
	//fmt.Println(err)
	//fmt.Println(messageUsers)

	//qUser := helper.DbQuery.User
	//mUser := chat_model.User{}
	//sql := "SELECT `user`.`id`,`user`.`user_name`,`user`.`nickname`,`user`.`phone`,`user`.`avatar` " +
	//	"FROM `user` INNER JOIN `user_contact` ON (`user_contact`.`user_id` = `user`.`id` OR user_contact.`friend_user_id` = user.`id`) " +
	//	"WHERE user.`id` != ? AND `user`.`deleted_at` IS NULL"
	//res := helper.Db.Raw(sql, 1).Scan(&mUser)
	//if res.Error != nil {
	//	fmt.Println(res.Error)
	//	return
	//}
	//fmt.Println(mUser)
}

func TestRedis(t *testing.T) {
	c := &gin.Context{Request: &http.Request{}}
	//fmt.Println(helper.Redis.Lock(c, "a", time.Hour))
	//fmt.Println(helper.Redis.Lock(c, "b", time.Hour))
	//fmt.Println(helper.Redis.Lock(c, "c", time.Hour))

	//fmt.Println(helper.Redis.Lock(c, "a1", time.Hour))

	fmt.Println(helper.Redis.Del(c, "a", "a1"))
}
