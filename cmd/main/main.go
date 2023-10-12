package main

import (
	"github.com/maxik12233/image-process-microservice/internal/frontend"
	"github.com/maxik12233/image-process-microservice/internal/service"
	"github.com/maxik12233/image-process-microservice/pkg/config"
	"github.com/maxik12233/image-process-microservice/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.Init()
	if err != nil {
		panic(err)
	}

	log.Info("Settuping config...")
	_, err = config.Init()
	if err != nil {
		log.Error("Could't initialize config", zap.Error(err))
		panic(err)
	}

	log.Info("Settuping backend microservice...")
	backendServer := service.NewHTTPServer()

	log.Info("Settuping frontend application...")
	frontendServer := frontend.NewFrontendHTTPServer()

	log.Info("Starting frontend application server...")
	go frontendServer.ListenAndServe()

	log.Info("Starting backend microservice's server...")
	log.Fatal(backendServer.ListenAndServe().Error())
}
