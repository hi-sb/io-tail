package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hi-sb/io-tail/common/constants"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/log"
	"strings"
)


var groupMemberMolde = new(GroupMemberModel)

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
	// 0:普通用户  1 后端管理人员
	UserRole int `gorm:"type:int(2);not null;default:0"`
	// 公私钥
	PrvKey string
	PubKey string

}

// user briefly
// 用户简要，提供开放的访问
type UserBriefly struct {
	//昵称
	NickName string
	// 头像
	Avatar string
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

// 根据id获取用户信息
func (*UserModel) GetInfoById(ID string) *UserModel{
	user := new(UserModel)
	// 从redis获取
	result ,err := cache.RedisClient.HGet(constants.USER_BASE_INFO_REDIS_KEY,fmt.Sprintf(constants.USER_BASE_INFO_REDIS_PREFIX,ID)).Result()
	if err == nil && result != "" {
		err := json.Unmarshal([]byte(result), user)
		if err != nil {
			fmt.Println(err)
		}
		return user
	}
	err = mysql.DB.Where("id =?", ID).First(user).Error
	if err != nil {
		return nil
	}
	return user
}
// 根据ids获取用户信息
func (*UserModel) GetInfoByIds(ids *[]string)*[]UserModel{
	var users []UserModel
	// TODO 此处有可能有问题 有待进一步确认
	idArrayStr := strings.Replace(strings.Trim(fmt.Sprint(*ids), "[]"), " ", ",", -1)
	err := mysql.DB.Where("id in (?)", idArrayStr).Find(&users).Error
	if err != nil {
		return nil
	}
	return &users
}
// 根据手机号查询用户信息
func (*UserModel) GetInfoByPhone(phone string) *UserModel {
	user := new(UserModel)
	err := mysql.DB.Where("mobile_number =?", phone).First(user).Error
	if err != nil {
		return nil
	}
	return user
}

// 修改操作后刷新用户缓存
func (*UserModel) refushCache(ID string) {
	user := new(UserModel)
	err := mysql.DB.Where("id =?", ID).First(user).Error
	if err != nil {
		log.Log.Error(err)
	}
	// 缓存用户信息
	data,err := json.Marshal(user)
	if err == nil {
		_,err = cache.RedisClient.HSet(constants.USER_BASE_INFO_REDIS_KEY,fmt.Sprintf(constants.USER_BASE_INFO_REDIS_PREFIX,ID),data).Result()
		if err !=nil {
			println(err)
		}
	}

	// 刷新group-member缓存
	groupMemberMolde.RefreshCacheByMember(ID)
}
