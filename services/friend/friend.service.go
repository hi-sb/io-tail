package friend

import (
	"fmt"
	"github.com/axgle/pinyin"
	"github.com/emicklei/go-restful"
	"github.com/hi-sb/io-tail/core/cache"
	"github.com/hi-sb/io-tail/core/db/mysql"
	"github.com/hi-sb/io-tail/core/rest"
	"github.com/hi-sb/io-tail/core/syserr"
	"github.com/hi-sb/io-tail/services/user"
	"github.com/hi-sb/io-tail/utils"
	"sort"
)

type FriendService struct {
}

//地址
var friendService = new(FriendService)
var userModelService = new(user.UserModel)

const (
	// 好友列表redisKey前缀
	FRIEND_REDIS_PREFIX = "IO_TAIL_FRIEND_%s"
	// 好友黑名单（发送消息给某个好友查询是否被拉黑）
	FRIEND_BLACK_REDIS_PREFIX = "IO_TAIL_FRIEND_BLACK_%s"
)


// 添加好友
func (*FriendService) addFriend(request *restful.Request, response *restful.Response) {
	err := func() error {
		userId := utils.Strval(request.Attribute("currentUserId"))
		// 获取添加好友请求
		friendModel := new(FriendModel)
		err := request.ReadEntity(friendModel)
		if err != nil {
			return err
		}
		friendModel.ID = ""
		friendModel.UserID = userId
		friendModel.Bind()
		// 添加好友状态设置 (发送添加请求)
		friendModel.IsAgree = WAITING_AGREE // 设置状态为等待对方确认
		friendModel.IsBlack = IS_NOT_BLACK // 正常

		// 参数验证
		err = friendModel.Check()
		if err != nil {
			return err
		}

		// 验证是否重复添加
		var total = 0
		err = mysql.DB.Model(&FriendModel{}).Where("(user_id =? and friend_id = ?) or (friend_id =? and user_id = ?)",userId,friendModel.FriendID,friendModel.FriendID,userId).Count(&total).Error
		if err != nil {
			return err
		}
		if total >= 1 {
			return syserr.NewServiceError("好友请求已经发送,等待确认")
		}

		// 持久化好友请求
		err = mysql.DB.Create(friendModel).Error
		if err != nil {
			return err
		}
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 获取好友请求
func (*FriendService) getAddFriendReqList(request *restful.Request, response *restful.Response){
	frinedReqs, err := func() (*[]FriendAddReqModel,error) {
		userId := utils.Strval(request.Attribute("currentUserId"))
		// 获取好友请求列表
		var friendItem []FriendModel
		err := mysql.DB.Where("friend_id=? and is_agree=?", userId,WAITING_AGREE).Find(&friendItem).Error
		if err != nil {
			return nil, err
		}

		var friendReqs []FriendAddReqModel

		// 循环查询好友基本信息
		for _, friendModel := range friendItem {
			var friendReq FriendAddReqModel
			userInfo := userModelService.GetInfoById(friendModel.UserID)
			if userInfo != nil {
				friendReq.FriendID = userInfo.ID
				friendReq.Avatar = userInfo.Avatar
				friendReq.MobileNumber = userInfo.MobileNumber
				friendReq.NickName = userInfo.NickName
				friendReqs = append(friendReqs, friendReq)
			}
		}
		return &friendReqs,nil
	}()
	rest.WriteEntity(frinedReqs, err, response)
}

// 更新添加好友状态（拒绝/同意）
func (*FriendService) updateFriendIsAgree(request *restful.Request, response *restful.Response){
	err := func()error{
		userId := utils.Strval(request.Attribute("currentUserId"))

		// 获取添加好友请求
		addFReqModel := new(UpdateAddFReqModel)
		err := request.ReadEntity(addFReqModel)
		if err != nil {
			return err
		}

		// 更新添加状态 并持久化到MySQL And Redis
		friendModel := new(FriendModel)
		friendModel.ID = addFReqModel.ID
		friendModel.FriendID = userId
		friendModel.UserID = addFReqModel.ReqId
		friendModel.FtoURemark = addFReqModel.FtoURemark
		if addFReqModel.State == 1 {
			friendModel.IsAgree = AGREE_ADD
			// 同意添加 持久化到数据库
			err = mysql.DB.Where("id = ?", friendModel.ID).First(&FriendModel{}).Update(friendModel).First(friendModel).Error
			if err != nil {
				return err
			}
			// 互为好友 持久化到Redis
			userIdKEY := fmt.Sprintf(FRIEND_REDIS_PREFIX,userId)
			cache.RedisClient.SAdd(userIdKEY,addFReqModel.ReqId)
			friendIdKEY := fmt.Sprintf(FRIEND_REDIS_PREFIX,addFReqModel.ReqId)
			cache.RedisClient.SAdd(friendIdKEY,userId)
		}else {
			friendModel.IsAgree = NOT_AGREE_ADD
			//拒绝添加
			// 删除当前记录 以便再次发起添加请求
			err = mysql.DB.Delete(&friendModel).Error
			if err != nil {
				return err
			}
		}


		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 获取好友列表
func (this *FriendService) getFriendList(request *restful.Request, response *restful.Response){
	friendList,err := func() (*[]FriendAddReqModel,error){
		userId := utils.Strval(request.Attribute("currentUserId"))
		// 从redis获取当前用户的好友列表
		friendIDs, err := cache.RedisClient.SMembers(fmt.Sprintf(FRIEND_REDIS_PREFIX,userId)).Result()
		if err != nil {
			return nil,err
		}
		// 组装好友列表
		// friendInfos := userModelService.GetInfoByIds(&friendIDs)
		var friendReqs []FriendAddReqModel // 返回当前用户的好友列表

		for _, id := range friendIDs {
			user := userModelService.GetInfoById(id)
			var friendReq FriendAddReqModel
			friendReq.FriendID = user.ID
			friendReq.Avatar = user.Avatar
			friendReq.MobileNumber = user.MobileNumber
			friendReq.NickName = user.NickName

			// 查询昵称
			friendModel := new(FriendModel)
			err = mysql.DB.Where("(user_id =? and friend_id = ?) or (user_id =? and  friend_id= ?)",userId,user.ID,user.ID,userId).First(friendModel).Error
			if err != nil {
				return nil,err
			}
			// 组装昵称
			if user.ID == friendModel.UserID {
				friendReq.Remark = friendModel.FtoURemark
			}else if user.ID == friendModel.FriendID {
				friendReq.Remark = friendModel.UtoFRemark
			}

			// 获取字符串首字母
			if friendReq.Remark == "" {
				friendReq.Initial = string(this.checkAscII(int(pinyin.Convert(friendReq.NickName)[0])))
			}else{
				friendReq.Initial = string(this.checkAscII(int(pinyin.Convert(friendReq.Remark)[0])))
			}
			friendReqs = append(friendReqs, friendReq)
		}
		// 根据首字母排序、
		sort.Slice(friendReqs, func(i, j int) bool {
			return friendReqs[i].Initial > friendReqs[j].Initial
		})
		return &friendReqs,nil
	}()
	rest.WriteEntity(friendList, err, response)
}

// 拉黑 还原 好友
func (this *FriendService) pullBlackFriend(request *restful.Request, response *restful.Response){
	err := func() error{
		userId := utils.Strval(request.Attribute("currentUserId"))
		// 获取对好友 黑名单操作
		pullBlackModel := new(PullBlackModel)
		err := request.ReadEntity(pullBlackModel)
		if err != nil {
			return err
		}
		// 获取当前好友的关系
		friendModel := new(FriendModel)
		err = mysql.DB.Where("(user_id =? and friend_id = ?) or (user_id =? and  friend_id= ?)",userId,pullBlackModel.FriendID,pullBlackModel.FriendID,userId).First(friendModel).Error
		if err != nil {
			return err
		}
		friendModel = this.setIsBlack(friendModel,pullBlackModel.IsBlack,userId)
		// 更新isBlack状态
		err = mysql.DB.Where("id = ?", friendModel.ID).First(&FriendModel{}).Update(friendModel).First(friendModel).Error
		if err != nil {
			return err
		}
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 根据手机号搜索 未添加的好友
func (*FriendService) serchFriendByMobilePhone(request *restful.Request, response *restful.Response){
	friendInfo,err := func() (*FriendAddReqModel,error) {
		// 验证参数
		phone := request.PathParameter("phone")
		if !utils.VerifyMobileFormat(phone) {
			return nil,syserr.NewParameterError("手机号格式不正确")
		}

		userModel := userModelService.GetInfoByPhone(phone)
		if userModel == nil {
			return nil,syserr.NewServiceError("没有该用户")
		}

		var friendReq FriendAddReqModel
		friendReq.FriendID = userModel.ID
		friendReq.Avatar = userModel.Avatar
		friendReq.MobileNumber = userModel.MobileNumber
		friendReq.NickName = userModel.NickName
		return &friendReq,nil
	}()

	rest.WriteEntity(friendInfo, err, response)
}

// 删除好友关系
func (*FriendService) delFriend(request *restful.Request, response *restful.Response){
	err := func() error {
		userId := utils.Strval(request.Attribute("currentUserId"))
		// 验证参数
		friendId := request.PathParameter("friendId")
		if friendId == "" {
			return syserr.NewParameterError("参数不正确")
		}
		// 删除redis
		cache.RedisClient.SRem(fmt.Sprintf(FRIEND_REDIS_PREFIX,userId),friendId)
		cache.RedisClient.SRem(fmt.Sprintf(FRIEND_REDIS_PREFIX,friendId),userId)

		// 删除DB
		err := mysql.DB.Where("(user_id =? and friend_id = ?) or (user_id =? and  friend_id= ?)",userId,friendId,friendId,userId).Delete(&FriendModel{}).Error
		if err != nil {
			return err
		}
		return nil
	}()
	rest.WriteEntity(nil, err, response)
}

// 发送消息时候验证黑名单/
func (*FriendService) checkBlackList(request *restful.Request, response *restful.Response){
	isSend,err := func() (bool,error) {
		userId := utils.Strval(request.Attribute("currentUserId"))

		// 验证参数
		friendId := request.PathParameter("friendId")
		if friendId == "" {
			return false,syserr.NewParameterError("参数不正确")
		}

		// userID 发送消息给 friendId  验证 friendId的黑名单中是否有userId
		isMember,err := cache.RedisClient.SIsMember(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,friendId),userId).Result()
		if err != nil {
			return false,err
		}
		if !isMember {
			return true,nil
		}
		return false,nil
	}()
	rest.WriteEntity(isSend, err, response)
}




func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/friend")
	webService.Route(webService.POST("").To(friendService.addFriend))
	webService.Route(webService.GET("/add-friend-req/items").To(friendService.getAddFriendReqList))
	webService.Route(webService.GET("").To(friendService.getFriendList))
	webService.Route(webService.GET("/search/{phone}").To(friendService.serchFriendByMobilePhone))
	webService.Route(webService.GET("/check-send-msg/{friendId}").To(friendService.checkBlackList))
	webService.Route(webService.PUT("/update-friend-req").To(friendService.updateFriendIsAgree))
	webService.Route(webService.PUT("/black").To(friendService.pullBlackFriend))
	webService.Route(webService.DELETE("{friendId}").To(friendService.delFriend))
	binder.BindAdd()
}

// 设置拉黑值 并更新redis
func (*FriendService) setIsBlack(friendModel *FriendModel,status int,currentUserID string) *FriendModel{
	// 0 拉黑  1 正常
	if friendModel.UserID == currentUserID {  	// 如果当前用户是U 对F操作
		// 原始状态f拉黑u,u未拉黑f(10)  即将操作 u拉黑f  设置状态为互相拉黑(00)
		if friendModel.IsBlack == IS_BLACK_F_PULL_U && status == 0 {
			friendModel.IsBlack = IS_BLACK_EACH_OTHER
			// 将F加入U的黑名单
			cache.RedisClient.SAdd(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}
		// 原始状态 互相拉黑（00）   即将操作 u恢复对f的关系  设置状态为（10）
		if friendModel.IsBlack == IS_BLACK_EACH_OTHER && status == 1 {
			friendModel.IsBlack = IS_BLACK_F_PULL_U
			// 将F从U的黑名单的移除
			cache.RedisClient.SRem(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}

		// 原始状态 正常（11）   即将操作 u拉黑f    设置状态（01）
		if friendModel.IsBlack == IS_NOT_BLACK && status == 0 {
			friendModel.IsBlack = IS_BLACK_U_PULL_F
			// 将F加入U的黑名单
			cache.RedisClient.SAdd(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}
		// // 原始状态 （01）  即将操作 u未拉黑f  设置状态为正常（00）
		if friendModel.IsBlack == IS_BLACK_U_PULL_F && status == 1 {
			friendModel.IsBlack = IS_NOT_BLACK
			// 将F从U的黑名单的移除
			cache.RedisClient.SRem(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.FriendID)
		}

	} else if friendModel.FriendID == currentUserID {  // 如果当前用户是F 对U操作
		// 原始状态 f未拉黑u  u拉黑f （01）  即将操作  f拉黑u  设置状态为互相拉黑（00）
		if friendModel.IsBlack == IS_BLACK_U_PULL_F && status == 0 {
			friendModel.IsBlack = IS_BLACK_EACH_OTHER
			cache.RedisClient.SAdd(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}
		// 原始状态 f未拉黑u  u未拉黑f（11）  即将操作  f拉黑u  设置状态（10）
		if friendModel.IsBlack == IS_NOT_BLACK && status == 0 {
			friendModel.IsBlack = IS_BLACK_F_PULL_U
			cache.RedisClient.SAdd(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}

		// 原始状态 f拉黑u  u拉黑f（00）  即将操作  f未拉黑u  设置状态为 （01）
		if friendModel.IsBlack == IS_BLACK_EACH_OTHER && status == 1 {
			friendModel.IsBlack = IS_BLACK_U_PULL_F
			cache.RedisClient.SRem(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}

		// 原始状态 f拉黑u  u未拉黑f（10）  即将操作  f恢复拉黑u  设置状态为（11）
		if friendModel.IsBlack == IS_BLACK_F_PULL_U && status == 1 {
			friendModel.IsBlack = IS_NOT_BLACK
			cache.RedisClient.SRem(fmt.Sprintf(FRIEND_BLACK_REDIS_PREFIX,currentUserID),friendModel.UserID)
		}
	}

	return friendModel
}

// 首字母排序检查是否是a-z
func (*FriendService) checkAscII(ascValue int) int{
	if (ascValue>= 122 && ascValue <= 97) || (ascValue>= 65 && ascValue <= 90) {
		return ascValue
	}else{
		// # 的ASCII值
		return 35
	}
}
