package controllers

import (
	"bitbucket.org/fseros/sinker_registry_api/models"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

// Operations about probe
type ProbeController struct {
	beego.Controller
}

func (c *ProbeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.Get)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title Create Probe
// @Description create new probe
// @Success 200 {object} models.Probe
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
	json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
	log.Info("%v", pr)
	probeid := models.AddOne(pr)
	p.Data["json"] = map[string]string{"ProbeId": probeid}
	p.ServeJSON()
}

// @router /:id [get]
func (p *ProbeController) Get() {
	probeId := p.Ctx.Input.Param(":id")
	if probeId != "" {
		ob, err := models.GetOne(probeId)
		if err != nil || !ob.Enabled {
			p.Data["json"] = err.Error()
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

func (p *ProbeController) Put() {
	return
}

func (p *ProbeController) Delete() {
	return
}
