package models

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/abh/geoip"
	"github.com/asaskevich/govalidator"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type Provider int

const (
	DIGITALOCEAN = 1 + iota
	VULTR
	AWS
	GOOGLECLOUD
	LINODE
	ERROR
)

var providers = map[string]Provider{
	"Digital Ocean": DIGITALOCEAN,
	"Vultr":         VULTR,
	"AWS":           AWS,
	"Google Cloud":  GOOGLECLOUD,
	"Linode":        LINODE,
	"Bad Provider":  ERROR,
}

func (p Provider) String() string {
	for i, k := range providers {
		if k == p {
			return i
		}
	}
	return ""
}

func initializeGeoIP() *geoip.GeoIP {
	file := "/usr/share/GeoIP/GeoIPCity.dat"

	gi, err := geoip.Open(file)
	if err != nil {
		fmt.Printf("Could not open GeoIP database\n")
	}
	return gi
}

func ParseProvider(descr string) (Provider, bool) {
	if val, ok := providers[descr]; ok {
		return val, ok
	}
	return ERROR, false
}

// Model Struct
type Probe struct {
	ProbeID       string    `orm:"pk" json:"ProbeID"`
	FQDN          string    `orm:"size(100)" json:"fqdn"`
	Ipv4          string    `json:"ipv4"`
	Ipv6          string    `json:"ipv6"`
	Provider      string    `orm:"size(100)" json:"provider"`
	GeoLongitude  string    `json:"geolongitude"`
	GeoLatitude   string    `json:"geolatitude"`
	Country       string    `json:"country"`
	SSHPrivateKey string    `json:"SSHPrivateKey"`
	SSHPublicKey  string    `json:"SSHPublicKey"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var o orm.Ormer
var gip *geoip.GeoIP

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

func Validate(probe Probe) (bool, error) {
	if _, ok := ParseProvider(probe.Provider); !ok {
		return false, fmt.Errorf("invalid provider %s", probe.Provider)
	}
	if len(probe.Ipv4) < 7 {
		return false, fmt.Errorf("bad ipv4, too short for an ipv4 address %s", probe.Ipv4)
	}
	if !govalidator.IsIPv4(probe.Ipv4) {
		return false, fmt.Errorf("unable to parse this `%s` as ipv4", probe.Ipv4)
	}
	if probe.Ipv6 != "" {
		if !govalidator.IsIPv6(probe.Ipv6) {
			return false, fmt.Errorf("unable to parse this `%s` as ipv6", probe.Ipv6)
		}
	}
	if !govalidator.IsDNSName(probe.FQDN) {
		return false, fmt.Errorf("unable to parse this `%s` as FQDN", probe.FQDN)
	}
	if (probe.GeoLatitude != "NaN" && probe.GeoLatitude != "") && !govalidator.IsLatitude(probe.GeoLatitude) {
		return false, fmt.Errorf("invalid latitude `%s`", probe.GeoLatitude)
	}

	return true, nil

}
func (probe *Probe) SetDefaults() {
	probe.GeoLatitude = "NaN"
	probe.GeoLongitude = "NaN"
	probe.Country = ""
	probe.FQDN = ""
	probe.Enabled = false
	probe.Ipv4 = ""
	probe.Ipv6 = ""
	probe.Provider = ""
	probe.SSHPrivateKey = ""
	probe.SSHPublicKey = ""
	probe.CreatedAt = time.Now()
	probe.UpdatedAt = time.Now()

}

// TODO: probe fields should be validated
func AddOne(probe Probe) (ProbeID string, err error) {

	hashID := toHash(probe)
	probe.ProbeID = hashID
	probe.CreatedAt = time.Now()
	probe.UpdatedAt = time.Now()

	ok, err := Validate(probe)
	if !ok {
		return "", err
	}
	record := gip.GetRecord(probe.Ipv4)
	if record != nil {
		probe.GeoLatitude = fmt.Sprintf("%f", record.Latitude)
		probe.GeoLongitude = fmt.Sprintf("%f", record.Longitude)
		probe.Country = record.CountryName
	}

	fmt.Printf("%+v", record)
	ok, err = Validate(probe)
	if !ok {
		return "", err
	}

	_, err = o.Insert(&probe)
	if err != nil {
		log.Fatal(err)
	}
	return probe.ProbeID, nil
}

func GetOne(ProbeID string) (probe *Probe, err error) {
	var id = ProbeID
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

func GetByIPv4(ProbeIP string) ([]Probe, error) {
	var ip = ProbeIP
	var probes []Probe
	num, err := o.Raw("SELECT * FROM probe where ipv4 = ?", ip).QueryRows(&probes)
	if err == nil {
		fmt.Println("nums: ", num)
	}

	return probes, err

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
