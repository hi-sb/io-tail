package group

import (
	"encoding/json"
	"fmt"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/syserr"
)

type GroupMemberService struct {
}

var groupMemberService = new(GroupMemberService)


// 插入群成员
func (*GroupMemberService) insertMembers(model *GroupMemberModel) error {
	err := mysql.DB.Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询当前邀请者是否已经加入 没有加入则持久化
func (g *GroupMemberService) checkMemberAndJoin(groupID string, userID string) error {
	err := func() error {
		gmd := new(GroupMemberModel)
		//  cache
		data,err := cache.RedisClient.HGet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),userID).Result()
		if err == nil && data != "" {
			return nil
		}

		//  持久化
		gmd.GroupID = groupID
		gmd.GroupMermerID = userID
		gmd.GroupMemberRole = 0
		gmd.Bind()
		err = g.insertMembers(gmd)
		if err != nil {
			return syserr.NewServiceError("加入群聊失败")
		}

		userInfo := userService.GetInfoById(memberModel.GroupMermerID)
		if userInfo != nil{
			gmd.MobileNumber = userInfo.MobileNumber
			gmd.Avatar = userInfo.Avatar
			gmd.NickName = userInfo.NickName
			data,err := json.Marshal(gmd)
			if err == nil {
				cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),userInfo.ID, data)
			}
		}
		return err
	}()
	return err
}