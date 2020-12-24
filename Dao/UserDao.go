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

//更新用户字段通过用户名
func (d *UserDao) Update(username string, userinfo Model.Userinfo) error {
	_, err := d.Where("name = ?", username).Update(&userinfo)
	return err
}

//注册
func (d *UserDao) Register(userinfo Model.Userinfo) error {
	_, err := d.InsertOne(&userinfo)
	return err
}
