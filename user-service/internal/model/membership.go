package model

type Membership struct {
	Tiers []MembershipTier `json:"tiers"`
}

type MembershipTier struct {
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Rewards string `json:"rewards"`
}
