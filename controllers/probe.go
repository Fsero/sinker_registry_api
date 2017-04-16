package controllers

import (
	"bitbucket.org/fseros/sinker_registry_api/models"
	"encoding/json"
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

// @router / [post]
func (p *ProbeController) Post() {
	var pr models.Probe
	json.Unmarshal(p.Ctx.Input.RequestBody, &pr)
	probeid := models.AddOne(pr)
	p.Data["json"] = map[string]string{"ProbeId": probeid}
	p.ServeJSON()
}

// @router /:id [get]
func (p *ProbeController) Get() {
	probeId := p.Ctx.Input.Param(":probeId")
	if probeId != "" {
		ob, err := models.GetOne(probeId)
		if err != nil {
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
	probeId := p.Ctx.Input.Param(":probeId")
	models.Delete(probeId)
	p.Data["json"] = "delete success!"
	p.ServeJSON()
}
