package service

import (
	"time"
	"user-service/internal/model"

	"github.com/google/uuid"
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
	AddUserSubscription(string, string) error
	RemoveUserSubscription(string, string) error
	SaveNewPost(model.Post) error
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

func (s *Content) SubscribeCurrentUserToUser(subbingUserToken string, subscribeToUserId string) error {
	subbingUserId, err := jwt.ParseString(subbingUserToken, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return err
	}
	err = s.repo.AddUserSubscription(subbingUserId.Subject(), subscribeToUserId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Content) UnsubscribeFromUser(unsubbingUserToken string, unsubscribeFromUserId string) error {
	unsubbingUserId, err := jwt.ParseString(unsubbingUserToken, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return err
	}
	err = s.repo.RemoveUserSubscription(unsubbingUserId.Subject(), unsubscribeFromUserId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Content) SubmitPost(userToken string, postData model.PostSubmitRequest) error {
	userId, err := jwt.ParseString(userToken, jwt.WithVerify(jwa.HS256, []byte(s.secret)), jwt.WithValidate(true))
	if err != nil {
		return err
	}
	postID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	posterID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	post := model.Post{
		PostID:        postID,
		PostedAt:      time.Now(),
		PosterID:      posterID,
		Body:          postData.Body,
		PaywallLocked: postData.PaywallLocked,
		PaywallTier:   postData.PaywallTier,
		ImageRef:      postData.ImageRef,
	}
	err = s.repo.SaveNewPost(post)
	if err != nil {
		return err
	}
	return nil
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
