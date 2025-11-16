package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/httprouter"
	"pr-reviewer/internal/repo/postgres"
	"pr-reviewer/internal/service"
	"syscall"
	"time"
)

func Run(config config.Config) {
	uow, err := postgres.NewPostgresUoW(config.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer uow.Close()

	teamSvc := service.NewTeamService(uow)
	userSvs := service.NewUserService(uow)
	prSvs := service.NewPRService(uow)

	router := httprouter.NewRouter(teamSvc, userSvs, prSvs)
	addr := fmt.Sprintf(":%d", config.HTTPPort)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Println("starting server")
		err = server.ListenAndServe()
		log.Println("server error: ", err)
		cancel()
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	<-ctx.Done()
}
