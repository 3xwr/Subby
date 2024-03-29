package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"

	"user-service/internal/model"
)

const (
	authTimeout = time.Minute * 10
)

type Auth struct {
	logger *zerolog.Logger
	repo   AuthRepo
	secret []byte
}

func NewAuth(logger *zerolog.Logger, repo AuthRepo, secret []byte) *Auth {
	return &Auth{
		logger: logger,
		repo:   repo,
		secret: secret,
	}
}

type AuthRepo interface {
	GetUser(string, string) (*model.User, error)
	SaveToken(uuid.UUID, string) error
	SaveUser(string, string, string) error
	ChangePassword(userID uuid.UUID, oldPassword string, newPassword string) error
}

func (s *Auth) Authenticate(username string, password string) (string, string, error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = s.repo.SaveToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Auth) Register(username string, email string, password string) error {
	err := s.repo.SaveUser(username, email, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Auth) ChangePassword(userID uuid.UUID, oldPassword string, newPassword string) error {
	err := s.repo.ChangePassword(userID, oldPassword, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *Auth) generateAccessToken(id uuid.UUID) (string, error) {
	token := jwt.New()
	token.Set(jwt.SubjectKey, id.String())
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(authTimeout).Unix())
	tokenString, err := jwt.Sign(token, jwa.HS256, s.secret)
	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (s *Auth) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
