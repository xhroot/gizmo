package main

import (
	"github.com/xhroot/gizmo/examples/pubsub/api-sns-pub/service"

	"github.com/xhroot/gizmo/config"
	svr "github.com/xhroot/gizmo/server"
)

func main() {
	// showing 1 way of managing gizmo/config: importing from a local file
	cfg := config.NewConfig("./config.json")

	svr.Init("nyt-json-pub-proxy", cfg.Server)

	err := svr.Register(service.NewJSONPubService(cfg))
	if err != nil {
		svr.Log.Fatal("unable to register service: ", err)
	}

	err = svr.Run()
	if err != nil {
		svr.Log.Fatal("server encountered a fatal error: ", err)
	}
}
