package mini

import (
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/model"
)

type MiniService struct {
}

var miniService = new(MiniService)
var miniModelService = new(model.MiniModel)

// 创建小程序
func (*MiniService) createMini(request *restful.Request, response *restful.Response) {
	err := func() error {
		miniModel := new(model.MiniModel)
		err := request.ReadEntity(miniModel)
		if err != nil {
			return err
		}
		// 持久化并加入缓存
		err = miniModel.CreateAndJoinCache()
		if err != nil {
			return err
		}
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 根据id获取小程序基本信息
func (*MiniService) getOne(request *restful.Request, response *restful.Response) {
	miniInfo,err := func() (*model.MiniModel,error) {
		id := request.PathParameter("id")
		return miniModelService.FindByMiniId(id)
	}()
	rest.WriteEntity(miniInfo, err, response)
}

// 更新小程序
func (*MiniService) updateMini(request *restful.Request, response *restful.Response){
	err := func() error {
		miniModel := new(model.MiniModel)
		err:= request.ReadEntity(miniModel)
		if err != nil {
			return err
		}
		// update and Join Cache
		err = miniModel.UpdateAndJoinCache()
		return err
	}()
	rest.WriteEntity(nil,err,response)
}

// 删除小程序
func (*MiniService) delOne(request *restful.Request, response *restful.Response){
	err := func() error {
		id := request.PathParameter("id")
		return miniModelService.RemoveByMiniId(id)
	}()
	rest.WriteEntity(nil,err,response)
}

func init(){
	binder, webService := rest.NewJsonWebServiceBinder("/mini")
	webService.Route(webService.GET("{id}").To(miniService.getOne))
	binder.BindAdd()

	binderAdmin, webServiceAdmin := rest.NewJsonWebServiceBinder("/admin/mini")
	webServiceAdmin.Route(webServiceAdmin.POST("").To(miniService.createMini))
	webServiceAdmin.Route(webServiceAdmin.GET("{id}").To(miniService.getOne))
	webServiceAdmin.Route(webServiceAdmin.PUT("").To(miniService.updateMini))
	webServiceAdmin.Route(webServiceAdmin.DELETE("{id}").To(miniService.delOne))
	binderAdmin.BindAdd()

}
