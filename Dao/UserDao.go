package Dao

import (
	"Tieba/Model"
	"Tieba/Tool"
)

type UserDao struct {
	*Tool.Orm
}

func (d *UserDao) ChangeState(state int, name string) error {
	user := Model.Userinfo{State: state}
	_, err :=d.Where("name = ?", name).Update(&user)
	return err
}

//查询用户状态
func (d *UserDao) FindState(name string) (int, error) {
	userinfo := new(Model.Userinfo)
	_, err := d.Where("name=?", name).Get(userinfo)
	if err != nil {
		return -1, err
	}

	return userinfo.State, nil
}

//未激活注册
func (d *UserDao) Register(userinfo Model.Userinfo) error {
	_, err := d.InsertOne(&userinfo)
	return err
}
