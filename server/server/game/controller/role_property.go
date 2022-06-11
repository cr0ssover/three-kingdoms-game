package controller

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/logic"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var DefaultRoleProperty = &roleProperty{}

type roleProperty struct {
}

func (r *roleProperty) Router(router *net.Router) {
	g := router.NewGroup("role")
	g.AddRouter("myProperty", r.myProperty)
}

func (r *roleProperty) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	var err error
	// 加载 资源 城池 建筑 部队 武将
	propertyRsp := &model.MyRolePropertyRsp{}
	// 获取角色
	roleValue, err := req.Conn.GetProperty("role")
	if err != nil {
		logger.Warn("获取角色失败,err: ", err)
	}
	role := roleValue.(*data.Role)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	// 获取城池
	propertyRsp.Citys, err = logic.RoleCityService.GetRoleCitys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 获取建筑
	propertyRsp.MRBuilds, err = logic.RoleBuildService.GetRoleBuilds(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 获取角色资源
	propertyRsp.RoleRes, err = logic.RoleService.GetRoleRes(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 获取武将
	propertyRsp.Generals, err = logic.GeneralService.GetGenerals(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 获取军队
	propertyRsp.Armys, err = logic.ArmyService.GetArmys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rsp.Body.Code = errcode.OK
	rsp.Body.Msg = propertyRsp
}
