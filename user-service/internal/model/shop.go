package model

import "github.com/google/uuid"

type ShopItem struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	Description string    `json:"description"`
	ImageRef    string    `json:"image_ref"`
}

type AddItemRequest struct {
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	ImageRef    string `json:"image_ref"`
}

type GetShopRequest struct {
	OwnerID uuid.UUID `json:"owner_id"`
}

type DeleteItemRequest struct {
	ItemID uuid.UUID `json:"item_id"`
}
