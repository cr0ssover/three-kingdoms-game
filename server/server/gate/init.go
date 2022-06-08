package gate

import (
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/gate/controller"
)

var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	controller.GateHandler.Router(Router)
}
