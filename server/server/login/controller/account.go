package controller

import (
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct{}

func (a *Account) Router(r *net.Router) {
	g := r.NewGroup("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{}
	loginRes.UID = 1
	loginRes.Username = "admin"
	loginRes.Session = "sss"
	loginRes.Password = ""
	rsp.Body.Msg = loginRes
}
