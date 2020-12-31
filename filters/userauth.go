package filters

import (
	"strings"

	"alphawolf.com/alphaservice/controllers"
	"github.com/astaxie/beego/context"
)

func FilterUserAuth(ctx *context.Context) {
	if strings.HasPrefix(ctx.Input.URL(), "/v1/auth") ||
		strings.HasPrefix(ctx.Input.URL(), "/v1/demo") ||
		strings.HasPrefix(ctx.Input.URL(), "/v1/ssh/ws/terminal") ||
		strings.HasPrefix(ctx.Input.URL(), "/v1/ws") ||
		strings.HasPrefix(ctx.Input.URL(), "/v1/webConfig") ||
		strings.HasPrefix(ctx.Input.URL(), "/v1/grafana") {
		return
	}

	_, ok := ctx.Input.Session("uid").(int)
	if !ok {
		ctx.Input.Context.SetCookie("AlphaServiceSessionId", "")
		ctx.Input.Context.SetCookie("userInfo", "")
		ctx.Output.SetStatus(302)
		_ = ctx.Output.JSON(controllers.GetErrMsg("未授权"), true, true)
		return
	}
}
