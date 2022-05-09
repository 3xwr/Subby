package service

import (
	"user-service/internal/model"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

type Content struct {
	logger *zerolog.Logger
	repo   ContentRepo
	secret string
}

type ContentRepo interface {
	GetUserSubs(string) ([]string, error)
	GetUserFeed(string, int) ([]model.Post, error)
}

func NewContent(logger *zerolog.Logger, repo ContentRepo, secret string) *Content {
	return &Content{
		logger: logger,
		repo:   repo,
		secret: secret,
	}
}

const (
	defaultPostAmount = 15
)

func (s *Content) GetUserSubscriptionList(token string) ([]string, error) {
	userId, err := jwt.ParseString(token, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	subs, err := s.repo.GetUserSubs(userId.Subject())
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *Content) GetPostsFeedByID(token string) ([]model.Post, error) {
	userId, err := jwt.ParseString(token, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetUserFeed(userId.Subject(), defaultPostAmount)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
