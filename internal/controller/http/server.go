package http

import (
	"evrone_go_hw_5_1/config"
	"github.com/gorilla/mux"
	"net/http"
)

type HttpServer struct {
	cfg   *config.Config
	greet string
}

func NewHttpServer(cfg *config.Config) *HttpServer {
	return &HttpServer{greet: "Hello", cfg: cfg}
}

func (s *HttpServer) Save(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Save user"))
}

func (s *HttpServer) FindByID(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("find user: " + mux.Vars(request)["id"]))
}

func (s *HttpServer) FindAll(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("find All"))
}

func (s *HttpServer) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("delete user: " + mux.Vars(request)["id"]))
}
