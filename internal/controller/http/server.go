package http

import (
	"encoding/json"
	"errors"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/entity/dto"
	"evrone_go_hw_5_1/internal/repo"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Server struct {
	cfg         *config.Config
	userService UserService
}

func NewServer(cfg *config.Config, userService UserService) *Server {
	return &Server{cfg: cfg, userService: userService}
}

func (s *Server) Save(writer http.ResponseWriter, request *http.Request) {
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

	savedUser, err := s.userService.CreateUser(request.Context(), userRequest.Name, userRequest.Email, entity.UserRole(userRequest.Role))
	if err != nil {
		slog.Error("Не удалось сохранить пользователя", slog.String("error", err.Error()), slog.String("email", userRequest.Email))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	userResponseObject := dto.UserResponseBody{
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

func (s *Server) FindByID(writer http.ResponseWriter, request *http.Request) {
	userId := mux.Vars(request)["id"]

	user, err := s.userService.GetUser(request.Context(), userId)
	if err != nil {
		if errors.Is(err, &repo.ErrorUserNotFound{}) {
			s.respondWithError(writer, http.StatusNotFound, err.Error())
			return
		}

		slog.Error("Не загрузить пользователя", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := dto.UserResponseBody{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	userResponseBody, err := json.Marshal(userResponse)
	if err != nil {
		slog.Error("Не удалось сформировать тело ответа", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writer.Write(userResponseBody)
}

func (s *Server) FindAll(writer http.ResponseWriter, request *http.Request) {
	users, err := s.userService.ListUsers(request.Context())
	if err != nil {
		slog.Error("Не удалось загрузить пользователей", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	usersResponseList := make([]dto.UserResponseBody, 0, len(users))
	for i := range users {
		user := dto.UserResponseBody{
			ID:    users[i].ID,
			Name:  users[i].Name,
			Email: users[i].Email,
			Role:  string(users[i].Role),
		}
		usersResponseList = append(usersResponseList, user)
	}

	usersResponse := dto.UsersResponseBody{Data: usersResponseList}

	usersResponseBody, err := json.Marshal(usersResponse)
	if err != nil {
		slog.Error("Не удалось сформировать тело ответа", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writer.Write(usersResponseBody)
}

func (s *Server) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	userId := mux.Vars(request)["id"]

	err := s.userService.RemoveUser(request.Context(), userId)
	if err != nil {
		if errors.Is(err, &repo.ErrorUserNotFound{}) {
			s.respondWithError(writer, http.StatusNotFound, err.Error())
			return
		}

		slog.Error("Произошла ошибка при удалении пользователя", slog.String("error", err.Error()))
		s.respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) respondWithError(writer http.ResponseWriter, code int, message string) {
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
