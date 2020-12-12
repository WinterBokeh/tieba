package Tool

import (
	"Tieba/Model"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

type Orm struct {
	*xorm.Engine
}

var DbEngine *Orm

//初始化orm
func init() {
	cfg := GetCfg().Database
	source := cfg.User +":" + cfg.Password + "@tcp(" +cfg.Host + ":" + cfg.Port + ")/" +cfg.DbName +"?charset="+cfg.Charset
	engine, err :=xorm.NewEngine(cfg.Driver, source)
	if err != nil {
		fmt.Println("getEngineErr: ", err)
		panic(err)
	}

	err = engine.Sync2(new(Model.Userinfo))
	if err != nil {
		fmt.Println("syncErr: ", err)
		panic(err)
	}

	engine.ShowSQL(cfg.ShowSql)
	orm := new(Orm)
	orm.Engine = engine

	DbEngine = orm

}
