package logger

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func LoggerMiddleware(h http.HandlerFunc, l zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := fmt.Sprintf("Handling %v %v request. Host: %v", r.URL, r.Method, r.Host)
		l.Info(log, zap.Any("RequestBody", r.Body))

		h.ServeHTTP(w, r)
	}
}
