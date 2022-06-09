package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web"
	"github.com/gin-gonic/gin"
)

func main() {
	host := config.File.Section("web_server").Key("host").MustString("127.0.0.1")
	port := config.File.Section("web_server").Key("port").MustString("8088")

	router := gin.Default()
	web.Init(router)

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", host, port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

}
