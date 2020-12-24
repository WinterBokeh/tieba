package Service

import (
	"Tieba/Dao"
	"Tieba/Model"
	"Tieba/Tool"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserService struct {

}

//更改用户邮箱服务
func (u *UserService) ChangeEmail(username, newEmail string) error {
	ud := Dao.UserDao{Tool.DbEngine}
	userinfo := Model.Userinfo{Email: newEmail}
	err := ud.Update(username, userinfo)
	return err
}

//通过用户名查邮箱
func (u *UserService) GetEmailByName(username string) (string, error) {
	ud := Dao.UserDao{Tool.DbEngine}
	userinfo, err := ud.QueryByName(username)
	return userinfo.Email, err
}

//更改个性签名
func (u *UserService) ChangeStatement(statement, username string) error {
	ud := Dao.UserDao{Tool.DbEngine}
	userinfo := Model.Userinfo{Statement: statement}
	err := ud.Update(username, userinfo)
	return err
}

//构建一个jwt，包括用户名, 邮箱
func (u *UserService) CreateToken(name string, email string, ExpireTime int64) (string, error) {
	JwtCfg := Tool.GetCfg().Jwt
	mySigningKey := []byte(JwtCfg.SigningKey)

	claims := Model.MyCustomClaims{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ExpireTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

//解析Token
func (u *UserService) ParseToken(tokenString string) (*Model.MyCustomClaims, error) {
	JwtCfg := Tool.GetCfg().Jwt
	mySigningKey := []byte(JwtCfg.SigningKey)
	token, err := jwt.ParseWithClaims(tokenString, &Model.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if clams, ok := token.Claims.(*Model.MyCustomClaims); ok && token.Valid {
		return clams, nil
	} else {
		return nil, err
	}
}

//登录
func (u *UserService) Login(name, pwd string) (bool, error) {
	ud := Dao.UserDao{Tool.DbEngine}
	userinfo, err := ud.QueryByName(name)
	if err != nil {
		return false, err
	}

	m5 := md5.New()
	m5.Write([]byte(pwd))
	m5.Write([]byte(userinfo.Salt))
	st := m5.Sum(nil)
	hashPwd := hex.EncodeToString(st)

	//密码错误
	if hashPwd != userinfo.Pwd {
		return false, nil
	}

	return true, nil
}

//注册服务
func (u *UserService) Register(userinfo Model.Userinfo) error {
	ud := Dao.UserDao{Tool.DbEngine}
	err := ud.Register(userinfo)
	return err
}

//redis检验验证码
func (u *UserService) VerifyJudgeCodeFromRedis(ctx *gin.Context, value string, inputCode string) (bool, error) {
	redisConn := Tool.GetRedisConn()
	cmd := redisConn.Get(ctx, value)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	if cmd.Val() != inputCode {
		return false, nil
	}

	return true, nil
}

//发送验证码,并放入redis中
func (u *UserService) SendCode(ctx *gin.Context, email string) (string, error) {
	emailCfg := Tool.GetCfg().Email
	auth := smtp.PlainAuth("", emailCfg.ServiceEmail, emailCfg.ServicePwd, emailCfg.SmtpHost)
	to := []string{email}

	fmt.Println("EMAIL", email)

	rand.Seed(time.Now().Unix())
	code := rand.Intn(10000)
	str := fmt.Sprintf("From:%v\r\nTo:%v\r\nSubject:tieba注册验证码\r\n\r\n您的验证码为：%d\r\n请在10分钟内完成验证", emailCfg.ServiceEmail, email, code)
	msg := []byte(str)
	err := smtp.SendMail(emailCfg.SmtpHost + ":" + emailCfg.SmtpPort, auth, emailCfg.ServiceEmail, to, msg)
	if err != nil {
		return "", err
	}


	redisConn := Tool.GetRedisConn()
	redisConn.Set(ctx, email, strconv.Itoa(code), time.Minute * 10)

	return strconv.Itoa(code), nil
}
