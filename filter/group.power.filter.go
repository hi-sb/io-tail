package filter

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/services/group"
	"io/ioutil"
	"strings"
)

// 群管理员 群主权限验证
func groupPowerFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	// 需要验证的url集合
	urlMap := map[string]string {
		"_group_global_forbidden_words":"/group/global/forbidden/words",
		"_group_global_notice":"/group/global/notice",
		"_group-member_remove":"/group-member/remove",
		"_group-member_join":"/group-member/join",
		"_group-member_admin ":"/group-member/admin",
	}
	// 当前请求的URI
	uri := strings.Replace(fmt.Sprintf("%s", req.Request.URL),"/","_",-1)
	if urlMap[uri] != "" {
		err := checkGroupPower(req)
		if err != nil {
			rest.WriteEntity(nil,err,resp)
			return
		}
	}
	chain.ProcessFilter(req, resp)
}

// 验证是否是群主 或者 管理员
func checkGroupPower(req *restful.Request) error {
	userId := req.Attribute(CURRENT_USER)
	bodyBytes, _ := ioutil.ReadAll(req.Request.Body)
	req.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 读取参数
	groupModel := new(group.GroupModel)
	err := req.ReadEntity(groupModel)
	if err != nil {
		return err
	}
	if groupModel.ID == ""  {
		return syserr.NewParameterError("参数缺失")
	}

	groupMemberModel := new(group.GroupMemberModel)
	err = mysql.DB.Where("group_id = ? and group_member_id = ?",groupModel.ID,userId).Find(groupMemberModel).Error
	if err !=nil {
		return err
	}
	if groupMemberModel.GroupMemberRole == 0 {
		return syserr.NewPermissionErr("对不起，你没有权限操作")
	}

	req.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}