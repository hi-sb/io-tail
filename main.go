package main

import (
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/lock"
	"github.com/hi-sb/io-tail/core/topic"
	_ "github.com/hi-sb/io-tail/services/sms"
	_ "github.com/hi-sb/io-tail/services/user"
	"net/http"
	"os"
)

// log logo
func printASCIILogo() {
	logo := `
			 ___   _______    _______  _______  ___   ___     
			|   | |       |  |       ||   _   ||   | |   |    
			|   | |   _   |  |_     _||  |_|  ||   | |   |    
			|   | |  | |  |    |   |  |       ||   | |   |    
			|   | |  |_|  |    |   |  |       ||   | |   |___ 
			|   | |       |    |   |  |   _   ||   | |       |
			|___| |_______|    |___|  |__| |__||___| |_______|
	`
	fmt.Println(logo)
}


// http service
func httpService(httpAddr *string) {
	fmt.Println("Start http server listen build addr is ", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, restful.DefaultContainer); err != nil {
		panic(err)
	}
}

// mysql service
func MysqlService(addr *string, showSql *bool) {
	// add &parseTime=True
	issucc := mysql.GetInstance().InitDataPool(*addr+"&parseTime=True", *showSql)
	if !issucc {
		os.Exit(1)
	}
}


// redis lock service
func RedisLockService(url *string, pw *string) {
	lock.InitRedisLock(*url, *pw)
}

func main() {
	buildAddr := flag.String("build", ":7654", "server http buildAddr")
	dataPath := flag.String("dataPath", "./data", " message data path")
	redis := flag.String("redis", "127.0.0.1:6379", "redis connect host and port ")
	redisPass := flag.String("redisPass", "", "redis connect password")
	showSql := flag.Bool("showSql", true, "show sql is :true or false")
	mysqlAddr := flag.String("mysqlAddr", "", "mysql addr root:xxx@tcp(127.0.0.1:3306)/dbname?charset=utf8")
	flag.Parse()
	printASCIILogo()
	topic.SetDataPath(*dataPath)
	cache.RedisServiceClientInit(*redis, *redisPass)
	RedisLockService(redis, redisPass)
	MysqlService(mysqlAddr, showSql)
	httpService(buildAddr)
}
