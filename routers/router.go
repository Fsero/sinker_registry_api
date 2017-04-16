package routers

import (
	"bitbucket.org/fseros/sinker_registry_api/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/probe",
			beego.NSInclude(
				&controllers.ProbeController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
