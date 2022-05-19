package model

import "github.com/google/uuid"

type Membership struct {
	OwnerID uuid.UUID        `json:"owner_id"`
	Tiers   []MembershipTier `json:"tiers"`
}

type MembershipTier struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Price    int64     `json:"price"`
	Rewards  string    `json:"rewards"`
	ImageRef *string   `json:"image_ref,omitempty"`
}

type SubscribeToMembershipTierRequest struct {
	UserID uuid.UUID `json:"user_id"`
	TierID uuid.UUID `json:"tier_id"`
}

type CreateMembershipRequest struct {
	Tiers []CreateTierRequest `json:"tiers"`
}

type MembershipIDByOwnerIDRequest struct {
	OwnerID uuid.UUID `json:"owner_id"`
}

type MembershipIdResponse struct {
	MembershipID uuid.UUID `json:"membership_id"`
}

type DeleteTierRequest struct {
	TierID uuid.UUID `json:"tier_id"`
}

type UserSubscribedTier struct {
	TierID uuid.UUID `json:"tier_id"`
}

type CreateTierRequest struct {
	Name     string  `json:"name"`
	Price    int64   `json:"price"`
	Rewards  string  `json:"rewards"`
	ImageRef *string `json:"image_ref,omitempty"`
}

type MembershipRequest struct {
	MembershipID uuid.UUID `json:"membership_id"`
}
