package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	RegisterPath = "/register"
)

type Register struct {
	logger  *zerolog.Logger
	service RegisterService
}

type RegisterService interface {
	Register(string, string) error
}

func NewRegister(logger *zerolog.Logger, srv RegisterService) *Register {
	return &Register{
		logger:  logger,
		service: srv,
	}
}

func (h *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &model.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}

	//TODO: Check password strength
	err = h.service.Register(req.Username, req.Password)
	if err != nil {
		h.logger.Error().Err(err).Msg("Registration error")
		if err.Error() == model.UserExistsError {
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "User with this username already exists"})
			return
		}
		writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
		return
	}

	resp := &model.RegisterResponse{
		Username: req.Username,
		Created:  true,
	}

	writeResponse(w, http.StatusOK, resp)
}
