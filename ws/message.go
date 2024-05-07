package ws

import (
	"GoChatServer/consts"
	"GoChatServer/dal/model/chat_model"
	"GoChatServer/dal/query/chat_query"
	"GoChatServer/dal/structs"
	"GoChatServer/helper"
	"GoChatServer/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// Message 消息
type Message struct {
	MessageId  int64             `json:"messageId"`  //消息ID
	Type       int32             `json:"type"`       //消息类型
	Sender     int64             `json:"sender"`     //消息发送的用户ID
	Receiver   int64             `json:"receiver"`   //消息接收的用户ID
	RoomId     int64             `json:"roomId"`     //消息关联的群ID
	Data       string            `json:"content"`    //消息内容
	Time       string            `json:"time"`       //消息
	SenderInfo *structs.UserItem `json:"senderInfo"` //消息发送人信息 TODO 后期关键信息去掉
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

// NewMessageText 实例化一个文本类型的消息
func NewMessageText(messageId int64, data string, sender int64, receiver int64, roomId int64) *Message {
	return &Message{
		MessageId: messageId,
		Type:      consts.MessageTypeEntryGroup,
		Sender:    sender,
		Receiver:  receiver,
		RoomId:    roomId,
		Data:      data,
		Time:      time.Now().Local().Format(time.DateTime),
	}
}

// NewEntryGroupMessage 返回一个入群消息
func NewEntryGroupMessage(data string) string {
	message := NewMessageText(0, data, 0, 0, 0)
	return message.ToString()
}

// HandleMessageSaveAndSend 处理用户消息保存和发送
func HandleMessageSaveAndSend(wsMessage string, sender int64) (messageData Message, err error) {
	err = json.Unmarshal([]byte(wsMessage), &messageData)
	if err != nil {
		return messageData, err
	}

	mSenderInfo, err := service.User.GetMessageUserById(sender)
	if err != nil {
		return messageData, fmt.Errorf("用戶不存在")
	}

	//消息内容
	mMessage := chat_model.Message{
		Sender:  mSenderInfo.ID,
		RoomID:  messageData.RoomId,
		Content: messageData.Data,
		Type:    consts.MessageTypeNormal,
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
		users, err := service.User.GetMessageReceiverUsers(messageData.RoomId, messageData.Receiver)
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

	// 根据消息关联的用户发送消息
	messageJson, _ := json.Marshal(messageData)
	for _, messageUser := range messageUsers {
		messageReceiverUserId := messageUser.Receiver
		fmt.Println("消息准备发送给好友：", sender, messageReceiverUserId, string(messageJson))
		//判断是否为好友，不是好友不给发私聊消息
		if yes, _ := service.User.IsFriendContact(&gin.Context{Request: nil}, sender, messageReceiverUserId); yes > 0 {
			fmt.Println("消息准备发送给好友：ok....")
			go IM.SendMessageByUserId(messageJson, messageReceiverUserId)
		} else {
			helper.Logger.Infof("消息接收用户还不是好友[%d]，无法给用户发送消息: %s", messageReceiverUserId, messageJson)
		}
	}

	return messageData, nil
}
