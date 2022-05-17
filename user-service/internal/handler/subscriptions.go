package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

const (
	SubscriptionsPath = "/subscriptions"
	SubscribePath     = "/subscribe"
	UnsubscribePath   = "/unsubscribe"
	CheckSubPath      = "/checksub"
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
	CheckSubscribe(uuid.UUID, uuid.UUID) (bool, error)
}

func (h *Subscriptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDToken, err := getUserID(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
		return
	}
	if r.Method == http.MethodGet {
		subs, err := h.service.GetUserSubscriptionList(userIDToken)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
			return
		}
		writeResponse(w, http.StatusOK, subs)
	}

	if r.Method == http.MethodPost {
		userId, err := jwt.ParseString(userIDToken)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid incoming data")
			writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
			return
		}
		if r.URL.String() == SubscribePath {
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
		if r.URL.String() == CheckSubPath {
			var u model.CheckSubscriptionRequest
			err = json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			subscribed, err := h.service.CheckSubscribe(u.Subscriber, u.Subscribed)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			resp := model.CheckSubscriptionResponse{
				Subscribed: subscribed,
			}
			writeResponse(w, http.StatusOK, resp)
		}
	}
}
