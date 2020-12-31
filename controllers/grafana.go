package controllers

import (
	"alphawolf.com/alpha/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"strings"
)

// GrafanaController operations for doc
type GrafanaController struct {
	beego.Controller
}

type DashboardDataReplaceMsg struct {
	OldKeys    []string
	NewKeys    []string
	AddSuffix  string
	DelSuffix  string
	MatchQuote bool
}

// Dashboard data content replace
// @Title Dashboard data content replace
// @Description Dashboard data content replace
// @Param	body		body 	controllers.DashboardDataReplaceMsg	true		"body for replace message"
// @Success 200 {string} "replace success"
// @Failure 403 {string} "error message"
// @router /dashboard/data/replace [post]
func (c *GrafanaController) DashboardDataReplace() {
	msg := &DashboardDataReplaceMsg{}
	err := util.JsonIter.Unmarshal([]byte(c.Ctx.Input.RequestBody), msg)
	if err != nil {
		c.Data["json"] = GetErrMsg(err.Error())
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	type DashBoardData struct {
		Id   int
		Data string
	}

	dds := make([]DashBoardData, 0)

	o := orm.NewOrm()
	_, err = o.Raw("select id,data from grafana.dashboard").QueryRows(&dds)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Ctx.Output.Body([]byte("query dashboard table failed, " + err.Error()))
		return
	}

	oLen := len(msg.OldKeys)
	nLen := len(msg.NewKeys)

	if oLen < 1 {
		c.Ctx.Output.SetStatus(403)
		c.Ctx.Output.Body([]byte("no old keys input"))
		return
	}

	dataChanged := false

	if nLen > 0 {
		mLen := math.Min(float64(oLen), float64(nLen))
		for j := 0; j < len(dds); j++ {
			for i := 0; i < int(mLen); i++ {
				dataChanged = true
				if msg.MatchQuote {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, "\""+msg.OldKeys[i]+"\"", "\""+msg.NewKeys[i]+"\"")
				} else {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, msg.OldKeys[i], msg.NewKeys[i])
				}

			}
		}
	}

	if len(msg.AddSuffix) > 0 {
		for j := 0; j < len(dds); j++ {
			for i := 0; i < oLen; i++ {
				dataChanged = true
				if msg.MatchQuote {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, "\""+msg.OldKeys[i]+"\"", "\""+msg.OldKeys[i]+msg.AddSuffix+"\"")
				} else {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, msg.OldKeys[i], msg.OldKeys[i]+msg.AddSuffix)
				}
			}
		}
	}

	if len(msg.DelSuffix) > 0 {
		for j := 0; j < len(dds); j++ {
			for i := 0; i < oLen; i++ {
				dataChanged = true
				if msg.MatchQuote {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, "\""+msg.OldKeys[i]+"\"", "\""+strings.TrimSuffix(msg.OldKeys[i], msg.DelSuffix)+"\"")
				} else {
					dds[j].Data = strings.ReplaceAll(dds[j].Data, msg.OldKeys[i], strings.TrimSuffix(msg.OldKeys[i], msg.DelSuffix))
				}
			}
		}
	}

	//update data
	if dataChanged {
		for _, dd := range dds {
			_, err = o.Raw("update grafana.dashboard set data = ? where id = ?", dd.Data, dd.Id).Exec()
			if err != nil {
				c.Ctx.Output.SetStatus(403)
				c.Ctx.Output.Body([]byte("update dashboard failed, " + err.Error()))
				return
			}
		}
	}

	return
}
