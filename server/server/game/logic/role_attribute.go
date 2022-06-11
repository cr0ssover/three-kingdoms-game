package logic

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var RoleAttrService = &roleAttrService{}

type roleAttrService struct{}

func (r *roleAttrService) Init(rid int) error {
	roleAttr := &data.RoleAttribute{}
	ok, err := db.Engine.Table(roleAttr).Where("rid=?", rid).Get(roleAttr)
	if err != nil {
		logger.Warn("角色属性查询失败,err:", err)
		return common.New(errcode.DBError, "角色属性查询失败")
	}

	if !ok {
		// 没有数据,进行初始化
		roleAttr.RId = rid
		roleAttr.ParentId = 0
		roleAttr.UnionId = 0
		_, err = db.Engine.Table(roleAttr).Insert(roleAttr)
		if err != nil {
			logger.Warn("角色属性数据插入失败,err", err)
			return common.New(errcode.DBError, "角色属性数据插入失败")
		}
	}

	return nil
}
