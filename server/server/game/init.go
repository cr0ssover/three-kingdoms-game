package game

import (
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	gameConfig "github.com/cr0ssover/three-kingdoms-game/server/server/game/config"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/controller"
)

var Router = &net.Router{}

func Init() {
	db.InitDB()
	// 加载基础配置
	gameConfig.Base.Load()
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
}
