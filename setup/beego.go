package setup

import (
	"alphawolf.com/alphaservice/filters"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/mysql"
)

func InitBeego() {
	beego.BConfig.AppName = "AlphaService"
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Listen.Graceful = true

	beego.BConfig.WebConfig.StaticDir["static"] = "static"

	runMode := beego.AppConfig.String("runmode")
	if len(runMode) < 1 {
		beego.BConfig.RunMode = "prod"
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.BConfig.ServerName = "AlphaService"

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = true
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 0
	beego.BConfig.WebConfig.Session.SessionAutoSetCookie = true
	//beego.BConfig.WebConfig.Session.SessionProvider = "memory"
	beego.BConfig.WebConfig.Session.SessionName = "AlphaServiceSessionId"
	//set session
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.String("orm::connStr")

	//globalSessions, _ = session.NewManager("mysql")
	//go globalSessions.GC()

	beego.InsertFilter("/*", beego.FinishRouter, filters.FilterResLogger, false)
	beego.InsertFilter("/*", beego.BeforeRouter, filters.FilterReqLogger)
	beego.InsertFilter("/*", beego.BeforeRouter, filters.FilterRequestID)
	beego.InsertFilter("/*", beego.BeforeRouter, filters.FilterUserAuth)

}
