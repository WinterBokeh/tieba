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
	engine.POST("/register", u.register)
	engine.GET("/code/:name", u.judgeCode)
	engine.POST("/login", u.login)
	engine.GET("/logout", u.logout)
}

func (u *UserController) logout(ctx *gin.Context) {
	value := Tool.CheckLogin(ctx)
	if value == "" {
		Tool.Failed(ctx, "未登录")
		return
	}

	redisConn := Tool.GetRedisConn()
	redisConn.Del("isLogin")
	Tool.Success(ctx, "退出登录成功")
}

//登录
func (u *UserController) login(ctx *gin.Context) {
	name := ctx.PostForm("name")
	pwd := ctx.PostForm("pwd")

	//未验证
	us := Service.UserService{}
	state, err := us.GetUserState(name)
	if err != nil {
		fmt.Println("GetUserStateErr: ", err)
		Tool.Failed(ctx, "服务器错误")
		return
	}

	if state == 0 {
		Tool.Failed(ctx, "该账号未激活或不存在")
		return
	}

	//已经处于登录状态
	isLogin := Tool.CheckLogin(ctx)
	if isLogin != "" {
		Tool.Failed(ctx, "您已经处于登录状态，无需再次登录")
		return
	}

	//愉快登录
	flag, err := us.Login(name, pwd)
	if err != nil {
		fmt.Println("loginErr:", err)
		Tool.Failed(ctx, "服务器错误")
		return
	}

	if flag == false {
		Tool.Failed(ctx, "密码错误")
		return
	}

	redisConn := Tool.GetRedisConn()
	redisConn.Set("isLogin", name, time.Hour)

	Tool.Success(ctx, "登陆成功，有效期一小时")
}

//验证验证码
func (u *UserController) judgeCode(ctx *gin.Context) {
	name := ctx.Param("name")
	code := ctx.Query("code")

	//判断是否未激活状态
	us := Service.UserService{}
	state, err := us.GetUserState(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if state != 0 {
		Tool.Failed(ctx, "您已经处于激活状态，无需再次激活")
		return
	}
	//redis取激活码
	redisConn := Tool.GetRedisConn()
	originCode := redisConn.Get(name).Val()
	if code != originCode {
		Tool.Failed(ctx, "验证码错误！")
		return
	}

	err = us.ChangeUserState(1, name)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	Tool.Success(ctx, "验证成功！")
}

//注册服务
func (u *UserController) register(ctx *gin.Context) {
	//获取并解析用户表单
	var userParam Param.UserParam
	err := Tool.Decode(ctx.Request.Body, &userParam)
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

	//发送验证码
	us := new(Service.UserService)
	code, err := us.SendCode(userParam.Email)
	if err != nil {
		fmt.Println("SendCodeErr: ", err)
		Tool.Failed(ctx, "服务器发送验证码错误, 请尝试检查邮箱是否存在")
		return
	}

	//数据放入实体，并插入数据库
	var user Model.Userinfo
	user.RegDate = time.Now()
	user.Email = userParam.Email
	user.Name = userParam.Username
//	user.Pwd = userParam.Pwd
	user.Statement = userParam.Statement
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

	Tool.Success(ctx, "注册成功，请在10分钟内使用验证码激活账户")

	//验证码放入redis，并设置有效时间
	redisConn := Tool.GetRedisConn()
	redisConn.Set(userParam.Username, code, time.Minute*10)
}
