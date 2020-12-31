package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

// Update category items by id ...
// @Title Update category items by ids
// @Description Update category items by ids
// @Param	id			query 	string	true		"category id"
// @Param	body		body 	[]models.TargetCategoryItem	true		"body for itemTpl ids"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /updateItems [post]
func (c *TargetCategoryController) UpdateItems() {
	categoryId, err := c.GetInt("id")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	categoryItems := make([]models.TargetCategoryItem, 0)
	err = util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &categoryItems)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	} else {
		o := orm.NewOrm()
		o.Begin()
		_, err = o.Raw("delete from target_category_item where category_id = ?", categoryId).Exec()
		if err != nil {
			o.Rollback()
			c.Data["json"] = "delete target_category_item by id failed, " + err.Error()
			c.ServeJSON()
			return
		}

		_, err = o.InsertMulti(len(categoryItems), categoryItems)
		if err != nil {
			o.Rollback()
			c.Data["json"] = "multi insert into target_category_item by failed, " + err.Error()
			c.ServeJSON()
			return
		}

		err = o.Commit()
		if err != nil {
			o.Rollback()
			c.Data["json"] = "update category items commit failed, " + err.Error()
			c.ServeJSON()
			return
		}

		c.ServeJSON()
		return

	}

}

// Delete target category and category items,category target list ...
// @Title Delete target category and category items,category target list
// @Description Delete target category and category items,category target list
// @Param	body		body 	[]int	true		"category id array"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /remove [post]
func (c *TargetCategoryController) Remove() {
	ids := make([]int, 0)
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid id array parameter, " + err.Error()
		c.ServeJSON()
		return
	}

	if len(ids) < 1 {
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	roleId := c.Ctx.Input.Session("rid").(int)

	idsStr := impls.IntArrayJoin(ids)

	o := orm.NewOrm()
	sql := ""

	if roleId != 1 {
		tIds := make([]int, 0)
		sql = "select id from target_category where id in (" + idsStr + ") and user_id = ?"
		_, err = o.Raw(sql, userId).QueryRows(&tIds)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "query target_category failed, " + err.Error()
			c.ServeJSON()
			return
		}

		if len(tIds) < 1 {
			c.ServeJSON()
			return
		}

		idsStr = impls.IntArrayJoin(tIds)

	}

	o.Begin()
	//delete category item
	sql = "delete from target_category_item where category_id in (" + idsStr + ")"
	_, err = o.Raw(sql).Exec()
	if err != nil {
		o.Rollback()
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete target_category_item failed, " + err.Error()
		c.ServeJSON()
		return
	}

	//delete category detail
	sql = "delete from target_category_detail where category_id in (" + idsStr + ")"
	_, err = o.Raw(sql).Exec()
	if err != nil {
		o.Rollback()
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete target_category_detail failed, " + err.Error()
		c.ServeJSON()
		return
	}

	//delete category

	sql = "delete from target_category where id in (" + idsStr + ")"
	_, err = o.Raw(sql).Exec()
	if err != nil {
		o.Rollback()
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete target_category failed, " + err.Error()
		c.ServeJSON()
		return
	}

	err = o.Commit()
	if err != nil {
		o.Rollback()
		c.Data["json"] = "remove category commit failed, " + err.Error()
		c.ServeJSON()
		return
	}

	return

}

// Get category tree ...
// @Title category tree
// @Description category tree
// @Param	gid		query	string	false		"category_group_id"
// @Param	pid		query	string	false		"tree start pid"
// @Param	layer	query	string	false		"tree layer deep"
// @Param	body		body 	[]int	true		"category id array"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /categoryTree [get]
func (c *TargetCategoryController) CategoryTree() {
	gid, _ := c.GetInt("gid")
	pid, _ := c.GetInt("pid")
	layer, _ := c.GetInt("layer")
	if layer < 1 {
		layer = 10
	}

	categories := make([]*models.TargetCategory, 0)
	o := orm.NewOrm()
	_, err := o.Raw("select * from target_category where group_id = ? order by priority", gid).QueryRows(&categories)
	if err != nil {
		c.Data["json"] = "query target_category failed, " + err.Error()
		c.ServeJSON()
		return
	}

	nodes := getCategoryNodes(categories, pid, layer)

	c.Data["json"] = nodes
	c.ServeJSON()
	return
}

func getCategoryNodes(rows []*models.TargetCategory, pid int, layer int) []map[string]interface{} {
	nodes := make([]map[string]interface{}, 0)
	layer--
	if layer < 0 {
		return nodes
	}

	for _, r := range rows {
		if r.Pid == pid {
			node := make(map[string]interface{})
			node["id"] = r.Id
			node["name"] = r.Name
			node["checked"] = false
			node["pid"] = pid
			node["layer"] = layer

			//other option
			node["uid"] = r.UserId

			subNodes := getCategoryNodes(rows, r.Id, layer)
			if len(subNodes) > 0 {
				node["children"] = subNodes
			}

			nodes = append(nodes, node)

		}
	}

	return nodes
}
