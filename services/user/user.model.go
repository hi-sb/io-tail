package user

import (
	"errors"
	"github.com/hi-sb/io-tail/core/db"
)


// user model
type UserModel struct {
	// base model
	db.BaseModel
	//密码
	Password string `gorm:"type:varchar(32);not null"`
	// 手机号
	MobileNumber string `gorm:"type:varchar(15);not null"`

	//昵称
	NickName string
	// 头像
	Avatar string

	// 公私钥
	PrvKey string
	PubKey string

}

// 注册 model
type RegisterModel struct {
	//手机号
	MobileNumber string
	//verify code
	VerifyCode string
	//密码
//	Password string
}


// 快捷登录模型
type QuickLogin struct {
	//手机号
	MobileNumber string
	//verify code
	VerifyCode string
}


//检查
func (this *RegisterModel) Check() error {
	if len(this.MobileNumber) == 0 {
		return errors.New("手机号不能为空")
	}

	if len(this.VerifyCode) == 0 {
		return errors.New("验证码不能为空")
	}
	return nil
}
