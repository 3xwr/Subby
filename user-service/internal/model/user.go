package model

import "github.com/google/uuid"

type UserRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FullInfo bool      `json:"full_info"`
}

type UserIDRequest struct {
	Username string `json:"username"`
}

type UserIDResponse struct {
	Username string    `json:"username"`
	UserID   uuid.UUID `json:"user_id"`
}

type User struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	AvatarRef       string    `json:"avatar_ref"`
	SubscriberCount int64     `json:"subscriber_count,omitempty"`
}
