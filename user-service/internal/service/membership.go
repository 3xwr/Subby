package service

import (
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

type Membership struct {
	logger *zerolog.Logger
	repo   MembershipRepo
}

type MembershipRepo interface {
	GetMembershipInfo(string) (model.Membership, error)
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
