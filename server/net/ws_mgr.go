package net

import (
	"sync"
)

var Mgr = &WsMgr{
	userCache: make(map[int]WsConner),
}

type WsMgr struct {
	mux       sync.RWMutex
	userCache map[int]WsConner
}

func (m *WsMgr) UserLogin(conn WsConner, uid int, token string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	oldConn := m.userCache[uid]
	if oldConn != nil {
		if conn != oldConn {
			//通过旧客户端 有用户抢登录了
			oldConn.Push("robLogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
