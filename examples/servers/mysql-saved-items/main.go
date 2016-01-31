package main

import (
	"flag"

	"github.com/xhroot/gizmo/config"
	"github.com/xhroot/gizmo/server"
	_ "github.com/go-sql-driver/mysql"

	"github.com/xhroot/gizmo/examples/servers/mysql-saved-items/service"
)

func main() {
	// load from the local JSON file into a config.Config struct
	cfg := config.NewConfig("./config.json")
	flag.Parse()
	// SetServerOverrides will allow us to override some of the values in
	// the JSON file with CLI flags.
	config.SetServerOverrides(cfg.Server)

	// initialize Gizmo’s server with given configs
	server.Init("nyt-saved-items", cfg.Server)

	// instantiate a new ‘saved items service’ with our MySQL credentials
	svc, err := service.NewSavedItemsService(cfg.MySQL)
	if err != nil {
		server.Log.Fatal("unable to create saved items service: ", err)
	}

	// register our saved item service with the Gizmo server
	err = server.Register(svc)
	if err != nil {
		server.Log.Fatal("unable to register saved items service: ", err)
	}

	// run the Gizmo server
	err = server.Run()
	if err != nil {
		server.Log.Fatal("unable to run saved items service: ", err)
	}
}
