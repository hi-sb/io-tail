package group

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/utils"
	"strconv"
	"strings"
)

type GroupMemberService struct {
}

var groupMemberService = new(GroupMemberService)
var groupMemberModelService = new(GroupMemberModel)

// 邀请新用户加入群
func (this *GroupMemberService) newMemberJoin(request *restful.Request, response *restful.Response) {
	memberJoinResModel, err := func() (*NewMemberJoinResModel, error) {
		// 读取body
		joinModel := new(NewMemberJoinModel)
		err := request.ReadEntity(joinModel)
		if err != nil {
			return nil, err
		}

		// 查询群是否存在
		groupModel, err := groupModelService.GetGroupInfo(joinModel.GroupID)
		if err != nil {
			return nil, err
		}

		// 邀请的成员可能是多个
		members := strings.Split(joinModel.UserID, ",")
		memberList := list.New()
		for _, m := range members {
			// 查询当前邀请者是否已经加入 没有加入则持久化
			err = groupMemberModelService.checkMemberAndJoin(joinModel.GroupID, m)
			if err == nil {
				memberList.PushFront(m)
			}
		}

		//  加入成功  返回邀请者信息 被邀请者信息  当前群的基本信息 人数
		res := new(NewMemberJoinResModel)
		res.CurrentUser = userService.GetInfoById(utils.Strval(request.Attribute("currentUserId")))

		// 查询被邀请者
		var invitationUsers []GroupMemberModel
		for i := memberList.Front(); i != nil; i = i.Next() {
			gmd := new(GroupMemberModel)
			gmd.ID = groupModel.ID
			gmd.GroupID = joinModel.GroupID
			gmd.GroupMemberID = utils.Strval(i.Value)
			gmd.GroupMemberRole = 0
			user := userService.GetInfoById(utils.Strval(i.Value))
			if user != nil {
				gmd.MobileNumber = user.MobileNumber
				gmd.Avatar = user.Avatar
				gmd.NickName = user.NickName
				invitationUsers = append(invitationUsers, *gmd)
			}
		}
		res.InvitationUserArray = &invitationUsers

		res.GroupInfo = groupModel
		res.Count = groupMemberModelService.findMemberCountByGroupID(joinModel.GroupID)

		var buffer bytes.Buffer
		buffer.WriteString(res.GroupInfo.GroupName)
		buffer.WriteString("(")
		buffer.WriteString(strconv.Itoa(res.Count))
		buffer.WriteString(")")
		res.GroupInfo.GroupName = buffer.String()

		return res, nil
	}()
	rest.WriteEntity(memberJoinResModel, err, response)
}

// 群主或者管理员  从当前群组移除成员
func (*GroupMemberService) removeMember(request *restful.Request, response *restful.Response) {
	err := func() error {
		// 验证登录
		token := request.HeaderParameter(auth.AUTH_HEADER)
		userId, err := auth.GetUID(token)
		if userId == "" || err != nil {
			return errors.New("您还没有登录")
		}

		// 读取body
		rmModel := new(NewMemberJoinModel)
		err = request.ReadEntity(rmModel)
		if err != nil {
			return  err
		}

		// 验证当前用户在当前聊天组的角色
		currentUserGroupMember,err := groupMemberModelService.getGroupMemberByGroupIdAndMemberId(rmModel.GroupID,userId)
		if err != nil {
			return  err
		}
		// 普通成员无法剔除群成员
		if currentUserGroupMember.GroupMemberRole == 0 {
			return syserr.NewServiceError("你没有权限移除群成员")
		}
		// 删除群成员  redis db
		cache.RedisClient.HDel(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,rmModel.GroupID),rmModel.UserID)
		mysql.DB.Where("group_id = ? and group_member_id = ?",rmModel.GroupID,rmModel.UserID).Delete(GroupMemberModel{})
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 设置管理员
func (*GroupMemberService) setGroupAdmin(request *restful.Request, response *restful.Response){
	err := func() error {
		groupModelParams := new(GroupMemberModel)
		err := request.ReadEntity(groupModelParams)
		if err != nil {
			return err
		}
		// 验证群主\
		userId := utils.Strval(request.Attribute("currentUserId"))
		old,err := groupMemberModelService.getGroupMemberByGroupIdAndMemberId(groupModelParams.GroupID,userId)
		if err != nil {
			return err
		}
		if old.GroupMemberRole != 1 {
			return syserr.NewPermissionErr("对不起你么有权限设置管理员")
		}

		// 设置角色
		err = mysql.DB.Model(groupModelParams).Where("group_id = ? And group_member_id = ?",groupModelParams.GroupID,groupModelParams.GroupMemberID).UpdateColumn("group_member_role",groupModelParams.GroupMemberRole).Error
		if err != nil {
			return nil
		}
		// 刷新缓存
		groupMemberModelService.refushCacheGroupMemberInfo(groupModelParams.GroupID,groupModelParams.GroupMemberID)
		return nil
	}()
	rest.WriteEntity(nil,err,response)
}

// 群管理员设置成员昵称
func (*GroupMemberService) setMemberNickName(request *restful.Request, response *restful.Response){
	err := func() error {
		groupModelParams := new(GroupMemberModel)
		err := request.ReadEntity(groupModelParams)
		if err != nil {
			return err
		}
		// 设置昵称
		err = mysql.DB.Model(groupModelParams).Where("group_id = ? And group_member_id = ?",groupModelParams.GroupID,groupModelParams.GroupMemberID).UpdateColumn("group_member_nick_name",groupModelParams.GroupMemberNickName).Error
		if err != nil {
			return nil
		}
		// 刷新缓存
		groupMemberModelService.refushCacheGroupMemberInfo(groupModelParams.GroupID,groupModelParams.GroupMemberID)
		return nil
	}()
	rest.WriteEntity(nil,err,response)
}

// 推出群聊  普通成员不能退出群聊
func (*GroupMemberService) signOutGroupChat(request *restful.Request, response *restful.Response){
	err := func() error {
		groupID := request.PathParameter("groupID")
		memberID := utils.Strval(request.Attribute("currentUserId"))
		// 删除DB
		err := mysql.DB.Where("group_id =? and group_member_id = ?",groupID,memberID).Delete(&GroupMemberModel{}).Error
		if err != nil {
			return err
		}
		// 删除Redis
		cache.RedisClient.HDel(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX, groupID),memberID)
		return nil
	}()
	rest.WriteEntity(nil,err,response)
}


// 对某个成员设置禁言






func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/group-member")
	webService.Route(webService.POST("/join").To(groupMemberService.newMemberJoin))
	webService.Route(webService.DELETE("/remove").To(groupMemberService.removeMember))
	webService.Route(webService.PUT("/admin").To(groupMemberService.setGroupAdmin))
	webService.Route(webService.PUT("/nick-name").To(groupMemberService.setMemberNickName ))
	webService.Route(webService.DELETE("/{groupID}/sign-out").To(groupMemberService.signOutGroupChat ))
	binder.BindAdd()
}
