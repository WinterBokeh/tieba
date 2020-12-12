package Service

import (
	"Tieba/Dao"
	"Tieba/Model"
	"Tieba/Tool"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserService struct {

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

//更改用户状态
func (u *UserService) ChangeUserState(state int, name string) error {
	ud := Dao.UserDao{Tool.DbEngine}
	return ud.ChangeState(state, name)
}

//获取用户状态
func (u *UserService) GetUserState(name string) (int, error) {
	ud := Dao.UserDao{Tool.DbEngine}
	userinfo, err := ud.QueryByName(name)
	if err != nil {
		return -1, err
	}

	return userinfo.State, nil
}

func (u *UserService) Register(userinfo Model.Userinfo) error {
	ud := Dao.UserDao{Tool.DbEngine}
	err := ud.Register(userinfo)
	return err
}

//发送验证码
func (u *UserService) SendCode(email string) (string, error) {
	emailCfg := Tool.GetCfg().Email
	auth := smtp.PlainAuth("", emailCfg.ServiceEmail, emailCfg.ServicePwd, emailCfg.SmtpHost)
	to := []string{email}

	rand.Seed(time.Now().Unix())
	code := rand.Intn(10000)
	str := fmt.Sprintf("From:%v\r\nTo:%v\r\nSubject:tieba注册验证码\r\n\r\n您的验证码为：%d\r\n请在10分钟内完成验证", emailCfg.ServiceEmail, email, code)
	msg := []byte(str)
	err := smtp.SendMail(emailCfg.SmtpHost + ":" + emailCfg.SmtpPort, auth, emailCfg.ServiceEmail, to, msg)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(code), nil
}
