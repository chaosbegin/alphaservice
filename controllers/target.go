package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

// TargetController operations for Target
type TargetController struct {
	beego.Controller
}

// URLMapping ...
func (c *TargetController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Target
// @Param	body		body 	models.Target	true		"body for Target content"
// @Success 200 {int} models.Target
// @Failure 403 body is error message
// @router / [post]
func (c *TargetController) Post() {
	var v models.Target
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddTarget(&v); err == nil {
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
// @Description get Target by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Target
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TargetController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTargetById(id)
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
// @Description get Target
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Target
// @Failure 403
// @router / [get]
func (c *TargetController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

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

	l, err := models.GetAllTarget(userId, roleId, query, fields, sortby, order, offset, limit)
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
// @Description update the Target
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Target	true		"body for Target content"
// @Success 200 {object} models.Target
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TargetController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Target{Id: id}
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTargetById(&v); err == nil {
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
// @Description delete the Target
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TargetController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteTarget(id); err == nil {
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
