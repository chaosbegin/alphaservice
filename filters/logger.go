package filters

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

func FilterReqLogger(ctx *context.Context) {
	uid := ctx.Input.Session("uid")
	logs.Info("req|", ctx.Input.Host(), "|", uid, "|", ctx.Input.Method(), "|", ctx.Input.URI(), "|", getParamString(ctx.Input.Params()), "|", string(ctx.Input.RequestBody))
}

func getParamString(param map[string]string) string {
	s := ""
	if param != nil {
		for k, v := range param {
			s += k + ":" + v + ";"
		}
	}
	return s
}

func FilterResLogger(ctx *context.Context) {
	var uid interface{}
	if ctx.Input.URL() != "/v1/auth/logout" {
		uid = ctx.Input.Session("uid")
	}
	logs.Info("res|", ctx.Input.Host(), "|", uid, "|", ctx.Output.Context.ResponseWriter.Status, "|", ctx.Output.Context.ResponseWriter.Elapsed.Milliseconds(), "ms|", ctx.Input.Method(), "|", ctx.Input.URI())
}
