package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

type GetTargetOwnerMsg struct {
	Id   int
	Name string
}

// GetAllByTargetGroupId ...
// @Title Get All by target group id
// @Description get TargetOwner by target group id
// @Param	targetGroupId	query	string	false	"target group id"
// @Success 200 {object} controllers.GetTargetOwnerMsg
// @Failure 403
// @router /get [get]
func (c *TargetOwnerController) GetAllByTargetGroupId() {
	id, _ := c.GetInt("targetGroupId")
	if id < 1 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid targetGroupId parameter"
		c.ServeJSON()
		return
	}

	msgs := make([]GetTargetOwnerMsg, 0)

	o := orm.NewOrm()
	_, err := o.Raw("select a.id id,b.name name FROM target_owner a, user_group b where a.user_group_id = b.id and a.target_group_id = ?", id).QueryRows(&msgs)
	if err != nil && err != orm.ErrNoRows {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "query target_owner failed, " + err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = msgs
	c.ServeJSON()
	return

}

// DelByTargetOwnerIds ...
// @Title delete target owners by ids
// @Description delete target owners by ids
// @Param	body		body	models.ParamIds		true		""
// @Success 200 {string} success
// @Failure 403 error message
// @router /del [post]
func (c *TargetOwnerController) DelByTargetOwnerIds() {
	paramIds := models.ParamIds{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &paramIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid ids parameter, " + err.Error()
		c.ServeJSON()
		return
	}

	if len(paramIds.Ids) < 1 {
		c.Data["json"] = ""
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	_, err = o.Raw("delete from target_owner where id in (" + impls.IntArrayJoin(paramIds.Ids) + ")").Exec()
	if err != nil && err != orm.ErrNoRows {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete target_owner by ids failed, " + err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = ""
	c.ServeJSON()
	return

}

type AddTargetOwnerMsg struct {
	TargetGroupId int
	UserGroupIds  []int
}

// GetAllByTargetGroupId ...
// @Title Get All by target group id
// @Description get TargetOwner by target group id
// @Param	body		body	models.ParamIds		true		""
// @Success 200 {string} success
// @Failure 403 error message
// @router /add [post]
func (c *TargetOwnerController) Add() {
	msg := AddTargetOwnerMsg{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid add parameter, " + err.Error()
		c.ServeJSON()
		return
	}

	if len(msg.UserGroupIds) < 1 {
		c.Data["json"] = ""
		c.ServeJSON()
		return
	}

	_, err = models.GetTargetGroupById(msg.TargetGroupId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid target group id"
		c.ServeJSON()
		return
	}

	for _, id := range msg.UserGroupIds {
		to := &models.TargetOwner{
			TargetGroupId: msg.TargetGroupId,
			UserGroupId:   id,
			ApiId:         "",
		}

		_, err = models.AddTargetOwner(to)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "add target owner failed, " + err.Error()
			c.ServeJSON()
			return
		}

	}

	c.Data["json"] = ""
	c.ServeJSON()
	return

}
