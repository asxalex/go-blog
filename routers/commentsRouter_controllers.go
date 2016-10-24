package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["myblog/controllers:PostController"] = append(beego.GlobalControllerRouter["myblog/controllers:PostController"],
		beego.ControllerComments{
			Method: "GetIndex",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["myblog/controllers:PostController"] = append(beego.GlobalControllerRouter["myblog/controllers:PostController"],
		beego.ControllerComments{
			Method: "GetSpecifiedTags",
			Router: `/tags/:tagname`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["myblog/controllers:PostController"] = append(beego.GlobalControllerRouter["myblog/controllers:PostController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/posts/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["myblog/controllers:PostController"] = append(beego.GlobalControllerRouter["myblog/controllers:PostController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/posts/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
