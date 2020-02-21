package utils

import "strings"

// get host form by
// hostAndPort
func GetHost(hostAndPort string) string {
	//获取请求ip
	remoteAddr := strings.Split(hostAndPort, ":")
	var remoteHost string
	if len(remoteAddr) == 2 {
		remoteHost = remoteAddr[0]
	}
	return remoteHost
}
