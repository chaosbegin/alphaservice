package controllers

import (
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// Set target group pid ...
// @Title Set target group pid
// @Description Set target group pid
// @Param	id	query	string	true	"target group id"
// @Param	pid	query	string	true	"target group pid"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /setPid [get]
func (c *TargetGroupController) SetPid() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid target id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if id < 1 {
		c.Data["json"] = "invalid target id"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	pid, err := c.GetInt("pid")
	if err != nil {
		c.Data["json"] = "invalid target group pid, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	group, err := models.GetTargetGroupById(id)

	if err != nil || group == nil {
		c.Data["json"] = "invalid target group id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	pGroup, err := models.GetTargetGroupById(pid)
	if err != nil || pGroup == nil {
		c.Data["json"] = "invalid parent target group id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if group.Id == pid {
		c.Data["json"] = "不能将父组设置成自己！"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if len(group.ApiId) > 0 || len(pGroup.ApiId) > 0 {
		c.Data["json"] = "系统自动管理类目标组不支持此操作！"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	group.Pid = pid

	err = models.UpdateTargetGroupById(group)

	if err != nil {
		c.Data["json"] = "update target pid failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = ""
	c.ServeJSON()
	return
}

// DeleteByIds ...
// @Title Delete by ids
// @Description Delete by ids
// @Param	ids	query	string	true	"ids. e.g. 1,2 ..."
// @Success 200 {object} models.ItemGroup
// @Failure 403
// @router /delete [get]
func (c *TargetGroupController) DeleteByIds() {
	var ids []int

	// ids: 1,2,entity.3
	if v := c.GetString("ids"); v != "" {
		strArray := strings.Split(v, ",")
		for _, v := range strArray {
			i, err := strconv.Atoi(v)
			if err != nil {
				errMsg := "Invalid ids parameter: " + v
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

			if i == 1 { //deny delete the default node
				continue
			}

			ids = append(ids, i)
		}
	}

	if len(ids) < 1 {
		c.Data["json"] = GetErrMsg("parameter ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	var idsStr = ""
	for _, v := range ids {
		idsStr += "," + strconv.Itoa(v)
	}
	idsStr = idsStr[1:]

	o := orm.NewOrm()
	o.Begin()

	_, err := o.Raw("update target set group_id = 1 where group_id in (" + idsStr + ")").Exec()
	if err != nil {
		o.Rollback()
		errMsg := "update item_tmpl failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return

	}

	_, err = o.Raw("delete from target_group where id in (" + idsStr + ")").Exec()
	if err != nil {
		o.Rollback()
		errMsg := "delete item_group failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return

	}

	err = o.Commit()
	if err != nil {
		o.Rollback()
		errMsg := "delete target_group commit failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
}
