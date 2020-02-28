
## IO-TAIL 接口文档

### 公共响应头

|消息头|说明|
|-|-|
|Content-Type|只支持JSON格式，application/json; charset=utf-8|

### Token 请求头
> AUTH_TOKEN

### 错误返回格式

|参数名|类型|说明|
|-|-|-|
|Code|string|错误码|
|Message|string|错误信息|
|Body|string|响应体|


# 业务文档
> 说明： 完成的API  http://IP+端口+uri

## 一、登录鉴权

### 1.1 获取短信验证码

##### URI
> POST  /verify/sms
##### 请求参数
```
{
	"MobileNumber":"181*******"
}
```
##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "Data": "MTM2ODgwNTg5OTU=",
        "Id": "17b4c230774179c3083e20d1dc42a6f0"
    },
    "Success": true
}
```

### 1.2 登录注册二合一接口
##### URI
> POST  /user/login


##### 请求参数
```
{
	"MobileNumber":"181*******",
	"VerifyCode":"4709"
}
```

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdE51bSI6IiIsImV4cCI6MTU4MzAzMDU5OCwiaWF0IjoxNTgyNDI1Nzk4LCJpZCI6IjJkMWRkYmU3OTcyZjRkNWM5NTA2M2Q0OWFkYWNhMDMxIiwidHlwZSI6InVzZXIiLCJ1c2VybmFtZSI6IjEzNjg4MDU4OTk1In0.TiU1bCtQg-b-cH0B74l_0q9bI_QOfpW3Y8h1rhNjyHs",
    "Success": true
}
```

## 二、好友业务

### 2.1 添加好友
##### URI
> POST  /friend 

> 请求头 AUTH_TOKEN : TOKEN

##### 请求参数
```
{
	"FriendID":"",  // 被添加好友的用户ID
    "UtoFRemark":""  // 好友备注
}
```

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": null,
    "Success": true
}
```


### 2.2 获取当前用户添加好友的请求
##### URI
> POST  /friend/add-friend-req/items

> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": [
        {
            "FriendID": "bed8fb35819c473098d63aba5a8f71a8",
            "MobileNumber": "",
            "NickName": "",
            "Avatar": ""
        }
    ],
    "Success": true
}
```


### 2.3 更新添加好友状态（拒绝/同意）
##### URI
> PUT  /friend/update-friend-req

> 请求头 AUTH_TOKEN : TOKEN

##### 请求参数
```
{
	"ID":"b2f3819077a6459b9728d7e3b504764e",   // 当前请求数据库记录ID
	"State":1 ,   // 1：同意 0：拒绝
	"ReqId":"2d1ddbe7972f4d5c95063d49adaca031",   // 请求方用户ID 
    "FtoURemark":""    // 当前用户对请求方的备注
}
```

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": null,
    "Success": true
}
```

### 2.4 获取好友列表
##### URI
> GET  /friend

> 请求头 AUTH_TOKEN : TOKEN


##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": [
        {
            "FriendID": "bed8fb35819c473098d63aba5a8f71a8",
            "MobileNumber": "191******",
            "NickName": "191******",
            "Avatar": "",
            "Remark": "", // 备注
            "Initial":""  // 首字母

        }
    ],
    "Success": true
}
```


### 2.5 好友黑名单设置（将好友加入or移除黑名单）
##### URI
> PUT  /friend/black

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"FriendID":"b2f3819077a6459b9728d7e3b504764e",  
	"IsBlack":1    // 0 拉黑  1 正常

}
```
##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": null,
    "Success": true
}
```

### 2.6 根据手机号搜索好友（添加好友使用）
##### URI
> GET  /friend/{phone}

> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "FriendID": "2d1ddbe7972f4d5c95063d49adaca031",
        "MobileNumber": "",
        "NickName": "",
        "Avatar": "",
        "Remark": ""
    },
    "Success": true
}
```


### 2.7 删除好友
##### URI
> DELETE  /friend/{friendId}

> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": null,
    "Success": true
}
```


### 2.8 发送消息时验证是否在黑名单
##### URI
> GET  /friend/check-send-msg/{friendId}

> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": true,   // true 可以发送  false 不可以发送
    "Success": true
}
```


## 三、群组业务
### 3.1 创建群
##### URI
> POST  /group

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupName":"", // 默认名称 群聊(成员数)
	"GroupAnnouncement": "群公告",
	"GroupMembers":"1111,2222,3333"   // 成员ID（当前用户除外）
}
```
##### 响应内容
```
{
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "GroupModel": {
            "ID": "a05ce9eb171d46c8b8467e4ec354b949",
            "CreatedAt": "2020-02-27T15:09:49+08:00",
            "UpdatedAt": "2020-02-27T15:09:49+08:00",
            "GroupName": "群聊(4)",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "2d1ddbe7972f4d5c95063d49adaca031",
            "GroupChatStatus": 1    // 0:全体禁言  1:正常
        },
        "GroupMemberDetail": [
            {
                "ID": "",
                "CreatedAt": "0001-01-01T00:00:00Z",
                "UpdatedAt": "0001-01-01T00:00:00Z",
                "GroupID": "a05ce9eb171d46c8b8467e4ec354b949",
                "GroupMermerID": "",
                "GroupMermerNickName": "",   //群昵称
                "GroupMemberRole": 0,  //  0: 普通成员 1.群主  2。管理员
                "MobileNumber": "",
                "NickName": "",    // 用户昵称
                "Avatar": ""
            }
        ]
    },
    "Success": true
}
}
```

### 3.2 邀请新成员加入
##### URI
> POST  /group-member/join

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupID":"9dd4bc46ce7a44288587e04752a2bf68",
	"userID":"e4bf011da16c4eed91010df048f915c2"
}
```
##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "CurrentUser": {
            "ID": "aa88ac0cd9f94fc1adb6779ff3cf4cdf",
            "CreatedAt": "0001-01-01T00:00:00Z",
            "UpdatedAt": "0001-01-01T00:00:00Z",
            "MobileNumber": "",
            "NickName": "",
            "Avatar": "",
            "PrvKey": "",
            "PubKey": ""
        },
        "InvitationUser": {
            "ID": "e4bf011da16c4eed91010df048f915c2",
            "CreatedAt": "2020-02-28T10:34:40.8949178+08:00",
            "UpdatedAt": "2020-02-28T10:34:40.8949178+08:00",
            "MobileNumber": "",
            "NickName": "",
            "Avatar": "",
            "PrvKey": "",
            "PubKey": ""
        },
        "GroupInfo": {
            "ID": "9dd4bc46ce7a44288587e04752a2bf68",
            "CreatedAt": "2020-02-28T10:57:00+08:00",
            "UpdatedAt": "2020-02-28T10:57:00+08:00",
            "GroupName": "群聊(4)",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "aa88ac0cd9f94fc1adb6779ff3cf4cdf",
            "GroupChatStatus": 1
        },
        "Count": 4
    },
    "Success": true
}
```








