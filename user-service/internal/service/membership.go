package service

import (
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

type Membership struct {
	logger *zerolog.Logger
	repo   MembershipRepo
}

type MembershipRepo interface {
	GetMembershipInfo(string) (model.Membership, error)
	CreateMembership(model.Membership) error
	DeleteMembership(uuid.UUID, uuid.UUID) error
}

func NewMembership(logger *zerolog.Logger, repo MembershipRepo) *Membership {
	return &Membership{
		logger: logger,
		repo:   repo,
	}
}

func (s *Membership) GetMembershipInfo(membershipID string) (model.Membership, error) {
	membership, err := s.repo.GetMembershipInfo(membershipID)
	if err != nil {
		return model.Membership{}, err
	}
	return membership, nil
}

func (s *Membership) CreateMembership(mCreateRequest model.CreateMembershipRequest, token string) error {
	userId, err := jwt.ParseString(token)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	var tiers []model.MembershipTier
	for _, createRequestTier := range mCreateRequest.Tiers {
		tierID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		tier := model.MembershipTier{
			ID:      tierID,
			Name:    createRequestTier.Name,
			Price:   createRequestTier.Price,
			Rewards: createRequestTier.Rewards,
		}
		tiers = append(tiers, tier)
	}

	membership := model.Membership{
		OwnerID: ownerID,
		Tiers:   tiers,
	}

	err = s.repo.CreateMembership(membership)
	if err != nil {
		return err
	}

	return nil
}

func (s *Membership) DeleteMembership(token string, membershipID uuid.UUID) error {
	userId, err := jwt.ParseString(token)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	err = s.repo.DeleteMembership(ownerID, membershipID)
	if err != nil {
		return err
	}

	return nil
}
