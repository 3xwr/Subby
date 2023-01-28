package service

import (
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

type Shop struct {
	logger *zerolog.Logger
	repo   ShopRepo
}

type ShopRepo interface {
	GetUserItems(uuid.UUID) ([]model.ShopItem, error)
	AddItem(model.ShopItem) error
	DeleteItem(uuid.UUID, uuid.UUID) error
}

func NewShop(logger *zerolog.Logger, repo ShopRepo) *Shop {
	return &Shop{
		logger: logger,
		repo:   repo,
	}
}

func (s *Shop) GetUserShop(OwnerID uuid.UUID) ([]model.ShopItem, error) {
	items, err := s.repo.GetUserItems(OwnerID)
	if err != nil {
		return []model.ShopItem{}, err
	}
	return items, nil
}

func (s *Shop) AddItem(itemRequest model.AddItemRequest, userToken string) error {
	itemID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	userId, err := jwt.ParseString(userToken)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	item := model.ShopItem{
		ID:          itemID,
		OwnerID:     ownerID,
		Name:        itemRequest.Name,
		Price:       itemRequest.Price,
		Description: itemRequest.Description,
		ImageRef:    itemRequest.ImageRef,
	}
	err = s.repo.AddItem(item)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shop) DeleteItem(itemID uuid.UUID, userToken string) error {
	userId, err := jwt.ParseString(userToken)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	err = s.repo.DeleteItem(itemID, ownerID)
	if err != nil {
		return err
	}
	return nil
}
