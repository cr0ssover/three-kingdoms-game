package controller

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	gameConfig "github.com/cr0ssover/three-kingdoms-game/server/server/game/config"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/mitchellh/mapstructure"
)

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
	serverReq := &model.EnterServerReq{}
	serverRsp := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, serverReq)
	if err != nil {
		logger.Warn("参数有误,err", err)
		rsp.Body.Code = errcode.InvalidParam
		return
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
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
	role := &data.Role{}
	// 根据uid查询角色信息
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		logger.Warn("查询角色出错,err", err)
		rsp.Body.Code = errcode.DBError
		return
	}
	if !ok {
		logger.Warn("角色不存在")
		rsp.Body.Code = errcode.RoleNotExist
		return
	}

	// 查询角色资源
	rid := role.RId
	roleRes := &data.RoleRes{}
	ok, err = db.Engine.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
	if err != nil {
		logger.Warn("查询角色资源出错,err", err)
		rsp.Body.Code = errcode.DBError
		return
	}
	if !ok {
		// 如果查询数据为空，则初始化角色资源数据
		roleRes.RId = rid
		roleRes.Gold = gameConfig.Base.Role.Gold
		roleRes.Decree = gameConfig.Base.Role.Decree
		roleRes.Grain = gameConfig.Base.Role.Grain
		roleRes.Iron = gameConfig.Base.Role.Iron
		roleRes.Stone = gameConfig.Base.Role.Stone
		roleRes.Wood = gameConfig.Base.Role.Wood
		if _, err := db.Engine.Table(roleRes).Insert(roleRes); err != nil {
			logger.Warn("角色资源插入数据失败,err", err)
			rsp.Body.Code = errcode.DBError
			return
		}
	}

	// 给前端返回的数据
	serverRsp.RoleRes = roleRes.ToModel()
	serverRsp.Role = role.ToModel()
	serverRsp.Time = time.Now().UnixNano() / 1e6
	token, _ := utils.Award(rid)
	serverRsp.Token = token
	rsp.Body.Code = errcode.OK
	rsp.Body.Msg = serverRsp
	logger.Debug(rsp, rsp.Body.Msg)
}
