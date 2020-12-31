package controllers

import (
	"alphawolf.com/alpha/util"
	"crypto/tls"
	"time"

	"alphawolf.com/alphaservice/impls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// PwdMgrController operations for PwdMgr
type PwdMgrController struct {
	beego.Controller
}

// Get pwd for option id ...
// @Title Get pwd for option id
// @Description Get pwd for option id
// @Param	id		query 	string	true		"The pwd_option id"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /getPwd [get]
func (c *PwdMgrController) GetPwd() {
	id, _ := c.GetInt("id", -1)
	if id < 1 {
		c.Data["json"] = GetErrMsg("invalid id parameter")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	//roleId := c.Ctx.Input.Session("rid").(int)

	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}
	req := httplib.Post(masterApiAddr + "/cluster/pwdMgr/getPwd")

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "AlphaService",
		ConnectTimeout:   15 * time.Second,
		ReadWriteTimeout: 30 * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)

	body := make(map[string]interface{})
	body["UserId"] = userId
	body["OptionId"] = id
	bb, _ := util.JsonIter.Marshal(body)

	req.Body(bb)

	res, code, err := impls.CommonSrv.HttpReq(req)
	//logs.Info("res,code,err:",string(res),code,err.Error())
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(code)
	c.Ctx.Output.Body([]byte(res))
	c.ServeJSON()
	return

}
