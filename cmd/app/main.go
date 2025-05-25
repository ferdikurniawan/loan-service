package main

import (
	"log"

	"github.com/ferdikurniawan/loan-service/config"
	"github.com/ferdikurniawan/loan-service/internal/app"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	app.Run(cfg)
}
