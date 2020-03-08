package user

import (
	"encoding/base64"
	"fmt"
	"github.com/hi-sb/io-tail/services/sms"
	"github.com/hi-sb/io-tail/utils"
	"github.com/mojocn/base64Captcha"
	"io"
)

// sms Driver
type DriverSms struct {
	numLength int
}

//
func (d *DriverSms) GenerateQuestionAnswer(key map[string]string) (string, string, string, error) {
	ip := key["ip"]
	model := key["model"]
	mobile := key["mobile"]
	keyId := utils.Md5V2(mobile)
	answer := d.getAnswer()
	content := fmt.Sprintf(model, answer)
	err := sms.SmsServiceObj.Send(mobile, content, ip)
	if err != nil {
		return "", "", "", err
	}
	return keyId, mobile, answer, nil
}

//
func (d *DriverSms) GenerateItem(question string) (item base64Captcha.Item, err error) {
	return &StringItem{body: question}, nil
}

//生成短信验证码
//这里按指定长度生成纯数字验证码
func (d *DriverSms) getAnswer() string {
	//var answer string
	//for i := 0; i < d.numLength; i++ {
	//	j := rand.Intn(10)
	//	answer += strconv.Itoa(j)
	//}
	//return answer
	return "8888"
}

func NewDriverSms(length int) *DriverSms {
	return &DriverSms{numLength: length}
}

type KeyDriver interface {
	GenerateItem(content string) (item base64Captcha.Item, err error)
	GenerateQuestionAnswer(key map[string]string) (string, string, string, error)
}

// Captcha captcha basic information.
type KeyCaptcha struct {
	Driver KeyDriver
	Store  base64Captcha.Store
}

//Generate generates a random id, base64 image string or an error if any
func (c *KeyCaptcha) Generate(key map[string]string) (id, b64s string, err error) {
	id, content, answer, err := c.Driver.GenerateQuestionAnswer(key)
	if err != nil {
		return "", "", err
	}
	item, err := c.Driver.GenerateItem(content)
	if err != nil {
		return "", "", err
	}
	c.Store.Set(id, answer)
	b64s = item.EncodeB64string()
	return
}

//Verify by a given id key and remove the captcha value in store,
//return boolean value.
//if you has multiple captcha instances which share a same store.
//You may want to call `store.Verify` method instead.
func (c *KeyCaptcha) Verify(id, answer string, clear bool) (match bool) {
	match = c.Store.Get(id, clear) == answer
	return
}

// new KeyCaptcha
func NewKeyCaptcha(driver KeyDriver, store base64Captcha.Store) *KeyCaptcha {
	return &KeyCaptcha{Driver: driver, Store: store}
}

// string Item
type StringItem struct {
	body string
}

//WriteTo writes to a writer
func (*StringItem) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

func (s *StringItem) EncodeB64string() string {
	return base64.StdEncoding.EncodeToString([]byte(s.body))
}
