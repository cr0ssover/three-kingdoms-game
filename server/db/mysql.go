package db

import (
	"fmt"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func InitDB() {
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		logger.Error("读取数据库配置出错")
		panic(err)
	}
	host := mysqlConfig.Key("host").MustString("localhost")
	prot := mysqlConfig.Key("port").MustString("3306")
	username := mysqlConfig.Key("username").Value()
	password := mysqlConfig.Key("password").Value()
	dbname := mysqlConfig.Key("dbname").MustString("tkg")
	// 数据库空闲连接数
	maxIdle := mysqlConfig.Key("max_idle").MustInt(2)
	// 数据库最大连接数
	maxConn := mysqlConfig.Key("max_conn").MustInt(10)
	if maxIdle > maxConn {
		panic("数据库最大空闲连接数不能大于最大连接数")
	}
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		host,
		prot,
		dbname,
	)
	Engine, err = xorm.NewEngine("mysql", dbConn)
	if err != nil {
		logger.Error("数据库连接失败")
		panic(err)
	}
	err = Engine.Ping()
	if err != nil {
		logger.Error("数据库连接测试失败")
		panic(err)
	}
	Engine.SetMaxIdleConns(maxIdle)
	Engine.SetMaxOpenConns(maxConn)
	Engine.ShowSQL(true)
}
