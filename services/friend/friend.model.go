package friend

import "github.com/hi-sb/io-tail/core/db"

// 好友模型
type FriendModel struct {
	db.BaseModel
	// 用户ID
	UserID string `gorm:"type:varchar(32);not null"`

	// 好友ID
	FriendID string `gorm:"type:varchar(32);not null"`

	// 是否拉黑 11:正常   10 userId 拉黑 friendId   01 friendId 拉黑userId
	IsBlack int `gorm:"type:int(4);default:11"`

	// 添加好友方向 10: userid - friendId   01 friendId-userID
	FriendFrom int `gorm:"type:int(4);not null"`

	// 是否同意添加好友  0： 不同意   1：同意   2：待确认
	IsAgree int `gorm:"type:int(4);default:2"`
}
