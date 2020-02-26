package user

import (
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/lock"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/utils"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type UserService struct {
	//
}
//地址
var userService = new(UserService)


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

// 根据id获取用户信息
func (*UserService) GetInfoById(ID string)*UserModel{
	user := new(UserModel)
	err := mysql.DB.Where("id =?", ID).First(user).Error
	if err != nil {
		return nil
	}
	return user
}
// 根据ids获取用户信息
func (*UserService) GetInfoByIds(ids *[]string)*[]UserModel{
	var users []UserModel
	idArrayStr := strings.Replace(strings.Trim(fmt.Sprint(*ids), "[]"), " ", ",", -1)
	err := mysql.DB.Where("id in (?)", idArrayStr).Find(&users).Error
	if err != nil {
		return nil
	}
	return &users
}

// 根据手机号查询用户信息
func (*UserService) GetInfoByPhone(phone string) *UserModel {
	user := new(UserModel)
	err := mysql.DB.Where("mobile_number =?", phone).First(user).Error
	if err != nil {
		return nil
	}
	return user
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