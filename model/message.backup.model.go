package model

import (
	"github.com/hi-sb/io-tail/core/base"
	"github.com/hi-sb/io-tail/core/db"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/syserr"
)

var (
	friendModel = new(FriendModel)
	groupModel  = new(GroupModel)
)

//消息记录
type MessageBackup struct {
	db.BaseModel
	// form user id
	FormId string
	//to user id
	ToId string
	//昵称
	NickName string
	// 头像
	Avatar string
	// send time
	SendTime int64
	// message body
	Body string
	// message type
	ContentType string
}

//异步入库，断线可能导致消息备份失败
//但是该备份丢失不影响消息，消息不会丢失
func (*MessageBackup) AsyncSave(backup MessageBackup) {
	go func() {
		backup.Bind()
		mysql.DB.Create(backup)
	}()
}

//分页查询接口
// 分页查询 and 条件查询
//my id 我的id
// fId 好友id
func (*MessageBackup) PrivateMessagePage(page base.Pager, myId string, fId string, sendTime int64) (*base.Pager, error) {
	//验证好友关系
	isFriend := friendModel.CheckRelationship(myId, fId)
	if !isFriend {
		return nil, syserr.NewBaseErr("您还不是对方的好友")
	}
	if fId == "" {
		return nil, syserr.NewBaseErr("缺少必要的参数 fId")
	}
	var MessageBackupModels []MessageBackup
	// 查询
	err := mysql.DB.
		Limit(page.GetLimit()).
		Where("(form_id =? and to_id =?) or (form_id =? and to_id =?) and send_time <?",
			myId, fId,
			fId, myId,
			sendTime).
		Offset(page.GetOffset()).
		Order("send_time desc").
		Find(&MessageBackupModels).Error
	var total int64 = 0
	err = mysql.DB.Model(&MessageBackup{}).
		Where("(form_id =? and to_id =?) or (form_id =? and to_id =?) and send_time <?",
			myId, fId,
			fId, myId,
			sendTime).
		Count(&total).Error
	if err != nil {
		return nil, err
	}
	page.Total = total
	page.Body = MessageBackupModels
	return &page, nil
}

//分页查询群消息接口
// 分页查询 and 条件查询
//my id 我的id
// fId 好友id
func (*MessageBackup) GroupMessagePage(page base.Pager, myId string, groupId string, sendTime int64) (*base.Pager, error) {
	if groupId == "" {
		return nil, syserr.NewBaseErr("缺少必要的参数 groupId")
	}
	isGroupLife := groupModel.CheckGroupLife(groupId)
	if !isGroupLife {
		return nil, syserr.NewBaseErr("当前群不可用")
	}
	var MessageBackupModels []MessageBackup
	// 查询
	err := mysql.DB.
		Limit(page.GetLimit()).
		Where("to_id=? and send_time <?",
			groupId,
			sendTime).
		Offset(page.GetOffset()).
		Order("send_time desc").
		Find(&MessageBackupModels).Error
	var total int64 = 0
	err = mysql.DB.Model(&MessageBackup{}).
		Where("to_id=? and send_time <?",
		groupId,
		sendTime).Count(&total).Error
	if err != nil {
		return nil, err
	}
	page.Total = total
	page.Body = MessageBackupModels
	return &page, nil
}
