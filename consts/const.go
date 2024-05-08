package consts

// 时间格式化模板
const (
	DateYMDHIS = "20060102150405"
	DateYMD    = "20060102"
	DateY      = "2006"
)

// 登录用户的Context上下文缓存key
const (
	UserId       = "JwtUserId"
	UserPhone    = "JwtUserPhone"
	UserNickname = "JwtUserNickname"
)

// 链路跟踪
const (
	TraceId  = "traceId"
	TraceSql = "traceSql"
)

// 好友关系
const (
	UserFriendStatusIsApplying = 1 //好友申请中
	UserFriendStatusIsFriend   = 2 //好友
	UserFriendStatusIsReject   = 3 //拒绝
)

// 消息是否已读
const (
	MessageReadStatusNo  = 0 //未读
	MessageReadStatusYes = 1 //已读
)

const (
	MessageTypeText  = 1 //消息类型-普通文本消息
	MessageTypeImage = 2 //消息类型-图片文件
	MessageTypeAudio = 3 //消息类型-语音文件
	MessageTypeVideo = 4 //消息类型-视频文件
)

const (
	RoomTypeSingle = 1 //私聊
	RoomTypeGroup  = 2 //群聊
)

const (
	RoomUserIsMessageRemindYes = 1 //消息提醒
	RoomUserIsMessageRemindNo  = 2 //消息免打扰
	RoomUserIsTopYes           = 1 //置顶
	RoomUserIsTopNo            = 2 //取消置顶
)
