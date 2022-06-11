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
	// 加载地图配置
	gameConfig.MapBuildConfig.Load()
	initRouter()
}

func initRouter() {
	// 角色路由
	controller.DefaultRoleController.Router(Router)
	// 地图路由
	controller.DefaultNationMapController.Router(Router)
	// 角色属性
	controller.DefaultRoleProperty.Router(Router)
}
