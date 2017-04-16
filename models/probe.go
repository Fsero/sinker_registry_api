package models

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"

	_ "github.com/mattn/go-sqlite3"
)

// Model Struct
type Probe struct {
	ProbeId       string `orm:"pk"`
	FQDN          string `orm:"size(100)"`
	Ipv4          string
	Ipv6          string
	Provider      string `orm:"size(100)"`
	GeoLongitude  float64
	GeoLatitude   float64
	SshPrivateKey string
	SshPublicKey  string
	Enabled       bool
}

var o orm.Ormer

func init() {
	orm.RegisterModel(new(Probe))
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	o = orm.NewOrm()
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		log.Error(err)
	}
}

func getHash(probe Probe) string {

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, probe)
	hash := fmt.Sprintf("%x", sha512.Sum512_224(bin_buf.Bytes()))
	return hash
}

func AddOne(probe Probe) (ProbeId string) {
	hash := getHash(probe)
	probe.ProbeId = hash
	_, err := o.Insert(&probe)
	if err != nil {
		log.Fatal(err)
	}
	return probe.ProbeId
}

func GetOne(ProbeId string) (probe *Probe, err error) {
	var id string = ProbeId
	pr := Probe{ProbeId: id}
	err = o.Read(&pr)
	if err == orm.ErrNoRows {
		log.Fatalf("No result found for id %s", ProbeId)
		return nil, errors.New("ProbeId not found")
	} else if err == orm.ErrMissPK {
		log.Fatalf("No primary key found for id %s.", ProbeId)
		return nil, errors.New("ProbeId not found")
	} else {
		log.Debugf("%v", probe)
		return &pr, nil
	}

}

func GetAll() []*Probe {
	var probes []*Probe
	num, err := o.QueryTable("probe").All(&probes)
	log.Debugf("Returned Rows Num: %s, %s", num, err)
	return probes
}

func Update(ProbeId string, Score int64) (err error) {
	return errors.New("ProbeId Not Exist")
}

func Delete(ProbeId string) {
	return
}
