package controllers

import (
	"alphawolf.com/alpha/util"
	"bytes"
	"crypto/tls"
	"github.com/astaxie/beego/logs"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// ExecuteController operations for command execute
type ExecuteController struct {
	beego.Controller
}

// Execute item ...
// @Title Execute item
// @Description item
// @Param	body		body 	models.Item	true		"body for item"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /item [post]
func (c *ExecuteController) Item() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	item := &models.Item{}
	err := util.JsonIter.Unmarshal([]byte(c.Ctx.Input.RequestBody), item)
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	req := httplib.Post(masterApiAddr + "/cluster/execute/item")

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "AlphaService",
		ConnectTimeout:   time.Duration(item.ConnTimeout) * time.Second,
		ReadWriteTimeout: time.Duration(item.ExecTimeout) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)
	req.Param("uid", strconv.Itoa(userId))
	req.Param("rid", strconv.Itoa(roleId))

	req.Body(c.Ctx.Input.RequestBody)
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

// Execute internal request ...
// @Title Execute internal request
// @Description Execute internal request
// @Param	api		query 	string	true		"internal request api"
// @Param	conn_timeout		query 	string	false		"connect timeout for execute item"
// @Param	exec_timeout		query 	string	false		"execute timeout for execute item"
// @Param	body		body 	string	true		"body for internal parameter"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /internal [post]
func (c *ExecuteController) Internal() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

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

	api := c.GetString("api")
	//data := c.GetString("data")
	//logs.Trace("api:",api," data:",data)
	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	req := httplib.Post(masterApiAddr + api + "?uid=" + strconv.Itoa(userId) + "&rid=" + strconv.Itoa(roleId))

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "AlphaService",
		ConnectTimeout:   time.Duration(conn_timeout) * time.Second,
		ReadWriteTimeout: time.Duration(exec_timeout) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)

	req.Body(c.Ctx.Input.RequestBody)
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

// Client upload ...
// @Title Client upload
// @Description Client upload
// @Param	body		body	string		true		""
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /clientUpload [post]
func (c *ExecuteController) ClientUpload() {
	f, h, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("file", h.Filename)
	io.Copy(fileWriter, f)
	bodyWriter.Close()
	f.Close()

	header := make(map[string]string)
	header["Content-Type"] = bodyWriter.FormDataContentType()

	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	url := masterApiAddr + "/cluster/clientMgr/upload"
	req, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	out := impls.HttpReadBody(resp)

	if resp.StatusCode != 200 {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = out
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = out

	c.ServeJSON()
	return

}

// Client upload ...
// @Title Client upload
// @Description Client upload
// @Param	api		query 	string	true		"internal request api"
// @Param	body		body	string		true		""
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /upload [post]
func (c *ExecuteController) Upload() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	f, h, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("file", h.Filename)
	io.Copy(fileWriter, f)
	bodyWriter.Close()
	f.Close()

	header := make(map[string]string)
	header["Content-Type"] = bodyWriter.FormDataContentType()

	api := c.GetString("api")

	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	url := masterApiAddr + api + "?uid=" + strconv.Itoa(userId) + "&rid=" + strconv.Itoa(roleId)
	req, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	out := impls.HttpReadBody(resp)

	if resp.StatusCode != 200 {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = out
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = out

	c.ServeJSON()
	return

}
