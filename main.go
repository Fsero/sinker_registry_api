package main

import (
	_ "bitbucket.org/fseros/sinker_registry_api/routers"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
	log.SetLevel(log.WarnLevel)
}
