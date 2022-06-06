package model

import "time"

const (
	Login int8 = iota
	Logout
)

type LoginHistory struct {
	Id       int       `xorm:"int(0) pk autoincr 'id'"`
	UId      int       `xorm:"int(0) 'uid' notnull comment('用户UID')"`
	CTime    time.Time `xorm:"timestamp(0) 'ctime' notnull comment('登录时间')"`
	Ip       string    `xorm:"varchar(31) 'ip' notnull comment('ip')"`
	State    int8      `xorm:"int 'state' notnull default(0) comment('登录状态：0-登录 1-登出')"`
	Hardware string    `xorm:"varchar(64) 'hardware' notnull default('') comment('硬件信息')"`
}

type LoginLast struct {
	Id         int       `xorm:"int pk autoincr 'id'"`
	UId        int       `xorm:"int(0) 'uid' notnull comment('用户UID')"`
	LoginTime  time.Time `xorm:"timestamp(0) 'login_time' default(null) comment('登录时间')"`
	LogoutTime time.Time `xorm:"timestamp(0) 'logout_time' default(null) comment('登出时间')"`
	Ip         string    `xorm:"varchar(31) 'ip' notnull comment('ip')"`
	Session    string    `xorm:"varchar(255) 'session' default(null) comment('会话') "`
	IsLogout   int8      `xorm:"tinyint(0) 'is_logout' notnull default(0) comment('是否登出：1-logout 0-login')"`
	Hardware   string    `xorm:"varchar(64) 'hardware' notnull default('') comment('硬件信息')"`
}

func (*LoginHistory) TableName() string {
	return "login_history"
}

func (*LoginLast) TableName() string {
	return "login_last"
}
