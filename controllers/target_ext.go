package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// Get target item status ...
// @Title Get target item status
// @Description Get target item status
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /itemStatus [get]
func (c *TargetController) GetItemStatus() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	o := orm.NewOrm()
	type itemStatus struct {
		Id     int
		Status int
	}

	ss := make([]*itemStatus, 0)

	if roleId > 2 {
		targetGroupIds := make([]int, 0)
		_, err := o.Raw("select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)", userId).QueryRows(&targetGroupIds)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target group id by user failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		if len(targetGroupIds) == 0 {
			c.Data["json"] = GetDataMsg(ss)
			c.ServeJSON()
			return
		}

		_, err = o.Raw("select a.id,count(b.status) as status from target a left join item b on a.id = b.target_id and b.status != 1 and a.group_id in (" + impls.IntArrayJoin(targetGroupIds) + ") group by a.id having count(b.status) > 0").QueryRows(&ss)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target item status failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		c.Data["json"] = GetDataMsg(ss)
		c.ServeJSON()
		return
	} else {
		_, err := o.Raw("select a.id,count(b.status) as status from target a left join item b on a.id = b.target_id and b.status != 1 group by a.id having count(b.status) > 0").QueryRows(&ss)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target item status failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		c.Data["json"] = GetDataMsg(ss)
		c.ServeJSON()
		return
	}

	return
}

// Get target alert status ...
// @Title Get target alert status
// @Description Get target alert status
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /alertStatus [get]
func (c *TargetController) GetAlertStatus() {
	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	o := orm.NewOrm()
	type itemStatus struct {
		Id     int
		Status int
	}

	ss := make([]*itemStatus, 0)

	if roleId > 2 {
		targetGroupIds := make([]int, 0)
		_, err := o.Raw("select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)", userId).QueryRows(&targetGroupIds)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target group id by user failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		if len(targetGroupIds) == 0 {
			c.Data["json"] = GetDataMsg(ss)
			c.ServeJSON()
			return
		}

		_, err = o.Raw("select a.id,count(b.status) as status from target a left join alert b on a.id = b.target_id and b.status = 1 and a.group_id in (" + impls.IntArrayJoin(targetGroupIds) + ") group by a.id having count(b.status) > 0").QueryRows(&ss)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target item status failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		c.Data["json"] = GetDataMsg(ss)
		c.ServeJSON()
		return
	} else {
		_, err := o.Raw("select a.id,count(b.status) as status from target a left join alert b on a.id = b.target_id and b.status = 1 group by a.id having count(b.status) > 0").QueryRows(&ss)
		if err != nil {
			c.Data["json"] = GetErrMsg("Get target item status failed, " + err.Error())
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		c.Data["json"] = GetDataMsg(ss)
		c.ServeJSON()
		return
	}

	return
}

// Set target group_id ...
// @Title Set target group_id
// @Description Set target group_id
// @Param	id	query	string	true	"target id"
// @Param	group_id	query	string	true	"target group_id"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /setGroupId [get]
func (c *TargetController) SetGroupId() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid target id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if id < 1 {
		c.Data["json"] = "invalid target id"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	groupId, err := c.GetInt("group_id")
	if err != nil {
		c.Data["json"] = "invalid target groupId, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	target, err := models.GetTargetById(id)
	if err != nil || target == nil {
		c.Data["json"] = "invalid target id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	group, err := models.GetTargetGroupById(groupId)
	if err != nil || group == nil {
		c.Data["json"] = "get move to target group by id failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	pGroup, err := models.GetTargetGroupById(target.GroupId)
	if err != nil || pGroup == nil {
		c.Data["json"] = "get target group by id failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if len(group.ApiId) > 0 || len(pGroup.ApiId) > 0 {
		c.Data["json"] = "系统自动管理类目标组不支持此操作！"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	target.GroupId = groupId

	err = models.UpdateTargetById(target)

	if err != nil || target == nil {
		c.Data["json"] = "update target groupId failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = ""
	c.ServeJSON()
	return
}

// Get target item status ...
// @Title Get target item status
// @Description Get target item status
// @Param	pid		query	string	false		"tree start pid"
// @Param	layer	query	string	false		"tree layer deep"
// @Param	subNo	query 	string	false		"sub group no"
// @Param	groupOnly	query 	string	false		"query parent only"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /targetTree [get]
func (c *TargetController) TargetTree() {
	pid, _ := c.GetInt("pid")
	layer, _ := c.GetInt("layer")
	if layer < 1 {
		layer = 10
	}

	subNo, _ := c.GetInt("subNo")
	if subNo < 1 {
		subNo = 1000000
	}

	groupOnly, _ := c.GetBool("groupOnly", false)

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)
	groups := make([]*models.TargetGroup, 0)
	o := orm.NewOrm()
	sql := ""
	if roleId > 2 {
		sql = "select * from target_group where id in (select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = " + strconv.Itoa(userId) + "))"
	} else {
		sql = "select * from target_group"
	}

	_, err := o.Raw(sql).QueryRows(&groups)
	if err != nil {
		c.Data["json"] = "query target group failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	gids := make([]int, len(groups))
	pgids := make([]int, 0)
	gIdMap := make(map[int]int)
	for k, v := range groups {
		gids[k] = v.Id
		gIdMap[v.Id] = v.Id

	}

	for _, v := range groups {
		if v.Pid > 0 {
			_, ok := gIdMap[v.Pid]
			if !ok {
				pgids = append(pgids, v.Pid)
			}
		}
	}

	gids = util.UniqIntArray(gids)
	pgids = util.UniqIntArray(pgids)

	if len(pgids) > 0 {
		for i := 0; i < layer; i++ {
			tGroups := make([]*models.TargetGroup, 0)
			_, err = o.Raw("select * from target_group where id in (" + impls.IntArrayJoin(pgids) + ")").QueryRows(&tGroups)
			if err != nil {
				c.Data["json"] = "query parent target_group failed, " + err.Error()
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}

			if len(tGroups) > 0 {
				groups = append(groups, tGroups...)
				pgids = make([]int, 0)
				for _, g := range tGroups {
					if g.Pid > 0 {
						_, ok := gIdMap[g.Pid]
						if !ok {
							pgids = append(pgids, g.Pid)
						}
					}
				}

				pgids = util.UniqIntArray(pgids)

				if len(pgids) > 0 {
					continue
				} else {
					break
				}
			}
		}
	}

	targets := make([]*models.Target, 0)
	if !groupOnly {
		if roleId > 2 {
			if len(gids) > 0 {
				sql = "select * from target where group_id in (" + impls.IntArrayJoin(gids) + ")"
				_, err = o.Raw(sql).QueryRows(&targets)
				if err != nil {
					c.Data["json"] = "query target failed, " + err.Error()
					c.Ctx.Output.SetStatus(403)
					c.ServeJSON()
					return
				}
			}
		} else {
			sql = "select * from target"
			_, err = o.Raw(sql).QueryRows(&targets)
			if err != nil {
				c.Data["json"] = "query target failed, " + err.Error()
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
		}
	}

	//for _,t := range targets{
	//	logs.Info(*t)
	//}
	//
	//for _,g := range groups{
	//	logs.Info(*g)
	//}
	type IdCount struct {
		Id    int
		Total int
	}

	targetAlertMap := make(map[int]int)
	targetAlerts := make([]IdCount, 0)
	targetItemStatusMap := make(map[int]int)
	targetItemStatus := make([]IdCount, 0)

	//target alert
	_, err = o.Raw("select a.id id,count(*) total from target a,alert b where a.id = b.target_id and b.`status` = 1 group by a.id").QueryRows(&targetAlerts)
	if err != nil {
		c.Data["json"] = "query target alert count failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	//target item status
	_, err = o.Raw("select a.id id,count(*) total from target a,item b where a.id = b.target_id and b.`status` = 0 group by a.id").QueryRows(&targetItemStatus)
	if err != nil {
		c.Data["json"] = "query target alert count failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	for _, t := range targetAlerts {
		targetAlertMap[t.Id] = t.Total
	}

	for _, t := range targetItemStatus {
		targetItemStatusMap[t.Id] = t.Total
	}

	_, _, tree := c.GetTargetTree(false, pid, layer, subNo, groups, targets, targetAlertMap, targetItemStatusMap)

	c.Data["json"] = tree
	c.ServeJSON()
	return

}

func (c *TargetController) GetTargetTree(sub bool, pid int, layer int, subNo int, parentRows []*models.TargetGroup, subRows []*models.Target, targetAlert map[int]int, targetItemStatus map[int]int) (int, int, []map[string]interface{}) {
	rows := make([]map[string]interface{}, 0)
	l := layer
	l--
	if l < 0 {
		return 0, 0, rows
	}

	ac, ic := 0, 0
	if sub {
		for _, s := range subRows {
			//logs.Trace("s.GroupId:",s.GroupId," pid:",pid)
			if s.GroupId == pid {
				row := make(map[string]interface{})
				row["id"] = s.Id + subNo
				row["pid"] = pid
				row["layer"] = layer
				row["name"] = s.Name + "[" + s.Address + "]"
				row["checked"] = false
				row["show"] = true

				row["adminAddr"] = s.AdminAddress

				//alert count and item disable status
				ao, ok := targetAlert[s.Id]
				if ok {
					ac += ao
					row["alertCount"] = ao
				} else {
					row["alertCount"] = 0
				}

				io, ok := targetItemStatus[s.Id]
				if ok {
					ic += io
					row["diCount"] = io
				} else {
					row["diCount"] = 0
				}

				//logs.Trace("parent row:",row)

				rows = append(rows, row)
			}
		}

		//logs.Trace("sub rows:",rows)

		return ac, ic, rows
	} else {
		var lastP *models.TargetGroup
		for _, p := range parentRows {
			//logs.Trace("p.Pid:",p.Pid," pid:",pid)
			if pid < 1 {
				ac = 0
				ic = 0
			}

			if p.Pid == pid {
				row := make(map[string]interface{})
				row["id"] = p.Id
				row["pid"] = p.Pid
				row["layer"] = layer
				row["name"] = p.Name
				row["checked"] = false
				row["show"] = true

				row["autoScan"] = p.AutoScan
				row["apiId"] = p.ApiId

				ao, io, children := c.GetTargetTree(false, p.Id, l, subNo, parentRows, subRows, targetAlert, targetItemStatus)
				if len(children) > 0 {
					row["children"] = children

					ac += ao
					ic += io
					if lastP != nil {
						row["alertCount"] = ao
						row["diCount"] = io
					} else {
						row["alertCount"] = ac
						row["diCount"] = ic
					}

				} else {
					if len(subRows) > 0 {

						ao, io, children := c.GetTargetTree(true, p.Id, l, subNo, parentRows, subRows, targetAlert, targetItemStatus)
						//logs.Error("p.Name:",p.Name," ao:",ao," io:",io)
						row["children"] = children

						ac += ao
						ic += io
						row["alertCount"] = ac
						row["diCount"] = ic

						if lastP != nil {
							row["alertCount"] = ao
							row["diCount"] = io
						} else {
							row["alertCount"] = ac
							row["diCount"] = ic
						}

					}
				}
				//logs.Error("p.Name:",p.Name," l:",l," levelSeq:",levelSeq," ac:",ac," ic:",ic)
				rows = append(rows, row)

				lastP = p

			}
			//logs.Trace("parent rows:",rows)
		}
	}

	//logs.Trace("total rows:",rows)

	return ac, ic, rows
}

// Get target by category_id ...
// @Title Get target by category_id
// @Description Get target by category_id
// @Param	category_id		query	string	false		"category id"
// @Success 200 {object} []models.Target
// @Failure 403 {string} error message
// @router /byCategory [get]
func (c *TargetController) ByCategory() {
	categoryId, _ := c.GetInt("category_id")
	if categoryId < 1 {
		c.Data["json"] = "invalid category_id parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	o := orm.NewOrm()
	targets := make([]models.Target, 0)
	if roleId > 2 {
		_, err := o.Raw("select * from target where group_id in (select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)) and id in (select target_id from target_category_detail where category_id = ?)", userId, categoryId).QueryRows(&targets)
		if err != nil && err != orm.ErrNoRows {
			c.Data["json"] = "query target by category:" + strconv.Itoa(categoryId) + " failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	} else {
		_, err := o.Raw("select * from target where id in (select target_id from target_category_detail where category_id = ?)", categoryId).QueryRows(&targets)
		if err != nil && err != orm.ErrNoRows {
			c.Data["json"] = "query target by category:" + strconv.Itoa(categoryId) + " failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = targets
	c.ServeJSON()
	return
}

// Get panel by target id ...
// @Title Get panel by target id
// @Description Get panel by target id
// @Param	id 	query	string	true	"target id"
// @Success 200 {object} []impls.Panel
// @Failure 403
// @router /panel [get]
func (c *TargetController) Panel() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = "invalid target id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	ps, err := impls.GetPanelByTargetId(id)
	if err != nil {
		c.Data["json"] = "invalid target id, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ps
		c.ServeJSON()
		return
	}

}
