package logic

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var ArmyService = &armyService{}

type armyService struct{}

func (r *armyService) GetArmys(rid int) ([]model.Army, error) {
	armys := make([]data.Army, 0)
	army := &data.Army{}
	err := db.Engine.Table(army).Where("rid = ?", rid).Find(armys)
	if err != nil {
		logger.Warn("军队查询失败,err:", err)
		return nil, common.New(errcode.DBError, "军队查询失败")
	}
	mArmys := make([]model.Army, len(armys))
	for _, v := range armys {
		mArmys = append(mArmys, v.ToModel())
	}
	return mArmys, nil
}
