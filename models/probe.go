package models

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
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
	HETZNER
	ERROR
)

var providers = map[string]Provider{
	"Digital Ocean": DIGITALOCEAN,
	"Vultr":         VULTR,
	"AWS":           AWS,
	"Google Cloud":  GOOGLECLOUD,
	"Linode":        LINODE,
	"Hetzner":       HETZNER,
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
	SSHPrivateKey string    `json:"sshprivateKey"`
	SSHPublicKey  string    `json:"sshpublicKey"`
	TracesPath    string    `json:"tracespath"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProbeSSHKeys struct {
	Private string `json:"SSHPrivateKey"`
	Public  string `json:"SSHPublicKey"`
}

func toHash(probe Probe) (hashID string) {
	out, err := json.Marshal(probe)
	if err == nil {
		log.Debugf("%s", out)
	} else {
		log.Warningf("Json marshall failed %s", err)
	}
	sum := sha512.Sum512_224(out)
	log.Debugf("[models.toHash] checksum %s", sum)
	strsum := fmt.Sprintf("%s", sum)
	hash := base64.URLEncoding.EncodeToString([]byte(strsum))
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

	probes, err := GetByFQDN(probe.FQDN)
	log.Debugf("[models.probe.addOne]: number of probes %s found by FQDN", len(probes))

	if err == nil && len(probes) >= 1 {
		return false, fmt.Errorf("FQDN name already registered %s", probe.FQDN)
	} else if err != nil {
		log.Errorf("[models.probe.addOne]: Error querying database %s", err)
	}
	probes, err = GetByIPv4(probe.Ipv4)
	log.Debugf("[models.probe.addOne]: number of probes %s found by IP", len(probes))

	if err == nil && len(probes) >= 1 {
		return false, fmt.Errorf("IPv4 address already registered %s", probe.Ipv4)
	} else if err != nil {
		log.Errorf("[models.probe.addOne]: Error querying database %s", err)
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
	probe.TracesPath = "/var/log/traces"
	probe.CreatedAt = time.Now()
	probe.UpdatedAt = time.Now()

}

func AddOne(probe Probe) (ProbeID string, err error) {

	hashID := toHash(probe)
	probe.ProbeID = hashID
	probe.CreatedAt = time.Now()
	probe.UpdatedAt = time.Now()

	log.Infof("[models.AddOne] new probe %+v", probe)
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

func GetByID(ProbeID string) (probe *Probe, err error) {
	var id = ProbeID
	pr := Probe{ProbeID: id}
	err = o.Read(&pr)
	if err == orm.ErrNoRows {
		log.Warningf("[models.probe.GetByID]: No result found for id %s", ProbeID)
		return nil, errors.New("ProbeID not found")
	} else if err == orm.ErrMissPK {
		log.Warningf("[models.probe.GetByID]: No primary key found for id %s.", ProbeID)
		return nil, errors.New("ProbeID not found")
	} else {
		log.Debugf("[models.probe.GetByID]: %v", probe)
		return &pr, nil
	}

}

func GetByIPv4(ProbeIP string) ([]Probe, error) {
	var ip = ProbeIP
	var probes []Probe

	if !govalidator.IsIPv4(ip) {
		return probes, fmt.Errorf("Invalid IPv4 address provided %s", ip)
	}
	num, err := o.QueryTable("probe").Filter("enabled", true).Filter("Ipv4", ip).All(&probes)
	//num, err := o.Raw("SELECT * FROM probe where enabled = 1 and ipv4 = ?", ip).QueryRows(&probes)
	if err == nil {
		fmt.Println("nums: ", num)
	}
	return probes, err
}

func GetByFQDN(fqdn string) ([]Probe, error) {
	var probes []Probe
	if !govalidator.IsDNSName(fqdn) {
		return probes, fmt.Errorf("Invalid DNS name provided %s", fqdn)
	}
	num, err := o.QueryTable("probe").Filter("enabled", true).Filter("FQDN", fqdn).All(&probes)
	//num, err := o.Raw("SELECT * FROM probe where enabled = 1 and f_q_d_n = ?", fqdn).QueryRows(&probes)
	if err == nil {
		fmt.Println("nums: ", num)
	}
	return probes, err
}

func GetAll() []*Probe {
	var probes []*Probe
	num, err := o.QueryTable("probe").Filter("enabled", true).All(&probes)
	log.Debugf("[models.probe.GetAll]: Returned Rows Num: %s, %s", num, err)
	return probes
}

func Disable(ProbeID string) (*Probe, error) {
	log.Infof("[model.probe.Disable]: disabling probe %s", ProbeID)
	probe, err := GetByID(ProbeID)
	if err == nil {
		probe.Enabled = false
		probe.UpdatedAt = time.Now()
		_, err := o.Update(probe)
		if err != nil {
			return nil, err
		}
		return probe, nil
	}
	return nil, err
}

func Enable(ProbeID string) (*Probe, error) {
	log.Infof("[model.probe.Enable]: enabling probe %s", ProbeID)
	probe, err := GetByID(ProbeID)
	if err == nil {
		probe.Enabled = true
		probe.UpdatedAt = time.Now()
		_, err := o.Update(probe)
		if err != nil {
			return nil, err
		}
		return probe, nil
	}
	return nil, err

}
func Update(ProbeID string, Score int64) (err error) {
	return errors.New("ProbeID Not Exist")
}

func UploadSSH(ProbeID string, SSHPrivateKey string, SSHPublicKey string) (*Probe, error) {
	log.Infof("[model.probe.UploadSSH]: Uploading SSH %s", ProbeID)
	probe, err := GetByID(ProbeID)
	if err == nil {
		if !(govalidator.IsBase64(SSHPrivateKey) && govalidator.IsBase64(SSHPublicKey)) {
			log.Infof("[model.probe.UploadSSH] malformed base64 ssh key %s %s", SSHPrivateKey, SSHPublicKey)
			return nil, errors.New("Malformed base64 in SSH Keys")
		}

		probe.SSHPrivateKey = SSHPrivateKey
		probe.SSHPublicKey = SSHPublicKey
		probe.UpdatedAt = time.Now()
		_, err = o.Update(probe)
		log.Infof("[model.probe.UploadSSH]: Saving new key for probe %s", ProbeID)
		if err != nil {
			return nil, err
		}
		return probe, nil
	}
	return nil, err
}

func UpdateTracesPath(ProbeID string, traces_path string) (*Probe, error) {
	log.Infof("[model.probe.UpdateTracesPath]: updating traces path %s", ProbeID)
	probe, err := GetByID(ProbeID)
	if ok, _ := govalidator.IsFilePath(traces_path); !ok {
		return nil, fmt.Errorf("Invalid path for traces '%s'", traces_path)
	}

	if err == nil {
		probe.TracesPath = traces_path
		probe.UpdatedAt = time.Now()
		_, err = o.Update(probe)
		if err != nil {
			return nil, err
		}
		return probe, nil
	}
	return nil, err
}

func GetSSH(ProbeID string) (*ProbeSSHKeys, error) {
	log.Infof("[model.probe.GetSSH]: Getting SSH keys %s", ProbeID)
	probe, err := GetByID(ProbeID)
	keys := ProbeSSHKeys{Public: probe.SSHPublicKey, Private: probe.SSHPrivateKey}

	if err != nil {
		return nil, err
	}
	if keys.Private == "" || keys.Public == "" {
		return nil, fmt.Errorf("partial ssh content, unable to get ssh keys base64 encoded from '%s' '%s'", probe.SSHPublicKey, probe.SSHPrivateKey)
	}
	if !govalidator.IsBase64(keys.Public) || !govalidator.IsBase64(keys.Private) {
		return nil, fmt.Errorf("Invalid format for ssh keys, expect base64 encoding ones '%s' '%s'", keys.Public, keys.Private)
	}
	decodedPrivate, errPrivate := base64.URLEncoding.DecodeString(probe.SSHPrivateKey)
	if errPrivate != nil {
		return nil, errPrivate
	}
	decodedPublic, errPublic := base64.URLEncoding.DecodeString(probe.SSHPublicKey)
	if errPublic != nil {
		return nil, errPublic
	}
	keys.Private = fmt.Sprintf("%s", decodedPrivate)
	keys.Public = fmt.Sprintf("%s", decodedPublic)
	return &keys, nil
}

func Delete(ProbeID string) (bool, error) {
	log.Infof("[model.probe.Delete]: removing probe %s", ProbeID)
	probe, err := GetByID(ProbeID)
	if err == nil {
		probe.Enabled = true
		probe.UpdatedAt = time.Now()
		_, err = o.Delete(probe)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}
