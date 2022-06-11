package logic

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	gameConfig "github.com/cr0ssover/three-kingdoms-game/server/server/game/config"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
)

var RoleService = &roleService{}

type roleService struct{}

func (r *roleService) EnterServer(uid int, serverRsp *model.EnterServerRsp, conn net.WsConner) error {
	role := &data.Role{}
	// 根据uid查询角色信息
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		logger.Warn("查询角色出错,err", err)
		return common.New(errcode.DBError, "查询角色出错")
	}
	if !ok {
		logger.Warn("角色不存在")
		return common.New(errcode.DBError, "角色不存在")
	}

	// 查询角色资源
	rid := role.RId
	roleRes := &data.RoleRes{}
	ok, err = db.Engine.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
	if err != nil {
		logger.Warn("查询角色资源出错,err", err)
		return common.New(errcode.DBError, "查询角色资源出错")
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
			return common.New(errcode.DBError, "角色资源插入数据失败")
		}
	}

	// 给前端返回的数据赋值
	serverRsp.RoleRes = roleRes.ToModel()
	serverRsp.Role = role.ToModel()
	serverRsp.Time = time.Now().UnixNano() / 1e6
	token, _ := utils.Award(rid)
	serverRsp.Token = token
	conn.SetProperty("role", role)

	// 初始化玩家角色属性
	err = RoleAttrService.Init(serverRsp.Role.RId)
	if err != nil {
		return err
	}

	// 初始化城池
	err = RoleCityService.Init(role)
	if err != nil {
		return err
	}
	return nil
}

// 获取角色资源
func (*roleService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
	if err != nil {
		return model.RoleRes{}, common.New(errcode.DBError, "查询角色资源出错")
	}
	if ok {
		return roleRes.ToModel(), nil
	}
	return model.RoleRes{}, nil
}
