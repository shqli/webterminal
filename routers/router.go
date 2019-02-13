package routers

import (
	"webterminal/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router(controllers.ShellUrl,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlImgPng,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlJsJqueryMinJs,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlJsPopperMinJs,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlJsBootstrapMinJs,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlJsXtermMinJs,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlJsFullScreenMinJs,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlCssBootstrapMinCss,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlCssXtermMinCss,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlCssFullScreenMinCss,&controllers.ShellController{})
	beego.Router(controllers.ShellUrlShellWs, &controllers.ShellWsController{})
	beego.Router("/api/*",&controllers.ApiDeny{})
}
