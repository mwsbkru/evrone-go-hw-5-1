package http

import (
	"evrone_go_hw_5_1/config"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func Serve(server *Server, cfg *config.Config) {
	router := mux.NewRouter()

	router.HandleFunc("/users", server.Save).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", server.FindByID).Methods(http.MethodGet)
	router.HandleFunc("/users", server.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.DeleteByID).Methods(http.MethodDelete)

	router.Use(loggingMiddleware)
	router.Use(defaultHeadersMiddleware)

	srv := &http.Server{Handler: router, Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)}
	slog.Error(srv.ListenAndServe().Error())
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
