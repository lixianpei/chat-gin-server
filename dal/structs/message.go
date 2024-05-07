package structs

type MessageListItem struct {
	MessageId  int64                `gorm:"column:messageId" json:"messageId"`
	Sender     int64                `json:"sender"`
	Source     int32                `json:"source"`
	Type       int32                `json:"type"`
	Content    string               `json:"content"`
	CreatedAt  string               `gorm:"created_at" json:"createdAt"`
	SenderInfo *MessageListUserItem `json:"senderInfo" gorm:"-"`
}

type MessageListUserItem struct {
	UserID    int64  `gorm:"column:userId" json:"userId"` // 自增
	Phone     string `json:"phone"`                       // 用户手机号
	UserName  string `json:"user_name"`                   // 用户名称
	Nickname  string `json:"nickname"`                    // 用户昵称
	Gender    int32  `json:"gender"`                      // 性别
	Avatar    string `json:"avatar"`                      // 头像
	AvatarUrl string `json:"avatarUrl" gorm:"-"`
}
