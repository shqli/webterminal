package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type MainController struct {
	beego.Controller
}

type ApiDeny struct {
	beego.Controller
}

func (c *MainController) Get() {
	if c.Ctx.Request.URL.Path == "/" {
		http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, "web/", 302)
	} else {
		http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
	}
}

func (c *ApiDeny) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ApiDeny) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ApiDeny) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ApiDeny) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}
