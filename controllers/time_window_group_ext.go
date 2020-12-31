package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

// Get timeWindow item status ...
// @Title Get timeWindow item status
// @Description Get timeWindow item status
// @Param	pid		query	string	false		"tree start pid"
// @Param	layer	query	string	false		"tree layer deep"
// @Param	subNo	query 	string	false		"sub group no"
// @Param	groupOnly	query 	string	false		"query parent only"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /tree [get]
func (c *TimeWindowGroupController) TimeWindowTree() {
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

	groups := make([]*models.TimeWindowGroup, 0)
	o := orm.NewOrm()
	sql := "select * from time_window_group"

	_, err := o.Raw(sql).QueryRows(&groups)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = "query timeWindow group failed, " + err.Error()
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
			tGroups := make([]*models.TimeWindowGroup, 0)
			_, err = o.Raw("select * from time_window_group where id in (" + impls.IntArrayJoin(pgids) + ")").QueryRows(&tGroups)
			if err != nil {
				c.Data["json"] = "query parent timeWindow_group failed, " + err.Error()
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

	timeWindows := make([]*models.TimeWindow, 0)
	if !groupOnly {
		sql = "select * from time_window"
		_, err = o.Raw(sql).QueryRows(&timeWindows)
		if err != nil {
			c.Data["json"] = "query time_window failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	tree := c.GetTimeWindowTree(false, pid, layer, subNo, groups, timeWindows)

	c.Data["json"] = tree
	c.ServeJSON()
	return

}

func (c *TimeWindowGroupController) GetTimeWindowTree(sub bool, pid int, layer int, subNo int, parentRows []*models.TimeWindowGroup, subRows []*models.TimeWindow) []map[string]interface{} {
	rows := make([]map[string]interface{}, 0)
	l := layer
	l--
	if l < 0 {
		return rows
	}

	if sub {
		for _, s := range subRows {
			//logs.Trace("s.GroupId:",s.GroupId," pid:",pid)
			if s.GroupId == pid {
				row := make(map[string]interface{})
				row["id"] = s.Id + subNo
				row["pid"] = pid
				row["layer"] = layer
				row["name"] = s.Name
				row["checked"] = false
				row["show"] = true
				row["uid"] = s.UserId

				rows = append(rows, row)
			}
		}

		//logs.Trace("sub rows:",rows)

		return rows
	} else {
		for _, p := range parentRows {
			if p.Pid == pid {
				row := make(map[string]interface{})
				row["id"] = p.Id
				row["pid"] = p.Pid
				row["layer"] = layer
				row["name"] = p.Name
				row["checked"] = false
				row["show"] = true
				row["uid"] = p.UserId

				children := c.GetTimeWindowTree(false, p.Id, l, subNo, parentRows, subRows)
				if len(children) > 0 {
					row["children"] = children

				} else {
					if len(subRows) > 0 {

						children := c.GetTimeWindowTree(true, p.Id, l, subNo, parentRows, subRows)
						row["children"] = children
					}
				}
				rows = append(rows, row)
			}
		}
	}

	return rows
}
