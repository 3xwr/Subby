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
	GetMembershipIDByOwnerID(OwnerID uuid.UUID) (uuid.UUID, error)
	GetUserTiers(UserID uuid.UUID) ([]model.UserSubscribedTier, error)
	GetMembershipInfo(string) (model.Membership, error)
	SubscribeToMembershipTier(uuid.UUID, uuid.UUID) error
	CreateMembership(model.Membership) error
	DeleteMembership(uuid.UUID, uuid.UUID) error
	AddTier(model.MembershipTier, uuid.UUID) error
	DeleteTier(uuid.UUID, uuid.UUID) error
}

func NewMembership(logger *zerolog.Logger, repo MembershipRepo) *Membership {
	return &Membership{
		logger: logger,
		repo:   repo,
	}
}

func (s *Membership) GetUserTiers(UserID uuid.UUID) ([]model.UserSubscribedTier, error) {
	tiers, err := s.repo.GetUserTiers(UserID)
	if err != nil {
		return nil, err
	}
	return tiers, nil
}

func (s *Membership) GetMembershipIDByOwnerID(OwnerID uuid.UUID) (uuid.UUID, error) {
	id, err := s.repo.GetMembershipIDByOwnerID(OwnerID)
	if err != nil {
		return id, err
	}
	return id, nil
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
			ID:       tierID,
			Name:     createRequestTier.Name,
			Price:    createRequestTier.Price,
			Rewards:  createRequestTier.Rewards,
			ImageRef: createRequestTier.ImageRef,
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

func (s *Membership) SubscribeToMembershipTier(userID uuid.UUID, tierID uuid.UUID) error {
	err := s.repo.SubscribeToMembershipTier(userID, tierID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Membership) AddTier(tierCreateRequest model.CreateTierRequest, token string) error {
	userId, err := jwt.ParseString(token)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	TierID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	tier := model.MembershipTier{
		ID:       TierID,
		Name:     tierCreateRequest.Name,
		Price:    tierCreateRequest.Price,
		Rewards:  tierCreateRequest.Rewards,
		ImageRef: tierCreateRequest.ImageRef,
	}
	err = s.repo.AddTier(tier, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Membership) DeleteTier(tierID uuid.UUID, token string) error {
	userId, err := jwt.ParseString(token)
	if err != nil {
		return err
	}
	ownerID, err := uuid.Parse(userId.Subject())
	if err != nil {
		return err
	}
	err = s.repo.DeleteTier(tierID, ownerID)
	if err != nil {
		return err
	}
	return nil
}
