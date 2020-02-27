package group

import (
	"fmt"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/syserr"
	"strings"
)

// 群聊设置
type GroupModel struct {
	db.BaseModel

	// 群名称
	GroupName string

	// 群公告
	GroupAnnouncement string

	// 群主
	GreateUserID string `gorm:"type:varchar(32);not null"`

	// 群聊天状态 0:全体禁言  1:正常
	GroupChatStatus int `gorm:"type:int(2);not null;default:1"`

}

// 创建群模型
type CreateGroupModel struct {
	// 群名称
	GroupName string
	// 群公告
	GroupAnnouncement string
	// 群成员
	GroupMembers string
}

// 验证创建群模型参数
func (g *CreateGroupModel) checkParams() error {

	if len(g.GroupMembers) <= 1 {
		return syserr.NewParameterError("参数有误，不能创建群聊")
	}
	if g.GroupName == "" {
		g.GroupName = fmt.Sprintf( "群聊(%d)",len(strings.Split(g.GroupMembers, ","))+1)
		return nil
	}
	return nil
}

type GroupInfoAndMembersModel struct {
	GroupModel GroupModel
	GroupMemberDetail *[]GroupMemberModel
}

// 获取成员和成员基础信息
func (g *GroupModel) GetGroupInfoAndMembers(groupID string) (*GroupInfoAndMembersModel,error) {
	groupAndMemberInfo,err := func() (*GroupInfoAndMembersModel,error){
		// 群基础信息
		groupModel := new(GroupModel)
		err := mysql.DB.Where("id = ?",groupID).First(groupModel).Error
		if err !=nil {
			return nil,err
		}
		// 群成员list
		gmList,err := new(GroupMemberModel).GetMembersInfo(groupID)
		groupAndMemberInfo := new(GroupInfoAndMembersModel)
		groupAndMemberInfo.GroupModel = *groupModel
		groupAndMemberInfo.GroupMemberDetail = gmList
		return groupAndMemberInfo,nil
	}()
	return groupAndMemberInfo,err
}
