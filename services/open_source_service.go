package service

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/syserr"
)

const (
	// user source
	privateSource = "private_source"
	// open source
	publicSource = "public_source"
)

//
type OpenSource struct {
	//
	// open source name
	Name string
	//
	ProfilePhotoUrl string
	//
	Describe string
	// create name
	CreateName string
}

//
type OpenSourceService struct {

}


// create open source
// If public key is null .
// The content will not be encrypted when it is stored, and no password is required when it is join
//
func (*OpenSourceService) CreateOpenSource(openSource *OpenSource) (bool, error) {
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