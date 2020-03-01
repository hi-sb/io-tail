package model

import (
	"fmt"
	"github.com/hi-sb/io-tail/common/constants"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/syserr"
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



// 设置拉黑值 并更新redis
func (*PullBlackModel) SetIsBlack(friendModel *FriendModel,status int,currentUserID string) *FriendModel {
	// 0 拉黑  1 正常
	if friendModel.UserID == currentUserID {  	// 如果当前用户是U 对F操作
		// 原始状态f拉黑u,u未拉黑f(10)  即将操作 u拉黑f  设置状态为互相拉黑(00)
		if friendModel.IsBlack == constants.IS_BLACK_F_PULL_U && status == 0 {
			friendModel.IsBlack = constants.IS_BLACK_EACH_OTHER
			// 将F加入U的黑名单
			cache.RedisClient.SAdd(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}
		// 原始状态 互相拉黑（00）   即将操作 u恢复对f的关系  设置状态为（10）
		if friendModel.IsBlack == constants.IS_BLACK_EACH_OTHER && status == 1 {
			friendModel.IsBlack = constants.IS_BLACK_F_PULL_U
			// 将F从U的黑名单的移除
			cache.RedisClient.SRem(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}

		// 原始状态 正常（11）   即将操作 u拉黑f    设置状态（01）
		if friendModel.IsBlack == constants.IS_NOT_BLACK && status == 0 {
			friendModel.IsBlack = constants.IS_BLACK_U_PULL_F
			// 将F加入U的黑名单
			cache.RedisClient.SAdd(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}
		// // 原始状态 （01）  即将操作 u未拉黑f  设置状态为正常（00）
		if friendModel.IsBlack == constants.IS_BLACK_U_PULL_F && status == 1 {
			friendModel.IsBlack = constants.IS_NOT_BLACK
			// 将F从U的黑名单的移除
			cache.RedisClient.SRem(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}

	} else if friendModel.FriendID == currentUserID {  // 如果当前用户是F 对U操作
		// 原始状态 f未拉黑u  u拉黑f （01）  即将操作  f拉黑u  设置状态为互相拉黑（00）
		if friendModel.IsBlack == constants.IS_BLACK_U_PULL_F && status == 0 {
			friendModel.IsBlack = constants.IS_BLACK_EACH_OTHER
			cache.RedisClient.SAdd(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}
		// 原始状态 f未拉黑u  u未拉黑f（11）  即将操作  f拉黑u  设置状态（10）
		if friendModel.IsBlack == constants.IS_NOT_BLACK && status == 0 {
			friendModel.IsBlack = constants.IS_BLACK_F_PULL_U
			cache.RedisClient.SAdd(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}

		// 原始状态 f拉黑u  u拉黑f（00）  即将操作  f未拉黑u  设置状态为 （01）
		if friendModel.IsBlack == constants.IS_BLACK_EACH_OTHER && status == 1 {
			friendModel.IsBlack = constants.IS_BLACK_U_PULL_F
			cache.RedisClient.SRem(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}

		// 原始状态 f拉黑u  u未拉黑f（10）  即将操作  f恢复拉黑u  设置状态为（11）
		if friendModel.IsBlack == constants.IS_BLACK_F_PULL_U && status == 1 {
			friendModel.IsBlack = constants.IS_NOT_BLACK
			cache.RedisClient.SRem(fmt.Sprintf(constants.FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}
	}

	return friendModel
}

// 首字母排序检查是否是a-z
func (*FriendAddReqModel) CheckAscII(ascValue int) int{
	if (ascValue>= 122 && ascValue <= 97) || (ascValue>= 65 && ascValue <= 90) {
		return ascValue
	}else{
		// # 的ASCII值
		return 35
	}
}

