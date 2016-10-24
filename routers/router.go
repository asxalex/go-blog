package routers

import (
	"github.com/astaxie/beego"
	"myblog/controllers"
)

func init() {
	//beego.SetStaticPath("/static", "static")
	//beego.StaticDir["static"] = "static/"
	//beego.Router("/", &controllers.MainController{})
	//beego.Router("/posts/", &controllers.PostController{})
	beego.Include(&controllers.PostController{})
}
