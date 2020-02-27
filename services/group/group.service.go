package group

import (
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/jinzhu/gorm"
	"strings"
)

type GroupService struct {
}

var groupService = new(GroupService)


//  创建群
func (*GroupService) createGroup(request *restful.Request, response *restful.Response){
	groupInfoAndMembers,err := func() (*GroupInfoAndMembersModel,error){
		// 验证是否登录
		token := request.HeaderParameter(auth.AUTH_HEADER)
		userId, err := auth.GetUID(token)
		if userId == "" || err != nil {
			return nil,errors.New("您还没有登录")
		}

		// 读取参数
		createGroup := new(CreateGroupModel)
		err = request.ReadEntity(createGroup)
		if err != nil {
			return nil,err
		}
		// 验证参数
		err = createGroup.checkParams()
		if err != nil {
			return nil,err
		}

		/**
		创建群流程开始  1.创建主群信息  2.写入成员
		 */
		// 主群模型
		groupModel := new(GroupModel)
		groupModel.GreateUserID = userId
		groupModel.GroupName = createGroup.GroupName
		groupModel.GroupAnnouncement = createGroup.GroupAnnouncement
		groupModel.Bind()

		// 事务处理
		err = mysql.Transactional(func(tx *gorm.DB) error {
			members := strings.Split(createGroup.GroupMembers, ",")
			members = append(members, userId)
			// 持久化群成员
			for _,member := range members{
				groupMember := new(GroupMemberModel)
				groupMember.GroupID = groupModel.ID
				if member == userId {
					groupMember.GroupMermerID = member
					groupMember.GroupMemberRole = 1 // 设置群主
				}
				groupMember.GroupMermerID = member
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
		return new(GroupModel).GetGroupInfoAndMembers(groupModel.ID)
	}()
	rest.WriteEntity(groupInfoAndMembers, err, response)
}


func init(){
	binder, webService := rest.NewJsonWebServiceBinder("/group")
	webService.Route(webService.POST("").To(groupService.createGroup))
	binder.BindAdd()
}