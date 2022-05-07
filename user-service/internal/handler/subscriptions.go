package handler

import (
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	SubscriptionsPath = "/subscriptions"
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
}

func (h *Subscriptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
