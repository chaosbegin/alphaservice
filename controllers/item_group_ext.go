package controllers

import (
	"strconv"
	"strings"

	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// DeleteByIds ...
// @Title Delete by ids
// @Description Delete by ids
// @Param	ids	query	string	true	"ids. e.g. 1,2 ..."
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /delete [get]
func (c *ItemGroupController) DeleteByIds() {
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

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	if roleId > 2 {
		o.Begin()
		_, err := o.Raw("delete from item_tmpl where group_id in ("+idsStr+") and user_id = ?", userId).Exec()
		if err != nil {
			o.Rollback()
			errMsg := "update item_tmpl failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return

		}

		_, err = o.Raw("delete from item_group where id in ("+idsStr+") and user_id = ?", userId).Exec()
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
			errMsg := "delete item_group commit failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	} else if roleId == 1 {
		o.Begin()
		_, err := o.Raw("delete from item_tmpl where group_id in (" + idsStr + ")").Exec()
		if err != nil {
			o.Rollback()
			errMsg := "update item_tmpl failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return

		}

		_, err = o.Raw("delete from item_group where id in (" + idsStr + ")").Exec()
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
			errMsg := "delete item_group commit failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	} else {
		c.Data["json"] = GetErrMsg("Insufficient permissions")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
}

// Set item group pid ...
// @Title Set item group pid
// @Description Set item group pid
// @Param	id	query	string	true	"item group id"
// @Param	pid	query	string	true	"item group pid"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /setPid [get]
func (c *ItemGroupController) SetPid() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = GetErrMsg("Invalid item id, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if id < 1 {
		c.Data["json"] = GetErrMsg("Invalid item id")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	pid, err := c.GetInt("pid")
	if err != nil {
		c.Data["json"] = GetErrMsg("Invalid item group pid, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	group, err := models.GetItemGroupById(id)

	if err != nil || group == nil {
		c.Data["json"] = GetErrMsg("Invalid item id, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if roleId > 2 {
		if group.UserId != userId {
			c.Data["json"] = GetErrMsg("Insufficient permissions")
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	group.Pid = pid

	err = models.UpdateItemGroupById(group)

	if err != nil {
		c.Data["json"] = GetErrMsg("Update item pid failed, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
	return
}

// Get group item count ...
// @Title Get group item count
// @Description Get group item count
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /itemCount [get]
func (c *ItemGroupController) GetGroupItemCount() {
	o := orm.NewOrm()
	type itemCount struct {
		Id     int
		UserId int
		Num    int
	}

	cs := make([]*itemCount, 0)

	_, err := o.Raw("select a.id,a.user_id,count(b.id) as num from item_group a left join item_tmpl b on a.id = b.group_id group by a.id,a.user_id having count(b.id) > 0").QueryRows(&cs)
	if err != nil {
		c.Data["json"] = GetErrMsg("Get item count failed, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetDataMsg(cs)
	c.ServeJSON()
	return

}
