package ext

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/auth"
	"github.com/hi-sb/io-tail/cache"
	"github.com/hi-sb/io-tail/syserr"
	"strings"
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

//Here we only need an interface to get the public key by user name. Therefore your username must be unique.
func (*RedisExternalInterface) GetUserPublicKey(name string) (string, error) {
	id := strings.Split(name, "@")
	source, _ := cache.RedisClient.HGet(privateSource, id[0]).Result()
	if source == "" {
		return source, nil
	}
	sourceModel := new(Source)
	err := json.Unmarshal([]byte(source), sourceModel)
	if err != nil {
		return "", syserr.NewSysErr(err.Error())
	}
	return sourceModel.PublicKey, nil
}

//Check whether a user has write access to a resource, that is, whether messages can be sent to an identity
func (*RedisExternalInterface) CheckWritePermission(jwt *auth.JWT, name string) error {
	id := strings.Split(jwt.AtNum, "@")
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
	//Allow everyone to send messages
	if sourceModel.SourceType == InteractiveOpenSourceType {
		return nil
	}
	// Only creators are allowed to send messages
	if sourceModel.SourceType == SubscriptionOpenSourceType &&
		id[0] == name {
		return nil
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
	if openSource.SourceType != SubscriptionOpenSourceType &&
		openSource.SourceType != InteractiveOpenSourceType {
		return false, syserr.NewParameterError("SourceType err ")
	}
	privateSourceExists, _ := cache.RedisClient.HExists(privateSource, openSource.Name).Result()
	publicSourceExists, _ := cache.RedisClient.HExists(publicSource, openSource.Name).Result()
	if privateSourceExists || publicSourceExists {
		return false, syserr.NewParameterError("name already exist")
	}
	openSource.IsOpenSource = true
	if openSource.Nickname == "" {
		openSource.Nickname = openSource.Name
	}
	bytes, err := json.Marshal(openSource)
	if err != nil {
		return false, syserr.NewSysErr(err.Error())
	}
	return cache.RedisClient.HSet(publicSource, openSource.Name, string(bytes)).Result()
}

// get base data
// get source ( source or open source ) return source data json string
func (*RedisExternalInterface) GetSourceBaseData(name string) (interface{}, error) {
	privateSourceData, _ := cache.RedisClient.HGet(privateSource, name).Result()
	if privateSourceData != "" {
		source := new(Source)
		err := json.Unmarshal([]byte(privateSourceData), source)
		if err != nil {
			return nil, syserr.NewSysErr(err.Error())
		}
		return source, nil
	}
	publicSourceData, _ := cache.RedisClient.HGet(publicSource, name).Result()
	if publicSourceData != "" {
		openSource := new(OpenSource)
		err := json.Unmarshal([]byte(publicSourceData), openSource)
		if err != nil {
			return nil, syserr.NewSysErr(err.Error())
		}
		return openSource, nil
	}
	return "", syserr.NewSourceNotFound("source name not found")
}

// get private or public source rsa
// public key
func (*RedisExternalInterface) GetRsaPublicKey(name string) (*string,error) {
	privateSourceData, _ := cache.RedisClient.HGet(privateSource, name).Result()
	if privateSourceData != "" {
		source := new(Source)
		err := json.Unmarshal([]byte(privateSourceData), source)
		if err != nil {
			return nil, syserr.NewSysErr(err.Error())
		}
		return &source.PublicKey, nil
	}
	publicSourceData, _ := cache.RedisClient.HGet(publicSource, name).Result()
	if publicSourceData != "" {
		openSource := new(OpenSource)
		err := json.Unmarshal([]byte(publicSourceData), openSource)
		if err != nil {
			return nil, syserr.NewSysErr(err.Error())
		}
		return &openSource.PublicKey, nil
	}
	return nil, syserr.NewSourceNotFound("source name not found")
}
