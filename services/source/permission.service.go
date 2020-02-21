package source

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/syserr"
)

//
type PermissionService struct {

}


//检查写入权限
//对于个人消息和群消息来说，其实就是一个是私有资源
//而另一个是开放资源，写入一个私有资源，也就是给一个固定的人发送消息，那么需要验证发送者是否是被发送者的好友
//而开放资源则是写入一个群消息到话题，那么这个时候就需要验证该发送者是否加入了该群
func (*PermissionService) CheckWritePermission(jwt *auth.JWT, name string) error {
	//todo
	// open source and private source
	// check fd and contain
	//By default, all have send permission
	source, _ := cache.RedisClient.HGet(publicSource, name).Result()
	if source == "" {
		return nil
	}
	sourceModel := new(OpenSource)
	// public open source check source type
	err := json.Unmarshal([]byte(source), sourceModel)
	if err != nil {
		return syserr.NewSysErr(err.Error())
	}
	return syserr.NewPermissionErr("no permission send message")
}