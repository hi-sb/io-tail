package filter

import "github.com/emicklei/go-restful"


// 注册filter 执行顺序
func init() {
	restful.Filter(globalLogging)
	restful.Filter(globalAuthTokenFilter)
	restful.Filter(globalAdminFilter)
}