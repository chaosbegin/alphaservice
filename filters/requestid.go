package filters

import (
	"github.com/astaxie/beego/context"
	"github.com/satori/go.uuid"
)

func FilterRequestID(ctx *context.Context) {
	requestID := ctx.Request.Header.Get("X-Request-ID")

	if requestID == "" {
		uuid4, _ := uuid.NewV4()
		requestID = uuid4.String()
		ctx.Request.Header.Add("X-Request-ID", requestID)
	}

	ctx.Output.Header("X-Request-ID", requestID)
}
