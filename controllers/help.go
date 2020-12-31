package controllers

import (
	"alphawolf.com/alphaservice/impls"
	"github.com/astaxie/beego"
	"net/url"
)

// HelpController operations for doc
type HelpController struct {
	beego.Controller
}

// Get help html by key
// @Title Get help html by key
// @Description Get help doc
// @Param	key 	query	string	true	"help key"
// @Success 200 {string} "html text"
// @Failure 403 {string} "error message"
// @router /content [get]
func (c *HelpController) Content() {
	key, err := url.QueryUnescape(c.GetString("key"))
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Ctx.Output.Body([]byte("invalid key parameter, " + key))
		return
	}
	hBytes, ok := impls.HelpSrv.DocMap.Load(key)
	if ok {
		data := make(map[string]interface{})
		data["html"] = string(hBytes.([]byte))
		c.Data["json"] = data
		c.ServeJSON()
		return
	} else {
		c.Ctx.Output.SetStatus(403)
		c.Ctx.Output.Body([]byte("can't find the key:" + key))
		return
	}
}

// Get help toc
// @Title Get help toc
// @Description Get help toc
// @Success 200 {string} "toc json array"
// @Failure 403 {string} "error message"
// @router /toc [get]
func (c *HelpController) Toc() {
	toc := impls.HelpSrv.GetToc()
	c.Data["json"] = toc
	c.ServeJSON()
	return
}

// Refresh help content
// @Title Refresh help content
// @Description Refresh help content
// @Success 200 {string} ""
// @Failure 403 {string} "error message"
// @router /refresh [get]
func (c *HelpController) Refresh() {
	impls.HelpSrv.Refresh()
	c.Data["json"] = ""
	c.ServeJSON()
	return
}
