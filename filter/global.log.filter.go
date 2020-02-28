package filter

import (
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/log"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Log.Printf("%s:access-log-[%s] %s,%s", req.Request.Proto, req.Request.RemoteAddr, req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}


