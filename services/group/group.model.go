package group

import (
	"encoding/json"
	"fmt"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/log"
	"github.com/hi-sb/io-tail/core/syserr"
	"strconv"
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
		g.GroupName = "群聊"
		return nil
	}
	return nil
}


// 群信息和成员列表
type GroupInfoAndMembersModel struct {
	GroupModel GroupModel  // 群基础信息
	GroupMemberDetail *[]GroupMemberModel  // 群成员信息
}


/**
    isNewGroup: true 新建   false 缓存读取
	获取成员和成员基础信息
	如果是新创建的群  缓存到 redis  反之从redis读取
 */
func (g *GroupModel) GetGroupInfoAndMembers(groupID string,isNewGroup bool) (*GroupInfoAndMembersModel,error) {
	groupAndMemberInfo,err := func() (*GroupInfoAndMembersModel,error){
		// 群基础信息
		groupModel := new(GroupModel)
		// 如果是新建的群  缓存到redis
		if isNewGroup {
			err := mysql.DB.Where("id = ?",groupID).First(groupModel).Error
			if err !=nil {
				return nil,err
			}

			data,err := json.Marshal(groupModel)
			if err == nil {
				_,err = cache.RedisClient.Set(fmt.Sprintf(GROUP_BASE_INFO_REDIS_PREFIX,groupID),data,0).Result()
				if err !=nil {
					log.Log.Error(err)
				}
			}
		}else{  // 从缓存读取groupInfo
			jsonData,err := cache.RedisClient.Get(fmt.Sprintf(GROUP_BASE_INFO_REDIS_PREFIX,groupID)).Result()
			if err != nil || jsonData == "" {
				//
				err := mysql.DB.Where("id = ?",groupID).First(groupModel).Error
				if err !=nil {
					return nil,err
				}
			}else {
				err := json.Unmarshal([]byte(jsonData), groupModel)
				if err != nil {
					log.Log.Error(err)
					return nil,err
				}
			}
		}

		// 群成员list
		gmList, _ := groupMemberModelService.GetMembersInfo(groupID,isNewGroup)
		groupAndMemberInfo := new(GroupInfoAndMembersModel)
		groupAndMemberInfo.GroupModel = *groupModel
		groupAndMemberInfo.GroupMemberDetail = gmList
		num := len(*groupAndMemberInfo.GroupMemberDetail)
		groupAndMemberInfo.GroupModel.GroupName = groupAndMemberInfo.GroupModel.GroupName + "(" +strconv.Itoa(num) +")"
		return groupAndMemberInfo,nil
	}()
	return groupAndMemberInfo,err
}


// 获取群基本信息
func (g *GroupModel) GetGroupInfo(groupID string) (*GroupModel,error) {
	// 群基础信息
	groupModel := new(GroupModel)

	// read redis
	data,err :=cache.RedisClient.Get(fmt.Sprintf(GROUP_BASE_INFO_REDIS_PREFIX,groupID)).Result()
	if err == nil && data != "" {
		err := json.Unmarshal([]byte(data), groupModel)
		if err != nil {
			fmt.Println(err)
		}
		return groupModel,nil
	}
	// read DB
	err = mysql.DB.Where("id = ?",groupID).First(groupModel).Error
	if err !=nil {
		return nil,err
	}
	return groupModel,nil
}

// 更新群基本信息缓存
func (*GroupModel) updateGroupInfoCache(groupID string){
	// 群基础信息
	groupModel := new(GroupModel)
	// read DB
	err := mysql.DB.Where("id = ?",groupID).First(groupModel).Error
	if err !=nil {
		log.Log.Error(err)
	}

	data,err := json.Marshal(groupModel)
	if err == nil {
		_,err = cache.RedisClient.Set(fmt.Sprintf(GROUP_BASE_INFO_REDIS_PREFIX,groupID),data,0).Result()
		if err !=nil {
			log.Log.Error(err)
		}
	}



}