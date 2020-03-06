package source

import (
	"github.com/hi-sb/io-tail/core/auth"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/model"
)

var (
	friendModel      = new(model.FriendModel)
	groupModel       = new(model.GroupModel)
	groupMemberModel = new(model.GroupMemberModel)
)

//
type PermissionService struct {
}

//检查写入权限
//对于个人消息和群消息来说，其实就是一个是私有资源
//而另一个是开放资源，写入一个私有资源，也就是给一个固定的人发送消息，那么需要验证发送者是否是被发送者的好友
//而开放资源则是写入一个群消息到话题，那么这个时候就需要验证该发送者是否加入了该群
func (*PermissionService) CheckWritePermission(jwt *auth.JWT, name string) error {
	isFriend := friendModel.CheckRelationship(name, jwt.ID)
	if !isFriend {
		return syserr.NewBaseErr("您还不是对方的好友")
	}
	isFriendBlack := friendModel.CheckFriendBlack(name, jwt.ID)
	if isFriendBlack {
		return syserr.NewBaseErr("您已经被对方拉黑")
	}
	return nil
}

//检查写入权限
//对于个人消息和群消息来说，其实就是一个是私有资源
//而另一个是开放资源，写入一个私有资源，也就是给一个固定的人发送消息，那么需要验证发送者是否是被发送者的好友
//而开放资源则是写入一个群消息到话题，那么这个时候就需要验证该发送者是否加入了该群
func (*PermissionService) CheckGroupWritePermission(jwt *auth.JWT, name string) error {
	isGroupLife := groupModel.CheckGroupLife(name)
	if !isGroupLife {
		return syserr.NewBaseErr("当前群不可用")
	}
	isMemberIsTalk := groupMemberModel.CheckMemberIsTalk(name, jwt.ID)
	if isMemberIsTalk {
		return syserr.NewBaseErr("您已经被禁言")
	}
	return nil
}
