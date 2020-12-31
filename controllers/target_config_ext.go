package controllers

import (
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

// Get target config by target_id ...
// @Title Get target config by target_id
// @Description Get target config by target_id
// @Param	id		query 	string	true		"The key for target_id"
// @Success 200 {object} models.TargetConfig
// @Failure 403 :id is empty
// @router /query [get]
func (c *TargetConfigController) Query() {
	targetId, err := c.GetInt("id")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	if targetId < 1 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid target_id:" + c.GetString("id")
		c.ServeJSON()
		return
	}

	targetConfig := &models.TargetConfig{}
	o := orm.NewOrm()
	err = o.Raw("select * from target_config where target_id = ? limit 1", targetId).QueryRow(targetConfig)
	if err != nil {
		if err == orm.ErrNoRows {
			target, err := models.GetTargetById(targetId)
			if err != nil {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = "query target failed, " + err.Error()
				c.ServeJSON()
				return
			}
			targetConfig.TargetId = targetId
			targetConfig.Ip = target.Address
			targetConfig.Config = "{}"

			tid, err := models.AddTargetConfig(targetConfig)
			if err != nil {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = "add default target config failed, " + err.Error()
				c.ServeJSON()
				return
			}

			targetConfig.Id = int(tid)

			c.Data["json"] = targetConfig
			c.ServeJSON()
			return

		} else {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "query target config failed, " + err.Error()
			c.ServeJSON()
			return
		}

	}

	c.Data["json"] = targetConfig
	c.ServeJSON()
	return
}
