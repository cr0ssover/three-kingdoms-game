package web

import (
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web/controller"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	// 初始化数据库
	db.InitDB()
	// 处理跨域
	router.Use(middleware.Cors())
	initRouter(router)

}

func initRouter(router *gin.Engine) {
	router.GET("/account/register", controller.DefaultAccountConterller.Register)
}
