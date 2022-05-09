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
	GetPostsFeedByID(string) ([]model.Post, error)
}

func (h *Posts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
		return
	}
	posts, err := h.service.GetPostsFeedByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Unable to get posts")
		writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, posts)
}
