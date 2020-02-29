package group

import (
	"encoding/json"
	"fmt"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/log"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/services/user"
)

// 群成员
type GroupMemberModel struct {
	db.BaseModel
	// 群ID
	GroupID string `gorm:"type:varchar(32);not null"`
	// 成员
	GroupMemberID string `gorm:"type:varchar(32);not null"`
	// 成员在当前群的昵称
	GroupMemberNickName string `gorm:"type:varchar(255)"`
	// 成员角色  0: 普通成员 1.群主  2。管理员
	GroupMemberRole int `gorm:"type:int(2);not null;default:0"`
	// 是否被禁言  0: 正常发言 1:禁言
	IsForbidden int `gorm:"type:int(2);not null;default:0"`
	// 手机号
	MobileNumber string `gorm:"-"`
	//昵称
	NickName string `gorm:"-"`
	// 头像
	Avatar string `gorm:"-"`
}

var userService = new(user.UserService)

// 新用户加入群聊模型
type NewMemberJoinModel struct {
	GroupID string
	UserID  string
}

// 新用户加入群聊返回模型
type NewMemberJoinResModel struct {
	// 当前用户
	CurrentUser *user.UserModel
	// 被邀请者
	InvitationUserArray *[]GroupMemberModel
	// 群基础信息
	GroupInfo *GroupModel
	// 群人数
	Count int
}

// 持久化群成员到DB
func (*GroupMemberModel) insertMembers(model *GroupMemberModel) error {
	err := mysql.DB.Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取成员和成员基础信息
func (g *GroupMemberModel) GetMembersInfo(groupID string, isNewGroup bool) (*[]GroupMemberModel, error) {
	groupMemberDetails, err := func() (*[]GroupMemberModel, error) {
		// 新建的群 从库中读取 并加入缓存
		if isNewGroup {
			return g.getGroupMemberDetailsForDB(groupID)
		} else {
			// 已经创建好的群  从reids读取成员信息
			dataMap, err := cache.RedisClient.HGetAll(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID)).Result()
			if err != nil || dataMap == nil {
				// 读库
				return g.getGroupMemberDetailsForDB(groupID)
			}else {
				var groupMemberModels []GroupMemberModel
				for _, v := range dataMap {
					gmd := new(GroupMemberModel)
					err := json.Unmarshal([]byte(v), gmd)
					if err != nil {
						log.Log.Error(err)
					}
					groupMemberModels = append(groupMemberModels, *gmd)
				}
				// 如果缓存读取失败 读取DB
				if len(groupMemberModels) == 0 {
					return g.getGroupMemberDetailsForDB(groupID)
				}
				return &groupMemberModels,nil
			}
		}
	}()
	return groupMemberDetails, err
}

// 查询当前群的成员人数
func (g *GroupMemberModel) findMemberCountByGroupID(groupID string) int {
	// 返回当前群的人数
	memberArray, err := cache.RedisClient.HKeys(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID)).Result()
	if err != nil {
		// 从DB 统计 数量
		var total int = 0
		err = mysql.DB.Model(&GroupMemberModel{}).Where("group_id = ?", groupID).Count(&total).Error
		if err != nil {
			return 0
		}
		return total
	}
	return len(memberArray)
}

// 查询当前邀请者是否已经加入 没有加入则持久化
func (g *GroupMemberModel) checkMemberAndJoin(groupID string, userID string) error {
	err := func() error {
		gmd := new(GroupMemberModel)
		//  cache 如果当前用户在缓存中 证明已经在群中
		data, err := cache.RedisClient.HGet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID), userID).Result()
		if err == nil && data != "" {
			return syserr.NewServiceError("当前成员已经添加")
		}
		//  反之，不在群中，加入缓存并持久化
		gmd.GroupID = groupID
		gmd.GroupMemberID = userID
		gmd.GroupMemberRole = 0
		gmd.Bind()
		err = g.insertMembers(gmd)
		if err != nil {
			return syserr.NewServiceError("加入群聊失败")
		}

		userInfo := userService.GetInfoById(gmd.GroupMemberID)
		if userInfo != nil {
			gmd.MobileNumber = userInfo.MobileNumber
			gmd.Avatar = userInfo.Avatar
			gmd.NickName = userInfo.NickName
			data, err := json.Marshal(gmd)
			if err == nil {
				cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID), userInfo.ID, data)
			}
		}
		return nil
	}()
	return err
}

// 从数据库获取群成员信息 并加入缓存
func (g *GroupMemberModel) getGroupMemberDetailsForDB(groupID string) (*[]GroupMemberModel, error) {
	groupMemberDetails, err := func() (*[]GroupMemberModel, error) {
		var groupMemberModels []GroupMemberModel
		err := mysql.DB.Where("group_id = ?", groupID).Find(&groupMemberModels).Error
		if err != nil {
			return nil, err
		}
		var groupMemberDetails []GroupMemberModel
		for _, memberModel := range groupMemberModels {
			userInfo := userService.GetInfoById(memberModel.GroupMemberID)
			if userInfo != nil {
				memberModel.MobileNumber = userInfo.MobileNumber
				memberModel.Avatar = userInfo.Avatar
				memberModel.NickName = userInfo.NickName
				groupMemberDetails = append(groupMemberDetails, memberModel)
				data, err := json.Marshal(memberModel)
				if err == nil {
					cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID), userInfo.ID, data)
				}
			}
		}
		return &groupMemberDetails, nil
	}()
	return groupMemberDetails, err
}

// 根据 userID groupID 从db查询并刷新groupMemberInfo 缓存到redis
func (*GroupMemberModel) refushCacheGroupMemberInfo(groupID string,memberID string){
	groupMemberModel:= new( GroupMemberModel)
	err := mysql.DB.Where("group_id = ? and group_member_id=?", groupID,memberID).Find(groupMemberModel).Error
	if err != nil {
		log.Log.Error(err)
	}
	userInfo := userService.GetInfoById(memberID)
	if userInfo != nil {
		groupMemberModel.MobileNumber = userInfo.MobileNumber
		groupMemberModel.Avatar = userInfo.Avatar
		groupMemberModel.NickName = userInfo.NickName
		data, err := json.Marshal(groupMemberModel)
		if err == nil {
			cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID), userInfo.ID, data)
		}else {
			log.Log.Error(err)
		}
	}
}

// 根据 userID groupID 获取 groupMemberInfo  用户群角色验证
func (g *GroupMemberModel) getGroupMemberByGroupIdAndMemberId(groupID string,memberID string) (*GroupMemberModel, error){
	groupMemberInfo,err := func() (*GroupMemberModel, error){
		gmd := new(GroupMemberModel)
		jsonData, err := cache.RedisClient.HGet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),memberID).Result()
		err = json.Unmarshal([]byte(jsonData), gmd)
		if err != nil {
			// 查DB
			err = mysql.DB.Where("group_id = ? and group_member_id = ?",groupID,memberID).Find(gmd).Error
			if err != nil {
				return nil,err
			}else {
				return gmd,nil
			}
		}else{
			return gmd,nil
		}
	}()
	return groupMemberInfo,err
}


