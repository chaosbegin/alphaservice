package controllers

import (
	"alphawolf.com/alphaservice/impls"
	"alphawolf.com/alphaservice/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// SshController operations for ssh terminal
type SshController struct {
	beego.Controller
}

// Ssh terminal ...
// @Title Ssh terminal
// @Description Ssh terminal
// @Param	targetId		query 	string	true		"targetId"
// @Param	optionId		query 	string	true		"optionId"
// @Param	uuid            query 	string	true		"uuid"
// @Success 200 connected
// @Failure 403 disconnected
// @router /ws/terminal [get]
func (c *SshController) WsTerminal() {
	c.EnableRender = false

	targetId, err := c.GetInt("targetId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid targetId"
		c.ServeJSON()
		return
	}

	target, err := models.GetTargetById(targetId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid target"
		c.ServeJSON()
		return
	}

	optionId, err := c.GetInt("optionId")
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid optionId"
		c.ServeJSON()
		return
	}

	option, err := models.GetTargetOptionById(optionId)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid option"
		c.ServeJSON()
		return
	}

	uuid := c.GetString("uuid")

	userId := c.Ctx.Input.Session("uid").(int)

	user, err := models.GetUserById(userId)
	if err != nil {
		logs.Error("invalid login user, " + err.Error())
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "invalid login user"
		c.ServeJSON()
		return
	}

	userGroupIds := make([]int, 0)

	o := orm.NewOrm()
	_, err = o.Raw("select group_id from user_owner where user_id = ?", user.Id).QueryRows(&userGroupIds)
	if err != nil && err != orm.ErrNoRows {
		c.Ctx.Output.SetStatus(403)
		errMsg := "invalid login user group information"
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		c.ServeJSON()
		return
	}

	//logs.Trace("optionId:",optionId," targetId:",targetId," userId:",userId," userGroupIds:",userGroupIds)

	ok, err := impls.OperateAclSrv.ConnectCheck(option.Id, target.Id, target.GroupId, user.Id, userGroupIds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		errMsg := "connect check failed, " + err.Error()
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	} else if ok {
		c.Ctx.Output.SetStatus(403)
		errMsg := "deny connect"
		logs.Error(errMsg)
		c.Data["json"] = errMsg
		impls.OperateAuditSrv.Add(uuid, 1, user, target, option, c.Ctx.Request.RemoteAddr, errMsg)
		c.ServeJSON()
		return
	}

	logs.Trace("start upgrade to websocket...")

	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		logs.Error("websocket upgrade failed, " + err.Error())
		return
	}
	//defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()

	defer func() {
		impls.SshTerminalSrv.TerminalMap.Delete(uuid)
		logs.Info("client disconnected :" + remoteAddr)
	}()

	sshPty := impls.NewSshPty(uuid, conn, userGroupIds, user, target, option)

	impls.SshTerminalSrv.TerminalMap.Store(uuid, sshPty)

	err = sshPty.Run()
	if err != nil {
		logs.Error("run ssh pty failed, " + err.Error())
		conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
	}
	sshPty.Stop()

	return
}
