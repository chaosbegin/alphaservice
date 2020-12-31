package controllers

import (
	"strconv"
	"strings"

	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// GetItemTmplByPaging ...
// @Title Get ItemTmpl by paging
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
func (c *ItemController) GetItemByPaging() {

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

	page, err := models.GetItemByPage(query, fields, sortby, order, pageNo, pageSize)
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

// BulkDeleteItemTmpl ...
// @Title Bulk delete itemTmpl rows
// @Description Bulk delete itemTmpl rows
// @Param	ids 	query	string	true	"delete ids,etc 1,2,3 or 1|2|3"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /delete [get]
func (c *ItemController) BulkDeleteItem() {
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

		if len(ids) > 0 {
			if strings.ToLower(ids) == "all" {
				_, err := o.Raw("delete from item where target_group_id in (" + tgid_str + ")").Exec()
				if err != nil {
					errMsg := "delete all item table failed, " + err.Error()
					logs.Error(errMsg)
					c.Data["json"] = GetErrMsg(errMsg)
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			} else {
				_, err := o.Raw("delete from item where target_group_id in (" + tgid_str + ") and id in (" + ids + ")").Exec()
				if err != nil {
					errMsg := "delete item rows failed, " + err.Error()
					logs.Error(errMsg)
					c.Data["json"] = GetErrMsg(errMsg)
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			}

		} else {
			c.Data["json"] = GetErrMsg("ids can not been null")
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	} else if roleId == 1 {
		if len(ids) > 0 {
			if strings.ToLower(ids) == "all" {
				_, err := o.Raw("truncate table item").Exec()
				if err != nil {
					errMsg := "truncate item table failed, " + err.Error()
					logs.Error(errMsg)
					c.Data["json"] = GetErrMsg(errMsg)
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			} else {
				_, err := o.Raw("delete from item where id in (" + ids + ")").Exec()
				if err != nil {
					errMsg := "delete item rows failed, " + err.Error()
					logs.Error(errMsg)
					c.Data["json"] = GetErrMsg(errMsg)
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			}

		} else {

			c.Data["json"] = GetErrMsg("ids can not been null")
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
	return

}

// BulkUpdateItemTmpl ...
// @Title Bulk update itemTmpl rows
// @Description Bulk delete itemTmpl rows
// @Param	mode 	query	string	true	"update mode,1: set is_default to target id,2:cancel is_default to target id,3:move to target id,4:copy to target id"
// @Param	id 	query	string	true	"target id"
// @Param	ids 	query	string	true	"srouce ids,etc 1,2,3"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /update [get]
func (c *ItemController) BulkUpdateItem() {
	var mode int
	var id int
	var ids string

	if v, err := c.GetInt("mode"); err == nil {
		mode = v
	} else {
		c.Data["json"] = GetErrMsg("Invalid mode parameter")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if v, err := c.GetInt("id"); err == nil {
		id = v
	} else {
		c.Data["json"] = GetErrMsg("Invalid id parameter")
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

	o := orm.NewOrm()

	switch mode {
	case 1:
		_, err := o.Raw("update item set is_default = 1 where group_id = ? and id in ("+ids+")", id).Exec()
		if err != nil {
			errMsg := "update item_tmpl table failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	case 2:
		_, err := o.Raw("update item set is_default = 0 where group_id = ? and id in ("+ids+")", id).Exec()
		if err != nil {
			errMsg := "update item table failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	case 3:
		_, err := o.Raw("update item set target_id = ? where id in ("+ids+")", id).Exec()
		if err != nil {
			errMsg := "update item table failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	case 4:
		items := make([]models.Item, 0)
		_, err := o.Raw("select * from item where id in (" + ids + ")").QueryRows(&items)
		if err != nil {
			errMsg := "select item table failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		for i := 0; i < len(items); i++ {
			items[i].TargetId = id
			items[i].Id = 0
		}

		_, err = o.InsertMulti(len(items), &items)
		if err != nil {
			errMsg := "insert item table failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
	return

}

// Get panel by item id ...
// @Title Get panel by item id
// @Description Get panel by item id
// @Param	id 	query	string	true	"item id"
// @Success 200 {object} []impls.Panel
// @Failure 403
// @router /panel [get]
func (c *ItemController) Panel() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid item id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	ps, err := impls.GetPanelByItemId(id)
	if err != nil {
		c.Data["json"] = "invalid item id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ps
		c.ServeJSON()
		return
	}

}
