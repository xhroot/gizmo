package main

import (
	"github.com/xhroot/gizmo/examples/servers/simple/service"

	"github.com/xhroot/gizmo/config"
	"github.com/xhroot/gizmo/server"
)

func main() {
	// showing 1 way of managing gizmo/config: importing from the environment
	cfg := service.Config{Server: &config.Server{}}
	config.LoadEnvConfig(&cfg)
	config.LoadEnvConfig(cfg.Server)

	server.Init("nyt-simple-proxy", cfg.Server)

	err := server.Register(service.NewSimpleService(&cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
