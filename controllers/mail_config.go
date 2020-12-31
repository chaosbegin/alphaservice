package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

// MailConfigController operations for MailConfig
type MailConfigController struct {
	beego.Controller
}

// URLMapping ...
func (c *MailConfigController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create MailConfig
// @Param	body		body 	models.MailConfig	true		"body for MailConfig content"
// @Success 200 {int} models.MailConfig
// @Failure 403 body is error message
// @router / [post]
func (c *MailConfigController) Post() {
	var v models.MailConfig
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddMailConfig(&v); err == nil {
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
// @Description get MailConfig by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MailConfig
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MailConfigController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetMailConfigById(id)
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
// @Description get MailConfig
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MailConfig
// @Failure 403
// @router / [get]
func (c *MailConfigController) GetAll() {
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

	l, err := models.GetAllMailConfig(query, fields, sortby, order, offset, limit)
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
// @Description update the MailConfig
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MailConfig	true		"body for MailConfig content"
// @Success 200 {object} models.MailConfig
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MailConfigController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.MailConfig{Id: id}
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateMailConfigById(&v); err == nil {
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
// @Description delete the MailConfig
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MailConfigController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteMailConfig(id); err == nil {
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
