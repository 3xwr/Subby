package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	MembershipPath       = "/membership"
	CreateMembershipPath = "/createmembership"
)

type Membership struct {
	logger  *zerolog.Logger
	service MembershipService
}

func NewMembership(logger *zerolog.Logger, srv MembershipService) *Membership {
	return &Membership{
		logger:  logger,
		service: srv,
	}
}

type MembershipService interface {
	GetMembershipInfo(string) (model.Membership, error)
	CreateMembership(model.CreateMembershipRequest, string) error
	// DeleteMembership(string, string) error
	// AddTier(model.MembershipTier, string, string) error
	// DeleteTier(string) error
}

func (h *Membership) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDToken, err := getUserID(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
		return
	}
	if r.Method == http.MethodPost {
		if r.URL.String() == CreateMembershipPath {
			var mCreateRequest model.CreateMembershipRequest
			err := json.NewDecoder(r.Body).Decode(&mCreateRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.CreateMembership(mCreateRequest, userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
		}
		if r.URL.String() == MembershipPath {
			var mRequest model.MembershipRequest
			err := json.NewDecoder(r.Body).Decode(&mRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			membership, err := h.service.GetMembershipInfo(mRequest.MembershipID.String())
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
			writeResponse(w, http.StatusOK, membership)
		}

	}
}
