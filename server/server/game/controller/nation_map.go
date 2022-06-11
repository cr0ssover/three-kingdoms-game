package controller

import (
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	gameConfig "github.com/cr0ssover/three-kingdoms-game/server/server/game/config"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
)

var DefaultNationMapController = &nationMapController{}

type nationMapController struct {
}

func (r *nationMapController) Router(router *net.Router) {
	g := router.NewGroup("nationMap")
	g.AddRouter("config", r.config)

}

// 初始化地图配置
func (r *nationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	configRsp := &model.ConfigRsp{}
	cfgs := gameConfig.MapBuildConfig.Cfg
	configRsp.Confs = make([]model.Conf, len(cfgs))
	for i, v := range cfgs {
		configRsp.Confs[i].Type = v.Type
		configRsp.Confs[i].Name = v.Name
		configRsp.Confs[i].Level = v.Level
		configRsp.Confs[i].Defender = v.Defender
		configRsp.Confs[i].Durable = v.Durable
		configRsp.Confs[i].Grain = v.Grain
		configRsp.Confs[i].Iron = v.Iron
		configRsp.Confs[i].Stone = v.Stone
		configRsp.Confs[i].Wood = v.Wood
	}
	rsp.Body.Code = errcode.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Msg = configRsp
	rsp.Body.Name = req.Body.Name
}
