package model

import "github.com/google/uuid"

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
	Body string `json:"body"`
}
