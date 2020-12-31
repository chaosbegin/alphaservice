package controllers

import (
	"alphawolf.com/alpha/util"
	"github.com/astaxie/beego"
)

type DemoController struct {
	beego.Controller
}

// Demo http object ...
// @Title Demo http object
// @Description Demo http object
// @Param	body		body 	string	true		"body for demo json object"
// @Success 201 {int} models.Alert
// @Failure 403 body is empty
// @router /http/object [post]
func (c *DemoController) Object() {
	req := make(map[string]interface{})
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
	} else {
		c.Data["json"] = req
	}

	c.ServeJSON()
	return
}

// Demo http array ...
// @Title Demo http array
// @Description Demo http array
// @Param	body		body 	string	true		"body for demo json object"
// @Success 201 {int} models.Alert
// @Failure 403 body is empty
// @router /http/array [post]
func (c *DemoController) Array() {
	req := make([]map[string]interface{}, 0)
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		req = append(req, map[string]interface{}{"success": false, "message": err.Error(), "data": ""})
		c.Data["json"] = req
		c.Ctx.Output.SetStatus(403)
	} else {
		c.Data["json"] = req
	}

	c.ServeJSON()
	return
}
