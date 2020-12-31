package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type RespMsg struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func GetErrMsg(message string) RespMsg {
	return RespMsg{Success: false, Message: message}
}

func GetOkMsg() RespMsg {
	return RespMsg{Success: true, Message: "success"}
}

func GetDataMsg(data interface{}) RespMsg {
	return RespMsg{Success: true, Message: "success", Data: data}
}

func GetRespMsg(success bool, message string, data interface{}) RespMsg {
	return RespMsg{Success: success, Message: message, Data: data}
}

func RetErr(c interface{}, errMsg string) {
	logs.Error(errMsg)
	cc, ok := c.(*beego.Controller)
	if ok {
		cc.Data["json"] = RespMsg{Success: false, Message: errMsg}
		cc.Ctx.Output.SetStatus(403)
		cc.ServeJSON()
	}
	return
}

func RetOk(c interface{}) {
	cc, ok := c.(*beego.Controller)
	if ok {
		cc.Data["json"] = RespMsg{Success: true, Message: "success"}
		cc.ServeJSON()
	}
	return
}

func RetData(c interface{}, data interface{}) {
	cc, ok := c.(*beego.Controller)
	if ok {
		cc.Data["json"] = RespMsg{Success: true, Message: "success", Data: data}
		cc.ServeJSON()
	}
	return
}
