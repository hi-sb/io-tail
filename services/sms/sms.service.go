package sms

import (
	"errors"
	"fmt"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

//sms service
type SmsService struct {
}

var (
	httpClient = new(http.Client)
)

const (
	sendUrlModel = "http://mb345.com:999/WS/BatchSend2.aspx?CorpID=LKSDK0005555&Pwd=yrwl@188&Mobile=%s&Content=%s&Cell=&SendTime="
)

var SmsServiceObj = new(SmsService)

func (this *SmsService) Send(mobile string, content string, ip string) error {
	err := this.CheckSend(mobile, content, ip)
	if err != nil {
		return err
	}
	body := utils.Utf8ToGBK(content)
	body = url.QueryEscape(body)
	requestUrl := fmt.Sprintf(sendUrlModel, mobile, body)
	response, err := httpClient.Get(requestUrl)
	if err != nil {
		return err
	}
	byte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	responseBody := string(byte)
	smsLog := SmsLogModel{
		IP:           ip,
		MobileNumber: mobile,
		Content:      content,
		ResponseBody: responseBody,
	}
	smsLog.Bind()
	return mysql.DB.Create(&smsLog).Error
}

func (*SmsService) CheckSend(mobile string, content string, ip string) error {
	if len(mobile) == 0 {
		return errors.New("手机号不能为空")
	}
	if len(content) == 0 {
		return errors.New("发送内容不能为空")
	}
	if len(content) >= 1024 {
		return errors.New("发送内容超出长度")
	}
	return nil
}
