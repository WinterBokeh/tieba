# 接口文档

## RESTful API

| 请求类型 |      请求路径       |          功能描述           |
| :------: | :-----------------: | :-------------------------: |
|   POST   |  /api/verify/email  |       发送邮箱验证码        |
|   POST   | /api/user/register  |          用户注册           |
|   POST   |   /api/user/login   |          用户登录           |
|   GET    |  /api/verify/token  | 使用refreshToken获取新token |
|   PUT    | /api/user/statement |        更新个性签名         |
|   PUT    |   /api/user/email   |        更新用户邮箱         |

## 一些约定

- 默认情况下返回的json字段为status和data
- status为状态码，"0"为成功, "1"为失败，类型为string
- data为返回值，无特殊情况下可以直接输出到浏览器，对于一些前端需要处理错误的情况，会单独给出对应的data值
- 在token权限操作中，token失效会统一返回status: "1", data: "token失效" 

## 更新用户邮箱

### 输入参数

``application/x-www-form-urlencoded``

| 字段       | 说明       |
| ---------- | ---------- |
| token      | token      |
| newEmail   | 新邮箱     |
| verifyCode | 邮箱验证码 |

## 更新个性签名

### 输入参数

``application/x-www-form-urlencoded``

| 字段      | 说明     |
| --------- | -------- |
| statement | 个性签名 |
| token     | token    |

### 返回示例

```json
{
    "status": "1",
    "data": "token失效"
}
```

## 发送邮箱验证码

### 输入参数

``application/x-www-form-urlencoded``

| 字段  | 说明   |
| ----- | ------ |
| email | 邮箱名 |

## 用户注册

### 输入参数

``application/x-www-form-urlencoded``

| 字段       | 说明           |
| ---------- | -------------- |
| username   | 用户名         |
| pwd        | 密码           |
| email      | 邮箱           |
| verifyCode | 验证码(string) |

## 用户登录

### 输入参数

| 字段     | 说明   |
| -------- | ------ |
| username | 用户名 |
| password | 密码   |

### 返回参数

| 字段         | 说明         |
| ------------ | ------------ |
| status       | 状态码       |
| data         | 返回消息     |
| token        | 用户token    |
| refreshToken | refreshToken |

### 返回示例

```json
{
    "data": "登录成功",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdGVyIiwiZXhwIjoxNjA5MzgyMzIzfQ.o41B5f-CChwEez6b61J2Ca92FDBrXq74j4rQpxHuTpY",
    "status": "0",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdGVyIiwiZXhwIjoxNjA4Nzc3NjQzfQ.1iKjv1ydzM9a0KTo5OEJZI3FBLec3Ry9wDNjYquK_E0"
}
```

## 使用refreshToken获取新token

### 方法

GET

### 输入参数

| 字段         | 说明         |
| ------------ | ------------ |
| refreshToken | refreshToken |
| username     | 用户名       |

### 返回参数

| 字段   | 说明                                                         |
| ------ | ------------------------------------------------------------ |
| data   | 成功则为新的token，若refreshToken失效则为 "refreshToken失效" |
| status | 状态码                                                       |

