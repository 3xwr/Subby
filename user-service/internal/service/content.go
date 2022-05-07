package service

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"

	"user-service/internal/model"
)

type Content struct {
	logger *zerolog.Logger
	repo   ContentRepo
	secret string
}

type ContentRepo interface {
	GetUserPosts(token string) ([]model.Post, error)
}

func NewContent(logger *zerolog.Logger, repo ContentRepo, secret string) *Content {
	return &Content{
		logger: logger,
		repo:   repo,
		secret: secret,
	}
}

func (s *Content) GetUserFeed(token string) ([]model.Post, error) {
	//do some ordering??
	userId, err := jwt.ParseString(token, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return []model.Post{}, err
	}
	posts, err := s.repo.GetUserPosts(userId.Subject())
	if err != nil {
		return []model.Post{}, err
	}
	return posts, nil
}
