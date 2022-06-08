package main

import (
	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login"
)

func main() {
	host := config.File.Section("server").Key("host").MustString("127.0.0.1")
	port := config.File.Section("server").Key("port").MustString("8003")
	s := net.NewServer(host + ":" + port)
	s.NeedSecret(false)
	login.Init()
	s.Router(login.Router)
	s.Start()
}
