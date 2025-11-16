package main

import (
	"log"
	"pr-reviewer/internal/app"
	"pr-reviewer/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(*cfg)
}
