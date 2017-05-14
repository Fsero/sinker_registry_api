package controllers

import (
	"encoding/json"

	"fmt"

	"bitbucket.org/fseros/sinker_registry_api/models"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

// Operations about probe
type ProbeController struct {
	beego.Controller
}

func (p *ProbeController) URLMapping() {
	p.Mapping("Post", p.Post)
	p.Mapping("GetOne", p.Get)
	p.Mapping("GetAll", p.GetAll)
	p.Mapping("Put", p.Put)
	p.Mapping("Delete", p.Delete)
}

// @Title Create Probe
// @Description create new probe
// @Success 201 {object} models.Probe
// @Param  fqdn  body string true "fqdn address of the probe"
// @Param  ipv4  body string true "ipv4 address of the probe"
// @Param  ipv6  body string false "ipv6 address of the probe"
// @Param  provider  body string true "cloud provider of the probe"
// @Param  geolongitude  body float false "geolongitude of the probe"
// @Param  geolatitude  body float false "geolongitude of the probe"
// @Param  sshprivatekey  body string false "ssh private key of the probe"
// @Param  sshpublickey  body string false "ssh public key of the probe"
// @Param  enabled  body bool false "probe status" "true"
// @router / [post]
func (p *ProbeController) Post() {
	var pr models.Probe
	pr.SetDefaults()
	json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
	log.Debugf(" received %v via POST", pr)
	probeid, err := models.AddOne(pr)
	if err != nil {
		p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		p.Ctx.Output.SetStatus(400)
		p.ServeJSON()
		return
	}
	p.Data["json"] = map[string]string{"ProbeId": probeid}
	p.Ctx.Output.SetStatus(201)
	p.ServeJSON()
}

// @router /:id [get]
func (p *ProbeController) Get() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		ob, err := models.GetByID(ProbeID)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}

// @router / [get]
func (p *ProbeController) GetAll() {
	obs := models.GetAll()
	p.Data["json"] = obs
	p.ServeJSON()
}

// @router /ip/:ip [get]
func (p *ProbeController) GetByIP() {
	probeIP := p.Ctx.Input.Param(":ip")
	fmt.Printf("Looking for probes with ip %s", probeIP)
	if probeIP != "" {
		obs, err := models.GetByIPv4(probeIP)
		if err == nil {
			fmt.Printf("Found \n %+v", obs)
			newobs := make([]models.Probe, 0)
			for _, ob := range obs {
				if ob.Enabled {
					newobs = append(newobs, ob)
				}
			}
			fmt.Printf("Found newobs \n %+v", newobs)
			p.Data["json"] = newobs
			p.Ctx.Output.SetStatus(200)
			p.ServeJSON()
		} else {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
			p.Ctx.Output.SetStatus(500)
			p.ServeJSON()
			p.StopRun()
		}
	}
}

// @router /disable/?:id [put]
func (p *ProbeController) Disable() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.Disable]: disabling probe %s", ProbeID)
		ob, err := models.Disable(ProbeID)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}

func getIDbyQueryParamOrAsAParam(p *ProbeController) string {
	var ProbeID string
	ProbeID = p.GetString("id")
	if ProbeID == "" {
		ProbeID = p.Ctx.Input.Param(":id")
	}
	return ProbeID

}

// @router /enable/?:id [put]
func (p *ProbeController) Enable() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.Enable]: enabling probe %s", ProbeID)
		ob, err := models.Enable(ProbeID)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}

// @Title Updates traces path
// @Description updates traces path
// @Success 200 {object} models.Probe
// @Param  tracespath  body string false "traces path for probe"
// @router /tracespath/?:id [put]
func (p *ProbeController) UpdateTracesPath() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.UpdateTracesPath]: updating traces path for probe %s", ProbeID)
		var pr models.Probe
		json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
		ob, err := models.UpdateTracesPath(ProbeID, pr.TracesPath)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}

// @Title Updates ssh key
// @Description updates ssh key
// @Success 200 {object} models.Probe
// @Param  sshprivatekey  body string false "ssh private key of the probe"
// @Param  sshpublickey  body string false "ssh public key of the probe"
// @router /ssh/?:id [put]
func (p *ProbeController) UploadSSH() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.UploadSSH]: updating ssh key for probe %s", ProbeID)
		var pr models.Probe
		json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
		ob, err := models.UploadSSH(ProbeID, pr.SSHPrivateKey, pr.SSHPublicKey)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}

// @Title Updates ssh key
// @Description updates ssh key
// @Success 200 {object} models.Probe
// @Param  sshprivatekey  body string false "ssh private key of the probe"
// @Param  sshpublickey  body string false "ssh public key of the probe"
// @router /ssh/?:id [get]
func (p *ProbeController) GetSSH() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.GetSSH]: Getting ssh key for probe %s", ProbeID)
		keys, err := models.GetSSH(ProbeID)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = map[string]string{"ProbeId": ProbeID, "SSHPrivateKey": keys.Private, "SSHPublicKey": keys.Public}

		}
	}
	p.ServeJSON()
}

func (p *ProbeController) Put() {
	var pr models.Probe
	pr.SetDefaults()

	json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
	log.Debugf(" received %v via POST", pr)
	probeid, err := models.AddOne(pr)
	if err != nil {
		p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		p.Ctx.Output.SetStatus(400)
		p.ServeJSON()
		return
	}

	p.Data["json"] = map[string]string{"ProbeId": probeid}
	p.Ctx.Output.SetStatus(201)
	p.ServeJSON()
}

// @router /delete/?:id [put]
func (p *ProbeController) Delete() {
	ProbeID := getIDbyQueryParamOrAsAParam(p)
	if ProbeID != "" {
		log.Infof("[controllers.probe.Delete]: deleting probe %s", ProbeID)
		_, err := models.Delete(ProbeID)
		if err != nil {
			p.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		} else {
			p.Data["json"] = fmt.Sprintf("{ 'msg': 'deleted probe with id %s' }", ProbeID)
		}
	}
	p.ServeJSON()
}
