package data

// 角色资源表
type RoleRes struct {
	Id     int `xorm:"int(0) 'id' pk autoincr notnull comment('id')"`
	RId    int `xorm:"int(0) 'r_id' notnull comment('角色id')"`
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
