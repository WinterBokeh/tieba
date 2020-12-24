package Controller

import (
	"Tieba/Model"
	"Tieba/Param"
	"Tieba/Service"
	"Tieba/Tool"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type UserController struct {

}

func (u *UserController) Router(engine *gin.Engine) {
	engine.POST("/api/user/register", u.register)
	engine.POST("/api/user/login", u.login)
	engine.POST("/api/verify/email", u.sendEmailVerifyCode)
	engine.GET("/api/verify/token", u.getToken)
	engine.PUT("/api/user/statement", u.changeStatement)
	engine.PUT("/api/user/email", u.changeEmail)
}

//更改邮箱, 前端需提前调用对应api发送验证码
func (u *UserController) changeEmail(ctx *gin.Context) {
	token := ctx.PostForm("token")
	newEmail := ctx.PostForm("newEmail")
	verifyCode := ctx.PostForm("verifyCode")

	//解析token
	us := Service.UserService{}
	model, err := us.ParseToken(token)
	if err != nil {
		if err.Error()[:16] == "token is expired" {
			Tool.Failed(ctx, "token失效")
			return
		}
		fmt.Println("changeStatementParseTokenErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	username := model.Name
	email := model.Email
	//验证邮箱验证码
	flag, err := us.VerifyJudgeCodeFromRedis(ctx, email, verifyCode)
	if err != nil {
		fmt.Println("changeEmailVerifyJudgeCodeFromRedisErr", err)
		Tool.Failed(ctx, "系统错误")
		return
	}
	if flag == false {
		Tool.Failed(ctx, "验证码错误")
		return
	}

	err = us.ChangeEmail(username, newEmail)
	if err != nil {
		fmt.Println("ChangeEmailServiceErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	Tool.Success(ctx, "更改邮箱成功！")
}

//更新个性签名
func (u *UserController) changeStatement(ctx *gin.Context) {
	statement := ctx.PostForm("statement")
	token := ctx.PostForm("token")

	//解析token
	us := Service.UserService{}
	model, err := us.ParseToken(token)
	if err != nil {
		if err.Error()[:16] == "token is expired" {
			Tool.Failed(ctx, "token失效")
			return
		}
		fmt.Println("changeStatementParseTokenErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	username := model.Name
	err = us.ChangeStatement(statement, username)
	if err != nil {
		fmt.Println("ChangeStatementErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	Tool.Success(ctx, "修改个性签名成功")
}

//通过refreshToken刷新token
func (u *UserController) getToken(ctx *gin.Context) {
	refreshToken := ctx.Query("refreshToken")

	us := Service.UserService{}

	//判断refreshToken状态
	model, err := us.ParseToken(refreshToken)
	if err != nil {
		if err.Error()[:16] == "token is expired" {
			Tool.Failed(ctx, "refreshToken失效")
			return
		}

		fmt.Println("getTokenParseTokenErr:", err)
		Tool.Failed(ctx, "refreshToken不正确或系统错误")
		return

	}

	//创建新token
	newToken, err := us.CreateToken(model.Name, model.Email, 120)
	if err != nil {
		fmt.Println("getTokenCreateErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	Tool.Success(ctx, newToken)
}

//发送验证码到邮箱
func (u *UserController) sendEmailVerifyCode(ctx *gin.Context) {
	email := ctx.PostForm("email")

	us := Service.UserService{}
	_, err := us.SendCode(ctx, email)
	if err != nil {
		fmt.Println("SendCodeErr: ", err)
		Tool.Success(ctx, "发送失败")
		return
	}

	Tool.Success(ctx, "发送成功，有效期10分钟")
}

//登录
func (u *UserController) login(ctx *gin.Context) {
	name := ctx.PostForm("username")
	pwd := ctx.PostForm("password")

	us := Service.UserService{}

	//愉快登录
	flag, err := us.Login(name, pwd)
	if err != nil {
		fmt.Println("loginErr:", err)
		Tool.Failed(ctx, "服务器错误")
		return
	}

	if flag == false {
		Tool.Failed(ctx, "用户名或密码错误")
		return
	}

	//创建token, 有效期两分钟
	email, err := us.GetEmailByName(name)
	if err != nil {
		fmt.Println("GetEmailByNameErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	tokenString, err := us.CreateToken(name, email, 120)
	if err != nil {
		fmt.Println("CreateTokenErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	//创建refresh token, 有效期一周
	refreshToken, err := us.CreateToken(name, email, 604800)
	if err != nil {
		fmt.Println("CreateRefreshTokenErr:", err)
		Tool.Failed(ctx, "系统错误")
		return
	}

	ctx.JSON(200, gin.H{
		"status": "0",
		"data": "登录成功",
		"token": tokenString,
		"refreshToken": refreshToken,
	})

	//logHere

	//info, err := us.ParseToken(tokenString)
	//fmt.Println(info.Name)

}

//注册服务
func (u *UserController) register(ctx *gin.Context) {
	//获取并解析用户表单
	var userParam Param.UserParam
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		Tool.Failed(ctx, "参数解析失败")
		return
	}

	//检验用户名格式
	if len(userParam.Username) < 1 {
		Tool.Failed(ctx, "用户名至少两位")
		return
	}

	//检验密码格式
	if len(userParam.Pwd) < 8 {
		Tool.Failed(ctx, "密码必须大于8位")
		return
	}

	//检查验证码
	us := Service.UserService{}
	flag, err := us.VerifyJudgeCodeFromRedis(ctx, userParam.Email, userParam.VerifyCode)
	if err != nil {
		if err.Error() == "redis: nil" {
			Tool.Failed(ctx, "验证码未发送或已失效")
			return
		}
		fmt.Println("VerifyJudgeCodeFromRedisErr:", err)
		Tool.Failed(ctx, "校验验证码失败")
		return
	}

	if flag == false {
		Tool.Failed(ctx, "验证码错误")
		return
	}

	//数据放入实体，并插入数据库
	var user Model.Userinfo
	user.RegDate = time.Now()
	user.Email = userParam.Email
	user.Name = userParam.Username
	user.Salt = strconv.FormatInt(time.Now().Unix(), 10)
// 撒盐
	m5 := md5.New()
	m5.Write([]byte(userParam.Pwd))
	m5.Write([]byte(user.Salt))
	st := m5.Sum(nil)
	user.Pwd = hex.EncodeToString(st)

	err = us.Register(user)
	if err != nil {
		if err.Error()[:10] == "Error 1062" {
			Tool.Failed(ctx, "用户名重复")
			return
		}
		fmt.Println("registerErr:", err.Error())
		Tool.Failed(ctx, "系统错误")
		return
	}

	ctx.JSON(200, gin.H{
		"status": "0",
		"data": "注册成功",
	})
}
