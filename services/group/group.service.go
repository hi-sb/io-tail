package group

import (
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/model"
	"github.com/hi-sb/io-tail/utils"
	"github.com/jinzhu/gorm"
	"strings"
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
		return groupModelService.GetGroupInfoAndMembers(groupID,false)
	}()
	rest.WriteEntity(groupAndMemberInfo, err, response)
}

// 更新群公告
func (*GroupService) updateGroupNotice(request *restful.Request, response *restful.Response){
	err := func() error {
		groupModel := new(model.GroupModel)
		err := request.ReadEntity(groupModel)
		if err != nil {
			return err
		}

		if !(groupMemberModelService.CheckGroupRole(groupModel.ID,utils.Strval(utils.Strval(request.Attribute("currentUserId"))),false)){
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		err = mysql.DB.Model(groupModel).UpdateColumn("group_announcement",groupModel.GroupAnnouncement).Error
		if err != nil {
			return err
		}

		groupModelService.UpdateGroupInfoCache(groupModel.ID)

		return nil
	}()
	rest.WriteEntity(nil,err,response)
}

// 群禁言设置
func (*GroupService) updateGroupForbiddenStatus(request *restful.Request, response *restful.Response){
	err := func() error {
		groupModel := new(model.GroupModel)
		err := request.ReadEntity(groupModel)
		if err != nil {
			return err
		}

		if !(groupMemberModelService.CheckGroupRole(groupModel.ID,utils.Strval(utils.Strval(request.Attribute("currentUserId"))),false)){
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		// 验证状态有效性
		flag := groupModel.GroupChatStatus == 1 || groupModel.GroupChatStatus == 0

		if flag {
			err = mysql.DB.Model(groupModel).UpdateColumn("group_chat_status",groupModel.GroupChatStatus).Error
			if err != nil {
				return err
			}
		} else {
			return syserr.NewParameterError("参数有误")
		}
		groupModelService.UpdateGroupInfoCache(groupModel.ID)
		return nil
	}()
	rest.WriteEntity(nil,err,response)
}


// 解散群 并删除群成员
func (*GroupService) delGroupById(request *restful.Request, response *restful.Response){
	err := func() error {
		groupId := request.PathParameter("groupID")
		if groupId == ""{
			return syserr.NewParameterError("请求参数不能为空")
		}
		// 验证档期用户是否是群主
		if !(groupMemberModelService.CheckGroupRole(groupId,utils.Strval(utils.Strval(request.Attribute("currentUserId"))),true)){
			return syserr.NewPermissionErr("对不起，您没有权限操作")
		}

		/**
		  解散群成员 删除群信息  清除缓存
		 */
		return groupMemberModelService.DissolutionGroupAndClearCache(groupId)
	}()
	rest.WriteEntity(nil,err,response)
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/group")
	webService.Route(webService.POST("").To(groupService.createGroup))
	webService.Route(webService.GET("/{groupID}").To(groupService.findOne))
	webService.Route(webService.PUT("/global/notice").To(groupService.updateGroupNotice))
	webService.Route(webService.PUT("/global/forbidden/words").To(groupService.updateGroupForbiddenStatus))
	webService.Route(webService.DELETE("{groupID}").To(groupService.delGroupById))
	binder.BindAdd()
}
