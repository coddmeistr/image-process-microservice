package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxik12233/image-process-microservice/pkg/config"
	"github.com/maxik12233/image-process-microservice/pkg/logger"
)

func NewFrontendHTTPServer() *http.Server {
	router := mux.NewRouter()

	router = ConfigureRoutes(*router)

	addr := "localhost" + config.GetConfig().FrontendPort
	server := &http.Server{
		Handler: router,
		Addr:    addr,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.GetLogger().Info(fmt.Sprintf("Completed frontend server. Server gonna listen on: %v", addr))
	return server
}

func ConfigureRoutes(r mux.Router) *mux.Router {
	log := logger.GetLogger()

	r.HandleFunc("/", logger.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/frontend/index.html")
	}, *log)).Methods("GET")

	return &r
}
