package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

// Bulk add target_option  ...
// @Title Bulk add target_option
// @Description Bulk add target_option
// @Param	body		body 	[]models.TargetOption	true		"body for TargetOption content array"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /bulk_add [post]
func (c *TargetOptionController) BulkAdd() {
	options := make([]*models.TargetOption, 0)
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &options)
	if err == nil {
		o := orm.NewOrm()
		_, err = o.InsertMulti(len(options), options)
		if err == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = GetOkMsg()
			c.ServeJSON()
			return
		} else {
			c.Data["json"] = GetErrMsg(err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	return

}

// Get targetOption by targetId and itemId ...
// @Title Get targetOption by targetId and itemId
// @Description Get targetOption by targetId and itemId
// @Param	targetId		query 	string	true		"targetId"
// @Param	itemId		query 	string	true		"itemId"
// @Success 200 {object} []models.TargetOption
// @Failure 403 error message
// @router /getOptionByTargetItemId [get]
func (c *TargetOptionController) GetOptionByTargetItemId() {
	targetId, err := c.GetInt("targetId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "无效的目标"
		c.ServeJSON()
		return
	}

	itemId, err := c.GetInt("itemId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "无效的指标"
		c.ServeJSON()
		return
	}

	target, err := models.GetTargetById(targetId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "无效的目标"
		c.ServeJSON()
		return
	}

	item, err := models.GetItemById(itemId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "无效的指标"
		c.ServeJSON()
		return
	}

	targetOptions := make([]models.TargetOption, 0)
	o := orm.NewOrm()
	_, err = o.Raw("select * from target_option where type_id < 2 and item_type_id = ? and target_id = ?", item.ItemTypeId, target.Id).QueryRows(&targetOptions)
	if err != nil && err != orm.ErrNoRows {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "query target_option failed, " + err.Error()
		c.ServeJSON()
		return
	}

	//clear password
	for i := 0; i < len(targetOptions); i++ {
		targetOptions[i].Password = ""
	}

	c.Data["json"] = targetOptions
	c.ServeJSON()
	return
}
