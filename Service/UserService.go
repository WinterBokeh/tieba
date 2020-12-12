package Service

import (
	"Tieba/Dao"
	"Tieba/Model"
	"Tieba/Tool"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserService struct {

}

//更改用户状态
func (u *UserService) ChangeUserState(state int, name string) error {
	ud := Dao.UserDao{Tool.DbEngine}
	return ud.ChangeState(state, name)
}

//获取用户状态
func (u *UserService) GetUserState(name string) (int, error) {
	ud := Dao.UserDao{Tool.DbEngine}
	return ud.FindState(name)
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
