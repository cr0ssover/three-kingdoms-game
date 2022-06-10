package main

import (
	"fmt"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/game"
)

func main() {
	host := config.File.Section("game_proxy").Key("host").MustString("127.0.0.1")
	port := config.File.Section("game_proxy").Key("port").MustString("8001")

	s := net.NewServer(fmt.Sprintf("%s:%s", host, port))
	game.Init()
	s.Router(game.Router)
	s.Start()
}
