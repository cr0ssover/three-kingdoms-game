package login

import (
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login/controller"
)

// 初始化路由
var Router = net.NewRouter()

func Init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
