package handler

import (
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	PostsPath = "/posts"
)

type Posts struct {
	logger  *zerolog.Logger
	service PostsService
}

func NewPosts(logger *zerolog.Logger, srv PostsService) *Posts {
	return &Posts{
		logger:  logger,
		service: srv,
	}
}

type PostsService interface {
	GetUserSubscriptionList(string) ([]string, error)
}

func (h *Posts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add functionality, this is a skeleton
	userID, err := getUserID(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
		return
	}
	_ = userID
	resp := http.Response{}
	writeResponse(w, http.StatusOK, resp)
}
