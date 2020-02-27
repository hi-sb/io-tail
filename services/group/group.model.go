package group

import "github.com/hi-sb/io-tail/core/db"

// 群聊设置
type GroupModel struct {
	db.BaseModel

	// 群名称
	GroupName string

	// 群公告
	GroupAnnouncement string

	// 群主
	GreateUserID String `gorm:"type:varchar(32);not null"`

	// 群聊天状态 0:全体禁言  1:正常
	GroupChatStatus int `gorm:"type:int(2);not null;default:1"`


}
