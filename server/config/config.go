package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/ini.v1"
)

var configFile = "/conf/conf.ini"

var File *ini.File

func init() {
	// 获取获取配置路径
	_, projectPath, _, _ := runtime.Caller(0)
	projectPath, err := filepath.Abs(filepath.Dir(filepath.Dir(projectPath)))
	if err != nil {
		panic(err)
	}
	configPath := projectPath + configFile

	// 解析命令行参数
	var flagConfParam string
	flag.StringVar(&flagConfParam, "conf", configPath, "conf")
	testing.Init()
	flag.Parse()
	configPath = flagConfParam
	if err != nil {
		panic(err)
	}
	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	// 加载配置文件
	File, err = ini.Load(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
