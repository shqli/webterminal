package controllers
import (
	"github.com/astaxie/beego"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"fmt"
)

const(
	ShellUrl 			= "/api/v1.0/web/shell"
	ShellUrlPath 			= "/api/v1.0/web"
	ShellUrlImgPng			= ShellUrlPath  + 	"/static/img/favicon.png"
	ShellUrlJsJqueryMinJs      	= ShellUrlPath  +	"/static/js/jquery.min.js"
	ShellUrlJsPopperMinJs 		= ShellUrlPath  +	"/static/js/popper.min.js"
	ShellUrlJsBootstrapMinJs 	= ShellUrlPath  +	"/static/js/bootstrap.min.js"
	ShellUrlJsXtermMinJs 		= ShellUrlPath  +	"/static/js/xterm.min.js"
	ShellUrlJsFullScreenMinJs 	= ShellUrlPath  +	"/static/js/fullscreen.min.js"
	ShellUrlCssBootstrapMinCss 	= ShellUrlPath	+	"/static/css/bootstrap.min.css"
	ShellUrlCssXtermMinCss		= ShellUrlPath	+	"/static/css/xterm.min.css"
	ShellUrlCssFullScreenMinCss	= ShellUrlPath	+	"/static/css/fullscreen.min.css"
	ShellUrlShellWs 		= ShellUrlPath  +	"/shell/ws"
	ShellUrlShellStatic 		= ShellUrlPath 	+	"/static/"
)
type ShellController struct {
	beego.Controller
}
type shellInfoStruct  struct {
	Proto 	string 		`json:"proto"`
	IpAddr 	string 		`json:"ipaddr"`
	Port 	string 		`json:"port"`
	User 	string 		`json:"user"`
	Passwd 	string 		`json:"passwd"`
}
var upgrader = websocket.Upgrader{}
func (c *ShellController) Get() {
	fmt.Println(c.Ctx.Request.URL.Path)
	rethandle := func(err error,status int){
		reply := "ok"
		if err != nil {
			reply = err.Error()
		}
		c.Data["json"] = reply
		c.Ctx.Output.Status = status
		c.ServeJSON()
		return
	}
	switch c.Ctx.Request.URL.Path {
	case ShellUrl:
		proto := c.Ctx.Input.Query("proto")
		if proto == "" {
			rethandle(errors.New("Please Set Parameter 'proto'."),http.StatusOK)
			return
		}
		ipaddr := c.Ctx.Input.Query("ipaddr")
		if ipaddr == "" {
			rethandle(errors.New("Please Set Parameter 'ipaddr'."),http.StatusOK)
			return
		}
		port := c.Ctx.Input.Query("port")
		user := c.Ctx.Input.Query("user")
		passwd := c.Ctx.Input.Query("passwd")
		if proto == "ssh" && user == ""{
			rethandle(errors.New("Please Set Parameter 'user'."),http.StatusOK)
			return
		}
		shellinfo := shellInfoStruct{
			Proto:proto,
			IpAddr:ipaddr,
			Port:port,
			User:user,
			Passwd:passwd,
		}
		c.SetSession("keyinfo",shellinfo)
		c.TplName = "index.html"
	case ShellUrlJsJqueryMinJs,ShellUrlJsPopperMinJs,ShellUrlImgPng,
			ShellUrlJsBootstrapMinJs,ShellUrlJsXtermMinJs,ShellUrlJsFullScreenMinJs,
			ShellUrlCssBootstrapMinCss,ShellUrlCssXtermMinCss,ShellUrlCssFullScreenMinCss:
		handle := http.StripPrefix(ShellUrlShellStatic, http.FileServer(http.Dir("static")))
		handle.ServeHTTP(c.Ctx.ResponseWriter,c.Ctx.Request)
	}
}

func (c *ShellController) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ShellController) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ShellController) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

