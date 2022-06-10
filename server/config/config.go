package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"gopkg.in/ini.v1"
)

var configFile = "/conf/conf.ini"

var File *ini.File

func init() {
	var (
		projectPath   string
		configPath    string
		err           error
		flagConfParam string
	)
	// 读取环境变量
	projectPath = os.Getenv("TKG_PROJECT_PAHT")
	if projectPath == "" {
		_, projectPath, _, _ = runtime.Caller(0)
		projectPath, err = filepath.Abs(filepath.Dir(filepath.Dir(projectPath)))
		if err != nil {
			panic(err)
		}
		err = os.Setenv("TKG_PROJECT_PAHT", projectPath)
		if err != nil {
			logger.Error("环境变量设置异常")
		}
	}
	configPath = projectPath + configFile
	// 解析命令行参数
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
