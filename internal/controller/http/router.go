package http

import (
	"evrone_go_hw_5_1/config"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func Serve(server *HttpServer, cfg *config.Config) {
	router := mux.NewRouter()

	router.HandleFunc("/users", server.Save).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", server.FindByID).Methods(http.MethodGet)
	router.HandleFunc("/users", server.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.DeleteByID).Methods(http.MethodDelete)

	srv := &http.Server{Handler: router, Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)}
	slog.Error(srv.ListenAndServe().Error())
}
