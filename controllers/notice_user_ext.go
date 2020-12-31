package controllers

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/orm"
)

type NoticeUserBulkAddMsg struct {
	AlertLevel     int
	CategoryId     int
	ItemIds        []int
	ItemTplIds     []int
	TargetIds      []int
	TargetGroupIds []int
	UserIds        []int
}

// Post ...
// @Title Post
// @Description create NoticeUser
// @Param	body		body 	models.NoticeUser	true		"body for NoticeUser content"
// @Success 200 {string} success
// @Failure 403 body is error message
// @router /bulkAdd [post]
func (c *NoticeUserController) BulkAdd() {
	var msg NoticeUserBulkAddMsg
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid req msg parameter, " + err.Error()
		c.ServeJSON()
		return
	}

	userId := c.Ctx.Input.Session("uid").(int)
	noticeUsers := make([]*models.NoticeUser, 0)

	for _, v := range msg.ItemIds {
		for _, u := range msg.UserIds {
			nu := &models.NoticeUser{
				AlertLevel: msg.AlertLevel,
				CategoryId: msg.CategoryId,
				ItemId:     v,
				UserId:     u,
				CreateUid:  userId,
			}
			noticeUsers = append(noticeUsers, nu)
		}
	}

	for _, v := range msg.ItemTplIds {
		for _, u := range msg.UserIds {
			nu := &models.NoticeUser{
				AlertLevel: msg.AlertLevel,
				CategoryId: msg.CategoryId,
				ItemTplId:  v,
				UserId:     u,
				CreateUid:  userId,
			}
			noticeUsers = append(noticeUsers, nu)
		}
	}

	for _, v := range msg.TargetIds {
		for _, u := range msg.UserIds {
			nu := &models.NoticeUser{
				AlertLevel: msg.AlertLevel,
				CategoryId: msg.CategoryId,
				TargetId:   v,
				UserId:     u,
				CreateUid:  userId,
			}
			noticeUsers = append(noticeUsers, nu)
		}
	}

	for _, v := range msg.TargetGroupIds {
		for _, u := range msg.UserIds {
			nu := &models.NoticeUser{
				AlertLevel:    msg.AlertLevel,
				CategoryId:    msg.CategoryId,
				TargetGroupId: v,
				UserId:        u,
				CreateUid:     userId,
			}
			noticeUsers = append(noticeUsers, nu)
		}
	}

	o := orm.NewOrm()
	_, err = o.InsertMulti(len(noticeUsers), noticeUsers)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "insert to notice table failed, " + err.Error()
		c.ServeJSON()
		return
	}
	c.ServeJSON()
	return
}
