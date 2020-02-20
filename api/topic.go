package api

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/abstract"
	"github.com/hi-sb/io-tail/auth"
	"github.com/hi-sb/io-tail/body"
	"github.com/hi-sb/io-tail/cache"
	"github.com/hi-sb/io-tail/rest"
	service "github.com/hi-sb/io-tail/services"
	"github.com/hi-sb/io-tail/syserr"
	"github.com/hi-sb/io-tail/topic"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	topicApi          = new(TopicApi)
	filePathAdapter   = abstract.NewDefaultFilePathAdapter()
	readAndWrite      = abstract.NewDefaultReadAndWriteAdapter()
	permissionService = service.PermissionService{}
)

// send request
type SendRequest struct {
	// send time
	SendTime int64
	// message body
	Body string
	// message type
	ContentType string
}

// topic rest http service
type TopicApi struct {
}

//
// user private listen
func (topicApi *TopicApi) privateSourceListen(request *restful.Request, response *restful.Response) {
	errChan := make(chan error)
	readChan := make(chan *body.Message)
	openid, source, err := func() (string, string, error) {
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return "", "", err
		}
		source := request.PathParameter("source")
		// this service user
		// and id == source
		if JWT.AtNum != source || JWT.Type != auth.TokenTypeUser {
			return "", "", syserr.NewTokenAuthError("access denied")
		}
		offset := request.QueryParameter("offset")
		var offsetInt int64
		if offset != "" {
			offsetInt, err = strconv.ParseInt(offset, 10, 64)
		}
		if err != nil {
			return "", "", syserr.NewBadRequestErr("offset bad request")
		}
		tell := topic.NewDefaultTell(offsetInt)
		return JWT.AtNum, source, tell.TellMessage(topic.TellChan{Error: errChan, Reader: readChan}, request.Request)
	}()
	topicApi.tellChan(openid, source, errChan, readChan, response, err)
}

//
// open source listen
func (topicApi *TopicApi) publicSourceListen(request *restful.Request, response *restful.Response) {
	errChan := make(chan error)
	readChan := make(chan *body.Message)
	openid, source, err := func() (string, string, error) {
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return "", "", err
		}
		source := request.PathParameter("source")
		offset := request.QueryParameter("offset")
		var offsetInt int64
		if offset != "" {
			offsetInt, err = strconv.ParseInt(offset, 10, 64)
		}
		if err != nil {
			return "", "", syserr.NewBadRequestErr("offset bad request")
		}
		tell := topic.NewDefaultTell(offsetInt)
		return JWT.AtNum, source, tell.TellMessage(topic.TellChan{Error: errChan, Reader: readChan}, request.Request)
	}()
	topicApi.tellChan(openid, source, errChan, readChan, response, err)
}

// send
func (topicApi *TopicApi) tellChan(name string, source string, errChan chan error, readChan chan *body.Message, response *restful.Response, err error) {
	if err != nil {
		rest.WriteEntity(nil, err, response)
		return
	}
	flusher, _ := response.ResponseWriter.(http.Flusher)
	for {
		select {
		// tail message err
		case err := <-errChan:
			fmt.Println(err.Error())
			return
		// do message
		case message := <-readChan:
			data, err := readAndWrite.Encoding(message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			key := fmt.Sprintf("%s_offset_from_%s", name, source)
			cache.RedisClient.Set(key, message.Offset, 0)
			_, _ = fmt.Fprint(response.ResponseWriter, data)
			flusher.Flush()
		}
	}
}

// get resource offset
func (*TopicApi) offset(request *restful.Request, response *restful.Response) {
	offset, err := func() (int64, error) {
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return 0, err
		}
		source := request.PathParameter("source")
		key := fmt.Sprintf("%s_offset_from_%s", JWT.AtNum, source)
		value, _ := cache.RedisClient.Get(key).Result()
		if value == "" {
			value = "0"
		}
		return strconv.ParseInt(value, 10, 64)
	}()
	rest.WriteEntity(offset, err, response)
}

// send
func (*TopicApi) send(request *restful.Request, response *restful.Response) {
	err := func() error {
		sendRequest := new(SendRequest)
		err := request.ReadEntity(sendRequest)
		if err != nil {
			return syserr.NewBadRequestErr(err.Error())
		}
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return err
		}
		source := request.PathParameter("source")
		err = permissionService.CheckWritePermission(JWT, source)
		if err != nil {
			return err
		}
		path, err := filePathAdapter.Handle(request.Request.RequestURI)
		if err != nil {
			err = syserr.NewSysErr(err.Error())
			fmt.Println(err)
			return err
		}
		// check content type
		err = new(body.Message).CheckContentType(sendRequest.ContentType)
		if err != nil {
			return syserr.NewContentTypeErr(err.Error())
		}
		if sendRequest.SendTime == 0 {
			sendRequest.SendTime = time.Now().UnixNano() / 1e6
		}
		messageFile, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println(err)
			return syserr.NewSourceNotFound("topic not found")
		}
		defer messageFile.Close()
		var fromId = JWT.AtNum
		message := body.Message{
			FormId:      fromId,
			SendTime:    sendRequest.SendTime,
			Body:        sendRequest.Body,
			ContentType: sendRequest.ContentType,
		}
		body, err := readAndWrite.Encoding(&message)
		if err != nil {
			err = syserr.NewSysErr(err.Error())
			fmt.Println(err)
			return err
		}
		_, err = messageFile.WriteString(body)
		return err
	}()
	rest.WriteEntity(nil, err, response)
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/topic")
	webService.Route(webService.GET("offset/{source}").To(topicApi.offset))
	webService.Route(webService.GET("/{source}").To(topicApi.privateSourceListen))
	webService.Route(webService.PUT("/{source}").To(topicApi.send))
	webService.Route(webService.GET("open/{source}").To(topicApi.publicSourceListen))
	binder.BindAdd()
}
