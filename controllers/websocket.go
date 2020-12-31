package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"time"
)

// WebsocketController operations for ssh terminal
type WebsocketController struct {
	beego.Controller
}

// Websocket send ...
// @Title Websocket send
// @Description Websocket send
// @Param	id		query 	string	true		"id"
// @Param	name		query 	string	true		"name"
// @Success 200 connected
// @Failure 403 disconnected
// @router /send [get]
func (c *WebsocketController) WsSend() {
	c.EnableRender = false
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		logs.Error("websocket upgrade failed, " + err.Error())
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			logs.Error("websocket read message failed, " + err.Error())
			break
		}
		logs.Trace("websocket recevie message:", string(message))

		err = conn.WriteMessage(mt, message)
		if err != nil {
			logs.Error("websocket write message failed, " + err.Error())
			break
		}
	}

	logs.Trace("websocket exit...")
}

// Websocket terminal ...
// @Title Websocket terminal
// @Description Websocket terminal
// @Param	id		query 	string	true		"targetId"
// @Param	name		query 	string	true		"optionId"
// @Success 200 connected
// @Failure 403 disconnected
// @router /receive [get]
func (c *WebsocketController) WsReceive() {
	c.EnableRender = false
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		logs.Error("websocket upgrade failed, " + err.Error())
		return
	}
	defer conn.Close()

	for {
		//mt, message, err := conn.ReadMessage()
		//if err != nil {
		//	logs.Error("websocket read message failed, " + err.Error())
		//	break
		//}
		//logs.Trace("websocket recevie message:", string(message))
		conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"receiver":"good test !"}`))
		if err != nil {
			logs.Error("websocket write message failed, " + err.Error())
			break
		}

		time.Sleep(5 * time.Second)
	}

	logs.Trace("websocket exit...")
}
