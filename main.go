package main

import (
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	_ "github.com/hi-sb/io-tail/api"
	"github.com/hi-sb/io-tail/cache"
	"github.com/hi-sb/io-tail/topic"
	"net/http"
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


func main() {
	buildAddr := flag.String("build", ":7654", "server http buildAddr")
	dataPath := flag.String("dataPath", "./data", " message data path")
	redis := flag.String("redis", "127.0.0.1:6379", "redis connect host and port ")
	redisPass := flag.String("redisPass", "", "redis connect password")
	flag.Parse()
	printASCIILogo()
	topic.SetDataPath(*dataPath)
	cache.RedisServiceClientInit(*redis, *redisPass)
	httpService(buildAddr)
}
