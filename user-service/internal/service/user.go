package service

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"user-service/internal/model"
)

type User struct {
	logger *zerolog.Logger
	repo   UserRepo
}

func NewUser(logger *zerolog.Logger, repo UserRepo) *User {
	return &User{
		logger: logger,
		repo:   repo,
	}
}

type UserRepo interface {
	GetUserPublicData(userID uuid.UUID) (model.User, error)
}

func (s *User) GetUserPublicData(userID uuid.UUID) (model.User, error) {
	user, err := s.repo.GetUserPublicData(userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
