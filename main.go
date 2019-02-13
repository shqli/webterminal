package main

import (
	_ "webterminal/routers"
	"github.com/astaxie/beego"
)
func init(){
	beego.BConfig.Listen.Graceful = true
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600000
}

func main() {
	beego.Run()

}

