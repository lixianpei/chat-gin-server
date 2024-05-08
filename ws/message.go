package ws

import (
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/query/chat_query"
	"GoChatServer/dal/types"
	"GoChatServer/helper"
	"GoChatServer/service"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Message 消息
type Message struct {
	MessageId  int64           `json:"messageId"`  //消息ID
	Type       int32           `json:"type"`       //消息类型
	Sender     int64           `json:"sender"`     //消息发送的用户ID
	Receiver   int64           `json:"receiver"`   //消息接收的用户ID
	RoomId     int64           `json:"roomId"`     //消息关联的群ID
	Content    string          `json:"content"`    //消息内容 消息类型: 1-普通文本消息；2-图片文件；3-语音文件；4-视频文件； 除1类型外，其他三种文件类型消息内容格式：{"attachment_id":7,"filepath":"xxx.pdf"}
	Time       string          `json:"time"`       //消息
	SenderInfo *types.UserItem `json:"senderInfo"` //消息发送人信息
	RoomInfo   *types.RoomInfo `json:"roomInfo"`   //聊天室信息
}

// ToString 对消息格式化
func (m *Message) ToString() (messageStr string) {
	messageByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("message json.Marshal Fail：%s  【%v】", err.Error(), m)
		return ""
	}
	messageStr = string(messageByte)
	return messageStr
}

// HandleMessageSaveAndSend 处理用户消息保存和发送
func HandleMessageSaveAndSend(wsMessage string, sender int64) (messageData Message, err error) {
	err = json.Unmarshal([]byte(wsMessage), &messageData)
	if err != nil {
		return messageData, err
	}
	messageData.Content = strings.TrimSpace(messageData.Content)
	if messageData.RoomId == 0 || len(messageData.Content) == 0 {
		helper.Logger.Errorf("消息格式错误，不能正确处理: %s", wsMessage)
		return messageData, fmt.Errorf("消息格式错误，不能正确处理")
	}

	mSenderInfo, err := service.User.GetMessageUserById(sender)
	if err != nil {
		return messageData, fmt.Errorf("用戶不存在")
	}

	//消息内容
	mMessage := chat_model.Message{
		Sender:  mSenderInfo.ID,
		RoomID:  messageData.RoomId,
		Content: messageData.Content,
		Type:    messageData.Type,
	}
	fmt.Printf("mMessage....... %+v", mMessage)
	//消息关联的用户
	messageUsers := make([]*chat_model.MessageUser, 0)

	//消息入库
	err = helper.DbQuery.Transaction(func(tx *chat_query.Query) error {
		//保存消息
		err = tx.Message.Create(&mMessage)
		if err != nil {
			return err
		}

		//查询消息关联的用户
		users, err := service.User.GetMessageReceiverUsers(messageData.RoomId)
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
			err = tx.MessageUser.Create(messageUsers...)
			if err != nil {
				return err
			}
		}

		//消息关联的聊天会话数据更新
		qr := tx.Room
		_, err = qr.Where(qr.ID.Eq(messageData.RoomId)).Update(qr.LastMessageID, mMessage.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return messageData, err
	}

	//扩展ws的消息内容
	messageData.Time = time.Now().Local().Format(time.DateTime)
	messageData.MessageId = mMessage.ID
	messageData.SenderInfo = mSenderInfo
	messageData.Sender = sender

	qr := helper.DbQuery.Room
	messageData.RoomInfo = &types.RoomInfo{}
	_ = qr.Where(qr.ID.Eq(messageData.RoomId)).Scan(&messageData.RoomInfo)

	messageData.Content = helper.FormatFileMessageContent(messageData.Type, messageData.Content)

	// 根据消息关联的用户发送消息
	messageJson, _ := json.Marshal(messageData)
	for _, messageUser := range messageUsers {
		messageReceiverUserId := messageUser.Receiver
		//TODO 暂时针对自己发送的消息，不推送给自己，前端视为直接发送成功
		if messageReceiverUserId == messageData.Sender {
			continue
		}
		go IM.SendMessageByUserId(messageJson, messageReceiverUserId)
	}

	return messageData, nil
}
