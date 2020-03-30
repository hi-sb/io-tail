
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


### 测试IP： http：//148.70.231.222:7654

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


### 1.3 更新昵称/头像
##### URI
> POST  /user/update

##### 请求参数
```
{
	"NickName":"",
	"Avatar":""
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

### 1.4 设置后台管理人员
##### URI
> POST  /admin/user/update

##### 请求参数
```
{
	"ID":"",  // 用户id
	"UserRole":""   // 0:普通用户  1 后端管理人员
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







## 二、好友业务

### 2.1 添加好友
##### URI
> POST  /friend 

> 请求头 AUTH_TOKEN : TOKEN

##### 请求参数
```
{
	"FriendID":"e52781a030724f9080e88f0847caf400",
	"UtoFRemark":"444444",
	"FrUReason":"我是***请通过下"
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
            "ID": "00a8f708f445448a88af003195acfc4a",
            "FriendID": "bc85bb92176248f8ba6566c2d1efc27a",
            "MobileNumber": "",
            "NickName": "",
            "Avatar": "",
            "Remark": "我是Test请通过下",
            "Initial": "",
            "IsBlack": 3
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
            "ID": "",
            "FriendID": "016e487a13e5441786fa316880518cfe",
            "MobileNumber": "",
            "NickName": "",
            "Avatar": "https:///header.png",
            "Remark": "测试备注A",
            "Initial": "C",
            "IsBlack": 3    //0 互相拉黑  1:被拉黑  2：拉黑好友   3：关系正常

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
> GET  /friend/search/{phone}

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

### 3.2 更新群信息（名称and公告）
##### URI
> PUT /group
> 请求头 AUTH_TOKEN : TOKEN

```
{

	"ID":"be43b195a3bb4eb5abe73f246f8d9c47",
    "GroupName":"", // 默认名称 群聊(成员数)
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
            "GroupChatStatus": 1  // 0:全体禁言  1:正常
        },
        "GroupMemberDetail": [
            {
                "ID": "6f5fdaf3e8e2490ab6235c6a2997969a",
                "CreatedAt": "2020-02-29T13:50:11+08:00",
                "UpdatedAt": "2020-02-29T13:50:11+08:00",
                "GroupID": "15be653fb4b64bb1941739c7b8673e5a",
                "GroupMemberID": "23a463f1e3f5459ea252d3817682f2d9",
                "GroupMemberNickName": "测试昵称",
                "GroupMemberRole": 2,  // 0 普通成员   1：创建者   2：管理员
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



### 3.8 设置昵称
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



### 3.9 退出群聊
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

### 3.10 根据群昵称查询成员信息
##### URI
> POST  /group-member/nick-name

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"GroupId":"e781e6e03fd64e84897295eed219d619",
	"NickName":"阿迪斯"
}
```

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "ID": "c705e28ad46d4aa1bddf3aaf0de9004f",
        "CreatedAt": "2020-02-28T17:20:29+08:00",
        "UpdatedAt": "2020-02-28T17:20:29+08:00",
        "GroupID": "e781e6e03fd64e84897295eed219d619",
        "GroupMemberID": "e52781a030724f9080e88f0847caf400",
        "GroupMemberNickName": "阿迪斯",
        "GroupMemberRole": 1,
        "IsForbidden": 0,
        "MobileNumber": "",
        "NickName": "",
        "Avatar": ""
    },
    "Success": true
}
```



### 3.11 群主解散群
##### URI
> DELETE  /group/{groupID}

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


### 3.12 群消息验证 是否可以发消息
##### URI
> DELETE  /group//check/{groupID}

> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "对不起,当前群已经被解散",
    "Code": 1091,
    "Body": false,  // false 不能发消息  true 正常群聊天
    "Success": true
}
```




### 3.13 查询当前用户已经加入的群

##### URI
> GET  /group/join/items
> 请求头 AUTH_TOKEN : TOKEN

##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": [
        {
            "ID": "be43b195a3bb4eb5abe73f246f8d9c47",
            "CreatedAt": "2020-02-28T15:34:47+08:00",
            "UpdatedAt": "2020-02-28T15:34:47+08:00",
            "GroupName": "群聊",
            "GroupAnnouncement": "测试更新群公告哈哈",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 0
        },
        {
            "ID": "4d01bd5f3723426581ff962c3c545f36",
            "CreatedAt": "2020-02-28T17:15:09+08:00",
            "UpdatedAt": "2020-02-28T17:15:09+08:00",
            "GroupName": "群聊",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 1
        },
        {
            "ID": "f6b7b75c505e4a0087dbca6e1cd4ffd3",
            "CreatedAt": "2020-02-28T17:06:21+08:00",
            "UpdatedAt": "2020-02-28T17:06:21+08:00",
            "GroupName": "群聊",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 1
        },
        {
            "ID": "e781e6e03fd64e84897295eed219d619",
            "CreatedAt": "2020-02-28T17:20:29+08:00",
            "UpdatedAt": "2020-02-28T17:20:29+08:00",
            "GroupName": "群聊",
            "GroupAnnouncement": "群公告",
            "GreateUserID": "e52781a030724f9080e88f0847caf400",
            "GroupChatStatus": 1
        }
    ],
    "Success": true
}
```











## 四、小程序

### 4.1 创建小程序（管理员）
##### URI
> POST  /admin/mini

> 请求头 AUTH_TOKEN : TOKEN


```
{
	"MiniLogo":"MiniLogo",
	"MiniName":"MiniName",
	"MiniAddress":"MiniAddress",
	"MiniDesc":"MiniDesc",
	"MiniRemark":"MiniRemark",
	"MiniStatus":1,
	"MiniSort":0
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


### 4.2 根据ID获取小程序详情
##### URI
> GET  /admin/mini/{id}  // 管理端
> GET /mini/{id} //前端

> 请求头 AUTH_TOKEN : TOKEN


##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "ID": "629bc3f1e6114523a178b126000fd323",
        "CreatedAt": "2020-03-02T12:14:34.4241998+08:00",
        "UpdatedAt": "2020-03-02T12:14:34.4241998+08:00",
        "MiniLogo": "MiniLogo",
        "MiniName": "MiniName",
        "MiniAddress": "MiniAddress",
        "MiniDesc": "MiniDesc",
        "MiniRemark": "MiniRemark",
        "MiniStatus": 1,
        "MiniSort": 0
    },
    "Success": true
} "Success": true
}
```



### 4.3 更新小程序（管理端）
##### URI
> PUT  /admin/mini 

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"MiniLogo":"MiniLogo",
	"MiniName":"MiniName",
	"MiniAddress":"MiniAddress",
	"MiniDesc":"MiniDesc",
	"MiniRemark":"MiniRemark",
	"MiniStatus":1,
	"MiniSort":0
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



### 4.4 删除小程序（管理端）
##### URI
> DELETE  /admin/mini/{id} 

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




### 4.5 分页查询小程序列表
##### URI
> POST  /admin/mini/page

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"Page": 1,
	"PageSize": 10,
	 "Body": {
	 	
	 }
}
```


##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "Page": 1,
        "PageSize": 10,
        "Total": 3,
        "Body": [
            {
                "ID": "7ccdd265b8e5441ebb5d732f7ec8263e",
                "CreatedAt": "2020-03-02T13:00:11+08:00",
                "UpdatedAt": "2020-03-02T13:00:30+08:00",
                "MiniLogo": "MiniLogo",
                "MiniName": "MiniName",
                "MiniAddress": "MiniAddress_TEST",
                "MiniDesc": "MiniDesc",
                "MiniRemark": "MiniRemark_TEST",
                "MiniStatus": 1,
                "MiniSort": 12
            },
            {
                "ID": "1e81f3b297f8432db5ba0c8c74a6f8f9",
                "CreatedAt": "2020-03-02T13:00:10+08:00",
                "UpdatedAt": "2020-03-02T13:00:57+08:00",
                "MiniLogo": "MiniL11ogo",
                "MiniName": "MiniNam1e",
                "MiniAddress": "MiniAddress",
                "MiniDesc": "MiniDesc111",
                "MiniRemark": "MiniRemark",
                "MiniStatus": 1,
                "MiniSort": 12
            },
            {
                "ID": "3da176a34e7742a8b114438a7ad9b2d2",
                "CreatedAt": "2020-03-02T13:00:09+08:00",
                "UpdatedAt": "2020-03-02T13:00:09+08:00",
                "MiniLogo": "MiniLogo",
                "MiniName": "MiniName",
                "MiniAddress": "MiniAddress",
                "MiniDesc": "MiniDesc",
                "MiniRemark": "MiniRemark",
                "MiniStatus": 1,
                "MiniSort": 0
            }
        ]
    },
    "Success": true
}
```



### 4.6 分页查询小程序列表（前端）
##### URI
> POST  /mini/page

> 请求头 AUTH_TOKEN : TOKEN

```
{
	"Page": 1,
	"PageSize": 10,
	 "Body": {
	 	
	 }
}
```


##### 响应内容
```
{
    "Message": "OK",
    "Code": 200,
    "Body": {
        "Page": 1,
        "PageSize": 10,
        "Total": 3,
        "Body": [
            {
                "ID": "7ccdd265b8e5441ebb5d732f7ec8263e",
                "CreatedAt": "2020-03-02T13:00:11+08:00",
                "UpdatedAt": "2020-03-02T13:00:30+08:00",
                "MiniLogo": "MiniLogo",
                "MiniName": "MiniName",
                "MiniAddress": "MiniAddress_TEST",
                "MiniDesc": "MiniDesc",
                "MiniRemark": "MiniRemark_TEST",
                "MiniStatus": 1,
                "MiniSort": 12
            },
            {
                "ID": "1e81f3b297f8432db5ba0c8c74a6f8f9",
                "CreatedAt": "2020-03-02T13:00:10+08:00",
                "UpdatedAt": "2020-03-02T13:00:57+08:00",
                "MiniLogo": "MiniL11ogo",
                "MiniName": "MiniNam1e",
                "MiniAddress": "MiniAddress",
                "MiniDesc": "MiniDesc111",
                "MiniRemark": "MiniRemark",
                "MiniStatus": 1,
                "MiniSort": 12
            },
            {
                "ID": "3da176a34e7742a8b114438a7ad9b2d2",
                "CreatedAt": "2020-03-02T13:00:09+08:00",
                "UpdatedAt": "2020-03-02T13:00:09+08:00",
                "MiniLogo": "MiniLogo",
                "MiniName": "MiniName",
                "MiniAddress": "MiniAddress",
                "MiniDesc": "MiniDesc",
                "MiniRemark": "MiniRemark",
                "MiniStatus": 1,
                "MiniSort": 0
            }
        ]
    },
    "Success": true
}
```



###5.0 接收消息java 列子
```
public class Main {

    public static OkHttpClient okHttpClient = new OkHttpClient.Builder().readTimeout(10, TimeUnit.SECONDS).build();

    //监听地址，可以是私有消息的也可以是群消息的地址
    private static final String source = "http://127.0.0.1:7654/topic/private";

    //token
    private static final String token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdE51bSI6IiIsImV4cCI6MTU4NDQzMTM1MywiaWF0IjoxNTgzODI2NTUzLCJpZCI6IjAxNmU0ODdhMTNlNTQ0MTc4NmZhMzE2ODgwNTE4Y2ZlIiwidHlwZSI6InVzZXIiLCJ1c2VybmFtZSI6IjEzODE2ODgwMDAzIn0._R1RuZAYs3LG-w5Ls1uy3tDOuiGtYj01j7MV8dxkvXU";

    public static void main(String[] ags) {
        new MessageListener(source,token) {
            @Override
            public void onErr(Throwable err) {
                System.out.println(err.getMessage());
            }

            @Override
            public void onMessage(String message) {
                System.out.println(message);
            }
        }.listen();
    }
}

class MessageListenerException extends Exception {

    MessageListenerException(String message) {
        super(message);
    }
}

//做一个简单的 内部类
abstract class MessageListener extends Thread {

    private String source;

    private String token;

    public MessageListener(String source, String token) {
        this.token = token;
        this.source = source;
    }

    //收到消息
    public abstract void onMessage(String message);

    //错误回调
    public abstract void onErr(Throwable err);

    //开始监听
    public void listen() {
        this.start();
    }

    @Override
    public void run() {
        Request request = new Request.Builder().get().url(source).header("AUTH_TOKEN",token).build();
        try {
            Response response = Main.okHttpClient.newCall(request).execute();
            //判断响应码
            // 如果非200 则将流中的所有数据读出作为错误描述
            if (!response.isSuccessful()) {
                String body = response.body().string();
                throw new MessageListenerException(body);
            }
            InputStreamReader inputStreamReader = new InputStreamReader(response.body().byteStream());
            BufferedReader bufferedReader = new BufferedReader(inputStreamReader);
            String message;
            //按行读取
            // 每一行就是一条消息，当没有新的消息时则阻塞
            //当然如果采用非阻塞io，那么写法上不一样，不需要使用这个线程
            while ((message = bufferedReader.readLine()) != null) {
                onMessage(message);
            }
        } catch (MessageListenerException e) {
            onErr(e);
        } catch (IOException e) {
            onErr(new RuntimeException(e.getMessage() + ":" + "网络错误，应该定时重连等操作"));
        }
    }
}

```


### 5.1 监听私有消息
##### URI

> GET /topic/private

> 请求头 AUTH_TOKEN : TOKEN

> query 参数: offset int 

> 描述：offset 为消息读取位置，当收到一条消息后，该消息内容将包含该字段的描述（也就是当前消息在话题中的位置）。
> 在进行消息监听的时候，如果传入该字段则表示从指定位置开始监听消息。如果不传入该字段则认为从昂头开始监听消息，也就是
>同时拉取所有历史消息。读取消息方式为，使用http协议访问该地址，然后按行读取响应流，每一行则是一条新的消息。如果该链接断开，
>应该重新创建链接进行监听（传入上一条消息offset，不然会拉取历史消息）。当http 状态码为200时则表示监听成功，非200则表示监听失败。监听
>失败时，应该将响应流中的所有内容读出，作为错误响应。

##### Content type ：
```

	// text message
	//文本消息，格式也是文本
	MessageTypeText string = "text/text"
	//语音消息，格式是base64
	// voice message
	MessageTypeBase64Voice string = "voice/base64"
	//图片消息，格式是base64
	// img message
	MessageTypeBase64Img string = "img/base64"
	//语音消息，格式是一个语音下载地址
	// voice message
	MessageTypeUrlVoice string = "voice/url"
	//图片消息 格式是一个图片下载地址
	// img message
	MessageTypeUrlImg string = "img/url"
	//系统通知消息，格式为文本
	// sys notify message
	MessageTypeNotify string = "notify/text"
	//添加好友消息，格式为一个json
	// Add friends
	MessageTypeAddFriends string = "add-friends/json"
    //心跳消息可以忽略
	MessageTypeHeartbeat string = "heartbeat/time-stamp"
	//被邀请 入群通知消息
	MessageTypeAddToGroup string ="add-to-group/json"
	//被踢出群通知消息
	MessageTypeExpelGroup string ="expel-group/json"

```

##### 非200 状态码错误响应格式

```

{
 "Message": "错误内容",
 "Code": 1016,
 "Body": null,
 "Success": false
}

```

##### 200 状态码正常消息格式

```
{
    	// form user id
    	FormId:"发送者的id",
    	// send time
    	SendTime:"发送时间",
    	// message body
    	Body:"消息内容",
    	// offset
    	Offset:"当前监听的话题位置",
    	// message type
    	ContentType:"消息体格式类型"
}

```

### 5.2 监听群消息
##### URI

> GET /topic/group/{source}

> 请求头 AUTH_TOKEN : TOKEN

> path 参数：source string 群id 

> query 参数: offset int 

> 描述：offset 为消息读取位置，当收到一条消息后，该消息内容将包含该字段的描述（也就是当前消息在话题中的位置）。
> 在进行消息监听的时候，如果传入该字段则表示从指定位置开始监听消息。如果不传入该字段则认为从昂头开始监听消息，也就是
>同时拉取所有历史消息。读取消息方式为，使用http协议访问该地址，然后按行读取响应流，每一行则是一条新的消息。如果该链接断开，
>应该重新创建链接进行监听（传入上一条消息offset，不然会拉取历史消息）。当http 状态码为200时则表示监听成功，非200则表示监听失败。监听
>失败时，应该将响应流中的所有内容读出，作为错误响应。

##### Content type ：
```

	// text message
	//文本消息，格式也是文本
	MessageTypeText string = "text/text"
	//语音消息，格式是base64
	// voice message
	MessageTypeBase64Voice string = "voice/base64"
	//图片消息，格式是base64
	// img message
	MessageTypeBase64Img string = "img/base64"
	//语音消息，格式是一个语音下载地址
	// voice message
	MessageTypeUrlVoice string = "voice/url"
	//图片消息 格式是一个图片下载地址
	// img message
	MessageTypeUrlImg string = "img/url"
	//系统通知消息，格式为文本
	// sys notify message
	MessageTypeNotify string = "notify/text"
	//添加好友消息，格式为一个json
	// Add friends
	MessageTypeAddFriends string = "add-friends/json"
    //心跳消息可以忽略
	MessageTypeHeartbeat string = "heartbeat/time-stamp"
	//被邀请 入群通知消息
	MessageTypeAddToGroup string ="add-to-group/json"
	//被踢出群通知消息
	MessageTypeExpelGroup string ="expel-group/json"
```


##### 非200 状态码错误响应格式

```

{
 "Message": "错误内容",
 "Code": 1016,
 "Body": null,
 "Success": false
}

```
##### 200 状态码正常消息格式
```
{
    	// form user id
    	FormId:"发送者的id",
        //昵称
        NickName: "",
        // 头像
        Avatar :"string",
    	// send time
    	SendTime:"发送时间",
    	// message body
    	Body:"消息内容",
    	// offset
    	Offset:"当前监听的话题位置",
    	// message type
    	ContentType:"消息体格式类型"
}

```

### 5.3 发送好友消息

##### URI

> PUT /topic/private/{source}

> 请求头 AUTH_TOKEN : TOKEN

> path 参数：source string  好友的id

##### Content type ：
```

	// text message
	//文本消息，格式也是文本
	MessageTypeText string = "text/text"
	//语音消息，格式是base64
	// voice message
	MessageTypeBase64Voice string = "voice/base64"
	//图片消息，格式是base64
	// img message
	MessageTypeBase64Img string = "img/base64"
	//语音消息，格式是一个语音下载地址
	// voice message
	MessageTypeUrlVoice string = "voice/url"
	//图片消息 格式是一个图片下载地址
	// img message
	MessageTypeUrlImg string = "img/url"
	//系统通知消息，格式为文本
	// sys notify message
	MessageTypeNotify string = "notify/text"
	//添加好友消息，格式为一个json
	// Add friends
	MessageTypeAddFriends string = "add-friends/json"
    //心跳消息可以忽略
	MessageTypeHeartbeat string = "heartbeat/time-stamp"
	//被邀请 入群通知消息
	MessageTypeAddToGroup string ="add-to-group/json"
	//被踢出群通知消息
	MessageTypeExpelGroup string ="expel-group/json"
```

##### 请求体

```
{
    Body:"消息内容"，
    ContentType:"内容格式"
}
```

#####  响应体 

```
{
    "Message": "OK",
    "Code": 200,
    "Success": true
}

```



### 5.4 发送群消息
##### URI

> PUT /topic/group/{source}

> 请求头 AUTH_TOKEN : TOKEN

> path 参数：source string 群id 

##### Content type ：
```

	// text message
	//文本消息，格式也是文本
	MessageTypeText string = "text/text"
	//语音消息，格式是base64
	// voice message
	MessageTypeBase64Voice string = "voice/base64"
	//图片消息，格式是base64
	// img message
	MessageTypeBase64Img string = "img/base64"
	//语音消息，格式是一个语音下载地址
	// voice message
	MessageTypeUrlVoice string = "voice/url"
	//图片消息 格式是一个图片下载地址
	// img message
	MessageTypeUrlImg string = "img/url"
	//系统通知消息，格式为文本
	// sys notify message
	MessageTypeNotify string = "notify/text"
	//添加好友消息，格式为一个json
	// Add friends
	MessageTypeAddFriends string = "add-friends/json"
	//heartbeat 心跳消息
    //心跳消息可以忽略
	MessageTypeHeartbeat string = "heartbeat/time-stamp"
	//被邀请 入群通知消息
	MessageTypeAddToGroup string ="add-to-group/json"
	//被踢出群通知消息
	MessageTypeExpelGroup string ="expel-group/json"
```

##### 请求体

```
{
    Body:"消息内容"，
    ContentType:"内容格式"
}
```

#####  响应体 

```
{
    "Message": "OK",
    "Code": 200,
    "Success": true
}

```

### 5.5 获取自己在指定话题的监听位置，也就是offset
##### URI

> GET /topic/offset/{source}

> 请求头 AUTH_TOKEN : TOKEN

> path 参数：source string 群id 或者自己的id（可通过解析token获得）

##### 响应体：
```
{
    "Message": "OK",
    "Code": 200,
    "Body":167616,//offset
    "Success": true
}
```

### 6.1 按用户id获取用户简要信息（昵称、头像）
##### URI

> GET /user/briefly/{id}

> 请求头 AUTH_TOKEN : TOKEN

> path 参数：id string 用户id

##### 响应体：
```
{
 "Message": "OK",
 "Code": 200,
 "Body": {
  "NickName": "13816880001",
  "Avatar": ""
 },
 "Success": true
}
```

### 7.1 对象存储接口

> 测试地址:148.70.231.222:6543

> POST /project/upload/{suffix}
>
> path 参数：suffix string 文件后缀 类似于jpg
>
> from body 参数 uploadFile 类型为文件类型

##### 响应体：
```
{
    "Message": "OK",
    "Code": 200,
    //地址
    "Body": "http://148.70.231.222:6543/20200312/8/4/3/a/854b33ade0824f0aa50d523ea60f68fc.jpg",
    "Success": true
}
```

