package group

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/services/user"
	"github.com/hi-sb/io-tail/utils"
	"strconv"
	"strings"
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
			err = groupMemberService.checkMemberAndJoin(joinModel.GroupID,m)
			if err == nil {
				memberList.PushFront(m)
			}
		}

		//  加入成功  返回邀请者信息 被邀请者信息  当前群的基本信息 人数
		res := new(NewMemberJoinResModel)
		res.CurrentUser = userService.GetInfoById(userId)

		// 查询被邀请者
		var invitationUsers []user.UserModel
		for i := memberList.Front(); i != nil; i = i.Next() {
			user := userService.GetInfoById(utils.Strval(i.Value))
			if  user != nil {
				user.Password = ""
				invitationUsers = append(invitationUsers, *user)
			}
		}
		res.InvitationUserArray = &invitationUsers

		res.GroupInfo = groupModel
		res.Count = this.findMemberCountByGroupID(joinModel.GroupID)

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

// 查询当前邀请者是否已经加入 没有加入则持久化
func (g *GroupMemberService) checkMemberAndJoin(groupID string, userID string) error {
	err := func() error {
		gmd := new(GroupMemberModel)
		//  cache 如果当前用户在缓存中 证明已经在群中
		data,err := cache.RedisClient.HGet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),userID).Result()
		if err == nil && data != "" {
			return nil
		}

		//  反之，不在群中，加入缓存并持久化
		gmd.GroupID = groupID
		gmd.GroupMermerID = userID
		gmd.GroupMemberRole = 0
		gmd.Bind()
		err = g.insertMembers(gmd)
		if err != nil {
			return syserr.NewServiceError("加入群聊失败")
		}

		userInfo := userService.GetInfoById(gmd.GroupMermerID)
		if userInfo != nil{
			gmd.MobileNumber = userInfo.MobileNumber
			gmd.Avatar = userInfo.Avatar
			gmd.NickName = userInfo.NickName
			data,err := json.Marshal(gmd)
			if err == nil {
				cache.RedisClient.HSet(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID),userInfo.ID, data)
			}
		}
		return nil
	}()
	return err
}

// 查询当前群的人数
func (g *GroupMemberService) findMemberCountByGroupID(groupID string) int{
	// 返回当前群的人数
	memberArray,err := cache.RedisClient.HKeys(fmt.Sprintf(GROUP_MEMBER_INFO_REDIS_PREFIX,groupID)).Result()
	if err != nil {
		// 从DB 统计
		var total int = 0
		err = mysql.DB.Model(&GroupMemberModel{}).Where("group_id = ?",groupID).Count(&total).Error
		if err != nil {
			return 0
		}
		return total
	}
	return len(memberArray)
}


func init(){
	binder, webService := rest.NewJsonWebServiceBinder("/group-member")
	webService.Route(webService.POST("/join").To(groupMemberService.newMemberJoin))
	binder.BindAdd()
}
