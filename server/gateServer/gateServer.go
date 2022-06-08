package main

import (
	"fmt"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/gate"
)

func main() {
	host := config.File.Section("gate_proxy").Key("host").MustString("127.0.0.1")
	port := config.File.Section("gate_proxy").Key("port").MustString("8004")

	s := net.NewServer(fmt.Sprintf("%s:%s", host, port))
	s.NeedSecret(true)
	gate.Init()
	s.Router(gate.Router)
	s.Start()
}
