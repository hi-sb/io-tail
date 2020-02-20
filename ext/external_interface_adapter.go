package ext

import "github.com/hi-sb/io-tail/config"

var (
	externalInterfaceRedisObj = new(RedisExternalInterface)
)

func GetExternalInterface() ExternalInterface {
	switch config.ExternalInterface {
	case config.ExternalInterfaceRedis:
		{
			return externalInterfaceRedisObj
		}
	default:
		{
			return externalInterfaceRedisObj
		}
	}
}
