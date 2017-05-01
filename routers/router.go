// @APIVersion 1.0.0
// @Title probe API
// @Description A simple probe registry API
// @Contact fsero
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
