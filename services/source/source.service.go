package source

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/abstract"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/body"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/core/topic"
	"github.com/hi-sb/io-tail/model"
	"net/http"
	"os"
	"strconv"
	"time"
)

//
type SourceService struct {
}

var (
	sourceService     = new(SourceService)
	filePathAdapter   = abstract.NewDefaultFilePathAdapter()
	readAndWrite      = abstract.NewDefaultReadAndWriteAdapter()
	permissionService = new(PermissionService)
)

//
// 监听私有资源
// 也就是监听自己的消息话题，当有人发送消息到该资源，那么会发送一个消息到监听者
func (sourceService *SourceService) privateSourceListen(request *restful.Request, response *restful.Response) {
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
		if JWT.AtNum != source {
			return "", "", syserr.NewTokenAuthError("拒绝访问")
		}
		offset := request.QueryParameter("offset")
		var offsetInt int64
		if offset != "" {
			offsetInt, err = strconv.ParseInt(offset, 10, 64)
		}
		if err != nil {
			return "", "", syserr.NewBadRequestErr("错误的参数 offset")
		}
		path, err := filePathAdapter.Handle(source)
		if err != nil {
			err = syserr.NewSysErr(err.Error())
			fmt.Println(err)
			return "", "", err
		}
		tell := topic.NewDefaultTell(offsetInt)
		return JWT.AtNum, source, tell.TellMessage(topic.TellChan{Error: errChan, Reader: readChan}, path, request.Request)
	}()
	sourceService.tellChan(openid, source, errChan, readChan, response, err)
}

//
//监听一个共有的资源
// 也就是群消息，群消息我们认为是一个公共的开放的资源
// 只要加入了该群则可以发送消息到该资源，也就是往该资源写入消息，此时监听者将收到一个消息
func (sourceService *SourceService) publicSourceListen(request *restful.Request, response *restful.Response) {
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
			return "", "", syserr.NewBadRequestErr("错误的offset参数")
		}
		tell := topic.NewDefaultTell(offsetInt)
		path, err := filePathAdapter.Handle(source)
		if err != nil {
			err = syserr.NewSysErr(err.Error())
			fmt.Println(err)
			return "", "", err
		}
		return JWT.AtNum, source, tell.TellMessage(topic.TellChan{Error: errChan, Reader: readChan}, path, request.Request)
	}()
	sourceService.tellChan(openid, source, errChan, readChan, response, err)
}

//通过tell监听资源文件
func (sourceService *SourceService) tellChan(name string, source string, errChan chan error, readChan chan *body.Message, response *restful.Response, err error) {
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

// 获取一个资源对应一个有权访问的访问者的offset
// 也就是话题文件的读取位置
// 一个资源对应一个话题文件，每个人对于一个话题文件都有一个读取位置。
func (sourceService *SourceService) offset(request *restful.Request, response *restful.Response) {
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

//发送消息到话题资源
//也就是往某个话题文件写入消息，这个时候需要验证权限，当向一个私有话题发送消息的时候
//会验证是否是对方的好友，如果是，那么则具有写入话题的权限。而共有话题，则会验证该用户是否加入了该群。
func (sourceService *SourceService) send(request *restful.Request, response *restful.Response) {
	err := func() error {
		sendRequest := new(model.SendRequest)
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
		return sourceService.SendMessage(JWT.ID, source, sendRequest)
	}()
	rest.WriteEntity(nil, err, response)
}

//发送消息到话题资源
//也就是往某个话题文件写入消息，这个时候需要验证权限，当向一个私有话题发送消息的时候
//会验证是否是对方的好友，如果是，那么则具有写入话题的权限。而共有话题，则会验证该用户是否加入了该群。
func (sourceService *SourceService) groupSend(request *restful.Request, response *restful.Response) {
	err := func() error {
		sendRequest := new(model.SendRequest)
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
		err = permissionService.CheckGroupWritePermission(JWT, source)
		if err != nil {
			return err
		}
		return sourceService.SendMessage(JWT.ID, source, sendRequest)
	}()
	rest.WriteEntity(nil, err, response)
}

func (sourceService *SourceService) SendMessage(fromId string, toId string, sendRequest *model.SendRequest) error {
	path, err := filePathAdapter.Handle(toId)
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
		return syserr.NewSourceNotFound("没有这样的群或用户")
	}
	defer messageFile.Close()
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
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/topic")
	webService.Route(webService.GET("offset/{source}").To(sourceService.offset))
	webService.Route(webService.GET("/{source}").To(sourceService.privateSourceListen))
	webService.Route(webService.PUT("/{source}").To(sourceService.send))
	webService.Route(webService.PUT("public/{source}").To(sourceService.groupSend))
	webService.Route(webService.GET("open/{source}").To(sourceService.publicSourceListen))
	binder.BindAdd()
}
