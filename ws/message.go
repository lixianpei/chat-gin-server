package ws

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	EventMessage    = "message"    //websocket事件类型-消息
	EventEntryGroup = "entryGroup" //websocket事件类型-用户加入群聊
	EventAddFriend  = "addFriend"  //websocket事件类型-加好友
	EventEnter      = "userEnter"  //websocket事件类型-用户上线
	EventExit       = "userExit"   //websocket事件类型-用户离线

	MessageTypeText   = "text"   //消息类型-文本类型
	MessageTypeBinary = "binary" //消息类型-二进制类型
)

// Message 消息
type Message struct {
	Event     string `json:"event"`      //事件
	MessageId int64  `json:"message_id"` //消息ID
	Type      string `json:"type"`       //消息类型
	Sender    string `json:"sender"`     //消息发送的用户ID
	Receiver  string `json:"receiver"`   //消息接收的用户ID
	GroupId   int64  `json:"group_id"`   //消息关联的群ID
	Data      string `json:"content"`    //消息内容
	Time      string `json:"time"`       //消息
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
func NewMessageText(event string, messageId int64, data string, sender string, receiver string, groupId int64) *Message {
	return &Message{
		Event:     event,
		MessageId: messageId,
		Type:      MessageTypeText,
		Sender:    sender,
		Receiver:  receiver,
		GroupId:   groupId,
		Data:      data,
		Time:      time.Now().Local().Format(time.DateTime),
	}
}

// NewEntryGroupMessage 返回一个入群消息
func NewEntryGroupMessage(data string) string {
	message := NewMessageText(EventEntryGroup, 0, data, "0", "0", 0)
	return message.ToString()
}
