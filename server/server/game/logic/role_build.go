package logic

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var RoleBuildService = &roleBuildService{}

type roleBuildService struct{}

func (r *roleBuildService) GetRoleBuilds(rid int) ([]model.MapRoleBuild, error) {
	roleBuilds := make([]data.MapRoleBuild, 0)
	roleBuild := &data.MapRoleBuild{}
	err := db.Engine.Table(roleBuild).Where("rid = ?", rid).Find(&roleBuilds)
	if err != nil {
		logger.Warn("建筑查询失败,err:", err)
		return nil, common.New(errcode.DBError, "建筑查询失败")
	}
	mRoleBuilds := make([]model.MapRoleBuild, len(roleBuilds))
	for _, v := range roleBuilds {
		mRoleBuilds = append(mRoleBuilds, v.ToModel())
	}
	return mRoleBuilds, nil
}
