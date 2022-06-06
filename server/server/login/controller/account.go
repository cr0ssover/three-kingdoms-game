package controller

import (
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login/proto"
	"github.com/mitchellh/mapstructure"
)

var DefaultAccount = &Account{}

type Account struct{}

func (a *Account) Router(r *net.Router) {
	g := r.NewGroup("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	/**
	1.用户 密码 硬件ID
	2.根据用户名查询User表得到数据
	3.进行密码比对
	4.保存用户登录记录
	5.保存用户最后一次登录信息
	6.客户端需要一个session,jwt
	*/
	loginReq := &proto.LoginReq{}
	// loginRsp := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{}
	loginRes.UID = 1
	loginRes.Username = "admin"
	loginRes.Session = "sss"
	loginRes.Password = ""
	rsp.Body.Msg = loginRes
}
