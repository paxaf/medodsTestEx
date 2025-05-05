package main

import (
	"fmt"
	"log"

	"github.com/paxaf/medodsTestEx/config"
	"github.com/paxaf/medodsTestEx/internal/app"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err, "failed to load cfg")
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err, "error creating app")
	}
	fmt.Println(cfg.APIServer)
	if err = app.Run(); err != nil {
		log.Fatal(err, "error running app")
	}
}
