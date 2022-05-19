package model

import (
	"time"

	"github.com/google/uuid"
)

type PostSubmitRequest struct {
	Body             *string    `json:"body"`
	MembershipLocked bool       `json:"membership_locked"`
	MembershipTier   *uuid.UUID `json:"membership_tier,omitempty"`
	ImageRef         *string    `json:"image_ref,omitempty"`
}

type PostDeleteRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

type GetUserPostsRequest struct {
	PosterID uuid.UUID `json:"poster_id"`
}

type Post struct {
	PostID           uuid.UUID  `json:"post_id"`
	PostedAt         time.Time  `json:"posted_at"`
	PosterID         uuid.UUID  `json:"poster_id"`
	PosterUsername   string     `json:"poster_username"`
	PosterAvatarRef  string     `json:"poster_avatar"`
	Body             *string    `json:"body"`
	MembershipLocked bool       `json:"membership_locked"`
	MembershipTier   *uuid.UUID `json:"membership_tier,omitempty"`
	ImageRef         *string    `json:"image_ref,omitempty"`
}
