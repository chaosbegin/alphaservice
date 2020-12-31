package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// Bulk add category items ...
// @Title Bulk add category items
// @Description Bulk add category items
// @Param	categoryId	query	string	false	"category id"
// @Param	order		query	string	false	"item start order"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /bulkAdd [post]
func (c *TargetCategoryItemController) BulkAdd() {
	categoryId, err := c.GetInt("categoryId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid categoryId, " + err.Error()
		c.ServeJSON()
		return
	}

	category, err := models.GetTargetCategoryById(categoryId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid categoryId, " + err.Error()
		c.ServeJSON()
		return
	}

	order, err := c.GetInt("order")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid order, " + err.Error()
		c.ServeJSON()
		return
	}

	tplIds := make([]int, 0)
	err = util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &tplIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid itemTplIds, " + err.Error()
		c.ServeJSON()
		return
	}

	for _, id := range tplIds {
		itemTpl, err := models.GetItemTmplById(id)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "invalid itemTplId: " + strconv.Itoa(id) + ", " + err.Error()
			c.ServeJSON()
			return
		}

		targetCategoryItem := &models.TargetCategoryItem{
			CategoryId: category.Id,
			TmplTypeId: itemTpl.TmplTypeId,
			ItemId:     id,
			Order:      order,
		}

		_, err = models.AddTargetCategoryItem(targetCategoryItem)
		if err != nil {
			c.Ctx.Output.SetStatus(403)
			c.Data["json"] = "add target category item failed, " + err.Error()
			c.ServeJSON()
			return
		}

		order++
	}

	c.ServeJSON()
	return
}

// Swap category item order ...
// @Title Swap category item order
// @Description Swap category item order
// @Param	id1		query	string	false	"first category item id"
// @Param	id2		query	string	false	"second category item id"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /swapOrder [post]
func (c *TargetCategoryItemController) SwapOrder() {
	id1, err := c.GetInt("id1")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid id1, " + err.Error()
		c.ServeJSON()
		return
	}

	id2, err := c.GetInt("id2")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid id2, " + err.Error()
		c.ServeJSON()
		return
	}

	item1, err := models.GetTargetCategoryItemById(id1)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid id1, " + err.Error()
		c.ServeJSON()
		return
	}

	item2, err := models.GetTargetCategoryItemById(id2)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid id2, " + err.Error()
		c.ServeJSON()
		return
	}

	tOrder := item1.Order
	item1.Order = item2.Order
	item2.Order = tOrder

	err = models.UpdateTargetCategoryItemById(item1)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "update item failed, " + err.Error()
		c.ServeJSON()
		return
	}

	err = models.UpdateTargetCategoryItemById(item2)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "update item failed, " + err.Error()
		c.ServeJSON()
		return
	}

	c.ServeJSON()
	return
}

// Bulk delete category items ...
// @Title Bulk delete category items
// @Description Bulk delete category items
// @Param	ids 	query	string	false	"category item ids,ex 1,2,3"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /bulkDel [post]
func (c *TargetCategoryItemController) BulkDel() {
	ids := c.GetString("ids")

	if len(ids) < 1 {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid ids"
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	_, err := o.Raw("delete from target_category_item where id in (" + ids + ")").Exec()
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "delete target_category_item failed, " + err.Error()
		c.ServeJSON()
		return
	}

	c.ServeJSON()
	return
}
