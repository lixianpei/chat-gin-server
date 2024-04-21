package ws

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	MessageTypeNormal     = "text"       //消息类型-普通文本消息
	MessageTypeEntryGroup = "entryGroup" //消息类型-用户加入群聊消息
	MessageTypeAddFriend  = "addFriend"  //消息类型-加好友消息
	MessageTypeBinary     = "binary"     //消息类型-二进制类型
	MessageTypeUserEntry  = "userEntry"  //消息类型-用户上线
	MessageTypeUserExit   = "userExit"   //消息类型-用户下线
)

// Message 消息
type Message struct {
	MessageId    int64  `json:"message_id"`    //消息ID
	Type         string `json:"type"`          //消息类型
	Sender       int64  `json:"sender"`        //消息发送的用户ID
	Receiver     int64  `json:"receiver"`      //消息接收的用户ID
	GroupId      int64  `json:"group_id"`      //消息关联的群ID
	Data         string `json:"content"`       //消息内容
	Time         string `json:"time"`          //消息
	SenderAvatar string `json:"sender_avatar"` //消息发送人的头像
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
func NewMessageText(messageId int64, data string, sender int64, receiver int64, groupId int64) *Message {
	return &Message{
		MessageId: messageId,
		Type:      MessageTypeEntryGroup,
		Sender:    sender,
		Receiver:  receiver,
		GroupId:   groupId,
		Data:      data,
		Time:      time.Now().Local().Format(time.DateTime),
	}
}

// NewEntryGroupMessage 返回一个入群消息
func NewEntryGroupMessage(data string) string {
	message := NewMessageText(0, data, 0, 0, 0)
	return message.ToString()
}
