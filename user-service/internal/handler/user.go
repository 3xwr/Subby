package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"user-service/internal/model"
)

const (
	UserDataPath = "/userdata"
)

type User struct {
	logger  *zerolog.Logger
	service UserService
}

type UserService interface {
	GetUserPublicData(userID uuid.UUID) (model.User, error)
}

func NewUser(logger *zerolog.Logger, srv UserService) *User {
	return &User{
		logger:  logger,
		service: srv,
	}
}

func (h *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.String() == UserDataPath {
			var userRequest model.UserRequest
			err := json.NewDecoder(r.Body).Decode(&userRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			user, err := h.service.GetUserPublicData(userRequest.UserID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			writeResponse(w, http.StatusOK, user)
		}
	}
}
