package controller

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/logic"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/mitchellh/mapstructure"
)

var DefaultRoleController = &roleController{}

type roleController struct {
}

func (r *roleController) Router(router *net.Router) {
	g := router.NewGroup("role")
	g.AddRouter("enterServer", r.enterServer)
}

func (r *roleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 进入游戏的逻辑
	// 验证session合法性
	serverReq := &model.EnterServerReq{}
	serverRsp := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, serverReq)
	if err != nil {
		logger.Warn("参数有误,err", err)
		rsp.Body.Code = errcode.InvalidParam
		return
	}

	// 获取session
	seesion := serverReq.Session
	// 解析session
	_, claim, err := utils.ParseToken(seesion)
	// 解析失败，seesion不合法
	if err != nil {
		logger.Warn("seesion不合法")
		rsp.Body.Code = errcode.SessionInvalid
		return
	}

	uid := claim.Uid
	// 初始化玩家数据
	err = logic.RoleService.EnterServer(uid, serverRsp, req.Conn)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = errcode.OK
	rsp.Body.Msg = serverRsp
}
