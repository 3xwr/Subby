package model

import (
	"time"

	"github.com/google/uuid"
)

type PostSubmitRequest struct {
	Body             *string `json:"body"`
	MembershipLocked bool    `json:"Membership_locked"`
	MembershipTier   *int    `json:"Membership_tier,omitempty"`
	ImageRef         *string `json:"image_ref,omitempty"`
}

type PostDeleteRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

type Post struct {
	PostID           uuid.UUID `json:"post_id"`
	PostedAt         time.Time `json:"posted_at"`
	PosterID         uuid.UUID `json:"poster_id"`
	PosterUsername   string    `json:"poster_username"`
	Body             *string   `json:"body"`
	MembershipLocked bool      `json:"Membership_locked"`
	MembershipTier   *int      `json:"Membership_tier,omitempty"`
	ImageRef         *string   `json:"image_ref,omitempty"`
}
