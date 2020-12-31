package controllers

import (
	"alphawolf.com/alpha/util"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"

	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego/logs"
)

//  Add user...
// @Title Add user
// @Description create User
// @Param	groupId	query	string	false	"user group id"
// @Param	userIds query	[]int	false	"user ids"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /addOwner [get]
func (c *UserController) AddOwner() {
	groupId, _ := c.GetInt("groupId")
	if groupId < 1 {
		c.Data["json"] = "invalid groupId"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userIdsString := c.GetString("userIds", "")

	userIds := strings.Split(userIdsString, ",")
	for _, u := range userIds {
		uid, err := strconv.Atoi(u)
		if err != nil {
			c.Data["json"] = "invalid user id:" + u
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
		uo := &models.UserOwner{
			GroupId: groupId,
			UserId:  uid,
		}

		_, err = models.AddUserOwner(uo)
		if err != nil {
			c.Data["json"] = "add user owner failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	c.ServeJSON()
	return
}

//  Add user...
// @Title Add user
// @Description create User
// @Param	groupId	query	string	false	"user group id"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /add [post]
func (c *UserController) Add() {
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId != 1 {
		c.Data["json"] = "权限不足"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	user := models.User{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = "invalid body parameter, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	groupId, _ := c.GetInt("groupId", -1)
	if groupId < 1 {
		c.Data["json"] = "invalid groupId"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	count := 0
	o := orm.NewOrm()
	err = o.Raw("select count(*) from user where login_name = ?", user.LoginName).QueryRow(&count)
	if err != nil {
		c.Data["json"] = "get user count failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if count > 0 {
		c.Data["json"] = "用户名已经存在，请重新输入"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}
	if len(user.Password) > 0 {
		user.Password, err = impls.PwdDecrypt(user.Password)
		if err != nil {
			c.Data["json"] = "密码错误"
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	user.Password = impls.CommonSrv.PwdHash(user.Password)

	uid, err := models.AddUser(&user)
	if err != nil {
		c.Data["json"] = "add user failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userOwner := &models.UserOwner{
		GroupId: groupId,
		UserId:  int(uid),
	}

	_, err = models.AddUserOwner(userOwner)
	if err != nil {
		c.Data["json"] = "add user owner failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()

	return
}

//  Add user...
// @Title Add user
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /noticeAdd [post]
func (c *UserController) NoticeAdd() {
	user := models.User{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = "invalid body parameter, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	count := 0
	o := orm.NewOrm()
	err = o.Raw("select count(*) from user where login_name = ?", user.LoginName).QueryRow(&count)
	if err != nil {
		c.Data["json"] = "get user count failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if count > 0 {
		c.Data["json"] = "用户名已经存在，请重新输入"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	user.Status = 0
	user.NoticeStatus = 1
	user.Password = ""

	uid, err := models.AddUser(&user)
	if err != nil {
		c.Data["json"] = "add user failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	groupId := 0
	err = o.Raw("select value from sys_config where id = 15").QueryRow(&groupId)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = "query default notice user group_id failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if groupId < 1 {
		groupId = 2
	}

	userOwner := &models.UserOwner{
		GroupId: groupId,
		UserId:  int(uid),
	}

	_, err = models.AddUserOwner(userOwner)
	if err != nil {
		c.Data["json"] = "add user owner failed, " + err.Error()
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.ServeJSON()

	return
}

// Remove user...
// @Title Remove user
// @Description Remove User
// @Param	mode	query	string	false	"remove mode ,1: delete user and all user owner,2: remove user from this group"
// @Param	groupId	query	string	false	"user group id"
// @Param	userIds	query	string	false	"user ids"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /remove [get]
func (c *UserController) Remove() {
	mode, _ := c.GetInt("mode", -1)
	if mode < 1 || mode > 2 {
		c.Data["json"] = "invalid mode parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}
	groupId, _ := c.GetInt("groupId", -1)
	if groupId < 0 {
		c.Data["json"] = "invalid mode parameter"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userIds := c.GetString("userIds")

	switch mode {
	case 1: //delete
		o := orm.NewOrm()
		_, err := o.Raw("delete from user_owner where user_id in ("+userIds+") and group_id = ?", groupId).Exec()
		if err != nil {
			c.Data["json"] = "delete user owner failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
		_, err = o.Raw("delete from user where id in (" + userIds + ")").Exec()
		if err != nil {
			c.Data["json"] = "delete user failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	case 2: //remove
		type RemoveCount struct {
			UserId int
			Count  int
		}
		rcs := make([]RemoveCount, 0)

		o := orm.NewOrm()

		_, err := o.Raw("delete from user_owner where user_id in ("+userIds+") and group_id = ?", groupId).Exec()
		if err != nil {
			c.Data["json"] = "delete user owner failed, " + err.Error()
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		_, err = o.Raw("select user_id,count(*) as count from user_owner where user_id in (" + userIds + ") group by user_id").QueryRows(&rcs)
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

type ReqChangePwdMsg struct {
	UserId int
	OldPwd string
	NewPwd string
}

// Change user password ...
// @Title Change user password
// @Description Change user password
// @Param	body		body 	controllers.ReqChangePwdMsg	true		"body for ReqChangePwdMsg"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /changePwd [post]
func (c *UserController) ChangePwd() {
	msg := ReqChangePwdMsg{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.Data["json"] = GetErrMsg("invalid input parameter, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	roleId := c.Ctx.Input.Session("rid").(int)
	userId := c.Ctx.Input.Session("uid").(int)
	if roleId != 1 {
		if userId < 1 || userId != msg.UserId {
			c.Data["json"] = GetErrMsg("Insufficient permissions")
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	if len(msg.NewPwd) < 8 || len(msg.OldPwd) < 8 {
		c.Data["json"] = GetErrMsg("invalid password")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if msg.UserId == userId {
		if len(msg.OldPwd) < 8 {
			c.Data["json"] = GetErrMsg("invalid password")
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	user, err := models.GetUserById(msg.UserId)
	if err != nil {
		c.Data["json"] = GetErrMsg("invalid userid, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if msg.UserId == userId {
		deOldPwd, err := impls.PwdDecrypt(msg.OldPwd)
		if err != nil {
			errMsg := "密码无效"
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

		if impls.CommonSrv.PwdHash(deOldPwd) != user.Password {
			errMsg := "用户名或密码不正确"
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}

	}

	deNewPwd, err := impls.PwdDecrypt(msg.NewPwd)
	if err != nil {
		errMsg := "密码无效"
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	user.Password = impls.CommonSrv.PwdHash(deNewPwd)

	err = models.UpdateUserById(user)

	if err != nil {
		errMsg := "修改用户密码失败," + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.Data["json"] = GetOkMsg()
	c.Ctx.Output.SetStatus(200)
	c.ServeJSON()

	return
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router /getInfo [get]
func (c *UserController) GetInfo() {
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId != 1 {
		c.Data["json"] = GetErrMsg("Insufficient permissions")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	fields := []string{"", "", "", "", ""}
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1|col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, "|")
	}
	// order: desc|asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, "|")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, "|") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = "Error: invalid query key/value pair"
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetUserTmplByPaging ...
// @Title Get UserTmpl by paging
// @Description Get Alert by paging
// @Param	groupId	query	string	false	"user group id"
// @Param	query	query	string	false	"Filter. e.g. col1:v1|col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1|col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1|col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc|asc ..."
// @Param	pageSize	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	pageNo		query	string	false	"Start position of page. Must be an integer"
// @Param	headPage		query	string	false	"Store page infomation to header. Must be an bool"
// @Success 200 {object} models.Alert
// @Failure 403
// @router /page [get]
func (c *UserController) GetUserByPaging() {
	roleId := c.Ctx.Input.Session("rid").(int)
	if roleId != 1 {
		c.Data["json"] = GetErrMsg("Insufficient permissions")
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	var groupId int
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var pageSize int64 = 10
	var pageNo int64
	var headPage bool = false

	// groupId
	groupId, _ = c.GetInt("groupId", -1)

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, "|")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("pageSize"); err == nil {
		pageSize = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("pageNo"); err == nil {
		pageNo = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, "|")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, "|")
	}
	// headFlag
	if v, err := c.GetBool("headPage"); err == nil {
		headPage = v
	}

	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, "|") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = GetErrMsg("Error: invalid query key/value pair")
				c.Ctx.Output.SetStatus(403)
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	page, err := models.GetUserByPage(groupId, query, fields, sortby, order, pageNo, pageSize)
	if err != nil && err != orm.ErrNoRows {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
	} else {
		if headPage {
			c.Ctx.Output.Header("Total", strconv.Itoa(int(page.Total)))
			c.Ctx.Output.Header("PageNo", strconv.Itoa(int(page.PageNo)))
			c.Ctx.Output.Header("PageSize", strconv.Itoa(int(page.PageSize)))
			c.Data["json"] = page.Rows
		} else {
			c.Data["json"] = GetDataMsg(page)
		}

	}

	c.ServeJSON()
	return

}
