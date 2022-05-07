package service

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

type Subscriptions struct {
	logger *zerolog.Logger
	repo   SubscriptionsRepo
	secret string
}

type SubscriptionsRepo interface {
	GetUserSubs(token string) ([]string, error)
}

func NewSubscriptions(logger *zerolog.Logger, repo SubscriptionsRepo, secret string) *Subscriptions {
	return &Subscriptions{
		logger: logger,
		repo:   repo,
		secret: secret,
	}
}

func (s *Subscriptions) GetUserSubscriptionList(token string) ([]string, error) {
	userId, err := jwt.ParseString(token, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetUserSubs(userId.Subject())
	if err != nil {
		return nil, err
	}
	return posts, nil
}
