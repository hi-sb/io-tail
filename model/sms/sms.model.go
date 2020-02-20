package model

import "github.com/hi-sb/io-tail/core/db"

//Sms model Log
type SmsLogModel struct {
	db.BaseModel
	//手机号
	MobileNumber string `gorm:"type:varchar(15);not null"`
	//内容
	Content string `gorm:"type:varchar(1024);not null"`
	// IP
	IP string `gorm:"type:varchar(32);not null"`
	// response Body
	ResponseBody string `gorm:"type:varchar(1024);not null"`
}
