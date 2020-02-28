package filter

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"strings"
)

const CURRENT_USER  =  "currentUserId"

// 全局登录验证
func globalAuthTokenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// 忽略验证的url集合
	urlMap := map[string]string {
		"_verify_sms":"/verify/sms",
		"_user_login":"/user/login",
	}
	// 当前请求的URI
	uri := strings.Replace(fmt.Sprintf("%s", req.Request.URL),"/","_",-1)

	if urlMap[uri] == "" {
		err := checkLogin(req)
		if err != nil {
			rest.WriteEntity(nil,err,resp)
			return
		}
	}
	chain.ProcessFilter(req, resp)
}


// 验证登录
func checkLogin(req *restful.Request) error{
	token := req.HeaderParameter(auth.AUTH_HEADER)
	userId, err := auth.GetUID(token)
	if userId == "" || err != nil {
		req.SetAttribute(CURRENT_USER,nil)
		return syserr.NewPermissionErr("您的签名已过期，请重新登录")
	}else{
		req.SetAttribute(CURRENT_USER,userId)
	}
	return err
}


