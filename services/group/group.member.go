package group

import (
	"encoding/json"
	"fmt"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/services/user"
)

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

	// 手机号
	MobileNumber string `gorm:"-"`
	//昵称
	NickName string  `gorm:"-"`
	// 头像
	Avatar string `gorm:"-"`

}

var userService = new(user.UserService)

// 获取成员和成员基础信息
func (g *GroupMemberModel) GetMembersInfo(groupID string) (*[]GroupMemberModel,error) {
	groupMemberDetails,err := func() (*[]GroupMemberModel,error){
		var groupMemberModels  []GroupMemberModel
		err := mysql.DB.Where("group_id = ?",groupID).Find(&groupMemberModels).Error
		if err !=nil {
			return nil,err
		}
		var groupMemberDetails []GroupMemberModel
		for _,memberModel := range groupMemberModels{
			var gmd GroupMemberModel
			gmd.GroupID = groupID
			userInfo := userService.GetInfoById(memberModel.GroupMermerID)
			if userInfo != nil{
				gmd.CreatedAt = memberModel.CreatedAt
				gmd.UpdatedAt = memberModel.UpdatedAt
				gmd.GroupMermerID = memberModel.GroupMermerID
				gmd.ID = memberModel.ID
				gmd.GroupMermerNickName = memberModel.NickName
				gmd.MobileNumber = userInfo.MobileNumber
				gmd.Avatar = userInfo.Avatar
				gmd.NickName = userInfo.NickName
				groupMemberDetails = append(groupMemberDetails, gmd)
				data,err := json.Marshal(gmd)
				if err == nil {
					cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),userInfo.ID, data)
				}
			}
		}
		return &groupMemberDetails,nil
	}()
	return groupMemberDetails,err
}



// 新用户加入群聊模型
type NewMemberJoinModel struct {
	GroupID string
	UserID string
}

// 新用户加入群聊返回模型
type  NewMemberJoinResModel struct {
	// 当前用户
	CurrentUser *user.UserModel
	// 被邀请者
	InvitationUser *user.UserModel
	// 群基础信息
	GroupInfo *GroupModel
	// 群人数
	Count int
}