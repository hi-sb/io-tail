package group

import (
	"bytes"
	"container/list"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/utils"
	"strconv"
	"strings"
)

type GroupMemberService struct {
}

var groupMemberService = new(GroupMemberService)


// 邀请新用户加入群
func (this *GroupMemberService) newMemberJoin(request *restful.Request, response *restful.Response){
	memberJoinResModel,err := func() (*NewMemberJoinResModel,error) {
		// 验证登录
		token := request.HeaderParameter(auth.AUTH_HEADER)
		userId, err := auth.GetUID(token)
		if userId == "" || err != nil {
			return nil,errors.New("您还没有登录")
		}

		// 读取body
		joinModel := new(NewMemberJoinModel)
		err = request.ReadEntity(joinModel)
		if err != nil {
			return nil,err
		}

		// 查询群是否存在
		groupModel,err := groupModelService.GetGroupInfo(joinModel.GroupID)
		if err != nil {
			return nil,err
		}

		// 邀请的成员可能是多个
		members := strings.Split(joinModel.UserID, ",")
		memberList := list.New()
		for _,m := range members{
			// 查询当前邀请者是否已经加入 没有加入则持久化
			err = groupMemberModelService.checkMemberAndJoin(joinModel.GroupID,m)
			if err == nil {
				memberList.PushFront(m)
			}
		}

		//  加入成功  返回邀请者信息 被邀请者信息  当前群的基本信息 人数
		res := new(NewMemberJoinResModel)
		res.CurrentUser = userService.GetInfoById(userId)

		// 查询被邀请者
		var invitationUsers []GroupMemberModel
		for i := memberList.Front(); i != nil; i = i.Next() {
			gmd := new(GroupMemberModel)
			gmd.ID = groupModel.ID
			gmd.GroupID = joinModel.GroupID
			gmd.GroupMermerID = utils.Strval(i.Value)
			gmd.GroupMemberRole = 0
			user := userService.GetInfoById(utils.Strval(i.Value))
			if  user != nil {
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

		return res,nil
	}()
	rest.WriteEntity(memberJoinResModel, err, response)
}

//查看当前群成员以及群信息
//func (*GroupMemberService) getMembers(request *restful.Request, response *restful.Response){
//	memberResModel,err := func() (*NewMemberJoinResModel,error) {
//		// 验证登录
//		token := request.HeaderParameter(auth.AUTH_HEADER)
//		userId, err := auth.GetUID(token)
//		if userId == "" || err != nil {
//			return nil,errors.New("您还没有登录")
//		}
//
//		// 读取body
//		groupID := request.PathParameter("groupID")
//		if groupID != "" {
//			return nil,syserr.NewParameterError("参数不正确")
//		}
//
//		// 查询群是否存在
//		groupModel,err := groupModelService.GetGroupInfo(groupID)
//		if err != nil {
//			return nil,err
//		}
//
//		//  当前群的的信息以及成员信息
//		res := new(NewMemberJoinResModel)
//
//		var memberArray []string
//		memberArray,err = cache.RedisClient.HKeys(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID)).Result()
//		if err != nil {
//			err = mysql.DB.Model(&GroupMemberModel{}).Where("group_id = ?",groupID).Count(&total).Error
//		}
//
//
//		// 查询群成员
//		var invitationUsers []user.UserModel
//		for i := memberList.Front(); i != nil; i = i.Next() {
//			user := userService.GetInfoById(utils.Strval(i.Value))
//			if  user != nil {
//				user.Password = ""
//				invitationUsers = append(invitationUsers, *user)
//			}
//		}
//		res.InvitationUserArray = &invitationUsers
//
//		res.GroupInfo = groupModel
//		res.Count = this.findMemberCountByGroupID(joinModel.GroupID)
//
//		var buffer bytes.Buffer
//		buffer.WriteString(res.GroupInfo.GroupName)
//		buffer.WriteString("(")
//		buffer.WriteString(strconv.Itoa(res.Count))
//		buffer.WriteString(")")
//		res.GroupInfo.GroupName = buffer.String()
//
//		return res,nil
//	}()
//	rest.WriteEntity(memberResModel, err, response)
//}











func init(){
	binder, webService := rest.NewJsonWebServiceBinder("/group-member")
	webService.Route(webService.POST("/join").To(groupMemberService.newMemberJoin))
//	webService.Route(webService.GET("/details/{groupID}").To(groupMemberService.newMemberJoin))
	binder.BindAdd()
}
