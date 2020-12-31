// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"alphawolf.com/alphaservice/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/demo",
			beego.NSInclude(
				&controllers.DemoController{},
			),
		),
		beego.NSNamespace("/alert",
			beego.NSInclude(
				&controllers.AlertController{},
			),
		),

		beego.NSNamespace("/alert_config",
			beego.NSInclude(
				&controllers.AlertConfigController{},
			),
		),

		beego.NSNamespace("/item",
			beego.NSInclude(
				&controllers.ItemController{},
			),
		),

		beego.NSNamespace("/itemGroup",
			beego.NSInclude(
				&controllers.ItemGroupController{},
			),
		),

		beego.NSNamespace("/itemTmpl",
			beego.NSInclude(
				&controllers.ItemTmplController{},
			),
		),
		beego.NSNamespace("/itemType",
			beego.NSInclude(
				&controllers.ItemTypeController{},
			),
		),
		beego.NSNamespace("/itemCategory",
			beego.NSInclude(
				&controllers.ItemCategoryController{},
			),
		),

		beego.NSNamespace("/target",
			beego.NSInclude(
				&controllers.TargetController{},
			),
		),

		beego.NSNamespace("/targetGroup",
			beego.NSInclude(
				&controllers.TargetGroupController{},
			),
		),

		beego.NSNamespace("/targetOwner",
			beego.NSInclude(
				&controllers.TargetOwnerController{},
			),
		),

		beego.NSNamespace("/targetOption",
			beego.NSInclude(
				&controllers.TargetOptionController{},
			),
		),

		beego.NSNamespace("/targetConfig",
			beego.NSInclude(
				&controllers.TargetConfigController{},
			),
		),

		beego.NSNamespace("/timeWindow",
			beego.NSInclude(
				&controllers.TimeWindowController{},
			),
		),
		beego.NSNamespace("/timeWindowGroup",
			beego.NSInclude(
				&controllers.TimeWindowGroupController{},
			),
		),
		beego.NSNamespace("/holiday",
			beego.NSInclude(
				&controllers.HolidayController{},
			),
		),
		beego.NSNamespace("/alertLevel",
			beego.NSInclude(
				&controllers.AlertLevelController{},
			),
		),
		beego.NSNamespace("/alertStatus",
			beego.NSInclude(
				&controllers.AlertStatusController{},
			),
		),
		beego.NSNamespace("/noticeConfig",
			beego.NSInclude(
				&controllers.NoticeConfigController{},
			),
		),
		beego.NSNamespace("/itemConnOpts",
			beego.NSInclude(
				&controllers.ItemConnOptionController{},
			),
		),
		beego.NSNamespace("/execute",
			beego.NSInclude(
				&controllers.ExecuteController{},
			),
		),
		beego.NSNamespace("/config",
			beego.NSInclude(
				&controllers.ConfigController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/userGroup",
			beego.NSInclude(
				&controllers.UserGroupController{},
			),
		),
		beego.NSNamespace("/userOwner",
			beego.NSInclude(
				&controllers.UserOwnerController{},
			),
		),
		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),
		beego.NSNamespace("/rights",
			beego.NSInclude(
				&controllers.RightsController{},
			),
		),
		beego.NSNamespace("/rightsCategory",
			beego.NSInclude(
				&controllers.RightsCategoryController{},
			),
		),
		beego.NSNamespace("/service",
			beego.NSInclude(
				&controllers.ServiceController{},
			),
		),
		beego.NSNamespace("/sysconfig",
			beego.NSInclude(
				&controllers.SysConfigController{},
			),
		),
		beego.NSNamespace("/targetCategory",
			beego.NSInclude(
				&controllers.TargetCategoryController{},
			),
		),
		beego.NSNamespace("/targetCategoryDetail",
			beego.NSInclude(
				&controllers.TargetCategoryDetailController{},
			),
		),
		beego.NSNamespace("/targetCategoryGroup",
			beego.NSInclude(
				&controllers.TargetCategoryGroupController{},
			),
		),
		beego.NSNamespace("/targetCategoryItem",
			beego.NSInclude(
				&controllers.TargetCategoryItemController{},
			),
		),
		beego.NSNamespace("/targetType",
			beego.NSInclude(
				&controllers.TargetTypeController{},
			),
		),
		beego.NSNamespace("/targetPanel",
			beego.NSInclude(
				&controllers.TargetPanelController{},
			),
		),
		beego.NSNamespace("/internalAlert",
			beego.NSInclude(
				&controllers.InternalAlertController{},
			),
		),
		beego.NSNamespace("/internalAlertType",
			beego.NSInclude(
				&controllers.InternalAlertTypeController{},
			),
		),
		beego.NSNamespace("/noticeUser",
			beego.NSInclude(
				&controllers.NoticeUserController{},
			),
		),
		beego.NSNamespace("/pwdGroup",
			beego.NSInclude(
				&controllers.PwdGroupController{},
			),
		),
		beego.NSNamespace("/pwdTarget",
			beego.NSInclude(
				&controllers.PwdTargetController{},
			),
		),
		beego.NSNamespace("/pwdOption",
			beego.NSInclude(
				&controllers.PwdOptionController{},
			),
		),
		beego.NSNamespace("/pwdMgr",
			beego.NSInclude(
				&controllers.PwdMgrController{},
			),
		),
		beego.NSNamespace("/hundsunTradingCenter",
			beego.NSInclude(
				&controllers.HundsunTradingCenterController{},
			),
		),
		beego.NSNamespace("/ssh",
			beego.NSInclude(
				&controllers.SshController{},
			),
		),
		beego.NSNamespace("/databaseAccess",
			beego.NSInclude(
				&controllers.DatabaseAccessController{},
			),
		),
		beego.NSNamespace("/webConfig",
			beego.NSInclude(
				&controllers.WebconsoleConfigController{},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebsocketController{},
			),
		),
		beego.NSNamespace("/client",
			beego.NSInclude(
				&controllers.ClientController{},
			),
		),
		beego.NSNamespace("/help",
			beego.NSInclude(
				&controllers.HelpController{},
			),
		),
		beego.NSNamespace("/grafana",
			beego.NSInclude(
				&controllers.GrafanaController{},
			),
		),
		beego.NSNamespace("/bugDemo",
			beego.NSInclude(
				&controllers.BugDemoController{},
			),
		),
	)
	beego.AddNamespace(ns)
}