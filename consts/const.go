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
