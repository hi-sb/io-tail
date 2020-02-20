package main

import (
	"flag"
	"fmt"
	_ "github.com/hi-sb/io-tail/api"
	"github.com/hi-sb/io-tail/cache"
	"github.com/hi-sb/io-tail/topic"
)

// log logo
func printASCIILogo() {
	logo := `
 _______           _____    
|__   __|         |  __ \     
   | | ___ _ __   | |  | | ___  __ _ _ __ ___  ___  ___ 
   | |/ _ \ '_ \  | |  | |/ _ \/ _' | '__/ _ \/ _ \/ __|
   | |  __/ | | | | |__| |  __/ (_| | | |  __/  __/\__ \
   |_|\___|_| |_| |_____/ \___|\__, |_|  \___|\___||___/
                                __/ |                   
                               |___/
	`
	fmt.Println(logo)
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
	cache.RedisClient.HSet("private_source", "huangxing", "{\"Name\":\"huangxing\",\"PublicKey\":\"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDdmE5ZTlUUHRQVHpPdTFMQ3NRaTllVkg3ZwpvaUF3V2xVQ210UlhBdHdxS2dtUjZZYnpaV1ZhMFliVC9WbDJKUFcrcFBKcnhwcHdZS0J5MDZPd0VCL0JVc1hOCkU2UDJnS2xYclJOaFE5Sk1PQjJ5eEM3UXdYV2ppTUVUdXUwNFVXeG9uN3RKL1l4cW5iblZGQlNXbmxhZ3M0eXkKWlNwWFg1ZnE2UlJ1UUFSYjFRSURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=\",\"Nickname\":\"大虾\",\"ProfilePhotoUrl\":\"http://www.17qq.com/img_qqtouxiang/88502343.jpeg\"}")
	httpService(buildAddr)
}
