package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

type MysqlConnectPool struct {
	Name string
}

var instance *MysqlConnectPool
var once sync.Once

var err error
var DB *gorm.DB

//初始化连接
func GetInstance() *MysqlConnectPool {
	once.Do(func() {
		instance = &MysqlConnectPool{}
	})
	return instance
}

/*
 * 初始化数据库连接(可在mail()适当位置调用)
 */
func (m *MysqlConnectPool) InitDataPool(addr string, showSql bool) (issucc bool) {
	// addr = addr + "&time_zone=Asia/Shanghai" //+ url.QueryEscape("Asia/Shanghai")
	// fmt.Printf(addr)
	DB, err = gorm.Open("mysql", addr)
	if err != nil {
		log.Fatal(err)
		return false
	}
	DB.LogMode(showSql)
	return true
}

/*
 *  对外获取数据库连接对象db
 */
func (m *MysqlConnectPool) GetMysqlDB() *gorm.DB {
	return DB
}

//获取当前 db 连接池
func GetDefaultDb() *gorm.DB {
	return GetInstance().GetMysqlDB()
}

