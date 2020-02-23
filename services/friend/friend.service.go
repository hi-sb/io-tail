package friend

type FriendService struct {

}

//地址
var friendService = new(FriendService)

// 添加好友

// 获取好友请求 Top10

// 拉黑好友

// 同意添加好友

// 获取好友列表

// 删除好友


//func init() {
//	binder, webService := rest.NewJsonWebServiceBinder("/friend")
//	//webService.Route(webService.GET("/{token}").To(userService.get))
//	//webService.Route(webService.POST("/login").To(userService.regOrlogin))
//	binder.BindAdd()
//}