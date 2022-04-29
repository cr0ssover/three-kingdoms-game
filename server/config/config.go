package config

import (
	"errors"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

var configFile = "/conf/conf.ini"

var File *ini.File

func init() {
	// 获取获取当前路径
	crruentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := crruentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
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
