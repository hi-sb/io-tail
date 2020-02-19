package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitee.com/saltlamp/im-service/auth"
	"gitee.com/saltlamp/im-service/cache"
	"gitee.com/saltlamp/im-service/config"
	"gitee.com/saltlamp/im-service/ext"
	"gitee.com/saltlamp/im-service/rest"
	"gitee.com/saltlamp/im-service/syserr"
	"gitee.com/saltlamp/im-service/utils"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
)

const (
	AuthRandomForm = "auth_%s_random_from_%s"
)

var (
	authApi           = new(AuthApi)
	externalInterface = ext.GetExternalInterface()
)

// CheckIdentity request body
type CheckCodeBody struct {
	// at num
	// Other message service user description
	AtNum string
	// code md5
	// consist of login code and  identity code
	// login:xxxxxxx or identity:xxxxxxx md5  value
	Code string
}

// auth
type AuthApi struct {
}

// get login code
func (authApi *AuthApi) loginCode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	randomKey, err := authApi.getRandomKey("login", name)
	rest.WriteEntity(randomKey, err, response)
}

// check randomCode code
//login this message service
func (authApi *AuthApi) login(request *restful.Request, response *restful.Response) {
	checkCodeBody := new(CheckCodeBody)
	err := authApi.checkRandomKeyMd5("login", checkCodeBody, request)
	var token string
	if err == nil {
		jwt := auth.JWT{
			AtNum:       checkCodeBody.AtNum,
			Type:     auth.TokenTypeUser,
			Duration: &config.UserDuration,
		}
		token, err = auth.CreateToken(&jwt)
	}
	rest.WriteEntity(token, err, response)
}

// get identity code
func (authApi *AuthApi) identityCode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	randomKey, err := authApi.getRandomKey("identity", name)
	rest.WriteEntity(randomKey, err, response)
}

// check randomCode code
//Federation auth, which is used by other nodes to confirm the identity of a user
func (authApi *AuthApi) checkIdentity(request *restful.Request, response *restful.Response) {
	checkCodeBody := new(CheckCodeBody)
	err := authApi.checkRandomKeyMd5("identity", checkCodeBody, request)
	rest.WriteEntity("ok", err, response)
}

// visitor Apply registration
//When we need to send messages to the current message service node, we need to register visitors.
//
//1. Federation verify the authenticity of visitors
//
//2. Issue a visitor token
func (authApi *AuthApi) visitorApplyRegistration(request *restful.Request, response *restful.Response) {
	randomKey, err := func() (string, error) {
		atNum := request.PathParameter("atNum")
		account, err := utils.GetAccount(atNum)
		if err != nil {
			return "", err
		}
		urlModel := "http://%s:%d/api/identity-code/%s"
		url := fmt.Sprintf(urlModel, account.MessageServiceHost, account.MessageServicePort, account.Name)
		res, err := http.Get(url)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		if res.StatusCode != 200 {
			return "", syserr.NewFederationAuthErr(string(body))
		}
		responseModel := new(rest.ResponseModel)
		err = json.Unmarshal(body, responseModel)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		if !responseModel.Success {
			return "", syserr.NewFederationAuthErr(responseModel.Message)
		}
		return responseModel.Body.(string), nil
	}()
	rest.WriteEntity(randomKey, err, response)
}

// visitor Apply registration
//When we need to send messages to the current message service node, we need to register visitors.
//
//1. Federation verify the authenticity of visitors
//
//2. Issue a visitor token
func (authApi *AuthApi) visitorRegistration(request *restful.Request, response *restful.Response) {
	token, err := func() (string, error) {
		visitorRegistrationBody := new(CheckCodeBody)
		err := request.ReadEntity(visitorRegistrationBody)
		// bad request
		if err != nil {
			return "", syserr.NewBadRequestErr(err.Error())
		}
		//atNum is null
		if visitorRegistrationBody.AtNum ==""{
			return "",syserr.NewParameterError("AtNum is nll")
		}
		//code is null
		if visitorRegistrationBody.Code ==""{
			return "",syserr.NewParameterError("Code is nll")
		}
		account, err := utils.GetAccount(visitorRegistrationBody.AtNum)
		if err != nil {
			return "", err
		}
		urlModel := "http://%s:%d/api/check-identity"
		// check request body
		checkCodeBody := CheckCodeBody{
			// other message service user
			AtNum: visitorRegistrationBody.AtNum,
			// visitor user identity check code .
			Code: visitorRegistrationBody.Code,
		}
		jsonBodyStr, err := json.Marshal(checkCodeBody)
		url := fmt.Sprintf(urlModel, account.MessageServiceHost, account.MessageServicePort)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBodyStr))
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		if res.StatusCode != 200 {
			return "", syserr.NewFederationAuthErr(string(body))
		}
		responseModel := new(rest.ResponseModel)
		err = json.Unmarshal(body, responseModel)
		if err != nil {
			return "", syserr.NewFederationAuthErr(err.Error())
		}
		if !responseModel.Success {
			return "", syserr.NewFederationAuthErr(responseModel.Message)
		}
		// return visitor token
		jwt := auth.JWT{
			AtNum:       visitorRegistrationBody.AtNum,
			Type:     auth.TokenTypeVisitor,
			Duration: &config.VisitorDuration,
		}
		return auth.CreateToken(&jwt)
	}()
	rest.WriteEntity(token, err, response)
}

//
func (authApi *AuthApi) checkRandomKeyMd5(codeType string, checkCodeBody *CheckCodeBody, request *restful.Request) error {
	err := request.ReadEntity(checkCodeBody)
	// bad request
	if err != nil {
		return syserr.NewBadRequestErr(err.Error())
	}
	//atNum is null
	if checkCodeBody.AtNum ==""{
		return syserr.NewParameterError("atNum is nll")
	}
	//code is null
	if checkCodeBody.Code ==""{
		return syserr.NewParameterError("Code is nll")
	}
	authRandomCacheKey := fmt.Sprintf(AuthRandomForm, codeType, checkCodeBody.AtNum)
	// Use random in session
	random, _ := cache.RedisClient.Get(authRandomCacheKey).Result()
	randomCode := fmt.Sprintf("%s:%s", codeType, random)
	randomCodeMd5 := utils.Md5V2(randomCode)
	// check success
	if randomCodeMd5 == checkCodeBody.Code {
		_ = cache.RedisClient.Del(authRandomCacheKey)
		return nil
	} else {
		return syserr.NewCheckRandomKeyError("Not recorded")
	}
}

// getCode
func (*AuthApi) getRandomKey(codeType string, name string) (string, error) {
	if name == "" {
		return "", syserr.NewIdIsNullAuthError("atNum is null")
	}
	value, _ := externalInterface.GetUserPublicKey(name)
	if value == "" {
		return "", syserr.NewUnboundAuthError("unbound user")
	}
	authRandomCacheKey := fmt.Sprintf(AuthRandomForm, codeType, name)
	// If the session exists
	// Use random in session
	random, _ := cache.RedisClient.Get(authRandomCacheKey).Result()
	if random == "" {
		random = utils.GetID()
	}
	publicKey, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", syserr.NewPublicKeyAuthError("public key err")
	}
	code := fmt.Sprintf("%s:%s", codeType, random)
	randomKey, err := utils.RsaEncryptAndBase64([]byte(code), publicKey)
	if err != nil {
		return "", syserr.NewPublicKeyAuthError("public key err")
	}
	cache.RedisClient.Set(authRandomCacheKey, random, config.AuthCodeTimeOut)
	return randomKey, nil
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/auth")
	webService.Route(webService.GET("/login-code/{name}").To(authApi.loginCode))
	webService.Route(webService.GET("/identity-code/{name}").To(authApi.identityCode))
	webService.Route(webService.POST("/login").To(authApi.login))
	webService.Route(webService.POST("/check-identity").To(authApi.checkIdentity))
	webService.Route(webService.GET("/visitor-apply-registration/{atNum}").To(authApi.visitorApplyRegistration))
	webService.Route(webService.POST("/visitor-registration").To(authApi.visitorRegistration))
	binder.BindAdd()
}
