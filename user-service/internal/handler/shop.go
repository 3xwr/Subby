package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	ShopPath = "/shop"
)

type Shop struct {
	logger  *zerolog.Logger
	service ShopService
}

func NewShop(logger *zerolog.Logger, srv ShopService) *Shop {
	return &Shop{
		logger:  logger,
		service: srv,
	}
}

type ShopService interface {
	GetUserShop(uuid.UUID)
}

func (h *Shop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// userIDToken, err := getUserID(r)
	// if err != nil {
	// 	h.logger.Error().Err(err).Msg("Invalid incoming data")
	// 	writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
	// 	return
	// }
}
