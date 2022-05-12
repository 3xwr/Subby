package repository

import (
	"context"
	"database/sql"
	"fmt"
	"user-service/internal/model"

	"github.com/google/uuid"
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

func (db *Membership) CreateMembership(membership model.Membership) error {
	ctx := context.Background()
	membershipID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO memberships (id, owner_id) VALUES ($1, $2)", membershipID, membership.OwnerID)
	if err != nil {
		fmt.Println("memberships insert error")
		tx.Rollback()
		return err
	}

	for _, tier := range membership.Tiers {
		_, err = tx.ExecContext(ctx, "INSERT INTO tiers (id, name, price, rewards, membership_id) VALUES ($1, $2, $3, $4, $5)", tier.ID, tier.Name, tier.Price, tier.Rewards, membershipID)
		if err != nil {
			fmt.Println("tiers insert error")
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
