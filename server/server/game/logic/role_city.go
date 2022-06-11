package logic

import (
	"math/rand"
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	gameConfig "github.com/cr0ssover/three-kingdoms-game/server/server/game/config"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model/data"
)

var RoleCityService = &roleCityService{}

type roleCityService struct{}

// 初始化角色城池
func (r *roleCityService) Init(role *data.Role) error {
	roleCity := &data.MapRoleCity{}
	ok, err := db.Engine.Table(roleCity).Where("rid = ?", role.RId).Get(roleCity)
	if err != nil {
		logger.Warn("角色城池查询失败,err:", err)
		return common.New(errcode.DBError, "角色城池查询失败")
	}

	if !ok {
		// 没有数据,进行初始化
		roleCity.X = rand.Intn(MapWith)  // 随机横坐标
		roleCity.Y = rand.Intn(MapHight) //随机纵坐标
		// *TODO 城池初始化时N范围内不能存在其他城池

		// *
		roleCity.RId = role.RId
		roleCity.CreatedAt = time.Now()
		roleCity.CurDurable = gameConfig.Base.City.Durable
		roleCity.Name = role.Nickname
		_, err = db.Engine.Table(roleCity).Insert(roleCity)
		if err != nil {
			logger.Warn("角色城池数据插入失败,err", err)
			return common.New(errcode.DBError, "角色城池数据插入失败")
		}
		// *TODO 城池基础设施初始化

		// *
	}

	return nil
}

func (r *roleCityService) GetRoleCitys(rid int) ([]model.MapRoleCity, error) {
	roleCitys := make([]data.MapRoleCity, 0)
	roleCity := &data.MapRoleCity{}
	err := db.Engine.Table(roleCity).Where("rid = ?", rid).Find(&roleCitys)
	if err != nil {
		logger.Warn("角色城池查询失败,err:", err)
		return nil, common.New(errcode.DBError, "角色城池查询失败")
	}
	mRoleCitys := make([]model.MapRoleCity, len(roleCitys))
	for _, v := range roleCitys {
		mRoleCitys = append(mRoleCitys, v.ToModel())
	}
	return mRoleCitys, nil
}
