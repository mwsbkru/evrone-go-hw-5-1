package http

import (
	"evrone_go_hw_5_1/config"
	"fmt"
	"log/slog"
	"net/http"
)

// Serve configures routes and runs http-server
func Serve(server *Server, cfg *config.Config) {
	router := http.NewServeMux()

	router.HandleFunc("POST /users", server.Save)
	router.HandleFunc("GET /users/{id}", server.FindByID)
	router.HandleFunc("GET /users", server.FindAll)
	router.HandleFunc("DELETE /users/{id}", server.DeleteByID)

	withMiddlewaresRouter := loggingMiddleware(defaultHeadersMiddleware(router))

	srv := &http.Server{Handler: withMiddlewaresRouter, Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)}

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("srv.ListenAndServe", "err", err)
	}
}

func defaultHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("Начало обработки запроса", slog.String("method", request.Method), slog.String("url", request.RequestURI))
		next.ServeHTTP(writer, request)
	})
}
