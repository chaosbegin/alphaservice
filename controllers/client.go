package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"strings"
)

// ClientController operations for Alert
type ClientController struct {
	beego.Controller
}

type ClientBinInfo struct {
	Name      string
	Version   string
	OS        string
	ARCH      string
	BuildTime string
	FullName  string
}

// Download client file...
// @Title Download client file
// @Description Download client file
// @Param	name 	query	string	true	"client setup file name"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /download [get]
func (c *ClientController) Download() {
	name := c.GetString("name")

	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "\\", "")
	name = strings.ReplaceAll(name, "*", "")

	if len(name) < 1 {
		c.Data["json"] = "invalid file name"
		c.Ctx.Output.SetStatus(403)
		c.ServeJSON()
		return
	}

	url := "static/client"
	url = filepath.Join(url, name)

	c.Ctx.Output.Download(url, name)
	return
}

// Download client file...
// @Title Download client file
// @Description Download client file
// @Param	os 	query	string	true	"client os"
// @Param	version 	query	string	true	"client version"
// @Success 200 {object} controllers.RespMsg
// @Failure 403 {string} {Success:false,Message:"错误信息..."}
// @router /list [get]
func (c *ClientController) List() {
	exePath, _ := os.Executable()
	pwd := filepath.Dir(exePath)
	downloadDir := filepath.Join(pwd, "static", "client")

	binInfos := make([]*ClientBinInfo, 0)

	filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logs.Error("path:", path, " name:", info.Name())

			name := info.Name()
			ext := filepath.Ext(name)
			cols := strings.Split(name, "-")
			if len(cols) != 5 {
				logs.Error("invalid client setup file:", path)
				return nil
			}

			binInfo := &ClientBinInfo{
				Name:      cols[0],
				OS:        cols[1],
				ARCH:      cols[2],
				Version:   cols[3],
				BuildTime: strings.TrimSuffix(cols[4], ext),
				FullName:  name,
			}

			binInfos = append(binInfos, binInfo)
		}
		return nil
	})

	c.Data["json"] = binInfos
	c.ServeJSON()
	return
}
