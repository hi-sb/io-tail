package rest

import (
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/emicklei/go-restful"
)

// err model
type ResponseModel struct {
	Message string
	Code    int
	Body    interface{}
	Success bool
}

func WriteEntity(obj interface{}, err error, response *restful.Response) {
	responseModel := ResponseModel{Message: "OK", Code: 200, Success: true, Body: obj}
	if err != nil {
		responseModel.Body = nil
		// default error code
		responseModel.Code = 10000
		responseModel.Success = false
		responseModel.Message = err.Error()
		// is base err
		if baseErr, ok := err.(syserr.BaseErrorInterface); ok {
			responseModel.Code = baseErr.Code()
		}
		_ = response.WriteHeaderAndEntity(500, responseModel)
		return
	}
	_ = response.WriteEntity(responseModel)
}


// 特殊需要
func WriteErrAndEntity(obj interface{}, err error, response *restful.Response) {
	responseModel := ResponseModel{Message: "OK", Code: 200, Success: true, Body: obj}
	if err != nil {
		responseModel.Message = err.Error()
		// is base err
		if baseErr, ok := err.(syserr.BaseErrorInterface); ok {
			responseModel.Code = baseErr.Code()
		}
		_ = response.WriteHeaderAndEntity(500, responseModel)
		return
	}
	_ = response.WriteEntity(responseModel)
}
