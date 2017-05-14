package models

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/abh/geoip"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Probe))
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	o = orm.NewOrm()
	forced, verbose := false, true
	err := orm.RunSyncdb("default", forced, verbose)
	if err != nil {
		log.Error(err)
	}
	gip = initializeGeoIP()
}

func initializeGeoIP() *geoip.GeoIP {
	file := "/usr/share/GeoIP/GeoIPCity.dat"

	gi, err := geoip.Open(file)
	if err != nil {
		fmt.Printf("Could not open GeoIP database\n")
	}
	return gi
}

var o orm.Ormer
var gip *geoip.GeoIP
