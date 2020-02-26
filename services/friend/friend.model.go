package friend

import (
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/syserr"
)


const (
	IS_NOT_BLACK string = "11"  // 正常
	IS_BLACK_F_PULL_U string = "10"  // u 拉黑 f
	IS_BLACK_U_PULL_F string = "01"  // f 拉黑 u
	IS_BLACK_EACH_OTHER string = "00" // 互相拉黑

	AGREE_ADD int = 11   // 互为好友
	NOT_AGREE_ADD int = 10   // 对方拒绝 删除记录
	WAITING_AGREE int = 13  // 等待同意

)


// 好友模型
type FriendModel struct {
	db.BaseModel
	// 用户ID
	UserID string `gorm:"type:varchar(32);not null"`

	// 好友ID
	FriendID string `gorm:"type:varchar(32);not null"`

	// 是否拉黑  userid -> friendId   11:正常   10 userId 拉黑 friendId   01 friendId 拉黑userId  00：互相拉黑
	IsBlack string `gorm:"type:varchar(4);default:11";json:"-"`

	// 是否同意添加好友   userid -> friendId   10 f拒  11 互为好友 13等待确认
	IsAgree int `gorm:"type:int(4);default:13";json:"-"`
	// 好友备注
	UtoFRemark string `gorm:"type:varchar(32)"`  // userid -> friendId 的备注
	FtoURemark string `gorm:"type:varchar(32)`   // friendId --> userId 的备注
}


// 添加好友参数验证
func (f *FriendModel) Check() error{
	if f.FriendID == "" || len(f.FriendID) != 32 {
		return syserr.NewParameterError("参数不正确")
	}
	// 防止自己添加自己
	if f.FriendID == f.UserID {
		return syserr.NewParameterError("好友信息不对，请确认")
	}
	return nil
}

// 更新好友添加请求体
type UpdateAddFReqModel struct {
	ID string
	State int  // 1:同意  0:拒绝
	ReqId string   // 请求添加者
	FtoURemark string
}

// 拉黑/还原拉黑好友模型
type PullBlackModel struct {
	FriendID string
	IsBlack int  // 0 拉黑  1 正常
}

// 好友请求返回模型
type FriendAddReqModel struct {
	FriendID string
	// 手机号
	MobileNumber string
	//昵称
	NickName string
	// 头像
	Avatar string
	//备注
	Remark string
	// 首字母
	Initial string
}

