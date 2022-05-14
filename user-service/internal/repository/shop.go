package repository

import (
	"database/sql"
	"user-service/internal/model"

	"github.com/google/uuid"
)

type Shop struct {
	*sql.DB
}

func NewShop(db *sql.DB) *Shop {
	return &Shop{db}
}

func (db *Shop) GetUserItems(OwnerID uuid.UUID) ([]model.ShopItem, error) {
	rows, err := db.Query("SELECT id, owner_id, name, price, description, image_ref FROM shop_items WHERE owner_id=$1", OwnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userItems []model.ShopItem

	for rows.Next() {
		var userItem model.ShopItem
		if err := rows.Scan(&userItem.ID, &userItem.OwnerID, &userItem.Name, &userItem.Price, &userItem.Description, &userItem.ImageRef); err != nil {
			return nil, err
		}
		userItems = append(userItems, userItem)
	}

	return userItems, nil
}
