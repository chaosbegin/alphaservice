package controllers

import (
	"alphawolf.com/alpha/util"
	"strconv"

	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// ServiceController operations for voice notice
type AuthController struct {
	beego.Controller
}

type LoginInfo struct {
	UserName string
	Password string
}

// Login
// @Title Login
// @Description Login
// @Param	body		body 	controllers.LoginInfo	true		"body for login content"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /login [post]
func (c *AuthController) Login() {
	//_,ok := c.GetSession("user_info").(string)
	//if ok {
	//	c.Data["json"] = GetDataMsg("已经登录")
	//	c.ServeJSON()
	//	return
	//}
	//c.StartSession()

	loginInfo := &LoginInfo{}
	err := util.JsonIter.Unmarshal(c.Ctx.Input.RequestBody, &loginInfo)
	if err != nil {
		errMsg := "Parse LoginInfo failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if !impls.CommonSrv.UserNameRegx.MatchString(loginInfo.UserName) {
		errMsg := "无效的用户名: " + loginInfo.UserName
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	user := &models.User{}

	o := orm.NewOrm()
	err = o.Raw("select * from user where login_name = ?", loginInfo.UserName).QueryRow(user)

	if err != nil {
		if err == orm.ErrNoRows {
			errMsg := "用户名或密码不正确"
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		} else {
			errMsg := "Query user failed, " + err.Error()
			logs.Error(errMsg)
			c.Data["json"] = GetErrMsg(errMsg)
			c.Ctx.Output.SetStatus(403)
			c.ServeJSON()
			return
		}
	}

	if user.Status != 1 {
		errMsg := "用户被限制登录"
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	//logs.Trace(impls.CommonSrv.PwdHash(loginInfo.Password))

	pwd, err := impls.PwdDecrypt(loginInfo.Password)
	if err != nil {
		errMsg := "密码解密失败"
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	if impls.CommonSrv.PwdHash(pwd) != user.Password {
		errMsg := "用户名或密码不正确"
		logs.Error(errMsg)
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	c.SetSession("uid", user.Id)
	c.SetSession("rid", user.RoleId)

	type UserInfo struct {
		Id        int
		GroupIds  []int
		RoleId    int
		Name      string
		LoginName string
		Mobile    string
		Openid    string
		Email     string
		Acl       []string
	}
	userInfo := UserInfo{}

	_, err = o.Raw("select group_id from user_owner where user_id = ?", user.Id).QueryRows(&userInfo.GroupIds)
	if err != nil {
		errMsg := "获取用户组信息失败"
		logs.Error(errMsg, ","+err.Error())
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userInfo.Id = user.Id
	userInfo.LoginName = user.LoginName
	userInfo.Name = user.Name
	userInfo.Mobile = user.Mobile
	userInfo.Openid = user.Openid
	userInfo.Email = user.Email
	userInfo.RoleId = user.RoleId

	acl := make([]string, 0)

	acl, err = c.getAcl(user.Id, user.RoleId, userInfo.GroupIds)
	if err != nil {
		errMsg := "获取用户权限列表失败"
		logs.Error(errMsg, ", ", err.Error())
		c.Data["json"] = GetErrMsg(errMsg)
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	userInfo.Acl = acl

	c.Data["json"] = GetDataMsg(userInfo)
	c.ServeJSON()
	return

}

func (c *AuthController) getAcl(uid int, roleId int, gids []int) ([]string, error) {
	acls := make([]string, 0)

	o := orm.NewOrm()
	_, err := o.Raw("select b.key from user_right a,rights b where a.right_id = b.id and a.user_id = ?", uid).QueryRows(&acls)

	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}

	if len(gids) < 1 {
		return acls, nil
	}

	gids_str := ""
	for _, v := range gids {
		preGid := v
		gids_str += strconv.Itoa(preGid) + ","

		for i := 0; i < 5; i++ {
			group, err := models.GetUserGroupById(preGid)
			if err != nil {
				if err == orm.ErrNoRows {
					break
				}

				return nil, err
			}

			if group.Pid < 1 {
				break
			}
			gids_str += strconv.Itoa(group.Pid) + ","

			preGid = group.Pid

		}
	}

	gids_str = gids_str[:len(gids_str)-1]

	//logs.Trace("gids:",gids)

	groupAcls := make([]string, 0)
	_, err = o.Raw("select distinct(b.key) from user_group_right a,rights b  where a.right_id = b.id and a.group_id in (" + gids_str + ")").QueryRows(&groupAcls)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	//logs.Trace("gAcls:",groupAcls)
	acls = append(acls, groupAcls...)

	//role acls
	roleAcls := make([]string, 0)
	_, err = o.Raw("select distinct(b.key) from user_role_right a,rights b  where a.right_id = b.id and b.id = a.right_id and a.role_id = ?", roleId).QueryRows(&roleAcls)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	acls = append(acls, roleAcls...)

	aclMap := make(map[string]int, 0)
	for _, a := range acls {
		aclMap[a] = 1
	}

	fixAcls := make([]string, 0)
	for k, _ := range aclMap {
		fixAcls = append(fixAcls, k)
	}

	//logs.Trace("fixAcls:", fixAcls)

	return fixAcls, nil
}

// Logout
// @Title Logout
// @Description Logout
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /logout [get]
func (c *AuthController) Logout() {
	c.DestroySession()
	c.Data["json"] = GetOkMsg()
	c.ServeJSON()
	return

}
