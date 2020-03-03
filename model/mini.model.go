package model

import (
	"encoding/json"
	"github.com/hi-sb/io-tail/common/constants"
	"github.com/hi-sb/io-tail/core/base"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/syserr"
)

// 小程序模型
type MiniModel struct {

	db.BaseModel

	// 小程序图标
	MiniLogo string

	// 小程序名称
	MiniName string

	// 小程序地址
	MiniAddress string

	// 游戏介绍
	MiniDesc string

	//备注
	MiniRemark string

	// 状态 1:启用 0：停用
	MiniStatus  int `gorm:"type:int(2);not null;default:1"`

	// 排序
	MiniSort int `gorm:"type:int(2);not null;default:0"`
}


// 参数验证
func (*MiniModel) checkCreate(m *MiniModel) error {
	if m.MiniAddress == "" {
		return syserr.NewParameterError("小程序地址不能为空")
	}
	if m.MiniName == "" {
		return syserr.NewParameterError("小程序名称不能为空")
	}
	if m.MiniLogo == "" {
		return syserr.NewParameterError("小程序Logo不能为空")
	}
	return nil
}

// 持久化并加入缓存
func (m *MiniModel) CreateAndJoinCache() error {
	// 参数验证
	err := m.checkCreate(m)
	if  err != nil {
		return err
	}
	m.ID = ""
	m.Bind()
	err = mysql.DB.Create(m).Error
	if err != nil {
		return err
	}
	// 加入缓存
	m.saveOrUpdateCache(m)
	return nil
}

// 创建或者更新缓存
func (*MiniModel) saveOrUpdateCache(m *MiniModel){
	data, err := json.Marshal(m)
	if err == nil {
		cache.RedisClient.HSet(constants.MINI_PROGRAM_HKEY, m.ID, data)
	}
}

// 根据id获取MiniModel
func (m *MiniModel) FindByMiniId(ID string) (*MiniModel,error){
	mini := new(MiniModel)
	jsonData,err := cache.RedisClient.HGet(constants.MINI_PROGRAM_HKEY,ID).Result()
	if jsonData !="" && err == nil {
		err = json.Unmarshal([]byte(jsonData), mini)
		return mini,nil
	}
	if err !=nil {
		// 查询DB
		err = mysql.DB.Model(mini).First(mini).Error
		if err != nil {
			return nil,err
		}else {
			m.saveOrUpdateCache(mini)
		}
		return mini,nil
	}
	return mini,nil
}

// 更新并刷新缓存
func (m *MiniModel) UpdateAndJoinCache() error {
	// 验证参数
	err := m.checkCreate(m)
	if err != nil {
		return err
	}
	err = mysql.DB.Save(m).Error
	if err != nil {
		return err
	}
	m.saveOrUpdateCache(m)
	return nil
}

// 从cache and db 移除
func (*MiniModel) RemoveByMiniId(id string) error {
	cache.RedisClient.HDel(constants.MINI_PROGRAM_HKEY, id)
	return mysql.DB.Where("id=?",id).Delete(&MiniModel{}).Error
}

// 分页查询 and 条件查询

func (*MiniModel) FindOptionsPage(page base.Pager) (*base.Pager, error) {
	var miniArray []MiniModel

	// 查询
	err := mysql.DB.
		Limit(page.GetLimit()).
		Offset(page.GetOffset()).
		Order("created_at desc").
		Order("mini_sort desc").
		Find(&miniArray).Error
	var total int64 = 0
	err = mysql.DB.Model(&MiniModel{}).Count(&total).Error
	if err != nil {
		return nil, err
	}
	page.Total = total
	page.Body = miniArray
	return &page, nil
}