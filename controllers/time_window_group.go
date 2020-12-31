package controllers

import (
	"alphawolf.com/alphaservice/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// TimeWindowGroupController operations for TimeWindowGroup
type TimeWindowGroupController struct {
	beego.Controller
}

// URLMapping ...
func (c *TimeWindowGroupController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create TimeWindowGroup
// @Param	body		body 	models.TimeWindowGroup	true		"body for TimeWindowGroup content"
// @Success 200 {int} models.TimeWindowGroup
// @Failure 403 body is error message
// @router / [post]
func (c *TimeWindowGroupController) Post() {
	var v models.TimeWindowGroup
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddTimeWindowGroup(&v); err == nil {
			c.Data["json"] = v
		} else {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = err.Error()
		}
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get TimeWindowGroup by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TimeWindowGroup
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TimeWindowGroupController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTimeWindowGroupById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get TimeWindowGroup
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TimeWindowGroup
// @Failure 403
// @router / [get]
func (c *TimeWindowGroupController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1|col2|entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, "|")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1|col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, "|")
	}
	// order: desc|asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, "|")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, "|") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = "Error: invalid query key/value pair"
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllTimeWindowGroup(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the TimeWindowGroup
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.TimeWindowGroup	true		"body for TimeWindowGroup content"
// @Success 200 {object} models.TimeWindowGroup
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TimeWindowGroupController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.TimeWindowGroup{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTimeWindowGroupById(&v); err == nil {
		} else {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = err.Error()
		}
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the TimeWindowGroup
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TimeWindowGroupController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	twc := 0
	o := orm.NewOrm()
	err := o.Raw("select count(*) from time_window where group_id = ?", id).QueryRow(&twc)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	if twc > 0 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "时间窗口组不为空，请先删除组内时间窗口！"
		c.ServeJSON()
		return
	}

	if err := models.DeleteTimeWindowGroup(id); err == nil {
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
