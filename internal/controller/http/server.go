package http

import (
	"encoding/json"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/entity/dto"
	"evrone_go_hw_5_1/internal/usecase"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	cfg         *config.Config
	userService *usecase.UserService
	greet       string
}

func NewHttpServer(cfg *config.Config, userService *usecase.UserService) *HttpServer {
	return &HttpServer{greet: "Hello", cfg: cfg, userService: userService}
}

func (s *HttpServer) Save(writer http.ResponseWriter, request *http.Request) {
	var userRequest dto.SaveUserRequestBody
	err := json.NewDecoder(request.Body).Decode(&userRequest)
	if err != nil {
		slog.Error("Не удалось распарсить тело запроса", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	if !userRequest.IsValid() {
		s.respondWithError(writer, http.StatusBadRequest, "Переданы не верные данные пользователя")
		return
	}

	savedUser, err := s.userService.CreateUser(userRequest.Name, userRequest.Email, entity.UserRole(userRequest.Role))
	if err != nil {
		slog.Error("Не удалось сохранить пользователя", slog.String("error", err.Error()), slog.String("email", userRequest.Email))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	userResponseObject := dto.SaveUserResponseBody{
		ID:    savedUser.ID,
		Name:  savedUser.Name,
		Email: savedUser.Email,
		Role:  string(savedUser.Role),
	}

	userResponseBody, err := json.Marshal(userResponseObject)
	if err != nil {
		slog.Error("Не удалось сформировать тело ответа", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writer.Write(userResponseBody)
}

func (s *HttpServer) FindByID(writer http.ResponseWriter, request *http.Request) {
	s.respondWithError(writer, http.StatusBadRequest, mux.Vars(request)["id"])
	return
}

func (s *HttpServer) FindAll(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("find All"))
}

func (s *HttpServer) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("delete user: " + mux.Vars(request)["id"]))
}

func (s *HttpServer) respondWithError(writer http.ResponseWriter, code int, message string) {
	errorObject := dto.ErrorResponse{
		Code:    code,
		Message: message,
	}

	responseBody, err := json.Marshal(errorObject)
	if err != nil {
		slog.Error("Не удалось сформировать тело ошибки", slog.String("error", err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(code)
	writer.Write(responseBody)
}
