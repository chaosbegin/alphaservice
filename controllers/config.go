package controllers

import (
	"alphawolf.com/alpha/util"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// ConfigController operations for Config
type ConfigController struct {
	beego.Controller
}

type ConfigQueryMsg struct {
	TargetId    int    `json:"target_id,omitempty"`
	TargetIp    string `json:"target_ip,omitempty"`
	QueryString string `json:"query_string,omitempty"`
}

// Query target config data ...
// @Title Query target config data ...
// @Description Query target config data
// @Param	body		body 	controllers.ConfigQueryMsg	true		"body for query content"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /query [post]
func (c *ConfigController) Query() {
	var msg ConfigQueryMsg
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		errMsg := "Invalid body parameter, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	data := ""
	if msg.TargetId > 0 {
		if len(msg.QueryString) > 0 {
			err = o.Raw("select JSON_EXTRACT(data,'"+msg.QueryString+"') from config_data where id = ?", msg.TargetId).QueryRow(&data)
			if err != nil && err != orm.ErrNoRows {
				errMsg := "Get config data by id:" + strconv.Itoa(msg.TargetId) + " failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

		} else {
			err = o.Raw("select data from config_data where id = ?", msg.TargetId).QueryRow(&data)
			if err != nil && err != orm.ErrNoRows {
				errMsg := "Get config data by id:" + strconv.Itoa(msg.TargetId) + " failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}

	} else if len(msg.TargetIp) > 0 {
		if len(msg.QueryString) > 0 {
			err = o.Raw("select JSON_EXTRACT(a.data,'"+msg.QueryString+"') from config_data a,target b where a.target_id = b.id and b.address = ? limit 1", msg.TargetIp).QueryRow(&data)
			if err != nil && err != orm.ErrNoRows {
				errMsg := "Get config data by ip:" + msg.TargetIp + " failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		} else {
			err = o.Raw("select a.data from config_data a,target b where a.target_id = b.id and b.address = ? limit 1", msg.TargetIp).QueryRow(&data)
			if err != nil && err != orm.ErrNoRows {
				errMsg := "Get config data by ip:" + msg.TargetIp + " failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}

	} else {
		errMsg := "Please provide at least one of target_id or target_ip."
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetDataMsg(data)
	c.Ctx.Output.SetStatus(200)
	c.ServeJSON()
	return

}
