package constants


const (
	// key
	USER_BASE_INFO_REDIS_KEY = "USER_BASE_INFO"
	//  field
	USER_BASE_INFO_REDIS_PREFIX = "USER_BASE_INFO_"


	IS_NOT_BLACK string = "11"  // 正常
	IS_BLACK_F_PULL_U string = "10"  // f 拉黑 u
	IS_BLACK_U_PULL_F string = "01"  // u 拉黑 f
	IS_BLACK_EACH_OTHER string = "00" // 互相拉黑

	AGREE_ADD int = 11   // 互为好友
	NOT_AGREE_ADD int = 10   // 对方拒绝 删除记录
	WAITING_AGREE int = 13  // 等待同意


	// 好友列表redisKey前缀
	FRIEND_REDIS_PREFIX = "IO_TAIL_FRIEND_%s"
	// 好友黑名单（发送消息给某个好友查询是否被拉黑）
	FRIEND_BLACK_REDIS_PREFIX = "IO_TAIL_FRIEND_BLACK_%s"


	// 群基础信息
	GROUP_BASE_INFO_REDIS_PREFIX = "GROUP_BASE_INFO_%s"
	// 群成员
	GROUP_MEMBER_INFO_REDIS_PREFIX = "GROUP_MEMBER_INFO_%s"

	// user source
	PRIVATE_SOURCE = "private_source"
	// open source
	PUBLIC_SOURCE = "public_source"


	// 小程序缓存key
	MINI_PROGRAM_HKEY = "MINI_PROGRAM"


)