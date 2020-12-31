package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

// DatabaseAccessController operations for database access
type DatabaseAccessController struct {
	beego.Controller
}

type ExecuteOutput struct {
	Status  int
	Result  []map[string]interface{}
	Columns []string
	RawOut  string
	Cost    int64
}

type SqlResMsg struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []*ExecuteOutput `json:"data,omitempty"`
}

// Execute sql ...
// @Title Execute sql
// @Description Execute sql
// @Param	targetId		query 	string	true		"targetId"
// @Param	optionId		query 	string	true		"optionId"
// @Param	uuid            query 	string	true		"uuid"
// @Param	body		body 	string	true		"sql"
// @Success 200 connected
// @Failure 403 disconnected
// @router /sql [post]
func (c *DatabaseAccessController) ExecuteSql() {
	c.EnableRender = false

	targetId, err := c.GetInt("targetId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid targetId"
		c.ServeJSON()
		return
	}

	target, err := models.GetTargetById(targetId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid target"
		c.ServeJSON()
		return
	}

	optionId, err := c.GetInt("optionId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid optionId"
		c.ServeJSON()
		return
	}

	option, err := models.GetTargetOptionById(optionId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid option"
		c.ServeJSON()
		return
	}

	uuid := c.GetString("uuid")

	userId := c.Ctx.Input.Session("uid").(int)

	user, err := models.GetUserById(userId)
	if err != nil {
		logs.Error("invalid login user, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid login user"
		c.ServeJSON()
		return
	}

	userGroupIds := make([]int, 0)

	o := orm.NewOrm()
	_, err = o.Raw("select group_id from user_owner where user_id = ?", user.Id).QueryRows(&userGroupIds)
	if err != nil && err != orm.ErrNoRows {
		c.Ctx.Output.SetStatus(403)
		errMsg := "invalid login user group information"
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		c.ServeJSON()
		return
	}

	//logs.Trace("optionId:",optionId," targetId:",targetId," userId:",userId," userGroupIds:",userGroupIds)

	ok, err := impls.OperateAclSrv.ConnectCheck(option.Id, target.Id, target.GroupId, user.Id, userGroupIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		errMsg := "connect check failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	} else if ok {
		c.Ctx.Output.SetStatus(403)
		errMsg := "deny connect"
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	sql := string(c.Ctx.Input.RequestBody)

	ok, err = impls.OperateAclSrv.SqlInputCheck(sql, option.Id, target.Id, target.GroupId, user.Id, userGroupIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		errMsg := err.Error()
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	} else if !ok {
		c.Ctx.Output.SetStatus(403)
		errMsg := "拒绝非法sql执行"
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	impls.OperateAuditSrv.Add(uuid, 3, user, target, option, c.Ctx.Request.RemoteAddr, sql)

	item := &models.Item{
		AutoIgnore:     0,
		CategoryId:     2,
		CmdType:        1,
		Command:        sql,
		ConnTimeout:    5,
		ConnectOptions: option.ConnectOptions,
		Dbname:         option.Dbname,
		ExecTimeout:    15,
		Host:           target.Address,
		ItemTypeId:     option.ItemTypeId,
		ParseMode:      "array",
		Password:       option.Password,
		Port:           option.Port,
		ServiceId:      option.ServiceId,
		Username:       option.Username,
	}

	itemBytes, _ := util.JsonIter.Marshal(item)

	masterApiAddr, err := impls.GlobalConfig.GetMasterApiAddr()
	if err != nil {
		c.Data["json"] = err.Error()
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

	req.Body(itemBytes)
	res, code, err := impls.CommonSrv.HttpReq(req)
	if err != nil {
		c.Data["json"] = err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if code != 200 {
		c.Data["json"] = string(res)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	logs.Trace("sql res:", res)

	sqlRes := SqlResMsg{}
	err = util.JsonIter.Unmarshal([]byte(res), &sqlRes)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		errMsg := err.Error()
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	if !sqlRes.Success {
		c.Ctx.Output.SetStatus(403)
		errMsg := sqlRes.Message
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	if len(sqlRes.Data) == 0 {
		c.Ctx.Output.SetStatus(code)
		c.Ctx.Output.Body([]byte(res))
		c.ServeJSON()
		return
	}

	data, err := impls.OperateAclSrv.SqlOutputCheck(sqlRes.Data[0].Result, option.Id, target.Id, target.GroupId, user.Id, userGroupIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		errMsg := err.Error()
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	logs.Trace("out data:", data)

	sqlRes.Data[0].Result = data
	sqlRes.Data[0].Columns = impls.OperateAclSrv.GetCols(data)

	resBytes, _ := util.JsonIter.Marshal(sqlRes)

	impls.OperateAuditSrv.Add(uuid, 3, user, target, option, c.Ctx.Request.RemoteAddr, string(resBytes))

	c.Ctx.Output.SetStatus(code)
	c.Ctx.Output.Body([]byte(resBytes))
	c.ServeJSON()
	return

}
