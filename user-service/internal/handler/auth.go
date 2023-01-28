package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"

	"user-service/internal/model"
)

const (
	AuthPath           = "/auth"
	ChangePasswordPath = "/changepassword"
)

type Auth struct {
	logger  *zerolog.Logger
	service AuthService
}

type AuthService interface {
	Authenticate(string, string) (string, string, error)
	ChangePassword(userID uuid.UUID, oldPassword string, newPassword string) error
}

func NewAuth(logger *zerolog.Logger, srv AuthService) *Auth {
	return &Auth{
		logger:  logger,
		service: srv,
	}
}

func (h *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == ChangePasswordPath {
		userIDToken, err := getUserID(r)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
			return
		}
		userID, err := jwt.ParseString(userIDToken)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
			return
		}
		userUUID, err := uuid.Parse(userID.Subject())
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
			return
		}
		var changePasswordRequest model.UserChangePasswordRequest
		err = json.NewDecoder(r.Body).Decode(&changePasswordRequest)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
			return
		}
		err = h.service.ChangePassword(userUUID, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
			return
		}
	} else {
		req := &model.AuthRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
			return
		}

		accessToken, refreshToken, err := h.service.Authenticate(req.Username, req.Password)
		if err != nil {
			h.logger.Error().Err(err).Msg("Authentication error")
			writeResponse(w, http.StatusForbidden, model.Error{Error: "Forbidden"})
			return
		}

		res := &model.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		writeResponse(w, http.StatusOK, res)
	}
}
