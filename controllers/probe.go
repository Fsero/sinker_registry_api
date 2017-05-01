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

func (c *ProbeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.Get)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
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
	log.Info("%v", pr)
	probeid, err := models.AddOne(pr)
	if err != nil {
		p.Data["json"] = err.Error()
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

// @router /ip/:ip [get]
func (p *ProbeController) GetByIP() {
	probeIP := p.Ctx.Input.Param(":id")
	if probeIP != "" {
		obs, err := models.GetByIPv4(probeIP)
		fmt.Printf("%+v", obs)
		if err != nil {
			newobs := make([]models.Probe, len(obs))
			for _, ob := range obs {
				if ob.Enabled {
					newobs = append(newobs, ob)
				}
			}
			p.Data["json"] = newobs
		} else {
			p.Data["json"] = err
			p.Ctx.Output.SetStatus(500)
			p.ServeJSON()
			return
		}

	}
	p.ServeJSON()
}

func (p *ProbeController) Put() {
	return
}

func (p *ProbeController) Delete() {
	return
}
