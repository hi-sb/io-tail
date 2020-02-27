package group

import "github.com/hi-sb/io-tail/core/db"

// 群成员
type GroupMemberModel struct {
	db.BaseModel

	// 群ID
	GroupID string `gorm:"type:varchar(32);not null"`

	// 成员
	GroupMermerID string  `gorm:"type:varchar(32);not null"`

	// 成员在当前群的昵称
	GroupMermerNickName string  `gorm:"type:varchar(255)"`

	// 成员角色  0: 普通成员 1.群主  2。管理员
	GroupMemberRole int `gorm:"type:int(2);not null;default:0"`

}
