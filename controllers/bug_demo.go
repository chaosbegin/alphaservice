package controllers

import (
	"alphawolf.com/alphaservice/impls"
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"time"
)

// BugDemoController operations for bug demo
type BugDemoController struct {
	beego.Controller
}

// Execute internal request ...
// @Title Execute internal request
// @Description Execute internal request
// @Param	conn_timeout		query 	string	false		"connect timeout for execute item"
// @Param	exec_timeout		query 	string	false		"execute timeout for execute item"
// @Param	body		body 	string	true		"body for internal parameter"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 error string
// @router /run [post]
func (c *BugDemoController) Run() {
	conn_timeout, err := c.GetInt("conn_timeout")
	if err != nil {
		conn_timeout = 20
	}

	if conn_timeout < 1 {
		conn_timeout = 20
	} else {
		conn_timeout += 5
	}

	exec_timeout, err := c.GetInt("exec_timeout")
	if err != nil {
		exec_timeout = 65
	}

	if conn_timeout < 1 {
		exec_timeout = 65
	} else {
		exec_timeout += 5
	}

	req := httplib.Post(beego.AppConfig.String("bug_demo::url"))

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "alphaservice",
		ConnectTimeout:   time.Duration(conn_timeout) * time.Second,
		ReadWriteTimeout: time.Duration(exec_timeout) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)

	req.Body(c.Ctx.Input.RequestBody)
	res, code, err := impls.CommonSrv.HttpReq(req)
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
