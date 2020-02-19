package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
)


// http service
func httpService(httpAddr *string) {
	fmt.Println("Start http server listen build addr is ", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, restful.DefaultContainer); err != nil {
		panic(err)
	}
}
