package controllers

import (
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// Get rights by roleId ...
// @Title Get rights by roleId
// @Description Get rights by roleId
// @Param	roleId	query	string	false	"role id"
// @Success 200 {string} success message
// @Failure 403 {string} error message
// @router /list [get]
func (c *RightsController) List() {
	roleId, _ := c.GetInt("roleId", -1)
	if roleId < 0 {
		c.Data["json"] = "invalid roleId parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	type RightsList struct {
		Id           int
		CategoryName string
		Name         string
	}

	ml := make([]RightsList, 0)
	o := orm.NewOrm()
	_, err := o.Raw("select b.id,c.name as 'category_name',b.name as 'name' from user_role_right a,rights b,rights_category c where a.right_id = b.id and b.category_id = c.id and role_id = ?", roleId).QueryRows(&ml)
	if err != nil {
		c.Data["json"] = "query right failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = ml
	c.ServeJSON()
	return
}

// Remove by roleId ...
// @Title Remove by roleId
// @Description Remove by roleId
// @Param	roleId	query	string	false	"role id"
// @Param	rightIds	query	string	false	"rights id"
// @Success 200 {string} success message
// @Failure 403 {string} error message
// @router /removeByRole [get]
func (c *RightsController) RemoveByRole() {
	roleId, _ := c.GetInt("roleId", -1)
	if roleId < 0 {
		c.Data["json"] = "invalid roleId parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}
	rightIds := c.GetString("rightIds")
	if len(rightIds) < 1 {
		c.Data["json"] = "invalid rightIds parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	_, err := o.Raw("delete from user_role_right where right_id in ("+rightIds+") and role_id = ?", roleId).Exec()
	if err != nil {
		c.Data["json"] = "remove right by roleId failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.ServeJSON()
	return

}

// Add right by roleId ...
// @Title Add right by roleId
// @Description Add right by roleId
// @Param	roleId	query	string	false	"role id"
// @Param	rightIds	query	string	false	"rights id"
// @Success 200 {string} success message
// @Failure 403 {string} error message
// @router /addByRole [get]
func (c *RightsController) AddByRole() {
	roleId, _ := c.GetInt("roleId", -1)
	if roleId < 0 {
		c.Data["json"] = "invalid roleId parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}
	rightIdsStr := c.GetString("rightIds")
	if len(rightIdsStr) < 1 {
		c.Data["json"] = "invalid rightIds parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	rightIds := strings.Split(rightIdsStr, ",")
	for _, v := range rightIds {
		id, err := strconv.Atoi(v)
		if err != nil {
			c.Data["json"] = "invalid right id:" + v
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		urr := &models.UserRoleRight{
			RightId: id,
			RoleId:  roleId,
		}

		_, err = models.AddUserRoleRight(urr)
		if err != nil {
			c.Data["json"] = "add user role right failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	}

	c.ServeJSON()
	return

}
