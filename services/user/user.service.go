package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/lock"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type UserService struct {
	//
}
//地址
var userService = new(UserService)

const (
	// key
	USER_BASE_INFO_REDIS_KEY = "USER_BASE_INFO"
	//  field
	USER_BASE_INFO_REDIS_PREFIX = "USER_BASE_INFO_"
)

//用token 获取用户信息
func (*UserService) get(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	JWT, err := auth.GetJWT(token)
	user := new(UserModel)
	if err == nil {
		err = mysql.DB.Where("id =?", JWT.ID).First(user).Error
	}
	rest.WriteEntity(user, err, response)
}

// 注册并登陆
func (this *UserService) regOrlogin(request *restful.Request, response *restful.Response) {
	userModel, err := func() (*UserModel, error) {
		registerModel := new(RegisterModel)
		err := request.ReadEntity(registerModel)
		if err != nil {
			return nil, err
		}
		err = registerModel.Check()
		if err != nil {
			return nil, err
		}
		verifyId := utils.Md5V2(registerModel.MobileNumber)
		isVerify := verifyService.Verify(verifyId, registerModel.VerifyCode, true)
		if !isVerify {
			return nil, errors.New("验证码错误")
		}
		userModel := new(UserModel)
		userModel.MobileNumber = registerModel.MobileNumber
		userModel.NickName = registerModel.MobileNumber
		return userModel, nil
	}()
	if err == nil {
		// 判断当前手机号已经持久化
		var user UserModel
		mysql.DB.Where("mobile_number = ?", userModel.MobileNumber).First(&user)

		if user.ID == "" {
			// 注册
			err = mysql.Transactional(func(tx *gorm.DB) error {
				sync := lock.GetSync("register:" + userModel.MobileNumber)
				err := sync.Lock()
				if err != nil {
					return err
				}
				defer sync.Unlock()
				userModel.Bind()
				return tx.Create(userModel).Error
			})
		}else {
			userModel.ID = user.ID
		}
		// 缓存用户信息
		data,err := json.Marshal(userModel)
		if err == nil {
			_,err = cache.RedisClient.HSet(USER_BASE_INFO_REDIS_KEY,fmt.Sprintf(USER_BASE_INFO_REDIS_PREFIX,userModel.ID),data).Result()
			if err !=nil {
				println(err)
			}
		}
	}
	// 完成注册 并登录
	var token string
	if err == nil {
		jwt := auth.JWT{
			UserName: userModel.MobileNumber,
			ID:       userModel.ID,
			Duration:  time.Minute*10080,
		}
		token, err = auth.CreateToken(&jwt)
	}
	rest.WriteEntity(token, err, response)
}




func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/user")
	webService.Route(webService.GET("/{token}").To(userService.get))
	webService.Route(webService.POST("/login").To(userService.regOrlogin))
	binder.BindAdd()
}