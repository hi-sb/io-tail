package rest

import "github.com/emicklei/go-restful"


// get json WebService
func getJsonWebService(root string) *restful.WebService {
	return new(restful.WebService).Path(root).Produces(restful.MIME_JSON)
}

// json WebService
type JsonWebServiceBinder struct {
	webService *restful.WebService
}

// new create JsonWebService
func NewJsonWebServiceBinder(root string) (*JsonWebServiceBinder, *restful.WebService) {
	webService := getJsonWebService(root)
	return &JsonWebServiceBinder{webService: webService}, webService
}

// Add to Container
func (jsonWebServiceBinder *JsonWebServiceBinder) BindAdd() {
	restful.Add(jsonWebServiceBinder.webService)
}
