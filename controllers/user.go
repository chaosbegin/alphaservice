package controllers

import (
	"alphawolf.com/alpha/util"
	"strconv"
	"strings"

	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"

	"github.com/astaxie/beego"
)

// UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId != 1 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Insufficient permissions"
		c.ServeJSON()
		return
	}

	var v models.User
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		pwd, err := impls.PwdDecrypt(v.Password)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "Invalid password"
			c.ServeJSON()
			return
		}

		v.Password = impls.CommonSrv.PwdHash(pwd)

		if _, err := models.AddUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
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
	return
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	roleId := c.Ctx.Input.Session("rid").(int)
	userId := c.Ctx.Input.Session("uid").(int)

	if roleId != 1 {
		if id < 1 || userId != id {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "Insufficient permissions"
			c.ServeJSON()
			return
		}

	}

	v, err := models.GetUserById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	} else {
		v.Password = ""
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = v
		}

	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
	var fields []string
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId == 1 {
		fields = []string{"Id", "RoleId", "Name", "Status", "NoticeStatus", "LoginName", "Mobile", "Openid", "Email"}
	} else {
		fields = []string{"Id", "Name"}
	}

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

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
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
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	roleId := c.Ctx.Input.Session("rid").(int)
	userId := c.Ctx.Input.Session("uid").(int)
	if roleId != 1 {
		if userId < 1 || userId != id {
			c.Data["json"] = GetErrMsg("Insufficient permissions")
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	v := models.User{Id: id}
	if err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		tUser, err := models.GetUserById(v.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = err.Error()
		} else {
			v.Password = tUser.Password
			if err := models.UpdateUserById(&v); err == nil {
			} else {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = err.Error()
			}
		}

	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId != 1 {
		c.Data["json"] = GetErrMsg("Insufficient permissions")
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUser(id); err == nil {
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
