
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
    "Message": "OK",
    "Code": 200,
    "Body": {
        "GroupModel": {
            "ID": "be43b195a3bb4eb5abe73f246f8d9c47",
            "CreatedAt": "2020-02-28T15:34:47+08:00",
            "UpdatedAt": "2020-02-28T15:34:47+08:00",
            "GroupName": "群聊(3)",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 1
        },
        "GroupMemberDetail": [
            {
                "ID": "3a9566c0d1b94b0e9120ec64334fd042",
                "CreatedAt": "2020-02-28T15:34:47+08:00",
                "UpdatedAt": "2020-02-28T15:34:47+08:00",
                "GroupID": "be43b195a3bb4eb5abe73f246f8d9c47",
                "GroupMermerID": "23a463f1e3f5459ea252d3817682f2d9",
                "GroupMermerNickName": "",
                "GroupMemberRole": 0,
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            },
            {
                "ID": "61d3c9bf8e1a4944a67ec76242aea139",
                "CreatedAt": "2020-02-28T15:34:47+08:00",
                "UpdatedAt": "2020-02-28T15:34:47+08:00",
                "GroupID": "be43b195a3bb4eb5abe73f246f8d9c47",
                "GroupMermerID": "ef0966125d9a475b8b02ef1b732298f1",
                "GroupMermerNickName": "",
                "GroupMemberRole": 0,
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            },
            {
                "ID": "95db3b530b56462ea849ced1fd1cfd34",
                "CreatedAt": "2020-02-28T15:34:47+08:00",
                "UpdatedAt": "2020-02-28T15:34:47+08:00",
                "GroupID": "be43b195a3bb4eb5abe73f246f8d9c47",
                "GroupMermerID": "e52781a030724f9080e88f0847caf400",
                "GroupMermerNickName": "",
                "GroupMemberRole": 0,
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            }
        ]
    },
    "Success": true
}
```

### 3.2 更新群公告
##### URI
> PUT /group/global/notice

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupName":"", // 默认名称 群聊(成员数)
	"GroupAnnouncement": "群公告",
	"GroupMembers":"1111,2222,3333"   // 成员ID（当前用户除外）
}
```

```
{
	"ID":"be43b195a3bb4eb5abe73f246f8d9c47",
	"GroupAnnouncement":"测试更新群公告1333311"
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

### 3.3 获取群基本信息 以及成员列表
##### URI
> GET /group/{groupID}

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "GroupModel": {
            "ID": "be43b195a3bb4eb5abe73f246f8d9c47",
            "CreatedAt": "2020-02-28T15:34:47+08:00",
            "UpdatedAt": "2020-02-28T15:34:47+08:00",
            "GroupName": "群聊(4)",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 1
        },
        "GroupMemberDetail": [
            {
                "ID": "6f5fdaf3e8e2490ab6235c6a2997969a",
                "CreatedAt": "2020-02-29T13:50:11+08:00",
                "UpdatedAt": "2020-02-29T13:50:11+08:00",
                "GroupID": "15be653fb4b64bb1941739c7b8673e5a",
                "GroupMemberID": "23a463f1e3f5459ea252d3817682f2d9",
                "GroupMemberNickName": "测试昵称",
                "GroupMemberRole": 2,
                "IsForbidden": 1,   // 是否禁言  0: 正常发言 1:禁言
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            },
            {
                "ID": "b550170d084a4fad9a096f853db6e2e6",
                "CreatedAt": "2020-02-29T13:50:11+08:00",
                "UpdatedAt": "2020-02-29T13:50:11+08:00",
                "GroupID": "15be653fb4b64bb1941739c7b8673e5a",
                "GroupMemberID": "ef0966125d9a475b8b02ef1b732298f1",
                "GroupMemberNickName": "",
                "GroupMemberRole": 0,
                "IsForbidden": 0,
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            }
        ]
    },
    "Success": true
}
```

### 3.4 邀请新成员加入
##### URI
> POST  /group-member/join

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupID":"a405dffc9bcd4d76a244ddbc66810662",
	"UserID":"e4bf011da16c4eed91010df048f915c2,1be2fc22d1a64e29bbcdeaf99747ba2c"   // 多个用逗号隔开
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
            "Avatar": ""
        },
        "InvitationUserArray": [
            {
                "ID": "1be2fc22d1a64e29bbcdeaf99747ba2c",
                "CreatedAt": "2020-02-27T17:23:22.2107817+08:00",
                "UpdatedAt": "2020-02-27T17:23:22.2107817+08:00",
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            },
            {
                "ID": "e4bf011da16c4eed91010df048f915c2",
                "CreatedAt": "2020-02-28T10:34:40.8949178+08:00",
                "UpdatedAt": "2020-02-28T10:34:40.8949178+08:00",
                "MobileNumber": "",
                "NickName": "",
                "Avatar": ""
            }
        ],
        "GroupInfo": {
            "ID": "a405dffc9bcd4d76a244ddbc66810662",
            "CreatedAt": "2020-02-28T11:23:51+08:00",
            "UpdatedAt": "2020-02-28T11:23:51+08:00",
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

### 3.5 剔除成员
##### URI
> DELETE  /group-member/remove

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupID":"a405dffc9bcd4d76a244ddbc66810662",
	"userID":"e4bf011da16c4eed91010df048f915c2,1be2fc22d1a64e29bbcdeaf99747ba2c"   // 多个用逗号隔开
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


### 3.6 更新禁言状态
##### URI
> PUT  /group/global/forbidden/words

> 请求头 AUTH_TOKEN : TOKEN

```
{
	
    "ID":"be43b195a3bb4eb5abe73f246f8d9c47",
    "GroupChatStatus":0   // 1：正常  0：全体禁言
    
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


### 3.7 设置管理员
##### URI
> PUT  /group-member/admin

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupID":"15be653fb4b64bb1941739c7b8673e5a",  
	"GroupMemberID":"23a463f1e3f5459ea252d3817682f2d9",
	"GroupMemberRole":2
	
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



### 3.7 设置昵称
##### URI
> PUT  /group-member/nick-name

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupID":"15be653fb4b64bb1941739c7b8673e5a",
	"GroupMemberID":"23a463f1e3f5459ea252d3817682f2d9",
	"GroupMemberNickName":"测试昵称"
	
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



### 3.8 退出群聊
##### URI
> DELETE  /group-member//sign-out

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupMemberID":"23a463f1e3f5459ea252d3817682f2d9",
	"GroupID":"15be653fb4b64bb1941739c7b8673e5a"
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











