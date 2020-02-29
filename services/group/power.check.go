package group

import "github.com/hi-sb/io-tail/core/log"


// 验证当前用户和所在group中的角色
func CheckGroupRole(groupID string ,userID string) bool {
	groupMemberModel,err := new(GroupMemberModel).getGroupMemberByGroupIdAndMemberId(groupID,userID)
	if err != nil {
		log.Log.Error(err)
		return false
	}
	if groupMemberModel.GroupMemberRole != 0 {
		return true
	}
	return false
}
