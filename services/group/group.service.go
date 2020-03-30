package group

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/body"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/model"
	"github.com/hi-sb/io-tail/utils"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type GroupService struct {
}

var groupService = new(GroupService)
var groupModelService = new(model.GroupModel)

//  创建群
func (*GroupService) createGroup(request *restful.Request, response *restful.Response) {
	groupInfoAndMembers, err := func() (*model.GroupInfoAndMembersModel, error) {
		userId := utils.Strval(request.Attribute("currentUserId"))

		// 读取参数
		createGroup := new(model.CreateGroupModel)
		err := request.ReadEntity(createGroup)
		if err != nil {
			return nil, err
		}
		// 验证参数
		err = createGroup.CheckParams()
		if err != nil {
			return nil, err
		}

		mems := strings.Split(createGroup.GroupMembers, ",")
		if len(mems) <= 2 {
			return nil,syserr.NewParameterError("群成员必须大于2个人")
		}

		/**
		创建群流程开始  1.创建主群信息  2.写入成员
		*/
		// 主群模型
		groupModel := new(model.GroupModel)
		groupModel.GreateUserID = userId
		groupModel.GroupName = createGroup.GroupName
		groupModel.GroupAnnouncement = createGroup.GroupAnnouncement
		groupModel.Bind()

		// 事务处理
		err = mysql.Transactional(func(tx *gorm.DB) error {
			members := strings.Split(createGroup.GroupMembers, ",")
			members = append(members, userId)
			// 持久化群成员
			for _, member := range members {
				groupMember := new(model.GroupMemberModel)
				groupMember.GroupID = groupModel.ID
				if member == userId {
					groupMember.GroupMemberID = member
					groupMember.GroupMemberRole = 1 // 设置群主
				}
				groupMember.GroupMemberID = member
				groupMember.Bind()
				err = tx.Create(groupMember).Error
				if err != nil {
					return err
				}
				addGroupStringByte, _ := json.Marshal(groupMember)
				addGroupSendRequest := model.SendRequest{
					SendTime:    time.Now().Unix(),
					Body:        string(addGroupStringByte),
					ContentType: body.MessageTypeAddToGroup,
				}
				//发送加入消息
				go message.SendMessage("-1", groupMember.GroupMemberID, &addGroupSendRequest)
			}

			// 持久化群
			err = tx.Create(groupModel).Error
			if err != nil {
				return err
			}

			return nil
		})

		// 返回初始化后的群信息
		return groupModelService.GetGroupInfoAndMembers(groupModel.ID, true)
	}()
	rest.WriteEntity(groupInfoAndMembers, err, response)
}

//  获取当前群信息 以及群成员
func (*GroupService) findOne(request *restful.Request, response *restful.Response) {
	groupAndMemberInfo, err := func() (*model.GroupInfoAndMembersModel, error) {
		// 读取body
		groupID := request.PathParameter("groupID")
		if groupID == "" {
			return nil, syserr.NewParameterError("参数不正确")
		}
		return groupModelService.GetGroupInfoAndMembers(groupID, false)
	}()
	rest.WriteEntity(groupAndMemberInfo, err, response)
}

// 更新群公告
func (*GroupService) updateGroupNotice(request *restful.Request, response *restful.Response) {
	err := func() error {
		groupModel := new(model.GroupModel)
		err := request.ReadEntity(groupModel)
		if err != nil {
			return err
		}

		if !(groupMemberModelService.CheckGroupRole(groupModel.ID, utils.Strval(utils.Strval(request.Attribute("currentUserId"))), false)) {
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		if groupModel.GroupAnnouncement != ""  {
			err = mysql.DB.Model(groupModel).UpdateColumn("group_announcement", groupModel.GroupAnnouncement).Error
		}

		if  groupModel.GroupName != ""{
			err = mysql.DB.Model(groupModel).UpdateColumn("group_name", groupModel.GroupName).Error
		}

		if err != nil {
			return err
		}


		groupModelService.UpdateGroupInfoCache(groupModel.ID)

		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 群禁言设置
func (*GroupService) updateGroupForbiddenStatus(request *restful.Request, response *restful.Response) {
	err := func() error {
		groupModel := new(model.GroupModel)
		err := request.ReadEntity(groupModel)
		if err != nil {
			return err
		}

		if !(groupMemberModelService.CheckGroupRole(groupModel.ID, utils.Strval(utils.Strval(request.Attribute("currentUserId"))), false)) {
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		// 验证状态有效性
		flag := groupModel.GroupChatStatus == 1 || groupModel.GroupChatStatus == 0

		if flag {
			err = mysql.DB.Model(groupModel).UpdateColumn("group_chat_status", groupModel.GroupChatStatus).Error
			if err != nil {
				return err
			}
		} else {
			return syserr.NewParameterError("参数有误")
		}
		groupModelService.UpdateGroupInfoCache(groupModel.ID)
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 解散群 并删除群成员
func (*GroupService) delGroupById(request *restful.Request, response *restful.Response) {
	err := func() error {
		groupId := request.PathParameter("groupID")
		if groupId == "" {
			return syserr.NewParameterError("请求参数不能为空")
		}
		// 验证档期用户是否是群主
		if !(groupMemberModelService.CheckGroupRole(groupId, utils.Strval(utils.Strval(request.Attribute("currentUserId"))), true)) {
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		/**
		  解散群成员 删除群信息  清除缓存
		*/
		return groupMemberModelService.DissolutionGroupAndClearCache(groupId)
	}()
	rest.WriteEntity(nil, err, response)
}

/***
 群消息验证
	1: 验证当前群组的生命状态
	2: 验证当前群组的会话状态
	3: 验证当前用户是否被禁言
  	返回 true: 正常对话  false: 不能说话
*/
func (*GroupService) checkDialogueStatus(request *restful.Request, response *restful.Response) {
	result, err := func() (bool, error) {
		groupId := request.PathParameter("groupID")
		if groupId == "" {
			return false, syserr.NewParameterError("参数有误")
		}

		// 验证当前群的生命状态
		groupInfo, err := groupModelService.GetGroupInfo(groupId)
		if err != nil {
			if err.Error() =="record not found" {
				return false, syserr.NewServiceError("对不起,当前群已经被解散")
			}
			return false, err
		}
		if err == nil && groupInfo == nil {
			return false, nil
		}
		// 验证当前群的会话状态
		if groupInfo.GroupChatStatus == 0 {
			return false, nil
		}

		// 验证当前用户在当前群组是否被禁言
		groupMemberInfo, err := groupMemberModelService.GetGroupMemberByGroupIdAndMemberId(groupId, utils.Strval(utils.Strval(request.Attribute("currentUserId"))))
		if err != nil {
			if err.Error() =="record not found" {
				return false, syserr.NewServiceError("对不起,你已经被管理员请出当前聊天群")
			}
			return false, err
		}

		if err == nil && groupMemberInfo == nil {
			return false, nil
		}

		if groupMemberInfo.IsForbidden == 0 {
			return true, nil
		}

		return false, syserr.NewServiceError("对不起，您已经被禁言")
	}()
	rest.WriteErrAndEntity(result, err, response)
}

// 根据UserId 获取已经加入的群
func (*GroupService) findGroupIdsByUserId(request *restful.Request, response *restful.Response) {
	groupArray,err := func()(*[]model.GroupModel,error) {
		return groupModelService.GetGroupsByUserId(utils.Strval(utils.Strval(request.Attribute("currentUserId"))))
	}()
	rest.WriteEntity(groupArray,err,response)
}


func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/group")
	webService.Route(webService.POST("").To(groupService.createGroup))
	webService.Route(webService.GET("/{groupID}").To(groupService.findOne))
	webService.Route(webService.PUT("").To(groupService.updateGroupNotice))
	webService.Route(webService.PUT("/global/forbidden/words").To(groupService.updateGroupForbiddenStatus))
	webService.Route(webService.DELETE("{groupID}").To(groupService.delGroupById))
	webService.Route(webService.GET("/check/{groupID}").To(groupService.checkDialogueStatus))
	webService.Route(webService.GET("/join/items").To(groupService.findGroupIdsByUserId))
	binder.BindAdd()
}
