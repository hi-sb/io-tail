package ext

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/auth"
	"github.com/hi-sb/io-tail/cache"
	"github.com/hi-sb/io-tail/syserr"
)

const (
	// user source
	privateSource = "private_source"
	// open source
	publicSource = "public_source"
)

// external_interface
// Interface definitions that allow external extensions.
type RedisExternalInterface struct {
}

//Check whether a user has write access to a resource, that is, whether messages can be sent to an identity
func (*RedisExternalInterface) CheckWritePermission(jwt *auth.JWT, name string) error {
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

// create open source
// If public key is null .
// The content will not be encrypted when it is stored, and no password is required when it is join
//
func (*RedisExternalInterface) CreateOpenSource(openSource *OpenSource) (bool, error) {
	if openSource.Name == "" {
		return false, syserr.NewParameterError("name is null")
	}
	publicSourceExists, _ := cache.RedisClient.HExists(publicSource, openSource.Name).Result()
	if publicSourceExists {
		return false, syserr.NewParameterError("name already exist")
	}
	bytes, err := json.Marshal(openSource)
	if err != nil {
		return false, syserr.NewSysErr(err.Error())
	}
	return cache.RedisClient.HSet(publicSource, openSource.Name, string(bytes)).Result()
}