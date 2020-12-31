package controllers

import (
	"strconv"
	"strings"

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
// @Param	hideDefault		query	string	false	"Hide default item for single item group. Must be an int"
// @Param	groupId			query	string	false	"Hide default item for single item group. Must be an int"
// @Param	hideDefault		query	string	false	"Hide default item for single item group. Must be an int"
// @Success 200 {object} models.Alert
// @Failure 403
// @router /page [get]
func (c *ItemTmplController) GetItemTmplByPaging() {

	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var pageSize int64 = 10
	var pageNo int64
	var headPage bool = false
	var hideDefault int = 0
	var groupId int = 0

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

	// hideDefault
	if v, err := c.GetInt("hideDefault"); err == nil {
		hideDefault = v
	}

	// groupId
	if v, err := c.GetInt("groupId"); err == nil {
		groupId = v
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

	page, err := models.GetItemTmplByPage(hideDefault, groupId, userId, query, fields, sortby, order, pageNo, pageSize)
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
func (c *ItemTmplController) BulkDeleteItemTmpl() {
	var ids string
	if v := c.GetString("ids"); v != "" {
		ids = strings.Replace(v, "|", ",", -1)
	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if len(ids) < 1 {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	o := orm.NewOrm()

	if roleId > 2 {
		if strings.ToLower(ids) == "all" {
			_, err := o.Raw("delete from item_tmpl where user_id = ?", userId).Exec()
			if err != nil {
				errMsg := "delete item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		} else {
			_, err := o.Raw("delete from item_tmpl where id in ("+ids+") and user_id = ?", userId).Exec()
			if err != nil {
				errMsg := "delete item_tmpl table rows failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}
	} else if roleId == 1 {
		if strings.ToLower(ids) == "all" {
			_, err := o.Raw("truncate table item_tmpl").Exec()
			if err != nil {
				errMsg := "truncate item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		} else {
			_, err := o.Raw("delete from item_tmpl where id in (" + ids + ")").Exec()
			if err != nil {
				errMsg := "delete item_tmpl rows failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
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
func (c *ItemTmplController) BulkUpdateItemTmpl() {
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

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	if roleId > 2 {
		switch mode {
		case 1: //set default
			_, err := o.Raw("update item_tmpl set is_default = 1 where group_id = ? and id in ("+ids+") and user_id = ?", id, userId).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 2: //cancel default
			_, err := o.Raw("update item_tmpl set is_default = 0 where group_id = ? and id in ("+ids+") and user_id = ?", id, userId).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 3: //move group
			_, err := o.Raw("update item_tmpl set group_id = ? where id in ("+ids+") and user_id = ?", id, userId).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 4: //copy to new group
			itemTmpls := make([]models.ItemTmpl, 0)
			_, err := o.Raw("select * from item_tmpl where id in (" + ids + ")").QueryRows(&itemTmpls)
			if err != nil {
				errMsg := "select item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

			for i := 0; i < len(itemTmpls); i++ {
				itemTmpls[i].GroupId = id
				itemTmpls[i].Id = 0
				itemTmpls[i].UserId = userId
			}

			_, err = o.InsertMulti(len(itemTmpls), &itemTmpls)
			if err != nil {
				errMsg := "insert item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

		}

	} else if roleId == 1 {
		switch mode {
		case 1: //set default
			_, err := o.Raw("update item_tmpl set is_default = 1 where group_id = ? and id in ("+ids+")", id).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 2: //cancel default
			_, err := o.Raw("update item_tmpl set is_default = 0 where group_id = ? and id in ("+ids+")", id).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 3: //move group
			_, err := o.Raw("update item_tmpl set group_id = ? where id in ("+ids+")", id).Exec()
			if err != nil {
				errMsg := "update item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		case 4: //copy to new group
			itemTmpls := make([]models.ItemTmpl, 0)
			_, err := o.Raw("select * from item_tmpl where id in (" + ids + ")").QueryRows(&itemTmpls)
			if err != nil {
				errMsg := "select item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

			for i := 0; i < len(itemTmpls); i++ {
				itemTmpls[i].GroupId = id
				itemTmpls[i].Id = 0
				itemTmpls[i].UserId = userId
			}

			_, err = o.InsertMulti(len(itemTmpls), &itemTmpls)
			if err != nil {
				errMsg := "insert item_tmpl table failed, " + err.Error()
				logs.Error(errMsg)
				c.Data["json"] = GetErrMsg(errMsg)
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

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

// Get items distinct item type..
// @Title Get items distinct item type
// @Description Get items distinct item type
// @Param	ids 	query	string	true	"query item ids,etc 1,2,3 or 1|2|3"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /itemtype [get]
func (c *ItemTmplController) ItemTypes() {
	var ids string
	if v := c.GetString("ids"); v != "" {
		ids = strings.Replace(v, "|", ",", -1)
	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	type ItemType struct {
		Id   int
		Name string
	}

	itemTypes := make([]*ItemType, 0)

	o := orm.NewOrm()
	_, err := o.Raw("select distinct(b.id) as id ,b.type_name as name from item_tmpl a,item_type b where a.item_type_id = b.id and b.id not in (select id from item_type where no_auth = 1) and a.id in (" + ids + ")").QueryRows(&itemTypes)
	if err != nil {
		if err != orm.ErrNoRows {
			c.Data["json"] = GetErrMsg("Query distinct item type failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = GetDataMsg(itemTypes)
	c.ServeJSON()
	return

}

// Get itemTpl applied targets...
// @Title Get itemTpl applied targets
// @Description Get itemTpl applied targets
// @Param	ids 	query	string	true	"query itemTpl ids,etc 1,2,3 or 1|2|3"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /appliedTarget [get]
func (c *ItemTmplController) AppliedTarget() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)
	var ids string
	if v := c.GetString("ids"); v != "" {
		ids = strings.Replace(v, "|", ",", -1)
	} else {
		c.Data["json"] = GetErrMsg("ids can not been null")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	type ItemTargets struct {
		Name    string
		Address string
		Count   int
	}

	itemTargets := make([]*ItemTargets, 0)

	o := orm.NewOrm()
	if roleId == 1 {
		_, err := o.Raw("select a.name,a.address,count(*) as count from target a,item b where a.id = b.target_id and b.tmpl_id in (" + ids + ") group by a.name,a.address").QueryRows(&itemTargets)
		if err != nil && err != orm.ErrNoRows {
			c.Data["json"] = GetErrMsg("Query applied targets failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	} else {
		_, err := o.Raw("select a.name,a.address,count(*) as count from target a,item b where a.id = b.target_id and a.group_id in ( select target_group_id from target_owner where user_group_id in ( select group_id from user_owner where user_id = ?)) and b.tmpl_id in ("+ids+") group by a.name,a.address", userId).QueryRows(&itemTargets)
		if err != nil && err != orm.ErrNoRows {
			c.Data["json"] = GetErrMsg("Query applied targets failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = GetDataMsg(itemTargets)
	c.ServeJSON()
	return

}

// Get itemTpl by category_id ...
// @Title Get itemTpl by category_id
// @Description Get itemTpl by category_id
// @Param	id 	query	string	true	"category_id"
// @Success 200 {int} []models.ItemTmpl
// @Failure 403 body is error message
// @router /getItemTplByCategoryId [get]
func (c *ItemTmplController) GetItemTplByCategoryId() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid id parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	itemTpls := make([]models.ItemTmpl, 0)
	o := orm.NewOrm()
	_, err = o.Raw("select * from item_tmpl where id in (select item_id from target_category_item  where category_id = ?) order by tmpl_type_id", id).QueryRows(&itemTpls)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = "Query itemTpl failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = itemTpls
	c.ServeJSON()
	return
}

// Check api item series exists ...
// @Title Check api item series exists ...
// @Description Check api item series exists ...
// @Param	series		query 	string	true		"The api item series"
// @Success 200 {object} models.Item
// @Failure 403 :id is empty
// @router /apiItemCheck [get]
func (c *ItemTmplController) ApiItemCheck() {
	series := c.GetString("series")
	series = strings.TrimSpace(series)
	if len(series) < 1 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid series name"
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	count := 0
	err := o.Raw("select count(*) from item_tmpl where item_type_id = 1000 and series = ?", series).QueryRow(&count)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "query item_tmpl by series failed, " + err.Error()
		c.ServeJSON()
		return
	}

	if count == 0 {
		c.Data["json"] = "not exists"
		c.ServeJSON()
		return
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "already exists"
		return
	}

	return
}

type ItemOwner struct {
	Id   int
	Name string
}

// Get itemTmpl owner ...
// @Title Get itemTmpl owner
// @Description Get itemTmpl owner
// @Param	ids 	query	string	true	"item ids, split by ','"
// @Success 200 {object} []controllers.ItemOwner
// @Failure 403
// @router /owner [get]
func (c *ItemTmplController) Owner() {
	ids := c.GetString("ids")
	if len(ids) < 1 {
		c.Data["json"] = "invalid item ids"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	ios := make([]*ItemOwner, 0)

	o := orm.NewOrm()
	_, err := o.Raw("select a.id as id,b.name as name from item_tmpl a,user b where a.user_id = b.id and a.id in (" + ids + ")").QueryRows(&ios)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = "query item owner failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = ios
	c.ServeJSON()
	return
}
