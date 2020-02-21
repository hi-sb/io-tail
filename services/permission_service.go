package service

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/syserr"
)

//
type PermissionService struct {

}


//Check whether a user has write access to a resource, that is, whether messages can be sent to an identity
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