package repository

import (
	"context"
	"database/sql"
	"user-service/internal/model"

	"github.com/google/uuid"
)

type Membership struct {
	*sql.DB
}

func NewMembership(db *sql.DB) *Membership {
	return &Membership{db}
}

func (db *Membership) GetMembershipIDByOwnerID(OwnerID uuid.UUID) (uuid.UUID, error) {
	var membershipID uuid.UUID
	err := db.QueryRow("SELECT id FROM memberships WHERE owner_id = $1", OwnerID).Scan(&membershipID)
	if err != nil {
		return membershipID, err
	}
	return membershipID, nil
}

func (db *Membership) GetUserTiers(UserID uuid.UUID) ([]model.UserSubscribedTier, error) {
	rows, err := db.Query("SELECT tier_id FROM members WHERE user_id=$1", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiers []model.UserSubscribedTier

	for rows.Next() {
		var tier model.UserSubscribedTier
		if err := rows.Scan(&tier.TierID); err != nil {
			return nil, err
		}
		tiers = append(tiers, tier)
	}
	return tiers, nil
}

func (db *Membership) GetMembershipInfo(membershipID string) (model.Membership, error) {
	var membership model.Membership
	err := db.QueryRow("SELECT owner_id FROM memberships WHERE id = $1", membershipID).Scan(&membership.OwnerID)
	if err != nil {
		return membership, err
	}
	rows, err := db.Query("SELECT id, name, price, rewards, image_ref FROM tiers WHERE membership_id=$1", membershipID)
	if err != nil {
		return membership, err
	}
	defer rows.Close()

	var tiers []model.MembershipTier

	for rows.Next() {
		var tier model.MembershipTier
		if err := rows.Scan(&tier.ID, &tier.Name, &tier.Price, &tier.Rewards, &tier.ImageRef); err != nil {
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
		tx.Rollback()
		return err
	}

	for _, tier := range membership.Tiers {
		_, err = tx.ExecContext(ctx, "INSERT INTO tiers (id, name, price, rewards, image_ref, membership_id) VALUES ($1, $2, $3, $4, $5, $6)", tier.ID, tier.Name, tier.Price, tier.Rewards, tier.ImageRef, membershipID)
		if err != nil {
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

func (db *Membership) DeleteMembership(ownerID uuid.UUID, membershipID uuid.UUID) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var membership model.Membership
	err = tx.QueryRow("SELECT owner_id FROM memberships WHERE id = $1 AND owner_id = $2", membershipID, ownerID).Scan(&membership.OwnerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM memberships WHERE id=$1 AND owner_id=$2", membershipID, ownerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *Membership) AddTier(tier model.MembershipTier, ownerID uuid.UUID) error {
	var membershipID uuid.UUID
	err := db.QueryRow("SELECT id FROM memberships WHERE owner_id = $1", ownerID).Scan(&membershipID)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO tiers (id, name, price, rewards, image_ref, membership_id) VALUES ($1, $2, $3, $4, $5, $6)", tier.ID, tier.Name, tier.Price, tier.Rewards, tier.ImageRef, membershipID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Membership) SubscribeToMembershipTier(userID uuid.UUID, tierID uuid.UUID) error {
	var u user
	err := db.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&u.id)
	if err != nil {
		return err
	}

	var t model.MembershipTier
	err = db.QueryRow("SELECT id FROM tiers WHERE id = $1", tierID).Scan(&t.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO members (user_id, tier_id) VALUES ($1,$2)", userID, tierID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Membership) DeleteTier(tierID uuid.UUID, ownerID uuid.UUID) error {
	var membershipID uuid.UUID
	err := db.QueryRow("SELECT id FROM memberships WHERE owner_id = $1", ownerID).Scan(&membershipID)
	if err != nil {
		return err
	}
	var tID uuid.UUID
	err = db.QueryRow("SELECT id FROM tiers WHERE id = $1", tierID).Scan(&tID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM tiers WHERE id=$1", tierID)
	if err != nil {
		return err
	}
	return nil
}
