package app

import (
	"fmt"
	"log"
	"net/http"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/httpRouter"
)

func Run(config config.Config) {
	// some creating
	router := httpRouter.NewRouter(nil, nil, nil)
	addr := fmt.Sprintf(":%d", config.HTTPPort)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
