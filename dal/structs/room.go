package structs

type RoomListItem struct {
	RoomId        int64            `gorm:"column:roomId" json:"roomId"`
	Title         string           `json:"title"`
	Type          int32            `json:"type"`
	UserCount     int32            `json:"userCount"`
	LastMessageId int64            `json:"lastMessageId"`
	LastMessage   *MessageListItem `json:"lastMessage" gorm:"-"`
	RoomUsers     []*RoomUserItem  `json:"users" gorm:"-"`
	AvatarUrls    []string         `gorm:"-" json:"avatarUrls"`
}

type RoomUserItem struct {
	RoomId    int64  `gorm:"column:roomId" json:"roomId"`
	UserID    int64  `json:"userId"`   // 用户ID
	UserName  string `json:"userName"` // 用户名称
	Nickname  string `json:"nickname"` // 用户昵称
	Gender    int32  `json:"gender"`   // 性别
	Avatar    string `json:"avatar"`   // 头像
	AvatarUrl string `json:"avatarUrl" gorm:"-"`
}

type RoomInfo struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Type          int32  `json:"type"`
	UserCount     int32  `json:"userCount"`
	LastMessageId int64  `json:"lastMessageId"`
}
