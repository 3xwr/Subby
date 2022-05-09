package model

import (
	"time"

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

type UploadResponse struct {
	FileAddress string `json:"file_address"`
	Uploaded    bool   `json:"uploaded"`
}

type SubSuccessResponse struct {
	Subscriber string `json:"subscriber"`
	Subscribed string `json:"subscribed"`
	SubSuccess bool   `json:"sub_success"`
}

type UnsubSuccessResponse struct {
	Subscriber   string `json:"subscriber"`
	Unsubscribed string `json:"unsubscribed"`
	UnsubSuccess bool   `json:"unsub_success"`
}

type PostSubmitRequest struct {
	Body          *string `json:"body"`
	PaywallLocked bool    `json:"paywall_locked"`
	PaywallTier   *int    `json:"paywall_tier,omitempty"`
	ImageRef      *string `json:"image_ref,omitempty"`
}

type Post struct {
	PostID        uuid.UUID `json:"post_id"`
	PostedAt      time.Time `json:"posted_at"`
	PosterID      uuid.UUID `json:"poster_id"`
	Body          *string   `json:"body"`
	PaywallLocked bool      `json:"paywall_locked"`
	PaywallTier   *int      `json:"paywall_tier,omitempty"`
	ImageRef      *string   `json:"image_ref,omitempty"`
}
