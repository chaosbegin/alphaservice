package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// RoleController operations for Role
type RoleController struct {
	beego.Controller
}

// URLMapping ...
func (c *RoleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Role
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 200 {int} models.Role
// @Failure 403 body is error message
// @router / [post]
func (c *RoleController) Post() {
	var v models.Role
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddRole(&v); err == nil {
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
// @Description get Role by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Role
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RoleController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRoleById(id)
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
// @Description get Role
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Role
// @Failure 403
// @router / [get]
func (c *RoleController) GetAll() {
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

	l, err := models.GetAllRole(query, fields, sortby, order, offset, limit)
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
// @Description update the Role
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 200 {object} models.Role
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RoleController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Role{Id: id}
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateRoleById(&v); err == nil {
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
// @Description delete the Role
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RoleController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	count := 0
	o := orm.NewOrm()
	err := o.Raw("select count(*) from user where role_id = ?", id).QueryRow(&count)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "query role count failed, " + err.Error()
		c.ServeJSON()
		return
	}

	if count > 0 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "角色还有用户在使用，无法删除"
		c.ServeJSON()
		return
	}

	_, err = o.Raw("delete from user_role_right where role_id = ?", id).Exec()
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete role rights failed, " + err.Error()
		c.ServeJSON()
		return
	}

	if err := models.DeleteRole(id); err == nil {
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
