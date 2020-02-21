package lock

import (
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	Redsync *redsync.Redsync
	//默认 10秒超时
	DefaultTimeOutOption = redsync.SetExpiry(10 * time.Second)
)

//初始化 服务
func InitRedisLock(url string, pw string) {
	p := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://" + url)
			if _, err := c.Do("AUTH", pw); err != nil {
				defer c.Close()
				return nil, err
			}
			return c, err
		},
	}
	Redsync = redsync.New([]redsync.Pool{p})
}

// 获取默认的 Redsync
func GetSync(lock string) *redsync.Mutex {
	return Redsync.NewMutex(lock, DefaultTimeOutOption)
}

// 获取指定超时的 Redsync
func GetSyncByOption(lock string, timeOutOption redsync.Option) *redsync.Mutex {
	return Redsync.NewMutex(lock, timeOutOption)
}
