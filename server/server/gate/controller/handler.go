package controller

import (
	"log"
	"strings"
	"sync"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}

type Handler struct {
	proxyMutex sync.RWMutex
	// 代理地址 -> 客户端连接（客户端ID） -> 连接
	proxyMap   map[string]map[int64]*net.ProxyClient
	loginProxy string
	gameProxy  string
}

func (h *Handler) Router(r *net.Router) {
	h.loginProxy = config.File.Section("gate_proxy").Key("login_proxy").MustString("ws://127.0.0.1:8003")
	h.gameProxy = config.File.Section("gate_proxy").Key("game_proxy").MustString("ws://127.0.0.1:8001")
	g := r.NewGroup("*")
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 路由转发
	name := req.Body.Name
	proxyStr := ""
	if isAccount(name) {
		proxyStr = h.loginProxy
	}

	if proxyStr == "" {
		rsp.Body.Code = errcode.ProxyNotInConnect
		return
	}

	// 查询proxyMap中查询是否有连接
	h.proxyMutex.RLock()
	_, ok := h.proxyMap[proxyStr]
	if !ok {
		h.proxyMap[proxyStr] = make(map[int64]*net.ProxyClient)
	}
	h.proxyMutex.RUnlock()

	cidValue, err := req.Conn.GetProperty("cid")
	if err != nil {
		log.Println("获取cid失败,err: ", err)
		rsp.Body.Code = errcode.InvalidParam
		return
	}
	cid := cidValue.(int64)
	proxy, ok := h.proxyMap[proxyStr][cid]
	if !ok {
		proxy = net.NewProxyClient(proxyStr)
		h.proxyMutex.Lock()
		h.proxyMap[proxyStr][cid] = proxy
		h.proxyMutex.Unlock()
		err := proxy.Connect()
		if err != nil {
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxyStr], cid)
			h.proxyMutex.Unlock()
			rsp.Body.Code = errcode.ProxyConnectError
			log.Println(err)
			return
		}
		proxy.SetProperty("cid", cid)
		proxy.SetProperty("proxy", proxyStr)
		proxy.SetProperty("gateConn", req.Conn)
		proxy.SetOnPush(h.onPush)
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	r, err := proxy.Send(req.Body.Name, req.Body.Msg)
	if err != nil {
		log.Println("发送异常，err: ", err)
		rsp.Body.Code = errcode.ProxyConnectError
		rsp.Body.Msg = nil
	}
	rsp.Body.Code = r.Code
	rsp.Body.Msg = r.Msg
}

func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}

func (h *Handler) onPush(conn *net.ClientConn, body *net.RspBody) {
	gc, err := conn.GetProperty("gateConn")
	if err != nil {
		log.Println("onPush gateConn,err: ", err)
		return
	}
	gateConn := gc.(net.WsConner)
	gateConn.Push(body.Name, body.Msg)
}
