package service

import (
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Shop struct {
	logger *zerolog.Logger
	repo   ShopRepo
}

type ShopRepo interface {
	GetUserItems(uuid.UUID) ([]model.ShopItem, error)
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
