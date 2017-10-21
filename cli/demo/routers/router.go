package routers

import (
	"github.com/todaychiji/ha/cli/demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
