package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog"
)

const (
	MembershipIDByOwnerIDPath     = "/membershipowner"
	MembershipPath                = "/membership"
	SubscribeToMembershipTierPath = "/tiersubscribe"
	GetUserTiersPath              = "/usertiers"
	CreateMembershipPath          = "/createmembership"
	DeleteMembershipPath          = "/deletemembership"
	AddTierPath                   = "/addtier"
	DeleteTierPath                = "/deletetier"
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
	GetMembershipIDByOwnerID(OwnerID uuid.UUID) (uuid.UUID, error)
	GetUserTiers(UserID uuid.UUID) ([]model.UserSubscribedTier, error)
	GetMembershipInfo(string) (model.Membership, error)
	SubscribeToMembershipTier(uuid.UUID, uuid.UUID) error
	CreateMembership(model.CreateMembershipRequest, string) error
	DeleteMembership(string, uuid.UUID) error
	AddTier(model.CreateTierRequest, string) error
	DeleteTier(uuid.UUID, string) error
}

func (h *Membership) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.String() == MembershipIDByOwnerIDPath {
			var idRequest model.MembershipIDByOwnerIDRequest
			err := json.NewDecoder(r.Body).Decode(&idRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			membershipID, err := h.service.GetMembershipIDByOwnerID(idRequest.OwnerID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
			writeResponse(w, http.StatusOK, model.MembershipIdResponse{MembershipID: membershipID})
		}
		if r.URL.String() == SubscribeToMembershipTierPath {
			var subRequest model.SubscribeToMembershipTierRequest
			err := json.NewDecoder(r.Body).Decode(&subRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.SubscribeToMembershipTier(subRequest.UserID, subRequest.TierID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Subscribe error"})
				return
			}
		}
		if r.URL.String() == AddTierPath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			var mCreateRequest model.CreateTierRequest
			err = json.NewDecoder(r.Body).Decode(&mCreateRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.AddTier(mCreateRequest, userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
		}
		if r.URL.String() == CreateMembershipPath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			var mCreateRequest model.CreateMembershipRequest
			err = json.NewDecoder(r.Body).Decode(&mCreateRequest)
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
		if r.URL.String() == DeleteMembershipPath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			var mRequest model.MembershipRequest
			err = json.NewDecoder(r.Body).Decode(&mRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.DeleteMembership(userIDToken, mRequest.MembershipID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
		}
		if r.URL.String() == DeleteTierPath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			var deleteTierRequest model.DeleteTierRequest
			err = json.NewDecoder(r.Body).Decode(&deleteTierRequest)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err = h.service.DeleteTier(deleteTierRequest.TierID, userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
				return
			}
		}

	}
	if r.Method == http.MethodGet {
		if r.URL.String() == GetUserTiersPath {
			userIDToken, err := getUserID(r)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			jwt, err := jwt.ParseString(userIDToken)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
				return
			}
			userUUID, err := uuid.Parse(jwt.Subject())
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusInternalServerError, model.Error{Error: "UUID parse error"})
				return
			}
			tiers, err := h.service.GetUserTiers(userUUID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Get user tiers error"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if len(tiers) == 0 {
				tiers = make([]model.UserSubscribedTier, 0)
			}
			b, err := json.Marshal(tiers)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(b))
		}
	}
}
