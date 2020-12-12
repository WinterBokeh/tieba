package Dao

import (
	"Tieba/Model"
	"Tieba/Tool"
)

type UserDao struct {
	*Tool.Orm
}


//根据用户名查询表
func (d *UserDao) QueryByName(name string) (*Model.Userinfo, error) {
	userinfo := new(Model.Userinfo)
	_, err := d.Where("name=?", name).Get(userinfo)
	if err != nil {
		return nil, err
	}

	return userinfo, nil
}

//更改用户状态
func (d *UserDao) ChangeState(state int, name string) error {
	user := Model.Userinfo{State: state}
	_, err :=d.Where("name = ?", name).Update(&user)
	return err
}

//未激活注册
func (d *UserDao) Register(userinfo Model.Userinfo) error {
	_, err := d.InsertOne(&userinfo)
	return err
}
