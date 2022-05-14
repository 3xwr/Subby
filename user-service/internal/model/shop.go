package model

import "github.com/google/uuid"

type ShopItem struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	Description string    `json:"rewards"`
	ImageRef    string    `json:"image_ref"`
}

type GetShopRequest struct {
	OwnerID uuid.UUID `json:"owner_id"`
}
