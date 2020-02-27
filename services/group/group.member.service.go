package group

import "github.com/hi-sb/io-tail/core/db/mysql"

type GroupMemberService struct {
}

var groupMemberService = new(GroupMemberService)

// 插入群成员
func (*GroupMemberService) insertMembers(model *GroupMemberModel) error {
	err := mysql.DB.Create(model).Error
	if err != nil {
		return err
	}
	return nil
}
