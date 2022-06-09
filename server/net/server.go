package net

import (
	"net/http"

	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/gorilla/websocket"
)

type server struct {
	addr       string
	router     *Router
	needSecret bool
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *Router) {
	s.router = router
}

// 需要加密
func (s *server) NeedSecret(needSecret bool) {
	s.needSecret = needSecret
}

// 启动服务
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http 升级websocket协议配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warn("websocket服务连接出错:", err)
	}
	logger.Info("websocket服务连接成功")
	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}
