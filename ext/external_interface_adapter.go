package ext

import "gitee.com/saltlamp/im-service/config"

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
