package data

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
)

type MapRoleCity struct {
	// mutex      sync.Mutex `xorm:"-"`
	CityId     int       `xorm:"int(0) pk autoincr 'city_id' comment('城池ID')"`
	RId        int       `xorm:"int(0) 'rid' comment('角色id')"`
	Name       string    `xorm:"varchar(100) 'name' comment('城池名称')" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	X          int       `xorm:"int(0) 'x' comment('x坐标')"`
	Y          int       `xorm:"int(0) 'y' comment('y坐标')"`
	IsMain     int8      `xorm:"int(0) 'is_main' comment('是否是主城 0- 1-')"`
	CurDurable int       `xorm:"int(0) 'cur_durable' comment('当前耐久')"`
	CreatedAt  time.Time `xorm:"timestamp(0) 'created_at' comment('城池创建时间')"`
	OccupyTime time.Time `xorm:"timestamp(0) 'occupy_time'  comment('占领时间')"`
}

func (*MapRoleCity) TableName() string {
	return "map_role_city"
}

func (m *MapRoleCity) ToModel() model.MapRoleCity {
	p := model.MapRoleCity{}
	p.CityId = m.CityId
	p.RId = m.RId
	p.Name = m.Name
	p.X = m.X
	p.Y = m.Y
	p.IsMain = (m.IsMain == 1)
	p.CurDurable = m.CurDurable
	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6

	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.Level = 1
	return p
}
