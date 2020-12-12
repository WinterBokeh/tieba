package Model

import "time"

type Userinfo struct {
	Uid       int64     `xorm:"pk autoincr" json:"uid"`
	Name      string    `xorm:"varchar(10) unique" json:"name"`
	Pwd       string    `xorm:"varchar(32)" json:"pwd"`
	Email     string    `xorm:"varchar(20)" json:"email"`
	RegDate   time.Time `xorm:"date" json:"reg_date"`
	Statement string    `xorm:"varchar(30)" json:"statement"`
	State     int       `xorm:"int" json:"state"`
	Power     int       `xorm:"int" json:"power"`
	Salt      string    `xorm:"varchar(10)" json:"salt"`
}
