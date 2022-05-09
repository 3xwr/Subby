package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

const (
	SubscriptionsPath = "/subscriptions"
	SubscribePath     = "/subscribe"
	UnsubscribePath   = "/unsubscribe"
)

type Subscriptions struct {
	logger  *zerolog.Logger
	service SubscriptionsService
}

func NewSubscriptions(logger *zerolog.Logger, srv SubscriptionsService) *Subscriptions {
	return &Subscriptions{
		logger:  logger,
		service: srv,
	}
}

type SubscriptionsService interface {
	GetUserSubscriptionList(string) ([]string, error)
	SubscribeCurrentUserToUser(string, string) error
	UnsubscribeFromUser(string, string) error
}

func (h *Subscriptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, err := getUserID(r)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
			return
		}
		subs, err := h.service.GetUserSubscriptionList(userID)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
			return
		}
		writeResponse(w, http.StatusOK, subs)
	}

	if r.Method == http.MethodPost {
		if r.URL.String() == SubscribePath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			userId, err := jwt.ParseString(userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}

			var u model.User
			err = json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.SubscribeCurrentUserToUser(userIDToken, u.ID.String())
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			resp := model.SubSuccessResponse{
				Subscriber: userId.Subject(),
				Subscribed: u.ID.String(),
				SubSuccess: true,
			}
			writeResponse(w, http.StatusOK, resp)
		}
		if r.URL.String() == UnsubscribePath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			userId, err := jwt.ParseString(userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}

			var u model.User
			err = json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.UnsubscribeFromUser(userIDToken, u.ID.String())
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			resp := model.UnsubSuccessResponse{
				Subscriber:   userId.Subject(),
				Unsubscribed: u.ID.String(),
				UnsubSuccess: true,
			}
			writeResponse(w, http.StatusOK, resp)
		}
	}
}
