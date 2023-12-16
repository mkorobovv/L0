package main

import (
	"log"

	"github.com/mkorobovv/L0/config"
	"github.com/mkorobovv/L0/internal/app"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error at starting: %v", err)
	}

	app.ServerRun(cfg)
}
