package model

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	UserExistsError = "user with this username already exists"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Error struct {
	Error string `json:"error"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Created  bool   `json:"created"`
}

type UserContentResponse struct {
	Token string `json:"token"`
}

type Post struct {
	PostID        uuid.UUID      `json:"post_id"`
	PosterID      uuid.UUID      `json:"poster_id"`
	Body          sql.NullString `json:"body"`
	PaywallLocked bool           `json:"paywall_locked"`
	PaywallTier   sql.NullInt64  `json:"paywall_tier,omitempty"`
	ImageRef      sql.NullString `json:"image_ref,omitempty"`
}
