package repository

import (
	"database/sql"
	"user-service/internal/model"
)

type Membership struct {
	*sql.DB
}

func NewMembership(db *sql.DB) *Membership {
	return &Membership{db}
}

func (db *Membership) GetMembershipInfo(membershipID string) (model.Membership, error) {
	var membership model.Membership
	err := db.QueryRow("SELECT owner_id FROM memberships WHERE id = $1", membershipID).Scan(&membership.OwnerID)
	if err != nil {
		return membership, err
	}
	rows, err := db.Query("SELECT id, name, price, rewards FROM tiers WHERE membership_id=$1", membershipID)
	if err != nil {
		return membership, err
	}
	defer rows.Close()

	var tiers []model.MembershipTier

	for rows.Next() {
		var tier model.MembershipTier
		if err := rows.Scan(&tier.ID, &tier.Name, &tier.Price, &tier.Rewards); err != nil {
			return membership, err
		}
		tiers = append(tiers, tier)
	}

	membership.Tiers = tiers

	return membership, nil
}
