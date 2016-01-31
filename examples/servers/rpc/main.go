package main

import (
	"github.com/xhroot/gizmo/examples/servers/rpc/service"

	"github.com/xhroot/gizmo/config"
	"github.com/xhroot/gizmo/server"
)

func main() {
	// showing 1 way of managing gizmo/config: importing from a local file
	var cfg *service.Config
	config.LoadJSONFile("./config.json", &cfg)

	server.Init("nyt-rpc-proxy", cfg.Server)

	err := server.Register(service.NewRPCService(cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
