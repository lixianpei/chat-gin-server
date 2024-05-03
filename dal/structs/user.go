package structs

// UserItem 用户信息
type UserItem struct {
	ID            int64  `json:"id"`              // 自增
	Phone         string `json:"phone"`           // 用户手机号
	UserName      string `json:"user_name"`       // 用户名称
	Nickname      string `json:"nickname"`        // 用户昵称
	Gender        int32  `json:"gender"`          // 性别
	Avatar        string `json:"avatar"`          // 头像
	LastLoginTime string `json:"last_login_time"` // 最后登录时间
	UnreadCount   int64  `json:"unread_count"`
	AvatarUrl     string `json:"avatarUrl" gorm:"-"`
}
