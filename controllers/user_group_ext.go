package controllers

import (
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

// Remove groups and users...
// @Title Remove groups and users
// @Description Remove groups and users
// @Param	groupIds  query	string	false	"user group ids"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /remove [get]
func (c *UserGroupController) Remove() {
	groupIds := c.GetString("groupIds")
	if len(groupIds) < 1 {
		c.Data["json"] = "invalid groupIds parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	userIds := make([]int, 0)
	_, err := o.Raw("select DISTINCT(user_id) from user_owner where group_id in (" + groupIds + ")").QueryRows(&userIds)
	if err != nil {
		c.Data["json"] = "get user ids failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	_, err = o.Raw("delete from user_owner where group_id in (" + groupIds + ")").Exec()
	if err != nil {
		c.Data["json"] = "delete user owner failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	_, err = o.Raw("delete from user_group where id in (" + groupIds + ")").Exec()
	if err != nil {
		c.Data["json"] = "delete user_group failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if len(userIds) > 0 {
		type RemoveCount struct {
			UserId int
			Count  int
		}
		rcs := make([]RemoveCount, 0)

		_, err = o.Raw("select user_id,count(*) as count from user_owner where user_id in (" + impls.IntArrayJoin(userIds) + ") group by user_id").QueryRows(&rcs)
		if err != nil {
			c.Data["json"] = "get user owner count failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		for _, v := range rcs {
			if v.Count == 0 {
				err = models.DeleteUser(v.UserId)
				if err != nil {
					c.Data["json"] = "delete user failed, " + err.Error()
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			}
		}
	}

	c.ServeJSON()
	return
}
