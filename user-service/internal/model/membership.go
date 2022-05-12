package model

import "github.com/google/uuid"

type Membership struct {
	OwnerID uuid.UUID        `json:"owner_id"`
	Tiers   []MembershipTier `json:"tiers"`
}

type MembershipTier struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Price   int64     `json:"price"`
	Rewards string    `json:"rewards"`
}

type CreateMembershipRequest struct {
	Tiers []CreateTierRequest `json:"tiers"`
}

type DeleteTierRequest struct {
	TierID uuid.UUID `json:"tier_id"`
}

type CreateTierRequest struct {
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Rewards string `json:"rewards"`
}

type MembershipRequest struct {
	MembershipID uuid.UUID `json:"membership_id"`
}
