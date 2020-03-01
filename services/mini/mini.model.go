package mini

import "github.com/hi-sb/io-tail/core/db"

// 小程序模型
type MiniModel struct {

	db.BaseModel

	// 小程序图标
	MiniLogo string

	// 小程序名称
	MiniName string

	// 小程序地址
	MiniAddress string

	// 游戏介绍
	MiniDesc string

	//备注
	MiniRemark string

	// 状态 1:启用 0：停用
	MiniStatus  int `gorm:"type:int(2);not null;default:1"`

	// 排序
	MiniSort int `gorm:"type:int(2);not null;default:0"`
}
