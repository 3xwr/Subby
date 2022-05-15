package model

import "github.com/google/uuid"

type UserRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarRef string    `json:"avatar_ref"`
}
