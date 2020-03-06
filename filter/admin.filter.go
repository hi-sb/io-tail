package filter

import (
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/model"
	"strings"
)


var userModelService = new(model.UserModel)

// 全局登录验证
func globalAdminFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if strings.HasPrefix(req.Request.URL.String(),"/admin") {
		result,err := checkUserRole(req)
		if err != nil {
			rest.WriteEntity(nil,err,resp)
			return
		}

		if !result {
			rest.WriteEntity(nil,err,resp)
			return
		}

	}
	chain.ProcessFilter(req, resp)
}

// 验证角色
func checkUserRole(req *restful.Request) (bool,error){
	token := req.HeaderParameter(auth.AUTH_HEADER)
	userId, err := auth.GetUID(token)
	if userId == "" || err != nil {
		req.SetAttribute(CURRENT_USER,nil)
		return false,syserr.NewPermissionErr("您的签名已过期，请重新登录")
	}else{
		req.SetAttribute(CURRENT_USER,userId)
		// 验证角色
		userInfo := userModelService.GetInfoById(userId)
		if userInfo == nil {
			return false,syserr.NewServiceError("对不起，您无权限访问")
		}

		if userInfo.UserRole == 1 {
			return true,nil
		}

		return false,syserr.NewServiceError("对不起，您无权限访问")
	}
}
