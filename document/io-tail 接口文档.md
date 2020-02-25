
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

### 1.获取短信验证码

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

### 2.登录注册二合一接口
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


### 3. 添加好友
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


### 4. 获取当前用户添加好友的请求
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


### 5. 更新添加好友状态（拒绝/同意）
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

### 6.获取好友列表
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
            "Remark": ""
        }
    ],
    "Success": true
}
```


### 7. 好友黑名单设置（将好友加入or移除黑名单）
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







