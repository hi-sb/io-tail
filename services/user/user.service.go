package user

import (
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/model"
	"github.com/hi-sb/io-tail/utils"
	"time"
)

type UserService struct {
	//
}

//地址
var userService = new(UserService)
var userModelService = new(model.UserModel)

//用token 获取用户信息
func (*UserService) get(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	JWT, err := auth.GetJWT(token)
	user := new(model.UserModel)
	if err == nil {
		err = mysql.DB.Where("id =?", JWT.ID).First(user).Error
	}
	rest.WriteEntity(user, err, response)
}

// 注册并登陆
func (this *UserService) regOrlogin(request *restful.Request, response *restful.Response) {
	userModel, err := func() (*model.UserModel, error) {
		registerModel := new(model.RegisterModel)
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
		userModel := new(model.UserModel)
		userModel.MobileNumber = registerModel.MobileNumber
		userModel.NickName = registerModel.MobileNumber
		userModel.UserRole = 0 // 设置为普通用户
		return userModel, nil
	}()
	if err == nil {
		// 判断当前手机号已经持久化
		var user model.UserModel
		mysql.DB.Where("mobile_number = ?", userModel.MobileNumber).First(&user)

		if user.ID == "" {
			// 注册
			//err = mysql.Transactional(func(tx *gorm.DB) error {
			//	sync := lock.GetSync("register:" + userModel.MobileNumber)
			//	err := sync.Lock()
			//	if err != nil {
			//		return err
			//	}
			//	defer sync.Unlock()
			//	userModel.Bind()
			//	return tx.Create(userModel).Error
			//})
			err = userModel.CreateAndJoinCache()
		} else {
			userModel.ID = user.ID
		}
		//// 缓存用户信息
		//data, err := json.Marshal(userModel)
		//if err == nil {
		//	_, err = cache.RedisClient.HSet(constants.USER_BASE_INFO_REDIS_KEY, fmt.Sprintf(constants.USER_BASE_INFO_REDIS_PREFIX, userModel.ID), data).Result()
		//	if err != nil {
		//		println(err)
		//	}
		//}
	}
	// 完成注册 并登录
	var token string
	if err == nil {
		jwt := auth.JWT{
			UserName: userModel.MobileNumber,
			ID:       userModel.ID,
			Duration: time.Minute * 10080,
		}
		token, err = auth.CreateToken(&jwt)
	}
	rest.WriteEntity(token, err, response)
}

// 更新用户信息（昵称头像）
func (*UserService) updateInfO(request *restful.Request, response *restful.Response) {
	err := func() error {
		userId := utils.Strval(request.Attribute("currentUserId"))
		userMode := new(model.UserModel)
		err := request.ReadEntity(userMode)
		if err != nil {
			return err
		}
		userMode.ID = userId
		if userMode.NickName != "" {
			mysql.DB.Model(userMode).UpdateColumn("nick_name", userMode.NickName)
		}
		if userMode.Avatar != "" {
			mysql.DB.Model(userMode).UpdateColumn("avatar", userMode.Avatar)
		}
		userModelService.RefushCache(userId)
		// 刷新缓存
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

//从缓存获取简要信息
func (*UserService) briefly(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	user := new(model.UserModel).GetInfoById(id)
	var userBriefly model.UserBriefly
	if user != nil {
		userBriefly = model.UserBriefly{
			NickName: user.NickName,
			Avatar:   user.Avatar,
		}
	}
	rest.WriteEntity(userBriefly, nil, response)
}

// 设置管理员
func (*UserService) setAdmin(request *restful.Request, response *restful.Response){
	err := func() error {
		setInfo := new(model.SetAdmin)
		err:= request.ReadEntity(setInfo)
		if err != nil {
			return err
		}
		err = mysql.DB.Model(model.UserModel{}).Where("id = ?",setInfo.ID).UpdateColumn("user_role", setInfo.UserRole).Error
		return err
	}()
	rest.WriteEntity(nil,err,response)
}

func (*UserService) initAdmin(request *restful.Request, response *restful.Response){
	userModelService.InitADMIN()
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/user")
	webService.Route(webService.GET("/briefly/{id}").To(userService.briefly))
	webService.Route(webService.GET("/{token}").To(userService.get))
	webService.Route(webService.POST("/login").To(userService.regOrlogin))
	webService.Route(webService.PUT("/update").To(userService.regOrlogin))
	webService.Route(webService.GET("/admin/init").To(userService.initAdmin))
	binder.BindAdd()


	adminBinder, adminWebService := rest.NewJsonWebServiceBinder("/admin/user")
	adminWebService.Route(webService.PUT("/update").To(userService.regOrlogin))
	adminBinder.BindAdd()



}

