package logic

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var GeneralService = &generalService{}

type generalService struct{}

func (g *generalService) GetGenerals(rid int) ([]model.General, error) {
	generals := make([]data.General, 0)
	general := &data.General{}
	err := db.Engine.Table(general).Where("rid = ?", rid).Find(&generals)
	if err != nil {
		logger.Warn("角色属性查询失败,err:", err)
		return nil, common.New(errcode.DBError, "角色属性查询失败")
	}
	mGenerals := make([]model.General, 0)
	for _, v := range generals {
		mGenerals = append(mGenerals, v.ToModel())
	}
	return mGenerals, nil
}
