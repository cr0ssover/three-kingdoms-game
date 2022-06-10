package data

import "time"

// 玩家表
type Role struct {
	RId        int       `xorm:"int(0) 'r_id' pk autoincr comment('roleId')"`
	UId        int       `xorm:"int(0) 'u_id' notnull comment('用户UID')"`
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
