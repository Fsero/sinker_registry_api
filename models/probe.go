package models

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

// Model Struct
type Probe struct {
	ProbeID       string    `orm:"pk" json:"ProbeID"`
	FQDN          string    `orm:"size(100)" json:"fqdn"`
	Ipv4          string    `json:"ipv4"`
	Ipv6          string    `json:"ipv6"`
	Provider      string    `orm:"size(100)" json:"provider"`
	GeoLongitude  float64   `json:"geolongitude"`
	GeoLatitude   float64   `json:"geolatitude"`
	SshPrivateKey string    `json:"sshprivatekey"`
	SshPublicKey  string    `json:"sshpublickey"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var o orm.Ormer

func init() {
	orm.RegisterModel(new(Probe))
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	o = orm.NewOrm()
	forced, verbose := true, true
	err := orm.RunSyncdb("default", forced, verbose)
	if err != nil {
		log.Error(err)
	}
}

func toHash(probe Probe) (hashID string) {
	out, err := json.Marshal(probe)
	if err == nil {
		log.Infof("%s", out)
	} else {
		log.Warningf("Json marshall failed %s", err)
	}
	sum := sha512.Sum512_224(out)
	log.Infof("Sum %v", sum)
	strsum := fmt.Sprintf("%s", sum)
	hash := base64.StdEncoding.EncodeToString([]byte(strsum))
	hashID = fmt.Sprintf("%s", hash)
	return hashID
}

// TODO: probe fields should be validated
func AddOne(probe Probe) (ProbeID string) {

	hashID := toHash(probe)
	probe.ProbeID = hashID
	probe.CreatedAt = time.Now()
	probe.UpdatedAt = time.Now()
	_, err := o.Insert(&probe)
	if err != nil {
		log.Fatal(err)
	}
	return probe.ProbeID
}

func GetOne(ProbeID string) (probe *Probe, err error) {
	var id string = ProbeID
	pr := Probe{ProbeID: id}
	err = o.Read(&pr)
	if err == orm.ErrNoRows {
		log.Warningf("No result found for id %s", ProbeID)
		return nil, errors.New("ProbeID not found")
	} else if err == orm.ErrMissPK {
		log.Warningf("No primary key found for id %s.", ProbeID)
		return nil, errors.New("ProbeID not found")
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

func Update(ProbeID string, Score int64) (err error) {
	return errors.New("ProbeID Not Exist")
}

func Delete(ProbeID string) {
	return
}
