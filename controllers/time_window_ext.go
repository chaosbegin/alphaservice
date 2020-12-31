package controllers

import "alphawolf.com/alphaservice/models"

// Set timeWindow group pid ...
// @Title Set timeWindow group pid
// @Description Set timeWindow group pid
// @Param	id	query	string	true	"timeWindow group id"
// @Param	pid	query	string	true	"timeWindow group pid"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /setPid [get]
func (c *TimeWindowController) SetPid() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid timeWindow id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if id < 1 {
		c.Data["json"] = "invalid timeWindow id"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	timeWindow, err := models.GetTimeWindowById(id)
	if err != nil || timeWindow == nil {
		c.Data["json"] = "invalid timeWindow id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	pid, err := c.GetInt("pid")
	if err != nil {
		c.Data["json"] = "invalid timeWindow group id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	group, err := models.GetTimeWindowGroupById(pid)
	if err != nil || group == nil {
		c.Data["json"] = "invalid timeWindow group id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	timeWindow.GroupId = pid
	err = models.UpdateTimeWindowById(timeWindow)
	if err != nil {
		c.Data["json"] = "update timeWindow failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.ServeJSON()
	return
}
