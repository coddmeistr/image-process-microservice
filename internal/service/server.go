package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxik12233/image-process-microservice/pkg/config"
	"github.com/maxik12233/image-process-microservice/pkg/logger"
)

func NewHTTPServer() *http.Server {
	router := mux.NewRouter()

	router = ConfigureRoutes(*router)

	addr := "localhost" + config.GetConfig().BackendPort
	server := &http.Server{
		Handler: router,
		Addr:    addr,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.GetLogger().Info(fmt.Sprintf("Completed server setup. Server gonna listen on: %v", addr))
	return server
}

func ConfigureRoutes(r mux.Router) *mux.Router {
	log := logger.GetLogger()

	src := NewService()
	endps := NewEndpoints(src)

	r.HandleFunc("/upload", logger.LoggerMiddleware(endps.MakeFileUploadEndpoint(src), *log)).Methods("POST")
	r.HandleFunc("/resize", logger.LoggerMiddleware(endps.MakeResizeEndpoint(src), *log)).Methods("POST")

	return &r
}
