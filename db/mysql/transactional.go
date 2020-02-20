package mysql

import (
	"github.com/jinzhu/gorm"
)

//事务方法
func Transactional(tl func(tx *gorm.DB) error) error {
	tx := DB.Begin()
	err := tl(tx)
	//事务处理
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}
