package data

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
)

// 玩家表
type Role struct {
	RId        int       `xorm:"int(0) 'rid' pk autoincr comment('roleId')"`
	UId        int       `xorm:"int(0) 'uid' notnull comment('用户UID')"`
	Nickname   string    `xorm:"varchar(100) 'nickname'  comment('昵称')" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Balance    int       `xorm:"int(0) 'balance' default('0') comment('余额')"`
	HeadId     int16     `xorm:"int 'headId' notnull   comment('头像ID')"`
	Sex        int8      `xorm:"int(0) 'sex' notnull   comment('性别 0-女 1-男')"`
	Profile    string    `xorm:"varchar(500) 'profile' notnull  comment('个人简介')"`
	LoginTime  time.Time `xorm:"timestamp(0) 'login_time' comment('登录时间')"`
	LogoutTime time.Time `xorm:"timestamp(0) 'logout_time' comment('登出时间')"`
	CreatedAt  time.Time `xorm:"timestamp(0) 'created_at' notnull  default('CURRENT_TIMESTAMP(0)')  comment('角色创建时间')"`
}

func (r *Role) TableName() string {
	return "role"
}

func (r *Role) ToModel() model.Role {
	p := model.Role{}
	p.Balance = r.Balance
	p.HeadId = r.HeadId
	p.Nickname = r.Nickname
	p.Profile = r.Profile
	p.RId = r.RId
	p.Sex = r.Sex
	p.UId = r.UId
	return p
}
