package controllers

import (
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// GetAlertByPaging ...
// @Title Get Alert by paging
// @Description Get Alert by paging
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	pageSize	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	pageNo		query	string	false	"Start position of page. Must be an integer"
// @Param	headPage		query	string	false	"Store page infomation to header. Must be an bool"
// @Success 200 {object} models.Alert
// @Failure 403
// @router /page [get]
func (c *AlertController) GetAlertByPaging() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var pageSize int64 = 10
	var pageNo int64
	var headPage bool = false

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, "|")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("pageSize"); err == nil {
		pageSize = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("pageNo"); err == nil {
		pageNo = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, "|")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, "|")
	}
	// headFlag
	if v, err := c.GetBool("headPage"); err == nil {
		headPage = v
	}

	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, "|") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = GetErrMsg("Error: invalid query key/value pair")
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	//logs.Trace("userId:",userId)
	//logs.Trace("roleId:",roleId)

	page, err := models.GetAlertByPage(userId, roleId, query, fields, sortby, order, pageNo, pageSize)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
	} else {
		if headPage {
			c.Ctx.Output.Header("Total", strconv.Itoa(int(page.Total)))
			c.Ctx.Output.Header("PageNo", strconv.Itoa(int(page.PageNo)))
			c.Ctx.Output.Header("PageSize", strconv.Itoa(int(page.PageSize)))
			c.Data["json"] = page.Rows
		} else {
			c.Data["json"] = GetDataMsg(page)
		}

	}

	c.ServeJSON()
	return

}

// UpdateAlertStatus ...
// @Title Update alert status
// @Description Update alert status
// @Param	uid	query	string	true	"user id"
// @Param	username	query	string	true	"user id"
// @Param	value	query	string	true	"set level or status value"
// @Param	ids 	query	string	true	"level or status ids"
// @Success 200 {object} models.Alert
// @Failure 403
// @router /update [get]
func (c *AlertController) UpdateAlertStatus() {
	var ids string
	var uid int64 = 0
	var username string
	var value int64 = 0

	// uid: 0 (default is 0)
	if v, err := c.GetInt64("uid"); err == nil {
		uid = v
	} else {
		c.Data["json"] = GetErrMsg("Invalid uid, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	// value: 0 (default is 0)
	if v, err := c.GetInt64("value"); err == nil {
		value = v
	} else {
		c.Data["json"] = GetErrMsg("Invalid value, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if v := c.GetString("username"); v != "" {
		username = v
	} else {
		c.Data["json"] = GetErrMsg("username can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if v := c.GetString("ids"); v != "" {
		ids = strings.Replace(v, "|", ",", -1)
	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	if roleId > 2 {
		o := orm.NewOrm()
		targetGroupIds := make([]int, 0)
		_, err := o.Raw("select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)", userId).QueryRows(&targetGroupIds)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target_group_ids by user_id failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		tgid_str := impls.IntArrayJoin(targetGroupIds)

		if ids == "all" {
			_, err := o.Raw("update alert set status = ?,confirm_uid = ?,confirm_username = ?,confirm_time = now() where target_group_id in ("+tgid_str+") and status = 1",
				value, uid, username).Exec()
			if err != nil {
				errMsg := "udpate alert level status failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

		} else {
			_, err := o.Raw("update alert set status = ?,confirm_uid = ?,confirm_username = ?,confirm_time = now() where target_group_id in ("+tgid_str+") and id in ("+ids+")",
				value, uid, username).Exec()
			if err != nil {
				errMsg := "udpate alert level status failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}

	} else if roleId == 1 {
		o := orm.NewOrm()
		if ids == "all" {
			_, err := o.Raw("update alert set status = ?,confirm_uid = ?,confirm_username = ?,confirm_time = now() where status = 1",
				value, uid, username).Exec()
			if err != nil {
				errMsg := "udpate alert level status failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

		} else {
			_, err := o.Raw("update alert set status = ?,confirm_uid = ?,confirm_username = ?,confirm_time = now() where id in ("+ids+")",
				value, uid, username).Exec()
			if err != nil {
				errMsg := "udpate alert level status failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}
	} else {
		RetErr(c, "Unauthorized")
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
	return

}

// BulkDeleteAlert ...
// @Title Bulk delete alert rows
// @Description Bulk delete alert rows
// @Param	status 	query	string	true	"delete alert status"
// @Param	ids 	query	string	true	"delete ids"
// @Success 200 {object} models.Alert
// @Failure 403
// @router /delete [get]
func (c *AlertController) BulkDeleteAlert() {
	// uid: 0 (default is 0)
	var status int
	if v, err := c.GetInt("status"); err == nil {
		status = v
	} else {
		c.Data["json"] = GetErrMsg("Invalid status, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	var ids string
	if v := c.GetString("ids"); v != "" {
		ids = strings.Replace(v, "|", ",", -1)
	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	o := orm.NewOrm()
	if len(ids) > 0 {
		if roleId > 2 {
			targetGroupIds := make([]int, 0)
			_, err := o.Raw("select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)", userId).QueryRows(&targetGroupIds)
			if err != nil {
				c.Data["json"] = GetErrMsg("Get target_group_ids by user_id failed, " + err.Error())
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

			tgid_str := impls.IntArrayJoin(targetGroupIds)

			if status > 0 {
				if ids == "all" {
					_, err := o.Raw("delete from alert where target_group_id in ("+tgid_str+") and status = ?", status).Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				} else {
					_, err := o.Raw("delete from alert where target_group_id in ("+tgid_str+") and id in ("+ids+") and status = ?", status).Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				}

			} else {
				if ids == "all" {
					_, err := o.Raw("delete from alert where target_group_id in (" + tgid_str + ")").Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				} else {
					_, err := o.Raw("delete from alert where target_group_id in (" + tgid_str + ") and id in (" + ids + ")").Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				}
			}

		} else if roleId == 1 {
			if status > 0 {
				if ids == "all" {
					_, err := o.Raw("delete from alert where status = ?", status).Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				} else {
					_, err := o.Raw("delete from alert where id in ("+ids+") and status = ?", status).Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				}

			} else {
				if ids == "all" {
					_, err := o.Raw("truncate table alert").Exec()
					if err != nil {
						RetErr(c, "truncate table alert failed, "+err.Error())
						return
					}
				} else {
					_, err := o.Raw("delete from alert where id in (" + ids + ")").Exec()
					if err != nil {
						RetErr(c, "delete alert by ids failed, "+err.Error())
						return
					}
				}
			}
		} else {
			RetErr(c, "Unauthorized")
			return
		}

	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
	return

}

// Get panel by alert id ...
// @Title Get panel by alert id
// @Description Get panel by alert id
// @Param	id 	query	string	true	"alert id"
// @Success 200 {object} []impls.Panel
// @Failure 403
// @router /panel [get]
func (c *AlertController) Panel() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid alert id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	ps, err := impls.GetPanelByAlertId(id)
	if err != nil {
		c.Data["json"] = "invalid alert id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ps
		c.ServeJSON()
		return
	}

}

// Get alert count by item ids ...
// @Title Get alert count by item ids
// @Description Get alert count by item ids
// @Param	ids 	query	string	true	"item ids,ex:1,2,3,4"
// @Success 200 {object} alert count array
// @Failure 403
// @router /alertCountByItems [get]
func (c *AlertController) AlertCountByItems() {
	ids := c.GetString("ids")
	if len(ids) < 1 {
		c.Data["json"] = "invalid item ids"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	type AlertCount struct {
		Id    int
		Total int
	}

	ac := make([]AlertCount, 0)

	o := orm.NewOrm()
	_, err := o.Raw("select item_id id,count(*) total from alert where status = 1 and item_id in (" + ids + ") group by item_id").QueryRows(&ac)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = "query alert count by item ids failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ac
		c.ServeJSON()
		return
	}

}
