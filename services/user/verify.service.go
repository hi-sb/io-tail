package user

import (
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/lock"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/utils"
	"strconv"
	"time"
)

var (
	//注册短信 model
	registerSmsModel = "您的注册验证码是:%s,请不要把验证码泄漏给其他人,如非本人请勿操作。"
	//验证码超时时间
	verifyCodeTimeOut = time.Second * 300
	// verify service
	verifyService = new(VerifyService)
	// 短信验证码
	driverSms = NewDriverSms(4)
	//mathImgCaptcha
	smsCaptcha = NewKeyCaptcha(driverSms, verifyService)
)


// verify
type VerifyService struct {
}



func (*VerifyService) Set(id string, value string) {
	fmt.Println(id, value)
	cache.RedisClient.Set(id, value, verifyCodeTimeOut)
}

func (*VerifyService) Get(id string, clear bool) string {
	stringCmd := cache.RedisClient.Get(id)
	value, _ := stringCmd.Result()
	if clear {
		//不耽误返回
		cache.RedisClient.Del(id)
	}
	return value
}

func (this *VerifyService) Verify(id, answer string, clear bool) bool {
	value := this.Get(id, clear)
	return value == answer
}

func (this *VerifyService) getSmsVerify(request *restful.Request, response *restful.Response) {
	verify, err := func() (*VerifyModel, error) {
		smsVerify := new(SmsVerify)
		err := request.ReadEntity(smsVerify)
		if err != nil {
			return nil, err
		}
		err = smsVerify.Check()
		if err != nil {
			return nil, err
		}
		//手机号
		mobile := smsVerify.MobileNumber
		//加锁
		// 避免并发发送
		sync := lock.GetSync("verify:" + mobile)
		err = sync.Lock()
		if err != nil {
			return nil, err
		}
		defer sync.Unlock()
		remoteHost := utils.GetHost(request.Request.RemoteAddr)
		ipCheckKey := fmt.Sprintf("sms_check_%s", remoteHost)
		//查看 限制 key 是否存在
		value, _ := cache.RedisClient.Get(ipCheckKey).Result()
		var sendNum int64
		if value != "" {
			sendNum, _ = strconv.ParseInt(value, 10, 64)
		}
		var sendErrTimeMemo = "您发送得太快了，请稍后再试"
		if sendNum < 50 {
			//更新 限制
			cache.RedisClient.Set(ipCheckKey, sendNum+1, time.Minute)
		} else {
			// 一分钟 内如果发送次数超过 50次
			// 那么直接 封ip 6小时
			cache.RedisClient.Set(ipCheckKey, sendNum+1, time.Hour*6)
			sendErrTimeMemo = "触发预警，该IP 将被限制发送6小时,请不要继续尝试"
		}
		if sendNum > 1 {
			return nil, errors.New(sendErrTimeMemo)
		}
		id := utils.Md5V2(mobile)
		v := smsCaptcha.Store.Get(id, false)
		if len(v) > 0 {
			return nil, errors.New("发送太过频繁请稍后再试")
		}
		keyMap := map[string]string{"mobile": mobile, "model": registerSmsModel, "ip": remoteHost}
		id, base64, err := smsCaptcha.Generate(keyMap)
		if err != nil {
			return nil, err
		}
		return &VerifyModel{Id: id, Data: base64}, nil
	}()
	rest.WriteEntity(verify, err, response)
}


func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/verify")
	webService.Route(webService.POST("/sms").To(verifyService.getSmsVerify))
	binder.BindAdd()
}
