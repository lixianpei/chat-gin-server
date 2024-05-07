package service

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/helper"
	"github.com/gin-gonic/gin"
)

var MessageService = &ms{}

type ms struct{}

type UnreadMessageUsers struct {
	Sender       int64
	UnreadCount  int64
	MaxMessageId int64
}

// GetUnreadMessageCount 获取未读消息总数
func (m *ms) GetUnreadMessageCount(c *gin.Context, userId int64, friends []int64) ([]*UnreadMessageUsers, error) {
	qm := helper.DbQuery.Message
	qmu := helper.DbQuery.MessageUser
	messageUsers := make([]*UnreadMessageUsers, 0)
	err := qm.WithContext(c.Request.Context()).
		Join(qmu, qmu.MessageID.EqCol(qm.ID)).
		Select(qm.Sender, qmu.IsRead.Count().As("unread_count"), qm.ID.Max().As("max_message_id")).
		Where(qmu.Receiver.Eq(userId)). //
		Where(qm.Sender.In(friends...)).
		Where(qmu.IsRead.Eq(consts.MessageReadStatusNo)).
		Group(qm.Sender).
		Order(qmu.IsRead.Count().Desc()).
		Scan(&messageUsers)
	return messageUsers, err
}

// GetLastMessage 获取最后一条未读消息
func (m *ms) GetLastMessage(c *gin.Context, userId int64) ([]*chat_model.Message, error) {
	messages := make([]*chat_model.Message, 0)
	sql := "SELECT m.* FROM message m WHERE m.id IN (SELECT MAX(m.`id`) AS max_message_id FROM message m INNER JOIN message_user mu ON m.id = mu.`message_id` WHERE mu.`receiver` = ? GROUP BY m.`sender`)"
	res := helper.Db.WithContext(c).Raw(sql, userId).Scan(&messages)
	return messages, res.Error
}
