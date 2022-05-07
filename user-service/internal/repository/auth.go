package repository

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"user-service/internal/model"

	"github.com/google/uuid"
)

type user struct {
	id   uuid.UUID
	name string
}

type Auth struct {
	*sql.DB
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{db}
}

func (db *Auth) GetUser(username string, password string) (*model.User, error) {
	var u user
	err := db.QueryRow("SELECT id, username FROM users WHERE username = $1 AND password = $2", username, hash(password)).Scan(&u.id, &u.name)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:   u.id,
		Name: u.name,
	}
	return user, nil
}

func (db *Auth) SaveToken(userID uuid.UUID, token string) error {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO tokens (id, token, user_id) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token=EXCLUDED.token", tokenID, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Auth) SaveUser(username string, password string) error {
	var name string
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if name != "" {
		return fmt.Errorf(model.UserExistsError)
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING", userID, username, hash(password))
	if err != nil {
		return err
	}
	return nil
}

func hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
