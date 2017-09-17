package feedify

import (
	"github.com/astaxie/beego"
	"github.com/roscopecoltran/feedify/contextor"
)

type Controller struct {
	beego.Controller
}

func (c *Controller) GetInput() *contextor.Input {
	return &contextor.Input{c.Ctx.Input}
}

func (c *Controller) GetCtx() *contextor.Context {
	return &contextor.Context{c.Ctx}
}

func (c *Controller) GetJsonData() interface{} {
	return c.Controller.Data["json"]
}
