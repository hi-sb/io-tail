package model

import (
	"errors"
	"github.com/hi-sb/io-tail/utils"
)

//验证码
type VerifyModel struct {
	// base 64 data
	Data string
	// id
	Id string
}

//发送短信验证码消息体
type SmsVerify struct {
	// num
	MobileNumber string
	// img Verify id
	VerifyId string
	// num
	VerifyNum string
}

// check
func (this *SmsVerify) Check() error {
	if len(this.MobileNumber) == 0 {
		return errors.New("手机号不能为空")
	}


	isMobile := utils.VerifyMobileFormat(this.MobileNumber)
	if !isMobile {
		return errors.New("请输入正确的手机号")
	}
	return nil
}
