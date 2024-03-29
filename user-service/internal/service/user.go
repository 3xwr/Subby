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
	GetUserIDByName(name string) (uuid.UUID, error)
	GetFullUserPublicData(userID uuid.UUID) (model.User, error)
	GetUserPrivateData(userID uuid.UUID) (model.User, error)
	ChangeUserAvatar(userID uuid.UUID, avatarRef string) error
}

func (s *User) GetUserPublicData(userID uuid.UUID) (model.User, error) {
	user, err := s.repo.GetUserPublicData(userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *User) GetFullUserPublicData(userID uuid.UUID) (model.User, error) {
	user, err := s.repo.GetFullUserPublicData(userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *User) GetUserPrivateData(userID uuid.UUID) (model.User, error) {
	user, err := s.repo.GetUserPrivateData(userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *User) ChangeUserAvatar(userID uuid.UUID, avatarRef string) error {
	err := s.repo.ChangeUserAvatar(userID, avatarRef)
	if err != nil {
		return err
	}
	return nil
}

func (s *User) GetUserID(name string) (uuid.UUID, error) {
	id, err := s.repo.GetUserIDByName(name)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
