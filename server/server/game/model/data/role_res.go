package data

import "github.com/cr0ssover/three-kingdoms-game/server/server/game/model"

// 角色资源表
type RoleRes struct {
	Id     int `xorm:"int(0) 'id' pk autoincr notnull comment('id')"`
	RId    int `xorm:"int(0) 'rid' notnull comment('角色id')"`
	Wood   int `xorm:"int(0) 'wood' notnull comment('木')"`
	Iron   int `xorm:"int(0) 'iron' notnull comment('铁')"`
	Stone  int `xorm:"int(0) 'stone' notnull comment('石头')"`
	Grain  int `xorm:"int(0) 'grain' notnull comment('粮食')"`
	Gold   int `xorm:"int(0) 'gold' notnull comment('金币')"`
	Decree int `xorm:"int(0) 'decree' notnull comment('令牌')"`
}

func (r *RoleRes) TableName() string {
	return "role_res"
}

func (r *RoleRes) ToModel() model.RoleRes {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	p.GoldYield = 1
	p.GrainYield = 1
	p.StoneYield = 1
	p.IronYield = 1
	p.WoodYield = 1
	return p
}
