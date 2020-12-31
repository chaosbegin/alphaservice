package setup

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"syscall"
)

type Program struct {
}

func (this *Program) Start(s service.Service) error {
	go this.run()
	return nil
}

func (this *Program) run() {
	InitDebug()
	InitBeego()
	err := InitOrm()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	err = impls.OperateAuditSrv.Start()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	impls.HelpSrv.Initialize()

	beego.Run()
}

func (this *Program) Stop(s service.Service) error {
	//syscall.Kill(syscall.Getpid(),syscall.SIGTERM)
	//if this.client != nil {
	//	this.client.Stop()
	//}
	util.Raise(syscall.SIGTERM)
	return nil
}
