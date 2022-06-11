package data

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
)

const (
	MapBuildSysFortress = 50 //系统要塞
	MapBuildSysCity     = 51 //系统城市
	MapBuildFortress    = 56 //玩家要塞
)

// 角色建筑表
type MapRoleBuild struct {
	Id         int       `xorm:"id pk autoincr"`
	RId        int       `xorm:"rid comment('角色ID')"`
	Type       int8      `xorm:"type comment('建筑类型')"`
	Level      int8      `xorm:"level comment('建筑等级')"`
	OPLevel    int8      `xorm:"op_level comment('建筑操作等级')"` //操作level
	X          int       `xorm:"x comment('x坐标')"`
	Y          int       `xorm:"y"`
	Name       string    `xorm:"name comment('名称') "`
	Wood       int       `xorm:"-"`
	Iron       int       `xorm:"-"`
	Stone      int       `xorm:"-"`
	Grain      int       `xorm:"-"`
	Defender   int       `xorm:"-"`
	CurDurable int       `xorm:"cur_durable comment('当前耐久') "`
	MaxDurable int       `xorm:"max_durable comment('最大耐久') "`
	OccupyTime time.Time `xorm:"occupy_time comment('占领时间') "`
	EndTime    time.Time `xorm:"end_time comment('建造、升级、拆除结束时间') "` //建造或升级完的时间
	GiveUpTime int64     `xorm:"give_up_time comment('放弃时间') "`
}

func (m *MapRoleBuild) TableName() string {
	return "map_role_build"
}

func (m *MapRoleBuild) ToModel() model.MapRoleBuild {
	p := model.MapRoleBuild{}
	p.RId = m.RId
	p.Type = m.Type
	p.Level = m.Level
	p.OPLevel = m.OPLevel
	p.X = m.X
	p.Y = m.Y
	p.Name = m.Name
	p.CurDurable = m.CurDurable
	p.MaxDurable = m.MaxDurable
	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
	p.EndTime = m.EndTime.UnixNano() / 1e6
	p.GiveUpTime = m.GiveUpTime

	p.RNick = "111"
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0

	return p
}
