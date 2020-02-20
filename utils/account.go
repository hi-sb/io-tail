package utils

import (
	"github.com/hi-sb/io-tail/syserr"
	"regexp"
	"strconv"
	"strings"
)

// @ num
type Account struct {
	// user name
	Name string
	// message service host
	MessageServiceHost string
	// port
	MessageServicePort int
}

// get account
func GetAccount(atNum string) (*Account, error) {
	if strings.Index(atNum, ":") < 0 {
		atNum += ":7654"
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+(:[0-9]{1,5})+$", atNum)
	if !match || err != nil {
		return nil, syserr.NewParameterError("wrong account format")
	}
	// fmt
	nameAndServiceHostAndPort := strings.Split(atNum, "@")
	var name = nameAndServiceHostAndPort[0]
	serviceHostAndPort := strings.Split(nameAndServiceHostAndPort[1], ":")
	var serviceHost = serviceHostAndPort[0]
	messageServicePort, _ := strconv.Atoi(serviceHostAndPort[1])
	return &Account{Name: name, MessageServiceHost: serviceHost, MessageServicePort: messageServicePort}, nil
}
