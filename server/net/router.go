package net

import (
	"strings"

	"github.com/cr0ssover/three-kingdoms-game/server/logger"
)

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type Router struct {
	group []*group
}

// 新建路由组
func (r *Router) NewGroup(prefix string) *group {
	for _, v := range r.group {
		if v.prefix == prefix {
			return v
		}
	}
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group, g)
	return g
}

func NewRouter() *Router {
	return &Router{}
}

// 注册路由
func (g *group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	// 路径 登录业务 account.login (account组标识)login 路由标识
	strs := strings.Split(req.Body.Name, ".")

	var (
		prefix string
		name   string
	)
	if len(strs) == 2 {
		prefix = strs[0] //前缀
		name = strs[1]   //路由名称
	}
	for _, g := range r.group {
		// 判断路由前缀是否相等，或是否为全放行
		if g.prefix == prefix || g.prefix == "*" {
			g.exec(name, req, rsp)
			break
		}
	}
}

// 执行处理函数
func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	} else {
		h = g.handlerMap["*"]
		if h != nil {
			h(req, rsp)
		} else {
			logger.Warnf("%s 路由未定义", name)
		}
	}
}
