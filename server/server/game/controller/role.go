package controller

import "github.com/cr0ssover/three-kingdoms-game/server/net"

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.NewGroup("role")
	g.AddRouter("enterServer", r.enterServer)
}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 进入游戏的逻辑
	// 验证session合法性
}
