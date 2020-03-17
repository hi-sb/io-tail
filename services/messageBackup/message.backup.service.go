package messageBackup

import (
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/base"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/model"
	"strconv"
)

var (
	messageBackup        = new(model.MessageBackup)
	messageBackupService = new(MessageBackupService)
)

type MessageBackupService struct {
}

func (*MessageBackupService) privateMessagePage(request *restful.Request, response *restful.Response) {
	page, err := func() (*base.Pager, error) {
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return nil, err
		}
		var page base.Pager
		err = request.ReadEntity(&page)
		if err != nil {
			return nil, err
		}
		fId := request.PathParameter("fId")
		sendTimeStr := request.PathParameter("sendTime")
		var sendTime int64
		if sendTimeStr != "" {
			sendTime, err = strconv.ParseInt(sendTimeStr, 10, 64)
		}
		return messageBackup.PrivateMessagePage(page, JWT.ID, fId, sendTime)
	}()
	rest.WriteEntity(page, err, response)
}

func (*MessageBackupService) groupMessagePage(request *restful.Request, response *restful.Response) {
	page, err := func() (*base.Pager, error) {
		token := request.HeaderParameter(auth.AUTH_HEADER)
		JWT, err := auth.GetJWT(token)
		if err != nil {
			return nil, err
		}
		var page base.Pager
		err = request.ReadEntity(&page)
		if err != nil {
			return nil, err
		}
		groupId := request.PathParameter("groupId")
		sendTimeStr := request.PathParameter("sendTime")
		var sendTime int64
		if sendTimeStr != "" {
			sendTime, err = strconv.ParseInt(sendTimeStr, 10, 64)
		}
		return messageBackup.GroupMessagePage(page, JWT.ID, groupId, sendTime)
	}()
	rest.WriteEntity(page, err, response)
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/message-backup")
	webService.Route(webService.POST("/private/{fId}/{sendTime}").To(messageBackupService.privateMessagePage))
	webService.Route(webService.POST("/group/{groupId}/{sendTime}").To(messageBackupService.groupMessagePage))
	binder.BindAdd()
}
