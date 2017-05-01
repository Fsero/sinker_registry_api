package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetByIP",
			Router: `/ip/:ip`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
