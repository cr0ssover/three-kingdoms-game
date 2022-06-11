package gameconfig

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cr0ssover/three-kingdoms-game/server/logger"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConfig = &mapBuildConf{}

// 基础配置文件路径
const mapBuildConfFile string = "/conf/game/map_build.json"

func (m *mapBuildConf) Load() {
	configPath := projectPath + mapBuildConfFile
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error("基础资源配置读取异常,err: ", err)
		panic(data)
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		logger.Warn("基础资源配置json格式有误,err: ", err)
		panic(err)
	}

}
