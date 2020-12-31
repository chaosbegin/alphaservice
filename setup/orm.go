package setup

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func InitOrm() error {
	//orm setting
	maxIdleConn, err := beego.AppConfig.Int("orm::maxIdleConn")
	if err != nil {
		return errors.New("maxIdleConn parameter is missing, " + err.Error())
	}

	maxOpenConn, err := beego.AppConfig.Int("orm::maxOpenConn")
	if err != nil {
		return errors.New("maxOpenConn parameter is missing, " + err.Error())
	}

	//orm.RegisterDriver()
	dbType := beego.AppConfig.String("orm::type")
	if dbType == "" {
		dbType = "mysql"
	}

	err = orm.RegisterDataBase("default", dbType, beego.AppConfig.String("orm::connStr"), maxIdleConn, maxOpenConn)
	if err != nil {
		return errors.New("Connect to database failed, " + err.Error())
	}

	maxLifeTime, err := beego.AppConfig.Int("orm::maxLifeTime")
	if err != nil {
		return errors.New("maxLifeTime parameter is missing, " + err.Error())
	}

	defaultDb, err := orm.GetDB("default")
	if err != nil {
		return errors.New("Get default database failed, " + err.Error())
	}

	defaultDb.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)

	ormDebug, err := beego.AppConfig.Bool("orm::debug")
	if err == nil && ormDebug {
		orm.Debug = true
	}

	return nil
}
