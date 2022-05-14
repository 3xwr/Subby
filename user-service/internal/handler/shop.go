package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	ShopPath    = "/shop"
	AddItemPath = "/additem"
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
	GetUserShop(uuid.UUID) ([]model.ShopItem, error)
	AddItem(model.ShopItem) error
}

func (h *Shop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// userIDToken, err := getUserID(r)
	// if err != nil {
	// 	h.logger.Error().Err(err).Msg("Invalid incoming data")
	// 	writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
	// 	return
	// }
	if r.Method == http.MethodPost {
		if r.URL.String() == ShopPath {
			var shopRequest model.GetShopRequest
			err := json.NewDecoder(r.Body).Decode(&shopRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			shopItems, err := h.service.GetUserShop(shopRequest.OwnerID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
			writeResponse(w, http.StatusOK, shopItems)
		}
		if r.URL.String() == AddItemPath {
			var shopItem model.ShopItem
			err := json.NewDecoder(r.Body).Decode(&shopItem)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.AddItem(shopItem)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
		}
	}
}
