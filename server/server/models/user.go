package models

import "time"

type User struct {
	UId      int       `xorm:"int(0) 'uid'  pk autoincr"`
	Username string    `xorm:"varchar(20) 'username' notnull unique comment('用户名')"`
	Passcode string    `xorm:"char(12) 'passcode' notnull unique comment('加密随机数')"`
	Passwd   string    `xorm:"char(64) 'passwd' notnull comment('密码')"`
	Status   int       `xorm:"tinyint(0) 'status' notnull default(0) comment('用户账户状态 0-默认 1-冻结 2-停号')"`
	Hardware string    `xorm:"varchar(64) 'hardware' notnull default('') comment('硬件信息')"`
	Ctime    time.Time `xorm:"timestamp(0) 'ctime' notnull default('2022-06-01 08:00:00')"`
	Mtime    time.Time `xorm:"timestamp(0) 'mtime' notnull default('2022-06-01 08:00:00')"`
	IsOnline bool      `xorm:"-"`
}
