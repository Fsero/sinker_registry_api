package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetByIP",
			Router: `/ip/:ip`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetByFQDN",
			Router: `/name/:fqdn`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Disable",
			Router: `/disable/?:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Enable",
			Router: `/enable/?:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "UpdateTracesPath",
			Router: `/tracespath/?:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "UploadSSH",
			Router: `/ssh/?:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "GetSSH",
			Router: `/ssh/?:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/sinker_registry_api/controllers:ProbeController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/delete/?:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

}
